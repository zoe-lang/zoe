package zoe

import (
	"fmt"
)

type prattTk struct {
	lbp int                                            // left binding power
	nud func(c *ZoeContext, tk *Token, lbp int) Node   // when landing on it as a value or prefix
	led func(c *ZoeContext, tk *Token, left Node) Node // when landing on it as an operator
}

func (c *ZoeContext) Peek(tk ...TokenKind) bool {
	if c.Current == nil {
		return false
	}
	kind := c.Current.Kind
	for _, t := range tk {
		if kind == t {
			return true
		}
	}
	return false
}

func (c *ZoeContext) Consume(tk TokenKind) bool {
	if c.Peek(tk) {
		c.advance()
		return true
	}
	return false
}

func (c *ZoeContext) Expect(tk TokenKind) *Token {
	if c.Current.Kind != tk {
		c.reportError(c.Current.Position, fmt.Sprintf(`expected '%s' but got '%s'`, tokstr[tk], c.Current.String()))
		return nil
	}
	res := c.Current
	c.advance()
	return res
}

// Expression is the standard Pratt parser Expression function
func (c *ZoeContext) Expression(rbp int) Node {
	// This is an error case, but has to be handled
	if c.isEof() {
		return c.End.CreateEof()
	}

	t, sym_cur := c.currentSym()
	c.advance()
	left := sym_cur.nud(c, t, rbp)

	// nud might have advanced without us knowing...
	if c.isEof() {
		return left
	}

	next, next_sym := c.currentSym()
	// c.debugtoken(next)

	for rbp < next_sym.lbp {
		// log.Print(c.Current.KindStr(), c.Current.Value(c.data))
		c.advance()
		left = next_sym.led(c, next, left)

		if c.isEof() {
			return left
		}

		next, next_sym = c.currentSym()
	}

	return left
}

// binary gives a regular binary operator that will attempt
// to build trees where the root node is the left-most operation node
// happening at its level.
func binary(tk TokenKind) {
	precedence := lbp
	s := &syms[tk]
	s.lbp = precedence

	s.led = func(lbp int) func(c *ZoeContext, tk *Token, left Node) Node {
		return func(c *ZoeContext, tk *Token, left Node) Node {
			res := tk.CreateOperation().AddOperands(left, c.Expression(lbp-1))
			res.TokenKind = tk.Kind
			return res
		}
	}(precedence)
}

func unary(tk TokenKind) {
	s := &syms[tk]
	rbp := lbp - 1
	s.nud = func(c *ZoeContext, tk *Token, _ int) Node {
		res := tk.CreateOperation().AddOperands(c.Expression(rbp))
		res.TokenKind = tk.Kind
		return res
		// return NewNode(nk, tk.Position, c.Expression(rbp))
	}
}

func led(tk TokenKind, fn func(c *ZoeContext, tk *Token, left Node) Node) {
	s := &syms[tk]
	s.lbp = lbp
	s.led = fn
}

func nud(tk TokenKind, fn func(c *ZoeContext, tk *Token, lbp int) Node) {
	s := &syms[tk]
	s.nud = fn
}

func terminal(tk TokenKind, create func(tk *Token) Node) {
	nud(tk, func(c *ZoeContext, tk *Token, _ int) Node { return create(tk) })
}

// // parseUntil calls expression several times until landing on a token
// func parseUntil(c *ZoeContext, nk NodeKind, lst *Token, until TokenKind, rbp int) Node {
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
// func parseList(c *ZoeContext, lst *Token, separator TokenKind, terminator TokenKind, produce func(c *ZoeContext) Node) Node {
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
