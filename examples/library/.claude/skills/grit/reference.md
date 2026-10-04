# Grit Framework: Detailed Reference

## API Conventions

### Response Format

`go
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

// Error. respond.Fail takes the status from the error catalogue, so a code
// cannot be paired with two different statuses in two different handlers.
respond.Fail(c, respond.CodeValidationError, "Email is required")
`

### Error Codes

| Code | HTTP Status | When |
|------|------------|------|
| `VALIDATION_ERROR` | 422 | Invalid input |
| `NOT_FOUND` | 404 | Resource missing |
| `UNAUTHORIZED` | 401 | Invalid JWT |
| `FORBIDDEN` | 403 | Insufficient role |
| `INTERNAL_ERROR` | 500 | Server error |
| `CONFLICT` | 409 | Duplicate key |

### Authentication

`
POST /api/auth/register  → { access_token, refresh_token }
POST /api/auth/login     → { access_token, refresh_token } or { totp_required, pending_token }
POST /api/auth/refresh   → New access_token from refresh_token
POST /api/auth/logout    → Invalidates refresh token
GET  /api/auth/me        → Current user (requires auth)
`

Access tokens: 15 minutes. Refresh tokens: 7 days.

### Two-Factor Authentication (TOTP)

If user has 2FA enabled and no trusted device cookie, login returns `{ totp_required: true, pending_token: "..." }`.
Client redirects to TOTP page, user enters 6-digit code from authenticator app.

`
POST /api/auth/totp/setup              → { secret, uri } (JWT required)
POST /api/auth/totp/enable             → { enabled, backup_codes } (JWT required)
POST /api/auth/totp/verify             → { user, tokens } (public, uses pending_token)
POST /api/auth/totp/backup-codes/verify → { user, tokens } (public, uses pending_token)
POST /api/auth/totp/disable            → Disable 2FA (JWT + password required)
GET  /api/auth/totp/status             → { enabled, backup_codes_remaining, trusted_devices }
`

TOTP: RFC 6238, HMAC-SHA1, 6 digits, 30s period. Backup codes: 10 bcrypt-hashed one-time codes.
Trusted devices: HttpOnly cookie, SHA-256 hashed token, 30-day sliding expiry.

### Route Groups

`go
public := router.Group("/api/auth")          // No auth
protected := router.Group("/api")            // Requires JWT
protected.Use(middleware.Auth(cfg.JWTSecret))
admin := protected.Group("/admin")           // Requires JWT + admin role
admin.Use(middleware.RequireRole("admin"))
`

---

## Go Model Pattern

`go
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
`

Rules:
- Always include ID, CreatedAt, UpdatedAt, DeletedAt
- `json:"-"` for DeletedAt (hidden from API)
- `binding:"required"` for required fields

---

## Admin Panel — Resource Definitions

`typescript
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
`

### Form Field Types

| Type | Component | Notes |
|------|----------|-------|
| `text` | Text input | prefix, suffix |
| `textarea` | Textarea | configurable rows |
| `number` | Number input | min, max, step |
| `select` | Dropdown | requires options |
| `date` / `datetime` | Picker | |
| `toggle` / `checkbox` | Boolean | |
| `image` | Image upload | react-dropzone |
| `richtext` | Tiptap WYSIWYG | |
| `relationship-select` | Searchable dropdown | belongs_to |
| `multi-relationship-select` | Multi-select tags | many_to_many |

### Form View Variants

`formView`: `modal` (default), `page`, `modal-steps`, `page-steps`

---

## Frontend Patterns

### React Query Hooks

`typescript
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
`

---

## Batteries (Services)

### File Storage

`go
storage.Upload(ctx, "uploads/2024/01/photo.jpg", reader, "image/jpeg")
url := storage.GetURL("uploads/2024/01/photo.jpg")
url, err := storage.GetSignedURL(ctx, key, 1*time.Hour)
`

### Email

MAIL_MAILER picks the driver (smtp, resend, mailgun, postmark, sendgrid, ses, log, failover); development sends to Mailhog.

`go
// A built-in template, sent now
mailer.Send(ctx, mail.SendOptions{
    To: "user@example.com", Subject: "Welcome!",
    Template: "welcome", Data: map[string]interface{}{"Name": "John"},
})

// Any message: several recipients, cc, bcc, a text part, attachments
mailer.SendMessage(ctx, &mail.Message{To: []string{"a@example.com"}, Subject: "Hi", HTML: "<p>Hi</p>", Text: "Hi"})

// Queued: the worker sends it and retries. Attachments over 256 KB are refused.
mail.Queue(ctx, svc.Jobs, &mail.Message{To: []string{"a@example.com"}, Subject: "Hi", HTML: "<p>Hi</p>"})
`

Templates: `welcome`, `password-reset`, `email-verification`, `notification`. Add your own with `grit generate mail OrderShipped` (internal/mail/templates/order_shipped.go: typed OrderShippedData, SendOrderShipped, QueueOrderShipped, listed in the admin Mail Preview). In tests, `mailtest.New()` records mail: `fake.AssertSent(t, to, subject)`.

### Background Jobs

`go
// Fire-and-forget enqueue (uses framework defaults: 5 max retries,
// exponential backoff 1s/2s/4s.../5min cap, 5min per-attempt timeout).
svc.Jobs.EnqueueSendEmail(ctx, "user@example.com", "Welcome", "welcome", data)

// Idempotency key — a retry of the same operation is deduped within
// the 24h window. Critical for charge/notify/order-confirmation work.
svc.Jobs.EnqueueSendEmail(ctx, user.Email, "Your account", "notification", data, jobs.EnqueueOption{
    IdempotencyKey: "account-notice:" + user.ID,
})

svc.Jobs.EnqueueProcessImage(ctx, uploadID, key, mimeType)
`

### Redis Cache

`go
cache.Set(ctx, "user:123", userData, 5*time.Minute)
cache.Get(ctx, "user:123", &user)
cache.Delete(ctx, "user:123")
`

### AI Integration (Vercel AI Gateway)

`go
result, err := ai.Complete(ctx, ai.CompletionRequest{Prompt: "Summarize..."})
ai.Stream(ctx, req, func(chunk string) { /* SSE */ })
`

One key, hundreds of models. Config: `AI_GATEWAY_API_KEY`, `AI_GATEWAY_MODEL` (e.g. `anthropic/claude-sonnet-4-6`).

### Security (Sentinel)

WAF, rate limiting, brute-force protection. Dashboard at `/sentinel/ui`.

### Observability (Pulse)

Request tracing, DB monitoring, metrics. Dashboard at `/pulse`. Prometheus at `/pulse/metrics`.

### API Documentation (gin-docs)

Zero-annotation OpenAPI 3.1 spec. Interactive UI at `/docs`.

---

## Naming Conventions

| Thing | Convention | Example |
|-------|-----------|---------|
| Go files | `snake_case.go` | `user_handler.go` |
| Go structs | `PascalCase` | `type PostHandler struct` |
| TS files | `kebab-case.ts` | `use-posts.ts` |
| React | `PascalCase.tsx` | `DataTable.tsx` |
| API routes | `/api/plural` | `/api/posts` |
| DB tables | `plural_snake` | `blog_posts` |
