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

func datadogAlert() alert.Alert {
	return alert.Alert{
		Port:      8080,
		Message:   "port opened",
		Level:     alert.Warn,
		Timestamp: time.Now(),
	}
}

func TestDataDog_SendsCorrectPayload(t *testing.T) {
	var received map[string]interface{}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		if r.Header.Get("DD-API-KEY") != "test-key" {
			t.Errorf("missing or wrong DD-API-KEY header")
		}
		w.WriteHeader(http.StatusAccepted)
	}))
	defer ts.Close()

	dd := notify.NewDataDog("test-key", ts.URL)
	if err := dd.Send(datadogAlert()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if received["alert_type"] != "warning" {
		t.Errorf("expected alert_type=warning, got %v", received["alert_type"])
	}
	tags, ok := received["tags"].([]interface{})
	if !ok || len(tags) == 0 {
		t.Errorf("expected tags in payload")
	}
}

func TestDataDog_InfoLevelUsesInfoAlertType(t *testing.T) {
	var received map[string]interface{}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewDecoder(r.Body).Decode(&received)
		w.WriteHeader(http.StatusAccepted)
	}))
	defer ts.Close()

	a := datadogAlert()
	a.Level = alert.Info
	dd := notify.NewDataDog("key", ts.URL)
	if err := dd.Send(a); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if received["alert_type"] != "info" {
		t.Errorf("expected alert_type=info, got %v", received["alert_type"])
	}
}

func TestDataDog_ErrorOnBadStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer ts.Close()

	dd := notify.NewDataDog("bad-key", ts.URL)
	if err := dd.Send(datadogAlert()); err == nil {
		t.Fatal("expected error on non-2xx status")
	}
}

func TestDataDog_DefaultTimeout(t *testing.T) {
	dd := notify.NewDataDog("key", "")
	if dd == nil {
		t.Fatal("expected non-nil DataDog notifier")
	}
}
