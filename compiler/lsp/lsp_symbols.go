package main

import (
	"context"

	zoe "github.com/ceymard/zoe/compiler"
	"github.com/creachadair/jrpc2/handler"
	"github.com/sourcegraph/go-lsp"
)

type SymbolTags struct {
}

// Not in the go-lsp package
type DocumentSymbol struct {
	Name           string           `json:"name"`
	Detail         string           `json:"detail,omitempty"`
	Kind           lsp.SymbolKind   `json:"kind"`
	Tags           []SymbolTags     `json:"tags,omitempty"`
	Deprecated     bool             `json:"deprecated,omitempty"`
	Range          lsp.Range        `json:"range"`
	SelectionRange lsp.Range        `json:"selectionRange"`
	Children       []DocumentSymbol `json:"children,omitempty"`
}

func init() {
	addHandler(func(l *LspConnection, mp handler.Map) {
		mp["textDocument/documentSymbol"] = handler.New(l.HandleFileSymbols)
	})
	Capabilities.DocumentSymbolProvider = true
}

func NodeSymbolKind(n zoe.Node) lsp.SymbolKind {
	// switch n.Kind() {
	// case zoe.NODE_VAR:
	// 	return lsp.SKVariable
	// case zoe.NODE_ENUM:
	// 	return lsp.SKEnum
	// case zoe.NODE_STRUCT:
	// 	return lsp.SKStruct
	// case zoe.NODE_TYPE:
	// 	return lsp.SKClass
	// case zoe.NODE_TRAIT:
	// 	return lsp.SKStruct
	// case zoe.NODE_FN:
	// 	return lsp.SKFunction
	// case zoe.NODE_IMPORT:
	// 	return lsp.SKModule
	// default:
	return lsp.SKProperty
	// }
}

func (l *LspConnection) HandleFileSymbols(_ context.Context, params lsp.DocumentSymbolParams) ([]DocumentSymbol, error) {

	// fname := zoe.InternedIds.Save(string(params.TextDocument.URI))
	// var file, ok = l.Solution.Files[string(params.TextDocument.URI)]
	// if !ok {
	// 	return nil, errors.New(`file not found`)
	// }

	// var process_scope func(zoe.Scope) []DocumentSymbol

	// process_scope = func(scope zoe.Scope) []DocumentSymbol {
	// 	var res = make([]DocumentSymbol, 0)
	// 	for _, name := range scope.Names() {
	// 		var sym = DocumentSymbol{
	// 			Name:           name.Name,
	// 			Kind:           NodeSymbolKind(name.Node),
	// 			Range:          name.Node.Range(),
	// 			SelectionRange: name.Node.Range(),
	// 		}
	// 		if name.Node.Scope() != scope {
	// 			sym.Children = process_scope(name.Node.Scope())
	// 		}
	// 		res = append(res, sym)
	// 	}
	// 	// for _, scope := range scope.AllNames()
	// 	return res
	// }
	// var res = process_scope(file.RootScope())

	return nil, nil
}
