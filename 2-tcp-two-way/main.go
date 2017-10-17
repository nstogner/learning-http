package main

import (
	"bufio"
	"io"
	"log"
	"net"
	"os"
	"strings"
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

// serve manages reading and writing to a connection.
func serve(c net.Conn) {
	defer c.Close()

	// The bufio Reader provides some nice convenience functions for reading
	// up until a particular character is found
	buf := bufio.NewReader(c)

	for {
		// Read up to and including the next newline character
		// (the second byte of a crlf)
		ln, err := buf.ReadString('\n')
		if err != nil {
			log.Println("unable to read from conn:", err)
			break
		}

		// NOTE: Telnet will send a crlf when you hit enter so we will strip
		// that off here
		cmd := strings.TrimSuffix(ln, "\r\n")
		handle(c, cmd)
	}
}

// handle accepts a Writer so that it can respond to a given parsed command
// string.
func handle(w io.Writer, cmd string) {
	switch cmd {
	case "BEEP":
		log.Print("beeping!")
		os.Stdout.Write([]byte("\u0007"))

		w.Write([]byte("ACCEPTED\r\n"))
	default:
		w.Write([]byte("REJECTED\r\n"))
	}
}
