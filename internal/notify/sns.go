package notify

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/user/portwatch/internal/alert"
)

// SNSNotifier publishes alerts to an AWS SNS-compatible HTTP endpoint.
// It is compatible with both real AWS SNS and local mock servers.
type SNSNotifier struct {
	topicARN string
	endpoint  string
	client    *http.Client
}

type snsPayload struct {
	TopicARN string `json:"TopicARN"`
	Subject  string `json:"Subject"`
	Message  string `json:"Message"`
}

// NewSNS creates a new SNSNotifier.
// endpoint is the HTTP URL of the SNS endpoint (e.g. https://sns.us-east-1.amazonaws.com).
// topicARN is the ARN of the target topic.
func NewSNS(endpoint, topicARN string, opts ...func(*SNSNotifier)) *SNSNotifier {
	s := &SNSNotifier{
		topicARN: topicARN,
		endpoint:  strings.TrimRight(endpoint, "/"),
		client:    &http.Client{Timeout: 10 * time.Second},
	}
	for _, o := range opts {
		o(s)
	}
	return s
}

// WithSNSHTTPClient overrides the default HTTP client.
func WithSNSHTTPClient(c *http.Client) func(*SNSNotifier) {
	return func(s *SNSNotifier) { s.client = c }
}

// Send publishes the alert to the configured SNS topic.
func (s *SNSNotifier) Send(ctx context.Context, a alert.Alert) error {
	subject := fmt.Sprintf("[portwatch] %s port %d", strings.ToUpper(a.Level), a.Port)
	body, err := json.Marshal(snsPayload{
		TopicARN: s.topicARN,
		Subject:  subject,
		Message:  a.Message,
	})
	if err != nil {
		return fmt.Errorf("sns: marshal payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.endpoint, strings.NewReader(string(body)))
	if err != nil {
		return fmt.Errorf("sns: build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("sns: send: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("sns: unexpected status %d", resp.StatusCode)
	}
	return nil
}
