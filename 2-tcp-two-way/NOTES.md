# Two-way TCP

- A TCP connection is bidirectional
- In Go, this translates to our `net.Conn` interface having both `Read` and `Write` methods. (`io.Reader` / `io.Writer` interfaces)

