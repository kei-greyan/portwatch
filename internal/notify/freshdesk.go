package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/patrickward/portwatch/internal/alert"
)

// freshdeskTicket represents the payload sent to the Freshdesk API.
type freshdeskTicket struct {
	Subject     string `json:"subject"`
	Description string `json:"description"`
	Email       string `json:"email"`
	Priority    int    `json:"priority"`
	Status      int    `json:"status"`
	Tags        []string `json:"tags,omitempty"`
}

// freshdeskPriority maps alert levels to Freshdesk ticket priorities.
// 1=Low, 2=Medium, 3=High, 4=Urgent
func freshdeskPriority(level alert.Level) int {
	if level == alert.LevelWarn {
		return 3 // High
	}
	return 2 // Medium
}

// FreshdeskNotifier sends alert notifications as tickets to a Freshdesk helpdesk.
type FreshdeskNotifier struct {
	domain     string
	apiKey     string
	requesterEmail string
	client     *http.Client
}

// NewFreshdesk creates a FreshdeskNotifier that opens tickets via the Freshdesk
// REST API. domain is the subdomain (e.g. "acme" for acme.freshdesk.com),
// apiKey is the Freshdesk API key, and requesterEmail is the email address
// associated with the ticket requester.
func NewFreshdesk(domain, apiKey, requesterEmail string, opts ...func(*FreshdeskNotifier)) *FreshdeskNotifier {
	n := &FreshdeskNotifier{
		domain:         domain,
		apiKey:         apiKey,
		requesterEmail: requesterEmail,
		client:         &http.Client{Timeout: 10 * time.Second},
	}
	for _, o := range opts {
		o(n)
	}
	return n
}

// WithFreshdeskHTTPClient overrides the default HTTP client used by the notifier.
func WithFreshdeskHTTPClient(c *http.Client) func(*FreshdeskNotifier) {
	return func(n *FreshdeskNotifier) {
		n.client = c
	}
}

// Send opens a Freshdesk ticket for the given alert.
func (n *FreshdeskNotifier) Send(a alert.Alert) error {
	if a.Timestamp.IsZero() {
		a.Timestamp = time.Now().UTC()
	}

	subject := fmt.Sprintf("[portwatch] %s – port %d/%s",
		a.Title, a.Port, a.Proto)

	description := fmt.Sprintf(
		"<p><strong>%s</strong></p><p>%s</p><ul><li>Port: %d</li><li>Protocol: %s</li><li>Host: %s</li><li>Time: %s</li></ul>",
		a.Title,
		a.Message,
		a.Port,
		a.Proto,
		a.Host,
		a.Timestamp.Format(time.RFC3339),
	)

	ticket := freshdeskTicket{
		Subject:     subject,
		Description: description,
		Email:       n.requesterEmail,
		Priority:    freshdeskPriority(a.Level),
		Status:      2, // Open
		Tags:        []string{"portwatch", a.Proto},
	}

	body, err := json.Marshal(ticket)
	if err != nil {
		return fmt.Errorf("freshdesk: marshal payload: %w", err)
	}

	url := fmt.Sprintf("https://%s.freshdesk.com/api/v2/tickets", n.domain)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("freshdesk: build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(n.apiKey, "X") // Freshdesk uses API key as username, "X" as password

	resp, err := n.client.Do(req)
	if err != nil {
		return fmt.Errorf("freshdesk: send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("freshdesk: unexpected status %d", resp.StatusCode)
	}
	return nil
}
