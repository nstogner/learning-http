package http

import (
	"io"
	"net"
)

type Handler interface {
	ServeHTTP(*ResponseWriter, *Request)
}

type ResponseWriter struct {
	Status  int
	Headers map[string]string
}

func (rw *ResponseWriter) Write(b []byte) (int, error) {
	// TODO
	return 0, nil
}

type Request struct {
	Method  string
	Path    string
	Proto   string
	Headers map[string]string

	keepalive bool
}

type Server struct {
	Addr    string
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
	req, err := readRequest(hc.netConn)
	if err != nil {
		// TODO: Send bad request
	}

	res := &ResponseWriter{}

	hc.handler.ServeHTTP(res, req)

	if req.keepalive {
		hc.handle()
	}
}

func readRequest(r io.Reader) (*Request, error) {
	var req Request

	// TODO: Parse request line

	for {
		// TODO: Parse headers
	}

	return &req, nil
}
