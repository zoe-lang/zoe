package main

func init() {
	handlers["textDocument/definition"] = HandleDefinition
}

func HandleDefinition(req *LspRequest) error {

	return nil
}
