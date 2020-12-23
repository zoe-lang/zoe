package main

import (
	"encoding/json"
	"fmt"

	zoe "github.com/ceymard/zoe/compiler"
	"github.com/sourcegraph/go-lsp"
)

func init() {
	handlers["textDocument/hover"] = HandleHover
}

func HandleHover(req *LspRequest) error {
	params := lsp.TextDocumentPositionParams{}
	if err := json.Unmarshal(req.RawParams(), &params); err != nil {
		return err
	}

	fname := zoe.InternedIds.Save(string(params.TextDocument.URI))
	file, ok := req.Conn.Solution.Files[fname]
	if !ok {
		req.Reply(lsp.Hover{
			Contents: []lsp.MarkedString{{
				Language: "zoe",
				Value:    fmt.Sprint("file '", fname, "'not found"),
			}},
		})
		return nil
	}

	pos, err2 := file.FindNodePosition(params.Position)
	if err2 != nil {
		req.Reply(lsp.Hover{
			Contents: []lsp.MarkedString{{
				Language: "zoe",
				Value:    err2.Error(),
			}},
		})
		return err2
	}

	res := lsp.Hover{}

	var last = pos[len(pos)-1]
	dbg := last.Debug()
	// log.Print(pos)
	res.Contents = []lsp.MarkedString{{
		Language: "zoe",
		Value:    dbg + " " + last.GetText(),
	}}

	req.Reply(res)
	return nil
}
