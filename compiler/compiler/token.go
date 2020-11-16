package zoe

import (
	"encoding/json"
	"fmt"
	"io"
)

type TokenKind int

type Position struct {
	// should I include the source code as well ?
	Context   *ZoeContext
	Start     uint32
	End       uint32
	Line      uint32
	Column    uint32
	EndLine   uint32
	EndColumn uint32
}

func (p Position) GetPosition() *Position {
	return &p
}

func (p *Position) GetText() string {
	return string(p.Context.data[p.Start:p.End])
}

type Positioned interface {
	GetPosition() *Position
}

type Token struct {
	Position
	Kind   TokenKind
	Next   *Token
	WsNext *Token
}

func (t *Token) ToSlice() []string {
	res := make([]string, 0)
	for t != nil {
		res = append(res, t.String())
		t = t.Next
	}
	return res
}

func (t *Token) Dump(w io.Writer) {
	_, _ = w.Write([]byte(t.GetText()))
}

func (t *Token) GetPosition() *Position {
	return &t.Position
}

func (t *Token) panicIfNot(k TokenKind) {
	if t.Kind != k {
		panic(`requested ` + t.KindStr() + ` but got ` + tokstr[k])
	}
}

func (t *Token) MarshalJSON() ([]byte, error) {
	if t == nil {
		return []byte("<nil>"), nil
	}
	return json.Marshal(map[string]interface{}{
		"Value": t.String(),
	})
}

func (t *Token) Is(tk TokenKind) bool {
	if t != nil && t.Kind == tk {
		return true
	}
	return false
}

func (t *Token) IsSkippable() bool {
	kind := t.Kind
	return kind == TK_WHITESPACE || kind == TK_COMMENT
}

func (t *Token) String() string {
	if t == nil {
		return "<nil>"
	}
	if t.Kind < 0 {
		return "ERROR"
	}
	if t.Kind == TK_EOF {
		// return "EOF"
		return "*"
	}
	// return fmt.Sprintf("%#v", string(z.Context.data[z.Start:z.End]))
	return t.Position.GetText()
}

func (t *Token) Debug() string {
	return fmt.Sprint(t.String(), ":", t.KindStr())
	// if z == nil {
	// 	return "<nil nil>"
	// }
	// return z.KindStr() + " " + z.String()
}

func (t *Token) KindStr() string {
	if t.Kind == -1 {
		return "FAKE"
	}
	return tokstr[t.Kind]

}
