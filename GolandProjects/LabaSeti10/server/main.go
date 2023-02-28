package main

import (
	"LabaSeti10/api/server"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
)

func main() {
	mux := httprouter.New()
	server.Register(mux)
	log.Fatal(http.ListenAndServe(":8080", mux))
}
