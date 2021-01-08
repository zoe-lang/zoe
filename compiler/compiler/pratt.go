package zoe

type prattTk struct {
	lbp int                                            // left binding power
	nud func(ctx Context, tk Tk, lbp int) (Tk, Node)   // when landing on it as a value or prefix
	led func(ctx Context, tk Tk, left Node) (Tk, Node) // when landing on it as an operator
}

// Expression is the standard Pratt parser Expression function
func Expression(ctx Context, tk Tk, rbp int) (Tk, Node) {
	// This is an error case, but has to be handled
	if tk.IsEof() {
		// error ?
		return tk, EmptyNode
	}

	sym_cur := tk.sym()
	tk, left := sym_cur.nud(ctx, tk, rbp)

	// nud might have advanced without us knowing...
	if tk.IsEof() {
		return tk, left
	}

	next_sym := tk.sym()

	for rbp < next_sym.lbp {
		// log.Print(c.Current.KindStr(), c.Current.Value(c.data))
		tk, left = next_sym.led(ctx, tk, left)

		if tk.IsEof() {
			return tk, left
		}

		next_sym = tk.sym()
	}

	return tk, left
}

func literal(tk TokenKind, nk AstNodeKind) {
	s := &syms[tk]
	s.nud = func(ctx Context, tk Tk, lbp int) (Tk, Node) {
		return tk.Next(), tk.createNode(ctx, nk)
	}
}

// binary gives a regular binary operator that will attempt
// to build trees where the root node is the left-most operation node
// happening at its level.
func binary(tk TokenKind, nk AstNodeKind) {
	precedence := lbp
	s := &syms[tk]
	s.lbp = precedence

	s.led = func(lbp int) func(ctx Context, tk Tk, left Node) (Tk, Node) {
		return func(ctx Context, tk Tk, left Node) (Tk, Node) {
			next, right := Expression(ctx, tk.Next(), lbp-1)
			return next, tk.createNode(ctx, nk, left, right) // b.createNodeFromToken(tk, nk, scope, left, right)
		}
	}(precedence)
}

func unary(tk TokenKind, nk AstNodeKind) {
	s := &syms[tk]
	rbp := lbp - 1
	s.nud = func(ctx Context, tk Tk, _ int) (Tk, Node) {
		next, right := Expression(ctx, tk.Next(), rbp)
		return next, tk.createNode(ctx, nk, right) // b.createNodeFromToken(tk, nk, scope, right)
		// return NewNode(nk, tk.Position, c.Expression(rbp))
	}
}

func led(tk TokenKind, fn func(ctx Context, tk Tk, left Node) (Tk, Node)) {
	s := &syms[tk]
	s.lbp = lbp
	s.led = fn
}

func nud(tk TokenKind, fn func(ctx Context, tk Tk, lbp int) (Tk, Node)) {
	s := &syms[tk]
	s.nud = fn
}

func tryParseList(ctx Context, iter Tk, openKind, closeKind, separatorKind TokenKind, mandatorySep bool, fn func(ctx Context, iter Tk) (Tk, Node)) (Tk, Node) {
	// It's always rbp 0
	var list = newList()
	var ok bool

	if iter, ok = iter.expect(openKind); !ok {
		return iter, EmptyNode
	}

	iter = iter.whileNotClosing(func(iter Tk) Tk {
		var ok bool
		var node Node
		var orig = iter

		iter, node = fn(ctx, iter)
		if !node.IsEmpty() {
			list.append(node)
		}

		// If there is a separator, we check for it there

		if iter, ok = iter.consume(separatorKind); !ok && mandatorySep && !iter.IsClosing() {
			iter.reportError(`expected '` + tokstr[separatorKind] + `' but got ` + iter.GetText())
			// ...
			// report an error if it was expected
		}

		if iter.pos == orig.pos {
			// forcibly advance the parser if we didn't get
			iter = iter.Next()
		}

		return iter
	})

	iter, _ = iter.expect(closeKind)
	// check for the closing token

	return iter, list.first
}
