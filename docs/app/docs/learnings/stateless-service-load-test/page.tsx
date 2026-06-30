import Link from 'next/link'
import { ArrowLeft, ArrowRight, Activity, AlertTriangle, Info, Check, Target } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'
import { CodeBlock } from '@/components/code-block'
import { getDocMetadata } from '@/config/docs-metadata'

export const metadata = getDocMetadata('/docs/learnings/stateless-service-load-test')

export default function StatelessServiceLoadTestPage() {
  return (
    <div className="min-h-screen bg-background isolate">
      <SiteHeader />
      <DocsSidebar />
      <main className="lg:pl-64">
        <div className="container max-w-screen-xl py-10 px-6">
          <div className="max-w-3xl">
            {/* Header */}
            <div className="mb-10">
              <span className="tag-mono text-primary/80 mb-3 inline-flex items-center gap-1.5">
                <Activity className="h-3.5 w-3.5" />
                Learnings · Challenge #1
              </span>
              <h1 className="text-4xl font-bold tracking-tight mb-4 leading-tight">
                Stateless service + k6 load test
              </h1>
              <p className="text-lg text-muted-foreground leading-relaxed">
                Scaffold a stateless Go API with <code>grit new myapp --api</code>,
                load-test the <code>/api/health</code> endpoint with k6, and
                commit a latency chart showing p50 / p95 / p99 of the run.
              </p>
            </div>

            {/* The brief */}
            <div className="mb-12 rounded-2xl border border-primary/20 bg-primary/[0.04] p-6">
              <div className="flex items-start gap-3 mb-3">
                <Target className="h-5 w-5 text-primary mt-0.5 shrink-0" />
                <div>
                  <p className="text-xs font-mono uppercase tracking-wider text-primary mb-1">The challenge</p>
                  <p className="text-foreground leading-relaxed">
                    Scaffold the repo; a stateless HTTP service + health-check endpoint;
                    load-test it (k6 or autocannon); record p50 / p95 / p99.
                  </p>
                </div>
              </div>
              <div className="flex items-start gap-3 pl-8 pt-2 border-t border-primary/10 mt-3">
                <Check className="h-4 w-4 text-emerald-500 mt-0.5 shrink-0" />
                <p className="text-sm text-foreground/80">
                  <span className="font-semibold">Milestone:</span>{' '}
                  A committed latency chart of the service under load.
                </p>
              </div>
            </div>

            {/* Concepts refresher */}
            <div className="mb-12">
              <h2 className="text-2xl font-semibold tracking-tight mb-3">What we&apos;re measuring (and why)</h2>
              <p className="text-muted-foreground leading-relaxed mb-4">
                Two terms drive the whole exercise. Get them right and the rest is mechanical.
              </p>

              <div className="rounded-xl border border-border/50 bg-card/40 p-5 mb-3">
                <h3 className="font-semibold mb-1.5">Stateless HTTP service</h3>
                <p className="text-sm text-muted-foreground leading-relaxed">
                  The server keeps no per-client memory between requests. Every request stands
                  alone — no session in RAM, no in-process counter, no &quot;remember me from last
                  time&quot;. That property is what makes load tests meaningful: throughput scales
                  horizontally, and a slow request is the service&apos;s fault, not state contention.
                </p>
              </div>

              <div className="rounded-xl border border-border/50 bg-card/40 p-5">
                <h3 className="font-semibold mb-1.5">p50 / p95 / p99 latency</h3>
                <p className="text-sm text-muted-foreground leading-relaxed">
                  Percentiles, not averages. <strong>p50</strong> is the median — half of requests
                  were faster, half slower. <strong>p95</strong> is the value 95% of requests beat.{' '}
                  <strong>p99</strong> is what 99% of requests beat — i.e. only 1 in 100 was
                  slower. Averages hide tail latency; percentiles expose it. A service with p50 of
                  5 ms and p99 of 2,000 ms is unusable for that unlucky 1% — and the average won&apos;t
                  show it.
                </p>
              </div>
            </div>

            {/* Prerequisites */}
            <div className="mb-12">
              <h2 className="text-2xl font-semibold tracking-tight mb-4">Prerequisites</h2>
              <ul className="space-y-2.5">
                {[
                  'Go 1.21+ installed (verify with `go version`)',
                  'Grit CLI installed (`go install github.com/MUKE-coder/grit/v3/cmd/grit@latest`) — needs to be from v3.24+ for SQLite support',
                  'k6 — install instructions in Step 4 below',
                  'Either curl or any HTTP client to sanity-check the endpoint',
                ].map((item) => (
                  <li key={item} className="flex items-start gap-2.5 text-sm text-muted-foreground">
                    <Check className="h-4 w-4 text-emerald-500 mt-0.5 shrink-0" strokeWidth={2.5} />
                    <span>{item}</span>
                  </li>
                ))}
              </ul>
            </div>

            {/* STEP 1 */}
            <Step n={1} title="Scaffold the stateless API">
              <p className="text-[14px] text-muted-foreground leading-relaxed mb-4">
                The <code>--api</code> flag tells Grit to produce a headless Go API kit — pure
                Gin + GORM, no frontend at all. That&apos;s exactly what we want: the smallest
                possible surface area to load-test.
              </p>

              <CodeBlock terminal code={`grit new bench-api --api\ncd bench-api`} className="mb-4" />

              <p className="text-[14px] text-muted-foreground leading-relaxed mb-3">
                The scaffolder creates a small monorepo with the Go API inside{' '}
                <code>apps/api/</code>. Trimmed to what matters for this exercise:
              </p>

              <CodeBlock language="text" filename="bench-api/" code={`bench-api/
├── .env                       ← DATABASE_URL, JWT_SECRET, etc. live here
├── .env.example
├── docker-compose.yml         ← Postgres + Redis + MinIO for local dev
├── grit.json
└── apps/
    └── api/
        ├── go.mod
        ├── .air.toml          ← hot reload via air
        ├── cmd/
        │   ├── server/        ← main.go is here — entry point for the API
        │   ├── migrate/
        │   └── seed/
        └── internal/
            ├── config/        ← loads .env, exposes typed Config struct
            ├── database/      ← Postgres connection + AutoMigrate
            ├── handlers/
            ├── middleware/
            ├── models/
            ├── routes/
            │   └── routes.go  ← the /api/health route lives here
            └── services/`} className="mb-4" />

              <Callout tone="warning">
                <strong>Two things to know up front:</strong>
                <ol className="mt-2 ml-4 space-y-1.5 list-decimal">
                  <li>
                    <code>--api</code> still produces a <strong>monorepo</strong> (<code>apps/api/</code>)
                    — not a flat single-folder project. The <code>main.go</code> entry point is at{' '}
                    <code>apps/api/cmd/server/main.go</code>.
                  </li>
                  <li>
                    The <code>.env</code> sits at <strong>project root</strong>{' '}
                    (<code>bench-api/.env</code>), and the config loader expects you to run the
                    server with the project root as your working directory. If you{' '}
                    <code>cd</code> into <code>apps/api/cmd/server</code> and run <code>go run .</code>,{' '}
                    <code>.env</code> won&apos;t be loaded and you&apos;ll see{' '}
                    <code>Failed to load config: DATABASE_URL is required</code>. Step 3 shows the
                    right invocation.
                  </li>
                </ol>
              </Callout>
            </Step>

            {/* STEP 2 */}
            <Step n={2} title="Tour the health-check endpoint">
              <p className="text-[14px] text-muted-foreground leading-relaxed mb-4">
                Grit pre-wires <code>/api/health</code> as a public, no-auth endpoint. Open{' '}
                <code>apps/api/internal/routes/routes.go</code> and find it:
              </p>

              <CodeBlock
                language="go"
                filename="apps/api/internal/routes/routes.go"
                code={`// Health check
r.GET("/api/health", func(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{
        "status":  "ok",
        "version": "0.1.0",
    })
})`}
                className="mb-4"
              />

              <p className="text-[14px] text-muted-foreground leading-relaxed mb-3">
                It&apos;s deliberately tiny — no DB hit, no auth, no allocation beyond the JSON
                response. That gives us a clean read of <em>the framework&apos;s</em> overhead
                (Gin + Go&apos;s net/http) without database or external service noise polluting the
                number.
              </p>

              <Callout tone="info">
                <strong>Why this endpoint is the right one to bench:</strong> it answers the
                question &quot;how fast can my framework hand off a request and serialize a tiny
                JSON response?&quot;. Once you add a DB query or external API call, you&apos;re measuring
                that, not the service. Start clean, then complicate.
              </Callout>
            </Step>

            {/* STEP 3 */}
            <Step n={3} title="Switch to SQLite &amp; run the API">
              <p className="text-[14px] text-muted-foreground leading-relaxed mb-3">
                The scaffold defaults to Postgres (via the docker-compose service it ships).
                For a benchmark we want zero infrastructure, so let&apos;s flip the connection to
                SQLite — Grit&apos;s database package supports both. Open{' '}
                <code>.env</code> at project root and edit the <code>DATABASE_URL</code> line:
              </p>

              <CodeBlock
                language="dotenv"
                filename="bench-api/.env"
                code={`# Database — Postgres (default) or SQLite
#   postgres://...        → Postgres (requires docker compose up -d postgres)
#   sqlite:./app.db       → SQLite file (no Docker, pure-Go driver)
#   sqlite::memory:       → SQLite in memory (great for tests; gone on restart)
DATABASE_URL=sqlite:./bench.db
APP_ENV=production

# Turn Sentinel (WAF) and Pulse (observability) OFF for the benchmark.
# Both sit in the request middleware chain — leaving them on means we'd be
# benchmarking them, not Gin. Re-enable them when you're done.
SENTINEL_ENABLED=false
PULSE_ENABLED=false`}
                className="mb-4"
              />

              <p className="text-[14px] text-muted-foreground leading-relaxed mb-3">
                Now run the server. The Go module lives at{' '}
                <code>apps/api/go.mod</code>, so <code>go run</code> needs to start there.
                <strong> Crucially</strong>: don&apos;t <code>cd</code> all the way into{' '}
                <code>cmd/server</code> — the config loader expects the working directory
                to be <code>apps/api/</code> so it can find <code>../../.env</code>.
              </p>

              <CodeBlock terminal code={`cd apps/api
go run ./cmd/server`} className="mb-4" />

              <p className="text-[14px] text-muted-foreground leading-relaxed mb-3">
                You should see Grit&apos;s startup banner,{' '}
                <code>Database connected successfully</code>, and a line like{' '}
                <code>listening on :8080</code>. In another terminal, prove it&apos;s alive:
              </p>

              <CodeBlock terminal code={`curl -i http://localhost:8080/api/health

HTTP/1.1 200 OK
Content-Type: application/json; charset=utf-8
Content-Length: 33

{"status":"ok","version":"0.1.0"}`} className="mb-4" />

              <Callout tone="warning">
                <strong>Don&apos;t run from <code>apps/api/cmd/server/</code>.</strong>{' '}
                If you <code>cd</code> in and run <code>go run .</code>, the config loader
                can&apos;t find <code>../../.env</code> (which resolves to <code>apps/api/.env</code>{' '}
                from there — wrong location) and exits with{' '}
                <code>Failed to load config: DATABASE_URL is required</code>. The fix is to
                run from <code>apps/api/</code> with <code>go run ./cmd/server</code>.
              </Callout>

              <Callout tone="info">
                <strong>Why we set <code>APP_ENV=production</code>:</strong> Gin runs in
                debug mode by default — it adds non-trivial per-request overhead (extra
                logging, route printing on startup, slower error rendering). Always bench in
                release mode. Restart the server after editing <code>.env</code>.
              </Callout>

              <Callout tone="info">
                <strong>Prefer Postgres?</strong> Skip the .env edit, keep the original{' '}
                <code>postgres://...</code> DSN, and run <code>docker compose up -d postgres</code>{' '}
                from project root before starting the API. The rest of the tutorial works
                identically — the latency numbers will be slightly different (Postgres has its
                own connection round-trip) but the methodology is the same.
              </Callout>
            </Step>

            {/* STEP 4 */}
            <Step n={4} title="Install k6">
              <p className="text-[14px] text-muted-foreground leading-relaxed mb-3">
                k6 is a single binary written in Go that runs JS test scripts. It&apos;s open source
                (Grafana Labs) and the de-facto standard for HTTP load testing.
              </p>

              <div className="rounded-xl border border-border/50 bg-card/40 p-5 mb-4">
                <p className="text-xs font-mono uppercase tracking-wider text-muted-foreground mb-2">macOS</p>
                <CodeBlock terminal code={`brew install k6`} className="mb-3" />

                <p className="text-xs font-mono uppercase tracking-wider text-muted-foreground mb-2">Windows (winget or Chocolatey)</p>
                <CodeBlock terminal code={`winget install k6 --source winget
# or
choco install k6`} className="mb-3" />

                <p className="text-xs font-mono uppercase tracking-wider text-muted-foreground mb-2">Linux (Debian/Ubuntu)</p>
                <CodeBlock terminal code={`sudo gpg -k && sudo gpg --no-default-keyring --keyring /usr/share/keyrings/k6-archive-keyring.gpg --keyserver hkp://keyserver.ubuntu.com:80 --recv-keys C5AD17C747E3415A3642D57D77C6C491D6AC1D69
echo "deb [signed-by=/usr/share/keyrings/k6-archive-keyring.gpg] https://dl.k6.io/deb stable main" | sudo tee /etc/apt/sources.list.d/k6.list
sudo apt-get update && sudo apt-get install k6`} />
              </div>

              <p className="text-[14px] text-muted-foreground leading-relaxed">
                Verify with <code>k6 version</code>. You want v0.50+ for the built-in HTML report
                we&apos;ll use later.
              </p>
            </Step>

            {/* STEP 5 */}
            <Step n={5} title="Write the smoke test">
              <p className="text-[14px] text-muted-foreground leading-relaxed mb-3">
                Before slamming the service with traffic, run a 30-second smoke test with one
                virtual user. This verifies the script works, the endpoint is reachable, and the
                response shape is what you expect.
              </p>

              <p className="text-[14px] text-muted-foreground leading-relaxed mb-3">
                Create a folder for your k6 scripts:
              </p>

              <CodeBlock terminal code={`mkdir -p loadtests && cd loadtests`} className="mb-4" />

              <p className="text-[14px] text-muted-foreground leading-relaxed mb-3">
                Then create <code>smoke.js</code>:
              </p>

              <CodeBlock
                language="javascript"
                filename="loadtests/smoke.js"
                code={`import http from 'k6/http'
import { check, sleep } from 'k6'

export const options = {
  vus: 1,           // one virtual user
  duration: '30s',  // for 30 seconds
  thresholds: {
    http_req_failed:   ['rate<0.01'],    // <1% requests can fail
    http_req_duration: ['p(95)<200'],    // p95 must beat 200ms
  },
}

export default function () {
  const res = http.get('http://localhost:8080/api/health')
  check(res, {
    'status is 200':            (r) => r.status === 200,
    'body has "status":"ok"':   (r) => r.body.includes('"status":"ok"'),
  })
  sleep(1)
}`}
                className="mb-4"
              />

              <p className="text-[14px] text-muted-foreground leading-relaxed mb-3">
                Run it:
              </p>

              <CodeBlock terminal code={`k6 run smoke.js`} className="mb-4" />

              <Callout tone="info">
                <strong>What to look for:</strong> at the end you&apos;ll see{' '}
                <code>checks_succeeded</code> at 100% and{' '}
                <code>http_req_duration</code> with low single-digit milliseconds for p50 and
                p95. If anything fails here, fix it before scaling up — the load test won&apos;t
                magically clarify a broken script.
              </Callout>
            </Step>

            {/* STEP 6 */}
            <Step n={6} title="Write the real load test">
              <p className="text-[14px] text-muted-foreground leading-relaxed mb-3">
                The load test ramps virtual users (VUs) up, holds a peak, then ramps down. That
                shape — ramp → plateau → ramp-down — is the canonical &quot;average load&quot; profile
                from k6&apos;s testing types catalog. It surfaces both steady-state behavior and
                what happens when VUs spin up.
              </p>

              <CodeBlock
                language="javascript"
                filename="loadtests/load.js"
                code={`import http from 'k6/http'
import { check } from 'k6'

export const options = {
  scenarios: {
    average_load: {
      executor: 'ramping-vus',
      startVUs: 0,
      stages: [
        { duration: '30s',  target: 50  },   // ramp to 50 VUs over 30s
        { duration: '1m30s', target: 50 },   // hold 50 VUs for 1m30s
        { duration: '30s',  target: 100 },   // ramp to 100 VUs over 30s
        { duration: '2m',   target: 100 },   // hold 100 VUs for 2m
        { duration: '30s',  target: 0   },   // ramp down
      ],
      gracefulRampDown: '10s',
    },
  },
  thresholds: {
    http_req_failed:   ['rate<0.01'],     // <1% requests can error
    http_req_duration: [
      'p(50)<50',                          // p50 under 50ms
      'p(95)<200',                         // p95 under 200ms
      'p(99)<500',                         // p99 under 500ms
    ],
  },
  summaryTrendStats: ['min', 'med', 'avg', 'p(95)', 'p(99)', 'max'],
}

export default function () {
  const res = http.get('http://localhost:8080/api/health')
  check(res, { 'status is 200': (r) => r.status === 200 })
}`}
                className="mb-4"
              />

              <Callout tone="info">
                <strong>What each part does:</strong> the <code>stages</code> array drives the
                VU count over time. <code>thresholds</code> turns the test into a pass/fail
                check — if p95 climbs over 200 ms, k6 exits with a non-zero code, which is what
                you want in CI. <code>summaryTrendStats</code> tells k6 which percentiles to
                print at the end.
              </Callout>
            </Step>

            {/* STEP 7 */}
            <Step n={7} title="Run the load test &amp; capture the data">
              <p className="text-[14px] text-muted-foreground leading-relaxed mb-3">
                Run with two outputs: the JSON sample stream (for charting later) and the JSON
                summary (for the final aggregated numbers).
              </p>

              <CodeBlock terminal code={`k6 run \\
  --out json=results.jsonl \\
  --summary-export=summary.json \\
  load.js`} className="mb-4" />

              <p className="text-[14px] text-muted-foreground leading-relaxed mb-3">
                You&apos;ll see live progress bars and rolling metrics in the terminal. When it
                finishes, two new files sit in the folder:
              </p>

              <ul className="space-y-1.5 text-sm text-muted-foreground mb-4 list-disc list-inside">
                <li><code>results.jsonl</code> — one JSON sample per request (used for the chart)</li>
                <li><code>summary.json</code> — aggregated metrics (used for the README table)</li>
              </ul>

              <p className="text-[14px] text-muted-foreground leading-relaxed mb-3">
                And in the terminal, the end-of-run summary block. The line you care about most:
              </p>

              <CodeBlock
                language="text"
                code={`http_req_duration..............: avg=3.42ms  min=0.31ms  med=2.81ms  max=121.6ms  p(95)=7.84ms  p(99)=18.2ms
http_reqs.....................: 24812  165.41/s
iteration_duration............: avg=3.78ms
vus...........................: 100    min=0    max=100`}
                className="mb-4"
              />

              <Callout tone="warning">
                <strong>Run on the same machine? Cap your expectations.</strong> Running k6
                and the service on one laptop costs you accuracy — they fight for the same CPU
                and you measure the worst of both. For a serious number, put k6 on a separate
                box (or another VM) on the same network.
              </Callout>
            </Step>

            {/* STEP 8 */}
            <Step n={8} title="Generate the latency chart">
              <p className="text-[14px] text-muted-foreground leading-relaxed mb-3">
                You have three solid options for the chart, in increasing order of effort and
                fidelity:
              </p>

              <div className="space-y-3 mb-5">
                <div className="rounded-xl border border-border/50 bg-card/40 p-5">
                  <div className="flex items-center gap-2 mb-2">
                    <span className="text-[10px] font-mono uppercase tracking-wider text-emerald-500 bg-emerald-500/10 border border-emerald-500/20 px-1.5 py-0.5 rounded-full">EASIEST</span>
                    <p className="font-semibold">k6&apos;s built-in HTML report</p>
                  </div>
                  <p className="text-sm text-muted-foreground leading-relaxed mb-3">
                    k6 v0.50+ ships a one-flag HTML report with percentile lines, request rate,
                    and error rate over time. Add a single env var when you run:
                  </p>
                  <CodeBlock terminal code={`K6_WEB_DASHBOARD=true K6_WEB_DASHBOARD_EXPORT=report.html \\
  k6 run --summary-export=summary.json load.js`} />
                  <p className="text-xs text-muted-foreground/80 mt-2">
                    <code>report.html</code> opens in any browser. Commit it directly — that&apos;s
                    your chart.
                  </p>
                </div>

                <div className="rounded-xl border border-border/50 bg-card/40 p-5">
                  <div className="flex items-center gap-2 mb-2">
                    <span className="text-[10px] font-mono uppercase tracking-wider text-amber-500 bg-amber-500/10 border border-amber-500/20 px-1.5 py-0.5 rounded-full">DIY</span>
                    <p className="font-semibold">Custom chart from results.jsonl</p>
                  </div>
                  <p className="text-sm text-muted-foreground leading-relaxed mb-3">
                    If you want a custom-styled chart in your README, write a short Node script
                    that bins the JSONL into seconds and renders an SVG. Drop this into{' '}
                    <code>loadtests/chart.mjs</code>:
                  </p>
                  <CodeBlock
                    language="javascript"
                    filename="loadtests/chart.mjs"
                    code={`import fs from 'node:fs'

const rows = fs.readFileSync('results.jsonl', 'utf8')
  .trim().split('\\n').map(JSON.parse)
  .filter(r => r.metric === 'http_req_duration' && r.type === 'Point')

// Bin by second since first sample
const start = new Date(rows[0].data.time).getTime()
const buckets = new Map()
for (const r of rows) {
  const t = Math.floor((new Date(r.data.time).getTime() - start) / 1000)
  if (!buckets.has(t)) buckets.set(t, [])
  buckets.get(t).push(r.data.value)
}

const pct = (xs, p) => {
  const s = [...xs].sort((a, b) => a - b)
  return s[Math.floor(s.length * p)] || 0
}

const series = [...buckets.entries()]
  .sort(([a], [b]) => a - b)
  .map(([t, vs]) => ({ t, p50: pct(vs, 0.5), p95: pct(vs, 0.95), p99: pct(vs, 0.99) }))

// emit an SVG line chart
const w = 800, h = 320, pad = 40
const maxY = Math.max(...series.flatMap(s => [s.p50, s.p95, s.p99])) * 1.1
const x = i => pad + (i / (series.length - 1)) * (w - pad * 2)
const y = v => h - pad - (v / maxY) * (h - pad * 2)
const line = key => series.map((s, i) => \`\${i ? 'L' : 'M'}\${x(i)},\${y(s[key])}\`).join(' ')

const svg = \`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 \${w} \${h}">
  <rect width="\${w}" height="\${h}" fill="#0a0a0f"/>
  <path d="\${line('p99')}" fill="none" stroke="#ff6b6b" stroke-width="2"/>
  <path d="\${line('p95')}" fill="none" stroke="#fdcb6e" stroke-width="2"/>
  <path d="\${line('p50')}" fill="none" stroke="#6c5ce7" stroke-width="2"/>
  <text x="\${pad}" y="\${pad - 10}" fill="#e8e8f0" font-family="monospace" font-size="14">Latency (ms) — p50 (purple) · p95 (yellow) · p99 (red)</text>
</svg>\`

fs.writeFileSync('latency.svg', svg)
console.log('wrote latency.svg')`}
                  />
                  <p className="text-xs text-muted-foreground/80 mt-2">
                    Run with <code>node chart.mjs</code> — out pops <code>latency.svg</code> ready
                    to commit and embed in your README.
                  </p>
                </div>

                <div className="rounded-xl border border-border/50 bg-card/40 p-5">
                  <div className="flex items-center gap-2 mb-2">
                    <span className="text-[10px] font-mono uppercase tracking-wider text-sky-500 bg-sky-500/10 border border-sky-500/20 px-1.5 py-0.5 rounded-full">PRO</span>
                    <p className="font-semibold">InfluxDB + Grafana (for repeat use)</p>
                  </div>
                  <p className="text-sm text-muted-foreground leading-relaxed mb-3">
                    For a permanent perf dashboard, point k6 at an InfluxDB and import the
                    official k6 Grafana dashboard (ID <code>2587</code>). Worth setting up once
                    when you&apos;ll be doing repeated runs; overkill for a one-shot.
                  </p>
                  <CodeBlock terminal code={`k6 run --out influxdb=http://localhost:8086/k6 load.js`} />
                </div>
              </div>

              <Callout tone="info">
                For the milestone, <strong>option 1 is enough</strong> — commit{' '}
                <code>report.html</code>. If you want the chart on your README too, run option 2
                and embed the SVG.
              </Callout>
            </Step>

            {/* STEP 9 */}
            <Step n={9} title="Make sense of the numbers">
              <p className="text-[14px] text-muted-foreground leading-relaxed mb-4">
                Numbers without interpretation are noise. Here&apos;s the cheat sheet for what
                each metric means and what &quot;good&quot; looks like for a tiny health endpoint on
                a local box.
              </p>

              <div className="overflow-hidden rounded-xl border border-border/50 mb-5">
                <table className="w-full text-sm">
                  <thead className="bg-muted/40">
                    <tr>
                      <th className="text-left px-4 py-2.5 font-semibold text-foreground/90">Metric</th>
                      <th className="text-left px-4 py-2.5 font-semibold text-foreground/90">Meaning</th>
                      <th className="text-left px-4 py-2.5 font-semibold text-foreground/90">Healthy</th>
                    </tr>
                  </thead>
                  <tbody className="divide-y divide-border/40">
                    <Row m="http_req_duration p50" mean="Median request time end-to-end" good="< 10 ms" />
                    <Row m="http_req_duration p95" mean="95% of requests beat this" good="< 50 ms" />
                    <Row m="http_req_duration p99" mean="99% of requests beat this" good="< 200 ms" />
                    <Row m="http_req_failed" mean="Fraction of requests that errored" good="< 0.1 %" />
                    <Row m="http_reqs / iterations" mean="Throughput (RPS)" good="As high as the box allows" />
                    <Row m="vus" mean="Concurrent virtual users at sample time" good="Matches your stages" />
                    <Row m="http_req_waiting" mean="Time waiting for the first byte (TTFB)" good="Should track p50 closely — if much higher, the server is slow to respond, not slow to send" />
                    <Row m="http_req_connecting" mean="TCP handshake time" good="Effectively 0 with keep-alive on" />
                  </tbody>
                </table>
              </div>

              <Callout tone="warning">
                <strong>Three traps to avoid:</strong>
                <ol className="mt-2 ml-4 space-y-1.5 list-decimal">
                  <li>
                    <strong>Reading averages.</strong> &quot;Avg 5 ms&quot; can hide a 5,000 ms p99.
                    Always look at p95 and p99.
                  </li>
                  <li>
                    <strong>Comparing across hardware.</strong> Numbers on your laptop tell you
                    your laptop&apos;s number. They don&apos;t map to a $5 VPS or a 64-core server.
                  </li>
                  <li>
                    <strong>Forgetting <code>http_req_failed</code>.</strong> 0% failures is
                    table-stakes. If latency looks great but errors are at 12%, your &quot;great
                    latency&quot; is just the fast failures.
                  </li>
                </ol>
              </Callout>
            </Step>

            {/* STEP 10 */}
            <Step n={10} title="Commit the milestone">
              <p className="text-[14px] text-muted-foreground leading-relaxed mb-3">
                You&apos;ve got everything. Wrap the deliverable into git:
              </p>

              <CodeBlock terminal code={`git add loadtests/load.js loadtests/smoke.js loadtests/report.html loadtests/summary.json
git commit -m "perf: k6 load test — p50/p95/p99 chart of /api/health"
git push`} className="mb-4" />

              <p className="text-[14px] text-muted-foreground leading-relaxed">
                Optional but high-value: add a <code>loadtests/README.md</code> with the
                resulting numbers in a table (pasted from <code>summary.json</code>) and the
                hardware you ran it on (CPU, RAM, network). Future-you will thank you when you
                re-run the benchmark after a release.
              </p>
            </Step>

            {/* What I learned */}
            <div className="mb-12 rounded-2xl border border-emerald-500/20 bg-emerald-500/[0.04] p-6">
              <p className="text-xs font-mono uppercase tracking-wider text-emerald-500 mb-3">What I learned</p>
              <ul className="space-y-2 text-sm text-foreground/85 leading-relaxed">
                <li>
                  <strong>Percentiles &gt; averages, every time.</strong> A 4 ms average that
                  hides a 2-second p99 is a production fire waiting to happen.
                </li>
                <li>
                  <strong>Gin in release mode is fast.</strong> Debug mode adds easily 2–3× to
                  the median on this endpoint. Always bench release.
                </li>
                <li>
                  <strong>Co-located bench is fine for relative numbers.</strong> If you&apos;re
                  comparing &quot;before vs after&quot; for one change, running k6 next to the
                  service is OK. For absolute numbers, separate boxes.
                </li>
                <li>
                  <strong>Thresholds turn k6 into CI.</strong> Once the test exits non-zero on a
                  p95 regression, you can wire it into GitHub Actions and catch perf bugs the
                  same way you catch unit-test failures.
                </li>
              </ul>
            </div>

            {/* Extensions */}
            <div className="mb-16">
              <h2 className="text-2xl font-semibold tracking-tight mb-3">Where to go next</h2>
              <ul className="space-y-2.5 text-sm text-muted-foreground leading-relaxed">
                <li className="flex items-start gap-2.5">
                  <ArrowRight className="h-4 w-4 text-primary mt-0.5 shrink-0" />
                  Bench an endpoint that hits the DB. <code>/api/users</code> with seeded
                  data shows you how GORM + Postgres add to the tail.
                </li>
                <li className="flex items-start gap-2.5">
                  <ArrowRight className="h-4 w-4 text-primary mt-0.5 shrink-0" />
                  Add Sentinel rate limiting and re-run. See where p99 starts climbing as the
                  limiter sheds requests.
                </li>
                <li className="flex items-start gap-2.5">
                  <ArrowRight className="h-4 w-4 text-primary mt-0.5 shrink-0" />
                  Move to a <strong>spike test</strong> — same setup, 5s ramp to 500 VUs. The
                  full k6 testing catalogue lives at{' '}
                  <Link href="/docs/testing" className="text-primary hover:underline">/docs/testing</Link>{' '}
                  — six pre-written tests for smoke, average, stress, spike, soak, and
                  breakpoint.
                </li>
                <li className="flex items-start gap-2.5">
                  <ArrowRight className="h-4 w-4 text-primary mt-0.5 shrink-0" />
                  Run the test against the deployed instance instead of localhost. Numbers on
                  your VPS are the numbers that matter.
                </li>
              </ul>
            </div>

            {/* Nav */}
            <div className="flex items-center justify-between pt-6 border-t border-border/30">
              <Button variant="ghost" size="sm" asChild className="text-muted-foreground/60 hover:text-foreground">
                <Link href="/docs/learnings" className="gap-1.5">
                  <ArrowLeft className="h-3.5 w-3.5" />
                  All Learnings
                </Link>
              </Button>
              <Button variant="ghost" size="sm" asChild className="text-muted-foreground/60 hover:text-foreground">
                <Link href="/docs/testing" className="gap-1.5">
                  Performance &amp; Pentest Testing
                  <ArrowRight className="h-3.5 w-3.5" />
                </Link>
              </Button>
            </div>
          </div>
        </div>
      </main>
    </div>
  )
}

/* ─── Step block ─────────────────────────────────────────────────── */
function Step({ n, title, children }: { n: number; title: React.ReactNode; children: React.ReactNode }) {
  return (
    <div className="relative pl-10 mb-14">
      <div className="absolute left-0 top-0 flex h-7 w-7 items-center justify-center rounded-md bg-primary/10 border border-primary/15 text-xs font-mono font-semibold text-primary">
        {n}
      </div>
      <h3 className="text-xl font-semibold mb-3 leading-snug">{title}</h3>
      {children}
    </div>
  )
}

/* ─── Callout ────────────────────────────────────────────────────── */
function Callout({ tone, children }: { tone: 'info' | 'warning'; children: React.ReactNode }) {
  const Icon = tone === 'warning' ? AlertTriangle : Info
  const styles =
    tone === 'warning'
      ? 'border-amber-500/30 bg-amber-500/[0.05] text-amber-50/85'
      : 'border-sky-500/30 bg-sky-500/[0.05] text-sky-50/85'
  const iconColor = tone === 'warning' ? 'text-amber-500' : 'text-sky-500'
  return (
    <div className={`rounded-xl border ${styles} p-4 my-2 flex gap-3 text-sm leading-relaxed`}>
      <Icon className={`h-4 w-4 ${iconColor} mt-0.5 shrink-0`} />
      <div className="text-foreground/85">{children}</div>
    </div>
  )
}

/* ─── Metric row ─────────────────────────────────────────────────── */
function Row({ m, mean, good }: { m: string; mean: string; good: string }) {
  return (
    <tr className="hover:bg-card/40 transition-colors">
      <td className="px-4 py-2.5 align-top font-mono text-[12.5px] text-primary/90">{m}</td>
      <td className="px-4 py-2.5 align-top text-[13px] text-foreground/80">{mean}</td>
      <td className="px-4 py-2.5 align-top text-[13px] text-muted-foreground">{good}</td>
    </tr>
  )
}
