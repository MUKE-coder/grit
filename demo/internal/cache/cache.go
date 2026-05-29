package cache

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
