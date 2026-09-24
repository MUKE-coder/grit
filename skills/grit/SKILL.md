---
name: grit
description: Build full-stack applications with the Grit framework (Go + Gin + GORM backend, Next.js or TanStack Router frontend, generated admin panel). Use this whenever the user asks to start, extend, debug or deploy a Grit project, mentions `grit new`, `grit generate resource`, or has a repo containing grit.config.ts or a .grit/ directory.
---

# Building with Grit

Grit is a full-stack meta-framework: a Go API (Gin + GORM), a Next.js or
TanStack Router frontend, and a generated admin panel, in one monorepo with
shared Zod schemas and TypeScript types.

**The single most important thing to understand:** Grit is a *code generator*,
not a runtime library. `grit generate resource Post` writes a Go model, a
service, a handler, a Zod schema, TypeScript types, React Query hooks and an
admin page, and injects the wiring into the router and the registries. That
generated code is yours: read it, edit it, delete it. An agent that hand-writes
those nine files instead of running the command produces something that looks
right, compiles, and is missing the injections that make it reachable.

So the loop is: **generate, then edit.** Never hand-write what the generator
owns.

---

## Step 0: install and verify

```bash
curl -fsSL https://gritframework.dev/install.sh | sh   # macOS, Linux, Git Bash
grit version
```

On Windows PowerShell: `irm https://gritframework.dev/install.ps1 | iex`

With a Go toolchain: `go install github.com/MUKE-coder/grit/v3/cmd/grit@latest`

Already installed? `grit update` self-updates the CLI. `grit upgrade`, run
inside a project, brings that project's templates up to the CLI's version, and
is a different command. Do not confuse them.

Prerequisites: Go 1.21+, Node 20+, pnpm, and Docker if you want Postgres, Redis,
MinIO and Mailhog locally. Full list: https://gritframework.dev/docs/prerequisites

---

## Step 1: choose the architecture before you scaffold

This is the one decision that is annoying to reverse. Ask the user if it is not
obvious from what they described.

| They want | Command | You get |
|---|---|---|
| A product with a public site *and* an admin back office | `grit new app --triple --next` | apps/web + apps/admin + apps/api |
| A product with no separate admin | `grit new app --double --next` | apps/web + apps/api |
| One deployable binary, SPA embedded | `grit new app --single` | Go binary with `go:embed` frontend |
| A backend for someone else's frontend | `grit new app --api` | Go API only |
| A phone app | `grit new app --mobile` | apps/api + apps/expo (React Native) |
| A native desktop app | `grit new-desktop app` | Wails + Go + React + SQLite |

Frontend flag: `--next` (Next.js App Router, default) or `--vite` (TanStack
Router, SPA, faster builds). Database: `--db postgres|mysql|sqlite|memory`.

`grit new .` scaffolds into the current directory and takes the project name
from the folder; add `--force` if it is not empty.

Architecture guide: https://gritframework.dev/docs/concepts/architecture-modes

Then:

```bash
cd app
docker compose up -d      # Postgres, Redis, MinIO, Mailhog
pnpm install
grit start                # or: pnpm dev
```

Register at http://localhost:3000, and the admin is on :3001.

---

## Step 2: model the domain with `grit generate resource`

One command per noun in the user's domain. Do this *before* writing any feature
code, because everything else hangs off these.

```bash
grit generate resource Post --fields "title:string,body:richtext,published:bool,slug:slug:title"
grit generate resource Comment --fields "body:text,post:belongs_to:Post,author:belongs_to:User"
```

### Field types

| Type | Go | TypeScript | Notes |
|---|---|---|---|
| `string` | `string` | `string` | required by default |
| `text` | `string` | `string` | GORM `type:text` |
| `richtext` | `string` | `string` | Tiptap editor in the admin |
| `int` `uint` `float` | `int` `uint` `float64` | `number` | `float` is for weights and ratings, never money |
| `bool` / `toggle` | `bool` | `boolean` | `toggle` is the friendlier alias |
| `date` `datetime` `time` | `*time.Time` | `string \| null` | |
| `money` | minor units + currency | | **use this for money.** 0.1 + 0.2 is not 0.3 |
| `email` `url` `domain` `tel` `country` `color` `percent` `rating` | `string` / number | | checked and normalised on every write; the admin renders the matching input. `tel:UG` sets a default country, `rating:10` the number of stars |
| `slug` | `string` | `string` | auto-unique, from another field |
| `select` `radio` | `string` | union | one value from a fixed list |
| `check` | JSON array | `string[]` | zero or more from a list |
| `file` / `files` | `FileRef` / `[]FileRef` | | upload widget, S3-backed; `file:image` limits what is accepted |
| `json` | JSON column | | |
| `string_array` | `datatypes.JSONSlice[string]` | `string[]` | |
| `belongs_to` | `string` | `string` | a UUID foreign key + index. **Not `uint`** |
| `one_to_one` | `string` | `string` | a `belongs_to` whose key is unique. Declared on the side holding the key |
| `many_to_many` | `[]string` | `string[]` | junction table + a picker in the admin |

**Modifiers:** `:unique`, `:optional`, `:encrypted` (AES-256-GCM at rest, on
string/text/richtext), `:slug:<source>`, `:belongs_to:<Model>`,
`:many_to_many:<Model>`, `:select:draft=Draft|paid=Paid`, `:file:image`

```bash
grit generate resource Order --fields "number:string:unique,total:money,status:select:draft=Draft|paid=Paid,customer:belongs_to:Customer,notes:text:optional"
```

Full reference: https://gritframework.dev/docs/concepts/field-types

### What that one command produced

Nine files and several injections. Know them, because these are the files you
will be editing for the rest of the project:

```
apps/api/internal/models/post.go        the GORM model + hooks
apps/api/internal/services/post.go      business logic lives HERE
apps/api/internal/handlers/post.go      thin: parse, call service, respond
apps/api/internal/routes/routes.go      <- injected, not created
packages/shared/src/schemas/post.ts     Zod, shared by every frontend
packages/shared/src/types/post.ts       TypeScript types
apps/web|admin/hooks/use-posts.ts       React Query hooks
apps/admin/app/(dashboard)/resources/posts/page.tsx
apps/admin/lib/resources/post.ts        the admin resource definition
```

https://gritframework.dev/docs/concepts/generated-files

### Other generators

```bash
grit generate resource Product -i                 # interactive
grit generate resource Post --from post.yaml      # from a spec file
grit generate resource Post --seed --faker --count 500   # and a seeder with fake rows
grit generate field Invoice status:select:draft=Draft|paid=Paid   # one column, in place
grit generate seeder Post Comment                 # seeders for existing resources
grit generate form Post                           # a multi-step form for a resource
grit generate table Post                          # a standalone table view
grit remove resource Post                         # deletes files AND reverses every injection
grit add role EDITOR                              # a new role, wired into the guards
grit add variants --resource Product              # product options and variants
grit sync                                         # Go types -> TypeScript + Zod
```

`grit generate` is also `grit g`. `grit generate field` handles scalar, select
and toggle columns in place and lets GORM add the column on the next
`grit migrate`; for a relationship, file, slug or array field, regenerate the
resource instead.

`grit remove resource` reverses injections properly. Deleting files by hand
leaves the router referencing a handler that no longer exists.

---

## Step 3: write the feature, in the right layer

```
handler  ->  service  ->  model
```

- **Handlers are thin.** Parse the request, call one service function, respond.
  No business logic, no direct `db.Where(...)` chains.
- **Services hold the logic.** They take a `*gorm.DB` and a context, return
  values and errors. They are what you unit test.
- **Models hold the shape and the hooks.** `BeforeCreate`, `BeforeUpdate`.
  Never import `services` from `models`: that is an import cycle, and the fix
  is to call the underlying helper directly.

Guides: [handlers](https://gritframework.dev/docs/backend/handlers) ·
[services](https://gritframework.dev/docs/backend/services) ·
[models](https://gritframework.dev/docs/backend/models) ·
[request lifecycle](https://gritframework.dev/docs/backend/request-lifecycle)

### The response format is not optional

Every endpoint answers in one of three shapes. The frontend, the generated
hooks and the error handling all assume it.

```jsonc
// one item
{ "data": { }, "message": "Post created successfully" }

// a list
{ "data": [ ], "meta": { "total": 100, "page": 1, "page_size": 20, "pages": 5 } }

// an error
{ "error": { "code": "VALIDATION_ERROR", "message": "Email is required",
             "details": { "email": "This field is required" } } }
```

Use the `respond` package rather than writing `c.JSON` error bodies by hand:
`respond.Fail(c, respond.CodeValidationError, "...")` picks the status code from
the code. https://gritframework.dev/docs/backend/errors

### Frontend rules

- Data fetching goes through the generated React Query hooks. No `fetch` in a
  component.
- Validate with the Zod schema from `@repo/shared/schemas`. It is the same
  schema the API validates against, which is the point.
- Next.js App Router only. Never Pages Router.
- Tailwind + the existing `components/ui/` primitives. No new CSS files.
- Never `any`.

---

## Step 4: reach for what is already there

The most common failure mode is an agent building something Grit ships. Check
before you write:

| Need | Already in the box |
|---|---|
| Login, register, refresh, roles | [authentication](https://gritframework.dev/docs/backend/authentication) |
| 2FA, passkeys, sessions, sign-in links | [account security](https://gritframework.dev/docs/backend/account-security) · [passkeys](https://gritframework.dev/docs/backend/passkeys) |
| Google / GitHub login | [oauth](https://gritframework.dev/docs/backend/oauth) |
| Permissions beyond roles | [rbac](https://gritframework.dev/docs/backend/rbac) · [authorization](https://gritframework.dev/docs/security/authorization) |
| File uploads, images, video | [storage](https://gritframework.dev/docs/batteries/storage) · [media](https://gritframework.dev/docs/batteries/media) |
| Transactional email | [email](https://gritframework.dev/docs/batteries/email) |
| Background work, retries | [jobs](https://gritframework.dev/docs/batteries/jobs) · [cron](https://gritframework.dev/docs/batteries/cron) |
| Caching | [caching](https://gritframework.dev/docs/batteries/caching) |
| LLM calls | [ai](https://gritframework.dev/docs/batteries/ai) |
| WebSockets / live updates | [realtime](https://gritframework.dev/docs/backend/realtime) |
| Receiving webhooks | [webhooks](https://gritframework.dev/docs/backend/webhooks) |
| Outgoing events, exactly once | [outbox](https://gritframework.dev/docs/backend/outbox) |
| Audit trail | [append-only](https://gritframework.dev/docs/backend/append-only) |
| Feature flags | [feature flags](https://gritframework.dev/docs/backend/feature-flags) |
| Money, tax, invoices | [money](https://gritframework.dev/docs/concepts/money) · [invoices](https://gritframework.dev/docs/backend/invoices) |
| Offline-capable clients | [offline sync](https://gritframework.dev/docs/concepts/offline-sync) |
| Multi-tenancy, impersonation, ⌘K, saved views | `grit plugin add <name>`, see [plugins](https://gritframework.dev/docs/plugins/overview) |
| Backups, GDPR export/erasure | [backups](https://gritframework.dev/docs/batteries/backups) · [compliance](https://gritframework.dev/docs/security/compliance) |

Admin panel: [resources](https://gritframework.dev/docs/admin/resources) ·
[data tables](https://gritframework.dev/docs/admin/datatable) ·
[forms](https://gritframework.dev/docs/admin/forms) ·
[relationships](https://gritframework.dev/docs/admin/relationships) ·
[widgets](https://gritframework.dev/docs/admin/widgets) ·
[custom pages](https://gritframework.dev/docs/admin/custom-pages)

---

## Step 5: verify, every time

Never report work as finished on the strength of having written it. Run:

```bash
cd apps/api
gofmt -l ./internal ./cmd     # must print nothing
go build ./... && go vet ./...
go test ./... -count=1

cd ../admin                   # and ../web
npx tsc --noEmit
pnpm lint
npx next build                # or: npx vite build

cd ../.. && grit doctor       # the mistakes that do not announce themselves
```

A generated project ships its own tests. Add to them rather than replacing
them. https://gritframework.dev/docs/testing

If something looks wrong at runtime rather than at build time:
`grit routes` lists what is actually mounted, `/pulse/ui` traces requests and
queries, `/studio` browses the database, `/docs` is the generated API reference.

---

## Conventions, in one place

| Thing | Convention |
|---|---|
| Go files | `snake_case.go` |
| Go exported / unexported | `GetUsers` / `parseToken` |
| TypeScript files | `kebab-case.ts` |
| React components | `PascalCase.tsx` |
| API routes | plural, lowercase: `/api/v1/posts` |
| Tables | plural, snake_case |
| Zod schemas | `PostSchema`, `CreatePostSchema` |
| Errors | `fmt.Errorf("context: %w", err)`, never swallowed with `_` |

https://gritframework.dev/docs/concepts/naming-conventions

---

## Mistakes that cost the most time

1. **Hand-writing a resource** instead of running the generator. The files are
   right and nothing is wired up.
2. **Deleting resource files by hand.** Use `grit remove resource`.
3. **Business logic in the handler.** It cannot be tested or reused, and the
   next generated handler will not match it.
4. **Building an in-the-box feature.** See step 4 first.
5. **`float` for money.** Use the money type.
6. **Importing `services` from `models`.** Import cycle.
7. **Editing generated wiring by hand** (`routes.go` injections, resource
   registries) when re-running the generator would do it correctly.
8. **Assuming an upgrade rewrote a file you own.** `grit upgrade` deliberately
   leaves project-owned files (`routes.go`, your models) alone and reports what
   it could not do. Read its output.
9. **Reporting success without building.** Step 5.

---

## Worked examples

Each of these is a complete application, start to finish:

- [Build your first Grit app](https://gritframework.dev/blog/build-your-first-grit-app)
- [Build a CRM](https://gritframework.dev/blog/build-a-crm-with-grit)
- [Build a storefront](https://gritframework.dev/blog/build-a-storefront-with-grit)
- [Build an invoice app](https://gritframework.dev/blog/build-an-invoice-app-with-grit)
- [Build a mobile app](https://gritframework.dev/blog/build-mobile-app-with-grit)
- [Build a desktop app](https://gritframework.dev/blog/build-desktop-app-with-grit)
- [Your table, our machinery](https://gritframework.dev/blog/your-table-our-machinery): how the generator thinks
- [Roles, permissions and automatic backups](https://gritframework.dev/blog/roles-permissions-and-automatic-backups)
- [Plugins: what they are and how to build one](https://gritframework.dev/blog/grit-plugins-what-they-are-and-how-to-build-one)
- [GDPR and access reviews](https://gritframework.dev/blog/gdpr-and-access-reviews)
- [Why I built Grit](https://gritframework.dev/blog/why-i-built-grit)

Tutorials: [blog](https://gritframework.dev/docs/tutorials/blog) ·
[contact app](https://gritframework.dev/docs/tutorials/contact-app) ·
[e-commerce](https://gritframework.dev/docs/tutorials/ecommerce) ·
[SaaS](https://gritframework.dev/docs/tutorials/saas) ·
[custom endpoints](https://gritframework.dev/docs/tutorials/custom-endpoints)

---

## Deploying

```bash
grit deploy --host user@server.com --domain myapp.com
```

SSH + systemd + Caddy with automatic TLS. Also documented for Railway, Render,
Fly.io, Coolify, Dokploy and plain VPS:
https://gritframework.dev/docs/deployment

Pre-flight list: https://gritframework.dev/docs/deployment/checklist

---

## Where to look when stuck

- Full docs: https://gritframework.dev/docs
- CLI reference: https://gritframework.dev/docs/cli
- Start here: https://gritframework.dev/docs/start
- Changelog: https://gritframework.dev/docs/changelog
- What is stable and what is not: https://gritframework.dev/docs/stability
- Machine-readable guide for LLMs: https://gritframework.dev/docs/ai-skill/llm-guide
- MCP server: https://gritframework.dev/docs/ai-workflows/mcp
