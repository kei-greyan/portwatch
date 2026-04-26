package config

import "fmt"

// GooglePubSubConfig holds settings for the Google Cloud Pub/Sub notifier.
type GooglePubSubConfig struct {
	// Enabled controls whether the notifier is active.
	Enabled bool `toml:"enabled" json:"enabled"`

	// ProjectID is the GCP project that owns the topic.
	ProjectID string `toml:"project_id" json:"project_id"`

	// TopicID is the Pub/Sub topic to publish to.
	TopicID string `toml:"topic_id" json:"topic_id"`

	// CredentialsFile is an optional path to a service-account JSON key.
	// When empty the SDK falls back to Application Default Credentials.
	CredentialsFile string `toml:"credentials_file,omitempty" json:"credentials_file,omitempty"`

	// TimeoutSeconds is the per-publish HTTP timeout (default: 10).
	TimeoutSeconds int `toml:"timeout_seconds,omitempty" json:"timeout_seconds,omitempty"`
}

func defaultGooglePubSubConfig() *GooglePubSubConfig {
	return &GooglePubSubConfig{
		Enabled:        false,
		TimeoutSeconds: 10,
	}
}

// ValidateGooglePubSub returns an error if cfg is invalid.
// A nil cfg or a disabled cfg is always valid.
func ValidateGooglePubSub(cfg *GooglePubSubConfig) error {
	if cfg == nil || !cfg.Enabled {
		return nil
	}
	if cfg.ProjectID == "" {
		return fmt.Errorf("googlepubsub: project_id must not be empty")
	}
	if cfg.TopicID == "" {
		return fmt.Errorf("googlepubsub: topic_id must not be empty")
	}
	if cfg.TimeoutSeconds <= 0 {
		return fmt.Errorf("googlepubsub: timeout_seconds must be positive")
	}
	return nil
}
