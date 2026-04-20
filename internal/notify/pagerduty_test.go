package notify

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/user/portwatch/internal/alert"
)

func pdAlert(level alert.Level, msg string) alert.Alert {
	return alert.Alert{
		Level:     level,
		Message:   msg,
		Timestamp: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
	}
}

func TestPagerDuty_SendsCorrectPayload(t *testing.T) {
	var received pagerDutyPayload

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		w.WriteHeader(http.StatusAccepted)
	}))
	defer srv.Close()

	pd := NewPagerDuty("test-routing-key")
	pd.url = srv.URL

	a := pdAlert(alert.Warn, "port 8080 opened")
	if err := pd.Send(a); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if received.RoutingKey != "test-routing-key" {
		t.Errorf("routing key = %q, want %q", received.RoutingKey, "test-routing-key")
	}
	if received.EventAction != "trigger" {
		t.Errorf("event_action = %q, want trigger", received.EventAction)
	}
	if received.Payload.Severity != pagerDutySeverityErr {
		t.Errorf("severity = %q, want %q", received.Payload.Severity, pagerDutySeverityErr)
	}
	if received.Payload.Summary != "port 8080 opened" {
		t.Errorf("summary = %q, want %q", received.Payload.Summary, "port 8080 opened")
	}
	if received.Payload.Source != "portwatch" {
		t.Errorf("source = %q, want portwatch", received.Payload.Source)
	}
}

func TestPagerDuty_InfoLevelUsesSeverityInfo(t *testing.T) {
	var received pagerDutyPayload

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewDecoder(r.Body).Decode(&received)
		w.WriteHeader(http.StatusAccepted)
	}))
	defer srv.Close()

	pd := NewPagerDuty("key")
	pd.url = srv.URL

	if err := pd.Send(pdAlert(alert.Info, "port 22 closed")); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if received.Payload.Severity != pagerDutySeverityInfo {
		t.Errorf("severity = %q, want %q", received.Payload.Severity, pagerDutySeverityInfo)
	}
}

func TestPagerDuty_ErrorOnBadStatus(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer srv.Close()

	pd := NewPagerDuty("bad-key")
	pd.url = srv.URL

	if err := pd.Send(pdAlert(alert.Warn, "test")); err == nil {
		t.Fatal("expected error for non-2xx status, got nil")
	}
}

func TestPagerDuty_DefaultTimeout(t *testing.T) {
	pd := NewPagerDuty("key")
	if pd.client.Timeout != defaultPDTimeout {
		t.Errorf("timeout = %v, want %v", pd.client.Timeout, defaultPDTimeout)
	}
}
