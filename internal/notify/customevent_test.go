package notify

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/celzero/portwatch/internal/alert"
)

func customEventAlert() alert.Alert {
	return alert.Alert{
		Port:      9200,
		Proto:     "tcp",
		Level:     alert.Warn,
		Message:   "port 9200/tcp opened",
		Timestamp: time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC),
	}
}

func TestCustomEvent_SendsCorrectPayload(t *testing.T) {
	var got customEventPayload
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&got); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	n := NewCustomEvent(srv.URL)
	if err := n.Send(customEventAlert()); err != nil {
		t.Fatalf("Send: %v", err)
	}

	if got.Event != "port.opened" {
		t.Errorf("event = %q; want port.opened", got.Event)
	}
	if got.Port != 9200 {
		t.Errorf("port = %d; want 9200", got.Port)
	}
	if got.Proto != "tcp" {
		t.Errorf("proto = %q; want tcp", got.Proto)
	}
	if got.Level != "warn" {
		t.Errorf("level = %q; want warn", got.Level)
	}
}

func TestCustomEvent_InfoLevelUsesPortClosedEvent(t *testing.T) {
	var got customEventPayload
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&got) //nolint:errcheck
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	a := customEventAlert()
	a.Level = alert.Info
	NewCustomEvent(srv.URL).Send(a) //nolint:errcheck

	if got.Event != "port.closed" {
		t.Errorf("event = %q; want port.closed", got.Event)
	}
}

func TestCustomEvent_ErrorOnBadStatus(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	if err := NewCustomEvent(srv.URL).Send(customEventAlert()); err == nil {
		t.Error("expected error on 500 response")
	}
}

func TestCustomEvent_DefaultTimeout(t *testing.T) {
	n := NewCustomEvent("http://localhost:9999").(*customEventNotifier)
	if n.client.Timeout == 0 {
		t.Error("expected non-zero HTTP client timeout")
	}
}

func TestCustomEvent_SetsTimestampWhenZero(t *testing.T) {
	var got customEventPayload
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&got) //nolint:errcheck
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	a := customEventAlert()
	a.Timestamp = time.Time{}
	NewCustomEvent(srv.URL).Send(a) //nolint:errcheck

	if got.Timestamp.IsZero() {
		t.Error("expected timestamp to be populated")
	}
}
