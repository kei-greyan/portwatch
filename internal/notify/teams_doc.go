// Package notify provides notification backends for portwatch alerts.
//
// # Microsoft Teams
//
// The Teams notifier delivers alerts to a Microsoft Teams channel using
// an Office 365 Connector (incoming webhook) URL.
//
// Configure an incoming webhook in Teams and pass the URL to NewTeams:
//
//	n := notify.NewTeams("https://outlook.office.com/webhook/...")
//	err := n.Send(a)
//
// Alert levels are reflected as card theme colours:
//   - Warn  → orange (#FFA500)
//   - Info  → green  (#2DC72D)
package notify
