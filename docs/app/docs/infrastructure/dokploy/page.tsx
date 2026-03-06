import Link from 'next/link'
import { ArrowLeft, ArrowRight } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'
import { CodeBlock } from '@/components/code-block'
import { getDocMetadata } from '@/config/docs-metadata'

export const metadata = getDocMetadata('/docs/infrastructure/dokploy')

export default function DokployPage() {
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
                Deploy with Dokploy
              </h1>
              <p className="text-lg text-muted-foreground leading-relaxed">
                Deploy your Grit application using Dokploy &mdash; a self-hosted PaaS that gives you
                the convenience of Vercel or Railway on your own VPS. Web dashboard, auto-SSL, GitHub
                integration, and Docker Compose support &mdash; all for free.
              </p>
            </div>

            <div className="prose-grit">

              {/* Table of Contents */}
              <div className="mb-12 rounded-xl border border-border/40 bg-card/30 p-6">
                <h2 className="text-lg font-semibold tracking-tight mb-3">On This Page</h2>
                <div className="grid gap-1.5 sm:grid-cols-2">
                  {[
                    { label: '1. What is Dokploy?', id: '#what-is-dokploy' },
                    { label: '2. Prerequisites', id: '#prerequisites' },
                    { label: '3. Install Dokploy', id: '#install' },
                    { label: '4. Create a Project', id: '#create-project' },
                    { label: '5. Docker Compose Config', id: '#compose' },
                    { label: '6. Environment Variables', id: '#env-vars' },
                    { label: '7. Configure Domains', id: '#domains' },
                    { label: '8. DNS Configuration', id: '#dns' },
                    { label: '9. Deploy', id: '#deploy' },
                    { label: '10. Database Backups', id: '#backups' },
                    { label: '11. Monitoring', id: '#monitoring' },
                    { label: '12. Tips & Troubleshooting', id: '#troubleshooting' },
                  ].map((item) => (
                    <a key={item.id} href={item.id} className="text-[13px] text-muted-foreground/70 hover:text-primary transition-colors py-0.5">
                      {item.label}
                    </a>
                  ))}
                </div>
              </div>

              {/* Step 1: What is Dokploy? */}
              <div id="what-is-dokploy" className="mb-12 scroll-mt-20">
                <div className="flex items-center gap-3 mb-4">
                  <div className="flex h-8 w-8 items-center justify-center rounded-md bg-primary/10 border border-primary/15 text-sm font-mono font-semibold text-primary shrink-0">
                    1
                  </div>
                  <h2 className="text-2xl font-semibold tracking-tight">
                    What is Dokploy?
                  </h2>
                </div>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  <a href="https://dokploy.com" className="text-primary hover:underline" target="_blank" rel="noopener noreferrer">Dokploy</a> is
                  an open-source, self-hosted Platform-as-a-Service (PaaS). Think of it as your own
                  private Vercel or Heroku running on a VPS you control. It provides:
                </p>
                <div className="rounded-lg border border-border/30 bg-card/30 overflow-hidden mb-6">
                  <table className="w-full text-sm">
                    <thead>
                      <tr className="border-b border-border/30 bg-accent/20">
                        <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Feature</th>
                        <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Description</th>
                      </tr>
                    </thead>
                    <tbody className="text-muted-foreground">
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-medium text-foreground/90">Web Dashboard</td>
                        <td className="px-4 py-2.5 text-xs">Visual UI for managing deployments, logs, domains, and environment variables</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-medium text-foreground/90">Docker Compose</td>
                        <td className="px-4 py-2.5 text-xs">Deploy multi-service apps with a single Compose file &mdash; perfect for Grit monorepos</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-medium text-foreground/90">Auto SSL</td>
                        <td className="px-4 py-2.5 text-xs">Automatic Let&apos;s Encrypt certificates for all your domains via Traefik</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-medium text-foreground/90">GitHub Integration</td>
                        <td className="px-4 py-2.5 text-xs">Auto-deploy on git push with webhook-based continuous deployment</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-medium text-foreground/90">Database Backups</td>
                        <td className="px-4 py-2.5 text-xs">Scheduled backups to S3-compatible storage with one-click restore</td>
                      </tr>
                      <tr>
                        <td className="px-4 py-2.5 font-medium text-foreground/90">Monitoring</td>
                        <td className="px-4 py-2.5 text-xs">Built-in logs viewer, resource usage graphs, and deployment history</td>
                      </tr>
                    </tbody>
                  </table>
                </div>

                <div className="p-4 rounded-lg border border-primary/20 bg-primary/5">
                  <p className="text-sm text-foreground/80 leading-relaxed">
                    <span className="font-semibold text-primary/90">Why Dokploy for Grit?</span>{' '}
                    Grit&apos;s monorepo (Go API + 2 Next.js apps + PostgreSQL + Redis) maps naturally to a
                    Docker Compose project. Dokploy deploys the entire stack from a single Compose file,
                    handles SSL for all services, and gives you a dashboard to manage everything &mdash;
                    no manual Caddy/Nginx config, no SSH for routine deploys.
                  </p>
                </div>
              </div>

              {/* Step 2: Prerequisites */}
              <div id="prerequisites" className="mb-12 scroll-mt-20">
                <div className="flex items-center gap-3 mb-4">
                  <div className="flex h-8 w-8 items-center justify-center rounded-md bg-primary/10 border border-primary/15 text-sm font-mono font-semibold text-primary shrink-0">
                    2
                  </div>
                  <h2 className="text-2xl font-semibold tracking-tight">
                    Prerequisites
                  </h2>
                </div>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Before you begin, make sure you have:
                </p>
                <ul className="space-y-2.5 mb-6">
                  {[
                    'A VPS with at least 2 vCPU / 4GB RAM (Hetzner CX22, DigitalOcean 4GB, or similar)',
                    'Ubuntu 22.04 or Debian 12 on the VPS (fresh install recommended)',
                    'A domain name (e.g., yourdomain.com) with access to DNS settings',
                    'Your Grit project pushed to a GitHub repository',
                    'SSH access to your VPS (root or sudo user)',
                  ].map((item) => (
                    <li key={item} className="flex items-start gap-2.5 text-[14px] text-muted-foreground">
                      <span className="text-primary mt-0.5 shrink-0">&#10003;</span>
                      {item}
                    </li>
                  ))}
                </ul>

                <div className="p-4 rounded-lg border border-primary/20 bg-primary/5">
                  <p className="text-sm text-foreground/80 leading-relaxed">
                    <span className="font-semibold text-primary/90">Server sizing:</span>{' '}
                    A Grit monorepo (Go API + 2 Next.js apps + PostgreSQL + Redis + MinIO) uses
                    around <span className="font-mono text-primary/80">2-3 GB</span> of RAM at idle.
                    We recommend <span className="font-mono text-primary/80">4 GB RAM</span> minimum
                    to leave room for builds and traffic spikes.
                  </p>
                </div>
              </div>

              {/* Step 3: Install Dokploy */}
              <div id="install" className="mb-12 scroll-mt-20">
                <div className="flex items-center gap-3 mb-4">
                  <div className="flex h-8 w-8 items-center justify-center rounded-md bg-primary/10 border border-primary/15 text-sm font-mono font-semibold text-primary shrink-0">
                    3
                  </div>
                  <h2 className="text-2xl font-semibold tracking-tight">
                    Install Dokploy
                  </h2>
                </div>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  SSH into your VPS and run the one-line installer. This sets up Docker, Traefik
                  (reverse proxy), and the Dokploy dashboard:
                </p>
                <div className="rounded-xl border border-border/40 bg-card/80 overflow-hidden mb-4">
                  <div className="flex items-center gap-2 px-4 py-2.5 border-b border-border/30 bg-accent/30">
                    <div className="flex items-center gap-1.5">
                      <div className="h-2.5 w-2.5 rounded-full bg-[#ff5f57]" />
                      <div className="h-2.5 w-2.5 rounded-full bg-[#febc2e]" />
                      <div className="h-2.5 w-2.5 rounded-full bg-[#28c840]" />
                    </div>
                    <span className="ml-2 text-[11px] font-mono text-muted-foreground/40">server</span>
                  </div>
                  <div className="p-5 font-mono text-sm space-y-2">
                    <div>
                      <span className="text-primary/50 select-none">$ </span>
                      <span className="text-foreground/80">curl -sSL https://dokploy.com/install.sh | sh</span>
                    </div>
                  </div>
                </div>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  The installation takes 2-5 minutes. Once complete, you can access the Dokploy
                  dashboard at:
                </p>
                <div className="rounded-lg border border-border/30 bg-card/30 p-4 mb-4">
                  <code className="text-sm font-mono text-primary">http://your-server-ip:3000</code>
                </div>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Create your admin account on first visit. Save these credentials &mdash;
                  you&apos;ll need them to manage your deployments.
                </p>

                <div className="p-4 rounded-lg border border-[hsl(38,90%,50%)]/20 bg-[hsl(38,90%,50%)]/5">
                  <p className="text-sm text-foreground/80 leading-relaxed">
                    <span className="font-semibold text-[hsl(38,90%,50%)]">Important:</span>{' '}
                    If port 3000 is already in use on your server, Dokploy will fail to start.
                    Make sure no other service is using port 3000 before installing. You can check
                    with <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">sudo lsof -i :3000</code>.
                  </p>
                </div>
              </div>

              {/* Step 4: Create a Project */}
              <div id="create-project" className="mb-12 scroll-mt-20">
                <div className="flex items-center gap-3 mb-4">
                  <div className="flex h-8 w-8 items-center justify-center rounded-md bg-primary/10 border border-primary/15 text-sm font-mono font-semibold text-primary shrink-0">
                    4
                  </div>
                  <h2 className="text-2xl font-semibold tracking-tight">
                    Create a Project
                  </h2>
                </div>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  In the Dokploy dashboard:
                </p>
                <ol className="space-y-3 mb-6">
                  {[
                    'Click "Projects" in the sidebar, then "Create Project"',
                    'Give it a name (e.g., "my-grit-app")',
                    'Inside the project, click "+ Create Service" and select "Compose"',
                    'Choose "GitHub" as the provider (you\'ll need to connect your GitHub account on first use)',
                    'Select your repository and branch (usually "main")',
                    'Set the Compose file path to docker-compose.prod.yml (or wherever your production compose file lives)',
                  ].map((step, i) => (
                    <li key={i} className="flex items-start gap-3 text-[14px] text-muted-foreground">
                      <span className="flex h-5 w-5 items-center justify-center rounded-full bg-primary/10 text-[11px] font-mono font-semibold text-primary shrink-0 mt-0.5">
                        {i + 1}
                      </span>
                      <span dangerouslySetInnerHTML={{ __html: step.replace(/"/g, '&quot;').replace(/'/g, '&#39;') }} />
                    </li>
                  ))}
                </ol>

                <div className="p-4 rounded-lg border border-primary/20 bg-primary/5">
                  <p className="text-sm text-foreground/80 leading-relaxed">
                    <span className="font-semibold text-primary/90">Tip:</span>{' '}
                    Dokploy also supports deploying from a raw Docker Compose YAML pasted directly
                    into the dashboard. This is useful for quick testing before connecting GitHub.
                  </p>
                </div>
              </div>

              {/* Step 5: Docker Compose Config */}
              <div id="compose" className="mb-12 scroll-mt-20">
                <div className="flex items-center gap-3 mb-4">
                  <div className="flex h-8 w-8 items-center justify-center rounded-md bg-primary/10 border border-primary/15 text-sm font-mono font-semibold text-primary shrink-0">
                    5
                  </div>
                  <h2 className="text-2xl font-semibold tracking-tight">
                    Docker Compose Configuration
                  </h2>
                </div>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Your <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">docker-compose.prod.yml</code> defines
                  all services. Here&apos;s a production-ready Compose file optimized for Dokploy:
                </p>
                <CodeBlock language="yaml" filename="docker-compose.prod.yml" code={`services:
  # Go API
  api:
    build:
      context: ./apps/api
      dockerfile: Dockerfile
    ports:
      - "8080:8080"
    env_file:
      - .env
    depends_on:
      postgres:
        condition: service_healthy
      redis:
        condition: service_healthy
    restart: unless-stopped

  # Next.js Web App
  web:
    build:
      context: .
      dockerfile: apps/web/Dockerfile
    ports:
      - "3000:3000"
    environment:
      - NEXT_PUBLIC_API_URL=\${API_URL}
    restart: unless-stopped

  # Next.js Admin Panel
  admin:
    build:
      context: .
      dockerfile: apps/admin/Dockerfile
    ports:
      - "3001:3000"
    environment:
      - NEXT_PUBLIC_API_URL=\${API_URL}
    restart: unless-stopped

  # PostgreSQL
  postgres:
    image: postgres:16-alpine
    volumes:
      - postgres_data:/var/lib/postgresql/data
    environment:
      POSTGRES_USER: \${DB_USER:-postgres}
      POSTGRES_PASSWORD: \${DB_PASSWORD:-password}
      POSTGRES_DB: \${DB_NAME:-gritapp}
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U postgres"]
      interval: 5s
      timeout: 5s
      retries: 5
    restart: unless-stopped

  # Redis
  redis:
    image: redis:7-alpine
    volumes:
      - redis_data:/var/lib/redis/data
    command: redis-server --appendonly yes
    healthcheck:
      test: ["CMD", "redis-cli", "ping"]
      interval: 5s
      timeout: 5s
      retries: 5
    restart: unless-stopped

  # MinIO (S3-compatible file storage)
  minio:
    image: minio/minio:latest
    volumes:
      - minio_data:/data
    environment:
      MINIO_ROOT_USER: \${STORAGE_ACCESS_KEY:-minioadmin}
      MINIO_ROOT_PASSWORD: \${STORAGE_SECRET_KEY:-minioadmin}
    command: server /data --console-address ":9001"
    ports:
      - "9000:9000"
      - "9001:9001"
    restart: unless-stopped

volumes:
  postgres_data:
  redis_data:
  minio_data:`} />

                <div className="p-4 rounded-lg border border-primary/20 bg-primary/5 mt-4">
                  <p className="text-sm text-foreground/80 leading-relaxed">
                    <span className="font-semibold text-primary/90">Note:</span>{' '}
                    Dokploy uses Traefik as a reverse proxy. You do <strong>not</strong> need Caddy or
                    Nginx &mdash; Dokploy handles routing and SSL automatically. The ports are exposed
                    so Dokploy can route traffic to each service.
                  </p>
                </div>

                <h3 className="text-lg font-semibold tracking-tight mt-8 mb-3">
                  Build Context &amp; Monorepo
                </h3>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Notice that the <strong>API</strong> service uses <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">context: ./apps/api</code> because
                  the Go module is self-contained. The <strong>Next.js</strong> services use <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">context: .</code> (repo
                  root) because they need access to the shared <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">packages/shared</code>, the
                  root <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">pnpm-lock.yaml</code>, and <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">pnpm-workspace.yaml</code>.
                </p>

                <h3 className="text-lg font-semibold tracking-tight mt-8 mb-3">
                  Dockerfiles
                </h3>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Grit scaffolds production-ready Dockerfiles with <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">grit new</code>.
                  Here&apos;s what each one looks like:
                </p>

                <h4 className="text-base font-semibold tracking-tight mb-3 text-foreground/80">
                  Go API Dockerfile
                </h4>
                <p className="text-muted-foreground leading-relaxed mb-3 text-sm">
                  Multi-stage build &mdash; compiles a static Go binary, then copies it into a minimal Alpine image (~15 MB final):
                </p>
                <CodeBlock language="dockerfile" filename="apps/api/Dockerfile" code={`# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build static binary
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server ./cmd/server

# Run stage
FROM alpine:3.19

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

COPY --from=builder /app/server .

EXPOSE 8080

CMD ["./server"]`} />

                <h4 className="text-base font-semibold tracking-tight mt-6 mb-3 text-foreground/80">
                  Next.js Dockerfile (Web &amp; Admin)
                </h4>
                <p className="text-muted-foreground leading-relaxed mb-3 text-sm">
                  Three-stage build &mdash; installs deps, builds with standalone output, then creates a minimal runner (~120 MB final).
                  The same pattern is used for both <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">apps/web</code> and <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">apps/admin</code>:
                </p>
                <CodeBlock language="dockerfile" filename="apps/web/Dockerfile" code={`# Build stage
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

                <div className="p-4 rounded-lg border border-[hsl(38,90%,50%)]/20 bg-[hsl(38,90%,50%)]/5 mt-4">
                  <p className="text-sm text-foreground/80 leading-relaxed">
                    <span className="font-semibold text-[hsl(38,90%,50%)]">Critical:</span>{' '}
                    The Next.js Dockerfiles depend on <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">output: &quot;standalone&quot;</code> being
                    set in your <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">next.config.ts</code>. Grit sets this automatically
                    for both web and admin apps. Without it, <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">.next/standalone</code> won&apos;t
                    be generated and the Docker build will fail. If you removed it, add it back:
                  </p>
                  <div className="mt-3">
                    <CodeBlock language="typescript" filename="next.config.ts" code={`const nextConfig: NextConfig = {
  output: "standalone",  // Required for Docker builds
  reactStrictMode: true,
};`} />
                  </div>
                </div>

                <div className="p-4 rounded-lg border border-primary/20 bg-primary/5 mt-4">
                  <p className="text-sm text-foreground/80 leading-relaxed">
                    <span className="font-semibold text-primary/90">How standalone works:</span>{' '}
                    With <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">output: &quot;standalone&quot;</code>, Next.js
                    creates a self-contained <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">server.js</code> that includes
                    only the required <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">node_modules</code> files &mdash;
                    reducing the Docker image from ~1 GB to ~120 MB. The runner stage
                    only needs the <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">standalone</code>, <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">static</code>,
                    and <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">public</code> folders &mdash; no full <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">node_modules</code>.
                  </p>
                </div>
              </div>

              {/* Step 6: Environment Variables */}
              <div id="env-vars" className="mb-12 scroll-mt-20">
                <div className="flex items-center gap-3 mb-4">
                  <div className="flex h-8 w-8 items-center justify-center rounded-md bg-primary/10 border border-primary/15 text-sm font-mono font-semibold text-primary shrink-0">
                    6
                  </div>
                  <h2 className="text-2xl font-semibold tracking-tight">
                    Environment Variables
                  </h2>
                </div>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  In the Dokploy dashboard, navigate to your Compose service and click the
                  &quot;Environment&quot; tab. Copy the template below, paste it in, and fill in
                  your values. Every variable has a descriptive placeholder &mdash; replace anything
                  that says <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">CHANGE_THIS</code> or <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">YOUR_</code>.
                </p>
                <CodeBlock language="bash" filename="Dokploy → Environment tab (paste this)" code={`# ═══════════════════════════════════════════════════════════════
# GRIT CLOUD — Production Environment (Dokploy)
# ═══════════════════════════════════════════════════════════════

# App
APP_NAME=gritcloud
APP_ENV=production
APP_PORT=8080
APP_URL=https://api.yourdomain.com

# ─── Database (Docker Postgres service in Dokploy) ───────────
# Create a Postgres service in Dokploy, then use the internal hostname.
# Format: postgres://USER:PASSWORD@SERVICE_NAME:5432/DB_NAME?sslmode=disable
DATABASE_URL=postgres://gritcloud:CHANGE_THIS_STRONG_PASSWORD@gritcloud-postgres:5432/gritcloud?sslmode=disable

# ─── Redis (Docker Redis service in Dokploy) ─────────────────
# Create a Redis service in Dokploy, then use the internal hostname.
REDIS_URL=redis://gritcloud-redis:6379

# ─── JWT Authentication ──────────────────────────────────────
# Generate with: openssl rand -hex 32
JWT_SECRET=CHANGE_THIS_GENERATE_WITH_OPENSSL_RAND_HEX_32
JWT_ACCESS_EXPIRY=15m
JWT_REFRESH_EXPIRY=168h

# ─── Storage (Cloudflare R2) ─────────────────────────────────
STORAGE_DRIVER=r2
R2_ENDPOINT=https://YOUR_CF_ACCOUNT_ID.r2.cloudflarestorage.com
R2_ACCESS_KEY=YOUR_R2_ACCESS_KEY_ID
R2_SECRET_KEY=YOUR_R2_SECRET_ACCESS_KEY
R2_BUCKET=gritcloud-uploads
R2_REGION=auto
R2_PUBLIC_URL=https://pub-XXXX.r2.dev

# ─── OAuth2 (Google + GitHub) ────────────────────────────────
GOOGLE_CLIENT_ID=YOUR_GOOGLE_CLIENT_ID
GOOGLE_CLIENT_SECRET=YOUR_GOOGLE_CLIENT_SECRET
GITHUB_CLIENT_ID=YOUR_GITHUB_CLIENT_ID
GITHUB_CLIENT_SECRET=YOUR_GITHUB_CLIENT_SECRET
GITHUB_WEBHOOK_SECRET=YOUR_GITHUB_WEBHOOK_SECRET
OAUTH_FRONTEND_URL=https://cloud.yourdomain.com

# ─── Email (Resend) ──────────────────────────────────────────
RESEND_API_KEY=re_YOUR_RESEND_API_KEY
MAIL_FROM=noreply@yourdomain.com

# ─── CORS (all frontend origins) ─────────────────────────────
CORS_ORIGINS=https://yourdomain.com,https://cloud.yourdomain.com,https://ai.yourdomain.com

# ─── AI Provider ─────────────────────────────────────────────
AI_PROVIDER=claude
AI_API_KEY=sk-ant-YOUR_CLAUDE_API_KEY
AI_MODEL=claude-sonnet-4-5-20250929

# ─── Stripe Billing ──────────────────────────────────────────
STRIPE_SECRET_KEY=sk_live_YOUR_STRIPE_SECRET_KEY
STRIPE_PUBLISHABLE_KEY=pk_live_YOUR_STRIPE_PUBLISHABLE_KEY
STRIPE_WEBHOOK_SECRET=whsec_YOUR_STRIPE_WEBHOOK_SECRET
STRIPE_PRICE_STARTER=price_YOUR_STARTER_PRICE_ID
STRIPE_PRICE_PRO=price_YOUR_PRO_PRICE_ID
STRIPE_PRICE_BUSINESS=price_YOUR_BUSINESS_PRICE_ID

# ─── GORM Studio (DB browser — disable or secure in prod) ───
GORM_STUDIO_ENABLED=false
GORM_STUDIO_USERNAME=admin
GORM_STUDIO_PASSWORD=CHANGE_THIS_STRONG_PASSWORD

# ─── Pulse (Performance monitoring) ──────────────────────────
PULSE_ENABLED=true
PULSE_USERNAME=admin
PULSE_PASSWORD=CHANGE_THIS_STRONG_PASSWORD

# ─── Sentinel (WAF / Rate limiting) ──────────────────────────
SENTINEL_ENABLED=true
SENTINEL_USERNAME=admin
SENTINEL_PASSWORD=CHANGE_THIS_STRONG_PASSWORD
SENTINEL_SECRET_KEY=CHANGE_THIS_GENERATE_WITH_OPENSSL_RAND_HEX_32

# ─── Next.js Frontend Build-time Vars ────────────────────────
# These are needed when building the Next.js apps (cloud, ai, web)
NEXT_PUBLIC_API_URL=https://api.yourdomain.com
NEXT_PUBLIC_APP_URL=https://yourdomain.com
NEXT_PUBLIC_STRIPE_PUBLISHABLE_KEY=pk_live_YOUR_STRIPE_PUBLISHABLE_KEY`} />

                <div className="p-4 rounded-lg border border-[hsl(38,90%,50%)]/20 bg-[hsl(38,90%,50%)]/5 mt-4">
                  <p className="text-sm text-foreground/80 leading-relaxed">
                    <span className="font-semibold text-[hsl(38,90%,50%)]">Important:</span>{' '}
                    Notice that database and Redis hostnames use Docker <strong>service names</strong> (e.g.,{' '}
                    <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">postgres</code>,{' '}
                    <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">redis</code>,{' '}
                    <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">minio</code>)
                    instead of <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">localhost</code>.
                    Docker Compose networking resolves service names automatically.
                  </p>
                </div>
              </div>

              {/* Step 7: Configure Domains */}
              <div id="domains" className="mb-12 scroll-mt-20">
                <div className="flex items-center gap-3 mb-4">
                  <div className="flex h-8 w-8 items-center justify-center rounded-md bg-primary/10 border border-primary/15 text-sm font-mono font-semibold text-primary shrink-0">
                    7
                  </div>
                  <h2 className="text-2xl font-semibold tracking-tight">
                    Configure Domains
                  </h2>
                </div>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Dokploy uses Traefik to route domains to services and automatically provisions
                  SSL certificates via Let&apos;s Encrypt. For each service, go to the
                  &quot;Domains&quot; tab and add a domain:
                </p>
                <div className="rounded-lg border border-border/30 bg-card/30 overflow-hidden mb-6">
                  <table className="w-full text-sm">
                    <thead>
                      <tr className="border-b border-border/30 bg-accent/20">
                        <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Service</th>
                        <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Domain</th>
                        <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Container Port</th>
                      </tr>
                    </thead>
                    <tbody className="text-muted-foreground">
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-medium text-foreground/90">api</td>
                        <td className="px-4 py-2.5 font-mono text-xs">api.yourdomain.com</td>
                        <td className="px-4 py-2.5 font-mono text-xs">8080</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-medium text-foreground/90">web</td>
                        <td className="px-4 py-2.5 font-mono text-xs">yourdomain.com</td>
                        <td className="px-4 py-2.5 font-mono text-xs">3000</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-medium text-foreground/90">admin</td>
                        <td className="px-4 py-2.5 font-mono text-xs">admin.yourdomain.com</td>
                        <td className="px-4 py-2.5 font-mono text-xs">3000</td>
                      </tr>
                      <tr>
                        <td className="px-4 py-2.5 font-medium text-foreground/90">minio (optional)</td>
                        <td className="px-4 py-2.5 font-mono text-xs">files.yourdomain.com</td>
                        <td className="px-4 py-2.5 font-mono text-xs">9001</td>
                      </tr>
                    </tbody>
                  </table>
                </div>

                <p className="text-muted-foreground leading-relaxed mb-4">
                  For each domain in Dokploy:
                </p>
                <ol className="space-y-2 mb-6">
                  {[
                    'Click "Add Domain" on the service',
                    'Enter the domain (e.g., api.yourdomain.com)',
                    'Set HTTPS to "Let\'s Encrypt" for automatic SSL',
                    'Set the container port (the port your service listens on inside the container)',
                    'Save and deploy',
                  ].map((step, i) => (
                    <li key={i} className="flex items-start gap-3 text-[14px] text-muted-foreground">
                      <span className="flex h-5 w-5 items-center justify-center rounded-full bg-primary/10 text-[11px] font-mono font-semibold text-primary shrink-0 mt-0.5">
                        {i + 1}
                      </span>
                      {step}
                    </li>
                  ))}
                </ol>

                <div className="p-4 rounded-lg border border-primary/20 bg-primary/5">
                  <p className="text-sm text-foreground/80 leading-relaxed">
                    <span className="font-semibold text-primary/90">SSL note:</span>{' '}
                    Let&apos;s Encrypt certificates take 1-2 minutes to provision. If you see
                    a &quot;Not Secure&quot; warning immediately after deploying, wait a moment and
                    refresh. The certificate will be issued automatically once DNS is propagated.
                  </p>
                </div>
              </div>

              {/* Step 8: DNS Configuration */}
              <div id="dns" className="mb-12 scroll-mt-20">
                <div className="flex items-center gap-3 mb-4">
                  <div className="flex h-8 w-8 items-center justify-center rounded-md bg-primary/10 border border-primary/15 text-sm font-mono font-semibold text-primary shrink-0">
                    8
                  </div>
                  <h2 className="text-2xl font-semibold tracking-tight">
                    DNS Configuration
                  </h2>
                </div>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  In your domain registrar or DNS provider (Cloudflare, Namecheap, etc.),
                  create <strong>A records</strong> pointing to your VPS IP address:
                </p>
                <div className="rounded-lg border border-border/30 bg-card/30 overflow-hidden mb-6">
                  <table className="w-full text-sm">
                    <thead>
                      <tr className="border-b border-border/30 bg-accent/20">
                        <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Type</th>
                        <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Name</th>
                        <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Value</th>
                        <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Purpose</th>
                      </tr>
                    </thead>
                    <tbody className="text-muted-foreground">
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-mono text-xs">A</td>
                        <td className="px-4 py-2.5 font-mono text-xs">@</td>
                        <td className="px-4 py-2.5 font-mono text-xs">YOUR_VPS_IP</td>
                        <td className="px-4 py-2.5 text-xs">Web app (yourdomain.com)</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-mono text-xs">A</td>
                        <td className="px-4 py-2.5 font-mono text-xs">api</td>
                        <td className="px-4 py-2.5 font-mono text-xs">YOUR_VPS_IP</td>
                        <td className="px-4 py-2.5 text-xs">Go API (api.yourdomain.com)</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-mono text-xs">A</td>
                        <td className="px-4 py-2.5 font-mono text-xs">admin</td>
                        <td className="px-4 py-2.5 font-mono text-xs">YOUR_VPS_IP</td>
                        <td className="px-4 py-2.5 text-xs">Admin panel (admin.yourdomain.com)</td>
                      </tr>
                      <tr>
                        <td className="px-4 py-2.5 font-mono text-xs">A</td>
                        <td className="px-4 py-2.5 font-mono text-xs">files</td>
                        <td className="px-4 py-2.5 font-mono text-xs">YOUR_VPS_IP</td>
                        <td className="px-4 py-2.5 text-xs">MinIO console (optional)</td>
                      </tr>
                    </tbody>
                  </table>
                </div>

                <div className="p-4 rounded-lg border border-[hsl(38,90%,50%)]/20 bg-[hsl(38,90%,50%)]/5">
                  <p className="text-sm text-foreground/80 leading-relaxed">
                    <span className="font-semibold text-[hsl(38,90%,50%)]">Cloudflare users:</span>{' '}
                    If you&apos;re using Cloudflare, set the proxy status to <strong>&quot;DNS only&quot;</strong> (grey
                    cloud) for Let&apos;s Encrypt to work. If you want Cloudflare&apos;s proxy (orange cloud),
                    set the SSL mode to &quot;Full (Strict)&quot; in Cloudflare and configure Dokploy
                    to use Cloudflare-issued certificates instead of Let&apos;s Encrypt.
                  </p>
                </div>
              </div>

              {/* Step 9: Deploy */}
              <div id="deploy" className="mb-12 scroll-mt-20">
                <div className="flex items-center gap-3 mb-4">
                  <div className="flex h-8 w-8 items-center justify-center rounded-md bg-primary/10 border border-primary/15 text-sm font-mono font-semibold text-primary shrink-0">
                    9
                  </div>
                  <h2 className="text-2xl font-semibold tracking-tight">
                    Deploy
                  </h2>
                </div>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Once your environment variables and domains are configured, hit the
                  &quot;Deploy&quot; button in the Dokploy dashboard. Dokploy will:
                </p>
                <ol className="space-y-2 mb-6">
                  {[
                    'Pull your code from GitHub',
                    'Build all Docker images (Go API + Next.js apps)',
                    'Start all services defined in your Compose file',
                    'Configure Traefik routes for your domains',
                    'Provision SSL certificates via Let\'s Encrypt',
                  ].map((step, i) => (
                    <li key={i} className="flex items-start gap-3 text-[14px] text-muted-foreground">
                      <span className="flex h-5 w-5 items-center justify-center rounded-full bg-primary/10 text-[11px] font-mono font-semibold text-primary shrink-0 mt-0.5">
                        {i + 1}
                      </span>
                      {step}
                    </li>
                  ))}
                </ol>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  The first build takes longer (5-10 minutes) because Docker needs to download base
                  images and build from scratch. Subsequent deploys are faster thanks to Docker layer caching.
                </p>

                <h3 className="text-lg font-semibold tracking-tight mb-3">
                  Auto-Deploy on Git Push
                </h3>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  To enable automatic deployments when you push to GitHub:
                </p>
                <ol className="space-y-2 mb-6">
                  {[
                    'Go to your Compose service in Dokploy',
                    'Click the "General" tab',
                    'Enable "Auto Deploy" and select your branch (e.g., main)',
                    'Dokploy will set up a GitHub webhook automatically',
                  ].map((step, i) => (
                    <li key={i} className="flex items-start gap-3 text-[14px] text-muted-foreground">
                      <span className="flex h-5 w-5 items-center justify-center rounded-full bg-primary/10 text-[11px] font-mono font-semibold text-primary shrink-0 mt-0.5">
                        {i + 1}
                      </span>
                      {step}
                    </li>
                  ))}
                </ol>
                <p className="text-muted-foreground leading-relaxed">
                  Now every <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">git push origin main</code> will
                  trigger a new deployment automatically.
                </p>
              </div>

              {/* Step 10: Database Backups */}
              <div id="backups" className="mb-12 scroll-mt-20">
                <div className="flex items-center gap-3 mb-4">
                  <div className="flex h-8 w-8 items-center justify-center rounded-md bg-primary/10 border border-primary/15 text-sm font-mono font-semibold text-primary shrink-0">
                    10
                  </div>
                  <h2 className="text-2xl font-semibold tracking-tight">
                    Database Backups
                  </h2>
                </div>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Dokploy has built-in backup support for database services. Alternatively, you
                  can run manual backups using <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">pg_dump</code>:
                </p>

                <h3 className="text-lg font-semibold tracking-tight mb-3">
                  Manual Backup
                </h3>
                <div className="rounded-xl border border-border/40 bg-card/80 overflow-hidden mb-6">
                  <div className="flex items-center gap-2 px-4 py-2.5 border-b border-border/30 bg-accent/30">
                    <div className="flex items-center gap-1.5">
                      <div className="h-2.5 w-2.5 rounded-full bg-[#ff5f57]" />
                      <div className="h-2.5 w-2.5 rounded-full bg-[#febc2e]" />
                      <div className="h-2.5 w-2.5 rounded-full bg-[#28c840]" />
                    </div>
                    <span className="ml-2 text-[11px] font-mono text-muted-foreground/40">server</span>
                  </div>
                  <div className="p-5 font-mono text-sm space-y-2">
                    <div className="text-muted-foreground/50 text-xs"># Find the postgres container name</div>
                    <div>
                      <span className="text-primary/50 select-none">$ </span>
                      <span className="text-foreground/80">docker ps --filter name=postgres --format {'{{.Names}}'}</span>
                    </div>
                    <div className="mt-4 text-muted-foreground/50 text-xs"># Create a backup</div>
                    <div>
                      <span className="text-primary/50 select-none">$ </span>
                      <span className="text-foreground/80">docker exec -t CONTAINER_NAME pg_dumpall -c -U postgres {'>'} backup_$(date +%Y%m%d).sql</span>
                    </div>
                    <div className="mt-4 text-muted-foreground/50 text-xs"># Restore from backup</div>
                    <div>
                      <span className="text-primary/50 select-none">$ </span>
                      <span className="text-foreground/80">{'cat backup_20260223.sql | docker exec -i CONTAINER_NAME psql -U postgres'}</span>
                    </div>
                  </div>
                </div>

                <h3 className="text-lg font-semibold tracking-tight mb-3">
                  Automated Backups with Cron
                </h3>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Set up a cron job on your server to back up daily:
                </p>
                <CodeBlock filename="crontab -e" code={`# Daily backup at 3 AM — keeps 7 days of backups
0 3 * * * docker exec -t $(docker ps --filter name=postgres -q) pg_dumpall -c -U postgres | gzip > /backups/grit_$(date +\\%Y\\%m\\%d).sql.gz && find /backups -name "grit_*.sql.gz" -mtime +7 -delete`} />
              </div>

              {/* Step 11: Monitoring */}
              <div id="monitoring" className="mb-12 scroll-mt-20">
                <div className="flex items-center gap-3 mb-4">
                  <div className="flex h-8 w-8 items-center justify-center rounded-md bg-primary/10 border border-primary/15 text-sm font-mono font-semibold text-primary shrink-0">
                    11
                  </div>
                  <h2 className="text-2xl font-semibold tracking-tight">
                    Monitoring
                  </h2>
                </div>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Dokploy provides built-in monitoring through its dashboard. On top of that,
                  Grit includes additional observability tools:
                </p>
                <div className="rounded-lg border border-border/30 bg-card/30 overflow-hidden mb-6">
                  <table className="w-full text-sm">
                    <thead>
                      <tr className="border-b border-border/30 bg-accent/20">
                        <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Tool</th>
                        <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Access</th>
                        <th className="text-left px-4 py-2.5 font-medium text-foreground/80">What It Shows</th>
                      </tr>
                    </thead>
                    <tbody className="text-muted-foreground">
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-medium text-foreground/90">Dokploy Dashboard</td>
                        <td className="px-4 py-2.5 font-mono text-xs">http://your-ip:3000</td>
                        <td className="px-4 py-2.5 text-xs">Container logs, CPU/memory usage, deployment history</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-medium text-foreground/90">Pulse</td>
                        <td className="px-4 py-2.5 font-mono text-xs">api.yourdomain.com/pulse</td>
                        <td className="px-4 py-2.5 text-xs">Request tracing, DB metrics, error tracking, alerting</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-medium text-foreground/90">Sentinel</td>
                        <td className="px-4 py-2.5 font-mono text-xs">api.yourdomain.com/sentinel/ui</td>
                        <td className="px-4 py-2.5 text-xs">Security threats, blocked IPs, rate limiting</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-medium text-foreground/90">GORM Studio</td>
                        <td className="px-4 py-2.5 font-mono text-xs">api.yourdomain.com/studio</td>
                        <td className="px-4 py-2.5 text-xs">Database browser, SQL editor, data export</td>
                      </tr>
                      <tr>
                        <td className="px-4 py-2.5 font-medium text-foreground/90">API Docs</td>
                        <td className="px-4 py-2.5 font-mono text-xs">api.yourdomain.com/docs</td>
                        <td className="px-4 py-2.5 text-xs">Auto-generated OpenAPI spec + interactive UI</td>
                      </tr>
                    </tbody>
                  </table>
                </div>

                <h3 className="text-lg font-semibold tracking-tight mb-3">
                  Viewing Logs
                </h3>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Click any service in the Dokploy dashboard to see real-time logs. You can also
                  view logs from the command line:
                </p>
                <div className="rounded-xl border border-border/40 bg-card/80 overflow-hidden mb-4">
                  <div className="flex items-center gap-2 px-4 py-2.5 border-b border-border/30 bg-accent/30">
                    <div className="flex items-center gap-1.5">
                      <div className="h-2.5 w-2.5 rounded-full bg-[#ff5f57]" />
                      <div className="h-2.5 w-2.5 rounded-full bg-[#febc2e]" />
                      <div className="h-2.5 w-2.5 rounded-full bg-[#28c840]" />
                    </div>
                    <span className="ml-2 text-[11px] font-mono text-muted-foreground/40">server</span>
                  </div>
                  <div className="p-5 font-mono text-sm space-y-2">
                    <div className="text-muted-foreground/50 text-xs"># Follow API logs</div>
                    <div>
                      <span className="text-primary/50 select-none">$ </span>
                      <span className="text-foreground/80">docker compose -f docker-compose.prod.yml logs -f api</span>
                    </div>
                    <div className="mt-3 text-muted-foreground/50 text-xs"># View all service logs</div>
                    <div>
                      <span className="text-primary/50 select-none">$ </span>
                      <span className="text-foreground/80">docker compose -f docker-compose.prod.yml logs --tail 100</span>
                    </div>
                  </div>
                </div>
              </div>

              {/* Step 12: Tips & Troubleshooting */}
              <div id="troubleshooting" className="mb-12 scroll-mt-20">
                <div className="flex items-center gap-3 mb-4">
                  <div className="flex h-8 w-8 items-center justify-center rounded-md bg-primary/10 border border-primary/15 text-sm font-mono font-semibold text-primary shrink-0">
                    12
                  </div>
                  <h2 className="text-2xl font-semibold tracking-tight">
                    Tips &amp; Troubleshooting
                  </h2>
                </div>

                <div className="space-y-6">
                  {[
                    {
                      q: 'Build fails with "out of memory"',
                      a: 'Next.js builds are memory-intensive. Make sure your VPS has at least 4GB RAM. You can also add swap space: sudo fallocate -l 2G /swapfile && sudo chmod 600 /swapfile && sudo mkswap /swapfile && sudo swapon /swapfile',
                    },
                    {
                      q: 'SSL certificate not provisioning',
                      a: 'Make sure your DNS A records are pointing to the correct VPS IP and have propagated (use dig yourdomain.com to verify). If using Cloudflare, set the proxy to "DNS only" (grey cloud). Let\'s Encrypt needs to reach your server directly.',
                    },
                    {
                      q: 'Port 3000 already in use',
                      a: 'Dokploy\'s dashboard uses port 3000 by default. If your web app also uses port 3000, they won\'t conflict because the web app runs inside a Docker container with its own network. Dokploy routes traffic based on domain, not port.',
                    },
                    {
                      q: 'Database connection refused',
                      a: 'Make sure DATABASE_URL uses the Docker service name (postgres) instead of localhost. Inside Docker Compose, services communicate using service names as hostnames.',
                    },
                    {
                      q: 'How do I SSH into a running container?',
                      a: 'Use the Dokploy dashboard\'s "Terminal" tab, or from the command line: docker exec -it CONTAINER_NAME /bin/sh',
                    },
                    {
                      q: 'How do I redeploy without code changes?',
                      a: 'Click the "Redeploy" button in the Dokploy dashboard, or use the Dokploy API webhook URL to trigger a deployment programmatically.',
                    },
                    {
                      q: 'MinIO bucket not found',
                      a: 'MinIO doesn\'t auto-create buckets. After deployment, access the MinIO console (port 9001), sign in with your credentials, and create the bucket manually. Or add a startup script that uses the mc client.',
                    },
                  ].map((item) => (
                    <div key={item.q} className="rounded-xl border border-border/30 bg-card/30 p-4">
                      <h3 className="text-sm font-semibold mb-2 text-foreground/90">{item.q}</h3>
                      <p className="text-xs text-muted-foreground/70 leading-relaxed">{item.a}</p>
                    </div>
                  ))}
                </div>

                <div className="p-4 rounded-lg border border-primary/20 bg-primary/5 mt-6">
                  <p className="text-sm text-foreground/80 leading-relaxed">
                    <span className="font-semibold text-primary/90">Production checklist:</span>{' '}
                    Before going live, make sure you&apos;ve changed all default passwords
                    (JWT secrets, database password, MinIO keys, Sentinel/Pulse passwords),
                    enabled HTTPS for all domains, and set up automated database backups.
                    See the full{' '}
                    <Link href="/docs/infrastructure/deployment#checklist" className="text-primary hover:underline">
                      production checklist
                    </Link>.
                  </p>
                </div>
              </div>

              {/* Nav */}
              <div className="flex items-center justify-between pt-6 border-t border-border/30">
                <Button variant="ghost" size="sm" asChild className="text-muted-foreground/60 hover:text-foreground">
                  <Link href="/docs/infrastructure/deployment" className="gap-1.5">
                    <ArrowLeft className="h-3.5 w-3.5" />
                    Deployment
                  </Link>
                </Button>
                <Button variant="ghost" size="sm" asChild className="text-muted-foreground/60 hover:text-foreground">
                  <Link href="/docs/design/theme" className="gap-1.5">
                    Theme &amp; Colors
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
