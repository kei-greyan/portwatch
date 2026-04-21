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

func splunkAlert() alert.Alert {
	return alert.Alert{
		Level:     alert.LevelWarn,
		Message:   "port 9200 opened",
		Port:      9200,
		Timestamp: time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC),
	}
}

func TestSplunk_SendsCorrectPayload(t *testing.T) {
	var captured map[string]any
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&captured); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		auth := r.Header.Get("Authorization")
		if auth != "Splunk test-token" {
			t.Errorf("unexpected Authorization header: %q", auth)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	n := notify.NewSplunk(ts.URL, "test-token")
	if err := n.Send(splunkAlert()); err != nil {
		t.Fatalf("Send: %v", err)
	}

	event, ok := captured["event"].(map[string]any)
	if !ok {
		t.Fatalf("event field missing or wrong type")
	}
	if event["message"] != "port 9200 opened" {
		t.Errorf("unexpected message: %v", event["message"])
	}
	if captured["sourcetype"] != "portwatch" {
		t.Errorf("unexpected sourcetype: %v", captured["sourcetype"])
	}
}

func TestSplunk_ErrorOnBadStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer ts.Close()

	n := notify.NewSplunk(ts.URL, "bad-token")
	if err := n.Send(splunkAlert()); err == nil {
		t.Fatal("expected error on 403, got nil")
	}
}

func TestSplunk_DefaultTimeout(t *testing.T) {
	n := notify.NewSplunk("http://localhost:19999", "tok")
	a := splunkAlert()
	// Should return an error (connection refused) rather than hang.
	if err := n.Send(a); err == nil {
		t.Fatal("expected connection error, got nil")
	}
}

func TestSplunk_SetsTimestampWhenZero(t *testing.T) {
	var captured map[string]any
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewDecoder(r.Body).Decode(&captured)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	a := splunkAlert()
	a.Timestamp = time.Time{}
	n := notify.NewSplunk(ts.URL, "tok")
	if err := n.Send(a); err != nil {
		t.Fatalf("Send: %v", err)
	}
	if captured["time"] == float64(0) {
		t.Error("expected non-zero timestamp")
	}
}
