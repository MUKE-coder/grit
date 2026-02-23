import Link from 'next/link'
import { ArrowLeft, ArrowRight } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'
import { CodeBlock } from '@/components/code-block'
import { getDocMetadata } from '@/config/docs-metadata'

export const metadata = getDocMetadata('/docs/infrastructure/deployment')

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
                managed cloud platforms. This guide walks you through a complete production
                deployment on a VPS, step by step.
              </p>
            </div>

            <div className="prose-grit">

              {/* Table of Contents */}
              <div className="mb-12 rounded-xl border border-border/40 bg-card/30 p-6">
                <h2 className="text-lg font-semibold tracking-tight mb-3">On This Page</h2>
                <div className="grid gap-1.5 sm:grid-cols-2">
                  {[
                    { label: '1. Choose a VPS Provider', id: '#providers' },
                    { label: '2. Initial Server Setup', id: '#server-setup' },
                    { label: '3. Install Docker', id: '#docker' },
                    { label: '4. Clone & Configure', id: '#clone' },
                    { label: '5. Deploy with Docker Compose', id: '#deploy' },
                    { label: '6. Set Up SSL with Caddy', id: '#ssl' },
                    { label: '7. Configure Firewall (UFW)', id: '#firewall' },
                    { label: '8. Automated Deployments (CI/CD)', id: '#cicd' },
                    { label: '9. Database Backups', id: '#backups' },
                    { label: '10. Monitoring & Logging', id: '#monitoring' },
                    { label: '11. Zero-Downtime Updates', id: '#updates' },
                    { label: '12. Production Checklist', id: '#checklist' },
                  ].map((item) => (
                    <a key={item.id} href={item.id} className="text-[13px] text-muted-foreground/70 hover:text-primary transition-colors py-0.5">
                      {item.label}
                    </a>
                  ))}
                </div>
              </div>

              {/* Step 1: Choose a VPS */}
              <div id="providers" className="mb-12 scroll-mt-20">
                <div className="flex items-center gap-3 mb-4">
                  <div className="flex h-8 w-8 items-center justify-center rounded-md bg-primary/10 border border-primary/15 text-sm font-mono font-semibold text-primary shrink-0">
                    1
                  </div>
                  <h2 className="text-2xl font-semibold tracking-tight">
                    Choose a VPS Provider
                  </h2>
                </div>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  A VPS (Virtual Private Server) is the most cost-effective way to deploy a Grit application.
                  You get full control over the server, and a $5-10/month VPS can handle thousands of users.
                </p>

                <div className="rounded-lg border border-border/30 bg-card/30 overflow-hidden mb-4">
                  <table className="w-full text-sm">
                    <thead>
                      <tr className="border-b border-border/30 bg-accent/20">
                        <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Provider</th>
                        <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Starting At</th>
                        <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Specs (entry)</th>
                        <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Best For</th>
                      </tr>
                    </thead>
                    <tbody className="text-muted-foreground">
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-medium text-foreground/90">Hetzner</td>
                        <td className="px-4 py-2.5 font-mono text-xs">$4/mo</td>
                        <td className="px-4 py-2.5 text-xs">2 vCPU, 4GB RAM, 40GB SSD</td>
                        <td className="px-4 py-2.5 text-xs">Best price-to-performance</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-medium text-foreground/90">DigitalOcean</td>
                        <td className="px-4 py-2.5 font-mono text-xs">$6/mo</td>
                        <td className="px-4 py-2.5 text-xs">1 vCPU, 1GB RAM, 25GB SSD</td>
                        <td className="px-4 py-2.5 text-xs">Great docs, managed DBs</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-medium text-foreground/90">Linode (Akamai)</td>
                        <td className="px-4 py-2.5 font-mono text-xs">$5/mo</td>
                        <td className="px-4 py-2.5 text-xs">1 vCPU, 1GB RAM, 25GB SSD</td>
                        <td className="px-4 py-2.5 text-xs">Good global coverage</td>
                      </tr>
                      <tr>
                        <td className="px-4 py-2.5 font-medium text-foreground/90">Vultr</td>
                        <td className="px-4 py-2.5 font-mono text-xs">$6/mo</td>
                        <td className="px-4 py-2.5 text-xs">1 vCPU, 1GB RAM, 25GB SSD</td>
                        <td className="px-4 py-2.5 text-xs">Many locations, hourly billing</td>
                      </tr>
                    </tbody>
                  </table>
                </div>

                <div className="p-4 rounded-lg border border-primary/20 bg-primary/5">
                  <p className="text-sm text-foreground/80 leading-relaxed">
                    <span className="font-semibold text-primary/90">Recommendation:</span>{' '}
                    For a Grit monorepo (Go API + 2 Next.js apps + PostgreSQL + Redis), we recommend at least
                    <span className="font-mono text-primary/80"> 2 vCPU / 4GB RAM</span>.
                    Hetzner&apos;s CX22 ($4/mo) or DigitalOcean&apos;s Basic ($24/mo for 2 vCPU / 4GB) are solid choices.
                  </p>
                </div>
              </div>

              {/* Step 2: Initial Server Setup */}
              <div id="server-setup" className="mb-12 scroll-mt-20">
                <div className="flex items-center gap-3 mb-4">
                  <div className="flex h-8 w-8 items-center justify-center rounded-md bg-primary/10 border border-primary/15 text-sm font-mono font-semibold text-primary shrink-0">
                    2
                  </div>
                  <h2 className="text-2xl font-semibold tracking-tight">
                    Initial Server Setup
                  </h2>
                </div>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  After provisioning your VPS with <strong>Ubuntu 22.04+</strong> (recommended), SSH in as root and
                  set up a non-root user with sudo access.
                </p>

                <h3 className="text-lg font-semibold tracking-tight mb-3">
                  2a. SSH into your server
                </h3>
                <div className="rounded-xl border border-border/40 bg-card/80 overflow-hidden mb-6">
                  <div className="flex items-center gap-2 px-4 py-2.5 border-b border-border/30 bg-accent/30">
                    <div className="flex items-center gap-1.5">
                      <div className="h-2.5 w-2.5 rounded-full bg-[#ff5f57]" />
                      <div className="h-2.5 w-2.5 rounded-full bg-[#febc2e]" />
                      <div className="h-2.5 w-2.5 rounded-full bg-[#28c840]" />
                    </div>
                    <span className="ml-2 text-[11px] font-mono text-muted-foreground/40">local machine</span>
                  </div>
                  <div className="p-5 font-mono text-sm space-y-2">
                    <div>
                      <span className="text-primary/50 select-none">$ </span>
                      <span className="text-foreground/80">ssh root@your-server-ip</span>
                    </div>
                  </div>
                </div>

                <h3 className="text-lg font-semibold tracking-tight mb-3">
                  2b. Create a deploy user
                </h3>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Running everything as root is a security risk. Create a dedicated deploy user:
                </p>
                <div className="rounded-xl border border-border/40 bg-card/80 overflow-hidden mb-6">
                  <div className="flex items-center gap-2 px-4 py-2.5 border-b border-border/30 bg-accent/30">
                    <div className="flex items-center gap-1.5">
                      <div className="h-2.5 w-2.5 rounded-full bg-[#ff5f57]" />
                      <div className="h-2.5 w-2.5 rounded-full bg-[#febc2e]" />
                      <div className="h-2.5 w-2.5 rounded-full bg-[#28c840]" />
                    </div>
                    <span className="ml-2 text-[11px] font-mono text-muted-foreground/40">server (as root)</span>
                  </div>
                  <div className="p-5 font-mono text-sm space-y-2">
                    <div className="text-muted-foreground/50 text-xs mb-2"># Create user and add to sudo group</div>
                    <div>
                      <span className="text-primary/50 select-none"># </span>
                      <span className="text-foreground/80">adduser deploy</span>
                    </div>
                    <div>
                      <span className="text-primary/50 select-none"># </span>
                      <span className="text-foreground/80">usermod -aG sudo deploy</span>
                    </div>
                    <div className="mt-3 text-muted-foreground/50 text-xs mb-2"># Copy SSH keys to the new user</div>
                    <div>
                      <span className="text-primary/50 select-none"># </span>
                      <span className="text-foreground/80">rsync --archive --chown=deploy:deploy ~/.ssh /home/deploy</span>
                    </div>
                    <div className="mt-3 text-muted-foreground/50 text-xs mb-2"># Now log out and SSH as deploy</div>
                    <div>
                      <span className="text-primary/50 select-none"># </span>
                      <span className="text-foreground/80">exit</span>
                    </div>
                  </div>
                </div>

                <div className="rounded-xl border border-border/40 bg-card/80 overflow-hidden mb-6">
                  <div className="flex items-center gap-2 px-4 py-2.5 border-b border-border/30 bg-accent/30">
                    <div className="flex items-center gap-1.5">
                      <div className="h-2.5 w-2.5 rounded-full bg-[#ff5f57]" />
                      <div className="h-2.5 w-2.5 rounded-full bg-[#febc2e]" />
                      <div className="h-2.5 w-2.5 rounded-full bg-[#28c840]" />
                    </div>
                    <span className="ml-2 text-[11px] font-mono text-muted-foreground/40">local machine</span>
                  </div>
                  <div className="p-5 font-mono text-sm">
                    <span className="text-primary/50 select-none">$ </span>
                    <span className="text-foreground/80">ssh deploy@your-server-ip</span>
                  </div>
                </div>

                <h3 className="text-lg font-semibold tracking-tight mb-3">
                  2c. Update packages and set timezone
                </h3>
                <div className="rounded-xl border border-border/40 bg-card/80 overflow-hidden mb-6">
                  <div className="flex items-center gap-2 px-4 py-2.5 border-b border-border/30 bg-accent/30">
                    <div className="flex items-center gap-1.5">
                      <div className="h-2.5 w-2.5 rounded-full bg-[#ff5f57]" />
                      <div className="h-2.5 w-2.5 rounded-full bg-[#febc2e]" />
                      <div className="h-2.5 w-2.5 rounded-full bg-[#28c840]" />
                    </div>
                    <span className="ml-2 text-[11px] font-mono text-muted-foreground/40">server (as deploy)</span>
                  </div>
                  <div className="p-5 font-mono text-sm space-y-2">
                    <div>
                      <span className="text-primary/50 select-none">$ </span>
                      <span className="text-foreground/80">sudo apt update &amp;&amp; sudo apt upgrade -y</span>
                    </div>
                    <div>
                      <span className="text-primary/50 select-none">$ </span>
                      <span className="text-foreground/80">sudo timedatectl set-timezone UTC</span>
                    </div>
                  </div>
                </div>
              </div>

              {/* Step 3: Install Docker */}
              <div id="docker" className="mb-12 scroll-mt-20">
                <div className="flex items-center gap-3 mb-4">
                  <div className="flex h-8 w-8 items-center justify-center rounded-md bg-primary/10 border border-primary/15 text-sm font-mono font-semibold text-primary shrink-0">
                    3
                  </div>
                  <h2 className="text-2xl font-semibold tracking-tight">
                    Install Docker
                  </h2>
                </div>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Docker is the only dependency you need on the server. Install Docker Engine and Docker Compose:
                </p>

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
                    <div className="text-muted-foreground/50 text-xs mb-2"># Install Docker using the official convenience script</div>
                    <div>
                      <span className="text-primary/50 select-none">$ </span>
                      <span className="text-foreground/80">curl -fsSL https://get.docker.com | sh</span>
                    </div>
                    <div className="mt-3 text-muted-foreground/50 text-xs mb-2"># Add your user to the docker group (no more sudo for docker)</div>
                    <div>
                      <span className="text-primary/50 select-none">$ </span>
                      <span className="text-foreground/80">sudo usermod -aG docker $USER</span>
                    </div>
                    <div className="mt-3 text-muted-foreground/50 text-xs mb-2"># Apply group changes (or log out and back in)</div>
                    <div>
                      <span className="text-primary/50 select-none">$ </span>
                      <span className="text-foreground/80">newgrp docker</span>
                    </div>
                    <div className="mt-3 text-muted-foreground/50 text-xs mb-2"># Verify installation</div>
                    <div>
                      <span className="text-primary/50 select-none">$ </span>
                      <span className="text-foreground/80">docker --version</span>
                    </div>
                    <div>
                      <span className="text-primary/50 select-none">$ </span>
                      <span className="text-foreground/80">docker compose version</span>
                    </div>
                  </div>
                </div>
              </div>

              {/* Step 4: Clone & Configure */}
              <div id="clone" className="mb-12 scroll-mt-20">
                <div className="flex items-center gap-3 mb-4">
                  <div className="flex h-8 w-8 items-center justify-center rounded-md bg-primary/10 border border-primary/15 text-sm font-mono font-semibold text-primary shrink-0">
                    4
                  </div>
                  <h2 className="text-2xl font-semibold tracking-tight">
                    Clone &amp; Configure
                  </h2>
                </div>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Clone your project to the server and set up production environment variables.
                  We recommend <code>/opt/myapp</code> as the deployment directory.
                </p>

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
                    <div className="text-muted-foreground/50 text-xs mb-2"># Clone your project</div>
                    <div>
                      <span className="text-primary/50 select-none">$ </span>
                      <span className="text-foreground/80">sudo mkdir -p /opt/myapp &amp;&amp; sudo chown deploy:deploy /opt/myapp</span>
                    </div>
                    <div>
                      <span className="text-primary/50 select-none">$ </span>
                      <span className="text-foreground/80">git clone https://github.com/you/myapp.git /opt/myapp</span>
                    </div>
                    <div>
                      <span className="text-primary/50 select-none">$ </span>
                      <span className="text-foreground/80">cd /opt/myapp</span>
                    </div>
                    <div className="mt-3 text-muted-foreground/50 text-xs mb-2"># Create production .env from example</div>
                    <div>
                      <span className="text-primary/50 select-none">$ </span>
                      <span className="text-foreground/80">cp .env.example .env</span>
                    </div>
                    <div>
                      <span className="text-primary/50 select-none">$ </span>
                      <span className="text-foreground/80">nano .env</span>
                    </div>
                  </div>
                </div>

                <h3 className="text-lg font-semibold tracking-tight mb-3">
                  Production .env
                </h3>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Update these critical values for production. Never use default credentials in production:
                </p>
                <CodeBlock language="bash" filename=".env" code={`APP_ENV=production
APP_PORT=8080

# Database — change password!
DATABASE_URL=postgres://grit:YOUR_STRONG_DB_PASSWORD@postgres:5432/myapp?sslmode=disable
POSTGRES_USER=grit
POSTGRES_PASSWORD=YOUR_STRONG_DB_PASSWORD
POSTGRES_DB=myapp

# Auth — generate a random 64-char secret
# openssl rand -hex 32
JWT_SECRET=GENERATE_A_RANDOM_64_CHAR_STRING

# Redis
REDIS_URL=redis://redis:6379

# Frontend URLs (update with your domains)
WEB_URL=https://yourdomain.com
ADMIN_URL=https://admin.yourdomain.com
CORS_ORIGINS=https://yourdomain.com,https://admin.yourdomain.com

# Storage (Cloudflare R2 recommended for production)
STORAGE_DRIVER=r2
R2_ENDPOINT=https://your-account.r2.cloudflarestorage.com
R2_ACCESS_KEY=your-access-key
R2_SECRET_KEY=your-secret-key
R2_BUCKET=myapp-uploads

# Email
RESEND_API_KEY=re_your_api_key
MAIL_FROM=noreply@yourdomain.com

# AI (optional)
AI_PROVIDER=claude
AI_API_KEY=sk-ant-xxxxx
AI_MODEL=claude-sonnet-4-5-20250929

# Disable GORM Studio in production
GORM_STUDIO_ENABLED=false`} />

                <div className="p-4 rounded-lg border border-warning/20 bg-warning/5 mb-6">
                  <p className="text-sm text-foreground/80 leading-relaxed">
                    <span className="font-semibold text-warning">Security tip:</span>{' '}
                    Generate a strong JWT secret with: <code className="text-xs bg-card/50 px-1.5 py-0.5 rounded">openssl rand -hex 32</code>.
                    Never commit your <code>.env</code> file to git.
                  </p>
                </div>

                <h3 className="text-lg font-semibold tracking-tight mb-3">
                  Environment Variables Reference
                </h3>
                <div className="rounded-lg border border-border/30 bg-card/30 overflow-hidden mb-6">
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
                        { var: 'DATABASE_URL', req: 'Yes', note: 'PostgreSQL connection string' },
                        { var: 'JWT_SECRET', req: 'Yes', note: 'Random 64+ character string. Never reuse.' },
                        { var: 'REDIS_URL', req: 'Yes', note: 'Redis connection URL' },
                        { var: 'CORS_ORIGINS', req: 'Yes', note: 'Comma-separated allowed frontend domains' },
                        { var: 'WEB_URL', req: 'Yes', note: 'Your web app URL (for CORS and redirects)' },
                        { var: 'ADMIN_URL', req: 'Yes', note: 'Your admin panel URL' },
                        { var: 'RESEND_API_KEY', req: 'If emails', note: 'Resend API key for transactional emails' },
                        { var: 'STORAGE_DRIVER', req: 'If uploads', note: '"r2", "b2", or "minio" for file storage' },
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

              {/* Step 5: Deploy */}
              <div id="deploy" className="mb-12 scroll-mt-20">
                <div className="flex items-center gap-3 mb-4">
                  <div className="flex h-8 w-8 items-center justify-center rounded-md bg-primary/10 border border-primary/15 text-sm font-mono font-semibold text-primary shrink-0">
                    5
                  </div>
                  <h2 className="text-2xl font-semibold tracking-tight">
                    Deploy with Docker Compose
                  </h2>
                </div>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Grit includes a production-ready <code>docker-compose.prod.yml</code> that builds optimized
                  containers for the Go API and Next.js apps, plus PostgreSQL and Redis.
                </p>

                <div className="rounded-xl border border-border/40 bg-card/80 overflow-hidden glow-purple-sm mb-6">
                  <div className="flex items-center gap-2 px-4 py-2.5 border-b border-border/30 bg-accent/30">
                    <div className="flex items-center gap-1.5">
                      <div className="h-2.5 w-2.5 rounded-full bg-[#ff5f57]" />
                      <div className="h-2.5 w-2.5 rounded-full bg-[#febc2e]" />
                      <div className="h-2.5 w-2.5 rounded-full bg-[#28c840]" />
                    </div>
                    <span className="ml-2 text-[11px] font-mono text-muted-foreground/40">server — /opt/myapp</span>
                  </div>
                  <div className="p-5 font-mono text-sm space-y-2">
                    <div className="text-muted-foreground/50 text-xs mb-2"># Build and start everything</div>
                    <div>
                      <span className="text-primary/50 select-none">$ </span>
                      <span className="text-foreground/80">docker compose -f docker-compose.prod.yml up -d --build</span>
                    </div>
                    <div className="mt-3 text-muted-foreground/50 text-xs mb-2"># Check that all services are running</div>
                    <div>
                      <span className="text-primary/50 select-none">$ </span>
                      <span className="text-foreground/80">docker compose -f docker-compose.prod.yml ps</span>
                    </div>
                    <div className="mt-3 text-muted-foreground/50 text-xs mb-2"># Test the API</div>
                    <div>
                      <span className="text-primary/50 select-none">$ </span>
                      <span className="text-foreground/80">curl http://localhost:8080/api/auth/me</span>
                    </div>
                  </div>
                </div>

                <p className="text-muted-foreground leading-relaxed mb-4">
                  Expected output: you should see all containers running (api, web, admin, postgres, redis).
                  The first build may take a few minutes as Docker builds the Go binary and Next.js apps.
                </p>

                <h3 className="text-lg font-semibold tracking-tight mb-3">
                  What the production build does
                </h3>
                <div className="space-y-2 mb-6">
                  {[
                    { label: 'Go API', desc: 'Multi-stage build: compiles to a single static binary (~15MB), runs in a minimal distroless container' },
                    { label: 'Next.js Web', desc: 'Builds with standalone output mode, creating a minimal Node.js server' },
                    { label: 'Next.js Admin', desc: 'Same as web, built as a separate service on port 3001' },
                    { label: 'PostgreSQL', desc: 'Persistent volume for data, configured with production settings' },
                    { label: 'Redis', desc: 'In-memory cache and job queue, configured with maxmemory and eviction policy' },
                  ].map((item) => (
                    <div key={item.label} className="flex gap-3 text-[13px]">
                      <span className="font-mono text-primary/70 shrink-0 w-24">{item.label}</span>
                      <span className="text-muted-foreground">{item.desc}</span>
                    </div>
                  ))}
                </div>

                <h3 className="text-lg font-semibold tracking-tight mb-3">
                  View logs
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
                    <div className="text-muted-foreground/50 text-xs mb-2"># Follow all logs</div>
                    <div>
                      <span className="text-primary/50 select-none">$ </span>
                      <span className="text-foreground/80">docker compose -f docker-compose.prod.yml logs -f</span>
                    </div>
                    <div className="mt-3 text-muted-foreground/50 text-xs mb-2"># Follow only API logs</div>
                    <div>
                      <span className="text-primary/50 select-none">$ </span>
                      <span className="text-foreground/80">docker compose -f docker-compose.prod.yml logs -f api</span>
                    </div>
                  </div>
                </div>
              </div>

              {/* Step 6: SSL/TLS */}
              <div id="ssl" className="mb-12 scroll-mt-20">
                <div className="flex items-center gap-3 mb-4">
                  <div className="flex h-8 w-8 items-center justify-center rounded-md bg-primary/10 border border-primary/15 text-sm font-mono font-semibold text-primary shrink-0">
                    6
                  </div>
                  <h2 className="text-2xl font-semibold tracking-tight">
                    Set Up SSL with Caddy
                  </h2>
                </div>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Caddy is a modern web server that automatically provisions SSL certificates from
                  Let&apos;s Encrypt. It&apos;s the easiest way to get HTTPS working.
                </p>

                <h3 className="text-lg font-semibold tracking-tight mb-3">
                  6a. Install Caddy
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
                      <span className="text-foreground/80">curl -1sLf &apos;https://dl.cloudsmith.io/public/caddy/stable/debian.deb.txt&apos; | sudo tee /etc/apt/sources.list.d/caddy-stable.list</span>
                    </div>
                    <div>
                      <span className="text-primary/50 select-none">$ </span>
                      <span className="text-foreground/80">sudo apt update &amp;&amp; sudo apt install caddy</span>
                    </div>
                  </div>
                </div>

                <h3 className="text-lg font-semibold tracking-tight mb-3">
                  6b. Configure the Caddyfile
                </h3>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Before editing, make sure your DNS A records point to your server&apos;s IP address.
                  Caddy will automatically obtain SSL certificates once DNS is configured.
                </p>
                <CodeBlock filename="/etc/caddy/Caddyfile" code={`# Go API
api.yourdomain.com {
    reverse_proxy localhost:8080
}

# Next.js Web App
yourdomain.com {
    reverse_proxy localhost:3000
}

# Admin Panel
admin.yourdomain.com {
    reverse_proxy localhost:3001
}`} />

                <h3 className="text-lg font-semibold tracking-tight mb-3">
                  6c. Start Caddy
                </h3>
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
                      <span className="text-foreground/80">sudo nano /etc/caddy/Caddyfile  # paste config above</span>
                    </div>
                    <div>
                      <span className="text-primary/50 select-none">$ </span>
                      <span className="text-foreground/80">sudo systemctl reload caddy</span>
                    </div>
                    <div>
                      <span className="text-primary/50 select-none">$ </span>
                      <span className="text-foreground/80">sudo systemctl enable caddy     # auto-start on boot</span>
                    </div>
                  </div>
                </div>
                <p className="text-sm text-muted-foreground/60 mb-6">
                  Caddy will automatically obtain and renew SSL certificates from Let&apos;s Encrypt.
                  After a few seconds, your site will be live at <code>https://yourdomain.com</code>.
                </p>

                <h3 className="text-lg font-semibold tracking-tight mb-3">
                  DNS Records
                </h3>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Add these DNS A records at your domain registrar or Cloudflare:
                </p>
                <div className="rounded-lg border border-border/30 bg-card/30 overflow-hidden mb-6">
                  <table className="w-full text-sm">
                    <thead>
                      <tr className="border-b border-border/30 bg-accent/20">
                        <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Type</th>
                        <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Name</th>
                        <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Value</th>
                      </tr>
                    </thead>
                    <tbody className="text-muted-foreground">
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-mono text-xs">A</td>
                        <td className="px-4 py-2.5 font-mono text-xs">@</td>
                        <td className="px-4 py-2.5 font-mono text-xs">your-server-ip</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-mono text-xs">A</td>
                        <td className="px-4 py-2.5 font-mono text-xs">api</td>
                        <td className="px-4 py-2.5 font-mono text-xs">your-server-ip</td>
                      </tr>
                      <tr>
                        <td className="px-4 py-2.5 font-mono text-xs">A</td>
                        <td className="px-4 py-2.5 font-mono text-xs">admin</td>
                        <td className="px-4 py-2.5 font-mono text-xs">your-server-ip</td>
                      </tr>
                    </tbody>
                  </table>
                </div>
              </div>

              {/* Step 7: Firewall */}
              <div id="firewall" className="mb-12 scroll-mt-20">
                <div className="flex items-center gap-3 mb-4">
                  <div className="flex h-8 w-8 items-center justify-center rounded-md bg-primary/10 border border-primary/15 text-sm font-mono font-semibold text-primary shrink-0">
                    7
                  </div>
                  <h2 className="text-2xl font-semibold tracking-tight">
                    Configure Firewall (UFW)
                  </h2>
                </div>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Lock down your server to only expose the ports you need. UFW (Uncomplicated Firewall) makes
                  this simple:
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
                    <div className="text-muted-foreground/50 text-xs mb-2"># Allow SSH (don&apos;t lock yourself out!)</div>
                    <div>
                      <span className="text-primary/50 select-none">$ </span>
                      <span className="text-foreground/80">sudo ufw allow OpenSSH</span>
                    </div>
                    <div className="mt-3 text-muted-foreground/50 text-xs mb-2"># Allow HTTP and HTTPS (for Caddy)</div>
                    <div>
                      <span className="text-primary/50 select-none">$ </span>
                      <span className="text-foreground/80">sudo ufw allow 80/tcp</span>
                    </div>
                    <div>
                      <span className="text-primary/50 select-none">$ </span>
                      <span className="text-foreground/80">sudo ufw allow 443/tcp</span>
                    </div>
                    <div className="mt-3 text-muted-foreground/50 text-xs mb-2"># Enable the firewall</div>
                    <div>
                      <span className="text-primary/50 select-none">$ </span>
                      <span className="text-foreground/80">sudo ufw enable</span>
                    </div>
                    <div>
                      <span className="text-primary/50 select-none">$ </span>
                      <span className="text-foreground/80">sudo ufw status</span>
                    </div>
                  </div>
                </div>
                <p className="text-sm text-muted-foreground/60">
                  This blocks all inbound traffic except SSH (22), HTTP (80), and HTTPS (443).
                  Database (5432) and Redis (6379) ports are only accessible from within Docker&apos;s network.
                </p>
              </div>

              {/* Step 8: CI/CD */}
              <div id="cicd" className="mb-12 scroll-mt-20">
                <div className="flex items-center gap-3 mb-4">
                  <div className="flex h-8 w-8 items-center justify-center rounded-md bg-primary/10 border border-primary/15 text-sm font-mono font-semibold text-primary shrink-0">
                    8
                  </div>
                  <h2 className="text-2xl font-semibold tracking-tight">
                    Automated Deployments (CI/CD)
                  </h2>
                </div>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Automate deployments so that every push to <code>main</code> deploys to production.
                  Here&apos;s a GitHub Actions workflow that SSHs into your server and redeploys:
                </p>

                <CodeBlock language="yaml" filename=".github/workflows/deploy.yml" code={`name: Deploy to VPS

on:
  push:
    branches: [main]

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - name: Deploy via SSH
        uses: appleboy/ssh-action@v1
        with:
          host: \${{ secrets.VPS_HOST }}
          username: deploy
          key: \${{ secrets.VPS_SSH_KEY }}
          script: |
            cd /opt/myapp
            git pull origin main
            docker compose -f docker-compose.prod.yml up -d --build
            docker image prune -f`} />

                <h3 className="text-lg font-semibold tracking-tight mb-3">
                  Set up GitHub Secrets
                </h3>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  In your GitHub repo, go to <strong>Settings &rarr; Secrets &rarr; Actions</strong> and add:
                </p>
                <div className="rounded-lg border border-border/30 bg-card/30 overflow-hidden mb-6">
                  <table className="w-full text-sm">
                    <thead>
                      <tr className="border-b border-border/30 bg-accent/20">
                        <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Secret</th>
                        <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Value</th>
                      </tr>
                    </thead>
                    <tbody className="text-muted-foreground">
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-mono text-xs text-primary/80">VPS_HOST</td>
                        <td className="px-4 py-2.5 text-xs">Your server IP (e.g., 164.90.131.42)</td>
                      </tr>
                      <tr>
                        <td className="px-4 py-2.5 font-mono text-xs text-primary/80">VPS_SSH_KEY</td>
                        <td className="px-4 py-2.5 text-xs">The deploy user&apos;s private SSH key (contents of ~/.ssh/id_ed25519)</td>
                      </tr>
                    </tbody>
                  </table>
                </div>

                <div className="p-4 rounded-lg border border-primary/20 bg-primary/5">
                  <p className="text-sm text-foreground/80 leading-relaxed">
                    <span className="font-semibold text-primary/90">Tip:</span>{' '}
                    Generate a dedicated deploy key for CI/CD instead of using your personal SSH key:
                    <code className="text-xs bg-card/50 px-1.5 py-0.5 rounded ml-1">ssh-keygen -t ed25519 -f deploy_key -C &quot;github-actions&quot;</code>.
                    Add the public key to <code>~/.ssh/authorized_keys</code> on the server.
                  </p>
                </div>
              </div>

              {/* Step 9: Database Backup */}
              <div id="backups" className="mb-12 scroll-mt-20">
                <div className="flex items-center gap-3 mb-4">
                  <div className="flex h-8 w-8 items-center justify-center rounded-md bg-primary/10 border border-primary/15 text-sm font-mono font-semibold text-primary shrink-0">
                    9
                  </div>
                  <h2 className="text-2xl font-semibold tracking-tight">
                    Database Backups
                  </h2>
                </div>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Always back up your production database. Losing data is irreversible.
                </p>

                <h3 className="text-lg font-semibold tracking-tight mb-3">
                  Manual backup
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
                    <div className="text-muted-foreground/50 text-xs mb-2"># Create a backup</div>
                    <div>
                      <span className="text-primary/50 select-none">$ </span>
                      <span className="text-foreground/80">docker compose -f docker-compose.prod.yml exec postgres pg_dump -U grit myapp &gt; backup.sql</span>
                    </div>
                    <div className="mt-3 text-muted-foreground/50 text-xs mb-2"># Restore from backup</div>
                    <div>
                      <span className="text-primary/50 select-none">$ </span>
                      <span className="text-foreground/80">docker compose -f docker-compose.prod.yml exec -T postgres psql -U grit myapp &lt; backup.sql</span>
                    </div>
                  </div>
                </div>

                <h3 className="text-lg font-semibold tracking-tight mb-3">
                  Automated daily backups (cron)
                </h3>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Set up a cron job to backup daily and keep the last 7 days:
                </p>
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
                    <div className="text-muted-foreground/50 text-xs mb-2"># Create backup directory</div>
                    <div>
                      <span className="text-primary/50 select-none">$ </span>
                      <span className="text-foreground/80">sudo mkdir -p /backups &amp;&amp; sudo chown deploy:deploy /backups</span>
                    </div>
                    <div className="mt-3 text-muted-foreground/50 text-xs mb-2"># Add to crontab</div>
                    <div>
                      <span className="text-primary/50 select-none">$ </span>
                      <span className="text-foreground/80">crontab -e</span>
                    </div>
                  </div>
                </div>
                <CodeBlock filename="crontab" code={`# Daily DB backup at 3 AM, keep last 7 days
0 3 * * * cd /opt/myapp && docker compose -f docker-compose.prod.yml exec -T postgres pg_dump -U grit myapp | gzip > /backups/myapp-$(date +\\%Y\\%m\\%d).sql.gz && find /backups -name "myapp-*.sql.gz" -mtime +7 -delete`} />

                <h3 className="text-lg font-semibold tracking-tight mb-3">
                  Managed database backups
                </h3>
                <p className="text-muted-foreground leading-relaxed">
                  If you use a managed PostgreSQL service (DigitalOcean Managed Databases, Supabase, Neon, etc.),
                  backups are handled automatically. This is the recommended approach for important production workloads.
                </p>
              </div>

              {/* Step 10: Monitoring */}
              <div id="monitoring" className="mb-12 scroll-mt-20">
                <div className="flex items-center gap-3 mb-4">
                  <div className="flex h-8 w-8 items-center justify-center rounded-md bg-primary/10 border border-primary/15 text-sm font-mono font-semibold text-primary shrink-0">
                    10
                  </div>
                  <h2 className="text-2xl font-semibold tracking-tight">
                    Monitoring &amp; Logging
                  </h2>
                </div>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Grit&apos;s Go API includes a built-in request logger middleware. For production, add proper
                  monitoring to catch issues before your users do:
                </p>
                <div className="space-y-3">
                  {[
                    { name: 'Docker Logs', desc: 'Built-in. Use "docker compose logs -f api" to stream API logs in real time. Combine with journald for persistence.', tier: 'Free' },
                    { name: 'Uptime Kuma', desc: 'Self-hosted uptime monitoring. Set up health checks against your API endpoint. Get notified via email, Slack, or Telegram.', tier: 'Free (self-hosted)' },
                    { name: 'Grafana + Prometheus', desc: 'For advanced metrics. Add a /metrics endpoint to your Go API and visualize request rates, latency, and error rates.', tier: 'Free (self-hosted)' },
                    { name: 'Sentry', desc: 'Error tracking for both Go API and Next.js frontends. Captures stack traces, context, and user info.', tier: 'Free tier available' },
                    { name: 'Better Stack', desc: 'Managed log aggregation and uptime monitoring. Stream Docker logs for search, alerting, and dashboards.', tier: 'Free tier available' },
                  ].map((item) => (
                    <div key={item.name} className="p-4 rounded-lg border border-border/30 bg-card/30 flex justify-between items-start gap-4">
                      <div>
                        <h3 className="text-sm font-semibold mb-1.5">{item.name}</h3>
                        <p className="text-xs text-muted-foreground/70 leading-relaxed">{item.desc}</p>
                      </div>
                      <span className="shrink-0 rounded-full bg-accent/10 px-2.5 py-0.5 text-[10px] font-medium text-primary/70">{item.tier}</span>
                    </div>
                  ))}
                </div>
              </div>

              {/* Step 11: Zero-Downtime Updates */}
              <div id="updates" className="mb-12 scroll-mt-20">
                <div className="flex items-center gap-3 mb-4">
                  <div className="flex h-8 w-8 items-center justify-center rounded-md bg-primary/10 border border-primary/15 text-sm font-mono font-semibold text-primary shrink-0">
                    11
                  </div>
                  <h2 className="text-2xl font-semibold tracking-tight">
                    Zero-Downtime Updates
                  </h2>
                </div>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  To update your running application with minimal downtime, use Docker&apos;s rolling update capability:
                </p>

                <div className="rounded-xl border border-border/40 bg-card/80 overflow-hidden mb-6">
                  <div className="flex items-center gap-2 px-4 py-2.5 border-b border-border/30 bg-accent/30">
                    <div className="flex items-center gap-1.5">
                      <div className="h-2.5 w-2.5 rounded-full bg-[#ff5f57]" />
                      <div className="h-2.5 w-2.5 rounded-full bg-[#febc2e]" />
                      <div className="h-2.5 w-2.5 rounded-full bg-[#28c840]" />
                    </div>
                    <span className="ml-2 text-[11px] font-mono text-muted-foreground/40">server — /opt/myapp</span>
                  </div>
                  <div className="p-5 font-mono text-sm space-y-2">
                    <div className="text-muted-foreground/50 text-xs mb-2"># Pull latest code</div>
                    <div>
                      <span className="text-primary/50 select-none">$ </span>
                      <span className="text-foreground/80">git pull origin main</span>
                    </div>
                    <div className="mt-3 text-muted-foreground/50 text-xs mb-2"># Rebuild and restart only changed services</div>
                    <div>
                      <span className="text-primary/50 select-none">$ </span>
                      <span className="text-foreground/80">docker compose -f docker-compose.prod.yml up -d --build</span>
                    </div>
                    <div className="mt-3 text-muted-foreground/50 text-xs mb-2"># Clean up old images to save disk space</div>
                    <div>
                      <span className="text-primary/50 select-none">$ </span>
                      <span className="text-foreground/80">docker image prune -f</span>
                    </div>
                  </div>
                </div>

                <p className="text-muted-foreground leading-relaxed mb-4">
                  Docker Compose detects which services have changed and only rebuilds and restarts those.
                  PostgreSQL and Redis keep running throughout. The typical downtime for an API update is under 5 seconds.
                </p>

                <div className="p-4 rounded-lg border border-primary/20 bg-primary/5">
                  <p className="text-sm text-foreground/80 leading-relaxed">
                    <span className="font-semibold text-primary/90">For true zero-downtime:</span>{' '}
                    Use Docker Compose&apos;s health checks combined with Caddy. Caddy will only route traffic to
                    healthy containers, so the old container serves requests until the new one is ready.
                  </p>
                </div>
              </div>

              {/* Step 12: Checklist */}
              <div id="checklist" className="mb-12 scroll-mt-20">
                <div className="flex items-center gap-3 mb-4">
                  <div className="flex h-8 w-8 items-center justify-center rounded-md bg-primary/10 border border-primary/15 text-sm font-mono font-semibold text-primary shrink-0">
                    12
                  </div>
                  <h2 className="text-2xl font-semibold tracking-tight">
                    Production Checklist
                  </h2>
                </div>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Before going live, make sure you have completed these steps:
                </p>
                <div className="space-y-2.5">
                  {[
                    { category: 'Security', items: [
                      'Generate a strong, unique JWT_SECRET (64+ characters)',
                      'Set APP_ENV=production',
                      'Set GORM_STUDIO_ENABLED=false',
                      'Use sslmode=require if using external database',
                      'Restrict CORS_ORIGINS to your actual domains',
                      'Change default database credentials from grit/grit',
                      'Enable SSL/TLS via Caddy',
                      'Configure UFW firewall (SSH + HTTP/HTTPS only)',
                      'Disable root SSH login',
                    ]},
                    { category: 'Infrastructure', items: [
                      'Set up automated daily database backups',
                      'Configure uptime monitoring (Uptime Kuma or Better Stack)',
                      'Set up error tracking (Sentry or similar)',
                      'Docker restart policies set to "unless-stopped"',
                      'DNS A records pointing to your server',
                      'SSL certificates working (check with browser)',
                      'Set up CI/CD for automated deployments',
                    ]},
                    { category: 'Performance', items: [
                      'Set appropriate database connection pool limits',
                      'Add indexes to frequently queried columns',
                      'Enable Redis caching for read-heavy endpoints',
                      'Use a CDN for static assets (Cloudflare recommended)',
                      'Test the app under load with hey or k6',
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

              {/* Cloud Platforms */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Alternative: Cloud Platforms
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  If you prefer managed infrastructure over a VPS, these platforms support Grit&apos;s architecture.
                  Deploy the Go API and Next.js apps as separate services:
                </p>
                <div className="space-y-4">
                  {[
                    {
                      name: 'Railway',
                      desc: 'Deploy from GitHub with zero config. Supports Go, Node.js, PostgreSQL, and Redis as managed services.',
                      fit: 'Best for solo developers and small teams.',
                    },
                    {
                      name: 'Render',
                      desc: 'Deploy Go as a Web Service, Next.js as a Web Service. Managed PostgreSQL and Redis available.',
                      fit: 'Best for side projects and startups.',
                    },
                    {
                      name: 'Fly.io',
                      desc: 'Container-based deployment with edge computing. Deploy Docker images directly. Built-in Postgres and Redis.',
                      fit: 'Best for apps needing global edge distribution.',
                    },
                    {
                      name: 'Coolify (Self-Hosted)',
                      desc: 'Open-source Heroku alternative you host on your own VPS. Docker-based, built-in databases, auto SSL. Free.',
                      fit: 'Best for PaaS convenience at VPS pricing.',
                    },
                    {
                      name: 'Dokploy (Self-Hosted)',
                      desc: 'Open-source PaaS with Docker Compose support, web dashboard, auto SSL via Traefik, and GitHub auto-deploy. Free.',
                      fit: 'Best for Docker Compose-based monorepos like Grit. See our dedicated guide.',
                      link: '/docs/infrastructure/dokploy',
                    },
                  ].map((platform: { name: string; desc: string; fit: string; link?: string }) => (
                    <div key={platform.name} className="p-4 rounded-lg border border-border/30 bg-card/30">
                      <h3 className="text-sm font-semibold mb-1.5">{platform.name}</h3>
                      <p className="text-xs text-muted-foreground/70 leading-relaxed mb-2">{platform.desc}</p>
                      <p className="text-xs text-primary/70 leading-relaxed">
                        {platform.fit}
                        {platform.link && (
                          <>
                            {' '}<Link href={platform.link} className="text-primary hover:underline font-medium">Read the guide &rarr;</Link>
                          </>
                        )}
                      </p>
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
