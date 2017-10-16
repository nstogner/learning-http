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
		// Accept blocks until there is an incoming connection
		conn, err := l.Accept()
		if err != nil {
			log.Print("unable to accept")
			break
		}

		serve(conn)
	}
}

// serve manages reading from a connection.
func serve(c net.Conn) {
	defer c.Close()

	// The bufio Reader provides some nice convenience functions for reading
	// up until a particular character is found
	r := bufio.NewReader(c)

	for {
		// Read up to and including the next newline character
		ln, err := r.ReadString('\n')
		if err != nil {
			log.Println("unable to read from conn:", err)
			break
		}

		log.Printf("read line: %q", ln)
	}
}
