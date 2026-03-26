package scaffold

import "fmt"

// gritSkillFile generates the .claude/skills/grit/SKILL.md for a scaffolded project.
// Follows the Claude Code skills format: https://code.claude.com/docs/en/skills
func gritSkillFile(opts Options) string {
	return fmt.Sprintf(`---
name: grit
description: >
  Grit framework conventions and patterns for Go + React full-stack monorepo projects.
  Use when modifying models, handlers, routes, schemas, types, resources, or admin panel
  components in a Grit project. Automatically loaded as background knowledge.
user-invocable: false
---

# Grit Framework

Grit is a full-stack meta-framework: **Go** (Gin + GORM + PostgreSQL) + **React/Next.js** (App Router + Tailwind + shadcn/ui) in a monorepo. Think Laravel + Filament, but Go + React.

**Batteries included:** file storage (S3), email (Resend), background jobs (asynq), cron, Redis cache, AI (Claude/OpenAI), security (Sentinel), observability (Pulse), auto-generated API docs (gin-docs).

For detailed API conventions, code patterns, and service documentation, see [reference.md](reference.md).

---

## CLI Commands

%[1]sash
grit new myapp                    # Full monorepo (Go API + Next.js web + admin)
grit new myapp --api              # Go API only
grit new myapp --full             # Everything + Expo mobile + docs site

grit generate resource Post --fields "title:string,content:text,published:bool"
grit generate resource Post --from post.yaml
grit generate resource Category -i   # Interactive mode

grit sync                             # Go types → TypeScript
grit add role MODERATOR               # Injects role into 7 locations
grit migrate                          # Run GORM AutoMigrate
grit seed                             # Create admin + demo users
grit upgrade                          # Update project to latest templates
grit update                           # Update Grit CLI itself
%[1]s

---

## Project Structure

%[1]s
%[2]s/
├── .env                          # Environment variables
├── docker-compose.yml            # PostgreSQL, Redis, MinIO, Mailhog
├── .claude/skills/grit/          # This skill — AI assistant guide
├── packages/shared/              # Zod schemas, TS types, constants
│   ├── schemas/                  # Zod validation (user.ts, etc.)
│   ├── types/                    # TypeScript interfaces
│   └── constants/                # API_ROUTES, ROLES, etc.
├── apps/
│   ├── api/                      # Go backend (Gin + GORM)
│   │   ├── cmd/server/main.go
│   │   └── internal/
│   │       ├── config/           # Loads .env
│   │       ├── database/         # GORM connection
│   │       ├── models/           # GORM models
│   │       ├── handlers/         # HTTP handlers
│   │       ├── services/         # Business logic
│   │       ├── middleware/       # Auth, CORS, logger, cache
│   │       ├── routes/           # Route registration
│   │       ├── mail/             # Email (Resend)
│   │       ├── storage/          # File storage (S3)
│   │       ├── jobs/             # Background jobs (asynq)
│   │       ├── cron/             # Scheduled tasks
│   │       ├── cache/            # Redis cache
│   │       └── ai/               # AI service
│   ├── web/                      # SaaS landing page (Next.js)
│   └── admin/                    # Filament-like admin panel
│       ├── components/           # Layout, tables, forms, widgets
│       ├── hooks/                # use-auth, use-resource, use-system
│       ├── resources/            # Resource definitions
│       └── lib/                  # defineResource(), icons
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

Creates **8 files** (model, service, handler, schema, types, hooks, resource def, admin page) and injects into **6 existing files** (models, routes, schemas, types, constants, resource registry) via marker comments.

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

## Marker Comments

Grit uses marker comments to inject generated code. **Never delete these:**

%[1]sgo
// grit:models          — models/user.go (AutoMigrate list)
// grit:handlers        — routes/routes.go (handler initialization)
// grit:routes:protected — routes/routes.go (protected route group)
// grit:routes:admin    — routes/routes.go (admin route group)
%[1]s

%[1]stypescript
// grit:schemas         — schemas/index.ts
// grit:types           — types/index.ts
// grit:api-routes      — constants/index.ts
// grit:resources       — resources/index.ts (imports)
// grit:resource-list   — resources/index.ts (registry array)
%[1]s

---

## Common Tasks

### Add a field to an existing resource

1. Add field to Go model (%[1]sapps/api/internal/models/<name>.go%[1]s)
2. Update handler if field needs special handling
3. Update Zod schema (%[1]spackages/shared/schemas/<name>.ts%[1]s)
4. Update TypeScript type (%[1]spackages/shared/types/<name>.ts%[1]s)
5. Update admin resource (%[1]sapps/admin/resources/<name>.ts%[1]s) — add column + form field
6. Restart API (GORM auto-migrates)

### Add a new API endpoint

1. Create/update handler in %[1]sapps/api/internal/handlers/%[1]s
2. Register route in %[1]sapps/api/internal/routes/routes.go%[1]s
3. Create React Query hook in %[1]sapps/web/hooks/%[1]s or %[1]sapps/admin/hooks/%[1]s

### Add a relationship

%[1]sgo
type Post struct {
    CategoryID uint     %[1]sjson:"category_id"%[1]s
    Category   Category %[1]sgorm:"foreignKey:CategoryID" json:"category,omitempty"%[1]s
}
// In handler: query.Preload("Category").Find(&posts)
%[1]s

---

## Critical Rules

1. **Never delete marker comments** (%[1]s// grit:*%[1]s)
2. **Follow the response format** — %[1]s{ data, message }%[1]s / %[1]s{ data, meta }%[1]s / %[1]s{ error: { code, message } }%[1]s
3. **Always handle errors in Go** — never ignore with %[1]s_%[1]s
4. **Keep the folder structure** — don't move files
5. **Use React Query** for all data fetching — no raw %[1]sfetch%[1]s
6. **Use Zod** for validation — shared between frontend and backend
7. **Use Tailwind + shadcn/ui** — no custom CSS files
8. **Use App Router** — never Pages Router
`, "`", opts.ProjectName)
}

// gritSkillReference generates .claude/skills/grit/reference.md with detailed
// API conventions, code patterns, admin panel docs, and service documentation.
func gritSkillReference(opts Options) string {
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

Access tokens: 15 minutes. Refresh tokens: 7 days. Auto-refresh via Axios interceptor.

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
POST /api/auth/totp/backup-codes       → Regenerate backup codes (JWT required)
DELETE /api/auth/totp/trusted-devices   → Revoke all trusted devices (JWT required)
%[1]s

TOTP: RFC 6238, HMAC-SHA1, 6 digits, 30s period, ±1 window. Backup codes: 10 bcrypt-hashed one-time codes.
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

Rules:
- Always include ID, CreatedAt, UpdatedAt, DeletedAt
- %[1]sjson:"-"%[1]s for DeletedAt (hidden from API)
- %[1]sbinding:"required"%[1]s for required fields
- %[1]sgorm:"size:255"%[1]s for strings, %[1]sgorm:"type:text"%[1]s for long text

---

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

  dashboard: {
    widgets: [{
      type: "stat", label: "Total Posts", endpoint: "/api/posts?page_size=1",
      icon: "FileText", color: "accent", format: "number", colSpan: 1,
    }],
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
| %[1]sradio%[1]s | Radio group | requires options |
| %[1]simage%[1]s | Image upload | react-dropzone |
| %[1]srichtext%[1]s | Tiptap WYSIWYG | |
| %[1]srelationship-select%[1]s | Searchable dropdown | belongs_to |
| %[1]smulti-relationship-select%[1]s | Multi-select tags | many_to_many |

### Form View Variants

%[1]sformView%[1]s: %[1]smodal%[1]s (default), %[1]spage%[1]s, %[1]smodal-steps%[1]s, %[1]spage-steps%[1]s

### Dropzone

%[1]stsx
<Dropzone
  variant="default"  // "default" | "compact" | "minimal" | "avatar" | "inline"
  maxFiles={5}
  maxSize={1024 * 1024 * 10}
  onFilesChange={setFiles}
  accept={{ "image/*": [] }}
/>
%[1]s

---

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

## Batteries (Services)

### File Storage

%[1]sgo
storage.Upload(ctx, "uploads/2024/01/photo.jpg", reader, "image/jpeg")
url := storage.GetURL("uploads/2024/01/photo.jpg")
url, err := storage.GetSignedURL(ctx, key, 1*time.Hour)
%[1]s

Endpoint: %[1]sPOST /api/uploads%[1]s (multipart, max 10MB)

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
// As middleware:
router.GET("/api/posts", middleware.CacheResponse(cache, 5*time.Minute), handler.List)
%[1]s

### AI Integration

%[1]sgo
result, err := ai.Complete(ctx, ai.CompletionRequest{Prompt: "Summarize..."})
ai.Stream(ctx, req, func(chunk string) { /* SSE */ })
%[1]s

### Security (Sentinel)

WAF, rate limiting, brute-force protection, anomaly detection. Dashboard at %[1]s/sentinel/ui%[1]s. Disable: %[1]sSENTINEL_ENABLED=false%[1]s.

### Observability (Pulse)

Request tracing, DB monitoring, runtime metrics, error tracking, health checks, Prometheus at %[1]s/pulse/metrics%[1]s. Dashboard at %[1]s/pulse%[1]s. Disable: %[1]sPULSE_ENABLED=false%[1]s.

### API Documentation (gin-docs)

Zero-annotation OpenAPI 3.1 spec. Interactive UI at %[1]s/docs%[1]s. Export to Postman/Insomnia.

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
AI_GATEWAY_API_KEY=your-key
AI_GATEWAY_MODEL=anthropic/claude-sonnet-4-6
SENTINEL_ENABLED=true
PULSE_ENABLED=true
%[1]s

## Docker Services

| Service | Port | Purpose |
|---------|------|---------|
| PostgreSQL | 5432 | Database |
| Redis | 6379 | Cache + job queue |
| MinIO | 9000/9001 | Local S3 storage |
| Mailhog | 1025/8025 | Email testing |

## Naming Conventions

| Thing | Convention | Example |
|-------|-----------|---------|
| Go files | %[1]ssnake_case.go%[1]s | %[1]suser_handler.go%[1]s |
| Go structs | %[1]sPascalCase%[1]s | %[1]stype PostHandler struct%[1]s |
| TS files | %[1]skebab-case.ts%[1]s | %[1]suse-posts.ts%[1]s |
| React | %[1]sPascalCase.tsx%[1]s | %[1]sDataTable.tsx%[1]s |
| API routes | %[1]s/api/plural%[1]s | %[1]s/api/posts%[1]s |
| DB tables | %[1]splural_snake%[1]s | %[1]sblog_posts%[1]s |

## Design System

| Token | Value | Usage |
|-------|-------|-------|
| %[1]sbg-primary%[1]s | %[1]s#0a0a0f%[1]s | Page background |
| %[1]sbg-secondary%[1]s | %[1]s#111118%[1]s | Card background |
| %[1]saccent%[1]s | %[1]s#6c5ce7%[1]s | Primary accent (purple) |
| %[1]ssuccess%[1]s | %[1]s#00b894%[1]s | Success states |
| %[1]sdanger%[1]s | %[1]s#ff6b6b%[1]s | Error states |

Fonts: **Onest** (UI), **JetBrains Mono** (code)
`, "`", opts.ProjectName)
}
