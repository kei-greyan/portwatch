package config

import (
	"testing"
)

func TestMigrate_NilReturnsError(t *testing.T) {
	if err := Migrate(nil); err == nil {
		t.Fatal("expected error for nil config")
	}
}

func TestMigrate_AlreadyCurrent(t *testing.T) {
	cfg := defaults()
	cfg.Version = currentVersion
	if err := Migrate(cfg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Version != currentVersion {
		t.Fatalf("version changed unexpectedly: got %d", cfg.Version)
	}
}

func TestMigrate_V0FillsStatePath(t *testing.T) {
	cfg := defaults()
	cfg.Version = 0
	cfg.StatePath = ""
	if err := Migrate(cfg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.StatePath == "" {
		t.Fatal("expected StatePath to be filled in after migration")
	}
	if cfg.Version != 1 {
		t.Fatalf("expected version 1, got %d", cfg.Version)
	}
}

func TestMigrate_V0PreservesExistingStatePath(t *testing.T) {
	cfg := defaults()
	cfg.Version = 0
	cfg.StatePath = "/custom/path.json"
	if err := Migrate(cfg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.StatePath != "/custom/path.json" {
		t.Fatalf("StatePath overwritten: got %s", cfg.StatePath)
	}
}

func TestMigrate_FutureVersionReturnsError(t *testing.T) {
	cfg := defaults()
	cfg.Version = currentVersion + 1
	if err := Migrate(cfg); err == nil {
		t.Fatal("expected error for future version")
	}
}

func TestMigrate_V0AdvancesToCurrentVersion(t *testing.T) {
	cfg := defaults()
	cfg.Version = 0
	if err := Migrate(cfg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Version != currentVersion {
		t.Fatalf("expected version %d after full migration, got %d", currentVersion, cfg.Version)
	}
}
