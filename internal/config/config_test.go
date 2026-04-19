package config_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/example/portwatch/internal/config"
)

func tmpPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "config.json")
}

func TestLoad_ReturnsDefaultsWhenMissing(t *testing.T) {
	cfg, err := config.Load("/nonexistent/portwatch-config-xyz.json")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Interval != config.DefaultInterval {
		t.Errorf("expected default interval, got %v", cfg.Interval)
	}
	if cfg.StatePath != config.DefaultStatePath {
		t.Errorf("expected default state path, got %s", cfg.StatePath)
	}
}

func TestLoad_ParsesValidFile(t *testing.T) {
	p := tmpPath(t)
	raw := `{"interval": 10000000000, "state_path": "/tmp/state.json", "alert_webhook": "http://example.com"}`
	if err := os.WriteFile(p, []byte(raw), 0o644); err != nil {
		t.Fatal(err)
	}
	cfg, err := config.Load(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Interval != 10*time.Second {
		t.Errorf("expected 10s, got %v", cfg.Interval)
	}
	if cfg.AlertWebhook != "http://example.com" {
		t.Errorf("unexpected webhook: %s", cfg.AlertWebhook)
	}
}

func TestLoad_ErrorOnCorruptFile(t *testing.T) {
	p := tmpPath(t)
	if err := os.WriteFile(p, []byte("not json"), 0o644); err != nil {
		t.Fatal(err)
	}
	_, err := config.Load(p)
	if err == nil {
		t.Fatal("expected error for corrupt file")
	}
}

func TestSave_PersistsAndLoads(t *testing.T) {
	p := tmpPath(t)
	orig := &config.Config{
		Interval:    5 * time.Second,
		StatePath:   "/tmp/s.json",
		IgnorePorts: []uint16{22, 80},
	}
	if err := config.Save(p, orig); err != nil {
		t.Fatalf("save error: %v", err)
	}
	loaded, err := config.Load(p)
	if err != nil {
		t.Fatalf("load error: %v", err)
	}
	data1, _ := json.Marshal(orig)
	data2, _ := json.Marshal(loaded)
	if string(data1) != string(data2) {
		t.Errorf("mismatch:\n got  %s\n want %s", data2, data1)
	}
}
