# Phase 1 — Claude Code Prompt

Paste everything between **PROMPT START** and **PROMPT END** into Claude Code.

---

## PROMPT START

I'm building **Grit** — a full-stack meta-framework that fuses Go (Gin + GORM) with Next.js (React + TypeScript) in a monorepo. Think Laravel's DX + Go's performance + React's ecosystem.

Before writing any code, read these three files in the project root:

1. **CLAUDE.md** — Quick context, tech stack, architecture rules, naming conventions, design system, API format. Read this first.
2. **GRIT.md** — Master specification. Everything about what Grit is, full folder structure, features, tech decisions.
3. **PHASES.md** — Development phases with checkboxes. We are building **Phase 1 only**.

Read all three files now before proceeding.

---

### YOUR TASK: Build Phase 1 — Foundation

Build everything listed under **Phase 1** in PHASES.md. Here's a summary of what needs to be built, in order:

#### Step 1: CLI Scaffolder (`grit` CLI tool)

Create a Go CLI tool using the `cobra` library in the `grit/` directory.

The `grit new <project-name>` command should:
- Validate the project name (lowercase, alphanumeric, hyphens)
- Create the full monorepo folder structure as defined in GRIT.md
- Write all template files (Go API, Next.js apps, shared package, Docker, configs)
- Initialize `go.mod` for the API
- Create `package.json` files for web, admin, and shared packages
- Create `pnpm-workspace.yaml` and `turbo.json`
- Create `.env` and `.env.example`
- Create `docker-compose.yml` (PostgreSQL 16, Redis 7, MinIO, Mailhog)
- Create `docker-compose.prod.yml` with multi-stage builds
- Create Dockerfiles for Go API and Next.js
- Create `.gitignore` files
- Print a colored success message with ASCII art and next steps

Support the `--api` flag to scaffold only the Go API (no frontend apps).

Print the Grit logo in ASCII art when the command runs.

#### Step 2: Go API (`apps/api/`)

Build the Go backend inside the scaffolded `apps/api/` directory:

**Configuration:**
- Load config from `.env` using `godotenv`
- Config struct with all settings (port, DB URL, JWT secret, Redis URL, S3 settings, etc.)

**Database:**
- GORM connection to PostgreSQL (configurable via DATABASE_URL)
- Auto-migration on startup
- Separate database.go with connection setup

**Models:**
- `User` model: ID, Name, Email, Password (bcrypt hashed), Role, Avatar, Active, EmailVerifiedAt, CreatedAt, UpdatedAt, DeletedAt (soft delete)
- `Role` as a string field with constants: "admin", "editor", "user"
- BeforeCreate hook for password hashing
- `CheckPassword(password string) bool` method

**Auth Handlers (`internal/handlers/auth.go`):**
- `POST /api/auth/register` — Validate (name, email, password), check email uniqueness, create user, return JWT tokens
- `POST /api/auth/login` — Validate email + password, return JWT tokens (access + refresh)
- `POST /api/auth/refresh` — Accept refresh token, return new access token
- `POST /api/auth/logout` — Invalidate refresh token
- `GET /api/auth/me` — Return current user (requires auth middleware)
- `POST /api/auth/forgot-password` — Accept email, generate reset token (log it for now, email in Phase 4)
- `POST /api/auth/reset-password` — Accept token + new password, reset it

**Auth Service (`internal/services/auth.go`):**
- JWT generation: access token (15min), refresh token (7 days)
- JWT validation and parsing
- Secret from `JWT_SECRET` env var

**Auth Middleware (`internal/middleware/auth.go`):**
- Extract Bearer token from Authorization header
- Validate JWT, extract user ID
- Load user from DB and attach to Gin context
- `RequireRole(roles ...string)` middleware for role-based access

**User Handlers (`internal/handlers/user.go`):**
- `GET /api/users` — List users with pagination, sorting, search (admin only)
- `GET /api/users/:id` — Get user by ID
- `PUT /api/users/:id` — Update user
- `DELETE /api/users/:id` — Soft delete user (admin only)

**Other Middleware:**
- CORS middleware (allow frontend origins from env)
- Request logging middleware (structured JSON)
- Recovery middleware (panic recovery)

**Routes:**
- Clean route registration in `internal/routes/routes.go`
- Public routes (auth)
- Protected routes (require JWT)
- Admin routes (require JWT + admin role)

**GORM Studio:**
- Import and mount GORM Studio at `/studio`
- Pass all models
- Protected with admin middleware

**Health Check:**
- `GET /api/health` — Returns `{ "status": "ok", "version": "0.1.0" }`

#### Step 3: Next.js Frontend (`apps/web/`)

Build the Next.js app inside `apps/web/`:

**Setup:**
- Next.js 14+ with App Router
- Tailwind CSS configured
- shadcn/ui initialized (install: button, input, card, label, toast, dropdown-menu, avatar, sheet, separator)
- Dark theme as default using CSS variables from CLAUDE.md design system
- Onest + JetBrains Mono fonts from Google Fonts

**API Client (`lib/api-client.ts`):**
- Axios instance with base URL from env (NEXT_PUBLIC_API_URL)
- Request interceptor: attach JWT access token from cookie/localStorage
- Response interceptor: on 401, attempt token refresh, retry original request
- Type-safe request/response helpers

**React Query Setup (`lib/query-client.ts`):**
- QueryClient with sensible defaults
- QueryClientProvider in root layout

**Auth Hooks (`hooks/use-auth.ts`):**
- `useLogin(email, password)` — mutation, stores tokens, redirects to dashboard
- `useRegister(name, email, password)` — mutation, auto-login after register
- `useLogout()` — mutation, clears tokens, redirects to login
- `useMe()` — query, fetches current user, used for auth state
- Token storage in cookies (httpOnly if possible) or localStorage

**Auth Pages:**
- `/login` — Dark themed login form with email + password. Links to register and forgot password. Shows validation errors. Beautiful design matching GORM Studio aesthetic.
- `/register` — Dark themed registration form with name, email, password, confirm password.
- `/forgot-password` — Email input form.
- All auth pages centered on screen with Grit logo above the form.

**Protected Layout (`(dashboard)/layout.tsx`):**
- Auth guard: redirect to `/login` if not authenticated
- Sidebar with navigation items:
  - Dashboard (icon: LayoutDashboard)
  - Users (icon: Users) — admin only
  - Settings (icon: Settings)
- Top navbar with:
  - Page title
  - User avatar + name dropdown
  - Logout option
- Collapsible sidebar on mobile (Sheet component)
- The sidebar should have the Grit logo at the top

**Dashboard Page:**
- Welcome message: "Welcome back, {name}"
- 4 stats cards (placeholder data): Total Users, Active Users, New This Month, Total Posts
- Stats cards should have icons, numbers, and subtle gradient backgrounds
- Below stats: "Getting Started" section with helpful links

**Responsive:** All pages must work on mobile.

#### Step 4: Admin Panel Shell (`apps/admin/`)

Build a minimal admin panel app in `apps/admin/`:

- Same Next.js + Tailwind + shadcn/ui setup as web app
- Same dark theme
- Admin layout with sidebar (matching the web app but with "Admin" badge)
- Dashboard page with placeholder widgets
- Users management page:
  - Data table showing users from `GET /api/users`
  - Columns: ID, Name, Email, Role, Active (badge), Created At
  - Pagination
  - Search
  - Edit and Delete actions
- This is the foundation that Phase 3 will expand into the full Filament-like system

#### Step 5: Shared Package (`packages/shared/`)

- `schemas/user.ts` — Zod schemas: LoginSchema, RegisterSchema, UpdateUserSchema
- `schemas/index.ts` — Re-exports all schemas
- `types/user.ts` — TypeScript interfaces: User, LoginRequest, RegisterRequest, AuthResponse
- `types/api.ts` — PaginatedResponse<T>, ApiError, ApiResponse<T>
- `types/index.ts` — Re-exports all types
- `constants/index.ts` — ROLES object, API_ROUTES object
- `package.json` with proper exports config

#### Step 6: Docker Setup

- `docker-compose.yml`:
  ```yaml
  services:
    postgres:    # PostgreSQL 16, port 5432, persistent volume
    redis:       # Redis 7, port 6379
    minio:       # MinIO, ports 9000/9001, default credentials
    mailhog:     # Mailhog, ports 1025/8025
  ```
- `docker-compose.prod.yml` with Go + Next.js containers
- `apps/api/Dockerfile` — Multi-stage Go build
- `apps/web/Dockerfile` — Next.js standalone build
- `apps/admin/Dockerfile` — Next.js standalone build

#### Step 7: Developer Experience

- `.env.example` with every variable documented with comments
- Root `README.md` with:
  - What is Grit
  - Quick start (3 commands)
  - Project structure overview
  - Links to docs
- `Makefile` or root `package.json` scripts for common tasks

---

### CRITICAL RULES

1. **Follow the folder structure in GRIT.md exactly.**
2. **Follow naming conventions in CLAUDE.md exactly.**
3. **Use ONLY the tech stack listed in CLAUDE.md.** Do not substitute libraries.
4. **Dark theme is mandatory.** Use the exact colors from CLAUDE.md.
5. **API responses must follow the format in CLAUDE.md** (data/meta/error structure).
6. **Test each piece as you build it.** Don't just write files — verify they work.
7. **Use conventional commits:** `feat:`, `fix:`, `docs:`, etc.
8. **Update CLAUDE.md "Current State" section** as you complete tasks.
9. **Check off completed items in PHASES.md** as you finish them.

### BUILD ORDER

Build in exactly this order:
1. CLI scaffolder first (so it generates the project structure)
2. Docker setup (so you can start PostgreSQL and Redis)
3. Go API with auth (core backend)
4. Shared package (types and schemas)
5. Next.js web app (frontend that talks to the API)
6. Admin panel shell (minimal admin UI)
7. GORM Studio integration
8. Developer experience polish (README, env files, scripts)

Start now. Begin by creating the CLI scaffolder.

## PROMPT END
