package config_test

import (
	"testing"
	"time"

	"github.com/example/portwatch/internal/config"
)

func validConfig() *config.Config {
	return &config.Config{
		Interval:  30 * time.Second,
		StatePath: "/tmp/state.json",
	}
}

func TestValidate_AcceptsValidConfig(t *testing.T) {
	if err := config.Validate(validConfig()); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestValidate_RejectsNil(t *testing.T) {
	if err := config.Validate(nil); err == nil {
		t.Error("expected error for nil config")
	}
}

func TestValidate_RejectsShortInterval(t *testing.T) {
	cfg := validConfig()
	cfg.Interval = 500 * time.Millisecond
	if err := config.Validate(cfg); err == nil {
		t.Error("expected error for interval < 1s")
	}
}

func TestValidate_RejectsEmptyStatePath(t *testing.T) {
	cfg := validConfig()
	cfg.StatePath = ""
	if err := config.Validate(cfg); err == nil {
		t.Error("expected error for empty state_path")
	}
}

func TestValidate_RejectsZeroInIgnorePorts(t *testing.T) {
	cfg := validConfig()
	cfg.IgnorePorts = []uint16{80, 0, 443}
	if err := config.Validate(cfg); err == nil {
		t.Error("expected error for port 0 in ignore_ports")
	}
}
