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

func lineAlert() alert.Alert {
	return alert.Alert{
		Level:     alert.LevelWarn,
		Message:   "port opened",
		Port:      8080,
		Proto:     "tcp",
		Timestamp: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
	}
}

func TestLine_SendsCorrectPayload(t *testing.T) {
	var gotAuth, gotContentType string
	var body map[string]string

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")
		gotContentType = r.Header.Get("Content-Type")
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	n := notify.NewLine("test-token", srv.URL)
	if err := n.Send(lineAlert()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if gotAuth != "Bearer test-token" {
		t.Errorf("auth header = %q, want %q", gotAuth, "Bearer test-token")
	}
	if gotContentType != "application/json" {
		t.Errorf("content-type = %q, want application/json", gotContentType)
	}
	if body["message"] == "" {
		t.Error("expected non-empty message in payload")
	}
}

func TestLine_ErrorOnBadStatus(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer srv.Close()

	n := notify.NewLine("bad-token", srv.URL)
	if err := n.Send(lineAlert()); err == nil {
		t.Fatal("expected error on 401, got nil")
	}
}

func TestLine_DefaultTimeout(t *testing.T) {
	n := notify.NewLine("tok", "")
	if n == nil {
		t.Fatal("expected non-nil notifier")
	}
}

func TestLine_SetsTimestampWhenZero(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	a := lineAlert()
	a.Timestamp = time.Time{}

	n := notify.NewLine("tok", srv.URL)
	if err := n.Send(a); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
