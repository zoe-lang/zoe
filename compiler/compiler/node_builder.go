package zoe

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// A list of nodes
type list struct {
	first Node
	last  Node
}

func newList() list {
	return list{}
}

func (f *list) append(node Node) {

	if node.IsEmpty() {
		return
	}

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
