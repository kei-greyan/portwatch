package config_test

import (
	"testing"

	"github.com/user/portwatch/internal/config"
)

func validGoogleChatConfig() *config.GoogleChatConfig {
	return &config.GoogleChatConfig{
		Enabled:    true,
		WebhookURL: "https://chat.googleapis.com/v1/spaces/ABC/messages?key=k&token=t",
	}
}

func TestValidateGoogleChat_NilAlwaysPasses(t *testing.T) {
	if err := config.ValidateGoogleChat(nil); err != nil {
		t.Fatalf("expected nil error for nil config, got %v", err)
	}
}

func TestValidateGoogleChat_DisabledAlwaysPasses(t *testing.T) {
	cfg := validGoogleChatConfig()
	cfg.Enabled = false
	cfg.WebhookURL = ""
	if err := config.ValidateGoogleChat(cfg); err != nil {
		t.Fatalf("expected nil error for disabled config, got %v", err)
	}
}

func TestValidateGoogleChat_AcceptsValidConfig(t *testing.T) {
	if err := config.ValidateGoogleChat(validGoogleChatConfig()); err != nil {
		t.Fatalf("expected nil error for valid config, got %v", err)
	}
}

func TestValidateGoogleChat_RejectsEmptyWebhookURL(t *testing.T) {
	cfg := validGoogleChatConfig()
	cfg.WebhookURL = ""
	if err := config.ValidateGoogleChat(cfg); err == nil {
		t.Error("expected error for empty webhook_url")
	}
}
