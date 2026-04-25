package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/user/portwatch/internal/alert"
)

// jiraPayload is the request body sent to the Jira REST API.
type jiraPayload struct {
	Fields jiraFields `json:"fields"`
}

type jiraFields struct {
	Project   jiraKey `json:"project"`
	IssueType jiraKey `json:"issuetype"`
	Summary   string  `json:"summary"`
	Description string `json:"description"`
	Priority  jiraKey `json:"priority"`
}

type jiraKey struct {
	Key string `json:"key"`
}

// jiraNotifier posts a Jira issue for each alert.
type jiraNotifier struct {
	baseURL    string
	projectKey string
	username   string
	token      string
	client     *http.Client
}

// NewJira returns a Notifier that creates Jira issues.
func NewJira(baseURL, projectKey, username, token string) Notifier {
	return &jiraNotifier{
		baseURL:    baseURL,
		projectKey: projectKey,
		username:   username,
		token:      token,
		client:     &http.Client{Timeout: 10 * time.Second},
	}
}

func (j *jiraNotifier) Send(a alert.Alert) error {
	priority := "High"
	if a.Level == alert.Info {
		priority = "Low"
	}

	body := jiraPayload{
		Fields: jiraFields{
			Project:     jiraKey{Key: j.projectKey},
			IssueType:   jiraKey{Key: "Bug"},
			Summary:     a.Title,
			Description: fmt.Sprintf("%s\n\nPort: %d/%s\nHost: %s", a.Message, a.Port, a.Proto, a.Host),
			Priority:    jiraKey{Key: priority},
		},
	}

	raw, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("jira: marshal: %w", err)
	}

	url := j.baseURL + "/rest/api/2/issue"
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(raw))
	if err != nil {
		return fmt.Errorf("jira: new request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(j.username, j.token)

	resp, err := j.client.Do(req)
	if err != nil {
		return fmt.Errorf("jira: do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("jira: unexpected status %d", resp.StatusCode)
	}
	return nil
}
