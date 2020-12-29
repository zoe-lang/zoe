package main

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"github.com/valyala/fastjson"
)

type LspRequestHandler func(req *LspRequest) error

var handlers = make(map[string]LspRequestHandler)

type LspRequest struct {
	Lock           sync.Mutex
	Conn           *LspConnection
	Id             int
	Method         string
	IsNotification bool
	Params         *fastjson.Value

	replied bool
}

func (r *LspRequest) RawParams() []byte {
	return r.Params.MarshalTo(nil)
}

func (r *LspRequest) Cancel() {
	r.Lock.Lock()
	defer r.Lock.Unlock()

	if r.replied {
		return
	}
	r.ReplyError(map[string]interface{}{
		"error":   -32800, // request cancelled
		"message": "canceled.",
	})
	r.replied = true
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
	r.Lock.Lock()
	defer r.Lock.Unlock()

	if r.replied {
		// do not reply if we were cancelled
		return
	}

	if r.IsNotification {
		log.Print(red("error "), r.Method, " is a notification and must not be replied to")
		return
	}

	mars, _ := json.Marshal(map[string]interface{}{
		"id":      r.Id,
		"jsonrpc": "2.0",
		"result":  val,
	})
	delete(r.Conn.RequestMap, r.Id)
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
	delete(r.Conn.RequestMap, r.Id)
	log.Print("replying error ", string(mars))
	if _, err := fmt.Fprintf(r.Conn, "Content-Length: %v\r\n\r\n%s", len(mars), mars); err != nil {
		log.Fatal("err not nil", err)
	}
}
