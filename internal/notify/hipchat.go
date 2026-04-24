package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/user/portwatch/internal/alert"
)

const defaultHipChatTimeout = 10 * time.Second

// hipChatPayload is the JSON body sent to the HipChat v2 API.
type hipChatPayload struct {
	Message       string `json:"message"`
	MessageFormat string `json:"message_format"`
	Color         string `json:"color"`
	Notify        bool   `json:"notify"`
}

// hipChatNotifier sends alerts to a HipChat room via the v2 REST API.
type hipChatNotifier struct {
	roomURL string // full URL: https://api.hipchat.com/v2/room/{id}/notification?auth_token=…
	client  *http.Client
}

// NewHipChat returns a Notifier that posts to the given HipChat room URL.
func NewHipChat(roomURL string, opts ...func(*hipChatNotifier)) *hipChatNotifier {
	n := &hipChatNotifier{
		roomURL: roomURL,
		client:  &http.Client{Timeout: defaultHipChatTimeout},
	}
	for _, o := range opts {
		o(n)
	}
	return n
}

// WithHipChatHTTPClient overrides the default HTTP client.
func WithHipChatHTTPClient(c *http.Client) func(*hipChatNotifier) {
	return func(n *hipChatNotifier) { n.client = c }
}

// Send delivers the alert to HipChat.
func (n *hipChatNotifier) Send(a alert.Alert) error {
	color := "yellow"
	if a.Level == alert.LevelWarn {
		color = "red"
	} else if a.Level == alert.LevelInfo {
		color = "green"
	}

	p := hipChatPayload{
		Message:       fmt.Sprintf("[portwatch] %s — port %d/%s", a.Message, a.Port, a.Proto),
		MessageFormat: "text",
		Color:         color,
		Notify:        a.Level == alert.LevelWarn,
	}

	body, err := json.Marshal(p)
	if err != nil {
		return fmt.Errorf("hipchat: marshal payload: %w", err)
	}

	resp, err := n.client.Post(n.roomURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("hipchat: post: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("hipchat: unexpected status %d", resp.StatusCode)
	}
	return nil
}
