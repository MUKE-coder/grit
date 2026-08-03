'use client'

import { useState } from 'react'
import type React from 'react'
import { Database, HardDrive, Zap, Check } from 'lucide-react'

/**
 * Databases, file storage and Redis — what Grit actually supports.
 *
 * Every provider, env var and default here was read out of the scaffold
 * templates, not assumed. Two things this section deliberately does NOT claim:
 * MySQL (there is no MySQL dialector — internal/scaffold/api_files.go picks
 * SQLite or Postgres by DSN shape and nothing else), and any Redis provider
 * "integration" beyond a URL.
 *
 * If you add a driver, add it here with its real env vars. A provider logo on
 * a marketing page that turns out to be unimplemented costs more trust than
 * the logo ever bought.
 */

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
  /** The single switch that picks a provider. */
  selector: { env: string; values: string }
  providers: Provider[]
  /** What the framework does with it, as plain statements. */
  uses: string[]
  footnote: string
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
    uses: [
      'Connection pool tuned on startup: 100 open, 10 idle, 30-minute max lifetime',
      'AutoMigrate wired to every generated model',
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
      'Uploads go browser-to-storage through a presigned URL — files never travel through your API. Switching provider is one env var; nothing in your handlers or components changes.',
    selector: { env: 'STORAGE_DRIVER', values: 'minio · s3 · r2 · b2' },
    providers: [
      {
        name: 'MinIO',
        detail: 'The default. docker compose brings it up on port 9002.',
        env: ['MINIO_ENDPOINT=http://localhost:9002', 'MINIO_ACCESS_KEY', 'MINIO_SECRET_KEY', 'MINIO_BUCKET'],
        note: 'Runs locally with minioadmin/minioadmin so uploads work before you have any cloud account.',
        recommended: true,
      },
      {
        name: 'AWS S3',
        detail: 'Leave the endpoint empty and the SDK finds the regional one.',
        env: ['S3_BUCKET', 'S3_REGION', 'S3_ACCESS_KEY', 'S3_SECRET_KEY'],
        note: 'Falls back to AWS_ACCESS_KEY_ID / AWS_SECRET_ACCESS_KEY / AWS_REGION, so an IAM role on EC2 or ECS needs no keys in .env at all.',
      },
      {
        name: 'Cloudflare R2',
        detail: 'Zero egress fees, region "auto".',
        env: ['R2_ENDPOINT=https://<account>.r2.cloudflarestorage.com', 'R2_BUCKET', 'R2_ACCESS_KEY', 'R2_SECRET_KEY'],
      },
      {
        name: 'Backblaze B2',
        detail: 'The cheapest per gigabyte of the four.',
        env: ['B2_ENDPOINT', 'B2_BUCKET', 'B2_ACCESS_KEY', 'B2_SECRET_KEY'],
      },
    ],
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
  const topic = TOPICS.find((t) => t.key === active) ?? TOPICS[0]

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

      <div className="grid lg:grid-cols-[1fr_20rem] gap-8 items-start">
        <div className="grid sm:grid-cols-2 gap-4">
          {topic.providers.map((p) => (
            <div
              key={p.name}
              className="rounded-xl border border-border/60 bg-card/40 p-4 flex flex-col"
            >
              <div className="flex items-center gap-2 mb-1.5">
                <h4 className="text-[15px] font-semibold">{p.name}</h4>
                {p.recommended && (
                  <span className="rounded-full bg-primary/10 border border-primary/20 px-2 py-0.5 text-[10px] font-mono uppercase tracking-wider text-primary">
                    default
                  </span>
                )}
              </div>
              <p className="text-[12.5px] text-muted-foreground leading-relaxed mb-3">{p.detail}</p>

              {/* Long DSNs scroll inside the card rather than widening the page. */}
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

        <div>
          <div className="rounded-xl border border-border/60 bg-card/40 p-4 mb-4">
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
          <p className="text-[11.5px] text-muted-foreground/70 leading-relaxed">{topic.footnote}</p>
        </div>
      </div>
    </div>
  )
}
