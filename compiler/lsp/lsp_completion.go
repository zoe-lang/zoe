package main

func init() {
	handlers["textDocument/Completion"] = HandleCompletion
}

func HandleCompletion(req *LspRequest) error {

	return nil
}
