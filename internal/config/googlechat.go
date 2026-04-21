package config

import "fmt"

// GoogleChatConfig holds settings for the Google Chat notifier.
type GoogleChatConfig struct {
	Enabled    bool   `json:"enabled" yaml:"enabled"`
	WebhookURL string `json:"webhook_url" yaml:"webhook_url"`
}

func defaultGoogleChatConfig() *GoogleChatConfig {
	return &GoogleChatConfig{
		Enabled:    false,
		WebhookURL: "",
	}
}

// ValidateGoogleChat returns an error if cfg is enabled but incomplete.
// A nil cfg is treated as disabled and is always valid.
func ValidateGoogleChat(cfg *GoogleChatConfig) error {
	if cfg == nil {
		return nil
	}
	if !cfg.Enabled {
		return nil
	}
	if cfg.WebhookURL == "" {
		return fmt.Errorf("googlechat: webhook_url must not be empty when enabled")
	}
	return nil
}
