import Link from 'next/link'
import { ArrowRight, ArrowLeft } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'
import { CodeBlock } from '@/components/code-block'

export default function DataTablePage() {
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
                DataTable
              </h1>
              <p className="text-lg text-muted-foreground leading-relaxed">
                The DataTable component is the primary way data is displayed in the admin panel.
                It supports server-side pagination, column sorting, filtering, search, custom
                cell renderers, row actions, and data export &mdash; all driven by your resource
                definition.
              </p>
            </div>

            <div className="prose-grit">
              {/* Server-side Pagination */}
              <h2>Server-Side Pagination</h2>
              <p>
                The DataTable never loads the entire dataset into memory. It communicates with
                your Go API using query parameters for page, page size, sort, and filter values.
                The API returns a paginated response with a <code>meta</code> object containing
                total count, current page, page size, and total pages.
              </p>
              <p>
                Pagination controls appear at the bottom of the table showing the current range
                (e.g. &quot;Showing 1-20 of 156&quot;), page navigation buttons, and a page size
                selector (10, 20, 50, 100 rows).
              </p>
            </div>

            {/* API request example */}
            <div className="mt-4 mb-8">
              <CodeBlock filename="HTTP request" code={`GET /api/posts?page=1&page_size=20&sort=created_at&order=desc&search=hello

Response:
{
  "data": [ ... ],
  "meta": {
    "total": 156,
    "page": 1,
    "page_size": 20,
    "pages": 8
  }
}`} />
            </div>

            <div className="prose-grit">
              <p>
                The default page size is controlled by <code>table.pageSize</code> in your
                resource definition (default: 20). Users can change the page size at runtime
                using the selector in the pagination bar.
              </p>

              {/* Sorting */}
              <h2>Column Sorting</h2>
              <p>
                Any column with <code>sortable: true</code> in its definition becomes
                clickable. Clicking a column header cycles through three states:
              </p>
              <ol>
                <li><strong>Ascending</strong> &mdash; small arrow pointing up appears next to the header.</li>
                <li><strong>Descending</strong> &mdash; arrow points down.</li>
                <li><strong>No sort</strong> &mdash; returns to the default sort order.</li>
              </ol>
              <p>
                Sorting sends <code>sort</code> and <code>order</code> query parameters to the
                API. Only single-column sorting is supported (clicking a new column clears the
                previous sort). You can set a default sort in the resource definition:
              </p>
            </div>

            <div className="mt-4 mb-8">
              <CodeBlock filename="Default sort" code={`table: {
  defaultSort: { key: 'created_at', direction: 'desc' },
  columns: [ ... ],
}`} />
            </div>

            <div className="prose-grit">
              {/* Filtering */}
              <h2>Column Filtering</h2>
              <p>
                Filters appear as a horizontal bar above the data table. Three filter types
                are supported, each rendering a different control:
              </p>

              <h3>Select Filter</h3>
              <p>
                Renders a dropdown with predefined options. Useful for status fields,
                categories, or any column with a fixed set of values.
              </p>
            </div>

            <div className="mt-4 mb-8">
              <CodeBlock filename="Select filter" code={`{ key: 'status', type: 'select', options: ['active', 'inactive', 'banned'] }`} />
            </div>

            <div className="prose-grit">
              <h3>Date Range Filter</h3>
              <p>
                Renders two date pickers (from and to) for filtering records within a time
                window. Sends <code>created_at_from</code> and <code>created_at_to</code> query
                parameters.
              </p>
            </div>

            <div className="mt-4 mb-8">
              <CodeBlock filename="Date range filter" code={`{ key: 'created_at', type: 'date-range', label: 'Created Date' }`} />
            </div>

            <div className="prose-grit">
              <h3>Number Range Filter</h3>
              <p>
                Renders two number inputs (min and max) for filtering numeric values. Useful
                for price ranges, quantities, or scores. Sends <code>amount_min</code> and
                <code>amount_max</code> query parameters.
              </p>
            </div>

            <div className="mt-4 mb-8">
              <CodeBlock filename="Number range filter" code={`{ key: 'amount', type: 'number-range', label: 'Amount' }`} />
            </div>

            <div className="prose-grit">
              <p>
                Active filters show a count badge on the filter bar and a &quot;Clear
                filters&quot; button appears when any filter is applied. Changing filters resets
                the page to 1.
              </p>

              {/* Show/Hide Columns */}
              <h2>Show/Hide Columns</h2>
              <p>
                A column visibility dropdown appears in the table toolbar. Users can toggle
                individual columns on or off. Columns with <code>hidden: true</code> in their
                definition are hidden by default but can be shown via the dropdown. Column
                visibility preferences are persisted in <code>localStorage</code> so they
                survive page reloads.
              </p>

              {/* Search */}
              <h2>Search</h2>
              <p>
                When <code>table.searchable</code> is <code>true</code> (the default), a search
                input appears in the table toolbar. Typing into it sends a <code>search</code> query
                parameter to the API. On the Go side, the search handler applies a
                <code>ILIKE</code> query across all columns marked with <code>searchable: true</code> in
                the resource definition.
              </p>
              <p>
                Search is debounced at 300ms to avoid excessive API calls while the user is
                typing.
              </p>

              {/* Custom Cell Renderers */}
              <h2>Custom Cell Renderers</h2>
              <p>
                The <code>format</code> and <code>badge</code> column options cover most use
                cases. Here are examples of each renderer:
              </p>
            </div>

            {/* Cell renderer examples */}
            <div className="mt-4 mb-8">
              <CodeBlock filename="Column renderer examples" code={`columns: [
  // Badge — colored pills for status values
  { key: 'status', label: 'Status', badge: {
    active:   { color: 'green', label: 'Active' },
    inactive: { color: 'gray',  label: 'Inactive' },
  }},

  // Currency — formatted as $1,234.50
  { key: 'amount', label: 'Amount', format: 'currency' },

  // Date — formatted as "Jan 15, 2026"
  { key: 'due_date', label: 'Due Date', format: 'date' },

  // Relative — formatted as "3 hours ago"
  { key: 'created_at', label: 'Created', format: 'relative' },

  // Boolean — green checkmark or red X
  { key: 'active', label: 'Active', format: 'boolean' },

  // Image — 32x32 rounded thumbnail
  { key: 'avatar', label: 'Avatar', format: 'image' },

  // Number — formatted with thousand separators
  { key: 'views', label: 'Views', format: 'number' },

  // Relation — displays a field from a related object
  { key: 'customer.name', label: 'Customer', relation: 'customer' },

  // Video — thumbnail with play overlay
  { key: 'preview', label: 'Preview', format: 'video' },

  // Link — clickable URL with hostname
  { key: 'website', label: 'Website', format: 'link' },

  // Email — clickable mailto link
  { key: 'email', label: 'Email', format: 'email' },

  // Color — swatch circle with hex value
  { key: 'color', label: 'Color', format: 'color' },
]`} />
            </div>

            <div className="prose-grit">
              <h2>Column Styling</h2>
              <p>
                Add the <code>className</code> property to any column definition to apply
                custom Tailwind CSS classes to every cell in that column. This wraps the
                rendered content in a <code>&lt;span&gt;</code> with your classes, so it
                works alongside any format type.
              </p>
            </div>

            <div className="mt-4 mb-8">
              <CodeBlock filename="Column className examples" code={`columns: [
  // Bold title column
  { key: 'title', label: 'Title', className: 'font-semibold text-foreground' },

  // Green currency column
  { key: 'price', label: 'Price', format: 'currency', className: 'text-success' },

  // Monospace code column
  { key: 'sku', label: 'SKU', className: 'font-mono text-xs tracking-wider' },

  // Truncated long text
  { key: 'description', label: 'Description', className: 'max-w-[200px] truncate' },
]`} />
            </div>

            <div className="prose-grit">
              {/* Row Actions */}
              <h2>Row Actions</h2>
              <p>
                Each row has an actions menu (three-dot icon) on the right side. The available
                actions are controlled by the <code>table.actions</code> array in your resource
                definition:
              </p>
              <ul>
                <li><strong>create</strong> &mdash; adds a &quot;New [Resource]&quot; button to the table toolbar that opens the create form modal.</li>
                <li><strong>edit</strong> &mdash; opens the edit form modal pre-filled with the row data.</li>
                <li><strong>delete</strong> &mdash; shows a confirmation dialog, then sends a DELETE request to the API.</li>
                <li><strong>view</strong> &mdash; navigates to a detail page for the resource.</li>
                <li><strong>export</strong> &mdash; adds CSV/JSON export buttons to the toolbar.</li>
              </ul>
              <p>
                Delete actions use optimistic updates via React Query &mdash; the row is
                removed from the table immediately and restored if the API call fails.
              </p>
            </div>

            <div className="mt-4 mb-8">
              <CodeBlock filename="Actions configuration" code={`table: {
  actions: ['create', 'edit', 'delete', 'export'],
  bulkActions: ['delete', 'export'],
  columns: [ ... ],
}`} />
            </div>

            <div className="prose-grit">
              {/* Bulk Actions */}
              <h2>Bulk Actions</h2>
              <p>
                When <code>bulkActions</code> are defined, each row gets a checkbox on the
                left side. Selecting one or more rows reveals a floating action bar at the
                bottom of the table with the configured bulk actions. For example, selecting
                5 rows and clicking &quot;Delete&quot; sends 5 DELETE requests in parallel.
              </p>

              {/* Empty State */}
              <h2>Empty State</h2>
              <p>
                When a resource has no data (or no results match the current filters), the
                DataTable shows a polished empty state with:
              </p>
              <ul>
                <li>An illustration matching the resource icon</li>
                <li>A message like &quot;No posts yet&quot;</li>
                <li>A &quot;Create your first post&quot; button that opens the create form modal</li>
              </ul>
              <p>
                The empty state maintains the table&apos;s full width and height so the layout
                does not collapse.
              </p>

              {/* Loading State */}
              <h2>Loading Skeleton</h2>
              <p>
                While data is being fetched, the DataTable renders a skeleton loader that
                matches the exact layout of the table &mdash; header row, column widths, and
                row heights are preserved. This prevents layout shift when data loads and
                gives users confidence that content is on its way.
              </p>
              <p>
                Subsequent page navigations (changing page, applying filters) show a subtle
                loading indicator in the table header instead of replacing the entire table
                with a skeleton. This keeps the current data visible while the new data loads.
              </p>

              {/* Export */}
              <h2>Export to CSV / JSON</h2>
              <p>
                When <code>&apos;export&apos;</code> is included in <code>table.actions</code>,
                export buttons appear in the table toolbar. Users can export the current
                filtered and sorted view as CSV or JSON. The export fetches all matching
                records from the API (not just the current page) and triggers a browser
                download.
              </p>
              <p>
                Column labels are used as CSV headers. Hidden columns are excluded from the
                export unless explicitly shown. Badge values export as their raw value (e.g.
                <code>&quot;paid&quot;</code>) rather than the display label.
              </p>

              {/* Responsive Behavior */}
              <h2>Responsive Behavior</h2>
              <p>
                On screens narrower than the table&apos;s natural width, the DataTable enables
                horizontal scrolling. The row actions column is sticky on the right side so
                it remains visible while scrolling. On mobile devices, the filter bar collapses
                into a &quot;Filters&quot; button that opens a sheet overlay.
              </p>
            </div>

            {/* Full working example */}
            <div className="mt-8 mb-8">
              <h2 className="text-2xl font-semibold tracking-tight mb-4">
                Full Table Configuration Example
              </h2>
              <CodeBlock filename="apps/admin/resources/invoices.ts (table section)" code={`table: {
  columns: [
    { key: 'id', label: 'ID', sortable: true, hidden: true },
    { key: 'number', label: 'Invoice #', sortable: true, searchable: true },
    { key: 'customer.name', label: 'Customer', relation: 'customer',
      searchable: true },
    { key: 'amount', label: 'Amount', format: 'currency', sortable: true },
    { key: 'status', label: 'Status', badge: {
      paid:    { color: 'green',  label: 'Paid' },
      pending: { color: 'yellow', label: 'Pending' },
      overdue: { color: 'red',    label: 'Overdue' },
    }},
    { key: 'due_date', label: 'Due Date', format: 'date', sortable: true },
    { key: 'created_at', label: 'Created', format: 'relative' },
  ],
  filters: [
    { key: 'status', type: 'select',
      options: ['paid', 'pending', 'overdue'] },
    { key: 'created_at', type: 'date-range' },
    { key: 'amount', type: 'number-range' },
  ],
  pageSize: 20,
  defaultSort: { key: 'created_at', direction: 'desc' },
  searchable: true,
  actions: ['create', 'edit', 'delete', 'view', 'export'],
  bulkActions: ['delete', 'export', 'mark-paid'],
}`} />
            </div>

            {/* Nav */}
            <div className="flex items-center justify-between pt-6 border-t border-border/30">
              <Button variant="ghost" size="sm" asChild className="text-muted-foreground/60 hover:text-foreground">
                <Link href="/docs/admin/resources" className="gap-1.5">
                  <ArrowLeft className="h-3.5 w-3.5" />
                  Resource Definitions
                </Link>
              </Button>
              <Button variant="ghost" size="sm" asChild className="text-muted-foreground/60 hover:text-foreground">
                <Link href="/docs/admin/forms" className="gap-1.5">
                  Form Builder
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
