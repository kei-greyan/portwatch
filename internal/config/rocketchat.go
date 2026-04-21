package config

import "fmt"

// RocketChatConfig holds configuration for the Rocket.Chat notifier.
type RocketChatConfig struct {
	Enabled    bool   `json:"enabled"`
	WebhookURL string `json:"webhook_url"`
}

func defaultRocketChatConfig() *RocketChatConfig {
	return &RocketChatConfig{
		Enabled:    false,
		WebhookURL: "",
	}
}

// ValidateRocketChat returns an error if the RocketChatConfig is invalid.
// A nil or disabled config is always valid.
func ValidateRocketChat(c *RocketChatConfig) error {
	if c == nil {
		return nil
	}
	if !c.Enabled {
		return nil
	}
	if c.WebhookURL == "" {
		return fmt.Errorf("rocketchat: webhook_url must not be empty when enabled")
	}
	return nil
}
