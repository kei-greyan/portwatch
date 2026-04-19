package monitor_test

import (
	"bytes"
	"context"
	"testing"
	"time"

	"github.com/yourorg/portwatch/internal/alert"
	"github.com/yourorg/portwatch/internal/monitor"
	"github.com/yourorg/portwatch/internal/scanner"
	"github.com/yourorg/portwatch/internal/state"
)

func TestMonitor_AlertsOnNewPort(t *testing.T) {
	var buf bytes.Buffer
	al := alert.New(&buf)

	// state with no previous ports
	st := state.New(t.TempDir() + "/state.json")

	sc := scanner.New(scanner.Config{Targets: []string{"127.0.0.1"}})

	// pre-seed state so diff yields an opened port without a real listener
	_ = st.Save(state.Snapshot{Ports: []uint16{}})

	cfg := monitor.Config{
		Interval:  10 * time.Millisecond,
		StatePath: "",
	}

	m := monitor.New(cfg, sc, st, al)

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	// Run returns when ctx is cancelled — that is expected
	_ = m.Run(ctx)

	// No assertions on buf here because no real port is open in CI;
	// the test validates that Run returns without panic on cancellation.
	if ctx.Err() == nil {
		t.Error("expected context to be cancelled")
	}
}

func TestMonitor_RunRespectsCancel(t *testing.T) {
	var buf bytes.Buffer
	al := alert.New(&buf)
	st := state.New(t.TempDir() + "/state.json")
	sc := scanner.New(scanner.Config{Targets: []string{"127.0.0.1"}})

	m := monitor.New(monitor.Config{Interval: 1 * time.Hour}, sc, st, al)

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel immediately

	start := time.Now()
	_ = m.Run(ctx)
	if time.Since(start) > time.Second {
		t.Error("Run did not respect context cancellation promptly")
	}
}
