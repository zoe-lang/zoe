package zoe

import "github.com/sourcegraph/go-lsp"

// Nodes are loaded into one big array to avoid too much work for the gc.
// An AST is not guaranteed to be correct.
// The type checker validates the AST as it type checks, and marks the IsError
// field of the nodes it traverses if it encounters errors.

type NodePosition int

func (n NodePosition) Node(file *File) Node {
	return Node{
		pos:  n,
		file: file,
	}
}

type AstNodeKind int
type Flag int

var EmptyNode = Node{}

const (
	FLAG_LOCAL Flag = 1 << iota
	FLAG_CONST
	FLAG_BLOCK_ENDS_EXECUTION // set when return is the last instruction of a block
	// used to know if an if when it has no else block sets the rest of the block as its else condition
)

// Identifiers are interned for faster lookup in a concurrent-safe interning store
// so that we can have namespaces that are maps of int32 (string id) => int32 (node position)

const (
	NODE_EMPTY AstNodeKind = iota // grey("~")

	// node used for declaration blocks such as namespaces, files, and implement blocks
	NODE_FILE // "file" ::: contents
	// block contains code blocks
	NODE_BLOCK // "block" ::: contents
	NODE_TUPLE // "tuple" ::: contents

	NODE_FN            // bblue("fn") 				::: name signature definition
	NODE_METHOD        // bblue("method")     ::: name signature definition
	NODE_TYPE          // bblue("type") 			::: name template typeexp block
	NODE_NAMESPACE     // bblue("namespace")  ::: name block
	NODE_VAR           // bblue("var")        ::: name typeexp assign
	NODE_SIGNATURE     // "signature"         ::: template args rettype
	NODE_RETURN        // "return" 						::: exp
	NODE_ENUM          // bblue("enum")       ::: varlist
	NODE_STRUCT        // bblue("struct")     ::: template varlist
	NODE_TRAIT         // bblue("trait")      ::: template methodlist
	NODE_UNION         // "union"							::: members
	NODE_ISO_BLOCK     // "iso"               ::: block
	NODE_ISO_TYPE      // "iso-type"          ::: type_expr
	NODE_STRING        // "str" ::: contents
	NODE_ARRAY_LITERAL // "array" ::: contents
	NODE_IF            // "if"                      ::: cond thenarm elsearm
	NODE_SWITCH        // "switch" ::: exp arms
	NODE_SWITCH_ARM    // "arm" ::: cond block
	NODE_FOR           // "for"   ::: vardecl rng block
	NODE_WHILE         // "while" ::: cond block
	NODE_IMPORT        // bblue("import") ::: module id exp
	NODE_IMPLEMENT     // bblue("implement") ::: id block

	NODE_UNA_ELLIPSIS // "..."
	NODE_UNA_NOT      // "!" ::: exp
	NODE_UNA_POINTER  // "ptr" ::: pointed
	NODE_UNA_REF      // "ref" ::: variable

	NODE_BIN_ASSIGN // "="
	NODE_BIN_PLUS   // "+"
	NODE_BIN_MIN    // "-"
	NODE_BIN_DIV    // "/"
	NODE_BIN_MUL    // "*"
	NODE_BIN_MOD    // "%"
	NODE_BIN_EQ     // "=="
	NODE_BIN_NEQ    // "!="
	NODE_BIN_GTEQ   // ">="
	NODE_BIN_GT     // ">"
	NODE_BIN_LTEQ   // "<="
	NODE_BIN_LT     // "<"
	NODE_BIN_LSHIFT // "<<"
	NODE_BIN_RSHIFT // ">>"
	NODE_BIN_BITAND // "&"
	NODE_BIN_BITOR  // "|"
	NODE_BIN_BITXOR // "^"
	NODE_BIN_OR     // "||"
	NODE_BIN_AND    // "&&"
	NODE_BIN_IS     // "is"
	NODE_BIN_CAST   // "cast"
	NODE_BIN_CALL   // "call"
	NODE_BIN_INDEX  // "index"
	NODE_BIN_DOT    // "."

	NODE_LIT_NONE   // mag("null")
	NODE_LIT_THIS   // mag("null")
	NODE_LIT_VOID   // mag("void")
	NODE_LIT_FALSE  // mag("false")
	NODE_LIT_TRUE   // mag("true")
	NODE_LIT_CHAR   // green(n.GetText())
	NODE_LIT_RAWSTR // green("'",n.GetText(),"'")
	NODE_LIT_NUMBER // mag(n.GetText())
	NODE_INTEGER    // mag(n.GetValue())
	NODE_ID         // cyan(GetInternedString(n.InternedString()))

	NODE__SIZE
)

type AstNode struct {
	Kind        AstNodeKind
	Range       TkRange // the range inside the source file. an enclosing node updates its range according to its internal nodes
	Scope       ScopePosition
	IsIncorrect bool // true if the node was tagged as being incorrect and thus should not be type checked
	Value       int  // can represent either a boolean (1 or 0), a node position, or a string id
	ArgLen      int8
	Args        [4]NodePosition // probably unused
	Next        NodePosition    // The next node position as defined by its parent node when inside a list (tuples, template params or blocks)
}

type Node struct {
	pos  NodePosition
	file *File
}

// op: Node(NODE_BINOP_PLUS, value: idx:FIRST) first(TYPE_IDENT, next: second) second(NODE_LIT_INT, next: 0)

func (n Node) ref() *AstNode {
	return &n.file.Nodes[n.pos]
}

func (n Node) IsEmpty() bool {
	return n.pos == 0
}

// SetIncorrect marks the node as being incorrect. The type checker should ignore these as this flag
// is set if the underlying node representation makes no sense.
func (n Node) SetIncorrect() {
	n.ref().IsIncorrect = true
}

// Gets the scope the node resides in
func (n Node) Scope() Scope {
	return Scope{
		pos:  n.ref().Scope,
		file: n.file,
	}
}

func (n Node) SetFlag(value Flag) {
	n.ref().Value &= int(value)
}

func (n Node) HasFlag(value Flag) bool {
	return n.ref().Value&int(value) != 0
}

func (n Node) SetNext(other Node) {
	n.ref().Next = other.pos
}

func (n Node) HasNext() bool {
	return !n.Next().IsEmpty()
}

func (n Node) Next() Node {
	return Node{
		pos:  n.ref().Next,
		file: n.file,
	}
}

func (n Node) Clone() Node {
	orig := n.ref()
	respos := NodePosition(len(n.file.Nodes))
	n.file.Nodes = append(n.file.Nodes, *orig)
	n.file.Nodes[respos].Next = 0
	res := Node{
		pos:  respos,
		file: n.file,
	}
	return res
}

///////////////////////////////////
// func (n Node) Range() Range {
// 	var ref = n.ref()
// 	var tks = n.file.Tokens[ref.Range.Start]
// 	var tke = n.file.Tokens[ref.Range.End]
// 	return Range{
// 		Start:     tks.Offset,
// 		End:       tke.Offset,
// 		Line:      tks.Line,
// 		LineEnd:   tke.Line,
// 		Column:    tks.Column,
// 		ColumnEnd: tke.Column,
// 	}
// }

func PositionInRange(pos lsp.Position, rng lsp.Range) bool {
	var st = rng.Start
	var ed = rng.End
	return pos.Line >= st.Line && pos.Line < ed.Line &&
		(pos.Line != st.Line || pos.Character >= st.Character) &&
		(pos.Line != ed.Line || pos.Character < ed.Character)
}

func (n Node) HasPosition(lsppos lsp.Position) bool {
	var rng = n.Range()
	return PositionInRange(lsppos, rng)
}

func (n Node) Range() lsp.Range {
	var t = n.ref().Range
	var tks = n.file.Tokens
	var st = tks[int(t.Start)]
	var ed = tks[int(t.End)]
	return lsp.Range{
		Start: lsp.Position{
			Line:      int(st.Line),
			Character: int(st.Column),
		},
		End: lsp.Position{
			Line:      int(ed.Line),
			Character: int(ed.Column),
		},
	}
}

func (n Node) ReportError(msg ...string) {
	n.file.reportError(n.Range(), msg...)
}

func (n Node) SetValue(v int) {
	n.ref().Value = v
}

func (n Node) GetValue() int {
	return n.ref().Value
}

func (n Node) SetArgs(args ...Node) {
	node := n.ref()
	node.ArgLen = int8(len(args))
	for i, a := range args {
		node.Args[i] = a.pos
		n.ExtendRangeFromNode(a)
	}
}

func (n Node) ArgLen() int {
	return int(n.ref().ArgLen)
}

func (n Node) GetArg(nb int) Node {
	return n.ref().Args[nb].Node(n.file)
}

func (n Node) Is(nk AstNodeKind) bool {
	if n == EmptyNode {
		return false
	}
	return n.ref().Kind == nk
}

func (n Node) IsAnyOf(nk ...AstNodeKind) bool {
	for _, nk := range nk {
		if n.Is(nk) {
			return true
		}
	}
	return false
}

func (n Node) expect(nk ...AstNodeKind) bool {
	if !n.IsAnyOf(nk...) {
		return false
	}
	return true
}

func (n Node) Kind() AstNodeKind {
	return n.ref().Kind
}

func (n Node) GetText() string {
	return n.file.GetNodeText(n.pos)
}

func (n Node) InternedString() InternedString {
	// FIXME: should we check to make sure we're holding an interned Id ?
	return InternedString(n.ref().Value)
}

func (n Node) SetInternedString(val InternedString) {
	n.ref().Value = int(val)
}

func (n Node) Extend(other Tk) {
	if n == EmptyNode {
		return
	}
	var ref = n.ref()
	if ref.Range.End < other.pos {
		ref.Range.End = other.pos
	}
	if ref.Range.Start > other.pos {
		ref.Range.Start = other.pos
	}
}

func (n Node) ExtendRangeFromNode(other Node) {
	if !other.IsEmpty() {
		var oref = other.ref()
		n.Extend(Tk{pos: oref.Range.Start})
		n.Extend(Tk{pos: oref.Range.End})
	}
}
