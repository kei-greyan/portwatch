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

func telegramAlert() alert.Alert {
	return alert.Alert{
		Title:     "Port opened",
		Host:      "localhost",
		Port:      8080,
		Proto:     "tcp",
		Level:     alert.Warn,
		Timestamp: time.Date(2024, 1, 2, 3, 4, 5, 0, time.UTC),
	}
}

func TestTelegram_SendsCorrectPayload(t *testing.T) {
	var got map[string]interface{}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&got); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	n := notify.NewTelegramWithBase("mytoken", "-100123", srv.URL)
	if err := n.Send(telegramAlert()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got["chat_id"] != "-100123" {
		t.Errorf("chat_id = %v, want -100123", got["chat_id"])
	}
	if got["parse_mode"] != "Markdown" {
		t.Errorf("parse_mode = %v, want Markdown", got["parse_mode"])
	}
	text, _ := got["text"].(string)
	if text == "" {
		t.Error("text must not be empty")
	}
}

func TestTelegram_ErrorOnBadStatus(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer srv.Close()

	n := notify.NewTelegramWithBase("bad", "123", srv.URL)
	if err := n.Send(telegramAlert()); err == nil {
		t.Fatal("expected error for 401 status")
	}
}

func TestTelegram_DefaultTimeout(t *testing.T) {
	n := notify.NewTelegram("tok", "cid")
	if n == nil {
		t.Fatal("NewTelegram returned nil")
	}
}
