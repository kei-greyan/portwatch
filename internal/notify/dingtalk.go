package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/user/portwatch/internal/alert"
)

// DingTalk sends alerts to a DingTalk group robot webhook.
type DingTalk struct {
	webhookURL string
	client     *http.Client
}

type dingtalkPayload struct {
	MsgType  string           `json:"msgtype"`
	Markdown dingtalkMarkdown `json:"markdown"`
}

type dingtalkMarkdown struct {
	Title string `json:"title"`
	Text  string `json:"text"`
}

// NewDingTalk creates a new DingTalk notifier using the provided webhook URL.
func NewDingTalk(webhookURL string) *DingTalk {
	return &DingTalk{
		webhookURL: webhookURL,
		client:     &http.Client{Timeout: 10 * time.Second},
	}
}

// Send dispatches an alert to the configured DingTalk webhook.
func (d *DingTalk) Send(a alert.Alert) error {
	if a.Timestamp.IsZero() {
		a.Timestamp = time.Now()
	}

	title := fmt.Sprintf("[portwatch] Port %d/%s %s", a.Port, a.Proto, a.Level)
	text := fmt.Sprintf(
		"**%s**\n\n- **Port:** %d/%s\n- **Level:** %s\n- **Time:** %s\n\n%s",
		title, a.Port, a.Proto, a.Level,
		a.Timestamp.Format(time.RFC3339), a.Message,
	)

	body, err := json.Marshal(dingtalkPayload{
		MsgType: "markdown",
		Markdown: dingtalkMarkdown{Title: title, Text: text},
	})
	if err != nil {
		return fmt.Errorf("dingtalk: marshal payload: %w", err)
	}

	resp, err := d.client.Post(d.webhookURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("dingtalk: http post: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("dingtalk: unexpected status %d", resp.StatusCode)
	}
	return nil
}
