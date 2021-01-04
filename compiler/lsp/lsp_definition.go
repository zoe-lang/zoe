package main

import (
	"context"
	"errors"
	"log"

	zoe "github.com/ceymard/zoe/compiler"
	"github.com/creachadair/jrpc2/handler"
	"github.com/sourcegraph/go-lsp"
)

func init() {
	addHandler(func(l *LspConnection, mp handler.Map) {
		mp["textDocument/definition"] = handler.New(l.HandleDefinition)

	})
	Capabilities.DefinitionProvider = true
	// Capabilities.TypeDefinitionProvider = true
}

// We should handle that as well.
// func HandleTypeDefinition(req *LspRequest) error {

// }

func (l *LspConnection) HandleDefinition(_ context.Context, params lsp.TextDocumentPositionParams) (*lsp.Location, error) {

	// fname := zoe.InternedIds.Save(string(params.TextDocument.URI))
	var file, ok = l.Solution.Files[string(params.TextDocument.URI)]
	if !ok {
		return nil, errors.New(`file not found`)
	}

	var pos = params.Position
	path, err := file.FindNodePosition(pos)

	// Now that we've got the path, we can try and find out if its a symbol, and whether
	// we should find a definition for it.

	if err != nil || len(path) == 0 {
		return nil, err
	}
	var last = path[len(path)-1]

	// TODO, we should check along the path to look for the first
	// instance of a "symbol" tree to try and resolve the symbol.
	if !last.Is(zoe.NODE_ID) {
		return nil, nil
	}

	if found, ok := last.FindDefinition(); ok {
		var loc = lsp.Location{}
		loc.URI = lsp.DocumentURI(found.File().Filename)
		loc.Range = found.Range()
		return &loc, nil
	}

	log.Print(path)

	return nil, nil
}
