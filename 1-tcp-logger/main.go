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

	for {
		ln, err := r.ReadString('\n')
		if err != nil {
			log.Println("unable to read from conn:", err)
			break
		}

		log.Printf("read line: %q", ln)
	}
}
