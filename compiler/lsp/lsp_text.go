package main

import "log"

func init() {
	handlers["textDocument/didOpen"] = HandleDidOpen
}

func HandleDidOpen(req *LspRequest) error {
	log.Print(req.Params.String())
	return nil
}
