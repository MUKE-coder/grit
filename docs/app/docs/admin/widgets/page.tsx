import Link from 'next/link'
import { ArrowRight, ArrowLeft } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'
import { CodeBlock } from '@/components/code-block'
import { getDocMetadata } from '@/config/docs-metadata'

export const metadata = getDocMetadata('/docs/admin/widgets')

export default function DashboardWidgetsPage() {
  return (
    <div className="min-h-screen bg-background isolate">
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
                all registered resources &mdash; both the preset widgets every resource
                gets automatically and any custom widgets you declare in the resource&apos;s
                <code>dashboard</code> section.
              </p>
              <p>
                The dashboard layout uses a responsive grid. Each widget claims a number of
                columns through its <code>colSpan</code> property (1&ndash;4), and the grid
                collapses to fewer columns on smaller screens:
              </p>
              <ul>
                <li><strong>Desktop (lg+)</strong> &mdash; 4-column grid; a <code>colSpan: 2</code> widget takes half the row.</li>
                <li><strong>Tablet (md)</strong> &mdash; 2-column grid; wide widgets span the full width.</li>
                <li><strong>Mobile (sm)</strong> &mdash; single column, all widgets stacked vertically.</li>
              </ul>
              <p>
                Widgets load their data independently using React Query, so the dashboard
                renders progressively &mdash; fast widgets appear immediately while slower
                ones show skeleton loaders.
              </p>

              {/* Preset widgets */}
              <h2>Preset Widgets (Opt-Out)</h2>
              <p>
                Every generated resource <strong>automatically gets a set of preset dashboard
                widgets</strong> &mdash; a <em>Total</em> stat with a sparkline and a
                <em>Latest N</em> activity list. You do not need to configure anything to get
                them. They are opt-<strong>out</strong>: to hide a resource&apos;s preset
                widgets from the dashboard, set <code>enabled: false</code> on its
                <code>dashboard</code> definition.
              </p>
            </div>

            <div className="mt-4 mb-8">
              <CodeBlock language="typescript" filename="Hide the preset widgets" code={`dashboard: {
  enabled: false,   // hides the preset Total + Latest N widgets for this resource
}`} />
            </div>

            <div className="prose-grit">
              <p>
                Declaring <code>widgets</code> does not disable the presets &mdash; your custom
                widgets render <em>alongside</em> the presets unless you also set
                <code>enabled: false</code>.
              </p>

              {/* Widget Types */}
              <h2>Widget Types</h2>
              <p>
                Every dashboard widget shares a single shape, <code>WidgetDefinition</code>.
                The <code>type</code> field selects the kind of widget, and for charts the
                <code>chartType</code> field selects the visualization. Import the types (and
                <code>defineResource</code>) from <code>@/lib/resource</code>.
              </p>
            </div>

            <div className="mt-4 mb-8">
              <CodeBlock language="typescript" filename="apps/admin/lib/resource.ts" code={`export type WidgetType = "stat" | "chart" | "activity";
export type ChartType = "line" | "bar" | "pie";
export type WidgetFormat = "number" | "currency" | "percentage";

export interface WidgetDefinition {
  type: WidgetType;
  label: string;
  endpoint?: string;         // where the widget fetches its data
  icon?: string;             // Lucide icon name
  color?: string;            // accent color
  format?: WidgetFormat;     // how "stat" values are formatted
  chartType?: ChartType;     // only for type "chart"
  limit?: number;            // e.g. how many rows for an "activity" widget
  colSpan?: 1 | 2 | 3 | 4;   // grid width
}

export interface DashboardDefinition {
  enabled?: boolean;              // false hides the preset per-resource widgets
  widgets?: WidgetDefinition[];   // custom widgets
}`} />
            </div>

            <div className="prose-grit">
              <p>
                The three widget types, and the <code>chartType</code> matrix for charts:
              </p>
            </div>

            <div className="mt-4 mb-8">
              <div className="rounded-lg border border-border/30 bg-card/30 overflow-hidden">
                <table className="w-full text-sm">
                  <thead>
                    <tr className="border-b border-border/30 bg-accent/20">
                      <th className="text-left px-4 py-2.5 font-medium text-foreground/80">type</th>
                      <th className="text-left px-4 py-2.5 font-medium text-foreground/80">chartType</th>
                      <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Renders</th>
                    </tr>
                  </thead>
                  <tbody className="text-muted-foreground">
                    {[
                      ['stat', '—', 'A single metric value (formatted by format), icon, and color.'],
                      ['chart', 'line', 'A time-series line chart.'],
                      ['chart', 'bar', 'A categorical bar chart.'],
                      ['chart', 'pie', 'A proportional pie chart.'],
                      ['activity', '—', 'A list of the latest limit records/events.'],
                    ].map(([type, chart, renders]) => (
                      <tr key={`${type}-${chart}`} className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-mono text-xs text-primary">{type}</td>
                        <td className="px-4 py-2.5 font-mono text-xs">{chart}</td>
                        <td className="px-4 py-2.5 text-xs">{renders}</td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            </div>

            <div className="prose-grit">
              {/* Stat widget */}
              <h3>Stat Widget</h3>
              <p>
                A compact card that displays a single metric fetched from its
                <code>endpoint</code>. The <code>format</code> property controls how the value
                is rendered: <code>&quot;number&quot;</code>, <code>&quot;currency&quot;</code>,
                or <code>&quot;percentage&quot;</code>.
              </p>
            </div>

            <div className="mt-4 mb-8">
              <CodeBlock language="typescript" filename="Stat widget" code={`{
  type: 'stat',
  label: 'Total Revenue',
  endpoint: '/api/orders/stats/revenue',
  format: 'currency',
  icon: 'DollarSign',
  color: 'green',
  colSpan: 1,
}`} />
            </div>

            <div className="prose-grit">
              {/* Chart widget */}
              <h3>Chart Widget</h3>
              <p>
                Charts render with <a href="https://recharts.org" target="_blank" rel="noreferrer">Recharts</a>.
                Set <code>type: &apos;chart&apos;</code> and pick a <code>chartType</code> of
                <code>&apos;line&apos;</code> (trends over time), <code>&apos;bar&apos;</code>
                (categorical comparisons), or <code>&apos;pie&apos;</code> (proportions of a
                whole). The widget fetches an array of data points from its
                <code>endpoint</code>.
              </p>
            </div>

            <div className="mt-4 mb-8">
              <CodeBlock language="typescript" filename="Chart widgets" code={`// Line chart — revenue over time
{
  type: 'chart',
  chartType: 'line',
  label: 'Revenue Over Time',
  endpoint: '/api/orders/stats/revenue-by-month',
  format: 'currency',
  color: 'purple',
  colSpan: 2,
}

// Bar chart — orders grouped by status
{
  type: 'chart',
  chartType: 'bar',
  label: 'Orders by Status',
  endpoint: '/api/orders/stats/by-status',
  color: 'blue',
  colSpan: 2,
}

// Pie chart — share of orders per category
{
  type: 'chart',
  chartType: 'pie',
  label: 'Orders by Category',
  endpoint: '/api/orders/stats/by-category',
  colSpan: 2,
}`} />
            </div>

            <div className="prose-grit">
              {/* Activity widget */}
              <h3>Activity Widget</h3>
              <p>
                The activity widget displays a chronological list of the most recent records
                or events. Use <code>limit</code> to control how many rows it shows; the
                widget requests them from its <code>endpoint</code>.
              </p>
            </div>

            <div className="mt-4 mb-8">
              <CodeBlock language="typescript" filename="Activity widget" code={`{
  type: 'activity',
  label: 'Recent Orders',
  endpoint: '/api/orders?sort=-created_at',
  limit: 10,
  colSpan: 2,
}`} />
            </div>

            <div className="prose-grit">
              {/* Grid Layout */}
              <h2>Grid Layout</h2>
              <p>
                The dashboard is a 4-column grid on desktop. Each widget&apos;s
                <code>colSpan</code> (1&ndash;4) decides how many columns it occupies; widgets
                flow left-to-right and wrap onto the next row when the current one fills. A
                typical layout &mdash; a row of four <code>colSpan: 1</code> stats, then two
                <code>colSpan: 2</code> charts, then a full-width activity feed &mdash; looks
                like this:
              </p>
            </div>

            <div className="mt-4 mb-8">
              <CodeBlock language="text" filename="Dashboard grid (colSpan)" code={`┌──────────┬──────────┬──────────┬──────────┐
│  stat    │  stat    │  stat    │  stat    │   colSpan: 1  ×4
│  (1)     │  (1)     │  (1)     │  (1)     │
├──────────┴──────────┼──────────┴──────────┤
│  chart · line       │  chart · bar        │   colSpan: 2  ×2
│  (2)                │  (2)                │
├─────────────────────┴─────────────────────┤
│  activity · Recent Orders                 │   colSpan: 4
│  (4)                                      │
└───────────────────────────────────────────┘`} />
            </div>

            <div className="prose-grit">
              <p>
                On tablet the grid collapses to 2 columns and on mobile to a single column, so
                a <code>colSpan: 2</code> widget becomes full-width and everything stacks.
              </p>

              {/* Custom Endpoints */}
              <h2>Custom Widget Endpoints</h2>
              <p>
                Each widget names an <code>endpoint</code>, and the admin fetches that URL
                directly with React Query &mdash; there is no query DSL or translation layer.
                Your Go handler decides what the widget shows; return the payload under a
                <code>data</code> key following the standard Grit response format.
              </p>
              <p>
                A <code>stat</code> widget expects a single value, a <code>chart</code> widget
                expects an array of <code>{'{ label, value }'}</code> points, and an
                <code>activity</code> widget expects an array of records.
              </p>
            </div>

            <div className="mt-4 mb-8">
              <CodeBlock filename="apps/api/internal/handlers/stats.go" code={`// GET /api/orders/stats/revenue  →  feeds a { type: 'stat' } widget
func (h *StatsHandler) GetRevenue(c *gin.Context) {
    total, err := h.service.SumOrderTotals(c)
    if err != nil {
        c.JSON(500, gin.H{"error": gin.H{"message": err.Error()}})
        return
    }

    c.JSON(200, gin.H{
        "data": total, // e.g. 84350.00 — rendered with format: 'currency'
    })
}

// GET /api/orders/stats/revenue-by-month  →  feeds a { type: 'chart' } widget
func (h *StatsHandler) GetRevenueByMonth(c *gin.Context) {
    points, err := h.service.RevenueByMonth(c)
    if err != nil {
        c.JSON(500, gin.H{"error": gin.H{"message": err.Error()}})
        return
    }

    // points: []gin.H{{"label": "Sep 2025", "value": 12400}, ...}
    c.JSON(200, gin.H{"data": points})
}`} />
            </div>

            <div className="prose-grit">
              <p>
                Register the endpoints in your routes file:
              </p>
            </div>

            <div className="mt-4 mb-8">
              <CodeBlock filename="apps/api/internal/routes/routes.go" code={`// Order stats endpoints
orders := api.Group("/orders/stats")
orders.Use(middleware.AuthMiddleware(), middleware.RequireRole("admin"))
{
    orders.GET("/revenue", statsHandler.GetRevenue)
    orders.GET("/revenue-by-month", statsHandler.GetRevenueByMonth)
    orders.GET("/by-status", statsHandler.GetOrdersByStatus)
}`} />
            </div>

            <div className="prose-grit">
              {/* Stat cards above tables */}
              <h2>Stat Cards Above the Table</h2>
              <p>
                Separate from dashboard widgets, every resource page can show a row of
                <strong>stat cards above its data table</strong>. These are configured with the
                resource&apos;s <code>stats</code> property, which accepts either a boolean or a
                <code>StatsConfig</code> object:
              </p>
              <ul>
                <li><strong>Omit <code>stats</code></strong> &mdash; you get 4 auto-generated cards (Total, This Week, This Month, Updated Recently).</li>
                <li><strong><code>stats: false</code></strong> &mdash; disables the stat cards for this resource page.</li>
                <li><strong><code>stats: {'{ cards: [...] }'}</code></strong> &mdash; fully custom cards.</li>
              </ul>
            </div>

            <div className="mt-4 mb-8">
              <CodeBlock language="typescript" filename="apps/admin/lib/resource.ts" code={`export interface StatsConfig {
  enabled?: boolean;
  cards?: StatCardConfig[];
}

export interface StatCardConfig {
  label: string;
  icon?: string;
  color?: "default" | "success" | "warning" | "danger" | "info";
  value?: string | number;         // a fixed value, or…
  endpoint?: string;               // …fetch the value from here
  field?: string;                  // which field in the response to read
  trend?: { value: number; direction: "up" | "down" };
}`} />
            </div>

            <div className="mt-4 mb-8">
              <CodeBlock language="typescript" filename="Custom stat cards" code={`stats: {
  cards: [
    {
      label: 'Total Orders',
      icon: 'ShoppingCart',
      color: 'default',
      endpoint: '/api/orders/stats/count',
      field: 'value',
    },
    {
      label: 'Revenue',
      icon: 'DollarSign',
      color: 'success',
      endpoint: '/api/orders/stats/revenue',
      field: 'value',
      trend: { value: 12.5, direction: 'up' },
    },
    {
      label: 'Pending',
      icon: 'Clock',
      color: 'warning',
      value: 18,
    },
  ],
}`} />
            </div>

            <div className="prose-grit">
              {/* Full example */}
              <h2>Full Resource Example</h2>
              <p>
                Putting it together &mdash; a resource with custom dashboard widgets and custom
                stat cards. The preset dashboard widgets are hidden here with
                <code>dashboard.enabled: false</code> so only the custom widgets show.
              </p>
            </div>

            <div className="mt-4 mb-8">
              <CodeBlock language="typescript" filename="apps/admin/resources/orders.ts" code={`import { defineResource } from '@/lib/resource'

export default defineResource({
  name: 'Order',
  slug: 'orders',
  endpoint: '/api/orders',
  icon: 'ShoppingCart',

  table: { /* ... columns and filters ... */ },
  form: { /* ... fields ... */ },

  // Custom stat cards above the orders table
  stats: {
    cards: [
      { label: 'Total Orders', icon: 'ShoppingCart', endpoint: '/api/orders/stats/count', field: 'value' },
      { label: 'Revenue', icon: 'DollarSign', color: 'success', endpoint: '/api/orders/stats/revenue', field: 'value', trend: { value: 12.5, direction: 'up' } },
      { label: 'Pending', icon: 'Clock', color: 'warning', value: 18 },
    ],
  },

  // Dashboard widgets (presets hidden via enabled: false)
  dashboard: {
    enabled: false,
    widgets: [
      { type: 'stat', label: 'Total Revenue', endpoint: '/api/orders/stats/revenue', format: 'currency', icon: 'DollarSign', color: 'green', colSpan: 1 },
      { type: 'stat', label: 'Total Orders', endpoint: '/api/orders/stats/count', format: 'number', icon: 'ShoppingCart', color: 'purple', colSpan: 1 },
      { type: 'chart', chartType: 'line', label: 'Revenue Over Time', endpoint: '/api/orders/stats/revenue-by-month', format: 'currency', color: 'purple', colSpan: 2 },
      { type: 'chart', chartType: 'bar', label: 'Orders by Status', endpoint: '/api/orders/stats/by-status', color: 'blue', colSpan: 2 },
      { type: 'activity', label: 'Recent Orders', endpoint: '/api/orders?sort=-created_at', limit: 10, colSpan: 4 },
    ],
  },
})`} />
            </div>

            <div className="prose-grit">
              {/* API Response Format */}
              <h2>Widget API Response Format</h2>
              <p>
                Widget endpoints return their payload under a <code>data</code> key. The shape
                of <code>data</code> depends on the widget type consuming it.
              </p>
            </div>

            <div className="mt-4 mb-8">
              <CodeBlock language="json" filename="Widget API responses" code={`// stat widget  →  GET /api/orders/stats/revenue
{ "data": 84350.00 }

// chart widget →  GET /api/orders/stats/revenue-by-month
{
  "data": [
    { "label": "Sep 2025", "value": 12400 },
    { "label": "Oct 2025", "value": 15800 },
    { "label": "Nov 2025", "value": 13200 },
    { "label": "Dec 2025", "value": 19500 }
  ]
}

// activity widget →  GET /api/orders?sort=-created_at
{
  "data": [
    { "id": 1247, "status": "paid", "total": 129.00, "created_at": "2026-02-11T14:30:00Z" },
    { "id": 1246, "status": "pending", "total": 84.50, "created_at": "2026-02-11T14:15:00Z" }
  ]
}`} />
            </div>

            <div className="prose-grit">
              {/* Styling */}
              <h2>Widget Styling</h2>
              <p>
                All widgets follow the Grit dark theme aesthetic. Cards have subtle borders
                (<code>border-border/40</code>), slightly elevated backgrounds
                (<code>bg-card/80</code>), and consistent padding. Charts use the purple accent
                color by default with gradient fills. Stat widgets and stat cards render their
                icon and accent using the <code>color</code> property.
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
