import Link from 'next/link'
import type { Metadata } from 'next'
import { GuidePage, Section, H3, Outcome, Note, Warn, Checklist, Faq, Related } from '@/components/guide'
import { CodeBlock } from '@/components/code-block'

export const metadata: Metadata = {
  title: 'Caching in Go: cache-aside, jittered TTLs and stampede protection',
  description:
    'The thundering herd, singleflight, why a fixed TTL synchronises your expiries, and invalidation you can actually reason about.',
  alternates: { canonical: 'https://gritframework.dev/guides/go-redis-caching-stampede' },
  openGraph: {
    title: 'Caching in Go: cache-aside, jittered TTLs and stampede protection',
    description: 'A cache that helps under load rather than failing under it.',
    url: 'https://gritframework.dev/guides/go-redis-caching-stampede',
    type: 'article',
  },
}

export default function Page() {
  return (
    <GuidePage
      slug="go-redis-caching-stampede"
      lead={
        <>
          <p className="mb-4">
            A cache is a correctness problem disguised as a performance improvement. The naive
            version works perfectly in development and has two failure modes that only appear
            under load, both of which make the system slower than having no cache at all.
          </p>
          <p>
            This covers cache-aside in Go, the stampede that follows an expiry, why a fixed TTL
            synchronises your misses, and the invalidation strategies in order of how much
            trouble they cause.
          </p>
        </>
      }
    >
      <Outcome>
        <p>
          A <code>Remember</code> helper with single-flight protection and jittered expiry.
        </p>
        <p>An answer to invalidation that does not rely on remembering every call site.</p>
        <p>A clear idea of what should not be cached, which is most things.</p>
      </Outcome>

      <Section id="measure" title="Measure before you cache">
        <p>
          Caching is the second answer. The first is an index, and it is right far more often
          than its reputation suggests: a query doing a sequential scan on a million rows is a
          missing index, not a caching opportunity, and caching it hides the problem while adding
          staleness.
        </p>
        <Checklist
          items={[
            <>
              <strong>Is it slow because of an N+1?</strong> Fix the query. A cached N+1 is five
              hundred cache lookups.
            </>,
            <>
              <strong>Is it slow because of a missing index?</strong> Add the index. It is exact,
              it has no staleness, and it helps every query against that column.
            </>,
            <>
              <strong>Is it genuinely expensive and read far more than it is written?</strong>{' '}
              Now cache it.
            </>,
          ]}
        />
      </Section>

      <Section id="aside" title="Cache-aside, and why the simple version is wrong">
        <CodeBlock
          code={`// The version in every tutorial.
func (s *Service) Product(ctx context.Context, id uint) (*Product, error) {
	if hit, err := s.cache.Get(ctx, key(id)); err == nil {
		return decode(hit)
	}
	p, err := s.repo.ByID(ctx, id)
	if err != nil {
		return nil, err
	}
	s.cache.Set(ctx, key(id), encode(p), 10*time.Minute)
	return p, nil
}`}
        />
        <p>Two problems, neither visible until there is traffic.</p>
        <H3>The stampede</H3>
        <p>
          The key expires. Five hundred requests arrive in the same second, all miss, and all run
          the expensive query. The database, which was comfortable because the cache was
          absorbing the reads, now receives five hundred of them at once. The slower it gets, the
          longer the window during which more requests miss.
        </p>
        <p>
          This is why a cache can make an outage worse rather than better: the system is sized
          for the cached load and the cache is exactly what disappears under pressure.
        </p>
        <H3>Synchronised expiry</H3>
        <p>
          Warm a thousand keys at deploy time with a fixed ten-minute TTL and all thousand expire
          in the same second, ten minutes later. The cache turns a steady load into a sawtooth,
          with a spike at a predictable interval.
        </p>
      </Section>

      <Section id="remember" title="One helper that gets both right">
        <CodeBlock
          filename="internal/cache/remember.go"
          code={`package cache

import (
	"context"
	"math/rand"
	"time"

	"golang.org/x/sync/singleflight"
)

type Cache struct {
	rdb    redis.UniversalClient
	group  singleflight.Group
}

// Remember returns the cached value, or computes it once and caches it.
//
// "Once" is the important word: singleflight collapses concurrent calls
// for the same key in this process into one execution of fn, and gives
// every caller the same result. Five hundred simultaneous misses become
// one query.
func Remember[T any](
	ctx context.Context, c *Cache, key string, ttl time.Duration,
	fn func(context.Context) (T, error),
) (T, error) {
	var zero T

	if raw, err := c.rdb.Get(ctx, key).Bytes(); err == nil {
		var out T
		if json.Unmarshal(raw, &out) == nil {
			return out, nil
		}
		// A value that will not decode is a value from an older shape.
		// Drop it rather than failing the request forever.
		c.rdb.Del(ctx, key)
	}

	v, err, _ := c.group.Do(key, func() (any, error) {
		val, err := fn(ctx)
		if err != nil {
			return nil, err
		}
		if raw, err := json.Marshal(val); err == nil {
			// Jitter: up to 20% on top of the TTL, so keys written
			// together do not expire together.
			c.rdb.Set(ctx, key, raw, withJitter(ttl))
		}
		return val, nil
	})
	if err != nil {
		return zero, err
	}
	return v.(T), nil
}

func withJitter(d time.Duration) time.Duration {
	return d + time.Duration(rand.Int63n(int64(d/5)))
}`}
        />
        <Warn title="singleflight is per process">
          <p>
            It collapses concurrent calls inside one instance. Ten instances still produce ten
            queries rather than five hundred, which is usually fine. When it is not, take a short
            distributed lock in Redis around the recompute, and have the losers wait briefly and
            re-read rather than computing. Do not reach for that first: it adds a failure mode,
            the lock holder dying, that singleflight does not have.
          </p>
        </Warn>
      </Section>

      <Section id="degrade" title="A cache that is down must not take the service with it">
        <CodeBlock
          code={`// Every cache call has a short timeout of its own, separate from the
// request's. A Redis that is slow should cost two milliseconds and a
// cache miss, not the whole request.
cctx, cancel := context.WithTimeout(ctx, 50*time.Millisecond)
defer cancel()

raw, err := c.rdb.Get(cctx, key).Bytes()
if err != nil && !errors.Is(err, redis.Nil) {
	// Log it as a miss, not as a failure. The source of truth is the
	// database, and it is still there.
	metrics.CacheError.Inc()
}`}
        />
        <p>
          The rule: a cache error is a cache miss. Code that returns a 500 because Redis is
          unavailable has made an optional dependency into a required one, which is the opposite
          of what a cache is for.
        </p>
      </Section>

      <Section id="invalidation" title="Invalidation, in order of how much trouble it causes">
        <H3>1. A short TTL and nothing else</H3>
        <p>
          Boring and correct. If sixty seconds of staleness is acceptable, this is the whole
          design and there is no invalidation code to get wrong. Start here and only move when
          somebody states a requirement that rules it out.
        </p>
        <H3>2. Delete on write</H3>
        <CodeBlock
          code={`func (s *ProductService) Update(ctx context.Context, p *Product) error {
	if err := s.repo.Save(ctx, p); err != nil {
		return fmt.Errorf("save product: %w", err)
	}
	// Delete, do not overwrite. Overwriting with the value you just
	// wrote races with a concurrent update and can leave the older of
	// the two in the cache.
	return s.cache.Del(ctx, productKey(p.ID), productListKey(p.CategoryID))
}`}
        />
        <p>
          The difficulty is remembering every key a write affects, and lists are where it breaks:
          a product appears in its own key, in a category listing, in a search result and in a
          dashboard count.
        </p>
        <H3>3. Version the key instead of deleting it</H3>
        <CodeBlock
          code={`// Bump a counter and every key derived from it is unreachable at once.
// No key enumeration, no scanning, and the old values expire on their
// own TTL.
func (c *Cache) listKey(ctx context.Context, categoryID uint) string {
	v, _ := c.rdb.Get(ctx, fmt.Sprintf("ver:category:%d", categoryID)).Int64()
	return fmt.Sprintf("products:cat:%d:v%d", categoryID, v)
}

func (c *Cache) BumpCategory(ctx context.Context, categoryID uint) error {
	return c.rdb.Incr(ctx, fmt.Sprintf("ver:category:%d", categoryID)).Err()
}`}
        />
        <p>
          This is the technique worth knowing. It turns &quot;delete every key that might contain
          this&quot; into one increment, and it is safe under concurrency because the version is
          part of the key rather than something checked separately.
        </p>
        <Warn title="Never KEYS, and be careful with SCAN">
          <p>
            <code>KEYS pattern</code> blocks Redis while it walks the entire keyspace. On a
            production instance that is a stall measured in seconds, during which every other
            command waits. <code>SCAN</code> is incremental and safe, and is still a poor
            invalidation strategy: key versioning avoids needing either.
          </p>
        </Warn>
        <H3>4. Pub/sub, for local caches</H3>
        <p>
          Only when you have an in-process cache in front of Redis. One instance writes, all
          instances hear about it and drop their local copy. It is a real technique with a real
          cost: a dropped message means one instance is stale indefinitely, so the local copy
          needs a short TTL anyway, which raises the question of whether it was worth it.
        </p>
      </Section>

      <Section id="keys" title="Keys, and the mistakes in them">
        <CodeBlock
          code={`// Namespaced, versioned, and the scope is IN the key.
"app:v2:product:8821"
"app:v2:org:412:invoices:page:1:status:draft"
"app:v2:user:88:permissions"`}
        />
        <Checklist
          items={[
            <>
              <strong>Put the tenant in the key.</strong> A cached list keyed by page number and
              not by organisation serves one customer&apos;s data to another, and it is a
              particularly bad bug because it is intermittent and looks like a database problem.
            </>,
            <>
              <strong>Put the version in the prefix.</strong> When the cached shape changes,
              bumping <code>v2</code> to <code>v3</code> invalidates everything instantly, which
              is otherwise a deploy that reads old structures into new types.
            </>,
            <>
              <strong>Include every input that changes the answer:</strong> filters, sort, page
              size, and the viewer&apos;s permissions if those affect what is returned.
            </>,
            <>
              <strong>Never cache a response that varies by user under a key that does not.</strong>{' '}
              This is the same mistake as the tenant one, and it is how a CDN serves one
              person&apos;s dashboard to everybody.
            </>,
          ]}
        />
      </Section>

      <Section id="what" title="What to cache, and what not to">
        <H3>Good candidates</H3>
        <Checklist
          items={[
            <>Expensive aggregates: counts, reports, dashboard figures.</>,
            <>Reference data that changes rarely: categories, settings, feature flags, plans.</>,
            <>Responses from third parties, which are slow and rate limited.</>,
            <>Permission sets, which are read on every request and written rarely.</>,
          ]}
        />
        <H3>Poor candidates</H3>
        <Checklist
          items={[
            <>
              <strong>Anything a user has just changed.</strong> Seeing a stale version of your
              own edit is the worst kind of staleness, because it reads as a bug rather than as a
              delay.
            </>,
            <>
              <strong>Data that is already fast.</strong> A primary key lookup on a warm table is
              under a millisecond. Caching it adds a network hop and a staleness window to save
              nothing.
            </>,
            <>
              <strong>Anything where stale means wrong:</strong> stock levels at checkout,
              account balances, permission denials.
            </>,
            <>
              <strong>Large objects read occasionally.</strong> They evict the small, hot entries
              that were doing the real work.
            </>,
          ]}
        />
      </Section>

      <Section id="operating" title="Operating it">
        <Checklist
          items={[
            <>
              <strong>Measure the hit rate per key prefix,</strong> not overall. An overall rate
              of 85% can hide a prefix at 5% that is pure overhead.
            </>,
            <>
              <strong>Set <code>maxmemory</code> and an eviction policy.</strong>{' '}
              <code>allkeys-lru</code> for a pure cache. Without a policy, Redis eventually
              refuses writes instead of evicting, which is a confusing outage.
            </>,
            <>
              <strong>Know whether you are also using Redis as a queue.</strong> If so, do not
              use <code>allkeys-lru</code>, because it will evict your jobs. Separate the
              instances, or use <code>volatile-lru</code> and set TTLs only on cache keys.
            </>,
            <>
              <strong>Watch the key count and the memory, not just the hit rate.</strong> A
              growing key count with a steady hit rate usually means keys that include something
              unbounded, a timestamp or a search string, which is a slow memory leak.
            </>,
          ]}
        />
      </Section>

      <Section id="with-grit" title="Generated">
        <CodeBlock
          code={`// Remember with single-flight and jitter, ready to use.
products, err := cache.Remember(ctx, c, key, 10*time.Minute,
	func(ctx context.Context) ([]Product, error) {
		return repo.ListFeatured(ctx)
	})`}
        />
        <CodeBlock
          terminal
          code={`grit add cache      # Redis client, Remember, the middleware and the keys`}
        />
        <p>
          There is also response caching as middleware for public read endpoints, keyed on the
          full URL. One caution that comes from this project&apos;s own history: that cache is
          keyed by URL and not by viewer, so an endpoint that varies by who is asking must not be
          on it. A writing or personalised endpoint on a URL-keyed cache group is a data leak
          wearing a performance improvement.
        </p>
        <p>
          <Link href="/docs/batteries/caching" className="text-accent hover:underline">
            The caching reference
          </Link>{' '}
          covers the helper, the middleware and the invalidation helpers.
        </p>
      </Section>

      <Faq
        slug="go-redis-caching-stampede"
        items={[
          {
            q: 'In-memory or Redis?',
            a: (
              <>
                In-memory is faster and is per instance, so each one has a different view and
                invalidation means telling all of them. Redis is shared, which makes invalidation
                a single delete, at the cost of a network hop. Start with Redis: a consistent
                cache is easier to reason about than a fast one, and the hop is well under a
                millisecond.
              </>
            ),
            text:
              'In-memory is faster but per instance, so each has a different view and invalidation must reach all of them. Redis is shared, making invalidation one delete at the cost of a sub-millisecond hop. Start with Redis, because a consistent cache is easier to reason about than a fast one.',
          },
          {
            q: 'How long should a TTL be?',
            a: (
              <>
                As long as the staleness you can defend to the person who will notice. Write it
                down as a sentence: &quot;the dashboard may be up to five minutes behind&quot;. If
                that sentence is acceptable, the TTL is five minutes. If no sentence is
                acceptable, the thing should not be cached. Add jitter regardless.
              </>
            ),
            text:
              'As long as the staleness you can defend to whoever notices. Write it as a sentence, such as the dashboard may be up to five minutes behind: if that is acceptable, the TTL is five minutes. If no such sentence is acceptable, do not cache it. Add jitter regardless.',
          },
          {
            q: 'Is singleflight enough, or do I need a distributed lock?',
            a: (
              <>
                Enough for almost everything. With ten instances, a popular key expiring produces
                ten queries rather than five hundred, and ten is not a stampede. A distributed
                lock is worth it when the recompute is genuinely heavy, tens of seconds, and it
                brings its own failure mode in the lock holder dying mid-compute.
              </>
            ),
            text:
              'Enough for almost everything: with ten instances, an expiring popular key produces ten queries rather than five hundred. A distributed lock is worth it only when the recompute takes tens of seconds, and it adds the failure mode of the lock holder dying mid-compute.',
          },
          {
            q: 'Why is my cache making things slower?',
            a: (
              <>
                Three usual causes. A hit rate low enough that most requests pay the lookup and
                then query anyway. Values large enough that serialising and transferring them
                costs more than the query did. Or a cache call with no timeout of its own, so a
                struggling Redis adds its latency to every request rather than being skipped.
              </>
            ),
            text:
              'Three usual causes: a hit rate low enough that most requests pay the lookup and query anyway, values large enough that serialising and transferring them costs more than the query, or cache calls without their own timeout, so a struggling Redis adds latency to every request instead of being skipped.',
          },
        ]}
      />

      <Related
        links={[
          {
            href: '/guides/go-horizontal-scaling',
            label: 'Running Go behind a load balancer',
            blurb: 'Where a shared cache stops being optional.',
          },
          {
            href: '/guides/gorm-tutorial-relationships',
            label: 'GORM relationships and N+1',
            blurb: 'The query problem to fix before reaching for a cache.',
          },
          {
            href: '/docs/batteries/caching',
            label: 'Caching reference',
            blurb: 'The generated helper, the middleware and the invalidation helpers.',
          },
          {
            href: '/docs/backend/pulse',
            label: 'Profiling with Pulse',
            blurb: 'Finding out what is actually slow before caching it.',
          },
        ]}
      />
    </GuidePage>
  )
}
