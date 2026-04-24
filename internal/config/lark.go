package config

import "fmt"

// LarkConfig holds settings for the Lark (Feishu) webhook notifier.
type LarkConfig struct {
	Enabled    bool   `toml:"enabled" json:"enabled"`
	WebhookURL string `toml:"webhook_url" json:"webhook_url"`
}

func defaultLarkConfig() *LarkConfig {
	return &LarkConfig{
		Enabled:    false,
		WebhookURL: "",
	}
}

// ValidateLark returns an error if cfg is enabled but misconfigured.
// A nil pointer is treated as disabled and always passes.
func ValidateLark(cfg *LarkConfig) error {
	if cfg == nil {
		return nil
	}
	if !cfg.Enabled {
		return nil
	}
	if cfg.WebhookURL == "" {
		return fmt.Errorf("lark: webhook_url must not be empty when enabled")
	}
	return nil
}
