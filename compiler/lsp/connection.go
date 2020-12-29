package main

import (
	"bytes"
	"io"
	"log"
	"regexp"
	"strconv"

	zoe "github.com/ceymard/zoe/compiler"
	"github.com/fatih/color"
	"github.com/valyala/fastjson"
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
	receivedShutdown bool
	Solution         *zoe.Solution
	RequestMap       map[int]*LspRequest
}

// The Zoe LSP should be capable of being multi user.
// This means the usage links (what symbol refer to what symbol) should be a "per-session" thing (?)
//   or at least should be aware that other versions of the same files may be open at the same time.
// Another use case could be the pooling of resources ; if the lsp is launched with all of its type checking,
// it should be able to create a binary on the fly.

var re_len = regexp.MustCompile(`Content-Length: (\d+)`)

func NewConnection(conn io.ReadWriteCloser) *LspConnection {
	return &LspConnection{
		ReadWriteCloser: conn,
		Solution:        zoe.NewSolution(),
		RequestMap:      make(map[int]*LspRequest),
	}
}

// message is really the JSON
// the buffer is guaranteed to be available
func (l *LspConnection) HandleMessage(message []byte) {

	// log.Print("-->", string(message))
	parsed, err := fastjson.ParseBytes(message)
	if err != nil {
		log.Print("!!!", string(message))
		log.Fatal("invalid json received")
	}
	// log.Print(string(message))

	is_notification := !parsed.Exists("id")

	id := parsed.GetInt("id") // we need the id to reply to the request
	method := string(parsed.GetStringBytes("method"))

	if method == "$/cancelRequest" {
		if req, ok := l.RequestMap[id]; ok {
			req.Cancel()
			return
		}
	}

	if hld, ok := handlers[method]; ok {
		req := LspRequest{
			Conn:           l,
			Id:             id,
			IsNotification: is_notification,
			Params:         parsed.Get("params"),
			Method:         method,
		}
		l.RequestMap[id] = &req
		if is_notification {
			log.Print(green("!"), " notified ", cyan(method))
		} else {
			log.Print(green("*"), " handling ", green(method))
		}
		if err := hld(&req); err != nil {
			delete(l.RequestMap, id)
			req.ReplyError(map[string]interface{}{
				"code":    -32700,
				"message": err.Error(),
			})
		} // should handle any errors returned by hld
	} else {
		log.Print(red("error "), "no handler found for method '", yel(method), "'")
	}

	// showmessage: { type: number, message: string }
	// error = 1, warning = 2, info = 3, log = 4

}

func (l *LspConnection) ProcessIncomingRequests() {
	// process incoming request to get the JSON bounds
	// JSON is valid as long as the processing request is handling it.

	// At first, we look for the headers. Once we know the length of the request, we use that
	// to read just enough to get the json.

	// When it gets the json, it just scans the asked method and dispatches it to the associated handler.

	var chunkSize = 2048
	// make the buffer big enough to handle the request.
	var buf = bytes.NewBuffer(make([]byte, 0))
	var chunk = make([]byte, chunkSize)

	var total_read = 0
	n, err := l.Read(chunk)
	for n > 0 && err == nil {
		total_read += n
		_, _ = buf.Write(chunk[:n])
		bytes := buf.Bytes()

		indices := re_len.FindSubmatchIndex(bytes)
		if indices == nil {
			log.Fatal("Ooops, didn't find content-length")
		}
		length, _ := strconv.Atoi(string(bytes[indices[2]:indices[3]]))
		start := indices[1] + 4

		for total_read < start+length {
			// fill up the remainder
			if start+length-total_read > chunkSize {
				n, err = l.Read(chunk)
			} else {
				n, err = l.Read(chunk[:start+length-total_read])
			}
			if n == 0 && err != nil {
				log.Fatal("WTF")
			}
			_, _ = buf.Write(chunk[:n])
			total_read += n
		}

		wholeBuf := buf.Bytes()
		// we've read it all, now we're getting parsing !
		for wholeBuf[start] != '{' {
			start++
		}

		packet := wholeBuf[start : start+length]
		rest := wholeBuf[start+length:]
		msg := make([]byte, len(packet))
		copy(msg, packet)
		// log.Print(string(msg))
		l.HandleMessage(msg)
		// parser := fastjson

		// This implementation is pretty naïve right now. It just looks for the header and watches how long the request oughts to be
		if err != nil {
			// Most likely connection was closed
			break
		}

		// keep the buffer
		buf.Reset()
		// log.Print("-->", string(rest))
		_, _ = buf.Write(rest) // write the left overs at the beginning
		total_read = len(rest)
		n, err = l.Read(chunk)
		// log.Print(n, err, string(chunk[:n]))
	}

	// This connection looks done, so we're closing it.
	log.Print(green("closed connection"))
	if e := l.Close(); e != nil {
		log.Print(e)
	}
}
