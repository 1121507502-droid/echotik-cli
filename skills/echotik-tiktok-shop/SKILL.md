---
name: echotik-tiktok-shop
version: 0.1.0
description: "Use for TikTok Shop seller intelligence through EchoTik: shop lists, shop rankings, category or region filtering, and competitor discovery."
metadata:
  requires:
    bins: ["echotik"]
---

# EchoTik TikTok Shop Intelligence

First read `../echotik-shared/SKILL.md` for authentication and error handling.

## Shortcuts

```bash
echotik shop +list --region US --category-id 601450 --page-size 20
echotik shop +rank --region US --type weekly --page-size 20
```

## Routing

- Use `shop +list` for broad seller discovery.
- Use `shop +rank` to identify leading or fast-growing shops.

## Agent practice

For competitive analysis, combine shop ranking with product ranking in the same region/category.
