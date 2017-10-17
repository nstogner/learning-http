package http

import (
	"bufio"
	"bytes"
	"io"
	"net"
)

// Handler responds to a HTTP request.
type Handler interface {
	// ServeHTTP takes a Response struct rather than a ResponseWriter interface
	// like the standard library to keep things simple.
	ServeHTTP(*Response, *Request)
}

// Response is used to construct a HTTP response.
type Response struct {
	Status  int
	Headers map[string]string

	// We know we need a buffer because we need to know what the Content-Length
	// header should be before we can write the body to the connection
	buf bytes.Buffer
}

func (res *Response) Write(b []byte) (int, error) {
	// TODO
	return 0, nil
}

// Request represents a HTTP request sent to a server.
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

// httpConn handles persistent HTTP connections.
type httpConn struct {
	netConn net.Conn
	handler Handler
}

// serve reads and responds to one or many HTTP requests off of a single
// connection.
func (hc *httpConn) serve() {
	br := bufio.NewReader(hc.netConn)

	for {
		req, err := readRequest(br)
		if err != nil {
			// TODO: Send bad request
			break
		}

		res := &Response{}

		hc.handler.ServeHTTP(res, req)

		// TODO: Send response back

		if !req.keepalive {
			break
		}
	}
}

// Server wraps a Handler and manages a network listener.
type Server struct {
	Handler Handler
}

// Serve accepts incoming HTTP connections and handles them in a new goroutine.
func (s *Server) Serve(l net.Listener) error {
	defer l.Close()

	for {
		nc, err := l.Accept()
		if err != nil {
			return err
		}

		hc := httpConn{nc, s.Handler}
		// Spawn off a goroutine so we can accept other connections
		go hc.serve()
	}
	return nil
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
