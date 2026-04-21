package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/user/portwatch/internal/alert"
)

// MattermostConfig holds configuration for the Mattermost notifier.
type MattermostConfig struct {
	WebhookURL string
	Channel    string
	Username   string
	Timeout    time.Duration
}

type mattermostNotifier struct {
	cfg    MattermostConfig
	client *http.Client
}

type mattermostPayload struct {
	Channel  string `json:"channel,omitempty"`
	Username string `json:"username,omitempty"`
	Text     string `json:"text"`
}

// NewMattermost creates a Notifier that sends alerts to a Mattermost
// incoming webhook.
func NewMattermost(cfg MattermostConfig) Notifier {
	timeout := cfg.Timeout
	if timeout == 0 {
		timeout = 10 * time.Second
	}
	return &mattermostNotifier{
		cfg:    cfg,
		client: &http.Client{Timeout: timeout},
	}
}

func (m *mattermostNotifier) Send(a alert.Alert) error {
	text := fmt.Sprintf("**[%s]** %s — port `%d` (%s)",
		a.Level, a.Message, a.Port, a.Proto)

	payload := mattermostPayload{
		Channel:  m.cfg.Channel,
		Username: m.cfg.Username,
		Text:     text,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("mattermost: marshal payload: %w", err)
	}

	resp, err := m.client.Post(m.cfg.WebhookURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("mattermost: post: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("mattermost: unexpected status %d", resp.StatusCode)
	}
	return nil
}
