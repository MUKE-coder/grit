// The guide set at /guides: problem-first articles for somebody searching for
// their problem rather than for Grit.
//
// This file holds the index-card data and the version each page was last
// verified against. It deliberately does NOT hold the pages' prose or their
// metadata title and description: those live in each page.tsx, where both of
// the site's search indexers can see them. A section that exists only as
// config is invisible to search, silently, and that has happened here before.
//
// `updated` is the Grit version whose output the page's commands were last run
// against, per page. It is not GRIT_VERSION, because filling it from the
// current version would claim every page is re-verified on every release.
// Re-run a page's commands, then bump its line.

export type Audience = 'beginner' | 'senior'

export interface Guide {
  /** Path under /guides, without the prefix. */
  slug: string
  /** The card title on the index. Fuller than the page's SEO title. */
  title: string
  /** One sentence on the index card. */
  blurb: string
  audience: Audience
  reading: string
  /** The Grit version the commands were last run against, or null if planned. */
  updated: string | null
}

export const GUIDES: Guide[] = [
  // Beginners, in the review's publishing order: the pages people hit when
  // they are stuck come first.
  {
    slug: 'first-go-react-app',
    title: 'Build your first full-stack app with Go and React',
    blurb:
      'Start to finish in about thirty minutes, with a checkpoint you can verify after every step.',
    audience: 'beginner',
    reading: '25 min',
    updated: 'v3.560.0',
  },
  {
    slug: 'common-go-docker-errors',
    title: 'Fix the errors you hit starting a Go and Docker project',
    blurb:
      'One section per error message, with the cause and the fix for Windows, macOS and Linux.',
    audience: 'beginner',
    reading: '14 min',
    updated: 'v3.560.0',
  },
  {
    slug: 'docker-compose-postgres-redis',
    title: 'Run Postgres, Redis and MinIO locally with Docker Compose',
    blurb:
      'What a container actually is, the compose file line by line, and the three things that go wrong on Windows.',
    audience: 'beginner',
    reading: '16 min',
    updated: 'v3.560.0',
  },
  {
    slug: 'nextjs-go-backend',
    title: 'Connect a Next.js frontend to a Go backend',
    blurb:
      'CORS, cookies and SameSite: why the login works in Postman and the browser throws the request away.',
    audience: 'beginner',
    reading: '18 min',
    updated: 'v3.560.0',
  },
  {
    slug: 'golang-jwt-authentication',
    title: 'JWT authentication in Go, explained',
    blurb:
      'Access tokens, refresh tokens, HttpOnly cookies, and the four mistakes that make an auth system fake.',
    audience: 'beginner',
    reading: '20 min',
    updated: 'v3.560.0',
  },
  {
    slug: 'gorm-tutorial-relationships',
    title: 'GORM relationships: belongs to, has many, many to many',
    blurb: 'Each relation drawn, the struct that produces it, and the SQL GORM writes for it.',
    audience: 'beginner',
    reading: '18 min',
    updated: 'v3.560.0',
  },
  {
    slug: 'gin-rest-api-tutorial',
    title: 'Build a REST API in Go with Gin',
    blurb:
      'Routing, JSON binding, validation and error responses, written out by hand before anything generates them.',
    audience: 'beginner',
    reading: '20 min',
    updated: 'v3.560.0',
  },
  {
    slug: 'deploy-go-app-vps-https',
    title: 'Deploy a Go and React app to a cheap VPS with HTTPS',
    blurb:
      'What a VPS, systemd and a reverse proxy are, DNS, certificates, and the first deploy end to end.',
    audience: 'beginner',
    reading: '22 min',
    updated: 'v3.560.0',
  },
  {
    slug: 'go-file-upload-s3-presigned',
    title: 'Upload files and images in Go to S3, R2 or MinIO',
    blurb: 'The presigned flow drawn out, why the browser uploads directly, and image renditions.',
    audience: 'beginner',
    reading: '17 min',
    updated: 'v3.560.0',
  },
  {
    slug: 'build-admin-dashboard-go-react',
    title: 'Build an admin dashboard without writing it from scratch',
    blurb: 'Tables, forms, filters and widgets, and what to declare rather than build.',
    audience: 'beginner',
    reading: '15 min',
    updated: 'v3.560.0',
  },
  {
    slug: 'rbac-explained-golang',
    title: 'Roles and permissions explained, with Go examples',
    blurb: 'What RBAC is in plain English, where the check belongs, and the row rule people miss.',
    audience: 'beginner',
    reading: '14 min',
    updated: 'v3.560.0',
  },
  {
    slug: 'go-background-jobs-asynq',
    title: 'Background jobs in Go: send email without slowing the API',
    blurb: 'Why a queue exists, what a retry costs, and testing mail locally before it is real.',
    audience: 'beginner',
    reading: '16 min',
    updated: 'v3.560.0',
  },
  {
    slug: 'go-vs-nodejs-backend',
    title: 'Go or Node.js for your backend: an honest comparison',
    blurb:
      'Where each one genuinely wins, what Go costs a JavaScript developer, and how to learn it fast.',
    audience: 'beginner',
    reading: '15 min',
    updated: 'v3.560.0',
  },
  {
    slug: 'build-saas-mvp-weekend',
    title: 'Build a SaaS MVP in a weekend with Go and React',
    blurb: 'A two-day plan: auth, tenants, billing, a landing page and a deploy, in that order.',
    audience: 'beginner',
    reading: '18 min',
    updated: 'v3.560.0',
  },
  {
    slug: 'ai-assistant-build-go-app',
    title: 'Use an AI coding assistant to build a Go app properly',
    blurb:
      'What to give it before it writes anything, the prompts that work, and how to review what comes back.',
    audience: 'beginner',
    reading: '16 min',
    updated: 'v3.560.0',
  },

  // Seniors.
  {
    slug: 'go-project-structure',
    title: 'Go project structure for production APIs',
    blurb:
      'Handler, service, repository: the layout, the trade-offs against hexagonal and DDD, and where transactions live.',
    audience: 'senior',
    reading: '20 min',
    updated: 'v3.560.0',
  },
  {
    slug: 'go-transactional-outbox',
    title: 'The transactional outbox pattern in Go with Postgres',
    blurb: 'Why dual writes lose events, a minimal outbox, idempotency keys and parked failures.',
    audience: 'senior',
    reading: '18 min',
    updated: 'v3.560.0',
  },
  {
    slug: 'go-saga-pattern',
    title: 'Sagas in Go: compensating transactions without a workflow engine',
    blurb: 'Do and undo steps, persisted state, and the point at which you should reach for Temporal.',
    audience: 'senior',
    reading: '18 min',
    updated: 'v3.560.0',
  },
  {
    slug: 'go-multi-tenancy',
    title: 'Multi-tenant SaaS in Go: row, schema or database per tenant',
    blurb:
      'A decision matrix, Postgres RLS against application scoping, and the cross-tenant leak test.',
    audience: 'senior',
    reading: '22 min',
    updated: 'v3.560.0',
  },
  {
    slug: 'go-authorization-policies-row-rules',
    title: 'Authorization in Go: policies answer 403, row rules answer 404',
    blurb: 'IDOR, enumeration, and why an out-of-scope row must be indistinguishable from a missing one.',
    audience: 'senior',
    reading: '18 min',
    updated: 'v3.560.0',
  },
  {
    slug: 'hash-chained-audit-log',
    title: 'Tamper-evident audit logs with hash chains in Postgres',
    blurb: 'Chaining, verification, pruning without breaking the chain, and what auditors ask for.',
    audience: 'senior',
    reading: '17 min',
    updated: 'v3.560.0',
  },
  {
    slug: 'gorm-zero-downtime-migrations',
    title: 'Zero-downtime Postgres migrations with GORM',
    blurb: 'Expand and contract, why AutoMigrate stops at development, and the DDL that takes a lock.',
    audience: 'senior',
    reading: '19 min',
    updated: 'v3.560.0',
  },
  {
    slug: 'go-redis-caching-stampede',
    title: 'Caching in Go: cache-aside, jittered TTLs and stampede protection',
    blurb: 'The thundering herd, singleflight, jitter, and invalidation you can reason about.',
    audience: 'senior',
    reading: '17 min',
    updated: 'v3.560.0',
  },
  {
    slug: 'go-money-currency-invoices',
    title: 'Money in Go: minor units, ISO 4217 and gap-free invoice numbers',
    blurb: 'Why float is wrong, what a zero-decimal currency breaks, and row-locked sequences.',
    audience: 'senior',
    reading: '16 min',
    updated: 'v3.560.0',
  },
  {
    slug: 'go-horizontal-scaling',
    title: 'Running Go behind a load balancer',
    blurb:
      'A stateless checklist, WebSockets over a Redis backplane, cron on exactly one node, read replicas.',
    audience: 'senior',
    reading: '20 min',
    updated: 'v3.560.0',
  },
  {
    slug: 'go-webhooks-hmac-idempotency',
    title: 'Webhooks in Go: HMAC verification, idempotency and replay protection',
    blurb: 'Constant-time compare, timestamp tolerance, a dedupe table, and retries that are safe.',
    audience: 'senior',
    reading: '17 min',
    updated: 'v3.560.0',
  },
  {
    slug: 'go-typescript-type-safety',
    title: 'End-to-end type safety between Go and TypeScript without GraphQL',
    blurb: 'tygo, OpenAPI codegen and Connect compared, and what a generated Zod schema buys you.',
    audience: 'senior',
    reading: '18 min',
    updated: 'v3.560.0',
  },
  {
    slug: 'go-full-stack-frameworks-compared',
    title: 'Choosing a Go full-stack approach in 2026',
    blurb:
      'PocketBase, Encore, Goravel, Buffalo and Grit: runtime backend, infrastructure from code, or generated code.',
    audience: 'senior',
    reading: '19 min',
    updated: 'v3.560.0',
  },
  {
    slug: 'owasp-top-10-golang',
    title: 'Securing a Go API against the OWASP Top 10',
    blurb: 'Each risk, the Go code that addresses it, and the check that keeps it addressed in CI.',
    audience: 'senior',
    reading: '24 min',
    updated: 'v3.560.0',
  },
  {
    slug: 'migrate-gin-gorm-to-grit',
    title: 'Move an existing Gin and GORM API into Grit',
    blurb:
      'One resource at a time, with both halves running, and a point at every step where you can stop.',
    audience: 'senior',
    reading: '18 min',
    updated: 'v3.562.0',
  },
  {
    slug: 'go-ai-rag-pgvector',
    title: 'Adding AI to a Go app: streaming, pgvector and RAG',
    blurb: 'Server-sent events, an embeddings job, hybrid search, and where the cost actually goes.',
    audience: 'senior',
    reading: '20 min',
  updated: 'v3.560.0',
  },
]

export const GUIDE_BY_SLUG: Record<string, Guide> = Object.fromEntries(
  GUIDES.map((g) => [g.slug, g]),
)

/** A guide by slug, throwing rather than rendering a page with no furniture. */
export function guide(slug: string): Guide {
  const found = GUIDE_BY_SLUG[slug]
  if (!found) {
    throw new Error(
      `No entry in config/guides.ts for "${slug}". Add one, or the page renders ` +
        `without its reading time and its verified version.`,
    )
  }
  return found
}

export const LIVE_GUIDES = GUIDES.filter((g) => g.updated !== null)
export const PLANNED_GUIDES = GUIDES.filter((g) => g.updated === null)
