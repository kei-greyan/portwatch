package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/user/portwatch/internal/alert"
)

const (
	pushbulletDefaultAPIURL = "https://api.pushbullet.com/v2/pushes"
	pushbulletDefaultTimeout = 10 * time.Second
)

type pushbulletPayload struct {
	Type  string `json:"type"`
	Title string `json:"title"`
	Body  string `json:"body"`
}

// PushbulletNotifier sends alerts via the Pushbullet API.
type PushbulletNotifier struct {
	apiKey string
	apiURL string
	client *http.Client
}

// NewPushbullet creates a PushbulletNotifier using the provided API key.
func NewPushbullet(apiKey string) *PushbulletNotifier {
	return &PushbulletNotifier{
		apiKey: apiKey,
		apiURL: pushbulletDefaultAPIURL,
		client: &http.Client{Timeout: pushbulletDefaultTimeout},
	}
}

// Send dispatches the alert to Pushbullet.
func (p *PushbulletNotifier) Send(a alert.Alert) error {
	if a.Timestamp.IsZero() {
		a.Timestamp = time.Now()
	}

	payload := pushbulletPayload{
		Type:  "note",
		Title: fmt.Sprintf("portwatch: %s", a.Level),
		Body:  fmt.Sprintf("[%s] port %d/%s — %s", a.Timestamp.Format(time.RFC3339), a.Port, a.Proto, a.Message),
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("pushbullet: marshal payload: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, p.apiURL, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("pushbullet: build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Access-Token", p.apiKey)

	resp, err := p.client.Do(req)
	if err != nil {
		return fmt.Errorf("pushbullet: send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("pushbullet: unexpected status %d", resp.StatusCode)
	}
	return nil
}
