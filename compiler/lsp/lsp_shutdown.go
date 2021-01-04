package main

import (
	"context"

	"github.com/creachadair/jrpc2/handler"
)

func init() {
	addHandler(func(l *LspConnection, mp handler.Map) {
		mp["shutdown"] = handler.New(l.HandleShutdown)
		mp["exit"] = handler.New(l.HandleExit)
	})
}

func (l *LspConnection) HandleShutdown(_ context.Context) error {
	// FIXME should probably do some cleanup...
	l.receivedShutdown = true
	return nil
}

func (l *LspConnection) HandleExit(_ context.Context) error {
	// if req.Conn.receivedShutdown {
	// 	os.Exit(0)
	// }
	// os.Exit(1)
	return nil
}
