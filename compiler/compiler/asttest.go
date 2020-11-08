package zoe

import (
	"bytes"
	"log"
	"regexp"
	"strings"

	"github.com/fatih/color"
)

var reNoSuperflousSpace = regexp.MustCompilePOSIX(`[ \n\t\r]+`)
var reBeforeSpace = regexp.MustCompilePOSIX(` (\}|\)|\])`)
var reAfterSpace = regexp.MustCompilePOSIX(`(\{|\(|\[) `)
var reAstComments = regexp.MustCompilePOSIX(`--[^\n]*`)

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
func (c *ZoeContext) TestFileAst() {
	for n, cmt := range c.DocCommentMap {
		// found it !
		str := []byte(cmt.String())
		str = []byte(strings.TrimSpace(string(str[3 : len(str)-2]))) // remove the comment marks
		// trim the fat, remove excess
		str = cleanup(str)
		p := color.NoColor
		color.NoColor = true

		var buf bytes.Buffer
		n.Dump(&buf)
		test := buf.String()

		color.NoColor = p

		if string(test) != string(str) {
			log.Print("expected: ", yel(string(str)))
			log.Print("result:   ", red(string(test)))
			log.Print("src: ", n.GetText(), "\n\n")
		} else {
			log.Print("ok: ", green(string(str)), "\n\n")
		}
	}

}
