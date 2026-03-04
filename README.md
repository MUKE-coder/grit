# Grit

**Go + React. Built with Grit.**

Grit is a full-stack meta-framework that fuses Go (Gin + GORM) with Next.js (React + TypeScript) in a monorepo. One command to scaffold a complete production-ready project with authentication, admin panel, code generation, file storage, email, background jobs, Redis caching, AI integration, and Docker setup.

## Install

```bash
go install github.com/MUKE-coder/grit/v2/cmd/grit@latest
```

To remove a previous version first:

```bash
rm $(which grit)
```

## Quick Start

```bash
# Create a new project
grit new myapp

# Start infrastructure (PostgreSQL, Redis, MinIO, Mailhog)
cd myapp
docker compose up -d

# Start the Go API
cd apps/api && go run cmd/server/main.go

# In another terminal — start the frontend
cd myapp && pnpm install && pnpm dev
```

Open http://localhost:3000 — register, log in, see the dashboard.

## What You Get

| Service        | URL                               |
| -------------- | --------------------------------- |
| Web App        | http://localhost:3000             |
| Admin Panel    | http://localhost:3001             |
| Go API         | http://localhost:8080             |
| GORM Studio    | http://localhost:8080/studio      |
| API Docs       | http://localhost:8080/docs        |
| Pulse          | http://localhost:8080/pulse/ui/   |
| Sentinel       | http://localhost:8080/sentinel/ui |
| Mailhog        | http://localhost:8025             |
| MinIO Console  | http://localhost:9001             |

## Project Structure (Generated)

```
myapp/
├── apps/
│   ├── api/                    # Go backend (Gin + GORM)
│   │   ├── cmd/server/         # Entry point
│   │   └── internal/
│   │       ├── config/         # Environment configuration
│   │       ├── database/       # GORM connection pool
│   │       ├── models/         # GORM models
│   │       ├── handlers/       # HTTP handlers
│   │       ├── services/       # Business logic
│   │       ├── middleware/     # Auth, CORS, gzip, request ID
│   │       ├── routes/         # Route definitions
│   │       ├── cache/          # Redis cache service
│   │       ├── storage/        # S3-compatible file storage
│   │       ├── mail/           # Resend email service
│   │       ├── jobs/           # asynq background jobs
│   │       ├── cron/           # asynq cron scheduler
│   │       └── ai/             # Claude / OpenAI integration
│   ├── web/                    # Next.js frontend (App Router)
│   └── admin/                  # Next.js admin panel
├── packages/
│   └── shared/                 # Zod schemas, TypeScript types, constants
├── docker-compose.yml          # PostgreSQL 16, Redis 7, MinIO, Mailhog
├── docker-compose.prod.yml     # Production Docker Compose
└── turbo.json                  # Monorepo task runner
```

## Features

### Foundation
- **JWT Authentication** — Register, login, refresh tokens, role-based access (ADMIN / EDITOR / USER)
- **OAuth2 Social Login** — Google and GitHub via `goth`
- **User Management** — CRUD with pagination, search, sorting
- **GORM Studio** — Visual database browser at `/studio`
- **API Documentation** — Auto-generated Scalar docs at `/docs`
- **Dark Theme** — Premium dark UI across all apps
- **Docker Ready** — Dev and production Docker Compose setups

### Code Generator
- **`grit generate resource <Name>`** — Full-stack resource in one command:
  - Go model with GORM tags + auto-migration
  - REST handler with pagination, search, sorting
  - Zod schema + TypeScript types in `packages/shared`
  - React Query hooks in admin and web
  - Admin panel page with DataTable
  - Automatic route injection via markers
- **`grit sync`** — Sync Go models → TypeScript types and Zod schemas
- **Inline fields**: `--fields "title:string,content:text,published:bool"`
- **YAML definition**: `--from post.yaml`

### Admin Panel
- **Resource system** — `defineResource()` for data-driven admin pages
- **DataTable** — Server-side pagination, sorting, filtering, column visibility, export CSV/JSON
- **Form builder** — 8+ field types with Zod validation, create/edit from same definition
- **Dashboard widgets** — Stats cards, line charts, bar charts, activity feed
- **Collapsible sidebar** — Auto-generated navigation from registered resources
- **Dark/light theme toggle**

### Batteries Included
- **File Storage** — S3/R2/B2/MinIO with image processing (thumbnails via background jobs)
- **Email** — Resend integration with HTML templates (welcome, reset, verify, notify)
- **Background Jobs** — Redis-backed `asynq` queue with admin dashboard (retry, clear, stats)
- **Cron Scheduler** — `asynq` Scheduler with admin task list
- **Redis Caching** — Get/Set/Delete/Flush with cache middleware for GET responses
- **AI Integration** — Claude and OpenAI APIs with SSE streaming

### Desktop (Wails)
- **`grit new-desktop <name>`** — Scaffold a native desktop app (Go + React + SQLite)
- **Wails bindings** — Go functions directly callable from React (no HTTP)
- **Auth, CRUD, Export** — Login/register, blog + contact CRUD, PDF/Excel export
- **GORM Studio** — Standalone database browser at `cmd/studio`
- **Resource generation** — `grit generate resource` works for desktop projects too
- **Custom title bar** — Frameless window with draggable title bar
- **Dark theme** — Grit dark theme with accent purple

### Security & Observability
- **Sentinel** — WAF, rate limiting per IP/route, auth shield, anomaly detection at `/sentinel/ui`
- **Pulse** — Request tracing, DB monitoring, metrics, health checks at `/pulse/ui/`
- **Gzip compression** — Response compression middleware (Best Speed level)
- **Request ID tracing** — `X-Request-ID` header injected on every request

## CLI Commands

```bash
# Project scaffolding
grit new <name>                    # Full monorepo (web + admin + API)
grit new <name> --api              # Go API only
grit new <name> --mobile           # Mobile-first (API + Expo)
grit new <name> --full             # Everything including Expo + docs site
grit new <name> --style modern     # Admin style: default|modern|minimal|glass
grit new-desktop <name>            # Native desktop app (Wails + Go + React)

# Code generation (works for both web and desktop projects)
grit generate resource <Name>                      # Interactive field prompts
grit generate resource <Name> --fields "..."      # Inline field definition
grit generate resource <Name> --from post.yaml    # YAML definition file
grit remove resource <Name>                       # Remove generated resource

# Type sync (web projects)
grit sync                           # Go models -> TypeScript + Zod

# Development
grit dev                            # Start all frontend apps (pnpm dev)
grit server                         # Start Go API server
grit start                          # Auto-detect project type, start dev server
grit compile                        # Build desktop executable (wails build)
grit studio                         # Open database browser (web + desktop)

# Database
grit migrate                        # Run database migrations
grit migrate --fresh                # Drop + re-migrate
grit seed                           # Run database seeder

# Meta
grit version                        # Print CLI version (2.0.0)
```

## Field Types (Code Generator)

| Type           | Go Type                        | TypeScript      | Notes                         |
| -------------- | ------------------------------ | --------------- | ----------------------------- |
| `string`       | `string`                       | `string`        | Required by default           |
| `text`         | `string`                       | `string`        | GORM `type:text`              |
| `richtext`     | `string`                       | `string`        | GORM `type:text`              |
| `int`          | `int`                          | `number`        |                               |
| `uint`         | `uint`                         | `number`        |                               |
| `float`        | `float64`                      | `number`        |                               |
| `bool`         | `bool`                         | `boolean`       |                               |
| `datetime`     | `*time.Time`                   | `string\|null`  |                               |
| `date`         | `*time.Time`                   | `string\|null`  |                               |
| `slug`         | `string`                       | `string`        | Auto-unique index             |
| `belongs_to`   | `uint`                         | `number`        | FK + index                    |
| `many_to_many` | `[]uint`                       | `number[]`      | Requires target model         |
| `string_array` | `datatypes.JSONSlice[string]`  | `string[]`      | GORM `type:json`              |

**Modifiers:** `:unique`, `:optional`, `:slug:<source>`, `:belongs_to:<Model>`, `:many_to_many:<Model>`

## Tech Stack

| Layer          | Technology                      |
| -------------- | ------------------------------- |
| Backend        | Go 1.21+ · Gin · GORM           |
| Frontend       | Next.js 14 (App Router) · React |
| Styling        | Tailwind CSS · shadcn/ui        |
| Database       | PostgreSQL (dev: Docker)        |
| Cache / Queue  | Redis · asynq                   |
| File Storage   | S3-compatible (MinIO/R2/B2)     |
| Email          | Resend                          |
| AI             | Anthropic Claude · OpenAI       |
| Validation     | Zod (shared between apps)       |
| Data Fetching  | TanStack Query                  |
| Monorepo       | Turborepo · pnpm                |
| DB Browser     | GORM Studio                     |
| Security       | Sentinel (WAF + rate limiting)  |
| Observability  | Pulse (tracing + metrics)       |

## License

MIT
