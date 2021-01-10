package main

import (
	"context"

	"github.com/creachadair/jrpc2/handler"
	"github.com/sourcegraph/go-lsp"
)

func init() {
	// handlers["textDocument/hover"] = HandleHover
	addHandler(func(l *LspConnection, mp handler.Map) {
		mp["textDocument/hover"] = handler.New(l.HandleHover)
	})
	Capabilities.HoverProvider = true
}

func (l *LspConnection) HandleHover(_ context.Context, params lsp.TextDocumentPositionParams) (*lsp.Hover, error) {
	return nil, nil
	// var fname = string(params.TextDocument.URI)
	// file, ok := l.Solution.Files[fname]
	// if !ok {
	// 	return &lsp.Hover{
	// 		Contents: []lsp.MarkedString{{
	// 			Language: "zoe",
	// 			Value:    fmt.Sprint("file '", fname, "'not found"),
	// 		}},
	// 	}, nil
	// }

	// pos, err2 := file.FindNodePosition(params.Position)
	// if err2 != nil {
	// 	return &lsp.Hover{
	// 		Contents: []lsp.MarkedString{{
	// 			Language: "zoe",
	// 			Value:    err2.Error(),
	// 		}},
	// 	}, err2
	// }

	// res := lsp.Hover{}

	// var last = pos[len(pos)-1]
	// dbg := last.Debug()
	// // log.Print(pos)
	// res.Contents = []lsp.MarkedString{{
	// 	Language: "zoe",
	// 	Value:    dbg + " " + last.GetText(),
	// }}

	// return &res, nil
}
