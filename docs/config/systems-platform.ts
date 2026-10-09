import type { SystemDesign } from './systems-types'

export const CODE_GENERATION: SystemDesign = {
  slug: 'code-generation',
  name: 'Code Generation',
  tagline:
    'Turning one field list into a model, a service, a handler, routes, a schema, types, hooks and an admin page, and being able to take all of it back out.',
  group: 'Platform',
  packages: ['internal/generate', 'internal/scaffold'],

  problem: {
    text: [
      'A new resource in a full-stack application is about eleven files. A Go model, a service, a handler, a route group, a Zod schema, a TypeScript type, a React Query hook, an admin resource definition, an admin page, a sidebar entry, a CSV importer. None of them is hard, all of them are the same every time, and writing them by hand takes an afternoon and produces four inconsistencies.',
      'The inconsistencies are the real cost rather than the afternoon. One resource paginates with a cap and the next does not. One scopes its list to the owner and the next forgets, which is a data leak rather than a style difference. One validates with Zod and the next trusts the form. Nobody notices, because each file is individually reasonable and the divergence is only visible across the set.',
      'A generator removes both costs, and introduces its own problem. It does not only write files: it has to inject into files that already exist. The route file gains a group, the sidebar gains a link, the shared package gains a schema, the admin registry gains an entry. Injection is where a generator stops being a template engine and becomes something that edits a codebase, and edits can be wrong in ways that do not fail loudly.',
      'Two specific ways, both of which happened here. An injection performed with a string replace against an anchor that has moved matches nothing, returns the input unchanged, and reports success. And every injection that has no matching removal leaves the project uncompilable the moment somebody removes the resource: a dangling import, a sidebar link to a deleted page.',
    ],
    capabilities: [
      'Take a resource name and a field list and produce the whole vertical slice.',
      'Inject into existing files at named anchors, idempotently.',
      'Remove a resource completely, undoing every injection.',
      'Map a Go field type to a GORM tag, a Zod validator and a TypeScript type, consistently.',
      'Regenerate TypeScript types and Zod schemas from the Go structs on demand.',
      'Add a field to an existing resource without regenerating it.',
      'Honour flags that change the shape: owned by user, append-only, public read, role-restricted.',
      'Write the same code whether a project is new or being upgraded.',
    ],
  },

  functional: [
    'grit generate resource with an inline field list or a definition file.',
    'Thirty field types, each with its Go type, GORM tag, Zod validator, TypeScript type and admin form control.',
    'Anchor comments in the scaffolded templates, which the generator injects before.',
    'grit remove, undoing every injection the generator makes.',
    'grit sync, regenerating TypeScript and Zod from the Go structs.',
    'grit add field, for a field on a resource that already exists.',
    'Generators for jobs, mail, factories and column packs.',
    'Flags: owned by user, append-only, public, role restriction, soft delete.',
    'A CSV importer per resource, resolving relationships by natural key.',
  ],

  nonFunctional: [
    {
      label: 'Idempotent',
      text: 'generating the same resource twice must not produce two routes. Every injection checks for its own output first.',
    },
    {
      label: 'Reversible',
      text: 'generate and remove are a pair. An injection added without its removal is a bug, because it makes a cleanup operation break the build.',
    },
    {
      label: 'Fails loudly',
      text: 'an import block is built as a list, never with a string replace on an anchor line, because a replace whose anchor moved matches nothing and reports success.',
    },
    {
      label: 'Output is idiomatic',
      text: 'generated code is read, edited and reviewed by people. It is formatted, commented and ordinary, not marked as untouchable.',
    },
    {
      label: 'One source for shared text',
      text: 'anything written both at generate time and at upgrade time is written by one function called from both, or the two drift.',
    },
    {
      label: 'Types stay in step',
      text: 'the response shape, the Zod schema and the TypeScript type change in one commit. That is what sync is for.',
    },
  ],

  capacity: {
    assumptions: [
      ['Files written per resource', '~11'],
      ['Files injected into per resource', '~6'],
      ['Field types supported', '30'],
      ['Resources in a mature project', '25 to 60'],
      ['Generation time', 'under a second'],
    ],
    estimates: [
      {
        label: 'What one command replaces',
        working: [
          '~11 files, ~900 lines of ordinary code',
          'plus 6 injections across existing files',
          'by hand: an afternoon, and four inconsistencies',
        ],
      },
      {
        label: 'Injection surface',
        working: [
          '6 anchors x 60 resources = 360 injections in a mature project',
          'each of which must be removable',
          'each of which must be idempotent',
        ],
        note: 'This is the number that makes removal a first-class requirement rather than a nicety. Three hundred and sixty edits that only go one way is a codebase that can only grow.',
      },
      {
        label: 'Type mapping combinations',
        working: [
          '30 field types x 5 outputs (GORM, Zod, TS, form control, importer)',
          '= 150 mappings',
        ],
        note: 'Which is why the rule is to generate one resource with all thirty types and compile it. A mapping that is wrong for money or json is invisible in a resource that uses neither.',
      },
      {
        label: 'Test surface',
        working: [
          '124 unit and integration tests across generate and scaffold',
          'covering pluralisation, type mapping, injection, removal and round trips',
        ],
      },
    ],
  },

  highLevel: {
    intro:
      'Parse a definition, write files from templates, inject at anchors, and keep a removal for every injection.',
    components: [
      {
        label: 'Definition parser',
        text: 'turns a field list, inline or from a file, into a typed description: names in every case the templates need, types, tags, relationships and flags.',
      },
      {
        label: 'Type mapper',
        text: 'one field type to its Go type, GORM tag, Zod validator, TypeScript type, admin control and importer handling. One table, because five tables drift.',
      },
      {
        label: 'Template writers',
        text: 'the files produced whole: the model, the service, the handler, the schema, the types, the hooks, the admin definition and page.',
      },
      {
        label: 'Injectors',
        text: 'inject before an anchor, or inline. The anchors are comments the scaffold emits, which makes them a contract between the two packages.',
      },
      {
        label: 'Removers',
        text: 'one per injection. Remove a line containing, remove inline text, remove a line block, remove a schema export block.',
      },
      {
        label: 'Sync',
        text: 'parses the Go structs and regenerates the TypeScript types and Zod schemas, so the three never need to be edited in step by hand.',
      },
      {
        label: 'Shared hunk functions',
        text: 'text needed at both generate time and upgrade time, written by one function called from both so a new project and an upgraded one cannot differ.',
      },
    ],
    flow: {
      title: 'One command, eleven files and six injections',
      nodes: [
        { id: 'cmd', label: 'grit generate', col: 0, row: 0 },
        { id: 'parse', label: 'Parse Definition', col: 1, row: 0, accent: true },
        { id: 'types', label: 'Type Mapper', col: 2, row: 0, accent: true },
        { id: 'write', label: 'Write Files', col: 3, row: 0 },
        { id: 'anchors', label: 'Find Anchors', col: 2, row: 1, accent: true },
        { id: 'inject', label: 'Inject', col: 3, row: 1, accent: true },
        { id: 'remove', label: 'grit remove', col: 0, row: 2 },
        { id: 'undo', label: 'Matching Removals', col: 1, row: 2, accent: true },
        { id: 'project', label: 'Project Compiles', col: 2, row: 2 },
      ],
      edges: [
        { from: 'cmd', to: 'parse', step: 1 },
        { from: 'parse', to: 'types', step: 2 },
        { from: 'types', to: 'write', step: 3 },
        { from: 'write', to: 'anchors', step: 4, bend: 'v' },
        { from: 'anchors', to: 'inject', step: 5 },
        { from: 'remove', to: 'undo', step: 6 },
        { from: 'undo', to: 'project', step: 7 },
        { from: 'inject', to: 'project', step: 8, bend: 'v' },
      ],
      steps: [
        'A field list arrives: a resource name and its fields, with types and modifiers, inline or from a definition file.',
        'It is parsed into every form the templates need. A name is wanted as singular, plural, PascalCase, camelCase, snake_case and kebab-case, and getting pluralisation wrong is how a route ends up at /api/v1/persons.',
        'Each field type is mapped to its five outputs through one table: the Go type, the GORM tag, the Zod validator, the TypeScript type and the admin form control. One table rather than five, because five would disagree within a month.',
        'The files that are written whole are written: model, service, handler, schema, types, hooks, admin definition, admin page, importer.',
        'The files that already exist are located by anchor. The anchors are comments the scaffold emits, which makes them a contract: renaming one means the scaffold stops emitting a hook the generator needs, and the generator then silently injects nothing.',
        'Injection happens before the anchor, or inline. An import block is assembled as a list rather than by replacing an anchor line, because a replace whose anchor has moved matches nothing, returns the input unchanged and reports success.',
        'Later, the resource is removed. Every injection has a matching removal, written in the same change that added the injection, because an injection with no removal leaves a dangling import or a link to a deleted page and breaks the build on an operation meant to clean up.',
        'The project compiles either way. That is the test: generate, compile, remove, compile. A round trip that leaves no trace is the only evidence the pair is complete.',
      ],
    },
    dataFlow: [
      'The definition is the single input. Everything about a resource that any layer needs is derived from it, so the layers cannot disagree.',
      'Anchors are load-bearing and shared between two packages. A new anchor goes into the scaffold template in the same commit as the generator code that looks for it.',
      'Flags change the shape rather than adding a wrapper. Owned by user means the list is scoped and the handlers check ownership, not that a middleware was added.',
      'Generated code is ordinary code. It is formatted, commented and expected to be edited, because a resource always needs something the generator did not predict.',
    ],
  },

  stack: [
    ['Input', 'an inline field list or a definition file'],
    ['Field types', '30, each mapped to five outputs'],
    ['Injection', 'before a named anchor, or inline'],
    ['Removal', 'one remover per injection, added together'],
    ['Sync', 'Go structs to TypeScript and Zod'],
    ['Formatting', 'gofmt on Go output, Prettier conventions on TypeScript'],
    ['Tests', '124 across generate and scaffold, including generate-remove round trips'],
  ],

  api: {
    groups: [
      {
        title: 'The commands',
        rows: [
          { method: 'POST', path: 'grit generate resource Post title:string body:text', what: 'The whole vertical slice' },
          { method: 'POST', path: 'grit generate resource Order --owned-by user', what: 'Scoped to its owner at every layer' },
          { method: 'DELETE', path: 'grit remove Post', what: 'Every file and every injection' },
          { method: 'PATCH', path: 'grit add field Post published:bool', what: 'One field across all layers' },
          { method: 'GET', path: 'grit sync', what: 'Regenerate TypeScript and Zod from Go' },
        ],
      },
    ],
    samples: [
      {
        title: 'What one command produces',
        language: 'bash',
        code: `$ grit generate resource Invoice number:string total:money due_on:date \\
    customer:belongs_to status:enum:draft,sent,paid --owned-by user

  apps/api/internal/models/invoice.go
  apps/api/internal/services/invoice_service.go
  apps/api/internal/handlers/invoice_handler.go
  apps/api/internal/handlers/invoice_import.go
  packages/shared/src/schemas/invoice.ts
  packages/shared/src/types/invoice.ts
  apps/web/src/hooks/use-invoices.ts
  apps/admin/src/resources/invoice.ts
  apps/admin/src/app/invoices/page.tsx

  injected: routes.go, models.go, admin sidebar, resource registry,
            shared index, admin navigation`,
      },
    ],
  },

  lowLevel: {
    classes: [
      {
        name: 'Generator.Run',
        file: 'internal/generate/generator.go',
        what: 'The end-to-end shape: parse, write, inject. The integration tests drive this rather than the pieces, because the pieces passing individually is not evidence the whole works.',
      },
      {
        name: 'injectBefore / injectInline',
        file: 'internal/generate/inject.go',
        what: 'The two injection primitives. Both check for their own output first, because generating the same resource twice must not produce two routes.',
      },
      {
        name: 'go_import.go',
        what: 'Builds an emitted import block as a list. The shape to copy, specifically because a strings.Replace on an anchor line fails silently when the anchor has moved.',
      },
      {
        name: 'remove.go',
        file: 'internal/generate/remove.go',
        what: 'removeLinesContaining, removeInlineText, removeLineBlock, removeSchemaExportBlock. One per injection kind, tested by a generate-then-remove round trip.',
      },
      {
        name: 'sync.go',
        what: 'Parses the Go structs and emits TypeScript types and Zod schemas. The reason the response shape, the schema and the type can change in one commit.',
      },
    ],
    principles: [
      {
        label: 'Every injection needs a removal, in the same change',
        text: 'they are a pair. Adding one half leaves a project that cannot clean up after itself, and the breakage appears during an operation that was supposed to help.',
      },
      {
        label: 'Never replace on an anchor line',
        text: 'a no-op replace reports success. Build the block as a list so a missing anchor is an error rather than a silent nothing.',
      },
      {
        label: 'Idempotent or a bug',
        text: 'check for your own output before injecting. Somebody will run the command twice.',
      },
      {
        label: 'One writer for shared text',
        text: 'anything emitted at both generate time and upgrade time comes from one function. Two copies of the same template diverge, and the divergence is between a new project and an upgraded one, which is the hardest pair to compare.',
      },
      {
        label: 'Generate all thirty types and compile',
        text: 'a type mapping that is wrong for money or json is invisible in a resource that uses neither, and the resource that uses them is somebody else’s project.',
      },
      {
        label: 'Never call a service from a model',
        text: 'models to services is an import cycle. Auto-numbering in a create hook calls the sequence directly, and the generator has to know that.',
      },
    ],
    patterns: [
      ['Template method', 'one definition driving every layer’s output'],
      ['Anchor injection', 'named comments as a contract between two packages'],
      ['Inverse operation', 'a remover paired with every injector'],
      ['Single mapping table', 'one field type to five outputs'],
      ['Round-trip test', 'generate, compile, remove, compile'],
    ],
  },

  scaling: [
    'Generation is a development-time operation. The only latency that matters is that it feels instant, and it does.',
    'The injection surface is what grows: six anchors times sixty resources is three hundred and sixty edits a mature project is carrying.',
    'Adding a layer means adding an injector and a remover and a test for the round trip. The cost of a new layer is paid once and then applies to every resource.',
    'Adding a field type means one row in the mapping table and a compile of the all-types resource. That is the cheap axis, deliberately.',
    'Upgrade is where this gets hard: existing projects need the new output applied to code they have edited, which is why shared hunk functions and a manifest of what has been modified exist.',
    'The real scaling limit is not the generator, it is that generated code is read by people. Output nobody wants to read is output they fork, and then the generator is irrelevant.',
  ],

  bottlenecks: {
    problems: [
      {
        label: 'Anchors that move',
        text: 'rename one and the generator silently injects nothing. Nothing errors, and the resource is half-generated.',
      },
      {
        label: 'Edited generated files',
        text: 'a project that changed a generated file cannot have the new version written over it, and the upgrade either loses the edit or skips the fix.',
      },
      {
        label: 'Injections without removals',
        text: 'the project stops compiling during a removal, which is the operation least likely to be tested and most likely to be run in a hurry.',
      },
      {
        label: 'Type mapping gaps',
        text: 'thirty types times five outputs is a hundred and fifty mappings, and the wrong one is invisible until a project uses that type.',
      },
      {
        label: 'Divergence between generate and upgrade',
        text: 'the same code emitted by two writers drifts, and the difference is between a fresh project and an upgraded one.',
      },
    ],
    improvements: [
      {
        label: 'Test the anchors exist',
        text: 'a test asserting every anchor the generator looks for is present in the scaffold output. A rename then fails a test rather than producing half a resource.',
      },
      {
        label: 'A manifest of what was modified',
        text: 'a hash per generated file, so an upgrade can report an edited file rather than silently replacing or skipping it.',
      },
      {
        label: 'Round-trip every injection in a test',
        text: 'generate, compile, remove, compile. It is the only check that a new injector has its remover.',
      },
      {
        label: 'One resource with every field type, compiled in CI',
        text: 'thirty types in one resource that has to build. A bad mapping fails the build rather than somebody’s project.',
      },
      {
        label: 'One function, called from both paths',
        text: 'generate and upgrade call the same writer. It is the only arrangement where the two cannot differ.',
      },
    ],
  },

  seeAlso: [
    { title: 'Generator reference', href: '/docs/cli/generate' },
    { title: 'API contract', href: '/docs/systems/api-contract' },
  ],
}

export const API_CONTRACT: SystemDesign = {
  slug: 'api-contract',
  name: 'API Contract',
  tagline:
    'One response envelope, one table of error codes, and the reason a client written against one endpoint works against all of them.',
  group: 'Platform',
  packages: ['internal/respond', 'internal/errorcodes', 'internal/access'],

  problem: {
    text: [
      'Six clients are built against this API: an admin panel, a web app, an Expo app, a desktop client, a generated SDK and whatever an operator writes next. Each one has code that reads a response, finds the data, finds the error, and decides what to show.',
      'If two endpoints answer in different shapes, that code has to handle both. In practice it does not: it handles whichever shape the developer was looking at, and the other one produces a blank screen. A handler answering in a different shape does not fail a test. It fails one screen, usually the one nobody opened before release.',
      'The subtle version of the problem is not the shape, it is the status codes, and specifically having two tables of them. The named helpers each carried a status, and a generated code table carried the same pairs, which is how one error code came back as 422 from a helper and 400 from a handler that wrote its own envelope. A client branching on the status saw two different outcomes for the same condition.',
      'There are also two distinctions that get muddled everywhere and are worth deciding once. 400 against 422: cannot parse against understood and refused. And 403 against 404 for a record somebody may not see, where 403 confirms the record exists and is therefore an information leak on any endpoint keyed by a guessable identifier.',
    ],
    capabilities: [
      'Answer every request in one envelope: data and message, or data and meta, or error.',
      'Carry a stable machine-readable code on every error, from one catalogue.',
      'Derive the status from the code, in one place, so the same condition always has the same status.',
      'Put per-field messages where a form can place them under their input.',
      'Never leak internal detail in a 500, while giving the operator the cause and a request identifier.',
      'Answer 404 rather than 403 for a record the caller may not see.',
      'Describe which routes exist and what each demands of a caller, as data.',
      'Generate a client from the contract rather than by hand.',
    ],
  },

  functional: [
    'An object at the top level of every response, never a bare array or string.',
    'Single reads: data, with an optional human message.',
    'Lists: data and meta, with total, page, page_size and pages always present.',
    'Cursor lists: next_cursor, has_more and mode.',
    'Errors: error with code, message and optional per-field details.',
    'A catalogue of error codes, each with its status, generating the code table.',
    'Response helpers, one per status, reading the status from the catalogue.',
    'A generated access table: every route, and what it demands of a caller.',
    'A version prefix on every route, from one constant.',
  ],

  nonFunctional: [
    {
      label: 'One table of code to status',
      text: 'there were two, which is how one code answered 422 from a helper and 400 from a hand-written envelope. One catalogue, and no handler writes an error body by hand.',
    },
    {
      label: 'Always an object',
      text: 'an array at the top level cannot gain a meta later without breaking every caller. The envelope is what makes the API extensible.',
    },
    {
      label: 'Zero is an answer',
      text: 'the pagination counts are always present. A response that omits total leaves every client doing arithmetic on undefined, which renders as a blank card rather than a nought.',
    },
    {
      label: 'Codes are stable, messages are not',
      text: 'clients branch on the code. The message is one sentence for a person and may change without notice.',
    },
    {
      label: 'A 500 never explains itself',
      text: 'a driver error can describe the schema and sometimes contains SQL. The caller gets a code and a sentence; the operator gets the cause and the request identifier that ties a report to a log line.',
    },
    {
      label: '404 over 403 for records',
      text: 'a 403 confirms existence. 403 is reserved for an action the caller may not perform on a record they can already see.',
    },
  ],

  capacity: {
    assumptions: [
      ['Clients built against the contract', '6'],
      ['Endpoints in a mature project', '200 to 500'],
      ['Error codes in the catalogue', '~40'],
      ['Envelope overhead', '~40 bytes per response'],
    ],
    estimates: [
      {
        label: 'What consistency is worth',
        working: [
          '6 clients x 1 response reader each, not 6 x N shapes',
          '1 error handler per client, branching on ~40 known codes',
          'instead of per-endpoint special cases',
        ],
        note: 'The saving is not bytes, it is that a client can be written once. Six clients times two hundred endpoints is where inconsistency actually costs.',
      },
      {
        label: 'Envelope overhead',
        working: [
          '~40 bytes of wrapper per response',
          'at 926 req/s = ~37 KB/s',
          'under gzip, effectively nothing',
        ],
      },
      {
        label: 'The cost of two status tables',
        working: [
          '1 code answering 2 different statuses',
          'x every client branching on status',
          '= an intermittent bug per client',
        ],
        note: 'This is the kind of defect that gets reported as "sometimes the form does not show the error", investigated three times and closed as unreproducible.',
      },
    ],
  },

  highLevel: {
    intro:
      'A catalogue generates a code table. Helpers read the table. Nothing writes an envelope by hand.',
    components: [
      {
        label: 'The envelope',
        text: 'data and message, data and meta, or error. Three shapes, and an object at the top level in all of them.',
      },
      {
        label: 'Error code catalogue',
        text: 'the one table of code to status. It generates the code table the runtime uses, which is what stops a second table existing.',
      },
      {
        label: 'Response helpers',
        text: 'one per status, each a line over a common failure writer, each reading its status from the catalogue.',
      },
      {
        label: 'Validation details',
        text: 'a field name to a sentence, so a form puts each message under its own input instead of showing one summary.',
      },
      {
        label: 'Access table',
        text: 'generated data describing every route and what it demands: authentication, a role, a permission, or nothing.',
      },
      {
        label: 'Version prefix',
        text: 'one constant. A prefix written out per route group is a prefix that will be inconsistent.',
      },
      {
        label: 'Generated client',
        text: 'the TypeScript types, the Zod schemas and the hooks, produced from the same definitions as the handlers.',
      },
    ],
    flow: {
      title: 'An error, from a handler to a form field',
      nodes: [
        { id: 'req', label: 'Request', col: 0, row: 0 },
        { id: 'handler', label: 'Handler', col: 1, row: 0 },
        { id: 'helper', label: 'respond.Validation', col: 2, row: 0, accent: true },
        { id: 'catalog', label: 'Code Catalogue', col: 3, row: 0, accent: true },
        { id: 'status', label: 'Status 422', col: 3, row: 1 },
        { id: 'envelope', label: 'Error Envelope', col: 2, row: 1, accent: true },
        { id: 'client', label: 'Client Reader', col: 1, row: 1, accent: true },
        { id: 'field', label: 'Message Under Input', col: 1, row: 2 },
        { id: 'log', label: 'Request ID + Log', col: 3, row: 2 },
      ],
      edges: [
        { from: 'req', to: 'handler', step: 1 },
        { from: 'handler', to: 'helper', step: 2 },
        { from: 'helper', to: 'catalog', step: 3 },
        { from: 'catalog', to: 'status', step: 4 },
        { from: 'status', to: 'envelope', step: 5 },
        { from: 'envelope', to: 'client', step: 6 },
        { from: 'client', to: 'field', step: 7 },
        { from: 'helper', to: 'log', step: 8, bend: 'v' },
      ],
      steps: [
        'A request arrives with a value that is well formed and wrong: an email already taken, a total that does not balance.',
        'The handler calls a response helper. It does not write an envelope, because a hand-written envelope is how a code came to carry two different statuses.',
        'The helper reads the status for its code out of the catalogue. The catalogue is the only table of code to status, which is the single most important property of the whole arrangement.',
        'The status is 422, because the request was understood and the answer is no. A 400 would mean it could not be parsed, and the distinction is the one that gets muddled most often.',
        'The envelope is written: a code the client branches on, a sentence for a person, and a map of field names to their own sentences.',
        'Every client has one function that reads this. It finds the code in the same place for every endpoint, which is what makes writing six clients tractable.',
        'The per-field details let the form put each message under its own input. A single summary message is the alternative, and it is the version users cannot act on.',
        'In the 500 case this diverges: the caller gets a generic sentence while the operator gets the real error logged with the method, the path and the request identifier. A driver error can describe the schema and sometimes contains SQL, so it never goes out.',
      ],
    },
    dataFlow: [
      'The code is for machines, the message is for people. A client that branches on message text breaks the first time somebody improves the wording.',
      'A record the caller may not see answers 404. A 403 confirms it exists, which is an information leak on any endpoint keyed by a guessable identifier.',
      'The pagination counts are always present, including zero, because a missing count is worse for every client than a zero one.',
      'The access table is generated, so the answer to "what does this route require" is data rather than a reading of middleware chains.',
    ],
  },

  stack: [
    ['Envelope', 'data and message, data and meta, or error'],
    ['Error codes', 'screaming snake case, from one catalogue'],
    ['Status mapping', 'generated from the catalogue, one table'],
    ['Helpers', 'one per status, no hand-written error bodies'],
    ['Versioning', 'a path prefix from one constant'],
    ['Access', 'a generated table of route requirements'],
    ['Client', 'generated types, schemas and hooks'],
  ],

  api: {
    groups: [
      {
        title: 'The status codes, and when each one is right',
        rows: [
          { method: 'GET', path: '200', what: 'Read, update, delete, and anything else that worked' },
          { method: 'POST', path: '201', what: 'A record was created' },
          { method: 'POST', path: '400', what: 'Cannot be parsed, or a required parameter is absent' },
          { method: 'GET', path: '401', what: 'No credentials, or they are not valid' },
          { method: 'POST', path: '403', what: 'An action this caller may not perform on a record they can see' },
          { method: 'GET', path: '404', what: 'No such record, or none this caller may see' },
          { method: 'PUT', path: '409', what: 'A duplicate, or a version mismatch' },
          { method: 'POST', path: '422', what: 'Understood, and the answer is no' },
          { method: 'GET', path: '429', what: 'Rate limited' },
          { method: 'GET', path: '500', what: 'We broke, and the response does not say how' },
        ],
      },
    ],
    samples: [
      {
        title: 'The helpers, which are the only way an envelope is written',
        language: 'go',
        code: `respond.OK(c, user)                        // 200 { data }
respond.OK(c, user, "Saved")               // 200 { data, message }
respond.Created(c, user, "User created")   // 201

respond.BadRequest(c, "page must be a number")     // 400
respond.Unauthorized(c, "")                        // 401
respond.Forbidden(c, "You cannot publish")         // 403
respond.NotFound(c, "User not found")              // 404
respond.Conflict(c, "That email is already taken") // 409
respond.Validation(c, "Check the form", fields)    // 422
respond.Internal(c, err)                           // 500, logged, not explained`,
      },
      {
        title: 'An error a form can act on',
        language: 'json',
        code: `{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Check the form",
    "details": {
      "email": "That email is already taken",
      "total": "Debits and credits do not balance"
    }
  }
}`,
      },
    ],
  },

  lowLevel: {
    classes: [
      {
        name: 'errorcodes.Catalog',
        file: 'internal/errorcodes/catalog.go',
        what: 'The one table of code to status. It generates the runtime code table, which is the mechanism that prevents a second table existing.',
      },
      {
        name: 'respond helpers',
        what: 'One per status, each a line over a shared failure writer. Short enough that writing an envelope by hand is never the easier option.',
      },
      {
        name: 'respond.Internal',
        what: 'Logs the error with the method, path and request identifier, and sends a generic message. The asymmetry is deliberate: the operator gets the cause, the caller gets a code.',
      },
      {
        name: 'access.Table',
        file: 'internal/access/access.go',
        what: 'Generated. Every route and what it demands of a caller, as data, so the question can be answered without reading middleware chains.',
      },
    ],
    principles: [
      {
        label: 'One table, not two',
        text: 'the duplicate status table is the defect this whole design is organised around. One code, one status, one place.',
      },
      {
        label: 'Always an object',
        text: 'a top-level array cannot gain a sibling field later. The envelope is what makes the API able to change.',
      },
      {
        label: 'Codes for machines, messages for people',
        text: 'the code is a contract and the message is copy. Keeping them separate lets the copy improve without breaking a client.',
      },
      {
        label: '404 by default for records',
        text: 'existence is information. The leak is free to prevent and expensive to notice.',
      },
      {
        label: 'Never explain a 500',
        text: 'the detail that helps an operator is the detail that helps an attacker. The request identifier is how both get what they need.',
      },
    ],
    patterns: [
      ['Envelope', 'one shape for every response'],
      ['Single source of truth', 'a catalogue generating the status table'],
      ['Facade', 'helpers in front of the envelope'],
      ['Code-first client generation', 'types and schemas from the handlers’ definitions'],
    ],
  },

  scaling: [
    'A consistent contract is what makes a second and third client cheap. Six clients against one shape is one reader written six times; against many shapes it is unbounded.',
    'Envelope overhead is negligible and compresses away. Nothing about consistency costs throughput.',
    'The version prefix is what allows a breaking change: a second version alongside the first, rather than a flag day.',
    'A generated client means a contract change is a compile error in every client that uses it, which is the only mechanism that scales past two consumers.',
    'The access table makes authorisation auditable as data, which is what lets a security review be a query rather than a reading exercise.',
    'The limit is discipline: one hand-written envelope reintroduces the inconsistency, which is why there are helpers for every case and a rule against writing error bodies directly.',
  ],

  bottlenecks: {
    problems: [
      {
        label: 'Hand-written envelopes',
        text: 'one handler writing its own error body is how the second status table came back, and it is invisible until a client branches on the status.',
      },
      {
        label: 'The 400 against 422 muddle',
        text: 'both are client errors and the distinction is semantic, so they get used interchangeably and clients cannot rely on either.',
      },
      {
        label: 'Breaking changes',
        text: 'a field removed or a type changed breaks every client, and the version prefix is the only escape.',
      },
      {
        label: 'Message text as a contract',
        text: 'a client branching on message text breaks the first time the copy is improved, and the break is in the client rather than where the change was made.',
      },
      {
        label: 'Pagination counts omitted',
        text: 'a missing total becomes undefined in arithmetic, which renders as a blank card rather than a zero and looks like a data problem.',
      },
    ],
    improvements: [
      {
        label: 'Lint for hand-written error bodies',
        text: 'a check that no handler writes a JSON error directly. It is a grep, and it protects the property the whole contract rests on.',
      },
      {
        label: 'Write the distinction down and test it',
        text: 'contract tests asserting a malformed body is 400 and a failed rule is 422, per endpoint. The semantics then stop being a judgement call.',
      },
      {
        label: 'Version additively',
        text: 'add fields, never remove or retype. A new version only for changes that cannot be additive.',
      },
      {
        label: 'Generate the client from the contract',
        text: 'then a contract change is a compile error rather than a runtime surprise in a client nobody is looking at.',
      },
      {
        label: 'Assert the counts in a test',
        text: 'a list test that checks total, page, page_size and pages are present on an empty result. One test, and it covers the whole class of blank-card bugs.',
      },
    ],
  },

  seeAlso: [
    { title: 'API design contract', href: '/docs/api/design' },
    { title: 'Pagination and search', href: '/docs/systems/pagination-and-search' },
    { title: 'Code generation', href: '/docs/systems/code-generation' },
  ],
}

export const FEATURE_FLAGS: SystemDesign = {
  slug: 'feature-flags',
  name: 'Feature Flags',
  tagline:
    'Turning a feature on for ten per cent of users without a deploy, and having the same user stay in the same group.',
  group: 'Platform',
  packages: ['internal/flags', 'internal/models'],

  problem: {
    text: [
      'A redesigned checkout is ready. Shipping it to everybody at once means that if it is wrong, the fix is a deploy, and the damage is everybody. Shipping it to ten per cent means the damage is ten per cent and the fix is a toggle.',
      'That requires the decision to be made at runtime rather than at build time, which rules out a constant and a compile-time flag. It also requires the decision to be stable per user: somebody who sees the new checkout on Monday must see it on Tuesday, or they are not a user of either version, they are a user of a flickering application.',
      'The naive implementation reads a row per check. A page rendering a dozen flagged components does a dozen queries, and the flag system becomes the slowest part of the request it was meant to make safe. The naive exposure tracking is worse: a goroutine and an insert per check turned a busy page into thousands of writes a second.',
      'And a flag is only half useful without knowing who saw what. A ten per cent rollout with no record of which users were in it cannot be measured, so the experiment produces a feeling rather than a result.',
    ],
    capabilities: [
      'Turn a feature on or off at runtime, without a deploy.',
      'Roll out to a percentage of users, stably per user.',
      'Serve several variants, for comparing two new versions rather than one.',
      'Check a flag without touching the database.',
      'Propagate a change in seconds, not on the next deploy.',
      'Record exposures, so a rollout can be measured.',
      'Target specific users, roles or tenants, not only percentages.',
      'Tell clients a flag changed, so they refetch rather than waiting.',
    ],
  },

  functional: [
    'A flags table with a key, a state, rules and variants.',
    'An engine holding the flags in memory, one per process.',
    'A background refresh on an interval, 30 seconds by default.',
    'An immediate refresh and a realtime broadcast on an administrative change.',
    'Percentage rollout bucketed per user and flag.',
    'Variants, for more than two arms.',
    'Exposure records, written by one writer rather than per check.',
    'An admin page to create, edit, roll out and delete.',
    'A client-side hook reading the evaluated flags for the current user.',
  ],

  nonFunctional: [
    {
      label: 'A flag check never touches the database',
      text: 'all flags are in memory. A page with a dozen checks must not be a page with a dozen queries, or the mechanism costs more than the feature it guards.',
    },
    {
      label: 'Sticky per user and flag',
      text: 'the bucket is a hash of the user identifier and the flag name, so a user always lands in the same bucket for a given flag and their assignment does not flicker between sessions.',
    },
    {
      label: 'Independent per flag',
      text: 'including the flag name in the hash means the ten per cent for one flag is not the same ten per cent as another. Otherwise one unlucky cohort gets every experiment.',
    },
    {
      label: 'Propagates in seconds',
      text: '30 seconds from the refresh, immediately on an administrative change. A flag that needs a deploy to take effect is a constant.',
    },
    {
      label: 'Exposure tracking never blocks',
      text: 'fire and forget to a single writer. A goroutine and an insert per check was thousands of writes a second on a busy page.',
    },
    {
      label: 'Fails to a known state',
      text: 'an unknown flag is off. The default has to be the safe one, because the failure case is a flag that was never created or was just deleted.',
    },
  ],

  capacity: {
    assumptions: [
      ['Flags in a mature project', '20 to 100'],
      ['Flag checks per request', '1 to 12'],
      ['Requests per second', '926 average, 4,630 peak'],
      ['Refresh interval', '30 seconds'],
      ['Flag row size', '~400 bytes'],
    ],
    estimates: [
      {
        label: 'Check cost',
        working: [
          'a read lock on a map, plus one SHA-256 over ~40 bytes',
          '= sub-microsecond',
          '926 req/s x 12 checks = ~11,000 checks/s, unmeasurable',
        ],
        note: 'Against a database read per check at ~1 ms, 11,000 checks a second would be eleven seconds of database time per second. That is the difference the in-memory cache makes.',
      },
      {
        label: 'Refresh load',
        working: [
          '1 query every 30 s per process',
          'x 12 processes = 0.4 queries/s',
          'returning ~100 rows of ~400 bytes',
        ],
      },
      {
        label: 'Memory',
        working: [
          '100 flags x ~400 bytes = 40 KB per process',
        ],
        note: 'Which is why holding all of them is obviously right rather than a trade-off worth analysing.',
      },
      {
        label: 'Exposure writes',
        working: [
          'naive: one goroutine and one INSERT per check = 11,000/s',
          'batched through one writer: ~10 inserts/s of 1,000 rows',
        ],
        note: 'The naive version is the thing that broke. A single writer with batching is three orders of magnitude fewer statements for the same information.',
      },
    ],
  },

  highLevel: {
    intro:
      'One engine per process holding every flag in memory, refreshed on a timer and on change, with bucketing by hash.',
    components: [
      {
        label: 'Flags table',
        text: 'the definition: a key, whether it is on, its rules, its rollout percentage and its variants.',
      },
      {
        label: 'Engine',
        text: 'the in-memory map, one per process, pre-warmed at boot behind a read-write lock. A flag check reads the map and never the database.',
      },
      {
        label: 'Refresher',
        text: 'a background loop pulling fresh state every 30 seconds, which is quick enough for an administrative change to feel immediate without polling hard.',
      },
      {
        label: 'Bucketer',
        text: 'SHA-256 of the user identifier and the flag name, modulo 100. Sticky per pair, and independent across flags because the name is in the hash.',
      },
      {
        label: 'Variant selector',
        text: 'the same bucket mapped onto several arms, for comparing two new versions rather than one against the old.',
      },
      {
        label: 'Exposure writer',
        text: 'one writer fed by a channel. Each check used to start a goroutine and an insert of its own, which a busy page turned into thousands a second.',
      },
      {
        label: 'Change broadcast',
        text: 'an administrative change refreshes the engine immediately and publishes a realtime event, so subscribed clients refetch instead of waiting out the interval.',
      },
    ],
    flow: {
      title: 'A flag check, and a change propagating',
      nodes: [
        { id: 'req', label: 'Request', col: 0, row: 0 },
        { id: 'check', label: 'flags.IsEnabled', col: 1, row: 0, accent: true },
        { id: 'map', label: 'In-Memory Map', col: 2, row: 0, accent: true },
        { id: 'bucket', label: 'Hash Bucket', col: 3, row: 0, accent: true },
        { id: 'variant', label: 'Variant', col: 3, row: 1 },
        { id: 'exposure', label: 'Exposure Writer', col: 2, row: 1, accent: true },
        { id: 'admin', label: 'Admin Change', col: 0, row: 2 },
        { id: 'refresh', label: 'Refresh + Broadcast', col: 1, row: 2, accent: true },
        { id: 'clients', label: 'Clients Refetch', col: 2, row: 2 },
      ],
      edges: [
        { from: 'req', to: 'check', step: 1 },
        { from: 'check', to: 'map', step: 2 },
        { from: 'map', to: 'bucket', step: 3 },
        { from: 'bucket', to: 'variant', step: 4 },
        { from: 'check', to: 'exposure', step: 5, bend: 'v' },
        { from: 'admin', to: 'refresh', step: 6 },
        { from: 'refresh', to: 'map', step: 7, bend: 'v' },
        { from: 'refresh', to: 'clients', step: 8 },
      ],
      steps: [
        'A request renders a page with a dozen flagged components, each asking whether its feature is on for this user.',
        'Each check reads the in-memory map under a read lock. No database round trip, which is the property that lets a page have twelve checks instead of one.',
        'For a percentage rollout, the bucket is SHA-256 of the user identifier and the flag name, modulo 100. The user identifier makes it sticky, so Monday’s assignment is Tuesday’s. The flag name makes flags independent, so the unlucky ten per cent of one experiment is not the unlucky ten per cent of all of them.',
        'With variants, the same bucket maps onto several arms rather than on and off, which is how two new versions get compared against each other instead of only against the old one.',
        'The exposure is recorded, fire and forget, through one writer fed by a channel. Every check starting its own goroutine and its own insert is what made a busy page thousands of writes a second.',
        'Separately, somebody changes a flag in the admin: enables it, moves it to twenty-five per cent, adds a variant.',
        'The engine refreshes immediately rather than waiting out the interval, and the timer continues as the fallback that keeps every other process in step within thirty seconds.',
        'A realtime event goes out, so clients holding evaluated flags refetch rather than showing a stale answer until their next navigation. An anonymous user buckets on a random per-request value, which is effectively random, so anything that needs stickiness without a login passes a session or device identifier instead.',
      ],
    },
    dataFlow: [
      'The user identifier and the flag name are both in the hash. Omitting the first makes assignment flicker; omitting the second makes one cohort receive every experiment.',
      'An unknown flag is off. That covers the window after a deploy that references a flag nobody created yet, and the window after a deletion.',
      'Exposures are a stream to one writer, not a write per check. The information is the same and the statement count is three orders of magnitude lower.',
      'An anonymous check is random per request unless the caller supplies a stable identifier, and the package says so rather than pretending otherwise.',
    ],
  },

  stack: [
    ['Storage', 'a flags table and an exposures table'],
    ['Evaluation', 'in memory, one engine per process'],
    ['Refresh', '30 seconds, plus immediate on change'],
    ['Bucketing', 'SHA-256 of user id and flag name, modulo 100'],
    ['Variants', 'several arms over the same bucket'],
    ['Exposure', 'fire and forget to a single batched writer'],
    ['Propagation', 'a realtime flag.updated event'],
  ],

  dataModel: {
    entities: [
      {
        name: 'feature_flags',
        fields: [
          ['key', 'the name used in code'],
          ['enabled', 'the master switch'],
          ['rules', 'JSON: rollout percentage, targeted users, roles, tenants'],
          ['variants', 'JSON: the arms and their weights'],
          ['description', 'what it is for, because a flag outlives whoever added it'],
        ],
      },
      {
        name: 'flag_exposures',
        fields: [
          ['flag_key', 'which flag'],
          ['user_id', 'who saw it'],
          ['variant', 'which arm they got'],
          ['created_at', 'when'],
        ],
        note: 'This is what turns a rollout into a measurement. Without it, a ten per cent experiment produces a feeling rather than a result.',
      },
    ],
  },

  api: {
    groups: [
      {
        title: 'Flags',
        rows: [
          { method: 'GET', path: '/api/v1/flags', what: 'Evaluated flags for the current caller' },
          { method: 'GET', path: '/api/v1/admin/flags', what: 'All flags with their rules' },
          { method: 'POST', path: '/api/v1/admin/flags', what: 'Create one' },
          { method: 'PUT', path: '/api/v1/admin/flags/:key', what: 'Change state, rollout or variants' },
          { method: 'DELETE', path: '/api/v1/admin/flags/:key', what: 'Remove it' },
        ],
      },
    ],
    samples: [
      {
        title: 'Two arms, or several',
        language: 'go',
        code: `if flags.IsEnabled(c, "new_checkout") {
    return h.newCheckout(c)
}

switch flags.Variant(c, "checkout_redesign") {
case "control":   return h.oldFlow(c)
case "variant_a": return h.newFlow(c)
case "variant_b": return h.alternateFlow(c)
}`,
      },
    ],
  },

  lowLevel: {
    classes: [
      {
        name: 'flags.Engine',
        file: 'internal/flags/flags.go',
        what: 'Owns the in-memory cache, one per process, pre-warmed at construction. Takes the realtime hub optionally, so broadcasts are available without being required.',
        methods: ['IsEnabled', 'Variant', 'Refresh'],
      },
      {
        name: 'DefaultRefreshInterval',
        what: '30 seconds. Quick enough that an administrative change feels immediate, slow enough not to poll the database for no reason.',
      },
      {
        name: 'Bucketing',
        what: 'SHA-256 of the user identifier, a separator and the flag name, modulo 100. Both inputs matter: the first for stickiness, the second for independence between flags.',
      },
      {
        name: 'Exposure channel',
        what: 'Feeds one writer. Each check used to start a goroutine and an insert of its own, which a busy page turned into thousands a second.',
      },
    ],
    principles: [
      {
        label: 'Evaluate in memory',
        text: 'a flag check has to be cheaper than the branch it guards, or nobody will put one on a hot path and the mechanism goes unused where it matters.',
      },
      {
        label: 'Hash both inputs',
        text: 'the user for stickiness, the flag for independence. Each omission produces a different and confusing failure.',
      },
      {
        label: 'Unknown means off',
        text: 'the default has to be the safe one, because the uncertain cases are a flag not yet created and a flag just deleted.',
      },
      {
        label: 'One writer for a high-frequency side effect',
        text: 'the exposure write happens as often as the check. A goroutine per check is not a smaller version of a writer, it is a different system with no backpressure.',
      },
      {
        label: 'A flag needs a description and an owner',
        text: 'flags outlive the people who add them, and a flag nobody can explain is never removed.',
      },
    ],
    patterns: [
      ['Consistent hashing', 'stable bucket assignment per user and flag'],
      ['Read-through cache with refresh', 'in-memory state on an interval'],
      ['Fire and forget', 'exposures that never block a check'],
      ['Single writer', 'high-frequency inserts batched through one path'],
      ['Pub/sub invalidation', 'a change broadcast so clients refetch'],
    ],
  },

  scaling: [
    'Checks are free: a map read and a hash, so flag count and check frequency both scale without cost.',
    'Refresh load is one query per process per interval, which is independent of traffic entirely.',
    'Memory is tens of kilobytes, which is why holding every flag is obviously right rather than a decision.',
    'Exposure writes are the only part that scales with traffic, and batching through one writer is what keeps them from being the dominant write load in the application.',
    'The thirty second window means processes can briefly disagree. For a rollout that is harmless; for a kill switch it is the reason the immediate refresh and the broadcast exist.',
    'The real scaling problem is organisational: flags accumulate, and a project with two hundred of them has two hundred untested combinations.',
  ],

  bottlenecks: {
    problems: [
      {
        label: 'Flag debt',
        text: 'a flag left in after its rollout finished is a permanent branch, and a hundred of them is a codebase with no single defined behaviour.',
      },
      {
        label: 'Propagation window',
        text: 'thirty seconds of disagreement between processes, which matters when the flag is being used to turn something off in a hurry.',
      },
      {
        label: 'Exposure volume',
        text: 'the exposure table grows with checks rather than with actions, which makes it the fastest growing table in the system.',
      },
      {
        label: 'Anonymous stickiness',
        text: 'with no user identifier the bucket is random per request, so an anonymous user sees a different arm on every page unless a stable identifier is passed.',
      },
      {
        label: 'Untested combinations',
        text: 'twenty independent flags is a million combinations, and the suite tests one of them.',
      },
    ],
    improvements: [
      {
        label: 'Expire flags deliberately',
        text: 'a review date on every flag and a report of the overdue ones. Removing a flag is part of shipping the feature, not a separate project.',
      },
      {
        label: 'Broadcast every change',
        text: 'immediate refresh plus a realtime event, with the interval as the fallback. The window then only applies when the broadcast is missed.',
      },
      {
        label: 'Sample or aggregate exposures',
        text: 'one row per user per flag per day rather than per check. The measurement is the same and the volume is bounded by users rather than by traffic.',
      },
      {
        label: 'Bucket anonymous users on a device identifier',
        text: 'a session or device value gives stickiness without a login, and the package should take it rather than silently being random.',
      },
      {
        label: 'Test both arms of what matters',
        text: 'not every combination, but both sides of any flag guarding a path with its own tests. The combinatorial problem is a reason to have fewer flags, not more tests.',
      },
    ],
  },

  seeAlso: [
    { title: 'Feature flags reference', href: '/docs/features/feature-flags' },
    { title: 'Caching', href: '/docs/systems/caching' },
  ],
}
