package zoe

import (
	"fmt"
	"io"
	"strings"
	"unicode"
	"unicode/utf8"
)

/*
	Nodes represent a (very) simplified parse tree that in the end kind of would look
	like lisp.
	The first argument is always a token from the source file, except when they are
	artificially inject by the parsing process, in which case their Kind will be
	the only

	Functions are always represented as
	- (fn (=> (-> (...args) return-type) block))
	or in its generic variant
	- (fn ([ (<generic args>) (=> (-> ...args) return-type)))

	variables probably should be parsed differently from the rest
	What about tuple assignment ?
	- (var <name> <type> <exp>)

	this is probably not what we want
	- (= (var <name>) <exp>)
	- (= (var <name> <type>) <exp>)

	types
	- (is (type <name>) <exp> (implements (...)) )

	All other expressions are
*/

type NodeKind string

const (
	NODE_TERMINAL = "terminal" // Id, Easy strings, Number, Boolean
	NODE_ERROR    = "error"

	NODE_FNDECL    = "decl:fn"
	NODE_VARDECL   = "decl:var"
	NODE_TYPEDECL  = "decl:type"
	NODE_NAMESPACE = "decl:namespace"

	NODE_TEMPLATE = "template"
	NODE_FNDEF    = "fndef"
	NODE_FNCALL   = "call"
	NODE_LOCAL    = "local"
	NODE_CONST    = "const"
	NODE_IMPORT   = "import"
	NODE_RETURN   = "return"

	// NODE_LIST denotes that the node's children are a list
	NODE_NAMESPACE_ACCESS = "::"
	NODE_COLON            = ":"
	NODE_BINOP_DOT        = "."
	NODE_AT               = "@"
	NODE_AS               = "as"
	NODE_ASSIGN           = "="
	NODE_IS               = "is"
	NODE_LT               = "<"
	NODE_GT               = ">"
	NODE_ELLIPSIS         = "..."
	NODE_PLUS             = "+"
	NODE_MIN              = "-"
	NODE_MUL              = "*"
	NODE_DIV              = "/"
	NODE_DEREF            = "deref"

	NODE_TYPED     = ":"
	NODE_INFER     = "infer"
	NODE_SIGNATURE = "signature"
	NODE_FN        = "fn"
	NODE_FOR       = "for"
	NODE_IF        = "if"

	NODE_FRAGMENT = "-fragment-" // this node should not appear in the final AST
	NODE_FILE     = "file"
	NODE_STR      = "str"
	NODE_LIST     = "lst"
	NODE_UNION    = "union"
	NODE_INDEX    = "index"
	NODE_BLOCK    = "block"
	NODE_DECLS    = "decls"
)

func NewTerminalNode(tk *Token) *Node {
	return &Node{
		Kind:     NODE_TERMINAL,
		Position: tk.Position,
		Token:    tk,
	}
}

func WrapNode(inKind NodeKind, node *Node) *Node {
	return &Node{
		Kind:     inKind,
		Position: node.Position,
		Children: []*Node{node},
	}
}

func NewNode(kind NodeKind, pos Position, chld ...*Node) *Node {
	res := &Node{
		Kind:     kind,
		Position: pos,
		Children: chld,
	}
	res.UpdatePosition()
	return res
}

func NewErrorNode(tk *Token, chld ...*Node) *Node {
	return &Node{
		Kind:     NODE_ERROR,
		Position: tk.Position,
		Token:    tk,
		Children: chld,
	}
}

type Node struct {
	// A node tagged as an error says that it is garbage, essentially.
	// The type checker will stop on nodes marked as errors, as they could crash
	// the compiler.
	IsError  bool
	Position Position
	Kind     NodeKind
	Token    *Token // only defined for literals
	Children []*Node
}

// func (n *Node) IsBinary()

func (n *Node) Left() *Node {
	return n.Children[0]
}

func (n *Node) Right() *Node {
	return n.Children[1]
}

func (n *Node) IsValidVariableName() bool {
	if n.Kind != NODE_TERMINAL || n.Token == nil || n.Token.Kind != TK_ID {
		return false
	}
	run, _ := utf8.DecodeRuneInString(n.Token.String())
	return unicode.IsLower(run)
}

func (n *Node) EnsureList() *Node {
	if !n.Is(NODE_LIST) {
		return WrapNode(NODE_LIST, n)
	}
	return n
}

func (n *Node) Is(nk NodeKind) bool {
	return nk == n.Kind
}

func (n *Node) ReportError(msg ...string) {
	n.IsError = true
	n.Position.Context.reportError(n.Position, msg...)
}

func (n *Node) IncludePosition(cpos *Position) {
	pos := &n.Position

	pos.Start = minInt(pos.Start, cpos.Start)
	pos.End = maxInt(pos.End, cpos.End)

	if cpos.Line < pos.Line {
		pos.Line = cpos.Line
		pos.Column = cpos.Column
	} else if cpos.Line == pos.Line {
		pos.Column = minInt(pos.Column, cpos.Column)
	}

	if cpos.EndLine > pos.EndLine {
		pos.EndLine = cpos.EndLine
		pos.EndColumn = cpos.EndColumn
	} else if cpos.EndLine == pos.EndLine {
		pos.EndColumn = maxInt(pos.EndColumn, cpos.EndColumn)
	}

}

func (n *Node) UpdatePosition() {
	for _, c := range n.Children {
		n.IncludePosition(&c.Position)
	}
}

func (n *Node) IsIdent() bool {
	if n.Kind != NODE_TERMINAL {
		return false
	}
	if n.Token != nil && n.Token.Kind == TK_ID {
		return true
	}
	return false
}

//////////////////// DUMP FUNCTIONS //////////////////

type ZoeWriter struct {
	io.Writer
}

func (z *ZoeWriter) Write(strs ...string) {
	for _, s := range strs {
		_, _ = z.Writer.Write([]byte(s))
	}
}

func (n *Node) debugChildren() string {
	if len(n.Children) == 0 {
		return ""
	}
	strs := make([]string, len(n.Children))
	for i, chld := range n.Children {
		strs[i] = chld.Debug()
	}
	return strings.Join(strs, " ")
}

func (n *Node) Debug() string {
	err := ""
	if n.IsError {
		err = red("!")
	}
	switch n.Kind {
	case NODE_ERROR:
		rest := ""
		if len(n.Children) > 0 {
			rest = " ..." + n.debugChildren()
		}
		return fmt.Sprintf("%s(%s %#v%s)", err, red(NODE_ERROR), n.Token.String(), rest)
	case NODE_TERMINAL:
		tkd := n.Token.Kind
		switch tkd {
		case TK_ID:
			return err + cyan(n.Token.String())
		case TK_RAWSTR:
			return err + green(n.Token.String())
		default:
			return err + yel(n.Token.String())
		}
	case NODE_BLOCK:
		return fmt.Sprint(err, mag("{"), n.debugChildren(), mag("}"))
	case NODE_LIST:
		return fmt.Sprint(err, mag("["), n.debugChildren(), mag("]"))
	default:
		return fmt.Sprint(err, grey("("), n.Kind, " ", n.debugChildren(), grey(")"))
	}
}

func (n *Node) String() string {
	c := n.Position.Context
	p := n.Position
	return string(c.data[p.Start:p.End])
}

// func (n *Node) WriteTree(w *ZoeWriter, indent string) {
// 	l := len(n.Children)
// 	// lst := "└ "
// 	switch n.Kind {
// 	case NODE_UNARY:
// 		w.Write(yel(`'`, n.Op().String(), `'`), "\n")
// 		w.Write(grey(indent + "╰ "))
// 		n.Subject().WriteTree(w, indent+"  ")
// 		return
// 	case NODE_BINOP:
// 		w.Write(grey("┌ "))
// 		n.Left().WriteTree(w, indent+"│ ")

// 		w.Write("\n", grey(indent), grey("╞ "), yel(`'`, n.Op().String(), `'`))
// 		w.Write("\n", grey(indent), grey("└ "))
// 		n.Right().WriteTree(w, indent+"  ")
// 		return
// 	case NODE_LIST:
// 		// lst = "╸ "
// 		// suffix := ""
// 		if l == 0 {
// 			w.Write(grey("(empty) "), mag(n.Token.String()))
// 		}

// 		pref := n.Op().String() + " "
// 		for i, c := range n.Children {
// 			indentmore := indent + "  "
// 			// pref := "╸ "

// 			if i > 0 || l == 0 {
// 				w.Write("\n", grey(indent))
// 			}
// 			w.Write(grey(pref))
// 			c.WriteTree(w, indentmore)
// 		}

// 	case NODE_LITERAL:
// 		w.Write(green(n.Token.String()))
// 	case NODE_ERROR:
// 		w.Write(red("! " + n.Token.String()))
// 	}

// }
