package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/user/portwatch/internal/alert"
)

type freshdeskPriority int

const (
	freshdeskPriorityLow    freshdeskPriority = 1
	freshdeskPriorityMedium freshdeskPriority = 2
	freshdeskPriorityHigh   freshdeskPriority = 3
	freshdeskPriorityUrgent freshdeskPriority = 4
)

type freshdeskNotifier struct {
	apiURL string
	apiKey string
	email  string
	client *http.Client
}

type freshdeskPayload struct {
	Subject     string            `json:"subject"`
	Description string            `json:"description"`
	Email       string            `json:"email"`
	Priority    freshdeskPriority `json:"priority"`
	Status      int               `json:"status"`
}

// NewFreshdesk returns a Notifier that opens a Freshdesk ticket for each alert.
func NewFreshdesk(apiURL, apiKey, requesterEmail string, opts ...func(*http.Client)) Notifier {
	c := &http.Client{Timeout: 10 * time.Second}
	for _, o := range opts {
		o(c)
	}
	return &freshdeskNotifier{apiURL: apiURL, apiKey: apiKey, email: requesterEmail, client: c}
}

// WithFreshdeskHTTPClient overrides the default HTTP client.
func WithFreshdeskHTTPClient(c *http.Client) func(*http.Client) {
	return func(existing *http.Client) {
		*existing = *c
	}
}

func (f *freshdeskNotifier) Send(a alert.Alert) error {
	pri := freshdeskPriorityHigh
	if a.Level == alert.Info {
		pri = freshdeskPriorityLow
	}

	body := freshdeskPayload{
		Subject:     fmt.Sprintf("[portwatch] %s", a.Title),
		Description: a.Message,
		Email:       f.email,
		Priority:    pri,
		Status:      2, // Open
	}

	raw, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("freshdesk: marshal: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, f.apiURL+"/api/v2/tickets", bytes.NewReader(raw))
	if err != nil {
		return fmt.Errorf("freshdesk: build request: %w", err)
	}
	req.SetBasicAuth(f.apiKey, "X")
	req.Header.Set("Content-Type", "application/json")

	resp, err := f.client.Do(req)
	if err != nil {
		return fmt.Errorf("freshdesk: send: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("freshdesk: unexpected status %d", resp.StatusCode)
	}
	return nil
}
