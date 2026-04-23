package config

import "fmt"

// SignalSciencesConfig holds credentials and routing information for the
// Signal Sciences / Fastly Next-Gen WAF custom event notifier.
type SignalSciencesConfig struct {
	Enabled         bool   `toml:"enabled"`
	APIURL          string `toml:"api_url"`
	CorpName        string `toml:"corp_name"`
	SiteName        string `toml:"site_name"`
	AccessKeyID     string `toml:"access_key_id"`
	SecretAccessKey string `toml:"secret_access_key"`
}

func defaultSignalSciencesConfig() *SignalSciencesConfig {
	return &SignalSciencesConfig{
		Enabled: false,
		APIURL:  "https://dashboard.signalsciences.net/api",
	}
}

// ValidateSignalSciences returns an error if cfg is enabled but incomplete.
func ValidateSignalSciences(cfg *SignalSciencesConfig) error {
	if cfg == nil || !cfg.Enabled {
		return nil
	}
	if cfg.APIURL == "" {
		return fmt.Errorf("signalsciences: api_url must not be empty")
	}
	if cfg.CorpName == "" {
		return fmt.Errorf("signalsciences: corp_name must not be empty")
	}
	if cfg.SiteName == "" {
		return fmt.Errorf("signalsciences: site_name must not be empty")
	}
	if cfg.AccessKeyID == "" {
		return fmt.Errorf("signalsciences: access_key_id must not be empty")
	}
	if cfg.SecretAccessKey == "" {
		return fmt.Errorf("signalsciences: secret_access_key must not be empty")
	}
	return nil
}
