import Link from 'next/link'
import { ArrowLeft, ArrowRight } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'
import { CodeBlock } from '@/components/code-block'
import { Diagram, DiagramBox, DiagramRow, DiagramArrow, DiagramLegend, FileTree } from '@/components/diagram'
import { getDocMetadata } from '@/config/docs-metadata'

export const metadata = getDocMetadata('/docs/concepts/architecture')

export default function ArchitecturePage() {
  return (
    <div className="min-h-screen bg-background isolate">
      <SiteHeader />
      <DocsSidebar />

      <main className="lg:pl-64">
        <div className="container max-w-screen-xl py-10 px-6">
          <div className="max-w-3xl">
            {/* Header */}
            <div className="mb-10">
              <span className="tag-mono text-primary/80 mb-3 block">Core Concepts</span>
              <h1 className="text-4xl font-bold tracking-tight mb-4">
                Architecture Overview
              </h1>
              <p className="text-lg text-muted-foreground leading-relaxed">
                Grit is a full-stack monorepo that fuses a Go backend with Next.js frontends,
                a shared type package, and an admin panel. This page explains how all the layers
                fit together.
              </p>
            </div>

            <div className="prose-grit">
              {/* Monorepo Layout */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Monorepo Structure
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Every Grit project is a Turborepo-powered monorepo managed by pnpm. The Go API lives
                  alongside the Next.js frontends and a shared TypeScript package. This structure allows
                  all apps to share validation schemas, types, and constants from a single source of truth.
                </p>
                <FileTree
                  title="myapp/"
                  nodes={[
                    { name: 'apps/', type: 'folder', depth: 0 },
                    { name: 'api/', type: 'folder', depth: 1, comment: 'Go backend (Gin + GORM)' },
                    { name: 'web/', type: 'folder', depth: 1, comment: 'Next.js main frontend' },
                    { name: 'admin/', type: 'folder', depth: 1, comment: 'Next.js admin panel' },
                    { name: 'packages/', type: 'folder', depth: 0 },
                    { name: 'shared/', type: 'folder', depth: 1, comment: 'Zod schemas, TS types, constants' },
                    { name: 'docker-compose.yml', type: 'file', depth: 0, comment: 'PostgreSQL, Redis, MinIO, Mailhog' },
                    { name: 'turbo.json', type: 'file', depth: 0, comment: 'Monorepo task orchestration' },
                    { name: 'pnpm-workspace.yaml', type: 'file', depth: 0, comment: 'Workspace definition' },
                    { name: '.env', type: 'file', depth: 0, comment: 'Environment variables' },
                  ]}
                />
                <p className="text-sm text-muted-foreground/60 mt-3">
                  Turborepo handles parallel builds and caching across apps. The <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">pnpm-workspace.yaml</code> file
                  links the frontend apps to the shared package, so importing <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">@shared/schemas</code> works
                  without publishing to npm.
                </p>
              </div>

              {/* Architecture Diagram */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  High-Level Architecture
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  The following diagram shows how the layers communicate. The browser talks to the
                  Next.js frontends, which talk to the Go API over REST. The Go API manages the
                  database, cache, file storage, job queue, and email.
                </p>
                <Diagram>
                  {/* Clients */}
                  <div className="mb-1 text-center text-[11px] font-mono uppercase tracking-wider text-muted-foreground/50">Browser</div>
                  <DiagramRow>
                    <DiagramBox tone="blue" title="Web App" sub=":3000" />
                    <DiagramBox tone="blue" title="Admin Panel" sub=":3001" />
                    <DiagramBox tone="cyan" title="GORM Studio" sub=":8080/studio" />
                  </DiagramRow>

                  <DiagramArrow label="REST + JWT" />

                  {/* API */}
                  <DiagramBox tone="primary" title="Go API \u2014 :8080" sub="Gin Router \u2192 Middleware \u2192 Handlers \u2192 Services" />
                  <div className="mt-2 grid gap-2 sm:grid-cols-3 text-[11px] text-muted-foreground/80">
                    <div className="rounded-md border border-border/40 bg-card/30 px-3 py-2"><span className="font-semibold text-foreground/80">Middleware</span> \u2014 CORS, Auth (JWT), Logger, Recovery</div>
                    <div className="rounded-md border border-border/40 bg-card/30 px-3 py-2"><span className="font-semibold text-foreground/80">Handlers</span> \u2014 thin HTTP layer, request/response only</div>
                    <div className="rounded-md border border-border/40 bg-card/30 px-3 py-2"><span className="font-semibold text-foreground/80">Services</span> \u2014 business logic, DB queries, validation</div>
                  </div>

                  <DiagramArrow />

                  {/* Data & services */}
                  <DiagramRow>
                    <DiagramBox tone="green" title="PostgreSQL" sub=":5434" />
                    <DiagramBox tone="rose" title="Redis" sub=":6380" />
                    <DiagramBox tone="amber" title="MinIO" sub=":9002" />
                    <DiagramBox tone="violet" title="Resend" sub="email" />
                    <DiagramBox tone="default" title="Jobs" sub="asynq" />
                  </DiagramRow>

                  <DiagramLegend
                    items={[
                      { tone: 'blue', label: 'Frontend (React)' },
                      { tone: 'primary', label: 'Go API' },
                      { tone: 'green', label: 'Data & services' },
                    ]}
                  />
                </Diagram>
              </div>

              {/* Go API Layer */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Go API Layer
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  The backend is a Go application built with Gin (HTTP router) and GORM (ORM).
                  It follows a layered architecture: routes register handlers, handlers call services,
                  and services interact with the database through GORM models.
                </p>

                <FileTree
                  title="apps/api/"
                  nodes={[
                    { name: 'cmd/server/main.go', type: 'file', depth: 0, comment: 'Entry point: config, DB, routes, Gin' },
                    { name: 'internal/', type: 'folder', depth: 0 },
                    { name: 'config/config.go', type: 'file', depth: 1, comment: 'Reads .env \u2192 typed config' },
                    { name: 'database/database.go', type: 'file', depth: 1, comment: 'GORM connection' },
                    { name: 'models/', type: 'folder', depth: 1, comment: 'GORM structs (User, Post, \u2026)' },
                    { name: 'handlers/', type: 'folder', depth: 1, comment: 'Gin handlers (HTTP request/response)' },
                    { name: 'services/', type: 'folder', depth: 1, comment: 'Business logic (queries, validation)' },
                    { name: 'middleware/', type: 'folder', depth: 1, comment: 'Auth, CORS, logging, rate limiting' },
                    { name: 'routes/routes.go', type: 'file', depth: 1, comment: 'Route registration + GORM Studio mount' },
                    { name: 'mail/', type: 'folder', depth: 1, comment: 'Resend email client + templates' },
                    { name: 'storage/', type: 'folder', depth: 1, comment: 'S3-compatible file storage' },
                    { name: 'jobs/', type: 'folder', depth: 1, comment: 'Background jobs (asynq + Redis)' },
                    { name: 'cron/', type: 'folder', depth: 1, comment: 'Cron scheduler' },
                    { name: 'cache/', type: 'folder', depth: 1, comment: 'Redis caching layer' },
                    { name: 'ai/', type: 'folder', depth: 1, comment: 'AI (Vercel AI Gateway)' },
                  ]}
                />

                <h3 className="text-xl font-semibold tracking-tight mt-8 mb-3">
                  Handler / Service Separation
                </h3>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Handlers are thin -- they parse HTTP requests, call a service method, and return
                  a JSON response. All business logic, database queries, and validation live in
                  the service layer. This separation makes code testable and reusable.
                </p>
                <CodeBlock language="go" filename="internal/handlers/post.go" code={`// Handler: thin HTTP layer
func (h *PostHandler) List(c *gin.Context) {
    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    search := c.Query("search")
    sortBy := c.DefaultQuery("sort_by", "created_at")

    items, total, pages, err := h.Service.List(params)
    if err != nil {
        c.JSON(500, gin.H{"error": gin.H{"code": "INTERNAL_ERROR", ...}})
        return
    }

    c.JSON(200, gin.H{
        "data": items,
        "meta": gin.H{"total": total, "page": page, "pages": pages},
    })
}`} />

                <CodeBlock language="go" filename="internal/services/post.go" className="mt-4 mb-6" code={`// Service: business logic + database queries
func (s *PostService) List(params PostListParams) ([]models.Post, int64, int, error) {
    query := s.DB.Model(&models.Post{})

    if params.Search != "" {
        query = query.Where("title ILIKE ? OR content ILIKE ?",
            "%"+params.Search+"%", "%"+params.Search+"%")
    }

    var total int64
    query.Count(&total)

    var items []models.Post
    offset := (params.Page - 1) * params.PageSize
    query.Order(params.SortBy + " " + params.SortOrder).
        Offset(offset).Limit(params.PageSize).Find(&items)

    pages := int(math.Ceil(float64(total) / float64(params.PageSize)))
    return items, total, pages, nil
}`} />

                <h3 className="text-xl font-semibold tracking-tight mt-8 mb-3">
                  Middleware Stack
                </h3>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Every request passes through a middleware chain before reaching a handler.
                  Grit ships with four built-in middleware layers:
                </p>
                <div className="rounded-lg border border-border/30 bg-card/30 overflow-hidden">
                  <table className="w-full text-sm">
                    <thead>
                      <tr className="border-b border-border/30 bg-accent/20">
                        <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Middleware</th>
                        <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Purpose</th>
                      </tr>
                    </thead>
                    <tbody className="text-muted-foreground">
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-mono text-xs">Recovery</td>
                        <td className="px-4 py-2.5">Catches panics and returns a 500 JSON error instead of crashing</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-mono text-xs">CORS</td>
                        <td className="px-4 py-2.5">Allows cross-origin requests from the frontend origins (localhost:3000, :3001)</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-mono text-xs">Logger</td>
                        <td className="px-4 py-2.5">Structured JSON logging with method, path, status, and duration</td>
                      </tr>
                      <tr>
                        <td className="px-4 py-2.5 font-mono text-xs">Auth (JWT)</td>
                        <td className="px-4 py-2.5">Validates JWT tokens, attaches user to context, enforces roles</td>
                      </tr>
                    </tbody>
                  </table>
                </div>
              </div>

              {/* Route Groups */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Route Groups
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  The Go API organizes routes into three groups with different levels of
                  authentication and authorization:
                </p>
                <div className="rounded-lg border border-border/30 bg-card/30 overflow-hidden">
                  <table className="w-full text-sm">
                    <thead>
                      <tr className="border-b border-border/30 bg-accent/20">
                        <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Group</th>
                        <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Auth</th>
                        <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Example Routes</th>
                      </tr>
                    </thead>
                    <tbody className="text-muted-foreground">
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-mono text-xs">public</td>
                        <td className="px-4 py-2.5">None</td>
                        <td className="px-4 py-2.5 font-mono text-xs">/api/health, /api/auth/login, /api/auth/register</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-mono text-xs">protected</td>
                        <td className="px-4 py-2.5">Valid JWT</td>
                        <td className="px-4 py-2.5 font-mono text-xs">/api/posts, /api/users/:id, /api/auth/me</td>
                      </tr>
                      <tr>
                        <td className="px-4 py-2.5 font-mono text-xs">admin</td>
                        <td className="px-4 py-2.5">JWT + admin role</td>
                        <td className="px-4 py-2.5 font-mono text-xs">DELETE /api/users/:id, DELETE /api/posts/:id</td>
                      </tr>
                    </tbody>
                  </table>
                </div>
                <p className="text-sm text-muted-foreground/60 mt-3">
                  When the code generator creates a new resource, it automatically injects CRUD routes
                  into the <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">protected</code> group
                  (list, get, create, update) and the delete route into
                  the <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">admin</code> group.
                </p>
              </div>

              {/* Next.js Frontend Layer */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Next.js Frontend Layer
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Both the web app and admin panel are Next.js 14+ applications using the App Router.
                  They share the same data-fetching pattern: React Query hooks call the Go API through
                  an Axios-based API client that automatically injects JWT tokens.
                </p>

                <h3 className="text-xl font-semibold tracking-tight mt-6 mb-3">
                  Web App (<code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">apps/web</code>)
                </h3>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  The main customer-facing frontend. Ships with authentication pages (login, register,
                  forgot-password), a protected dashboard layout with a sidebar, and React Query hooks
                  for data fetching. Styled with Tailwind CSS and shadcn/ui components.
                </p>

                <h3 className="text-xl font-semibold tracking-tight mt-6 mb-3">
                  Admin Panel (<code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">apps/admin</code>)
                </h3>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  A Filament-inspired admin dashboard. Developers define resources in TypeScript
                  (columns, form fields, filters, actions) and get full CRUD pages with data tables,
                  forms, and dashboard widgets. The admin panel features:
                </p>
                <ul className="space-y-2.5 mb-4">
                  {[
                    'Resource-based page generation from TypeScript definitions',
                    'Server-side paginated DataTable with sorting, filtering, and search',
                    'Form builder with text, number, select, date, toggle, and textarea fields',
                    'Dashboard with stats cards, charts (Recharts), and recent activity widgets',
                    'Auto-generated sidebar navigation from registered resources',
                    'Dark theme matching the GORM Studio aesthetic',
                  ].map((item) => (
                    <li key={item} className="flex items-start gap-2.5 text-[14px] text-muted-foreground">
                      <span className="text-primary mt-1">&#10003;</span>
                      {item}
                    </li>
                  ))}
                </ul>

                <h3 className="text-xl font-semibold tracking-tight mt-6 mb-3">
                  Data Fetching Pattern
                </h3>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  All API communication flows through React Query hooks. Components never call
                  fetch or Axios directly. This ensures consistent caching, background refetching,
                  loading states, and cache invalidation on mutations.
                </p>
                <CodeBlock language="tsx" filename="hooks/use-posts.ts" code={`// Generated React Query hook
export function usePosts({ page, search, sortBy, sortOrder } = {}) {
  return useQuery<PostsResponse>({
    queryKey: ["posts", { page, search, sortBy, sortOrder }],
    queryFn: async () => {
      const { data } = await apiClient.get(\`/api/posts?\${params}\`);
      return data;
    },
  });
}

export function useCreatePost() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (input) => apiClient.post("/api/posts", input),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["posts"] });
    },
  });
}`} />
              </div>

              {/* Shared Package */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Shared Package
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  The <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">packages/shared</code> directory
                  is the bridge between the Go backend and the TypeScript frontends. It contains three
                  categories of shared code:
                </p>
                <div className="rounded-lg border border-border/30 bg-card/30 overflow-hidden">
                  <table className="w-full text-sm">
                    <thead>
                      <tr className="border-b border-border/30 bg-accent/20">
                        <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Directory</th>
                        <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Contents</th>
                        <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Example</th>
                      </tr>
                    </thead>
                    <tbody className="text-muted-foreground">
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-mono text-xs">schemas/</td>
                        <td className="px-4 py-2.5">Zod validation schemas</td>
                        <td className="px-4 py-2.5 font-mono text-xs">CreatePostSchema, UpdatePostSchema</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-mono text-xs">types/</td>
                        <td className="px-4 py-2.5">TypeScript interfaces</td>
                        <td className="px-4 py-2.5 font-mono text-xs">{`interface Post { id: number; title: string; ... }`}</td>
                      </tr>
                      <tr>
                        <td className="px-4 py-2.5 font-mono text-xs">constants/</td>
                        <td className="px-4 py-2.5">Shared constants</td>
                        <td className="px-4 py-2.5 font-mono text-xs">API_ROUTES, ROLES, APP_CONFIG</td>
                      </tr>
                    </tbody>
                  </table>
                </div>
                <p className="text-sm text-muted-foreground/60 mt-3">
                  When you run <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">grit generate resource</code> or <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">grit sync</code>,
                  these files are auto-generated from your Go models, keeping the frontend in sync with the backend.
                </p>
              </div>

              {/* Request Flow */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Request Flow
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Here is the complete path of a typical authenticated API request, from the
                  user clicking a button to the data appearing on screen:
                </p>
                <div className="space-y-3">
                  {[
                    { step: '1', title: 'User interaction', desc: 'User clicks "Save" in the admin panel form. React calls the useCreatePost() mutation hook.' },
                    { step: '2', title: 'API client', desc: 'The hook calls apiClient.post("/api/posts", data). The browser automatically sends the HttpOnly grit_access cookie that the API set on login — JavaScript never reads the token, so XSS cannot exfiltrate it.' },
                    { step: '3', title: 'CORS middleware', desc: 'The Go API receives the request. CORS middleware validates the origin (localhost:3001) and allows it.' },
                    { step: '4', title: 'Auth middleware', desc: 'JWT middleware extracts the token, validates its signature and expiry, and attaches the authenticated user to the Gin context.' },
                    { step: '5', title: 'Handler', desc: 'The PostHandler.Create handler parses the JSON body using ShouldBindJSON, validates required fields via struct binding tags, and calls the service.' },
                    { step: '6', title: 'Service', desc: 'The PostService.Create method runs business logic (if any), then calls GORM to insert the record into PostgreSQL.' },
                    { step: '7', title: 'Database', desc: 'GORM translates the struct into an INSERT query, executes it against PostgreSQL, and populates the ID and timestamps.' },
                    { step: '8', title: 'Response', desc: 'The handler returns a JSON response: { "data": { ... }, "message": "Post created successfully" } with HTTP 201.' },
                    { step: '9', title: 'Cache invalidation', desc: 'React Query receives the response, triggers onSuccess, and invalidates the ["posts"] query key. The list refetches automatically.' },
                  ].map((item) => (
                    <div key={item.step} className="flex gap-3 p-3 rounded-lg border border-border/30 bg-card/30">
                      <div className="flex h-7 w-7 shrink-0 items-center justify-center rounded-md bg-primary/10 border border-primary/15 text-xs font-mono font-semibold text-primary">
                        {item.step}
                      </div>
                      <div>
                        <p className="text-sm font-medium">{item.title}</p>
                        <p className="text-xs text-muted-foreground/70 mt-0.5">{item.desc}</p>
                      </div>
                    </div>
                  ))}
                </div>
              </div>

              {/* GORM Studio */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  GORM Studio
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  GORM Studio is a full-featured visual database browser and editor embedded directly into the Go API.
                  It mounts at <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">/studio</code> and
                  gives you a web UI to browse, query, import/export data, generate Go models, and inspect
                  your database tables. All registered GORM models appear automatically.
                </p>
                <CodeBlock language="go" filename="internal/routes/routes.go" code={`// GORM Studio is mounted with all registered models
studio.Mount(router, db, []interface{}{
    &models.User{},
    &models.Post{},  // auto-injected by grit generate
    /* grit:studio */
}, studio.Config{Prefix: "/studio"})`} />
                <p className="text-sm text-muted-foreground/60 mt-3">
                  When you generate a new resource, the CLI automatically injects the model into
                  the GORM Studio mount call using the <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">/* grit:studio */</code> marker.
                  See the <Link href="/docs/infrastructure/database" className="text-primary hover:underline">Database</Link> page
                  for full GORM Studio documentation.
                </p>
              </div>

              {/* API Response Format */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  API Response Format
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  All API endpoints follow a consistent JSON response format. This predictability
                  allows React Query hooks and error handlers to work generically across every resource.
                </p>
                <div className="grid gap-4 sm:grid-cols-2">
                  <CodeBlock language="json" filename="success (single item)" code={`{
  "data": { ... },
  "message": "Post created"
}`} />
                  <CodeBlock language="json" filename="success (paginated list)" code={`{
  "data": [ ... ],
  "meta": {
    "total": 100,
    "page": 1,
    "page_size": 20,
    "pages": 5
  }
}`} />
                  <div className="sm:col-span-2">
                    <CodeBlock language="json" filename="error" code={`{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Title is required",
    "details": { "title": "This field is required" }
  }
}`} />
                  </div>
                </div>
              </div>

              {/* How Layers Communicate */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Authentication Flow
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Grit uses JWT (JSON Web Tokens) for authentication. The flow works as follows:
                </p>
                <ol className="space-y-2.5 mb-4 list-decimal list-inside">
                  {[
                    'User submits login form with email and password',
                    'Go API validates credentials and sets two HttpOnly cookies: grit_access (15 min, Secure, SameSite=Lax) and grit_refresh (7 days, Path=/api/auth)',
                    'Frontend uses withCredentials: true on its fetch/axios client — the browser attaches the cookies automatically. The JS code never sees the tokens, so XSS cannot read them.',
                    'When the access token expires, the response interceptor catches the 401, calls /api/auth/refresh once (browser sends grit_refresh), the API issues a fresh grit_access cookie, and the original request is retried.',
                    'Protected routes in the frontend check the /api/auth/me response (or a server-side check) and redirect to /login if missing.',
                    'The Go API auth middleware reads the token from the grit_access cookie (falling back to Authorization: Bearer for native mobile/desktop clients), validates the signature and expiry, and attaches the user to the Gin context.',
                    'Role-based middleware (RequireRole("admin")) checks the user\'s role from the context. CSRF middleware enforces an X-CSRF-Token header on state-changing methods to mitigate cookie-auth CSRF.',
                  ].map((item, i) => (
                    <li key={i} className="text-[14px] text-muted-foreground pl-1">
                      {item}
                    </li>
                  ))}
                </ol>
              </div>

              {/* Infrastructure */}
              <div className="mb-10">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Infrastructure Services
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Docker Compose provides all the infrastructure services needed for local development.
                  In production, you can swap these for managed services (e.g., AWS RDS for PostgreSQL,
                  ElastiCache for Redis, S3 for file storage).
                </p>
                <div className="rounded-lg border border-border/30 bg-card/30 overflow-hidden">
                  <table className="w-full text-sm">
                    <thead>
                      <tr className="border-b border-border/30 bg-accent/20">
                        <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Service</th>
                        <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Port</th>
                        <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Purpose</th>
                      </tr>
                    </thead>
                    <tbody className="text-muted-foreground">
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-mono text-xs">PostgreSQL 16</td>
                        <td className="px-4 py-2.5 font-mono text-xs">5432</td>
                        <td className="px-4 py-2.5">Primary database (GORM auto-migrates on startup)</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-mono text-xs">Redis 7</td>
                        <td className="px-4 py-2.5 font-mono text-xs">6379</td>
                        <td className="px-4 py-2.5">Caching, session storage, job queue backend (asynq)</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-mono text-xs">MinIO</td>
                        <td className="px-4 py-2.5 font-mono text-xs">9000 / 9001</td>
                        <td className="px-4 py-2.5">S3-compatible object storage for file uploads</td>
                      </tr>
                      <tr>
                        <td className="px-4 py-2.5 font-mono text-xs">Mailhog</td>
                        <td className="px-4 py-2.5 font-mono text-xs">1025 / 8025</td>
                        <td className="px-4 py-2.5">Local email testing (SMTP on 1025, web UI on 8025)</td>
                      </tr>
                    </tbody>
                  </table>
                </div>
              </div>
            </div>

            {/* Nav */}
            <div className="flex items-center justify-between pt-6 border-t border-border/30">
              <Button variant="ghost" size="sm" asChild className="text-muted-foreground/60 hover:text-foreground">
                <Link href="/docs/getting-started/configuration" className="gap-1.5">
                  <ArrowLeft className="h-3.5 w-3.5" />
                  Configuration
                </Link>
              </Button>
              <Button variant="ghost" size="sm" asChild className="text-muted-foreground/60 hover:text-foreground">
                <Link href="/docs/concepts/cli" className="gap-1.5">
                  CLI Commands
                  <ArrowRight className="h-3.5 w-3.5" />
                </Link>
              </Button>
            </div>
          </div>
        </div>
      </main>
    </div>
  )
}
