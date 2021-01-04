package main

import (
	"bytes"
	"context"
	"log"

	"github.com/creachadair/jrpc2/handler"
	"github.com/sourcegraph/go-lsp"
)

func init() {
	addHandler(func(l *LspConnection, mp handler.Map) {
		mp["textDocument/didOpen"] = handler.New(l.HandleDidOpen)
		mp["textDocument/didChange"] = handler.New(l.HandleDidChange)
	})

	Capabilities.TextDocumentSync = &lsp.TextDocumentSyncOptionsOrKind{
		Options: &lsp.TextDocumentSyncOptions{
			OpenClose: true,
			Change:    lsp.TDSKIncremental, // TODO: change that to incremental to avoid resending the whole file all the time.
		},
	}

}

func (l *LspConnection) HandleDidOpen(ctx context.Context, fsent lsp.DidOpenTextDocumentParams) error {

	log.Print("added file ", string(fsent.TextDocument.URI))
	f, err := l.Solution.AddFile(string(fsent.TextDocument.URI), fsent.TextDocument.Text, fsent.TextDocument.Version)
	if err == nil {
		diags := make([]lsp.Diagnostic, len(f.Errors))
		for i, e := range f.Errors {
			diags[i] = e.ToLspDiagnostic()
		}

		l.Server.Notify(ctx, "textDocument/publishDiagnostics", lsp.PublishDiagnosticsParams{
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
func (l *LspConnection) HandleDidChange(ctx context.Context, changes lsp.DidChangeTextDocumentParams) error {

	// try to parse the changes
	file, ok := l.Solution.Files[string(changes.TextDocument.URI)]
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
		log.Print(`In tokens[`, len(file.Tokens), `] range is `, offsetstart, `-`, offsetend, ` for len `, len(data)-1)
		_, _ = buf.Write(data[0:offsetstart])
		if offsetstart <= offsetend {
			_, _ = buf.Write([]byte(chg.Text))
		}
		_, _ = buf.Write(data[offsetend:])
		data = buf.Bytes()
	}

	// req.Conn.Solution.AddFile(string(changes.TextDocument.URI), changes.ContentChanges[0].Text, changes.TextDocument.Version)

	log.Print(string(data))
	f, err := l.Solution.AddFile(string(changes.TextDocument.URI), string(data), changes.TextDocument.Version)
	if err == nil {
		// if len(f.Errors) > 0 {
		diags := make([]lsp.Diagnostic, len(f.Errors))
		for i, e := range f.Errors {
			diags[i] = e.ToLspDiagnostic()
		}
		l.Server.Notify(ctx, "textDocument/publishDiagnostics", lsp.PublishDiagnosticsParams{
			URI:         changes.TextDocument.URI,
			Diagnostics: diags,
		})
		// }
		// pp.Print(f.Errors)
	}

	// pp.Print(changes)
	return nil
}
