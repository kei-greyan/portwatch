package notify

import (
	"fmt"
	"log/syslog"

	"github.com/user/portwatch/internal/alert"
)

// Syslog sends alerts to the local or remote syslog daemon.
type Syslog struct {
	writer *syslog.Writer
	tag    string
}

// NewSyslog creates a Syslog notifier. network and addr may be empty strings
// to use the local Unix socket. tag is the syslog program identifier.
func NewSyslog(network, addr, tag string) (*Syslog, error) {
	if tag == "" {
		tag = "portwatch"
	}
	w, err := syslog.Dial(network, addr, syslog.LOG_DAEMON|syslog.LOG_NOTICE, tag)
	if err != nil {
		return nil, fmt.Errorf("syslog dial: %w", err)
	}
	return &Syslog{writer: w, tag: tag}, nil
}

// Send writes the alert to syslog using a severity that matches the alert level.
func (s *Syslog) Send(a alert.Alert) error {
	msg := fmt.Sprintf("port=%d proto=%s message=%q", a.Port, a.Proto, a.Message)
	switch a.Level {
	case alert.LevelWarn:
		return s.writer.Warning(msg)
	case alert.LevelError:
		return s.writer.Err(msg)
	default:
		return s.writer.Info(msg)
	}
}

// Close releases the underlying syslog connection.
func (s *Syslog) Close() error {
	return s.writer.Close()
}
