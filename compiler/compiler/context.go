package zoe

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

type ZoeError struct {
	Position Position
	Message  string
}

func (err ZoeError) Print(w io.Writer) {
	fred.Fprint(w, err.Position.Context.Filename+" ")
	fgreen.Fprint(w, err.Position.Line)
	_, _ = w.Write([]byte(": " + err.Message + "\n"))
}

// ZoeContext holds the current parsing context
// also does the error handling stuff.
type ZoeContext struct {
	Start    *Token
	End      *Token
	Filename string
	Current  *Token
	Errors   []ZoeError
	data     []byte
	Root     *Node

	DocCommentMap   map[Node]*Token // contains the doc comments related to given nodes
	RootDocComments []*Token
}

func NewZoeContext(filename string) (*ZoeContext, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	data, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	ctx := &ZoeContext{
		Filename:        filename,
		Errors:          make([]ZoeError, 0),
		data:            append(data, '\x00'),
		DocCommentMap:   make(map[Node]*Token),
		RootDocComments: make([]*Token, 0),
	}

	lxerr := ctx.Lex()
	if lxerr != nil {
		return nil, lxerr
	}

	ctx.Current = ctx.Start
	if ctx.Current.IsSkippable() {
		ctx.advance()
	}

	return ctx, nil
}

func (c *ZoeContext) currentSym() (*Token, *prattTk) {
	res := c.Current
	return res, &syms[res.Kind]
}

func (c *ZoeContext) isEof() bool {
	return c.Current == nil
}

func (c *ZoeContext) advance() {
	if c.Current != nil {
		c.Current = c.Current.NextMeaningfulToken()
	}
}

// At the top level, just parse everything we can
func (c *ZoeContext) ParseFile() *Node {
	res := parseUntil(c, NODE_DECLS, &Token{Kind: TK_EOF}, TK_EOF, 0)
	c.Root = res
	return res
}

func (c *ZoeContext) reportErrorAtCurrentPosition(message ...string) {
	var pos Position
	if c.Current == nil {
		pos = c.End.Position
	} else {
		pos = c.Current.Position
	}
	c.reportError(pos, message...)
}

func (c *ZoeContext) reportError(pos Position, message ...string) {
	c.Errors = append(c.Errors, ZoeError{
		Position: pos,
		Message:  strings.Join(message, ""),
	})
}

func (c *ZoeContext) nodeError(tk *Token, left ...*Node) *Node {
	c.reportError(tk.Position, fmt.Sprintf(`unexpected '%s'`, tk.String()))
	return NewErrorNode(tk, left...)
}
