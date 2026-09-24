/**
 * The prompt behind every "Copy prompt to build with AI" button.
 *
 * It is a bootstrap, not an encyclopedia. The depth lives in the Grit skill
 * (skills/grit/SKILL.md, published at /skill.md), which this tells the agent to
 * install first. A prompt that tried to carry all of it would be too long to
 * paste into half the tools people actually use, and would go stale the moment
 * the framework moved, whereas the skill is fetched fresh every time.
 *
 * What this file owns: the order of operations, the curriculum of links, and
 * the rule that nothing is reported as done until it builds.
 */

export const BUILD_WITH_AI_PROMPT = `You are going to build a complete, production-ready application with the **Grit framework** (Go + Gin + GORM on the backend, Next.js or TanStack Router on the frontend, with a generated admin panel). Follow these steps in order.

## Step 1: Load the Grit skill before anything else

Grit is a code generator, and almost every mistake an AI makes with it comes from not knowing that. Load the skill first:

\`\`\`bash
npx skills add MUKE-coder/grit --skill grit
\`\`\`

If you have no skills CLI, fetch and read this instead, and keep it in context for the whole build:

    https://gritframework.dev/skill.md

Also read the machine-readable reference: https://gritframework.dev/docs/ai-skill/llm-guide

## Step 2: Install the CLI

\`\`\`bash
curl -fsSL https://gritframework.dev/install.sh | sh   # macOS, Linux, Git Bash
# Windows PowerShell:  irm https://gritframework.dev/install.ps1 | iex
# With a Go toolchain:  go install github.com/MUKE-coder/grit/v3/cmd/grit@latest
grit version
\`\`\`

Prerequisites: Go 1.21+, Node 20+, pnpm, Docker. https://gritframework.dev/docs/prerequisites

## Step 3: Understand the one thing that matters

**Grit generates code; it is not a runtime library you import.**

\`grit generate resource Post --fields "..."\` writes a Go model, a service, a handler, a Zod schema, TypeScript types, React Query hooks and an admin page, and *injects* the wiring into the router and the registries. That code is then mine to edit.

So the loop is **generate, then edit**. If you hand-write those files instead of running the command, you will produce something that compiles and is unreachable, because the injections are missing. This is the single most common way an agent wastes an afternoon on Grit.

Read these two before writing anything:
- https://gritframework.dev/docs/concepts/code-generation
- https://gritframework.dev/docs/concepts/generated-files

## Step 4: Ask me what I am building, then choose the architecture

Do not guess. Ask me for the product in a sentence, who uses it, and whether it needs a separate admin back office. Then pick:

| If | Command |
|---|---|
| Public site **and** an admin back office | \`grit new app --triple --next\` |
| Just an app, no separate admin | \`grit new app --double --next\` |
| One binary, SPA embedded | \`grit new app --single\` |
| Backend only | \`grit new app --api\` |
| Phone app | \`grit new app --mobile\` |
| Native desktop app | \`grit new-desktop app\` |

\`--next\` for Next.js App Router, \`--vite\` for TanStack Router. https://gritframework.dev/docs/concepts/architecture-modes

Then: \`cd app && docker compose up -d && pnpm install && grit start\`

## Step 5: Model the domain, one command per noun

Before any feature code. Get the field types right the first time; they are what the admin, the API and the validation are all built from.

\`\`\`bash
grit generate resource Customer --fields "name:string,email:email,phone:tel,notes:text:optional"
grit generate resource Invoice  --fields "number:string:unique,total:money,status:select:draft=Draft|paid=Paid,customer:belongs_to:Customer"
\`\`\`

Field type reference (read it, do not guess): https://gritframework.dev/docs/concepts/field-types

Two that are always worth stating: **\`money\` for money, never \`float\`** (https://gritframework.dev/docs/concepts/money), and **\`belongs_to\` is a UUID string**, not an integer.

## Step 6: Write features in the right layer

\`handler → service → model\`. Handlers parse and respond, services hold the logic, models hold the shape and the hooks.

- https://gritframework.dev/docs/backend/handlers
- https://gritframework.dev/docs/backend/services
- https://gritframework.dev/docs/backend/models
- https://gritframework.dev/docs/backend/request-lifecycle
- https://gritframework.dev/docs/backend/response-format (every endpoint answers in this shape)
- https://gritframework.dev/docs/backend/errors (use the \`respond\` package, do not hand-write error bodies)

Frontend: data through the generated React Query hooks, validation through the shared Zod schemas, App Router only, Tailwind and the existing \`components/ui\` primitives, no \`any\`.

- https://gritframework.dev/docs/frontend/hooks
- https://gritframework.dev/docs/frontend/shared-package
- https://gritframework.dev/docs/admin/resources
- https://gritframework.dev/docs/admin/datatable
- https://gritframework.dev/docs/admin/forms
- https://gritframework.dev/docs/admin/relationships
- https://gritframework.dev/docs/admin/widgets

## Step 7: Do not build what is already there

Check this list before writing a feature. Building one of these by hand is the second most common way to waste a day.

| Need | Already shipped |
|---|---|
| Auth, roles, refresh | https://gritframework.dev/docs/backend/authentication |
| 2FA, passkeys, sessions, sign-in links | https://gritframework.dev/docs/backend/account-security |
| Google / GitHub login | https://gritframework.dev/docs/backend/oauth |
| Permissions | https://gritframework.dev/docs/backend/rbac |
| Uploads, images, video | https://gritframework.dev/docs/batteries/storage |
| Email | https://gritframework.dev/docs/batteries/email |
| Background jobs, cron | https://gritframework.dev/docs/batteries/jobs |
| Caching | https://gritframework.dev/docs/batteries/caching |
| LLM calls | https://gritframework.dev/docs/batteries/ai |
| Realtime | https://gritframework.dev/docs/backend/realtime |
| Incoming webhooks | https://gritframework.dev/docs/backend/webhooks |
| Reliable outgoing events | https://gritframework.dev/docs/backend/outbox |
| Audit trail | https://gritframework.dev/docs/backend/append-only |
| Feature flags | https://gritframework.dev/docs/backend/feature-flags |
| Invoices, tax | https://gritframework.dev/docs/backend/invoices |
| Offline clients | https://gritframework.dev/docs/concepts/offline-sync |
| Backups, GDPR | https://gritframework.dev/docs/batteries/backups |
| Multi-tenancy, impersonation, ⌘K, saved views | \`grit plugin add <name>\`, see https://gritframework.dev/docs/plugins/overview |

## Step 8: Verify before you claim anything works

Never tell me a step is done on the strength of having written it. Run all of this, and paste what failed rather than describing it:

\`\`\`bash
cd apps/api && gofmt -l ./internal ./cmd && go build ./... && go vet ./... && go test ./... -count=1
cd ../admin && npx tsc --noEmit && pnpm lint && npx next build
cd ../.. && grit doctor
\`\`\`

When something is wrong at runtime rather than at build time: \`grit routes\` shows what is actually mounted, \`/pulse/ui\` traces requests and queries, \`/studio\` browses the database, \`/docs\` is the generated API reference.

Testing guide: https://gritframework.dev/docs/testing

## Step 9: Ship it

\`\`\`bash
grit deploy --host user@server.com --domain myapp.com
\`\`\`

https://gritframework.dev/docs/deployment · checklist: https://gritframework.dev/docs/deployment/checklist

---

## Worked examples, if you want to see the whole shape first

- Your first app: https://gritframework.dev/blog/build-your-first-grit-app
- A CRM: https://gritframework.dev/blog/build-a-crm-with-grit
- A storefront: https://gritframework.dev/blog/build-a-storefront-with-grit
- An invoice app: https://gritframework.dev/blog/build-an-invoice-app-with-grit
- A mobile app: https://gritframework.dev/blog/build-mobile-app-with-grit
- A desktop app: https://gritframework.dev/blog/build-desktop-app-with-grit
- How the generator thinks: https://gritframework.dev/blog/your-table-our-machinery
- Roles, permissions, backups: https://gritframework.dev/blog/roles-permissions-and-automatic-backups
- Writing a plugin: https://gritframework.dev/blog/grit-plugins-what-they-are-and-how-to-build-one
- GDPR and access reviews: https://gritframework.dev/blog/gdpr-and-access-reviews
- Why Grit exists: https://gritframework.dev/blog/why-i-built-grit

Tutorials: https://gritframework.dev/docs/tutorials · CLI reference: https://gritframework.dev/docs/cli · Full docs: https://gritframework.dev/docs

---

**Now start at Step 1. When you reach Step 4, stop and ask me what I am building.**
`

/** Roughly how long the prompt is, for the button's subtitle. */
export const BUILD_WITH_AI_PROMPT_WORDS = BUILD_WITH_AI_PROMPT.split(/\s+/).length
