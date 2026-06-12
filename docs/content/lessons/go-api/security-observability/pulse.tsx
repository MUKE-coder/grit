import { CodeBlock } from '@/components/code-block'
import { TipBox } from '@/components/course/tip-box'
import { Exercise } from '@/components/course/exercise'
import { KnowledgeCheck } from '@/components/course/knowledge-check'

export default function Lesson() {
  return (
    <>
      <p>
        You can&apos;t fix what you can&apos;t see. Pulse is Grit&apos;s
        observability dashboard — p50/p95/p99 latency, RPS per route,
        per-handler timings, slow-query log. No Datadog bill, no separate
        infra. Mounted at <code>/pulse/ui</code> in every Grit API.
      </p>

      <h2>What you see</h2>
      <ul>
        <li>Live request rate (last 1 min / hour / day)</li>
        <li>Latency percentiles per route — p50, p95, p99</li>
        <li>Error rates per route (4xx / 5xx breakdown)</li>
        <li>Slow query log (queries &gt;100ms)</li>
        <li>System metrics — Go runtime, goroutine count, GC stats</li>
      </ul>

      <h2>Why percentiles, not averages</h2>
      <p>
        An average hides outliers. If your endpoint serves 99 fast requests
        (10ms each) and 1 slow one (2,000ms), the average is 30ms — looks
        fine. But for that 1% of users, the experience is unusable. p99
        catches it.
      </p>
      <CodeBlock
        language="text"
        code={`Average     30ms     ← looks fine
p50         10ms     ← typical
p95         12ms     ← still fine
p99         2,000ms  ← someone is having a bad time`}
      />

      <h2>The four numbers to watch</h2>
      <ul>
        <li>
          <strong>p50</strong> — typical experience. Should be {`<`}100ms for
          most CRUD endpoints.
        </li>
        <li>
          <strong>p95</strong> — almost everyone. {`<`}500ms is a healthy
          ceiling.
        </li>
        <li>
          <strong>p99</strong> — the worst 1%. {`<`}2s; over that, something
          is wrong.
        </li>
        <li>
          <strong>Error rate</strong> — should be 0% for healthy endpoints.
          A spike means a regression.
        </li>
      </ul>

      <h2>The slow query log</h2>
      <p>
        Queries that take &gt;100ms get logged. Grit&apos;s Pulse dashboard
        groups them by SQL signature so you see <strong>which query</strong>{' '}
        is slow, not just <strong>that something is</strong>.
      </p>
      <CodeBlock
        language="text"
        code={`Slow queries (last hour):
1. SELECT * FROM orders WHERE customer_id = ?    avg 350ms, 12 calls   ← add an index
2. SELECT * FROM users JOIN ...                  avg 800ms, 1 call     ← spike, investigate
3. SELECT count(*) FROM activities               avg 1200ms, 4 calls   ← needs covering index`}
      />
      <p>
        The most common fix is an index. Add{' '}
        <code>gorm:&quot;index&quot;</code> to the field, run{' '}
        <code>grit migrate</code>, problem solved.
      </p>

      <TipBox tone="success">
        <strong>SQLite vs Postgres storage:</strong> Pulse can store its
        metrics in SQLite (zero infra) or Postgres (shared across replicas).
        Default is SQLite. Set <code>PULSE_STORAGE=postgres</code> when you
        scale out.
      </TipBox>

      <h2>What Pulse doesn&apos;t do</h2>
      <p>
        Pulse is in-process — same binary as your API. That has trade-offs:
      </p>
      <ul>
        <li>
          <strong>No cross-service tracing.</strong> If you have multiple Grit
          APIs talking to each other, you can&apos;t follow a request across
          them. Use OpenTelemetry for that.
        </li>
        <li>
          <strong>Metrics are per-replica.</strong> If you run 3 API
          containers, each has its own Pulse view. The Postgres backend can
          merge if you point all 3 at the same DB.
        </li>
        <li>
          <strong>Last 24 hours by default.</strong> Older data ages out. For
          long-term retention, export to Prometheus / Datadog.
        </li>
      </ul>
      <p>
        For most products, single-binary Pulse is enough. You can always
        graduate.
      </p>

      <h2>Prometheus export — when you grow up</h2>
      <p>
        Pulse exposes <code>/pulse/metrics</code> in Prometheus format. Add
        Prometheus + Grafana, scrape this endpoint, and you have proper
        long-term metrics. Use the Pulse UI for day-to-day; Prometheus for
        historical analysis.
      </p>

      <KnowledgeCheck
        question="Pulse shows your `/api/orders` endpoint at p50=15ms, p95=50ms, p99=1.8s. What's the most likely cause?"
        choices={[
          {
            label: 'Slow database query that fires occasionally',
            correct: true,
            feedback:
              "Right — a small fraction of requests hit a slow path. Most are fast (p95=50ms is healthy), but 1% take 1.8s. Check Pulse's slow query log — there's probably a missing index or an unbounded query.",
          },
          {
            label: 'The whole endpoint is slow — everything needs optimization',
            feedback:
              "Wrong — p50 and p95 are healthy. Most requests are fast. Only the tail is bad.",
          },
          {
            label: 'Cold cache',
            feedback:
              "Possible but would show up as a different pattern — first request slow, then fast for the cache window. Persistent p99 outliers point to a query / lock / external service issue.",
          },
          {
            label: 'Network — probably the client\'s connection',
            feedback:
              "Pulse measures server-side time only. Client network doesn't enter the picture.",
          },
        ]}
      />

      <Exercise
        prompt={
          <>
            <p>Load-test your API to populate Pulse:</p>
            <ol>
              <li>
                Make sure Pulse is on (<code>PULSE_ENABLED=true</code>).
              </li>
              <li>
                Run <code>ab</code> or <code>wrk</code> for 30 seconds against{' '}
                <code>/api/health</code>:
                <CodeBlock
                  terminal
                  code={`ab -n 1000 -c 10 http://localhost:8080/api/health`}
                />
              </li>
              <li>
                Open <code>/pulse/ui</code>. You should see the spike + the
                p50/p95/p99 numbers.
              </li>
              <li>
                Screenshot the route view and paste it in <code>notes.md</code>.
              </li>
            </ol>
          </>
        }
        hint={
          <>
            If <code>ab</code> isn&apos;t installed,{' '}
            <code>brew install apache-utils</code> on Mac. Or use{' '}
            <code>k6</code> (which has its own course entry in Learnings).
          </>
        }
        solution={
          <>
            <p>Healthy Pulse output for /api/health:</p>
            <CodeBlock
              language="text"
              code={`/api/health
RPS:    33.3
p50:    1.2 ms
p95:    3.4 ms
p99:    8.1 ms
Errors: 0`}
            />
            <p>
              Health checks are designed to be {`<`}10ms — anything else
              indicates middleware bloat or contention.
            </p>
          </>
        }
      />

      <h2>What&apos;s next</h2>
      <p>
        Last lesson of the chapter — the <strong>tamper-evident audit log</strong>.
        Every mutation is hashed into a chain so retroactive tampering
        breaks verification. The compliance team will love you.
      </p>
    </>
  )
}
