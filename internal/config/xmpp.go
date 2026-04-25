package config

import "fmt"

// XMPPConfig holds connection parameters for the XMPP notifier.
type XMPPConfig struct {
	Enabled  bool   `yaml:"enabled"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	From     string `yaml:"from"`
	Password string `yaml:"password"`
	To       string `yaml:"to"`
}

func defaultXMPPConfig() *XMPPConfig {
	return &XMPPConfig{
		Enabled: false,
		Port:    5222,
	}
}

// ValidateXMPP returns an error if cfg contains invalid XMPP settings.
// A nil or disabled config is always valid.
func ValidateXMPP(cfg *XMPPConfig) error {
	if cfg == nil || !cfg.Enabled {
		return nil
	}
	if cfg.Host == "" {
		return fmt.Errorf("xmpp: host must not be empty")
	}
	if cfg.Port <= 0 || cfg.Port > 65535 {
		return fmt.Errorf("xmpp: port %d is out of range", cfg.Port)
	}
	if cfg.From == "" {
		return fmt.Errorf("xmpp: from must not be empty")
	}
	if cfg.To == "" {
		return fmt.Errorf("xmpp: to must not be empty")
	}
	return nil
}
