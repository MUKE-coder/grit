import Link from 'next/link'
import { ArrowRight, ArrowLeft } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'
import { CodeBlock } from '@/components/code-block'
import { getDocMetadata } from '@/config/docs-metadata'

export const metadata = getDocMetadata('/docs/admin/overview')

export default function AdminOverviewPage() {
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
                Admin Overview
              </h1>
              <p className="text-lg text-muted-foreground leading-relaxed">
                Grit ships with a Filament-inspired admin panel that gives you a full-featured
                back-office from simple TypeScript definitions. Define your resources once and
                get data tables, forms, dashboards, and navigation automatically.
              </p>
            </div>

            <div className="prose-grit">
              {/* What is the Admin Panel */}
              <h2>What is the Admin Panel?</h2>
              <p>
                The Grit admin panel is a <strong>resource-based admin dashboard</strong> inspired
                by <a href="https://filamentphp.com" target="_blank" rel="noreferrer">Laravel Filament</a>.
                Instead of hand-coding every CRUD page, data table, and form, you describe your
                resources in a declarative TypeScript file and the admin panel renders everything
                for you &mdash; columns, filters, form fields, validation, dashboard widgets, and
                sidebar navigation.
              </p>
              <p>
                It lives in <code>apps/admin/</code> as a standalone Next.js application inside
                the Grit monorepo. It connects to your Go API over HTTP and uses React Query for
                data fetching, caching, and optimistic updates.
              </p>

              {/* How it Works */}
              <h2>How It Works</h2>
              <p>
                The workflow is straightforward:
              </p>
              <ol>
                <li><strong>Define a resource</strong> in <code>apps/admin/resources/</code> using the <code>defineResource()</code> helper.</li>
                <li><strong>Register it</strong> in <code>resources/index.ts</code> so the admin panel discovers it.</li>
                <li><strong>That&apos;s it.</strong> The sidebar link, data table page, create/edit forms, and dashboard widgets are generated automatically.</li>
              </ol>
              <p>
                When you run <code>grit generate resource Post</code>, the CLI creates both the
                Go backend artifacts (model, handler, service) and the admin resource definition
                file. You get a working admin page for your new resource with zero manual wiring.
              </p>
            </div>

            {/* Resource definition example */}
            <div className="mt-8 mb-8">
              <CodeBlock language="typescript" filename="apps/admin/resources/posts.ts" code={`import { defineResource } from '@grit/admin'

export default defineResource({
  name: 'Post',
  endpoint: '/api/posts',
  icon: 'FileText',

  table: {
    columns: [
      { key: 'id', label: 'ID', sortable: true },
      { key: 'title', label: 'Title', searchable: true },
      { key: 'status', label: 'Status', badge: {
        published: { color: 'green', label: 'Published' },
        draft: { color: 'yellow', label: 'Draft' },
      }},
      { key: 'created_at', label: 'Created', format: 'relative' },
    ],
    actions: ['create', 'edit', 'delete'],
  },

  form: {
    fields: [
      { key: 'title', label: 'Title', type: 'text', required: true },
      { key: 'content', label: 'Content', type: 'textarea' },
      { key: 'status', label: 'Status', type: 'select',
        options: ['published', 'draft'], default: 'draft' },
    ],
  },
})`} />
            </div>

            <div className="prose-grit">
              <p>
                The definition above produces a fully functional admin page with a sortable,
                searchable data table, create and edit modals with validation, badge-styled
                status column, and relative timestamps &mdash; all without writing any JSX.
              </p>

              {/* Folder Structure */}
              <h2>Admin Panel Folder Structure</h2>
              <p>
                The admin panel follows Next.js App Router conventions with a clear separation
                between layout components, reusable UI, data-fetching hooks, resource
                definitions, and utilities.
              </p>
            </div>

            <div className="mt-4 mb-8">
              <CodeBlock filename="apps/admin/" code={`apps/admin/
\u251c\u2500\u2500 app/
\u2502   \u251c\u2500\u2500 layout.tsx            # Root layout with providers
\u2502   \u251c\u2500\u2500 page.tsx              # Dashboard (admin home)
\u2502   \u2514\u2500\u2500 resources/
\u2502       \u251c\u2500\u2500 users/page.tsx     # User management page
\u2502       \u2514\u2500\u2500 posts/page.tsx     # Post management page
\u251c\u2500\u2500 components/
\u2502   \u251c\u2500\u2500 layout/               # Shell, sidebar, navbar
\u2502   \u2502   \u251c\u2500\u2500 admin-layout.tsx
\u2502   \u2502   \u251c\u2500\u2500 sidebar.tsx
\u2502   \u2502   \u2514\u2500\u2500 navbar.tsx
\u2502   \u251c\u2500\u2500 tables/               # DataTable system
\u2502   \u2502   \u251c\u2500\u2500 data-table.tsx
\u2502   \u2502   \u251c\u2500\u2500 columns.tsx
\u2502   \u2502   \u2514\u2500\u2500 filters.tsx
\u2502   \u251c\u2500\u2500 forms/                # Form builder system
\u2502   \u2502   \u251c\u2500\u2500 form-builder.tsx
\u2502   \u2502   \u251c\u2500\u2500 form-modal.tsx
\u2502   \u2502   \u2514\u2500\u2500 fields/
\u2502   \u2502       \u251c\u2500\u2500 text-field.tsx
\u2502   \u2502       \u251c\u2500\u2500 select-field.tsx
\u2502   \u2502       \u251c\u2500\u2500 date-field.tsx
\u2502   \u2502       \u2514\u2500\u2500 file-field.tsx
\u2502   \u2514\u2500\u2500 widgets/              # Dashboard widgets
\u2502       \u251c\u2500\u2500 stats-card.tsx
\u2502       \u251c\u2500\u2500 chart-widget.tsx
\u2502       \u2514\u2500\u2500 recent-activity.tsx
\u251c\u2500\u2500 hooks/                    # React Query data hooks
\u2502   \u251c\u2500\u2500 use-auth.ts
\u2502   \u251c\u2500\u2500 use-users.ts
\u2502   \u2514\u2500\u2500 use-posts.ts
\u251c\u2500\u2500 resources/                # Resource definitions
\u2502   \u251c\u2500\u2500 index.ts              # Registry (exports all resources)
\u2502   \u251c\u2500\u2500 users.ts
\u2502   \u2514\u2500\u2500 posts.ts
\u2514\u2500\u2500 lib/
    \u251c\u2500\u2500 api-client.ts         # Axios wrapper for Go API
    \u251c\u2500\u2500 auth.ts
    \u2514\u2500\u2500 utils.ts`} />
            </div>

            <div className="prose-grit">
              <h2>Built-in Features</h2>

              <h3>Collapsible Sidebar</h3>
              <p>
                The sidebar is auto-generated from your registered resources. Each resource
                appears as a navigation item with its configured Lucide icon and label. The
                sidebar supports collapsing to icon-only mode, active state highlighting, and
                a user profile section at the bottom with logout functionality.
              </p>

              <h3>Dark and Light Theme</h3>
              <p>
                The admin panel defaults to a premium dark theme inspired by tools like Linear
                and Vercel Dashboard. Colors use the Grit palette &mdash; deep
                backgrounds (<code>#0a0a0f</code>, <code>#111118</code>), purple
                accent (<code>#6c5ce7</code>), and carefully calibrated text hierarchy.
                A theme toggle in the sidebar lets users switch to light mode.
              </p>

              <h3>Responsive Layout</h3>
              <p>
                On desktop the sidebar is always visible. On tablets and mobile devices it
                collapses into a hamburger menu overlay. Data tables switch to horizontal
                scrolling and cards stack vertically. The entire admin panel is usable on
                any screen size.
              </p>

              <h3>Authentication Guard</h3>
              <p>
                The admin layout wraps all routes in an authentication guard. Unauthenticated
                users are redirected to the login page. Role-based access is enforced &mdash;
                only users with the <code>admin</code> role can access the admin panel by default.
                This is configurable per-resource via the <code>permissions</code> option in the
                resource definition.
              </p>

              <h2>System Pages</h2>
              <p>
                Beyond your custom resources, the admin panel includes several built-in system
                pages for managing infrastructure concerns:
              </p>

              <h3>Jobs Dashboard</h3>
              <p>
                Monitor your Redis-backed background job queue. The jobs dashboard shows live
                counts for pending, active, completed, and failed jobs. You can inspect
                individual failed jobs to see error details, retry them, or clear entire queues.
                Queue stats update in real time via polling.
              </p>

              <h3>File Browser</h3>
              <p>
                Browse and manage all files uploaded to your S3-compatible storage (AWS S3,
                Cloudflare R2, or local MinIO). The file browser displays a grid of thumbnails
                with file metadata &mdash; name, size, MIME type, upload date. You can delete
                files, copy URLs, and preview images directly in the admin panel.
              </p>

              <h3>Cron Viewer</h3>
              <p>
                View all registered cron tasks, their schedules (cron expressions), last run
                time, and next scheduled execution. The cron viewer gives you visibility into
                background scheduling without checking server logs.
              </p>

              <h3>Email Preview</h3>
              <p>
                Preview your Go HTML email templates rendered with sample data. This lets you
                iterate on email design without sending real emails. In development, emails
                sent through the mail service are also captured by Mailhog and visible at
                its web UI.
              </p>

              <h2>UI Aesthetic</h2>
              <p>
                The admin panel is designed to feel like a <strong>premium CRM or SaaS
                dashboard</strong>, not a generic CRUD generator. Key design principles:
              </p>
              <ul>
                <li><strong>Generous spacing</strong> &mdash; content breathes with consistent padding and margins.</li>
                <li><strong>Professional data density</strong> &mdash; tables show enough information without overwhelming the user.</li>
                <li><strong>Beautiful empty states</strong> &mdash; when a resource has no data, you see a polished illustration with a call to action, not a blank page.</li>
                <li><strong>Polished loading states</strong> &mdash; skeleton loaders match the exact layout of the content they replace.</li>
                <li><strong>Subtle animations</strong> &mdash; hover effects, transitions, and micro-interactions make the interface feel alive.</li>
                <li><strong>Typography hierarchy</strong> &mdash; DM Sans for UI text, JetBrains Mono for code and data, with carefully chosen font weights and sizes.</li>
              </ul>
            </div>

            {/* What's next */}
            <div className="mt-10 mb-10">
              <h2 className="text-2xl font-semibold tracking-tight mb-4">
                What&apos;s Next?
              </h2>
              <div className="space-y-3">
                {[
                  { label: 'Resource Definitions', href: '/docs/admin/resources', desc: 'Learn the full defineResource() API and configuration options.' },
                  { label: 'DataTable', href: '/docs/admin/datatable', desc: 'Server-side pagination, sorting, filtering, and custom cell renderers.' },
                  { label: 'Form Builder', href: '/docs/admin/forms', desc: 'All field types, validation, layout options, and create/edit modes.' },
                  { label: 'Dashboard & Widgets', href: '/docs/admin/widgets', desc: 'Stats cards, charts, activity feeds, and widget configuration.' },
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
                <Link href="/docs/backend/response-format" className="gap-1.5">
                  <ArrowLeft className="h-3.5 w-3.5" />
                  API Response Format
                </Link>
              </Button>
              <Button variant="ghost" size="sm" asChild className="text-muted-foreground/60 hover:text-foreground">
                <Link href="/docs/admin/resources" className="gap-1.5">
                  Resource Definitions
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
