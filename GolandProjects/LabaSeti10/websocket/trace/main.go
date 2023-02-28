package main

import (
	"LabaSeti10/api/trace"
	"github.com/gorilla/websocket"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
)

func main() {
	mux := httprouter.New()

	upgrade := &websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	trace.Register(mux, upgrade)

	log.Fatal(http.ListenAndServe("localhost:8383", mux))
}
