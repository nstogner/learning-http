package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"

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
	json.NewEncoder(os.Stderr).Encode(r)

	btys, _ := ioutil.ReadAll(r)
	fmt.Printf("read request body: %q\n", string(btys))

	w.Write([]byte("hey"))
}
