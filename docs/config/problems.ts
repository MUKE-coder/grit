/**
 * The twenty problems the talk is built on, each with what Grit does about it.
 *
 * The deck states a problem in six words and moves on, which is right for a
 * stage and useless to somebody who reads the slide afterwards and wants the
 * answer. Every problem slide links here, to the same problem with the
 * explanation, the commands and the honest limit attached.
 *
 * The wording of the problem itself is the canonical one from GRIT.md, reused
 * rather than rephrased: the talk, the docs and the specification should say
 * the same sentence, or the numbering stops meaning anything.
 */

export interface ProblemCommand {
  /** What this command is for, shown above it. Optional for a single line. */
  label?: string
  code: string
}

export interface ProblemLink {
  label: string
  href: string
}

export interface Problem {
  /** D1, B7. The identifier the talk and GRIT.md both use. */
  id: string
  /** Lowercased id, which is the URL. */
  slug: string
  group: 'developer' | 'client'
  /** The short name on the slide's sidebar. */
  label: string
  /** The slide's headline, verbatim. */
  headline: string
  /** The slide's one-paragraph statement of the problem. */
  symptom: string
  /** What actually goes wrong, in the detail a slide has no room for. */
  detail: string[]
  /** What Grit does about it. */
  answer: string[]
  commands: ProblemCommand[]
  /** What it does not do. Every entry has one, because every entry has one. */
  limit: string
  docs: ProblemLink[]
}

export const PROBLEMS: Problem[] = [
  {
    id: 'D1',
    slug: 'd1',
    group: 'developer',
    label: 'The copy-paste start',
    headline: 'You open the last project and start pasting.',
    symptom:
      'Auth. File storage. Email. You wrote them already. Writing JWT refresh from scratch a sixth time is worse than copying it.',
    detail: [
      'Every project starts by harvesting the last one: the auth handlers, the upload pipeline, the mailer, the middleware stack. It is faster than writing them again and it is how a bug travels. The fix you made in project four is not in projects one, two and three, and nobody remembers which ones are missing it.',
      'The copy also arrives without its reasons. The rate limit is there, the comment explaining why it is sixty seconds is not, so the next person tunes it to something that looks nicer.',
    ],
    answer: [
      'A new project has authentication, sessions, roles, uploads, mail, jobs, caching and the audit chain already in it, as framework code rather than as a starter template.',
      'They are upgraded by a command, so a fix made this month reaches a project started last year.',
      'None of it is opt-in. An opt-in security feature is a security feature most projects do not have.',
    ],
    commands: [
      { label: 'Everything below is already in this project', code: 'grit new billing' },
      { label: 'Take this month’s fixes into a project you started earlier', code: 'grit upgrade' },
    ],
    limit:
      'An upgrade cannot touch a file you edited: it reports it instead. The more framework code you rewrite, the more of the upgrade becomes a manual merge, which is the price of owning the code outright.',
    docs: [
      { label: 'What arrives in a new project', href: '/docs/getting-started/project-structure' },
      { label: 'Upgrading an existing project', href: '/docs/upgrades' },
    ],
  },
  {
    id: 'D2',
    slug: 'd2',
    group: 'developer',
    label: 'Wiring the monorepo',
    headline: 'A Go API beside Next.js. Or beside Expo.',
    symptom:
      'Workspaces, scripts, ports, env files, Docker. A day of undifferentiated work, at the exact moment enthusiasm is highest and patience is lowest.',
    detail: [
      'None of it is hard and all of it is fiddly: the workspace file, the shared package resolution, the proxy from the frontend to the API, the ports that collide, the four .env files that must agree, the compose file for Postgres and Redis.',
      'It is also the work that is least interesting on day one, when the idea is clearest, so it is where projects lose their first day and sometimes their momentum.',
    ],
    answer: [
      'One command produces the workspace, both frontends, the API, the shared package, the compose file and the environment files, already agreeing with each other.',
      'One command starts all of it together, in the right order, with the infrastructure.',
      'The architecture is a choice at creation: API only, one frontend, two, or a mobile app beside them.',
    ],
    commands: [
      { code: 'grit new billing' },
      { label: 'API, web, admin and the infrastructure, together', code: 'grit start' },
      { label: 'Or just the API', code: 'grit new billing --api' },
    ],
    limit:
      'The layout is a convention, not a configuration. A repository that must look different is one you will be adapting by hand, and the generator writes into apps/api, apps/web and apps/admin.',
    docs: [
      { label: 'Project structure', href: '/docs/getting-started/project-structure' },
      { label: 'Architecture modes', href: '/docs/concepts/architecture-modes' },
    ],
  },
  {
    id: 'D3',
    slug: 'd3',
    group: 'developer',
    label: 'Types written four times',
    headline: 'One invoice. Four descriptions of it.',
    symptom:
      'The Go struct, the TypeScript type, the validation schema, the mobile client. Kept in step by memory. Nothing fails to compile when they drift.',
    detail: [
      'The same record is described in the Go model, in a TypeScript interface, in a validation schema and again in whatever the mobile app holds. Four files, four authors, one shape.',
      'The failure is silent. Rename a column and the frontend keeps compiling, the form keeps submitting, and the field is quietly dropped on the way through. You find it when a customer asks why their phone number vanished.',
    ],
    answer: [
      'The Go model is the source. One command reads it and rewrites the TypeScript interface, the three Zod schemas and the admin resource fields.',
      'Renaming a column then makes the editor underline every screen that used the old name, before anything runs.',
      'The create schema drops the server-assigned columns and the update schema is the create schema partialled, so a PATCH cannot blank a field the form did not show.',
    ],
    commands: [
      { label: 'After changing a Go model', code: 'grit migrate\ngrit sync' },
      { label: 'Or let the generator write all of it from one description', code: 'grit generate resource Invoice --fields "number:string, total:money, due:date"' },
    ],
    limit:
      'Sync keeps the type, the schema and the admin in step. It does not update generated mobile or desktop screens; it names the fields they never saw and leaves them to you.',
    docs: [
      { label: 'The shared package', href: '/docs/frontend/shared-package' },
      { label: 'The type system', href: '/docs/concepts/type-system' },
    ],
  },
  {
    id: 'D4',
    slug: 'd4',
    group: 'developer',
    label: 'The same UI for every table',
    headline: 'CRUD hell.',
    symptom:
      'List, form, pagination, search, filters, sorting. Twenty tables, twenty times. It is not intellectually hard, which is exactly why it is expensive.',
    detail: [
      'Each table needs a list with server-side pagination, a sort that cannot be injected into, filters, a search, a create form, an edit form, validation wired to the API’s errors, a delete confirmation, and bulk actions once somebody asks.',
      'Doing it once is a week. Doing it twenty times is a quarter, and the twentieth is subtly different from the first because you learned something in between.',
    ],
    answer: [
      'Generating a resource writes the list, the detail page and the form, with pagination, sorting, filters, search, bulk actions, CSV import and Excel export already in them.',
      'The screen is a definition rather than a page, so changing a column is an edit to one file that survives regeneration.',
      'Every resource uses the same shared components, so an improvement to the table reaches all twenty screens.',
    ],
    commands: [
      { code: 'grit generate resource Invoice --fields "number:string, total:money, due:date, client:belongs_to"' },
      { label: 'Put a generated table or form in the public web app too', code: 'grit expose table Invoice --to apps/web/app/invoices/page.tsx' },
    ],
    limit:
      'It is an operator interface, not your product. A screen that joins four models with a chart and a date range is a page you write, using the generated hooks.',
    docs: [
      { label: 'The admin panel', href: '/docs/admin/overview' },
      { label: 'The data table', href: '/docs/admin/datatable' },
    ],
  },
  {
    id: 'D5',
    slug: 'd5',
    group: 'developer',
    label: 'Testing performance',
    headline: 'Every tool needs a day before it tells you anything.',
    symptom:
      'Load tests, profiling, tracing, query analysis. So they all wait for a quiet week, and the quiet week does not come.',
    detail: [
      'Profiling a Go service is well supported and still a project: wire pprof, get a load generator, produce realistic traffic, read a flame graph, correlate it with the queries that ran.',
      'Because the setup costs a day, it is done after the complaint rather than before it, which is the most expensive possible order.',
    ],
    answer: [
      'Request tracing, the queries each request ran, runtime metrics and errors with their request attached are mounted at /pulse in every project.',
      'The slow query list is sorted by total time rather than worst case, which is what finds the cheap query run ten thousand times.',
      'Allocation budgets ship as tests, so a change that allocates more on the hot path fails a build rather than being noticed in a year.',
    ],
    commands: [
      { label: 'Already running', code: 'grit start\n# http://localhost:8080/pulse' },
      { label: 'The budgets, as ordinary tests', code: 'cd apps/api && go test ./internal/handlers/ -run TestAllocationBudgets' },
    ],
    limit:
      'Pulse explains load; it does not generate it. A load generator belongs out of process, which is why the docs point at k6 rather than shipping one in the API.',
    docs: [
      { label: 'Pulse', href: '/docs/backend/pulse' },
      { label: 'Performance', href: '/docs/concepts/performance' },
    ],
  },
  {
    id: 'D6',
    slug: 'd6',
    group: 'developer',
    label: 'No window into the data',
    headline: 'Can I just look at my data?',
    symptom:
      'Browse the tables. Run a query. Draw the ERD. Export the schema. All possible, never already there.',
    detail: [
      'Answering “did that actually save” means a database client, credentials for an environment you may not have, and a connection that works from where you are standing.',
      'So the question gets answered with a log line instead, and half of all bug reports are a row that is exactly what the code wrote.',
    ],
    answer: [
      'A database browser is mounted at /studio in every project: tables, rows, filters, foreign keys followed, a query runner and an ERD drawn from the schema.',
      'It is the same tool in development and against a deployed database, so the question is answered the same way everywhere.',
    ],
    commands: [
      { label: 'Already running with the API', code: 'grit start\n# http://localhost:8080/studio' },
      { label: 'Or open it directly', code: 'grit studio' },
    ],
    limit:
      'Studio sits underneath the ownership scoping. A row you can see there and cannot fetch through the API is usually the scoping working, not a bug.',
    docs: [{ label: 'Studio', href: '/docs/systems/studio' }],
  },
  {
    id: 'D7',
    slug: 'd7',
    group: 'developer',
    label: 'The undocumented API',
    headline: 'What does this endpoint return?',
    symptom:
      'The honest answer is read the handler. Hand-written OpenAPI is accurate for one afternoon, and nobody is assigned to keep it true.',
    detail: [
      'API documentation written by hand is a second description of the same thing, which makes it a thing that drifts. It is correct on the day it is written and wrong by the end of the sprint.',
      'The cost lands on whoever is integrating: a mobile developer, a client’s team, or you in four months.',
    ],
    answer: [
      'The reference is generated from the registered routes and the request and response types, and served at /docs by the API itself.',
      'Generate a resource, reload the page, and its thirteen routes are there with their parameters and responses.',
      'It cannot drift, because nobody maintains it.',
    ],
    commands: [
      { label: 'Already running', code: 'grit start\n# http://localhost:8080/docs' },
      { label: 'Every route, in the terminal', code: 'grit routes' },
    ],
    limit:
      'It documents shape, not meaning. What a status transition implies for your business is prose you still write.',
    docs: [{ label: 'API documentation', href: '/docs/backend/api-docs' }],
  },
  {
    id: 'D8',
    slug: 'd8',
    group: 'developer',
    label: 'Import and export',
    headline: 'Can we upload a spreadsheet?',
    symptom:
      'Every client asks eventually. It looks like an afternoon until you meet column mapping, partial failure, forty thousand rows and the one date the parser hates.',
    detail: [
      'The afternoon version reads a CSV and inserts rows. The real version maps columns from a header nobody standardised, validates each row, reports row 17 without discarding rows 1 to 16, handles a file larger than memory, and tells the user what to fix.',
      'Export has the same shape: it must honour the filters the user is looking at, or the numbers in the file disagree with the numbers on the screen.',
    ],
    answer: [
      'Every generated resource gets CSV import and Excel export, with a dry run that validates before anything is written and per-row errors rather than one failure for the file.',
      'Export honours the filters currently applied to the list.',
      'Both are ordinary Go in your project: internal/handlers/<resource>_import.go and its service half.',
    ],
    commands: [
      { label: 'Generated with the resource', code: 'grit generate resource Invoice --fields "number:string, total:money"' },
      { label: 'Large imports, bounded', code: 'IMPORT_CONCURRENCY=4 grit start' },
    ],
    limit:
      'Column mapping is matched from the header. A file whose columns are named in another language, or not named at all, needs a mapping you write.',
    docs: [{ label: 'Import and export', href: '/docs/admin/datatable' }],
  },
  {
    id: 'D9',
    slug: 'd9',
    group: 'developer',
    label: 'No data to build against',
    headline: 'You test on ten rows. The client has fifty thousand.',
    symptom:
      'Writing a seeder that produces realistic volume is a day’s work nobody budgets, so every list page is judged on hand-typed rows.',
    detail: [
      'Everything in a list page works at ten rows. Sorting is instant, pagination has one page, the N+1 costs ten queries instead of ten thousand, and the missing index costs nothing. You find out at the same moment your client does.',
      'Development email has the same hole: it goes to a provider you cannot open, so the reset link that never arrives is a mystery rather than a line in a terminal.',
    ],
    answer: [
      'Every resource gets a factory and a seeder, so realistic volume is one command rather than a day.',
      'Seeding runs in batches with the heap flat, so a hundred thousand rows is a command and not an out-of-memory error.',
      'Mail has a log driver and a Mailhog default, so a reset link is never lost in development.',
    ],
    commands: [
      { label: 'Realistic volume, before you look at the screen', code: 'grit seed Invoice --count 25000' },
      { label: 'Mail you can actually read', code: 'MAIL_MAILER=log grit start' },
    ],
    limit:
      'A faker produces plausible rows, not your client’s distribution. A table where ninety percent of rows belong to one customer behaves differently, and that shape is yours to seed.',
    docs: [
      { label: 'Seeders', href: '/docs/backend/seeders' },
      { label: 'Email', href: '/docs/batteries/email' },
    ],
  },
  {
    id: 'D10',
    slug: 'd10',
    group: 'developer',
    label: 'Migrations with no way back',
    headline: 'A schema change on production, with no way back.',
    symptom:
      'Adding a column is the easy half. The missing half is the record: which statements in that run actually applied, and what undoes them.',
    detail: [
      'A migration that fails on statement four leaves a database in a state no version of the code expects. Without a record of what applied, recovery is somebody reading SQL at midnight and guessing.',
      'The forward migration is the part everybody writes. The run identifier, the rollback and the refusal to do something unsafe are the parts that are missing when it matters.',
    ],
    answer: [
      'Each migration run has an identifier and a record of the statements that applied, so a failed run can be reversed rather than reconstructed.',
      'A change that would lose data is refused rather than performed, and says what it would have dropped.',
      'Backup and restore are commands, so the state before a migration is recoverable.',
    ],
    commands: [
      { code: 'grit migrate' },
      { label: 'What ran, and what it would take back', code: 'grit migrate status\ngrit migrate down' },
      { label: 'Before anything irreversible', code: 'grit backup' },
    ],
    limit:
      'A rollback restores schema, not meaning. A migration that rewrote values on the way through needs a data migration to reverse it, and the framework cannot invent one.',
    docs: [
      { label: 'Migrations', href: '/docs/backend/migrations' },
      { label: 'Backups', href: '/docs/batteries/backups' },
    ],
  },
  {
    id: 'D11',
    slug: 'd11',
    group: 'developer',
    label: 'The frontend nobody wants',
    headline: 'You are fluent in Go. Then you open the React project.',
    symptom:
      'Nobody funds a frontend developer for an internal admin panel. So you learn just enough React to build screens you are not proud of, slowly, on every project.',
    detail: [
      'The backend is the part you are good at and the part the client cannot see. The admin panel is the part they judge the work by, and it is built by whoever is left.',
      'The result is screens that work and look like they were built under protest, and a developer who spent a third of the project not doing what they are best at.',
    ],
    answer: [
      'The admin panel is generated, dark-themed and consistent across every resource, from the same description that produced the API.',
      'Customising a screen is editing one TypeScript file that is never regenerated, so you change what you care about rather than building from scratch.',
      'The generated hooks are available to any page you write yourself, typed from the Go model.',
    ],
    commands: [
      { code: 'grit generate resource Invoice --fields "number:string, total:money"' },
      { label: 'The screen, as a file you own', code: '# apps/admin/app/_dashboard/resources/invoices/index.tsx' },
    ],
    limit:
      'It is an operator interface and it is opinionated. Reskinning it entirely is a real project, and you will be merging your changes against upgrades afterwards.',
    docs: [
      { label: 'Admin resources', href: '/docs/admin/resources' },
      { label: 'Custom pages', href: '/docs/admin/custom-pages' },
    ],
  },
  {
    id: 'D12',
    slug: 'd12',
    group: 'developer',
    label: 'Paying a model to write it',
    headline: 'Asking a model to write it is not free.',
    symptom:
      'Seventy-four thousand lines is about three quarters of a million output tokens, before the first correction. You pay it again on the next project, and it writes something different.',
    detail: [
      'A model can produce the monorepo, the API, the admin panel and the auth. It produces them again on the next project, at the same cost, in a different shape, with different names and a different idea of where a rule belongs.',
      'The variance is the expensive part, not the tokens. Two projects generated a month apart cannot share a fix, because they do not share a structure.',
    ],
    answer: [
      'The scaffolding is deterministic: the same description produces the same files, this month and next year.',
      'Tokens go on the feature rather than on the shape, because the shape is already there.',
      'The structure is documented and identical across projects, so a model working inside it has a convention to follow rather than one to invent.',
    ],
    commands: [
      { label: 'The same output every time', code: 'grit new billing\ngrit generate resource Invoice --fields "number:string, total:money"' },
      { label: 'Give your AI tool the project’s own conventions', code: 'grit mcp' },
    ],
    limit:
      'Determinism is about the shape, not about your domain. The rules that make the application yours are still written by somebody, and that somebody may well be a model.',
    docs: [
      { label: 'Working with AI', href: '/docs/ai-skill' },
      { label: 'Code generation', href: '/docs/concepts/code-generation' },
    ],
  },

  {
    id: 'B1',
    slug: 'b1',
    group: 'client',
    label: 'Can different people see different things?',
    headline: '“Can different people see different things?”',
    symptom:
      'Roles, then permissions, then per-record ownership. It is always one layer deeper than the first conversation suggests.',
    detail: [
      'The first answer is roles: an admin and a user. The second is permissions, because an editor may publish but not delete. The third is ownership, because a salesperson may read accounts, but only theirs.',
      'The third layer is the one that is missed, and it is the one that becomes an incident: changing a number in a URL and seeing somebody else’s invoice.',
    ],
    answer: [
      'Roles and permissions are in the project before you start, with the sidebar, the forms and the seeded permissions all knowing about a new role.',
      'Per-record ownership is a flag on the resource. Queries are scoped by the actor in the context, so the rule holds for the API, an export, an import and a job.',
      'A record you may not see answers 404 rather than 403, so IDs cannot be enumerated.',
    ],
    commands: [
      { label: 'A new role, everywhere it has to be known', code: 'grit add role MANAGER' },
      { label: 'Rows belong to the user who created them', code: 'grit generate resource Note --owned-by user --fields "title:string, body:text"' },
      { label: 'Or to an organisation', code: 'grit generate resource Invoice --tenant-owned --fields "number:string, total:money"' },
    ],
    limit:
      'Permission design is yours. That a regional manager sees their region, and what happens when somebody changes region, is domain modelling that nothing generates.',
    docs: [
      { label: 'Roles and permissions', href: '/docs/backend/rbac' },
      { label: 'Authorization', href: '/docs/security/authorization' },
    ],
  },
  {
    id: 'B2',
    slug: 'b2',
    group: 'client',
    label: 'What happens if we get hacked?',
    headline: '“What happens if we get hacked?”',
    symptom:
      'Rate limits. Brute-force lockouts. A second factor. Encryption for the columns that deserve it. A credible answer when a credential leaks.',
    detail: [
      'The question is asked once, early, and answered with a sentence. The work behind the sentence is weeks, and it is invisible, so it is the first thing cut when a feature is due.',
      'The sign-in form alone needs a rate limit, a lockout, a growing delay, and a locked account that answers exactly like a wrong password, or it becomes an oracle for discovering who has an account.',
    ],
    answer: [
      'Sessions are server-side, so a device can be revoked and a reused refresh token kills the chain it came from.',
      'The sign-in form has four layers, including the one people miss: a locked account answers like a wrong password, on the same timing.',
      'A security suite with a firewall, rate limiting and an anomaly log is mounted at /sentinel, with its threat log in your own database.',
    ],
    commands: [
      { label: 'Check the configuration before you deploy', code: 'grit doctor' },
      { label: 'Already running', code: 'grit start\n# http://localhost:8080/sentinel' },
    ],
    limit:
      'None of it protects a leaked admin password that nobody rotates, or a server nobody patches. A framework can close the application’s doors, not the operational ones.',
    docs: [
      { label: 'Authentication', href: '/docs/backend/authentication' },
      { label: 'The defender’s handbook', href: '/docs/security/defenders-handbook' },
    ],
  },
  {
    id: 'B3',
    slug: 'b3',
    group: 'client',
    label: 'Does it work without internet?',
    headline: '“Does it work without internet?”',
    symptom:
      'In a lot of rooms this is a nice-to-have you can talk a client out of. On this continent it is the requirement, and it is weeks of work.',
    detail: [
      'Offline means a local mirror of the data, a queue of writes made while disconnected, a replay when the connection returns, and a policy for what happens when the same record changed on both sides.',
      'Every one of those is a design decision, and getting the last one wrong silently loses somebody’s work.',
    ],
    answer: [
      'Offline support is a command: a local mirror, an outbox for writes, a replay, and a declared conflict policy per resource.',
      'The policy is part of the resource definition rather than a second file, so what syncs and how it resolves lives with the resource it describes.',
    ],
    commands: [
      { code: 'grit add offline' },
      { label: 'Declared with the resource', code: 'sync:\n  mode: offline_first\n  conflict: server_wins\n  max_offline_age: 72h' },
    ],
    limit:
      'Conflict resolution is a business decision wearing a technical hat. Server-wins and client-wins are both wrong for some records, and deciding which is which is yours.',
    docs: [{ label: 'Offline sync', href: '/docs/concepts/offline-sync' }],
  },
  {
    id: 'B4',
    slug: 'b4',
    group: 'client',
    label: 'How will we know when it breaks?',
    headline: '“How will we know when it breaks?”',
    symptom:
      'Errors, latency, resource use. And the slow page nobody can find, because no single query in it is slow.',
    detail: [
      'The honest failure is not a crash. It is a page that takes four seconds because it runs two hundred queries, none of which is slow enough to appear in a slow-query log.',
      'Without a trace that ties a request to the queries it ran, that page is diagnosed by guessing.',
    ],
    answer: [
      'Every request is traced with the queries it ran and what each cost, so the two hundred show up as two hundred.',
      'A health endpoint reports each dependency separately and distinguishes “off because nobody configured it” from “configured and broken”.',
      'Every request carries an ID, returned in a header and attached to every log line, so a support message becomes one search.',
    ],
    commands: [
      { label: 'Already running', code: 'grit start\n# http://localhost:8080/pulse' },
      { label: 'What a load balancer should read', code: 'curl localhost:8080/health' },
      { label: 'What the dev servers printed', code: 'grit logs --errors' },
    ],
    limit:
      'It tells you what is happening in the process. It does not page anybody: alerting is a service you still choose and wire up.',
    docs: [
      { label: 'Pulse', href: '/docs/backend/pulse' },
      { label: 'Health checks', href: '/docs/backend/health' },
    ],
  },
  {
    id: 'B5',
    slug: 'b5',
    group: 'client',
    label: 'Is any of this backed up?',
    headline: '“Is any of this backed up?”',
    symptom:
      'Asked once, answered with a yes, and verified by nobody until the morning it matters.',
    detail: [
      'A backup that has never been restored is a file, not a backup. The first restore attempt is not the moment to discover it is missing the uploads, or the schema, or three hours.',
      'The half nobody budgets is the restore: the command, the practice, and the knowledge of exactly what is lost between the last backup and the failure.',
    ],
    answer: [
      'Backup and restore are both commands, in every project, including the uploads and not only the database.',
      'The admin has a data page that shows what exists and when it was taken.',
      'Restoring into a scratch environment is the documented habit, not an emergency procedure.',
    ],
    commands: [
      { code: 'grit backup' },
      { label: 'The half that counts', code: 'grit restore backup-2026-10-22.zip' },
    ],
    limit:
      'Where the archives live and how long they are kept is yours. A backup on the same disk as the database is a copy, not a backup.',
    docs: [{ label: 'Backups', href: '/docs/batteries/backups' }],
  },
  {
    id: 'B6',
    slug: 'b6',
    group: 'client',
    label: 'What if it gets busy?',
    headline: '“What if it gets busy?”',
    symptom:
      'Caching, read replicas, background work, bounded memory. And knowing which of those you actually need, which is the hard part.',
    detail: [
      'Scaling is four different problems that share a word. A slow report, a thundering herd, a queue that cannot keep up and a process that runs out of memory each need a different answer.',
      'Buying the wrong one is worse than buying none: a cache in front of the wrong query hides the problem and adds a second one.',
    ],
    answer: [
      'Caching, background jobs, a queue, read-replica support and bounded batch processing are in the project already.',
      'Response caching is mounted only where its key is safe, on the public group, and nowhere else.',
      'Pulse tells you which of the four you are actually facing before you choose.',
    ],
    commands: [
      { label: 'Work that must not happen in a request', code: 'grit generate job SendInvoiceReminders --cron "0 9 * * *"' },
      { label: 'What the process is doing under load', code: 'grit start\n# http://localhost:8080/pulse' },
    ],
    limit:
      'This is the one the talk marks in a different colour. The tools are here; which one your bottleneck needs stays your call, and anybody who says otherwise is selling something.',
    docs: [
      { label: 'Scaling', href: '/docs/scaling' },
      { label: 'Caching', href: '/docs/batteries/caching' },
    ],
  },
  {
    id: 'B7',
    slug: 'b7',
    group: 'client',
    label: 'Who changed this?',
    headline: '“Who changed this?”',
    symptom:
      'Not a timeline that anyone with database access can quietly edit. An append-only record of every write, and a way to prove it has not been touched.',
    detail: [
      'An audit table that can be updated is a log that proves nothing. The question behind the question is whether the record can be trusted by somebody who does not trust you.',
      'The usual implementation is a trigger writing to a table that the same database user can edit, which answers the easy version of the question and not the real one.',
    ],
    answer: [
      'Writes are recorded through the event bus, so adding an audited action is not a line somebody has to remember to add to a handler.',
      'The log is hash-chained and append-only, enforced by the routes, by a GORM guard and by a database trigger, and it has a verify action.',
      'Resources can be append-only themselves, and reads of sensitive resources can be recorded too.',
    ],
    commands: [
      { label: 'A ledger that cannot be edited or deleted', code: 'grit generate resource LedgerEntry --append-only --fields "kind:string, amount:money"' },
      { label: 'Record who read it, as well as who wrote it', code: 'grit generate resource Record --audit-reads --fields "subject:string, notes:text"' },
    ],
    limit:
      'It proves the chain has not been altered. It cannot prove that what was written was true, and somebody with the server can still stop the process.',
    docs: [{ label: 'The audit log', href: '/docs/security/audit-log' }],
  },
  {
    id: 'B8',
    slug: 'b8',
    group: 'client',
    label: 'Can it take a ten megabyte photo?',
    headline: '“Our staff photograph things on their phones.”',
    symptom:
      'Rotated, ten megabytes each, with the GPS of a client site in the EXIF. Uploaded over a connection that drops.',
    detail: [
      'The naive version stores what arrives and serves it back at full size, so a list of fifty rows is half a gigabyte of thumbnails.',
      'It also stores what the photograph knows: the orientation flag that makes a portrait picture sideways, and the coordinates of the house it was taken in.',
    ],
    answer: [
      'Images are resized, re-encoded, stripped of metadata and given a thumbnail on the way in. A 5.3 MB phone photo stores at about 147 KB including the thumbnail.',
      'The stored reference carries the dimensions, the type and the thumbnail URL, so a page can reserve space and a list can use the small version.',
      'A decompression bomb is refused from its header, before any pixels are allocated.',
    ],
    commands: [
      { label: 'One field, the whole pipeline', code: 'grit generate resource Site --fields "name:string, photo:file, gallery:files"' },
      { label: 'No object store needed to start', code: 'STORAGE_DRIVER=local grit start' },
    ],
    limit:
      'Resumable uploads over a dropping connection are the one part this does not solve for you: large files go direct to storage with a presigned URL, and the retry behaviour is the client’s.',
    docs: [
      { label: 'Media and uploads', href: '/docs/batteries/media' },
      { label: 'Storage', href: '/docs/batteries/storage' },
    ],
  },
]

export const DEVELOPER_PROBLEMS = PROBLEMS.filter((p) => p.group === 'developer')
export const CLIENT_PROBLEMS = PROBLEMS.filter((p) => p.group === 'client')

export function getProblem(slug: string): Problem | undefined {
  return PROBLEMS.find((p) => p.slug === slug)
}
