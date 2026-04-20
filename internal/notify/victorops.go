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
	defaultVictorOpsTimeout = 10 * time.Second
	victorOpsCritical       = "CRITICAL"
	victorOpsWarning        = "WARNING"
	victorOpsInfo           = "INFO"
)

// VictorOps sends alert notifications to a VictorOps (Splunk On-Call) REST
// endpoint. Each alert is translated into a VictorOps message type based on
// the alert level: WARN maps to CRITICAL, INFO maps to INFO.
type VictorOps struct {
	endpointURL string
	routingKey  string
	client      *http.Client
}

// victorOpsPayload mirrors the VictorOps REST endpoint payload schema.
type victorOpsPayload struct {
	MessageType       string `json:"message_type"`
	EntityID          string `json:"entity_id"`
	EntityDisplayName string `json:"entity_display_name"`
	StateMessage      string `json:"state_message"`
	Timestamp         int64  `json:"timestamp"`
	MonitoringTool    string `json:"monitoring_tool"`
	Port              int    `json:"port"`
	Proto             string `json:"proto"`
}

// NewVictorOps returns a VictorOps notifier that posts to the given REST
// endpoint URL (which should already include the routing key path segment,
// e.g. https://alert.victorops.com/integrations/generic/.../alert/<key>/<routing>).
// Pass a nil client to use the default HTTP client with a sensible timeout.
func NewVictorOps(endpointURL string, client *http.Client) *VictorOps {
	if client == nil {
		client = &http.Client{Timeout: defaultVictorOpsTimeout}
	}
	return &VictorOps{
		endpointURL: endpointURL,
		client:      client,
	}
}

// Send delivers the alert to VictorOps. It returns an error if the HTTP
// request fails or the server responds with a non-2xx status code.
func (v *VictorOps) Send(a alert.Alert) error {
	if a.Timestamp.IsZero() {
		a.Timestamp = time.Now()
	}

	payload := victorOpsPayload{
		MessageType:       v.messageType(a),
		EntityID:          fmt.Sprintf("portwatch-%s-%d", a.Proto, a.Port),
		EntityDisplayName: fmt.Sprintf("Port %s/%d %s", a.Proto, a.Port, a.Message),
		StateMessage:      a.Message,
		Timestamp:         a.Timestamp.Unix(),
		MonitoringTool:    "portwatch",
		Port:              a.Port,
		Proto:             a.Proto,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("victorops: marshal payload: %w", err)
	}

	resp, err := v.client.Post(v.endpointURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("victorops: http post: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("victorops: unexpected status %d", resp.StatusCode)
	}
	return nil
}

// messageType maps an alert level to a VictorOps message type string.
func (v *VictorOps) messageType(a alert.Alert) string {
	switch a.Level {
	case alert.LevelWarn:
		return victorOpsCritical
	default:
		return victorOpsInfo
	}
}
