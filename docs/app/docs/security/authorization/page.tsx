import Link from 'next/link'
import { ArrowLeft, ArrowRight, ShieldCheck, AlertTriangle } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { DocsShell } from '@/components/docs/shell'
import { CodeBlock } from '@/components/code-block'

export const metadata = {
  title: 'Roles & Permissions — Grit',
  description:
    'How authorization works in Grit: permission keys, roles as bags of grants, guarding routes, and the admin UI.',
  alternates: { canonical: 'https://gritframework.dev/docs/security/authorization' },
}

export default function AuthorizationPage() {
  return (
    <DocsShell>
          <div className="mb-3 flex items-center gap-2">
            <ShieldCheck className="h-5 w-5 text-accent" />
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
staff.DELETE("/products/:id",
    middleware.RequireRole("perm:products.delete"),
    productHandler.Delete)

// Mixed: an ADMIN passes, and so does any role granted users.edit.
staff.PUT("/users/:id",
    middleware.RequireRole("ADMIN", "perm:users.edit"),
    userHandler.Update)

// Role-only still works. Nothing existing had to change.
admin.Use(middleware.RequireRole("ADMIN"))`} />

            <h2 id="staff-routes">Staff routes: how a custom role reaches the admin API</h2>
            <p>
              <code>routes.go</code> has two groups for the admin API, and which one a route is on
              decides who can reach it:
            </p>
            <ul>
              <li>
                <code>staff</code> admits anyone who holds at least one permission, and every route
                on it names the permission it needs. Users, roles, the audit log, jobs, backups,
                blogs and the dashboards live here, so a role granted <code>users.view</code> can
                list users and nothing more.
              </li>
              <li>
                <code>admin</code> requires the <code>ADMIN</code> role and nothing less. Routes that
                name no permission stay here (SSO connections, form sharing, writing settings,
                changing feature flags), and so does anything a plugin adds, so a route that forgets
                its permission fails closed.
              </li>
            </ul>
            <p>
              Generated resources follow the same split. Reading and writing a record is on the
              protected group; delete and bulk are on the staff group and ask for{' '}
              <code>&lt;resource&gt;.delete</code>, so a role can be allowed to delete products
              without being made an admin.
            </p>
            <CodeBlock language="go" filename="internal/routes/product_routes.go" code={`m.Staff.DELETE("/products/:id", middleware.RequireRole("ADMIN", "perm:products.delete"), h.Delete)
m.Staff.POST("/products/bulk", middleware.RequireRole("ADMIN", "perm:products.delete"), h.Bulk)`} />
            <p>
              Before v3.220.0 there was only the admin group. Every route on it sat behind the{' '}
              <code>ADMIN</code> role, so the <code>perm:</code> guards on them never ran, and a
              custom role was refused by every admin endpoint whatever it was granted. The admin
              panel showed an <code>EDITOR</code> pages the API then refused.
            </p>

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

            <h2 id="policies">Policies: a rule about one row</h2>
            <p>
              A permission answers &quot;may this role touch this resource at
              all&quot;, and it is checked before the handler runs. It cannot
              answer anything that depends on the record: editable while it is
              a draft, approvable by anyone but the person who raised it,
              deletable only once archived. Those are <strong>policies</strong>,
              and they live in <code>internal/policies</code>.
            </p>
            <p>
              <code>grit generate resource Invoice</code> writes{' '}
              <code>internal/policies/invoice.go</code> with the three abilities
              the service checks, every one commented out. A rule you do not
              write is not a rule that denies, so a project that ignores the
              package behaves exactly as it did before it had one.
            </p>
            <CodeBlock language="go" filename="internal/policies/invoice.go" code={`func registerInvoices() {
    authz.Define("invoices.update", func(actor authz.Actor, record any) authz.Decision {
        item, ok := record.(*models.Invoice)
        if !ok {
            // Another model reached this rule, which is a wiring mistake
            // rather than a refusal. Let the ability meant for it decide.
            return authz.Allow()
        }
        if item.Status == "paid" {
            return authz.Deny("This invoice is paid, so it is corrected with a credit note rather than edited")
        }
        return authz.Allow()
    })
}`} />
            <p>
              The three abilities are{' '}
              <code>&lt;resource&gt;.read</code>,{' '}
              <code>&lt;resource&gt;.update</code> and{' '}
              <code>&lt;resource&gt;.delete</code>. The service consults them
              around the row it loaded, so they hold on every path into that
              row: the REST handler, a background job, a CLI command and a test
              all go through the same service method. The reason a rule gives
              is the message the client receives, with a 403.
            </p>
            <p>
              An <code>ADMIN</code> is <strong>not</strong> exempt by default.
              That is deliberate: these rules are about the state of a record
              rather than the rank of a user, and &quot;this invoice is already
              paid&quot; is as true for an administrator as for anybody else.{' '}
              <code>authz.Before</code> adds an exemption for the abilities
              that want one.
            </p>
            <p>
              Bulk actions and the editable grid ask per row. A policy that
              refuses twelve of the selected rows refuses the whole action
              rather than applying nine of them, because a bulk action that
              reports twelve and changed nine is worse than one that says no
              and names what it would not touch.
            </p>

            <div className="my-8 rounded-xl border border-amber-500/30 bg-amber-500/5 p-5">
              <div className="mb-2 flex items-center gap-2">
                <AlertTriangle className="h-4 w-4 text-amber-500" />
                <strong className="text-foreground">
                  Before v3.451.0, no rule ran
                </strong>
              </div>
              <p className="mb-0 text-sm text-muted-foreground">
                <code>policies.Register()</code> had no call site. The package
                was generated, documented and injected into correctly, and
                nothing ever called the function that defines the rules, so
                every rule anybody wrote was registered into a map that stayed
                empty. <code>grit upgrade</code> adds the call. If you wrote a
                rule before that version, upgrade and then check it actually
                refuses what you expect, because until now it did not.
              </p>
            </div>

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

            <h2 id="edge-gate">The panel is not sent to a stranger</h2>
            <p>
              The same point, one layer out. Every admin page needs a session,
              and until v3.459.0 that was decided in the browser: the page
              arrived, React mounted, asked <code>/api/auth/me</code> who you
              were, and on a 401 replaced the route with the login. Correct, and
              it meant the panel was downloaded and run for somebody who had
              never signed in.
            </p>
            <p>
              A generated project gates it before that. A Next.js admin app has{' '}
              <code>apps/admin/proxy.ts</code>; a panel inside a web app is
              covered by that app&apos;s <code>proxy.ts</code>, which carries the
              same check as <code>adminGate</code>; a Vite admin does it in its
              router&apos;s <code>beforeLoad</code>.
            </p>
            <CodeBlock language="ts" code={`// apps/admin/proxy.ts: the pages a signed-out visitor may reach.
const PUBLIC_PATHS = [
  "/login", "/sign-up", "/forgot-password",
  "/reset-password", "/verify-email", "/magic-link", "/callback",
];

// Everything else needs one of two cookies, and neither holds a secret:
//   grit_admin_session  written on the admin's origin after a sign-in
//   grit_signed_in      set by the API, when the API shares this host`} />
            <p>
              The matcher denies by default, so a page you add later, with{' '}
              <code>grit generate</code> or a plugin, is gated the day it lands.
              Add a page a visitor must reach without signing in and add its path
              to <code>PUBLIC_PATHS</code>, or it will redirect to the login.
            </p>

            <div className="my-8 rounded-xl border border-amber-500/30 bg-amber-500/5 p-5">
              <div className="mb-2 flex items-center gap-2">
                <AlertTriangle className="h-4 w-4 text-amber-500" />
                <strong className="text-foreground">
                  The gate is not access control either
                </strong>
              </div>
              <p className="mb-0 text-sm text-muted-foreground">
                Those cookies are markers. Anyone can write one from the console,
                and what it buys them is the HTML of a panel whose every request
                the API refuses. The gate decides what gets sent; the API decides
                what is allowed. Do not put a check here and nowhere else.
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
            <p>
              A project scaffolded before v3.220.0 has no staff group, and{' '}
              <code>grit upgrade</code> does not restructure <code>routes.go</code>, so it says so
              instead. Its generated deletes stay ADMIN-only until you add the group. To add it,
              create the group beside the admin one, move each route you want grantable onto it
              with the permission it needs, and hand it to generated resources:
            </p>
            <CodeBlock language="go" filename="internal/routes/routes.go" code={`staff := v1.Group("")
staff.Use(middleware.APIKeyOrAuth(db, middleware.Auth(db, authService)))
staff.Use(middleware.RequireStaff())
{
    staff.GET("/users", middleware.RequireRole("ADMIN", "perm:users.view"), userHandler.List)
}

mountResources(&Mount{
    // ...
    Admin:     admin,
    Staff:     staff,
})`} />
            <p>
              <code>grit upgrade</code> does not rewrite <code>internal/middleware/auth.go</code> or{' '}
              <code>internal/routes/resources.go</code> either. Copy <code>RequireStaff</code> and{' '}
              <code>RequirePermissionFor</code> from a new project into the first, and add{' '}
              <code>Staff *gin.RouterGroup</code> to the <code>Mount</code> in the second. Regenerate a
              resource afterwards and its delete and bulk routes move to the staff group.
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
    </DocsShell>
  )
}
