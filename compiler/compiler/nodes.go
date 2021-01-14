package zoe

import (
	"github.com/sourcegraph/go-lsp"
)

type Node interface {
	Extend(n Node)
	ExtendTk(tk *Parser)
	ExtendPos(pos TokenPos)

	// Get the text
	GetName() *AstIdentifier
	GetChildren() []Node
	GetLspRange() lsp.Range
	GetTkRange() TkRange
	GetBytes() []byte
	GetText() string

	GetScope() *Scope

	// Register another node to this node.
	// This is where a node might decide to send a node to a scope or to its members
	Register(n Node, scope *Scope)

	ReportError(msg ...string)
}

type Names map[Name]Node

type noopNud struct{}

/////////////////////////////////////////////////////////
// Node location within the file

/*
	The base node type
*/
type nodeBase struct {
	File      *File
	Scope     *Scope
	HasErrors bool
	TkRange   TkRange
}

func (l *nodeBase) GetScope() *Scope {
	return l.Scope
}

func (l *nodeBase) GetName() *AstIdentifier {
	return nil
}

func (l *nodeBase) GetChildren() []Node {
	return nil
}

func (l *nodeBase) Extend(n Node) {
	if n == nil {
		return
	}
	l.TkRange.ExtendRange(n.GetTkRange())
}

func (l *nodeBase) ExtendTk(tk *Parser) {
	l.TkRange.ExtendTk(tk)
}

func (l *nodeBase) ExtendPos(pos TokenPos) {
	l.TkRange.ExtendPos(pos)
}

func (l *nodeBase) GetTkRange() TkRange {
	return l.TkRange
}

func (l *nodeBase) Register(_ Node, _ *Scope) {

}

/*
	ReportError reports an error on this node.
	Warning : this is not safe as long as the node is not yet "mounted".
	ReportError should thus be used mostly during typechecking or during
	the Ast creation.
*/
func (l *nodeBase) ReportError(msg ...string) {
	var file = l.GetFile()
	file.reportError(l.GetLspRange(), msg...)
}

/*
  GetFile returns the file this node is associated to
*/
func (l *nodeBase) GetFile() *File {
	return l.File
}

/*
	GetLspRange returns a range as defined by the Language Server Protocol that corresponds to this node.
*/
func (l *nodeBase) GetLspRange() lsp.Range {
	var file = l.GetFile()
	var t = l.GetTkRange()
	var tks = file.Tokens
	var st = tks[int(t.Start)]
	var ed = tks[int(t.End)]
	return lsp.Range{
		Start: lsp.Position{
			Line:      int(st.Line),
			Character: int(st.Column),
		},
		End: lsp.Position{
			Line:      int(ed.Line),
			Character: int(ed.Column),
		},
	}
}

/*
	GetBytes returns a []byte containing the text portion of this node in the source file.
*/
func (l *nodeBase) GetBytes() []byte {
	var file = l.GetFile()
	return file.GetTkRangeBytes(l.TkRange)
}

/*
	GetText returns the string corresponding to this node's text portion in the source file.
*/
func (l *nodeBase) GetText() string {
	var file = l.GetFile()
	return file.GetTkRangeString(l.TkRange)
}

/////////////////////////////////////////////////////////
// Declaration

type declaration struct {
	IsLocal  bool
	IsExtern bool
}

//////////////////////////////////////////////////////////

type varLike struct {
	named
	IsConst    bool
	IsEllipsis bool
	TypeExp    Node
	DefaultExp Node
}

type named struct {
	Name *AstIdentifier
}

func (n *named) GetName() *AstIdentifier {
	return n.Name
}

type templated struct {
	TemplateParams []*AstTemplateParam
}

type membered struct {
	Members Names
}

func (m *membered) AddMember(n Node) {
	if n == nil {
		// ??
		return // ???
	}
	var name = n.GetName().Name
	if _, ok := m.Members[name]; ok {
		n.ReportError("member '" + name.GetText() + "' is already defined")
		return
	}
	m.Members[name] = n

}

///////////////////////////////////////////////////////////

type AstFile struct {
	nodeBase
	membered
	File *File
}

///////////////////////////////////////////////////////////

type AstImport struct {
	nodeBase
	named
	Resolver      Node // Resolver is either a string, which will resolve to a module or a node containing an expression
	SubExpression Node
}

type AstImportModuleName struct {
	nodeBase
	ModuleName string
}

type AstVarDecl struct {
	// Symbol
	nodeBase
	varLike
	declaration
	Type       Node
	Expression Node
}

//////////////////////////////////////
// NAMESPACE

type AstNamespaceDecl struct {
	nodeBase
	named
	membered
}

///////////////////////////////////////
///////////////////////////////////////

type AstImplement struct {
	nodeBase
	membered
	Name Node // a path
}

type AstTemplateParam struct {
	nodeBase
	named
}

type AstTypeDecl struct {
	nodeBase
	named
	templated
	Members Names
}

type AstEnumDecl struct {
	AstTypeDecl
}

type AstUnionDecl struct {
	AstTypeDecl
	Types []Node
}

type AstStructDecl struct {
	AstTypeDecl
	membered
}

///////////// Functions

type AstFn struct {
	nodeBase
	named
	templated
	IsMethod   bool
	Args       []*AstVarDecl
	ReturnType Node
	Definition Node
}

//////////////
type AstBlock struct {
	nodeBase
	Statements []Node // Is a Vardecl an expression ?
}

type AstFnCall struct {
	nodeBase
	FnExp Node
	Args  []Node
}

type AstIndexCall struct {
	nodeBase
	Indexed Node
	Indices []Node
}

/////////////////////////////

type unaryOperation struct {
	nodeBase
	Operand Node
}

type AstDerefOp struct{ unaryOperation }
type AstPointerOp struct{ unaryOperation }
type AstReturnOp struct{ unaryOperation }
type AstTakeOp struct{ unaryOperation }

//////////////////////////////

type binaryOperation struct {
	nodeBase
	Left  Node
	Right Node
}

func (b *binaryOperation) SetLeft(left Node) {
	b.Left = left
}

func (b *binaryOperation) SetRight(right Node) {
	b.Right = right
}

type binOpNode interface {
	Node
	SetLeft(left Node)
	SetRight(right Node)
}

type AstMulBinOp struct{ binaryOperation }
type AstDivBinOp struct{ binaryOperation }
type AstAddBinOp struct{ binaryOperation }
type AstSubBinOp struct{ binaryOperation }
type AstModBinOp struct{ binaryOperation }

type AstPipeBinOp struct{ binaryOperation }
type AstAmpBinOp struct{ binaryOperation }
type AstLShiftBinOp struct{ binaryOperation }
type AstRShiftBinOp struct{ binaryOperation }

type AstAndBinOp struct{ binaryOperation }
type AstOrBinOp struct{ binaryOperation }
type AstGtBinOp struct{ binaryOperation }
type AstGteBinOp struct{ binaryOperation }
type AstLtBinOp struct{ binaryOperation }
type AstLteBinOp struct{ binaryOperation }
type AstEqBinOp struct{ binaryOperation }
type AstNeqBinOp struct{ binaryOperation }

type AstIsBinOp struct{ binaryOperation }
type AstIsNotBinOp struct{ binaryOperation }

type AstDotBinOp struct{ binaryOperation }

/////////////////////////////////////////////////////////////////////////
//  IDENTIFIER

type Literal struct {
	nodeBase
	noopNud
}

type AstNone struct{ Literal }
type AstTrue struct{ Literal }
type AstFalse struct{ Literal }
type AstIntLiteral struct{ Literal }
type AstStringLiteral struct{ Literal }

type AstIdentifier struct {
	nodeBase
	Name Name
}

func (id *AstIdentifier) GetName() *AstIdentifier {
	return id
}

type AstStringExp struct {
	nodeBase
	Components []Node
}

////////////////////////////////////////////////////////////////////////////
// CONTROL STRUCTURES

type AstIf struct {
	nodeBase
	ConditionExp Node
	ThenArm      Node
	ElseArm      Node
}

type AstWhile struct {
	nodeBase
	ConditionExp Node
	Body         Node
}
