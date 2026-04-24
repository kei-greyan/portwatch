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
	lineDefaultAPIURL = "https://notify-api.line.me/api/notify"
	lineDefaultTimeout = 10 * time.Second
)

// lineNotifier sends alerts to a LINE Notify channel.
type lineNotifier struct {
	token   string
	apiURL  string
	client  *http.Client
}

type linePayload struct {
	Message string `json:"message"`
}

// NewLine returns a Notifier that posts to LINE Notify using the given token.
func NewLine(token, apiURL string) Notifier {
	if apiURL == "" {
		apiURL = lineDefaultAPIURL
	}
	return &lineNotifier{
		token:  token,
		apiURL: apiURL,
		client: &http.Client{Timeout: lineDefaultTimeout},
	}
}

func (l *lineNotifier) Send(a alert.Alert) error {
	if a.Timestamp.IsZero() {
		a.Timestamp = time.Now()
	}

	msg := fmt.Sprintf("[%s] %s — port %d/%s",
		a.Level, a.Message, a.Port, a.Proto)

	body, err := json.Marshal(linePayload{Message: msg})
	if err != nil {
		return fmt.Errorf("line: marshal payload: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, l.apiURL, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("line: create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+l.token)

	resp, err := l.client.Do(req)
	if err != nil {
		return fmt.Errorf("line: send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("line: unexpected status %d", resp.StatusCode)
	}
	return nil
}
