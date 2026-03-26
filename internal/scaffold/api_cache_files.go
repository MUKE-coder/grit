package scaffold

import (
	"fmt"
	"path/filepath"
	"strings"
)

func writeCacheFiles(root string, opts Options) error {
	apiRoot := opts.APIRoot(root)
	module := opts.Module()

	files := map[string]string{
		filepath.Join(apiRoot, "internal", "cache", "cache.go"):      cacheServiceGo(),
		filepath.Join(apiRoot, "internal", "middleware", "cache.go"): cacheMiddlewareGo(),
	}

	for path, content := range files {
		content = strings.ReplaceAll(content, "{{MODULE}}", module)
		if err := writeFile(path, content); err != nil {
			return fmt.Errorf("writing %s: %w", path, err)
		}
	}

	return nil
}

func cacheServiceGo() string {
	return `package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// DefaultTTL is the default cache expiration time.
const DefaultTTL = 5 * time.Minute

// Cache provides a Redis-backed caching service.
type Cache struct {
	client *redis.Client
}

// New creates a new Cache instance connected to the given Redis URL.
func New(redisURL string) (*Cache, error) {
	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, fmt.Errorf("parsing redis URL: %w", err)
	}

	client := redis.NewClient(opts)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("connecting to redis: %w", err)
	}

	return &Cache{client: client}, nil
}

// Get retrieves a cached value and unmarshals it into dest.
// Returns false if the key does not exist.
func (c *Cache) Get(ctx context.Context, key string, dest interface{}) (bool, error) {
	val, err := c.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("cache get %q: %w", key, err)
	}

	if err := json.Unmarshal([]byte(val), dest); err != nil {
		return false, fmt.Errorf("cache unmarshal %q: %w", key, err)
	}

	return true, nil
}

// Set stores a value in the cache with the given TTL.
func (c *Cache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("cache marshal %q: %w", key, err)
	}

	if err := c.client.Set(ctx, key, data, ttl).Err(); err != nil {
		return fmt.Errorf("cache set %q: %w", key, err)
	}

	return nil
}

// Delete removes a key from the cache.
func (c *Cache) Delete(ctx context.Context, key string) error {
	if err := c.client.Del(ctx, key).Err(); err != nil {
		return fmt.Errorf("cache delete %q: %w", key, err)
	}
	return nil
}

// DeletePattern removes all keys matching a glob pattern.
func (c *Cache) DeletePattern(ctx context.Context, pattern string) error {
	iter := c.client.Scan(ctx, 0, pattern, 100).Iterator()
	for iter.Next(ctx) {
		if err := c.client.Del(ctx, iter.Val()).Err(); err != nil {
			return fmt.Errorf("cache delete pattern %q: %w", pattern, err)
		}
	}
	return iter.Err()
}

// Flush clears the entire cache.
func (c *Cache) Flush(ctx context.Context) error {
	return c.client.FlushDB(ctx).Err()
}

// Client returns the underlying Redis client for advanced operations.
func (c *Cache) Client() *redis.Client {
	return c.client
}

// Close closes the Redis connection.
func (c *Cache) Close() error {
	return c.client.Close()
}
`
}

func cacheMiddlewareGo() string {
	return `package middleware

import (
	"crypto/sha256"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"{{MODULE}}/internal/cache"
)

// CacheResponse caches GET request responses in Redis for the given duration.
// Only caches 200 OK responses. Skips caching if no cache service is available.
func CacheResponse(cacheService *cache.Cache, ttl time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		if cacheService == nil || c.Request.Method != http.MethodGet {
			c.Next()
			return
		}

		// Build cache key from URL + query params
		key := fmt.Sprintf("http:%x", sha256.Sum256([]byte(c.Request.URL.String())))

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
		writer := &responseCapture{ResponseWriter: c.Writer, body: make([]byte, 0)}
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
}

type cachedResponse struct {
	Status      int    ` + "`" + `json:"status"` + "`" + `
	ContentType string ` + "`" + `json:"content_type"` + "`" + `
	Body        []byte ` + "`" + `json:"body"` + "`" + `
}

type responseCapture struct {
	gin.ResponseWriter
	body   []byte
	status int
}

func (w *responseCapture) Write(b []byte) (int, error) {
	w.body = append(w.body, b...)
	return w.ResponseWriter.Write(b)
}

func (w *responseCapture) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}
`
}
