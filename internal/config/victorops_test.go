package config_test

import (
	"testing"

	"github.com/user/portwatch/internal/config"
)

func validVictorOpsConfig() *config.VictorOpsConfig {
	return &config.VictorOpsConfig{
		Enabled:    true,
		APIURL:     "https://alert.victorops.com/integrations/generic/20131114/alert/apikey/routingkey",
		RoutingKey: "default",
	}
}

func TestValidateVictorOps_NilAlwaysPasses(t *testing.T) {
	if err := config.ValidateVictorOps(nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidateVictorOps_DisabledAlwaysPasses(t *testing.T) {
	cfg := validVictorOpsConfig()
	cfg.Enabled = false
	cfg.APIURL = ""
	if err := config.ValidateVictorOps(cfg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidateVictorOps_AcceptsValidConfig(t *testing.T) {
	if err := config.ValidateVictorOps(validVictorOpsConfig()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidateVictorOps_RejectsEmptyAPIURL(t *testing.T) {
	cfg := validVictorOpsConfig()
	cfg.APIURL = ""
	if err := config.ValidateVictorOps(cfg); err == nil {
		t.Error("expected error for empty api_url")
	}
}

func TestValidateVictorOps_RejectsEmptyRoutingKey(t *testing.T) {
	cfg := validVictorOpsConfig()
	cfg.RoutingKey = ""
	if err := config.ValidateVictorOps(cfg); err == nil {
		t.Error("expected error for empty routing_key")
	}
}
