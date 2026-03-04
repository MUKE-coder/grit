import Link from 'next/link'
import { ArrowRight, Cloud, Database, Mail, HardDrive, Shield } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'
import { CodeBlock } from '@/components/code-block'
import { getDocMetadata } from '@/config/docs-metadata'

export const metadata = getDocMetadata('/docs/getting-started/create-without-docker')

export default function CreateWithoutDockerPage() {
  return (
    <div className="min-h-screen bg-background">
      <SiteHeader />
      <DocsSidebar />

      <main className="lg:pl-64">
        <div className="container max-w-screen-xl py-10 px-6">
          <div className="max-w-3xl">
            {/* Header */}
            <div className="mb-10">
              <span className="tag-mono text-primary/80 mb-3 block">Getting Started</span>
              <h1 className="text-4xl font-bold tracking-tight mb-4">
                Create without Docker
              </h1>
              <p className="text-lg text-muted-foreground leading-relaxed">
                Set up a Grit project using cloud services instead of Docker. Perfect for
                machines where Docker isn&apos;t available or when you prefer managed services.
              </p>
            </div>

            {/* Prerequisites */}
            <div className="prose-grit mb-8">
              <h2>Prerequisites</h2>
              <p>
                You don&apos;t need Docker at all. Just make sure you have these installed:
              </p>
              <ul>
                <li><strong>Go 1.21+</strong> &mdash; <code>go version</code></li>
                <li><strong>Node.js 18+</strong> &mdash; <code>node --version</code></li>
                <li><strong>pnpm</strong> &mdash; <code>npm install -g pnpm</code></li>
                <li><strong>Grit CLI</strong> &mdash; <code>go install github.com/MUKE-coder/grit/v2/cmd/grit@latest</code></li>
              </ul>
            </div>

            {/* Step 1: Create project */}
            <div className="prose-grit mb-8">
              <h2>Step 1: Create your project</h2>
              <p>
                Scaffold a new Grit project as usual:
              </p>
            </div>
            <CodeBlock terminal code={`grit new myapp`} className="mb-4" />
            <div className="prose-grit mb-8">
              <p>
                The CLI will generate all the files, including <code>.env</code>,{' '}
                <code>.env.cloud.example</code>, and <code>docker-compose.yml</code>. You&apos;ll
                ignore the Docker files and configure cloud services instead.
              </p>
            </div>

            {/* Step 2: Cloud services */}
            <div className="prose-grit mb-6">
              <h2>Step 2: Set up cloud services</h2>
              <p>
                Instead of running PostgreSQL, Redis, and MinIO in Docker containers, you&apos;ll
                use managed cloud services. All of these have generous free tiers.
              </p>
            </div>

            {/* Neon */}
            <div className="mb-6">
              <div className="rounded-xl border border-border/40 bg-card/30 overflow-hidden">
                <div className="px-5 py-3 border-b border-border/30 bg-accent/30 flex items-center gap-2">
                  <Database className="h-4 w-4 text-primary" />
                  <span className="text-sm font-semibold">PostgreSQL &mdash; Neon</span>
                  <a href="https://neon.tech" target="_blank" rel="noreferrer" className="ml-auto text-xs text-primary hover:underline">neon.tech &rarr;</a>
                </div>
                <div className="p-5 space-y-3">
                  <div className="prose-grit">
                    <ol>
                      <li>Create a free account at <a href="https://neon.tech" target="_blank" rel="noreferrer">neon.tech</a></li>
                      <li>Create a new project (pick a region close to you)</li>
                      <li>Copy the <strong>connection string</strong> from the dashboard</li>
                    </ol>
                    <p>It looks like this:</p>
                  </div>
                  <CodeBlock language="bash" code="DATABASE_URL=postgres://user:password@ep-xxx-xxx-123456.us-east-2.aws.neon.tech/neondb?sslmode=require" className="mb-0" />
                </div>
              </div>
            </div>

            {/* Upstash */}
            <div className="mb-6">
              <div className="rounded-xl border border-border/40 bg-card/30 overflow-hidden">
                <div className="px-5 py-3 border-b border-border/30 bg-accent/30 flex items-center gap-2">
                  <Shield className="h-4 w-4 text-red-400" />
                  <span className="text-sm font-semibold">Redis &mdash; Upstash</span>
                  <a href="https://upstash.com" target="_blank" rel="noreferrer" className="ml-auto text-xs text-primary hover:underline">upstash.com &rarr;</a>
                </div>
                <div className="p-5 space-y-3">
                  <div className="prose-grit">
                    <ol>
                      <li>Create a free account at <a href="https://upstash.com" target="_blank" rel="noreferrer">upstash.com</a></li>
                      <li>Create a new Redis database</li>
                      <li>Copy the <strong>Redis URL</strong> (TLS enabled)</li>
                    </ol>
                    <p>It looks like this:</p>
                  </div>
                  <CodeBlock language="bash" code="REDIS_URL=rediss://default:your-password@your-endpoint.upstash.io:6379" className="mb-0" />
                  <div className="prose-grit">
                    <p className="text-sm text-muted-foreground">
                      Note the <code>rediss://</code> (with double &apos;s&apos;) — that&apos;s TLS-encrypted Redis.
                    </p>
                  </div>
                </div>
              </div>
            </div>

            {/* Cloudflare R2 */}
            <div className="mb-6">
              <div className="rounded-xl border border-border/40 bg-card/30 overflow-hidden">
                <div className="px-5 py-3 border-b border-border/30 bg-accent/30 flex items-center gap-2">
                  <HardDrive className="h-4 w-4 text-orange-400" />
                  <span className="text-sm font-semibold">File Storage &mdash; Cloudflare R2</span>
                  <a href="https://dash.cloudflare.com" target="_blank" rel="noreferrer" className="ml-auto text-xs text-primary hover:underline">cloudflare.com &rarr;</a>
                </div>
                <div className="p-5 space-y-3">
                  <div className="prose-grit">
                    <ol>
                      <li>Go to <a href="https://dash.cloudflare.com" target="_blank" rel="noreferrer">Cloudflare Dashboard</a> &rarr; R2</li>
                      <li>Create a bucket (e.g. <code>myapp-uploads</code>)</li>
                      <li>Go to <strong>Manage R2 API Tokens</strong> &rarr; Create Token</li>
                      <li>Copy the endpoint, access key, and secret key</li>
                    </ol>
                  </div>
                  <CodeBlock language="bash" code={`STORAGE_DRIVER=r2
R2_ENDPOINT=https://your-account-id.r2.cloudflarestorage.com
R2_ACCESS_KEY=your-r2-access-key-id
R2_SECRET_KEY=your-r2-secret-access-key
R2_BUCKET=myapp-uploads
R2_REGION=auto`} className="mb-0" />
                  <div className="prose-grit">
                    <p className="text-sm text-muted-foreground">
                      Alternatively, you can use <strong>Backblaze B2</strong> — set <code>STORAGE_DRIVER=b2</code> and
                      configure the B2 variables instead.
                    </p>
                  </div>
                </div>
              </div>
            </div>

            {/* Resend */}
            <div className="mb-8">
              <div className="rounded-xl border border-border/40 bg-card/30 overflow-hidden">
                <div className="px-5 py-3 border-b border-border/30 bg-accent/30 flex items-center gap-2">
                  <Mail className="h-4 w-4 text-blue-400" />
                  <span className="text-sm font-semibold">Email &mdash; Resend</span>
                  <a href="https://resend.com" target="_blank" rel="noreferrer" className="ml-auto text-xs text-primary hover:underline">resend.com &rarr;</a>
                </div>
                <div className="p-5 space-y-3">
                  <div className="prose-grit">
                    <ol>
                      <li>Sign up at <a href="https://resend.com" target="_blank" rel="noreferrer">resend.com</a></li>
                      <li>Verify your domain (or use the sandbox for testing)</li>
                      <li>Create an API key</li>
                    </ol>
                  </div>
                  <CodeBlock language="bash" code={`RESEND_API_KEY=re_your_api_key_here
MAIL_FROM=noreply@yourdomain.com`} className="mb-0" />
                </div>
              </div>
            </div>

            {/* Step 3: Configure .env */}
            <div className="prose-grit mb-6">
              <h2>Step 3: Configure your .env</h2>
              <p>
                Grit scaffolds a <code>.env.cloud.example</code> file with all the cloud variables
                pre-filled as placeholders. Copy it and fill in your keys:
              </p>
            </div>
            <CodeBlock terminal code={`cd myapp
cp .env.cloud.example .env`} className="mb-4" />
            <div className="prose-grit mb-4">
              <p>
                Then open <code>.env</code> and replace the placeholder values with your actual
                cloud service credentials. Here&apos;s the full template:
              </p>
            </div>
            <CodeBlock language="bash" filename=".env" code={`# App
APP_NAME=myapp
APP_ENV=development
APP_PORT=8080
APP_URL=http://localhost:8080

# Database (Neon)
DATABASE_URL=postgres://user:password@ep-xxx.neon.tech/neondb?sslmode=require

# JWT
JWT_SECRET=change-me-to-a-random-string-at-least-32-chars
JWT_ACCESS_EXPIRY=15m
JWT_REFRESH_EXPIRY=168h

# Redis (Upstash)
REDIS_URL=rediss://default:your-password@your-endpoint.upstash.io:6379

# Storage (Cloudflare R2)
STORAGE_DRIVER=r2
R2_ENDPOINT=https://your-account-id.r2.cloudflarestorage.com
R2_ACCESS_KEY=your-r2-access-key-id
R2_SECRET_KEY=your-r2-secret-access-key
R2_BUCKET=myapp-uploads
R2_REGION=auto

# Email (Resend)
RESEND_API_KEY=re_your_api_key_here
MAIL_FROM=noreply@yourdomain.com

# CORS
CORS_ORIGINS=http://localhost:3000,http://localhost:3001

# GORM Studio
GORM_STUDIO_ENABLED=true
GORM_STUDIO_USERNAME=admin
GORM_STUDIO_PASSWORD=change-me-in-prod`} className="mb-8" />

            {/* Step 4: Run */}
            <div className="prose-grit mb-6">
              <h2>Step 4: Start developing</h2>
              <p>
                With cloud services configured, you can run everything locally without Docker:
              </p>
            </div>

            <div className="space-y-4 mb-8">
              <div className="rounded-lg border border-border/30 bg-card/30 p-4">
                <h4 className="text-sm font-semibold mb-2">1. Install dependencies</h4>
                <CodeBlock terminal code={`cd myapp
pnpm install
cd apps/api && go mod tidy && cd ../..`} className="mb-0" />
              </div>

              <div className="rounded-lg border border-border/30 bg-card/30 p-4">
                <h4 className="text-sm font-semibold mb-2">2. Start the Go API</h4>
                <CodeBlock terminal code={`cd apps/api
go run cmd/api/main.go`} className="mb-0" />
                <p className="text-xs text-muted-foreground mt-2">
                  Runs on <code>http://localhost:8080</code>. Auto-migrates the database on first run.
                </p>
              </div>

              <div className="rounded-lg border border-border/30 bg-card/30 p-4">
                <h4 className="text-sm font-semibold mb-2">3. Start the frontend apps (in separate terminals)</h4>
                <CodeBlock terminal code={`# Terminal 2 — Web app
pnpm --filter web dev

# Terminal 3 — Admin panel
pnpm --filter admin dev`} className="mb-0" />
                <p className="text-xs text-muted-foreground mt-2">
                  Web on <code>http://localhost:3000</code>, Admin on <code>http://localhost:3001</code>
                </p>
              </div>
            </div>

            {/* Comparison */}
            <div className="prose-grit mb-6">
              <h2>Docker vs Cloud: when to use which</h2>
            </div>
            <div className="mb-8">
              <div className="rounded-xl border border-border/40 overflow-hidden">
                <table className="w-full text-sm">
                  <thead>
                    <tr className="border-b border-border/30 bg-accent/30">
                      <th className="text-left px-4 py-2.5 font-semibold"></th>
                      <th className="text-left px-4 py-2.5 font-semibold">Docker (default)</th>
                      <th className="text-left px-4 py-2.5 font-semibold">Cloud services</th>
                    </tr>
                  </thead>
                  <tbody className="divide-y divide-border/20">
                    <tr>
                      <td className="px-4 py-2.5 font-medium text-muted-foreground">Setup</td>
                      <td className="px-4 py-2.5"><code>docker compose up -d</code></td>
                      <td className="px-4 py-2.5">Sign up for 3-4 services</td>
                    </tr>
                    <tr>
                      <td className="px-4 py-2.5 font-medium text-muted-foreground">Cost</td>
                      <td className="px-4 py-2.5">Free (local)</td>
                      <td className="px-4 py-2.5">Free tiers available</td>
                    </tr>
                    <tr>
                      <td className="px-4 py-2.5 font-medium text-muted-foreground">Offline</td>
                      <td className="px-4 py-2.5">Works offline</td>
                      <td className="px-4 py-2.5">Needs internet</td>
                    </tr>
                    <tr>
                      <td className="px-4 py-2.5 font-medium text-muted-foreground">Resources</td>
                      <td className="px-4 py-2.5">Uses local RAM/CPU</td>
                      <td className="px-4 py-2.5">Lightweight locally</td>
                    </tr>
                    <tr>
                      <td className="px-4 py-2.5 font-medium text-muted-foreground">Best for</td>
                      <td className="px-4 py-2.5">Most developers</td>
                      <td className="px-4 py-2.5">Low-spec machines, CI/CD</td>
                    </tr>
                  </tbody>
                </table>
              </div>
            </div>

            {/* Tip */}
            <div className="mb-8">
              <div className="rounded-xl border border-primary/30 bg-primary/5 p-6">
                <div className="flex items-start gap-3">
                  <Cloud className="h-5 w-5 text-primary mt-0.5 shrink-0" />
                  <div className="prose-grit">
                    <p className="font-semibold mb-1">Pro tip: use cloud for production too</p>
                    <p>
                      If you set up cloud services for development, you already have your production
                      infrastructure ready. Just create separate Neon projects and R2 buckets for
                      staging and production, and you&apos;re good to deploy.
                    </p>
                  </div>
                </div>
              </div>
            </div>

            {/* Nav */}
            <div className="flex flex-wrap gap-3 mt-12 pt-6 border-t border-border/30">
              <Button variant="outline" asChild className="border-border/60 bg-transparent hover:bg-accent/50">
                <Link href="/docs/getting-started/troubleshooting">
                  Troubleshooting
                </Link>
              </Button>
              <Button asChild className="glow-purple-sm ml-auto">
                <Link href="/docs/concepts/architecture">
                  Architecture
                  <ArrowRight className="ml-1.5 h-4 w-4" />
                </Link>
              </Button>
            </div>
          </div>
        </div>
      </main>
    </div>
  )
}
