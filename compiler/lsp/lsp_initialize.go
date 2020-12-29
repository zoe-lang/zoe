package main

import (
	"github.com/sourcegraph/go-lsp"
)

func init() {
	handlers["initialize"] = HandleInitialize
	handlers["initialized"] = HandleInitialized
}

var Capabilities = lsp.ServerCapabilities{}

func HandleInitialize(req *LspRequest) error {

	// log.Print(req.Params.String())

	req.Reply(lsp.InitializeResult{
		Capabilities: Capabilities,
	})

	return nil
}

// This is called by vscode after a response to initialize was sent.
// For now we're not handling it, we just assume everything is right
// and that we can simply reply to stuff.
func HandleInitialized(_ *LspRequest) error {
	// for now we do nothing here
	return nil
}
