package zoe

// We have symbols in scopes
//
// A symbol can be
//	 - an import
//   - a variable
//   - a function / method
//   - a named scope (like a file or a namespace)
//   - a type definition (struct / type / union / enum), which is itself a namespace when resolved directly
//   - a template !
//
// A symbol always has members !
// What about scopes ? Should scopes live separately ?
// In a file, the root scope is available to all. But outside the file, the root scope
// becomes a module, which has members !
//
// a file module, when a name registers ;
//   - registers the name in its associated scope
//   - *also* registers the name in its members
//
// All symbols have types associated with them, that can be retrived using GetType() on them.
// A Type's type is a TypeDefinition.
// A type may have members, which are struct members and methods.
//
// A type is associated to a namespace, where *methods* may be located that can
// be called on a type (with instance as first argument) or be called an instance.
//
// We distinguish two types of symbol scopes ;
//   - the ones where instructions are being executed, in the body of functions
//   - the ones where only definitions reside (and some execution ?)

type Names map[Name]Symbol

func (n *Names) RegisterSymbol(sym Symbol) {
	(*n)[sym.Name()] = sym
}

// Nameholders can be scopes or types
// Scopes just add the symbol to their map, while types can filter stuff and
// send to their namespace
type Nameholder interface {
	RegisterSymbol(name Name, sym Symbol)
}

///////////////////////////////////////////////////
// Symbol

type Symbol interface {
	Name() Name
	Scope() Scope
	Node() Node
}

type symbolBase struct {
	name Name
	node Node
}

func (s *symbolBase) Name() Name {
	return s.name
}

func (s *symbolBase) Scope() Scope {
	return s.node.Scope()
}

func (s *symbolBase) Node() Node {
	return s.node
}

////////////////////////////////////////////////////////
// Types with members

type Membered struct {
	names Names
}

func (m *Membered) GetMembers() Names {
	return m.names
}

// Namespace holds names
type Namespace struct {
	names Names
}

////////////////////////////////////////////////////////////////
// All the type definitions need the following

type typedefBase struct {
	symbolBase
	Namespace
	Implements []*Implement
}

func (b *typedefBase) RegisterImplement(impl *Implement) {
	b.Implements = append(b.Implements, impl)
	impl.ParentType = b
}

type TypeDefinition interface {
	RegisterImplement(impl *Implement)
	// TryRegisterSymbol(sym Symbol) bool
}

///////////////////////////////////////////////////////////////////

type NamespaceDef struct {
	typedefBase
	Membered
}

// StructDef gets its members
type StructDef struct {
	typedefBase
	Membered
}

// When struct receives a name, it must check if it should add it to its members,
// to its namespace, or to its associated scope (its namespace is its scope)

type Implement struct {
	Membered
	ParentType TypeDefinition // implements always have an associated type
	Node       Node
}

type EnumDef struct {
	typedefBase
	Membered // there needs to be a check done somewhere ?
}

type UnionDef struct {
	// shit shit, type expressions...
	TypeDefinition
	Union Union
}

////// Function definition

type FnDef struct {
	symbolBase
}

////////

type Union struct {
	TypeDefinition
	Membered
	Types []TypeExpression
}

// Invalid type is the result
type InvalidType struct {
	TypeDefinition
}

//////////////////////////////////////////////////////
// Pure symbols

type Import struct {
	symbolBase
}

type Variable struct {
	symbolBase
}

//////////////////////////////////////////////////////
//

// TypeExpression
type TypeExpression interface {
	GetMembers() Names
}
