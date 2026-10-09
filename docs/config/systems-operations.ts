import type { SystemDesign } from './systems-types'

export const OBSERVABILITY: SystemDesign = {
  slug: 'observability',
  name: 'Observability',
  tagline:
    'Being able to answer what happened, where the time went, and whether the thing is actually working, from outside the process.',
  group: 'Operations',
  packages: ['internal/tracing', 'internal/health', 'internal/middleware'],

  problem: {
    text: [
      'A user reports that saving an order took eleven seconds at about ten past four. There is no way to answer that from the code. The question needs evidence collected before anybody knew it would be asked, which is the whole premise of observability: the instrumentation has to already be there, because the interesting failures are the ones nobody predicted.',
      'Three different questions get conflated into one word. What happened is logs, and a log line without a request identifier cannot be joined to the other forty lines from the same request. Where did the time go is tracing, and a trace crossing service boundaries is the one artifact that cannot be reconstructed from three sets of logs with three different request identifiers. And is it working is health, which sounds like the easy one and is where the subtlest mistake lives.',
      'That mistake is a boolean. Every component here used to report one, and that boolean covered two opposite situations: Redis is down, and this deployment has no Redis. An operator cannot tell those apart, and neither could the dashboard, so it guessed, differently per component. Redis read one way, email another, and storage was not reported at all. A monitoring system that cannot distinguish broken from absent trains the people reading it to ignore it.',
    ],
    capabilities: [
      'Give every request an identifier, in the response and in every log line it causes.',
      'Make that identifier the same string as the trace identifier, so one lookup serves both.',
      'Emit spans to a collector when one is configured, and cost nothing when it is not.',
      'Report component health in four states, not two.',
      'Let anything wired into the application register its own probe without editing the handler.',
      'Decide the overall status from the component states by one rule in one place.',
      'Record who did what, as a queryable log rather than as text.',
      'Expose counters for the things that fail silently: dropped realtime messages, queue depth.',
    ],
  },

  functional: [
    'A request identifier middleware, generating one or honouring an inbound header.',
    'The identifier in the X-Request-ID response header and in structured log lines.',
    'OpenTelemetry spans per request, exported when an OTLP endpoint is configured and inert otherwise.',
    'The request identifier and the trace identifier being the same value.',
    'A health registry: a name and a probe function, registered from anywhere.',
    'Four component states: ok, degraded, off, unknown.',
    'An overall status derived from the components, where only degraded lowers it.',
    'An activity log of user actions, with the actor, the action, the target and the time.',
    'Hub and queue counters in the health response.',
  ],

  nonFunctional: [
    {
      label: 'Four states, not a boolean',
      text: 'absent and broken are different facts and must be different values. Collapsing them means the dashboard has to guess, and it will guess differently in each place.',
    },
    {
      label: 'Unknown is not degraded',
      text: 'not knowing is not the same as being broken. Treating a probe that could not run as a failure is how an on-call rota learns to ignore a page.',
    },
    {
      label: 'One joining key',
      text: 'the request identifier and the trace identifier are the same string. Two identifiers for one request means every investigation starts with a translation step.',
    },
    {
      label: 'Inert when unconfigured',
      text: 'tracing with no endpoint costs nothing and logs nothing. Instrumentation that requires a collector to exist is instrumentation that gets removed.',
    },
    {
      label: 'Probes are cheap and bounded',
      text: 'health is polled by a load balancer every few seconds. A probe that runs an expensive query turns monitoring into load.',
    },
    {
      label: 'The rule lives in one place',
      text: 'how component states become an overall status is written once, so a new component cannot quietly invent its own interpretation.',
    },
  ],

  capacity: {
    assumptions: [
      ['Requests per second', '926 average, 4,630 peak'],
      ['Log lines per request', '3 to 8'],
      ['Span bytes per request', '~1.5 KB at full sampling'],
      ['Health poll interval', '5 seconds, from each load balancer'],
      ['Activity log rows per day', '~200,000'],
      ['Log retention', '14 days'],
    ],
    estimates: [
      {
        label: 'Log volume',
        working: [
          '926 req/s x 5 lines x ~400 bytes = ~1.85 MB/s',
          'x 86,400 = ~160 GB/day',
          'x 14 days = ~2.2 TB retained',
        ],
        note: 'Log volume is the cost that surprises people, and it is why levels and sampling exist. Debug logging left on in production is this number times several.',
      },
      {
        label: 'Trace volume at full sampling',
        working: [
          '926 req/s x 1.5 KB = ~1.4 MB/s',
          'x 86,400 = ~120 GB/day',
        ],
        note: 'Which is why sampling exists. At 10% head sampling that is 12 GB/day, and tail sampling keeps the slow and failed requests that are the ones worth having.',
      },
      {
        label: 'Health probe load',
        working: [
          '4 load balancer targets x every 5 s = 0.8 polls/s',
          'each running 5 probes = 4 probe executions/s',
        ],
        note: 'Negligible, as long as each probe is a ping rather than a count. A probe doing SELECT COUNT(*) turns this into real database load forever.',
      },
      {
        label: 'Activity log growth',
        working: [
          '200,000 rows/day x ~500 bytes = ~100 MB/day',
          'x 365 = ~36 GB/year',
        ],
        note: 'It is in the application database because it is queried alongside application data. That makes it the fastest growing table in most projects.',
      },
    ],
  },

  highLevel: {
    intro:
      'Three mechanisms with one joining key: an identifier per request, spans around the work, and a registry of probes.',
    components: [
      {
        label: 'Request identifier',
        text: 'middleware generating one or accepting an inbound header. In the response and in every log line, so forty lines become one request.',
      },
      {
        label: 'Tracing',
        text: 'a span per request, with child spans around database calls and outbound requests. Exported to an OTLP collector when one is configured, inert otherwise.',
      },
      {
        label: 'Structured logs',
        text: 'key and value rather than a sentence, so a log store can filter by status, route and identifier instead of matching text.',
      },
      {
        label: 'Health registry',
        text: 'a name and a function. Snapshot runs them all. A plugin or anything wired in main registers itself without the handler being edited.',
      },
      {
        label: 'Four-state status',
        text: 'ok, degraded, off, unknown, with the derivation rule in one place: only degraded lowers the overall status.',
      },
      {
        label: 'Activity log',
        text: 'who did what, as rows. A hash chain over them, so the record can be shown not to have been edited.',
      },
      {
        label: 'Counters',
        text: 'what the realtime hub delivered, dropped and failed to publish, and what the queues hold, in the health response. These are the things that fail without raising an error.',
      },
    ],
    flow: {
      title: 'One request, instrumented three ways',
      nodes: [
        { id: 'client', label: 'Client', col: 0, row: 0 },
        { id: 'reqid', label: 'Request ID', col: 1, row: 0, accent: true },
        { id: 'span', label: 'Root Span', col: 2, row: 0, accent: true },
        { id: 'handler', label: 'Handler', col: 3, row: 0 },
        { id: 'logs', label: 'Structured Logs', col: 1, row: 1, accent: true },
        { id: 'child', label: 'Child Spans', col: 2, row: 1 },
        { id: 'activity', label: 'Activity Log', col: 3, row: 1 },
        { id: 'otlp', label: 'OTLP Collector', col: 2, row: 2 },
        { id: 'probe', label: 'Health Snapshot', col: 0, row: 2 },
      ],
      edges: [
        { from: 'client', to: 'reqid', step: 1 },
        { from: 'reqid', to: 'span', step: 2 },
        { from: 'span', to: 'handler', step: 3 },
        { from: 'reqid', to: 'logs', step: 4 },
        { from: 'handler', to: 'child', step: 5, bend: 'h' },
        { from: 'handler', to: 'activity', step: 6 },
        { from: 'child', to: 'otlp', step: 7 },
        { from: 'probe', to: 'logs', step: 8, bend: 'h' },
      ],
      steps: [
        'A request arrives. If it carries a request identifier from an upstream service, that one is kept, so a trace does not start over at every hop.',
        'A root span opens, and its trace identifier is the request identifier. One string joins the log store and the trace store, which removes the translation step from every investigation.',
        'The handler runs. It does not know it is being traced or logged, because instrumentation that each handler participates in is instrumentation that half of them forget.',
        'Every log line carries the identifier, the route, the status and the duration as fields rather than as prose, so the log store can filter rather than grep.',
        'Child spans wrap the database calls and the outbound requests. This is what answers where the eleven seconds went, and it answers it without anybody having predicted the question.',
        'Actions worth attributing are written to the activity log: who did what, to what, when. A hash chain over the rows means the record can be shown not to have been altered afterwards.',
        'Spans go to the collector when an OTLP endpoint is configured. With none, the whole path is inert: no buffering, no errors, no cost. Instrumentation that demands infrastructure gets deleted.',
        'Separately, the health endpoint runs every registered probe and reports four states per component. Only degraded lowers the overall status, because a component that is off was never asked for and one that is unknown has not accused anybody of anything.',
      ],
    },
    dataFlow: [
      'The request identifier propagates outbound as a header, so a downstream service continues the same trace rather than starting one.',
      'Health probes are cheap by rule. A ping, a connection check, a cached count. A probe doing real work becomes load, forever, every five seconds.',
      'The activity log is in the application database because it is queried with application data. That is also why it grows faster than anything else.',
      'Counters are in the health response rather than only in logs, because a realtime hub dropping messages or a queue backing up produces no error anywhere.',
    ],
  },

  stack: [
    ['Tracing', 'OpenTelemetry, exported on OTEL_EXPORTER_OTLP_ENDPOINT'],
    ['Correlation', 'X-Request-ID, identical to the trace id'],
    ['Logging', 'structured, levelled'],
    ['Health', 'a probe registry and four states'],
    ['Profiling', 'Pulse, per-request'],
    ['Audit', 'an activity log with a SHA-256 hash chain'],
    ['Unconfigured', 'tracing inert, health still answers'],
  ],

  api: {
    groups: [
      {
        title: 'Operational endpoints',
        rows: [
          { method: 'GET', path: '/api/health', what: 'Component states, overall status, hub and queue counters' },
          { method: 'GET', path: '/api/v1/admin/activity', what: 'The activity log, filtered and paginated' },
          { method: 'GET', path: '/api/v1/admin/activity/integrity', what: 'Verify the hash chain' },
        ],
      },
    ],
    samples: [
      {
        title: 'A health response that distinguishes absent from broken',
        language: 'json',
        code: `{
  "status": "ok",
  "components": {
    "database": { "state": "ok", "detail": "postgres 16.2" },
    "cache":    { "state": "off", "detail": "REDIS_URL not set" },
    "storage":  { "state": "ok", "detail": "local disk at ./storage" },
    "mail":     { "state": "off", "detail": "MAIL_MAILER=log" },
    "queue":    { "state": "unknown", "detail": "no poll taken yet" }
  },
  "realtime": { "sockets": 412, "delivered": 918244, "dropped": 3 }
}`,
      },
    ],
  },

  lowLevel: {
    classes: [
      {
        name: 'health.Register',
        file: 'internal/health/health.go',
        what: 'Takes a name and a probe function. The registry is what lets a plugin report itself without the handler knowing it exists.',
      },
      {
        name: 'health.Snapshot',
        what: 'Runs every probe and returns the component states. The handler merges them into the response and adds nothing of its own.',
      },
      {
        name: 'health.Overall',
        what: 'Derives the overall status. Only degraded lowers it; off and unknown do not. One function, so a new component cannot invent its own rule.',
      },
      {
        name: 'tracing.Middleware',
        file: 'internal/tracing/middleware.go',
        what: 'The per-request span, which also makes the request identifier and the trace identifier the same string.',
      },
      {
        name: 'tracing.Init',
        file: 'internal/tracing/tracing.go',
        what: 'Sets up the exporter when an endpoint is configured and returns a no-op provider otherwise, so nothing downstream branches on whether tracing is on.',
      },
    ],
    principles: [
      {
        label: 'Four states beat a boolean',
        text: 'the two-state version made the dashboard guess, and it guessed differently per component. Putting the distinction in the response left the page nothing to infer.',
      },
      {
        label: 'Not knowing is not failing',
        text: 'unknown exists so that a probe which could not run does not page anybody. Alert fatigue is a monitoring failure, not an operator failure.',
      },
      {
        label: 'One key joins everything',
        text: 'the request identifier is the trace identifier. Two identifiers for one request doubles the work of every investigation.',
      },
      {
        label: 'A registry, not a switch statement',
        text: 'components register themselves. The alternative is a handler that has to be edited for every new component, which is where storage got left out.',
      },
      {
        label: 'Inert by default',
        text: 'tracing with no collector costs nothing. The moment instrumentation requires infrastructure, a project without that infrastructure removes the instrumentation.',
      },
    ],
    patterns: [
      ['Correlation identifier', 'one key joining logs, traces and responses'],
      ['Registry', 'self-registering probes instead of a central list'],
      ['Four-state health', 'absent, broken, unknown and fine as distinct facts'],
      ['No-op provider', 'instrumentation that is free when unconfigured'],
      ['Hash chain', 'an audit log that can be shown to be unedited'],
    ],
  },

  scaling: [
    'Logging is the dominant cost and it is linear in traffic. Levels and sampling are the only controls that matter, and debug left on in production is the single most common way to pay for it.',
    'Traces need sampling above modest traffic. Head sampling is cheap and loses the interesting requests; tail sampling keeps the slow and failed ones and costs a collector that buffers.',
    'Health probes are polled forever, so a probe that does real work is a permanent load increase. Caching a probe result for a few seconds is almost always correct.',
    'The activity log grows faster than any other table and eventually needs partitioning or archival to cold storage.',
    'The hash chain means verification walks the rows, which gets slower as the log grows. Verifying the recent tail is the practical answer, with full verification as an offline job.',
    'With nothing configured the whole system still works: health answers, logs go to stdout, tracing is inert. A single-binary deployment is not asked to run a collector.',
  ],

  bottlenecks: {
    problems: [
      {
        label: 'Log volume and cost',
        text: 'at scale the log bill can exceed the compute bill, and the usual cause is a debug level or a line per database query.',
      },
      {
        label: 'Expensive health probes',
        text: 'a probe that counts rows runs several times a second forever, and it is load that was added to detect load.',
      },
      {
        label: 'Cardinality explosion',
        text: 'a metric labelled with a user identifier or a path containing one creates a time series per value, and metric stores fall over on cardinality rather than volume.',
      },
      {
        label: 'Sampling losing the evidence',
        text: 'head sampling at 10% means the eleven second request is 90% likely not to have been recorded, which is exactly the request that was worth recording.',
      },
      {
        label: 'Activity log growth',
        text: 'the fastest growing table, in the same database as the application data it is competing with.',
      },
    ],
    improvements: [
      {
        label: 'Sample logs, not just traces',
        text: 'log every error and a fraction of successes. The successes are almost all identical, and the errors are all different.',
      },
      {
        label: 'Cache probe results',
        text: 'a few seconds of staleness in a health response is not a problem. Permanent query load to produce it is.',
      },
      {
        label: 'Template the labels',
        text: 'label by route pattern rather than by path. The pattern is a fixed set; the path is unbounded.',
      },
      {
        label: 'Tail sample on duration and status',
        text: 'buffer the spans, keep the slow and the failed. It costs a collector and it keeps precisely the traces anybody will ever look at.',
      },
      {
        label: 'Partition and archive the activity log',
        text: 'by month, with old partitions moved to object storage. The hash chain is per partition so verification stays bounded.',
      },
    ],
  },

  seeAlso: [
    { title: 'Observability reference', href: '/docs/operations/observability' },
    { title: 'Audit logging', href: '/docs/systems/audit-logging' },
  ],
}

export const AUDIT: SystemDesign = {
  slug: 'audit-logging',
  name: 'Audit Logging',
  tagline:
    'Recording who did what, in a form that can be shown afterwards not to have been edited.',
  group: 'Operations',
  packages: ['internal/audit', 'internal/appendonly', 'internal/models'],

  problem: {
    text: [
      'Somebody deleted a customer record. Six weeks later it matters who, and when, and what the record said. Without an audit log there is no answer: the row is gone, the application logs rotated after fourteen days, and the only remaining evidence is a backup that shows the record existed and nothing about its removal.',
      'An audit log solves that, and then introduces a question that an ordinary log does not have to answer: can the log itself be trusted? An audit trail that an administrator can edit proves nothing against an administrator, which is a large share of what an audit trail exists for. The same applies to the application: a bug or an injected statement that can delete audit rows makes every other row meaningless, because there is no way to tell a complete log from a pruned one.',
      'The answer is a hash chain. Each entry includes the hash of the one before it, so changing or removing any entry invalidates every entry after it, and that invalidation is detectable with one pass. It does not prevent tampering, because nothing in the database can. It makes tampering visible, which is a different and achievable goal.',
      'Getting that right turned out to depend on a detail nobody expects. The chain originally hashed a timestamp more precise than Postgres and MySQL actually store, so the value written and the value read back differed and the chain failed verification on its very first row, on every such deployment.',
    ],
    capabilities: [
      'Record every action worth attributing: the actor, the action, the target, the time and the context.',
      'Chain the entries, so an alteration anywhere is detectable.',
      'Verify the chain on demand, and say where it breaks.',
      'Make a table append-only, so updates and deletes are refused at the database layer.',
      'Record a reseal, when a chain genuinely has to be rebuilt, as an entry in the chain.',
      'Query the log like any other resource: filter by actor, action, target and date.',
      'Keep the record through a soft delete, so the deletion itself is the thing recorded.',
    ],
  },

  functional: [
    'An activity log row per recorded action, with actor, method, path, target and time.',
    'A SHA-256 hash of each entry including the previous entry’s hash.',
    'A verification pass reporting whether the chain holds and where it first fails.',
    'An append-only guard, installed on registered models, refusing updates and deletes.',
    'A reseal operation, recorded as a SECURITY entry naming itself.',
    'An admin page over the log, with filters and the integrity check.',
    'Hashed values normalised to what the column actually stores, so precision does not break the chain.',
  ],

  nonFunctional: [
    {
      label: 'Tamper evident, not tamper proof',
      text: 'nothing inside the database can stop a sufficiently privileged write. Making it detectable is achievable and is what the chain delivers.',
    },
    {
      label: 'Append only',
      text: 'the guard is on the model, so an update or delete is refused wherever it comes from: a handler, a job, a migration, the studio.',
    },
    {
      label: 'Precision-safe hashing',
      text: 'every hashed value is normalised to what the column stores. Postgres keeps microseconds and MySQL milliseconds, and a nanosecond timestamp in the hash input made verification fail on row one.',
    },
    {
      label: 'Verifiable in bounded time',
      text: 'verification is a single ordered pass. It gets slower as the log grows, which is why the recent tail is what gets checked routinely.',
    },
    {
      label: 'Queryable',
      text: 'the log is a table with indexes, not a text file. The question is always "what did this person do" or "what happened to this record", and both need an index.',
    },
    {
      label: 'Honest about a reseal',
      text: 'a chain that genuinely must be rebuilt records the rebuild in itself. A silent reseal would make the chain worthless.',
    },
  ],

  capacity: {
    assumptions: [
      ['Audited actions per day', '~200,000'],
      ['Entry size', '~500 bytes with context'],
      ['Hash cost', 'microseconds per entry'],
      ['Verification throughput', '~200,000 rows/second'],
      ['Retention', '7 years in a regulated context, 1 year otherwise'],
    ],
    estimates: [
      {
        label: 'Growth',
        working: [
          '200,000 x 500 bytes = ~100 MB/day',
          'x 365 = ~36 GB/year',
          'x 7 years = ~250 GB',
        ],
        note: 'Which is why partitioning matters. A single table of a quarter of a terabyte that is append-only and queried by date is the textbook case for it.',
      },
      {
        label: 'Write cost',
        working: [
          'read the previous hash: one indexed read',
          'SHA-256 over ~500 bytes: microseconds',
          'insert: one row',
        ],
        note: 'The previous-hash read is the real cost, and it serialises writes to the log: two entries cannot be appended concurrently without one of them having a stale predecessor.',
      },
      {
        label: 'Verification time',
        working: [
          '73M rows after a year',
          'at ~200,000 rows/s = ~6 minutes',
          'the last 10,000 rows: well under a second',
        ],
        note: 'Full verification is an offline job. Verifying the tail is the thing that can be done from a button in the admin.',
      },
    ],
  },

  highLevel: {
    intro:
      'An append-only table, a hash over each entry that includes its predecessor, and a verification pass.',
    components: [
      {
        label: 'Activity log model',
        text: 'the table, declared like any other so migrations and backups cover it. Actor, action, target, time, context and the chain fields.',
      },
      {
        label: 'Chain writer',
        text: 'reads the previous hash, computes this entry’s, inserts. Normalises every hashed value to what the column stores.',
      },
      {
        label: 'Verifier',
        text: 'walks the entries in order, recomputing each hash, and reports the first row where the chain breaks.',
      },
      {
        label: 'Append-only guard',
        text: 'a GORM hook installed on registered models, refusing updates and deletes. Installed in every project whether or not anything uses it yet, so the feature works the day somebody needs it.',
      },
      {
        label: 'Reseal',
        text: 'rebuilds a chain that cannot verify, and records that it did so as an entry in the new chain.',
      },
      {
        label: 'Admin view',
        text: 'the log with filters, and the integrity check as a button.',
      },
    ],
    flow: {
      title: 'An entry joining the chain, and the chain being checked',
      nodes: [
        { id: 'action', label: 'Audited Action', col: 0, row: 0 },
        { id: 'writer', label: 'Chain Writer', col: 1, row: 0, accent: true },
        { id: 'prev', label: 'Previous Hash', col: 2, row: 0, accent: true },
        { id: 'hash', label: 'SHA-256', col: 3, row: 0, accent: true },
        { id: 'norm', label: 'Normalise Values', col: 2, row: 1, accent: true },
        { id: 'row', label: 'Append Row', col: 3, row: 1 },
        { id: 'guard', label: 'Append-Only Guard', col: 1, row: 1, accent: true },
        { id: 'verify', label: 'Verify Chain', col: 3, row: 2 },
        { id: 'report', label: 'First Break', col: 2, row: 2 },
      ],
      edges: [
        { from: 'action', to: 'writer', step: 1 },
        { from: 'writer', to: 'prev', step: 2 },
        { from: 'prev', to: 'hash', step: 3 },
        { from: 'writer', to: 'norm', step: 4, bend: 'v' },
        { from: 'hash', to: 'row', step: 5 },
        { from: 'guard', to: 'row', step: 6, bend: 'h' },
        { from: 'row', to: 'verify', step: 7 },
        { from: 'verify', to: 'report', step: 8 },
      ],
      steps: [
        'Something happens that is worth attributing: a record deleted, a role granted, a setting changed, a sign-in from a new device.',
        'The chain writer reads the hash of the most recent entry. This is what links the new entry to everything before it, and it is also why appends serialise.',
        'The entry’s own hash is computed over its fields together with the previous hash. Changing any earlier entry changes every hash after it, which is the entire mechanism.',
        'Every hashed value is first normalised to what the column will actually store. Postgres keeps microseconds and MySQL milliseconds, and hashing a nanosecond timestamp meant the value read back never matched the value hashed. The chain failed on its first row on every such deployment until this was fixed.',
        'The row is appended, carrying its own hash and its predecessor’s.',
        'The append-only guard refuses any update or delete on the table. It is a model hook rather than a service rule, so it applies to a handler, a job, a migration and the database browser alike.',
        'Verification walks the entries in order and recomputes each hash. A chain that holds proves no entry was altered or removed since it was written.',
        'A break is reported with the row it starts at, which is the useful answer: everything before it is intact, and something happened at that point. A genuine rebuild is possible, and it writes a SECURITY entry naming itself, because a silent reseal would make the whole chain worthless.',
      ],
    },
    dataFlow: [
      'The chain links entries by hash, so the detection property survives a backup and restore: the restored log verifies, and it verifies as of the backup.',
      'Append-only is enforced at the model layer, which is the only layer every writer passes through.',
      'A soft delete of an audited record keeps the record and adds an entry about the deletion, so the deletion is itself the audited fact.',
      'Normalising hashed values to column precision is not a detail. It is the difference between a chain that verifies on Postgres and one that never has.',
    ],
  },

  stack: [
    ['Hash', 'SHA-256 over the entry plus the previous hash'],
    ['Storage', 'a table in the application database'],
    ['Enforcement', 'an append-only GORM guard on registered models'],
    ['Verification', 'an ordered pass reporting the first break'],
    ['Recovery', 'a reseal, recorded in the chain it rebuilds'],
    ['Admin', 'a filtered log view and an integrity button'],
  ],

  dataModel: {
    entities: [
      {
        name: 'activity_logs',
        fields: [
          ['id', 'primary key'],
          ['user_id', 'who, or empty for an unauthenticated action'],
          ['method', 'the kind of action, including SECURITY for chain events'],
          ['path', 'what was acted on'],
          ['target_type / target_id', 'the record, where there is one'],
          ['context', 'JSON: address, agent, changed fields'],
          ['created_at', 'stored at the precision the engine supports'],
          ['entry_hash', 'SHA-256 of this entry plus the previous hash'],
          ['prev_hash', 'the predecessor, empty on the first row'],
        ],
        note: 'created_at is in the hash, which is why its precision matters so much. It is normalised to what the column stores before hashing, not after.',
      },
    ],
    storage: [
      'In the application database because it is queried with application data, by actor and by target.',
      'Append-only in practice and enforced in code, so a restore of a backup produces a log that verifies as of that backup.',
      'The fastest growing table in most projects, and the clearest candidate for partitioning by month.',
    ],
  },

  api: {
    groups: [
      {
        title: 'Audit log',
        rows: [
          { method: 'GET', path: '/api/v1/admin/activity', what: 'Filter by actor, action, target and date' },
          { method: 'GET', path: '/api/v1/admin/activity/integrity', what: 'Verify the chain, report the first break' },
          { method: 'POST', path: '/api/v1/admin/activity/reseal', what: 'Rebuild a broken chain, recorded in itself' },
        ],
      },
    ],
  },

  lowLevel: {
    classes: [
      {
        name: 'audit.Append',
        file: 'internal/audit/chain.go',
        what: 'Writes one chained entry. Reads the previous hash, normalises the hashed values to column precision, computes and inserts.',
      },
      {
        name: 'audit.VerifyChain',
        what: 'Walks the entries recomputing hashes. Returns a status and the first row that does not match, because "it is broken" without a position is not actionable.',
      },
      {
        name: 'appendonly.Install',
        file: 'internal/appendonly/appendonly.go',
        what: 'Registers the guard for a model. Installed in every project from the start: the two calls that install it are no-ops with nothing registered, which is what lets the append-only flag work on the day somebody uses it rather than after a round of wiring.',
      },
      {
        name: 'Reseal',
        what: 'Rebuilds a chain and records a SECURITY entry naming itself. A reseal of a chain that already verifies is refused, because there would be no reason for it.',
      },
    ],
    principles: [
      {
        label: 'Detect rather than prevent',
        text: 'nothing in the database can stop a privileged write. A chain makes the write visible, which is the achievable version of the goal.',
      },
      {
        label: 'Enforce where every writer passes',
        text: 'a model hook covers handlers, jobs, migrations and the database browser. A service rule covers the callers who remembered.',
      },
      {
        label: 'Hash what will be stored',
        text: 'the hash input must survive the round trip through the column. This one cost every Postgres and MySQL deployment its chain until it was found.',
      },
      {
        label: 'Record the exception',
        text: 'a reseal is written into the chain. An audit system whose own repairs are invisible is not an audit system.',
      },
      {
        label: 'Ship the mechanism before the use',
        text: 'the append-only package is in every project, inert. Infrastructure that has to be added before a feature can be used is infrastructure the feature does without.',
      },
    ],
    patterns: [
      ['Hash chain', 'each entry committing to its predecessor'],
      ['Append-only store', 'updates and deletes refused at the model layer'],
      ['Tamper evidence', 'detection in place of prevention'],
      ['Event sourcing, partially', 'the log of what happened kept beside the current state'],
    ],
  },

  scaling: [
    'Writes serialise on reading the previous hash, which caps the append rate. For audit volumes that is not a constraint; for per-request logging of everything it would be.',
    'Verification is linear in the log, so full verification becomes an offline job and the routine check is of the recent tail.',
    'Partitioning by month keeps queries and verification bounded, with the chain verified per partition and the boundary hash linking them.',
    'The table is the fastest growing in most projects, and the retention period is usually a compliance decision rather than a technical one.',
    'Archiving old partitions to object storage keeps the database small while keeping the record, which is what long retention periods actually require.',
    'Hashing itself never matters: SHA-256 over 500 bytes is microseconds, and nothing about the cost of the chain is the hash.',
  ],

  bottlenecks: {
    problems: [
      {
        label: 'Serialised appends',
        text: 'every entry needs the previous hash, so the chain has a single writer by construction and a burst of audited actions queues.',
      },
      {
        label: 'Verification time',
        text: 'a full pass over tens of millions of rows takes minutes, which makes it unusable as an on-demand check from a button.',
      },
      {
        label: 'A break with no explanation',
        text: 'the chain says where it first fails and cannot say why. A restore from backup, a replication artefact and a deliberate edit all look identical.',
      },
      {
        label: 'Table growth under long retention',
        text: 'seven years of audit data in the operational database competes with the application for everything.',
      },
      {
        label: 'Nothing anchored outside',
        text: 'an attacker who can rewrite the whole table can recompute the whole chain. Internal consistency is not external proof.',
      },
    ],
    improvements: [
      {
        label: 'Chain per partition',
        text: 'monthly chains linked by their boundary hashes. Appends within a month still serialise, but verification and archival both become bounded.',
      },
      {
        label: 'Verify the tail routinely, the whole thing nightly',
        text: 'the last few thousand rows from a button, the full pass as a scheduled job that alerts on a break.',
      },
      {
        label: 'Record the operational events too',
        text: 'a restore, a reseal and a migration written into the chain turn an unexplained break into an explained one.',
      },
      {
        label: 'Archive to object storage',
        text: 'closed partitions moved out with their verification result, keeping the record without keeping the rows in the hot database.',
      },
      {
        label: 'Anchor the chain externally',
        text: 'publish the periodic head hash somewhere outside the system, to a log service or a signed note. That is what makes the chain evidence against somebody who controls the database.',
      },
    ],
  },

  seeAlso: [
    { title: 'Audit log reference', href: '/docs/security/audit-log' },
    { title: 'Observability', href: '/docs/systems/observability' },
  ],
}

export const FILE_STORAGE: SystemDesign = {
  slug: 'file-storage',
  name: 'File Storage',
  tagline:
    'Accepting an upload, keeping it somewhere that survives a redeploy, and serving it back to the people allowed to see it.',
  group: 'Operations',
  packages: ['internal/storage', 'internal/files', 'internal/media'],

  problem: {
    text: [
      'A user uploads an avatar. Writing it to the container’s filesystem works until the container restarts, and then the avatar is gone and so is every other upload. Writing it to the database as bytes works until the table is sixty gigabytes and every backup takes an hour. Neither is a storage design; both are what happens when storage is not designed.',
      'Object storage is the answer, and it brings its own set of decisions that are easy to get wrong in ways that only show up later. A browser that uploads through the API means every file passes through a request, holding a connection for its duration, bounded by the request body limit and the proxy timeout. A file served through the API means every read costs an application request. A file served directly from a public bucket means access control has been given up entirely, which is fine for an avatar and not for an invoice.',
      'And there is the local case. A small deployment should not need a bucket. It should be able to write files to a directory and serve them, with the same code, so that a single binary and SQLite remains a real option rather than a demo.',
      'That last part was broken here for a long time in an interesting way: the local driver worked, and the fallback that was supposed to select it when no bucket was configured never fired, so a project with no MinIO got errors instead of a working directory.',
    ],
    capabilities: [
      'Store a file through one interface, whichever backend is configured.',
      'Support S3-compatible buckets and a local directory with the same code.',
      'Accept an upload through the API, with a size limit and a type check.',
      'Offer a pre-signed URL so a browser can upload straight to the bucket.',
      'Serve a private file through an authorised route, and a public one directly.',
      'Derive thumbnails and sized variants without blocking the upload.',
      'Record every file as a row, so it can be listed, attributed and deleted.',
      'Delete the object when the row goes, and not before.',
    ],
  },

  functional: [
    'A Disk interface: put, get, delete, exists, URL.',
    'Drivers for S3, MinIO, R2 and B2, and a local directory driver.',
    'Driver selection from STORAGE_DRIVER, with local as a working default.',
    'Multipart upload handling with a configured size limit.',
    'Content type validation against an allowlist, by sniffed type and not only by extension.',
    'Pre-signed PUT URLs valid for an hour, with a clear error from drivers that cannot.',
    'An uploads table: key, size, type, owner, created time.',
    'Image processing in a background job, producing sized variants.',
    'An admin files page listing, previewing and deleting.',
  ],

  nonFunctional: [
    {
      label: 'Survives a redeploy',
      text: 'the container filesystem is scratch space. Anything that must outlive a deploy goes to a bucket or to a mounted volume, and that is a decision made once rather than per feature.',
    },
    {
      label: 'Local must actually work',
      text: 'a project with no bucket configured writes to a directory and serves from the API. Not as a degraded mode, as a supported deployment.',
    },
    {
      label: 'Private by default',
      text: 'a bucket is not public unless a file is meant to be. The default must be the safe one, because the unsafe one is invisible until it is indexed.',
    },
    {
      label: 'Bounded uploads',
      text: 'a size limit and a type allowlist, enforced in the handler. Without them an upload endpoint is a way to fill a disk.',
    },
    {
      label: 'Derived work is asynchronous',
      text: 'thumbnailing happens in a job. An upload that waits for image processing is an upload that times out on a large file.',
    },
    {
      label: 'The row is the record',
      text: 'a file is a row that references an object. Orphaned objects and orphaned rows are both bugs, and the row is what makes either detectable.',
    },
  ],

  capacity: {
    assumptions: [
      ['Uploads per day', '~12,000'],
      ['Average upload size', '1.8 MB'],
      ['Reads per upload over its life', '~40'],
      ['Thumbnail variants per image', '3'],
      ['Retention', 'indefinite'],
    ],
    estimates: [
      {
        label: 'Storage growth',
        working: [
          '12,000 x 1.8 MB = ~21 GB/day of originals',
          'variants add ~15%: ~24 GB/day',
          'x 365 = ~8.8 TB/year',
        ],
        note: 'Which is the argument for lifecycle rules rather than for deleting things. Moving year-old objects to infrequent access costs nothing in code and most of the bill.',
      },
      {
        label: 'Upload bandwidth through the API',
        working: [
          '12,000/day is ~0.14/s average, but bursts to ~5/s',
          '5 x 1.8 MB = 9 MB/s through the API',
          'each one holding a request for its duration',
        ],
        note: 'This is what pre-signed URLs remove. The bytes go browser to bucket, and the API handles only the authorisation and the row.',
      },
      {
        label: 'Read bandwidth',
        working: [
          '12,000 x 40 reads = 480,000 reads/day',
          'x 1.8 MB = ~860 GB/day egress',
        ],
        note: 'Egress is usually the largest line on an object storage bill, and a CDN in front of it is the only thing that changes the number materially.',
      },
      {
        label: 'Thumbnail work',
        working: [
          '12,000 images x 3 variants = 36,000 operations/day',
          'at ~400 ms each = 4 hours of CPU',
          'spread across workers, a fraction of one core',
        ],
      },
    ],
  },

  highLevel: {
    intro:
      'One interface, several drivers, a row per file, and the derived work in a job.',
    components: [
      {
        label: 'Disk interface',
        text: 'what every driver implements. New code takes a Disk rather than a concrete store, which is what makes the local case first class instead of special.',
      },
      {
        label: 'Drivers',
        text: 'S3-compatible for s3, minio, r2 and b2, and local for a directory on this machine served from the API.',
      },
      {
        label: 'Storage service',
        text: 'holds one Disk, chosen by configuration, and reports which driver it is. The report matters: "local disk at ./storage" in a health response answers a question that otherwise takes an hour.',
      },
      {
        label: 'Upload handler',
        text: 'multipart, with a size limit and a type allowlist. Writes the object and the row in that order, so a failed write never leaves a row pointing at nothing.',
      },
      {
        label: 'Pre-signed URLs',
        text: 'a PUT URL valid for an hour so the browser uploads directly. Drivers that cannot say so with a named error rather than failing obscurely.',
      },
      {
        label: 'Image processing',
        text: 'a job producing sized variants. Asynchronous because a large image takes seconds and an upload should not.',
      },
      {
        label: 'Uploads table',
        text: 'the row per file: key, size, type, owner. What makes a file listable, attributable and deletable.',
      },
    ],
    flow: {
      title: 'An upload, two ways, and the read path',
      nodes: [
        { id: 'browser', label: 'Browser', col: 0, row: 0 },
        { id: 'handler', label: 'Upload Handler', col: 1, row: 0 },
        { id: 'check', label: 'Size + Type Check', col: 2, row: 0, accent: true },
        { id: 'disk', label: 'Disk', col: 3, row: 0, accent: true },
        { id: 'presign', label: 'Presigned PUT', col: 1, row: 1, accent: true },
        { id: 'row', label: 'Upload Row', col: 2, row: 1 },
        { id: 'job', label: 'Thumbnail Job', col: 3, row: 1 },
        { id: 'serve', label: 'Authorised Read', col: 2, row: 2, accent: true },
        { id: 'cdn', label: 'CDN or API', col: 3, row: 2 },
      ],
      edges: [
        { from: 'browser', to: 'handler', step: 1 },
        { from: 'handler', to: 'check', step: 2 },
        { from: 'check', to: 'disk', step: 3 },
        { from: 'handler', to: 'presign', step: 4, bend: 'v' },
        { from: 'disk', to: 'row', step: 5, bend: 'v' },
        { from: 'row', to: 'job', step: 6 },
        { from: 'row', to: 'serve', step: 7, bend: 'v' },
        { from: 'serve', to: 'cdn', step: 8 },
      ],
      steps: [
        'A browser uploads. For a small file that goes through the API, which is simple and costs a held request for the duration of the transfer.',
        'The handler checks the size against the configured limit and the content type against an allowlist, by sniffing the bytes rather than trusting the extension, because the extension is supplied by the uploader.',
        'The object is written to the configured Disk: a bucket, or a directory on this machine. The same code either way, which is what makes a single-binary deployment a real option.',
        'The alternative path, for large files: the API returns a pre-signed PUT valid for an hour and the browser uploads straight to the bucket. The API authorises and records; it never carries the bytes. Drivers that cannot pre-sign return a named error rather than failing in some other way.',
        'A row is written after the object exists, never before. A row pointing at an object that was never written is a broken link; an object with no row is garbage that a sweep can find.',
        'Derived work is queued. Three thumbnail sizes take a couple of seconds, which is fine in a worker and not in a request.',
        'A read is authorised against the row. This is where private and public diverge: a private file is checked and streamed or given a short-lived signed URL, a public one is served directly.',
        'Public files go through a CDN, because egress is the largest part of an object storage bill and caching is the only thing that changes it. Private files are served by the API or by signed URLs with a short life.',
      ],
    },
    dataFlow: [
      'Object first, row second. The reverse ordering produces rows pointing at nothing, which is the harder of the two failure modes to detect.',
      'Deleting removes the row and the object, in that order, so a failure leaves an orphaned object rather than a broken reference.',
      'Keys are generated, not taken from the upload. A user-supplied filename is a path traversal attempt waiting to happen, and it is kept as a display name only.',
      'The driver name is reported in health, because "which store is this project using" is otherwise answered by reading configuration on a machine nobody has access to.',
    ],
  },

  stack: [
    ['Interface', 'one Disk, several drivers'],
    ['Object stores', 'S3, MinIO, Cloudflare R2, Backblaze B2'],
    ['Local', 'a directory, served from the API, a supported deployment'],
    ['Selection', 'STORAGE_DRIVER'],
    ['Direct upload', 'pre-signed PUT, one hour'],
    ['Processing', 'sized variants in a background job'],
    ['Record', 'an uploads table'],
  ],

  api: {
    groups: [
      {
        title: 'Uploads',
        rows: [
          { method: 'POST', path: '/api/v1/uploads', what: 'Multipart upload through the API' },
          { method: 'POST', path: '/api/v1/uploads/presign', what: 'A pre-signed PUT for a direct browser upload' },
          { method: 'GET', path: '/api/v1/uploads', what: 'List, scoped to the caller' },
          { method: 'GET', path: '/api/v1/uploads/:id', what: 'Authorised read of a private file' },
          { method: 'DELETE', path: '/api/v1/uploads/:id', what: 'Remove the row and the object' },
        ],
      },
    ],
  },

  dataModel: {
    entities: [
      {
        name: 'uploads',
        fields: [
          ['id', 'primary key'],
          ['key', 'the object key, generated, never user supplied'],
          ['original_name', 'what the uploader called it, for display only'],
          ['content_type', 'as sniffed, not as claimed'],
          ['size', 'bytes'],
          ['user_id', 'who uploaded it, for ownership scoping'],
          ['variants', 'JSON: the derived sizes and their keys'],
          ['created_at', 'when'],
        ],
        note: 'original_name is display only and is never part of a path. Using it as a key is how a path traversal gets in.',
      },
    ],
    storage: [
      'Objects live in the bucket or the directory; only metadata is in the database. Bytes in a database make every backup the size of every upload.',
      'Lifecycle rules on the bucket move old objects to cheaper classes, which is a configuration change rather than a code one.',
      'The local driver’s directory must be a mounted volume in a container, or it is scratch space with a longer name.',
    ],
  },

  lowLevel: {
    classes: [
      {
        name: 'storage.Disk',
        file: 'internal/storage/disk.go',
        what: 'The interface every driver implements. New code takes a Disk, which is what keeps the local case from being a special case.',
        methods: ['Put', 'Get', 'Delete', 'Exists', 'URL'],
      },
      {
        name: 'storage.NewLocal',
        file: 'internal/storage/local.go',
        what: 'Keeps files in a directory and serves them from the API. Always worked; the fallback that selects it when no bucket is configured did not fire until v3.370.0.',
      },
      {
        name: 'storage.Storage',
        what: 'Holds one Disk chosen by STORAGE_DRIVER and reports which. The report is in the health response, which is where somebody finds out what a deployment is actually doing.',
        methods: ['Driver', 'PresignPutURL'],
      },
      {
        name: 'ErrPresignUnsupported',
        what: 'Returned by drivers that cannot pre-sign, such as local. A named error so the caller can fall back to uploading through the API rather than guessing.',
      },
    ],
    principles: [
      {
        label: 'One interface, no special cases',
        text: 'the local driver implements the same interface as S3. The moment local is a branch rather than a driver, it stops being tested.',
      },
      {
        label: 'Generate the key',
        text: 'a user-supplied filename in a path is a traversal. Keep it as a display name and generate the key.',
      },
      {
        label: 'Object before row',
        text: 'the ordering decides which failure you get. An orphaned object is findable; a row pointing at nothing is a broken page.',
      },
      {
        label: 'Report the driver',
        text: 'the configuration is on a machine nobody can read. The health response is where the answer belongs.',
      },
      {
        label: 'Name what a driver cannot do',
        text: 'a specific error for an unsupported pre-sign lets the caller fall back. A generic failure makes it look like a bug.',
      },
    ],
    patterns: [
      ['Strategy', 'the Disk interface with a driver per store'],
      ['Pre-signed URL', 'the bytes bypass the application entirely'],
      ['Metadata and blob split', 'rows in the database, objects in the store'],
      ['Asynchronous derivation', 'variants built in a job after the upload'],
    ],
  },

  scaling: [
    'Object storage scales without the application participating, which is most of why it is the answer.',
    'Pre-signed uploads remove the API from the data path, so upload capacity stops being an application concern at all.',
    'Egress is the dominant cost at volume, and a CDN in front of public files is the only change that moves it significantly.',
    'Private files served through the API cost a request per read. Short-lived signed URLs give most of the benefit of direct serving while keeping authorisation.',
    'Thumbnail work scales with the worker pool and is cheap. Generating variants on demand and caching them is the alternative, and it trades storage for latency on first view.',
    'The local driver scales to one machine and its disk, which is the correct limit for the deployment it exists to serve.',
  ],

  bottlenecks: {
    problems: [
      {
        label: 'Uploads through the API',
        text: 'every upload holds a request, bounded by the body limit and the proxy timeout. A large file fails in a way that looks like a network problem.',
      },
      {
        label: 'Egress cost',
        text: 'at hundreds of gigabytes a day, bandwidth is the bill, and nothing in the application code is where it is decided.',
      },
      {
        label: 'Orphans in both directions',
        text: 'objects with no row accumulate silently from failed requests; rows with no object break pages.',
      },
      {
        label: 'Public buckets',
        text: 'a bucket made public to make serving simple is a bucket whose entire contents are enumerable by anybody who guesses the naming scheme.',
      },
      {
        label: 'Unbounded uploads',
        text: 'no size limit and no type check turns an upload endpoint into a way to fill a disk and host arbitrary content.',
      },
    ],
    improvements: [
      {
        label: 'Pre-sign everything large',
        text: 'the browser uploads to the bucket, the API authorises and records. Removes the body limit, the timeout and the bandwidth from the application.',
      },
      {
        label: 'CDN the public prefix',
        text: 'cache at the edge. It is a configuration change and it is the only thing that materially changes the egress bill.',
      },
      {
        label: 'Sweep for orphans',
        text: 'a scheduled job comparing keys to rows in both directions, reporting rather than deleting until the counts are understood.',
      },
      {
        label: 'Keep the bucket private, sign the reads',
        text: 'short-lived signed URLs for private files. Authorisation stays in the application and the bytes still come from the store.',
      },
      {
        label: 'Limit and sniff',
        text: 'a size cap and an allowlist checked against the sniffed type. The extension is supplied by the uploader and means nothing.',
      },
    ],
  },

  seeAlso: [
    { title: 'Storage reference', href: '/docs/batteries/storage' },
    { title: 'Background jobs', href: '/docs/systems/background-jobs' },
  ],
}

export const IMPORT_EXPORT: SystemDesign = {
  slug: 'import-and-export',
  name: 'Import and Export',
  tagline:
    'Taking a spreadsheet somebody has and turning it into rows, and giving them a file back, without either one taking the application down.',
  group: 'Operations',
  packages: ['internal/imports', 'internal/export'],

  problem: {
    text: [
      'Every application that replaces a spreadsheet has to start by reading one. The first thing a new customer asks is whether they can bring their data, and the answer decides whether they adopt the thing. The last thing they ask, often years later, is whether they can get it out.',
      'A naive importer reads a file into memory, loops over the rows and saves each one. That works for the test file with twelve rows and fails for the real one with ninety thousand, in several ways at once: memory, because the file and the parsed rows are both resident; the connection pool, because the loop holds a connection for its whole run and several concurrent imports drain the pool that every request shares; and partial failure, because row 40,000 is invalid and the first 39,999 are already committed with no record of what happened.',
      'Then there is the part that is specific to importing rather than to bulk writing. A spreadsheet refers to things by name, not by identifier. The contacts file says the group is "Suppliers", and the database wants a group identifier. Resolving that per row is one query per row, which is ninety thousand queries, and resolving it wrongly is worse: treating a lookup failure as a missing record silently creates duplicates of things that already exist under a slightly different name.',
      'Export has a smaller version of the same shape. Building the whole file in memory before sending it works until the result set is large, and then it does not.',
    ],
    capabilities: [
      'Read a large file without holding it all in memory.',
      'Write in batches, not row by row.',
      'Limit how many imports run at once, because each one holds a connection for its whole run.',
      'Resolve a relationship by its natural key, with a cache so a repeated name costs one lookup.',
      'Distinguish "not found" from "the lookup failed", and fail the row on the second.',
      'Report progress while running, and a per-row error list afterwards.',
      'Validate each row the way the API would, rather than inventing a second set of rules.',
      'Stream an export, so the response starts before the file is finished.',
    ],
  },

  functional: [
    'CSV import per resource, generated with the resource.',
    'Streaming read, batched write.',
    'A concurrency gate, from IMPORT_CONCURRENCY, that makes imports take turns.',
    'A progress callback, so the interface can show where it is.',
    'Natural key resolution for belongs-to columns, with a per-run cache.',
    'A missing related record created; any other lookup error failing the row.',
    'Per-row validation against the model’s binding rules.',
    'An error report naming the row number and what was wrong with it.',
    'Export to CSV and XLSX, streamed.',
    'An import modal and an export menu in the admin.',
  ],

  nonFunctional: [
    {
      label: 'Bounded memory',
      text: 'the file is streamed and the writes are batched, so a ninety thousand row import uses the memory of a batch rather than of a file.',
    },
    {
      label: 'Imports take turns',
      text: 'each one holds a database connection and writes in batches for its whole run. A burst of them drained the pool that every request shares, so they are gated.',
    },
    {
      label: 'A lookup failure is not a miss',
      text: 'treating a database error as "no such record" silently creates duplicates. The distinction is the difference between an import that works and one that quietly corrupts.',
    },
    {
      label: 'One set of validation rules',
      text: 'rows are validated against the same binding tags the API uses. A second set of rules for imports is a second set to keep in step.',
    },
    {
      label: 'Row-level reporting',
      text: 'a failed import that says "failed" is useless. It must say which rows and why, because the person fixing it has the spreadsheet open.',
    },
    {
      label: 'Streamed export',
      text: 'the response starts before the file is complete, so a large export does not need the whole thing in memory or a long silence.',
    },
  ],

  capacity: {
    assumptions: [
      ['Largest realistic import', '~90,000 rows'],
      ['Columns per row', '12'],
      ['Batch size', '500 rows'],
      ['Belongs-to columns resolved by name', '1 to 3'],
      ['Distinct related values', '~200'],
      ['Concurrent imports allowed', '1 by default'],
    ],
    estimates: [
      {
        label: 'Batched versus per row',
        working: [
          '90,000 individual inserts at ~1 ms = 90 seconds',
          '180 batches of 500 at ~25 ms = ~4.5 seconds',
        ],
        note: 'Twenty times faster, and it holds one connection for four seconds rather than for ninety. The second number is the one that matters to everybody else using the application.',
      },
      {
        label: 'Name resolution, cached and not',
        working: [
          'uncached: 90,000 rows x 2 columns = 180,000 lookups',
          'cached by distinct value: ~400 lookups',
        ],
        note: 'A map from name to identifier for the run. The cache is per run rather than global because a concurrent import may be creating the very records this one is looking up.',
      },
      {
        label: 'Memory',
        working: [
          'streamed read: one row at a time',
          'plus a batch of 500 x ~1 KB = ~500 KB',
          'plus the resolution cache: ~200 entries',
        ],
        note: 'Under a megabyte, independent of file size. Reading the file into memory first would be the file size plus the parsed representation, which is several times larger.',
      },
      {
        label: 'Connection pool pressure',
        working: [
          'pool of 25',
          '5 concurrent ungated imports each holding one for its run',
          '= 20% of the pool held for minutes, during business hours',
        ],
        note: 'This is the failure the gate exists for. It presents as the whole application being slow, with nothing in the logs about imports.',
      },
    ],
  },

  highLevel: {
    intro:
      'Stream the file, gate the concurrency, resolve names through a cache, write in batches, report per row.',
    components: [
      {
        label: 'Concurrency gate',
        text: 'imports wait for each other. The waiting is reported to the interface, so a queued import looks queued rather than stuck.',
      },
      {
        label: 'Streaming reader',
        text: 'one row at a time from the file, so memory is a function of the batch rather than of the upload.',
      },
      {
        label: 'Name resolver',
        text: 'a closure per belongs-to column, with a map from natural key to identifier. A missing related record is created; any other error fails the row.',
      },
      {
        label: 'Row validator',
        text: 'the model’s binding rules, so an imported row has to satisfy what an API-created row satisfies.',
      },
      {
        label: 'Batch writer',
        text: 'accumulates and inserts in batches, inside a transaction per batch, so a failure costs one batch rather than the run.',
      },
      {
        label: 'Progress and errors',
        text: 'a callback for position and a per-row error list with row numbers, which is the only form of report anybody can act on.',
      },
      {
        label: 'Export streamer',
        text: 'reads in batches and writes to the response as it goes, for CSV and XLSX.',
      },
    ],
    flow: {
      title: 'A ninety thousand row import',
      nodes: [
        { id: 'upload', label: 'CSV Upload', col: 0, row: 0 },
        { id: 'gate', label: 'Concurrency Gate', col: 1, row: 0, accent: true },
        { id: 'read', label: 'Stream Rows', col: 2, row: 0, accent: true },
        { id: 'resolve', label: 'Resolve Names', col: 3, row: 0, accent: true },
        { id: 'cache', label: 'Key Cache', col: 3, row: 1 },
        { id: 'validate', label: 'Validate Row', col: 2, row: 1, accent: true },
        { id: 'batch', label: 'Batch Insert', col: 1, row: 1, accent: true },
        { id: 'errors', label: 'Row Errors', col: 2, row: 2 },
        { id: 'report', label: 'Report', col: 0, row: 2 },
      ],
      edges: [
        { from: 'upload', to: 'gate', step: 1 },
        { from: 'gate', to: 'read', step: 2 },
        { from: 'read', to: 'resolve', step: 3 },
        { from: 'resolve', to: 'cache', step: 4 },
        { from: 'resolve', to: 'validate', step: 5, bend: 'h' },
        { from: 'validate', to: 'batch', step: 6 },
        { from: 'validate', to: 'errors', step: 7, bend: 'v' },
        { from: 'errors', to: 'report', step: 8 },
      ],
      steps: [
        'A file arrives. It is stored as an upload like any other file, so the import can be retried against the same bytes rather than needing a second upload.',
        'The import waits for the gate. One at a time by default, because each import holds a database connection and writes in batches for its whole run, and a burst of them drained the pool every request shares. The wait is reported, so a queued import does not look like a hung one.',
        'Rows are read one at a time. The file is never resident, so a file ten times larger costs the same memory and ten times the time.',
        'A belongs-to column referring to a group by name is resolved through a per-run map. A distinct name costs one lookup; the other four hundred rows mentioning it cost none.',
        'A missing related record is created. Any other error from the lookup fails the row instead of being treated as a miss, because a database error read as "not found" silently creates duplicates of records that already exist.',
        'The row is validated against the model’s binding rules, the same rules the API applies, and accumulated into the current batch. Each batch is one transaction, so a failure costs five hundred rows rather than ninety thousand.',
        'A row that fails is recorded with its number and its reason, and the run continues. Stopping on the first bad row means ninety thousand rows get imported in ninety attempts.',
        'The report lists what was created, what was skipped, and every failed row with its number and its error, which is what somebody with the spreadsheet open can act on.',
      ],
    },
    dataFlow: [
      'The resolution cache is per run. A global one would be stale as soon as another import created a related record, and the staleness would present as rows attached to the wrong parent.',
      'A batch is a transaction. Not the whole run, because a ninety thousand row transaction holds locks for its duration, and not a row, because that is the per-row cost the batching exists to remove.',
      'Validation reuses the model’s binding tags. A separate rule set for imports drifts, and the drift means an import can create a row the API would have refused.',
      'Export streams: read a batch, write a batch, flush. The client starts receiving before the query has finished.',
    ],
  },

  stack: [
    ['Import format', 'CSV, generated per resource'],
    ['Export formats', 'CSV and XLSX, streamed'],
    ['Reading', 'streamed, one row at a time'],
    ['Writing', 'batches of 500, a transaction each'],
    ['Concurrency', 'gated, IMPORT_CONCURRENCY, default 1'],
    ['Relationships', 'resolved by natural key through a per-run cache'],
    ['Validation', 'the model’s own binding tags'],
    ['Admin', 'an import modal and an export menu per resource'],
  ],

  api: {
    groups: [
      {
        title: 'Per resource',
        rows: [
          { method: 'POST', path: '/api/v1/contacts/import', what: 'Upload a CSV and start an import' },
          { method: 'GET', path: '/api/v1/contacts/import/:id', what: 'Progress, then the per-row report' },
          { method: 'GET', path: '/api/v1/contacts/export?format=csv', what: 'Streamed CSV of the current filters' },
          { method: 'GET', path: '/api/v1/contacts/export?format=xlsx', what: 'The same as a spreadsheet' },
        ],
      },
    ],
  },

  lowLevel: {
    classes: [
      {
        name: 'imports.Wait',
        file: 'internal/imports/imports.go',
        what: 'The gate. Returns a release function and reports the wait, so an import that is queued says so instead of appearing to hang.',
      },
      {
        name: 'Name resolver closure',
        what: 'Generated per belongs-to column. Caches by natural key, creates a missing related record, and fails the row on any other error rather than passing it off as a miss.',
      },
      {
        name: 'Batch writer',
        what: 'Accumulates to the batch size and inserts in one transaction. The batch is the unit of failure, which is why it is neither a row nor the whole run.',
      },
      {
        name: 'Export streamer',
        what: 'Reads in batches and writes to the response, flushing as it goes, for CSV and XLSX.',
      },
    ],
    principles: [
      {
        label: 'Generate the importer with the resource',
        text: 'a hand-written importer per resource diverges from the model it imports into. Generating both from the same definition keeps them in step.',
      },
      {
        label: 'Share the code between generate and upgrade',
        text: 'the resolver and the gate are written by functions used by both the generator and the upgrade repair, so a new project and an upgraded one cannot differ.',
      },
      {
        label: 'Never read an error as absence',
        text: 'this is the subtle one. A failed lookup treated as "not found" creates duplicates, and nothing in the run reports a problem.',
      },
      {
        label: 'Continue and report',
        text: 'one bad row must not end the run. The useful output is a list of row numbers, because the person fixing it is looking at the spreadsheet.',
      },
      {
        label: 'Gate what holds a connection',
        text: 'an import is a long-running database client. Letting several run at once is a pool outage that shows up as general slowness.',
      },
    ],
    patterns: [
      ['Streaming', 'memory bounded by the batch, not the file'],
      ['Batch processing', 'the batch as the unit of work and of failure'],
      ['Identity map', 'the per-run natural key cache'],
      ['Bulkhead', 'a concurrency gate protecting the shared pool'],
      ['Error collection', 'continue and report rather than stop on first'],
    ],
  },

  scaling: [
    'Import throughput is batch size times batch rate, and batch size has diminishing returns past a few hundred rows while lock duration keeps growing.',
    'The gate means import throughput does not scale with concurrent imports, deliberately. Serialised and predictable beats parallel and pool-exhausting.',
    'Name resolution is the cost that scales with distinct values rather than rows, which is why the cache turns 180,000 queries into 400.',
    'A very large import belongs in a background job rather than a request, so the browser is not what has to stay open for four minutes.',
    'Export streaming means the response size is unbounded without the memory being unbounded, which is the only way a hundred thousand row export works.',
    'XLSX cannot stream as freely as CSV because of its structure, so very large exports should prefer CSV and say so.',
  ],

  bottlenecks: {
    problems: [
      {
        label: 'Connection pool exhaustion',
        text: 'the original failure. Several concurrent imports each holding a connection for minutes present as the whole application being slow.',
      },
      {
        label: 'Duplicate creation from misread errors',
        text: 'a lookup error treated as a miss creates a second "Suppliers" group, and nothing reports anything wrong.',
      },
      {
        label: 'A run that fails halfway',
        text: 'batches are committed, so a failure at row 60,000 leaves 59,500 rows written and the user unsure what to do with the file.',
      },
      {
        label: 'Spreadsheet data quality',
        text: 'dates in three formats, numbers with thousands separators, trailing spaces, and a header row that does not match the columns.',
      },
      {
        label: 'Browser-held imports',
        text: 'a four minute import in a request is a four minute request, and a closed tab or a proxy timeout ends it.',
      },
    ],
    improvements: [
      {
        label: 'Keep the gate, report the wait',
        text: 'serialising is correct. Telling the user they are queued is what makes it acceptable.',
      },
      {
        label: 'Distinguish every lookup outcome',
        text: 'found, not found, and failed are three results. Collapsing the last two is how an import corrupts data quietly.',
      },
      {
        label: 'Make imports idempotent on a natural key',
        text: 'upsert on the key rather than insert, so rerunning the same file after a half-failure converges instead of duplicating.',
      },
      {
        label: 'Validate the whole file first',
        text: 'a dry run reporting every bad row before anything is written. The user fixes the spreadsheet once rather than discovering problems in stages.',
      },
      {
        label: 'Run large imports as a job',
        text: 'upload, queue, and let the user close the tab. Progress comes from polling the run, not from a held request.',
      },
    ],
  },

  seeAlso: [
    { title: 'Import and export reference', href: '/docs/features/import-export' },
    { title: 'Background jobs', href: '/docs/systems/background-jobs' },
  ],
}
