# Gap Analysis v2 — Maintainer's Review & Execution Plan

**Reviewed against:** the live codebase at **v3.106.0** (the analysis was written against v3.104).
**Method:** every claim below was checked against source, not accepted on trust. Verification commands are noted so this can be re-run when it goes stale.

---

## Verdict at a glance

| # | Analysis item | Gap real? | My call |
|---|---|---|---|
| 1 | `grit ui list` / `add` | **Premise stale** | **Blocked** — see below |
| 2 | Nielsen pass checklist | **Yes, and worse than stated** | **Do first** — expand it |
| 3 | UUIDv7 swap | Yes | Do — but via a helper, not 31 edits |
| 4 | `golangci-lint` config | Yes — nothing shipped | Do |
| 5 | MCP server | Yes | Do — **promote**, it's the strategic one |
| 6 | PgBouncer in prod compose | Yes | Do |
| 7 | `grit test` | Yes — no such command | Do |
| 8 | Queue depth + SLO tile | Yes | Defer — real work, not "already have the data" |
| 9 | Internal event bus | Yes | **Defer hard** — YAGNI |
| 10 | Request-lifecycle doc | Yes — middleware page has zero lifecycle content | Do (cheap) |
| 11 | Partial/functional indexes | Yes | Defer — no user is near that scale |
| 12 | GraphQL plugin | Yes | Don't build until asked |
| — | **Two live bugs found in session** | **Yes** | **Do first** — ahead of everything |

---

## Where I disagree with the analysis

### 1. Item #1 is blocked, not "the cheapest win"

The analysis says `grit ui add` is "a thin CLI wrapper around the existing `/r.json` endpoint." That endpoint **does not exist in scaffolded projects**. It was deliberately removed in v3.31.78, and `CLAUDE.md` explicitly says *"Do not re-add these to the scaffold."*

```bash
grep -rn '"/r.json"\|/r/:name\|UIRegistry' internal/scaffold/*.go   # → empty
grep -n '"ui"' cmd/grit/main.go                                     # → empty
```

So this isn't a thin wrapper over something that works — it's **"host the registry somewhere"** first. That's a different, larger job, and it needs a decision from you: *is grit-ui hosted at a stable public URL?* Until that's answered, item #1 cannot be top of the list.

### 2. The Nielsen item is underrated — and the evidence is stronger than the document knows

The analysis cites four bugs as heuristic violations. Two of them (`window.prompt` for access reviews, missing back buttons) **I fixed in v3.106.0**. The third — native `confirm()` for destructive actions — is still live in **six places**, and I added one of them *yesterday* while building SSO:

```bash
grep -rn "window.confirm" internal/scaffold/*.go | wc -l   # → 6
```

That's the important finding. A checklist alone would not have stopped me: I had the themed `ConfirmModal` component available and reached for `window.confirm` anyway, because it was one line. **A checklist without fixing the existing violations just documents the drift.** So this item becomes two things: write the checklist *and* replace the six sites.

### 3. UUIDv7 is right, but it is not a "drop-in swap"

The analysis calls it "a drop-in swap at the ID-generation function." In reality there are **31 call sites**, and the APIs differ:

```go
uuid.New()          // returns UUID
uuid.NewV7()        // returns (UUID, error)  ← different signature
```

Doing this as 31 inline edits means 31 new error paths in `BeforeCreate` hooks. The correct shape is a single `models.NewID() string` helper that swallows the error with a v4 fallback, then one mechanical substitution. Also worth stating in the docs, which the analysis omits: **v7 IDs leak creation time** — that is the point, but it's a disclosure change for anyone using IDs in public URLs.

### 4. The event bus should be deferred harder

The analysis ranks it #9 and hedges appropriately, but I'd go further: this is the only item that introduces a **new architectural primitive**. Grit's existing seams (jobs, webhooks, realtime) already work; the pain is hypothetical until someone actually needs three subscribers on one event. Building it now means maintaining a pub/sub layer with no proven consumer. Revisit when a real feature needs it.

### 5. Item 5.4 (`grit-search` / logical replication) is out of scope for this repo

```bash
grep -rn "grit-search\|meilisearch" internal/plugin/*.go   # → empty
```

It lives in the separate plugins repo. Worth a note there, not a line item here.

---

## What the analysis misses

### Two live bugs on the documented happy path

Both found while building SSO, both verified in source. These outrank most of the roadmap because they affect existing users today.

**A. Postgres-only SQL in two workers — fails on every SQLite project.**

```
internal/scaffold/api_jobs_files.go:433
  DELETE FROM users WHERE deleted_at IS NOT NULL AND deleted_at < NOW() - INTERVAL '30 days'
internal/scaffold/api_useractivity_files.go:368
  Where("created_at > NOW() - INTERVAL '24 hours'")
```

Observed live: `SQL logic error: near "'30 days'": syntax error`. SQLite is a **first-class supported database** (it's the quick-start path in CLAUDE.md). On any SQLite project the user-cleanup job fails silently every run, and the 24-hour activity query errors. Fix: compute the cutoff in Go and bind it as a parameter — portable across both drivers.

**B. `User.Active` carries `gorm:"default:true"` — the same class of bug I just fixed in SSO.**

GORM omits zero-valued fields from an INSERT when the column has a default, so `Active: false` cannot be stored on create; it silently comes back `true`. Proven with a diagnostic test this session. Scope is narrower than the SSO one (the admin's *disable* flow uses an update, which works), but creating a user as inactive silently produces an active account. Worth auditing every `default:true` bool in the scaffold at the same time.

### SCIM

The *other* enterprise review called SSO the "single biggest enterprise blocker" and paired it with SCIM. This analysis doesn't mention either. SSO shipped in v3.106.0; SCIM (directory-driven deprovisioning) is its natural follow-on — when someone is removed from the corporate directory, their access should die without an admin remembering. **No user has asked for it.** Flag it, don't build it.

---

## The plan

### Phase 1 — Correctness first (small, evidence-backed)

1. **Fix the two SQLite-breaking queries.** Bind a Go-computed cutoff instead of `NOW() - INTERVAL`. Add a test that runs the cleanup worker against SQLite.
2. **Audit `default:true` booleans.** Fix `User.Active`; sweep for the same pattern elsewhere. This bug class has now bitten twice.
3. **Nielsen pass, done properly.** Ten one-line checks in `GRIT_STYLE_GUIDE.md`, *plus* replace all six `window.confirm` sites with the existing themed `ConfirmModal`. The checklist without the cleanup is theatre.

*Rationale: these are defects, not features. Phase 1 ships as one patch release.*

### Phase 2 — Cheap, high-leverage infrastructure

4. **`.golangci.yml` in every scaffold**, tuned to the *100 Go Mistakes* categories the changelog shows actually happening (goroutine leaks, missing context cancellation, unchecked errors). Wire it into the generated CI workflow.
5. **UUIDv7 behind `models.NewID()`.** One helper, one substitution across 31 sites, a note in the docs about time disclosure.
6. **PgBouncer in `docker-compose.prod.yml`.** Boring, prevents an invisible failure mode for anyone running API + asynq workers.
7. **Request-lifecycle doc.** The middleware page exists but contains no ordered pipeline. One diagram of what runs in what order and where to hook in.

### Phase 3 — The strategic item

8. **`grit mcp serve`.** I'd promote this above the analysis's #5 ranking: Grit's stated differentiator is "vibe-coding ready," and this is the single item that most advances it. But I'd push back on "mostly wiring" — it needs a transport, a tool schema, and an auth story. Scope it as: `routes`, `schema` (via Studio's introspection), `openapi`, and `recent_errors` (via Pulse). Four read-only tools, stdio transport first.
9. **`grit test`.** One command, one report, skipping runners that don't apply to the architecture mode.

### Deferred, with reasons

| Item | Why not now |
|---|---|
| `grit ui list/add` | **Blocked** — needs a hosted registry URL from you |
| Queue depth + SLO tile | Real Pulse work; the data is there but the UI/threshold logic isn't "free" |
| Event bus | No proven second subscriber — YAGNI |
| Partial/functional indexes | Matters at ~100M rows; nobody is close |
| GraphQL plugin | Build when someone asks |
| SCIM | Natural follow-on to SSO, but zero demand |
| Sharding, multi-cloud, gRPC, DI container | Correctly excluded by the analysis |

---

## The one question I need answered

**Is `grit-ui` hosted at a stable public URL** (e.g. `gritframework.dev/r`)?

- **Yes** → `grit ui list/add` becomes genuinely cheap and jumps to Phase 2.
- **No** → the real task is "host the registry," which is a bigger piece of work and should be scoped separately.

Everything else in this plan can proceed without further input.
