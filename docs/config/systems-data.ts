import type { SystemDesign } from './systems-types'

export const CACHING: SystemDesign = {
  slug: 'caching',
  name: 'Caching',
  tagline:
    'Serving a repeated answer without asking again, while keeping a cache failure a slowdown rather than an outage.',
  group: 'Data',
  packages: ['internal/cache', 'internal/middleware'],

  problem: {
    text: [
      'Most read traffic asks the same questions. A dashboard count, a settings lookup, a permission set, a list that changes hourly and is read thousands of times an hour. Recomputing each of them from the database is work the database did not need to do.',
      'A cache is easy to add and easy to get wrong in four specific ways. It can serve a stale answer after a write. It can become a dependency, so that losing it takes the application down rather than slowing it. It can stampede, where a hot key expires and every concurrent request rebuilds it at once. And a deploy can write ten thousand keys in the same second, which then all expire in the same second and hand the database the whole working set back at once.',
      'Each of those has a known remedy, and the value of a cache layer in a framework is that the remedies are already in it rather than being rediscovered per project.',
    ],
    capabilities: [
      'Read through: ask for a value, get it from the cache or compute and store it, in one call.',
      'Fall through to the loader when the cache is unavailable, so an outage is a slowdown.',
      'Spread expiry times so a synchronised write does not become a synchronised miss.',
      'Let one caller rebuild a hot key while the others wait briefly rather than all rebuilding.',
      'Invalidate by key and by prefix when the underlying data changes.',
      'Run with no cache at all, because a small deployment should not need Redis.',
    ],
  },

  functional: [
    'A Remember call taking a key, a time to live, a destination and a loader.',
    'Get, Set, Delete and DeleteByPrefix for the cases that are not read-through.',
    'A response cache for whole GET responses on routes that opt in.',
    'Automatic invalidation on write for generated resources, by prefix.',
    'Jitter of up to ten per cent added to every time to live.',
    'Single-flight rebuilding of an expired key.',
    'Complete absence as a supported configuration: no Redis URL means no cache and no errors.',
  ],

  nonFunctional: [
    {
      label: 'Never load-bearing',
      text: 'every read path works with the cache removed. That is what makes it acceptable for a cache to be a network dependency at all.',
    },
    {
      label: 'Correct before fast',
      text: 'a write invalidates before it returns. Serving a value the caller just changed is worse than not caching.',
    },
    {
      label: 'Keyed completely',
      text: 'a key includes everything that varies the answer, including the owner and the tenant. A key that omits them is a data leak, not a performance bug.',
    },
    {
      label: 'Bounded',
      text: 'everything has a time to live. A cache with no expiry is a second database with no migrations.',
    },
    {
      label: 'Observable',
      text: 'hit rate and eviction are visible. A cache nobody measures is a cache nobody can tune.',
    },
  ],

  capacity: {
    assumptions: [
      ['Read requests per second', '926 average, 4,630 peak'],
      ['Cacheable share of reads', '~60%'],
      ['Target hit rate', '85%'],
      ['Average cached value', '4 KB'],
      ['Distinct hot keys', '~200,000'],
      ['Default time to live', '5 minutes'],
    ],
    estimates: [
      {
        label: 'Database load avoided',
        working: [
          '926 reads/s x 60% cacheable = 556 cacheable reads/s',
          'at 85% hit rate = 473 reads/s served from cache',
          'database sees 453 reads/s instead of 926',
        ],
        note: 'Roughly half the read load. The number that matters is not the hit rate on its own but the hit rate times the cacheable share.',
      },
      {
        label: 'Memory',
        working: [
          '200,000 keys x 4 KB = 800 MB',
          'plus Redis overhead of ~100 bytes/key = ~820 MB',
          'a 1 GB instance holds the working set with room to spare',
        ],
      },
      {
        label: 'Stampede without protection',
        working: [
          'a hot key read 200 times/second, time to live expires',
          'without single-flight: 200 concurrent rebuilds of the same query',
          'with it: 1 rebuild, 199 short waits',
        ],
        note: 'The failure is worst exactly when the key is most valuable, which is why this is not an optimisation but a correctness feature of the cache layer.',
      },
      {
        label: 'Synchronised expiry',
        working: [
          'a deploy warms 10,000 keys in one second with a 5 minute life',
          'without jitter: 10,000 misses in the same second, five minutes later',
          'with up to 10% jitter: spread over 30 seconds',
        ],
      },
    ],
  },

  highLevel: {
    intro:
      'One read-through call covers most uses. The response cache and the invalidation hooks are layered on top of the same store.',
    components: [
      {
        label: 'Cache client',
        text: 'the store interface. Backed by Redis when configured and by nothing when not, so the absence case needs no branches at the call sites.',
      },
      {
        label: 'Remember',
        text: 'the read-through call. Cache-aside with the three things a hand-written version usually misses: fall-through on failure, jittered expiry and single-flight rebuilds.',
      },
      {
        label: 'Response cache',
        text: 'middleware that stores a whole GET response for routes that opt in, keyed by path, query and identity.',
      },
      {
        label: 'Invalidation hooks',
        text: 'generated resources delete their prefix on create, update and delete, so a write cannot leave a stale list behind.',
      },
      {
        label: 'Key builder',
        text: 'the convention that makes keys include the owner and the tenant, so no key can be shared across a boundary.',
      },
    ],
    flow: {
      title: 'A read-through cache with a write invalidating it',
      nodes: [
        { id: 'client', label: 'Client', col: 0, row: 0 },
        { id: 'handler', label: 'Handler', col: 1, row: 0 },
        { id: 'remember', label: 'Remember', col: 2, row: 0, accent: true },
        { id: 'redis', label: 'Cache Store', col: 3, row: 0 },
        { id: 'flight', label: 'Single Flight', col: 2, row: 1, accent: true },
        { id: 'loader', label: 'Loader', col: 3, row: 1 },
        { id: 'db', label: 'Database', col: 3, row: 2 },
        { id: 'invalidate', label: 'Invalidate Prefix', col: 1, row: 2, accent: true },
        { id: 'write', label: 'Write Path', col: 0, row: 2 },
      ],
      edges: [
        { from: 'client', to: 'handler', step: 1 },
        { from: 'handler', to: 'remember', step: 2 },
        { from: 'remember', to: 'redis', step: 3 },
        { from: 'remember', to: 'flight', step: 4 },
        { from: 'flight', to: 'loader', step: 5 },
        { from: 'loader', to: 'db', step: 6 },
        { from: 'write', to: 'invalidate', step: 7 },
        { from: 'invalidate', to: 'redis', step: 8, bend: 'v' },
      ],
      steps: [
        'A read request arrives on a route whose answer is worth caching.',
        'The handler calls Remember with a key, a lifetime, a destination and a function that can compute the value.',
        'The cache is asked. A hit returns immediately. A cache that is unreachable is treated as a miss rather than an error, which is what keeps it from being load-bearing.',
        'On a miss, single-flight ensures one caller runs the loader while the others wait for its result instead of all running the same query.',
        'The loader computes the value, usually a query the handler would have run anyway.',
        'The result is stored with the requested lifetime plus up to ten per cent of jitter, so keys written together do not expire together.',
        'Separately, a write to the underlying resource deletes the cache prefix for it.',
        'Invalidation happens before the write returns, so a client that reads immediately after writing cannot see the value it just replaced.',
      ],
    },
    dataFlow: [
      'Keys carry everything that varies the answer. For an owned resource that includes the owner, and in a multitenant deployment the organisation.',
      'Nothing is cached without a lifetime. An entry that never expires is state that no migration will ever fix.',
      'Invalidation is by prefix rather than by enumerating keys, because the set of list keys for a resource is unbounded in query parameters.',
      'A cache miss and a cache outage take the same path, so there is no error handling at the call site to get wrong.',
    ],
  },

  stack: [
    ['Store', 'Redis, optional'],
    ['Pattern', 'cache-aside, read-through'],
    ['Default lifetime', '5 minutes, jittered by up to 10%'],
    ['Stampede control', 'single-flight per key'],
    ['Invalidation', 'by key and by prefix, on write'],
    ['Absence', 'no Redis URL disables cache, jobs and cron, with no errors'],
  ],

  dataModel: {
    intro:
      'There is no schema. The design question is the key, and the convention is what keeps it correct.',
    entities: [
      {
        name: 'key convention',
        fields: [
          ['resource:id', 'a single record'],
          ['resource:list:<hash of query>', 'a list, hashed because the parameter space is unbounded'],
          ['user:<id>:...', 'anything that varies per user'],
          ['org:<id>:...', 'anything that varies per tenant'],
          ['perm:role:<id>', 'resolved permission sets, invalidated when a role changes'],
        ],
        note: 'Everything that changes the answer appears in the key. The commonest cache bug in any system is a key that omits the owner.',
      },
    ],
  },

  lowLevel: {
    classes: [
      {
        name: 'cache.Client',
        file: 'internal/cache/cache.go',
        what: 'The store. Get, Set, Delete, DeleteByPrefix. Every method tolerates the backend being absent or unreachable.',
        methods: ['Get', 'Set', 'Delete', 'DeleteByPrefix'],
      },
      {
        name: 'cache.Remember',
        file: 'internal/cache/remember.go',
        what: 'The call that should be used by default. Generic over the value type, so the destination is typed rather than an interface the caller asserts.',
      },
      {
        name: 'Response cache middleware',
        what: 'Caches whole GET responses for opted-in routes. Opt-in rather than automatic, because a response that varies by identity is easy to cache wrongly.',
      },
    ],
    principles: [
      {
        label: 'Make the correct call the shortest one',
        text: 'Remember is one line and handles the four classic mistakes. Get and Set are available and longer to write, which is the right incentive.',
      },
      {
        label: 'A dependency that may be absent',
        text: 'the client treats unreachable as empty. There is no error to handle at the call site and therefore no error handling to get wrong.',
      },
      {
        label: 'Invalidate at the write, not on a timer',
        text: 'relying on a short lifetime to hide staleness means choosing between stale data and no caching. Invalidating removes the choice.',
      },
    ],
    patterns: [
      ['Cache aside', 'read through, write around, invalidate on change'],
      ['Single flight', 'one rebuild per key, the rest wait'],
      ['Jittered expiry', 'spreading synchronised writes so they do not become synchronised misses'],
      ['Null object', 'the no-op cache when none is configured'],
    ],
  },

  scaling: [
    'A cache hit removes a database round trip, which is the cheapest capacity available until the hit rate stops improving.',
    'Hit rate times cacheable share is the real number. A 99% hit rate on 5% of traffic is worth less than 80% on 60%.',
    'Memory grows with the working set, not with total data. Sizing follows from distinct hot keys, not from table sizes.',
    'Redis is single-threaded per instance. Very hot single keys are the limit, which single-flight helps with on the application side.',
    'Prefix invalidation is a scan on some backends. Keeping prefixes narrow keeps that cost bounded.',
    'With no cache configured the application is slower and entirely correct, which is the right default for a small deployment.',
  ],

  bottlenecks: {
    problems: [
      {
        label: 'Keys that omit the owner',
        text: 'the worst cache bug there is, because it does not look like a bug. It looks like one user occasionally seeing another user’s data.',
      },
      {
        label: 'The cache becoming load-bearing',
        text: 'a path that errors when the cache is unreachable has turned an optional component into a required one, usually without anybody deciding to.',
      },
      {
        label: 'Invalidation gaps',
        text: 'a write path that does not go through the service does not invalidate, and the stale value lives for a full lifetime.',
      },
      {
        label: 'Unbounded key growth',
        text: 'list keys hashed from arbitrary query parameters can be generated without limit by a caller who varies them.',
      },
    ],
    improvements: [
      {
        label: 'Build keys, do not write them',
        text: 'a helper that takes the context and appends the owner and tenant makes the complete key the easy one.',
      },
      {
        label: 'Test with the cache off',
        text: 'running the suite with no cache configured proves no path depends on it, which is the property that is otherwise assumed.',
      },
      {
        label: 'Invalidate in the service',
        text: 'putting the hook where every writer passes, rather than in the handler, closes the gap for jobs and commands.',
      },
      {
        label: 'Whitelist the query parameters',
        text: 'hashing only the parameters that actually vary the answer bounds the key space regardless of what a caller sends.',
      },
    ],
  },

  seeAlso: [
    { title: 'Caching reference', href: '/docs/batteries/caching' },
    { title: 'Background jobs', href: '/docs/systems/background-jobs' },
  ],
}

export const CONCURRENCY: SystemDesign = {
  slug: 'concurrency-control',
  name: 'Concurrency Control',
  tagline:
    'Two people editing the same record, and the second one finding out rather than silently winning.',
  group: 'Data',
  packages: ['internal/concurrency', 'internal/services'],

  problem: {
    text: [
      'Two operators open the same invoice. One changes the due date, the other changes the amount. Both save. The second write overwrites the first, and nobody is told. The first operator’s change is simply gone, and they will not find out until they look again, if they ever do.',
      'This is the lost update problem, and it is invisible in testing because it needs two people and a few seconds. In production it happens constantly, and it is reported as "the system lost my change", which is indistinguishable from a bug in saving.',
      'The two families of solution are pessimistic locking, where the first reader takes a lock and the second waits, and optimistic locking, where both proceed and the second write is rejected if the record moved. For a web application the second is almost always right: holding a database lock across a user thinking about a form is not viable.',
    ],
    capabilities: [
      'Give every record a version that changes on every write.',
      'Let a client declare which version it read, and reject a write against a different one.',
      'Tell the rejected client which version the record is at now, so it can show the difference.',
      'Leave clients that do not participate working exactly as before.',
      'Apply to the generated update and patch paths without each one implementing it.',
      'Carry the same guarantee into offline clients that sync later.',
    ],
  },

  functional: [
    'A version column on every generated model, incremented by the database on update.',
    'The version returned in responses and as an ETag header.',
    'An If-Match header carrying the version the client read.',
    'A 409 with the current version when the record has moved.',
    'Unconditional writes still accepted, for clients that do not opt in.',
    'The same check in the offline sync path, where the gap between read and write is hours.',
  ],

  nonFunctional: [
    {
      label: 'No held locks',
      text: 'nothing is locked while a human is deciding. A form open overnight costs nothing.',
    },
    {
      label: 'Atomic',
      text: 'the version check is part of the update statement, not a read followed by a write, or the race simply moves.',
    },
    {
      label: 'Opt-in',
      text: 'a client that sends no If-Match behaves as it always did. Making it mandatory would break every existing caller.',
    },
    {
      label: 'Actionable',
      text: 'a conflict response carries the current version and the current row, so the interface can show what changed rather than just refusing.',
    },
    {
      label: 'Uniform',
      text: 'the same mechanism on every resource, because a guarantee that applies to some records is one nobody can rely on.',
    },
  ],

  capacity: {
    assumptions: [
      ['Writes per second', '~90 at peak across all resources'],
      ['Share of writes using If-Match', 'all admin edits, most API clients'],
      ['Conflict rate on contended records', '0.1% to 5%, depending on the workload'],
      ['Version column', '4 bytes'],
    ],
    estimates: [
      {
        label: 'Cost of the check',
        working: [
          'the predicate is added to the existing UPDATE',
          'rows affected is already returned',
          'additional cost: zero queries, one integer comparison',
        ],
        note: 'Optimistic locking is free on the happy path. That is the entire argument for it over anything that holds a lock.',
      },
      {
        label: 'Conflict handling',
        working: [
          '90 writes/s x 1% conflict = ~1 conflict/second',
          'each costs one extra read to report the current version',
        ],
        note: 'The extra read happens only on the rejected path, which is rare by construction.',
      },
      {
        label: 'Storage',
        working: ['4 bytes per row x 10,000,000 rows = 40 MB'],
      },
    ],
  },

  highLevel: {
    intro:
      'The version lives on the row, travels in the response, comes back in a header, and becomes a predicate on the update.',
    components: [
      {
        label: 'Version column',
        text: 'an integer on every generated model, raised by the database as part of the update rather than by the application, so two concurrent writers cannot both read the same number.',
      },
      {
        label: 'Precondition parser',
        text: 'reads If-Match and turns it into a predicate, or into nothing when the header is absent.',
      },
      {
        label: 'Conditional update',
        text: 'the write, with the version predicate attached. Rows affected is the answer.',
      },
      {
        label: 'Conflict reporter',
        text: 'when nothing was updated, reads the row to find out which version it is at and returns both that and the row.',
      },
      {
        label: 'ETag support',
        text: 'the version as an HTTP entity tag, so standard clients and caches participate without bespoke fields.',
      },
    ],
    flow: {
      title: 'Two writers, one record',
      nodes: [
        { id: 'a', label: 'Client A', col: 0, row: 0 },
        { id: 'handler', label: 'Update Handler', col: 1, row: 0 },
        { id: 'pre', label: 'Precondition', col: 2, row: 0, accent: true },
        { id: 'service', label: 'Service', col: 3, row: 0 },
        { id: 'b', label: 'Client B', col: 0, row: 1 },
        { id: 'db', label: 'Conditional UPDATE', col: 2, row: 1, accent: true },
        { id: 'conflict', label: '409 + Current Version', col: 3, row: 1 },
        { id: 'ok', label: '200 + New Version', col: 2, row: 2 },
      ],
      edges: [
        { from: 'a', to: 'handler', step: 1 },
        { from: 'handler', to: 'pre', step: 2 },
        { from: 'pre', to: 'service', step: 3 },
        { from: 'service', to: 'db', step: 4, bend: 'v' },
        { from: 'db', to: 'ok', step: 5 },
        { from: 'b', to: 'db', step: 6, bend: 'h' },
        { from: 'db', to: 'conflict', step: 7, dashed: true },
      ],
      steps: [
        'Client A reads the invoice. The response carries version 7, both in the body and as an ETag.',
        'Client A submits a change with If-Match naming version 7. Client B, who read the same row, is still editing.',
        'The precondition is parsed into a predicate. A request without the header produces no predicate and proceeds unconditionally.',
        'The service issues an UPDATE with WHERE id = ? AND version = 7. The version is raised by the database as part of the same statement.',
        'One row is affected. Client A gets 200 and the new version, 8.',
        'Client B now submits, still naming version 7. The same statement runs.',
        'No rows match, because the row is at version 8. Rather than reporting a generic failure, the service reads the row and returns 409 with the current version and the current values, so the interface can show what changed and offer to merge.',
      ],
    },
    dataFlow: [
      'The version is raised in the UPDATE statement itself. An application that read, incremented and wrote would have reintroduced the race it is trying to remove.',
      'A missing If-Match means no predicate. The write succeeds, which keeps older clients working.',
      'The conflict response includes the current row, because a client that only learns it failed can do nothing useful except reload and lose the user’s typing.',
      'Offline clients carry the version they last saw, so a sync that happens hours later is checked the same way.',
    ],
  },

  stack: [
    ['Strategy', 'optimistic concurrency control'],
    ['Version', 'an integer column, raised by the database'],
    ['Transport', 'If-Match and ETag, plus the version in the body'],
    ['Conflict status', '409, with the current version and row'],
    ['Scope', 'every generated resource, and the sync path'],
  ],

  dataModel: {
    entities: [
      {
        name: 'any generated resource',
        fields: [
          ['version', 'integer, not null, default 1, raised on every update'],
        ],
        note: 'Raised by a BeforeUpdate hook that sets it to version + 1 as a SQL expression rather than a value read in Go, which is what makes it atomic.',
      },
    ],
  },

  api: {
    groups: [
      {
        title: 'Conditional writes',
        rows: [
          { method: 'GET', path: '/api/v1/invoices/:id', what: 'Returns the version and an ETag' },
          { method: 'PUT', path: '/api/v1/invoices/:id', what: 'Honours If-Match; 409 when the version moved' },
          { method: 'PATCH', path: '/api/v1/invoices/:id', what: 'The same, for partial writes' },
        ],
      },
    ],
    samples: [
      {
        title: 'A conditional write',
        language: 'http',
        code: `PUT /api/v1/invoices/01a11e82-d0bc-7367-9c39-5e9873953dde
If-Match: "7"
Content-Type: application/json

{ "amount": 1250, "due_on": "2026-11-01" }`,
      },
      {
        title: 'A conflict that the interface can act on',
        language: 'json',
        code: `{
  "error": {
    "code": "WRITE_CONFLICT",
    "message": "This record changed since you loaded it",
    "details": {
      "your_version": 7,
      "current_version": 8
    }
  },
  "data": {
    "id": "01a11e82-d0bc-7367-9c39-5e9873953dde",
    "amount": 1100,
    "due_on": "2026-10-28",
    "version": 8
  }
}`,
      },
    ],
  },

  lowLevel: {
    classes: [
      {
        name: 'concurrency.Precondition',
        file: 'internal/concurrency/concurrency.go',
        what: 'The parsed If-Match. Carries a scope that adds the predicate and a Missed method that interprets rows-affected.',
        methods: ['Scope', 'Missed'],
      },
      {
        name: 'concurrency.Tag',
        what: 'The ETag for a version. One function, so the format in the response and the format accepted in the request cannot drift.',
      },
      {
        name: 'ErrConflict',
        what: 'Carries the version the row is actually at, so the handler can answer with something more useful than a refusal.',
      },
      {
        name: 'BeforeUpdate',
        what: 'The GORM hook that raises the version as a SQL expression. Written as an expression rather than a Go value on purpose.',
      },
    ],
    principles: [
      {
        label: 'Check and write in one statement',
        text: 'a read-then-write check has the same race it is meant to prevent, just narrower. The predicate belongs in the UPDATE.',
      },
      {
        label: 'Optional by default',
        text: 'the feature is available to clients that want it and invisible to those that do not, so it could be added without a breaking change.',
      },
      {
        label: 'Report the state, not just the failure',
        text: 'a 409 carrying the current row lets the interface offer a merge. One carrying only a code forces a reload and loses the user’s work.',
      },
    ],
    patterns: [
      ['Optimistic concurrency', 'version checked in the write predicate'],
      ['Compare and swap', 'the database-level equivalent, as one statement'],
      ['HTTP preconditions', 'If-Match and ETag rather than a bespoke field'],
    ],
  },

  scaling: [
    'There is no additional query on the happy path and no lock held across a request, so this costs nothing at any write rate.',
    'Conflicts cost one extra read, on a path that is rare by construction.',
    'Contention is a property of the workload, not the mechanism. A single row written by many clients will conflict often, and the answer there is a different data model rather than a different locking strategy.',
    'The approach works unchanged across replicas and across an offline client syncing hours later, because the version lives on the row rather than in a session.',
    'Four bytes a row is immaterial at any table size.',
  ],

  bottlenecks: {
    problems: [
      {
        label: 'Hot rows',
        text: 'a counter row that everybody updates conflicts constantly, and optimistic locking turns that into a retry storm rather than a queue.',
      },
      {
        label: 'Clients that ignore the conflict',
        text: 'a client that retries a 409 by reloading and resubmitting without showing the user what changed has reimplemented last-write-wins with extra steps.',
      },
      {
        label: 'Partial writes across records',
        text: 'the version protects one row. An operation spanning three rows can still half-succeed unless it is in a transaction.',
      },
      {
        label: 'Version not surfaced',
        text: 'if a client never sees the version it cannot send If-Match, and the protection is present but unused.',
      },
    ],
    improvements: [
      {
        label: 'Model hot values differently',
        text: 'a counter belongs in an append-and-aggregate shape or an atomic increment, not in a row that every writer must version-check.',
      },
      {
        label: 'Show the difference',
        text: 'the conflict response carries the current row precisely so the interface can present both and let the user choose.',
      },
      {
        label: 'Wrap multi-row operations',
        text: 'a transaction around the whole operation makes the group atomic, with the per-row versions still preventing lost updates inside it.',
      },
      {
        label: 'Return the version everywhere',
        text: 'including it in list responses as well as reads means a client editing from a list has it without a second request.',
      },
    ],
  },

  seeAlso: [
    { title: 'Offline sync', href: '/docs/concepts/offline-sync' },
    { title: 'API contract', href: '/docs/systems/api-contract' },
  ],
}
