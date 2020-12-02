package zoe

import (
	"fmt"
)

type prattTk struct {
	lbp int                                                                            // left binding power
	nud func(b *nodeBuilder, scope Scope, tk TokenPos, lbp int) NodePosition           // when landing on it as a value or prefix
	led func(b *nodeBuilder, scope Scope, tk TokenPos, left NodePosition) NodePosition // when landing on it as an operator
}

var closingTokens = [TK__MAX]bool{}

func init() {
	for i := range closingTokens {
		closingTokens[i] = true
	}
	closingTokens[TK_RBRACE] = false
	closingTokens[TK_RPAREN] = false
	closingTokens[TK_RBRACKET] = false
}

// asLongAsNotClosingToken replies true as long as the current token is not
// a closing token such as ')', '}' or ']'
// This is to be in for loops when inside of balanced expressions.
// When we have for instance an opening '(', we *still* want to break the loop
// if the current token is ']', because that means there is a syntax error and
// the ']' is most likely trying to close a balancing [ started before the parenthesized
// group and the ')' has been omitted by the user. In effect, this is as if we "inserted"
// the closing token, which will be reporting as missing by an expect right after
// the loop.
func (b *nodeBuilder) asLongAsNotClosingToken() bool {
	cur := b.current
	if cur >= b.tokensLen {
		return false
	}
	return closingTokens[b.tokens[cur].Kind]
}

func (b *nodeBuilder) currentTokenIs(tk ...TokenKind) bool {
	if b.current >= b.tokensLen {
		return false
	}
	kind := b.tokens[b.current].Kind
	for _, t := range tk {
		if kind == t {
			return true
		}
	}

	return false
}

func (b *nodeBuilder) consume(tk TokenKind) TokenPos {
	if b.currentTokenIs(tk) {
		cur := b.current
		b.advance()
		return cur
	}
	return 0
}

func (b *nodeBuilder) expectNoAdvance(tk TokenKind) TokenPos {
	if b.current >= b.tokensLen {
		b.reportErrorAtToken(b.current-1, `unexpected end of file`)
		return 0
	}
	cur := b.current
	if b.tokens[cur].Kind != tk {
		b.reportErrorAtToken(b.current, fmt.Sprintf(`expected '%s' but got '%s'`, tokstr[tk], b.getTokenText(cur)))
		return 0
	}
	return cur
}

func (b *nodeBuilder) expect(tk TokenKind) TokenPos {
	tok := b.expectNoAdvance(tk)
	if tok != 0 {
		b.advance()
	}
	return tok
}

func (b *nodeBuilder) currentSym() *prattTk {
	if b.current >= b.tokensLen {
		return nil
	}
	cur := &b.tokens[b.current]
	return &syms[cur.Kind]
}

func (b *nodeBuilder) advance() {
	b.current++
}

func (b *nodeBuilder) isEof() bool {
	return b.current >= b.tokensLen
}

// Expression is the standard Pratt parser Expression function
func (b *nodeBuilder) Expression(scope Scope, rbp int) NodePosition {
	// This is an error case, but has to be handled
	if b.isEof() {
		// error ?
		return 0
	}

	sym_cur := b.currentSym()
	cur := b.current
	b.advance()
	left := sym_cur.nud(b, scope, cur, rbp)

	// nud might have advanced without us knowing...
	if b.isEof() {
		return left
	}

	cur = b.current
	next_sym := b.currentSym()

	for rbp < next_sym.lbp {
		// log.Print(c.Current.KindStr(), c.Current.Value(c.data))
		b.advance()
		left = next_sym.led(b, scope, cur, left)

		if b.isEof() {
			return left
		}

		cur = b.current
		next_sym = b.currentSym()
	}

	return left
}

func literal(tk TokenKind, nk AstNodeKind) {
	s := &syms[tk]
	s.nud = func(b *nodeBuilder, scope Scope, tk TokenPos, lbp int) NodePosition {
		return b.createNodeFromToken(tk, nk, scope)
	}
}

// binary gives a regular binary operator that will attempt
// to build trees where the root node is the left-most operation node
// happening at its level.
func binary(tk TokenKind, nk AstNodeKind) {
	precedence := lbp
	s := &syms[tk]
	s.lbp = precedence

	s.led = func(lbp int) func(b *nodeBuilder, scope Scope, tk TokenPos, left NodePosition) NodePosition {
		return func(b *nodeBuilder, scope Scope, tk TokenPos, left NodePosition) NodePosition {
			right := b.Expression(scope, lbp-1)
			return b.createNodeFromToken(tk, nk, scope, left, right)
		}
	}(precedence)
}

func unary(tk TokenKind, nk AstNodeKind) {
	s := &syms[tk]
	rbp := lbp - 1
	s.nud = func(b *nodeBuilder, scope Scope, tk TokenPos, _ int) NodePosition {
		right := b.Expression(scope, rbp)
		return b.createNodeFromToken(tk, nk, scope, right)
		// return NewNode(nk, tk.Position, c.Expression(rbp))
	}
}

func led(tk TokenKind, fn func(b *nodeBuilder, scope Scope, tk TokenPos, left NodePosition) NodePosition) {
	s := &syms[tk]
	s.lbp = lbp
	s.led = fn
}

func nud(tk TokenKind, fn func(b *nodeBuilder, scope Scope, tk TokenPos, lbp int) NodePosition) {
	s := &syms[tk]
	s.nud = fn
}

// // parseUntil calls expression several times until landing on a token
// func parseUntil(c *File, nk NodeKind, lst *Token, until TokenKind, rbp int) Node {
// 	res := make([]Node, 0)
// 	iter := c.Current
// 	for iter != nil {
// 		if iter.Kind == until {
// 			c.advance()
// 			break
// 		}
// 		res = append(res, c.Expression(rbp))
// 		iter = c.Current
// 	}
// 	// check that iter is nil for potential error ?
// 	return NewNode(nk, lst.Position, res...)
// }

// parse a terminated corresponding

// parse a list
// func parseList(c *File, lst *Token, separator TokenKind, terminator TokenKind, produce func(c *File) Node) Node {
// 	res := make([]Node, 0)
// 	iter := c.Current
// 	for iter != nil {
// 		if iter.Kind == terminator {
// 			c.advance()
// 			break
// 		}
// 		res = append(res, produce(c))
// 		if c.Peek(separator) {
// 			c.advance()
// 		}
// 		iter = c.Current
// 	}
// 	return NewListNode(lst, res...)
// }
