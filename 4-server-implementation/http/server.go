package http

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"net"
	"strings"
	"time"
)

const (
	http10 = "HTTP/1.0"
	http11 = "HTTP/1.1"
)

type Handler interface {
	ServeHTTP(*ResponseWriter, *Request)
}

type ResponseWriter struct {
	Status  int
	Headers map[string]string

	proto string

	buf  bytes.Buffer
	conn net.Conn
}

func (rw *ResponseWriter) Write(b []byte) (int, error) {
	return rw.buf.Write(b)
}

func (rw *ResponseWriter) send() error {
	if err := rw.sendHeaders(); err != nil {
		return err
	}
	if _, err := rw.buf.WriteTo(rw.conn); err != nil {
		return err
	}
	return nil
}

func (rw *ResponseWriter) sendHeaders() error {
	statusText, ok := statusTitles[rw.Status]
	if !ok {
		return fmt.Errorf("unsupported status code: %v", rw.Status)
	}

	// https://www.w3.org/Protocols/rfc2616/rfc2616-sec6.html
	statusline := fmt.Sprintf("%s %v %s"+crlf, rw.proto, rw.Status, statusText)
	if _, err := rw.conn.Write([]byte(statusline)); err != nil {
		return err
	}

	dateline := "Date: " + time.Now().UTC().Format("Mon, 02 Jan 2006 15:04:05 GMT") + crlf
	if _, err := rw.conn.Write([]byte(dateline)); err != nil {
		return err
	}

	for k, v := range rw.Headers {
		line := fmt.Sprintf("%s: %s", k, v)
		if _, err := rw.conn.Write([]byte(line + crlf)); err != nil {
			return err
		}
	}

	if _, err := rw.conn.Write([]byte(crlf)); err != nil {
		return err
	}

	return nil
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
		// TODO: Send bad req
		return
	}
	defer func() {
		if req.keepalive {
			hc.handle()
		} else {
			hc.netConn.Close()
		}
	}()

	res := &ResponseWriter{
		Status:  200,
		Headers: make(map[string]string),
		proto:   req.Proto,
		conn:    hc.netConn,
	}

	hc.handler.ServeHTTP(res, req)
	if err := res.send(); err != nil {
		req.keepalive = false
	}
}

func readRequest(r io.Reader) (*Request, error) {
	req := Request{
		Headers: make(map[string]string),
	}

	br := bufio.NewReader(r)

	// First line
	if ln0, err := readHTTPLine(br); err == nil {
		var ok bool
		if req.Method, req.Path, req.Proto, ok = parseRequestLine(ln0); !ok {
			return nil, errors.New("malformed request")
		}
	}

	// Headers
	for {
		ln, err := readHTTPLine(br)
		if err != nil {
			return nil, err
		}

		if len(ln) == 0 {
			break
		}

		if key, val, ok := parseHeaderLine(ln); ok {
			req.Headers[key] = val
		}
	}

	// Keep alive
	req.keepalive = shouldKeepAlive(req.Proto, req.Headers["connection"])

	return &req, nil
}

func shouldKeepAlive(proto, connHeader string) bool {
	switch proto {
	case http10:
		if connHeader == "keep-alive" {
			return true
		}
		return false
	default:
		if connHeader == "close" {
			return false
		}
		return true
	}
}

func parseRequestLine(ln string) (method, path, proto string, ok bool) {
	s := strings.Split(ln, " ")
	if len(s) != 3 {
		return
	}
	return s[0], s[1], s[2], true
}

func parseHeaderLine(ln string) (key, val string, ok bool) {
	s := strings.Split(ln, ":")
	if len(s) != 2 {
		return
	}
	return strings.ToLower(s[0]), strings.TrimSpace(s[1]), true
}

func readHTTPLine(br *bufio.Reader) (string, error) {
	ln, err := br.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSuffix(ln, "\r\n"), nil
}

const crlf = "\r\n"

var statusTitles = map[int]string{
	200: "OK",
	201: "Created",
	202: "Accepted",
	203: "Non-Authoritative Information",
	204: "No Content",
}
