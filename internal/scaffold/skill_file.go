package scaffold

import "fmt"

func gritSkillFile(opts Options) string {
	return fmt.Sprintf(`# Grit Framework — LLM Skill Guide

> **This document teaches AI assistants (Claude, Cursor, Kilo Code, etc.) how to work with the Grit framework.** Read this file completely before writing any code in a Grit project.

---

## What is Grit?

Grit is a full-stack meta-framework that combines **Go** (backend) + **React/Next.js** (frontend) in a monorepo. It provides:

- A **CLI tool** (%[1]s) that scaffolds entire projects and generates full-stack resources
- A **Go API** with Gin + GORM + PostgreSQL
- A **Next.js web app** with App Router + Tailwind + shadcn/ui
- A **Filament-like admin panel** with resource definitions, DataTables, forms, and widgets
- **Batteries included**: file storage (S3), email (Resend), background jobs (asynq), cron, Redis caching, AI integration (Claude/OpenAI)
- A **shared package** with Zod schemas, TypeScript types, and constants

**Think of it as:** Laravel + Filament, but with Go + React instead of PHP + Blade.

---

## Quick Reference — CLI Commands

%[1]sash
# Create a new project
grit new myapp                    # Full monorepo (Go API + Next.js web + admin)
grit new myapp --api              # Go API only
grit new myapp --full             # Everything + Expo mobile + docs site

# Generate a resource (full-stack CRUD)
grit generate resource Post --fields "title:string,content:text,published:bool"
grit generate resource Post --fields "title:string,slug:string:unique,views:int"
grit generate resource Post --from post.yaml
grit generate resource Category -i   # Interactive mode

# Sync Go types to TypeScript
grit sync

# Upgrade existing project to latest templates
grit upgrade                          # Updates admin, web, configs
grit upgrade --force                  # Overwrite without prompting
%[1]s

---

## Project Structure

After running %[1]s new %[2]s%[1]s, you get:

%[1]s
%[2]s/
├── .env                          # Environment variables (DB, Redis, S3, AI keys)
├── docker-compose.yml            # PostgreSQL, Redis, MinIO, Mailhog
├── turbo.json                    # Monorepo task orchestration
├── pnpm-workspace.yaml           # Workspace definition
├── GRIT_SKILL.md                 # This file — AI assistant guide
│
├── packages/shared/              # Shared between frontend and backend
│   ├── schemas/                  # Zod validation schemas
│   │   ├── index.ts              # Re-exports all schemas
│   │   └── user.ts               # UserSchema, LoginSchema, etc.
│   ├── types/                    # TypeScript interfaces
│   │   ├── index.ts              # Re-exports all types
│   │   ├── user.ts               # User, CreateUserInput, etc.
│   │   └── api.ts                # ApiResponse, PaginatedResponse, ApiError
│   └── constants/                # Shared constants
│       └── index.ts              # API_ROUTES, ROLES, etc.
│
├── apps/
│   ├── api/                      # Go backend
│   │   ├── go.mod
│   │   ├── cmd/server/main.go    # Entry point
│   │   └── internal/
│   │       ├── config/config.go       # Loads .env into Config struct
│   │       ├── database/database.go   # GORM connection setup
│   │       ├── models/                # GORM models
│   │       │   ├── user.go            # User model + AutoMigrate
│   │       │   └── upload.go          # Upload model
│   │       ├── handlers/              # HTTP request handlers
│   │       │   ├── auth.go            # Register, Login, Refresh, Logout, Me
│   │       │   ├── user.go            # CRUD for users
│   │       │   ├── upload.go          # File upload endpoints
│   │       │   ├── ai.go              # AI completion/chat/stream
│   │       │   ├── jobs.go            # Job queue admin endpoints
│   │       │   └── cron.go            # Cron task listing
│   │       ├── services/              # Business logic layer
│   │       │   └── auth.go            # JWT generation, validation
│   │       ├── middleware/             # Gin middleware
│   │       │   ├── auth.go            # JWT auth + role-based access
│   │       │   ├── cors.go            # CORS configuration
│   │       │   ├── logger.go          # Request logging
│   │       │   └── cache.go           # Response caching middleware
│   │       ├── routes/routes.go       # Route registration
│   │       ├── mail/                  # Email service (Resend)
│   │       ├── storage/               # File storage (S3/R2/MinIO)
│   │       ├── jobs/                  # Background jobs (asynq)
│   │       ├── cron/cron.go           # Scheduled tasks
│   │       ├── cache/cache.go         # Redis cache service
│   │       └── ai/ai.go              # AI service (Claude/OpenAI)
│   │
│   ├── web/                      # SaaS landing page (Next.js)
│   │   ├── app/
│   │   │   ├── layout.tsx             # Root layout
│   │   │   └── page.tsx               # Landing page with hero, features, CTA
│   │   └── lib/
│   │       └── utils.ts               # Utility functions
│   │
│   └── admin/                    # Admin panel (Filament-like)
│       ├── app/
│       │   ├── layout.tsx             # Root layout (Providers, no sidebar)
│       │   ├── page.tsx               # Redirect to /dashboard or /login
│       │   ├── (auth)/               # Auth pages (no sidebar)
│       │   │   ├── login/page.tsx
│       │   │   ├── sign-up/page.tsx
│       │   │   └── forgot-password/page.tsx
│       │   └── (dashboard)/          # Protected pages (with sidebar)
│       │       ├── layout.tsx         # AdminLayout wrapper
│       │       ├── dashboard/page.tsx # Dashboard with widgets
│       │       ├── resources/         # Resource pages
│       │       │   └── users/page.tsx
│       │       └── system/           # System pages
│       │           ├── jobs/page.tsx
│       │           ├── files/page.tsx
│       │           ├── cron/page.tsx
│       │           └── mail/page.tsx
│       ├── components/
│       │   ├── layout/                # Sidebar, Navbar
│       │   ├── tables/                # DataTable, filters, pagination
│       │   ├── forms/                 # FormBuilder, field components
│       │   ├── widgets/               # StatsCard, ChartWidget
│       │   └── resource/              # ResourcePage renderer
│       ├── hooks/
│       │   ├── use-auth.ts            # Admin auth hooks
│       │   ├── use-resource.ts        # Generic CRUD hooks
│       │   └── use-system.ts          # Hooks for jobs/files/cron
│       ├── resources/                # Resource definitions
│       │   ├── index.ts               # Resource registry
│       │   └── users.ts              # Users resource definition
│       └── lib/
│           ├── resource.ts            # defineResource() + types
│           └── icons.ts               # Lucide icon map
%[1]s

---

## How to Work with Grit Projects

### Adding a New Resource

The most common task. Run:

%[1]sash
grit generate resource Post --fields "title:string,content:text,published:bool,views:int"
%[1]s

This creates **8 files** and injects into **10 existing files**:

**New files created:**
| File | What it contains |
|------|-----------------|
| %[1]sapps/api/internal/models/post.go%[1]s | GORM model with fields, timestamps, soft delete |
| %[1]sapps/api/internal/services/post.go%[1]s | List (paginated, searchable, sortable), GetByID, Create, Update, Delete |
| %[1]sapps/api/internal/handlers/post.go%[1]s | HTTP handlers for GET/POST/PUT/DELETE |
| %[1]spackages/shared/schemas/post.ts%[1]s | Zod schemas: CreatePostSchema, UpdatePostSchema |
| %[1]spackages/shared/types/post.ts%[1]s | TypeScript interface: Post |
| %[1]sapps/web/hooks/use-posts.ts%[1]s | React Query hooks: usePosts, useCreatePost, etc. |
| %[1]sapps/admin/resources/posts.ts%[1]s | Resource definition: table columns, form fields, widgets |
| %[1]sapps/admin/app/(dashboard)/resources/posts/page.tsx%[1]s | Admin page that renders the resource |

**Existing files modified (via marker injection):**
| File | What's injected |
|------|----------------|
| %[1]smodels/user.go%[1]s | %[1]s&Post{}%[1]s added to AutoMigrate |
| %[1]sroutes/routes.go%[1]s | Handler init + CRUD routes registered |
| %[1]sschemas/index.ts%[1]s | Schema exports added |
| %[1]stypes/index.ts%[1]s | Type export added |
| %[1]sconstants/index.ts%[1]s | API route constants added |
| %[1]sresources/index.ts%[1]s | Resource imported and registered |

### Supported Field Types

| Type | Go Type | TypeScript | Zod | Form Input |
|------|---------|-----------|-----|-----------|
| %[1]sstring%[1]s | %[1]sstring%[1]s | %[1]sstring%[1]s | %[1]sz.string()%[1]s | Text input |
| %[1]stext%[1]s | %[1]sstring%[1]s | %[1]sstring%[1]s | %[1]sz.string()%[1]s | Textarea |
| %[1]sint%[1]s | %[1]sint%[1]s | %[1]snumber%[1]s | %[1]sz.number().int()%[1]s | Number input |
| %[1]suint%[1]s | %[1]suint%[1]s | %[1]snumber%[1]s | %[1]sz.number().int().nonneg()%[1]s | Number input |
| %[1]sfloat%[1]s | %[1]sfloat64%[1]s | %[1]snumber%[1]s | %[1]sz.number()%[1]s | Number input |
| %[1]sbool%[1]s | %[1]sbool%[1]s | %[1]sboolean%[1]s | %[1]sz.boolean()%[1]s | Toggle switch |
| %[1]sdatetime%[1]s | %[1]s*time.Time%[1]s | %[1]sstring | null%[1]s | %[1]sz.string().nullable()%[1]s | Datetime picker |
| %[1]sdate%[1]s | %[1]s*time.Time%[1]s | %[1]sstring | null%[1]s | %[1]sz.string().nullable()%[1]s | Date picker |

### Field Modifiers

Append modifiers after the type with colons:

%[1]sash
grit generate resource Post --fields "title:string,slug:string:unique,email:string:required,bio:text:optional"
%[1]s

| Modifier | Effect |
|----------|--------|
| %[1]sunique%[1]s | Adds %[1]sgorm:"uniqueIndex"%[1]s to the Go model |
| %[1]srequired%[1]s | Marks the field as required (string fields are required by default) |
| %[1]soptional%[1]s | Marks the field as optional (overrides default required for strings) |

### Understanding Markers

Grit uses **marker comments** to know where to inject code. Never delete these:

%[1]sgo
// In models/user.go — where new models are added to AutoMigrate
// grit:models          <-- DON'T DELETE THIS

// In routes/routes.go
// grit:handlers            <-- DON'T DELETE THIS
// grit:routes:protected    <-- DON'T DELETE THIS
// grit:routes:admin        <-- DON'T DELETE THIS
%[1]s

%[1]stypescript
// In schemas/index.ts
// grit:schemas              <-- DON'T DELETE THIS

// In types/index.ts
// grit:types                <-- DON'T DELETE THIS

// In constants/index.ts
// grit:api-routes           <-- DON'T DELETE THIS

// In resources/index.ts
// grit:resources            <-- DON'T DELETE THIS
// grit:resource-list        <-- DON'T DELETE THIS
%[1]s

---

## API Conventions

### Response Format

**Always** use this format in handlers:

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

### Standard Error Codes

| Code | HTTP Status | When |
|------|------------|------|
| %[1]sVALIDATION_ERROR%[1]s | 422 | Invalid input data |
| %[1]sNOT_FOUND%[1]s | 404 | Resource doesn't exist |
| %[1]sUNAUTHORIZED%[1]s | 401 | Missing or invalid JWT |
| %[1]sFORBIDDEN%[1]s | 403 | Insufficient role/permissions |
| %[1]sINTERNAL_ERROR%[1]s | 500 | Unexpected server error |
| %[1]sCONFLICT%[1]s | 409 | Duplicate key (e.g., unique email) |

### Authentication Flow

The API uses JWT tokens:

%[1]s
POST /api/auth/register  → Creates user, returns { access_token, refresh_token }
POST /api/auth/login     → Returns { access_token, refresh_token }
POST /api/auth/refresh   → New access_token from refresh_token
POST /api/auth/logout    → Invalidates refresh token
GET  /api/auth/me        → Returns current user (requires auth header)
%[1]s

Access tokens expire in 15 minutes. Refresh tokens last 7 days. The frontend automatically refreshes via Axios interceptor.

### Route Groups

%[1]sgo
// Public routes — no auth required
public := router.Group("/api/auth")

// Protected routes — requires valid JWT
protected := router.Group("/api")
protected.Use(middleware.Auth(cfg.JWTSecret))

// Admin routes — requires JWT + admin role
admin := protected.Group("/admin")
admin.Use(middleware.RequireRole("admin"))
%[1]s

---

## Admin Panel — Resource Definitions

The admin panel uses a **resource definition system**. Each resource is defined in TypeScript:

%[1]stypescript
// apps/admin/resources/posts.ts
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
      { key: "views", label: "Views", sortable: true, format: "number" },
      { key: "created_at", label: "Created", sortable: true, format: "relative" },
    ],
    defaultSort: { key: "created_at", direction: "desc" },
    filters: [
      { key: "published", label: "Published", type: "boolean" },
    ],
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

  dashboard: {
    widgets: [
      {
        type: "stat",
        label: "Total Posts",
        endpoint: "/api/posts?page_size=1",
        icon: "FileText",
        color: "accent",
        format: "number",
        colSpan: 1,
      },
    ],
  },
});
%[1]s

### Form Field Types

| Type | Component | Notes |
|------|----------|-------|
| %[1]stext%[1]s | Text input | Supports prefix, suffix |
| %[1]stextarea%[1]s | Multi-line textarea | Configurable rows |
| %[1]snumber%[1]s | Number input | Supports min, max, step |
| %[1]sselect%[1]s | Dropdown select | Requires options array |
| %[1]sdate%[1]s | Date picker | |
| %[1]sdatetime%[1]s | Datetime picker | |
| %[1]stoggle%[1]s | On/off switch | |
| %[1]scheckbox%[1]s | Checkbox | |
| %[1]sradio%[1]s | Radio group | Requires options array |
| %[1]simage%[1]s | Image upload (drag & drop) | Uses react-dropzone, uploads to /api/uploads |

---

## Working with Go Models

### Model Pattern

%[1]sgo
package models

import (
    "time"
    "gorm.io/gorm"
)

type Post struct {
    ID        uint           %[1]sgorm:"primarykey" json:"id"%[1]s
    Title     string         %[1]sgorm:"size:255;not null" json:"title" binding:"required"%[1]s
    Slug      string         %[1]sgorm:"size:255;uniqueIndex" json:"slug"%[1]s
    Content   string         %[1]sgorm:"type:text" json:"content"%[1]s
    Published bool           %[1]sgorm:"default:false" json:"published"%[1]s
    Views     int            %[1]sjson:"views"%[1]s
    UserID    uint           %[1]sjson:"user_id"%[1]s
    User      User           %[1]sgorm:"foreignKey:UserID" json:"user,omitempty"%[1]s
    CreatedAt time.Time      %[1]sjson:"created_at"%[1]s
    UpdatedAt time.Time      %[1]sjson:"updated_at"%[1]s
    DeletedAt gorm.DeletedAt %[1]sgorm:"index" json:"-"%[1]s
}
%[1]s

Key rules:
- Always include %[1]sID%[1]s, %[1]sCreatedAt%[1]s, %[1]sUpdatedAt%[1]s, %[1]sDeletedAt%[1]s
- Use %[1]sjson:"-"%[1]s for DeletedAt (soft delete, hidden from API)
- Use %[1]sbinding:"required"%[1]s for required fields in handlers
- Use %[1]sgorm:"size:255"%[1]s for string fields, %[1]sgorm:"type:text"%[1]s for long text

---

## Frontend Patterns

### React Query Hooks

All data fetching uses auto-generated hooks:

%[1]stypescript
// apps/web/hooks/use-posts.ts (auto-generated by grit generate)
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
      const { data } = await apiClient.get(%[1]s/api/posts?${params}%[1]s);
      return data;
    },
  });
}
%[1]s

---

## Batteries (Phase 4 Services)

### File Storage

%[1]sgo
// Upload a file
storage.Upload(ctx, "uploads/2024/01/photo.jpg", reader, "image/jpeg")

// Get public URL
url := storage.GetURL("uploads/2024/01/photo.jpg")

// Get signed URL (temporary access)
url, err := storage.GetSignedURL(ctx, key, 1*time.Hour)
%[1]s

API endpoint: %[1]sPOST /api/uploads%[1]s (multipart form, max 10MB)

### Email

%[1]sgo
mailer.Send(ctx, mail.SendOptions{
    To:       "user@example.com",
    Subject:  "Welcome!",
    Template: "welcome",
    Data:     map[string]interface{}{"Name": "John"},
})
%[1]s

Built-in templates: %[1]swelcome%[1]s, %[1]spassword-reset%[1]s, %[1]semail-verification%[1]s, %[1]snotification%[1]s

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

// As middleware
router.GET("/api/posts", middleware.CacheResponse(cache, 5*time.Minute), handler.List)
%[1]s

### AI Integration

%[1]sgo
result, err := ai.Complete(ctx, ai.CompletionRequest{
    Prompt: "Summarize this article...",
})

// Streaming
ai.Stream(ctx, req, func(chunk string) { /* SSE */ })
%[1]s

---

## Configuration (.env)

%[1]sash
DATABASE_URL=postgres://postgres:password@localhost:5432/%[2]s?sslmode=disable
PORT=8080
JWT_SECRET=your-secret-key
WEB_URL=http://localhost:3000
ADMIN_URL=http://localhost:3001
REDIS_URL=redis://localhost:6379
STORAGE_ENDPOINT=http://localhost:9000
STORAGE_BUCKET=%[2]s
STORAGE_ACCESS_KEY=minioadmin
STORAGE_SECRET_KEY=minioadmin
RESEND_API_KEY=re_xxxxx
AI_PROVIDER=claude
AI_API_KEY=sk-ant-xxxxx
%[1]s

---

## Docker Services

%[1]sash
docker compose up -d    # Start all services
%[1]s

| Service | Port | Purpose |
|---------|------|---------|
| PostgreSQL | 5432 | Database |
| Redis | 6379 | Cache + job queue |
| MinIO | 9000/9001 | Local S3 storage |
| Mailhog | 1025/8025 | Email testing |

---

## Naming Conventions

| Thing | Convention | Example |
|-------|-----------|---------|
| Go files | %[1]ssnake_case.go%[1]s | %[1]suser_handler.go%[1]s |
| Go structs | %[1]sPascalCase%[1]s | %[1]stype PostHandler struct%[1]s |
| TypeScript files | %[1]skebab-case.ts%[1]s | %[1]suse-posts.ts%[1]s |
| React components | %[1]sPascalCase.tsx%[1]s | %[1]sDataTable.tsx%[1]s |
| API routes | %[1]s/api/plural%[1]s | %[1]s/api/posts%[1]s |
| DB tables | %[1]splural_snake%[1]s | %[1]sblog_posts%[1]s |

---

## Design System

### Colors (Dark Theme Default)

| Token | Value | Usage |
|-------|-------|-------|
| %[1]sbg-primary%[1]s | %[1]s#0a0a0f%[1]s | Page background |
| %[1]sbg-secondary%[1]s | %[1]s#111118%[1]s | Card background |
| %[1]saccent%[1]s | %[1]s#6c5ce7%[1]s | Primary accent (purple) |
| %[1]ssuccess%[1]s | %[1]s#00b894%[1]s | Success states |
| %[1]sdanger%[1]s | %[1]s#ff6b6b%[1]s | Error states |

### Fonts

- **UI:** DM Sans (400, 500, 600, 700)
- **Code:** JetBrains Mono (400, 500, 600)

---

## Common Tasks for AI Assistants

### "Add a new field to an existing resource"

1. Update the Go model (%[1]sapps/api/internal/models/<name>.go%[1]s) — add the field
2. Update the handler if the field needs special handling
3. Update the Zod schema (%[1]spackages/shared/schemas/<name>.ts%[1]s)
4. Update the TypeScript type (%[1]spackages/shared/types/<name>.ts%[1]s)
5. Update the admin resource definition (%[1]sapps/admin/resources/<name>.ts%[1]s) — add column + form field
6. Restart the API (GORM auto-migrates on startup)

### "Add a new API endpoint"

1. Create or update handler in %[1]sapps/api/internal/handlers/%[1]s
2. Register the route in %[1]sapps/api/internal/routes/routes.go%[1]s
3. Create React Query hook in %[1]sapps/web/hooks/%[1]s or %[1]sapps/admin/hooks/%[1]s

### "Add a relationship between resources"

In the Go model:
%[1]sgo
type Post struct {
    CategoryID uint     %[1]sjson:"category_id"%[1]s
    Category   Category %[1]sgorm:"foreignKey:CategoryID" json:"category,omitempty"%[1]s
}
%[1]s

In the handler, preload related data:
%[1]sgo
query.Preload("Category").Find(&posts)
%[1]s

---

## Important: Don't Break These

1. **Never delete marker comments** (%[1]s// grit:models%[1]s, %[1]s// grit:routes:protected%[1]s, etc.)
2. **Always follow the response format** (%[1]s{ data, message }%[1]s or %[1]s{ data, meta }%[1]s or %[1]s{ error: { code, message } }%[1]s)
3. **Always handle errors in Go** — never ignore with %[1]s_%[1]s
4. **Keep the folder structure** — don't move files to non-standard locations
5. **Use React Query for all data fetching** — no raw %[1]sfetch%[1]s in components
6. **Use Zod for validation** — shared between frontend and backend
7. **Use Tailwind + shadcn/ui for styling** — no custom CSS files
8. **Use App Router** — never Pages Router in Next.js

---

## GORM Studio

The API embeds GORM Studio, a visual database browser, at %[1]s/studio%[1]s.

Access: %[1]shttp://localhost:8080/studio%[1]s

---

*Built with Grit. Go + React. Built with Grit.*
`, "`", opts.ProjectName)
}
