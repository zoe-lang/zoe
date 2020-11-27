package main

import (
	"encoding/json"

	"github.com/k0kubun/pp"
	"github.com/sourcegraph/go-lsp"
)

func init() {
	handlers["textDocument/didOpen"] = HandleDidOpen
	handlers["textDocument/didChange"] = HandleDidChange
}

func HandleDidOpen(req *LspRequest) error {
	fsent := lsp.DidOpenTextDocumentParams{}
	if err := json.Unmarshal(req.RawParams(), &fsent); err != nil {
		return err
	}
	pp.Print(fsent)
	return nil
}

func HandleDidChange(req *LspRequest) error {

	// try to parse the changes

	changes := lsp.DidChangeTextDocumentParams{}
	// log.Print(string(req.RawParams()))
	if err := json.Unmarshal(req.RawParams(), &changes); err != nil {
		return err
	}
	// pp.Print(changes)
	return nil
}
