import Link from 'next/link'
import { ArrowRight, ArrowLeft } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'
import { CodeBlock } from '@/components/code-block'
import { LaneFlow } from '@/components/lane-flow'
import { getDocMetadata } from '@/config/docs-metadata'

export const metadata = getDocMetadata('/docs/admin/resources')

export default function ResourceDefinitionsPage() {
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
                Resource Definitions
              </h1>
              <p className="text-lg text-muted-foreground leading-relaxed">
                Resources are the building blocks of the Grit admin panel. Define your data
                model once in TypeScript and get a complete admin interface &mdash; data table,
                forms, filters, sidebar navigation, and dashboard widgets.
              </p>
              <LaneFlow
                id="admin-resource"
                lanes={['Resource definition', 'Powers']}
                nodes={[
                  { id: 'res', lane: 0, row: 1, title: 'defineResource', sub: 'columns · fields · filters', tone: 'primary' },
                  { id: 'table', lane: 1, row: 0, title: 'DataTable', sub: 'from columns', tone: 'cyan' },
                  { id: 'form', lane: 1, row: 1, title: 'FormBuilder', sub: 'from fields', tone: 'blue' },
                  { id: 'filter', lane: 1, row: 2, title: 'Filter bar', sub: 'from filters', tone: 'amber' },
                  { id: 'nav', lane: 1, row: 3, title: 'Sidebar', sub: 'from label + icon', tone: 'green' },
                ]}
                edges={[
                  { from: 'res', to: 'table', tone: 'cyan' },
                  { from: 'res', to: 'form', label: 'feeds', tone: 'blue' },
                  { from: 'res', to: 'filter', tone: 'amber' },
                  { from: 'res', to: 'nav', tone: 'green' },
                ]}
                legend={[
                  { tone: 'primary', label: 'One definition object' },
                  { tone: 'cyan', label: 'Every admin surface' },
                ]}
                caption="Each part of the definition object feeds a different admin surface"
              />
            </div>

            <div className="prose-grit">
              {/* defineResource API */}
              <h2>The defineResource() API</h2>
              <p>
                Every admin resource is created with the <code>defineResource()</code> function.
                It accepts a single configuration object that describes the resource name,
                API endpoint, table columns, form fields, and optional dashboard widgets.
              </p>
            </div>

            {/* Type reference */}
            <div className="mt-4 mb-8">
              <CodeBlock filename="ResourceConfig type" code={`interface ResourceConfig {
  // Identity
  name: string              // Singular name, PascalCase (e.g. "Invoice")
  slug?: string             // URL slug, auto-derived if omitted (e.g. "invoices")
  endpoint: string          // Go API base URL (e.g. "/api/invoices")
  icon: string              // Lucide icon name (e.g. "FileText")
  label?: {                 // Display labels (auto-derived from name)
    singular: string        // "Invoice"
    plural: string          // "Invoices"
  }

  // How the create/edit form is presented (optional)
  formView?: 'sheet' | 'modal' | 'page' | 'modal-steps' | 'page-steps'

  // Table configuration
  table: TableConfig

  // Form configuration
  form: FormConfig

  // Dashboard widgets (optional)
  dashboard?: DashboardConfig

  // Stats cards above the table (optional; omit for 4 auto defaults, false to disable)
  stats?: StatsConfig | boolean

  // Sidebar nav grouping — resources sharing a key collapse under one header
  group?: string

  // Access control (optional)
  adminOnly?: boolean       // Hide from the sidebar for non-ADMIN/EDITOR users
}`} />
            </div>

            <div className="prose-grit">
              <h2>Resource Configuration</h2>

              <h3>Name, Slug, and Endpoint</h3>
              <p>
                The <code>name</code> field is the singular PascalCase name of your resource
                (e.g. <code>&quot;Invoice&quot;</code>). Grit auto-derives the plural form,
                URL slug, and display labels from it. You can override any of these with
                the <code>slug</code> and <code>label</code> fields.
              </p>
              <p>
                The <code>endpoint</code> must match the API route registered in your Go
                backend. For a resource named <code>&quot;Invoice&quot;</code>, the endpoint is
                typically <code>/api/invoices</code>. The admin panel appends <code>/:id</code> for
                single-item operations automatically.
              </p>

              <h3>Icon</h3>
              <p>
                Pass any <a href="https://lucide.dev/icons" target="_blank" rel="noreferrer">Lucide icon</a> name
                as a string. The sidebar renders it next to the resource label. Common
                choices: <code>&quot;Users&quot;</code>, <code>&quot;FileText&quot;</code>, <code>&quot;ShoppingCart&quot;</code>, <code>&quot;CreditCard&quot;</code>, <code>&quot;Mail&quot;</code>.
              </p>

              {/* Table Configuration */}
              <h2>Table Configuration</h2>
              <p>
                The <code>table</code> object controls how data appears in the resource&apos;s
                data table &mdash; which columns are shown, how they are formatted, what
                filters are available, and which row actions are enabled.
              </p>
            </div>

            {/* Table config type */}
            <div className="mt-4 mb-8">
              <CodeBlock filename="TableConfig type" code={`interface TableConfig {
  columns: ColumnDef[]
  filters?: FilterDef[]
  pageSize?: number            // Default: 20
  defaultSort?: {
    key: string
    direction: 'asc' | 'desc'
  }
  searchable?: boolean         // Enable global search (default: true)
  actions?: Action[]           // 'create' | 'edit' | 'delete' | 'view' | 'export'
  bulkActions?: BulkAction[]   // 'delete' | 'export'
  dateFilter?: {               // Date-window filter on the list page (default on)
    enabled?: boolean
    field?: string             // Column to filter (default: 'created_at')
    label?: string
  }
  export?: false | {           // Toolbar download menu (CSV / JSON / Excel, default all on)
    csv?: boolean
    json?: boolean
    excel?: boolean
  }
  import?: false | {           // Excel import button + modal (default on)
    excel?: boolean
    fields?: string[]
  }
}`} />
            </div>

            <div className="prose-grit">
              <h3>Column Definitions</h3>
              <p>
                Each column maps a field from the API response to a table column. The
                column definition controls sorting, searching, formatting, and custom
                rendering.
              </p>
            </div>

            {/* Column type */}
            <div className="mt-4 mb-8">
              <CodeBlock filename="ColumnDef type" code={`interface ColumnDef {
  key: string                  // JSON field name (supports dot notation: "customer.name")
  label: string                // Column header text
  sortable?: boolean           // Allow sorting by this column
  searchable?: boolean         // Include in global search
  format?: ColumnFormat        // Display format
  badge?: BadgeConfig          // Badge-style rendering for status fields
  hidden?: boolean             // Hidden by default (show/hide toggle)
  width?: string               // Fixed column width
  className?: string           // Tailwind classes applied to every cell
  cell?: (row) => ReactNode    // Custom renderer (overrides format/badge)
}

type ColumnFormat =
  | 'text'                     // Plain text (default)
  | 'currency'                 // Currency ($1,234.00)
  | 'boolean'                  // Check/X icon
  | 'date'                     // Formatted date (Jan 15, 2026)
  | 'relative'                 // Relative time (3 hours ago)
  | 'badge'                    // Colored badge
  | 'image'                    // Thumbnail image
  | 'video' | 'file' | 'files' // Media / attachment previews
  | 'link' | 'email' | 'color' // URL, mailto, and color swatch
  | 'richtext'                 // HTML stripped to a plain-text preview
  | 'user'                     // Avatar + name + email stacked

interface BadgeConfig {
  [value: string]: {
    color: 'green' | 'yellow' | 'red' | 'blue' | 'purple' | 'gray'
    label: string
  }
}`} />
            </div>

            <div className="prose-grit">
              <h3>Column Format Types</h3>
              <p>
                The <code>format</code> property determines how the cell value is rendered:
              </p>
            </div>

            {/* Format types table */}
            <div className="mt-4 mb-8">
              <div className="rounded-lg border border-border/30 bg-card/30 overflow-hidden">
                <table className="w-full text-sm">
                  <thead>
                    <tr className="border-b border-border/30 bg-accent/20">
                      <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Format</th>
                      <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Input</th>
                      <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Rendered As</th>
                    </tr>
                  </thead>
                  <tbody className="text-muted-foreground">
                    {[
                      ['text', '"Hello World"', 'Hello World'],
                      ['currency', '99.5', '$99.50'],
                      ['boolean', 'true / false', 'Green check / Red X icon'],
                      ['date', '"2026-01-15T..."', 'Jan 15, 2026'],
                      ['relative', '"2026-01-15T..."', '3 weeks ago'],
                      ['badge', '"active"', 'Colored pill with label'],
                      ['image', '"https://..."', '32x32 rounded thumbnail'],
                    ].map(([format, input, output]) => (
                      <tr key={format} className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-mono text-xs text-primary">{format}</td>
                        <td className="px-4 py-2.5 font-mono text-xs">{input}</td>
                        <td className="px-4 py-2.5">{output}</td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            </div>

            <div className="prose-grit">
              <h3>Badge Columns</h3>
              <p>
                For status-like fields, use the <code>badge</code> property instead of
                (or in addition to) <code>format</code>. It maps each possible value to a
                colored label:
              </p>
            </div>

            <div className="mt-4 mb-8">
              <CodeBlock filename="Badge column example" code={`{
  key: 'status',
  label: 'Status',
  badge: {
    paid:    { color: 'green',  label: 'Paid' },
    pending: { color: 'yellow', label: 'Pending' },
    overdue: { color: 'red',    label: 'Overdue' },
  },
}`} />
            </div>

            <div className="prose-grit">
              <h3>Filters</h3>
              <p>
                Filters appear above the data table and let users narrow down results.
                Three filter types are available:
              </p>
            </div>

            <div className="mt-4 mb-8">
              <CodeBlock filename="FilterDef type" code={`interface FilterDef {
  key: string                  // Field to filter on
  type: 'select' | 'date-range' | 'number-range'
  label?: string               // Display label (defaults to key)
  options?: string[]           // For 'select' type
}

// Example filters
filters: [
  { key: 'status', type: 'select', options: ['paid', 'pending', 'overdue'] },
  { key: 'created_at', type: 'date-range' },
  { key: 'amount', type: 'number-range' },
]`} />
            </div>

            <div className="prose-grit">
              {/* Form Configuration */}
              <h2>Form Configuration</h2>
              <p>
                The <code>form</code> object defines the fields that appear in create and
                edit modals. Grit supports a wide range of field types and validates input
                using Zod schemas from the shared package.
              </p>
            </div>

            <div className="mt-4 mb-8">
              <CodeBlock filename="FormConfig type" code={`interface FormConfig {
  fields: FieldDef[]
  layout?: 'single' | 'two-column'  // Default: 'single'
}

interface FieldDef {
  key: string                  // JSON field name
  label: string                // Display label
  type: FieldType              // Input type
  required?: boolean           // Required field (default: false)
  placeholder?: string         // Placeholder text
  defaultValue?: unknown       // Default value for create mode
  options?: { label: string; value: string }[]
  min?: number                 // For number type
  max?: number                 // For number type
  step?: number                // For number type
  rows?: number                // For textarea type
  colSpan?: 1 | 2              // Column span in two-column layout
}

type FieldType =
  | 'text'
  | 'textarea'
  | 'number'
  | 'select'
  | 'date'
  | 'datetime'
  | 'toggle'
  | 'checkbox'
  | 'radio'
  | 'image'                    // Single image upload
  | 'images'                   // Multiple image upload
  | 'video'                    // Single video upload
  | 'videos'                   // Multiple video upload
  | 'file'                     // Single file upload
  | 'files'                    // Multiple file upload
  | 'richtext'                 // Rich text editor
  | 'relationship-select'      // Select from related resource
  | 'multi-relationship-select' // Multi-select from related resource`} />
            </div>

            <div className="prose-grit">
              <h3>Form Field Types</h3>
              <p>
                Each field type renders a different input component. Here is a quick reference:
              </p>
            </div>

            <div className="mt-4 mb-8">
              <div className="rounded-lg border border-border/30 bg-card/30 overflow-hidden">
                <table className="w-full text-sm">
                  <thead>
                    <tr className="border-b border-border/30 bg-accent/20">
                      <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Type</th>
                      <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Component</th>
                      <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Extra Props</th>
                    </tr>
                  </thead>
                  <tbody className="text-muted-foreground">
                    {[
                      ['text', 'Text input', 'placeholder'],
                      ['textarea', 'Multi-line textarea', 'rows, placeholder'],
                      ['number', 'Numeric input', 'min, max, step'],
                      ['select', 'Dropdown select', 'options'],
                      ['date', 'Date picker', '--'],
                      ['datetime', 'Date & time picker', '--'],
                      ['toggle', 'Toggle switch', '--'],
                      ['checkbox', 'Checkbox', '--'],
                      ['radio', 'Radio button group', 'options'],
                      ['image', 'Single image upload', 'accept'],
                      ['images', 'Multiple image upload', 'accept'],
                      ['video', 'Single video upload', 'accept'],
                      ['videos', 'Multiple video upload', 'accept'],
                      ['file', 'Single file upload', 'accept'],
                      ['files', 'Multiple file upload', 'accept'],
                      ['richtext', 'Rich text editor', '--'],
                      ['relationship-select', 'Resource select', 'relatedEndpoint, displayField'],
                      ['multi-relationship-select', 'Resource multi-select', 'relatedEndpoint, displayField, relationshipKey'],
                    ].map(([type, component, props]) => (
                      <tr key={type} className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-mono text-xs text-primary">{type}</td>
                        <td className="px-4 py-2.5">{component}</td>
                        <td className="px-4 py-2.5 font-mono text-xs">{props}</td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            </div>

            <div className="prose-grit">
              {/* Dashboard Configuration */}
              <h2>Dashboard Configuration</h2>
              <p>
                Each resource can optionally define widgets that appear on the admin
                dashboard. Widgets pull data from your Go API and display stats, charts,
                or activity feeds.
              </p>
            </div>

            <div className="mt-4 mb-8">
              <CodeBlock filename="DashboardConfig type" code={`interface DashboardConfig {
  enabled?: boolean            // false hides this resource's preset widgets
  widgets?: WidgetDef[]
}

interface WidgetDef {
  type: 'stat' | 'chart' | 'activity'
  label: string
  endpoint?: string            // Go API URL the widget reads (e.g. "/api/orders?page_size=1")
  icon?: string                // Lucide icon name
  color?: string               // Accent color for the widget
  format?: 'number' | 'currency' | 'percentage'
  chartType?: 'line' | 'bar' | 'pie'   // For type: 'chart'
  colSpan?: 1 | 2 | 3 | 4
}`} />
            </div>

            <div className="prose-grit">
              {/* Complete Example */}
              <h2>Complete Example</h2>
              <p>
                Here is a full resource definition for a <strong>Posts</strong> resource with
                table columns, filters, form fields, and dashboard widgets:
              </p>
            </div>

            {/* Full example */}
            <div className="mt-4 mb-8">
              <CodeBlock language="typescript" filename="apps/admin/resources/posts.ts" code={`import { defineResource } from '@/lib/resource'

export const postsResource = defineResource({
  name: 'Post',
  slug: 'posts',
  endpoint: '/api/posts',
  icon: 'FileText',
  label: {
    singular: 'Blog Post',
    plural: 'Blog Posts',
  },
  adminOnly: true,

  table: {
    columns: [
      { key: 'title', label: 'Title', sortable: true, searchable: true },
      { key: 'author.name', label: 'Author' },
      { key: 'category', label: 'Category', sortable: true },
      { key: 'status', label: 'Status', badge: {
        published: { color: 'green', label: 'Published' },
        draft:     { color: 'yellow', label: 'Draft' },
        archived:  { color: 'gray', label: 'Archived' },
      }},
      { key: 'views', label: 'Views', sortable: true },
      { key: 'published_at', label: 'Published', format: 'date' },
      { key: 'created_at', label: 'Created', format: 'relative' },
    ],
    filters: [
      { key: 'status', type: 'select', options: ['published', 'draft', 'archived'] },
      { key: 'category', type: 'select', options: ['tech', 'design', 'business'] },
      { key: 'published_at', type: 'date-range' },
    ],
    pageSize: 25,
    defaultSort: { key: 'created_at', direction: 'desc' },
    actions: ['create', 'edit', 'delete', 'export'],
    bulkActions: ['delete', 'export'],
    // Toolbar download menu (CSV / JSON / Excel) + Excel import are on by default.
  },

  form: {
    layout: 'two-column',
    fields: [
      { key: 'title', label: 'Title', type: 'text', required: true,
        placeholder: 'Enter post title', colSpan: 2 },
      { key: 'slug', label: 'Slug', type: 'text',
        placeholder: 'auto-generated-from-title' },
      { key: 'category', label: 'Category', type: 'select',
        options: [
          { label: 'Tech', value: 'tech' },
          { label: 'Design', value: 'design' },
          { label: 'Business', value: 'business' },
        ] },
      { key: 'content', label: 'Content', type: 'richtext', colSpan: 2 },
      { key: 'excerpt', label: 'Excerpt', type: 'textarea', rows: 3,
        colSpan: 2 },
      { key: 'status', label: 'Status', type: 'select',
        options: [
          { label: 'Published', value: 'published' },
          { label: 'Draft', value: 'draft' },
          { label: 'Archived', value: 'archived' },
        ], defaultValue: 'draft' },
      { key: 'featured', label: 'Featured Post', type: 'toggle' },
      { key: 'cover_image', label: 'Cover Image', type: 'file', colSpan: 2 },
    ],
  },

  dashboard: {
    widgets: [
      { type: 'stat', label: 'Total Posts', icon: 'FileText', color: 'accent',
        endpoint: '/api/posts?page_size=1', format: 'number' },
      { type: 'stat', label: 'Published', icon: 'CheckCircle', color: 'success',
        endpoint: '/api/posts?status=published&page_size=1', format: 'number' },
    ],
  },
})`} />
            </div>

            <div className="prose-grit">
              {/* Resource Registry */}
              <h2>Resource Registry</h2>
              <p>
                After creating a resource definition file, you need to register it in the
                resource registry. This file is the single source of truth for all admin
                resources &mdash; the sidebar, router, and dashboard all read from it.
              </p>
            </div>

            <div className="mt-4 mb-8">
              <CodeBlock language="typescript" filename="apps/admin/resources/index.ts" code={`import users from './users'
import posts from './posts'
import invoices from './invoices'

// All registered resources — sidebar and routes are generated from this array
export const resources = [users, posts, invoices]

// Helper to look up a resource by slug
export function getResource(slug: string) {
  return resources.find((r) => r.slug === slug)
}`} />
            </div>

            <div className="prose-grit">
              <p>
                When you run <code>grit generate resource</code>, the CLI automatically
                adds the import and registration to this file using marker-based code
                injection. You never need to edit it manually unless you want to reorder
                the sidebar items.
              </p>

              <h2>How Resources Auto-Register in the Sidebar</h2>
              <p>
                The admin sidebar component imports the <code>resources</code> array from the
                registry and renders a navigation link for each one. Each resource&apos;s
                <code>icon</code> and <code>label.plural</code> (or auto-derived name) are
                used as the sidebar text. The active item is highlighted based on the
                current URL path.
              </p>
              <p>
                The order of resources in the registry array determines their order in the
                sidebar. System pages (Dashboard, Jobs, Files, Settings) are rendered
                separately and always appear at fixed positions.
              </p>
            </div>

            {/* Nav */}
            <div className="flex items-center justify-between pt-6 mt-10 border-t border-border/30">
              <Button variant="ghost" size="sm" asChild className="text-muted-foreground/60 hover:text-foreground">
                <Link href="/docs/admin/overview" className="gap-1.5">
                  <ArrowLeft className="h-3.5 w-3.5" />
                  Admin Overview
                </Link>
              </Button>
              <Button variant="ghost" size="sm" asChild className="text-muted-foreground/60 hover:text-foreground">
                <Link href="/docs/admin/datatable" className="gap-1.5">
                  DataTable
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
