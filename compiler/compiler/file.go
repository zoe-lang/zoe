package zoe

import (
	"io"
	"io/ioutil"
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
	Filename    string
	Tokens      []Token
	Nodes       NodeArray
	RootNodePos NodePosition
	scopes      []concreteScope
	Version     int

	Errors []ZoeError
	data   []byte

	current *Token
	tkpos   uint32

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

func (f *File) isEof() bool {
	return f.current == nil
}

func (f *File) reportError(pos Positioned, message ...string) {
	f.Errors = append(f.Errors, ZoeError{
		File:    f,
		Range:   *pos.GetPosition(),
		Message: strings.Join(message, ""),
	})
}

func (f *File) createNodeBuilder() *nodeBuilder {

	b := nodeBuilder{
		file:          f,
		tokens:        f.Tokens,
		tokensLen:     TokenPos(len(f.Tokens)),
		doccommentMap: f.DocCommentMap,
	}

	b.createNode(Range{}, NODE_EMPTY, f.RootScope())

	return &b
}

// Find a node that matches a given range
func (f *File) FindNodePosition(lsppos *lsp.Position) (NodePosition, error) {
	pos := f.RootNodePos
	nodes := f.Nodes

	// log.Print(lsppos.Line+1, ":", lsppos.Character+1)
search:
	for nodes[pos].Range.HasPosition(lsppos) {
		// First check in the node's children
		for _, chld := range nodes[pos].Args {
			if chld == EmptyNode {
				continue
			}
			if nodes[chld].Range.HasPosition(lsppos) {
				// log.Print(f.NodeDebug(chld))
				pos = chld
				continue search
			}

			// Then check in its list
			other := nodes[chld].Next
			for other != EmptyNode {
				// log.Print(f.NodeDebug(other))
				if nodes[other].Range.HasPosition(lsppos) {
					pos = other
					continue search
				}
				other = nodes[other].Next
			}

		}

		break
	}

	return pos, nil
}
