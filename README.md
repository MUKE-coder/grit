# Grit

**Go + React. Built with Grit.**

Grit is a full-stack meta-framework that fuses Go (Gin + GORM) with Next.js (React + TypeScript) in a monorepo. One command to scaffold a complete project with authentication, admin panel, database browser, and Docker setup.

## Remove grit Old version

```bash
rm $(which grit)
```

## Install

```bash
go install github.com/MUKE-coder/grit/cmd/grit@latest
```

## Quick Start

```bash
# Create a new project
grit new myapp

# Start infrastructure (PostgreSQL, Redis, MinIO, Mailhog)
cd myapp
docker compose up -d

# Start the Go API
cd apps/api
go mod tidy
go run cmd/server/main.go

# In another terminal — start the frontend
cd myapp
pnpm install
cd apps/web && pnpm dev
```

Open http://localhost:3000 — register a user, log in, and see the dashboard.

## What You Get

| Service       | URL                          |
| ------------- | ---------------------------- |
| Web App       | http://localhost:3000        |
| Admin Panel   | http://localhost:3001        |
| Go API        | http://localhost:8080        |
| GORM Studio   | http://localhost:8080/studio |
| Mailhog       | http://localhost:8025        |
| MinIO Console | http://localhost:9001        |

## Project Structure (Generated)

```
myapp/
├── apps/
│   ├── api/              # Go backend (Gin + GORM)
│   │   ├── cmd/server/   # Entry point
│   │   └── internal/     # Config, models, handlers, middleware, services, routes
│   ├── web/              # Next.js frontend (App Router)
│   └── admin/            # Next.js admin panel
├── packages/
│   └── shared/           # Zod schemas, TypeScript types, constants
├── docker-compose.yml    # PostgreSQL 16, Redis 7, MinIO, Mailhog
└── turbo.json            # Monorepo task runner
```

## Features

- **JWT Authentication** — Register, login, refresh tokens, role-based access
- **User Management** — CRUD with pagination, search, sorting
- **GORM Studio** — Visual database browser at `/studio`
- **Dark Theme** — Premium dark UI across all apps
- **Admin Panel** — Data tables, stats cards, user management
- **Shared Types** — Zod schemas + TypeScript types shared between apps
- **Docker Ready** — Dev and production Docker Compose setups
- **API-Only Mode** — `grit new myapp --api` for backend-only projects

## Tech Stack

| Layer         | Technology                      |
| ------------- | ------------------------------- |
| Backend       | Go + Gin + GORM                 |
| Frontend      | Next.js 14 (App Router) + React |
| Styling       | Tailwind CSS + shadcn/ui        |
| Database      | PostgreSQL                      |
| Cache         | Redis                           |
| Validation    | Zod                             |
| Data Fetching | React Query (TanStack Query)    |
| Monorepo      | Turborepo + pnpm                |
| DB Browser    | GORM Studio                     |

## CLI Commands

```bash
grit new <name>         # Scaffold full monorepo
grit new <name> --api   # Scaffold Go API only
grit version            # Print CLI version
```

## License

MIT
