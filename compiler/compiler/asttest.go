package zoe

import (
	"regexp"
)

var reNoSuperflousSpace = regexp.MustCompilePOSIX(`[ \n\t\r]+`)
var reBeforeSpace = regexp.MustCompilePOSIX(` (\}|\)|\])`)
var reAfterSpace = regexp.MustCompilePOSIX(`(\{|\(|\[) `)
var reAstComments = regexp.MustCompilePOSIX(`--[^\n]*`)

// func (n Node) Debug() string {
// 	// rng := n.ref().Range
// 	color.NoColor = true
// 	// var res = fmt.Sprintf("(%v)", n.DebugName(), n.Repr(), n.ref().Range.Start, n.ref().Range.End)
// 	color.NoColor = false
// 	return "node?"
// }

func cleanup(str []byte) []byte {
	s := reAstComments.ReplaceAllLiteral(str, []byte{})
	s = reNoSuperflousSpace.ReplaceAllFunc(str, func(b []byte) []byte {
		return []byte{' '}
	})
	s = reBeforeSpace.ReplaceAllFunc(s, func(b []byte) []byte {
		return []byte{b[1]}
	})
	s = reAfterSpace.ReplaceAllFunc(s, func(b []byte) []byte {
		return []byte{b[0]}
	})
	return s
}

// For every node at the root of the file, we're going to try
// and find a matching doc comment containing an AST dump.
// We'll then parse it, dump it, and compare the dump to the dump of the corresponding AST.
func (f *File) TestFileAst() {
	// for node, cmt := range f.DocCommentMap {
	// 	// found it !
	// }

}
