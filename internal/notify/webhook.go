package notify

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/user/portwatch/internal/alert"
)

// WebhookNotifier sends alerts as JSON POST requests to a URL.
type WebhookNotifier struct {
	url    string
	client *http.Client
}

// NewWebhook creates a WebhookNotifier targeting the given URL.
func NewWebhook(url string, timeout time.Duration) *WebhookNotifier {
	if timeout <= 0 {
		timeout = 10 * time.Second
	}
	return &WebhookNotifier{
		url: url,
		client: &http.Client{Timeout: timeout},
	}
}

type webhookPayload struct {
	Level   string    `json:"level"`
	Message string    `json:"message"`
	Port    uint16    `json:"port"`
	Proto   string    `json:"proto"`
	Time    time.Time `json:"time"`
}

// Send dispatches the alert to the configured webhook endpoint.
func (w *WebhookNotifier) Send(ctx context.Context, a alert.Alert) error {
	payload := webhookPayload{
		Level:   a.Level,
		Message: a.Message,
		Port:    a.Port,
		Proto:   a.Proto,
		Time:    a.Timestamp,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("webhook: marshal: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, w.url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("webhook: build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := w.client.Do(req)
	if err != nil {
		return fmt.Errorf("webhook: send: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return fmt.Errorf("webhook: unexpected status %d", resp.StatusCode)
	}
	return nil
}
