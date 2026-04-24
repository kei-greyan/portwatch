package config

import "fmt"

// AmplitudeConfig holds settings for the Amplitude notifier.
type AmplitudeConfig struct {
	Enabled bool   `toml:"enabled"`
	APIKey  string `toml:"api_key"`
	APIURL  string `toml:"api_url"`
}

func defaultAmplitudeConfig() *AmplitudeConfig {
	return &AmplitudeConfig{
		Enabled: false,
		APIURL:  "https://api2.amplitude.com/2/httpapi",
	}
}

// ValidateAmplitude returns an error if the Amplitude config is invalid.
func ValidateAmplitude(c *AmplitudeConfig) error {
	if c == nil {
		return nil
	}
	if !c.Enabled {
		return nil
	}
	if c.APIKey == "" {
		return fmt.Errorf("amplitude: api_key must not be empty")
	}
	if c.APIURL == "" {
		return fmt.Errorf("amplitude: api_url must not be empty")
	}
	return nil
}
