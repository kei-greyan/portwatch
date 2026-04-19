package metrics_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/portwatch/internal/metrics"
)

func TestReporter_Write_ContainsExpectedKeys(t *testing.T) {
	m := metrics.New()
	m.RecordScan(7)
	m.RecordAlert()

	r := metrics.NewReporter(m)
	var buf bytes.Buffer
	if err := r.Write(&buf); err != nil {
		t.Fatalf("Write returned error: %v", err)
	}

	out := buf.String()
	for _, key := range []string{"uptime", "scans_total", "alerts_total", "ports_open", "last_scan_at"} {
		if !strings.Contains(out, key) {
			t.Errorf("output missing key %q\n%s", key, out)
		}
	}
}

func TestReporter_Write_ReflectsCounters(t *testing.T) {
	m := metrics.New()
	m.RecordScan(42)
	m.RecordAlert()
	m.RecordAlert()

	r := metrics.NewReporter(m)
	var buf bytes.Buffer
	_ = r.Write(&buf)
	out := buf.String()

	if !strings.Contains(out, "42") {
		t.Errorf("expected ports_open=42 in output:\n%s", out)
	}
	if !strings.Contains(out, "2") {
		t.Errorf("expected alerts_total=2 in output:\n%s", out)
	}
}

func TestReporter_Write_LastScanNeverWhenNoScans(t *testing.T) {
	m := metrics.New()
	r := metrics.NewReporter(m)
	var buf bytes.Buffer
	_ = r.Write(&buf)

	if !strings.Contains(buf.String(), "never") {
		t.Errorf("expected 'never' for last_scan_at before any scan")
	}
}
