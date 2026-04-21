package config

import "fmt"

// MattermostConfig holds Mattermost notifier settings persisted in the
// portwatch configuration file.
type MattermostConfig struct {
	Enabled    bool   `json:"enabled"`
	WebhookURL string `json:"webhook_url"`
	Channel    string `json:"channel,omitempty"`
	Username   string `json:"username,omitempty"`
}

func defaultMattermostConfig() *MattermostConfig {
	return &MattermostConfig{
		Enabled:  false,
		Username: "portwatch",
	}
}

// ValidateMattermost returns an error if cfg contains invalid values.
// A nil pointer or a disabled config is always considered valid.
func ValidateMattermost(cfg *MattermostConfig) error {
	if cfg == nil || !cfg.Enabled {
		return nil
	}
	if cfg.WebhookURL == "" {
		return fmt.Errorf("mattermost: webhook_url must not be empty")
	}
	return nil
}
