package config

import "testing"

func validProwlConfig() *ProwlConfig {
	return &ProwlConfig{
		Enabled: true,
		APIKey:  "abc123",
		AppName: "portwatch",
	}
}

func TestValidateProwl_NilAlwaysPasses(t *testing.T) {
	if err := ValidateProwl(nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidateProwl_DisabledAlwaysPasses(t *testing.T) {
	c := validProwlConfig()
	c.Enabled = false
	c.APIKey = ""
	if err := ValidateProwl(c); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidateProwl_AcceptsValidConfig(t *testing.T) {
	if err := ValidateProwl(validProwlConfig()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidateProwl_RejectsEmptyAPIKey(t *testing.T) {
	c := validProwlConfig()
	c.APIKey = ""
	if err := ValidateProwl(c); err == nil {
		t.Fatal("expected error for empty api_key")
	}
}

func TestValidateProwl_RejectsEmptyAppName(t *testing.T) {
	c := validProwlConfig()
	c.AppName = ""
	if err := ValidateProwl(c); err == nil {
		t.Fatal("expected error for empty app_name")
	}
}
