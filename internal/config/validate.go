package config

import (
	"errors"
	"time"
)

// Validate checks that cfg holds sensible values and returns a descriptive
// error for the first violation found.
func Validate(cfg *Config) error {
	if cfg == nil {
		return errors.New("config must not be nil")
	}
	if cfg.Interval < time.Second {
		return errors.New("interval must be at least 1s")
	}
	if cfg.StatePath == "" {
		return errors.New("state_path must not be empty")
	}
	for _, p := range cfg.IgnorePorts {
		if p == 0 {
			return errors.New("ignore_ports contains invalid port 0")
		}
	}
	return nil
}
