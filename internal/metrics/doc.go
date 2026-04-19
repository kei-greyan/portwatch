// Package metrics tracks runtime statistics for portwatch.
//
// It records scan counts, alert counts, and uptime. A Reporter
// formats snapshots as human-readable key/value pairs, and an
// HTTP handler exposes them as JSON for external scraping.
//
// Usage:
//
//	m := metrics.New()
//	m.RecordScan()
//	m.RecordAlert()
//	snap := m.Snapshot()
package metrics
