package config

import "fmt"

// PushoverConfig holds configuration for the Pushover notifier.
type PushoverConfig struct {
	Enabled bool   `json:"enabled"`
	Token   string `json:"token"`
	UserKey string `json:"user_key"`
}

func defaultPushoverConfig() *PushoverConfig {
	return &PushoverConfig{
		Enabled: false,
	}
}

// ValidatePushover returns an error if cfg contains invalid values.
// A nil or disabled config is always valid.
func ValidatePushover(cfg *PushoverConfig) error {
	if cfg == nil || !cfg.Enabled {
		return nil
	}
	if cfg.Token == "" {
		return fmt.Errorf("pushover: token must not be empty")
	}
	if cfg.UserKey == "" {
		return fmt.Errorf("pushover: user_key must not be empty")
	}
	return nil
}
