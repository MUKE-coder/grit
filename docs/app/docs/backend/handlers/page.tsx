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
                Handlers are the HTTP layer of your Grit API. They receive requests, validate input,
                call services or the database, and return JSON responses. Handlers should stay thin --
                delegate business logic to services.
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
                  { id: 'req', lane: 0, row: 0, title: 'HTTP Request', sub: 'GET /api/posts', tone: 'blue', badge: 1 },
                  { id: 'router', lane: 1, row: 0, title: 'Gin Router', sub: 'matches route', tone: 'primary', badge: 2 },
                  { id: 'mw', lane: 1, row: 1, title: 'Middleware', sub: 'CORS · Auth · Log', tone: 'primary', badge: 3 },
                  { id: 'handler', lane: 1, row: 2, title: 'Handler', sub: 'thin HTTP layer', tone: 'primary', badge: 4 },
                  { id: 'service', lane: 1, row: 3, title: 'Service', sub: 'business logic', tone: 'primary', badge: 5 },
                  { id: 'gorm', lane: 1, row: 4, title: 'GORM Model', sub: 'query builder', tone: 'primary', badge: 6 },
                  { id: 'pg', lane: 2, row: 4, title: 'PostgreSQL', sub: ':5434', tone: 'green' },
                  { id: 'resp', lane: 0, row: 5, title: 'JSON response', sub: '{ data, meta }', tone: 'blue', badge: 7 },
                ]}
                edges={[
                  { from: 'req', to: 'router', label: 'route', tone: 'blue' },
                  { from: 'router', to: 'mw', tone: 'primary' },
                  { from: 'mw', to: 'handler', tone: 'primary' },
                  { from: 'handler', to: 'service', label: 'calls', tone: 'primary' },
                  { from: 'service', to: 'gorm', tone: 'primary' },
                  { from: 'gorm', to: 'pg', label: 'SQL', tone: 'green' },
                  { from: 'pg', to: 'resp', label: 'JSON', dashed: true, tone: 'blue' },
                ]}
                legend={[
                  { tone: 'blue', label: 'HTTP' },
                  { tone: 'primary', label: 'Go layers' },
                  { tone: 'green', label: 'Data' },
                ]}
                caption="Request flows down the Go layers; the JSON response bubbles back up"
              />

              {/* ── Handler Pattern ─────────────────────────────── */}
              <h2 id="handler-pattern">Handler Pattern</h2>
              <p>
                Every handler in Grit is a struct with a <code>DB</code> field (and optionally other
                dependencies). Methods on the struct correspond to HTTP endpoints.
              </p>
              <CodeBlock language="go" filename="apps/api/internal/handlers/post.go" code={`package handlers

import (
    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
)

// PostHandler handles post-related endpoints.
type PostHandler struct {
    DB *gorm.DB
}

// Each method on the struct is a handler function:
// func (h *PostHandler) List(c *gin.Context)    { ... }
// func (h *PostHandler) GetByID(c *gin.Context) { ... }
// func (h *PostHandler) Create(c *gin.Context)  { ... }
// func (h *PostHandler) Update(c *gin.Context)  { ... }
// func (h *PostHandler) Delete(c *gin.Context)  { ... }`} />
              <p>
                Handlers are instantiated in <code>routes/routes.go</code> and wired to their endpoints:
              </p>
              <CodeBlock language="go" filename="apps/api/internal/routes/routes.go" code={`postHandler := &handlers.PostHandler{DB: db}

// Wire to routes
protected.GET("/posts", postHandler.List)
protected.GET("/posts/:id", postHandler.GetByID)
protected.POST("/posts", postHandler.Create)
protected.PUT("/posts/:id", postHandler.Update)
protected.DELETE("/posts/:id", postHandler.Delete)`} />

              {/* ── Request Binding ─────────────────────────────── */}
              <h2 id="request-binding">Request Binding with Gin</h2>
              <p>
                Gin&apos;s <code>ShouldBindJSON</code> method parses the request body into a Go struct
                and validates it using <code>binding</code> struct tags. If validation fails, it
                returns an error that you can send back to the client.
              </p>
              <CodeBlock language="go" filename="request_binding.go" code={`type createPostRequest struct {
    Title     string \`json:"title" binding:"required,min=3,max=255"\`
    Body      string \`json:"body" binding:"required"\`
    Published bool   \`json:"published"\`
}

func (h *PostHandler) Create(c *gin.Context) {
    var req createPostRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusUnprocessableEntity, gin.H{
            "error": gin.H{
                "code":    "VALIDATION_ERROR",
                "message": err.Error(),
            },
        })
        return
    }

    // req is now validated and ready to use
}`} />
              <p>
                Define request structs as private types (lowercase first letter) inside the handler
                file. This keeps them close to the handler that uses them and prevents external access.
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

              {/* ── Pagination Pattern ─────────────────────────────── */}
              <h2 id="pagination">Pagination Pattern</h2>
              <p>
                All list endpoints in Grit use a consistent pagination pattern. Query parameters
                control the page number and page size, and the response includes a <code>meta</code> object
                with pagination details.
              </p>
              <CodeBlock language="go" filename="pagination.go" code={`func (h *PostHandler) List(c *gin.Context) {
    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

    // Clamp values
    if page < 1 {
        page = 1
    }
    if pageSize < 1 || pageSize > 100 {
        pageSize = 20
    }

    query := h.DB.Model(&models.Post{})

    // Count total records
    var total int64
    query.Count(&total)

    // Fetch paginated results
    var posts []models.Post
    offset := (page - 1) * pageSize
    query.Offset(offset).Limit(pageSize).Find(&posts)

    pages := int(math.Ceil(float64(total) / float64(pageSize)))

    c.JSON(http.StatusOK, gin.H{
        "data": posts,
        "meta": gin.H{
            "total":     total,
            "page":      page,
            "page_size": pageSize,
            "pages":     pages,
        },
    })
}`} />
              <p>The client calls the endpoint like this:</p>
              <CodeBlock terminal code="GET /api/posts?page=2&page_size=10" />

              {/* ── Search, Sort & Filter ─────────────────────────────── */}
              <h2 id="search-sort-filter">Search, Sort & Filter</h2>
              <p>
                The Grit handler pattern supports search, sort, and filter out of the box.
                The built-in <code>UserHandler.List</code> demonstrates the full pattern:
              </p>
              <CodeBlock language="go" filename="search_sort_filter.go" code={`func (h *UserHandler) List(c *gin.Context) {
    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
    search := c.Query("search")
    sortBy := c.DefaultQuery("sort_by", "created_at")
    sortOrder := c.DefaultQuery("sort_order", "desc")

    // Validate sort order (prevent SQL injection)
    if sortOrder != "asc" && sortOrder != "desc" {
        sortOrder = "desc"
    }

    // Whitelist allowed sort columns
    allowedSorts := map[string]bool{
        "id": true, "first_name": true, "last_name": true,
        "email": true, "role": true, "created_at": true,
    }
    if !allowedSorts[sortBy] {
        sortBy = "created_at"
    }

    query := h.DB.Model(&models.User{})

    // Search across multiple fields
    if search != "" {
        query = query.Where(
            "first_name ILIKE ? OR last_name ILIKE ? OR email ILIKE ?",
            "%"+search+"%", "%"+search+"%", "%"+search+"%",
        )
    }

    // Count then paginate
    var total int64
    query.Count(&total)

    var users []models.User
    offset := (page - 1) * pageSize
    query.Order(sortBy + " " + sortOrder).
        Offset(offset).
        Limit(pageSize).
        Find(&users)

    pages := int(math.Ceil(float64(total) / float64(pageSize)))

    c.JSON(http.StatusOK, gin.H{
        "data": users,
        "meta": gin.H{
            "total":     total,
            "page":      page,
            "page_size": pageSize,
            "pages":     pages,
        },
    })
}`} />

              <div className="rounded-lg border border-border/30 bg-card/30 overflow-hidden mb-6">
                <table className="w-full text-sm">
                  <thead>
                    <tr className="border-b border-border/30 bg-accent/20">
                      <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Query Param</th>
                      <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Default</th>
                      <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Description</th>
                    </tr>
                  </thead>
                  <tbody className="text-muted-foreground">
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5 font-mono text-xs">page</td>
                      <td className="px-4 py-2.5 font-mono text-xs">1</td>
                      <td className="px-4 py-2.5">Page number (1-based)</td>
                    </tr>
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5 font-mono text-xs">page_size</td>
                      <td className="px-4 py-2.5 font-mono text-xs">20</td>
                      <td className="px-4 py-2.5">Records per page (max 100)</td>
                    </tr>
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5 font-mono text-xs">search</td>
                      <td className="px-4 py-2.5 font-mono text-xs">(empty)</td>
                      <td className="px-4 py-2.5">Full-text search across first_name + last_name + email</td>
                    </tr>
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5 font-mono text-xs">sort_by</td>
                      <td className="px-4 py-2.5 font-mono text-xs">created_at</td>
                      <td className="px-4 py-2.5">Column to sort by (whitelisted)</td>
                    </tr>
                    <tr>
                      <td className="px-4 py-2.5 font-mono text-xs">sort_order</td>
                      <td className="px-4 py-2.5 font-mono text-xs">desc</td>
                      <td className="px-4 py-2.5">Sort direction: asc or desc</td>
                    </tr>
                  </tbody>
                </table>
              </div>

              {/* ── Full CRUD Handler ─────────────────────────────── */}
              <h2 id="crud-handler">Full CRUD Handler</h2>
              <p>
                Here is a complete handler with all five CRUD operations. This is the pattern that
                <code>grit generate resource</code> produces for every new resource.
              </p>

              <h3 id="create">Create</h3>
              <CodeBlock language="go" filename="handlers/post.go -- Create" code={`type createPostRequest struct {
    Title     string \`json:"title" binding:"required,min=3,max=255"\`
    Body      string \`json:"body" binding:"required"\`
    Published bool   \`json:"published"\`
}

func (h *PostHandler) Create(c *gin.Context) {
    var req createPostRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusUnprocessableEntity, gin.H{
            "error": gin.H{
                "code":    "VALIDATION_ERROR",
                "message": err.Error(),
            },
        })
        return
    }

    // Get authenticated user from context
    userID, _ := c.Get("user_id")

    post := models.Post{
        Title:     req.Title,
        Body:      req.Body,
        Published: req.Published,
        AuthorID:  userID.(string),
    }

    if err := h.DB.Create(&post).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": gin.H{
                "code":    "INTERNAL_ERROR",
                "message": "Failed to create post",
            },
        })
        return
    }

    c.JSON(http.StatusCreated, gin.H{
        "data":    post,
        "message": "Post created successfully",
    })
}`} />

              <h3 id="get-by-id">GetByID</h3>
              <CodeBlock language="go" filename="handlers/post.go -- GetByID" code={`func (h *PostHandler) GetByID(c *gin.Context) {
    id := c.Param("id")

    var post models.Post
    if err := h.DB.Preload("Author").First(&post, id).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{
            "error": gin.H{
                "code":    "NOT_FOUND",
                "message": "Post not found",
            },
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "data": post,
    })
}`} />

              <h3 id="update">Update</h3>
              <CodeBlock language="go" filename="handlers/post.go -- Update" code={`func (h *PostHandler) Update(c *gin.Context) {
    id := c.Param("id")

    var post models.Post
    if err := h.DB.First(&post, id).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{
            "error": gin.H{
                "code":    "NOT_FOUND",
                "message": "Post not found",
            },
        })
        return
    }

    var req struct {
        Title     string \`json:"title"\`
        Body      string \`json:"body"\`
        Published *bool  \`json:"published"\`
    }

    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusUnprocessableEntity, gin.H{
            "error": gin.H{
                "code":    "VALIDATION_ERROR",
                "message": err.Error(),
            },
        })
        return
    }

    // Build updates map (only include non-zero fields)
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

    if err := h.DB.Model(&post).Updates(updates).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": gin.H{
                "code":    "INTERNAL_ERROR",
                "message": "Failed to update post",
            },
        })
        return
    }

    h.DB.First(&post, id) // reload

    c.JSON(http.StatusOK, gin.H{
        "data":    post,
        "message": "Post updated successfully",
    })
}`} />
              <p>
                Notice the use of a <strong>pointer</strong> for boolean fields (<code>*bool</code>).
                This lets you distinguish between &quot;not sent&quot; (<code>nil</code>) and &quot;sent as false&quot;.
                Without the pointer, Go&apos;s zero value (<code>false</code>) would always overwrite the field.
              </p>

              <h3 id="delete">Delete</h3>
              <CodeBlock language="go" filename="handlers/post.go -- Delete" code={`func (h *PostHandler) Delete(c *gin.Context) {
    id := c.Param("id")

    var post models.Post
    if err := h.DB.First(&post, id).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{
            "error": gin.H{
                "code":    "NOT_FOUND",
                "message": "Post not found",
            },
        })
        return
    }

    if err := h.DB.Delete(&post).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": gin.H{
                "code":    "INTERNAL_ERROR",
                "message": "Failed to delete post",
            },
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "message": "Post deleted successfully",
    })
}`} />
              <p>
                Because the <code>Post</code> model includes <code>gorm.DeletedAt</code>, calling
                <code>db.Delete()</code> performs a <strong>soft delete</strong>. The row remains in the
                database with a <code>deleted_at</code> timestamp, but is excluded from all future queries.
              </p>

              {/* ── Best Practices ─────────────────────────────── */}
              <h2 id="best-practices">Best Practices</h2>
              <ul>
                <li>
                  <strong>Keep handlers thin.</strong> A handler should parse input, call a service or the DB, and
                  return a response. If your handler exceeds 50 lines, extract logic into a service.
                </li>
                <li>
                  <strong>Always validate sort columns.</strong> Use a whitelist of allowed column names to prevent
                  SQL injection through the <code>sort_by</code> parameter.
                </li>
                <li>
                  <strong>Use pointers for optional boolean/numeric fields</strong> in update requests.
                  This distinguishes &quot;not provided&quot; from &quot;set to zero/false.&quot;
                </li>
                <li>
                  <strong>Return consistent error responses.</strong> Always use the standard error envelope
                  with <code>code</code> and <code>message</code> fields.
                </li>
                <li>
                  <strong>Preload relationships only when needed.</strong> Use <code>db.Preload(&quot;Author&quot;)</code> in
                  GetByID but not necessarily in List to keep list queries fast.
                </li>
                <li>
                  <strong>Clamp pagination values.</strong> Always enforce <code>page {'>'}= 1</code> and
                  <code>pageSize {'<'}= 100</code> to prevent abuse.
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
