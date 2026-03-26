# Grit v3.0 — Multi-Architecture Refactor Plan

## Summary

Transform Grit from a fixed triple-architecture scaffolder into a flexible multi-architecture, multi-frontend CLI with 5 architecture modes and 2 frontend frameworks.

## Architecture Matrix

| Mode | Structure | Frontend Options | Deploy |
|------|-----------|-----------------|--------|
| **Single** | Go API + embedded React SPA (go:embed) | TanStack Router / Next.js | Single binary |
| **Double** | Turborepo (apps/web + apps/api) | TanStack Router / Next.js | 2 containers |
| **Triple** | Turborepo (apps/web + apps/admin + apps/api) | TanStack Router / Next.js | 3 containers |
| **API Only** | Go API only | None | Single binary |
| **API + Expo** | Turborepo (apps/expo + apps/api) | React Native | 2 processes |

## Interactive CLI Flow

```
$ grit new my-app

? Select architecture:
  > Single App (Go API + embedded React SPA — one binary)
    Double (Web + API monorepo)
    Triple (Web + Admin + API monorepo)
    API Only
    API + Mobile (Expo)

? Select frontend framework:
  > TanStack Router (Vite — fast builds, small bundle)
    Next.js (SSR, SEO, App Router)
```

Flags for non-interactive: `grit new my-app --single --vite`, `grit new my-app --triple --next`

---

## Phase 1: Interactive CLI + Options Refactor

**Complexity: Medium | Est: 1-2 days**

### Changes
- Replace boolean flags (`APIOnly`, `IncludeExpo`, `MobileOnly`, `Full`) with `Architecture` and `Frontend` enums
- Add interactive prompts using `github.com/charmbracelet/huh`
- Keep old flags as aliases for backward compatibility
- New flags: `--single`, `--double`, `--triple`, `--vite`, `--next`

### New Options Struct
```go
type Architecture string
const (
    ArchSingle Architecture = "single"
    ArchDouble Architecture = "double"
    ArchTriple Architecture = "triple"
    ArchAPI    Architecture = "api"
    ArchMobile Architecture = "mobile"
)

type Frontend string
const (
    FrontendNext     Frontend = "next"
    FrontendTanStack Frontend = "tanstack"
)
```

### ShouldInclude* Mapping

| Method | single | double | triple | api | mobile |
|--------|--------|--------|--------|-----|--------|
| ShouldIncludeWeb() | false | true | true | false | false |
| ShouldIncludeAdmin() | false | false | true | false | false |
| ShouldIncludeFrontend() | true | true | true | false | false |
| ShouldIncludeSingleSPA() | true | false | false | false | false |
| ShouldUseTurborepo() | false | true | true | false | true |
| ShouldIncludeShared() | false | true | true | false | true |
| ShouldIncludeExpo() | false | false | false | false | true |

### Files
- [x] `go.mod` — add `github.com/charmbracelet/huh`
- [ ] `internal/scaffold/scaffold.go` — refactor Options, update ShouldInclude*, update Run()
- [ ] `internal/prompt/prompt.go` — new package for interactive prompts
- [ ] `cmd/grit/main.go` — interactive flow, new flags, backward compat
- [ ] Tests — update all test construction of Options

---

## Phase 2: TanStack Router Frontend

**Complexity: HIGHEST | Est: 5-7 days**

### Changes
Create parallel TanStack Router templates for every Next.js frontend file.

### Key Differences

| Aspect | Next.js | TanStack Router |
|--------|---------|-----------------|
| Routing | app/ directory | src/routes/ via plugin |
| Layouts | layout.tsx | __root.tsx + _layout.tsx |
| Build | next build | vite build |
| Dev | next dev | vite dev |
| "use client" | Required | Not needed |
| Output | .next/ | dist/ |
| SSR | Built-in | SPA only |

### TanStack Web Structure
```
apps/web/
├── src/
│   ├── routes/
│   │   ├── __root.tsx
│   │   ├── index.tsx
│   │   └── blog/
│   ├── components/
│   ├── hooks/
│   ├── lib/
│   └── main.tsx
├── index.html
├── vite.config.ts
├── tailwind.config.ts
└── package.json
```

### Files
- [ ] `internal/scaffold/web_tanstack_files.go` — TanStack web app templates
- [ ] `internal/scaffold/admin_tanstack_files.go` — TanStack admin templates
- [ ] `internal/scaffold/admin_tanstack_layout_files.go` — sidebar, navbar
- [ ] `internal/scaffold/admin_tanstack_table_files.go` — DataTable (no "use client")
- [ ] `internal/scaffold/admin_tanstack_form_files.go` — FormBuilder
- [ ] `internal/scaffold/admin_tanstack_widget_files.go` — dashboard widgets
- [ ] `internal/scaffold/admin_tanstack_system_files.go` — system pages
- [ ] `internal/scaffold/admin_tanstack_resource_files.go` — resource definitions
- [ ] `internal/scaffold/admin_tanstack_style_files.go` — style variants
- [ ] `internal/scaffold/shared_tanstack_files.go` — shared types (Vite-compatible)
- [ ] `internal/scaffold/scaffold.go` — dispatcher for frontend framework
- [ ] Tests — verify all TanStack files generated correctly

---

## Phase 3: Single App Architecture

**Complexity: HIGH | Est: 3-4 days**

### Single App Structure
```
my-app/
├── main.go                    # go:embed + Gin server
├── internal/                  # Go backend (same structure as api/)
│   ├── config/
│   ├── database/
│   ├── models/
│   ├── handlers/
│   ├── services/
│   ├── middleware/
│   └── routes/
├── frontend/                  # React SPA (TanStack Router or Next.js static)
│   ├── src/
│   ├── package.json
│   └── vite.config.ts
├── go.mod                     # Module: my-app (not my-app/apps/api)
├── .env
├── docker-compose.yml
└── Makefile
```

### Embed Pattern
```go
//go:embed frontend/dist/*
var frontendFS embed.FS

// Serve API at /api/*, SPA fallback for everything else
```

### Dev Mode
- Go API on :8080
- Vite on :5173 (standard dev server)
- Go proxies non-API requests to Vite in dev mode

### Files
- [ ] `internal/scaffold/single_scaffold.go` — RunSingle() orchestrator
- [ ] `internal/scaffold/single_root_files.go` — main.go with embed, go.mod, .env, Makefile
- [ ] `internal/scaffold/single_go_files.go` — Go backend (reuse api_files.go templates)
- [ ] `internal/scaffold/single_frontend_files.go` — React SPA scaffold
- [ ] `internal/scaffold/single_test.go` — tests
- [ ] `internal/project/detect.go` — add ProjectSingle detection

### Key Challenge
API files currently use module path `my-app/apps/api`. For single app it's just `my-app`. Must parameterize the module path (already uses `{{MODULE}}` placeholder).

---

## Phase 4: Double Architecture

**Complexity: LOW | Est: 0.5 days**

Turborepo with apps/web + apps/api (no admin). Essentially current triple minus admin.

### Files
- [ ] `internal/scaffold/root_files.go` — conditionalize grit.config.ts, package.json
- [ ] `internal/scaffold/docker_files.go` — remove admin service when double
- [ ] `cmd/grit/main.go` — update printSuccess()

---

## Phase 5: Goravel-Inspired CLI Commands

**Complexity: Medium | Est: 2-3 days**

### 5a. `grit routes` — List registered API routes
Parse routes.go and print table of method, path, handler, middleware.

### 5b. `grit down` / `grit up` — Maintenance mode
Write/delete `.maintenance` file. Add maintenance middleware to scaffolded API.

### 5c. `grit deploy` — Production deployment
1. Build Go binary + frontend
2. SSH to target server
3. Upload binary + assets
4. Configure systemd service
5. Configure Caddy reverse proxy with auto-TLS

### 5d. Service provider pattern
Refactor how batteries are initialized in scaffolded main.go. Each battery becomes a registerable provider.

### Files
- [ ] `internal/routeparser/parser.go` — parse Gin route registrations
- [ ] `internal/deploy/deploy.go` — deployment orchestrator
- [ ] `internal/deploy/ssh.go` — SSH client
- [ ] `internal/deploy/systemd.go` — systemd unit template
- [ ] `internal/deploy/caddy.go` — Caddyfile template
- [ ] `cmd/grit/main.go` — routesCmd(), downCmd(), upCmd(), deployCmd()
- [ ] `internal/scaffold/api_files.go` — add maintenance middleware

---

## Phase 6: Update Code Generator

**Complexity: HIGH | Est: 3-4 days**

### grit.json Project Manifest
Written at scaffold time, read by generator:
```json
{
  "architecture": "triple",
  "frontend": "next",
  "version": "3.0.0"
}
```

### Generator Dispatch
```go
switch info.Architecture {
case "single":
    g.writeSingleFrontendResource(info)
case "double":
    g.writeWebResource(info)
case "triple":
    g.writeWebResource(info)
    g.writeAdminResource(info)
}
```

### Files
- [ ] `internal/project/detect.go` — read grit.json, detect architecture+frontend
- [ ] `internal/generate/generator.go` — dispatch based on architecture+frontend
- [ ] `internal/generate/templates_tanstack.go` — TanStack resource templates
- [ ] `internal/generate/single.go` — single-app resource generation
- [ ] All scaffold files — write grit.json at scaffold time

---

## Implementation Order

```
Phase 1 (Interactive CLI) ──> Phase 4 (Double) ──> Phase 2 (TanStack) ──> Phase 3 (Single) ──> Phase 6 (Generator)
                                                         |
                                                   Phase 5 (CLI Commands) — can run in parallel
```

**Total estimate: 15-21 working days**

## Version Plan
- Phase 1 + 4: v3.0.0 (breaking: new Options, interactive CLI)
- Phase 2: v3.1.0 (TanStack Router option)
- Phase 3: v3.2.0 (Single app architecture)
- Phase 5: v3.3.0 (deploy, routes, maintenance mode)
- Phase 6: v3.4.0 (generator for all architectures)
