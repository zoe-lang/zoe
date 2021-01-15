package zoe

type ZoeType interface {
	Definition() Node

	// Methods that may be used by comptime
	HasTrait(trait *TypeTrait) bool
}

type TypeTrait struct{}

type TypePointer struct{}
type TypeFloat struct{}
type TypeInteger struct{}

type TypeInstance struct {
}

/*
  Namespaces but also types when used as their name return a TypeNamespace
*/
type TypeNamespace struct {
	Members Names
}

type InvalidType struct{}
