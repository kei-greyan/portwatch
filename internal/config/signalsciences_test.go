package config_test

import (
	"testing"

	"github.com/user/portwatch/internal/config"
)

func validSignalSciencesConfig() *config.SignalSciencesConfig {
	return &config.SignalSciencesConfig{
		Enabled:         true,
		APIURL:          "https://dashboard.signalsciences.net/api",
		CorpName:        "mycorp",
		SiteName:        "mysite",
		AccessKeyID:     "keyid",
		SecretAccessKey: "secret",
	}
}

func TestValidateSignalSciences_NilAlwaysPasses(t *testing.T) {
	if err := config.ValidateSignalSciences(nil); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestValidateSignalSciences_DisabledAlwaysPasses(t *testing.T) {
	cfg := validSignalSciencesConfig()
	cfg.Enabled = false
	if err := config.ValidateSignalSciences(cfg); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestValidateSignalSciences_AcceptsValidConfig(t *testing.T) {
	if err := config.ValidateSignalSciences(validSignalSciencesConfig()); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestValidateSignalSciences_RejectsEmptyAPIURL(t *testing.T) {
	cfg := validSignalSciencesConfig()
	cfg.APIURL = ""
	if err := config.ValidateSignalSciences(cfg); err == nil {
		t.Error("expected error for empty api_url")
	}
}

func TestValidateSignalSciences_RejectsEmptyCorpName(t *testing.T) {
	cfg := validSignalSciencesConfig()
	cfg.CorpName = ""
	if err := config.ValidateSignalSciences(cfg); err == nil {
		t.Error("expected error for empty corp_name")
	}
}

func TestValidateSignalSciences_RejectsEmptySiteName(t *testing.T) {
	cfg := validSignalSciencesConfig()
	cfg.SiteName = ""
	if err := config.ValidateSignalSciences(cfg); err == nil {
		t.Error("expected error for empty site_name")
	}
}

func TestValidateSignalSciences_RejectsEmptyAccessKeyID(t *testing.T) {
	cfg := validSignalSciencesConfig()
	cfg.AccessKeyID = ""
	if err := config.ValidateSignalSciences(cfg); err == nil {
		t.Error("expected error for empty access_key_id")
	}
}

func TestValidateSignalSciences_RejectsEmptySecretAccessKey(t *testing.T) {
	cfg := validSignalSciencesConfig()
	cfg.SecretAccessKey = ""
	if err := config.ValidateSignalSciences(cfg); err == nil {
		t.Error("expected error for empty secret_access_key")
	}
}
