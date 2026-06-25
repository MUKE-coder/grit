import { CodeBlock } from '@/components/code-block'
import { TipBox } from '@/components/course/tip-box'
import { Exercise } from '@/components/course/exercise'
import { KnowledgeCheck } from '@/components/course/knowledge-check'

export default function Lesson() {
  return (
    <>
      <p>
        The dashboard widgets lesson covered the dashboard shell —
        stat cards, charts, activity feed — wired by hand. That
        pattern doesn&apos;t scale: every new resource you generate
        deserves its own snapshot on the dashboard, but copy-pasting
        a stat card per resource is a chore. <strong>Per-resource
        dashboard widgets</strong> (shipped in{' '}
        <strong>v3.31.44</strong>) make that automatic: every
        registered resource gets a Total + 30-day sparkline tile and
        a Latest N preview, all scoped by the dashboard&apos;s
        DateFilter.
      </p>

      <h2>What the dashboard looks like now</h2>
      <p>
        Below the existing Quick Access section a new{' '}
        <em>By resource</em> band renders one row per registered
        resource. Each row is two widgets: a stat card on the left
        (Total + 30-day sparkline) and a Latest 5 list on the right
        (the resource&apos;s most recent rows, columns auto-picked).
        A single DateFilter at the top of the dashboard scopes every
        row in lockstep.
      </p>

      <CodeBlock
        language="text"
        code={`┌─ Dashboard ────────────────────────────────────────────────┐
│ Good morning, Alex.                          [Last 7 days▾] │
│                                                            │
│ [Users] [Events 24h] [Notifications] [Resources]           │
│ [── Activity chart ──] [── Severity ──]                    │
│ [── Recent activity feed ──]                               │
│ [Quick access: Users][Categories][Products][System hub]    │
│                                                            │
│ BY RESOURCE              sparkline always last 30 days     │
│ ┌────────────────┐  ┌──────────────────────────────────┐  │
│ │  Total Products│  │ Latest Products                  │  │
│ │       42       │  │ Name: Widget A · Status: active  │  │
│ │  ▁▂▃▆▇▆▇▆▇▆▇▆ │  │ Name: Widget B · Status: draft   │  │
│ │ Last 30 days   │  │ ...                              │  │
│ └────────────────┘  └──────────────────────────────────┘  │
│ (one row per resource that opts in)                        │
└────────────────────────────────────────────────────────────┘`}
      />

      <h2>How a request flows</h2>
      <p>
        One endpoint per resource, dispatched server-side. The web
        widget calls{' '}
        <code>GET /api/admin/dashboard/resource-stats/:resource</code>{' '}
        with whichever date params the DateFilter is active for. The
        handler delegates to a service function that switches on the
        resource name and calls a single reflective helper:
      </p>

      <CodeBlock
        language="go"
        filename="apps/api/internal/services/resource_stats_dispatch.go"
        code={`func ComputeResourceStats(db *gorm.DB, resourceName string, filter ResourceStatsFilter) (*ResourceStats, error) {
    if filter.LatestLimit <= 0 {
        filter.LatestLimit = 10
    }
    switch resourceName {
    case "users":
        return reflectiveResourceStats(db, resourceName, &models.User{}, filter)
    case "blogs":
        return reflectiveResourceStats(db, resourceName, &models.Blog{}, filter)
    // grit:resource-stats:dispatch
    default:
        return nil, fmt.Errorf("dashboard stats not registered for %q", resourceName)
    }
}`}
      />
      <p>
        The dispatcher is also the security boundary: only resources
        registered here are reachable. A compromised admin token
        can&apos;t dump arbitrary tables by guessing resource names
        in the URL — anything not in the switch returns a 400.
      </p>

      <TipBox tone="info">
        Each <code>grit generate resource Foo</code> injects a new
        case at the marker. You never edit this file by hand for a
        normal resource; the generator owns it. If you ever want to
        opt a resource <em>out</em> of dashboard widgets, do it on
        the admin side (see &ldquo;Opting out&rdquo; below) rather
        than removing the dispatch case — the case is also used by
        the system-wide stats endpoint.
      </TipBox>

      <h2>Three pieces of data per resource</h2>
      <p>
        The helper computes three things in one round-trip:
      </p>
      <ul>
        <li>
          <strong>Total</strong> — row count within the active date
          range. No range = all-time count.
        </li>
        <li>
          <strong>30-day sparkline</strong> — one bucket per
          calendar day for the last 30 days. Always 30 days
          regardless of the active filter (see next section for
          why).
        </li>
        <li>
          <strong>Latest N</strong> — up to 10 newest rows in the
          active range. Returned via JSON round-trip so{' '}
          <code>json:&quot;-&quot;</code> tags are honoured —{' '}
          <code>PasswordHash</code> on User never reaches the
          response.
        </li>
      </ul>

      <CodeBlock
        language="json"
        code={`{
  "data": {
    "resource": "products",
    "total": 42,
    "series": [
      { "date": "2026-05-27", "count": 0 },
      { "date": "2026-05-28", "count": 1 },
      ... 30 entries ...
    ],
    "latest": [
      { "id": "01H...", "name": "Widget A", "status": "active", "created_at": "..." },
      ...
    ]
  }
}`}
      />

      <h2>Why the sparkline ignores the DateFilter</h2>
      <p>
        The first instinct is &ldquo;everything should scope to the
        filter.&rdquo; The Total + Latest list do — but the sparkline
        is fixed at 30 days on purpose. Under the &ldquo;Today&rdquo;
        preset a filter-scoped sparkline would collapse to a single
        bar with no trend information. The sparkline is meant to
        answer &ldquo;is this resource growing?&rdquo; not
        &ldquo;what&apos;s inside the filter?&rdquo; — those are
        different questions, and one widget shouldn&apos;t pretend to
        answer both.
      </p>

      <h2>Opting out per resource</h2>
      <p>
        Most resources are worth showing on the dashboard. A few
        aren&apos;t — internal-only tables, audit-log rows, sync
        scratch data. Hide a resource&apos;s widgets by setting{' '}
        <code>dashboard.enabled = false</code> in the resource
        definition:
      </p>

      <CodeBlock
        language="ts"
        filename="apps/admin/resources/internal-note.ts"
        code={`import { defineResource } from "@/lib/resource";

export const internalNoteResource = defineResource({
  name: "Internal Note",
  slug: "internal-notes",
  endpoint: "/api/internal-notes",
  icon: "FileText",
  table: { /* ... */ },
  form: { /* ... */ },
  dashboard: { enabled: false },  // hide widgets from main dashboard
});`}
      />
      <p>
        Default is enabled — a newly generated resource gets the
        widgets without any extra config. The flag is{' '}
        <em>opt-out</em>, not opt-in, because the cost of forgetting
        to opt in is higher than the cost of occasionally opting
        out: a brand-new resource not showing up on the dashboard
        looks like a bug, while an extra row on the dashboard is
        easily ignored.
      </p>

      <h2>The Latest table&apos;s column heuristics</h2>
      <p>
        The Latest widget shows at most three columns per row —
        it&apos;s a glance, not a full data table. The column picker
        prefers <code>name</code>, <code>title</code>,{' '}
        <code>subject</code>, <code>email</code>, <code>status</code>,
        and <code>price</code> in that order, then falls back to the
        first non-id text columns from the resource&apos;s{' '}
        <code>table.columns</code> config. Resources with no
        recognisable columns fall through to showing the row id in
        a monospace font.
      </p>

      <TipBox tone="warning">
        Columns are picked from the resource&apos;s declared{' '}
        <code>table.columns</code>, not from whatever the API
        returns. If you add a new field on the Go model and want it
        to surface in the dashboard preview, you also need to add it
        to the <code>table.columns</code> array (or regenerate the
        resource so the column lands automatically).
      </TipBox>

      <h2>Backward compatibility</h2>
      <p>
        Pre-v3.31.44 projects don&apos;t have the new dispatch file,
        the handler, or the marker. The generator detects this and
        prints a one-line warning per generate instead of failing —
        the dashboard widget for that resource will render &ldquo;not
        registered&rdquo; until the file is added. To upgrade an
        existing project:
      </p>
      <ol>
        <li>
          Copy{' '}
          <code>apps/api/internal/services/resource_stats_dispatch.go</code>{' '}
          and <code>apps/api/internal/handlers/resource_stats.go</code>{' '}
          from a fresh scaffold (or the framework repo&apos;s{' '}
          <code>internal/scaffold/api_resource_stats_files.go</code>).
        </li>
        <li>
          Register the handler + route in your routes file:{' '}
          <code>resourceStatsHandler := &amp;handlers.ResourceStatsHandler{`{DB: db}`}</code>{' '}
          and{' '}
          <code>admin.GET(&quot;/admin/dashboard/resource-stats/:resource&quot;, resourceStatsHandler.Get)</code>.
        </li>
        <li>
          Copy{' '}
          <code>apps/admin/components/dashboard/*.tsx</code> and the{' '}
          <code>By resource</code> section in the dashboard page.
        </li>
        <li>
          Re-run <code>grit generate</code> on each existing
          resource (or hand-add cases to{' '}
          <code>ComputeResourceStats</code>).
        </li>
      </ol>

      <KnowledgeCheck
        question="You set the dashboard DateFilter to 'Last 7 days', but the Products sparkline still shows 30 bars. Bug?"
        choices={[
          {
            label: 'Yes — the React Query key isn\'t including the date range',
            feedback:
              'Wrong — the React Query key does include the date params, which is why the Total and Latest list refetch. The sparkline\'s window is set on the server side and ignored from the filter on purpose.',
          },
          {
            label: 'No — the sparkline is intentionally always 30 days regardless of the filter, to keep the trend shape stable across all presets',
            correct: true,
            feedback:
              'Correct. Under "Today" or "1 day" the filter-scoped sparkline would be a single bar with no trend information. The widget answers two questions: "how many in range?" (Total) and "is this resource growing?" (sparkline). The filter controls the first; the second is fixed.',
          },
          {
            label: 'Yes — buildDailySeries should accept the filter',
            feedback:
              'Wrong — buildDailySeries deliberately uses its own 30-day cutoff and doesn\'t take a filter. Hooking the filter in would defeat the design.',
          },
          {
            label: 'The browser is showing a stale cached chart',
            feedback:
              'Wrong — React Query refetches the whole stat bundle on filter change. If the sparkline data was older it would show; here it\'s simply not scoped by the filter.',
          },
        ]}
      />

      <Exercise
        prompt={
          <>
            <p>
              In the project from the last exercise, see the new
              widgets end-to-end:
            </p>
            <ol>
              <li>
                Open the admin dashboard. Scroll past Quick Access —
                you should see the <em>By resource</em> band with one
                row per resource. Hover the sparkline; the tooltip
                should show &ldquo;N new&rdquo; per day.
              </li>
              <li>
                Toggle the dashboard DateFilter to &quot;Last 7
                days&quot;. The Total numbers should update; the
                sparkline should not. Confirm in the Network tab:{' '}
                <code>resource-stats/products?created_since=7d</code>{' '}
                fires per resource.
              </li>
              <li>
                In the admin, create one new Category. Wait the React
                Query refetch interval (60 s) or hard-refresh. Both
                the Total and the Latest list for Category should
                update.
              </li>
              <li>
                Add{' '}
                <code>{`dashboard: { enabled: false }`}</code> to one
                resource definition. Reload the dashboard — that
                resource&apos;s row should disappear from the band.
              </li>
            </ol>
            <p>
              Screenshot the dashboard with the &quot;By resource&quot;
              band visible and paste it into <code>notes.md</code>.
            </p>
          </>
        }
        hint={
          <>
            The fastest way to see the data shape is to hit the API
            directly with a valid admin cookie:
            <CodeBlock
              language="bash"
              code={`curl -s --cookie "grit_access=<jwt>" \\
  "http://localhost:8080/api/admin/dashboard/resource-stats/products?created_since=7d" | jq`}
            />
            Compare the <code>total</code> to what you see on the
            stat card. If they differ, the React Query cache is
            stale (give it 60 s) or the date range param isn&apos;t
            being forwarded.
          </>
        }
        solution={
          <>
            <p>
              The widget request keys are{' '}
              <code>{`["dashboard", "resource-stats", slug, params]`}</code>{' '}
              and{' '}
              <code>{`["dashboard", "resource-latest", slug, params]`}</code>.
              Changing the DateFilter rewrites <code>params</code>{' '}
              (via <code>dateRangeToQueryParams</code>) and both
              keys invalidate, triggering a refetch in lockstep.
            </p>
            <p>
              Opting out by setting{' '}
              <code>{`dashboard: { enabled: false }`}</code> filters
              the resource out of the band at render time:{' '}
              <code>{`resources.filter((r) => r.dashboard?.enabled !== false)`}</code>.
              The data endpoint still works (the server-side
              dispatch case is untouched) — only the widget is
              hidden from the dashboard.
            </p>
          </>
        }
      />

      <h2>What&apos;s next</h2>
      <p>
        The dashboard now scales linearly with your resource count —
        every new <code>grit generate</code> adds one row to the{' '}
        <em>By resource</em> band automatically. The next lesson
        ({' '}
        <code>web-nextjs/admin-panel/define-resource</code>) goes
        deeper on the resource definition itself, including the{' '}
        <code>table</code> + <code>form</code> halves and how they
        feed back into both the resource page and the dashboard
        widgets you just saw.
      </p>
    </>
  )
}
