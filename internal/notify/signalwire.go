package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/user/portwatch/internal/alert"
)

const signalWireDefaultTimeout = 10 * time.Second

// signalWireNotifier sends alerts via the SignalWire SMS REST API.
type signalWireNotifier struct {
	spaceURL string
	projectID string
	apiToken string
	from string
	to string
	client *http.Client
}

type signalWirePayload struct {
	From string `json:"from"`
	To   string `json:"to"`
	Body string `json:"body"`
}

// NewSignalWire returns a Notifier that sends SMS messages through SignalWire.
func NewSignalWire(spaceURL, projectID, apiToken, from, to string) Notifier {
	return &signalWireNotifier{
		spaceURL:  spaceURL,
		projectID: projectID,
		apiToken:  apiToken,
		from:      from,
		to:        to,
		client:    &http.Client{Timeout: signalWireDefaultTimeout},
	}
}

func (s *signalWireNotifier) Send(a alert.Alert) error {
	if a.Timestamp.IsZero() {
		a.Timestamp = time.Now()
	}

	body := fmt.Sprintf("[portwatch] %s port %d/%s at %s",
		a.Message, a.Port, a.Proto, a.Timestamp.Format(time.RFC3339))

	p := signalWirePayload{From: s.from, To: s.to, Body: body}
	b, err := json.Marshal(p)
	if err != nil {
		return fmt.Errorf("signalwire: marshal payload: %w", err)
	}

	url := fmt.Sprintf("https://%s/api/laml/2010-04-01/Accounts/%s/Messages.json",
		s.spaceURL, s.projectID)

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(b))
	if err != nil {
		return fmt.Errorf("signalwire: build request: %w", err)
	}
	req.SetBasicAuth(s.projectID, s.apiToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("signalwire: send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("signalwire: unexpected status %d", resp.StatusCode)
	}
	return nil
}
