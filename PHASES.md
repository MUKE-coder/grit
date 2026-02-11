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

- [ ] Template engine for Go and TypeScript file generation
- [ ] `grit generate resource <Name>` command:
  - Prompts for fields (name, type, required, unique) interactively
  - Or accepts a definition file: `grit generate resource --from post.yaml`
- [ ] Template files for each generated artifact
- [ ] Smart pluralization (Post → posts, Category → categories)
- [ ] Automatic import management (add imports to existing files)
- [ ] Automatic route registration (append to routes.go)

### 2.2 Generated Go Artifacts

- [ ] **Model** (`internal/models/<name>.go`):
  - GORM struct with proper tags
  - Relationships (belongs_to, has_many)
  - Timestamps and soft deletes
- [ ] **Handler** (`internal/handlers/<name>.go`):
  - `GET /api/<names>` — List with pagination, sorting, filtering, search
  - `GET /api/<names>/:id` — Get by ID with relations
  - `POST /api/<names>` — Create with validation
  - `PUT /api/<names>/:id` — Update with validation
  - `DELETE /api/<names>/:id` — Soft delete
  - `POST /api/<names>/bulk-delete` — Bulk delete
  - `GET /api/<names>/export` — Export as CSV/JSON
- [ ] **Service** (`internal/services/<name>.go`):
  - Business logic layer between handler and model
  - Reusable query scopes (pagination, filtering, sorting)
- [ ] Automatic migration registration

### 2.3 Generated Frontend Artifacts

- [ ] **Zod schema** (`packages/shared/schemas/<name>.ts`):
  - Create schema, update schema, list response schema
  - Proper Zod types matching Go types
- [ ] **TypeScript types** (`packages/shared/types/<name>.ts`):
  - Full type with all fields
  - Create/Update DTOs
  - List response with pagination
- [ ] **React Query hooks** (`apps/web/hooks/use-<names>.ts`):
  - `use<Names>()` — Paginated list query with sorting/filtering
  - `use<Name>(id)` — Single item query
  - `useCreate<Name>()` — Create mutation
  - `useUpdate<Name>()` — Update mutation
  - `useDelete<Name>()` — Delete mutation
  - Automatic cache invalidation on mutations

### 2.4 Type Sync Command

- [ ] `grit sync` command:
  - Reads all Go models in `internal/models/`
  - Generates corresponding TypeScript types and Zod schemas
  - Maps Go types → TypeScript types (uint → number, time.Time → string, etc.)
  - Handles relationships and nested types
  - Handles enums/constants

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

- [ ] Collapsible sidebar with:
  - Logo/brand area
  - Navigation items auto-generated from registered resources
  - Icon support (Lucide icons)
  - Active state highlighting
  - User profile section at bottom
  - Dark/light theme toggle
- [ ] Top navbar:
  - Breadcrumbs
  - Search (global search across resources)
  - Notifications dropdown
  - User menu (profile, settings, logout)
- [ ] Responsive layout (sidebar collapses on mobile)
- [ ] Beautiful page transitions/animations

### 3.2 Resource System

- [ ] `defineResource()` API (see GRIT.md for example)
- [ ] Resource registry (`resources/index.ts`)
- [ ] Auto-generated routes from resources
- [ ] Resource configuration:
  - Name, endpoint, icon
  - Table columns, filters, actions
  - Form fields, validation
  - Dashboard widgets
  - Permissions (which roles can access)
- [ ] Relationship handling in resources

### 3.3 DataTable Component

- [ ] Server-side pagination (communicates with Go API)
- [ ] Column sorting (click header to sort)
- [ ] Column filtering:
  - Text search (global and per-column)
  - Select/dropdown filters
  - Date range filters
  - Number range filters
  - Boolean toggle filters
- [ ] Column features:
  - Resizable columns
  - Show/hide columns toggle
  - Sticky first column (checkbox + ID)
  - Custom cell renderers (badge, currency, date, relative time, image, boolean)
- [ ] Row selection with checkboxes
- [ ] Bulk actions toolbar (delete, export, custom)
- [ ] Row actions (edit, delete, view, custom)
- [ ] Empty state with illustration
- [ ] Loading skeleton
- [ ] Export to CSV / JSON
- [ ] Responsive (horizontal scroll on mobile)
- [ ] Keyboard navigation (arrow keys)

### 3.4 Form Builder

- [ ] Form modal and full-page form views
- [ ] Field types:
  - Text input (with prefix/suffix)
  - Textarea
  - Number (with min/max, step)
  - Select / Multi-select
  - Combobox (searchable select)
  - Date picker
  - Date range picker
  - DateTime picker
  - Toggle / Switch
  - Checkbox / Checkbox group
  - Radio group
  - File upload (single + multiple, drag and drop)
  - Image upload with preview
  - Rich text editor
  - JSON editor
  - Color picker
  - Relation field (searchable dropdown that queries the related API)
  - Repeater (dynamic list of sub-fields)
- [ ] Validation:
  - Zod-based validation from shared schemas
  - Real-time validation (on blur and on change)
  - Server-side error display
- [ ] Form layout:
  - Single column, two column, tabbed
  - Section groups with headers
  - Conditional visibility (show field based on another field's value)
- [ ] Create and edit modes from the same form definition
- [ ] Auto-populated defaults

### 3.5 Dashboard & Widgets

- [ ] Dashboard page as the admin home
- [ ] Widget types:
  - **Stats card** — Number + label + change percentage + icon
  - **Line chart** — Time series data
  - **Bar chart** — Categorical data
  - **Pie/Donut chart** — Proportional data
  - **Recent activity** — List of recent events
  - **Table widget** — Mini data table (e.g., top 5 customers)
  - **Custom widget** — Render any React component
- [ ] Widget grid layout (responsive, configurable)
- [ ] Widgets fetch data from the Go API
- [ ] Charting library: Recharts (already available in React artifacts)
- [ ] Animated counters for stats cards

### 3.6 Admin Theme

- [ ] Dark theme (default) matching GORM Studio:
  - Background: `#0a0a0f`, `#111118`, `#1a1a24`
  - Accent: `#6c5ce7`
  - Fonts: DM Sans + JetBrains Mono
- [ ] Light theme option
- [ ] Theme toggle in sidebar
- [ ] CRM-inspired aesthetic:
  - Generous spacing
  - Subtle animations
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

- [ ] Storage abstraction layer (`internal/storage/storage.go`):
  - Interface: `Upload()`, `Download()`, `Delete()`, `GetURL()`, `GetSignedURL()`
  - S3 driver (AWS S3, Cloudflare R2)
  - Local driver (MinIO in dev, local filesystem fallback)
  - Configuration via `.env` (STORAGE_DRIVER, S3_BUCKET, S3_REGION, etc.)
- [ ] Upload handler: `POST /api/uploads`
  - File size limits, allowed MIME types (configurable)
  - Returns file URL and metadata
- [ ] Image processing on upload:
  - Thumbnail generation
  - Resize to max dimensions
- [ ] File upload React component:
  - Drag and drop
  - Progress bar
  - Preview (images)
  - Multiple file support
- [ ] File management in admin panel:
  - Browse uploaded files
  - Delete files
  - View file details

### 4.2 Email System

- [ ] Mail service (`internal/mail/mailer.go`):
  - Resend client integration
  - Send method: `SendMail(to, subject, template, data)`
  - HTML templates using Go `html/template`
- [ ] Built-in email templates:
  - Welcome email
  - Password reset
  - Email verification
  - Notification
- [ ] Template preview in admin panel (dev only)
- [ ] Email configuration via `.env`

### 4.3 Background Jobs

- [ ] Job queue system (`internal/jobs/`):
  - Redis-backed queue using `asynq` library
  - Job definition interface: `Process(ctx, payload) error`
  - Enqueue: `jobs.Enqueue("send-email", payload)`
  - Job priorities (critical, default, low)
  - Retry with exponential backoff
  - Dead letter queue
- [ ] Built-in jobs:
  - Send email job
  - Process image job
  - Cleanup expired tokens job
- [ ] Job dashboard in admin panel:
  - Queue stats (pending, active, completed, failed)
  - View failed jobs with error details
  - Retry failed jobs
  - Clear queues

### 4.4 Cron Scheduler

- [ ] Cron service (`internal/cron/cron.go`):
  - Schedule tasks with cron expressions
  - Built-in tasks: cleanup expired tokens, prune old logs
  - Easy registration: `cron.Register("0 * * * *", cleanupTokens)`
- [ ] Cron status in admin panel

### 4.5 Redis Caching

- [ ] Cache service (`internal/cache/cache.go`):
  - `Get(key)`, `Set(key, value, ttl)`, `Delete(key)`, `Flush()`
  - JSON serialization for complex values
  - Cache middleware for API responses
- [ ] Available via `grit add cache`

### 4.6 AI Integration

- [ ] AI service (`internal/ai/ai.go`):
  - Anthropic Claude API client
  - OpenAI API client
  - Simple interface: `ai.Complete(prompt)`, `ai.Chat(messages)`
- [ ] Frontend: Vercel AI SDK pre-configured
- [ ] Chat UI component in shared components
- [ ] Streaming support
- [ ] Configuration via `.env` (AI_PROVIDER, AI_API_KEY, AI_MODEL)

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

- [ ] Documentation website (Nextra or Mintlify or custom)
- [ ] Pages:
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
- [ ] Interactive examples
- [ ] Copy-paste code blocks
- [ ] Search functionality
- [ ] Dark theme docs site

### 5.2 Testing

- [ ] Go API tests:
  - Unit tests for services
  - Integration tests for handlers
  - Auth flow tests
  - Database tests (SQLite in-memory)
- [ ] Frontend tests:
  - Component tests (React Testing Library)
  - Hook tests
  - E2E tests (Playwright) for auth flow and admin panel
- [ ] CLI tests:
  - Scaffold command output validation
  - Generate command output validation
- [ ] CI/CD:
  - GitHub Actions for Go tests, frontend tests, linting
  - Automated build verification

### 5.3 Performance

- [ ] Go API:
  - Connection pooling optimization
  - Query optimization (preloading, indexes)
  - Response compression (gzip)
  - Caching headers
- [ ] Frontend:
  - Image optimization
  - Bundle analysis and code splitting
  - Lazy loading for admin panel routes
- [ ] Benchmarks:
  - API response times
  - Concurrent connection handling
  - Comparison with Laravel/Next.js

### 5.4 Launch Preparation

- [ ] README.md with GIF demo
- [ ] Landing page at `gritframework.dev`
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
