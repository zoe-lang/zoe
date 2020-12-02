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

func (b *nodeBuilder) createNodeFromToken(tk TokenPos, nk AstNodeKind, scope Scope, children ...NodePosition) NodePosition {
	res := b.createNode(b.tokens[tk].Range, nk, scope)
	node := &b.file.Nodes[res]
	l := len(children)
	if l > 0 {
		node.ArgLen = int8(l)
		for i, chld := range children {
			node.Args[i] = chld
			b.extendsNodeRangeFromNode(res, chld)
		}
	}
	return res
}

func (b *nodeBuilder) createNodeFromCurrentToken(nk AstNodeKind, scope Scope) NodePosition {
	return b.createNodeFromToken(b.current, nk, scope)
}

// createNode creates a node in the underlying node slice
// it should not be called directly. Rather, the different
// create...() provide a safer way to create the ast as it should
// provide checks and balances.
func (b *nodeBuilder) createNode(rng Range, kind AstNodeKind, scope Scope) NodePosition {
	// maybe we should handle here the capacity of the node arrays ?
	l := NodePosition(len(b.file.Nodes))
	b.file.Nodes = append(b.file.Nodes, AstNode{Kind: kind, Range: rng, Scope: scope.pos})
	return l
}

// Doesn't create an empty node
func (b *nodeBuilder) createIfToken(tk TokenKind, fn func(tk TokenPos) NodePosition) NodePosition {
	if b.currentTokenIs(tk) {
		cur := b.current
		b.advance()
		return fn(cur)
	}
	return 0
}

// Doesn't create an empty node
func (b *nodeBuilder) createExpectToken(tk TokenKind, fn func(tk TokenPos) NodePosition) NodePosition {
	if b.currentTokenIs(tk) {
		cur := b.current
		b.advance()
		return fn(cur)
	}
	b.reportErrorAtToken(b.current, "expected "+tokstr[tk]+" but got '"+b.getTokenText(b.current)+"'")
	return 0
}

func (b *nodeBuilder) createAndExpectOrEmpty(tk TokenKind, fn func(tk TokenPos) NodePosition) NodePosition {
	res := b.createIfToken(tk, fn)
	if res == 0 {
		b.reportErrorAtToken(b.current, "expected '", tokstr[tk], "'")
		return EmptyNode
	}
	return res
}

func (b *nodeBuilder) createIfTokenOrEmpty(tk TokenKind, fn func(tk TokenPos) NodePosition) NodePosition {
	res := b.createIfToken(tk, fn)
	if res == 0 {
		return EmptyNode
	}
	return res
}

func (b *nodeBuilder) nodeIs(ni NodePosition, nk AstNodeKind) bool {
	if ni == 0 {
		return false
	}
	return nk == b.file.Nodes[ni].Kind
}

func (b *nodeBuilder) extendNodeRange(ni NodePosition, rng Range) {
	b.file.Nodes[ni].Range.Extend(rng)
}

func (b *nodeBuilder) extendsNodeRangeFromNode(ni NodePosition, other NodePosition) {
	for other != EmptyNode {
		o := &b.file.Nodes[other]
		b.extendNodeRange(ni, o.Range)

		if o.Next == EmptyNode {
			break
		}

		other = o.Next
	}
}

func (b *nodeBuilder) extendRangeFromToken(ni NodePosition, tk TokenPos) {
	b.extendNodeRange(ni, b.tokens[tk].Range)
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// func (b *nodeBuilder) setNodeChildren(node NodePosition, children ...NodePosition) {
// 	app := b.appender(node)
// 	for _, c := range children {
// 		app.append(c)
// 	}
// }

func (b *nodeBuilder) createIdNode(tk TokenPos, scope Scope) NodePosition {
	idstr := InternedIds.Save(b.getTokenText(tk))
	idnode := b.createNodeFromToken(tk, NODE_ID, scope)
	b.file.Nodes[idnode].Value = idstr
	return idnode
}

func (b *nodeBuilder) fragment() *fragment {
	return &fragment{builder: b}
}

func (b *nodeBuilder) getTokenText(tk TokenPos) string {
	return b.file.GetTokenText(tk)
}

func (b *nodeBuilder) createBinOp(tk TokenPos, kind AstNodeKind, scope Scope, left NodePosition, right NodePosition) NodePosition {
	res := b.createNodeFromToken(tk, kind, scope, left, right)
	return res
}

func (b *nodeBuilder) createUnaryOp(tk TokenPos, kind AstNodeKind, scope Scope, left NodePosition) NodePosition {
	res := b.createNodeFromToken(tk, kind, scope, left)
	return res
}

// cloneNode shallow clones a node, mostly to help it have a different next,
// as "children" of nodes are not slices but just positions inside
// the node array.
func (b *nodeBuilder) cloneNode(pos NodePosition) NodePosition {
	var n = b.file.Nodes[pos]
	l := len(b.file.Nodes)
	n.Next = 0
	b.file.Nodes = append(b.file.Nodes, n)
	return NodePosition(l)
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// A list of nodes
type fragment struct {
	builder *nodeBuilder
	first   NodePosition
	last    NodePosition
}

func (f *fragment) append(pos NodePosition) {
	if f.first == EmptyNode {
		f.first = pos
		f.last = pos
		return
	}

	nodes := f.builder.file.Nodes
	target := &nodes[f.last]

	target.Next = pos

	for nodes[pos].Next != EmptyNode {
		pos = nodes[pos].Next
	}

	f.last = pos
}
