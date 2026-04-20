package config

// EmailConfig holds optional SMTP notification settings.
type EmailConfig struct {
	Enabled  bool     `json:"enabled"`
	Host     string   `json:"host"`
	Port     int      `json:"port"`
	Username string   `json:"username,omitempty"`
	Password string   `json:"password,omitempty"`
	From     string   `json:"from"`
	To       []string `json:"to"`
}

func defaultEmailConfig() EmailConfig {
	return EmailConfig{
		Enabled: false,
		Host:    "localhost",
		Port:    25,
	}
}

// ValidateEmail returns an error if the EmailConfig is enabled but incomplete.
func ValidateEmail(e *EmailConfig) error {
	if e == nil || !e.Enabled {
		return nil
	}
	if e.Host == "" {
		return &ValidationError{Field: "email.host", Reason: "must not be empty when email is enabled"}
	}
	if e.Port <= 0 || e.Port > 65535 {
		return &ValidationError{Field: "email.port", Reason: "must be between 1 and 65535"}
	}
	if e.From == "" {
		return &ValidationError{Field: "email.from", Reason: "must not be empty when email is enabled"}
	}
	if len(e.To) == 0 {
		return &ValidationError{Field: "email.to", Reason: "must have at least one recipient"}
	}
	return nil
}

// ValidationError describes a configuration validation failure.
type ValidationError struct {
	Field  string
	Reason string
}

func (v *ValidationError) Error() string {
	return v.Field + ": " + v.Reason
}
