// Package notify provides Notifier implementations for various alerting
// backends.
//
// # Google Chat
//
// NewGoogleChat creates a notifier that sends alerts to a Google Chat space
// via an incoming webhook URL. To obtain a webhook URL:
//
//  1. Open the Google Chat space.
//  2. Click the space name → Manage webhooks.
//  3. Add a webhook and copy the URL.
//
// Example configuration:
//
//	googlechat:
//	  enabled: true
//	  webhook_url: "https://chat.googleapis.com/v1/spaces/.../messages?key=...&token=..."
package notify
