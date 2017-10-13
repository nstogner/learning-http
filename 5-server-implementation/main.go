package main

import (
	"io/ioutil"
	"log"
	"net"

	"github.com/nstogner/learning-http/5-server-implementation/http"
)

func main() {
	s := http.Server{
		Handler: handler{},
	}

	l, err := net.Listen("tcp", ":7000")
	if err != nil {
		panic(err)
	}
	if err := s.Serve(l); err != nil {
		panic(err)
	}
}

type handler struct{}

func (h handler) ServeHTTP(w *http.ResponseWriter, r *http.Request) {
	log.Println("serving http")

	btys, _ := ioutil.ReadAll(r.Body)
	log.Printf("read request body: %q\n", string(btys))

	w.Write([]byte("howdy"))
}
