package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/user/portwatch/internal/alert"
)

// teamsPayload is the adaptive card message format for Microsoft Teams
// incoming webhooks (Office 365 Connector card schema).
type teamsPayload struct {
	Type       string         `json:"@type"`
	Context    string         `json:"@context"`
	ThemeColor string         `json:"themeColor"`
	Summary    string         `json:"summary"`
	Sections   []teamsSection `json:"sections"`
}

type teamsSection struct {
	ActivityTitle string       `json:"activityTitle"`
	Facts         []teamsFact  `json:"facts"`
	Markdown      bool         `json:"markdown"`
}

type teamsFact struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// Teams delivers alerts to a Microsoft Teams channel via an incoming webhook URL.
type Teams struct {
	webhookURL string
	client     *http.Client
}

// NewTeams returns a Teams notifier that posts to the given webhook URL.
func NewTeams(webhookURL string) *Teams {
	return &Teams{
		webhookURL: webhookURL,
		client:     &http.Client{Timeout: 10 * time.Second},
	}
}

// Send dispatches the alert to the configured Teams webhook.
func (t *Teams) Send(a alert.Alert) error {
	color := "2DC72D" // green – info
	if a.Level == alert.Warn {
		color = "FFA500" // orange – warning
	}

	payload := teamsPayload{
		Type:       "MessageCard",
		Context:    "http://schema.org/extensions",
		ThemeColor: color,
		Summary:    a.Message,
		Sections: []teamsSection{
			{
				ActivityTitle: fmt.Sprintf("portwatch – %s", a.Message),
				Facts: []teamsFact{
					{Name: "Port", Value: fmt.Sprintf("%d", a.Port)},
					{Name: "Level", Value: string(a.Level)},
					{Name: "Time", Value: a.Timestamp.UTC().Format(time.RFC3339)},
				},
				Markdown: false,
			},
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("teams: marshal payload: %w", err)
	}

	resp, err := t.client.Post(t.webhookURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("teams: http post: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("teams: unexpected status %d", resp.StatusCode)
	}
	return nil
}
