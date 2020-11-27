package main

import (
	"log"

	"github.com/sourcegraph/go-lsp"
)

func init() {
	handlers["initialize"] = HandleInitialize
	handlers["initialized"] = HandleInitialized
}

func HandleInitialize(req *LspRequest) error {

	log.Print(req.Params.String())

	req.Reply(lsp.InitializeResult{
		Capabilities: lsp.ServerCapabilities{
			TextDocumentSync: &lsp.TextDocumentSyncOptionsOrKind{
				Options: &lsp.TextDocumentSyncOptions{
					OpenClose: true,
					Change:    lsp.TDSKIncremental,
				},
			},
			CompletionProvider: &lsp.CompletionOptions{
				TriggerCharacters: []string{"."},
				ResolveProvider:   true,
			},
		},
	})

	return nil
}

func HandleInitialized(req *LspRequest) error {
	return nil
}
