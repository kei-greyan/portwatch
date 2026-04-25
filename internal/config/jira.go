package config

import "fmt"

// JiraConfig holds settings for the Jira notifier.
type JiraConfig struct {
	Enabled    bool   `toml:"enabled"`
	BaseURL    string `toml:"base_url"`
	ProjectKey string `toml:"project_key"`
	Username   string `toml:"username"`
	Token      string `toml:"token"`
}

func defaultJiraConfig() *JiraConfig {
	return &JiraConfig{
		Enabled: false,
	}
}

// ValidateJira returns an error if cfg contains invalid Jira settings.
func ValidateJira(cfg *JiraConfig) error {
	if cfg == nil {
		return nil
	}
	if !cfg.Enabled {
		return nil
	}
	if cfg.BaseURL == "" {
		return fmt.Errorf("jira: base_url must not be empty")
	}
	if cfg.ProjectKey == "" {
		return fmt.Errorf("jira: project_key must not be empty")
	}
	if cfg.Username == "" {
		return fmt.Errorf("jira: username must not be empty")
	}
	if cfg.Token == "" {
		return fmt.Errorf("jira: token must not be empty")
	}
	return nil
}
