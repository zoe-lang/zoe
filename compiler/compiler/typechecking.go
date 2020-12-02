package zoe

// Name resolution
//
// Filenames scope (import statements look there when their module is a string)
//
// Identifier resolution ;
//   Builtin scope (names and identifiers defined by the compiler, such as Int, Int32, $, ...)
//   	 Core scope (this is the real top level scope, has String, Map, Set, Slice, because they're coded in zoe)
//       -> current file toplevel
//         -> ... subscopes.
//
// Since strings are interned, maybe use these ids ?
// ... probably...

// File names are "identifiers" that are requested similarly to identifiers
//
// Symbolic links *must* be converted to their real file names

// Declared symbols remember where they are used
// Idents remember where their declared symbols are
// idents are hashed (?)

// Are identifiers globally stored ? We have a string interner, and it is global, so yeah.

// Goto definition is simply a matter of going up the scope chain and going back to the definition. It is expected to be
// very fast.
//
// Find implementations is however a little more involved, because we have to go look into all the open files if there
// is the possibility of the given symbol to be used.

// Context records the links between symbol uses and their corresponding declaration in the same file
// or in other files
type Context struct {
}
