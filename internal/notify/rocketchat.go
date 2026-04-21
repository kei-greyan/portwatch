package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/patrickdappollonio/portwatch/internal/alert"
)

const defaultRocketChatTimeout = 10 * time.Second

// RocketChat sends alerts to a Rocket.Chat incoming webhook.
type RocketChat struct {
	webhookURL string
	client     *http.Client
}

type rocketChatPayload struct {
	Text        string             `json:"text"`
	Attachments []rcAttachment     `json:"attachments,omitempty"`
}

type rcAttachment struct {
	Title string `json:"title"`
	Text  string `json:"text"`
	Color string `json:"color"`
}

// NewRocketChat returns a Notifier that posts to a Rocket.Chat webhook URL.
func NewRocketChat(webhookURL string) *RocketChat {
	return &RocketChat{
		webhookURL: webhookURL,
		client:     &http.Client{Timeout: defaultRocketChatTimeout},
	}
}

// Send delivers the alert to Rocket.Chat.
func (r *RocketChat) Send(a alert.Alert) error {
	color := "#e74c3c" // red for warn
	if a.Level == alert.LevelInfo {
		color = "#2ecc71" // green for info
	}

	payload := rocketChatPayload{
		Text: fmt.Sprintf("*portwatch*: %s", a.Message),
		Attachments: []rcAttachment{
			{
				Title: fmt.Sprintf("Port %d — %s", a.Port, a.Level),
				Text:  a.Message,
				Color: color,
			},
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("rocketchat: marshal payload: %w", err)
	}

	resp, err := r.client.Post(r.webhookURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("rocketchat: post: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("rocketchat: unexpected status %d", resp.StatusCode)
	}
	return nil
}
