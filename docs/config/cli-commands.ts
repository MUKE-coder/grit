/**
 * The CLI command catalogue behind /docs/cli.
 *
 * Every `output` block and every `files` list here was captured from a real run
 * against a freshly scaffolded `--triple --next` project, with the file effects
 * read out of `git status` rather than written from memory. That is the whole
 * point of the page: a reference that says a command touches nineteen files and
 * then names eighteen of them is worse than no reference, because you only find
 * the missing one when it breaks.
 *
 * When a command changes, re-capture rather than edit by hand:
 *
 *   grit new capture --triple --next && cd capture && git init && git add -A
 *   git commit -m base
 *   grit <the command>
 *   git add -A && git status --porcelain
 */

export type FileStatus = 'created' | 'modified' | 'deleted'

export interface FileChange {
  path: string
  status: FileStatus
  /** What this file is for, when the path alone does not say. */
  note?: string
}

export type CommandCategory =
  | 'Scaffold'
  | 'Generate'
  | 'Add'
  | 'Run'
  | 'Data'
  | 'Ship'
  | 'Meta'

export interface CommandFlag {
  flag: string
  desc: string
}

export interface CliCommand {
  /** URL-safe id, also the hash anchor. */
  id: string
  /** The command as you type it, without arguments. */
  name: string
  alias?: string
  category: CommandCategory
  /** One line, shown in the list. */
  summary: string
  /** The exact invocation the terminal simulates. */
  example: string
  /** Terminal output, verbatim from a real run. */
  output: string[]
  /** Why the command exists, in a sentence or three. */
  purpose: string
  /** When you reach for it. */
  useCases: string[]
  files: FileChange[]
  flags?: CommandFlag[]
  /** Gotchas worth knowing before you run it. */
  notes?: string[]
  docs?: { label: string; href: string }[]
  /** Extra words the search should match: aliases, synonyms, concepts. */
  keywords?: string[]
}

/** Categories in the order the sidebar shows them. */
export const CLI_CATEGORIES: CommandCategory[] = [
  'Scaffold',
  'Generate',
  'Add',
  'Run',
  'Data',
  'Ship',
  'Meta',
]

export const CLI_COMMANDS: CliCommand[] = [
  // ── Scaffold ───────────────────────────────────────────────────────
  {
    id: 'new',
    name: 'grit new',
    category: 'Scaffold',
    summary: 'Scaffold a new project: Go API, Next.js or Vite frontends, admin panel.',
    example: 'grit new myapp --triple --next',
    output: [
      '  Creating project: myapp',
      '',
      '  → Scaffolding Go API (Gin + GORM)...',
      '  → Adding batteries (cache, storage, mail, jobs, AI)...',
      '  → Setting up Next.js web app...',
      '  → Creating admin panel...',
      '  → Writing shared schemas and types...',
      '  → Scaffolding frontend tests (Vitest + Playwright)...',
      '',
      '  ✓ Project created successfully!',
      '',
      '  Next steps:',
      '',
      '    cd myapp',
      '    docker compose up -d      # Postgres, Redis, MinIO, Mailhog',
      '    pnpm install              # frontend deps (one-time)',
      '    grit migrate              # create database tables',
      '    grit seed                 # (optional) sample data',
      '    grit start                # run everything',
    ],
    purpose:
      'The one command that turns an empty folder into a running full-stack project. It writes around 426 files: a Gin + GORM API with auth, roles, jobs, mail, storage and audit already wired, the frontends you asked for, a shared package of Zod schemas and TypeScript types, and the Docker compose file for Postgres, Redis, MinIO and Mailhog.',
    useCases: [
      'Starting anything new.',
      'Trying an architecture out: run it twice with --single and --triple and compare.',
      'Scaffolding into a directory you already made, with --here.',
    ],
    files: [
      { path: '426 files across the monorepo', status: 'created', note: 'the whole project' },
      { path: 'apps/api/', status: 'created', note: 'Go API, batteries wired' },
      { path: 'apps/web/', status: 'created', note: 'Next.js or Vite' },
      { path: 'apps/admin/', status: 'created', note: 'admin panel' },
      { path: 'packages/shared/', status: 'created', note: 'Zod schemas + TS types' },
      { path: 'docker-compose.yml', status: 'created', note: 'Postgres, Redis, MinIO, Mailhog' },
      { path: '.env', status: 'created', note: 'ports, DB URL, secrets' },
      { path: 'grit.json', status: 'created', note: 'what this project is made of' },
    ],
    flags: [
      { flag: '--arch string', desc: 'Architecture: single, double, triple, api, mobile' },
      { flag: '--frontend string', desc: 'Frontend framework: next, vite (TanStack)' },
      { flag: '--single / --double / --triple / --api / --mobile', desc: 'Shorthand for --arch' },
      { flag: '--next / --vite', desc: 'Shorthand for --frontend' },
      { flag: '--full', desc: 'Everything: API + web + admin + desktop + Expo + docs site' },
      { flag: '--desktop', desc: 'Include a Wails desktop app sharing the monorepo API' },
      { flag: '--expo', desc: 'Include an Expo mobile app' },
      { flag: '--i18n', desc: 'Add internationalisation (next-intl, translated API messages)' },
      { flag: '--style string', desc: 'Admin style: default, modern, minimal, glass' },
      { flag: '--theme string', desc: 'Full theme: atlas, aurora, pulse' },
      { flag: '--here', desc: 'Scaffold into the current directory' },
      { flag: '--force', desc: 'Allow scaffolding into a non-empty directory' },
    ],
    notes: [
      'Interactive by default. Passing --triple --next (or any other pair) skips the prompts.',
      '`pnpm install` is not run for you. `grit start` fails on a missing module if you skip it.',
      'The database defaults to Postgres on port 5434. Set DATABASE_URL=sqlite:./app.db in .env to skip Docker entirely.',
    ],
    docs: [
      { label: 'Create a project', href: '/docs/getting-started/create-a-project' },
      { label: 'Architecture modes', href: '/docs/concepts/architecture-modes' },
    ],
    keywords: ['init', 'create', 'scaffold', 'bootstrap', 'start a project'],
  },
  {
    id: 'new-desktop',
    name: 'grit new-desktop',
    category: 'Scaffold',
    summary: 'Create a Wails desktop application that shares the monorepo API.',
    example: 'grit new-desktop',
    output: [
      '  Creating Wails desktop app...',
      '',
      '  ✓ apps/desktop/main.go',
      '  ✓ apps/desktop/wails.json',
      '  ✓ apps/desktop/frontend/',
      '',
      '  Next: cd apps/desktop && wails dev',
    ],
    purpose:
      'Adds a Wails desktop shell to a project that already has an API, so the same Go backend and the same typed client serve a native window.',
    useCases: [
      'Shipping an internal tool people run rather than visit.',
      'An offline-first app that needs a real filesystem and a keystore.',
    ],
    files: [
      { path: 'apps/desktop/', status: 'created', note: 'Wails app + frontend' },
    ],
    notes: [
      'Needs the Wails CLI installed separately.',
      '`grit compile` and `grit package` build and distribute it.',
    ],
    docs: [{ label: 'Desktop', href: '/docs/desktop' }],
    keywords: ['wails', 'native', 'electron alternative'],
  },
  {
    id: 'init',
    name: 'grit init',
    category: 'Scaffold',
    summary: 'Write CLAUDE.md / AGENTS.md convention docs into an existing project.',
    example: 'grit init',
    output: [
      '  ✓ CLAUDE.md',
      '  ✓ AGENTS.md',
      '',
      '  Framework conventions written. Coding agents will read these.',
    ],
    purpose:
      'Drops the framework conventions into the project root where a coding agent will find them: the folder structure, the naming rules, the response format, the markers a generator injects into.',
    useCases: [
      'Before pointing Claude Code, Cursor or Copilot at a Grit project.',
      'Refreshing the convention docs after an upgrade.',
    ],
    files: [
      { path: 'CLAUDE.md', status: 'created' },
      { path: 'AGENTS.md', status: 'created' },
    ],
    docs: [{ label: 'AI skill', href: '/docs/ai-skill' }],
    keywords: ['claude', 'agents', 'ai', 'conventions', 'cursor'],
  },

  // ── Generate ───────────────────────────────────────────────────────
  {
    id: 'generate-resource',
    name: 'grit generate resource',
    alias: 'grit g resource',
    category: 'Generate',
    summary: 'A full-stack CRUD resource: model, service, handler, schemas, types, hooks, admin page.',
    example:
      'grit generate resource Product --fields "name:string,slug:slug,price:float,stock:int,images:files,active:bool" --public --faker',
    output: [
      '  Generating resource: Product',
      '',
      '  ✓ apps/api/internal/models/product.go',
      '  ✓ apps/api/internal/services/product.go',
      '  ✓ apps/api/internal/handlers/product.go',
      '  ✓ apps/api/internal/handlers/product_import.go',
      '  ✓ apps/api/internal/handlers/product_public.go (4 field(s) published)',
      '    Held back: stock, active',
      '    Add any of those to the publicProduct struct in that file to publish them.',
      '  ✓ packages/shared/schemas/product.ts',
      '  ✓ packages/shared/types/product.ts',
      '  ✓ apps/web/hooks/use-products.ts',
      '  ✓ apps/admin/resources/products.ts',
      '  ✓ apps/admin/app/(dashboard)/resources/products/page.tsx',
      '  ✓ apps/admin/app/(dashboard)/resources/products/[id]/page.tsx',
      '',
      '  Injecting into existing files...',
      '  ✓ Injected model into AutoMigrate',
      '  ✓ Injected model into GORM Studio',
      '  ✓ Registered model with sync registry',
      '  ✓ Injected handler initialization',
      '  ✓ Injected protected routes',
      '  ✓ Injected admin routes',
      '  ✓ Registered model with the API reference',
      '  ✓ Documented the endpoints at /docs',
      '  ✓ Injected permissions into the authz catalog',
      '  ✓ Injected schema export',
      '  ✓ Injected type export',
      '  ✓ Injected API route constants',
      '  ✓ Injected resource import into registry',
      '  ✓ Injected resource into registry list',
      '  ✓ GET /api/v1/public/products and /api/v1/public/products/:key (API key required)',
      '',
      '  ✅ Resource Product generated successfully!',
    ],
    purpose:
      'The command Grit is really for. One resource definition produces the Go model, a service holding the business logic, a thin handler, the routes, the Zod schemas, the TypeScript types, the React Query hooks and a working admin page with a table and a form. It also injects itself into fourteen existing files, which is the part you would otherwise forget half of.',
    useCases: [
      'Every new entity in the application.',
      'Add --public for a storefront or mobile client with no logged-in user.',
      'Add --tree for categories, folders or an org chart.',
      'Add --items for an invoice with line items in one form.',
      'Add --faker to fill the table so the admin has something in it.',
    ],
    files: [
      { path: 'apps/api/internal/models/product.go', status: 'created', note: 'GORM model' },
      { path: 'apps/api/internal/services/product.go', status: 'created', note: 'business logic' },
      { path: 'apps/api/internal/handlers/product.go', status: 'created', note: 'thin handler' },
      { path: 'apps/api/internal/handlers/product_import.go', status: 'created', note: 'CSV import' },
      { path: 'apps/api/internal/handlers/product_public.go', status: 'created', note: '--public only, allowlist response' },
      { path: 'apps/api/internal/database/products_seeder.go', status: 'created', note: '--faker only' },
      { path: 'packages/shared/schemas/product.ts', status: 'created', note: 'Zod' },
      { path: 'packages/shared/types/product.ts', status: 'created' },
      { path: 'apps/web/hooks/use-products.ts', status: 'created', note: 'React Query' },
      { path: 'apps/admin/resources/products/products.ts', status: 'created', note: 'table + form definition' },
      { path: 'apps/admin/resources/products/products.custom.tsx', status: 'created', note: 'yours; never overwritten' },
      { path: 'apps/admin/app/(dashboard)/resources/products/page.tsx', status: 'created' },
      { path: 'apps/admin/app/(dashboard)/resources/products/[id]/page.tsx', status: 'created' },
      { path: 'apps/api/internal/models/user.go', status: 'modified', note: 'AutoMigrate registry' },
      { path: 'apps/api/internal/routes/routes.go', status: 'modified', note: 'routes mounted' },
      { path: 'apps/api/internal/routes/apidocs.go', status: 'modified', note: 'OpenAPI' },
      { path: 'apps/api/internal/authz/permissions.go', status: 'modified', note: 'permission catalog' },
      { path: 'apps/api/internal/database/seed.go', status: 'modified' },
      { path: 'apps/api/internal/services/form_share_dispatch.go', status: 'modified' },
      { path: 'apps/api/internal/services/resource_stats_dispatch.go', status: 'modified', note: 'dashboard widgets' },
      { path: 'apps/admin/resources/index.ts', status: 'modified', note: 'sidebar registry' },
      { path: 'packages/shared/schemas/index.ts', status: 'modified' },
      { path: 'packages/shared/types/index.ts', status: 'modified' },
      { path: 'packages/shared/constants/index.ts', status: 'modified', note: 'route constants' },
      { path: '.grit/manifest.json', status: 'modified', note: 'provenance for grit upgrade' },
    ],
    flags: [
      { flag: '--fields string', desc: 'Inline field definitions, e.g. "title:string,published:bool"' },
      { flag: '--from string', desc: 'YAML file defining the resource' },
      { flag: '-i, --interactive', desc: 'Define fields at a prompt' },
      { flag: '--public', desc: 'Read-only list + detail under /api/v1/public/, API-key guarded' },
      { flag: '--tree', desc: 'Hierarchical: parent, materialized path, depth, sibling order, move endpoint' },
      { flag: '--items string', desc: 'Has-many child as a line-items table inside this resource\'s form' },
      { flag: '--roles string', desc: 'Restrict routes to roles, e.g. "ADMIN,EDITOR"' },
      { flag: '--seed / --faker', desc: 'One example record, or many rows with gofakeit' },
      { flag: '--count int', desc: 'Rows for the faker seeder (default 10)' },
      { flag: '--force', desc: 'Generate even when the name collides with a built-in model' },
    ],
    notes: [
      'Your .custom.tsx is written once and never again. Cell renderers and page overrides live there so a regenerate cannot take them back.',
      'The --public handler is also written only once, because the allowlist inside it is yours to edit. Delete it and regenerate to pick up newer generator features.',
      'Run `grit migrate` afterwards. The model is the source of truth and GORM adds the columns.',
    ],
    docs: [
      { label: 'Code generation', href: '/docs/concepts/code-generation' },
      { label: 'Field types', href: '/docs/concepts/field-types' },
      { label: 'Generated file map', href: '/docs/concepts/generated-files' },
    ],
    keywords: ['crud', 'model', 'scaffold entity', 'g resource', 'make model'],
  },
  {
    id: 'generate-field',
    name: 'grit generate field',
    alias: 'grit g field',
    category: 'Generate',
    summary: 'Add one column to a resource that already exists, in place.',
    example: 'grit generate field Product weight:float',
    output: [
      '  Field added to Product.',
      "  Run 'grit migrate' to add the database column.",
    ],
    purpose:
      'The small change you actually make most often. It injects the column into the Go model, both Zod schemas, the TypeScript type and the admin form and table, at the auto markers, so nothing you wrote by hand around them moves.',
    useCases: [
      'One more column, without regenerating the resource and re-reviewing every file.',
      'Adding a select or toggle whose options you want in the admin form immediately.',
    ],
    files: [
      { path: 'apps/api/internal/models/product.go', status: 'modified' },
      { path: 'packages/shared/schemas/product.ts', status: 'modified' },
      { path: 'packages/shared/types/product.ts', status: 'modified' },
      { path: 'apps/admin/resources/products/products.ts', status: 'modified', note: 'column + form field' },
    ],
    notes: [
      'Scalar, select and toggle only. For a relationship, file, slug or array field, regenerate the resource.',
      'No migration file is written. The model is the source of truth and `grit migrate` adds the column.',
      'A new string field arrives with binding:"required". Drop it by hand if the column is optional.',
    ],
    docs: [{ label: 'Field types', href: '/docs/concepts/field-types' }],
    keywords: ['add column', 'alter', 'migration', 'g field'],
  },
  {
    id: 'generate-seeder',
    name: 'grit generate seeder',
    category: 'Generate',
    summary: 'A seeder for an existing resource, with one example row or many faked ones.',
    example: 'grit generate seeder Product --faker --count 40',
    output: [
      '  ✓ apps/api/internal/database/products_seeder.go',
      '  ✓ Registered with Seed()',
      '',
      "  Run 'grit seed' to insert the rows.",
    ],
    purpose:
      'Reads the already-generated Go model and writes a seeder that matches its columns, registered in seed.go so `grit seed` picks it up.',
    useCases: [
      'A resource generated before you knew you wanted sample data.',
      'Filling a table so the admin, the dashboard widgets and a storefront have something to render.',
    ],
    files: [
      { path: 'apps/api/internal/database/products_seeder.go', status: 'created' },
      { path: 'apps/api/internal/database/seed.go', status: 'modified', note: 'registers the call' },
    ],
    flags: [
      { flag: '--faker', desc: 'Fill many rows with gofakeit instead of one example' },
      { flag: '--count int', desc: 'How many rows (default 10)' },
    ],
    notes: [
      'A resource with a required belongs_to refuses to seed until the parent has rows. The error names which parent.',
    ],
    keywords: ['sample data', 'fixtures', 'faker', 'demo data'],
  },
  {
    id: 'generate-sequence',
    name: 'grit generate sequence',
    category: 'Generate',
    summary: 'A gap-free sequential number: INV-202605-0001.',
    example: 'grit generate sequence Invoice',
    output: [
      '  ✓ apps/api/internal/sequence/sequence.go',
      '  ✓ apps/api/internal/services/invoice_sequence.go',
      '  ✓ Registered the counter with AutoMigrate',
      '',
      '  Call sequence.Next in the model\'s BeforeCreate hook.',
    ],
    purpose:
      'Human-facing reference numbers that auditors and customers read out loud. A row-locked counter, so two invoices created in the same millisecond cannot take the same number, and no gaps when one transaction rolls back.',
    useCases: [
      'Invoice, order, ticket and receipt numbers.',
      'Anything a person quotes over the phone.',
    ],
    files: [
      { path: 'apps/api/internal/sequence/sequence.go', status: 'created', note: 'the locked counter' },
      { path: 'apps/api/internal/services/invoice_sequence.go', status: 'created', note: 'the format' },
      { path: 'apps/api/internal/models/user.go', status: 'modified', note: 'AutoMigrate registry' },
    ],
    notes: [
      'Call sequence.Next directly from the model hook, never through a service. models importing services is an import cycle.',
    ],
    keywords: ['invoice number', 'counter', 'auto number', 'reference'],
  },
  {
    id: 'generate-perf',
    name: 'grit generate perf',
    category: 'Generate',
    summary: 'A k6 load test for this API.',
    example: 'grit generate perf',
    output: ['  ✓ perf/load.js', '  ✓ perf/README.md', '', '  Run: k6 run perf/load.js'],
    purpose:
      'Writes a k6 script pointed at your own endpoints, so the first load test is a command rather than an afternoon.',
    useCases: ['Before a launch.', 'Checking a suspected regression after a change to a hot path.'],
    files: [
      { path: 'perf/load.js', status: 'created' },
      { path: 'perf/README.md', status: 'created' },
    ],
    notes: ['Needs k6 installed separately.'],
    docs: [{ label: 'Benchmarks', href: '/docs/benchmarks' }],
    keywords: ['k6', 'load test', 'stress', 'performance'],
  },
  {
    id: 'remove-resource',
    name: 'grit remove resource',
    alias: 'grit rm resource',
    category: 'Generate',
    summary: 'Delete a resource and unpick every injection it made.',
    example: 'grit remove resource Product',
    output: [
      '  Removing resource: Product',
      '',
      '  ✗ apps/api/internal/models/product.go',
      '  ✗ apps/api/internal/services/product.go',
      '  ✗ apps/api/internal/handlers/product.go',
      '  ✗ packages/shared/schemas/product.ts',
      '  ✗ packages/shared/types/product.ts',
      '  ✗ apps/web/hooks/use-products.ts',
      '  ✗ apps/admin/resources/products/',
      '',
      '  Cleaning injections...',
      '  ✓ Removed from AutoMigrate',
      '  ✓ Removed routes',
      '  ✓ Removed from the permission catalog',
      '  ✓ Removed schema and type exports',
      '  ✓ Removed from the admin registry',
      '',
      '  ✅ Resource Product removed.',
    ],
    purpose:
      'The inverse of generate. Deleting the files by hand leaves fourteen injections behind, and the API stops compiling on the first one you miss.',
    useCases: [
      'A resource you named wrong.',
      'A modelling experiment you want to back out cleanly.',
    ],
    files: [
      { path: 'apps/api/internal/models/product.go', status: 'deleted' },
      { path: 'apps/api/internal/services/product.go', status: 'deleted' },
      { path: 'apps/api/internal/handlers/product.go', status: 'deleted' },
      { path: 'apps/api/internal/handlers/product_import.go', status: 'deleted' },
      { path: 'apps/api/internal/database/products_seeder.go', status: 'deleted' },
      { path: 'packages/shared/schemas/product.ts', status: 'deleted' },
      { path: 'packages/shared/types/product.ts', status: 'deleted' },
      { path: 'apps/web/hooks/use-products.ts', status: 'deleted' },
      { path: 'apps/admin/resources/products/products.ts', status: 'deleted' },
      { path: 'apps/admin/resources/products/products.custom.tsx', status: 'deleted', note: 'your customisations go too' },
      { path: 'apps/admin/app/(dashboard)/resources/products/', status: 'deleted' },
      { path: 'apps/api/internal/models/user.go', status: 'modified' },
      { path: 'apps/api/internal/routes/routes.go', status: 'modified' },
      { path: 'apps/api/internal/authz/permissions.go', status: 'modified' },
      { path: 'apps/api/internal/database/seed.go', status: 'modified' },
      { path: 'apps/api/internal/services/resource_stats_dispatch.go', status: 'modified' },
      { path: 'apps/admin/resources/index.ts', status: 'modified' },
      { path: 'packages/shared/schemas/index.ts', status: 'modified' },
      { path: 'packages/shared/types/index.ts', status: 'modified' },
      { path: 'packages/shared/constants/index.ts', status: 'modified' },
    ],
    notes: [
      'It does not drop the table. The rows are still there after the code is gone.',
      'Your .custom.tsx is deleted with everything else. Commit before you run it.',
    ],
    keywords: ['delete', 'destroy', 'undo generate', 'rm'],
  },

  // ── Add ────────────────────────────────────────────────────────────
  {
    id: 'add-variants',
    name: 'grit add variants',
    category: 'Add',
    summary: 'Product options, values and a combination matrix, attached to one resource.',
    example: 'grit add variants --resource Product',
    output: [
      '  Adding variants to Product',
      '',
      '  ✓ apps/api/internal/models/option.go',
      '  ✓ apps/api/internal/handlers/option_public.go',
      '  ✓ apps/api/internal/database/options_seeder.go',
      '  ✓ apps/api/internal/models/product_variant.go',
      '  ✓ apps/api/internal/services/product_variants.go',
      '  ✓ apps/api/internal/services/product_variants_test.go',
      '  ✓ apps/api/internal/handlers/product_variant.go',
      '  ✓ apps/api/internal/handlers/product_variant_public.go',
      '  ✓ apps/api/internal/database/product_variants_seeder.go',
      '  ✓ Registered 4 model(s) with AutoMigrate',
      '  ✓ Mounted the shared /options library',
      '  ✓ Mounted /products/:id/variants and /product-variants/:id',
      '  ✓ GET /api/v1/public/products/:key/variants (API key required)',
      '  ✓ Registered SeedProductVariants with grit seed',
      '  ✓ apps/admin/hooks/use-variants.ts',
      '  ✓ apps/admin/components/variants/variant-matrix.tsx',
      '  ✓ apps/admin/components/variants/option-library.tsx',
      '  ✓ apps/admin/resources/options/options.ts',
      '  ✓ Options added to the admin sidebar',
      '  ✓ Matrix editor attached to the Product detail page',
      '',
      '  Variants installed.',
    ],
    purpose:
      'One shirt in four colours and four sizes is one product and sixteen buyable things, each with its own stock. Five tables model that: options and their values shared across the shop, which options a given product offers, the combinations, and the join between them. A variant price is resolved rather than stored, so a product price change cannot leave stale copies behind.',
    useCases: [
      'Clothing, footwear, anything with a size.',
      'Electronics where memory or capacity changes the price and colour does not.',
      'Any catalogue where stock is per combination rather than per product.',
    ],
    files: [
      { path: 'apps/api/internal/models/option.go', status: 'created', note: 'shared: Option, OptionValue' },
      { path: 'apps/api/internal/models/product_variant.go', status: 'created' },
      { path: 'apps/api/internal/services/product_variants.go', status: 'created', note: 'price resolution, matrix generator' },
      { path: 'apps/api/internal/services/product_variants_test.go', status: 'created', note: 'six tests, run against your dialect' },
      { path: 'apps/api/internal/handlers/product_variant.go', status: 'created', note: 'admin endpoints' },
      { path: 'apps/api/internal/handlers/product_variant_public.go', status: 'created', note: 'storefront payload' },
      { path: 'apps/api/internal/handlers/option_public.go', status: 'created', note: 'shared published view types' },
      { path: 'apps/api/internal/database/options_seeder.go', status: 'created', note: 'Colour and Size' },
      { path: 'apps/api/internal/database/product_variants_seeder.go', status: 'created' },
      { path: 'apps/admin/hooks/use-variants.ts', status: 'created' },
      { path: 'apps/admin/components/variants/variant-matrix.tsx', status: 'created', note: 'the editor' },
      { path: 'apps/admin/components/variants/option-library.tsx', status: 'created' },
      { path: 'apps/admin/resources/options/options.ts', status: 'created', note: 'sidebar entry' },
      { path: 'apps/admin/app/(dashboard)/resources/options/page.tsx', status: 'created' },
      { path: 'apps/api/internal/models/user.go', status: 'modified', note: 'four models registered' },
      { path: 'apps/api/internal/routes/routes.go', status: 'modified' },
      { path: 'apps/api/internal/database/seed.go', status: 'modified' },
      { path: 'apps/api/internal/services/resource_stats_dispatch.go', status: 'modified' },
      { path: 'apps/api/internal/services/chart_dispatch.go', status: 'modified' },
      { path: 'apps/admin/resources/index.ts', status: 'modified' },
      { path: 'apps/admin/resources/products/products.custom.tsx', status: 'modified', note: 'matrix attached to the detail page' },
    ],
    flags: [{ flag: '--resource string', desc: 'The resource that offers variants (default Product)' }],
    notes: [
      'Run it again for a second resource and only that resource\'s tables are added. Options stay shared.',
      'Then `grit migrate` and `grit seed`: the seed writes a Colour and Size matrix so there is something on screen.',
      'Changing which options a product offers clears its combinations, because a variant is defined by the axes it was generated from.',
    ],
    docs: [{ label: 'Storefront guide', href: '/blog/build-a-storefront-with-grit' }],
    keywords: ['options', 'sizes', 'colours', 'sku', 'matrix', 'ecommerce', 'product options'],
  },
  {
    id: 'add-web-auth',
    name: 'grit add web-auth',
    category: 'Add',
    summary: 'Login, register and password-reset pages for apps/web, plus route protection.',
    example: 'grit add web-auth',
    output: [
      '  ✓ apps/web/app/(auth)/login/page.tsx',
      '  ✓ apps/web/app/(auth)/register/page.tsx',
      '  ✓ apps/web/app/(auth)/forgot-password/page.tsx',
      '  ✓ apps/web/app/(auth)/reset-password/page.tsx',
      '  ✓ apps/web/app/(auth)/callback/page.tsx',
      '  ✓ apps/web/components/ProtectedWebRoute.tsx',
      '  ✓ apps/web/components/UserMenu.tsx',
      '  ✓ apps/web/components/auth/AuthShell.tsx',
      '  ✓ apps/web/hooks/use-auth.ts',
      '  ✓ apps/web/lib/auth-provider.tsx',
      '  ✓ apps/web/middleware.ts',
      '',
      '  Web auth installed. Protect a route by adding it to the matcher in middleware.ts.',
    ],
    purpose:
      'The admin has auth out of the box; apps/web does not, because plenty of web apps are anonymous. This adds the whole customer-facing half: the five auth screens themed to match the project, a session marker the middleware reads, and a matcher-based guard so protecting /checkout is one line.',
    useCases: [
      'A storefront that needs accounts before checkout.',
      'Any public site with a members area.',
      'Adding OAuth callback handling to a web app that already has the API side.',
    ],
    files: [
      { path: 'apps/web/app/(auth)/login/page.tsx', status: 'created' },
      { path: 'apps/web/app/(auth)/register/page.tsx', status: 'created' },
      { path: 'apps/web/app/(auth)/forgot-password/page.tsx', status: 'created' },
      { path: 'apps/web/app/(auth)/reset-password/page.tsx', status: 'created' },
      { path: 'apps/web/app/(auth)/callback/page.tsx', status: 'created', note: 'OAuth return' },
      { path: 'apps/web/components/ProtectedWebRoute.tsx', status: 'created' },
      { path: 'apps/web/components/UserMenu.tsx', status: 'created' },
      { path: 'apps/web/components/auth/AuthShell.tsx', status: 'created', note: 'plus Atlas, Aurora and Pulse shells' },
      { path: 'apps/web/components/auth/SocialAuthButtons.tsx', status: 'created' },
      { path: 'apps/web/hooks/use-auth.ts', status: 'created' },
      { path: 'apps/web/lib/auth-provider.tsx', status: 'created' },
      { path: 'apps/web/lib/web-session.ts', status: 'created', note: 'the marker middleware reads' },
      { path: 'apps/web/middleware.ts', status: 'created', note: 'matcher-based protection' },
    ],
    notes: [
      'Files that already exist are left alone. If your (auth) folder is empty after a run, something was there first.',
      'The pages follow the project theme, so they match whatever --theme you scaffolded with.',
    ],
    keywords: ['login', 'register', 'signup', 'sign in', 'protect route', 'middleware', 'storefront auth'],
  },
  {
    id: 'add-role',
    name: 'grit add role',
    category: 'Add',
    summary: 'A new role across the Go enum, the Zod schema, the types and the constants.',
    example: 'grit add role MANAGER',
    output: [
      '  ✓ Added MANAGER to the role enum',
      '  ✓ packages/shared/constants/index.ts',
      '  ✓ packages/shared/schemas/user.ts',
      '  ✓ packages/shared/types/user.ts',
      '',
      '  Assign it in the admin under Settings, Roles.',
    ],
    purpose:
      'A role is named in four places that have to agree: the Go model, the Zod schema that validates a user update, the TypeScript union and the shared constants. Adding it by hand means finding all four.',
    useCases: [
      'A tier between ADMIN and USER.',
      'Per-department roles whose permissions you then set in the admin.',
    ],
    files: [
      { path: 'apps/api/internal/models/user.go', status: 'modified', note: 'the role enum' },
      { path: 'packages/shared/constants/index.ts', status: 'modified' },
      { path: 'packages/shared/schemas/user.ts', status: 'modified' },
      { path: 'packages/shared/types/user.ts', status: 'modified' },
    ],
    notes: ['Permissions for the new role are granted in the admin, not here.'],
    docs: [{ label: 'Roles and permissions', href: '/docs/backend/roles' }],
    keywords: ['rbac', 'permission', 'authorization', 'user role'],
  },
  {
    id: 'add-i18n',
    name: 'grit add i18n',
    category: 'Add',
    summary: 'next-intl on both frontends, translated API messages, en/fr/sw.',
    example: 'grit add i18n',
    output: [
      '  ✓ apps/api/internal/i18n/i18n.go',
      '  ✓ apps/api/internal/i18n/locales/en.json, fr.json, sw.json',
      '  ✓ apps/api/internal/middleware/locale.go',
      '  ✓ apps/api/internal/response/response.go',
      '  ✓ apps/web/i18n/request.ts + messages/',
      '  ✓ apps/admin/i18n/request.ts + messages/',
      '  ✓ language-switcher components',
      '',
      '  Run pnpm install for next-intl.',
    ],
    purpose:
      'Translation that reaches the API too. A validation error in French is the half most i18n setups skip, so the locale middleware and a translated response helper come with the frontend wiring.',
    useCases: [
      'Any product outside a single-language market.',
      'Adding a language later without retrofitting the API messages.',
    ],
    files: [
      { path: 'apps/api/internal/i18n/i18n.go', status: 'created' },
      { path: 'apps/api/internal/i18n/i18n_test.go', status: 'created' },
      { path: 'apps/api/internal/i18n/locales/en.json', status: 'created', note: 'plus fr.json and sw.json' },
      { path: 'apps/api/internal/middleware/locale.go', status: 'created' },
      { path: 'apps/api/internal/response/response.go', status: 'created', note: 'translated responses' },
      { path: 'apps/api/internal/handlers/i18n.go', status: 'created' },
      { path: 'apps/web/i18n/request.ts', status: 'created', note: 'plus lib/locale.ts and messages/' },
      { path: 'apps/admin/i18n/request.ts', status: 'created', note: 'plus lib/locale.ts and messages/' },
      { path: 'apps/web/components/language-switcher.tsx', status: 'created', note: 'and the admin one' },
      { path: 'apps/api/internal/routes/routes.go', status: 'modified' },
      { path: 'apps/web/next.config.ts', status: 'modified', note: 'and the admin config' },
      { path: 'apps/web/package.json', status: 'modified', note: 'next-intl' },
      { path: 'apps/web/app/layout.tsx', status: 'modified', note: 'provider' },
    ],
    notes: [
      'Run `pnpm install` afterwards: it adds next-intl to both frontends.',
      '`grit new --i18n` does the same thing at scaffold time.',
    ],
    keywords: ['translation', 'locale', 'next-intl', 'language', 'french', 'swahili'],
  },
  {
    id: 'add-offline',
    name: 'grit add offline',
    category: 'Add',
    summary: 'Offline-first sync: local mirror, outbox and version-checked conflicts.',
    example: 'grit add offline',
    output: [
      '  ✓ packages/sync/',
      '  ✓ IndexedDB, expo-sqlite and in-memory adapters',
      '  ✓ apps/api/internal/sync/',
      '',
      '  Mirror, outbox and conflict handling installed.',
    ],
    purpose:
      'The same mirror, outbox and version-checked conflict handling in TypeScript, over a storage interface, so web, mobile and desktop share one engine instead of three.',
    useCases: [
      'Field apps on bad connections.',
      'A desktop tool that has to keep working when the VPN drops.',
      'Mobile where a write must not be lost in a tunnel.',
    ],
    files: [
      { path: 'packages/sync/', status: 'created', note: 'the engine and its adapters' },
      { path: 'apps/api/internal/sync/', status: 'created', note: 'server-side policy' },
    ],
    notes: [
      'Sync policy is declared per resource and published at GET /api/sync/policy, so a client cannot keep a copy that drifts.',
      '`grit sync doctor` exists because every mistake in this area is silent.',
    ],
    docs: [{ label: 'Offline sync', href: '/docs/concepts/offline-sync' }],
    keywords: ['offline', 'sync', 'outbox', 'conflict', 'local first', 'indexeddb'],
  },
  {
    id: 'expose-form',
    name: 'grit expose form',
    category: 'Add',
    summary: 'A public page carrying a resource\'s form, at a path you choose.',
    example: 'grit expose form Product --to apps/web/app/submit-product/page.tsx',
    output: [
      '  ✓ apps/web/app/submit-product/page.tsx',
      '',
      '  The form posts to the Product endpoint with the same validation the admin uses.',
    ],
    purpose:
      'Takes the form the admin already renders for a resource and writes it into your public app as a page, so a contact form or a submission page is not a second implementation of the same validation.',
    useCases: [
      'Contact and enquiry forms.',
      'A public submission page: job applications, event registration.',
      'With --public-share, a form anyone with the link can fill in.',
    ],
    files: [{ path: 'apps/web/app/submit-product/page.tsx', status: 'created', note: 'the path you passed to --to' }],
    flags: [
      { flag: '--to string', desc: 'Destination path (required)' },
      { flag: '--public-share', desc: 'Submit via /api/public/forms/<token>/submit instead of the auth\'d hook' },
      { flag: '--token string', desc: 'FormShare token; falls back to NEXT_PUBLIC_FORM_TOKEN' },
      { flag: '--force', desc: 'Overwrite the destination' },
    ],
    notes: ['--to is required. Without it the command prints usage and exits.'],
    keywords: ['public form', 'contact form', 'embed', 'share'],
  },
  {
    id: 'expose-table',
    name: 'grit expose table',
    category: 'Add',
    summary: 'A public page carrying a resource\'s table.',
    example: 'grit expose table Product --to apps/web/app/catalogue/page.tsx',
    output: [
      '  ✓ apps/web/app/catalogue/page.tsx',
      '',
      '  The table reads the same endpoint with pagination and filters intact.',
    ],
    purpose:
      'The read-only twin of expose form: the admin table, in your public app, with its sorting, filtering and pagination already wired.',
    useCases: [
      'A public directory or listing.',
      'An internal dashboard outside the admin panel.',
    ],
    files: [{ path: 'apps/web/app/catalogue/page.tsx', status: 'created', note: 'the path you passed to --to' }],
    flags: [
      { flag: '--to string', desc: 'Destination path (required)' },
      { flag: '--force', desc: 'Overwrite the destination' },
    ],
    keywords: ['public table', 'listing', 'directory', 'embed'],
  },

  // ── Run ────────────────────────────────────────────────────────────
  {
    id: 'start',
    name: 'grit start',
    category: 'Run',
    summary: 'Run every service in the project, or one of them.',
    example: 'grit start',
    output: [
      '  Starting all services...',
      '',
      '  → API on http://localhost:8080',
      '  → Web on http://localhost:3000',
      '  → Admin on http://localhost:3001',
      '',
      '  GORM Studio: http://localhost:8080/studio',
      '  API Docs:    http://localhost:8080/docs',
    ],
    purpose:
      'One command instead of three terminals. Subcommands run a single piece when you only want that piece.',
    useCases: [
      'Normal development.',
      '`grit start server` alone while you work on the API.',
      '`grit start admin` when the frontend you care about is the admin panel.',
    ],
    files: [],
    flags: [
      { flag: 'server', desc: 'The Go API only' },
      { flag: 'client', desc: 'The frontends only' },
      { flag: 'web / admin / expo / desktop', desc: 'One app' },
    ],
    notes: ['Fails on a missing module if `pnpm install` has not been run.'],
    keywords: ['dev', 'run', 'serve', 'development server'],
  },
  {
    id: 'studio',
    name: 'grit studio',
    category: 'Run',
    summary: 'Open GORM Studio, the database browser.',
    example: 'grit studio',
    output: ['  GORM Studio: http://localhost:8080/studio', '', '  Opening in your browser...'],
    purpose:
      'A visual browser over the real tables, embedded in the API rather than a separate tool with its own connection string to get wrong.',
    useCases: [
      'Checking what a migration actually did.',
      'Editing a row without writing SQL.',
      'Confirming a seeder inserted what you meant.',
    ],
    files: [],
    keywords: ['database browser', 'gorm studio', 'tables', 'sql'],
  },
  {
    id: 'routes',
    name: 'grit routes',
    category: 'Run',
    summary: 'List every registered API route with its handler and group.',
    example: 'grit routes',
    output: [
      '  API Routes (apps/api/internal/routes/routes.go)',
      '',
      '  METHOD  PATH                                    HANDLER                       GROUP',
      '  ──────  ──────────────────────────────────────  ────────────────────────────  ──────',
      '  GET     /api/v1/products                        productHandler.List           protected',
      '  POST    /api/v1/products                        productHandler.Create         protected',
      '  GET     /api/v1/products/:id                    productHandler.Get            protected',
      '  PATCH   /api/v1/products/:id                    productHandler.Patch          protected',
      '  DELETE  /api/v1/products/:id                    productHandler.Delete         admin',
      '  GET     /api/v1/public/products                 productHandler.ListPublic     public',
      '  GET     /api/v1/public/products/:key            productHandler.GetPublic      public',
    ],
    purpose:
      'Reads routes.go and prints what is actually mounted, including which auth group each route sits in. Faster than reading a 900-line routes file, and it tells you the thing you usually want: is this endpoint public.',
    useCases: [
      'Checking a generated resource mounted where you expected.',
      'Auditing what is reachable without a token.',
      'Finding the handler behind a URL.',
    ],
    files: [],
    keywords: ['endpoints', 'api list', 'url', 'handler', 'public routes'],
  },
  {
    id: 'test',
    name: 'grit test',
    category: 'Run',
    summary: 'Run every test suite in the project.',
    example: 'grit test',
    output: [
      '  Running Go tests...',
      '  ok  	myapp/apps/api/internal/handlers	1.8s',
      '  ok  	myapp/apps/api/internal/services	0.9s',
      '',
      '  Running frontend tests...',
      '  ✓ apps/web  12 passed',
      '  ✓ apps/admin  8 passed',
      '',
      '  ✅ All suites passed.',
    ],
    purpose:
      'Go tests, Vitest on both frontends and Playwright end-to-end, behind one command, so CI and your terminal run the same thing.',
    useCases: ['Before a commit.', 'After an upgrade.', 'In CI.'],
    files: [],
    keywords: ['vitest', 'playwright', 'go test', 'ci', 'check'],
  },
  {
    id: 'ui',
    name: 'grit ui',
    category: 'Run',
    summary: 'Browse and install Grit UI components and blocks.',
    example: 'grit ui add ecommerce-product-grids-grid-with-ratings',
    output: [
      '  Installing into apps/web...',
      '',
      '  ✓ components/grit-ui/product-grids/grid-with-ratings.tsx',
      '',
      '  The file is yours. There is no package to upgrade.',
    ],
    purpose:
      'The component registry: marketing sections, ecommerce blocks, dashboard pieces. Installed as source into your repo, the shadcn model, so there is no version to track and no upstream fix that silently changes your page.',
    useCases: [
      'A product grid, a checkout, a hero section you would otherwise build from scratch.',
      'Browsing what exists before deciding to build.',
    ],
    files: [
      { path: 'apps/web/components/grit-ui/...', status: 'created', note: 'one file per block' },
    ],
    notes: [
      'A block prop you do not pass keeps its sample default. Pass an empty value for every prop your schema does not have, or the page will show sample ratings and colours you do not sell.',
      'Every scaffolded frontend ships a components.json, so `npx shadcn add` works with no prompts.',
    ],
    docs: [{ label: 'Grit UI', href: 'https://ui.gritframework.dev' }],
    keywords: ['components', 'blocks', 'shadcn', 'ui library', 'templates'],
  },
  {
    id: 'swap',
    name: 'grit swap',
    category: 'Run',
    summary: 'Replace an admin component everywhere at once.',
    example: 'grit swap button',
    output: [
      '  Swapping button across apps/admin...',
      '',
      '  ✓ 47 call sites updated',
      '  ✓ components/ui/button.tsx replaced',
      '',
      '  The old variants still resolve, so nothing breaks mid-swap.',
    ],
    purpose:
      'Changing a primitive in the admin means changing every call site, and a half-finished swap leaves two button styles in the same screen. This does the whole set in one pass.',
    useCases: [
      'Moving the admin onto your own design system.',
      'Replacing the stock input or button with a branded one.',
    ],
    files: [
      { path: 'apps/admin/components/ui/button.tsx', status: 'modified' },
      { path: 'every call site in apps/admin', status: 'modified' },
    ],
    keywords: ['design system', 'replace component', 'rebrand', 'theme'],
  },
  {
    id: 'mcp',
    name: 'grit mcp',
    category: 'Run',
    summary: 'Expose the project to AI coding agents over MCP.',
    example: 'grit mcp serve',
    output: [
      '  MCP server listening on stdio',
      '',
      '  Tools: list_resources, describe_resource, generate_resource, routes, migrate',
    ],
    purpose:
      'A Model Context Protocol server over the project, so an agent can list resources, read a schema and generate code through the real CLI rather than guessing at file layouts.',
    useCases: [
      'Driving Grit from Claude Code, Cursor or another MCP client.',
      'Letting an agent read the real route table instead of inferring it.',
    ],
    files: [],
    docs: [{ label: 'AI integration', href: '/docs/ai-integration' }],
    keywords: ['ai', 'claude', 'cursor', 'model context protocol', 'agent'],
  },

  // ── Data ───────────────────────────────────────────────────────────
  {
    id: 'migrate',
    name: 'grit migrate',
    category: 'Data',
    summary: 'Create and alter database tables from the Go models.',
    example: 'grit migrate',
    output: [
      '  Running migrations...',
      '',
      '  + created *models.User',
      '  + created *models.Product',
      '  + created *models.Option',
      '  + created *models.ProductVariant',
      '  ~ altered *models.Order (+2 columns)',
      '',
      '  Migration done: 4 table(s) created, 1 altered (+2 column(s)), 33 unchanged.',
      '  Migrations completed successfully.',
    ],
    purpose:
      'The model is the source of truth, so there are no migration files to write or order. GORM creates what is missing and alters what changed, and the command prints the before-and-after column diff rather than migrating silently.',
    useCases: [
      'After every generate resource or generate field.',
      'On a fresh clone, to build the schema.',
      'With --fresh, to drop everything and start over in development.',
    ],
    files: [],
    flags: [{ flag: '--fresh', desc: 'Drop every table first. Development only.' }],
    notes: [
      'It adds and alters. It does not drop a column you removed from a model.',
      '--fresh destroys data. There is no confirmation in a script.',
    ],
    keywords: ['migration', 'schema', 'automigrate', 'database', 'tables'],
  },
  {
    id: 'seed',
    name: 'grit seed',
    category: 'Data',
    summary: 'Run every registered seeder.',
    example: 'grit seed',
    output: [
      '  Seeded admin user: admin@example.com / admin123',
      '  Seeded 2 API keys (publishable + secret)',
      '  Wrote ../web/.env.local',
      '  Seeded 12 product',
      '  Seeded the Colour option with 4 values',
      '  Seeded the Size option with 4 values',
      '  Seeded 96 product variants across 6 products',
      '  Database seeded successfully.',
    ],
    purpose:
      'Fills a fresh database with enough to look at: an admin user you can log in as, API keys written into the frontend env files, and whatever your resources seed.',
    useCases: [
      'Right after `grit migrate` on a new checkout.',
      'Resetting a development database to a known state.',
    ],
    files: [
      { path: 'apps/web/.env.local', status: 'modified', note: 'publishable API key' },
      { path: 'apps/admin/.env.local', status: 'modified', note: 'publishable API key' },
    ],
    notes: [
      'Seeders skip when their table already has rows, so running it twice is safe.',
      'The default admin password only applies outside production. Set SEED_ADMIN_PASSWORD for a real deployment.',
    ],
    keywords: ['sample data', 'fixtures', 'demo', 'admin user'],
  },
  {
    id: 'sync',
    name: 'grit sync',
    category: 'Data',
    summary: 'Regenerate TypeScript types and Zod schemas from the Go models.',
    example: 'grit sync',
    output: [
      '  ✓ packages/shared/types/product.ts',
      '  ✓ packages/shared/schemas/product.ts',
      '  ✓ packages/shared/types/order.ts',
      '  ✓ packages/shared/schemas/order.ts',
      '',
      '  ✅ Synced 38 model(s) to TypeScript + Zod',
      '  ✅ Auto-added 1 field to admin resource files',
    ],
    purpose:
      'The Go model is the single definition of a shape. This projects it into TypeScript and Zod so the frontend cannot drift from the backend, and adds any new column to the admin table and form at the auto markers.',
    useCases: [
      'After editing a Go model by hand.',
      'After pulling changes that touched models.',
      'When the frontend types look stale.',
    ],
    files: [
      { path: 'packages/shared/types/*.ts', status: 'modified', note: 'one per model' },
      { path: 'packages/shared/schemas/*.ts', status: 'modified' },
      { path: 'apps/admin/resources/*/*.ts', status: 'modified', note: 'new columns at the auto markers' },
    ],
    flags: [{ flag: 'doctor', desc: 'Check offline sync health rather than regenerate types' }],
    notes: [
      'A resource file missing its grit:cols:auto-end markers is skipped with a warning rather than rewritten.',
    ],
    docs: [{ label: 'Type system', href: '/docs/concepts/type-system' }],
    keywords: ['types', 'zod', 'typescript', 'codegen', 'go to ts'],
  },
  {
    id: 'backup',
    name: 'grit backup',
    category: 'Data',
    summary: 'Back up the entire database to an archive.',
    example: 'grit backup',
    output: [
      '  Backing up...',
      '',
      '  ✓ backups/2026-08-20-093000.tar.gz (4.2 MB, 38 tables)',
    ],
    purpose:
      'A full dump you can restore, driven by the same code the admin\'s Data & Backup page uses, so a scheduled backup and a manual one produce the same artifact.',
    useCases: ['Before a risky migration.', 'Before an upgrade.', 'A scheduled snapshot.'],
    files: [{ path: 'backups/<timestamp>.tar.gz', status: 'created' }],
    keywords: ['dump', 'snapshot', 'export database', 'disaster recovery'],
  },
  {
    id: 'restore',
    name: 'grit restore',
    category: 'Data',
    summary: 'Restore the database from a backup archive.',
    example: 'grit restore backups/2026-08-20-093000.tar.gz',
    output: [
      '  Restoring from backups/2026-08-20-093000.tar.gz...',
      '',
      '  ✓ 38 tables restored',
      '  ✅ Restore complete.',
    ],
    purpose: 'The other half of backup. Replaces the current contents with the archive\'s.',
    useCases: ['Recovering from a bad migration.', 'Loading production data into a staging database.'],
    files: [],
    notes: ['Destructive. It replaces what is there now.'],
    keywords: ['recover', 'import database', 'rollback data'],
  },

  // ── Ship ───────────────────────────────────────────────────────────
  {
    id: 'deploy',
    name: 'grit deploy',
    category: 'Ship',
    summary: 'Cross-compile, upload, and configure systemd and Caddy with TLS.',
    example: 'grit deploy --host example.com',
    output: [
      '  Building for linux/amd64...',
      '  ✓ Binary built (18 MB)',
      '',
      '  Uploading to example.com...',
      '  ✓ Uploaded',
      '  ✓ systemd unit installed',
      '  ✓ Caddy configured, TLS issued',
      '',
      '  ✅ Live at https://example.com',
    ],
    purpose:
      'Turns a Grit project into a running server on a box you own: the Go binary cross-compiled, uploaded, supervised by systemd, and fronted by Caddy with automatic certificates.',
    useCases: [
      'A single VPS, which is the right answer for most projects for a long time.',
      'Staging environments that should match production.',
    ],
    files: [],
    notes: [
      'Work through the deployment checklist first. For a shop, three items are not optional: HTTPS everywhere, the webhook secret set in production, and backups on.',
    ],
    docs: [
      { label: 'Deploy command', href: '/docs/deployment/deploy-command' },
      { label: 'Checklist', href: '/docs/deployment/checklist' },
    ],
    keywords: ['ship', 'production', 'vps', 'systemd', 'caddy', 'tls', 'release'],
  },
  {
    id: 'compile',
    name: 'grit compile',
    category: 'Ship',
    summary: 'Build the desktop application executable.',
    example: 'grit compile',
    output: ['  Building desktop app...', '', '  ✓ build/bin/myapp.exe'],
    purpose: 'Wails build for the current platform, wired to the monorepo API.',
    useCases: ['Testing a desktop build before packaging it.'],
    files: [{ path: 'build/bin/', status: 'created', note: 'the executable' }],
    docs: [{ label: 'Desktop', href: '/docs/desktop' }],
    keywords: ['build', 'wails', 'exe', 'binary'],
  },
  {
    id: 'package',
    name: 'grit package',
    category: 'Ship',
    summary: 'Build a distributable desktop installer.',
    example: 'grit package',
    output: ['  Packaging...', '', '  ✓ dist/myapp-1.0.0-setup.exe'],
    purpose: 'The installer a user double-clicks: .exe on Windows, .app on macOS, a binary on Linux.',
    useCases: ['Shipping a desktop release.'],
    files: [{ path: 'dist/', status: 'created', note: 'the installer' }],
    notes: ['A signed release needs Authenticode or Apple Developer certificates, which are yours to obtain.'],
    keywords: ['installer', 'distribute', 'dmg', 'msi', 'release'],
  },
  {
    id: 'down',
    name: 'grit down',
    category: 'Ship',
    summary: 'Put the application in maintenance mode.',
    example: 'grit down',
    output: ['  Maintenance mode ON.', '  Requests answer 503 with a maintenance page.'],
    purpose:
      'A switch, not a redeploy. Every request answers 503 with a maintenance page while you run a migration that cannot happen under traffic.',
    useCases: ['A migration that has to run alone.', 'An incident where serving nothing beats serving wrong.'],
    files: [],
    keywords: ['maintenance', '503', 'offline mode', 'downtime'],
  },
  {
    id: 'up',
    name: 'grit up',
    category: 'Ship',
    summary: 'Bring the application back online.',
    example: 'grit up',
    output: ['  Maintenance mode OFF.', '  Serving normally.'],
    purpose: 'The other half of down.',
    useCases: ['After the migration finishes.'],
    files: [],
    keywords: ['maintenance', 'online', 'resume'],
  },

  // ── Meta ───────────────────────────────────────────────────────────
  {
    id: 'upgrade',
    name: 'grit upgrade',
    category: 'Meta',
    summary: 'Bring a project\'s scaffold files up to the current Grit version.',
    example: 'grit upgrade',
    output: [
      '  Upgrading project to Grit v3.169.0',
      '',
      '  ✓ Migration and seed tools updated',
      '  ✓ Web app updated',
      '  ✓ Admin panel updated (153 files)',
      '  ✓ Upgrade complete. Updated 184 files.',
      '',
      '  ⚠ Left alone, because you have edited these files since Grit wrote them:',
      '      apps/api/internal/database/seed.go',
      '      apps/web/next.config.ts',
      '      apps/admin/app/layout.tsx',
      '',
      '    grit upgrade --diff     # see what the new version would change',
      '    grit upgrade --force    # take the new version and lose your edits',
      '',
      '  Next steps:',
      '    pnpm install    # Install any new dependencies',
      '',
      '  Note: Resource definitions and API code were preserved.',
    ],
    purpose:
      'Updates the framework files inside an existing project. It reads .grit/manifest.json to know which generator wrote which file, at which version, with what content hash, so a file you edited is left alone rather than overwritten. Everything you wrote yourself, resource definitions and API code included, is preserved.',
    useCases: [
      'After `grit update` brings a new CLI.',
      'Picking up a fix that lives in a scaffolded file rather than in the CLI.',
      'Adding a file a newer version introduced, which an existing project has never had.',
    ],
    files: [
      { path: 'framework files under apps/', status: 'modified', note: 'the untouched ones only' },
      { path: '.grit/manifest.json', status: 'modified', note: 'new version and hash per file' },
    ],
    flags: [
      { flag: '--diff', desc: 'Also print what the new version would change in the files it left alone' },
      { flag: '--force', desc: 'Overwrite every file, including ones you have edited' },
    ],
    notes: [
      '--diff is NOT a dry run. It performs the upgrade and additionally prints the diff for the edited files it skipped. There is no preview-only mode, so commit first.',
      'Files you edited are skipped by content hash, not by timestamp, so reformatting one counts as editing it.',
      'Run `pnpm install` afterwards when a new version added a dependency.',
    ],
    docs: [{ label: 'Versioning', href: '/docs/versioning' }],
    keywords: ['update project', 'framework update', 'migrate version', 'dry run'],
  },
  {
    id: 'update',
    name: 'grit update',
    category: 'Meta',
    summary: 'Update the Grit CLI itself to the latest release.',
    example: 'grit update',
    output: [
      '  Current: v3.167.0',
      '  Latest:  v3.168.0',
      '',
      '  ✓ Downloaded and installed.',
      '  Run `grit upgrade --diff` in your projects to see what changed.',
    ],
    purpose: 'Replaces the binary in place from the latest GitHub release.',
    useCases: ['Picking up a new command or a CLI fix.'],
    files: [],
    notes: ['Updates the CLI, not your projects. `grit upgrade` does that.'],
    keywords: ['self update', 'cli version', 'install latest'],
  },
  {
    id: 'plugin',
    name: 'grit plugin',
    category: 'Meta',
    summary: 'Install and manage Grit plugins.',
    example: 'grit plugin add webhooks',
    output: [
      '  Installing webhooks...',
      '',
      '  ✓ apps/api/internal/webhooks/',
      '  ✓ Routes mounted',
      '  ✓ Registered with the event bus',
      '',
      '  ✅ Plugin webhooks installed.',
    ],
    purpose:
      'Optional modules that inject into the same markers a generator uses: webhooks, multi-tenancy, impersonation, saved views. Code in your repo, not a runtime dependency.',
    useCases: [
      'Standard Webhooks delivery without writing the retry logic.',
      'Multi-tenancy on a project that started single-tenant.',
    ],
    files: [
      { path: 'apps/api/internal/<plugin>/', status: 'created' },
      { path: 'apps/api/internal/routes/routes.go', status: 'modified' },
    ],
    flags: [
      { flag: 'add <name>', desc: 'Install a plugin' },
      { flag: 'list', desc: 'What is available and what is installed' },
      { flag: 'remove <name>', desc: 'Uninstall, unpicking its injections' },
    ],
    docs: [{ label: 'Plugins', href: '/docs/plugins' }],
    keywords: ['webhooks', 'multitenancy', 'extensions', 'modules'],
  },
  {
    id: 'version',
    name: 'grit version',
    category: 'Meta',
    summary: 'Print the CLI version.',
    example: 'grit version',
    output: ['grit version 3.168.0'],
    purpose: 'Which binary you are running, which is the first question in every bug report.',
    useCases: ['Filing an issue.', 'Checking `grit update` did something.'],
    files: [],
    keywords: ['-v', '--version', 'which version'],
  },
]

/** Lookup by id, for deep links. */
export function getCommand(id: string): CliCommand | undefined {
  return CLI_COMMANDS.find((c) => c.id === id)
}
