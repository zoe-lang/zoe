package zoe

type ZoeType interface {
	// Returns the Node containing the type's true definition
	Definition() Node

	// Methods that may be used by comptime
	HasTrait(trait *TypeTrait) bool
}

type TypeAlias struct{}
type TypeTrait struct{}

type TypePointer struct{}
type TypeFloat struct{}
type TypeInteger struct{}

type TypeInstance struct {
}

/*
	Unknown type is the type assigned by default when a variable
	does not specify its type and until the type is known, or when handling
	templates for which the types must be inferred.
*/
type UnknownType struct{}

/*
  Namespaces but also types when used as their name return a TypeNamespace
*/
type TypeNamespace struct {
	Members Names
}

type InvalidType struct{}
