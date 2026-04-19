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

func slackAlert() alert.Alert {
	return alert.Alert{
		Level:   "WARN",
		Message: "port opened",
		Port:    8080,
		Proto:   "tcp",
	}
}

func TestSlack_SendsFormattedMessage(t *testing.T) {
	var got map[string]string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&got); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	s := notify.NewSlack(ts.URL, time.Second)
	if err := s.Send(slackAlert()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["text"] == "" {
		t.Error("expected non-empty text field")
	}
}

func TestSlack_ErrorOnBadStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	s := notify.NewSlack(ts.URL, time.Second)
	if err := s.Send(slackAlert()); err == nil {
		t.Error("expected error on 500 response")
	}
}

func TestSlack_DefaultTimeout(t *testing.T) {
	// passing zero timeout should not panic and should use default
	s := notify.NewSlack("http://localhost:0", 0)
	if s == nil {
		t.Fatal("expected non-nil notifier")
	}
}
