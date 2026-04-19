package scanner

import (
	"net"
	"strconv"
	"testing"
	"time"
)

// startListener opens a TCP listener on a random port and returns the port number and a stop func.
func startListener(t *testing.T) (int, func()) {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to start listener: %v", err)
	}
	_, portStr, _ := net.SplitHostPort(ln.Addr().String())
	port, _ := strconv.Atoi(portStr)
	return port, func() { ln.Close() }
}

func TestScan_DetectsOpenPort(t *testing.T) {
	port, stop := startListener(t)
	defer stop()

	s := New("127.0.0.1", []int{port}, time.Second)
	results, err := s.Scan()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 open port, got %d", len(results))
	}
	if results[0].Number != port {
		t.Errorf("expected port %d, got %d", port, results[0].Number)
	}
	if results[0].Protocol != "tcp" {
		t.Errorf("expected protocol tcp, got %s", results[0].Protocol)
	}
}

func TestScan_ClosedPortNotReported(t *testing.T) {
	// Use a port that is almost certainly closed.
	s := New("127.0.0.1", []int{1}, 200*time.Millisecond)
	results, err := s.Scan()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 0 {
		t.Errorf("expected 0 open ports, got %d", len(results))
	}
}
