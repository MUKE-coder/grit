import Link from 'next/link'
import { ArrowLeft, ArrowRight } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'
import { CodeBlock } from '@/components/code-block'
import { getDocMetadata } from '@/config/docs-metadata'

export const metadata = getDocMetadata('/docs/ai-skill/llm-guide')

export default function LLMGuidePage() {
  return (
    <div className="min-h-screen bg-background">
      <SiteHeader />
      <DocsSidebar />

      <main className="lg:pl-64">
        <div className="container max-w-screen-xl py-10 px-6">
          <div className="max-w-3xl">
            {/* Header */}
            <div className="mb-10">
              <span className="tag-mono text-primary/80 mb-3 block">For AI Assistants</span>
              <h1 className="text-4xl font-bold tracking-tight mb-4">
                Grit Framework — Complete LLM Reference
              </h1>
              <p className="text-lg text-muted-foreground leading-relaxed">
                Everything an AI assistant needs to build production-quality apps with
                the Grit framework. Read this once and you will understand the stack,
                the conventions, every CLI command, all features, and the rules that
                must never be broken.
              </p>
              <div className="mt-4 p-4 rounded-lg border border-primary/20 bg-primary/5">
                <p className="text-sm text-muted-foreground leading-relaxed">
                  <strong className="text-primary/90">For AI tools:</strong> This page is the canonical
                  reference for building with Grit. When in doubt, consult this page before generating
                  any code. Every convention listed here is enforced by the CLI code generator — deviating
                  from them will break <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">grit generate</code> and{' '}
                  <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">grit sync</code>.
                </p>
              </div>
            </div>

            <div className="prose-grit">

              {/* ── 1. What is Grit ── */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">1. What is Grit?</h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Grit is a <strong className="text-foreground/80">full-stack meta-framework</strong> that fuses a Go backend with Next.js frontends
                  inside a Turborepo monorepo. It is opinionated, batteries-included, and optimised for
                  AI-assisted development. A single CLI command scaffolds a complete production-ready project.
                </p>
                <div className="grid grid-cols-1 sm:grid-cols-2 gap-3 mb-4">
                  {[
                    { label: 'Backend', value: 'Go 1.24 · Gin · GORM · PostgreSQL · Redis' },
                    { label: 'Frontend', value: 'Next.js 15 · React 19 · TypeScript · Tailwind CSS' },
                    { label: 'Admin panel', value: 'Custom Filament-like panel scaffolded with Next.js' },
                    { label: 'Shared types', value: 'packages/shared — Zod schemas + TypeScript types' },
                    { label: 'Monorepo', value: 'Turborepo + pnpm workspaces' },
                    { label: 'Infra', value: 'Docker Compose · MinIO · Mailhog (dev) · S3/R2/Resend (prod)' },
                  ].map((item) => (
                    <div key={item.label} className="p-3 rounded-lg border border-border/30 bg-card/30">
                      <p className="text-xs font-semibold text-muted-foreground uppercase tracking-wider mb-1">{item.label}</p>
                      <p className="text-xs font-mono text-foreground/70">{item.value}</p>
                    </div>
                  ))}
                </div>
              </div>

              {/* ── 2. Project Structure ── */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">2. Project Structure</h2>
                <CodeBlock language="bash" filename="Grit monorepo layout" code={`myapp/
├── apps/
│   ├── api/                          # Go backend
│   │   └── internal/
│   │       ├── config/               # DB, Redis, env, storage
│   │       ├── middleware/           # Gzip, CORS, RequestID, Logger, Auth, Cache
│   │       ├── models/               # GORM structs + model registry
│   │       │   └── models.go         # // GRIT:MODELS marker — NEVER remove
│   │       ├── handlers/             # Thin HTTP handlers (call services)
│   │       ├── services/             # Business logic + GORM queries
│   │       ├── routes/               # Route registration
│   │       │   └── routes.go         # // GRIT:ROUTES marker — NEVER remove
│   │       ├── jobs/                 # asynq background jobs
│   │       └── seeders/              # DB seed data
│   ├── web/                          # Next.js public website
│   │   ├── app/                      # App Router pages
│   │   ├── components/
│   │   └── lib/api-client.ts         # Axios client + uploadFile()
│   └── admin/                        # Next.js admin panel
│       ├── app/(dashboard)/          # Protected admin pages
│       │   └── [resource]/
│       │       └── page.tsx          # Uses defineResource()
│       ├── components/               # Shared UI (DataTable, FormBuilder)
│       ├── hooks/                    # React Query hooks (generated)
│       └── lib/api-client.ts         # Axios client + uploadFile()
├── packages/
│   └── shared/
│       └── src/
│           ├── schemas/              # Zod schemas (generated by grit sync)
│           ├── types/                # TypeScript interfaces
│           └── index.ts              # Re-exports everything
├── docker-compose.yml                # PostgreSQL, Redis, MinIO, Mailhog
├── docker-compose.prod.yml           # Production config
├── turbo.json
├── pnpm-workspace.yaml
├── grit.config.ts                    # Project config
├── GRIT_SKILL.md                     # AI context file (auto-generated)
└── .env                              # Environment variables`} />
              </div>

              {/* ── 3. CLI Commands ── */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">3. All CLI Commands</h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Install the CLI once with{' '}
                  <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">go install github.com/MUKE-coder/grit/cmd/grit@latest</code>.
                  Every command is idempotent — safe to re-run.
                </p>
                <div className="space-y-4">
                  {[
                    {
                      cmd: 'grit new <app-name>',
                      desc: 'Scaffold a complete new Grit project: Go API, Next.js web, Next.js admin, shared package, Docker Compose, .env, GRIT_SKILL.md, and all configuration files.',
                      example: 'grit new myapp',
                    },
                    {
                      cmd: 'grit generate resource <Name> [fields...]',
                      desc: 'Generate a full-stack resource: Go model + GORM migration, handler, service, routes, Zod schema, TypeScript types, React Query hook, and admin page — all wired together.',
                      example: 'grit generate resource Product name:string price:float64 description:text category_id:uint:fk:Category is_active:bool',
                    },
                    {
                      cmd: 'grit generate resource <Name> --no-admin',
                      desc: 'Generate the Go model, handler, service, and routes — but skip the admin panel page and hooks. Useful for API-only resources.',
                      example: 'grit generate resource Webhook url:string secret:string --no-admin',
                    },
                    {
                      cmd: 'grit sync',
                      desc: 'Regenerate the shared/src/schemas and shared/src/types from all Go models. Run this after manually editing a Go model to keep TypeScript in sync.',
                      example: 'grit sync',
                    },
                    {
                      cmd: 'grit dev',
                      desc: 'Start all apps in development mode: Go API with hot reload (air), Next.js web dev server, Next.js admin dev server, all in parallel.',
                      example: 'grit dev',
                    },
                    {
                      cmd: 'grit start',
                      desc: 'Start all apps in production mode (no hot reload). Equivalent to running each app\'s start command in parallel.',
                      example: 'grit start',
                    },
                    {
                      cmd: 'grit migrate',
                      desc: 'Run GORM AutoMigrate for all registered models. Safe to run repeatedly — adds new columns/tables without dropping existing data.',
                      example: 'grit migrate',
                    },
                    {
                      cmd: 'grit seed',
                      desc: 'Run the database seeders defined in apps/api/internal/seeders/. Creates the default admin user and sample data.',
                      example: 'grit seed',
                    },
                    {
                      cmd: 'grit add role <ROLE_NAME>',
                      desc: 'Add a new role to the RBAC system in 7 locations: Go constants, Go middleware, Zod schema, TypeScript type, admin dropdown, seed file, and GRIT_SKILL.md.',
                      example: 'grit add role MODERATOR',
                    },
                    {
                      cmd: 'grit remove resource <Name>',
                      desc: 'Remove a previously generated resource: deletes Go files, removes model from registry, removes routes, and removes admin page.',
                      example: 'grit remove resource Product',
                    },
                    {
                      cmd: 'grit studio',
                      desc: 'Open GORM Studio — a browser-based database browser. Runs at http://localhost:PORT/studio.',
                      example: 'grit studio',
                    },
                    {
                      cmd: 'grit upgrade',
                      desc: 'Upgrade the Grit CLI to the latest version.',
                      example: 'grit upgrade',
                    },
                    {
                      cmd: 'grit version',
                      desc: 'Print the installed CLI version.',
                      example: 'grit version',
                    },
                  ].map((item) => (
                    <div key={item.cmd} className="rounded-xl border border-border/40 bg-card/50 overflow-hidden">
                      <div className="px-4 py-2.5 border-b border-border/20 bg-accent/20 flex items-center gap-2">
                        <code className="text-xs font-mono font-semibold text-primary/80">{item.cmd}</code>
                      </div>
                      <div className="p-4">
                        <p className="text-xs text-muted-foreground/70 leading-relaxed mb-3">{item.desc}</p>
                        <div className="rounded-lg border border-border/20 bg-background/60 px-3 py-2">
                          <span className="text-primary/40 font-mono text-xs select-none">$ </span>
                          <span className="font-mono text-xs text-foreground/60">{item.example}</span>
                        </div>
                      </div>
                    </div>
                  ))}
                </div>
              </div>

              {/* ── 4. Field Types ── */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">4. Field Types for Code Generation</h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  When running <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">grit generate resource</code>,
                  fields follow the format <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">name:type[:fk:RelatedModel]</code>.
                </p>
                <div className="rounded-xl border border-border/40 bg-card/50 overflow-hidden">
                  <table className="w-full text-sm">
                    <thead>
                      <tr className="border-b border-border/30 bg-accent/20">
                        <th className="px-4 py-2.5 text-left text-xs font-semibold text-muted-foreground uppercase">Type</th>
                        <th className="px-4 py-2.5 text-left text-xs font-semibold text-muted-foreground uppercase">Go Type</th>
                        <th className="px-4 py-2.5 text-left text-xs font-semibold text-muted-foreground uppercase">Zod / TS</th>
                        <th className="px-4 py-2.5 text-left text-xs font-semibold text-muted-foreground uppercase">Admin field</th>
                      </tr>
                    </thead>
                    <tbody className="divide-y divide-border/20 text-xs font-mono">
                      {[
                        { type: 'string', go: 'string', zod: 'z.string()', admin: 'text input' },
                        { type: 'text', go: 'string', zod: 'z.string()', admin: 'textarea' },
                        { type: 'richtext', go: 'string', zod: 'z.string()', admin: 'rich text editor' },
                        { type: 'int', go: 'int', zod: 'z.number().int()', admin: 'number input' },
                        { type: 'uint', go: 'uint', zod: 'z.number().int().nonnegative()', admin: 'number input' },
                        { type: 'float64', go: 'float64', zod: 'z.number()', admin: 'number input' },
                        { type: 'bool', go: 'bool', zod: 'z.boolean()', admin: 'toggle / checkbox' },
                        { type: 'time', go: 'time.Time', zod: 'z.string().datetime()', admin: 'date picker' },
                        { type: 'image', go: 'string', zod: 'z.string().url()', admin: 'image upload (dropzone)' },
                        { type: 'file', go: 'string', zod: 'z.string().url()', admin: 'file upload (dropzone)' },
                        { type: 'enum:A,B,C', go: 'string', zod: 'z.enum(["A","B","C"])', admin: 'select dropdown' },
                        { type: 'uint:fk:Model', go: 'uint + Model field', zod: 'z.number().int()', admin: 'belongs_to select' },
                        { type: '[]uint:m2m:Model', go: '[]Model (GORM M2M)', zod: 'z.array(z.number())', admin: 'many_to_many multi-select' },
                      ].map((row) => (
                        <tr key={row.type}>
                          <td className="px-4 py-2 text-primary/80">{row.type}</td>
                          <td className="px-4 py-2 text-foreground/60">{row.go}</td>
                          <td className="px-4 py-2 text-foreground/60">{row.zod}</td>
                          <td className="px-4 py-2 text-muted-foreground/60">{row.admin}</td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>
              </div>

              {/* ── 5. Code Patterns ── */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">5. Code Patterns</h2>

                <h3 className="text-lg font-semibold tracking-tight mb-3">Backend: Handler → Service → Model</h3>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Handlers are thin controllers. They parse the request, call a service method,
                  and return JSON. All database logic lives in services.
                </p>
                <CodeBlock language="go" filename="apps/api/internal/handlers/product_handler.go" code={`type ProductHandler struct {
    Service *services.ProductService
    Storage *config.Storage  // only if resource has image/file fields
}

func (h *ProductHandler) List(c *gin.Context) {
    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
    products, total, pages, err := h.Service.List(page, pageSize)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": gin.H{"code": "INTERNAL_ERROR", "message": "Failed to fetch products"},
        })
        return
    }
    c.JSON(http.StatusOK, gin.H{
        "data": products,
        "meta": gin.H{"total": total, "page": page, "page_size": pageSize, "pages": pages},
    })
}

func (h *ProductHandler) Create(c *gin.Context) {
    var req struct {
        Name        string  \`json:"name" binding:"required"\`
        Price       float64 \`json:"price" binding:"required"\`
        Description string  \`json:"description"\`
    }
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusUnprocessableEntity, gin.H{
            "error": gin.H{"code": "VALIDATION_ERROR", "message": err.Error()},
        })
        return
    }
    product, err := h.Service.Create(req.Name, req.Price, req.Description)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": gin.H{"code": "INTERNAL_ERROR", "message": "Failed to create product"},
        })
        return
    }
    c.JSON(http.StatusCreated, gin.H{"data": product, "message": "Product created successfully"})
}`} />

                <h3 className="text-lg font-semibold tracking-tight mb-3 mt-8">Frontend: Component → Hook → API</h3>
                <CodeBlock language="tsx" filename="apps/admin/app/(dashboard)/products/page.tsx" code={`'use client'
import { useProducts, useCreateProduct } from '@/hooks/use-products'
import { DataTable } from '@/components/data-table'
import { productResource } from './_resource'

export default function ProductsPage() {
  const { data, isLoading } = useProducts()
  const createProduct = useCreateProduct()

  return (
    <DataTable
      resource={productResource}
      data={data?.data ?? []}
      total={data?.meta?.total ?? 0}
      isLoading={isLoading}
      onCreate={(values) => createProduct.mutateAsync(values)}
    />
  )
}`} />
              </div>

              {/* ── 6. API Response Format ── */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">6. API Response Format</h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Every API response follows this exact format. Never deviate from it —
                  the frontend hooks and admin components depend on this structure.
                </p>
                <div className="space-y-4">
                  <CodeBlock language="json" filename="Success — single record" code={`{
  "data": { "id": 1, "name": "Widget", "price": 29.99 },
  "message": "Product created successfully"
}`} />
                  <CodeBlock language="json" filename="Success — paginated list" code={`{
  "data": [ { "id": 1, "name": "Widget" }, { "id": 2, "name": "Gadget" } ],
  "meta": {
    "total": 42,
    "page": 1,
    "page_size": 20,
    "pages": 3
  }
}`} />
                  <CodeBlock language="json" filename="Error response" code={`{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "name is required",
    "details": { "field": "name", "tag": "required" }
  }
}`} />
                  <div className="rounded-xl border border-border/40 bg-card/50 overflow-hidden">
                    <table className="w-full text-xs">
                      <thead>
                        <tr className="border-b border-border/30 bg-accent/20">
                          <th className="px-4 py-2 text-left font-semibold text-muted-foreground uppercase">Situation</th>
                          <th className="px-4 py-2 text-left font-semibold text-muted-foreground uppercase">HTTP Status</th>
                          <th className="px-4 py-2 text-left font-semibold text-muted-foreground uppercase">Error Code</th>
                        </tr>
                      </thead>
                      <tbody className="divide-y divide-border/20 font-mono">
                        {[
                          { situation: 'Missing required field', status: '422', code: 'VALIDATION_ERROR' },
                          { situation: 'Record not found', status: '404', code: 'NOT_FOUND' },
                          { situation: 'Not authenticated', status: '401', code: 'UNAUTHORIZED' },
                          { situation: 'Insufficient role', status: '403', code: 'FORBIDDEN' },
                          { situation: 'Internal/DB error', status: '500', code: 'INTERNAL_ERROR' },
                          { situation: 'Conflict (duplicate)', status: '409', code: 'CONFLICT' },
                        ].map((row) => (
                          <tr key={row.situation}>
                            <td className="px-4 py-2 text-foreground/60">{row.situation}</td>
                            <td className="px-4 py-2 text-amber-500/80">{row.status}</td>
                            <td className="px-4 py-2 text-primary/70">{row.code}</td>
                          </tr>
                        ))}
                      </tbody>
                    </table>
                  </div>
                </div>
              </div>

              {/* ── 7. Code Markers ── */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">7. Code Markers — NEVER Delete</h2>
                <div className="p-4 rounded-lg border border-red-500/20 bg-red-500/5 mb-4">
                  <p className="text-sm text-red-400/80 leading-relaxed font-medium">
                    These comments are used by the Grit CLI to inject code. Removing them breaks
                    <code className="mx-1 text-xs font-mono bg-red-500/10 px-1.5 py-0.5 rounded">grit generate</code>
                    and
                    <code className="mx-1 text-xs font-mono bg-red-500/10 px-1.5 py-0.5 rounded">grit add role</code>.
                    Never remove, rename, or move them. Always place new code <em>between</em> the markers.
                  </p>
                </div>
                <CodeBlock language="go" filename="apps/api/internal/models/models.go" code={`var RegisteredModels = []interface{}{
    // GRIT:MODELS — do not remove this comment
    &User{},
    &Upload{},
    &Blog{},
    &Product{}, // grit generate adds new models here
    // END GRIT:MODELS
}`} />
                <CodeBlock language="go" filename="apps/api/internal/routes/routes.go" code={`// GRIT:ROUTES — do not remove this comment
uploadHandler := &handlers.UploadHandler{...}
productHandler := &handlers.ProductHandler{...}
// END GRIT:ROUTES

// GRIT:PROTECTED_ROUTES — do not remove this comment
protected.GET("/products", productHandler.List)
protected.POST("/products", productHandler.Create)
// END GRIT:PROTECTED_ROUTES`} />
              </div>

              {/* ── 8. Naming Conventions ── */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">8. Naming Conventions</h2>
                <div className="rounded-xl border border-border/40 bg-card/50 overflow-hidden">
                  <table className="w-full text-sm">
                    <thead>
                      <tr className="border-b border-border/30 bg-accent/20">
                        <th className="px-4 py-3 text-left text-xs font-semibold text-muted-foreground uppercase">What</th>
                        <th className="px-4 py-3 text-left text-xs font-semibold text-muted-foreground uppercase">Convention</th>
                        <th className="px-4 py-3 text-left text-xs font-semibold text-muted-foreground uppercase">Example</th>
                      </tr>
                    </thead>
                    <tbody className="divide-y divide-border/20 font-mono text-xs">
                      {[
                        { what: 'Go files', conv: 'snake_case', ex: 'product_handler.go' },
                        { what: 'Go structs / types', conv: 'PascalCase', ex: 'ProductHandler, CreateProductReq' },
                        { what: 'Go functions', conv: 'camelCase (methods), PascalCase (exported)', ex: 'List(), GetByID()' },
                        { what: 'DB table names', conv: 'plural snake_case (GORM default)', ex: 'products, order_items' },
                        { what: 'DB column names', conv: 'snake_case (GORM default)', ex: 'created_at, category_id' },
                        { what: 'API routes', conv: 'plural lowercase, kebab for multi-word', ex: '/api/products, /api/blog-posts' },
                        { what: 'TypeScript files', conv: 'kebab-case', ex: 'use-products.ts, product-form.tsx' },
                        { what: 'React components', conv: 'PascalCase', ex: 'ProductForm, DataTable' },
                        { what: 'React hooks', conv: 'camelCase with "use" prefix', ex: 'useProducts, useCreateProduct' },
                        { what: 'Zod schemas', conv: 'PascalCase + Schema suffix', ex: 'CreateProductSchema' },
                        { what: 'TypeScript types', conv: 'PascalCase', ex: 'Product, CreateProductInput' },
                        { what: 'Environment variables', conv: 'SCREAMING_SNAKE_CASE', ex: 'DATABASE_URL, JWT_SECRET' },
                      ].map((row) => (
                        <tr key={row.what}>
                          <td className="px-4 py-2 text-foreground/70">{row.what}</td>
                          <td className="px-4 py-2 text-primary/70">{row.conv}</td>
                          <td className="px-4 py-2 text-muted-foreground/60">{row.ex}</td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>
              </div>

              {/* ── 9. All Batteries ── */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">9. Built-in Batteries</h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Grit ships these services pre-wired. They are enabled in the Go API and injected
                  into handlers via the{' '}
                  <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">Services</code> struct.
                </p>
                <div className="space-y-3">
                  {[
                    {
                      name: 'File Storage (S3-compatible)',
                      pkg: 'aws-sdk-go-v2/service/s3',
                      desc: 'Upload, delete, get signed URLs, and generate presigned PUT URLs. Dev uses MinIO; prod uses Cloudflare R2 or AWS S3. Presigned uploads bypass the API for large files.',
                      routes: 'POST /api/uploads/presign, POST /api/uploads/complete, GET /api/uploads, DELETE /api/uploads/:id',
                    },
                    {
                      name: 'Email (Resend)',
                      pkg: 'resend-go',
                      desc: 'Send transactional emails with Go HTML templates. Built-in templates: welcome, password reset, email verification, generic notification. Dev uses Mailhog.',
                      routes: 'Internal service — called from auth handlers and background jobs',
                    },
                    {
                      name: 'Background Jobs (asynq)',
                      pkg: 'hibiken/asynq',
                      desc: 'Redis-backed async job queue. Workers run in a separate goroutine pool. Built-in jobs: image processing, email sending. Add custom jobs in apps/api/internal/jobs/.',
                      routes: 'Admin UI at /jobs (if enabled)',
                    },
                    {
                      name: 'Cron Scheduler (asynq)',
                      pkg: 'hibiken/asynq',
                      desc: 'Define recurring tasks with cron expressions. Tasks are enqueued to the same asynq queue. Monitored through the background jobs dashboard.',
                      routes: 'Configured in apps/api/internal/config/cron.go',
                    },
                    {
                      name: 'Redis Caching',
                      pkg: 'redis/go-redis/v9',
                      desc: 'Cache arbitrary values with TTL. Cache middleware for caching full API responses by URL. Cache invalidation via prefix or key.',
                      routes: 'Middleware: middleware.Cache(duration). Service: cache.Set/Get/Delete',
                    },
                    {
                      name: 'AI Integration (Claude / OpenAI)',
                      pkg: 'anthropic-sdk-go, openai-go',
                      desc: 'Pluggable AI provider. Supports streaming responses. Configure ANTHROPIC_API_KEY or OPENAI_API_KEY in .env.',
                      routes: 'POST /api/ai/chat (streaming)',
                    },
                    {
                      name: 'Security (Sentinel)',
                      pkg: 'MUKE-coder/sentinel',
                      desc: 'WAF, rate limiting (per-IP), brute-force protection, anomaly detection, IP geolocation, threat dashboard. ExcludeRoutes: /pulse/*, /sentinel/*, /docs/*, /studio/*.',
                      routes: 'Dashboard at /sentinel',
                    },
                    {
                      name: 'Observability (Pulse)',
                      pkg: 'MUKE-coder/pulse',
                      desc: 'Request tracing, DB query monitoring, runtime metrics (goroutines, memory, GC), error tracking, health checks, Prometheus export, alerting.',
                      routes: 'Dashboard at /pulse',
                    },
                    {
                      name: 'API Docs (gin-docs)',
                      pkg: 'gin-docs',
                      desc: 'Auto-generated OpenAPI spec from Gin routes. Interactive Scalar/Swagger UI. No annotations needed — routes are introspected at startup.',
                      routes: 'Docs at /docs',
                    },
                    {
                      name: 'DB Browser (GORM Studio)',
                      pkg: 'MUKE-coder/gorm-studio',
                      desc: 'Browser-based database browser for PostgreSQL. View tables, run queries, inspect records. Dev-only — disable in production.',
                      routes: 'Browser at /studio',
                    },
                  ].map((item) => (
                    <div key={item.name} className="rounded-xl border border-border/40 bg-card/50 overflow-hidden">
                      <div className="px-4 py-2.5 border-b border-border/20 bg-accent/20 flex items-center justify-between">
                        <span className="text-sm font-semibold text-foreground/80">{item.name}</span>
                        <code className="text-[10px] font-mono text-muted-foreground/50">{item.pkg}</code>
                      </div>
                      <div className="p-4 space-y-2">
                        <p className="text-xs text-muted-foreground/70 leading-relaxed">{item.desc}</p>
                        <p className="text-[11px] font-mono text-muted-foreground/40">{item.routes}</p>
                      </div>
                    </div>
                  ))}
                </div>
              </div>

              {/* ── 10. Performance ── */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">10. Built-in Performance Optimisations</h2>
                <div className="grid grid-cols-1 sm:grid-cols-2 gap-3">
                  {[
                    { label: 'Gzip middleware', desc: 'All responses compressed at BestSpeed. 60-80% smaller JSON.' },
                    { label: 'Request ID tracing', desc: 'X-Request-ID on every request. Included in all log lines.' },
                    { label: 'DB connection pool', desc: 'MaxOpen:100, MaxIdle:10, Lifetime:30m, IdleTime:10m.' },
                    { label: 'Cache-Control headers', desc: 'Public blog endpoints: 5min list, 1hr single post.' },
                    { label: 'Presigned URL uploads', desc: 'Files go direct to S3/R2 — API never touched by large uploads.' },
                    { label: 'Async background jobs', desc: 'Emails, thumbnails, webhooks run in asynq workers.' },
                    { label: 'Redis response cache', desc: 'Entire API responses cached by URL with configurable TTL.' },
                    { label: 'Server Components (Next.js)', desc: 'Data fetching on server. Zero JS bundle for data layers.' },
                    { label: 'ISR / revalidate', desc: 'Public pages cached at CDN edge. Revalidate in background.' },
                    { label: 'React Query caching', desc: 'staleTime: 30s in admin. Instant back-navigation.' },
                    { label: 'next/image', desc: 'WebP/AVIF, lazy load, correct sizing, CDN-cached.' },
                    { label: 'Turborepo cache', desc: 'CI builds replay cached output. 4 min → <30 sec.' },
                  ].map((item) => (
                    <div key={item.label} className="p-3 rounded-lg border border-border/30 bg-card/30">
                      <p className="text-xs font-semibold text-primary/80 mb-1">{item.label}</p>
                      <p className="text-xs text-muted-foreground/60 leading-relaxed">{item.desc}</p>
                    </div>
                  ))}
                </div>
              </div>

              {/* ── 11. Golden Rules ── */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">11. Golden Rules — Never Break These</h2>
                <div className="p-4 rounded-lg border border-amber-500/20 bg-amber-500/5 mb-6">
                  <p className="text-sm text-amber-400/80 leading-relaxed">
                    These rules are non-negotiable. Violating them will cause silent failures,
                    broken code generation, or corrupted project state.
                  </p>
                </div>
                <div className="space-y-3">
                  {[
                    {
                      rule: 'Never remove GRIT: markers',
                      detail: '// GRIT:MODELS, // END GRIT:MODELS, // GRIT:ROUTES, // END GRIT:ROUTES — these are injection points for the CLI. Deleting them permanently breaks grit generate and grit add role.',
                    },
                    {
                      rule: 'Never use multipart/form-data for file uploads',
                      detail: 'Grit uses presigned URL uploads. Files go directly to S3/R2/MinIO. The uploadFile() function in lib/api-client.ts handles the 3-step flow. Never POST a FormData object to the API for file uploads.',
                    },
                    {
                      rule: 'Always use the standard error response shape',
                      detail: 'Every error must return { "error": { "code": "...", "message": "..." } }. Frontend hooks and admin components check for this exact shape. A different error format will break the UI.',
                    },
                    {
                      rule: 'Always register new models in models.go between the markers',
                      detail: 'GORM AutoMigrate only runs for models in the RegisteredModels slice. A model not registered will never create its table.',
                    },
                    {
                      rule: 'Always register new routes in routes.go between the markers',
                      detail: 'Routes added outside the markers are still valid Go, but grit remove resource will not be able to clean them up.',
                    },
                    {
                      rule: 'Keep handlers thin — no DB logic in handlers',
                      detail: 'Handlers parse requests and call services. All GORM queries go in service files. This is required for the handler pattern to be consistent with generated code.',
                    },
                    {
                      rule: 'Import from @shared/schemas, not from individual apps',
                      detail: 'TypeScript types and Zod schemas live in packages/shared. Both web and admin apps import from there. Never duplicate schemas between apps.',
                    },
                    {
                      rule: 'Run grit sync after manually editing Go models',
                      detail: 'The shared package is generated from Go structs. If you add a field to a Go model manually, run grit sync to regenerate the Zod schema and TypeScript type.',
                    },
                    {
                      rule: 'Never force-push to main',
                      detail: 'The main branch is the source of truth for the CLI install (go install ...@latest). Force-pushing can corrupt the module cache for all users.',
                    },
                    {
                      rule: 'Disable GORM Studio and Pulse in production',
                      detail: 'These dashboards expose internal DB and metrics data. Set GORM_STUDIO_ENABLED=false and PULSE_ENABLED=false in production environment variables.',
                    },
                  ].map((item, i) => (
                    <div key={item.rule} className="flex gap-3 p-4 rounded-lg border border-border/30 bg-card/30">
                      <div className="flex h-6 w-6 shrink-0 items-center justify-center rounded-md bg-amber-500/10 border border-amber-500/20 text-xs font-mono font-semibold text-amber-500/80 mt-0.5">
                        {i + 1}
                      </div>
                      <div>
                        <p className="text-sm font-semibold text-foreground/80 mb-1">{item.rule}</p>
                        <p className="text-xs text-muted-foreground/60 leading-relaxed">{item.detail}</p>
                      </div>
                    </div>
                  ))}
                </div>
              </div>

              {/* ── 12. Quick Reference ── */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">12. Quick Build Reference</h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Copy-paste recipes for the most common tasks.
                </p>
                <div className="space-y-4">
                  <CodeBlock language="bash" filename="Start a new project" code={`# Install CLI (once)
go install github.com/MUKE-coder/grit/cmd/grit@latest

# Create project
grit new myapp
cd myapp
cp .env.example .env  # fill in your values

# Start Docker services (PostgreSQL, Redis, MinIO, Mailhog)
docker compose up -d

# Migrate + seed
grit migrate
grit seed

# Start dev servers (Go API + web + admin)
grit dev`} />
                  <CodeBlock language="bash" filename="Add a resource (full-stack)" code={`# Generates model, handler, service, routes, Zod schema, TS types, hook, admin page
grit generate resource Invoice \\
  number:string \\
  amount:float64 \\
  due_date:time \\
  status:enum:draft,sent,paid,overdue \\
  customer_id:uint:fk:Customer \\
  notes:text

# Then re-run migration
grit migrate`} />
                  <CodeBlock language="bash" filename="Add a new role" code={`# Updates Go constants, middleware, Zod enum, TypeScript type,
# admin dropdown, seed file, and GRIT_SKILL.md — all in one command
grit add role MODERATOR`} />
                  <CodeBlock language="bash" filename="Sync Go → TypeScript" code={`# After manually editing a Go model, regenerate the shared package
grit sync`} />
                </div>
              </div>

              {/* Nav */}
              <div className="flex items-center justify-between pt-6 border-t border-border/30">
                <Button variant="ghost" size="sm" asChild className="text-muted-foreground/60 hover:text-foreground">
                  <Link href="/docs/ai-skill" className="gap-1.5">
                    <ArrowLeft className="h-3.5 w-3.5" />
                    LLM Skill Guide
                  </Link>
                </Button>
                <Button variant="ghost" size="sm" asChild className="text-muted-foreground/60 hover:text-foreground">
                  <Link href="/docs/changelog" className="gap-1.5">
                    Changelog
                    <ArrowRight className="h-3.5 w-3.5" />
                  </Link>
                </Button>
              </div>

            </div>
          </div>
        </div>
      </main>
    </div>
  )
}
