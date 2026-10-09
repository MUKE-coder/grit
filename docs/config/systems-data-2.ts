import type { SystemDesign } from './systems-types'

export const MIGRATIONS: SystemDesign = {
  slug: 'schema-migrations',
  name: 'Schema Migrations',
  tagline:
    'Changing the shape of a live database, and being able to undo a run in a framework that has no migration files.',
  group: 'Data',
  packages: ['internal/migrate', 'internal/models'],

  problem: {
    text: [
      'The schema changes constantly early on and keeps changing forever. A new column, a new table, an index somebody needed at three in the morning. Every one of those changes has to be applied to a database that already holds data, in an order, on every machine that runs the application, and in production exactly once.',
      'The usual answer is a directory of numbered migration files, each with an up and a down. It works, and it costs: every model change needs a hand-written SQL file that duplicates what the model already says, the two drift, and the down script is written by somebody who will never run it and is therefore usually wrong.',
      'Grit derives the schema from the models with GORM AutoMigrate instead, which removes the duplication and the drift. That creates a different problem. With no up script there is nothing to write a down script against, so a bad run looks unrecoverable. It is not, because of a fact about AutoMigrate that can be used: it only ever adds. It creates a table, adds a column, builds an index. It never drops a column. The reverse of a run does not have to be written, it can be computed.',
    ],
    capabilities: [
      'Apply the schema the models describe to a database at any prior state.',
      'Be safe to run twice: the second run is a no-op, not an error.',
      'Record what each run actually changed, rather than what it was asked to change.',
      'Undo a run by dropping exactly what that run added, newest first.',
      'Refuse to undo the run that built the schema in the first place.',
      'Never drop data without printing what will be lost and being told to go ahead.',
      'Work the same on SQLite, Postgres and MySQL.',
    ],
  },

  functional: [
    'A migrate command that applies the models to the database.',
    'Idempotent runs: nothing to do means nothing done and an exit code of zero.',
    'A schema snapshot taken before and after each run.',
    'Per-run history of the tables, columns and indexes added.',
    'A migrate down command that reverses a named run or the latest one.',
    'A dry run that prints the statements without issuing them.',
    'The first run against an empty database marked as a baseline, which down refuses by default.',
    'Seeding as a separate, repeatable command rather than part of migration.',
  ],

  nonFunctional: [
    {
      label: 'Idempotent',
      text: 'running twice does nothing the second time. Anything else makes deployment pipelines conditional, and conditional deployment steps get skipped.',
    },
    {
      label: 'Honest',
      text: 'the history records what the database reported changing, taken from before and after snapshots, not what the code intended.',
    },
    {
      label: 'Loud about loss',
      text: 'dropping a column loses every value written into it. The command prints each statement and will not proceed without an explicit flag.',
    },
    {
      label: 'Portable',
      text: 'three database engines, same commands. The snapshot reads from the information the driver exposes rather than from engine-specific SQL.',
    },
    {
      label: 'Ordered',
      text: 'a rollback reverses in the opposite order to application, because an index belongs to a column and a column to a table.',
    },
  ],

  capacity: {
    assumptions: [
      ['Tables in a mature generated project', '25 to 60'],
      ['Columns per table', '8 to 25'],
      ['Migration runs per month in active development', '20 to 100'],
      ['Changes per run', '0 to 15'],
      ['Rows in the largest table', 'up to tens of millions'],
    ],
    estimates: [
      {
        label: 'History growth',
        working: [
          '100 runs/month x 5 changes/run = 500 rows/month',
          'x 12 months = 6,000 rows/year',
          'at ~120 bytes/row = under 1 MB/year',
        ],
        note: 'The history is small enough that pruning it is not worth the code. Keeping every run forever is a feature: it is the record of how the schema got here.',
      },
      {
        label: 'Snapshot cost',
        working: [
          '60 tables, each read from the driver metadata',
          'two snapshots per run, one before and one after',
          'measured in tens of milliseconds, not seconds',
        ],
      },
      {
        label: 'The expensive change',
        working: [
          'adding a nullable column: metadata only, instant at any table size',
          'adding an index on 10,000,000 rows: minutes, and locking on some engines',
          'adding a not-null column with a default: a table rewrite on older engines',
        ],
        note: 'The cost is a property of the change, not of the tool. Which is why the dry run matters: it is the only chance to notice an index build before it happens on a live table.',
      },
    ],
  },

  highLevel: {
    intro:
      'Snapshot, migrate, snapshot, diff, record. The rollback reads the record and inverts each entry.',
    components: [
      {
        label: 'Model registry',
        text: 'the list of models to migrate. Generated resources append to it, so a new resource is migrated because it exists rather than because somebody remembered.',
      },
      {
        label: 'Schema snapshot',
        text: 'what the database currently has: tables, their columns, their indexes. Taken from driver metadata so one implementation covers three engines.',
      },
      {
        label: 'AutoMigrate',
        text: 'GORM applying the models. The property the whole design rests on: it adds and never drops.',
      },
      {
        label: 'Diff',
        text: 'after minus before. Each entry is a kind, a table and a name, and each kind has an exact inverse.',
      },
      {
        label: 'Run history',
        text: 'two tables. A run, with when it happened and which version of the framework did it, and its changes.',
      },
      {
        label: 'Rollback',
        text: 'reads a run and issues the inverse of each change in reverse order, after printing them all.',
      },
    ],
    flow: {
      title: 'A migration run, and the rollback it makes possible',
      nodes: [
        { id: 'cmd', label: 'grit migrate', col: 0, row: 0 },
        { id: 'before', label: 'Snapshot Before', col: 1, row: 0, accent: true },
        { id: 'auto', label: 'AutoMigrate', col: 2, row: 0 },
        { id: 'after', label: 'Snapshot After', col: 3, row: 0, accent: true },
        { id: 'diff', label: 'Diff', col: 3, row: 1, accent: true },
        { id: 'history', label: 'Run History', col: 2, row: 1 },
        { id: 'down', label: 'grit migrate down', col: 0, row: 2 },
        { id: 'invert', label: 'Invert + Confirm', col: 1, row: 2, accent: true },
        { id: 'db', label: 'Database', col: 2, row: 2 },
      ],
      edges: [
        { from: 'cmd', to: 'before', step: 1 },
        { from: 'before', to: 'auto', step: 2 },
        { from: 'auto', to: 'after', step: 3 },
        { from: 'after', to: 'diff', step: 4 },
        { from: 'diff', to: 'history', step: 5 },
        { from: 'down', to: 'invert', step: 6 },
        { from: 'invert', to: 'history', step: 7, bend: 'v' },
        { from: 'invert', to: 'db', step: 8, bend: 'v' },
      ],
      steps: [
        'The command runs, from the API directory in development or from anywhere as a built binary, loading the same configuration the server does.',
        'The current schema is read: every table, its columns, its indexes. An empty result here means this run is the baseline.',
        'AutoMigrate applies the models. It creates what is missing and leaves what exists, which is what makes a second run a no-op.',
        'The schema is read again, the same way.',
        'After minus before gives the changes. They are recorded against a run, with the time and the framework version, so the history says what happened rather than what was requested.',
        'Later, a rollback is asked for: the latest run, or one named by its identifier.',
        'The run is loaded with its changes. A run marked as the baseline is refused, because undoing the run that built the schema drops every table in it.',
        'Each change is printed with the statement that will reverse it. Nothing is issued until the operator confirms, because an added column that has been in production for a week has a week of data in it that the drop will take with it.',
      ],
    },
    dataFlow: [
      'The diff is computed from the database, not predicted from the models. A column AutoMigrate declined to add is not in the history, which is correct: the rollback should not try to drop it.',
      'Run identifiers are timestamps, pushed past the latest recorded run when the clock would otherwise collide, so two runs in the same second get distinct identifiers.',
      'The inverse order matters. An index is dropped before its column, a column before its table.',
      'Seeding is deliberately not part of this. A migration changes shape, a seed writes rows, and conflating them makes the migration non-repeatable.',
    ],
  },

  stack: [
    ['Applier', 'GORM AutoMigrate'],
    ['Snapshot', 'driver metadata, one implementation for all three engines'],
    ['History', 'two tables, migration_runs and migration_run_changes'],
    ['Reversal', 'computed from the diff, not hand-written'],
    ['Engines', 'SQLite, PostgreSQL, MySQL'],
    ['Commands', 'grit migrate, grit migrate down, grit seed'],
  ],

  dataModel: {
    entities: [
      {
        name: 'migration_runs',
        fields: [
          ['id', 'string, a timestamp, unique per run'],
          ['applied_at', 'when it ran'],
          ['grit_version', 'which framework version applied it'],
          ['baseline', 'true for the run that created the schema'],
        ],
        note: 'grit_version is here because a schema problem is often a framework upgrade, and knowing which version made a change narrows that in one query.',
      },
      {
        name: 'migration_run_changes',
        fields: [
          ['run_id', 'the run it belongs to'],
          ['kind', 'table, column or index'],
          ['on_table', 'the table it happened to'],
          ['name', 'the column or index name, empty for a table'],
        ],
        note: 'Three kinds, because AutoMigrate makes three kinds of change and each has exactly one inverse. A fourth kind would mean the design has a gap.',
      },
    ],
  },

  lowLevel: {
    classes: [
      {
        name: 'migrate.Change',
        file: 'internal/migrate/history.go',
        what: 'One thing a run added: a kind, a table, a name. Deliberately the smallest shape that can be inverted without interpretation.',
      },
      {
        name: 'migrate.Record',
        what: 'Writes a run and its changes in one transaction, with the framework version and whether it is the baseline.',
      },
      {
        name: 'Snapshot',
        what: 'Reads the schema into a comparable structure. The same function is called twice per run, so the before and after cannot be read differently.',
      },
      {
        name: 'nextRunID',
        what: 'A timestamp, moved past the latest recorded run when it would collide. Two migrations in the same second previously shared an identifier, and a rollback then reversed both.',
      },
    ],
    principles: [
      {
        label: 'Derive, do not duplicate',
        text: 'the model is the schema. A migration file restating it is a second source of truth, and two sources of truth is one more than is useful.',
      },
      {
        label: 'Record the effect, not the intent',
        text: 'taking the diff from the database means the history is true even when AutoMigrate did something other than what was expected.',
      },
      {
        label: 'Destructive operations ask',
        text: 'the rollback prints every statement and needs a flag. Making it easy to run would make it easy to run by accident.',
      },
    ],
    patterns: [
      ['Declarative schema', 'the models describe the target, the tool reaches it'],
      ['Computed inverse', 'the down migration derived from the observed diff'],
      ['Baseline marker', 'the one run that is not reversible, flagged rather than special-cased'],
      ['Idempotent apply', 'convergence to a target state rather than a sequence of steps'],
    ],
  },

  scaling: [
    'Migration is a deploy-time operation, so throughput is irrelevant and latency only matters because it holds up a deploy.',
    'The snapshot is proportional to the number of tables, not the number of rows, so it does not get slower as the data grows.',
    'The expensive changes are index builds and table rewrites on large tables, and those are properties of the change. The dry run is the place to catch them.',
    'Multiple instances starting at once can all try to migrate. Running the command as its own deploy step rather than at server boot avoids that entirely, and is what the generated pipelines do.',
    'The history grows by a few thousand rows a year, which never needs managing.',
    'AutoMigrate never dropping anything means an old instance rolling alongside a new one keeps working, because the columns it knows about are all still there.',
  ],

  bottlenecks: {
    problems: [
      {
        label: 'Changes AutoMigrate will not make',
        text: 'renaming a column, changing a type narrowly, adding a not-null column to a populated table. These need a deliberate step, and the tool not doing them silently is better than doing them wrongly.',
      },
      {
        label: 'Data migrations',
        text: 'a shape change often needs the existing rows reshaped too, and nothing about the schema diff knows that.',
      },
      {
        label: 'Concurrent migration at boot',
        text: 'several instances starting together can race, and on some engines the loser gets an error rather than a no-op.',
      },
      {
        label: 'Rollback loses data',
        text: 'computable does not mean safe. Dropping a column that has been live for a week drops a week of writes.',
      },
    ],
    improvements: [
      {
        label: 'Treat a rename as add, backfill, drop',
        text: 'three deliberate steps, each reversible, instead of one operation the tool cannot do. Slower and survivable.',
      },
      {
        label: 'Put data changes in a seed or a job',
        text: 'they are repeatable, idempotent code with tests, which a migration file is not a good place for.',
      },
      {
        label: 'Migrate as a deploy step',
        text: 'one process, before the new instances start. This is what the generated pipelines and Dockerfiles do, and it removes the race rather than handling it.',
      },
      {
        label: 'Snapshot before rolling back',
        text: 'a database backup immediately before a down run turns an unrecoverable mistake into a restore.',
      },
    ],
  },

  seeAlso: [
    { title: 'Migrations reference', href: '/docs/database/migrations' },
    { title: 'Code generation', href: '/docs/systems/code-generation' },
  ],
}

export const PAGINATION: SystemDesign = {
  slug: 'pagination-and-search',
  name: 'Pagination and Search',
  tagline:
    'Returning a slice of a large table with a total, a sort, a filter and a search box, without the query getting slower as the table grows.',
  group: 'Data',
  packages: ['internal/paginate', 'internal/services'],

  problem: {
    text: [
      'Every list endpoint in the application is the same endpoint. Give me page three of the orders, sorted by date, filtered to the ones that are unpaid, matching the text someone typed, and tell me how many there are in total so I can draw the page numbers.',
      'Written by hand per resource, that becomes forty slightly different implementations. Some validate the sort column and some interpolate it into the SQL. Some cap the page size and some let a caller ask for a million rows. Some count with the filters applied and some count the whole table, which makes the page count wrong. The variations are invisible until one of them is the one with the injection.',
      'There is also a limit nobody hits until they do. Offset pagination asks the database to produce and discard every row before the one you want, so page five thousand reads five thousand pages worth of rows to return twenty. The fix is keyset pagination, which is strictly better for scrolling and cannot produce page numbers, so both have to exist and the choice has to be per request.',
    ],
    capabilities: [
      'Bind page, size, sort, order, search and filters from the query string once, for every resource.',
      'Validate the sort and filter columns against a whitelist, never interpolate what a caller sent.',
      'Cap the page size, so no request can ask for the whole table.',
      'Return the total and the page count, computed over the same filters as the page.',
      'Offer keyset pagination for deep lists, where offset stops being viable.',
      'Answer the extra counts a dashboard needs on the same request rather than one request per card.',
      'Search across the columns the resource says are searchable.',
    ],
  },

  functional: [
    'Default page 1, default size 20, maximum size 100.',
    'Default sort created_at descending, overridable per resource.',
    'Sort column checked against the resource whitelist; anything else falls back to the default.',
    'Filters from query parameters, applied only for whitelisted columns.',
    'An IN filter for id columns, opted into per column rather than inferred.',
    'Date range filtering on a named column.',
    'Full text search across declared columns.',
    'Cursor mode, requested explicitly or implied by sending a cursor.',
    'Extra counts, a time series and a per-value breakdown, all on the list request.',
  ],

  nonFunctional: [
    {
      label: 'No interpolation, ever',
      text: 'a sort column is the one place a list endpoint cannot use a bound parameter, which makes it the one place a whitelist is mandatory rather than advisable.',
    },
    {
      label: 'Bounded by construction',
      text: 'the maximum page size is enforced in the binder, not asked of each handler. A handler cannot forget it.',
    },
    {
      label: 'Consistent envelope',
      text: 'every list returns data and meta in the same shape, so one client-side hook works for every resource.',
    },
    {
      label: 'Zero is an answer',
      text: 'the counts are serialised without omitempty. An empty table reporting no total at all leaves every client doing arithmetic on undefined.',
    },
    {
      label: 'Honest about mode',
      text: 'a cursor response says so, because total, page and pages are all zero in cursor mode and a client cannot otherwise distinguish that from an empty table.',
    },
  ],

  capacity: {
    assumptions: [
      ['List requests per second', '~550 at peak'],
      ['Default page size', '20'],
      ['Rows in a large table', '10,000,000'],
      ['Indexed sort column', 'created_at'],
      ['Share of requests on page 1', '~85%'],
    ],
    estimates: [
      {
        label: 'Offset cost by depth',
        working: [
          'page 1: LIMIT 20 OFFSET 0, reads 20 rows',
          'page 100: LIMIT 20 OFFSET 1,980, reads 2,000 rows',
          'page 10,000: LIMIT 20 OFFSET 199,980, reads 200,000 rows',
        ],
        note: 'The cost is linear in the page number. Fine for an admin table nobody pages past thirty, fatal for an export loop or an infinite scroll.',
      },
      {
        label: 'Keyset cost by depth',
        working: [
          'WHERE (created_at, id) < (?, ?) ORDER BY created_at DESC, id DESC LIMIT 20',
          'an index seek plus 20 rows, at any depth',
          'page 10,000 costs what page 1 costs',
        ],
      },
      {
        label: 'The count is the expensive part',
        working: [
          'COUNT(*) with filters on 10,000,000 rows: hundreds of milliseconds',
          'the page itself: single-digit milliseconds on an indexed sort',
          'so the total, not the page, dominates a list request',
        ],
        note: 'Which is why cursor mode does not count unless asked, and why the dashboard counts are opt-in per request rather than always computed.',
      },
      {
        label: 'Payload',
        working: [
          '20 rows x ~1.5 KB = 30 KB per page',
          '550 requests/s x 30 KB = ~16 MB/s before compression',
          'with gzip, roughly 4 MB/s',
        ],
      },
    ],
  },

  highLevel: {
    intro:
      'One binder turns the query string into a validated structure. One list function turns that structure plus a whitelist into a query and an envelope.',
    components: [
      {
        label: 'Params',
        text: 'the normalised query state. Produced from the request once, with everything already clamped to its limits.',
      },
      {
        label: 'Config',
        text: 'what this resource allows: which columns are sortable, filterable and searchable, and what the default sort is. The whitelist lives with the resource, not in the binder.',
      },
      {
        label: 'Query builder',
        text: 'applies search, filters, the date range and the sort, then either an offset and a limit or a keyset predicate.',
      },
      {
        label: 'Counter',
        text: 'the total over the same filters as the page, plus whatever extra counts were asked for.',
      },
      {
        label: 'Cursor codec',
        text: 'encodes the sort value and the id of the last row as an opaque token, so pages stay stable when rows are inserted mid-pagination.',
      },
      {
        label: 'Meta',
        text: 'the envelope. Page and pages in offset mode, next cursor and has-more in cursor mode, and the mode itself so the client can tell.',
      },
    ],
    flow: {
      title: 'A list request, from query string to envelope',
      nodes: [
        { id: 'client', label: 'Client', col: 0, row: 0 },
        { id: 'bind', label: 'Bind Params', col: 1, row: 0, accent: true },
        { id: 'config', label: 'Resource Config', col: 2, row: 0 },
        { id: 'build', label: 'Build Query', col: 3, row: 0, accent: true },
        { id: 'scope', label: 'Ownership Scope', col: 2, row: 1, accent: true },
        { id: 'db', label: 'Database', col: 3, row: 1 },
        { id: 'count', label: 'Count + Extras', col: 3, row: 2 },
        { id: 'meta', label: 'Meta Envelope', col: 1, row: 2, accent: true },
        { id: 'resp', label: 'data + meta', col: 0, row: 2 },
      ],
      edges: [
        { from: 'client', to: 'bind', step: 1 },
        { from: 'bind', to: 'config', step: 2 },
        { from: 'config', to: 'build', step: 3 },
        { from: 'build', to: 'scope', step: 4, bend: 'v' },
        { from: 'scope', to: 'db', step: 5 },
        { from: 'db', to: 'count', step: 6 },
        { from: 'count', to: 'meta', step: 7, bend: 'h' },
        { from: 'meta', to: 'resp', step: 8 },
      ],
      steps: [
        'A list request arrives with page, size, sort, order, search and whatever filters the interface offers.',
        'Everything is bound and clamped. A size of 5,000 becomes 100, a page of minus one becomes 1, and the reserved words are separated from the arbitrary query parameters so the whitelist can be applied to the arbitrary ones and only those.',
        'The resource config says what is allowed. A sort column not on the list is not an error, it is the default sort, because a stale bookmark should return rows rather than a 400.',
        'The query is assembled: the search across the searchable columns, the whitelisted filters, the date range, the validated sort.',
        'The ownership scope is applied last, by the same code the single-record read uses. A list is the easiest place to leak another tenant’s rows, so this is not something the list builder is trusted to remember.',
        'The page is read. In offset mode that is a limit and an offset; in cursor mode a keyset predicate on the sort value and the id.',
        'The total is counted over the same filters, and any extra counts, series or breakdowns the request asked for are computed alongside it rather than as separate requests.',
        'The envelope is assembled. Offset mode gets page and pages, cursor mode gets a next cursor and has-more and says that it is cursor mode, and both get the data under the same key.',
      ],
    },
    dataFlow: [
      'The count uses the same filters and search as the page. Counting the unfiltered table is the commonest way to get a page count wrong, and it is wrong in the direction that shows empty pages.',
      'The cursor encodes the sort value and the id together, because a sort value is not unique and a cursor on the value alone either skips or repeats rows.',
      'Search is applied to declared columns only. Searching everything means searching columns with indexes nobody built.',
      'The IN filter splits on commas and is therefore only enabled for id columns, where a comma never appears in a value. For anything a person types, "Smith, John" is one value.',
    ],
  },

  stack: [
    ['Binder', 'one function, every list endpoint'],
    ['Defaults', 'page 1, size 20, maximum 100, sort created_at desc'],
    ['Modes', 'offset by default, keyset on request'],
    ['Safety', 'column whitelists per resource, bound parameters everywhere else'],
    ['Envelope', 'data plus meta, identical across resources'],
    ['Extras', 'counts, series and breakdown on the same request'],
  ],

  api: {
    groups: [
      {
        title: 'The list surface, identical for every resource',
        rows: [
          { method: 'GET', path: '/api/v1/orders?page=2&page_size=50', what: 'Offset pagination, capped at 100' },
          { method: 'GET', path: '/api/v1/orders?sort_by=total&sort_order=asc', what: 'Sort, validated against the whitelist' },
          { method: 'GET', path: '/api/v1/orders?search=invoice', what: 'Search across the declared columns' },
          { method: 'GET', path: '/api/v1/orders?status=unpaid', what: 'Filter, whitelisted per column' },
          { method: 'GET', path: '/api/v1/orders?mode=cursor', what: 'Keyset pagination, flat cost at any depth' },
          { method: 'GET', path: '/api/v1/orders?counts=created_7d', what: 'Extra totals beside the page' },
        ],
      },
    ],
    samples: [
      {
        title: 'An offset response',
        language: 'json',
        code: `{
  "data": [ { "id": "01a1...", "total": 1250 } ],
  "meta": {
    "total": 4821,
    "page": 2,
    "page_size": 50,
    "pages": 97,
    "counts": { "created_7d": 112 }
  }
}`,
      },
      {
        title: 'A cursor response, which has no page numbers on purpose',
        language: 'json',
        code: `{
  "data": [ { "id": "01a1...", "total": 1250 } ],
  "meta": {
    "total": 0,
    "page": 0,
    "page_size": 20,
    "pages": 0,
    "mode": "cursor",
    "next_cursor": "eyJ2IjoiMjAyNi0xMC0wOSIsImkiOiIwMWEx...",
    "has_more": true
  }
}`,
      },
    ],
  },

  lowLevel: {
    classes: [
      {
        name: 'paginate.Params',
        file: 'internal/paginate/paginate.go',
        what: 'The bound query state. Reserved words are parsed into fields; everything else goes into a separate map so the whitelist applies to untrusted input and only to untrusted input.',
      },
      {
        name: 'paginate.Config',
        what: 'Per resource: sortable, filterable, searchable, the defaults, and whether cursor mode and the total are available.',
      },
      {
        name: 'paginate.List',
        what: 'Generic over the row type. Takes a query, params and config, returns the typed result and the meta.',
        methods: ['List'],
      },
      {
        name: 'paginate.Meta',
        what: 'The envelope. The counts have no omitempty, deliberately: a missing total is worse than a zero one.',
      },
    ],
    principles: [
      {
        label: 'Validate in one place',
        text: 'forty handlers each checking a sort column means thirty-nine chances to get it right and one not to.',
      },
      {
        label: 'Fall back rather than refuse',
        text: 'an unknown sort column returns the default order. A 400 for a stale bookmark is correct and useless.',
      },
      {
        label: 'Make the expensive thing opt-in',
        text: 'the count dominates the request, so cursor mode does not count and the dashboard extras are asked for explicitly.',
      },
      {
        label: 'Say which mode answered',
        text: 'the zeroes in a cursor response are indistinguishable from an empty table without it.',
      },
    ],
    patterns: [
      ['Offset pagination', 'page numbers and a total, linear cost in depth'],
      ['Keyset pagination', 'an opaque cursor, flat cost, no page numbers'],
      ['Whitelist validation', 'for the one input that cannot be a bound parameter'],
      ['Query object', 'one normalised structure instead of a dozen request fields'],
    ],
  },

  scaling: [
    'An indexed sort column is the single thing that makes a list endpoint fast. Without one, every page is a sort of the whole filtered set.',
    'Offset cost grows with the page number. That is acceptable for a human browsing and unacceptable for a loop, which is what cursor mode is for.',
    'The count is usually the slowest part of a list request, so the cheapest optimisation available is not computing it when nobody is showing it.',
    'The composite index that matters is the sort column plus the id, which is exactly what the keyset predicate uses.',
    'Capping the page size at 100 bounds the worst request anybody can send, which bounds memory and bandwidth per request by construction.',
    'List responses cache well because they are read far more than written, and prefix invalidation on write keeps them honest.',
  ],

  bottlenecks: {
    problems: [
      {
        label: 'Deep offset',
        text: 'page five thousand reads a hundred thousand rows to return twenty, and an export loop walks every page.',
      },
      {
        label: 'Counting on every request',
        text: 'a filtered count on a large table can cost more than the page by an order of magnitude, on every single request.',
      },
      {
        label: 'Unstable pages',
        text: 'rows inserted while somebody pages through offset results shift everything, so a row can be seen twice or missed.',
      },
      {
        label: 'Search without an index',
        text: 'a leading-wildcard LIKE across several columns is a full scan, and it looks fine on a thousand rows.',
      },
    ],
    improvements: [
      {
        label: 'Cursor mode for anything that walks',
        text: 'exports, infinite scroll, sync. Flat cost and stable pages, at the price of page numbers nobody was showing.',
      },
      {
        label: 'Cache or skip the total',
        text: 'an approximate total from statistics, a cached total on a short lifetime, or no total at all where the interface does not show one.',
      },
      {
        label: 'Index the sort key with the id',
        text: 'the composite index serves both the offset sort and the keyset predicate, and makes the tie-break deterministic.',
      },
      {
        label: 'Move search to the engine built for it',
        text: 'a Postgres full text index on a generated column, or a dedicated search service once the table is large enough that LIKE is the bottleneck.',
      },
    ],
  },

  seeAlso: [
    { title: 'Pagination reference', href: '/docs/api/pagination' },
    { title: 'Data isolation', href: '/docs/systems/data-isolation' },
  ],
}

export const ENCRYPTION: SystemDesign = {
  slug: 'encryption-at-rest',
  name: 'Field Encryption at Rest',
  tagline:
    'Making a column unreadable without the key, while the handler and the model carry on treating it as a string.',
  group: 'Data',
  packages: ['internal/crypto', 'internal/models'],

  problem: {
    text: [
      'Some columns hold things you must store and display but should never be legible straight out of the database. Personal notes on a client, a third-party token, a contact detail, an identity number. Database-level encryption protects the files on disk, which guards against a stolen volume and nothing else: anybody with a connection sees plaintext, and so does every backup dump, every replica and every analyst who was given read access for an afternoon in 2024.',
      'Encrypting those columns in the application closes that. The cost is usually that the model gets noisy: an encrypt call before every write, a decrypt after every read, a bug the first time somebody adds a field and forgets one, and a schema that can never adopt encryption later without a migration.',
      'A custom column type removes all of that. The field is declared as an encrypted string, the driver encrypts on the way in and decrypts on the way out, and nothing else in the application knows. What remains is one real constraint that cannot be designed away, and the honest thing is to state it rather than hide it.',
    ],
    capabilities: [
      'Encrypt a named field transparently, with no calls in the handler or the service.',
      'Use authenticated encryption, so a tampered value fails rather than decrypting to nonsense.',
      'Use a fresh nonce per write, so identical values do not produce identical ciphertext.',
      'Version the scheme, so the algorithm can change later without ambiguity about existing rows.',
      'Pass values through unencrypted when no key is configured, so the feature can be adopted without a migration.',
      'Refuse to start on a key of the wrong length, rather than quietly running without encryption.',
      'Encrypt the rows that are already there, in batches, when the feature is turned on.',
    ],
  },

  functional: [
    'An EncryptedString type usable as a model field.',
    'AES-256-GCM with a random 12-byte nonce per write.',
    'A version prefix on every stored value.',
    'A process-wide key loaded from configuration as base64 that decodes to 32 bytes.',
    'An empty key disabling encryption, with values stored as plaintext.',
    'A non-empty key of the wrong length failing at startup.',
    'A backfill that encrypts existing plaintext in a column, in batches.',
    'JSON serialisation of the plaintext, so an API response is unaffected.',
  ],

  nonFunctional: [
    {
      label: 'Authenticated, not just encrypted',
      text: 'GCM means a modified ciphertext fails to open. Encryption without authentication protects confidentiality and not integrity, and integrity is what you need when the threat is somebody with write access to the database.',
    },
    {
      label: 'Non-deterministic',
      text: 'a new nonce per write means the same value stored twice looks different. That is required for confidentiality and it is exactly why such a column cannot be queried by equality.',
    },
    {
      label: 'Invisible to the application',
      text: 'the handler, the service and the Zod schema see a string. If using it required thinking about it, some field somewhere would be missed.',
    },
    {
      label: 'Adoptable later',
      text: 'no key means passthrough, so turning the feature on is a configuration change and a backfill rather than a schema migration.',
    },
    {
      label: 'Loud on misconfiguration',
      text: 'a 16-byte key where 32 is needed stops the process. Running without the encryption that was asked for is the worse outcome.',
    },
  ],

  capacity: {
    assumptions: [
      ['Encrypted fields per project', '2 to 10'],
      ['Reads touching an encrypted field per second', '~200 at peak'],
      ['Average plaintext length', '200 bytes'],
      ['AES-GCM throughput with hardware support', '>1 GB/s per core'],
      ['Overhead per value', '12-byte nonce + 16-byte tag + prefix, then base64'],
    ],
    estimates: [
      {
        label: 'CPU',
        working: [
          '200 reads/s x 200 bytes = 40 KB/s decrypted',
          'against >1 GB/s available per core',
          'under 0.01% of one core',
        ],
        note: 'With AES-NI the cost is not measurable against a database round trip. The reason not to encrypt a column is never performance, it is that you need to query it.',
      },
      {
        label: 'Storage',
        working: [
          '200 bytes plaintext + 28 bytes nonce and tag = 228 bytes',
          'base64 encoded: 304 bytes',
          'plus the 7-byte prefix: ~311 bytes, about 1.55x',
        ],
      },
      {
        label: 'Backfill',
        working: [
          '1,000,000 rows, batches of 500',
          '2,000 batches, each a read, an encrypt and an update',
          'at ~50 ms per batch: about 100 seconds',
        ],
        note: 'Batched because a single update statement over a million rows holds a transaction open for the duration, and because a backfill that can be resumed is better than one that must complete.',
      },
    ],
  },

  highLevel: {
    intro:
      'A custom column type sitting in the driver interface, so the encryption happens between the model and the database and nowhere else.',
    components: [
      {
        label: 'EncryptedString',
        text: 'the field type. Implements the driver value and scan interfaces, which is where the encryption and decryption live.',
      },
      {
        label: 'Key holder',
        text: 'the process-wide key, loaded once at startup behind a lock. Nil means disabled, which is a supported state rather than an error.',
      },
      {
        label: 'Versioned envelope',
        text: 'the prefix that tags a value as encrypted and says which scheme wrote it. A value without the prefix is plaintext, which is how passthrough and partial backfill coexist.',
      },
      {
        label: 'AEAD',
        text: 'AES-256-GCM. One construction, chosen once, rather than a configurable cipher suite nobody is qualified to configure per project.',
      },
      {
        label: 'Backfill',
        text: 'walks a column in batches, encrypting anything without the prefix and leaving anything with it alone, so it is safe to run twice.',
      },
      {
        label: 'JSON hook',
        text: 'serialises the plaintext, so an encrypted field is an ordinary string in an API response and in the generated TypeScript type.',
      },
    ],
    flow: {
      title: 'A write and a read through an encrypted column',
      nodes: [
        { id: 'handler', label: 'Handler', col: 0, row: 0 },
        { id: 'model', label: 'Model Field', col: 1, row: 0 },
        { id: 'value', label: 'Value()', col: 2, row: 0, accent: true },
        { id: 'aead', label: 'AES-256-GCM', col: 3, row: 0, accent: true },
        { id: 'key', label: 'Field Key', col: 3, row: 1 },
        { id: 'col', label: 'Column: enc:v1:...', col: 2, row: 1 },
        { id: 'scan', label: 'Scan()', col: 1, row: 1, accent: true },
        { id: 'json', label: 'JSON: plaintext', col: 0, row: 1 },
        { id: 'backfill', label: 'Backfill', col: 3, row: 2 },
      ],
      edges: [
        { from: 'handler', to: 'model', step: 1 },
        { from: 'model', to: 'value', step: 2 },
        { from: 'value', to: 'aead', step: 3 },
        { from: 'key', to: 'aead', step: 4 },
        { from: 'aead', to: 'col', step: 5, bend: 'v' },
        { from: 'col', to: 'scan', step: 6 },
        { from: 'scan', to: 'json', step: 7 },
        { from: 'backfill', to: 'col', step: 8, bend: 'v' },
      ],
      steps: [
        'A handler binds a request and sets the field. It is a string type, so nothing here is different from any other field.',
        'The row is saved. GORM asks the field for its database value, which is the hook the whole design hangs on.',
        'A fresh 12-byte nonce is generated for this write. Reusing a nonce with the same key is the one catastrophic mistake available in GCM, so it is never derived from the value or a counter.',
        'The key is read from the process-wide holder. A nil key means encryption is disabled and the value is returned as it is, which is what lets a project adopt this later.',
        'The nonce, the ciphertext and the authentication tag are concatenated, base64 encoded and prefixed with the scheme version. The column now holds something opaque, and the prefix is what will let a future version rotate the algorithm without guessing how an old row was written.',
        'On a read, the column value comes back. Without the prefix it is plaintext from before the backfill and is used as it is; with the prefix it is opened, and a failure to open is an error rather than a fallback, because a value that will not authenticate has been tampered with or is encrypted under a different key.',
        'The field is plaintext from here on. The API response, the generated TypeScript type and the admin form all see an ordinary string.',
        'Separately, turning the feature on runs a backfill over the existing rows, in batches, skipping anything already prefixed so it can be rerun after an interruption.',
      ],
    },
    dataFlow: [
      'Encryption happens in the driver interface, which means it happens for every writer: the handler, a job, a seed, a bulk import. There is no path to the column that bypasses it.',
      'The prefix makes the column self-describing. Plaintext and ciphertext can coexist in it during a backfill, which is what makes the backfill resumable.',
      'JSON carries plaintext deliberately. Encrypting the API response would mean the client needs the key, which would mean the key is on the client.',
      'An encrypted column cannot be sorted, searched or filtered by equality, because the ciphertext for the same value differs every time. That is the price, and it is the reason this is for data you keep and show rather than for keys and lookup columns.',
    ],
  },

  stack: [
    ['Cipher', 'AES-256-GCM'],
    ['Nonce', '12 bytes, fresh per write, from crypto/rand'],
    ['Key', '32 bytes, base64 in configuration, process-wide'],
    ['Envelope', 'enc:v1: prefix, then base64'],
    ['Integration', 'driver.Valuer and sql.Scanner on a string type'],
    ['Absent key', 'passthrough, so adoption needs no migration'],
  ],

  dataModel: {
    intro:
      'There is no new table. The design is entirely in the column type, which is the point.',
    entities: [
      {
        name: 'any model field',
        fields: [
          ['notes', 'crypto.EncryptedString, stored as text'],
          ['stored form', 'enc:v1:<base64 of nonce + ciphertext + tag>'],
          ['size', 'roughly 1.55x the plaintext'],
        ],
        note: 'Text rather than bytes, so the column is readable in any client and the prefix is visible when somebody is working out what they are looking at.',
      },
    ],
    storage: [
      'The key never goes in the database. It comes from the environment or a secret manager, which is what separates it from the ciphertext it protects.',
      'Backups contain ciphertext. A leaked dump without the key is not a breach of those columns, which is the whole objective.',
      'A replica or an analytics connection sees ciphertext too, because the encryption is above the database rather than inside it.',
    ],
  },

  lowLevel: {
    classes: [
      {
        name: 'crypto.EncryptedString',
        file: 'internal/crypto/field.go',
        what: 'A string type with Value and Scan. Everything else about it is ordinary, which is why a model using it reads like a model that does not.',
        methods: ['Value', 'Scan', 'MarshalJSON', 'UnmarshalJSON'],
      },
      {
        name: 'crypto.InitFieldKey',
        what: 'Loads the base64 key once at startup. Empty disables encryption; non-empty and the wrong length is a hard error.',
      },
      {
        name: 'crypto.Encrypt / Decrypt',
        what: 'The envelope functions. Decrypt returns a value without the prefix unchanged, which is what makes plaintext and ciphertext coexist during a backfill.',
      },
      {
        name: 'crypto.EncryptExisting',
        what: 'The batched backfill. Skips prefixed values, so it is idempotent and resumable.',
      },
    ],
    principles: [
      {
        label: 'Put it where nothing can bypass it',
        text: 'in the driver interface, every writer is covered. In a service method, only the writers that call that method are.',
      },
      {
        label: 'One construction, not a choice',
        text: 'AES-256-GCM, fixed. A configurable cipher suite is an invitation to configure it wrongly in a project that has no cryptographer.',
      },
      {
        label: 'Version from the first write',
        text: 'the prefix costs seven bytes and is the only thing that makes rotating the scheme later possible rather than archaeological.',
      },
      {
        label: 'Fail loudly on the key, quietly on its absence',
        text: 'no key is a decision. A malformed key is a mistake, and the difference should be the difference between starting and not.',
      },
    ],
    patterns: [
      ['Transparent column type', 'encryption in Value and Scan rather than at the call sites'],
      ['Versioned envelope', 'a scheme tag on every stored value'],
      ['AEAD', 'confidentiality and integrity from one construction'],
      ['Idempotent backfill', 'skip what is already done, so it can be rerun'],
    ],
  },

  scaling: [
    'AES-GCM with hardware support is far faster than the database round trip that delivered the row, so the CPU cost does not appear in a profile.',
    'Storage grows by about 1.55x on the encrypted columns only, which is immaterial unless the column is large.',
    'The real scaling constraint is not throughput, it is that the column is unqueryable. A field that needs sorting, searching or an equality lookup cannot be encrypted this way, and that has to be decided before the data exists.',
    'The backfill is batched, so it runs against a live table without holding a long transaction, and can be interrupted and resumed.',
    'A key rotation is a second backfill: decrypt under the old key, encrypt under the new, with the version prefix telling the two apart.',
    'No new table and no new service means nothing extra to scale, replicate or operate.',
  ],

  bottlenecks: {
    problems: [
      {
        label: 'The column cannot be queried',
        text: 'non-deterministic ciphertext rules out equality, ordering, prefix search and unique constraints. This is the constraint that actually bites, and it bites after the data exists.',
      },
      {
        label: 'Key loss is data loss',
        text: 'there is no recovery. A lost key makes every encrypted column permanently unreadable, and that is the feature working as designed.',
      },
      {
        label: 'The key in the wrong place',
        text: 'a key committed to the repository or sitting in the same backup as the database provides no protection at all while looking exactly like protection.',
      },
      {
        label: 'Mixed state after a partial backfill',
        text: 'plaintext and ciphertext in one column is correct during the backfill and misleading if it is never finished.',
      },
    ],
    improvements: [
      {
        label: 'Keep a blind index for lookups',
        text: 'a deterministic HMAC of the normalised value in a second column gives exact-match lookup while the encrypted column stays non-deterministic.',
      },
      {
        label: 'Use a managed key service',
        text: 'KMS or Vault, with the field key wrapped rather than stored. Rotation and audit come with it, and the key stops being an environment variable somebody can print.',
      },
      {
        label: 'Separate the key from every backup',
        text: 'different system, different credentials, different retention. If the dump and the key travel together, nothing was gained.',
      },
      {
        label: 'Finish the backfill and then verify',
        text: 'a count of unprefixed values in the column is one query and is the only way to know the state is not mixed.',
      },
    ],
  },

  seeAlso: [
    { title: 'Field encryption reference', href: '/docs/security/field-encryption' },
    { title: 'Audit logging', href: '/docs/systems/audit-logging' },
  ],
}
