import Link from 'next/link'
import { ArrowLeft, ArrowRight, Building2, AlertTriangle, ShieldCheck } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'
import { CodeBlock } from '@/components/code-block'

export const metadata = {
  title: 'Multi-tenancy plugin — Grit',
  description:
    'Organizations, per-organization roles, and automatic query scoping that fails closed.',
  alternates: { canonical: 'https://gritframework.dev/docs/plugins/multitenant' },
}

export default function MultitenantPage() {
  return (
    <div className="min-h-screen bg-background">
      <SiteHeader />
      <DocsSidebar />

      <main className="lg:pl-64">
        <div className="mx-auto max-w-3xl px-6 py-12">
          <div className="mb-3 flex items-center gap-2">
            <Building2 className="h-5 w-5 text-primary" />
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

            <p>Then just query normally — scoping is applied for you:</p>
            <CodeBlock language="go" code={`// Only the active organization's invoices.
db.WithContext(ctx).Find(&invoices)

// OrgID is stamped on insert; you never set it.
db.WithContext(ctx).Create(&Invoice{Amount: 100})`} />

            <div className="my-8 rounded-xl border border-primary/30 bg-primary/5 p-5">
              <div className="mb-2 flex items-center gap-2">
                <ShieldCheck className="h-4 w-4 text-primary" />
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

            <h2>Roles per organization</h2>
            <p>
              Membership carries the role, so someone is Editor in one organization and
              Viewer in another. It reuses the{' '}
              <Link href="/docs/security/authorization">roles</Link> you already have &mdash;
              no parallel permission system.
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
