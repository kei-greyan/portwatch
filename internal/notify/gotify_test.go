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

func gotifyAlert(level alert.Level) alert.Alert {
	return alert.Alert{
		Level:     level,
		Message:   "port 8080 opened",
		Timestamp: time.Now(),
	}
}

func TestGotify_SendsCorrectPayload(t *testing.T) {
	var received map[string]interface{}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		if r.URL.Query().Get("token") != "testtoken" {
			t.Errorf("expected token 'testtoken', got %q", r.URL.Query().Get("token"))
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	g := notify.NewGotify(srv.URL, "testtoken",
		notify.WithGotifyHTTPClient(srv.Client()))

	if err := g.Send(gotifyAlert(alert.Warn)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if received["message"] != "port 8080 opened" {
		t.Errorf("unexpected message: %v", received["message"])
	}
	if priority, ok := received["priority"].(float64); !ok || priority != 8 {
		t.Errorf("expected priority 8 for Warn, got %v", received["priority"])
	}
}

func TestGotify_InfoLevelUsesLowerPriority(t *testing.T) {
	var received map[string]interface{}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewDecoder(r.Body).Decode(&received)
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	g := notify.NewGotify(srv.URL, "tok",
		notify.WithGotifyHTTPClient(srv.Client()))

	if err := g.Send(gotifyAlert(alert.Info)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if priority, ok := received["priority"].(float64); !ok || priority != 5 {
		t.Errorf("expected priority 5 for Info, got %v", received["priority"])
	}
}

func TestGotify_ErrorOnBadStatus(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer srv.Close()

	g := notify.NewGotify(srv.URL, "bad",
		notify.WithGotifyHTTPClient(srv.Client()))

	if err := g.Send(gotifyAlert(alert.Warn)); err == nil {
		t.Fatal("expected error on 401, got nil")
	}
}

func TestGotify_DefaultTimeout(t *testing.T) {
	g := notify.NewGotify("http://localhost", "tok")
	if g == nil {
		t.Fatal("expected non-nil notifier")
	}
}
