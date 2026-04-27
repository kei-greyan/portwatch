package notify

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/example/portwatch/internal/alert"
)

func prowlAlert() alert.Alert {
	return alert.Alert{
		Title:   "Port Opened",
		Message: "TCP port 8080 is now open",
		Level:   alert.LevelWarn,
		At:      time.Now(),
	}
}

func TestProwl_SendsCorrectPayload(t *testing.T) {
	var gotForm url.Values
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		gotForm = r.PostForm
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	a := prowlAlert()
	n := &prowl{
		apiKey:  "test-key",
		apiURL:  ts.URL,
		appName: "portwatch",
		client:  ts.Client(),
	}

	if err := n.Send(a); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got := gotForm.Get("apikey"); got != "test-key" {
		t.Errorf("apikey: got %q, want %q", got, "test-key")
	}
	if got := gotForm.Get("application"); got != "portwatch" {
		t.Errorf("application: got %q, want %q", got, "portwatch")
	}
	if got := gotForm.Get("event"); got != a.Title {
		t.Errorf("event: got %q, want %q", got, a.Title)
	}
	if got := gotForm.Get("priority"); got != "1" {
		t.Errorf("priority: got %q, want \"1\" for LevelWarn", got)
	}
}

func TestProwl_InfoLevelUsesNormalPriority(t *testing.T) {
	var gotPriority string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		gotPriority = r.PostFormValue("priority")
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	a := prowlAlert()
	a.Level = alert.LevelInfo
	n := &prowl{apiKey: "k", apiURL: ts.URL, appName: "pw", client: ts.Client()}

	if err := n.Send(a); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotPriority != "0" {
		t.Errorf("priority: got %q, want \"0\" for LevelInfo", gotPriority)
	}
}

func TestProwl_ErrorOnBadStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer ts.Close()

	n := &prowl{apiKey: "bad", apiURL: ts.URL, appName: "pw", client: ts.Client()}
	if err := n.Send(prowlAlert()); err == nil {
		t.Fatal("expected error for 401 response")
	}
}

func TestProwl_DefaultTimeout(t *testing.T) {
	n := NewProwl("key", "app").(*prowl)
	if n.client.Timeout != prowlDefaultTimeout {
		t.Errorf("timeout: got %v, want %v", n.client.Timeout, prowlDefaultTimeout)
	}
}
