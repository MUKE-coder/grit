import { CodeBlock } from '@/components/code-block'
import { TipBox } from '@/components/course/tip-box'
import { Exercise } from '@/components/course/exercise'
import { KnowledgeCheck } from '@/components/course/knowledge-check'

export default function Lesson() {
  return (
    <>
      <p>
        Redis is an in-memory key-value store. In Grit, it&apos;s the
        cache that sits between your handler and the database — saving
        you a query on hot paths. This lesson is what it does, where
        the code lives, and how to add caching to your own endpoint.
      </p>

      <h2>Why we need it</h2>
      <p>
        A list endpoint that runs the same query 1000 times per minute
        wastes the DB. With a 60-second cache, that&apos;s 1 DB query
        per minute instead of 1000. The DB is happier, p95 latency
        drops, and you don&apos;t need a bigger DB to scale.
      </p>

      <h2>Where it lives</h2>
      <pre className="not-prose text-xs leading-relaxed bg-bg-elevated border border-border rounded-lg p-4 overflow-x-auto">{`apps/api/internal/cache/
├── cache.go         ← Service struct + Get / Set / Delete
└── middleware.go    ← Optional HTTP-level cache (response caching)`}</pre>

      <h2>How it&apos;s implemented</h2>
      <CodeBlock
        language="go"
        filename="apps/api/internal/cache/cache.go (simplified)"
        code={`package cache

import (
  "context"
  "encoding/json"
  "time"
  "github.com/redis/go-redis/v9"
)

type Service struct {
  rdb *redis.Client
}

func New(url string) (*Service, error) {
  opt, err := redis.ParseURL(url)
  if err != nil { return nil, err }
  return &Service{rdb: redis.NewClient(opt)}, nil
}

func (s *Service) Get(ctx context.Context, key string, dest any) (bool, error) {
  raw, err := s.rdb.Get(ctx, key).Bytes()
  if err == redis.Nil { return false, nil }   // cache miss — not an error
  if err != nil { return false, err }
  return true, json.Unmarshal(raw, dest)
}

func (s *Service) Set(ctx context.Context, key string, val any, ttl time.Duration) error {
  b, err := json.Marshal(val)
  if err != nil { return err }
  return s.rdb.Set(ctx, key, b, ttl).Err()
}

func (s *Service) Del(ctx context.Context, keys ...string) error {
  return s.rdb.Del(ctx, keys...).Err()
}`}
      />
      <p>
        Three operations, JSON serialization baked in, cache-miss
        treated as &quot;not an error, just no data&quot;. That last
        detail matters: a cache miss is a normal case, not an
        exception.
      </p>

      <h2>How you call it from a service</h2>
      <CodeBlock
        language="go"
        filename="apps/api/internal/services/product_service.go"
        code={`func (s *ProductService) ListFeatured(ctx context.Context) ([]models.Product, error) {
  const key = "products:featured:v1"

  // 1. try cache
  var cached []models.Product
  if hit, _ := s.cache.Get(ctx, key, &cached); hit {
    return cached, nil
  }

  // 2. miss — go to DB
  var products []models.Product
  if err := s.db.WithContext(ctx).Where("is_featured = ?", true).Find(&products).Error; err != nil {
    return nil, err
  }

  // 3. populate cache for next time (1-minute TTL)
  _ = s.cache.Set(ctx, key, products, time.Minute)

  return products, nil
}`}
      />
      <p>
        Three lines you&apos;ll see in every cached-read service:
      </p>
      <ul>
        <li>
          <strong>Try cache.</strong> If hit, return immediately. The
          DB never gets touched.
        </li>
        <li>
          <strong>Hit DB.</strong> Normal query.
        </li>
        <li>
          <strong>Populate cache.</strong> So the NEXT request is a hit.
        </li>
      </ul>

      <h2>Key design — the silent skill</h2>
      <CodeBlock
        language="go"
        code={`// Good keys:
"products:featured:v1"                  // collection
"product:" + strconv.Itoa(id) + ":v1"  // single item
"user:" + uid + ":notes:v1"           // user-scoped

// Bad keys:
"prod123"          // unclear, collision-prone
"q=SELECT...""    // raw query — too brittle
"abc-" + time.Now() // includes time → never hits!`}
      />
      <ul>
        <li>
          <strong>Versioned (<code>:v1</code>).</strong> When you
          change the value shape, bump the version. Old keys expire
          naturally; no manual invalidation.
        </li>
        <li>
          <strong>Hierarchical with <code>:</code></strong>. Lets you
          group + scan keys for debugging in <code>redis-cli</code>.
        </li>
        <li>
          <strong>Deterministic.</strong> Same input → same key.
          Otherwise you never hit.
        </li>
      </ul>

      <h2>Invalidation — the hard part</h2>
      <p>
        &quot;There are only two hard problems in computer science:
        cache invalidation and naming things.&quot; Three patterns to
        know:
      </p>

      <h3>1. TTL — let it expire</h3>
      <p>
        Set a short TTL (30s, 60s, 5m) and accept stale data for that
        long. Simple, robust, hides bugs. Right answer for most cases.
      </p>

      <h3>2. Explicit Delete on write</h3>
      <CodeBlock
        language="go"
        code={`func (s *ProductService) Update(ctx context.Context, p Product) error {
  if err := s.db.WithContext(ctx).Save(&p).Error; err != nil { return err }
  _ = s.cache.Del(ctx, "products:featured:v1", "product:" + p.ID + ":v1")
  return nil
}`}
      />
      <p>
        After a write, blow away the cached read. Next read repopulates.
        Works well when you know exactly which keys are stale.
      </p>

      <h3>3. Pattern-based deletion</h3>
      <p>
        For broader invalidation (delete every product cache for a
        user), use <code>SCAN</code> + <code>DEL</code> in a loop. Avoid
        in production at scale — keep keys scoped instead.
      </p>

      <h2>HTTP-level cache middleware</h2>
      <CodeBlock
        language="go"
        filename="apps/api/internal/cache/middleware.go (sketch)"
        code={`func ResponseCache(c *Service, ttl time.Duration) gin.HandlerFunc {
  return func(ctx *gin.Context) {
    if ctx.Request.Method != "GET" { ctx.Next(); return }

    key := "http:" + ctx.Request.URL.String()
    var body []byte
    if hit, _ := c.Get(ctx, key, &body); hit {
      ctx.Header("X-Cache", "HIT")
      ctx.Data(200, "application/json", body)
      ctx.Abort()
      return
    }

    // capture the response so we can cache it
    writer := newCaptureWriter(ctx.Writer)
    ctx.Writer = writer
    ctx.Next()

    if writer.Status() == 200 {
      _ = c.Set(ctx, key, writer.Body(), ttl)
      ctx.Header("X-Cache", "MISS")
    }
  }
}`}
      />
      <p>
        Wire this on any GET endpoint where the response is identical
        for the same URL across users. The{' '}
        <code>X-Cache: HIT/MISS</code> header lets you verify in the
        browser DevTools.
      </p>

      <TipBox tone="warning">
        <strong>Don&apos;t cache user-specific responses with HTTP
        middleware.</strong> If <code>/api/me</code> returns user A&apos;s
        profile and you cache by URL, user B gets user A&apos;s data.
        Use SERVICE-level caching with the user ID in the key for
        anything personalised.
      </TipBox>

      <h2>How to modify the cache battery</h2>
      <ul>
        <li>
          <strong>Change the default TTL</strong> — find the service or
          middleware that sets <code>time.Minute</code> and adjust.
        </li>
        <li>
          <strong>Swap Redis for an alternative</strong> (Dragonfly,
          KeyDB) — they speak the Redis protocol, so just change{' '}
          <code>REDIS_URL</code>.
        </li>
        <li>
          <strong>Add metrics</strong> — wrap{' '}
          <code>Get</code> / <code>Set</code> with a Prometheus counter
          for hit ratio. One of the highest-signal metrics you can add.
        </li>
        <li>
          <strong>Disable in dev</strong> — set{' '}
          <code>REDIS_URL=&quot;&quot;</code> and Grit&apos;s service
          falls through to no-op (always returns miss). Useful when
          debugging stale-data bugs.
        </li>
      </ul>

      <h2>Local dev wiring</h2>
      <CodeBlock
        language="yaml"
        filename="docker-compose.yml (excerpt)"
        code={`redis:
  image: redis:7-alpine
  ports: ["6379:6379"]`}
      />
      <CodeBlock
        language="env"
        filename=".env (local)"
        code={`REDIS_URL=redis://localhost:6380`}
      />
      <p>
        In production, point <code>REDIS_URL</code> at Upstash, Redis
        Cloud, ElastiCache, or self-hosted. Same code, no change.
      </p>

      <KnowledgeCheck
        question="You add a 60-second cache to /api/products/featured. A teammate complains: &quot;I edited a product but it took a minute to show as featured.&quot; What's the cleanest fix?"
        choices={[
          {
            label: 'Drop the TTL to 1 second',
            feedback:
              "Almost defeats the purpose of caching. You'd still see stale data for a second; the DB still gets hit constantly.",
          },
          {
            label: 'Add an explicit cache.Del in the Update service so the next read repopulates with the fresh data',
            correct: true,
            feedback:
              "Right — TTL is the safety net; explicit invalidation on write is the precise fix. Both can co-exist: invalidate on write, fall back to TTL if you miss a write path.",
          },
          {
            label: 'Remove caching entirely',
            feedback:
              "You'd kill performance to fix a UX edge case. Invalidation is the better tradeoff.",
          },
          {
            label: 'Cache the write too',
            feedback:
              "Caching writes is mostly nonsense — writes need to land in the DB so subsequent reads get truth.",
          },
        ]}
      />

      <Exercise
        prompt={
          <>
            <p>Add a cache to a real endpoint:</p>
            <ol>
              <li>Pick a GET endpoint with no per-user data (e.g., the products list).</li>
              <li>In the service, wrap the DB call with Get / Set / 30s TTL.</li>
              <li>
                In the Update + Delete services for the same model,
                add a <code>cache.Del</code>.
              </li>
              <li>
                Test: hit the list endpoint twice with{' '}
                <code>-v</code> via curl. Second request should be
                faster. (You can also add a temporary log line in the
                service to confirm cache hit / miss.)
              </li>
              <li>
                Test invalidation: edit a row, immediately re-fetch
                the list, confirm the change shows up.
              </li>
            </ol>
          </>
        }
        hint={
          <>
            For the speedup measurement, add{' '}
            <code>time curl ...</code> or just look at the{' '}
            <code>X-Response-Time</code> header Grit logs by default.
            On a tiny DB the delta is small (microseconds); on a real
            production list, 50ms vs 2ms.
          </>
        }
        solution={
          <>
            <p>
              You should see ~10x speedup on a cache hit and an
              instantly-updating list after a write. That&apos;s the
              cache pattern at its best.
            </p>
          </>
        }
      />

      <h2>What&apos;s next</h2>
      <p>
        Next lesson — <strong>S3 file storage</strong>. Uploads, signed
        URLs, image processing — without taking on AWS&apos; whole
        learning curve.
      </p>
    </>
  )
}
