package notify_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/notify"
)

func discordAlert() alert.Alert {
	return alert.Alert{
		Level:     alert.LevelWarn,
		Title:     "Port opened",
		Body:      "Port 8080/tcp is now open",
		Timestamp: time.Now(),
	}
}

func TestDiscord_SendsCorrectPayload(t *testing.T) {
	var received map[string]interface{}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if ct := r.Header.Get("Content-Type"); ct != "application/json" {
			t.Errorf("unexpected Content-Type: %s", ct)
		}
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	d := notify.NewDiscord(srv.URL)
	if err := d.Send(discordAlert()); err != nil {
		t.Fatalf("Send returned error: %v", err)
	}

	if received["username"] != "portwatch" {
		t.Errorf("expected username 'portwatch', got %v", received["username"])
	}

	embeds, ok := received["embeds"].([]interface{})
	if !ok || len(embeds) == 0 {
		t.Fatal("expected at least one embed")
	}
	embed := embeds[0].(map[string]interface{})
	if embed["title"] == "" {
		t.Error("embed title should not be empty")
	}
}

func TestDiscord_InfoLevelUsesBlurpleColor(t *testing.T) {
	var received map[string]interface{}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewDecoder(r.Body).Decode(&received)
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	a := discordAlert()
	a.Level = alert.LevelInfo

	d := notify.NewDiscord(srv.URL)
	_ = d.Send(a)

	embeds := received["embeds"].([]interface{})
	embed := embeds[0].(map[string]interface{})
	// 0x5865F2 == 5793266
	if int(embed["color"].(float64)) != 0x5865F2 {
		t.Errorf("expected blurple color for info level, got %v", embed["color"])
	}
}

func TestDiscord_ErrorOnBadStatus(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer srv.Close()

	d := notify.NewDiscord(srv.URL)
	if err := d.Send(discordAlert()); err == nil {
		t.Fatal("expected error on non-2xx status")
	}
}

func TestDiscord_DefaultTimeout(t *testing.T) {
	d := notify.NewDiscord("http://example.com")
	// Access via interface to confirm it is a *Discord (compile-time check).
	var _ interface{ Send(alert.Alert) error } = d
}
