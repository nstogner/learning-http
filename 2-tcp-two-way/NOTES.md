# Two-way TCP

- A TCP connection is bidirectional
- In Go, this translates to our `net.Conn` having `Read` and `Write` methods. (Satisfying the `io.Reader` and `io.Writer` interfaces)

