// Package monitor orchestrates the port-watch cycle:
// it periodically invokes the scanner, computes a diff against the
// previously persisted state, emits alerts for any changes, and
// writes the new snapshot back to the state store.
//
// Typical usage:
//
//	sc := scanner.New(scanner.Config{Targets: []string{"127.0.0.1"}})
//	st := state.New("/var/lib/portwatch/state.json")
//	al := alert.New(os.Stdout)
//	m  := monitor.New(monitor.Config{Interval: 30 * time.Second}, sc, st, al)
//	if err := m.Run(ctx); err != nil && !errors.Is(err, context.Canceled) {
//		log.Fatal(err)
//	}
package monitor
