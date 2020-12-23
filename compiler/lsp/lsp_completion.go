package main

import (
	"encoding/json"
	"errors"
	"log"

	zoe "github.com/ceymard/zoe/compiler"
	"github.com/sourcegraph/go-lsp"
)

func init() {
	handlers["textDocument/completion"] = HandleCompletion
}

func HandleCompletion(req *LspRequest) error {

	var params = lsp.CompletionParams{}
	if err := json.Unmarshal(req.RawParams(), &params); err != nil {
		return err
	}

	fname := zoe.InternedIds.Save(string(params.TextDocument.URI))
	var file, ok = req.Conn.Solution.Files[fname]
	if !ok {
		return errors.New(`file not found`)
	}

	var pos = params.Position
	var path, err = file.FindNodePosition(pos)
	log.Print(path)

	if err != nil {
		return err
	}

	var result = make([]lsp.CompletionItem, 0)

	if len(path) > 0 {
		var last = path[len(path)-1]
		for _, name := range last.Scope().AllNames() {
			result = append(result, lsp.CompletionItem{
				Label: name.Name,
				Kind:  lsp.CIKProperty,
			})
		}
	}
	req.Reply(result)
	// log.Print(last.Scope().Find())

	// req.Reply([]lsp.CompletionItem{
	// 	{
	// 		Label:  path[len(path)-1].Debug(),
	// 		Kind:   lsp.CIKProperty,
	// 		Detail: "Ooooh yeeah",
	// 	},
	// })

	return nil
}
