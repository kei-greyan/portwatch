// Package ratelimit provides a simple per-key rate limiter used to
// suppress repeated alerts for the same port within a cooldown window.
package ratelimit

import (
	"sync"
	"time"
)

// Limiter suppresses duplicate events for the same key within a window.
type Limiter struct {
	mu      sync.Mutex
	window  time.Duration
	last    map[string]time.Time
	nowFunc func() time.Time
}

// New creates a Limiter with the given cooldown window.
func New(window time.Duration) *Limiter {
	return &Limiter{
		window:  window,
		last:    make(map[string]time.Time),
		nowFunc: time.Now,
	}
}

// Allow returns true if the event for key should be allowed through,
// i.e. it has not been seen within the cooldown window. Calling Allow
// with a new or expired key updates the recorded time.
func (l *Limiter) Allow(key string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := l.nowFunc()
	if t, ok := l.last[key]; ok && now.Sub(t) < l.window {
		return false
	}
	l.last[key] = now
	return true
}

// Reset clears the recorded time for key, allowing the next event
// through immediately regardless of the window.
func (l *Limiter) Reset(key string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	delete(l.last, key)
}
