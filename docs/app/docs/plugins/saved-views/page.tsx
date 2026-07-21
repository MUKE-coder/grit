import Link from 'next/link'
import { ArrowLeft, ArrowRight, Bookmark, Filter, Table, ListFilter } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'
import { CodeBlock } from '@/components/code-block'

export const metadata = {
  title: 'Saved views plugin — Grit',
  description:
    'Let each user capture a resource table’s filters, sort, and search as a named view they can return to in one click.',
  alternates: { canonical: 'https://gritframework.dev/docs/plugins/saved-views' },
}

export default function SavedViewsPage() {
  return (
    <div className="min-h-screen bg-background">
      <SiteHeader />
      <DocsSidebar />

      <main className="lg:pl-64">
        <div className="mx-auto max-w-3xl px-6 py-12">
          <div className="mb-3 flex items-center gap-2">
            <Bookmark className="h-5 w-5 text-primary" />
            <span className="font-mono text-xs uppercase tracking-wider text-muted-foreground">
              Plugin
            </span>
          </div>

          <h1 className="mb-4 font-display text-4xl font-bold tracking-tight">
            Saved views
          </h1>
          <p className="mb-10 text-lg text-muted-foreground">
            Capture a table&apos;s filters, sort, search, and date range as a named view
            &mdash; then switch between them in one click. Per user, per resource.
          </p>

          <div className="prose-grit">
            <CodeBlock language="bash" code={`grit plugin add saved-views
cd apps/api && go run cmd/migrate/main.go`} />

            <h2>The problem it solves</h2>
            <p>
              You keep re-applying the same filters and sort every time you open a table
              &mdash; &quot;active users sorted by signup&quot;, &quot;unpaid invoices this
              month&quot;. It&apos;s a small friction that repeats dozens of times a day.
              Saved views let you capture that arrangement once and reapply it instantly.
            </p>
            <p>
              In some admin ecosystems this exact feature ships as a paid
              &quot;advanced tables&quot; add-on. In Grit it&apos;s a plugin you install
              when you want it.
            </p>

            <div className="my-8 rounded-xl border border-primary/30 bg-primary/5 p-5">
              <div className="mb-2 flex items-center gap-2">
                <ListFilter className="h-4 w-4 text-primary" />
                <strong className="text-foreground">A saved view is just a saved query string</strong>
              </div>
              <p className="mb-0 text-sm text-muted-foreground">
                Grit&apos;s resource tables already sync their entire state &mdash; filters,
                sort, search, date range &mdash; to the URL query string. So a saved view is
                nothing more than that query string, stored per user, per resource. Applying
                one is a navigation; the table restores itself from the URL. Nothing in the
                DataTable itself has to change.
              </p>
            </div>

            <h2>How it works</h2>
            <p>
              Because the table&apos;s state already lives in the URL, the plugin only needs
              to store and replay a query string.
            </p>

            <h3>Backend</h3>
            <p>
              A <code>SavedView</code> model with fields <code>id</code>,{' '}
              <code>user_id</code>, <code>resource</code> (the resource slug),{' '}
              <code>name</code>, <code>query</code> (the raw URL query string), and{' '}
              <code>created_at</code>:
            </p>
            <CodeBlock language="go" code={`type SavedView struct {
    ID        string
    UserID    string
    Resource  string // the resource slug, e.g. "invoices"
    Name      string
    Query     string // "status=published&sort=created_at&order=desc"
    CreatedAt time.Time
}`} />
            <p>A handler exposes three endpoints:</p>
            <CodeBlock language="text" code={`GET    /api/saved-views?resource=<slug>   list the caller's views
POST   /api/saved-views                    create
DELETE /api/saved-views/:id                delete`} />
            <p>
              Every query is scoped to the signed-in user&apos;s id, so a user only ever sees
              or deletes their own views &mdash; a crafted <code>id</code> can&apos;t touch
              someone else&apos;s.
            </p>

            <h3>Frontend</h3>
            <p>
              A <strong>SavedViews</strong> bar renders above every generic resource table,
              injected at a resource-page toolbar marker. It shows a chip per saved view plus
              a <em>Save current view</em> button. Saving captures the current URL query
              (minus transient params like an open create/edit drawer), prompts for a name,
              and POSTs it. Clicking a chip navigates to the resource with that saved query;
              each chip has a small <code>&times;</code> to delete it.
            </p>

            <h2>Install</h2>
            <p>
              The plugin works on the <strong>triple</strong> and <strong>full</strong>{' '}
              architectures &mdash; it needs an admin app to render the views bar. Because it
              adds a table, run migrations after installing:
            </p>
            <CodeBlock language="bash" code={`grit plugin add saved-views
cd apps/api && go run cmd/migrate/main.go`} />

            <h2>Using it</h2>
            <p>
              Open a resource, set some filters, sort, or search, then click{' '}
              <em>Save current view</em> and give it a name. It appears as a chip. Click the
              chip anytime to reapply that exact view; click the <code>&times;</code> on a
              chip to remove it.
            </p>
            <div className="my-8 rounded-xl border border-border/60 bg-muted/20 p-5">
              <div className="mb-2 flex items-center gap-2">
                <Table className="h-4 w-4 text-muted-foreground" />
                <strong className="text-foreground">Where the bar shows up</strong>
              </div>
              <p className="mb-0 text-sm text-muted-foreground">
                The saved-views bar appears on <strong>generated resources</strong>, which
                use the generic resource page. Built-in resources with fully custom pages may
                not show it, since they don&apos;t render the toolbar marker the bar injects
                into.
              </p>
            </div>

            <h2>Use cases</h2>
            <p>
              Any filter-and-sort combination you return to is worth saving:
            </p>
            <ul>
              <li>
                <Filter className="mr-1 inline h-4 w-4 text-primary align-text-bottom" />
                <strong>Active users</strong> &mdash; status filter plus a sort by signup
                date.
              </li>
              <li>
                <strong>Unpaid invoices this month</strong> &mdash; a status filter and a date
                range, captured together.
              </li>
              <li>
                <strong>Published posts by date</strong> &mdash; the segment you review every
                morning, one click away.
              </li>
            </ul>
            <p>
              Views are <strong>per user</strong>, so yours don&apos;t affect other admins,
              and each person builds the set of segments that matches how they work. They&apos;re
              ideal for quickly toggling between different slices of the same table.
            </p>
          </div>

          <div className="mt-12 flex items-center justify-between border-t border-border/40 pt-6">
            <Button asChild variant="ghost">
              <Link href="/docs/plugins/command-palette">
                <ArrowLeft className="mr-2 h-4 w-4" />
                Command palette
              </Link>
            </Button>
            <Button asChild variant="ghost">
              <Link href="/docs/plugins/authoring">
                Writing a plugin
                <ArrowRight className="ml-2 h-4 w-4" />
              </Link>
            </Button>
          </div>
        </div>
      </main>
    </div>
  )
}
