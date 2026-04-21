package config

import "fmt"

// SyslogConfig holds configuration for the syslog notifier.
type SyslogConfig struct {
	Enabled bool   `toml:"enabled"`
	Network string `toml:"network"` // "tcp", "udp", or "" for local
	Addr    string `toml:"addr"`    // "host:port" or "" for local
	Tag     string `toml:"tag"`
}

func defaultSyslogConfig() *SyslogConfig {
	return &SyslogConfig{
		Enabled: false,
		Network: "",
		Addr:    "",
		Tag:     "portwatch",
	}
}

// ValidateSyslog returns an error if cfg contains invalid values.
// A nil or disabled config is always valid.
func ValidateSyslog(cfg *SyslogConfig) error {
	if cfg == nil || !cfg.Enabled {
		return nil
	}
	if cfg.Network != "" && cfg.Network != "tcp" && cfg.Network != "udp" {
		return fmt.Errorf("syslog: network must be \"tcp\", \"udp\", or empty; got %q", cfg.Network)
	}
	if cfg.Network != "" && cfg.Addr == "" {
		return fmt.Errorf("syslog: addr is required when network is %q", cfg.Network)
	}
	if cfg.Tag == "" {
		return fmt.Errorf("syslog: tag must not be empty")
	}
	return nil
}
