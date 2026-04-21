package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/user/portwatch/internal/alert"
)

const defaultTelegramTimeout = 10 * time.Second

// TelegramNotifier sends alerts to a Telegram chat via the Bot API.
type TelegramNotifier struct {
	token   string
	chatID  string
	apiBase string
	client  *http.Client
}

type telegramPayload struct {
	ChatID    string `json:"chat_id"`
	Text      string `json:"text"`
	ParseMode string `json:"parse_mode"`
}

// NewTelegram creates a TelegramNotifier for the given bot token and chat ID.
func NewTelegram(token, chatID string) *TelegramNotifier {
	return &TelegramNotifier{
		token:   token,
		chatID:  chatID,
		apiBase: "https://api.telegram.org",
		client:  &http.Client{Timeout: defaultTelegramTimeout},
	}
}

// Send dispatches the alert as a Telegram message.
func (t *TelegramNotifier) Send(a alert.Alert) error {
	emoji := "ℹ️"
	if a.Level == alert.Warn {
		emoji = "⚠️"
	}
	text := fmt.Sprintf("%s *%s* — port `%d` (%s)\n%s",
		emoji, a.Level, a.Port, a.Proto, a.Message)

	payload := telegramPayload{
		ChatID:    t.chatID,
		Text:      text,
		ParseMode: "Markdown",
	}
	body, err := json.Marshal(payload)
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
