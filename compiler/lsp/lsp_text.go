package main

import (
	"encoding/json"
	"log"

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

	f, err := req.Conn.Solution.AddFile(string(fsent.TextDocument.URI), fsent.TextDocument.Text, fsent.TextDocument.Version)
	if err == nil {
		if len(f.Errors) > 0 {
			diags := make([]lsp.Diagnostic, len(f.Errors))
			for i, e := range f.Errors {
				diags[i] = e.ToLspDiagnostic()
			}
			req.Notify("textDocument/publishDiagnostics", lsp.PublishDiagnosticsParams{
				URI:         fsent.TextDocument.URI,
				Diagnostics: diags,
			})
		}
		// pp.Print(f.Errors)
	} else {
		log.Print(err)
	}
	// pp.Print(fsent)
	return nil
}

// HandleDidChange should be throttled and operate on incremental updates instead of
// full text synchronization.
func HandleDidChange(req *LspRequest) error {

	// try to parse the changes

	changes := lsp.DidChangeTextDocumentParams{}
	// log.Print(string(req.RawParams()))
	if err := json.Unmarshal(req.RawParams(), &changes); err != nil {
		return err
	}
	req.Conn.Solution.AddFile(string(changes.TextDocument.URI), changes.ContentChanges[0].Text, changes.TextDocument.Version)

	f, err := req.Conn.Solution.AddFile(string(changes.TextDocument.URI), changes.ContentChanges[0].Text, changes.TextDocument.Version)
	if err == nil {
		// if len(f.Errors) > 0 {
		diags := make([]lsp.Diagnostic, len(f.Errors))
		for i, e := range f.Errors {
			diags[i] = e.ToLspDiagnostic()
		}
		req.Notify("textDocument/publishDiagnostics", lsp.PublishDiagnosticsParams{
			URI:         changes.TextDocument.URI,
			Diagnostics: diags,
		})
		// }
		// pp.Print(f.Errors)
	}

	// pp.Print(changes)
	return nil
}
