import Link from 'next/link'
import { ArrowRight, ArrowLeft } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'
import { CodeBlock } from '@/components/code-block'
import { Callout } from '@/components/callout'
import { LaneFlow } from '@/components/lane-flow'
import { getDocMetadata } from '@/config/docs-metadata'

export const metadata = getDocMetadata('/docs/admin/roles')

export default function AdminRolesPage() {
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
                Roles &amp; Permissions UI
              </h1>
              <p className="text-lg text-muted-foreground leading-relaxed">
                Every scaffolded admin ships a permission editor at{' '}
                <code className="text-sm font-mono text-primary bg-primary/10 rounded px-1.5 py-0.5">
                  /system/roles
                </code>
                . It lists roles, edits their grants through a CRUD matrix, and hides
                nav entries and buttons the current user has no permission to use.
              </p>
              <LaneFlow
                id="admin-roles"
                lanes={['Admin UI', 'Go API']}
                nodes={[
                  { id: 'editor', lane: 0, row: 0, title: 'Roles screen', sub: 'CRUD matrix', tone: 'primary' },
                  { id: 'hook', lane: 0, row: 2, title: 'usePermissions()', sub: 'can(...)', tone: 'cyan' },
                  { id: 'grants', lane: 1, row: 0, title: 'roles.grants', sub: 'JSON column', tone: 'green' },
                  { id: 'guard', lane: 1, row: 2, title: 'RequireRole', sub: 'enforced server-side', tone: 'amber' },
                ]}
                edges={[
                  { from: 'editor', to: 'grants', label: 'save', tone: 'green' },
                  { from: 'grants', to: 'hook', label: '/auth/permissions', tone: 'cyan' },
                  { from: 'hook', to: 'guard', tone: 'amber' },
                ]}
              />
            </div>

            <Callout type="warning" title="The UI is a courtesy, not the boundary">
              Hiding a button only stops it from being clicked. Every route is enforced
              server-side by{' '}
              <code className="text-xs font-mono text-primary bg-primary/10 rounded px-1.5 py-0.5">
                RequireRole
              </code>
              , and the API rejects a request whether or not the button was rendered.
              Treat the admin screen as a convenience layer over{' '}
              <Link href="/docs/security/authorization" className="text-primary hover:underline">
                the authorization model
              </Link>
              .
            </Callout>

            {/* The editor */}
            <div className="mb-12 mt-12">
              <h2 className="text-2xl font-semibold tracking-tight mb-4">The permission editor</h2>
              <p className="text-muted-foreground leading-relaxed mb-4">
                Permissions are grouped by module, then by feature. Each level has a
                tri-state checkbox, so you can grant an entire module without ticking
                every row underneath it.
              </p>
              <ul className="space-y-3 text-muted-foreground leading-relaxed mb-6">
                <li>
                  <strong className="text-foreground">CRUD matrix per feature.</strong>{' '}
                  Actions a feature doesn&apos;t declare render as an em dash rather than
                  a checkbox that would silently do nothing when ticked.
                </li>
                <li>
                  <strong className="text-foreground">Live counter and filter.</strong>{' '}
                  A &ldquo;granted / total&rdquo; count updates as you edit, and a text
                  filter narrows a long catalog.
                </li>
                <li>
                  <strong className="text-foreground">Copy permissions from.</strong>{' '}
                  Start a new role from an existing one instead of re-ticking a matrix.
                </li>
                <li>
                  <strong className="text-foreground">Built-in roles are locked.</strong>{' '}
                  Their names can&apos;t be edited and the delete button is hidden. The
                  server enforces both independently.
                </li>
              </ul>
              <Callout type="tip" title="Wildcards survive a round trip">
                Selection is seeded from the server&apos;s expanded permission list, then
                collapsed back to wildcards on save. A role granted{' '}
                <code className="text-xs font-mono text-primary bg-primary/10 rounded px-1.5 py-0.5">
                  products.*
                </code>{' '}
                keeps that wildcard, so it automatically inherits any action added to{' '}
                <code className="text-xs font-mono text-primary bg-primary/10 rounded px-1.5 py-0.5">
                  products
                </code>{' '}
                later &mdash; editing a role in the UI never silently freezes it to
                today&apos;s action list.
              </Callout>
            </div>

            {/* usePermissions */}
            <div className="mb-12">
              <h2 className="text-2xl font-semibold tracking-tight mb-4">
                Gating your own UI
              </h2>
              <p className="text-muted-foreground leading-relaxed mb-4">
                The{' '}
                <code className="text-sm font-mono text-primary bg-primary/10 rounded px-1.5 py-0.5">
                  usePermissions()
                </code>{' '}
                hook reads the current user&apos;s effective grants once and caches them.
                Use it to hide actions the API would reject anyway.
              </p>
              <CodeBlock
                language="tsx"
                filename="app/(dashboard)/products/page.tsx"
                code={`import { usePermissions } from "@/hooks/use-permissions";

export default function ProductsPage() {
  const { can, isSuper, isLoading } = usePermissions();

  return (
    <>
      <ProductTable />

      {/* exact permission */}
      {can("products.delete") && <BulkDeleteButton />}

      {/* any action on the resource */}
      {can("products.*") && <ProductToolbar />}
    </>
  );
}`}
              />
              <p className="text-muted-foreground leading-relaxed mt-4">
                <code className="text-sm font-mono text-primary bg-primary/10 rounded px-1.5 py-0.5">
                  can()
                </code>{' '}
                returns <strong className="text-foreground">false while permissions are
                still loading</strong>, and returns true for everything when the user is a
                super admin. Gated UI therefore stays hidden until grants are known rather
                than flashing into view and disappearing &mdash; a flash of forbidden UI
                looks broken and leaks the shape of the admin to users who can&apos;t use it.
              </p>
            </div>

            {/* Nav gating */}
            <div className="mb-12">
              <h2 className="text-2xl font-semibold tracking-tight mb-4">
                Navigation is gated the same way
              </h2>
              <p className="text-muted-foreground leading-relaxed mb-4">
                Sidebar entries declare the permission they need. An entry whose
                permission the user lacks is not rendered, so the nav never links to a
                page the API will refuse.
              </p>
              <CodeBlock
                language="tsx"
                filename="components/layout/sidebar.tsx"
                code={`{ href: "/system/roles", label: "Roles & permissions",
  iconKey: "ShieldCheck", adminOnly: true, requires: "roles.view" },`}
              />
              <p className="text-muted-foreground leading-relaxed mt-4">
                Generated resources register their own entry and their own{' '}
                <code className="text-sm font-mono text-primary bg-primary/10 rounded px-1.5 py-0.5">
                  &lt;resource&gt;.&lt;action&gt;
                </code>{' '}
                permissions when you run{' '}
                <code className="text-sm font-mono text-primary bg-primary/10 rounded px-1.5 py-0.5">
                  grit generate resource
                </code>
                , and both are removed again by{' '}
                <code className="text-sm font-mono text-primary bg-primary/10 rounded px-1.5 py-0.5">
                  grit remove resource
                </code>
                .
              </p>
            </div>

            {/* Both admins */}
            <div className="mb-12">
              <h2 className="text-2xl font-semibold tracking-tight mb-4">
                Identical in both admins
              </h2>
              <p className="text-muted-foreground leading-relaxed">
                The Next.js admin and the Vite/TanStack admin render the{' '}
                <em>same</em> roles component, transformed at scaffold time rather than
                written twice. That is deliberate: a permission editor that disagreed
                between the two frontends would be a security-shaped bug, not a cosmetic
                one.
              </p>
            </div>

            {/* Related */}
            <div className="mb-12">
              <h2 className="text-2xl font-semibold tracking-tight mb-4">Related</h2>
              <div className="grid gap-2">
                {[
                  { label: 'Roles & Permissions', href: '/docs/security/authorization', desc: 'The permission model, wildcards, and server-side guards.' },
                  { label: 'Admin Resources', href: '/docs/admin/resources', desc: 'Define the screens the permissions apply to.' },
                  { label: 'Multi-tenancy', href: '/docs/plugins/multitenant', desc: 'Per-organization roles via the multitenant plugin.' },
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
                <Link href="/docs/admin/widgets" className="gap-1.5">
                  <ArrowLeft className="h-3.5 w-3.5" />
                  Dashboard &amp; Widgets
                </Link>
              </Button>
              <Button variant="ghost" size="sm" asChild className="text-muted-foreground/60 hover:text-foreground">
                <Link href="/docs/security/authorization" className="gap-1.5">
                  Authorization
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
