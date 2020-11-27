package main

import "github.com/sourcegraph/go-lsp"

func init() {
	handlers["textDocument/completion"] = HandleCompletion
}

func HandleCompletion(req *LspRequest) error {
	req.Reply([]lsp.CompletionItem{
		{
			Label:  "Dodododo",
			Kind:   lsp.CIKProperty,
			Detail: "Ooooh yeeah",
		},
	})

	return nil
}
