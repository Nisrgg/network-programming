package connececting_test

import (
	"context"
	"net"
	"syscall"
	"testing"
	"time"
)

func TestDialContext(t *testing.T) {
	dl := time.Now().Add(5 * time.Second)                         // 1
	ctx, cancel := context.WithDeadline(context.Background(), dl) // 2
	defer cancel()                                                // 3
	var d net.Dialer                                              // DialContext is a method on a Dialer
	d.Control = func(_, _ string, _ syscall.RawConn) error {      // 4
		time.Sleep(5*time.Second + time.Millisecond) // Sleep long enough to reach the
		return nil                                   // context's deadline
	}
	conn, err := d.DialContext(ctx, "tcp", "10.0.0.0:80") // 5
	if err == nil {
		conn.Close()
		t.Fatal("connection did not time out")
	}
	nErr, ok := err.(net.Error)
	if !ok {
		t.Error(err)
	} else {
		if !nErr.Timeout() {
			t.Errorf("error is not a timeout: %v", err)
		}
	}
	if ctx.Err() != context.DeadlineExceeded { // 6
		t.Errorf("expected deadline exceeded; actual: %v", ctx.Err())
	}
}
