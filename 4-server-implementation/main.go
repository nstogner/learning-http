package main

import (
	"encoding/json"
	"log"
	"net"
	"os"
	"time"

	"github.com/nstogner/learning-http/4-server-implementation/http"
)

func main() {
	s := http.Server{
		Addr:    ":8080",
		Handler: handler{},
	}

	l, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}
	if err := s.Serve(l); err != nil {
		panic(err)
	}
}

type handler struct{}

func (h handler) ServeHTTP(w *http.ResponseWriter, r *http.Request) {
	json.NewEncoder(os.Stderr).Encode(r)
	log.Println("serving http")
	time.Sleep(10 * time.Second)
}
