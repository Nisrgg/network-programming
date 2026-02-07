package connececting_test

import (
	"net"
	"testing"
)

func TestListener(t *testing.T) {
	listener, err := net.Listen("tcp", ":") // tcp4 for all ipv4 and tcp6 for ipv6 address
	if err != nil {                         // leaving ip address or port empty will make it listen to all
		t.Fatal(err) // unicast and anycast IP address on a random port number
	}
	defer func() { _ = listener.Close() }()
	t.Logf("bound to %q", listener.Addr())
}
