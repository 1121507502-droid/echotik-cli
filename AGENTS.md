# AGENTS.md

## Build

```bash
go test ./...
go build -o bin/echotik .
```

## Project Shape

- `cmd/` contains Cobra commands.
- `internal/client` owns EchoTik HTTP calls and Basic Auth.
- `internal/core` owns config resolution.
- `internal/output` owns JSON envelopes and typed CLI errors.
- `skills/` contains Codex skills copied into `~/.codex/skills`.

## Rules

- Do not print credentials.
- Keep stdout machine-readable JSON for commands.
- Return non-zero exit codes for structured errors.
- Prefer agent-friendly shortcuts over exposing raw endpoint details.
