package metrics_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/user/portwatch/internal/metrics"
)

func TestHandler_ReturnsJSON(t *testing.T) {
	m := metrics.New()
	h := m.Handler()

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
		t.Fatalf("expected application/json, got %s", ct)
	}
	var snap map[string]any
	if err := json.NewDecoder(rec.Body).Decode(&snap); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
}

func TestHandler_ReflectsScans(t *testing.T) {
	m := metrics.New()
	m.RecordScan()
	m.RecordScan()

	h := m.Handler()
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	h.ServeHTTP(rec, req)

	var snap map[string]any
	if err := json.NewDecoder(rec.Body).Decode(&snap); err != nil {
		t.Fatalf("decode: %v", err)
	}
	scans, ok := snap["total_scans"].(float64)
	if !ok || scans != 2 {
		t.Fatalf("expected total_scans=2, got %v", snap["total_scans"])
	}
}
