import Link from 'next/link'
import { ArrowLeft, ArrowRight } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'
import { CodeBlock } from '@/components/code-block'

export default function ServicesPage() {
  return (
    <div className="min-h-screen bg-background">
      <SiteHeader />
      <DocsSidebar />
      <main className="lg:pl-64">
        <div className="container max-w-screen-xl py-10 px-6">
          <div className="max-w-3xl">
            {/* Header */}
            <div className="mb-10">
              <span className="tag-mono text-primary/80 mb-3 block">Backend (Go API)</span>
              <h1 className="text-4xl font-bold tracking-tight mb-4">
                Services
              </h1>
              <p className="text-lg text-muted-foreground leading-relaxed">
                Services contain your application&apos;s business logic. They sit between handlers
                and the database, making your code testable, reusable, and maintainable.
                Handlers should be thin -- services should be fat.
              </p>
            </div>

            <div className="prose-grit">
              {/* ── When to Use Services ─────────────────────────────── */}
              <h2 id="when-to-use">When to Use Services vs. Handlers</h2>
              <p>
                Not every handler needs a service. Here is the rule of thumb:
              </p>

              <div className="rounded-lg border border-border/30 bg-card/30 overflow-hidden mb-6">
                <table className="w-full text-sm">
                  <thead>
                    <tr className="border-b border-border/30 bg-accent/20">
                      <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Scenario</th>
                      <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Where</th>
                    </tr>
                  </thead>
                  <tbody className="text-muted-foreground">
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5">Simple CRUD (fetch, save, delete)</td>
                      <td className="px-4 py-2.5">Handler is fine</td>
                    </tr>
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5">Multiple DB operations in one request</td>
                      <td className="px-4 py-2.5">Use a service with a transaction</td>
                    </tr>
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5">Logic shared between handlers</td>
                      <td className="px-4 py-2.5">Extract into a service</td>
                    </tr>
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5">Complex query building (filters, joins)</td>
                      <td className="px-4 py-2.5">Service or helper function</td>
                    </tr>
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5">External API calls (email, storage, AI)</td>
                      <td className="px-4 py-2.5">Dedicated service</td>
                    </tr>
                    <tr>
                      <td className="px-4 py-2.5">Business rules and validation beyond binding tags</td>
                      <td className="px-4 py-2.5">Service layer</td>
                    </tr>
                  </tbody>
                </table>
              </div>

              <p>
                Grit ships with two built-in services: <code>AuthService</code> (JWT token operations)
                and a pattern you can follow for any new service.
              </p>

              {/* ── Service Pattern ─────────────────────────────── */}
              <h2 id="service-pattern">Service Pattern</h2>
              <p>
                A Grit service is a struct with a <code>DB</code> field (and any other dependencies)
                plus methods that contain business logic. Services live in
                <code>apps/api/internal/services/</code>.
              </p>
              <CodeBlock language="go" filename="apps/api/internal/services/post.go" code={`package services

import (
    "fmt"
    "math"

    "gorm.io/gorm"

    "myapp/apps/api/internal/models"
)

// PostService handles business logic for posts.
type PostService struct {
    DB *gorm.DB
}

// NewPostService creates a new PostService.
func NewPostService(db *gorm.DB) *PostService {
    return &PostService{DB: db}
}`} />
              <p>
                Inject the service into your handler:
              </p>
              <CodeBlock language="go" filename="routes/routes.go" code={`postService := services.NewPostService(db)
postHandler := &handlers.PostHandler{
    DB:      db,
    Service: postService,
}`} />

              {/* ── ListParams ─────────────────────────────── */}
              <h2 id="list-params">ListParams Struct</h2>
              <p>
                For list operations, define a <code>ListParams</code> struct that encapsulates all
                pagination, search, sort, and filter parameters. This keeps service method
                signatures clean and makes it easy to add new filters.
              </p>
              <CodeBlock language="go" filename="services/post.go" code={`// ListParams holds pagination, search, and sort parameters.
type ListParams struct {
    Page      int
    PageSize  int
    Search    string
    SortBy    string
    SortOrder string
    Filters   map[string]string // e.g., {"status": "published"}
}

// ListResult holds paginated query results.
type ListResult struct {
    Data  interface{} \`json:"data"\`
    Total int64       \`json:"total"\`
    Page  int         \`json:"page"\`
    Size  int         \`json:"page_size"\`
    Pages int         \`json:"pages"\`
}

// ClampDefaults ensures pagination values are within safe bounds.
func (p *ListParams) ClampDefaults() {
    if p.Page < 1 {
        p.Page = 1
    }
    if p.PageSize < 1 || p.PageSize > 100 {
        p.PageSize = 20
    }
    if p.SortOrder != "asc" && p.SortOrder != "desc" {
        p.SortOrder = "desc"
    }
    if p.SortBy == "" {
        p.SortBy = "created_at"
    }
}`} />

              {/* ── Query Building ─────────────────────────────── */}
              <h2 id="query-building">Query Building</h2>
              <p>
                Services build GORM queries step by step. This pattern keeps complex queries
                readable and composable.
              </p>
              <CodeBlock language="go" filename="services/post.go -- List" code={`// AllowedSorts defines which columns can be sorted on.
var postAllowedSorts = map[string]bool{
    "id": true, "title": true, "created_at": true, "published": true,
}

// List returns a paginated, filtered list of posts.
func (s *PostService) List(params ListParams) (*ListResult, error) {
    params.ClampDefaults()

    // Validate sort column against whitelist
    if !postAllowedSorts[params.SortBy] {
        params.SortBy = "created_at"
    }

    query := s.DB.Model(&models.Post{})

    // ── Search ──────────────────────────────────────
    if params.Search != "" {
        query = query.Where(
            "title ILIKE ? OR body ILIKE ?",
            "%"+params.Search+"%",
            "%"+params.Search+"%",
        )
    }

    // ── Filters ─────────────────────────────────────
    if status, ok := params.Filters["status"]; ok {
        switch status {
        case "published":
            query = query.Where("published = ?", true)
        case "draft":
            query = query.Where("published = ?", false)
        }
    }

    if authorID, ok := params.Filters["author_id"]; ok {
        query = query.Where("author_id = ?", authorID)
    }

    // ── Count ───────────────────────────────────────
    var total int64
    if err := query.Count(&total).Error; err != nil {
        return nil, fmt.Errorf("counting posts: %w", err)
    }

    // ── Fetch ───────────────────────────────────────
    var posts []models.Post
    offset := (params.Page - 1) * params.PageSize

    err := query.
        Order(params.SortBy + " " + params.SortOrder).
        Offset(offset).
        Limit(params.PageSize).
        Preload("Author").
        Find(&posts).Error

    if err != nil {
        return nil, fmt.Errorf("fetching posts: %w", err)
    }

    pages := int(math.Ceil(float64(total) / float64(params.PageSize)))

    return &ListResult{
        Data:  posts,
        Total: total,
        Page:  params.Page,
        Size:  params.PageSize,
        Pages: pages,
    }, nil
}`} />
              <p>
                The handler becomes much simpler when it delegates to a service:
              </p>
              <CodeBlock language="go" filename="handlers/post.go -- using service" code={`func (h *PostHandler) List(c *gin.Context) {
    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

    params := services.ListParams{
        Page:      page,
        PageSize:  pageSize,
        Search:    c.Query("search"),
        SortBy:    c.DefaultQuery("sort_by", "created_at"),
        SortOrder: c.DefaultQuery("sort_order", "desc"),
        Filters: map[string]string{
            "status":    c.Query("status"),
            "author_id": c.Query("author_id"),
        },
    }

    result, err := h.Service.List(params)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": gin.H{
                "code":    "INTERNAL_ERROR",
                "message": "Failed to fetch posts",
            },
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "data": result.Data,
        "meta": gin.H{
            "total":     result.Total,
            "page":      result.Page,
            "page_size": result.Size,
            "pages":     result.Pages,
        },
    })
}`} />

              {/* ── Business Logic Examples ─────────────────────────────── */}
              <h2 id="business-logic">Business Logic Examples</h2>
              <p>
                Services are the right place for any logic that goes beyond simple CRUD.
                Here are common patterns.
              </p>

              <h3 id="publish-post">Publishing a Post</h3>
              <p>
                A publish action might need to validate the post, update its status, and send a
                notification. All of this belongs in a service:
              </p>
              <CodeBlock language="go" filename="services/post.go -- Publish" code={`// Publish marks a post as published after validation.
func (s *PostService) Publish(postID uint) (*models.Post, error) {
    var post models.Post
    if err := s.DB.First(&post, postID).Error; err != nil {
        return nil, fmt.Errorf("post not found: %w", err)
    }

    if post.Published {
        return nil, fmt.Errorf("post is already published")
    }

    if len(post.Title) < 10 {
        return nil, fmt.Errorf("title must be at least 10 characters to publish")
    }

    if len(post.Body) < 100 {
        return nil, fmt.Errorf("body must be at least 100 characters to publish")
    }

    post.Published = true
    if err := s.DB.Save(&post).Error; err != nil {
        return nil, fmt.Errorf("publishing post: %w", err)
    }

    return &post, nil
}`} />

              <h3 id="user-stats">Aggregation / Statistics</h3>
              <CodeBlock language="go" filename="services/user.go -- Stats" code={`// UserStats holds aggregated user statistics.
type UserStats struct {
    Total       int64 \`json:"total"\`
    Active      int64 \`json:"active"\`
    Admins      int64 \`json:"admins"\`
    NewThisWeek int64 \`json:"new_this_week"\`
}

// GetStats returns aggregated user statistics.
func (s *UserService) GetStats() (*UserStats, error) {
    var stats UserStats

    if err := s.DB.Model(&models.User{}).Count(&stats.Total).Error; err != nil {
        return nil, fmt.Errorf("counting total users: %w", err)
    }

    s.DB.Model(&models.User{}).Where("active = ?", true).Count(&stats.Active)
    s.DB.Model(&models.User{}).Where("role = ?", "admin").Count(&stats.Admins)
    s.DB.Model(&models.User{}).
        Where("created_at >= NOW() - INTERVAL '7 days'").
        Count(&stats.NewThisWeek)

    return &stats, nil
}`} />

              {/* ── Transaction Handling ─────────────────────────────── */}
              <h2 id="transactions">Transaction Handling</h2>
              <p>
                When a service method performs multiple database operations that must succeed or
                fail together, wrap them in a GORM transaction. If any step returns an error,
                the entire transaction is rolled back.
              </p>
              <CodeBlock language="go" filename="services/order.go -- CreateOrder" code={`// CreateOrder creates an order and decrements product stock atomically.
func (s *OrderService) CreateOrder(order *models.Order, items []models.OrderItem) error {
    return s.DB.Transaction(func(tx *gorm.DB) error {
        // Step 1: Create the order
        if err := tx.Create(order).Error; err != nil {
            return fmt.Errorf("creating order: %w", err)
        }

        // Step 2: Create order items and decrement stock
        for i := range items {
            items[i].OrderID = order.ID

            if err := tx.Create(&items[i]).Error; err != nil {
                return fmt.Errorf("creating order item: %w", err)
            }

            // Decrement stock
            result := tx.Model(&models.Product{}).
                Where("id = ? AND stock >= ?", items[i].ProductID, items[i].Quantity).
                Update("stock", gorm.Expr("stock - ?", items[i].Quantity))

            if result.Error != nil {
                return fmt.Errorf("updating stock: %w", result.Error)
            }
            if result.RowsAffected == 0 {
                return fmt.Errorf("insufficient stock for product %d", items[i].ProductID)
            }
        }

        // Step 3: Calculate total
        var total float64
        for _, item := range items {
            total += item.Price * float64(item.Quantity)
        }
        if err := tx.Model(order).Update("total", total).Error; err != nil {
            return fmt.Errorf("updating order total: %w", err)
        }

        return nil // commit
    })
}`} />
              <p>Key points about GORM transactions:</p>
              <ul>
                <li>Use <code>tx</code> (the transaction handle) for all queries inside the callback, not <code>s.DB</code>.</li>
                <li>If the callback returns <code>nil</code>, the transaction commits.</li>
                <li>If the callback returns an error, the transaction rolls back automatically.</li>
                <li>If a panic occurs inside the callback, GORM recovers and rolls back.</li>
              </ul>

              {/* ── Auth Service ─────────────────────────────── */}
              <h2 id="auth-service">Built-in AuthService</h2>
              <p>
                Grit ships with an <code>AuthService</code> that handles all JWT token operations.
                It is the canonical example of a well-structured service.
              </p>
              <CodeBlock language="go" filename="apps/api/internal/services/auth.go" code={`type AuthService struct {
    Secret        string
    AccessExpiry  time.Duration
    RefreshExpiry time.Duration
}

type TokenPair struct {
    AccessToken  string \`json:"access_token"\`
    RefreshToken string \`json:"refresh_token"\`
    ExpiresAt    int64  \`json:"expires_at"\`
}

type Claims struct {
    UserID uint   \`json:"user_id"\`
    Email  string \`json:"email"\`
    Role   string \`json:"role"\`
    jwt.RegisteredClaims
}

// GenerateTokenPair creates access + refresh tokens.
func (s *AuthService) GenerateTokenPair(
    userID uint, email, role string,
) (*TokenPair, error) {
    accessToken, expiresAt, err := s.generateToken(
        userID, email, role, s.AccessExpiry,
    )
    if err != nil {
        return nil, fmt.Errorf("generating access token: %w", err)
    }

    refreshToken, _, err := s.generateToken(
        userID, email, role, s.RefreshExpiry,
    )
    if err != nil {
        return nil, fmt.Errorf("generating refresh token: %w", err)
    }

    return &TokenPair{
        AccessToken:  accessToken,
        RefreshToken: refreshToken,
        ExpiresAt:    expiresAt,
    }, nil
}

// ValidateToken parses and validates a JWT token.
func (s *AuthService) ValidateToken(tokenString string) (*Claims, error) {
    token, err := jwt.ParseWithClaims(
        tokenString, &Claims{},
        func(token *jwt.Token) (interface{}, error) {
            if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                return nil, fmt.Errorf("unexpected signing method")
            }
            return []byte(s.Secret), nil
        },
    )
    if err != nil {
        return nil, fmt.Errorf("parsing token: %w", err)
    }

    claims, ok := token.Claims.(*Claims)
    if !ok || !token.Valid {
        return nil, fmt.Errorf("invalid token")
    }

    return claims, nil
}`} />

              {/* ── Best Practices ─────────────────────────────── */}
              <h2 id="best-practices">Best Practices</h2>
              <ul>
                <li>
                  <strong>One service per resource.</strong> <code>PostService</code>, <code>UserService</code>,
                  <code>OrderService</code> -- each in its own file.
                </li>
                <li>
                  <strong>Return errors, not HTTP codes.</strong> Services should return Go errors.
                  The handler decides the HTTP status code.
                </li>
                <li>
                  <strong>Wrap errors with context.</strong> Use <code>fmt.Errorf(&quot;context: %w&quot;, err)</code> so
                  error messages tell you where things went wrong.
                </li>
                <li>
                  <strong>Use transactions for multi-step operations.</strong> If one step fails, everything
                  rolls back cleanly.
                </li>
                <li>
                  <strong>Keep services independent.</strong> A service should not import another service.
                  If two services need to collaborate, the handler orchestrates them.
                </li>
                <li>
                  <strong>Validate business rules here.</strong> Binding tags handle field-level validation.
                  Services handle business rules like &quot;a post must have at least 100 characters to be published.&quot;
                </li>
              </ul>
            </div>

            {/* Nav */}
            <div className="flex items-center justify-between pt-6 mt-10 border-t border-border/30">
              <Button variant="ghost" size="sm" asChild className="text-muted-foreground/60 hover:text-foreground">
                <Link href="/docs/backend/handlers" className="gap-1.5">
                  <ArrowLeft className="h-3.5 w-3.5" />
                  Handlers
                </Link>
              </Button>
              <Button variant="ghost" size="sm" asChild className="text-muted-foreground/60 hover:text-foreground">
                <Link href="/docs/backend/middleware" className="gap-1.5">
                  Middleware
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
