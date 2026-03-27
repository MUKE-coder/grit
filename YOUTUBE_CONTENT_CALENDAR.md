# Grit Framework — YouTube Content Calendar (30 Days)

**Channel:** https://www.youtube.com/@GritFramework
**Goal:** Awareness, SEO, practical course walkthroughs
**Format:** Under 10 minutes each, "coding along with me" style
**Schedule:** One video per day, 30 consecutive days
**Structure:** 19 courses (9 web + 5 desktop + 5 mobile), each split into 2-3 videos

---

## Content Strategy

### SEO Keywords (target in titles, descriptions, tags)
- Primary: grit framework, go react framework, full-stack go, go web framework
- Secondary: go crud generator, go jwt auth, go admin panel, go file upload s3, go background jobs
- Long-tail: "build saas with go", "go react monorepo", "wails desktop app", "expo react native go backend"
- Comparison: "go vs laravel", "gin vs next.js", "go react vs next.js full stack"

### Posting Schedule
- Upload daily at 9:00 AM UTC (consistent time builds algorithm trust)
- Publish Monday through Sunday, no gaps
- Pin a playlist "Grit Framework — Full Course" on the channel page
- Add end screens linking to the next video in the series

### Shorts Strategy
- Extract one 30-60 second clip from each video (the "hook" or a single impressive demo moment)
- Post each Short 4-6 hours after the main video goes live
- Title format for Shorts: "Build [X] in [Y] seconds with Grit" or "One command = [result]"
- Shorts drive discovery; long-form drives retention and subscribers

### Thumbnail Consistency Guidelines
- **Background:** Dark navy gradient (#0b1120 to #0a0a0f)
- **Primary text:** Bold white or sky blue (#38bdf8), max 4-5 words, positioned left or center
- **Secondary text:** Smaller, muted (#9090a8), bottom edge
- **Visual element:** Always include one of: terminal screenshot, code snippet, app screenshot, or architecture diagram on the right side
- **Accent glow:** Soft purple (#6c5ce7) or sky blue (#38bdf8) bloom behind the visual element
- **Course badge:** Small colored pill in top-left corner — "WEB 0", "WEB 1", "DESKTOP 3", etc. — color-coded per section (blue for web, green for desktop, orange for mobile)
- **No faces** in thumbnails — developer/code-focused aesthetic
- **Font:** Bold sans-serif (Inter Black or similar), always ALL CAPS for primary text
- **Resolution:** 1280x720, safe zone: keep text away from edges by 80px

---

## WEEK 1: GRIT WEB — Introduction & Setup (Days 1-7)

---

### Day 1 — Course 0, Part 1: What is Grit?

**Title:** "The Setup Tax is Killing Your Projects — Grit Fixes It"
**Description:** Every Go project starts with the same grind: wiring auth, building an admin panel, setting up file storage, configuring Docker. Grit eliminates all of it. In this first video, we explain the problem Grit solves, why it exists, and what makes it different from every other Go framework.
**Purpose:** Awareness — establish the problem, introduce the framework, hook developers who feel the pain
**Timeline:**
- 0:00 — Hook: "You've started this project before. Auth, admin, storage, email — again."
- 0:45 — The setup tax: what it costs you every time
- 2:00 — What is Grit? One CLI, full-stack Go + React
- 3:30 — Why Go + React? The case for this stack
- 5:00 — What you get out of the box (feature flythrough)
- 6:30 — Who is Grit for? (Laravel devs, Next.js devs, Go devs)
- 8:00 — What this course series covers
- 9:00 — Subscribe + next video teaser

**Thumbnail Prompt:** YouTube thumbnail (1280x720). Dark navy gradient background (#0b1120 to #0a0a0f). Large bold white text "SETUP TAX" with a red strikethrough line across it, positioned center-left. Right side shows a faded montage of repetitive boilerplate code (auth, config, docker files) with a purple glow (#6c5ce7) bloom behind it. Small sky blue text at bottom: "There's a better way." Top-left: small blue pill badge "WEB 0". Clean, developer-focused.

---

### Day 2 — Course 0, Part 2: Tech Stack & Architecture

**Title:** "Grit Tech Stack Explained — Go, React, and 5 Architectures"
**Description:** Grit gives you five architecture modes: Triple, Double, Single, API-only, and Mobile. This video walks through the full tech stack — Go, Gin, GORM, Next.js, TanStack, Tailwind, shadcn/ui — and explains when to use each architecture mode for your project.
**Purpose:** Awareness / SEO — "go react tech stack", "go web app architecture", decision guide
**Timeline:**
- 0:00 — Hook: "Not every project needs the same structure"
- 0:40 — The tech stack layer by layer (Go, Gin, GORM, Postgres, Redis)
- 2:00 — Frontend stack (Next.js, TanStack Router, Tailwind, shadcn/ui)
- 3:15 — Monorepo structure (Turborepo + pnpm)
- 4:30 — Architecture mode 1: Triple (Web + Admin + API)
- 5:15 — Architecture mode 2: Double (Web + API)
- 5:50 — Architecture mode 3: Single (one binary)
- 6:30 — Architecture mode 4: API-only
- 7:00 — Architecture mode 5: Mobile (API + Expo)
- 7:45 — Decision flowchart: which mode should you pick?
- 8:45 — Philosophy: conventions over configuration
- 9:30 — Next video: we install and scaffold

**Thumbnail Prompt:** YouTube thumbnail (1280x720). Dark navy background. Center: five colored boxes in a horizontal row labeled "SINGLE", "DOUBLE", "TRIPLE", "API", "MOBILE", each with a different icon (one app, two apps, three apps, server, phone). Bold white text above: "PICK YOUR ARCHITECTURE". Purple and sky blue color scheme. Top-left: blue pill badge "WEB 0". Bottom-right: small "Go + React" text in muted gray.

---

### Day 3 — Course 1, Part 1: Install & Scaffold

**Title:** "Install Grit and Scaffold Your First App in 5 Minutes"
**Description:** Hands on the keyboard. We install Go, install Grit, run the interactive CLI, and scaffold a full-stack project. By the end of this video, you'll have a project folder with a Go API, Next.js frontend, admin panel, Docker config, and shared types — all generated from one command.
**Purpose:** SEO / Practical — "install grit framework", "go react tutorial", "grit new"
**Timeline:**
- 0:00 — Hook: "One command. Full-stack app. Let's go."
- 0:30 — Prerequisites check (Go, Node, pnpm, Docker)
- 1:15 — Install Grit: `go install` command
- 1:45 — Run `grit new myapp` — interactive prompts walkthrough
- 3:00 — Project structure tour: root files (docker-compose, turbo.json, pnpm-workspace)
- 4:00 — Tour: apps/api/ — Go backend structure
- 5:00 — Tour: apps/web/ — Next.js frontend structure
- 6:00 — Tour: apps/admin/ — Admin panel structure
- 7:00 — Tour: packages/shared/ — Zod schemas, TypeScript types
- 8:00 — Tour: grit.config.ts — project configuration
- 9:00 — Next: starting Docker and running the app

**Thumbnail Prompt:** YouTube thumbnail (1280x720). Dark navy background. Left side: bold yellow text "5 MIN SETUP" in large font. Right side: terminal window showing `grit new myapp` with colorful ASCII art output and a sky blue glow behind it. Bottom text in muted gray: "Go + React Full-Stack App". Top-left: blue pill badge "WEB 1".

---

### Day 4 — Course 1, Part 2: Docker, Dev Servers & App Tour

**Title:** "Start Your Grit App — Docker, Dev Servers, and Built-in Tools"
**Description:** We spin up Docker services (Postgres, Redis, MinIO, Mailhog), start the dev servers, and take a full tour of the running application. You'll see the web app, admin panel, API docs, GORM Studio, Pulse monitoring, and Sentinel rate limiting — all running out of the box.
**Purpose:** Practical guide — show the full developer experience, demonstrate value
**Timeline:**
- 0:00 — Hook: "Your app is scaffolded. Let's bring it to life."
- 0:30 — Docker compose up: starting Postgres, Redis, MinIO, Mailhog
- 1:30 — pnpm install and pnpm dev — starting all services
- 2:30 — Web app tour (localhost:3000): auth pages, dashboard
- 3:30 — Admin panel tour (localhost:3001): layout, sidebar, data tables
- 4:30 — API endpoints tour (localhost:8080): health, auth routes
- 5:30 — GORM Studio: visual database browser
- 6:15 — API Docs: auto-generated endpoint reference
- 7:00 — Pulse: real-time monitoring dashboard
- 7:45 — Sentinel: rate limiting dashboard
- 8:30 — Mailhog: email testing interface
- 9:15 — Recap: everything running from one scaffold

**Thumbnail Prompt:** YouTube thumbnail (1280x720). Dark navy background. Left side: bold white text "FULL APP TOUR". Right side: a grid of four small app screenshots (web dashboard, admin panel, GORM Studio, terminal) with a soft purple glow behind the grid. Bottom text: "Everything runs out of the box". Top-left: blue pill badge "WEB 1".

---

### Day 5 — Course 2, Part 1: First Resource Generation

**Title:** "Generate Full-Stack CRUD in 10 Seconds with Grit"
**Description:** The code generator is Grit's superpower. We run one command to generate a complete Product resource — Go model, service, handler, Zod schema, TypeScript types, React Query hooks, and admin page. Then we walk through every generated file to understand what just happened.
**Purpose:** SEO / Practical — "go crud generator", "generate rest api go", core feature demo
**Timeline:**
- 0:00 — Hook: "One command. Eight files. Full-stack CRUD."
- 0:30 — The command: `grit generate resource Product --fields "name:string price:float"`
- 1:15 — What it generated: file tree overview
- 2:00 — Go model: struct, GORM tags, auto-migration
- 3:00 — Go service: CRUD business logic, pagination
- 4:00 — Go handler: REST endpoints, request validation
- 5:00 — Zod schema: validation rules generated from Go types
- 5:45 — TypeScript types: interfaces matching Go structs
- 6:30 — React Query hooks: useProducts, useCreateProduct, etc.
- 7:15 — Admin page: DataTable and form, already wired up
- 8:00 — Route injection: how markers work behind the scenes
- 9:00 — Next: field types, modifiers, and advanced generation

**Thumbnail Prompt:** YouTube thumbnail (1280x720). Dark navy background. Left side: bold sky blue text "1 COMMAND" on top line, bold white text "8 FILES" on bottom line. Right side: a code editor showing a file tree with highlighted generated files (model.go, handler.go, schema.ts, hooks.ts) with a purple glow. Bottom text: "Full-Stack CRUD Generator". Top-left: blue pill badge "WEB 2".

---

### Day 6 — Course 2, Part 2: Field Types & Modifiers

**Title:** "Grit Field Types — Strings, Slugs, Enums, and Modifiers"
**Description:** Grit supports 12+ field types and modifiers like :unique, :optional, and :default. We walk through every type — string, text, int, float, bool, datetime, enum, json, slug — and show how modifiers change the generated code in Go, Zod, and TypeScript.
**Purpose:** SEO / Practical — "gorm field types", "go struct tags", reference guide
**Timeline:**
- 0:00 — Hook: "Your fields control everything — model, validation, and UI"
- 0:30 — Basic types: string, text, int, float, bool
- 1:30 — Date types: datetime, date, time
- 2:30 — Special types: enum, json, slug
- 3:30 — The slug field: auto-generated from another field
- 4:15 — Modifier: :unique — what it changes in GORM and Zod
- 5:00 — Modifier: :optional — nullable fields, pointer types in Go
- 5:45 — Modifier: :default — default values in migration and schema
- 6:30 — Combining modifiers: `email:string:unique:optional`
- 7:15 — Live demo: generate a Blog resource with mixed fields
- 8:15 — Inspecting the generated code: Go struct tags, Zod rules
- 9:00 — Next: relationships and removing resources

**Thumbnail Prompt:** YouTube thumbnail (1280x720). Dark navy background. Left side: bold white text "FIELD TYPES" with a list of colored type labels below (string, int, float, bool, enum, slug) each in small colored pills. Right side: a Go struct code snippet showing field declarations with GORM tags, sky blue syntax highlighting. Top-left: blue pill badge "WEB 2".

---

### Day 7 — Course 2, Part 3: Relationships & Resource Removal

**Title:** "Grit Relationships — belongs_to, many_to_many, and Clean Undo"
**Description:** Resources don't live in isolation. We add belongs_to and many_to_many relationships, see how Grit generates foreign keys, join tables, and nested queries. Then we use `grit remove` to cleanly undo a resource — every file, every injection, gone.
**Purpose:** Practical guide — relationships are the make-or-break feature of any generator
**Timeline:**
- 0:00 — Hook: "Real apps have relationships. Grit handles them."
- 0:30 — belongs_to: adding a Category to our Product
- 1:30 — What belongs_to generates: foreign key, preload, nested response
- 2:30 — many_to_many: adding Tags to our Product
- 3:30 — What many_to_many generates: join table, association methods
- 4:30 — Frontend: how relationships appear in forms and tables
- 5:30 — Type sync: `grit sync` — keeping Go and TypeScript in lockstep
- 6:30 — Removing a resource: `grit remove resource Product`
- 7:15 — What removal cleans up: files, route injections, imports
- 8:00 — Round-trip demo: generate, verify, remove, verify clean state
- 8:45 — Recap: the full generation lifecycle
- 9:15 — Next week: authentication deep dive

**Thumbnail Prompt:** YouTube thumbnail (1280x720). Dark navy background. Center: a diagram showing three connected boxes — "Product", "Category", "Tag" — with arrows labeled "belongs_to" and "many_to_many" in sky blue. Bold white text at top: "RELATIONSHIPS". Subtle purple glow behind the diagram. Bottom text: "Generate, Connect, Remove". Top-left: blue pill badge "WEB 2".

---

## WEEK 2: GRIT WEB — Auth, Admin, Storage (Days 8-14)

---

### Day 8 — Course 3, Part 1: JWT Authentication Flow

**Title:** "JWT Auth in Go — Registration, Login, Tokens Explained"
**Description:** Grit ships with a complete JWT authentication system. We walk through the registration flow, login flow, access and refresh tokens, middleware protection, and how the frontend stores and refreshes tokens automatically.
**Purpose:** SEO — "go jwt authentication", "gin jwt middleware", "go register login api"
**Timeline:**
- 0:00 — Hook: "Auth is the first thing you need and the last thing you want to build"
- 0:30 — What Grit generates: auth handler, JWT service, middleware
- 1:30 — Registration flow: request validation, password hashing, token generation
- 2:30 — Login flow: credential check, token pair creation
- 3:30 — Access tokens vs refresh tokens: why both, expiration strategy
- 4:30 — Auth middleware: protecting routes, extracting user from token
- 5:30 — Token refresh: how the frontend auto-refreshes expired tokens
- 6:30 — Testing with curl: register, login, access protected route
- 7:30 — Frontend auth context: useAuth hook, token storage
- 8:30 — Logout and token invalidation
- 9:00 — Next: roles, 2FA, and social login

**Thumbnail Prompt:** YouTube thumbnail (1280x720). Dark navy background. Left side: bold white text "JWT AUTH" with a lock icon in sky blue above it. Right side: a flow diagram showing User -> Login -> Access Token -> Protected Route, with arrows in sky blue. Bottom text: "Go + React Authentication". Top-left: blue pill badge "WEB 3".

---

### Day 9 — Course 3, Part 2: Roles, 2FA & Trusted Devices

**Title:** "RBAC Roles and TOTP 2FA — Securing Your Grit App"
**Description:** We add role-based access control with admin, editor, and user roles. Then we enable TOTP two-factor authentication with QR code setup, backup codes, and trusted device management. Your Grit app goes from basic auth to production-grade security.
**Purpose:** SEO — "go rbac", "totp 2fa go", "two factor authentication golang"
**Timeline:**
- 0:00 — Hook: "Basic auth isn't enough. Let's add roles and 2FA."
- 0:30 — Role-based access control: how Grit implements RBAC
- 1:30 — Built-in roles: admin, editor, user — what each can do
- 2:30 — Role middleware: restricting routes by role
- 3:30 — Custom roles: adding your own role definitions
- 4:15 — TOTP 2FA: what it is, why it matters
- 5:00 — Setup flow: generating secret, QR code, verification
- 6:00 — Login with 2FA: the two-step authentication flow
- 7:00 — Backup codes: generation, storage, one-time use
- 7:45 — Trusted devices: skip 2FA on recognized devices
- 8:30 — Frontend: 2FA setup page, verification modal
- 9:15 — Next: OAuth2 social login

**Thumbnail Prompt:** YouTube thumbnail (1280x720). Dark navy background. Left side: bold white text "ROLES + 2FA". Right side: a shield icon with a checkmark in sky blue, surrounded by small icons representing roles (admin crown, editor pencil, user silhouette). A QR code graphic partially visible behind the shield with purple glow. Top-left: blue pill badge "WEB 3".

---

### Day 10 — Course 3, Part 3: OAuth2 Social Login

**Title:** "Add Google & GitHub Login to Your Go App in Minutes"
**Description:** Grit includes OAuth2 social login out of the box. We configure Google and GitHub providers, walk through the OAuth2 flow (redirect, callback, account linking), and show the frontend social login buttons — all pre-built and ready to customize.
**Purpose:** SEO — "go oauth2", "google login golang", "github oauth go gin"
**Timeline:**
- 0:00 — Hook: "Your users expect social login. Here's how Grit does it."
- 0:30 — OAuth2 flow overview: redirect, consent, callback, token
- 1:30 — Configuring Google OAuth: client ID, secret, redirect URI
- 2:30 — Configuring GitHub OAuth: same pattern, different provider
- 3:30 — The callback handler: what happens when the user returns
- 4:30 — Account linking: connecting social accounts to existing users
- 5:30 — New user flow: auto-registration from social profile
- 6:15 — Frontend: social login buttons, loading states, error handling
- 7:15 — Custom roles for social users: default role assignment
- 8:00 — Adding more providers: the provider interface pattern
- 8:45 — Full demo: sign up with GitHub, see the user in admin
- 9:15 — Next week: the admin panel deep dive

**Thumbnail Prompt:** YouTube thumbnail (1280x720). Dark navy background. Center: Google "G" logo and GitHub octocat logo side by side, connected to a Grit app icon with arrows. Bold white text above: "SOCIAL LOGIN". Sky blue glow behind the logos. Bottom text: "OAuth2 in Minutes". Top-left: blue pill badge "WEB 3".

---

### Day 11 — Course 4, Part 1: Admin Panel & DataTable

**Title:** "Grit Admin Panel — Dashboard Widgets and DataTable Deep Dive"
**Description:** The admin panel ships fully functional. We tour the dashboard with its stats cards, charts, and activity feed. Then we deep dive into DataTable — sorting, filtering, column selection, pagination, bulk actions — all generated automatically for every resource.
**Purpose:** Practical guide — admin panels are a major selling point
**Timeline:**
- 0:00 — Hook: "Every app needs an admin panel. Grit builds it for you."
- 0:30 — Admin panel overview: layout, sidebar, navigation
- 1:30 — Dashboard widgets: stats cards with trends
- 2:30 — Dashboard charts: configuring chart data
- 3:15 — Activity feed: recent actions log
- 4:00 — DataTable introduction: generated for every resource
- 4:45 — Sorting: click column headers, multi-sort
- 5:30 — Filtering: text search, column-specific filters
- 6:15 — Column selection: show/hide columns, persist preference
- 7:00 — Pagination: server-side, configurable page size
- 7:45 — Bulk actions: select rows, bulk delete
- 8:30 — Customizing DataTable: adding custom columns, actions
- 9:15 — Next: FormBuilder and multi-step forms

**Thumbnail Prompt:** YouTube thumbnail (1280x720). Dark navy background (#0a0a0f). Left side: bold white text "ADMIN PANEL". Right side: a screenshot of the Grit admin dashboard showing stats cards and a data table with rows of data, purple accent highlights on active elements. Sky blue glow behind the screenshot. Top-left: blue pill badge "WEB 4".

---

### Day 12 — Course 4, Part 2: FormBuilder & Variants

**Title:** "Grit FormBuilder — Multi-Step Forms and Style Variants"
**Description:** FormBuilder generates forms from your resource definition — text inputs, selects, date pickers, rich text editors, file uploads. We build a multi-step form, explore the three style variants (default, compact, floating), and show how to use FormBuilder as a standalone component.
**Purpose:** Practical guide — forms are tedious; showing automation saves time
**Timeline:**
- 0:00 — Hook: "Forms are boring to build. So we automated them."
- 0:30 — FormBuilder overview: what it generates from your fields
- 1:30 — Field types: text, textarea, select, date, toggle, file upload
- 2:30 — Rich text editor: built-in for text fields
- 3:15 — Validation: Zod schemas driving form validation
- 4:00 — Multi-step forms: splitting fields across steps
- 5:00 — Step navigation: progress indicator, validation per step
- 5:45 — Style variant: default — standard form layout
- 6:15 — Style variant: compact — dense, admin-friendly
- 6:45 — Style variant: floating — floating labels, modern feel
- 7:30 — Standalone usage: using FormBuilder outside the admin
- 8:15 — Custom fields: extending FormBuilder with your own inputs
- 9:00 — Next: file storage and image processing

**Thumbnail Prompt:** YouTube thumbnail (1280x720). Dark navy background. Left side: bold white text "FORMBUILDER" with three small form mockups stacked below (default, compact, floating label styles). Right side: a multi-step form UI with step indicators and form fields, purple accent on the active step. Top-left: blue pill badge "WEB 4".

---

### Day 13 — Course 5, Part 1: File Storage & Presigned URLs

**Title:** "S3 File Storage in Go — MinIO, Presigned URLs, Setup Guide"
**Description:** Grit includes a complete file storage system using S3-compatible storage. We set up MinIO locally, explain presigned URLs for secure direct uploads, and walk through the storage service — upload, download, delete, list — all pre-built in your scaffolded app.
**Purpose:** SEO — "go s3 upload", "minio golang", "presigned url go"
**Timeline:**
- 0:00 — Hook: "File uploads shouldn't take a week to set up"
- 0:30 — Storage architecture: why S3-compatible, why presigned URLs
- 1:30 — MinIO setup: already in your Docker Compose
- 2:15 — MinIO console tour: buckets, objects, access keys
- 3:00 — Storage service walkthrough: the Go code Grit generates
- 4:00 — Presigned URLs explained: how they work, why they're secure
- 5:00 — Upload flow: frontend requests URL, uploads directly to storage
- 6:00 — Download flow: generating presigned download URLs
- 6:45 — List and delete operations
- 7:30 — Upload model: tracking files in the database
- 8:15 — Environment config: bucket names, regions, credentials
- 9:00 — Next: the upload UI and image processing

**Thumbnail Prompt:** YouTube thumbnail (1280x720). Dark navy background. Left side: bold white text "FILE STORAGE" with an upload cloud icon in sky blue. Right side: a diagram showing Browser -> Presigned URL -> S3 Bucket, with arrows and the MinIO logo. Purple glow behind the diagram. Top-left: blue pill badge "WEB 5".

---

### Day 14 — Course 5, Part 2: Upload UI & Cloud Migration

**Title:** "Image Processing and Switching to Cloudflare R2 or AWS S3"
**Description:** We build the upload flow end-to-end: drag-and-drop UI, progress bar, image processing (resize, thumbnail generation), and file type validation. Then we switch from local MinIO to Cloudflare R2 — just change three environment variables.
**Purpose:** Practical guide — "cloudflare r2 golang", "go image processing", real-world upload flow
**Timeline:**
- 0:00 — Hook: "Local dev with MinIO, production with R2 — same code"
- 0:30 — Upload component: drag-and-drop, file selection, preview
- 1:30 — Progress tracking: upload progress bar implementation
- 2:15 — File type validation: allowed types, size limits
- 3:00 — Image processing: automatic resize on upload
- 3:45 — Thumbnail generation: creating smaller versions
- 4:30 — Processing pipeline: upload -> process -> store -> record
- 5:30 — Switching to Cloudflare R2: three env vars to change
- 6:15 — R2 configuration: account ID, access key, bucket
- 7:00 — Switching to AWS S3: same pattern, different credentials
- 7:45 — Backblaze B2: another option, same interface
- 8:15 — Production tips: CDN, caching, signed URLs expiration
- 9:00 — Next week: background jobs and AI

**Thumbnail Prompt:** YouTube thumbnail (1280x720). Dark navy background. Left side: bold white text "UPLOAD + CLOUD". Right side: three cloud provider logos (MinIO, Cloudflare R2, AWS S3) connected by arrows to a single Go gopher icon, showing interchangeability. Sky blue glow behind the logos. Bottom text: "Same Code, Any Provider". Top-left: blue pill badge "WEB 5".

---

## WEEK 3: GRIT WEB — Jobs, AI, Deploy (Days 15-21)

---

### Day 15 — Course 6, Part 1: Background Jobs with Asynq

**Title:** "Background Jobs in Go with Asynq — Workers, Queues, Retries"
**Description:** Grit includes asynq for background job processing. We explore the built-in workers (email, image processing, cleanup), create a custom job from scratch, and monitor everything from the admin jobs dashboard. No external job servers — Redis handles it all.
**Purpose:** SEO — "go background jobs", "asynq golang tutorial", "redis job queue go"
**Timeline:**
- 0:00 — Hook: "Some work doesn't belong in the request cycle"
- 0:30 — Why background jobs: email, image processing, reports, cleanup
- 1:30 — Asynq overview: how it uses Redis as a job broker
- 2:15 — Built-in workers: email worker, image processing worker, cleanup worker
- 3:15 — Job lifecycle: enqueue, process, retry, dead letter
- 4:15 — Creating a custom job: task definition, handler, registration
- 5:15 — Enqueuing jobs: from handlers, from services, with options
- 6:00 — Retry configuration: max retries, backoff strategy
- 6:45 — Queue priorities: critical, default, low
- 7:30 — Admin jobs dashboard: monitoring active, pending, failed jobs
- 8:15 — Inspecting failed jobs: error details, retry button
- 9:00 — Next: cron scheduler and email service

**Thumbnail Prompt:** YouTube thumbnail (1280x720). Dark navy background. Left side: bold white text "BACKGROUND JOBS". Right side: a queue visualization showing tasks flowing from "Enqueue" through "Process" to "Complete" with retry arrows, in sky blue and purple. Redis logo small in the corner. Top-left: blue pill badge "WEB 6".

---

### Day 16 — Course 6, Part 2: Cron Scheduler & Email Service

**Title:** "Cron Jobs and Email Templates in Go — Scheduled Tasks Made Easy"
**Description:** We set up the cron scheduler for recurring tasks — daily cleanup, weekly reports, session pruning. Then we explore the email service: Resend integration, four pre-built HTML templates, and Mailhog for local testing. Schedule it, send it, preview it.
**Purpose:** SEO — "go cron jobs", "send email golang", "resend go tutorial"
**Timeline:**
- 0:00 — Hook: "Cron jobs and emails — two things every app needs, zero fun to build"
- 0:30 — Cron scheduler: how Grit integrates asynq's periodic tasks
- 1:30 — Built-in cron jobs: cleanup expired sessions, prune old files
- 2:15 — Adding a custom cron job: schedule expression, handler
- 3:00 — Cron expressions explained: @daily, @weekly, custom patterns
- 3:45 — Admin cron dashboard: viewing scheduled tasks, next run times
- 4:30 — Email service overview: Resend as the provider
- 5:15 — Email templates: welcome, password reset, verification, notification
- 6:15 — Template structure: Go html/template with inline CSS
- 7:00 — Sending email from code: service method, template data
- 7:45 — Mailhog: testing emails locally without sending real mail
- 8:30 — Mail preview page: viewing templates in the admin panel
- 9:00 — Next: AI integration

**Thumbnail Prompt:** YouTube thumbnail (1280x720). Dark navy background. Left side: bold white text "CRON + EMAIL". Right side: split visual — top half shows a clock/calendar icon with cron expressions, bottom half shows an email template preview with purple and sky blue accents. Top-left: blue pill badge "WEB 6".

---

### Day 17 — Course 7, Part 1: AI Gateway Setup

**Title:** "Add AI to Your Go App — Claude and OpenAI Gateway"
**Description:** Grit includes an AI gateway that supports Claude and OpenAI models behind a unified API. We set up the gateway, make completion requests, and explore the chat endpoint — all through Grit's pre-built AI service. Bring your own API key, pick your model, and go.
**Purpose:** SEO — "go ai integration", "claude api golang", "openai go tutorial"
**Timeline:**
- 0:00 — Hook: "Every app needs AI now. Here's the easiest way to add it."
- 0:30 — AI gateway architecture: unified API over multiple providers
- 1:30 — Configuration: API keys, model selection, environment variables
- 2:15 — The AI service: what Grit generates for you
- 3:00 — Complete endpoint: single prompt, single response
- 4:00 — Making a completion request: curl demo
- 4:45 — Chat endpoint: multi-turn conversation with history
- 5:30 — Chat request structure: messages array, system prompt
- 6:15 — Model switching: Claude to OpenAI with one config change
- 7:00 — Token counting and rate limiting
- 7:45 — Error handling: API errors, fallback strategies
- 8:30 — Admin AI dashboard: usage monitoring
- 9:00 — Next: SSE streaming and building a chat UI

**Thumbnail Prompt:** YouTube thumbnail (1280x720). Dark navy background. Left side: bold white text "AI GATEWAY". Right side: Claude and OpenAI logos connected by arrows to a central Go gopher icon acting as a gateway. Purple and sky blue glow effects. Bottom text: "One API, Any Model". Top-left: blue pill badge "WEB 7".

---

### Day 18 — Course 7, Part 2: SSE Streaming & Chat UI

**Title:** "Build a Streaming Chat UI with Go SSE and React"
**Description:** We implement server-sent events for real-time AI streaming — tokens arrive one at a time, just like ChatGPT. Then we build a React chat interface with message history, typing indicators, and model selection. The full AI chat experience in your own app.
**Purpose:** SEO — "go sse streaming", "server sent events golang", "react chat ui ai"
**Timeline:**
- 0:00 — Hook: "Streaming responses, just like ChatGPT — in your Grit app"
- 0:30 — SSE explained: server-sent events vs WebSockets
- 1:30 — Streaming endpoint: how the Go handler sends token chunks
- 2:30 — Frontend SSE client: EventSource, parsing chunks
- 3:30 — Building the chat UI: message list, input area, send button
- 4:30 — Message rendering: markdown support, code blocks
- 5:15 — Typing indicator: showing "AI is thinking..." animation
- 6:00 — Conversation history: storing and sending previous messages
- 6:45 — Model selector: switching between Claude and OpenAI models
- 7:30 — System prompt configuration: customizing AI behavior
- 8:15 — Full demo: multi-turn conversation with streaming
- 9:00 — Next: deployment

**Thumbnail Prompt:** YouTube thumbnail (1280x720). Dark navy background. Left side: bold white text "STREAMING AI". Right side: a chat interface mockup showing a conversation with streaming text (partially rendered last message), sky blue message bubbles on dark background. Cursor blinking effect on the last line. Top-left: blue pill badge "WEB 7".

---

### Day 19 — Course 8, Part 1: Deployment Concepts & Build Pipeline

**Title:** "Deploy Your Grit App — Build Pipeline and grit deploy"
**Description:** Time to go live. We cover deployment concepts: building the Go binary, generating the frontend static assets, and using the `grit deploy` command to package everything. We walk through the build pipeline step by step — from source code to production artifacts.
**Purpose:** Practical guide — "deploy go app", "go production build", deployment is where frameworks fail
**Timeline:**
- 0:00 — Hook: "Building it is the fun part. Deploying it is where frameworks fail."
- 0:30 — Deployment overview: what needs to happen
- 1:30 — Building the Go API: cross-compilation, binary output
- 2:30 — Building the frontend: Next.js build, static export options
- 3:15 — Building the admin panel: same process, separate artifact
- 4:00 — The grit deploy command: what it does under the hood
- 4:45 — Environment variables: production config, secrets management
- 5:30 — Database migrations: running in production
- 6:15 — Health checks: the /health endpoint and readiness probes
- 7:00 — Build pipeline: CI/CD with GitHub Actions (generated workflow)
- 8:00 — Artifact overview: what you're deploying, where it goes
- 9:00 — Next: systemd, Caddy, and Docker deployment

**Thumbnail Prompt:** YouTube thumbnail (1280x720). Dark navy background. Left side: bold white text "DEPLOY". Right side: a pipeline diagram showing Code -> Build -> Test -> Deploy with arrows, each stage in a colored box (blue, purple, sky blue, green). A rocket icon at the end. Top-left: blue pill badge "WEB 8".

---

### Day 20 — Course 8, Part 2: systemd, Caddy & Docker Deploy

**Title:** "Go to Production — systemd, Caddy Reverse Proxy, Docker Deploy"
**Description:** Three deployment strategies in one video. We deploy with systemd and Caddy on a VPS (the simplest path), then show Docker Compose deployment for teams, and cover maintenance mode for zero-downtime updates. Your Grit app, live on the internet.
**Purpose:** SEO — "deploy go app vps", "caddy reverse proxy go", "docker compose production"
**Timeline:**
- 0:00 — Hook: "Three ways to deploy. Pick the one that fits."
- 0:30 — Strategy 1: VPS with systemd — the simplest deployment
- 1:15 — Writing the systemd service file
- 2:00 — Caddy reverse proxy: automatic HTTPS, simple config
- 2:45 — Caddy configuration: domain, proxy pass, static files
- 3:30 — Starting and managing the service: enable, start, logs
- 4:15 — Strategy 2: Docker Compose deployment
- 5:00 — Production docker-compose: API, frontend, Postgres, Redis
- 5:45 — Docker networking: service communication
- 6:30 — Strategy 3: maintenance mode for updates
- 7:15 — Zero-downtime update process
- 8:00 — Monitoring in production: logs, health checks, alerts
- 8:45 — Your app is live: full demo of deployed application
- 9:15 — Next: bonus SaaS build

**Thumbnail Prompt:** YouTube thumbnail (1280x720). Dark navy background. Left side: bold white text "PRODUCTION". Right side: three deployment icons stacked — a server rack (systemd), a Caddy logo, and a Docker whale — connected by dotted lines. Green "LIVE" indicator dot in top right. Top-left: blue pill badge "WEB 8".

---

### Day 21 — BONUS: Build a Complete SaaS with Grit

**Title:** "Build a Complete SaaS in 30 Minutes with Grit Framework"
**Description:** Everything we've learned, all in one build. We scaffold a SaaS project, generate resources, configure auth with roles and 2FA, set up file storage, add background jobs, integrate AI, and deploy — a real, production-ready SaaS application built from scratch with Grit.
**Purpose:** Awareness / SEO — "build saas with go", "go saas tutorial", capstone video
**Timeline:**
- 0:00 — Hook: "A complete SaaS. From zero to deployed. Let's build."
- 0:30 — Project scaffold: `grit new saas-app` with Triple mode
- 1:15 — Generate resources: Users, Projects, Tasks, Comments
- 2:30 — Relationships: Projects has_many Tasks, Tasks has_many Comments
- 3:30 — Auth configuration: roles (admin, member), 2FA enabled
- 4:30 — File storage: project attachments with image processing
- 5:30 — Background jobs: email notifications on task assignment
- 6:15 — AI integration: task summarization with Claude
- 7:00 — Admin panel: managing users, monitoring jobs
- 7:45 — Deployment: `grit deploy` and Docker Compose
- 8:30 — Live demo: the finished SaaS running in production
- 9:15 — Recap: every Grit feature, one project, 30 minutes

**Thumbnail Prompt:** YouTube thumbnail (1280x720). Dark navy background. Center: bold white text "BUILD A SAAS" on top, "30 MINUTES" in sky blue below. Right side: a polished SaaS dashboard screenshot with purple accent elements. Left bottom: small icons representing auth, storage, jobs, AI, deploy in a row. Top-left: yellow pill badge "BONUS".

---

## WEEK 4: GRIT DESKTOP & MOBILE (Days 22-30)

---

### Day 22 — Desktop Course 1: First Desktop App

**Title:** "Build a Desktop App with Go and Wails — Grit Desktop"
**Description:** Grit isn't just for web. We scaffold a desktop app with `grit new-desktop`, explore the Wails + Go + React stack, set up SQLite for local storage, and run the app in dev mode. A native desktop application with web technologies — no Electron, no bloat.
**Purpose:** SEO — "go desktop app", "wails tutorial", "wails golang react"
**Timeline:**
- 0:00 — Hook: "Desktop apps with Go. No Electron. No bloat."
- 0:30 — What is Grit Desktop? Wails v2 + Go + React
- 1:15 — Why Wails over Electron: native webview, small binary, Go backend
- 2:00 — Scaffold: `grit new-desktop my-app`
- 2:45 — Project structure: Go backend, React frontend, Wails config
- 3:30 — Go files: main.go, app.go, Wails bootstrap
- 4:15 — Database: GORM with SQLite, auto-migration
- 5:00 — Models: User, Blog, Contact — pre-generated
- 5:45 — Frontend: React + TypeScript + Tailwind + TanStack Router
- 6:30 — Dev mode: `wails dev` — live reload, DevTools
- 7:15 — Tour: the running desktop app
- 8:00 — Wails bindings: how Go functions become JavaScript calls
- 9:00 — Next: CRUD operations and bindings

**Thumbnail Prompt:** YouTube thumbnail (1280x720). Dark navy background. Left side: bold white text "DESKTOP APP" with "No Electron" in smaller red strikethrough text below. Right side: a desktop window frame showing a polished dark app UI with sidebar navigation. Go gopher icon in the corner. Top-left: green pill badge "DESKTOP 1".

---

### Day 23 — Desktop Course 2: Desktop CRUD & Bindings

**Title:** "Desktop CRUD with Wails Bindings and TanStack Router"
**Description:** We generate a resource for our desktop app, explore how Wails bindings expose Go functions to the React frontend, and set up TanStack Router for client-side navigation. Full CRUD — create, read, update, delete — with Go doing the heavy lifting and React rendering the UI.
**Purpose:** Practical guide — "wails bindings", "wails crud", "tanstack router wails"
**Timeline:**
- 0:00 — Hook: "Call Go functions directly from React. No REST API needed."
- 0:30 — Wails bindings explained: Go methods become TypeScript functions
- 1:30 — Blog service: CRUD methods bound to the frontend
- 2:30 — Generated TypeScript bindings: auto-typed from Go structs
- 3:30 — TanStack Router setup: routes, layouts, navigation
- 4:30 — Blog list page: calling Go, rendering in React
- 5:30 — Create form: calling the Go create method
- 6:15 — Edit and delete: update and remove operations
- 7:00 — Contact service: second resource, same pattern
- 7:45 — Error handling: Go errors surfacing in the UI
- 8:30 — Data flow recap: React -> Wails Bridge -> Go -> SQLite -> back
- 9:00 — Next: custom UI and frameless window

**Thumbnail Prompt:** YouTube thumbnail (1280x720). Dark navy background. Left side: bold white text "WAILS CRUD". Right side: a split view showing Go code on top and React/TypeScript code on bottom, connected by an arrow labeled "Bindings" in sky blue. Purple glow between them. Top-left: green pill badge "DESKTOP 2".

---

### Day 24 — Desktop Course 3: Custom UI & Frameless Window

**Title:** "Frameless Desktop Window with Custom Title Bar and Themes"
**Description:** We make the desktop app look native and polished. Custom frameless window with a draggable title bar, window controls (minimize, maximize, close), dark and light theme toggle, and 12 shadcn-style UI components — all pre-built in your Grit Desktop scaffold.
**Purpose:** Practical guide — "wails frameless window", "custom title bar wails", UI polish
**Timeline:**
- 0:00 — Hook: "Make your desktop app look like it belongs on the platform"
- 0:30 — Frameless window: removing the default title bar
- 1:15 — Custom title bar: draggable area, window controls
- 2:00 — Window controls: minimize, maximize, close — calling Wails runtime
- 2:45 — Draggable regions: CSS flags for drag behavior
- 3:30 — Sidebar navigation: collapsible, icon-based, route-aware
- 4:15 — App layout: sidebar + header + content area
- 5:00 — Theme toggle: dark and light mode with system detection
- 5:45 — shadcn-style UI components: button, input, card, dialog
- 6:30 — DataTable component: sorting, filtering for desktop data
- 7:15 — FormBuilder: reusable form generation
- 8:00 — Draggable panels: resizable split views
- 8:45 — Full UI tour: the polished desktop experience
- 9:15 — Next: export features and distribution

**Thumbnail Prompt:** YouTube thumbnail (1280x720). Dark navy background. Left side: bold white text "FRAMELESS UI". Right side: a desktop app screenshot showing a frameless window with custom title bar, sidebar, and dark theme. Window control buttons (red, yellow, green dots) visible in the title bar. Purple accent glow. Top-left: green pill badge "DESKTOP 3".

---

### Day 25 — Desktop Courses 4-5: Export & Distribution

**Title:** "PDF, Excel, CSV Export + Build and Distribute Your Desktop App"
**Description:** Two courses in one video. First: export data from your desktop app to PDF, Excel, and CSV using Go libraries — one-click export buttons in the UI. Second: build your Wails app into native binaries for Windows, macOS, and Linux, and prepare for distribution.
**Purpose:** Practical guide — "go pdf export", "wails build", "distribute wails app"
**Timeline:**
- 0:00 — Hook: "Export anything. Ship everywhere."
- 0:30 — Export overview: PDF, Excel, CSV from Go services
- 1:15 — PDF export: generating formatted PDF reports
- 2:00 — Excel export: creating .xlsx files with styled sheets
- 2:45 — CSV export: simple, universal data export
- 3:30 — Export UI: one-click buttons, file save dialog
- 4:15 — Export service architecture: Go handles generation, Wails saves
- 5:00 — Building for production: `wails build`
- 5:30 — Windows build: .exe, installer options
- 6:00 — macOS build: .app bundle, code signing basics
- 6:30 — Linux build: binary, AppImage
- 7:00 — Build configuration: wails.json options, app icon, metadata
- 7:45 — Binary size: what to expect, optimization tips
- 8:15 — Distribution: where to host, update strategies
- 8:45 — Recap: from scaffold to shipped desktop app
- 9:15 — Next: mobile development with Grit

**Thumbnail Prompt:** YouTube thumbnail (1280x720). Dark navy background. Left side: bold white text "EXPORT + SHIP". Right side: three file icons (PDF in red, Excel in green, CSV in blue) on top, and three platform icons (Windows, macOS, Linux) on bottom, separated by a horizontal line. Purple glow. Top-left: green pill badge "DESKTOP 4-5".

---

### Day 26 — Mobile Course 1: First Mobile App

**Title:** "Build a Mobile App with Go Backend — Grit + Expo"
**Description:** Grit's mobile mode scaffolds an Expo React Native app connected to your Go API. We create the project, explore the file structure, configure the development environment, and run the app on a physical device and emulator. Full-stack mobile, powered by Go.
**Purpose:** SEO — "expo react native go backend", "go mobile app", "expo go tutorial"
**Timeline:**
- 0:00 — Hook: "Your Go backend. A native mobile app. Let's connect them."
- 0:30 — Grit Mobile overview: Expo + Go API architecture
- 1:15 — Scaffold: `grit new my-mobile-app --mobile` (API + Expo)
- 2:00 — Project structure: Go API in apps/api, Expo in apps/mobile
- 2:45 — Expo project tour: app directory, components, hooks, config
- 3:30 — Development setup: Expo Go, iOS simulator, Android emulator
- 4:15 — Running on device: scanning QR code with Expo Go
- 5:00 — Running on emulator: iOS and Android side by side
- 5:45 — API connection: base URL configuration, environment-based
- 6:30 — Project config: app.json, EAS configuration
- 7:15 — Navigation structure: Expo Router file-based routing
- 8:00 — Shared types: same Zod schemas, same TypeScript types
- 8:45 — Dev workflow: hot reload, debugging, logs
- 9:15 — Next: auth and navigation

**Thumbnail Prompt:** YouTube thumbnail (1280x720). Dark navy background. Left side: bold white text "MOBILE APP" with "Go + Expo" in sky blue below. Right side: a phone mockup showing a mobile app screen with dark theme, next to a terminal running the Go API. Purple glow behind the phone. Top-left: orange pill badge "MOBILE 1".

---

### Day 27 — Mobile Course 2: Auth & Navigation

**Title:** "Mobile Auth with SecureStore and Expo Router Navigation"
**Description:** We implement mobile authentication — login, register, token storage with Expo SecureStore, and automatic token refresh. Then we set up Expo Router with tab navigation, stack navigation, and auth-guarded routes. The mobile app now has real user sessions.
**Purpose:** SEO — "expo securestore auth", "expo router tabs", "react native jwt auth"
**Timeline:**
- 0:00 — Hook: "Mobile auth is different. Tokens go in SecureStore, not cookies."
- 0:30 — Auth architecture: mobile vs web token storage
- 1:15 — SecureStore: encrypted storage for tokens on device
- 2:00 — Login screen: form, validation, API call to Go backend
- 2:45 — Register screen: same pattern, different endpoint
- 3:30 — Token management: storing, retrieving, refreshing
- 4:15 — Auth context: useAuth hook for the mobile app
- 5:00 — Auto-refresh: intercepting 401s, refreshing tokens silently
- 5:45 — Expo Router: file-based routing for mobile
- 6:30 — Tab navigation: bottom tabs with icons
- 7:15 — Stack navigation: push/pop within tabs
- 7:45 — Auth guard: redirecting unauthenticated users to login
- 8:30 — Protected routes: checking auth state before rendering
- 9:00 — Next: API integration with TanStack Query

**Thumbnail Prompt:** YouTube thumbnail (1280x720). Dark navy background. Left side: bold white text "MOBILE AUTH". Right side: a phone mockup showing a login screen with email and password fields, dark theme. A lock icon with a checkmark above the phone. Sky blue accent on form elements. Top-left: orange pill badge "MOBILE 2".

---

### Day 28 — Mobile Course 3: API Integration & Offline

**Title:** "TanStack Query on Mobile — Pull to Refresh, Offline Cache"
**Description:** We connect the mobile app to the Go API using TanStack Query — the same data fetching library used in the web app. Pull-to-refresh, infinite scroll, optimistic updates, and offline caching. Your mobile app stays fast and responsive even with a flaky connection.
**Purpose:** SEO — "react native tanstack query", "expo offline cache", "pull to refresh react native"
**Timeline:**
- 0:00 — Hook: "Same React Query skills. Native mobile experience."
- 0:30 — TanStack Query setup: QueryClient, provider, config
- 1:15 — Fetching data: useQuery with the Go API
- 2:00 — Custom hooks: useProducts, useUser — same pattern as web
- 2:45 — Pull-to-refresh: connecting RefreshControl to refetch
- 3:30 — Infinite scroll: paginated lists with useInfiniteQuery
- 4:15 — Creating data: useMutation for POST requests
- 5:00 — Optimistic updates: instant UI, reconcile with server
- 5:45 — Error handling: retry logic, error boundaries
- 6:30 — Offline caching: persisting query cache to AsyncStorage
- 7:15 — Stale data: showing cached data while fetching fresh
- 8:00 — Network detection: online/offline indicators
- 8:30 — Full demo: browsing, creating, pull-to-refresh, going offline
- 9:00 — Next: push notifications and app store

**Thumbnail Prompt:** YouTube thumbnail (1280x720). Dark navy background. Left side: bold white text "API + OFFLINE". Right side: a phone mockup showing a data list with a pull-to-refresh spinner, and a small "offline" badge in the corner. Arrows between phone and a cloud labeled "Go API". Top-left: orange pill badge "MOBILE 3".

---

### Day 29 — Mobile Courses 4-5: Push Notifications & App Store

**Title:** "Push Notifications and App Store Submission with Expo EAS"
**Description:** Two courses, one video. First: push notifications with Expo — registering tokens, sending from the Go backend, handling taps. Second: building with EAS, configuring app metadata, and submitting to the Apple App Store and Google Play Store. From code to store listing.
**Purpose:** Practical guide — "expo push notifications", "eas build", "submit app to app store"
**Timeline:**
- 0:00 — Hook: "Notifications bring users back. The app store brings users in."
- 0:30 — Push notification architecture: Expo push service + Go backend
- 1:15 — Registering push tokens: getting the device token on app start
- 2:00 — Storing tokens: saving to Go API, associating with user
- 2:45 — Sending notifications: Go service calling Expo push API
- 3:30 — Notification handling: foreground, background, tap actions
- 4:15 — Notification types: transactional, marketing, silent data push
- 5:00 — EAS Build: cloud builds for iOS and Android
- 5:30 — Build configuration: eas.json profiles (development, preview, production)
- 6:00 — iOS build: certificates, provisioning profiles, TestFlight
- 6:45 — Android build: signing key, AAB format, internal testing
- 7:15 — App metadata: icons, screenshots, descriptions
- 7:45 — EAS Submit: automated submission to both stores
- 8:15 — Review process: what to expect, common rejection reasons
- 8:45 — Your app is live: from Grit scaffold to app store listing
- 9:15 — Next: series wrap-up

**Thumbnail Prompt:** YouTube thumbnail (1280x720). Dark navy background. Left side: bold white text "NOTIFY + SHIP". Right side: a phone showing a push notification banner at top, with Apple App Store and Google Play Store logos below it. Green "PUBLISHED" indicator. Top-left: orange pill badge "MOBILE 4-5".

---

### Day 30 — Channel Wrap-Up: What's Next

**Title:** "30 Days of Grit — What's Next for the Framework and Channel"
**Description:** We built web apps, desktop apps, and mobile apps — all powered by Grit and Go. This final video recaps the journey, previews the Grit roadmap (Grit Cloud, marketplace, new generators), and invites the community to contribute, request features, and join the Discord.
**Purpose:** Awareness / Community — retention, community building, roadmap hype
**Timeline:**
- 0:00 — Hook: "30 days. 19 courses. 3 platforms. Here's what's next."
- 0:30 — Recap: what we built across 30 videos
- 1:15 — Web platform recap: 9 courses, auth to deployment
- 2:00 — Desktop platform recap: 5 courses, Wails to distribution
- 2:45 — Mobile platform recap: 5 courses, Expo to app store
- 3:30 — Most popular videos: what resonated with the community
- 4:15 — Grit roadmap: what's coming next
- 5:00 — Grit Cloud: hosted deployment platform (preview)
- 5:45 — Marketplace: community templates and plugins
- 6:15 — New generators: more resource types, more field types
- 7:00 — Community: Discord, GitHub discussions, contributing
- 7:30 — How to contribute: issues, PRs, component library
- 8:00 — Feature requests: what do YOU want Grit to build?
- 8:30 — Thank you: to everyone who followed along
- 9:00 — Call to action: subscribe, star the repo, join Discord

**Thumbnail Prompt:** YouTube thumbnail (1280x720). Dark navy background. Center: bold white text "WHAT'S NEXT" with the Grit logo above it, glowing with a purple and sky blue gradient aura. Below the text: three small icons representing web (browser), desktop (monitor), and mobile (phone) in a row. Bottom text: "30 Days Complete". Celebratory but clean — no confetti, just refined glow effects. Top-left: white pill badge "FINALE".

---

## Summary

| Week | Days | Focus | Courses Covered |
|------|------|-------|-----------------|
| 1 | 1-7 | Grit Web — Intro & Setup | Course 0 (2 videos), Course 1 (2 videos), Course 2 (3 videos) |
| 2 | 8-14 | Grit Web — Auth, Admin, Storage | Course 3 (3 videos), Course 4 (2 videos), Course 5 (2 videos) |
| 3 | 15-21 | Grit Web — Jobs, AI, Deploy | Course 6 (2 videos), Course 7 (2 videos), Course 8 (2 videos), Bonus (1 video) |
| 4 | 22-30 | Desktop & Mobile | Desktop 1-5 (4 videos), Mobile 1-5 (4 videos), Wrap-up (1 video) |

**Total:** 30 videos across 30 days covering all 19 courses (9 web + 5 desktop + 5 mobile)
