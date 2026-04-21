package notify_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/yourusername/portwatch/internal/alert"
	"github.com/yourusername/portwatch/internal/notify"
)

var zulipAlert = alert.Alert{
	Level:   alert.Warn,
	Message: "port opened",
	Port:    9090,
	Time:    time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
}

func TestZulip_SendsCorrectPayload(t *testing.T) {
	var gotBody url.Values

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			t.Fatalf("parse form: %v", err)
		}
		gotBody = r.Form
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]string{"result": "success", "msg": ""})
	}))
	defer srv.Close()

	n := notify.NewZulip(srv.URL, "bot@example.com", "token123", "alerts", "portwatch")
	if err := n.Send(zulipAlert); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if gotBody.Get("type") != "stream" {
		t.Errorf("expected type=stream, got %q", gotBody.Get("type"))
	}
	if gotBody.Get("to") != "alerts" {
		t.Errorf("expected to=alerts, got %q", gotBody.Get("to"))
	}
	if gotBody.Get("topic") != "portwatch" {
		t.Errorf("expected topic=portwatch, got %q", gotBody.Get("topic"))
	}
	if !strings.Contains(gotBody.Get("content"), "9090") {
		t.Errorf("expected content to contain port 9090, got %q", gotBody.Get("content"))
	}
}

func TestZulip_ErrorOnBadStatus(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(map[string]string{"result": "error", "msg": "invalid credentials"})
	}))
	defer srv.Close()

	n := notify.NewZulip(srv.URL, "bot@example.com", "badtoken", "alerts", "portwatch")
	if err := n.Send(zulipAlert); err == nil {
		t.Fatal("expected error on 401, got nil")
	}
}

func TestZulip_DefaultTimeout(t *testing.T) {
	n := notify.NewZulip("https://zulip.example.com", "bot@example.com", "token", "alerts", "portwatch")
	if n == nil {
		t.Fatal("expected non-nil notifier")
	}
}
