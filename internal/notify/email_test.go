package notify_test

import (
	"io"
	"net"
	"net/smtp"
	"net/textproto"
	"strings"
	"testing"
	"time"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/notify"
)

// smtpServer is a minimal fake SMTP server for testing.
type smtpServer struct {
	ln      net.Listener
	Received []string
}

func newSMTPServer(t *testing.T) *smtpServer {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	s := &smtpServer{ln: ln}
	go s.serve()
	return s
}

func (s *smtpServer) Addr() string { return s.ln.Addr().String() }

func (s *smtpServer) serve() {
	conn, err := s.ln.Accept()
	if err != nil {
		return
	}
	defer conn.Close()
	tc := textproto.NewConn(conn)
	_ = tc.PrintfLine("220 localhost SMTP")
	for {
		line, err := tc.ReadLine()
		if err == io.EOF {
			break
		}
		s.Received = append(s.Received, line)
		upper := strings.ToUpper(line)
		switch {
		case strings.HasPrefix(upper, "EHLO"), strings.HasPrefix(upper, "HELO"):
			_ = tc.PrintfLine("250 OK")
		case upper == "DATA":
			_ = tc.PrintfLine("354 Start")
		case line == ".":
			_ = tc.PrintfLine("250 OK")
		case upper == "QUIT":
			_ = tc.PrintfLine("221 Bye")
			return
		default:
			_ = tc.PrintfLine("250 OK")
		}
	}
}

func TestEmail_SendsAlert(t *testing.T) {
	_ = smtp.PlainAuth // ensure import used
	ts := newSMTPServer(t)
	defer ts.ln.Close()

	host, portStr, _ := net.SplitHostPort(ts.Addr())
	var port int
	fmt.Sscanf(portStr, "%d", &port)

	cfg := notify.EmailConfig{
		Host: host,
		Port: port,
		From: "portwatch@example.com",
		To:   []string{"admin@example.com"},
	}
	n := notify.NewEmail(cfg)
	a := alert.Alert{
		Level:     alert.Warn,
		Port:      8080,
		Message:   "port opened",
		Timestamp: time.Now(),
	}
	if err := n.Send(a); err != nil {
		t.Fatalf("Send: %v", err)
	}
}
