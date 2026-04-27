package notify

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/example/portwatch/internal/alert"
)

const (
	prowlDefaultAPIURL = "https://api.prowlapp.com/publicapi/add"
	prowlDefaultTimeout = 10 * time.Second
)

// prowl sends push notifications to iOS devices via the Prowl API.
type prowl struct {
	apiKey  string
	apiURL  string
	appName string
	client  *http.Client
}

// NewProwl returns a Notifier that delivers alerts via the Prowl push
// notification service. apiKey must be a valid Prowl API key.
func NewProwl(apiKey, appName string) Notifier {
	return &prowl{
		apiKey:  apiKey,
		apiURL:  prowlDefaultAPIURL,
		appName: appName,
		client:  &http.Client{Timeout: prowlDefaultTimeout},
	}
}

func (p *prowl) Send(a alert.Alert) error {
	priority := prowlPriority(a)

	params := url.Values{}
	params.Set("apikey", p.apiKey)
	params.Set("application", p.appName)
	params.Set("event", a.Title)
	params.Set("description", a.Message)
	params.Set("priority", fmt.Sprintf("%d", priority))

	resp, err := p.client.PostForm(p.apiURL, params)
	if err != nil {
		return fmt.Errorf("prowl: post: %w", err)
	}
	defer resp.Body.Close()
	io.Copy(io.Discard, resp.Body)

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("prowl: unexpected status %d", resp.StatusCode)
	}
	return nil
}

// prowlPriority maps alert levels to Prowl priority values.
// Prowl priorities: -2 (very low) to 2 (emergency).
func prowlPriority(a alert.Alert) int {
	switch a.Level {
	case alert.LevelWarn:
		return 1
	case alert.LevelInfo:
		return 0
	default:
		return 0
	}
}
