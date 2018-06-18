package zoe

// For every node at the root of the file, we're going to try
// and find a matching doc comment containing an AST dump.
// We'll then parse it, dump it, and compare the dump to the dump of the corresponding AST.
func (c *ZoeContext) TestFileAst() {
	// How do I get the doc comment that is tied to a node ?
	// Should I just add a DocComment field right next to a node ?
	// Should there be a map of doc comment to Node ? <- this is probably a better idea, less messy.
	// But wait, how does one find a doc comment to a particular symbol ?
}
