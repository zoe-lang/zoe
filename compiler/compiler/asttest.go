package zoe

import (
	"log"
	"regexp"
	"strings"

	"github.com/fatih/color"
)

var reNoSuperflousSpace = regexp.MustCompilePOSIX(`[ \n\t\r]+`)
var reBeforeSpace = regexp.MustCompilePOSIX(` \)`)
var reAstComments = regexp.MustCompilePOSIX(`--[^\n]*`)

func cleanup(str []byte) []byte {
	s := reAstComments.ReplaceAllLiteral(str, []byte{})
	s = reNoSuperflousSpace.ReplaceAllFunc(str, func(b []byte) []byte {
		return []byte{' '}
	})
	s = reBeforeSpace.ReplaceAllFunc(s, func(b []byte) []byte {
		return []byte{')'}
	})
	return s
}

func (c *ZoeContext) TestFileAst() {
	c.testFileAst(c.Root)
}

// For every node at the root of the file, we're going to try
// and find a matching doc comment containing an AST dump.
// We'll then parse it, dump it, and compare the dump to the dump of the corresponding AST.
func (c *ZoeContext) testFileAst(n *Node) {
	if cmt, ok := c.DocCommentMap[n]; ok {
		// found it !
		str := []byte(cmt.String())
		str = []byte(strings.TrimSpace(string(str[3 : len(str)-2]))) // remove the comment marks
		// trim the fat, remove excess
		str = cleanup(str)
		p := color.NoColor
		color.NoColor = true
		test := cleanup([]byte(n.Debug()))
		color.NoColor = p

		if string(test) != string(str) {
			log.Print("expected: ", yel(string(str)))
			log.Print("result:   ", red(string(test)))
			log.Print("src: ", n.String())
		} else {
			log.Print("ok: ", green(string(str)))
		}
	}

	for _, ch := range n.Children {
		c.testFileAst(ch)
	}
}
