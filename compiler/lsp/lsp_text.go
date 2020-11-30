package main

import (
	"encoding/json"

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

	f, err := req.Conn.Solution.AddFile(string(fsent.TextDocument.URI), fsent.TextDocument.Text)
	if err == nil {
		if len(f.Errors) > 0 {
			diags := make([]lsp.Diagnostic, len(f.Errors))
			for i, e := range f.Errors {
				diags[i].Message = e.Message
				diags[i].Range = lsp.Range{
					Start: lsp.Position{Line: int(e.Range.Line - 1), Character: int(e.Range.Column - 1)},
					End:   lsp.Position{Line: int(e.Range.Line - 1), Character: int(e.Range.Column)},
				}
			}
			req.Notify("textDocument/publishDiagnostics", lsp.PublishDiagnosticsParams{
				URI:         fsent.TextDocument.URI,
				Diagnostics: diags,
			})
		}
		// pp.Print(f.Errors)
	}
	// pp.Print(fsent)
	return nil
}

func HandleDidChange(req *LspRequest) error {

	// try to parse the changes

	changes := lsp.DidChangeTextDocumentParams{}
	// log.Print(string(req.RawParams()))
	if err := json.Unmarshal(req.RawParams(), &changes); err != nil {
		return err
	}
	req.Conn.Solution.AddFile(string(changes.TextDocument.URI), changes.ContentChanges[0].Text)
	// pp.Print(changes)
	return nil
}
