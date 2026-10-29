import { CodeBlock } from '@/components/code-block'
import { TipBox } from '@/components/course/tip-box'
import { Exercise } from '@/components/course/exercise'
import { KnowledgeCheck } from '@/components/course/knowledge-check'

export default function Lesson() {
  return (
    <>
      <p>
        Load tests in a Grit API find the same handful of things, over and
        over. Knowing the shortlist does not replace measuring, but it tells
        you where to point the measurement, and it stops you optimising the
        handler when the problem is the pool in front of it.
      </p>

      <h2>One: the connection pool, and the surprise in it</h2>
      <p>
        This is the most common one and the least intuitive, so it is worth
        real numbers. Grit sizes the pool in{' '}
        <code>internal/database/database.go</code>:
      </p>

      <CodeBlock
        language="go"
        filename="apps/api/internal/database/database.go"
        code={`maxOpen := getEnvInt("DB_MAX_OPEN_CONNS", 25)
maxIdle := getEnvInt("DB_MAX_IDLE_CONNS", maxOpen)   // idle defaults to open
sqlDB.SetMaxIdleConns(maxIdle)
sqlDB.SetMaxOpenConns(maxOpen)
sqlDB.SetConnMaxLifetime(30 * time.Minute)
sqlDB.SetConnMaxIdleTime(10 * time.Minute)`}
      />

      <p>
        Idle defaulting to open is the part that matters. When the idle limit
        is lower than the open limit, a request past the idle limit returns
        its connection to a full pool, so the connection is{' '}
        <strong>closed</strong>, and the next request opens a new one, which
        makes Postgres fork a backend process. Under concurrency that is a
        connection storm.
      </p>

      <p>Measured with k6 at 50 VUs, 4 CPUs per container, single row reads:</p>

      <CodeBlock
        language="text"
        code={`idle=10    ~810 req/s     Postgres ~840% CPU    API 196%
idle=100  ~2720 req/s     Postgres ~300% CPU    API ~300%

Same binary. Same query. Same hardware.`}
      />

      <p>
        Three times the throughput from one number, and note where the cost
        showed up: as <em>database</em> CPU. Nothing in the application
        looked wrong, which is why this is the bottleneck people profile
        around for a week.
      </p>

      <TipBox tone="warning">
        <p>
          Bigger is not automatically better. The total across processes is
          what the server sees: every replica opens a pool, and Sentinel
          opens a second one in each, all against{' '}
          <code>max_connections</code>, which is 100 on a default Postgres.
          Size it as{' '}
          <code>
            (max_connections - reserve) / replicas - SENTINEL_DB_MAX_OPEN_CONNS
          </code>
          . The defaults, 25 and 5, are that sum for three replicas with ten
          left over for migrations and psql.
        </p>
      </TipBox>

      <h2>Two: N+1 queries</h2>
      <p>
        One query for the list, then one per row for a relation. At a
        request a second it is invisible. At 100 virtual users it is the
        whole problem, because a 20 item page is 21 round trips and every
        one of them holds a connection from the pool above.
      </p>

      <CodeBlock
        language="go"
        filename="DELIBERATELY BROKEN"
        code={`// 1 query, then N more inside the loop.
db.Find(&posts)
for i := range posts {
    db.Where("post_id = ?", posts[i].ID).Find(&posts[i].Comments)
}

// 2 queries, whatever N is.
db.Preload("Comments").Find(&posts)`}
      />

      <p>
        You do not have to spot these by reading. Pulse counts queries per
        request and flags the pattern, which is the difference between
        knowing a route is slow and knowing it issued 43 queries.
      </p>

      <h2>Three: a cache that is cold, or absent</h2>
      <p>
        Two different faults that look the same in a graph. A cold cache
        gives you a slow first minute and a fast rest of the run, which is
        why the ramp stage is excluded from the baseline. A missing cache
        gives you a flat, uniformly mediocre line that never improves.
      </p>

      <p>
        The one to watch for is the cache that is evicted under load. If the
        working set no longer fits, the hit rate falls exactly when traffic
        rises, and latency goes up faster than load does. The signature is a
        p95 that climbs non linearly with virtual users.
      </p>

      <h2>Four: GC pauses, which are usually not the problem</h2>
      <p>
        Go&apos;s collector is concurrent, and pauses are typically well
        under a millisecond. They are worth suspecting only when p99 is bad
        while p50 is fine <em>and</em> allocation rate is high: a handler
        allocating a large slice per request, JSON encoding a big response,
        reading a whole file into memory.
      </p>

      <p>
        The fix is almost never tuning <code>GOGC</code>. It is allocating
        less, and the next course finds out where with a heap profile.
      </p>

      <h2>Where to look, in order</h2>
      <CodeBlock
        language="text"
        code={`1. Is the DATABASE's CPU high?        -> pool, N+1, or a missing index
2. Is the API's CPU high?             -> serialisation, crypto, your own code
3. Are both low and latency high?     -> waiting: pool exhaustion, a lock,
                                         an external call, a context deadline
4. Is memory climbing?                -> the soak test's territory`}
      />

      <p>
        The third line is the one people skip. Latency with no CPU anywhere
        is the system waiting, and waiting does not show up in a CPU profile
        at all.
      </p>

      <Exercise
        prompt={
          <>
            <p>
              <strong>One.</strong> Run the average load test against a list
              endpoint with a relation and watch Pulse&apos;s queries per
              request for that route.
            </p>
            <p>
              <strong>Two.</strong> Set <code>DB_MAX_IDLE_CONNS=2</code> in{' '}
              <code>.env</code>, restart, and run the same test. Compare
              throughput and database CPU.
            </p>
            <p>
              <strong>Three.</strong> Put it back. Then classify your own
              system against the four line list above.
            </p>
          </>
        }
        hint={
          <p>
            In step two, watch where the cost lands. The interesting part is
            not that it got slower, it is which process got busier.
          </p>
        }
        solution={
          <>
            <p>
              Step one usually produces a bigger number than people expect
              the first time they look. A query count of 21 on a 20 row page
              is the classic, and <code>Preload</code> takes it to 2.
            </p>
            <p>
              Step two reproduces the measurement above: throughput falls,
              and database CPU rises while the API gets less busy. That is
              the signature worth memorising, because every instinct says to
              look at the application and the application is not the problem.
            </p>
          </>
        }
      />

      <KnowledgeCheck
        question="Under load, p95 is 900ms. CPU is 15% on the API container and 20% on the database. What is happening?"
        choices={[
          {
            label: 'Something is waiting: pool exhaustion, a lock, or an external call',
            correct: true,
            feedback:
              'Latency with no CPU anywhere is time spent blocked, and a CPU profile shows nothing because nothing is running. Check pool wait time, then anything the handler calls over a network.',
          },
          {
            label: 'The API needs more CPU',
            feedback:
              'It is using 15% of what it has. More would sit idle alongside it.',
          },
          {
            label: 'GC pauses',
            feedback:
              'A collector busy enough to add 900ms would show up as CPU. At 15% it is not the explanation.',
          },
          {
            label: 'The load generator is saturated',
            feedback:
              'Worth ruling out, and it would make k6 report latency the server never saw. Check the generator CPU, then come back to pool wait time.',
          },
        ]}
      />
    </>
  )
}
