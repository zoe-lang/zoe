package zoe

import (
	"strconv"
)

type ScopePosition int

func (sp ScopePosition) Handler(file *File) Scope {
	return Scope{
		pos:  sp,
		file: file,
	}
}

type concreteScope struct {
	Parent ScopePosition
	Owner  NodePosition
	Names  map[InternedString]NodePosition
}

func (f *File) newScope() ScopePosition {
	pos := len(f.scopes)
	f.scopes = append(f.scopes, concreteScope{
		Names:  make(map[InternedString]NodePosition),
		Parent: -1,
	})
	return ScopePosition(pos)
}

func (f *File) RootScope() Scope {
	return Scope{
		pos:  0,
		file: f,
	}
}

type Scope struct {
	pos  ScopePosition
	file *File
}

func (sh Scope) setOwner(node Node) {
	sh.file.scopes[sh.pos].Owner = node.pos
}

func (sh Scope) getOwner() NodePosition {
	return sh.file.scopes[sh.pos].Owner
}

func (sh Scope) setParent(parent ScopePosition) {
	sh.file.scopes[sh.pos].Parent = parent
}

func (sh Scope) subScope() Scope {
	newscope := sh.file.newScope()

	h := Scope{
		pos:  newscope,
		file: sh.file,
	}

	h.setParent(sh.pos)

	return h
}

// Find a name in the scope
func (sh Scope) Find(name InternedString) (Node, bool) {
	sc := sh.file.scopes[sh.pos]
	if node, ok := sc.Names[name]; ok {
		return Node{pos: node, file: sh.file}, true
	}

	return EmptyNode, false
}

func (sh Scope) FindRecursive(name InternedString) (Node, bool) {
	if node, ok := sh.Find(name); ok {
		return node, true
	}

	sc := sh.file.scopes[sh.pos]
	if sc.Parent != -1 {
		return sc.Parent.Handler(sh.file).Find(name)
	}

	return EmptyNode, false
}

func (sh Scope) addSymbolFromIdNode(idnode Node, pos Node) {
	// idn := sh.file.Nodes[idnode]
	if !idnode.Is(NODE_ID) {
		sh.file.reportError(idnode.Range(), "is not an identifier")
		return
	}
	sh.addSymbol(idnode.InternedString(), pos.pos)
}

// Add a symbol to a scope
func (sh Scope) addSymbol(name InternedString, pos NodePosition) {
	s := &sh.file.scopes[sh.pos]
	// node := b.nodes[pos]
	// if node.Kind != NODE_ID {
	// 	b.reportErrorAtPosition(pos, "COMPILER ERROR not an id but was added to scope")
	// 	return
	// }

	// value := InternedString(node.Value)

	if orig, ok := sh.Find(name); ok {
		// we do not set that variable since it already existed in one of our parent scope.
		// note ; the choice was made to not allow shadowing to avoid footguns, since every
		// Zoe module needs to explicitely import other symbols (except maybe for core, which will then pollute)
		sh.file.reportError(sh.file.Nodes[pos].Range, "identifier '", GetInternedString(name), "' was already defined at line ", strconv.Itoa(int(orig.Range().Line)))
		return
	}

	s.Names[name] = pos
}
