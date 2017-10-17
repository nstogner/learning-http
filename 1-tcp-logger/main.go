package main

import (
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
		// Accept blocks until there is an incoming connection.
		conn, err := l.Accept()
		if err != nil {
			log.Println("unable to accept:", err)
			break
		}

		serve(conn)
	}
}

// serve manages reading from a connection.
func serve(c net.Conn) {
	defer c.Close()

	// Create a buffer of length = 1.
	// Try experimenting with different lengths.
	buf := make([]byte, 1)

	for {
		if _, err := c.Read(buf); err != nil {
			log.Println("unable to read from conn:", err)
			return
		}

		log.Printf("buffer bits: %08b: %q", buf, buf)
	}
}
