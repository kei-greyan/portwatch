package notify

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/patrickdappollonio/portwatch/internal/alert"
)

func matrixAlert() alert.Alert {
	return alert.Alert{
		Level:   alert.LevelWarn,
		Message: "port 9090 opened",
		Port:    9090,
		At:      time.Now(),
	}
}

func TestMatrix_SendsCorrectPayload(t *testing.T) {
	var got matrixPayload

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer test-token" {
			t.Errorf("missing or wrong Authorization header: %q", r.Header.Get("Authorization"))
		}
		if err := json.NewDecoder(r.Body).Decode(&got); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	m := NewMatrix(srv.URL, "!room:example.com", "test-token")
	a := matrixAlert()

	if err := m.Send(a); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got.MsgType != "m.text" {
		t.Errorf("msgtype = %q, want m.text", got.MsgType)
	}
	if got.Body == "" {
		t.Error("body must not be empty")
	}
	if got.Format != "org.matrix.custom.html" {
		t.Errorf("format = %q, want org.matrix.custom.html", got.Format)
	}
	if got.FormattedBody == "" {
		t.Error("formatted_body must not be empty")
	}
}

func TestMatrix_ErrorOnBadStatus(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer srv.Close()

	m := NewMatrix(srv.URL, "!room:example.com", "bad-token")
	if err := m.Send(matrixAlert()); err == nil {
		t.Fatal("expected error on 403, got nil")
	}
}

func TestMatrix_DefaultTimeout(t *testing.T) {
	m := NewMatrix("https://matrix.example.com", "!room:example.com", "tok")
	if m.client.Timeout != defaultMatrixTimeout {
		t.Errorf("timeout = %v, want %v", m.client.Timeout, defaultMatrixTimeout)
	}
}
