package scaffold

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"path/filepath"
)

// randomHex returns 2*n hex characters of cryptographically random bytes.
// Used to seed JWT, Sentinel and Pulse secrets at scaffold time so a fresh
// `grit new` project boots cleanly in production mode without the user
// having to manually replace placeholder credentials.
func randomHex(n int) string {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		// crypto/rand should never fail on a real OS; if it did the scaffold
		// would be unusable. Fall back to a clearly-fake string so the user
		// at least sees something obvious to change rather than a silent
		// weak default.
		return "REPLACE_ME_crypto_rand_failed"
	}
	return hex.EncodeToString(b)
}

func writeRootFiles(root string, opts Options) error {
	files := map[string]string{
		filepath.Join(root, ".env"):               envFile(opts),
		filepath.Join(root, ".env.example"):       envFile(opts), // Same as .env — serves as documentation for other devs
		filepath.Join(root, ".gitignore"):         rootGitignore(),
		filepath.Join(root, "README.md"):          readmeFile(opts),
		filepath.Join(root, "grit.json"):          gritJSON(opts),
		filepath.Join(root, ".claude", "skills", "grit", "SKILL.md"):      gritSkillFile(opts),
		filepath.Join(root, ".claude", "skills", "grit", "reference.md"):  gritSkillReference(opts),
	}

	// Prettier config (all architectures with frontend)
	if opts.Architecture != ArchAPI {
		files[filepath.Join(root, ".prettierrc")] = prettierConfig()
		files[filepath.Join(root, ".prettierignore")] = prettierIgnore()
	}

	if opts.ShouldUseTurborepo() {
		files[filepath.Join(root, "pnpm-workspace.yaml")] = pnpmWorkspace()
		files[filepath.Join(root, "turbo.json")] = turboJSON()
		files[filepath.Join(root, "package.json")] = rootPackageJSON(opts)
		files[filepath.Join(root, "grit.config.ts")] = gritConfig(opts)
		files[filepath.Join(root, "postcss.config.mjs")] = rootPostCSSConfig()
	}

	for path, content := range files {
		if err := writeFile(path, content); err != nil {
			return fmt.Errorf("writing %s: %w", path, err)
		}
	}

	return nil
}

func envFile(opts Options) string {
	// Generated per-scaffold so APP_ENV=production works out of the box.
	// Rotate any of these any time with `openssl rand -hex 32`.
	jwtSecret := randomHex(32)
	sentinelPassword := randomHex(16)
	sentinelSecretKey := randomHex(32)
	pulsePassword := randomHex(16)
	postgresPassword := randomHex(24) // strong default; new project just works

	return fmt.Sprintf(`# %s — Environment Variables

# App
APP_NAME=%s
APP_ENV=development
APP_PORT=8080
APP_URL=http://localhost:8080

# ─── Database (Postgres) ────────────────────────────────────────────────
# Single source of truth: edit ONLY the POSTGRES_* values below.
#   - docker-compose.yml reads them via ${VAR} substitution
#   - the Go API builds DATABASE_URL from these parts when DATABASE_URL is empty
# POSTGRES_PASSWORD is generated per-scaffold so 'grit migrate' works
# without any editing on a fresh project.
POSTGRES_USER=grit
POSTGRES_PASSWORD=%s
POSTGRES_DB=%s
POSTGRES_HOST=localhost
# 5434 (not the default 5432) avoids host collisions — see docker-compose.yml.
# The container still listens on 5432 inside the Docker network, which is
# why docker-compose.prod.yml overrides POSTGRES_PORT back to 5432 for
# inter-container traffic.
POSTGRES_PORT=5434

# Override the connection string ONLY if you're pointing at an external
# Postgres (Neon, Supabase, RDS) or want to use SQLite. When set, this
# wins over the POSTGRES_* parts above.
#   DATABASE_URL=postgres://user:pass@host:5432/db?sslmode=require
#   DATABASE_URL=sqlite:./app.db           # pure-Go driver, no CGO
#   DATABASE_URL=sqlite::memory:           # gone on restart, great for tests
# DATABASE_URL=

# JWT — generated at scaffold time. Rotate with: openssl rand -hex 32
JWT_SECRET=%s
JWT_ACCESS_EXPIRY=15m
JWT_REFRESH_EXPIRY=168h

# OAuth2 — Social Login (Google + GitHub)
GOOGLE_CLIENT_ID=
GOOGLE_CLIENT_SECRET=
GITHUB_CLIENT_ID=
GITHUB_CLIENT_SECRET=
OAUTH_FRONTEND_URL=http://localhost:3001

# Redis
REDIS_URL=redis://localhost:6379

# Public API URL — baked into Next.js bundles at build time
API_URL=http://localhost:8080

# Storage — Which provider to use: minio, s3, r2, b2
STORAGE_DRIVER=minio

# MinIO (local development — used when STORAGE_DRIVER=minio)
MINIO_ENDPOINT=http://localhost:9000
MINIO_ACCESS_KEY=minioadmin
MINIO_SECRET_KEY=minioadmin
MINIO_BUCKET=%s-uploads
MINIO_REGION=us-east-1
MINIO_USE_SSL=false

# AWS S3 (used when STORAGE_DRIVER=s3)
# Leave S3_ENDPOINT empty to use the AWS regional default.
# S3_ACCESS_KEY + S3_SECRET_KEY fall back to AWS_ACCESS_KEY_ID +
# AWS_SECRET_ACCESS_KEY (and S3_REGION to AWS_REGION) so an IAM role
# attached to your EC2 / ECS / Lambda Just Works.
S3_ENDPOINT=
S3_ACCESS_KEY=
S3_SECRET_KEY=
S3_BUCKET=
S3_REGION=us-east-1

# Cloudflare R2 (used when STORAGE_DRIVER=r2)
R2_ENDPOINT=
R2_ACCESS_KEY=
R2_SECRET_KEY=
R2_BUCKET=
R2_REGION=auto

# Backblaze B2 (used when STORAGE_DRIVER=b2)
B2_ENDPOINT=
B2_ACCESS_KEY=
B2_SECRET_KEY=
B2_BUCKET=
B2_REGION=us-west-004

# Email (Resend)
RESEND_API_KEY=re_your_api_key
MAIL_FROM=noreply@%s.dev

# CORS
CORS_ORIGINS=http://localhost:3000,http://localhost:3001

# GORM Studio
GORM_STUDIO_ENABLED=true
GORM_STUDIO_USERNAME=admin
GORM_STUDIO_PASSWORD=studio

# AI — Vercel AI Gateway (one key, hundreds of models)
AI_GATEWAY_API_KEY=                           # Get from vercel.com/ai-gateway
AI_GATEWAY_MODEL=anthropic/claude-sonnet-4-6  # provider/model format
AI_GATEWAY_URL=https://ai-gateway.vercel.sh/v1

# Two-Factor Authentication (TOTP)
TOTP_ISSUER=%s

# Observability — Pulse performance monitoring dashboard
# Pulse refuses to mount in APP_ENV=production with the literal default
# password "pulse" — the value below is generated at scaffold time so the
# gate is satisfied. Rotate with: openssl rand -hex 16
PULSE_ENABLED=true
PULSE_USERNAME=admin
PULSE_PASSWORD=%s

# Security — Sentinel WAF, rate limiting, threat detection
# Same as Pulse: Sentinel refuses to mount in production with default
# credentials. SECRET_KEY needs at least 32 bytes of entropy — both values
# below are generated per scaffold. Rotate with: openssl rand -hex 16 (password)
# / openssl rand -hex 32 (secret key).
SENTINEL_ENABLED=true
SENTINEL_USERNAME=admin
SENTINEL_PASSWORD=%s
SENTINEL_SECRET_KEY=%s
`,
		opts.ProjectName, opts.ProjectName, // banner + APP_NAME
		postgresPassword, opts.ProjectName, // POSTGRES_PASSWORD + POSTGRES_DB
		jwtSecret,
		opts.ProjectName, opts.ProjectName, opts.ProjectName, // MINIO_BUCKET + MAIL_FROM + TOTP_ISSUER
		pulsePassword, sentinelPassword, sentinelSecretKey,
	)
}

func envExampleFile(opts Options) string {
	return `# App — General application settings
APP_NAME=myapp              # Application name
APP_ENV=development         # Environment: development, staging, production
APP_PORT=8080               # API server port
APP_URL=http://localhost:8080

# ─── Database (Postgres) ────────────────────────────────────────────────
# Single source of truth. Edit ONLY the POSTGRES_* values below — both
# docker-compose.yml and the Go API read them. ` + "`grit new`" + ` generates a
# strong random POSTGRES_PASSWORD per project so a fresh scaffold runs
# without any editing.
POSTGRES_USER=grit
POSTGRES_PASSWORD=change-me          # MUST change in production
POSTGRES_DB=myapp
POSTGRES_HOST=localhost              # ` + "`postgres`" + ` inside docker-compose.prod.yml
POSTGRES_PORT=5434                   # host port; 5432 inside docker network

# Override the connection string ONLY for external Postgres (Neon,
# Supabase, RDS) or SQLite. When set, this wins over the POSTGRES_*
# parts above.
#   DATABASE_URL=postgres://user:pass@host:5432/db?sslmode=require
#   DATABASE_URL=sqlite:./app.db
#   DATABASE_URL=sqlite::memory:
# DATABASE_URL=

# JWT — Authentication tokens
JWT_SECRET=change-me-in-production   # MUST change in production
JWT_ACCESS_EXPIRY=15m                # Access token lifetime
JWT_REFRESH_EXPIRY=168h              # Refresh token lifetime (7 days)

# OAuth2 — Social Login (Google + GitHub)
# Google: https://console.cloud.google.com/apis/credentials
GOOGLE_CLIENT_ID=                    # Google OAuth 2.0 Client ID
GOOGLE_CLIENT_SECRET=                # Google OAuth 2.0 Client Secret
# GitHub: https://github.com/settings/developers
GITHUB_CLIENT_ID=                    # GitHub OAuth App Client ID
GITHUB_CLIENT_SECRET=                # GitHub OAuth App Client Secret
OAUTH_FRONTEND_URL=http://localhost:3001  # Where to redirect after OAuth

# Redis — Cache and job queue
REDIS_URL=redis://localhost:6379

# Public API URL — baked into Next.js bundles at build time
API_URL=https://api.example.com

# Storage — Active driver: minio, s3, r2, or b2
STORAGE_DRIVER=minio                 # Change to "s3", "r2", or "b2" to switch providers

# MinIO — Local S3-compatible storage (default for development)
MINIO_ENDPOINT=http://localhost:9000
MINIO_ACCESS_KEY=minioadmin
MINIO_SECRET_KEY=minioadmin
MINIO_BUCKET=myapp-uploads
MINIO_REGION=us-east-1
MINIO_USE_SSL=false

# AWS S3 — Native AWS object storage (most common production choice)
# Get credentials: IAM Console → Users → Security credentials → Access keys.
# Better still: attach an IAM role to your EC2 / ECS / Lambda and leave
# S3_ACCESS_KEY + S3_SECRET_KEY empty — the SDK auto-discovers role creds.
# Leave S3_ENDPOINT EMPTY so the SDK uses the AWS regional endpoint and
# virtual-hosted-style addressing (required for new buckets).
S3_ENDPOINT=                         # Empty = AWS default (s3.<region>.amazonaws.com)
S3_ACCESS_KEY=                       # Falls back to AWS_ACCESS_KEY_ID
S3_SECRET_KEY=                       # Falls back to AWS_SECRET_ACCESS_KEY
S3_BUCKET=myapp-uploads
S3_REGION=us-east-1                  # Falls back to AWS_REGION

# Cloudflare R2 — S3-compatible object storage with zero egress fees
# Get credentials: Cloudflare Dashboard → R2 → Manage R2 API Tokens
R2_ENDPOINT=https://YOUR_ACCOUNT_ID.r2.cloudflarestorage.com
R2_ACCESS_KEY=                       # R2 Access Key ID
R2_SECRET_KEY=                       # R2 Secret Access Key
R2_BUCKET=myapp-uploads
R2_REGION=auto                       # Always "auto" for R2

# Backblaze B2 — S3-compatible object storage, low cost per GB
# Get credentials: B2 Cloud Storage → App Keys
B2_ENDPOINT=https://s3.us-west-004.backblazeb2.com
B2_ACCESS_KEY=                       # B2 keyID
B2_SECRET_KEY=                       # B2 applicationKey
B2_BUCKET=myapp-uploads
B2_REGION=us-west-004               # Must match your bucket region

# Email — Resend integration
RESEND_API_KEY=re_your_api_key
MAIL_FROM=noreply@myapp.dev

# CORS — Allowed frontend origins (comma-separated)
CORS_ORIGINS=http://localhost:3000,http://localhost:3001

# GORM Studio — Visual database browser
GORM_STUDIO_ENABLED=true
GORM_STUDIO_USERNAME=admin              # Login username for the Studio UI
GORM_STUDIO_PASSWORD=studio             # Login password for the Studio UI

# AI — Vercel AI Gateway (one key, hundreds of models)
AI_GATEWAY_API_KEY=                           # Get from vercel.com/ai-gateway
AI_GATEWAY_MODEL=anthropic/claude-sonnet-4-6  # provider/model (e.g. openai/gpt-5.4, google/gemini-2.5-pro)
AI_GATEWAY_URL=https://ai-gateway.vercel.sh/v1

# Two-Factor Authentication (TOTP)
TOTP_ISSUER=myapp                    # App name shown in authenticator apps

# Observability — Pulse (performance monitoring, request tracing, error tracking)
PULSE_ENABLED=true                   # Set to "false" to disable Pulse entirely
PULSE_USERNAME=admin                 # Dashboard login username
PULSE_PASSWORD=pulse                 # Dashboard login password (change in production!)

# Security — Sentinel (WAF, rate limiting, threat detection)
SENTINEL_ENABLED=true                # Set to "false" to disable Sentinel entirely
SENTINEL_USERNAME=admin              # Dashboard login username
SENTINEL_PASSWORD=sentinel           # Dashboard login password (change in production!)
SENTINEL_SECRET_KEY=change-me        # Secret for dashboard JWT sessions
`
}

func envCloudExampleFile(opts Options) string {
	return fmt.Sprintf(`# %s — Cloud Environment Variables
#
# Use this file if you DON'T have Docker and want to use cloud services instead.
# Copy this to .env and fill in your keys:
#
#   cp .env.cloud.example .env
#
# No Docker required — just your API keys.

# ─── App ───────────────────────────────────────────────
APP_NAME=%s
APP_ENV=development
APP_PORT=8080
APP_URL=http://localhost:8080

# ─── Database (Neon — https://neon.tech) ───────────────
# Create a free project at neon.tech, copy the connection string
DATABASE_URL=postgres://user:password@ep-xxx-xxx-123456.us-east-2.aws.neon.tech/neondb?sslmode=require

# ─── JWT ───────────────────────────────────────────────
JWT_SECRET=change-me-to-a-random-string-at-least-32-chars
JWT_ACCESS_EXPIRY=15m
JWT_REFRESH_EXPIRY=168h

# ─── Redis (Upstash — https://upstash.com) ─────────────
# Create a free Redis database at upstash.com, copy the Redis URL
REDIS_URL=rediss://default:your-password@your-endpoint.upstash.io:6379

# ─── Storage ──────────────────────────────────────────
# Active driver: r2 or b2 (no minio in cloud mode)
STORAGE_DRIVER=r2

# ─── Cloudflare R2 (https://dash.cloudflare.com) ─────
# Dashboard → R2 → Create Bucket → Manage R2 API Tokens
R2_ENDPOINT=https://your-account-id.r2.cloudflarestorage.com
R2_ACCESS_KEY=your-r2-access-key-id
R2_SECRET_KEY=your-r2-secret-access-key
R2_BUCKET=%s-uploads
R2_REGION=auto

# ─── Backblaze B2 (https://backblaze.com/cloud-storage) ─
# B2 Cloud Storage → Buckets → Create → App Keys → Add Key
B2_ENDPOINT=https://s3.us-west-004.backblazeb2.com
B2_ACCESS_KEY=your-b2-key-id
B2_SECRET_KEY=your-b2-application-key
B2_BUCKET=%s-uploads
B2_REGION=us-west-004

# ─── Email (Resend — https://resend.com) ───────────────
# Sign up at resend.com, verify your domain, grab your API key
RESEND_API_KEY=re_your_api_key_here
MAIL_FROM=noreply@yourdomain.com

# ─── CORS ──────────────────────────────────────────────
CORS_ORIGINS=http://localhost:3000,http://localhost:3001

# ─── GORM Studio ──────────────────────────────────────
GORM_STUDIO_ENABLED=true
GORM_STUDIO_USERNAME=admin               # Login username for the Studio UI
GORM_STUDIO_PASSWORD=change-me-in-prod   # Login password — CHANGE THIS in production!

# ─── AI (Vercel AI Gateway) ──────────────────────────
AI_GATEWAY_API_KEY=your-gateway-key  # Get from vercel.com/ai-gateway
AI_GATEWAY_MODEL=anthropic/claude-sonnet-4-6  # provider/model format
AI_GATEWAY_URL=https://ai-gateway.vercel.sh/v1

# ─── TOTP (Two-Factor Authentication) ────────────────
TOTP_ISSUER=%s

# ─── Observability (Pulse) ────────────────────────────
PULSE_ENABLED=true
PULSE_USERNAME=admin
PULSE_PASSWORD=change-me-in-production

# ─── Security (Sentinel) ─────────────────────────────
SENTINEL_ENABLED=true
SENTINEL_USERNAME=admin
SENTINEL_PASSWORD=change-me-in-production
SENTINEL_SECRET_KEY=generate-a-random-string-here
`, opts.ProjectName, opts.ProjectName, opts.ProjectName, opts.ProjectName, opts.ProjectName)
}

func rootGitignore() string {
	return `# Dependencies
node_modules/
.pnpm-store/

# Build output
.next/
out/
dist/
build/
tmp/

# Go
*.exe
*.exe~
*.dll
*.so
*.dylib
*.test
*.out
migrate.exe

# Environment
.env
.env.local
.env.*.local

# IDE
.vscode/
.idea/
*.swp
*.swo
*~

# OS
.DS_Store
Thumbs.db

# Docker volumes
postgres-data/
redis-data/
minio-data/

# Turborepo
.turbo/

# Sentinel (WAF database)
sentinel.db
sentinel.db-shm
sentinel.db-wal

# Testing
e2e/test-results/
e2e/playwright-report/
coverage/

# Debug
*.log
npm-debug.log*
pnpm-debug.log*
`
}

func prettierConfig() string {
	return `{
  "semi": false,
  "singleQuote": true,
  "trailingComma": "es5",
  "tabWidth": 2,
  "printWidth": 100,
  "bracketSpacing": true,
  "arrowParens": "always",
  "endOfLine": "lf",
  "plugins": ["prettier-plugin-tailwindcss"]
}
`
}

func prettierIgnore() string {
	return `node_modules/
.next/
dist/
build/
out/
coverage/
.turbo/
pnpm-lock.yaml
*.min.js
*.min.css
`
}

func pnpmWorkspace() string {
	return `packages:
  - "apps/*"
  - "packages/*"
`
}

func turboJSON() string {
	return `{
  "$schema": "https://turbo.build/schema.json",
  "tasks": {
    "build": {
      "dependsOn": ["^build"],
      "outputs": [".next/**", "!.next/cache/**", "dist/**"]
    },
    "dev": {
      "cache": false,
      "persistent": true
    },
    "lint": {
      "dependsOn": ["^build"]
    },
    "type-check": {
      "dependsOn": ["^build"]
    }
  }
}
`
}

func rootPostCSSConfig() string {
	return `export default {
  plugins: {
    tailwindcss: {},
    autoprefixer: {},
  },
};
`
}

func rootPackageJSON(opts Options) string {
	scripts := fmt.Sprintf(`    "dev": "turbo dev",
    "build": "turbo build",
    "lint": "turbo lint",
    "type-check": "turbo type-check",
    "dev:api": "cd apps/api && air",`)

	if opts.ShouldIncludeWeb() {
		scripts += `
    "dev:web": "cd apps/web && pnpm dev",`
	}
	if opts.ShouldIncludeAdmin() {
		scripts += `
    "dev:admin": "cd apps/admin && pnpm dev",`
	}
	if opts.ShouldIncludeExpo() {
		scripts += `
    "dev:expo": "cd apps/expo && npx expo start",`
	}
	if opts.ShouldIncludeDesktop() {
		scripts += `
    "dev:desktop": "cd apps/desktop && wails dev",
    "build:desktop": "cd apps/desktop && wails build",`
	}
	if opts.ShouldIncludeDocs() {
		scripts += `
    "dev:docs": "cd apps/docs && pnpm dev",`
	}

	scripts += `
    "docker:up": "docker compose up -d",
    "docker:down": "docker compose down",
    "docker:logs": "docker compose logs -f",
    "test": "turbo test",
    "test:e2e": "playwright test",
    "test:e2e:ui": "playwright test --ui"`

	return fmt.Sprintf(`{
  "name": "%s",
  "private": true,
  "scripts": {
%s
  },
  "devDependencies": {
    "@playwright/test": "^1.48.0",
    "turbo": "^2.0.0"
  },
  "packageManager": "pnpm@10.0.0"
}
`, opts.ProjectName, scripts)
}

func gritConfig(opts Options) string {
	style := opts.Style
	if style == "" {
		style = "default"
	}
	config := fmt.Sprintf(`// Grit Framework Configuration
export default {
  name: "%s",
  style: "%s",
  api: {
    port: 8080,
    prefix: "/api",
  },`, opts.ProjectName, style)

	if opts.ShouldIncludeWeb() {
		config += `
  web: {
    port: 3000,
  },`
	}
	if opts.ShouldIncludeAdmin() {
		config += `
  admin: {
    port: 3001,
  },`
	}
	if opts.ShouldIncludeExpo() {
		config += `
  expo: {
    scheme: "` + opts.ProjectName + `",
  },`
	}
	if opts.ShouldIncludeDesktop() {
		config += `
  desktop: {
    framework: "wails",
    port: 5174,
  },`
	}
	if opts.ShouldIncludeDocs() {
		config += `
  docs: {
    port: 3002,
  },`
	}

	config += `
};
`
	return config
}

func readmeFile(opts Options) string {
	return fmt.Sprintf(`# %s

Built with [Grit](https://gritframework.dev) — Go + React. Built with Grit.

## Quick Start

`+"```bash"+`
# 1. Install Air for Go hot reloading
go install github.com/air-verse/air@latest

# 2. Start infrastructure (PostgreSQL, Redis, MinIO, Mailhog)
docker compose up -d

# 3. Install frontend dependencies
pnpm install

# 4. Start all services (API auto-reloads on file changes)
pnpm dev
`+"```"+`

## Project Structure

`+"```"+`
%s/
├── apps/
│   ├── api/          # Go backend (Gin + GORM)
│   ├── web/          # Next.js frontend
│   └── admin/        # Next.js admin panel
├── packages/
│   └── shared/       # Shared types, schemas, constants
├── docker-compose.yml
└── turbo.json
`+"```"+`

## Services

| Service       | URL                          |
|---------------|------------------------------|
| API           | http://localhost:8080         |
| GORM Studio   | http://localhost:8080/studio  |
| Web App       | http://localhost:3000         |
| Admin Panel   | http://localhost:3001         |
| PostgreSQL    | localhost:5432               |
| Redis         | localhost:6379               |
| MinIO Console | http://localhost:9001         |
| Mailhog       | http://localhost:8025         |

## Development

`+"```bash"+`
# Run Go API with hot reload
cd apps/api && air

# Run Next.js web app
cd apps/web && pnpm dev

# Run admin panel
cd apps/admin && pnpm dev

# Run all services via Turborepo
pnpm dev
`+"```"+`

## No Docker? No Problem

If you can't run Docker, use cloud services instead:

`+"```bash"+`
cp .env.cloud.example .env
`+"```"+`

Then fill in your keys for:
- **[Neon](https://neon.tech)** — PostgreSQL (free tier)
- **[Upstash](https://upstash.com)** — Redis (free tier)
- **[Cloudflare R2](https://dash.cloudflare.com)** — File storage (free tier)
- **[Resend](https://resend.com)** — Email (free tier)

No Docker needed — just your API keys and `+"``"+`go run`+"``"+`.

## Tech Stack

- **Backend:** Go + Gin + GORM
- **Frontend:** Next.js 14+ (App Router) + React + TypeScript
- **Styling:** Tailwind CSS + shadcn/ui
- **Database:** PostgreSQL
- **Cache:** Redis
- **Monorepo:** Turborepo + pnpm
- **Validation:** Zod (shared schemas)
- **Data Fetching:** React Query (TanStack Query)

---

*Built with Grit v%s*
`, opts.ProjectName, opts.ProjectName, opts.Version)
}

func gritJSON(opts Options) string {
	return fmt.Sprintf(`{
  "architecture": "%s",
  "frontend": "%s",
  "version": "%s"
}
`, string(opts.Architecture), string(opts.Frontend), opts.Version)
}
