package main

import (
	"context"
	"log"
	"net"
	"os"

	"github.com/intel-go/fastjson"
	"github.com/k0kubun/pp"
	"github.com/osamingo/jsonrpc"
	"github.com/sourcegraph/go-lsp"
)

var t lsp.Position

// for now we'll do an stdio lsp
// but there should also be a way to create it over file socket or TCP socket
const SockAddr = "/tmp/echo.sock"

type InitializeHandler struct {
}

func (h InitializeHandler) ServeJSONRPC(c context.Context, params *fastjson.RawMessage) (interface{}, *jsonrpc.Error) {

	var p lsp.InitializeParams
	if err := jsonrpc.Unmarshal(params, &p); err != nil {
		return nil, err
	}

	pp.Print(p)
	log.Print("got here...")

	return lsp.InitializeResult{
		Capabilities: lsp.ServerCapabilities{},
	}, nil
}

func main() {
	if err := os.RemoveAll(SockAddr); err != nil {
		log.Fatal(err)
	}

	mr := jsonrpc.NewMethodRepository()
	mr.RegisterMethod("initialize", InitializeHandler{}, lsp.InitializeParams{}, lsp.InitializeResult{})

	l, err := net.Listen("unix", SockAddr)
	if err != nil {
		log.Fatal("listen error:", err)
	}
	defer l.Close()

	for {
		// Accept new connections, dispatching them to echoServer
		// in a goroutine.
		conn, err := l.Accept()
		if err != nil {
			log.Fatal("accept error:", err)
		}
		log.Printf("Client connected [%s]", conn.RemoteAddr().Network())
		theconn := NewConnection(conn)
		go theconn.ProcessIncomingRequests()
	}
}
