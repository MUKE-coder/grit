import Link from 'next/link'
import { ArrowLeft, ArrowRight } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'
import { CodeBlock } from '@/components/code-block'
import { Callout } from '@/components/callout'
import { getDocMetadata } from '@/config/docs-metadata'

export const metadata = getDocMetadata('/docs/scaling')

const stages: { n: string; name: string; grit: string; you: string; state: 'free' | 'flag' | 'yours' }[] = [
  { n: '0', name: 'Measure first', grit: 'Request percentiles, connection use, slowest queries, cache hit rate, and a verdict.', you: 'Run grit scale against production.', state: 'free' },
  { n: '1', name: 'One server, one database', grit: 'Config from env, one database module, uploads in object storage, /health. A Grit app is born finished here.', you: 'Nothing.', state: 'free' },
  { n: '2', name: 'Vertical scaling', grit: 'Go uses every core in one process. There is no cluster module to add.', you: 'Buy a bigger machine.', state: 'free' },
  { n: '3', name: 'Horizontal + load balancer', grit: 'Stateless by construction, graceful shutdown on SIGTERM, a health endpoint the balancer can poll.', you: 'Set an instance count. The balancer is your platform’s.', state: 'free' },
  { n: '4', name: 'Stateless servers', grit: 'Sessions in the database, uploads in object storage, cache invalidation across replicas, and cron elected to one instance.', you: 'Nothing. The three-emails bug cannot happen.', state: 'free' },
  { n: '5', name: 'Connection pooling', grit: 'pgbouncer already in docker-compose.prod.yml, per-instance pool limits, and a doctor check that does the arithmetic.', you: 'Point DATABASE_URL at the pooler.', state: 'flag' },
  { n: '6', name: 'Indexes, then read replicas', grit: 'Indexes on foreign keys as generated. Replica routing, read-your-own-writes, and a lag probe.', you: 'Set DATABASE_REPLICA_URLS. No code changes.', state: 'flag' },
  { n: '7', name: 'Caching', grit: 'cache.Remember: cache-aside with fallback, jittered TTLs, stampede protection and a hit rate.', you: 'Decide what may be stale. Nobody can decide that for you.', state: 'flag' },
  { n: '8', name: 'Queues and background jobs', grit: 'asynq workers, retries with backoff, a jobs dashboard, cron, and a transactional outbox.', you: 'Move slow work into a job.', state: 'free' },
  { n: '9', name: 'Sharding', grit: 'Nothing, deliberately.', you: 'Archive, partition, or move to Citus first. Almost nobody needs this.', state: 'yours' },
]

export default function ScalingPage() {
  return (
    <div className="min-h-screen bg-background isolate">
      <SiteHeader />
      <DocsSidebar />

      <main className="lg:pl-64">
        <div className="container max-w-screen-xl py-10 px-6">
          <div className="max-w-3xl">
            <div className="mb-10">
              <span className="tag-mono text-primary/80 mb-3 block">Operations</span>
              <h1 className="text-4xl font-bold tracking-tight mb-4">Scaling</h1>
              <p className="text-lg text-muted-foreground leading-relaxed">
                Ten stages, from one server to sharding. A Grit app starts at the far side of five
                of them, and the rule for the rest is the same as everywhere else: do not add a
                box until something is actually breaking, then add exactly one.
              </p>
            </div>

            <div className="prose-grit">
              <h2 id="where">Where does a Grit app start?</h2>
              <p>
                Most scaling guides begin with work you have to do first: move config to
                environment variables, put database access behind one module, get uploads off the
                local disk, add a health endpoint. That is the whole of Stage 1, and a generated
                Grit project has all of it on the first commit.
              </p>
              <p>
                It goes further than that. Sessions are rows, not process memory. Cron is elected
                to a single instance, so the daily digest is sent once rather than once per
                replica. The server drains in-flight requests on <code>SIGTERM</code>. Jobs retry
                with backoff and there is a transactional outbox for the ones that must not be
                lost. Those are Stages 4 and 8, and they are there before you have a user.
              </p>

              <div className="not-prose my-8 overflow-hidden rounded-xl border border-border">
                <table className="w-full text-sm">
                  <thead className="bg-muted/40">
                    <tr>
                      <th className="px-3 py-2.5 text-left font-medium">Stage</th>
                      <th className="px-3 py-2.5 text-left font-medium">What Grit does</th>
                      <th className="px-3 py-2.5 text-left font-medium">What you do</th>
                    </tr>
                  </thead>
                  <tbody>
                    {stages.map((s) => (
                      <tr key={s.n} className="border-t border-border align-top">
                        <td className="px-3 py-2.5 whitespace-nowrap">
                          <span className="font-mono text-xs text-muted-foreground">{s.n}</span>{' '}
                          <span className="font-medium">{s.name}</span>
                        </td>
                        <td className="px-3 py-2.5 text-muted-foreground">{s.grit}</td>
                        <td className="px-3 py-2.5 text-muted-foreground">{s.you}</td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>

              <h2 id="scale">Which stage are you at?</h2>
              <p>
                Point it at the deployment that is struggling. These numbers describe wherever they
                are measured, and a laptop under no load is always healthy. The admin panel shows
                the same report; see <a href="#readiness">below</a>.
              </p>
              <CodeBlock
                terminal
                code={`grit scale --api https://api.yourapp.com --token $GRIT_ADMIN_TOKEN`}
              />
              <CodeBlock
                language="text"
                code={`  Requests
    p50      3ms   p95      7ms   p99     13ms   max    164ms
    65 requests in the window, 5.0/s since start

  Database
    11 of 100 connections in use (11%), pool max 25 per instance
    No read replicas

  Healthy
    p99 13ms over 65 requests, 11 of 100 database connections in use.

  Do this next: Nothing. Adding infrastructure now buys complexity and no speed.`}
              />
              <p>
                It names one thing, and most of the time that thing is nothing. A tool that lists
                six possible improvements has handed the hardest part of the job, choosing, back
                to you.
              </p>

              <Callout type="note" title="Why the API measures and the CLI only renders">
                Percentiles, connection counts and query statistics are only true where the load
                is. Measuring in the CLI would describe your laptop. The API exposes{' '}
                <code>GET /api/v1/scale</code> (admin only, because it reports connection counts
                and query shapes) and <code>grit scale</code> reads it over the network.
              </Callout>

              <h2 id="readiness">The same thing, in the admin panel</h2>
              <p>
                <code>grit scale</code> answers it from a terminal. The admin&apos;s{' '}
                <strong>Observability</strong> page answers it for everyone else, at the top of the
                page, refreshed every minute: the verdict, and then all ten stages with the state of
                each on this deployment.
              </p>
              <CodeBlock
                language="text"
                code={`Scaling readiness                        9 of 10 stages handled

  Healthy
  p99 13ms over 65 requests, 11 of 100 database connections in use.
  Do this next: Nothing. Adding infrastructure now buys complexity and no speed.

  0  Measure first               p50 3ms, p95 7ms, p99 13ms over 65 requests
  1  One server, one database    config from the environment, one database
                                 module (postgres), uploads in object storage
  2  Vertical scaling            8 cores visible, GOMAXPROCS 8
  3  Horizontal + load balancer  4 instances declared
  4  Stateless servers           sessions are rows, uploads on s3, cron
                                 elected to one instance through Redis
  5  Connection pooling          11 of 100 connections in use, pool max 25
  6  Indexes, then replicas      no replicas, and none needed until reads peg
                                 the primary with the indexes already right
  7  Caching                     hit rate 94% over 20,431 lookups
  8  Queues and background jobs  asynq workers, retries with backoff, cron
                                 and a transactional outbox
  9  Sharding                    largest table orders, about 412,000 rows`}
              />
              <p>
                Every line is read off the running deployment rather than off a list of features. A
                tick against Stage 4 means this app was observed keeping sessions in the database
                and uploads in object storage; it is not a claim about what the framework can do.
                The three measured stages, 5, 6 and 7, go amber on the same thresholds the verdict
                uses, so the panel cannot show green for the stage the verdict is calling out.
              </p>
              <p>
                Two things it will tell you about that are easy to miss until a second instance
                exists. <code>STORAGE_DRIVER=local</code> keeps uploads on one machine&apos;s disk,
                so the second instance serves 404s for half of them. And SQLite serialises writes,
                which makes every scaling question after Stage 1 have the same answer.
              </p>
              <Callout type="note" title="Why a hollow tick is still a good answer">
                Stages 6 and 7 usually read &quot;ready, not needed yet&quot;. That is the correct
                state for almost every application, and the panel says so rather than leaving a gap
                that looks like something missing. Replica routing and{' '}
                <code>cache.Remember</code> are in the project either way; they are one environment
                variable and one function call from being in force.
              </Callout>

              <h2 id="replicas">Stage 6: read replicas</h2>
              <p>One environment variable, and no code changes:</p>
              <CodeBlock
                language="bash"
                code={`DATABASE_REPLICA_URLS=postgres://...replica-1,postgres://...replica-2`}
              />
              <p>
                Reads go to the replicas from the next boot. Writes, and everything inside a
                transaction including its reads, stay on the primary. That last rule is the one
                that matters and the one hand-rolled splits get wrong: a balance check inside the
                transaction that debits the balance must not read a replica, and here it cannot,
                whatever the handler was written to do.
              </p>
              <p>For the reads whose answer decides a write, outside a transaction:</p>
              <CodeBlock
                language="go"
                code={`// Never stale: the answer decides whether to sell the seat.
database.Primary(h.DB).First(&seat, "id = ?", id)

// Staleness is the point: a report nobody acts on in the next second.
database.Replica(h.DB).Find(&monthlyTotals)`}
              />

              <h3 id="ryow">Read your own writes</h3>
              <p>
                The classic replica bug: you post, the feed loads from a replica that has not
                caught up, your own post is missing, and you post it again. Grit pins a person to
                the primary for five seconds after they write, using a cookie rather than a shared
                store: it travels with the person who wrote, costs no lookup, and cannot itself be
                stale. Generated list and detail handlers use it. Hand-written ones should too:
              </p>
              <CodeBlock
                language="go"
                code={`database.ForRequest(c, h.DB).Find(&posts)`}
              />

              <h2 id="cache">Stage 7: caching</h2>
              <CodeBlock
                language="go"
                code={`var followers int64
err := cache.Remember(ctx, svc.Cache, "user:"+id+":followers", time.Minute, &followers,
    func() (int64, error) {
        var n int64
        err := h.DB.Model(&models.Follow{}).Where("followee_id = ?", id).Count(&n).Error
        return n, err
    })`}
              />
              <p>
                <code>Remember</code> is cache-aside with the three things a hand-written version
                misses. A cache read failure falls through to the loader, so Redis being down makes
                the app slower rather than broken. TTLs are jittered, so ten thousand keys written
                by one deploy do not all expire in the same second. And when a hot key expires, one
                caller rebuilds it while the rest wait briefly and read the result, instead of every
                concurrent request running the same expensive query.
              </p>
              <p>
                Invalidate from the write that makes the value wrong, and keep the TTL as the
                safety net for the invalidation somebody forgets:
              </p>
              <CodeBlock language="go" code={`cache.Forget(ctx, svc.Cache, "user:"+id+":followers")`} />

              <Callout type="warning" title="The one thing no framework can do for you">
                Deciding what is allowed to be wrong, and for how long. A follower count may be a
                minute stale. An account balance at the moment of a debit may not. Caches and
                replicas make an app faster and slightly wrong on purpose, and which parts are
                allowed to be wrong is a product decision. Write the list down.
              </Callout>

              <h2 id="pool">Stage 5: the arithmetic</h2>
              <p>
                Connection exhaustion is the only failure here that arrives as errors rather than
                slowness: the API returns 500s while every dashboard looks healthy. It is also pure
                arithmetic, so <code>grit doctor</code> does it:
              </p>
              <CodeBlock
                language="text"
                code={`⚠ database  8 instances x DB_MAX_OPEN_CONNS=25 is 200 connections, and
            max_connections is 100. Spikes will fail with "sorry, too many
            clients already" while every dashboard looks healthy
            put pgbouncer in front, or set DB_MAX_OPEN_CONNS=10`}
              />
              <p>
                Tell it how many instances you run with <code>APP_INSTANCES</code>. That is the
                number people get wrong: a pool of 25 is comfortable on one machine and fatal on
                eight, and nothing else in the config mentions the other seven.
              </p>

              <h2 id="sharding">Stage 9: sharding</h2>
              <p>
                Grit does nothing here, deliberately. Sharding is the most expensive tool in the
                box and most applications never need it. Before it: archive cold rows, use
                Postgres native partitioning, move analytics to a warehouse and search to a search
                engine, or let Citus or CockroachDB shard for you. If you genuinely outgrow one
                primary, the shard key is the decision that matters, and it should be the column
                your queries already filter by.
              </p>

              <h2 id="order">The order is the whole thing</h2>
              <p>
                Every stage trades something: money, complexity, or correctness. Replicas and
                caches make the app faster and slightly wrong on purpose. Queues make{' '}
                <em>done</em> mean <em>promised</em>. The skill is not knowing the names of the
                boxes; it is knowing which one you need next and what it costs. That is why{' '}
                <code>grit scale</code> names one.
              </p>
            </div>

            <div className="mt-12 flex items-center justify-between border-t border-border pt-6">
              <Button variant="ghost" asChild>
                <Link href="/docs/deployment">
                  <ArrowLeft className="mr-2 h-4 w-4" />
                  Deployment
                </Link>
              </Button>
              <Button variant="ghost" asChild>
                <Link href="/docs/security/doctor">
                  Project audit (grit doctor)
                  <ArrowRight className="ml-2 h-4 w-4" />
                </Link>
              </Button>
            </div>
          </div>
        </div>
      </main>
    </div>
  )
}
