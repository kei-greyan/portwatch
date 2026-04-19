package notify_test

import (
	"bytes"
	"errors"
	"testing"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/notify"
)

type failNotifier struct{ err error }

func (f *failNotifier) Notify(_ alert.Alert) error { return f.err }

func TestMulti_DispatchesToAll(t *testing.T) {
	var b1, b2 bytes.Buffer
	m := notify.NewMulti(notify.New(&b1), notify.New(&b2))

	a := alert.PortOpened(9000, "tcp")
	if err := m.Notify(a); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if b1.Len() == 0 || b2.Len() == 0 {
		t.Error("expected both notifiers to receive the alert")
	}
}

func TestMulti_CollectsErrors(t *testing.T) {
	sentinel := errors.New("send failed")
	m := notify.NewMulti(
		&failNotifier{err: sentinel},
		&failNotifier{err: sentinel},
	)

	err := m.Notify(alert.PortOpened(80, "tcp"))
	if err == nil {
		t.Fatal("expected error")
	}
	if !errors.Is(err, sentinel) {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestMulti_PartialFailure_ContinuesDelivery(t *testing.T) {
	var buf bytes.Buffer
	m := notify.NewMulti(
		&failNotifier{err: errors.New("oops")},
		notify.New(&buf),
	)

	_ = m.Notify(alert.PortOpened(443, "tcp"))

	if buf.Len() == 0 {
		t.Error("second notifier should still receive alert despite first failing")
	}
}
