package scanner

import (
	"fmt"
	"net"
	"time"
)

// Port represents an open port found during a scan.
type Port struct {
	Protocol string
	Number   int
	Address  string
}

// Scanner scans for open TCP/UDP ports on a given host.
type Scanner struct {
	Host    string
	Ports   []int
	Timeout time.Duration
}

// New creates a new Scanner with default settings.
func New(host string, ports []int, timeout time.Duration) *Scanner {
	return &Scanner{
		Host:    host,
		Ports:   ports,
		Timeout: timeout,
	}
}

// Scan checks all configured ports and returns the ones that are open.
func (s *Scanner) Scan() ([]Port, error) {
	var open []Port

	for _, p := range s.Ports {
		addr := fmt.Sprintf("%s:%d", s.Host, p)
		conn, err := net.DialTimeout("tcp", addr, s.Timeout)
		if err != nil {
			// port is closed or unreachable
			continue
		}
		conn.Close()
		open = append(open, Port{
			Protocol: "tcp",
			Number:   p,
			Address:  s.Host,
		})
	}

	return open, nil
}
