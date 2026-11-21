import Link from 'next/link'
import { DocsShell } from '@/components/docs/shell'
import { getDocMetadata } from '@/config/docs-metadata'
import { GuideText } from '@/components/guide-text'

export const metadata = getDocMetadata('/docs/glossary')

/**
 * Every Grit term in one sentence.
 *
 * The terms live in this file rather than in config/, deliberately. Two
 * indexers read the documentation, the CLI's (`grit docs`) and the website's,
 * and a section that exists only as config is invisible to them. The CLI's
 * reader does parse object literals inside a page, so a const array here is
 * seen; a separate config file is not.
 *
 * Where TERMINOLOGY.md settles a word, these definitions follow it rather than
 * paraphrasing. That file is the arbiter for contributors; this page is the
 * same vocabulary for the people using it.
 *
 * Each entry gets an id so anything can link to the exact term:
 * /docs/glossary#row-rule
 */

interface Term {
  /** The anchor, so /docs/glossary#battery is a stable address. */
  id: string
  term: string
  /** One sentence. If it needs two, the second says what people get wrong. */
  definition: string
  /** Where to read more. */
  href?: string
  /** Words people search for instead. */
  also?: string
}

interface Group {
  id: string
  title: string
  blurb: string
  terms: Term[]
}

const GROUPS: Group[] = [
  {
    id: 'shape',
    title: 'The shape of a project',
    blurb: 'What `grit new` produces, and the words for its parts.',
    terms: [
      {
        id: 'project',
        term: 'Project',
        definition:
          'Everything one `grit new` produces: the root directory and every app inside it.',
        href: '/docs/getting-started/project-structure',
      },
      {
        id: 'app',
        term: 'App',
        definition:
          'One deployable unit inside a project, such as `apps/api`, `apps/web` or `apps/admin`. A project has several; an app is one of them.',
        also: 'service',
      },
      {
        id: 'package',
        term: 'Package',
        definition:
          'A shared library under packages/, such as `packages/shared`, which holds the types and schemas both the API and the frontend use. Not an app, and not a Go module.',
        href: '/docs/frontend/shared-package',
      },
      {
        id: 'architecture',
        term: 'Architecture',
        definition:
          'Which apps a project has: triple, double, single, api or mobile. Chosen at `grit new` and awkward to change afterwards, so it is the first question the CLI asks.',
        href: '/docs/concepts/architecture',
        also: 'tier, mode',
      },
      {
        id: 'triple',
        term: 'Triple',
        definition:
          'The Go API, a customer-facing Next.js web app and a separate Next.js admin, in one monorepo. The default, and the right choice when the operator console should not look like the customer app.',
        href: '/docs/concepts/architecture',
      },
      {
        id: 'double',
        term: 'Double',
        definition:
          'The Go API plus one frontend, with the admin as a route group at /admin inside it.',
        href: '/docs/concepts/architecture',
      },
      {
        id: 'single',
        term: 'Single',
        definition:
          'One frontend that is also the project root, with the Go API beside it and the admin as a section of the router. The frontend is at the root, not in `apps/web`.',
        href: '/docs/concepts/architecture',
      },
      {
        id: 'frontend',
        term: 'Frontend (the choice)',
        definition:
          'Which router and bundler the React apps use: `next` for Next.js, or `tanstack` for TanStack Router on Vite.',
        href: '/docs/frontend/tanstack-router',
        also: 'vite, spa',
      },
      {
        id: 'style',
        term: 'Style',
        definition:
          'Which admin dashboard layout is written: default, modern, minimal or glass. Decided once, at generation.',
        href: '/docs/concepts/styles',
      },
      {
        id: 'theme',
        term: 'Theme',
        definition:
          'Which colour variables the admin uses, such as aurora or emerald. Unlike a style, a theme can be changed at runtime.',
        href: '/docs/design',
      },
    ],
  },
  {
    id: 'making',
    title: 'Making code',
    blurb: 'The five verbs, which are not interchangeable.',
    terms: [
      {
        id: 'scaffold',
        term: 'Scaffold',
        definition:
          'What `grit new` does: write a whole project from nothing. Only ever used for the whole project, never for one resource.',
        href: '/docs/concepts/code-generation',
      },
      {
        id: 'generate',
        term: 'Generate',
        definition:
          'What `grit generate resource` does: add one resource to a project that already exists, by writing new files and injecting lines into existing ones.',
        href: '/docs/concepts/code-generation',
      },
      {
        id: 'upgrade',
        term: 'Upgrade',
        definition:
          'What `grit upgrade` does: bring an existing project up to the CLI version you have installed. Nothing to do with the database.',
        href: '/docs/upgrades',
      },
      {
        id: 'migrate',
        term: 'Migrate',
        definition:
          'Database schema only. `grit migrate` applies migration files; it never touches your application code.',
        href: '/docs/backend/migrations',
      },
      {
        id: 'sync',
        term: 'Sync',
        definition:
          'What `grit sync` does: regenerate the TypeScript types, Zod schemas and error codes from the Go models, so the two sides cannot disagree.',
        href: '/docs/frontend/shared-package',
      },
      {
        id: 'resource',
        term: 'Resource',
        definition:
          'One thing your application stores, such as Invoice or Contact, and everything generated for it: the model, service, handler, routes, types, schema, hooks and admin screens.',
        href: '/docs/concepts/code-generation',
        also: 'entity, model',
      },
      {
        id: 'field-type',
        term: 'Field type',
        definition:
          'What a field is, in the generator’s vocabulary: string, money, image, belongs_to, enum and about thirty others. The field type decides the database column, the Go type, the validation, the form input and the table cell all at once.',
        href: '/docs/concepts/field-types',
      },
      {
        id: 'definition-file',
        term: 'Definition file',
        definition:
          'A YAML file describing a resource, used instead of a long --fields string when a resource has many fields or you want it in version control.',
        href: '/docs/concepts/code-generation',
      },
      {
        id: 'injection',
        term: 'Injection',
        definition:
          'Adding a line to a file your project already has, at a named anchor, so a newly generated resource appears in the route table and the sidebar. It is why `grit remove` can take it out again.',
      },
      {
        id: 'anchor',
        term: 'Anchor',
        definition:
          'A marker comment in a generated file, such as `// grit:models`, that tells the generator where to inject. Deleting one does not break the build; it quietly stops future resources appearing there.',
      },
      {
        id: 'framework-owned',
        term: 'Framework-owned file',
        definition:
          'A file Grit writes and maintains, which travels on every upgrade. Editing one means your change is overwritten the next time you upgrade.',
        href: '/docs/upgrades',
      },
      {
        id: 'manifest',
        term: 'Manifest',
        definition:
          'The record in .grit/ of every file Grit has written and what it looked like. It is how `grit upgrade` can tell a file you edited from one you have not touched, and leave yours alone.',
        href: '/docs/upgrades',
      },
      {
        id: 'repair',
        term: 'Repair',
        definition:
          'A targeted rewrite applied during an upgrade to a file the project may have edited, where replacing the whole file would throw your changes away.',
        href: '/docs/upgrades',
      },
    ],
  },
  {
    id: 'admin',
    title: 'The admin panel',
    blurb: 'The screens an operator uses, and how they are described.',
    terms: [
      {
        id: 'admin-panel',
        term: 'Admin panel',
        definition:
          'The operator-facing screens over your data, in whichever of the three shapes your architecture gives you.',
        href: '/docs/admin/overview',
      },
      {
        id: 'admin-app',
        term: 'Admin app',
        definition:
          'The admin as its own application at `apps/admin`, with its own port and container. Only a triple project has one.',
        href: '/docs/admin/overview',
      },
      {
        id: 'embedded-panel',
        term: 'Embedded panel',
        definition:
          'The admin as a route group at /admin inside the web app, which is what a double project has.',
        href: '/docs/admin/overview',
      },
      {
        id: 'resource-definition',
        term: 'Resource definition',
        definition:
          'The file that describes an admin screen: which columns, which filters, which form fields. Change the definition and the table, the form and the detail view all follow.',
        href: '/docs/admin/resources',
      },
      {
        id: 'widget',
        term: 'Widget',
        definition:
          'A dashboard block: a statistic, a chart or an activity feed. A good one answers a question somebody asks at the start of the day rather than showing a number.',
        href: '/docs/admin/widgets',
      },
      {
        id: 'system-pages',
        term: 'System pages',
        definition:
          'The operational screens that come with the admin: jobs, files, cron, mail preview, trash, sessions, API keys and the route table.',
        href: '/docs/admin/overview',
      },
      {
        id: 'bulk-actions',
        term: 'Bulk actions',
        definition:
          'Creating or editing many rows at once through a grid on a list screen, with each row validated against the model’s own rules so a bad row reports itself rather than failing the batch.',
        href: '/docs/admin/datatable',
      },
      {
        id: 'saved-view',
        term: 'Saved view',
        definition:
          'A named combination of filters, sort and columns on a list screen, so an operator does not rebuild the same query every morning.',
        href: '/docs/plugins/saved-views',
      },
      {
        id: 'trash',
        term: 'Trash',
        definition:
          'The screen listing soft-deleted records, from which they can be restored. It exists because somebody will delete the wrong row.',
        href: '/docs/admin/trash',
      },
    ],
  },
  {
    id: 'data',
    title: 'Data and the backend',
    blurb: 'The layers a request passes through, and the patterns in them.',
    terms: [
      {
        id: 'model',
        term: 'Model',
        definition:
          'A Go struct that maps to a database table, with tags describing the column, the JSON name and the validation rules.',
        href: '/docs/backend/models',
      },
      {
        id: 'handler',
        term: 'Handler',
        definition:
          'The function that turns an HTTP request into a service call and the result back into a response. It holds no business logic.',
        href: '/docs/backend/handlers',
      },
      {
        id: 'service',
        term: 'Service',
        definition:
          'Where the business rules live, and where a transaction begins and ends. Called by handlers, by jobs and by the CLI alike.',
        href: '/docs/backend/handlers',
      },
      {
        id: 'seeder',
        term: 'Seeder',
        definition:
          'Code that puts starting data into a database: an admin account, reference data, or realistic rows for development. Run with `grit seed`.',
        href: '/docs/backend/models',
      },
      {
        id: 'automigrate',
        term: 'AutoMigrate',
        definition:
          'GORM comparing your structs to the database and adding what is missing. Right for development, and not for production, because it never drops or renames and cannot be reviewed before it runs.',
        href: '/guides/gorm-zero-downtime-migrations',
      },
      {
        id: 'soft-delete',
        term: 'Soft delete',
        definition:
          'Marking a row deleted rather than removing it, so it disappears from queries and can be restored from the trash.',
        href: '/docs/admin/trash',
      },
      {
        id: 'sequence',
        term: 'Sequence',
        definition:
          'A gap-free counter for things like invoice numbers, implemented with a locked row rather than a database sequence, because a database sequence does not roll back.',
        href: '/guides/go-money-currency-invoices',
      },
      {
        id: 'outbox',
        term: 'Outbox',
        definition:
          'A table that holds events written in the same transaction as the data they describe, so an event cannot be lost when the process dies between the write and the publish.',
        href: '/docs/backend/outbox',
      },
      {
        id: 'saga',
        term: 'Saga',
        definition:
          'A multi-step operation where each step has an undo, so a failure partway through runs the undos backwards. Used when the steps span services and no single transaction covers them.',
        href: '/docs/backend/sagas',
      },
      {
        id: 'append-only',
        term: 'Append-only table',
        definition:
          'A table the application may insert into and never update or delete, enforced by database grants and a trigger. Used for audit entries.',
        href: '/docs/backend/append-only',
      },
      {
        id: 'n-plus-one',
        term: 'N+1 query',
        definition:
          'One query for a list and then one more for each row in it, usually from reading a relation inside a loop. Invisible locally and expensive in production.',
        href: '/guides/gorm-tutorial-relationships',
      },
    ],
  },
  {
    id: 'security',
    title: 'Access and security',
    blurb: 'Who is calling, and what they may touch.',
    terms: [
      {
        id: 'permission',
        term: 'Permission',
        definition:
          'One thing a user may do, such as `invoice.delete`. Code checks the permission rather than the role, so changing policy does not mean editing handlers.',
        href: '/docs/backend/rbac',
      },
      {
        id: 'role',
        term: 'Role',
        definition:
          'A named bundle of permissions, such as admin, editor or viewer, assigned to a user.',
        href: '/docs/backend/rbac',
      },
      {
        id: 'policy',
        term: 'Policy',
        definition:
          'A function that sees both the user and the record and decides whether the action is allowed. Used for rules a role cannot express, such as "the author may edit their own post for thirty minutes".',
        href: '/docs/security/policies',
      },
      {
        id: 'row-rule',
        term: 'Row rule',
        definition:
          'A scope that is part of the query, so a record belonging to somebody else cannot be returned at all. A permission answers 403 and a row rule answers 404, because a 403 would confirm the record exists.',
        href: '/docs/security/row-rules',
      },
      {
        id: 'owned-by',
        term: 'owned-by',
        definition:
          'The generator flag that puts the owner into every query a resource produces, including the list, the bulk actions and the export, so the rule does not depend on remembering it.',
        href: '/docs/security/row-rules',
      },
      {
        id: 'confidential',
        term: 'Confidential resource',
        definition:
          'A resource an administrator cannot read either. Stricter than ownership, and deliberately not something the role system can grant.',
        href: '/docs/security/authorization',
      },
      {
        id: 'tenant',
        term: 'Tenant',
        definition:
          'One customer’s slice of a shared system, usually an organisation. Every tenant-owned table carries its id, and every query filters on it.',
        href: '/docs/plugins/multitenant',
        also: 'organisation, workspace',
      },
      {
        id: 'session',
        term: 'Session',
        definition:
          'A row backing one refresh token, so signing out of a device means something. It is what makes per-device revoke and an active-sessions screen possible.',
        href: '/docs/backend/authentication',
      },
      {
        id: 'refresh-rotation',
        term: 'Refresh rotation',
        definition:
          'Issuing a new refresh token every time one is used, so a token that arrives twice has been stolen and the whole family is revoked.',
        href: '/docs/backend/authentication',
      },
      {
        id: 'api-key',
        term: 'API key',
        definition:
          'A credential for a machine rather than a person: scoped, stored hashed, rotatable, with the last use recorded.',
        href: '/docs/backend/public-api',
      },
      {
        id: 'audit-log',
        term: 'Audit log',
        definition:
          'A record of who did what, chained so that altering or removing an entry is detectable. Tamper-evident rather than tamper-proof.',
        href: '/docs/security/audit-log',
      },
      {
        id: 'route-table',
        term: 'Route table',
        definition:
          'The static list in `internal/access` of every route and what guards it. It exists because a web framework cannot report its own middleware chain at runtime, so without it nothing can answer whether an endpoint is protected.',
        href: '/docs/security/authorization',
      },
    ],
  },
  {
    id: 'batteries',
    title: 'Batteries and extensions',
    blurb: 'The parts you turn on rather than build.',
    terms: [
      {
        id: 'battery',
        term: 'Battery',
        definition:
          'A capability that comes with a generated project rather than being installed: caching, the job queue, file storage, email, cron, AI. "Batteries included" means these are already wired.',
        href: '/docs/batteries',
      },
      {
        id: 'plugin',
        term: 'Plugin',
        definition:
          'Something you add to a project with `grit plugin add`, which writes real files into your repository rather than adding a dependency that hides behind an interface.',
        href: '/plugins',
      },
      {
        id: 'blueprint',
        term: 'Blueprint',
        definition:
          'A whole application built with Grit that you start from, rather than a feature you add. A separate repository, not part of the scaffold.',
        href: '/blueprints',
      },
      {
        id: 'job',
        term: 'Job',
        definition:
          'Work done outside the request, on a queue, so a slow email provider does not slow your API down. Retried with backoff, and visible in the admin.',
        href: '/docs/batteries/jobs',
        also: 'worker, task, background job',
      },
      {
        id: 'cron',
        term: 'Cron',
        definition:
          'Work on a schedule. It must run on exactly one instance however many you deploy, which is why the scheduler is not simply started inside the API.',
        href: '/docs/batteries/cron',
      },
      {
        id: 'storage-driver',
        term: 'Storage driver',
        definition:
          'Where uploaded files go: local disk, MinIO, S3, Cloudflare R2 or Backblaze B2. One environment variable chooses, and the application code is the same for all of them.',
        href: '/docs/batteries/storage',
      },
      {
        id: 'presigned-url',
        term: 'Presigned URL',
        definition:
          'A time-limited permission to upload one file to one key, signed by your server so the browser can send the bytes straight to storage without them passing through your API.',
        href: '/guides/go-file-upload-s3-presigned',
      },
      {
        id: 'realtime',
        term: 'Realtime',
        definition:
          'Pushing updates to a connected client. Across more than one instance it needs a Redis backplane, because a connection belongs to exactly one process.',
        href: '/docs/backend/realtime',
      },
      {
        id: 'vector-field',
        term: 'Vector field',
        definition:
          'A field storing an embedding in Postgres with pgvector, so a similarity search can find passages that mean the same thing rather than ones that share words.',
        href: '/docs/batteries/ai',
      },
    ],
  },
  {
    id: 'tools',
    title: 'Tools and companions',
    blurb: 'The separate products a generated project uses.',
    terms: [
      {
        id: 'doctor',
        term: 'grit doctor',
        definition:
          'An offline check of a project against the things that are commonly wrong. It reads files and compiled-in floors and never the network, which is what makes it usable in CI.',
        href: '/docs/cli#doctor',
      },
      {
        id: 'studio',
        term: 'GORM Studio',
        definition:
          'A database browser embedded in the generated API at /studio. A separate project of ours, so a bug in it is an issue on its repository.',
        href: '/docs/systems',
      },
      {
        id: 'pulse',
        term: 'Pulse',
        definition:
          'Request observability: what each request did, which queries it ran and where the time went. A separate project of ours.',
        href: '/docs/backend/pulse',
      },
      {
        id: 'sentinel',
        term: 'Sentinel',
        definition:
          'The rate limiter and web application firewall in front of the API. A separate project of ours.',
        href: '/docs/security/rate-limits',
      },
      {
        id: 'scout',
        term: 'Scout',
        definition:
          'A security scanner, fed the generated route table so it knows which endpoints are meant to be protected. Opt-in with `grit add scout`.',
        href: '/docs/security/scout',
      },
      {
        id: 'grit-ui',
        term: 'Grit UI',
        definition:
          'A standalone React component library. Not part of a generated project, and deliberately not added to the scaffold.',
        href: 'https://ui.gritframework.dev',
      },
      {
        id: 'release-channel',
        term: 'Release channel',
        definition:
          'Whether a version is stable or a prerelease. A version is stable unless it carries an -rc, -beta or -alpha suffix, and `grit update` only moves you between stable releases.',
        href: '/docs/stability',
      },
    ],
  },
]

function TermCard({ t }: { t: Term }) {
  const external = t.href?.startsWith('http')
  return (
    <div id={t.id} className="scroll-mt-24 border-b border-border/40 py-5 last:border-0">
      <dt className="mb-1.5 flex flex-wrap items-baseline gap-x-3">
        <span className="font-display text-base font-semibold text-foreground">{t.term}</span>
        {t.also && (
          <span className="text-xs text-muted-foreground/70">also called: {t.also}</span>
        )}
        <a
          href={`#${t.id}`}
          className="ml-auto font-mono text-[11px] text-muted-foreground/40 transition-colors hover:text-accent"
          aria-label={`Link to ${t.term}`}
        >
          #{t.id}
        </a>
      </dt>
      <dd className="text-[15px] leading-relaxed text-muted-foreground">
        <GuideText text={t.definition} />
        {t.href &&
          (external ? (
            <a
              href={t.href}
              target="_blank"
              rel="noreferrer"
              className="ml-2 whitespace-nowrap text-accent hover:underline"
            >
              Read more
            </a>
          ) : (
            <Link href={t.href} className="ml-2 whitespace-nowrap text-accent hover:underline">
              Read more
            </Link>
          ))}
      </dd>
    </div>
  )
}

export default function GlossaryPage() {
  const total = GROUPS.reduce((n, g) => n + g.terms.length, 0)

  return (
    <DocsShell>
      <div className="mb-10">
        <span className="tag-mono text-accent/80 mb-3 block">Reference</span>
        <h1 className="mb-4 text-4xl font-bold tracking-tight">Glossary</h1>
        <p className="text-lg leading-relaxed text-muted-foreground">
          Every word Grit uses that is not ordinary Go or React, in one sentence each. {total}{' '}
          terms. Each one has its own link, so anything can point at the exact definition:{' '}
          <code className="font-mono text-[13px]">/docs/glossary#row-rule</code>.
        </p>
      </div>

      {/* Jump list. Short enough to scan, which is the point of a glossary. */}
      <nav aria-label="Term groups" className="mb-12 rounded-xl border border-border/60 bg-card/30 p-5">
        <ul className="flex flex-wrap gap-x-5 gap-y-2 text-sm">
          {GROUPS.map((g) => (
            <li key={g.id}>
              <a href={`#${g.id}`} className="text-accent hover:underline">
                {g.title}
              </a>
              <span className="ml-1.5 text-xs text-muted-foreground/60">{g.terms.length}</span>
            </li>
          ))}
        </ul>
      </nav>

      {GROUPS.map((g) => (
        <section key={g.id} id={g.id} className="mb-14 scroll-mt-24">
          <h2 className="mb-1.5 font-display text-2xl font-semibold tracking-tight">{g.title}</h2>
          <p className="mb-4 text-sm text-muted-foreground">
            <GuideText text={g.blurb} />
          </p>
          <dl className="rounded-xl border border-border/50 px-5">
            {g.terms.map((t) => (
              <TermCard key={t.id} t={t} />
            ))}
          </dl>
        </section>
      ))}

      <section className="border-t border-border/60 pt-8">
        <h2 className="mb-3 font-display text-xl font-semibold tracking-tight">
          A word that is not here
        </h2>
        <p className="max-w-2xl text-sm leading-relaxed text-muted-foreground">
          Three are deliberately left out, because they are in circulation without an agreed
          meaning and using them without saying which sense you mean causes confusion:{' '}
          <strong>the site</strong> (the marketing pages, or the whole web app),{' '}
          <strong>the dashboard</strong> (the admin&apos;s landing screen, or its whole
          authenticated area) and <strong>tier</strong> (fine with a number attached, misleading
          alone). Say which one you mean.
        </p>
      </section>
    </DocsShell>
  )
}
