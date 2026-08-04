import Link from 'next/link'
import { ArrowLeft, ArrowRight, Gauge, AlertTriangle, ScrollText, Scale } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'
import { CodeBlock } from '@/components/code-block'
import { getDocMetadata } from '@/config/docs-metadata'
import {
  FRAMEWORKS,
  SCENARIOS,
  UNMEASURED,
  ratio,
  isAppBound,
  isDbBound,
  gritCpu,
} from '@/config/benchmarks'

export const metadata = getDocMetadata('/docs/benchmarks')

/* Every number on this page comes from config/benchmarks.ts, which is also what
   the homepage chart reads — there is one set of figures on this site, not two.
   Each row is a head-to-head pair: three repetitions per scenario, medians, zero
   failed requests, the same 10,000 rows restored before every single run. */

export default function BenchmarksPage() {
  return (
    <div className="min-h-screen bg-background isolate">
      <SiteHeader />
      <DocsSidebar />
      <main className="lg:pl-64">
        <div className="container max-w-screen-xl py-10 px-6">
          <div className="max-w-3xl">
            <div className="mb-10">
              <span className="tag-mono text-primary/80 mb-3 block">Benchmarks</span>
              <h1 className="text-4xl font-bold tracking-tight mb-4">
                The same CRUD API, built {FRAMEWORKS.length + 1} times
              </h1>
              <p className="text-lg text-muted-foreground leading-relaxed">
                One products table, four scenarios, one shared Postgres, identical container
                limits, and every framework on its own ecosystem&apos;s ORM &mdash; GORM, Eloquent,
                Prisma, Drizzle. Each framework was run head to head against Grit, back to back,
                three times, and each has a page showing exactly how to reproduce it.
              </p>
            </div>

            {/* ── Why pairs ────────────────────────────────────────── */}
            <section className="mb-12">
              <div className="flex gap-3 rounded-xl border border-primary/25 bg-primary/[0.04] px-5 py-4">
                <Scale className="h-5 w-5 shrink-0 text-primary mt-0.5" />
                <div>
                  <h2 className="font-semibold text-foreground mb-1.5">
                    Read the ratios, not the raw numbers
                  </h2>
                  <p className="text-sm text-muted-foreground leading-relaxed">
                    Grit&apos;s single-row read measured <strong>6,600</strong> req/s in the Bun
                    pair, <strong>4,392</strong> in the Encore pair and <strong>1,635</strong> in
                    the Express pair &mdash; from an identical binary, as hours of write scenarios
                    accumulated in Postgres. The machine drifts across a session. Within a pair
                    both sides ran minutes apart under the same conditions, so the ratio survives;
                    lining Bun&apos;s 3,196 up against Encore&apos;s 438 would imply a shared
                    baseline that does not exist. Every table below is grouped by pair for that
                    reason.
                  </p>
                </div>
              </div>
            </section>

            {/* ── Results ──────────────────────────────────────────── */}
            <section className="mb-12">
              <h2 className="text-2xl font-semibold tracking-tight mb-4 flex items-center gap-2">
                <Gauge className="h-5 w-5 text-primary" />
                Results
              </h2>
              <p className="text-muted-foreground leading-relaxed mb-6">
                50 concurrent users, 30 seconds per run, 4 CPUs and 2 GB per app container. Every
                app shares one Postgres on 8 CPUs &mdash; deliberately more than any of them gets.
                Medians of three repetitions, never a best run.
              </p>

              <div className="space-y-8">
                {FRAMEWORKS.map((fw) => (
                  <div key={fw.slug}>
                    <div className="flex flex-wrap items-baseline justify-between gap-2 mb-2.5">
                      <h3 className="text-lg font-semibold">
                        Grit vs{' '}
                        <Link
                          href={`/docs/benchmarks/${fw.slug}`}
                          className="text-primary hover:underline"
                        >
                          {fw.name}
                        </Link>
                      </h3>
                      <span className="text-xs font-mono text-muted-foreground">{fw.version}</span>
                    </div>

                    <div className="overflow-x-auto rounded-xl border border-border/40">
                      <table className="w-full text-sm border-collapse">
                        <thead>
                          <tr className="border-b border-border/40 bg-muted/30">
                            <th className="text-left px-4 py-2 font-medium">Scenario</th>
                            <th className="text-right px-4 py-2 font-medium">Grit</th>
                            <th className="text-right px-4 py-2 font-medium">{fw.name}</th>
                            <th className="text-right px-4 py-2 font-medium">Ratio</th>
                          </tr>
                        </thead>
                        <tbody>
                          {SCENARIOS.map((s) => {
                            const r = fw.results[s.id]
                            const bound = (x: Parameters<typeof isAppBound>[0]) =>
                              isAppBound(x)
                                ? 'own ceiling'
                                : isDbBound(x)
                                  ? 'DB-bound'
                                  : 'neither saturated'
                            return (
                              <tr
                                key={s.id}
                                className="border-b border-border/20 last:border-0 align-top"
                              >
                                <td className="px-4 py-2.5">
                                  <div className="font-mono text-xs text-primary">{s.id}</div>
                                  <div className="text-xs text-muted-foreground mt-0.5">
                                    {s.request}
                                  </div>
                                </td>
                                <td className="px-4 py-2.5 text-right">
                                  <div className="font-semibold tabular-nums">
                                    {r.gritRps.toLocaleString()}
                                  </div>
                                  <div className="text-[11px] text-muted-foreground">
                                    {r.gritMedian} &middot; {bound(gritCpu(r))}
                                  </div>
                                </td>
                                <td className="px-4 py-2.5 text-right">
                                  <div className="font-semibold tabular-nums">
                                    {r.rps.toLocaleString()}
                                  </div>
                                  <div className="text-[11px] text-muted-foreground">
                                    {r.median} &middot; {bound(r)}
                                  </div>
                                </td>
                                <td
                                  className={
                                    'px-4 py-2.5 text-right font-semibold tabular-nums ' +
                                    (ratio(r) < 1 ? 'text-warning' : 'text-primary')
                                  }
                                >
                                  {ratio(r).toFixed(2)}&times;
                                </td>
                              </tr>
                            )
                          })}
                        </tbody>
                      </table>
                    </div>
                  </div>
                ))}
              </div>

              {UNMEASURED.length > 0 && (
                <div className="mt-8 rounded-xl border border-border/40 p-5">
                  <h3 className="font-semibold mb-2.5">Not measured yet, and why</h3>
                  <ul className="space-y-2.5 text-sm text-muted-foreground leading-relaxed">
                    {UNMEASURED.map((u) => (
                      <li key={u.slug}>
                        <Link
                          href={`/docs/benchmarks/${u.slug}`}
                          className="font-medium text-foreground hover:text-primary"
                        >
                          {u.name}
                        </Link>{' '}
                        &mdash; {u.reason}
                      </li>
                    ))}
                  </ul>
                </div>
              )}
            </section>

            {/* ── Where Grit loses ────────────────────────────────── */}
            <section className="mb-12">
              <h2 className="text-2xl font-semibold tracking-tight mb-4">
                Where Grit loses
              </h2>
              <p className="text-muted-foreground leading-relaxed mb-4">
                <strong>Bun inserts faster than Grit</strong> &mdash; 3,274 req/s against 1,568, a
                bit over twice. That row is on this page for the same reason the others are.
              </p>
              <p className="text-muted-foreground leading-relaxed">
                The cause is GORM&apos;s implicit transaction: every single write is{' '}
                <code>BEGIN</code> + <code>INSERT</code> + <code>COMMIT</code>, three round trips
                where Drizzle sends one. v3.133.0 adds{' '}
                <code>DB_SKIP_DEFAULT_TRANSACTION=true</code> to turn it off, worth roughly a third
                of write throughput &mdash; and leaves it <em>off</em> by default. The generator
                emits models with relations, and saving a parent with children is several INSERTs;
                without the wrapping transaction a failure halfway leaves an invoice holding some
                of its line items. It would make this table look better. It is not worth that.
              </p>
            </section>

            {/* ── Reading it honestly ─────────────────────────────── */}
            <section className="mb-12">
              <h2 className="text-2xl font-semibold tracking-tight mb-4">
                What the CPU column is for
              </h2>
              <p className="text-muted-foreground leading-relaxed mb-4">
                A throughput number only means something when you know what stopped it. Where a
                framework&apos;s own container saturated &mdash; Laravel at 412&ndash;419% of its
                400% allowance, Express at 435&ndash;476%, Bun at 403&ndash;405% &mdash; the figure
                is that framework&apos;s genuine ceiling on this hardware.
              </p>
              <p className="text-muted-foreground leading-relaxed mb-4">
                The read-heavy rows are different. On <code>list</code> and <code>mixed</code> Grit
                used about 120&ndash;155% while Postgres sat near 900%: the database ran out first,
                not the framework. Those rows are a <strong>floor</strong>, and the real gap is
                wider than the ratio shows. Encore.ts is the mirror image &mdash; it never
                saturated anything, sitting at 155&ndash;211% with Postgres at 30&ndash;200%, so
                something else bounded it (most likely the Drizzle / node-postgres path rather than
                its Rust HTTP layer) and its numbers are a floor for Encore too.
              </p>
              <p className="text-muted-foreground leading-relaxed">
                The <code>list</code> endpoint runs a <code>COUNT(*)</code> over the whole table on
                every request, on every framework. That is a fair comparison &mdash; they all do
                the same work &mdash; but at 50 concurrent users it is mostly a measurement of
                Postgres. For the cleanest read of framework overhead, look at <code>show</code>:
                one indexed lookup, one JSON encode, nothing else in the way.
              </p>
            </section>

            {/* ── What the benchmark found ────────────────────────── */}
            <section className="mb-12">
              <h2 className="text-2xl font-semibold tracking-tight mb-4 flex items-center gap-2">
                <AlertTriangle className="h-5 w-5 text-warning" />
                Three bugs this found in Grit
              </h2>
              <p className="text-muted-foreground leading-relaxed mb-4">
                The first draft of these numbers was much worse, and the reason was Grit&apos;s
                fault every time. Building a benchmark you intend to publish is an unusually good
                way to find out your own defaults are wrong.
              </p>

              <div className="rounded-xl border border-border/40 p-5 mb-4">
                <h3 className="font-semibold mb-2">Connection-pool churn (fixed in v3.131.0)</h3>
                <p className="text-sm text-muted-foreground leading-relaxed mb-3">
                  The scaffold shipped <code>SetMaxIdleConns(10)</code> alongside{' '}
                  <code>SetMaxOpenConns(100)</code>. Past ten concurrent requests, a connection
                  handed back to a full idle pool is <em>closed</em> &mdash; and the next request
                  makes Postgres fork a new backend. Under load that is a connection storm, and it
                  shows up as database CPU rather than as anything you would think to look for in
                  the application.
                </p>
                <p className="text-sm text-muted-foreground leading-relaxed">
                  Single-row reads went from ~810 to ~2,720 req/s on the one-line change. Idle now
                  defaults to Open, tunable via <code>DB_MAX_OPEN_CONNS</code> and{' '}
                  <code>DB_MAX_IDLE_CONNS</code>.
                </p>
              </div>

              <div className="rounded-xl border border-border/40 p-5 mb-4">
                <h3 className="font-semibold mb-2">
                  <code>REDIS_URL=</code> did not disable Redis (fixed in v3.132.0)
                </h3>
                <p className="text-sm text-muted-foreground leading-relaxed">
                  <code>getEnv</code> treated an empty value as unset and returned the default, so
                  the asynq worker started anyway, failed to dial, and retried in a tight loop
                  &mdash; burning CPU on reconnects with nothing in the logs but dial errors.
                  Setting it empty now genuinely turns cache, jobs, worker and cron off, and says
                  so once at boot.
                </p>
              </div>

              <div className="rounded-xl border border-border/40 p-5">
                <h3 className="font-semibold mb-2">
                  No prepared-statement cache (fixed in v3.133.0)
                </h3>
                <p className="text-sm text-muted-foreground leading-relaxed">
                  GORM re-sent and Postgres re-planned every query. <code>PrepareStmt</code> is now
                  on by default, so a query that runs a thousand times is planned once per
                  connection. Turn it off with <code>DB_PREPARED_STATEMENTS=false</code> if you run
                  pgbouncer in transaction mode, where server-side prepared statements do not
                  survive between requests.
                </p>
              </div>
            </section>

            {/* ── Fairness ────────────────────────────────────────── */}
            <section className="mb-12">
              <h2 className="text-2xl font-semibold tracking-tight mb-4">
                How the comparison is kept fair
              </h2>
              <p className="text-muted-foreground leading-relaxed mb-4">
                A benchmark is only worth publishing if the loser was given every reasonable
                advantage. These rules apply to every pair; the framework-specific decisions are
                spelled out on each guide page.
              </p>
              <ul className="space-y-2.5 text-muted-foreground leading-relaxed list-disc pl-5">
                <li>
                  <strong>Every framework uses its ecosystem&apos;s ORM</strong> &mdash; GORM,
                  Eloquent, Prisma, Drizzle. Not raw SQL. Hand-written SQL on one side against an
                  ORM on the other measures the ORM, and nobody ships the framework that way.
                </li>
                <li>
                  <strong>Every framework runs in production shape</strong>, not its dev server:
                  php-fpm rather than <code>artisan serve</code>, a cluster worker per CPU for
                  Node, a production build for Next.js.
                </li>
                <li>
                  <strong>Identical schemas</strong>, down to the index on <code>deleted_at</code>.
                  A benchmark where one side has an index the other lacks measures the index.
                </li>
                <li>
                  <strong>The handlers match</strong> &mdash; same default page size and cap, same
                  searchable columns, same sortable allow-list, same{' '}
                  <code>{'{data, meta}'}</code> envelope, same version bump on update.
                </li>
                <li>
                  <strong>No auth on either side.</strong> With a token in play, part of what you
                  measure is JWT parsing rather than the request path.
                </li>
                <li>
                  <strong>One app runs at a time</strong>, the other stopped rather than idle, and
                  a warm-up run is discarded so nobody pays for the other&apos;s cold start.
                </li>
                <li>
                  <strong>The database is truncated and reloaded before every single run</strong>,
                  and the row count is verified at 10,000 before the run is allowed to start.
                </li>
                <li>
                  <strong>Medians of three</strong>, never the best run, and any run with a single
                  failed request is thrown out rather than published.
                </li>
              </ul>
            </section>

            {/* ── Harness gotchas ─────────────────────────────────── */}
            <section className="mb-12">
              <h2 className="text-2xl font-semibold tracking-tight mb-4 flex items-center gap-2">
                <ScrollText className="h-5 w-5 text-primary" />
                Three ways this benchmark lied before it was fixed
              </h2>

              <p className="text-muted-foreground leading-relaxed mb-4">
                All three are worth knowing if you build one of these yourself, because none of
                them announces itself &mdash; you get plausible numbers that are simply wrong.
              </p>

              <div className="rounded-xl border border-border/40 p-5 mb-4">
                <h3 className="font-semibold mb-2">The write scenario poisoned the read ones</h3>
                <p className="text-sm text-muted-foreground leading-relaxed">
                  Inserts persist. So <code>list</code> was running <code>COUNT(*)</code> against a
                  table that grew for the whole session &mdash; 345,680 rows on the Grit side
                  against 30,255 on Laravel&apos;s, because Grit writes faster and therefore
                  polluted its own table harder, then paid for it on every read. Grit&apos;s{' '}
                  <code>list</code> measured 20 req/s that way. With a verified reset before every
                  run it measures 821.
                </p>
              </div>

              <div className="rounded-xl border border-border/40 p-5 mb-4">
                <h3 className="font-semibold mb-2">k6 must run inside the container network</h3>
                <p className="text-sm text-muted-foreground leading-relaxed">
                  On Docker Desktop for Windows a published port goes through a userland proxy, and
                  the proxy hits its ceiling before either framework does. Measured through it: 432
                  req/s with a 3.94 ms floor. Container-to-container, same test: 740 req/s at 1.11
                  ms. Everything on this page is container-to-container.
                </p>
              </div>

              <div className="rounded-xl border border-border/40 p-5">
                <h3 className="font-semibold mb-2">
                  A run that fails completely still prints a number
                </h3>
                <p className="text-sm text-muted-foreground leading-relaxed">
                  When an app was not listening, k6 recorded thousands of instant connection
                  refusals and summarised them as a perfectly plausible request rate with a
                  sub-millisecond median. Nothing in the summary line says &ldquo;all of these were
                  errors.&rdquo; The harness now health-checks before every measurement and refuses
                  to publish any run whose failure rate is above zero.
                </p>
              </div>
            </section>

            {/* ── Reproduce ───────────────────────────────────────── */}
            <section className="mb-12">
              <h2 className="text-2xl font-semibold tracking-tight mb-4">Reproduce it</h2>
              <p className="text-muted-foreground leading-relaxed mb-4">
                The whole harness is in <code>benchmarks/</code> in the Grit repository &mdash;
                compose file, seed data, k6 script and every bench app. One pair per invocation, by
                design: running two at once means they share a Postgres under load and both sets of
                numbers are worthless.
              </p>
              <CodeBlock
                language="bash"
                code={`git clone https://github.com/MUKE-coder/grit && cd grit/benchmarks

# the Grit side is scaffolded, not vendored — build it yourself
grit new grit-bench --api
(cd grit-bench && grit generate resource Product \\
  name:string sku:string:unique price:float stock:int active:bool)

docker compose up -d postgres

./pair.sh bun          # resets all databases, runs 3 x 4 scenarios, one app at a time
python pair-report.py bun`}
              />
              <p className="text-sm text-muted-foreground mt-4 leading-relaxed">
                <code>pair-report.py</code> prints the app and database CPU for every row and
                labels which side saturated, so you can tell a real ceiling from a database limit.
                It refuses to report any scenario where a request failed.
              </p>
            </section>

            {/* ── Per-framework guides ────────────────────────────── */}
            <section className="mb-12">
              <h2 className="text-2xl font-semibold tracking-tight mb-4">
                Step-by-step guides
              </h2>
              <div className="grid gap-2 sm:grid-cols-2">
                {[...FRAMEWORKS, ...UNMEASURED].map((f) => (
                  <Link
                    key={f.slug}
                    href={`/docs/benchmarks/${f.slug}`}
                    className="rounded-lg border border-border/40 px-4 py-3 hover:border-primary/40 hover:bg-muted/30 transition-colors"
                  >
                    <div className="font-medium text-foreground">Grit vs {f.name}</div>
                    <div className="text-xs text-muted-foreground mt-0.5 line-clamp-1">
                      {'tagline' in f ? f.tagline : 'Method published, numbers pending'}
                    </div>
                  </Link>
                ))}
              </div>
            </section>

            <div className="flex items-center justify-between pt-6 border-t border-border/40">
              <Button variant="ghost" asChild>
                <Link href="/docs/deployment">
                  <ArrowLeft className="mr-2 h-4 w-4" />
                  Deployment
                </Link>
              </Button>
              <Button variant="ghost" asChild>
                <Link href="/docs/testing">
                  Performance testing
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
