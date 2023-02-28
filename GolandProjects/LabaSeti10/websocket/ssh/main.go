package main

import (
	"LabaSeti10/api/sshWeb"
	"github.com/gorilla/websocket"
	"github.com/julienschmidt/httprouter"
	"golang.org/x/crypto/ssh"
	"log"
	"net/http"
)

func main() {
	mux := httprouter.New()

	config := &ssh.ClientConfig{
		User: "test",
		Auth: []ssh.AuthMethod{
			ssh.Password("SDHBCXdsedfs222"),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	upgrade := &websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	sshWeb.Register(mux, config, upgrade)

	log.Fatal(http.ListenAndServe("localhost:8282", mux))
}
