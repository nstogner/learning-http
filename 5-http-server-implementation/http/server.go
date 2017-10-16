package http

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
	"time"
)

const (
	http10 = "HTTP/1.0"
)

// statusTitles map HTTP status codes to their titles. This is handy for
// sending the response header.
var statusTitles = map[int]string{
	200: "OK",
	201: "Created",
	202: "Accepted",
	203: "Non-Authoritative Information",
	204: "No Content",
	// TODO: More status codes
}

// Handler responds to a HTTP request.
type Handler interface {
	ServeHTTP(*ResponseWriter, *Request)
}

// ResponseWriter is used to construct a HTTP response.
// TODO: Change name b/c it might be confused with the std lib interface
// (ResponseRecorder?)
type ResponseWriter struct {
	Status  int
	Headers map[string]string

	proto string
	buf   bytes.Buffer
}

// Write writes data to a buffer which is later flushed to the network
// connection.
func (rw *ResponseWriter) Write(b []byte) (int, error) {
	return rw.buf.Write(b)
}

// writeTo writes an HTTP response with headers and buffered body to a writer.
func (rw *ResponseWriter) writeTo(w io.Writer) error {
	if err := rw.writeHeadersTo(w); err != nil {
		return err
	}

	if _, err := rw.buf.WriteTo(w); err != nil {
		return err
	}

	return nil
}

// writeHeadersTo writes HTTP headers to a writer.
func (rw *ResponseWriter) writeHeadersTo(w io.Writer) error {
	statusText, ok := statusTitles[rw.Status]
	if !ok {
		return fmt.Errorf("unsupported status code: %v", rw.Status)
	}

	rw.Headers["Date"] = time.Now().UTC().Format("Mon, 02 Jan 2006 15:04:05 GMT")
	rw.Headers["Content-Length"] = strconv.Itoa(rw.buf.Len())

	// https://www.w3.org/Protocols/rfc2616/rfc2616-sec6.html
	headers := fmt.Sprintf("%s %v %s\r\n", rw.proto, rw.Status, statusText)
	for k, v := range rw.Headers {
		headers += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	headers += "\r\n"

	if _, err := w.Write([]byte(headers)); err != nil {
		return err
	}

	return nil
}

// Request represents a HTTP request sent to a server.
type Request struct {
	Method  string
	URI     string
	Proto   string
	Headers map[string]string

	Body io.Reader

	keepalive bool
}

// httpConn handles persistent HTTP connections.
type httpConn struct {
	netConn net.Conn
	handler Handler
}

// serve reads and responds to one or many HTTP requests off of a single TCP
// connection.
func (hc *httpConn) serve() {
	defer hc.netConn.Close()

	br := bufio.NewReader(hc.netConn)

	for {
		req, err := readRequest(br)
		if err != nil {
			const bad = "HTTP/1.1 400 Bad Request\r\nConnection: close\r\n\r\n"
			hc.netConn.Write([]byte(bad))
			return
		}

		res := ResponseWriter{
			Status:  200,
			Headers: make(map[string]string),
			proto:   req.Proto,
		}

		hc.handler.ServeHTTP(&res, req)

		if err := res.writeTo(hc.netConn); err != nil {
			return
		}

		if !req.keepalive {
			return
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

		// Spawn off a goroutine so we can accept other connections.
		go hc.serve()
	}
}

// readRequest generates a Request object by parsing text from a bufio.Reader.
func readRequest(buf *bufio.Reader) (*Request, error) {
	req := Request{
		Headers: make(map[string]string),
	}

	// Read the HTTP request line (first line).
	if ln0, err := readHTTPLine(buf); err == nil {
		var ok bool
		if req.Method, req.URI, req.Proto, ok = parseRequestLine(ln0); !ok {
			return nil, fmt.Errorf("malformed request line: %q", ln0)
		}
	}

	// Read each subsequent header.
	for {
		ln, err := readHTTPLine(buf)
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

	// Limit the body to the number of bytes specified by Content-Length.
	var cl int64
	if str, ok := req.Headers["content-length"]; ok {
		var err error
		if cl, err = strconv.ParseInt(str, 10, 64); err != nil {
			return nil, err
		}
	}
	req.Body = &io.LimitedReader{R: buf, N: cl}

	// Determine if connection should be closed after request.
	req.keepalive = shouldKeepAlive(req.Proto, req.Headers["connection"])

	return &req, nil
}

// shouldKeepAlive determines whether a connection should be kept alive or
// closed based on the protocol version and "Connection" header.
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

// parseRequestLine attempts to parse the initial line of an HTTP request.
func parseRequestLine(ln string) (method, uri, proto string, ok bool) {
	s := strings.Split(ln, " ")
	if len(s) != 3 {
		return
	}

	return s[0], s[1], s[2], true
}

// parseHeaderLine attempts to parse a standard HTTP header, e.g.
// "Content-Type: application/json".
func parseHeaderLine(ln string) (key, val string, ok bool) {
	s := strings.SplitN(ln, ":", 2)
	if len(s) != 2 {
		return
	}

	return strings.ToLower(s[0]), strings.TrimSpace(s[1]), true
}

// readHTTPLine reads up to a newline feed and strips off the trailing crlf.
func readHTTPLine(buf *bufio.Reader) (string, error) {
	ln, err := buf.ReadString('\n')
	if err != nil {
		return "", err
	}

	return strings.TrimSuffix(ln, "\r\n"), nil
}
