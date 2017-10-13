package http

import (
	"bufio"
	"bytes"
	"io"
	"net"
)

type Handler interface {
	ServeHTTP(*ResponseWriter, *Request)
}

type ResponseWriter struct {
	Status  int
	Headers map[string]string

	// We know we need a buffer because we need to know what the Content-Length
	// header should be before we can write the body to the connection
	buf bytes.Buffer
}

func (rw *ResponseWriter) Write(b []byte) (int, error) {
	// TODO
	return 0, nil
}

type Request struct {
	// Parsed out header fields
	Method  string
	URI     string
	Proto   string
	Headers map[string]string

	// A way to read bytes from the body
	Body io.Reader

	// Should the connection be terminated after the response is sent?
	keepalive bool
}

type Server struct {
	Handler Handler
}

func (s *Server) Serve(l net.Listener) error {
	defer l.Close()

	for {
		nc, err := l.Accept()
		if err != nil {
			return err
		}

		hc := httpConn{nc, s.Handler}
		go hc.handle()
	}
	return nil
}

type httpConn struct {
	netConn net.Conn
	handler Handler
}

func (hc *httpConn) handle() {
	br := bufio.NewReader(hc.netConn)

	for {
		req, err := readRequest(br)
		if err != nil {
			// TODO: Send bad request
			break
		}

		res := &ResponseWriter{}

		hc.handler.ServeHTTP(res, req)

		// TODO: Send response back

		if !req.keepalive {
			break
		}
	}
}

func readRequest(r io.Reader) (*Request, error) {
	var req Request

	// TODO: Parse request line

	for {
		// TODO: Parse headers
	}

	// TODO: Assign body reader

	return &req, nil
}
