package zoe

func (f *File) Parse() {
	b := f.createNodeBuilder()
	f.RootNodePos = b.parseFile()
	f.Nodes = b.nodes
}

// At the top level, just parse everything we can
func (b *nodeBuilder) parseFile() NodePosition {
	// the error node at position 0
	b.createEmptyNode()

	res := b.createNode(Range{}, NODE_FILE) // should it be a file ?
	app := b.appender(res)
	for !b.isEof() {
		r := b.Expression(0)
		app.append(r)
	}
	return res
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

	nud(KW_NAMESPACE, func(b *nodeBuilder, tk TokenPos, lbp int) NodePosition {
		res := b.createNodeFromToken(tk, NODE_DECL_NMSP)
		// res.Block = res.CreateBlock()
		name := b.Expression(0)
		b.expect(TK_LBRACKET)
		block := parseBlock(b, tk, 0)

		// should be b.createNamespace(tk.Range, name, block)
		b.setNodeChildren(res, name, block)
		return res
	})

	// { , a block
	nud(TK_LBRACKET, parseBlock)

	// The doc comment forwards the results but sets itself first on the node that resulted
	// Doc comments whose next meaningful token are other doc comments or the end of the file
	// are added at the module level
	nud(TK_DOCCOMMENT, func(b *nodeBuilder, tk TokenPos, lbp int) NodePosition {
		tkpos := b.current - 1    // the parsed token position
		next := b.Expression(lbp) // forward the current lbp to the expression
		b.doccommentMap[next] = tkpos
		return next
	})

	// nud(KW_FOR, parseFor)
	nud(KW_IF, parseIf)

	nud(KW_VAR, parseVar)

	// nud(KW_IMPORT, parseImport)

	lbp += 2

	// return ...
	// will return an empty node if
	nud(KW_RETURN, func(c *nodeBuilder, tk TokenPos, lbp int) NodePosition {
		var res NodePosition
		// do not try to get next expression is return is immediately followed
		// by } or ]
		if c.currentTokenIs(TK_RPAREN, TK_RBRACKET) {
			// return can only return nothing if it is at the end of a block or expression
			res = c.createEmptyNode()
		} else {
			res = c.Expression(lbp)
		}

		return c.createNodeFromToken(tk, NODE_RETURN, res)
	})

	lbp += 2

	nud(KW_TEMPLATE, parseTemplate)

	nud(KW_TYPE, parseTypeDecl)
	nud(KW_STRUCT, parseStruct)

	// ; is a separator that creates a fragment
	lbp_semicolon = lbp
	led(TK_SEMICOLON, parseSemiColon)

	lbp += 2

	// , creates a tuple
	// maybe it should be handled in the different places where comma is expected,
	// which is to say in lists like (, , ) or [, ,]
	// comma_lbp := lbp
	// led(TK_COMMA, func(c *nodeBuilder, tk TokenPos, left NodePosition) NodePosition {
	// 	if c.currentTokenIs(TK_RBRACE, TK_RBRACKET, TK_RPAREN) {
	// 		return left
	// 	}
	// 	right := c.Expression(comma_lbp) // this should give us a right associative tree
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
	binary(TK_EQ, NODE_BIN_EQ)
	binary(KW_IS, NODE_BIN_IS)

	// fn eats up the expression right next to it
	nud(KW_FN, parseFn)

	lbp += 2

	// unary(KW_LOCAL)
	// unary(KW_CONST)

	lbp += 2

	binary(TK_LT, NODE_BIN_LT)
	lbp_gt = lbp
	binary(TK_GT, NODE_BIN_GT)
	// unary(TK_ELLIPSIS) // ???

	lbp += 2

	// conflict with bitwise or !
	// how the hell am I supposed to tell the difference between the two ?
	// led(TK_PIPE, func(c *nodeBuilder, tk TokenPos, left NodePosition) NodePosition {
	// 	right := c.Expression(lbp)
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

	lbp += 2

	binary(KW_IS, NODE_BIN_IS)

	lbp += 2

	// parseParens()
	// When used right next to an expression, then paren is a function call
	// handleParens(NODE_LIST, NODE_FNCALL, TK_LPAREN, TK_RPAREN, true)

	lbp += 2

	// surrounding(NODE_BLOCK, NODE_BLOCK, TK_LBRACKET, TK_RBRACKET, false)

	// the index operator
	// nud(TK_LBRACE, parseLbraceNud)

	lbp += 2

	// led(TK_FATARROW, parseFnFatArrow)
	// binary(NODE_FNDEF, TK_FATARROW)

	lbp += 2

	// led(TK_ARROW, parseArrow)

	lbp += 2

	binary(TK_DOT, NODE_BIN_DOT)
	binary(TK_AT, NODE_BIN_CAST)
	// binary(KW_AS)

	lbp += 2

	// lbp_colcol := lbp
	binary(TK_COLCOL, NODE_BIN_NMSP)

	lbp += 2
	// all the terminals. Lbp was raised, but this is not necessary

	nud(TK_QUOTE, parseQuote)

	literal(KW_TRUE, NODE_LIT_TRUE)
	literal(KW_FALSE, NODE_LIT_FALSE)
	literal(KW_NULL, NODE_LIT_NULL)
	literal(KW_VOID, NODE_LIT_VOID)
	literal(TK_CHAR, NODE_LIT_CHAR)
	literal(TK_NUMBER, NODE_LIT_NUMBER)
	literal(TK_RAWSTR, NODE_LIT_RAWSTR)

	nud(TK_ID, func(b *nodeBuilder, tk TokenPos, lbp int) NodePosition {
		return b.createIdNode(tk)
	})

}

///////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////

var syms = make([]prattTk, TK__MAX) // Far more than necessary

func nudError(b *nodeBuilder, tk TokenPos, rbp int) NodePosition {
	b.reportErrorAtToken(tk, `unexpected '`, b.getTokenText(tk), `'`)
	return b.Expression(rbp)
}

func ledError(b *nodeBuilder, tk TokenPos, left NodePosition) NodePosition {
	b.reportErrorAtToken(tk, `unexpected '`, b.getTokenText(tk), `'`)
	return left
}

///////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////

// Group expression between two parseParens tokens
// nk : the node of the list returned when nud
// lednk : the node of the list returned when led
// opening : the opening token
// closing : the closing token
// reduce : whether it is allowed to reduce the list to the single expression
// func parseParens() {
// 	s := &syms[TK_LPAREN]
// 	s.lbp = lbp

// 	parseParen := func(c *nodeBuilder, tk TokenPos) NodePosition {
// 		if c.currentTokenIs(TK_RPAREN) {
// 			// () is the empty tuple
// 			res := tk.CreateTuple()
// 			res.ExtendPosition(c.expect(TK_RPAREN))
// 			return res
// 		}

// 		exp := c.Expression(0)

// 		_, is_tuple := exp.(*Tuple)
// 		if is_tuple {
// 			// we only include the parenthesis in the position if exp is a tuple
// 			exp.ExtendPosition(tk)
// 		}

// 		if !c.currentTokenIs(TK_RPAREN) {
// 			c.reportError(tk.Position, `missing closing ')'`)
// 		} else {
// 			if is_tuple {
// 				exp.ExtendPosition(c.current)
// 			}
// 			c.advance()
// 		}

// 		return exp
// 	}

// 	//
// 	s.led = func(c *nodeBuilder, tk TokenPos, left NodePosition) NodePosition {
// 		exp := parseParen(c, tk).EnsureTuple()

// 		// On our left, we may have an FnDef or FnDecl, in which case we will graft
// 		// ourselves onto them
// 		switch v := left.(type) {
// 		case *FnDecl:
// 			v.FnDef.EnsureSignature(func(v *Signature) {
// 				v.SetArgs(exp.ToVars())
// 			})
// 			return v
// 		case *FnDef:
// 			return v.EnsureSignature(func(s *Signature) {
// 				s.SetArgs(exp.ToVars())
// 			})
// 		}

// 		// Otherwise, this is a plain fncall
// 		return tk.CreateFnCall().SetLeft(left).SetArgs(exp)
// 	}

// 	//
// 	s.nud = func(c *nodeBuilder, tk TokenPos, _ int) NodePosition {
// 		res := parseParen(c, tk)
// 		return res
// 	}
// }

// Handle import
// func parseImport(c *nodeBuilder, tk TokenPos, _ int) NodePosition {
// 	import_exp := c.Expression(0)
// 	// log.Print(module_or_namespace.GetText())

// 	switch v := import_exp.(type) {
// 	case *Operation:
// 		ident, is_ident := v.Right().(*BaseIdent)
// 		if !v.Is(KW_AS) || !is_ident {
// 			v.ReportError("expected 'as' <ident>")
// 		}
// 		return tk.CreateImport().SetPath(v.Left()).SetAs(ident)
// 	case *FnCall:
// 		res := tk.CreateFragment()
// 		mod := v.Left
// 		for _, a := range v.Args.Children {
// 			switch im := a.(type) {
// 			case *Operation:
// 				if ident, ok := im.Right().(*BaseIdent); ok {
// 					res.AddChildren(im.CreateImport().SetPath(mod).SetSubPath(im.Left()).SetAs(ident))
// 				}
// 				continue
// 			case *BaseIdent:
// 				res.AddChildren(im.CreateImport().SetPath(mod).SetSubPath(im).SetAs(im))
// 				continue
// 			}
// 			a.ReportError(`invalid import statement`)
// 		}
// 		return res
// 		// return tk.CreateImportList().SetPath(v.Left).SetNames(v.Args)
// 	case *String:
// 		return tk.CreateImport().SetPath(v)
// 	}

// 	import_exp.ReportError(`invalid import statement`)
// 	return import_exp
// }

///////////////////////////////////////////////////////
// "
func parseQuote(c *nodeBuilder, tk TokenPos, _ int) NodePosition {
	str := c.createNodeFromToken(tk, NODE_STRING)
	app := c.appender(str)
	for !c.currentTokenIs(TK_QUOTE) {
		app.append(c.Expression(0))
	}
	if tk2 := c.expect(TK_QUOTE); tk2 != 0 {
		c.extendRangeFromToken(str, tk2)
	}
	// this should transform the result to a string
	return str
}

////////////////////////////////////////////////////////
// ->
// It has a higher precedence that a function call, and thus attemps to transform
// a tuple to arguments.
//
// Unlike =>, which accepts a single identifier as its left operand,
// -> *requires* the left member to be a tuple, even if an empty one.
//
// It may be followed by a definition with => (that it handles itself)
// Or by a { which will return a block
// func parseArrow(c *nodeBuilder, _ TokenPos, left NodePosition) NodePosition {
// 	// left contains the fndef or fndecl
// 	right := c.Expression(0)
// 	var block *Block
// 	var ok bool

// 	if c.currentTokenIs(TK_LBRACKET) {
// 		bk := c.current
// 		if block, ok = c.Expression(0).(*Block); !ok {
// 			c.reportError(bk, `expected a block`)
// 		}
// 	}

// 	// left is necessarily a tuple. any other type is an error
// 	switch v := left.(type) {
// 	case *FnDecl:
// 		v.ExtendPosition(right)
// 		v.FnDef.Signature.SetReturnTypeExp(right)
// 		if block != nil {
// 			v.FnDef.SetDefinition(block)
// 		}
// 		return v
// 	case *FnDef:
// 		v.ExtendPosition(right)
// 		v.Signature.SetReturnTypeExp(right)
// 		if block != nil {
// 			v.SetDefinition(block)
// 		}
// 		return v
// 	default:
// 		left.ReportError(`the left side of '->' must be a function definition`)
// 		return right
// 	}
// }

// func parseFnFatArrow(c *nodeBuilder, tk TokenPos, left NodePosition) NodePosition {
// 	// left is a list of arguments
// 	// right of => is the implementation of the function

// 	impl := c.Expression(0) // it is a block or a single expression
// 	var block *Block
// 	var ok bool
// 	if block, ok = impl.(*Block); !ok {
// 		block = tk.CreateBlock().AddChildren(impl)
// 	}

// 	switch v := left.(type) {
// 	case *FnDef:
// 		v.SetDefinition(block)
// 		return v
// 	}

// 	left.ReportError(`left hand side of '=>' must be a lambda function definition`)
// 	return left
// }

/////////////////////////////////////////////////////
// Special handling for if block
func parseIf(c *nodeBuilder, tk TokenPos, _ int) NodePosition {
	cond := c.Expression(0) // can be a block. this could be confusing.
	then := c.Expression(0) // most likely, a block.
	els := c.createIfTokenOrEmpty(KW_ELSE, func(tk TokenPos) NodePosition {
		return c.Expression(0)
	})

	node := c.createNodeFromToken(tk, NODE_IF)
	c.setNodeChildren(node, cond, then, els)
	return node
}

/////////////////////////////////////////////////////
// Special handling for fn
func parseFn(c *nodeBuilder, tk TokenPos, _ int) NodePosition {

	name := c.createIfTokenOrEmpty(TK_ID, func(tk TokenPos) NodePosition {
		return c.createIdNode(tk)
	})

	args := c.createNodeFromCurrentToken(NODE_ARGS)
	app := c.appender(args)
	c.expect(TK_LPAREN)
	for !c.currentTokenIs(TK_RPAREN, TK_ARROW, TK_FATARROW) {
		arg := parseVar(c, c.current, 0)
		if args != 0 {
			app.append(arg)
		}
		c.consume(TK_COMMA)
		// test for comma presence
	}
	c.expect(TK_RPAREN)

	defarrow := c.createIfToken(TK_FATARROW, func(tk TokenPos) NodePosition {
		return c.Expression(0)
	})

	if defarrow != 0 {
		// this is a lambda function where the return type is to be inferred.
		// it also has a body
		// FIXME what about the generics ????
		return c.createNodeFromToken(tk, NODE_DECL_FN, name, args, c.createEmptyNode(), defarrow) // empty type
	}

	rettype := c.createIfTokenOrEmpty(TK_ARROW, func(tk TokenPos) NodePosition {
		return c.Expression(0)
	})

	var blk NodePosition
	if c.currentTokenIs(TK_LBRACE) {
		c.advance()
		blk = parseBlock(c, c.current, 0)
	} else {
		blk = c.createEmptyNode()
	}

	return c.createNodeFromToken(tk, NODE_DECL_FN, name, args, rettype, blk)
}

// parseBlock parses a block of code
func parseBlock(b *nodeBuilder, tk TokenPos, _ int) NodePosition {
	blk := b.createNodeFromCurrentToken(NODE_BLOCK)
	app_blk := b.appender(blk)

	for !b.currentTokenIs(TK_RBRACKET) {
		for b.consume(TK_SEMICOLON) != 0 {
			// advance as much as we can if we have semi colons in the input
		}

		if b.isEof() {
			break
		}

		app_blk.append(b.Expression(0))
	}

	if tk := b.expect(TK_RBRACKET); tk != 0 {
		b.extendRangeFromToken(blk, b.current-1) // FIXME, this is ugly
	}

	return blk
}

// parseTemplate parses a template declaration
// it begins with a list of arguments with optional default values
func parseTemplate(c *nodeBuilder, tk TokenPos, _ int) NodePosition {

	// tup := c.Expression(0).EnsureTuple()

	// // ensure args is a tuple containing variable declarations.
	// args := tup.ToVars()
	// tpl.SetArgs(args)

	// // Where clause would come here, most likely

	// templated := c.Expression(0)
	// // log.Printf("%s (%T)", templated.GetText(), templated)
	// switch v := templated.(type) {
	// case *FnDef:
	// 	v.SetTemplate(tpl)
	// case *FnDecl:
	// 	v.EnsureFnDef(func(f *FnDef) { f.SetTemplate(tpl) })
	// case *TypeDecl:
	// 	v.SetTemplate(tpl)
	// default:
	// 	templated.ReportError("template blocks must be followed by 'fn' or 'type'")
	// }

	// return templated
	return 0
}

// parseTypeDecl parses a type declaration
func parseTypeDecl(c *nodeBuilder, tk TokenPos, _ int) NodePosition {
	name := c.createAndExpectOrEmpty(TK_ID, func(tk TokenPos) NodePosition {
		return c.createIdNode(tk)
	})

	if c.consume(KW_IS) == 0 {
		c.reportErrorAtToken(c.current, `expected 'is' after type declaration`)
	}

	// there might be a pipe here. We don't have to parse a union afterwards because
	// if there is only one type, it doesn't matter.
	c.consume(TK_PIPE)

	typdef := c.Expression(0)
	// raise an error if there is no typedef ?

	typdecl := c.createNodeFromToken(tk, NODE_DECL_TYPE)
	c.setNodeChildren(typdecl, name, typdef)
	return typdecl
}

func parseSemiColon(c *nodeBuilder, tk TokenPos, left NodePosition) NodePosition {
	// ??? this probably shouldn't happen ?
	return 0
}

func parseStruct(c *nodeBuilder, tk TokenPos, _ int) NodePosition {
	stru := c.createNodeFromToken(tk, NODE_STRUCT)
	c.expect(TK_LPAREN)

	app := c.appender(stru)
	for !c.currentTokenIs(TK_RPAREN) {
		// parse a var declaration
		v := parseVar(c, c.current, 0)
		// FIXME ensure a field has a type declaration, as struct
		// fields should have them whether they have defaults or not
		app.append(v)
	}

	if tk := c.expect(TK_RPAREN); tk != 0 {
		c.extendRangeFromToken(stru, tk)
	}

	return stru
}

// parse a variable statement, but also a variable declaration inside
// an argument list of a function signature
func parseVar(c *nodeBuilder, tk TokenPos, rbp int) NodePosition {
	// first, try to scan the ident
	// this may fail, for dubious reasons
	var ident NodePosition
	if tkident := c.expect(TK_ID); tkident != 0 {
		ident = c.createIdNode(tkident)
	} else {
		ident = c.createEmptyNode()
	}

	// there might be a type expression right after the name declaration
	typenode := c.createIfTokenOrEmpty(TK_COLON, func(tk TokenPos) NodePosition {
		// We scan above '=' level to avoid eating it
		return c.ExpressionTokenRbp(TK_EQ)
	})

	// default value !
	var expnode NodePosition
	if c.consume(TK_EQ) != 0 {
		expnode = c.Expression(0)
	} else {
		expnode = c.createEmptyNode()
	}

	varnode := c.createNodeFromToken(tk, NODE_DECL_VAR)
	c.setNodeChildren(varnode, ident, typenode, expnode)
	return varnode
	// Try to parse VAR ourselves
}
