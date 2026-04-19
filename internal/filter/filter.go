// Package filter provides port filtering logic for portwatch.
// It allows users to specify ports or ranges to ignore during monitoring.
package filter

import (
	"fmt"
	"strconv"
	"strings"
)

// Rule represents a single port filter rule (exact port or range).
type Rule struct {
	Low  uint16
	High uint16
}

// Filter holds compiled filter rules.
type Filter struct {
	rules []Rule
}

// New creates a Filter from a slice of rule strings.
// Each string is either a single port ("22") or a range ("1000-2000").
func New(specs []string) (*Filter, error) {
	f := &Filter{}
	for _, s := range specs {
		r, err := parseRule(s)
		if err != nil {
			return nil, fmt.Errorf("filter: invalid rule %q: %w", s, err)
		}
		f.rules = append(f.rules, r)
	}
	return f, nil
}

// Ignored returns true if the given port matches any filter rule.
func (f *Filter) Ignored(port uint16) bool {
	for _, r := range f.rules {
		if port >= r.Low && port <= r.High {
			return true
		}
	}
	return false
}

func parseRule(s string) (Rule, error) {
	s = strings.TrimSpace(s)
	if idx := strings.IndexByte(s, '-'); idx >= 0 {
		lo, err1 := parsePort(s[:idx])
		hi, err2 := parsePort(s[idx+1:])
		if err1 != nil || err2 != nil {
			return Rule{}, fmt.Errorf("bad range")
		}
		if lo > hi {
			return Rule{}, fmt.Errorf("low > high")
		}
		return Rule{Low: lo, High: hi}, nil
	}
	p, err := parsePort(s)
	if err != nil {
		return Rule{}, err
	}
	return Rule{Low: p, High: p}, nil
}

func parsePort(s string) (uint16, error) {
	n, err := strconv.ParseUint(strings.TrimSpace(s), 10, 16)
	if err != nil {
		return 0, err
	}
	return uint16(n), nil
}
