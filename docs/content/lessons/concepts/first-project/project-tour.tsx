import { CodeBlock } from '@/components/code-block'
import { TipBox } from '@/components/course/tip-box'
import { Exercise } from '@/components/course/exercise'
import { KnowledgeCheck } from '@/components/course/knowledge-check'

export default function Lesson() {
  return (
    <>
      <p>
        Your <code>my-first-grit/</code> folder has a lot in it — easily 100+
        files. Don&apos;t panic. The shape is predictable and once you learn it
        here, every Grit project (yours, a teammate&apos;s, an open-source one)
        feels familiar. This tour walks every folder + the files inside, top
        to bottom, in the order they&apos;ll start to matter to you.
      </p>

      <h2>Top-level — what&apos;s in the project root</h2>
      <CodeBlock
        language="text"
        filename="my-first-grit/"
        code={`my-first-grit/
├── apps/                    # All runnable programs live here
│   ├── api/                 # Go backend (Gin + GORM)
│   ├── web/                 # Public Next.js site
│   └── admin/               # Next.js admin panel
├── packages/
│   ├── shared/              # TS types + Zod schemas web + admin both import
│   └── grit-ui/             # The 100-component Grit UI registry (local copy)
├── tests/
│   └── k6/                  # Load tests (smoke, load, stress, spike, soak)
├── e2e/                     # Playwright end-to-end browser tests
├── .claude/                 # Claude Code agent skill for this project
├── .github/                 # GitHub Actions CI workflows
├── docker-compose.yml       # Local dev infra (Postgres, Redis, MinIO, Mailhog)
├── docker-compose.prod.yml  # Production stack (behind your reverse proxy)
├── grit.config.ts           # Ports, paths, names — read by the CLI
├── grit.json                # Architecture + frontend + version metadata
├── package.json             # Workspace scripts (dev / build / test)
├── pnpm-workspace.yaml      # Which folders are workspace members
├── turbo.json               # Build cache + dependency graph for Turborepo
├── postcss.config.mjs       # PostCSS config (shared by web + admin)
├── playwright.config.ts     # Playwright runner config
├── README.md                # Project-specific setup notes
├── .env                     # Real secrets (gitignored)
├── .env.example             # Template — commit this, edit .env from it
├── .gitignore               # Files git ignores
├── .dockerignore            # Files docker ignores during build
├── .prettierrc              # Code formatter rules
└── .prettierignore          # Files prettier skips`}
      />
      <p>
        Each root file has a single job — config, secrets, or describes the
        workspace. There&apos;s no &quot;mystery&quot; file at the root: if you
        see something unfamiliar, it&apos;s almost always an editor / tool
        config (Prettier, Docker, PostCSS) and you can ignore it on day one.
      </p>

      <h2>apps/api — the Go backend (this is where most code lives)</h2>
      <CodeBlock
        language="text"
        filename="apps/api/"
        code={`apps/api/
├── cmd/                     # Entry points (one main.go per binary)
│   ├── server/main.go       # The HTTP API process — run this for the API
│   ├── migrate/main.go      # CLI: run AutoMigrate against the configured DB
│   └── seed/main.go         # CLI: seed dev data (admin user, sample blogs)
├── internal/                # All Go packages — see the table below
├── Dockerfile               # Multi-stage build of cmd/server (alpine runtime)
├── go.mod                   # Go module declaration + direct deps
└── go.sum                   # Pinned checksums for every dep`}
      />
      <p>
        Three commands, three entry points. The interesting work is in{' '}
        <code>internal/</code>, which is laid out as 25+ packages — one per
        concern. Here&apos;s the rough map (you don&apos;t need to memorise
        it; come back to this table as you build):
      </p>
      <div className="rounded-lg border border-border/30 bg-card/30 overflow-hidden my-5 text-sm">
        <table className="w-full">
          <thead>
            <tr className="border-b border-border/30 bg-accent/20">
              <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Package</th>
              <th className="text-left px-4 py-2.5 font-medium text-foreground/80">What lives there</th>
            </tr>
          </thead>
          <tbody className="text-muted-foreground">
            <tr className="border-b border-border/20"><td className="px-4 py-2 font-mono text-xs">handlers/</td><td className="px-4 py-2">HTTP handlers — one file per resource (user.go, blog_handler.go, upload.go, ...). Thin: parse request, call service, format response.</td></tr>
            <tr className="border-b border-border/20"><td className="px-4 py-2 font-mono text-xs">services/</td><td className="px-4 py-2">Business logic. Handlers call services; services own the rules and the DB calls.</td></tr>
            <tr className="border-b border-border/20"><td className="px-4 py-2 font-mono text-xs">models/</td><td className="px-4 py-2">GORM structs (User, Upload, Blog, ...). One file per model.</td></tr>
            <tr className="border-b border-border/20"><td className="px-4 py-2 font-mono text-xs">routes/</td><td className="px-4 py-2"><code>routes.go</code> wires every handler to its URL + middleware. The router map.</td></tr>
            <tr className="border-b border-border/20"><td className="px-4 py-2 font-mono text-xs">middleware/</td><td className="px-4 py-2">Auth, CORS, security headers, CSRF, rate-limit, logger, recovery. Each is a small file.</td></tr>
            <tr className="border-b border-border/20"><td className="px-4 py-2 font-mono text-xs">config/</td><td className="px-4 py-2">Loads .env into a typed <code>Config</code> struct. Single source of truth for settings.</td></tr>
            <tr className="border-b border-border/20"><td className="px-4 py-2 font-mono text-xs">database/</td><td className="px-4 py-2">Opens the Postgres connection, runs AutoMigrate, seeds. GORM lives here.</td></tr>
            <tr className="border-b border-border/20"><td className="px-4 py-2 font-mono text-xs">cache/</td><td className="px-4 py-2">Redis cache service + middleware. Set / Get / Delete with TTL.</td></tr>
            <tr className="border-b border-border/20"><td className="px-4 py-2 font-mono text-xs">storage/</td><td className="px-4 py-2">S3-compatible client. Works with MinIO locally, R2 / S3 in prod.</td></tr>
            <tr className="border-b border-border/20"><td className="px-4 py-2 font-mono text-xs">mail/</td><td className="px-4 py-2">Resend client + HTML templates (welcome, password reset, ...).</td></tr>
            <tr className="border-b border-border/20"><td className="px-4 py-2 font-mono text-xs">jobs/</td><td className="px-4 py-2">Background job queue (asynq). Workers + idempotency-aware enqueue helpers.</td></tr>
            <tr className="border-b border-border/20"><td className="px-4 py-2 font-mono text-xs">cron/</td><td className="px-4 py-2">Scheduled tasks (nightly cleanup, weekly digest).</td></tr>
            <tr className="border-b border-border/20"><td className="px-4 py-2 font-mono text-xs">ai/</td><td className="px-4 py-2">Claude + OpenAI unified client. Streaming chat, embeddings.</td></tr>
            <tr className="border-b border-border/20"><td className="px-4 py-2 font-mono text-xs">totp/</td><td className="px-4 py-2">Two-factor auth: TOTP, backup codes, trusted devices.</td></tr>
            <tr className="border-b border-border/20"><td className="px-4 py-2 font-mono text-xs">authz/</td><td className="px-4 py-2">Ownership / IDOR helpers — the &quot;does user X own resource Y?&quot; check.</td></tr>
            <tr className="border-b border-border/20"><td className="px-4 py-2 font-mono text-xs">safefetch/</td><td className="px-4 py-2">SSRF-safe HTTP client for any URL the user supplies.</td></tr>
            <tr className="border-b border-border/20"><td className="px-4 py-2 font-mono text-xs">paginate/</td><td className="px-4 py-2">Generic pagination params + response shape.</td></tr>
            <tr className="border-b border-border/20"><td className="px-4 py-2 font-mono text-xs">respond/</td><td className="px-4 py-2">Consistent JSON envelope helpers: respond.OK, respond.Error.</td></tr>
            <tr className="border-b border-border/20"><td className="px-4 py-2 font-mono text-xs">realtime/</td><td className="px-4 py-2">WebSocket hub for live updates (notifications, list refresh).</td></tr>
            <tr className="border-b border-border/20"><td className="px-4 py-2 font-mono text-xs">webhooks/</td><td className="px-4 py-2">Signed-payload sender for outbound webhooks (Stripe, Twilio).</td></tr>
            <tr className="border-b border-border/20"><td className="px-4 py-2 font-mono text-xs">sync/</td><td className="px-4 py-2"><code>grit sync</code> support — Go types → TypeScript on disk.</td></tr>
            <tr className="border-b border-border/20"><td className="px-4 py-2 font-mono text-xs">flags/</td><td className="px-4 py-2">Feature flags. Toggle features at runtime without a deploy.</td></tr>
            <tr className="border-b border-border/20"><td className="px-4 py-2 font-mono text-xs">audit/</td><td className="px-4 py-2">Activity log writer + the hash-chain integrity check.</td></tr>
            <tr className="border-b border-border/20"><td className="px-4 py-2 font-mono text-xs">export/</td><td className="px-4 py-2">PDF + Excel + CSV generation helpers.</td></tr>
            <tr className="border-b border-border/20"><td className="px-4 py-2 font-mono text-xs">pdf/</td><td className="px-4 py-2">The PDF rendering primitives export/ uses.</td></tr>
            <tr><td className="px-4 py-2 font-mono text-xs">docs/</td><td className="px-4 py-2">Auto-generated OpenAPI spec served at <code>/docs</code> via Scalar.</td></tr>
          </tbody>
        </table>
      </div>

      <TipBox tone="info">
        <strong>The mental model:</strong> when you generate a resource with
        <code> grit generate resource Customer</code>, it touches THREE files
        — <code>internal/models/customer.go</code>,
        <code> internal/services/customer.go</code>,
        <code> internal/handlers/customer.go</code> — plus one line in
        <code> routes/routes.go</code>. Everything else (cache, storage, mail,
        auth) is shared infra you can call from any service. Memorise the
        first 5 packages in the table; the rest you&apos;ll discover as you
        need them.
      </TipBox>

      <h2>apps/web — the public Next.js site</h2>
      <CodeBlock
        language="text"
        filename="apps/web/"
        code={`apps/web/
├── app/                     # App Router routes (page.tsx per route)
├── components/              # Reusable React components
├── hooks/                   # React Query hooks (use-users.ts, use-blogs.ts, ...)
├── lib/                     # api-client.ts, utils.ts, env helpers
├── public/                  # Static assets served as-is (favicon, images)
├── __tests__/               # Component / unit tests (Vitest + RTL)
├── next.config.ts           # Next.js config (App Router, image domains, ...)
├── tailwind.config.ts       # Tailwind theme + tokens (matches admin)
├── postcss.config.js        # PostCSS pipeline (Tailwind + autoprefixer)
├── tsconfig.json            # TS compiler config (path aliases, strict mode)
├── package.json             # Web-specific deps + scripts
├── vitest.config.ts         # Test runner
├── vitest.setup.ts          # Test setup (jest-dom matchers, fake timers)
└── Dockerfile               # Multi-stage build → nginx-served bundle on :3000`}
      />
      <p>
        Standard Next.js 14+ App Router. The two folders to memorise:
        <code> app/</code> for routes (each <code>page.tsx</code> is a URL)
        and <code>hooks/</code> for React Query — every API call goes through
        a hook, not a raw <code>fetch</code> inside a component.
      </p>

      <h2>apps/admin — the Filament-style admin panel</h2>
      <CodeBlock
        language="text"
        filename="apps/admin/"
        code={`apps/admin/
├── app/                     # Routes (auth, dashboard, system pages)
├── resources/               # ★ defineResource() per model — auto-CRUD pages
├── components/              # Admin-specific layout + widgets
├── hooks/                   # React Query hooks (admin endpoints)
├── lib/                     # API client (uses cookie auth + CSRF header)
├── public/                  # Static assets
├── __tests__/               # Vitest tests
├── next.config.ts           # Next.js config
├── tailwind.config.ts       # Tailwind theme — same tokens as apps/web
├── postcss.config.js
├── tsconfig.json
├── package.json             # Admin-specific deps
├── vitest.config.ts
├── vitest.setup.ts
└── Dockerfile`}
      />
      <p>
        Same skeleton as apps/web. The unique folder is <code>resources/</code>
        — each file is a <code>defineResource()</code> call that auto-builds
        list + create + edit + delete pages from a model. We&apos;ll spend a
        full chapter here in the Web (Next.js) course.
      </p>

      <h2>packages/shared — the type bridge</h2>
      <CodeBlock
        language="text"
        filename="packages/shared/"
        code={`packages/shared/
├── types/                   # TypeScript types (generated from Go via grit sync)
├── schemas/                 # Zod schemas (also generated, used for form validation)
├── constants/               # Shared route paths, role names, enums
├── tsconfig.json
└── package.json             # Published to the workspace as @<project>/shared`}
      />
      <p>
        When you run <code>grit sync</code>, Go structs in{' '}
        <code>apps/api/internal/models/*.go</code> become TypeScript types in{' '}
        <code>packages/shared/types/</code>. Both <code>apps/web</code> and
        <code> apps/admin</code> import these — so a wrong field name in your
        React code is a compile error, not a 3am bug.
      </p>

      <h2>packages/grit-ui — the 100-component library</h2>
      <CodeBlock
        language="text"
        filename="packages/grit-ui/"
        code={`packages/grit-ui/
├── registry.json            # Master index — every component's name + path
├── registry/                # Per-component JSON (metadata) + TSX (source)
└── package.json`}
      />
      <p>
        A shadcn-compatible local copy of the Grit UI registry (marketing
        sections, auth screens, SaaS dashboards, e-commerce, layouts). You
        copy components into <code>apps/web/components/</code> or{' '}
        <code>apps/admin/components/</code> as you need them; the originals
        stay here as reference.
      </p>

      <h2>tests/k6 — load tests</h2>
      <CodeBlock
        language="text"
        filename="tests/k6/"
        code={`tests/k6/
├── smoke.js                 # 10 VUs / 30s — sanity check on every PR
├── average-load.js          # Expected production traffic, 5-10 min
├── stress.js                # Ramp until the API breaks
├── spike.js                 # Sudden 10× surge (marketing email scenario)
├── soak.js                  # Moderate load for hours — catches leaks
├── breakpoint.js            # Find the cliff
├── lib/                     # Reused: auth helpers, thresholds
└── README.md                # How to run each one + thresholds explained`}
      />
      <p>
        The K6 Load Testing course walks all of these. They&apos;re ready to
        run today — just point them at your local API.
      </p>

      <h2>e2e — Playwright browser tests</h2>
      <CodeBlock
        language="text"
        filename="e2e/"
        code={`e2e/
├── auth.spec.ts             # Register → login → logout → forgot password
└── admin.spec.ts            # Admin login → DataTable → resource CRUD`}
      />
      <p>
        Headless browser tests of the real frontends against the real API.
        Run with <code>pnpm test:e2e</code>. Adds ~30 seconds to your CI but
        catches whole categories of bugs unit tests miss.
      </p>

      <h2>.claude — Claude Code agent skill</h2>
      <CodeBlock
        language="text"
        filename=".claude/"
        code={`.claude/
└── skills/grit/SKILL.md     # Project-specific instructions for Claude Code`}
      />
      <p>
        If you use Claude Code (Anthropic&apos;s CLI agent), it reads this
        file to learn the conventions of THIS project. Auto-generated; you
        can edit it to teach Claude your team&apos;s preferences.
      </p>

      <h2>.github — CI workflows</h2>
      <CodeBlock
        language="text"
        filename=".github/"
        code={`.github/
└── workflows/               # GitHub Actions YAML — runs on every push/PR
    ├── ci.yml               # Go tests + race + coverage + cross-platform build
    └── release.yml          # Tag-triggered binary release`}
      />

      <h2>The root config files — what each one is for</h2>
      <ul>
        <li>
          <code>grit.config.ts</code> — ports, paths, project name. The CLI
          reads this when you run <code>grit dev</code>, <code>grit deploy</code>, etc.
        </li>
        <li>
          <code>grit.json</code> — architecture mode + frontend + the Grit
          CLI version that scaffolded the project. Used by{' '}
          <code>grit upgrade</code> to know what to regenerate.
        </li>
        <li>
          <code>package.json</code> (root) — workspace-level scripts:{' '}
          <code>dev</code>, <code>build</code>, <code>test</code>,{' '}
          <code>lint</code>. These are aliases that Turbo dispatches into each
          app.
        </li>
        <li>
          <code>pnpm-workspace.yaml</code> — tells pnpm which folders are
          workspace members. Adding a new app means adding it here.
        </li>
        <li>
          <code>turbo.json</code> — Turborepo&apos;s pipeline. Says
          &quot;before <code>build</code>, build the things this depends on,
          and cache the result&quot;.
        </li>
        <li>
          <code>docker-compose.yml</code> — local infrastructure (Postgres,
          Redis, MinIO, Mailhog). Ports are bound to <code>127.0.0.1</code>{' '}
          only — your laptop, nothing on the LAN can reach them.
        </li>
        <li>
          <code>docker-compose.prod.yml</code> — production stack behind a
          reverse proxy. None of the services bind to a public port; traffic
          arrives via Traefik / Caddy / nginx / Dokploy.
        </li>
        <li>
          <code>.env</code> — real secrets, gitignored. Generated with crypto-
          random JWT + session secrets.
        </li>
        <li>
          <code>.env.example</code> — same keys as <code>.env</code> but with
          placeholder values. Commit this; teammates copy it to <code>.env</code>{' '}
          and fill in their own secrets.
        </li>
        <li>
          <code>.gitignore</code>, <code>.dockerignore</code> — what git and
          docker should skip. Standard defaults.
        </li>
        <li>
          <code>.prettierrc</code>, <code>.prettierignore</code> — TS / TSX
          formatting rules so the team is consistent.
        </li>
        <li>
          <code>postcss.config.mjs</code> — root PostCSS config (Tailwind plugin).
          Inherited by both <code>apps/web</code> and <code>apps/admin</code>.
        </li>
        <li>
          <code>playwright.config.ts</code> — test runner config for{' '}
          <code>e2e/</code>. Tells Playwright which browsers to test, where the
          dev server lives, where to put trace files.
        </li>
        <li>
          <code>README.md</code> — your project-specific notes (boot order,
          quirks, deploy targets). Edit this as you go.
        </li>
      </ul>

      <TipBox tone="success">
        <strong>You don&apos;t need to read all of this today.</strong>{' '}
        Bookmark this lesson. Come back when you hit a folder you don&apos;t
        recognise. By the time you ship your second Grit project the layout
        is muscle memory.
      </TipBox>

      <KnowledgeCheck
        question="A teammate asks: 'Where do I put the business logic for handling refunds?' What's your answer?"
        choices={[
          {
            label: 'apps/api/internal/handlers/refund_handler.go',
            feedback:
              "Wrong — handlers are intentionally thin. They just parse input and return responses. Business logic belongs in services.",
          },
          {
            label: 'apps/api/internal/services/refund_service.go',
            correct: true,
            feedback:
              "Right — services own the business logic. The handler calls the service, the service does the work, the handler shapes the response.",
          },
          {
            label: 'apps/api/internal/models/refund.go',
            feedback:
              "Wrong — models are data shape only (GORM struct + tags). No business logic.",
          },
          {
            label: 'apps/web/lib/refund.ts',
            feedback:
              "Wrong — apps/web is the public Next.js frontend. Business logic on the backend stays on the backend.",
          },
        ]}
      />

      <Exercise
        prompt={
          <>
            <p>
              Open your <code>my-first-grit</code> in your editor and answer
              in <code>notes.md</code>:
            </p>
            <ol>
              <li>
                How many files are under <code>apps/api/internal/handlers/</code>?
              </li>
              <li>
                Open <code>apps/api/internal/handlers/user.go</code>. How many
                exported functions does it have? Does it import a service?
              </li>
              <li>
                Open <code>apps/api/internal/routes/routes.go</code> and find
                the line where <code>user</code> routes are mounted. What URL
                prefix is used?
              </li>
              <li>
                Open <code>packages/shared/types/</code>. Pick any{' '}
                <code>.ts</code> file. Does the comment at the top say
                anything about being generated?
              </li>
            </ol>
            <p>
              The point: practice navigating before the next lesson, where we
              boot everything up.
            </p>
          </>
        }
        hint={
          <>
            For (3), the prefix is <code>/api</code> and the user routes are
            mounted as <code>/api/users</code>. For (4), generated files
            start with <code>// Code generated by grit sync. DO NOT EDIT.</code>
          </>
        }
        solution={
          <>
            <p>You should have noticed:</p>
            <ul>
              <li>
                ~20 files under <code>handlers/</code> — one per resource plus
                tests.
              </li>
              <li>
                <code>user.go</code> has 5 exported functions (List, Get,
                Create, Update, Delete) and imports{' '}
                <code>internal/services</code>.
              </li>
              <li>
                Routes are grouped under <code>/api</code> with{' '}
                <code>/api/users</code> on the authed group.
              </li>
              <li>
                Generated files explicitly say so — never hand-edit them; edit
                the Go model and re-run <code>grit sync</code>.
              </li>
            </ul>
          </>
        }
      />

      <h2>What&apos;s next</h2>
      <p>
        You know the layout. Next lesson is a Docker primer — many learners
        get stuck on the &quot;run <code>docker compose up</code>&quot; step
        because Docker itself feels foreign. We&apos;ll fix that before
        actually starting the dev servers in lesson 4.
      </p>
    </>
  )
}
