package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/user/portwatch/internal/alert"
)

const larkDefaultTimeout = 10 * time.Second

// larkPayload is the message card payload sent to a Lark/Feishu webhook.
type larkPayload struct {
	MsgType string      `json:"msg_type"`
	Content larkContent `json:"content"`
}

type larkContent struct {
	Text string `json:"text"`
}

// larkNotifier sends alerts to a Lark (Feishu) incoming webhook.
type larkNotifier struct {
	webhookURL string
	client     *http.Client
}

// NewLark returns a Notifier that posts alerts to the given Lark webhook URL.
func NewLark(webhookURL string) Notifier {
	return &larkNotifier{
		webhookURL: webhookURL,
		client:     &http.Client{Timeout: larkDefaultTimeout},
	}
}

func (l *larkNotifier) Send(a alert.Alert) error {
	text := fmt.Sprintf("[%s] %s — port %d/%s",
		a.Level, a.Message, a.Port, a.Proto)

	payload := larkPayload{
		MsgType: "text",
		Content: larkContent{Text: text},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("lark: marshal payload: %w", err)
	}

	resp, err := l.client.Post(l.webhookURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("lark: post: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("lark: unexpected status %d", resp.StatusCode)
	}
	return nil
}
