package main

import (
	"bufio"
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

	ln, err := r.ReadString('\n')
	if err != nil {
		log.Println("unable to read from conn:", err)
		return
	}

	cmd := strings.TrimSuffix(ln, "\r\n")
	handle(c, cmd)
}

func handle(c net.Conn, cmd string) {
	switch cmd {
	case "BEEP":
		log.Print("beeping!")
		os.Stdout.Write([]byte("\u0007"))

		c.Write([]byte("ACCEPTED\r\n"))
	default:
		c.Write([]byte("REJECTED\r\n"))
	}
}
