# 3. Minimal HTTP Server

We will define minimal here as a server that can successfully respond to a simple curl GET request. We will start out by observing what happens when we issue this request to our original server that logs all of the lines it receives over TCP:

```sh
go run ./1-read-tcp/main.go
```

```sh
# In another shell
curl localhost:7000/abc
```

