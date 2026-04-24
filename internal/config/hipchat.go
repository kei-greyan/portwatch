package config

import "fmt"

// HipChatConfig holds settings for the HipChat notifier.
type HipChatConfig struct {
	Enabled  bool   `toml:"enabled" json:"enabled"`
	RoomURL  string `toml:"room_url" json:"room_url"`
}

func defaultHipChatConfig() *HipChatConfig {
	return &HipChatConfig{
		Enabled: false,
		RoomURL: "",
	}
}

// ValidateHipChat returns an error if the HipChat configuration is invalid.
// A nil or disabled config is always considered valid.
func ValidateHipChat(c *HipChatConfig) error {
	if c == nil || !c.Enabled {
		return nil
	}
	if c.RoomURL == "" {
		return fmt.Errorf("hipchat: room_url must not be empty")
	}
	return nil
}
