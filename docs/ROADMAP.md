# EchoTik CLI + Skills Agent Roadmap

## Current Status

| Area | Status | Notes |
|---|---|---|
| npm distribution | Implemented | `echotik-cli` publishes a small npm package and downloads platform binaries from GitHub Releases in `postinstall`. |
| Basic auth and config | Implemented | `echotik config set-credential`, env credentials, `auth status`, and `doctor` are available. |
| Command model | Implemented | New command shape is `echotik <entity> <capability> <operation>`. Old `+search/+list/+rank` commands were removed in `v0.2.0`. |
| Product basic | Implemented | `product basic search/list/detail`. Search maps `--keyword` to EchoTik `sk`. |
| Product analytics | Implemented | `product analytics creators/videos/lives/trends`. Relations are emitted where a source product ID is known. |
| Product leaderboard | Implemented | `product leaderboard top`. |
| Shop basic | Implemented | `shop basic search/list/detail`. Search uses EchoTik general search. |
| Shop analytics | Implemented | `shop analytics products/creators/videos/lives/trends`. Relations are emitted where a source seller ID is known. |
| Shop leaderboard | Implemented | `shop leaderboard top`. |
| Creator basic | Implemented | `creator basic search/detail`. |
| Creator analytics | Implemented | `creator analytics products/videos/lives/trends`. |
| Creator leaderboard | Implemented | `creator leaderboard top`. |
| Video basic | Implemented | `video basic search/detail`. |
| Video analytics | Implemented | `video analytics products/comments/trends/media`. |
| Video leaderboard | Implemented | `video leaderboard top`. |
| Live basic | Implemented | `live basic search/detail`. |
| Live analytics | Partially implemented | Command placeholders exist, but return `unsupported_operation` until EchoTik exposes or we map reliable endpoints. |
| Live leaderboard | Not implemented | Placeholder returns `unsupported_operation`; no confirmed EchoTik live leaderboard endpoint yet. |
| Category entity | Not implemented | Category endpoints exist in docs, but no `category basic/analytics/leaderboard` command group yet. |
| Media resolve/download | Implemented | `media resolve` and `media download` produce artifacts and local manifests. |
| Unified output contract | Implemented | New commands return `records`, `entities`, `relations`, `artifacts`, and `raw`. |
| Entity/relation normalization | Partially implemented | Generic ID extraction exists. More endpoint-specific normalization is still needed for nested realtime payloads. |
| Bundled skills | Implemented | Product/shop/creator/video/live/media/shared skills are bundled in npm package. |
| Skill installer | Implemented | `echotik skills list/path/install codex` installs bundled skills into Codex. |
| High-level analyze workflows | Not implemented | No `echotik analyze ...` commands yet. Agents currently compose lower-level commands through skills. |
| Welcome banner | Implemented | Installation and `echotik welcome` show a terminal banner using compatible ANSI colors. |

## Direction

EchoTik CLI should become an agent-friendly TikTok data execution layer, not just a thin wrapper around API endpoints.

The central design is:

```text
User question
  -> Agent skill chooses an entity and capability
  -> CLI validates parameters and calls EchoTik APIs
  -> CLI normalizes records, entities, relations, and artifacts
  -> Agent analyzes the structured result
```

Agents should not need to memorize raw EchoTik endpoint names or fragile API parameter names. The CLI owns parameter mapping, validation, pagination defaults, media download handling, and the output contract.

## Entity And Capability Model

Every major EchoTik data object follows the same three capability layers:

```text
basic       discovery and object details
analytics   relationships, trends, and driver analysis
leaderboard rankings and opportunity discovery
```

Entities:

```text
product
shop
creator
video
live
category
media
```

Agent routing rule:

```text
Find opportunities      -> leaderboard or basic search
Understand one object   -> basic detail
Explain performance     -> analytics
Fetch media assets      -> analytics media or media download
Build a report          -> leaderboard -> basic -> analytics -> synthesis
```

## CLI Target Shape

The long-term CLI shape is:

```bash
echotik product basic search
echotik product basic list
echotik product basic detail
echotik product analytics creators
echotik product analytics videos
echotik product analytics lives
echotik product analytics trends
echotik product leaderboard top

echotik shop basic search
echotik shop basic list
echotik shop basic detail
echotik shop analytics products
echotik shop analytics creators
echotik shop analytics videos
echotik shop analytics lives
echotik shop analytics trends
echotik shop leaderboard top

echotik creator basic search
echotik creator basic detail
echotik creator analytics products
echotik creator analytics videos
echotik creator analytics lives
echotik creator analytics trends
echotik creator leaderboard top

echotik video basic search
echotik video basic detail
echotik video analytics products
echotik video analytics comments
echotik video analytics trends
echotik video analytics media
echotik video leaderboard top

echotik live basic search
echotik live basic detail
echotik live analytics products
echotik live analytics trends
echotik live analytics media
echotik live leaderboard top

echotik category basic list
echotik category basic children

echotik media resolve
echotik media download
```

The current implementation covers most of this shape except category and fully mapped live analytics/leaderboard.

## Output Contract

All agent-facing data commands should return:

```json
{
  "ok": true,
  "data": {
    "records": [],
    "entities": [],
    "relations": [],
    "artifacts": [],
    "raw": {}
  },
  "meta": {
    "entity": "product",
    "capability": "basic",
    "operation": "search",
    "freshness": "realtime",
    "path": "/api/...",
    "params": {}
  },
  "source": "echotik"
}
```

Field intent:

| Field | Purpose |
|---|---|
| `records` | The main result set for the current command. |
| `entities` | Extracted `product`, `shop`, `creator`, `video`, `live`, `category`, or `media` objects. |
| `relations` | Cross-entity links such as product-to-video or shop-to-product. |
| `artifacts` | Media URLs, local files, covers, videos, and download metadata. |
| `raw` | Original EchoTik response to avoid losing fields during normalization. |

This contract is the main bridge between CLI execution and agent analysis.

## Relation Model

Relations should be explicit whenever a command starts from a known object ID.

Core relation types:

```text
product -> shop
product -> creator
product -> video
product -> live

shop -> product
shop -> creator
shop -> video
shop -> live

creator -> product
creator -> video
creator -> live

video -> product
video -> comment
video -> media

live -> product
live -> creator
live -> media
```

Current normalization extracts common IDs generically. Future work should add endpoint-specific extraction for deeply nested realtime payloads, especially product realtime search and video/media payloads.

## Skills Strategy

Skills should teach agents task paths, not raw API endpoints.

Bundled skills:

```text
echotik-shared
echotik-product-intelligence
echotik-shop-intelligence
echotik-creator-intelligence
echotik-video-intelligence
echotik-live-intelligence
echotik-media-assets
```

Each intelligence skill should keep the same sections:

```text
basic
analytics
leaderboard
```

Each skill should answer:

```text
When to use this skill
Which CLI command to call
Which parameters are required
How to resolve missing IDs
How to follow relations
When to download media
How to summarize useful business insights
```

The CLI currently provides:

```bash
echotik skills list
echotik skills path
echotik skills install codex
```

## Media Pipeline

Media must be treated as first-class artifacts because TikTok/TOS URLs may expire.

Current commands:

```bash
echotik media resolve --url "https://..."
echotik media download --url "https://..." --output ./assets
```

Artifact target shape:

```json
{
  "type": "video",
  "sourceUrl": "...",
  "downloadUrl": "...",
  "localPath": "./assets/video.mp4",
  "expiresAt": "...",
  "sha256": "..."
}
```

Next improvements:

```text
batch download
URL expiry detection
refresh expired media URLs through upstream EchoTik commands
deduplication by SHA256
manifest merge instead of overwrite
```

## Roadmap

### Phase 1: Stabilize v0.2.x

- Keep the command model stable.
- Polish `echotik skills install codex`.
- Improve welcome banner compatibility across terminals.
- Add more tests for command tree and output contracts.
- Add endpoint-specific normalization for product/shop/creator/video known responses.

### Phase 2: Complete entity coverage

- Add `category basic list/children` for product and shop category trees.
- Map live analytics if reliable EchoTik endpoints are available.
- Add live leaderboard only after a real endpoint is confirmed.
- Add real-time product/share-link resolution as a first-class helper.

### Phase 3: Improve analysis readiness

- Add richer relation names and endpoint-specific relation extraction.
- Add score-friendly summarized fields for agents, while keeping `raw`.
- Add pagination helpers and batch command patterns.
- Add time range defaults for trend commands where safe.

### Phase 4: Add high-level workflows

Potential commands:

```bash
echotik analyze market --keyword sunscreen --region US
echotik analyze product --product-id ...
echotik analyze shop --seller-id ...
echotik analyze creator --creator-id ...
echotik analyze video --video-id ...
```

These commands should orchestrate lower-level commands and return a report-ready dataset, not just raw API output.

### Phase 5: Third-party agent distribution

- Keep Codex support first.
- Document how other agents can install bundled skills manually.
- Consider installers for other local agent skill directories only after their conventions are stable.

## Acceptance Criteria For Future Work

Every new data command should:

- Follow `echotik <entity> <capability> <operation>`.
- Validate required parameters before making network calls.
- Return the unified data contract.
- Preserve the raw EchoTik response.
- Emit useful `meta.entity`, `meta.capability`, `meta.operation`, `meta.path`, and `meta.params`.
- Add relations when the command is explicitly about a known source object.
- Add artifacts when media URLs are present.
- Be documented in README and the relevant skill.
- Have at least one CLI structure or output-contract test when behavior is non-trivial.

