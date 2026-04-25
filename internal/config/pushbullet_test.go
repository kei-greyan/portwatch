package config_test

import (
	"testing"

	"github.com/user/portwatch/internal/config"
)

func validPushbulletConfig() *config.PushbulletConfig {
	return &config.PushbulletConfig{
		Enabled: true,
		APIKey:  "test-api-key",
		APIURL:  "https://api.pushbullet.com/v2/pushes",
	}
}

func TestValidatePushbullet_NilAlwaysPasses(t *testing.T) {
	if err := config.ValidatePushbullet(nil); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestValidatePushbullet_DisabledAlwaysPasses(t *testing.T) {
	cfg := validPushbulletConfig()
	cfg.Enabled = false
	cfg.APIKey = ""
	if err := config.ValidatePushbullet(cfg); err != nil {
		t.Fatalf("expected nil for disabled config, got %v", err)
	}
}

func TestValidatePushbullet_AcceptsValidConfig(t *testing.T) {
	if err := config.ValidatePushbullet(validPushbulletConfig()); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestValidatePushbullet_RejectsEmptyAPIKey(t *testing.T) {
	cfg := validPushbulletConfig()
	cfg.APIKey = ""
	if err := config.ValidatePushbullet(cfg); err == nil {
		t.Fatal("expected error for empty api_key, got nil")
	}
}

func TestValidatePushbullet_RejectsEmptyAPIURL(t *testing.T) {
	cfg := validPushbulletConfig()
	cfg.APIURL = ""
	if err := config.ValidatePushbullet(cfg); err == nil {
		t.Fatal("expected error for empty api_url, got nil")
	}
}
