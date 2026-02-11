# CLAUDE.md — Grit Framework Project Context

> **This file exists to help Claude Code (and other AI agents) maintain context across sessions. Read this FIRST before doing any work.**

---

## What is This Project?

**Grit** is a full-stack meta-framework that combines:
- **Go backend** (Gin web framework + GORM ORM)
- **Next.js frontend** (React + App Router + TypeScript)
- **Admin panel** (Filament-like resource-based admin dashboard)
- All in a **monorepo** with shared types and validation

**Tagline:** Go + React. Built with Grit.

**Creator:** MUKE-coder (GitHub: github.com/MUKE-coder)

**Predecessor:** This project evolved from GORM Studio (github.com/MUKE-coder/gorm-studio), a visual database browser for GORM that is now embedded within Grit.

---

## Key Documents — READ THESE

| Document | Location | Purpose |
|----------|----------|---------|
| **GRIT.md** | `/GRIT.md` | Master specification — what Grit is, features, architecture, folder structure, tech stack, monetization. THE SOURCE OF TRUTH. |
| **PHASES.md** | `/PHASES.md` | Development phases with checkboxes. Build one phase at a time. Never skip ahead. |
| **This file** | `/CLAUDE.md` | Quick context for AI agents. You're reading it. |

**Before writing ANY code, read GRIT.md and PHASES.md to understand what you're building.**

---

## Current State

<!-- UPDATE THIS SECTION AS PHASES ARE COMPLETED -->

**Current Phase:** Phase 4 — Batteries
**Status:** Complete
**Last Updated:** 2026-02-11

### What's Been Built
- [x] CLI Scaffolder (`grit new <project-name>` with `--api` flag, name validation, ASCII art)
- [x] Go API: config, database, models (User, Upload), auth handlers, JWT service, middleware (auth, CORS, logger, cache), routes with Services struct
- [x] GORM Studio integrated at `/studio`
- [x] Next.js Web App: auth pages (login, register, forgot-password), dashboard layout with sidebar, stats cards
- [x] Admin Panel: layout with sidebar, navbar, dashboard with stats, users management page with data table
- [x] Shared Package: Zod schemas, TypeScript types (User, Upload, API), constants (routes for uploads, AI, admin)
- [x] Docker Setup: docker-compose.yml (PostgreSQL, Redis, MinIO, Mailhog), docker-compose.prod.yml, Dockerfiles
- [x] Dev Experience: README, .env.example, turbo.json, pnpm-workspace.yaml, root package.json

- [x] Phase 2 — Code Generator: `grit generate resource` (Go model, service, handler, Zod schemas, TS types, React hooks, admin resource definition + page), `grit sync` (Go types → TypeScript)
- [x] Phase 3 — Admin Panel: runtime resource definitions via `defineResource()`, advanced DataTable (sort, filter, select, pagination), FormBuilder (8 field types), dashboard widgets (stats, charts, activity), collapsible sidebar with Lucide icons, dark/light theme toggle, resource registry
- [x] Phase 4 — Batteries: Redis cache service + middleware, S3 file storage (MinIO/R2/B2) + image processing + upload handler, Resend email service + 4 HTML templates, asynq background jobs (email/image/cleanup workers) + admin dashboard, asynq cron scheduler, AI integration (Claude + OpenAI with streaming), admin System pages (Jobs, Files, Cron, Mail Preview)

### What's In Progress
- Nothing — Phase 4 complete

### What's Next
- [ ] Phase 5: Polish & Launch (docs, tests, performance)

---

## Tech Stack (Do Not Deviate)

| Layer | Technology | Notes |
|-------|-----------|-------|
| Backend | **Go 1.21+** | Use `go mod` for dependencies |
| Web framework | **Gin** | Not Echo, not Fiber, not Chi |
| ORM | **GORM** | Not sqlx, not sqlc, not ent |
| Database | **PostgreSQL** (prod), **SQLite** (quick start/testing) | |
| Frontend | **Next.js 14+** with App Router | NOT Pages Router |
| Styling | **Tailwind CSS** + **shadcn/ui** | NOT Material UI, NOT Chakra |
| Data fetching | **React Query (TanStack Query)** | NOT SWR, NOT Apollo |
| Validation | **Zod** | Shared between frontend and generated from Go types |
| Monorepo | **Turborepo** + **pnpm** | NOT npm, NOT yarn |
| Cache/Queue | **Redis** | Using `asynq` for job queues |
| File storage | **S3-compatible** (AWS S3, Cloudflare R2, MinIO) | |
| Email | **Resend** | |
| Containerization | **Docker** + **Docker Compose** | |
| DB browser | **GORM Studio** | Our own tool, embedded in the API |

---

## Architecture Rules

### Folder Structure
The folder structure is defined in GRIT.md. Follow it exactly. Here's the abbreviated version:

```
project-root/
├── grit.config.ts
├── docker-compose.yml
├── packages/shared/          # Zod schemas, TS types, constants
├── apps/
│   ├── api/                  # Go backend (Gin + GORM)
│   │   ├── cmd/server/       # Entry point
│   │   └── internal/         # All Go code (models, handlers, services, middleware, etc.)
│   ├── web/                  # Next.js main frontend
│   └── admin/                # Next.js admin panel
└── grit/                     # CLI tool (Go)
```

### Naming Conventions

| Thing | Convention | Example |
|-------|-----------|---------|
| Go files | snake_case | `user_handler.go` |
| Go structs | PascalCase | `type User struct` |
| Go functions | PascalCase (exported), camelCase (unexported) | `GetUsers`, `parseToken` |
| TypeScript files | kebab-case | `use-users.ts`, `api-client.ts` |
| React components | PascalCase files | `DataTable.tsx`, `StatsCard.tsx` |
| API routes | plural, lowercase | `/api/users`, `/api/posts` |
| Database tables | plural, snake_case | `users`, `blog_posts` |
| Zod schemas | PascalCase + Schema | `UserSchema`, `CreatePostSchema` |
| CSS classes | Tailwind utilities | No custom CSS unless absolutely necessary |

### Code Style

**Go:**
- Always handle errors explicitly. Never ignore errors with `_`.
- Use `fmt.Errorf("context: %w", err)` for error wrapping.
- Keep handlers thin — business logic goes in services.
- Use struct tags: `gorm:"..."`, `json:"..."`, `binding:"..."`.
- Group imports: stdlib, external, internal.

**TypeScript/React:**
- Use functional components only. No class components.
- Use React hooks. State with `useState`, effects with `useEffect`.
- All data fetching through React Query hooks. No `fetch` in components.
- Validate all API inputs with Zod.
- Export types explicitly. Use `interface` for objects, `type` for unions.

**Both:**
- Meaningful variable names. No single letters except in loops.
- Comments for WHY, not WHAT.
- Keep functions small (<50 lines preferred).

### Design System

**Theme (Dark Mode Default):**
```
--bg-primary:    #0a0a0f
--bg-secondary:  #111118
--bg-tertiary:   #1a1a24
--bg-elevated:   #22222e
--bg-hover:      #2a2a38
--border:        #2a2a3a
--text-primary:  #e8e8f0
--text-secondary:#9090a8
--text-muted:    #606078
--accent:        #6c5ce7  (purple)
--accent-hover:  #7c6cf7
--success:       #00b894
--danger:        #ff6b6b
--warning:       #fdcb6e
--info:          #74b9ff
```

**Fonts:**
- UI: `DM Sans` (weights: 400, 500, 600, 700)
- Code: `JetBrains Mono` (weights: 400, 500, 600)

**Design Feel:** Premium CRM / dark mode SaaS tool. Not generic Bootstrap. Not Material Design. Think Linear, Vercel Dashboard, or Raycast — dark, polished, fast.

---

## API Response Format

All API endpoints must follow this format:

### Success (single item):
```json
{
  "data": { ... },
  "message": "User created successfully"
}
```

### Success (list with pagination):
```json
{
  "data": [ ... ],
  "meta": {
    "total": 100,
    "page": 1,
    "page_size": 20,
    "pages": 5
  }
}
```

### Error:
```json
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Email is required",
    "details": {
      "email": "This field is required"
    }
  }
}
```

HTTP status codes: 200 (OK), 201 (Created), 400 (Bad Request), 401 (Unauthorized), 403 (Forbidden), 404 (Not Found), 422 (Validation Error), 500 (Server Error).

---

## Git Commit Convention

Use conventional commits:
```
feat: add user authentication handlers
fix: resolve JWT token refresh race condition
docs: update API reference for user endpoints
refactor: extract pagination logic to shared service
test: add integration tests for auth flow
chore: update Go dependencies
style: format code with gofmt
```

---

## Common Pitfalls — Avoid These

1. **Don't use Pages Router in Next.js.** We use App Router exclusively.
2. **Don't install dependencies not in the tech stack** without explicit approval.
3. **Don't create custom CSS files.** Use Tailwind utilities and shadcn/ui.
4. **Don't put business logic in handlers.** Handlers call services, services contain logic.
5. **Don't hardcode values.** Use `.env` variables and config structs.
6. **Don't skip error handling in Go.** Every error must be handled.
7. **Don't use `any` type in TypeScript.** Define proper types.
8. **Don't deviate from the folder structure.** It's the foundation of the framework's conventions.
9. **Don't build features from a later phase.** Follow PHASES.md sequentially.
10. **Don't compromise on the dark theme aesthetic.** Every UI must look premium.

---

## Quick Reference Commands

```bash
# Development
grit dev                          # Start all services
grit generate resource <n>     # Generate full-stack resource
grit migrate                      # Run DB migrations
grit sync                         # Sync Go types → TypeScript
grit studio                       # Open GORM Studio

# Docker
docker compose up -d              # Start infrastructure (DB, Redis, MinIO)
docker compose down               # Stop everything

# Go API
cd apps/api && go run cmd/server/main.go   # Run API directly
cd apps/api && air                          # Run with hot reload

# Frontend
cd apps/web && pnpm dev           # Run Next.js web app
cd apps/admin && pnpm dev         # Run admin panel

# Monorepo
pnpm install                      # Install all dependencies
turbo build                       # Build all apps
turbo dev                         # Dev mode for all apps
```

---

## Session Checklist for AI Agents

When starting a new session:

1. ✅ Read this file (CLAUDE.md)
2. ✅ Read GRIT.md for full specification
3. ✅ Read PHASES.md to know the current phase
4. ✅ Check the "Current State" section above for progress
5. ✅ Look at existing code to understand what's been built
6. ✅ Continue from where the last session left off
7. ✅ Update the "Current State" section when you complete tasks
8. ✅ Commit work with conventional commit messages

---

*Last context update: 2026-02-11 — Phase 3 complete. Admin panel with runtime resource definitions, advanced DataTable, FormBuilder, dashboard widgets, theme toggle, Lucide icons, resource registry.*
