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
	"errors"
	"fmt"
	"log"
	"strings"
	"sync"
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
// QuietDriverLogs collapses the Redis driver's own output.
//
// go-redis writes pool failures straight to stderr, once per attempt per pool,
// so a project started without Redis running greets you with a wall of
// identical "failed to dial after 5 attempts" lines before anything of ours has
// said a word. Silencing it entirely would hide real trouble later, so each
// distinct message is printed once and repeats are dropped.
//
// Called by the server before it connects, and it covers asynq too, which uses
// the same driver underneath.
func QuietDriverLogs() {
	redis.SetLogger(collapsingLogger{seen: map[string]bool{}})
}

type collapsingLogger struct {
	mu   *sync.Mutex
	seen map[string]bool
}

func (l collapsingLogger) Printf(ctx context.Context, format string, v ...interface{}) {
	if l.mu != nil {
		l.mu.Lock()
		defer l.mu.Unlock()
	}
	if l.seen[format] {
		return
	}
	l.seen[format] = true
	// The driver's own messages already start with "redis:", so the prefix is
	// only added when it is missing.
	prefix := "redis: "
	if strings.HasPrefix(format, "redis:") {
		prefix = ""
	}
	log.Printf(prefix+format+" (further identical messages suppressed)", v...)
}

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
	if errors.Is(err, redis.Nil) {
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
	return cacheMiddlewareNew
}
