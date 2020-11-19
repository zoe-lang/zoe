package zoe

type TokenKind int
type TokenPos int

type Range struct {
	// should I include the source code as well ?
	Start  uint32
	End    uint32
	Line   uint32
	Column uint32
}

func (p Range) GetPosition() *Range {
	return &p
}

func (r *Range) Extend(other Range) {
	if other.Line == 0 {
		return
	}
	if r.Line == 0 {
		*r = other
		return
	}

	if r.Line == other.Line {
		r.Column = minInt(r.Column, other.Column)
	} else {
		r.Column = other.Column
	}
	r.Start = minInt(r.Start, other.Start)
	r.Line = minInt(r.Line, other.Line)
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
