package config

import "fmt"

// KafkaConfig holds settings for the Kafka REST Proxy notifier.
type KafkaConfig struct {
	Enabled  bool   `toml:"enabled"`
	ProxyURL string `toml:"proxy_url"`
	Topic    string `toml:"topic"`
}

func defaultKafkaConfig() *KafkaConfig {
	return &KafkaConfig{
		Enabled:  false,
		ProxyURL: "",
		Topic:    "portwatch-alerts",
	}
}

// ValidateKafka returns an error if cfg contains invalid Kafka settings.
// A nil or disabled config is always valid.
func ValidateKafka(cfg *KafkaConfig) error {
	if cfg == nil || !cfg.Enabled {
		return nil
	}
	if cfg.ProxyURL == "" {
		return fmt.Errorf("kafka: proxy_url must not be empty")
	}
	if cfg.Topic == "" {
		return fmt.Errorf("kafka: topic must not be empty")
	}
	return nil
}
