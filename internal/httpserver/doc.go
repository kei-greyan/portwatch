// Package httpserver provides a small HTTP server used by portwatch to
// expose operational endpoints.
//
// Endpoints:
//
//	 GET /metrics  – JSON snapshot of runtime metrics (via metrics.Metrics.Handler)
//	 GET /health   – simple liveness probe returning 200 OK
//
// The server supports graceful shutdown and enforces conservative read/write
// timeouts to prevent resource leaks.
package httpserver
