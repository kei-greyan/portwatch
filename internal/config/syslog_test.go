package config

import "testing"

func validSyslogConfig() *SyslogConfig {
	c := defaultSyslogConfig()
	c.Enabled = true
	return c
}

func TestValidateSyslog_NilAlwaysPasses(t *testing.T) {
	if err := ValidateSyslog(nil); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestValidateSyslog_DisabledAlwaysPasses(t *testing.T) {
	cfg := validSyslogConfig()
	cfg.Enabled = false
	cfg.Tag = ""
	if err := ValidateSyslog(cfg); err != nil {
		t.Fatalf("expected nil for disabled config, got %v", err)
	}
}

func TestValidateSyslog_AcceptsLocalSocket(t *testing.T) {
	cfg := validSyslogConfig()
	cfg.Network = ""
	cfg.Addr = ""
	if err := ValidateSyslog(cfg); err != nil {
		t.Fatalf("local socket config should be valid, got %v", err)
	}
}

func TestValidateSyslog_AcceptsRemoteTCP(t *testing.T) {
	cfg := validSyslogConfig()
	cfg.Network = "tcp"
	cfg.Addr = "logs.example.com:514"
	if err := ValidateSyslog(cfg); err != nil {
		t.Fatalf("remote TCP config should be valid, got %v", err)
	}
}

func TestValidateSyslog_RejectsUnknownNetwork(t *testing.T) {
	cfg := validSyslogConfig()
	cfg.Network = "unix"
	cfg.Addr = "/dev/log"
	if err := ValidateSyslog(cfg); err == nil {
		t.Fatal("expected error for unknown network")
	}
}

func TestValidateSyslog_RejectsMissingAddr(t *testing.T) {
	cfg := validSyslogConfig()
	cfg.Network = "udp"
	cfg.Addr = ""
	if err := ValidateSyslog(cfg); err == nil {
		t.Fatal("expected error when network set but addr empty")
	}
}

func TestValidateSyslog_RejectsEmptyTag(t *testing.T) {
	cfg := validSyslogConfig()
	cfg.Tag = ""
	if err := ValidateSyslog(cfg); err == nil {
		t.Fatal("expected error for empty tag")
	}
}
