package connececting_test

import (
	"io"
	"net"
	"testing"
)

func TestDial(t *testing.T) {
	// Create a listener on a random port.
	listener, err := net.Listen("tcp", "127.0.0.1:")
	if err != nil {
		t.Fatal(err)
	}
	done := make(chan struct{})

	go func() { // 1
		defer func() { done <- struct{}{} }()

		for {
			conn, err := listener.Accept() // 2
			if err != nil {
				t.Log(err)
				return
			}
			go func(c net.Conn) { // 3
				defer func() {
					c.Close()
					done <- struct{}{}
				}()

				buf := make([]byte, 1024)
				for {
					n, err := c.Read(buf) // 4
					if err != nil {
						if err != io.EOF {
							t.Error(err)
						}
						return
					}
					t.Logf("received: %q", buf[:n])
				}
			}(conn)
		}
	}()
	conn, err := net.Dial("tcp", listener.Addr().String()) // 5, 6, 7
	if err != nil {
		t.Fatal(err)
	}
	conn.Close() // 8
	<-done
	listener.Close() // 9
	<-done
}
