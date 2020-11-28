package zoe

import (
	"bytes"
	"io"
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

func (f *File) PrintNode(w io.Writer, iter NodePosition) {
	n := &f.Nodes[iter]
	if n.ArgLen > 0 {
		_, _ = w.Write([]byte{'('})
	}
	f.PrintNodeRepr(w, iter)
	for i := 0; i < n.ArgLen; i++ {
		f.PrintNodeArg(w, n.Args[i])
	}
	if n.ArgLen > 0 {
		_, _ = w.Write([]byte{')'})
	}
}

func (f *File) PrintNodeArg(w io.Writer, iter NodePosition) {
	if iter == EmptyNode {
		_, _ = w.Write([]byte{' ', '~'})
		return
	}
	_, _ = w.Write([]byte{' '})
	if f.Nodes[iter].Next != EmptyNode {
		f.PrintNodeList(w, iter)
	} else {
		f.PrintNode(w, iter)
	}
}

func (f *File) PrintNodeList(w io.Writer, iter NodePosition) {
	_, _ = w.Write([]byte("["))
	first := true
	for iter != 0 {
		if !first {
			_, _ = w.Write([]byte(" "))
		} else {
			first = false
		}
		f.PrintNode(w, iter)
		iter = f.Nodes[iter].Next
	}
	_, _ = w.Write([]byte("]"))
}

func (f *File) PrintNodeString(n NodePosition) string {
	var buf bytes.Buffer
	f.PrintNode(&buf, n)
	return buf.String()
}

// For every node at the root of the file, we're going to try
// and find a matching doc comment containing an AST dump.
// We'll then parse it, dump it, and compare the dump to the dump of the corresponding AST.
func (f *File) TestFileAst() {
	for n, cmt := range f.DocCommentMap {
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

		test := cleanup([]byte(f.PrintNodeString(n)))

		color.NoColor = p

		if string(test) != string(str) {
			log.Print("src: ", grey(f.GetNodeText(n)), "\n")
			log.Print("expected: ", yel(string(str)))
			log.Print("result:   ", red(string(test)))
		} else {
			log.Print("src: ", grey(f.GetNodeText(n)), "\n")
			log.Print("ok: ", green(string(str)), "\n")
		}
	}

}
