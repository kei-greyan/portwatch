package notify

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/user/portwatch/internal/alert"
)

func ntfyAlert() alert.Alert {
	return alert.Alert{
		Level:   alert.LevelWarn,
		Message: "port 8080/tcp opened",
		At:      time.Now(),
	}
}

func TestNtfy_SendsCorrectPayload(t *testing.T) {
	var gotPath, gotPriority, gotTitle, gotBody string

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotPriority = r.Header.Get("Priority")
		gotTitle = r.Header.Get("Title")
		buf := make([]byte, 512)
		n, _ := r.Body.Read(buf)
		gotBody = string(buf[:n])
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	a := ntfyAlert()
	n := NewNtfy(NtfyConfig{
		BaseURL: srv.URL,
		Topic:   "portwatch",
		client:  srv.Client(),
	})

	if err := n.Send(a); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotPath != "/portwatch" {
		t.Errorf("path = %q, want /portwatch", gotPath)
	}
	if gotPriority != "high" {
		t.Errorf("priority = %q, want high", gotPriority)
	}
	if gotTitle == "" {
		t.Error("expected non-empty Title header")
	}
	if gotBody != a.Message {
		t.Errorf("body = %q, want %q", gotBody, a.Message)
	}
}

func TestNtfy_InfoLevelUsesDefaultPriority(t *testing.T) {
	var gotPriority string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPriority = r.Header.Get("Priority")
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	a := ntfyAlert()
	a.Level = alert.LevelInfo
	n := NewNtfy(NtfyConfig{BaseURL: srv.URL, Topic: "portwatch", client: srv.Client()})
	if err := n.Send(a); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotPriority != "default" {
		t.Errorf("priority = %q, want default", gotPriority)
	}
}

func TestNtfy_ErrorOnBadStatus(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer srv.Close()

	n := NewNtfy(NtfyConfig{BaseURL: srv.URL, Topic: "portwatch", client: srv.Client()})
	if err := n.Send(ntfyAlert()); err == nil {
		t.Fatal("expected error on 403, got nil")
	}
}

func TestNtfy_DefaultTimeout(t *testing.T) {
	n := NewNtfy(NtfyConfig{BaseURL: "https://ntfy.sh", Topic: "portwatch"})
	nn := n.(*ntfyNotifier)
	if nn.cfg.client.Timeout != 10*time.Second {
		t.Errorf("timeout = %v, want 10s", nn.cfg.client.Timeout)
	}
}

func TestNtfy_SendsBearerToken(t *testing.T) {
	var gotAuth string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	n := NewNtfy(NtfyConfig{BaseURL: srv.URL, Topic: "portwatch", Token: "secret", client: srv.Client()})
	if err := n.Send(ntfyAlert()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotAuth != "Bearer secret" {
		t.Errorf("Authorization = %q, want 'Bearer secret'", gotAuth)
	}
}
