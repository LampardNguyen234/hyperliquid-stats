# Hyperliquid Stats

A comprehensive command-line tool for fetching and analyzing Hyperliquid cryptocurrency volume data and vault statistics.

## Installation

### Prerequisites

- Go 1.24.4 or higher
- [Golang installation guide](https://go.dev/doc/install)

### From Source

```bash
git clone https://github.com/LampardNguyen234/hyperliquid-stats.git
cd hyperliquid-stats
go mod tidy
go build -o hyperliquid-stats
```

### Quick Test

```bash
./hyperliquid-stats --help
```

## Commands Overview

| Command | Aliases | Description |
|---------|---------|-------------|
| `largest-volume` | `large-vol`, `lvol` | Fetch largest users by USD volume |
| `largest-trade-count` | `largest-trades`, `ltc` | Fetch largest users by trade count |
| `daily-volume` | `daily`, `dvol` | Fetch daily USD volume data |
| `daily-volume-by-user` | `daily-by-user`, `duvol` | Fetch daily USD volume for specific users |
| `get-vault` | `vaults`, `vault` | Fetch open vaults with HLP priority |
| `vault-volume` | `vault-vol`, `vvol` | Fetch comprehensive vault volume information |

## Usage Examples

### Basic Usage

```bash
# Get top 10 largest users by volume
./hyperliquid-stats largest-volume --count 10

# Get largest users by trade count
./hyperliquid-stats largest-trade-count --count 15

# Get daily volume data for the last 7 days
./hyperliquid-stats daily-volume --range 7D

# Get all open vaults (HLP vaults shown first)
./hyperliquid-stats get-vault --count 50
```

### Advanced Vault Analysis

```bash
# Get comprehensive vault volume summary
./hyperliquid-stats vault-volume --summary

# Get top 20 vaults sorted by daily volume (HLP priority maintained)
./hyperliquid-stats vault-volume --count 20 --sort-by day

# Get only HLP vault volumes with 10 concurrent workers
./hyperliquid-stats vault-volume --hlp --workers 10

# Get specific vault volume data
./hyperliquid-stats vault-volume --address 0x123...abc
```

### Date Range Filtering

```bash
# Get daily volume data for last 30 days
./hyperliquid-stats daily-volume --range 30D

# Get daily volume for specific date range
./hyperliquid-stats daily-volume --from-date 2024-01-01 --to-date 2024-01-31

# Get user-specific volume for last 3 months
./hyperliquid-stats daily-volume-by-user --user 0x123...abc --range 3M
```

### Vault Filtering and Sorting

```bash
# Get vaults with TVL above $100,000
./hyperliquid-stats get-vault --min-tvl 100000

# Get vault volumes sorted by all-time volume
./hyperliquid-stats vault-volume --sort-by all-time --count 25

# Summary of vault ecosystem
./hyperliquid-stats vault-volume --summary
```

## Global Configuration

### Command Line Flags

```bash
# Custom API endpoints
./hyperliquid-stats --base-url "https://custom-api.com" --info-url "https://custom-info.com" [command]

# Using config file
./hyperliquid-stats --config /path/to/config.yaml [command]
```

### Configuration File

Create `~/.hype-stats.yaml` (or specify with `--config`):

```yaml
base_url: "https://d2v1fiwobg9w6.cloudfront.net"
info_url: "https://api.hyperliquid.xyz/info"
format: "table"
```

## Command Reference

### `largest-volume`

Fetch and display the largest users by USD volume.

```bash
./hyperliquid-stats largest-volume [flags]
```

**Flags:**
- `-c, --count int`: Number of users to display (default: 25)

**Example:**
```bash
./hyperliquid-stats largest-volume --count 50
```

### `largest-trade-count`

Fetch and display the largest users by trade count.

```bash
./hyperliquid-stats largest-trade-count [flags]
```

**Flags:**
- `-c, --count int`: Number of users to display (default: 25)

**Example:**
```bash
./hyperliquid-stats largest-trades --count 100
```

### `daily-volume`

Fetch and display daily USD volume data.

```bash
./hyperliquid-stats daily-volume [flags]
```

**Flags:**
- `-c, --count int`: Number of entries to display (default: 25)
- `-r, --range string`: Time range (e.g., 7D, 30D, 3M, 1Y)
- `--from-date string`: Start date (YYYY-MM-DD format)
- `--to-date string`: End date (YYYY-MM-DD format)
- `-s, --sort string`: Sort order - "asc" or "desc" (default: "desc")

**Examples:**
```bash
# Last 7 days of volume data
./hyperliquid-stats daily-volume --range 7D

# Specific date range
./hyperliquid-stats daily-volume --from-date 2024-01-01 --to-date 2024-01-31

# Last 50 entries, sorted ascending
./hyperliquid-stats daily-volume --count 50 --sort asc
```

### `daily-volume-by-user`

Fetch and display daily USD volume data for a specific user.

```bash
./hyperliquid-stats daily-volume-by-user [flags]
```

**Flags:**
- `-u, --user string`: Filter data for specific user
- `-c, --count int`: Number of entries to display (default: 25)
- `-r, --range string`: Time range (e.g., 7D, 30D, 3M, 1Y)
- `--from-date string`: Start date (YYYY-MM-DD format)
- `--to-date string`: End date (YYYY-MM-DD format)
- `-s, --sort string`: Sort order - "asc" or "desc" (default: "desc")

**Examples:**
```bash
# Get user's volume for last 30 days
./hyperliquid-stats daily-volume-by-user --user 0x123...abc --range 30D

# Get user's volume for specific date range
./hyperliquid-stats duvol --user 0x123...abc --from-date 2024-01-01 --to-date 2024-01-31
```

### `get-vault`

Fetch and display open vaults with HLP priority sorting.

```bash
./hyperliquid-stats get-vault [flags]
```

**Flags:**
- `-c, --count int`: Number of vaults to display (default: 100)
- `--min-tvl float`: Minimum TVL threshold (default: 50,000)
- `--desc`: Sort TVL in descending order (default: true)

**Examples:**
```bash
# Get top 50 vaults by TVL
./hyperliquid-stats get-vault --count 50

# Get vaults with TVL above $1M
./hyperliquid-stats vaults --min-tvl 1000000

# Get all qualifying vaults
./hyperliquid-stats vault --count 0
```

### `vault-volume`

Fetch and display comprehensive vault volume information.

```bash
./hyperliquid-stats vault-volume [flags]
```

**Flags:**
- `--address string`: Specific vault address to fetch
- `-c, --count int`: Number of vaults to display (0 for all)
- `--hlp`: Show only HLP vaults
- `--sort-by string`: Sort by "tvl", "day", "week", "month", "all-time" (default: "tvl")
- `--summary`: Display aggregated summary with totals and top 10
- `-w, --workers int`: Number of concurrent workers (default: 5)

**Examples:**
```bash
# Get comprehensive summary of vault ecosystem
./hyperliquid-stats vault-volume --summary

# Get top 20 vaults by daily volume
./hyperliquid-stats vault-volume --sort-by day --count 20

# Get only HLP vault data with high concurrency
./hyperliquid-stats vvol --hlp --workers 20

# Get specific vault details
./hyperliquid-stats vault-volume --address 0x123...abc

# Get all vault volumes (may take time)
./hyperliquid-stats vault-volume --count 0 --workers 10
```

## Summary Mode Features

The `--summary` flag for `vault-volume` provides:

### HLP Vaults Summary
- Total volumes across all time periods (day, week, month, all-time)
- Total TVL of all HLP vaults
- Count of HLP vaults included

### Non-HLP Vaults Summary  
- Total volumes across all time periods
- Total TVL of all regular vaults
- Count of regular vaults included

### Top 10 Vaults by TVL
- Ranked list of highest TVL vaults
- HLP vaults prioritized in ranking
- Shows TVL and daily volume for context
- All values displayed in millions ($M)

**Example Summary Output:**
```bash
./hyperliquid-stats vault-volume --summary
```

## Performance and Concurrency

### Rate Limiting
- Automatic retry logic
- Built-in rate limiting compliance
- Error handling for temporary failures

## Architecture

### Project Structure
```
├── cmd/                    # CLI commands
│   ├── root.go            # Root command with global flags
│   ├── largest.go         # Largest volume users
│   ├── largest_trade_count.go  # Largest trade count users  
│   ├── daily.go           # Daily volume by user
│   ├── daily_volume.go    # Daily volume aggregate
│   ├── get_vault.go       # Vault listing
│   └── vault_volume.go    # Vault volume analysis
├── internal/
│   ├── api/               # API client and data types
│   │   ├── client.go      # HTTP client with concurrent support
│   │   ├── types.go       # Response structures
│   │   └── vault_volume.go # Vault-specific data types
│   └── config/            # Configuration management
│       └── config.go      # Config struct and defaults
├── pkg/
│   └── common/            # Shared utilities
│       └── table_formatter.go  # Table formatting wrapper
└── main.go               # Entry point
```
