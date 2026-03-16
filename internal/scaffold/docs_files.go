package scaffold

import (
	"fmt"
	"path/filepath"
)

func writeDocsFiles(root string, opts Options) error {
	docsRoot := filepath.Join(root, "apps", "docs")

	files := map[string]string{
		// Config
		filepath.Join(docsRoot, "package.json"):      docsPackageJSON(opts),
		filepath.Join(docsRoot, "tsconfig.json"):      docsTSConfig(),
		filepath.Join(docsRoot, "next.config.mjs"):    docsNextConfig(),
		filepath.Join(docsRoot, "tailwind.config.js"): docsTailwindConfig(),
		filepath.Join(docsRoot, "postcss.config.mjs"): docsPostCSSConfig(),
		filepath.Join(docsRoot, "source.config.ts"):   docsSourceConfig(),

		// App
		filepath.Join(docsRoot, "app", "source.ts"):                       docsAppSource(),
		filepath.Join(docsRoot, "app", "global.css"):                      docsGlobalCSS(),
		filepath.Join(docsRoot, "app", "layout.tsx"):                      docsRootLayout(opts),
		filepath.Join(docsRoot, "app", "page.tsx"):                        docsHomePage(),
		filepath.Join(docsRoot, "app", "api", "search", "route.ts"):       docsSearchRoute(),
		filepath.Join(docsRoot, "app", "docs", "layout.tsx"):              docsDocsLayout(),
		filepath.Join(docsRoot, "app", "docs", "[[...slug]]", "page.tsx"): docsSlugPage(),

		// Content — Core
		filepath.Join(docsRoot, "content", "docs", "index.mdx"):           docsContentIndex(opts),
		filepath.Join(docsRoot, "content", "docs", "getting-started.mdx"): docsContentGettingStarted(opts),
		filepath.Join(docsRoot, "content", "docs", "meta.json"):           docsContentMeta(),

		// Content — CLI
		filepath.Join(docsRoot, "content", "docs", "cli", "commands.mdx"): docsContentCLI(),
		filepath.Join(docsRoot, "content", "docs", "cli", "meta.json"):    docsContentCLIMeta(),

		// Content — API
		filepath.Join(docsRoot, "content", "docs", "api", "authentication.mdx"):        docsContentAuth(),
		filepath.Join(docsRoot, "content", "docs", "api", "migrations-and-seeding.mdx"): docsContentMigrationsSeeding(),
		filepath.Join(docsRoot, "content", "docs", "api", "meta.json"):                  docsContentAPIMeta(),

		// Content — Admin
		filepath.Join(docsRoot, "content", "docs", "admin", "overview.mdx"):  docsContentAdminOverview(),
		filepath.Join(docsRoot, "content", "docs", "admin", "resources.mdx"): docsContentAdminResources(),
		filepath.Join(docsRoot, "content", "docs", "admin", "meta.json"):     docsContentAdminMeta(),

		// Content — Batteries
		filepath.Join(docsRoot, "content", "docs", "batteries", "overview.mdx"): docsContentBatteries(),
		filepath.Join(docsRoot, "content", "docs", "batteries", "meta.json"):    docsContentBatteriesMeta(),

		// Content — About
		filepath.Join(docsRoot, "content", "docs", "about.mdx"): docsContentAbout(),
	}

	for path, content := range files {
		if err := writeFile(path, content); err != nil {
			return fmt.Errorf("writing %s: %w", path, err)
		}
	}

	return nil
}

func docsPackageJSON(opts Options) string {
	return fmt.Sprintf(`{
  "name": "%s-docs",
  "version": "1.0.0",
  "private": true,
  "scripts": {
    "dev": "next dev --port 3002",
    "build": "next build",
    "start": "next start --port 3002"
  },
  "dependencies": {
    "next": "^15.1.0",
    "react": "^19.0.0",
    "react-dom": "^19.0.0",
    "fumadocs-core": "^14.0.0",
    "fumadocs-ui": "^14.0.0",
    "fumadocs-mdx": "^11.0.0"
  },
  "devDependencies": {
    "@types/node": "^22.0.0",
    "@types/react": "^19.0.0",
    "@types/react-dom": "^19.0.0",
    "typescript": "^5.7.0",
    "tailwindcss": "^3.4.0",
    "postcss": "^8.4.0",
    "autoprefixer": "^10.4.0"
  }
}
`, opts.ProjectName)
}

func docsTSConfig() string {
	return `{
  "compilerOptions": {
    "target": "ES2017",
    "lib": ["dom", "dom.iterable", "esnext"],
    "allowJs": true,
    "skipLibCheck": true,
    "strict": true,
    "noEmit": true,
    "esModuleInterop": true,
    "module": "esnext",
    "moduleResolution": "bundler",
    "resolveJsonModule": true,
    "isolatedModules": true,
    "jsx": "preserve",
    "incremental": true,
    "plugins": [{ "name": "next" }],
    "paths": {
      "@/*": ["./*"]
    }
  },
  "include": ["next-env.d.ts", "**/*.ts", "**/*.tsx", "**/*.mdx", ".source/**/*.ts"],
  "exclude": ["node_modules"]
}
`
}

func docsNextConfig() string {
	return `import { createMDX } from "fumadocs-mdx/next";

const withMDX = createMDX();

/** @type {import('next').NextConfig} */
const config = {
  reactStrictMode: true,
};

export default withMDX(config);
`
}

func docsTailwindConfig() string {
	return `const { createPreset } = require("fumadocs-ui/tailwind-plugin");

/** @type {import('tailwindcss').Config} */
module.exports = {
  darkMode: "class",
  content: [
    "./components/**/*.{ts,tsx}",
    "./app/**/*.{ts,tsx}",
    "./content/**/*.{md,mdx}",
    "./node_modules/fumadocs-ui/dist/**/*.js",
  ],
  presets: [
    createPreset({
      preset: "ocean",
    }),
  ],
  theme: {
    extend: {
      colors: {
        accent: {
          DEFAULT: "#6c5ce7",
          foreground: "#ffffff",
        },
      },
    },
  },
};
`
}

func docsSourceConfig() string {
	return `import { defineDocs, defineConfig } from "fumadocs-mdx/config";

export const { docs, meta } = defineDocs({
  dir: "content/docs",
});

export default defineConfig();
`
}

func docsPostCSSConfig() string {
	return `export default {
  plugins: {
    tailwindcss: {},
    autoprefixer: {},
  },
};
`
}

func docsAppSource() string {
	return `import { docs, meta } from "@/.source";
import { createMDXSource } from "fumadocs-mdx";
import { loader } from "fumadocs-core/source";

export const source = loader({
  baseUrl: "/docs",
  source: createMDXSource(docs, meta),
});
`
}

func docsSearchRoute() string {
	return `import { source } from "@/app/source";
import { createFromSource } from "fumadocs-core/search/server";

export const { GET } = createFromSource(source);
`
}

func docsGlobalCSS() string {
	return `@tailwind base;
@tailwind components;
@tailwind utilities;
`
}

func docsRootLayout(opts Options) string {
	return fmt.Sprintf(`import "./global.css";
import { RootProvider } from "fumadocs-ui/provider";
import type { Metadata } from "next";

export const metadata: Metadata = {
  title: {
    template: "%%s | %s Docs",
    default: "%s Documentation",
  },
  description: "Documentation for %s — Go + React. Built with Grit.",
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="en" className="dark" suppressHydrationWarning>
      <body>
        <RootProvider
          theme={{
            enabled: true,
            defaultTheme: "dark",
          }}
        >
          {children}
        </RootProvider>
      </body>
    </html>
  );
}
`, opts.ProjectName, opts.ProjectName, opts.ProjectName)
}

func docsHomePage() string {
	return `import Link from "next/link";

export default function HomePage() {
  return (
    <main className="flex min-h-screen flex-col items-center justify-center bg-fd-background text-fd-foreground">
      <div className="text-center max-w-2xl px-6">
        <h1 className="text-5xl font-bold mb-4">
          Documentation
        </h1>
        <p className="text-lg text-fd-muted-foreground mb-8">
          Everything you need to build with Grit — the full-stack Go + React framework.
        </p>
        <Link
          href="/docs"
          className="inline-flex items-center px-6 py-3 rounded-lg bg-fd-primary text-fd-primary-foreground font-medium hover:bg-fd-primary/90 transition-colors"
        >
          Get Started
        </Link>
      </div>
    </main>
  );
}
`
}

func docsDocsLayout() string {
	return `import { DocsLayout } from "fumadocs-ui/layouts/docs";
import type { ReactNode } from "react";
import { source } from "@/app/source";

export default function Layout({ children }: { children: ReactNode }) {
  return (
    <DocsLayout
      tree={source.pageTree}
      nav={{
        title: (
          <span className="flex items-center gap-2 font-bold">
            Grit Docs
            <span className="rounded-md bg-fd-primary/10 px-1.5 py-0.5 text-xs font-medium text-fd-primary">
              v0.13.0
            </span>
          </span>
        ),
      }}
    >
      {children}
    </DocsLayout>
  );
}
`
}

func docsSlugPage() string {
	return `import { source } from "@/app/source";
import {
  DocsPage,
  DocsBody,
  DocsTitle,
  DocsDescription,
} from "fumadocs-ui/page";
import { notFound } from "next/navigation";
import defaultMdxComponents from "fumadocs-ui/mdx";

export default async function Page(props: {
  params: Promise<{ slug?: string[] }>;
}) {
  const params = await props.params;
  const page = source.getPage(params.slug);
  if (!page) notFound();

  const MDX = page.data.body;

  return (
    <DocsPage toc={page.data.toc} full={page.data.full}>
      <DocsTitle>{page.data.title}</DocsTitle>
      <DocsDescription>{page.data.description}</DocsDescription>
      <DocsBody>
        <MDX components={{ ...defaultMdxComponents }} />
      </DocsBody>
    </DocsPage>
  );
}

export async function generateStaticParams() {
  return source.generateParams();
}

export async function generateMetadata(props: {
  params: Promise<{ slug?: string[] }>;
}) {
  const params = await props.params;
  const page = source.getPage(params.slug);
  if (!page) notFound();

  return {
    title: page.data.title,
    description: page.data.description,
  };
}
`
}

func docsContentIndex(opts Options) string {
	return fmt.Sprintf(`---
title: Introduction
description: Welcome to the %s documentation. Learn how to build full-stack applications with Go and React.
---

## What is Grit?

Grit is a full-stack meta-framework that fuses **Go** (Gin + GORM) with **Next.js** (React + TypeScript) in a monorepo. One command to scaffold a complete project with authentication, admin panel, database browser, and Docker setup.

## Features

- **JWT Authentication** — Register, login, refresh tokens, role-based access
- **Admin Panel** — Full auth (login, sign-up, forgot-password), resource CRUD with DataTable, FormBuilder, ViewModal, toast notifications
- **Code Generator** — `+"`grit generate resource`"+` creates Go model, handler, service, Zod schemas, TypeScript types, React hooks, and admin page
- **Landing Page** — SaaS marketing page with Framer Motion animations
- **GORM Studio** — Visual database browser at `+"`/studio`"+`
- **Batteries Included** — Redis cache, S3 storage, email (Resend), background jobs (asynq), cron scheduler, AI integration (Claude + OpenAI)
- **Shared Types** — Zod schemas + TypeScript types shared between apps
- **Docker Ready** — Dev and production Docker Compose setups
- **Upgrade Command** — `+"`grit upgrade`"+` updates existing projects to the latest scaffold templates
- **Self-Update** — `+"`grit update`"+` updates the CLI itself to the latest release

## Tech Stack

| Layer | Technology |
|-------|-----------|
| Backend | Go + Gin + GORM |
| Frontend | Next.js (App Router) + React |
| Styling | Tailwind CSS + shadcn/ui |
| Database | PostgreSQL |
| Cache | Redis |
| Queue | asynq (Redis-based) |
| Storage | S3-compatible (MinIO, R2, B2) |
| Email | Resend |
| Validation | Zod |
| Data Fetching | React Query (TanStack Query) |
| Monorepo | Turborepo + pnpm |

## Next Steps

- [Getting Started](/docs/getting-started) — Install Grit and create your first project
- [CLI Commands](/docs/cli/commands) — All available CLI commands
- [Admin Panel](/docs/admin/overview) — Dashboard, auth, and resource management
- [Authentication](/docs/api/authentication) — How the JWT auth system works
- [Batteries](/docs/batteries/overview) — Cache, storage, email, jobs, cron, and AI
`, opts.ProjectName)
}

func docsContentGettingStarted(opts Options) string {
	return fmt.Sprintf(`---
title: Getting Started
description: Install Grit and create your first full-stack project in minutes.
---

**Current version: v0.13.0**

## Installation

Install the Grit CLI globally:

`+"```bash"+`
go install github.com/MUKE-coder/grit/v2/cmd/grit@latest
`+"```"+`

Verify the installation:

`+"```bash"+`
grit version
# grit version 0.13.0
`+"```"+`

## Quick Start

`+"```bash"+`
grit new %s            # Scaffold a new project
cd %s
docker compose up -d   # Start PostgreSQL, Redis, MinIO, Mailhog
pnpm install           # Install frontend dependencies
pnpm dev               # Start all services
`+"```"+`

## Create a New Project

`+"```bash"+`
grit new %s
`+"```"+`

This creates a full monorepo with:

- **Go API** — JWT auth, user management, GORM Studio, and batteries (cache, storage, email, jobs, cron, AI)
- **Web app** — SaaS marketing landing page with Framer Motion animations, linking to the admin panel for login/sign-up
- **Admin panel** — Full authentication (login, sign-up, forgot-password), dashboard, resource management with DataTable, FormBuilder, ViewModal, toast notifications
- **Shared package** — Zod schemas and TypeScript types shared between apps
- **Docker Compose** — PostgreSQL, Redis, MinIO, and Mailhog

## Start Development

`+"```bash"+`
# Start infrastructure
cd %s
docker compose up -d

# Install frontend dependencies
pnpm install

# Start all services
pnpm dev
`+"```"+`

## Available Services

| Service | URL |
|---------|-----|
| Web (Landing Page) | http://localhost:3000 |
| Admin Panel | http://localhost:3001 |
| Go API | http://localhost:8080 |
| API Docs (gin-docs) | http://localhost:8080/docs |
| Pulse (Observability) | http://localhost:8080/pulse/ui/ |
| GORM Studio | http://localhost:8080/studio |
| Mailhog | http://localhost:8025 |
| MinIO Console | http://localhost:9001 |

## Project Structure

`+"```"+`
%s/
├── apps/
│   ├── api/          # Go backend (Gin + GORM)
│   ├── web/          # Landing page (Next.js)
│   └── admin/        # Admin panel (Next.js)
├── packages/
│   └── shared/       # Zod schemas, TS types, constants
├── docker-compose.yml
├── turbo.json
└── pnpm-workspace.yaml
`+"```"+`

## CLI Flags

`+"```bash"+`
grit new myapp            # Full monorepo (api + web + admin + shared)
grit new myapp --api      # Go API only
grit new myapp --expo     # Full monorepo + Expo mobile app
grit new myapp --mobile   # API + Expo mobile app only
grit new myapp --full     # Everything + documentation site
`+"```"+`

## Generate a Resource

Once your project is set up, generate full-stack CRUD resources:

`+"```bash"+`
grit generate resource Post --fields "title:string,content:text,published:bool"
`+"```"+`

This creates: Go model, handler, service, Zod schemas, TypeScript types, React Query hooks, and an admin resource page — all wired up automatically.

## Upgrade an Existing Project

After updating the CLI, upgrade your project's scaffold files to match:

`+"```bash"+`
# Update the CLI first
go install github.com/MUKE-coder/grit/v2/cmd/grit@latest

# Then upgrade your project
cd %s
grit upgrade
`+"```"+`

This regenerates admin components, web landing page, Docker configs, and documentation while preserving your resource definitions, API code, and environment variables.

Use `+"`--force`"+` to overwrite all files without prompting:

`+"```bash"+`
grit upgrade --force
`+"```"+`

After upgrading, run `+"`grit sync`"+` to regenerate TypeScript types from your Go models.

## No Docker?

If you can't run Docker, use cloud services:

`+"```bash"+`
cp .env.cloud.example .env
`+"```"+`

Fill in your keys for [Neon](https://neon.tech) (Postgres), [Upstash](https://upstash.com) (Redis), [Cloudflare R2](https://dash.cloudflare.com) (storage), and [Resend](https://resend.com) (email).
`, opts.ProjectName, opts.ProjectName, opts.ProjectName, opts.ProjectName, opts.ProjectName, opts.ProjectName)
}

func docsContentAuth() string {
	return `---
title: Authentication
description: JWT-based authentication with access and refresh tokens.
---

## Overview

Grit uses **JWT (JSON Web Tokens)** for authentication with a dual-token system:

- **Access Token** — Short-lived (15 minutes), sent with every request
- **Refresh Token** — Long-lived (7 days), used to obtain new access tokens

## Endpoints

### Register

` + "```bash" + `
POST /api/auth/register
Content-Type: application/json

{
  "name": "John Doe",
  "email": "john@example.com",
  "password": "securepassword"
}
` + "```" + `

### Login

` + "```bash" + `
POST /api/auth/login
Content-Type: application/json

{
  "email": "john@example.com",
  "password": "securepassword"
}
` + "```" + `

**Response:**

` + "```json" + `
{
  "data": {
    "user": {
      "id": "uuid",
      "name": "John Doe",
      "email": "john@example.com",
      "role": "user"
    },
    "access_token": "eyJhbG...",
    "refresh_token": "eyJhbG..."
  },
  "message": "Login successful"
}
` + "```" + `

### Refresh Token

` + "```bash" + `
POST /api/auth/refresh
Content-Type: application/json

{
  "refresh_token": "eyJhbG..."
}
` + "```" + `

### Get Current User

` + "```bash" + `
GET /api/auth/me
Authorization: Bearer <access_token>
` + "```" + `

### Logout

` + "```bash" + `
POST /api/auth/logout
Authorization: Bearer <access_token>
` + "```" + `

## Using Authentication in Requests

Include the access token in the ` + "`Authorization`" + ` header:

` + "```bash" + `
curl -H "Authorization: Bearer <access_token>" http://localhost:8080/api/users
` + "```" + `

## Role-Based Access

Users have a ` + "`role`" + ` field that can be ` + "`user`" + ` or ` + "`admin`" + `. Protected routes use the ` + "`RequireRole`" + ` middleware:

` + "```go" + `
// Only admins can access this route
admin := r.Group("/admin")
admin.Use(middleware.RequireAuth(db))
admin.Use(middleware.RequireRole("admin"))
` + "```" + `
`
}

func docsContentMeta() string {
	return `{
  "title": "Documentation",
  "pages": [
    "index",
    "getting-started",
    "...cli",
    "...admin",
    "...api",
    "...batteries",
    "about"
  ]
}
`
}

func docsContentCLIMeta() string {
	return `{
  "title": "CLI",
  "pages": ["commands"]
}
`
}

func docsContentAPIMeta() string {
	return `{
  "title": "API",
  "pages": ["authentication", "migrations-and-seeding"]
}
`
}

func docsContentAdminMeta() string {
	return `{
  "title": "Admin Panel",
  "pages": ["overview", "resources"]
}
`
}

func docsContentBatteriesMeta() string {
	return `{
  "title": "Batteries",
  "pages": ["overview"]
}
`
}

func docsContentCLI() string {
	return `---
title: CLI Commands
description: Complete reference for all Grit CLI commands.
---

## grit new

Scaffold a new Grit project.

` + "```bash" + `
grit new <project-name> [flags]
` + "```" + `

**Flags:**

| Flag | Description |
|------|-------------|
| --api | Scaffold Go API only (no frontend) |
| --expo | Include Expo React Native app |
| --mobile | API + Expo app only |
| --full | Everything including docs site |

**Examples:**

` + "```bash" + `
grit new my-saas              # Full monorepo
grit new my-saas --api        # Backend only
grit new my-saas --full       # With docs site
` + "```" + `

## grit generate resource

Generate a full-stack CRUD resource with a single command.

` + "```bash" + `
grit generate resource <Name> [flags]
` + "```" + `

**Aliases:** grit g resource

**Flags:**

| Flag | Description |
|------|-------------|
| --fields | Inline field definitions |
| --from | YAML file with field definitions |
| -i, --interactive | Interactive field builder |

**Generated files:**

- Go model (internal/models/name.go)
- Go service (internal/services/name_service.go)
- Go handler (internal/handlers/name_handler.go)
- Zod schemas (packages/shared/schemas/name.ts)
- TypeScript types (packages/shared/types/name.ts)
- React Query hooks (apps/admin/hooks/use-plural.ts)
- Admin resource definition (apps/admin/resources/plural.ts)
- Admin resource page (apps/admin/app/(dashboard)/resources/plural/page.tsx)

**Field types:**

| Type | Go Type | Example |
|------|---------|---------|
| string | string | title:string |
| text | string | content:text |
| int | int | views:int |
| float | float64 | price:float |
| bool | bool | published:bool |
| date | time.Time | due_date:date |
| email | string | contact:email |
| url | string | website:url |
| image | string | avatar:image |
| select | string | status:select:draft,published,archived |

**Field modifiers:**

Append modifiers after the type with a colon:

` + "```bash" + `
grit g resource Post --fields "title:string:required,slug:string:unique,views:int:default=0"
` + "```" + `

| Modifier | Description |
|----------|-------------|
| required | NOT NULL constraint |
| unique | UNIQUE constraint |
| optional | Nullable field |
| default=X | Default value |

## grit sync

Sync Go model types to TypeScript types and Zod schemas.

` + "```bash" + `
grit sync
` + "```" + `

Parses all Go model files in apps/api/internal/models/ and regenerates the corresponding TypeScript types and Zod schemas in packages/shared/.

## grit upgrade

Upgrade an existing Grit project to the latest scaffold templates.

` + "```bash" + `
grit upgrade [flags]
` + "```" + `

**Flags:**

| Flag | Description |
|------|-------------|
| -f, --force | Overwrite all files without prompting |

**What gets updated:**

- Admin panel components (layout, tables, forms, widgets, resource UI)
- Web landing page
- Docker and root configuration files
- Documentation (if present)
- Expo app (if present)

**What is preserved:**

- Your resource definitions (resources/*.ts)
- API handlers, services, and models
- Environment variables (.env)
- Custom code you've written

Run grit sync after upgrading to regenerate TypeScript types.

## grit migrate

Run database migrations using GORM AutoMigrate.

` + "```bash" + `
grit migrate [flags]
` + "```" + `

**Flags:**

| Flag | Description |
|------|-------------|
| --fresh | Drop all tables before migrating |

**Examples:**

` + "```bash" + `
grit migrate            # Run migrations
grit migrate --fresh    # Drop all tables and re-migrate
` + "```" + `

The migrate command connects to the database (using DATABASE_URL from .env) and runs GORM AutoMigrate for all registered models. Use --fresh to drop all tables first — useful for resetting during development.

## grit seed

Populate the database with initial data.

` + "```bash" + `
grit seed
` + "```" + `

Seeds the database with:

- **Admin user** — admin@example.com / password (role: admin)
- **Demo users** — Sample accounts for testing

The seeder is idempotent — it skips records that already exist. Add your own seeders in apps/api/internal/database/seed.go.

## grit update

Update the Grit CLI itself to the latest version. This removes the current binary and installs the newest release from GitHub.

` + "```bash" + `
grit update
` + "```" + `

This runs ` + "`go install github.com/MUKE-coder/grit/v2/cmd/grit@latest`" + ` under the hood, replacing the old binary with the latest version. Run ` + "`grit version`" + ` afterwards to verify.

> **Note:** ` + "`grit update`" + ` updates the **CLI tool**. To update your **project's scaffold files**, use ` + "`grit upgrade`" + ` instead.

## grit version

Print the installed Grit CLI version.

` + "```bash" + `
grit version
` + "```" + `
`
}

func docsContentMigrationsSeeding() string {
	return `---
title: Migrations & Seeding
description: Database migrations with GORM AutoMigrate and seeding with initial data.
---

## Overview

Grit uses **GORM AutoMigrate** for database schema management. Migrations run automatically on API startup and can also be triggered manually via the CLI. Seeding populates the database with initial data for development and production.

## Migrations

### How It Works

GORM AutoMigrate creates tables, adds missing columns, and creates indexes based on your Go model definitions. It does **not** delete columns or drop tables — it's additive only.

All models are registered in apps/api/internal/models/user.go via the AutoMigrate function:

` + "```go" + `
func AutoMigrate(db *gorm.DB) error {
    models := []interface{}{
        &User{},
        &Upload{},
        // grit:models  <-- new models injected here by grit generate
    }

    for _, model := range models {
        if err := db.AutoMigrate(model); err != nil {
            log.Printf("Warning: migration error for %T: %v (skipping)", model, err)
        }
    }

    return nil
}
` + "```" + `

### Running Migrations

Migrations run automatically when the API starts. You can also run them manually:

` + "```bash" + `
# Run migrations
grit migrate

# Drop all tables and re-migrate (development only!)
grit migrate --fresh
` + "```" + `

### Fresh Migrations

The --fresh flag drops **all tables** in the database and re-creates them from scratch. This is useful during development when you've made breaking schema changes.

**Warning:** This destroys all data. Never use --fresh in production.

` + "```bash" + `
# Reset everything and re-seed
grit migrate --fresh
grit seed
` + "```" + `

### Adding New Models

When you run grit generate resource, it automatically:

1. Creates a new Go model file
2. Injects the model into AutoMigrate via the // grit:models marker
3. Migrations run on next API startup or grit migrate

To add a model manually:

` + "```go" + `
// apps/api/internal/models/post.go
type Post struct {
    ID        uint           ` + "`" + `gorm:"primarykey" json:"id"` + "`" + `
    Title     string         ` + "`" + `gorm:"size:255;not null" json:"title"` + "`" + `
    Content   string         ` + "`" + `gorm:"type:text" json:"content"` + "`" + `
    Published bool           ` + "`" + `gorm:"default:false" json:"published"` + "`" + `
    CreatedAt time.Time      ` + "`" + `json:"created_at"` + "`" + `
    UpdatedAt time.Time      ` + "`" + `json:"updated_at"` + "`" + `
    DeletedAt gorm.DeletedAt ` + "`" + `gorm:"index" json:"-"` + "`" + `
}
` + "```" + `

Then add it to AutoMigrate in models/user.go above the // grit:models marker.

## Seeding

### Running Seeders

` + "```bash" + `
grit seed
` + "```" + `

This creates:

| Account | Email | Password | Role |
|---------|-------|----------|------|
| Admin | admin@example.com | password | admin |
| Jane Cooper | jane@example.com | password | editor |
| Robert Fox | robert@example.com | password | user |
| Emily Davis | emily@example.com | password | user |
| Michael Chen | michael@example.com | password | user (inactive) |

### How Seeders Work

Seeders live in apps/api/internal/database/seed.go. The main Seed function calls individual seeder functions:

` + "```go" + `
func Seed(db *gorm.DB) error {
    if err := seedAdminUser(db); err != nil {
        return err
    }
    if err := seedDemoUsers(db); err != nil {
        return err
    }
    // grit:seeders  <-- add your own seeders above this marker
    return nil
}
` + "```" + `

### Adding Custom Seeders

Add new seeder functions in seed.go:

` + "```go" + `
func seedCategories(db *gorm.DB) error {
    categories := []models.Category{
        {Name: "Technology", Slug: "technology"},
        {Name: "Business", Slug: "business"},
        {Name: "Design", Slug: "design"},
    }

    for _, c := range categories {
        var count int64
        db.Model(&models.Category{}).Where("slug = ?", c.Slug).Count(&count)
        if count > 0 {
            continue
        }
        if err := db.Create(&c).Error; err != nil {
            return err
        }
    }

    return nil
}
` + "```" + `

Then call it from the Seed function:

` + "```go" + `
if err := seedCategories(db); err != nil {
    return fmt.Errorf("seeding categories: %w", err)
}
` + "```" + `

### Idempotent Seeders

Seeders check if records already exist before creating them. This means you can run grit seed multiple times without creating duplicates.

## Common Workflow

` + "```bash" + `
# Initial setup
docker compose up -d          # Start PostgreSQL
grit migrate                  # Create tables
grit seed                     # Populate with test data

# After adding a new resource
grit generate resource Post --fields "title:string,content:text"
grit migrate                  # Apply new table

# Reset during development
grit migrate --fresh          # Drop everything
grit seed                     # Re-populate
` + "```" + `
`
}

func docsContentAdminOverview() string {
	return `---
title: Admin Panel Overview
description: A complete admin dashboard with authentication, resource management, and system pages.
---

## Overview

The Grit admin panel is a full-featured Next.js application at apps/admin/. It includes authentication, a dashboard, resource CRUD management, and system administration pages.

## Architecture

The admin panel uses **Next.js App Router** with route groups to separate authenticated and public pages:

` + "```" + `
apps/admin/app/
+-- layout.tsx                        # Root layout (Providers, no sidebar)
+-- page.tsx                          # Redirect: /dashboard or /login
+-- (auth)/
|   +-- login/page.tsx                # Login page
|   +-- sign-up/page.tsx              # Registration page
|   +-- forgot-password/page.tsx      # Password reset
+-- (dashboard)/
    +-- layout.tsx                    # AdminLayout (sidebar + navbar)
    +-- dashboard/page.tsx            # Dashboard home
    +-- resources/users/page.tsx      # Users resource page
    +-- system/
        +-- jobs/page.tsx             # Background jobs
        +-- files/page.tsx            # File browser
        +-- cron/page.tsx             # Cron scheduler
        +-- mail/page.tsx             # Mail preview
` + "```" + `

## Authentication Pages

The admin panel includes split-screen auth pages with a branded left panel and form on the right:

- **Login** (/login) — Email/password with password visibility toggle
- **Sign Up** (/sign-up) — Name, email, password, confirm password
- **Forgot Password** (/forgot-password) — Email-based password reset

Auth pages use the (auth) route group and do NOT include the sidebar/navbar layout.

## Dashboard

The dashboard page (/dashboard) displays:

- Welcome greeting with the logged-in user's name
- Resource cards showing all registered resources with icons and record counts
- Quick links to system pages and GORM Studio

## Layout

The admin layout includes:

- **Sidebar** — Collapsible with Lucide icons, dashboard link, dynamic resource navigation (read from the resource registry), system pages section
- **Navbar** — Breadcrumbs, theme toggle (dark/light), user menu with logout
- **Theme** — Dark mode by default with a premium design system

## Toast Notifications

All CRUD operations display toast notifications via Sonner:

- **Create** — "Created successfully"
- **Update** — "Updated successfully"
- **Delete** — "Deleted successfully"
- **Error** — "Something went wrong"

## System Pages

Under /system/*, the admin includes:

- **Jobs** — Background job dashboard (asynq)
- **Files** — S3/MinIO file browser
- **Cron** — Scheduled task viewer
- **Mail** — Email template preview
`
}

func docsContentAdminResources() string {
	return `---
title: Resource Management
description: How resources work in the Grit admin panel — DataTable, FormBuilder, ViewModal, and more.
---

## Overview

Resources are the core building block of the admin panel. Each resource represents a database entity (e.g., Users, Posts, Products) and is defined using defineResource().

## Defining a Resource

Create a resource definition in resources/plural.ts:

` + "```typescript" + `
import { defineResource } from "@/lib/resource";

export const postsResource = defineResource({
  name: "Post",
  plural: "posts",
  apiRoute: "/api/posts",
  icon: "FileText",
  columns: [
    { key: "title", label: "Title", sortable: true, searchable: true },
    { key: "status", label: "Status", format: "badge", sortable: true },
    { key: "createdAt", label: "Created", format: "date", sortable: true },
  ],
  fields: [
    { key: "title", label: "Title", type: "text", required: true },
    { key: "content", label: "Content", type: "textarea" },
    { key: "status", label: "Status", type: "select", options: [
      { label: "Draft", value: "draft" },
      { label: "Published", value: "published" },
    ]},
  ],
  defaultSort: { key: "createdAt", direction: "desc" },
  actions: ["create", "edit", "delete", "view"],
});
` + "```" + `

## Resource Page

Each resource gets a thin page wrapper in app/(dashboard)/resources/plural/page.tsx:

` + "```typescript" + `
"use client";
import { ResourcePage } from "@/components/resource/ResourcePage";
import { postsResource } from "@/resources/posts";

export default function PostsPage() {
  return <ResourcePage resource={postsResource} />;
}
` + "```" + `

The ResourcePage component handles everything: DataTable, search, filters, pagination, create/edit modals, view modal, and delete confirmation.

## DataTable

The DataTable component supports:

- **Sorting** — Click column headers to sort
- **Search** — Full-text search across searchable columns
- **Filters** — Column-specific filters
- **Selection** — Checkbox selection with bulk actions (bulk delete)
- **Pagination** — Page size selector and page navigation
- **Actions** — View (eye icon), Edit (pencil icon), Delete (trash icon) per row

## FormBuilder and FormModal

The FormBuilder renders forms dynamically based on the resource fields definition:

| Field Type | Component | Description |
|-----------|-----------|-------------|
| text | TextField | Single-line text input |
| textarea | TextAreaField | Multi-line text input |
| number | NumberField | Numeric input |
| select | SelectField | Dropdown selection |
| date | DateField | Date picker |
| toggle | ToggleField | Boolean toggle switch |
| checkbox | CheckboxField | Checkbox input |
| radio | RadioField | Radio button group |
| image | Dropzone | Image upload with drag-and-drop |

The FormModal wraps FormBuilder in a slide-over modal for create/edit operations.

## ViewModal

The ViewModal displays a read-only view of a resource record:

- Two-column layout for field labels and values
- Formatted values based on column format (badge, date, boolean, image)
- "Edit" button to switch to edit mode
- Clean, professional design

## Confirm Modal

Delete operations use a custom confirmation modal instead of the browser confirm() dialog:

- Danger variant with red confirm button
- Custom title and description
- Loading state during deletion
- Works for both single-item and bulk delete

## Column Formats

Available column formats for the DataTable:

| Format | Description |
|--------|-------------|
| text | Plain text (default) |
| badge | Colored badge |
| date | Formatted date |
| boolean | Check/X icon |
| image | Thumbnail image |
| email | Email link |
| url | External link |
| number | Formatted number |
| currency | Currency value |

## Generating Resources

Use the CLI to generate a complete resource:

` + "```bash" + `
grit generate resource Post --fields "title:string,content:text,status:select:draft,published"
` + "```" + `

This creates all Go backend code, shared types, React hooks, resource definition, and admin page — fully wired up and ready to use.
`
}

func docsContentBatteries() string {
	return `---
title: Batteries Overview
description: Built-in services for cache, storage, email, background jobs, cron, and AI.
---

## Overview

Grit comes with "batteries included" — production-ready services that are scaffolded with every project. All services follow a graceful degradation pattern: if the backing service (Redis, MinIO, etc.) is unavailable, the app logs a warning but continues running.

## Redis Cache

**Files:** internal/cache/, internal/middleware/cache.go

A Redis-based caching service with HTTP middleware:

` + "```go" + `
// Cache a value
cache.Set(ctx, "key", value, 5*time.Minute)

// Get a cached value
val, err := cache.Get(ctx, "key")

// Delete a cached value
cache.Delete(ctx, "key")
` + "```" + `

The cache middleware automatically caches GET responses:

` + "```go" + `
r.GET("/api/posts", middleware.Cache(cacheService, 5*time.Minute), handler.GetPosts)
` + "```" + `

## S3 File Storage

**Files:** internal/storage/, internal/handlers/upload_handler.go

S3-compatible file storage with image processing:

- Upload files to MinIO (dev) or S3/R2/B2 (production)
- Automatic image resizing with configurable dimensions
- Secure pre-signed URLs for downloads
- Admin file browser at /system/files

Upload endpoint: POST /api/uploads (multipart form data)

## Email (Resend)

**Files:** internal/mail/

Email service powered by Resend with 4 HTML templates:

- **Welcome** — New user registration
- **Password Reset** — Reset link email
- **Notification** — Generic notifications
- **Invoice** — Payment/billing emails

` + "```go" + `
mailer.Send(ctx, mail.Message{
    To:       "user@example.com",
    Subject:  "Welcome!",
    Template: "welcome",
    Data:     map[string]string{"name": "John"},
})
` + "```" + `

Preview emails at /system/mail in the admin panel.

## Background Jobs (asynq)

**Files:** internal/jobs/

Redis-based background job processing:

- **Email jobs** — Send emails asynchronously
- **Image jobs** — Process uploaded images
- **Cleanup jobs** — Remove expired data

` + "```go" + `
// Enqueue a job
client.Enqueue(jobs.NewEmailJob("user@example.com", "welcome"))
` + "```" + `

Monitor jobs at /system/jobs in the admin panel.

## Cron Scheduler

**Files:** internal/cron/

Scheduled task execution built on asynq:

- Define cron tasks with standard cron expressions
- View scheduled tasks at /system/cron in the admin panel

## AI Integration

**Files:** internal/ai/

Dual AI provider support:

- **Claude** (Anthropic) — via the Anthropic API
- **OpenAI** — via the OpenAI API

Features:
- Text generation with streaming support
- Configurable via environment variables
- Admin handler for AI operations

` + "```bash" + `
# .env configuration (Vercel AI Gateway)
AI_GATEWAY_API_KEY=your-key        # One key for all models
AI_GATEWAY_MODEL=anthropic/claude-sonnet-4-6  # provider/model format
` + "```" + `

## Configuration

All batteries are configured via environment variables in .env:

` + "```bash" + `
# Redis
REDIS_URL=localhost:6379

# S3 Storage
S3_ENDPOINT=localhost:9000
S3_ACCESS_KEY=minioadmin
S3_SECRET_KEY=minioadmin
S3_BUCKET=uploads
S3_REGION=us-east-1

# Email
RESEND_API_KEY=re_...
MAIL_FROM=noreply@example.com

# AI (Vercel AI Gateway — one key, hundreds of models)
AI_GATEWAY_API_KEY=your-key
AI_GATEWAY_MODEL=anthropic/claude-sonnet-4-6
` + "```" + `

## Docker Services

All backing services are included in docker-compose.yml:

| Service | Port | Purpose |
|---------|------|---------|
| PostgreSQL | 5432 | Database |
| Redis | 6379 | Cache + Queue |
| MinIO | 9000/9001 | File storage |
| Mailhog | 8025 | Email testing |
`
}

func docsContentAbout() string {
	return `---
title: About the Creator
description: Meet the developer behind Grit.
---

## Muke Johnbaptist (JB)

**Fullstack Developer | Tech Entrepreneur | Educator**

Grit is created by **Muke Johnbaptist** (JB), a self-taught tech entrepreneur, educator, and founder based in **Kampala, Uganda**. JB is passionate about building digital solutions and empowering others to succeed in tech.

## Background

JB is entirely self-taught with no formal computer science degree. He transitioned from learning to code independently into building a thriving business that serves clients globally. His philosophy: **you don't need a computer science degree to build a successful tech career** — anyone can learn, anyone can build, and anyone can succeed.

## What He Does

- **Founder & CEO** of [Desishub Technologies](https://desishub.com) — a company that develops websites, mobile applications, and AI-powered business solutions across Africa and beyond
- **Founder** of Nextjs Academy — an online learning platform with 9+ courses
- **YouTube Creator** — 14,600+ subscribers on the [JB Web Developer](https://youtube.com/@JBWEBDEVELOPER) channel
- **Educator** — Trained 100+ students through in-person coding bootcamps

## Achievements

- Generated **$15,000+** in revenue through online courses and digital products
- Created **9+ online courses** through Nextjs Academy
- Built multiple **enterprise-level systems** including school management platforms and member management systems
- Currently maintaining **13 active projects**

## Tech Stack

| Layer | Technologies |
|-------|-------------|
| Frontend | React, Next.js, TypeScript, Tailwind CSS, shadcn/ui |
| Backend | Node.js, Hono.js, Go, Express.js, Laravel |
| Databases | PostgreSQL, MongoDB, Redis |
| Tools | Docker, Git, Prisma ORM, GORM, Resend, React Query |

## Connect

- [GitHub](https://github.com/MUKE-coder) — @MUKE-coder
- [YouTube](https://youtube.com/@JBWEBDEVELOPER) — @JBWEBDEVELOPER
- [X (Twitter)](https://x.com/MJohnbaptist) — @MJohnbaptist
- [LinkedIn](https://linkedin.com/in/muke-johnbaptist) — muke-johnbaptist
- [TikTok](https://tiktok.com/@jbdesishub) — @jbdesishub
- [Website](https://jb.desishub.com) — jb.desishub.com

---

*Grit is built with the belief that developers deserve better tools — tools that are fast to set up, production-ready, and don't compromise on quality.*
`
}
