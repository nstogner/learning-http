# 1. Reading from a TCP connection

This simple server listens for TCP connections and logs anything that is sent to it before closing the connection to the client.

```sh
telnet localhost 7000
Hey there server!<enter>
```

Note what is logged in the server. Telnet sends a carriage return `\r` and a newline feed `\n` (usually referred to as CRLF online).

