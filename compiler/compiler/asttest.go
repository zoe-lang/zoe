package zoe

import (
	"log"
	"regexp"
	"strings"

	"github.com/fatih/color"
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
	for node, cmt := range f.DocCommentMap {
		// found it !
		str := []byte(f.GetTokenText(cmt))
		if str[1] == '?' {
			str = []byte(strings.TrimSpace(string(str[2:]))) // remove the comment marks
		} else {
			str = []byte(strings.TrimSpace(string(str[3 : len(str)-2]))) // remove the comment marks
		}
		// trim the fat, remove excess
		str = cleanup(str)
		p := color.NoColor
		color.NoColor = true

		// log.Printf("%#v", node)
		var node_text = node.GetText()
		test := cleanup([]byte(node_text))

		color.NoColor = p

		if string(test) != string(str) {
			log.Print("src: ", grey(node_text), "\n")
			log.Print("expected: ", yel(string(str)))
			log.Print("result:   ", red(string(test)))
		} else {
			log.Print("src: ", grey(node_text), "\n")
			log.Print("ok: ", green(string(str)), "\n")
		}
	}

}
