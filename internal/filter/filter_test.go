package filter

import (
	"testing"
)

func TestNew_AcceptsValidRules(t *testing.T) {
	_, err := New([]string{"22", "80", "1000-2000"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNew_RejectsInvalidRule(t *testing.T) {
	_, err := New([]string{"not-a-port"})
	if err == nil {
		t.Fatal("expected error for invalid rule")
	}
}

func TestNew_RejectsInvertedRange(t *testing.T) {
	_, err := New([]string{"2000-1000"})
	if err == nil {
		t.Fatal("expected error for inverted range")
	}
}

func TestIgnored_ExactMatch(t *testing.T) {
	f, _ := New([]string{"22"})
	if !f.Ignored(22) {
		t.Error("expected port 22 to be ignored")
	}
	if f.Ignored(23) {
		t.Error("expected port 23 not to be ignored")
	}
}

func TestIgnored_RangeMatch(t *testing.T) {
	f, _ := New([]string{"1000-2000"})
	for _, p := range []uint16{1000, 1500, 2000} {
		if !f.Ignored(p) {
			t.Errorf("expected port %d to be ignored", p)
		}
	}
	for _, p := range []uint16{999, 2001} {
		if f.Ignored(p) {
			t.Errorf("expected port %d not to be ignored", p)
		}
	}
}

func TestIgnored_EmptyFilter(t *testing.T) {
	f, _ := New(nil)
	if f.Ignored(80) {
		t.Error("empty filter should not ignore any port")
	}
}

func TestIgnored_MultipleRules(t *testing.T) {
	f, _ := New([]string{"22", "443", "8000-8080"})
	for _, p := range []uint16{22, 443, 8000, 8040, 8080} {
		if !f.Ignored(p) {
			t.Errorf("expected port %d to be ignored", p)
		}
	}
	if f.Ignored(80) {
		t.Error("expected port 80 not to be ignored")
	}
}
