# Grit Framework — LLM Skill Guide

> **This document teaches AI assistants (Claude, Cursor, Kilo Code, etc.) how to work with the Grit framework.** Read this file completely before writing any code in a Grit project.

---

## What is Grit?

Grit is a full-stack meta-framework that combines **Go** (backend) + **React/Next.js** (frontend) in a monorepo. It provides:

- A **CLI tool** (`grit`) that scaffolds entire projects and generates full-stack resources
- A **Go API** with Gin + GORM + PostgreSQL
- A **Next.js web app** with App Router + Tailwind + shadcn/ui
- A **Filament-like admin panel** with resource definitions, DataTables, forms, and widgets
- **Native desktop apps** with Wails (Go + React + SQLite) — `grit new-desktop`
- **Social login** out of the box — Google + GitHub OAuth2 via Gothic, with account linking
- **Batteries included**: file storage (S3), email (Resend), background jobs (asynq), cron, Redis caching, AI integration (Claude/OpenAI), security (Sentinel), observability (Pulse)
- **Auto-generated API docs** via gin-docs — zero-annotation OpenAPI spec with interactive UI
- A **shared package** with Zod schemas, TypeScript types, and constants

**Think of it as:** Laravel + Filament, but with Go + React instead of PHP + Blade.

---

## Quick Reference — CLI Commands

```bash
# Create a new project
grit new myapp                    # Full monorepo (Go API + Next.js web + admin)
grit new myapp --api              # Go API only
grit new myapp --full             # Everything + Expo mobile + docs site

# Generate a resource (full-stack CRUD)
grit generate resource Post --fields "title:string,content:text,published:bool"
grit generate resource Post --fields "title:string,slug:string:unique,views:int"
grit generate resource Invoice --from invoice.yaml
grit generate resource Category -i   # Interactive mode

# Sync Go types to TypeScript
grit sync

# Add a new role to the project
grit add role MODERATOR             # Injects into Go, Zod, TypeScript, admin

# Desktop projects
grit new-desktop mydesktop          # Scaffold Wails + Go + React desktop app
grit start                          # Auto-detect project type, run wails dev
grit compile                        # Build native executable (wails build)
grit studio                         # Open database browser (both web + desktop)

# Development (web projects)
grit dev                            # Start all services with hot-reload
grit start server                   # Start Go API only
grit start client                   # Start frontend only
```

---

## Project Structure

After running `grit new myapp`, you get:

```
myapp/
├── .env                          # Environment variables (DB, Redis, S3, AI keys)
├── docker-compose.yml            # PostgreSQL, Redis, MinIO, Mailhog
├── turbo.json                    # Monorepo task orchestration
├── pnpm-workspace.yaml           # Workspace definition
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
│   │       │   └── upload.go          # Upload model (Phase 4)
│   │       ├── handlers/              # HTTP request handlers
│   │       │   ├── auth.go            # Register, Login, Refresh, Logout, Me, ForgotPassword, ResetPassword, OAuth
│   │       │   ├── user.go            # CRUD for users
│   │       │   ├── upload.go          # File upload endpoints
│   │       │   ├── ai.go             # AI completion/chat/stream
│   │       │   ├── jobs.go           # Job queue admin endpoints
│   │       │   └── cron.go           # Cron task listing
│   │       ├── services/              # Business logic layer
│   │       │   └── auth.go            # JWT generation, validation, hashing
│   │       ├── middleware/            # Gin middleware
│   │       │   ├── auth.go            # JWT auth + role-based access
│   │       │   ├── cors.go            # CORS configuration
│   │       │   ├── logger.go          # Request logging
│   │       │   └── cache.go           # Response caching middleware
│   │       ├── routes/routes.go       # Route registration + Services struct
│   │       ├── mail/                  # Email service
│   │       │   ├── mailer.go          # Resend API client
│   │       │   └── templates.go       # HTML email templates
│   │       ├── storage/               # File storage
│   │       │   ├── storage.go         # S3-compatible client (AWS/R2/MinIO)
│   │       │   └── image.go           # Image resize + thumbnail
│   │       ├── jobs/                  # Background jobs
│   │       │   ├── client.go          # Job enqueue client
│   │       │   └── workers.go         # Job handlers
│   │       ├── cron/cron.go           # Scheduled tasks
│   │       ├── cache/cache.go         # Redis cache service
│   │       └── ai/ai.go              # AI service (Claude/OpenAI)
│   │
│   │   # Mounted tools (auto-configured in routes.go):
│   │   # /docs         → gin-docs API documentation (OpenAPI 3.1 + Scalar UI)
│   │   # /studio       → GORM Studio database browser
│   │   # /sentinel/*   → Sentinel security dashboard (WAF, rate limiting)
│   │   # /pulse        → Pulse observability dashboard (tracing, metrics)
│   │
│   ├── web/                      # Next.js main frontend
│   │   ├── app/
│   │   │   ├── layout.tsx             # Root layout + providers
│   │   │   ├── error.tsx              # Error boundary (client component)
│   │   │   ├── not-found.tsx          # Custom 404 page
│   │   │   ├── global-error.tsx       # Root error boundary (replaces entire HTML)
│   │   │   ├── (auth)/               # Auth pages (login, register, forgot-password, callback)
│   │   │   └── (dashboard)/          # Protected pages with sidebar
│   │   ├── components/ui/            # shadcn/ui components
│   │   ├── hooks/                    # React Query hooks (auto-generated)
│   │   └── lib/
│   │       ├── api-client.ts          # Axios instance with JWT interceptor
│   │       └── query-client.ts        # React Query setup
│   │
│   └── admin/                    # Admin panel (Filament-like)
│       ├── app/
│       │   ├── layout.tsx             # Admin shell (sidebar + navbar)
│       │   ├── page.tsx               # Dashboard with widgets
│       │   ├── error.tsx              # Error boundary (client component)
│       │   ├── not-found.tsx          # Custom 404 page
│       │   ├── global-error.tsx       # Root error boundary (replaces entire HTML)
│       │   ├── resources/             # Resource pages (auto-generated)
│       │   │   └── users/page.tsx
│       │   └── system/               # System management pages
│       │       ├── jobs/page.tsx       # Job queue dashboard
│       │       ├── files/page.tsx      # File browser
│       │       ├── cron/page.tsx       # Cron tasks viewer
│       │       ├── mail/page.tsx       # Email template preview
│       │       └── security/page.tsx   # Sentinel security dashboard
│       ├── components/
│       │   ├── layout/                # Sidebar, Navbar, AdminLayout
│       │   ├── tables/                # DataTable, filters, pagination
│       │   ├── forms/                 # FormBuilder, field components
│       │   ├── widgets/               # StatsCard, ChartWidget, etc.
│       │   └── resource/              # ResourcePage renderer
│       ├── hooks/
│       │   ├── use-auth.ts            # Admin auth hooks
│       │   ├── use-resource.ts        # Generic CRUD hooks for any resource
│       │   └── use-system.ts          # Hooks for jobs/files/cron
│       ├── resources/                # Resource definitions
│       │   ├── index.ts               # Resource registry
│       │   └── users.ts              # Users resource definition
│       └── lib/
│           ├── resource.ts            # defineResource() + types
│           └── icons.ts               # Lucide icon map
```

---

## How to Work with Grit Projects

### Adding a New Resource

The most common task. Run:

```bash
grit generate resource Post --fields "title:string,content:text,published:bool,views:int"
```

This creates **8 files** and injects into **10 existing files**:

**New files created:**
| File | What it contains |
|------|-----------------|
| `apps/api/internal/models/post.go` | GORM model with fields, timestamps, soft delete |
| `apps/api/internal/services/post.go` | List (paginated, searchable, sortable), GetByID, Create, Update, Delete |
| `apps/api/internal/handlers/post.go` | HTTP handlers for GET/POST/PUT/DELETE |
| `packages/shared/schemas/post.ts` | Zod schemas: CreatePostSchema, UpdatePostSchema |
| `packages/shared/types/post.ts` | TypeScript interface: Post |
| `apps/web/hooks/use-posts.ts` | React Query hooks: usePosts, useGetPost, useCreatePost, useUpdatePost, useDeletePost |
| `apps/admin/resources/posts.ts` | Resource definition: table columns, form fields, widgets |
| `apps/admin/app/resources/posts/page.tsx` | Admin page that renders the resource |

**Existing files modified (via marker injection):**
| File | What's injected |
|------|----------------|
| `models/user.go` | `&Post{}` added to AutoMigrate |
| `routes/routes.go` | Handler init + CRUD routes registered |
| `schemas/index.ts` | Schema exports added |
| `types/index.ts` | Type export added |
| `constants/index.ts` | API route constants added |
| `resources/index.ts` | Resource imported and registered |

### Supported Field Types

| Type | Go Type | TypeScript | Zod | Form Input |
|------|---------|-----------|-----|-----------|
| `string` | `string` | `string` | `z.string()` | Text input |
| `text` | `string` | `string` | `z.string()` | Textarea |
| `int` | `int` | `number` | `z.number().int()` | Number input |
| `uint` | `uint` | `number` | `z.number().int().nonneg()` | Number input |
| `float` | `float64` | `number` | `z.number()` | Number input |
| `bool` | `bool` | `boolean` | `z.boolean()` | Toggle switch |
| `datetime` | `*time.Time` | `string \| null` | `z.string().nullable()` | Datetime picker |
| `date` | `*time.Time` | `string \| null` | `z.string().nullable()` | Date picker |
| `richtext` | `string` | `string` | `z.string()` | Tiptap WYSIWYG editor |
| `slug` | `string` | `string` | Auto-generated | Hidden from forms |
| `string_array` | `datatypes.JSONSlice[string]` | `string[]` | `z.array(z.string())` | Tag input |
| `belongs_to:Resource` | `uint` (FK) | `number` | `z.number()` | Relationship select |
| `many_to_many:Resource` | Junction table | `number[]` | `z.array(z.number())` | Multi-select |

### Field Modifiers

Append modifiers after the type with colons:

```bash
grit generate resource Post --fields "title:string,slug:string:unique,email:string:required,bio:text:optional"
```

| Modifier | Effect |
|----------|--------|
| `unique` | Adds `gorm:"uniqueIndex"` to the Go model |
| `required` | Marks the field as required (string fields are required by default) |
| `optional` | Marks the field as optional (overrides default required for strings) |

### Understanding Markers

Grit uses **marker comments** to know where to inject code. Never delete these:

```go
// In models/user.go — where new models are added to AutoMigrate
db.AutoMigrate(
    &User{},
    &Upload{},
    // grit:models          <-- DON'T DELETE THIS
)

// In routes/routes.go — where handlers are initialized
userHandler := &handlers.UserHandler{DB: db}
// grit:handlers            <-- DON'T DELETE THIS

// In routes/routes.go — where protected routes are registered
protected.GET("/users", userHandler.List)
// grit:routes:protected    <-- DON'T DELETE THIS

// In routes/routes.go — where admin routes are registered
admin.DELETE("/users/:id", userHandler.Delete)
// grit:routes:admin        <-- DON'T DELETE THIS
```

```typescript
// In schemas/index.ts
export { CreateUserSchema, UpdateUserSchema } from "./user";
// grit:schemas              <-- DON'T DELETE THIS

// In types/index.ts
export type { User } from "./user";
// grit:types                <-- DON'T DELETE THIS

// In constants/index.ts — API route definitions
USERS: { LIST: "/api/users", ... },
// grit:api-routes           <-- DON'T DELETE THIS

// In resources/index.ts
import { usersResource } from "./users";
// grit:resources            <-- DON'T DELETE THIS

export const resources: ResourceDefinition[] = [
    usersResource,
    // grit:resource-list    <-- DON'T DELETE THIS
];
```

---

## API Conventions

### Response Format

**Always** use this format in handlers:

```go
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
        "details": gin.H{
            "email": "This field is required",
        },
    },
})
```

### Standard Error Codes

| Code | HTTP Status | When |
|------|------------|------|
| `VALIDATION_ERROR` | 422 | Invalid input data |
| `NOT_FOUND` | 404 | Resource doesn't exist |
| `UNAUTHORIZED` | 401 | Missing or invalid JWT |
| `FORBIDDEN` | 403 | Insufficient role/permissions |
| `INTERNAL_ERROR` | 500 | Unexpected server error |
| `CONFLICT` | 409 | Duplicate key (e.g., unique email) |

### Authentication Flow

The API uses JWT tokens:

```
POST /api/auth/register  → Creates user, returns { access_token, refresh_token }
POST /api/auth/login     → Returns { access_token, refresh_token }
POST /api/auth/refresh   → New access_token from refresh_token
POST /api/auth/logout    → Invalidates refresh token
GET  /api/auth/me        → Returns current user (requires auth header)
```

Access tokens expire in 15 minutes. Refresh tokens last 7 days. The frontend automatically refreshes via Axios interceptor.

### Social Login (OAuth2)

Grit includes Google and GitHub OAuth2 authentication via [Gothic](https://github.com/markbates/goth). Social login buttons appear on all login and register pages across all admin style variants.

**OAuth Flow:**
```
1. User clicks "Sign in with Google" → browser navigates to /api/auth/oauth/google
2. Backend (Gothic) redirects to Google consent screen
3. Google redirects back to /api/auth/oauth/google/callback
4. Backend completes auth, finds/creates user in DB, generates JWT tokens
5. Backend redirects to OAUTH_FRONTEND_URL/auth/callback?access_token=XXX&refresh_token=XXX
6. Frontend callback page stores tokens in cookies, fetches /api/auth/me, redirects to dashboard
```

**OAuth Routes:**
```
GET /api/auth/oauth/google          → Begins Google OAuth flow
GET /api/auth/oauth/google/callback → Google OAuth callback
GET /api/auth/oauth/github          → Begins GitHub OAuth flow
GET /api/auth/oauth/github/callback → GitHub OAuth callback
```

**Account Linking:**
- If an OAuth email matches an existing user → link accounts (sets GoogleID/GithubID), login
- If no existing user → create new user (Provider="google"/"github", empty password), login
- OAuth-only users have an empty password — they can set one later via profile
- Email/password login returns a helpful error for password-less accounts: "Please sign in with Google/GitHub"

**User Model Fields (for OAuth):**
```go
Provider string `gorm:"size:50;default:'local'" json:"provider"`  // "local", "google", "github"
GoogleID string `gorm:"size:255" json:"-"`
GithubID string `gorm:"size:255" json:"-"`
```

**Setup:** Add your OAuth credentials to `.env`:
```bash
GOOGLE_CLIENT_ID=your-google-client-id
GOOGLE_CLIENT_SECRET=your-google-client-secret
GITHUB_CLIENT_ID=your-github-client-id
GITHUB_CLIENT_SECRET=your-github-client-secret
OAUTH_FRONTEND_URL=http://localhost:3001
```

Providers are conditionally initialized — if credentials are empty, that provider is simply skipped (no crash).

### Route Groups

```go
// Public routes — no auth required
public := router.Group("/api/auth")

// Protected routes — requires valid JWT
protected := router.Group("/api")
protected.Use(middleware.Auth(cfg.JWTSecret))

// Admin routes — requires JWT + admin role
admin := protected.Group("/admin")
admin.Use(middleware.RequireRole("admin"))
```

---

## Admin Panel — Resource Definitions

The admin panel uses a **resource definition system**. Each resource is defined in TypeScript:

```typescript
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
    actions: { create: true, edit: true, delete: true, export: true },
  },

  form: {
    fields: [
      { key: "title", label: "Title", type: "text", required: true, placeholder: "Enter title..." },
      { key: "content", label: "Content", type: "textarea", required: true },
      { key: "published", label: "Published", type: "toggle", default: false },
    ],
    layout: "single",
  },

  dashboard: {
    widgets: [
      {
        type: "stat",
        label: "Total Posts",
        endpoint: "/api/posts",
        valueKey: "meta.total",
        icon: "FileText",
        color: "accent",
      },
    ],
  },
});
```

### Column Format Types

| Format | Renders as |
|--------|-----------|
| `text` | Plain text (default) |
| `number` | Formatted number |
| `currency` | $1,234.56 |
| `boolean` | Green/red badge |
| `date` | Formatted date |
| `relative` | "2 hours ago" |
| `badge` | Colored badge |
| `image` | Thumbnail image |

### Form Field Types

| Type | Component |
|------|----------|
| `text` | Text input |
| `textarea` | Multi-line textarea |
| `number` | Number input with step |
| `select` | Dropdown select |
| `date` | Date picker |
| `toggle` | On/off switch |
| `checkbox` | Checkbox |
| `radio` | Radio group |
| `image` | Image upload (drag & drop via Dropzone) |
| `richtext` | Tiptap WYSIWYG editor (bold, italic, headings, lists, code, links) |
| `relationship-select` | Searchable dropdown for belongs_to relationships |
| `multi-relationship-select` | Multi-select tags for many_to_many relationships |

### Form View Variants

Set `formView` in the resource definition to change how forms render:

| View | Description |
|------|-----------|
| `modal` | Default — form opens in a modal dialog |
| `page` | Full-page form at a dedicated URL |
| `modal-steps` | Multi-step wizard in a modal with step indicators |
| `page-steps` | Multi-step wizard on a full page with progress bar |

For multi-step forms, add `steps` and assign each field to a step:
```typescript
form: {
  formView: "modal-steps",
  steps: [
    { key: "basic", label: "Basic Info" },
    { key: "details", label: "Details" },
  ],
  fields: [
    { key: "title", label: "Title", type: "text", step: "basic" },
    { key: "content", label: "Content", type: "richtext", step: "details" },
  ],
}
```

### Dropzone Component

A reusable file upload component with 5 variants:

```tsx
import { Dropzone } from "@/components/ui/dropzone";

<Dropzone
  variant="default"     // "default" | "compact" | "minimal" | "avatar" | "inline"
  maxFiles={5}
  maxSize={1024 * 1024 * 10} // 10MB
  onFilesChange={setFiles}
  accept={{
    "image/*": [],
    "application/pdf": [],
  }}
/>
```

| Variant | Description |
|---------|------------|
| `default` | Full drop zone with icon, text, and file preview cards |
| `compact` | Single-line drop area with file chips |
| `minimal` | Small button-style upload trigger |
| `avatar` | Circular avatar upload with preview |
| `inline` | Card-style with browse button |

### Adding Resources to the Registry

After creating a resource definition, register it:

```typescript
// apps/admin/resources/index.ts
import { usersResource } from "./users";
import { postsResource } from "./posts";
// grit:resources

import type { ResourceDefinition } from "@/lib/resource";

export const resources: ResourceDefinition[] = [
  usersResource,
  postsResource,
  // grit:resource-list
];
```

The sidebar automatically picks up registered resources.

---

## Working with Go Models

### Model Pattern

```go
package models

import (
    "time"
    "gorm.io/gorm"
)

type Post struct {
    ID        uint           `gorm:"primarykey" json:"id"`
    Title     string         `gorm:"size:255;not null" json:"title" binding:"required"`
    Content   string         `gorm:"type:text" json:"content"`
    Published bool           `gorm:"default:false" json:"published"`
    Views     int            `json:"views"`
    UserID    uint           `json:"user_id"`
    User      User           `gorm:"foreignKey:UserID" json:"user,omitempty"`
    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
    DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
```

Key rules:
- Always include `ID`, `CreatedAt`, `UpdatedAt`, `DeletedAt`
- Use `json:"-"` for DeletedAt (soft delete, hidden from API)
- Use `binding:"required"` for required fields in handlers
- Use `gorm:"size:255"` for string fields, `gorm:"type:text"` for long text

### Handler Pattern

```go
package handlers

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
    "myapp/apps/api/internal/models"
)

type PostHandler struct {
    DB *gorm.DB
}

func (h *PostHandler) List(c *gin.Context) {
    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
    search := c.Query("search")
    sortBy := c.DefaultQuery("sort_by", "created_at")
    sortDir := c.DefaultQuery("sort_dir", "desc")

    query := h.DB.Model(&models.Post{})

    // Search across string fields
    if search != "" {
        query = query.Where("title ILIKE ?", "%"+search+"%")
    }

    var total int64
    query.Count(&total)

    var posts []models.Post
    offset := (page - 1) * pageSize
    query.Order(sortBy + " " + sortDir).Offset(offset).Limit(pageSize).Find(&posts)

    pages := int(math.Ceil(float64(total) / float64(pageSize)))

    c.JSON(http.StatusOK, gin.H{
        "data": posts,
        "meta": gin.H{
            "total": total, "page": page, "page_size": pageSize, "pages": pages,
        },
    })
}

func (h *PostHandler) Create(c *gin.Context) {
    var req struct {
        Title     string `json:"title" binding:"required"`
        Content   string `json:"content"`
        Published bool   `json:"published"`
    }

    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusUnprocessableEntity, gin.H{
            "error": gin.H{"code": "VALIDATION_ERROR", "message": err.Error()},
        })
        return
    }

    post := models.Post{
        Title:     req.Title,
        Content:   req.Content,
        Published: req.Published,
    }

    if err := h.DB.Create(&post).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": gin.H{"code": "INTERNAL_ERROR", "message": "Failed to create post"},
        })
        return
    }

    c.JSON(http.StatusCreated, gin.H{
        "data":    post,
        "message": "Post created successfully",
    })
}
```

### Service Pattern (for complex business logic)

```go
package services

import "gorm.io/gorm"

type PostService struct {
    DB *gorm.DB
}

type ListParams struct {
    Page     int
    PageSize int
    Search   string
    SortBy   string
    SortDir  string
}

func (s *PostService) List(params ListParams) ([]models.Post, int64, error) {
    // Query building, search, pagination logic
}
```

---

## Frontend Patterns

### React Query Hooks

All data fetching uses auto-generated hooks:

```typescript
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
      const { data } = await apiClient.get(`/api/posts?${params}`);
      return data;
    },
  });
}

export function useCreatePost() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: async (data: CreatePostInput) => {
      const res = await apiClient.post("/api/posts", data);
      return res.data;
    },
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ["posts"] }),
  });
}
```

### Using Hooks in Components

```tsx
"use client";

import { usePosts, useDeletePost } from "@/hooks/use-posts";

export default function PostsPage() {
  const { data, isLoading, error } = usePosts({ page: 1 });
  const { mutate: deletePost } = useDeletePost();

  if (isLoading) return <div>Loading...</div>;
  if (error) return <div>Error: {error.message}</div>;

  return (
    <div>
      {data.data.map((post) => (
        <div key={post.id}>
          <h2>{post.title}</h2>
          <button onClick={() => deletePost(post.id)}>Delete</button>
        </div>
      ))}
      <p>Total: {data.meta.total}</p>
    </div>
  );
}
```

### API Client

```typescript
// apps/web/lib/api-client.ts
import axios from "axios";

export const apiClient = axios.create({
  baseURL: process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080",
  headers: { "Content-Type": "application/json" },
});

// JWT interceptor — automatically attaches token
apiClient.interceptors.request.use((config) => {
  const token = localStorage.getItem("access_token");
  if (token) config.headers.Authorization = `Bearer ${token}`;
  return config;
});

// Refresh interceptor — auto-refreshes on 401
apiClient.interceptors.response.use(
  (response) => response,
  async (error) => {
    if (error.response?.status === 401) {
      // Try to refresh the token...
    }
    return Promise.reject(error);
  }
);
```

---

## Batteries (Phase 4 Services)

### File Storage

```go
// Upload a file
storage.Upload(ctx, "uploads/2024/01/photo.jpg", reader, "image/jpeg")

// Download a file
reader, err := storage.Download(ctx, "uploads/2024/01/photo.jpg")

// Get public URL
url := storage.GetURL("uploads/2024/01/photo.jpg")

// Get signed URL (temporary access)
url, err := storage.GetSignedURL(ctx, key, 1*time.Hour)

// Delete
storage.Delete(ctx, key)
```

API endpoint: `POST /api/uploads` (multipart form, max 10MB)

### Email

```go
// Send templated email
mailer.Send(ctx, mail.SendOptions{
    To:       "user@example.com",
    Subject:  "Welcome!",
    Template: "welcome",
    Data:     map[string]interface{}{"Name": "John", "AppName": "MyApp"},
})

// Send raw HTML
mailer.SendRaw(ctx, "user@example.com", "Subject", "<h1>Hello</h1>")
```

Built-in templates: `welcome`, `password-reset`, `email-verification`, `notification`

### Background Jobs

```go
// Enqueue a job
jobs.EnqueueSendEmail("user@example.com", "Welcome", "welcome", data)
jobs.EnqueueProcessImage(uploadID, key, mimeType)
jobs.EnqueueTokensCleanup()
```

Admin dashboard at `/system/jobs` shows queue stats, active/pending/failed jobs.

### Redis Cache

```go
// Cache a value
cache.Set(ctx, "user:123", userData, 5*time.Minute)

// Get a value
var user User
err := cache.Get(ctx, "user:123", &user)

// Delete
cache.Delete(ctx, "user:123")

// Delete by pattern
cache.DeletePattern(ctx, "user:*")

// Use as middleware
router.GET("/api/posts", middleware.CacheResponse(cache, 5*time.Minute), handler.List)
```

### AI Integration

```go
// Single completion
result, err := ai.Complete(ctx, ai.CompletionRequest{
    Prompt: "Summarize this article...",
    Model:  "claude-sonnet-4-5-20250929",
})

// Streaming
ai.Stream(ctx, req, func(chunk string) {
    // Send chunk to client via SSE
})
```

API endpoints:
- `POST /api/ai/complete` — Single prompt completion
- `POST /api/ai/chat` — Multi-turn conversation
- `POST /api/ai/stream` — Server-sent events streaming

### Security (Sentinel)

Grit ships with [Sentinel](https://github.com/MUKE-coder/sentinel), a production-grade security suite:

- **WAF** — SQL injection, XSS, path traversal, command injection detection
- **Rate Limiting** — Per-IP request throttling
- **Brute-Force Protection** — Account lockout after failed login attempts
- **Anomaly Detection** — Unusual traffic pattern alerts
- **IP Geolocation** — Country-level blocking
- **Security Headers** — HSTS, CSP, X-Frame-Options
- **Dashboard** — Real-time threat monitoring at `/sentinel/ui`
- **ExcludePaths** — `/pulse/*`, `/studio/*`, `/sentinel/*`, `/docs/*` are excluded from rate limiting by default

Enabled by default. Disable with `SENTINEL_ENABLED=false`.

### Observability (Pulse)

Grit ships with [Pulse](https://github.com/MUKE-coder/pulse), a self-hosted observability SDK:

- **Request Tracing** — HTTP method, path, status, latency for every request
- **Database Monitoring** — Query count, slow queries, connection pool stats
- **Runtime Metrics** — Memory, goroutines, GC stats
- **Error Tracking** — Captures panics and 5xx responses
- **Health Checks** — Built-in DB check + custom checks (Redis, etc.)
- **Alerting** — High error rate, slow response, memory threshold alerts
- **Prometheus** — `/pulse/metrics` endpoint for Grafana integration
- **Dashboard** — Embedded React dashboard at `/pulse`

Enabled by default. Disable with `PULSE_ENABLED=false`.

### API Documentation (gin-docs)

Auto-generated API docs powered by [gin-docs](https://github.com/MUKE-coder/gin-docs):

- Introspects Gin routes and GORM models automatically — zero annotations needed
- Generates OpenAPI 3.1 spec at `/docs/openapi.json`
- Interactive UI at `/docs` (Scalar theme by default, Swagger available)
- Export to Postman or Insomnia

No manual spec maintenance — routes and models are discovered at startup.

---

## Configuration (.env)

```bash
# Database
DATABASE_URL=postgres://postgres:password@localhost:5432/myapp?sslmode=disable

# Server
PORT=8080
JWT_SECRET=your-secret-key

# Frontend URLs
WEB_URL=http://localhost:3000
ADMIN_URL=http://localhost:3001

# Redis
REDIS_URL=redis://localhost:6379

# Storage (S3-compatible)
STORAGE_ENDPOINT=http://localhost:9000
STORAGE_BUCKET=myapp
STORAGE_ACCESS_KEY=minioadmin
STORAGE_SECRET_KEY=minioadmin
STORAGE_REGION=us-east-1

# Email
RESEND_API_KEY=re_xxxxx
MAIL_FROM=noreply@myapp.com

# AI
AI_PROVIDER=claude            # or: openai
AI_API_KEY=sk-ant-xxxxx
AI_MODEL=claude-sonnet-4-5-20250929

# OAuth2 — Social Login
GOOGLE_CLIENT_ID=
GOOGLE_CLIENT_SECRET=
GITHUB_CLIENT_ID=
GITHUB_CLIENT_SECRET=
OAUTH_FRONTEND_URL=http://localhost:3001

# Security (Sentinel)
SENTINEL_ENABLED=true
SENTINEL_USERNAME=admin
SENTINEL_PASSWORD=sentinel

# Observability (Pulse)
PULSE_ENABLED=true
PULSE_USERNAME=admin
PULSE_PASSWORD=pulse
```

---

## Docker Services

### Development

```bash
docker compose up -d    # Start all services
```

| Service | Port | Purpose |
|---------|------|---------|
| PostgreSQL | 5432 | Database |
| Redis | 6379 | Cache + job queue |
| MinIO | 9000/9001 | Local S3 storage |
| Mailhog | 1025/8025 | Email testing |

### Production (`docker-compose.prod.yml`)

The production compose file uses `expose` instead of `ports` (services communicate via Docker network, not host ports), `env_file` for secrets, a named bridge network, and includes a MinIO service. Next.js Dockerfiles accept `NEXT_PUBLIC_API_URL` as a build arg.

### Deploy without Docker

You can also run Grit without Docker by using cloud services:
- **Database:** Neon (PostgreSQL) — `neon.tech`
- **Redis:** Upstash — `upstash.com`
- **Storage:** Cloudflare R2 or Backblaze B2
- **Email:** Resend — `resend.com`

Copy `.env.cloud.example` to `.env` and fill in your cloud credentials.

---

## Naming Conventions

| Thing | Convention | Example |
|-------|-----------|---------|
| Go files | `snake_case.go` | `user_handler.go` |
| Go structs | `PascalCase` | `type PostHandler struct` |
| Go exported | `PascalCase` | `func GetUsers()` |
| Go unexported | `camelCase` | `func parseToken()` |
| TypeScript files | `kebab-case.ts` | `use-posts.ts` |
| React components | `PascalCase.tsx` | `DataTable.tsx` |
| API routes | `/api/plural` | `/api/posts` |
| DB tables | `plural_snake` | `blog_posts` |
| Zod schemas | `PascalCase + Schema` | `CreatePostSchema` |
| Resource slugs | `kebab-case` | `blog-posts` |

---

## Design System

### Colors (Dark Theme Default)

| Token | Value | Usage |
|-------|-------|-------|
| `bg-primary` | `#0a0a0f` | Page background |
| `bg-secondary` | `#111118` | Card background |
| `bg-tertiary` | `#1a1a24` | Elevated surfaces |
| `border` | `#2a2a3a` | Borders |
| `text-primary` | `#e8e8f0` | Main text |
| `text-secondary` | `#9090a8` | Secondary text |
| `text-muted` | `#606078` | Muted text |
| `accent` | `#6c5ce7` | Primary accent (purple) |
| `success` | `#00b894` | Success states |
| `danger` | `#ff6b6b` | Error states |
| `warning` | `#fdcb6e` | Warning states |

### Fonts

- **UI:** Onest (400, 500, 600, 700)
- **Code:** JetBrains Mono (400, 500, 600)

---

## Common Tasks for AI Assistants

### "Add a new field to an existing resource"

1. Update the Go model (`apps/api/internal/models/<name>.go`) — add the field with GORM + JSON tags
2. Update the handler if the field needs special handling
3. Update the Zod schema (`packages/shared/schemas/<name>.ts`) — add to Create/Update schemas
4. Update the TypeScript type (`packages/shared/types/<name>.ts`)
5. Update the admin resource definition (`apps/admin/resources/<name>.ts`) — add column + form field
6. Restart the API (GORM auto-migrates on startup)

### "Add a new API endpoint"

1. Create or update handler in `apps/api/internal/handlers/`
2. Register the route in `apps/api/internal/routes/routes.go`
3. Create React Query hook in `apps/web/hooks/` or `apps/admin/hooks/`

### "Add a relationship between resources"

In the Go model:
```go
type Post struct {
    // ...
    CategoryID uint     `json:"category_id"`
    Category   Category `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
}
```

In the handler, preload related data:
```go
query.Preload("Category").Find(&posts)
```

### "Add a new role"

```bash
grit add role MODERATOR
```

This injects the new role into 7 locations:
1. Go constants (`models/user.go` — role constant)
2. Zod enum (`schemas/user.ts` — UserRole enum)
3. TypeScript union (`types/user.ts` — Role type)
4. Roles constant (`constants/index.ts` — ROLES array)
5. Admin badge mapping (role color in DataTable)
6. Admin filter options (role filter dropdown)
7. Admin form options (role select in user form)

### "Customize the admin table"

Edit the resource definition in `apps/admin/resources/<name>.ts`:
```typescript
table: {
    columns: [
        // Add/remove/reorder columns
        { key: "title", label: "Title", sortable: true, searchable: true },
        { key: "category.name", label: "Category", format: "badge" },
    ],
}
```

### "Add a custom widget to the dashboard"

In the resource definition:
```typescript
dashboard: {
    widgets: [
        { type: "stat", label: "Total Posts", endpoint: "/api/posts", icon: "FileText" },
        { type: "chart", label: "Posts Over Time", endpoint: "/api/stats/posts", chartType: "line" },
    ],
}
```

---

## Important: Don't Break These

1. **Never delete marker comments** (`// grit:models`, `// grit:routes:protected`, etc.)
2. **Always follow the response format** (`{ data, message }` or `{ data, meta }` or `{ error: { code, message } }`)
3. **Always handle errors in Go** — never ignore with `_` except for deliberate fire-and-forget
4. **Keep the folder structure** — don't move files to non-standard locations
5. **Use React Query for all data fetching** — no raw `fetch` in components
6. **Use Zod for validation** — shared between frontend and backend
7. **Use Tailwind + shadcn/ui for styling** — no custom CSS files
8. **Use App Router** — never Pages Router in Next.js

---

## Mounted Dashboards

The Go API mounts several dashboards automatically:

| Path | Tool | Purpose |
|------|------|---------|
| `/studio` | GORM Studio | Visual database browser — view/edit data, run SQL, export schemas |
| `/docs` | gin-docs | Auto-generated API documentation with interactive UI |
| `/sentinel/ui` | Sentinel | Security dashboard — threats, blocked IPs, rate limits |
| `/pulse` | Pulse | Observability dashboard — request tracing, metrics, errors |

---

## Type Sync

```bash
grit sync
```

Reads Go models in `apps/api/internal/models/` using AST parsing and generates TypeScript types + Zod schemas in `packages/shared/`. Useful when you manually add models without using `grit generate`.

---

## Desktop Projects (Wails)

Grit v2.0 supports building **native desktop applications** using [Wails](https://wails.io/) — Go backend with React frontend compiled into a single executable.

### Creating a Desktop Project

```bash
grit new-desktop myapp
```

This scaffolds:

```
myapp/
├── main.go                    # Wails entry point
├── app.go                     # App struct with bound methods
├── wails.json                 # Wails configuration
├── go.mod                     # Go module
├── .env                       # Environment variables
├── internal/
│   ├── config/config.go       # Config loader (SQLite/Postgres)
│   ├── db/db.go               # GORM connection + AutoMigrate
│   ├── models/                # GORM models (user, blog, contact, types)
│   └── service/               # Service layer (auth, blog, contact, export)
├── cmd/studio/main.go         # Standalone database browser (port 4000)
└── frontend/
    ├── package.json           # React + Vite + TanStack Router + TanStack Query
    ├── src/
    │   ├── index.css          # Grit dark theme
    │   ├── components/
    │   │   ├── layout/        # TitleBar, Sidebar, AppLayout
    │   │   └── ui/            # Button, Input, Card, Select, etc.
    │   ├── routes/            # File-based routes (TanStack Router)
    │   │   ├── __root.tsx     # Root route
    │   │   ├── _layout.tsx    # Auth guard + sidebar layout
    │   │   └── _layout/       # Protected page routes
    │   ├── hooks/             # useAuth, useTheme
    │   └── lib/               # utils, query-client
    └── wailsjs/               # Auto-generated Go bindings (by Wails)
```

### Desktop vs Web Architecture

| Aspect | Web (`grit new`) | Desktop (`grit new-desktop`) |
|--------|-------------------|------------------------------|
| Backend | Gin HTTP server | Wails Go bindings (direct function calls) |
| Frontend | Next.js (App Router) | Vite + React + TanStack Router (file-based) |
| Database | PostgreSQL | SQLite (default) or Postgres |
| Routing | BrowserRouter | Hash history (createHashHistory) |
| API calls | fetch / axios | `wailsjs/go/main/App` imports |
| Structure | Turbo monorepo | Single Go module |
| Auth | JWT tokens | Local session (bcrypt) |
| Deploy | Cloud / Docker | Native .exe / .app / binary |

### Desktop CLI Commands

```bash
grit start         # Auto-detects desktop project, runs `wails dev`
grit compile       # Builds native executable via `wails build`
grit studio        # Opens GORM Studio at http://localhost:4000
grit version       # Prints CLI version (2.0.0)
```

### Desktop Resource Generation

Works the same as web, but generates desktop-specific files:

```bash
grit generate resource Product --fields "name:string,price:float,published:bool"
```

**Files created:**
- `internal/models/product.go` — GORM model struct
- `internal/service/product.go` — Service with List, ListAll, GetByID, Create, Update, Delete
- `frontend/src/routes/_layout/products.index.tsx` — List route with search, pagination, export
- `frontend/src/routes/_layout/products.new.tsx` — Create form route
- `frontend/src/routes/_layout/products.$id.edit.tsx` — Edit form route

**Injections (10 total):**
1. `db.go` — `&models.Product{}` in AutoMigrate (`// grit:models`)
2. `main.go` — Service initialization (`// grit:service-init`)
3. `main.go` — Service passed to NewApp (`/* grit:app-args */`)
4. `app.go` — Service field on App struct (`// grit:fields`)
5. `app.go` — Constructor parameter (`/* grit:constructor-params */`)
6. `app.go` — Constructor assignment (`/* grit:constructor-assign */`)
7. `app.go` — 7 bound methods: Get, GetAll, Create, Update, Delete, ExportPDF, ExportExcel (`// grit:methods`)
8. `types.go` — Input struct (`// grit:input-types`)
9. `cmd/studio/main.go` — Model in Studio AutoMigrate (`// grit:studio-models`)
10. `sidebar.tsx` — Nav icon + nav item (`// grit:nav-icons`, `// grit:nav`)

**Remove:**
```bash
grit remove resource Product
```

### Desktop Markers Reference

These markers are inserted by `grit new-desktop` and used by `grit generate resource`:

| File | Marker | Injection Type |
|------|--------|---------------|
| `internal/db/db.go` | `// grit:models` | Line before |
| `main.go` | `// grit:service-init` | Line before |
| `main.go` | `/* grit:app-args */` | Inline |
| `app.go` | `// grit:imports` | Line before |
| `app.go` | `// grit:fields` | Line before |
| `app.go` | `/* grit:constructor-params */` | Inline |
| `app.go` | `/* grit:constructor-assign */` | Inline |
| `app.go` | `// grit:methods` | Line before |
| `internal/models/types.go` | `// grit:input-types` | Line before |
| `cmd/studio/main.go` | `// grit:studio-models` | Line before |
| `frontend/src/components/layout/sidebar.tsx` | `// grit:nav-icons` | Line before |
| `frontend/src/components/layout/sidebar.tsx` | `// grit:nav` | Line before |

### Project Detection

Grit CLI auto-detects project type:
- **Desktop:** `wails.json` exists at project root
- **Web:** `turbo.json` + `apps/api` directory exist at project root

Commands like `grit start`, `grit studio`, and `grit generate resource` automatically adapt their behavior based on the detected project type.

---

*Built with Grit. Go + React. Built with Grit.*
