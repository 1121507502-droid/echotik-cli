---
name: echotik-tiktok-product
version: 0.1.0
description: "Use for TikTok Shop product intelligence through EchoTik: product search, product lists, rankings, category or region filtering, and product-market analysis."
metadata:
  requires:
    bins: ["echotik"]
---

# EchoTik TikTok Product Intelligence

First read `../echotik-shared/SKILL.md` for authentication and error handling.

## Shortcuts

```bash
echotik product +search --keyword "sunscreen" --region US --page-size 20
echotik product +list --region US --category-id 601450 --page-size 20
echotik product +rank --region US --type daily --page-size 20
```

## Routing

- Use `product +search` when the user needs realtime keyword search. The CLI maps `--keyword` to EchoTik's required `sk` query parameter.
- Use `product +list` when broad offline discovery is acceptable or realtime search is unstable.
- Use `product +rank` for best-selling or trending product analysis.

## Agent practice

Prefer region filters for market analysis. If the user does not specify a region, ask for one or default only when the task clearly implies a market.
