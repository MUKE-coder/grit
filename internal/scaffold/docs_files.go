package scaffold

import (
	"fmt"
	"path/filepath"
)

func writeDocsFiles(root string, opts Options) error {
	docsRoot := filepath.Join(root, "apps", "docs")

	files := map[string]string{
		filepath.Join(docsRoot, "package.json"):                           docsPackageJSON(opts),
		filepath.Join(docsRoot, "tsconfig.json"):                         docsTSConfig(),
		filepath.Join(docsRoot, "next.config.mjs"):                       docsNextConfig(),
		filepath.Join(docsRoot, "tailwind.config.js"):                    docsTailwindConfig(),
		filepath.Join(docsRoot, "postcss.config.mjs"):                    docsPostCSSConfig(),
		filepath.Join(docsRoot, "source.config.ts"):                      docsSourceConfig(),
		filepath.Join(docsRoot, "app", "source.ts"):                      docsAppSource(),
		filepath.Join(docsRoot, "app", "global.css"):                     docsGlobalCSS(),
		filepath.Join(docsRoot, "app", "layout.tsx"):                     docsRootLayout(opts),
		filepath.Join(docsRoot, "app", "page.tsx"):                       docsHomePage(),
		filepath.Join(docsRoot, "app", "docs", "layout.tsx"):             docsDocsLayout(),
		filepath.Join(docsRoot, "app", "docs", "[[...slug]]", "page.tsx"): docsSlugPage(),
		filepath.Join(docsRoot, "content", "docs", "index.mdx"):          docsContentIndex(opts),
		filepath.Join(docsRoot, "content", "docs", "getting-started.mdx"): docsContentGettingStarted(opts),
		filepath.Join(docsRoot, "content", "docs", "api", "authentication.mdx"): docsContentAuth(),
		filepath.Join(docsRoot, "content", "docs", "meta.json"):          docsContentMeta(),
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
      nav={{ title: "Grit Docs" }}
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
- **User Management** — CRUD with pagination, search, sorting
- **GORM Studio** — Visual database browser at ` + "`/studio`" + `
- **Admin Panel** — Data tables, stats cards, user management
- **Shared Types** — Zod schemas + TypeScript types shared between apps
- **Docker Ready** — Dev and production Docker Compose setups

## Tech Stack

| Layer | Technology |
|-------|-----------|
| Backend | Go + Gin + GORM |
| Frontend | Next.js (App Router) + React |
| Styling | Tailwind CSS + shadcn/ui |
| Database | PostgreSQL |
| Cache | Redis |
| Validation | Zod |
| Data Fetching | React Query (TanStack Query) |
| Monorepo | Turborepo + pnpm |

## Next Steps

- [Getting Started](/docs/getting-started) — Install Grit and create your first project
- [Authentication](/docs/api/authentication) — Learn how the auth system works
`, opts.ProjectName)
}

func docsContentGettingStarted(opts Options) string {
	return fmt.Sprintf(`---
title: Getting Started
description: Install Grit and create your first full-stack project in minutes.
---

## Installation

Install the Grit CLI globally:

` + "```bash" + `
go install github.com/MUKE-coder/grit/cmd/grit@latest
` + "```" + `

## Create a New Project

` + "```bash" + `
grit new %s
` + "```" + `

This creates a full monorepo with:

- **Go API** with JWT auth, user management, and GORM Studio
- **Next.js web app** with login, register, and dashboard
- **Admin panel** with data tables and user management
- **Shared package** with Zod schemas and TypeScript types
- **Docker Compose** with PostgreSQL, Redis, MinIO, and Mailhog

## Start Development

` + "```bash" + `
# Start infrastructure
cd %s
docker compose up -d

# Install frontend dependencies
pnpm install

# Start all services
pnpm dev
` + "```" + `

## Available Services

| Service | URL |
|---------|-----|
| Web App | http://localhost:3000 |
| Admin Panel | http://localhost:3001 |
| Go API | http://localhost:8080 |
| GORM Studio | http://localhost:8080/studio |
| Mailhog | http://localhost:8025 |
| MinIO Console | http://localhost:9001 |

## CLI Flags

` + "```bash" + `
grit new myapp            # Full monorepo (api + web + admin + shared)
grit new myapp --api      # Go API only
grit new myapp --expo     # Full monorepo + Expo mobile app
grit new myapp --mobile   # API + Expo mobile app only
grit new myapp --full     # Everything + documentation site
` + "```" + `

## No Docker?

If you can't run Docker, use cloud services:

` + "```bash" + `
cp .env.cloud.example .env
` + "```" + `

Fill in your keys for [Neon](https://neon.tech) (Postgres), [Upstash](https://upstash.com) (Redis), [Cloudflare R2](https://dash.cloudflare.com) (storage), and [Resend](https://resend.com) (email).
`, opts.ProjectName, opts.ProjectName)
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
  "pages": ["index", "getting-started", "---API---", "api/authentication"]
}
`
}
