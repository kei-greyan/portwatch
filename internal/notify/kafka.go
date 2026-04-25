package notify

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/user/portwatch/internal/alert"
)

// kafkaPayload is the JSON structure sent to a Kafka REST Proxy endpoint.
type kafkaPayload struct {
	Records []kafkaRecord `json:"records"`
}

type kafkaRecord struct {
	Value kafkaValue `json:"value"`
}

type kafkaValue struct {
	Timestamp time.Time `json:"timestamp"`
	Level     string    `json:"level"`
	Port      int       `json:"port"`
	Proto     string    `json:"proto"`
	Message   string    `json:"message"`
}

// kafkaNotifier sends alerts to a Kafka REST Proxy.
type kafkaNotifier struct {
	proxyURL string
	topic    string
	client   *http.Client
}

// NewKafka returns a Notifier that publishes alerts via the Kafka REST Proxy.
func NewKafka(proxyURL, topic string) Notifier {
	return &kafkaNotifier{
		proxyURL: strings.TrimRight(proxyURL, "/"),
		topic:    topic,
		client: &http.Client{
			Timeout: 10 * time.Second,
			Transport: &http.Transport{
				DialContext: (&net.Dialer{Timeout: 5 * time.Second}).DialContext,
			},
		},
	}
}

func (k *kafkaNotifier) Send(ctx context.Context, a alert.Alert) error {
	if a.Timestamp.IsZero() {
		a.Timestamp = time.Now().UTC()
	}

	payload := kafkaPayload{
		Records: []kafkaRecord{
			{
				Value: kafkaValue{
					Timestamp: a.Timestamp,
					Level:     a.Level,
					Port:      a.Port,
					Proto:     a.Proto,
					Message:   a.Message,
				},
			},
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("kafka: marshal payload: %w", err)
	}

	url := fmt.Sprintf("%s/topics/%s", k.proxyURL, k.topic)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, strings.NewReader(string(body)))
	if err != nil {
		return fmt.Errorf("kafka: build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/vnd.kafka.json.v2+json")

	resp, err := k.client.Do(req)
	if err != nil {
		return fmt.Errorf("kafka: send: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("kafka: unexpected status %d", resp.StatusCode)
	}
	return nil
}
