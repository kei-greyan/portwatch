// Package notify provides notifier implementations for dispatching portwatch
// alerts to external services.
//
// # Signal Sciences / Fastly Next-Gen WAF
//
// NewSignalSciences creates a notifier that posts custom events to the Signal
// Sciences API (https://dashboard.signalsciences.net/api). Authentication uses
// the x-api-user / x-api-token header pair issued from the Signal Sciences
// dashboard.
//
// Configuration fields:
//
//	api_url          – base URL of the Signal Sciences API (default: https://dashboard.signalsciences.net/api)
//	corp_name        – Signal Sciences corporation slug
//	site_name        – Signal Sciences site slug
//	access_key_id    – API access key ID
//	secret_access_key – API secret access key
package notify
