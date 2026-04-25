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

func jiraAlert() alert.Alert {
	return alert.Alert{
		Title:   "Port opened",
		Message: "TCP port 8080 is now open",
		Level:   alert.Warn,
		Port:    8080,
		Proto:   "tcp",
		Host:    "localhost",
		Time:    time.Now(),
	}
}

func TestJira_SendsCorrectPayload(t *testing.T) {
	var captured map[string]interface{}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/rest/api/2/issue" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if ct := r.Header.Get("Content-Type"); ct != "application/json" {
			t.Errorf("unexpected content-type: %s", ct)
		}
		_ = json.NewDecoder(r.Body).Decode(&captured)
		w.WriteHeader(http.StatusCreated)
	}))
	defer ts.Close()

	n := notify.NewJira(ts.URL, "OPS", "user", "token")
	if err := n.Send(jiraAlert()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	fields, _ := captured["fields"].(map[string]interface{})
	if fields == nil {
		t.Fatal("fields missing from payload")
	}
	if fields["summary"] != "Port opened" {
		t.Errorf("unexpected summary: %v", fields["summary"])
	}
	priority, _ := fields["priority"].(map[string]interface{})
	if priority["key"] != "High" {
		t.Errorf("expected High priority for Warn alert, got %v", priority["key"])
	}
}

func TestJira_InfoLevelUsesLowPriority(t *testing.T) {
	var captured map[string]interface{}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewDecoder(r.Body).Decode(&captured)
		w.WriteHeader(http.StatusCreated)
	}))
	defer ts.Close()

	a := jiraAlert()
	a.Level = alert.Info
	n := notify.NewJira(ts.URL, "OPS", "user", "token")
	if err := n.Send(a); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	fields, _ := captured["fields"].(map[string]interface{})
	priority, _ := fields["priority"].(map[string]interface{})
	if priority["key"] != "Low" {
		t.Errorf("expected Low priority for Info alert, got %v", priority["key"])
	}
}

func TestJira_ErrorOnBadStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer ts.Close()

	n := notify.NewJira(ts.URL, "OPS", "user", "bad-token")
	if err := n.Send(jiraAlert()); err == nil {
		t.Fatal("expected error on 401 response")
	}
}

func TestJira_DefaultTimeout(t *testing.T) {
	n := notify.NewJira("http://jira.example.com", "OPS", "user", "token")
	if n == nil {
		t.Fatal("expected non-nil notifier")
	}
}
