package zoe

import (
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/pkg/errors"
	"github.com/sourcegraph/go-lsp"
)

type ZoeError struct {
	File    *File
	Range   Range
	Message string
}

func (err ZoeError) Print(w io.Writer) {
	fred.Fprint(w, err.File.Filename+" ")
	fgreen.Fprint(w, err.Range.Line)
	_, _ = w.Write([]byte(": " + err.Message + "\n"))
}

func (err ZoeError) ToLspDiagnostic() lsp.Diagnostic {
	d := lsp.Diagnostic{}
	d.Message = err.Message
	d.Range.Start.Line = int(err.Range.Line - 1)
	d.Range.Start.Character = int(err.Range.Column - 1)
	d.Range.End.Line = int(err.Range.LineEnd - 1)
	d.Range.End.Character = int(err.Range.ColumnEnd - 1)
	return d
}

// File holds the current parsing context
// also does the error handling stuff.
type File struct {
	Filename string

	Tokens []Token
	scopes []concreteScope
	Nodes  []AstNode

	RootNode Node
	Version  int

	Errors []ZoeError
	data   []byte

	DocCommentMap map[NodePosition]TokenPos // node position => token position
}

func NewFileFromContents(filename string, contents []byte) (*File, error) {
	data := []byte(contents)

	ctx := &File{
		Filename:      filename,
		Errors:        make([]ZoeError, 0),
		data:          append(data, '\x00'),
		DocCommentMap: make(map[NodePosition]TokenPos),
		scopes:        make([]concreteScope, 0),
		// RootDocComments: make([]*Token, 0),
	}
	// create the root scope.
	ctx.newScope()

	lxerr := ctx.Lex()
	if lxerr != nil {
		return ctx, errors.Wrap(lxerr, "lexing failed")
	}

	return ctx, nil

}

// NewFile
func NewFile(filename string) (*File, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	data, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	return NewFileFromContents(filename, data)
}

func (f *File) GetTokenText(tk TokenPos) string {
	return f.GetRangeText(f.Tokens[tk].Range)
}

func (f *File) GetRangeText(p Range) string {
	return string(f.data[p.Start:p.End])
}

func (f *File) GetNodeText(n NodePosition) string {
	return f.GetRangeText(f.Nodes[n].Range)
}

func (f *File) reportError(pos Positioned, message ...string) {

	f.Errors = append(f.Errors, ZoeError{
		File:    f,
		Range:   *pos.GetPosition(),
		Message: strings.Join(message, ""),
	})
	// f.Errors[len(f.Errors)-1].Print(os.Stderr)
}

func (f *File) createNode(rng Range, kind AstNodeKind, scope Scope, children ...Node) Node {
	// maybe we should handle here the capacity of the node arrays ?
	l := NodePosition(len(f.Nodes))
	f.Nodes = append(f.Nodes, AstNode{Kind: kind, Range: rng, Scope: scope.pos})

	cl := len(children)
	if cl > 0 {
		node := &f.Nodes[l]
		node.ArgLen = int8(cl)
		for i, chld := range children {
			node.Args[i] = chld.pos
			if !chld.IsEmpty() {
				node.Range.Extend(chld.Range())
			}
		}
	}

	return Node{
		file: f,
		pos:  l,
	}
}

// Find a node that matches a given range
func (f *File) FindNodePosition(lsppos *lsp.Position) (Node, error) {
	node := f.RootNode
	// nodes := f.Nodes

	log.Print(lsppos.Line+1, ":", lsppos.Character+1)
search:
	for node.HasPosition(lsppos) {
		log.Print(node.Debug())
		// First check in the node's children
		for _, chl := range node.ref().Args {
			chld := chl.Node(f)
			if chld.IsEmpty() {
				continue
			}
			if chld.HasPosition(lsppos) {
				// log.Print(f.NodeDebug(chld))
				node = chld
				continue search
			}

			// Then check in its list
			other := chld.Next()
			for !other.IsEmpty() {
				// log.Print(f.NodeDebug(other))
				if other.HasPosition(lsppos) {
					node = other
					continue search
				}
				other = other.Next()
			}

		}

		break
	}

	return node, nil
}
