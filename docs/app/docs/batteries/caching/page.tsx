import Link from 'next/link'
import { ArrowRight, ArrowLeft } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'
import { CodeBlock } from '@/components/code-block'

export default function CachingPage() {
  return (
    <div className="min-h-screen bg-background">
      <SiteHeader />
      <DocsSidebar />

      <main className="lg:pl-64">
        <div className="container max-w-screen-xl py-10 px-6">
          <div className="max-w-3xl">
            {/* Header */}
            <div className="mb-10">
              <span className="tag-mono text-primary/80 mb-3 block">Batteries</span>
              <h1 className="text-4xl font-bold tracking-tight mb-4">
                Redis Caching
              </h1>
              <p className="text-lg text-muted-foreground leading-relaxed">
                Grit provides a Redis-backed caching service with JSON serialization, TTL configuration,
                pattern-based deletion, and a Gin middleware for automatic HTTP response caching.
                Speed up expensive queries and reduce database load with a few lines of code.
              </p>
            </div>

            <div className="prose-grit">
              {/* Cache Service */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Cache Service
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  The cache service at <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">internal/cache/cache.go</code> wraps
                  the <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">go-redis/v9</code> client with
                  convenience methods for storing and retrieving JSON-serialized values.
                </p>

                <CodeBlock filename="internal/cache/cache.go" code={`// DefaultTTL is the default cache expiration time.
const DefaultTTL = 5 * time.Minute

// Cache provides a Redis-backed caching service.
type Cache struct {
    client *redis.Client
}

// New creates a new Cache instance connected to the given Redis URL.
func New(redisURL string) (*Cache, error)

// Get retrieves a cached value and unmarshals it into dest.
// Returns false if the key does not exist.
func (c *Cache) Get(ctx context.Context, key string, dest interface{}) (bool, error)

// Set stores a value in the cache with the given TTL.
func (c *Cache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error

// Delete removes a key from the cache.
func (c *Cache) Delete(ctx context.Context, key string) error

// DeletePattern removes all keys matching a glob pattern.
func (c *Cache) DeletePattern(ctx context.Context, pattern string) error

// Flush clears the entire cache.
func (c *Cache) Flush(ctx context.Context) error

// Client returns the underlying Redis client for advanced operations.
func (c *Cache) Client() *redis.Client

// Close closes the Redis connection.
func (c *Cache) Close() error`} />
              </div>

              {/* JSON Serialization */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  JSON Serialization
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Values are automatically serialized to JSON when stored and deserialized when retrieved.
                  You can cache any Go struct, slice, map, or primitive type.
                </p>

                <CodeBlock filename="cache-examples.go" code={`// Cache a single struct
user := models.User{ID: 1, Name: "John", Email: "john@example.com"}
err := cache.Set(ctx, "user:1", user, 10*time.Minute)

// Retrieve it
var cachedUser models.User
found, err := cache.Get(ctx, "user:1", &cachedUser)
if found {
    fmt.Println(cachedUser.Name) // "John"
}

// Cache a slice
users := []models.User{...}
err = cache.Set(ctx, "users:page:1", users, 5*time.Minute)

// Retrieve the slice
var cachedUsers []models.User
found, err = cache.Get(ctx, "users:page:1", &cachedUsers)

// Cache a map
stats := map[string]int{"total": 100, "active": 42}
err = cache.Set(ctx, "stats:users", stats, 30*time.Second)

// Cache a simple string
err = cache.Set(ctx, "config:motd", "Welcome!", 24*time.Hour)`} />
              </div>

              {/* TTL Configuration */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  TTL Configuration
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Every cache entry requires a TTL (time-to-live). The default is 5 minutes.
                  Choose TTL values based on how frequently the data changes and how stale the
                  data can be.
                </p>

                <div className="rounded-lg border border-border/30 bg-card/30 overflow-hidden">
                  <table className="w-full text-sm">
                    <thead>
                      <tr className="border-b border-border/30 bg-accent/20">
                        <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Data Type</th>
                        <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Suggested TTL</th>
                        <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Reasoning</th>
                      </tr>
                    </thead>
                    <tbody className="text-muted-foreground">
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5">Dashboard stats</td>
                        <td className="px-4 py-2.5 font-mono text-xs">30s - 1m</td>
                        <td className="px-4 py-2.5">Frequently accessed, changes often</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5">List queries</td>
                        <td className="px-4 py-2.5 font-mono text-xs">1m - 5m</td>
                        <td className="px-4 py-2.5">Moderate change frequency</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5">User profiles</td>
                        <td className="px-4 py-2.5 font-mono text-xs">5m - 15m</td>
                        <td className="px-4 py-2.5">Rarely changes, frequently read</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5">Configuration</td>
                        <td className="px-4 py-2.5 font-mono text-xs">1h - 24h</td>
                        <td className="px-4 py-2.5">Almost never changes</td>
                      </tr>
                      <tr>
                        <td className="px-4 py-2.5">External API data</td>
                        <td className="px-4 py-2.5 font-mono text-xs">5m - 1h</td>
                        <td className="px-4 py-2.5">Reduce API calls, respect rate limits</td>
                      </tr>
                    </tbody>
                  </table>
                </div>
              </div>

              {/* CacheResponse Middleware */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  CacheResponse Middleware
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  The <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">CacheResponse</code> middleware
                  at <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">internal/middleware/cache.go</code> automatically
                  caches GET request responses. It hashes the full URL (including query parameters) to
                  generate cache keys, and sets an <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">X-Cache: HIT/MISS</code> header
                  for debugging.
                </p>

                <CodeBlock filename="internal/middleware/cache.go" code={`// CacheResponse caches GET request responses in Redis.
// Only caches 200 OK responses. Skips if no cache service available.
func CacheResponse(cacheService *cache.Cache, ttl time.Duration) gin.HandlerFunc

// How it works:
// 1. Generate cache key from URL: sha256(request.URL.String())
// 2. Check cache: if HIT -> return cached response (X-Cache: HIT)
// 3. If MISS -> capture response, serve it, then cache it
// 4. Only cache 200 OK responses with non-empty bodies`} />

                <h3 className="text-lg font-semibold tracking-tight mb-3">Using the Middleware</h3>

                <CodeBlock filename="internal/routes/routes.go" code={`import "myapp/apps/api/internal/middleware"

// Apply cache middleware to specific routes
api := router.Group("/api")
{
    // Cache the products list for 2 minutes
    api.GET("/products",
        middleware.CacheResponse(cacheService, 2*time.Minute),
        productHandler.List,
    )

    // Cache individual product for 5 minutes
    api.GET("/products/:id",
        middleware.CacheResponse(cacheService, 5*time.Minute),
        productHandler.GetByID,
    )

    // Do NOT cache mutations
    api.POST("/products", productHandler.Create)
    api.PUT("/products/:id", productHandler.Update)
    api.DELETE("/products/:id", productHandler.Delete)
}`} />
              </div>

              {/* Cache Key Patterns */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Cache Key Patterns
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Use consistent key patterns to make cache invalidation predictable.
                  The <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">DeletePattern()</code> method
                  accepts glob patterns, making it easy to clear all keys for a resource.
                </p>

                <CodeBlock filename="cache-key-patterns.go" code={`// Recommended key patterns:
"user:{id}"             // Single resource: user:42
"users:page:{page}"     // Paginated list: users:page:1
"users:count"           // Aggregation
"stats:dashboard"       // Dashboard data
"config:{key}"          // Configuration values

// Set with pattern-aware keys
cache.Set(ctx, "user:42", user, 10*time.Minute)
cache.Set(ctx, "users:page:1", users, 5*time.Minute)
cache.Set(ctx, "users:page:2", users, 5*time.Minute)

// Delete a single key
cache.Delete(ctx, "user:42")

// Delete all keys matching a pattern
cache.DeletePattern(ctx, "users:*")  // Clears all user cache

// Flush the entire cache (use with caution)
cache.Flush(ctx)`} />
              </div>

              {/* When to Use Caching */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  When to Use Caching
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Not everything should be cached. Here are guidelines for when caching adds value
                  versus when it adds unnecessary complexity.
                </p>

                <div className="space-y-3 mb-6">
                  <div className="p-4 rounded-lg border border-green-500/20 bg-green-500/5">
                    <h4 className="text-sm font-semibold text-green-400 mb-2">Good candidates for caching</h4>
                    <ul className="space-y-1.5">
                      {[
                        'Dashboard statistics (computed aggregations)',
                        'Public product/content listings',
                        'Configuration values loaded from database',
                        'External API responses (weather, exchange rates)',
                        'User profile data (read-heavy, write-rare)',
                      ].map((item) => (
                        <li key={item} className="flex items-start gap-2 text-[13px] text-muted-foreground">
                          <span className="text-green-400 mt-0.5">&#10003;</span>
                          {item}
                        </li>
                      ))}
                    </ul>
                  </div>

                  <div className="p-4 rounded-lg border border-red-500/20 bg-red-500/5">
                    <h4 className="text-sm font-semibold text-red-400 mb-2">Not ideal for caching</h4>
                    <ul className="space-y-1.5">
                      {[
                        'User-specific data that changes per request',
                        'Real-time data (chat messages, live feeds)',
                        'Write-heavy endpoints (mutations)',
                        'Data that must always be fresh (payment status)',
                        'Authenticated admin panels with small user bases',
                      ].map((item) => (
                        <li key={item} className="flex items-start gap-2 text-[13px] text-muted-foreground">
                          <span className="text-red-400 mt-0.5">&#10007;</span>
                          {item}
                        </li>
                      ))}
                    </ul>
                  </div>
                </div>
              </div>

              {/* Full Service Example */}
              <div className="mb-10">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Full Service Example
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Here is a complete example of using the cache in a service layer with
                  read-through caching and cache invalidation on writes.
                </p>

                <CodeBlock filename="internal/services/product.go" code={`type ProductService struct {
    DB    *gorm.DB
    Cache *cache.Cache
}

func (s *ProductService) GetByID(ctx context.Context, id uint) (*models.Product, error) {
    key := fmt.Sprintf("product:%d", id)

    // Try cache first
    var product models.Product
    found, err := s.Cache.Get(ctx, key, &product)
    if err == nil && found {
        return &product, nil // Cache HIT
    }

    // Cache MISS -- query database
    if err := s.DB.First(&product, id).Error; err != nil {
        return nil, err
    }

    // Store in cache for next time
    _ = s.Cache.Set(ctx, key, product, 10*time.Minute)

    return &product, nil
}

func (s *ProductService) Update(ctx context.Context, id uint, updates map[string]interface{}) error {
    if err := s.DB.Model(&models.Product{}).Where("id = ?", id).Updates(updates).Error; err != nil {
        return err
    }

    // Invalidate the cached product AND the list cache
    _ = s.Cache.Delete(ctx, fmt.Sprintf("product:%d", id))
    _ = s.Cache.DeletePattern(ctx, "products:*")

    return nil
}`} />
              </div>
            </div>

            {/* Nav */}
            <div className="flex items-center justify-between pt-6 border-t border-border/30">
              <Button variant="ghost" size="sm" asChild className="text-muted-foreground/60 hover:text-foreground">
                <Link href="/docs/batteries/cron" className="gap-1.5">
                  <ArrowLeft className="h-3.5 w-3.5" />
                  Cron Scheduler
                </Link>
              </Button>
              <Button variant="ghost" size="sm" asChild className="text-muted-foreground/60 hover:text-foreground">
                <Link href="/docs/batteries/ai" className="gap-1.5">
                  AI Integration
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
