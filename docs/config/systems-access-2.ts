import type { SystemDesign } from './systems-types'

export const MULTITENANCY: SystemDesign = {
  slug: 'multitenancy',
  name: 'Multitenancy',
  tagline:
    'One deployment serving many organisations, where a missing predicate is a data breach rather than a bug.',
  group: 'Access',
  packages: ['plugin: multitenant', 'internal/authz', 'internal/middleware'],

  problem: {
    text: [
      'Selling the same application to many organisations has three possible shapes. A database per tenant is the safest and the most expensive to operate: a hundred customers means a hundred migrations. A schema per tenant is cheaper but still multiplies the operational surface. A shared schema with a tenant column is the cheapest and the most dangerous, because the entire isolation guarantee is one predicate that a developer can forget.',
      'Grit takes the third shape and spends its effort on making the predicate impossible to forget rather than on hoping. That is the only defensible version of the cheap option.',
      'The second problem is the one nobody plans for: an operator needs to see a customer’s data to answer a support ticket. Doing that by turning off the predicate is how a support tool becomes the breach. Impersonation has to be a thing the system knows it is doing.',
    ],
    capabilities: [
      'Resolve which organisation a request belongs to, from the identity rather than from the request.',
      'Narrow every query on a tenant-owned table, in the layer that builds queries.',
      'Stamp the tenant on creation and refuse to let it be changed afterwards.',
      'Let a user belong to more than one organisation and switch between them explicitly.',
      'Let an operator act inside a tenant deliberately, with that fact recorded.',
      'Combine with per-user ownership, so "my invoices inside my organisation" is expressible.',
    ],
  },

  functional: [
    'An embedded tenant.Owned on models that belong to an organisation, which is what marks them.',
    'Tenant resolution in middleware, from the membership on the identity.',
    'Query scoping applied in the service, alongside and independent of owner scoping.',
    'Membership records joining users to organisations, with a role per membership.',
    'An explicit organisation switch, which changes the tenant on subsequent requests.',
    'Operator impersonation of a tenant, time-boxed and audited.',
  ],

  nonFunctional: [
    {
      label: 'Never from the request',
      text: 'the tenant is derived from the authenticated membership. A header or a body field naming an organisation is ignored, because otherwise it is an access control bypass with a nice name.',
    },
    {
      label: 'Fails closed',
      text: 'a request whose tenant cannot be resolved matches nothing. There is no state in which an unresolved tenant means "all of them".',
    },
    {
      label: 'Immutable after creation',
      text: 'the tenant column is not writable through update or patch. Moving a row between organisations is an explicit operation, not a field edit.',
    },
    {
      label: 'Opt-in per model',
      text: 'a shared reference table scoped by accident becomes invisible to every query, with no error to explain it. Marking is deliberate.',
    },
    {
      label: 'Impersonation is visible',
      text: 'an operator inside a tenant is a recorded fact with a beginning and an end, not an invisible capability.',
    },
  ],

  capacity: {
    assumptions: [
      ['Organisations', '5,000'],
      ['Users per organisation', 'median 8, p99 400'],
      ['Rows per tenant-owned table', 'highly skewed; the largest tenant is often 100x the median'],
      ['Requests per second', '926 average'],
    ],
    estimates: [
      {
        label: 'Skew',
        working: [
          '5,000 tenants, 10,000,000 rows',
          'median tenant: ~500 rows',
          'largest tenant: ~1,000,000 rows, 10% of the table',
        ],
        note: 'This skew is the defining capacity property of shared-schema multitenancy. An index that works for the median tenant can still be a sequential scan for the largest one, and the largest one is usually the most important customer.',
      },
      {
        label: 'Index shape',
        working: [
          'an index on (tenant_id, created_at) serves list queries',
          'the median tenant reads 500 rows, the largest reads a page of 1,000,000',
          'pagination must be keyset, not offset, or the largest tenant pays for every page',
        ],
      },
      {
        label: 'Membership lookups',
        working: [
          '5,000 orgs x 8 median members = ~40,000 memberships',
          'loaded once per request with the identity, cached per user',
        ],
      },
    ],
  },

  highLevel: {
    intro:
      'Tenant scoping sits beside owner scoping rather than replacing it. A row can be narrowed by organisation, by user, or by both.',
    components: [
      {
        label: 'tenant.Owned',
        text: 'an embedded struct carrying the organisation id. A model that embeds it is tenant-scoped; one that does not is shared.',
      },
      {
        label: 'Tenant resolver',
        text: 'middleware that puts the current organisation on the context, from the membership on the identity and the user’s active selection.',
      },
      {
        label: 'Tenant scope',
        text: 'the query predicate, applied in the service. Composes with owner scoping rather than competing with it.',
      },
      {
        label: 'Membership store',
        text: 'who belongs to which organisation and in what role. One user may have several.',
      },
      {
        label: 'Impersonation',
        text: 'an operator-initiated, time-boxed entry into a tenant, which appears in the audit log of that tenant.',
      },
    ],
    flow: {
      title: 'A request inside an organisation',
      nodes: [
        { id: 'client', label: 'Client', col: 0, row: 0 },
        { id: 'auth', label: 'Auth Middleware', col: 1, row: 0 },
        { id: 'tenant', label: 'Tenant Resolver', col: 2, row: 0, accent: true },
        { id: 'members', label: 'Membership Store', col: 3, row: 0 },
        { id: 'service', label: 'Service', col: 2, row: 1, accent: true },
        { id: 'scope', label: 'Tenant + Owner Scope', col: 3, row: 1, accent: true },
        { id: 'db', label: 'Database', col: 2, row: 2 },
        { id: 'audit', label: 'Audit Log', col: 3, row: 2 },
      ],
      edges: [
        { from: 'client', to: 'auth', step: 1 },
        { from: 'auth', to: 'tenant', step: 2 },
        { from: 'tenant', to: 'members', step: 3 },
        { from: 'tenant', to: 'service', step: 4 },
        { from: 'service', to: 'scope', step: 5 },
        { from: 'service', to: 'db', step: 6 },
        { from: 'db', to: 'audit', step: 7, dashed: true },
      ],
      steps: [
        'The request arrives with an access token. Nothing in it names an organisation, and if it did it would be ignored.',
        'Authentication produces the identity. The tenant is not yet known.',
        'The resolver reads the user’s memberships and their active selection, and puts the organisation on the context. A user with no membership resolves to nothing.',
        'The handler calls the service exactly as it would in a single-tenant project. It does not mention the organisation.',
        'The service applies both scopes. The tenant predicate comes from the context; the owner predicate, if the resource is also user-owned, comes from the same place.',
        'The narrowed query runs. An unresolved tenant produces a predicate that matches nothing rather than one that is omitted.',
        'Writes are recorded in the tenant’s audit log, including the fact that an operator was impersonating, when that is the case.',
      ],
    },
    dataFlow: [
      'The organisation id is written at creation from the context and is not in the writable column set afterwards.',
      'Shared reference data, such as a currency table, is deliberately not tenant-scoped. Marking it would make it invisible to everybody.',
      'An impersonating operator carries both identities: their own, for the audit record, and the tenant, for the scope.',
      'Switching organisation is an explicit call that changes the active membership. It does not change any row.',
    ],
  },

  stack: [
    ['Isolation model', 'shared schema with a tenant column'],
    ['Marking', 'an embedded struct on the model'],
    ['Resolution', 'from the authenticated membership, never the request'],
    ['Enforcement', 'a query predicate in the service layer'],
    ['Indexes', 'composite, leading with the tenant column'],
    ['Availability', 'a plugin, because a single-tenant project should not carry the machinery'],
  ],

  dataModel: {
    entities: [
      {
        name: 'organizations',
        fields: [
          ['id', 'UUIDv7'],
          ['name, slug', 'what it is called'],
          ['created_at', 'when it was onboarded'],
        ],
      },
      {
        name: 'memberships',
        note: 'A user may belong to several. The active one is a per-user selection, not a property of the membership.',
        fields: [
          ['user_id', 'the person'],
          ['organization_id', 'the organisation'],
          ['role', 'their role inside this organisation, which may differ per membership'],
        ],
      },
      {
        name: 'any tenant-owned resource',
        fields: [
          ['organization_id', 'the tenant, indexed first in every composite index'],
          ['user_id', 'optionally also user-owned, for the both case'],
        ],
      },
    ],
    storage: [
      'One database, one schema. The cost of that choice is paid in discipline about the predicate, and the benefit is one migration rather than five thousand.',
      'Every index on a tenant-owned table leads with the tenant column. An index that does not is useless to a scoped query.',
      'Keyset pagination rather than offset, because the largest tenant would otherwise pay a growing cost for each page.',
    ],
  },

  api: {
    groups: [
      {
        title: 'Organisations',
        rows: [
          { method: 'GET', path: '/api/v1/organizations', what: 'The caller’s memberships' },
          { method: 'POST', path: '/api/v1/organizations/switch', what: 'Change the active organisation' },
          { method: 'GET', path: '/api/v1/organizations/members', what: 'Who is in the current one' },
          { method: 'POST', path: '/api/v1/organizations/invite', what: 'Invite somebody into it' },
          { method: 'DELETE', path: '/api/v1/organizations/members/:id', what: 'Remove a membership' },
        ],
      },
      {
        title: 'Impersonation',
        rows: [
          { method: 'POST', path: '/api/v1/admin/impersonate', what: 'Enter a tenant as an operator, time-boxed' },
          { method: 'POST', path: '/api/v1/admin/impersonate/stop', what: 'Leave it' },
        ],
      },
    ],
  },

  lowLevel: {
    classes: [
      {
        name: 'tenant.Owned',
        what: 'The embedded marker. Its presence on a model is what the generator and the scope both key on, so there is no separate list to keep in step.',
      },
      {
        name: 'TenantResolver',
        what: 'Middleware. Puts the organisation on the context from the membership, and refuses the request when a user has none.',
      },
      {
        name: 'ScopeTenant',
        what: 'The predicate, with the same three-case shape as owner scoping: exempt, impossible, narrowed. The impossible case is first in the code for the same reason.',
      },
      {
        name: 'Impersonation',
        what: 'A bounded capability: an operator, a tenant, an expiry, and an audit record written at both ends.',
      },
    ],
    principles: [
      {
        label: 'Opt in, never infer',
        text: 'a model is tenant-scoped because somebody said so. Inferring it from a column name would silently scope a shared table and make it vanish.',
      },
      {
        label: 'Compose, do not replace',
        text: 'tenant scope and owner scope are separate predicates that both apply. Collapsing them into one would make "my rows in my organisation" inexpressible.',
      },
      {
        label: 'The dangerous path is explicit',
        text: 'crossing a tenant boundary is impersonation, which has a name, a duration and a log entry, rather than a flag that disables a predicate.',
      },
    ],
    patterns: [
      ['Discriminator column', 'shared schema, one column deciding visibility'],
      ['Ambient context', 'the tenant travels with the request rather than through every signature'],
      ['Scoped repository', 'the service is the only thing that builds a query'],
      ['Break-glass with an audit trail', 'impersonation as a recorded, bounded act'],
    ],
  },

  scaling: [
    'The tenant predicate improves selectivity, so a correctly indexed multitenant query is faster than the unscoped equivalent.',
    'Tenant skew is the thing to watch. The largest tenant can be a tenth of the table, and a plan that suits the median can be a scan for them.',
    'Keyset pagination keeps the cost per page constant regardless of tenant size, which offset pagination does not.',
    'A single very large tenant can be moved to its own database later without changing application code, because the predicate is already there.',
    'Memberships are small and cached per user. Resolving a tenant does not query on the hot path.',
  ],

  bottlenecks: {
    problems: [
      {
        label: 'The forgotten predicate',
        text: 'a hand-written query that does not go through the service sees every tenant. This is the entire risk of the shared-schema model, concentrated in one failure.',
      },
      {
        label: 'Tenant skew',
        text: 'one customer with a tenth of the rows experiences a different application from everybody else, and they are usually the customer who notices.',
      },
      {
        label: 'Noisy neighbours',
        text: 'a heavy report run by one tenant consumes connections and cache that every other tenant shares.',
      },
      {
        label: 'Impersonation drift',
        text: 'a support tool that can enter any tenant without expiry or record becomes the softest way into every customer’s data.',
      },
      {
        label: 'Untested combinations',
        text: 'tenant scoping combined with tree structures or public endpoints is where the predicate is most likely to be dropped, and least likely to be covered.',
      },
    ],
    improvements: [
      {
        label: 'Test the breach, not the feature',
        text: 'a fixture with two tenants, asserting that each cannot see the other, verified to fail when the predicate is removed.',
      },
      {
        label: 'Index leading with the tenant',
        text: 'every composite index starts with the tenant column, so the largest tenant gets the same plan as the smallest.',
      },
      {
        label: 'Per-tenant limits',
        text: 'rate limiting keyed by organisation as well as by address, so one tenant cannot consume the shared capacity.',
      },
      {
        label: 'Time-box and record impersonation',
        text: 'a session with an expiry and an entry in the tenant’s own audit log, so the customer can see it happened.',
      },
    ],
  },

  seeAlso: [
    { title: 'Data isolation', href: '/docs/systems/data-isolation' },
    { title: 'Authorization', href: '/docs/systems/authorization' },
  ],
}

export const RATE_LIMITING: SystemDesign = {
  slug: 'rate-limiting',
  name: 'Rate Limiting and Abuse Control',
  tagline:
    'Keeping one caller from consuming the capacity of all of them, and making a password-guessing attack cost something.',
  group: 'Access',
  packages: ['Sentinel', 'internal/middleware'],

  problem: {
    text: [
      'An API with no limit serves whoever asks loudest. A script with a loop will take as much capacity as the hardware can give it, and the users it starves are indistinguishable from the attacker until somebody reads a log.',
      'There are three different problems wearing the same name. Protecting capacity means limiting how often anybody can call anything. Protecting an account means limiting how often somebody can guess its password, which has to be counted per account rather than per address or a botnet defeats it. Protecting an expensive endpoint means limiting the specific operations that cost real money, like sending an email or generating a PDF.',
      'They need different keys, different windows and different responses, which is why one global limit is never enough.',
    ],
    capabilities: [
      'Limit requests per address, with a window short enough to matter and long enough not to punish a burst.',
      'Limit sign-in attempts per account, independently of where they come from.',
      'Limit specific expensive operations more tightly than ordinary reads.',
      'Tell a rejected caller when to come back, rather than simply refusing.',
      'Keep working when the shared counter store is unavailable, without either opening up or shutting down.',
      'Make the limits visible, so an operator can see who is being limited and why.',
    ],
  },

  functional: [
    'Per-address limiting on the whole API, configurable per route group.',
    'Account lockout after a number of failed sign-ins within a window, counted per account.',
    'Tighter limits on sign-in, password reset, invitation and other endpoints that send mail.',
    'Standard rate limit headers on every response, and Retry-After on a refusal.',
    'A dashboard showing current limits, recent refusals and the addresses involved.',
    'Limits that survive a restart, because an attacker can simply wait for a deploy otherwise.',
  ],

  nonFunctional: [
    {
      label: 'Cheap on the happy path',
      text: 'the overwhelming majority of requests are not limited. The check must cost microseconds, not a database round trip.',
    },
    {
      label: 'Counted before the expensive work',
      text: 'the limiter runs before password hashing and before the handler, or the attack still costs the server what it was meant to save.',
    },
    {
      label: 'Honest refusals',
      text: 'a 429 with Retry-After lets a well-behaved client back off correctly instead of retrying immediately and making it worse.',
    },
    {
      label: 'Degrades sensibly',
      text: 'if the shared counter is unreachable, the limiter falls back to a local one rather than failing open to unlimited or closed to nothing.',
    },
    {
      label: 'Not a membership oracle',
      text: 'account lockout applies to real accounts. An unknown address must not behave differently, or the endpoint reveals who is registered.',
    },
  ],

  capacity: {
    assumptions: [
      ['Requests per second', '926 average, 4,630 peak'],
      ['Distinct client addresses per minute', '~30,000'],
      ['Default limit', '100 requests per minute per address'],
      ['Sign-in limit', '10 failures per account per 15 minutes'],
      ['Counter entry', '~64 bytes'],
    ],
    estimates: [
      {
        label: 'Counter memory',
        working: [
          '30,000 active addresses x 64 bytes = ~2 MB',
          'plus sign-in counters for accounts under attack, small',
        ],
        note: 'The state is tiny. Rate limiting is not a storage problem; it is a coordination problem between replicas.',
      },
      {
        label: 'Check cost',
        working: [
          'in-process counter: one map lookup and an increment, ~50 nanoseconds',
          'shared counter: one round trip, ~0.3 ms',
        ],
        note: 'The difference is the whole design question. A per-process counter is free and wrong by a factor of the replica count; a shared one is correct and adds a hop to every request.',
      },
      {
        label: 'Effective limit across replicas',
        working: [
          '100/minute per address, enforced per process',
          'x 6 replicas = up to 600/minute actually allowed',
        ],
        note: 'This is why a shared store matters once there is more than one replica: otherwise the configured number is not the enforced number.',
      },
    ],
  },

  highLevel: {
    intro:
      'Three limiters with different keys, all running before the handler. The first protects capacity, the second protects accounts, the third protects specific costs.',
    components: [
      {
        label: 'Request limiter',
        text: 'per address, applied across the API. The broad one, sized so that normal use never meets it.',
      },
      {
        label: 'Account lockout',
        text: 'per account, counting failed sign-ins. Independent of address, which is what makes it work against a distributed attack.',
      },
      {
        label: 'Endpoint limiter',
        text: 'tighter rules on the routes that cost money or send mail, keyed by whichever identifier fits that route.',
      },
      {
        label: 'Counter store',
        text: 'shared when available so the configured limit is the enforced limit, with a local fallback when it is not.',
      },
      {
        label: 'Dashboard',
        text: 'what is being limited, how often, and from where. A limiter nobody can observe is a limiter nobody can tune.',
      },
    ],
    flow: {
      title: 'Where the limiters sit relative to the work',
      nodes: [
        { id: 'client', label: 'Client', col: 0, row: 0 },
        { id: 'limiter', label: 'Request Limiter', col: 1, row: 0, accent: true },
        { id: 'store', label: 'Counter Store', col: 2, row: 0 },
        { id: 'local', label: 'Local Fallback', col: 3, row: 0 },
        { id: 'reject', label: '429 + Retry-After', col: 1, row: 1 },
        { id: 'lockout', label: 'Account Lockout', col: 2, row: 1, accent: true },
        { id: 'handler', label: 'Handler', col: 3, row: 1 },
        { id: 'password', label: 'Password Hash', col: 3, row: 2 },
        { id: 'dash', label: 'Dashboard', col: 2, row: 2 },
      ],
      edges: [
        { from: 'client', to: 'limiter', step: 1 },
        { from: 'limiter', to: 'store', step: 2 },
        { from: 'store', to: 'local', step: 3, dashed: true },
        { from: 'limiter', to: 'reject', step: 4, dashed: true },
        { from: 'limiter', to: 'lockout', step: 5 },
        { from: 'lockout', to: 'handler', step: 6 },
        { from: 'handler', to: 'password', step: 7 },
        { from: 'store', to: 'dash', step: 8, dashed: true },
      ],
      steps: [
        'The request arrives. Nothing has been parsed beyond the route and the client address.',
        'The request limiter increments the counter for that address and reads the current count.',
        'If the shared store is unreachable the limiter uses a per-process counter instead. The limit becomes approximate rather than absent, which is the right trade in an outage.',
        'Over the limit, the request ends here with 429 and a Retry-After naming when the window resets. No handler ran and nothing was parsed.',
        'On a sign-in route the account lockout check runs next, keyed by the submitted address. An account already barred is refused before anything expensive happens.',
        'Only now does the handler run.',
        'Password hashing happens last, which is the point of the ordering: an attack that is going to be rejected is rejected before it costs sixty milliseconds of CPU.',
        'Counters feed a dashboard, so the limits can be tuned against what is actually happening rather than guessed at.',
      ],
    },
    dataFlow: [
      'Counters are keyed by address for the broad limit and by account for lockout. Mixing the two produces a limiter that a botnet defeats and that punishes an office behind one address.',
      'An unknown email address never increments an account counter, or anybody could lock out an address they can guess.',
      'Rate limit headers are set on successful responses too, so a client can slow down before it is refused.',
      'The counter store is the only shared state. Everything else is per-request.',
    ],
  },

  stack: [
    ['Request limiting', 'Sentinel, a sliding window per key'],
    ['Counter store', 'Redis when configured, in-process otherwise'],
    ['Account lockout', 'database-backed, so it survives a restart'],
    ['Headers', 'X-RateLimit-Limit, X-RateLimit-Remaining, Retry-After'],
    ['Status', '429 Too Many Requests'],
    ['Observability', 'the Sentinel dashboard, mounted in the API'],
  ],

  dataModel: {
    entities: [
      {
        name: 'login_attempts',
        note: 'Durable on purpose. An in-memory counter is reset by a deploy, which is a window an attacker can simply wait for.',
        fields: [
          ['email', 'the account being guessed at, indexed'],
          ['failed_count', 'attempts inside the current window'],
          ['locked_until', 'when it becomes usable again'],
          ['last_attempt_at', 'drives the window'],
        ],
      },
    ],
    storage: [
      'Request counters are ephemeral and belong in a store with expiry. Losing them on a restart costs one window of accuracy.',
      'Account lockout is durable, because the whole point is that it cannot be cleared by waiting for a deploy.',
    ],
  },

  api: {
    groups: [
      {
        title: 'What a limited caller sees',
        rows: [
          { method: 'GET', path: 'any route', what: 'Rate limit headers on every response, not only refusals' },
          { method: 'POST', path: '/api/v1/auth/login', what: '429 when limited, 423 when the account is locked' },
        ],
      },
      {
        title: 'Operations',
        rows: [
          { method: 'GET', path: '/sentinel/ui', what: 'The dashboard: current limits and recent refusals' },
          { method: 'GET', path: '/api/v1/admin/security/lockouts', what: 'Accounts currently barred' },
          { method: 'DELETE', path: '/api/v1/admin/security/lockouts/:email', what: 'Clear one, for support' },
        ],
      },
    ],
    samples: [
      {
        title: 'A refusal that tells the client what to do',
        language: 'http',
        code: `HTTP/1.1 429 Too Many Requests
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 0
Retry-After: 37

{
  "error": {
    "code": "RATE_LIMITED",
    "message": "Too many requests. Try again in 37 seconds."
  }
}`,
      },
    ],
  },

  lowLevel: {
    classes: [
      {
        name: 'Sentinel middleware',
        what: 'The request limiter. Runs before everything, keyed by address, with per-route-group configuration.',
      },
      {
        name: 'LoginAttemptService',
        what: 'The account counter. Increments only on a failed attempt against a real account, and reads before any hashing is paid for.',
        methods: ['IsLocked', 'RecordFailure', 'Clear'],
      },
      {
        name: 'CounterStore',
        what: 'An interface with two implementations, shared and local. The fallback is a decision the limiter makes, not an outage the request sees.',
      },
    ],
    principles: [
      {
        label: 'Reject before you pay',
        text: 'the ordering of the middleware chain is the security property. A limiter after the expensive work limits nothing that matters.',
      },
      {
        label: 'Key by the thing you are protecting',
        text: 'capacity is per address, accounts are per account, cost is per operation. One key cannot serve all three.',
      },
      {
        label: 'Degrade, do not fail',
        text: 'an unreachable counter store makes the limit approximate. Failing open removes protection; failing closed turns a cache outage into an outage.',
      },
      {
        label: 'Make it observable',
        text: 'a limit nobody can see is tuned by guesswork and discovered when a customer complains.',
      },
    ],
    patterns: [
      ['Sliding window', 'smoother than a fixed window at the boundary'],
      ['Token bucket', 'for endpoints that should tolerate a burst and then slow down'],
      ['Circuit fallback', 'local counting when the shared store is unavailable'],
      ['Middleware ordering', 'cheap checks first, expensive work last'],
    ],
  },

  scaling: [
    'A shared counter store makes the configured limit the enforced limit across replicas. Without it, the real limit is the configured one multiplied by the replica count.',
    'The counter check adds one round trip to the shared store. On a local network that is a fraction of a millisecond, which is the price of the limit being true.',
    'Counters expire with their window, so the store does not grow with traffic over time.',
    'Per-tenant keys matter in a multitenant deployment, or one customer can consume the capacity that everybody shares.',
    'Limits are configured per route group, so an expensive export can be a hundred times tighter than a list without changing the global number.',
  ],

  bottlenecks: {
    problems: [
      {
        label: 'Shared addresses',
        text: 'an office or a mobile carrier presents one address for thousands of people, who then share one limit and all get refused together.',
      },
      {
        label: 'The counter store as a dependency',
        text: 'putting a network call in front of every request makes the limiter a new thing that can be down.',
      },
      {
        label: 'Lockout as a denial of service',
        text: 'if an attacker can lock any account by guessing wrong ten times, lockout is a weapon pointed at your users.',
      },
      {
        label: 'Window boundary bursts',
        text: 'a fixed window lets a caller spend the whole allowance at the end of one and the start of the next, for double the intended rate.',
      },
    ],
    improvements: [
      {
        label: 'Key authenticated traffic by user',
        text: 'once a caller is identified, limiting by user rather than by address removes the shared-address problem entirely.',
      },
      {
        label: 'Local fallback with a short leash',
        text: 'the limiter keeps working on a per-process counter when the store is away, and says so in the dashboard rather than silently.',
      },
      {
        label: 'Delay rather than lock',
        text: 'an increasing delay after each failure costs an attacker far more than it costs a user who mistyped, without handing anybody a lockout weapon.',
      },
      {
        label: 'Sliding windows',
        text: 'counting over a moving interval removes the boundary burst without the bookkeeping of a full leaky bucket.',
      },
    ],
  },

  seeAlso: [
    { title: 'Authentication', href: '/docs/systems/authentication' },
    { title: 'Observability', href: '/docs/systems/observability' },
  ],
}
