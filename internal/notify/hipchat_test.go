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

func hipChatAlert() alert.Alert {
	return alert.Alert{
		Message:   "new port detected",
		Port:      8080,
		Proto:     "tcp",
		Level:     alert.LevelWarn,
		Timestamp: time.Now(),
	}
}

func TestHipChat_SendsCorrectPayload(t *testing.T) {
	var got map[string]interface{}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&got); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	n := notify.NewHipChat(srv.URL)
	if err := n.Send(hipChatAlert()); err != nil {
		t.Fatalf("Send: %v", err)
	}

	if got["color"] != "red" {
		t.Errorf("color = %q, want red", got["color"])
	}
	if got["notify"] != true {
		t.Errorf("notify = %v, want true", got["notify"])
	}
	msg, _ := got["message"].(string)
	if msg == "" {
		t.Error("message should not be empty")
	}
}

func TestHipChat_InfoLevelUsesGreenColor(t *testing.T) {
	var got map[string]interface{}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&got) //nolint:errcheck
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	a := hipChatAlert()
	a.Level = alert.LevelInfo
	n := notify.NewHipChat(srv.URL)
	if err := n.Send(a); err != nil {
		t.Fatalf("Send: %v", err)
	}
	if got["color"] != "green" {
		t.Errorf("color = %q, want green", got["color"])
	}
}

func TestHipChat_ErrorOnBadStatus(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer srv.Close()

	n := notify.NewHipChat(srv.URL)
	if err := n.Send(hipChatAlert()); err == nil {
		t.Error("expected error on 401")
	}
}

func TestHipChat_DefaultTimeout(t *testing.T) {
	n := notify.NewHipChat("http://example.com")
	_ = n // construction must not panic
}
