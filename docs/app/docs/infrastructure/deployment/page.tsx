import Link from 'next/link'
import { ArrowLeft, ArrowRight } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'

export default function DeploymentPage() {
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
                Deployment
              </h1>
              <p className="text-lg text-muted-foreground leading-relaxed">
                Deploy your Grit application anywhere Docker runs &mdash; from a $5 VPS to
                managed cloud platforms. Go compiles to a single binary, making deployment fast and predictable.
              </p>
            </div>

            <div className="prose-grit">
              {/* Docker Deployment */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Docker Deployment
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  The simplest way to deploy is with the production Docker Compose file. It builds
                  and runs your entire stack &mdash; Go API, Next.js web, admin panel, PostgreSQL, and Redis.
                </p>

                <div className="flex items-center gap-3 mb-4">
                  <div className="flex h-7 w-7 items-center justify-center rounded-md bg-primary/10 border border-primary/15 text-xs font-mono font-semibold text-primary">
                    1
                  </div>
                  <h3 className="text-lg font-semibold tracking-tight">
                    Set production environment variables
                  </h3>
                </div>
                <div className="ml-10 rounded-xl border border-border/40 bg-card/80 overflow-hidden mb-6">
                  <div className="flex items-center gap-2 px-4 py-2.5 border-b border-border/30 bg-accent/30">
                    <span className="text-[11px] font-mono text-muted-foreground/40">.env.production</span>
                  </div>
                  <pre className="p-5 text-sm font-mono text-foreground/80 overflow-x-auto">{`APP_ENV=production
APP_PORT=8080

# Use strong, unique values in production
DATABASE_URL=postgres://user:strongpassword@postgres:5432/myapp?sslmode=require
REDIS_URL=redis://redis:6379
JWT_SECRET=generate-a-64-char-random-string-here

# Storage (Cloudflare R2 or AWS S3 recommended)
STORAGE_DRIVER=r2
R2_ENDPOINT=https://your-account.r2.cloudflarestorage.com
R2_ACCESS_KEY=your-access-key
R2_SECRET_KEY=your-secret-key
R2_BUCKET=myapp-uploads

# Email
RESEND_API_KEY=re_your_production_api_key
MAIL_FROM=noreply@yourdomain.com

# CORS (your production domains)
CORS_ORIGINS=https://yourdomain.com,https://admin.yourdomain.com

# Disable GORM Studio in production
GORM_STUDIO_ENABLED=false`}</pre>
                </div>

                <div className="flex items-center gap-3 mb-4">
                  <div className="flex h-7 w-7 items-center justify-center rounded-md bg-primary/10 border border-primary/15 text-xs font-mono font-semibold text-primary">
                    2
                  </div>
                  <h3 className="text-lg font-semibold tracking-tight">
                    Build and start production containers
                  </h3>
                </div>
                <div className="ml-10 rounded-xl border border-border/40 bg-card/80 overflow-hidden glow-purple-sm mb-6">
                  <div className="flex items-center gap-2 px-4 py-2.5 border-b border-border/30 bg-accent/30">
                    <div className="flex items-center gap-1.5">
                      <div className="h-2.5 w-2.5 rounded-full bg-[#ff5f57]" />
                      <div className="h-2.5 w-2.5 rounded-full bg-[#febc2e]" />
                      <div className="h-2.5 w-2.5 rounded-full bg-[#28c840]" />
                    </div>
                    <span className="ml-2 text-[11px] font-mono text-muted-foreground/40">terminal</span>
                  </div>
                  <div className="p-5 font-mono text-sm space-y-2">
                    <div>
                      <span className="text-primary/50 select-none">$ </span>
                      <span className="text-foreground/80">docker compose -f docker-compose.prod.yml up -d --build</span>
                    </div>
                  </div>
                </div>

                <div className="flex items-center gap-3 mb-4">
                  <div className="flex h-7 w-7 items-center justify-center rounded-md bg-primary/10 border border-primary/15 text-xs font-mono font-semibold text-primary">
                    3
                  </div>
                  <h3 className="text-lg font-semibold tracking-tight">
                    Verify everything is running
                  </h3>
                </div>
                <div className="ml-10 rounded-xl border border-border/40 bg-card/80 overflow-hidden mb-6">
                  <div className="flex items-center gap-2 px-4 py-2.5 border-b border-border/30 bg-accent/30">
                    <div className="flex items-center gap-1.5">
                      <div className="h-2.5 w-2.5 rounded-full bg-[#ff5f57]" />
                      <div className="h-2.5 w-2.5 rounded-full bg-[#febc2e]" />
                      <div className="h-2.5 w-2.5 rounded-full bg-[#28c840]" />
                    </div>
                    <span className="ml-2 text-[11px] font-mono text-muted-foreground/40">terminal</span>
                  </div>
                  <div className="p-5 font-mono text-sm space-y-2">
                    <div>
                      <span className="text-primary/50 select-none">$ </span>
                      <span className="text-foreground/80">docker compose -f docker-compose.prod.yml ps</span>
                    </div>
                    <div>
                      <span className="text-primary/50 select-none">$ </span>
                      <span className="text-foreground/80">curl http://localhost:8080/api/health</span>
                    </div>
                  </div>
                </div>
              </div>

              {/* VPS Deployment */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  VPS Deployment
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  A VPS (Virtual Private Server) is the most cost-effective way to deploy a Grit application.
                  Services like DigitalOcean, Hetzner, or Linode offer capable servers starting at $5/month.
                </p>

                <h3 className="text-lg font-semibold tracking-tight mb-3 mt-6">
                  Recommended VPS Providers
                </h3>
                <div className="rounded-lg border border-border/30 bg-card/30 overflow-hidden mb-6">
                  <table className="w-full text-sm">
                    <thead>
                      <tr className="border-b border-border/30 bg-accent/20">
                        <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Provider</th>
                        <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Starting At</th>
                        <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Best For</th>
                      </tr>
                    </thead>
                    <tbody className="text-muted-foreground">
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-medium text-foreground/90">Hetzner</td>
                        <td className="px-4 py-2.5 font-mono text-xs">$4/mo</td>
                        <td className="px-4 py-2.5">Best price-to-performance ratio</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-medium text-foreground/90">DigitalOcean</td>
                        <td className="px-4 py-2.5 font-mono text-xs">$6/mo</td>
                        <td className="px-4 py-2.5">Great docs, managed databases available</td>
                      </tr>
                      <tr>
                        <td className="px-4 py-2.5 font-medium text-foreground/90">Linode (Akamai)</td>
                        <td className="px-4 py-2.5 font-mono text-xs">$5/mo</td>
                        <td className="px-4 py-2.5">Good global coverage, reliable</td>
                      </tr>
                    </tbody>
                  </table>
                </div>

                <h3 className="text-lg font-semibold tracking-tight mb-3 mt-6">
                  Server Setup
                </h3>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  After provisioning a VPS (Ubuntu 22.04+ recommended), install Docker and deploy:
                </p>
                <div className="rounded-xl border border-border/40 bg-card/80 overflow-hidden mb-6">
                  <div className="flex items-center gap-2 px-4 py-2.5 border-b border-border/30 bg-accent/30">
                    <div className="flex items-center gap-1.5">
                      <div className="h-2.5 w-2.5 rounded-full bg-[#ff5f57]" />
                      <div className="h-2.5 w-2.5 rounded-full bg-[#febc2e]" />
                      <div className="h-2.5 w-2.5 rounded-full bg-[#28c840]" />
                    </div>
                    <span className="ml-2 text-[11px] font-mono text-muted-foreground/40">terminal — VPS setup</span>
                  </div>
                  <div className="p-5 font-mono text-sm space-y-2">
                    <div className="text-muted-foreground/50 text-xs mb-2"># Install Docker</div>
                    <div>
                      <span className="text-primary/50 select-none">$ </span>
                      <span className="text-foreground/80">curl -fsSL https://get.docker.com | sh</span>
                    </div>
                    <div className="mt-3 text-muted-foreground/50 text-xs mb-2"># Clone your project</div>
                    <div>
                      <span className="text-primary/50 select-none">$ </span>
                      <span className="text-foreground/80">git clone https://github.com/you/myapp.git</span>
                    </div>
                    <div>
                      <span className="text-primary/50 select-none">$ </span>
                      <span className="text-foreground/80">cd myapp</span>
                    </div>
                    <div className="mt-3 text-muted-foreground/50 text-xs mb-2"># Create .env and deploy</div>
                    <div>
                      <span className="text-primary/50 select-none">$ </span>
                      <span className="text-foreground/80">cp .env.example .env</span>
                    </div>
                    <div>
                      <span className="text-primary/50 select-none">$ </span>
                      <span className="text-foreground/80">nano .env  # Edit with production values</span>
                    </div>
                    <div>
                      <span className="text-primary/50 select-none">$ </span>
                      <span className="text-foreground/80">docker compose -f docker-compose.prod.yml up -d --build</span>
                    </div>
                  </div>
                </div>
              </div>

              {/* Cloud Platforms */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Cloud Platform Deployment
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  If you prefer managed infrastructure, several cloud platforms support Grit&apos;s architecture.
                  You will need to deploy the Go API and Next.js apps as separate services.
                </p>
                <div className="space-y-4">
                  {[
                    {
                      name: 'Railway',
                      desc: 'Deploy from GitHub with zero config. Supports Go, Node.js, PostgreSQL, and Redis as managed services. Great developer experience with automatic deployments on push.',
                      fit: 'Best for solo developers and small teams who want simplicity.',
                    },
                    {
                      name: 'Render',
                      desc: 'Deploy Go as a Web Service, Next.js as a Static Site or Web Service. Managed PostgreSQL and Redis available. Free tier for experiments.',
                      fit: 'Best for side projects and startups with a budget.',
                    },
                    {
                      name: 'Fly.io',
                      desc: 'Container-based deployment with edge computing. Deploy your Docker images directly. Fly Postgres and Upstash Redis integrations built in.',
                      fit: 'Best for apps that need global edge distribution.',
                    },
                    {
                      name: 'Coolify (Self-Hosted)',
                      desc: 'Open-source alternative to Heroku/Netlify that you host on your own VPS. Docker-based deployments, built-in databases, and automatic SSL. Free.',
                      fit: 'Best for developers who want PaaS convenience with VPS pricing.',
                    },
                  ].map((platform) => (
                    <div key={platform.name} className="p-4 rounded-lg border border-border/30 bg-card/30">
                      <h3 className="text-sm font-semibold mb-1.5">{platform.name}</h3>
                      <p className="text-xs text-muted-foreground/70 leading-relaxed mb-2">{platform.desc}</p>
                      <p className="text-xs text-primary/70 leading-relaxed">{platform.fit}</p>
                    </div>
                  ))}
                </div>
              </div>

              {/* Environment Variables */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Production Environment Variables
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  These are the critical environment variables you must set for a production deployment:
                </p>
                <div className="rounded-lg border border-border/30 bg-card/30 overflow-hidden">
                  <table className="w-full text-sm">
                    <thead>
                      <tr className="border-b border-border/30 bg-accent/20">
                        <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Variable</th>
                        <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Required</th>
                        <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Notes</th>
                      </tr>
                    </thead>
                    <tbody className="text-muted-foreground">
                      {[
                        { var: 'APP_ENV', req: 'Yes', note: 'Set to "production"' },
                        { var: 'DATABASE_URL', req: 'Yes', note: 'Full PostgreSQL connection string with sslmode=require' },
                        { var: 'JWT_SECRET', req: 'Yes', note: 'Random 64+ character string. Never reuse across environments.' },
                        { var: 'REDIS_URL', req: 'Yes', note: 'Redis connection URL for caching and job queues' },
                        { var: 'CORS_ORIGINS', req: 'Yes', note: 'Comma-separated list of allowed frontend domains' },
                        { var: 'RESEND_API_KEY', req: 'If emails', note: 'Resend API key for transactional emails' },
                        { var: 'STORAGE_DRIVER', req: 'If uploads', note: 'r2, b2, or minio for file storage' },
                        { var: 'GORM_STUDIO_ENABLED', req: 'No', note: 'Set to "false" in production' },
                      ].map((item) => (
                        <tr key={item.var} className="border-b border-border/20 last:border-b-0">
                          <td className="px-4 py-2.5 font-mono text-xs text-primary/80">{item.var}</td>
                          <td className="px-4 py-2.5 font-mono text-xs">{item.req}</td>
                          <td className="px-4 py-2.5 text-xs">{item.note}</td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>
              </div>

              {/* SSL/TLS */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  SSL/TLS Setup
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Use a reverse proxy like Nginx or Caddy in front of your Grit application for SSL termination.
                  Caddy is recommended because it handles certificate provisioning automatically via Let&apos;s Encrypt.
                </p>
                <div className="rounded-xl border border-border/40 bg-card/80 overflow-hidden mb-6">
                  <div className="flex items-center gap-2 px-4 py-2.5 border-b border-border/30 bg-accent/30">
                    <span className="text-[11px] font-mono text-muted-foreground/40">Caddyfile</span>
                  </div>
                  <pre className="p-5 text-sm font-mono text-foreground/80 overflow-x-auto">{`# API
api.yourdomain.com {
    reverse_proxy localhost:8080
}

# Web app
yourdomain.com {
    reverse_proxy localhost:3000
}

# Admin panel
admin.yourdomain.com {
    reverse_proxy localhost:3001
}`}</pre>
                </div>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Install and start Caddy:
                </p>
                <div className="rounded-xl border border-border/40 bg-card/80 overflow-hidden mb-4">
                  <div className="flex items-center gap-2 px-4 py-2.5 border-b border-border/30 bg-accent/30">
                    <div className="flex items-center gap-1.5">
                      <div className="h-2.5 w-2.5 rounded-full bg-[#ff5f57]" />
                      <div className="h-2.5 w-2.5 rounded-full bg-[#febc2e]" />
                      <div className="h-2.5 w-2.5 rounded-full bg-[#28c840]" />
                    </div>
                    <span className="ml-2 text-[11px] font-mono text-muted-foreground/40">terminal</span>
                  </div>
                  <div className="p-5 font-mono text-sm space-y-2">
                    <div>
                      <span className="text-primary/50 select-none">$ </span>
                      <span className="text-foreground/80">sudo apt install -y debian-keyring debian-archive-keyring apt-transport-https</span>
                    </div>
                    <div>
                      <span className="text-primary/50 select-none">$ </span>
                      <span className="text-foreground/80">curl -1sLf &apos;https://dl.cloudsmith.io/public/caddy/stable/gpg.key&apos; | sudo gpg --dearmor -o /usr/share/keyrings/caddy-stable-archive-keyring.gpg</span>
                    </div>
                    <div>
                      <span className="text-primary/50 select-none">$ </span>
                      <span className="text-foreground/80">sudo apt update &amp;&amp; sudo apt install caddy</span>
                    </div>
                    <div>
                      <span className="text-primary/50 select-none">$ </span>
                      <span className="text-foreground/80">sudo systemctl enable --now caddy</span>
                    </div>
                  </div>
                </div>
                <p className="text-sm text-muted-foreground/60">
                  Caddy will automatically obtain and renew SSL certificates from Let&apos;s Encrypt.
                  Make sure your DNS A records point to your server&apos;s IP address before starting Caddy.
                </p>
              </div>

              {/* Database Backup */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Database Backup Strategies
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Always back up your production database. Here are common approaches:
                </p>

                <h3 className="text-lg font-semibold tracking-tight mb-3 mt-6">
                  Manual Backup with pg_dump
                </h3>
                <div className="rounded-xl border border-border/40 bg-card/80 overflow-hidden mb-6">
                  <div className="flex items-center gap-2 px-4 py-2.5 border-b border-border/30 bg-accent/30">
                    <div className="flex items-center gap-1.5">
                      <div className="h-2.5 w-2.5 rounded-full bg-[#ff5f57]" />
                      <div className="h-2.5 w-2.5 rounded-full bg-[#febc2e]" />
                      <div className="h-2.5 w-2.5 rounded-full bg-[#28c840]" />
                    </div>
                    <span className="ml-2 text-[11px] font-mono text-muted-foreground/40">terminal</span>
                  </div>
                  <div className="p-5 font-mono text-sm space-y-2">
                    <div className="text-muted-foreground/50 text-xs mb-2"># Create a backup</div>
                    <div>
                      <span className="text-primary/50 select-none">$ </span>
                      <span className="text-foreground/80">docker compose exec postgres pg_dump -U grit myapp &gt; backup.sql</span>
                    </div>
                    <div className="mt-3 text-muted-foreground/50 text-xs mb-2"># Restore from backup</div>
                    <div>
                      <span className="text-primary/50 select-none">$ </span>
                      <span className="text-foreground/80">docker compose exec -T postgres psql -U grit myapp &lt; backup.sql</span>
                    </div>
                  </div>
                </div>

                <h3 className="text-lg font-semibold tracking-tight mb-3 mt-6">
                  Automated Daily Backups (Cron)
                </h3>
                <div className="rounded-xl border border-border/40 bg-card/80 overflow-hidden mb-4">
                  <div className="flex items-center gap-2 px-4 py-2.5 border-b border-border/30 bg-accent/30">
                    <span className="text-[11px] font-mono text-muted-foreground/40">crontab -e</span>
                  </div>
                  <pre className="p-5 text-sm font-mono text-foreground/80 overflow-x-auto">{`# Daily backup at 3 AM, keep last 7 days
0 3 * * * cd /opt/myapp && docker compose exec -T postgres pg_dump -U grit myapp | gzip > /backups/myapp-$(date +\\%Y\\%m\\%d).sql.gz && find /backups -name "myapp-*.sql.gz" -mtime +7 -delete`}</pre>
                </div>

                <h3 className="text-lg font-semibold tracking-tight mb-3 mt-6">
                  Managed Database Backups
                </h3>
                <p className="text-muted-foreground leading-relaxed">
                  If you use a managed PostgreSQL service (DigitalOcean Managed Databases, Supabase, Neon, etc.),
                  backups are handled automatically. This is the recommended approach for production workloads.
                </p>
              </div>

              {/* Monitoring */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Monitoring &amp; Logging
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Grit&apos;s Go API includes a built-in request logger middleware that logs every HTTP request with
                  method, path, status code, and duration. For production monitoring, consider these tools:
                </p>
                <div className="space-y-3">
                  {[
                    { name: 'Docker Logs', desc: 'Built-in. Use "docker compose logs -f api" to stream API logs in real time.' },
                    { name: 'Uptime Kuma', desc: 'Self-hosted uptime monitoring. Set up health checks against /api/health to get notified of downtime.' },
                    { name: 'Grafana + Prometheus', desc: 'For advanced metrics. Add a /metrics endpoint to your Go API and visualize request rates, latency, and error rates.' },
                    { name: 'Sentry', desc: 'Error tracking for both Go API and Next.js frontends. Captures stack traces, context, and user info for every error.' },
                    { name: 'Better Stack (Logtail)', desc: 'Managed log aggregation. Stream Docker logs for search, alerting, and dashboards.' },
                  ].map((item) => (
                    <div key={item.name} className="p-4 rounded-lg border border-border/30 bg-card/30">
                      <h3 className="text-sm font-semibold mb-1.5">{item.name}</h3>
                      <p className="text-xs text-muted-foreground/70 leading-relaxed">{item.desc}</p>
                    </div>
                  ))}
                </div>
              </div>

              {/* Production Checklist */}
              <div className="mb-10">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Production Checklist
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Before going live, make sure you have completed these steps:
                </p>
                <div className="space-y-2.5">
                  {[
                    { category: 'Security', items: [
                      'Generate a strong, unique JWT_SECRET (64+ characters)',
                      'Set APP_ENV=production',
                      'Set GORM_STUDIO_ENABLED=false',
                      'Use sslmode=require in DATABASE_URL',
                      'Restrict CORS_ORIGINS to your actual domains',
                      'Change default database credentials from grit/grit',
                      'Enable SSL/TLS via Caddy or Nginx',
                    ]},
                    { category: 'Infrastructure', items: [
                      'Set up automated database backups',
                      'Configure uptime monitoring on /api/health',
                      'Set up error tracking (Sentry or similar)',
                      'Enable Docker restart policies (unless-stopped)',
                      'Point DNS records to your server IP',
                    ]},
                    { category: 'Performance', items: [
                      'Set appropriate database connection pool limits',
                      'Add indexes to frequently queried columns',
                      'Enable Redis caching for read-heavy endpoints',
                      'Use a CDN for static assets (Cloudflare recommended)',
                    ]},
                  ].map((section) => (
                    <div key={section.category} className="mb-6">
                      <h3 className="text-sm font-semibold mb-3 text-foreground/90">{section.category}</h3>
                      <div className="space-y-2">
                        {section.items.map((item) => (
                          <div key={item} className="flex items-start gap-2.5 text-[13px] text-muted-foreground pl-2">
                            <div className="h-4 w-4 rounded border border-border/40 bg-card/30 mt-0.5 shrink-0" />
                            {item}
                          </div>
                        ))}
                      </div>
                    </div>
                  ))}
                </div>
              </div>

              {/* Nav */}
              <div className="flex items-center justify-between pt-6 border-t border-border/30">
                <Button variant="ghost" size="sm" asChild className="text-muted-foreground/60 hover:text-foreground">
                  <Link href="/docs/infrastructure/database" className="gap-1.5">
                    <ArrowLeft className="h-3.5 w-3.5" />
                    Database &amp; Migrations
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
