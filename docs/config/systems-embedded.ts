import type { SystemDesign } from './systems-types'

export const GORM_STUDIO: SystemDesign = {
  slug: 'gorm-studio',
  name: 'GORM Studio',
  tagline:
    'A database browser that ships inside the binary, reads the schema from both the database and the models, and cannot be talked into dropping a table.',
  group: 'Embedded',
  packages: ['github.com/MUKE-coder/gorm-studio/studio', 'v1.1.1'],

  interview: {
    intro:
      '"Design an internal admin tool" is a real prompt and a surprisingly good one, because the interesting content is entirely about constraints rather than features. These are the questions it raises.',
    questions: [
      {
        q: 'Why build a database browser into the application rather than using an existing client?',
        a: 'Because of where it is needed. A staging container has no desktop, no installed client, and no database port open to anybody laptop. Customer-hosted instances cannot have one opened either. In exactly the situations where looking at the data would settle an argument in thirty seconds, the usual tools are the ones that are not there.',
        see: 'problem',
      },
      {
        q: 'That sounds like a database console on the public internet.',
        a: 'It is, which is why nearly all of the design is constraints. Authentication in front, read-only and no SQL editor unless a specific task needs otherwise, a per-table policy, and loud warnings at startup for each dangerous combination. The library cannot refuse to mount without breaking the laptop case it is good at, so the unsafe default is made noisy rather than impossible.',
        see: 'requirements',
      },
      {
        q: 'How do you stop somebody dropping a table through the SQL box?',
        a: 'Refuse categories rather than sanitising strings, because there is no safe escaping of DROP TABLE. A blocked leading keyword list applies whether or not the instance is read-only, and it includes the non-obvious ones: VACUUM is on it because VACUUM INTO writes an arbitrary file, which is data exfiltration wearing a maintenance command.',
        see: 'high-level-design',
      },
      {
        q: 'Checking the first keyword is easy to get around. "SELECT 1; DELETE FROM users".',
        a: 'Which is why the statement count is checked first. Comments are stripped, statements are split, and anything that is not exactly one statement is rejected. That single rule closes the whole class, where keyword matching on the first word never would. A common table expression hiding DML inside it needs whole-word matching to catch as well.',
        see: 'high-level-design',
      },
      {
        q: 'Some tables hold secrets. How do you handle those?',
        a: 'Hide them, and make hidden mean hidden in every path: absent from the schema, 404 on a direct request, left out of exports, and refused by the SQL editor when a statement names one. A table omitted from three of the four and present in the fourth is worse than one never hidden, because somebody is relying on it.',
        see: 'requirements',
      },
      {
        q: 'You have a scope function limiting what an operator sees. Does the SQL editor respect it?',
        a: 'No, and saying so is the point. A scope filters table queries and cannot filter arbitrary SQL, so leaving both enabled gives an operator a way around the boundary. Documenting that limitation, and advising that the editor be disabled when a scope is set, is the difference between a known constraint and a breach.',
        see: 'low-level-design',
      },
      {
        q: 'What stops an import from taking the process down?',
        a: 'Three caps applied together: a byte limit checked before the body is read, a row limit, and a timeout. The ordering matters most. A size check performed after buffering the upload has already lost, because the decompression bomb has gone off by then.',
        see: 'capacity',
      },
      {
        q: 'What is the operational risk once it is running?',
        a: 'That it shares the application connection pool. A sort on an unindexed column over ten million rows is a long query holding a connection a request wanted, and it is being run by somebody investigating a problem, which is when the system is already under strain. A separate pool with a low maximum and a statement timeout contains it.',
        see: 'bottlenecks',
      },
      {
        q: 'Someone edits a row here. What is different from editing it through the application?',
        a: 'Every validation, hook and audit the application would have applied is skipped, so the row can end up in a state no code path could produce. The audit callback records that it happened, not that it was correct. That is why writes through this tool belong in the same audit log as everything else, and why read-only is the right default.',
        see: 'bottlenecks',
      },
      {
        q: 'Why read the schema from both the database and the model structs?',
        a: 'Because neither alone describes the schema a developer is thinking about. The database knows the column types and the foreign keys; the models know the relationship names, the Go types and the many-to-many join tables. Reconciling the two is what lets the tool offer navigation along a relationship rather than just a list of columns.',
        see: 'high-level-design',
      },
    ],
  },

  problem: {
    text: [
      'Every developer needs to look at the data. Did the migration run, why does that row have a null in it, what did the importer actually write. The usual answers are a separate desktop client that each person installs and configures, or psql, which does not exist in a container, or an admin panel built for the business rather than for the rows.',
      'None of them works where it is most needed. A staging deployment has no desktop on it. A container has no client installed. A customer-hosted instance cannot have a database port opened to a laptop. In exactly the situations where looking at the data would settle an argument in thirty seconds, the usual tools are the ones not there.',
      'A browser mounted inside the application solves that, and immediately creates a worse problem than the one it solved. It is a route that can read every table, and if it runs SQL it is a route that can drop them. Mounted without thought it is a database console on the public internet. The interesting part of this design is not the grid, it is everything that constrains it.',
      'There is also a schema problem that is specific to Go. The database knows the column types and the foreign keys. The GORM models know the relationship names, the Go types, and the many-to-many join tables. Neither alone describes the schema a developer is thinking about, so the tool has to read both and reconcile them.',
    ],
    capabilities: [
      'Discover the schema from the database and from the model structs, and merge them.',
      'Browse, filter, sort and page any table without writing SQL.',
      'Follow a relationship from a row to the rows on the other side of it.',
      'Create, edit and delete records, when that has been allowed.',
      'Run a single read query, with the statements that cannot be made safe refused outright.',
      'Hide tables that must never be exposed, and make others readable but not writable.',
      'Export the schema and the data, and import both back.',
      'Generate Go model structs from an existing database.',
      'Run with no extra service, no installed client and no open database port.',
    ],
  },

  functional: [
    'Mount on a Gin router with a database handle and a list of models.',
    'Schema introspection, plus reflection over the registered model structs.',
    'A paginated grid with sorting, search and column filters.',
    'Relationship navigation for has one, has many, belongs to and many to many.',
    'Record create, edit and delete, and bulk delete.',
    'A raw SQL editor that accepts exactly one statement.',
    'A blocked keyword list that applies whether or not the instance is read-only.',
    'Per-table policy: hidden tables and read-only tables.',
    'An optional scope function, applied to every table query.',
    'An audit callback for every mutating action.',
    'Export to SQL, JSON, YAML, DBML, CSV, and an ERD as PNG or PDF.',
    'Import from SQL, JSON, CSV, Excel and Go struct files, bounded in size, rows and time.',
    'Per-IP rate limiting on the SQL and import endpoints.',
  ],

  nonFunctional: [
    {
      label: 'Safe by refusal, not by intention',
      text: 'DROP, ALTER, TRUNCATE, CREATE, GRANT, ATTACH and the rest are refused by keyword before anything is parsed. A tool that relies on the operator not typing the dangerous thing is not a safe tool.',
    },
    {
      label: 'One statement, always',
      text: 'the editor splits and counts. Two statements is an error, which is what closes the "harmless select, semicolon, delete" shape.',
    },
    {
      label: 'Warns when it is dangerous',
      text: 'mounting with no authentication, with writes enabled and no authentication, or with the SQL editor and no authentication each log a warning at startup. The default is usable; the unsafe default is loud.',
    },
    {
      label: 'Bounded imports',
      text: 'a size cap of 32 MiB, a row cap of 100,000 and a 30 second timeout, all applied before the body is read. An import endpoint without them is a decompression bomb away from an outage.',
    },
    {
      label: 'Honest about what it cannot scope',
      text: 'a scope function filters table queries and cannot filter raw SQL. The documentation says to disable the editor when a scope is set rather than implying the boundary holds.',
    },
    {
      label: 'No CGo, no extra service',
      text: 'a pure Go SQLite driver and an embedded frontend. One binary, which is the only reason it is present on the staging box at all.',
    },
  ],

  capacity: {
    assumptions: [
      ['Tables in a mature project', '25 to 60'],
      ['Rows in the largest table', 'up to tens of millions'],
      ['Grid page size', '25 to 100 rows'],
      ['Concurrent operators', '1 to 5'],
      ['Import cap', '32 MiB, 100,000 rows, 30 seconds'],
    ],
    estimates: [
      {
        label: 'Schema discovery',
        working: [
          '60 tables introspected, plus reflection over the model list',
          'once per page load, cacheable',
          'tens of milliseconds',
        ],
      },
      {
        label: 'A grid page',
        working: [
          'LIMIT 50 on an indexed sort: single-digit milliseconds',
          'the row count for the pager: a COUNT over the filter',
          'on 10,000,000 rows the count dominates, by a lot',
        ],
        note: 'The same arithmetic as any list endpoint. The difference is that a browser invites someone to sort by an unindexed column on the biggest table in the system, which no application screen would ever do.',
      },
      {
        label: 'An unbounded export',
        working: [
          '10,000,000 rows x ~1 KB = ~10 GB',
          'streamed, so memory stays flat',
          'but it is 10 GB over one HTTP connection',
        ],
        note: 'Streaming is what makes this survivable rather than fast. It is still a lot of database read on a connection pool the application is sharing.',
      },
      {
        label: 'Import worst case',
        working: [
          '32 MiB compressed upload',
          'capped at 100,000 rows and 30 seconds',
          'rejected before the body is read when over the byte cap',
        ],
        note: 'Rejecting on Content-Length rather than after buffering is the part that matters. Checking after reading means the bomb already went off.',
      },
    ],
  },

  highLevel: {
    intro:
      'One mount call. A schema reader that merges two sources, a query layer that every read passes through, and a set of guards that no path can go around.',
    components: [
      {
        label: 'Mount',
        text: 'takes the router, the database handle, the model list and a config. Registers the routes under a prefix and logs a warning for each dangerous combination it was given.',
      },
      {
        label: 'Schema reader',
        text: 'introspects the database for columns, types and keys, and reflects over the model structs for relationship names, Go types and join tables. Neither source describes the whole schema.',
      },
      {
        label: 'Table policy',
        text: 'hidden tables are omitted from the schema, 404 on direct request, left out of exports, and refused by the SQL editor when a statement names one. Read-only tables can be browsed and never mutated.',
      },
      {
        label: 'Scope function',
        text: 'applied to every table query, for an operator who should see only part of the data. It cannot apply to raw SQL, and that is stated rather than glossed.',
      },
      {
        label: 'SQL classifier',
        text: 'strips comments, splits statements, refuses anything that is not exactly one, checks the leading keyword against the blocked list, and decides read or write.',
      },
      {
        label: 'Import pipeline',
        text: 'size cap, row cap and timeout, with a parser per format. The caps come first.',
      },
      {
        label: 'Audit callback',
        text: 'every mutating action: a row written, a bulk delete, a raw SQL write, an import. The record of what was done through the back door.',
      },
      {
        label: 'Embedded frontend',
        text: 'the whole UI compiled into the binary, served with a content security policy that forbids framing and pins script origins.',
      },
    ],
    flow: {
      title: 'A query typed into the SQL editor',
      nodes: [
        { id: 'operator', label: 'Operator', col: 0, row: 0 },
        { id: 'auth', label: 'Auth Middleware', col: 1, row: 0 },
        { id: 'limit', label: 'Rate Limit', col: 2, row: 0 },
        { id: 'classify', label: 'Classify SQL', col: 3, row: 0, accent: true },
        { id: 'blocked', label: 'Blocked Keywords', col: 3, row: 1, accent: true },
        { id: 'policy', label: 'Table Policy', col: 2, row: 1, accent: true },
        { id: 'refuse', label: '403 Refused', col: 1, row: 1 },
        { id: 'run', label: 'Single Statement', col: 2, row: 2, accent: true },
        { id: 'audit', label: 'Audit Callback', col: 3, row: 2 },
      ],
      edges: [
        { from: 'operator', to: 'auth', step: 1 },
        { from: 'auth', to: 'limit', step: 2 },
        { from: 'limit', to: 'classify', step: 3 },
        { from: 'classify', to: 'blocked', step: 4 },
        { from: 'blocked', to: 'policy', step: 5 },
        { from: 'policy', to: 'refuse', step: 6, dashed: true },
        { from: 'policy', to: 'run', step: 7 },
        { from: 'run', to: 'audit', step: 8 },
      ],
      steps: [
        'An operator opens the studio on a staging deployment and types a query. There is no database client on that machine and no port open to it, which is the whole reason this exists.',
        'The auth middleware runs first, if one was configured. If none was, the startup log already said so three times, because a database browser with no authentication is the failure this design can warn about and cannot prevent.',
        'The SQL and import endpoints are rate limited per client address, separately from anything the host application does. They are the two that are expensive enough to be worth limiting on their own.',
        'The statement is classified: comments stripped, statements split, and anything that is not exactly one statement rejected. That single check closes the "innocent select, semicolon, delete from users" shape, which no amount of keyword matching on the first word would catch.',
        'The leading keyword is checked against the blocked list, which applies whether or not the instance is read-only. VACUUM and REINDEX are on it alongside DROP and ALTER, because VACUUM INTO writes an arbitrary file and is therefore data exfiltration wearing a maintenance command.',
        'A statement naming a hidden table is refused, and so is a write to a read-only one. Hidden means hidden everywhere: absent from the schema, 404 on a direct request, omitted from exports, and refused here.',
        'What is left is one statement that is allowed to run. If a scope function is configured it does not apply here, which is why the configuration documentation says to disable this editor when a scope is set rather than letting an operator assume the boundary holds.',
        'A write is reported to the audit callback. Changes made through a database browser are exactly the changes that have no application log line, so this is the only record there will be.',
      ],
    },
    dataFlow: [
      'The schema comes from two places and neither is sufficient. The database has the columns and the foreign keys; the models have the relationship names, the Go types and the many-to-many join tables.',
      'Reads go through the same query layer, so the scope function and the table policy apply to the grid, the relationship view and the exports alike.',
      'Exports stream. Building a ten gigabyte file in memory before sending it is the obvious implementation and it is an out-of-memory kill.',
      'Imports are capped before the body is read, which is the difference between refusing a decompression bomb and absorbing one.',
    ],
  },

  stack: [
    ['Language', 'Go, no CGo'],
    ['Host', 'Gin, mounted under a prefix'],
    ['Databases', 'SQLite, PostgreSQL, MySQL'],
    ['Schema', 'driver introspection plus struct reflection'],
    ['Frontend', 'embedded in the binary, CSP pinned, framing denied'],
    ['Guards', 'blocked keywords, one statement, table policy, scope, caps'],
    ['Testing', 'race detection and parser fuzzing across all three engines'],
  ],

  api: {
    groups: [
      {
        title: 'Mounted under the configured prefix',
        rows: [
          { method: 'GET', path: '/studio', what: 'The embedded browser' },
          { method: 'GET', path: '/studio/api/schema', what: 'Tables, columns, relationships, minus hidden tables' },
          { method: 'GET', path: '/studio/api/tables/:table/rows', what: 'Paginated, sorted, filtered, scoped' },
          { method: 'POST', path: '/studio/api/tables/:table/rows', what: 'Create, unless read-only' },
          { method: 'DELETE', path: '/studio/api/tables/:table/rows/:id', what: 'Delete, audited' },
          { method: 'POST', path: '/studio/api/sql', what: 'One statement, classified and rate limited' },
          { method: 'POST', path: '/studio/api/import', what: 'Capped in bytes, rows and seconds' },
          { method: 'GET', path: '/studio/api/export', what: 'Streamed, in any supported format' },
        ],
      },
    ],
    samples: [
      {
        title: 'How Grit mounts it, and the three switches that matter',
        language: 'go',
        code: `studioCfg := studio.Config{
    Prefix:     "/studio",
    ReadOnly:   cfg.GORMStudioReadOnly,
    DisableSQL: cfg.GORMStudioDisableSQL,
}
if cfg.GORMStudioUsername != "" && cfg.GORMStudioPassword != "" {
    studioCfg.AuthMiddleware = gin.BasicAuth(gin.Accounts{
        cfg.GORMStudioUsername: cfg.GORMStudioPassword,
    })
}
studio.Mount(r, db, []interface{}{
    &models.User{}, &models.Upload{}, &models.Blog{},
}, studioCfg)`,
      },
      {
        title: 'Hiding what must never be browsed',
        language: 'go',
        code: `studio.Config{
    TablePolicy: studio.TablePolicy{
        // Absent from the schema, 404 on request, left out of exports,
        // and the SQL editor refuses any statement naming one.
        Hidden:   []string{"sessions", "payment_tokens", "api_keys"},
        ReadOnly: []string{"activity_logs", "migration_runs"},
    },
    AuditLogger: studio.DefaultAuditLogger,
}`,
      },
    ],
  },

  lowLevel: {
    classes: [
      {
        name: 'studio.Mount',
        what: 'The entry point. Registers routes and logs a warning for each dangerous combination: no auth, writes without auth, SQL without auth.',
      },
      {
        name: 'classifySQL',
        file: 'studio/handlers.go',
        what: 'Strips comments, splits statements, refuses anything that is not exactly one, then reports whether it is blocked and whether it is a read.',
      },
      {
        name: 'blockedSQLKeywords',
        what: 'Leading keywords never allowed, read-only or not. VACUUM and REINDEX are in it because VACUUM INTO writes an arbitrary file.',
      },
      {
        name: 'containsSQLWord',
        what: 'Whole-word matching, used to catch DML smuggled inside a common table expression where the leading keyword is innocent.',
      },
      {
        name: 'studio.TablePolicy',
        what: 'Hidden and read-only tables. Hidden applies to the schema, direct requests, exports and the SQL editor, because a table that is hidden in only three of the four is not hidden.',
      },
      {
        name: 'Config.Scope',
        what: 'A per-request filter on table queries. Documented as not covering raw SQL, with the advice to disable the editor when it is set.',
      },
    ],
    principles: [
      {
        label: 'Refuse categories, do not sanitise strings',
        text: 'there is no safe escaping of DROP TABLE. The blocked list refuses whole classes of statement and the single-statement rule removes the way around it.',
      },
      {
        label: 'Hidden has to mean hidden in every path',
        text: 'the schema, the row endpoint, the export and the SQL editor. A table omitted from three of them and present in the fourth is worse than one that was never hidden, because somebody is relying on it.',
      },
      {
        label: 'Say what a guard does not cover',
        text: 'the scope function cannot filter raw SQL. Documenting that is the difference between a known limitation and a breach.',
      },
      {
        label: 'Cap before reading',
        text: 'an import limit checked after buffering the upload has already lost. The byte cap is applied to the request before the body is consumed.',
      },
      {
        label: 'Warn loudly about the unsafe default',
        text: 'it mounts with no authentication because that is right for a laptop, and it says so three times at startup because it is not right anywhere else.',
      },
    ],
    patterns: [
      ['Embedded admin', 'the tool inside the process it inspects'],
      ['Allowlist and blocklist', 'categories refused, not inputs escaped'],
      ['Policy object', 'per-table visibility and mutability as data'],
      ['Reflection plus introspection', 'two schema sources reconciled'],
      ['Streaming export', 'bounded memory over an unbounded result'],
    ],
  },

  scaling: [
    'This is an operator tool, so concurrency is a handful of people and throughput never matters.',
    'What does matter is that it shares the application’s connection pool. A sort on an unindexed column over ten million rows is a long query holding a connection that a request wanted.',
    'Exports stream, so memory stays flat while the database read does not. A large export during business hours is felt by users.',
    'Schema discovery is proportional to table count and worth caching, since it runs on every page load and the schema changes at deploy time.',
    'Imports are capped in three dimensions at once, because any one of them alone leaves a way to tie up the process.',
    'The honest scaling answer is that it belongs on staging and on a locked-down production route, not as a general-purpose query tool. Nothing in the design makes it a replacement for a read replica and a proper client.',
  ],

  bottlenecks: {
    problems: [
      {
        label: 'Mounted with no authentication',
        text: 'the single catastrophic failure. The library warns three times at startup and cannot do more, because refusing to mount would break the laptop case it is good at.',
      },
      {
        label: 'Sharing the application connection pool',
        text: 'an expensive browse competes with live traffic, and it is run by someone who is investigating a problem, which is when the system is already under strain.',
      },
      {
        label: 'The SQL editor defeats the scope',
        text: 'a scope function filters table queries and cannot filter raw SQL, so leaving both on gives an operator a way around the boundary.',
      },
      {
        label: 'Count on a huge table',
        text: 'the pager needs a total, and a filtered count over tens of millions of rows is far more expensive than the page itself.',
      },
      {
        label: 'Writes with no application context',
        text: 'a row edited here skips every validation, hook and audit the application would have applied. The audit callback records that it happened, not that it was correct.',
      },
    ],
    improvements: [
      {
        label: 'Always mount behind auth, and prefer read-only',
        text: 'basic auth at minimum, the application’s own admin check where there is one, and read-only plus no SQL editor unless a specific task needs otherwise.',
      },
      {
        label: 'Give it its own database connection',
        text: 'a separate pool, with a low maximum and a statement timeout, so a careless browse cannot take connections from live traffic.',
      },
      {
        label: 'Disable the editor whenever a scope is set',
        text: 'the two features are mutually exclusive in practice, and the configuration should be treated as if it enforced that.',
      },
      {
        label: 'Hide the tables that hold secrets',
        text: 'encrypted columns, tokens, credentials. Hidden rather than read-only, because the risk is reading them.',
      },
      {
        label: 'Treat the audit callback as required in production',
        text: 'wire it to the application’s audit log, not to stdout. A change made here is otherwise the one change with no record.',
      },
    ],
  },

  seeAlso: [
    { title: 'GORM Studio in Grit', href: '/docs/database/studio' },
    { title: 'Schema migrations', href: '/docs/systems/schema-migrations' },
    { title: 'Audit logging', href: '/docs/systems/audit-logging' },
  ],
}

export const PULSE: SystemDesign = {
  slug: 'pulse',
  name: 'Pulse',
  tagline:
    'Request tracing, query profiling, runtime metrics and alerting that run inside the process, with no collector, no Redis and no second deployment.',
  group: 'Embedded',
  packages: ['github.com/MUKE-coder/pulse/pulse', 'v1.2.0'],

  interview: {
    intro:
      '"Design a metrics and monitoring platform" is a classic prompt. This is the in-process version of that problem, which changes the answers in interesting ways: no collector to scale, and a hard rule that the measurement must not cost what it is measuring.',
    questions: [
      {
        q: 'Why build this in-process rather than running Prometheus and a collector?',
        a: 'For the application that is one container and a database, a monitoring stack is bigger and harder to operate than the thing it monitors, so it does not get installed and the answer to "why was it slow" stays "nobody knows". In-process costs a dependency instead of an architecture. It stops being the right trade once there are enough replicas that you need a combined view.',
        see: 'problem',
      },
      {
        q: 'What stops the profiler from slowing down the requests it measures?',
        a: 'Storing a metric is a buffer push or a queue append, and the call returns immediately. Nothing on the request path waits on storage. The moment measurement is a measurable fraction of request time, every latency number it reports includes itself, which makes the tool actively misleading rather than merely expensive.',
        see: 'requirements',
      },
      {
        q: 'How do you keep memory bounded under unexpected traffic?',
        a: 'A fixed-size ring buffer, so the bound is the data structure rather than a cleanup goroutine that is supposed to keep up. The oldest entry is overwritten rather than memory growing. That is a much stronger guarantee than any policy, and it is what lets the thing run in a container with a limit.',
        see: 'requirements',
      },
      {
        q: 'You sample traces to control cost. Does that not make your error rate wrong?',
        a: 'It would, which is why sampling and counting are separated. Traces are sampled; per-minute rollups count every single request. Totals, error rates and objectives are therefore exact at any sample rate. Computing an error rate from a ten per cent sample gives an answer with no stated error bar that people then act on, which is worse than having no number.',
        see: 'capacity',
      },
      {
        q: 'Your retention says 24 hours but the buffer holds two minutes of traffic. Who tells the user?',
        a: 'The tool does. It reports when buffer capacity rather than the retention setting is what limits history, along with dropped writes and unclean restarts. A monitoring system that silently holds two minutes while its settings claim a day is actively misleading, and of all tools it is the one that must not hide its own gaps.',
        see: 'bottlenecks',
      },
      {
        q: 'How do you find an N plus one query?',
        a: 'By tallying repeated query patterns within a single request and flagging when the same shape runs more than a handful of times. It is the single most common cause of a slow page and the hardest to see in a log of individual queries, because each one is fast. Recording the file and line that issued the query is what turns "this query is slow" into "this line is slow".',
        see: 'high-level-design',
      },
      {
        q: 'Capturing request bodies sounds risky.',
        a: 'It is, in two ways. It is a security risk, handled by redacting passwords, tokens, card numbers and keys before anything is stored rather than before it is displayed. And it is a correctness risk: in one version here the error middleware read the body and put back only the first four kilobytes, so every request with a Content-Length reached the handler truncated. Mobile and curl uploads failed while browsers, which send chunked, did not.',
        see: 'bottlenecks',
      },
      {
        q: 'You have four replicas. Where is the dashboard?',
        a: 'There are four of them, and whichever one the load balancer picks is what you see. That is the structural limit of an in-process design and the point at which a real collector earns its cost. The bridge is an OpenTelemetry exporter: keep the in-process detail, ship spans out for the combined view.',
        see: 'scaling',
      },
      {
        q: 'What should alerting look like?',
        a: 'Two-phase firing so a single spike does not page anybody, a cooldown so one problem is not fifty notifications, and burn-rate alerting against an objective rather than a raw threshold. A threshold on a noisy metric pages people for nothing, and the second time that happens the alert stops being read.',
        see: 'low-level-design',
      },
      {
        q: 'Memory or SQLite for storage?',
        a: 'Memory is a ring buffer with the highest write ceiling and loses everything on restart. SQLite appends to a queue committed in batches by one writer, so it survives a restart at a lower ceiling. For an application restarted more often than it is busy, which is most of them, surviving the restart is worth more than the throughput.',
        see: 'capacity',
      },
    ],
  },

  problem: {
    text: [
      'Something is slow. Answering why needs per-request latency, the queries each request ran, how long each took, where in the code they came from, what the heap and the goroutine count were doing at the time, and which of those is unusual. That is a known problem with known answers, and all of the standard ones are a second deployment: a collector, a time series database, a dashboard, and a budget line.',
      'For a large team that is correct. For the application that is one container and a database, it is a monitoring stack that is bigger and harder to run than the thing it monitors, so it does not get installed, and the answer to "why was it slow" stays "nobody knows".',
      'An in-process profiler changes the trade. It costs a dependency instead of an architecture. What it must not cost is the thing it is measuring: a monitor that adds latency to the request path, or that grows without bound until it is the reason for the outage, has made the problem it was bought to solve.',
      'There is a subtler requirement too. The moment you sample to control cost, every total computed from the samples is wrong, and an error rate that is wrong in an unknown direction is worse than no error rate. So the counts and the traces have to be separated: sample what you store, count everything.',
    ],
    capabilities: [
      'Trace every request with a latency, a trace identifier and a route.',
      'Capture every database query with its duration and the file and line that issued it.',
      'Detect N plus one patterns within a single request.',
      'Sample continuously: heap, goroutines, garbage collection pauses.',
      'Capture panics and errors with a stack trace, fingerprinted so duplicates collapse.',
      'Run health checks on a schedule and answer Kubernetes probes.',
      'Alert on a threshold without firing on a single spike.',
      'Compute service level objectives from exact counts, not from samples.',
      'Never block the request path on storage.',
      'Say when its own buffers, rather than the retention setting, are what limits history.',
    ],
  },

  functional: [
    'Middleware for tracing and for error and panic capture.',
    'A GORM plugin capturing queries, callers, N plus one patterns and pool statistics.',
    'A runtime sampler on a five second interval, with goroutine leak detection.',
    'Memory ring buffer storage, or SQLite for survival across restarts.',
    'Per-minute rollups of every request, independent of the sample rate.',
    'A threshold alert engine with two-phase firing and a cooldown.',
    'Multiwindow burn-rate alerting on service level objectives.',
    'Secret redaction over captured bodies, errors and outbound URLs.',
    'Notification channels: Slack, Discord, email, and webhooks signed with HMAC.',
    'A Prometheus exposition endpoint, and JSON or CSV export.',
    'A WebSocket feed with per-client channel subscriptions.',
    'An embedded React dashboard of twelve pages.',
  ],

  nonFunctional: [
    {
      label: 'Never in the request path',
      text: 'storing a metric is a ring buffer push or a queue append. The request never waits for storage, because a profiler that adds latency is measuring itself.',
    },
    {
      label: 'Bounded by construction',
      text: 'fixed-size ring buffers with an atomic head index. Memory does not grow with traffic, which is what lets it run in a container with a limit.',
    },
    {
      label: 'Exact counts despite sampling',
      text: 'per-minute rollups count every request. Totals, error rates, objectives and load-test comparisons are therefore independent of the sample rate and of how much the buffer held.',
    },
    {
      label: 'Honest about its own limits',
      text: 'it alerts when buffer capacity rather than the retention setting is what bounds history, when writes were dropped, and after an unclean restart. A monitor quietly losing data is worse than no monitor.',
    },
    {
      label: 'Redacts before storing',
      text: 'passwords, tokens, card numbers, JWTs and keys are stripped from bodies, error messages and outbound URLs whatever the field is called. A profiler is otherwise a credential store nobody meant to build.',
    },
    {
      label: 'No CGo, no external service',
      text: 'pure Go SQLite and an embedded frontend. The reason it is installed is that installing it is one line.',
    },
  ],

  capacity: {
    assumptions: [
      ['Requests per second', '926 average, 4,630 peak'],
      ['Queries per request', '4 to 12'],
      ['Request record', '~500 bytes'],
      ['Query record', '~300 bytes'],
      ['Runtime sample interval', '5 seconds'],
      ['Default retention', '24 hours'],
    ],
    estimates: [
      {
        label: 'Per-request overhead',
        working: [
          'a trace identifier, a timer, and a ring buffer push',
          'the push is an atomic index plus a short write lock',
          'microseconds against a request measured in milliseconds',
        ],
        note: 'The design requirement is that this stays invisible. The moment it is a percentage of request time, every latency number it reports includes itself.',
      },
      {
        label: 'Memory at full sampling',
        working: [
          '926 req/s x 500 bytes = ~460 KB/s of request records',
          'plus 926 x 8 queries x 300 bytes = ~2.2 MB/s of query records',
          'a ring buffer sized for 100,000 requests holds ~50 MB and about 2 minutes',
        ],
        note: 'Which is exactly why the capacity warning exists. At this rate the buffer, not the 24 hour retention setting, is what decides how much history there is, and without being told nobody would know.',
      },
      {
        label: 'Sampling, and what it does not touch',
        working: [
          'at SampleRate 0.1: 93 traces/s stored instead of 926',
          'memory and query volume fall by ten',
          'the per-minute rollups still count all 926',
        ],
        note: 'This is the central design decision. Sampling controls storage; the rollups keep totals and error rates exact. Computing an error rate from a 10% sample would give an answer with no stated error bar, which people would then act on.',
      },
      {
        label: 'SQLite throughput',
        working: [
          'appends go to a queue, committed in batches by one writer',
          'one writer, so writes serialise, which SQLite needs anyway',
          'survives restarts, at a lower ceiling than the ring buffer',
        ],
      },
      {
        label: 'Runtime sampling',
        working: [
          'one sample every 5 seconds = 17,280/day',
          'at ~200 bytes = ~3.5 MB/day',
        ],
      },
    ],
  },

  highLevel: {
    intro:
      'Middleware and a database plugin feed a storage layer that never blocks, and everything else reads from it.',
    components: [
      {
        label: 'Tracing middleware',
        text: 'a trace identifier, a timer and a route label per request. Sampled for storage, counted always.',
      },
      {
        label: 'Error middleware',
        text: 'recovers panics, captures stack traces and request context, fingerprints the error so repeats collapse into one entry with a count.',
      },
      {
        label: 'GORM plugin',
        text: 'every query with its duration and the file and line that issued it, plus per-request tallies of repeated patterns for N plus one detection, plus pool statistics.',
      },
      {
        label: 'Runtime sampler',
        text: 'heap, goroutines and garbage collection on a five second tick, with a leak threshold.',
      },
      {
        label: 'Storage',
        text: 'a generic ring buffer in memory, or a batched SQLite writer. Both are append-and-return; neither makes a request wait.',
      },
      {
        label: 'Aggregator',
        text: 'per-minute rollups of every request, which is what makes totals and objectives exact while traces are sampled.',
      },
      {
        label: 'Alert engine',
        text: 'threshold rules with two-phase firing so a single spike does not page anybody, a cooldown, and signed webhook delivery.',
      },
      {
        label: 'Capacity watch',
        text: 'notices when the buffer rather than the retention setting bounds history, when writes were dropped, and when the last shutdown was unclean.',
      },
    ],
    flow: {
      title: 'A request measured, and what it costs the request',
      nodes: [
        { id: 'req', label: 'Request', col: 0, row: 0 },
        { id: 'trace', label: 'Tracing MW', col: 1, row: 0, accent: true },
        { id: 'handler', label: 'Handler', col: 2, row: 0 },
        { id: 'gorm', label: 'GORM Plugin', col: 3, row: 0, accent: true },
        { id: 'n1', label: 'N+1 Tally', col: 3, row: 1 },
        { id: 'redact', label: 'Redact', col: 2, row: 1, accent: true },
        { id: 'buffer', label: 'Ring Buffer', col: 1, row: 1, accent: true },
        { id: 'rollup', label: 'Per-Minute Rollup', col: 1, row: 2, accent: true },
        { id: 'dash', label: 'Dashboard + Alerts', col: 2, row: 2 },
      ],
      edges: [
        { from: 'req', to: 'trace', step: 1 },
        { from: 'trace', to: 'handler', step: 2 },
        { from: 'handler', to: 'gorm', step: 3 },
        { from: 'gorm', to: 'n1', step: 4 },
        { from: 'gorm', to: 'redact', step: 5 },
        { from: 'redact', to: 'buffer', step: 6 },
        { from: 'trace', to: 'rollup', step: 7, bend: 'v' },
        { from: 'rollup', to: 'dash', step: 8 },
      ],
      steps: [
        'A request arrives. The tracing middleware starts a timer and attaches a trace identifier, which also goes out in the response so a complaint can be tied to a record.',
        'The handler runs. Nothing about it knows it is being measured, which is the point: instrumentation each handler participates in is instrumentation half of them forget.',
        'Every query the handler issues passes through the GORM plugin, which times it and records the file and line it came from. The caller is what turns "this query is slow" into "this line is slow".',
        'Repeated query patterns inside one request are tallied. Five or more of the same shape is an N plus one, which is the single most common cause of a slow page and the hardest to see from a log of individual queries, because each one is fast.',
        'Captured context is redacted before it is stored: passwords, tokens, card numbers, JWTs, keys, whatever the field is named. Without this a profiler becomes a store of credentials that nobody decided to build and nobody is guarding.',
        'The records are pushed into a fixed-size ring buffer, or appended to a queue that one writer commits in batches. Either way the call returns immediately; the request does not wait on storage, because a profiler that adds latency is measuring itself.',
        'Separately, and regardless of the sample rate, the request is counted into a per-minute rollup. This is the decision that makes the numbers trustworthy: traces are sampled, counts are complete, so an error rate is a fact rather than an estimate with no error bar.',
        'The dashboard, the alert engine and the objective burn rates all read from the rollups for anything that is a total, and from the buffers for anything that is an example. When the buffer rather than the retention setting is what limits history, it says so.',
      ],
    },
    dataFlow: [
      'Sampling decides what is stored, never what is counted. Every total in the product comes from the rollups.',
      'Redaction happens before storage, not before display, so a leaked secret is never written down in the first place.',
      'The ring buffer is a fixed slice with an atomic head. Bounded memory is not a tuning option, it is the data structure.',
      'The SQLite path trades peak write throughput for surviving a restart, which is the right trade for an application that is restarted more often than it is busy.',
    ],
  },

  stack: [
    ['Language', 'Go, no CGo'],
    ['Host', 'Gin and GORM'],
    ['Storage', 'in-memory ring buffers, or pure Go SQLite'],
    ['Default retention', '24 hours'],
    ['Defaults', 'sample rate 1.0, slow request 1s, slow query 200ms, N+1 at 5'],
    ['Runtime sampling', 'every 5 seconds, leak threshold 100 goroutines'],
    ['Export', 'Prometheus, JSON, CSV, OpenTelemetry'],
    ['Dashboard', '12 pages embedded in the binary'],
  ],

  api: {
    groups: [
      {
        title: 'Mounted under the configured prefix',
        rows: [
          { method: 'GET', path: '/pulse/ui', what: 'The embedded dashboard' },
          { method: 'GET', path: '/pulse/api/requests', what: 'Traces, filtered by route, status and latency' },
          { method: 'GET', path: '/pulse/api/db/n1', what: 'N plus one findings grouped by route and pattern' },
          { method: 'GET', path: '/pulse/api/slos', what: 'Compliance, error budget and burn rates' },
          { method: 'GET', path: '/pulse/metrics', what: 'Prometheus exposition' },
          { method: 'GET', path: '/pulse/live', what: 'Kubernetes liveness' },
          { method: 'GET', path: '/pulse/ready', what: 'Kubernetes readiness, composite over health checks' },
        ],
      },
    ],
    samples: [
      {
        title: 'How Grit mounts it, including the bug the comment is there to prevent',
        language: 'go',
        code: `pulseOpts := []pulse.Option{
    pulse.WithAppName(cfg.AppName),
    pulse.WithCredentials(cfg.PulseUsername, cfg.PulsePassword),
    pulse.WithExcludePaths("/studio/*", "/sentinel/*", "/docs/*", "/pulse/*"),
    pulse.WithPrometheus(),
    // Request-body capture needs v1.0.1 or later. Before that the error
    // middleware read the body and put back only the first 4 KB, so every
    // request carrying a Content-Length reached the handler truncated:
    // uploads from mobile and curl failed while browsers, which send
    // chunked, did not. Pin below v1.0.1 and you want
    // pulse.WithRequestBodyCaptureDisabled() back.
}
if cfg.PulseStorage == "sqlite" && cfg.PulseStorageDSN != "" {
    pulseOpts = append(pulseOpts, pulse.WithSQLite(cfg.PulseStorageDSN))
}
p := pulse.Mount(context.Background(), r, db, pulseOpts...)`,
      },
    ],
  },

  lowLevel: {
    classes: [
      {
        name: 'RingBuffer[T]',
        what: 'O(1) append into a fixed-size slice behind an atomic head index and a short write lock. The reason memory is bounded regardless of traffic.',
      },
      {
        name: 'GORM plugin',
        file: 'pulse/gorm_plugin.go',
        what: 'Times every query, records the caller, and keeps per-request tallies of repeated patterns. Tallies have an idle expiry so an abandoned request cannot leak one.',
      },
      {
        name: 'Aggregator',
        file: 'pulse/aggregator.go',
        what: 'Per-minute rollups of every request. The separation between what is sampled and what is counted lives here.',
      },
      {
        name: 'Capacity watch',
        file: 'pulse/capacity.go',
        what: 'Reports when a ring buffer rather than the retention setting is what bounds history. A monitor that silently holds two minutes of data while its settings claim a day is actively misleading.',
      },
      {
        name: 'Pointer config fields',
        what: '*bool and *float64 so "not set" and "explicitly zero" are different. A sample rate of 0 and an unset sample rate must not mean the same thing.',
      },
    ],
    principles: [
      {
        label: 'Measuring must not be measurable',
        text: 'push and return. Every design choice in the storage layer exists so the request path never waits.',
      },
      {
        label: 'Sample the examples, count the totals',
        text: 'sampling is necessary and making totals approximate is not. Rollups keep the numbers exact at a cost that does not depend on the sample rate.',
      },
      {
        label: 'Bounded by the data structure',
        text: 'a ring buffer cannot grow. That is a stronger guarantee than a cleanup goroutine that is supposed to keep up.',
      },
      {
        label: 'Report your own failures',
        text: 'dropped writes, truncated history and unclean restarts are surfaced. A monitoring tool that hides its own gaps is the one tool that must not.',
      },
      {
        label: 'Redact on the way in',
        text: 'before storage, not before display. The alternative is a database of captured secrets with a filter in front of it.',
      },
    ],
    patterns: [
      ['Ring buffer', 'bounded retention as a data structure'],
      ['Batched single writer', 'durability without blocking the request'],
      ['Fingerprinting', 'duplicate errors collapsed into one entry with a count'],
      ['Two-phase alerting', 'a threshold must hold, not just be touched'],
      ['Multiwindow burn rate', 'the Google SRE workbook objective alerting'],
    ],
  },

  scaling: [
    'Overhead per request is microseconds and does not grow with traffic, because the work is a timer and a buffer push.',
    'Memory is fixed by the buffer sizes. The question is never how much memory it will use, it is how much history that memory buys, which is what the capacity warning answers.',
    'Sampling is the lever for storage volume, and it costs nothing in accuracy because the rollups count everything.',
    'The SQLite backend serialises on one writer, which is the right shape for SQLite and a lower ceiling than memory. Choose it for surviving restarts, not for throughput.',
    'Everything is per process. Behind several replicas each has its own dashboard and its own view, which is the main structural limit of an in-process design and the point at which a real collector earns its cost.',
    'The OpenTelemetry exporter is the bridge out: keep the in-process dashboard for the single-instance case and ship spans to a collector when there are enough instances for that to be worth running.',
  ],

  bottlenecks: {
    problems: [
      {
        label: 'Per-process view',
        text: 'with four replicas there are four dashboards and no combined picture, and the one you reach through the load balancer is whichever replica answered.',
      },
      {
        label: 'Buffer rather than retention',
        text: 'at high traffic a memory buffer holds minutes while the configured retention says a day. It is surfaced, and it still surprises people.',
      },
      {
        label: 'Capture that changes behaviour',
        text: 'reading a request body to record it is a real risk to the request. In v1.0.0 it truncated every body with a Content-Length to 4 KB, so mobile and curl uploads failed while browsers did not.',
      },
      {
        label: 'Dashboard exposure',
        text: 'it holds request bodies, errors and query text. Redaction reduces what it contains; it does not make the dashboard safe to leave open.',
      },
      {
        label: 'Alert noise',
        text: 'threshold rules on a noisy metric page people for nothing, and the second time that happens the alert stops being read.',
      },
    ],
    improvements: [
      {
        label: 'Export to a collector once there are replicas',
        text: 'the OpenTelemetry path gives one combined view. Keep Pulse mounted for the per-process detail it is good at.',
      },
      {
        label: 'Size buffers from the actual rate',
        text: 'requests per second times the history you want. Or move to SQLite, which trades write ceiling for a retention setting that means what it says.',
      },
      {
        label: 'Pin the version when capture is on',
        text: 'v1.2.0 or later wraps the body and passes it through whole. Below v1.0.1, turn request body capture off.',
      },
      {
        label: 'Authenticate the dashboard and exclude it from itself',
        text: 'credentials always, and exclude the monitoring routes from tracing so the tool does not fill its own buffers with views of itself.',
      },
      {
        label: 'Prefer burn-rate alerts to thresholds',
        text: 'multiwindow burn rate fires on budget actually being consumed rather than on an instant being bad, which is most of the difference between an alert people act on and one they mute.',
      },
    ],
  },

  seeAlso: [
    { title: 'Observability', href: '/docs/systems/observability' },
    { title: 'Pulse in Grit', href: '/docs/operations/pulse' },
  ],
}

export const SENTINEL: SystemDesign = {
  slug: 'sentinel',
  name: 'Sentinel',
  tagline:
    'A web application firewall, threat intelligence and audit layer that mounts in one call, and the accuracy numbers to say how often it is wrong.',
  group: 'Embedded',
  packages: ['github.com/MUKE-coder/sentinel/v2', 'v2.6.0'],

  interview: {
    intro:
      '"Design a web application firewall" and "design a DDoS protection layer" are real prompts, and this is the in-application version. The questions worth preparing are about false positives, about keeping the request path cheap, and about what an attack does to your own analysis.',
    questions: [
      {
        q: 'What is the dominant failure mode of a firewall inside an application?',
        a: 'Blocking legitimate traffic, not missing an attack. A pattern that matches real requests takes the site down for the people it was protecting, and it does so quietly, as a fraction of users who cannot complete something. This library shipped that twice: an SSRF pattern matched 0.0.0.0 inside browser version strings like Chrome/140.0.0.0, blocking every Chrome user, and a bare double hyphen matched inside base64url cookies, rejecting about one session in ten at random.',
        see: 'problem',
      },
      {
        q: 'So how do you know your rules are any good?',
        a: 'Measure both rates against a fixed corpus and pin them in continuous integration. Here that is 89 legitimate-but-suspicious requests alongside 56 attacks, giving ten per cent false positives and one hundred per cent detection, down from thirty-eight per cent false positives. Without a corpus, a rule change is a guess with production as the test, and the accuracy of the whole thing is simply unknown.',
        see: 'requirements',
      },
      {
        q: 'How do you roll it out without breaking everyone?',
        a: 'Log mode first. Run it against real traffic for a week and look at what it would have refused, then switch to blocking once the false positives are understood. The generated projects do exactly this, logging in development and blocking in production, and exclude the authenticated rich-text routes whose bodies are legitimately full of markup.',
        see: 'high-level-design',
      },
      {
        q: 'An attacker double-encodes the payload. Does your pattern still match?',
        a: 'Only if you decode until there is nothing left to decode. Decoding once means a double-encoded payload passes through unexamined, which was a real bypass here. Every layer has to be unwrapped across the path, the parameters and form bodies, and then every layer scanned.',
        see: 'high-level-design',
      },
      {
        q: 'How do you keep inspection off the critical path?',
        a: 'Split the decision from the analysis. Blocking is synchronous and cheap: decode, match, decide. Everything else, scoring, actor profiling, reputation lookups, anomaly detection and the audit write, goes through a bounded buffer to background workers. That is what makes the expensive half affordable.',
        see: 'high-level-design',
      },
      {
        q: 'What happens to that buffer during an actual attack?',
        a: 'It fills, which is exactly when losing the records is worst. So drops are counted and surfaced on the dashboard rather than silently discarded. An attack large enough to overwhelm the analysis must not also erase the evidence of itself, and the only way to know it happened is a counter somebody is watching.',
        see: 'capacity',
      },
      {
        q: 'What bounds the worst-case cost of inspecting one request?',
        a: 'A body size cap, with oversized bodies rejected rather than partly scanned. Without one, the attacker chooses how much processor time each request costs by choosing how large a body to send. Scanning the first ten kilobytes and passing the rest is the other failure, since it is a bypass rather than a limit.',
        see: 'requirements',
      },
      {
        q: 'Behind a load balancer, how do you know the client address?',
        a: 'By walking the forwarded chain from the right, past the proxies you trust. Taking the leftmost entry trusts a value the client wrote, so anybody can choose their own address and walk past rate limits, address blocks and the lockout. The safe default is to ignore the header entirely unless trusted proxy ranges are configured.',
        see: 'high-level-design',
      },
      {
        q: 'You wrote an exclusion for /api/blogs/:id and it does nothing. Why?',
        a: 'Because exclusions match the real request path, not the router template, so it matches the literal string ":id" and never "/api/blogs/123". It is dead config that reads as correct, which is the worst kind. The general defence is a validator that refuses configuration that compiles and does nothing, and the specific one is exporting the matcher so a test can assert a concrete production path is excluded.',
        see: 'bottlenecks',
      },
      {
        q: 'Should a security layer fail open or closed?',
        a: 'It depends on what failed, and the distinction is the interesting part. If your own dependency is unavailable, such as the shared counter store, allow the request: a limiter that takes the site down when Redis blinks is an outage with a security rationale. If the caller is hostile and a pattern matched, refuse. Our infrastructure failing and the caller attacking are different situations and should not share a policy.',
        see: 'requirements',
      },
    ],
  },

  problem: {
    text: [
      'An application on the public internet is scanned within hours of getting a DNS record. Most of it is automated and most of it is pointless, and the small remainder is the reason the rest cannot be ignored: SQL injection against a parameter somebody forgot to bind, a path traversal in a file endpoint, credential stuffing against the login, an open redirect in a callback.',
      'The platform answers are a managed firewall in front of the application, which costs money and sees only what the edge sees, or a library, which sees the parsed request and the authenticated user and runs for the price of a dependency.',
      'A firewall inside the application has one dominant failure mode, and it is not missing an attack. It is blocking legitimate traffic. A pattern that matches real requests takes the site down for the people it was protecting, and it does it quietly, as a fraction of users who cannot complete something. Sentinel shipped exactly that twice: an SSRF pattern matched `0.0.0.0` inside browser version strings like `Chrome/140.0.0.0`, which blocked every Chrome user, and a bare `--` matched inside base64url cookies, which rejected about one session in ten at random.',
      'That is why the interesting number here is not detection but false positives, and why the accuracy is measured against a corpus and pinned in CI rather than asserted.',
    ],
    capabilities: [
      'Inspect the path, query, body and headers for known attack shapes.',
      'Run in log mode or block mode, with configurable strictness per rule family.',
      'Decode layered encoding before matching, so one more percent sign is not a bypass.',
      'Limit requests per address, per user, per route and globally, shared across replicas.',
      'Lock out brute force against a login, with an optional CAPTCHA tier.',
      'Profile threat actors and check addresses against reputation feeds.',
      'Keep an audit log that can be shown not to have been edited.',
      'Report accuracy against a fixed corpus, so a rule change has a measurable cost.',
      'Refuse configuration that compiles and silently does nothing.',
      'Do all of it without the request waiting on any of the analysis.',
    ],
  },

  functional: [
    'Mount on a Gin router, with MountE returning an error rather than killing the host.',
    'Detection for SQL injection, XSS, path traversal, command injection, SSRF and XXE.',
    'Four sensitivity levels per rule family, and custom rules with a block or log action.',
    'Recursive decoding of layered encoding across path, parameters and form bodies.',
    'Route exclusions with wildcard and globstar patterns.',
    'A body inspection cap, with oversized bodies rejected rather than partly scanned.',
    'Rate limiting with fixed window, sliding window and token bucket strategies.',
    'A shared counter store interface, with a Redis implementation and a conformance suite.',
    'Auth shield: failed login budget, lockout, optional CAPTCHA.',
    'Trusted proxy handling that walks the forwarded chain from the right.',
    'A GORM plugin recording database changes, panic-guarded so auditing cannot break a write.',
    'A hash-chained audit log with optional HMAC and a verify endpoint.',
    'An asynchronous event pipeline with a ring buffer and visible drop counts.',
    'Live config: dashboard changes that survive a restart and reach every replica.',
    'Compliance reports carrying provenance and a warning when the data cannot support them.',
    'An embedded dashboard of fifteen pages, plus alerting to Slack, email, webhook and PagerDuty.',
  ],

  nonFunctional: [
    {
      label: 'False positives are the headline number',
      text: 'with the default rules, 10% false positives and 100% detection against a corpus of 89 legitimate-but-suspicious requests and 56 attacks, pinned in CI. It was 38% and 92%. A firewall whose accuracy is not measured is a firewall whose accuracy is unknown.',
    },
    {
      label: 'Analysis never blocks the request',
      text: 'events go to a ring buffer and are processed by background workers. The decision to block is synchronous and cheap; everything after it is not.',
    },
    {
      label: 'Drops are visible',
      text: 'the pipeline counts what it could not buffer, and the count is on the dashboard. An attack large enough to overwhelm analysis must not also hide itself.',
    },
    {
      label: 'Dead config is refused',
      text: 'config that compiles, reads as correct and does nothing was the common thread behind four separate issues. Mount validates and logs, and the validator can be called to fail a deploy.',
    },
    {
      label: 'Shared state across replicas',
      text: 'counters and lockouts in Redis, because counted per process behind N instances a client gets N times every limit and N times the failed-login budget.',
    },
    {
      label: 'Fails open on infrastructure',
      text: 'if Redis is down, requests are allowed. A rate limiter that takes the site down when its counter store is unavailable has become the outage.',
    },
    {
      label: 'No published secret',
      text: 'an unset dashboard key is random per process, the default password works from localhost only, and WebSocket handshakes must be same-origin.',
    },
  ],

  capacity: {
    assumptions: [
      ['Requests per second', '926 average, 4,630 peak'],
      ['Attack traffic', '1% to 20%, bursty'],
      ['Body inspection cap', '1 MB as Grit configures it'],
      ['Pipeline ring buffer', '10,000 events'],
      ['Audit retention', '365 days'],
    ],
    estimates: [
      {
        label: 'Inspection cost',
        working: [
          'decode, then match a compiled pattern set against path, query and body',
          'tens of microseconds for a normal request',
          'proportional to body size, which is why the cap matters',
        ],
        note: 'The cap is the whole story on worst case. Without one, an attacker chooses how much CPU each request costs by choosing how large a body to send.',
      },
      {
        label: 'What a 10% false positive rate means',
        working: [
          '926 req/s, suppose 2% look legitimate-but-suspicious',
          '18.5 req/s through the suspicious corpus shape',
          'at 10% false positives: ~1.9 legitimate requests blocked per second',
        ],
        note: 'Which is why block mode belongs behind route exclusions for the endpoints that carry rich text, and why Grit runs log mode in development. The number is small and it is not zero, and the people it hits cannot tell you why.',
      },
      {
        label: 'Pipeline headroom',
        working: [
          '10,000 event buffer',
          'at 926 events/s a worker stall of 10 s fills it',
          'past that, events are dropped and counted',
        ],
        note: 'An attack is exactly when the event rate spikes and exactly when losing the record is worst, so the drop count is on the dashboard rather than in a log.',
      },
      {
        label: 'Rate limit cost',
        working: [
          'in-memory: a map lookup and some arithmetic',
          'Redis: one round trip per limited request, down from two in v2.6.0',
          '~0.3 ms added when the store is shared',
        ],
      },
      {
        label: 'Audit volume',
        working: [
          'database changes plus dashboard actions plus logins',
          '~50,000 entries/day at this traffic, at ~600 bytes',
          '= ~30 MB/day, ~11 GB over the 365 day retention',
        ],
      },
    ],
  },

  highLevel: {
    intro:
      'A synchronous decision on the request path, and everything expensive behind a buffer.',
    components: [
      {
        label: 'Middleware chain',
        text: 'client address resolution, rate limit, WAF, auth shield. The only parts that can block, and all of them cheap.',
      },
      {
        label: 'Decoder',
        text: 'unwraps layered encoding before matching. Values were decoded once, so a double-encoded payload passed straight through.',
      },
      {
        label: 'Detector',
        text: 'the pattern set, scoped per rule to where it may match. Patterns that scan headers as well as parameters are how browser version strings came to look like SSRF.',
      },
      {
        label: 'Scorer and classifier',
        text: 'severity, a CVSS vector per threat type, and sensitivity levels that change what counts as a match.',
      },
      {
        label: 'Counter store',
        text: 'the interface behind rate limiting and the failed-login budget. In process by default, Redis for shared, with a conformance suite for other implementations.',
      },
      {
        label: 'Event pipeline',
        text: 'a buffered channel and background workers. Everything that is not the block decision happens here, with drops counted.',
      },
      {
        label: 'Intelligence',
        text: 'threat actor profiling, reputation lookups and anomaly detection, all fed by the pipeline rather than by the request.',
      },
      {
        label: 'Live config',
        text: 'settings changed from the dashboard, stored and propagated to every replica within the sync interval.',
      },
      {
        label: 'Audit chain',
        text: 'hash-chained entries with optional HMAC, and an endpoint that verifies them.',
      },
    ],
    flow: {
      title: 'A request inspected, blocked, and analysed afterwards',
      nodes: [
        { id: 'client', label: 'Client', col: 0, row: 0 },
        { id: 'ip', label: 'Trusted Proxy Walk', col: 1, row: 0, accent: true },
        { id: 'limit', label: 'Rate Limit', col: 2, row: 0, accent: true },
        { id: 'counters', label: 'Counter Store', col: 3, row: 0 },
        { id: 'decode', label: 'Decode Layers', col: 2, row: 1, accent: true },
        { id: 'waf', label: 'Pattern Match', col: 3, row: 1, accent: true },
        { id: 'block', label: '403 or Log', col: 1, row: 1 },
        { id: 'pipeline', label: 'Event Pipeline', col: 2, row: 2, accent: true },
        { id: 'intel', label: 'Intelligence + Audit', col: 3, row: 2 },
      ],
      edges: [
        { from: 'client', to: 'ip', step: 1 },
        { from: 'ip', to: 'limit', step: 2 },
        { from: 'limit', to: 'counters', step: 3 },
        { from: 'limit', to: 'decode', step: 4 },
        { from: 'decode', to: 'waf', step: 5 },
        { from: 'waf', to: 'block', step: 6, dashed: true },
        { from: 'waf', to: 'pipeline', step: 7 },
        { from: 'pipeline', to: 'intel', step: 8 },
      ],
      steps: [
        'A request arrives. Before anything else the client address is resolved by walking the forwarded chain from the right past the configured trusted proxies. Taking the leftmost entry lets any client write its own address and walk past rate limits, address blocks and the lockout, which was the v2.2.2 fix.',
        'The rate limiter counts the request against its keys: address, user, route, global. The strategy is a real sliding window by default, approximated from two fixed windows, with token bucket and fixed window also available.',
        'Counting happens in the shared store when one is configured. Per process, behind four replicas, a client got four times every limit. If the store is unreachable the request is allowed, because a limiter that fails closed is an outage with a security rationale.',
        'The firewall decodes before it matches. Leftover encoding in the path, parameters and form bodies is unwrapped, because a single decode pass meant a double-encoded payload went through unexamined.',
        'The pattern set runs, each rule scoped to where it is allowed to match. That scoping is the lesson from two production incidents: an unscoped host pattern matched browser version strings, and a bare statement terminator matched inside session cookies.',
        'A match blocks or logs, depending on mode. Grit logs in development and blocks in production, and excludes the authenticated rich-text admin routes, whose bodies are legitimately full of markup that the XSS rules would otherwise flag on every save.',
        'The event goes to a ring buffer and the request continues. Nothing after this point is on the request path, which is what lets the analysis be expensive.',
        'Background workers score it, profile the actor, check reputation feeds, run anomaly detection, write the audit entry into the hash chain, and push it to the dashboard. If the buffer was full the event is dropped and the drop is counted, because an attack large enough to overwhelm the analysis must not also erase the evidence of itself.',
      ],
    },
    dataFlow: [
      'The block decision is synchronous and everything else is not. That split is what makes a firewall affordable in the request path.',
      'Every pattern declares where it may match. A pattern with no declared scope will eventually match something in a header that nobody considered.',
      'Counters and lockouts are shared state when a store is configured, and the semantics are pinned by a conformance suite so another implementation behaves the same way.',
      'Data sent to an AI provider is redacted by default: query values, bodies, full addresses and personal data are stripped first.',
    ],
  },

  stack: [
    ['Language', 'Go 1.24+, no CGo'],
    ['Host', 'Gin, mounted with one call'],
    ['Storage', 'SQLite, PostgreSQL'],
    ['Shared counters', 'Redis, single, sentinel or cluster'],
    ['Limiter strategies', 'sliding window (default), fixed window, token bucket'],
    ['Pipeline', 'a 10,000 event ring buffer with background workers'],
    ['Audit', 'hash chain, optional HMAC, 365 day retention'],
    ['Accuracy', '10% false positive, 100% detection, pinned in CI'],
    ['Dashboard', '15 pages embedded, 50+ API endpoints'],
  ],

  api: {
    groups: [
      {
        title: 'Mounted under /sentinel',
        rows: [
          { method: 'GET', path: '/sentinel/ui', what: 'The embedded dashboard' },
          { method: 'GET', path: '/sentinel/api/threats', what: 'Threat events, with severity and CVSS' },
          { method: 'GET', path: '/sentinel/api/audit-logs/verify', what: 'Verify the hash chain' },
          { method: 'GET', path: '/sentinel/api/performance/overview', what: 'Latency, error rates and pipeline drops' },
          { method: 'POST', path: '/sentinel/api/blocks', what: 'Block an address, 24 hours by default' },
          { method: 'DELETE', path: '/sentinel/api/settings/live', what: 'Discard dashboard overrides everywhere' },
          { method: 'POST', path: '/sentinel/csp-report', what: 'CSP violations into the same dashboard' },
        ],
      },
    ],
    samples: [
      {
        title: 'How Grit mounts it, and why each line is there',
        language: 'go',
        code: `// Counted per process, N replicas gave a client N times every limit.
var counters sentinel.CounterStore
if svc.Cache != nil {
    counters = redisstore.New(svc.Cache.Client())
}

// MountE, so a misconfiguration in dev does not kill the host.
err := sentinel.MountE(r, db, sentinel.Config{
    Counters: counters,
    WAF: sentinel.WAFConfig{
        Enabled: true,
        // Log in development, block in production.
        Mode: mode,
        // Empty means ignore X-Forwarded-For entirely, which is the
        // safe default. Populate it only behind a known proxy.
        TrustedProxies:      cfg.SentinelTrustedProxies,
        MaxBodyBytes:        1 * 1024 * 1024,
        RejectOversizedBody: true,
    },
})`,
      },
      {
        title: 'Asserting your own paths against your exclusions',
        language: 'go',
        code: `// Exclusions are matched against the real request path, not gin's route
// template, so "/api/blogs/:id" matches the literal ":id" and nothing else.
// That was silent dead config. Use a subtree match, and test it.
m := sentinel.NewRouteMatcher(cfg.WAF.ExcludeRoutes)
if !m.Matches("/v1/payments/collect") {
    t.Fatal("payments endpoint is not excluded from the WAF")
}`,
      },
    ],
  },

  lowLevel: {
    classes: [
      {
        name: 'sentinel.MountE',
        what: 'Mounts and returns an error. The earlier Mount killed the host process on a library failure, which is the wrong trade for a security add-on.',
      },
      {
        name: 'sentinel.ValidateConfig',
        what: 'Typed issues for every silent-dead-config trap: unmatchable route patterns, unknown storage drivers that fall back to memory, broken regexes, alert sinks with no credentials, unreachable CAPTCHA tiers. Run by Mount, callable to fail a deploy.',
      },
      {
        name: 'core.CounterStore',
        file: 'core/counters.go',
        what: 'The interface behind rate limiting and the lockout. The caller passes the clock so one request’s operations share a reading, and package countertest checks an implementation against the in-memory semantics.',
      },
      {
        name: 'pipeline.Pipeline',
        file: 'pipeline/pipeline.go',
        what: 'A buffered channel, background workers, and atomic counters for emitted and dropped. Default capacity 10,000.',
      },
      {
        name: 'detection',
        what: 'Decoding, patterns, scoring and sensitivity, with a corpus test that pins the false positive and detection rates in CI.',
      },
      {
        name: 'middleware.ClientIP',
        what: 'Exported, because the dashboard login limiter used gin’s version, which trusts the forwarded header from anyone, and could therefore be reset by rotating a header.',
      },
    ],
    principles: [
      {
        label: 'Measure the false positives',
        text: 'a corpus of legitimate-but-suspicious requests alongside the attacks, and both rates pinned in CI. Without it a rule change is a guess with production as the test.',
      },
      {
        label: 'Scope every pattern',
        text: 'each rule declares where it may match. Both of the incidents that blocked real users were unscoped patterns matching somewhere nobody considered.',
      },
      {
        label: 'Config that does nothing must be an error',
        text: 'four separate issues shared one cause: configuration that compiled, read as correct and was never consulted. A validator is the only defence against that class.',
      },
      {
        label: 'Fail open on infrastructure, closed on attack',
        text: 'Redis unreachable allows the request; a matched attack pattern does not. The distinction is between our dependency failing and the caller being hostile.',
      },
      {
        label: 'Keep the request path cheap',
        text: 'decide synchronously, analyse asynchronously, and count what the buffer could not take.',
      },
    ],
    patterns: [
      ['Chain of responsibility', 'middleware each able to refuse'],
      ['Async pipeline', 'a bounded buffer between decision and analysis'],
      ['Strategy', 'three limiter algorithms behind one interface'],
      ['Conformance suite', 'a shared test proving another counter store behaves the same'],
      ['Hash chain', 'tamper-evident audit entries'],
      ['Corpus testing', 'accuracy as a pinned number rather than a claim'],
    ],
  },

  scaling: [
    'Inspection is per request and cheap, so it scales with the application. The body cap is what bounds the worst case, and without it the attacker picks the cost.',
    'The asynchronous pipeline is what keeps analysis off the request path, and its buffer is the thing that fills under attack, which is when it matters.',
    'Shared counters cost one Redis round trip per limited request, halved in v2.6.0. That is the price of limits that mean the same thing behind every replica.',
    'Live config propagates within the sync interval, five seconds, so a block made on one replica applies on the others without a deploy.',
    'Audit and threat data grow steadily and have a retention setting, and at a year of retention the audit table is the largest thing the library owns.',
    'The dashboard is per process like any embedded UI, so the combined picture across replicas comes from shared storage rather than from the page.',
  ],

  bottlenecks: {
    problems: [
      {
        label: 'False positives',
        text: 'the dominant risk, and the one that presents as a fraction of users unable to do something rather than as an error anybody sees.',
      },
      {
        label: 'Exclusions that match nothing',
        text: 'patterns are matched against the real path, not the route template, so an exclusion written as "/api/blogs/:id" matches the literal string and is dead config that looks correct.',
      },
      {
        label: 'Pipeline saturation under attack',
        text: 'the event rate spikes exactly when the records matter most, and past the buffer they are dropped.',
      },
      {
        label: 'Proxy configuration',
        text: 'trusting the forwarded header wrongly either lets clients spoof their address or makes every request appear to come from the load balancer, and the two mistakes look identical from inside.',
      },
      {
        label: 'An embedded dashboard holding attack data',
        text: 'it contains request payloads and actor profiles, so it is itself a target and must not be reachable with default credentials.',
      },
    ],
    improvements: [
      {
        label: 'Run log mode first, then block',
        text: 'a week in log mode on real traffic shows what would have been refused. Grit does this by default in development for the same reason.',
      },
      {
        label: 'Test the exclusions against real paths',
        text: 'the route matcher is exported so a test can assert that a concrete production path is excluded. That turns dead config into a failing test.',
      },
      {
        label: 'Alert on the drop counter',
        text: 'it is on the performance overview. Watching it is how you learn that an attack was large enough that the record is incomplete.',
      },
      {
        label: 'Be explicit about proxies',
        text: 'leave the trusted list empty to ignore the forwarded header entirely, or populate it with the actual proxy ranges. Never leave it to a default.',
      },
      {
        label: 'Set the dashboard secret and password',
        text: 'an unset key is random per process, so sessions end on every restart, and the default password works only from localhost. Both are deliberate and both need a real value in production.',
      },
    ],
  },

  seeAlso: [
    { title: 'Rate limiting and abuse control', href: '/docs/systems/rate-limiting' },
    { title: 'Audit logging', href: '/docs/systems/audit-logging' },
    { title: 'Sentinel in Grit', href: '/docs/security/sentinel' },
  ],
}
