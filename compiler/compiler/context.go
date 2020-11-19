package zoe

import (
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/pkg/errors"
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

// File holds the current parsing context
// also does the error handling stuff.
type File struct {
	Filename    string
	Tokens      []Token
	Nodes       NodeArray
	RootNodePos NodePosition
	Errors      []ZoeError
	data        []byte

	current *Token
	tkpos   uint32

	DocCommentMap map[NodePosition]TokenPos // node position => token position
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

	ctx := &File{
		Filename:      filename,
		Errors:        make([]ZoeError, 0),
		data:          append(data, '\x00'),
		DocCommentMap: make(map[NodePosition]TokenPos),
		// RootDocComments: make([]*Token, 0),
	}

	lxerr := ctx.Lex()
	if lxerr != nil {
		return ctx, errors.Wrap(lxerr, "lexing failed")
	}

	return ctx, nil
}

func (c *File) GetTokenText(tk TokenPos) string {
	return c.GetRangeText(c.Tokens[tk].Range)
}

func (c *File) GetRangeText(p Range) string {
	return string(c.data[p.Start:p.End])
}

func (c *File) GetNodeText(n NodePosition) string {
	return c.GetRangeText(c.Nodes[n].Range)
}

func (c *File) isEof() bool {
	return c.current == nil
}

func (c *File) reportError(pos Positioned, message ...string) {
	c.Errors = append(c.Errors, ZoeError{
		File:    c,
		Range:   *pos.GetPosition(),
		Message: strings.Join(message, ""),
	})
}

func (f *File) createNodeBuilder() nodeBuilder {
	return nodeBuilder{
		file:          f,
		tokens:        f.Tokens,
		tokensLen:     TokenPos(len(f.Tokens)),
		doccommentMap: f.DocCommentMap,
	}
}
