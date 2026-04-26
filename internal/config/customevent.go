package config

import "fmt"

// CustomEventConfig holds settings for the custom-event HTTP notifier.
type CustomEventConfig struct {
	// Enabled controls whether the notifier is active.
	Enabled bool `toml:"enabled" json:"enabled"`

	// URL is the endpoint that receives POST requests with event payloads.
	URL string `toml:"url" json:"url"`
}

func defaultCustomEventConfig() *CustomEventConfig {
	return &CustomEventConfig{
		Enabled: false,
		URL:     "",
	}
}

// ValidateCustomEvent returns an error if cfg contains invalid settings.
// A nil cfg or a disabled config is always valid.
func ValidateCustomEvent(cfg *CustomEventConfig) error {
	if cfg == nil || !cfg.Enabled {
		return nil
	}
	if cfg.URL == "" {
		return fmt.Errorf("config: customevent: url must not be empty")
	}
	return nil
}
