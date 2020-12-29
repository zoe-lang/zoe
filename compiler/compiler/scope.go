package zoe

import (
	"strconv"
)

////////////////////////

type ScopePosition int

func (sp ScopePosition) Handler(file *File) Scope {
	return Scope{
		pos:  sp,
		file: file,
	}
}

/////////////////////////

type concreteScope struct {
	Parent        ScopePosition
	Owner         NodePosition
	EndsExecution bool
	Names         map[InternedString]NodePosition // this is probably not what we want
}

//////////////////////////

type Scope struct {
	pos  ScopePosition
	file *File
}

// get the concrete scope reference
func (s Scope) ref() *concreteScope {
	return &s.file.scopes[s.pos]
}

// Create a new scope in the file
func (f *File) newScope() Scope {
	pos := len(f.scopes)
	f.scopes = append(f.scopes, concreteScope{
		Names:  make(map[InternedString]NodePosition),
		Parent: -1,
	})
	return Scope{pos: ScopePosition(pos), file: f}
}

func (f *File) RootScope() Scope {
	return Scope{
		pos:  0,
		file: f,
	}
}

func (s Scope) getOwner() NodePosition {
	return s.ref().Owner
}

func (s Scope) setOwner(np Node) {
	s.ref().Owner = np.pos
}

func (s Scope) setParent(parent Scope) {
	s.ref().Parent = parent.pos
}

// Create a sub scope, that will go look into its parent
func (s Scope) subScope() Scope {
	newscope := s.file.newScope()
	var parent = s
	newscope.setParent(parent)
	return newscope
}

func (s Scope) Parent() Scope {
	return Scope{
		pos:  s.ref().Parent,
		file: s.file,
	}
}

type ScopeName struct {
	Name string
	Node Node
}

// Get all the names available to a given scope
func (s Scope) AllNames() []ScopeName {
	var res = make([]ScopeName, 0)
	for s.pos != -1 {
		var ref = s.ref()
		for i, n := range ref.Names {
			res = append(res, ScopeName{Name: GetInternedString(i), Node: n.Node(s.file)})
		}
		s = ref.Parent.Handler(s.file)
	}
	return res
}

// Find a name in the scope
func (s Scope) Find(name InternedString) (Node, bool) {
	sc := s.ref()
	if node, ok := sc.Names[name]; ok {
		return Node{pos: node, file: s.file}, true
	}

	return EmptyNode, false
}

func (s Scope) FindRecursive(name InternedString) (Node, bool) {
	if node, ok := s.Find(name); ok {
		return node, true
	}

	sc := s.ref()
	if sc.Parent != -1 {
		return sc.Parent.Handler(s.file).Find(name)
	}

	return EmptyNode, false
}

func (s Scope) addSymbolFromIdNodeForce(idnode Node, target Node) {
	// idn := sh.file.Nodes[idnode]
	if !idnode.Is(NODE_ID) {
		s.file.reportError(idnode.Range(), "is not an identifier")
		return
	}

	var name = idnode.InternedString()

	sc := &s.file.scopes[s.pos]

	sc.Names[name] = target.pos
}

func (s Scope) addSymbolFromIdNode(idnode Node, target Node) {
	// idn := sh.file.Nodes[idnode]
	if !idnode.Is(NODE_ID) {
		s.file.reportError(idnode.Range(), "is not an identifier")
		return
	}

	var name = idnode.InternedString()

	sc := &s.file.scopes[s.pos]

	if orig, ok := s.FindRecursive(name); ok {
		// we do not set that variable since it already existed in one of our parent scope.
		// note ; the choice was made to not allow shadowing to avoid footguns, since every
		// Zoe module needs to explicitely import other symbols (except maybe for core, which will then pollute)
		s.file.reportError(idnode.Range(), "identifier '", GetInternedString(name), "' was already defined at line ", strconv.Itoa(int(orig.Range().Start.Line+1)))
		return
	}

	sc.Names[name] = target.pos
}
