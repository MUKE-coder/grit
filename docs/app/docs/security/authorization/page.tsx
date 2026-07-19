import Link from 'next/link'
import { ArrowLeft, ArrowRight, ShieldCheck, AlertTriangle } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'
import { CodeBlock } from '@/components/code-block'

export const metadata = {
  title: 'Roles & Permissions — Grit',
  description:
    'How authorization works in Grit: permission keys, roles as bags of grants, guarding routes, and the admin UI.',
  alternates: { canonical: 'https://gritframework.dev/docs/security/authorization' },
}

export default function AuthorizationPage() {
  return (
    <div className="min-h-screen bg-background">
      <SiteHeader />
      <DocsSidebar />

      <main className="lg:pl-64">
        <div className="mx-auto max-w-3xl px-6 py-12">
          <div className="mb-3 flex items-center gap-2">
            <ShieldCheck className="h-5 w-5 text-primary" />
            <span className="font-mono text-xs uppercase tracking-wider text-muted-foreground">
              Security
            </span>
          </div>

          <h1 className="mb-4 font-display text-4xl font-bold tracking-tight">
            Roles &amp; Permissions
          </h1>
          <p className="mb-10 text-lg text-muted-foreground">
            Grit guards routes by <strong>permission</strong>, not by role name. Roles are
            named bags of permissions you edit at runtime, so adding a role never means
            editing routes and redeploying.
          </p>

          <div className="prose-grit">
            <h2>The shape of a permission</h2>
            <p>
              A permission key is <code>&lt;resource&gt;.&lt;action&gt;</code>, where the
              action is one of <code>create</code>, <code>view</code>, <code>edit</code> or{' '}
              <code>delete</code>:
            </p>
            <CodeBlock language="text" code={`products.create
users.delete
uploads.view`} />

            <p>Grants may use wildcards:</p>
            <CodeBlock language="text" code={`*             every permission (superuser)
products.*    every action on products
*.view        view on every resource`} />

            <h2>Roles hold grants</h2>
            <p>
              Three roles are seeded on first migrate: <code>ADMIN</code> (holds{' '}
              <code>*</code>), <code>EDITOR</code> (content) and <code>USER</code> (no
              administrative permissions). Seeding is idempotent and{' '}
              <strong>never overwrites</strong> grants you have edited.
            </p>
            <p>
              Grants are stored <em>as authored</em>. If a role holds{' '}
              <code>products.*</code> and you later add a new action to that resource, the
              role picks it up automatically — nothing to re-grant.
            </p>

            <h2>Guarding a route</h2>
            <p>
              <code>RequireRole</code> accepts role names, permissions (prefixed{' '}
              <code>perm:</code>), or both. Access is granted if <strong>any</strong>{' '}
              argument matches:
            </p>
            <CodeBlock language="go" filename="internal/routes/routes.go" code={`// Permission-based (preferred)
admin.DELETE("/products/:id",
    middleware.RequireRole("perm:products.delete"),
    productHandler.Delete)

// Mixed — useful while migrating. Either passes.
admin.PUT("/users/:id",
    middleware.RequireRole("ADMIN", "perm:users.edit"),
    userHandler.Update)

// Role-only still works. Nothing existing had to change.
admin.Use(middleware.RequireRole("ADMIN"))`} />

            <p>
              Inside a handler, use <code>authz.Can</code> for conditional logic:
            </p>
            <CodeBlock language="go" code={`if authz.Can(h.DB, userID, "products.delete") {
    // show the destructive action
}`} />

            <h2>Generated resources register themselves</h2>
            <p>
              <code>grit generate resource Product</code> adds{' '}
              <code>products.create|view|edit|delete</code> to the catalog automatically, so
              a new resource is grantable immediately.{' '}
              <code>grit remove resource</code> takes them back out. Hand-written
              permissions belong in <code>coreModules()</code> in{' '}
              <code>internal/authz/permissions.go</code>; machine-written entries live in{' '}
              <code>generatedModules()</code> between the marker comments.
            </p>

            <h2>Managing roles</h2>
            <p>
              The admin panel has a <strong>Roles &amp; permissions</strong> screen under
              System (<code>/system/roles</code>): create roles, tick permissions in a tree
              with a CRUD matrix, copy permissions from an existing role, and see how many
              users hold each one.
            </p>
            <p>
              Built-in roles can have their permissions edited, but they cannot be renamed
              or deleted — routes and the upgrade path resolve them by name. This is
              enforced by the API, not just hidden in the UI.
            </p>

            <h2>In the frontend</h2>
            <p>
              The API returns the caller&apos;s permissions already expanded, so the client
              never re-implements wildcard matching:
            </p>
            <CodeBlock language="tsx" code={`const { can, isSuper } = usePermissions()

{can("products.delete") && <DeleteButton />}
{can("products.*") && <ProductsMenu />}`} />

            <div className="my-8 rounded-xl border border-amber-500/30 bg-amber-500/5 p-5">
              <div className="mb-2 flex items-center gap-2">
                <AlertTriangle className="h-4 w-4 text-amber-500" />
                <strong className="text-foreground">
                  Hiding UI is not access control
                </strong>
              </div>
              <p className="mb-0 text-sm text-muted-foreground">
                <code>can()</code> is a convenience for hiding buttons and nav items. It is
                not a security boundary — anyone can call your API directly. Always guard
                the route as well.
              </p>
            </div>

            <h2>Upgrading an existing app</h2>
            <p>Nothing breaks, by construction:</p>
            <ul>
              <li>
                every existing <code>RequireRole(&quot;ADMIN&quot;)</code> call site keeps
                working — the signature did not change
              </li>
              <li>
                if a user has no role assignment yet, grants resolve from the legacy{' '}
                <code>users.role</code> string, so admins do not lose access
              </li>
              <li>
                changing a user&apos;s role in the admin keeps the assignment in step, so
                the dropdown still takes effect
              </li>
            </ul>
            <p>
              Adopt permissions route by route: add a <code>perm:</code> argument alongside
              the role name, then drop the role name once every caller has a role that
              grants it.
            </p>

            <h2>Multiple roles per user</h2>
            <p>
              Assignment is many-to-many. The admin&apos;s user form edits a single primary
              role; assign several with:
            </p>
            <CodeBlock language="bash" code={`PUT /api/users/:id/roles
{ "role_ids": ["...", "..."] }`} />
            <p>A user&apos;s grants are the union of every role they hold.</p>
          </div>

          <div className="mt-12 flex items-center justify-between border-t border-border/40 pt-6">
            <Button asChild variant="ghost">
              <Link href="/docs/security">
                <ArrowLeft className="mr-2 h-4 w-4" />
                Security Guide
              </Link>
            </Button>
            <Button asChild variant="ghost">
              <Link href="/docs/batteries/security">
                Sentinel
                <ArrowRight className="ml-2 h-4 w-4" />
              </Link>
            </Button>
          </div>
        </div>
      </main>
    </div>
  )
}
