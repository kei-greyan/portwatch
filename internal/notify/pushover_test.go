package notify

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/user/portwatch/internal/alert"
)

func pushoverAlert(level alert.Level) alert.Alert {
	return alert.Alert{
		Level:   level,
		Message: "port 8080 opened",
		At:      time.Now(),
	}
}

func TestPushover_SendsCorrectPayload(t *testing.T) {
	var got pushoverPayload
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&got); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	n := NewPushover(PushoverConfig{
		Token:   "tok123",
		UserKey: "user456",
		APIURL:  ts.URL,
	})

	a := pushoverAlert(alert.LevelWarn)
	if err := n.Send(a); err != nil {
		t.Fatalf("Send: %v", err)
	}

	if got.Token != "tok123" {
		t.Errorf("token = %q, want tok123", got.Token)
	}
	if got.User != "user456" {
		t.Errorf("user = %q, want user456", got.User)
	}
	if got.Message != a.Message {
		t.Errorf("message = %q, want %q", got.Message, a.Message)
	}
	if got.Priority != pushoverPriorityHigh {
		t.Errorf("priority = %d, want %d", got.Priority, pushoverPriorityHigh)
	}
}

func TestPushover_InfoLevelUsesNormalPriority(t *testing.T) {
	var got pushoverPayload
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&got) //nolint:errcheck
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	n := NewPushover(PushoverConfig{APIURL: ts.URL})
	if err := n.Send(pushoverAlert(alert.LevelInfo)); err != nil {
		t.Fatalf("Send: %v", err)
	}
	if got.Priority != pushoverPriorityNormal {
		t.Errorf("priority = %d, want %d", got.Priority, pushoverPriorityNormal)
	}
}

func TestPushover_ErrorOnBadStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer ts.Close()

	n := NewPushover(PushoverConfig{APIURL: ts.URL})
	if err := n.Send(pushoverAlert(alert.LevelWarn)); err == nil {
		t.Fatal("expected error on non-200 status")
	}
}

func TestPushover_DefaultTimeout(t *testing.T) {
	n := NewPushover(PushoverConfig{Token: "t", UserKey: "u"}).(*pushoverNotifier)
	if n.cfg.Client.Timeout != pushoverDefaultTimeout {
		t.Errorf("timeout = %v, want %v", n.cfg.Client.Timeout, pushoverDefaultTimeout)
	}
}
