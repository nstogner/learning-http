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
			log.Println("unable to accept", err)
			break
		}

		serve(conn)
	}
}

// serve manages reading and writing to a connection.
func serve(c net.Conn) {
	defer c.Close()

	buf := bufio.NewReader(c)

	// Read request line
	// e.g. "GET /abc HTTP/1.1"
	ln0, err := buf.ReadString('\n')
	if err != nil {
		return
	}
	log.Printf("read request line: %q", ln0)

	// Read headers
	// e.g. "Content-Type: application/json"
	for {
		ln, err := buf.ReadString('\n')
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

	// Ignore the request body for now

	// Write response
	c.Write([]byte("HTTP/1.1 200 OK\r\nContent-Length: 0\r\n\r\n"))
}
