// Builds the "teach an LLM Grit from zero to 100%" prompt for the AI Integration
// wizard. The content is a distilled, accurate version of the SKILL.md that ships
// in every generated project (internal/scaffold/skill_file.go) plus the framework
// conventions from CLAUDE.md / GRIT.md — kept in sync by hand.

export type ClientId = 'api' | 'website' | 'admin' | 'mobile' | 'desktop'
export type Frontend = 'next' | 'vite'
export type PluginId = 'multitenant' | 'impersonate' | 'command-palette' | 'saved-views'

export interface PromptInput {
  clients: ClientId[]
  frontend: Frontend
  theme: string
  plugins: PluginId[]
  projectName: string
}

const PLUGIN_META: Record<PluginId, { summary: string }> = {
  multitenant: { summary: 'Organizations, per-org roles, and automatic query scoping that fails closed.' },
  impersonate: { summary: 'An admin signs in as another user, with an audit trail and one-click return.' },
  'command-palette': { summary: '⌘K navigation across the admin, built from the resource registry. Frontend-only.' },
  'saved-views': { summary: "Per-user named table views (filters + sort), built on the tables' URL state." },
}

// Map the selected clients + frontend into the recommended `grit new` command
// plus a few honest variants, so the reader sees the shape of the CLI.
export function deriveCommands(input: Pick<PromptInput, 'clients' | 'frontend' | 'theme' | 'projectName'>): {
  primary: string
  variants: { label: string; command: string }[]
  archLabel: string
} {
  const { clients, frontend, theme, projectName } = input
  const name = projectName || 'my-app'
  const has = (c: ClientId) => clients.includes(c)
  const fe = frontend === 'vite' ? ' --vite' : ' --next'
  const themeFlag = theme && theme !== 'atlas' ? ` --theme ${theme}` : ''

  let arch: string
  let archLabel: string
  if (has('website') && has('admin')) {
    arch = 'triple'
    archLabel = 'Triple — web + admin + API'
  } else if (has('admin')) {
    arch = 'triple'
    archLabel = 'Triple — admin comes with the web app'
  } else if (has('website')) {
    arch = 'double'
    archLabel = 'Double — web + API'
  } else if (has('mobile') && !has('website') && !has('admin') && !has('desktop')) {
    arch = 'mobile'
    archLabel = 'Mobile — API + Expo'
  } else {
    arch = 'api'
    archLabel = 'API — Go backend only'
  }

  const webish = arch === 'triple' || arch === 'double'
  let primary = `grit new ${name} --${arch}`
  if (webish) primary += fe
  if (webish) primary += themeFlag
  if (has('mobile') && arch !== 'mobile') primary += ' --expo'
  if (has('desktop')) primary += ' --desktop'

  const variants: { label: string; command: string }[] = []

  if (arch === 'triple') {
    variants.push({ label: 'Add a docs site too (full)', command: primary.replace('--triple', '--full') })
  }
  if (webish && frontend === 'next') {
    variants.push({ label: 'Vite / TanStack Router instead of Next.js', command: primary.replace('--next', '--vite') })
  }
  if (webish && frontend === 'vite') {
    variants.push({ label: 'Next.js instead of Vite', command: primary.replace('--vite', '--next') })
  }
  if (has('website') && !has('admin') && !has('mobile') && !has('desktop')) {
    variants.push({
      label: 'One self-contained binary (embedded SPA)',
      command: `grit new ${name} --single${frontend === 'vite' ? ' --vite' : ''}`,
    })
  }
  if (webish) {
    const alt = theme === 'aurora' ? 'pulse' : 'aurora'
    const tail =
      (has('mobile') && arch !== 'mobile' ? ' --expo' : '') + (has('desktop') ? ' --desktop' : '')
    variants.push({
      label: `Try the ${alt} theme`,
      command: `grit new ${name} --${arch}${fe} --theme ${alt}${tail}`,
    })
  }

  return { primary, variants, archLabel }
}

export function buildFrameworkPrompt(input: PromptInput): string {
  const { clients, frontend, plugins, projectName } = input
  const name = projectName || 'my-app'
  const { primary } = deriveCommands(input)
  const has = (c: ClientId) => clients.includes(c)

  const pluginAdds =
    plugins.length > 0
      ? plugins.map((p) => `cd ${name} && grit plugin add ${p}`).join('\n')
      : ''

  const feName = frontend === 'vite' ? 'Vite + TanStack Router (SPA)' : 'Next.js (App Router)'
  const feRouting =
    frontend === 'vite'
      ? `- Routes live in \`src/routes/\`. File naming: \`index.tsx\` (list), \`new.tsx\` (create), \`$id.tsx\` (detail), \`$id.edit.tsx\` (edit).
- \`Route.useParams()\` for type-safe params. Navigation: \`import { Link, useNavigate } from '@tanstack/react-router'\`.
- Uses \`createHashHistory()\` so the same build runs inside the Wails desktop shell.`
      : `- Routes live in \`app/\`. File naming: \`page.tsx\`, \`layout.tsx\`, \`loading.tsx\`, \`error.tsx\`.
- Server Components by default; add \`'use client'\` for interactivity. Navigation: \`import { useRouter } from 'next/navigation'\`, \`import Link from 'next/link'\`.`

  const clientLines: string[] = []
  if (has('api')) clientLines.push('- **Go API** (Gin + GORM) — the source of truth. Every other client consumes it.')
  if (has('website')) clientLines.push(`- **Web app** — ${feName}, in \`apps/web/\` (or \`frontend/\` for a single-binary build).`)
  if (has('admin')) clientLines.push('- **Admin panel** — a generated, Filament-like dashboard in `apps/admin/` with data tables, forms, widgets, and a roles editor.')
  if (has('mobile')) clientLines.push('- **Mobile app** — Expo / React Native in `apps/expo/`, sharing the same API and types.')
  if (has('desktop')) clientLines.push('- **Desktop app** — a Wails native window in `apps/desktop/` that embeds the React app and can run offline-first against a local SQLite mirror.')

  const pluginSection =
    plugins.length > 0
      ? `\n## Plugins in this project\n\nThis project installs these first-party plugins. A Grit plugin generates real, reversible code into the repo (recorded in \`.grit/plugins.lock.json\`; \`grit plugin remove <name>\` replays it backwards). Treat the generated plugin code as normal app code you can read and edit.\n\n${plugins
          .map((p) => `- **${p}** — ${PLUGIN_META[p].summary}`)
          .join('\n')}\n`
      : `\n## Plugins\n\nNo plugins are installed. If you later need one, plugins generate reversible code into the repo: \`grit plugin list\`, then \`grit plugin add <name>\`. Available: multitenant, impersonate, command-palette, saved-views. Removal is derived from a lockfile, so it's clean.\n`

  const mobileSection = has('mobile')
    ? `\n## Mobile specifics (Expo)\n\n- Token storage uses \`expo-secure-store\` (encrypted) — NEVER AsyncStorage for tokens.\n- Uploads use \`fetch\` + \`FormData\` with NO explicit Content-Type header (let the runtime set the multipart boundary).\n- The API base URL is derived from the Metro host so a physical device can reach your machine — don't hardcode \`localhost\`.\n- Data fetching is React Query, same as web.\n`
    : ''

  const desktopSection = has('desktop')
    ? `\n## Desktop specifics (Wails)\n\n- The desktop app embeds the React frontend in a native window and talks to the same Go API.\n- It's offline-first: a local SQLite mirror syncs to the server when connected. Auth tokens live in the OS keychain, not localStorage.\n- \`grit generate resource\` fans out to desktop CRUD screens too.\n`
    : ''

  return `# You are building with Grit. Read this before writing any code.

Grit is a **full-stack meta-framework**: a **Go** backend (Gin + GORM) plus a
**React** frontend and a generated admin panel, wired together in one monorepo
and driven by a CLI. Think "Laravel/Rails developer experience, but the backend
is Go and the admin and clients are generated for you." Tagline: *Go + React.
Built with Grit.*

**The single most important rule:** Grit generates code — you do not hand-write
CRUD. When the user asks for a new entity (a Post, a Product, an Invoice), you run
\`grit generate resource\` and then customize the generated files. Hand-writing a
model + handler + service + schema + types + hooks + admin page by hand is the
wrong instinct here and fights the framework.

## This project's stack

Scaffold it with:

\`\`\`bash
${primary}
\`\`\`

Clients in this project:

${clientLines.join('\n')}
${pluginAdds ? `\nThen install the plugins:\n\n\`\`\`bash\n${pluginAdds}\n\`\`\`\n` : ''}
Start everything (Docker infra + all apps):

\`\`\`bash
cd ${name}
docker compose up -d          # PostgreSQL, Redis, MinIO, Mailhog
grit migrate                  # create tables + seed default roles
grit dev                      # run every app in the monorepo
\`\`\`

Frontend: **${feName}**.
${feRouting}

## The tech stack (do not substitute)

- **Backend:** Go 1.21+, **Gin** (web), **GORM** (ORM). PostgreSQL in prod, SQLite for tests.
- **Frontend:** ${feName} + **Tailwind CSS** + **shadcn/ui**. Data fetching is **TanStack Query (React Query)** — never raw \`fetch\` in components.
- **Validation:** **Zod**, shared between frontend and generated from the Go types.
- **Monorepo:** Turborepo + pnpm.
- **Infra:** Redis (cache + \`asynq\` jobs), S3-compatible storage (MinIO/R2/S3), Resend (email), Docker.

## Project layout

\`\`\`
${name}/
├── grit.json                 # project manifest (architecture, frontend)
├── docker-compose.yml        # Postgres, Redis, MinIO, Mailhog
├── packages/shared/          # Zod schemas, TS types, constants (shared by all clients)
└── apps/
    ├── api/                  # Go backend
    │   ├── cmd/server/       # entry point
    │   └── internal/         # models, handlers, services, middleware, routes, ...
    ├── web/                  # ${feName}
    ├── admin/                # generated admin panel
    ├── expo/                 # mobile (if present)
    └── desktop/              # Wails desktop (if present)
\`\`\`

(A \`--single\` build has no \`apps/\` — the Go code is at the root and the SPA lives in \`frontend/\`, embedded with \`go:embed\`.)

## Generating a resource — the core workflow

\`\`\`bash
grit generate resource Post --fields "title:string,content:text,published:bool,views:int"
\`\`\`

From ONE command this generates, and wires up:

- **Go:** \`internal/models/post.go\` (GORM model), \`internal/services/post_service.go\`, \`internal/handlers/post_handler.go\`, and route registration.
- **Shared:** a Zod schema and TypeScript types in \`packages/shared\`.
- **Frontend:** React Query hooks (\`use-posts.ts\`) and, for triple/full, an admin resource definition + page.
- **Multi-client:** with mobile/desktop present, CRUD screens for those too.

It injects into existing files at **marker comments** — \`// grit:models\`, \`// grit:handlers\`, \`// grit:routes:protected\`, \`// grit:seeders\`, \`// grit:sync\`, and others. **Never delete these markers** — they're how generation and \`grit remove resource\` work.

### Field types

| Grit type | Go | TypeScript | Form control |
|-----------|-----|-----------|--------------|
| \`string\` | string | string | text |
| \`text\` | string | string | textarea |
| \`int\` / \`float\` | int / float64 | number | number |
| \`bool\` | bool | boolean | toggle |
| \`date\` / \`datetime\` | time.Time | string | date picker |
| \`image\` / \`file\` | string | string | upload |
| \`belongs_to:User\` | FK + relation | id + object | searchable select |

**Modifiers:** append \`:unique\`, \`:required\`, or \`:optional\` after the type, e.g. \`email:string:unique\`.

### Removing a resource

\`grit remove resource Post\` deletes every generated file and reverses every injection — the inverse of generation. Use it instead of deleting files by hand.

## CLI reference

\`\`\`bash
grit new <name> [--single|--double|--triple|--full|--api|--mobile] [--next|--vite] [--desktop] [--expo] [--theme atlas|aurora|pulse]
grit generate resource <Name> --fields "..."   # generate a full-stack resource
grit remove resource <Name>                     # reverse a generation
grit sync                                       # regenerate TS types + Zod from Go models
grit migrate                                    # GORM AutoMigrate + seed default roles
grit dev                                        # run all apps
grit add role <NAME>                            # register a new role across the app
grit plugin add|remove|list|info <name>         # manage plugins
grit studio                                     # open the GORM Studio DB browser
\`\`\`

## Adding a field to an existing resource

There is no "add field" command. Edit the model, then regenerate the derived code:

1. Add the field to the Go model, e.g. \`SKU string \\\`gorm:"size:64" json:"sku"\\\`\`.
2. \`grit sync\` — regenerates the shared TS types + Zod schemas from the Go structs.
3. \`grit migrate\` — adds the DB column.
4. Add the field to the admin resource definition (\`apps/admin/resources/<name>.ts\`) — its \`fields\` (form) and \`table.columns\` — since those are hand-editable, not regenerated.

## Batteries — already included, don't re-implement

Every Grit app ships with, and you should USE rather than rebuild:

- **Auth** — JWT in HttpOnly cookies + bearer for native clients; register/login/refresh/logout.
- **RBAC** — a real roles table + \`user_roles\`, granular \`resource.action\` permissions with wildcards, and a permission editor in the admin. Guard routes with \`middleware.RequireRole("ADMIN")\` or \`middleware.RequireRole("perm:posts.edit")\`. Gate UI with the \`usePermissions()\` hook's \`can(...)\`.
- **File storage** — S3/R2/B2/MinIO with presigned uploads + image processing.
- **Email** — Resend + HTML templates.
- **Background jobs** — \`asynq\` on Redis, with an admin dashboard.
- **Cron** — scheduled tasks.
- **Cache** — Redis cache service + middleware.
- **AI** — Claude / OpenAI with streaming.
- **2FA / TOTP**, **hash-chained audit log**, **feature flags**, **rate limiting / security (Sentinel)**.
- **Automatic database backups** — manual + scheduled dumps to storage, with \`grit restore\` (see the Data & Backup docs).
- **GORM Studio** — a built-in DB browser at \`/studio\`.
${pluginSection}${mobileSection}${desktopSection}
## API response format (follow exactly)

Success (single): \`{ "data": { ... }, "message": "..." }\`
Success (list): \`{ "data": [ ... ], "meta": { "total": 100, "page": 1, "page_size": 20, "pages": 5 } }\`
Error: \`{ "error": { "code": "VALIDATION_ERROR", "message": "...", "details": { ... } } }\`

Status codes: 200, 201, 400, 401, 403, 404, 422 (validation), 500.

## Conventions (match the generated code)

- **Go files** snake_case (\`user_handler.go\`); **structs** PascalCase; handlers stay thin, business logic in services; always handle errors (wrap with \`fmt.Errorf("context: %w", err)\`), never \`_\`-ignore them.
- **TS files** kebab-case; functional components + hooks only; all data fetching through React Query; validate inputs with Zod; no \`any\`.
- **Tables** plural snake_case; **API routes** plural lowercase (\`/api/posts\`); **Zod schemas** PascalCase + \`Schema\`.
- **Styling** is Tailwind utilities + shadcn/ui. Dark, premium aesthetic — think Linear / Vercel dashboard. No custom CSS files.

## Common pitfalls — avoid these

1. Don't hand-write CRUD — run \`grit generate resource\`.
2. Don't delete \`// grit:*\` marker comments.
3. Don't put business logic in handlers — it belongs in services.
4. Don't use raw \`fetch\` in components — use the generated React Query hooks.
5. Don't hardcode config — use \`.env\` + config structs.
6. Don't add dependencies outside the stack above without a reason.
7. Don't use the Pages Router — it's the App Router (Next) or TanStack Router (Vite).

## How to work on this project

1. **Plan first.** Restate the data model as resources and their fields.
2. **Generate.** For each entity, run \`grit generate resource <Name> --fields "..."\`, then \`grit migrate\`.
3. **Customize.** Edit the generated model/service/handler/admin definition for anything the generator can't infer (relationships, custom endpoints, validation).
4. **Wire the UI.** Use the generated hooks; gate on permissions; follow the response format.
5. **Verify.** \`go build ./...\` for the API; build the frontend; run \`grit dev\`.

Build the product. The scaffolding, auth, admin, and batteries are already done — your job is the domain logic on top.
`
}
