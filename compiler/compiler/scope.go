package zoe

import "strconv"

type ScopePosition int
type ScopeType int

const (
	SCOPE_NAMESPACE ScopeType = iota + 1
	SCOPE_FUNCTION
	SCOPE_IMPLEMENT
	SCOPE_BLOCK
)

type Scope struct {
	Parent ScopePosition
	Names  map[int]NodePosition
	Kind   ScopeType
}

func (b *nodeBuilder) ScopeNew(kind ScopeType) ScopePosition {
	pos := len(b.scopes)
	b.scopes = append(b.scopes, Scope{
		Names:  make(map[int]NodePosition),
		Kind:   kind,
		Parent: -1,
	})

	return ScopePosition(pos)
}

func (b *nodeBuilder) ScopeNewSub(parent ScopePosition, kind ScopeType) ScopePosition {
	pos := b.ScopeNew(kind)
	b.scopes[pos].Parent = parent
	return pos
}

// Add a symbol to a scope
func (b *nodeBuilder) ScopeAddSymbol(scope ScopePosition, pos NodePosition) {
	s := &b.scopes[scope]
	node := b.nodes[pos]
	if node.Kind != NODE_ID {
		b.reportErrorAtPosition(pos, "COMPILER ERROR not an id but was added to scope")
		return
	}

	value := node.Value

	if orig, ok := b.ScopeFindPosition(scope, value); ok {
		// we do not set that variable since it already existed in one of our parent scope.
		// note ; the choice was made to not allow shadowing to avoid footguns, since every
		// Zoe module needs to explicitely import other symbols (except maybe for core, which will then pollute)
		b.reportErrorAtPosition(pos, "symbol '", InternedIds.Get(value), "' was already defined at line ", strconv.Itoa(int(b.nodes[orig].Range.Line)))
		return
	}
	s.Names[value] = pos
}

// FindPosition
func (b *nodeBuilder) ScopeFindPosition(pos ScopePosition, name int) (NodePosition, bool) {
	s := b.scopes[pos]
	if pos, ok := s.Names[name]; ok {
		return pos, true
	}
	if s.Parent != -1 {
		return b.ScopeFindPosition(s.Parent, name)
	}
	return EmptyNode, false
}
