# echotik-cli

Agent-friendly CLI for EchoTik TikTok Shop intelligence APIs.

## Quick Start

Install from npm:

```bash
npm install -g echotik-cli
```

The installer prints a one-time EchoTik pixel welcome when npm shows lifecycle
script output. If your npm hides install scripts, use:

```bash
npm install -g echotik-cli --foreground-scripts
```

```bash
echotik config set-credential
echotik auth status
echotik product +search --keyword "sunscreen" --region US
echotik product +list --region US --page-size 20
echotik product +rank --region US --type daily
echotik shop +list --region US
echotik shop +rank --region US --type weekly
```

You can also use environment variables:

```bash
export ECHOTIK_USERNAME="..."
export ECHOTIK_PASSWORD="..."
export ECHOTIK_BASE_URL="https://open.echotik.live"
```

## Commands

- `echotik config set-credential` stores Basic Auth credentials.
- `echotik auth status` validates credentials with a lightweight product-list request.
- `echotik doctor` checks local setup.
- `echotik api <method> <path>` calls raw EchoTik endpoints.
- `echotik product +search` searches realtime TikTok Shop products. `--keyword` is sent to EchoTik as the required `sk` query parameter.
- `echotik product +list` reads the offline product library.
- `echotik product +rank` reads product rankings.
- `echotik shop +list` reads the offline seller library.
- `echotik shop +rank` reads seller rankings.
- `echotik welcome` shows the EchoTik pixel logo animation.

## Output Contract

Successful commands print:

```json
{
  "ok": true,
  "data": {},
  "meta": {},
  "source": "echotik"
}
```

Failed commands print:

```json
{
  "ok": false,
  "error": {
    "type": "authentication_error",
    "message": "missing EchoTik username or password",
    "hint": "run: echotik config set-credential"
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

This repository starts with the same product shape as agent-first CLIs such as
`lark-cli`: native CLI core, npm wrapper, JSON envelopes, high-level shortcuts,
raw API fallback, and bundled agent skills.

## GitHub Release

Create a tag to trigger the release workflow:

```bash
git tag v0.1.0
git push origin main --tags
```

The workflow builds binaries for macOS, Linux, and Windows and uploads them to
the GitHub Release.
