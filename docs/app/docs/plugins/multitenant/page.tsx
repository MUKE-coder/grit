import Link from 'next/link'
import { ArrowLeft, ArrowRight, Building2, AlertTriangle, ShieldCheck } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'
import { CodeBlock } from '@/components/code-block'
import { PluginSteps, RelatedPlugins } from '@/components/plugin-detail'
import type { PluginStep, RelatedPlugin } from '@/components/plugin-detail'

export const metadata = {
  title: 'Multi-tenancy plugin — Grit',
  description:
    'Organizations, per-organization roles, and automatic query scoping that fails closed.',
  alternates: { canonical: 'https://gritframework.dev/docs/plugins/multitenant' },
}

const STEPS: PluginStep[] = [
  {
    title: 'Install it',
    run: 'grit plugin add multitenant',
    then: 'Writes the Organization and OrganizationMember models, the tenant package holding the GORM callbacks, the organisation screens in the admin, and the switcher in the navbar.',
  },
  {
    title: 'Run the migration',
    run: 'cd apps/api && go run cmd/migrate/main.go',
    then: 'Creates the two tables. Nothing else in the project has changed shape yet.',
  },
  {
    title: 'Make a model tenant-owned',
    run: 'type Deal struct {\n	ID    string `gorm:"primarykey;size:36" json:"id"`\n	Title string `json:"title"`\n\n	tenant.Owned   // <- this line is the whole opt-in\n}',
    language: 'go',
    then: 'Embedding tenant.Owned adds OrgID and brings the model under scoping. Until a model has it, it stays shared across organisations, which is the right default for a currency table and the wrong one for a deal.',
  },
  {
    title: 'Migrate again, and backfill',
    run: 'cd apps/api && go run cmd/migrate/main.go',
    then: 'The new org_id column starts empty on existing rows, and an empty org_id matches no organisation. Rows that predate this are invisible to every scoped query until you give them one, so backfill before anybody looks.',
  },
  {
    title: 'Tell the API which organisation is active',
    run: 'curl -H "X-Organization-ID: org_123" localhost:8080/api/v1/deals',
    then: 'The admin and web clients send this header automatically once an organisation is selected. Without it a scoped query fails closed with NO_ORGANIZATION rather than returning everything.',
  },
  {
    title: 'Opt one query out, visibly',
    run: '// A platform-wide report, deliberately across organisations.\nvar total int64\ntenant.Unscoped(db).Model(&models.Deal{}).Count(&total)',
    language: 'go',
    then: 'Scoping is automatic so that forgetting is impossible; the escape hatch is explicit so that crossing tenants shows up in review. Read internal/tenant/tenant.go before using it.',
  },
]

const RELATED: RelatedPlugin[] = [
  { name: 'Impersonate', href: '/docs/plugins/impersonate', why: 'Once data is per-organisation, the fastest way to answer "what does this member actually see" is to be them for a minute.' },
  { name: 'Kanban boards', href: '/docs/plugins/kanban', why: 'Boards become per-organisation the moment Board, BoardColumn and BoardCard embed tenant.Owned. Nothing else changes.' },
  { name: 'Saved views', href: '/docs/plugins/saved-views', why: 'Views are already per-user; with organisations in play they are the natural way to keep one team\'s filters apart from another\'s.' },
]

export default function MultitenantPage() {
  return (
    <div className="min-h-screen bg-background">
      <SiteHeader />
      <DocsSidebar />

      <main className="lg:pl-64">
        <div className="mx-auto max-w-3xl px-6 py-12">
          <div className="mb-3 flex items-center gap-2">
            <Building2 className="h-5 w-5 text-accent" />
            <span className="font-mono text-xs uppercase tracking-wider text-muted-foreground">
              Plugin
            </span>
          </div>

          <h1 className="mb-4 font-display text-4xl font-bold tracking-tight">
            Multi-tenancy
          </h1>
          <p className="mb-10 text-lg text-muted-foreground">
            One user, many organizations, a different role in each &mdash; with query
            scoping that can&apos;t be forgotten.
          </p>

          <div className="prose-grit">
            <CodeBlock language="bash" code={`grit plugin add multitenant
cd apps/api && go run cmd/migrate/main.go`} />

            <h2>Why a plugin</h2>
            <p>
              Most apps are single-tenant, and putting <code>org_id</code> on every table is
              a schema decision no framework should make for you. Install it when you need
              it.
            </p>

            <h2>Marking a model as tenant-owned</h2>
            <p>
              Embed <code>tenant.Owned</code>. That&apos;s what opts a table into isolation:
            </p>
            <CodeBlock language="go" code={`type Invoice struct {
    ID     string
    Amount int
    tenant.Owned   // adds OrgID + the index
}`} />

            <div className="not-prose my-5 rounded-xl border border-amber-500/30 bg-amber-500/[0.05] p-4 text-sm leading-relaxed">
              <p className="text-foreground">
                <strong className="text-amber-400">Generated resources are shared until you do this.</strong>
              </p>
              <p className="mt-2 text-muted-foreground">
                <code>grit generate resource</code> does not know this plugin is
                installed, so a model it writes has no <code>OrgID</code> and no
                scoping. That is the right default for a shared table like a country
                list, and the wrong one for anything belonging to a customer. Open the
                model and embed <code>tenant.Owned</code> for each resource that
                belongs to an organization, and remember it again after any{' '}
                <code>--force</code> regeneration.
              </p>
              <p className="mt-2 text-muted-foreground">
                Scoping applies per model, so a table without <code>OrgID</code> is
                simply never filtered. Nothing warns you: the endpoint works, returns
                rows, and returns everyone&apos;s.
              </p>
            </div>

            <p>Then just query normally — scoping is applied for you:</p>
            <CodeBlock language="go" code={`// Only the active organization's invoices.
db.WithContext(ctx).Find(&invoices)

// OrgID is stamped on insert; you never set it.
db.WithContext(ctx).Create(&Invoice{Amount: 100})`} />

            <div className="my-8 rounded-xl border border-accent/30 bg-accent/5 p-5">
              <div className="mb-2 flex items-center gap-2">
                <ShieldCheck className="h-4 w-4 text-accent" />
                <strong className="text-foreground">Scoping fails closed</strong>
              </div>
              <p className="mb-0 text-sm text-muted-foreground">
                A query against a tenant-owned model with <em>no</em> active organization
                returns an error &mdash; it does not quietly return every tenant&apos;s
                rows. Hand-written scoping fails the other way: forget one{' '}
                <code>WHERE</code> and you get a silent cross-tenant leak that no test
                catches, because the query returns <em>more</em> rows rather than failing.
              </p>
            </div>

            <h2>Cross-tenant work</h2>
            <p>
              Platform admin screens and background sweeps need to see everything. Opt out
              explicitly &mdash; each call is deliberate and greppable:
            </p>
            <CodeBlock language="go" code={`tenant.Unscoped(db).Find(&allInvoices)`} />

            <h2>How the active organization is chosen</h2>
            <p>
              From the <code>X-Organization-ID</code> header, or automatically when the user
              belongs to exactly one. <strong>No subdomains</strong> &mdash; they force DNS
              and TLS decisions on every deployment and make local development awkward.
            </p>
            <CodeBlock language="bash" code={`curl https://api.example.com/api/invoices \\
  -H "Authorization: Bearer ..." \\
  -H "X-Organization-ID: 01H..."`} />
            <p>
              Membership is <strong>always verified against the database</strong>. Trusting
              the header alone would let any signed-in user read another tenant&apos;s data
              by editing one request header.
            </p>
            <p>
              Requesting an organization you don&apos;t belong to returns <code>404</code>,
              not <code>403</code> &mdash; a 403 confirms it exists and lets an attacker
              enumerate tenants.
            </p>

            <h2>When there is no active organization</h2>
            <p>
              A tenant-owned row cannot be read or written without one, and the scoping
              callback fails closed rather than guessing. The caller is told so:{' '}
              <code>400</code> with <Link href="/docs/backend/errors">the code</Link>{' '}
              <code>NO_ORGANIZATION</code>, which happens when they belong to no
              organization, or to several and named none.
            </p>
            <CodeBlock
              language="json"
              filename="400 Bad Request"
              code={`{
  "error": {
    "code": "NO_ORGANIZATION",
    "message": "no active organization: send it as X-Organization-ID, or join one"
  }
}`}
            />
            <p>
              Until v3.233.0 this arrived as a 500 saying &quot;Failed to fetch deals&quot;,
              because the error reached the generic handler and nothing knew what it meant.
            </p>

            <h2>Roles per organization</h2>
            <p>
              Membership carries a role, so somebody can be an editor in one organization and
              a viewer in another. It reuses the{' '}
              <Link href="/docs/security/authorization">roles</Link> you already have, with no
              parallel permission system.
            </p>
            <CodeBlock
              language="bash"
              code={`# Give Dave the "Deal Closer" role, in this organization only
curl -X POST https://api.example.com/api/v1/organizations/01H.../members \\
  -H "Authorization: Bearer ..." \\
  -d '{"user_id": "01H...", "role_id": "01H..."}'`}
            />
            <p>
              The semantics, because they are a choice: an organization role <strong>adds</strong>{' '}
              permissions inside that organization. Platform roles, the ones assigned with{' '}
              <code>PUT /users/:id/roles</code>, stay global. So an organization can grant but
              cannot take away, and a project with no membership roles behaves exactly as it did.
            </p>
            <p>
              This is what v3.233.0 fixed: the membership role was stored, put on the request
              context, and read by nothing, so being an administrator of one organization meant
              being one in every organization you belonged to.
            </p>

            <h2>Single sign-on does not join an organization</h2>
            <p>
              A user provisioned by SSO or a social login belongs to no organization, so every
              tenant-owned endpoint answers <code>NO_ORGANIZATION</code> for them until somebody
              puts them in one. That is deliberate: which organization a new user joins is a
              policy the framework cannot guess, and guessing would be a cross-tenant bug rather
              than a convenience. Map it where you know it, from the email domain, the SSO
              connection, or a group the identity provider releases.{' '}
              <code>grit doctor</code> reminds you when both are present.
            </p>

            <h2>Endpoints</h2>
            <CodeBlock language="text" code={`GET    /api/organizations                       orgs you belong to
POST   /api/organizations                       create (you become owner)
GET    /api/organizations/:id/members
POST   /api/organizations/:id/members
DELETE /api/organizations/:id/members/:userId`} />

            <div className="my-8 rounded-xl border border-amber-500/30 bg-amber-500/5 p-5">
              <div className="mb-2 flex items-center gap-2">
                <AlertTriangle className="h-4 w-4 text-amber-500" />
                <strong className="text-foreground">Adding this to an existing app</strong>
              </div>
              <p className="mb-0 text-sm text-muted-foreground">
                Existing rows have no <code>org_id</code>. Before adding{' '}
                <code>tenant.Owned</code> to a model that already has data, create a default
                organization and backfill it &mdash; otherwise those rows become invisible
                to every scoped query.
              </p>
            </div>
          </div>

          <PluginSteps steps={STEPS} intro={'Taking a single-tenant project to one where a user belongs to several organisations and holds a different role in each. Nothing is scoped until you say which models are owned, which is step 3.'} />

          <RelatedPlugins items={RELATED} />

          <div className="mt-12 flex items-center justify-between border-t border-border/40 pt-6">
            <Button asChild variant="ghost">
              <Link href="/docs/plugins/authoring">
                <ArrowLeft className="mr-2 h-4 w-4" />
                Plugins
              </Link>
            </Button>
            <Button asChild variant="ghost">
              <Link href="/docs/security/authorization">
                Roles &amp; Permissions
                <ArrowRight className="ml-2 h-4 w-4" />
              </Link>
            </Button>
          </div>
        </div>
      </main>
    </div>
  )
}
