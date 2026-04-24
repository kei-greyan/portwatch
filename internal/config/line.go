package config

import "fmt"

// LineConfig holds settings for the LINE Notify integration.
type LineConfig struct {
	Enabled bool   `toml:"enabled"`
	Token   string `toml:"token"`
	APIURL  string `toml:"api_url"`
}

func defaultLineConfig() *LineConfig {
	return &LineConfig{
		Enabled: false,
		APIURL:  "https://notify-api.line.me/api/notify",
	}
}

// ValidateLine returns an error if the LineConfig is invalid.
// A nil or disabled config is always valid.
func ValidateLine(c *LineConfig) error {
	if c == nil || !c.Enabled {
		return nil
	}
	if c.Token == "" {
		return fmt.Errorf("line: token must not be empty")
	}
	if c.APIURL == "" {
		return fmt.Errorf("line: api_url must not be empty")
	}
	return nil
}
