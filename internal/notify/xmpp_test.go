package notify

import (
	"errors"
	"testing"
	"time"

	"github.com/user/portwatch/internal/alert"
)

// fakeXMPPClient captures the last message sent.
type fakeXMPPClient struct {
	to   string
	body string
	err  error
}

func (f *fakeXMPPClient) Send(to, body string) error {
	f.to = to
	f.body = body
	return f.err
}
func (f *fakeXMPPClient) Close() error { return nil }

func xmppAlert() alert.Alert {
	return alert.Alert{
		Port:      5222,
		Proto:     "tcp",
		Message:   "new port opened",
		Level:     alert.Warn,
		Timestamp: time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC),
	}
}

func TestXMPP_SendsCorrectPayload(t *testing.T) {
	client := &fakeXMPPClient{}
	n := &xmppNotifier{
		to: "admin@example.com",
		dial: func(_, _ string, _, _ string) (XMPPClient, error) {
			return client, nil
		},
	}

	if err := n.Send(xmppAlert()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if client.to != "admin@example.com" {
		t.Errorf("recipient = %q, want admin@example.com", client.to)
	}
	for _, want := range []string{"5222", "tcp", "new port opened", "portwatch"} {
		if !contains(client.body, want) {
			t.Errorf("body %q missing %q", client.body, want)
		}
	}
}

func TestXMPP_ErrorOnDialFailure(t *testing.T) {
	n := &xmppNotifier{
		to: "admin@example.com",
		dial: func(_, _ string, _, _ string) (XMPPClient, error) {
			return nil, errors.New("connection refused")
		},
	}
	if err := n.Send(xmppAlert()); err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestXMPP_ErrorOnSendFailure(t *testing.T) {
	client := &fakeXMPPClient{err: errors.New("send failed")}
	n := &xmppNotifier{
		to: "admin@example.com",
		dial: func(_, _ string, _, _ string) (XMPPClient, error) {
			return client, nil
		},
	}
	if err := n.Send(xmppAlert()); err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestXMPP_SetsTimestampWhenZero(t *testing.T) {
	client := &fakeXMPPClient{}
	n := &xmppNotifier{
		to: "admin@example.com",
		dial: func(_, _ string, _, _ string) (XMPPClient, error) {
			return client, nil
		},
	}
	a := xmppAlert()
	a.Timestamp = time.Time{}
	if err := n.Send(a); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if client.body == "" {
		t.Error("body should not be empty")
	}
}

// contains is a small helper shared across notify tests.
func contains(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(s) > 0 && stringContains(s, sub))
}

func stringContains(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
