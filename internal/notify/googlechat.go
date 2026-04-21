package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/user/portwatch/internal/alert"
)

const googleChatDefaultTimeout = 10 * time.Second

// GoogleChat sends alerts to a Google Chat webhook.
type GoogleChat struct {
	webhookURL string
	client     *http.Client
}

// NewGoogleChat returns a Notifier that posts to the given Google Chat
// incoming webhook URL.
func NewGoogleChat(webhookURL string) *GoogleChat {
	return &GoogleChat{
		webhookURL: webhookURL,
		client:     &http.Client{Timeout: googleChatDefaultTimeout},
	}
}

type googleChatPayload struct {
	Text string `json:"text"`
}

// Send delivers the alert to Google Chat.
func (g *GoogleChat) Send(a alert.Alert) error {
	emoji := "🔴"
	if a.Level == alert.Info {
		emoji = "🟢"
	}

	payload := googleChatPayload{
		Text: fmt.Sprintf("%s *portwatch* | %s | port %d (%s)",
			emoji, a.Message, a.Port, a.Proto),
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("googlechat: marshal payload: %w", err)
	}

	resp, err := g.client.Post(g.webhookURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("googlechat: post: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("googlechat: unexpected status %d", resp.StatusCode)
	}
	return nil
}
