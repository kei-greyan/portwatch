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

func googleChatAlert() alert.Alert {
	return alert.Alert{
		Level:   alert.Warn,
		Message: "port opened",
		Port:    8080,
		Proto:   "tcp",
		Time:    time.Now(),
	}
}

func TestGoogleChat_SendsCorrectPayload(t *testing.T) {
	var got map[string]string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&got); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	n := notify.NewGoogleChat(srv.URL)
	if err := n.Send(googleChatAlert()); err != nil {
		t.Fatalf("Send: %v", err)
	}

	if got["text"] == "" {
		t.Error("expected non-empty text field")
	}
	for _, substr := range []string{"8080", "tcp", "port opened"} {
		if !containsStr(got["text"], substr) {
			t.Errorf("text %q missing %q", got["text"], substr)
		}
	}
}

func TestGoogleChat_InfoLevelUsesGreenEmoji(t *testing.T) {
	var got map[string]string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewDecoder(r.Body).Decode(&got)
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	a := googleChatAlert()
	a.Level = alert.Info
	n := notify.NewGoogleChat(srv.URL)
	_ = n.Send(a)

	if !containsStr(got["text"], "🟢") {
		t.Errorf("expected green emoji for info level, got %q", got["text"])
	}
}

func TestGoogleChat_ErrorOnBadStatus(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	n := notify.NewGoogleChat(srv.URL)
	if err := n.Send(googleChatAlert()); err == nil {
		t.Error("expected error on 500 response")
	}
}

func TestGoogleChat_DefaultTimeout(t *testing.T) {
	n := notify.NewGoogleChat("https://example.com")
	if n == nil {
		t.Fatal("expected non-nil notifier")
	}
}

// containsStr is a helper shared across notify tests.
func containsStr(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr ||
		len(s) > 0 && containsSubstr(s, substr))
}

func containsSubstr(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
