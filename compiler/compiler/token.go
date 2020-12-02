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

type Tk struct {
	pos  TokenPos
	file *File
}

func (tk Tk) ref() *Token {
	return &tk.file.Tokens[tk.pos]
}

func (tk Tk) Range() Range {
	return tk.ref().Range
}

func (tk Tk) IsEof() bool {
	return int(tk.pos) >= len(tk.file.Tokens)
}

func (tk Tk) Is(tkind TokenKind) bool {
	return tk.file.Tokens[tk.pos].Kind == tkind
}

func (tk Tk) Peek(kind ...TokenKind) bool {
	n := tk.Next()
	if n.IsEof() {
		return false
	}
	for _, k := range kind {
		if n.Is(k) {
			return true
		}
	}
	return false
}

func (tk Tk) consume(kind TokenKind, fn ...func(tk Tk)) (Tk, bool) {
	if !tk.Is(kind) {
		return tk, false
	}
	next := tk.Next()
	for _, f := range fn {
		f(tk)
	}
	return next, true
}

func (tk Tk) expect(kind TokenKind, fn ...func(tk Tk)) (Tk, bool) {
	if !tk.Is(kind) {
		tk.reportError("expected " + tokstr[kind] + " but got '" + "'")
		return tk, false
	}
	next := tk.Next()
	for _, f := range fn {
		f(tk)
	}
	return next, true
}

func (tk Tk) expectCommaIfNot(kind ...TokenKind) Tk {
	var kd = tk.ref().Kind
	for _, k := range kind {
		if kd == k {
			return tk // we don't move
		}
	}
	tk, _ = tk.expect(TK_COMMA)
	return tk
}

func (tk Tk) GetText() string {
	r := tk.ref().Range
	return string(tk.file.data[r.Start:r.End])
}

func (tk Tk) Next() Tk {
	return Tk{
		pos:  tk.pos + 1,
		file: tk.file,
	}
}

/////////////////////////////////////////////
////

var closingTokens = [TK__MAX]bool{}

func init() {
	for i := range closingTokens {
		closingTokens[i] = true
	}
	closingTokens[TK_RBRACE] = false
	closingTokens[TK_RPAREN] = false
	closingTokens[TK_RBRACKET] = false
}

// IsClosing is true if the current token is a closing token
// such as ), ] or }
func (tk Tk) IsClosing() bool {
	if tk.IsEof() {
		return true // EOF closes, but it's usually an error
	}
	kind := tk.ref().Kind
	return closingTokens[kind]
}

func (tk Tk) sym() *prattTk {
	return &syms[tk.ref().Kind]
}

//////////////////////////////////////////////
//

func (tk Tk) reportError(msg ...string) {
	tk.file.reportError(tk.ref().Range, msg...)
}

//////////////////////////////////////////////
//

func (tk Tk) createNode(scope Scope, nk AstNodeKind, args ...Node) Node {
	return tk.file.createNode(tk.Range(), nk, scope, args...) // ????
}

func (tk Tk) createIdNode(scope Scope) Node {
	idstr := SaveInternedString(tk.GetText())
	idnode := tk.createNode(scope, NODE_ID)
	idnode.SetInternedString(idstr)
	// b.file.Nodes[idnode].Value = idstr
	return idnode
}

func (tk Tk) createBinOp(scope Scope, kind AstNodeKind, left Node, right Node) Node {
	return tk.createNode(scope, kind, left, right)
}

func (tk Tk) createUnaryOp(scope Scope, kind AstNodeKind, left Node) Node {
	return tk.createNode(scope, kind, left)
}
