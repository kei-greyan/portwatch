package config

import (
	"errors"
	"time"
)

// WebhookBatchConfig holds settings for the batch webhook notifier.
type WebhookBatchConfig struct {
	Enabled    bool          `toml:"enabled"`
	URL        string        `toml:"url"`
	Timeout    time.Duration `toml:"timeout"`
	MaxBatch   int           `toml:"max_batch"`
}

func defaultWebhookBatchConfig() WebhookBatchConfig {
	return WebhookBatchConfig{
		Enabled:  false,
		Timeout:  10 * time.Second,
		MaxBatch: 50,
	}
}

// ValidateWebhookBatch returns an error if cfg contains invalid values.
// A nil pointer or disabled config is always valid.
func ValidateWebhookBatch(cfg *WebhookBatchConfig) error {
	if cfg == nil || !cfg.Enabled {
		return nil
	}
	if cfg.URL == "" {
		return errors.New("webhook_batch: url must not be empty")
	}
	if cfg.Timeout < 0 {
		return errors.New("webhook_batch: timeout must be non-negative")
	}
	if cfg.MaxBatch <= 0 {
		return errors.New("webhook_batch: max_batch must be greater than zero")
	}
	return nil
}
