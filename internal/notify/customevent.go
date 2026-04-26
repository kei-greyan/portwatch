package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/celzero/portwatch/internal/alert"
)

// customEventPayload is the JSON body sent to a custom event endpoint.
type customEventPayload struct {
	Event     string    `json:"event"`
	Port      int       `json:"port"`
	Proto     string    `json:"proto"`
	Level     string    `json:"level"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
}

// customEventNotifier posts a structured event to an arbitrary HTTP endpoint.
type customEventNotifier struct {
	url    string
	client *http.Client
}

// NewCustomEvent returns a Notifier that posts alert data as a structured
// JSON event to the given URL. It is useful for integrating portwatch with
// in-house event pipelines that do not match any of the named notifiers.
func NewCustomEvent(url string) Notifier {
	return &customEventNotifier{
		url:    url,
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

func (n *customEventNotifier) Send(a alert.Alert) error {
	if a.Timestamp.IsZero() {
		a.Timestamp = time.Now().UTC()
	}

	eventName := "port.opened"
	if a.Level == alert.Info {
		eventName = "port.closed"
	}

	payload := customEventPayload{
		Event:     eventName,
		Port:      a.Port,
		Proto:     a.Proto,
		Level:     string(a.Level),
		Message:   a.Message,
		Timestamp: a.Timestamp,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("customevent: marshal: %w", err)
	}

	resp, err := n.client.Post(n.url, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("customevent: post: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("customevent: unexpected status %d", resp.StatusCode)
	}
	return nil
}
