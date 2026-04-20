package notify

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/user/portwatch/internal/alert"
)

func teamsAlert(level alert.Level) alert.Alert {
	return alert.Alert{
		Port:      8080,
		Message:   "port 8080 opened",
		Level:     level,
		Timestamp: time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC),
	}
}

func TestTeams_SendsCorrectPayload(t *testing.T) {
	var received teamsPayload

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if ct := r.Header.Get("Content-Type"); ct != "application/json" {
			t.Errorf("expected application/json, got %s", ct)
		}
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	n := NewTeams(srv.URL)
	if err := n.Send(teamsAlert(alert.Warn)); err != nil {
		t.Fatalf("Send returned error: %v", err)
	}

	if received.Type != "MessageCard" {
		t.Errorf("expected @type MessageCard, got %q", received.Type)
	}
	if received.ThemeColor != "FFA500" {
		t.Errorf("expected orange theme for Warn, got %q", received.ThemeColor)
	}
	if len(received.Sections) == 0 {
		t.Fatal("expected at least one section")
	}
	facts := received.Sections[0].Facts
	if len(facts) != 3 {
		t.Fatalf("expected 3 facts, got %d", len(facts))
	}
	if facts[0].Value != "8080" {
		t.Errorf("expected port fact 8080, got %q", facts[0].Value)
	}
}

func TestTeams_InfoLevelUsesGreenColor(t *testing.T) {
	var received teamsPayload

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewDecoder(r.Body).Decode(&received)
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	n := NewTeams(srv.URL)
	if err := n.Send(teamsAlert(alert.Info)); err != nil {
		t.Fatalf("Send returned error: %v", err)
	}
	if received.ThemeColor != "2DC72D" {
		t.Errorf("expected green theme for Info, got %q", received.ThemeColor)
	}
}

func TestTeams_ErrorOnBadStatus(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	n := NewTeams(srv.URL)
	if err := n.Send(teamsAlert(alert.Warn)); err == nil {
		t.Fatal("expected error for 500 status, got nil")
	}
}

func TestTeams_DefaultTimeout(t *testing.T) {
	n := NewTeams("https://example.com/webhook")
	if n.client.Timeout != 10*time.Second {
		t.Errorf("expected 10s timeout, got %v", n.client.Timeout)
	}
}
