---
name: grit
description: >
  Grit framework conventions and patterns for this triple project.
  Use when modifying models, handlers, routes, schemas, types, or components.
  Automatically loaded as background knowledge.
user-invocable: false
---

# Grit Framework — Triple (Web + Admin + API) + Next.js

This project uses a **triple** architecture: Go API + React web frontend + admin panel in a Turborepo monorepo. Three apps that share types and schemas.

It uses **Next.js** (App Router) for the frontend — SSR, SSG, ISR, SEO-optimized.

**Batteries included:** file storage (S3), email (Resend), background jobs (asynq), cron, Redis cache, AI (Vercel AI Gateway), security (Sentinel), observability (Pulse), auto-generated API docs (gin-docs).

For detailed API conventions, code patterns, and service documentation, see [reference.md](reference.md).

---

## CLI Commands

`ash
# Code generation
grit generate resource Post --fields "title:string,content:text,published:bool"
grit generate resource Post --from post.yaml
grit remove resource Post               # Cleanly removes all generated files + injections

# Development
grit start                            # Start dev servers
grit sync                             # Go types → TypeScript
grit add role MODERATOR               # Injects role into 7 locations
grit migrate                          # Run GORM AutoMigrate
grit seed                             # Create admin + demo users

# Operations
grit routes                           # List all API routes
grit down                             # Enable maintenance mode (503)
grit up                               # Disable maintenance mode
grit deploy --host user@server --domain myapp.com  # Production deploy

# Updates
grit upgrade                          # Update project to latest templates
grit update                           # Update Grit CLI itself
`

---

## Project Structure

`
library/
├── packages/shared/              # Zod schemas, TS types, constants
├── apps/
│   ├── api/                      # Go backend (Gin + GORM)
│   │   ├── cmd/server/main.go
│   │   └── internal/             # config, database, models, handlers, services, middleware, routes
│   ├── web/                      # Next.js frontend (App Router)
│   └── admin/                    # Next.js admin panel (App Router)
├── .env
├── docker-compose.yml
└── turbo.json
`

**Mounted dashboards** (auto-configured in routes.go):
- `/docs` — API documentation (gin-docs, OpenAPI 3.1)
- `/studio` — Database browser (GORM Studio)
- `/sentinel/ui` — Security dashboard (WAF, rate limiting)
- `/pulse` — Observability (tracing, metrics)

---

## Generating Resources

`ash
grit generate resource Post --fields "title:string,content:text,published:bool,views:int"
`

Creates model, service, handler, schema, types, hooks, and injects into existing files via marker comments.

### Field Types

| Type | Go | TypeScript | Form |
|------|----|-----------|------|
| `string` | `string` | `string` | Text input |
| `text` | `string` | `string` | Textarea |
| `int` / `uint` / `float` | `int` / `uint` / `float64` | `number` | Number input |
| `bool` | `bool` | `boolean` | Toggle |
| `datetime` / `date` | `*time.Time` | `string | null` | Picker |
| `richtext` | `string` | `string` | Tiptap editor |
| `slug` | `string` | `string` | Auto-generated |
| `string_array` | `JSONSlice[string]` | `string[]` | Tag input |
| `email` / `url` / `domain` | `string` | `string` | Email, web address, bare domain |
| `tel` (`phone:tel:UG`) | `string` (E.164) | `string` | Country picker + as-you-type number |
| `country` | `string` (ISO alpha-2) | `string` | Searchable country picker |
| `color` / `time` | `string` | `string` | Colour picker, time input |
| `percent` / `rating` (`score:rating:10`) | `float64` / `int` | `number` | Number with %, stars |
| `json` | `datatypes.JSON` | `unknown` | JSON editor |
| `belongs_to:X` | `uint` (FK) | `number` | Relationship select |
| `many_to_many:X` | Junction table | `number[]` | Multi-select |

**Modifiers:** `:unique`, `:required`, `:optional` (append after type).

---

## Marker Comments

Grit uses marker comments to inject generated code. **Never delete these:**

`go
// grit:models          — models/user.go (AutoMigrate list)
// grit:handlers        — routes/routes.go (handler initialization)
// grit:routes:protected — routes/routes.go (protected route group)
// grit:routes:admin    — routes/routes.go (admin route group)
`

`typescript
// grit:schemas         — schemas/index.ts
// grit:types           — types/index.ts
// grit:api-routes      — constants/index.ts
`
`typescript
// grit:resources       — resources/index.ts (imports)
// grit:resource-list   — resources/index.ts (registry array)
`
---

## Frontend Routing

Routes live in app/. File naming: page.tsx (route), layout.tsx (layout), loading.tsx, error.tsx.
Use 'use client' directive for client components. Server Components by default.
Navigation: import { useRouter } from 'next/navigation', import Link from 'next/link'.

---

## Common Tasks

### Add a field to an existing resource

1. Add field to Go model (apps/api/internal/models/<name>.go)
2. Update handler if field needs special handling
3. Update Zod schema (packages/shared/schemas/<name>.ts)
4. Update TypeScript type (packages/shared/types/<name>.ts)
5. Update admin resource (apps/admin/resources/<name>.ts) — add column + form field
6. Restart API (GORM auto-migrates)

### Add a new API endpoint

1. Create/update handler in the handlers directory
2. Register route in routes.go
3. Create React Query hook if frontend needs it

### Add a relationship

`go
type Post struct {
    CategoryID uint     // json:"category_id"
    Category   Category // gorm:"foreignKey:CategoryID" json:"category,omitempty"
}
// In handler: query.Preload("Category").Find(&posts)
`

---

## Performance & Production Hygiene

Grit ships with audited primitives — but app code can still introduce
CPU burn / memory leaks / lock contention. Follow these rules so the
patterns the framework hands you stay efficient at production scale.

### Hot-path rules (request handlers + middleware)

- **Never spawn unbounded goroutines per-request.** If you need
  background work, push to a buffered channel + a fixed-size worker
  pool. Pattern shipped: `ActivityLogger` in
  `internal/middleware/activity.go` — single writer, bounded
  channel, drop on overflow.
- **Never hold a mutex across slow operations.** Read shared state
  under the lock, copy what you need, release, then do the work.
  Pattern shipped: `flags.Engine.evaluate` — copies the flag
  struct under RLock, runs all rules unlocked.
- **Never load a whole table into memory.** Use
  `paginate.List` for list endpoints,
  `FindInBatches` for exports, `sync.VerifyChain`
  pattern (cursor on (created_at, id)) for verification scans.
- **Never call SHA-256 in a hot path** unless cryptographic. Cache
  keys, ID buckets, and hash maps should use `fnv.New64a` or
  `xxhash` — a 50x speedup with no correctness loss. Pattern
  shipped: cache middleware uses fnv.
- **Never write to the same row from multiple goroutines** without a
  serialization point. The audit hash chain serializes by design via
  the single-writer worker. If you need ad-hoc chained writes, use
  `audit.AppendChained` (it takes the FOR UPDATE lock).
- **Never use `time.Now().UnixNano() % N` for randomness** — it's
  biased toward the call frequency. Use `crypto/rand` or
  `math/rand`. Pattern shipped: `flags.bucketFor` for
  anonymous users uses crypto/rand.

### Database rules

- **N+1 queries**: always `Preload("Association")` when you'll
  access an association inside a loop. The resource generator emits
  Preloads for declared `belongs_to` fields automatically.
- **No SELECT * on large tables**: use `.Select("id", "name")`
  when you only need a few columns.
- **No queries inside loops**: batch with `IN(?)` or
  `CreateInBatches`. The webhook receiver and sync push handler
  do batch inserts; mirror that pattern.
- **Indexes**: every column in a `WHERE` or `ORDER BY`
  on a hot endpoint needs an index. Use the
  `gorm:"index"` tag on the model.
- **Keep transactions short**: never hold a transaction open across
  network calls (S3 upload, webhook callback, AI completion). Take the
  lock, flush, release, then do the slow work.

### Background job rules

- **Job idempotency**: every handler must be safe to re-run. Use the
  job's own ID as a dedup key (the asynq scaffold does this). For
  webhooks, the `(provider, external_id)` UNIQUE constraint
  enforces it for free.
- **Job timeouts**: pass a `context.WithTimeout` through every
  job handler. Asynq propagates it via `ctx.Done()` — the job
  is killed if it runs over.
- **Distributed locks for cron**: if you run >1 instance of the API,
  cron jobs need a lock or they double-fire. Use Redis `SETNX`
  with a TTL or asynq's `UniqueOpt`.
- **Bounded retries**: asynq retries with exponential backoff by
  default; cap `MaxRetry` explicitly per job type so a poison
  message doesn't burn the queue forever.
- **Per-queue worker limits**: don't share one worker pool across
  `default` + `critical` + `low` queues. Set
  per-queue concurrency in `jobs.go` so a low-priority backlog
  can't starve criticals.

### Logging rules

- **No `log.Printf` in tight loops.** Sample (e.g.
  `if i%1000 == 0 { log… }`) or move to debug level.
- **Never log full request/response bodies** — they may contain PII,
  secrets, or megabyte-scale uploads. The activity logger stores
  SHA-256 of the body, not the body itself. Follow that pattern.
- **Use structured logging at info+** — easier to grep, alert on, and
  ship to a log aggregator.

### Memory rules

- **Caches need a TTL**: never cache without expiration. Redis is the
  default cache; `cache.Set(ctx, key, val, ttl)` requires a TTL.
- **Channels need a bound**: `make(chan T, N)` with a real N,
  never `make(chan T)` for high-volume work.
- **Maps that only grow are leaks** — if you keep a per-user cache,
  evict on a TTL or use an LRU library.

### What's already audited

The Grit framework's own primitives have been performance-reviewed
through v3.22 — the patterns below are safe to call at scale:

- `ActivityLogger` middleware: bounded channel + single writer
- `audit.VerifyChain`: chunked walk with context cancellation
- `paginate.List`: limit-bounded queries, opt-in cursor mode
- `flags.Engine`: in-memory cache, copy-then-release locking
- `webhook receiver`: fire-and-forget dispatch, atomic retry
- `sync engine`: bounded squashed outbox, batched push

If you suspect CPU burn in production, the framework's own primitives
are unlikely to be the cause — check application code first
(unbounded goroutines, missing Preloads, sync work in handlers).

---

## Critical Rules

1. **Never delete marker comments** (`// grit:*`)
2. **Follow the response format** — `{ data, message }` / `{ data, meta }` / `{ error: { code, message } }`
3. **Always handle errors in Go** — never ignore with `_`
4. **Keep the folder structure** — don't move files
5. **Use React Query** for all data fetching — no raw fetch
6. **Use Zod** for validation — shared between frontend and backend
7. **Use Tailwind + shadcn/ui** — no custom CSS files
8. **Use App Router** — never Pages Router
9. **Report bugs** — if you encounter a bug, open an issue at https://github.com/MUKE-coder/grit/issues with the error message and steps to reproduce
