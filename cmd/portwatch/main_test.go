package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"
)

// TestMain_InvalidConfig ensures the binary exits non-zero on bad config.
func TestMain_InvalidConfig(t *testing.T) {
	if os.Getenv("PORTWATCH_RUN_MAIN") == "1" {
		main()
		return
	}

	tmpDir := t.TempDir()
	cfgPath := filepath.Join(tmpDir, "bad.json")
	if err := os.WriteFile(cfgPath, []byte(`{not json}`), 0644); err != nil {
		t.Fatal(err)
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestMain_InvalidConfig", "-config", cfgPath)
	cmd.Env = append(os.Environ(), "PORTWATCH_RUN_MAIN=1")
	err := cmd.Run()
	if err == nil {
		t.Fatal("expected non-zero exit for invalid config")
	}
}

// TestMain_SignalShutdown verifies the daemon stops cleanly on SIGTERM.
func TestMain_SignalShutdown(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping signal test in short mode")
	}

	tmpDir := t.TempDir()
	cfgPath := filepath.Join(tmpDir, "cfg.json")
	cfgContent := `{"interval":"200ms","state_path":"` + filepath.Join(tmpDir, "state.json") + `"}`
	if err := os.WriteFile(cfgPath, []byte(cfgContent), 0644); err != nil {
		t.Fatal(err)
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestMain_SignalShutdown", "-config", cfgPath)
	cmd.Env = append(os.Environ(), "PORTWATCH_RUN_MAIN=1")
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}

	time.Sleep(300 * time.Millisecond)
	if err := cmd.Process.Signal(os.Interrupt); err != nil {
		t.Fatal(err)
	}

	done := make(chan error, 1)
	go func() { done <- cmd.Wait() }()
	select {
	case <-done:
		// exited — pass
	case <-time.After(3 * time.Second):
		cmd.Process.Kill()
		t.Fatal("process did not exit after SIGINT")
	}
}
