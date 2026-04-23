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

func sigSciAlert() alert.Alert {
	return alert.Alert{
		Message: "port 9200 opened",
		Level: "warn",
		Timestamp: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
	}
}

func TestSignalSciences_SendsCorrectPayload(t *testing.T) {
	var got map[string]interface{}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&got); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		if r.Header.Get("x-api-user") != "keyid" {
			t.Errorf("expected x-api-user=keyid, got %s", r.Header.Get("x-api-user"))
		}
		if r.Header.Get("x-api-token") != "secret" {
			t.Errorf("expected x-api-token=secret, got %s", r.Header.Get("x-api-token"))
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	n := notify.NewSignalSciences(ts.URL, "mycorp", "mysite", "keyid", "secret")
	if err := n.Send(sigSciAlert()); err != nil {
		t.Fatalf("Send: %v", err)
	}

	if got["event"] != "portwatch" {
		t.Errorf("event = %v, want portwatch", got["event"])
	}
	if got["message"] != "port 9200 opened" {
		t.Errorf("message = %v, want 'port 9200 opened'", got["message"])
	}
	if got["level"] != "warn" {
		t.Errorf("level = %v, want warn", got["level"])
	}
}

func TestSignalSciences_ErrorOnBadStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer ts.Close()

	n := notify.NewSignalSciences(ts.URL, "corp", "site", "k", "s")
	if err := n.Send(sigSciAlert()); err == nil {
		t.Error("expected error on 403, got nil")
	}
}

func TestSignalSciences_DefaultTimeout(t *testing.T) {
	n := notify.NewSignalSciences("http://localhost", "c", "s", "k", "t")
	if n == nil {
		t.Fatal("expected non-nil notifier")
	}
}

func TestSignalSciences_SetsTimestampWhenZero(t *testing.T) {
	var got map[string]interface{}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewDecoder(r.Body).Decode(&got)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	n := notify.NewSignalSciences(ts.URL, "c", "s", "k", "t")
	a := alert.Alert{Message: "test", Level: "info"} // zero timestamp
	if err := n.Send(a); err != nil {
		t.Fatalf("Send: %v", err)
	}
	if got["timestamp"] == "" || got["timestamp"] == nil {
		t.Error("expected non-empty timestamp")
	}
}
