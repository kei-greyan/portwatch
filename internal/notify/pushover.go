package notify

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/user/portwatch/internal/alert"
)

const (
	pushoverAPIURL        = "https://api.pushover.net/1/messages.json"
	pushoverDefaultTimeout = 10 * time.Second
	pushoverPriorityHigh  = 1
	pushoverPriorityNormal = 0
)

// PushoverConfig holds credentials for the Pushover notification service.
type PushoverConfig struct {
	Token   string
	UserKey string
	APIURL  string // override for testing
	Client  *http.Client
}

type pushoverNotifier struct {
	cfg PushoverConfig
}

type pushoverPayload struct {
	Token    string `json:"token"`
	User     string `json:"user"`
	Title    string `json:"title"`
	Message  string `json:"message"`
	Priority int    `json:"priority"`
}

// NewPushover returns a Notifier that sends alerts via Pushover.
func NewPushover(cfg PushoverConfig) Notifier {
	if cfg.APIURL == "" {
		cfg.APIURL = pushoverAPIURL
	}
	if cfg.Client == nil {
		cfg.Client = &http.Client{Timeout: pushoverDefaultTimeout}
	}
	return &pushoverNotifier{cfg: cfg}
}

func (p *pushoverNotifier) Send(a alert.Alert) error {
	priority := pushoverPriorityNormal
	if a.Level == alert.LevelWarn {
		priority = pushoverPriorityHigh
	}

	payload := pushoverPayload{
		Token:    p.cfg.Token,
		User:     p.cfg.UserKey,
		Title:    fmt.Sprintf("portwatch – %s", a.Level),
		Message:  a.Message,
		Priority: priority,
	}

	b, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("pushover: marshal payload: %w", err)
	}

	resp, err := p.cfg.Client.Post(p.cfg.APIURL, "application/json", strings.NewReader(string(b)))
	if err != nil {
		return fmt.Errorf("pushover: http post: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("pushover: unexpected status %d", resp.StatusCode)
	}
	return nil
}
