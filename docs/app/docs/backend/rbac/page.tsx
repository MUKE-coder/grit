import Link from 'next/link'
import { ArrowRight, ArrowLeft } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'
import { CodeBlock } from '@/components/code-block'
import { getDocMetadata } from '@/config/docs-metadata'

export const metadata = getDocMetadata('/docs/backend/rbac')

export default function RBACPage() {
  return (
    <div className="min-h-screen bg-background">
      <SiteHeader />
      <DocsSidebar />

      <main className="lg:pl-64">
        <div className="container max-w-screen-xl py-10 px-6">
          <div className="max-w-3xl">
            {/* Header */}
            <div className="mb-10">
              <span className="tag-mono text-primary/80 mb-3 block">Backend</span>
              <h1 className="text-4xl font-bold tracking-tight mb-4">
                RBAC &amp; Roles
              </h1>
              <p className="text-lg text-muted-foreground leading-relaxed">
                Grit includes a role-based access control (RBAC) system with three built-in roles.
                Routes can be restricted to specific roles using middleware, and the admin panel
                adapts its UI based on the authenticated user&apos;s role.
              </p>
            </div>

            <div className="prose-grit">
              {/* Overview */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Overview
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Every user has a <code>role</code> field stored as an uppercase string. Grit ships
                  with three roles:
                </p>

                <div className="overflow-hidden rounded-xl border border-border/40">
                  <table className="w-full text-sm">
                    <thead>
                      <tr className="border-b border-border/30 bg-accent/20">
                        <th className="px-4 py-3 text-left font-medium text-foreground/80">Role</th>
                        <th className="px-4 py-3 text-left font-medium text-foreground/80">Value</th>
                        <th className="px-4 py-3 text-left font-medium text-foreground/80">Description</th>
                      </tr>
                    </thead>
                    <tbody className="divide-y divide-border/20">
                      <tr>
                        <td className="px-4 py-3 font-mono text-primary">ADMIN</td>
                        <td className="px-4 py-3 font-mono text-muted-foreground">&quot;ADMIN&quot;</td>
                        <td className="px-4 py-3 text-muted-foreground">Full access. Can manage all resources, users, and system settings.</td>
                      </tr>
                      <tr>
                        <td className="px-4 py-3 font-mono text-primary">EDITOR</td>
                        <td className="px-4 py-3 font-mono text-muted-foreground">&quot;EDITOR&quot;</td>
                        <td className="px-4 py-3 text-muted-foreground">Can manage content resources. Sees the full admin panel but cannot delete users or access system pages.</td>
                      </tr>
                      <tr>
                        <td className="px-4 py-3 font-mono text-primary">USER</td>
                        <td className="px-4 py-3 font-mono text-muted-foreground">&quot;USER&quot;</td>
                        <td className="px-4 py-3 text-muted-foreground">Standard user. Sees only their profile page in the admin panel.</td>
                      </tr>
                    </tbody>
                  </table>
                </div>
              </div>

              {/* Role Constants */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Role Constants
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Roles are defined as Go constants in the User model and as TypeScript types in the
                  shared package:
                </p>

                <h3 className="text-lg font-medium mb-3">Go (Backend)</h3>
                <CodeBlock
                  language="go"
                  filename="apps/api/internal/models/user.go"
                  code={`const (
    RoleAdmin  = "ADMIN"
    RoleEditor = "EDITOR"
    RoleUser   = "USER"
)

type User struct {
    ID        uint   \`gorm:"primaryKey" json:"id"\`
    FirstName string \`gorm:"size:100;not null" json:"first_name"\`
    LastName  string \`gorm:"size:100;not null" json:"last_name"\`
    Email     string \`gorm:"uniqueIndex;size:255;not null" json:"email"\`
    Password  string \`gorm:"size:255;not null" json:"-"\`
    Role      string \`gorm:"size:20;default:USER" json:"role"\`
    // ...
}`}
                />

                <h3 className="text-lg font-medium mb-3 mt-6">TypeScript (Frontend)</h3>
                <CodeBlock
                  language="typescript"
                  filename="packages/shared/types/user.ts"
                  code={`export interface User {
  id: number;
  first_name: string;
  last_name: string;
  email: string;
  role: "ADMIN" | "EDITOR" | "USER";
  avatar: string;
  job_title: string;
  bio: string;
  active: boolean;
  // ...
}`}
                />

                <CodeBlock
                  language="typescript"
                  filename="packages/shared/constants/index.ts"
                  code={`export const ROLES = {
  ADMIN: "ADMIN",
  EDITOR: "EDITOR",
  USER: "USER",
} as const;`}
                />
              </div>

              {/* RequireRole Middleware */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  RequireRole Middleware
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  The <code>RequireRole</code> middleware checks the authenticated user&apos;s role
                  against a list of allowed roles. It must be used after the <code>Auth</code> middleware.
                </p>

                <CodeBlock
                  language="go"
                  filename="apps/api/internal/middleware/auth.go"
                  code={`// RequireRole returns middleware that restricts access to specific roles.
func RequireRole(roles ...string) gin.HandlerFunc {
    return func(c *gin.Context) {
        userRole, exists := c.Get("user_role")
        if !exists {
            c.JSON(http.StatusForbidden, gin.H{
                "error": gin.H{
                    "code":    "FORBIDDEN",
                    "message": "Access denied",
                },
            })
            c.Abort()
            return
        }

        role := userRole.(string)
        for _, allowed := range roles {
            if role == allowed {
                c.Next()
                return
            }
        }

        c.JSON(http.StatusForbidden, gin.H{
            "error": gin.H{
                "code":    "FORBIDDEN",
                "message": "Insufficient permissions",
            },
        })
        c.Abort()
    }
}`}
                />

                <p className="text-muted-foreground leading-relaxed mb-4">
                  You can pass one or more roles:
                </p>

                <CodeBlock
                  language="go"
                  code={`// Single role
admin.Use(middleware.RequireRole("ADMIN"))

// Multiple roles
editors.Use(middleware.RequireRole("ADMIN", "EDITOR"))`}
                />
              </div>

              {/* Default Route Groups */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Default Route Groups
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Grit scaffolds three route groups in <code>routes.go</code>:
                </p>

                <div className="overflow-hidden rounded-xl border border-border/40">
                  <table className="w-full text-sm">
                    <thead>
                      <tr className="border-b border-border/30 bg-accent/20">
                        <th className="px-4 py-3 text-left font-medium text-foreground/80">Group</th>
                        <th className="px-4 py-3 text-left font-medium text-foreground/80">Marker</th>
                        <th className="px-4 py-3 text-left font-medium text-foreground/80">Access</th>
                      </tr>
                    </thead>
                    <tbody className="divide-y divide-border/20">
                      <tr>
                        <td className="px-4 py-3 font-medium">Protected</td>
                        <td className="px-4 py-3 font-mono text-xs text-muted-foreground">// grit:routes:protected</td>
                        <td className="px-4 py-3 text-muted-foreground">Any authenticated user</td>
                      </tr>
                      <tr>
                        <td className="px-4 py-3 font-medium">Admin</td>
                        <td className="px-4 py-3 font-mono text-xs text-muted-foreground">// grit:routes:admin</td>
                        <td className="px-4 py-3 text-muted-foreground">ADMIN only</td>
                      </tr>
                      <tr>
                        <td className="px-4 py-3 font-medium">Custom</td>
                        <td className="px-4 py-3 font-mono text-xs text-muted-foreground">// grit:routes:custom</td>
                        <td className="px-4 py-3 text-muted-foreground">Role-restricted (via <code>--roles</code> flag)</td>
                      </tr>
                    </tbody>
                  </table>
                </div>

                <p className="text-muted-foreground leading-relaxed mt-4 mb-4">
                  By default, <code>grit generate resource</code> places CRUD routes in the protected group
                  and DELETE in the admin group. When you use the <code>--roles</code> flag, all routes
                  are placed in the custom group instead.
                </p>

                <CodeBlock
                  language="go"
                  filename="apps/api/internal/routes/routes.go"
                  code={`// Protected routes (any authenticated user)
protected := r.Group("/api")
protected.Use(middleware.Auth(db, authService))
{
    // Profile routes (any authenticated user)
    profile := protected.Group("/profile")
    {
        profile.GET("", userHandler.GetProfile)
        profile.PUT("", userHandler.UpdateProfile)
        profile.DELETE("", userHandler.DeleteProfile)
    }

    // Default resource routes go here
    // grit:routes:protected

    // Admin routes (ADMIN only)
    admin := protected.Group("")
    admin.Use(middleware.RequireRole("ADMIN"))
    {
        // grit:routes:admin
    }

    // Custom role-restricted routes
    // grit:routes:custom
}`}
                />
              </div>

              {/* Using the --roles Flag */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Using the --roles Flag
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  When generating a resource, you can restrict all its routes to specific roles
                  using the <code>--roles</code> flag:
                </p>

                <CodeBlock
                  terminal
                  code='grit generate resource Post --fields "title:string,content:text,published:bool" --roles "ADMIN,EDITOR"'
                />

                <p className="text-muted-foreground leading-relaxed mb-4">
                  This generates a self-contained route group with <code>RequireRole</code> middleware:
                </p>

                <CodeBlock
                  language="go"
                  filename="Generated route group"
                  code={`// Posts routes (restricted to ADMIN, EDITOR)
postsGroup := protected.Group("/posts")
postsGroup.Use(middleware.RequireRole("ADMIN", "EDITOR"))
{
    postsGroup.GET("", postHandler.List)
    postsGroup.GET("/:id", postHandler.GetByID)
    postsGroup.POST("", postHandler.Create)
    postsGroup.PUT("/:id", postHandler.Update)
    postsGroup.DELETE("/:id", postHandler.Delete)
}`}
                />

                <div className="p-4 rounded-lg border border-primary/20 bg-primary/5 mb-6">
                  <p className="text-sm text-foreground/80 mb-0">
                    <strong>Without --roles:</strong> GET, GET/:id, POST, and PUT go in the protected group (any auth user),
                    and DELETE goes in the admin group (ADMIN only).
                    <br />
                    <strong>With --roles:</strong> All five CRUD routes are placed in a single group restricted to
                    the specified roles.
                  </p>
                </div>
              </div>

              {/* Profile Page */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Profile Page
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Users with the <code>USER</code> role are redirected from the dashboard to a dedicated
                  profile page at <code>/profile</code>. The profile page includes:
                </p>
                <ul className="list-disc list-inside text-muted-foreground space-y-1 mb-4">
                  <li><strong>Personal Information</strong> &mdash; Edit first name, last name, and email</li>
                  <li><strong>Professional Information</strong> &mdash; Job title and bio</li>
                  <li><strong>Change Password</strong> &mdash; Update account password</li>
                  <li><strong>Danger Zone</strong> &mdash; Delete account with confirmation dialog</li>
                </ul>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  The sidebar adapts based on role:
                </p>
                <div className="overflow-hidden rounded-xl border border-border/40">
                  <table className="w-full text-sm">
                    <thead>
                      <tr className="border-b border-border/30 bg-accent/20">
                        <th className="px-4 py-3 text-left font-medium text-foreground/80">Role</th>
                        <th className="px-4 py-3 text-left font-medium text-foreground/80">Sidebar Shows</th>
                      </tr>
                    </thead>
                    <tbody className="divide-y divide-border/20">
                      <tr>
                        <td className="px-4 py-3 font-mono text-primary">ADMIN / EDITOR</td>
                        <td className="px-4 py-3 text-muted-foreground">Dashboard, Resources, System pages, Profile</td>
                      </tr>
                      <tr>
                        <td className="px-4 py-3 font-mono text-primary">USER</td>
                        <td className="px-4 py-3 text-muted-foreground">Profile only</td>
                      </tr>
                    </tbody>
                  </table>
                </div>

                <h3 className="text-lg font-medium mb-3 mt-6">Backend Profile Endpoints</h3>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Profile endpoints use the authenticated user&apos;s ID from the JWT context &mdash; no
                  <code>:id</code> parameter needed:
                </p>

                <div className="overflow-hidden rounded-xl border border-border/40">
                  <table className="w-full text-sm">
                    <thead>
                      <tr className="border-b border-border/30 bg-accent/20">
                        <th className="px-4 py-3 text-left font-medium text-foreground/80">Method</th>
                        <th className="px-4 py-3 text-left font-medium text-foreground/80">Endpoint</th>
                        <th className="px-4 py-3 text-left font-medium text-foreground/80">Description</th>
                      </tr>
                    </thead>
                    <tbody className="divide-y divide-border/20">
                      <tr>
                        <td className="px-4 py-3 font-mono text-sm">GET</td>
                        <td className="px-4 py-3 font-mono text-sm">/api/profile</td>
                        <td className="px-4 py-3 text-muted-foreground">Get current user&apos;s profile</td>
                      </tr>
                      <tr>
                        <td className="px-4 py-3 font-mono text-sm">PUT</td>
                        <td className="px-4 py-3 font-mono text-sm">/api/profile</td>
                        <td className="px-4 py-3 text-muted-foreground">Update profile (name, email, password, job title, bio)</td>
                      </tr>
                      <tr>
                        <td className="px-4 py-3 font-mono text-sm">DELETE</td>
                        <td className="px-4 py-3 font-mono text-sm">/api/profile</td>
                        <td className="px-4 py-3 text-muted-foreground">Delete account</td>
                      </tr>
                    </tbody>
                  </table>
                </div>
              </div>

              {/* Adding New Roles */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Adding New Roles
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Use the <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">grit add role</code> command
                  to add a new role across all project files in one step:
                </p>
                <CodeBlock
                  language="bash"
                  code={`grit add role MODERATOR`}
                />
                <p className="text-muted-foreground leading-relaxed mb-4 mt-4">
                  This automatically updates <strong className="text-foreground/90">7 locations</strong>:
                </p>
                <ul className="space-y-2 mb-6">
                  {[
                    'Go model constants (models/user.go)',
                    'Zod schema enum (schemas/user.ts)',
                    'TypeScript union type (types/user.ts)',
                    'ROLES constants object (constants/index.ts)',
                    'Admin badge configuration (resources/users.ts)',
                    'Admin table filter options',
                    'Admin form select options',
                  ].map((item) => (
                    <li key={item} className="flex items-start gap-2.5 text-[14px] text-muted-foreground">
                      <span className="mt-1.5 h-1.5 w-1.5 rounded-full bg-primary/50 shrink-0" />
                      <span>{item}</span>
                    </li>
                  ))}
                </ul>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  After adding a role, you can use it in route restrictions:
                </p>
                <CodeBlock
                  language="go"
                  code={`moderators := protected.Group("/reports")
moderators.Use(middleware.RequireRole("ADMIN", "MODERATOR"))
{
    moderators.GET("", reportHandler.List)
}`}
                />
                <p className="text-sm text-muted-foreground/60 mt-4">
                  You may also want to update the sidebar visibility logic in{' '}
                  <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">components/layout/sidebar.tsx</code>{' '}
                  if the new role should have admin-level navigation access.
                </p>
              </div>

              {/* Navigation */}
              <div className="mt-16 flex items-center justify-between border-t border-border/30 pt-8">
                <Link href="/docs/backend/seeders">
                  <Button variant="outline" size="sm" className="gap-2">
                    <ArrowLeft className="h-4 w-4" />
                    Seeders
                  </Button>
                </Link>
                <Link href="/docs/admin/overview">
                  <Button variant="outline" size="sm" className="gap-2">
                    Admin Overview
                    <ArrowRight className="h-4 w-4" />
                  </Button>
                </Link>
              </div>
            </div>
          </div>
        </div>
      </main>
    </div>
  )
}
