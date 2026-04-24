package notify

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/user/portwatch/internal/alert"
)

func larkAlert() alert.Alert {
	return alert.Alert{
		Level:   alert.Warn,
		Message: "new port detected",
		Port:    8080,
		Proto:   "tcp",
		SentAt:  time.Now(),
	}
}

func TestLark_SendsCorrectPayload(t *testing.T) {
	var received larkPayload
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	n := NewLark(srv.URL)
	if err := n.Send(larkAlert()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if received.MsgType != "text" {
		t.Errorf("msg_type = %q, want \"text\"", received.MsgType)
	}
	if received.Content.Text == "" {
		t.Error("content.text should not be empty")
	}
}

func TestLark_ErrorOnBadStatus(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	n := NewLark(srv.URL)
	if err := n.Send(larkAlert()); err == nil {
		t.Error("expected error on 500 response")
	}
}

func TestLark_DefaultTimeout(t *testing.T) {
	n := NewLark("http://example.com").(*larkNotifier)
	if n.client.Timeout != larkDefaultTimeout {
		t.Errorf("timeout = %v, want %v", n.client.Timeout, larkDefaultTimeout)
	}
}
