// Package notify provides notifier implementations for dispatching alerts
// to various destinations.
//
// # PagerDuty
//
// The PagerDuty notifier sends alerts using the PagerDuty Events API v2
// (https://developer.pagerduty.com/docs/ZG9jOjExMDI5NTgw-send-an-alert-event).
//
// Usage:
//
//	pd := notify.NewPagerDuty("your-integration-routing-key")
//	err := pd.Send(a)
//
// Warn-level alerts are mapped to PagerDuty severity "error" to trigger
// an incident. Info-level alerts use severity "info" and will not page
// on-call responders by default.
package notify
