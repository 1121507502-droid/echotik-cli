---
name: echotik-live-intelligence
version: 0.2.0
description: "Use for TikTok live stream intelligence through EchoTik: live search and live stream details."
metadata:
  requires:
    bins: ["echotik"]
---

# EchoTik Live Intelligence

First read `../echotik-shared/SKILL.md`.

## Basic

```bash
echotik live basic search --keyword "beauty" --region US
echotik live basic detail --room-id "123" --creator-id "456"
```

## Analytics

Live analytics commands are reserved. Prefer product, shop, or creator analytics when tracing live relationships:

```bash
echotik product analytics lives --product-id "123"
echotik shop analytics lives --seller-id "123"
echotik creator analytics lives --creator-id "123"
```

## Leaderboard

No live leaderboard endpoint is exposed in the current EchoTik docs. Use live search plus product/shop/creator analytics as the fallback.
