---
title: "Full-stack CRUD in one command"
subtitle: "grit generate resource writes the Go model, service, handler, routes, Zod schema, TypeScript types, React hooks, and an admin page — all at once."
series: "The Daily Grit"
edition: 3
date: 2026-07-03
readingTime: "5 min"
author: "Muke JohnBaptist"
tags: [grit, crud, code-generation, golang, react]
canonical: "https://gritframework.dev/blog/full-stack-crud-in-one-command"
---

You have a running app. Now let's add a feature — the whole vertical slice — with
a single command.

## One command, both sides of the stack

```bash
grit generate resource Product --fields "name:string,price:float,stock:int"
```

Grit emits everything a `Product` needs, front to back:

```
✓ apps/api/internal/models/product.go        # GORM model
✓ apps/api/internal/services/product.go       # business logic
✓ apps/api/internal/handlers/product.go       # HTTP handlers
✓ packages/shared/schemas/product.ts          # Zod schema
✓ packages/shared/types/product.ts            # TypeScript types
✓ apps/web/hooks/use-products.ts              # React Query hooks
✓ apps/admin/resources/products.ts            # admin resource + page
✓ Injected model, routes, and registry
```

## The backend it writes

Thin handler, real Go, ownership-scoped query:

```go
func (h *ProductHandler) List(c *gin.Context) {
    var products []models.Product
    h.DB.
        Where("user_id = ?", c.GetString("user_id")).
        Order("created_at desc").
        Find(&products)

    c.JSON(http.StatusOK, gin.H{"data": products})
}
```

Every list endpoint follows Grit's response contract, so clients are predictable:

```json
{ "data": [ ... ], "meta": { "total": 100, "page": 1, "page_size": 20, "pages": 5 } }
```

## The frontend it writes

A typed React Query hook you can drop into any component — no `fetch`, no manual
types:

```ts
export function useProducts() {
  return useQuery({
    queryKey: ['products'],
    queryFn: async () => {
      const res = await api.get('/api/products')
      return res.data.data as Product[]
    },
  })
}
```

And an **admin page** — a sortable, filterable, paginated data table plus a form
builder — generated from the same field definitions. Nothing to wire by hand.

## Field types that cover real apps

`grit generate` understands far more than strings and ints:

`text`, `int`, `float`, `bool`, `date`, `datetime`, `uuid`, `json`, `email`,
`url`, `phone`, `slug`, `file`, `files`, plus relationships (`belongs_to`,
`many_to_many`). Add modifiers inline:

```bash
grit generate resource Post --fields "title:string,slug:string:unique,views:int"
```

Restrict routes to roles when you need to:

```bash
grit generate resource Invoice --fields "number:string,total:float" --roles "ADMIN,ACCOUNTANT"
```

## Keep types honest

Edit a Go model later? Re-sync the frontend types and Zod schemas from the Go
structs — one source of truth:

```bash
grit sync
```

Changed your mind about a resource? `grit remove resource Product` deletes the
files and reverses every injection cleanly.

## Why this matters

The generated code isn't a black box — it's idiomatic Go and React you'd have
written yourself, following one set of conventions. You spend your time on the
20% that's your product, not the 80% that's the same every time.

**Tomorrow:** Grit vs. Next.js — when API routes are enough, and when you want a
real Go backend behind them.

*Go + React. Built with Grit.*
