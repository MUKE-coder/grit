package cache

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"math/rand"
	"sync/atomic"
	"time"

	"github.com/redis/go-redis/v9"
)

// Hit rates, so the question "is the cache working" has a number.
//
// Counters rather than a metrics library: two atomics cost nothing on the hot
// path, and a cache whose hit rate nobody can see is a cache nobody can tune.
// /api/v1/health reports them.
var (
	hits   atomic.Int64
	misses atomic.Int64
)

// Stats is the hit rate since the process started.
func Stats() (hit, miss int64, ratio float64) {
	h, m := hits.Load(), misses.Load()
	if h+m == 0 {
		return 0, 0, 0
	}
	return h, m, float64(h) / float64(h+m)
}

// Remember returns the cached value for key, or computes it, stores it and
// returns that.
//
// The cache-aside pattern, with the three things a hand-written version usually
// misses:
//
//  1. A cache failure falls through to the loader. Redis being down makes the
//     app slower, not broken, and that difference is the whole reason a cache
//     is allowed to be a dependency at all.
//
//  2. Jittered TTLs. Ten thousand keys written in the same second by a deploy
//     expire in the same second without it, and the database gets all of them
//     back at once. Up to 10% is added to each.
//
//  3. Stampede protection. When a hot key expires, one caller rebuilds it and
//     the rest wait briefly and read the new value, instead of every concurrent
//     request running the same expensive query.
//
// Usage:
//
//	var count int64
//	err := cache.Remember(ctx, c, "user:"+id+":followers", time.Minute, &count, func() (int64, error) {
//	    var n int64
//	    err := db.Model(&models.Follow{}).Where("followee_id = ?", id).Count(&n).Error
//	    return n, err
//	})
//
// dest must be a pointer. The loader's return value is what gets cached.
func Remember[T any](
	ctx context.Context,
	c *Cache,
	key string,
	ttl time.Duration,
	dest *T,
	loader func() (T, error),
) error {
	if c == nil {
		v, err := loader()
		if err != nil {
			return err
		}
		*dest = v
		return nil
	}

	if found, err := c.Get(ctx, key, dest); err == nil && found {
		hits.Add(1)
		return nil
	} else if err != nil {
		// Cache down is not app down. Say so once and carry on to the loader.
		log.Printf("cache: reading %s failed, falling through to the source: %v", key, err)
	}
	misses.Add(1)

	// One rebuilder per key. The lock is short and self-expiring, so a process
	// that dies mid-rebuild delays the next attempt by seconds rather than
	// leaving the key uncacheable.
	lockKey := "lock:" + key
	locked, lockErr := c.Client().SetNX(ctx, lockKey, "1", 10*time.Second).Result()
	if lockErr != nil {
		// No lock available means no stampede protection, not no answer.
		locked = true
	}

	if !locked {
		// Somebody else is rebuilding. Wait for them rather than joining in.
		if v, ok := waitForRebuild[T](ctx, c, key); ok {
			hits.Add(1)
			*dest = v
			return nil
		}
		// They did not finish in time. Load it rather than waiting longer: a
		// slow answer beats a missing one.
	} else {
		defer func() {
			if err := c.Client().Del(context.WithoutCancel(ctx), lockKey).Err(); err != nil &&
				!errors.Is(err, redis.Nil) {
				log.Printf("cache: releasing the lock on %s failed: %v", key, err)
			}
		}()
	}

	value, err := loader()
	if err != nil {
		return err
	}
	*dest = value

	// Written with the request's deadline removed: the value was computed at
	// the caller's expense and throwing it away because the client hung up
	// means the next request pays again.
	if err := c.Set(context.WithoutCancel(ctx), key, value, jitter(ttl)); err != nil {
		log.Printf("cache: storing %s failed: %v", key, err)
	}
	return nil
}

// waitForRebuild polls briefly for the value another caller is building.
//
// Polling rather than pub/sub on purpose: the wait is under a second, and a
// subscription per cache miss is a connection per cache miss.
func waitForRebuild[T any](ctx context.Context, c *Cache, key string) (T, bool) {
	var zero T
	for i := 0; i < 10; i++ {
		select {
		case <-ctx.Done():
			return zero, false
		case <-time.After(50 * time.Millisecond):
		}
		var v T
		if found, err := c.Get(ctx, key, &v); err == nil && found {
			return v, true
		}
	}
	return zero, false
}

// Forget drops one key. Call it from the write that makes the value wrong.
//
// TTL and delete together, not one or the other: the delete keeps the value
// fresh, and the TTL is the safety net for the invalidation somebody forgets.
func Forget(ctx context.Context, c *Cache, keys ...string) {
	if c == nil || len(keys) == 0 {
		return
	}
	if err := c.Client().Del(ctx, keys...).Err(); err != nil && !errors.Is(err, redis.Nil) {
		log.Printf("cache: forgetting %v failed: %v", keys, err)
	}
}

// jitter spreads expiry so a class of keys written together does not expire
// together. Up to 10% later, never earlier.
func jitter(ttl time.Duration) time.Duration {
	if ttl <= 0 {
		return ttl
	}
	// math/rand, not crypto/rand: this decides how long a cache entry lives,
	// and an attacker who could predict it learns when a key expires, which is
	// a thing the TTL already tells them.
	//nolint:gosec // G404: TTL spread, not a secret
	//#nosec G404
	return ttl + time.Duration(rand.Int63n(int64(ttl)/10+1))
}

// MarshalCheck fails loudly at startup if a type cannot round-trip through the
// cache. Values are stored as JSON, so an unexported field or a channel comes
// back empty rather than erroring, which is the kind of bug that looks like a
// cache miss that never hits.
func MarshalCheck(v any) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}
	if len(b) == 0 || string(b) == "{}" {
		return errors.New("value marshals to nothing: unexported fields are not cached")
	}
	return nil
}
