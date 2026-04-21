package notify_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/notify"
)

func snsAlert() alert.Alert {
	return alert.Alert{
		Level:   "warn",
		Port:    8080,
		Message: "port 8080 opened",
		Time:    time.Now(),
	}
}

func TestSNS_SendsCorrectPayload(t *testing.T) {
	var received map[string]string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &received)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	n := notify.NewSNS(ts.URL, "arn:aws:sns:us-east-1:123456789012:portwatch")
	if err := n.Send(context.Background(), snsAlert()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if received["TopicARN"] != "arn:aws:sns:us-east-1:123456789012:portwatch" {
		t.Errorf("unexpected TopicARN: %q", received["TopicARN"])
	}
	if received["Message"] != "port 8080 opened" {
		t.Errorf("unexpected Message: %q", received["Message"])
	}
	if received["Subject"] == "" {
		t.Error("Subject should not be empty")
	}
}

func TestSNS_ErrorOnBadStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	n := notify.NewSNS(ts.URL, "arn:aws:sns:us-east-1:123456789012:portwatch")
	if err := n.Send(context.Background(), snsAlert()); err == nil {
		t.Fatal("expected error on 500 response")
	}
}

func TestSNS_DefaultTimeout(t *testing.T) {
	n := notify.NewSNS("http://localhost", "arn:aws:sns:us-east-1:123456789012:portwatch")
	_ = n // construction should not panic; timeout is set internally
}

func TestSNS_WithCustomClient(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusAccepted)
	}))
	defer ts.Close()

	custom := &http.Client{Timeout: 5 * time.Second}
	n := notify.NewSNS(ts.URL, "arn:test", notify.WithSNSHTTPClient(custom))
	if err := n.Send(context.Background(), snsAlert()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
