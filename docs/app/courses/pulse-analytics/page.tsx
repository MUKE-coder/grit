import Link from 'next/link'
import { CheckCircle2 } from 'lucide-react'
import { SiteHeader } from '@/components/site-header'
import { GridFrame } from '@/components/grid-frame'
import { CodeBlock, Challenge, Note, Tip, Definition, Code, CourseNav, CourseFooter } from '@/components/course-components'
import type { Metadata } from 'next'

export const metadata: Metadata = {
  title: 'Performance Analytics with Pulse: Tracing, Metrics & Monitoring',
  description: 'Deep dive into Pulse observability. Learn request tracing, database monitoring, runtime metrics, error tracking, health checks, Prometheus export, and performance baselining.',
}

export default function PulseAnalyticsCourse() {
  return (
    <div className="min-h-screen bg-background isolate">
      <SiteHeader />
      <GridFrame />

      <main className="max-w-3xl mx-auto px-6 py-12">
        {/* Breadcrumb */}
        <div className="flex items-center gap-2 text-sm text-muted-foreground mb-8">
          <Link href="/courses" className="hover:text-foreground transition-colors">Courses</Link>
          <span>/</span>
          <span className="text-foreground">Pulse Analytics</span>
        </div>

        {/* Header */}
        <div className="mb-10">
          <div className="flex items-center gap-2 mb-3">
            <span className="text-xs font-mono font-medium text-primary bg-primary/10 px-2 py-0.5 rounded">Standalone Course</span>
            <span className="text-xs text-muted-foreground">~30 min</span>
            <span className="text-xs text-muted-foreground">•</span>
            <span className="text-xs text-muted-foreground">10 challenges</span>
          </div>
          <h1 className="text-3xl md:text-4xl font-bold text-foreground mb-4">
            Performance Analytics with Pulse: Tracing, Metrics & Monitoring
          </h1>
          <p className="text-lg text-muted-foreground leading-relaxed">
            You cannot improve what you cannot measure. When your API is slow, which endpoint is the
            bottleneck? When memory usage spikes, what caused it? When errors happen in production,
            how do you find them? In this course, you will learn how Pulse — Grit{"'"}s built-in
            observability SDK — gives you complete visibility into your running application.
          </p>
        </div>

        <hr className="border-border/40 mb-10" />

        {/* ═══ Section 1: What is Observability? ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">What is Observability?</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Observability is the ability to understand what your application is doing internally by
            examining its external outputs. In production, you cannot attach a debugger or add print
            statements. You need systems that continuously collect and present data about your
            application{"'"}s behavior.
          </p>

          <Definition term="Observability">
            The measure of how well you can understand a system{"'"}s internal state from its external
            outputs. A highly observable system lets you answer questions like {'"'}why is this endpoint
            slow?{'"'} or {'"'}what caused that error?{'"'} without deploying new code or restarting
            the application.
          </Definition>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Observability rests on three pillars:
          </p>

          <Definition term="Logs">
            Timestamped text records of events: {'"'}User 42 logged in at 10:03:22,{'"'} {'"'}Database
            query took 340ms,{'"'} {'"'}Error: connection refused.{'"'} Logs tell you WHAT happened.
            They are the most basic form of observability.
          </Definition>

          <Definition term="Metrics">
            Numerical measurements over time: request count, response time (p50, p95, p99), memory
            usage, CPU utilization, active database connections. Metrics tell you HOW MUCH is
            happening. They are cheap to collect and good for alerting.
          </Definition>

          <Definition term="Traces">
            End-to-end records of a single request as it flows through your system. A trace shows
            every step: the middleware, the handler, the database query, the external API call — with
            timing for each step. Traces tell you WHERE time is being spent.
          </Definition>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Most observability solutions (Datadog, New Relic, Grafana Cloud) are external services
            that cost money and send your data to third-party servers. Pulse is different — it{"'"}s
            built into your Grit application and runs locally. Your data stays with you.
          </p>

          <Challenge number={1} title="Three Pillars">
            <p>In your own words, explain the difference between logging, metrics, and tracing.
            Give a concrete example of when you would use each one. For instance: when would a
            log be more useful than a metric? When would a trace be more useful than a log?</p>
          </Challenge>
        </section>

        {/* ═══ Section 2: What is Pulse? ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">What is Pulse?</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Pulse is Grit{"'"}s built-in observability SDK. It{"'"}s embedded directly in your API
            as middleware and provides a dashboard at <Code>/pulse/ui</Code>. No external services,
            no API keys, no monthly bills — just start your Grit API and open the dashboard.
          </p>

          <p className="text-muted-foreground leading-relaxed mb-4">
            What Pulse provides:
          </p>

          <ul className="space-y-2 text-muted-foreground mb-4">
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Request tracing</strong> — every API request timed and recorded with endpoint, method, status, and duration</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Database monitoring</strong> — SQL query tracking, duration, and connection pool status</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Runtime metrics</strong> — goroutine count, memory usage, garbage collection stats</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Error tracking</strong> — automatic error capture with stack traces and frequency</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Health checks</strong> — <Code>GET /pulse/health</Code> for load balancers and uptime monitoring</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Prometheus export</strong> — <Code>/pulse/metrics</Code> for integration with Grafana</li>
          </ul>

          <Tip>
            Pulse is designed for self-hosting. Your request data, error logs, and performance metrics
            never leave your infrastructure. This makes it suitable for applications with strict data
            residency requirements (GDPR, HIPAA) where sending telemetry to external services is not
            permitted.
          </Tip>

          <Challenge number={2} title="Open the Dashboard">
            <p>Start your Grit API and open <Code>http://localhost:8080/pulse/ui</Code> in your
            browser. What sections does the dashboard have? Take note of the main navigation
            areas — you will explore each one in the following sections.</p>
          </Challenge>
        </section>

        {/* ═══ Section 3: Request Tracing ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Request Tracing</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Every HTTP request to your API is automatically traced by Pulse. For each request,
            Pulse records:
          </p>

          <ul className="space-y-2 text-muted-foreground mb-4">
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Endpoint</strong> — the URL path (e.g., <Code>/api/users</Code>)</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Method</strong> — GET, POST, PUT, DELETE</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Status code</strong> — 200, 201, 400, 401, 404, 500</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Duration</strong> — how long the request took in milliseconds</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Timestamp</strong> — when the request was made</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Client IP</strong> — the origin of the request</li>
          </ul>

          <p className="text-muted-foreground leading-relaxed mb-4">
            The dashboard displays this data in a table, sorted by most recent. You can filter by
            endpoint, status code, or time range. The most useful view is sorting by duration — this
            instantly reveals your slowest endpoints.
          </p>

          <CodeBlock filename="What You See in Pulse">
{`┌─────────┬───────────────────┬────────┬──────────┬─────────────────┐
│ Method  │ Endpoint          │ Status │ Duration │ Timestamp       │
├─────────┼───────────────────┼────────┼──────────┼─────────────────┤
│ POST    │ /api/auth/login   │ 200    │ 142ms    │ 10:03:22.481    │
│ GET     │ /api/users        │ 200    │ 23ms     │ 10:03:21.102    │
│ GET     │ /api/products     │ 200    │ 89ms     │ 10:03:20.847    │
│ POST    │ /api/products     │ 201    │ 31ms     │ 10:03:19.223    │
│ GET     │ /api/users/5      │ 404    │ 4ms      │ 10:03:18.991    │
│ DELETE  │ /api/products/3   │ 200    │ 12ms     │ 10:03:17.445    │
└─────────┴───────────────────┴────────┴──────────┴─────────────────┘`}
          </CodeBlock>

          <Note>
            A response time under 100ms is good. Under 50ms is excellent. If you see endpoints
            consistently over 200ms, they likely have inefficient database queries — the Database
            Monitoring section (next) will help you find them.
          </Note>

          <Challenge number={3} title="Find the Slowest Endpoint">
            <p>Make at least 20 requests to different endpoints in your API — use the Scalar docs,
            curl, or the frontend. Then check the Pulse dashboard. Sort by duration. Which endpoint
            is the slowest? How much slower is it compared to the fastest? Can you guess why?</p>
          </Challenge>
        </section>

        {/* ═══ Section 4: Database Monitoring ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Database Monitoring</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Most API performance issues are database issues. A single N+1 query can turn a 10ms
            endpoint into a 2-second endpoint. Pulse tracks every SQL query executed by GORM:
          </p>

          <ul className="space-y-2 text-muted-foreground mb-4">
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Query text</strong> — the actual SQL that was executed</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Duration</strong> — how long the query took</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Rows affected</strong> — how many rows were returned or modified</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Frequency</strong> — how many times this query pattern was executed</li>
          </ul>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Pulse also shows your connection pool status:
          </p>

          <CodeBlock filename="Connection Pool Metrics">
{`Database Connection Pool
─────────────────────────
Idle connections:    5
Active connections:  2
Max connections:     25
Wait count:          0     ← requests waiting for a connection
Max idle time:       10m`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            If {'"'}Wait count{'"'} is greater than zero, your pool is too small — requests are
            waiting for an available database connection. If {'"'}Idle connections{'"'} is consistently
            at max, your pool is too large — those connections consume memory on the database server.
          </p>

          <Tip>
            Look for query patterns, not individual queries. If you see <Code>SELECT * FROM products
            WHERE id = ?</Code> executed 100 times in a single request, you have an N+1 problem.
            The fix is usually a GORM <Code>Preload()</Code> to load relationships in a single query.
          </Tip>

          <Challenge number={4} title="Monitor Your Queries">
            <p>Generate a resource with many records (e.g., seed 100 products). Make a request
            to list them with pagination. Check the Pulse database monitoring section. How many
            SQL queries were executed for that single API request? What is the slowest query?</p>
          </Challenge>
        </section>

        {/* ═══ Section 5: Runtime Metrics ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Runtime Metrics</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Go provides rich runtime information through its <Code>runtime</Code> package. Pulse
            collects and displays these metrics so you can understand how your application is using
            system resources.
          </p>

          <Definition term="Goroutine">
            Go{"'"}s lightweight thread. Your API uses goroutines to handle concurrent requests —
            each incoming HTTP request is processed in its own goroutine. Unlike OS threads (which
            consume ~1MB of memory each), goroutines start at just 2KB. Thousands of goroutines
            are normal and expected.
          </Definition>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Key runtime metrics Pulse tracks:
          </p>

          <div className="overflow-x-auto mb-6">
            <table className="w-full text-sm border border-border/40 rounded-lg overflow-hidden">
              <thead>
                <tr className="bg-muted/20">
                  <th className="text-left px-4 py-3 font-semibold text-foreground border-b border-border/40">Metric</th>
                  <th className="text-left px-4 py-3 font-semibold text-foreground border-b border-border/40">What it Means</th>
                  <th className="text-left px-4 py-3 font-semibold text-foreground border-b border-border/40">Normal Range</th>
                </tr>
              </thead>
              <tbody className="text-muted-foreground">
                <tr className="border-b border-border/20">
                  <td className="px-4 py-3 font-medium text-foreground">Goroutines</td>
                  <td className="px-4 py-3">Number of active goroutines</td>
                  <td className="px-4 py-3">10-100 (idle), 100-1000+ (under load)</td>
                </tr>
                <tr className="border-b border-border/20">
                  <td className="px-4 py-3 font-medium text-foreground">Heap Memory</td>
                  <td className="px-4 py-3">Memory allocated for Go objects</td>
                  <td className="px-4 py-3">20-200 MB depending on data size</td>
                </tr>
                <tr className="border-b border-border/20">
                  <td className="px-4 py-3 font-medium text-foreground">System Memory</td>
                  <td className="px-4 py-3">Total memory obtained from the OS</td>
                  <td className="px-4 py-3">Higher than heap (includes stack, GC overhead)</td>
                </tr>
                <tr className="border-b border-border/20">
                  <td className="px-4 py-3 font-medium text-foreground">GC Pause</td>
                  <td className="px-4 py-3">Time the garbage collector pauses your app</td>
                  <td className="px-4 py-3">&lt;1ms is good, &gt;10ms is concerning</td>
                </tr>
                <tr>
                  <td className="px-4 py-3 font-medium text-foreground">GC Runs</td>
                  <td className="px-4 py-3">Number of garbage collection cycles</td>
                  <td className="px-4 py-3">Frequency depends on allocation rate</td>
                </tr>
              </tbody>
            </table>
          </div>

          <Note>
            A steadily increasing goroutine count that never decreases is a goroutine leak. This
            usually means a goroutine is blocked waiting for something that will never happen (a
            channel that is never closed, a lock that is never released). Pulse helps you spot
            these leaks by showing goroutine count over time.
          </Note>

          <Challenge number={5} title="Check Your Runtime">
            <p>Open the Pulse dashboard and find the runtime metrics section. How many goroutines
            is your API running right now? What is the current memory usage? Make 50 rapid requests
            and check again — did the goroutine count increase? Did it come back down?</p>
          </Challenge>
        </section>

        {/* ═══ Section 6: Error Tracking ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Error Tracking</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            When your API returns a 4xx or 5xx status code, Pulse captures the error automatically.
            For each error, it records:
          </p>

          <ul className="space-y-2 text-muted-foreground mb-4">
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Error message</strong> — what went wrong</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Stack trace</strong> — where in the code the error occurred</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Endpoint</strong> — which URL triggered the error</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Frequency</strong> — how many times this error has occurred</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Last occurrence</strong> — when it happened most recently</li>
          </ul>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Errors are grouped by type and endpoint. If the same validation error occurs 500 times
            on the registration endpoint, it shows as one entry with a count of 500 — not 500
            separate entries. This makes it easy to prioritize: fix the most frequent errors first.
          </p>

          <CodeBlock filename="Error Tracking View">
{`┌───────────────────────────┬──────────┬───────┬───────────────────┐
│ Error                     │ Endpoint │ Count │ Last Seen         │
├───────────────────────────┼──────────┼───────┼───────────────────┤
│ record not found          │ GET /api │ 23    │ 2 minutes ago     │
│                           │ /users/  │       │                   │
│ validation: email required│ POST /api│ 8     │ 15 minutes ago    │
│                           │ /auth/   │       │                   │
│ duplicate key: email      │ POST /api│ 3     │ 1 hour ago        │
│                           │ /auth/   │       │                   │
└───────────────────────────┴──────────┴───────┴───────────────────┘`}
          </CodeBlock>

          <Challenge number={6} title="Trigger and Track an Error">
            <p>Cause an intentional error: request a resource with a non-existent ID (e.g.,
            <Code>GET /api/users/99999</Code>). Check the Pulse error tracking section. Does the
            error appear? What information does Pulse capture about it? Trigger it 5 more times
            and watch the count increase.</p>
          </Challenge>
        </section>

        {/* ═══ Section 7: Health Checks ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Health Checks</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            The health check endpoint at <Code>GET /pulse/health</Code> returns the current status
            of your API and its dependencies. This endpoint is used by:
          </p>

          <ul className="space-y-2 text-muted-foreground mb-4">
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Load balancers</strong> — to route traffic only to healthy instances</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Kubernetes</strong> — as a liveness probe to restart unhealthy pods</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Uptime monitors</strong> — services like UptimeRobot or Pingdom that alert you when your API goes down</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Deployment scripts</strong> — to verify a new deployment is healthy before switching traffic</li>
          </ul>

          <CodeBlock filename="Health Check Response">
{`GET /pulse/health

{
  "status": "healthy",
  "uptime": "4h 23m 12s",
  "checks": {
    "database": "connected",
    "redis": "connected",
    "storage": "connected"
  },
  "version": "1.0.0"
}`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            If any dependency is unhealthy (database disconnected, Redis unreachable), the status
            changes to {'"'}degraded{'"'} or {'"'}unhealthy{'"'} and the HTTP status code changes
            from 200 to 503 (Service Unavailable).
          </p>

          <Challenge number={7} title="Check Your Health">
            <p>Call <Code>GET /pulse/health</Code> using your browser or curl. What does the response
            contain? What status code is returned? What happens if you stop your database (Docker
            container) and call the health endpoint again?</p>
          </Challenge>
        </section>

        {/* ═══ Section 8: Prometheus Export ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Prometheus Export</h2>

          <Definition term="Prometheus">
            An open-source monitoring system that collects metrics from your applications by
            {'"'}scraping{'"'} HTTP endpoints at regular intervals (typically every 15 seconds).
            It stores the data as time series and provides a query language (PromQL) for analysis.
            Grafana is the standard tool for visualizing Prometheus data as dashboards and graphs.
          </Definition>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Pulse exposes a <Code>/pulse/metrics</Code> endpoint in Prometheus format. This is a
            plain-text format that Prometheus understands:
          </p>

          <CodeBlock filename="GET /pulse/metrics (excerpt)">
{`# HELP http_requests_total Total number of HTTP requests
# TYPE http_requests_total counter
http_requests_total{method="GET",endpoint="/api/users",status="200"} 142
http_requests_total{method="POST",endpoint="/api/auth/login",status="200"} 38
http_requests_total{method="GET",endpoint="/api/products",status="200"} 67

# HELP http_request_duration_seconds HTTP request duration in seconds
# TYPE http_request_duration_seconds histogram
http_request_duration_seconds_bucket{endpoint="/api/users",le="0.01"} 120
http_request_duration_seconds_bucket{endpoint="/api/users",le="0.05"} 138
http_request_duration_seconds_bucket{endpoint="/api/users",le="0.1"} 141

# HELP go_goroutines Number of goroutines
# TYPE go_goroutines gauge
go_goroutines 42

# HELP go_memstats_alloc_bytes Number of bytes allocated
# TYPE go_memstats_alloc_bytes gauge
go_memstats_alloc_bytes 2.4813e+07`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            To use this with Grafana, add a Prometheus data source pointing at
            <Code> http://your-api:8080/pulse/metrics</Code> and build dashboards with panels for
            request rate, error rate, response time percentiles, and resource usage.
          </p>

          <Tip>
            You do not need Prometheus and Grafana to use Pulse. The built-in dashboard at
            <Code> /pulse/ui</Code> covers most needs. The Prometheus endpoint is for teams that
            already have a monitoring stack and want to integrate Grit{"'"}s metrics into their
            existing dashboards.
          </Tip>

          <Challenge number={8} title="Read the Metrics">
            <p>Visit <Code>http://localhost:8080/pulse/metrics</Code> in your browser. Can you read
            the Prometheus format? Find the total request count, the goroutine count, and the memory
            usage. What does the <Code># TYPE</Code> line tell you about each metric (counter vs
            gauge vs histogram)?</p>
          </Challenge>
        </section>

        {/* ═══ Section 9: Configuration ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Configuration</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Pulse is configured through environment variables in your <Code>.env</Code> file:
          </p>

          <div className="overflow-x-auto mb-6">
            <table className="w-full text-sm border border-border/40 rounded-lg overflow-hidden">
              <thead>
                <tr className="bg-muted/20">
                  <th className="text-left px-4 py-3 font-semibold text-foreground border-b border-border/40">Variable</th>
                  <th className="text-left px-4 py-3 font-semibold text-foreground border-b border-border/40">Default</th>
                  <th className="text-left px-4 py-3 font-semibold text-foreground border-b border-border/40">Purpose</th>
                </tr>
              </thead>
              <tbody className="text-muted-foreground">
                <tr className="border-b border-border/20">
                  <td className="px-4 py-3 font-medium text-foreground"><Code>PULSE_ENABLED</Code></td>
                  <td className="px-4 py-3"><Code>true</Code></td>
                  <td className="px-4 py-3">Enable or disable Pulse entirely</td>
                </tr>
                <tr className="border-b border-border/20">
                  <td className="px-4 py-3 font-medium text-foreground"><Code>PULSE_USERNAME</Code></td>
                  <td className="px-4 py-3"><Code>admin</Code></td>
                  <td className="px-4 py-3">Username for the Pulse dashboard</td>
                </tr>
                <tr>
                  <td className="px-4 py-3 font-medium text-foreground"><Code>PULSE_PASSWORD</Code></td>
                  <td className="px-4 py-3"><Code>pulse</Code></td>
                  <td className="px-4 py-3">Password for the Pulse dashboard</td>
                </tr>
              </tbody>
            </table>
          </div>

          <CodeBlock filename=".env">
{`# Pulse Configuration
PULSE_ENABLED=true
PULSE_USERNAME=admin
PULSE_PASSWORD=your-secure-password-here`}
          </CodeBlock>

          <Note>
            Always change the default Pulse credentials in production. The dashboard exposes
            sensitive information about your API{"'"}s internals — response times, error messages,
            database queries, memory usage. This data is valuable to attackers for reconnaissance.
          </Note>

          <Challenge number={9} title="Secure Your Dashboard">
            <p>Change the Pulse password in your <Code>.env</Code> file. Restart the API. Try
            accessing <Code>/pulse/ui</Code> — are you prompted for credentials? Verify that the
            old password no longer works and the new one does.</p>
          </Challenge>
        </section>

        {/* ═══ Section 10: Summary + Final Challenge ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Summary</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Let{"'"}s review what Pulse provides:
          </p>

          <div className="overflow-x-auto mb-6">
            <table className="w-full text-sm border border-border/40 rounded-lg overflow-hidden">
              <thead>
                <tr className="bg-muted/20">
                  <th className="text-left px-4 py-3 font-semibold text-foreground border-b border-border/40">Feature</th>
                  <th className="text-left px-4 py-3 font-semibold text-foreground border-b border-border/40">Endpoint</th>
                  <th className="text-left px-4 py-3 font-semibold text-foreground border-b border-border/40">Purpose</th>
                </tr>
              </thead>
              <tbody className="text-muted-foreground">
                <tr className="border-b border-border/20">
                  <td className="px-4 py-3 font-medium text-foreground">Dashboard</td>
                  <td className="px-4 py-3"><Code>/pulse/ui</Code></td>
                  <td className="px-4 py-3">Visual overview of all metrics</td>
                </tr>
                <tr className="border-b border-border/20">
                  <td className="px-4 py-3 font-medium text-foreground">Request Tracing</td>
                  <td className="px-4 py-3">Dashboard</td>
                  <td className="px-4 py-3">Find slow endpoints</td>
                </tr>
                <tr className="border-b border-border/20">
                  <td className="px-4 py-3 font-medium text-foreground">Database Monitor</td>
                  <td className="px-4 py-3">Dashboard</td>
                  <td className="px-4 py-3">Find slow queries, pool status</td>
                </tr>
                <tr className="border-b border-border/20">
                  <td className="px-4 py-3 font-medium text-foreground">Runtime Metrics</td>
                  <td className="px-4 py-3">Dashboard</td>
                  <td className="px-4 py-3">Goroutines, memory, GC</td>
                </tr>
                <tr className="border-b border-border/20">
                  <td className="px-4 py-3 font-medium text-foreground">Error Tracking</td>
                  <td className="px-4 py-3">Dashboard</td>
                  <td className="px-4 py-3">Grouped errors with frequency</td>
                </tr>
                <tr className="border-b border-border/20">
                  <td className="px-4 py-3 font-medium text-foreground">Health Check</td>
                  <td className="px-4 py-3"><Code>/pulse/health</Code></td>
                  <td className="px-4 py-3">Load balancer / uptime monitoring</td>
                </tr>
                <tr>
                  <td className="px-4 py-3 font-medium text-foreground">Prometheus</td>
                  <td className="px-4 py-3"><Code>/pulse/metrics</Code></td>
                  <td className="px-4 py-3">Integration with Grafana</td>
                </tr>
              </tbody>
            </table>
          </div>

          <Challenge number={10} title="Final Challenge: Performance Baseline">
            <p>Establish a performance baseline for your API. Record the response time of these
            5 key operations (run each one 5 times and take the average):</p>
            <ol className="mt-2 space-y-2 list-decimal list-inside">
              <li><Code>POST /api/auth/login</Code> — authentication</li>
              <li><Code>GET /api/users</Code> — list with pagination</li>
              <li><Code>POST /api/products</Code> — resource creation</li>
              <li><Code>GET /pulse/health</Code> — health check</li>
              <li><Code>GET /api/products?page=1&page_size=50</Code> — large page size</li>
            </ol>
            <p className="mt-3">Write down the average response time for each. This is your
            baseline. As you add features, run these same tests to detect performance regressions.
            If login goes from 140ms to 800ms after a change, you know something is wrong.</p>
          </Challenge>
        </section>

        {/* ═══ Footer ═══ */}
        <CourseFooter />

        <div className="mt-8">
          <CourseNav
            prev={{ href: '/courses', label: 'All Courses' }}
            next={{ href: '/courses', label: 'More Courses' }}
          />
        </div>
      </main>
    </div>
  )
}
