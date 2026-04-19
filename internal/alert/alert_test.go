package alert_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/user/portwatch/internal/alert"
)

func TestSend_WritesFormattedAlert(t *testing.T) {
	var buf bytes.Buffer
	n := alert.New(&buf)

	a := alert.Alert{
		Timestamp: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		Level:     alert.LevelWarn,
		Port:      8080,
		Message:   "unexpected port opened",
	}

	if err := n.Send(a); err != nil {
		t.Fatalf("Send returned error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "WARN") {
		t.Errorf("expected WARN in output, got: %s", out)
	}
	if !strings.Contains(out, "port=8080") {
		t.Errorf("expected port=8080 in output, got: %s", out)
	}
	if !strings.Contains(out, "unexpected port opened") {
		t.Errorf("expected message in output, got: %s", out)
	}
}

func TestSend_SetsTimestampIfZero(t *testing.T) {
	var buf bytes.Buffer
	n := alert.New(&buf)

	a := alert.Alert{Level: alert.LevelInfo, Port: 22, Message: "port closed"}
	if err := n.Send(a); err != nil {
		t.Fatalf("Send returned error: %v", err)
	}

	if buf.Len() == 0 {
		t.Error("expected non-empty output")
	}
}

func TestPortOpened_ReturnsWarnAlert(t *testing.T) {
	a := alert.PortOpened(9000)
	if a.Level != alert.LevelWarn {
		t.Errorf("expected LevelWarn, got %s", a.Level)
	}
	if a.Port != 9000 {
		t.Errorf("expected port 9000, got %d", a.Port)
	}
}

func TestPortClosed_ReturnsInfoAlert(t *testing.T) {
	a := alert.PortClosed(9000)
	if a.Level != alert.LevelInfo {
		t.Errorf("expected LevelInfo, got %s", a.Level)
	}
}
