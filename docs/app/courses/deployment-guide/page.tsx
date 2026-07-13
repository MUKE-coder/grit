import Link from 'next/link'
import { CheckCircle2 } from 'lucide-react'
import { SiteHeader } from '@/components/site-header'
import { GridFrame } from '@/components/grid-frame'
import { CodeBlock, Challenge, Note, Tip, Definition, Code, CourseNav, CourseFooter } from '@/components/course-components'
import type { Metadata } from 'next'

export const metadata: Metadata = {
  title: 'Deployment Guide: Dokploy, Orbita, VPS & Vercel — Grit Course',
  description: 'Learn four deployment methods for Grit applications: grit deploy (direct VPS), Dokploy (self-hosted PaaS), Orbita (Grit-native deployment), and Vercel (for Next.js frontends).',
}

export default function DeploymentGuideCourse() {
  return (
    <div className="relative min-h-screen bg-background isolate">
      <SiteHeader />
      <GridFrame />

      <main className="max-w-3xl mx-auto px-6 py-12">
        {/* Breadcrumb */}
        <div className="flex items-center gap-2 text-sm text-muted-foreground mb-8">
          <Link href="/courses" className="hover:text-foreground transition-colors">Courses</Link>
          <span>/</span>
          <span className="text-foreground">Deployment Guide</span>
        </div>

        {/* Header */}
        <div className="mb-10">
          <div className="flex items-center gap-2 mb-3">
            <span className="text-xs font-mono font-medium text-primary bg-primary/10 px-2 py-0.5 rounded">Standalone Course</span>
            <span className="text-xs text-muted-foreground">~30 min</span>
            <span className="text-xs text-muted-foreground">•</span>
            <span className="text-xs text-muted-foreground">14 challenges</span>
          </div>
          <h1 className="text-3xl md:text-4xl font-bold text-foreground mb-4">
            Deployment Guide: Dokploy, Orbita, VPS & Vercel
          </h1>
          <p className="text-lg text-muted-foreground leading-relaxed">
            Your app works locally. Now you need to put it on the internet. In this course, you will learn
            four deployment methods — from the simplest one-command deploy to full PaaS setups — so you
            can choose the right approach for your project.
          </p>
        </div>

        <div className="my-4 rounded-lg border border-primary/20 bg-primary/5 px-4 py-3">
          <p className="text-sm text-muted-foreground">
            <strong className="text-foreground">Reference docs:</strong>{' '}
            <a href="/docs/infrastructure/deployment" className="text-primary hover:underline">Deployment →</a>
          </p>
        </div>

        <hr className="border-border/40 mb-10" />

        {/* ═══ Section 1: Deployment Options ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Deployment Options Overview</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Grit apps can be deployed in multiple ways depending on your needs, budget, and technical
            comfort level. Here are the four methods you will learn:
          </p>

          <div className="overflow-x-auto mb-4">
            <table className="w-full text-sm border border-border/40 rounded-lg overflow-hidden">
              <thead>
                <tr className="bg-muted/20">
                  <th className="text-left px-4 py-3 font-semibold text-foreground border-b border-border/40">Method</th>
                  <th className="text-left px-4 py-3 font-semibold text-foreground border-b border-border/40">What It Is</th>
                  <th className="text-left px-4 py-3 font-semibold text-foreground border-b border-border/40">Best For</th>
                </tr>
              </thead>
              <tbody className="text-muted-foreground">
                <tr className="border-b border-border/20">
                  <td className="px-4 py-3 font-medium text-foreground">grit deploy</td>
                  <td className="px-4 py-3">Direct VPS deployment via SCP</td>
                  <td className="px-4 py-3">Quick deploys, solo developers</td>
                </tr>
                <tr className="border-b border-border/20">
                  <td className="px-4 py-3 font-medium text-foreground">Dokploy</td>
                  <td className="px-4 py-3">Self-hosted PaaS with web UI</td>
                  <td className="px-4 py-3">Teams, auto-deploy on push</td>
                </tr>
                <tr className="border-b border-border/20">
                  <td className="px-4 py-3 font-medium text-foreground">Orbita</td>
                  <td className="px-4 py-3">Grit-native deployment platform</td>
                  <td className="px-4 py-3">Grit projects, seamless integration</td>
                </tr>
                <tr>
                  <td className="px-4 py-3 font-medium text-foreground">Vercel</td>
                  <td className="px-4 py-3">Managed hosting for Next.js</td>
                  <td className="px-4 py-3">Next.js frontend (API on VPS)</td>
                </tr>
              </tbody>
            </table>
          </div>

          <Note>
            All four methods assume you have a VPS (Virtual Private Server) for the Go API. You can get
            one for $4-6/month from providers like Hetzner, DigitalOcean, or Linode. Vercel is the only
            method where part of your app runs on managed infrastructure.
          </Note>

          <Challenge number={1} title="Choose Your Factors">
            <p>What factors would you consider when choosing a deployment method? Think about: cost, complexity, automation, team size, frequency of deploys, and whether you want a web UI to manage your server.</p>
          </Challenge>
        </section>

        {/* ═══ Section 2: Method 1 — grit deploy ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Method 1: grit deploy (Direct VPS)</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            The simplest deployment method. One command builds your app, uploads it to your server,
            and configures everything automatically:
          </p>

          <CodeBlock filename="Terminal">
{`grit deploy --host deploy@server.com --domain myapp.com`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Here{"'"}s what <Code>grit deploy</Code> does behind the scenes:
          </p>

          <ul className="space-y-2 text-muted-foreground mb-4">
            <li className="flex gap-2"><span className="text-primary">1.</span> <strong className="text-foreground">Builds the Go binary</strong> — cross-compiles for Linux (your server{"'"}s OS)</li>
            <li className="flex gap-2"><span className="text-primary">2.</span> <strong className="text-foreground">Uploads via SCP</strong> — securely copies the binary to your server</li>
            <li className="flex gap-2"><span className="text-primary">3.</span> <strong className="text-foreground">Creates a systemd service</strong> — so your app auto-restarts on crash or reboot</li>
            <li className="flex gap-2"><span className="text-primary">4.</span> <strong className="text-foreground">Configures Caddy</strong> — reverse proxy with automatic HTTPS (Lets Encrypt)</li>
            <li className="flex gap-2"><span className="text-primary">5.</span> <strong className="text-foreground">Starts the service</strong> — your app is live at your domain</li>
          </ul>

          <Definition term="systemd Service">
            A Linux background process manager. When you register your app as a systemd service, the
            operating system ensures it stays running — restarting it automatically if it crashes or
            when the server reboots. Think of it as a {'"'}keep alive{'"'} system for your app.
          </Definition>

          <Definition term="Caddy">
            A modern web server and reverse proxy that automatically manages HTTPS certificates via
            Lets Encrypt. It sits in front of your Go app, handles SSL termination, and forwards
            requests to your app on its internal port. Zero SSL configuration needed.
          </Definition>

          <Tip>
            Before your first deploy, make sure you can SSH into your server: <Code>ssh deploy@server.com</Code>.
            If that works, <Code>grit deploy</Code> will work too.
          </Tip>

          <Challenge number={2} title="List the Deploy Steps">
            <p>Without looking above, list the 5 steps that <Code>grit deploy</Code> performs. For each step, explain why it{"'"}s necessary.</p>
          </Challenge>
        </section>

        {/* ═══ Section 3: Method 2 — Dokploy ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Method 2: Dokploy</h2>

          <Definition term="Dokploy">
            A self-hosted Platform as a Service (PaaS) you install on your own VPS. It gives you a
            Heroku-like experience — connect a GitHub repo, push code, and Dokploy builds and deploys
            automatically. It runs Docker containers behind the scenes and provides a web dashboard for
            managing your apps, databases, domains, and SSL certificates. Free and open-source.
            Website: <Code>https://dokploy.com</Code>
          </Definition>

          <p className="text-muted-foreground leading-relaxed mb-4">
            While <Code>grit deploy</Code> is a one-time push, Dokploy gives you a full deployment pipeline:
          </p>

          <ul className="space-y-2 text-muted-foreground mb-4">
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Git push to deploy</strong> — Dokploy watches your repo and auto-deploys on push</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Web dashboard</strong> — manage apps, view logs, set environment variables from a browser</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Docker-based</strong> — each app runs in an isolated container</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Auto-SSL</strong> — HTTPS certificates managed automatically</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Database management</strong> — spin up PostgreSQL, Redis, and more from the UI</li>
          </ul>

          <p className="text-muted-foreground leading-relaxed mb-4">
            To set up Dokploy on your VPS:
          </p>

          <CodeBlock filename="Terminal (on your VPS)">
{`# Install Dokploy (one command)
curl -sSL https://dokploy.com/install.sh | sh`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            After installation, access the Dokploy dashboard at <Code>https://your-server-ip:3000</Code>.
            From there, connect your GitHub account, select your repository, and Dokploy handles the rest.
          </p>

          <Challenge number={3} title="Explore Dokploy">
            <p>Visit <Code>dokploy.com</Code>. What features does it offer? How does it compare to Heroku or Railway? What databases can it manage?</p>
          </Challenge>
        </section>

        {/* ═══ Section 4: Dokploy Configuration ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Dokploy Configuration</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Dokploy uses Docker to build and run your app. Grit scaffolds a production-ready Dockerfile
            that uses multi-stage builds to keep the final image small:
          </p>

          <CodeBlock filename="Dockerfile">
{`# Stage 1: Build the Go binary
FROM golang:1.24 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server

# Stage 2: Create a minimal runtime image
FROM alpine:latest
COPY --from=builder /app/server /server
EXPOSE 8080
CMD ["/server"]`}
          </CodeBlock>

          <Definition term="Multi-Stage Build">
            A Docker technique that uses multiple <Code>FROM</Code> statements in a single Dockerfile.
            The first stage (builder) has all the build tools — Go compiler, source code, dependencies.
            The final stage only contains the compiled binary. This dramatically reduces the image size:
            the builder stage might be 1 GB, but the final image is only 15-20 MB.
          </Definition>

          <p className="text-muted-foreground leading-relaxed mb-4">
            For the full stack (API + database + Redis), use the production Docker Compose file:
          </p>

          <CodeBlock filename="docker-compose.prod.yml">
{`version: '3.8'
services:
  api:
    build: .
    ports:
      - "8080:8080"
    env_file: .env
    depends_on:
      - db
      - redis
  db:
    image: postgres:16-alpine
    environment:
      POSTGRES_DB: myapp
      POSTGRES_USER: myapp
      POSTGRES_PASSWORD: \${DB_PASSWORD}
    volumes:
      - pgdata:/var/lib/postgresql/data
  redis:
    image: redis:7-alpine
    volumes:
      - redisdata:/data
volumes:
  pgdata:
  redisdata:`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            In Dokploy, you can configure environment variables through the web UI instead of managing
            <Code>.env</Code> files on the server. Set your database credentials, JWT secret, API keys,
            and other sensitive values securely through the dashboard.
          </p>

          <Challenge number={4} title="Read the Dockerfile">
            <p>Look at the Dockerfile above. What{"'"}s the difference between the builder stage and the final stage? Why does the final stage use <Code>alpine:latest</Code> instead of <Code>golang:1.24</Code>? What would happen to the image size if you used only one stage?</p>
          </Challenge>
        </section>

        {/* ═══ Section 5: Method 3 — Orbita ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Method 3: Orbita</h2>

          <Definition term="Orbita">
            A deployment and server management platform created by the same developer who built Grit
            (MUKE-coder). It{"'"}s designed to work seamlessly with Grit projects — understanding the
            project structure, environment variables, and deployment patterns out of the box. Open-source
            and self-hosted. GitHub: <Code>https://github.com/MUKE-coder/orbita</Code>
          </Definition>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Orbita is purpose-built for Grit{"'"}s ecosystem. While Dokploy is a general-purpose PaaS,
            Orbita understands Grit conventions:
          </p>

          <ul className="space-y-2 text-muted-foreground mb-4">
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">One-click deploy</strong> — connect your server, select a Grit project, deploy</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Auto-SSL</strong> — HTTPS certificates managed automatically via Lets Encrypt</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Environment management</strong> — configure production .env through the UI</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Log viewing</strong> — stream application logs from the dashboard</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Server management</strong> — manage multiple servers and apps from one place</li>
          </ul>

          <p className="text-muted-foreground leading-relaxed mb-4">
            The deployment flow with Orbita:
          </p>

          <ul className="space-y-2 text-muted-foreground mb-4">
            <li className="flex gap-2"><span className="text-primary">1.</span> Install Orbita on your VPS (or use a separate management server)</li>
            <li className="flex gap-2"><span className="text-primary">2.</span> Add your server{"'"}s SSH credentials in the Orbita dashboard</li>
            <li className="flex gap-2"><span className="text-primary">3.</span> Connect your GitHub repository</li>
            <li className="flex gap-2"><span className="text-primary">4.</span> Configure environment variables and domain</li>
            <li className="flex gap-2"><span className="text-primary">5.</span> Click Deploy — Orbita builds and ships your app</li>
          </ul>

          <Challenge number={5} title="Explore Orbita">
            <p>Visit the Orbita GitHub repo at <Code>https://github.com/MUKE-coder/orbita</Code>. What deployment features does it support? How does the setup process compare to Dokploy?</p>
          </Challenge>
        </section>

        {/* ═══ Section 6: Orbita vs Dokploy ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Orbita vs Dokploy</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Both are self-hosted PaaS tools, but they serve slightly different purposes:
          </p>

          <div className="overflow-x-auto mb-4">
            <table className="w-full text-sm border border-border/40 rounded-lg overflow-hidden">
              <thead>
                <tr className="bg-muted/20">
                  <th className="text-left px-4 py-3 font-semibold text-foreground border-b border-border/40">Feature</th>
                  <th className="text-left px-4 py-3 font-semibold text-foreground border-b border-border/40">Orbita</th>
                  <th className="text-left px-4 py-3 font-semibold text-foreground border-b border-border/40">Dokploy</th>
                </tr>
              </thead>
              <tbody className="text-muted-foreground">
                <tr className="border-b border-border/20">
                  <td className="px-4 py-3 font-medium text-foreground">Grit integration</td>
                  <td className="px-4 py-3">Native — built for Grit</td>
                  <td className="px-4 py-3">Generic — works with any Docker app</td>
                </tr>
                <tr className="border-b border-border/20">
                  <td className="px-4 py-3 font-medium text-foreground">Ecosystem</td>
                  <td className="px-4 py-3">Same creator as Grit</td>
                  <td className="px-4 py-3">Independent project</td>
                </tr>
                <tr className="border-b border-border/20">
                  <td className="px-4 py-3 font-medium text-foreground">Use case</td>
                  <td className="px-4 py-3">Grit-first deployment</td>
                  <td className="px-4 py-3">Any Docker-based app</td>
                </tr>
                <tr className="border-b border-border/20">
                  <td className="px-4 py-3 font-medium text-foreground">Web UI</td>
                  <td className="px-4 py-3">Yes</td>
                  <td className="px-4 py-3">Yes</td>
                </tr>
                <tr className="border-b border-border/20">
                  <td className="px-4 py-3 font-medium text-foreground">Auto-SSL</td>
                  <td className="px-4 py-3">Yes</td>
                  <td className="px-4 py-3">Yes</td>
                </tr>
                <tr>
                  <td className="px-4 py-3 font-medium text-foreground">Database management</td>
                  <td className="px-4 py-3">Yes</td>
                  <td className="px-4 py-3">Yes (more options)</td>
                </tr>
              </tbody>
            </table>
          </div>

          <Tip>
            If you{"'"}re deploying only Grit apps and want the tightest integration, use Orbita. If you{"'"}re
            managing multiple projects (some non-Grit) on the same server, Dokploy{"'"}s broader compatibility
            may be more useful.
          </Tip>

          <Challenge number={6} title="Compare the Two">
            <p>You have 3 apps to deploy: a Grit SaaS, a Python Flask API, and a static marketing site. Which deployment tool (Orbita or Dokploy) would you use for each, and why?</p>
          </Challenge>
        </section>

        {/* ═══ Section 7: Method 4 — Vercel ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Method 4: Vercel (Frontend Only)</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            For Triple or Double architecture projects with a Next.js frontend, you can deploy the
            frontend to Vercel while keeping the Go API on your VPS. This gives you the best of both
            worlds — Vercel{"'"}s CDN and edge network for the frontend, your own server for the API.
          </p>

          <CodeBlock filename="Terminal">
{`# Step 1: Deploy the Go API to your VPS
grit deploy --host deploy@server.com --domain api.myapp.com

# Step 2: Deploy the Next.js frontend to Vercel
cd apps/web && vercel`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            The critical configuration step: tell the Vercel-hosted frontend where to find your API.
            In Vercel{"'"}s dashboard, set this environment variable:
          </p>

          <CodeBlock filename="Vercel Environment Variables">
{`NEXT_PUBLIC_API_URL=https://api.myapp.com`}
          </CodeBlock>

          <Note>
            This split deployment only works with Next.js frontends. If you{"'"}re using TanStack Router
            (Vite), the frontend is embedded in the Go binary via <Code>go:embed</Code> and served by
            the Go server directly. There{"'"}s nothing to deploy to Vercel.
          </Note>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Advantages of the Vercel + VPS split:
          </p>

          <ul className="space-y-2 text-muted-foreground mb-4">
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Global CDN</strong> — frontend served from edge nodes closest to the user</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Automatic previews</strong> — every PR gets a preview deployment URL</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Zero-config</strong> — Vercel auto-detects Next.js and configures everything</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Free tier</strong> — generous free plan for personal projects</li>
          </ul>

          <Challenge number={7} title="The Connection Point">
            <p>What environment variable connects the Vercel frontend to your API? Why does the variable name start with <Code>NEXT_PUBLIC_</Code>? What would happen if you forgot to set it?</p>
          </Challenge>
        </section>

        {/* ═══ Section 8: DNS Configuration ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">DNS Configuration</h2>

          <Definition term="DNS (Domain Name System)">
            The system that translates human-readable domain names (like myapp.com) into IP addresses
            (like 167.99.100.50) that computers use to find each other on the internet. When you type
            a URL in your browser, DNS looks up the IP address for that domain and connects you to the
            right server.
          </Definition>

          <p className="text-muted-foreground leading-relaxed mb-4">
            After buying a domain, you need to point it at your server by creating DNS records in your
            domain registrar{"'"}s dashboard (Namecheap, Cloudflare, GoDaddy, etc.):
          </p>

          <div className="overflow-x-auto mb-4">
            <table className="w-full text-sm border border-border/40 rounded-lg overflow-hidden">
              <thead>
                <tr className="bg-muted/20">
                  <th className="text-left px-4 py-3 font-semibold text-foreground border-b border-border/40">Type</th>
                  <th className="text-left px-4 py-3 font-semibold text-foreground border-b border-border/40">Name</th>
                  <th className="text-left px-4 py-3 font-semibold text-foreground border-b border-border/40">Value</th>
                  <th className="text-left px-4 py-3 font-semibold text-foreground border-b border-border/40">Purpose</th>
                </tr>
              </thead>
              <tbody className="text-muted-foreground">
                <tr className="border-b border-border/20">
                  <td className="px-4 py-3 font-medium text-foreground">A</td>
                  <td className="px-4 py-3">@</td>
                  <td className="px-4 py-3">167.99.100.50</td>
                  <td className="px-4 py-3">Root domain (myapp.com)</td>
                </tr>
                <tr className="border-b border-border/20">
                  <td className="px-4 py-3 font-medium text-foreground">A</td>
                  <td className="px-4 py-3">api</td>
                  <td className="px-4 py-3">167.99.100.50</td>
                  <td className="px-4 py-3">API subdomain (api.myapp.com)</td>
                </tr>
                <tr>
                  <td className="px-4 py-3 font-medium text-foreground">CNAME</td>
                  <td className="px-4 py-3">www</td>
                  <td className="px-4 py-3">myapp.com</td>
                  <td className="px-4 py-3">www subdomain redirects to root</td>
                </tr>
              </tbody>
            </table>
          </div>

          <Note>
            DNS changes can take up to 48 hours to propagate worldwide, but usually take 5-30 minutes.
            You can check propagation status at <Code>dnschecker.org</Code>.
          </Note>

          <Challenge number={8} title="Plan Your DNS Records">
            <p>If your server IP is <Code>167.99.100.50</Code> and your domain is <Code>myapp.com</Code>, what DNS records would you create for: (1) the main site at myapp.com, (2) the API at api.myapp.com, (3) www.myapp.com redirecting to myapp.com?</p>
          </Challenge>
        </section>

        {/* ═══ Section 9: Production Environment Variables ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Production Environment Variables</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Your development <Code>.env</Code> file has placeholder values that are not safe for production.
            Before deploying, you need to update several critical variables:
          </p>

          <CodeBlock filename=".env (production)">
{`# App
APP_ENV=production
APP_PORT=8080

# Database
DB_DRIVER=postgres
DB_HOST=localhost
DB_PORT=5432
DB_USER=myapp
DB_PASSWORD=<strong-random-password>
DB_NAME=myapp_production

# Auth
JWT_SECRET=<64-character-random-string>

# Storage
STORAGE_DRIVER=s3
S3_BUCKET=myapp-uploads
S3_REGION=us-east-1
S3_ACCESS_KEY=<your-access-key>
S3_SECRET_KEY=<your-secret-key>

# Email
RESEND_API_KEY=re_<your-real-key>
EMAIL_FROM=noreply@myapp.com

# Redis
REDIS_URL=redis://localhost:6379`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Critical changes from development to production:
          </p>

          <ul className="space-y-2 text-muted-foreground mb-4">
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">APP_ENV</strong> — set to <Code>production</Code> (disables debug logging, enables security headers)</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">JWT_SECRET</strong> — must be a long, random string (not {'"'}secret{'"'} or {'"'}changeme{'"'})</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">DB_PASSWORD</strong> — a strong, unique password (not {'"'}password{'"'})</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">STORAGE_DRIVER</strong> — switch from local to S3, R2, or B2 for persistent file storage</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">RESEND_API_KEY</strong> — your real Resend API key (development uses Mailhog)</li>
          </ul>

          <Tip>
            Generate strong random strings for JWT_SECRET and passwords using: <Code>openssl rand -hex 32</Code>
          </Tip>

          <Challenge number={9} title="Write a Production .env">
            <p>Write a complete production <Code>.env</Code> file for your app. Use placeholder values (not real secrets), but make sure every required variable is present. How many variables changed from your development .env?</p>
          </Challenge>
        </section>

        {/* ═══ Section 10: SSL/TLS ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">SSL/TLS Certificates</h2>

          <Definition term="SSL/TLS Certificate">
            A digital certificate that encrypts the connection between a user{"'"}s browser and your server.
            It{"'"}s what makes the padlock icon appear in the browser{"'"}s address bar and enables HTTPS
            (instead of HTTP). Without SSL, data (including passwords) is transmitted in plain text and
            can be intercepted.
          </Definition>

          <p className="text-muted-foreground leading-relaxed mb-4">
            The good news: you almost never need to manage SSL certificates manually. Each deployment
            method handles it automatically:
          </p>

          <ul className="space-y-2 text-muted-foreground mb-4">
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">grit deploy</strong> — Caddy auto-provisions Lets Encrypt certificates</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Dokploy</strong> — built-in Lets Encrypt support through the dashboard</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Orbita</strong> — automatic SSL certificate management</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Vercel</strong> — SSL included automatically for all deployments</li>
          </ul>

          <Definition term="Lets Encrypt">
            A free, automated certificate authority that provides SSL/TLS certificates at no cost. Before
            Lets Encrypt (launched 2015), SSL certificates cost $50-200/year. Now, tools like Caddy
            and Dokploy automatically request, install, and renew certificates from Lets Encrypt
            with zero configuration.
          </Definition>

          <Challenge number={10} title="Check Your SSL">
            <p>After deploying your app, visit <Code>https://www.ssllabs.com/ssltest/</Code> and enter your domain. What grade does your site receive? An A or A+ is expected with Caddy{"'"}s default configuration.</p>
          </Challenge>
        </section>

        {/* ═══ Section 11: Monitoring After Deployment ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Monitoring After Deployment</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Deploying your app is only half the job. You also need to know when something goes wrong.
            There are two types of monitoring:
          </p>

          <Definition term="Application Monitoring">
            Internal monitoring that tracks your app{"'"}s performance and errors from the inside. Metrics like
            response times, error rates, database query speed, and memory usage. Grit{"'"}s Pulse feature
            provides this — it{"'"}s a built-in dashboard that shows your app{"'"}s health.
          </Definition>

          <Definition term="Uptime Monitoring">
            External monitoring that checks whether your app is reachable from the outside. A service
            pings your app{"'"}s URL every few minutes and alerts you (via email, SMS, or Slack) if it{"'"}s
            down. Think of it as a friend who checks your website every minute and calls you if it{"'"}s broken.
          </Definition>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Recommended monitoring setup:
          </p>

          <ul className="space-y-2 text-muted-foreground mb-4">
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Pulse</strong> — built into Grit, shows response times, error rates, and active connections</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">UptimeRobot</strong> (free) — pings your health endpoint every 5 minutes, alerts on downtime</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Better Stack</strong> (free tier) — uptime monitoring + incident management + status pages</li>
          </ul>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Every Grit app has a health endpoint at <Code>/api/health</Code> that returns the app{"'"}s
            status. Point your uptime monitor at this endpoint.
          </p>

          <Challenge number={11} title="Set Up Uptime Monitoring">
            <p>Create a free account at <Code>uptimerobot.com</Code>. Add a monitor for your deployed API{"'"}s health endpoint (<Code>https://api.myapp.com/api/health</Code>). Set the check interval to 5 minutes. What alerts did you configure?</p>
          </Challenge>
        </section>

        {/* ═══ Section 12: Summary + Final Challenge ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">What You Learned</h2>

          <ul className="space-y-2 text-muted-foreground mb-6">
            <li className="flex gap-2"><CheckCircle2 className="h-4 w-4 text-green-500 shrink-0 mt-1" /> Four deployment methods: grit deploy, Dokploy, Orbita, and Vercel</li>
            <li className="flex gap-2"><CheckCircle2 className="h-4 w-4 text-green-500 shrink-0 mt-1" /> How <Code>grit deploy</Code> builds, uploads, and configures your app in one command</li>
            <li className="flex gap-2"><CheckCircle2 className="h-4 w-4 text-green-500 shrink-0 mt-1" /> Dokploy as a self-hosted PaaS with Docker and auto-deploy</li>
            <li className="flex gap-2"><CheckCircle2 className="h-4 w-4 text-green-500 shrink-0 mt-1" /> Orbita as a Grit-native deployment platform</li>
            <li className="flex gap-2"><CheckCircle2 className="h-4 w-4 text-green-500 shrink-0 mt-1" /> Splitting frontend (Vercel) and API (VPS) deployment</li>
            <li className="flex gap-2"><CheckCircle2 className="h-4 w-4 text-green-500 shrink-0 mt-1" /> DNS configuration for domains and subdomains</li>
            <li className="flex gap-2"><CheckCircle2 className="h-4 w-4 text-green-500 shrink-0 mt-1" /> Production environment variables and security</li>
            <li className="flex gap-2"><CheckCircle2 className="h-4 w-4 text-green-500 shrink-0 mt-1" /> Automatic SSL with Lets Encrypt</li>
            <li className="flex gap-2"><CheckCircle2 className="h-4 w-4 text-green-500 shrink-0 mt-1" /> Monitoring with Pulse and uptime services</li>
          </ul>

          <Challenge number={12} title="Deploy End-to-End">
            <p>Choose ONE deployment method and deploy a Grit app from start to finish:</p>
            <ol className="list-decimal list-inside space-y-1 mt-2">
              <li>Scaffold a project with <Code>grit new deploy-test --double --next</Code></li>
              <li>Generate at least one resource</li>
              <li>Deploy using your chosen method</li>
              <li>Configure DNS for your domain</li>
              <li>Verify HTTPS works</li>
              <li>Set up uptime monitoring</li>
              <li>Document every step you took</li>
            </ol>
          </Challenge>

          <Challenge number={13} title="Compare Two Methods">
            <p>Deploy the same app using two different methods (e.g., <Code>grit deploy</Code> and Dokploy). Compare the experience: setup time, ease of re-deployment, log access, and environment variable management. Which do you prefer and why?</p>
          </Challenge>

          <Challenge number={14} title="Production Checklist">
            <p>Create a deployment checklist for your team. Include: environment variables to change, DNS records to create, SSL verification, uptime monitoring setup, and a rollback plan. Would this checklist change depending on the deployment method?</p>
          </Challenge>
        </section>

        {/* Footer */}
        <CourseFooter />

        {/* Navigation */}
        <div className="mt-8">
          <CourseNav
            prev={{ href: '/courses', label: 'All Courses' }}
            next={{ href: '/courses', label: 'More Courses' }}
          />
        </div>
      </main>
    </div>
  )
}
