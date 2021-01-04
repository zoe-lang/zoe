package zoe

// A type has two components ; its definition and its namespace / methods

type NodePath []Node

func (n Node) ResolveMember() Node {
	return Node{}
}

// For a given expression, find the symbol it is referencing. This could
// be a variable, a type, or a member.
// This is generally a prelude to type finding.
// It crosses files boundaries when on an import.
func (n Node) FindDefinition() (Node, bool) {
	if !n.Is(NODE_ID) {
		panic("compiler error, the node should be an Id node")
	}
	var is = n.InternedString()
	// log.Print(is, " -> ", GetInternedString(is))
	if found, ok := n.Scope().FindRecursive(is); ok {
		// This is where we should check whether found is actually an import.
		return found, true
	}
	return Node{}, false
}

// Find the type of a variable
func (np NodePath) FindTypeDefinition() Node {
	return Node{}
}

////////////////////////////////////////////////////////////////////////////////
///
///
///

// Type interface represents something that holds a Type
type IType interface {
	IsDefined() bool
}

type TypeBase struct{}

func (tb *TypeBase) IsDefined() bool {
	return tb != nil
}

type TypeModule struct {
	TypeBase
	Members map[InternedString]IType
}

type TypeTemplate struct {
	TypeBase
	// Is it really IType, since we're going to have to resolve ?
	Args []IType
}

type TypeFunction struct {
	TypeBase
	Args       []IType
	ReturnType IType
}

type TypeEnum struct {
	// Also, values ?
	TypeBase
	Members map[InternedString]IType
}

type TypeStruct struct {
	TypeBase
	Fields map[InternedString]IType
}

type TypePointer struct {
	TypeBase
	Pointed IType
}

type TypeUnion struct {
	TypeBase
	Members map[InternedString]IType
}

type TypeLiteral struct {
	TypeBase
}

var (
	None = &TypeLiteral{}
	Void = &TypeLiteral{}
)
