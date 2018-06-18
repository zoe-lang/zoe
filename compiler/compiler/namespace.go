package zoe

// The namespace is the core of the program

type TypeInstance struct {
	Declaration *TypeDeclaration
	Args        []*TypeInstance
}

type FunctionInstance struct {
	Declaration *FunctionDeclaration
	Args        []*TypeInstance
}

type TypeDeclaration struct {
	IsLocal     bool
	Id          string // a unique identifier accross all programs
	GenericArgs []*VariableDeclaration
	Instance    *TypeInstance // ??
}

type VariableDeclaration struct {
	IsLocal   bool
	TypeExpr  *Node // can be nil
	ValueExpr *Node // can be nil

	ResolvedType *TypeInstance // nil by default
}

type FunctionDeclaration struct {
	IsLocal        bool
	GenericArgs    []*VariableDeclaration
	Args           []*VariableDeclaration
	ReturnTypeExpr *Node // can be nil
	Body           *Node // can be nil, like for traits ?
}

type Namespace struct {
	Functions map[string]*FunctionDeclaration
	Variables map[string]*VariableDeclaration
	Types     map[string]*TypeDeclaration
}
