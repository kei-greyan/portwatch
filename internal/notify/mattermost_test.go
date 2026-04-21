package notify_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/notify"
)

func mattermostAlert() alert.Alert {
	return alert.Alert{
		Level:   alert.LevelWarn,
		Message: "new port detected",
		Port:    8080,
		Proto:   "tcp",
	}
}

func TestMattermost_SendsCorrectPayload(t *testing.T) {
	var received map[string]interface{}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Errorf("decode body: %v", err)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	n := notify.NewMattermost(notify.MattermostConfig{
		WebhookURL: srv.URL,
		Channel:    "#alerts",
		Username:   "portwatch",
	})

	if err := n.Send(mattermostAlert()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if received["channel"] != "#alerts" {
		t.Errorf("channel = %v, want #alerts", received["channel"])
	}
	if received["username"] != "portwatch" {
		t.Errorf("username = %v, want portwatch", received["username"])
	}
	text, _ := received["text"].(string)
	if text == "" {
		t.Error("text field is empty")
	}
}

func TestMattermost_ErrorOnBadStatus(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	n := notify.NewMattermost(notify.MattermostConfig{WebhookURL: srv.URL})
	if err := n.Send(mattermostAlert()); err == nil {
		t.Fatal("expected error for non-2xx status")
	}
}

func TestMattermost_DefaultTimeout(t *testing.T) {
	// Ensure zero Timeout is replaced with a sensible default (no panic).
	n := notify.NewMattermost(notify.MattermostConfig{
		WebhookURL: "http://127.0.0.1:0",
		Timeout:    0,
	})
	// Send will fail due to bad URL; we only care that it doesn't hang.
	done := make(chan struct{})
	go func() {
		_ = n.Send(mattermostAlert())
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(15 * time.Second):
		t.Fatal("Send blocked longer than expected")
	}
}

func TestMattermost_TextContainsPortAndProto(t *testing.T) {
	var received map[string]interface{}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Errorf("decode body: %v", err)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	n := notify.NewMattermost(notify.MattermostConfig{WebhookURL: srv.URL})
	if err := n.Send(mattermostAlert()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	text, _ := received["text"].(string)
	for _, want := range []string{"8080", "tcp"} {
		if !contains(text, want) {
			t.Errorf("text %q does not contain %q", text, want)
		}
	}
}

// contains reports whether substr appears within s.
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr ||
		len(substr) == 0 ||
		(len(s) > 0 && stringContains(s, substr)))
}

func stringContains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
