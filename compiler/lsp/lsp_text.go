package main

import (
	"bytes"
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

	log.Print("added file ", string(fsent.TextDocument.URI))
	f, err := req.Conn.Solution.AddFile(string(fsent.TextDocument.URI), fsent.TextDocument.Text, fsent.TextDocument.Version)
	if err == nil {
		diags := make([]lsp.Diagnostic, len(f.Errors))
		for i, e := range f.Errors {
			diags[i] = e.ToLspDiagnostic()
		}
		req.Notify("textDocument/publishDiagnostics", lsp.PublishDiagnosticsParams{
			URI:         fsent.TextDocument.URI,
			Diagnostics: diags,
		})
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

	file, ok := req.Conn.Solution.Files[string(changes.TextDocument.URI)]
	if !ok {
		log.Print("file not found ", string(changes.TextDocument.URI))
		// return errors.New(`file not found: ` + string(changes.TextDocument.URI))
		return nil
	}

	var data = file.GetData()
	// log.Print("original (", file.Version, "): ", string(data))
	for _, chg := range changes.ContentChanges {
		var buf bytes.Buffer
		var offsetstart = file.GetOffsetForPosition(chg.Range.Start)
		var offsetend = file.GetOffsetForPosition(chg.Range.End)
		// log.Print(`REPLACING `, offsetstart, `-`, offsetend)
		// log.Print(`In tokens[`, len(file.Tokens), `] range is `, offsetstart, `-`, offsetend, ` for len `, len(data)-1)
		_, _ = buf.Write(data[0:offsetstart])
		if offsetstart <= offsetend {
			_, _ = buf.Write([]byte(chg.Text))
		}
		_, _ = buf.Write(data[offsetend:])
		data = buf.Bytes()
	}

	// req.Conn.Solution.AddFile(string(changes.TextDocument.URI), changes.ContentChanges[0].Text, changes.TextDocument.Version)

	// log.Print(data)
	f, err := req.Conn.Solution.AddFile(string(changes.TextDocument.URI), string(data), changes.TextDocument.Version)
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
