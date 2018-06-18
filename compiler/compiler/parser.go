package zoe

import (
	"fmt"
)

var lbp_equal = 0
var lbp_commas = 0
var rbp_arrow = 0
var lbp = 2

func init() {
	for i := range syms {
		syms[i].nud = nudError
		syms[i].led = ledError
	}

	// The doc comment forwards the results but sets itself first on the node that resulted
	// Doc comments whose next meaningful token are other doc comments or the end of the file
	// are added at the module level
	nud(TK_DOCCOMMENT, func(c *ZoeContext, tk *Token, rbp int) *Node {
		next := c.Expression(rbp)
		nt := tk.NextMeaningfulToken()
		if nt == nil || nt.Is(TK_DOCCOMMENT) {
			c.RootDocComments = append(c.RootDocComments, tk)
		} else {
			c.DocCommentMap[next] = tk
		}
		return next
	})

	nud(KW_FOR, parseFor)
	nud(KW_IF, parseIf)
	unary(NODE_VAR, KW_VAR)
	nud(KW_IMPORT, parseImport)

	lbp += 2

	prefix(NODE_RETURN, KW_RETURN) // FIXME should check for ('}' / 'else' / '|')

	lbp += 2

	unary(NODE_TYPE, KW_TYPE)
	list(TK_SEMICOLON, NODE_BLOCK, false) // cannot start with a semicolon
	list(TK_COMMA, NODE_LIST, false, TK_RPAREN, TK_RBRACKET)

	lbp += 2
	lbp_equal = lbp

	binary(NODE_ASSIGN, TK_EQ) // what about the precedence
	binary(NODE_IS, KW_IS)

	// fn eats up the expression right next to it
	nud(KW_FN, parseFn)

	lbp += 2

	binary(NODE_COLON, TK_COLON)
	unary(NODE_LOCAL, KW_LOCAL)
	unary(NODE_CONST, KW_CONST)

	lbp += 2

	binary(NODE_LT, TK_LT)
	binary(NODE_GT, TK_GT)
	unary(NODE_ELLIPSIS, TK_ELLIPSIS)

	lbp += 2

	list(TK_PIPE, NODE_UNION, true)

	lbp += 2

	prefix(NODE_PLUS, TK_PLUS)
	prefix(NODE_MIN, TK_MIN)
	binary(NODE_MIN, TK_MIN)
	binary(NODE_PLUS, TK_PLUS)

	lbp += 2

	binary(NODE_MUL, TK_STAR)
	binary(NODE_DIV, TK_DIV)

	lbp += 2

	binary(NODE_IS, KW_IS)

	lbp += 2

	// When used right next to an expression, then paren is a function call
	surrounding(NODE_LIST, NODE_FNCALL, TK_LPAREN, TK_RPAREN)

	lbp += 2

	surrounding(NODE_BLOCK, NODE_BLOCK, TK_LBRACKET, TK_RBRACKET)
	surrounding(NODE_LIST, NODE_INDEX, TK_LBRACE, TK_RBRACE)

	lbp += 2
	rbp_arrow = lbp - 1

	led(TK_ARROW, parseArrow)
	binary(NODE_FNDEF, TK_FATARROW)

	lbp += 2

	binary(NODE_NAMESPACE_ACCESS, TK_COLCOL)
	binary(NODE_BINOP_DOT, TK_DOT)
	binary(NODE_AT, TK_AT)

	binary(NODE_AS, KW_AS)

	lbp += 2
	// all the terminals. Lbp was raised, but this is not necessary

	nud(TK_QUOTE, parseQuote)

	terminal(
		TK_STAR, // * when used as a nud is the dereference operator as in something.*
		TK_NUMBER,
		TK_INSTR,
		TK_ID,
		TK_RAWSTR,
	)

}

///////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////

var syms = make([]prattTk, 256) // Far more than necessary

func nudError(c *ZoeContext, tk *Token, _ int) *Node {
	return c.nodeError(tk)
}

func ledError(c *ZoeContext, tk *Token, left *Node) *Node {
	return c.nodeError(tk, left)
}

///////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////

// Handle import
func parseImport(c *ZoeContext, tk *Token, _ int) *Node {
	if !c.Peek(TK_RAWSTR) {
		c.reportErrorAtCurrentPosition(`import expects a raw string as the module name`)
	}
	name := c.Current
	c.advance()
	if c.Consume(KW_AS) {
		exp := c.Expression(0)
		return NewNode(NODE_IMPORT, tk.Position, NewTerminalNode(name), exp)
	}
	exp := c.Expression(0)
	if !exp.Is(NODE_LIST) {

	}
	return NewNode(NODE_IMPORT, tk.Position, NewTerminalNode(name), exp)
}

///////////////////////////////////////////////////////
// "
func parseQuote(c *ZoeContext, tk *Token, _ int) *Node {
	// this should transform the result to a string
	return parseUntil(c, NODE_STR, tk, TK_QUOTE, 0)
}

////////////////////////////////////////////////////////
// ->
func parseArrow(c *ZoeContext, tk *Token, left *Node) *Node {
	// left contains the list parenthesis
	res := NewNode(NODE_SIGNATURE, tk.Position, left, c.Expression(rbp_arrow))
	if c.Peek(TK_LBRACKET) {
		blk := c.Expression(0)
		return NewNode(NODE_FNDEF, tk.Position, res, blk)
	}
	return res
}

/////////////////////////////////////////////////////
// FOR block
func parseFor(c *ZoeContext, tk *Token, _ int) *Node {
	first := c.Expression(0)
	exp := c.Expression(0)
	return NewNode(NODE_FOR, tk.Position, first, exp)
	// return New
}

/////////////////////////////////////////////////////
// Special handling for if block
func parseIf(c *ZoeContext, tk *Token, _ int) *Node {
	cond := c.Expression(0)
	then := c.Expression(0)
	if c.Consume(KW_ELSE) {
		els := c.Expression(0)
		return NewNode(NODE_IF, tk.Position, cond, then, els)
	}
	return NewNode(NODE_IF, tk.Position, cond, then)
}

/////////////////////////////////////////////////////
// Special handling for fn
func parseFn(c *ZoeContext, tk *Token, _ int) *Node {
	if c.Peek(TK_ID) {
		idtk := c.Current
		id := NewTerminalNode(idtk)
		c.advance()
		if c.Peek(TK_FATARROW) {
			// fn a => ..., we have to reset the parser
			c.Current = idtk
			return NewNode(NODE_FNDEF, tk.Position, c.Expression(0))
		}
		return NewNode(NODE_FNDEF, tk.Position, id, c.Expression(0))
	}
	return c.Expression(0) //NewNode(NODE_FNDEF, tk.Position, c.Expression(0))
}

//////////////////////// KWPARSER //////////////////////////
// The kwParser is just a tool to help us parse the beginning of a definition
// since the keywords could appear in a funny order or be used incorrectly,
// which it should check

var kwDeclaration = newKwParser(KW_VAR, KW_FN, KW_LOCAL, KW_TYPE)

func newKwParser(kws ...TokenKind) func() *kwParser {
	allowed := make(map[TokenKind]struct{})
	for _, k := range kws {
		allowed[k] = struct{}{}
	}
	return func() *kwParser {
		return &kwParser{
			mp:      make(map[TokenKind]*Token),
			allowed: allowed,
			dups:    make([]*Token, 0),
		}
	}
}

type kwParser struct {
	mp      map[TokenKind]*Token
	allowed map[TokenKind]struct{}
	dups    []*Token
}

func (k *kwParser) isAllowed(tk *Token) bool {
	if _, ok := k.allowed[tk.Kind]; ok {
		return true
	}
	return false
}

func (k *kwParser) addToken(c *ZoeContext, tk *Token) bool {
	if k.isAllowed(tk) {
		if _, ok := k.mp[tk.Kind]; ok {
			k.dups = append(k.dups, tk)
			c.reportError(tk.Position, fmt.Sprintf(`unnecessary '%s'`, tk.String()))
		} else {
			k.mp[tk.Kind] = tk
		}
		return true
	}
	return false
}

func (k *kwParser) parse(c *ZoeContext) {
	for c.Current != nil {
		if !k.addToken(c, c.Current) {
			break
		}
		c.advance()
	}
}

func (k *kwParser) has(tk TokenKind) bool {
	if _, ok := k.mp[tk]; ok {
		return true
	}
	return false
}
