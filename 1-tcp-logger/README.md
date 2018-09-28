# 1. Reading from a TCP connection

## Networking Layers

- OSI Model (7-layer) (Open Systems Interconnection)
- TCP/IP Model (4-layer)
  - Also called: Internet protocol suite
  - Also called: DoD (Department of Defense) Model (Development was funded by US government: DARPA)

```
| OSI Model             | TCP/IP Model            | Protocol |
|-----------------------|-------------------------|----------|
| 7. Application Layer  |                         | HTTP     |
| 6. Presentation Layer | 4. Application Layer    |          |
| 5. Session Layer      |                         |          |
|-----------------------|-------------------------|----------|
| 4. Transport Layer    | 3. Transport Layer      | TCP      |
|-----------------------|-------------------------|----------|
| 3. Network Layer      | 2. Internet Layer       |          |
|-----------------------|-------------------------|----------|
| 2. Data Link Layer    |                         |          |
| 1. Physical Layer     | 1. Network Access Layer |          |
|-----------------------|-------------------------|----------|
```

NOTES:

- Usually when people refer to a layer by number, they are using the OSI model
- Usually TCP/IP layers are referred to by name
- Some protocols dont neatly fit into a given layer (i.e. TCP in layer 4, but it also deals with layer 5 "Session Layer")

**We are going to build out a library that implements the HTTP protocol using the abstractions provided by TCP.**

## What is TCP?

- Stands for Transmission Control Protocol
- Has the concept of a connection between a client and a server
- Handles control flow (it tries not to send more data than the receiver can handle)
- Handles congestion control (using a slow-start algorithm)
- Allows for the reliable transfer of data (bytes of information) across the network
- It protects against data corruption by using checksums of the data that gets sent

Side Note: How does it relate to UDP?

- UDP stands for User Datagram Protocol
- UDP also checksums data to provide data integrity
- UDP does not redeliver dropped packets

## Tools

What is `telnet`?

- When referring to `telnet` here, we are talking about the CLI tool
- There is a `telnet` text-based protocol
- We will use `telnet` to make TCP connections and send text over them 
- When we press <enter> `telnet` will send what we typed followed by two bytes: `\r\n`

## Code

This simple server listens for TCP connections and logs anything that is sent to it before closing the connection to the client.

```sh
go run ./1-tcp-logger/main.go
```

```sh
# In another shell
telnet localhost 7000
Hey there server!<enter>
```

Note what is logged in the server. The `telnet` tool sends a carriage return `\r` and a newline feed `\n` (usually referred to as CRLF online).
