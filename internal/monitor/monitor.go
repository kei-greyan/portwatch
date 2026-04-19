package monitor

import (
	"context"
	"log"
	"time"

	"github.com/yourorg/portwatch/internal/alert"
	"github.com/yourorg/portwatch/internal/scanner"
	"github.com/yourorg/portwatch/internal/state"
)

// Config controls Monitor behaviour.
type Config struct {
	Interval time.Duration
	StatePath string
}

// Monitor ties together scanning, diffing, alerting and state persistence.
type Monitor struct {
	cfg     Config
	scanner *scanner.Scanner
	store   *state.Store
	alerter *alert.Alerter
}

// New constructs a Monitor from the provided config and dependencies.
func New(cfg Config, sc *scanner.Scanner, st *state.Store, al *alert.Alerter) *Monitor {
	return &Monitor{cfg: cfg, scanner: sc, store: st, alerter: al}
}

// Run starts the monitoring loop and blocks until ctx is cancelled.
func (m *Monitor) Run(ctx context.Context) error {
	ticker := time.NewTicker(m.cfg.Interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			if err := m.tick(); err != nil {
				log.Printf("monitor tick error: %v", err)
			}
		}
	}
}

func (m *Monitor) tick() error {
	current, err := m.scanner.Scan()
	if err != nil {
		return err
	}

	prev, err := m.store.Load()
	if err != nil {
		return err
	}

	diff := alert.Diff(prev.Ports, current)
	for _, p := range diff.Opened {
		m.alerter.Send(alert.PortOpened(p))
	}
	for _, p := range diff.Closed {
		m.alerter.Send(alert.PortClosed(p))
	}

	return m.store.Save(state.Snapshot{Ports: current})
}
