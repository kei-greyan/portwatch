package config

import "fmt"

// TelegramConfig holds settings for the Telegram Bot notifier.
type TelegramConfig struct {
	Enabled bool   `toml:"enabled"`
	Token   string `toml:"token"`
	ChatID  string `toml:"chat_id"`
}

func defaultTelegramConfig() *TelegramConfig {
	return &TelegramConfig{
		Enabled: false,
	}
}

// ValidateTelegram returns an error if cfg is enabled but incomplete.
func ValidateTelegram(cfg *TelegramConfig) error {
	if cfg == nil {
		return nil
	}
	if !cfg.Enabled {
		return nil
	}
	if cfg.Token == "" {
		return fmt.Errorf("telegram: token must not be empty")
	}
	if cfg.ChatID == "" {
		return fmt.Errorf("telegram: chat_id must not be empty")
	}
	return nil
}
