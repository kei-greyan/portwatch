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

func pushbulletAlert() alert.Alert {
	return alert.Alert{
		Level:     alert.LevelWarn,
		Message:   "port opened",
		Port:      8080,
		Proto:     "tcp",
		Timestamp: time.Date(2024, 1, 2, 3, 4, 5, 0, time.UTC),
	}
}

func TestPushbullet_SendsCorrectPayload(t *testing.T) {
	var got map[string]string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Access-Token") != "test-key" {
			t.Errorf("missing or wrong Access-Token header")
		}
		if err := json.NewDecoder(r.Body).Decode(&got); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	n := notify.NewPushbullet("test-key")
	n.(*notify.PushbulletNotifier).APIURL = ts.URL // exposed for testing

	if err := n.Send(pushbulletAlert()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got["type"] != "note" {
		t.Errorf("expected type=note, got %q", got["type"])
	}
	if got["title"] == "" {
		t.Error("expected non-empty title")
	}
	if got["body"] == "" {
		t.Error("expected non-empty body")
	}
}

func TestPushbullet_ErrorOnBadStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer ts.Close()

	n := notify.NewPushbullet("bad-key")
	n.(*notify.PushbulletNotifier).APIURL = ts.URL

	if err := n.Send(pushbulletAlert()); err == nil {
		t.Fatal("expected error on 401, got nil")
	}
}

func TestPushbullet_DefaultTimeout(t *testing.T) {
	n := notify.NewPushbullet("key")
	pb := n.(*notify.PushbulletNotifier)
	if pb.Client.Timeout == 0 {
		t.Error("expected non-zero default timeout")
	}
}

func TestPushbullet_SetsTimestampWhenZero(t *testing.T) {
	var body map[string]string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewDecoder(r.Body).Decode(&body)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	n := notify.NewPushbullet("key")
	n.(*notify.PushbulletNotifier).APIURL = ts.URL

	a := pushbulletAlert()
	a.Timestamp = time.Time{}
	if err := n.Send(a); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if body["body"] == "" {
		t.Error("expected body to contain timestamp")
	}
}
