package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/user/portwatch/internal/alert"
)

const defaultWhatsAppTimeout = 10 * time.Second

// whatsAppNotifier sends alerts via the WhatsApp Business Cloud API.
type whatsAppNotifier struct {
	phoneNumberID string
	token         string
	recipient     string
	apiURL        string
	client        *http.Client
}

type whatsAppPayload struct {
	MessagingProduct string            `json:"messaging_product"`
	To               string            `json:"to"`
	Type             string            `json:"type"`
	Text             whatsAppText      `json:"text"`
}

type whatsAppText struct {
	Body string `json:"body"`
}

// NewWhatsApp returns a Notifier that delivers alerts via WhatsApp Business Cloud API.
func NewWhatsApp(phoneNumberID, token, recipient string) Notifier {
	return &whatsAppNotifier{
		phoneNumberID: phoneNumberID,
		token:         token,
		recipient:     recipient,
		apiURL:        fmt.Sprintf("https://graph.facebook.com/v18.0/%s/messages", phoneNumberID),
		client:        &http.Client{Timeout: defaultWhatsAppTimeout},
	}
}

func (w *whatsAppNotifier) Send(a alert.Alert) error {
	body := fmt.Sprintf("[portwatch] %s — port %d/%s on %s",
		a.Message, a.Port, a.Proto, a.Host)

	payload := whatsAppPayload{
		MessagingProduct: "whatsapp",
		To:               w.recipient,
		Type:             "text",
		Text:             whatsAppText{Body: body},
	}

	b, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("whatsapp: marshal payload: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, w.apiURL, bytes.NewReader(b))
	if err != nil {
		return fmt.Errorf("whatsapp: create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+w.token)

	resp, err := w.client.Do(req)
	if err != nil {
		return fmt.Errorf("whatsapp: send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("whatsapp: unexpected status %d", resp.StatusCode)
	}
	return nil
}
