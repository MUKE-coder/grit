import Link from 'next/link'
import { ArrowLeft, ArrowRight, Layers } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'
import { CodeBlock } from '@/components/code-block'
import { getDocMetadata } from '@/config/docs-metadata'

export const metadata = getDocMetadata('/docs/deployment/multiple-instances')

export default function MultipleInstancesPage() {
  return (
    <div className="min-h-screen bg-background">
      <SiteHeader />
      <DocsSidebar />

      <main className="lg:pl-64">
        <div className="mx-auto max-w-3xl px-6 py-12">
          <div className="mb-3 flex items-center gap-2">
            <Layers className="h-5 w-5 text-primary" />
            <span className="font-mono text-xs uppercase tracking-wider text-muted-foreground">
              Deployment
            </span>
          </div>

          <h1 className="mb-4 font-display text-4xl font-bold tracking-tight">
            Running more than one instance
          </h1>
          <p className="mb-6 text-lg text-muted-foreground">
            A Grit API is one Go binary, and you can run as many copies of it as you like behind a
            load balancer, pointed at the same database and the same Redis. This page is what that
            takes, what is already shared between the copies, and the little that is still counted
            per copy.
          </p>

          <div className="prose-grit">
            <h2 id="requirements">What every copy needs</h2>
            <ul>
              <li>
                <strong>The same database and the same Redis.</strong> Everything below that is shared
                is shared through one of the two.
              </li>
              <li>
                <strong>The same <code>JWT_SECRET</code>.</strong> A token issued by one copy is checked
                by whichever copy the next request lands on.
              </li>
              <li>
                <strong>The same <code>FIELD_ENCRYPTION_KEY</code>,</strong> if you use{' '}
                <code>:encrypted</code> fields, or a value written by one copy cannot be read by another.
              </li>
              <li>
                <strong>Migrations run once per deploy,</strong> not once per copy. The server does not
                migrate on start, so run <code>grit migrate</code> (or the migrate command in your
                image) as a deploy step before the new copies take traffic.
              </li>
            </ul>
            <p>
              No sticky sessions are needed, for HTTP or for WebSockets. Point the load balancer&apos;s
              health check at <code>/api/health</code>.
            </p>

            <h2 id="shared">What is shared between copies</h2>
            <ul>
              <li>
                <strong>Sessions and refresh tokens</strong> live in the database, so signing out of all
                devices signs out everywhere.
              </li>
              <li>
                <strong>The response cache and idempotency keys</strong> live in Redis, so a retried
                request with the same <code>Idempotency-Key</code> replays whichever copy it reaches.
              </li>
              <li>
                <strong>Per-API-key rate limits</strong> are counted in Redis, so a key limited to 100
                requests a minute gets 100, however many copies are running.
              </li>
              <li>
                <strong>Realtime events</strong> cross copies through a Redis backplane. A user whose
                socket is on one copy receives an event caused on another. Measured: one event reached
                1,000 sockets on one copy, caused by a write on another, in 0.16 seconds.
              </li>
              <li>
                <strong>The tamper-evident activity log</strong> is written under a database lock, so
                copies never fork the hash chain. Verified with 50 concurrent writes across two copies.
              </li>
              <li>
                <strong>Background jobs</strong> are queued in Redis, and any copy&apos;s worker takes
                any job.
              </li>
              <li>
                <strong>Permission changes.</strong> Each copy caches who may do what. A role change on
                one copy reaches every other within about a second, through a small{' '}
                <code>cluster_generations</code> table. Before v3.218.0 a revoked permission kept
                working on the other copies until they restarted.
              </li>
              <li>
                <strong>SSO connections.</strong> Each copy builds its identity-provider connections at
                start-up and rebuilds them within a couple of seconds of another copy saving one.
                Before v3.218.0 a connection created on one copy was unknown to the others.
              </li>
              <li>
                <strong>Scheduled jobs run once,</strong> on whichever copy holds a Redis lock. It is a
                30-second lease, renewed every 10 seconds, so if that copy dies another takes over
                within half a minute. Before v3.218.0 every copy ran every scheduled job.
              </li>
            </ul>
            <CodeBlock language="text" code={`grit:cron:leader   the copy running the scheduled jobs (a 30s lease)
cluster_generations   one row per shared cache: authz, sso`} />

            <h2 id="per-copy">What is still counted per copy</h2>
            <ul>
              <li>
                <strong>Sentinel&apos;s per-IP rate limits.</strong> Sentinel keeps those counters in
                memory, so with three copies a client spreading its requests gets three times each
                limit, including the login limit. Tracked in{' '}
                <a
                  href="https://github.com/MUKE-coder/sentinel/issues/18"
                  className="text-primary hover:underline"
                >
                  Sentinel #18
                </a>
                . Until it changes, set the limits for one copy&apos;s share of the traffic, and lean
                on per-API-key limits, which are shared, for anything that must hold exactly.
              </li>
              <li>
                <strong>Pulse</strong> keeps its default store in memory, so each copy&apos;s dashboard
                shows that copy&apos;s traffic.
              </li>
            </ul>

            <h2 id="your-code">In your own code</h2>
            <p>
              Anything you keep in a package-level variable is per copy. If you cache something that
              other copies can change, use the same mechanism the framework does:
            </p>
            <CodeBlock language="go" code={`// Where the data changes:
cluster.Bump(db, "price-list")

// Where it is cached, checked at most once a second:
var prices = cluster.NewWatch(db, "price-list", time.Second)

if prices.Changed() {
    reloadPriceList()
}`} />
          </div>

          <div className="mt-12 flex items-center justify-between border-t border-border/40 pt-6">
            <Button asChild variant="ghost">
              <Link href="/docs/deployment/checklist">
                <ArrowLeft className="mr-2 h-4 w-4" />
                Go-live checklist
              </Link>
            </Button>
            <Button asChild variant="ghost">
              <Link href="/docs/deployment/deploy-command">
                Deploy command
                <ArrowRight className="ml-2 h-4 w-4" />
              </Link>
            </Button>
          </div>
        </div>
      </main>
    </div>
  )
}
