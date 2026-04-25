package notify

import (
	"fmt"
	"net"
	"time"

	"github.com/user/portwatch/internal/alert"
)

// XMPPClient is a minimal interface for sending XMPP messages, allowing test injection.
type XMPPClient interface {
	Send(to, body string) error
	Close() error
}

// xmppNotifier sends alerts via XMPP (Jabber).
type xmppNotifier struct {
	host   string
	port   int
	from   string
	to     string
	dial   func(host string, port int, from, password string) (XMPPClient, error)
	passwd string
}

// NewXMPP returns a Notifier that delivers alerts over XMPP.
func NewXMPP(host string, port int, from, password, to string) Notifier {
	return &xmppNotifier{
		host:   host,
		port:   port,
		from:   from,
		passwd: password,
		to:     to,
		dial:   defaultXMPPDial,
	}
}

func (x *xmppNotifier) Send(a alert.Alert) error {
	if a.Timestamp.IsZero() {
		a.Timestamp = time.Now().UTC()
	}
	body := fmt.Sprintf("[portwatch] %s | port %d/%s | %s",
		a.Timestamp.Format(time.RFC3339), a.Port, a.Proto, a.Message)

	client, err := x.dial(x.host, x.port, x.from, x.passwd)
	if err != nil {
		return fmt.Errorf("xmpp: dial: %w", err)
	}
	defer client.Close()

	if err := client.Send(x.to, body); err != nil {
		return fmt.Errorf("xmpp: send: %w", err)
	}
	return nil
}

// defaultXMPPDial is a placeholder that verifies TCP reachability only.
// Real deployments should replace this with a proper XMPP library.
func defaultXMPPDial(host string, port int, _, _ string) (XMPPClient, error) {
	addr := fmt.Sprintf("%s:%d", host, port)
	conn, err := net.DialTimeout("tcp", addr, 5*time.Second)
	if err != nil {
		return nil, err
	}
	return &tcpXMPPStub{conn: conn}, nil
}

type tcpXMPPStub struct{ conn net.Conn }

func (s *tcpXMPPStub) Send(_, _ string) error { return nil }
func (s *tcpXMPPStub) Close() error           { return s.conn.Close() }
