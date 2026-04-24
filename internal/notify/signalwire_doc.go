// Package notify provides notification backends for portwatch alerts.
//
// SignalWire notifier sends SMS alerts via the SignalWire REST API,
// which is compatible with the Twilio API surface. Configure it with
// your SignalWire Space URL, Project ID, API token, and the from/to
// phone numbers.
//
// Example configuration:
//
//	notify.NewSignalWire(
//		"https://example.signalwire.com",
//		"project-id",
//		"api-token",
//		"+15550001111", // from
//		"+15559998888", // to
//	)
package notify
