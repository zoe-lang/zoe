package main

import (
	"context"

	zoe "github.com/ceymard/zoe/compiler"
	"github.com/creachadair/jrpc2/handler"
	"github.com/sourcegraph/go-lsp"
)

func init() {

	addHandler(func(l *LspConnection, mp handler.Map) {
		mp["textDocument/completion"] = handler.New(l.HandleCompletion)
	})

	Capabilities.CompletionProvider = &lsp.CompletionOptions{
		TriggerCharacters: []string{".", "::"},
	}

}

func NodeCIKKind(_ zoe.Node) lsp.CompletionItemKind {
	return lsp.CIKProperty
}

func (l *LspConnection) HandleCompletion(_ context.Context, params lsp.CompletionParams) ([]lsp.CompletionItem, error) {
	return nil, nil

	// fname := zoe.InternedIds.Save(string(params.TextDocument.URI))
	// var file, ok = l.Solution.Files[string(params.TextDocument.URI)]
	// if !ok {
	// 	return nil, errors.New(`file not found`)
	// }

	// var pos = params.Position
	// var path, err = file.FindNodePosition(pos)
	// log.Print(path)

	// if err != nil {
	// 	return nil, err
	// }

	// var result = make([]lsp.CompletionItem, 0)

	// if len(path) > 0 {
	// 	var last = path[len(path)-1]
	// 	log.Print(last.Debug())
	// 	for _, name := range last.Scope().AllNames() {
	// 		var docstr = ""
	// 		if doc, ok := name.Node.DocComment(); ok {
	// 			log.Print("found the ocmment")
	// 			docstr = doc.GetText()
	// 		}

	// 		result = append(result, lsp.CompletionItem{
	// 			Label:         name.Name,
	// 			Kind:          NodeCIKKind(name.Node),
	// 			Documentation: docstr,
	// 			Detail:        docstr,
	// 		})
	// 	}
	// }

	// return result, nil
}
