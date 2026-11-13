import Link from 'next/link'
import { ArrowRight, ArrowLeft } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { DocsShell } from '@/components/docs/shell'
import { CodeBlock } from '@/components/code-block'
import { Callout } from '@/components/callout'
import { getDocMetadata } from '@/config/docs-metadata'

export const metadata = getDocMetadata('/docs/backend/db-callbacks')

export default function DBCallbacksPage() {
  return (
    <DocsShell>
      <div className="mb-10">
        <span className="tag-mono text-accent/80 mb-3 block">Backend</span>
        <h1 className="text-4xl font-bold tracking-tight mb-4">Database callbacks</h1>
        <p className="text-lg text-muted-foreground leading-relaxed">
          Code that runs around every statement, registered in one place and wired onto every
          connection. Field encryption, HTML sanitising and the append-only guard already work this
          way; <code>internal/dbchain</code> is how yours joins them.
        </p>
      </div>

      <div className="prose-grit">
        <h2>Why a registry rather than db.Callback()</h2>
        <p>
          GORM lets you register a callback on a <code>*gorm.DB</code> whenever you like, and
          nothing stops you doing that. Three things go wrong when you do.
        </p>
        <ul>
          <li>
            <strong>It lands on one connection.</strong> The handle you happened to hold, not the
            one the test suite opens, or the desktop app&apos;s local mirror, or the one a worker
            made. The application then behaves differently from its own tests, and nothing says
            why.
          </li>
          <li>
            <strong>Its order is whatever it is.</strong> A callback registered before{' '}
            <code>crypto:encrypt_map_values</code> sees plaintext; one after it sees ciphertext.
            Both compile.
          </li>
          <li>
            <strong>Nothing can tell you what is on the chain,</strong> which is the first question
            when one of them turns out to be slow.
          </li>
        </ul>

        <h2>Registering a link</h2>
        <CodeBlock
          language="go"
          filename="internal/hooks/region.go"
          code={`func init() {
    dbchain.MustRegister(dbchain.Link{
        Name:   "myapp:stamp_region",
        Op:     dbchain.Create,
        When:   dbchain.After,
        Anchor: "fieldtypes:create",
        Why:    "fills region from the actor's organization",
        Fn:     stampRegion,
    })
}

func stampRegion(tx *gorm.DB) {
    // tx.Statement.Context carries the request, the actor and the
    // organization, the same way it does inside a service.
}`}
        />
        <p>
          <code>Register</code> touches no connection. <code>dbchain.Apply(db)</code> does that, and{' '}
          <code>database.Connect</code> already calls it, so a link registered from an{' '}
          <code>init</code> is on every handle opened afterwards. Call <code>Apply</code> yourself
          anywhere else you open one.
        </p>

        <h2>Ordering</h2>
        <p>
          <code>When</code> and <code>Anchor</code> are the ordering, together: run before or after
          a named callback. One constraint, not a list, because that is all GORM has.
          <code>Anchor</code> left empty means the statement itself, which is what you want unless
          your link has to see a value another one writes.
        </p>
        <Callout type="warning">
          An anchor that names nothing is refused. GORM accepts <code>After(&quot;nope&quot;)</code>{' '}
          without complaint and gives you no ordering at all, which is a callback running in the
          wrong place and saying nothing about it. <code>Apply</code> checks every anchor before it
          registers anything, so a typo leaves the connection as it was and tells you.
        </Callout>
        <p>The anchors you can name: GORM&apos;s own, and this framework&apos;s.</p>
        <CodeBlock
          language="text"
          code={`gorm:create   gorm:query   gorm:update   gorm:delete   gorm:row   gorm:raw
gorm:before_create   gorm:after_create   gorm:preload   ...

crypto:encrypt_map_values   sanitize:html_create     sanitize:html_update
fieldtypes:create           fieldtypes:update
appendonly:update           appendonly:delete
paginate:count_create       paginate:count_update
paginate:count_delete       paginate:count_raw
sync:soft_delete_touches_updated_at`}
        />

        <h2>Recording the queries a run makes</h2>
        <p>
          Set <code>GRIT_QUERY_LOG</code> and the chain writes one JSON line per query: the request
          that issued it, the matched route, the SQL with its placeholders, the row count, and the
          line of your application it came from. No code change, and nothing is recorded when the
          variable is unset.
        </p>
        <CodeBlock
          language="bash"
          code={`GRIT_QUERY_LOG=.grit/queries.jsonl go test ./...
# or run the app and use it
GRIT_QUERY_LOG=.grit/queries.jsonl go run ./cmd/server`}
        />
        <CodeBlock
          language="json"
          filename=".grit/queries.jsonl"
          code={`{"at":"2026-11-13T12:29:03Z","request":"1794572943-6149904","route":"GET /api/v1/posts","op":"query","table":"users","sql":"SELECT * FROM \`users\` WHERE id = ?","rows":1,"caller":"internal/services/post.go:88"}`}
        />
        <p>
          The SQL keeps its placeholders. The shape is what a repeated query has in common, and the
          values are both the part that differs and the part that would put an email address in a
          file somebody pastes into an issue.
        </p>

        <h2>Reading it back</h2>
        <p>
          <code>grit doctor</code> reads the file and reports the same query issued many times
          inside one request, which is what an N+1 looks like from the outside.
        </p>
        <CodeBlock
          language="bash"
          code={`grit doctor
grit doctor --queries ci-run.jsonl   # somewhere else`}
        />
        <CodeBlock
          language="text"
          code={`⚠ GET /api/v1/posts runs the same query 23 times in one request,
  from internal/services/post.go:88 (in 14 requests):
  SELECT * FROM \`users\` WHERE id = ?
    load them in one query instead: Preload the relation, or fetch the ids
    and use a single IN (...).
    (query-repeated-in-one-request)`}
        />
        <p>
          It is silent when there is no log. It also says so when the log is older than the code,
          because a report from before a change reads exactly the same as one from after it, and
          that is the whole danger of an artefact on disk.
        </p>
        <Callout type="tip">
          This answers a different question from{' '}
          <Link href="/docs/backend/pulse">Pulse</Link>, which profiles a request while it is
          happening and has a dashboard. This writes a file that <code>grit doctor</code> reads in
          CI, with nothing running. Use Pulse to see a slow request; use this to fail a build.
        </Callout>

        <h2>The inventory</h2>
        <p>
          <code>dbchain.Links()</code> is the chain, sorted, so printing it twice gives the same
          answer and a diff of it means something.
        </p>

        <Callout type="warning">
          Recording takes a mutex and writes a line per statement. That is fine for a test run or a
          development session and not for a busy server: leave <code>GRIT_QUERY_LOG</code> unset in
          production.
        </Callout>
      </div>

      <div className="mt-12 flex items-center justify-between border-t border-border/60 pt-6">
        <Button variant="ghost" asChild>
          <Link href="/docs/backend/pulse">
            <ArrowLeft className="mr-2 h-4 w-4" />
            Pulse
          </Link>
        </Button>
        <Button variant="ghost" asChild>
          <Link href="/docs/security/doctor">
            grit doctor
            <ArrowRight className="ml-2 h-4 w-4" />
          </Link>
        </Button>
      </div>
    </DocsShell>
  )
}
