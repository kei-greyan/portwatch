package config

import "fmt"

// SplunkConfig holds settings for the Splunk HEC notifier.
type SplunkConfig struct {
	Enabled  bool   `toml:"enabled"`
	Endpoint string `toml:"endpoint"`
	Token    string `toml:"token"`
}

func defaultSplunkConfig() *SplunkConfig {
	return &SplunkConfig{
		Enabled:  false,
		Endpoint: "",
		Token:    "",
	}
}

// ValidateSplunk returns an error if the SplunkConfig is invalid.
// A nil or disabled config is always valid.
func ValidateSplunk(c *SplunkConfig) error {
	if c == nil || !c.Enabled {
		return nil
	}
	if c.Endpoint == "" {
		return fmt.Errorf("splunk: endpoint must not be empty")
	}
	if c.Token == "" {
		return fmt.Errorf("splunk: token must not be empty")
	}
	return nil
}
