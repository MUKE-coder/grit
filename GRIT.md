# GRIT — Go + React. Built with Grit.

## What is Grit?

Grit is a **full-stack meta-framework** that fuses a **Go backend** (Gin + GORM) with a **Next.js + React frontend** in a single monorepo. It ships with batteries included — authentication, admin panel generator, file storage, email, background jobs, and AI — all wired together out of the box.

Think of it as **Laravel's developer experience, but with Go's performance and React's frontend ecosystem.**

One command to start:
```bash
grit new myapp
```

One command to generate full-stack resources:
```bash
grit generate resource Invoice
# Creates: Go model, API handler, React Query hook, Zod schema, admin page, form
```

One command to deploy:
```bash
grit deploy
```

---

## The Problem

### The JavaScript Ecosystem is Fragmented

Building a modern full-stack app in 2025/2026 requires stitching together 15+ tools: Next.js, Prisma, NextAuth, tRPC, Zod, React Query, Tailwind, S3 SDK, Resend, Redis client, Bull queues, admin panel library... Every new project starts with 2-3 days of boilerplate. Every developer makes different choices, creating inconsistent codebases that are hard to maintain.

### Go Lacks a Batteries-Included Framework

Go has incredible performance, simplicity, and deployment story. But the Go web ecosystem is fragmented — Gin, Echo, Fiber, Chi for routing; GORM, sqlc, sqlx for databases; dozens of auth libraries. There's no equivalent to Laravel or Rails that gives you everything wired together. Go developers spend more time plumbing than building.

### No Admin Panel Generator for Go + React

Laravel has **Filament** — a tool that lets you generate beautiful admin dashboards, data tables, forms, and widgets from simple PHP definitions. The Go + React world has nothing comparable. Every team builds their admin panel from scratch, wasting weeks on repetitive CRUD interfaces.

### The Gap Between Backend and Frontend

In most full-stack apps, the backend and frontend are disconnected. Change a Go struct? Now manually update your TypeScript types. Add pagination to an API? Now wire it up separately in React. There's no framework that keeps both sides in sync automatically.

### Serverless Isn't Always the Answer

Next.js pushes developers toward serverless (Vercel, AWS Lambda). For many applications — CRMs, SaaS tools, internal dashboards, apps with WebSockets, background jobs, or heavy processing — a real backend server is superior. Go gives you a single binary that handles 100k+ concurrent connections, deploys anywhere, and costs a fraction of serverless.

---

## Who is Grit For?

### Primary Audience

- **Freelancers and agencies** who take client projects and need to ship fast. Grit gives them a complete stack that works out of the box — auth, admin panel, file uploads, email — without spending days on setup.

- **Solo SaaS developers** who want Go's performance for their backend but React's ecosystem for their frontend. Grit eliminates the glue code between the two.

- **Small to mid-size teams (2-15 developers)** building internal tools, CRMs, dashboards, or SaaS products. The monorepo structure and code generation scale well with team size.

### Secondary Audience

- **Laravel developers** who want to move to Go but miss the DX. Grit is the closest thing to Laravel's philosophy in the Go ecosystem.

- **Next.js developers** who are frustrated with serverless limitations and want a real backend. Grit gives them Go without forcing them to leave React.

- **Agencies building client admin panels** — the Filament-like tool alone is worth adopting Grit for.

### Who Grit is NOT For

- **Large enterprises** with established stacks and dedicated platform teams. They don't need opinionated frameworks.
- **Developers who only want backend or only frontend.** Grit's value is in the full-stack integration.
- **Minimalists** who prefer picking every library themselves. Grit is opinionated by design.

---

## Core Philosophy

### 1. Convention Over Configuration
There is ONE way to do things in Grit. One auth system. One state management approach. One folder structure. Opinions are features. This means any developer can jump into any Grit project and immediately understand it.

### 2. Code Generation Over Runtime Magic
Instead of complex runtime coupling between Go and React, Grit uses a CLI code generator. `grit generate resource Post` creates all the files — Go model, handler, React hook, Zod schema, admin page. The generated code is readable, editable, and debuggable. No black boxes.

### 3. Batteries Included, Optionally Removable
Auth, file storage, email, queues, cron, AI, and admin panel ship with every project. But they're modular — you can remove what you don't need. The framework should be great for a freelancer building a simple app AND for a team building a complex SaaS.

### 4. Beautiful by Default
Every UI component — the admin panel, data tables, forms, login pages — ships with a polished dark theme inspired by GORM Studio. First impressions matter. When a developer runs `grit new` and sees the result, they should think "this looks like a real product."

### 5. Monorepo Native
Go backend, Next.js frontend, shared schemas, and (optionally) Expo mobile app — all in one repository. Shared types, shared validation, shared constants. One PR to change the full stack. This is how modern teams work.

### 6. Vibe Coding Ready
Grit's strict conventions and predictable structure make it ideal for AI-assisted development. An AI can understand the entire project structure, generate resources, and modify code confidently because Grit follows the same patterns everywhere.

---

## Features

### Framework Core

| Feature | Description |
|---------|-------------|
| **CLI Scaffolder** | `grit new myapp` creates a complete monorepo with Go API, Next.js frontend, admin panel, Docker setup, and shared packages. Flags: `--api` (Go only), `--full` (includes Expo mobile app) |
| **Code Generator** | `grit generate resource <Name>` creates Go model + migration, API handler with full CRUD, React Query hooks, Zod validation schema, admin panel page with table + form |
| **Monorepo Structure** | Turborepo-powered monorepo with `apps/api` (Go), `apps/web` (Next.js), `apps/admin` (Next.js admin panel), `packages/shared` (Zod schemas, TypeScript types) |
| **Docker Setup** | Production-ready `docker-compose.yml` with Go API, Next.js, PostgreSQL, Redis, MinIO (S3-compatible), and Mailhog for development |
| **Type Sync** | Go structs auto-generate TypeScript types and Zod schemas. Change the Go model, run `grit sync`, and the frontend types update automatically |

### Authentication & Authorization

| Feature | Description |
|---------|-------------|
| **JWT Auth** | Pre-configured JWT authentication with access + refresh tokens |
| **Role-Based Access** | Built-in role system (admin, editor, user, custom roles) with middleware |
| **Login/Register Pages** | Pre-built, themed auth pages in the frontend |
| **Password Reset** | Email-based password reset flow wired to Resend |
| **Session Management** | Token refresh, logout, device tracking |
| **OAuth Ready** | Extensible auth system for adding Google, GitHub, etc. |

### Admin Panel (Filament-like)

This is the flagship feature — a Filament-style admin panel generator for Go + React.

| Feature | Description |
|---------|-------------|
| **Resource Definitions** | Define a resource in TypeScript — columns, form fields, filters, actions — and get a full admin page |
| **DataTable** | Server-side paginated, sortable, filterable, searchable data tables with column visibility toggle |
| **Form Builder** | Auto-generated forms from resource definitions — text, number, select, date, file upload, rich text, JSON editor, relations |
| **Dashboard Widgets** | Stats cards, charts (line, bar, pie), recent activity, custom widgets |
| **Sidebar Navigation** | Auto-generated sidebar from registered resources |
| **Bulk Actions** | Select multiple rows → delete, export, custom actions |
| **Filters** | Column filters, date range, select filters, search |
| **Dark Theme** | Beautiful dark theme consistent with GORM Studio aesthetic |
| **CRM Feel** | The admin panel should feel like a premium CRM — not a boring CRUD generator |

### File Storage

| Feature | Description |
|---------|-------------|
| **S3/R2 Abstraction** | Unified API for AWS S3, Cloudflare R2, MinIO (local dev) |
| **File Upload Component** | Drag-and-drop file upload React component wired to the Go backend |
| **Image Processing** | Thumbnail generation, resize on upload |
| **Signed URLs** | Pre-signed URLs for secure file access |

### Email

| Feature | Description |
|---------|-------------|
| **Resend Integration** | Pre-configured email sending via Resend |
| **Email Templates** | Go template-based HTML email templates |
| **Transactional Emails** | Welcome, password reset, notification emails ready to use |

### Background Jobs & Scheduling

| Feature | Description |
|---------|-------------|
| **Job Queue** | Redis-backed job queue for background processing |
| **Cron Jobs** | Built-in cron scheduler for recurring tasks |
| **Job Dashboard** | View queued, processing, failed, and completed jobs in the admin panel |

### AI Integration

| Feature | Description |
|---------|-------------|
| **AI SDK** | Pre-wired AI SDK (Vercel AI SDK on frontend, Go client on backend) |
| **Chat Component** | Ready-to-use chat UI component |
| **Structured Output** | AI responses with Zod schema validation |

### Database

| Feature | Description |
|---------|-------------|
| **GORM** | Pre-configured GORM with PostgreSQL (production) and SQLite (quick start) |
| **Migrations** | `grit migrate` runs GORM auto-migrations |
| **GORM Studio** | Embedded visual database browser at `/studio` |
| **Seeding** | Database seeder for development data |

### Developer Experience

| Feature | Description |
|---------|-------------|
| **Hot Reload** | Go backend (Air) + Next.js hot reload in development |
| **Type Safety** | End-to-end type safety from Go structs → TypeScript types → Zod schemas → React components |
| **Error Handling** | Consistent error format across the entire stack |
| **Logging** | Structured JSON logging with request tracing |
| **Environment Management** | `.env` files with validation, separate dev/staging/prod configs |

---

## Tech Stack

| Layer | Technology | Why |
|-------|-----------|-----|
| **Backend** | Go + Gin + GORM | Performance, simplicity, single binary deployment, strong typing |
| **Frontend** | Next.js 14+ (App Router) | React ecosystem, SSR/SSG, file-based routing, huge community |
| **Admin Panel** | Next.js + Custom Components | Same stack as the main frontend, consistent DX |
| **Validation** | Zod (frontend) + Go struct tags (backend) | Shared validation logic, type-safe |
| **Data Fetching** | React Query (TanStack Query) | Caching, background refetch, optimistic updates, server state management |
| **Styling** | Tailwind CSS + shadcn/ui | Utility-first, consistent, customizable, dark mode built-in |
| **Database** | PostgreSQL (prod) + SQLite (dev) | Battle-tested, scalable, GORM supports both |
| **Cache** | Redis | Session storage, job queues, caching |
| **File Storage** | S3 / Cloudflare R2 / MinIO | Universal object storage with local dev support |
| **Email** | Resend | Simple API, great DX, reliable delivery |
| **Monorepo** | Turborepo | Fast builds, dependency caching, task orchestration |
| **Package Manager** | pnpm | Fast, disk-efficient, strict dependency resolution |
| **Containerization** | Docker + Docker Compose | Consistent dev/prod environments |
| **DB Browser** | GORM Studio | Visual database management (built by us!) |

---

## Folder Structure

```
myapp/
├── grit.config.ts                  # Framework configuration
├── docker-compose.yml              # Full dev environment
├── docker-compose.prod.yml         # Production setup
├── turbo.json                      # Monorepo task config
├── pnpm-workspace.yaml             # Workspace definition
├── .env                            # Environment variables
├── .env.example                    # Template
│
├── packages/
│   └── shared/                     # Shared between all apps
│       ├── package.json
│       ├── schemas/                # Zod schemas (source of truth for validation)
│       │   ├── user.ts
│       │   ├── post.ts
│       │   └── index.ts
│       ├── types/                  # TypeScript types (auto-generated from Go)
│       │   ├── user.ts
│       │   ├── post.ts
│       │   └── index.ts
│       └── constants/              # Shared constants (roles, permissions, config)
│           └── index.ts
│
├── apps/
│   ├── api/                        # Go backend
│   │   ├── go.mod
│   │   ├── go.sum
│   │   ├── cmd/
│   │   │   └── server/
│   │   │       └── main.go         # Entry point
│   │   ├── internal/
│   │   │   ├── config/             # App configuration
│   │   │   │   └── config.go
│   │   │   ├── database/           # DB connection + migrations
│   │   │   │   └── database.go
│   │   │   ├── models/             # GORM models
│   │   │   │   ├── user.go
│   │   │   │   └── post.go
│   │   │   ├── handlers/           # HTTP handlers (controllers)
│   │   │   │   ├── auth.go
│   │   │   │   ├── user.go
│   │   │   │   └── post.go
│   │   │   ├── middleware/         # Auth, CORS, logging, rate limit
│   │   │   │   ├── auth.go
│   │   │   │   ├── cors.go
│   │   │   │   └── logger.go
│   │   │   ├── services/           # Business logic
│   │   │   │   ├── auth.go
│   │   │   │   └── user.go
│   │   │   ├── routes/             # Route registration
│   │   │   │   └── routes.go
│   │   │   ├── mail/               # Email templates + sending
│   │   │   │   ├── mailer.go
│   │   │   │   └── templates/
│   │   │   ├── storage/            # File storage abstraction
│   │   │   │   └── storage.go
│   │   │   ├── jobs/               # Background jobs
│   │   │   │   └── jobs.go
│   │   │   └── cron/               # Scheduled tasks
│   │   │       └── cron.go
│   │   └── studio/                 # GORM Studio (embedded)
│   │
│   ├── web/                        # Next.js main frontend
│   │   ├── package.json
│   │   ├── next.config.ts
│   │   ├── tailwind.config.ts
│   │   ├── app/
│   │   │   ├── layout.tsx          # Root layout with providers
│   │   │   ├── page.tsx            # Landing page
│   │   │   ├── (auth)/             # Auth route group
│   │   │   │   ├── login/page.tsx
│   │   │   │   ├── register/page.tsx
│   │   │   │   └── forgot-password/page.tsx
│   │   │   └── (dashboard)/        # Protected routes
│   │   │       ├── layout.tsx
│   │   │       └── dashboard/page.tsx
│   │   ├── components/
│   │   │   ├── ui/                 # shadcn/ui components
│   │   │   └── shared/             # App-specific components
│   │   ├── hooks/                  # React Query hooks (auto-generated)
│   │   │   ├── use-auth.ts
│   │   │   ├── use-users.ts
│   │   │   └── use-posts.ts
│   │   └── lib/
│   │       ├── api-client.ts       # Axios/fetch wrapper pointing to Go API
│   │       ├── auth.ts             # Auth utilities
│   │       └── utils.ts
│   │
│   └── admin/                      # Admin panel (Next.js)
│       ├── package.json
│       ├── next.config.ts
│       ├── app/
│       │   ├── layout.tsx          # Admin layout with sidebar
│       │   ├── page.tsx            # Dashboard with widgets
│       │   └── resources/
│       │       ├── users/page.tsx   # User management page
│       │       └── posts/page.tsx   # Post management page
│       ├── components/
│       │   ├── layout/             # Admin shell, sidebar, navbar
│       │   │   ├── admin-layout.tsx
│       │   │   ├── sidebar.tsx
│       │   │   └── navbar.tsx
│       │   ├── tables/             # DataTable components
│       │   │   ├── data-table.tsx
│       │   │   ├── columns.tsx
│       │   │   └── filters.tsx
│       │   ├── forms/              # Form builder components
│       │   │   ├── form-builder.tsx
│       │   │   ├── fields/
│       │   │   │   ├── text-field.tsx
│       │   │   │   ├── select-field.tsx
│       │   │   │   ├── date-field.tsx
│       │   │   │   └── file-field.tsx
│       │   │   └── form-modal.tsx
│       │   └── widgets/            # Dashboard widgets
│       │       ├── stats-card.tsx
│       │       ├── chart-widget.tsx
│       │       └── recent-activity.tsx
│       └── resources/              # Resource definitions (THE MAGIC)
│           ├── index.ts            # Resource registry
│           ├── users.ts            # User resource definition
│           └── posts.ts            # Post resource definition
│
└── grit/                           # CLI tool (Go)
    ├── go.mod
    ├── cmd/
    │   └── grit/
    │       └── main.go
    ├── internal/
    │   ├── scaffold/               # Project scaffolding
    │   ├── generate/               # Resource code generation
    │   ├── migrate/                # Migration runner
    │   └── templates/              # Go + TS file templates
    └── templates/                  # Project templates
        ├── api/                    # Go API template files
        ├── web/                    # Next.js template files
        ├── admin/                  # Admin panel template files
        └── shared/                 # Shared package template files
```

---

## CLI Commands

| Command | Description |
|---------|-------------|
| `grit new <name>` | Scaffold a new Grit project (Go API + Next.js + Admin) |
| `grit new <name> --api` | Scaffold only the Go API |
| `grit new <name> --full` | Scaffold with Expo mobile app included |
| `grit generate resource <Name>` | Generate full-stack resource (model, handler, hooks, schema, admin page) |
| `grit generate model <Name>` | Generate only a Go model |
| `grit generate handler <Name>` | Generate only a Go handler |
| `grit migrate` | Run database migrations |
| `grit migrate:fresh` | Drop all tables and re-migrate |
| `grit seed` | Run database seeders |
| `grit sync` | Sync Go types → TypeScript types + Zod schemas |
| `grit add queue` | Add Redis queue support |
| `grit add cache` | Add Redis caching |
| `grit add cron` | Add cron job scheduler |
| `grit add mobile` | Add Expo app to the monorepo |
| `grit dev` | Start all services in development mode |
| `grit build` | Build all apps for production |
| `grit deploy` | Deploy to production (Docker-based) |
| `grit studio` | Open GORM Studio in browser |

---

## Resource Definition Example

This is how the Filament-like admin panel works. You define a resource once, and Grit generates everything:

```typescript
// apps/admin/resources/invoices.ts
import { defineResource } from '@grit/admin';

export default defineResource({
  name: 'Invoice',
  endpoint: '/api/invoices',
  icon: 'FileText',

  table: {
    columns: [
      { key: 'id', label: 'ID', sortable: true },
      { key: 'number', label: 'Invoice #', sortable: true, searchable: true },
      { key: 'customer.name', label: 'Customer', relation: 'customer' },
      { key: 'amount', label: 'Amount', format: 'currency' },
      { key: 'status', label: 'Status', badge: {
        paid: { color: 'green', label: 'Paid' },
        pending: { color: 'yellow', label: 'Pending' },
        overdue: { color: 'red', label: 'Overdue' },
      }},
      { key: 'due_date', label: 'Due Date', format: 'date' },
      { key: 'created_at', label: 'Created', format: 'relative' },
    ],
    filters: [
      { key: 'status', type: 'select', options: ['paid', 'pending', 'overdue'] },
      { key: 'created_at', type: 'date-range' },
      { key: 'amount', type: 'number-range' },
    ],
    actions: ['create', 'edit', 'delete', 'export'],
    bulkActions: ['delete', 'export', 'mark-paid'],
  },

  form: {
    fields: [
      { key: 'number', label: 'Invoice Number', type: 'text', required: true },
      { key: 'customer_id', label: 'Customer', type: 'relation', resource: 'customers', displayKey: 'name' },
      { key: 'amount', label: 'Amount', type: 'number', prefix: '$' },
      { key: 'status', label: 'Status', type: 'select', options: ['paid', 'pending', 'overdue'], default: 'pending' },
      { key: 'due_date', label: 'Due Date', type: 'date' },
      { key: 'notes', label: 'Notes', type: 'richtext' },
      { key: 'attachments', label: 'Files', type: 'file', multiple: true },
    ],
  },

  dashboard: {
    widgets: [
      { type: 'stat', label: 'Total Revenue', query: 'sum:amount', format: 'currency' },
      { type: 'stat', label: 'Pending', query: 'count:status=pending', color: 'yellow' },
      { type: 'chart', label: 'Revenue Over Time', type: 'line', query: 'sum:amount:by:month' },
    ],
  },
});
```

The corresponding Go code (`grit generate resource Invoice`) produces:

```go
// apps/api/internal/models/invoice.go
type Invoice struct {
    ID         uint      `gorm:"primarykey" json:"id"`
    Number     string    `gorm:"size:50;uniqueIndex;not null" json:"number"`
    CustomerID uint      `gorm:"not null;index" json:"customer_id"`
    Customer   Customer  `gorm:"foreignKey:CustomerID" json:"customer,omitempty"`
    Amount     float64   `gorm:"not null" json:"amount"`
    Status     string    `gorm:"size:20;default:pending" json:"status"`
    DueDate    time.Time `json:"due_date"`
    Notes      string    `gorm:"type:text" json:"notes"`
    CreatedAt  time.Time `json:"created_at"`
    UpdatedAt  time.Time `json:"updated_at"`
}
```

```go
// apps/api/internal/handlers/invoice.go
// Full CRUD handler with pagination, filtering, sorting, search
// POST /api/invoices, GET /api/invoices, GET /api/invoices/:id,
// PUT /api/invoices/:id, DELETE /api/invoices/:id
```

```typescript
// packages/shared/schemas/invoice.ts
import { z } from 'zod';

export const InvoiceSchema = z.object({
  number: z.string().min(1),
  customer_id: z.number().int().positive(),
  amount: z.number().positive(),
  status: z.enum(['paid', 'pending', 'overdue']).default('pending'),
  due_date: z.string().datetime(),
  notes: z.string().optional(),
});
```

---

## Monetization Strategy

### 1. Open Core (Primary Revenue)

The framework is **free and open source**. Revenue comes from:

| Tier | Price | What's Included |
|------|-------|-----------------|
| **Community (Free)** | $0 | Full framework, CLI, basic admin panel (tables, forms, CRUD), auth, file storage, email, GORM Studio |
| **Pro** | $99/year per developer | Advanced admin components: charts, dashboard builder, kanban boards, calendar views, drag-and-drop form builder, audit logs, activity timeline, advanced filters, PDF export, multi-tenancy support |
| **Team** | $299/year per team | Everything in Pro + priority support, private Discord channel, team license, white-label admin panel (remove Grit branding) |

### 2. Starter Kit Marketplace (Quick Revenue)

Sell pre-built starter kits built on Grit:

| Kit | Price | Description |
|-----|-------|-------------|
| **SaaS Starter** | $149 | Multi-tenant SaaS with billing (Stripe), subscription management, user dashboard, landing page |
| **CRM Starter** | $199 | Contact management, deal pipeline, activity tracking, email integration |
| **E-Commerce Starter** | $199 | Product catalog, cart, checkout, order management, inventory |
| **Agency Starter** | $99 | Project management, client portal, invoicing, time tracking |
| **Blog/CMS Starter** | $79 | Content management, markdown editor, categories, SEO, RSS |

### 3. Plugin Marketplace (Community Revenue)

Once Grit has adoption, open a marketplace for community-built plugins:

- Admin panel themes ($29-$49)
- Custom form field types ($19-$39)
- Dashboard widget packs ($29-$49)
- Auth provider plugins (Google, GitHub, SAML) ($19-$29)
- Payment integrations ($39-$59)
- Grit takes a **20% commission** on all marketplace sales

### 4. Hosted Platform (Long-term)

`grit deploy` pushes to Grit Cloud:

| Plan | Price | Includes |
|------|-------|---------|
| **Hobby** | $0 | 1 project, shared resources, `*.grit.app` domain |
| **Pro** | $20/month | 3 projects, custom domains, auto-scaling |
| **Business** | $50/month | Unlimited projects, priority builds, SLA |

### Revenue Timeline

- **Month 1-3**: Free — focus on GitHub stars, community, content
- **Month 3-6**: Launch Pro tier ($99/year) and first starter kit ($149)
- **Month 6-12**: Launch marketplace, add more starter kits
- **Year 2+**: Launch hosted platform

### Realistic Revenue Targets

- 1,000 GitHub stars → 50 Pro subscribers → **$4,950/year**
- 5,000 GitHub stars → 200 Pro subscribers + 100 starter kit sales → **$34,600/year**
- 10,000+ stars → 500 Pro + marketplace + hosting → **$100,000+/year**

---

## Competitive Advantages

| Advantage | Why It Matters |
|-----------|---------------|
| **Go backend** | 10-50x faster than Node.js/PHP for CPU-bound work. Single binary deployment. Low memory usage. Go is the fastest-growing backend language. |
| **React frontend** | Largest frontend ecosystem. Most developers know React. Massive component library ecosystem. |
| **Admin panel generator** | Nothing like Filament exists for Go + React. This alone is a reason to adopt Grit. |
| **Monorepo with type sync** | Change a Go struct → TypeScript types update automatically. True end-to-end type safety. |
| **GORM Studio built-in** | Visual database browser that no other Go framework offers. |
| **One command setup** | `grit new myapp` gives you everything. Compare that to 2-3 days of boilerplate with manual setup. |
| **Convention over configuration** | Any developer can jump into any Grit project and understand it immediately. |
| **Vibe coding ready** | AI can generate Grit resources confidently because the patterns are predictable and consistent. |

---

## Competitive Landscape

| Tool | Language | Admin Panel | Full-Stack | Monorepo | How Grit Differs |
|------|----------|-------------|------------|----------|-----------------|
| **Laravel + Filament** | PHP | ✅ Excellent | ✅ | ❌ | Grit uses Go (faster) + React (better frontend ecosystem) |
| **Next.js + Prisma** | JS/TS | ❌ | Partial | ❌ | Grit has a real backend, admin panel, and batteries |
| **BHVR** | JS/TS | ❌ | ✅ | ✅ | Grit has admin panel, Go backend, more batteries |
| **t3-stack** | JS/TS | ❌ | Partial | ❌ | Grit has Go backend, admin panel, file storage, queues |
| **Refine** | JS/TS | ✅ | ❌ Frontend only | ❌ | Grit includes the backend |
| **PocketBase** | Go | Basic | ❌ | ❌ | Grit is extensible, has React frontend, admin panel |
| **Buffalo** | Go | ❌ | ❌ | ❌ | Dead project, no React, no admin panel |
| **Django + Admin** | Python | ✅ Basic | ✅ | ❌ | Grit has Go speed, React frontend, modern DX |
| **Rails + ActiveAdmin** | Ruby | ✅ Basic | ✅ | ❌ | Grit has Go speed, React frontend, modern admin panel |

---

## Success Metrics

### Launch (Month 1)
- [ ] 500+ GitHub stars
- [ ] 50+ projects created with `grit new`
- [ ] 10+ community contributors
- [ ] Documentation site live
- [ ] 3+ YouTube tutorials

### Growth (Month 3-6)
- [ ] 2,000+ GitHub stars
- [ ] 20+ Pro subscribers
- [ ] 5+ community plugins
- [ ] Featured in Go Weekly / React Newsletter
- [ ] Conference talk submitted

### Scale (Month 6-12)
- [ ] 5,000+ GitHub stars
- [ ] 100+ Pro subscribers
- [ ] 20+ marketplace plugins
- [ ] First starter kit sales
- [ ] Full-time revenue potential

---

## Brand Identity

| Element | Value |
|---------|-------|
| **Name** | Grit |
| **Tagline** | Go + React. Built with Grit. |
| **Full Name** | Grit Framework |
| **CLI Command** | `grit` |
| **GitHub** | `github.com/MUKE-coder/grit` |
| **Website** | `gritframework.dev` |
| **npm Scope** | `@grit/` |
| **Go Module** | `github.com/MUKE-coder/grit` |
| **Color Palette** | Primary: `#6c5ce7` (purple), Dark: `#0a0a0f`, Success: `#00b894`, Danger: `#ff6b6b` |
| **Fonts** | DM Sans (UI), JetBrains Mono (code) |
| **Theme** | Dark-first, premium, CRM-inspired |

---

## Key Technical Decisions

| Decision | Choice | Rationale |
|----------|--------|-----------|
| Go router | Gin | Most popular Go web framework, best middleware ecosystem |
| ORM | GORM | Most popular Go ORM, auto-migration, hooks, relationships |
| Frontend | Next.js App Router | Industry standard, SSR/SSG, file-based routing |
| Styling | Tailwind + shadcn/ui | Most popular utility-first CSS, great component library |
| Data fetching | React Query | Industry standard for server state, caching, mutations |
| Validation | Zod | Type-safe validation that works with TypeScript |
| Monorepo | Turborepo | Fast, simple, well-supported |
| Package manager | pnpm | Fast, strict, disk-efficient |
| Database | PostgreSQL | Industry standard, GORM supports it natively |
| Cache/Queue | Redis | Industry standard for caching and job queues |
| File storage | S3-compatible | Works with AWS, Cloudflare R2, MinIO |
| Email | Resend | Simple API, great DX, good free tier |
| Deployment | Docker | Universal, consistent, works everywhere |

---

*This document is the source of truth for the Grit project. All development, documentation, and marketing should align with the vision described here.*
