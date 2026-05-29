import type { Metadata } from 'next'

const base = 'https://gritframework.dev'

interface DocPage {
  title: string
  description: string
}

// Every doc page with its SEO title and description
export const docsMetadata: Record<string, DocPage> = {
  // Introduction
  '/docs': {
    title: 'Introduction',
    description:
      'Get started with Grit, the full-stack meta-framework that combines Go (Gin + GORM) with React (Next.js) and a Filament-like admin panel.',
  },

  // Getting Started
  '/docs/stack-selector': {
    title: 'Stack Selector — Pick the Right Grit Combo',
    description:
      'Eleven Grit architecture combos with use cases, exact commands, trade-offs, and a capability matrix. Helps humans (and AI assistants) pick the right scaffold for what they\'re building — web portal, offline desktop, multi-platform, API only, and everything in between.',
  },
  '/docs/getting-started/installation': {
    title: 'Installation',
    description:
      'Install Grit CLI and set up your development environment. Requires Go 1.21+, Node.js 18+, pnpm, and Docker.',
  },
  '/docs/getting-started/quick-start': {
    title: 'Quick Start',
    description:
      'Create your first Grit project in under 5 minutes. Scaffold a full-stack app with Go API, Next.js frontend, and admin panel.',
  },
  '/docs/getting-started/configuration': {
    title: 'Configuration',
    description:
      'Configure your Grit project with environment variables, database connections, JWT secrets, Redis, S3 storage, and more.',
  },
  '/docs/getting-started/philosophy': {
    title: 'Philosophy',
    description:
      'The design principles behind Grit: convention over configuration, batteries included, type safety, and developer experience.',
  },
  '/docs/getting-started/project-structure': {
    title: 'Project Structure',
    description:
      'Understand the Grit monorepo structure: apps/api (Go), apps/web (Next.js), apps/admin (Next.js), and packages/shared.',
  },
  '/docs/getting-started/troubleshooting': {
    title: 'Troubleshooting',
    description:
      'Common issues and solutions when working with Grit projects, including Docker, database, and build errors.',
  },
  '/docs/getting-started/create-without-docker': {
    title: 'Create without Docker',
    description:
      'Set up a Grit project using cloud services (Neon, Upstash, Cloudflare R2, Resend) instead of Docker. No Docker required.',
  },
  '/docs/getting-started/cli-cheatsheet': {
    title: 'CLI Cheatsheet',
    description:
      'Complete Grit CLI reference: every command, flag, and field type. Quick-copy recipes for scaffolding, code generation, migrations, seeding, and more.',
  },

  // Prerequisites
  '/docs/prerequisites/golang': {
    title: 'Go for Grit Developers',
    description:
      'Learn Go fundamentals for building Grit backends: variables, structs, functions, error handling, interfaces, pointers, goroutines, Gin routing, GORM, middleware, and JWT authentication.',
  },
  '/docs/prerequisites/nextjs': {
    title: 'Next.js for Grit Developers',
    description:
      'Learn Next.js fundamentals for building Grit frontends: App Router, server/client components, data fetching, routing, and React Query.',
  },
  '/docs/prerequisites/docker': {
    title: 'Docker for Grit Developers',
    description:
      'Learn Docker fundamentals for running Grit infrastructure: containers, images, docker-compose, PostgreSQL, Redis, and MinIO.',
  },

  // Concepts
  '/docs/concepts/architecture': {
    title: 'Architecture',
    description:
      'Understand Grit architecture: monorepo layout, Go API with handler-service-model pattern, Next.js frontend, shared types, and code generation.',
  },
  '/docs/concepts/architecture-modes': {
    title: 'Architecture Modes',
    description:
      'Grit supports 5 architecture modes: triple (web + admin + API), double (web + API), single (embedded SPA), API only, and mobile. Choose the one that fits your project.',
  },
  '/docs/concepts/architecture-modes/triple': {
    title: 'Triple Architecture: Web + Admin + API',
    description:
      'The default and most feature-rich Grit architecture. Turborepo monorepo with 3 apps (web, admin, API) sharing types via packages/shared. Best for SaaS, marketplaces, and content platforms.',
  },
  '/docs/concepts/architecture-modes/double': {
    title: 'Double Architecture: Web + API',
    description:
      'Streamlined Turborepo monorepo with 2 apps (web + API). Admin features live as role-protected routes in the web app. Best for blogs, portfolios, and simpler SaaS applications.',
  },
  '/docs/concepts/architecture-modes/single': {
    title: 'Single Architecture: One Binary',
    description:
      'Single architecture in Grit: one Go binary serves both API and embedded React frontend via go:embed. Flat project structure, no monorepo, simplest deployment.',
  },
  '/docs/concepts/architecture-modes/api-only': {
    title: 'API Only Architecture: Headless Go Backend',
    description:
      'API-only architecture in Grit: pure Go backend with no frontend. All batteries included (auth, storage, email, jobs, AI). Perfect for mobile backends and microservices.',
  },
  '/docs/concepts/architecture-modes/mobile': {
    title: 'Mobile Architecture: API + Expo React Native',
    description:
      'Mobile architecture in Grit: Go API paired with Expo React Native in a Turborepo monorepo. SecureStore for tokens, Expo Router for navigation, shared types between backend and mobile.',
  },
  '/docs/concepts/architecture-modes/multi-client': {
    title: 'Multi-Client: API + Mobile + Desktop',
    description:
      'Combine --desktop with other flags to scaffold a shared-API multi-client SaaS. One Go backend serving Expo mobile and Wails desktop, with optional web and admin. The Linear / Notion / Slack pattern.',
  },
  '/docs/concepts/cli': {
    title: 'CLI Commands',
    description:
      'Complete reference for Grit CLI commands: grit new, grit generate resource, grit sync, grit add role, grit start, and grit remove.',
  },
  '/docs/concepts/code-generation': {
    title: 'Code Generation',
    description:
      'How Grit code generation works: generating full-stack resources with models, handlers, services, Zod schemas, TypeScript types, hooks, and admin pages.',
  },
  '/docs/concepts/naming-conventions': {
    title: 'Naming Conventions',
    description:
      'Naming conventions in Grit: Go files (snake_case), TypeScript files (kebab-case), React components (PascalCase), API routes (plural lowercase).',
  },
  '/docs/concepts/styles': {
    title: 'Style Variants',
    description:
      'Choose from 4 admin panel style variants in Grit: default, modern, minimal, and glass themes.',
  },
  '/docs/concepts/type-system': {
    title: 'Type System',
    description:
      'How Grit shares types between Go and TypeScript: Go structs to Zod schemas to TypeScript interfaces, keeping frontend and backend in sync.',
  },
  '/docs/concepts/performance': {
    title: 'Performance',
    description:
      'All performance optimisations built into every Grit project: Gzip, Request ID, connection pool tuning, Cache-Control, presigned uploads, background jobs, Redis caching, Server Components, ISR, React Query, next/image, and Turborepo.',
  },

  // Backend
  '/docs/backend/authentication': {
    title: 'Authentication',
    description:
      'Implement JWT authentication in Grit: login, register, token refresh, password hashing with bcrypt, and protected routes.',
  },
  '/docs/backend/oauth': {
    title: 'Social Login (OAuth2)',
    description:
      'Set up Google and GitHub OAuth2 social login in Grit: provider configuration, callback URLs, account linking, and production deployment.',
  },
  '/docs/backend/handlers': {
    title: 'Handlers',
    description:
      'Write Gin HTTP handlers in Grit: request parsing, validation with binding tags, JSON responses, pagination, and error handling.',
  },
  '/docs/backend/middleware': {
    title: 'Middleware',
    description:
      'Built-in Grit middleware: authentication, CORS, logging, rate limiting, cache, and how to write custom Gin middleware.',
  },
  '/docs/backend/migrations': {
    title: 'Migrations',
    description:
      'Database migrations in Grit with GORM AutoMigrate: adding fields, creating tables, and managing schema changes.',
  },
  '/docs/backend/models': {
    title: 'Models',
    description:
      'Define GORM models in Grit: struct tags, field types, relationships (belongs_to, many_to_many), hooks, and soft deletes.',
  },
  '/docs/backend/rbac': {
    title: 'RBAC',
    description:
      'Role-based access control in Grit: ADMIN, EDITOR, USER roles, RequireRole middleware, role-restricted routes, and grit add role.',
  },
  '/docs/backend/api-docs': {
    title: 'API Documentation',
    description:
      'Auto-generated API documentation in Grit with gin-docs: zero-annotation OpenAPI spec, interactive Scalar/Swagger UI, Postman/Insomnia export, and GORM model schemas.',
  },
  '/docs/backend/pulse': {
    title: 'Pulse (Observability)',
    description:
      'Self-hosted observability for Grit APIs with Pulse: request tracing, database monitoring, runtime metrics, error tracking, health checks, alerting, and Prometheus export.',
  },
  '/docs/security': {
    title: 'Security Guide — OWASP Top 10:2025 defences',
    description:
      'How Grit defends every category of the OWASP Top 10:2025 by default — broken access control / IDOR via authz.MustOwn, SSRF via the safefetch package, injection via parameterised GORM queries, XSS via React escaping + CSP, auth flaws via Sentinel rate-limiting + JWT algorithm pinning + TOTP, security misconfiguration via the SecurityHeaders middleware, supply-chain via Dependabot + govulncheck, and tamper-evident audit logging.',
  },
  '/docs/testing': {
    title: 'Performance & Security Testing',
    description:
      'Run k6 load tests (smoke / average-load / stress / spike / soak / breakpoint) and methodology-driven penetration tests against a Grit app. Includes the 5-phase pentest methodology, the attack catalogue, CVSS scoring, and the audit-report structure that justifies a $2k–$10k engagement.',
  },
  '/docs/demo': {
    title: 'Demo Application — Grit Motors',
    description:
      'A live, full-stack Grit demo built on grit new --single --vite: motorcycle dealership management, POS, loans + repayments, daily-boda fleet tracking, multi-tenant RBAC, in-app Security (Sentinel) + Observability (Pulse) dashboards. Source on GitHub, login pre-filled, database reset nightly.',
  },
  '/docs/backend/response-format': {
    title: 'API Response Format',
    description:
      'Standard API response format in Grit: success responses with data/message, paginated lists with meta, and error responses with codes.',
  },
  '/docs/backend/seeders': {
    title: 'Seeders',
    description:
      'Seed your Grit database with initial data: admin users, sample records, and the built-in blog example with posts.',
  },
  '/docs/backend/services': {
    title: 'Services',
    description:
      'The service pattern in Grit: business logic separation, GORM queries, pagination, filtering, and the Services struct.',
  },

  // Frontend
  '/docs/frontend/hooks': {
    title: 'React Hooks',
    description:
      'Generated React Query hooks in Grit: useList, useGet, useCreate, useUpdate, useDelete for every resource with type safety.',
  },
  '/docs/frontend/shared-package': {
    title: 'Shared Package',
    description:
      'The packages/shared module in Grit: Zod schemas, TypeScript types, API route constants shared between web and admin apps.',
  },
  '/docs/frontend/web-app': {
    title: 'Web App',
    description:
      'The Next.js web app in Grit: App Router pages, authentication flow, dashboard layout, API client, and React Query setup.',
  },

  // Admin Panel
  '/docs/admin/overview': {
    title: 'Admin Panel Overview',
    description:
      'Grit admin panel: a Filament-like dashboard with runtime resource definitions, DataTable, FormBuilder, widgets, and dark/light theme.',
  },
  '/docs/admin/resources': {
    title: 'Resources',
    description:
      'Define admin resources in Grit with defineResource(): columns, filters, sorting, search, forms, and permissions.',
  },
  '/docs/admin/datatable': {
    title: 'DataTable',
    description:
      'Advanced DataTable in Grit admin: sorting, filtering, search, pagination, column visibility, row selection, and custom cell renderers.',
  },
  '/docs/admin/forms': {
    title: 'Forms',
    description:
      'FormBuilder in Grit admin: text, number, select, date, toggle, checkbox, radio, textarea, richtext, and relationship fields.',
  },
  '/docs/admin/multi-step-forms': {
    title: 'Multi-Step Forms',
    description:
      'Multi-step forms in Grit: modal-steps and page-steps variants with step indicators, per-step validation, and progress tracking.',
  },
  '/docs/admin/relationships': {
    title: 'Relationships',
    description:
      'Relationship fields in Grit: belongs_to and many_to_many with searchable select dropdowns, multi-select tags, and automatic foreign keys.',
  },
  '/docs/admin/widgets': {
    title: 'Dashboard Widgets',
    description:
      'Dashboard widgets in Grit admin: StatsCard, ChartWidget (Recharts), ActivityWidget, and WidgetGrid for building custom dashboards.',
  },
  '/docs/admin/standalone-usage': {
    title: 'Standalone Usage',
    description:
      'Use Grit components (DataTable, FormBuilder, FormStepper) on any page in web or admin apps without the resource system.',
  },

  // Batteries
  '/docs/batteries/ai': {
    title: 'AI Integration',
    description:
      'AI integration in Grit: Claude and OpenAI support with streaming responses, configurable providers, and an AI handler.',
  },
  '/docs/batteries/caching': {
    title: 'Caching',
    description:
      'Redis caching in Grit: cache service, cache middleware for API responses, TTL configuration, and cache invalidation.',
  },
  '/docs/batteries/cron': {
    title: 'Cron Jobs',
    description:
      'Cron scheduling in Grit with asynq: define recurring tasks, cron expressions, admin dashboard for monitoring schedules.',
  },
  '/docs/batteries/email': {
    title: 'Email',
    description:
      'Send emails in Grit with Resend: welcome, password reset, verification, and notification templates with HTML layouts.',
  },
  '/docs/batteries/jobs': {
    title: 'Background Jobs',
    description:
      'Background job processing in Grit with asynq and Redis: email jobs, image processing, cleanup workers, and admin monitoring.',
  },
  '/docs/batteries/security': {
    title: 'Security',
    description:
      'Security in Grit with Sentinel: WAF, rate limiting, brute-force protection, anomaly detection, IP geolocation, and threat dashboard.',
  },
  '/docs/batteries/storage': {
    title: 'File Storage',
    description:
      'S3-compatible file storage in Grit: upload handler, image processing, MinIO for development, Cloudflare R2 or AWS S3 for production.',
  },

  // Plugins
  '/docs/plugins': {
    title: 'Plugins',
    description:
      'Grit plugins: drop-in Go packages for WebSockets, Stripe payments, OAuth, notifications, search, video processing, conferencing, webhooks, i18n, and data export.',
  },

  // Infrastructure
  '/docs/infrastructure/database': {
    title: 'Database',
    description:
      'Database setup in Grit: PostgreSQL for production, SQLite for development, GORM Studio visual browser, and connection configuration.',
  },
  '/docs/infrastructure/deployment': {
    title: 'Deployment',
    description:
      'Deploy Grit projects: Docker production builds, environment configuration, database setup, and hosting options.',
  },
  '/docs/infrastructure/dokploy': {
    title: 'Deploy with Dokploy',
    description:
      'Deploy your Grit application with Dokploy: self-hosted PaaS with Docker Compose, auto-SSL, GitHub integration, and a web dashboard on your own VPS.',
  },
  '/docs/infrastructure/docker': {
    title: 'Docker',
    description:
      'Docker setup in Grit: docker-compose for PostgreSQL, Redis, MinIO, and Mailhog. Production Dockerfiles for Go API and Next.js apps.',
  },
  '/docs/infrastructure/docker-cheatsheet': {
    title: 'Docker Cheatsheet',
    description:
      'Quick reference for Docker commands used with Grit: container management, volumes, networking, and troubleshooting.',
  },

  // Design
  '/docs/design/theme': {
    title: 'Theme',
    description:
      'Grit design system: dark mode default, color palette, typography (Onest + JetBrains Mono), and component styling with Tailwind CSS.',
  },

  // Tutorials
  '/docs/tutorials/contact-app': {
    title: 'Your First App',
    description:
      'Step-by-step tutorial: build a contact manager with Grit. Create Group and Contact resources, explore the admin panel, GORM Studio, and API docs.',
  },
  '/docs/tutorials/blog': {
    title: 'Build a Blog Tutorial',
    description:
      'Step-by-step tutorial: build a full-stack blog with Grit including Go API, Next.js pages, admin panel, and SEO-friendly URLs.',
  },
  '/docs/tutorials/ecommerce': {
    title: 'Build an E-Commerce App',
    description:
      'Tutorial: build a full-stack e-commerce application with Grit including products, categories, orders, and admin management.',
  },
  '/docs/tutorials/learn': {
    title: 'Learn Grit Step by Step',
    description:
      'Beginner tutorial: learn Grit from scratch by building a complete full-stack application with Go API and React frontend.',
  },
  '/docs/tutorials/product-catalog': {
    title: 'Build a Product Catalog',
    description:
      'Tutorial: build a product catalog with Grit using code generation, multi-step forms, standalone DataTable, and FormBuilder.',
  },
  '/docs/tutorials/saas': {
    title: 'Build a SaaS App',
    description:
      'Tutorial: build a multi-tenant SaaS application with Grit including authentication, billing, teams, and admin dashboard.',
  },

  // AI Workflows
  '/docs/ai-workflows/claude': {
    title: 'Using Grit with Claude',
    description:
      'How to use Claude AI to build Grit projects faster: prompting strategies, code generation, and AI-assisted development.',
  },
  '/docs/ai-workflows/antigravity': {
    title: 'Using Grit with Antigravity',
    description:
      'How to use Antigravity AI assistant with Grit projects for faster development and code generation.',
  },
  '/docs/ai-skill': {
    title: 'AI Skill',
    description:
      'The Grit AI skill: teach AI assistants about Grit conventions, architecture, and code patterns for better code generation.',
  },
  '/docs/ai-skill/llm-guide': {
    title: 'Complete LLM Reference',
    description:
      'The complete Grit reference for AI assistants and LLMs: framework overview, full project structure, every CLI command, all field types, code patterns, API conventions, code markers, naming rules, all batteries, performance features, and the golden rules never to break.',
  },

  // Desktop (Wails)
  '/docs/desktop': {
    title: 'Desktop App Development',
    description:
      'Build native desktop applications with Grit and Wails. Scaffold complete Wails projects with Go backend, React frontend, SQLite, authentication, and CRUD.',
  },
  '/docs/desktop/getting-started': {
    title: 'Getting Started with Desktop',
    description:
      'Prerequisites, scaffolding, and development workflow for Grit desktop apps using Wails, Go, and React.',
  },
  '/docs/desktop/first-app': {
    title: 'Your First Desktop App',
    description:
      'Step-by-step tutorial: build a Task Manager desktop app with Grit and Wails. Scaffold, generate resources, browse with GORM Studio, and compile for distribution.',
  },
  '/docs/desktop/resource-generation': {
    title: 'Desktop Resource Generation',
    description:
      'Generate full-stack CRUD resources for desktop apps. Models, services, TanStack Router route files, and 10 automatic code injections.',
  },
  '/docs/desktop/offline': {
    title: 'Building an Offline-First Desktop App',
    description:
      'Use Grit\'s built-in sync engine to ship desktop apps that work fully offline. Local SQLite mirror, outbox with squash semantics, manual Sync button, field-level conflict resolution, versioned writes — Git workflow applied to your data.',
  },
  '/docs/desktop/building': {
    title: 'Building & Distribution',
    description:
      'Compile Grit desktop apps into native executables. Cross-platform builds, NSIS installers, and distribution tips.',
  },
  '/docs/desktop/pos-app': {
    title: 'Build a POS App',
    description:
      'Advanced tutorial: build a Point of Sale desktop application with Grit Desktop. Product catalog, sales transactions, receipts, inventory tracking, and daily reports in a single native binary.',
  },
  '/docs/desktop/project-ideas': {
    title: '20 Desktop Project Ideas',
    description:
      'Ready-to-build desktop application ideas with Grit and Wails. Each project includes resources, field definitions, and grit generate commands to get started immediately.',
  },
  '/docs/desktop/llm-reference': {
    title: 'Desktop LLM Reference',
    description:
      'The complete Grit Desktop reference for AI assistants and LLMs: architecture, Wails bindings, CLI commands, resource generation, field types, code markers, DataTable, FormBuilder, building executables, and the golden rules.',
  },

  // Changelog
  '/docs/changelog': {
    title: 'Changelog',
    description:
      'All notable changes to Grit: new features, bug fixes, and breaking changes for each release.',
  },

  // Course
  '/course': {
    title: 'Grit Fullstack Course',
    description:
      'Master fullstack development with Grit. Build 5 real-world apps, complete 50+ assignments, and learn Go, React, Next.js from zero to production.',
  },

  // Playground
  '/playground': {
    title: 'Go Playground',
    description:
      'Interactive Go code editor powered by the official Go Playground API. Write, run, and share Go code directly in your browser.',
  },
}

// Helper to generate Metadata for a doc page
export function getDocMetadata(path: string): Metadata {
  const page = docsMetadata[path]
  if (!page) return {}

  return {
    title: page.title,
    description: page.description,
    alternates: {
      canonical: `${base}${path}`,
    },
    openGraph: {
      title: `${page.title} | Grit`,
      description: page.description,
      url: `${base}${path}`,
      type: 'article',
    },
  }
}
