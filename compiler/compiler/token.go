package zoe

import (
	"fmt"
	"math"

	"github.com/sourcegraph/go-lsp"
)

type TokenKind uint32
type TokenPos uint32

//////////////////////////////////////

type TkRange struct {
	Start TokenPos
	End   TokenPos
}

func newTkRange() TkRange {
	return TkRange{Start: math.MaxUint32, End: math.MaxUint32}
}

func (t *TkRange) ExtendPos(p TokenPos) {
	if p == math.MaxUint32 {
		return
	}
	if t.Start > p {
		t.Start = p
	}
	if t.End < p {
		t.End = p
	}
}

func (t *TkRange) ExtendTk(tk *Parser) {
	t.ExtendPos(tk.pos)
}

func (t *TkRange) ExtendRange(rng TkRange) {
	t.ExtendPos(rng.Start)
	t.ExtendPos(rng.End)
}

///////////////////////////////////

// type Range struct {
// 	// should I include the source code as well ?
// 	Start     uint32
// 	End       uint32
// 	Line      uint32
// 	Column    uint32
// 	LineEnd   uint32
// 	ColumnEnd uint32
// }

// func (r Range) HasPosition(p *lsp.Position) bool {
// 	line := uint32(p.Line) // lsp is 0 based, but we're 1-based
// 	char := uint32(p.Character)

// 	if line < r.Line || line > r.LineEnd || line == r.Line && char < r.Column || line == r.LineEnd && char >= r.ColumnEnd {
// 		return false
// 	}

// 	return true
// }

// func (r Range) GetPosition() *Range {
// 	return &r
// }

// func (r Range) ToLspRange() *lsp.Range {
// 	return &lsp.Range{
// 		Start: lsp.Position{
// 			Line:      int(r.Line - 1),
// 			Character: int(r.Column - 1),
// 		},
// 		End: lsp.Position{
// 			Line:      int(r.LineEnd - 1),
// 			Character: int(r.ColumnEnd - 1),
// 		},
// 	}
// }

// func (r *Range) Extend(other Range) {
// 	if other.Line == 0 {
// 		// do not extend from a buggy range
// 		return
// 	}
// 	if r.Line == 0 {
// 		// take the other range as our own if we didn't exist
// 		*r = other
// 		return
// 	}

// 	if r.Line == other.Line {
// 		// if we're on the same line, the final column is the left-most one
// 		r.Column = minInt(r.Column, other.Column)
// 	} else {
// 		if other.Line < r.Line {
// 			r.Column = other.Column
// 			r.Line = other.Line
// 		}
// 	}
// 	if r.LineEnd == other.LineEnd {
// 		r.ColumnEnd = maxInt(r.ColumnEnd, other.ColumnEnd)
// 	} else {
// 		if other.LineEnd > r.LineEnd {
// 			r.LineEnd = other.LineEnd
// 			r.ColumnEnd = other.ColumnEnd
// 		}
// 	}

// 	r.Start = minInt(r.Start, other.Start)
// 	r.End = maxInt(r.End, other.End)
// }

// type Positioned interface {
// 	GetPosition() *Range
// }

type Token struct {
	Kind   TokenKind
	Offset uint32
	Length uint32
	Line   uint32
	Column uint32
}

type Parser struct {
	pos     TokenPos
	prev    TokenPos
	file    *File
	binding bindingPower
}

func (parser *Parser) AsRange() TkRange {
	if parser.IsEof() {
		return TkRange{Start: parser.pos, End: parser.pos}
	}
	return TkRange{Start: parser.pos, End: parser.pos + 1}
}

func (parser *Parser) Kind() TokenKind {
	return parser.ref().Kind
}

func (parser *Parser) Line() int {
	return int(parser.ref().Line)
}

func (parser *Parser) Column() int {
	return int(parser.ref().Column)
}

func (parser *Parser) Offset() int {
	return int(parser.ref().Offset)
}

func (parser *Parser) Length() int {
	return int(parser.ref().Length)
}

func (parser *Parser) Range() lsp.Range {
	return lsp.Range{
		Start: lsp.Position{
			Line:      parser.Line(),
			Character: parser.Column(),
		},
		End: lsp.Position{
			Line:      parser.Line(),
			Character: parser.Column() + parser.Length(),
		},
	}
}

func (parser *Parser) ref() *Token {
	return &parser.file.Tokens[parser.pos]
}

func (parser *Parser) IsEof() bool {
	return parser.ref().Kind == TK_EOF
	// return int(tk.pos) >= len(tk.file.Tokens)
}

func (parser *Parser) Is(tkind TokenKind) bool {
	if parser.IsEof() {
		return false
	}
	return parser.file.Tokens[parser.pos].Kind == tkind
}

func (parser *Parser) expectClosing(opening Parser, fn ...func()) {
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
	if !parser.Is(ckind) {
		opening.reportError("missing closing token")
		parser.reportError("missing closing token for '" + opening.GetText() + "'")
	}
	for _, f := range fn {
		f()
	}
	parser.Advance()
}

func (parser *Parser) shouldBe(kind TokenKind) bool {
	if !parser.Is(kind) {
		parser.reportError("expected " + tokstr[kind] + " but got '" + parser.GetText() + "'")
		return false
	}
	return true
}

func (parser *Parser) should(kind TokenKind) bool {
	if !parser.Is(kind) {
		parser.reportError("expected " + tokstr[kind] + " but got '" + parser.GetText() + "'")
		return false
	}
	return true
}

/*
	expect expects the provided token at the current position and
	calls fn on itself if found, and will advance the parser if
	the callback did not do so.
*/
func (parser *Parser) expect(kind TokenKind, fn func()) bool {
	if !parser.should(kind) {
		return false
	}
	var curpos = parser.pos
	fn()
	if curpos == parser.pos {
		parser.Advance()
	}
	return true
}

func (parser *Parser) consume(kind TokenKind) bool {
	return parser.expect(kind, func() {})
}

func (parser *Parser) ifToken(kind TokenKind, fn func()) bool {
	if !parser.Is(kind) {
		return false
	}
	fn()
	return true
}

func (parser *Parser) expectCommaIfNot(kind ...TokenKind) bool {
	var kd = parser.ref().Kind
	for _, k := range kind {
		if kd == k {
			return false
		}
	}
	parser.expect(TK_COMMA, func() {})
	return true
}

func (parser *Parser) GetText() string {
	if parser.IsEof() {
		return "<EOF>"
	}
	return parser.file.GetTokenText(parser.pos)
}

func (parser *Parser) isSkippable() bool {
	var kind = parser.file.Tokens[parser.pos].Kind
	return kind == TK_WHITESPACE || kind == TK_COMMENT
}

func (parser *Parser) Advance() {
	if parser.IsEof() {
		return
	}
	parser.prev = parser.pos
	parser.pos++
	for parser.isSkippable() {
		parser.pos++
	}
}

func (parser *Parser) advanceIf(kind TokenKind) bool {
	if parser.Is(kind) {
		parser.Advance()
		return true
	}
	return false
}

/////////////////////////////////////////////
////

var closingTokens = [TK__MAX]bool{}
var closerTokens = [TK__MAX]TokenKind{}

func init() {
	closerTokens[int(TK_LBRACE)] = TK_RBRACE
	closerTokens[int(TK_LBRACKET)] = TK_RBRACKET
	closerTokens[int(TK_LPAREN)] = TK_RPAREN
	closerTokens[int(TK_QUOTE)] = TK_QUOTE

	closingTokens[TK_RBRACE] = true
	closingTokens[TK_RPAREN] = true
	closingTokens[TK_RBRACKET] = true
}

// IsClosing is true if the current token is a closing token
// such as ), ] or }
func (parser *Parser) IsClosing() bool {
	if parser.IsEof() {
		return true // EOF closes, but it's usually an error
	}
	kind := parser.ref().Kind
	return closingTokens[kind]
}

func (parser *Parser) while(cond func() bool, fn func()) {
	var current = parser.pos
	for cond() {
		fn()
		if parser.pos == current {
			// We can't allow the parser to stay on the same token,
			// so we advance it. Most like, this is due to an error that
			// happened in `fn` so it already should have been reported.
			parser.Advance()
		}
	}
}

// Execute a function as long as the current token is not a closing token.
func (parser *Parser) whileNotClosing(fn func()) {
	parser.while(
		func() bool { return !parser.IsClosing() },
		fn,
	)
}

func (parser *Parser) parseEnclosedSeparatedByComma(fn func()) {
	var open = parser.Kind()
	var close = closerTokens[open]
	if close == 0 {
		// if we get here it means the compiler sent us here on another token than (, [ or {
		panic("this should not happen")
	}

	parser.Advance()
	parser.whileNotClosing(func() {
		fn()
		if !parser.Is(close) {
			parser.consume(TK_COMMA)
		}
	})

	parser.consume(close)
}

func (parser *Parser) parseEnclosed(fn func()) {
	var open = parser.Kind()
	var close = closerTokens[open]
	if close == 0 {
		// if we get here it means the compiler sent us here on another token than (, [ or {
		panic("this should not happen")
	}

	parser.Advance()

	parser.whileNotClosing(func() {
		fn()
	})

	parser.consume(close)
}

func (parser *Parser) whileNotEof(fn func()) {
	parser.while(
		func() bool { return !parser.IsEof() },
		fn,
	)
}

func (parser *Parser) whileNot(kind TokenKind, fn func()) {
	parser.while(
		func() bool { return !parser.IsEof() && !parser.Is(kind) },
		fn,
	)
}

// func (tk Tk) sym() *prattTk {
// 	return &syms[tk.ref().Kind]
// }

//////////////////////////////////////////////
//

func (parser *Parser) reportError(msg ...string) {
	parser.file.reportError(parser.Range(), msg...)
}

//////////////////////////////////////////////
//

func (parser *Parser) Debug() string {
	if parser.IsEof() {
		return "T[EOF]"
	}
	return fmt.Sprintf(`T[%s '%s' @%v:%v]`, tokstr[parser.ref().Kind], parser.GetText(), parser.Line()+1, parser.Column()+1)
}

func (parser *Parser) CreateIdentifier() *AstIdentifier {
	return parser.createAstIdentifier()
}
