package notify

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/yourorg/portwatch/internal/alert"
)

func batchAlerts() []alert.Alert {
	return []alert.Alert{
		{Level: alert.Warn, Message: "port 8080 opened", Port: 8080, Proto: "tcp"},
		{Level: alert.Info, Message: "port 9090 closed", Port: 9090, Proto: "tcp"},
	}
}

func TestWebhookBatch_SendsJSONArray(t *testing.T) {
	var received []alert.Alert
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Errorf("decode body: %v", err)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	nb := NewWebhookBatch(srv.URL, 5*time.Second)
	if err := nb.SendBatch(batchAlerts()); err != nil {
		t.Fatalf("SendBatch: %v", err)
	}
	if len(received) != 2 {
		t.Fatalf("expected 2 alerts, got %d", len(received))
	}
	if received[0].Port != 8080 {
		t.Errorf("expected port 8080, got %d", received[0].Port)
	}
}

func TestWebhookBatch_SetsTimestampWhenZero(t *testing.T) {
	var received []alert.Alert
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&received) //nolint:errcheck
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	alerts := []alert.Alert{{Level: alert.Warn, Message: "test", Port: 22, Proto: "tcp"}}
	nb := NewWebhookBatch(srv.URL, 5*time.Second)
	_ = nb.SendBatch(alerts)
	if len(received) == 0 || received[0].Timestamp.IsZero() {
		t.Error("expected timestamp to be set")
	}
}

func TestWebhookBatch_ErrorOnBadStatus(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	nb := NewWebhookBatch(srv.URL, 5*time.Second)
	err := nb.SendBatch(batchAlerts())
	if err == nil {
		t.Fatal("expected error on 500 status")
	}
}

func TestWebhookBatch_EmptyAlertsIsNoop(t *testing.T) {
	called := false
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	nb := NewWebhookBatch(srv.URL, 5*time.Second)
	if err := nb.SendBatch(nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if called {
		t.Error("expected no HTTP call for empty alert slice")
	}
}

func TestWebhookBatch_DefaultTimeout(t *testing.T) {
	nb := NewWebhookBatch("http://example.com", 0)
	if nb.client.Timeout != 10*time.Second {
		t.Errorf("expected 10s default timeout, got %v", nb.client.Timeout)
	}
}
