package zoe

import (
	"io"
)

func setNode(parent Node, target *Node, assign Node) {
	parent.ExtendPosition(assign)
	*target = assign
}

func appendNode(parent Node, target *[]Node, assign Node) {
	parent.ExtendPosition(assign)
	*target = append(*target, assign)
}

type Node interface {
	DumpString() string
	Dump(w io.Writer)
	SetError()
	GetPosition() *Position
	GetText() string
	ExtendPosition(other Positioned)
	ReportError(msg ...string)
	EnsureTuple() *Tuple
}

type NodeBase struct {
	Position
	IsError bool
}

//////////////////// LIST NODES //////////////////////

/////////////////////////////////////////////////////////////////////////

/////////////////////////////////////////////////////////////////////////
/////////////////////////////////////////////////////////////////////////

// NODEBASE

func (n *NodeBase) SetError() {
	n.IsError = true
}

func (n *NodeBase) ReportError(msg ...string) {
	n.IsError = true
	n.Context.reportError(n.Position, msg...)
}

func (n *NodeBase) GetPosition() *Position {
	return &n.Position
}

// ExtendPosition extends the position of the node
func (n *NodeBase) ExtendPosition(otherp Positioned) {
	other := otherp.GetPosition()
	if other.Line == 0 {
		// do not try to absorb a faulty position
		return
	}

	pos := &n.Position

	if pos.Line == 0 {
		// In the case that this position is faulty, give it the other one
		*pos = *other
		return
	}

	pos.Start = minInt(pos.Start, other.Start)
	pos.End = maxInt(pos.End, other.End)

	if other.Line < pos.Line {
		pos.Line = other.Line
		pos.Column = other.Column
	} else if other.Line == pos.Line {
		pos.Column = minInt(pos.Column, other.Column)
	}

	if other.EndLine > pos.EndLine {
		pos.EndLine = other.EndLine
		pos.EndColumn = other.EndColumn
	} else if other.EndLine == pos.EndLine {
		pos.EndColumn = maxInt(pos.EndColumn, other.EndColumn)
	}
}

/////////////////////////////////////////////////////////////////////////////

type Namespace struct {
	NodeBase
	Identifier Node
	Block      *Block
}

type Fragment struct {
	NodeBase
	Children []Node
}

type TypeDecl struct {
	NodeBase
	Ident    *BaseIdent
	Template *Template
	Def      Node
}

type Union struct {
	NodeBase
	TypeExprs []Node
}

type ImportAs struct {
	NodeBase
	Path Node
	As   *BaseIdent
}

type ImportList struct {
	NodeBase
	Path  Node
	Names *Tuple
}

type Var struct {
	NodeBase
	Ident   *BaseIdent
	TypeExp Node
	Exp     Node // Exp can be potentially empty
}

type Operation struct {
	NodeBase
	Token    *Token
	Operands []Node
}

func (o *Operation) Is(tk TokenKind) bool {
	return o.Token.Kind == tk
}

func (o *Operation) IsBinary() bool {
	return len(o.Operands) == 2
}

func (o *Operation) IsUnary() bool {
	return len(o.Operands) == 1
}

func (o *Operation) Left() Node {
	return o.Operands[0]
}

func (o *Operation) Right() Node {
	return o.Operands[1]
}

///////////////////////////////////////////////////////////////
// TEMPLATE

type Template struct {
	NodeBase
	Args *VarTuple
	// maybe probably here would go a where clause
}

///////////////////////////////////////////////////////////////
// FNDEF

type FnDef struct {
	NodeBase
	Template   *Template
	Signature  *Signature
	Definition *Block
}

type Signature struct {
	NodeBase
	Args          *VarTuple
	ReturnTypeExp Node
}

type FnCall struct {
	NodeBase
	Left Node
	Args *Tuple
}

type GetIndex struct {
	NodeBase
	Left  Node
	Index Node
}

type SetIndex struct {
	NodeBase
	Left  Node
	Index Node
	Value Node
}

type If struct {
	NodeBase
	Cond Node
	Then Node
	Else Node
}

///////////////////////////////////////////////////////////////
// FNDECL

type FnDecl struct {
	NodeBase
	Ident *BaseIdent
	FnDef *FnDef
}

///////////////////////////////////////////////////////////////
// BLOCK

type Tuple struct {
	NodeBase
	Children []Node
}

type VarTuple struct {
	NodeBase
	Vars []*Var
}

// Transforms a tuple to arguments
func (t *Tuple) ToVars() *VarTuple {
	tup := &VarTuple{}
	for _, c := range t.Children {
		v := coerceToVar(c)
		if v != nil {
			tup.AddVars(v)
		}
	}
	return tup
}

///////////////////////////////////////////////////////////////
// BLOCK

type Block struct {
	NodeBase
	Children []Node
}

///////////////////////////////////////////////////////////////
// RETURN

type Return struct {
	NodeBase
	Expr Node
}

///////////////////////////////////////////////////////////////
// ID

type Ident struct {
	NodeBase
	Path []*BaseIdent
}

type BaseIdent struct {
	NodeBase
}

///////////////////////////////////////////////////////////////
// STRING

////////////////////////////////////////////////////////////////////////
//

/////////////////////////////////////////////////

type Null struct {
	NodeBase
}

type False struct {
	NodeBase
}

type True struct {
	NodeBase
}

type String struct {
	NodeBase
}

type Integer struct {
	NodeBase
}

type Float struct {
	NodeBase
}

type Eof struct {
	NodeBase
}

func coerceToVar(n Node) *Var {
	switch v := n.(type) {
	case *Operation:
		if v.Is(TK_EQ) && len(v.Operands) == 2 {
			// this is perfect, check that left is
			if id, ok := v.Operands[0].(*BaseIdent); ok {
				return n.GetPosition().CreateVar().SetExp(v.Operands[1]).SetIdent(id)
			}
		}
	case *BaseIdent:
		res := &Var{}
		return res.SetIdent(v)
	case *Var:
		return v
	}
	n.ReportError(`variable declaration expected`)
	return nil
}
