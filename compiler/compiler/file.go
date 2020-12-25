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
	Range   lsp.Range
	Message string
}

func (err ZoeError) Print(w io.Writer) {
	fred.Fprint(w, err.File.Filename+" ")
	fgreen.Fprint(w, err.Range.Start.Line+1)
	_, _ = w.Write([]byte(": " + err.Message + "\n"))
}

func (err ZoeError) ToLspDiagnostic() lsp.Diagnostic {
	d := lsp.Diagnostic{}
	d.Message = err.Message
	d.Range = err.Range
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

func (f *File) GetData() []byte {
	return f.data
}

func (f *File) GetOffsetForPosition(pos lsp.Position) int {
	// We're going to use the token positions to find the
	// real offset, since they're ordered.
	var tokens = f.Tokens
	var tkpos = len(tokens) / 2
	var size = len(tokens) / 2 // Size is going to be divided by 2 all the time

	// log.Print(pos)
	var data = f.data
	for {
		var tk = tokens[tkpos]

		if tkpos >= len(tokens)-1 {
			// log.Print(tkpos, ` (`, size, `) `, tk.Line, tk.Column, " --> ", int(tk.Offset)+(pos.Character-int(tk.Column)))
			return int(tk.Offset) + (pos.Character - int(tk.Column))
		}

		var ntk = tokens[tkpos+1]
		size = size / 2
		if size == 0 {
			size = 1
		}

		// log.Print(tkpos, ` (`, size, `) `, tk.Line, tk.Column, " - ", ntk.Line, ntk.Column)

		// We're before the current token
		if pos.Line < int(tk.Line) || (pos.Line == int(tk.Line) && pos.Character < int(tk.Column)) {
			tkpos = tkpos - size
			continue
		}

		// We're after the current token
		if pos.Line > int(ntk.Line) || (pos.Line == int(ntk.Line) && pos.Character >= int(ntk.Column)) {
			tkpos = tkpos + size
			continue
		}

		// We're on the right token, but the token might span lines, so now we go until the position
		// manually
		var off = int(tk.Offset)
		var line = int(tk.Line)
		var col = int(tk.Column)
		for line < pos.Line || col < pos.Character {
			col = col + 1
			if data[off] == '\n' {
				line = line + 1
				col = 0
			}
			off = off + 1
		}
		return off

	}
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
	var t = Tk{pos: tk, file: f}
	if t.IsEof() {
		return "<EOF>"
	}
	var tt = t.ref()

	return string(f.data[int(tt.Offset) : int(tt.Offset)+int(tt.Length)])
}

func (f *File) GetNodeText(np NodePosition) string {
	var n = np.Node(f)
	var rng = n.ref().Range
	return string(f.data[int(f.Tokens[rng.Start].Offset):int(f.Tokens[int(rng.End)].Offset)])
}

func (f *File) reportError(rng lsp.Range, message ...string) {
	f.Errors = append(f.Errors, ZoeError{
		File:    f,
		Range:   rng,
		Message: strings.Join(message, ""),
	})
	// f.Errors[len(f.Errors)-1].Print(os.Stderr)
}

func (f *File) createNode(tk Tk, kind AstNodeKind, scope Scope, children ...Node) Node {
	// maybe we should handle here the capacity of the node arrays ?
	l := NodePosition(len(f.Nodes))
	f.Nodes = append(f.Nodes, AstNode{Kind: kind, Range: TkRange{Start: tk.pos, End: tk.pos + 1}, Scope: scope.pos})

	cl := len(children)
	if cl > 0 {
		node := &f.Nodes[l]
		node.ArgLen = int8(cl)

		for i, chld := range children {
			node.Args[i] = chld.pos
			if !chld.IsEmpty() {
				node.Range.ExtendNode(chld)
			}
		}
	}

	return Node{
		file: f,
		pos:  l,
	}
}

// Find a node that matches a given range
func (f *File) FindNodePosition(lsppos lsp.Position) ([]Node, error) {
	node := f.RootNode
	var path = make([]Node, 0)
	path = append(path, node)
	// nodes := f.Nodes

	// log.Print(lsppos.Line+1, ":", lsppos.Character+1)
search:
	for node.HasPosition(lsppos) {
		// log.Print(node.Debug())
		// check in the node's children
		for _, chl := range node.ref().Args {
			chld := chl.Node(f)
			if chld.IsEmpty() {
				continue
			}
			if chld.HasPosition(lsppos) {
				// log.Print(f.NodeDebug(chld))
				node = chld
				path = append(path, node)
				continue search
			}

			// Then check in its list
			other := chld.Next()
			for !other.IsEmpty() {
				// log.Print(f.NodeDebug(other))
				if other.HasPosition(lsppos) {
					node = other
					path = append(path, node)
					continue search
				}
				other = other.Next()
			}

		}

		break
	}

	return path, nil
}
