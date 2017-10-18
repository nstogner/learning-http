package main

import (
	"io/ioutil"
	"log"
	"net"

	"github.com/nstogner/learning-http/5-http-implementation/http"
)

func main() {
	s := http.Server{
		Handler: handler{},
	}

	l, err := net.Listen("tcp", ":7000")
	if err != nil {
		log.Fatalln("unable to listen:", err)
	}
	if err := s.Serve(l); err != nil {
		log.Fatalln("unable to serve:", err)
	}
}

// handler implements the http.Handler interface.
type handler struct{}

func (h handler) ServeHTTP(w *http.Response, r *http.Request) {
	log.Println("serving http")

	btys, _ := ioutil.ReadAll(r.Body)
	log.Printf("read request body: %q\n", string(btys))

	w.Write([]byte("howdy"))
}
