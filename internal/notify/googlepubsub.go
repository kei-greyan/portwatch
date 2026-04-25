package notify

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/user/portwatch/internal/alert"
)

const (
	googlePubSubDefaultTimeout = 10 * time.Second
	googlePubSubAPIBase        = "https://pubsub.googleapis.com/v1"
)

// GooglePubSub sends alert notifications to a Google Cloud Pub/Sub topic
// using the REST API with an API key or service account bearer token.
type GooglePubSub struct {
	projectID string
	topicID   string
	token     string
	client    *http.Client
}

type pubSubMessage struct {
	Data       string            `json:"data"`
	Attributes map[string]string `json:"attributes,omitempty"`
}

type pubSubPublishRequest struct {
	Messages []pubSubMessage `json:"messages"`
}

// NewGooglePubSub creates a Notifier that publishes alerts to the given
// Google Cloud Pub/Sub topic. token is a Bearer token (e.g. from a service
// account access token or gcloud auth print-access-token).
func NewGooglePubSub(projectID, topicID, token string, opts ...func(*GooglePubSub)) *GooglePubSub {
	g := &GooglePubSub{
		projectID: projectID,
		topicID:   topicID,
		token:     token,
		client:    &http.Client{Timeout: googlePubSubDefaultTimeout},
	}
	for _, o := range opts {
		o(g)
	}
	return g
}

// WithGooglePubSubHTTPClient overrides the default HTTP client.
func WithGooglePubSubHTTPClient(c *http.Client) func(*GooglePubSub) {
	return func(g *GooglePubSub) { g.client = c }
}

// Send publishes the alert as a base64-encoded JSON message to Pub/Sub.
func (g *GooglePubSub) Send(ctx context.Context, a alert.Alert) error {
	if a.Timestamp.IsZero() {
		a.Timestamp = time.Now()
	}

	payload, err := json.Marshal(a)
	if err != nil {
		return fmt.Errorf("googlepubsub: marshal alert: %w", err)
	}

	encoded := base64.StdEncoding.EncodeToString(payload)

	body, err := json.Marshal(pubSubPublishRequest{
		Messages: []pubSubMessage{
			{
				Data: encoded,
				Attributes: map[string]string{
					"level": a.Level,
					"port":  fmt.Sprintf("%d", a.Port),
					"proto": a.Proto,
				},
			},
		},
	})
	if err != nil {
		return fmt.Errorf("googlepubsub: marshal request: %w", err)
	}

	url := fmt.Sprintf("%s/projects/%s/topics/%s:publish",
		googlePubSubAPIBase, g.projectID, g.topicID)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("googlepubsub: build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+g.token)

	resp, err := g.client.Do(req)
	if err != nil {
		return fmt.Errorf("googlepubsub: send: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("googlepubsub: unexpected status %d", resp.StatusCode)
	}
	return nil
}
