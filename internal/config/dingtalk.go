package config

import "fmt"

// DingTalkConfig holds settings for the DingTalk notifier.
type DingTalkConfig struct {
	Enabled    bool   `toml:"enabled"`
	WebhookURL string `toml:"webhook_url"`
}

func defaultDingTalkConfig() *DingTalkConfig {
	return &DingTalkConfig{
		Enabled:    false,
		WebhookURL: "",
	}
}

// ValidateDingTalk returns an error if the DingTalk configuration is invalid.
// A nil or disabled config is always considered valid.
func ValidateDingTalk(c *DingTalkConfig) error {
	if c == nil || !c.Enabled {
		return nil
	}
	if c.WebhookURL == "" {
		return fmt.Errorf("dingtalk: webhook_url must not be empty")
	}
	return nil
}
