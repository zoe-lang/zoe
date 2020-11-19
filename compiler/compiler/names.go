package zoe

// Name resolution
// This is how a given id (only the left most of a . sequence) or a :: operation gets resolved.
// The type checker is the one that resolves the . operations, as they always need the types

// Do all the files share a common pool of names ? Does it need to be "garbage collected" (probably) ?

// How are the mapping done between definition and usage, since we need to track those for go to definition and usage ?
// How are the types associated to symbols ?

// the type checker is the one in charge of tracking the usage of a given definition.
