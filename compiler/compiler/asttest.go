package zoe

import (
	"bytes"
	"fmt"
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

func (n Node) Debug() string {
	rng := n.ref().Range
	return fmt.Sprintf("(%v:%s %v:%v - %v:%v [%v] -> %v)", n.pos, n.Repr(), rng.Line, rng.Column, rng.LineEnd, rng.ColumnEnd, n.ref().Args, n.Next)
}

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

func (f *File) PrintNode(w io.Writer, iter Node) {
	// n := &f.Nodes[iter]
	if iter.Is(NODE_BLOCK) {
		f.PrintNodeList(w, iter.GetArg(0), [2]byte{'{', '}'})
		return
	}

	al := iter.ArgLen()
	if al > 0 {
		_, _ = w.Write([]byte{'('})
	}
	_, _ = w.Write([]byte(iter.Repr()))
	for i := 0; i < al; i++ {
		f.PrintNodeArg(w, iter.GetArg(i))
	}
	if al > 0 {
		_, _ = w.Write([]byte{')'})
	}
}

func (f *File) PrintNodeArg(w io.Writer, iter Node) {
	if iter.IsEmpty() {
		_, _ = w.Write([]byte{' ', '~'})
		return
	}
	_, _ = w.Write([]byte{' '})
	if !iter.Next().IsEmpty() {
		f.PrintNodeList(w, iter, [2]byte{'[', ']'})
	} else {
		f.PrintNode(w, iter)
	}
}

func (f *File) PrintNodeList(w io.Writer, iter Node, pair [2]byte) {
	_, _ = w.Write([]byte{pair[0]})
	first := true
	for !iter.IsEmpty() {
		if !first {
			_, _ = w.Write([]byte(" "))
		} else {
			first = false
		}
		f.PrintNode(w, iter)
		iter = iter.Next()
	}
	_, _ = w.Write([]byte{pair[1]})
}

func (f *File) PrintNodeString(n Node) string {
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

		test := cleanup([]byte(f.PrintNodeString(n.Node(f))))

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
