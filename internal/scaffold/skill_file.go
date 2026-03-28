package scaffold

import "fmt"

// gritSkillFile generates the .claude/skills/grit/SKILL.md for a scaffolded project.
// Tailored to the specific architecture and frontend the user chose.
func gritSkillFile(opts Options) string {
	bt := "`"

	// Architecture-specific description
	archDesc := ""
	switch opts.Architecture {
	case ArchSingle:
		archDesc = "a **single-app** architecture: one Go binary serves both the API and the embedded React SPA (via go:embed). Think Laravel or embedded Next.js — one deployment, one URL, one binary."
	case ArchDouble:
		archDesc = "a **double** architecture: Go API + React web frontend in a Turborepo monorepo. Two apps that share types and schemas."
	case ArchTriple:
		archDesc = "a **triple** architecture: Go API + React web frontend + admin panel in a Turborepo monorepo. Three apps that share types and schemas."
	case ArchAPI:
		archDesc = "an **API-only** architecture: pure Go backend with no frontend. Perfect for headless APIs, mobile backends, and microservices."
	case ArchMobile:
		archDesc = "a **mobile** architecture: Go API + Expo React Native app in a Turborepo monorepo. Two apps that share types and schemas."
	default:
		archDesc = "a **triple** architecture: Go API + React web frontend + admin panel in a Turborepo monorepo."
	}

	// Frontend-specific description
	frontendDesc := ""
	frontendRouting := ""
	switch opts.Frontend {
	case FrontendTanStack:
		frontendDesc = "**TanStack Router** (Vite) for the frontend — SPA with file-based routing, fast HMR, smaller bundles."
		frontendRouting = `Routes live in src/routes/. File naming: index.tsx (list), new.tsx (create), $id.tsx (detail), $id.edit.tsx (edit).
Use Route.useParams() for type-safe params. Uses createHashHistory() for Wails desktop compatibility.
Navigation: import { Link, useNavigate } from '@tanstack/react-router'.`
	default:
		frontendDesc = "**Next.js** (App Router) for the frontend — SSR, SSG, ISR, SEO-optimized."
		frontendRouting = `Routes live in app/. File naming: page.tsx (route), layout.tsx (layout), loading.tsx, error.tsx.
Use 'use client' directive for client components. Server Components by default.
Navigation: import { useRouter } from 'next/navigation', import Link from 'next/link'.`
	}

	// Project structure based on architecture
	projectStructure := ""
	switch opts.Architecture {
	case ArchSingle:
		projectStructure = fmt.Sprintf(`%[1]s/
├── main.go                       # Entry point (go:embed frontend/dist)
├── internal/
│   ├── config/                   # Loads .env
│   ├── database/                 # GORM connection
│   ├── models/                   # GORM models
│   ├── handlers/                 # HTTP handlers
│   ├── services/                 # Business logic
│   ├── middleware/               # Auth, CORS, logger, cache, maintenance
│   ├── routes/                   # Route registration
│   ├── mail/                     # Email (Resend)
│   ├── storage/                  # File storage (S3)
│   ├── jobs/                     # Background jobs (asynq)
│   ├── cache/                    # Redis cache
│   └── ai/                       # AI service (Vercel AI Gateway)
├── frontend/                     # React app (Vite + TanStack Router)
│   ├── src/routes/               # File-based routes
│   ├── src/components/           # React components
│   ├── src/hooks/                # Custom hooks
│   └── src/lib/                  # API client, utilities
├── .env                          # Environment variables
└── docker-compose.yml            # PostgreSQL, Redis, MinIO, Mailhog`, opts.ProjectName)
	case ArchAPI:
		projectStructure = fmt.Sprintf(`%[1]s/
├── apps/api/                     # Go backend (Gin + GORM)
│   ├── cmd/server/main.go
│   └── internal/                 # config, database, models, handlers, services, middleware, routes
├── .env
└── docker-compose.yml`, opts.ProjectName)
	case ArchDouble:
		webDir := "web/"
		if opts.Frontend == FrontendTanStack {
			webDir = "web/                      # React app (Vite + TanStack Router)"
		} else {
			webDir = "web/                      # Next.js frontend (App Router)"
		}
		projectStructure = fmt.Sprintf(`%[1]s/
├── packages/shared/              # Zod schemas, TS types, constants
├── apps/
│   ├── api/                      # Go backend (Gin + GORM)
│   │   ├── cmd/server/main.go
│   │   └── internal/             # config, database, models, handlers, services, middleware, routes
│   └── %[2]s
├── .env
├── docker-compose.yml
└── turbo.json                    # Turborepo configuration`, opts.ProjectName, webDir)
	case ArchMobile:
		projectStructure = fmt.Sprintf(`%[1]s/
├── packages/shared/              # Zod schemas, TS types, constants
├── apps/
│   ├── api/                      # Go backend (Gin + GORM)
│   │   ├── cmd/server/main.go
│   │   └── internal/             # config, database, models, handlers, services, middleware, routes
│   └── expo/                     # Expo React Native app
│       ├── app/                  # Expo Router screens
│       ├── components/           # React Native components
│       └── hooks/                # Custom hooks
├── .env
├── docker-compose.yml
└── turbo.json`, opts.ProjectName)
	default: // Triple
		adminDir := "admin/"
		webDir := "web/"
		if opts.Frontend == FrontendTanStack {
			webDir = "web/                      # React app (Vite + TanStack Router)"
			adminDir = "admin/                    # Admin panel (Vite + TanStack Router)"
		} else {
			webDir = "web/                      # Next.js frontend (App Router)"
			adminDir = "admin/                    # Next.js admin panel (App Router)"
		}
		projectStructure = fmt.Sprintf(`%[1]s/
├── packages/shared/              # Zod schemas, TS types, constants
├── apps/
│   ├── api/                      # Go backend (Gin + GORM)
│   │   ├── cmd/server/main.go
│   │   └── internal/             # config, database, models, handlers, services, middleware, routes
│   ├── %[2]s
│   └── %[3]s
├── .env
├── docker-compose.yml
└── turbo.json`, opts.ProjectName, webDir, adminDir)
	}

	// Architecture-specific common tasks
	addFieldPaths := ""
	switch opts.Architecture {
	case ArchSingle:
		addFieldPaths = `1. Add field to Go model (internal/models/<name>.go)
2. Update handler if field needs special handling
3. Update Zod schema (frontend/src/schemas/<name>.ts)
4. Update TypeScript type (frontend/src/types/<name>.ts)
5. Restart API (GORM auto-migrates)`
	case ArchAPI:
		addFieldPaths = `1. Add field to Go model (apps/api/internal/models/<name>.go)
2. Update handler if field needs special handling
3. Restart API (GORM auto-migrates)`
	default:
		addFieldPaths = `1. Add field to Go model (apps/api/internal/models/<name>.go)
2. Update handler if field needs special handling
3. Update Zod schema (packages/shared/schemas/<name>.ts)
4. Update TypeScript type (packages/shared/types/<name>.ts)
5. Update admin resource (apps/admin/resources/<name>.ts) — add column + form field
6. Restart API (GORM auto-migrates)`
	}

	// Architecture-specific critical rules
	frontendRules := ""
	switch opts.Architecture {
	case ArchAPI:
		frontendRules = "" // No frontend rules for API-only
	case ArchSingle:
		if opts.Frontend == FrontendTanStack {
			frontendRules = `5. **Use React Query** for all data fetching — no raw fetch
6. **Use Zod** for validation
7. **Use Tailwind + shadcn/ui** — no custom CSS files
8. **Use TanStack Router** — file-based routes in src/routes/`
		} else {
			frontendRules = `5. **Use React Query** for all data fetching — no raw fetch
6. **Use Zod** for validation
7. **Use Tailwind + shadcn/ui** — no custom CSS files
8. **Use App Router** — never Pages Router`
		}
	default:
		if opts.Frontend == FrontendTanStack {
			frontendRules = `5. **Use React Query** for all data fetching — no raw fetch
6. **Use Zod** for validation — shared between frontend and backend
7. **Use Tailwind + shadcn/ui** — no custom CSS files
8. **Use TanStack Router** — file-based routes in src/routes/`
		} else {
			frontendRules = `5. **Use React Query** for all data fetching — no raw fetch
6. **Use Zod** for validation — shared between frontend and backend
7. **Use Tailwind + shadcn/ui** — no custom CSS files
8. **Use App Router** — never Pages Router`
		}
	}

	// Build the marker comments section (only for architectures that have them)
	markerSection := ""
	if opts.Architecture != ArchAPI {
		tsMarkers := ""
		if opts.Architecture != ArchSingle || opts.Frontend != "" {
			tsMarkers = fmt.Sprintf(`
%[1]stypescript
// grit:schemas         — schemas/index.ts
// grit:types           — types/index.ts
// grit:api-routes      — constants/index.ts
%[1]s`, bt)
			if opts.ShouldIncludeAdmin() {
				tsMarkers += fmt.Sprintf(`
%[1]stypescript
// grit:resources       — resources/index.ts (imports)
// grit:resource-list   — resources/index.ts (registry array)
%[1]s`, bt)
			}
		}

		markerSection = fmt.Sprintf(`
## Marker Comments

Grit uses marker comments to inject generated code. **Never delete these:**

%[1]sgo
// grit:models          — models/user.go (AutoMigrate list)
// grit:handlers        — routes/routes.go (handler initialization)
// grit:routes:protected — routes/routes.go (protected route group)
// grit:routes:admin    — routes/routes.go (admin route group)
%[1]s
%[2]s
---
`, bt, tsMarkers)
	} else {
		markerSection = fmt.Sprintf(`
## Marker Comments

Grit uses marker comments to inject generated code. **Never delete these:**

%[1]sgo
// grit:models          — models/user.go (AutoMigrate list)
// grit:handlers        — routes/routes.go (handler initialization)
// grit:routes:protected — routes/routes.go (protected route group)
// grit:routes:admin    — routes/routes.go (admin route group)
%[1]s

---
`, bt)
	}

	// Single app specific notes
	singleAppNote := ""
	if opts.Architecture == ArchSingle {
		singleAppNote = fmt.Sprintf(`
## Single App — go:embed

This project uses %[1]sgo:embed%[1]s to embed the frontend build output into the Go binary:

%[1]sgo
//go:embed frontend/dist/*
var frontendFS embed.FS
%[1]s

In **development**, run Go and Vite separately — Vite proxies API calls to Go:
- Go API: %[1]sgo run main.go%[1]s (port 8080)
- Vite dev: %[1]scd frontend && pnpm dev%[1]s (port 5173, proxies /api → 8080)

For **production**, build the frontend first, then the Go binary:
%[1]sash
cd frontend && pnpm build    # Outputs to frontend/dist/
go build -o myapp main.go   # Embeds frontend/dist/ into binary
./myapp                      # Serves API + frontend on port 8080
%[1]s

---
`, bt)
	}

	// Mobile-specific notes
	mobileNote := ""
	if opts.Architecture == ArchMobile {
		mobileNote = fmt.Sprintf(`
## Mobile — Expo React Native

The Expo app lives in %[1]sapps/expo/%[1]s and uses Expo Router for file-based navigation.

**Development:**
- Start API: %[1]scd apps/api && go run cmd/server/main.go%[1]s
- Start Expo: %[1]scd apps/expo && npx expo start%[1]s
- Scan QR code with Expo Go on your phone

**Token storage:** Use %[1]sexpo-secure-store%[1]s (encrypted) — NOT AsyncStorage.

**API URL:** On physical devices, %[1]slocalhost%[1]s won't work — use your machine's local IP.

---
`, bt)
	}

	return fmt.Sprintf(`---
name: grit
description: >
  Grit framework conventions and patterns for this %[3]s project.
  Use when modifying models, handlers, routes, schemas, types, or components.
  Automatically loaded as background knowledge.
user-invocable: false
---

# Grit Framework — %[4]s

This project uses %[5]s

It uses %[6]s

**Batteries included:** file storage (S3), email (Resend), background jobs (asynq), cron, Redis cache, AI (Vercel AI Gateway), security (Sentinel), observability (Pulse), auto-generated API docs (gin-docs).

For detailed API conventions, code patterns, and service documentation, see [reference.md](reference.md).

---

## CLI Commands

%[1]sash
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
%[1]s

---

## Project Structure

%[1]s
%[7]s
%[1]s

**Mounted dashboards** (auto-configured in routes.go):
- %[1]s/docs%[1]s — API documentation (gin-docs, OpenAPI 3.1)
- %[1]s/studio%[1]s — Database browser (GORM Studio)
- %[1]s/sentinel/ui%[1]s — Security dashboard (WAF, rate limiting)
- %[1]s/pulse%[1]s — Observability (tracing, metrics)

---

## Generating Resources

%[1]sash
grit generate resource Post --fields "title:string,content:text,published:bool,views:int"
%[1]s

Creates model, service, handler, schema, types, hooks, and injects into existing files via marker comments.

### Field Types

| Type | Go | TypeScript | Form |
|------|----|-----------|------|
| %[1]sstring%[1]s | %[1]sstring%[1]s | %[1]sstring%[1]s | Text input |
| %[1]stext%[1]s | %[1]sstring%[1]s | %[1]sstring%[1]s | Textarea |
| %[1]sint%[1]s / %[1]suint%[1]s / %[1]sfloat%[1]s | %[1]sint%[1]s / %[1]suint%[1]s / %[1]sfloat64%[1]s | %[1]snumber%[1]s | Number input |
| %[1]sbool%[1]s | %[1]sbool%[1]s | %[1]sboolean%[1]s | Toggle |
| %[1]sdatetime%[1]s / %[1]sdate%[1]s | %[1]s*time.Time%[1]s | %[1]sstring | null%[1]s | Picker |
| %[1]srichtext%[1]s | %[1]sstring%[1]s | %[1]sstring%[1]s | Tiptap editor |
| %[1]sslug%[1]s | %[1]sstring%[1]s | %[1]sstring%[1]s | Auto-generated |
| %[1]sstring_array%[1]s | %[1]sJSONSlice[string]%[1]s | %[1]sstring[]%[1]s | Tag input |
| %[1]sbelongs_to:X%[1]s | %[1]suint%[1]s (FK) | %[1]snumber%[1]s | Relationship select |
| %[1]smany_to_many:X%[1]s | Junction table | %[1]snumber[]%[1]s | Multi-select |

**Modifiers:** %[1]s:unique%[1]s, %[1]s:required%[1]s, %[1]s:optional%[1]s (append after type).

---
%[8]s%[9]s%[10]s
## Frontend Routing

%[11]s

---

## Common Tasks

### Add a field to an existing resource

%[12]s

### Add a new API endpoint

1. Create/update handler in the handlers directory
2. Register route in routes.go
3. Create React Query hook if frontend needs it

### Add a relationship

%[1]sgo
type Post struct {
    CategoryID uint     // json:"category_id"
    Category   Category // gorm:"foreignKey:CategoryID" json:"category,omitempty"
}
// In handler: query.Preload("Category").Find(&posts)
%[1]s

---

## Critical Rules

1. **Never delete marker comments** (%[1]s// grit:*%[1]s)
2. **Follow the response format** — %[1]s{ data, message }%[1]s / %[1]s{ data, meta }%[1]s / %[1]s{ error: { code, message } }%[1]s
3. **Always handle errors in Go** — never ignore with %[1]s_%[1]s
4. **Keep the folder structure** — don't move files
%[13]s
9. **Report bugs** — if you encounter a bug, open an issue at https://github.com/MUKE-coder/grit/issues with the error message and steps to reproduce
`,
		bt,                  // [1] backtick
		opts.ProjectName,    // [2] project name (unused in this format but kept for consistency)
		opts.Architecture,   // [3] architecture name
		archTitle(opts),     // [4] human-readable title
		archDesc,            // [5] architecture description
		frontendDesc,        // [6] frontend description
		projectStructure,    // [7] project tree
		markerSection,       // [8] marker comments
		singleAppNote,       // [9] single app note (empty if not single)
		mobileNote,          // [10] mobile note (empty if not mobile)
		frontendRouting,     // [11] routing info
		addFieldPaths,       // [12] add field steps
		frontendRules,       // [13] frontend-specific rules
	)
}

func archTitle(opts Options) string {
	arch := string(opts.Architecture)
	fe := ""
	switch opts.Frontend {
	case FrontendTanStack:
		fe = " + TanStack Router (Vite)"
	case FrontendNext:
		fe = " + Next.js"
	}

	switch opts.Architecture {
	case ArchSingle:
		return "Single App" + fe
	case ArchDouble:
		return "Double (Web + API)" + fe
	case ArchTriple:
		return "Triple (Web + Admin + API)" + fe
	case ArchAPI:
		return "API Only"
	case ArchMobile:
		return "Mobile (API + Expo)"
	default:
		return arch + fe
	}
}

// gritSkillReference generates .claude/skills/grit/reference.md with detailed
// API conventions, code patterns, admin panel docs, and service documentation.
// This file is the same across architectures (API conventions don't change).
func gritSkillReference(opts Options) string {
	bt := "`"

	// Only include admin panel section for triple architecture
	adminSection := ""
	if opts.ShouldIncludeAdmin() {
		adminSection = fmt.Sprintf(`
## Admin Panel — Resource Definitions

%[1]stypescript
import { defineResource } from "@/lib/resource";

export const postsResource = defineResource({
  name: "Post",
  slug: "posts",
  endpoint: "/api/posts",
  icon: "FileText",
  label: { singular: "Post", plural: "Posts" },

  table: {
    columns: [
      { key: "id", label: "ID", sortable: true, format: "number" },
      { key: "title", label: "Title", sortable: true, searchable: true },
      { key: "published", label: "Published", format: "boolean" },
      { key: "created_at", label: "Created", sortable: true, format: "relative" },
    ],
    defaultSort: { key: "created_at", direction: "desc" },
    filters: [{ key: "published", label: "Published", type: "boolean" }],
    pageSize: 20,
    searchable: true,
    actions: ["create", "edit", "delete"],
  },

  form: {
    fields: [
      { key: "title", label: "Title", type: "text", required: true },
      { key: "content", label: "Content", type: "textarea", required: true },
      { key: "published", label: "Published", type: "toggle", defaultValue: false },
      { key: "cover", label: "Cover Image", type: "image" },
    ],
    layout: "single",
  },
});
%[1]s

### Form Field Types

| Type | Component | Notes |
|------|----------|-------|
| %[1]stext%[1]s | Text input | prefix, suffix |
| %[1]stextarea%[1]s | Textarea | configurable rows |
| %[1]snumber%[1]s | Number input | min, max, step |
| %[1]sselect%[1]s | Dropdown | requires options |
| %[1]sdate%[1]s / %[1]sdatetime%[1]s | Picker | |
| %[1]stoggle%[1]s / %[1]scheckbox%[1]s | Boolean | |
| %[1]simage%[1]s | Image upload | react-dropzone |
| %[1]srichtext%[1]s | Tiptap WYSIWYG | |
| %[1]srelationship-select%[1]s | Searchable dropdown | belongs_to |
| %[1]smulti-relationship-select%[1]s | Multi-select tags | many_to_many |

### Form View Variants

%[1]sformView%[1]s: %[1]smodal%[1]s (default), %[1]spage%[1]s, %[1]smodal-steps%[1]s, %[1]spage-steps%[1]s

---
`, bt)
	}

	// Frontend patterns section (skip for API-only)
	frontendSection := ""
	if opts.Architecture != ArchAPI {
		frontendSection = fmt.Sprintf(`
## Frontend Patterns

### React Query Hooks

%[1]stypescript
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { apiClient } from "@/lib/api-client";

export function usePosts({ page = 1, pageSize = 20, search = "" } = {}) {
  return useQuery({
    queryKey: ["posts", { page, pageSize, search }],
    queryFn: async () => {
      const params = new URLSearchParams({
        page: String(page),
        page_size: String(pageSize),
        ...(search && { search }),
      });
      const { data } = await apiClient.get("/api/posts?" + params);
      return data;
    },
  });
}
%[1]s

---
`, bt)
	}

	return fmt.Sprintf(`# Grit Framework — Detailed Reference

## API Conventions

### Response Format

%[1]sgo
// Success (single item)
c.JSON(http.StatusOK, gin.H{
    "data":    item,
    "message": "Item retrieved successfully",
})

// Success (paginated list)
c.JSON(http.StatusOK, gin.H{
    "data": items,
    "meta": gin.H{
        "total":     total,
        "page":      page,
        "page_size": pageSize,
        "pages":     pages,
    },
})

// Error
c.JSON(http.StatusBadRequest, gin.H{
    "error": gin.H{
        "code":    "VALIDATION_ERROR",
        "message": "Email is required",
    },
})
%[1]s

### Error Codes

| Code | HTTP Status | When |
|------|------------|------|
| %[1]sVALIDATION_ERROR%[1]s | 422 | Invalid input |
| %[1]sNOT_FOUND%[1]s | 404 | Resource missing |
| %[1]sUNAUTHORIZED%[1]s | 401 | Invalid JWT |
| %[1]sFORBIDDEN%[1]s | 403 | Insufficient role |
| %[1]sINTERNAL_ERROR%[1]s | 500 | Server error |
| %[1]sCONFLICT%[1]s | 409 | Duplicate key |

### Authentication

%[1]s
POST /api/auth/register  → { access_token, refresh_token }
POST /api/auth/login     → { access_token, refresh_token } or { totp_required, pending_token }
POST /api/auth/refresh   → New access_token from refresh_token
POST /api/auth/logout    → Invalidates refresh token
GET  /api/auth/me        → Current user (requires auth)
%[1]s

Access tokens: 15 minutes. Refresh tokens: 7 days.

### Two-Factor Authentication (TOTP)

If user has 2FA enabled and no trusted device cookie, login returns %[1]s{ totp_required: true, pending_token: "..." }%[1]s.
Client redirects to TOTP page, user enters 6-digit code from authenticator app.

%[1]s
POST /api/auth/totp/setup              → { secret, uri } (JWT required)
POST /api/auth/totp/enable             → { enabled, backup_codes } (JWT required)
POST /api/auth/totp/verify             → { user, tokens } (public, uses pending_token)
POST /api/auth/totp/backup-codes/verify → { user, tokens } (public, uses pending_token)
POST /api/auth/totp/disable            → Disable 2FA (JWT + password required)
GET  /api/auth/totp/status             → { enabled, backup_codes_remaining, trusted_devices }
%[1]s

TOTP: RFC 6238, HMAC-SHA1, 6 digits, 30s period. Backup codes: 10 bcrypt-hashed one-time codes.
Trusted devices: HttpOnly cookie, SHA-256 hashed token, 30-day sliding expiry.

### Route Groups

%[1]sgo
public := router.Group("/api/auth")          // No auth
protected := router.Group("/api")            // Requires JWT
protected.Use(middleware.Auth(cfg.JWTSecret))
admin := protected.Group("/admin")           // Requires JWT + admin role
admin.Use(middleware.RequireRole("admin"))
%[1]s

---

## Go Model Pattern

%[1]sgo
type Post struct {
    ID        uint           // gorm:"primarykey" json:"id"
    Title     string         // gorm:"size:255;not null" json:"title" binding:"required"
    Slug      string         // gorm:"size:255;uniqueIndex" json:"slug"
    Content   string         // gorm:"type:text" json:"content"
    Published bool           // gorm:"default:false" json:"published"
    UserID    uint           // json:"user_id"
    User      User           // gorm:"foreignKey:UserID" json:"user,omitempty"
    CreatedAt time.Time      // json:"created_at"
    UpdatedAt time.Time      // json:"updated_at"
    DeletedAt gorm.DeletedAt // gorm:"index" json:"-"
}
%[1]s

Rules:
- Always include ID, CreatedAt, UpdatedAt, DeletedAt
- %[1]sjson:"-"%[1]s for DeletedAt (hidden from API)
- %[1]sbinding:"required"%[1]s for required fields

---
%[2]s%[3]s
## Batteries (Services)

### File Storage

%[1]sgo
storage.Upload(ctx, "uploads/2024/01/photo.jpg", reader, "image/jpeg")
url := storage.GetURL("uploads/2024/01/photo.jpg")
url, err := storage.GetSignedURL(ctx, key, 1*time.Hour)
%[1]s

### Email

%[1]sgo
mailer.Send(ctx, mail.SendOptions{
    To: "user@example.com", Subject: "Welcome!",
    Template: "welcome", Data: map[string]interface{}{"Name": "John"},
})
%[1]s

Templates: %[1]swelcome%[1]s, %[1]spassword-reset%[1]s, %[1]semail-verification%[1]s, %[1]snotification%[1]s

### Background Jobs

%[1]sgo
jobs.EnqueueSendEmail("user@example.com", "Welcome", "welcome", data)
jobs.EnqueueProcessImage(uploadID, key, mimeType)
%[1]s

### Redis Cache

%[1]sgo
cache.Set(ctx, "user:123", userData, 5*time.Minute)
cache.Get(ctx, "user:123", &user)
cache.Delete(ctx, "user:123")
%[1]s

### AI Integration (Vercel AI Gateway)

%[1]sgo
result, err := ai.Complete(ctx, ai.CompletionRequest{Prompt: "Summarize..."})
ai.Stream(ctx, req, func(chunk string) { /* SSE */ })
%[1]s

One key, hundreds of models. Config: %[1]sAI_GATEWAY_API_KEY%[1]s, %[1]sAI_GATEWAY_MODEL%[1]s (e.g. %[1]santhropic/claude-sonnet-4-6%[1]s).

### Security (Sentinel)

WAF, rate limiting, brute-force protection. Dashboard at %[1]s/sentinel/ui%[1]s.

### Observability (Pulse)

Request tracing, DB monitoring, metrics. Dashboard at %[1]s/pulse%[1]s. Prometheus at %[1]s/pulse/metrics%[1]s.

### API Documentation (gin-docs)

Zero-annotation OpenAPI 3.1 spec. Interactive UI at %[1]s/docs%[1]s.

---

## Naming Conventions

| Thing | Convention | Example |
|-------|-----------|---------|
| Go files | %[1]ssnake_case.go%[1]s | %[1]suser_handler.go%[1]s |
| Go structs | %[1]sPascalCase%[1]s | %[1]stype PostHandler struct%[1]s |
| TS files | %[1]skebab-case.ts%[1]s | %[1]suse-posts.ts%[1]s |
| React | %[1]sPascalCase.tsx%[1]s | %[1]sDataTable.tsx%[1]s |
| API routes | %[1]s/api/plural%[1]s | %[1]s/api/posts%[1]s |
| DB tables | %[1]splural_snake%[1]s | %[1]sblog_posts%[1]s |
`, bt, adminSection, frontendSection)
}
