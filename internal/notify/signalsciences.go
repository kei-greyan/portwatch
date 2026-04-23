package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/user/portwatch/internal/alert"
)

const defaultSignalSciencesTimeout = 10 * time.Second

// SignalSciences sends alerts to the Signal Sciences / Fastly Next-Gen WAF
// custom event endpoint.
type SignalSciences struct {
	apiURL string
	corpName string
	siteName string
	accessKeyID string
	secretAccessKey string
	client *http.Client
}

type signalSciencesPayload struct {
	Event string `json:"event"`
	Message string `json:"message"`
	Level string `json:"level"`
	Timestamp string `json:"timestamp"`
}

// NewSignalSciences returns a Notifier that posts events to the Signal Sciences
// custom events API.
func NewSignalSciences(apiURL, corpName, siteName, accessKeyID, secretAccessKey string) *SignalSciences {
	return &SignalSciences{
		apiURL: apiURL,
		corpName: corpName,
		siteName: siteName,
		accessKeyID: accessKeyID,
		secretAccessKey: secretAccessKey,
		client: &http.Client{Timeout: defaultSignalSciencesTimeout},
	}
}

// Send delivers an alert to the Signal Sciences custom event endpoint.
func (s *SignalSciences) Send(a alert.Alert) error {
	if a.Timestamp.IsZero() {
		a.Timestamp = time.Now()
	}

	payload := signalSciencesPayload{
		Event: "portwatch",
		Message: a.Message,
		Level: a.Level,
		Timestamp: a.Timestamp.UTC().Format(time.RFC3339),
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("signalsciences: marshal payload: %w", err)
	}

	url := fmt.Sprintf("%s/v0/corps/%s/sites/%s/events", s.apiURL, s.corpName, s.siteName)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("signalsciences: build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-user", s.accessKeyID)
	req.Header.Set("x-api-token", s.secretAccessKey)

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("signalsciences: send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("signalsciences: unexpected status %d", resp.StatusCode)
	}
	return nil
}
