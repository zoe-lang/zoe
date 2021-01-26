package zoe

type ScopeContext uint8

const (
	scopeInstructions ScopeContext = iota
	scopeArguments
	scopeType
	scopeFile
)

type Scope struct {
	Context ScopeContext
	Parent  *Scope
	Names   Names
}

func (s *Scope) isFunctionArguments() bool {
	return s.Context == scopeArguments
}

func (s *Scope) isInstructions() bool {
	return s.Context == scopeInstructions
}

func (s *Scope) isType() bool {
	return s.Context == scopeType
}

func (s *Scope) subScope(kind ScopeContext) *Scope {
	return &Scope{
		Parent:  s,
		Names:   make(Names),
		Context: kind,
	}
}

func (s *Scope) SetParent(parent *Scope) {
	s.Parent = parent
}

func (s *Scope) Add(node Node) {
	var name = node.GetName()
	if name == nil {

	}
}

func (s *Scope) ResolveSymbolForName(name Name) Node {
	return nil
}

/*
	Resolve an expression and get the node that is described by a symbol
*/
func (s *Scope) ResolveSymbolForExpression(exp Node) Node {
	switch v := exp.(type) {
	case *AstIdentifier:
		// If it is an identifier, get its name
		return v.Scope.ResolveSymbolForName(v.Name)
	case *AstDotBinOp:
		var left = s.ResolveSymbolForExpression(v.Left)
		if left == nil {
			return nil
		}
		// Once we have left, we treat it as a namespace
	}
	return nil
}
