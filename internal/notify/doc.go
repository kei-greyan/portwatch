// Package notify provides pluggable notification backends for portwatch.
//
// The core abstraction is the Notifier interface:
//
//	type Notifier interface {
//		Notify(a alert.Alert) error
//	}
//
// Built-in implementations:
//   - LogNotifier – writes human-readable lines to any io.Writer.
//   - Multi        – fans out to multiple Notifiers, collecting errors.
//
// Use NewMulti to combine several backends so that every alert is
// delivered to all configured destinations.
package notify
