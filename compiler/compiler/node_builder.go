package zoe

const ErrorNode NodePosition = 0
const EmptyNode NodePosition = 1

type nodeBuilder struct {
	file          *File
	nodes         NodeArray
	tokens        []Token
	current       TokenPos
	tokensLen     TokenPos
	doccommentMap map[NodePosition]TokenPos
}

func (b *nodeBuilder) reportError(rng Range, msg ...string) {
	b.file.reportError(rng, msg...)
}

func (b *nodeBuilder) reportErrorAtToken(tk TokenPos, msg ...string) {
	rng := b.tokens[tk]
	b.file.reportError(rng, msg...)
}

func (b *nodeBuilder) createEmptyNode() NodePosition {
	return b.createNode(Range{}, NODE_EMPTY)
}

func (b *nodeBuilder) createNodeFromToken(tk TokenPos, nk AstNodeKind, children ...NodePosition) NodePosition {
	res := b.createNode(b.tokens[tk].Range, nk)
	if len(children) > 0 {
		b.setNodeChildren(res, children...)
	}
	return res
}

func (b *nodeBuilder) createNodeFromCurrentToken(nk AstNodeKind) NodePosition {
	return b.createNodeFromToken(b.current, nk)
}

// createNode creates a node in the underlying node slice
// it should not be called directly. Rather, the different
// create...() provide a safer way to create the ast as it should
// provide checks and balances.
func (b *nodeBuilder) createNode(rng Range, kind AstNodeKind) NodePosition {
	// maybe we should handle here the capacity of the node arrays ?
	l := NodePosition(len(b.nodes))
	b.nodes = append(b.nodes, AstNode{Kind: kind, Range: rng})
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
		return b.createEmptyNode()
	}
	return res
}

func (b *nodeBuilder) createIfTokenOrEmpty(tk TokenKind, fn func(tk TokenPos) NodePosition) NodePosition {
	res := b.createIfToken(tk, fn)
	if res == 0 {
		return b.createEmptyNode()
	}
	return res
}

func (b *nodeBuilder) nodeIs(ni NodePosition, nk AstNodeKind) bool {
	if ni == 0 {
		return false
	}
	return nk == b.nodes[ni].Kind
}

func (b *nodeBuilder) extendNodeRange(ni NodePosition, rng Range) {
	b.nodes[ni].Range.Extend(rng)
}

func (b *nodeBuilder) extendRangeFromToken(ni NodePosition, tk TokenPos) {
	b.extendNodeRange(ni, b.tokens[tk].Range)
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// repeat a block until cbk returns 0 (empty node)
func (b *nodeBuilder) repeat(rng Range, nk AstNodeKind, cbk func() NodePosition) NodePosition {
	res := b.createNode(rng, nk)
	app := b.appender(res)
	for true {
		pos := cbk()
		if pos == 0 {
			break
		}
		app.append(pos)
	}
	return res
}

func (b *nodeBuilder) setNodeChildren(node NodePosition, children ...NodePosition) {
	app := b.appender(node)
	for _, c := range children {
		app.append(c)
	}
}

func (b *nodeBuilder) createIdNode(tk TokenPos) NodePosition {
	idstr := internedIds.Save(b.getTokenText(tk))
	idnode := b.createNodeFromToken(tk, NODE_ID)
	b.nodes[idnode].Value = idstr
	return idnode
}

func (b *nodeBuilder) appender(from NodePosition) *appender {
	return &appender{builder: b, first: from, pos: from}
}

func (b *nodeBuilder) fragmenter() *fragment {
	return &fragment{builder: b}
}

func (b *nodeBuilder) getTokenText(tk TokenPos) string {
	return b.file.GetTokenText(tk)
}

// cloneNode shallow clones a node, mostly to help it have a different next
func (b *nodeBuilder) cloneNode(pos NodePosition) NodePosition {
	var n = b.nodes[pos]
	l := len(b.nodes)
	n.Next = 0
	b.nodes = append(b.nodes, n)
	return NodePosition(l)
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// appender is a type used to update list of nodes
type appender struct {
	builder *nodeBuilder
	first   NodePosition
	pos     NodePosition
}

func (a *appender) append(pos NodePosition) {
	first := a.first
	nodes := a.builder.nodes
	target := &nodes[a.pos]
	// FIXME check if pos already has a Next (this shouldn't be the case, unless
	// we have a fragment).
	if first == a.pos {
		target.Value = int(pos) // the start position
	} else {
		target.Next = pos
	}
	for nodes[pos].Next != 0 {
		pos = nodes[pos].Next
	}
	a.pos = pos
	nodes[first].Range.Extend(nodes[pos].Range)
}

//// More or less the same as appender

type fragment struct {
	builder *nodeBuilder
	first   NodePosition
	last    NodePosition
}

func (f *fragment) append(pos NodePosition) {
	if f.first == 0 {
		f.first = pos
		f.last = pos
		return
	}

	nodes := f.builder.nodes
	target := &nodes[f.last]

	target.Next = pos

	for nodes[pos].Next != 0 {
		pos = nodes[pos].Next
	}

	f.last = pos
}
