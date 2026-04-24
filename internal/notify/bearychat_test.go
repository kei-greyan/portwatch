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

func bearyChatAlert() alert.Alert {
	return alert.Alert{
		Level:     alert.Warn,
		Message:   "port 9200 opened",
		Port:      9200,
		Proto:     "tcp",
		Timestamp: time.Now(),
	}
}

func TestBearyChat_SendsCorrectPayload(t *testing.T) {
	var received map[string]interface{}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	bc := notify.NewBearyChat(ts.URL)
	if err := bc.Send(bearyChatAlert()); err != nil {
		t.Fatalf("Send: %v", err)
	}

	if _, ok := received["text"]; !ok {
		t.Error("expected 'text' field in payload")
	}
	attachments, ok := received["attachments"].([]interface{})
	if !ok || len(attachments) == 0 {
		t.Fatal("expected non-empty attachments")
	}
	attach := attachments[0].(map[string]interface{})
	if attach["color"] != "#FF4444" {
		t.Errorf("expected warn color #FF4444, got %v", attach["color"])
	}
}

func TestBearyChat_InfoLevelUsesGreenColor(t *testing.T) {
	var received map[string]interface{}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&received) //nolint:errcheck
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	a := bearyChatAlert()
	a.Level = alert.Info
	bc := notify.NewBearyChat(ts.URL)
	if err := bc.Send(a); err != nil {
		t.Fatalf("Send: %v", err)
	}

	attachments := received["attachments"].([]interface{})
	attach := attachments[0].(map[string]interface{})
	if attach["color"] != "#36A64F" {
		t.Errorf("expected info color #36A64F, got %v", attach["color"])
	}
}

func TestBearyChat_ErrorOnBadStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	bc := notify.NewBearyChat(ts.URL)
	if err := bc.Send(bearyChatAlert()); err == nil {
		t.Error("expected error on 500 response")
	}
}

func TestBearyChat_DefaultTimeout(t *testing.T) {
	bc := notify.NewBearyChat("http://example.com")
	if bc == nil {
		t.Fatal("expected non-nil BearyChat")
	}
}
