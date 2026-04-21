package notify

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/yourusername/portwatch/internal/alert"
)

const defaultZulipTimeout = 10 * time.Second

// ZulipNotifier sends alerts to a Zulip stream via the Zulip REST API.
type ZulipNotifier struct {
	siteURL  string
	botEmail string
	botToken string
	stream   string
	topic    string
	client   *http.Client
}

type zulipResponse struct {
	Result string `json:"result"`
	Msg    string `json:"msg"`
}

// NewZulip creates a ZulipNotifier that posts to the given stream and topic.
func NewZulip(siteURL, botEmail, botToken, stream, topic string) *ZulipNotifier {
	return &ZulipNotifier{
		siteURL:  strings.TrimRight(siteURL, "/"),
		botEmail: botEmail,
		botToken: botToken,
		stream:   stream,
		topic:    topic,
		client:   &http.Client{Timeout: defaultZulipTimeout},
	}
}

// Send delivers the alert as a Zulip stream message.
func (z *ZulipNotifier) Send(a alert.Alert) error {
	body := url.Values{}
	body.Set("type", "stream")
	body.Set("to", z.stream)
	body.Set("topic", z.topic)
	body.Set("content", fmt.Sprintf("**[%s]** %s — port %d", a.Level, a.Message, a.Port))

	endpoint := z.siteURL + "/api/v1/messages"
	req, err := http.NewRequest(http.MethodPost, endpoint, strings.NewReader(body.Encode()))
	if err != nil {
		return fmt.Errorf("zulip: build request: %w", err)
	}
	req.SetBasicAuth(z.botEmail, z.botToken)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := z.client.Do(req)
	if err != nil {
		return fmt.Errorf("zulip: send: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var zr zulipResponse
		_ = json.NewDecoder(resp.Body).Decode(&zr)
		return fmt.Errorf("zulip: unexpected status %d: %s", resp.StatusCode, zr.Msg)
	}
	return nil
}
