package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/user/portwatch/internal/alert"
)

const zendeskDefaultTimeout = 10 * time.Second

// zendeskPayload represents a Zendesk ticket creation request.
type zendeskPayload struct {
	Ticket zendeskTicket `json:"ticket"`
}

type zendeskTicket struct {
	Subject  string        `json:"subject"`
	Comment  zendeskComment `json:"comment"`
	Priority string        `json:"priority"`
	Tags     []string      `json:"tags"`
}

type zendeskComment struct {
	Body string `json:"body"`
}

// ZendeskNotifier sends alerts as Zendesk tickets via the Tickets API.
type ZendeskNotifier struct {
	subdomain string
	email     string
	token     string
	client    *http.Client
}

// NewZendesk creates a ZendeskNotifier that posts tickets to the given subdomain.
func NewZendesk(subdomain, email, token string) *ZendeskNotifier {
	return &ZendeskNotifier{
		subdomain: subdomain,
		email:     email,
		token:     token,
		client:    &http.Client{Timeout: zendeskDefaultTimeout},
	}
}

// Send delivers an alert as a Zendesk ticket.
func (z *ZendeskNotifier) Send(a alert.Alert) error {
	priority := "normal"
	if a.Level == alert.Warn {
		priority = "high"
	}

	body := fmt.Sprintf("Port %d/%s on host %s changed state: %s\nTimestamp: %s",
		a.Port, a.Proto, a.Host, a.Message, a.Timestamp.Format(time.RFC3339))

	payload := zendeskPayload{
		Ticket: zendeskTicket{
			Subject:  fmt.Sprintf("[portwatch] %s", a.Message),
			Comment:  zendeskComment{Body: body},
			Priority: priority,
			Tags:     []string{"portwatch", a.Proto},
		},
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("zendesk: marshal payload: %w", err)
	}

	url := fmt.Sprintf("https://%s.zendesk.com/api/v2/tickets.json", z.subdomain)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("zendesk: build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(z.email+"/token", z.token)

	resp, err := z.client.Do(req)
	if err != nil {
		return fmt.Errorf("zendesk: send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("zendesk: unexpected status %d", resp.StatusCode)
	}
	return nil
}
