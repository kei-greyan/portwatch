package notify

import (
	"bytes"
	"fmt"
	"net/http"
	"time"

	"github.com/user/portwatch/internal/alert"
)

// NtfyConfig holds configuration for the ntfy.sh notification channel.
type NtfyConfig struct {
	// BaseURL is the ntfy server URL, e.g. https://ntfy.sh or a self-hosted instance.
	BaseURL string
	// Topic is the ntfy topic to publish to.
	Topic string
	// Token is an optional Bearer token for authenticated topics.
	Token string

	client *http.Client
}

type ntfyNotifier struct {
	cfg NtfyConfig
}

// NewNtfy returns a Notifier that publishes alerts to an ntfy topic.
func NewNtfy(cfg NtfyConfig) Notifier {
	if cfg.client == nil {
		cfg.client = &http.Client{Timeout: 10 * time.Second}
	}
	return &ntfyNotifier{cfg: cfg}
}

func (n *ntfyNotifier) Send(a alert.Alert) error {
	url := fmt.Sprintf("%s/%s", n.cfg.BaseURL, n.cfg.Topic)

	body := []byte(a.Message)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("ntfy: build request: %w", err)
	}

	req.Header.Set("Content-Type", "text/plain")
	req.Header.Set("Title", fmt.Sprintf("portwatch – %s", a.Level))
	if a.Level == alert.LevelWarn {
		req.Header.Set("Priority", "high")
	} else {
		req.Header.Set("Priority", "default")
	}
	if n.cfg.Token != "" {
		req.Header.Set("Authorization", "Bearer "+n.cfg.Token)
	}

	resp, err := n.cfg.client.Do(req)
	if err != nil {
		return fmt.Errorf("ntfy: send: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("ntfy: unexpected status %d", resp.StatusCode)
	}
	return nil
}
