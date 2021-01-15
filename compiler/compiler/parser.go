package zoe

type nudNode interface {
	Node
	nud(iter *Parser, scope *Scope)
}

type ledNode interface {
	Node
	led(iter *Parser, scope *Scope, left Node)
}

/*
	Parse parses the file.
*/
func (f *File) Parse() {
	// /var file = AstFileNew(f)
	// f.RootNode = file

	var scope = &Scope{Names: make(Names), Context: scopeFile}
	f.RootScope = scope

	var parser = Parser{
		pos:  0,
		prev: 0,
		file: f,
	}
	var root = parser.createAstNamespaceDecl(scope)
	f.RootNode = root

	if parser.isSkippable() {
		parser.Advance()
		parser.prev = parser.pos
	}

	parser.parseUntilEof(func() {
		var node = parser.Expression(scope, 0)
		root.Register(node, scope)
	})

}

type bindingPower struct {
	left  int
	right int
}

var rbps [TK__MAX]bindingPower
var rbp = 2 // the base rbp

func prefix(tk TokenKind) {
	rbps[tk] = bindingPower{right: rbp}
}

func suffix(tk TokenKind) {
	rbps[tk] = bindingPower{left: rbp}
}

func leftAssoc(tk TokenKind) {
	rbps[tk] = bindingPower{left: rbp, right: rbp}
}

func rightAssoc(tk TokenKind) {
	rbps[tk] = bindingPower{left: rbp, right: rbp - 1}
}

// __ augments the priority
func __() {
	rbp = rbp + 2
}

/*
	Set up the operator precedence.
*/
func init() {

	leftAssoc(TK_EQ)          //   =
	leftAssoc(TK_PLUSEQ)      //   +=
	leftAssoc(TK_STAREQ)      //   *=
	leftAssoc(TK_MINEQ)       //   -=
	leftAssoc(TK_DIVEQ)       //   /=
	leftAssoc(TK_MODEQ)       //   %=
	leftAssoc(TK_PIPEEQ)      //   |=
	leftAssoc(TK_AMPEQ)       //   &=
	leftAssoc(TK_QUESTIONEQ)  //   ?=
	__()                      //
	leftAssoc(KW_IS)          //   is
	leftAssoc(KW_ISNOT)       //   is not
	__()                      //
	leftAssoc(TK_EQEQ)        //   ==
	leftAssoc(TK_NOTEQ)       //   !=
	__()                      //
	leftAssoc(TK_EXCLAM)      //   !
	__()                      //
	leftAssoc(TK_PIPEPIPE)    //   ||
	leftAssoc(TK_AMPAMP)      //   &&
	__()                      //
	leftAssoc(TK_LT)          //   <
	leftAssoc(TK_LTE)         //   <=
	leftAssoc(TK_GT)          //   >
	leftAssoc(TK_GTE)         //   >=
	__()                      //
	leftAssoc(TK_PIPE)        //   |
	__()                      //
	leftAssoc(TK_RSHIFT)      //   >>
	leftAssoc(TK_LSHIFT)      //   <<
	__()                      //
	leftAssoc(TK_PLUS)        //   +
	leftAssoc(TK_MIN)         //   -
	__()                      //
	leftAssoc(TK_STAR)        //   *
	leftAssoc(TK_DIV)         //   /
	leftAssoc(TK_MOD)         //   %
	__()                      //
	suffix(TK_PLUSPLUS)       //   ++
	suffix(TK_MINMIN)         //   --
	__()                      //
	leftAssoc(TK_AT)          //   @    // it is not really left assoc, it's just to set it at that prio
	__()                      //
	suffix(TK_LPAREN)         //   function call, only suffix since nud is parenthesized
	__()                      //
	leftAssoc(TK_COLCOL)      //   ::
	__()                      //
	suffix(TK_LBRACE)         //   index expression, only suffix since nud is array litteral
	__()                      //
	leftAssoc(TK_DOT)         //   .
	leftAssoc(TK_QUESTIONDOT) //   ?.
}

/*
	Expression is the Pratt parser expression function.
	It implements a very simple algorithm for operator precedence resolution
	based on the concept of right and left "binding power".
*/
func (parser *Parser) Expression(scope *Scope, rbp int) Node {

	if parser.IsEof() {
		parser.reportError("unexpected end of file")
		return nil
	}

	var left = parser.Nud(scope, rbp)

	// nud might have advanced without us knowing...
	if parser.IsEof() {
		return left
	}

	// var next_sym = tk.sym()

	for rbp < parser.Lbp() {
		// log.Print(c.Current.KindStr(), c.Current.Value(c.data))
		left = parser.Led(scope, left)

		if parser.IsEof() {
			return left
		}
	}

	return left
}

/*
	Nud chooses the right method depending on the token
*/
func (parser *Parser) Nud(scope *Scope, rbp int) Node {
	var (
		node  nudNode
		start = parser.pos
	)
	parser.binding = rbps[int(parser.Kind())]

	switch parser.Kind() {

	case KW_THIS:
		node = parser.createAstThisLiteral(scope)

	case KW_VOID:
		node = parser.createAstVoidLiteral(scope)

	case TK_QUOTE:
		node = parser.createAstStringExp(scope)

	case TK_RAWSTR:
		node = parser.createAstStringLiteral(scope)

	case TK_CHAR:
		node = parser.createAstCharLiteral(scope)

	case KW_IF:
		node = parser.createAstIf(scope)

	case KW_ENUM:
		node = parser.createAstEnumDecl(scope)

	case KW_STRUCT:
		node = parser.createAstStructDecl(scope)

	case KW_TRAIT:
		node = parser.createAstTraitDecl(scope)

	case KW_TYPE:
		node = parser.createAstTypeAliasDecl(scope)

	case KW_SWITCH:
		node = parser.createAstSwitch(scope)

	case KW_WHILE:
		node = parser.createAstWhile(scope)

	case KW_IMPORT:
		node = parser.createAstImport(scope)

	case KW_ISO:
		node = parser.createAstIso(scope)

	case TK_DOCCOMMENT:
		var cmtpos = parser.pos
		parser.Advance()
		var nnode = parser.Expression(scope, rbp)
		parser.file.SetComment(nnode, cmtpos)
		return nnode

	case TK_LBRACE:
		// Array lparseral
		// unused for now
		parser.reportError("array literals are not implemented for now")
		parser.Advance()
		return nil

	case KW_FN, KW_METHOD:
		node = parser.createAstFn(scope)

	case KW_VAR, KW_CONST:
		node = parser.createAstVarDecl(scope)

	case TK_LPAREN:
		// parenthesized expression has to be handled differently since the node it returns

	case KW_NAMESPACE:
		node = parser.createAstNamespaceDecl(scope)

	case KW_RETURN:
		node = parser.createAstReturnOp(scope)
	case KW_TAKE:
		node = parser.createAstTakeOp(scope)

	case KW_NONE:
		node = parser.createAstNone(scope)
	case KW_FALSE:
		node = parser.createAstFalse(scope)
	case KW_TRUE:
		node = parser.createAstTrue(scope)

	case TK_NUMBER:
		node = parser.createAstIntLiteral(scope)

	case TK_ID:
		node = parser.createAstIdentifier(scope)

	case TK_LBRACKET:
		node = parser.createAstBlock(scope)

	case TK_AT:
		node = parser.createAstPointerOp(scope)

	default:
		parser.reportError("unexpected token '", parser.GetText(), "'")
		// Do not stay on the error and give a chance to the parser to advance.
		parser.Advance()
		return nil
	}

	if node == nil {
		return node
	}

	// Should extend position of the node !
	node.ExtendPos(start)
	// Advance before nud, so that calls to expression et. al will work.
	node.nud(parser, scope)

	// Assume the parser advanced

	node.ExtendPos(parser.pos)
	return node
}

func (parser *Parser) Led(scope *Scope, left Node) Node {
	var (
		node  ledNode
		start = parser.pos
	)

	switch parser.Kind() {

	case TK_COLCOL:
		node = parser.createAstCastBinOp(scope)

	case TK_AT:
		node = parser.createAstDerefOp(scope)

	case TK_DOT:
		node = parser.createAstDotBinOp(scope)

	case TK_LBRACE:
		node = parser.createAstIndexCall(scope)

	case TK_LPAREN:
		node = parser.createAstFnCall(scope)

	case TK_GT:
		node = parser.createAstGtBinOp(scope)

	case TK_GTE:
		node = parser.createAstGteBinOp(scope)

	case TK_LT:
		node = parser.createAstLtBinOp(scope)

	case TK_LTE:
		node = parser.createAstLteBinOp(scope)

	case TK_EQ, TK_AMPEQ, TK_PIPEEQ, TK_RSHIFTEQ, TK_LSHIFTEQ, TK_DIVEQ, TK_PLUSEQ, TK_STAREQ, TK_MINEQ, TK_MODEQ:
		node = parser.createAstEqBinOp(scope)

	case TK_AMP:
		node = parser.createAstAmpBinOp(scope)
	case TK_PIPE:
		node = parser.createAstPipeBinOp(scope)
	case TK_RSHIFT:
		node = parser.createAstRShiftBinOp(scope)
	case TK_LSHIFT:
		node = parser.createAstLShiftBinOp(scope)

	case TK_DIV:
		node = parser.createAstDivBinOp(scope)
	case TK_PLUS:
		node = parser.createAstAddBinOp(scope)
	case TK_STAR:
		node = parser.createAstMulBinOp(scope)
	case TK_MIN:
		node = parser.createAstSubBinOp(scope)
	case TK_MOD:
		node = parser.createAstModBinOp(scope)

	case TK_AMPAMP:
		node = parser.createAstAndBinOp(scope)
	case TK_PIPEPIPE:
		node = parser.createAstOrBinOp(scope)

	case KW_IS:
		node = parser.createAstIsBinOp(scope)
	case KW_ISNOT:
		node = parser.createAstIsNotBinOp(scope)

	default:
		parser.reportError("unexpected token '", parser.GetText(), "'")
		parser.Advance()
		return left
	}

	node.Extend(left)

	node.led(parser, scope, left)

	if parser.pos > start {
		node.ExtendPos(parser.pos)
	}
	return node
}

func (parser *Parser) Lbp() int {
	parser.binding = rbps[int(parser.Kind())]
	return parser.binding.left
}

/////////////////////////////////////////////////////////////////////////////////
//
//                       NUD CASES
//

/*
	nud for components that don't have to do processing because they're lparserals.
*/
func (n *noopNud) nud(parser *Parser, _ *Scope) {
	parser.Advance()
}

/*
	Identifier
*/
func (id *AstIdentifier) nud(parser *Parser, _ *Scope) {
	parser.Advance()
}

/*
	Return statement
*/
func (ret *AstReturnOp) nud(parser *Parser, scope *Scope) {
	parser.Advance()
	if !parser.IsClosing() {
		ret.Operand = parser.Expression(scope, 0)
	}
}

/*
	All prefix unary operators
*/
func (una *unaryOperation) nud(parser *Parser, scope *Scope) {
	parser.Advance()
	una.Operand = parser.Expression(scope, parser.binding.right)
}

/*
  Parse a block { } when it comes as nud.
*/
func (blk *AstBlock) nud(parser *Parser, scope *Scope) {
	// ???
	var blkscope = scope.subScope(scopeInstructions)
	parser.parseEnclosed(func() {
		var exp = parser.Expression(blkscope, 0)
		blk.Register(exp, blkscope)
	})
}

func (parser *Parser) parseBlock(scope *Scope) *AstBlock {
	var blk = parser.createAstBlock(scope)
	blk.nud(parser, scope)
	blk.ExtendPos(parser.pos)
	return blk
}

/*
	Parse a namespace declaration
*/
func (nm *AstNamespaceDecl) nud(parser *Parser, scope *Scope) {
	parser.Advance()
	// A parent scope should be set...
	var nmscope = scope.subScope(scopeType)

	// Try to parse the identifier
	parser.expect(TK_ID, func() {
		nm.Name = parser.createAstIdentifier(scope)
	})

	if parser.should(TK_LBRACKET) {
		parser.parseEnclosed(func() {
			var xp = parser.Expression(nmscope, 0)
			nm.Register(xp, scope)
		})
	}

}

func (vl *varLike) parseAfterKeyword(parser *Parser, scope *Scope) {
	var start = parser.pos
	if parser.advanceIf(TK_ELLIPSIS) {
		if !scope.isFunctionArguments() {
			parser.reportError("ellipsis is only in function arguments")
		}
		vl.IsEllipsis = true
	}

	parser.expect(TK_ID, func() {
		vl.Name = parser.createAstIdentifier(scope)
	})

	if parser.advanceIf(TK_COLON) {
		vl.TypeExp = parser.Expression(scope, rbps[TK_EQ].right+1)
	}

	if parser.advanceIf(TK_EQ) {
		vl.DefaultExp = parser.Expression(scope, 0)
	}

	if parser.pos == start {
		parser.Advance()
	}
}

/*
	Parse something that looks like a variable declaration.
*/
func (vl *varLike) nud(parser *Parser, scope *Scope) {

	if parser.Is(KW_CONST) {
		vl.IsConst = true
		parser.Advance()
	} else if parser.Is(KW_VAR) {
		parser.Advance()
	}

	vl.parseAfterKeyword(parser, scope)

}

/*
	Parse a function prototype
*/
func (fn *AstFn) nud(parser *Parser, scope *Scope) {

	if parser.Is(KW_METHOD) {
		if !scope.isType() {
			fn.ReportError("methods can only be defined inside of types")
		}
		// This should not happen inside a regular block...
		// Should the error be set here ?
		fn.IsMethod = true
	}
	parser.Advance()

	if parser.Is(TK_ID) {
		fn.parseName(parser, scope)
	}

	var argscope = scope.subScope(scopeArguments)

	fn.parseTemplate(parser, argscope)

	if parser.should(TK_LPAREN) {
		parser.parseEnclosedSeparatedByComma(func() {
			var arg = parser.createAstVarDecl(scope)
			arg.nud(parser, argscope)
			argscope.Add(arg)
			fn.Args = append(fn.Args, arg)
		})
	}

	if parser.advanceIf(TK_ARROW) {
		fn.ReturnType = parser.Expression(scope, 0)
	}

	if parser.Is(TK_LBRACKET) {
		// this is a block
		fn.Definition = parser.Expression(argscope, 0)
	}
}

func (tpl *templated) parseTemplate(parser *Parser, scope *Scope) {
	if !parser.Is(TK_LBRACE) {
		return
	}

	parser.parseEnclosedSeparatedByComma(func() {
		var tpl = parser.createAstTemplateParam(scope)
		tpl.nud(parser, scope)
		scope.Add(tpl)
	})
}

func (tpl *AstTemplateParam) nud(parser *Parser, scope *Scope) {
	if parser.Is(TK_ID) {
		tpl.Name = parser.createAstIdentifier(scope)
	}
	parser.Advance()
}

/*
	Parse an import statement
*/
func (fn *AstImport) nud(parser *Parser, scope *Scope) {
	parser.Advance()

	if parser.Is(TK_RAWSTR) {
		var name = parser.createAstImportModuleName(scope)
		name.ModuleName = parser.GetText() // FIXME we may need to remove quotes !
		fn.Resolver = name
		parser.Advance()
	} else {
		fn.Resolver = parser.Expression(scope, rbps[int(TK_DOT)].right-1) // we don't want to consume parents
	}

	if parser.advanceIf(KW_AS) {
		// import 'whatever' as name
		parser.expect(TK_ID, func() {
			fn.Name = parser.createAstIdentifier(scope)
		})
	} else {
		// import 'whatever' ( symbol, exp as name )
		if !parser.should(TK_LPAREN) {
			return
		}

		parser.parseEnclosedSeparatedByComma(func() {
			var _ = parser.Expression(scope, 0)
			if parser.advanceIf(KW_AS) {
				if parser.should(TK_ID) {
					var _ = parser.createAstIdentifier(scope)
					parser.Advance()
				}
			}
		})
	}
}

func (as *AstTypeDecl) parseTypeDeclBody(parser *Parser, scope *Scope) {
	if !parser.Is(TK_LBRACKET) {
		return
	}

	var typescope = scope.subScope(scopeType)
	// Will parse everything inside the type.
	parser.parseEnclosed(func() {
		_ = parser.Expression(typescope, 0)
	})
}

/*
	Try to parse a name for a named component
*/
func (nam *named) parseName(parser *Parser, scope *Scope) {
	if parser.should(TK_ID) {
		nam.Name = parser.createAstIdentifier(scope)
		parser.Advance()
	}
}

/*
	Type statement
*/
func (typ *AstTypeAliasDecl) nud(parser *Parser, scope *Scope) {
	parser.Advance()
	typ.parseName(parser, scope)
	typ.parseTemplate(parser, scope)

	if parser.Is(TK_LPAREN) {
		parser.parseEnclosedSeparatedByPipe(func() {
			var xp = parser.Expression(scope, 0)
			typ.TypeExps = append(typ.TypeExps, xp)
		})
	}

	if parser.Is(TK_LBRACKET) {
		typ.parseTypeDeclBody(parser, scope)
	}
}

/*
	Parse an enum
*/
func (en *AstEnumDecl) nud(parser *Parser, scope *Scope) {
	parser.Advance()

	en.parseName(parser, scope)
	// enums may not have templates ?

	if parser.should(TK_LPAREN) {
		parser.parseEnclosedSeparatedByComma(func() {
			var v = parser.createAstVarDecl(scope)
			v.parseAfterKeyword(parser, scope)
			en.Fields = append(en.Fields, v)
			en.AddMember(v)
		})
	}

	en.parseTypeDeclBody(parser, scope)
}

/*
	Parse a struct statement
*/
func (st *AstStructDecl) nud(parser *Parser, scope *Scope) {
	parser.Advance()

	st.parseName(parser, scope)
	st.parseTemplate(parser, scope)

	if parser.Is(TK_LPAREN) {
		parser.parseEnclosedSeparatedByComma(func() {
			var v = parser.createAstVarDecl(scope)
			v.parseAfterKeyword(parser, scope)
		})
	}

	st.parseTypeDeclBody(parser, scope)
}

/*
	Parse a trait statement
*/
func (st *AstTraitDecl) nud(parser *Parser, scope *Scope) {
	parser.Advance()

	st.parseName(parser, scope)
	st.parseTemplate(parser, scope)

	if parser.Is(TK_LPAREN) {
		parser.parseEnclosedSeparatedByComma(func() {
			var v = parser.createAstVarDecl(scope)
			v.parseAfterKeyword(parser, scope)
		})
	}

	st.parseTypeDeclBody(parser, scope)
}

/*
	Parse the if statement
*/
func (aif *AstIf) nud(parser *Parser, scope *Scope) {
	parser.Advance()
	var thenscope = scope.subScope(scopeInstructions)
	aif.ConditionExp = parser.Expression(thenscope, 0)

	if parser.should(TK_LBRACKET) {
		aif.ThenArm = parser.parseBlock(thenscope)
	}

	if parser.advanceIf(KW_ELSE) {
		if parser.should(TK_LBRACKET) {
			aif.ElseArm = parser.parseBlock(scope)
		}
	}
}

func (wh *AstWhile) nud(parser *Parser, scope *Scope) {
	parser.Advance()
	wh.ConditionExp = parser.Expression(scope, 0)
	if parser.should(TK_LBRACKET) {
		var blk = parser.createAstBlock(scope)
		blk.nud(parser, scope)
		wh.Body = blk
	}
}

func (swi *AstSwitch) nud(parser *Parser, scope *Scope) {
	parser.Advance()
	var sc = scope.subScope(scopeInstructions)
	swi.ConditionExp = parser.Expression(scope, 0)
	if parser.should(TK_LBRACKET) {
		parser.parseEnclosedSeparatedByComma(func() {
			if parser.Is(KW_ELSE) {
				if swi.ElseArm != nil {
					parser.reportError("redefinition of else statement")
				}
				parser.Advance()
				parser.consume(TK_ARROW)
				swi.ElseArm = parser.Expression(sc, 0)
			} else {
				var arm = parser.createAstSwitchArm(scope)
				arm.ConditionExp = parser.Expression(sc, 0)
				parser.consume(TK_ARROW)
				arm.Body = parser.Expression(sc, 0)
			}
		})
	}
}

func (str *AstStringExp) nud(parser *Parser, scope *Scope) {
	parser.parseEnclosed(func() {
		str.Components = append(str.Components, parser.Expression(scope, 0))
	})
}

/////////////////////////////////////////////////////////////////////////////
//
//                     LED ITEMS
//

/*
	We get here on '('
*/
func (fn *AstFnCall) led(parser *Parser, scope *Scope, left Node) {

	fn.FnExp = left
	parser.parseEnclosedSeparatedByComma(func() {
		fn.Args = append(fn.Args, parser.Expression(scope, 0))
	})
}

/*
  Parse a binary operator
*/
func (bn *binaryOperation) led(parser *Parser, scope *Scope, left Node) {
	parser.Advance()
	bn.Left = left
	bn.Right = parser.Expression(scope, parser.binding.right)
}

/*
	Parse all =, &=, ... assignement nodes, where we don't bother
	creating ast nodes for each of them.
*/
func (eq *AstEqBinOp) led(parser *Parser, scope *Scope, left Node) {
	var kind = parser.Kind()
	var rbp = parser.binding.right
	var rnode binOpNode

	switch kind {
	case TK_AMPEQ:
		rnode = parser.createAstAmpBinOp(scope)
	case TK_PIPEEQ:
		rnode = parser.createAstPipeBinOp(scope)
	case TK_RSHIFTEQ:
		rnode = parser.createAstRShiftBinOp(scope)
	case TK_LSHIFTEQ:
		rnode = parser.createAstLShiftBinOp(scope)

	case TK_DIVEQ:
		rnode = parser.createAstDivBinOp(scope)
	case TK_PLUSEQ:
		rnode = parser.createAstAddBinOp(scope)
	case TK_STAREQ:
		rnode = parser.createAstMulBinOp(scope)
	case TK_MINEQ:
		rnode = parser.createAstSubBinOp(scope)
	case TK_MODEQ:
		rnode = parser.createAstModBinOp(scope)

	default:
		// no rnode !
	}
	parser.Advance()

	var right = parser.Expression(scope, rbp)
	eq.Left = left

	if rnode != nil {
		rnode.SetLeft(left)
		rnode.SetRight(right)
		right = rnode
	}

	eq.Right = right
}

/*
	Suffix operation in led position
*/
func (una *unaryOperation) led(parser *Parser, scope *Scope, left Node) {
	parser.Advance()
	una.Operand = left
}

func (idx *AstIndexCall) led(parser *Parser, scope *Scope, left Node) {
	idx.Indexed = left
	parser.parseEnclosedSeparatedByComma(func() {
		var exp = parser.Expression(scope, 0)
		idx.Indices = append(idx.Indices, exp)
	})
}
