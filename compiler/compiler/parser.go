package zoe

import "log"

var lbp_equal = 0
var lbp_commas = 0
var rbp_arrow = 0
var lbp_gt = 0
var lbp = 2

func init() {
	for i := range syms {
		syms[i].nud = nudError
		syms[i].led = ledError
	}

	// @ in first position, typically after the fn keyword
	// nud(TK_AT, func(c *ZoeContext, tk *Token, rbp int) *Node {
	// 	if c.Peek(TK_LT) { //

	// 	}
	// })

	nud(TK_LBRACKET, func(c *ZoeContext, tk *Token, rbp int) *Node {
		res := make([]*Node, 0)

		for !c.Peek(TK_RBRACKET) {
			if c.isEof() {
				c.reportErrorAtCurrentPosition(`unexpected end of file`)
				break
			}
			res = append(res, c.Expression(0))
		}
		c.Consume(TK_RBRACKET)
		return NewNode(NODE_BLOCK, tk.Position, res...)
	})

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

	nud(KW_VAR, func(c *ZoeContext, tk *Token, rbp int) *Node {
		res := c.Expression(0)
		return normalizeVarDecls(res)
	})

	nud(KW_IMPORT, parseImport)

	lbp += 2

	prefix(NODE_RETURN, KW_RETURN) // FIXME should check for ('}' / 'else' / '|')

	lbp += 2

	nud(KW_TEMPLATE, func(c *ZoeContext, tk *Token, rbp int) *Node {
		args := c.Expression(0).EnsureList()
		// Where clause would come here, most likely
		templated := c.Expression(0)
		if templated.Is(NODE_TYPEDECL) || templated.Is(NODE_FNDECL) {
			exp := templated.Children[1]
			tpl := NewNode(NODE_TEMPLATE, tk.Position, args, exp)
			tpl.UpdatePosition()
			templated.Children[1] = tpl
			templated.UpdatePosition()
			return templated
		}
		c.reportError(tk.Position, `template must precede a 'type' or 'fn' declaration`)
		// this is an error
		return NewNode(NODE_TEMPLATE, tk.Position, args, templated)
	})

	// unary(NODE_TYPE, KW_TYPE)
	nud(KW_TYPE, func(c *ZoeContext, tk *Token, rbp int) *Node {
		name := c.Expect(TK_ID)
		if !c.Consume(KW_IS) {
			c.reportErrorAtCurrentPosition(`expected 'is' after type name`)
		}
		typdef := c.Expression(0)
		return NewNode(NODE_TYPEDECL, tk.Position, NewTerminalNode(name), typdef)
	})

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
	lbp_gt = lbp
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
	surrounding(NODE_LIST, NODE_FNCALL, TK_LPAREN, TK_RPAREN, true)

	lbp += 2

	// surrounding(NODE_BLOCK, NODE_BLOCK, TK_LBRACKET, TK_RBRACKET, false)

	// the index operator
	nud(TK_LBRACE, parseLbraceNud)

	lbp += 2
	rbp_arrow = lbp - 1

	led(TK_ARROW, parseFnSignature)
	led(TK_FATARROW, parseFnFatArrow)
	// binary(NODE_FNDEF, TK_FATARROW)

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

var syms = make([]prattTk, TK__MAX) // Far more than necessary

func nudError(c *ZoeContext, tk *Token, _ int) *Node {
	return c.nodeError(tk)
}

func ledError(c *ZoeContext, tk *Token, left *Node) *Node {
	return c.nodeError(tk, left)
}

///////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////
// Handling of var decl nodes

func normalizeVarDeclNode(node *Node) *Node {
	if node.Is(NODE_COLON) {
		ident := node.Left()
		typ := node.Right()
		if typ.Is(NODE_ASSIGN) {
			return NewNode(NODE_VARDECL, node.Position, ident, typ.Left(), typ.Right())
		}
		return NewNode(NODE_VARDECL, node.Position, ident, typ)
	} else if node.Is(NODE_ASSIGN) {
		// No
		return NewNode(NODE_VARDECL, node.Position, node.Left(), NewNode(NODE_INFER, node.Position), node.Right())
	} else if node.IsIdent() {
		return NewNode(NODE_VARDECL, node.Position, node, NewNode(NODE_INFER, node.Position))
	}

	node.ReportError(`this expression is not a valid declaration`)
	return node
}

func normalizeVarDecls(node *Node) *Node {
	if node.Is(NODE_LIST) {
		vars := make([]*Node, len(node.Children))
		for i, c := range node.Children {
			vars[i] = normalizeVarDeclNode(c)
		}
		return NewNode(NODE_LIST, node.Position, vars...)
	}
	return normalizeVarDeclNode(node)
}

///////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////

///
func parseAfterAt(c *ZoeContext) *Node {
	var right *Node
	if c.Consume(TK_LT) { // opening <
		lst := []*Node{c.Expression(lbp_gt)}
		for c.Consume(TK_COMMA) {
			lst = append(lst, c.Expression(lbp_gt))
		}
		if !c.Consume(TK_GT) {
			c.reportErrorAtCurrentPosition(`expected '>'`)
		}
		if len(lst) > 1 {
			return NewNode(NODE_LIST, lst[0].Position, lst...)
		}
		right = lst[0]
	} else {
		right = c.Expression(0)
	}
	return right
}

// Handle [ ] in nud position, when assigning a function to a variable
func parseLbraceNud(c *ZoeContext, tk *Token, _ int) *Node {
	var xp = c.Expression(0)
	if !c.Consume(TK_RBRACE) {
		c.reportError(c.Current.Position, `expected ']'`)
	}
	return NewNode(NODE_INDEX, tk.Position, xp)
}

// Handle [] as an operator, where it can be
func parseLbraceLed(c *ZoeContext, tk *Token, left *Node) *Node {
	var xp = c.Expression(0)
	if !c.Consume(TK_RBRACE) {
		c.reportError(c.Current.Position, `expected ']'`)
		return NewNode(NODE_LIST, tk.Position, xp)
	}
	return NewNode(NODE_INDEX, tk.Position, left, xp)
}

// Handle import
func parseImport(c *ZoeContext, tk *Token, _ int) *Node {
	if !c.Peek(TK_RAWSTR) {
		c.reportErrorAtCurrentPosition(`import expects a raw string as the module name`)
	}
	name := NewTerminalNode(c.Current)
	c.advance()
	if c.Consume(KW_AS) {
		exp := c.Expression(0)
		return NewNode(NODE_IMPORT, tk.Position, name, exp)
	}
	exp := c.Expression(0)
	if !exp.Is(NODE_LIST) {
		log.Print(exp.Kind, " - ", exp.String(), exp.Token.Debug(), " - ")
		return NewNode(NODE_IMPORT, tk.Position, name, NewNode(NODE_LIST, tk.Position, exp))
	}
	return NewNode(NODE_IMPORT, tk.Position, name, exp)
}

///////////////////////////////////////////////////////
// "
func parseQuote(c *ZoeContext, tk *Token, _ int) *Node {
	// this should transform the result to a string
	return parseUntil(c, NODE_STR, tk, TK_QUOTE, 0)
}

////////////////////////////////////////////////////////
// ->
// This is the signature operator
// It may be followed by a definition with => (that it handles itself)
// Or by a { which will return a block
func parseFnSignature(c *ZoeContext, tk *Token, left *Node) *Node {
	// left contains the list parenthesis

	right := c.Expression(0)
	res := NewNode(NODE_SIGNATURE, tk.Position, left.EnsureList(), right)

	if c.Peek(TK_LBRACKET) {
		// if we are followed by a {, it means this function has a body
		blk := c.Expression(0)
		return NewNode(NODE_FNDEF, tk.Position, res, blk)
	}

	return res
}

func parseFnFatArrow(c *ZoeContext, tk *Token, left *Node) *Node {
	// left is a list of arguments
	// right of => is the implementation of the function

	impl := c.Expression(0) // it is a block or a single expression

	if !impl.Is(NODE_BLOCK) {
		impl = WrapNode(NODE_BLOCK, impl)
	}

	if !left.Is(NODE_SIGNATURE) {
		// TODO: check that what's incoming can be a signature
		if left.Is(":") || left.Is(NODE_TERMINAL) || left.Is(NODE_LIST) {
			return NewNode(
				NODE_FNDEF,
				tk.Position,
				NewNode(NODE_SIGNATURE, left.Position,
					left.EnsureList(),
					NewNode(NODE_INFER, left.Position),
				),
				impl,
			)
		}
		c.reportError(tk.Position, `unexpected '=>', found `, string(left.Kind))
		// if !left.Is(":") && !left.IsToken(TK_ID) && !left.Is
		return NewNode(NODE_ERROR, tk.Position, left, impl)
	}

	// at this stage, we have a node signature and a block, so we just report it a
	// function definition
	return NewNode(NODE_FNDEF, tk.Position, left, impl)
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
			// this feels really hacky
			c.Current = idtk
			return NewNode(NODE_FNDECL, tk.Position, c.Expression(0))
		}
		return NewNode(NODE_FNDECL, tk.Position, id, c.Expression(0))
	}

	return c.Expression(0) //NewNode(NODE_FNDEF, tk.Position, c.Expression(0))

}
