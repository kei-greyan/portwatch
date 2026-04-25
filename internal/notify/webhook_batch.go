package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/yourorg/portwatch/internal/alert"
)

// WebhookBatch sends multiple alerts as a single JSON array payload to a
// webhook endpoint. This is useful when many port changes occur simultaneously
// and you want to reduce the number of outbound HTTP requests.
type WebhookBatch struct {
	url    string
	client *http.Client
}

// NewWebhookBatch returns a WebhookBatch notifier that POSTs a JSON array of
// alerts to the given URL.
func NewWebhookBatch(url string, timeout time.Duration) *WebhookBatch {
	if timeout == 0 {
		timeout = 10 * time.Second
	}
	return &WebhookBatch{
		url:    url,
		client: &http.Client{Timeout: timeout},
	}
}

// SendBatch POSTs all alerts as a JSON array in a single HTTP request.
func (w *WebhookBatch) SendBatch(alerts []alert.Alert) error {
	if len(alerts) == 0 {
		return nil
	}
	for i := range alerts {
		if alerts[i].Timestamp.IsZero() {
			alerts[i].Timestamp = time.Now().UTC()
		}
	}
	body, err := json.Marshal(alerts)
	if err != nil {
		return fmt.Errorf("webhook_batch: marshal: %w", err)
	}
	resp, err := w.client.Post(w.url, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("webhook_batch: post: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("webhook_batch: unexpected status %d", resp.StatusCode)
	}
	return nil
}
