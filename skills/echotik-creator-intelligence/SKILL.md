---
name: echotik-creator-intelligence
version: 0.2.0
description: "Use for TikTok creator intelligence through EchoTik: creator discovery, details, products, videos, lives, trends, and creator leaderboards."
metadata:
  requires:
    bins: ["echotik"]
---

# EchoTik Creator Intelligence

First read `../echotik-shared/SKILL.md`.

## Basic

```bash
echotik creator basic search --keyword "beauty" --region US
echotik creator basic detail --creator-id "123"
echotik creator basic detail --unique-id "handle"
```

## Analytics

```bash
echotik creator analytics products --creator-id "123"
echotik creator analytics videos --creator-id "123"
echotik creator analytics lives --creator-id "123"
echotik creator analytics trends --creator-id "123" --start-date 2026-01-01 --end-date 2026-01-31
```

Use analytics for creator-product fit, content performance, live commerce activity, and growth tracking.

## Leaderboard

```bash
echotik creator leaderboard top --region US --date 2026-01-01 --type daily
```

Use leaderboard first when the user asks to discover creators rather than analyze a known creator.
