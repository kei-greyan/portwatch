package config

import "fmt"

// NtfyConfig holds configuration for the ntfy notification channel.
type NtfyConfig struct {
	Enabled bool   `toml:"enabled"`
	BaseURL string `toml:"base_url"`
	Topic   string `toml:"topic"`
	// Token is optional; used for authenticated/private topics.
	Token string `toml:"token,omitempty"`
}

func defaultNtfyConfig() *NtfyConfig {
	return &NtfyConfig{
		Enabled: false,
		BaseURL: "https://ntfy.sh",
	}
}

// ValidateNtfy returns an error if the NtfyConfig is invalid.
// A nil or disabled config is always valid.
func ValidateNtfy(c *NtfyConfig) error {
	if c == nil || !c.Enabled {
		return nil
	}
	if c.BaseURL == "" {
		return fmt.Errorf("ntfy: base_url must not be empty")
	}
	if c.Topic == "" {
		return fmt.Errorf("ntfy: topic must not be empty")
	}
	return nil
}
