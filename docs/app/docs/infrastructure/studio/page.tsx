import Link from 'next/link'
import { ArrowLeft, ArrowRight } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { DocsShell } from '@/components/docs/shell'
import { CodeBlock } from '@/components/code-block'
import { Callout } from '@/components/callout'
import { getDocMetadata } from '@/config/docs-metadata'

export const metadata = getDocMetadata('/docs/infrastructure/studio')

const features: { title: string; href: string; desc: string }[] = [
  { title: 'Schema Discovery', href: '#schema-discovery', desc: 'Reads the database and your GORM structs, and shows you where they disagree' },
  { title: 'Browse & Filter', href: '#browse', desc: 'Paginated grid, column sorting, full-text search' },
  { title: 'Editing Records', href: '#crud', desc: 'Create, edit and delete through forms, with bulk delete' },
  { title: 'Relationships', href: '#relationships', desc: 'Follow a foreign key instead of copying an id into another query' },
  { title: 'SQL Editor', href: '#sql', desc: 'Read/write detection and DDL blocking, and it can be turned off entirely' },
  { title: 'Export', href: '#export', desc: 'Schema as SQL, JSON, YAML, DBML or an ERD; data as JSON, CSV or INSERTs' },
  { title: 'Import', href: '#import', desc: 'Data from JSON, CSV, SQL or Excel; schema from SQL, JSON, YAML or DBML' },
  { title: 'Go Code Generation', href: '#codegen', desc: 'Structs from the schema, and tables from a .go file' },
  { title: 'Access Control', href: '#access', desc: 'Row scoping, hidden and read-only tables, and an audit logger' },
  { title: 'Authentication', href: '#auth', desc: 'A real login page, and the token never touches browser storage' },
]

export default function StudioPage() {
  return (
    <DocsShell>
            <div className="mb-10">
              <span className="tag-mono text-accent/80 mb-3 block">Infrastructure</span>
              <h1 className="text-4xl font-bold tracking-tight mb-4">
                GORM Studio (Database Browser)
              </h1>
              <p className="text-lg text-muted-foreground leading-relaxed">
                A database browser that runs inside your API, reached at <code>/studio</code>.
                Tables, rows, relationships, a query box, import and export. It is mounted in
                every Grit project, because the moments you most want to look at the data are a
                staging container and a customer&apos;s install, and those are exactly where a
                desktop database client is not.
              </p>
            </div>

            <div className="prose-grit">
              <h2 className="text-2xl font-semibold tracking-tight mb-4">What it gives you</h2>
              <div className="grid gap-3 sm:grid-cols-2 mb-10">
                {features.map((f) => (
                  <a
                    key={f.title}
                    href={f.href}
                    className="rounded-lg border border-border/30 bg-card/30 px-4 py-3 block hover:border-accent/40 transition-colors"
                  >
                    <h3 className="text-sm font-semibold text-foreground mb-1">{f.title}</h3>
                    <p className="text-xs text-muted-foreground/60 leading-relaxed">{f.desc}</p>
                  </a>
                ))}
              </div>

              <Callout type="warning" title="It is a database console, reachable over HTTP">
                Nearly all of Studio&apos;s design is constraints, for that reason. Put
                authentication in front of it, keep the SQL editor off unless a task needs it,
                hide the tables that should never be read, and serve it over HTTPS.{' '}
                <Link href="/docs/security/doctor" className="text-accent hover:underline">
                  <code>grit doctor</code>
                </Link>{' '}
                fails when it is mounted with no login or on a default password.
              </Callout>

              <h2 className="text-2xl font-semibold tracking-tight mt-12 mb-4">
                Every feature, in detail
              </h2>

              <h3 id="schema-discovery" className="text-xl font-semibold tracking-tight mb-3">
                Schema Discovery
              </h3>
              <p className="text-muted-foreground leading-relaxed mb-4">
                Studio reads your schema from two places: the database itself, and your GORM
                model structs by reflection. Having both is the useful part, because the gap
                between them is where migration bugs live. A column the struct declares and the
                database does not have is a migration that never ran; a column the database has
                and no struct mentions is one a <code>DropColumn</code> forgot.
              </p>
              <p className="text-muted-foreground leading-relaxed mb-4">
                SQLite, PostgreSQL and MySQL are supported. SQLite goes through a pure-Go driver,
                so there is no CGO and no build toolchain to install.
              </p>

              <h3 id="browse" className="text-xl font-semibold tracking-tight mb-3">
                Browse and Filter
              </h3>
              <p className="text-muted-foreground leading-relaxed mb-4">
                A paginated grid per table, with column sorting and full-text search across the
                row. It is a data grid rather than a report builder: the question it answers well
                is &quot;what does this row actually contain&quot;, which is the question you
                usually have at the point you reach for a database client.
              </p>

              <h3 id="crud" className="text-xl font-semibold tracking-tight mb-3">
                Editing Records
              </h3>
              <p className="text-muted-foreground leading-relaxed mb-4">
                Create, edit and delete through modal forms, and select several rows for a bulk
                delete. Deletion goes through a styled confirmation rather than the browser&apos;s
                own <code>confirm()</code>, which matters more than it sounds: the native dialog
                is easy to dismiss by reflex.
              </p>
              <p className="text-muted-foreground leading-relaxed mb-4">
                Writes here bypass your handlers and therefore your validation, your hooks and
                your audit log. That is the point when you are repairing data, and a hazard when
                you are not, which is what the read-only policy below is for.
              </p>

              <h3 id="relationships" className="text-xl font-semibold tracking-tight mb-3">
                Relationship Navigation
              </h3>
              <p className="text-muted-foreground leading-relaxed mb-4">
                Foreign keys are shown as links you can follow, across <code>has_one</code>,{' '}
                <code>has_many</code>, <code>belongs_to</code> and <code>many_to_many</code>. In
                practice this is what makes the tool fast to use: tracing an order to its
                customer to their addresses is three clicks rather than three queries with ids
                copied between them.
              </p>

              <h3 id="sql" className="text-xl font-semibold tracking-tight mb-3">
                Raw SQL Editor
              </h3>
              <p className="text-muted-foreground leading-relaxed mb-4">
                A query box, with read and write detected automatically and DDL blocked, so a
                stray <code>DROP TABLE</code> is refused rather than run. A table listed as
                hidden is refused here too, not only in the grid.
              </p>
              <p className="text-muted-foreground leading-relaxed mb-4">
                It cannot be row-scoped, because arbitrary SQL cannot be. If you use the scoping
                below to keep one tenant out of another&apos;s data, turn the SQL editor off with{' '}
                <code>DisableSQL: true</code>, or the scoping is decoration.
              </p>

              <h3 id="export" className="text-xl font-semibold tracking-tight mb-3">
                Export
              </h3>
              <p className="text-muted-foreground leading-relaxed mb-4">
                Schema as SQL DDL, JSON, YAML, DBML, or an ERD diagram as PNG or PDF. Data as
                JSON, CSV (zipped, one file per table) or SQL <code>INSERT</code> statements.
              </p>
              <p className="text-muted-foreground leading-relaxed mb-4">
                CSV export escapes formula injection, so a cell beginning <code>=</code> or{' '}
                <code>+</code> arrives in a spreadsheet as text rather than as something that
                runs. A field an attacker controls is an attack on whoever opens the export,
                which is usually not the person who was attacked.
              </p>

              <h3 id="import" className="text-xl font-semibold tracking-tight mb-3">
                Import
              </h3>
              <p className="text-muted-foreground leading-relaxed mb-4">
                Data from JSON, CSV, SQL or Excel. Schema from SQL, JSON, YAML or DBML, which
                creates the tables. Useful for seeding an environment from a sample, and the
                fastest way to ruin one, so it belongs behind the same authentication as
                everything else here.
              </p>

              <h3 id="codegen" className="text-xl font-semibold tracking-tight mb-3">
                Go Code Generation
              </h3>
              <p className="text-muted-foreground leading-relaxed mb-4">
                Two directions. Generate Go structs from an existing schema, with the GORM tags
                that match, which is how you adopt a database you did not create. Or upload a{' '}
                <code>.go</code> file of struct definitions and have the tables created from it.
              </p>
              <p className="text-muted-foreground leading-relaxed mb-4">
                In a Grit project you usually want{' '}
                <Link href="/docs/concepts/code-generation" className="text-accent hover:underline">
                  <code>grit generate resource</code>
                </Link>{' '}
                instead, which writes the model plus the service, handler, routes, schemas, types
                and admin screen. Studio&apos;s generator gives you the struct alone.
              </p>

              <h3 id="access" className="text-xl font-semibold tracking-tight mb-3">
                Access Control
              </h3>
              <p className="text-muted-foreground leading-relaxed mb-4">
                Three controls, and they compose:
              </p>
              <CodeBlock
                language="go"
                code={`studio.Mount(router, db, models, studio.Config{
    DisableSQL: true, // the SQL editor cannot be scoped
    Scope: func(c *gin.Context, table string, tx *gorm.DB) *gorm.DB {
        if table == "orders" {
            return tx.Where("tenant_id = ?", c.GetString("tenant_id"))
        }
        return tx
    },
    TablePolicy: studio.TablePolicy{
        Hidden:   []string{"payment_tokens"}, // never exposed
        ReadOnly: []string{"audit_log"},      // browsable, not editable
    },
    AuditLogger: studio.DefaultAuditLogger, // log every mutation
})`}
              />
              <p className="text-muted-foreground leading-relaxed mt-4 mb-4">
                <strong>Scope</strong> narrows what a request may see, per table, from whatever
                the request carries. <strong>Hidden</strong> removes a table entirely, including
                from the SQL editor. <strong>ReadOnly</strong> lets a table be browsed but not
                changed, which is the right setting for anything whose value is that nobody edited
                it: an audit log, an append-only ledger.{' '}
                <strong>AuditLogger</strong> records every mutation made through Studio, because
                a change made here leaves no trace in your application&apos;s own log.
              </p>

              <h3 id="auth" className="text-xl font-semibold tracking-tight mb-3">
                Authentication
              </h3>
              <p className="text-muted-foreground leading-relaxed mb-4">
                With <code>AuthMiddleware</code> configured, Studio serves a themed login page
                rather than the browser&apos;s native auth popup. The credential token is kept in
                memory only and never written to browser storage, so a refresh asks again and a
                shared machine does not keep your database open for the next person.
              </p>
              <p className="text-muted-foreground leading-relaxed mb-4">
                In a Grit project this is <code>GORM_STUDIO_USERNAME</code> and{' '}
                <code>GORM_STUDIO_PASSWORD</code>. The scaffold adds the middleware only when
                both are set, which means forgetting one leaves Studio open;{' '}
                <Link href="/docs/security/doctor" className="text-accent hover:underline">
                  <code>grit doctor</code>
                </Link>{' '}
                reports exactly that, and reports a default password as an error in production.
              </p>

              <h2 className="text-2xl font-semibold tracking-tight mt-12 mb-4">
                Turning it off
              </h2>
              <p className="text-muted-foreground leading-relaxed mb-4">
                <code>GORM_STUDIO_ENABLED=false</code> does not mount it at all. Worth doing on a
                deployment where nobody is going to use it: the safest console is the one that is
                not there.
              </p>

              <div className="flex items-center justify-between mt-16 pt-8 border-t border-border">
                <Button variant="ghost" asChild>
                  <Link href="/docs/infrastructure/database">
                    <ArrowLeft className="mr-2 h-4 w-4" />
                    Database &amp; Migrations
                  </Link>
                </Button>
                <Button variant="ghost" asChild>
                  <Link href="/docs/backend/pulse">
                    Pulse (Observability)
                    <ArrowRight className="ml-2 h-4 w-4" />
                  </Link>
                </Button>
              </div>
            </div>
    </DocsShell>
  )
}
