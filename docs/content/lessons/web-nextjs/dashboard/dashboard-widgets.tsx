import { CodeBlock } from '@/components/code-block'
import { TipBox } from '@/components/course/tip-box'
import { Exercise } from '@/components/course/exercise'
import { KnowledgeCheck } from '@/components/course/knowledge-check'

export default function Lesson() {
  return (
    <>
      <p>
        Logged-in user landing surface. Three sections: stat cards across
        the top, a chart in the middle, an activity feed on the side. The
        Grit scaffold ships these as starter widgets — this lesson covers
        how they wire up.
      </p>

      <h2>The dashboard shell</h2>
      <CodeBlock
        language="tsx"
        filename="apps/web/app/(app)/dashboard/page.tsx"
        code={`import { StatCards } from './stat-cards'
import { RevenueChart } from './revenue-chart'
import { ActivityFeed } from './activity-feed'
import { apiFetch } from '@/lib/api'

export default async function DashboardPage() {
  const stats = await apiFetch('/api/dashboard/stats').then((r) => r.json())

  return (
    <div className="grid grid-cols-1 lg:grid-cols-3 gap-4">
      <div className="lg:col-span-3">
        <StatCards stats={stats.data} />
      </div>
      <div className="lg:col-span-2">
        <RevenueChart />
      </div>
      <div>
        <ActivityFeed />
      </div>
    </div>
  )
}`}
      />
      <p>
        Server component — fetches stats on the server, passes to the
        client widgets that need interactivity (chart hover, feed
        polling).
      </p>

      <h2>Stat cards</h2>
      <CodeBlock
        language="tsx"
        filename="apps/web/app/(app)/dashboard/stat-cards.tsx"
        code={`import { formatCurrency } from '@workspace/shared/utils/currency'

interface Stats {
  users: number
  revenue: number
  orders: number
  conversion: number
}

export function StatCards({ stats }: { stats: Stats }) {
  return (
    <div className="grid grid-cols-2 md:grid-cols-4 gap-3">
      <Card label="Users"       value={stats.users.toLocaleString()} />
      <Card label="Revenue"     value={formatCurrency(stats.revenue)} />
      <Card label="Orders"      value={stats.orders.toLocaleString()} />
      <Card label="Conversion"  value={(stats.conversion * 100).toFixed(1) + '%'} />
    </div>
  )
}

function Card({ label, value }) {
  return (
    <div className="rounded-xl border p-4 bg-card">
      <p className="text-xs text-muted-foreground">{label}</p>
      <p className="text-2xl font-semibold mt-1">{value}</p>
    </div>
  )
}`}
      />
      <p>
        Pure server component. No JS shipped to the browser for the
        cards.
      </p>

      <h2>The chart — Recharts (client component)</h2>
      <CodeBlock
        language="tsx"
        filename="apps/web/app/(app)/dashboard/revenue-chart.tsx"
        code={`'use client'
import { useQuery } from '@tanstack/react-query'
import { LineChart, Line, XAxis, YAxis, Tooltip, ResponsiveContainer } from 'recharts'

export function RevenueChart() {
  const { data } = useQuery({
    queryKey: ['revenue-30d'],
    queryFn: () => fetch('/api/dashboard/revenue/30d').then(r => r.json()),
  })
  return (
    <div className="rounded-xl border p-4 bg-card">
      <h3 className="font-semibold mb-3">Revenue, last 30 days</h3>
      <ResponsiveContainer width="100%" height={240}>
        <LineChart data={data?.data ?? []}>
          <XAxis dataKey="date" />
          <YAxis />
          <Tooltip />
          <Line type="monotone" dataKey="revenue" stroke="#6c5ce7" />
        </LineChart>
      </ResponsiveContainer>
    </div>
  )
}`}
      />
      <p>
        Recharts is the practical pick — small bundle, declarative,
        composable. Charts are inherently interactive so they live in
        client components.
      </p>

      <TipBox tone="info">
        <strong>Split server vs. client by interactivity, not page.</strong>{' '}
        The page itself is server; chart is client; cards are server. Each
        piece ships only as much JS as it needs.
      </TipBox>

      <h2>The activity feed</h2>
      <CodeBlock
        language="tsx"
        filename="apps/web/app/(app)/dashboard/activity-feed.tsx"
        code={`'use client'
import { useQuery } from '@tanstack/react-query'

export function ActivityFeed() {
  const { data } = useQuery({
    queryKey: ['activity'],
    queryFn: () => fetch('/api/activity?limit=10').then(r => r.json()),
    refetchInterval: 30_000,  // poll every 30s for fresh events
  })
  return (
    <div className="rounded-xl border p-4 bg-card">
      <h3 className="font-semibold mb-3">Recent activity</h3>
      <ul className="space-y-2 text-sm">
        {(data?.data ?? []).map((a) => (
          <li key={a.id}>
            <span className="text-muted-foreground">{a.created_at}</span>{' '}
            {a.event}
          </li>
        ))}
      </ul>
    </div>
  )
}`}
      />
      <p>
        Polls every 30 seconds for new events. For real-time, swap polling
        for the Grit WebSocket plugin (covered in the plugins course).
      </p>

      <KnowledgeCheck
        question="Your dashboard initially loads fast but slows over time as users add data. Which widget is the most likely culprit?"
        choices={[
          {
            label: 'Stat cards — they aggregate everything',
            correct: true,
            feedback:
              "Right — `COUNT(*)` over a growing table is the classic slowdown. Fix with materialised views, summary tables updated nightly, or Redis-cached aggregates with periodic refresh.",
          },
          {
            label: 'The chart — Recharts is slow',
            feedback:
              "Recharts handles thousands of points fine. The chart fetches at most 30 days worth — bounded data.",
          },
          {
            label: 'The activity feed — polling',
            feedback:
              "Polling 30s is light. The query is `LIMIT 10` so it stays fast.",
          },
          {
            label: 'React Query caching',
            feedback:
              "React Query speeds things up, doesn't slow them down. Not the culprit.",
          },
        ]}
      />

      <Exercise
        prompt={
          <>
            <p>For chapter 3&apos;s assignment, build the dashboard:</p>
            <ol>
              <li>Sign up a new user (from last lesson).</li>
              <li>Add 3 stat cards backed by real API endpoints.</li>
              <li>Add the activity feed pulling from <code>/api/activity</code>.</li>
              <li>
                Screenshot the dashboard. Paste it in <code>notes.md</code>.
              </li>
            </ol>
          </>
        }
        hint={
          <>
            If <code>/api/dashboard/stats</code> doesn&apos;t exist on your
            scaffolded API yet, build it: a handler that returns{' '}
            <code>{`{ users: int, revenue: decimal, orders: int }`}</code>{' '}
            counted from the DB.
          </>
        }
        solution={
          <>
            <p>
              Should look like:
            </p>
            <CodeBlock
              language="text"
              code={`[ Users: 124 ]  [ Revenue: $84,321 ]  [ Orders: 1,284 ]  [ Conv: 3.2% ]

╭─────────────╮  ╭──────────────╮
│  Revenue 30 │  │ Recent       │
│   chart     │  │ activity     │
╰─────────────╯  ╰──────────────╯`}
            />
          </>
        }
      />

      <h2>What&apos;s next</h2>
      <p>
        Chapter 4 — <strong>The Admin Panel</strong>. Filament-style CRUD
        pages built from one <code>defineResource()</code> call.
      </p>
    </>
  )
}
