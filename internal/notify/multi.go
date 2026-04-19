package notify

import (
	"errors"

	"github.com/user/portwatch/internal/alert"
)

// Multi fans out a single alert to multiple Notifiers.
type Multi struct {
	notifiers []Notifier
}

// NewMulti returns a Multi that dispatches to each provided Notifier.
func NewMulti(notifiers ...Notifier) *Multi {
	return &Multi{notifiers: notifiers}
}

// Notify sends the alert to all notifiers, collecting any errors.
func (m *Multi) Notify(a alert.Alert) error {
	var errs []error
	for _, n := range m.notifiers {
		if err := n.Notify(a); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}
