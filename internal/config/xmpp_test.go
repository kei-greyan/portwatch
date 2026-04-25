package config

import "testing"

func validXMPPConfig() *XMPPConfig {
	return &XMPPConfig{
		Enabled:  true,
		Host:     "xmpp.example.com",
		Port:     5222,
		From:     "portwatch@example.com",
		Password: "secret",
		To:       "admin@example.com",
	}
}

func TestValidateXMPP_NilAlwaysPasses(t *testing.T) {
	if err := ValidateXMPP(nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidateXMPP_DisabledAlwaysPasses(t *testing.T) {
	cfg := validXMPPConfig()
	cfg.Enabled = false
	if err := ValidateXMPP(cfg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidateXMPP_AcceptsValidConfig(t *testing.T) {
	if err := ValidateXMPP(validXMPPConfig()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidateXMPP_RejectsEmptyHost(t *testing.T) {
	cfg := validXMPPConfig()
	cfg.Host = ""
	if err := ValidateXMPP(cfg); err == nil {
		t.Fatal("expected error for empty host")
	}
}

func TestValidateXMPP_RejectsInvalidPort(t *testing.T) {
	cfg := validXMPPConfig()
	cfg.Port = 0
	if err := ValidateXMPP(cfg); err == nil {
		t.Fatal("expected error for port 0")
	}
}

func TestValidateXMPP_RejectsEmptyFrom(t *testing.T) {
	cfg := validXMPPConfig()
	cfg.From = ""
	if err := ValidateXMPP(cfg); err == nil {
		t.Fatal("expected error for empty from")
	}
}

func TestValidateXMPP_RejectsEmptyTo(t *testing.T) {
	cfg := validXMPPConfig()
	cfg.To = ""
	if err := ValidateXMPP(cfg); err == nil {
		t.Fatal("expected error for empty to")
	}
}
