package config_test

import (
	"testing"

	"github.com/user/portwatch/internal/config"
)

func validTelegramConfig() *config.TelegramConfig {
	return &config.TelegramConfig{
		Enabled: true,
		Token:   "123456:ABC-DEFxxx",
		ChatID:  "-1001234567890",
	}
}

func TestValidateTelegram_NilAlwaysPasses(t *testing.T) {
	if err := config.ValidateTelegram(nil); err != nil {
		t.Fatalf("expected nil error for nil config, got %v", err)
	}
}

func TestValidateTelegram_DisabledAlwaysPasses(t *testing.T) {
	c := validTelegramConfig()
	c.Enabled = false
	c.Token = ""
	c.ChatID = ""
	if err := config.ValidateTelegram(c); err != nil {
		t.Fatalf("expected nil error for disabled config, got %v", err)
	}
}

func TestValidateTelegram_AcceptsValidConfig(t *testing.T) {
	if err := config.ValidateTelegram(validTelegramConfig()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidateTelegram_RejectsEmptyToken(t *testing.T) {
	c := validTelegramConfig()
	c.Token = ""
	if err := config.ValidateTelegram(c); err == nil {
		t.Fatal("expected error for empty token, got nil")
	}
}

func TestValidateTelegram_RejectsEmptyChatID(t *testing.T) {
	c := validTelegramConfig()
	c.ChatID = ""
	if err := config.ValidateTelegram(c); err == nil {
		t.Fatal("expected error for empty chat_id, got nil")
	}
}
