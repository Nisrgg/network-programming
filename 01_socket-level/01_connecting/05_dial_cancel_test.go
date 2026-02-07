package connececting_test

import (
	"context"
	"net"
	"syscall"
	"testing"
	"time"
)

func TestDialContextCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background()) // 1
	sync := make(chan struct{})
	go func() { // 2
		defer func() { sync <- struct{}{} }()
		var d net.Dialer
		d.Control = func(_, _ string, _ syscall.RawConn) error {
			time.Sleep(time.Second)
			return nil
		}
		conn, err := d.DialContext(ctx, "tcp", "10.0.0.1:80")
		if err != nil {
			t.Log(err)
			return
		}
		conn.Close()
		t.Error("connection did not time out")
	}()
	cancel() // 3
	<-sync
	if ctx.Err() != context.Canceled { // 4
		t.Errorf("expected canceled context; actual: %q", ctx.Err())
	}
}
