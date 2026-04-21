package config

import "fmt"

// DataDogConfig holds settings for the Datadog notifier.
type DataDogConfig struct {
	Enabled bool   `toml:"enabled"`
	APIKey  string `toml:"api_key"`
	APIURL  string `toml:"api_url"`
}

func defaultDataDogConfig() *DataDogConfig {
	return &DataDogConfig{
		Enabled: false,
		APIURL:  "https://api.datadoghq.com/api/v1/events",
	}
}

// ValidateDataDog returns an error if the DataDog config is invalid.
// A nil or disabled config is always valid.
func ValidateDataDog(c *DataDogConfig) error {
	if c == nil || !c.Enabled {
		return nil
	}
	if c.APIKey == "" {
		return fmt.Errorf("datadog: api_key must not be empty")
	}
	if c.APIURL == "" {
		return fmt.Errorf("datadog: api_url must not be empty")
	}
	return nil
}
