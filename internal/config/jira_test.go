package config_test

import (
	"testing"

	"github.com/user/portwatch/internal/config"
)

func validJiraConfig() *config.JiraConfig {
	return &config.JiraConfig{
		Enabled:    true,
		BaseURL:    "https://jira.example.com",
		ProjectKey: "OPS",
		Username:   "admin",
		Token:      "secret-token",
	}
}

func TestValidateJira_NilAlwaysPasses(t *testing.T) {
	if err := config.ValidateJira(nil); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestValidateJira_DisabledAlwaysPasses(t *testing.T) {
	cfg := validJiraConfig()
	cfg.Enabled = false
	cfg.BaseURL = ""
	if err := config.ValidateJira(cfg); err != nil {
		t.Fatalf("expected nil for disabled config, got %v", err)
	}
}

func TestValidateJira_AcceptsValidConfig(t *testing.T) {
	if err := config.ValidateJira(validJiraConfig()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidateJira_RejectsEmptyBaseURL(t *testing.T) {
	cfg := validJiraConfig()
	cfg.BaseURL = ""
	if err := config.ValidateJira(cfg); err == nil {
		t.Fatal("expected error for empty base_url")
	}
}

func TestValidateJira_RejectsEmptyProjectKey(t *testing.T) {
	cfg := validJiraConfig()
	cfg.ProjectKey = ""
	if err := config.ValidateJira(cfg); err == nil {
		t.Fatal("expected error for empty project_key")
	}
}

func TestValidateJira_RejectsEmptyToken(t *testing.T) {
	cfg := validJiraConfig()
	cfg.Token = ""
	if err := config.ValidateJira(cfg); err == nil {
		t.Fatal("expected error for empty token")
	}
}
