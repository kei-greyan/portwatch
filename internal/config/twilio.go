package config

import "fmt"

// TwilioConfig holds credentials and addressing for Twilio SMS notifications.
type TwilioConfig struct {
	Enabled    bool   `toml:"enabled"`
	AccountSID string `toml:"account_sid"`
	AuthToken  string `toml:"auth_token"`
	From       string `toml:"from"`
	To         string `toml:"to"`
}

func defaultTwilioConfig() *TwilioConfig {
	return &TwilioConfig{
		Enabled: false,
	}
}

// ValidateTwilio returns an error if the Twilio configuration is invalid.
// A nil or disabled config is always valid.
func ValidateTwilio(c *TwilioConfig) error {
	if c == nil || !c.Enabled {
		return nil
	}
	if c.AccountSID == "" {
		return fmt.Errorf("twilio: account_sid must not be empty")
	}
	if c.AuthToken == "" {
		return fmt.Errorf("twilio: auth_token must not be empty")
	}
	if c.From == "" {
		return fmt.Errorf("twilio: from must not be empty")
	}
	if c.To == "" {
		return fmt.Errorf("twilio: to must not be empty")
	}
	return nil
}
