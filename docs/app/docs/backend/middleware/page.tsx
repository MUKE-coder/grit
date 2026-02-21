import Link from 'next/link'
import { ArrowLeft, ArrowRight } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'
import { CodeBlock } from '@/components/code-block'
import { getDocMetadata } from '@/config/docs-metadata'

export const metadata = getDocMetadata('/docs/backend/middleware')

export default function MiddlewarePage() {
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
                Middleware
              </h1>
              <p className="text-lg text-muted-foreground leading-relaxed">
                Middleware functions intercept HTTP requests before they reach your handlers.
                Grit ships with four built-in middleware: Auth (JWT validation), CORS,
                Logger, and Cache. You can also create custom middleware.
              </p>
            </div>

            <div className="prose-grit">
              {/* ── Middleware Order ─────────────────────────────── */}
              <h2 id="middleware-order">Middleware Order & Registration</h2>
              <p>
                Middleware is registered in <code>routes/routes.go</code>. The order matters --
                middleware executes in the order it is registered.
              </p>
              <CodeBlock filename="apps/api/internal/routes/routes.go" code={`func Setup(db *gorm.DB, cfg *config.Config, svc *Services) *gin.Engine {
    r := gin.New()

    // ── Global middleware (applied to ALL routes) ────────
    r.Use(middleware.Logger())       // 1. Log every request
    r.Use(gin.Recovery())            // 2. Recover from panics
    r.Use(middleware.CORS(cfg.CORSOrigins)) // 3. Handle CORS

    // ── Public routes (no auth required) ────────────────
    auth := r.Group("/api/auth")
    {
        auth.POST("/register", authHandler.Register)
        auth.POST("/login", authHandler.Login)
    }

    // ── Protected routes (auth middleware applied) ──────
    protected := r.Group("/api")
    protected.Use(middleware.Auth(db, authService)) // 4. Require JWT
    {
        protected.GET("/posts", postHandler.List)
    }

    // ── Admin routes (auth + role check) ────────────────
    admin := r.Group("/api")
    admin.Use(middleware.Auth(db, authService))     // 4. Require JWT
    admin.Use(middleware.RequireRole("admin"))       // 5. Require admin role
    {
        admin.GET("/users", userHandler.List)
        admin.DELETE("/users/:id", userHandler.Delete)
    }

    return r
}`} />

              <div className="p-4 rounded-lg border border-primary/20 bg-primary/5 mb-6">
                <p className="text-sm text-foreground/80 mb-0">
                  <strong>Execution order:</strong> Logger runs first, then Recovery, then CORS,
                  then Auth (if on a protected route), then RequireRole (if on an admin route),
                  then finally the handler itself.
                </p>
              </div>

              {/* ── Auth Middleware ─────────────────────────────── */}
              <h2 id="auth-middleware">Auth Middleware</h2>
              <p>
                The Auth middleware extracts the JWT token from the <code>Authorization</code> header,
                validates it, loads the user from the database, and stores the user in the
                Gin context for downstream handlers.
              </p>
              <CodeBlock filename="apps/api/internal/middleware/auth.go" code={`func Auth(db *gorm.DB, authService *services.AuthService) gin.HandlerFunc {
    return func(c *gin.Context) {
        // 1. Extract the Authorization header
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.JSON(http.StatusUnauthorized, gin.H{
                "error": gin.H{
                    "code":    "UNAUTHORIZED",
                    "message": "Authorization header is required",
                },
            })
            c.Abort()
            return
        }

        // 2. Parse "Bearer <token>"
        parts := strings.SplitN(authHeader, " ", 2)
        if len(parts) != 2 || parts[0] != "Bearer" {
            c.JSON(http.StatusUnauthorized, gin.H{
                "error": gin.H{
                    "code":    "UNAUTHORIZED",
                    "message": "Invalid authorization header format",
                },
            })
            c.Abort()
            return
        }

        // 3. Validate the JWT token
        claims, err := authService.ValidateToken(parts[1])
        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{
                "error": gin.H{
                    "code":    "UNAUTHORIZED",
                    "message": "Invalid or expired token",
                },
            })
            c.Abort()
            return
        }

        // 4. Load user from database
        var user models.User
        if err := db.First(&user, claims.UserID).Error; err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{
                "error": gin.H{
                    "code":    "UNAUTHORIZED",
                    "message": "User not found",
                },
            })
            c.Abort()
            return
        }

        // 5. Check if the account is active
        if !user.Active {
            c.JSON(http.StatusForbidden, gin.H{
                "error": gin.H{
                    "code":    "ACCOUNT_DISABLED",
                    "message": "Your account has been disabled",
                },
            })
            c.Abort()
            return
        }

        // 6. Store user data in context for handlers
        c.Set("user", user)
        c.Set("user_id", user.ID)
        c.Set("user_role", user.Role)
        c.Next()
    }
}`} />
              <p>
                After this middleware runs, handlers can access the authenticated user:
              </p>
              <CodeBlock filename="accessing user in handler" code={`func (h *PostHandler) Create(c *gin.Context) {
    // Get the full user object
    user, _ := c.Get("user")
    currentUser := user.(models.User)

    // Or get individual fields
    userID, _ := c.Get("user_id")    // uint
    role, _ := c.Get("user_role")     // string
}`} />

              {/* ── RequireRole ─────────────────────────────── */}
              <h2 id="require-role">RequireRole Middleware</h2>
              <p>
                <code>RequireRole</code> checks if the authenticated user has one of the specified roles.
                It must be used <strong>after</strong> the Auth middleware (which sets <code>user_role</code>
                in the context).
              </p>
              <CodeBlock filename="apps/api/internal/middleware/auth.go" code={`// RequireRole checks if the user has one of the required roles.
func RequireRole(roles ...string) gin.HandlerFunc {
    return func(c *gin.Context) {
        userRole, exists := c.Get("user_role")
        if !exists {
            c.JSON(http.StatusUnauthorized, gin.H{
                "error": gin.H{
                    "code":    "UNAUTHORIZED",
                    "message": "Not authenticated",
                },
            })
            c.Abort()
            return
        }

        role, ok := userRole.(string)
        if !ok {
            c.JSON(http.StatusInternalServerError, gin.H{
                "error": gin.H{
                    "code":    "INTERNAL_ERROR",
                    "message": "Invalid user role",
                },
            })
            c.Abort()
            return
        }

        for _, r := range roles {
            if role == r {
                c.Next()
                return
            }
        }

        c.JSON(http.StatusForbidden, gin.H{
            "error": gin.H{
                "code":    "FORBIDDEN",
                "message": "You do not have permission to access this resource",
            },
        })
        c.Abort()
    }
}`} />
              <p>Usage examples:</p>
              <CodeBlock filename="role_examples.go" code={`// Admin only
admin.Use(middleware.RequireRole("admin"))

// Admin or editor
editors.Use(middleware.RequireRole("admin", "editor"))

// Any authenticated user (no RequireRole needed, just Auth middleware)`} />

              {/* ── CORS Middleware ─────────────────────────────── */}
              <h2 id="cors-middleware">CORS Middleware</h2>
              <p>
                The CORS middleware allows your Next.js frontend (running on a different port)
                to make API requests to the Go backend. It reads the allowed origins from
                the <code>CORS_ORIGINS</code> environment variable.
              </p>
              <CodeBlock filename="apps/api/internal/middleware/cors.go" code={`func CORS(allowedOrigins []string) gin.HandlerFunc {
    originsMap := make(map[string]bool)
    for _, origin := range allowedOrigins {
        originsMap[origin] = true
    }

    return func(c *gin.Context) {
        origin := c.GetHeader("Origin")

        if originsMap[origin] {
            c.Header("Access-Control-Allow-Origin", origin)
        }

        c.Header("Access-Control-Allow-Methods",
            "GET, POST, PUT, PATCH, DELETE, OPTIONS")
        c.Header("Access-Control-Allow-Headers",
            "Origin, Content-Type, Accept, Authorization")
        c.Header("Access-Control-Allow-Credentials", "true")
        c.Header("Access-Control-Max-Age", "86400")

        if c.Request.Method == http.MethodOptions {
            c.AbortWithStatus(http.StatusNoContent)
            return
        }

        c.Next()
    }
}`} />
              <p>
                Configure allowed origins in your <code>.env</code> file:
              </p>
              <CodeBlock language="bash" filename=".env" code={`CORS_ORIGINS=http://localhost:3000,http://localhost:3001`} />

              {/* ── Logger Middleware ─────────────────────────────── */}
              <h2 id="logger-middleware">Logger Middleware</h2>
              <p>
                The Logger middleware logs every HTTP request with its status code, method, path,
                client IP, and response latency. It runs before all other middleware so it
                captures the total request time.
              </p>
              <CodeBlock filename="apps/api/internal/middleware/logger.go" code={`func Logger() gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()
        path := c.Request.URL.Path
        query := c.Request.URL.RawQuery

        c.Next()

        latency := time.Since(start)
        status := c.Writer.Status()
        method := c.Request.Method
        clientIP := c.ClientIP()

        if query != "" {
            path = path + "?" + query
        }

        log.Printf("[%d] %s %s | %s | %v",
            status, method, path, clientIP, latency,
        )
    }
}`} />
              <p>Example output:</p>
              <CodeBlock language="bash" code={`[200] GET /api/users?page=1 | 127.0.0.1 | 3.241ms
[201] POST /api/posts | 127.0.0.1 | 12.507ms
[401] GET /api/auth/me | 127.0.0.1 | 0.128ms`} />

              {/* ── Cache Middleware ─────────────────────────────── */}
              <h2 id="cache-middleware">Cache Middleware</h2>
              <p>
                The <code>CacheResponse</code> middleware caches the full HTTP response for GET requests in Redis.
                Subsequent identical requests are served from cache with an <code>X-Cache: HIT</code> header.
              </p>
              <CodeBlock filename="apps/api/internal/middleware/cache.go" code={`func CacheResponse(
    cacheService *cache.Cache, ttl time.Duration,
) gin.HandlerFunc {
    return func(c *gin.Context) {
        // Only cache GET requests
        if cacheService == nil || c.Request.Method != http.MethodGet {
            c.Next()
            return
        }

        // Build cache key from the full URL (path + query)
        key := fmt.Sprintf("http:%x",
            sha256.Sum256([]byte(c.Request.URL.String())),
        )

        // Try to serve from cache
        var cached cachedResponse
        found, err := cacheService.Get(c.Request.Context(), key, &cached)
        if err == nil && found {
            c.Header("X-Cache", "HIT")
            c.Data(cached.Status, cached.ContentType, cached.Body)
            c.Abort()
            return
        }

        // Capture the response
        writer := &responseCapture{
            ResponseWriter: c.Writer,
            body: make([]byte, 0),
        }
        c.Writer = writer
        c.Header("X-Cache", "MISS")

        c.Next()

        // Cache successful responses
        if writer.status == http.StatusOK && len(writer.body) > 0 {
            resp := cachedResponse{
                Status:      writer.status,
                ContentType: writer.Header().Get("Content-Type"),
                Body:        writer.body,
            }
            _ = cacheService.Set(c.Request.Context(), key, resp, ttl)
        }
    }
}`} />
              <p>Apply it to specific routes:</p>
              <CodeBlock filename="cache_usage.go" code={`// Cache the posts list for 5 minutes
protected.GET("/posts",
    middleware.CacheResponse(svc.Cache, 5*time.Minute),
    postHandler.List,
)

// Cache individual post for 10 minutes
protected.GET("/posts/:id",
    middleware.CacheResponse(svc.Cache, 10*time.Minute),
    postHandler.GetByID,
)`} />

              {/* ── Custom Middleware ─────────────────────────────── */}
              <h2 id="custom-middleware">Creating Custom Middleware</h2>
              <p>
                A Gin middleware is any function that returns <code>gin.HandlerFunc</code>.
                Here is the pattern for creating your own:
              </p>
              <CodeBlock filename="apps/api/internal/middleware/rate_limit.go" code={`package middleware

import (
    "net/http"
    "sync"
    "time"

    "github.com/gin-gonic/gin"
)

// RateLimit creates a simple in-memory rate limiter.
// maxRequests is the maximum number of requests allowed
// within the given window duration per IP address.
func RateLimit(maxRequests int, window time.Duration) gin.HandlerFunc {
    type client struct {
        count    int
        lastSeen time.Time
    }

    var mu sync.Mutex
    clients := make(map[string]*client)

    // Clean up old entries periodically
    go func() {
        for {
            time.Sleep(window)
            mu.Lock()
            for ip, c := range clients {
                if time.Since(c.lastSeen) > window {
                    delete(clients, ip)
                }
            }
            mu.Unlock()
        }
    }()

    return func(c *gin.Context) {
        ip := c.ClientIP()

        mu.Lock()
        cl, exists := clients[ip]
        if !exists {
            clients[ip] = &client{count: 1, lastSeen: time.Now()}
            mu.Unlock()
            c.Next()
            return
        }

        if time.Since(cl.lastSeen) > window {
            cl.count = 1
            cl.lastSeen = time.Now()
            mu.Unlock()
            c.Next()
            return
        }

        cl.count++
        cl.lastSeen = time.Now()

        if cl.count > maxRequests {
            mu.Unlock()
            c.JSON(http.StatusTooManyRequests, gin.H{
                "error": gin.H{
                    "code":    "RATE_LIMITED",
                    "message": "Too many requests, please try again later",
                },
            })
            c.Abort()
            return
        }

        mu.Unlock()
        c.Next()
    }
}`} />
              <p>Register it on routes that need rate limiting:</p>
              <CodeBlock filename="routes.go" code={`// Limit auth endpoints to 10 requests per minute per IP
auth := r.Group("/api/auth")
auth.Use(middleware.RateLimit(10, 1*time.Minute))
{
    auth.POST("/login", authHandler.Login)
    auth.POST("/register", authHandler.Register)
}`} />

              {/* ── Middleware Anatomy ─────────────────────────────── */}
              <h2 id="middleware-anatomy">Middleware Anatomy</h2>
              <p>Every Gin middleware follows the same structure:</p>
              <CodeBlock filename="middleware_anatomy.go" code={`func MyMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // ── BEFORE the handler ──────────────────────
        // Check conditions, set context values, etc.

        // Option A: Continue to next middleware/handler
        c.Next()

        // ── AFTER the handler ───────────────────────
        // Log results, clean up, etc.

        // Option B: Stop the chain (use instead of c.Next())
        // c.Abort()
        // c.JSON(http.StatusForbidden, gin.H{...})
    }
}`} />
              <ul>
                <li>
                  <strong><code>c.Next()</code></strong> passes control to the next middleware or handler
                  in the chain. Code after <code>c.Next()</code> runs after the handler returns.
                </li>
                <li>
                  <strong><code>c.Abort()</code></strong> stops the chain. No further middleware or handlers
                  will execute.
                </li>
                <li>
                  <strong><code>c.Set(key, value)</code></strong> stores data in the context for downstream
                  middleware and handlers.
                </li>
                <li>
                  <strong><code>c.Get(key)</code></strong> retrieves data stored by upstream middleware.
                </li>
              </ul>

              {/* ── Summary Table ─────────────────────────────── */}
              <h2 id="summary">Built-in Middleware Summary</h2>
              <div className="rounded-lg border border-border/30 bg-card/30 overflow-hidden mb-6">
                <table className="w-full text-sm">
                  <thead>
                    <tr className="border-b border-border/30 bg-accent/20">
                      <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Middleware</th>
                      <th className="text-left px-4 py-2.5 font-medium text-foreground/80">File</th>
                      <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Purpose</th>
                    </tr>
                  </thead>
                  <tbody className="text-muted-foreground">
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5 font-mono text-xs">Logger()</td>
                      <td className="px-4 py-2.5 font-mono text-xs">middleware/logger.go</td>
                      <td className="px-4 py-2.5">Logs every request with status, method, path, latency</td>
                    </tr>
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5 font-mono text-xs">CORS(origins)</td>
                      <td className="px-4 py-2.5 font-mono text-xs">middleware/cors.go</td>
                      <td className="px-4 py-2.5">Sets CORS headers, handles preflight OPTIONS</td>
                    </tr>
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5 font-mono text-xs">Auth(db, svc)</td>
                      <td className="px-4 py-2.5 font-mono text-xs">middleware/auth.go</td>
                      <td className="px-4 py-2.5">Validates JWT, loads user, sets context</td>
                    </tr>
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5 font-mono text-xs">RequireRole(roles...)</td>
                      <td className="px-4 py-2.5 font-mono text-xs">middleware/auth.go</td>
                      <td className="px-4 py-2.5">Checks user role against allowed list</td>
                    </tr>
                    <tr>
                      <td className="px-4 py-2.5 font-mono text-xs">CacheResponse(cache, ttl)</td>
                      <td className="px-4 py-2.5 font-mono text-xs">middleware/cache.go</td>
                      <td className="px-4 py-2.5">Caches GET responses in Redis</td>
                    </tr>
                  </tbody>
                </table>
              </div>
            </div>

            {/* Nav */}
            <div className="flex items-center justify-between pt-6 mt-10 border-t border-border/30">
              <Button variant="ghost" size="sm" asChild className="text-muted-foreground/60 hover:text-foreground">
                <Link href="/docs/backend/services" className="gap-1.5">
                  <ArrowLeft className="h-3.5 w-3.5" />
                  Services
                </Link>
              </Button>
              <Button variant="ghost" size="sm" asChild className="text-muted-foreground/60 hover:text-foreground">
                <Link href="/docs/backend/authentication" className="gap-1.5">
                  Authentication
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
