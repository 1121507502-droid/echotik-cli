---
name: echotik-shared
version: 0.1.0
description: "Use when setting up echotik CLI, configuring EchoTik Basic Auth credentials, checking authentication, handling EchoTik API errors, or deciding whether to use offline or realtime TikTok Shop data."
metadata:
  requires:
    bins: ["echotik"]
---

# echotik shared rules

Use `echotik` for EchoTik TikTok Shop intelligence APIs.

## Setup

Run:

```bash
echotik config set-credential
echotik auth status
```

The CLI uses EchoTik Basic Auth. Credentials can also come from:

```bash
ECHOTIK_USERNAME=...
ECHOTIK_PASSWORD=...
ECHOTIK_BASE_URL=https://open.echotik.live
```

## Output

Commands return JSON envelopes. Treat `ok: true` as success and read `data`. Treat `ok: false` as a structured error and follow `error.hint`.

## Data freshness

- Offline product/shop library commands expose `meta.freshness = "offline_t_plus_1"`.
- Realtime search exposes `meta.freshness = "realtime"` and may need retry/backoff if EchoTik risk control or server errors occur.

## Error handling

- `authentication_error`: ask the user to configure credentials with `echotik config set-credential`.
- `rate_limit`: retry with backoff.
- `server_error` or realtime failures: retry with backoff, or use offline list/ranking commands when realtime freshness is not required.
