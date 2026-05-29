/**
 * Per-kit "starter prompt" generator.
 *
 * Each kit on /docs/tech-kits/* renders a copyable prompt the user pastes
 * into claude.ai. Claude responds with the four planning files
 * (project-description.md, project-phases.md, design-style-guide.md,
 * prompt.md) which the user then feeds to Claude Code via the kit's CLI
 * command.
 *
 * The same builder powers the /docs/ai-integration wizard — it just lets
 * the user pick a kit + a use-case slice and customizes the prompt body.
 */

export type KitSlug =
  | 'single'
  | 'single-vite'
  | 'double'
  | 'triple'
  | 'api'
  | 'mobile'
  | 'desktop'

export interface KitSpec {
  slug: KitSlug
  label: string
  command: string
  /** Bullet list of what the scaffold produces. */
  stack: string[]
  /** Docs URLs that anchor this kit. */
  docsRefs: string[]
  /** Short note steering Claude when the idea is mismatched. */
  fitGuidance: string
}

export const KITS: Record<KitSlug, KitSpec> = {
  single: {
    slug: 'single',
    label: 'Single — Go + embedded SPA',
    command: 'grit new my-app --single',
    stack: [
      'A single Go binary (no separate frontend server)',
      'Next.js (App Router) frontend embedded into the Go binary via `//go:embed`',
      'Postgres in prod, SQLite for local dev (GORM)',
      'Tailwind CSS + shadcn/ui',
      'JWT + OAuth (Google, GitHub) + TOTP 2FA',
      'Asynq background jobs + cron scheduler',
      'S3-compatible file storage (R2, MinIO, S3)',
      'Resend email + HTML templates',
      'AI Gateway (100+ models, one API key, streaming)',
      'In-app **Sentinel** (security WAF) + **Pulse** (observability) admin pages',
      'Audit log with tamper-evident SHA-256 hash chain',
    ],
    docsRefs: [
      'https://gritframework.dev/docs/tech-kits/single',
      'https://gritframework.dev/docs/concepts/architecture-modes/single',
      'https://gritframework.dev/docs/getting-started/quick-start',
    ],
    fitGuidance:
      'If the idea really needs a separate admin panel for non-engineering staff, recommend the `--triple` kit instead.',
  },
  'single-vite': {
    slug: 'single-vite',
    label: 'Single + Vite — TanStack Router SPA',
    command: 'grit new my-app --single --vite',
    stack: [
      'A single Go binary (no separate frontend server)',
      'Vite + React + TanStack Router (file-based routes, type-safe links)',
      'Pure client-side SPA — embedded into the Go binary via `//go:embed`',
      'Postgres in prod, SQLite for local dev (GORM)',
      'Tailwind CSS + shadcn/ui',
      'JWT auth helper (`lib/auth.ts`) with refresh-on-401 baked in',
      'Sub-second cold dev start (no Next.js compile step)',
      'Asynq jobs, S3 storage, Resend mail, AI Gateway — same batteries',
      'In-app **Sentinel** + **Pulse** admin pages',
    ],
    docsRefs: [
      'https://gritframework.dev/docs/tech-kits/single-vite',
      'https://gritframework.dev/docs/concepts/architecture-modes/single',
    ],
    fitGuidance:
      'If the idea needs SEO / SSR (marketing pages, blogs), recommend the `--single` (Next.js) kit instead.',
  },
  double: {
    slug: 'double',
    label: 'Double — Web + API monorepo',
    command: 'grit new my-app --double',
    stack: [
      'A Turborepo monorepo with `apps/web` (Next.js) + `apps/api` (Go)',
      '`packages/shared` for Zod schemas + TypeScript types shared between web and API',
      '`grit sync` keeps TS types in sync with Go structs',
      'CORS pre-wired between web and API',
      'Deploy each app on its own schedule (separate Dockerfiles)',
      'JWT + OAuth + 2FA, Asynq jobs, S3 storage, Resend mail, AI Gateway',
      'In-app **Sentinel** + **Pulse** admin pages on the API',
    ],
    docsRefs: [
      'https://gritframework.dev/docs/tech-kits/double',
      'https://gritframework.dev/docs/concepts/architecture-modes/double',
    ],
    fitGuidance:
      'If the idea has non-engineering admin users (support, ops), recommend the `--triple` kit so the admin panel is a separate app with its own auth scope.',
  },
  triple: {
    slug: 'triple',
    label: 'Triple — Web + Admin + API',
    command: 'grit new my-app --triple',
    stack: [
      'A Turborepo monorepo with `apps/web` + `apps/admin` + `apps/api`',
      '`apps/web` — public Next.js marketing + product surface',
      '`apps/admin` — Filament-style admin panel with `defineResource()`',
      '`apps/api` — Go (Gin + GORM) backend',
      '`packages/shared` for Zod schemas + TypeScript types',
      'RBAC + invitation flow, audit log, in-app Sentinel + Pulse',
      'Asynq jobs, S3 storage, Resend mail, AI Gateway — same batteries',
      'Recommended starter for any SaaS with multi-tenant admin needs',
    ],
    docsRefs: [
      'https://gritframework.dev/docs/tech-kits/triple',
      'https://gritframework.dev/docs/concepts/architecture-modes/triple',
      'https://gritframework.dev/docs/getting-started/quick-start',
    ],
    fitGuidance:
      'If the project is a solo / indie ship with no separate admin staff, the `--single` kit is lighter; flag this in project-description.md.',
  },
  api: {
    slug: 'api',
    label: 'API — Go backend only',
    command: 'grit new my-app --api',
    stack: [
      'Pure Gin + GORM Go API. No frontend.',
      'OpenAPI 3.0 docs auto-served at `/docs` (Scalar UI)',
      'JWT + OAuth + 2FA, RBAC, audit log',
      'Asynq jobs + cron scheduler, S3 storage, Resend mail, AI Gateway',
      'In-app **Sentinel** + **Pulse** admin pages',
      'Smallest scaffold; fastest to deploy',
      'Ideal for: mobile clients, Discord/Slack bots, native desktop, headless services',
    ],
    docsRefs: [
      'https://gritframework.dev/docs/tech-kits/api',
      'https://gritframework.dev/docs/concepts/architecture-modes/api-only',
      'https://gritframework.dev/docs/backend/api-docs',
    ],
    fitGuidance:
      'If the project also needs a web admin (most do), recommend the `--triple` kit instead — the admin panel ships with RBAC and invitations.',
  },
  mobile: {
    slug: 'mobile',
    label: 'Mobile — Expo + API',
    command: 'grit new my-app --mobile',
    stack: [
      'Expo (React Native) frontend in `apps/mobile`',
      'Go API in `apps/api` (Gin + GORM)',
      '`packages/shared` for Zod schemas + TS types shared between mobile and API',
      '`grit sync` keeps types aligned',
      'Mobile-friendly auth: refresh tokens persisted in AsyncStorage, biometric-ready',
      'EAS Build configuration + OTA updates ready',
      'Push notifications scaffolded (Expo + APNs/FCM)',
      'In-app **Sentinel** + **Pulse** on the API',
    ],
    docsRefs: [
      'https://gritframework.dev/docs/tech-kits/mobile',
      'https://gritframework.dev/docs/concepts/architecture-modes/mobile',
    ],
    fitGuidance:
      'If the app also needs a web companion, recommend triple-mode (`grit new my-app --triple --mobile`) so the same API powers web + admin + mobile.',
  },
  desktop: {
    slug: 'desktop',
    label: 'Desktop — Wails + GORM',
    command: 'grit new-desktop my-app',
    stack: [
      'Standalone Wails v2 desktop app (not a monorepo)',
      'Go backend bound to React frontend — call Go methods like JS functions',
      'GORM with SQLite (default) or Postgres',
      'Local auth via bcrypt — no server required',
      'PDF + Excel export built in (go-pdf, excelize)',
      'Frameless window, custom titlebar, draggable panels, dark/light theme',
      'Single binary per platform: macOS / Windows / Linux',
      'Optional `--triple` mode adds outbox-based server sync',
    ],
    docsRefs: [
      'https://gritframework.dev/docs/tech-kits/desktop',
      'https://gritframework.dev/docs/concepts/architecture-modes/multi-client',
    ],
    fitGuidance:
      'If the app must run on the web too, recommend triple-mode (`grit new my-app --triple --desktop`) so the same Go core ships as a server + desktop binary.',
  },
}

/**
 * Build the canonical four-files prompt for a given kit. Inserts an optional
 * use-case sentence the wizard collects (e.g. "I'm building a SaaS for X").
 */
export function buildKitPrompt(slug: KitSlug, opts?: { useCase?: string; idea?: string }): string {
  const kit = KITS[slug]
  const useCaseLine = opts?.useCase
    ? `\n\n**Project shape:** ${opts.useCase}\n`
    : ''
  const ideaBlock =
    opts?.idea && opts.idea.trim().length > 0
      ? opts.idea.trim()
      : '[YOUR IDEA — replace this block with a detailed paragraph describing what your app does, who it\'s for, the core user flows, and any reference apps or inspirations you\'d like to draw from.]'

  return `I want to build the project below using the **Grit framework**.

## What Grit is

Grit is a full-stack meta-framework that scaffolds production-ready apps in one command. It pairs a **Go (Gin + GORM)** backend with React on the front and ships every battery — auth, payments, jobs, mail, file storage, AI, observability, security WAF, and a Filament-style admin panel — out of the box.

- Docs: https://gritframework.dev/docs
- Repo: https://github.com/MUKE-coder/grit

## The kit I want to use — \`${kit.label}\`

I'm using the **${kit.slug}** kit, which produces:

${kit.stack.map((b) => `- ${b}`).join('\n')}

The scaffold command I'll run is:

\`\`\`
${kit.command}
\`\`\`
${useCaseLine}
Reference docs for this kit:

${kit.docsRefs.map((u) => `- ${u}`).join('\n')}

## What I need from you

Generate exactly **four files** I'll save to my project root before scaffolding.

### 1. \`project-description.md\` — the full project briefing

A complete, opinionated brief on the project. Cover:

- One-line positioning (what it is in 12 words)
- Target users (3–5 concrete personas with names and motivations)
- The core problem and why existing alternatives fall short
- Top 5 features ranked by importance, with one sentence each on the user value
- Non-goals (what we are NOT building — anti-scope is as important as scope)
- Success metrics (3 concrete numbers we'd judge success against)
- Inspirations / references (apps to learn from + what to copy / avoid)
- Tech-stack alignment — explain how each Grit battery (auth, jobs, mail, storage, AI, admin) maps to the project's needs

### 2. \`project-phases.md\` — the phased build plan

A sequenced roadmap with checkbox tasks. Use this structure:

\`\`\`
## Phase 1 — Scaffold & Skeleton
- [ ] Run \`${kit.command}\`
- [ ] Configure \`.env\` (database, JWT, SMTP, S3, ...)
- ...

## Phase 2 — Core Data Model
- [ ] Define User model extensions
- [ ] \`grit generate resource Order\`
- ...

## Phase 3 — Authentication & RBAC
## Phase 4 — Core User Flows
## Phase 5 — Admin Panel (or mobile screens / desktop pages — match the kit)
## Phase 6 — Payments / Webhooks (if applicable)
## Phase 7 — Polish & UX
## Phase 8 — Deploy
\`\`\`

Rules for the phases doc:

- Each phase has 5–12 checkbox tasks
- Tasks must be small enough to complete in <30 minutes each
- Reference Grit CLI commands (\`grit generate ...\`, \`grit migrate\`, \`grit sync\`) where appropriate
- Each phase ends with a "Verify" step listing what should work before moving on

### 3. \`design-style-guide.md\` — visual design system

Opinionated and specific. Cover:

- Brand voice (3 adjectives + tagline)
- Color palette — **all with hex codes** — in both dark and light mode: bg primary/secondary/tertiary/elevated, border, text primary/secondary/muted, accent (primary brand color), success, danger, warning, info
- Typography: UI font + code font (specific Google Fonts), type scale (xs through 4xl with px values)
- Spacing scale (4-px base)
- Border radius scale
- Key component patterns: buttons (primary, ghost, outline), cards (elevated, flat), data tables, forms, toasts, modals
- Page-level patterns: navbar, sidebar, dashboard layout, empty states, loading states
- Animation principles (durations, easings)
- Reference apps (e.g., "navbar feels like Linear; tables feel like Vercel dashboard")

### 4. \`prompt.md\` — the Claude Code prompt

This is the file I'll paste into Claude Code. It should:

- Open with a one-paragraph project recap
- Tell Claude Code to **read** \`project-description.md\`, \`project-phases.md\`, and \`design-style-guide.md\` **first** and hold those as load-bearing context
- Reference the Grit docs URLs above
- Tell Claude Code to scaffold with \`${kit.command}\` if the project directory doesn't exist yet
- Tell Claude Code to work through \`project-phases.md\` **one phase at a time**, checking off boxes as it goes, and to **stop and report** at the end of each phase for review before starting the next
- Specify code conventions: snake_case Go files, kebab-case TS files, no \`any\` type, every error handled in Go, business logic lives in services not handlers, all data fetching through React Query / TanStack Query
- Specify the Grit API response shape it must use (\`{ data, meta }\` for lists, \`{ data, message }\` for single, \`{ error: { code, message, details } }\` for errors)
- Tell Claude Code to run \`grit sync\` after every Go model change
- Tell Claude Code to lean on \`design-style-guide.md\` for any UI work — never invent colors or spacing on the fly

## Critical instructions for you (planning Claude)

- Do NOT write Go or React code. You are planning only.
- Be opinionated. If I gave you a vague idea, make concrete choices and explain them in one sentence.
- For every feature you propose, ask yourself "what specific Grit primitive ships this?" and name it.
- ${kit.fitGuidance}
- Output the four files as separate fenced code blocks, each with a \`filename:\` annotation, in order: project-description.md → project-phases.md → design-style-guide.md → prompt.md.

## My idea

${ideaBlock}
`
}
