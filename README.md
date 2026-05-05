# logslice

Fast structured log filter and time-range extractor for large log files.

## Installation

```bash
go install github.com/yourusername/logslice@latest
```

## Usage

Extract logs within a time range:

```bash
logslice --from "2024-01-15T08:00:00Z" --to "2024-01-15T09:00:00Z" --file app.log
```

Filter by log level and output as JSON:

```bash
logslice --level error --from "2024-01-15T08:00:00Z" --file app.log --format json
```

Pipe from stdin:

```bash
cat app.log | logslice --level warn --from "2024-01-15T08:00:00Z"
```

### Flags

| Flag | Description | Default |
|------|-------------|---------|
| `--file` | Path to log file | stdin |
| `--from` | Start of time range (RFC3339) | — |
| `--to` | End of time range (RFC3339) | now |
| `--level` | Filter by log level (info, warn, error) | all |
| `--format` | Output format: `text` or `json` | text |
| `--field` | Filter by arbitrary key=value field | — |

### Example Output

```
2024-01-15T08:23:11Z ERROR service=auth msg="invalid token" user_id=42
2024-01-15T08:45:02Z ERROR service=db msg="connection timeout" retry=3
```

## Performance

`logslice` uses binary search on sorted log files to jump directly to the target time range, making it efficient even on multi-gigabyte log files.

## License

MIT © 2024 yourusername