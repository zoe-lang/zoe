package zoe

import "io"

func setNode(parent Node, target *Node, assign Node) {
	parent.ExtendPosition(assign)
	*target = assign
}

func appendNode(parent Node, target *[]Node, assign Node) {
	parent.ExtendPosition(assign)
	*target = append(*target, assign)
}

type Node interface {
	Dump(w io.Writer)
	SetError()
	GetPosition() *Position
	GetText() string
	ExtendPosition(other Positioned)
	ReportError(msg ...string)
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
	Children []Node
}

type Fragment struct {
	NodeBase
	Children []Node
}

type TypeDecl struct {
	NodeBase
	Template *Template
	Ident    *Ident
	Def      Node
}

type Union struct {
	NodeBase
	TypeExprs []Node
}

type Var struct {
	NodeBase
	Ident   *Ident
	TypeExp Node
	Exp     Node // Exp can be potentially empty
}

type Operation struct {
	NodeBase
	TokenKind TokenKind
	Operands  []Node
}

///////////////////////////////////////////////////////////////
// TEMPLATE

type Template struct {
	NodeBase
	Args []*Var
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
	Args          []*Var
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
	Ident *Ident
	FnDef *FnDef
}

///////////////////////////////////////////////////////////////
// BLOCK

type Tuple struct {
	NodeBase
	Children []Node
}

// Transforms a tuple to arguments
func (t *Tuple) ToVars() []*Var {
	res := make([]*Var, 0)
	for _, c := range t.Children {
		v := coerceToVar(c)
		if v != nil {
			res = append(res, v)
		}
	}
	return res
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
		if v.TokenKind == TK_EQ && len(v.Operands) == 2 {
			// this is perfect, check that left is
			if id, ok := v.Operands[0].(*Ident); ok {
				return n.GetPosition().CreateVar().SetExp(v.Operands[1]).SetIdent(id)
			}
		}
	case *Ident:
		res := &Var{}
		return res.SetIdent(v)
	case *Var:
		return v
	}
	n.ReportError(`variable declaration expected`)
	return nil
}
