// Package notify provides notifier implementations for delivering portwatch
// alerts to various destinations.
//
// SlackNotifier posts alert messages to a Slack incoming webhook URL.
// Create one with NewSlack, then pass it to a Multi notifier or use it
// directly as a Notifier.
//
// Example:
//
//	sn := notify.NewSlack(os.Getenv("SLACK_WEBHOOK"), 5*time.Second)
//	sn.Send(alert.Alert{Level: "WARN", Message: "port opened", Port: 443, Proto: "tcp"})
package notify
