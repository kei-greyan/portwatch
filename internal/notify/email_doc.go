// Package notify provides notification backends for portwatch alerts.
//
// # Email notifier
//
// NewEmail creates a notifier that delivers alerts over SMTP.
// Configure it with an EmailConfig specifying the SMTP host, port,
// optional credentials, sender address, and one or more recipients.
//
// Example:
//
//	n := notify.NewEmail(notify.EmailConfig{
//		Host: "smtp.example.com",
//		Port: 587,
//		Username: "user",
//		Password: "secret",
//		From:     "portwatch@example.com",
//		To:       []string{"ops@example.com"},
//	})
package notify
