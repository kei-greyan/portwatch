package ratelimit

import (
	"testing"
	"time"
)

func TestAllow_FirstCallPermitted(t *testing.T) {
	l := New(time.Minute)
	if !l.Allow("80/tcp") {
		t.Fatal("expected first call to be allowed")
	}
}

func TestAllow_SecondCallWithinWindowBlocked(t *testing.T) {
	l := New(time.Minute)
	l.Allow("80/tcp")
	if l.Allow("80/tcp") {
		t.Fatal("expected second call within window to be blocked")
	}
}

func TestAllow_CallAfterWindowPermitted(t *testing.T) {
	now := time.Now()
	l := New(time.Minute)
	l.nowFunc = func() time.Time { return now }
	l.Allow("80/tcp")

	// advance past the window
	l.nowFunc = func() time.Time { return now.Add(2 * time.Minute) }
	if !l.Allow("80/tcp") {
		t.Fatal("expected call after window to be allowed")
	}
}

func TestAllow_IndependentKeys(t *testing.T) {
	l := New(time.Minute)
	l.Allow("80/tcp")
	if !l.Allow("443/tcp") {
		t.Fatal("expected different key to be allowed independently")
	}
}

func TestReset_AllowsImmediateRetry(t *testing.T) {
	l := New(time.Minute)
	l.Allow("80/tcp")
	l.Reset("80/tcp")
	if !l.Allow("80/tcp") {
		t.Fatal("expected allow after reset")
	}
}

func TestReset_UnknownKeyIsNoop(t *testing.T) {
	l := New(time.Minute)
	// should not panic
	l.Reset("9999/tcp")
}
