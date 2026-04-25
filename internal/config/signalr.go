package config

import "fmt"

// SignalRConfig holds settings for the Azure SignalR notifier.
type SignalRConfig struct {
	Enabled    bool   `toml:"enabled"`
	HubURL     string `toml:"hub_url"`
	APIKey     string `toml:"api_key"`
}

func defaultSignalRConfig() *SignalRConfig {
	return &SignalRConfig{
		Enabled: false,
		HubURL:  "",
		APIKey:  "",
	}
}

// ValidateSignalR returns an error if cfg is non-nil, enabled, and missing
// required fields.
func ValidateSignalR(cfg *SignalRConfig) error {
	if cfg == nil {
		return nil
	}
	if !cfg.Enabled {
		return nil
	}
	if cfg.HubURL == "" {
		return fmt.Errorf("signalr: hub_url must not be empty")
	}
	return nil
}
