package zoe

func (f *File) Parse() {
	b := f.createNodeBuilder()
	f.RootNodePos = b.parseFile()
}

// At the top level, just parse everything we can
func (b *nodeBuilder) parseFile() NodePosition {
	scope := b.file.RootScope()
	file := b.createNode(Range{}, NODE_FILE, scope)

	app := b.fragment()
	for !b.isEof() {
		r := b.Expression(scope, 0)
		app.append(r)
	}

	if app.first != EmptyNode {
		f := &b.file.Nodes[file]
		f.ArgLen = 1
		f.Args[0] = app.first
		b.extendsNodeRangeFromNode(file, app.first)
	}
	return file
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

	nud(TK_LPAREN, func(b *nodeBuilder, scope Scope, tk TokenPos, lbp int) NodePosition {
		// We are going to check if we have several components to the paren, or just
		// one, in which case we just send it back.
		// an empty () parenthesis block is an error as it doesn't mean anything.

		exp := b.Expression(scope, 0)
		// check if we end with a parenthesis
		if pos := b.consume(TK_RPAREN); pos != 0 {
			b.extendRangeFromToken(exp, pos)
			return exp
		}

		// If we didn't encounter ), we want a comma
		b.expect(TK_COMMA)

		app := b.fragment()
		app.append(exp)

		for b.asLongAsNotClosingToken() {
			exp := b.Expression(scope, 0)
			if !b.currentTokenIs(TK_RPAREN) {
				b.expect(TK_COMMA) // there can be a comma
			}
			app.append(exp)
		}
		tup := b.createTuple(tk, scope, app.first)
		if tok := b.expect(TK_RPAREN); tok != 0 {
			b.extendRangeFromToken(tup, tok)
		}
		return tup
	})

	nud(KW_NAMESPACE, func(b *nodeBuilder, scope Scope, tk TokenPos, lbp int) NodePosition {
		// res.Block = res.CreateBlock()
		name := b.Expression(scope, 0)
		b.expect(TK_LBRACKET)
		nmsp_scope := scope.subScope()
		block := parseBlock(b, nmsp_scope, tk, 0)

		nmsp := b.createNamespace(tk, scope, name, block)
		nmsp_scope.setOwner(nmsp)
		return nmsp
	})

	// { , a block
	nud(TK_LBRACKET, parseBlock)

	nud(TK_LBRACE, func(b *nodeBuilder, scope Scope, tk TokenPos, _ int) NodePosition {
		// function call !
		fragment := b.fragment()

		for b.asLongAsNotClosingToken() {
			exp := b.Expression(scope, 0)
			fragment.append(exp)
			if !b.currentTokenIs(TK_RBRACE) {
				b.consume(TK_COMMA)
			}
		}
		array := b.createArrayLiteral(tk, scope, fragment.first)
		if tk := b.expect(TK_RBRACE); tk != 0 {
			b.extendRangeFromToken(array, tk)
		}

		return array
	})

	// The doc comment forwards the results but sets itself first on the node that resulted
	// Doc comments whose next meaningful token are other doc comments or the end of the file
	// are added at the module level
	nud(TK_DOCCOMMENT, func(b *nodeBuilder, scope Scope, tk TokenPos, lbp int) NodePosition {
		tkpos := b.current - 1           // the parsed token position
		next := b.Expression(scope, lbp) // forward the current lbp to the expression
		b.doccommentMap[next] = tkpos
		return next
	})

	nud(KW_FOR, parseFor)
	nud(KW_WHILE, parseWhile)
	nud(KW_IF, parseIf)

	nud(KW_VAR, parseVar)

	nud(KW_CONST, func(b *nodeBuilder, scope Scope, tk TokenPos, lbp int) NodePosition {
		va := parseVar(b, scope, tk, lbp)
		b.file.Nodes[va].SetFlag(FLAG_CONST)
		return va
	})

	nud(KW_IMPORT, parseImport)

	lbp += 2

	// return ...
	// will return an empty node if
	nud(KW_RETURN, func(c *nodeBuilder, scope Scope, tk TokenPos, lbp int) NodePosition {
		var res NodePosition
		// do not try to get next expression is return is immediately followed
		// by } or ]
		if c.currentTokenIs(TK_RPAREN, TK_RBRACKET) {
			// return can only return nothing if it is at the end of a block or expression
			res = EmptyNode
		} else {
			res = c.Expression(scope, lbp)
		}

		return c.createNodeFromToken(tk, NODE_RETURN, scope, res)
	})

	lbp += 2

	nud(KW_TYPE, parseTypeDecl)
	nud(KW_STRUCT, parseStruct)

	lbp += 2

	// , creates a tuple
	// maybe it should be handled in the different places where comma is expected,
	// which is to say in lists like (, , ) or [, ,]
	// comma_lbp := lbp
	// led(TK_COMMA, func(c *nodeBuilder, scope Scope, tk TokenPos, left NodePosition) NodePosition {
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
	// led(TK_PIPE, func(c *nodeBuilder, scope Scope, tk TokenPos, left NodePosition) NodePosition {
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

	led(TK_PLUSPLUS, func(b *nodeBuilder, scope Scope, tk TokenPos, left NodePosition) NodePosition {
		return b.createUnaryOp(tk, NODE_UNA_PLUSPLUS, scope, left)
	})
	led(TK_MINMIN, func(b *nodeBuilder, scope Scope, tk TokenPos, left NodePosition) NodePosition {
		return b.createUnaryOp(tk, NODE_UNA_MINMIN, scope, left)
	})

	lbp += 2

	// surrounding(NODE_BLOCK, NODE_BLOCK, TK_LBRACKET, TK_RBRACKET, false)

	// the index operator
	// nud(TK_LBRACE, parseLbraceNud)
	nud(TK_STAR, func(b *nodeBuilder, scope Scope, tk TokenPos, lbp int) NodePosition {
		typeexpr := b.Expression(scope, syms[TK_MINMIN].lbp+1)
		if typeexpr == EmptyNode {
			b.reportErrorAtToken(tk, "expected * to be followed by a type name")
		}
		return b.createUnaPointer(tk, scope, typeexpr)
	})

	nud(TK_AMP, func(b *nodeBuilder, scope Scope, tk TokenPos, lbp int) NodePosition {
		expr := b.Expression(scope, syms[TK_MINMIN].lbp+1)
		if expr == EmptyNode {
			b.reportErrorAtToken(tk, "expected & to be followed by an expression")
		}
		return b.createUnaRef(tk, scope, expr)
	})

	lbp += 2

	// led(TK_FATARROW, parseFnFatArrow)
	// binary(NODE_FNDEF, TK_FATARROW)

	led(TK_LBRACE, func(b *nodeBuilder, scope Scope, tk TokenPos, left NodePosition) NodePosition {
		// function call !
		fragment := b.fragment()

		for b.asLongAsNotClosingToken() {
			exp := b.Expression(scope, 0)
			fragment.append(exp)
			b.consume(TK_COMMA)
		}
		index := b.createBinOp(tk, NODE_BIN_INDEX, scope, left, fragment.first)
		if tk := b.expect(TK_RBRACE); tk != 0 {
			b.extendRangeFromToken(index, tk)
		}
		return index
	})

	lbp += 2

	led(TK_LPAREN, func(b *nodeBuilder, scope Scope, tk TokenPos, left NodePosition) NodePosition {
		// function call !
		fragment := b.fragment()

		for b.asLongAsNotClosingToken() {
			exp := b.Expression(scope, 0)
			fragment.append(exp)
			b.consume(TK_COMMA)
		}
		call := b.createBinOp(tk, NODE_BIN_CALL, scope, left, fragment.first)
		if tk := b.expect(TK_RPAREN); tk != 0 {
			b.extendRangeFromToken(call, tk)
		}
		return call
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

	nud(TK_ID, func(b *nodeBuilder, scope Scope, tk TokenPos, lbp int) NodePosition {
		return b.createIdNode(tk, scope)
	})

}

///////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////

var syms = make([]prattTk, TK__MAX) // Far more than necessary

func nudError(b *nodeBuilder, scope Scope, tk TokenPos, rbp int) NodePosition {
	b.reportErrorAtToken(tk, `unexpected '`, b.getTokenText(tk), `'`)
	return b.Expression(scope, rbp)
}

func ledError(b *nodeBuilder, _ Scope, tk TokenPos, left NodePosition) NodePosition {
	b.reportErrorAtToken(tk, `unexpected '`, b.getTokenText(tk), `'`)
	return left
}

///////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////

func parseImport(b *nodeBuilder, scope Scope, tk TokenPos, _ int) NodePosition {
	// import is always (imp module subexp name)
	// module is either a string or a path expression
	var mod NodePosition

	mod = b.createIfToken(TK_RAWSTR, func(tk TokenPos) NodePosition {
		return b.createNodeFromToken(tk, NODE_LIT_RAWSTR, scope)
	})

	if mod == 0 {
		mod = b.Expression(scope, syms[TK_DOT].lbp-1)
	}

	if as := b.consume(KW_AS); as != 0 {
		name := b.createAndExpectOrEmpty(TK_ID, func(tk TokenPos) NodePosition {
			return b.createIdNode(tk, scope)
		})

		imp := b.createNodeFromToken(tk, NODE_IMPORT, scope, mod, name, EmptyNode)

		if name != EmptyNode {
			scope.addSymbolFromIdNode(name, imp)
		}

		return imp
	}

	if b.consume(TK_LPAREN) == 0 {
		b.reportErrorAtToken(tk, "malformed import expression, expected '(' or 'as'")
		return EmptyNode
	}

	fragment := b.fragment()
	for b.asLongAsNotClosingToken() {
		mod2 := b.cloneNode(mod)
		cur := b.current
		path := b.Expression(scope, syms[TK_DOT].lbp-1) // we want the tk_dots

		if b.consume(KW_AS) != 0 {
			as := b.createAndExpectOrEmpty(TK_ID, func(tk TokenPos) NodePosition {
				return b.createIdNode(tk, scope)
			})
			imp := b.createNodeFromToken(cur, NODE_IMPORT, scope, mod2, as, path)

			if as != EmptyNode {
				scope.addSymbolFromIdNode(as, imp)
			}
			fragment.append(imp)
		} else {
			id2 := b.createIdNode(cur, scope)

			imp := b.createNodeFromToken(cur, NODE_IMPORT, scope, mod2, id2, path)
			scope.addSymbolFromIdNode(id2, imp)
			fragment.append(imp)
		}
		b.consume(TK_COMMA)
	}
	b.expect(TK_RPAREN)

	return fragment.first
}

///////////////////////////////////////////////////////
// "
func parseQuote(c *nodeBuilder, scope Scope, tk TokenPos, _ int) NodePosition {
	fragment := c.fragment()
	for !c.isEof() && !c.currentTokenIs(TK_QUOTE) {
		fragment.append(c.Expression(scope, 0))
	}
	str := c.createString(tk, scope, fragment.first)
	if tk2 := c.expect(TK_QUOTE); tk2 != 0 {
		c.extendRangeFromToken(str, tk2)
	}
	// this should transform the result to a string
	return str
}

/////////////////////////////////////////////////////
// Special handling for if block
func parseIf(c *nodeBuilder, scope Scope, tk TokenPos, _ int) NodePosition {
	cond := c.Expression(scope, 0) // can be a block. this could be confusing.
	c.expectNoAdvance(TK_LBRACKET)
	then := c.Expression(scope, 0) // most likely, a block.
	els := c.createIfTokenOrEmpty(KW_ELSE, func(tk TokenPos) NodePosition {
		c.expectNoAdvance(TK_LBRACKET)
		return c.Expression(scope, 0)
	})

	return c.createIf(tk, scope, cond, then, els)
}

/////////////////////////////////////////////////////
// Special handling for fn
func parseFn(b *nodeBuilder, scope Scope, tk TokenPos, _ int) NodePosition {

	fnscope := scope.subScope()

	name := b.createIfTokenOrEmpty(TK_ID, func(tk TokenPos) NodePosition {
		return b.createIdNode(tk, scope)
	})

	tpl := b.createIfTokenOrEmpty(TK_LBRACE, func(tk TokenPos) NodePosition {
		return parseTemplate(b, scope, tk, 0)
	})

	// args := c.createNodeFromCurrentToken(NODE_TUPLE)
	args := b.fragment()
	b.expect(TK_LPAREN)
	for !b.currentTokenIs(TK_RPAREN, TK_ARROW, TK_FATARROW) {
		arg := parseVar(b, fnscope, b.current, 0)
		if arg != 0 {
			args.append(arg)
		} else {
			b.advance()
		}
		if !b.currentTokenIs(TK_RPAREN) {
			b.expect(TK_COMMA)
		}
		// test for comma presence
	}
	b.expect(TK_RPAREN)

	rettype := b.createIfTokenOrEmpty(TK_ARROW, func(tk TokenPos) NodePosition {
		return b.Expression(fnscope, 0)
	})

	signature := b.createSignature(tk, scope, tpl, args.first, rettype)

	var blk NodePosition
	if b.currentTokenIs(TK_LBRACKET) {
		b.advance()
		blk = parseBlock(b, fnscope, b.current, 0)
		return b.createFn(tk, scope, name, signature, blk)
	}

	return signature
}

// parseBlock parses a block of code
func parseBlock(b *nodeBuilder, scope Scope, tk TokenPos, _ int) NodePosition {
	// blk := b.createNodeFromToken(tk, NODE_BLOCK)
	app_blk := b.fragment()

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
func parseTemplate(b *nodeBuilder, scope Scope, tk TokenPos, _ int) NodePosition {
	// tpl := b.createNodeFromToken(tk, NODE_TEMPLATE)
	fragment := b.fragment()

	for b.asLongAsNotClosingToken() { // missing WHERE
		v := b.createExpectToken(TK_ID, func(tk TokenPos) NodePosition {
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
func parseTypeDecl(b *nodeBuilder, scope Scope, tk TokenPos, _ int) NodePosition {
	name := b.createAndExpectOrEmpty(TK_ID, func(tk TokenPos) NodePosition {
		return b.createIdNode(tk, scope)
	})

	tpl := b.createIfTokenOrEmpty(TK_LBRACE, func(tk TokenPos) NodePosition {
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

func parseStruct(c *nodeBuilder, scope Scope, tk TokenPos, _ int) NodePosition {
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

func parseFor(b *nodeBuilder, scope Scope, tk TokenPos, _ int) NodePosition {
	forscope := scope.subScope()
	decl := b.Expression(forscope, 0) // this is where a var is created
	b.expect(KW_IN)
	inexp := b.Expression(scope, 0)
	b.expectNoAdvance(TK_LBRACKET) // needs an opening '{'
	block := b.Expression(forscope, 0)
	return b.createFor(tk, scope, decl, inexp, block)
}

func parseWhile(b *nodeBuilder, scope Scope, tk TokenPos, _ int) NodePosition {
	whilescope := scope.subScope()
	cond := b.Expression(whilescope, 0)
	b.expectNoAdvance(TK_LBRACKET)
	block := b.Expression(whilescope, 0)
	return b.createWhile(tk, scope, cond, block)
}

// parse a variable statement, but also a variable declaration inside
// an argument list of a function signature
func parseVar(c *nodeBuilder, scope Scope, tk TokenPos, _ int) NodePosition {
	// first, try to scan the ident
	// this may fail, for dubious reasons
	var ident NodePosition
	if tkident := c.expect(TK_ID); tkident != 0 {
		ident = c.createIdNode(tkident, scope)
	} else {
		ident = EmptyNode
	}

	// there might be a type expression right after the name declaration
	typenode := c.createIfTokenOrEmpty(TK_COLON, func(tk TokenPos) NodePosition {
		// We scan above '=' level to avoid eating it if there is a default value
		// right after the type declaration
		return c.Expression(scope, syms[TK_EQ].lbp+1) // anything above =
	})

	// default value !
	var expnode NodePosition
	if c.consume(TK_EQ) != 0 {
		expnode = c.Expression(scope, 0)
	} else {
		expnode = EmptyNode
	}

	if c.current == tk {
		// no var here
		return 0
	}

	return c.createVar(tk, scope, ident, typenode, expnode)
	// Try to parse VAR ourselves
}
