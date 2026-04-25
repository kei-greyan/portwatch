package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/yourorg/portwatch/internal/alert"
)

const defaultSignalRTimeout = 10 * time.Second

// SignalR sends alerts to an Azure SignalR Service REST endpoint.
type SignalR struct {
	hub    string
	apiKey string
	client *http.Client
}

type signalRPayload struct {
	Target    string        `json:"target"`
	Arguments []interface{} `json:"arguments"`
}

type signalRMessage struct {
	Level   string    `json:"level"`
	Port    int       `json:"port"`
	Proto   string    `json:"proto"`
	Message string    `json:"message"`
	Time    time.Time `json:"time"`
}

// NewSignalR returns a Notifier that broadcasts alerts via Azure SignalR.
// hub is the full REST broadcast URL, e.g.
// https://<name>.service.signalr.net/api/v1/hubs/<hubname>.
func NewSignalR(hub, apiKey string, opts ...func(*http.Client)) *SignalR {
	c := &http.Client{Timeout: defaultSignalRTimeout}
	for _, o := range opts {
		o(c)
	}
	return &SignalR{hub: hub, apiKey: apiKey, client: c}
}

// Send broadcasts the alert to all connected SignalR clients.
func (s *SignalR) Send(a alert.Alert) error {
	if a.Time.IsZero() {
		a.Time = time.Now()
	}
	msg := signalRMessage{
		Level:   a.Level,
		Port:    a.Port,
		Proto:   a.Proto,
		Message: a.Message,
		Time:    a.Time,
	}
	payload := signalRPayload{
		Target:    "portAlert",
		Arguments: []interface{}{msg},
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("signalr: marshal: %w", err)
	}
	req, err := http.NewRequest(http.MethodPost, s.hub, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("signalr: build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	if s.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+s.apiKey)
	}
	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("signalr: send: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return fmt.Errorf("signalr: unexpected status %d", resp.StatusCode)
	}
	return nil
}
