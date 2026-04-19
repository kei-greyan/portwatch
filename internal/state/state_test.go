package state_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/yourorg/portwatch/internal/state"
)

func tmpPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "state.json")
}

func TestLoad_ReturnsEmptyWhenMissing(t *testing.T) {
	s := state.New(tmpPath(t))
	snap, err := s.Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(snap.Ports) != 0 {
		t.Errorf("expected empty ports, got %v", snap.Ports)
	}
}

func TestSave_PersistsAndLoads(t *testing.T) {
	s := state.New(tmpPath(t))
	want := state.Snapshot{
		Ports:      []uint16{80, 443, 8080},
		RecordedAt: time.Now().Truncate(time.Second),
	}
	if err := s.Save(want); err != nil {
		t.Fatalf("save error: %v", err)
	}
	got, err := s.Load()
	if err != nil {
		t.Fatalf("load error: %v", err)
	}
	if len(got.Ports) != len(want.Ports) {
		t.Errorf("ports length mismatch: want %d got %d", len(want.Ports), len(got.Ports))
	}
	for i, p := range want.Ports {
		if got.Ports[i] != p {
			t.Errorf("port[%d]: want %d got %d", i, p, got.Ports[i])
		}
	}
}

func TestSave_SetsTimestampIfZero(t *testing.T) {
	s := state.New(tmpPath(t))
	if err := s.Save(state.Snapshot{Ports: []uint16{22}}); err != nil {
		t.Fatalf("save error: %v", err)
	}
	got, _ := s.Load()
	if got.RecordedAt.IsZero() {
		t.Error("expected RecordedAt to be set automatically")
	}
}

func TestLoad_ErrorOnCorruptFile(t *testing.T) {
	p := tmpPath(t)
	if err := os.WriteFile(p, []byte("not json{"), 0o600); err != nil {
		t.Fatal(err)
	}
	s := state.New(p)
	_, err := s.Load()
	if err == nil {
		t.Error("expected error for corrupt file")
	}
}
