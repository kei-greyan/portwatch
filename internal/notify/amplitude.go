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
	amplitudeDefaultTimeout = 10 * time.Second
	amplitudeDefaultAPIURL  = "https://api2.amplitude.com/2/httpapi"
)

type amplitudeEvent struct {
	UserID    string                 `json:"user_id"`
	EventType string                 `json:"event_type"`
	EventTime int64                  `json:"time"`
	EventProp map[string]interface{} `json:"event_properties"`
}

type amplitudePayload struct {
	APIKey string           `json:"api_key"`
	Events []amplitudeEvent `json:"events"`
}

// AmplitudeNotifier sends port-change alerts as Amplitude events.
type AmplitudeNotifier struct {
	apiKey string
	apiURL string
	client *http.Client
}

// NewAmplitude returns a notifier that posts events to Amplitude.
func NewAmplitude(apiKey, apiURL string) *AmplitudeNotifier {
	if apiURL == "" {
		apiURL = amplitudeDefaultAPIURL
	}
	return &AmplitudeNotifier{
		apiKey: apiKey,
		apiURL: apiURL,
		client: &http.Client{Timeout: amplitudeDefaultTimeout},
	}
}

// Send dispatches the alert to Amplitude.
func (a *AmplitudeNotifier) Send(al alert.Alert) error {
	if al.Timestamp.IsZero() {
		al.Timestamp = time.Now()
	}

	eventType := "port_opened"
	if al.Level == alert.Info {
		eventType = "port_closed"
	}

	payload := amplitudePayload{
		APIKey: a.apiKey,
		Events: []amplitudeEvent{
			{
				UserID:    "portwatch",
				EventType: eventType,
				EventTime: al.Timestamp.UnixMilli(),
				EventProp: map[string]interface{}{
					"port":    al.Port,
					"proto":   al.Proto,
					"message": al.Message,
					"level":   al.Level.String(),
				},
			},
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("amplitude: marshal payload: %w", err)
	}

	resp, err := a.client.Post(a.apiURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("amplitude: post event: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("amplitude: unexpected status %d", resp.StatusCode)
	}
	return nil
}
