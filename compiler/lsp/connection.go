package main

import (
	"encoding/json"
	"fmt"
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

type LspRequestHandler func(req *LspRequest) error

var handlers = make(map[string]LspRequestHandler)

type LspRequest struct {
	Conn           *LspConnection
	Id             int
	Method         string
	IsNotification bool
	Params         *fastjson.Value
}

func (r *LspRequest) RawParams() []byte {
	return r.Params.MarshalTo(nil)
}

func (r *LspRequest) ReplyEmpty() error {
	return nil
}

func (r *LspRequest) Notify(method string, params interface{}) {
	mars, _ := json.Marshal(map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  method,
		"params":  params,
	})

	if _, err := fmt.Fprintf(r.Conn, "Content-Length: %v\r\n\r\n%s", len(mars), mars); err != nil {
		log.Fatal("err not nil", err)
	}
}

// Reply to the request if it was one.
func (r *LspRequest) Reply(val interface{}) {
	if r.IsNotification {
		log.Print(red("error "), r.Method, " is a notification and must not be replied to")
		return
	}

	mars, _ := json.Marshal(map[string]interface{}{
		"id":      r.Id,
		"jsonrpc": "2.0",
		"result":  val,
	})
	log.Print("replying ", string(mars))
	if _, err := fmt.Fprintf(r.Conn, "Content-Length: %v\r\n\r\n%s", len(mars), mars); err != nil {
		log.Fatal("err not nil", err)
	}
}

// Send back an error to the requester
func (r *LspRequest) ReplyError(val interface{}) {
	if r.IsNotification {
		log.Print(red("error "), r.Method, " is a notification and must not be replied to")
		return
	}

	mars, _ := json.Marshal(map[string]interface{}{
		"id":      r.Id,
		"jsonrpc": "2.0",
		"error":   val,
	})
	log.Print("replying error ", string(mars))
	if _, err := fmt.Fprintf(r.Conn, "Content-Length: %v\r\n\r\n%s", len(mars), mars); err != nil {
		log.Fatal("err not nil", err)
	}
}

// LspConnection is in charge of reading a request and sending back the results
// It also holds a reference to a compiler session, where it will manipulate the
// files as they get edited.
type LspConnection struct {
	io.ReadWriteCloser
	receivedShutdown bool
	Solution         *zoe.Solution
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
	}
}

// message is really the JSON
// the buffer is guaranteed to be available
func (l *LspConnection) HandleMessage(message []byte) {

	// log.Print("-->", string(message))
	parsed, err := fastjson.ParseBytes(message)
	if err != nil {
		log.Fatal("invalid json received")
	}
	// log.Print(string(message))

	is_notification := !parsed.Exists("id")

	id := parsed.GetInt("id") // we need the id to reply to the request
	method := string(parsed.GetStringBytes("method"))

	if hld, ok := handlers[method]; ok {
		req := LspRequest{
			Conn:           l,
			Id:             id,
			IsNotification: is_notification,
			Params:         parsed.Get("params"),
			Method:         method,
		}
		if is_notification {
			log.Print(green("!"), " notified ", cyan(method))
		} else {
			log.Print(green("*"), " handling ", green(method))
		}
		if err := hld(&req); err != nil {
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

	// make the buffer big enough to handle the request.
	var buf = make([]byte, 128*1024)
	var total_read = 0
	n, err := l.Read(buf)
	for n > 0 && err == nil {
		total_read += n

		indices := re_len.FindSubmatchIndex(buf)
		if indices == nil {
			// this shouldn't happen and means that we have faulty input.
			log.Fatal("Ooops, didn't find content-length")
		}
		length, _ := strconv.Atoi(string(buf[indices[2]:indices[3]]))
		start := indices[1]

		for total_read < start+length {
			// fill up the remainder
			n, err = l.Read(buf[start : start+length])
			if n == 0 && err != nil {
				log.Fatal("WTF")
			}
			total_read += n
		}

		// log.Print("length", length, " -- ", indices[1])

		for buf[start] != '{' {
			start++
		}

		packet := buf[start : start+length]
		msg := make([]byte, len(packet))
		copy(msg, packet)
		l.HandleMessage(msg)
		total_read = 0
		// parser := fastjson

		// This implementation is pretty naïve right now. It just looks for the header and watches how long the request oughts to be
		if err != nil {
			// Most likely connection was closed
			break
		}

		n, err = l.Read(buf)
		// log.Print(n, err)
	}

	// This connection looks done, so we're closing it.
	log.Print(green("closed connection"))
	if e := l.Close(); e != nil {
		log.Print(e)
	}
}
