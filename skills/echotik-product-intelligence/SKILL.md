---
name: echotik-product-intelligence
version: 0.2.0
description: "Use for TikTok Shop product intelligence through EchoTik: product discovery, product details, product-related creators/videos/lives/trends, and product leaderboards."
metadata:
  requires:
    bins: ["echotik"]
---

# EchoTik Product Intelligence

First read `../echotik-shared/SKILL.md`.

## Basic

Use when the user needs to find or identify products.

```bash
echotik product basic search --keyword "sunscreen" --region US --count 10
echotik product basic list --region US --category-id 601450 --page-size 20
echotik product basic detail --product-id "123"
```

If the user gives a product share link instead of an ID, resolve it through the raw API only if needed, then continue with `product basic detail`.

## Analytics

Use when the user asks why a product is performing, who drives it, or what content/live streams are linked.

```bash
echotik product analytics creators --product-id "123"
echotik product analytics videos --product-id "123"
echotik product analytics lives --product-id "123"
echotik product analytics trends --product-id "123" --start-date 2026-01-01 --end-date 2026-01-31
```

Read `relations` to connect the product to creators, videos, and live streams. Use `artifacts` when video covers or playback URLs are present.

## Leaderboard

Use for opportunity discovery and market scans.

```bash
echotik product leaderboard top --region US --date 2026-01-01 --type daily --page-size 20
```

After a leaderboard result, use `product basic detail` and then analytics commands for the top candidates.
