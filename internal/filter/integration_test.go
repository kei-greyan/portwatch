package filter_test

import (
	"testing"

	"github.com/yourorg/portwatch/internal/filter"
)

// TestFilter_ConfigIgnoredPorts simulates loading ignored_ports from config
// and verifying the filter behaves correctly end-to-end.
func TestFilter_ConfigIgnoredPorts(t *testing.T) {
	// Typical ignored_ports from a config file
	specs := []string{"22", "53", "443", "32768-60999"}

	f, err := filter.New(specs)
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	expectIgnored := []uint16{22, 53, 443, 32768, 45000, 60999}
	for _, p := range expectIgnored {
		if !f.Ignored(p) {
			t.Errorf("port %d should be ignored", p)
		}
	}

	expectAllowed := []uint16{80, 8080, 3000, 61000}
	for _, p := range expectAllowed {
		if f.Ignored(p) {
			t.Errorf("port %d should not be ignored", p)
		}
	}
}
