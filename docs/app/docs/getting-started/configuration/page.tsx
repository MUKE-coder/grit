import Link from 'next/link'
import { ArrowRight } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'
import { CodeBlock } from '@/components/code-block'
import { getDocMetadata } from '@/config/docs-metadata'

export const metadata = getDocMetadata('/docs/getting-started/configuration')

export default function ConfigurationPage() {
  return (
    <div className="min-h-screen bg-background isolate">
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
                <li><strong><code>.env</code></strong> -- Your actual configuration, generated at scaffold time with strong random secrets (JWT, Postgres, Pulse, Sentinel). This file is gitignored and never committed.</li>
                <li><strong><code>.env.example</code></strong> -- Documented template with all variables and sensible defaults. Committed to git.</li>
                <li><strong><code>.env.cloud.example</code></strong> -- Template for cloud-only setups (Neon, Upstash, R2/B2) when you do not have Docker.</li>
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
                  { variable: 'APP_ENV', default: 'development', desc: 'Controls logging verbosity, GORM Studio visibility, and error detail level. Set to production in deployed environments to disable debug features. Note: Pulse and Sentinel refuse to mount in production if their passwords are still the literal defaults.' },
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
                Grit uses PostgreSQL as its primary database, connected through GORM. Rather than
                a single connection string, the scaffold treats the <code>POSTGRES_*</code> parts
                as the single source of truth -- both <code>docker-compose.yml</code> and the Go
                API read them. The Go API builds the connection DSN from these parts automatically
                when <code>DATABASE_URL</code> is left empty.
              </p>
            </div>

            <div className="mb-10">
              <CodeBlock language="bash" filename=".env" code={`# Database (Postgres) — single source of truth
# Edit ONLY the POSTGRES_* values below. docker-compose.yml reads them via
# \${VAR} substitution and the Go API builds the DSN from these parts when
# DATABASE_URL is empty.
POSTGRES_USER=grit
POSTGRES_PASSWORD=change-me          # generated per-scaffold; MUST change in production
POSTGRES_DB=myapp
POSTGRES_HOST=localhost              # "postgres" inside docker-compose.prod.yml
POSTGRES_PORT=5434                   # host port; 5432 inside the docker network

# Override the connection string ONLY for external Postgres (Neon, Supabase,
# RDS) or SQLite. When set, this wins over the POSTGRES_* parts above.
#   DATABASE_URL=postgres://user:pass@host:5432/db?sslmode=require
#   DATABASE_URL=sqlite:./app.db           # pure-Go driver, no CGO
#   DATABASE_URL=sqlite::memory:           # gone on restart, great for tests
# DATABASE_URL=`} />
              <div className="mt-4 space-y-3">
                {[
                  { variable: 'POSTGRES_USER', default: 'grit', desc: 'Postgres role the API connects as. Matches the user created by docker-compose.yml.' },
                  { variable: 'POSTGRES_PASSWORD', default: '(generated)', desc: 'Password for the Postgres role. grit new generates a strong random value per project so a fresh scaffold runs without editing. Change it in production.' },
                  { variable: 'POSTGRES_DB', default: 'Project name', desc: 'Name of the database to connect to. Defaults to your project name.' },
                  { variable: 'POSTGRES_HOST', default: 'localhost', desc: 'Database host. Stays localhost for local Docker; docker-compose.prod.yml overrides it to "postgres" (the service name) for inter-container traffic.' },
                  { variable: 'POSTGRES_PORT', default: '5434', desc: 'Host port for Postgres. Grit uses 5434 (not the default 5432) to avoid collisions with a Postgres you may already run locally. The container still listens on 5432 inside the Docker network.' },
                  { variable: 'DATABASE_URL', default: 'optional', desc: 'Optional full connection string. Leave commented to build the DSN from the POSTGRES_* parts. Set it to point at an external Postgres (postgres://... with sslmode=require) or to use SQLite (sqlite:./app.db or sqlite::memory:). When set, it wins over the POSTGRES_* parts.' },
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

            {/* JWT Config */}
            <div className="prose-grit">
              <h2>JWT Authentication</h2>
              <p>
                Grit uses JWT tokens for authentication with separate access and refresh tokens.
                The access token is short-lived (15 minutes) and the refresh token lasts 7 days.
                <code>JWT_SECRET</code> is generated at scaffold time.
              </p>
            </div>

            <div className="mb-10">
              <CodeBlock language="bash" filename=".env" code={`# JWT — generated at scaffold time. Rotate with: openssl rand -hex 32
JWT_SECRET=change-me-in-production
JWT_ACCESS_EXPIRY=15m        # Access token lifetime
JWT_REFRESH_EXPIRY=168h      # Refresh token lifetime (7 days)`} />
              <div className="mt-4 space-y-3">
                {[
                  { variable: 'JWT_SECRET', default: '(generated)', desc: 'The secret key used to sign and verify JWT tokens. Generated per-scaffold (openssl rand -hex 32). Rotate it in production. Both access and refresh tokens use this same secret.' },
                  { variable: 'JWT_ACCESS_EXPIRY', default: '15m', desc: 'How long access tokens are valid. Uses Go duration format: 15m (15 minutes), 1h (1 hour), etc. Keep short for security. The frontend automatically refreshes expired tokens.' },
                  { variable: 'JWT_REFRESH_EXPIRY', default: '168h', desc: 'How long refresh tokens are valid. 168h is 7 days. Users must re-login after this period.' },
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

            {/* OAuth / Social Login */}
            <div className="prose-grit">
              <h2>OAuth &amp; Social Login</h2>
              <p>
                Grit ships with Google and GitHub OAuth2 sign-in. The OAuth API routes are always
                registered server-side; <code>SOCIAL_AUTH_ENABLED</code> only controls whether the
                social buttons render on the auth pages. Leave the client IDs and secrets empty to
                run without social login.
              </p>
            </div>

            <div className="mb-10">
              <CodeBlock language="bash" filename=".env" code={`# OAuth2 — Social Login (Google + GitHub)
# Google: https://console.cloud.google.com/apis/credentials
GOOGLE_CLIENT_ID=                    # Google OAuth 2.0 Client ID
GOOGLE_CLIENT_SECRET=                # Google OAuth 2.0 Client Secret
# GitHub: https://github.com/settings/developers
GITHUB_CLIENT_ID=                    # GitHub OAuth App Client ID
GITHUB_CLIENT_SECRET=                # GitHub OAuth App Client Secret
OAUTH_FRONTEND_URL=http://localhost:3001  # Where to redirect after OAuth

# Toggle the social buttons + "or continue with" divider on auth pages.
# Mirrored to the browser bundle as NEXT_PUBLIC_SOCIAL_AUTH_ENABLED.
SOCIAL_AUTH_ENABLED=true`} />
              <div className="mt-4 space-y-3">
                {[
                  { variable: 'GOOGLE_CLIENT_ID / GOOGLE_CLIENT_SECRET', default: '(your keys)', desc: 'Google OAuth 2.0 credentials. Create them in the Google Cloud Console under APIs & Services > Credentials.' },
                  { variable: 'GITHUB_CLIENT_ID / GITHUB_CLIENT_SECRET', default: '(your keys)', desc: 'GitHub OAuth App credentials. Create an OAuth App under GitHub Settings > Developers.' },
                  { variable: 'OAUTH_FRONTEND_URL', default: 'http://localhost:3001', desc: 'The frontend URL the API redirects back to after a successful OAuth exchange. Defaults to the admin panel dev port.' },
                  { variable: 'SOCIAL_AUTH_ENABLED', default: 'true', desc: 'Show or hide the Google/GitHub buttons on auth pages. Set to false to remove them and the "or continue with" divider. The OAuth routes stay registered either way. Mirrored to the browser as NEXT_PUBLIC_SOCIAL_AUTH_ENABLED.' },
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
REDIS_URL=redis://localhost:6380`} />
              <div className="mt-4 space-y-3">
                <div className="rounded-lg border border-border/30 bg-card/30 px-4 py-3">
                  <div className="flex items-center gap-2 mb-1">
                    <code className="text-sm font-mono text-primary/70 font-medium">REDIS_URL</code>
                    <span className="text-[10px] font-mono text-muted-foreground/40 bg-accent/50 px-1.5 py-0.5 rounded">redis://localhost:6380</span>
                  </div>
                  <p className="text-xs text-muted-foreground/60 leading-relaxed">
                    Redis connection URL. For local Docker, Grit maps Redis to host port 6380 (not the default 6379) to avoid collisions with a local Redis. For cloud Redis (Upstash), use the <code className="text-xs font-mono text-primary/50">rediss://</code> protocol (with double s) and include the password: <code className="text-xs font-mono text-primary/50">rediss://default:password@endpoint:6379</code>.
                  </p>
                </div>
              </div>
            </div>

            {/* Frontend & Public URLs */}
            <div className="prose-grit">
              <h2>Frontend &amp; Public URLs</h2>
              <p>
                URLs baked into the Next.js web and admin bundles at build time. These tell the
                clients where to reach the API and the admin panel.
              </p>
            </div>

            <div className="mb-10">
              <CodeBlock language="bash" filename=".env" code={`# Public API URL — baked into Next.js bundles at build time
API_URL=http://localhost:8080

# Admin panel URL — surfaced in the web app's navbar + landing-page dev links.
# Set to your production admin origin (e.g. https://admin.example.com) before shipping.
NEXT_PUBLIC_ADMIN_URL=http://localhost:3001`} />
              <div className="mt-4 space-y-3">
                {[
                  { variable: 'API_URL', default: 'http://localhost:8080', desc: 'The public API base URL baked into the Next.js bundles at build time. Point it at your deployed API origin (e.g. https://api.example.com) for production builds.' },
                  { variable: 'NEXT_PUBLIC_ADMIN_URL', default: 'http://localhost:3001', desc: 'The admin panel URL surfaced in the web app navbar and landing-page dev links. Defaults to the admin dev port; set to your production admin origin before shipping.' },
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

            {/* Storage Config */}
            <div className="prose-grit">
              <h2>File Storage</h2>
              <p>
                Grit supports four S3-compatible storage providers: MinIO (local development), AWS
                S3, Cloudflare R2, and Backblaze B2. The <code>STORAGE_DRIVER</code> variable
                controls which provider is active -- only the variables for the active driver need
                to be set.
              </p>
            </div>

            <div className="mb-10">
              <CodeBlock language="bash" filename=".env" code={`# Storage — Which provider to use: minio, s3, r2, b2
STORAGE_DRIVER=minio

# MinIO — Local S3-compatible storage (default for development)
MINIO_ENDPOINT=http://localhost:9002
MINIO_ACCESS_KEY=minioadmin
MINIO_SECRET_KEY=minioadmin
MINIO_BUCKET=myapp-uploads
MINIO_REGION=us-east-1
MINIO_USE_SSL=false

# AWS S3 (used when STORAGE_DRIVER=s3)
# Leave S3_ENDPOINT empty to use the AWS regional default. S3_ACCESS_KEY /
# S3_SECRET_KEY / S3_REGION fall back to AWS_ACCESS_KEY_ID / AWS_SECRET_ACCESS_KEY
# / AWS_REGION, so an IAM role on your EC2 / ECS / Lambda Just Works.
S3_ENDPOINT=
S3_ACCESS_KEY=
S3_SECRET_KEY=
S3_BUCKET=
S3_REGION=us-east-1

# Cloudflare R2 (used when STORAGE_DRIVER=r2)
R2_ENDPOINT=
R2_ACCESS_KEY=               # R2 Access Key ID
R2_SECRET_KEY=               # R2 Secret Access Key
R2_BUCKET=
R2_REGION=auto               # Always "auto" for R2

# Backblaze B2 (used when STORAGE_DRIVER=b2)
B2_ENDPOINT=
B2_ACCESS_KEY=               # B2 keyID
B2_SECRET_KEY=               # B2 applicationKey
B2_BUCKET=
B2_REGION=us-west-004        # Must match your bucket region`} />
              <div className="mt-4 space-y-3">
                {[
                  { variable: 'STORAGE_DRIVER', default: 'minio', desc: 'Which storage provider to use. Options: minio (local dev with Docker), s3 (AWS S3), r2 (Cloudflare R2), b2 (Backblaze B2). Only the variables for the active driver need to be set.' },
                  { variable: 'MINIO_ENDPOINT', default: 'http://localhost:9002', desc: 'MinIO server URL. Grit maps MinIO to host port 9002 (API) and 9003 (web console) to avoid the default 9000/9001 colliding with other services.' },
                  { variable: 'MINIO_ACCESS_KEY / MINIO_SECRET_KEY', default: 'minioadmin', desc: 'Default MinIO credentials. These match the Docker Compose configuration. Change in production if running your own MinIO instance.' },
                  { variable: 'S3_ENDPOINT', default: '(empty)', desc: 'AWS S3 endpoint. Leave empty to use the AWS regional default (s3.<region>.amazonaws.com) with virtual-hosted-style addressing. Only used when STORAGE_DRIVER=s3.' },
                  { variable: 'S3_ACCESS_KEY / S3_SECRET_KEY / S3_REGION', default: '(IAM fallback)', desc: 'AWS credentials and region. Fall back to AWS_ACCESS_KEY_ID / AWS_SECRET_ACCESS_KEY / AWS_REGION, so leaving them empty lets an attached IAM role supply the credentials automatically.' },
                  { variable: 'R2_ENDPOINT', default: '(your account)', desc: 'Cloudflare R2 endpoint. Format: https://ACCOUNT_ID.r2.cloudflarestorage.com. Only used when STORAGE_DRIVER=r2.' },
                  { variable: 'R2_ACCESS_KEY / R2_SECRET_KEY', default: '(your keys)', desc: 'Cloudflare R2 API credentials. Create an API token in the Cloudflare Dashboard under R2 > Manage R2 API Tokens. R2_REGION is always "auto".' },
                  { variable: 'B2_ENDPOINT', default: '(your region)', desc: 'Backblaze B2 S3-compatible endpoint. Format depends on your bucket region. B2_ACCESS_KEY is the keyID and B2_SECRET_KEY is the applicationKey. Only used when STORAGE_DRIVER=b2.' },
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
MAIL_FROM=noreply@myapp.dev

# Support inbox — every ticket opened in /system/support is emailed here
# (when RESEND_API_KEY is set). Leave empty in dev to skip email-out.
SUPPORT_EMAIL=`} />
              <div className="mt-4 space-y-3">
                {[
                  { variable: 'RESEND_API_KEY', default: 're_your_api_key', desc: 'Your Resend API key. Get one at resend.com/api-keys. In development, emails are sent to Mailhog regardless of this key.' },
                  { variable: 'MAIL_FROM', default: 'noreply@myapp.dev', desc: 'The sender email address for all outgoing emails. Must be a verified domain in your Resend account for production use.' },
                  { variable: 'SUPPORT_EMAIL', default: '(empty)', desc: 'Destination inbox for support tickets opened in /system/support. Tickets are emailed here when RESEND_API_KEY is set. Leave empty in dev to skip sending.' },
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
                frontend origins are allowed to make requests. The Wails desktop webview does not
                need an entry -- its origin is matched by host (<code>wails.localhost</code>) on any port.
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
                API server, behind its own basic-auth login.
              </p>
            </div>

            <div className="mb-10">
              <CodeBlock language="bash" filename=".env" code={`# GORM Studio — Visual database browser
GORM_STUDIO_ENABLED=true
GORM_STUDIO_USERNAME=admin              # Login username for the Studio UI
GORM_STUDIO_PASSWORD=studio             # Login password for the Studio UI`} />
              <div className="mt-4 space-y-3">
                {[
                  { variable: 'GORM_STUDIO_ENABLED', default: 'true', desc: 'Enable or disable GORM Studio. Set to true in development for visual database browsing. Set to false in production to disable the browser and prevent unauthorized access to your data.' },
                  { variable: 'GORM_STUDIO_USERNAME', default: 'admin', desc: 'Login username for the Studio UI.' },
                  { variable: 'GORM_STUDIO_PASSWORD', default: 'studio', desc: 'Login password for the Studio UI. Change it if you keep Studio enabled outside local development.' },
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

            {/* AI Config */}
            <div className="prose-grit">
              <h2>AI Integration</h2>
              <p>
                Grit ships with built-in AI support via Vercel AI Gateway (one key, hundreds of models).
                The AI service provides text completion and streaming endpoints.
              </p>
            </div>

            <div className="mb-10">
              <CodeBlock language="bash" filename=".env" code={`# AI — Vercel AI Gateway (one key, hundreds of models)
AI_GATEWAY_API_KEY=                           # Get from vercel.com/ai-gateway
AI_GATEWAY_MODEL=anthropic/claude-sonnet-4-6  # provider/model format
AI_GATEWAY_URL=https://ai-gateway.vercel.sh/v1`} />
              <div className="mt-4 space-y-3">
                {[
                  { variable: 'AI_GATEWAY_API_KEY', default: '(your key)', desc: 'Your Vercel AI Gateway API key. Get one at vercel.com/ai-gateway. A single key gives you access to all providers (Anthropic, OpenAI, Google, and more). Leave empty if you do not use AI features.' },
                  { variable: 'AI_GATEWAY_MODEL', default: 'anthropic/claude-sonnet-4-6', desc: 'The model to use, in provider/model format. Examples: anthropic/claude-sonnet-4-6, openai/gpt-5.4, google/gemini-2.5-pro.' },
                  { variable: 'AI_GATEWAY_URL', default: 'https://ai-gateway.vercel.sh/v1', desc: 'The Vercel AI Gateway endpoint URL. This is the same for all providers and models. You should not need to change this.' },
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

            {/* TOTP Config */}
            <div className="prose-grit">
              <h2>Two-Factor Authentication</h2>
              <p>
                Grit supports TOTP-based two-factor authentication. The issuer name appears in
                authenticator apps (e.g., Google Authenticator, Authy) alongside the user&apos;s account.
              </p>
            </div>

            <div className="mb-10">
              <CodeBlock language="bash" filename=".env" code={`# Two-Factor Authentication (TOTP)
TOTP_ISSUER=myapp`} />
              <div className="mt-4 space-y-3">
                <div className="rounded-lg border border-border/30 bg-card/30 px-4 py-3">
                  <div className="flex items-center gap-2 mb-1">
                    <code className="text-sm font-mono text-primary/70 font-medium">TOTP_ISSUER</code>
                    <span className="text-[10px] font-mono text-muted-foreground/40 bg-accent/50 px-1.5 py-0.5 rounded">myapp</span>
                  </div>
                  <p className="text-xs text-muted-foreground/60 leading-relaxed">
                    The issuer name displayed in authenticator apps when users set up 2FA. Set this to your application or company name.
                  </p>
                </div>
              </div>
            </div>

            {/* Observability / Pulse */}
            <div className="prose-grit">
              <h2>Observability (Pulse)</h2>
              <p>
                Pulse is Grit&apos;s built-in performance monitoring, request-tracing, and
                error-tracking dashboard. Its password is generated at scaffold time -- Pulse
                refuses to mount in <code>APP_ENV=production</code> while the password is still the
                literal default <code>pulse</code>.
              </p>
            </div>

            <div className="mb-10">
              <CodeBlock language="bash" filename=".env" code={`# Observability — Pulse performance monitoring dashboard
PULSE_ENABLED=true                   # Set to "false" to disable Pulse entirely
PULSE_USERNAME=admin                 # Dashboard login username
PULSE_PASSWORD=pulse                 # Generated per-scaffold. Rotate: openssl rand -hex 16`} />
              <div className="mt-4 space-y-3">
                {[
                  { variable: 'PULSE_ENABLED', default: 'true', desc: 'Turn the Pulse dashboard on or off. Set to false to disable Pulse entirely.' },
                  { variable: 'PULSE_USERNAME', default: 'admin', desc: 'Login username for the Pulse dashboard.' },
                  { variable: 'PULSE_PASSWORD', default: '(generated)', desc: 'Login password for the Pulse dashboard. Generated per-scaffold so the production gate is satisfied out of the box. Rotate with openssl rand -hex 16.' },
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

            {/* Security / Sentinel */}
            <div className="prose-grit">
              <h2>Security (Sentinel)</h2>
              <p>
                Sentinel is Grit&apos;s built-in WAF, rate limiter, and threat-detection layer with
                its own dashboard. Like Pulse, it refuses to mount in production while its
                credentials are still the defaults. Both <code>SENTINEL_PASSWORD</code> and
                <code>SENTINEL_SECRET_KEY</code> are generated at scaffold time.
              </p>
            </div>

            <div className="mb-10">
              <CodeBlock language="bash" filename=".env" code={`# Security — Sentinel WAF, rate limiting, threat detection
SENTINEL_ENABLED=true                # Set to "false" to disable Sentinel entirely
SENTINEL_USERNAME=admin              # Dashboard login username
SENTINEL_PASSWORD=sentinel           # Generated per-scaffold. Rotate: openssl rand -hex 16
SENTINEL_SECRET_KEY=change-me        # Generated per-scaffold (>=32 bytes). Rotate: openssl rand -hex 32`} />
              <div className="mt-4 space-y-3">
                {[
                  { variable: 'SENTINEL_ENABLED', default: 'true', desc: 'Turn Sentinel on or off. Set to false to disable the WAF, rate limiting, and threat detection entirely.' },
                  { variable: 'SENTINEL_USERNAME', default: 'admin', desc: 'Login username for the Sentinel dashboard.' },
                  { variable: 'SENTINEL_PASSWORD', default: '(generated)', desc: 'Login password for the Sentinel dashboard. Generated per-scaffold. Rotate with openssl rand -hex 16.' },
                  { variable: 'SENTINEL_SECRET_KEY', default: '(generated)', desc: 'Secret used to sign the Sentinel dashboard JWT sessions. Needs at least 32 bytes of entropy; generated per-scaffold. Rotate with openssl rand -hex 32.' },
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

            {/* Theme */}
            <div className="prose-grit">
              <h2>Theme</h2>
              <p>
                Picks the visual identity for the auth pages and dashboard. The value is mirrored
                into <code>NEXT_PUBLIC_THEME</code> for the web and admin clients in
                <code>next.config</code>, so server components render the right theme without a
                flash of unstyled content.
              </p>
            </div>

            <div className="mb-10">
              <CodeBlock language="bash" filename=".env" code={`# Theme — visual identity for auth pages + dashboard
#   atlas  — split-screen, team/organisation, Inter (default)
#   aurora — centered Clerk-style, consumer SaaS, Geist
#   pulse  — split-screen with carousel, ecommerce/brand, Onest + DM Serif
THEME=atlas`} />
              <div className="mt-4 space-y-3">
                <div className="rounded-lg border border-border/30 bg-card/30 px-4 py-3">
                  <div className="flex items-center gap-2 mb-1">
                    <code className="text-sm font-mono text-primary/70 font-medium">THEME</code>
                    <span className="text-[10px] font-mono text-muted-foreground/40 bg-accent/50 px-1.5 py-0.5 rounded">atlas</span>
                  </div>
                  <p className="text-xs text-muted-foreground/60 leading-relaxed">
                    Chosen with <code className="text-xs font-mono text-primary/50">--theme</code> at scaffold time (defaults to atlas). Options: <code className="text-xs font-mono text-primary/50">atlas</code> (split-screen, team/organisation, Inter), <code className="text-xs font-mono text-primary/50">aurora</code> (centered Clerk-style, consumer SaaS, Geist), and <code className="text-xs font-mono text-primary/50">pulse</code> (split-screen with carousel, ecommerce/brand, Onest + DM Serif). Mirrored to the clients as NEXT_PUBLIC_THEME.
                  </p>
                </div>
              </div>
            </div>

            {/* Local Service Ports */}
            <div className="prose-grit">
              <h2>Local Service Ports</h2>
              <p>
                The services started by Docker Compose and the dev servers each bind their own host
                port. Grit shifts the infrastructure ports off their defaults (Postgres, Redis,
                MinIO) to avoid colliding with services you may already run locally.
              </p>
            </div>

            <div className="mb-10">
              <table className="w-full border-collapse">
                <thead>
                  <tr className="border-b border-border">
                    <th className="text-left text-sm font-semibold text-foreground pb-3 pr-4">Service</th>
                    <th className="text-left text-sm font-semibold text-foreground pb-3 pr-4">Dev URL</th>
                    <th className="text-left text-sm font-semibold text-foreground pb-3">Host Port</th>
                  </tr>
                </thead>
                <tbody>
                  {[
                    { app: 'Go API', url: 'http://localhost:8080', port: '8080' },
                    { app: 'GORM Studio', url: 'http://localhost:8080/studio', port: '8080' },
                    { app: 'Web App (Next.js)', url: 'http://localhost:3000', port: '3000' },
                    { app: 'Admin Panel (Next.js)', url: 'http://localhost:3001', port: '3001' },
                    { app: 'PostgreSQL', url: 'localhost:5434', port: '5434' },
                    { app: 'Redis', url: 'localhost:6380', port: '6380' },
                    { app: 'MinIO (API)', url: 'http://localhost:9002', port: '9002' },
                    { app: 'MinIO Console', url: 'http://localhost:9003', port: '9003' },
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

# Database (Postgres) — edit these; DATABASE_URL is optional
POSTGRES_USER=grit
POSTGRES_PASSWORD=change-me
POSTGRES_DB=myapp
POSTGRES_HOST=localhost
POSTGRES_PORT=5434
# DATABASE_URL=postgres://user:pass@host:5432/db?sslmode=require
# DATABASE_URL=sqlite:./app.db
# DATABASE_URL=sqlite::memory:

# JWT
JWT_SECRET=change-me-in-production
JWT_ACCESS_EXPIRY=15m
JWT_REFRESH_EXPIRY=168h

# OAuth2 — Social Login (Google + GitHub)
GOOGLE_CLIENT_ID=
GOOGLE_CLIENT_SECRET=
GITHUB_CLIENT_ID=
GITHUB_CLIENT_SECRET=
OAUTH_FRONTEND_URL=http://localhost:3001
SOCIAL_AUTH_ENABLED=true

# Redis
REDIS_URL=redis://localhost:6380

# Frontend & Public URLs
API_URL=http://localhost:8080
NEXT_PUBLIC_ADMIN_URL=http://localhost:3001

# Storage — minio, s3, r2, b2
STORAGE_DRIVER=minio

# MinIO (local dev)
MINIO_ENDPOINT=http://localhost:9002
MINIO_ACCESS_KEY=minioadmin
MINIO_SECRET_KEY=minioadmin
MINIO_BUCKET=myapp-uploads
MINIO_REGION=us-east-1
MINIO_USE_SSL=false

# AWS S3
S3_ENDPOINT=
S3_ACCESS_KEY=
S3_SECRET_KEY=
S3_BUCKET=
S3_REGION=us-east-1

# Cloudflare R2
R2_ENDPOINT=
R2_ACCESS_KEY=
R2_SECRET_KEY=
R2_BUCKET=
R2_REGION=auto

# Backblaze B2
B2_ENDPOINT=
B2_ACCESS_KEY=
B2_SECRET_KEY=
B2_BUCKET=
B2_REGION=us-west-004

# Email
RESEND_API_KEY=re_your_api_key
MAIL_FROM=noreply@myapp.dev
SUPPORT_EMAIL=

# CORS
CORS_ORIGINS=http://localhost:3000,http://localhost:3001

# GORM Studio
GORM_STUDIO_ENABLED=true
GORM_STUDIO_USERNAME=admin
GORM_STUDIO_PASSWORD=studio

# AI (Vercel AI Gateway)
AI_GATEWAY_API_KEY=
AI_GATEWAY_MODEL=anthropic/claude-sonnet-4-6
AI_GATEWAY_URL=https://ai-gateway.vercel.sh/v1

# Two-Factor Authentication (TOTP)
TOTP_ISSUER=myapp

# Observability (Pulse)
PULSE_ENABLED=true
PULSE_USERNAME=admin
PULSE_PASSWORD=pulse

# Security (Sentinel)
SENTINEL_ENABLED=true
SENTINEL_USERNAME=admin
SENTINEL_PASSWORD=sentinel
SENTINEL_SECRET_KEY=change-me

# Theme — atlas, aurora, pulse
THEME=atlas`} />
            </div>

            <div className="prose-grit">
              <h2>Production Checklist</h2>
              <p>
                Before deploying to production, make sure you have addressed these configuration items:
              </p>
              <ul>
                <li>Change <code>JWT_SECRET</code> to a strong random string (at least 32 characters)</li>
                <li>Set <code>APP_ENV=production</code> to disable debug logging and GORM Studio</li>
                <li>Set <code>GORM_STUDIO_ENABLED=false</code> (or change its username/password) to lock down the database browser</li>
                <li>Change <code>POSTGRES_PASSWORD</code> (or set <code>DATABASE_URL</code> to your managed Postgres with <code>sslmode=require</code>)</li>
                <li>Update <code>CORS_ORIGINS</code> to your production domain names</li>
                <li>Set <code>STORAGE_DRIVER</code> to <code>s3</code>, <code>r2</code>, or <code>b2</code> and configure the cloud credentials</li>
                <li>Set <code>RESEND_API_KEY</code> with your production API key, verify your sender domain, and set <code>SUPPORT_EMAIL</code></li>
                <li>Update <code>REDIS_URL</code> to your production Redis instance</li>
                <li>Set <code>API_URL</code> and <code>NEXT_PUBLIC_ADMIN_URL</code> to your production origins</li>
                <li>Rotate <code>PULSE_PASSWORD</code>, <code>SENTINEL_PASSWORD</code>, and <code>SENTINEL_SECRET_KEY</code> off the literal defaults (they gate mounting in production)</li>
                <li>Fill in OAuth credentials if you use social login, or set <code>SOCIAL_AUTH_ENABLED=false</code></li>
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
