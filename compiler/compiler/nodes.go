package zoe

// Nodes are loaded into one big array to avoid too much work for the gc.
// An AST is not guaranteed to be correct.
// The type checker validates the AST as it type checks, and marks the IsError
// field of the nodes it traverses if it encounters errors.

type NodePosition int
type AstNodeKind int
type Flag int

const EmptyNode NodePosition = 0

type NodeArray []AstNode

const (
	FLAG_LOCAL Flag = 1 << iota
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
	NODE_TYPE          // bblue("type") 			::: name template typeexp
	NODE_NAMESPACE     // bblue("namespace")  ::: name block
	NODE_VAR           // bblue("var")        ::: name typeexp assign
	NODE_SIGNATURE     // "signature"         ::: template args rettype
	NODE_RETURN        // "return" 						::: exp
	NODE_STRUCT        // "struct"            ::: varlist
	NODE_UNION         // "union"							::: members
	NODE_STRING        // "str" ::: contents
	NODE_ARRAY_LITERAL // "array" ::: contents
	NODE_IF            // "if"                      ::: cond thenarm elsearm
	NODE_WHILE         // "while"
	NODE_IMPORT        // bblue("import")
	NODE_UNA_ELLIPSIS  // "..."
	NODE_UNA_PLUS      // "+"
	NODE_UNA_MIN       // "-"
	NODE_UNA_NOT       // "!"
	NODE_UNA_BITNOT    // "~"
	NODE_BIN_ASSIGN    // "="
	NODE_BIN_PLUS      // "+"
	NODE_BIN_MIN       // "-"
	NODE_BIN_DIV       // "/"
	NODE_BIN_MUL       // "*"
	NODE_BIN_MOD       // "%"
	NODE_BIN_EQ        // "=="
	NODE_BIN_NEQ       // "!="
	NODE_BIN_GTEQ      // ">="
	NODE_BIN_GT        // ">"
	NODE_BIN_LTEQ      // "<="
	NODE_BIN_LT        // "<"
	NODE_BIN_LSHIFT    // "<<"
	NODE_BIN_RSHIFT    // ">>"
	NODE_BIN_BITANDEQ  // "&="
	NODE_BIN_BITAND    // "&"
	NODE_BIN_BITOR     // "|"
	NODE_BIN_BITXOR    // "^"
	NODE_BIN_OR        // "||"
	NODE_BIN_AND       // "&&"
	NODE_BIN_IS        // "is"
	NODE_BIN_CAST      // "cast"
	NODE_BIN_CALL      // "call"
	NODE_BIN_INDEX     // "index"
	NODE_BIN_DOT       // "."

	NODE_LIT_NULL   // mag("null")
	NODE_LIT_VOID   // mag("void")
	NODE_LIT_FALSE  // mag("false")
	NODE_LIT_TRUE   // mag("true")
	NODE_LIT_CHAR   // green(f.GetRangeText(n.Range))
	NODE_LIT_RAWSTR // green("'",f.GetRangeText(n.Range),"'")
	NODE_LIT_NUMBER // mag(f.GetRangeText(n.Range))
	NODE_ID         // cyan(internedIds.Get(n.Value))

	NODE__SIZE
)

type AstNode struct {
	Kind        AstNodeKind
	Range       Range // the range inside the source file. an enclosing node updates its range according to its internal nodes
	Scope       ScopePosition
	IsIncorrect bool // true if the node was tagged as being incorrect and thus should not be type checked
	Value       int  // can represent either a boolean (1 or 0), a node position, or a string id
	ArgLen      int
	Args        [4]NodePosition // probably unused
	Next        NodePosition    // The next node position as defined by its parent node when inside a list (tuples, template params or blocks)
}

// op: Node(NODE_BINOP_PLUS, value: idx:FIRST) first(TYPE_IDENT, next: second) second(NODE_LIT_INT, next: 0)

// SetIncorrect marks the node as being incorrect. The type checker should ignore these as this flag
// is set if the underlying node representation makes no sense.
func (a *AstNode) SetIncorrect() {
	a.IsIncorrect = true
}

func (a *AstNode) SetFlag(value Flag) {
	a.Value &= int(value)
}

func (a *AstNode) HasFlag(value Flag) bool {
	return a.Value&int(value) != 0
}
