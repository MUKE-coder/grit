# Job Portal -- Single App Architecture + TanStack Router (Vite)

A full-featured job portal built with Grit using the **Single** architecture:
one Go binary that serves both the API and the React frontend via `go:embed`.

- No Turborepo. No `apps/` directory. Flat project structure.
- One binary, one port, one deployment unit.
- Like Next.js fullstack or Laravel -- everything in one place.

## Why Single?

| | Triple / Double | Single |
|---|---|---|
| Deployment | 2-3 containers + nginx | **1 binary** |
| Project structure | Turborepo monorepo | **Flat (main.go at root)** |
| Module path | `job-portal/apps/api` | **`job-portal`** |
| Dev mode | Turbo orchestrates | **2 terminals: Go + Vite** |
| Production | Go + nginx + static | **One Go binary (frontend embedded)** |
| Best for | Large SaaS, teams | **Internal tools, microservices, solo devs** |

Choose Single when you want the simplest possible deployment. Build the Go
binary, upload it to a server, run it. The frontend is embedded in the binary
and served on the same port as the API.

## Quick Start

### Prerequisites
- Go 1.21+
- Node.js 18+
- pnpm 8+
- Docker (for PostgreSQL, Redis, MinIO, Mailhog)

### Setup

1. Install Grit CLI:
   ```bash
   go install github.com/MUKE-coder/grit/v3/cmd/grit@latest
   ```

2. Create the project:
   ```bash
   grit new job-portal --single --vite
   ```

3. Start Docker services:
   ```bash
   cd job-portal
   docker compose up -d
   ```

4. Generate resources:
   ```bash
   grit generate resource Category --fields "name:string:unique,slug:slug:name"
   grit generate resource Company --fields "name:string,slug:slug:name,description:richtext,logo:string,website:string:optional,email:string"
   grit generate resource Job --fields "title:string,description:richtext,company_id:belongs_to:Company,location:string,salary_min:float,salary_max:float,type:string,experience:string,active:bool,deadline:date"
   grit generate resource Application --fields "job_id:belongs_to:Job,applicant_name:string,applicant_email:string,resume:string,cover_letter:text:optional,status:string"
   ```

5. Run migrations and seed:
   ```bash
   grit migrate
   grit seed
   ```

6. Start development servers:
   ```bash
   # Terminal 1: Go API
   go run main.go

   # Terminal 2: Vite dev server
   cd frontend && pnpm install && pnpm dev
   ```

7. Open:
   - Web: http://localhost:5173 (Vite proxies API calls to :8080)
   - API: http://localhost:8080
   - GORM Studio: http://localhost:8080/studio
   - Mailhog: http://localhost:8025

## Project Structure

```
job-portal/
├── main.go                        # Entry point: go:embed + Gin server
├── go.mod                         # Module: job-portal
├── go.sum
├── .env                           # Environment variables
├── .env.example
├── Makefile                       # dev, build, clean targets
├── internal/                      # Go backend
│   ├── config/                    # Environment config
│   ├── database/                  # GORM setup + migrations
│   ├── models/                    # Job, Company, Application, Category, User
│   ├── handlers/                  # REST handlers
│   ├── services/                  # Business logic
│   ├── middleware/                # Auth, CORS, cache, rate limiting
│   └── routes/                    # Route registration
├── frontend/                      # React SPA (TanStack Router)
│   ├── index.html
│   ├── vite.config.ts
│   ├── tailwind.config.ts
│   ├── tsconfig.json
│   ├── package.json
│   └── src/
│       ├── main.tsx               # App entry point
│       ├── routes/                # File-based routing (TanStack Router)
│       │   ├── __root.tsx         # Root layout (navbar, footer)
│       │   ├── index.tsx          # Landing page
│       │   ├── jobs/
│       │   │   ├── index.tsx      # Job listings
│       │   │   └── $id.tsx        # Job detail
│       │   ├── companies/
│       │   │   ├── index.tsx      # Company directory
│       │   │   └── $slug.tsx      # Company profile
│       │   ├── apply/
│       │   │   └── $jobId.tsx     # Application form
│       │   ├── auth/
│       │   │   ├── login.tsx
│       │   │   └── register.tsx
│       │   └── admin/             # Role-protected management pages
│       │       ├── index.tsx      # Admin dashboard
│       │       ├── jobs.tsx
│       │       ├── companies.tsx
│       │       ├── applications.tsx
│       │       └── categories.tsx
│       ├── components/            # Shared UI components
│       ├── hooks/                 # React Query hooks
│       └── lib/                   # API client, utils
├── docker-compose.yml             # Dev services (Postgres, Redis, MinIO)
└── docker-compose.prod.yml        # Production deployment
```

Note: no `apps/`, no `packages/`, no `turbo.json`, no `pnpm-workspace.yaml`.
Types are defined in `frontend/src/lib/types.ts` and validated with Zod
in `frontend/src/lib/schemas.ts`.

## How go:embed Works

The `main.go` file embeds the built frontend into the Go binary:

```go
package main

import (
    "embed"
    "io/fs"
    "net/http"

    "job-portal/internal/config"
    "job-portal/internal/database"
    "job-portal/internal/routes"

    "github.com/gin-gonic/gin"
)

//go:embed frontend/dist/*
var frontendFS embed.FS

func main() {
    cfg := config.Load()
    db := database.Connect(cfg)

    r := gin.Default()

    // API routes under /api/*
    routes.Register(r, db, cfg)

    // Serve embedded frontend for everything else
    distFS, _ := fs.Sub(frontendFS, "frontend/dist")
    fileServer := http.FileServer(http.FS(distFS))
    r.NoRoute(gin.WrapH(fileServer))

    r.Run(":" + cfg.Port)
}
```

In **development**, the `frontend/dist` directory does not exist, so you run
the Vite dev server separately on port 5173. The Vite config proxies `/api`
requests to Go on port 8080.

In **production**, you build the frontend first (`cd frontend && pnpm build`),
then build the Go binary (`go build -o job-portal main.go`). The resulting
binary contains everything -- API logic and frontend assets.

## Development Workflow

```bash
# Terminal 1: Go API with hot reload (using air)
air
# Or without air:
go run main.go

# Terminal 2: Vite dev server
cd frontend && pnpm dev
```

The Vite config includes a proxy so API calls from the browser go to Go:

```typescript
// frontend/vite.config.ts
export default defineConfig({
  server: {
    port: 5173,
    proxy: {
      "/api": "http://localhost:8080",
    },
  },
});
```

During development, open http://localhost:5173 (Vite serves the frontend
with HMR, proxies API calls to Go).

## Production Build

```bash
# 1. Build the frontend
cd frontend && pnpm install && pnpm build

# 2. Build the Go binary (embeds frontend/dist)
cd .. && go build -o job-portal main.go

# 3. Run the single binary
./job-portal
# Serves API + frontend on port 8080
```

That is it. One binary. Copy it to any server and run it.

## Features Demonstrated

### Authentication
- Email/password registration and login
- OAuth2 social login (Google, GitHub)
- JWT access tokens (15 min) + refresh tokens (7 days)
- TOTP two-factor authentication with backup codes
- Role-based access: ADMIN (full access), EDITOR (manage jobs), USER (apply)

### Data Management
- DataTable with server-side pagination, sorting, filtering
- FormBuilder with rich text editor for job descriptions
- File upload for company logos and resumes
- Application status workflow: Pending -> Reviewed -> Interview -> Accepted/Rejected
- Export to CSV/Excel

### Public Website
- Landing page with hero, featured jobs, and category cards
- Job listings with search, location filter, and salary range filter
- Company profiles with job listings
- Application form with resume upload
- Client-side routing (TanStack Router SPA)

### Emails
- Welcome email on registration
- Application confirmation email
- Application status update notifications

### TanStack Router Specifics
- File-based routing in `frontend/src/routes/`
- Type-safe route params
- `beforeLoad` guards for auth and role checks
- `__root.tsx` for root layout
- Vite for fast HMR and small production bundles

## Environment Variables

See `.env.example` for the full list.

Key variables:
- `APP_URL` -- Single URL for the whole app (API + frontend)
- `DATABASE_URL` -- PostgreSQL connection string
- `JWT_SECRET` -- Must be changed in production
- `STORAGE_ENDPOINT` -- MinIO (dev) or S3/R2 (prod)
- `RESEND_API_KEY` -- For sending emails
- `GOOGLE_CLIENT_ID` -- For OAuth (optional)
- `GITHUB_CLIENT_ID` -- For OAuth (optional)

## Deployment

### Option 1: Single binary on any VPS

```bash
# Build
cd frontend && pnpm build && cd ..
go build -o job-portal main.go

# Copy to server
scp job-portal deploy@server.com:/opt/job-portal/
scp .env deploy@server.com:/opt/job-portal/

# Run
ssh deploy@server.com
cd /opt/job-portal && ./job-portal
```

Use systemd to keep it running, Caddy or nginx for TLS.

### Option 2: grit deploy (VPS)

```bash
grit deploy --host deploy@server.com --domain jobs.example.com
```

### Option 3: Docker Compose

```bash
docker compose -f docker-compose.prod.yml up -d
```

This runs just the Go binary (with embedded frontend) + PostgreSQL + Redis +
MinIO. No separate frontend container.

## LLM Instructions

If you are an AI assistant working on this project:
1. Read `.claude/skills/grit/SKILL.md` for Grit conventions
2. Use `grit generate resource` for new resources
3. Never delete `// grit:*` marker comments
4. Follow the API response format: `{ data, message }` / `{ data, meta }` / `{ error: { code, message } }`
5. Use React Query for data fetching, Zod for validation
6. The module path is `job-portal` (not `job-portal/apps/api`)
7. Frontend code is in `frontend/` (not `apps/web/`)
8. Report bugs at https://github.com/MUKE-coder/grit/issues

## License

MIT
