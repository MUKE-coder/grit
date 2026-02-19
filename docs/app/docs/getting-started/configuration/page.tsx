import Link from 'next/link'
import { ArrowRight } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'
import { CodeBlock } from '@/components/code-block'

export default function ConfigurationPage() {
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
                Configuration
              </h1>
              <p className="text-lg text-muted-foreground leading-relaxed">
                A complete reference for every environment variable in your Grit project. All
                configuration is done through the <code className="text-sm font-mono bg-accent/80 px-1.5 py-0.5 rounded text-primary">.env</code> file
                at the project root.
              </p>
            </div>

            <div className="prose-grit">
              <h2>Environment Files</h2>
              <p>
                Every Grit project includes three environment files:
              </p>
              <ul>
                <li><strong><code>.env</code></strong> -- Your actual configuration. This file is gitignored and never committed.</li>
                <li><strong><code>.env.example</code></strong> -- Documented template with all variables and sensible defaults. Committed to git.</li>
                <li><strong><code>.env.cloud.example</code></strong> -- Template for cloud-only setups (Neon, Upstash, R2) when you do not have Docker.</li>
              </ul>
              <p>
                The Go API loads these variables at startup using <code>godotenv</code> and parses
                them into a typed <code>Config</code> struct. Environment variables are read
                once at startup and are available throughout the application.
              </p>
            </div>

            {/* App Config */}
            <div className="prose-grit">
              <h2>Application</h2>
              <p>
                General application settings that control the server behavior.
              </p>
            </div>

            <div className="mb-10">
              <CodeBlock language="bash" filename=".env" code={`# App — General application settings
APP_NAME=myapp              # Application name (used in emails, logs)
APP_ENV=development         # Environment: development, staging, production
APP_PORT=8080               # API server port
APP_URL=http://localhost:8080`} />
              <div className="mt-4 space-y-3">
                {[
                  { variable: 'APP_NAME', default: 'Project name', desc: 'Used as the application title in email templates, log entries, and the admin panel header. Set to your project name during scaffolding.' },
                  { variable: 'APP_ENV', default: 'development', desc: 'Controls logging verbosity, GORM Studio visibility, and error detail level. Set to production in deployed environments to disable debug features.' },
                  { variable: 'APP_PORT', default: '8080', desc: 'The port the Go API server listens on. The frontend apps proxy API requests to this port.' },
                  { variable: 'APP_URL', default: 'http://localhost:8080', desc: 'The full URL of the API server. Used for generating absolute URLs in emails and file storage signed URLs.' },
                ].map((item) => (
                  <div key={item.variable} className="rounded-lg border border-border/30 bg-card/30 px-4 py-3">
                    <div className="flex items-center gap-2 mb-1">
                      <code className="text-sm font-mono text-primary/70 font-medium">{item.variable}</code>
                      <span className="text-[10px] font-mono text-muted-foreground/40 bg-accent/50 px-1.5 py-0.5 rounded">{item.default}</span>
                    </div>
                    <p className="text-xs text-muted-foreground/60 leading-relaxed">{item.desc}</p>
                  </div>
                ))}
              </div>
            </div>

            {/* Database Config */}
            <div className="prose-grit">
              <h2>Database</h2>
              <p>
                Grit uses PostgreSQL as its primary database, connected through GORM. The
                connection string follows the standard PostgreSQL URI format.
              </p>
            </div>

            <div className="mb-10">
              <CodeBlock language="bash" filename=".env" code={`# Database — PostgreSQL connection
DATABASE_URL=postgres://grit:grit@localhost:5432/myapp?sslmode=disable`} />
              <div className="mt-4 space-y-3">
                <div className="rounded-lg border border-border/30 bg-card/30 px-4 py-3">
                  <div className="flex items-center gap-2 mb-1">
                    <code className="text-sm font-mono text-primary/70 font-medium">DATABASE_URL</code>
                  </div>
                  <p className="text-xs text-muted-foreground/60 leading-relaxed mb-2">
                    PostgreSQL connection string. Format: <code className="text-xs font-mono text-primary/50">postgres://user:password@host:port/database?sslmode=disable</code>
                  </p>
                  <p className="text-xs text-muted-foreground/60 leading-relaxed">
                    For local development with Docker, the default credentials are <code className="text-xs font-mono text-primary/50">grit:grit</code> on port 5432. For cloud databases like Neon, use the connection string provided by your provider and set <code className="text-xs font-mono text-primary/50">sslmode=require</code>.
                  </p>
                </div>
              </div>
            </div>

            {/* JWT Config */}
            <div className="prose-grit">
              <h2>JWT Authentication</h2>
              <p>
                Grit uses JWT tokens for authentication with separate access and refresh tokens.
                The access token is short-lived (15 minutes) and the refresh token lasts 7 days.
              </p>
            </div>

            <div className="mb-10">
              <CodeBlock language="bash" filename=".env" code={`# JWT — Authentication tokens
JWT_SECRET=your-super-secret-jwt-key-change-in-production
JWT_ACCESS_EXPIRY=15m        # Access token lifetime
JWT_REFRESH_EXPIRY=168h      # Refresh token lifetime (7 days)`} />
              <div className="mt-4 space-y-3">
                {[
                  { variable: 'JWT_SECRET', default: 'Random string', desc: 'The secret key used to sign and verify JWT tokens. MUST be changed in production. Use a random string of at least 32 characters. Both access and refresh tokens use this same secret.' },
                  { variable: 'JWT_ACCESS_EXPIRY', default: '15m', desc: 'How long access tokens are valid. Uses Go duration format: 15m (15 minutes), 1h (1 hour), etc. Keep short for security. The frontend automatically refreshes expired tokens.' },
                  { variable: 'JWT_REFRESH_EXPIRY', default: '168h', desc: 'How long refresh tokens are valid. 168h is 7 days. Users must re-login after this period. Refresh tokens are single-use and rotated on each refresh.' },
                ].map((item) => (
                  <div key={item.variable} className="rounded-lg border border-border/30 bg-card/30 px-4 py-3">
                    <div className="flex items-center gap-2 mb-1">
                      <code className="text-sm font-mono text-primary/70 font-medium">{item.variable}</code>
                      <span className="text-[10px] font-mono text-muted-foreground/40 bg-accent/50 px-1.5 py-0.5 rounded">{item.default}</span>
                    </div>
                    <p className="text-xs text-muted-foreground/60 leading-relaxed">{item.desc}</p>
                  </div>
                ))}
              </div>
            </div>

            {/* Redis Config */}
            <div className="prose-grit">
              <h2>Redis</h2>
              <p>
                Redis is used for response caching and background job queues (via Asynq). A single
                Redis instance handles both use cases.
              </p>
            </div>

            <div className="mb-10">
              <CodeBlock language="bash" filename=".env" code={`# Redis — Cache and job queue
REDIS_URL=redis://localhost:6379`} />
              <div className="mt-4 space-y-3">
                <div className="rounded-lg border border-border/30 bg-card/30 px-4 py-3">
                  <div className="flex items-center gap-2 mb-1">
                    <code className="text-sm font-mono text-primary/70 font-medium">REDIS_URL</code>
                    <span className="text-[10px] font-mono text-muted-foreground/40 bg-accent/50 px-1.5 py-0.5 rounded">redis://localhost:6379</span>
                  </div>
                  <p className="text-xs text-muted-foreground/60 leading-relaxed">
                    Redis connection URL. For local Docker, uses port 6379 with no authentication. For cloud Redis (Upstash), use the <code className="text-xs font-mono text-primary/50">rediss://</code> protocol (with double s) and include the password: <code className="text-xs font-mono text-primary/50">rediss://default:password@endpoint:6379</code>.
                  </p>
                </div>
              </div>
            </div>

            {/* Storage Config */}
            <div className="prose-grit">
              <h2>File Storage</h2>
              <p>
                Grit supports three S3-compatible storage providers: MinIO (local development),
                Cloudflare R2, and Backblaze B2. The <code>STORAGE_DRIVER</code> variable controls
                which provider is active.
              </p>
            </div>

            <div className="mb-10">
              <CodeBlock language="bash" filename=".env" code={`# Storage — Active driver: minio, r2, or b2
STORAGE_DRIVER=minio

# MinIO — Local S3-compatible storage (default for development)
MINIO_ENDPOINT=http://localhost:9000
MINIO_ACCESS_KEY=minioadmin
MINIO_SECRET_KEY=minioadmin
MINIO_BUCKET=myapp-uploads
MINIO_REGION=us-east-1
MINIO_USE_SSL=false

# Cloudflare R2 — S3-compatible object storage
R2_ENDPOINT=https://YOUR_ACCOUNT_ID.r2.cloudflarestorage.com
R2_ACCESS_KEY=               # R2 Access Key ID
R2_SECRET_KEY=               # R2 Secret Access Key
R2_BUCKET=myapp-uploads
R2_REGION=auto               # Always "auto" for R2

# Backblaze B2 — S3-compatible object storage
B2_ENDPOINT=https://s3.us-west-004.backblazeb2.com
B2_ACCESS_KEY=               # B2 keyID
B2_SECRET_KEY=               # B2 applicationKey
B2_BUCKET=myapp-uploads
B2_REGION=us-west-004        # Must match your bucket region`} />
              <div className="mt-4 space-y-3">
                {[
                  { variable: 'STORAGE_DRIVER', default: 'minio', desc: 'Which storage provider to use. Options: minio (local dev with Docker), r2 (Cloudflare R2), b2 (Backblaze B2). Only the variables for the active driver need to be set.' },
                  { variable: 'MINIO_ENDPOINT', default: 'http://localhost:9000', desc: 'MinIO server URL. The Docker Compose file starts MinIO on port 9000 with a web console on port 9001.' },
                  { variable: 'MINIO_ACCESS_KEY / MINIO_SECRET_KEY', default: 'minioadmin', desc: 'Default MinIO credentials. These match the Docker Compose configuration. Change in production if running your own MinIO instance.' },
                  { variable: 'R2_ENDPOINT', default: '(your account)', desc: 'Your Cloudflare R2 endpoint. Get this from the Cloudflare Dashboard under R2 > Overview. Format: https://ACCOUNT_ID.r2.cloudflarestorage.com' },
                  { variable: 'R2_ACCESS_KEY / R2_SECRET_KEY', default: '(your keys)', desc: 'Cloudflare R2 API credentials. Create an API token in the Cloudflare Dashboard under R2 > Manage R2 API Tokens.' },
                  { variable: 'B2_ENDPOINT', default: '(your region)', desc: 'Backblaze B2 S3-compatible endpoint. Format depends on your bucket region. Find it in your B2 bucket settings.' },
                ].map((item) => (
                  <div key={item.variable} className="rounded-lg border border-border/30 bg-card/30 px-4 py-3">
                    <div className="flex items-center gap-2 mb-1">
                      <code className="text-sm font-mono text-primary/70 font-medium">{item.variable}</code>
                      <span className="text-[10px] font-mono text-muted-foreground/40 bg-accent/50 px-1.5 py-0.5 rounded">{item.default}</span>
                    </div>
                    <p className="text-xs text-muted-foreground/60 leading-relaxed">{item.desc}</p>
                  </div>
                ))}
              </div>
            </div>

            {/* Email Config */}
            <div className="prose-grit">
              <h2>Email</h2>
              <p>
                Grit uses <a href="https://resend.com" target="_blank" rel="noreferrer">Resend</a> for
                transactional emails -- welcome emails, password resets, and notifications. In
                development, emails are caught by Mailhog (accessible at <code>http://localhost:8025</code>).
              </p>
            </div>

            <div className="mb-10">
              <CodeBlock language="bash" filename=".env" code={`# Email — Resend integration
RESEND_API_KEY=re_your_api_key
MAIL_FROM=noreply@myapp.dev`} />
              <div className="mt-4 space-y-3">
                {[
                  { variable: 'RESEND_API_KEY', default: 're_your_api_key', desc: 'Your Resend API key. Get one at resend.com/api-keys. In development, emails are sent to Mailhog regardless of this key.' },
                  { variable: 'MAIL_FROM', default: 'noreply@myapp.dev', desc: 'The sender email address for all outgoing emails. Must be a verified domain in your Resend account for production use.' },
                ].map((item) => (
                  <div key={item.variable} className="rounded-lg border border-border/30 bg-card/30 px-4 py-3">
                    <div className="flex items-center gap-2 mb-1">
                      <code className="text-sm font-mono text-primary/70 font-medium">{item.variable}</code>
                      <span className="text-[10px] font-mono text-muted-foreground/40 bg-accent/50 px-1.5 py-0.5 rounded">{item.default}</span>
                    </div>
                    <p className="text-xs text-muted-foreground/60 leading-relaxed">{item.desc}</p>
                  </div>
                ))}
              </div>
            </div>

            {/* CORS Config */}
            <div className="prose-grit">
              <h2>CORS</h2>
              <p>
                Cross-Origin Resource Sharing configuration. The Go API needs to know which
                frontend origins are allowed to make requests.
              </p>
            </div>

            <div className="mb-10">
              <CodeBlock language="bash" filename=".env" code={`# CORS — Allowed frontend origins (comma-separated)
CORS_ORIGINS=http://localhost:3000,http://localhost:3001`} />
              <div className="mt-4 space-y-3">
                <div className="rounded-lg border border-border/30 bg-card/30 px-4 py-3">
                  <div className="flex items-center gap-2 mb-1">
                    <code className="text-sm font-mono text-primary/70 font-medium">CORS_ORIGINS</code>
                    <span className="text-[10px] font-mono text-muted-foreground/40 bg-accent/50 px-1.5 py-0.5 rounded">http://localhost:3000,http://localhost:3001</span>
                  </div>
                  <p className="text-xs text-muted-foreground/60 leading-relaxed">
                    Comma-separated list of allowed origins. In development, port 3000 is the web app and port 3001 is the admin panel. In production, set this to your actual domain names (e.g., <code className="text-xs font-mono text-primary/50">https://myapp.com,https://admin.myapp.com</code>).
                  </p>
                </div>
              </div>
            </div>

            {/* GORM Studio Config */}
            <div className="prose-grit">
              <h2>GORM Studio</h2>
              <p>
                The embedded visual database browser. Accessible at <code>/studio</code> on the
                API server.
              </p>
            </div>

            <div className="mb-10">
              <CodeBlock language="bash" filename=".env" code={`# GORM Studio — Visual database browser
GORM_STUDIO_ENABLED=true`} />
              <div className="mt-4 space-y-3">
                <div className="rounded-lg border border-border/30 bg-card/30 px-4 py-3">
                  <div className="flex items-center gap-2 mb-1">
                    <code className="text-sm font-mono text-primary/70 font-medium">GORM_STUDIO_ENABLED</code>
                    <span className="text-[10px] font-mono text-muted-foreground/40 bg-accent/50 px-1.5 py-0.5 rounded">true</span>
                  </div>
                  <p className="text-xs text-muted-foreground/60 leading-relaxed">
                    Enable or disable GORM Studio. Set to <code className="text-xs font-mono text-primary/50">true</code> in development for visual database browsing. Set to <code className="text-xs font-mono text-primary/50">false</code> in production to disable the database browser and prevent unauthorized access to your data.
                  </p>
                </div>
              </div>
            </div>

            {/* AI Config */}
            <div className="prose-grit">
              <h2>AI Integration</h2>
              <p>
                Grit ships with built-in AI support for both Anthropic Claude and OpenAI. The
                AI service provides text completion and streaming endpoints.
              </p>
            </div>

            <div className="mb-10">
              <CodeBlock language="bash" filename=".env" code={`# AI — Text generation (Claude or OpenAI)
AI_PROVIDER=claude           # "claude" or "openai"
AI_API_KEY=                  # Your API key (sk-ant-... or sk-...)
AI_MODEL=claude-sonnet-4-5-20250929`} />
              <div className="mt-4 space-y-3">
                {[
                  { variable: 'AI_PROVIDER', default: 'claude', desc: 'Which AI provider to use. Options: claude (Anthropic) or openai. The AI service uses the appropriate API format based on this setting.' },
                  { variable: 'AI_API_KEY', default: '(your key)', desc: 'Your AI provider API key. For Anthropic, starts with sk-ant-. For OpenAI, starts with sk-. Leave empty if you do not use AI features.' },
                  { variable: 'AI_MODEL', default: 'claude-sonnet-4-5-20250929', desc: 'The model to use for completions. For Claude: claude-sonnet-4-5-20250929, claude-opus-4-6, etc. For OpenAI: gpt-4o, gpt-4o-mini, etc.' },
                ].map((item) => (
                  <div key={item.variable} className="rounded-lg border border-border/30 bg-card/30 px-4 py-3">
                    <div className="flex items-center gap-2 mb-1">
                      <code className="text-sm font-mono text-primary/70 font-medium">{item.variable}</code>
                      <span className="text-[10px] font-mono text-muted-foreground/40 bg-accent/50 px-1.5 py-0.5 rounded">{item.default}</span>
                    </div>
                    <p className="text-xs text-muted-foreground/60 leading-relaxed">{item.desc}</p>
                  </div>
                ))}
              </div>
            </div>

            {/* Frontend URLs */}
            <div className="prose-grit">
              <h2>Frontend URLs</h2>
              <p>
                The frontend applications run on their own ports in development. These are
                configured in each app&apos;s <code>package.json</code> and referenced in the
                CORS configuration.
              </p>
            </div>

            <div className="mb-10">
              <table className="w-full border-collapse">
                <thead>
                  <tr className="border-b border-border">
                    <th className="text-left text-sm font-semibold text-foreground pb-3 pr-4">Application</th>
                    <th className="text-left text-sm font-semibold text-foreground pb-3 pr-4">Dev URL</th>
                    <th className="text-left text-sm font-semibold text-foreground pb-3">Port</th>
                  </tr>
                </thead>
                <tbody>
                  {[
                    { app: 'Go API', url: 'http://localhost:8080', port: '8080' },
                    { app: 'GORM Studio', url: 'http://localhost:8080/studio', port: '8080' },
                    { app: 'Web App (Next.js)', url: 'http://localhost:3000', port: '3000' },
                    { app: 'Admin Panel (Next.js)', url: 'http://localhost:3001', port: '3001' },
                    { app: 'PostgreSQL', url: 'localhost:5432', port: '5432' },
                    { app: 'Redis', url: 'localhost:6379', port: '6379' },
                    { app: 'MinIO', url: 'http://localhost:9000', port: '9000' },
                    { app: 'MinIO Console', url: 'http://localhost:9001', port: '9001' },
                    { app: 'Mailhog', url: 'http://localhost:8025', port: '8025' },
                  ].map((row) => (
                    <tr key={row.app} className="border-b border-border/50">
                      <td className="text-sm text-foreground py-2.5 pr-4">{row.app}</td>
                      <td className="text-sm text-muted-foreground py-2.5 pr-4 font-mono text-xs text-primary/60">{row.url}</td>
                      <td className="text-sm text-muted-foreground py-2.5 font-mono">{row.port}</td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>

            {/* Complete Reference */}
            <div className="prose-grit">
              <h2>Complete .env Reference</h2>
              <p>
                Here is every environment variable in a single block for quick copy-paste:
              </p>
            </div>

            <div className="mb-10">
              <CodeBlock language="bash" filename=".env.example" code={`# App
APP_NAME=myapp
APP_ENV=development
APP_PORT=8080
APP_URL=http://localhost:8080

# Database
DATABASE_URL=postgres://grit:grit@localhost:5432/myapp?sslmode=disable

# JWT
JWT_SECRET=change-me-in-production
JWT_ACCESS_EXPIRY=15m
JWT_REFRESH_EXPIRY=168h

# Redis
REDIS_URL=redis://localhost:6379

# Storage
STORAGE_DRIVER=minio
MINIO_ENDPOINT=http://localhost:9000
MINIO_ACCESS_KEY=minioadmin
MINIO_SECRET_KEY=minioadmin
MINIO_BUCKET=myapp-uploads
MINIO_REGION=us-east-1
MINIO_USE_SSL=false

# Email
RESEND_API_KEY=re_your_api_key
MAIL_FROM=noreply@myapp.dev

# CORS
CORS_ORIGINS=http://localhost:3000,http://localhost:3001

# GORM Studio
GORM_STUDIO_ENABLED=true

# AI
AI_PROVIDER=claude
AI_API_KEY=
AI_MODEL=claude-sonnet-4-5-20250929`} />
            </div>

            <div className="prose-grit">
              <h2>Production Checklist</h2>
              <p>
                Before deploying to production, make sure you have addressed these configuration items:
              </p>
              <ul>
                <li>Change <code>JWT_SECRET</code> to a strong random string (at least 32 characters)</li>
                <li>Set <code>APP_ENV=production</code> to disable debug logging and GORM Studio</li>
                <li>Set <code>GORM_STUDIO_ENABLED=false</code> to disable the database browser</li>
                <li>Update <code>DATABASE_URL</code> to point to your production database</li>
                <li>Update <code>CORS_ORIGINS</code> to your production domain names</li>
                <li>Set <code>STORAGE_DRIVER</code> to <code>r2</code> or <code>b2</code> and configure cloud credentials</li>
                <li>Set <code>RESEND_API_KEY</code> with your production API key and verify your sender domain</li>
                <li>Update <code>REDIS_URL</code> to your production Redis instance</li>
                <li>Set <code>APP_URL</code> to your production API domain</li>
              </ul>
            </div>

            {/* Nav */}
            <div className="flex flex-wrap gap-3 mt-12 pt-6 border-t border-border/30">
              <Button variant="outline" asChild className="border-border/60 bg-transparent hover:bg-accent/50">
                <Link href="/docs/getting-started/project-structure">
                  Project Structure
                </Link>
              </Button>
              <Button asChild className="glow-purple-sm ml-auto">
                <Link href="/docs/concepts/architecture">
                  Architecture Overview
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
