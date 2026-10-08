# Grit Framework: Detailed Reference

## API Conventions

### Response Format

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

// Error. respond.Fail takes the status from the error catalogue, so a code
// cannot be paired with two different statuses in two different handlers.
respond.Fail(c, respond.CodeValidationError, "Email is required")
```

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

```
POST /api/auth/register  → { access_token, refresh_token }
POST /api/auth/login     → { access_token, refresh_token } or { totp_required, pending_token }
POST /api/auth/refresh   → New access_token from refresh_token
POST /api/auth/logout    → Invalidates refresh token
GET  /api/auth/me        → Current user (requires auth)
```

Access tokens: 15 minutes. Refresh tokens: 7 days.

### Two-Factor Authentication (TOTP)

If user has 2FA enabled and no trusted device cookie, login returns `{ totp_required: true, pending_token: "..." }`.
Client redirects to TOTP page, user enters 6-digit code from authenticator app.

```
POST /api/auth/totp/setup              → { secret, uri } (JWT required)
POST /api/auth/totp/enable             → { enabled, backup_codes } (JWT required)
POST /api/auth/totp/verify             → { user, tokens } (public, uses pending_token)
POST /api/auth/totp/backup-codes/verify → { user, tokens } (public, uses pending_token)
POST /api/auth/totp/disable            → Disable 2FA (JWT + password required)
GET  /api/auth/totp/status             → { enabled, backup_codes_remaining, trusted_devices }
```

TOTP: RFC 6238, HMAC-SHA1, 6 digits, 30s period. Backup codes: 10 bcrypt-hashed one-time codes.
Trusted devices: HttpOnly cookie, SHA-256 hashed token, 30-day sliding expiry.

### Route Groups

```go
public := router.Group("/api/auth")          // No auth
protected := router.Group("/api")            // Requires JWT
protected.Use(middleware.Auth(cfg.JWTSecret))
admin := protected.Group("/admin")           // Requires JWT + admin role
admin.Use(middleware.RequireRole("admin"))
```

---

## Money

A `price:money` field is `money.Money`, never a float:

```go
type Money struct {
    Amount   int64  // MINOR units: 2800 is $28.00
    Currency string
}
```

- `money.New(2800, "USD")` to build one, `money.FromMajor(28.00, "USD")` when the figure came from a person typing it.
- `a.Add(b)` returns `(Money, error)` and errors on mixed currencies.
  Check it even in a one-currency app: the day it has two, an ignored error
  is a silently wrong total.
- `m.MulInt(n)` for a line total. `m.Major()` for display only.

GORM embeds it as two columns with the field name as prefix, so the column
is `price_amount` and NOT `price`. That is what you sort and filter on:

```
GET /api/v1/public/products?sort_by=price_amount&sort_order=asc
GET /api/v1/public/products?price_amount_min=5000
```

On the frontend, `formatMoney(m)` from `@repo/shared/types`. It knows UGX
and JPY have no minor unit; a hardcoded `amount / 100` shows a 50,000
shilling price as 500, and nobody notices until the first Ugandan customer
complains.

## The public API (--public)

`grit generate resource X --public` adds a read-only surface for callers with
no logged-in user: a storefront, a mobile app, a public directory.

```
GET /api/v1/public/<plural>               list, search, sort, filter
GET /api/v1/public/<plural>/:key          by slug, or by id when there is none
GET /api/v1/public/<plural>/:key/related  same collection, newest first
```

- Guarded by a publishable API key. Not secrecy, since the key ships inside
  your app: identification, a rate-limit bucket per key, per-endpoint and
  per-origin narrowing, and the ability to turn one client off without a
  deploy.
- The response is an allowlist in `internal/handlers/<x>_public.go`, not the
  model. That file is NEVER overwritten when the resource is regenerated, so
  a column you want published goes in the struct and in `toPublicX` by hand.
- Relations are held back on purpose: publishing one publishes a whole
  related record nobody vetted. Filter by the foreign key id instead, which
  the handler already allows.

### A public endpoint that writes

There is no generator for one. When you write it (a cart, a sign-up, a
review), three rules:

- Do NOT mount it on the `publicAPI` group. That group has a response cache
  keyed on the URL and nothing else, so one caller's answer is served to
  every caller. Use your own group with the same `middleware.RequireAPIKey`.
- Set `Cache-Control: private` on the response anyway, so moving the route
  later cannot make it wrong.
- Never take a price, a total or a discount from the request body. Take an id
  and a quantity, and read every figure from the database.

## Variants

`grit add variants --resource Product` installs options, option values, the
per-product offered options, the variants and the join between them, plus a
matrix editor on the product's admin page and one public endpoint:

```
GET /api/v1/public/products/:key/variants
```

One response, with everything a picker needs: the options to draw (values,
swatches, price deltas), every combination (id, sku, price, in_stock,
option_value_ids), and the price range.

- One request, not three. Fetching options, then variants, then a price per
  swatch click is three round trips per interaction on your busiest page.
- The price in it is the server's. Do not re-apply option deltas in the
  browser: a second implementation of that arithmetic is a second answer to
  "what does this cost".
- A product with no variants gets empty lists and a range of its own price,
  so one component renders both cases. Most of a real catalogue has no
  options at all.
- Match a selection to a variant by comparing the set of `option_value_ids`.
  The values are not nested on each variant because they are already in the
  options list, and nesting them would send the same objects twice.

## Go Model Pattern

```go
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
```

Rules:
- Always include ID, CreatedAt, UpdatedAt, DeletedAt
- `json:"-"` for DeletedAt (hidden from API)
- `binding:"required"` for required fields

---

## Admin Panel — Resource Definitions

```typescript
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
```

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

### Entering Many Rows

Do not build a screen for this. Every list page already has it:

- **Bulk Create** opens a grid of blank rows, created in one transaction.
- **Bulk Edit** opens the selected rows as a grid and sends only changed cells.

Both come from the resource definition, need no code, and follow the create and
update permissions. Hide the button with `table: { bulkCreate: false }`.

The grid holds a cell per field, so `richtext`, `line-items`,
`json` and `multi-relationship-select` stay in the form. The API
validates each row against the model's `binding:` tags, so a required
column left empty is refused per row rather than written.

---

## Frontend Patterns

### React Query Hooks

```typescript
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
```

---

## Batteries (Services)

### File Storage

```go
storage.Upload(ctx, "uploads/2024/01/photo.jpg", reader, "image/jpeg")
url := storage.GetURL("uploads/2024/01/photo.jpg")
url, err := storage.GetSignedURL(ctx, key, 1*time.Hour)
```

### Email

MAIL_MAILER picks the driver (smtp, resend, mailgun, postmark, sendgrid, ses, log, failover); development sends to Mailhog.

```go
// A built-in template, sent now
mailer.Send(ctx, mail.SendOptions{
    To: "user@example.com", Subject: "Welcome!",
    Template: "welcome", Data: map[string]interface{}{"Name": "John"},
})

// Any message: several recipients, cc, bcc, a text part, attachments
mailer.SendMessage(ctx, &mail.Message{To: []string{"a@example.com"}, Subject: "Hi", HTML: "<p>Hi</p>", Text: "Hi"})

// Queued: the worker sends it and retries. Attachments over 256 KB are refused.
mail.Queue(ctx, svc.Jobs, &mail.Message{To: []string{"a@example.com"}, Subject: "Hi", HTML: "<p>Hi</p>"})
```

Templates: `welcome`, `password-reset`, `email-verification`, `notification`. Add your own with `grit generate mail OrderShipped` (internal/mail/templates/order_shipped.go: typed OrderShippedData, SendOrderShipped, QueueOrderShipped, listed in the admin Mail Preview). In tests, `mailtest.New()` records mail: `fake.AssertSent(t, to, subject)`.

### Background Jobs

```go
// Fire-and-forget enqueue (uses framework defaults: 5 max retries,
// exponential backoff 1s/2s/4s.../5min cap, 5min per-attempt timeout).
svc.Jobs.EnqueueSendEmail(ctx, "user@example.com", "Welcome", "welcome", data)

// Idempotency key — a retry of the same operation is deduped within
// the 24h window. Critical for charge/notify/order-confirmation work.
svc.Jobs.EnqueueSendEmail(ctx, user.Email, "Your account", "notification", data, jobs.EnqueueOption{
    IdempotencyKey: "account-notice:" + user.ID,
})

svc.Jobs.EnqueueProcessImage(ctx, uploadID, key, mimeType)
```

### Redis Cache

```go
cache.Set(ctx, "user:123", userData, 5*time.Minute)
cache.Get(ctx, "user:123", &user)
cache.Delete(ctx, "user:123")
```

### AI Integration (Vercel AI Gateway)

```go
result, err := ai.Complete(ctx, ai.CompletionRequest{Prompt: "Summarize..."})
ai.Stream(ctx, req, func(chunk string) { /* SSE */ })
```

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
