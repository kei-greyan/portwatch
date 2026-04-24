// Package notify provides notification backends for portwatch alerts.
//
// PagerDuty notifier sends alerts to the PagerDuty Events API v2.
// It maps alert severity (warn → critical, info → info) to PagerDuty
// severity levels and uses the port number as the dedup key so that
// repeated alerts for the same port are deduplicated automatically.
package notify
