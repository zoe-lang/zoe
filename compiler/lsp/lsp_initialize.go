package main

import (
	"context"
	"log"

	"github.com/creachadair/jrpc2/handler"
	"github.com/sourcegraph/go-lsp"
)

func init() {
	addHandler(func(l *LspConnection, mp handler.Map) {
		mp["initialize"] = handler.New(l.HandleInitialize)
		mp["initialized"] = handler.New(l.HandleInitialized)
		mp["$cancelRequest"] = handler.New(l.cancelRequestHandler)
	})

	// handlers["initialize"] = HandleInitialize
	// handlers["initialized"] = HandleInitialized
}

func (l *LspConnection) cancelRequestHandler(_ context.Context, _ lsp.None) error {
	log.Printf("cancel request")
	return nil
}

var Capabilities = lsp.ServerCapabilities{}

func (l *LspConnection) HandleInitialize(_ context.Context, _ lsp.InitializeParams) lsp.InitializeResult {

	// log.Print(req.Params.String())

	return lsp.InitializeResult{
		Capabilities: Capabilities,
	}

}

// This is called by vscode after a response to initialize was sent.
// For now we're not handling it, we just assume everything is right
// and that we can simply reply to stuff.
func (l *LspConnection) HandleInitialized(_ context.Context, _ lsp.None) error {

	// for now we do nothing here
	return nil
}
