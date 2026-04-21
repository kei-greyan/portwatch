package config

import "fmt"

// TelegramConfig holds settings for the Telegram notifier.
type TelegramConfig struct {
	Enabled bool   `yaml:"enabled"`
	Token   string `yaml:"token"`
	ChatID  string `yaml:"chat_id"`
}

func defaultTelegramConfig() *TelegramConfig {
	return &TelegramConfig{
		Enabled: false,
	}
}

// ValidateTelegram returns an error if the Telegram configuration is invalid.
// A nil or disabled config is always valid.
func ValidateTelegram(c *TelegramConfig) error {
	if c == nil || !c.Enabled {
		return nil
	}
	if c.Token == "" {
		return fmt.Errorf("telegram: token must not be empty")
	}
	if c.ChatID == "" {
		return fmt.Errorf("telegram: chat_id must not be empty")
	}
	return nil
}
