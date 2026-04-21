package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/user/portwatch/internal/alert"
)

// GotifyNotifier sends alerts to a self-hosted Gotify server.
type GotifyNotifier struct {
	baseURL  string
	token    string
	client   *http.Client
}

type gotifyPayload struct {
	Title    string `json:"title"`
	Message  string `json:"message"`
	Priority int    `json:"priority"`
}

// NewGotify constructs a GotifyNotifier.
// baseURL should be the root URL of the Gotify server (e.g. "https://gotify.example.com").
// token is the application token used to authenticate the push request.
func NewGotify(baseURL, token string, opts ...func(*GotifyNotifier)) *GotifyNotifier {
	g := &GotifyNotifier{
		baseURL: baseURL,
		token:   token,
		client:  &http.Client{Timeout: 10 * time.Second},
	}
	for _, o := range opts {
		o(g)
	}
	return g
}

// WithGotifyHTTPClient overrides the default HTTP client.
func WithGotifyHTTPClient(c *http.Client) func(*GotifyNotifier) {
	return func(g *GotifyNotifier) { g.client = c }
}

// Send delivers the alert to the Gotify server.
func (g *GotifyNotifier) Send(a alert.Alert) error {
	priority := 5 // default: normal
	if a.Level == alert.Warn {
		priority = 8
	}

	payload := gotifyPayload{
		Title:    fmt.Sprintf("portwatch: %s", a.Level),
		Message:  a.Message,
		Priority: priority,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("gotify: marshal payload: %w", err)
	}

	url := fmt.Sprintf("%s/message?token=%s", g.baseURL, g.token)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("gotify: build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := g.client.Do(req)
	if err != nil {
		return fmt.Errorf("gotify: send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("gotify: unexpected status %d", resp.StatusCode)
	}
	return nil
}
