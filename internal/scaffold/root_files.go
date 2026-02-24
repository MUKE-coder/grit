package scaffold

import (
	"fmt"
	"path/filepath"
)

func writeRootFiles(root string, opts Options) error {
	files := map[string]string{
		filepath.Join(root, ".env"):               envFile(opts),
		filepath.Join(root, ".env.example"):       envExampleFile(opts),
		filepath.Join(root, ".env.cloud.example"): envCloudExampleFile(opts),
		filepath.Join(root, ".gitignore"):         rootGitignore(),
		filepath.Join(root, "README.md"):          readmeFile(opts),
		filepath.Join(root, "GRIT_SKILL.md"):      gritSkillFile(opts),
	}

	if !opts.APIOnly {
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
	return fmt.Sprintf(`# %s — Environment Variables

# App
APP_NAME=%s
APP_ENV=development
APP_PORT=8080
APP_URL=http://localhost:8080

# Database
DATABASE_URL=postgres://grit:grit@localhost:5432/%s?sslmode=disable

# JWT
JWT_SECRET=your-super-secret-jwt-key-change-in-production
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

# Storage — Which provider to use: minio, r2, b2
STORAGE_DRIVER=minio

# MinIO (local development — used when STORAGE_DRIVER=minio)
MINIO_ENDPOINT=http://localhost:9000
MINIO_ACCESS_KEY=minioadmin
MINIO_SECRET_KEY=minioadmin
MINIO_BUCKET=%s-uploads
MINIO_REGION=us-east-1
MINIO_USE_SSL=false

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

# AI — Provider: claude, openai, or gemini
AI_PROVIDER=claude
AI_API_KEY=                          # sk-ant-... (Claude), sk-... (OpenAI), or AIza... (Gemini)
AI_MODEL=claude-sonnet-4-5-20250929  # Model to use

# Observability — Pulse performance monitoring dashboard
PULSE_ENABLED=true
PULSE_USERNAME=admin
PULSE_PASSWORD=pulse

# Security — Sentinel WAF, rate limiting, threat detection
SENTINEL_ENABLED=true
SENTINEL_USERNAME=admin
SENTINEL_PASSWORD=sentinel
SENTINEL_SECRET_KEY=change-me-in-production
`, opts.ProjectName, opts.ProjectName, opts.ProjectName, opts.ProjectName, opts.ProjectName)
}

func envExampleFile(opts Options) string {
	return `# App — General application settings
APP_NAME=myapp              # Application name
APP_ENV=development         # Environment: development, staging, production
APP_PORT=8080               # API server port
APP_URL=http://localhost:8080

# Database — PostgreSQL connection
DATABASE_URL=postgres://grit:grit@localhost:5432/myapp?sslmode=disable

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

# Storage — Active driver: minio, r2, or b2
STORAGE_DRIVER=minio                 # Change to "r2" or "b2" to switch providers

# MinIO — Local S3-compatible storage (default for development)
MINIO_ENDPOINT=http://localhost:9000
MINIO_ACCESS_KEY=minioadmin
MINIO_SECRET_KEY=minioadmin
MINIO_BUCKET=myapp-uploads
MINIO_REGION=us-east-1
MINIO_USE_SSL=false

# Cloudflare R2 — S3-compatible object storage
# Get credentials: Cloudflare Dashboard → R2 → Manage R2 API Tokens
R2_ENDPOINT=https://YOUR_ACCOUNT_ID.r2.cloudflarestorage.com
R2_ACCESS_KEY=                       # R2 Access Key ID
R2_SECRET_KEY=                       # R2 Secret Access Key
R2_BUCKET=myapp-uploads
R2_REGION=auto                       # Always "auto" for R2

# Backblaze B2 — S3-compatible object storage
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

# AI — Text generation (Claude, OpenAI, or Gemini)
AI_PROVIDER=claude                   # "claude", "openai", or "gemini"
AI_API_KEY=                          # sk-ant-... (Claude), sk-... (OpenAI), or AIza... (Gemini)
AI_MODEL=claude-sonnet-4-5-20250929  # Model to use

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

# ─── AI (Claude, OpenAI, or Gemini) ──────────────────
AI_PROVIDER=claude                   # "claude", "openai", or "gemini"
AI_API_KEY=your-api-key-here         # sk-ant-... (Claude), sk-... (OpenAI), or AIza... (Gemini)
AI_MODEL=claude-sonnet-4-5-20250929  # Model to use

# ─── Observability (Pulse) ────────────────────────────
PULSE_ENABLED=true
PULSE_USERNAME=admin
PULSE_PASSWORD=change-me-in-production

# ─── Security (Sentinel) ─────────────────────────────
SENTINEL_ENABLED=true
SENTINEL_USERNAME=admin
SENTINEL_PASSWORD=change-me-in-production
SENTINEL_SECRET_KEY=generate-a-random-string-here
`, opts.ProjectName, opts.ProjectName, opts.ProjectName, opts.ProjectName)
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

# Debug
*.log
npm-debug.log*
pnpm-debug.log*
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
	if opts.ShouldIncludeDocs() {
		scripts += `
    "dev:docs": "cd apps/docs && pnpm dev",`
	}

	scripts += `
    "docker:up": "docker compose up -d",
    "docker:down": "docker compose down",
    "docker:logs": "docker compose logs -f"`

	return fmt.Sprintf(`{
  "name": "%s",
  "private": true,
  "scripts": {
%s
  },
  "devDependencies": {
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

*Built with Grit v0.13.0*
`, opts.ProjectName, opts.ProjectName)
}
