package notify_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/notify"
)

func kafkaAlert() alert.Alert {
	return alert.Alert{
		Timestamp: time.Date(2024, 1, 2, 3, 4, 5, 0, time.UTC),
		Level:     alert.LevelWarn,
		Port:      9092,
		Proto:     "tcp",
		Message:   "port opened",
	}
}

func TestKafka_SendsCorrectPayload(t *testing.T) {
	var got map[string]interface{}
	var gotContentType string

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotContentType = r.Header.Get("Content-Type")
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &got)
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	n := notify.NewKafka(srv.URL, "portwatch-alerts")
	if err := n.Send(context.Background(), kafkaAlert()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if gotContentType != "application/vnd.kafka.json.v2+json" {
		t.Errorf("content-type = %q, want application/vnd.kafka.json.v2+json", gotContentType)
	}

	records, ok := got["records"].([]interface{})
	if !ok || len(records) != 1 {
		t.Fatalf("expected 1 record, got %v", got["records"])
	}

	val := records[0].(map[string]interface{})["value"].(map[string]interface{})
	if val["port"].(float64) != 9092 {
		t.Errorf("port = %v, want 9092", val["port"])
	}
	if val["level"] != alert.LevelWarn {
		t.Errorf("level = %v, want %s", val["level"], alert.LevelWarn)
	}
}

func TestKafka_ErrorOnBadStatus(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
	}))
	defer srv.Close()

	n := notify.NewKafka(srv.URL, "portwatch-alerts")
	if err := n.Send(context.Background(), kafkaAlert()); err == nil {
		t.Fatal("expected error on 503, got nil")
	}
}

func TestKafka_DefaultTimeout(t *testing.T) {
	n := notify.NewKafka("http://localhost:9999", "test")
	if n == nil {
		t.Fatal("expected non-nil notifier")
	}
}

func TestKafka_SetsTimestampWhenZero(t *testing.T) {
	var got map[string]interface{}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &got)
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	a := kafkaAlert()
	a.Timestamp = time.Time{}

	n := notify.NewKafka(srv.URL, "portwatch-alerts")
	if err := n.Send(context.Background(), a); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	records := got["records"].([]interface{})
	val := records[0].(map[string]interface{})["value"].(map[string]interface{})
	ts, ok := val["timestamp"].(string)
	if !ok || ts == "" {
		t.Error("expected non-empty timestamp")
	}
}
