// Package config handles loading, saving, validating, and migrating
// portwatch configuration from a TOML file.
package config

import (
	"encoding/json"
	"errors"
	"os"
	"time"
)

// CurrentVersion is the schema version produced by this release.
const CurrentVersion = 1

// Config holds all portwatch runtime configuration.
type Config struct {
	Version      int           `json:"version"`
	Interval     time.Duration `json:"interval"`
	StatePath    string        `json:"state_path"`
	IgnoredPorts []string      `json:"ignored_ports"`
}

func defaults() *Config {
	return &Config{
		Version:      CurrentVersion,
		Interval:     30 * time.Second,
		StatePath:    "/var/lib/portwatch/state.json",
		IgnoredPorts: []string{},
	}
}

// Load reads config from path, returning defaults when the file is absent.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return defaults(), nil
	}
	if err != nil {
		return nil, err
	}
	cfg := defaults()
	if err := json.Unmarshal(data, cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

// Save writes cfg to path as JSON.
func Save(path string, cfg *Config) error {
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}
