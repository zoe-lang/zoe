package main

import (
	"log"
	"net"
	"os"
)

// for now we'll do an stdio lsp
// but there should also be a way to create it over file socket or TCP socket
const SockAddr = "/tmp/echo.sock"

func main() {
	if err := os.RemoveAll(SockAddr); err != nil {
		log.Fatal(err)
	}

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
		go theconn.Listen()
	}
}
