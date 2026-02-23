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

  // Backend
  '/docs/backend/authentication': {
    title: 'Authentication',
    description:
      'Implement JWT authentication in Grit: login, register, token refresh, password hashing with bcrypt, and protected routes.',
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
      'Grit design system: dark mode default, color palette, typography (DM Sans + JetBrains Mono), and component styling with Tailwind CSS.',
  },

  // Tutorials
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

  // Changelog
  '/docs/changelog': {
    title: 'Changelog',
    description:
      'All notable changes to Grit: new features, bug fixes, and breaking changes for each release.',
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
