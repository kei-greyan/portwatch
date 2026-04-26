package config_test

import (
	"testing"

	"github.com/example/portwatch/internal/config"
)

func validGooglePubSubConfig() *config.GooglePubSubConfig {
	return &config.GooglePubSubConfig{
		Enabled:        true,
		ProjectID:      "my-project",
		TopicID:        "portwatch-alerts",
		TimeoutSeconds: 10,
	}
}

func TestValidateGooglePubSub_NilAlwaysPasses(t *testing.T) {
	if err := config.ValidateGooglePubSub(nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidateGooglePubSub_DisabledAlwaysPasses(t *testing.T) {
	cfg := validGooglePubSubConfig()
	cfg.Enabled = false
	cfg.ProjectID = ""
	if err := config.ValidateGooglePubSub(cfg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidateGooglePubSub_AcceptsValidConfig(t *testing.T) {
	if err := config.ValidateGooglePubSub(validGooglePubSubConfig()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidateGooglePubSub_RejectsEmptyProjectID(t *testing.T) {
	cfg := validGooglePubSubConfig()
	cfg.ProjectID = ""
	if err := config.ValidateGooglePubSub(cfg); err == nil {
		t.Fatal("expected error for empty project_id")
	}
}

func TestValidateGooglePubSub_RejectsEmptyTopicID(t *testing.T) {
	cfg := validGooglePubSubConfig()
	cfg.TopicID = ""
	if err := config.ValidateGooglePubSub(cfg); err == nil {
		t.Fatal("expected error for empty topic_id")
	}
}

func TestValidateGooglePubSub_RejectsZeroTimeout(t *testing.T) {
	cfg := validGooglePubSubConfig()
	cfg.TimeoutSeconds = 0
	if err := config.ValidateGooglePubSub(cfg); err == nil {
		t.Fatal("expected error for zero timeout_seconds")
	}
}
