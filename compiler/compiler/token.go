package zoe

import "github.com/sourcegraph/go-lsp"

type TokenKind int
type TokenPos int

type Range struct {
	// should I include the source code as well ?
	Start     uint32
	End       uint32
	Line      uint32
	Column    uint32
	LineEnd   uint32
	ColumnEnd uint32
}

func (r Range) HasPosition(p *lsp.Position) bool {
	line := uint32(p.Line + 1) // lsp is 0 based, but we're 1-based
	char := uint32(p.Character + 1)

	if line < r.Line || line > r.LineEnd || line == r.Line && char < r.Column || line == r.LineEnd && char >= r.ColumnEnd {
		return false
	}

	return true
}

func (r Range) GetPosition() *Range {
	return &r
}

func (r Range) ToLspRange() *lsp.Range {
	return &lsp.Range{
		Start: lsp.Position{
			Line:      int(r.Line - 1),
			Character: int(r.Column - 1),
		},
		End: lsp.Position{
			Line:      int(r.LineEnd - 1),
			Character: int(r.ColumnEnd - 1),
		},
	}
}

func (r *Range) Extend(other Range) {
	if other.Line == 0 {
		// do not extend from a buggy range
		return
	}
	if r.Line == 0 {
		// take the other range as our own if we didn't exist
		*r = other
		return
	}

	if r.Line == other.Line {
		// if we're on the same line, the final column is the left-most one
		r.Column = minInt(r.Column, other.Column)
	} else {
		if other.Line < r.Line {
			r.Column = other.Column
			r.Line = other.Line
		}
	}
	if r.LineEnd == other.LineEnd {
		r.ColumnEnd = maxInt(r.ColumnEnd, other.ColumnEnd)
	} else {
		if other.LineEnd > r.LineEnd {
			r.LineEnd = other.LineEnd
			r.ColumnEnd = other.ColumnEnd
		}
	}

	r.Start = minInt(r.Start, other.Start)
	r.End = maxInt(r.End, other.End)
}

type Positioned interface {
	GetPosition() *Range
}

type Token struct {
	Kind TokenKind
	Range
}

func (t Token) getSym() *prattTk {
	return &syms[t.Kind]
}

func (t *Token) panicIfNot(k TokenKind) {
	if t.Kind != k {
		panic(`requested ` + t.KindStr() + ` but got ` + tokstr[k])
	}
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

func (t *Token) KindStr() string {
	if t.Kind == -1 {
		return "FAKE"
	}
	return tokstr[t.Kind]

}
