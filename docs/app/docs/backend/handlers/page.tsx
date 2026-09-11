import Link from 'next/link'
import { ArrowLeft, ArrowRight } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'
import { CodeBlock } from '@/components/code-block'
import { LaneFlow } from '@/components/lane-flow'
import { getDocMetadata } from '@/config/docs-metadata'

export const metadata = getDocMetadata('/docs/backend/handlers')

export default function HandlersPage() {
  return (
    <div className="min-h-screen bg-background isolate">
      <SiteHeader />
      <DocsSidebar />
      <main className="lg:pl-64">
        <div className="container max-w-screen-xl py-10 px-6">
          <div className="max-w-3xl">
            {/* Header */}
            <div className="mb-10">
              <span className="tag-mono text-primary/80 mb-3 block">Backend (Go API)</span>
              <h1 className="text-4xl font-bold tracking-tight mb-4">
                Handlers
              </h1>
              <p className="text-lg text-muted-foreground leading-relaxed">
                Handlers are the HTTP layer of your Grit API. A generated handler reads the request,
                calls its <Link href="/docs/backend/services">service</Link>, and writes the answer.
                It runs no query of its own: every read and write belongs to the service, so the same
                logic is there for a background job, a command or a test that has no request at all.
              </p>
            </div>

            <div className="prose-grit">
              {/* ── Request Lifecycle Diagram ─────────────────────── */}
              <p>
                Every request walks the same path down through the Go layers and returns back up as
                a JSON response:
              </p>
              <LaneFlow
                id="lifecycle"
                lanes={['Client', 'Go API', 'Data']}
                groups={[{ lane: 1, rows: [0, 4], label: 'Request pipeline', tone: 'primary' }]}
                nodes={[
                  { id: 'req', lane: 0, row: 0, title: 'HTTP Request', sub: 'GET /api/v1/posts', tone: 'blue', badge: 1 },
                  { id: 'router', lane: 1, row: 0, title: 'Gin Router', sub: 'matches route', tone: 'primary', badge: 2 },
                  { id: 'mw', lane: 1, row: 1, title: 'Middleware', sub: 'CORS · Auth · Log', tone: 'primary', badge: 3 },
                  { id: 'handler', lane: 1, row: 2, title: 'Handler', sub: 'bind · respond', tone: 'primary', badge: 4 },
                  { id: 'service', lane: 1, row: 3, title: 'Service', sub: 'every query', tone: 'primary', badge: 5 },
                  { id: 'gorm', lane: 1, row: 4, title: 'GORM Model', sub: 'query builder', tone: 'primary', badge: 6 },
                  { id: 'pg', lane: 2, row: 4, title: 'PostgreSQL', sub: ':5434', tone: 'green' },
                  { id: 'resp', lane: 0, row: 5, title: 'JSON response', sub: '{ data, meta }', tone: 'blue', badge: 7 },
                ]}
                edges={[
                  { from: 'req', to: 'router', label: 'route', tone: 'blue' },
                  { from: 'router', to: 'mw', tone: 'primary' },
                  { from: 'mw', to: 'handler', tone: 'primary' },
                  { from: 'handler', to: 'service', label: 'ctx + input', tone: 'primary' },
                  { from: 'service', to: 'gorm', tone: 'primary' },
                  { from: 'gorm', to: 'pg', label: 'SQL', tone: 'green' },
                  { from: 'pg', to: 'resp', label: 'JSON', dashed: true, tone: 'blue' },
                ]}
                legend={[
                  { tone: 'blue', label: 'HTTP' },
                  { tone: 'primary', label: 'Go layers' },
                  { tone: 'green', label: 'Data' },
                ]}
                caption="The handler talks HTTP; the service talks to the database"
              />

              {/* ── Handler Pattern ─────────────────────────────── */}
              <h2 id="handler-pattern">Handler Pattern</h2>
              <p>
                A generated handler is a struct with the handler&apos;s dependencies and three small
                helpers. Its other methods are the endpoints.
              </p>
              <CodeBlock language="go" filename="apps/api/internal/handlers/post.go" code={`// PostHandler serves the post endpoints. It reads the request, asks
// services.PostService, and writes the answer; it runs no query of its
// own, so everything a route does is also available to a job or a test.
type PostHandler struct {
    DB *gorm.DB
}

// service is the post service over this handler's database.
func (h *PostHandler) service() *services.PostService {
    return &services.PostService{DB: h.DB}
}

// ctx is the request's context with the caller on it, which the service
// scopes owned rows by. The organization the multitenant plugin resolved and
// the cancellation when the client goes away travel with it.
func (h *PostHandler) ctx(c *gin.Context) context.Context {
    return authz.WithActor(c.Request.Context(), authz.ActorOf(c))
}

// fail answers an error from the service: a version conflict with the version
// the record is at, a missing row with 404, a broken rule with 422 and its
// message, and anything else as an opaque 500 that is logged.
func (h *PostHandler) fail(c *gin.Context, err error, fallback string) {
    var conflict *concurrency.ErrConflict
    switch {
    case errors.As(err, &conflict):
        concurrency.WriteConflict(c, conflict.Current)
    case errors.Is(err, gorm.ErrRecordNotFound):
        c.JSON(http.StatusNotFound, gin.H{
            "error": gin.H{"code": "NOT_FOUND", "message": "Post not found"},
        })
    default:
        respond.WriteError(c, err, fallback)
    }
}`} />
              <p>
                The handler keeps its <code>DB</code> field and builds the service from it on each
                call, so the routes file constructs it exactly as it always has. The generated routes
                (excerpt):
              </p>
              <CodeBlock language="go" filename="apps/api/internal/routes/post_routes.go" code={`m.Protected.GET("/posts", h.List)
m.Protected.GET("/posts/export", h.Export)
m.Protected.GET("/posts/:id", h.GetByID)
m.Protected.POST("/posts", h.Create)
m.Protected.PUT("/posts/:id", h.Update)
m.Protected.PATCH("/posts/:id", h.Patch)
m.Staff.DELETE("/posts/:id", middleware.RequireRole("ADMIN", "perm:posts.delete"), h.Delete)
m.Staff.POST("/posts/bulk", middleware.RequireRole("ADMIN", "perm:posts.delete"), h.Bulk)`} />

              {/* ── What stays in the handler ─────────────────────────────── */}
              <h2 id="what-stays">What Stays in the Handler</h2>
              <p>
                Everything that is about HTTP, and nothing that is about the data:
              </p>
              <div className="rounded-lg border border-border/30 bg-card/30 overflow-hidden mb-6">
                <table className="w-full text-sm">
                  <thead>
                    <tr className="border-b border-border/30 bg-accent/20">
                      <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Handler</th>
                      <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Service</th>
                    </tr>
                  </thead>
                  <tbody className="text-muted-foreground">
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5">Binding and validating the body</td>
                      <td className="px-4 py-2.5">Every query, and every transaction</td>
                    </tr>
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5">Turning the request into a row or an update map</td>
                      <td className="px-4 py-2.5">What may be searched, sorted, filtered and patched</td>
                    </tr>
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5">Reading <code>If-Match</code>, writing <code>ETag</code></td>
                      <td className="px-4 py-2.5">Refusing a write against a stale version</td>
                    </tr>
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5">Putting the caller on the context</td>
                      <td className="px-4 py-2.5">Scoping owned rows to that caller</td>
                    </tr>
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5">Mapping an error to a status code</td>
                      <td className="px-4 py-2.5">Returning the error: not found, a conflict, a broken rule</td>
                    </tr>
                    <tr>
                      <td className="px-4 py-2.5">Rendering: JSON, CSV, XLSX, PDF; emitting the activity event</td>
                      <td className="px-4 py-2.5">Handing over rows, a batch at a time for an export</td>
                    </tr>
                  </tbody>
                </table>
              </div>

              {/* ── Request Binding ─────────────────────────────── */}
              <h2 id="request-binding">Request Binding with Gin</h2>
              <p>
                Gin&apos;s <code>ShouldBindJSON</code> parses the body into a struct and validates it
                with its <code>binding</code> tags. The generated <code>Create</code> binds, builds the
                row, and hands it to the service:
              </p>
              <CodeBlock language="go" filename="handlers/post.go (Create)" code={`// CreatePostRequest is the JSON body accepted by POST /posts.
type CreatePostRequest struct {
    Title     string \`json:"title" binding:"required"\`
    Body      string \`json:"body"\`
    Published bool   \`json:"published"\`
}

func (h *PostHandler) Create(c *gin.Context) {
    var req CreatePostRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusUnprocessableEntity, gin.H{
            "error": gin.H{"code": "VALIDATION_ERROR", "message": err.Error()},
        })
        return
    }

    item := models.Post{
        Title:     req.Title,
        Body:      req.Body,
        Published: req.Published,
    }

    if err := h.service().Create(h.ctx(c), &item); err != nil {
        h.fail(c, err, "Failed to create post")
        return
    }

    events.Emitted(c, "posts", "Post", "created", item.ID, item.Title, "", nil, item)

    c.JSON(http.StatusCreated, gin.H{
        "data":    item,
        "message": "Post created successfully",
    })
}`} />
              <p>
                Request structs are named types, not anonymous ones, so the API reference can reflect
                over them and document the body.
              </p>

              {/* ── Validation Tags ─────────────────────────────── */}
              <h2 id="validation">Validation with Binding Tags</h2>
              <p>
                Gin uses the <code>go-playground/validator</code> library under the hood.
                Here are the most commonly used binding tags:
              </p>
              <div className="rounded-lg border border-border/30 bg-card/30 overflow-hidden mb-6">
                <table className="w-full text-sm">
                  <thead>
                    <tr className="border-b border-border/30 bg-accent/20">
                      <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Tag</th>
                      <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Description</th>
                      <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Example</th>
                    </tr>
                  </thead>
                  <tbody className="text-muted-foreground">
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5 font-mono text-xs">required</td>
                      <td className="px-4 py-2.5">Field must be present and non-zero</td>
                      <td className="px-4 py-2.5 font-mono text-xs">{`binding:"required"`}</td>
                    </tr>
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5 font-mono text-xs">email</td>
                      <td className="px-4 py-2.5">Must be a valid email address</td>
                      <td className="px-4 py-2.5 font-mono text-xs">{`binding:"required,email"`}</td>
                    </tr>
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5 font-mono text-xs">min=N</td>
                      <td className="px-4 py-2.5">Minimum length (string) or value (number)</td>
                      <td className="px-4 py-2.5 font-mono text-xs">{`binding:"min=3"`}</td>
                    </tr>
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5 font-mono text-xs">max=N</td>
                      <td className="px-4 py-2.5">Maximum length (string) or value (number)</td>
                      <td className="px-4 py-2.5 font-mono text-xs">{`binding:"max=255"`}</td>
                    </tr>
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5 font-mono text-xs">gt=N</td>
                      <td className="px-4 py-2.5">Greater than (for numbers)</td>
                      <td className="px-4 py-2.5 font-mono text-xs">{`binding:"gt=0"`}</td>
                    </tr>
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5 font-mono text-xs">oneof=a b c</td>
                      <td className="px-4 py-2.5">Must be one of the listed values</td>
                      <td className="px-4 py-2.5 font-mono text-xs">{`binding:"oneof=admin editor user"`}</td>
                    </tr>
                    <tr>
                      <td className="px-4 py-2.5 font-mono text-xs">url</td>
                      <td className="px-4 py-2.5">Must be a valid URL</td>
                      <td className="px-4 py-2.5 font-mono text-xs">{`binding:"url"`}</td>
                    </tr>
                  </tbody>
                </table>
              </div>

              {/* ── Lists ─────────────────────────────── */}
              <h2 id="search-sort-filter">Pagination, Search, Sort and Filter</h2>
              <p>
                The handler reads the query string with <code>paginate.Bind</code> and passes it on.
                Which columns a client may search, sort and filter by is the service&apos;s decision,
                whitelisted in its list config, because each name ends up in SQL.
              </p>
              <CodeBlock language="go" filename="handlers/post.go (List)" code={`func (h *PostHandler) List(c *gin.Context) {
    res, err := h.service().List(h.ctx(c), paginate.Bind(c), c.Query("archived"))
    if err != nil {
        h.fail(c, err, "Failed to fetch posts")
        return
    }
    c.JSON(http.StatusOK, res) // { data, meta }
}`} />
              <div className="rounded-lg border border-border/30 bg-card/30 overflow-hidden mb-6">
                <table className="w-full text-sm">
                  <thead>
                    <tr className="border-b border-border/30 bg-accent/20">
                      <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Query Param</th>
                      <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Description</th>
                    </tr>
                  </thead>
                  <tbody className="text-muted-foreground">
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5 font-mono text-xs">page, page_size</td>
                      <td className="px-4 py-2.5">Page number (1-based) and rows per page, clamped</td>
                    </tr>
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5 font-mono text-xs">search</td>
                      <td className="px-4 py-2.5">Matched against the service&apos;s <code>Searchable</code> columns</td>
                    </tr>
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5 font-mono text-xs">sort_by, sort_order</td>
                      <td className="px-4 py-2.5">A column in <code>Sortable</code>, and asc or desc</td>
                    </tr>
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5 font-mono text-xs">any column, e.g. status=paid</td>
                      <td className="px-4 py-2.5">An equality filter, if the column is in <code>Filterable</code></td>
                    </tr>
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5 font-mono text-xs">created_from, created_to</td>
                      <td className="px-4 py-2.5">A date window, both ends inclusive</td>
                    </tr>
                    <tr>
                      <td className="px-4 py-2.5 font-mono text-xs">archived</td>
                      <td className="px-4 py-2.5"><code>true</code> for archived rows only, <code>all</code> for both</td>
                    </tr>
                  </tbody>
                </table>
              </div>

              {/* ── Full CRUD Handler ─────────────────────────────── */}
              <h2 id="crud-handler">Reads and Writes</h2>

              <h3 id="get-by-id">GetByID</h3>
              <CodeBlock language="go" filename="handlers/post.go (GetByID)" code={`func (h *PostHandler) GetByID(c *gin.Context) {
    item, err := h.service().GetByID(h.ctx(c), c.Param("id"))
    if err != nil {
        h.fail(c, err, "Failed to load post")
        return
    }

    // The version to send back as If-Match when saving.
    c.Header("ETag", concurrency.Tag(item.Version))
    c.JSON(http.StatusOK, gin.H{"data": item})
}`} />

              <h3 id="update">Update</h3>
              <p>
                The handler turns the typed request into the columns to write, and passes on the
                version the client read. The service writes only if the row is still at that version.
              </p>
              <CodeBlock language="go" filename="handlers/post.go (Update)" code={`// UpdatePostRequest is the JSON body accepted by PUT /posts/:id.
// Every field is optional: only what the client sends is applied.
type UpdatePostRequest struct {
    Title     string \`json:"title"\`
    Body      string \`json:"body"\`
    Published *bool  \`json:"published"\`
}

func (h *PostHandler) Update(c *gin.Context) {
    var req UpdatePostRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusUnprocessableEntity, gin.H{
            "error": gin.H{"code": "VALIDATION_ERROR", "message": err.Error()},
        })
        return
    }

    updates := map[string]interface{}{}
    if req.Title != "" {
        updates["title"] = req.Title
    }
    if req.Body != "" {
        updates["body"] = req.Body
    }
    if req.Published != nil {
        updates["published"] = *req.Published
    }

    item, err := h.service().Update(h.ctx(c), c.Param("id"), updates, concurrency.FromRequest(c))
    if err != nil {
        h.fail(c, err, "Failed to update post") // 409 on a stale If-Match
        return
    }

    c.Header("ETag", concurrency.Tag(item.Version))
    c.JSON(http.StatusOK, gin.H{
        "data":    item,
        "message": "Post updated successfully",
    })
}`} />
              <p>
                Notice the <strong>pointer</strong> for the boolean (<code>*bool</code>). It tells
                &quot;not sent&quot; (<code>nil</code>) from &quot;sent as false&quot;; without it,
                Go&apos;s zero value would overwrite the field on every save.
              </p>

              <h3 id="delete">Delete</h3>
              <CodeBlock language="go" filename="handlers/post.go (Delete)" code={`func (h *PostHandler) Delete(c *gin.Context) {
    item, err := h.service().Delete(h.ctx(c), c.Param("id"))
    if err != nil {
        h.fail(c, err, "Failed to delete post")
        return
    }

    events.Emitted(c, "posts", "Post", "deleted", item.ID, item.Title, "", item, nil)
    c.JSON(http.StatusOK, gin.H{"message": "Post deleted successfully"})
}`} />
              <p>
                Models carry <code>gorm.DeletedAt</code>, so this is a <strong>soft delete</strong>:
                the row stays with a <code>deleted_at</code> timestamp and drops out of every query.
              </p>

              {/* ── Errors ─────────────────────────────── */}
              <h2 id="errors">Errors</h2>
              <p>
                The service returns Go errors and never an HTTP status. <code>fail</code> decides the
                status:
              </p>
              <div className="rounded-lg border border-border/30 bg-card/30 overflow-hidden mb-6">
                <table className="w-full text-sm">
                  <thead>
                    <tr className="border-b border-border/30 bg-accent/20">
                      <th className="text-left px-4 py-2.5 font-medium text-foreground/80">The service returns</th>
                      <th className="text-left px-4 py-2.5 font-medium text-foreground/80">The client gets</th>
                    </tr>
                  </thead>
                  <tbody className="text-muted-foreground">
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5 font-mono text-xs">*concurrency.ErrConflict</td>
                      <td className="px-4 py-2.5">409 <code>VERSION_CONFLICT</code>, with the version the row is at</td>
                    </tr>
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5 font-mono text-xs">gorm.ErrRecordNotFound</td>
                      <td className="px-4 py-2.5">404, also for somebody else&apos;s row on an owned resource</td>
                    </tr>
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5 font-mono text-xs">respond.Rule(&quot;...&quot;)</td>
                      <td className="px-4 py-2.5">422 with that message</td>
                    </tr>
                    <tr>
                      <td className="px-4 py-2.5 font-mono text-xs">anything else</td>
                      <td className="px-4 py-2.5">500 with the fallback message; the error is logged, not sent</td>
                    </tr>
                  </tbody>
                </table>
              </div>

              {/* ── Existing projects ─────────────────────────────── */}
              <h2 id="existing-projects">Resources Generated Before v3.224.0</h2>
              <p>
                Until v3.224.0 the generated handler ran its own queries, and the service beside it
                was never called. <code>grit upgrade</code> does not rewrite your API code, so a
                resource generated before then keeps that handler until you regenerate it. The
                CSV import, the public read endpoints and the tree endpoints still query from their
                own handler files; they move to the service next.
              </p>

              {/* ── Best Practices ─────────────────────────────── */}
              <h2 id="best-practices">Best Practices</h2>
              <ul>
                <li>
                  <strong>No query in a handler.</strong> If a handler needs data, it asks a service
                  method for it. A rule written in a handler is a rule a job can skip.
                </li>
                <li>
                  <strong>Pass <code>h.ctx(c)</code>, not <code>context.Background()</code>.</strong> The
                  caller, the organization and the cancellation are on the request&apos;s context, and a
                  service given anything else cannot see them.
                </li>
                <li>
                  <strong>Use pointers for optional fields</strong> in update requests, so &quot;not
                  provided&quot; is not &quot;set to zero&quot;.
                </li>
                <li>
                  <strong>Send errors through <code>fail</code>.</strong> It keeps the standard error
                  envelope, and it keeps database errors out of responses.
                </li>
              </ul>
            </div>

            {/* Nav */}
            <div className="flex items-center justify-between pt-6 mt-10 border-t border-border/30">
              <Button variant="ghost" size="sm" asChild className="text-muted-foreground/60 hover:text-foreground">
                <Link href="/docs/backend/models" className="gap-1.5">
                  <ArrowLeft className="h-3.5 w-3.5" />
                  Models & Database
                </Link>
              </Button>
              <Button variant="ghost" size="sm" asChild className="text-muted-foreground/60 hover:text-foreground">
                <Link href="/docs/backend/services" className="gap-1.5">
                  Services
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
