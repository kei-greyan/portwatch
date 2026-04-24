package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/user/portwatch/internal/alert"
)

const bearyChatDefaultTimeout = 10 * time.Second

type bearyChatPayload struct {
	Text        string `json:"text"`
	Attachments []bearyChatAttachment `json:"attachments,omitempty"`
}

type bearyChatAttachment struct {
	Title string `json:"title"`
	Text  string `json:"text"`
	Color string `json:"color"`
}

// BearyChat sends alert notifications to a BearyChat incoming webhook.
type BearyChat struct {
	webhookURL string
	client     *http.Client
}

// NewBearyChat creates a new BearyChat notifier targeting the given webhook URL.
func NewBearyChat(webhookURL string) *BearyChat {
	return &BearyChat{
		webhookURL: webhookURL,
		client:     &http.Client{Timeout: bearyChatDefaultTimeout},
	}
}

// Send delivers the alert to the configured BearyChat webhook.
func (b *BearyChat) Send(a alert.Alert) error {
	color := "#FF4444"
	if a.Level == alert.Info {
		color = "#36A64F"
	}

	payload := bearyChatPayload{
		Text: fmt.Sprintf("**[portwatch]** %s", a.Message),
		Attachments: []bearyChatAttachment{
			{
				Title: fmt.Sprintf("Port %d/%s", a.Port, a.Proto),
				Text:  a.Message,
				Color: color,
			},
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("bearychat: marshal payload: %w", err)
	}

	resp, err := b.client.Post(b.webhookURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("bearychat: post: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("bearychat: unexpected status %d", resp.StatusCode)
	}
	return nil
}
