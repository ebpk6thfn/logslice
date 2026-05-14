# logslice

Fast log file slicer that extracts time-bounded segments from large log files without loading them fully into memory.

---

## Installation

```bash
go install github.com/yourusername/logslice@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/logslice.git
cd logslice
go build -o logslice .
```

## Usage

```bash
logslice [flags] <logfile>
```

### Flags

| Flag | Description | Example |
|------|-------------|---------|
| `--from` | Start timestamp (RFC3339) | `2024-01-15T08:00:00Z` |
| `--to` | End timestamp (RFC3339) | `2024-01-15T09:00:00Z` |
| `--format` | Log timestamp format | `2006-01-02T15:04:05` |
| `--out` | Output file (default: stdout) | `slice.log` |

### Example

Extract one hour of logs from a large application log:

```bash
logslice --from 2024-01-15T08:00:00Z --to 2024-01-15T09:00:00Z app.log
```

Write the output to a file:

```bash
logslice --from 2024-01-15T08:00:00Z --to 2024-01-15T09:00:00Z --out slice.log app.log
```

## How It Works

`logslice` uses binary search over the log file to locate the start and end positions of the target time range, then streams only the matching lines — no matter how large the file is.

## Requirements

- Go 1.21+

## License

MIT © 2024 yourusername