package notify

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/yourorg/portwatch/internal/alert"
)

const twilioDefaultTimeout = 10 * time.Second

// TwilioNotifier sends SMS alerts via the Twilio REST API.
type TwilioNotifier struct {
	accountSID string
	authToken  string
	from       string
	to         string
	client     *http.Client
}

// NewTwilio constructs a TwilioNotifier.
func NewTwilio(accountSID, authToken, from, to string) *TwilioNotifier {
	return &TwilioNotifier{
		accountSID: accountSID,
		authToken:  authToken,
		from:       from,
		to:         to,
		client:     &http.Client{Timeout: twilioDefaultTimeout},
	}
}

// Send dispatches an SMS for the given alert.
func (t *TwilioNotifier) Send(a alert.Alert) error {
	body := url.Values{}
	body.Set("From", t.from)
	body.Set("To", t.to)
	body.Set("Body", fmt.Sprintf("[portwatch] %s port %d/%s", a.Level, a.Port, a.Proto))

	apiURL := fmt.Sprintf("https://api.twilio.com/2010-04-01/Accounts/%s/Messages.json", t.accountSID)
	req, err := http.NewRequest(http.MethodPost, apiURL, strings.NewReader(body.Encode()))
	if err != nil {
		return fmt.Errorf("twilio: build request: %w", err)
	}
	req.SetBasicAuth(t.accountSID, t.authToken)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := t.client.Do(req)
	if err != nil {
		return fmt.Errorf("twilio: send: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var e struct {
			Message string `json:"message"`
		}
		_ = json.NewDecoder(resp.Body).Decode(&e)
		return fmt.Errorf("twilio: unexpected status %d: %s", resp.StatusCode, e.Message)
	}
	return nil
}
