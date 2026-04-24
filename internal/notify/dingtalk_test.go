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

func dingtalkAlert() alert.Alert {
	return alert.Alert{
		Port:      8080,
		Proto:     "tcp",
		Level:     alert.Warn,
		Message:   "port opened unexpectedly",
		Timestamp: time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC),
	}
}

func TestDingTalk_SendsCorrectPayload(t *testing.T) {
	var received map[string]interface{}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	dt := notify.NewDingTalk(srv.URL)
	if err := dt.Send(dingtalkAlert()); err != nil {
		t.Fatalf("Send returned error: %v", err)
	}

	if received["msgtype"] != "markdown" {
		t.Errorf("expected msgtype=markdown, got %v", received["msgtype"])
	}
	md, ok := received["markdown"].(map[string]interface{})
	if !ok {
		t.Fatal("expected markdown field to be an object")
	}
	if md["title"] == "" {
		t.Error("expected non-empty title")
	}
	if md["text"] == "" {
		t.Error("expected non-empty text")
	}
}

func TestDingTalk_ErrorOnBadStatus(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	dt := notify.NewDingTalk(srv.URL)
	if err := dt.Send(dingtalkAlert()); err == nil {
		t.Fatal("expected error on 500 response")
	}
}

func TestDingTalk_DefaultTimeout(t *testing.T) {
	dt := notify.NewDingTalk("https://example.com/webhook")
	if dt == nil {
		t.Fatal("expected non-nil DingTalk notifier")
	}
}

func TestDingTalk_SetsTimestampWhenZero(t *testing.T) {
	var received map[string]interface{}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewDecoder(r.Body).Decode(&received)
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	a := dingtalkAlert()
	a.Timestamp = time.Time{}

	dt := notify.NewDingTalk(srv.URL)
	if err := dt.Send(a); err != nil {
		t.Fatalf("Send returned error: %v", err)
	}
}
