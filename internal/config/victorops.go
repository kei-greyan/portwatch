package config

import "fmt"

// VictorOpsConfig holds settings for the VictorOps (Splunk On-Call) notifier.
type VictorOpsConfig struct {
	// Enabled controls whether VictorOps notifications are sent.
	Enabled bool `json:"enabled"`

	// APIURL is the REST endpoint URL including the routing key path segment.
	// Example: https://alert.victorops.com/integrations/generic/20131114/alert/<api_key>/<routing_key>
	APIURL string `json:"api_url"`

	// RoutingKey identifies the escalation policy to invoke.
	RoutingKey string `json:"routing_key"`
}

func defaultVictorOpsConfig() *VictorOpsConfig {
	return &VictorOpsConfig{
		Enabled:    false,
		APIURL:     "",
		RoutingKey: "default",
	}
}

// ValidateVictorOps returns an error if cfg contains invalid settings.
// A nil or disabled config is always valid.
func ValidateVictorOps(cfg *VictorOpsConfig) error {
	if cfg == nil || !cfg.Enabled {
		return nil
	}
	if cfg.APIURL == "" {
		return fmt.Errorf("victorops: api_url must not be empty when enabled")
	}
	if cfg.RoutingKey == "" {
		return fmt.Errorf("victorops: routing_key must not be empty when enabled")
	}
	return nil
}
