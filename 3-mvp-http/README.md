# 3. Minimal HTTP Server

We will define "minimal HTTP server" as a server that can successfully respond to a simple curl GET request. We will start out by observing what happens when we issue this request to our TCP logger server:

```sh
go run ./1-tcp-logger/main.go
```

```sh
# In another shell
curl localhost:7000/abc
```

