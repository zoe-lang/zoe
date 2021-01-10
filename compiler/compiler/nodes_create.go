package zoe

func AstFileNew(file *File) *AstFile {
	return &AstFile{
		File:     file,
		nodeBase: nodeBase{TkRange: newTkRange(), File: file},
	}
}

func (parser *Parser) createNodeBase() nodeBase {
	return nodeBase{TkRange: parser.AsRange(), File: parser.file}
}

func (sc *Scope) SetParent(parent *Scope) {
	sc.Parent = parent
}
