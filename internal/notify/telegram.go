package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/yourorg/portwatch/internal/alert"
)

const telegramDefaultTimeout = 10 * time.Second

type telegramNotifier struct {
	token   string
	chatID  string
	client  *http.Client
	apiBase string
}

type telegramMessage struct {
	ChatID    string `json:"chat_id"`
	Text      string `json:"text"`
	ParseMode string `json:"parse_mode"`
}

// NewTelegram returns a Notifier that sends alerts via the Telegram Bot API.
func NewTelegram(token, chatID string) Notifier {
	return &telegramNotifier{
		token:   token,
		chatID:  chatID,
		client:  &http.Client{Timeout: telegramDefaultTimeout},
		apiBase: "https://api.telegram.org",
	}
}

func (t *telegramNotifier) Send(a alert.Alert) error {
	if a.Timestamp.IsZero() {
		a.Timestamp = time.Now()
	}

	emoji := "\u26a0\ufe0f"
	if a.Level == alert.Info {
		emoji = "\u2139\ufe0f"
	}

	text := fmt.Sprintf("%s *%s*\n`%s:%d/%s`\n%s",
		emoji,
		a.Title,
		a.Host, a.Port, a.Proto,
		a.Timestamp.UTC().Format(time.RFC3339),
	)

	msg := telegramMessage{
		ChatID:    t.chatID,
		Text:      text,
		ParseMode: "Markdown",
	}

	body, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("telegram: marshal: %w", err)
	}

	url := fmt.Sprintf("%s/bot%s/sendMessage", t.apiBase, t.token)
	resp, err := t.client.Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("telegram: post: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("telegram: unexpected status %d", resp.StatusCode)
	}
	return nil
}
