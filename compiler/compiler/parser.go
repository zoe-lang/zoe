package zoe

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

	// { , a block
	nud(TK_LBRACKET, parseBlock)

	// The doc comment forwards the results but sets itself first on the node that resulted
	// Doc comments whose next meaningful token are other doc comments or the end of the file
	// are added at the module level
	nud(TK_DOCCOMMENT, parseDocComment)

	// nud(KW_FOR, parseFor)
	nud(KW_IF, parseIf)

	nud(KW_VAR, func(c *ZoeContext, tk *Token, rbp int) Node {
		res := c.Expression(0)
		if tup, ok := res.(*Tuple); ok {
			vars := tup.ToVars()
			frag := tk.CreateFragment()
			for _, v := range vars {
				frag.AddChildren(v)
			}
			return frag
		}
		v := coerceToVar(res)
		if v != nil {
			return v
		}
		res.ReportError(`var must be followed by a variable declaration`)
		return res
		// Try to parse VAR ourselves
	})

	// nud(KW_IMPORT, parseImport)

	lbp += 2

	unary(KW_RETURN) // FIXME should check for ('}' / 'else' / '|')

	lbp += 2

	nud(KW_TEMPLATE, parseTemplate)

	nud(KW_TYPE, parseTypeDef)

	// ; is a separator that creates a fragment
	lbp_semicolon = lbp
	led(TK_SEMICOLON, parseSemiColon)

	lbp += 2

	// , creates a tuple
	comma_lbp := lbp
	led(TK_COMMA, func(c *ZoeContext, tk *Token, left Node) Node {
		right := c.Expression(comma_lbp) // this should give us a right associative tree
		switch v := left.(type) {
		case *Tuple:
			return v.AddChildren(right)
		default:
			return tk.CreateTuple().AddChildren(left, v)
		}
	})

	lbp += 2
	lbp_equal = lbp

	// =
	binary(TK_EQ)
	binary(KW_IS)

	// fn eats up the expression right next to it
	nud(KW_FN, parseFn)

	lbp += 2

	binary(TK_COLON)
	unary(KW_LOCAL)
	unary(KW_CONST)

	lbp += 2

	binary(TK_LT)
	lbp_gt = lbp
	binary(TK_GT)
	unary(TK_ELLIPSIS)

	lbp += 2

	// |
	nud(TK_PIPE, func(c *ZoeContext, tk *Token, lbp int) Node {
		right := c.Expression(lbp - 1)
		if union, ok := right.(*Union); ok {
			return union
		}
		tk.Context.reportError(tk, `a '|' must always lead an union`)
		return right
	})

	led(TK_PIPE, func(c *ZoeContext, tk *Token, left Node) Node {
		right := c.Expression(lbp)
		if v, ok := left.(*Union); ok {
			return v.AddTypeExprs(right)
		}
		return tk.CreateUnion().AddTypeExprs(left, right)
	})

	lbp += 2

	unary(TK_PLUS)
	unary(TK_MIN)
	binary(TK_MIN)
	binary(TK_PLUS)

	lbp += 2

	binary(TK_STAR)
	binary(TK_DIV)

	lbp += 2

	binary(KW_IS)

	lbp += 2

	handleParens()
	// When used right next to an expression, then paren is a function call
	// handleParens(NODE_LIST, NODE_FNCALL, TK_LPAREN, TK_RPAREN, true)

	lbp += 2

	// surrounding(NODE_BLOCK, NODE_BLOCK, TK_LBRACKET, TK_RBRACKET, false)

	// the index operator
	// nud(TK_LBRACE, parseLbraceNud)

	lbp += 2

	led(TK_FATARROW, parseFnFatArrow)
	// binary(NODE_FNDEF, TK_FATARROW)

	lbp += 2

	led(TK_ARROW, parseArrow)

	lbp += 2

	binary(TK_COLCOL)
	binary(TK_DOT)
	binary(TK_AT)

	binary(KW_AS)

	lbp += 2
	// all the terminals. Lbp was raised, but this is not necessary

	// nud(TK_QUOTE, parseQuote)

	terminal(TK_NUMBER, func(tk *Token) Node {
		return tk.CreateInteger()
	})

	terminal(TK_RAWSTR, func(tk *Token) Node {
		return tk.CreateString()
	})

	terminal(TK_ID, func(tk *Token) Node {
		return tk.CreateIdent()
	})

}

///////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////

var syms = make([]prattTk, TK__MAX) // Far more than necessary

func nudError(c *ZoeContext, tk *Token, rbp int) Node {
	c.reportError(tk, `unexpected '`, tk.String(), `'`)
	return c.Expression(rbp)
}

func ledError(c *ZoeContext, tk *Token, left Node) Node {
	c.reportError(tk, `unexpected '`, tk.String(), `'`)
	return left
}

///////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////

// Group expression between two handleParens tokens
// nk : the node of the list returned when nud
// lednk : the node of the list returned when led
// opening : the opening token
// closing : the closing token
// reduce : whether it is allowed to reduce the list to the single expression
func handleParens() {
	s := &syms[TK_LPAREN]
	s.lbp = lbp

	//
	s.led = func(c *ZoeContext, tk *Token, left Node) Node {
		res := tk.CreateTuple()
		for !c.Consume(TK_RPAREN) {
			if c.isEof() {
				c.reportErrorAtCurrentPosition(`unexpected end of file`)
				// contents = append(contents, c.EOF())
				break
			}
			res.AddChildren(c.Expression(0))
		}

		// On our left, we may have an FnDef or FnDecl, in which case we will graft
		// ourselves on them
		switch v := left.(type) {
		case *FnDecl:
			v.FnDef.Signature.AddArgs(res.ToVars()...)
		case *FnDef:
			v.Signature.AddArgs(res.ToVars()...)
			return v
		}

		return tk.CreateFnCall().SetLeft(left).SetArgs(res)
	}

	//
	s.nud = func(c *ZoeContext, tk *Token, _ int) Node {
		res := tk.CreateTuple()
		for !c.Consume(TK_RPAREN) {
			if c.isEof() {
				c.reportErrorAtCurrentPosition(`unexpected end of file`)
				// contents = append(contents, c.EOF())
				break
			}
			res.AddChildren(c.Expression(0))
		}
		return res
	}
}

// Handle import
// func parseImport(c *ZoeContext, tk *Token, _ int) Node {
// 	if !c.Peek(TK_RAWSTR) {
// 		c.reportErrorAtCurrentPosition(`import expects a raw string as the module name`)
// 	}
// 	name := NewTerminalNode(c.Current)
// 	c.advance()
// 	if c.Consume(KW_AS) {
// 		exp := c.Expression(0)
// 		return NewNode(NODE_IMPORT, tk.Position, name, exp)
// 	}
// 	exp := c.Expression(0)
// 	if !exp.Is(NODE_LIST) {
// 		log.Print(exp.Kind, " - ", exp.String(), exp.Token.Debug(), " - ")
// 		return NewNode(NODE_IMPORT, tk.Position, name, NewNode(NODE_LIST, tk.Position, exp))
// 	}
// 	return NewNode(NODE_IMPORT, tk.Position, name, exp)
// }

///////////////////////////////////////////////////////
// "
// func parseQuote(c *ZoeContext, tk *Token, _ int) Node {
// 	// this should transform the result to a string
// 	return parseUntil(c, NODE_STR, tk, TK_QUOTE, 0)
// }

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
func parseArrow(c *ZoeContext, _ *Token, left Node) Node {
	// left contains the fndef or fndecl
	right := c.Expression(0)
	var block *Block
	var ok bool

	if c.Peek(TK_LBRACKET) {
		bk := c.Current
		if block, ok = c.Expression(0).(*Block); !ok {
			c.reportError(bk, `expected a block`)
		}
	}

	// left is necessarily a tuple. any other type is an error
	switch v := left.(type) {
	case *FnDecl:
		v.FnDef.Signature.SetReturnTypeExp(right)
		if block != nil {
			v.FnDef.SetDefinition(block)
		}
		return v
	case *FnDef:
		v.Signature.SetReturnTypeExp(right)
		if block != nil {
			v.SetDefinition(block)
		}
		return v
	default:
		left.ReportError(`the left side of '->' must be a function definition`)
		return right
	}
}

func parseFnFatArrow(c *ZoeContext, tk *Token, left Node) Node {
	// left is a list of arguments
	// right of => is the implementation of the function

	impl := c.Expression(0) // it is a block or a single expression
	var block *Block
	var ok bool
	if block, ok = impl.(*Block); !ok {
		block = tk.CreateBlock().AddChildren(impl)
	}

	switch v := left.(type) {
	case *FnDef:
		v.SetDefinition(block)
		return v
	}

	left.ReportError(`left hand side of '=>' must be a lambda function definition`)
	return left
}

/////////////////////////////////////////////////////
// Special handling for if block
func parseIf(c *ZoeContext, tk *Token, _ int) Node {
	cond := c.Expression(0)
	then := c.Expression(0)
	node := tk.CreateIf().SetCond(cond).SetThen(then)
	if c.Consume(KW_ELSE) {
		els := c.Expression(0)
		return node.SetElse(els)
	}
	return node
}

/////////////////////////////////////////////////////
// Special handling for fn
func parseFn(c *ZoeContext, tk *Token, _ int) Node {

	// fake_signature := NewNode(NODE_SIGNATURE, tk.Position)
	// Should we create a fake fndef and signature ?
	def := tk.CreateFnDef()

	if c.Peek(TK_ID) {
		id := c.Current.CreateIdent()
		c.advance()

		return tk.CreateFnDecl().SetIdent(id).SetFnDef(def)
	}

	return def //NewNode(NODE_FNDEF, tk.Position, c.Expression(0))

}

func parseBlock(c *ZoeContext, tk *Token, _ int) Node {
	contents := make([]Node, 0)

	for !c.Peek(TK_RBRACKET) {
		if c.isEof() {
			break
		}
		contents = append(contents, c.Expression(0))
	}

	res := tk.CreateBlock().AddChildren(contents...)

	if tk := c.Expect(TK_RBRACKET); tk != nil {
		res.ExtendPosition(tk)
	}

	return res
}

func parseDocComment(c *ZoeContext, tk *Token, rbp int) Node {
	next := c.Expression(rbp)
	nt := tk.NextMeaningfulToken()
	if nt == nil || nt.Is(TK_DOCCOMMENT) {
		c.RootDocComments = append(c.RootDocComments, tk)
	} else {
		c.DocCommentMap[next] = tk
	}
	return next
}

func parseTemplate(c *ZoeContext, tk *Token, _ int) Node {
	tpl := &Template{}
	tpl.ExtendPosition(tk)

	targs := c.Expression(0)

	// ensure args is a tuple containing variable declarations.
	if tup, ok := targs.(*Tuple); ok {
		args := tup.ToVars()

		if len(args) > 0 {
			// FIXME verify the nomenclature of the Idents ($T, $expr, ...)
			// FIXME verify the variables do not have a type.
			tpl.AddArgs(args...)
		}
		// We have our template arguments, which is dandy
		// Template arguments only accept = for default, so we're going to check that
		// there is no type expression.
	} else {
		targs.ReportError(`template should be followed by a list of template arguments`)
	}

	// Where clause would come here, most likely

	templated := c.Expression(0)
	switch v := templated.(type) {
	case *FnDef:
		v.Template = tpl
	case *FnDecl:
		v.FnDef.Template = tpl
	case *TypeDecl:
		v.Template = tpl
	default:
		templated.ReportError("template blocks must be followed by 'fn' or 'type'")
	}

	return templated
}

func parseTypeDef(c *ZoeContext, tk *Token, _ int) Node {
	name := c.Expect(TK_ID)

	if !c.Consume(KW_IS) {
		c.reportErrorAtCurrentPosition(`expected 'is' after type name`)
	}

	typdef := c.Expression(0)
	if name == nil {
		typdef.SetError() // Set error flag
		return typdef
	}

	return tk.CreateTypeDecl().SetIdent(name.CreateIdent()).SetDef(typdef)
}

func parseSemiColon(c *ZoeContext, tk *Token, left Node) Node {
	right := c.Expression(lbp_semicolon) // this should give us a right associative tree
	switch v := left.(type) {
	case *Fragment:
		return v.AddChildren(right)
	default:
		return tk.CreateFragment().AddChildren(left, v)
	}
}
