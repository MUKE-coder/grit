# Grit Framework — YouTube Content Calendar (30 Days)

**Channel:** https://www.youtube.com/@GritFramework
**Goal:** Awareness, SEO, practical guides
**Format:** Under 10 minutes, fast-paced, no fluff
**Schedule:** One video per day

---

## WEEK 1: INTRODUCTION & CORE CONCEPTS

### Day 1 — Channel Launch / Framework Introduction
**Title:** "Grit: The Full-Stack Go + React Framework (Everything You Need to Know)"
**Description:** Meet Grit — the framework that gives you auth, admin panel, file storage, email, background jobs, AI, and deployment in one CLI command. This is the only full-stack Go + React framework you'll ever need. Install it in 10 seconds, scaffold a production app in 30.
**Purpose:** Awareness — first impression, hook developers, establish what Grit is
**Timeline:**
- 0:00 — Hook: "What if one command gave you everything?"
- 0:30 — What is Grit? (Go + React meta-framework)
- 1:30 — Live demo: grit new myapp (interactive CLI)
- 3:00 — Tour of generated project (folders, files, what's inside)
- 5:00 — What ships by default (quick feature flythrough)
- 7:00 — Who is Grit for? (Laravel devs, Next.js devs, Go devs)
- 8:30 — Install command + docs link
- 9:00 — What's coming on this channel
**Thumbnail Prompt:** Create a YouTube thumbnail (1280x720) with dark navy background (#0b1120). Large bold white text "GRIT" on the left. Right side shows a terminal window with "grit new myapp" command. Sky blue (#38bdf8) accent glow behind the terminal. Bottom right: "Go + React Framework" in smaller text. Clean, developer-focused, no face.

---

### Day 2 — Installation & First Project
**Title:** "Install Grit and Build Your First App in 5 Minutes"
**Description:** Step-by-step: install Go, install Grit, scaffold your first project, start Docker, run the app. In 5 minutes you'll have a full-stack app with auth, database, and admin panel running on localhost.
**Purpose:** SEO — "install grit framework", "go react tutorial", practical onboarding
**Timeline:**
- 0:00 — Prerequisites check (Go, Node, pnpm, Docker)
- 1:00 — Install Grit: go install command
- 1:30 — grit new myapp (show interactive prompts)
- 2:30 — docker compose up -d (start services)
- 3:00 — pnpm install && pnpm dev
- 3:30 — Tour: web app (localhost:3000)
- 4:30 — Tour: admin panel (localhost:3001)
- 5:30 — Tour: API endpoints (localhost:8080)
- 6:30 — Tour: GORM Studio, API Docs, Pulse, Sentinel
- 8:00 — Recap + next video teaser
**Thumbnail Prompt:** YouTube thumbnail (1280x720), dark background. Left side: "5 MIN SETUP" in large bold yellow text. Right side: split screen showing a terminal on top and a web app dashboard on bottom. Small Grit logo in corner. Text at bottom: "Go + React Full-Stack App".

---

### Day 3 — Architecture Modes Explained
**Title:** "5 Ways to Build with Grit (Pick Your Architecture)"
**Description:** Grit lets you choose your architecture: Triple (Web + Admin + API), Double (Web + API), Single (one binary like Laravel), API-only, or Mobile (API + Expo). Learn which one is right for your project.
**Purpose:** Awareness — unique selling point, no other framework does this
**Timeline:**
- 0:00 — Hook: "Not every project needs the same structure"
- 0:30 — Triple: when and why (SaaS, platforms)
- 2:00 — Double: when and why (simpler apps, blogs)
- 3:00 — Single: when and why (Laravel devs, microservices)
- 4:30 — API-only: when and why (mobile backends, headless)
- 5:30 — Mobile: when and why (React Native apps)
- 6:30 — Desktop bonus: grit new-desktop
- 7:30 — Decision flowchart: "Which should YOU pick?"
- 9:00 — Quick demo: scaffold each mode
**Thumbnail Prompt:** YouTube thumbnail (1280x720), dark background. Center: 5 colored boxes in a row labeled "SINGLE", "DOUBLE", "TRIPLE", "API", "MOBILE". Title text above: "PICK YOUR ARCHITECTURE". Each box has a different icon (one app, two apps, three apps, server, phone). Sky blue and purple color scheme.

---

### Day 4 — Next.js vs TanStack Router
**Title:** "Next.js or TanStack Router? Choosing Your Grit Frontend"
**Description:** Grit supports two frontends: Next.js (SSR, SEO, App Router) and TanStack Router (Vite, SPA, fast builds). Here's how they compare and when to use each one.
**Purpose:** SEO — "next.js vs tanstack router", "vite react router", decision guide
**Timeline:**
- 0:00 — Hook: "Two frontends, same backend — which one?"
- 0:30 — Next.js overview: SSR, App Router, SEO
- 2:00 — TanStack Router overview: Vite, SPA, speed
- 3:30 — Side-by-side: build time, bundle size, DX
- 5:00 — When to use Next.js (marketing sites, SEO-heavy)
- 6:00 — When to use TanStack Router (dashboards, internal tools)
- 7:00 — Demo: same app, both frontends
- 8:30 — Can you mix them? (Yes — Next.js web + TanStack admin)
**Thumbnail Prompt:** YouTube thumbnail (1280x720), dark background. Left side: Next.js logo with "SSR" label. Right side: Vite logo with "SPA" label. Center: "VS" in large bold text with a lightning bolt. Bottom text: "Which Frontend for Your Grit App?". Split design with blue on left, purple on right.

---

### Day 5 — Code Generator Deep Dive
**Title:** "Generate a Full-Stack CRUD in 10 Seconds (Grit Code Generator)"
**Description:** Watch the code generator create a Go model, service, handler, Zod schema, TypeScript types, React hooks, and admin page — all from one command. Then we'll look at every generated file.
**Purpose:** Practical guide — core feature demo, show the magic
**Timeline:**
- 0:00 — Hook: "One command, 8 files, full-stack CRUD"
- 0:30 — The command: grit generate resource Product --fields "..."
- 1:30 — What it generated (file tree walkthrough)
- 2:30 — Go model: struct, tags, migration
- 3:30 — Go service: CRUD business logic
- 4:00 — Go handler: REST endpoints with pagination
- 4:30 — Zod schema + TypeScript types
- 5:00 — React Query hooks
- 5:30 — Admin page: DataTable + forms
- 6:30 — Route injection: how markers work
- 7:30 — Field types and modifiers overview
- 8:30 — grit remove: clean undo
- 9:00 — Interactive mode demo (-i flag)
**Thumbnail Prompt:** YouTube thumbnail (1280x720), dark background. Terminal showing "grit generate resource Product" on left. Right side: fan of generated files (model.go, handler.go, schema.ts, page.tsx) spreading out. Bold text: "ONE COMMAND = FULL CRUD". Green accent color for the terminal.

---

### Day 6 — Authentication System Tour
**Title:** "JWT + 2FA + OAuth — Grit's Auth System Explained"
**Description:** Every Grit project ships with JWT authentication, TOTP two-factor auth with backup codes, and OAuth2 social login. Walk through the entire auth flow from register to 2FA setup.
**Purpose:** SEO — "go jwt authentication", "totp golang", feature showcase
**Timeline:**
- 0:00 — Hook: "Most frameworks give you basic auth. Grit gives you bank-grade security."
- 0:30 — Registration flow (API + frontend)
- 1:30 — Login + JWT tokens (access + refresh)
- 2:30 — Role-based access control (ADMIN, EDITOR, USER)
- 3:30 — TOTP setup: QR code, authenticator app
- 5:00 — Backup codes: 10 one-time recovery codes
- 6:00 — Trusted devices: 30-day cookie
- 7:00 — OAuth2: Google + GitHub login
- 8:00 — Custom roles: grit add role MODERATOR
- 9:00 — Code walkthrough: where auth lives in the project
**Thumbnail Prompt:** YouTube thumbnail (1280x720), dark background. Large lock icon with a shield in sky blue. Text: "JWT + 2FA + OAuth". Subtitle: "Authentication That Ships by Default". Small QR code graphic in corner. Professional security-themed design.

---

### Day 7 — Admin Panel Showcase
**Title:** "The Admin Panel That Rivals Laravel Nova (Built-in, Free)"
**Description:** Grit's admin panel includes DataTable with pagination/sorting/filtering, FormBuilder with 8+ field types, dashboard widgets, 4 style variants, and system pages. All generated, all customizable.
**Purpose:** Awareness — show the admin panel quality, compare to paid tools
**Timeline:**
- 0:00 — Hook: "This admin panel ships free with every Grit project"
- 0:30 — Dashboard: stats cards, charts, activity feed
- 2:00 — DataTable: pagination, sorting, filtering, column visibility
- 3:30 — FormBuilder: text, select, date, richtext, file upload
- 4:30 — Multi-step forms: modal-steps + page-steps
- 5:30 — Resource system: defineResource() pattern
- 6:30 — Style variants: default, modern, minimal, glass
- 7:30 — System pages: Jobs, Files, Cron, Mail Preview
- 8:30 — Standalone usage: use DataTable/FormBuilder anywhere
**Thumbnail Prompt:** YouTube thumbnail (1280x720), dark background. Screenshot of a dark-themed admin dashboard with stats cards and a data table. Bold text overlay: "FREE ADMIN PANEL". Subtitle: "Ships with Every Grit Project". Gold/amber accent to suggest premium quality.

---

## WEEK 2: BATTERIES & FEATURES

### Day 8 — File Storage & Uploads
**Title:** "Presigned URL Uploads in Go (S3, R2, MinIO) — Grit Batteries"
**Description:** How Grit handles file uploads: presigned URLs that bypass the API, works with AWS S3, Cloudflare R2, or MinIO locally. Includes image processing via background jobs.
**Purpose:** SEO — "go s3 upload", "presigned url golang", practical guide
**Timeline:**
- 0:00 — Hook: "Your uploads shouldn't go through your API"
- 0:30 — The problem: traditional uploads vs presigned URLs
- 1:30 — How it works: browser → S3 directly, then record in DB
- 3:00 — Configuration: .env storage variables
- 4:00 — MinIO for local development
- 5:00 — Switching to Cloudflare R2 or AWS S3
- 6:00 — Image processing: background job thumbnails
- 7:00 — Upload progress tracking in the frontend
- 8:00 — Admin file management page
**Thumbnail Prompt:** YouTube thumbnail (1280x720), dark background. Cloud upload icon with arrows pointing to S3, R2, and MinIO logos. Text: "PRESIGNED UPLOADS". Subtitle: "Skip the API, Upload Direct". Blue and orange color scheme.

---

### Day 9 — Background Jobs & Cron
**Title:** "Background Jobs + Cron Scheduler in Go (Redis + asynq)"
**Description:** How Grit uses Redis and asynq for background job processing and scheduled tasks. Includes built-in workers for email, image processing, and cleanup — plus an admin dashboard.
**Purpose:** SEO — "golang background jobs", "asynq tutorial", "redis job queue go"
**Timeline:**
- 0:00 — Hook: "Don't make users wait for email sends"
- 0:30 — What are background jobs? When to use them
- 1:30 — asynq overview: Redis-backed Go job queue
- 2:30 — Built-in workers: email, image, cleanup
- 3:30 — Creating a custom job
- 5:00 — Cron scheduler: recurring tasks
- 6:00 — Admin dashboard: job status, retries, failures
- 7:00 — Error handling and retry strategies
- 8:00 — Production considerations
**Thumbnail Prompt:** YouTube thumbnail (1280x720), dark background. Conveyor belt illustration with job cards moving along it. Redis logo in corner. Text: "BACKGROUND JOBS IN GO". Subtitle: "Redis + asynq". Red and blue color scheme.

---

### Day 10 — Email Service
**Title:** "Transactional Emails in Go with Resend (Templates Included)"
**Description:** Grit includes a complete email service with Resend integration and 4 HTML templates. Welcome emails, password resets, verification, notifications — all working out of the box.
**Purpose:** SEO — "resend golang", "send email go", practical guide
**Timeline:**
- 0:00 — Hook: "Every app needs email. Here's how Grit handles it."
- 0:30 — Resend overview: modern email API
- 1:30 — Configuration: API key, sender address
- 2:30 — Built-in templates: welcome, reset, verify, notify
- 3:30 — Mailhog for local development
- 4:30 — Sending from a handler/service
- 5:30 — Sending via background jobs
- 6:30 — Custom templates: adding your own
- 7:30 — Admin mail preview page
**Thumbnail Prompt:** YouTube thumbnail (1280x720), dark background. Email envelope icon with code brackets inside. Text: "EMAIL IN GO". Subtitle: "Resend + HTML Templates". Purple and white color scheme.

---

### Day 11 — AI Integration
**Title:** "Add AI to Your Go App in 2 Minutes (Vercel AI Gateway)"
**Description:** Grit integrates with Vercel AI Gateway — one API key for Claude, GPT, Gemini, and hundreds more. Completions, multi-turn chat, and SSE streaming all built in.
**Purpose:** SEO — "vercel ai gateway", "ai golang", "openai go api", trending topic
**Timeline:**
- 0:00 — Hook: "One API key. Hundreds of AI models."
- 0:30 — The old way: separate keys for every provider
- 1:30 — Vercel AI Gateway: unified endpoint
- 2:30 — Configuration: 3 env vars
- 3:30 — Complete endpoint: POST /api/ai/complete
- 4:30 — Chat endpoint: multi-turn conversations
- 5:30 — Streaming: Server-Sent Events
- 6:30 — Switching models: just change one string
- 7:30 — Building a chat UI with React
**Thumbnail Prompt:** YouTube thumbnail (1280x720), dark background. Brain/AI icon connected to multiple provider logos (Anthropic, OpenAI, Google) via lines through a gateway. Text: "AI IN 2 MINUTES". Subtitle: "One Key, Hundreds of Models". Gradient blue to purple.

---

### Day 12 — Security with Sentinel
**Title:** "WAF + Rate Limiting for Your Go API (Sentinel Security)"
**Description:** Every Grit project ships with Sentinel — a Web Application Firewall with rate limiting, brute-force protection, anomaly detection, and a real-time threat dashboard.
**Purpose:** SEO — "golang rate limiting", "web application firewall go", security
**Timeline:**
- 0:00 — Hook: "Your API is under attack. You just don't know it yet."
- 0:30 — What is Sentinel? WAF + rate limiter + threat detection
- 1:30 — Rate limiting: per IP, per route configuration
- 3:00 — Brute-force protection on auth endpoints
- 4:00 — Security headers: HSTS, CSP, X-Frame-Options
- 5:00 — Anomaly detection: unusual patterns
- 6:00 — Threat dashboard tour (/sentinel/ui)
- 7:00 — Configuration: env vars and customization
- 8:00 — ExcludePaths: health checks, studio, docs
**Thumbnail Prompt:** YouTube thumbnail (1280x720), dark background. Shield icon with a lock, surrounded by blocked attack arrows. Text: "SECURE YOUR GO API". Subtitle: "WAF + Rate Limiting Built In". Red warning accents on dark blue.

---

### Day 13 — Observability with Pulse
**Title:** "Monitor Your Go API: Request Tracing, Metrics, Health Checks (Pulse)"
**Description:** Pulse gives you request tracing, database monitoring, runtime metrics, error tracking, health checks, and Prometheus export — all embedded in your Grit API.
**Purpose:** SEO — "golang observability", "go api monitoring", "prometheus go"
**Timeline:**
- 0:00 — Hook: "Can you see what your API is doing right now?"
- 0:30 — Pulse overview: self-hosted observability
- 1:30 — Request tracing: timing every endpoint
- 2:30 — Database monitoring: slow queries, connection pool
- 3:30 — Runtime metrics: goroutines, memory, GC
- 4:30 — Error tracking: catch and log errors
- 5:30 — Health checks: /pulse/health endpoint
- 6:30 — Prometheus export: integrate with Grafana
- 7:30 — Dashboard tour (/pulse/ui)
**Thumbnail Prompt:** YouTube thumbnail (1280x720), dark background. Dashboard mockup with charts, graphs, and metrics. Text: "MONITOR YOUR GO API". Subtitle: "Tracing + Metrics + Alerts". Green and blue color scheme suggesting health/monitoring.

---

### Day 14 — GORM Studio & API Docs
**Title:** "Visual Database Browser + Auto-Generated API Docs (GORM Studio + gin-docs)"
**Description:** GORM Studio lets you browse tables, edit data, run SQL, and export schemas. gin-docs auto-generates OpenAPI 3.1 documentation. Both embedded in every Grit project.
**Purpose:** SEO — "gorm studio", "gin api documentation", "openapi golang"
**Timeline:**
- 0:00 — Hook: "Two tools you'll use every day"
- 0:30 — GORM Studio: visual database browser
- 1:30 — Browse tables, view records, inline editing
- 2:30 — Raw SQL editor
- 3:00 — Schema export: SQL, JSON, YAML, ERD
- 4:00 — Data import/export: CSV, JSON, XLSX
- 4:30 — gin-docs: auto-generated API documentation
- 5:30 — OpenAPI 3.1 spec (no annotations needed)
- 6:30 — Interactive Scalar UI at /docs
- 7:30 — Postman/Insomnia export
- 8:00 — Both tools in one project, zero config
**Thumbnail Prompt:** YouTube thumbnail (1280x720), dark background. Split screen: database table view on left, API docs on right. Text: "DB BROWSER + API DOCS". Subtitle: "Zero Config, Ships by Default". Teal and white color scheme.

---

## WEEK 3: PRACTICAL BUILDS

### Day 15 — Build a Blog API in 5 Minutes
**Title:** "Build a Blog API with Go in 5 Minutes (Grit Generate)"
**Description:** Use grit generate to create a complete blog API with title, slug, content, published status, and categories. From zero to REST API with pagination, search, and filtering.
**Purpose:** SEO — "go blog api", "golang rest api tutorial", practical build
**Timeline:**
- 0:00 — Hook: "Blog API in 5 minutes. No boilerplate."
- 0:30 — The fields: title, slug, content, published, category
- 1:00 — Generate command with field types
- 2:00 — Test in API docs: create, list, get, update, delete
- 3:30 — Pagination and search working out of the box
- 4:30 — Admin panel: manage posts
- 5:30 — Add categories (belongs_to relationship)
- 7:00 — Public endpoints vs protected endpoints
- 8:00 — Deploy this to production (teaser)
**Thumbnail Prompt:** YouTube thumbnail (1280x720), dark background. Blog icon (document with lines) next to a terminal. Text: "BLOG API IN 5 MIN". Subtitle: "Go + Grit Code Generator". Timer graphic showing 5:00.

---

### Day 16 — Build a SaaS Dashboard
**Title:** "Build a SaaS Dashboard with Go + React (Full Tutorial)"
**Description:** Scaffold a triple architecture project, generate User, Subscription, and Invoice resources, customize the admin dashboard with stats cards and charts.
**Purpose:** SEO — "saas dashboard tutorial", "go react saas", practical build
**Timeline:**
- 0:00 — Hook: "SaaS dashboard from scratch in one video"
- 0:30 — Scaffold: grit new saas-app --triple --next
- 1:30 — Generate: User profile, Subscription, Invoice resources
- 3:00 — Customize dashboard: stats cards
- 4:30 — Add charts: revenue, signups
- 5:30 — Role-based access: admin vs user views
- 7:00 — Stripe integration (grit-stripe plugin teaser)
- 8:00 — What's next: billing, webhooks
**Thumbnail Prompt:** YouTube thumbnail (1280x720), dark background. SaaS dashboard mockup with stats cards, charts, and sidebar. Text: "SAAS DASHBOARD". Subtitle: "Go + React + Grit". Dollar sign icon in sky blue.

---

### Day 17 — Build an E-commerce API
**Title:** "E-commerce API in 10 Minutes: Products, Categories, Orders (Grit)"
**Description:** Generate Product, Category, and Order resources with relationships. Add image uploads for products. Build a complete e-commerce backend.
**Purpose:** SEO — "ecommerce api golang", "go rest api ecommerce", practical build
**Timeline:**
- 0:00 — Hook: "E-commerce backend. 10 minutes. Let's go."
- 0:30 — Plan: Product, Category, Order models
- 1:30 — Generate Category (name, slug, description)
- 2:30 — Generate Product (name, price, belongs_to Category, image)
- 4:00 — Generate Order (status, total, user relationship)
- 5:30 — Test the API: create products, place orders
- 7:00 — Admin panel: manage inventory
- 8:00 — Frontend: product listing page
- 9:00 — What's next: Stripe checkout, cart logic
**Thumbnail Prompt:** YouTube thumbnail (1280x720), dark background. Shopping cart icon with code brackets. Text: "E-COMMERCE API". Subtitle: "Products + Orders + Categories". Orange and white color scheme.

---

### Day 18 — Build a Task Manager (Single App)
**Title:** "Build a Task Manager with Go (Single Binary App)"
**Description:** Use the single architecture to build a task manager where the Go binary serves both the API and the React frontend. One binary, one deployment.
**Purpose:** SEO — "golang single binary app", "go embed react", practical build
**Timeline:**
- 0:00 — Hook: "One binary. API + frontend. Deploy anywhere."
- 0:30 — Scaffold: grit new tasks --single --vite
- 1:30 — Project structure tour (no monorepo, flat)
- 2:30 — Generate Task resource (title, description, done, priority)
- 4:00 — Frontend: task list with filters
- 5:30 — Mark complete, edit, delete
- 6:30 — Build: go build (frontend embedded via go:embed)
- 7:30 — Run the single binary
- 8:30 — Deploy: scp binary to server, done
**Thumbnail Prompt:** YouTube thumbnail (1280x720), dark background. Single cube/box icon representing one binary. Text: "ONE BINARY APP". Subtitle: "Go + React = Single File". Minimalist design, green accent.

---

### Day 19 — Build a Desktop App
**Title:** "Build a Desktop App with Go + React (Grit + Wails)"
**Description:** Use grit new-desktop to create a native desktop application with Go backend, React frontend, SQLite database, local auth, and PDF export.
**Purpose:** SEO — "wails tutorial", "go desktop app", "golang gui", practical build
**Timeline:**
- 0:00 — Hook: "Desktop apps with Go. Yes, really."
- 0:30 — grit new-desktop todo-app
- 1:30 — Project structure: Wails + React + SQLite
- 2:30 — Run: wails dev (hot reload)
- 3:30 — Tour: login, dashboard, CRUD
- 4:30 — Generate a new resource
- 5:30 — Custom title bar, frameless window
- 6:30 — Export: PDF and Excel
- 7:30 — Build: grit compile (native executable)
- 8:30 — Distribution: Windows, macOS, Linux
**Thumbnail Prompt:** YouTube thumbnail (1280x720), dark background. Desktop app window mockup with dark UI. Text: "DESKTOP APP IN GO". Subtitle: "Wails + React + SQLite". Desktop monitor icon, purple accents.

---

### Day 20 — Deploy to Production
**Title:** "Deploy Your Go App in One Command (SSH + systemd + Caddy)"
**Description:** Use grit deploy to cross-compile your app, upload via SSH, create a systemd service, and configure Caddy with automatic HTTPS. From localhost to production in one command.
**Purpose:** SEO — "deploy go app", "golang production", "caddy reverse proxy", practical guide
**Timeline:**
- 0:00 — Hook: "Code to HTTPS in one command"
- 0:30 — What grit deploy does (5-step pipeline)
- 1:30 — Prerequisites: VPS, SSH access, domain
- 2:30 — The command: grit deploy --host --domain
- 3:30 — Watch it work: build, upload, systemd, Caddy
- 5:00 — Verify: visit https://myapp.com
- 6:00 — What's running: systemd service explained
- 7:00 — Caddy config: auto-TLS, headers, logging
- 8:00 — Updates: re-run grit deploy
- 9:00 — Alternative: Docker deployment
**Thumbnail Prompt:** YouTube thumbnail (1280x720), dark background. Rocket launching from a terminal into a cloud. Text: "DEPLOY IN ONE COMMAND". Subtitle: "SSH + systemd + Auto-HTTPS". Gradient from blue to green.

---

### Day 21 — Grit Plugins Overview
**Title:** "10 Plugins That Supercharge Your Grit App"
**Description:** Overview of all 10 official Grit plugins: WebSockets, Stripe, OAuth, notifications, search, video, conferencing, webhooks, i18n, and export.
**Purpose:** Awareness — plugin ecosystem, show extensibility
**Timeline:**
- 0:00 — Hook: "Core features + plugins = unlimited possibilities"
- 0:30 — What are Grit plugins? (standalone Go packages)
- 1:00 — WebSockets: real-time chat, live dashboards
- 2:00 — Stripe: subscriptions, checkout, webhooks
- 3:00 — OAuth: Google, GitHub, Discord
- 3:30 — Notifications: in-app, push (FCM), SMS
- 4:30 — Search: Meilisearch full-text search
- 5:00 — Video: upload, transcode, HLS streaming
- 5:30 — Conference: WebRTC video calls
- 6:00 — Webhooks: outgoing events with HMAC
- 6:30 — i18n: multi-language support
- 7:00 — Export: PDF, Excel, CSV
- 7:30 — How to install and use a plugin
- 8:30 — Community plugins: how to create your own
**Thumbnail Prompt:** YouTube thumbnail (1280x720), dark background. Grid of 10 colorful plugin icons arranged in 2 rows of 5. Text: "10 PLUGINS". Subtitle: "WebSockets, Stripe, AI, Video & More". Rainbow accent colors on dark background.

---

## WEEK 4: ADVANCED TOPICS & TIPS

### Day 22 — Grit vs Laravel
**Title:** "Grit vs Laravel: Go + React Framework Compared to PHP King"
**Description:** Honest comparison: CLI scaffolding, code generation, admin panel, auth, deployment, performance, ecosystem. When to choose Grit, when to stick with Laravel.
**Purpose:** SEO — "go vs php", "laravel alternative", "golang framework comparison"
**Timeline:**
- 0:00 — Hook: "Can a Go framework match Laravel's DX?"
- 0:30 — CLI: grit new vs laravel new
- 1:30 — Code generation: grit generate vs artisan make
- 3:00 — Admin: Grit admin vs Laravel Nova/Filament
- 4:00 — Auth: JWT+TOTP vs Laravel Breeze/Fortify
- 5:00 — Deployment: grit deploy vs Laravel Forge
- 6:00 — Performance: Go vs PHP benchmarks
- 7:00 — Ecosystem: Composer vs Go modules
- 8:00 — Verdict: when to use which
**Thumbnail Prompt:** YouTube thumbnail (1280x720), dark background. Grit logo on left, Laravel logo on right, "VS" in center. Text: "GRIT vs LARAVEL". Subtitle: "Go + React vs PHP". Split design blue/red.

---

### Day 23 — Grit for Next.js Developers
**Title:** "Why Next.js Developers Should Try Grit (Go Backend > Serverless)"
**Description:** If you love Next.js but hate serverless limitations, Grit gives you a proper Go backend with the same React frontend. Real database, real caching, real background jobs.
**Purpose:** SEO — "next.js backend", "next.js alternative backend", audience expansion
**Timeline:**
- 0:00 — Hook: "Love Next.js? Hate serverless cold starts?"
- 0:30 — The problem: API routes limitations
- 1:30 — What Grit adds: Go backend with everything
- 3:00 — Same React, same TanStack Query, same Zod
- 4:00 — But now: real PostgreSQL, Redis, background jobs
- 5:30 — File uploads that actually work (presigned URLs)
- 6:30 — Admin panel for free
- 7:30 — Deploy to any VPS (not just Vercel)
- 8:30 — Migration guide: adding Grit to existing Next.js
**Thumbnail Prompt:** YouTube thumbnail (1280x720), dark background. Next.js logo with an arrow pointing to Grit logo. Text: "NEXT.JS + GO BACKEND". Subtitle: "Better Than Serverless". Blue gradient.

---

### Day 24 — Environment Variables & Configuration
**Title:** "Every Grit Environment Variable Explained (Complete .env Guide)"
**Description:** Walk through every environment variable in a Grit project: database, auth, storage, email, AI, security, observability. What each one does and how to configure for production.
**Purpose:** SEO — "golang environment variables", "go env configuration", practical reference
**Timeline:**
- 0:00 — Hook: "Your .env file has 30+ variables. Here's what they all do."
- 0:30 — Core: APP_NAME, APP_ENV, APP_PORT, APP_URL
- 1:30 — Database: DATABASE_URL patterns
- 2:30 — Auth: JWT_SECRET, TOTP_ISSUER
- 3:00 — Storage: STORAGE_DRIVER and S3 config
- 4:00 — Email: RESEND_API_KEY, MAIL_FROM
- 4:30 — AI: AI_GATEWAY_API_KEY, AI_GATEWAY_MODEL
- 5:30 — OAuth: Google + GitHub client IDs
- 6:00 — Security: Sentinel config
- 6:30 — Observability: Pulse config
- 7:00 — GORM Studio: username, password
- 7:30 — Production vs development differences
- 8:30 — .env.example and .env.cloud.example
**Thumbnail Prompt:** YouTube thumbnail (1280x720), dark background. .env file icon with key-value pairs visible. Text: "EVERY ENV VAR EXPLAINED". Subtitle: "Complete Grit Configuration Guide". Green terminal-style text on dark.

---

### Day 25 — Relationships & Advanced Fields
**Title:** "Database Relationships in Grit: belongs_to, many_to_many, slug"
**Description:** How to use belongs_to, many_to_many, slug, and string_array field types. Build a blog with categories (belongs_to) and tags (many_to_many).
**Purpose:** SEO — "gorm relationships", "golang belongs to", practical guide
**Timeline:**
- 0:00 — Hook: "Real apps have relationships. Here's how Grit handles them."
- 0:30 — belongs_to: foreign key relationships
- 2:00 — Demo: Post belongs_to Category
- 3:30 — many_to_many: junction tables
- 5:00 — Demo: Post many_to_many Tag
- 6:30 — slug: auto-generated URL-friendly strings
- 7:30 — string_array: JSON arrays in PostgreSQL
- 8:30 — Field modifiers: :unique, :optional
- 9:00 — Admin panel: relationship select fields
**Thumbnail Prompt:** YouTube thumbnail (1280x720), dark background. Database diagram showing connected tables with arrows. Text: "RELATIONSHIPS IN GO". Subtitle: "belongs_to + many_to_many". Blue connected nodes design.

---

### Day 26 — Custom Roles & Permissions
**Title:** "Role-Based Access Control in Go (RBAC with Grit)"
**Description:** How Grit handles roles: built-in ADMIN/EDITOR/USER, adding custom roles with grit add role, protecting routes per role, and admin-only endpoints.
**Purpose:** SEO — "golang rbac", "go role based access", practical guide
**Timeline:**
- 0:00 — Hook: "Not every user should see everything"
- 0:30 — Built-in roles: ADMIN, EDITOR, USER
- 1:30 — Role middleware: RequireRole("ADMIN")
- 3:00 — Adding custom roles: grit add role MODERATOR
- 4:00 — What grit add role changes (Go + TypeScript + Zod)
- 5:00 — Role-restricted resource generation: --roles flag
- 6:00 — Admin panel: user role management
- 7:00 — Frontend: conditional rendering by role
- 8:00 — Best practices: least privilege, role hierarchy
**Thumbnail Prompt:** YouTube thumbnail (1280x720), dark background. Three user silhouettes at different heights (hierarchy). Text: "RBAC IN GO". Subtitle: "Roles + Permissions + Middleware". Purple and blue gradient.

---

### Day 27 — Maintenance Mode & Route Listing
**Title:** "grit down, grit up, grit routes — 3 Commands You Need to Know"
**Description:** Maintenance mode (grit down/up) and route listing (grit routes). Simple but powerful CLI commands for operations and debugging.
**Purpose:** SEO — "golang maintenance mode", practical tips, short and useful
**Timeline:**
- 0:00 — Hook: "Three commands that make your life easier"
- 0:30 — grit routes: see every API endpoint
- 2:00 — Reading the route table: method, path, handler, group
- 3:00 — grit down: maintenance mode
- 4:00 — What happens: .maintenance file, 503 responses
- 5:00 — Use case: during deployments, migrations
- 6:00 — grit up: back online
- 6:30 — Scripting: use in deploy pipelines
- 7:00 — Bonus: grit version, grit update
**Thumbnail Prompt:** YouTube thumbnail (1280x720), dark background. Three terminal commands stacked vertically: "grit routes", "grit down", "grit up". Text: "3 MUST-KNOW COMMANDS". Green, yellow, green color coding (up, warning, up).

---

### Day 28 — Redis Caching Strategies
**Title:** "Redis Caching in Go: Cache Middleware, TTL, Invalidation (Grit)"
**Description:** How to use Grit's built-in Redis cache service: GET/SET operations, cache middleware for API responses, TTL configuration, and cache invalidation patterns.
**Purpose:** SEO — "golang redis cache", "go api caching", practical guide
**Timeline:**
- 0:00 — Hook: "Your API is fast. Caching makes it instant."
- 0:30 — Redis cache service overview
- 1:30 — Get/Set/Delete operations
- 2:30 — Cache middleware: automatic response caching
- 4:00 — TTL configuration per key
- 5:00 — Cache invalidation: when to clear
- 6:00 — Patterns: cache-aside, write-through
- 7:00 — Admin: cache management
- 8:00 — Production: Redis configuration tips
**Thumbnail Prompt:** YouTube thumbnail (1280x720), dark background. Redis logo with speed/lightning effects. Text: "REDIS CACHING IN GO". Subtitle: "Middleware + TTL + Invalidation". Red Redis color with blue accents.

---

### Day 29 — Grit UI Components
**Title:** "100 Free React Components That Ship With Your Grit App"
**Description:** Tour the Grit UI component registry: 100 shadcn-compatible components across marketing, auth, SaaS, ecommerce, and layout categories. Install with one command.
**Purpose:** SEO — "free react components", "shadcn components", showcase
**Timeline:**
- 0:00 — Hook: "100 components. Free. Already in your project."
- 0:30 — What is Grit UI? shadcn-compatible registry
- 1:30 — Marketing components: heroes, features, pricing
- 3:00 — Auth components: login, register, OTP forms
- 4:00 — SaaS components: dashboards, billing, settings
- 5:00 — Ecommerce components: products, cart, checkout
- 6:00 — Layout components: navbars, sidebars, footers
- 7:00 — Installing: npx shadcn@latest add --url
- 8:00 — Customizing: it's just Tailwind + React
**Thumbnail Prompt:** YouTube thumbnail (1280x720), dark background. Grid of small UI component previews (cards, forms, navbars). Text: "100 FREE COMPONENTS". Subtitle: "shadcn-Compatible React UI Kit". Colorful component grid on dark.

---

### Day 30 — What's Next for Grit
**Title:** "Grit Roadmap: What's Coming in v4.0 and Beyond"
**Description:** The future of Grit: community feedback, planned features, plugin marketplace, Grit Cloud, and how to contribute. Plus a recap of everything we covered this month.
**Purpose:** Awareness — community building, set expectations, call to action
**Timeline:**
- 0:00 — Hook: "30 days of Grit. Here's what's next."
- 0:30 — Month recap: what we built and learned
- 2:00 — Community: GitHub stars, issues, contributions
- 3:00 — Planned: more plugins (payments, analytics)
- 4:00 — Planned: plugin marketplace
- 5:00 — Planned: Grit Cloud (hosted platform)
- 6:00 — Planned: more architecture modes, frontend options
- 7:00 — How to contribute: issues, PRs, plugins
- 8:00 — How to support: star, share, sponsor
- 9:00 — Thank you + subscribe CTA
**Thumbnail Prompt:** YouTube thumbnail (1280x720), dark background. Road/path stretching into the horizon with Grit logo at the end. Text: "GRIT ROADMAP". Subtitle: "What's Coming in v4.0". Blue gradient horizon, futuristic feel.

---

## CONTENT STRATEGY NOTES

### SEO Keywords Targeted:
- golang full stack framework
- go react framework
- go rest api tutorial
- gorm tutorial
- gin framework tutorial
- next.js go backend
- laravel alternative golang
- deploy go app
- golang authentication jwt
- vercel ai gateway
- wails desktop app
- golang background jobs
- redis golang cache
- shadcn components free

### Posting Schedule:
- Upload daily at 9:00 AM UTC
- Premiere format for first 3 videos (build hype)
- Community posts on off-days with code snippets
- Shorts from each video (30-60 second clips of the best moments)

### Thumbnail Consistency:
- Always dark background (#0b1120)
- Bold white text, sky blue (#38bdf8) accents
- One main visual element (terminal, dashboard, icon)
- Subtitle in smaller text below main title
- No faces unless doing a talking-head video
- Consistent placement: title left/center, visual right
