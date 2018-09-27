# 3. Minimal HTTP Server

## What is a protocol anyways?

- In terms of networking, a protocol is just a set of rules that endpoints use to communicate
- These rules specify the manner in which data is sent and formatted

HTTP Protocol

- There are multiple versions of HTTP
- A HTTP 1.0/1.1 request is nothing more than a few lines of text. (We will focus on 1.0 b/c it is simpler)
- HTTP 2.0 is a binary format (We will not touch on this)
- For each HTTP request, there is a single HTTP response
- The structure of an HTTP 1.0 request and response is defined by [RFC 1945](https://tools.ietf.org/html/rfc1945)
- [Wikipedia](https://en.wikipedia.org/wiki/Hypertext_Transfer_Protocol) has a good example
- Each line is seperated by a Carraige Return byte (CR = `\r`) followed by a Line Feed byte (LF = `\n`)

## Request Structure

| Line            | Example                 | Format                       |
|-----------------|-------------------------|------------------------------|
| Request Line    | `GET /users HTTP/1.1`   | `<method> <path> <protocol>` |
| Header(s)       | `Host: www.example.com` | `<key>: <value>`             |
| Empty Line      | ``                      | N/A                          |
| Body (Optional) | `Hello from client!`    | N/A                          |

We can parse this structure using the following logic:

With a line determined by: everything read up to the CRLF (`\r\n`)...

1. Split the first line of the request on spaces:

```
split[0] = Method
split[1] = Path
split[2] = Protocol
```

2. For each line following (that isnt empty), split on the first `:` (max length of 2):

```
split[0] = Header key
split[1] = Header value
```

Note: HTTP allows for duplicates of the same header keys. This is why Go stores headers as: `map[string][]string` rather than `map[string]string`.

3. When a empty line is encountered, we know we have reached the end of the headers.

4. Body: TODO

## Response Structure

| Line            | Example                               | Format                                   |
|-----------------|---------------------------------------|------------------------------------------|
| Response Line   | `HTTP/1.1 200 OK`                     | `<protocol> <status-code> <status-text>` |
| Header(s)       | `Date: Mon, 23 May 2005 22:38:34 GMT` | `<key>: <value>`                         |
| Empty Line      | ``                                    | N/A                                      |
| Body (Optional) | `Hello from server!`                  | N/A                                      |

## Code

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

