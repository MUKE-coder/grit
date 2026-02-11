import Link from 'next/link'
import { ArrowRight, ArrowLeft } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'

export default function DashboardWidgetsPage() {
  return (
    <div className="min-h-screen bg-background">
      <SiteHeader />
      <DocsSidebar />

      <main className="lg:pl-64">
        <div className="container max-w-screen-xl py-10 px-6">
          <div className="max-w-3xl">
            {/* Header */}
            <div className="mb-10">
              <span className="tag-mono text-primary/80 mb-3 block">Admin Panel</span>
              <h1 className="text-4xl font-bold tracking-tight mb-4">
                Dashboard & Widgets
              </h1>
              <p className="text-lg text-muted-foreground leading-relaxed">
                The admin dashboard is the home page of the admin panel. It displays a
                collection of widgets &mdash; stats cards, charts, and activity feeds &mdash;
                assembled from your resource definitions and custom API endpoints.
              </p>
            </div>

            <div className="prose-grit">
              {/* Dashboard Page */}
              <h2>Dashboard Page</h2>
              <p>
                When an admin user opens the admin panel, the first page they see is the
                dashboard (<code>apps/admin/app/page.tsx</code>). It aggregates widgets from
                all registered resources that define a <code>dashboard</code> section, plus
                any global widgets you configure.
              </p>
              <p>
                The dashboard layout uses a responsive grid that adapts to the screen size:
              </p>
              <ul>
                <li><strong>Desktop (lg+)</strong> &mdash; 4-column grid for stats cards, 2-column grid for charts.</li>
                <li><strong>Tablet (md)</strong> &mdash; 2-column grid for stats, full-width charts.</li>
                <li><strong>Mobile (sm)</strong> &mdash; single column, all widgets stacked vertically.</li>
              </ul>
              <p>
                Widgets load their data independently using React Query, so the dashboard
                renders progressively &mdash; fast widgets appear immediately while slower
                ones show skeleton loaders.
              </p>

              {/* Widget Types */}
              <h2>Widget Types</h2>

              {/* Stats Card */}
              <h3>Stats Card</h3>
              <p>
                A compact card that displays a single metric. Stats cards are the most common
                widget type and typically appear in a row of 3-4 across the top of the
                dashboard.
              </p>
            </div>

            {/* Stats card config */}
            <div className="mt-4 mb-8">
              <div className="rounded-xl border border-border/40 bg-card/80 overflow-hidden">
                <div className="flex items-center gap-2 px-4 py-2.5 border-b border-border/30 bg-accent/30">
                  <span className="text-[11px] font-mono text-muted-foreground/40">StatsCard properties</span>
                </div>
                <pre className="p-5 text-sm font-mono text-foreground/80 overflow-x-auto">{`interface StatsCardWidget {
  type: 'stat'
  label: string              // Display label (e.g. "Total Revenue")
  query: string              // Server query (e.g. "sum:amount")
  format?: 'number' | 'currency' | 'percent'
  color?: string             // Accent color: 'purple' | 'green' | 'blue' | 'yellow' | 'red'
  icon?: string              // Lucide icon name
  changeQuery?: string       // Query for period-over-period change %
}`}</pre>
              </div>
            </div>

            <div className="prose-grit">
              <p>
                Each stats card renders:
              </p>
              <ul>
                <li><strong>Value</strong> &mdash; the main metric, formatted according to the <code>format</code> property (plain number, currency with <code>$</code> prefix, or percentage).</li>
                <li><strong>Label</strong> &mdash; descriptive text below the value.</li>
                <li><strong>Change indicator</strong> &mdash; if <code>changeQuery</code> is set, a green upward arrow or red downward arrow with the percentage change compared to the previous period.</li>
                <li><strong>Icon</strong> &mdash; a Lucide icon on the right side of the card, rendered in the accent color.</li>
                <li><strong>Color</strong> &mdash; a subtle colored accent on the left border or background of the card.</li>
              </ul>
            </div>

            <div className="mt-4 mb-8">
              <div className="rounded-xl border border-border/40 bg-card/80 overflow-hidden glow-purple-sm">
                <div className="flex items-center gap-2 px-4 py-2.5 border-b border-border/30 bg-accent/30">
                  <span className="text-[11px] font-mono text-muted-foreground/40">Stats card examples</span>
                </div>
                <pre className="p-5 text-sm font-mono text-foreground/80 overflow-x-auto">{`dashboard: {
  widgets: [
    {
      type: 'stat',
      label: 'Total Revenue',
      query: 'sum:amount',
      format: 'currency',
      icon: 'DollarSign',
      color: 'green',
      changeQuery: 'sum:amount:change:month',
    },
    {
      type: 'stat',
      label: 'Total Users',
      query: 'count',
      format: 'number',
      icon: 'Users',
      color: 'blue',
    },
    {
      type: 'stat',
      label: 'Pending Orders',
      query: 'count:status=pending',
      format: 'number',
      icon: 'Clock',
      color: 'yellow',
    },
    {
      type: 'stat',
      label: 'Conversion Rate',
      query: 'custom:conversion-rate',
      format: 'percent',
      icon: 'TrendingUp',
      color: 'purple',
    },
  ],
}`}</pre>
              </div>
            </div>

            <div className="prose-grit">
              {/* Line Chart */}
              <h3>Line Chart</h3>
              <p>
                Line charts display time-series data using <a href="https://recharts.org" target="_blank" rel="noreferrer">Recharts</a>.
                They are ideal for showing trends over time &mdash; revenue per month, new
                users per week, or page views per day.
              </p>
            </div>

            <div className="mt-4 mb-8">
              <div className="rounded-xl border border-border/40 bg-card/80 overflow-hidden">
                <div className="flex items-center gap-2 px-4 py-2.5 border-b border-border/30 bg-accent/30">
                  <span className="text-[11px] font-mono text-muted-foreground/40">LineChart properties</span>
                </div>
                <pre className="p-5 text-sm font-mono text-foreground/80 overflow-x-auto">{`interface LineChartWidget {
  type: 'chart'
  chartType: 'line'
  label: string              // Chart title
  query: string              // e.g. "sum:amount:by:month" or "count:by:week"
  format?: 'number' | 'currency'
  color?: string             // Line color
  height?: number            // Chart height in px (default: 300)
}`}</pre>
              </div>
            </div>

            <div className="prose-grit">
              <p>
                The chart component renders a smooth line with a gradient fill area below it.
                Hovering over a data point shows a tooltip with the exact value and date.
                The X-axis shows time labels (months, weeks, or days depending on the query
                granularity) and the Y-axis auto-scales to fit the data range.
              </p>
            </div>

            <div className="mt-4 mb-8">
              <div className="rounded-xl border border-border/40 bg-card/80 overflow-hidden">
                <div className="flex items-center gap-2 px-4 py-2.5 border-b border-border/30 bg-accent/30">
                  <span className="text-[11px] font-mono text-muted-foreground/40">Line chart example</span>
                </div>
                <pre className="p-5 text-sm font-mono text-foreground/80 overflow-x-auto">{`{
  type: 'chart',
  chartType: 'line',
  label: 'Revenue Over Time',
  query: 'sum:amount:by:month',
  format: 'currency',
  color: 'purple',
  height: 320,
}`}</pre>
              </div>
            </div>

            <div className="prose-grit">
              {/* Bar Chart */}
              <h3>Bar Chart</h3>
              <p>
                Bar charts display categorical data &mdash; comparisons between categories,
                status breakdowns, or grouped counts. They use the same Recharts library and
                share the same tooltip and axis styling as line charts.
              </p>
            </div>

            <div className="mt-4 mb-8">
              <div className="rounded-xl border border-border/40 bg-card/80 overflow-hidden">
                <div className="flex items-center gap-2 px-4 py-2.5 border-b border-border/30 bg-accent/30">
                  <span className="text-[11px] font-mono text-muted-foreground/40">BarChart properties</span>
                </div>
                <pre className="p-5 text-sm font-mono text-foreground/80 overflow-x-auto">{`interface BarChartWidget {
  type: 'chart'
  chartType: 'bar'
  label: string              // Chart title
  query: string              // e.g. "count:by:category" or "sum:amount:by:status"
  format?: 'number' | 'currency'
  color?: string             // Bar color
  height?: number            // Chart height in px (default: 300)
}`}</pre>
              </div>
            </div>

            <div className="mt-4 mb-8">
              <div className="rounded-xl border border-border/40 bg-card/80 overflow-hidden">
                <div className="flex items-center gap-2 px-4 py-2.5 border-b border-border/30 bg-accent/30">
                  <span className="text-[11px] font-mono text-muted-foreground/40">Bar chart example</span>
                </div>
                <pre className="p-5 text-sm font-mono text-foreground/80 overflow-x-auto">{`{
  type: 'chart',
  chartType: 'bar',
  label: 'Posts by Category',
  query: 'count:by:category',
  format: 'number',
  color: 'blue',
}`}</pre>
              </div>
            </div>

            <div className="prose-grit">
              {/* Recent Activity */}
              <h3>Recent Activity</h3>
              <p>
                The activity widget displays a chronological list of recent events &mdash; new
                records created, records updated, users logged in, jobs completed, etc. Each
                event shows a timestamp, an icon, a description, and optionally a link to
                the related resource.
              </p>
            </div>

            <div className="mt-4 mb-8">
              <div className="rounded-xl border border-border/40 bg-card/80 overflow-hidden">
                <div className="flex items-center gap-2 px-4 py-2.5 border-b border-border/30 bg-accent/30">
                  <span className="text-[11px] font-mono text-muted-foreground/40">Activity widget properties</span>
                </div>
                <pre className="p-5 text-sm font-mono text-foreground/80 overflow-x-auto">{`interface ActivityWidget {
  type: 'activity'
  label: string              // Widget title (e.g. "Recent Activity")
  query: string              // e.g. "recent:10" (last 10 events)
  height?: number            // Max height in px (scrollable)
}

// API response format for activity events
interface ActivityEvent {
  id: string
  type: 'created' | 'updated' | 'deleted' | 'login' | 'custom'
  resource: string           // Resource name (e.g. "Post")
  description: string        // "John created a new post"
  user?: { name: string; avatar?: string }
  timestamp: string          // ISO 8601
  link?: string              // Optional link to the resource
}`}</pre>
              </div>
            </div>

            <div className="mt-4 mb-8">
              <div className="rounded-xl border border-border/40 bg-card/80 overflow-hidden">
                <div className="flex items-center gap-2 px-4 py-2.5 border-b border-border/30 bg-accent/30">
                  <span className="text-[11px] font-mono text-muted-foreground/40">Activity widget example</span>
                </div>
                <pre className="p-5 text-sm font-mono text-foreground/80 overflow-x-auto">{`{
  type: 'activity',
  label: 'Recent Activity',
  query: 'recent:15',
  height: 400,
}`}</pre>
              </div>
            </div>

            <div className="prose-grit">
              {/* Widget Grid Layout */}
              <h2>Widget Grid Layout</h2>
              <p>
                The dashboard arranges widgets in a responsive grid. The layout logic follows
                these rules:
              </p>
              <ul>
                <li><strong>Stats cards</strong> are grouped together and rendered in a 4-column row (lg), 2-column (md), or 1-column (sm).</li>
                <li><strong>Charts</strong> take half the grid width on desktop (2 charts side by side) and full width on smaller screens.</li>
                <li><strong>Activity widgets</strong> take half the grid width on desktop and full width on smaller screens.</li>
              </ul>
              <p>
                Widgets from different resources are merged and grouped by type. All stats
                cards from all resources appear in one row at the top, followed by charts,
                then activity feeds. This creates a cohesive dashboard rather than
                resource-isolated sections.
              </p>

              {/* Data Fetching */}
              <h2>Widget Data Fetching</h2>
              <p>
                Each widget fetches its data independently from the Go API using React Query.
                The <code>query</code> string in the widget definition is translated to an
                API call:
              </p>
            </div>

            <div className="mt-4 mb-8">
              <div className="rounded-lg border border-border/30 bg-card/30 overflow-hidden">
                <table className="w-full text-sm">
                  <thead>
                    <tr className="border-b border-border/30 bg-accent/20">
                      <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Widget Query</th>
                      <th className="text-left px-4 py-2.5 font-medium text-foreground/80">API Call</th>
                    </tr>
                  </thead>
                  <tbody className="text-muted-foreground">
                    {[
                      ['count', 'GET /api/posts/stats?metric=count'],
                      ['sum:amount', 'GET /api/invoices/stats?metric=sum&field=amount'],
                      ['count:status=pending', 'GET /api/orders/stats?metric=count&status=pending'],
                      ['sum:amount:by:month', 'GET /api/invoices/stats?metric=sum&field=amount&group=month'],
                      ['count:by:category', 'GET /api/posts/stats?metric=count&group=category'],
                      ['recent:10', 'GET /api/activity?limit=10'],
                      ['custom:conversion-rate', 'GET /api/stats/conversion-rate'],
                    ].map(([query, api]) => (
                      <tr key={query} className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-mono text-xs text-primary">{query}</td>
                        <td className="px-4 py-2.5 font-mono text-xs">{api}</td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            </div>

            <div className="prose-grit">
              <p>
                Widget data is cached by React Query with a default stale time of 30 seconds.
                Stats cards poll for fresh data every 60 seconds in the background so the
                dashboard stays reasonably current without excessive API calls.
              </p>

              {/* Custom Widget Endpoints */}
              <h2>Custom Widget Endpoints</h2>
              <p>
                For metrics that cannot be expressed as simple aggregations, use the
                <code>custom:</code> prefix in the query string. This tells the widget to call a
                custom API endpoint that you implement in Go:
              </p>
            </div>

            <div className="mt-4 mb-8">
              <div className="rounded-xl border border-border/40 bg-card/80 overflow-hidden">
                <div className="flex items-center gap-2 px-4 py-2.5 border-b border-border/30 bg-accent/30">
                  <span className="text-[11px] font-mono text-muted-foreground/40">apps/api/internal/handlers/stats.go</span>
                </div>
                <pre className="p-5 text-sm font-mono text-foreground/80 overflow-x-auto">{`// Custom endpoint for conversion rate widget
func (h *StatsHandler) GetConversionRate(c *gin.Context) {
    totalVisitors, err := h.service.CountVisitors(c)
    if err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }

    totalSignups, err := h.service.CountSignups(c)
    if err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }

    rate := float64(totalSignups) / float64(totalVisitors) * 100

    c.JSON(200, gin.H{
        "data": gin.H{
            "value": rate,
            "label": "Conversion Rate",
        },
    })
}`}</pre>
              </div>
            </div>

            <div className="prose-grit">
              <p>
                Register the custom endpoint in your routes file:
              </p>
            </div>

            <div className="mt-4 mb-8">
              <div className="rounded-xl border border-border/40 bg-card/80 overflow-hidden">
                <div className="flex items-center gap-2 px-4 py-2.5 border-b border-border/30 bg-accent/30">
                  <span className="text-[11px] font-mono text-muted-foreground/40">apps/api/internal/routes/routes.go</span>
                </div>
                <pre className="p-5 text-sm font-mono text-foreground/80 overflow-x-auto">{`// Stats endpoints
stats := api.Group("/stats")
stats.Use(middleware.AuthMiddleware(), middleware.RequireRole("admin"))
{
    stats.GET("/conversion-rate", statsHandler.GetConversionRate)
    stats.GET("/monthly-mrr", statsHandler.GetMonthlyMRR)
}`}</pre>
              </div>
            </div>

            <div className="prose-grit">
              {/* Adding Widgets to Resources */}
              <h2>Adding Widgets to Resource Definitions</h2>
              <p>
                Each resource can define its own widgets in the <code>dashboard</code> section.
                These widgets are automatically included on the admin dashboard alongside
                widgets from other resources.
              </p>
            </div>

            <div className="mt-4 mb-8">
              <div className="rounded-xl border border-border/40 bg-card/80 overflow-hidden glow-purple-sm">
                <div className="flex items-center gap-2 px-4 py-2.5 border-b border-border/30 bg-accent/30">
                  <span className="text-[11px] font-mono text-muted-foreground/40">apps/admin/resources/orders.ts</span>
                </div>
                <pre className="p-5 text-sm font-mono text-foreground/80 overflow-x-auto">{`import { defineResource } from '@grit/admin'

export default defineResource({
  name: 'Order',
  endpoint: '/api/orders',
  icon: 'ShoppingCart',

  table: { /* ... columns and filters ... */ },
  form: { /* ... fields ... */ },

  dashboard: {
    widgets: [
      // Stats cards
      {
        type: 'stat',
        label: 'Total Orders',
        query: 'count',
        icon: 'ShoppingCart',
        color: 'purple',
      },
      {
        type: 'stat',
        label: 'Revenue',
        query: 'sum:total',
        format: 'currency',
        icon: 'DollarSign',
        color: 'green',
        changeQuery: 'sum:total:change:month',
      },
      {
        type: 'stat',
        label: 'Pending Orders',
        query: 'count:status=pending',
        icon: 'Clock',
        color: 'yellow',
      },

      // Charts
      {
        type: 'chart',
        chartType: 'line',
        label: 'Revenue Over Time',
        query: 'sum:total:by:month',
        format: 'currency',
        color: 'purple',
      },
      {
        type: 'chart',
        chartType: 'bar',
        label: 'Orders by Status',
        query: 'count:by:status',
        color: 'blue',
      },

      // Activity feed
      {
        type: 'activity',
        label: 'Recent Orders',
        query: 'recent:10',
      },
    ],
  },
})`}</pre>
              </div>
            </div>

            <div className="prose-grit">
              {/* API Response Format */}
              <h2>Widget API Response Format</h2>
              <p>
                The Go API must return widget data in these formats:
              </p>

              <h3>Stats Response</h3>
            </div>

            <div className="mt-4 mb-8">
              <div className="rounded-xl border border-border/40 bg-card/80 overflow-hidden">
                <div className="flex items-center gap-2 px-4 py-2.5 border-b border-border/30 bg-accent/30">
                  <span className="text-[11px] font-mono text-muted-foreground/40">Stats API response</span>
                </div>
                <pre className="p-5 text-sm font-mono text-foreground/80 overflow-x-auto">{`// GET /api/orders/stats?metric=count
{
  "data": {
    "value": 1247,
    "change": 12.5,         // +12.5% vs previous period (optional)
    "previous": 1108        // previous period value (optional)
  }
}

// GET /api/orders/stats?metric=sum&field=total
{
  "data": {
    "value": 84350.00,
    "change": -3.2
  }
}`}</pre>
              </div>
            </div>

            <div className="prose-grit">
              <h3>Chart Response</h3>
            </div>

            <div className="mt-4 mb-8">
              <div className="rounded-xl border border-border/40 bg-card/80 overflow-hidden">
                <div className="flex items-center gap-2 px-4 py-2.5 border-b border-border/30 bg-accent/30">
                  <span className="text-[11px] font-mono text-muted-foreground/40">Chart API response</span>
                </div>
                <pre className="p-5 text-sm font-mono text-foreground/80 overflow-x-auto">{`// GET /api/orders/stats?metric=sum&field=total&group=month
{
  "data": [
    { "label": "Sep 2025", "value": 12400 },
    { "label": "Oct 2025", "value": 15800 },
    { "label": "Nov 2025", "value": 13200 },
    { "label": "Dec 2025", "value": 19500 },
    { "label": "Jan 2026", "value": 17800 },
    { "label": "Feb 2026", "value": 21300 }
  ]
}

// GET /api/posts/stats?metric=count&group=category
{
  "data": [
    { "label": "Tech",     "value": 42 },
    { "label": "Design",   "value": 28 },
    { "label": "Business", "value": 15 }
  ]
}`}</pre>
              </div>
            </div>

            <div className="prose-grit">
              <h3>Activity Response</h3>
            </div>

            <div className="mt-4 mb-8">
              <div className="rounded-xl border border-border/40 bg-card/80 overflow-hidden">
                <div className="flex items-center gap-2 px-4 py-2.5 border-b border-border/30 bg-accent/30">
                  <span className="text-[11px] font-mono text-muted-foreground/40">Activity API response</span>
                </div>
                <pre className="p-5 text-sm font-mono text-foreground/80 overflow-x-auto">{`// GET /api/activity?limit=10
{
  "data": [
    {
      "id": "act_001",
      "type": "created",
      "resource": "Order",
      "description": "New order #1247 placed by John Doe",
      "user": { "name": "John Doe", "avatar": "https://..." },
      "timestamp": "2026-02-11T14:30:00Z",
      "link": "/resources/orders/1247"
    },
    {
      "id": "act_002",
      "type": "updated",
      "resource": "Invoice",
      "description": "Invoice INV-089 marked as paid",
      "user": { "name": "Admin" },
      "timestamp": "2026-02-11T14:15:00Z"
    }
  ]
}`}</pre>
              </div>
            </div>

            <div className="prose-grit">
              {/* Styling */}
              <h2>Widget Styling</h2>
              <p>
                All widgets follow the Grit dark theme aesthetic. Cards have subtle borders
                (<code>border-border/40</code>), slightly elevated backgrounds
                (<code>bg-card/80</code>), and consistent padding. Charts use the purple accent
                color by default with gradient fills. Stats cards have a thin colored left
                border matching their <code>color</code> property.
              </p>
              <p>
                Skeleton loaders match the exact dimensions of each widget type, preventing
                layout shift during initial load. Error states display a subtle error message
                inside the widget area without breaking the grid layout.
              </p>
            </div>

            {/* What's next */}
            <div className="mt-10 mb-10">
              <h2 className="text-2xl font-semibold tracking-tight mb-4">
                What&apos;s Next?
              </h2>
              <div className="space-y-3">
                {[
                  { label: 'Web App (Frontend)', href: '/docs/frontend/web-app', desc: 'Learn about the Next.js frontend application and its structure.' },
                  { label: 'React Query Hooks', href: '/docs/frontend/hooks', desc: 'Auto-generated data fetching hooks for your resources.' },
                  { label: 'File Storage', href: '/docs/batteries/storage', desc: 'S3-compatible file uploads and management.' },
                  { label: 'Background Jobs', href: '/docs/batteries/jobs', desc: 'Redis-backed job queue with admin dashboard.' },
                ].map((item) => (
                  <Link
                    key={item.href}
                    href={item.href}
                    className="flex items-center justify-between p-4 rounded-lg border border-border/30 bg-card/30 hover:bg-card/60 hover:border-primary/20 transition-all group"
                  >
                    <div>
                      <p className="text-sm font-medium group-hover:text-primary transition-colors">{item.label}</p>
                      <p className="text-xs text-muted-foreground/60 mt-0.5">{item.desc}</p>
                    </div>
                    <ArrowRight className="h-4 w-4 text-muted-foreground/40 group-hover:text-primary transition-colors" />
                  </Link>
                ))}
              </div>
            </div>

            {/* Nav */}
            <div className="flex items-center justify-between pt-6 border-t border-border/30">
              <Button variant="ghost" size="sm" asChild className="text-muted-foreground/60 hover:text-foreground">
                <Link href="/docs/admin/forms" className="gap-1.5">
                  <ArrowLeft className="h-3.5 w-3.5" />
                  Form Builder
                </Link>
              </Button>
              <Button variant="ghost" size="sm" asChild className="text-muted-foreground/60 hover:text-foreground">
                <Link href="/docs/frontend/web-app" className="gap-1.5">
                  Web App
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
