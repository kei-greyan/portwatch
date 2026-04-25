package notify_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/notify"
)

func freshdeskAlert() alert.Alert {
	return alert.Alert{
		Title:   "Port 9200 opened",
		Message: "TCP port 9200 is now open",
		Level:   alert.Warn,
		At:      time.Now(),
	}
}

func TestFreshdesk_SendsCorrectPayload(t *testing.T) {
	var gotAuth, gotContentType string
	var gotBody []byte

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")
		gotContentType = r.Header.Get("Content-Type")
		gotBody, _ = io.ReadAll(r.Body)
		w.WriteHeader(http.StatusCreated)
	}))
	defer srv.Close()

	n := notify.NewFreshdesk(srv.URL, "mykey", "ops@example.com")
	if err := n.Send(freshdeskAlert()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if gotContentType != "application/json" {
		t.Errorf("content-type = %q, want application/json", gotContentType)
	}
	if gotAuth == "" {
		t.Error("expected Authorization header to be set")
	}
	if !bytes.Contains(gotBody, []byte("Port 9200 opened")) {
		t.Errorf("body missing title: %s", gotBody)
	}
}

func TestFreshdesk_InfoLevelUsesLowPriority(t *testing.T) {
	var gotBody []byte
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotBody, _ = io.ReadAll(r.Body)
		w.WriteHeader(http.StatusCreated)
	}))
	defer srv.Close()

	a := freshdeskAlert()
	a.Level = alert.Info
	n := notify.NewFreshdesk(srv.URL, "k", "a@b.com")
	_ = n.Send(a)

	if !bytes.Contains(gotBody, []byte(`"priority":1`)) {
		t.Errorf("expected low priority (1) for info alert, got: %s", gotBody)
	}
}

func TestFreshdesk_ErrorOnBadStatus(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer srv.Close()

	n := notify.NewFreshdesk(srv.URL, "bad", "a@b.com")
	if err := n.Send(freshdeskAlert()); err == nil {
		t.Error("expected error on 401 response")
	}
}

func TestFreshdesk_DefaultTimeout(t *testing.T) {
	n := notify.NewFreshdesk("http://localhost", "k", "a@b.com")
	_ = n // construction must not panic
}
