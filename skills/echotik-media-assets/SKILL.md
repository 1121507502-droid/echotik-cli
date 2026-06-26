---
name: echotik-media-assets
version: 0.2.0
description: "Use for EchoTik media assets: resolving image/video/TOS URLs, downloading assets, and preparing local files for multimodal analysis."
metadata:
  requires:
    bins: ["echotik"]
---

# EchoTik Media Assets

First read `../echotik-shared/SKILL.md`.

## Resolve

Use when a command returns a cover, playback, download, or TOS URL and the agent needs an artifact object.

```bash
echotik media resolve --url "https://..."
```

## Download

Use when the user asks to save images/videos locally or when a multimodal analysis needs local files.

```bash
echotik media download --url "https://..." --output ./assets
```

If a TikTok/TOS URL fails, it may have expired. Re-run the upstream EchoTik command to obtain a fresh URL, then download again.
