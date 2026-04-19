package metrics

import (
	"fmt"
	"io"
	"text/tabwriter"
	"time"
)

// Reporter writes human-readable metric summaries.
type Reporter struct {
	m *Metrics
}

// NewReporter returns a Reporter backed by the given Metrics.
func NewReporter(m *Metrics) *Reporter {
	return &Reporter{m: m}
}

// Write formats the current snapshot as a table to w.
func (r *Reporter) Write(w io.Writer) error {
	s := r.m.Snapshot()
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)

	uptime := time.Since(s.UptimeStart).Truncate(time.Second)
	lastScan := "never"
	if !s.LastScanAt.IsZero() {
		lastScan = s.LastScanAt.Format(time.RFC3339)
	}

	rows := []struct{ k, v string }{
		{"uptime", uptime.String()},
		{"scans_total", fmt.Sprintf("%d", s.ScansTotal)},
		{"alerts_total", fmt.Sprintf("%d", s.AlertsTotal)},
		{"ports_open", fmt.Sprintf("%d", s.PortsOpen)},
		{"last_scan_at", lastScan},
	}

	for _, row := range rows {
		if _, err := fmt.Fprintf(tw, "%s\t%s\n", row.k, row.v); err != nil {
			return err
		}
	}
	return tw.Flush()
}
