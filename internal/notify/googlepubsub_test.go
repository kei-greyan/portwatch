package notify_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/example/portwatch/internal/alert"
	"github.com/example/portwatch/internal/notify"
)

func pubSubAlert() alert.Alert {
	return alert.Alert{
		Level:   alert.Warn,
		Message: "port 9200 opened (tcp)",
		Port:    9200,
		Proto:   "tcp",
		At:      time.Unix(1700000000, 0),
	}
}

func TestGooglePubSub_SendsCorrectPayload(t *testing.T) {
	var received map[string]interface{}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"messageIds":["1"]}`)) 
	}))
	defer srv.Close()

	n := notify.NewGooglePubSub("test-project", "test-topic",
		notify.WithGooglePubSubHTTPClient(srv.Client()),
		notify.WithGooglePubSubBaseURL(srv.URL),
	)

	if err := n.Send(context.Background(), pubSubAlert()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	messages, ok := received["messages"].([]interface{})
	if !ok || len(messages) == 0 {
		t.Fatal("expected messages array in payload")
	}
}

func TestGooglePubSub_ErrorOnBadStatus(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer srv.Close()

	n := notify.NewGooglePubSub("test-project", "test-topic",
		notify.WithGooglePubSubHTTPClient(srv.Client()),
		notify.WithGooglePubSubBaseURL(srv.URL),
	)

	if err := n.Send(context.Background(), pubSubAlert()); err == nil {
		t.Fatal("expected error on 403, got nil")
	}
}

func TestGooglePubSub_DefaultTimeout(t *testing.T) {
	n := notify.NewGooglePubSub("p", "t")
	if n == nil {
		t.Fatal("expected non-nil notifier")
	}
}
