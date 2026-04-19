package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/user/portwatch/internal/alert"
)

// SlackNotifier sends alerts to a Slack incoming webhook URL.
type SlackNotifier struct {
	webhookURL string
	client     *http.Client
}

type slackPayload struct {
	Text string `json:"text"`
}

// NewSlack creates a SlackNotifier that posts to the given Slack webhook URL.
func NewSlack(webhookURL string, timeout time.Duration) *SlackNotifier {
	if timeout <= 0 {
		timeout = 5 * time.Second
	}
	return &SlackNotifier{
		webhookURL: webhookURL,
		client:     &http.Client{Timeout: timeout},
	}
}

// Send delivers the alert as a formatted Slack message.
func (s *SlackNotifier) Send(a alert.Alert) error {
	text := fmt.Sprintf("[%s] %s — port %d (%s)",
		a.Level, a.Message, a.Port, a.Proto)

	body, err := json.Marshal(slackPayload{Text: text})
	if err != nil {
		return fmt.Errorf("slack: marshal payload: %w", err)
	}

	resp, err := s.client.Post(s.webhookURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("slack: post: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("slack: unexpected status %d", resp.StatusCode)
	}
	return nil
}
