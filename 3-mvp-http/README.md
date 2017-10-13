# 3. Minimal HTTP Server

We will define a "minimal HTTP server" as a server that can successfully respond to a simple curl GET request. We will start out by observing what happens when we issue this request to our TCP logger server:

```sh
go run ./1-tcp-logger/main.go
```

```sh
# In another shell
curl localhost:7000/abc -v
```

We can observe that request header lines are seperated by crlf's and the header ends in two back-to-back crlf's. Hint: response headers are the same, with the exception of the first line. An example response header line: `HTTP/1.1 200 OK`.

## Bonus

- Parse the request line (the first line) to extract the HTTP method, URI, and protocol version.

