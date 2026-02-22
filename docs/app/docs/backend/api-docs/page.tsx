import Link from 'next/link'
import { ArrowRight, ArrowLeft } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'
import { CodeBlock } from '@/components/code-block'
import { getDocMetadata } from '@/config/docs-metadata'

export const metadata = getDocMetadata('/docs/backend/api-docs')

export default function APIDocsPage() {
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
                API Documentation
              </h1>
              <p className="text-lg text-muted-foreground leading-relaxed">
                Every Grit project ships with auto-generated, interactive API documentation powered by{' '}
                <a href="https://github.com/MUKE-coder/gin-docs" className="text-primary hover:underline" target="_blank" rel="noopener noreferrer">gin-docs</a>.
                No annotations, no comments, no manual spec files &mdash; your docs are generated
                directly from your Gin routes and GORM models.
              </p>
            </div>

            <div className="prose-grit">
              {/* How It Works */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  How It Works
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  When you scaffold a project with <code>grit new</code>, gin-docs is mounted in your
                  routes file with a single function call. It introspects your Gin router and GORM models
                  at startup to generate a complete OpenAPI 3.1 specification.
                </p>

                <div className="rounded-xl border border-border/40 bg-accent/5 p-4 mb-6">
                  <p className="text-sm text-muted-foreground mb-0">
                    <strong className="text-foreground/90">Zero annotations.</strong>{' '}
                    gin-docs reads your route definitions, struct tags (<code>json</code>, <code>binding</code>,{' '}
                    <code>gorm</code>, <code>docs</code>), and handler signatures to generate documentation
                    automatically. No <code>// @Summary</code> comments needed.
                  </p>
                </div>
              </div>

              {/* What You Get */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  What You Get
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  After running your API, these endpoints are available:
                </p>

                <div className="overflow-x-auto mb-6">
                  <table className="w-full text-sm">
                    <thead>
                      <tr className="border-b border-border/30">
                        <th className="text-left py-2 pr-4 text-foreground/80 font-medium">Endpoint</th>
                        <th className="text-left py-2 text-foreground/80 font-medium">Description</th>
                      </tr>
                    </thead>
                    <tbody>
                      <tr className="border-b border-border/20">
                        <td className="py-2 pr-4"><code>GET /docs</code></td>
                        <td className="py-2 text-muted-foreground">Interactive Scalar UI</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="py-2 pr-4"><code>GET /docs?ui=swagger</code></td>
                        <td className="py-2 text-muted-foreground">Swagger UI alternative</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="py-2 pr-4"><code>GET /docs/openapi.json</code></td>
                        <td className="py-2 text-muted-foreground">OpenAPI 3.1 spec (JSON)</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="py-2 pr-4"><code>GET /docs/openapi.yaml</code></td>
                        <td className="py-2 text-muted-foreground">OpenAPI 3.1 spec (YAML)</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="py-2 pr-4"><code>GET /docs/export/postman</code></td>
                        <td className="py-2 text-muted-foreground">Postman Collection v2.1</td>
                      </tr>
                      <tr>
                        <td className="py-2 pr-4"><code>GET /docs/export/insomnia</code></td>
                        <td className="py-2 text-muted-foreground">Insomnia export</td>
                      </tr>
                    </tbody>
                  </table>
                </div>
              </div>

              {/* Configuration */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Configuration
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  The mount call in <code>routes.go</code> configures everything:
                </p>

                <CodeBlock
                  language="go"
                  filename="internal/routes/routes.go"
                  code={`// API Documentation (gin-docs)
gindocs.Mount(r, db, gindocs.Config{
    Title:       cfg.AppName + " API",
    Description: "REST API built with Grit — Go + React meta-framework.",
    Version:     "1.0.0",
    UI:          gindocs.UIScalar,     // or gindocs.UISwagger
    ScalarTheme: "kepler",
    Models:      []interface{}{&models.User{}, &models.Upload{}, &models.Blog{}},
    Auth:        gindocs.AuthBearer,
})`}
                />

                <h3 className="text-lg font-semibold mt-8 mb-3">Config Options</h3>

                <div className="overflow-x-auto mb-6">
                  <table className="w-full text-sm">
                    <thead>
                      <tr className="border-b border-border/30">
                        <th className="text-left py-2 pr-4 text-foreground/80 font-medium">Option</th>
                        <th className="text-left py-2 pr-4 text-foreground/80 font-medium">Default</th>
                        <th className="text-left py-2 text-foreground/80 font-medium">Description</th>
                      </tr>
                    </thead>
                    <tbody>
                      <tr className="border-b border-border/20">
                        <td className="py-2 pr-4"><code>Prefix</code></td>
                        <td className="py-2 pr-4 text-muted-foreground"><code>/docs</code></td>
                        <td className="py-2 text-muted-foreground">URL path for the docs UI</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="py-2 pr-4"><code>Title</code></td>
                        <td className="py-2 pr-4 text-muted-foreground"><code>API Documentation</code></td>
                        <td className="py-2 text-muted-foreground">API title shown in the UI</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="py-2 pr-4"><code>Version</code></td>
                        <td className="py-2 pr-4 text-muted-foreground"><code>1.0.0</code></td>
                        <td className="py-2 text-muted-foreground">API version</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="py-2 pr-4"><code>UI</code></td>
                        <td className="py-2 pr-4 text-muted-foreground"><code>UISwagger</code></td>
                        <td className="py-2 text-muted-foreground"><code>UIScalar</code> or <code>UISwagger</code></td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="py-2 pr-4"><code>ScalarTheme</code></td>
                        <td className="py-2 pr-4 text-muted-foreground"><code>kepler</code></td>
                        <td className="py-2 text-muted-foreground">Scalar UI theme (kepler, moon, purple, saturn, etc.)</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="py-2 pr-4"><code>Models</code></td>
                        <td className="py-2 pr-4 text-muted-foreground"><code>[]</code></td>
                        <td className="py-2 text-muted-foreground">GORM models for schema generation</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="py-2 pr-4"><code>Auth</code></td>
                        <td className="py-2 pr-4 text-muted-foreground"><code>AuthNone</code></td>
                        <td className="py-2 text-muted-foreground"><code>AuthBearer</code> for JWT authentication</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="py-2 pr-4"><code>DevMode</code></td>
                        <td className="py-2 pr-4 text-muted-foreground"><code>false</code></td>
                        <td className="py-2 text-muted-foreground">Regenerate spec on every request</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="py-2 pr-4"><code>ReadOnly</code></td>
                        <td className="py-2 pr-4 text-muted-foreground"><code>false</code></td>
                        <td className="py-2 text-muted-foreground">Disable &ldquo;Try It Out&rdquo; in the UI</td>
                      </tr>
                      <tr>
                        <td className="py-2 pr-4"><code>ExcludeRoutes</code></td>
                        <td className="py-2 pr-4 text-muted-foreground"><code>[]</code></td>
                        <td className="py-2 text-muted-foreground">Glob patterns to exclude routes</td>
                      </tr>
                    </tbody>
                  </table>
                </div>
              </div>

              {/* GORM Model Schemas */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  GORM Model Schemas
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  When you pass GORM models to the <code>Models</code> config, gin-docs automatically
                  generates three schema variants for each model:
                </p>
                <ul className="space-y-2 text-muted-foreground mb-6">
                  <li><strong className="text-foreground/80">Full Model</strong> &mdash; includes ID and timestamps, used for responses</li>
                  <li><strong className="text-foreground/80">Create Variant</strong> &mdash; excludes auto-generated fields (ID, timestamps)</li>
                  <li><strong className="text-foreground/80">Update Variant</strong> &mdash; all fields optional, for PATCH operations</li>
                </ul>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  This means your struct tags drive the documentation:
                </p>

                <CodeBlock
                  language="go"
                  filename="internal/models/user.go"
                  code={`type User struct {
    ID        uint      \`json:"id" gorm:"primarykey"\`
    FirstName string    \`json:"first_name" binding:"required,min=2"\`
    LastName  string    \`json:"last_name" binding:"required,min=2"\`
    Email     string    \`json:"email" binding:"required,email" gorm:"uniqueIndex"\`
    Password  string    \`json:"-"\`
    Role      string    \`json:"role" binding:"oneof=ADMIN EDITOR USER" gorm:"default:'USER'"\`
    Active    bool      \`json:"active" gorm:"default:true"\`
    CreatedAt time.Time \`json:"created_at" gorm:"autoCreateTime"\`
    UpdatedAt time.Time \`json:"updated_at" gorm:"autoUpdateTime"\`
}`}
                />

                <div className="rounded-xl border border-border/40 bg-accent/5 p-4 mt-4 mb-6">
                  <p className="text-sm text-muted-foreground mb-0">
                    <strong className="text-foreground/90">Struct tags matter:</strong>{' '}
                    <code>binding:&quot;required&quot;</code> marks fields as required,{' '}
                    <code>binding:&quot;oneof=A B C&quot;</code> creates enums,{' '}
                    <code>json:&quot;-&quot;</code> hides fields (like passwords),{' '}
                    <code>gorm:&quot;primarykey&quot;</code> and <code>gorm:&quot;autoCreateTime&quot;</code> mark
                    fields as read-only in create/update variants.
                  </p>
                </div>
              </div>

              {/* Route Customization */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Customizing Route Documentation
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  While gin-docs infers most documentation automatically, you can customize
                  individual routes using the fluent builder API:
                </p>

                <CodeBlock
                  language="go"
                  filename="internal/routes/routes.go"
                  code={`docs := gindocs.Mount(r, db, gindocs.Config{...})

// Override a specific route's documentation
docs.Route("POST /api/users").
    Summary("Create a new user").
    Description("Creates a user account with the given details.").
    RequestBody(CreateUserInput{}).
    Response(201, User{}, "User created successfully").
    Tags("Users")`}
                />

                <p className="text-muted-foreground leading-relaxed mt-4 mb-4">
                  Or use inline middleware for per-handler documentation:
                </p>

                <CodeBlock
                  language="go"
                  filename="internal/routes/routes.go"
                  code={`r.POST("/api/users",
    gindocs.Doc(gindocs.DocConfig{
        Summary:     "Create user",
        RequestBody: CreateUserInput{},
        Response:    User{},
    }),
    userHandler.Create,
)`}
                />
              </div>

              {/* Adding Models for New Resources */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Adding Models for New Resources
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  When you generate a new resource with <code>grit generate resource</code>,
                  add the model to the <code>Models</code> slice in the gin-docs config:
                </p>

                <CodeBlock
                  language="go"
                  filename="internal/routes/routes.go"
                  code={`gindocs.Mount(r, db, gindocs.Config{
    // ...
    Models: []interface{}{
        &models.User{},
        &models.Upload{},
        &models.Blog{},
        &models.Product{},   // ← add new models here
        &models.Category{},
    },
})`}
                />
              </div>

              {/* Excluding Routes */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Excluding Routes
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  You can hide internal or admin-only routes from the documentation:
                </p>

                <CodeBlock
                  language="go"
                  code={`gindocs.Mount(r, db, gindocs.Config{
    // ...
    ExcludeRoutes:    []string{"/api/admin/*"},
    ExcludePrefixes:  []string{"/studio", "/sentinel"},
})`}
                />
              </div>

              {/* Switching UI */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Switching Between UIs
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Grit defaults to the Scalar UI (modern, dark theme), but you can switch
                  to Swagger UI in config:
                </p>

                <CodeBlock
                  language="go"
                  code={`// Use Swagger UI instead
gindocs.Config{
    UI: gindocs.UISwagger,
}`}
                />

                <p className="text-muted-foreground leading-relaxed mt-4">
                  You can also switch at runtime by adding a query parameter:{' '}
                  <code>/docs?ui=swagger</code> or <code>/docs?ui=scalar</code>.
                </p>
              </div>

              {/* Exporting */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Exporting Specs
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Use the export endpoints to import your API into other tools:
                </p>
                <ul className="space-y-2 text-muted-foreground mb-6">
                  <li><code>GET /docs/openapi.json</code> &mdash; for any OpenAPI-compatible tool</li>
                  <li><code>GET /docs/openapi.yaml</code> &mdash; YAML format</li>
                  <li><code>GET /docs/export/postman</code> &mdash; import directly into Postman</li>
                  <li><code>GET /docs/export/insomnia</code> &mdash; import directly into Insomnia</li>
                </ul>
              </div>
            </div>

            {/* Navigation */}
            <div className="flex items-center justify-between pt-6 border-t border-border/30">
              <Button variant="ghost" size="sm" asChild>
                <Link href="/docs/backend/rbac">
                  <ArrowLeft className="h-3.5 w-3.5 mr-1" />
                  RBAC &amp; Roles
                </Link>
              </Button>
              <Button variant="ghost" size="sm" asChild>
                <Link href="/docs/admin/overview">
                  Admin Overview
                  <ArrowRight className="h-3.5 w-3.5 ml-1" />
                </Link>
              </Button>
            </div>
          </div>
        </div>
      </main>
    </div>
  )
}
