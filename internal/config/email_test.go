package config_test

import (
	"testing"

	"github.com/user/portwatch/internal/config"
)

func validEmailConfig() *config.EmailConfig {
	return &config.EmailConfig{
		Enabled:  true,
		Host:     "smtp.example.com",
		Port:     587,
		From:     "portwatch@example.com",
		To:       []string{"admin@example.com"},
	}
}

func TestValidateEmail_DisabledAlwaysPasses(t *testing.T) {
	cfg := &config.EmailConfig{Enabled: false}
	if err := config.ValidateEmail(cfg); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestValidateEmail_NilAlwaysPasses(t *testing.T) {
	if err := config.ValidateEmail(nil); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestValidateEmail_AcceptsValidConfig(t *testing.T) {
	if err := config.ValidateEmail(validEmailConfig()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidateEmail_RejectsEmptyHost(t *testing.T) {
	cfg := validEmailConfig()
	cfg.Host = ""
	if err := config.ValidateEmail(cfg); err == nil {
		t.Fatal("expected error for empty host")
	}
}

func TestValidateEmail_RejectsInvalidPort(t *testing.T) {
	cfg := validEmailConfig()
	cfg.Port = 0
	if err := config.ValidateEmail(cfg); err == nil {
		t.Fatal("expected error for port 0")
	}
}

func TestValidateEmail_RejectsEmptyFrom(t *testing.T) {
	cfg := validEmailConfig()
	cfg.From = ""
	if err := config.ValidateEmail(cfg); err == nil {
		t.Fatal("expected error for empty from")
	}
}

func TestValidateEmail_RejectsEmptyTo(t *testing.T) {
	cfg := validEmailConfig()
	cfg.To = nil
	if err := config.ValidateEmail(cfg); err == nil {
		t.Fatal("expected error for empty to list")
	}
}
