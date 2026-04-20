package notify

import (
	"fmt"
	"net/smtp"
	"strings"
	"time"

	"github.com/user/portwatch/internal/alert"
)

// EmailConfig holds SMTP configuration for email notifications.
type EmailConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
	To       []string
}

type emailNotifier struct {
	cfg EmailConfig
}

// NewEmail returns a Notifier that sends alerts via SMTP.
func NewEmail(cfg EmailConfig) Notifier {
	return &emailNotifier{cfg: cfg}
}

func (e *emailNotifier) Send(a alert.Alert) error {
	if a.Timestamp.IsZero() {
		a.Timestamp = time.Now()
	}

	subject := fmt.Sprintf("[portwatch] %s port %d", a.Level, a.Port)
	body := fmt.Sprintf("To: %s\r\nFrom: %s\r\nSubject: %s\r\n\r\n%s\n\nPort:    %d\nTime:    %s\n",
		strings.Join(e.cfg.To, ", "),
		e.cfg.From,
		subject,
		a.Message,
		a.Port,
		a.Timestamp.Format(time.RFC3339),
	)

	addr := fmt.Sprintf("%s:%d", e.cfg.Host, e.cfg.Port)
	var auth smtp.Auth
	if e.cfg.Username != "" {
		auth = smtp.PlainAuth("", e.cfg.Username, e.cfg.Password, e.cfg.Host)
	}

	return smtp.SendMail(addr, auth, e.cfg.From, e.cfg.To, []byte(body))
}
