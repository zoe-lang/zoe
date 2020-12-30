package main

import (
	"encoding/json"
	"errors"
	"log"

	zoe "github.com/ceymard/zoe/compiler"
	"github.com/sourcegraph/go-lsp"
)

func init() {
	handlers["textDocument/definition"] = HandleDefinition
	Capabilities.DefinitionProvider = true
}

func HandleDefinition(req *LspRequest) error {

	var params lsp.TextDocumentPositionParams
	if err := json.Unmarshal(req.RawParams(), &params); err != nil {
		return err
	}

	// fname := zoe.InternedIds.Save(string(params.TextDocument.URI))
	var file, ok = req.Conn.Solution.Files[string(params.TextDocument.URI)]
	if !ok {
		return errors.New(`file not found`)
	}

	var pos = params.Position
	path, err := file.FindNodePosition(pos)

	// Now that we've got the path, we can try and find out if its a symbol, and whether
	// we should find a definition for it.

	if err != nil || len(path) == 0 {
		return err
	}
	var last = path[len(path)-1]
	if !last.Is(zoe.NODE_ID) {
		req.Reply(nil)
		return nil
	}

	if found, ok := last.FindDefinition(); ok {
		var loc = lsp.Location{}
		loc.URI = lsp.DocumentURI(found.File().Filename)
		loc.Range = found.Range()
		req.Reply(loc)
		return nil
	}

	log.Print(path)

	req.Reply(nil)
	return nil
}
