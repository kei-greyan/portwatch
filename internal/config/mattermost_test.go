package config_test

import (
	"testing"

	"github.com/user/portwatch/internal/config"
)

func validMattermostConfig() *config.MattermostConfig {
	return &config.MattermostConfig{
		Enabled:    true,
		WebhookURL: "https://mattermost.example.com/hooks/abc123",
		Channel:    "#ops",
		Username:   "portwatch",
	}
}

func TestValidateMattermost_NilAlwaysPasses(t *testing.T) {
	if err := config.ValidateMattermost(nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidateMattermost_DisabledAlwaysPasses(t *testing.T) {
	cfg := validMattermostConfig()
	cfg.Enabled = false
	cfg.WebhookURL = ""
	if err := config.ValidateMattermost(cfg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidateMattermost_AcceptsValidConfig(t *testing.T) {
	if err := config.ValidateMattermost(validMattermostConfig()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidateMattermost_RejectsEmptyWebhookURL(t *testing.T) {
	cfg := validMattermostConfig()
	cfg.WebhookURL = ""
	if err := config.ValidateMattermost(cfg); err == nil {
		t.Fatal("expected error for empty webhook_url")
	}
}
