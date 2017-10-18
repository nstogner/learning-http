package http_test

import (
	"bytes"
	"io/ioutil"
	"net"
	stdhttp "net/http"
	"strconv"
	"testing"

	"github.com/nstogner/learning-http/5-http-implementation/http"
)

func TestServe(t *testing.T) {
	const (
		reqBody = `{"abc":123}`
		resBody = `{"xyz":456}`
	)

	// Start a listener on any free port.
	l, err := net.Listen("tcp", "")
	if err != nil {
		t.Fatal("unable to listen:", err)
	}

	var th testHandler
	th.response.status = 201
	th.response.body = []byte(resBody)
	server := http.Server{
		Handler: &th,
	}
	go func() {
		if err := server.Serve(l); err != nil {
			t.Fatal("unable to serve:", err)
		}
	}()

	_, port, err := net.SplitHostPort(l.Addr().String())
	if err != nil {
		t.Fatal("unable to split host and port:", err)
	}
	resp, err := stdhttp.Post("http://localhost:"+port, "application/json", bytes.NewReader([]byte(reqBody)))
	if err != nil {
		t.Fatal("post failed:", err)
	}

	// Check that the response set by the testHandler match the response that
	// was received.
	if th.response.status != resp.StatusCode {
		t.Fatalf("expected status code %s, got: %s", th.response.status, resp.StatusCode)
	}

	// Check that expected response headers are set.
	if len(resp.Header.Get("Date")) == 0 {
		t.Fatal("expected header 'Date' to be set")
	}
	if cl, err := strconv.Atoi(resp.Header.Get("Content-Length")); err == nil {
		if exp := len(resBody); cl != exp {
			t.Fatalf("expected header 'Content-Length' = %v, got: %v", exp, cl)
		}
	} else {
		t.Fatal("unable to parse 'Content-Length' header:", err)
	}

	// Check that the fields that were parsed by the Server match what was
	// requested.
	if exp := "POST"; th.request.method != exp {
		t.Fatalf("expected method %s, got: %s", exp, th.request.method)
	}
	if string(th.request.body) != reqBody {
		t.Fatalf("expected body '%s', got: '%s'", reqBody, string(th.request.body))
	}
}

// testHandler records parsed request fields and sets predefined response
// fields.
type testHandler struct {
	request struct {
		method string
		body   []byte
	}
	response struct {
		status int
		body   []byte
	}
}

// ServeHTTP satisfies the http.Handler interface.
func (th *testHandler) ServeHTTP(res *http.Response, req *http.Request) {
	th.request.method = req.Method

	// Ignore error because if there is an error it will be caught by body
	// comparison in test.
	th.request.body, _ = ioutil.ReadAll(req.Body)

	res.Status = th.response.status
	res.Write(th.response.body)
}
