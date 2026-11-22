import Link from 'next/link'
import type { Metadata } from 'next'
import { GuidePage, Section, H3, Outcome, Note, Warn, Checklist, Faq, Related } from '@/components/guide'
import { CodeBlock } from '@/components/code-block'

export const metadata: Metadata = {
  title: 'Running Go behind a load balancer',
  description:
    'The stateless checklist, WebSockets over a Redis backplane, cron on exactly one node, read replicas and read-your-writes, and connection pool arithmetic.',
  alternates: { canonical: 'https://gritframework.dev/guides/go-horizontal-scaling' },
  openGraph: {
    title: 'Running Go behind a load balancer',
    description: 'What breaks when one instance becomes three, and what to do about each.',
    url: 'https://gritframework.dev/guides/go-horizontal-scaling',
    type: 'article',
  },
}

export default function Page() {
  return (
    <GuidePage
      slug="go-horizontal-scaling"
      lead={
        <>
          <p className="mb-4">
            Going from one instance to three breaks a specific and predictable set of things:
            anything held in memory, anything on local disk, anything that assumed it was the
            only one running, and anything that counted on a user reaching the same process
            twice.
          </p>
          <p>
            None of it is hard once named. This is the list, in the order you will meet it, with
            the fix for each and the arithmetic that decides how many instances you can actually
            run.
          </p>
        </>
      }
    >
      <Outcome>
        <p>A checklist you can run against an existing service before adding the second node.</p>
        <p>WebSockets that work when the sender and the recipient are on different instances.</p>
        <p>Scheduled work that happens once rather than once per instance.</p>
      </Outcome>

      <Section id="stateless" title="The stateless checklist">
        <Checklist
          items={[
            <>
              <strong>Sessions.</strong> In-memory session storage means a user is signed in on
              one instance and anonymous on the other two. Use a signed token or a shared store;
              sticky sessions are a workaround that fails the moment an instance restarts.
            </>,
            <>
              <strong>Uploaded files.</strong> A file written to local disk exists on one
              instance. The next request lands elsewhere and gets a 404. Object storage, always,
              from the first deploy.
            </>,
            <>
              <strong>In-memory caches and rate limit counters.</strong> Three instances means
              three counters, so a limit of 100 is actually 300. Redis, or accept the
              multiplication knowingly.
            </>,
            <>
              <strong>Generated files.</strong> A PDF written to <code>/tmp</code> and served
              from a second request is the same problem as uploads, and is easier to miss because
              it works in testing.
            </>,
            <>
              <strong>Background goroutines started at boot.</strong> They now run on every
              instance. Sometimes harmless, sometimes three emails.
            </>,
            <>
              <strong>Secrets and configuration from the environment,</strong> identical across
              instances. A per-instance JWT secret means a token issued by one is rejected by the
              others, which presents as random sign-outs.
            </>,
          ]}
        />
        <Note title="The quick test">
          <p>
            Start two instances locally on different ports, put anything in front of them that
            alternates, and use the application. Every one of the above fails within a few
            minutes, and failing locally is considerably cheaper than failing after the second
            node is deployed.
          </p>
        </Note>
      </Section>

      <Section id="shutdown" title="Graceful shutdown, which rolling deploys require">
        <CodeBlock
          code={`srv := &http.Server{Addr: cfg.Addr, Handler: r}

go func() {
	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal().Err(err).Msg("listen")
	}
}()

stop := make(chan os.Signal, 1)
signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
<-stop

// Stop reporting ready BEFORE closing, so the load balancer removes
// this instance from rotation while it is still serving. Without this
// gap, requests are routed to a process that is already shutting down,
// and the deploy produces a handful of 502s every time.
readiness.Store(false)
time.Sleep(cfg.DrainDelay)      // one or two health check intervals

ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()
if err := srv.Shutdown(ctx); err != nil {
	log.Error().Err(err).Msg("forced shutdown")
}`}
        />
        <p>
          The drain delay is the part that is usually missing. <code>Shutdown</code> stops
          accepting new connections immediately, but the load balancer does not know that until
          its next health check, and everything it routes in between fails.
        </p>
        <CodeBlock
          code={`// Two endpoints, two meanings. Conflating them is why a database
// outage can get every instance restarted at once.
r.GET("/healthz", func(c *gin.Context) {
	// Liveness: is this process working? Do NOT check the database:
	// a database outage would then restart every instance, which helps
	// nobody and loses all your warm connections.
	c.Status(http.StatusOK)
})

r.GET("/readyz", func(c *gin.Context) {
	// Readiness: should traffic come here right now? This one does
	// check dependencies, and says no while draining.
	if !readiness.Load() {
		c.Status(http.StatusServiceUnavailable)
		return
	}
	ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
	defer cancel()
	if err := db.WithContext(ctx).Exec("SELECT 1").Error; err != nil {
		c.Status(http.StatusServiceUnavailable)
		return
	}
	c.Status(http.StatusOK)
})`}
        />
      </Section>

      <Section id="websockets" title="WebSockets across instances">
        <p>
          A WebSocket is a connection to <em>one</em> process. User A connects to instance 1,
          user B to instance 3, and a message from A reaches nobody, because instance 1 holds no
          connection to B.
        </p>
        <p>
          The fix is a backplane: publish to Redis, every instance subscribes, each delivers to
          the connections it happens to hold.
        </p>
        <CodeBlock
          filename="internal/realtime/hub.go"
          code={`// Local connections, by room.
type Hub struct {
	mu    sync.RWMutex
	rooms map[string]map[*Conn]struct{}
	rdb   redis.UniversalClient
}

// Publish goes to Redis, never directly to local connections, so that
// one code path serves both the local and the remote case. A fast path
// that skips Redis for local recipients is where "it works for some
// users" bugs come from.
func (h *Hub) Publish(ctx context.Context, room string, msg []byte) error {
	return h.rdb.Publish(ctx, "rt:"+room, msg).Err()
}

// Every instance runs this, and delivers to whatever it holds.
func (h *Hub) subscribe(ctx context.Context, room string) {
	sub := h.rdb.Subscribe(ctx, "rt:"+room)
	for m := range sub.Channel() {
		h.mu.RLock()
		for c := range h.rooms[room] {
			select {
			case c.send <- []byte(m.Payload):
			default:
				// A client that is not reading must not block the hub.
				// Drop it: a slow consumer that holds the lock is how
				// one bad connection stalls every other one.
				close(c.send)
			}
		}
		h.mu.RUnlock()
	}
}`}
        />
        <Checklist
          items={[
            <>
              <strong>Redis pub/sub is fire and forget.</strong> A message published while an
              instance is restarting is lost. For anything that must arrive, write it to the
              database first and have the client reconcile on reconnect.
            </>,
            <>
              <strong>Raise the idle timeout on the load balancer,</strong> or it will close
              WebSockets after sixty seconds and your users will see constant reconnects.
            </>,
            <>
              <strong>Send pings.</strong> Intermediaries drop idle connections silently, and the
              client finds out only when a message fails to arrive.
            </>,
            <>
              <strong>Consider server-sent events instead.</strong> If the traffic is one-way,
              SSE is plain HTTP, survives proxies that mangle WebSocket upgrades, and reconnects
              by itself.
            </>,
          ]}
        />
      </Section>

      <Section id="cron" title="Scheduled work, exactly once">
        <p>
          Run the scheduler inside the application and scale to three instances, and every
          nightly job runs three times. With an email, that is visible. With a billing run, it is
          an incident.
        </p>
        <H3>Option one: a separate single-instance process</H3>
        <p>
          The simplest correct answer. One deployment, one replica, running only the scheduler.
          It is a single point of failure for scheduled work, and for most systems a job running
          late is much better than a job running three times.
        </p>
        <H3>Option two: leader election</H3>
        <CodeBlock
          code={`// Whoever holds the lock runs the schedule. The lock expires, so a
// dead leader is replaced within one TTL rather than never.
func (l *Leader) campaign(ctx context.Context) {
	t := time.NewTicker(5 * time.Second)
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			// SET NX EX: set only if absent, with an expiry. The expiry
			// is what makes a crashed leader recoverable.
			ok, err := l.rdb.SetNX(ctx, "leader:cron", l.id, 15*time.Second).Result()
			if err == nil && ok {
				l.become(true)
				continue
			}
			// Already leader? Extend, but only if the value is still
			// ours, which needs a script to be atomic.
			if l.isLeader() && l.extend(ctx) == nil {
				continue
			}
			l.become(false)
		}
	}
}`}
        />
        <Warn title="A lock with a TTL is not a guarantee">
          <p>
            A leader that pauses for longer than the TTL, during a long garbage collection or a
            paused container, can still believe it is leader while another instance has taken
            over. For scheduling, a rare double run is usually acceptable. For anything where it
            is not, make the job itself idempotent, which is the real defence, and treat the lock
            as an optimisation.
          </p>
        </Warn>
        <H3>Option three: let the queue arbitrate</H3>
        <p>
          Enqueue with a deterministic id, <code>digest:2026-11-22</code>, and let the unique
          constraint reject the duplicates. Every instance may enqueue; only one job exists. This
          is the most robust of the three because the guarantee is in the database rather than in
          a timing assumption.
        </p>
      </Section>

      <Section id="pool" title="The arithmetic that actually limits you">
        <p>
          Postgres has a connection limit, often around a hundred. Every instance has its own
          pool, and so does every worker.
        </p>
        <CodeBlock
          language="text"
          code={`6 API instances   × 25 connections = 150
3 worker instances × 10 connections =  30
                                      ---
                                      180   against max_connections = 100

Half of your instances cannot connect, intermittently, and the error
says "too many clients" at a moment unrelated to the deploy that caused
it.`}
        />
        <CodeBlock
          code={`sqlDB, err := db.DB()
if err != nil {
	return fmt.Errorf("sql handle: %w", err)
}
// Size the pool from max_connections divided by the number of
// processes, with headroom for migrations and for somebody connecting
// with psql during an incident.
sqlDB.SetMaxOpenConns(cfg.DBMaxOpen)       // not the default of unlimited
sqlDB.SetMaxIdleConns(cfg.DBMaxOpen / 2)
sqlDB.SetConnMaxLifetime(30 * time.Minute) // so a failover is noticed
sqlDB.SetConnMaxIdleTime(5 * time.Minute)`}
        />
        <p>
          Past a handful of instances, put PgBouncer in transaction mode in front. It multiplexes
          many application connections onto few database ones, and the cost is that session-level
          state, including <code>SET</code> outside a transaction and prepared statements, stops
          working as expected. That is exactly why{' '}
          <code>set_config(&apos;app.org_id&apos;, ..., true)</code> with the transaction-local
          flag matters if you use row-level security.
        </p>
      </Section>

      <Section id="replicas" title="Read replicas and read-your-writes">
        <p>
          Replicas scale reads and introduce lag, usually milliseconds and occasionally seconds.
          The failure that follows is specific and infuriating: a user saves something, the next
          request reads a replica that has not caught up, and their change has apparently
          vanished.
        </p>
        <CodeBlock
          code={`// A simple and effective rule: after a write, that user reads the
// primary for a short window. Correct, cheap, and obvious to the next
// reader of the code.
func (r *Router) For(ctx context.Context, userID uint) *gorm.DB {
	if r.recentlyWrote(ctx, userID) {
		return r.primary
	}
	return r.replica
}

func (r *Router) MarkWrite(ctx context.Context, userID uint) {
	r.rdb.Set(ctx, fmt.Sprintf("rw:%d", userID), "1", 5*time.Second)
}`}
        />
        <Checklist
          items={[
            <>
              <strong>Everything inside a transaction goes to the primary,</strong> always. A
              read-modify-write split across primary and replica is a lost update.
            </>,
            <>
              <strong>Send reports and exports to a replica deliberately.</strong> They are the
              queries that benefit most and tolerate lag best.
            </>,
            <>
              <strong>Monitor the lag,</strong> and fall back to the primary when it exceeds a
              threshold. A replica minutes behind is worse than no replica.
            </>,
          ]}
        />
      </Section>

      <Section id="order" title="The order to do this in">
        <p>
          Horizontal scaling is usually not the first answer, and reaching for it early buys
          complexity instead of capacity.
        </p>
        <Checklist
          items={[
            <>
              <strong>1. Measure.</strong> Find out what is slow. It is usually a handful of
              queries, and often one.
            </>,
            <>
              <strong>2. Fix the queries.</strong> Indexes and N+1s. This is routinely a tenfold
              improvement for an afternoon&apos;s work.
            </>,
            <>
              <strong>3. Cache what is expensive and read often.</strong>
            </>,
            <>
              <strong>4. Scale vertically.</strong> Doubling a server is one change with no new
              failure modes, and modern machines are very large.
            </>,
            <>
              <strong>5. Then scale horizontally,</strong> which is also what buys you
              availability rather than only throughput.
            </>,
          ]}
        />
        <p>
          A single Go instance on modest hardware serves thousands of requests a second against
          well-indexed queries. If you are not there yet, the problem is not the instance count.
        </p>
      </Section>

      <Section id="with-grit" title="What a generated project already does">
        <p>
          Graceful shutdown with a drain delay, separate liveness and readiness endpoints,
          configured pool limits, Redis-backed sessions and rate limits, object storage rather
          than local disk, and the realtime hub on a Redis backplane.
        </p>
        <CodeBlock
          terminal
          code={`grit scale          # the stateless checklist, run against this project`}
        />
        <p>
          It reports what would break at more than one instance: a rate limiter with an in-memory
          store, a scheduler inside the API binary, a file written to local disk, a pool size
          that multiplies past <code>max_connections</code>.
        </p>
        <p>
          <Link href="/docs/scaling" className="text-accent hover:underline">
            The scaling documentation
          </Link>{' '}
          walks the stages from one box to many.
        </p>
      </Section>

      <Faq
        slug="go-horizontal-scaling"
        items={[
          {
            q: 'When should I add a second instance?',
            a: (
              <>
                For availability, as soon as downtime costs something: one instance means a
                deploy, a crash or a host failure is an outage. For throughput, only after
                measuring, because the usual bottleneck is a query rather than the process. Those
                are two different reasons and only the first is urgent early.
              </>
            ),
            text:
              'For availability, as soon as downtime costs something, since one instance makes every deploy, crash or host failure an outage. For throughput, only after measuring, because the usual bottleneck is a query rather than the process.',
          },
          {
            q: 'Do I need sticky sessions?',
            a: (
              <>
                No, and wanting them is a signal that state is in the wrong place. They fail when
                an instance restarts, make deploys uneven and prevent a load balancer from
                balancing. Move the state to a token or a shared store instead. WebSockets are
                the one legitimate exception, and a backplane handles that better than stickiness
                does.
              </>
            ),
            text:
              'No, and needing them signals state in the wrong place. They fail when an instance restarts, make deploys uneven and stop the balancer balancing. Move state to a token or shared store. WebSockets are the exception, and a Redis backplane handles that better than stickiness.',
          },
          {
            q: 'How many connections should each instance open?',
            a: (
              <>
                Divide the database&apos;s limit by the number of processes and leave headroom
                for migrations and for a person with psql. Twenty-five per instance against a
                hundred is already four instances with nothing spare. Past that, PgBouncer in
                transaction mode, and remember it changes how session-level settings behave.
              </>
            ),
            text:
              'Divide the database limit by the number of processes and leave headroom for migrations and manual access. Twenty-five per instance against a limit of one hundred is already four instances with nothing spare. Past that use PgBouncer in transaction mode, noting it changes how session-level settings behave.',
          },
          {
            q: 'Why does my cron job run three times?',
            a: (
              <>
                Because the scheduler is inside the application binary and you have three of
                them. Run the scheduler as its own single-replica process, elect a leader, or
                enqueue with a deterministic id so the unique constraint rejects the duplicates.
                The third is the most robust, because the guarantee is in the database rather
                than in a timing assumption.
              </>
            ),
            text:
              'Because the scheduler is inside the application binary and there are three of them. Run it as its own single-replica process, elect a leader, or enqueue with a deterministic id so a unique constraint rejects duplicates. The last is most robust, since the guarantee is in the database rather than a timing assumption.',
          },
        ]}
      />

      <Related
        links={[
          {
            href: '/guides/go-redis-caching-stampede',
            label: 'Caching and stampede protection',
            blurb: 'The shared cache that replaces the in-memory one.',
          },
          {
            href: '/guides/deploy-go-app-vps-https',
            label: 'Deploy to a VPS',
            blurb: 'The single-instance deployment this is the next step from.',
          },
          {
            href: '/docs/scaling',
            label: 'Scaling',
            blurb: 'The stages from one box to many, with what changes at each.',
          },
          {
            href: '/docs/backend/realtime',
            label: 'Realtime',
            blurb: 'The hub, the backplane and reconnection.',
          },
        ]}
      />
    </GuidePage>
  )
}
