package config

import "fmt"

// PushbulletConfig holds settings for the Pushbullet notifier.
type PushbulletConfig struct {
	Enabled bool   `toml:"enabled" json:"enabled"`
	APIKey  string `toml:"api_key" json:"api_key"`
	APIURL  string `toml:"api_url" json:"api_url"`
}

func defaultPushbulletConfig() *PushbulletConfig {
	return &PushbulletConfig{
		Enabled: false,
		APIURL:  "https://api.pushbullet.com/v2/pushes",
	}
}

// ValidatePushbullet returns an error if cfg contains invalid values.
// A nil or disabled config is always valid.
func ValidatePushbullet(cfg *PushbulletConfig) error {
	if cfg == nil || !cfg.Enabled {
		return nil
	}
	if cfg.APIKey == "" {
		return fmt.Errorf("pushbullet: api_key must not be empty")
	}
	if cfg.APIURL == "" {
		return fmt.Errorf("pushbullet: api_url must not be empty")
	}
	return nil
}
