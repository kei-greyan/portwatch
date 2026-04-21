package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/user/portwatch/internal/alert"
)

const defaultSplunkTimeout = 10 * time.Second

// splunkEvent is the HEC (HTTP Event Collector) payload.
type splunkEvent struct {
	Time       float64        `json:"time"`
	SourceType string         `json:"sourcetype"`
	Event      map[string]any `json:"event"`
}

// splunkNotifier sends alerts to a Splunk HEC endpoint.
type splunkNotifier struct {
	endpoint string
	token    string
	client   *http.Client
}

// NewSplunk returns a Notifier that forwards alerts to Splunk via HEC.
func NewSplunk(endpoint, token string) Notifier {
	return &splunkNotifier{
		endpoint: endpoint,
		token:    token,
		client:   &http.Client{Timeout: defaultSplunkTimeout},
	}
}

func (s *splunkNotifier) Send(a alert.Alert) error {
	if a.Timestamp.IsZero() {
		a.Timestamp = time.Now()
	}

	payload := splunkEvent{
		Time:       float64(a.Timestamp.UnixNano()) / 1e9,
		SourceType: "portwatch",
		Event: map[string]any{
			"level":   a.Level,
			"message": a.Message,
			"port":    a.Port,
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("splunk: marshal payload: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, s.endpoint, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("splunk: build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Splunk "+s.token)

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("splunk: send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("splunk: unexpected status %d", resp.StatusCode)
	}
	return nil
}
