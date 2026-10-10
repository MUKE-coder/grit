import type { SystemDesign } from './systems-types'

/* Access: what a known caller may touch. */

export const AUTHORIZATION: SystemDesign = {
  slug: 'authorization',
  name: 'Authorization',
  tagline:
    'Deciding what a known caller may do, in a way that fails closed and can be listed route by route.',
  group: 'Access',
  packages: ['internal/authz', 'internal/access', 'internal/middleware'],

  interview: {
    intro:
      '"Design an access control system" is a standard prompt, and it is also the follow-up that ends most authentication questions. Interviewers are looking for more than an admin flag, and for an answer to "how do you know it is right".',
    questions: [
      {
        q: 'What is the difference between authentication and authorization?',
        a: 'Authentication is who you are. Authorization is what you may do. They are separate systems and merging them causes real breaches, because a token can be entirely valid and still not entitle its holder to the record they just asked for. In this design the auth middleware produces an identity and makes no decision about permission; a second layer reads that identity and decides.',
        see: 'problem',
      },
      {
        q: 'Admins can do everything and users can do some things. Is that an access control model?',
        a: 'No, and an interviewer asking this is checking whether you will say so. It collapses the moment a product needs an editor who can publish but not delete, or a finance role that sees invoices and nothing else. The usual models are role-based access control, where permissions attach to roles and roles to users, attribute-based, where the decision is computed from properties of the request, and resource-based, where the record itself names who may touch it. This design uses the first, with ownership handled separately.',
        see: 'problem',
      },
      {
        q: 'Where do you put the check?',
        a: 'In one place that every request passes through, which here is a middleware driven by a registry of route requirements. Checks scattered across handlers are the same question asked forty times, and the fortieth is the one somebody forgets. A registry also makes the whole policy a list you can read in one sitting, which is what turns a security review from an exercise in hunting into a query.',
        see: 'high-level-design',
      },
      {
        q: 'What is the default when a route has no rule?',
        a: 'Denied. This is the single most important line in the design. If an unlisted route is open, every new endpoint is public until somebody remembers to protect it, and the failure is silent. If an unlisted route is closed, the failure is a 403 in development that somebody fixes in a minute. Default deny makes forgetting loud rather than dangerous.',
        see: 'requirements',
      },
      {
        q: 'Should this also decide whether the user owns the record?',
        a: 'No, and keeping the two apart matters. Role permission answers "may someone like you edit invoices", and ownership answers "is this particular invoice yours". Only the second needs the record, so only the second can be decided after it is loaded. Merging them produces the classic bug where a user who legitimately holds an edit permission can edit everybody else records.',
        see: 'low-level-design',
      },
      {
        q: 'How does the user interface know which buttons to show?',
        a: 'From the same source as the enforcement, fetched as the effective permissions of the current user, never from a list hardcoded in the frontend. Two copies of the rules drift, and the drift is visible either as a button that errors when clicked or, worse, as a hidden button for something the backend would happily allow. The frontend is a convenience; the server decides.',
        see: 'api',
      },
      {
        q: 'What does the check cost on every request?',
        a: 'A map lookup, because the permission set is resolved once and carried on the identity rather than queried per check. That matters: an authorization layer that costs a database read per request will be worked around, and a layer people work around is not a layer. The capacity section gives the numbers.',
        see: 'capacity',
      },
      {
        q: 'A role changes. When do the people in it notice?',
        a: 'Not instantly, if permissions are cached or carried in a token, and that window is the thing to be explicit about. The options are to invalidate the cached set on write, to keep the carried copy short-lived, or to look up on every request and pay for it. Saying "it is eventually consistent and here is the window" is a much better answer than implying a change is instant when it is not.',
        see: 'bottlenecks',
      },
      {
        q: 'How do you stop roles multiplying until nobody can explain them?',
        a: 'By reviewing roles rather than users, and by resisting a role per customer request. The failure mode is well known: a role is added for one person, never removed, and after two years the permission model is a list nobody dares change. Keeping permissions fine-grained and roles few, and periodically reading the whole registry, is the practical defence.',
        see: 'bottlenecks',
      },
      {
        q: 'Your design has an administrator who bypasses the checks. Is that not the hole?',
        a: 'It is, and it is deliberate, so the thing to do is bound it rather than pretend it is not there. A bypass that is implicit and unlogged is how a support engineer reads customer data with no record. A bypass that is explicit, narrow and written to an audit log is a known risk with evidence attached. An interviewer raising this wants to hear the second answer.',
        see: 'bottlenecks',
      },
    ],
  },

  problem: {
    text: [
      'Authentication answers who is calling. Authorization answers whether they may. The two get conflated constantly, and the result is an application where the only check is "is there a token", which means any signed-in user can call any endpoint.',
      'The naive fix is a role check in each handler. It works until there are two hundred routes, at which point nobody can answer "who can delete an invoice" without reading two hundred functions, and a route added last week has no check at all because the person who added it did not know there was a convention.',
      'What is wanted is a permission model fine enough to be useful, a place that knows every route and the permission it needs, and a default that refuses rather than allows when a route is not listed. The last part is the one that matters: an authorization system whose failure mode is "allow" is not an authorization system.',
    ],
    capabilities: [
      'Express permissions finely enough that "may read invoices" and "may delete invoices" are different answers.',
      'Group permissions into roles, so people are assigned a job rather than forty checkboxes.',
      'Attach a required permission to every route, in one place that can be listed and reviewed.',
      'Refuse a route that nobody has classified, rather than allowing it.',
      'Let an administrator bypass the lot, because somebody has to be able to fix things.',
      'Answer the frontend’s question "what may this user see" without the frontend becoming the enforcement point.',
    ],
  },

  functional: [
    'Three built-in roles, ADMIN, EDITOR and USER, with permissions attached to each.',
    'A permission per resource and verb, generated with the resource: invoices.view, invoices.create, invoices.update, invoices.delete.',
    'Custom roles created at runtime from the admin, holding any set of permissions.',
    'A generated registry naming every route and the permission that guards it.',
    'Middleware that reads the identity from the request context and the requirement from the registry.',
    'An endpoint the frontend calls to learn its own permissions, used to hide what cannot be used.',
    'An administrator short-circuit, recorded in the audit log when it is exercised.',
  ],

  nonFunctional: [
    {
      label: 'Fails closed',
      text: 'a route with no entry in the registry is refused. The cost of forgetting is a 403 in testing, not an open endpoint in production.',
    },
    {
      label: 'Auditable as a list',
      text: 'the registry is a generated file. "Who can delete an invoice" is answered by reading one table, not by searching handlers.',
    },
    {
      label: 'Cheap',
      text: 'a permission check is a set membership test against data already loaded with the identity. It does not query.',
    },
    {
      label: 'Not the frontend’s job',
      text: 'the API enforces. The permission endpoint exists so the UI can hide a button, never so it can decide.',
    },
    {
      label: 'Separable from ownership',
      text: 'a permission says which kind of thing you may touch. Ownership says which rows. Conflating them produces a model that cannot express "every user may read their own invoices".',
    },
  ],

  capacity: {
    assumptions: [
      ['Routes in a generated project', '~200, growing with each resource'],
      ['Permissions', '4 per resource plus ~40 system permissions'],
      ['Roles', '3 built-in plus a handful of custom'],
      ['Permissions per role', '10 to 200'],
      ['Requests needing a check', 'all authenticated traffic, ~926/second average'],
    ],
    estimates: [
      {
        label: 'Check cost',
        working: [
          'permissions are loaded once with the identity, as a set',
          'a check is one hash lookup, ~20 nanoseconds',
          '926 checks/second x 20 ns = negligible',
        ],
        note: 'The cost is in loading the set, not testing it, which is why the set travels with the identity rather than being fetched per check.',
      },
      {
        label: 'Identity payload',
        working: [
          '200 permissions x ~24 bytes = ~5 KB per identity',
          'loaded once per request from the role, cached per role',
        ],
        note: 'Small enough to carry, large enough that it is worth caching per role rather than rebuilding for every request.',
      },
      {
        label: 'Registry size',
        working: ['~200 routes x one line = a generated file of a few hundred lines'],
        note: 'Regenerated whenever routes change, so it cannot drift from the routes it describes.',
      },
    ],
  },

  highLevel: {
    intro:
      'The identity arrives from authentication. The registry says what this route needs. The middleware compares the two. Nothing in a handler decides.',
    components: [
      {
        label: 'Actor',
        text: 'the identity on the request context: user id, role, whether they are an administrator, and the permission set. A value, not a database handle.',
      },
      {
        label: 'Access registry',
        text: 'a generated map from route to required permission, covering every route the router knows. Regenerated whenever routes change.',
      },
      {
        label: 'Authorization middleware',
        text: 'looks the route up, tests the actor’s set, and refuses with 403 if the permission is missing or the route is unlisted.',
      },
      {
        label: 'Role store',
        text: 'roles and their permissions, including custom roles made at runtime. Cached per role, invalidated when a role changes.',
      },
      {
        label: 'Permission endpoint',
        text: 'returns the caller’s own set so the interface can hide what will be refused. Advisory only.',
      },
    ],
    flow: {
      title: 'An authenticated request being authorized',
      nodes: [
        { id: 'client', label: 'Client', col: 0, row: 0 },
        { id: 'authmw', label: 'Auth Middleware', col: 1, row: 0 },
        { id: 'authzmw', label: 'Authz Middleware', col: 2, row: 0, accent: true },
        { id: 'registry', label: 'Access Registry', col: 3, row: 0 },
        { id: 'actor', label: 'Actor on Context', col: 1, row: 1 },
        { id: 'roles', label: 'Role Store', col: 2, row: 1, accent: true },
        { id: 'deny', label: '403 Forbidden', col: 3, row: 1 },
        { id: 'handler', label: 'Handler', col: 2, row: 2 },
        { id: 'service', label: 'Service', col: 3, row: 2 },
      ],
      edges: [
        { from: 'client', to: 'authmw', step: 1 },
        { from: 'authmw', to: 'actor', step: 2, bend: 'v' },
        { from: 'authmw', to: 'authzmw', step: 3 },
        { from: 'authzmw', to: 'registry', step: 4 },
        { from: 'authzmw', to: 'roles', step: 5 },
        { from: 'roles', to: 'deny', step: 6, dashed: true },
        { from: 'authzmw', to: 'handler', step: 7, bend: 'v' },
        { from: 'handler', to: 'service', step: 8 },
      ],
      steps: [
        'The request arrives with an access token. Authorization has not happened yet and nothing downstream has run.',
        'The auth middleware verifies the token and puts an actor on the context: user id, role, administrator flag and permission set.',
        'Control passes to the authorization middleware, which is the only component that decides.',
        'It looks up the matched route in the generated registry. A route with no entry is refused, because an unclassified route is a route nobody has thought about.',
        'The required permission is tested against the actor’s set. An administrator passes without the test, and that bypass is recorded.',
        'A missing permission ends here with 403. The handler never runs, so it cannot leak the existence of the thing by its error message.',
        'On success the handler runs, with the same actor still on the context.',
        'The service narrows the query by ownership where the resource is owned. Permission said which kind of thing; ownership says which rows.',
      ],
    },
    dataFlow: [
      'The permission set is resolved from the role once and travels with the identity. A handler asking twice costs nothing.',
      'The registry is generated from the route definitions, so it cannot describe a route that does not exist or miss one that does.',
      'Role changes invalidate the cached permission set for that role. A user whose role is revoked loses access on their next request, not at their next sign-in.',
      'The permission endpoint reads the same set the middleware tests, so the interface and the enforcement cannot disagree.',
    ],
  },

  stack: [
    ['Model', 'role-based, with permissions as the unit'],
    ['Permission naming', 'resource.verb, generated with the resource'],
    ['Registry', 'a generated Go file, rebuilt when routes change'],
    ['Enforcement point', 'middleware, before the handler'],
    ['Default', 'deny; an unlisted route is refused'],
    ['Role storage', 'the primary database, cached per role'],
  ],

  dataModel: {
    entities: [
      {
        name: 'roles',
        fields: [
          ['id', 'UUIDv7'],
          ['name', 'ADMIN, EDITOR, USER, or a custom one'],
          ['is_system', 'built-in roles cannot be deleted'],
        ],
      },
      {
        name: 'permissions',
        fields: [
          ['key', 'invoices.view, users.delete, system.view'],
          ['description', 'what it allows, shown in the admin'],
          ['feature', 'the group it belongs to, for the permission picker'],
        ],
      },
      {
        name: 'role_permissions',
        note: 'The join. Changing a row here changes what a role may do on the next request.',
        fields: [
          ['role_id', 'the role'],
          ['permission_key', 'the permission it grants'],
        ],
      },
    ],
    storage: [
      'Roles and permissions are small, slow-changing and read constantly, which is the textbook case for caching by role rather than by user.',
      'The registry is not in the database. It is generated code, so it is reviewed in a pull request alongside the routes it guards.',
    ],
  },

  api: {
    groups: [
      {
        title: 'Introspection',
        rows: [
          { method: 'GET', path: '/api/v1/auth/permissions', what: 'The caller’s own permission set' },
        ],
      },
      {
        title: 'Administration',
        rows: [
          { method: 'GET', path: '/api/v1/roles', what: 'Every role and its permissions' },
          { method: 'POST', path: '/api/v1/roles', what: 'Create a custom role' },
          { method: 'PUT', path: '/api/v1/roles/:id', what: 'Change which permissions it grants' },
          { method: 'DELETE', path: '/api/v1/roles/:id', what: 'Remove a custom role' },
          { method: 'GET', path: '/api/v1/permissions', what: 'The catalogue, grouped by feature' },
        ],
      },
    ],
    samples: [
      {
        title: 'What the frontend asks for',
        language: 'json',
        code: `{
  "data": {
    "role": "EDITOR",
    "is_admin": false,
    "permissions": [
      "invoices.view",
      "invoices.create",
      "invoices.update",
      "contacts.view",
      "contacts.create"
    ]
  }
}`,
      },
      {
        title: 'A refusal says nothing about what exists',
        language: 'json',
        code: `{
  "error": {
    "code": "FORBIDDEN",
    "message": "You do not have permission to access this resource"
  }
}`,
      },
    ],
  },

  lowLevel: {
    classes: [
      {
        name: 'Actor',
        file: 'internal/authz/actor.go',
        what: 'The identity as a value: user id, role, administrator flag, permission set. Taking a context rather than a request is what lets a background job carry the same identity.',
        methods: ['ActorFrom', 'UserIDFrom', 'Can'],
      },
      {
        name: 'access.Registry',
        file: 'internal/access/registry.go',
        what: 'Generated. One entry per route, naming the permission it needs. The file is the answer to "who can do what", in a form a reviewer can read.',
      },
      {
        name: 'RequirePermission',
        file: 'internal/middleware/authz.go',
        what: 'The enforcement point. The only place a 403 originates for a permission reason.',
      },
      {
        name: 'RoleService',
        what: 'Resolves a role to its permission set and caches it, invalidating on change.',
        methods: ['PermissionsFor', 'Invalidate'],
      },
    ],
    principles: [
      {
        label: 'Single point of enforcement',
        text: 'handlers do not check permissions. If they did, the registry would be a description rather than the truth, and the two would drift.',
      },
      {
        label: 'Deny by default',
        text: 'the registry is exhaustive and an unlisted route is refused. Forgetting produces a loud failure in development rather than a quiet hole in production.',
      },
      {
        label: 'Separate concerns that look alike',
        text: 'permissions are about kinds of things, ownership is about rows. Keeping them apart is what allows "any user may read their own invoices" to be expressible.',
      },
      {
        label: 'Generate what must not drift',
        text: 'the registry is derived from the routes rather than maintained beside them, so a new route cannot be missed.',
      },
    ],
    patterns: [
      ['Role-based access control', 'permissions grouped into roles, roles assigned to users'],
      ['Policy registry', 'one generated table from route to requirement'],
      ['Middleware chain', 'authenticate, then authorize, then handle'],
      ['Cache aside', 'permission sets by role, invalidated on write'],
    ],
  },

  scaling: [
    'Checks are in-memory set tests against data already loaded, so authorization adds no queries to the request path.',
    'Permission sets are cached per role, not per user, so a million users with three roles means three cached sets.',
    'A role change invalidates one cache entry and takes effect on the next request, without signing anybody out.',
    'The registry is compiled in. Looking a route up is a map read and does not touch storage at any traffic level.',
    'Growth is in routes and resources, not in traffic. A project with a thousand routes has a larger registry and the same per-request cost.',
  ],

  bottlenecks: {
    problems: [
      {
        label: 'Role explosion',
        text: 'custom roles multiply until there is one per person, at which point the model has become per-user permissions with extra indirection.',
      },
      {
        label: 'The administrator bypass',
        text: 'an account that passes every check is a single credential that opens everything, and it is usually the account with the weakest operational hygiene.',
      },
      {
        label: 'Stale permission sets',
        text: 'a cached set that is not invalidated leaves a revoked user with access until the entry expires.',
      },
      {
        label: 'Frontend drift',
        text: 'a UI that hides buttons by its own rules rather than by the permission endpoint will eventually hide something that is allowed, or show something that is not.',
      },
      {
        label: 'Permission without ownership',
        text: 'granting invoices.view to every user, on a resource that is not owner-scoped, means every user can read every invoice. The two systems have to be used together.',
      },
    ],
    improvements: [
      {
        label: 'Review roles, not users',
        text: 'the admin shows which users hold each role, so the question at review time is "should these twelve people be editors", which is answerable.',
      },
      {
        label: 'Audit the bypass',
        text: 'every administrator action is recorded with the fact that it bypassed a check, so the privileged path is visible rather than assumed.',
      },
      {
        label: 'Invalidate on write',
        text: 'the role store invalidates its own cache when a role changes, rather than relying on an expiry nobody tuned.',
      },
      {
        label: 'One source for the UI',
        text: 'the interface hides on the permission endpoint and nothing else, so a change to a role is reflected without a frontend deploy.',
      },
    ],
  },

  seeAlso: [
    { title: 'Data isolation', href: '/docs/systems/data-isolation' },
    { title: 'Roles and permissions reference', href: '/docs/backend/rbac' },
  ],
}

export const DATA_ISOLATION: SystemDesign = {
  slug: 'data-isolation',
  name: 'Data Isolation',
  tagline:
    'Making "your invoices" mean yours, in the query rather than in the handler, so a forgotten check is not a breach.',
  group: 'Access',
  packages: ['internal/authz', 'internal/services'],

  interview: {
    intro:
      'This is the "your API returned another user invoice" question, and it is asked because it is one of the most common real vulnerabilities. The interviewer wants the fix to be structural rather than a reminder to be careful.',
    questions: [
      {
        q: 'A signed-in user requests /invoices/4821, which belongs to somebody else. What happens?',
        a: 'They get a 404. The request is authenticated, so a check that only asks "is this caller signed in" passes and hands over the record. This is usually listed as the most common serious API vulnerability, under the name broken object level authorization, and it is common because the check that is missing is the one that needs the record in hand.',
        see: 'problem',
      },
      {
        q: 'Should that be a 403 or a 404?',
        a: 'A 404, and the reasoning is worth giving. A 403 confirms the record exists, which on any endpoint with guessable ids lets somebody map your data by walking numbers and reading the status codes. Reserve 403 for an action the caller may not perform on a record they can already see. For a record they may not see at all, the honest answer is that it does not exist as far as they are concerned.',
        see: 'requirements',
      },
      {
        q: 'Where do you put the ownership check?',
        a: 'In the query, not after it. Filtering in the handler means the record was already loaded, so any path that forgets the filter has the data in memory and usually returns it. Putting the predicate into the query means a record that is not yours is never fetched, and a forgotten check produces an empty result rather than a leak. That is the difference between a bug and a breach.',
        see: 'high-level-design',
      },
      {
        q: 'How do you make it impossible to forget?',
        a: 'By making the safe path the default one. The scope is applied inside the shared service layer that generated handlers call, so a new resource is isolated because it exists, not because somebody remembered. The remaining risk is a hand-written endpoint that queries directly, which is why the rule is that such endpoints go through the same service.',
        see: 'low-level-design',
      },
      {
        q: 'What happens when there is no identity on the request at all?',
        a: 'The scope matches nothing rather than everything. This is the detail that decides whether a mistake is survivable. A scope that quietly becomes a no-op when it cannot find an actor turns a middleware ordering error into "every user sees every record". Failing closed turns the same error into an empty list, which somebody notices immediately and nobody gets hurt by.',
        see: 'requirements',
      },
      {
        q: 'How do admins see everything without the scope blocking them?',
        a: 'The scope is skipped for an administrator, deliberately and in one place. The important part is that it is one place: a bypass implemented per endpoint will be applied inconsistently, and the inconsistency is invisible until somebody finds the endpoint where it was not. One bypass that is easy to find is also easy to audit.',
        see: 'low-level-design',
      },
      {
        q: 'What about reports and aggregates?',
        a: 'They are the usual leak, because a count or a sum is written as a direct query and nobody thinks of a number as data. A total revenue figure computed without the scope tells a user exactly how much everybody else is making. Aggregates go through the same scoped service as the lists.',
        see: 'bottlenecks',
      },
      {
        q: 'And caching?',
        a: 'A cache key that omits the owner is a data leak wearing a performance feature. The first user to request a resource populates the entry and the next user gets it. This is a nasty one because it presents intermittently, depending on who warmed the cache, and is almost impossible to reproduce deliberately. Everything that varies the answer belongs in the key.',
        see: 'bottlenecks',
      },
      {
        q: 'How do you test for this?',
        a: 'Write the test as the attack. Two users, one record each, then assert that user A gets a 404 on user B record for read, update and delete. Tests built from the normal flow never request somebody else id, so they can pass completely while the system is wide open.',
        see: 'low-level-design',
      },
      {
        q: 'Does the extra predicate cost anything?',
        a: 'Almost nothing, provided the index leads with the owner column. The filter makes queries more selective, not less, so it usually makes them faster. The failure case is an index on the sort column alone, where the database scans rows belonging to everybody before discarding them.',
        see: 'capacity',
      },
    ],
  },

  problem: {
    text: [
      'Permissions decide what kind of thing a caller may touch. They cannot decide which rows. A permission model alone gives every user with invoices.view the ability to read every invoice in the system, which is the single most common serious bug in web applications and has its own name: broken object level authorization.',
      'The reason it is so common is that the correct check is easy to write and easy to forget. Every list endpoint must narrow, every read by id must verify, every update and delete must verify, and the failure is silent: the endpoint works, the tests pass, and the data is wrong only from somebody else’s point of view.',
      'The fix is to move the check off the handler and into the layer that builds the query, so that a developer has to go out of their way to write an unscoped read rather than having to remember to scope one.',
    ],
    capabilities: [
      'Narrow a list query to the caller’s rows without the handler saying anything.',
      'Refuse a read, update or delete of a row the caller does not own.',
      'Answer 404 rather than 403 for a row they may not see, so the response does not confirm it exists.',
      'Let an administrator see everything, because the operator console reads the same endpoints.',
      'Match nothing at all when the caller is unknown, rather than opening up.',
      'Apply in a background job and a command, not only in a request.',
    ],
  },

  functional: [
    'A --owned-by flag on resource generation that adds the owner column and the scoping.',
    'Scoping applied inside the service, in the function that builds the query.',
    'The owner stamped from the caller on create, never read from the request body.',
    'The owner column not writable through update or patch.',
    'An administrator exempt from the narrowing, on the same endpoints.',
    'A guessed id answering 404, not 403.',
  ],

  nonFunctional: [
    {
      label: 'Default deny',
      text: 'a context with no actor matches no rows. A scope that opened up when it could not tell who was asking would read as protection and provide none.',
    },
    {
      label: 'In the service, not the handler',
      text: 'the handler is one caller of the service. A job, a command and a webhook are others, and all of them get the same narrowing.',
    },
    {
      label: 'Not leaky in errors',
      text: 'the difference between "does not exist" and "not yours" is not observable, because 403 on the second is itself a disclosure.',
    },
    {
      label: 'Not forgeable',
      text: 'the owner comes from the authenticated identity. A body field with the same name is ignored.',
    },
    {
      label: 'Testable as a breach',
      text: 'the property is "Ada cannot see Basil’s row", which is a test that fails loudly when the scoping is removed.',
    },
  ],

  capacity: {
    assumptions: [
      ['Owned resources in a typical project', '3 to 10'],
      ['Rows per owned table', 'up to tens of millions'],
      ['Rows per user', 'tens to thousands'],
      ['List queries', 'the majority of authenticated reads'],
    ],
    estimates: [
      {
        label: 'Query selectivity',
        working: [
          '10,000,000 invoices across 200,000 users',
          '= ~50 rows per user',
          'an index on (user_id, created_at) turns a full scan into ~50 rows read',
        ],
        note: 'The scoping is not only correctness, it is the difference between a sorted scan of ten million rows and an index read of fifty. The owner column is indexed for that reason.',
      },
      {
        label: 'Cost of the check on a read by id',
        working: [
          'one row loaded, one comparison against the actor',
          'no extra query: the ownership column is already on the row',
        ],
      },
    ],
  },

  highLevel: {
    intro:
      'Two primitives do the whole job. One narrows a query, the other judges a loaded row. Everything else is making sure they are the only way through.',
    components: [
      {
        label: 'ScopeOwned',
        text: 'takes a query and a column and returns a query narrowed to the actor. An administrator gets it back untouched; an unknown caller gets a query that matches nothing.',
      },
      {
        label: 'Owns',
        text: 'takes a loaded row and says whether the actor may see it. Used on read by id, update and delete.',
      },
      {
        label: 'Generated service',
        text: 'where both are called. The generator emits the calls, so an owned resource is scoped the moment it exists rather than when somebody remembers.',
      },
      {
        label: 'Create-time stamping',
        text: 'the owner is set from the context in the create path and ignored in the request, so it cannot be forged.',
      },
      {
        label: 'Isolation tests',
        text: 'generated alongside, asserting that one user cannot reach another’s rows and that the operator still can.',
      },
    ],
    flow: {
      title: 'A list and a read by id, both narrowed',
      nodes: [
        { id: 'client', label: 'Client', col: 0, row: 0 },
        { id: 'handler', label: 'Handler', col: 1, row: 0 },
        { id: 'service', label: 'Service', col: 2, row: 0, accent: true },
        { id: 'scope', label: 'ScopeOwned', col: 3, row: 0, accent: true },
        { id: 'db', label: 'Database', col: 2, row: 1 },
        { id: 'owns', label: 'Owns', col: 3, row: 1, accent: true },
        { id: 'notfound', label: '404 Not Found', col: 3, row: 2 },
        { id: 'rows', label: 'Rows Returned', col: 2, row: 2 },
      ],
      edges: [
        { from: 'client', to: 'handler', step: 1 },
        { from: 'handler', to: 'service', step: 2 },
        { from: 'service', to: 'scope', step: 3 },
        { from: 'service', to: 'db', step: 4 },
        { from: 'db', to: 'owns', step: 5, bend: 'h' },
        { from: 'owns', to: 'notfound', step: 6, dashed: true },
        { from: 'db', to: 'rows', step: 7 },
      ],
      steps: [
        'The request arrives already authenticated and already past the permission check. Permission said "may read invoices"; nothing yet has said which.',
        'The handler calls the service. It passes the context and no ownership argument, because it does not know about ownership.',
        'The service builds the query and passes it through ScopeOwned, which adds the predicate. An administrator gets the query unchanged; a context with no actor gets one that matches nothing.',
        'The narrowed query reaches the database. The owner column is indexed, so this is also the difference between reading fifty rows and scanning ten million.',
        'On a read by id the row comes back first and is judged by Owns, because the primary key lookup does not carry the predicate.',
        'A row the actor does not own is reported as not found. Answering 403 would confirm that an invoice with that id exists, which is the disclosure the check exists to prevent.',
        'Otherwise the rows are returned. The handler never saw the decision, which is what keeps it from being forgotten in a new endpoint.',
      ],
    },
    dataFlow: [
      'The owner column is written once, from the identity, at creation. It is excluded from the update and patch paths, so it cannot be reassigned by a request.',
      'The column name reaching the query builder is a fixed string from the generator, never anything derived from input.',
      'Administrators are exempt by the same primitive, so the operator console uses the same endpoints as the customer and there is no second code path to keep in step.',
      'A background job that acts for a user carries the same actor on its context, so the narrowing applies there too.',
    ],
  },

  stack: [
    ['Scoping primitive', 'a query modifier applied in the service'],
    ['Ownership column', 'a foreign key to users, indexed'],
    ['Identity source', 'the actor on context.Context'],
    ['Exemption', 'the administrator flag on the actor'],
    ['Failure mode', 'match nothing when the caller is unknown'],
    ['Verification', 'generated isolation tests over HTTP'],
  ],

  dataModel: {
    intro:
      'There is no table for this system. It is a column on the resources that have an owner, plus the discipline about where it is read and written.',
    entities: [
      {
        name: 'any owned resource',
        fields: [
          ['user_id', 'the owner, indexed, set at creation from the identity'],
          ['...', 'the resource’s own columns'],
        ],
        note: 'The index matters as much as the column. Without (user_id, created_at) a scoped list is a filtered sort of the whole table.',
      },
    ],
  },

  lowLevel: {
    classes: [
      {
        name: 'ScopeOwned',
        file: 'internal/authz/actor.go',
        what: 'The query narrowing. Three outcomes: untouched for an administrator, impossible for an unknown caller, narrowed for everybody else. The order of those cases is the security property.',
      },
      {
        name: 'Owns',
        file: 'internal/authz/actor.go',
        what: 'The row judgement, for the paths where a primary key lookup has already happened.',
      },
      {
        name: 'Ownable',
        what: 'A one-method interface: GetOwnerID. A model satisfies it by having an owner, which is how the generator knows the resource is scoped.',
      },
      {
        name: 'OwnsOr404',
        file: 'internal/security/...',
        what: 'The handler-side convenience that converts a failed ownership check into a not-found rather than a forbidden.',
      },
    ],
    principles: [
      {
        label: 'Make the safe path the easy path',
        text: 'the generator writes the scoping. A developer has to delete a line to produce an unscoped read, which is a visible act in review rather than an omission.',
      },
      {
        label: 'Fail closed, loudly',
        text: 'an unknown actor produces a query that matches nothing. The endpoint returns an empty list, which is noticed, rather than everything, which is not.',
      },
      {
        label: 'Enforce where every caller passes',
        text: 'the service is the narrowest point that all callers share. The handler is not, because jobs and commands do not go through it.',
      },
      {
        label: 'Do not leak through status codes',
        text: 'the distinction between forbidden and missing is itself information, so it is not offered.',
      },
    ],
    patterns: [
      ['Query scope', 'a predicate added by a shared function rather than by each call site'],
      ['Guard clause', 'ownership judged before the row is returned or mutated'],
      ['Ambient context', 'the identity travels on the context, so background work inherits it'],
    ],
  },

  scaling: [
    'The predicate makes queries more selective, not less. On a large table the scoped version is dramatically faster than the unscoped one it replaces.',
    'The composite index on the owner column and the sort column is what keeps a scoped list from degrading into a sort of the whole table.',
    'Nothing is cached, because a cache keyed without the owner is exactly the bug this system prevents.',
    'The check costs no additional round trip. The ownership column arrives with the row that was being read anyway.',
    'As a table grows, scoping becomes more valuable rather than less: the work per request stays proportional to one user’s data.',
  ],

  bottlenecks: {
    problems: [
      {
        label: 'A hand-written endpoint that forgets',
        text: 'the generator scopes what it generates. A custom endpoint added later is where this breaks, and it breaks silently.',
      },
      {
        label: 'Aggregates and reports',
        text: 'a count or a sum written as raw SQL bypasses the service and the scope with it, and the number it returns is the whole table.',
      },
      {
        label: 'Caching by id',
        text: 'a cache key that does not include the owner will serve one user’s row to another, turning a correct system into an incorrect one.',
      },
      {
        label: 'Relationships that cross the fence',
        text: 'an owned invoice pointing at a customer that is not owned means a user can enumerate customers through the relationship.',
      },
      {
        label: 'Missing index',
        text: 'the scope is still correct without an index on the owner column, and still slow. Correctness hides the performance problem until the table is large.',
      },
    ],
    improvements: [
      {
        label: 'Isolation tests that would catch a breach',
        text: 'tests asserting that one user cannot see another’s rows, verified to fail when the scoping line is deleted. A test that passes either way proves nothing.',
      },
      {
        label: 'Route the aggregates through the service',
        text: 'counts and sums take the same scoped query builder, so a report cannot see further than a list.',
      },
      {
        label: 'Owner in the cache key',
        text: 'every cache entry for an owned resource is keyed by owner as well as id, which makes the cross-user case impossible rather than unlikely.',
      },
      {
        label: 'Index with the sort',
        text: 'a composite index on owner and the default sort column, created with the resource rather than added after the first slow query.',
      },
    ],
  },

  seeAlso: [
    { title: 'Authorization', href: '/docs/systems/authorization' },
    { title: 'Multitenancy', href: '/docs/systems/multitenancy' },
  ],
}
