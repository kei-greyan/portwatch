// Package filter implements port ignore-list filtering for portwatch.
//
// Rules are expressed as individual ports or inclusive ranges:
//
//	"22"        – ignore port 22
//	"1000-2000" – ignore ports 1000 through 2000 inclusive
//
// A Filter is constructed once from configuration and consulted by the
// monitor before raising alerts, so that well-known or expected ports
// can be silently skipped.
package filter
