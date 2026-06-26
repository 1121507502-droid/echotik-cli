---
name: echotik-shop-intelligence
version: 0.2.0
description: "Use for TikTok Shop seller intelligence through EchoTik: shop discovery, shop details, products, creators, videos, lives, trends, and shop leaderboards."
metadata:
  requires:
    bins: ["echotik"]
---

# EchoTik Shop Intelligence

First read `../echotik-shared/SKILL.md`.

## Basic

Use when the user needs to find or identify shops.

```bash
echotik shop basic search --keyword "skincare" --region US
echotik shop basic list --region US --category-id 601450 --page-size 20
echotik shop basic detail --seller-id "123"
```

## Analytics

Use for competitor analysis, assortment analysis, content mix, and sales-driver tracing.

```bash
echotik shop analytics products --seller-id "123"
echotik shop analytics creators --seller-id "123"
echotik shop analytics videos --seller-id "123"
echotik shop analytics lives --seller-id "123"
echotik shop analytics trends --seller-id "123" --start-date 2026-01-01 --end-date 2026-01-31
```

Use `relations` to connect the shop to products, creators, videos, and live streams.

## Leaderboard

Use when the user asks for leading shops, fast-growing shops, or market competitors.

```bash
echotik shop leaderboard top --region US --date 2026-01-01 --type weekly --page-size 20
```
