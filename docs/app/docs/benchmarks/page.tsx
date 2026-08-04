import Link from 'next/link'
import { ArrowLeft, ArrowRight, Gauge, AlertTriangle, ScrollText } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'
import { CodeBlock } from '@/components/code-block'
import { getDocMetadata } from '@/config/docs-metadata'

export const metadata = getDocMetadata('/docs/benchmarks')

/* Every number on this page comes from results/final in the benchmark harness:
   three repetitions per scenario, medians, zero failed requests. The CPU
   columns are not decoration — they are what separates a real ceiling from a
   number that only shows where the database gave out. */

const RESULTS = [
  {
    scenario: 'show',
    request: 'GET /products/:id',
    grit: { rps: '4,722', med: '8.4 ms', p95: '21.8 ms', cpu: '267%' },
    laravel: { rps: '113', med: '136.3 ms', p95: '230.1 ms', cpu: '412%' },
    ratio: '41.8×',
    bound: 'app',
  },
  {
    scenario: 'write',
    request: 'POST /products',
    grit: { rps: '2,425', med: '18.2 ms', p95: '39.9 ms', cpu: '291%' },
    laravel: { rps: '96', med: '119.3 ms', p95: '306.7 ms', cpu: '419%' },
    ratio: '25.4×',
    bound: 'app',
  },
  {
    scenario: 'list',
    request: 'GET /products?page=N',
    grit: { rps: '821', med: '52.2 ms', p95: '130.1 ms', cpu: '130%' },
    laravel: { rps: '96', med: '149.5 ms', p95: '252.3 ms', cpu: '415%' },
    ratio: '8.6×',
    bound: 'db',
  },
  {
    scenario: 'mixed',
    request: '85% list / 10% show / 5% write',
    grit: { rps: '748', med: '58.2 ms', p95: '148.4 ms', cpu: '126%' },
    laravel: { rps: '95', med: '119.5 ms', p95: '226.5 ms', cpu: '412%' },
    ratio: '7.9×',
    bound: 'db',
  },
]

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
              <h1 className="text-4xl font-bold tracking-tight mb-4">Grit vs Laravel</h1>
              <p className="text-lg text-muted-foreground leading-relaxed">
                The same public CRUD resource, built twice. Same Postgres, same container
                limits, same 10,000 rows, same k6 script. Three repetitions, medians reported,
                zero failed requests on either side. The harness is in the repo so you can
                disagree with it in public.
              </p>
            </div>

            {/* ── Results ─────────────────────────────────────────── */}
            <section className="mb-12">
              <h2 className="text-2xl font-semibold tracking-tight mb-4 flex items-center gap-2">
                <Gauge className="h-5 w-5 text-primary" />
                Results
              </h2>

              <div className="overflow-x-auto rounded-xl border border-border/40">
                <table className="w-full text-sm border-collapse">
                  <thead>
                    <tr className="border-b border-border/40 bg-muted/30">
                      <th className="text-left px-4 py-2 font-medium">Scenario</th>
                      <th className="text-right px-4 py-2 font-medium">Grit</th>
                      <th className="text-right px-4 py-2 font-medium">Laravel 13</th>
                      <th className="text-right px-4 py-2 font-medium">Ratio</th>
                    </tr>
                  </thead>
                  <tbody>
                    {RESULTS.map((r) => (
                      <tr key={r.scenario} className="border-b border-border/20 last:border-0">
                        <td className="px-4 py-3 align-top">
                          <div className="font-mono text-xs text-primary">{r.scenario}</div>
                          <div className="text-xs text-muted-foreground mt-0.5">{r.request}</div>
                        </td>
                        <td className="px-4 py-3 text-right align-top">
                          <div className="font-semibold">{r.grit.rps} req/s</div>
                          <div className="text-xs text-muted-foreground mt-0.5">
                            {r.grit.med} median · {r.grit.p95} p95
                          </div>
                          <div className="text-xs text-muted-foreground">{r.grit.cpu} CPU</div>
                        </td>
                        <td className="px-4 py-3 text-right align-top">
                          <div className="font-semibold">{r.laravel.rps} req/s</div>
                          <div className="text-xs text-muted-foreground mt-0.5">
                            {r.laravel.med} median · {r.laravel.p95} p95
                          </div>
                          <div className="text-xs text-muted-foreground">{r.laravel.cpu} CPU</div>
                        </td>
                        <td className="px-4 py-3 text-right align-top">
                          <span className="font-semibold text-primary">{r.ratio}</span>
                          <div className="text-xs text-muted-foreground mt-0.5">
                            {r.bound === 'app' ? 'true ceiling' : 'Grit is a floor'}
                          </div>
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>

              <p className="text-sm text-muted-foreground mt-4 leading-relaxed">
                50 concurrent users, 30 seconds per run, 4 CPUs and 2 GB per container. Both
                sides were given 8 CPUs of Postgres to share &mdash; deliberately more than
                either app gets.
              </p>
            </section>

            {/* ── Reading it honestly ─────────────────────────────── */}
            <section className="mb-12">
              <h2 className="text-2xl font-semibold tracking-tight mb-4">
                What the CPU column is for
              </h2>
              <p className="text-muted-foreground leading-relaxed mb-4">
                A throughput number only means something when you know what stopped it. In every
                run Laravel was the thing that saturated &mdash; it sat at 412&ndash;419% of its
                400% allowance, which makes those figures its genuine ceiling on this hardware.
              </p>
              <p className="text-muted-foreground leading-relaxed mb-4">
                Grit&apos;s two read-heavy rows are different. On <code>list</code> and{' '}
                <code>mixed</code> Grit used about 130% while Postgres sat near 850%: the
                database ran out first, not the framework. Those rows are a{' '}
                <strong>floor</strong>. Given a larger database Grit would go higher, and the
                real gap is wider than 8.6× and 7.9× suggest.
              </p>
              <p className="text-muted-foreground leading-relaxed">
                The <code>list</code> endpoint runs a <code>COUNT(*)</code> over the whole table
                on every request, on both sides. That is a fair comparison &mdash; both
                frameworks do the same work &mdash; but at 50 concurrent users it is mostly a
                measurement of Postgres. If you want the cleanest read of framework overhead,
                look at <code>show</code>: one indexed lookup, one JSON encode, nothing else in
                the way.
              </p>
            </section>

            {/* ── What the benchmark found ────────────────────────── */}
            <section className="mb-12">
              <h2 className="text-2xl font-semibold tracking-tight mb-4 flex items-center gap-2">
                <AlertTriangle className="h-5 w-5 text-warning" />
                Two bugs this found in Grit
              </h2>
              <p className="text-muted-foreground leading-relaxed mb-4">
                The first draft of these numbers was much worse, and the reason was Grit&apos;s
                fault both times. Building a benchmark you intend to publish is an unusually good
                way to find your own defaults are wrong.
              </p>

              <div className="rounded-xl border border-border/40 p-5 mb-4">
                <h3 className="font-semibold mb-2">Connection-pool churn (fixed in v3.131.0)</h3>
                <p className="text-sm text-muted-foreground leading-relaxed mb-3">
                  The scaffold shipped <code>SetMaxIdleConns(10)</code> alongside{' '}
                  <code>SetMaxOpenConns(100)</code>. Past ten concurrent requests, a connection
                  handed back to a full idle pool is <em>closed</em> &mdash; and the next request
                  makes Postgres fork a new backend. Under load that is a connection storm, and
                  it shows up as database CPU rather than as anything you would think to look for
                  in the application.
                </p>
                <p className="text-sm text-muted-foreground leading-relaxed">
                  Single-row reads went from ~810 to ~2,720 req/s on the one-line change. Idle
                  now defaults to Open, tunable via <code>DB_MAX_OPEN_CONNS</code> and{' '}
                  <code>DB_MAX_IDLE_CONNS</code>.
                </p>
              </div>

              <div className="rounded-xl border border-border/40 p-5">
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
            </section>

            {/* ── Fairness ────────────────────────────────────────── */}
            <section className="mb-12">
              <h2 className="text-2xl font-semibold tracking-tight mb-4">
                How the comparison is kept fair
              </h2>
              <ul className="space-y-2.5 text-muted-foreground leading-relaxed list-disc pl-5">
                <li>
                  <strong>Laravel runs in production shape</strong> &mdash; nginx + php-fpm with
                  32 workers, opcache on with tracing JIT, <code>APP_DEBUG=false</code>, and{' '}
                  <code>artisan optimize</code>. Not <code>artisan serve</code>, which is a
                  single-threaded dev server and would make the whole thing meaningless.
                </li>
                <li>
                  <strong>Not Octane.</strong> Octane keeps the framework booted between requests
                  and is considerably faster, but it is opt-in and not what most Laravel apps
                  run. Benchmarking it and calling it &ldquo;Laravel&rdquo; would be dishonest.
                </li>
                <li>
                  <strong>The controller matches Grit&apos;s generated handler</strong> &mdash;
                  same page size and cap, same searchable columns, same sortable allow-list, same{' '}
                  <code>{'{data, meta}'}</code> envelope, same version bump on update. Plain
                  Eloquent rather than API Resources, because Resources would add a layer
                  Grit&apos;s handler has no equivalent of.
                </li>
                <li>
                  <strong>Identical schemas</strong>, down to the index on <code>deleted_at</code>.
                  A benchmark where one side has an index the other lacks measures the index.
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
                  <strong>Medians of three</strong>, never the best run. Publishing a best run is
                  how a benchmark becomes something nobody else can reproduce.
                </li>
              </ul>
            </section>

            {/* ── Harness gotchas ─────────────────────────────────── */}
            <section className="mb-12">
              <h2 className="text-2xl font-semibold tracking-tight mb-4 flex items-center gap-2">
                <ScrollText className="h-5 w-5 text-primary" />
                Two ways this benchmark lied before it was fixed
              </h2>

              <p className="text-muted-foreground leading-relaxed mb-4">
                Both are worth knowing if you build one of these yourself, because neither
                announces itself &mdash; you get plausible numbers that are simply wrong.
              </p>

              <div className="rounded-xl border border-border/40 p-5 mb-4">
                <h3 className="font-semibold mb-2">The write scenario poisoned the read ones</h3>
                <p className="text-sm text-muted-foreground leading-relaxed">
                  Inserts persist. So <code>list</code> was running <code>COUNT(*)</code> against
                  a table that grew for the whole session &mdash; 345,680 rows on the Grit side
                  against 30,255 on Laravel&apos;s, because Grit writes faster and therefore
                  polluted its own table harder, then paid for it on every read. Grit&apos;s{' '}
                  <code>list</code> measured 20 req/s that way. With a reset before every run it
                  measures 821. The harness now truncates and reloads the same 10,000 rows before
                  each run.
                </p>
              </div>

              <div className="rounded-xl border border-border/40 p-5">
                <h3 className="font-semibold mb-2">k6 must run inside the container network</h3>
                <p className="text-sm text-muted-foreground leading-relaxed">
                  On Docker Desktop for Windows a published port goes through a userland proxy,
                  and the proxy hits its ceiling before either framework does. Measured through
                  it: 432 req/s with a 3.94 ms floor. Container-to-container, same test: 740
                  req/s at 1.11 ms. Everything on this page is container-to-container.
                </p>
              </div>
            </section>

            {/* ── Reproduce ───────────────────────────────────────── */}
            <section className="mb-12">
              <h2 className="text-2xl font-semibold tracking-tight mb-4">Reproduce it</h2>
              <CodeBlock
                language="bash"
                code={`# one Postgres, two databases, identical seed
docker compose up -d postgres
docker compose exec -T postgres psql -U bench -d bench \\
  -c "CREATE DATABASE bench_grit;" -c "CREATE DATABASE bench_laravel;"

(cd grit-bench && grit migrate)
docker compose up -d laravel
docker compose exec -T laravel php artisan migrate --force
docker compose exec -T laravel php artisan optimize

for db in bench_grit bench_laravel; do
  docker compose exec -T postgres psql -q -U bench -d $db < seed/products.sql
done

./final.sh            # 3 reps x 4 scenarios x 2 apps, one app at a time
python aggregate.py   # medians, plus which rows are DB-bound`}
              />
              <p className="text-sm text-muted-foreground mt-4 leading-relaxed">
                <code>aggregate.py</code> prints the app and database CPU for every row and
                labels each scenario app-bound or DB-bound. If you only take one thing from this
                page, take that habit: a throughput number without the CPU beside it is a number
                you cannot interpret.
              </p>
            </section>

            {/* ── Caveats ─────────────────────────────────────────── */}
            <section className="mb-12">
              <h2 className="text-2xl font-semibold tracking-tight mb-4">
                What this does not tell you
              </h2>
              <ul className="space-y-2.5 text-muted-foreground leading-relaxed list-disc pl-5">
                <li>
                  <strong>It is one machine.</strong> A 14-core Windows box running Docker
                  Desktop, with other containers alive on it. Ratios travel better than absolute
                  numbers; run it yourself on the hardware you actually deploy to.
                </li>
                <li>
                  <strong>Four endpoints is not an application.</strong> Real apps have N+1
                  queries, cache layers, queue workers and third-party calls, and those usually
                  dominate long before framework overhead does.
                </li>
                <li>
                  <strong>Throughput is rarely why you pick a framework.</strong> Laravel has an
                  ecosystem two decades deep. If 96 req/s per 4 CPUs is comfortably above what
                  you need &mdash; and for a great many applications it is &mdash; this page is
                  not an argument for switching.
                </li>
                <li>
                  <strong>The list and mixed rows are floors, not ceilings.</strong> Stated again
                  because it is the easiest thing on this page to quote wrongly.
                </li>
              </ul>
            </section>

            <div className="flex items-center justify-between pt-6 border-t border-border/40">
              <Button variant="ghost" asChild>
                <Link href="/docs/testing">
                  <ArrowLeft className="mr-2 h-4 w-4" />
                  Performance &amp; Pentest Testing
                </Link>
              </Button>
              <Button variant="ghost" asChild>
                <Link href="/docs/infrastructure">
                  Infrastructure
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
