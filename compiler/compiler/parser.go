package zoe

func (f *File) Parse() {
	_, f.RootNode = f.parseFile()
	// control that we got to the last token ???
}

// At the top level, just parse everything we can
func (f *File) parseFile() (Tk, Node) {
	scope := f.RootScope()
	tk := Tk{
		pos:  0,
		file: f,
	}
	if tk.isSkippable() {
		tk = tk.Next()
	}
	file := f.createNode(tk, NODE_FILE, scope)

	app := newList()
	tk.whileNotEof(func(iter Tk) Tk {
		var node Node
		iter, node = Expression(scope, iter, 0)
		if !node.IsEmpty() {
			file.ExtendRangeFromNode(node)
			app.append(node)
		}
		return iter
	})

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

	//
	// Parse a parenthesized expression.
	//
	nud(TK_LPAREN, func(scope Scope, tk Tk, lbp int) (Tk, Node) {
		// We are going to check if we have several components to the paren, or just
		// one, in which case we just send it back.
		// an empty () parenthesis block is an error as it doesn't mean anything.

		iter, exp := Expression(scope, tk.Next(), 0)
		// check if we end with a parenthesis
		if next, ok := iter.consume(TK_RPAREN); ok {
			exp.Extend(tk)
			return next, exp
		}

		// If we didn't encounter ), we want a comma
		iter, _ = iter.expect(TK_COMMA)

		app := newList()
		app.append(exp)

		iter = iter.whileNotClosing(func(iter Tk) Tk {
			iter, exp = Expression(scope, iter, 0)
			if !iter.Is(TK_RPAREN) {
				iter, _ = iter.expect(TK_COMMA) // there can be a comma
			}
			app.append(exp)
			return iter
		})

		tup := tk.createTuple(scope, app.first)
		iter, _ = iter.expect(TK_RPAREN, func(tk Tk) { tup.Extend(tk) })

		return iter, tup
	})

	//
	// Namespace declaration (might go away)
	//
	nud(KW_NAMESPACE, func(scope Scope, tk Tk, lbp int) (Tk, Node) {
		var iter = tk.Next()
		// res.Block = res.CreateBlock()
		var name Node
		iter, name = Expression(scope, iter, 0)

		var block Node
		var nmsp_scope = scope.subScope()

		if _, ok := iter.expect(TK_LBRACKET); ok {
			iter, block = parseBlock(nmsp_scope, iter, 0)
		}

		var nmsp = tk.createNamespace(nmsp_scope, name, block)
		return iter, nmsp
	})

	// { , a block
	nud(TK_LBRACKET, parseBlock)

	nud(TK_LBRACE, func(scope Scope, tk Tk, _ int) (Tk, Node) {
		// function call !
		var iter = tk.Next()

		fragment := newList()
		iter = iter.whileNotClosing(func(iter Tk) Tk {
			var exp Node
			iter, exp = Expression(scope, iter, 0)
			fragment.append(exp)
			if !iter.Is(TK_RBRACE) {
				iter, _ = iter.consume(TK_COMMA)
			}
			return iter
		})

		array := tk.createArrayLiteral(scope, fragment.first)

		iter, _ = iter.expect(TK_RBRACE, func(tk Tk) {
			array.Extend(tk)
		})

		return iter, array
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
	nud(KW_SWITCH, parseSwitch)

	nud(KW_LOCAL, func(scope Scope, tk Tk, lbp int) (Tk, Node) {
		next, v := Expression(scope, tk.Next(), lbp)
		if v.expect(NODE_VAR) {
			v.Extend(tk)
		} else {
			tk.reportError("expected a variable declaration")
		}
		return next, v
	})

	nud(KW_VAR, func(scope Scope, tk Tk, lbp int) (Tk, Node) {
		next, v := parseVar(scope, tk.Next(), lbp)
		if v.expect(NODE_VAR) {
			v.Extend(tk)
		} else {
			tk.reportError("expected a variable declaration")
		}
		return next, v
	})

	nud(KW_CONST, func(scope Scope, tk Tk, lbp int) (Tk, Node) {
		next, va := parseVar(scope, tk.Next(), lbp)
		if !va.IsEmpty() {
			va.SetFlag(FLAG_CONST)
		}
		return next, va
	})

	nud(KW_IMPORT, parseImport)

	nud(KW_IMPLEMENT, func(scope Scope, tk Tk, lbp int) (Tk, Node) {
		var iter = tk.Next()
		var namexp Node
		iter, namexp = Expression(scope, iter, 0)

		var blk Node
		if iter.Is(TK_LBRACKET) {
			iter, blk = parseBlock(scope, iter, 0)
		}

		return iter, tk.createImplement(scope, namexp, blk)
	})

	lbp += 2

	nud(KW_TAKE, func(scope Scope, tk Tk, lbp int) (Tk, Node) {
		var res Node
		iter := tk
		// do not try to get next expression is return is immediately followed
		// by } or ]
		if tk.Peek(TK_RPAREN, TK_RBRACKET) {
			// return can only return nothing if it is at the end of a block or expression
			res = EmptyNode
		} else {
			iter, res = Expression(scope, tk.Next(), lbp)
		}

		return iter, tk.createTake(scope, res)
	})

	// return ...
	// will return an empty node if
	nud(KW_RETURN, func(scope Scope, tk Tk, lbp int) (Tk, Node) {
		var res Node
		iter := tk
		// do not try to get next expression is return is immediately followed
		// by } or ]
		if tk.Peek(TK_RPAREN, TK_RBRACKET) {
			// return can only return nothing if it is at the end of a block or expression
			res = EmptyNode
		} else {
			iter, res = Expression(scope, tk.Next(), lbp)
		}

		return iter, tk.createReturn(scope, res)
	})

	nud(KW_TYPE, parseTypeDecl)
	nud(KW_STRUCT, parseStruct)
	nud(KW_TRAIT, parseTrait)
	nud(KW_ENUM, parseEnum)

	lbp += 2
	lbp_equal = lbp

	// =
	binary(TK_EQ, NODE_BIN_ASSIGN)
	binary(KW_IS, NODE_BIN_IS)

	// fn eats up the expression right next to it
	nud(KW_FN, parseFn)
	nud(KW_METHOD, parseFn)

	lbp += 2

	led(TK_COLON, parseVarEnd)

	// unary(KW_LOCAL)
	// unary(KW_CONST)
	binary(TK_EQEQ, NODE_BIN_EQ)
	binary(TK_NOTEQ, NODE_BIN_NEQ)

	lbp_eq := lbp
	nud(TK_EXCLAM, func(scope Scope, tk Tk, lbp int) (Tk, Node) {
		next, exp := Expression(scope, tk.Next(), lbp_eq)
		node := tk.createUnaNot(scope, exp)
		return next, node
	})

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

	lbp_addition := lbp
	// The + prefix operator, which is essentially a noop
	nud(TK_PLUS, func(scope Scope, tk Tk, lbp int) (Tk, Node) {
		return Expression(scope, tk.Next(), lbp_addition)
	})

	// The - prefix operator, which gets converted as a multiplication by -1
	nud(TK_MIN, func(scope Scope, tk Tk, lbp int) (Tk, Node) {
		var next, exp = Expression(scope, tk.Next(), lbp_addition)
		// create a node for -1
		var min_one = tk.file.createNode(tk, NODE_INTEGER, scope)
		min_one.SetValue(-1) // a forced integer
		bin := tk.createBinOp(scope, NODE_BIN_MUL, min_one, exp)
		return next, bin
	})

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
		var next = tk.Next()
		var one = tk.createNode(scope, NODE_INTEGER)
		one.SetValue(1)
		var addition = tk.createNode(scope, NODE_BIN_PLUS, one, left)
		var assign = tk.createNode(scope, NODE_BIN_ASSIGN, left, addition)
		return next, assign
	})

	led(TK_MINMIN, func(scope Scope, tk Tk, left Node) (Tk, Node) {
		var next = tk.Next()
		var one = tk.createNode(scope, NODE_INTEGER)
		one.SetValue(-1)
		var addition = tk.createNode(scope, NODE_BIN_PLUS, one, left)
		var assign = tk.createNode(scope, NODE_BIN_ASSIGN, left, addition)
		return next, assign
	})

	lbp += 2

	// surrounding(NODE_BLOCK, NODE_BLOCK, TK_LBRACKET, TK_RBRACKET, false)

	// the index operator
	// nud(TK_LBRACE, parseLbraceNud)
	nud(TK_STAR, func(scope Scope, tk Tk, lbp int) (Tk, Node) {
		iter := tk.Next()
		iter, typeexpr := Expression(scope, iter, syms[TK_MINMIN].lbp+1)
		if typeexpr.IsEmpty() {
			tk.reportError("expected * to be followed by a type name")
		}
		return iter, tk.createUnaPointer(scope, typeexpr)
	})

	nud(TK_AMP, func(scope Scope, tk Tk, lbp int) (Tk, Node) {
		iter := tk.Next()
		iter, expr := Expression(scope, iter, syms[TK_MINMIN].lbp+1)
		if expr.IsEmpty() {
			tk.reportError("expected & to be followed by an expression")
		}
		return iter, tk.createUnaRef(scope, expr)
	})

	lbp += 2

	// led(TK_FATARROW, parseFnFatArrow)
	// binary(NODE_FNDEF, TK_FATARROW)

	led(TK_LBRACE, func(scope Scope, tk Tk, left Node) (Tk, Node) {
		var iter = tk.Next()
		var fragment = newList()

		iter = iter.whileNotClosing(func(iter Tk) Tk {
			var exp Node
			iter, exp = Expression(scope, iter, 0)
			fragment.append(exp)
			iter = iter.expectCommaIfNot(TK_RBRACE)
			return iter
		})

		var index = tk.createBinOp(scope, NODE_BIN_INDEX, left, fragment.first)
		iter, _ = iter.expect(TK_RBRACE, func(tk Tk) {
			index.Extend(tk)
		})

		return iter, index
	})

	lbp += 2

	led(TK_LPAREN, func(scope Scope, tk Tk, left Node) (Tk, Node) {
		// function call !
		var iter = tk.Next()
		var fragment = newList()

		iter = iter.whileNotClosing(func(iter Tk) Tk {
			var exp Node
			iter, exp = Expression(scope, iter, 0)

			fragment.append(exp)
			iter = iter.expectCommaIfNot(TK_RPAREN)
			return iter
		})

		var call = tk.createBinOp(scope, NODE_BIN_CALL, left, fragment.first)

		iter, _ = iter.expect(TK_RPAREN, func(tk Tk) {
			call.Extend(tk)
		})

		return iter, call
	})

	lbp += 2

	binary(TK_DOT, NODE_BIN_DOT)
	// binary(KW_AS)

	lbp += 2
	binary(TK_COLCOL, NODE_BIN_CAST)
	lbp += 2

	nud(TK_QUOTE, parseQuote)
	nud(KW_ISO, func(scope Scope, tk Tk, lbp int) (Tk, Node) {
		next := tk.Next()
		if next.Is(TK_LBRACKET) {
			var blk Node
			next, blk = parseBlock(scope, next, 0)
			block := tk.createIsoBlock(scope, blk)
			return next, block
		}
		if next.Is(TK_LBRACE) {
			var exp Node
			next, exp = parseBlock(scope, next, 0)
			iso_expr := tk.createIsoType(scope, exp)
			return next, iso_expr
		}

		next.reportError(`iso expects either a {block} or a [type] expression`)
		// We still advance the parser to make sure that we don't get stuck in a loop
		return next, EmptyNode
	})

	literal(KW_TRUE, NODE_LIT_TRUE)
	literal(KW_FALSE, NODE_LIT_FALSE)
	literal(KW_NONE, NODE_LIT_NONE)
	literal(KW_VOID, NODE_LIT_VOID)
	literal(KW_THIS, NODE_LIT_THIS)
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
	if tk.IsClosing() {
		return tk, EmptyNode
	}
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
	var iter = tk.Next()
	var mod Node
	var ok bool

	if iter.Is(TK_RAWSTR) {
		mod = iter.createNode(scope, NODE_LIT_RAWSTR)
		iter = iter.Next()
	} else {
		iter, mod = Expression(scope, iter, syms[TK_DOT].lbp-1)
	}

	if iter, ok = iter.consume(KW_AS); ok {
		var name Node
		if iter.Is(TK_ID) {
			name = iter.createIdNode(scope)
			iter = iter.Next()
		}

		var imp = tk.createImport(scope, mod, name, EmptyNode)

		if !name.IsEmpty() {
			// Add the import to the current scope.
			scope.addSymbolFromIdNode(name, imp)
		}

		return iter, imp
	}

	if iter, ok = iter.consume(TK_LPAREN); !ok {
		iter.reportError("malformed import expression, expected '(' or 'as'")
		return iter, EmptyNode
	}

	fragment := newList()
	iter = iter.whileNotClosing(func(iter Tk) Tk {

		mod2 := mod.Clone()

		var path Node
		iter, path = Expression(scope, iter, syms[TK_DOT].lbp-1) // we want the tk_dots

		if iter, ok = iter.consume(KW_AS); ok {
			var as Node
			var prev = iter

			if iter.Is(TK_ID) {
				as = iter.createIdNode(scope)
				iter = iter.Next()
			}

			var imp = prev.createImport(scope, mod2, as, path)
			if !as.IsEmpty() {
				scope.addSymbolFromIdNode(as, imp)
			}

			fragment.append(imp)
		} else {
			var id2 = path.Clone()
			var imp = iter.createImport(scope, mod2, id2, path)

			scope.addSymbolFromIdNode(id2, imp)
			fragment.append(imp)
		}
		iter = iter.expectCommaIfNot(TK_RPAREN)
		return iter
	})

	iter, _ = iter.expect(TK_RPAREN)

	return iter, fragment.first
}

///////////////////////////////////////////////////////
// "
func parseQuote(scope Scope, tk Tk, _ int) (Tk, Node) {
	iter := tk.Next()
	fragment := newList()
	iter = iter.whileNot(TK_QUOTE, func(iter Tk) Tk {
		var exp Node
		iter, exp = Expression(scope, iter, 0)
		fragment.append(exp)
		return iter
	})

	str := tk.createString(scope, fragment.first)
	iter, _ = iter.expect(TK_QUOTE, func(tk Tk) {
		str.Extend(tk)
	})

	// this should transform the result to a string
	return iter, str
}

/////////////////////////////////////////////////////
// Special handling for if block
func parseIf(scope Scope, tk Tk, _ int) (Tk, Node) {
	var ok bool
	var has_else bool
	iter := tk.Next()

	// Inside the if condition, we want to know if there were some is operators associated
	// to some identifiers, so that we can build unions for them.
	// How to do that, through the scope ?
	iter, cond := Expression(scope, iter, 0) // can be a block. this could be confusing.

	iter.expect(TK_LBRACKET)

	var thenscope = scope.subScope()
	var elsescope = scope.subScope()

	// We need to check in the then block if it gives back control once it is done
	// or if it stops execution in the current scope. If it does, and there is no else
	// block, then the rest of the parent block is in fact the else condition.
	iter, then := Expression(thenscope, iter, 0) // most likely, a block.

	var els Node
	iter, has_else = iter.consume(KW_ELSE)

	if has_else {
		iter.expect(TK_LBRACKET)
	}

	if iter, ok = iter.consume(KW_ELSE); ok {
		iter.expect(TK_LBRACKET)
		iter, els = Expression(elsescope, iter, 0)
	}

	return iter, tk.createIf(scope, cond, then, els)
}

//
func parseSwitch(scope Scope, tk Tk, _ int) (Tk, Node) {
	var iter = tk.Next()

	// 1. Get the expression we're switching on.
	var switchexp Node
	iter, switchexp = Expression(scope, iter, 0)

	// 2. Parse all the matching arms
	iter, _ = iter.expect(TK_LBRACKET)

	var list = newList()
	iter = iter.whileNotClosing(func(iter Tk) Tk {
		var first = iter
		// iter, _ = iter.expect(TK_PIPE)

		var armexp Node
		var iselse bool
		if !iter.Is(KW_ELSE) {
			iter, armexp = Expression(scope, iter, 0)
		} else {
			iselse = true
			iter = iter.Next()
		}

		iter, _ = iter.expect(TK_ARROW)

		var thenexp Node
		iter, thenexp = Expression(scope, iter, 0)

		if iselse {
			iter, _ = iter.consume(TK_COMMA)
		} else {
			iter, _ = iter.expect(TK_COMMA)
		}

		var arm = first.createSwitchArm(scope, armexp, thenexp)
		list.append(arm)
		return iter
	})

	var sw = tk.createSwitch(scope, switchexp, list.first)

	if iter.shouldBe(TK_RBRACKET) {
		// extend the range of the switch
		sw.Extend(iter)
		iter = iter.Next()
	}
	return iter, sw
}

/////////////////////////////////////////////////////
// Special handling for fn
func parseFn(scope Scope, tk Tk, _ int) (Tk, Node) {

	var fnscope = scope.subScope()

	var iter = tk.Next()

	// Function name, may not exist
	var name Node
	iter, _ = iter.consume(TK_ID, func(tk Tk) {
		name = tk.createIdNode(scope)
	})

	if name != EmptyNode {
		scope.addSymbolFromIdNode(name, name)
	}

	// Template arguments, may not exist
	var tpl Node
	if iter.Is(TK_LBRACE) {
		iter, tpl = parseTemplate(fnscope, iter, 0)
	}

	// Function arguments, mandatory
	var args = newList()
	iter, _ = iter.expect(TK_LPAREN)
	iter = iter.whileNotClosing(func(iter Tk) Tk {

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
		return iter
	})

	iter, _ = iter.expect(TK_RPAREN)

	// Return type, may not exist
	var rettype Node
	if iter.Is(TK_ARROW) {
		iter = iter.Next()
		iter, rettype = Expression(fnscope, iter, 0)
	}

	// The signature node
	signature := tk.createSignature(fnscope, tpl, args.first, rettype)

	var result = signature

	// Function definition
	var blk Node
	if iter.Is(TK_LBRACKET) {
		iter, blk = parseBlock(fnscope, iter, 0)
		// should register the function somewhere in scope, no ?
		result = tk.createFn(fnscope, name, signature, blk)
	}

	if name != EmptyNode {
		var is_method = tk.Is(KW_METHOD)
		if is_method {

		} else {
			scope.addSymbolFromIdNodeForce(name, result)
		}
	}

	return iter, result
}

func parseUntilClosing(scope Scope, tk Tk, _ int) (Tk, Node) {
	var iter = tk
	var fragment = newList()
	iter = iter.whileNotClosing(func(iter Tk) Tk {
		for iter.Is(TK_SEMICOLON) {
			iter = iter.Next()
		}

		if iter.IsEof() {
			return iter
		}

		var exp Node
		iter, exp = Expression(scope, iter, 0)
		fragment.append(exp)
		return iter
	})
	return iter, fragment.first
}

// parseBlock parses a block of code
func parseBlock(scope Scope, tk Tk, _ int) (Tk, Node) {
	var iter = tk.Next()
	var subscope = scope.subScope()

	var first Node
	iter, first = parseUntilClosing(scope, iter, 0)

	block := tk.createBlock(subscope, first)

	iter = iter.expectClosing(tk, func(tk Tk) {
		block.Extend(tk)
	})

	return iter, block
}

// parseTemplate parses a template declaration, which is enclosed between [ ]
// it is expected that '[' has been consumed, and that tk is '['
func parseTemplate(scope Scope, tk Tk, _ int) (Tk, Node) {
	// tpl := b.createNodeFromToken(tk, NODE_TEMPLATE)
	var iter = tk.Next()
	var fragment = newList()

	iter = iter.whileNotClosing(func(iter Tk) Tk {
		iter, node := Expression(scope, iter, 0)
		if node.expect(NODE_VAR, NODE_ID) {
			fragment.append(node)
		}

		iter = iter.expectCommaIfNot(TK_RBRACE)
		return iter
	})
	iter, _ = iter.expect(TK_RBRACE)
	return iter, fragment.first
}

// parseTypeDecl parses a type declaration
func parseTypeDecl(scope Scope, tk Tk, _ int) (Tk, Node) {

	var iter = tk.Next()
	var typescope = scope.subScope()

	var name Node
	iter, _ = iter.expect(TK_ID, func(tk Tk) {
		name = tk.createIdNode(scope)
	})

	var tpl Node
	if iter.Is(TK_LBRACE) {
		iter, tpl = parseTemplate(scope, iter, 0)
	}

	iter, _ = iter.expect(KW_IS)

	// there might be a pipe here. We don't have to parse a union afterwards because
	// if there is only one type, it doesn't matter.
	iter, _ = iter.consume(TK_PIPE)

	var typdef Node
	iter, typdef = Expression(typescope, iter, 0)

	var blk Node
	if iter.Is(TK_LBRACKET) {
		iter, blk = Expression(typescope, iter, 0)
	}

	var typ = tk.createType(typescope, name, tpl, typdef, blk)

	// register the type to the scope
	if name != EmptyNode {
		scope.addSymbolFromIdNode(name, typ)
	}

	return iter, typ
}

//
// TRAIT
//
func parseTrait(scope Scope, tk Tk, _ int) (Tk, Node) {
	var iter = tk.Next()

	var traitscope = scope.subScope()

	var name Node
	iter, _ = iter.expect(TK_ID, func(tk Tk) {
		name = tk.createIdNode(scope)
	})

	var template Node
	if iter.Is(TK_LBRACE) {
		iter, template = parseTemplate(scope, iter, 0)
	}

	var fields Node
	iter, fields = tryParseList(scope, iter, TK_LBRACKET, TK_RBRACKET, TK_COMMA, false, func(scope Scope, iter Tk) (Tk, Node) {

		var node Node
		iter, node = Expression(traitscope, iter, 0)

		return iter, node
	})

	var trait = tk.createTrait(traitscope, template, fields)
	trait.Extend(iter)

	if !name.IsEmpty() {
		scope.addSymbolFromIdNode(name, trait)
	}

	return iter, trait
}

//
//	STRUCT
//
func parseStruct(scope Scope, tk Tk, _ int) (Tk, Node) {
	var iter = tk.Next()

	var strscope = scope.subScope()

	var name Node
	iter, _ = iter.expect(TK_ID, func(tk Tk) {
		name = tk.createIdNode(scope)
	})

	var template Node
	if iter.Is(TK_LBRACE) {
		iter, template = parseTemplate(scope, iter, 0)
	}

	var fields Node
	iter, fields = tryParseList(strscope, iter, TK_LBRACKET, TK_RBRACKET, TK_COMMA, false, func(scope Scope, iter Tk) (Tk, Node) {
		var member Node
		iter, member = Expression(strscope, iter, 0)

		return iter, member
	})

	var stru = tk.createStruct(strscope, template, fields)
	stru.Extend(iter)

	if !name.IsEmpty() {
		scope.addSymbolFromIdNode(name, stru)
	}

	return iter, stru
}

func parseEnum(scope Scope, tk Tk, _ int) (Tk, Node) {
	var iter = tk.Next()

	var strscope = scope.subScope()

	var name Node
	iter, _ = iter.expect(TK_ID, func(tk Tk) {
		name = tk.createIdNode(scope)
	})

	var fields Node
	iter, fields = tryParseList(scope, iter, TK_LBRACKET, TK_RBRACKET, TK_COMMA, false, func(scope Scope, iter Tk) (Tk, Node) {
		var variable Node
		iter, variable = parseVar(strscope, iter, 0)

		if variable.IsEmpty() {
			iter.reportError("expected a field declaration")
		}

		return iter, variable
	})

	var stru = tk.createEnum(scope, fields)
	stru.Extend(iter)

	if !name.IsEmpty() {
		scope.addSymbolFromIdNode(name, stru)
	}

	return iter, stru
}

func parseFor(scope Scope, tk Tk, _ int) (Tk, Node) {
	var iter = tk.Next()
	var forscope = scope.subScope()

	// the "var ..." part of the for loop
	var decl Node
	iter, decl = Expression(forscope, iter, 0)

	iter, _ = iter.expect(KW_IN)

	var inexp Node
	iter, inexp = Expression(scope, iter, 0)

	iter.expect(TK_LBRACKET)

	var block Node
	iter, block = Expression(forscope, iter, 0)

	var fornode = tk.createFor(scope, decl, inexp, block)
	// I should add the subscope to the for node !
	return iter, fornode
}

func parseWhile(scope Scope, tk Tk, _ int) (Tk, Node) {
	var whilescope = scope.subScope()
	var iter = tk.Next()

	var cond Node
	iter, cond = Expression(whilescope, iter, 0)

	iter.expect(TK_LBRACKET)

	var block Node
	iter, block = Expression(whilescope, iter, 0)

	var whilenode = tk.createWhile(scope, cond, block)
	// FIXME add the subscope to the while node
	return iter, whilenode
}

// We got here on :
func parseVarEnd(scope Scope, tk Tk, left Node) (Tk, Node) {
	var iter = tk.Next()
	var typ Node
	iter, typ = Expression(scope, iter, syms[TK_EQ].lbp+1)

	var eql Node
	if iter.Is(TK_EQ) {
		iter, eql = Expression(scope, iter, 0)
	}

	var v = tk.createVar(scope, left, typ, eql)
	return iter, v
}

// parse a variable statement, but also a variable declaration inside
// an argument list of a function signature
func parseVar(scope Scope, tk Tk, _ int) (Tk, Node) {
	var iter = tk
	var ok bool

	// first, try to scan the ident
	// this may fail, for dubious reasons
	var ident Node
	iter, _ = iter.expect(TK_ID, func(tk Tk) {
		ident = tk.createIdNode(scope)
	})

	// An optional type definition
	var typenode Node
	if iter, ok = iter.consume(TK_COLON); ok {
		// there is a type expression
		iter, typenode = Expression(scope, iter, syms[TK_EQ].lbp+1)
	}

	// An optional default value
	var expnode Node
	if iter, ok = iter.consume(TK_EQ); ok {
		iter, expnode = Expression(scope, iter, 0)
	}

	if iter.pos == tk.pos {
		// no variable was found since all of the above expression may fail
		// so we still advance the parser in the hopes of not failing

		return iter.Next(), EmptyNode
	}

	var varnode = tk.createVar(scope, ident, typenode, expnode)
	if !ident.IsEmpty() {
		scope.addSymbolFromIdNode(ident, varnode)
	}

	return iter, varnode
	// Try to parse VAR ourselves
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// A list of nodes
type list struct {
	first Node
	last  Node
}

func newList() list {
	return list{}
}

func (f *list) append(node Node) {

	if node.IsEmpty() {
		return
	}

	if f.first.IsEmpty() {
		f.first = node
		f.last = node
		return
	}

	f.last.SetNext(node)
	for node.HasNext() {
		node = node.Next()
	}

	f.last = node
}
