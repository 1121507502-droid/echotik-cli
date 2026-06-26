# echotik-cli

Agent-friendly CLI for EchoTik TikTok data APIs.

## Install

```bash
npm install -g echotik-cli
```

To show the install-time welcome output:

```bash
npm install -g echotik-cli --foreground-scripts
```

## Setup

```bash
echotik config set-credential
echotik auth status
echotik doctor
```

Credentials can also come from:

```bash
export ECHOTIK_USERNAME="..."
export ECHOTIK_PASSWORD="..."
export ECHOTIK_BASE_URL="https://open.echotik.live"
```

## Command Model

Commands use:

```bash
echotik <entity> <capability> <operation>
```

Capabilities:

- `basic`: discovery and object details.
- `analytics`: related objects, trends, and relationship analysis.
- `leaderboard`: ranking and opportunity discovery.

Examples:

```bash
echotik product basic search --keyword "sunscreen" --region US --count 10
echotik product basic list --region US --page-size 20
echotik product basic detail --product-id "123"
echotik product analytics creators --product-id "123"
echotik product analytics videos --product-id "123"
echotik product analytics lives --product-id "123"
echotik product analytics trends --product-id "123" --start-date 2026-01-01 --end-date 2026-01-31
echotik product leaderboard top --region US --date 2026-01-01 --type daily

echotik shop basic search --keyword "skincare" --region US
echotik shop basic list --region US --page-size 20
echotik shop basic detail --seller-id "123"
echotik shop analytics products --seller-id "123"
echotik shop analytics creators --seller-id "123"
echotik shop analytics videos --seller-id "123"
echotik shop analytics lives --seller-id "123"
echotik shop analytics trends --seller-id "123" --start-date 2026-01-01 --end-date 2026-01-31
echotik shop leaderboard top --region US --date 2026-01-01 --type weekly

echotik creator basic search --keyword "beauty" --region US
echotik creator basic detail --creator-id "123"
echotik creator analytics products --creator-id "123"
echotik creator analytics videos --creator-id "123"
echotik creator analytics lives --creator-id "123"
echotik creator analytics trends --creator-id "123" --start-date 2026-01-01 --end-date 2026-01-31
echotik creator leaderboard top --region US --date 2026-01-01

echotik video basic search --keyword "sunscreen" --region US
echotik video basic detail --video-id "123"
echotik video analytics products --video-id "123"
echotik video analytics comments --video-id "123"
echotik video analytics trends --video-id "123"
echotik video analytics media --url "https://www.tiktok.com/..."
echotik video leaderboard top --region US --date 2026-01-01

echotik live basic search --keyword "beauty" --region US
echotik live basic detail --room-id "123" --creator-id "456"

echotik media resolve --url "https://..."
echotik media download --url "https://..." --output ./assets
```

The old `+search`, `+list`, and `+rank` commands were removed in `v0.2.0`.

## Agent Output Contract

Successful commands print:

```json
{
  "ok": true,
  "data": {
    "records": [],
    "entities": [],
    "relations": [],
    "artifacts": [],
    "raw": {}
  },
  "meta": {
    "entity": "product",
    "capability": "basic",
    "operation": "search",
    "freshness": "realtime",
    "path": "/api/...",
    "params": {}
  },
  "source": "echotik"
}
```

Failed commands print structured errors:

```json
{
  "ok": false,
  "error": {
    "type": "validation_error",
    "message": "--keyword is required",
    "hint": "example: echotik product basic search --keyword sunscreen --region US"
  },
  "source": "echotik"
}
```

## Development

```bash
go mod tidy
go test ./...
go build -o bin/echotik .
node scripts/run.js doctor
```
