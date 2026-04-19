package config

import "fmt"

// currentVersion is the schema version this binary understands.
const currentVersion = 1

// Migrate upgrades cfg in-place from its declared Version to currentVersion.
// It returns an error if the version is unknown or a downgrade is attempted.
func Migrate(cfg *Config) error {
	if cfg == nil {
		return fmt.Errorf("migrate: nil config")
	}
	if cfg.Version > currentVersion {
		return fmt.Errorf("migrate: config version %d is newer than supported version %d",
			cfg.Version, currentVersion)
	}
	for cfg.Version < currentVersion {
		if err := migrateStep(cfg); err != nil {
			return err
		}
	}
	return nil
}

// migrateStep advances cfg by exactly one version.
func migrateStep(cfg *Config) error {
	switch cfg.Version {
	case 0:
		// v0 → v1: StatePath default was empty; fill it in.
		if cfg.StatePath == "" {
			cfg.StatePath = defaults().StatePath
		}
		cfg.Version = 1
		return nil
	default:
		return fmt.Errorf("migrate: no migration path from version %d", cfg.Version)
	}
}
