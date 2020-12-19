package zoe

import (
	"fmt"

	"github.com/sourcegraph/go-lsp"
)

type TokenKind uint32
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
	Kind   TokenKind
	Offset uint32
	Length uint32
	Line   uint32
	Column uint32
}

type Tk struct {
	pos  TokenPos
	file *File
}

func (tk Tk) ref() *Token {
	return &tk.file.Tokens[tk.pos]
}

func (tk Tk) IsEof() bool {
	return tk.ref().Kind == TK_EOF
	// return int(tk.pos) >= len(tk.file.Tokens)
}

func (tk Tk) Is(tkind TokenKind) bool {
	if tk.IsEof() {
		return false
	}
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

func (tk Tk) expectClosing(opening Tk, fn ...func(tk Tk)) Tk {
	var okind = opening.ref().Kind
	var ckind TokenKind
	switch okind {
	case TK_LPAREN:
		ckind = TK_RPAREN
	case TK_LBRACKET:
		ckind = TK_RBRACKET
	case TK_LBRACE:
		ckind = TK_RBRACE
	default:
		panic(tokstr[okind] + " has no corresponding closing token, this is a compiler bug")
	}
	if !tk.Is(ckind) {
		opening.reportError("missing closing token")
		tk.reportError("missing closing token for '" + opening.GetText() + "'")
		return tk
	}
	for _, f := range fn {
		f(tk)
	}
	return tk.Next()
}

func (tk Tk) shouldBe(kind TokenKind) bool {
	if !tk.Is(kind) {
		tk.reportError("expected " + tokstr[kind] + " but got '" + tk.GetText() + "'")
		return false
	}
	return true
}

func (tk Tk) expect(kind TokenKind, fn ...func(tk Tk)) (Tk, bool) {
	if !tk.Is(kind) {
		tk.reportError("expected " + tokstr[kind] + " but got '" + tk.GetText() + "'")
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

func (tk Tk) Range() Range {
	var next Tk
	if tk.IsEof() {
		next = tk
	} else {
		next = tk.Next()
	}

	var ref = tk.ref()
	var nextref = next.ref()
	return Range{
		Start:     ref.Offset,
		End:       nextref.Offset,
		Line:      ref.Line,
		Column:    ref.Column,
		LineEnd:   nextref.Line,
		ColumnEnd: nextref.Column,
	}
}

func (tk Tk) GetText() string {
	if tk.IsEof() {
		return "<EOF>"
	}
	var tokens = tk.file.Tokens
	var t = tokens[tk.pos]

	return string(tk.file.data[int(t.Offset) : int(t.Offset)+int(t.Length)])
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
	closingTokens[TK_RBRACE] = true
	closingTokens[TK_RPAREN] = true
	closingTokens[TK_RBRACKET] = true
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

func (tk Tk) while(cond func(iter Tk) bool, fn func(iter Tk) Tk) Tk {
	var iter = tk
	for cond(iter) {
		var res = fn(iter)
		if res.pos == iter.pos {
			// We can't allow the parser to stay on the same token,
			// so we advance it. Most like, this is due to an error that
			// happened in `fn` so it already should have been reported.
			iter = iter.Next()
		} else {
			iter = res
		}
	}
	return iter
}

// Execute a function as long as the current token is not a closing token.
func (tk Tk) whileNotClosing(fn func(iter Tk) Tk) Tk {
	return tk.while(
		func(iter Tk) bool { return !iter.IsClosing() },
		fn,
	)
}

func (tk Tk) whileNotEof(fn func(iter Tk) Tk) Tk {
	return tk.while(
		func(iter Tk) bool { return !iter.IsEof() },
		fn,
	)
}

func (tk Tk) whileNot(kind TokenKind, fn func(iter Tk) Tk) Tk {
	return tk.while(
		func(iter Tk) bool { return !iter.IsEof() && !iter.Is(kind) },
		fn,
	)
}

func (tk Tk) sym() *prattTk {
	return &syms[tk.ref().Kind]
}

//////////////////////////////////////////////
//

func (tk Tk) reportError(msg ...string) {
	tk.file.reportError(tk.Range(), msg...)
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

func (tk Tk) Debug() string {
	if tk.IsEof() {
		return "T[EOF]"
	}
	return fmt.Sprintf(`T[%s '%s' @%v:%v]`, tokstr[tk.ref().Kind], tk.GetText(), tk.Range().Line, tk.Range().Column)
}
