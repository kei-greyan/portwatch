package config

import "fmt"

// ProwlConfig holds configuration for the Prowl push notification notifier.
type ProwlConfig struct {
	Enabled bool   `toml:"enabled"`
	APIKey  string `toml:"api_key"`
	AppName string `toml:"app_name"`
}

func defaultProwlConfig() *ProwlConfig {
	return &ProwlConfig{
		Enabled: false,
		AppName: "portwatch",
	}
}

// ValidateProwl returns an error if c contains invalid Prowl settings.
// A nil or disabled config is always valid.
func ValidateProwl(c *ProwlConfig) error {
	if c == nil || !c.Enabled {
		return nil
	}
	if c.APIKey == "" {
		return fmt.Errorf("prowl: api_key must not be empty")
	}
	if c.AppName == "" {
		return fmt.Errorf("prowl: app_name must not be empty")
	}
	return nil
}
