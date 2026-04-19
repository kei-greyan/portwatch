// Package notify provides notification backends for portwatch alerts.
package notify

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/user/portwatch/internal/alert"
)

// Notifier sends an alert to a destination.
type Notifier interface {
	Notify(a alert.Alert) error
}

// LogNotifier writes alerts to an io.Writer in a human-readable format.
type LogNotifier struct {
	w io.Writer
}

// New returns a LogNotifier that writes to w.
// If w is nil, os.Stdout is used.
func New(w io.Writer) *LogNotifier {
	if w == nil {
		w = os.Stdout
	}
	return &LogNotifier{w: w}
}

// Notify formats and writes the alert.
func (n *LogNotifier) Notify(a alert.Alert) error {
	ts := a.Timestamp
	if ts.IsZero() {
		ts = time.Now().UTC()
	}
	_, err := fmt.Fprintf(n.w, "%s [%s] port %d/%-3s — %s\n",
		ts.Format(time.RFC3339),
		a.Level,
		a.Port,
		a.Proto,
		a.Message,
	)
	return err
}
