package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/user/portwatch/internal/alert"
)

const (
	pagerDutyEventURL    = "https://events.pagerduty.com/v2/enqueue"
	defaultPDTimeout     = 10 * time.Second
	pagerDutySeverityErr = "error"
	pagerDutySeverityInfo = "info"
)

type pagerDutyPayload struct {
	RoutingKey  string          `json:"routing_key"`
	EventAction string          `json:"event_action"`
	Payload     pdInnerPayload  `json:"payload"`
}

type pdInnerPayload struct {
	Summary   string `json:"summary"`
	Severity  string `json:"severity"`
	Timestamp string `json:"timestamp"`
	Source    string `json:"source"`
}

// PagerDuty sends alerts to PagerDuty via the Events API v2.
type PagerDuty struct {
	routingKey string
	client     *http.Client
	url        string
}

// NewPagerDuty creates a new PagerDuty notifier with the given integration routing key.
func NewPagerDuty(routingKey string) *PagerDuty {
	return &PagerDuty{
		routingKey: routingKey,
		client:     &http.Client{Timeout: defaultPDTimeout},
		url:        pagerDutyEventURL,
	}
}

// Send dispatches an alert to PagerDuty.
func (p *PagerDuty) Send(a alert.Alert) error {
	severity := pagerDutySeverityInfo
	if a.Level == alert.Warn {
		severity = pagerDutySeverityErr
	}

	if a.Timestamp.IsZero() {
		a.Timestamp = time.Now()
	}

	body := pagerDutyPayload{
		RoutingKey:  p.routingKey,
		EventAction: "trigger",
		Payload: pdInnerPayload{
			Summary:   a.Message,
			Severity:  severity,
			Timestamp: a.Timestamp.UTC().Format(time.RFC3339),
			Source:    "portwatch",
		},
	}

	data, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("pagerduty: marshal payload: %w", err)
	}

	resp, err := p.client.Post(p.url, "application/json", bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("pagerduty: http post: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("pagerduty: unexpected status %d", resp.StatusCode)
	}
	return nil
}
