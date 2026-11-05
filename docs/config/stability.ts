/*
 * What each subsystem's stability badge says, and the only place that says it.
 *
 * This lived in docs/app/docs/stability/page.tsx, where it was the page's own
 * local array. The capability matrix made its own call about the same
 * subsystems and the two disagreed: two-factor, passkeys and RBAC read beta on
 * the stability page while counting as ready-on-first-run in the matrix, which
 * is the project telling a reader two things depending on which page they
 * found (grit#166).
 *
 * A capability row points here by id. stabilityOf resolves it, and
 * scripts/check-stability-pointers.py fails the build on a pointer that names
 * no area, so the two cannot drift apart again by a rename.
 */

export type Status = 'stable' | 'beta' | 'new'

export const statusStyle: Record<Status, string> = {
  stable: 'bg-emerald-500/10 text-emerald-400 border-emerald-500/20',
  beta: 'bg-amber-500/10 text-amber-400 border-amber-500/20',
  new: 'bg-sky-500/10 text-sky-400 border-sky-500/20',
}

export type StabilityArea = {
  /** stable across this file and the capability matrix, which points at it */
  id: string
  area: string
  status: Status
  guaranteed: string
  proof: string
  yours: string
  history?: string
}

// One row per subsystem. The rule for this table: a claim in "Grit guarantees"
// has to name the test that would fail if it stopped being true. Where there is
// no such test, the claim belongs in "Your responsibility" instead, however
// much the code looks right.
export const STABILITY_AREAS: StabilityArea[] = [
  {
    id: 'auth',
    area: 'Auth: JWT, sessions, refresh rotation',
    status: 'stable',
    guaranteed:
      'Every refresh token is backed by a session row, rotation detects replay, and a password change signs out every device.',
    proof: 'handlers/auth_test.go, services/session_test.go, services/password_reset_test.go (41 handler tests, 75 service tests ship in your project)',
    yours: 'Set JWT_SECRET to 32+ random bytes and keep it out of the repository. grit doctor fails if it is short, empty or a placeholder.',
    history: 'Two logins in the same second once produced identical refresh tokens (fixed in v3.86.0: every token carries a jti).',
  },
  {
    id: 'twofa',
    area: 'Two-factor, passkeys, recovery codes',
    status: 'beta',
    guaranteed: 'TOTP is RFC 6238 (HMAC-SHA1 by the spec), recovery codes are single-use, and passkey challenges are bound to a session.',
    proof: 'services/passkey_test.go, services/recovery_test.go, crypto/field_test.go',
    yours: 'Decide your own lockout and recovery policy, and test the flow with the authenticator app your users actually have.',
  },
  {
    id: 'rbac',
    area: 'RBAC: roles, permissions, the staff group',
    status: 'beta',
    guaranteed:
      'A custom role reaches only the endpoints its permissions name, a revoked permission stops working everywhere, and a route that names no permission stays ADMIN-only, so one that forgets fails closed.',
    proof: 'authz/grants_test.go, authz/permissions_test.go, handlers/role_test.go, services/user_test.go, plus a live check that a staff account holding notes.delete still cannot delete another user’s row',
    yours: 'Name a permission on every route you add yourself. The generator does it for generated resources; a hand-written route is yours.',
    history: 'A custom role could not reach a single admin endpoint until v3.220.0; a revoked permission kept working on other replicas until v3.218.0.',
  },
  {
    id: 'crud',
    area: 'Generated CRUD: services, handlers, pagination',
    status: 'stable',
    guaranteed:
      'Handlers run no queries; the service owns every read and write, takes a context, and whitelists what a client may search, sort, filter and patch, because those names reach SQL.',
    proof: 'internal/generate tests in the CLI (492 test functions across 99 files), an import-versus-use check on every generated file, and 57 live checks against a running app on Postgres 15, 16 and 17',
    yours: 'Put your own queries in the service, not the handler, so a job and a route cannot disagree.',
    history: 'Handlers ran all of their own queries until v3.224.0, and the service written beside them was never called.',
  },
  {
    id: 'owned',
    area: 'Owned resources (--owned-by)',
    status: 'stable',
    guaranteed:
      'List, export, get, PDF, update, patch, delete and bulk are all scoped to the caller; somebody else’s row is 404 rather than 403; the owner comes from the session and never from the request body. ADMIN is exempt.',
    proof: 'Live checks cover every one of those paths, including the staff-permission case, plus generator tests that fail if a method loses its scoping. grit doctor reports the same thing in your project.',
    yours: 'Add --owned-by when rows belong to a user. A resource that merely references a user is not scoped, and grit doctor asks about it.',
    history: 'The export, PDF, patch and bulk routes were unscoped once: a second ordinary account printed another user’s record and rewrote it.',
  },
  {
    id: 'locking',
    area: 'Optimistic locking (If-Match)',
    status: 'stable',
    guaranteed:
      'A read returns the version as an ETag, a write against a stale version is 409 naming the current one, and of twenty simultaneous writes on one version exactly one lands.',
    proof: 'concurrency/concurrency_test.go in your project, generator tests, and a 20-way race in the live suite on all three Postgres versions',
    yours: 'Send If-Match from your own clients. Without it the last write still wins, by design.',
    history: 'Added in v3.219.0, after two people saving the same record meant the second silently won.',
  },
  {
    id: 'encryption',
    area: 'Encryption at rest (:encrypted fields)',
    status: 'beta',
    guaranteed:
      'An encrypted column is ciphertext in the database and plaintext through the API, and it is kept out of search, sort and filter lists, where ciphertext would match nothing.',
    proof: 'crypto/field_test.go, crypto/map_update_test.go, and a live check that reads the column straight out of Postgres',
    yours:
      'Set FIELD_ENCRYPTION_KEY. Without it the value is stored as it is, and nothing fails: grit doctor makes this an error, and it is the first thing it found in a real project.',
    history: 'A plaintext leak was fixed in v3.212.0. Re-audit before storing real PII.',
  },
  {
    id: 'multitenancy',
    area: 'Multitenancy (plugin, tenant.Owned)',
    status: 'beta',
    guaranteed:
      'A tenant-owned model is scoped by the organization on the request context on every authenticated route group, a query with no organization is refused rather than answered (400 NO_ORGANIZATION), the X-Organization-ID header is checked against membership, and a role held through a membership grants inside that organization only.',
    proof:
      'A CI job scaffolds a project with the multitenant and impersonate plugins and runs 26 checks across the seams on a real Postgres: cross-tenant reads, an admin with no organization, an organization and an owner on one row, a membership role used in the wrong organization, and an impersonated session inside a tenant. It also asserts the middleware order that one of those bugs was.',
    yours:
      'Add tenant.Owned to every model that belongs to an organization. Decide which organization a user provisioned by SSO joins: nothing does it for you, and grit doctor says so. Combining tenancy with --tree or --public is still not covered by the fixture.',
    history:
      'Isolation bugs were fixed in v3.197 and v3.198. v3.233.0 fixed three found by building tenancy, roles and impersonation together: every DELETE on a tenant-owned row answered 500 because the staff group never resolved an organization, a membership role granted nothing at all, and “no active organization” was an opaque 500.',
  },
  {
    id: 'money',
    area: 'Money (exact, multi-currency)',
    status: 'stable',
    guaranteed: 'Amounts are integer minor units with their currency, two money columns per row do not collide, and sums stay exact.',
    proof: 'money/money_test.go (13 tests), generator tests for the embedded columns and whitelists, and a live round-trip check',
    yours: 'Pick the currency per row deliberately; Grit does not convert between currencies.',
  },
  {
    id: 'importexport',
    area: 'CSV/XLSX import and export',
    status: 'beta',
    guaranteed:
      'An export streams every matching row; an import runs in the background with a job to poll, resolves relations, skips duplicates, and an owned resource’s rows belong to whoever imported them unless an ADMIN names another owner.',
    proof: 'Live checks for both, including the owned cases, and generator tests for the importer’s columns',
    yours: 'Validate the files your users upload. The parser is hardened and bounded, and it is still parsing a stranger’s spreadsheet.',
    history: 'Every export was an empty 200 once. The importer created a user for every unrecognised email until that was removed.',
  },
  {
    id: 'trees',
    area: 'Trees (--tree)',
    status: 'beta',
    guaranteed:
      'A move carries its subtree, a reorder keeps the parent the node has, and a node cannot be moved under its own descendant.',
    proof: 'services/category_tree_test.go ships with the resource (10 tests), plus live checks for the move, the reorder and the refusal',
    yours: 'Run the rebuild endpoint once after adding --tree to a table that already had rows.',
    history: 'Dragging a node that predated --tree once added one to the depth of every row in the table.',
  },
  {
    id: 'publicapi',
    area: 'Public API (--public)',
    status: 'beta',
    guaranteed:
      'Public reads are an allowlist struct rather than the model, they are guarded by an API key, archived rows are excluded, and cost-shaped columns are held back by default.',
    proof: 'Live checks for the key, the allowlist, the filters and the archived case; grit doctor reports a held-back column that was published by hand',
    yours:
      'Review the allowlist when you add fields; it is your file and the generator will not rewrite it. Public responses are cached by URL for cache.public_ttl_seconds (60 by default), so an archived row can linger in a cached list for up to a minute.',
  },
  {
    id: 'appendonly',
    area: 'Append-only resources',
    status: 'beta',
    guaranteed: 'Rows are created and read, never changed or deleted: a GORM guard, no mutating routes, and a database trigger that refuses the write.',
    proof: 'appendonly/appendonly_test.go, and grit doctor reports an append-only resource that still mounts a write',
    yours: 'Run grit migrate so the trigger exists. The guard alone is not the database refusing.',
  },
  {
    id: 'realtime',
    area: 'Realtime (WebSockets, Redis backplane)',
    status: 'beta',
    guaranteed: 'Events reach a user on every replica through the backplane, and a closed socket is cleaned up.',
    proof: 'realtime/backplane_test.go, cluster/cluster_test.go',
    yours: 'Run Redis if you run more than one replica. Without it, events stay in the process that emitted them.',
    history: 'Cross-replica and revocation bugs were fixed in v3.193 and v3.196.',
  },
  {
    id: 'outbox',
    area: 'Durable events (outbox)',
    status: 'beta',
    guaranteed: 'A durable subscriber’s work is committed with the row that triggered it, and the relay delivers it once.',
    proof: 'outbox/outbox_test.go (12 tests)',
    yours: 'Start the relay. grit upgrade wires it for projects that predate it.',
  },
  {
    id: 'errors',
    area: 'Error responses',
    status: 'new',
    guaranteed:
      'Every error is the same envelope, each code is documented with what it means and what a client should do, and a code always arrives with the same status. The typed codes in Go and the union type in TypeScript are generated from one catalogue, so they cannot disagree.',
    proof: 'internal/errorcodes walks every template: a code that is not catalogued, or that is returned with a status other than its own, fails the build',
    yours:
      'Codes your own handlers invent are yours to document. respond.Rule covers a business rule without adding one.',
    history:
      'Shipped in v3.232.0, which reconciled 31 sites. VALIDATION_ERROR was 422 in thirty-eight handlers and 400 in twenty-five, INVALID_TOKEN was both 401 and 400, and an upload that did not exist answered 400.',
  },
  {
    id: 'migrations',
    area: 'Migrations and rollback',
    status: 'new',
    guaranteed:
      'Every run records what it changed, taken from the schema before and after, and grit migrate down drops exactly that: indexes, then columns, then tables. A drop already done is skipped, so an interrupted rollback can be re-run. The run that built the schema is a baseline and is refused.',
    proof: 'internal/migrate tests ship in every project, and the live suite adds a column with grit generate field, rolls it back on a real Postgres and checks the column is gone and the table is not',
    yours:
      'Read the statements before you answer the prompt. A rollback drops columns, and the data in them is not recoverable: this undoes a schema change, not a deployment.',
    history:
      'Shipped in v3.230.0. The first version counted SQLite’s own sqlite_sequence table as part of the schema, so on SQLite no run would have been recognised as the baseline and a rollback would have offered to drop everything.',
  },
  {
    id: 'backups',
    area: 'Backups and restore',
    status: 'beta',
    guaranteed: 'A backup streams row by row, and a restore writes parents before children.',
    proof: 'The CLI’s backup tests and the generated project’s backup command',
    yours: 'Restore into a scratch database and look at the rows. A backup nobody has restored is a hope.',
    history: 'Restore correctness bugs were fixed in v3.207.0 and v3.213.0.',
  },
  {
    id: 'offline',
    area: 'Offline sync (desktop)',
    status: 'new',
    guaranteed: 'Registered models mirror to the client, and conflicts resolve by the policy you set.',
    proof: 'The sync engine’s tests in the CLI, and grit sync doctor reports a model that cannot work (no version column, an allowlist naming a column that does not exist)',
    yours: 'Run grit sync doctor, and test the flows your users will have offline.',
  },
  {
    id: 'workflows',
    area: 'Workflows (status state machines)',
    status: 'new',
    guaranteed: 'Only declared transitions are allowed, they are guarded by role, and a transition hook commits with the move or not at all.',
    proof: 'The generated workflow service’s tests, plus generator tests for the transition route',
    yours: 'Decide who may make each transition. The generator gives you the shape, not your policy.',
    history: 'Shipped in v3.221.0: before it, a workflow could change a status and nothing else.',
  },
  {
    id: 'flags',
    area: 'Feature flags',
    status: 'new',
    guaranteed: 'A flag can target attributes, a percentage and a date window, and flags.IsEnabled works from anywhere, answering false before the engine starts.',
    proof: 'flags/flags_test.go ships in every project and runs against a real engine',
    yours: 'Manage flags through the API: there is no admin screen for them yet.',
    history: 'Shipped in v3.223.0, when the documented package-level API did not exist at all.',
  },
  {
    id: 'desktop-mobile',
    area: 'Desktop (Wails) and mobile (Expo)',
    status: 'new',
    guaranteed: 'Both share the monorepo API and its generated types.',
    proof: 'The desktop module builds in CI; the Expo app is built by the nightly canary',
    yours:
      'Signing and store submission. A signed desktop release has never been verified end to end, because that needs Authenticode and Apple Developer certificates this project does not own.',
  },
]

const BY_ID: Record<string, StabilityArea> = Object.fromEntries(
  STABILITY_AREAS.map((a) => [a.id, a]),
)

/** The area a capability row points at, or undefined when it points nowhere. */
export function stabilityOf(id: string | undefined): StabilityArea | undefined {
  return id ? BY_ID[id] : undefined
}

/** How many areas sit at each status. Used in prose, so it cannot go stale. */
export function stabilityCounts(): Record<Status, number> {
  const counts: Record<Status, number> = { stable: 0, beta: 0, new: 0 }
  for (const area of STABILITY_AREAS) counts[area.status] += 1
  return counts
}
