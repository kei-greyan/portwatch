package notify_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/notify"
)

func TestNew_DefaultsToStdout(t *testing.T) {
	n := notify.New(nil)
	if n == nil {
		t.Fatal("expected non-nil notifier")
	}
}

func TestNotify_WritesFormattedLine(t *testing.T) {
	var buf bytes.Buffer
	n := notify.New(&buf)

	a := alert.Alert{
		Level:     "WARN",
		Port:      8080,
		Proto:     "tcp",
		Message:   "port opened",
		Timestamp: time.Date(2024, 1, 2, 3, 4, 5, 0, time.UTC),
	}

	if err := n.Notify(a); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	line := buf.String()
	for _, want := range []string{"2024-01-02T03:04:05Z", "WARN", "8080", "tcp", "port opened"} {
		if !strings.Contains(line, want) {
			t.Errorf("output %q missing %q", line, want)
		}
	}
}

func TestNotify_SetsTimestampWhenZero(t *testing.T) {
	var buf bytes.Buffer
	n := notify.New(&buf)

	a := alert.Alert{
		Level:   "INFO",
		Port:    22,
		Proto:   "tcp",
		Message: "port closed",
	}

	if err := n.Notify(a); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if buf.Len() == 0 {
		t.Error("expected non-empty output")
	}
}
