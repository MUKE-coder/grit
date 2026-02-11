package scaffold

import (
	"fmt"
	"path/filepath"
)

func writeRootFiles(root string, opts Options) error {
	files := map[string]string{
		filepath.Join(root, ".env"):         envFile(opts),
		filepath.Join(root, ".env.example"): envExampleFile(opts),
		filepath.Join(root, ".gitignore"):   rootGitignore(),
		filepath.Join(root, "README.md"):    readmeFile(opts),
	}

	if !opts.APIOnly {
		files[filepath.Join(root, "pnpm-workspace.yaml")] = pnpmWorkspace()
		files[filepath.Join(root, "turbo.json")] = turboJSON()
		files[filepath.Join(root, "package.json")] = rootPackageJSON(opts)
		files[filepath.Join(root, "grit.config.ts")] = gritConfig(opts)
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

# Redis
REDIS_URL=redis://localhost:6379

# S3 / MinIO
S3_ENDPOINT=http://localhost:9000
S3_ACCESS_KEY=minioadmin
S3_SECRET_KEY=minioadmin
S3_BUCKET=%s-uploads
S3_REGION=us-east-1
S3_USE_SSL=false

# Email (Resend)
RESEND_API_KEY=re_your_api_key
MAIL_FROM=noreply@%s.dev

# CORS
CORS_ORIGINS=http://localhost:3000,http://localhost:3001

# GORM Studio
GORM_STUDIO_ENABLED=true
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

# Redis — Cache and job queue
REDIS_URL=redis://localhost:6379

# S3 / MinIO — File storage
S3_ENDPOINT=http://localhost:9000    # MinIO for local dev
S3_ACCESS_KEY=minioadmin
S3_SECRET_KEY=minioadmin
S3_BUCKET=myapp-uploads
S3_REGION=us-east-1
S3_USE_SSL=false

# Email — Resend integration
RESEND_API_KEY=re_your_api_key
MAIL_FROM=noreply@myapp.dev

# CORS — Allowed frontend origins (comma-separated)
CORS_ORIGINS=http://localhost:3000,http://localhost:3001

# GORM Studio — Visual database browser
GORM_STUDIO_ENABLED=true
`
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

func rootPackageJSON(opts Options) string {
	return fmt.Sprintf(`{
  "name": "%s",
  "private": true,
  "scripts": {
    "dev": "turbo dev",
    "build": "turbo build",
    "lint": "turbo lint",
    "type-check": "turbo type-check",
    "dev:web": "cd apps/web && pnpm dev",
    "dev:admin": "cd apps/admin && pnpm dev",
    "dev:api": "cd apps/api && go run cmd/server/main.go",
    "docker:up": "docker compose up -d",
    "docker:down": "docker compose down",
    "docker:logs": "docker compose logs -f"
  },
  "devDependencies": {
    "turbo": "^2.0.0"
  },
  "packageManager": "pnpm@10.0.0"
}
`, opts.ProjectName)
}

func gritConfig(opts Options) string {
	return fmt.Sprintf(`// Grit Framework Configuration
export default {
  name: "%s",
  api: {
    port: 8080,
    prefix: "/api",
  },
  web: {
    port: 3000,
  },
  admin: {
    port: 3001,
  },
};
`, opts.ProjectName)
}

func readmeFile(opts Options) string {
	return fmt.Sprintf(`# %s

Built with [Grit](https://gritframework.dev) — Go + React. Built with Grit.

## Quick Start

`+"```bash"+`
# 1. Start infrastructure (PostgreSQL, Redis, MinIO, Mailhog)
docker compose up -d

# 2. Install frontend dependencies
pnpm install

# 3. Start all services
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

*Built with Grit v0.1.0*
`, opts.ProjectName, opts.ProjectName)
}
