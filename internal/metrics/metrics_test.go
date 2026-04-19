package metrics_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/metrics"
)

func TestNew_SetsUptimeStart(t *testing.T) {
	before := time.Now()
	m := metrics.New()
	after := time.Now()

	s := m.Snapshot()
	if s.UptimeStart.Before(before) || s.UptimeStart.After(after) {
		t.Errorf("UptimeStart out of range: %v", s.UptimeStart)
	}
}

func TestRecordScan_IncrementsCounter(t *testing.T) {
	m := metrics.New()
	m.RecordScan(5)
	m.RecordScan(3)

	s := m.Snapshot()
	if s.ScansTotal != 2 {
		t.Errorf("expected ScansTotal=2, got %d", s.ScansTotal)
	}
	if s.PortsOpen != 3 {
		t.Errorf("expected PortsOpen=3, got %d", s.PortsOpen)
	}
	if s.LastScanAt.IsZero() {
		t.Error("expected LastScanAt to be set")
	}
}

func TestRecordAlert_IncrementsCounter(t *testing.T) {
	m := metrics.New()
	m.RecordAlert()
	m.RecordAlert()
	m.RecordAlert()

	s := m.Snapshot()
	if s.AlertsTotal != 3 {
		t.Errorf("expected AlertsTotal=3, got %d", s.AlertsTotal)
	}
}

func TestSnapshot_IsConsistent(t *testing.T) {
	m := metrics.New()
	m.RecordScan(10)
	m.RecordAlert()

	s := m.Snapshot()
	if s.ScansTotal != 1 || s.AlertsTotal != 1 || s.PortsOpen != 10 {
		t.Errorf("unexpected snapshot: %+v", s)
	}
}
