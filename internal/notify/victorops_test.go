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

func victorOpsAlert() alert.Alert {
	return alert.Alert{
		Level:   alert.Warn,
		Message: "port 8080 opened",
		Port:    8080,
		At:      time.Date(2024, 1, 2, 3, 4, 5, 0, time.UTC),
	}
}

func TestVictorOps_SendsCorrectPayload(t *testing.T) {
	var got map[string]any
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&got); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	v := notify.NewVictorOps(ts.URL, "routingKey")
	if err := v.Send(victorOpsAlert()); err != nil {
		t.Fatalf("Send: %v", err)
	}

	if got["message_type"] != "CRITICAL" {
		t.Errorf("message_type = %v, want CRITICAL", got["message_type"])
	}
	if got["entity_display_name"] == "" {
		t.Error("entity_display_name should not be empty")
	}
	if got["state_message"] == "" {
		t.Error("state_message should not be empty")
	}
}

func TestVictorOps_InfoLevelUsesInfoMessageType(t *testing.T) {
	var got map[string]any
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewDecoder(r.Body).Decode(&got)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	a := victorOpsAlert()
	a.Level = alert.Info
	v := notify.NewVictorOps(ts.URL, "routingKey")
	_ = v.Send(a)

	if got["message_type"] != "INFO" {
		t.Errorf("message_type = %v, want INFO", got["message_type"])
	}
}

func TestVictorOps_ErrorOnBadStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	v := notify.NewVictorOps(ts.URL, "routingKey")
	if err := v.Send(victorOpsAlert()); err == nil {
		t.Error("expected error on non-2xx status")
	}
}

func TestVictorOps_DefaultTimeout(t *testing.T) {
	v := notify.NewVictorOps("http://example.com", "key")
	if v == nil {
		t.Fatal("NewVictorOps returned nil")
	}
}
