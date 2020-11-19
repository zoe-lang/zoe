package zoe

// Nodes are loaded into one big array to avoid too much work for the gc.

type NodePosition int
type AstNodeKind int
type Flag int

type NodeArray []AstNode

const (
	FLAG_LOCAL Flag = 1 << iota
)

// Identifiers are interned for faster lookup in a concurrent-safe interning store
// so that we can have namespaces that are maps of int32 (string id) => int32 (node position)

const (
	NODE_UNAOP AstNodeKind = 1 << 10
	NODE_BINOP AstNodeKind = 1 << 12
)

const (
	NODE_EMPTY     AstNodeKind = iota // "<>"
	NODE_ID                           // {{ str(n.Value) }}
	NODE_LITERAL                      // {{ lit(n.Value) }}
	NODE_DECL_FN                      // "decl-fn" signature block
	NODE_DECL_TYPE                    // "decl-type" template? expr
	NODE_DECL_NMSP                    // "decl-namespace" ident decl_block
	NODE_DECL_VAR                     // "decl-var" ident expr? expr?

	NODE_RETURN // "return"
	NODE_STRUCT // "struct"
	NODE_UNION  // "union"

	// node used for declaration blocks such as namespaces, files, and implement blocks
	NODE_DECL_BLOCK // "{{" <- "}}"
	// block contains code blocks
	NODE_BLOCK  // "{" <- "}" ?
	NODE_TUPLE  // "[" <- "]"
	NODE_ARGS   // "args"
	NODE_STRING // "str"

	NODE_FN        // "fn"
	NODE_SIGNATURE // "signature"
	NODE_TEMPLATE  // "template"
	NODE_IF        // "if"
	NODE_WHILE     // "while"

	NODE_LIT_NULL  // "null"
	NODE_LIT_VOID  // "void"
	NODE_LIT_FALSE // "false"
	NODE_LIT_TRUE  // "true"
	NODE_LIT_CHAR  //
	NODE_LIT_RAWSTR
	NODE_LIT_NUMBER
)

const (
	NODE_UNA_PLUS   AstNodeKind = NODE_UNAOP + iota // "+"
	NODE_UNA_MIN                                    // "-"
	NODE_UNA_NOT                                    // "!"
	NODE_UNA_BITNOT                                 // "~"

	NODE_BIN_ASSIGN    AstNodeKind = NODE_BINOP + iota // "="
	NODE_BIN_PLUS                                      // "+"
	NODE_BIN_MIN                                       // "-"
	NODE_BIN_DIV                                       // "/"
	NODE_BIN_MUL                                       // "*"
	NODE_BIN_MOD                                       // "%"
	NODE_BIN_EQ                                        // "=="
	NODE_BIN_NEQ                                       // "!="
	NODE_BIN_GTEQ                                      // ">="
	NODE_BIN_GT                                        // ">"
	NODE_BIN_LTEQ                                      // "<="
	NODE_BIN_LT                                        // "<"
	NODE_BIN_LSHIFT                                    // "<<"
	NODE_BIN_RSHIFT                                    // ">>"
	NODE_BIN_BITANDEQ                                  // "&="
	NODE_BIN_BITAND                                    // "&"
	NODE_BIN_BITOR                                     // "|"
	NODE_BIN_BITXOR                                    // "^"
	NODE_BIN_OR                                        // "||"
	NODE_BIN_AND                                       // "&&"
	NODE_BIN_IS                                        // "is"
	NODE_BIN_CAST                                      // "cast"
	NODE_BIN_CALL                                      // "call"
	NODE_BIN_TPLACCESS                                 // "tplcall"
	NODE_BIN_DOT                                       // "."
	NODE_BIN_NMSP                                      // "::"
)

//
const (
	LIT_FALSE uint32 = iota
	LIT_TRUE
	LIT_NULL
	LIT_VOID
)

type AstNode struct {
	Kind        AstNodeKind
	Range       Range        // the range inside the source file. an enclosing node updates its range according to its internal nodes
	IsIncorrect bool         // true if the node was tagged as being incorrect and thus should not be type checked
	Value       int          // can represent either a boolean (1 or 0), a node position, or a string id
	Next        NodePosition // The next node position as defined by its parent node
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
