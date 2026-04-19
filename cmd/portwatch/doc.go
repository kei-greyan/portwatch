// Command portwatch is a lightweight daemon that monitors open TCP/UDP ports
// and emits structured alerts whenever ports are opened or closed.
//
// Usage:
//
//	portwatch [-config <path>]
//
// Flags:
//
//	-config  Optional path to a JSON configuration file. When omitted, built-in
//	         defaults are used (1 minute scan interval, state stored in
//	         ~/.portwatch/state.json).
//
// Signals:
//
//	SIGINT / SIGTERM  Gracefully stop the daemon.
//
// Alerts are written as JSON lines to stdout; redirect or pipe them to your
// preferred log aggregator.
package main
