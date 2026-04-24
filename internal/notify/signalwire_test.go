package notify

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/user/portwatch/internal/alert"
)

func signalWireAlert() alert.Alert {
	return alert.Alert{
		Level:     alert.Warn,
		Message:   "port opened",
		Port:      2222,
		Proto:     "tcp",
		Timestamp: time.Date(2024, 1, 2, 3, 4, 5, 0, time.UTC),
	}
}

func TestSignalWire_SendsCorrectPayload(t *testing.T) {
	var got signalWirePayload
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&got); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		w.WriteHeader(http.StatusCreated)
	}))
	defer ts.Close()

	n := &signalWireNotifier{
		spaceURL:  ts.Listener.Addr().String(),
		projectID: "proj123",
		apiToken:  "token456",
		from:      "+15550001111",
		to:        "+15559998888",
		client:    ts.Client(),
	}
	// Override URL to hit test server directly.
	n.spaceURL = ts.Listener.Addr().String()

	a := signalWireAlert()
	if err := n.Send(a); err != nil {
		t.Fatalf("Send: %v", err)
	}

	if got.To != "+15559998888" {
		t.Errorf("To = %q, want +15559998888", got.To)
	}
	if got.From != "+15550001111" {
		t.Errorf("From = %q, want +15550001111", got.From)
	}
	if got.Body == "" {
		t.Error("Body is empty")
	}
}

func TestSignalWire_ErrorOnBadStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer ts.Close()

	n := &signalWireNotifier{
		spaceURL:  ts.Listener.Addr().String(),
		projectID: "proj",
		apiToken:  "tok",
		from:      "+1",
		to:        "+2",
		client:    ts.Client(),
	}

	if err := n.Send(signalWireAlert()); err == nil {
		t.Fatal("expected error on 403, got nil")
	}
}

func TestSignalWire_DefaultTimeout(t *testing.T) {
	n := NewSignalWire("space.signalwire.com", "proj", "tok", "+1", "+2")
	sw, ok := n.(*signalWireNotifier)
	if !ok {
		t.Fatal("unexpected type")
	}
	if sw.client.Timeout != signalWireDefaultTimeout {
		t.Errorf("Timeout = %v, want %v", sw.client.Timeout, signalWireDefaultTimeout)
	}
}
