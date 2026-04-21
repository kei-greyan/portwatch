package config_test

import (
	"testing"

	"github.com/yourusername/portwatch/internal/config"
)

var validRocketChatConfig = &config.RocketChatConfig{
	Enabled:    true,
	WebhookURL: "https://rocketchat.example.com/hooks/abc123",
}

func TestValidateRocketChat_NilAlwaysPasses(t *testing.T) {
	if err := config.ValidateRocketChat(nil); err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
}

func TestValidateRocketChat_DisabledAlwaysPasses(t *testing.T) {
	cfg := &config.RocketChatConfig{Enabled: false, WebhookURL: ""}
	if err := config.ValidateRocketChat(cfg); err != nil {
		t.Fatalf("expected nil error for disabled config, got %v", err)
	}
}

func TestValidateRocketChat_AcceptsValidConfig(t *testing.T) {
	if err := config.ValidateRocketChat(validRocketChatConfig); err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
}

func TestValidateRocketChat_RejectsEmptyWebhookURL(t *testing.T) {
	cfg := &config.RocketChatConfig{Enabled: true, WebhookURL: ""}
	if err := config.ValidateRocketChat(cfg); err == nil {
		t.Fatal("expected error for empty webhook URL, got nil")
	}
}
