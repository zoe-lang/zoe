package zoe

type nodeBuilder struct {
	file          *File
	tokens        []Token
	current       TokenPos
	tokensLen     TokenPos
	doccommentMap map[NodePosition]TokenPos
}

func (b *nodeBuilder) reportError(rng Range, msg ...string) {
	b.file.reportError(rng, msg...)
}

func (b *nodeBuilder) reportErrorAtPosition(pos NodePosition, msg ...string) {
	b.file.reportError(b.file.Nodes[pos].Range, msg...)
}

func (b *nodeBuilder) reportErrorAtToken(tk TokenPos, msg ...string) {
	if tk >= b.tokensLen {
		tk = b.tokensLen - 1
	}
	rng := b.tokens[tk]
	b.file.reportError(rng, msg...)
}

// createNode creates a node in the underlying node slice
// it should not be called directly. Rather, the different
// create...() provide a safer way to create the ast as it should
// provide checks and balances.

// Doesn't create an empty node

// func (b *nodeBuilder) extendNodeRange(ni NodePosition, rng Range) {
// 	b.file.Nodes[ni].Range.Extend(rng)
// }

// func (b *nodeBuilder) extendsNodeRangeFromNode(ni NodePosition, other NodePosition) {
// 	for other != EmptyNode {
// 		o := &b.file.Nodes[other]
// 		b.extendNodeRange(ni, o.Range)

// 		if o.Next == EmptyNode {
// 			break
// 		}

// 		other = o.Next
// 	}
// }

// func (b *nodeBuilder) extendRangeFromToken(ni NodePosition, tk TokenPos) {
// 	b.extendNodeRange(ni, b.tokens[tk].Range)
// }

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// func (b *nodeBuilder) setNodeChildren(node NodePosition, children ...NodePosition) {
// 	app := b.appender(node)
// 	for _, c := range children {
// 		app.append(c)
// 	}
// }

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// A list of nodes
type fragment struct {
	first Node
	last  Node
}

func newFragment() fragment {
	return fragment{}
}

func (f *fragment) append(node Node) {
	if f.first.IsEmpty() {
		f.first = node
		f.last = node
		return
	}

	f.last.SetNext(node)
	for node.HasNext() {
		node = node.Next()
	}

	f.last = node
}
