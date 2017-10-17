package http_test

import (
	"bytes"
	"io/ioutil"
	"net"
	stdhttp "net/http"
	"testing"

	"github.com/nstogner/learning-http/4-http-outline/http"
)

func TestServe(t *testing.T) {
	const body = `{"abc":123}`

	l, err := net.Listen("tcp", "")
	if err != nil {
		t.Fatal("unable to listen:", err)
	}

	var th testHandler
	th.response.status = 201
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
	resp, err := stdhttp.Post("http://localhost:"+port, "application/json", bytes.NewReader([]byte(body)))
	if err != nil {
		t.Fatal("post failed:", err)
	}
	if resp.StatusCode != th.response.status {
		t.Fatalf("expected status code %s, got: %s", th.response.status, resp.StatusCode)
	}
	if string(th.request.body) != body {
		t.Fatalf("expected body '%s', got: '%s'", body, string(th.request.body))
	}
}

type testHandler struct {
	request struct {
		method string
		body   []byte
	}
	response struct {
		status int
	}
}

func (th *testHandler) ServeHTTP(res *http.Response, req *http.Request) {
	th.request.method = req.Method

	// Ignore error because if there is an error it will be caught by body
	// comparison in test.
	th.request.body, _ = ioutil.ReadAll(req.Body)

	res.Status = th.response.status
}
