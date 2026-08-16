# Grit — Project Phases

This document breaks the Grit framework development into 5 phases. Each phase builds on the previous one. **AI agents (Claude Code) should complete one phase fully before moving to the next.**

---

## Phase 1 — Foundation (Weeks 1-3)

**Goal:** A working monorepo that scaffolds via CLI, with a Go API (auth included), Next.js frontend with login/register, shared schemas, Docker setup, and GORM Studio.

**Success Criteria:** Run `grit new myapp`, then `cd myapp && docker compose up`, open the browser, register a user, log in, see a dashboard, and browse the database via GORM Studio.

### 1.1 CLI Scaffolder

- [x] Create the `grit` CLI tool in Go using `cobra` library
- [x] `grit new <project-name>` command that:
  - Creates the full monorepo folder structure (see GRIT.md for structure)
  - Initializes `go.mod` for the API
  - Initializes `package.json` files for web, admin, and shared packages
  - Creates `pnpm-workspace.yaml`
  - Creates `turbo.json`
  - Creates `.env` and `.env.example` with sensible defaults
  - Creates `docker-compose.yml` (PostgreSQL, Redis, MinIO, Mailhog)
  - Creates `.gitignore` files
  - Runs `go mod tidy` and `pnpm install`
  - Prints success message with next steps
- [x] `grit new <project-name> --api` flag (scaffolds only the Go API, no frontend)
- [x] Colored CLI output with ASCII art logo
- [x] Validate project name (lowercase, alphanumeric, hyphens only)

### 1.2 Go API Boilerplate

- [x] Entry point: `apps/api/cmd/server/main.go`
- [x] Configuration loading from `.env` using `godotenv`
- [x] Database connection setup with GORM + PostgreSQL driver
- [x] Auto-migration on startup
- [x] Gin router setup with middleware:
  - CORS middleware (configured for frontend origins)
  - JSON logging middleware
  - Recovery middleware
- [x] Route registration in `internal/routes/routes.go`
- [x] Health check endpoint: `GET /api/health`
- [x] Graceful shutdown handling

### 1.3 Authentication System

- [x] **User model** (`internal/models/user.go`):
  - ID, Name, Email, Password (hashed), Role, Avatar, Active, CreatedAt, UpdatedAt, DeletedAt
  - Password hashing with bcrypt on BeforeCreate hook
  - `CheckPassword()` method
- [x] **Role constants**: admin, editor, user (as string constants on User model)
- [x] **Auth handlers** (`internal/handlers/auth.go`):
  - `POST /api/auth/register` — Register with name, email, password
  - `POST /api/auth/login` — Login with email + password, returns JWT tokens
  - `POST /api/auth/refresh` — Refresh access token using refresh token
  - `POST /api/auth/logout` — Invalidate refresh token
  - `GET /api/auth/me` — Get current user profile
  - `POST /api/auth/forgot-password` — Send password reset email
  - `POST /api/auth/reset-password` — Reset password with token
- [x] **Auth middleware** (`internal/middleware/auth.go`):
  - JWT validation middleware
  - Role-based access middleware (`RequireRole("admin")`)
  - Extracts user from token and attaches to Gin context
- [x] **JWT service** (`internal/services/auth.go`):
  - Generate access token (15 min expiry)
  - Generate refresh token (7 day expiry)
  - Validate and parse tokens
  - Token secret from environment variable
- [x] **User handlers** (`internal/handlers/user.go`):
  - `GET /api/users` — List users (admin only, paginated)
  - `GET /api/users/:id` — Get user by ID
  - `PUT /api/users/:id` — Update user
  - `DELETE /api/users/:id` — Soft delete user

### 1.4 GORM Studio Integration

- [x] Embed GORM Studio in the API
- [x] Mount at `/studio` route
- [x] Pass all registered models to GORM Studio
- [x] Enable by default in development, configurable in `.env`

### 1.5 Next.js Frontend (Web App)

- [x] Next.js 14+ with App Router
- [x] Tailwind CSS + shadcn/ui setup
- [x] Dark theme as default (matching GORM Studio aesthetic)
- [x] API client (`lib/api-client.ts`):
  - Axios instance pointed at Go API
  - Automatic JWT token injection from cookies/localStorage
  - Token refresh interceptor
  - Error handling wrapper
- [x] React Query setup (`lib/query-client.ts`)
- [x] Auth hooks (`hooks/use-auth.ts`):
  - `useLogin()`, `useRegister()`, `useLogout()`, `useMe()`
  - Store tokens securely
  - Redirect on 401
- [x] Auth pages:
  - `/login` — Email + password form, dark themed
  - `/register` — Name + email + password form
  - `/forgot-password` — Email input
- [x] Protected layout (`(dashboard)/layout.tsx`):
  - Auth guard — redirect to login if not authenticated
  - Sidebar navigation
  - Top navbar with user avatar + logout
- [x] Dashboard page (`(dashboard)/dashboard/page.tsx`):
  - Welcome message with user name
  - Placeholder stats cards
  - Placeholder recent activity
- [x] Responsive design (mobile-friendly)

### 1.6 Shared Package

- [x] `packages/shared/schemas/user.ts` — Zod schema for user registration/login
- [x] `packages/shared/types/user.ts` — TypeScript User type
- [x] `packages/shared/types/api.ts` — API response types (pagination, error format)
- [x] `packages/shared/constants/index.ts` — Roles, API routes, config

### 1.7 Docker Setup

- [x] `docker-compose.yml` for development:
  - PostgreSQL 16 with persistent volume
  - Redis 7
  - MinIO (S3-compatible, for local file storage testing)
  - Mailhog (local email testing)
- [x] `docker-compose.prod.yml` for production:
  - Multi-stage Go build (small final image)
  - Next.js standalone build
  - PostgreSQL + Redis
- [x] `Dockerfile` for Go API (multi-stage build)
- [x] `Dockerfile` for Next.js web app
- [x] `.dockerignore` files

### 1.8 Developer Experience

- [x] Sensible `.env.example` with all required variables documented
- [x] README.md with quick start instructions

### Phase 1 Deliverables

- Working `grit` CLI that scaffolds the full project
- Go API with JWT auth and role-based access
- Next.js app with login, register, and protected dashboard
- GORM Studio at `/studio`
- Docker Compose for dev and prod
- Shared Zod schemas and TypeScript types
- Hot reload development setup

---

## Phase 2 — Code Generator (Weeks 4-6)

**Goal:** `grit generate resource <Name>` creates a full-stack resource — Go model, handler, React Query hook, Zod schema, and admin panel page — all wired together.

**Success Criteria:** Run `grit generate resource Post`, then immediately browse Posts in the admin panel with full CRUD, pagination, sorting, and filtering.

### 2.1 Code Generator Engine

- [x] Template engine for Go and TypeScript file generation (strings.NewReplacer with named placeholders)
- [x] `grit generate resource <Name>` command:
  - Prompts for fields (name, type, required, unique) interactively (`-i`)
  - Or accepts a definition file: `grit generate resource --from post.yaml`
  - Or inline fields: `--fields "title:string,content:text,published:bool"`
- [x] Template functions for each generated artifact
- [x] Smart pluralization (Post → posts, Category → categories)
- [x] Marker-based code injection into existing files
- [x] Automatic route registration (inject into routes.go)

### 2.2 Generated Go Artifacts

- [x] **Model** (`internal/models/<name>.go`):
  - GORM struct with proper tags
  - Timestamps and soft deletes
- [x] **Handler** (`internal/handlers/<name>.go`):
  - `GET /api/<names>` — List with pagination, sorting, filtering, search
  - `GET /api/<names>/:id` — Get by ID
  - `POST /api/<names>` — Create with validation
  - `PUT /api/<names>/:id` — Update with validation
  - `DELETE /api/<names>/:id` — Soft delete
- [x] **Service** (`internal/services/<name>.go`):
  - Business logic layer between handler and model
  - Reusable query scopes (pagination, filtering, sorting)
- [x] Automatic migration registration (AutoMigrate + GORM Studio injection)

### 2.3 Generated Frontend Artifacts

- [x] **Zod schema** (`packages/shared/schemas/<name>.ts`):
  - Create schema, update schema
  - Proper Zod types matching Go types
- [x] **TypeScript types** (`packages/shared/types/<name>.ts`):
  - Full type with all fields
  - Create/Update input types inferred from Zod
- [x] **React Query hooks** (`apps/admin/hooks/use-<names>.ts` + `apps/web/hooks/`):
  - `use<Names>()` — Paginated list query with sorting/filtering
  - `useGet<Name>(id)` — Single item query
  - `useCreate<Name>()` — Create mutation
  - `useUpdate<Name>()` — Update mutation
  - `useDelete<Name>()` — Delete mutation
  - Automatic cache invalidation on mutations
- [x] **Admin page** (`apps/admin/app/resources/<names>/page.tsx`):
  - DataTable with search, pagination, delete actions
  - Auto-injected into admin sidebar navigation

### 2.4 Type Sync Command

- [x] `grit sync` command:
  - Reads all Go models in `internal/models/` using Go AST parser
  - Generates corresponding TypeScript types and Zod schemas
  - Maps Go types → TypeScript types (uint → number, time.Time → string, etc.)
  - Skips User model (has custom schemas)

### Phase 2 Deliverables

- Working `grit generate resource` command
- Full-stack resource generation (Go + TypeScript)
- `grit sync` for Go → TypeScript type generation
- Generated code is clean, readable, and editable

---

## Phase 3 — Admin Panel (Weeks 7-12)

**Goal:** A Filament-like admin panel where developers define resources in TypeScript and get beautiful, functional admin pages with data tables, forms, dashboards, and widgets.

**Success Criteria:** Define a resource with 10+ columns, 8+ form fields, 3+ filters, and 2+ widgets. The resulting admin page should look like a premium CRM — not a generic CRUD tool.

### 3.1 Admin Layout Shell

- [x] Collapsible sidebar with:
  - Logo/brand area
  - Navigation items auto-generated from registered resources
  - Icon support (Lucide icons)
  - Active state highlighting
  - User profile section at bottom
  - Dark/light theme toggle
- [x] Top navbar:
  - Breadcrumbs
  - Search (global search across resources)
  - User menu (profile, settings, logout)
- [x] Responsive layout (sidebar collapses on mobile)

### 3.2 Resource System

- [x] `defineResource()` API (see GRIT.md for example)
- [x] Resource registry (`resources/index.ts`)
- [x] Auto-generated routes from resources
- [x] Resource configuration:
  - Name, endpoint, icon
  - Table columns, filters, actions
  - Form fields, validation
  - Dashboard widgets
  - Permissions (which roles can access)

### 3.3 DataTable Component

- [x] Server-side pagination (communicates with Go API)
- [x] Column sorting (click header to sort)
- [x] Column filtering:
  - Text search (global and per-column)
  - Select/dropdown filters
  - Boolean toggle filters
- [x] Column features:
  - Show/hide columns toggle
  - Custom cell renderers (badge, currency, date, relative time, image, boolean)
- [x] Row actions (edit, delete, view, custom)
- [x] Empty state with illustration
- [x] Loading skeleton
- [x] Export to CSV / JSON
- [x] Responsive (horizontal scroll on mobile)

### 3.4 Form Builder

- [x] Form modal and full-page form views
- [x] Field types:
  - Text input
  - Textarea
  - Number (with min/max, step)
  - Select
  - Date picker
  - Toggle / Switch
  - Checkbox
  - Radio group
- [x] Validation:
  - Zod-based validation from shared schemas
  - Real-time validation (on blur and on change)
  - Server-side error display
- [x] Form layout:
  - Single column, two column
  - Section groups with headers
- [x] Create and edit modes from the same form definition
- [x] Auto-populated defaults

### 3.5 Dashboard & Widgets

- [x] Dashboard page as the admin home
- [x] Widget types:
  - **Stats card** — Number + label + change percentage + icon
  - **Line chart** — Time series data
  - **Bar chart** — Categorical data
  - **Recent activity** — List of recent events
- [x] Widget grid layout (responsive, configurable)
- [x] Widgets fetch data from the Go API
- [x] Charting library: Recharts

### 3.6 Admin Theme

- [x] Dark theme (default) matching GORM Studio:
  - Background: `#0a0a0f`, `#111118`, `#1a1a24`
  - Accent: `#6c5ce7`
  - Fonts: Onest + JetBrains Mono
- [x] Light theme option
- [x] Theme toggle in sidebar
- [x] CRM-inspired aesthetic:
  - Generous spacing
  - Professional data density
  - Beautiful empty states
  - Polished loading states

### Phase 3 Deliverables

- Complete admin panel layout shell
- Resource definition system
- DataTable with server-side pagination, sorting, filtering
- Form builder with all field types
- Dashboard with widgets (stats, charts, activity)
- Dark/light theme
- Beautiful, CRM-quality aesthetic

---

## Phase 4 — Batteries (Weeks 13-16)

**Goal:** Add file storage, email, background jobs, cron, Redis caching, and AI integration. All pre-configured and wired to the admin panel.

### 4.1 File Storage

- [x] Storage abstraction layer (`internal/storage/storage.go`):
  - Interface: `Upload()`, `Download()`, `Delete()`, `GetURL()`, `GetSignedURL()`
  - S3 driver (AWS S3, Cloudflare R2, Backblaze B2)
  - Local driver (MinIO in dev)
  - Configuration via `.env` (STORAGE_DRIVER, S3_BUCKET, S3_REGION, etc.)
- [x] Upload handler: `POST /api/uploads`
  - File size limits, allowed MIME types (configurable)
  - Returns file URL and metadata
- [x] Image processing on upload:
  - Thumbnail generation (via background job)
  - Resize to max dimensions
- [x] File upload in admin panel:
  - Drag and drop
  - Image previews
  - Multiple file support
- [x] File management in admin panel:
  - Browse uploaded files (grid view)
  - Delete files
  - View file details

### 4.2 Email System

- [x] Mail service (`internal/mail/mailer.go`):
  - Resend client integration (raw net/http, no SDK)
  - Send method: `Send(ctx, SendOptions)`
  - HTML templates using Go `html/template`
- [x] Built-in email templates:
  - Welcome email
  - Password reset
  - Email verification
  - Notification
- [x] Template preview in admin panel
- [x] Email configuration via `.env`

### 4.3 Background Jobs

- [x] Job queue system (`internal/jobs/`):
  - Redis-backed queue using `asynq` library
  - Enqueue: `client.EnqueueSendEmail()`, `client.EnqueueProcessImage()`
  - Job priorities (critical, default, low)
  - Retry with configurable max attempts
- [x] Built-in jobs:
  - Send email job
  - Process image job (thumbnail generation)
  - Cleanup expired tokens job
- [x] Job dashboard in admin panel:
  - Queue stats (pending, active, completed, failed)
  - View jobs by status with error details
  - Retry failed jobs
  - Clear queues

### 4.4 Cron Scheduler

- [x] Cron service (`internal/cron/cron.go`):
  - asynq Scheduler for cron-like scheduling
  - Built-in tasks: cleanup expired tokens (hourly)
  - `// grit:cron-tasks` marker for injection
- [x] Cron status in admin panel (list registered tasks)

### 4.5 Redis Caching

- [x] Cache service (`internal/cache/cache.go`):
  - `Get(key)`, `Set(key, value, ttl)`, `Delete(key)`, `DeletePattern()`, `Flush()`
  - JSON serialization for complex values
  - Cache middleware for API GET responses (`CacheResponse()`)

### 4.6 AI Integration

- [x] AI service (`internal/ai/ai.go`):
  - Anthropic Claude API client (raw net/http)
  - OpenAI API client (raw net/http)
  - Simple interface: `ai.Complete(ctx, req)`, `ai.Stream(ctx, req, handler)`
- [x] Streaming support (SSE via Gin)
- [x] API endpoints: `/api/ai/complete`, `/api/ai/chat`, `/api/ai/stream`
- [x] Configuration via `.env` (AI_PROVIDER, AI_API_KEY, AI_MODEL)

### Phase 4 Deliverables

- File storage with S3/R2/MinIO
- Email system with Resend + templates
- Background job queue with Redis
- Cron scheduler
- Redis caching
- AI integration (Anthropic + OpenAI)
- All features visible/manageable in admin panel

---

## Phase 5 — Polish & Launch (Weeks 17-20)

**Goal:** Documentation site, testing, performance optimization, and public launch.

### 5.1 Documentation Site

> **Shipped.** These boxes were stale — the site has been live at
> gritframework.dev for a long time and covers considerably more than the list
> below (desktop, mobile, plugins, governance, security, tech kits, tutorials,
> a stack selector and a changelog, none of which were planned here).

- [x] Documentation website — custom Next.js app in `docs/`
- [x] Pages:
  - Getting Started (5-minute quick start)
  - Installation
  - Project Structure
  - Configuration
  - Authentication & Authorization
  - Admin Panel (resource definitions, tables, forms, widgets)
  - Code Generator (CLI commands, templates)
  - File Storage
  - Email
  - Background Jobs
  - Cron
  - Caching
  - AI Integration
  - GORM Studio
  - Deployment (Docker, VPS, cloud)
  - API Reference
  - Contributing Guide
  - FAQ
- [x] Interactive examples (live showcases: auth, theme, resource definitions, dropzone)
- [x] Copy-paste code blocks (`components/code-block.tsx`)
- [x] Search functionality (`components/search-dialog.tsx`)
- [x] Dark theme docs site (`components/theme-provider.tsx`)

### 5.2 Testing

- [x] Go CLI/generator tests:
  - Unit tests for `internal/generate` package (pluralize, field types, definition parsing)
  - 31 test cases covering GoType, TSType, ZodType, GORMTag, field helpers, FKColumnName, RelatedModelName, UIHelpers, ParseInlineFields, LoadFromYAML, ValidFieldTypes
- [x] Go CLI/scaffold tests:
  - 13 tests for `internal/scaffold` package — ValidateProjectName, ValidateStyle, ShouldInclude* helpers, createDirectories (api-only/default/full modes), writeAPIFiles (module substitution, key files exist), writeFile helper
  - Fixed pre-existing build errors: `%3C` URL-encoding in `fmt.Sprintf` template, trailing `\n` in `color.Println` calls
- [x] Go CLI/generator inject + sync tests:
  - `inject_test.go`: injectBefore (inserts before marker, marker-not-found, idempotent, missing file), injectInline (inserts inline, errors), guessLucideIcon (12 name→icon cases)
  - `sync_test.go`: goTypeToTS (15 type mappings), goTypeToZod (10 type+tag combos), extractTag (7 cases incl. omitempty), isAutoField, buildTSType, buildZodSchema (auto fields excluded, .optional() in update schema), parseGoStructs (simple struct, multiple structs, invalid Go file, missing file), round-trip test
  - 68 total test cases (including subtests) — all passing
- [x] Go CLI/generator integration tests:
  - `generator_test.go`: Names() (simple, compound, irregular plurals), writeGoModel (basic fields, slug + helpers.go, belongs_to), writeGoService (module substitution, search where), writeGoHandler (CRUD methods, no placeholders), writeZodSchema (Create/Update/type exports), writeTSTypes (interface fields)
  - Full `Generator.Run()` integration: basic resource (all files created + all 8 injections verified), idempotent (markers preserved), compound names (BlogPost → blog_post.go), role-restricted routes
  - `remove_test.go`: removeLinesContaining (removes matching, error on missing), removeInlineText, removeLineBlock (handler init blocks), removeSchemaExportBlock (multi-line export blocks), generate+remove round-trip (user.go and types/index.ts restored to original)
  - **56 total passing tests** across both packages — 100% `go vet` clean
- [x] Go API tests:
  - Integration tests for handlers
  - Auth flow tests
  - Database tests (SQLite in-memory)
- [x] Frontend tests:
  - Component tests (React Testing Library + Vitest) scaffolded into web + admin apps
  - Utility unit tests (cn, formatCurrency, truncate)
  - E2E tests (Playwright) for auth flow and admin panel navigation
- [x] CI/CD:
  - GitHub Actions `ci.yml` — test + race detector + coverage + cross-platform build (linux/darwin/windows × amd64/arm64)
  - GitHub Actions `release.yml` — tag-triggered release with binaries + auto release notes

### 5.3 Performance

- [x] Go API:
  - Connection pooling: `SetMaxIdleConns(10)`, `SetMaxOpenConns(100)`, `SetConnMaxLifetime(30m)` in scaffolded `database.go`
  - Response compression: `Gzip()` middleware (gzip.BestSpeed + Vary header) in all scaffolded APIs
  - Request ID tracing: `RequestID()` middleware for log correlation across requests
  - Rate limiting: Sentinel — per-IP (100 req/min) + per-route limits on auth endpoints
  - Benchmarks: 7 `BenchmarkXxx` functions in `internal/generate/bench_test.go` — run with `go test -bench=.`
- [x] Frontend:
  - Image optimization (Next.js `<Image>` — already in web app template)
  - Bundle analysis and code splitting (admin lazy routes + @next/bundle-analyzer)
- [x] Additional benchmarks:
  - API response time benchmarks (httptest-based)
  - BenchmarkHealthCheck, BenchmarkAuthLogin, BenchmarkAuthRegister in scaffolded API

### 5.4 Launch Preparation

- [ ] README.md with GIF demo
- [x] Landing page at `gritframework.dev` (`docs/app/page.tsx`)
- [ ] YouTube tutorial: "Build a SaaS in 10 Minutes with Grit"
- [ ] Blog post: "Why I Built Grit — Go + React Framework with Admin Panels"
- [ ] Product Hunt listing
- [ ] Dev.to article
- [ ] Twitter/X thread
- [ ] LinkedIn post
- [ ] TikTok demo video
- [ ] Reddit posts (r/golang, r/reactjs, r/webdev)
- [ ] Hacker News: Show HN
- [ ] Go Weekly / React Newsletter submission
- [ ] Discord community server

### Phase 5 Deliverables

- Documentation website
- Comprehensive test suite
- Performance optimized
- Public launch across all channels
- Community infrastructure (Discord, GitHub Discussions)

---

## Phase 6 — Feature Completion ✅ complete (v3.120.0 → v3.130.0)

> **Read this first when asked "what's remaining".** Tick a box the moment the
> work is verified and pushed, not when the code is written. Every item here is
> a real gap found by using the framework, not a wish-list.
>
> All eight workstreams are shipped. The only sub-item still open is 6.8's
> end-to-end check of a *desktop* signed release, which needs Authenticode and
> Apple Developer certificates the project does not own — the workflow is
> written and the API half signs on every tag.
>
> **Definition of done for every sub-task below:** `go build ./...` passes,
> `go test ./internal/...` passes, a project is scaffolded fresh, **both admin
> variants build** (`--triple` with `next build`, `--triple --vite` with
> `vite build`), the behaviour is driven over HTTP or in a browser, and the
> release is tagged and pushed. See [Release testing](#) — skipping the
> both-variants check is what broke v3.119.0.

### 6.1 — Audit log surface ✅ shipped v3.120.0
- [x] `/system/audit` page listing every chained entry (method, status, duration, digest)
- [x] Verify-chain button wired to `/admin/activity/integrity`
- [x] Registered for Next.js **and** TanStack variants + System Hub tile
- [x] Weekly `audit:prune` retention job that re-anchors the chain
- [x] Verified by editing a row directly in SQLite — the check named that row

### 6.2 — Email verification ✅ shipped v3.121.0
- [x] `EmailVerificationToken` model (hash-only, 48h, single-use, records the address)
- [x] `POST /auth/verify-email` (public) + `POST /auth/verify-email/send` (authenticated)
- [x] Verification mail on register, sent off the request path
- [x] `REQUIRE_EMAIL_VERIFICATION` login gate, off by default
- [x] `/verify-email` page + resend banner, both admin variants
- [x] 7 generated tests, including the changed-address case

### 6.3 — Account lockout ✅ shipped v3.122.0
- [x] `failed_login_count` / `locked_until` on the user
- [x] Lock after `LOGIN_MAX_ATTEMPTS`, checked before the password comparison
- [x] Unknown emails never count (else it becomes a DoS tool)
- [x] `POST /users/:id/unlock` for admins, written to the activity log
- [x] Verified: 3 fails → 401, 4th → 429, unknown email never locks, unlock resets

### 6.4 — API-key authentication ✅ shipped v3.123.0
- [x] `APIKey` model: hashed secret, prefix for display, name, last_used_at, expires_at, revoked_at
- [x] Issuance endpoint returning the raw key **once** (server keeps only the hash)
- [x] `middleware.APIKeyAuth` that coexists with JWT — same `user_id` in context
- [x] Scope/permission binding, reusing the existing `perm:<key>` catalogue
- [x] Revoke + rotate endpoints
- [x] Admin page: list, create (one-time reveal), revoke
- [x] Registered for both admin variants
- [x] Generated tests: hash-only storage, revoked key rejected, expired key rejected
- [x] Verified over HTTP end to end

### 6.5 — OpenAPI coverage for built-in endpoints ✅ shipped (v3.130.0)
> Measured against a live `/docs/openapi.json` from a scaffolded project, not
> estimated: **134 of 134** operations now carry a response schema, up from 29.
> Request bodies: 117 of 134, and the other 17 are exactly the operations that
> have no JSON body to describe.
- [x] Every operation has a response schema (134/134)
- [x] Request bodies wherever there is one to describe (117/134)
- [x] The 17 without one are correct as they stand: 14 bodyless action POSTs
      (logout, revoke-all, totp/setup, retry, replay, unlock, close, reopen,
      read, read-all, generate, backup-codes, verify-email/send, review
      complete), `POST /uploads` (multipart, not JSON), `POST /webhooks/:provider`
      (an arbitrary third-party payload) and the SAML ACS callback (form-encoded)
- [x] 16 handlers that bound `var req struct{…}` inline now bind named exported
      types — gindocs reflects over a type, and routes.go is a different package,
      so an anonymous or unexported struct gives it nothing to read
- [x] 16 already-named request types exported for the same reason
- [x] Verified by compiling a scaffolded project and counting the live spec

> Fixed while verifying: the generated project's own auth tests failed on
> `grit new`. `newTestDB` migrated only `models.User`, so registering could not
> write its session or activity row, and the duplicate-email case surfaced as a
> 500 instead of a 409. `email_verification_tokens` is deliberately still not
> migrated there — that write happens on a goroutine, and creating the table
> turns a harmless early return into a background writer racing the test on one
> in-memory SQLite connection.
- [x] Generated project: `go test ./...` is green across all ten packages

### 6.6 — 2FA on Expo and desktop ✅ shipped (v3.125.0 challenge, v3.129.0 enrolment)
> The QR arrives as a PNG data URI from the API, so neither client ships a QR
> encoder — `<img>` and React Native `<Image>` both render it directly.
- [x] Expo: TOTP challenge on login (code, backup code, trust device)
- [x] Expo: enrol screen at `app/two-factor.tsx`, reached from Settings → Security
- [x] Desktop: TOTP challenge in the Wails client
- [x] Desktop: enrol section on the profile page, between Password and Delete
- [x] Backup codes on both: shown once, and the panel will not close until they
      have been saved. Mobile offers Share rather than the clipboard — it is the
      affordance phones actually have for getting text into a password manager,
      and it avoids adding expo-clipboard for one screen
- [x] Verified end to end in both running clients with live TOTP codes: setup →
      QR → verify → enabled → backup codes, and the disable path
- [x] Fixed en route: `expo-secure-store` has no web implementation, so every
      scaffolded Expo app threw `getValueWithKeyAsync is not a function` on the
      first render under `pnpm web`. `lib/secure-store.ts` now wraps it —
      unchanged on native, localStorage on web

### 6.7 — Swap Phase 3 ✅ shipped (v3.127.0 button, v3.128.0 input)
> The earlier count of 185 was wrong — it came from a regex that broke on `=>`
> inside `onClick`, and it counted `bg-accent/10` badges as buttons. A
> brace-aware scan finds **142** inline `bg-accent` buttons across all templates,
> **91** with a literal class string, and **46** of those in the admin.
>
> The refactor routes class strings through `buttonClasses()` rather than
> restructuring call sites into `<Button>`. The slot exports it for exactly this
> reason, and it leaves the element, handlers, spinner logic and children alone —
> so the diff is one attribute per site instead of a rewrite.
>
> Sizes are inferred from height, never width: a button that changes height
> shifts the row it sits in, while a few pixels of horizontal padding go
> unnoticed.
- [x] Adopt the slot at the 46 literal-className `bg-accent` sites in the admin
- [x] Slot reach measured, not assumed: **2 → 23** emitted files call `buttonClasses()`
- [x] Decided deliberately what does NOT follow the slot — themed auth pages
      (they use `var(--auth-primary)` so atlas/aurora/pulse can restyle them;
      routing those through the slot would fight the theme system), accent-tint
      badges, and buttons whose className is a runtime expression
- [x] Regression guard: `TestButtonClassesImportsMatchUsage` catches a call with
      no import *and* an import with no call (the wrong-emitted-file failure),
      across both admin variants
- [x] Budget test pins the remaining inline count at 13 so it can only go down
- [x] Verified visually: `/system/roles` "New role" goes `rounded-lg` → pill
      under `grit swap button glow-ring`, with no layout shift
- [x] Both variants rebuilt from a fresh scaffold (`next build`, `vite build`)
- [x] Adopt the swappable Input primitive — 40 sites across 13 files (v3.128.0)
- [x] Decided what does NOT follow the input slot: checkboxes and radios (the
      slot sets `w-full`, wrong for a 16px box), fields carrying an explicit
      width that would fight that same `w-full`, and toolbar chrome like the
      page-size select, which is a control rather than a form field
- [x] Two conflicts resolved by hand rather than mechanically: a redundant
      `text-text-primary` next to the slot's `text-foreground` (both resolve to
      `var(--text-primary)`), and the confirm-to-delete field, which now uses
      the slot's own `invalid` state instead of `focus:border-danger` fighting
      `focus:border-accent`
- [x] Verified visually: every field in the New-form-share modal restyles
      together under `grit swap input soft-filled`
- [x] `grit swap` round-trip verified clean (swap + revert, with type-check)

### 6.8 — Release hardening 🟡 partial (v3.126.0)
> Correction: Grit's OWN release.yml already had SBOM, keyless cosign and
> provenance. The gap was that scaffolded projects had no release workflow at
> all — the framework met a bar its output could not.
- [x] SBOM generation in the generated project's release workflow
- [x] Signed release artifacts (cosign keyless — no secrets required)
- [x] Build provenance attestation
- [x] Authenticode + notarization steps, guarded on secrets, warning when absent
- [x] Verified a real signed release end to end on the v3.130.0 tag: the
      workflow published binaries, `SHA256SUMS`, an SPDX SBOM, a certificate and
      a signature, and `cosign verify-blob` returns **Verified OK** against the
      exact workflow identity and tag:

      cosign verify-blob SHA256SUMS \
        --certificate SHA256SUMS.pem --signature SHA256SUMS.sig \
        --certificate-identity-regexp \
          "https://github.com/MUKE-coder/grit/.github/workflows/release.yml@refs/tags/v3.130.0" \
        --certificate-oidc-issuer "https://token.actions.githubusercontent.com"

- [ ] Desktop installers signed end to end — needs Authenticode and Apple
      Developer certificates this project does not own. The steps are written
      and guarded, so a repo with the secrets gets signed installers and one
      without gets a build warning instead of a failure.

---

## Phase 7 — Beyond Filament parity (in progress)

Two independent reviews landed in August 2026: a gap analysis against Filament
(17 items, `grit-framework-gap-analysis.md`) and a product/engineering review
(31 recommendations, `grit-framework-product-engineering-review.md`). This is
the merged list, ordered by what unblocks what. Where the two reviews agreed
without talking to each other, the item is marked **[both]**.

Both reviews were partly out of date, because they read the public site. Fixing
that is item 7.1 for exactly that reason.

### 7.0 Shipped

- [x] **MySQL support** (gap 3.6) - v3.146.0. DSN branch, `parseTime` default,
      `jsonb` tags removed, RETURNING made runtime-conditional via
      `internal/database/dialect.go`, boolean query filters coerced against the
      model schema. Verified against MySQL 8.4.

### 7.1 Credibility (days)

- [x] Correct the public docs and README (gap 3.11) - v3.148.0. Email
      verification and API keys documented in Authentication, MySQL named in the
      README stack table, and the mobile offline page no longer claims offline
      writes are desktop-only.
- [ ] `grit doctor` **[both]** (new 4.4 / rec 23): Go and Node versions, DB,
      Redis, storage, mail, route and type sync, plugin integrity, weak JWT
      secret, wildcard CORS in production.

### 7.2 The generator engine **[both]** (gap 4.1 / recs 21, 22)

Ranked above the plugin registry, against the gap analysis's own sequencing.
Three bugs in one week came from its absence: `grit upgrade` wrote 31 files
nothing imports, `/bulk` 404'd on every existing project because upgrade never
regenerates API code, and a generated handler referenced a `dialect.go` the
project had never been given. A plugin ecosystem on top of a generator that
cannot upgrade its own output multiplies that by the number of plugins.

- [x] `.grit/manifest.json` provenance: which generator wrote which file, at
      which version, and its content hash (`internal/manifest`).
- [x] `grit upgrade` stops overwriting files you have edited. Before this it
      overwrote everything unconditionally: `writeUpgradeFiles` took a `force`
      parameter and never read it.
- [x] `grit upgrade --diff` shows what the new version would change in the
      files it left alone; `--force` keeps the old behaviour.
- [ ] `grit upgrade --resource Product`: regenerate one resource's files
      through the same three-way logic.
- [ ] Per-resource generated tests (rec 19), so a generator change is provably
      safe before it ships.

### 7.3 Ecosystem **[both]** (gaps 3.1, 3.2, 4.4 / recs 19, 20)

- [ ] `grit plugin add github.com/you/grit-x` from a git ref.
- [ ] A JSON plugin index and `grit plugin search`.
- [ ] `grit plugin info <module>`: compatibility, last-tested Grit version,
      known issues, shown at install time rather than buried in a README.
- [ ] Fix the `uint` vs UUID `user_id` mismatch in the official `grit-plugins`
      packages before any tutorial promotes them.

### 7.4 The application layer (product review, section 4)

The review's real contribution: Grit has webhooks, asynq, realtime and audit as
four separate consumers with no shared bus.

- [ ] **Domain events** (rec 7). The substrate everything else hangs off.
- [ ] **Workflows / state machines** **[both]** (gap 3.3 / recs 6, 30):
      `status:enum:draft,sent,paid:workflow` generating transition guards, an
      admin badge and timeline, and transition events.
- [ ] **Notifications and inbox** (recs 8, 9), a consumer of domain events.
- [ ] **Settings registry** (rec 29): `setting.Define("invoice.prefix")` with
      typed access, admin UI and env fallback.
- [ ] **Application-wide search** (rec 10) with provider adapters: Postgres FTS
      first, Meilisearch and Typesense as plugins.

### 7.5 SaaS (recs 14, 15, 27, 31)

Billing, entitlements, payment adapters, invitations, usage limits, automation
rules. A plugin bundle, not core, exactly as the review's own section 36 argues.

- [ ] Payments abstraction with adapters (Stripe, Paystack, Flutterwave, MTN).
- [ ] Subscriptions, plans, entitlements.
- [ ] Automation rules (when / if / then) on top of domain events.
- [ ] Starter kits per vertical **[both]** (gap 4.2 / rec 18).

### 7.6 Reach

- [ ] i18n at the generator, so every new resource emits keys rather than
      English (gap 3.7).
- [x] Extract the offline-sync engine out of `apps/desktop` **[both]**
      (gap 3.10 / rec 26) - v3.148.0. `grit add offline` installs
      `packages/sync`: the same mirror, outbox and version-checked conflicts in
      TypeScript, over a storage interface with IndexedDB, expo-sqlite and
      in-memory adapters. Verified by running 45 checks against a mock server.
- [ ] Media abstraction over storage (rec 16).
- [ ] `grit audit a11y` against the generated admin (gap 4.3).
- [ ] GraphQL as a plugin over the generated handlers (gap 3.5).
- [ ] Preview environments per branch (gap 3.12).
- [ ] Visual application explorer, read/write against the same
      `defineResource` files, not a separate runtime **[both]**
      (gap 3.4 / rec 20).

### 7.7 Not a feature (gap 3.9)

- [ ] Co-maintainers with commit access, an RFC process for community plugins,
      and named sponsors. One maintainer is the adoption blocker no feature
      fixes.

---

## Phase Summary

| Phase | Duration | Focus | Key Deliverable |
|-------|----------|-------|-----------------|
| **Phase 1** | Weeks 1-3 | Foundation | CLI + Go API with auth + Next.js + Docker + GORM Studio |
| **Phase 2** | Weeks 4-6 | Code Generator | `grit generate resource` full-stack code generation |
| **Phase 3** | Weeks 7-12 | Admin Panel | Filament-like resource-based admin panel with tables, forms, widgets |
| **Phase 4** | Weeks 13-16 | Batteries | File storage, email, jobs, cron, cache, AI |
| **Phase 5** | Weeks 17-20 | Launch | Docs, tests, performance, public launch |

---

## Rules for AI Agents

1. **Complete one phase fully before moving to the next.** Do not start Phase 2 until every checkbox in Phase 1 is done.
2. **Test every feature as you build it.** Don't just write code — verify it works.
3. **Follow the folder structure exactly.** The conventions in GRIT.md are non-negotiable.
4. **Generated code must be clean and readable.** A developer should be able to understand and modify any generated file.
5. **The dark theme aesthetic is mandatory.** Every UI component must match the GORM Studio design language.
6. **Commit frequently** with conventional commit messages (feat:, fix:, docs:, refactor:, test:).
7. **If something is unclear, refer to GRIT.md.** That document is the source of truth.

---

*When building with Claude Code, always read CLAUDE.md first for project context.*
