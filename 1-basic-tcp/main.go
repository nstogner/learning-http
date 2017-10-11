package main

import (
	"bytes"
	"log"
	"net"
)

func main() {
	l, err := net.Listen("tcp", ":7000")
	if err != nil {
		log.Fatalln("unable to listen:", err)
	}
	defer l.Close()

	log.Print("listening")

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Println("unable to accept")
			break
		}

		go serve(conn)
	}
}

func serve(c net.Conn) {
	defer c.Close()
	for {
		btys := make([]byte, 30)
		n, err := c.Read(btys)
		if err != nil {
			break
		}
		log.Printf("read %v bytes: %q", n, string(btys))

		if bytes.Contains(btys, []byte("C")) {
			log.Print("read 'C': closing connection")
			break
		}
	}
}
