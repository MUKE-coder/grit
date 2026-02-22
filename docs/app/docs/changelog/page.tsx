import Link from 'next/link'
import { ArrowLeft } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'
import { getDocMetadata } from '@/config/docs-metadata'

export const metadata = getDocMetadata('/docs/changelog')

export default function ChangelogPage() {
  return (
    <div className="min-h-screen bg-background">
      <SiteHeader />
      <DocsSidebar />

      <main className="lg:pl-64">
        <div className="container max-w-screen-xl py-10 px-6">
          <div className="max-w-3xl">
            {/* Header */}
            <div className="mb-10">
              <span className="tag-mono text-primary/80 mb-3 block">Release History</span>
              <h1 className="text-4xl font-bold tracking-tight mb-4">
                Changelog
              </h1>
              <p className="text-lg text-muted-foreground leading-relaxed">
                All notable changes to Grit are documented here. Each release includes new features,
                bug fixes, and any breaking changes you need to be aware of.
              </p>
            </div>

            {/* v0.17.0 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-accent/15 px-3 py-1 text-sm font-semibold text-primary">
                  v0.17.0
                </span>
                <span className="text-sm text-muted-foreground">February 22, 2026</span>
              </div>

              <div className="prose-grit">
                <h3>Features</h3>
                <ul>
                  <li>
                    <strong>API Documentation (gin-docs)</strong> &mdash; Replaced hand-written Scalar/OpenAPI
                    spec with{' '}
                    <a href="https://github.com/MUKE-coder/gin-docs" className="text-primary hover:underline" target="_blank" rel="noopener noreferrer">gin-docs</a>,
                    a zero-annotation API documentation generator. Routes and GORM models are introspected
                    automatically to produce an OpenAPI 3.1 spec with interactive Scalar or Swagger UI,
                    plus Postman and Insomnia export.
                  </li>
                  <li>
                    <strong>Dark/Light mode for Go Playground</strong> &mdash; The playground now follows the
                    site-wide theme toggle, switching between VS Code dark and light CodeMirror themes.
                  </li>
                  <li>
                    <strong>Umami Analytics</strong> &mdash; Optional visitor analytics via self-hosted Umami,
                    configured with <code>NEXT_PUBLIC_UMAMI_WEBSITE_ID</code> environment variable.
                  </li>
                </ul>

                <h3>Documentation</h3>
                <ul>
                  <li>
                    New{' '}
                    <Link href="/docs/backend/api-docs" className="text-primary hover:underline">API Documentation</Link> page
                    covering gin-docs configuration, GORM model schemas, route customization, UI switching, and spec export
                  </li>
                  <li>Full SEO + AEO implementation: sitemap, robots.txt, JSON-LD structured data, per-page metadata</li>
                </ul>

                <h3>Infrastructure</h3>
                <ul>
                  <li>Added Dockerfile for docs site deployment (Next.js standalone output)</li>
                  <li>Google Search Console verification</li>
                </ul>
              </div>
            </div>

            {/* v0.16.0 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-muted/50 px-3 py-1 text-sm font-semibold text-muted-foreground">
                  v0.16.0
                </span>
                <span className="text-sm text-muted-foreground">February 21, 2026</span>
              </div>

              <div className="prose-grit">
                <h3>Features</h3>
                <ul>
                  <li>
                    <strong>Go Playground</strong> &mdash; Interactive code editor at{' '}
                    <Link href="/playground" className="text-primary hover:underline">/playground</Link> with
                    Go syntax highlighting, code execution via the official Go Playground API, example snippets,
                    share links, and keyboard shortcuts (Ctrl+Enter to run).
                  </li>
                  <li>
                    <strong>GORM Studio updated</strong> &mdash; Updated to latest version with raw SQL editor,
                    schema export (SQL/JSON/YAML/DBML/ERD), data import/export (JSON/CSV/SQL/XLSX),
                    and Go model generation from database schema.
                  </li>
                </ul>

                <h3>Documentation</h3>
                <ul>
                  <li>
                    <Link href="/docs/prerequisites/golang" className="text-primary hover:underline">Go for Grit Developers</Link> &mdash;
                    comprehensive rewrite with 22 sections covering methods, Gin routing, middleware, CORS,
                    handler/service architecture, GORM CRUD, migrations, seeding, JWT auth flow, and RBAC
                  </li>
                  <li>Fixed right-side table of contents for the Go prerequisites page</li>
                  <li>New Middleware and CORS sections added to Go guide</li>
                </ul>
              </div>
            </div>

            {/* v0.15.0 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-muted/50 px-3 py-1 text-sm font-semibold text-muted-foreground">
                  v0.15.0
                </span>
                <span className="text-sm text-muted-foreground">February 20, 2026</span>
              </div>

              <div className="prose-grit">
                <h3>Features</h3>
                <ul>
                  <li>
                    <strong>Security (Sentinel)</strong> &mdash; Every <code>grit new</code> project now ships with
                    a production-grade security suite powered by{' '}
                    <Link href="https://github.com/MUKE-coder/sentinel" className="text-primary hover:underline">Sentinel</Link>.
                    Includes WAF, rate limiting, brute-force protection, anomaly detection, IP geolocation,
                    security headers, and a real-time threat dashboard at <code>/sentinel/ui</code>.
                    See <Link href="/docs/batteries/security" className="text-primary hover:underline">Security docs</Link>.
                  </li>
                  <li>
                    <strong>Admin security page</strong> &mdash; New System &rarr; Security page in the admin panel
                    embeds the Sentinel dashboard for monitoring threats without leaving the admin UI.
                  </li>
                </ul>

                <h3>Documentation</h3>
                <ul>
                  <li>New: <Link href="/docs/batteries/security" className="text-primary hover:underline">Security (Sentinel)</Link> documentation page</li>
                  <li>Migrated getting-started pages (Installation, Quick Start, Troubleshooting) to use CodeBlock component</li>
                  <li>Added prerequisite learning pages for Go, Next.js, and Docker</li>
                </ul>
              </div>
            </div>

            {/* v0.14.0 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-muted/50 px-3 py-1 text-sm font-semibold text-muted-foreground">
                  v0.14.0
                </span>
                <span className="text-sm text-muted-foreground">February 18, 2026</span>
              </div>

              <div className="prose-grit">
                <h3>Features</h3>
                <ul>
                  <li>
                    <strong>Multi-step forms</strong> &mdash; New <code>formView: &quot;modal-steps&quot;</code> and{' '}
                    <code>&quot;page-steps&quot;</code> variants with horizontal/vertical step indicators,
                    per-step validation, progress bar, and clickable step navigation.
                    See <Link href="/docs/admin/multi-step-forms" className="text-primary hover:underline">Multi-Step Forms</Link>.
                  </li>
                  <li>
                    <strong>Standalone component usage</strong> &mdash; FormBuilder, FormStepper, and DataTable
                    can now be used on any page in both web and admin apps without the resource system.
                    See <Link href="/docs/admin/standalone-usage" className="text-primary hover:underline">Standalone Usage</Link>.
                  </li>
                  <li>
                    <strong>Richtext field type</strong> &mdash; New <code>richtext</code> field with Tiptap WYSIWYG
                    editor (bold, italic, headings, lists, code blocks, links, undo/redo).
                  </li>
                  <li>
                    <strong><code>string_array</code> field type</strong> &mdash; Store arrays of strings
                    using <code>datatypes.JSONSlice[string]</code>. Works with PostgreSQL and SQLite.
                    Maps to <code>string[]</code> in TypeScript and <code>z.array(z.string())</code> in Zod.
                  </li>
                  <li>
                    <strong>Built-in blog example</strong> &mdash; <code>grit new</code> now scaffolds a complete
                    blog with model, service, handler, seed data, public web pages, and admin resource definition.
                  </li>
                  <li>
                    <strong>Sidebar user avatar</strong> &mdash; Admin sidebar shows the current user&apos;s avatar
                    with a dropdown menu for profile and logout.
                  </li>
                  <li>
                    <strong>Profile avatar upload</strong> &mdash; Profile page now supports avatar image upload.
                  </li>
                  <li>
                    <strong><code>react-hook-form</code> in web app</strong> &mdash; Web app scaffold now includes{' '}
                    <code>react-hook-form</code> as a dependency, enabling standalone FormBuilder usage out of the box.
                  </li>
                </ul>

                <h3>Bug Fixes</h3>
                <ul>
                  <li>
                    <strong>Scalar API docs crash</strong> &mdash; Fixed <code>c.String</code> treating HTML as
                    a format string. Now uses <code>c.Data</code> to avoid panics when Scalar HTML
                    contains <code>%</code> characters in CSS/JS.
                  </li>
                  <li>
                    <strong>Blog route conflict</strong> &mdash; Admin blog CRUD routes moved
                    from <code>/api/blogs</code> to <code>/api/admin/blogs</code> to avoid conflict
                    with public blog routes.
                  </li>
                  <li>
                    <strong>Select dropdown styling</strong> &mdash; Fixed relationship select dropdown
                    rendering behind modals using portal-based positioning.
                  </li>
                </ul>

                <h3>Documentation</h3>
                <ul>
                  <li>New: <Link href="/docs/tutorials/product-catalog" className="text-primary hover:underline">Build a Product Catalog</Link> tutorial &mdash; resource generation, multi-step forms, standalone DataTable &amp; FormBuilder</li>
                  <li>New: <Link href="/docs/admin/multi-step-forms" className="text-primary hover:underline">Multi-Step Forms</Link> guide</li>
                  <li>New: <Link href="/docs/admin/standalone-usage" className="text-primary hover:underline">Standalone Usage</Link> guide</li>
                  <li>New: Changelog page</li>
                  <li>Updated CLI Commands, Code Generation, Quick Start, Resources, Shared Package, Web App, Seeders, and Forms pages</li>
                </ul>
              </div>
            </div>

            {/* v0.12.0 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-muted/50 px-3 py-1 text-sm font-semibold text-muted-foreground">
                  v0.12.0
                </span>
                <span className="text-sm text-muted-foreground">February 2026</span>
              </div>

              <div className="prose-grit">
                <h3>Features</h3>
                <ul>
                  <li>
                    <strong>Relationship support</strong> &mdash; New <code>belongs_to</code> and{' '}
                    <code>many_to_many</code> field types for the code generator. Automatically creates
                    foreign keys, junction tables, and relationship-aware form fields.
                  </li>
                  <li>
                    <strong>Relationship select fields</strong> &mdash; New <code>relationship-select</code> and{' '}
                    <code>multi-relationship-select</code> form field components with search, portal-based
                    dropdowns, and tag-based multi-select.
                  </li>
                  <li>
                    <strong>Beginner tutorial</strong> &mdash; &quot;Learn Grit Step by Step&quot; tutorial
                    walking through building a full-stack app from scratch.
                  </li>
                </ul>
              </div>
            </div>

            {/* v0.11.0 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-muted/50 px-3 py-1 text-sm font-semibold text-muted-foreground">
                  v0.11.0
                </span>
                <span className="text-sm text-muted-foreground">February 2026</span>
              </div>

              <div className="prose-grit">
                <h3>Features</h3>
                <ul>
                  <li>
                    <strong>Full-page form view</strong> &mdash; New <code>formView: &quot;page&quot;</code> option
                    renders forms as dedicated pages instead of modals.
                  </li>
                  <li>
                    <strong><code>slug</code> field type</strong> &mdash; Auto-generates URL-friendly slugs with
                    unique suffixes. Excluded from create/update forms and Zod schemas.
                  </li>
                  <li>
                    <strong>DataTable column customization</strong> &mdash; Hide/show columns, column visibility
                    toggle in table toolbar.
                  </li>
                  <li>
                    <strong><code>grit start</code> commands</strong> &mdash; <code>grit start client</code> and{' '}
                    <code>grit start server</code> for running frontend and API separately.
                  </li>
                </ul>
              </div>
            </div>

            {/* v0.10.0 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-muted/50 px-3 py-1 text-sm font-semibold text-muted-foreground">
                  v0.10.0
                </span>
                <span className="text-sm text-muted-foreground">January 2026</span>
              </div>

              <div className="prose-grit">
                <h3>Features</h3>
                <ul>
                  <li>
                    <strong>Style variants</strong> &mdash; <code>--style</code> flag for <code>grit new</code> with
                    4 admin panel styles: default, modern, minimal, and glass.
                  </li>
                  <li>
                    <strong>Air hot reloading</strong> &mdash; Go API development with automatic rebuild on file
                    changes using Air.
                  </li>
                  <li>
                    <strong><code>grit remove resource</code></strong> &mdash; Remove a generated resource and
                    clean up all injected code (model, handler, routes, schemas, types, hooks, admin pages).
                  </li>
                  <li>
                    <strong>AI workflow docs</strong> &mdash; Guides for using Grit with Claude and Antigravity AI assistants.
                  </li>
                </ul>
              </div>
            </div>

            {/* Nav */}
            <div className="flex items-center justify-between pt-6 border-t border-border/30">
              <Button variant="ghost" size="sm" asChild className="text-muted-foreground/60 hover:text-foreground">
                <Link href="/docs" className="gap-1.5">
                  <ArrowLeft className="h-3.5 w-3.5" />
                  Introduction
                </Link>
              </Button>
              <div />
            </div>
          </div>
        </div>
      </main>
    </div>
  )
}
