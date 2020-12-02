package zoe

func (f *File) Parse() {
	_, f.RootNode = f.parseFile()
	// control that we got to the last token ???
}

// At the top level, just parse everything we can
func (f *File) parseFile() (Tk, Node) {
	scope := f.RootScope()
	file := f.createNode(Range{}, NODE_FILE, scope)
	tk := Tk{
		pos:  0,
		file: f,
	}

	app := newFragment()
	for !tk.IsEof() {
		var node Node
		tk, node = Expression(scope, tk, 0)
		app.append(node)
	}

	if !app.first.IsEmpty() {
		file.SetArgs(app.first)
	}

	return tk, file
}

var lbp_equal = 0
var lbp_commas = 0
var lbp_semicolon = 0

// var rbp_arrow = 0
var lbp_gt = 0
var lbp = 2

func init() {
	for i := range syms {
		syms[i].nud = nudError
		syms[i].led = ledError
	}

	nud(TK_LPAREN, func(scope Scope, tk Tk, lbp int) (Tk, Node) {
		// We are going to check if we have several components to the paren, or just
		// one, in which case we just send it back.
		// an empty () parenthesis block is an error as it doesn't mean anything.

		next, exp := Expression(scope, tk.Next(), 0)
		// check if we end with a parenthesis
		if next, ok := next.consume(TK_RPAREN); ok {
			exp.ExtendRange(tk.Range())
			return next, exp
		}

		// If we didn't encounter ), we want a comma
		next, _ = next.expect(TK_COMMA)

		app := newFragment()
		app.append(exp)

		for !next.IsClosing() {
			next, exp = Expression(scope, next, 0)
			if !next.Is(TK_RPAREN) {
				next, _ = next.expect(TK_COMMA) // there can be a comma
			}
			app.append(exp)
		}

		tup := tk.createTuple(scope, app.first)
		next, _ = next.expect(TK_RPAREN, func(tk Tk) { tup.ExtendRange(tk.Range()) })

		return next, tup
	})

	nud(KW_NAMESPACE, func(scope Scope, tk Tk, lbp int) (Tk, Node) {
		// res.Block = res.CreateBlock()
		next, name := Expression(scope, tk.Next(), 0)
		next, _ = tk.expect(TK_LBRACKET)

		nmsp_scope := scope.subScope()
		next, block := parseBlock(nmsp_scope, next, 0)

		nmsp := tk.createNamespace(scope, name, block)
		nmsp_scope.setOwner(nmsp)
		return next, nmsp
	})

	// { , a block
	nud(TK_LBRACKET, parseBlock)

	nud(TK_LBRACE, func(scope Scope, tk Tk, _ int) (Tk, Node) {
		// function call !
		next := tk.Next()

		fragment := newFragment()
		for !next.IsClosing() {
			var exp Node
			next, exp = Expression(scope, next, 0)
			fragment.append(exp)
			if !next.Is(TK_RBRACE) {
				next, _ = next.consume(TK_COMMA)
			}
		}
		array := tk.createArrayLiteral(scope, fragment.first)

		next, _ = next.expect(TK_RBRACE, func(tk Tk) {
			array.ExtendRange(tk.Range())
		})

		return next, array
	})

	// The doc comment forwards the results but sets itself first on the node that resulted
	// Doc comments whose next meaningful token are other doc comments or the end of the file
	// are added at the module level
	nud(TK_DOCCOMMENT, func(scope Scope, tk Tk, lbp int) (Tk, Node) {
		// tkpos := b.current - 1           // the parsed token position
		next, node := Expression(scope, tk.Next(), lbp) // forward the current lbp to the expression
		tk.file.DocCommentMap[node.pos] = tk.pos
		return next, node
	})

	nud(KW_FOR, parseFor)
	nud(KW_WHILE, parseWhile)
	nud(KW_IF, parseIf)

	nud(KW_VAR, parseVar)

	nud(KW_CONST, func(scope Scope, tk Tk, lbp int) (Tk, Node) {
		next, va := parseVar(scope, tk.Next(), lbp)
		va.SetFlag(FLAG_CONST)
		return next, va
	})

	nud(KW_IMPORT, parseImport)

	lbp += 2

	// return ...
	// will return an empty node if
	nud(KW_RETURN, func(scope Scope, tk Tk, lbp int) (Tk, Node) {
		var res Node
		next := tk
		// do not try to get next expression is return is immediately followed
		// by } or ]
		if tk.Peek(TK_RPAREN, TK_RBRACKET) {
			// return can only return nothing if it is at the end of a block or expression
			res = EmptyNode
		} else {
			next, res = Expression(scope, tk.Next(), lbp)
		}

		return next, tk.createReturn(scope, res)
	})

	lbp += 2

	nud(KW_TYPE, parseTypeDecl)
	nud(KW_STRUCT, parseStruct)

	lbp += 2

	// , creates a tuple
	// maybe it should be handled in the different places where comma is expected,
	// which is to say in lists like (, , ) or [, ,]
	// comma_lbp := lbp
	// led(TK_COMMA, func(scope Scope, tk Tk, left Node) (Tk, Node) {
	// 	if c.currentTokenIs(TK_RBRACE, TK_RBRACKET, TK_RPAREN) {
	// 		return left
	// 	}
	// 	right := c.Expression(scope, comma_lbp) // this should give us a right associative tree
	// 	switch v := left.(type) {
	// 	case *Tuple:
	// 		return v.AddChildren(right)
	// 	default:
	// 		return tk.CreateTuple().AddChildren(left, right)
	// 	}
	// })

	lbp += 2
	lbp_equal = lbp

	// =
	binary(TK_EQ, NODE_BIN_ASSIGN)
	binary(KW_IS, NODE_BIN_IS)

	// fn eats up the expression right next to it
	nud(KW_FN, parseFn)
	nud(KW_METHOD, parseFn)

	lbp += 2

	// unary(KW_LOCAL)
	// unary(KW_CONST)
	binary(TK_EQEQ, NODE_BIN_EQ)
	binary(TK_NOTEQ, NODE_BIN_NEQ)

	lbp += 2

	binary(TK_LT, NODE_BIN_LT)
	lbp_gt = lbp
	binary(TK_GT, NODE_BIN_GT)
	// unary(TK_ELLIPSIS) // ???

	lbp += 2

	binary(TK_PIPE, NODE_BIN_BITOR)
	// conflict with bitwise or !
	// how the hell am I supposed to tell the difference between the two ?
	// led(TK_PIPE, func(scope Scope, tk Tk, left Node) (Tk, Node) {
	// 	right := c.Expression(scope, lbp)
	// 	if v, ok := left.(*Union); ok {
	// 		return v.AddTypeExprs(right)
	// 	}
	// 	return tk.CreateUnion().AddTypeExprs(left, right)
	// })

	lbp += 2

	unary(TK_PLUS, NODE_UNA_PLUS)
	unary(TK_MIN, NODE_UNA_MIN)
	binary(TK_MIN, NODE_BIN_MIN)
	binary(TK_PLUS, NODE_BIN_PLUS)

	lbp += 2

	binary(TK_STAR, NODE_BIN_MUL)
	binary(TK_DIV, NODE_BIN_DIV)
	binary(TK_MOD, NODE_BIN_MOD)

	lbp += 2

	binary(KW_IS, NODE_BIN_IS)
	unary(TK_ELLIPSIS, NODE_UNA_ELLIPSIS)

	lbp += 2

	// parseParens()
	// When used right next to an expression, then paren is a function call
	// handleParens(NODE_LIST, NODE_FNCALL, TK_LPAREN, TK_RPAREN, true)

	led(TK_PLUSPLUS, func(scope Scope, tk Tk, left Node) (Tk, Node) {
		return tk.Next(), tk.createUnaryOp(scope, NODE_UNA_PLUSPLUS, left)
	})
	led(TK_MINMIN, func(scope Scope, tk Tk, left Node) (Tk, Node) {
		return tk.Next(), tk.createUnaryOp(scope, NODE_UNA_MINMIN, left)
	})

	lbp += 2

	// surrounding(NODE_BLOCK, NODE_BLOCK, TK_LBRACKET, TK_RBRACKET, false)

	// the index operator
	// nud(TK_LBRACE, parseLbraceNud)
	nud(TK_STAR, func(scope Scope, tk Tk, lbp int) (Tk, Node) {
		next := tk.Next()
		next, typeexpr := Expression(scope, next, syms[TK_MINMIN].lbp+1)
		if typeexpr.IsEmpty() {
			tk.reportError("expected * to be followed by a type name")
		}
		return next, tk.createUnaPointer(scope, typeexpr)
	})

	nud(TK_AMP, func(scope Scope, tk Tk, lbp int) (Tk, Node) {
		next := tk.Next()
		next, expr := Expression(scope, next, syms[TK_MINMIN].lbp+1)
		if expr.IsEmpty() {
			tk.reportError("expected & to be followed by an expression")
		}
		return next, tk.createUnaRef(scope, expr)
	})

	lbp += 2

	// led(TK_FATARROW, parseFnFatArrow)
	// binary(NODE_FNDEF, TK_FATARROW)

	led(TK_LBRACE, func(scope Scope, tk Tk, left Node) (Tk, Node) {
		// function call !
		fragment := newFragment()
		next := tk.Next()

		for !next.IsClosing() {
			var exp Node
			next, exp = Expression(scope, next, 0)
			fragment.append(exp)
			next, _ = next.consume(TK_COMMA)
		}

		index := tk.createBinOp(scope, NODE_BIN_INDEX, left, fragment.first)
		next, _ = next.expect(TK_RBRACE, func(tk Tk) {
			index.ExtendRange(tk.Range())
		})

		return next, index
	})

	lbp += 2

	led(TK_LPAREN, func(scope Scope, tk Tk, left Node) (Tk, Node) {
		// function call !
		next := tk.Next()
		fragment := newFragment()

		for !next.IsClosing() {
			var exp Node
			next, exp = Expression(scope, next, 0)
			fragment.append(exp)
			next, _ = next.consume(TK_COMMA)
		}

		call := tk.createBinOp(scope, NODE_BIN_CALL, left, fragment.first)
		next, _ = next.expect(TK_RPAREN, func(tk Tk) {
			call.ExtendRange(tk.Range())
		})

		return next, call
	})

	lbp += 2

	binary(TK_DOT, NODE_BIN_DOT)
	// binary(KW_AS)

	lbp += 2
	binary(TK_AT, NODE_BIN_CAST)
	lbp += 2

	nud(TK_QUOTE, parseQuote)

	literal(KW_TRUE, NODE_LIT_TRUE)
	literal(KW_FALSE, NODE_LIT_FALSE)
	literal(KW_NULL, NODE_LIT_NULL)
	literal(KW_VOID, NODE_LIT_VOID)
	literal(TK_CHAR, NODE_LIT_CHAR)
	literal(TK_NUMBER, NODE_LIT_NUMBER)
	literal(TK_RAWSTR, NODE_LIT_RAWSTR)

	nud(TK_ID, func(scope Scope, tk Tk, lbp int) (Tk, Node) {
		return tk.Next(), tk.createIdNode(scope)
	})

}

///////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////

var syms = make([]prattTk, TK__MAX) // Far more than necessary

func nudError(scope Scope, tk Tk, rbp int) (Tk, Node) {
	tk.reportError(`unexpected '`, tk.GetText(), `'`)
	return Expression(scope, tk.Next(), rbp)
}

func ledError(_ Scope, tk Tk, left Node) (Tk, Node) {
	tk.reportError(`unexpected '`, tk.GetText(), `'`)
	return tk.Next(), left
}

///////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////

func parseImport(scope Scope, tk Tk, _ int) (Tk, Node) {
	// import is always (imp module subexp name)
	// module is either a string or a path expression
	var mod Node
	var iter = tk.Next()
	var ok bool

	if iter.Is(TK_RAWSTR) {
		mod = iter.createNode(scope, NODE_LIT_RAWSTR)
		iter = iter.Next()
	} else {
		iter, mod = Expression(scope, iter, syms[TK_DOT].lbp-1)
	}

	if next, ok := iter.consume(KW_AS); ok {
		var name Node
		if next.Is(TK_ID) {
			name = next.createIdNode(scope)
			next = next.Next()
		}
		imp := tk.createImport(scope, mod, name, EmptyNode)

		if !name.IsEmpty() {
			// Add the import to the current scope.
			scope.addSymbolFromIdNode(name, imp)
		}

		return next, imp
	}

	if iter, ok = iter.consume(TK_LPAREN); !ok {
		iter.reportError("malformed import expression, expected '(' or 'as'")
		return iter, EmptyNode
	}

	fragment := newFragment()
	for !iter.IsClosing() {
		mod2 := mod.Clone()

		var path Node
		iter, path = Expression(scope, iter, syms[TK_DOT].lbp-1) // we want the tk_dots

		if iter, ok = iter.consume(KW_AS); ok {

			var as Node
			prev := iter
			if iter.Is(TK_ID) {
				as = iter.createIdNode(scope)
				iter = iter.Next()
			}

			imp := prev.createImport(scope, mod2, as, path)
			if !as.IsEmpty() {
				scope.addSymbolFromIdNode(as, imp)
			}

			fragment.append(imp)
		} else {
			id2 := iter.createIdNode(scope)
			imp := iter.createImport(scope, mod2, id2, path)
			scope.addSymbolFromIdNode(id2, imp)
			iter = iter.Next()
			fragment.append(imp)
		}
		iter, _ = iter.consume(TK_COMMA)
	}
	iter, _ = iter.expect(TK_RPAREN)

	return iter, fragment.first
}

///////////////////////////////////////////////////////
// "
func parseQuote(scope Scope, tk Tk, _ int) (Tk, Node) {
	iter := tk.Next()
	fragment := newFragment()
	for !iter.IsEof() && !iter.Is(TK_QUOTE) {
		var exp Node
		iter, exp = Expression(scope, iter, 0)
		fragment.append(exp)
	}

	str := tk.createString(scope, fragment.first)
	iter, _ = iter.expect(TK_QUOTE, func(tk Tk) {
		str.ExtendRange(tk.Range())
	})

	// this should transform the result to a string
	return iter, str
}

/////////////////////////////////////////////////////
// Special handling for if block
func parseIf(scope Scope, tk Tk, _ int) (Tk, Node) {
	var ok bool
	iter := tk.Next()

	iter, cond := Expression(scope, iter, 0) // can be a block. this could be confusing.

	iter.expect(TK_LBRACKET)
	iter, then := Expression(scope, iter, 0) // most likely, a block.

	var els Node
	if iter, ok = iter.consume(KW_ELSE); ok {
		iter.expect(TK_LBRACKET)
		iter, els = Expression(scope, iter, 0)
	}

	return iter, tk.createIf(scope, cond, then, els)
}

/////////////////////////////////////////////////////
// Special handling for fn
func parseFn(scope Scope, tk Tk, _ int) (Tk, Node) {

	fnscope := scope.subScope()
	iter := tk.Next()

	// Function name, may not exist
	var name Node
	iter, _ = iter.expect(TK_ID, func(tk Tk) {
		name = tk.createIdNode(scope)
	})

	// Template arguments, may not exist
	var tpl Node
	if iter.Is(TK_LBRACE) {
		iter, tpl = parseTemplate(scope, iter, 0)
	}

	// Function arguments, mandatory
	args := newFragment()
	iter, _ = iter.expect(TK_LPAREN)
	for !iter.IsClosing() {
		var arg Node
		iter, arg = parseVar(fnscope, iter, 0)
		if !arg.IsEmpty() {
			args.append(arg)
		} else {
			iter = iter.Next()
		}

		// test for comma presence if not the end of the arguments
		if !iter.Is(TK_RPAREN) {
			iter, _ = iter.expect(TK_COMMA)
		}
	}
	iter, _ = iter.expect(TK_RPAREN)

	// Return type, may not exist
	var rettype Node
	if iter.Is(TK_ARROW) {
		iter = iter.Next()
		iter, rettype = Expression(fnscope, iter, 0)
	}

	// The signature node
	signature := tk.createSignature(scope, tpl, args.first, rettype)

	// Function definition
	var blk Node
	if iter.Is(TK_LBRACKET) {
		iter, blk = parseBlock(fnscope, iter, 0)
		// should register the function somewhere in scope, no ?
		return iter, tk.createFn(scope, name, signature, blk)
	}

	return iter, signature
}

// parseBlock parses a block of code
func parseBlock(scope Scope, tk Tk, _ int) (Tk, Node) {
	// blk := b.createNodeFromToken(tk, NODE_BLOCK)
	app_blk := newFragment()

	for b.asLongAsNotClosingToken() {
		for b.consume(TK_SEMICOLON) != 0 {
			// advance as much as we can if we have semi colons in the input
		}

		if b.isEof() {
			break
		}

		app_blk.append(b.Expression(scope, 0))
	}

	block := b.createBlock(tk, scope, app_blk.first)

	if tk := b.expect(TK_RBRACKET); tk != 0 {
		b.extendRangeFromToken(block, tk) // FIXME, this is ugly
	}

	return block
}

// parseTemplate parses a template declaration, which is enclosed between [ ]
// it is expected that '[' has been consumed, and that tk is '['
func parseTemplate(scope Scope, tk Tk, _ int) (Tk, Node) {
	// tpl := b.createNodeFromToken(tk, NODE_TEMPLATE)
	fragment := newFragment()

	for b.asLongAsNotClosingToken() { // missing WHERE
		v := b.createExpectToken(TK_ID, func(tk Tk) (Tk, Node) {
			return b.createIdNode(tk, scope)
		})
		if v == 0 {
			b.reportErrorAtToken(b.current, "expected a template variable declaration")
			b.advance()
		} else {
			fragment.append(v)
		}
		b.consume(TK_COMMA)
	}
	b.expect(TK_RBRACE)
	return fragment.first
}

// parseTypeDecl parses a type declaration
func parseTypeDecl(scope Scope, tk Tk, _ int) (Tk, Node) {
	name := b.createAndExpectOrEmpty(TK_ID, func(tk Tk) (Tk, Node) {
		return b.createIdNode(tk, scope)
	})

	tpl := b.createIfTokenOrEmpty(TK_LBRACE, func(tk Tk) (Tk, Node) {
		return parseTemplate(b, scope, tk, 0)
	})

	if b.consume(KW_IS) == 0 {
		b.reportErrorAtToken(b.current, `expected 'is' after type declaration`)
	}

	// there might be a pipe here. We don't have to parse a union afterwards because
	// if there is only one type, it doesn't matter.
	b.consume(TK_PIPE)

	typdef := b.Expression(scope, 0)
	// b.file.Nodes = b.nodes
	// log.Print("!!!", b.file.NodeDebug(typdef))
	// raise an error if there is no typedef ?

	typ := b.createType(tk, scope, name, tpl, typdef)
	if name != EmptyNode {
		scope.addSymbolFromIdNode(name, typ)
	}
	return typ
}

func parseStruct(scope Scope, tk Tk, _ int) (Tk, Node) {
	// stru := c.createNodeFromToken(tk, NODE_STRUCT)
	c.expect(TK_LPAREN)

	fragment := c.fragment()
	for !c.isEof() && !c.currentTokenIs(TK_RPAREN, TK_RBRACE, TK_RBRACKET) {
		// parse a var declaration
		v := parseVar(c, scope, c.current, 0)
		if v == 0 {
			c.advance()
		} else {
			// FIXME ensure a field has a type declaration, as struct
			// fields should have them whether they have defaults or not
			fragment.append(v)
		}
		if !c.currentTokenIs(TK_RPAREN) {
			c.consume(TK_COMMA)
		}
	}

	stru := c.createStruct(tk, scope, fragment.first)
	if tk := c.expect(TK_RPAREN); tk != 0 {
		c.extendRangeFromToken(stru, tk)
	}

	return stru
}

func parseFor(scope Scope, tk Tk, _ int) (Tk, Node) {
	forscope := scope.subScope()
	decl := b.Expression(forscope, 0) // this is where a var is created
	b.expect(KW_IN)
	inexp := b.Expression(scope, 0)
	b.expectNoAdvance(TK_LBRACKET) // needs an opening '{'
	block := b.Expression(forscope, 0)
	return b.createFor(tk, scope, decl, inexp, block)
}

func parseWhile(scope Scope, tk Tk, _ int) (Tk, Node) {
	whilescope := scope.subScope()
	cond := b.Expression(whilescope, 0)
	b.expectNoAdvance(TK_LBRACKET)
	block := b.Expression(whilescope, 0)
	return b.createWhile(tk, scope, cond, block)
}

// parse a variable statement, but also a variable declaration inside
// an argument list of a function signature
func parseVar(scope Scope, tk Tk, _ int) (Tk, Node) {
	// first, try to scan the ident
	// this may fail, for dubious reasons
	var ident NodePosition
	if tkident := c.expect(TK_ID); tkident != 0 {
		ident = c.createIdNode(tkident, scope)
	} else {
		ident = EmptyNode
	}

	// there might be a type expression right after the name declaration
	typenode := c.createIfTokenOrEmpty(TK_COLON, func(tk Tk) (Tk, Node) {
		// We scan above '=' level to avoid eating it if there is a default value
		// right after the type declaration
		return c.Expression(scope, syms[TK_EQ].lbp+1) // anything above =
	})

	// default value !
	var expnode Node
	if c.consume(TK_EQ) != 0 {
		expnode = c.Expression(scope, 0)
	} else {
		expnode = c.file.emptyNode()
	}

	if c.current == tk {
		// no var here
		return 0
	}

	return c.createVar(tk, scope, ident, typenode, expnode)
	// Try to parse VAR ourselves
}
