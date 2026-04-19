package alert

import (
	"fmt"
	"io"
	"os"
	"time"
)

// Level represents the severity of an alert.
type Level string

const (
	LevelInfo  Level = "INFO"
	LevelWarn  Level = "WARN"
	LevelError Level = "ERROR"
)

// Alert holds information about a port state change.
type Alert struct {
	Timestamp time.Time
	Level     Level
	Port      int
	Message   string
}

// Notifier sends alerts to a destination.
type Notifier struct {
	out io.Writer
}

// New creates a Notifier that writes to the given writer.
// If w is nil, os.Stdout is used.
func New(w io.Writer) *Notifier {
	if w == nil {
		w = os.Stdout
	}
	return &Notifier{out: w}
}

// Send formats and writes an alert.
func (n *Notifier) Send(a Alert) error {
	if a.Timestamp.IsZero() {
		a.Timestamp = time.Now()
	}
	_, err := fmt.Fprintf(
		n.out,
		"[%s] %s port=%d msg=%q\n",
		a.Timestamp.Format(time.RFC3339),
		a.Level,
		a.Port,
		a.Message,
	)
	return err
}

// PortOpened returns an Alert indicating a new open port was detected.
func PortOpened(port int) Alert {
	return Alert{
		Timestamp: time.Now(),
		Level:     LevelWarn,
		Port:      port,
		Message:   "unexpected port opened",
	}
}

// PortClosed returns an Alert indicating a previously open port was closed.
func PortClosed(port int) Alert {
	return Alert{
		Timestamp: time.Now(),
		Level:     LevelInfo,
		Port:      port,
		Message:   "port closed",
	}
}
