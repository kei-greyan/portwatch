package config

import "fmt"

// FreshdeskConfig holds settings for the Freshdesk notifier.
type FreshdeskConfig struct {
	Enabled         bool   `toml:"enabled"`
	APIURL          string `toml:"api_url"`
	APIKey          string `toml:"api_key"`
	RequesterEmail  string `toml:"requester_email"`
}

func defaultFreshdeskConfig() *FreshdeskConfig {
	return &FreshdeskConfig{
		Enabled: false,
		APIURL:  "https://yoursubdomain.freshdesk.com",
	}
}

// ValidateFreshdesk returns an error if cfg is enabled but incomplete.
func ValidateFreshdesk(cfg *FreshdeskConfig) error {
	if cfg == nil || !cfg.Enabled {
		return nil
	}
	if cfg.APIURL == "" {
		return fmt.Errorf("freshdesk: api_url must not be empty")
	}
	if cfg.APIKey == "" {
		return fmt.Errorf("freshdesk: api_key must not be empty")
	}
	if cfg.RequesterEmail == "" {
		return fmt.Errorf("freshdesk: requester_email must not be empty")
	}
	return nil
}
