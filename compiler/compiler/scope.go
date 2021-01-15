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
