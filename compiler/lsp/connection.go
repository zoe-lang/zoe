package main

import (
	"io"
	"log"
	"os"
	"regexp"

	zoe "github.com/ceymard/zoe/compiler"
	"github.com/creachadair/jrpc2"
	"github.com/creachadair/jrpc2/channel"
	"github.com/creachadair/jrpc2/handler"
	"github.com/fatih/color"
)

var fred = color.New(color.FgRed, color.Bold)
var red = fred.SprintFunc()
var fgreen = color.New(color.FgGreen)
var green = fgreen.SprintFunc()
var cyan = color.New(color.FgCyan).SprintFunc()
var yel = color.New(color.FgYellow).SprintFunc()
var mag = color.New(color.FgMagenta).SprintFunc()
var blue = color.New(color.FgBlue).SprintFunc()
var grey = color.New(color.Faint).SprintFunc()
var bblue = color.New(color.FgHiBlue, color.Bold).SprintFunc()

// LspConnection is in charge of reading a request and sending back the results
// It also holds a reference to a compiler session, where it will manipulate the
// files as they get edited.
type LspConnection struct {
	io.ReadWriteCloser
	Server           *jrpc2.Server
	receivedShutdown bool
	Solution         *zoe.Solution
}

type handlerBuilder func(l *LspConnection, mp handler.Map)

var builders = make([]handlerBuilder, 0)

func addHandler(bld handlerBuilder) {
	builders = append(builders, bld)
}

// The Zoe LSP should be capable of being multi user.
// This means the usage links (what symbol refer to what symbol) should be a "per-session" thing (?)
//   or at least should be aware that other versions of the same files may be open at the same time.
// Another use case could be the pooling of resources ; if the lsp is launched with all of its type checking,
// it should be able to create a binary on the fly.

var re_len = regexp.MustCompile(`Content-Length: (\d+)`)

func NewConnection(conn io.ReadWriteCloser) *LspConnection {
	var l = &LspConnection{
		ReadWriteCloser: conn,
		Solution:        zoe.NewSolution(),
	}

	var mp = make(handler.Map)
	for _, hld := range builders {
		hld(l, mp)
	}

	l.Server = jrpc2.NewServer(mp,
		&jrpc2.ServerOptions{
			AllowPush: true,
			Logger:    jrpc2.StdLogger(log.New(os.Stderr, "-- ", 0)),
		})

	return l
}

func (l *LspConnection) Listen() error {
	defer func() {
		r := recover()
		if r != nil {
			log.Printf("listen recovered panic: %s", r)
		}
	}()

	log.Printf("hlb-langserver listening")
	s := l.Server.Start(channel.Header("")(l.ReadWriteCloser, l.ReadWriteCloser))
	return s.Wait()

}
