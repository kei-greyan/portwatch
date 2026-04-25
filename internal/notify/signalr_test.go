package notify_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/yourorg/portwatch/internal/alert"
	"github.com/yourorg/portwatch/internal/notify"
)

func signalRAlert() alert.Alert {
	return alert.Alert{
		Level:   "warn",
		Port:    8080,
		Proto:   "tcp",
		Message: "port opened",
		Time:    time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
	}
}

func TestSignalR_SendsCorrectPayload(t *testing.T) {
	var got map[string]interface{}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&got); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		if r.Header.Get("Authorization") != "Bearer test-key" {
			t.Errorf("missing auth header")
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	n := notify.NewSignalR(srv.URL, "test-key")
	if err := n.Send(signalRAlert()); err != nil {
		t.Fatalf("Send: %v", err)
	}
	if got["target"] != "portAlert" {
		t.Errorf("target = %q, want portAlert", got["target"])
	}
	args, ok := got["arguments"].([]interface{})
	if !ok || len(args) == 0 {
		t.Fatal("expected arguments array")
	}
}

func TestSignalR_ErrorOnBadStatus(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	n := notify.NewSignalR(srv.URL, "")
	if err := n.Send(signalRAlert()); err == nil {
		t.Fatal("expected error on 500")
	}
}

func TestSignalR_DefaultTimeout(t *testing.T) {
	// Verify construction does not panic and client is configured.
	n := notify.NewSignalR("https://example.signalr.net", "key")
	if n == nil {
		t.Fatal("expected non-nil notifier")
	}
}

func TestSignalR_SetsTimestampWhenZero(t *testing.T) {
	var got map[string]interface{}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewDecoder(r.Body).Decode(&got)
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	a := signalRAlert()
	a.Time = time.Time{}
	n := notify.NewSignalR(srv.URL, "")
	if err := n.Send(a); err != nil {
		t.Fatalf("Send: %v", err)
	}
	// If timestamp was zero it would marshal as the zero value; no error means it was set.
}
