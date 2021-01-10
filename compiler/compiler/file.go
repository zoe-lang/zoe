package zoe

import (
	"io"
	"io/ioutil"
	"os"
	"strings"
	"sync"

	"github.com/pkg/errors"
	"github.com/sourcegraph/go-lsp"
)

////////////////////////////////////////////////////
// FLOODGATE
////////////////////////////////////////////////////

/*
	FloodGate is a sync primitive whose purpose is to make sure
	a file does not read the contents of another while
	it is being processed.
*/
type FloodGate struct {
	isOpen bool
	cond   *sync.Cond
}

func NewFloodGate() *FloodGate {
	var m = sync.Mutex{}
	return &FloodGate{
		cond: sync.NewCond(&m),
	}
}

func (fg *FloodGate) Wait() {
	fg.cond.L.Lock()
	for !fg.isOpen {
		fg.cond.Wait()
	}
	fg.cond.L.Unlock()
}

func (fg *FloodGate) Open() {
	fg.cond.L.Lock()
	if !fg.isOpen {
		fg.isOpen = true
		fg.cond.Broadcast()
	}
	fg.cond.L.Unlock()
}

func (fg *FloodGate) Close() {
	fg.cond.L.Lock()
	fg.isOpen = false
	fg.cond.L.Unlock()
}

////////////////////////////////////////////////////
// COMPILER ERRORS
////////////////////////////////////////////////////

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

////////////////////////////////////////////////////
// FILE
////////////////////////////////////////////////////

// File holds the current parsing context
// A File instance may become obsolete as the editor sends new informations about it.
// Since type checking happens in its own goroutine, its result might not be needed anymore.
type File struct {
	Filename string

	Tokens []Token

	RootNode  Node
	RootScope *Scope
	Version   int

	DoneParsing *FloodGate

	Errors []ZoeError
	data   []byte

	DocCommentMap map[Node]TokenPos // node position => token position
}

// GetData returns the bytes of the file without the artifical null byte added
// during compilation.
func (f *File) GetData() []byte {
	if len(f.data) == 0 {
		return []byte{}
	}
	return f.data[0 : len(f.data)-1]
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
			if data[off] == '\n' {
				line = line + 1
				col = 0
			} else {
				col = col + 1
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
		DocCommentMap: make(map[Node]TokenPos),
		DoneParsing:   NewFloodGate(),
		// RootDocComments: make([]*Token, 0),
	}

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

func (f *File) SetComment(node Node, cmt TokenPos) {
	if node == nil {
		// This can happen for unrelated reasons
		return
	}
	f.DocCommentMap[node] = cmt
}

func (f *File) GetTokenText(tk TokenPos) string {
	var t = Parser{pos: tk, file: f}
	if t.IsEof() {
		return "<EOF>"
	}
	var tt = t.ref()

	return string(f.data[int(tt.Offset) : int(tt.Offset)+int(tt.Length)])
}

func (f *File) GetTkRangeBytes(rng TkRange) []byte {
	return f.data[int(f.Tokens[rng.Start].Offset):int(f.Tokens[int(rng.End)].Offset)]
}

func (f *File) GetTkRangeString(rng TkRange) string {
	return string(f.GetTkRangeBytes(rng))
}

func (f *File) reportError(rng lsp.Range, message ...string) {
	f.Errors = append(f.Errors, ZoeError{
		File:    f,
		Range:   rng,
		Message: strings.Join(message, ""),
	})
	// f.Errors[len(f.Errors)-1].Print(os.Stderr)
}
