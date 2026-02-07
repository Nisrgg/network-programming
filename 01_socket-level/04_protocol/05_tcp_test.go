package protocol

import (
	"io"
	"net"
	"sync"
	"testing"
	"time"
)

// TestTCPConnFunctionalities demonstrates specific *net.TCPConn methods
// that are not available on the generic net.Conn interface.
func TestTCPConnFunctionalities(t *testing.T) {
	// 1. Setup: Create a TCP Listener
	// We use ListenTCP specifically to get a *net.TCPListener
	addr, _ := net.ResolveTCPAddr("tcp", "127.0.0.1:0") // Port 0 picks a random free port
	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		t.Fatal(err)
	}
	defer listener.Close()

	var wg sync.WaitGroup
	wg.Add(2) // We will wait for both Client and Server goroutines

	// --- SERVER GOROUTINE ---
	go func() {
		defer wg.Done()
		t.Logf("[Server] Listening on %s", listener.Addr())

		// AcceptTCP returns *net.TCPConn, not the generic net.Conn
		conn, err := listener.AcceptTCP()
		if err != nil {
			t.Error(err)
			return
		}
		defer conn.Close() // Ensure full cleanup eventually

		// --- Functionality 1: KeepAlive ---
		// Helps detect if the client has silently disappeared (e.g., power loss).
		// Not available on generic net.Conn.
		err = conn.SetKeepAlive(true)
		if err != nil {
			t.Error("Failed to set KeepAlive:", err)
		}
		conn.SetKeepAlivePeriod(3 * time.Minute)
		t.Log("[Server] Functionality: KeepAlive enabled")

		// --- Functionality 2: SetReadBuffer ---
		// interacting with the OS kernel's socket receive buffer.
		// Useful for tuning high-throughput applications.
		err = conn.SetReadBuffer(4096) // Set to 4KB
		if err != nil {
			t.Error("Failed to set read buffer:", err)
		}
		t.Log("[Server] Functionality: ReadBuffer size set")

		// Read Loop
		buf := make([]byte, 1024)
		for {
			n, err := conn.Read(buf)
			if err == io.EOF {
				t.Log("[Server] Client sent EOF (CloseWrite received)")
				break
			}
			if err != nil {
				t.Error("Read error:", err)
				break
			}
			t.Logf("[Server] Received data: %s", string(buf[:n]))
		}
	}()

	// --- CLIENT GOROUTINE ---
	go func() {
		defer wg.Done()

		// Small sleep to ensure server is listening
		time.Sleep(100 * time.Millisecond)

		serverAddr := listener.Addr().(*net.TCPAddr)

		// DialTCP returns *net.TCPConn
		conn, err := net.DialTCP("tcp", nil, serverAddr)
		if err != nil {
			t.Error(err)
			return
		}
		defer conn.Close()

		t.Log("[Client] Connected to server")

		// --- Functionality 3: SetNoDelay (Nagle's Algorithm) ---
		// true = Disable Nagle. Send packets immediately (low latency).
		// false = Enable Nagle. Buffer small packets (bandwidth efficiency).
		// This is critical for real-time apps (games, SSH).
		err = conn.SetNoDelay(true)
		if err != nil {
			t.Error("Failed to set NoDelay:", err)
		}
		t.Log("[Client] Functionality: Nagle's Algorithm disabled (NoDelay=true)")

		// --- Functionality 4: SetWriteBuffer ---
		// Tuning the OS kernel's socket send buffer.
		err = conn.SetWriteBuffer(4096)
		if err != nil {
			t.Error("Failed to set write buffer:", err)
		}
		t.Log("[Client] Functionality: WriteBuffer size set")

		// Send Data
		conn.Write([]byte("Hello TCP World!"))

		// --- Functionality 5: CloseWrite (Half-Close) ---
		// This is a TCP-specific feature. It says "I am done writing,
		// but I can still read if you have anything to say."
		// It sends a FIN packet to the server, resulting in io.EOF on the server side.
		err = conn.CloseWrite()
		if err != nil {
			t.Error("Failed to CloseWrite:", err)
		}
		t.Log("[Client] Functionality: CloseWrite called (Half-Closed)")

		// Note: We could still conn.Read() here if the server sent data back.
	}()

	// Wait for the test interaction to complete
	wg.Wait()
}
