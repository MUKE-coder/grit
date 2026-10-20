#import "../template.typ": *

= Build your own job queue

Chapter forty three used a queue. This chapter builds one, on the database you
already have, which for most applications is the right answer: no Redis, no
second thing to operate, and jobs that can be enqueued *inside the same
transaction* as the work that caused them.

The hard parts are not the obvious ones. Storing a job is a row. Handing that
row to exactly one worker, and getting it back when that worker dies, is the
engineering.

== The table

```sql
CREATE TABLE jobs (
  id           BIGSERIAL PRIMARY KEY,
  queue        TEXT        NOT NULL DEFAULT 'default',
  type         TEXT        NOT NULL,
  payload      JSONB       NOT NULL,

  state        TEXT        NOT NULL DEFAULT 'pending',
  attempts     INT         NOT NULL DEFAULT 0,
  max_attempts INT         NOT NULL DEFAULT 5,

  run_at       TIMESTAMPTZ NOT NULL DEFAULT now(),
  locked_at    TIMESTAMPTZ,
  locked_by    TEXT,
  last_error   TEXT,
  unique_key   TEXT,
  created_at   TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX jobs_claim ON jobs (queue, state, run_at)
  WHERE state = 'pending';

CREATE UNIQUE INDEX jobs_unique_key ON jobs (unique_key)
  WHERE unique_key IS NOT NULL AND state <> 'done';
```

Two index decisions carry most of the performance.

*The partial index.* Claiming only ever looks for pending rows, and a finished
job is dead weight in the index. A partial index stays small even when the
table holds millions of completed jobs.

*The partial unique index* is idempotency: a job with a `unique_key` that is
already pending or running cannot be enqueued twice, and the database enforces
it rather than the application checking first and racing.

== Enqueue, inside your transaction

```go
func (q *Queue) Enqueue(ctx context.Context, tx *sql.Tx,
    j Job) error {

    _, err := tx.ExecContext(ctx, `
        INSERT INTO jobs (queue, type, payload, run_at, unique_key)
        VALUES ($1, $2, $3, $4, $5)
        ON CONFLICT DO NOTHING`,
        j.Queue, j.Type, j.Payload, j.RunAt, j.UniqueKey)
    return err
}
```

#bridge[
  That `tx` parameter is the whole reason to put a queue in your database.

  With Redis, enqueueing and committing are two systems, which is chapter
  forty four's publish-or-commit problem and needs an outbox. Here the job row
  and the business row commit together or not at all: the outbox *is* the
  queue.

  It is the strongest argument for a database-backed queue, and it disappears
  the moment your throughput needs something specialised.
]

== Claiming work, which is the whole problem

Several workers poll the same table. Two must never get the same row, and a
worker waiting on a lock is a worker doing nothing.

```sql
UPDATE jobs SET
  state     = 'running',
  locked_at = now(),
  locked_by = $1,
  attempts  = attempts + 1
WHERE id IN (
  SELECT id FROM jobs
  WHERE queue = $2 AND state = 'pending' AND run_at <= now()
  ORDER BY run_at
  LIMIT $3
  FOR UPDATE SKIP LOCKED
)
RETURNING id, type, payload, attempts;
```

`FOR UPDATE SKIP LOCKED` is the line that makes this work, and it exists in
Postgres for exactly this use case.

#gotcha[
  Without `SKIP LOCKED`, ten workers all select the same oldest rows, nine
  block on the first one's lock, and your queue processes jobs one at a time
  with nine idle workers holding transactions open.

  With it, each worker takes the rows nobody else has locked and moves on.

  MySQL 8 has it too. SQLite does not, and serialises writers anyway, so there
  a single claiming goroutine handing work to a pool is the honest design.
]

Claim a batch rather than one row. The round trip dominates at small sizes, and
ten or twenty jobs per claim is usually an order of magnitude more throughput.

== When a worker dies

The job is `running`, `locked_by` names a process that no longer exists, and
nothing will ever finish it. The lease is what fixes this:

```sql
-- The reaper, every 30 seconds
UPDATE jobs
SET state = 'pending', locked_at = NULL, locked_by = NULL
WHERE state = 'running'
  AND locked_at < now() - interval '5 minutes';
```

And a worker on a long job keeps its lease alive:

```go
func (w *Worker) heartbeat(ctx context.Context, id int64) {
    t := time.NewTicker(30 * time.Second)
    defer t.Stop()
    for {
        select {
        case <-ctx.Done():
            return
        case <-t.C:
            w.db.ExecContext(ctx,
                `UPDATE jobs SET locked_at = now() WHERE id = $1`,
                id)
        }
    }
}
```

#gotcha[
  The lease is why delivery is *at least once* and can never be exactly once.

  A worker that finishes the work and dies before marking the job done will
  have that job reclaimed and run again. There is no ordering of those two
  steps that removes the window, which is the same impossibility chapter forty
  three described for mail.

  So the handler must be idempotent, and the lease duration must exceed your
  longest job. A five minute lease on a job that takes ten minutes means the
  job is reclaimed and run again *while it is still running*, which is the
  worst outcome available: two workers doing the same work at the same time.
]

== Failure, backoff and the dead letter

```go
func (w *Worker) fail(ctx context.Context, j Job, err error) {
    if j.Attempts >= j.MaxAttempts {
        w.db.ExecContext(ctx, `
            UPDATE jobs SET state = 'dead', last_error = $1
            WHERE id = $2`, err.Error(), j.ID)
        return
    }
    delay := backoff(j.Attempts)
    w.db.ExecContext(ctx, `
        UPDATE jobs SET state = 'pending', run_at = $1,
            last_error = $2, locked_at = NULL, locked_by = NULL
        WHERE id = $3`, time.Now().Add(delay), err.Error(), j.ID)
}

func backoff(attempt int) time.Duration {
    base := time.Duration(math.Pow(2, float64(attempt))) *
        time.Second
    if base > 10*time.Minute {
        base = 10 * time.Minute
    }
    // Jitter, and it is not optional.
    jitter := time.Duration(rand.Int63n(int64(base / 2)))
    return base/2 + jitter
}
```

#gotcha[
  The jitter is the part people skip, and it is the part that matters during an
  incident.

  A downstream service goes down for a minute. Five thousand jobs fail at once.
  Without jitter, all five thousand retry at exactly the same instant, four
  seconds later, and again sixteen seconds later. The service comes back, gets
  hit by five thousand simultaneous requests, and goes down again. You have
  built a mechanism that prevents recovery.

  Spread the retries. Full jitter, half jitter, anything: the shape matters far
  more than the formula.
]

A dead job is not deleted. It sits in the table with its error, where somebody
can look at it, fix the cause and requeue it. A queue with no dead letter state
loses work silently, which is the one thing a queue must never do.

== Polling, and the latency it costs

Polling every second means up to a second of latency for every job. For a
reminder email that is irrelevant; for a user waiting on a thumbnail it is not.

Postgres gives you a way out:

```sql
NOTIFY jobs_default;
```

```go
// The worker waits on the notification OR a timer, whichever
// comes first. The timer is the floor: notifications can be
// missed, and polling is the thing that makes it correct.
select {
case <-notified:
case <-time.After(pollInterval):
}
```

Notifications are an optimisation, not a guarantee. Keeping the poll as a floor
means a missed notification costs latency rather than a stuck job.

== Fairness, and the queue that starves

With several queues, the naive loop takes everything from the first queue
before looking at the second, so a flood of image jobs stops every email.

Weighted selection, by starting from a random position in a weighted list, is
enough: `default` weight 5, `email` 3, `images` 2. Over any reasonable window
each queue gets its share, and no queue can be starved by another.

== Knowing it works

Four numbers, and the second is the one to alert on:

#table(
  columns: (auto, 1fr),
  stroke: 0.4pt + rgb("#e3e6ea"),
  inset: 5pt,
  [Counts by state], [Pending, running, dead. The shape of the backlog],
  [Oldest pending age], [*The single best queue alert there is*],
  [Duration by type], [Which jobs are getting slower],
  [Dead per hour], [A rise means something downstream broke],
)

Depth alone is misleading: ten thousand jobs enqueued a second ago is fine, and
one job pending for an hour is an outage. Age catches both.

== Shutting down without losing work

```go
func (w *Worker) Shutdown(ctx context.Context) error {
    w.stopClaiming()                 // no new jobs
    done := make(chan struct{})
    go func() { w.wg.Wait(); close(done) }()

    select {
    case <-done:
        return nil
    case <-ctx.Done():
        // Out of time: release the leases we hold so another
        // worker picks them up now rather than in five minutes.
        return w.releaseHeldJobs()
    }
}
```

That last branch is worth the few lines. A deploy that releases its leases
makes the next worker pick those jobs up immediately, instead of them sitting
`running` until the reaper notices.

== When this is the wrong design

Be honest about the ceiling. A database-backed queue handles a few thousand
jobs a second on decent hardware, which is more than most applications ever
produce. Beyond that, the claim query becomes the bottleneck and the table
becomes a hot spot.

#table(
  columns: (auto, 1fr),
  stroke: 0.4pt + rgb("#e3e6ea"),
  inset: 5pt,
  [Database queue], [You already have a database, you want transactional
    enqueue, and volume is moderate],
  [Redis queue], [High throughput, low latency, jobs that are not tied to a
    transaction],
  [Kafka or similar], [Streams, many consumers, replay, retention. A
    different tool for a different problem],
)

#project(title: "Project 53: a queue you would trust")[
  In `part6-from-scratch/57-queue`, build it against Postgres in Docker.

  *One: the schema and enqueue*, with the partial indexes and transactional
  enqueue taking a `*sql.Tx`.

  *Two: claiming* with `FOR UPDATE SKIP LOCKED`, batched. Run eight workers
  against ten thousand jobs and assert every job ran exactly once.

  *Three: prove `SKIP LOCKED` matters.* Remove it, run the same test, and
  record the throughput difference and what the workers were doing.

  *Four: the dead worker.* Kill a worker mid-job with `kill -9`. Confirm the
  job is stuck `running`, then watch the reaper return it. Measure the
  recovery time.

  *Five: the lease trap.* Set the lease to ten seconds and run a job that takes
  thirty. Observe two workers running it at once, and write down what that
  would have done if the job charged a card.

  *Six: backoff.* Make a job fail every time and plot its retry times. Then
  remove the jitter, fail five thousand jobs at once, and graph the arrivals.
  The spike is the chapter's point.

  *Seven: fairness.* Flood one queue and confirm another still progresses.

  *Eight: observability.* A `/queue/stats` endpoint with the four numbers, and an
  alert rule on oldest pending age. State the threshold and why.
]

#tryit[
  Enqueue a job inside a transaction that then rolls back.

  Nothing should be queued. Now do the same against a Redis queue from chapter
  forty three and watch the job run for a row that does not exist. That one
  experiment is the clearest demonstration of why the outbox exists.
]

#assignment(name: "Job Queue - Assignment One")[
  Build the queue and prove each guarantee with an experiment, not a claim.

  *One.* The full implementation: transactional enqueue, batched claiming with
  `SKIP LOCKED`, leases and a reaper, jittered backoff, a dead letter state,
  weighted fairness and graceful shutdown.

  *Two.* `GUARANTEES.md`, one short section each: exactly once or at least
  once and why, what happens when a worker dies, what happens when a job
  outlives its lease, and what happens when the database is briefly down.

  *Three.* The eight-worker correctness test, with the count of jobs run
  exactly once out of ten thousand.

  *Four.* The jitter graph, before and after, even if it is ASCII.

  *Five.* One paragraph on when you would replace this with Redis or Kafka for
  your own workload, with the number that would make you switch.
]
