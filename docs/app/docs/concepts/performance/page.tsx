import Link from 'next/link'
import { ArrowLeft, ArrowRight, Zap, Server, Globe, Database, Shield, Clock, Image, Code2, Layers, Wind } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'
import { CodeBlock } from '@/components/code-block'
import { getDocMetadata } from '@/config/docs-metadata'

export const metadata = getDocMetadata('/docs/concepts/performance')

export default function PerformancePage() {
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
                Performance
              </h1>
              <p className="text-lg text-muted-foreground leading-relaxed">
                Every Grit project is scaffolded with production-grade performance optimisations
                out of the box — both in the Go API and the Next.js frontends. No extra
                configuration needed.
              </p>
            </div>

            <div className="prose-grit">

              {/* ── BLOCK 1: Backend ─────────────────────────────── */}
              <div className="mb-16">
                <div className="flex items-center gap-3 mb-6">
                  <div className="flex h-10 w-10 items-center justify-center rounded-xl bg-primary/10 border border-primary/20">
                    <Server className="h-5 w-5 text-primary" />
                  </div>
                  <div>
                    <h2 className="text-2xl font-bold tracking-tight">Backend Performance</h2>
                    <p className="text-sm text-muted-foreground">Go API — Gin + GORM</p>
                  </div>
                </div>

                <div className="space-y-10">

                  {/* Gzip */}
                  <div>
                    <div className="flex items-center gap-2 mb-3">
                      <Wind className="h-4 w-4 text-primary/70" />
                      <h3 className="text-lg font-semibold tracking-tight">Gzip Response Compression</h3>
                    </div>
                    <p className="text-muted-foreground leading-relaxed mb-4">
                      All API responses are compressed with Gzip automatically. The middleware
                      checks for <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">Accept-Encoding: gzip</code> and
                      wraps the Gin response writer to compress output at{' '}
                      <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">BestSpeed</code> — the
                      sweet spot between CPU cost and payload size. JSON payloads typically shrink by
                      60–80%, dramatically reducing bandwidth on paginated list endpoints.
                    </p>
                    <CodeBlock language="go" filename="apps/api/internal/middleware/middleware.go" code={`func Gzip() gin.HandlerFunc {
    return func(c *gin.Context) {
        if !strings.Contains(c.GetHeader("Accept-Encoding"), "gzip") {
            c.Next()
            return
        }
        gz, _ := gzip.NewWriterLevel(c.Writer, gzip.BestSpeed)
        defer gz.Close()
        c.Header("Content-Encoding", "gzip")
        c.Header("Vary", "Accept-Encoding")
        c.Writer = &gzipResponseWriter{ResponseWriter: c.Writer, Writer: gz}
        c.Next()
    }
}

// Registered globally — applies to every route
r.Use(middleware.Gzip())`} />
                  </div>

                  {/* Request ID */}
                  <div>
                    <div className="flex items-center gap-2 mb-3">
                      <Code2 className="h-4 w-4 text-primary/70" />
                      <h3 className="text-lg font-semibold tracking-tight">Request ID Tracing</h3>
                    </div>
                    <p className="text-muted-foreground leading-relaxed mb-4">
                      Every request receives a unique{' '}
                      <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">X-Request-ID</code> header.
                      If the upstream proxy or client already sends one it is echoed back; otherwise a new ID is
                      generated from the nanosecond timestamp. The ID is stored in Gin&apos;s context and
                      included in every structured log line — making it trivial to trace a specific
                      request through logs, Pulse, and your reverse proxy.
                    </p>
                    <CodeBlock language="go" filename="apps/api/internal/middleware/middleware.go" code={`func RequestID() gin.HandlerFunc {
    return func(c *gin.Context) {
        id := c.GetHeader("X-Request-ID")
        if id == "" {
            id = fmt.Sprintf("%d-%d", time.Now().UnixNano(), rand.Int63())
        }
        c.Set("request_id", id)
        c.Header("X-Request-ID", id)
        c.Next()
    }
}

// Logger includes the request_id in every log line
log.Printf("[%s] %s %s %d", requestID, c.Request.Method, c.Request.URL.Path, statusCode)`} />
                  </div>

                  {/* DB Connection Pool */}
                  <div>
                    <div className="flex items-center gap-2 mb-3">
                      <Database className="h-4 w-4 text-primary/70" />
                      <h3 className="text-lg font-semibold tracking-tight">Database Connection Pool</h3>
                    </div>
                    <p className="text-muted-foreground leading-relaxed mb-4">
                      The scaffold configures GORM&apos;s underlying{' '}
                      <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">database/sql</code> pool
                      with four tuned settings. Without these, Go&apos;s default behaviour is to open
                      unlimited connections and keep them forever — which exhausts Postgres under load
                      and causes stale connection errors after network interruptions.
                    </p>
                    <CodeBlock language="go" filename="apps/api/internal/config/database.go" code={`sqlDB, _ := db.DB()

sqlDB.SetMaxIdleConns(10)          // keep 10 idle connections warm
sqlDB.SetMaxOpenConns(100)         // never open more than 100 connections
sqlDB.SetConnMaxLifetime(30 * time.Minute) // recycle to prevent stale connections
sqlDB.SetConnMaxIdleTime(10 * time.Minute) // evict idle connections sooner`} />
                    <div className="mt-4 grid grid-cols-1 sm:grid-cols-2 gap-3">
                      {[
                        { label: 'MaxIdleConns: 10', desc: 'Warm pool ready for burst traffic without cold-start latency' },
                        { label: 'MaxOpenConns: 100', desc: 'Prevents connection exhaustion on the Postgres server' },
                        { label: 'ConnMaxLifetime: 30m', desc: 'Recycles connections — fixes stale connection errors after failover' },
                        { label: 'ConnMaxIdleTime: 10m', desc: 'Evicts idle connections to free Postgres resources during quiet periods' },
                      ].map((item) => (
                        <div key={item.label} className="p-3 rounded-lg border border-border/30 bg-card/30">
                          <p className="text-xs font-mono font-semibold text-primary/80 mb-1">{item.label}</p>
                          <p className="text-xs text-muted-foreground/70 leading-relaxed">{item.desc}</p>
                        </div>
                      ))}
                    </div>
                  </div>

                  {/* Cache-Control */}
                  <div>
                    <div className="flex items-center gap-2 mb-3">
                      <Clock className="h-4 w-4 text-primary/70" />
                      <h3 className="text-lg font-semibold tracking-tight">Cache-Control Headers on Public Endpoints</h3>
                    </div>
                    <p className="text-muted-foreground leading-relaxed mb-4">
                      Public read endpoints (blog list and single post by slug) emit{' '}
                      <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">Cache-Control</code> headers
                      so CDNs and edge caches (Cloudflare, Vercel Edge, Nginx) can serve them without
                      hitting the Go API at all. The single-post endpoint uses a longer TTL because
                      published content changes infrequently.
                    </p>
                    <CodeBlock language="go" filename="apps/api/internal/handlers/blog_handler.go" code={`// ListPublished — paginated public blog list
func (h *BlogHandler) ListPublished(c *gin.Context) {
    // ... query
    c.Header("Cache-Control", "public, max-age=300")  // 5 minutes
    c.JSON(http.StatusOK, gin.H{"data": blogs, "meta": meta})
}

// GetBySlug — single published post
func (h *BlogHandler) GetBySlug(c *gin.Context) {
    // ... query
    c.Header("Cache-Control", "public, max-age=3600") // 1 hour
    c.JSON(http.StatusOK, gin.H{"data": blog})
}`} />
                  </div>

                  {/* Presigned uploads */}
                  <div>
                    <div className="flex items-center gap-2 mb-3">
                      <Zap className="h-4 w-4 text-primary/70" />
                      <h3 className="text-lg font-semibold tracking-tight">Presigned URL Uploads (Bypass the API)</h3>
                    </div>
                    <p className="text-muted-foreground leading-relaxed mb-4">
                      File uploads never pass through the Go API. The browser asks the API for a
                      presigned PUT URL, then uploads the binary directly to S3/R2/MinIO. This
                      eliminates request body size limits imposed by reverse proxies (Nginx, Traefik),
                      removes Go memory pressure from large uploads, and allows XHR progress tracking.
                    </p>
                    <div className="rounded-xl border border-border/40 bg-card/50 p-5 mb-4">
                      <div className="space-y-3">
                        {[
                          { step: '1', label: 'Browser → POST /api/uploads/presign', desc: 'Returns a time-limited presigned PUT URL + storage key + public URL' },
                          { step: '2', label: 'Browser → PUT directly to R2/S3/MinIO', desc: 'Binary is sent straight to storage — Go API is not involved. XHR tracks progress.' },
                          { step: '3', label: 'Browser → POST /api/uploads/complete', desc: 'API creates the Upload DB record and enqueues thumbnail processing job' },
                        ].map((item) => (
                          <div key={item.step} className="flex items-start gap-3">
                            <div className="flex h-6 w-6 shrink-0 items-center justify-center rounded-md bg-primary/10 border border-primary/15 text-xs font-mono font-semibold text-primary mt-0.5">
                              {item.step}
                            </div>
                            <div>
                              <p className="text-xs font-mono font-semibold text-foreground/80 mb-0.5">{item.label}</p>
                              <p className="text-xs text-muted-foreground/60 leading-relaxed">{item.desc}</p>
                            </div>
                          </div>
                        ))}
                      </div>
                    </div>
                  </div>

                  {/* Background Jobs */}
                  <div>
                    <div className="flex items-center gap-2 mb-3">
                      <Layers className="h-4 w-4 text-primary/70" />
                      <h3 className="text-lg font-semibold tracking-tight">Async Background Jobs</h3>
                    </div>
                    <p className="text-muted-foreground leading-relaxed mb-4">
                      Slow operations — image thumbnail generation, welcome emails, PDF reports,
                      webhook delivery — are pushed to a Redis-backed asynq queue so the HTTP
                      handler returns immediately. Workers process jobs concurrently in separate
                      goroutines. Retries, dead-letter queues, and job monitoring are available
                      through the built-in admin dashboard.
                    </p>
                    <CodeBlock language="go" filename="apps/api/internal/handlers/upload_handler.go" code={`// Handler returns 201 in < 1ms. Thumbnail is generated asynchronously.
func (h *UploadHandler) CompleteUpload(c *gin.Context) {
    // ... save Upload record
    if strings.HasPrefix(upload.MimeType, "image/") {
        h.Jobs.EnqueueImageProcessing(upload.ID, upload.StorageKey)
    }
    c.JSON(http.StatusCreated, gin.H{"data": upload})
}`} />
                  </div>

                  {/* Redis Caching */}
                  <div>
                    <div className="flex items-center gap-2 mb-3">
                      <Database className="h-4 w-4 text-amber-500/70" />
                      <h3 className="text-lg font-semibold tracking-tight">Redis Caching Layer</h3>
                    </div>
                    <p className="text-muted-foreground leading-relaxed mb-4">
                      The scaffold includes a Redis cache service and a Gin middleware that caches
                      entire API responses by URL. Hot endpoints like product lists or homepage data
                      are served from memory in under a millisecond. The cache middleware skips
                      authenticated routes automatically so user-specific data is never cached.
                    </p>
                    <CodeBlock language="go" filename="apps/api/internal/routes/routes.go" code={`// Cache public product list for 5 minutes
public.GET("/products", middleware.Cache(5*time.Minute), productHandler.List)

// Authenticated routes — cache skipped automatically
protected.GET("/orders", orderHandler.List)`} />
                  </div>

                  {/* Rate Limiting + Security */}
                  <div>
                    <div className="flex items-center gap-2 mb-3">
                      <Shield className="h-4 w-4 text-primary/70" />
                      <h3 className="text-lg font-semibold tracking-tight">Rate Limiting & WAF (Sentinel)</h3>
                    </div>
                    <p className="text-muted-foreground leading-relaxed mb-4">
                      Sentinel is Grit&apos;s built-in security suite. It acts as a Web Application Firewall,
                      rate limiter, and brute-force shield — protecting the API from abuse that could
                      degrade performance for legitimate users. Internal dashboards (<code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">/pulse</code>,{' '}
                      <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">/sentinel</code>,{' '}
                      <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">/docs</code>,{' '}
                      <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">/studio</code>) are
                      excluded from rate limiting so health checks never trigger false positives.
                    </p>
                    <CodeBlock language="go" filename="apps/api/internal/config/sentinel.go" code={`sentinel.Init(sentinel.Config{
    RateLimit: &sentinel.RateLimitConfig{
        Enabled:       true,
        RequestsPerSecond: 10,
        Burst:         20,
        ExcludeRoutes: []string{"/pulse/*", "/sentinel/*", "/docs/*", "/studio/*"},
    },
    WAF:            &sentinel.WAFConfig{Enabled: true},
    BruteForce:     &sentinel.BruteForceConfig{Enabled: true, MaxAttempts: 5},
})`} />
                  </div>

                </div>
              </div>

              {/* ── BLOCK 2: Frontend ─────────────────────────────── */}
              <div className="mb-12">
                <div className="flex items-center gap-3 mb-6">
                  <div className="flex h-10 w-10 items-center justify-center rounded-xl bg-emerald-500/10 border border-emerald-500/20">
                    <Globe className="h-5 w-5 text-emerald-500" />
                  </div>
                  <div>
                    <h2 className="text-2xl font-bold tracking-tight">Frontend Performance</h2>
                    <p className="text-sm text-muted-foreground">Next.js 15 — App Router + React</p>
                  </div>
                </div>

                <div className="space-y-10">

                  {/* Server Components */}
                  <div>
                    <div className="flex items-center gap-2 mb-3">
                      <Server className="h-4 w-4 text-emerald-500/70" />
                      <h3 className="text-lg font-semibold tracking-tight">React Server Components by Default</h3>
                    </div>
                    <p className="text-muted-foreground leading-relaxed mb-4">
                      Every page in the Next.js web app is a React Server Component unless it
                      explicitly opts in with <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">&apos;use client&apos;</code>.
                      Server Components fetch data directly on the server before streaming HTML to
                      the browser — zero JavaScript for data fetching, no loading spinners on initial
                      render, and no client-server waterfalls. Only interactive UI (forms, modals,
                      dropdowns) becomes a Client Component.
                    </p>
                    <CodeBlock language="tsx" filename="apps/web/app/blog/page.tsx" code={`// Server Component — fetches on the server, streams HTML
// No JS bundle cost. No useEffect. No loading state.
export default async function BlogPage() {
  const posts = await fetch(\`\${process.env.API_URL}/api/blogs/published\`, {
    next: { revalidate: 300 }, // ISR: revalidate every 5 minutes
  }).then((r) => r.json())

  return <BlogList posts={posts.data} />
}`} />
                  </div>

                  {/* Incremental Static Regen */}
                  <div>
                    <div className="flex items-center gap-2 mb-3">
                      <Zap className="h-4 w-4 text-emerald-500/70" />
                      <h3 className="text-lg font-semibold tracking-tight">Incremental Static Regeneration (ISR)</h3>
                    </div>
                    <p className="text-muted-foreground leading-relaxed mb-4">
                      Public content pages (blogs, product catalogs, landing pages) use ISR via
                      Next.js&apos;s{' '}
                      <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">next: {'{ revalidate }'}</code> option.
                      The page is rendered once and cached at the CDN edge. Subsequent visitors get
                      the cached HTML in milliseconds. The cache revalidates in the background after
                      the TTL expires — users always get fast responses even while content refreshes.
                    </p>
                    <CodeBlock language="tsx" filename="apps/web/app/blog/[slug]/page.tsx" code={`// generateStaticParams pre-builds all published post pages at build time
export async function generateStaticParams() {
  const posts = await fetch(\`\${process.env.API_URL}/api/blogs/published\`)
    .then((r) => r.json())
  return posts.data.map((p: { slug: string }) => ({ slug: p.slug }))
}

export default async function PostPage({ params }: { params: { slug: string } }) {
  const post = await fetch(\`\${process.env.API_URL}/api/blogs/slug/\${params.slug}\`, {
    next: { revalidate: 3600 }, // re-check every hour
  }).then((r) => r.json())

  return <PostContent post={post.data} />
}`} />
                  </div>

                  {/* React Query */}
                  <div>
                    <div className="flex items-center gap-2 mb-3">
                      <Database className="h-4 w-4 text-emerald-500/70" />
                      <h3 className="text-lg font-semibold tracking-tight">React Query — Smart Client Caching</h3>
                    </div>
                    <p className="text-muted-foreground leading-relaxed mb-4">
                      All admin panel data fetching uses React Query (TanStack Query). Responses are
                      cached in memory so navigating back to a previously visited page is instant.
                      Mutations automatically invalidate the relevant query so lists refresh without
                      a full page reload. Generated hooks follow a consistent pattern across all resources.
                    </p>
                    <CodeBlock language="tsx" filename="apps/admin/hooks/use-products.ts (generated)" code={`// Generated by: grit generate resource Product
export function useProducts(page = 1) {
  return useQuery({
    queryKey: ['products', page],
    queryFn: () => apiClient.get(\`/api/products?page=\${page}\`).then((r) => r.data),
    staleTime: 30_000, // treat data as fresh for 30 s
  })
}

export function useCreateProduct() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (data) => apiClient.post('/api/products', data).then((r) => r.data),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['products'] }),
  })
}`} />
                  </div>

                  {/* Image optimisation */}
                  <div>
                    <div className="flex items-center gap-2 mb-3">
                      <Image className="h-4 w-4 text-emerald-500/70" />
                      <h3 className="text-lg font-semibold tracking-tight">Next.js Image Optimisation</h3>
                    </div>
                    <p className="text-muted-foreground leading-relaxed mb-4">
                      The scaffolded web app uses the{' '}
                      <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">next/image</code> component
                      throughout. Images are automatically converted to WebP/AVIF, served at the
                      correct size for the user&apos;s device, lazy-loaded by default, and cached at the
                      CDN layer. The{' '}
                      <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">next.config</code> is
                      pre-configured with{' '}
                      <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">remotePatterns</code> for
                      your storage domain (R2, S3, MinIO) so remote images work without{' '}
                      <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">unoptimized</code>.
                    </p>
                    <CodeBlock language="tsx" filename="apps/web/app/blog/[slug]/page.tsx" code={`import Image from 'next/image'

<Image
  src={post.image}           // remote R2 / S3 URL
  alt={post.title}
  width={1200}
  height={630}
  priority                   // eager-load above-the-fold hero
  className="rounded-xl object-cover"
/>`} />
                  </div>

                  {/* Turborepo */}
                  <div>
                    <div className="flex items-center gap-2 mb-3">
                      <Zap className="h-4 w-4 text-amber-500/70" />
                      <h3 className="text-lg font-semibold tracking-tight">Turborepo Build Cache</h3>
                    </div>
                    <p className="text-muted-foreground leading-relaxed mb-4">
                      Grit projects are managed by Turborepo. Build outputs are hashed and cached
                      locally (and optionally remotely). If neither the source nor its dependencies
                      changed, Turbo replays the cached output in milliseconds instead of re-running
                      the build. On a typical GritCMS-scale project this cuts CI build time from
                      4+ minutes to under 30 seconds on the second run.
                    </p>
                    <CodeBlock language="json" filename="turbo.json" code={`{
  "tasks": {
    "build": {
      "dependsOn": ["^build"],
      "outputs": [".next/**", "!.next/cache/**", "dist/**"]
    },
    "dev": {
      "cache": false,
      "persistent": true
    }
  }
}`} />
                  </div>

                  {/* Code splitting */}
                  <div>
                    <div className="flex items-center gap-2 mb-3">
                      <Layers className="h-4 w-4 text-emerald-500/70" />
                      <h3 className="text-lg font-semibold tracking-tight">Automatic Code Splitting</h3>
                    </div>
                    <p className="text-muted-foreground leading-relaxed mb-4">
                      Next.js App Router automatically splits the JavaScript bundle per route — users
                      only download code for the page they are visiting. Heavy admin components
                      (rich-text editor, chart library, data grid) are dynamically imported with{' '}
                      <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">next/dynamic</code> so
                      they don&apos;t inflate the initial bundle. Combined with Server Components, the
                      JS sent to the browser is kept to an absolute minimum.
                    </p>
                    <CodeBlock language="tsx" filename="apps/admin/components/rich-text-editor.tsx" code={`import dynamic from 'next/dynamic'

// Loaded only when the component is actually rendered
const RichTextEditor = dynamic(() => import('@/components/editor'), {
  loading: () => <div className="h-40 animate-pulse rounded-lg bg-accent/30" />,
  ssr: false, // editor requires browser APIs
})`} />
                  </div>

                </div>
              </div>

              {/* Summary table */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  What You Get Out of the Box
                </h2>
                <div className="rounded-xl border border-border/40 bg-card/50 overflow-hidden">
                  <table className="w-full text-sm">
                    <thead>
                      <tr className="border-b border-border/30 bg-accent/20">
                        <th className="px-4 py-3 text-left text-xs font-semibold text-muted-foreground uppercase tracking-wider">Optimisation</th>
                        <th className="px-4 py-3 text-left text-xs font-semibold text-muted-foreground uppercase tracking-wider">Layer</th>
                        <th className="px-4 py-3 text-left text-xs font-semibold text-muted-foreground uppercase tracking-wider">Benefit</th>
                      </tr>
                    </thead>
                    <tbody className="divide-y divide-border/20">
                      {[
                        { name: 'Gzip middleware', layer: 'Backend', benefit: '60-80% smaller API responses' },
                        { name: 'Request ID tracing', layer: 'Backend', benefit: 'Correlate logs across services' },
                        { name: 'Connection pool tuning', layer: 'Backend', benefit: 'No stale connections under load' },
                        { name: 'Cache-Control headers', layer: 'Backend', benefit: 'CDN-cacheable public endpoints' },
                        { name: 'Presigned URL uploads', layer: 'Backend', benefit: 'Bypass API for large files' },
                        { name: 'Background jobs (asynq)', layer: 'Backend', benefit: 'Non-blocking async operations' },
                        { name: 'Redis response cache', layer: 'Backend', benefit: 'Sub-millisecond hot reads' },
                        { name: 'Sentinel rate limiting', layer: 'Backend', benefit: 'Protect API from abuse' },
                        { name: 'Server Components', layer: 'Frontend', benefit: 'Zero JS for data fetching' },
                        { name: 'ISR / revalidate', layer: 'Frontend', benefit: 'CDN-cached public pages' },
                        { name: 'React Query caching', layer: 'Frontend', benefit: 'Instant back-navigation in admin' },
                        { name: 'next/image', layer: 'Frontend', benefit: 'WebP, lazy load, correct sizing' },
                        { name: 'Turborepo cache', layer: 'Frontend', benefit: 'Fast CI and local builds' },
                        { name: 'Code splitting', layer: 'Frontend', benefit: 'Minimal JS per route' },
                      ].map((row) => (
                        <tr key={row.name}>
                          <td className="px-4 py-2.5 font-mono text-[12px] text-foreground/80">{row.name}</td>
                          <td className="px-4 py-2.5">
                            <span className={`inline-flex items-center rounded-md px-2 py-0.5 text-[10px] font-medium ${
                              row.layer === 'Backend'
                                ? 'bg-primary/10 text-primary border border-primary/20'
                                : 'bg-emerald-500/10 text-emerald-500 border border-emerald-500/20'
                            }`}>
                              {row.layer}
                            </span>
                          </td>
                          <td className="px-4 py-2.5 text-xs text-muted-foreground/70">{row.benefit}</td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>
              </div>

              {/* Nav */}
              <div className="flex items-center justify-between pt-6 border-t border-border/30">
                <Button variant="ghost" size="sm" asChild className="text-muted-foreground/60 hover:text-foreground">
                  <Link href="/docs/concepts/naming-conventions" className="gap-1.5">
                    <ArrowLeft className="h-3.5 w-3.5" />
                    Naming Conventions
                  </Link>
                </Button>
                <Button variant="ghost" size="sm" asChild className="text-muted-foreground/60 hover:text-foreground">
                  <Link href="/docs/backend/middleware" className="gap-1.5">
                    Middleware
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
