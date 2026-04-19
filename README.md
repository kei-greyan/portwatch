# portwatch

Lightweight daemon that monitors open ports and alerts on unexpected changes.

## Installation

```bash
go install github.com/yourusername/portwatch@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/portwatch.git && cd portwatch && go build ./...
```

## Usage

Start the daemon with a configuration file:

```bash
portwatch --config /etc/portwatch/config.yaml
```

Example `config.yaml`:

```yaml
interval: 30s
alert:
  email: admin@example.com
baseline:
  - 22
  - 80
  - 443
```

portwatch will scan open ports at the specified interval and send alerts whenever a port outside the baseline is detected open — or an expected port goes missing.

### Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--config` | `config.yaml` | Path to config file |
| `--interval` | `60s` | Scan interval |
| `--verbose` | `false` | Enable verbose logging |

### Example Alert

```
[ALERT] Unexpected port open: 8080 (2024-01-15 03:22:11)
[ALERT] Expected port closed: 443 (2024-01-15 03:22:11)
```

## License

MIT © 2024 portwatch contributors