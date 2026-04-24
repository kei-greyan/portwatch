package config

import "fmt"

// WhatsAppConfig holds credentials for the WhatsApp Business Cloud API notifier.
type WhatsAppConfig struct {
	Enabled       bool   `toml:"enabled"`
	PhoneNumberID string `toml:"phone_number_id"`
	Token         string `toml:"token"`
	Recipient     string `toml:"recipient"`
}

func defaultWhatsAppConfig() *WhatsAppConfig {
	return &WhatsAppConfig{
		Enabled: false,
	}
}

// ValidateWhatsApp returns an error if the WhatsApp config is enabled but incomplete.
func ValidateWhatsApp(c *WhatsAppConfig) error {
	if c == nil {
		return nil
	}
	if !c.Enabled {
		return nil
	}
	if c.PhoneNumberID == "" {
		return fmt.Errorf("whatsapp: phone_number_id must not be empty")
	}
	if c.Token == "" {
		return fmt.Errorf("whatsapp: token must not be empty")
	}
	if c.Recipient == "" {
		return fmt.Errorf("whatsapp: recipient must not be empty")
	}
	return nil
}
