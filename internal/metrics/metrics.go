// Package metrics tracks runtime counters for portwatch.
package metrics

import (
	"sync"
	"time"
)

// Snapshot holds a point-in-time copy of all counters.
type Snapshot struct {
	ScansTotal   uint64
	AlertsTotal  uint64
	PortsOpen    int
	LastScanAt   time.Time
	UptimeStart  time.Time
}

// Metrics is a thread-safe counter store.
type Metrics struct {
	mu          sync.RWMutex
	scansTotal  uint64
	alertsTotal uint64
	portsOpen   int
	lastScanAt  time.Time
	uptimeStart time.Time
}

// New returns a Metrics instance with the uptime clock started.
func New() *Metrics {
	return &Metrics{uptimeStart: time.Now()}
}

// RecordScan increments the scan counter and records the open port count.
func (m *Metrics) RecordScan(openPorts int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.scansTotal++
	m.portsOpen = openPorts
	m.lastScanAt = time.Now()
}

// RecordAlert increments the alert counter.
func (m *Metrics) RecordAlert() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.alertsTotal++
}

// Snapshot returns a consistent copy of current counters.
func (m *Metrics) Snapshot() Snapshot {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return Snapshot{
		ScansTotal:  m.scansTotal,
		AlertsTotal: m.alertsTotal,
		PortsOpen:   m.portsOpen,
		LastScanAt:  m.lastScanAt,
		UptimeStart: m.uptimeStart,
	}
}
