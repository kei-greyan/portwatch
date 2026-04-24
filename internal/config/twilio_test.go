package config_test

import (
	"testing"

	"github.com/yourorg/portwatch/internal/config"
)

func validTwilioConfig() *config.TwilioConfig {
	return &config.TwilioConfig{
		Enabled:    true,
		AccountSID: "ACxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
		AuthToken:  "secret",
		From:       "+15550001111",
		To:         "+15559998888",
	}
}

func TestValidateTwilio_NilAlwaysPasses(t *testing.T) {
	if err := config.ValidateTwilio(nil); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestValidateTwilio_DisabledAlwaysPasses(t *testing.T) {
	c := validTwilioConfig()
	c.Enabled = false
	if err := config.ValidateTwilio(c); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestValidateTwilio_AcceptsValidConfig(t *testing.T) {
	if err := config.ValidateTwilio(validTwilioConfig()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidateTwilio_RejectsEmptyAccountSID(t *testing.T) {
	c := validTwilioConfig()
	c.AccountSID = ""
	if err := config.ValidateTwilio(c); err == nil {
		t.Fatal("expected error for empty account_sid")
	}
}

func TestValidateTwilio_RejectsEmptyAuthToken(t *testing.T) {
	c := validTwilioConfig()
	c.AuthToken = ""
	if err := config.ValidateTwilio(c); err == nil {
		t.Fatal("expected error for empty auth_token")
	}
}

func TestValidateTwilio_RejectsEmptyFrom(t *testing.T) {
	c := validTwilioConfig()
	c.From = ""
	if err := config.ValidateTwilio(c); err == nil {
		t.Fatal("expected error for empty from")
	}
}

func TestValidateTwilio_RejectsEmptyTo(t *testing.T) {
	c := validTwilioConfig()
	c.To = ""
	if err := config.ValidateTwilio(c); err == nil {
		t.Fatal("expected error for empty to")
	}
}
