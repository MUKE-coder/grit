import type { SystemDesign } from './systems-types'

export const BACKGROUND_JOBS: SystemDesign = {
  slug: 'background-jobs',
  name: 'Background Jobs',
  tagline:
    'Getting slow work out of the request path, and having it survive a retry, a duplicate and a deploy.',
  group: 'Delivery',
  packages: ['internal/jobs', 'internal/handlers'],

  problem: {
    text: [
      'A request that sends an email, resizes an image, builds a PDF or calls a third-party API is a request held open by something that has nothing to do with the response. The user waits for work they did not ask to wait for, a connection and a database handle are held for seconds, and if the third party is slow then so is the application.',
      'The obvious fix is a goroutine, and it is wrong in four ways that all look fine in development. Nothing bounds it, so a burst of sign-ups is a burst of goroutines each holding a database connection and a provider socket. Nothing retries it, so a provider down for sixty seconds loses that work permanently. A deploy in the middle drops whatever was in flight. And no error surfaces, because from the request handler’s point of view nothing failed: it already returned 200.',
      'This project shipped exactly that bug. Every email left the API from a bare goroutine started by the request, while the job client that handles all of it properly was never called from anywhere.',
      'A job queue replaces the goroutine with a durable record. The work is written down, a worker picks it up, failures retry with backoff, and the thing that is holding up delivery is visible rather than lost.',
    ],
    capabilities: [
      'Enqueue work from a handler and return immediately.',
      'Survive a process restart: queued work is not in memory.',
      'Retry with exponential backoff, a bounded number of times.',
      'Time out a stuck task so it does not hold a worker slot forever.',
      'Refuse a duplicate enqueue of the same logical task within a window.',
      'Separate queues by urgency, so a backlog of thumbnails does not delay a password reset.',
      'Show what is queued, running, retrying and dead, with the error.',
      'Run with no queue at all, in a deployment that has no Redis.',
    ],
  },

  functional: [
    'An enqueue call per job type, typed rather than a map of strings.',
    'Workers registered per task type, started as their own process or alongside the API.',
    'Five retry attempts by default, with exponential backoff.',
    'A five minute task timeout by default, overridable per enqueue.',
    'An idempotency key with a 24 hour window, deduplicating repeat enqueues.',
    'Delayed enqueue, for work that should happen later rather than now.',
    'Three queues: critical, default and low.',
    'Completed tasks retained for 24 hours so a success can be inspected, not just a failure.',
    'An admin dashboard over the queues, with retry and delete.',
  ],

  nonFunctional: [
    {
      label: 'Durable',
      text: 'enqueued work is in Redis, not in a goroutine. A deploy mid-flight loses nothing that was accepted.',
    },
    {
      label: 'Bounded',
      text: 'worker concurrency is configured, so load is a queue depth rather than a connection exhaustion. This is the single biggest difference from a goroutine.',
    },
    {
      label: 'At least once',
      text: 'a task can run twice if a worker dies after doing the work and before acknowledging. Handlers must be idempotent, and the framework says so rather than implying exactly-once.',
    },
    {
      label: 'Observable',
      text: 'a failed job and its error are in the dashboard. The failure mode of the goroutine version was that nobody ever found out.',
    },
    {
      label: 'Optional',
      text: 'with no Redis configured, jobs, cron and the cache are all disabled and the application runs. The alternative is requiring Redis for a project that wanted SQLite and a single binary.',
    },
  ],

  capacity: {
    assumptions: [
      ['Jobs enqueued per second', '~45 average, 280 peak'],
      ['Average job duration', '400 ms'],
      ['Slowest common job', 'image processing, 3 to 8 seconds'],
      ['Worker concurrency per process', '10'],
      ['Payload size', '~1 KB'],
      ['Retry attempts', '5, exponential'],
    ],
    estimates: [
      {
        label: 'Worker capacity',
        working: [
          '10 concurrent slots / 0.4 s average = 25 jobs/s per worker process',
          '280 jobs/s at peak / 25 = 12 worker processes',
          'or fewer processes with higher concurrency per process',
        ],
        note: 'Concurrency is the number to turn, not the process count, until the database connection pool becomes the limit. Each concurrent job may hold a connection.',
      },
      {
        label: 'Queue memory',
        working: [
          'a backlog of 100,000 jobs x 1 KB = 100 MB',
          'plus completed retention: 45 jobs/s x 86,400 s = ~3.9M/day',
          'at 24 hours retention and 1 KB that is ~3.9 GB',
        ],
        note: 'Retention is the number that surprises people. It is what makes a success inspectable, and it is also what fills Redis. Lower it per task type for anything high volume.',
      },
      {
        label: 'Backoff schedule',
        working: [
          'attempt 1 immediate, then roughly 1 s, 4 s, 15 s, 60 s',
          'five attempts spans about 80 seconds of downstream outage',
          'longer outages need a higher MaxRetry on that task type',
        ],
        note: 'Five was chosen to ride out a short blip without filling the dead queue with hours-old work that is no longer worth doing.',
      },
      {
        label: 'Database connections',
        working: [
          '12 workers x 10 concurrency = 120 potential connections',
          'plus the API processes',
          'against a pool that is usually configured at 25 per process',
        ],
        note: 'This is the real constraint on scaling workers, and it is the thing that breaks when somebody raises concurrency to fix a backlog.',
      },
    ],
  },

  highLevel: {
    intro:
      'A client that writes the task, a broker that holds it, workers that run it, and a dashboard over the broker.',
    components: [
      {
        label: 'Job client',
        text: 'what a handler calls. One typed method per job type, so a payload cannot be wrong at the call site and discovered in the worker.',
      },
      {
        label: 'Broker',
        text: 'Redis, through asynq. Holds the pending, scheduled, retrying, completed and dead sets.',
      },
      {
        label: 'Workers',
        text: 'a process with a handler registered per task type and a concurrency setting. Runs as its own container in production and alongside the dev server locally.',
      },
      {
        label: 'Retry policy',
        text: 'attempts, backoff and timeout, with defaults that each task type can override. Set at enqueue rather than in the worker, so the caller decides how important the work is.',
      },
      {
        label: 'Idempotency',
        text: 'a key plus a window. The same logical task enqueued twice inside the window is refused, and the refusal is a success for the caller because the original is already on its way.',
      },
      {
        label: 'Dashboard',
        text: 'the admin view over the queues: depth, failures with their errors, retry and delete.',
      },
    ],
    flow: {
      title: 'A handler enqueuing work, and a worker running it',
      nodes: [
        { id: 'client', label: 'Client', col: 0, row: 0 },
        { id: 'handler', label: 'Handler', col: 1, row: 0 },
        { id: 'enqueue', label: 'Enqueue', col: 2, row: 0, accent: true },
        { id: 'redis', label: 'Redis Queue', col: 3, row: 0 },
        { id: 'dedupe', label: 'Idempotency Key', col: 2, row: 1, accent: true },
        { id: 'worker', label: 'Worker Pool', col: 3, row: 1, accent: true },
        { id: 'task', label: 'Task Handler', col: 3, row: 2 },
        { id: 'retry', label: 'Retry or Dead', col: 2, row: 2, accent: true },
        { id: 'admin', label: 'Jobs Dashboard', col: 0, row: 2 },
      ],
      edges: [
        { from: 'client', to: 'handler', step: 1 },
        { from: 'handler', to: 'enqueue', step: 2 },
        { from: 'enqueue', to: 'dedupe', step: 3 },
        { from: 'enqueue', to: 'redis', step: 4 },
        { from: 'redis', to: 'worker', step: 5 },
        { from: 'worker', to: 'task', step: 6 },
        { from: 'task', to: 'retry', step: 7 },
        { from: 'retry', to: 'admin', step: 8, bend: 'h' },
      ],
      steps: [
        'A request arrives that would otherwise do slow work inline: a sign-up that sends a welcome email, an upload that needs thumbnails.',
        'The handler calls the typed enqueue method and returns. The response does not wait for the work, and the user does not wait for a mail provider.',
        'If the caller supplied an idempotency key, a duplicate inside the window is refused. The caller treats that as success, because the original task is already queued and running it twice is the thing being avoided.',
        'The task is written to Redis with its retry count, its timeout, its queue and its retention. It is now durable: a deploy, a crash or a restart does not lose it.',
        'A worker process pulls from the queues in priority order. Critical before default before low, so a backlog of thumbnails never delays a password reset.',
        'The registered handler for that task type runs, inside the timeout. It must be safe to run twice, because at-least-once delivery is what a durable queue buys and exactly-once is not available.',
        'On failure the task goes back with exponential backoff, up to five attempts by default. After the last attempt it moves to the dead set rather than disappearing.',
        'The dashboard shows depth, failures and the error text, with retry and delete. The whole point of the design over a goroutine is that a failure is a row somebody can see rather than a log line nobody wrote.',
      ],
    },
    dataFlow: [
      'A payload is small and references data rather than carrying it. An image job carries an upload identifier, not the image.',
      'Nothing in a payload can be assumed still to exist. By the time a task runs the row may have been deleted, so a handler loads and tolerates absence.',
      'Retry decisions are made at enqueue, not in the worker, because the caller is the one who knows whether this work is worth five attempts.',
      'Completed tasks are retained so a success can be inspected. That retention is also the main consumer of Redis memory at volume.',
    ],
  },

  stack: [
    ['Queue', 'asynq over Redis'],
    ['Queues', 'critical, default, low'],
    ['Defaults', '5 retries, 5 minute timeout, 24 hour idempotency window'],
    ['Retention', 'completed tasks kept 24 hours'],
    ['Workers', 'their own process in production, with the dev server locally'],
    ['Dashboard', 'admin pages over the broker'],
    ['Absent Redis', 'jobs, cron and cache all disabled, application runs'],
  ],

  api: {
    groups: [
      {
        title: 'Queue administration',
        rows: [
          { method: 'GET', path: '/api/v1/admin/jobs', what: 'Queue depths and worker state' },
          { method: 'GET', path: '/api/v1/admin/jobs/failed', what: 'Dead tasks with their last error' },
          { method: 'POST', path: '/api/v1/admin/jobs/:id/retry', what: 'Requeue one dead task' },
          { method: 'DELETE', path: '/api/v1/admin/jobs/:id', what: 'Discard a dead task' },
        ],
      },
    ],
    samples: [
      {
        title: 'Enqueuing from a handler, with the duplicate case handled',
        language: 'go',
        code: `err := h.Jobs.EnqueueSendEmail(ctx, jobs.SendEmailPayload{
    To:       user.Email,
    Template: "welcome",
    Data:     map[string]any{"name": user.Name},
}, jobs.EnqueueOption{
    Queue:          "critical",
    IdempotencyKey: "welcome:" + user.ID,
})
// A duplicate inside the window is not a failure: the original is on its way.
if err != nil && !errors.Is(err, jobs.ErrDuplicateTask) {
    return fmt.Errorf("queueing welcome mail: %w", err)
}`,
      },
    ],
  },

  lowLevel: {
    classes: [
      {
        name: 'jobs.Client',
        file: 'internal/jobs/jobs.go',
        what: 'Wraps the broker client with one typed method per job type. Typed on purpose: a string task name and a map payload moves every mistake from compile time to the worker.',
        methods: ['EnqueueSendEmail', 'EnqueueProcessImage', 'EnqueueCleanup'],
      },
      {
        name: 'jobs.EnqueueOption',
        what: 'Per-task overrides: queue, retries, timeout, delay, retention, idempotency key and window. Zero values fall back to the package defaults.',
      },
      {
        name: 'jobs.ErrDuplicateTask',
        what: 'Returned when the same idempotency key is enqueued inside its window. Named so callers can distinguish it, because treating it as an error would make a retry look like a failure.',
      },
      {
        name: 'mail_dispatch.go',
        file: 'internal/handlers/mail_dispatch.go',
        what: 'The single place mail leaves a handler. Added after the review found every email going out from a bare goroutine while the job client sat unused.',
      },
    ],
    principles: [
      {
        label: 'One way out of a handler',
        text: 'all mail goes through one dispatch function. Twelve handlers each starting a goroutine is twelve chances to get bounding, retry and shutdown wrong.',
      },
      {
        label: 'Typed payloads',
        text: 'a wrong field is a compile error rather than a nil map access inside a worker at three in the morning.',
      },
      {
        label: 'Say at-least-once out loud',
        text: 'documenting the real guarantee is what makes handlers idempotent. Implying exactly-once produces handlers that charge a card twice.',
      },
      {
        label: 'Defaults tuned for safety',
        text: 'five retries and a five minute timeout are about surviving a blip without hoarding stale work, not about maximum throughput.',
      },
    ],
    patterns: [
      ['Work queue', 'a durable record instead of an in-process goroutine'],
      ['Exponential backoff', 'retries that do not hammer a struggling dependency'],
      ['Dead letter', 'exhausted tasks kept rather than dropped'],
      ['Idempotency key', 'deduplicating the enqueue, not just the execution'],
      ['Priority queues', 'urgency separated so a backlog is not a blockage'],
    ],
  },

  scaling: [
    'Workers scale horizontally: more processes consume the same queues with no coordination.',
    'The first number to turn is concurrency per process, and the first thing it breaks is the database connection pool. Worker concurrency times worker count is a connection count.',
    'Priority queues mean a backlog is localised. Thumbnails piling up does not delay a password reset, which is the difference between a slow feature and a broken one.',
    'Completed retention dominates Redis memory at volume. It should be lowered per task type for anything high frequency.',
    'A long running job should be split. One task that takes an hour is a task that cannot be retried cheaply and that blocks a worker slot for an hour.',
    'With no Redis the application degrades to doing nothing in the background, which is correct for a small single-binary deployment and documented rather than surprising.',
  ],

  bottlenecks: {
    problems: [
      {
        label: 'Connection pool exhaustion',
        text: 'raising worker concurrency to clear a backlog takes connections from the API, and the symptom shows up in request latency rather than in the queue.',
      },
      {
        label: 'Non-idempotent handlers',
        text: 'at-least-once plus a handler that charges a card is a double charge, and it happens on the day the worker is killed mid-task.',
      },
      {
        label: 'Payloads carrying data',
        text: 'an image in a payload is an image in Redis, several times over with retries, and a queue that falls over on memory rather than depth.',
      },
      {
        label: 'Poison tasks',
        text: 'a task that fails deterministically burns five attempts and a backoff schedule every time, and a flood of them crowds out real work.',
      },
      {
        label: 'Redis as a single point',
        text: 'the broker is the durable store. Losing it loses the queue, and no amount of worker redundancy helps.',
      },
    ],
    improvements: [
      {
        label: 'Size the pool for workers plus API',
        text: 'count the real maximum: worker processes times concurrency, plus API processes times their pool. Then set the database limit above it.',
      },
      {
        label: 'Make handlers idempotent with a key',
        text: 'a processed-tasks table or a unique constraint on the effect, so the second run is a no-op rather than a second charge.',
      },
      {
        label: 'Pass references, not payloads',
        text: 'an identifier and let the handler load. The row is the source of truth and it may have changed since the enqueue anyway.',
      },
      {
        label: 'Separate queue for the known-flaky',
        text: 'a task type that fails often belongs on its own queue, where its backoff cannot crowd out anything that matters.',
      },
      {
        label: 'Persist and replicate Redis',
        text: 'append-only persistence and a replica, because the queue is durable storage whether or not it was planned as such.',
      },
    ],
  },

  seeAlso: [
    { title: 'Background jobs reference', href: '/docs/batteries/jobs' },
    { title: 'Scheduled tasks', href: '/docs/systems/scheduled-tasks' },
    { title: 'Transactional outbox', href: '/docs/systems/transactional-outbox' },
  ],
}

export const SCHEDULING: SystemDesign = {
  slug: 'scheduled-tasks',
  name: 'Scheduled Tasks',
  tagline:
    'Running something every night, exactly once, across however many replicas happen to be up.',
  group: 'Delivery',
  packages: ['internal/cron', 'internal/jobs'],

  problem: {
    text: [
      'Every application accumulates recurring work. Expire old sessions at midnight. Send the digest on Monday morning. Prune soft-deleted rows after thirty days. Reconcile with the payment provider hourly.',
      'A system crontab does it, outside the application, in a different language, with no access to the models, no logging the application can see, and no record on any machine that is not that machine. It also stops existing the moment deployment becomes containers.',
      'A ticker inside the process is worse in a specific way that only appears in production. It works perfectly on one instance. Scale to three replicas and the nightly job runs three times: three digests to each subscriber, three reconciliation passes, three prunes racing each other. Nothing errors, because each replica did exactly what it was told.',
      'The requirement is not "run on a schedule". It is "run on a schedule, once, regardless of how many copies of the application are running", and that needs something that can be held by one replica at a time.',
    ],
    capabilities: [
      'Declare a schedule in the application, next to the work it runs.',
      'Run each occurrence exactly once across all replicas.',
      'Survive the replica holding the schedule dying, by having another take over.',
      'Enqueue a job rather than doing the work inline, so a slow task does not delay the next tick.',
      'Show the schedule, the last run and the next, in the admin.',
      'Allow a manual run, for the case where last night failed.',
      'Run on standard cron expressions, because that is the notation everybody already knows.',
    ],
  },

  functional: [
    'Schedules declared in Go with a cron expression and a task type.',
    'A scheduler that enqueues rather than executes, so the work goes through the job queue.',
    'A distributed lock, so exactly one replica schedules.',
    'Automatic takeover when the lock holder goes away.',
    'An admin page listing each task, its expression, its last run and its next.',
    'A trigger button, enqueuing one occurrence immediately.',
    'The scheduler disabled with no Redis, logged plainly rather than failing.',
  ],

  nonFunctional: [
    {
      label: 'Exactly once per occurrence',
      text: 'this is the whole requirement. A digest sent three times is worse than a digest not sent, because the second is noticed and fixed.',
    },
    {
      label: 'Survives a replica loss',
      text: 'the lock is leased, not held. A replica that dies releases it by not renewing, and another picks it up without an operator.',
    },
    {
      label: 'Scheduling is not executing',
      text: 'the scheduler enqueues. Running the work inline means a task that takes ten minutes delays the next tick and a crash loses the occurrence.',
    },
    {
      label: 'Visible',
      text: 'the last run and the next are in the admin, because the failure mode of scheduled work is that it quietly stops and nobody notices for a fortnight.',
    },
    {
      label: 'Honest about time zones',
      text: 'the expression is evaluated in a declared zone. "Midnight" without one is midnight in whatever zone the container happened to get.',
    },
  ],

  capacity: {
    assumptions: [
      ['Scheduled tasks per project', '5 to 20'],
      ['Finest useful granularity', 'one minute'],
      ['API replicas', '2 to 12'],
      ['Lock lease', 'tens of seconds, renewed'],
      ['Heaviest task', 'nightly cleanup over millions of rows'],
    ],
    estimates: [
      {
        label: 'Scheduler cost',
        working: [
          '20 tasks, evaluated once a minute',
          '= 1,200 expression evaluations an hour',
          'each one an arithmetic comparison',
        ],
        note: 'Negligible. The scheduler is cheap precisely because it only enqueues. All the real cost is in the workers.',
      },
      {
        label: 'Lock traffic',
        working: [
          '12 replicas, each attempting or renewing every few seconds',
          'a handful of Redis operations per second, total',
        ],
      },
      {
        label: 'The midnight spike',
        working: [
          '8 tasks all written as 0 0 * * *',
          'all enqueue in the same second',
          'the worker pool absorbs them as a short backlog, not a stall',
        ],
        note: 'This is the argument for enqueuing rather than executing: eight simultaneous occurrences are a queue depth of eight, not eight things running in the scheduler at once.',
      },
      {
        label: 'Occurrences missed by downtime',
        working: [
          'no replica up at midnight = one occurrence missed',
          'the schedule does not backfill: the next run is the next occurrence',
        ],
        note: 'Deliberate. A nightly task catching up on four missed nights at once is usually worse than skipping them, and the manual trigger exists for when it is not.',
      },
    ],
  },

  highLevel: {
    intro:
      'Every replica runs a scheduler. Exactly one of them holds the lock, and only the holder enqueues.',
    components: [
      {
        label: 'Schedule registry',
        text: 'the declared tasks: an expression, a task type and a payload. In Go, next to the work, so adding one is a code change that gets reviewed.',
      },
      {
        label: 'Scheduler',
        text: 'evaluates the expressions and enqueues what is due. Runs in every replica, acts in one.',
      },
      {
        label: 'Distributed lock',
        text: 'a leased key in Redis. Whoever holds it schedules; whoever does not, waits. The lease is what makes a dead holder recoverable without intervention.',
      },
      {
        label: 'Job queue',
        text: 'where the occurrence goes. The scheduler’s output is a queued job, which inherits retries, timeouts and the dashboard from the queue.',
      },
      {
        label: 'Admin page',
        text: 'the tasks, their expressions, their last and next runs, and a trigger.',
      },
    ],
    flow: {
      title: 'Three replicas, one occurrence',
      nodes: [
        { id: 'r1', label: 'Replica 1', col: 0, row: 0 },
        { id: 'lock', label: 'Cron Lock', col: 1, row: 0, accent: true },
        { id: 'sched', label: 'Scheduler', col: 2, row: 0, accent: true },
        { id: 'queue', label: 'Job Queue', col: 3, row: 0 },
        { id: 'r2', label: 'Replica 2', col: 0, row: 1 },
        { id: 'r3', label: 'Replica 3', col: 0, row: 2 },
        { id: 'worker', label: 'Worker', col: 3, row: 1 },
        { id: 'admin', label: 'Cron Admin', col: 2, row: 2 },
        { id: 'run', label: 'Task Runs Once', col: 3, row: 2 },
      ],
      edges: [
        { from: 'r1', to: 'lock', step: 1 },
        { from: 'r2', to: 'lock', step: 2, bend: 'h', dashed: true },
        { from: 'r3', to: 'lock', step: 3, bend: 'h', dashed: true },
        { from: 'lock', to: 'sched', step: 4 },
        { from: 'sched', to: 'queue', step: 5 },
        { from: 'queue', to: 'worker', step: 6 },
        { from: 'worker', to: 'run', step: 7 },
        { from: 'admin', to: 'queue', step: 8, bend: 'v' },
      ],
      steps: [
        'Every replica starts a scheduler and tries to take the cron lock. One of them gets it.',
        'Replica 2 does not. It keeps trying on an interval, so that if the holder disappears it is already a candidate rather than needing a deploy.',
        'Replica 3 likewise. The scheduler is running in all three and acting in one, which is what makes takeover automatic rather than operational.',
        'The holder renews its lease while it keeps running. A lease rather than a permanent key is the whole mechanism: a replica that is killed stops renewing, the key expires, and the next attempt succeeds.',
        'When an expression comes due the holder enqueues a job. It does not run the work, so a ten minute cleanup does not delay the next evaluation and a crash during it loses nothing.',
        'A worker picks the job up like any other, with the same retries, the same timeout and the same dashboard.',
        'The task runs once, because only one replica enqueued it.',
        'A manual trigger from the admin enqueues the same job directly, for the night that failed or the digest somebody needs now.',
      ],
    },
    dataFlow: [
      'The lock is a lease with an expiry, renewed by the holder. That is what distinguishes "another replica takes over in thirty seconds" from "a human notices tomorrow".',
      'Schedules do not backfill. A missed occurrence stays missed, and the manual trigger is the deliberate way to make one up.',
      'The scheduler enqueues, which means every property the job queue has applies to scheduled work for free: retries, timeout, dead letter, the dashboard.',
      'Expressions are evaluated in a declared zone, so a nightly task does not drift by an hour twice a year.',
    ],
  },

  stack: [
    ['Scheduler', 'asynq periodic tasks'],
    ['Coordination', 'a leased lock in Redis'],
    ['Notation', 'standard cron expressions'],
    ['Granularity', 'one minute'],
    ['Execution', 'enqueued to the job queue, never run inline'],
    ['Admin', 'a page listing tasks, last run, next run, and a trigger'],
    ['Absent Redis', 'scheduler disabled, logged'],
  ],

  api: {
    groups: [
      {
        title: 'Schedule administration',
        rows: [
          { method: 'GET', path: '/api/v1/admin/cron', what: 'Declared tasks with last and next run' },
          { method: 'POST', path: '/api/v1/admin/cron/:name/run', what: 'Enqueue one occurrence now' },
        ],
      },
    ],
  },

  lowLevel: {
    classes: [
      {
        name: 'cron.Register',
        file: 'internal/cron/cron.go',
        what: 'Declares a schedule: an expression, a task type and a payload. In code next to the work rather than in a configuration file, so adding one goes through review.',
      },
      {
        name: 'cron.Start',
        what: 'Starts the scheduler in this replica. Competes for the lock and acts only while holding it. With no Redis it logs that scheduling is off and returns.',
      },
      {
        name: 'CronHandler',
        what: 'The admin endpoints: list the tasks, trigger one. The trigger enqueues the same job the scheduler would, so a manual run is not a different code path.',
      },
    ],
    principles: [
      {
        label: 'Schedule, do not execute',
        text: 'the scheduler enqueues and nothing more. Everything else follows from that: a slow task cannot delay a tick, and a scheduler crash cannot lose work that was already queued.',
      },
      {
        label: 'Assume more than one replica',
        text: 'a design that is correct on one instance and wrong on three is a design that is wrong, because the third instance arrives without a decision being made.',
      },
      {
        label: 'Lease, do not hold',
        text: 'a permanent lock needs a human when its holder dies. An expiring lease needs nobody.',
      },
      {
        label: 'Do not catch up',
        text: 'missed occurrences stay missed, and a manual trigger makes one up deliberately. Automatic backfill means a deploy outage sends four days of digests at once.',
      },
    ],
    patterns: [
      ['Leader election', 'by leased lock, so exactly one replica schedules'],
      ['Scheduler and executor split', 'the tick enqueues, the worker runs'],
      ['Declarative schedule', 'expressions in code, next to the work'],
      ['Manual override', 'the same enqueue path, triggered by a person'],
    ],
  },

  scaling: [
    'Adding replicas does not add occurrences, which is the single property the design exists to provide.',
    'The scheduler itself never becomes a bottleneck, because its work is an arithmetic comparison per task per minute.',
    'All the real cost is in the workers, which scale independently of the schedule.',
    'Tasks written as midnight all fire at midnight. Spreading the expressions turns one spike into a flat few minutes, and costs nothing.',
    'A heavy nightly task over millions of rows should process in batches and be resumable, because it will eventually not finish inside its window.',
    'Takeover time is the lease duration. Shorter means faster recovery and more lock traffic, and tens of seconds is the right end of that trade for nightly work.',
  ],

  bottlenecks: {
    problems: [
      {
        label: 'Clock skew between replicas',
        text: 'a few seconds of drift can mean an occurrence evaluated twice at a boundary, with the lock the only thing preventing a double enqueue.',
      },
      {
        label: 'A task that outgrows its window',
        text: 'a nightly cleanup that takes nine hours overlaps the next night, and two copies of it then compete over the same rows.',
      },
      {
        label: 'Silent stoppage',
        text: 'the characteristic failure of scheduled work is that it stops and nobody finds out. No error is raised by something not happening.',
      },
      {
        label: 'Everything at midnight',
        text: 'eight tasks on the same expression make a spike that can push a worker pool into a backlog that lasts into the morning.',
      },
      {
        label: 'Time zone and daylight saving',
        text: 'an hourly task runs twice or not at all on the two days a year the clock changes, in whichever zone the expression is evaluated.',
      },
    ],
    improvements: [
      {
        label: 'Make tasks idempotent by occurrence',
        text: 'key the work on the date it is for, so a second enqueue for the same night is a no-op. Cheap, and it makes the lock a performance feature rather than a correctness one.',
      },
      {
        label: 'Batch and checkpoint long tasks',
        text: 'process in batches, record progress, and let the next occurrence resume. Then overrunning is slow rather than broken.',
      },
      {
        label: 'Alert on a missing run',
        text: 'a heartbeat per task and an alert when the last run is older than the expression allows. This is the only way the silent failure becomes visible.',
      },
      {
        label: 'Spread the expressions',
        text: 'move the eight midnight tasks to eight different minutes. It is a one-line change and it removes the spike entirely.',
      },
      {
        label: 'Schedule in UTC',
        text: 'and convert for display. The twice-yearly ambiguity disappears, at the cost of a nightly task happening an hour later for half the year.',
      },
    ],
  },

  seeAlso: [
    { title: 'Cron reference', href: '/docs/batteries/cron' },
    { title: 'Background jobs', href: '/docs/systems/background-jobs' },
  ],
}

export const OUTBOX: SystemDesign = {
  slug: 'transactional-outbox',
  name: 'Transactional Outbox',
  tagline:
    'Saving a row and telling the world about it as one atomic act, when the database and the message broker cannot share a transaction.',
  group: 'Delivery',
  packages: ['internal/outbox', 'internal/models'],

  problem: {
    text: [
      'An order is placed. A row is written and a webhook goes out. Those are two systems, and there is no transaction that spans both, so one of them happens first and the other can fail.',
      'Publish first, then commit: the webhook succeeds, the commit fails, and a downstream system has been told about an order that does not exist. It acts on it. There is nothing to reconcile against, because the row was never written.',
      'Commit first, then publish: the commit succeeds, the process is killed before the publish, and the order exists with nobody told. Nothing is logged anywhere, because from the process’s point of view nothing failed.',
      'Both orderings are wrong, and which one is wrong in a worse way depends on the integration. The outbox pattern removes the choice: the message is written to a table in the same transaction as the business data, so it commits or rolls back with it, and a separate relay delivers committed messages afterwards. Either both happened or neither did.',
    ],
    capabilities: [
      'Write a message in the caller’s transaction, so it shares the fate of the data.',
      'Refuse to enqueue outside a transaction, because that is the bug with extra steps.',
      'Deliver committed messages from a relay, separately from the request.',
      'Claim a message before delivering, so two relays do not deliver it twice.',
      'Retry a failed delivery with backoff, and stop after enough attempts.',
      'Keep delivered messages as a record of what was sent and when.',
      'Be honest that delivery is at-least-once, and say what consumers must do about it.',
    ],
  },

  functional: [
    'An enqueue that takes the transaction and returns an error if it is not one.',
    'A message row with a type, a payload, a status and an attempt count.',
    'Four statuses: pending, claimed, delivered, failed.',
    'A relay that claims pending messages, delivers them and records the outcome.',
    'Backoff between attempts, and a cap after which a message is marked failed.',
    'The table declared as a model, so migrations create it and backups include it.',
    'Delivered messages retained, as the audit trail of what left the system.',
  ],

  nonFunctional: [
    {
      label: 'Atomic with the data',
      text: 'the message is in the same transaction as the row. This is the entire property, and everything else is in service of it.',
    },
    {
      label: 'At least once, stated plainly',
      text: 'a message can be delivered twice if the process dies between the send and the status update. Consumers must be idempotent, and a dedup key on the receiving end is what that means in practice.',
    },
    {
      label: 'Ordered enough',
      text: 'messages are relayed oldest first. Strict global ordering is not offered, because it would mean a single relay and no parallelism.',
    },
    {
      label: 'Durable',
      text: 'the outbox is a database table, so it inherits the database’s durability, replication and backups rather than needing its own.',
    },
    {
      label: 'Refuses the wrong usage',
      text: 'enqueuing without a transaction is rejected rather than accepted. It would work in every test and fail only in production, which is the worst possible failure mode to allow.',
    },
  ],

  capacity: {
    assumptions: [
      ['Messages per second', '~30 average, 200 peak'],
      ['Payload size', '~2 KB of JSON'],
      ['Relay poll interval', 'seconds'],
      ['Delivery target latency', 'under a second at normal load'],
      ['Retention of delivered rows', '30 days'],
    ],
    estimates: [
      {
        label: 'Table growth',
        working: [
          '30 messages/s x 86,400 = ~2.6M rows/day',
          'at ~2.2 KB per row = ~5.7 GB/day',
          'x 30 days retention = ~170 GB',
        ],
        note: 'This is the number that decides whether the outbox needs pruning, and at this volume it very much does. A smaller project at 1 message/s is 190 MB a month and needs nothing.',
      },
      {
        label: 'Relay throughput',
        working: [
          'batch of 100 claimed per poll',
          'delivery at ~50 ms each, 10 in parallel = ~0.5 s per batch',
          '= ~200 messages/s per relay process',
        ],
        note: 'One relay keeps up with peak. More than one requires the claim step to be correct, which is why claiming is a status transition and not a read.',
      },
      {
        label: 'Added write cost',
        working: [
          'one extra INSERT inside an existing transaction',
          'no extra round trip to a broker in the request path',
          'the transaction is marginally longer, not fundamentally slower',
        ],
      },
      {
        label: 'Delivery latency',
        working: [
          'commit to next poll: up to the poll interval',
          'plus delivery time',
          'so seconds, not milliseconds',
        ],
        note: 'The outbox trades latency for correctness. Anything that genuinely needs sub-second delivery needs a different mechanism, and usually does not actually need it.',
      },
    ],
  },

  highLevel: {
    intro:
      'One extra insert inside the transaction, and one process reading committed rows.',
    components: [
      {
        label: 'Outbox table',
        text: 'declared as a model like any other, so AutoMigrate creates it and the backup writer includes it. Calling code still refers to it through the outbox package.',
      },
      {
        label: 'Enqueue',
        text: 'writes the message using the caller’s transaction. Takes the transaction as an argument and errors if handed anything else.',
      },
      {
        label: 'Claim',
        text: 'moves a batch from pending to claimed in one statement. This is what makes more than one relay safe.',
      },
      {
        label: 'Relay',
        text: 'delivers claimed messages and records the result. Its own process, or a scheduled job, independent of the request that created the message.',
      },
      {
        label: 'Status machine',
        text: 'pending, claimed, delivered, failed. Strings rather than an enum so a person reading the table can see what happened without a lookup.',
      },
    ],
    flow: {
      title: 'A row and its message, committing together',
      nodes: [
        { id: 'handler', label: 'Handler', col: 0, row: 0 },
        { id: 'tx', label: 'Transaction', col: 1, row: 0, accent: true },
        { id: 'row', label: 'Order Row', col: 2, row: 0 },
        { id: 'msg', label: 'Outbox Row', col: 3, row: 0, accent: true },
        { id: 'commit', label: 'Commit', col: 1, row: 1, accent: true },
        { id: 'relay', label: 'Relay', col: 2, row: 1, accent: true },
        { id: 'claim', label: 'Claim Batch', col: 3, row: 1 },
        { id: 'target', label: 'Webhook Target', col: 3, row: 2 },
        { id: 'status', label: 'Delivered or Retry', col: 2, row: 2 },
      ],
      edges: [
        { from: 'handler', to: 'tx', step: 1 },
        { from: 'tx', to: 'row', step: 2 },
        { from: 'tx', to: 'msg', step: 3, bend: 'h' },
        { from: 'tx', to: 'commit', step: 4 },
        { from: 'commit', to: 'relay', step: 5 },
        { from: 'relay', to: 'claim', step: 6 },
        { from: 'claim', to: 'target', step: 7, bend: 'v' },
        { from: 'target', to: 'status', step: 8 },
      ],
      steps: [
        'A handler starts a transaction, because the outbox needs one and will say so if it does not get one.',
        'The business data is written: the order, its lines, whatever the operation is.',
        'The message is written in the same transaction. This is the whole design. From here, the row and the announcement of the row have exactly one fate between them.',
        'The transaction commits, or it does not. If it rolls back, the message rolls back with it, so nothing was announced that did not happen.',
        'Some time later, within seconds, the relay looks for pending messages. It is a separate process, so a crash in the request path after the commit changes nothing.',
        'A batch is claimed: pending to claimed in one statement. A second relay reading at the same moment sees nothing to claim, which is what makes running two of them safe.',
        'Each message is delivered. The attempt count rises so a persistently failing target backs off rather than being retried in a tight loop.',
        'The outcome is recorded. Delivered stays as the record of what was sent; a failure goes back for another attempt until the cap, after which it is marked failed and waits for a person. A process killed between the send and this update means the message is delivered again later, which is exactly why consumers need to be idempotent.',
      ],
    },
    dataFlow: [
      'The message payload is a snapshot of what was true at commit time. It does not re-read the row at delivery, because the row may have changed and the message is about what happened, not about what is.',
      'Claiming is a status transition rather than a read, so concurrent relays partition the work instead of duplicating it.',
      'Delivered rows are kept. They are the answer to "did we tell them, and when", which is the question asked during every integration dispute.',
      'At-least-once is a consequence of the crash window between sending and recording. It cannot be closed without a transaction spanning the database and the target, which is the thing that does not exist.',
    ],
  },

  stack: [
    ['Storage', 'a database table, in the application database'],
    ['Atomicity', 'the caller’s transaction'],
    ['Statuses', 'pending, claimed, delivered, failed'],
    ['Relay', 'its own process or a scheduled job'],
    ['Guarantee', 'at least once'],
    ['Ordering', 'oldest first, not strictly global'],
  ],

  dataModel: {
    entities: [
      {
        name: 'outbox_messages',
        fields: [
          ['id', 'primary key'],
          ['type', 'what happened, for example order.created'],
          ['payload', 'JSON, a snapshot at commit time'],
          ['status', 'pending, claimed, delivered or failed'],
          ['attempts', 'how many deliveries have been tried'],
          ['last_error', 'why the last attempt failed'],
          ['created_at', 'when the transaction committed it'],
          ['delivered_at', 'when it was accepted'],
        ],
        note: 'The index that matters is on status and created_at, because every relay poll is a query for the oldest pending rows and it runs several times a minute forever.',
      },
    ],
    storage: [
      'In the application database on purpose. A separate store would need its own transaction, which is the problem being solved.',
      'The table grows with every message and needs pruning at volume. Delivered rows older than the retention period are the ones to go.',
      'It is included in backups because it is declared as a model, which also means a restore brings undelivered messages back and they will be delivered.',
    ],
  },

  lowLevel: {
    classes: [
      {
        name: 'outbox.Enqueue',
        file: 'internal/outbox/outbox.go',
        what: 'Writes a message using the transaction it is handed. The signature takes the transaction rather than a database handle, so the correct usage is the only one that compiles cleanly.',
      },
      {
        name: 'outbox.ErrNoTransaction',
        what: 'Returned when enqueue is given something that is not a transaction. Refused rather than tolerated, because enqueuing outside a transaction is the commit-then-publish bug and it fails only in production.',
      },
      {
        name: 'outbox.Message',
        what: 'An alias for the model type. The table is declared in internal/models with every other table so migrations and backups cover it, while calling code still reads outbox.Message.',
      },
      {
        name: 'outbox.Relay',
        file: 'internal/outbox/relay.go',
        what: 'Claims a batch, delivers each message, records the outcome, backs off on failure.',
      },
    ],
    principles: [
      {
        label: 'Refuse the usage that fails silently',
        text: 'enqueue without a transaction is not a degraded version of the right thing, it is the original bug. Returning an error is better than accepting it.',
      },
      {
        label: 'One table, declared like any other',
        text: 'putting it in models means AutoMigrate, backups and the studio all cover it with no special cases.',
      },
      {
        label: 'Claim before delivering',
        text: 'a status transition rather than a read is the difference between two relays sharing the work and two relays doing all of it.',
      },
      {
        label: 'Say at-least-once',
        text: 'the crash window is real and cannot be closed. Documenting it is what makes consumers idempotent instead of surprised.',
      },
    ],
    patterns: [
      ['Transactional outbox', 'the message committed with the data'],
      ['Polling publisher', 'a relay reading committed rows'],
      ['Claim check', 'a status transition partitioning work between relays'],
      ['At-least-once delivery', 'with idempotent consumers as the stated requirement'],
    ],
  },

  scaling: [
    'The write path cost is one extra insert inside a transaction that already exists, which is as cheap as this problem gets.',
    'Relay throughput scales by running more relays, which is safe because claiming is a status transition.',
    'The table is the thing that needs managing. At thousands of messages a second it grows by gigabytes a day and needs a prune job.',
    'The index on status and created_at is not optional. Without it, every poll is a scan of a table that only grows.',
    'Delivery latency is bounded by the poll interval. Shortening it costs queries and buys seconds, and a listen-notify trigger removes the polling entirely when the latency genuinely matters.',
    'Because the outbox lives in the application database, it inherits replication and backup rather than adding a second durable system to operate.',
  ],

  bottlenecks: {
    problems: [
      {
        label: 'Table growth',
        text: 'retaining delivered messages is what makes the audit trail useful and what fills the disk. Both are true at once.',
      },
      {
        label: 'Polling latency',
        text: 'a message waits up to the poll interval before anybody looks at it, which is fine for a webhook and not for anything interactive.',
      },
      {
        label: 'A stuck claim',
        text: 'a relay that dies after claiming leaves rows in claimed that no other relay will touch, and they stay there until something notices.',
      },
      {
        label: 'Duplicate delivery',
        text: 'the crash window between sending and recording guarantees it will happen eventually, and a consumer that is not idempotent will act twice.',
      },
      {
        label: 'Head of line blocking',
        text: 'one target that is down and retrying can monopolise a relay batch, delaying messages for targets that are perfectly healthy.',
      },
    ],
    improvements: [
      {
        label: 'Prune delivered rows on a schedule',
        text: 'a nightly job deleting delivered messages past the retention window, in batches. The retention length is an integration decision, not a technical one.',
      },
      {
        label: 'Use listen and notify where latency matters',
        text: 'the commit signals the relay directly, so delivery is immediate, with the poll kept as the fallback that makes it correct.',
      },
      {
        label: 'Time out claims',
        text: 'treat a row claimed longer than a threshold as pending again. The cost is a possible duplicate delivery, which at-least-once already allows for.',
      },
      {
        label: 'Require a dedup key at the consumer',
        text: 'the message identifier in the delivery, and a uniqueness check on the receiving side. This is the only real answer to duplicates.',
      },
      {
        label: 'Partition by target',
        text: 'claim per destination so a failing target backs off on its own without holding up anybody else’s messages.',
      },
    ],
  },

  seeAlso: [
    { title: 'Outbox reference', href: '/docs/concepts/outbox' },
    { title: 'Webhooks', href: '/docs/systems/webhooks' },
    { title: 'Background jobs', href: '/docs/systems/background-jobs' },
  ],
}
