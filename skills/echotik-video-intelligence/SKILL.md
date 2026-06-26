---
name: echotik-video-intelligence
version: 0.2.0
description: "Use for TikTok video intelligence through EchoTik: video search, details, related products, comments, trends, media URLs, and video leaderboards."
metadata:
  requires:
    bins: ["echotik"]
---

# EchoTik Video Intelligence

First read `../echotik-shared/SKILL.md`.

## Basic

```bash
echotik video basic search --keyword "sunscreen" --region US
echotik video basic detail --video-id "123"
```

## Analytics

```bash
echotik video analytics products --video-id "123"
echotik video analytics comments --video-id "123"
echotik video analytics trends --video-id "123"
echotik video analytics media --url "https://www.tiktok.com/..."
```

Use `video analytics media` before downloading video assets, then pass the returned artifact URL to `echotik media download`.

## Leaderboard

```bash
echotik video leaderboard top --region US --date 2026-01-01 --type daily
```
