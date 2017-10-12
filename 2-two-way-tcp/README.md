# 2. Writing to TCP Connections

In this example, we define a simple protocol for talking to our server.

The client can issue a command by sending a command string plus a CRLF ("\r\n") over the connection. By default we have implemented a single command: "BEEP". Ideally this will cause our computer to beep if all goes well.

```sh
go run ./2-two-way-tcp/main.go
```

```sh
# In another shell
telnet localhost 7000
BEEP<enter>
```

The server will ether respond with "REJECTED" or "ACCEPTED" depending on whether it understands the command.

## Bonus

- Add a "CLOSE" command that will close the connection to the client

