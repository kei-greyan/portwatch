package config

import (
	"encoding/json"
	"os"
	"time"
)

// Default values.
const (
	DefaultInterval  = 30 * time.Second
	DefaultStatePath = "/var/lib/portwatch/state.json"
)

// Config holds runtime configuration for portwatch.
type Config struct {
	// Interval between port scans.
	Interval time.Duration `json:"interval"`

	// StatePath is the file used to persist port state across restarts.
	StatePath string `json:"state_path"`

	// AlertWebhook is an optional URL to POST alerts to.
	AlertWebhook string `json:"alert_webhook,omitempty"`

	// IgnorePorts lists ports that should never trigger alerts.
	IgnorePorts []uint16 `json:"ignore_ports,omitempty"`
}

// Load reads a JSON config file from path. Missing file returns defaults.
func Load(path string) (*Config, error) {
	cfg := defaults()

	f, err := os.Open(path)
	if os.IsNotExist(err) {
		return cfg, nil
	}
	if err != nil {
		return nil, err
	}
	defer f.Close()

	if err := json.NewDecoder(f).Decode(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

// Save writes the config as JSON to path.
func Save(path string, cfg *Config) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(cfg)
}

func defaults() *Config {
	return &Config{
		Interval:  DefaultInterval,
		StatePath: DefaultStatePath,
	}
}
