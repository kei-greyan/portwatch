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
		datadogDefaultURL     = "https://api.datadoghq.com/api/v1/events"
		datadogDefaultTimeout = 10 * time.Second
)

// DataDog sends alerts to the Datadog Events API.
type DataDog struct {
	apiKey  string
	apiURL  string
	client  *http.Client
}

type datadogPayload struct {
	Title      string   `json:"title"`
	Text       string   `json:"text"`
	AlertType  string   `json:"alert_type"`
	Tags       []string `json:"tags,omitempty"`
	SourceType string   `json:"source_type_name"`
}

// NewDataDog creates a Datadog notifier using the given API key.
// apiURL may be empty to use the default Datadog endpoint.
func NewDataDog(apiKey, apiURL string) *DataDog {
	if apiURL == "" {
		apiURL = datadogDefaultURL
	}
	return &DataDog{
		apiKey: apiKey,
		apiURL: apiURL,
		client: &http.Client{Timeout: datadogDefaultTimeout},
	}
}

// Send delivers an alert to Datadog as an event.
func (d *DataDog) Send(a alert.Alert) error {
	alertType := "info"
	if a.Level == alert.Warn {
		alertType = "warning"
	}

	payload := datadogPayload{
		Title:      fmt.Sprintf("portwatch: %s", a.Message),
		Text:       fmt.Sprintf("Port %d — %s", a.Port, a.Message),
		AlertType:  alertType,
		Tags:       []string{fmt.Sprintf("port:%d", a.Port), "source:portwatch"},
		SourceType: "portwatch",
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("datadog: marshal payload: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, d.apiURL, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("datadog: create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("DD-API-KEY", d.apiKey)

	resp, err := d.client.Do(req)
	if err != nil {
		return fmt.Errorf("datadog: send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("datadog: unexpected status %d", resp.StatusCode)
	}
	return nil
}
