package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/user/portwatch/internal/alert"
)

const defaultDiscordTimeout = 10 * time.Second

// discordPayload is the JSON body sent to a Discord webhook.
type discordPayload struct {
	Username string         `json:"username"`
	Embeds   []discordEmbed `json:"embeds"`
}

type discordEmbed struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Color       int    `json:"color"`
}

// discordColors maps alert levels to Discord embed color integers.
var discordColors = map[alert.Level]int{
	alert.LevelInfo: 0x5865F2, // blurple
	alert.LevelWarn: 0xFEE75C, // yellow
}

// Discord sends alerts to a Discord channel via an incoming webhook URL.
type Discord struct {
	webhookURL string
	client     *http.Client
}

// NewDiscord creates a Discord notifier that posts to the given webhook URL.
func NewDiscord(webhookURL string) *Discord {
	return &Discord{
		webhookURL: webhookURL,
		client:     &http.Client{Timeout: defaultDiscordTimeout},
	}
}

// Send delivers the alert to the configured Discord webhook.
func (d *Discord) Send(a alert.Alert) error {
	color, ok := discordColors[a.Level]
	if !ok {
		color = 0x99AAB5 // grey fallback
	}

	payload := discordPayload{
		Username: "portwatch",
		Embeds: []discordEmbed{
			{
				Title:       fmt.Sprintf("[%s] %s", a.Level, a.Title),
				Description: a.Body,
				Color:       color,
			},
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("discord: marshal payload: %w", err)
	}

	resp, err := d.client.Post(d.webhookURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("discord: post: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("discord: unexpected status %d", resp.StatusCode)
	}
	return nil
}
