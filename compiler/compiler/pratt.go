package zoe

import (
	"fmt"
)

type prattTk struct {
	lbp int                                              // left binding power
	nud func(c *ZoeContext, tk *Token, rbp int) *Node    // when landing on it as a value or prefix
	led func(c *ZoeContext, tk *Token, left *Node) *Node // when landing on it as an operator
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

func (c *ZoeContext) EOF() *Node {
	tk := &Token{Kind: TK_EOF}
	return NewErrorNode(tk)
}

func (c *ZoeContext) Expect(tk TokenKind) *Token {
	if c.Current.Kind != tk {
		c.reportError(c.Current.Position, fmt.Sprintf(`unexpected '%s'`, c.Current.String()))
		return nil
	}
	res := c.Current
	c.advance()
	return res
}

// Expression is the standard Pratt parser Expression function
func (c *ZoeContext) Expression(rbp int) *Node {
	// This is an error case, but has to be handled
	if c.isEof() {
		return c.EOF()
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

func prefix(nk NodeKind, tk TokenKind) {
	s := &syms[tk]
	rbp := lbp - 1
	s.nud = func(c *ZoeContext, tk *Token, _ int) *Node {
		return NewNode(nk, tk.Position, c.Expression(rbp))
	}
}

// binary gives a regular binary operator that will attempt
// to build trees where the root node is the left-most operation node
// happening at its level.
func binary(nk NodeKind, tk TokenKind) {
	precedence := lbp
	s := &syms[tk]
	s.lbp = precedence
	s.led = func(c *ZoeContext, tk *Token, left *Node) *Node {
		// log.Print(c.Current.Debug(c.data))
		return NewNode(nk, tk.Position, left, c.Expression(precedence-1))
	}
}

func unary(nk NodeKind, tk TokenKind) {
	s := &syms[tk]
	rbp := lbp - 1
	s.nud = func(c *ZoeContext, tk *Token, _ int) *Node {
		return NewNode(nk, tk.Position, c.Expression(rbp))
	}
}

func led(tk TokenKind, fn func(c *ZoeContext, tk *Token, left *Node) *Node) {
	s := &syms[tk]
	s.lbp = lbp
	s.led = fn
}

func nud(tk TokenKind, fn func(c *ZoeContext, tk *Token, rbp int) *Node) {
	s := &syms[tk]
	s.nud = fn
}

func terminal(tk ...TokenKind) {
	for _, t := range tk {
		nud(t, func(c *ZoeContext, tk *Token, _ int) *Node { return NewTerminalNode(tk) })
	}
}

// Group expression between two surrounding tokens
// nk : the node of the list returned when nud
// lednk : the node of the list returned when led
// opening : the opening token
// closing : the closing token
// reduce : whether it is allowed to reduce the list to the single expression
func surrounding(nk NodeKind, lednk NodeKind, opening TokenKind, closing TokenKind, reduce bool) {
	s := &syms[opening]
	s.lbp = lbp
	s.led = func(c *ZoeContext, tk *Token, left *Node) *Node {
		contents := make([]*Node, 0)
		for !c.Consume(closing) {
			if c.isEof() {
				contents = append(contents, c.EOF())
				break
			}
			contents = append(contents, c.Expression(0))
		}
		if reduce && len(contents) == 1 && !c.Peek(TK_ARROW, TK_FATARROW) {
			return contents[0]
		}
		return NewNode(lednk, tk.Position, left, NewNode(nk, tk.Position, contents...))
	}

	s.nud = func(c *ZoeContext, tk *Token, _ int) *Node {
		contents := make([]*Node, 0)
		for !c.Consume(closing) {
			if c.isEof() {
				contents = append(contents, c.EOF())
				break
			}
			contents = append(contents, c.Expression(0))
		}
		if reduce && len(contents) == 1 && !c.Peek(TK_ARROW, TK_FATARROW) {
			return contents[0]
		}
		return NewNode(nk, tk.Position, contents...)
	}
}

func list(kind TokenKind, nk NodeKind, allowLeading bool, trailing ...TokenKind) {
	rbp := lbp - 1
	s := &syms[kind]
	trailings := make([]bool, 256)
	for _, t := range trailing {
		trailings[t] = true
	}

	s.lbp = lbp
	s.nud = func(c *ZoeContext, tk *Token, _ int) *Node {
		if !allowLeading {
			c.reportError(tk.Position, `'`, tk.String(), `' cannot come first in an expression`)
		}
		return c.Expression(rbp)
	}

	s.led = func(c *ZoeContext, tk *Token, left *Node) *Node {
		if c.Current != nil && trailings[c.Current.Kind] {
			return left
		}
		next := c.Expression(rbp)
		if next.Kind == nk {
			next.Token = tk
			next.Children = append([]*Node{left}, next.Children...)
			return next
		}
		return NewNode(nk, tk.Position, left, next)
	}
}

// parseUntil calls expression several times until landing on a token
func parseUntil(c *ZoeContext, nk NodeKind, lst *Token, until TokenKind, rbp int) *Node {
	res := make([]*Node, 0)
	iter := c.Current
	for iter != nil {
		if iter.Kind == until {
			c.advance()
			break
		}
		res = append(res, c.Expression(rbp))
		iter = c.Current
	}
	// check that iter is nil for potential error ?
	return NewNode(nk, lst.Position, res...)
}

// parse a terminated corresponding

// parse a list
// func parseList(c *ZoeContext, lst *Token, separator TokenKind, terminator TokenKind, produce func(c *ZoeContext) *Node) *Node {
// 	res := make([]*Node, 0)
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
