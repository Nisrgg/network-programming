# Network Programming in Go

This repository documents my hands-on journey learning network programming
using Go. The code is written while studying *Network Programming with Go*
by Adam Woodbeck (No Starch Press), with a focus on understanding how
networked systems work at a practical level rather than building production
frameworks.

---

## Requirements

- Go 1.20 or newer
- Linux is required for Unix Domain Socket credential and authentication
  examples
- Some tests assume access to free local ports

---

## Running the Code

Run all tests:

```sh
go test ./...
```

Some networking tests are timing- or environment-sensitive and may behave
differently across systems.