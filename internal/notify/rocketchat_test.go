package notify_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/patrickdappollonio/portwatch/internal/alert"
	"github.com/patrickdappollonio/portwatch/internal/notify"
)

func rocketChatAlert() alert.Alert {
	return alert.Alert{
		Port:      9200,
		Level:     alert.LevelWarn,
		Message:   "port 9200 opened",
		Timestamp: time.Now(),
	}
}

func TestRocketChat_SendsCorrectPayload(t *testing.T) {
	var received map[string]interface{}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &received)
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	n := notify.NewRocketChat(srv.URL)
	if err := n.Send(rocketChatAlert()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if received["text"] == nil {
		t.Fatal("expected text field in payload")
	}
	attachments, ok := received["attachments"].([]interface{})
	if !ok || len(attachments) == 0 {
		t.Fatal("expected at least one attachment")
	}
	attach := attachments[0].(map[string]interface{})
	if attach["color"] != "#e74c3c" {
		t.Errorf("expected warn color #e74c3c, got %v", attach["color"])
	}
}

func TestRocketChat_InfoLevelUsesGreenColor(t *testing.T) {
	var received map[string]interface{}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &received)
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	a := rocketChatAlert()
	a.Level = alert.LevelInfo

	n := notify.NewRocketChat(srv.URL)
	if err := n.Send(a); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	attachments := received["attachments"].([]interface{})
	attach := attachments[0].(map[string]interface{})
	if attach["color"] != "#2ecc71" {
		t.Errorf("expected info color #2ecc71, got %v", attach["color"])
	}
}

func TestRocketChat_ErrorOnBadStatus(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	n := notify.NewRocketChat(srv.URL)
	if err := n.Send(rocketChatAlert()); err == nil {
		t.Fatal("expected error on 500 response")
	}
}

func TestRocketChat_DefaultTimeout(t *testing.T) {
	n := notify.NewRocketChat("http://example.com")
	if n == nil {
		t.Fatal("expected non-nil notifier")
	}
}
