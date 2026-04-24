package notify_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/notify"
)

func amplitudeAlert() alert.Alert {
	return alert.Alert{
		Message:   "port 9200 opened",
		Level:     alert.Warn,
		Port:      9200,
		Proto:     "tcp",
		Timestamp: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
	}
}

func TestAmplitude_SendsCorrectPayload(t *testing.T) {
	var received map[string]interface{}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &received)
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	n := notify.NewAmplitude("test-api-key", srv.URL)
	if err := n.Send(amplitudeAlert()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if received["api_key"] != "test-api-key" {
		t.Errorf("expected api_key=test-api-key, got %v", received["api_key"])
	}

	events, ok := received["events"].([]interface{})
	if !ok || len(events) != 1 {
		t.Fatalf("expected 1 event, got %v", received["events"])
	}

	ev := events[0].(map[string]interface{})
	if ev["event_type"] != "port_opened" {
		t.Errorf("expected event_type=port_opened, got %v", ev["event_type"])
	}

	props := ev["event_properties"].(map[string]interface{})
	if int(props["port"].(float64)) != 9200 {
		t.Errorf("expected port=9200, got %v", props["port"])
	}
}

func TestAmplitude_InfoLevelUsesPortClosed(t *testing.T) {
	var received map[string]interface{}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &received)
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	al := amplitudeAlert()
	al.Level = alert.Info

	n := notify.NewAmplitude("key", srv.URL)
	_ = n.Send(al)

	events := received["events"].([]interface{})
	ev := events[0].(map[string]interface{})
	if ev["event_type"] != "port_closed" {
		t.Errorf("expected port_closed, got %v", ev["event_type"])
	}
}

func TestAmplitude_ErrorOnBadStatus(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer srv.Close()

	n := notify.NewAmplitude("bad-key", srv.URL)
	if err := n.Send(amplitudeAlert()); err == nil {
		t.Error("expected error on non-2xx status")
	}
}

func TestAmplitude_DefaultTimeout(t *testing.T) {
	n := notify.NewAmplitude("key", "")
	if n == nil {
		t.Fatal("expected non-nil notifier")
	}
}
