package zoe

type prattTk struct {
	lbp int                                            // left binding power
	nud func(scope Scope, tk Tk, lbp int) (Tk, Node)   // when landing on it as a value or prefix
	led func(scope Scope, tk Tk, left Node) (Tk, Node) // when landing on it as an operator
}

// Expression is the standard Pratt parser Expression function
func Expression(scope Scope, tk Tk, rbp int) (Tk, Node) {
	// This is an error case, but has to be handled
	if tk.IsEof() {
		// error ?
		return tk, EmptyNode
	}

	sym_cur := tk.sym()
	tk, left := sym_cur.nud(scope, tk, rbp)

	// nud might have advanced without us knowing...
	if tk.IsEof() {
		return tk, left
	}

	next_sym := tk.sym()

	for rbp < next_sym.lbp {
		// log.Print(c.Current.KindStr(), c.Current.Value(c.data))
		tk, left = next_sym.led(scope, tk, left)

		if tk.IsEof() {
			return tk, left
		}

		next_sym = tk.sym()
	}

	return tk, left
}

func literal(tk TokenKind, nk AstNodeKind) {
	s := &syms[tk]
	s.nud = func(scope Scope, tk Tk, lbp int) (Tk, Node) {
		return tk.Next(), tk.createNode(scope, nk)
	}
}

// binary gives a regular binary operator that will attempt
// to build trees where the root node is the left-most operation node
// happening at its level.
func binary(tk TokenKind, nk AstNodeKind) {
	precedence := lbp
	s := &syms[tk]
	s.lbp = precedence

	s.led = func(lbp int) func(scope Scope, tk Tk, left Node) (Tk, Node) {
		return func(scope Scope, tk Tk, left Node) (Tk, Node) {
			next, right := Expression(scope, tk.Next(), lbp-1)
			return next, tk.createNode(scope, nk, left, right) // b.createNodeFromToken(tk, nk, scope, left, right)
		}
	}(precedence)
}

func unary(tk TokenKind, nk AstNodeKind) {
	s := &syms[tk]
	rbp := lbp - 1
	s.nud = func(scope Scope, tk Tk, _ int) (Tk, Node) {
		next, right := Expression(scope, tk.Next(), rbp)
		return next, tk.createNode(scope, nk, right) // b.createNodeFromToken(tk, nk, scope, right)
		// return NewNode(nk, tk.Position, c.Expression(rbp))
	}
}

func led(tk TokenKind, fn func(scope Scope, tk Tk, left Node) (Tk, Node)) {
	s := &syms[tk]
	s.lbp = lbp
	s.led = fn
}

func nud(tk TokenKind, fn func(scope Scope, tk Tk, lbp int) (Tk, Node)) {
	s := &syms[tk]
	s.nud = fn
}
