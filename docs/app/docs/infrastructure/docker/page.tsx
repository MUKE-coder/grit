import Link from 'next/link'
import { ArrowLeft, ArrowRight } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'
import { CodeBlock } from '@/components/code-block'
import { getDocMetadata } from '@/config/docs-metadata'

export const metadata = getDocMetadata('/docs/infrastructure/docker')

export default function DockerSetupPage() {
  return (
    <div className="min-h-screen bg-background">
      <SiteHeader />
      <DocsSidebar />

      <main className="lg:pl-64">
        <div className="container max-w-screen-xl py-10 px-6">
          <div className="max-w-3xl">
            {/* Header */}
            <div className="mb-10">
              <span className="tag-mono text-primary/80 mb-3 block">Infrastructure</span>
              <h1 className="text-4xl font-bold tracking-tight mb-4">
                Docker Setup
              </h1>
              <p className="text-lg text-muted-foreground leading-relaxed">
                Grit uses Docker Compose to run all infrastructure services locally.
                One command gives you PostgreSQL, Redis, MinIO, and Mailhog &mdash; ready to go.
              </p>
            </div>

            {/* Overview */}
            <div className="prose-grit">
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Overview
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Every Grit project ships with two Docker Compose files:
                </p>
                <ul className="space-y-2.5 mb-6">
                  {[
                    { file: 'docker-compose.yml', desc: 'Development infrastructure (databases, caching, storage, email)' },
                    { file: 'docker-compose.prod.yml', desc: 'Production deployment (includes API, web, admin containers)' },
                  ].map((item) => (
                    <li key={item.file} className="flex items-start gap-2.5 text-[14px] text-muted-foreground">
                      <span className="text-primary mt-1">&#10003;</span>
                      <span>
                        <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">{item.file}</code>{' '}
                        &mdash; {item.desc}
                      </span>
                    </li>
                  ))}
                </ul>
                <p className="text-muted-foreground leading-relaxed">
                  You do <strong>not</strong> need Docker to run the Go API or Next.js apps directly &mdash; Docker
                  is only required for the infrastructure services (PostgreSQL, Redis, etc.). The API and frontends
                  run natively during development for faster iteration.
                </p>
              </div>

              {/* docker-compose.yml */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Development Compose File
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  The development <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">docker-compose.yml</code> spins
                  up four services. Start them all with a single command:
                </p>
                <div className="rounded-xl border border-border/40 bg-card/80 overflow-hidden glow-purple-sm mb-6">
                  <div className="flex items-center gap-2 px-4 py-2.5 border-b border-border/30 bg-accent/30">
                    <div className="flex items-center gap-1.5">
                      <div className="h-2.5 w-2.5 rounded-full bg-[#ff5f57]" />
                      <div className="h-2.5 w-2.5 rounded-full bg-[#febc2e]" />
                      <div className="h-2.5 w-2.5 rounded-full bg-[#28c840]" />
                    </div>
                    <span className="ml-2 text-[11px] font-mono text-muted-foreground/40">terminal</span>
                  </div>
                  <div className="p-5 font-mono text-sm">
                    <span className="text-primary/50 select-none">$ </span>
                    <span className="text-foreground/80">docker compose up -d</span>
                  </div>
                </div>

                <CodeBlock language="yaml" filename="docker-compose.yml" code={`services:
  postgres:
    image: postgres:16-alpine
    container_name: myapp-postgres
    restart: unless-stopped
    ports:
      - "5432:5432"
    environment:
      POSTGRES_USER: grit
      POSTGRES_PASSWORD: grit
      POSTGRES_DB: myapp
    volumes:
      - postgres-data:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U grit"]
      interval: 5s
      timeout: 5s
      retries: 5

  redis:
    image: redis:7-alpine
    container_name: myapp-redis
    restart: unless-stopped
    ports:
      - "6379:6379"
    volumes:
      - redis-data:/data
    healthcheck:
      test: ["CMD", "redis-cli", "ping"]
      interval: 5s
      timeout: 5s
      retries: 5

  minio:
    image: minio/minio
    container_name: myapp-minio
    restart: unless-stopped
    ports:
      - "9000:9000"
      - "9001:9001"
    environment:
      MINIO_ROOT_USER: minioadmin
      MINIO_ROOT_PASSWORD: minioadmin
    volumes:
      - minio-data:/data
    command: server /data --console-address ":9001"

  mailhog:
    image: mailhog/mailhog
    container_name: myapp-mailhog
    restart: unless-stopped
    ports:
      - "1025:1025"
      - "8025:8025"

volumes:
  postgres-data:
  redis-data:
  minio-data:`} />
              </div>

              {/* Service details table */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Service Details
                </h2>
                <div className="rounded-lg border border-border/30 bg-card/30 overflow-hidden">
                  <table className="w-full text-sm">
                    <thead>
                      <tr className="border-b border-border/30 bg-accent/20">
                        <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Service</th>
                        <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Port(s)</th>
                        <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Credentials</th>
                        <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Purpose</th>
                      </tr>
                    </thead>
                    <tbody className="text-muted-foreground">
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-medium text-foreground/90">PostgreSQL</td>
                        <td className="px-4 py-2.5 font-mono text-xs">5432</td>
                        <td className="px-4 py-2.5 font-mono text-xs">grit / grit</td>
                        <td className="px-4 py-2.5">Primary database</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-medium text-foreground/90">Redis</td>
                        <td className="px-4 py-2.5 font-mono text-xs">6379</td>
                        <td className="px-4 py-2.5 font-mono text-xs">No auth</td>
                        <td className="px-4 py-2.5">Cache, sessions, job queues</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-medium text-foreground/90">MinIO</td>
                        <td className="px-4 py-2.5 font-mono text-xs">9000 / 9001</td>
                        <td className="px-4 py-2.5 font-mono text-xs">minioadmin / minioadmin</td>
                        <td className="px-4 py-2.5">S3-compatible file storage</td>
                      </tr>
                      <tr>
                        <td className="px-4 py-2.5 font-medium text-foreground/90">Mailhog</td>
                        <td className="px-4 py-2.5 font-mono text-xs">1025 / 8025</td>
                        <td className="px-4 py-2.5 font-mono text-xs">No auth</td>
                        <td className="px-4 py-2.5">Email testing (SMTP + Web UI)</td>
                      </tr>
                    </tbody>
                  </table>
                </div>
              </div>

              {/* Accessing services */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Accessing Services
                </h2>
                <div className="space-y-4">
                  {[
                    {
                      title: 'MinIO Console',
                      url: 'http://localhost:9001',
                      desc: 'Web-based file browser for your S3-compatible storage. Create buckets, upload files, manage access policies. Login with minioadmin / minioadmin.',
                    },
                    {
                      title: 'Mailhog UI',
                      url: 'http://localhost:8025',
                      desc: 'Catches all outgoing emails from your application. View HTML emails, check headers, and test email flows without sending real emails.',
                    },
                    {
                      title: 'PostgreSQL',
                      url: 'localhost:5432',
                      desc: 'Connect using any database client (pgAdmin, TablePlus, DBeaver). Connection string: postgres://grit:grit@localhost:5432/myapp?sslmode=disable',
                    },
                    {
                      title: 'Redis',
                      url: 'localhost:6379',
                      desc: 'Connect with redis-cli or any Redis GUI client (RedisInsight, Medis). No authentication required in development.',
                    },
                  ].map((service) => (
                    <div key={service.title} className="p-4 rounded-lg border border-border/30 bg-card/30">
                      <div className="flex items-center justify-between mb-1.5">
                        <h3 className="text-sm font-semibold">{service.title}</h3>
                        <code className="text-xs font-mono text-primary/70 bg-primary/5 px-2 py-0.5 rounded-md border border-primary/10">
                          {service.url}
                        </code>
                      </div>
                      <p className="text-xs text-muted-foreground/70 leading-relaxed">{service.desc}</p>
                    </div>
                  ))}
                </div>
              </div>

              {/* Production compose */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Production Compose File
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  The <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">docker-compose.prod.yml</code> builds
                  and runs your entire application stack including the Go API, Next.js web app, and admin panel
                  alongside PostgreSQL and Redis.
                </p>

                <CodeBlock language="yaml" filename="docker-compose.prod.yml" code={`services:
  api:
    build:
      context: ./apps/api
      dockerfile: Dockerfile
    container_name: myapp-api
    restart: unless-stopped
    ports:
      - "8080:8080"
    environment:
      APP_ENV: production
      DATABASE_URL: postgres://grit:grit@postgres:5432/myapp?sslmode=disable
      REDIS_URL: redis://redis:6379
      JWT_SECRET: \${JWT_SECRET}
    depends_on:
      postgres:
        condition: service_healthy
      redis:
        condition: service_healthy

  web:
    build:
      context: .
      dockerfile: apps/web/Dockerfile
    container_name: myapp-web
    restart: unless-stopped
    ports:
      - "3000:3000"
    environment:
      NEXT_PUBLIC_API_URL: http://api:8080

  admin:
    build:
      context: .
      dockerfile: apps/admin/Dockerfile
    container_name: myapp-admin
    restart: unless-stopped
    ports:
      - "3001:3000"
    environment:
      NEXT_PUBLIC_API_URL: http://api:8080

  postgres:
    image: postgres:16-alpine
    container_name: myapp-postgres
    restart: unless-stopped
    environment:
      POSTGRES_USER: grit
      POSTGRES_PASSWORD: grit
      POSTGRES_DB: myapp
    volumes:
      - postgres-data:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U grit"]
      interval: 5s
      timeout: 5s
      retries: 5

  redis:
    image: redis:7-alpine
    container_name: myapp-redis
    restart: unless-stopped
    volumes:
      - redis-data:/data
    healthcheck:
      test: ["CMD", "redis-cli", "ping"]
      interval: 5s
      timeout: 5s
      retries: 5

volumes:
  postgres-data:
  redis-data:`} />

                <p className="text-muted-foreground leading-relaxed mb-4">
                  Deploy to production with:
                </p>
                <div className="rounded-xl border border-border/40 bg-card/80 overflow-hidden glow-purple-sm">
                  <div className="flex items-center gap-2 px-4 py-2.5 border-b border-border/30 bg-accent/30">
                    <div className="flex items-center gap-1.5">
                      <div className="h-2.5 w-2.5 rounded-full bg-[#ff5f57]" />
                      <div className="h-2.5 w-2.5 rounded-full bg-[#febc2e]" />
                      <div className="h-2.5 w-2.5 rounded-full bg-[#28c840]" />
                    </div>
                    <span className="ml-2 text-[11px] font-mono text-muted-foreground/40">terminal</span>
                  </div>
                  <div className="p-5 font-mono text-sm">
                    <span className="text-primary/50 select-none">$ </span>
                    <span className="text-foreground/80">docker compose -f docker-compose.prod.yml up -d --build</span>
                  </div>
                </div>
              </div>

              {/* Dockerfiles */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Dockerfiles
                </h2>

                <h3 className="text-lg font-semibold tracking-tight mb-3 mt-6">
                  Go API (Multi-Stage Build)
                </h3>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  The API Dockerfile uses a multi-stage build. The first stage compiles the Go binary
                  with all dependencies, and the second stage copies only the binary into a minimal
                  Alpine image. The final image is typically under 20MB.
                </p>
                <CodeBlock filename="apps/api/Dockerfile" code={`# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build binary
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server ./cmd/server

# Run stage
FROM alpine:3.19

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

COPY --from=builder /app/server .

EXPOSE 8080

CMD ["./server"]`} />

                <h3 className="text-lg font-semibold tracking-tight mb-3 mt-6">
                  Next.js (Standalone Build)
                </h3>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  The Next.js Dockerfile also uses a multi-stage build. It installs dependencies,
                  builds the app with standalone output, and runs the production server as a non-root user.
                  Both the web and admin apps share this same Dockerfile pattern.
                </p>
                <CodeBlock filename="apps/web/Dockerfile" code={`# Build stage
FROM node:20-alpine AS base

RUN corepack enable

# Install dependencies
FROM base AS deps
WORKDIR /app

COPY package.json pnpm-lock.yaml pnpm-workspace.yaml ./
COPY apps/web/package.json ./apps/web/
COPY packages/shared/package.json ./packages/shared/

RUN pnpm install --frozen-lockfile

# Build
FROM base AS builder
WORKDIR /app

COPY --from=deps /app/node_modules ./node_modules
COPY --from=deps /app/apps/web/node_modules ./apps/web/node_modules
COPY --from=deps /app/packages/shared/node_modules ./packages/shared/node_modules
COPY . .

RUN pnpm --filter web build

# Run
FROM base AS runner
WORKDIR /app

ENV NODE_ENV=production

RUN addgroup --system --gid 1001 nodejs
RUN adduser --system --uid 1001 nextjs

COPY --from=builder /app/apps/web/.next/standalone ./
COPY --from=builder /app/apps/web/.next/static ./apps/web/.next/static
COPY --from=builder /app/apps/web/public ./apps/web/public

USER nextjs

EXPOSE 3000

ENV PORT=3000
ENV HOSTNAME="0.0.0.0"

CMD ["node", "apps/web/server.js"]`} />
              </div>

              {/* Common Commands */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Common Commands
                </h2>
                <div className="space-y-3">
                  {[
                    { cmd: 'docker compose up -d', desc: 'Start all development services in the background' },
                    { cmd: 'docker compose down', desc: 'Stop and remove all containers' },
                    { cmd: 'docker compose down -v', desc: 'Stop all containers and delete volumes (resets all data)' },
                    { cmd: 'docker compose logs -f postgres', desc: 'Follow logs for a specific service' },
                    { cmd: 'docker compose logs -f', desc: 'Follow logs for all services' },
                    { cmd: 'docker compose ps', desc: 'Show running containers and their status' },
                    { cmd: 'docker compose restart redis', desc: 'Restart a specific service' },
                    { cmd: 'docker compose exec postgres psql -U grit myapp', desc: 'Open a psql shell inside the PostgreSQL container' },
                    { cmd: 'docker compose exec redis redis-cli', desc: 'Open a Redis CLI session inside the container' },
                  ].map((item) => (
                    <div key={item.cmd} className="rounded-xl border border-border/40 bg-card/80 overflow-hidden">
                      <div className="flex items-center gap-2 px-4 py-2.5 border-b border-border/30 bg-accent/30">
                        <div className="flex items-center gap-1.5">
                          <div className="h-2.5 w-2.5 rounded-full bg-[#ff5f57]" />
                          <div className="h-2.5 w-2.5 rounded-full bg-[#febc2e]" />
                          <div className="h-2.5 w-2.5 rounded-full bg-[#28c840]" />
                        </div>
                        <span className="ml-2 text-[11px] font-mono text-muted-foreground/40">terminal</span>
                      </div>
                      <div className="px-5 py-3 font-mono text-sm">
                        <div>
                          <span className="text-primary/50 select-none">$ </span>
                          <span className="text-foreground/80">{item.cmd}</span>
                        </div>
                        <p className="text-xs text-muted-foreground/50 mt-1">{item.desc}</p>
                      </div>
                    </div>
                  ))}
                </div>
              </div>

              {/* Volumes */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Data Persistence
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Docker named volumes keep your data safe across container restarts:
                </p>
                <div className="rounded-lg border border-border/30 bg-card/30 overflow-hidden">
                  <table className="w-full text-sm">
                    <thead>
                      <tr className="border-b border-border/30 bg-accent/20">
                        <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Volume</th>
                        <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Mounted To</th>
                        <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Contains</th>
                      </tr>
                    </thead>
                    <tbody className="text-muted-foreground">
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-mono text-xs">postgres-data</td>
                        <td className="px-4 py-2.5 font-mono text-xs">/var/lib/postgresql/data</td>
                        <td className="px-4 py-2.5">All database tables and data</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-mono text-xs">redis-data</td>
                        <td className="px-4 py-2.5 font-mono text-xs">/data</td>
                        <td className="px-4 py-2.5">Cached data, sessions, job queue state</td>
                      </tr>
                      <tr>
                        <td className="px-4 py-2.5 font-mono text-xs">minio-data</td>
                        <td className="px-4 py-2.5 font-mono text-xs">/data</td>
                        <td className="px-4 py-2.5">Uploaded files and bucket data</td>
                      </tr>
                    </tbody>
                  </table>
                </div>
                <p className="text-sm text-muted-foreground/60 mt-3">
                  Use <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">docker compose down -v</code> to
                  delete all volumes and reset to a clean state. This is useful when you want to start fresh
                  or if your database schema has diverged.
                </p>
              </div>

              {/* Troubleshooting */}
              <div className="mb-10">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Troubleshooting
                </h2>
                <div className="space-y-4">
                  {[
                    {
                      q: 'Port already in use',
                      a: 'Another process is using port 5432, 6379, etc. Stop the conflicting process or change the port mapping in docker-compose.yml. For example, change "5432:5432" to "5433:5432" and update your .env DATABASE_URL accordingly.',
                    },
                    {
                      q: 'Container keeps restarting',
                      a: 'Check the logs with "docker compose logs <service>". Common causes: incorrect credentials, corrupted volume data. Try "docker compose down -v && docker compose up -d" for a clean start.',
                    },
                    {
                      q: 'Cannot connect from API to PostgreSQL',
                      a: 'Ensure the API uses "localhost" (not the container name) when running outside Docker. The connection string in .env should be: postgres://grit:grit@localhost:5432/myapp?sslmode=disable',
                    },
                    {
                      q: 'MinIO bucket not found',
                      a: 'MinIO starts with no buckets. Open the MinIO console at http://localhost:9001, login with minioadmin/minioadmin, and create your bucket. Or set MINIO_DEFAULT_BUCKETS in the compose file.',
                    },
                    {
                      q: 'Docker Compose V1 vs V2',
                      a: 'Grit uses "docker compose" (V2, no hyphen). If you see errors, make sure Docker Desktop is updated. The old "docker-compose" (V1, with hyphen) is deprecated.',
                    },
                  ].map((item) => (
                    <div key={item.q} className="p-4 rounded-lg border border-border/30 bg-card/30">
                      <h3 className="text-sm font-semibold mb-1.5">{item.q}</h3>
                      <p className="text-xs text-muted-foreground/70 leading-relaxed">{item.a}</p>
                    </div>
                  ))}
                </div>
              </div>

              {/* Nav */}
              <div className="flex items-center justify-between pt-6 border-t border-border/30">
                <Button variant="ghost" size="sm" asChild className="text-muted-foreground/60 hover:text-foreground">
                  <Link href="/docs/batteries/ai" className="gap-1.5">
                    <ArrowLeft className="h-3.5 w-3.5" />
                    AI Integration
                  </Link>
                </Button>
                <Button variant="ghost" size="sm" asChild className="text-muted-foreground/60 hover:text-foreground">
                  <Link href="/docs/infrastructure/database" className="gap-1.5">
                    Database &amp; Migrations
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
