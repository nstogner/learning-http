package main

import (
	"bufio"
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

		serve(conn)
	}
}

func serve(c net.Conn) {
	defer c.Close()

	r := bufio.NewReader(c)

	// Read request line
	ln0, err := r.ReadString('\n')
	if err != nil {
		return
	}
	log.Printf("read request line: %q", ln0)

	// Read headers
	for {
		ln, err := r.ReadString('\n')
		if err != nil {
			break
		}

		// An empty line with crlf marks the end of the headers
		if ln == "\r\n" {
			break
		}

		// Must be a header if we have reached this point
		log.Printf("read request header: %q", ln)
	}

	// Ignore body

	// Write response
	c.Write([]byte("HTTP/1.1 200 OK\r\nContent-Length: 0\r\n\r\n"))
}
