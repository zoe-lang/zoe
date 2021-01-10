package zoe

/*
	Register a node to the file.
	Here, the file checks whether the node being registered is allowed in this context.
*/
func (astf *AstFile) Register(n Node, scope *Scope) {
	// log.Printf("%#v", n)
	if n == nil {
		// Error case that may already have been handled
		return
	}

	var name = n.GetName()
	if name == nil {
		n.ReportError("this statement is disallowed in a file context")
		return
	}

	if v, ok := n.(*AstVarDecl); ok {
		if !v.IsConst {
			n.ReportError("only constants are permitted outside function bodies")
		}
	}

	scope.Add(n)
	// FIXME add to the namespace
}

/*
	Register a node onto a type
*/
func (tp *AstTypeDecl) Register(n Node, _ *Scope) {
	var name = n.GetName()
	if name == nil {
		n.ReportError("expected a declaration")
	}
	tp.Members[name.Name] = n
}

/*

 */
func (st *AstStructDecl) Register(n Node, scope *Scope) {
	st.nodeBase.Register(n, scope)
	st.AstTypeDecl.Register(n, scope)
}

func (en *AstEnumDecl) Register(n Node, scope *Scope) {
	en.nodeBase.Register(n, scope)
	en.AstTypeDecl.Register(n, scope)
}
