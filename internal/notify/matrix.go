package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/patrickdappollonio/portwatch/internal/alert"
)

const defaultMatrixTimeout = 10 * time.Second

// matrixPayload is the request body sent to the Matrix client-server API.
type matrixPayload struct {
	MsgType       string `json:"msgtype"`
	Body          string `json:"body"`
	FormattedBody string `json:"formatted_body,omitempty"`
	Format        string `json:"format,omitempty"`
}

// Matrix sends alert notifications to a Matrix room via the client-server API.
type Matrix struct {
	homeserver string
	roomID     string
	token      string
	client     *http.Client
}

// NewMatrix creates a Matrix notifier. homeserver is the base URL (e.g.
// "https://matrix.example.com"), roomID is the full room ID (e.g.
// "!abc123:example.com"), and token is the access token.
func NewMatrix(homeserver, roomID, token string) *Matrix {
	return &Matrix{
		homeserver: homeserver,
		roomID:     roomID,
		token:      token,
		client:     &http.Client{Timeout: defaultMatrixTimeout},
	}
}

// Send dispatches an alert to the configured Matrix room.
func (m *Matrix) Send(a alert.Alert) error {
	body := matrixPayload{
		MsgType:       "m.text",
		Body:          fmt.Sprintf("[%s] %s", a.Level, a.Message),
		Format:        "org.matrix.custom.html",
		FormattedBody: fmt.Sprintf("<b>[%s]</b> %s", a.Level, a.Message),
	}

	raw, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("matrix: marshal payload: %w", err)
	}

	url := fmt.Sprintf("%s/_matrix/client/v3/rooms/%s/send/m.room.message",
		m.homeserver, m.roomID)

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(raw))
	if err != nil {
		return fmt.Errorf("matrix: build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+m.token)

	resp, err := m.client.Do(req)
	if err != nil {
		return fmt.Errorf("matrix: send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("matrix: unexpected status %d", resp.StatusCode)
	}
	return nil
}
