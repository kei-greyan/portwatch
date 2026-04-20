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

const (
	opsGenieDefault * time.Second
	opsGenieAPIURL         = "https://api.opsgenie.com/v2/alertsn
// OpsGenie sends alerts to the OpsGenie Alerts API.
type OpsGenie struct {
	apiKey  string
	url     string
	client  *http.Client
}

type opsGeniePayload struct {
	Message     string            `json:"message"`
	Description string            `json:"description"`
	Priority    string            `json:"priority"`
	Tags        []string          `json:"tags,omitempty"`
	Details     map[string]string `json:"details,omitempty"`
}

// NewOpsGenie returns a Notifier that posts to the OpsGenie Alerts API.
// apiKey must be a valid OpsGenie API integration key.
func NewOpsGenie(apiKey string, opts ...func(*OpsGenie)) *OpsGenie {
	og := &OpsGenie{
		apiKey: apiKey,
		url:    opsGenieAPIURL,
		client: &http.Client{Timeout: opsGenieDefaultTimeout},
	}
	for _, o := range opts {
		o(og)
	}
	return og
}

// Send dispatches the alert to OpsGenie.
func (og *OpsGenie) Send(ctx context.Context, a alert.Alert) error {
	priority := "P3"
	if a.Level == alert.Warn {
		priority = "P2"
	}

	payload := opsGeniePayload{
		Message:     a.Title,
		Description: a.Body,
		Priority:    priority,
		Tags:        []string{"portwatch"},
		Details: map[string]string{
			"port": fmt.Sprintf("%d", a.Port),
			"host": a.Host,
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("opsgenie: marshal payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, og.url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("opsgenie: create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "GenieKey "+og.apiKey)

	resp, err := og.client.Do(req)
	if err != nil {
		return fmt.Errorf("opsgenie: send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("opsgenie: unexpected status %d", resp.StatusCode)
	}
	return nil
}
