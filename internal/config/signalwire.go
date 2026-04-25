package config

import "fmt"

// SignalWireConfig holds configuration for the SignalWire SMS notifier.
type SignalWireConfig struct {
	// Enabled controls whether SignalWire notifications are active.
	Enabled bool `toml:"enabled" json:"enabled"`

	// SpaceURL is the SignalWire space URL, e.g. "example.signalwire.com".
	SpaceURL string `toml:"space_url" json:"space_url"`

	// ProjectID is the SignalWire project identifier.
	ProjectID string `toml:"project_id" json:"project_id"`

	// APIToken is the SignalWire API token used for authentication.
	APIToken string `toml:"api_token" json:"api_token"`

	// From is the SignalWire phone number to send messages from.
	From string `toml:"from" json:"from"`

	// To is the destination phone number to receive alerts.
	To string `toml:"to" json:"to"`
}

// defaultSignalWireConfig returns a SignalWireConfig with safe defaults.
func defaultSignalWireConfig() *SignalWireConfig {
	return &SignalWireConfig{
		Enabled: false,
	}
}

// ValidateSignalWire checks that cfg contains all required fields when enabled.
// A nil cfg or a disabled cfg is always valid.
func ValidateSignalWire(cfg *SignalWireConfig) error {
	if cfg == nil || !cfg.Enabled {
		return nil
	}
	if cfg.SpaceURL == "" {
		return fmt.Errorf("signalwire: space_url must not be empty")
	}
	if cfg.ProjectID == "" {
		return fmt.Errorf("signalwire: project_id must not be empty")
	}
	if cfg.APIToken == "" {
		return fmt.Errorf("signalwire: api_token must not be empty")
	}
	if cfg.From == "" {
		return fmt.Errorf("signalwire: from must not be empty")
	}
	if cfg.To == "" {
		return fmt.Errorf("signalwire: to must not be empty")
	}
	return nil
}
