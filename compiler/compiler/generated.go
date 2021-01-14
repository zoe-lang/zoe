// Code generated by a lame .js file, DO NOT EDIT.

package zoe



func (parser *Parser) createAstFile() *AstFile {
  var res = &AstFile{}
  res.nodeBase = parser.createNodeBase()
  return res
}



func (parser *Parser) createAstImport() *AstImport {
  var res = &AstImport{}
  res.nodeBase = parser.createNodeBase()
  return res
}



func (n *AstImport) GetName() *AstIdentifier {
  return n.named.GetName()
}


func (parser *Parser) createAstImportModuleName() *AstImportModuleName {
  var res = &AstImportModuleName{}
  res.nodeBase = parser.createNodeBase()
  return res
}



func (parser *Parser) createAstVarDecl() *AstVarDecl {
  var res = &AstVarDecl{}
  res.nodeBase = parser.createNodeBase()
  return res
}



func (n *AstVarDecl) GetName() *AstIdentifier {
  return n.named.GetName()
}


func (parser *Parser) createAstNamespaceDecl() *AstNamespaceDecl {
  var res = &AstNamespaceDecl{}
  res.nodeBase = parser.createNodeBase()
  return res
}



func (n *AstNamespaceDecl) GetName() *AstIdentifier {
  return n.named.GetName()
}


func (parser *Parser) createAstImplement() *AstImplement {
  var res = &AstImplement{}
  res.nodeBase = parser.createNodeBase()
  return res
}



func (parser *Parser) createAstTemplateParam() *AstTemplateParam {
  var res = &AstTemplateParam{}
  res.nodeBase = parser.createNodeBase()
  return res
}



func (n *AstTemplateParam) GetName() *AstIdentifier {
  return n.named.GetName()
}


func (parser *Parser) createAstEnumDecl() *AstEnumDecl {
  var res = &AstEnumDecl{}
  res.nodeBase = parser.createNodeBase()
  return res
}



func (n *AstEnumDecl) GetName() *AstIdentifier {
  return n.named.GetName()
}


func (parser *Parser) createAstUnionDecl() *AstUnionDecl {
  var res = &AstUnionDecl{}
  res.nodeBase = parser.createNodeBase()
  return res
}



func (n *AstUnionDecl) GetName() *AstIdentifier {
  return n.named.GetName()
}


func (parser *Parser) createAstStructDecl() *AstStructDecl {
  var res = &AstStructDecl{}
  res.nodeBase = parser.createNodeBase()
  return res
}



func (n *AstStructDecl) GetName() *AstIdentifier {
  return n.named.GetName()
}


func (parser *Parser) createAstFn() *AstFn {
  var res = &AstFn{}
  res.nodeBase = parser.createNodeBase()
  return res
}



func (n *AstFn) GetName() *AstIdentifier {
  return n.named.GetName()
}


func (parser *Parser) createAstBlock() *AstBlock {
  var res = &AstBlock{}
  res.nodeBase = parser.createNodeBase()
  return res
}



func (parser *Parser) createAstFnCall() *AstFnCall {
  var res = &AstFnCall{}
  res.nodeBase = parser.createNodeBase()
  return res
}



func (parser *Parser) createAstIndexCall() *AstIndexCall {
  var res = &AstIndexCall{}
  res.nodeBase = parser.createNodeBase()
  return res
}



func (parser *Parser) createAstDerefOp() *AstDerefOp {
  var res = &AstDerefOp{}
  res.nodeBase = parser.createNodeBase()
  return res
}



func (parser *Parser) createAstPointerOp() *AstPointerOp {
  var res = &AstPointerOp{}
  res.nodeBase = parser.createNodeBase()
  return res
}



func (parser *Parser) createAstReturnOp() *AstReturnOp {
  var res = &AstReturnOp{}
  res.nodeBase = parser.createNodeBase()
  return res
}



func (parser *Parser) createAstTakeOp() *AstTakeOp {
  var res = &AstTakeOp{}
  res.nodeBase = parser.createNodeBase()
  return res
}



func (parser *Parser) createAstMulBinOp() *AstMulBinOp {
  var res = &AstMulBinOp{}
  res.nodeBase = parser.createNodeBase()
  return res
}



func (parser *Parser) createAstDivBinOp() *AstDivBinOp {
  var res = &AstDivBinOp{}
  res.nodeBase = parser.createNodeBase()
  return res
}



func (parser *Parser) createAstAddBinOp() *AstAddBinOp {
  var res = &AstAddBinOp{}
  res.nodeBase = parser.createNodeBase()
  return res
}



func (parser *Parser) createAstSubBinOp() *AstSubBinOp {
  var res = &AstSubBinOp{}
  res.nodeBase = parser.createNodeBase()
  return res
}



func (parser *Parser) createAstModBinOp() *AstModBinOp {
  var res = &AstModBinOp{}
  res.nodeBase = parser.createNodeBase()
  return res
}



func (parser *Parser) createAstPipeBinOp() *AstPipeBinOp {
  var res = &AstPipeBinOp{}
  res.nodeBase = parser.createNodeBase()
  return res
}



func (parser *Parser) createAstAmpBinOp() *AstAmpBinOp {
  var res = &AstAmpBinOp{}
  res.nodeBase = parser.createNodeBase()
  return res
}



func (parser *Parser) createAstLShiftBinOp() *AstLShiftBinOp {
  var res = &AstLShiftBinOp{}
  res.nodeBase = parser.createNodeBase()
  return res
}



func (parser *Parser) createAstRShiftBinOp() *AstRShiftBinOp {
  var res = &AstRShiftBinOp{}
  res.nodeBase = parser.createNodeBase()
  return res
}



func (parser *Parser) createAstAndBinOp() *AstAndBinOp {
  var res = &AstAndBinOp{}
  res.nodeBase = parser.createNodeBase()
  return res
}



func (parser *Parser) createAstOrBinOp() *AstOrBinOp {
  var res = &AstOrBinOp{}
  res.nodeBase = parser.createNodeBase()
  return res
}



func (parser *Parser) createAstGtBinOp() *AstGtBinOp {
  var res = &AstGtBinOp{}
  res.nodeBase = parser.createNodeBase()
  return res
}



func (parser *Parser) createAstGteBinOp() *AstGteBinOp {
  var res = &AstGteBinOp{}
  res.nodeBase = parser.createNodeBase()
  return res
}



func (parser *Parser) createAstLtBinOp() *AstLtBinOp {
  var res = &AstLtBinOp{}
  res.nodeBase = parser.createNodeBase()
  return res
}



func (parser *Parser) createAstLteBinOp() *AstLteBinOp {
  var res = &AstLteBinOp{}
  res.nodeBase = parser.createNodeBase()
  return res
}



func (parser *Parser) createAstEqBinOp() *AstEqBinOp {
  var res = &AstEqBinOp{}
  res.nodeBase = parser.createNodeBase()
  return res
}



func (parser *Parser) createAstNeqBinOp() *AstNeqBinOp {
  var res = &AstNeqBinOp{}
  res.nodeBase = parser.createNodeBase()
  return res
}



func (parser *Parser) createAstIsBinOp() *AstIsBinOp {
  var res = &AstIsBinOp{}
  res.nodeBase = parser.createNodeBase()
  return res
}



func (parser *Parser) createAstIsNotBinOp() *AstIsNotBinOp {
  var res = &AstIsNotBinOp{}
  res.nodeBase = parser.createNodeBase()
  return res
}



func (parser *Parser) createAstDotBinOp() *AstDotBinOp {
  var res = &AstDotBinOp{}
  res.nodeBase = parser.createNodeBase()
  return res
}



func (parser *Parser) createAstNone() *AstNone {
  var res = &AstNone{}
  res.nodeBase = parser.createNodeBase()
  return res
}



func (parser *Parser) createAstTrue() *AstTrue {
  var res = &AstTrue{}
  res.nodeBase = parser.createNodeBase()
  return res
}



func (parser *Parser) createAstFalse() *AstFalse {
  var res = &AstFalse{}
  res.nodeBase = parser.createNodeBase()
  return res
}



func (parser *Parser) createAstIntLiteral() *AstIntLiteral {
  var res = &AstIntLiteral{}
  res.nodeBase = parser.createNodeBase()
  return res
}



func (parser *Parser) createAstStringLiteral() *AstStringLiteral {
  var res = &AstStringLiteral{}
  res.nodeBase = parser.createNodeBase()
  return res
}



func (parser *Parser) createAstThisLiteral() *AstThisLiteral {
  var res = &AstThisLiteral{}
  res.nodeBase = parser.createNodeBase()
  return res
}



func (parser *Parser) createAstNoneLiteral() *AstNoneLiteral {
  var res = &AstNoneLiteral{}
  res.nodeBase = parser.createNodeBase()
  return res
}



func (parser *Parser) createAstIdentifier() *AstIdentifier {
  var res = &AstIdentifier{}
  res.nodeBase = parser.createNodeBase()
  return res
}



func (parser *Parser) createAstStringExp() *AstStringExp {
  var res = &AstStringExp{}
  res.nodeBase = parser.createNodeBase()
  return res
}



func (parser *Parser) createAstIf() *AstIf {
  var res = &AstIf{}
  res.nodeBase = parser.createNodeBase()
  return res
}



func (parser *Parser) createAstWhile() *AstWhile {
  var res = &AstWhile{}
  res.nodeBase = parser.createNodeBase()
  return res
}


