'use client'

import { useState } from 'react'
import type React from 'react'
import { Database, HardDrive, Zap, Check, AlertTriangle } from 'lucide-react'
import { CodeBlock } from '@/components/code-block'

/**
 * Databases, file storage and Redis — what Grit actually supports.
 *
 * Every provider, env var and default here was read out of the scaffold
 * templates. Two things this section deliberately does NOT claim: MySQL (the
 * connector picks SQLite or Postgres by DSN shape and has no third dialector)
 * and any Redis "integration" beyond a URL.
 *
 * The storage tab carries setup steps rather than just variable names because
 * the failure people actually hit is procedural: uploads succeed and images
 * never render, because the object URL points at an S3 API endpoint that only
 * answers signed requests. That is called out on the R2 and S3 panels — it is
 * the single most expensive hour a new user loses.
 */

/* ── Storage providers ──────────────────────────────────────────────── */

interface StorageProvider {
  key: string
  name: string
  tagline: string
  env: string
  steps: string[]
  /** CORS or bucket policy the browser upload needs. */
  policy?: { title: string; code: string }
  warning?: string
}

const STORAGE_PROVIDERS: StorageProvider[] = [
  {
    key: 'minio',
    name: 'MinIO',
    tagline: 'The default. Already running after docker compose up.',
    env: `STORAGE_DRIVER=minio

MINIO_ENDPOINT=http://localhost:9002
MINIO_ACCESS_KEY=minioadmin
MINIO_SECRET_KEY=minioadmin
MINIO_BUCKET=myapp-uploads
MINIO_REGION=us-east-1
MINIO_USE_SSL=false`,
    steps: [
      'Nothing to sign up for — docker compose up -d starts MinIO on port 9002.',
      'The bucket is created on first boot; the console is at localhost:9003 (minioadmin / minioadmin).',
      'Uploads work immediately. MinIO serves objects from the same host it takes API calls on, so no public-URL setting is needed.',
    ],
  },
  {
    key: 'r2',
    name: 'Cloudflare R2',
    tagline: 'Zero egress fees. Also the one with the famous gotcha.',
    env: `STORAGE_DRIVER=r2

R2_ENDPOINT=https://<account-id>.r2.cloudflarestorage.com
R2_ACCESS_KEY=...
R2_SECRET_KEY=...
R2_BUCKET=myapp-uploads
R2_REGION=auto

# REQUIRED for images to display — see the warning
R2_PUBLIC_URL=https://pub-<hash>.r2.dev`,
    steps: [
      'Cloudflare dashboard → R2 → Create bucket. Location "Automatic" is fine.',
      'R2 → API → Manage R2 API tokens → Create API token, with Object Read & Write on that bucket.',
      'Copy the Access Key ID, Secret Access Key and the S3 endpoint it shows you — that endpoint is R2_ENDPOINT.',
      'Bucket → Settings → Public Development URL → Enable. Copy the pub-<hash>.r2.dev origin into R2_PUBLIC_URL. (A custom domain works too, and is what you want in production.)',
      'Bucket → Settings → CORS policy → add the rule opposite, with your real origins.',
    ],
    policy: {
      title: 'R2 CORS policy',
      code: `[
  {
    "AllowedOrigins": [
      "http://localhost:3000",
      "http://localhost:3001",
      "https://app.yourdomain.com"
    ],
    "AllowedMethods": ["GET", "PUT"],
    "AllowedHeaders": ["*"],
    "ExposeHeaders": ["ETag"],
    "MaxAgeSeconds": 3000
  }
]`,
    },
    warning:
      'If uploads succeed but images never render, this is why: the R2 endpoint is the S3 API host and only answers SigV4-signed requests, so an <img> pointed at it gets a 401. It looks like CORS and is not. Set R2_PUBLIC_URL to the bucket’s public origin and stored URLs switch to it.',
  },
  {
    key: 's3',
    name: 'AWS S3',
    tagline: 'The one everyone already has an account for.',
    env: `STORAGE_DRIVER=s3

S3_BUCKET=myapp-uploads
S3_REGION=us-east-1
S3_ACCESS_KEY=...
S3_SECRET_KEY=...

# Leave S3_ENDPOINT empty — the SDK picks the regional host.
S3_ENDPOINT=

# Optional: a CloudFront domain, or the bucket host if reads are public
S3_PUBLIC_URL=https://cdn.yourdomain.com`,
    steps: [
      'S3 → Create bucket, pick a region. Leave "Block all public access" on if you will serve through CloudFront or presigned reads.',
      'IAM → Users → Create user, attach a policy scoped to that bucket, and create an access key.',
      'Bucket → Permissions → CORS → paste the rule opposite with your real origins.',
      'If you want plain public reads, turn off "Block all public access" and add the bucket policy opposite, then set S3_PUBLIC_URL.',
      'On EC2/ECS/Lambda you can skip the keys entirely — the config falls back to AWS_ACCESS_KEY_ID / AWS_SECRET_ACCESS_KEY / AWS_REGION, so an IAM role is enough.',
    ],
    policy: {
      title: 'S3 CORS + public-read policy',
      code: `// Permissions → CORS
[
  {
    "AllowedHeaders": ["*"],
    "AllowedMethods": ["GET", "PUT", "HEAD"],
    "AllowedOrigins": ["https://app.yourdomain.com"],
    "ExposeHeaders": ["ETag"],
    "MaxAgeSeconds": 3000
  }
]

// Permissions → Bucket policy (only if reads are public)
{
  "Version": "2012-10-17",
  "Statement": [{
    "Sid": "PublicReadGetObject",
    "Effect": "Allow",
    "Principal": "*",
    "Action": "s3:GetObject",
    "Resource": "arn:aws:s3:::myapp-uploads/*"
  }]
}`,
    },
    warning:
      'Same trap as R2, one step earlier: with "Block all public access" left on and no CDN, objects upload fine and every <img> returns 403. Either front the bucket with CloudFront and set S3_PUBLIC_URL, or make reads public with the bucket policy opposite.',
  },
  {
    key: 'b2',
    name: 'Backblaze B2',
    tagline: 'The cheapest per gigabyte of the four.',
    env: `STORAGE_DRIVER=b2

B2_ENDPOINT=https://s3.us-west-004.backblazeb2.com
B2_ACCESS_KEY=...
B2_SECRET_KEY=...
B2_BUCKET=myapp-uploads
B2_REGION=us-west-004

# Set if the bucket is public or behind a CDN
B2_PUBLIC_URL=`,
    steps: [
      'Backblaze → B2 Cloud Storage → Create a Bucket. Choose Public if you want browsers to read objects directly.',
      'App Keys → Add a New Application Key, scoped to that bucket. Save the keyID and applicationKey.',
      'Copy the S3-compatible endpoint shown on the bucket — it encodes the region, and B2_REGION must match it.',
      'Bucket → CORS Rules → allow GET and PUT from your origins.',
      'For a public bucket, set B2_PUBLIC_URL to the friendly URL Backblaze shows (or your CDN domain).',
    ],
  },
]

/* ── Redis and database code samples ────────────────────────────────── */

const REDIS_CODE = `// Cache — the service is on every handler that needs it.
// Nil-safe: with REDIS_URL unset, svc.Cache is nil and callers skip it.
var products []models.Product
if err := h.Cache.Get(ctx, "products:featured", &products); err != nil {
    h.DB.Where("featured = ?", true).Find(&products)
    h.Cache.Set(ctx, "products:featured", products, cache.DefaultTTL)
}

// Cache a whole route instead, keyed by URL. Adds X-Cache: HIT/MISS.
products.GET("", middleware.CacheResponse(svc.Cache, 5*time.Minute),
    productHandler.List)

// Invalidate by prefix after a write
h.Cache.DeletePattern(ctx, "products:*")

// Background job — enqueued here, run by the worker pool
h.Jobs.EnqueueEmail(ctx, jobs.EmailPayload{
    To:       user.Email,
    Subject:  "Welcome",
    Template: "welcome",
})

// A recurring job. Registered once; the admin lists it at /system/cron.
cron.Register("0 3 * * *", "reports:nightly", asynq.MaxRetry(2))`

const DB_CODE = `// One URL. The driver is chosen by the DSN's shape.
//   DATABASE_URL=sqlite:./app.db
//   DATABASE_URL=sqlite::memory:
//   DATABASE_URL=postgres://user:pass@host:5432/db?sslmode=require

// Models are ordinary GORM structs — generated, then yours.
type Product struct {
    ID         string    \`gorm:"primaryKey" json:"id"\`
    Name       string    \`gorm:"not null" json:"name"\`
    CategoryID string    \`json:"category_id"\`
    Category   *Category \`json:"category,omitempty"\`
    Price      float64   \`json:"price"\`
    CreatedAt  time.Time \`json:"created_at"\`
}

// grit migrate runs AutoMigrate across every registered model,
// and reports what it created or altered rather than doing it silently:
//   + created *models.Product
//   ~ altered *models.Invoice (+2 columns)`

/* ── Topics ─────────────────────────────────────────────────────────── */

interface Provider {
  name: string
  detail: string
  env: string[]
  note?: string
  recommended?: boolean
}

interface Topic {
  key: string
  label: string
  icon: React.ComponentType<{ className?: string }>
  headline: string
  blurb: string
  selector: { env: string; values: string }
  uses: string[]
  footnote: string
  /** Database uses simple provider cards; storage uses nested tabs. */
  providers?: Provider[]
  code?: string
}

const TOPICS: Topic[] = [
  {
    key: 'database',
    label: 'Database',
    icon: Database,
    headline: 'Postgres in production, SQLite when you just want to run it',
    blurb:
      'One DATABASE_URL. The driver is chosen by the shape of the DSN, so moving from a local file to a managed Postgres is an env change and nothing else. GORM is the ORM either way, so your models and migrations do not care.',
    selector: { env: 'DATABASE_URL', values: 'sqlite:… or a postgres:// URL' },
    providers: [
      {
        name: 'PostgreSQL',
        detail: 'The default, and what docker compose brings up locally.',
        env: ['DATABASE_URL=postgres://user:pass@host:5432/db?sslmode=require'],
        note: 'Any Postgres-compatible host works — it is a standard DSN, not a per-vendor integration.',
        recommended: true,
      },
      {
        name: 'SQLite',
        detail: 'A file, or :memory: for tests. Pure-Go driver, so no CGO.',
        env: ['DATABASE_URL=sqlite:./app.db', 'DATABASE_URL=sqlite::memory:'],
        note: 'Ideal for the first five minutes and for the generated Go test suite.',
      },
    ],
    code: DB_CODE,
    uses: [
      'Connection pool tuned on startup: 100 open, 10 idle, 30-minute max lifetime',
      'AutoMigrate wired to every generated model, reporting what changed',
      'GORM Studio at /studio to browse and edit rows',
      'Sessions, roles, permissions and the audit log are all ordinary tables',
    ],
    footnote:
      'MySQL is not supported. The connector picks SQLite or Postgres by DSN prefix — there is no third dialector to fall back to.',
  },
  {
    key: 'storage',
    label: 'File storage',
    icon: HardDrive,
    headline: 'S3-compatible, so the bill is yours to choose',
    blurb:
      'Uploads go browser-to-storage through a presigned URL — files never travel through your API. Switching provider is one env var; nothing in your handlers or components changes. Pick a provider below for the exact setup steps.',
    selector: { env: 'STORAGE_DRIVER', values: 'minio · s3 · r2 · b2' },
    uses: [
      'Presigned direct upload, then a completion call that records the row',
      'Per-field accept lists enforced on the server, not just in the browser',
      'Thumbnails and image processing queued as background jobs',
      'Orphaned uploads swept nightly by a scheduled job',
    ],
    footnote:
      'All four speak the S3 API, which is why one storage client covers them. A provider that does not is not on this list.',
  },
  {
    key: 'redis',
    label: 'Redis',
    icon: Zap,
    headline: 'One URL, five jobs',
    blurb:
      'Redis is not a checkbox in Grit — four subsystems depend on it and a fifth reports on it. All of them read the same REDIS_URL, and all of them keep working without it, minus the feature.',
    selector: { env: 'REDIS_URL', values: 'redis://… or rediss://… for TLS' },
    providers: [
      {
        name: 'Local Redis',
        detail: 'docker compose runs redis:7-alpine on host port 6380.',
        env: ['REDIS_URL=redis://localhost:6380'],
        note: 'Port 6380 rather than 6379 on purpose — it avoids clashing with a Redis you already have installed.',
        recommended: true,
      },
      {
        name: 'Any managed Redis',
        detail: 'Upstash, Redis Cloud, ElastiCache, Railway — a URL is a URL.',
        env: ['REDIS_URL=rediss://default:<password>@<host>:6379'],
        note: 'Use the rediss:// scheme for TLS. There is no provider-specific code to configure.',
      },
    ],
    code: REDIS_CODE,
    uses: [
      'Response cache — GET responses keyed by URL, 5-minute default TTL, X-Cache: HIT/MISS',
      'Idempotency — replays an Idempotency-Key request instead of double-charging, 24-hour window',
      'Job queue — asynq, 10 workers, critical/default/low priorities, exponential backoff',
      'Cron scheduler — token cleanup hourly, orphan uploads nightly, backups every 30 minutes',
      'The admin job inspector and the Redis health check',
    ],
    footnote:
      'Leave REDIS_URL unset and the API still boots: caching, idempotency, jobs and cron each log a warning and switch off. Only the admin job screens hard-fail, with a 503.',
  },
]

export function InfraShowcase() {
  const [active, setActive] = useState(TOPICS[0].key)
  const [provider, setProvider] = useState(STORAGE_PROVIDERS[0].key)
  const topic = TOPICS.find((t) => t.key === active) ?? TOPICS[0]
  const sp = STORAGE_PROVIDERS.find((p) => p.key === provider) ?? STORAGE_PROVIDERS[0]

  return (
    <div>
      <div role="tablist" aria-label="Infrastructure" className="flex flex-wrap gap-2 mb-8">
        {TOPICS.map((t) => {
          const selected = t.key === active
          return (
            <button
              key={t.key}
              type="button"
              role="tab"
              aria-selected={selected}
              onClick={() => setActive(t.key)}
              className={`inline-flex items-center gap-2 rounded-full border px-4 py-2 text-sm font-medium transition-colors ${
                selected
                  ? 'border-primary/40 bg-primary/10 text-foreground'
                  : 'border-border/60 text-muted-foreground hover:text-foreground hover:border-border'
              }`}
            >
              <t.icon className="h-4 w-4" />
              {t.label}
            </button>
          )
        })}
      </div>

      <div className="mb-8 max-w-3xl">
        <h3 className="text-xl md:text-2xl font-semibold mb-3 leading-snug">{topic.headline}</h3>
        <p className="text-sm md:text-base text-muted-foreground leading-relaxed mb-4">{topic.blurb}</p>
        <div className="inline-flex flex-wrap items-center gap-2 rounded-lg border border-border/60 bg-card/40 px-3 py-2">
          <span className="text-[10.5px] font-mono uppercase tracking-wider text-muted-foreground/70">
            Switch with
          </span>
          <code className="text-[12px] font-mono text-foreground/90">{topic.selector.env}</code>
          <span className="text-[11.5px] text-muted-foreground">{topic.selector.values}</span>
        </div>
      </div>

      {/* ── Storage: a nested tab per provider, each with real steps ── */}
      {topic.key === 'storage' ? (
        <>
          <div role="tablist" aria-label="Storage provider" className="flex flex-wrap gap-2 mb-6">
            {STORAGE_PROVIDERS.map((p) => {
              const selected = p.key === provider
              return (
                <button
                  key={p.key}
                  type="button"
                  role="tab"
                  aria-selected={selected}
                  onClick={() => setProvider(p.key)}
                  className={`rounded-lg border px-3.5 py-1.5 text-[13px] font-medium transition-colors ${
                    selected
                      ? 'border-primary/40 bg-primary/10 text-foreground'
                      : 'border-border/50 text-muted-foreground hover:text-foreground'
                  }`}
                >
                  {p.name}
                </button>
              )
            })}
          </div>

          <p className="text-sm text-muted-foreground mb-6">{sp.tagline}</p>

          <div className="grid lg:grid-cols-2 gap-6 items-start">
            <div>
              <div className="rounded-xl border border-border/60 bg-card/40 overflow-hidden mb-4">
                <div className="px-4 py-2.5 border-b border-border/50 text-[10.5px] font-mono uppercase tracking-wider text-muted-foreground/70">
                  .env
                </div>
                <div className="overflow-x-auto p-4">
                  <pre className="text-[11.5px] font-mono text-foreground/85 leading-relaxed whitespace-pre">
                    {sp.env}
                  </pre>
                </div>
              </div>

              {sp.warning && (
                <div className="rounded-xl border border-amber-500/30 bg-amber-500/[0.07] p-4">
                  <div className="flex items-center gap-2 mb-2">
                    <AlertTriangle className="h-4 w-4 text-amber-500 shrink-0" />
                    <span className="text-[12px] font-semibold text-amber-200/90">
                      Uploads work, images do not show?
                    </span>
                  </div>
                  <p className="text-[12px] text-muted-foreground leading-relaxed">{sp.warning}</p>
                </div>
              )}
            </div>

            <div>
              <div className="rounded-xl border border-border/60 bg-card/40 p-4 mb-4">
                <div className="text-[10.5px] font-mono uppercase tracking-wider text-muted-foreground/70 mb-3">
                  Where the values come from
                </div>
                <ol className="space-y-2.5">
                  {sp.steps.map((step, i) => (
                    <li key={step} className="flex gap-2.5">
                      <span className="mt-0.5 flex h-4.5 w-4.5 min-w-[1.125rem] shrink-0 items-center justify-center rounded-full bg-primary/10 text-[10px] font-mono font-medium text-primary">
                        {i + 1}
                      </span>
                      <span className="text-[12.5px] text-foreground/80 leading-relaxed">{step}</span>
                    </li>
                  ))}
                </ol>
              </div>

              {sp.policy && (
                <div className="rounded-xl border border-border/60 bg-card/40 overflow-hidden">
                  <div className="px-4 py-2.5 border-b border-border/50 text-[10.5px] font-mono uppercase tracking-wider text-muted-foreground/70">
                    {sp.policy.title}
                  </div>
                  <div className="overflow-x-auto p-4 max-h-72">
                    <pre className="text-[11px] font-mono text-foreground/80 leading-relaxed whitespace-pre">
                      {sp.policy.code}
                    </pre>
                  </div>
                </div>
              )}
            </div>
          </div>

          <div className="grid sm:grid-cols-2 gap-x-8 gap-y-2.5 mt-8 pt-6 border-t border-border/40">
            {topic.uses.map((u) => (
              <div key={u} className="flex gap-2 text-[12.5px] text-foreground/80 leading-relaxed">
                <Check className="h-3.5 w-3.5 shrink-0 mt-0.5 text-primary" strokeWidth={2.5} />
                {u}
              </div>
            ))}
          </div>
          <p className="text-[11.5px] text-muted-foreground/70 mt-5">{topic.footnote}</p>
        </>
      ) : (
        /* ── Database and Redis: cards plus a code sample ── */
        <div className="grid lg:grid-cols-2 gap-8 items-start">
          <div>
            <div className="grid sm:grid-cols-2 gap-4 mb-6">
              {topic.providers?.map((p) => (
                <div key={p.name} className="rounded-xl border border-border/60 bg-card/40 p-4 flex flex-col">
                  <div className="flex items-center gap-2 mb-1.5">
                    <h4 className="text-[15px] font-semibold">{p.name}</h4>
                    {p.recommended && (
                      <span className="rounded-full bg-primary/10 border border-primary/20 px-2 py-0.5 text-[10px] font-mono uppercase tracking-wider text-primary">
                        default
                      </span>
                    )}
                  </div>
                  <p className="text-[12.5px] text-muted-foreground leading-relaxed mb-3">{p.detail}</p>
                  <div className="overflow-x-auto rounded-lg bg-background/60 border border-border/40 p-2.5 mb-3">
                    <pre className="text-[11px] font-mono text-foreground/80 leading-relaxed whitespace-pre">
                      {p.env.join('\n')}
                    </pre>
                  </div>
                  {p.note && (
                    <p className="text-[11.5px] text-muted-foreground/80 leading-relaxed mt-auto">{p.note}</p>
                  )}
                </div>
              ))}
            </div>

            <div className="rounded-xl border border-border/60 bg-card/40 p-4">
              <div className="text-[10.5px] font-mono uppercase tracking-wider text-muted-foreground/70 mb-3">
                What Grit does with it
              </div>
              <ul className="space-y-2.5">
                {topic.uses.map((u) => (
                  <li key={u} className="flex gap-2 text-[12.5px] text-foreground/80 leading-relaxed">
                    <Check className="h-3.5 w-3.5 shrink-0 mt-0.5 text-primary" strokeWidth={2.5} />
                    {u}
                  </li>
                ))}
              </ul>
            </div>
            <p className="text-[11.5px] text-muted-foreground/70 mt-4 leading-relaxed">{topic.footnote}</p>
          </div>

          <div className="rounded-xl overflow-hidden border border-border bg-card/40">
            <div className="flex items-center gap-2 px-4 py-2.5 bg-card/60 border-b border-border/60">
              <span className="text-[11.5px] font-mono text-muted-foreground">
                {topic.key === 'redis' ? 'using it in your handlers' : 'internal/models/product.go'}
              </span>
            </div>
            <CodeBlock
              key={topic.key}
              code={topic.code ?? ''}
              language="go"
              className="!border-0 !rounded-none !shadow-none !bg-transparent dark:!bg-transparent !m-0"
            />
          </div>
        </div>
      )}
    </div>
  )
}
