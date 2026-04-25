// Package notify provides notifier implementations for portwatch alerts.
//
// # Jira Notifier
//
// The Jira notifier creates a Jira issue via the REST API (v2) whenever
// portwatch detects an unexpected port change.
//
// Warn-level alerts (port opened) are filed with High priority; Info-level
// alerts (port closed) use Low priority. The issue type is always "Bug".
//
// Configuration fields:
//
//	[jira]
//	enabled     = true
//	base_url    = "https://jira.example.com"   # no trailing slash
//	project_key = "OPS"
//	username    = "bot@example.com"
//	token       = "your-api-token"
package notify
