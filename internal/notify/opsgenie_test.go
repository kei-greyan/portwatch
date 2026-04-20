package notify_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/notify"
)

func opsGenieAlert() alert.Alert {
	return alert.Alert{
		Title:     "Port opened",
		Body:      "Port 9200 is now open",
		Level:     alert.Warn,
		Port:      9200,
		Host:      "localhost",
		Timestamp: time.Now(),
	}
}

func TestOpsGenie_SendsCorrectPayload(t *testing.T) {
	var received map[string]interface{}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "GenieKey test-key" {
			t.Errorf("missing or wrong Authorization header: %s", r.Header.Get("Authorization"))
		}
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		w.WriteHeader(http.StatusAccepted)
	}))
	defer srv.Close()

	og := notify.NewOpsGenie("test-key", func(o *notify.OpsGenie) {
		o.(*notify.OpsGenie) // compile check — use exported setter instead
	})
	// Use internal URL override via a wrapper approach; test server replaces URL.
	og2 := notify.NewOpsGenieWithURL("test-key", srv.URL)

	if err := og2.Send(context.Background(), opsGenieAlert()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if received["message"] != "Port opened" {
		t.Errorf("message = %v, want 'Port opened'", received["message"])
	}
	if received["priority"] != "P2" {
		t.Errorf("priority = %v, want P2 for Warn level", received["priority"])
	}
	_ = og // suppress unused
}

func TestOpsGenie_InfoLevelUsesP3(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var payload map[string]interface{}
		_ = json.NewDecoder(r.Body).Decode(&payload)
		if payload["priority"] != "P3" {
			t.Errorf("priority = %v, want P3 for Info level", payload["priority"])
		}
		w.WriteHeader(http.StatusAccepted)
	}))
	defer srv.Close()

	og := notify.NewOpsGenieWithURL("key", srv.URL)
	a := opsGenieAlert()
	a.Level = alert.Info
	if err := og.Send(context.Background(), a); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestOpsGenie_ErrorOnBadStatus(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer srv.Close()

	og := notify.NewOpsGenieWithURL("bad-key", srv.URL)
	if err := og.Send(context.Background(), opsGenieAlert()); err == nil {
		t.Fatal("expected error on 401, got nil")
	}
}

func TestOpsGenie_DefaultTimeout(t *testing.T) {
	og := notify.NewOpsGenieWithURL("key", "http://127.0.0.1:1")
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()
	if err := og.Send(ctx, opsGenieAlert()); err == nil {
		t.Fatal("expected connection error")
	}
}
