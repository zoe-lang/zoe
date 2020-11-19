package zoe

import (
	"fmt"
)

type prattTk struct {
	lbp int                                                               // left binding power
	nud func(b *nodeBuilder, tk TokenPos, lbp int) NodePosition           // when landing on it as a value or prefix
	led func(b *nodeBuilder, tk TokenPos, left NodePosition) NodePosition // when landing on it as an operator
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

func (b *nodeBuilder) expect(tk TokenKind) TokenPos {
	if b.current >= b.tokensLen {
		b.reportErrorAtToken(b.current-1, `unexpected end of file`)
		return 0
	}
	cur := b.current
	if b.tokens[cur].Kind != tk {
		b.reportErrorAtToken(b.current, fmt.Sprintf(`expected '%s' but got '%s'`, tokstr[tk], b.getTokenText(cur)))
		return 0
	}
	b.advance()
	return cur
}

func (b *nodeBuilder) currentSym() *prattTk {
	if b.current >= b.tokensLen {
		return nil
	}
	cur := &b.tokens[b.current]
	return &syms[cur.Kind]
}

// ExpressionToken parses an expression at the rbp specified by token
func (b *nodeBuilder) ExpressionTokenRbp(tk TokenKind) NodePosition {
	s := syms[tk]
	return b.Expression(s.lbp + 1) // fixme ???
}

func (b *nodeBuilder) advance() {
	b.current++
}

func (b *nodeBuilder) isEof() bool {
	return b.current >= b.tokensLen
}

// Expression is the standard Pratt parser Expression function
func (b *nodeBuilder) Expression(rbp int) NodePosition {
	// This is an error case, but has to be handled
	if b.isEof() {
		// error ?
		return 0
	}

	sym_cur := b.currentSym()
	cur := b.current
	b.advance()
	left := sym_cur.nud(b, cur, rbp)

	// nud might have advanced without us knowing...
	if b.isEof() {
		return left
	}

	cur = b.current
	next_sym := b.currentSym()

	for rbp < next_sym.lbp {
		// log.Print(c.Current.KindStr(), c.Current.Value(c.data))
		b.advance()
		left = next_sym.led(b, cur, left)

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
	s.nud = func(b *nodeBuilder, tk TokenPos, lbp int) NodePosition {
		return b.createNodeFromToken(tk, nk)
	}
}

// binary gives a regular binary operator that will attempt
// to build trees where the root node is the left-most operation node
// happening at its level.
func binary(tk TokenKind, nk AstNodeKind) {
	precedence := lbp
	s := &syms[tk]
	s.lbp = precedence

	s.led = func(lbp int) func(b *nodeBuilder, tk TokenPos, left NodePosition) NodePosition {
		return func(b *nodeBuilder, tk TokenPos, left NodePosition) NodePosition {
			right := b.Expression(lbp - 1)
			return b.createNodeFromToken(tk, nk, left, right)
		}
	}(precedence)
}

func unary(tk TokenKind, nk AstNodeKind) {
	s := &syms[tk]
	rbp := lbp - 1
	s.nud = func(b *nodeBuilder, tk TokenPos, _ int) NodePosition {
		right := b.Expression(rbp)
		return b.createNodeFromToken(tk, nk, right)
		// return NewNode(nk, tk.Position, c.Expression(rbp))
	}
}

func led(tk TokenKind, fn func(b *nodeBuilder, tk TokenPos, left NodePosition) NodePosition) {
	s := &syms[tk]
	s.lbp = lbp
	s.led = fn
}

func nud(tk TokenKind, fn func(b *nodeBuilder, tk TokenPos, lbp int) NodePosition) {
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
