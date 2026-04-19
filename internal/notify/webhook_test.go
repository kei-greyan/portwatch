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

func TestWebhook_SendsJSON(t *testing.T) {
	var received map[string]interface{}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if ct := r.Header.Get("Content-Type"); ct != "application/json" {
			t.Errorf("unexpected content-type: %s", ct)
		}
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	wh := notify.NewWebhook(ts.URL, 5*time.Second)
	a := alert.Alert{
		Level:     "warn",
		Message:   "port opened",
		Port:      8080,
		Proto:     "tcp",
		Timestamp: time.Now(),
	}
	if err := wh.Send(context.Background(), a); err != nil {
		t.Fatalf("Send: %v", err)
	}
	if received["level"] != "warn" {
		t.Errorf("expected level=warn, got %v", received["level"])
	}
	if received["proto"] != "tcp" {
		t.Errorf("expected proto=tcp, got %v", received["proto"])
	}
}

func TestWebhook_ErrorOnBadStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	wh := notify.NewWebhook(ts.URL, 5*time.Second)
	err := wh.Send(context.Background(), alert.Alert{Timestamp: time.Now()})
	if err == nil {
		t.Fatal("expected error on 500 response")
	}
}

func TestWebhook_DefaultTimeout(t *testing.T) {
	wh := notify.NewWebhook("http://localhost:0", 0)
	if wh == nil {
		t.Fatal("expected non-nil notifier")
	}
}
