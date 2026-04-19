// Package config provides loading, saving, and validation of portwatch
// runtime configuration.
//
// Configuration is stored as a JSON file and supports the following fields:
//
//   - interval      – duration between consecutive port scans (default 30s)
//   - state_path    – path to the persistent state file
//   - alert_webhook – optional HTTP endpoint to receive alert payloads
//   - ignore_ports  – list of ports that will never trigger alerts
//
// Example config file:
//
//	{
//	  "interval": 60000000000,
//	  "state_path": "/var/lib/portwatch/state.json",
//	  "ignore_ports": [22, 80, 443]
//	}
package config
