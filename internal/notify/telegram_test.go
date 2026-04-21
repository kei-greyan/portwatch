package notify_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/notify"
)

func telegramAlert() alert.Alert {
	return alert.Alert{
		Level:     alert.Warn,
		Port:      9200,
		Proto:     "tcp",
		Message:   "port opened unexpectedly",
		Timestamp: time.Now(),
	}
}

func TestTelegram_SendsCorrectPayload(t *testing.T) {
	var received map[string]interface{}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	n := notify.NewTelegram("bot-token", "chat-123")
	n.(*notify.TelegramNotifier).SetAPIBase(srv.URL)

	if err := n.Send(telegramAlert()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if received["chat_id"] != "chat-123" {
		t.Errorf("chat_id = %v, want chat-123", received["chat_id"])
	}
	text, _ := received["text"].(string)
	if !strings.Contains(text, "9200") {
		t.Errorf("text does not contain port: %q", text)
	}
	if received["parse_mode"] != "Markdown" {
		t.Errorf("parse_mode = %v, want Markdown", received["parse_mode"])
	}
}

func TestTelegram_ErrorOnBadStatus(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer srv.Close()

	n := notify.NewTelegram("bad-token", "chat-123")
	n.(*notify.TelegramNotifier).SetAPIBase(srv.URL)

	if err := n.Send(telegramAlert()); err == nil {
		t.Fatal("expected error on 401, got nil")
	}
}

func TestTelegram_DefaultTimeout(t *testing.T) {
	n := notify.NewTelegram("tok", "cid")
	tn := n.(*notify.TelegramNotifier)
	if tn.Timeout() != 10*time.Second {
		t.Errorf("timeout = %v, want 10s", tn.Timeout())
	}
}
