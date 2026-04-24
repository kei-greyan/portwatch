package config

import "testing"

func validLarkConfig() *LarkConfig {
	return &LarkConfig{
		Enabled:    true,
		WebhookURL: "https://open.larksuite.com/open-apis/bot/v2/hook/abc123",
	}
}

func TestValidateLark_NilAlwaysPasses(t *testing.T) {
	if err := ValidateLark(nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidateLark_DisabledAlwaysPasses(t *testing.T) {
	cfg := validLarkConfig()
	cfg.Enabled = false
	cfg.WebhookURL = ""
	if err := ValidateLark(cfg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidateLark_AcceptsValidConfig(t *testing.T) {
	if err := ValidateLark(validLarkConfig()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidateLark_RejectsEmptyWebhookURL(t *testing.T) {
	cfg := validLarkConfig()
	cfg.WebhookURL = ""
	if err := ValidateLark(cfg); err == nil {
		t.Error("expected error for empty webhook_url")
	}
}
