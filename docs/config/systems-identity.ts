import type { SystemDesign } from './systems-types'

/* Identity: who the caller is. Authentication, sessions, MFA, passkeys. */

export const AUTHENTICATION: SystemDesign = {
  slug: 'authentication',
  name: 'Authentication',
  tagline:
    'Proving who the caller is, on every request, without asking them for a password more than once.',
  group: 'Identity',
  packages: ['internal/handlers', 'internal/password', 'internal/middleware', 'internal/crypto'],

  interview: {
    intro:
      '"Design an authentication system" is one of the most common system design interview prompts, and the follow-up questions are remarkably consistent between companies. Here are the ones that come up, each with a short answer. Every one of them is worked through somewhere below, so if you can follow this page you can answer the question.',
    questions: [
      {
        q: 'Design a login system. Where do you start?',
        a: 'With requirements, not with tokens. Ask who is signing in (people in a browser, a mobile app, another service), whether you need to sign somebody out instantly, whether a second factor is required, and roughly how many users there are. Those four answers decide almost everything else. Only then write the functional list: register, sign in, stay signed in, sign out, reset a forgotten password. Jumping straight to "I would use JWTs" is the mistake interviewers are watching for, because it answers a question nobody asked yet.',
        see: 'requirements',
      },
      {
        q: 'How do you store passwords?',
        a: 'Never as text, and never with a fast hash like SHA-256. You use a function built to be slow, such as bcrypt, scrypt or Argon2. Slow is the feature: a function that takes 60 milliseconds instead of a microsecond makes guessing a stolen hash roughly 60,000 times more expensive. The function also mixes in a salt, a random value stored next to the hash, so two people with the same password get different hashes and an attacker cannot crack them all at once with a precomputed table. bcrypt generates and stores the salt inside its own output, so there is no separate column to manage.',
        see: 'data-model',
      },
      {
        q: 'Sessions or JWTs? Which would you pick and why?',
        a: 'This is a trade-off question, and naming a winner is the wrong answer. A session is a row in your database; the client holds a meaningless id, and you look it up on every request. That costs a database read per request and makes signing somebody out instant, because you delete the row. A JWT is a self-contained token the server can check with mathematics alone, so there is no lookup, but you cannot un-issue it before it expires. The design here uses both, deliberately: a short-lived JWT for the per-request check, and a session row behind the long-lived token so revocation still works. You get stateless verification and instant sign-out, at the cost of one database write when a token is renewed.',
        see: 'problem',
      },
      {
        q: 'A JWT cannot be cancelled once issued. So how do you log someone out?',
        a: 'You make the uncancellable thing short-lived, and put the cancellable thing behind it. The access token lives 15 minutes and nothing can stop it inside that window. The refresh token lives 7 days, and it only works if a matching session row exists in the database. Signing out marks that row revoked, so the next renewal fails and the user is out within at most 15 minutes. If a product truly needs instant revocation of the access token too, the usual answer is a deny list of token ids in Redis checked on every request, which trades the per-request lookup back in.',
        see: 'high-level-design',
      },
      {
        q: 'Users should stay signed in for a week, but be signed out after a period of inactivity. How?',
        a: 'Two separate clocks on the same session row, which is a question people often get half right. The idle timeout moves forward every time the session is used, so a person who keeps working never gets kicked out. The absolute timeout never moves, so a session dies at a fixed point no matter how active it is. Here the idle timeout is 7 days and the absolute timeout is 30. An unused sign-in dies in a week; a constantly used one still dies in a month. Without the absolute clock a session that is touched often enough lives forever.',
        see: 'data-model',
      },
      {
        q: 'Where should a browser keep the token? localStorage or a cookie?',
        a: 'A cookie with the HttpOnly flag, which means JavaScript cannot read it. Anything in localStorage is readable by any script on the page, so one cross-site scripting bug hands over every token. The cookie approach has its own weakness, cross-site request forgery, where another site makes the browser send your cookie along with a request you did not intend; the defence is the SameSite attribute plus checking the origin. Here the refresh token is an HttpOnly cookie scoped to the /api/auth path, so it is not even attached to ordinary API calls, and the access token is sent in the Authorization header where no browser attaches it automatically.',
        see: 'high-level-design',
      },
      {
        q: 'How do you stop someone guessing passwords?',
        a: 'Count failures per account and bar it temporarily once there are too many: here ten failures inside a window bars the account for fifteen minutes. The important detail is that the check runs before the password is hashed. If you hash first and check the limit afterwards, an attacker can still force you to burn 60 milliseconds of processor time per guess, which is a denial of service rather than a break-in but takes you down just the same.',
        see: 'high-level-design',
      },
      {
        q: 'What if the attacker spreads the guesses across a thousand different IP addresses?',
        a: 'Then a limit counted per IP address does nothing, because no single address reaches it. This is why the lockout counts failures per account. The mirror-image mistake is to lock out the IP address instead: an attacker who can guess your office address can then deliberately lock out everybody in your building. Count per account for the lockout, and keep a per-address rate limit as a separate, looser control on top.',
        see: 'high-level-design',
      },
      {
        q: 'How do you avoid telling an attacker which email addresses are registered?',
        a: 'By making "no such account" and "wrong password" indistinguishable, in both the message and the time taken. Returning different errors turns the login form into a lookup service for valid addresses. Timing matters as much as wording: if an unknown address returns instantly while a real one takes 60 milliseconds because a hash was computed, the difference is measurable, so the unknown-address path does the same wasted work. The same rule applies to the forgotten-password endpoint, which answers identically whether or not the address exists.',
        see: 'requirements',
      },
      {
        q: 'What is refresh token rotation, and how would you detect a stolen token?',
        a: 'Rotation means every renewal issues a new refresh token and invalidates the one presented, so any given token is usable exactly once. Detection falls out of that for free: the session row remembers the hash of the token it just replaced, so if the old token is presented again you know there are two copies in the world. One of them is a thief. The system cannot tell which, so the safe response is to revoke the whole session and make the real user sign in again. The subtlety worth raising in an interview is the false alarm: a client that retries after a dropped response looks exactly like a replay, which is why a short grace window for an immediate retry is a common refinement.',
        see: 'data-model',
      },
      {
        q: 'How do other services check the token without calling your auth service?',
        a: 'That is the property a JWT is bought for. The token carries its claims (who the user is, their role, when it expires) and a signature over them. Any service holding the key can recompute that signature and know the contents were not altered, using only processor time and the clock. Nothing is called and nothing is looked up. With a shared secret, as here, every service needs that secret, so the usual step for a larger estate is to sign with a private key and let services verify with the matching public one, which means a leaked verifier cannot mint tokens.',
        see: 'high-level-design',
      },
      {
        q: 'What is the difference between authentication and authorization?',
        a: 'Authentication is who you are; authorization is what you may do. They are different systems and conflating them is a classic source of real breaches: a token can be perfectly valid and still not entitle its holder to the invoice they just requested. Here the auth middleware does exactly one job, turning a token into an identity on the request context, and it deliberately makes no decision about permission. The next layer reads that identity and decides. Keeping them apart is why an ownership check cannot be accidentally satisfied by the mere fact that somebody is signed in.',
        see: 'low-level-design',
      },
      {
        q: 'What happens to signed-in users when the database goes down?',
        a: 'They keep working, for up to 15 minutes. This is the quiet benefit of stateless verification: checking an access token touches no storage, so every already-authenticated read keeps succeeding through a database outage. What stops is anything needing a write: new sign-ins, renewals, and sign-outs. The failure is therefore gradual rather than instant, and the window is exactly the access token lifetime, which is a useful thing to be able to state precisely when an interviewer asks what breaks first.',
        see: 'scaling',
      },
      {
        q: 'How do you rotate the signing secret without signing everybody out?',
        a: 'With one key you cannot, which is the point of the question. Changing it invalidates every token in existence at once. The fix is to verify against two keys while signing with only one: add the new key as an accepted verifier, switch signing over to it, and once one access token lifetime has passed, nothing signed with the old key is still valid and it can be dropped. Rotation then takes 15 minutes instead of logging out every user.',
        see: 'bottlenecks',
      },
      {
        q: 'How does multi-factor authentication fit into the login flow?',
        a: 'Sign-in stops halfway. Instead of returning tokens, a correct password returns a short-lived pending token that proves the first factor passed and nothing else. The client sends that back along with the code from the authenticator app, and only then are real tokens issued. The important part is that the pending token carries no authority: it cannot be used on any other endpoint, so an attacker holding a correct password is still holding nothing useful.',
        see: 'api',
      },
      {
        q: 'How do you handle a forgotten password safely?',
        a: 'Email a single-use token, store only its hash, and expire it in minutes rather than days. Single-use means it stops working the moment it is used, so an emailed link sitting in a mailbox is not a permanent key to the account. Storing only the hash means a database leak does not hand over working reset links. And the endpoint answers identically whether or not the address exists, for the same enumeration reason as the login form.',
        see: 'api',
      },
      {
        q: 'How would you size this? How many servers, how much storage?',
        a: 'Work from active users rather than registered ones, because dormant accounts cost storage and nothing else. Multiply through to get requests per second, then check the two expensive operations separately: password hashing, which is deliberately slow and therefore caps sign-ins, and token verification, which is nearly free and therefore does not cap anything. Every figure in the capacity section below shows where it came from, so you can change one assumption and redo the arithmetic.',
        see: 'capacity',
      },
    ],
  },

  problem: {
    text: [
      'Every endpoint past the public ones needs to know which user is calling. The obvious answer is to send the password with every request, which means the client has to store it and hand it over a thousand times a day. That is wrong for a reason worth stating plainly: a password is a long-lived secret, and the more places and times it travels, the more chances there are to lose it.',
      'The second answer is a session: the server keeps a row for each signed-in person, the client holds a meaningless identifier, and the server looks it up on every request. This works, and plenty of large systems do exactly this. The cost is that it puts a database read in front of every single thing the application does.',
      'The third answer is a token that carries its own proof. A JSON Web Token (JWT) is a small piece of JSON, such as "user 123, role editor, expires at 09:15", with a signature over it. Any server holding the key can recompute that signature and know nobody tampered with the contents, without looking anything up. Verification becomes arithmetic instead of a query.',
      'That buys speed and gives up control. Once a token is issued it is valid until it expires, and there is no way to take it back. If somebody steals it, or a user signs out, or an account is suspended, the token carries on working regardless.',
      'So this design uses both, and the split is the central idea of the page. A short-lived access token (15 minutes) does the per-request work with no lookup. A long-lived refresh token (7 days) renews it, and that one is backed by a row in the database, so it can be revoked. You get stateless verification on the hot path and real revocation where it matters, at the cost of one database write each time a token is renewed.',
      'The hard parts are not the signing, which is a library call. They are what happens when a token is stolen, what happens when a user signs out on one device but not another, and what happens when somebody tries ten thousand passwords against one account from a thousand different addresses.',
    ],
    capabilities: [
      'Register a user, hash the password with a function that is slow on purpose, and never store the plain text.',
      'Issue a credential a request can carry, and verify it without reading the database.',
      'Renew that credential without a second sign-in, and notice when a renewal is a replay of an old one.',
      'Revoke one device without signing the user out everywhere, and revoke everywhere when a password changes.',
      'Resist password guessing spread across many addresses, without letting an attacker lock out a user whose address they can guess.',
      'Hand the caller to the authorization system as an identity, so that layer never has to parse a token.',
    ],
  },

  functional: [
    'Registration with email and password, and an optional requirement that the address be confirmed before the first sign-in.',
    'Sign-in returning an access token and a refresh token, with the refresh token also set as an HttpOnly cookie scoped to /api/auth.',
    'Token refresh that rotates the refresh token and invalidates the one presented.',
    'Sign-out for the current device, and sign-out for every device.',
    'Password reset by emailed single-use token, and password change for a signed-in user.',
    'Social sign-in through Google and GitHub, landing on the same user record.',
    'A context-level identity every handler and service reads, rather than each one parsing the Authorization header.',
  ],

  nonFunctional: [
    {
      label: 'Verification cost',
      text: 'checking an access token must not touch the database. It is a signature check and a clock comparison, both of which take microseconds, and that is what makes it affordable to do on every single request.',
    },
    {
      label: 'Blast radius',
      text: 'this is the question of how much damage one stolen thing can do. A stolen access token expires in 15 minutes. A stolen refresh token is usable once, and using it after the real client already has tells the server something is wrong.',
    },
    {
      label: 'Storage safety',
      text: 'somebody who steals a copy of the database must not be able to sign in with it. Passwords are stored as bcrypt hashes and refresh tokens as SHA-256 digests, so neither column contains anything that can be replayed.',
    },
    {
      label: 'Availability',
      text: 'signing in may depend on the database. Verifying a token may not. An application therefore stays readable through a database outage for as long as its access tokens live.',
    },
    {
      label: 'No membership oracle',
      text: 'an oracle here means anything that answers a question an attacker should not get to ask. A wrong password and an unknown address answer the same way, in the same time, so the login form cannot be used to discover which addresses are registered.',
    },
  ],

  capacity: {
    assumptions: [
      [
        'Registered users',
        '1,000,000 (a successful mid-market product; pick your own number and every sum below rescales with it)',
      ],
      [
        'Daily active users',
        '200,000 (20% of registered, a normal ratio for a tool people use on working days; a consumer social app would be higher, a tax product far lower)',
      ],
      [
        'Sign-ins per active user per day',
        '1.2 (most people are already signed in and never see the form; the 0.2 covers a second device or a session that expired overnight)',
      ],
      [
        'API requests per active user per day',
        '400 (one screen typically makes 5 to 15 API calls, and an engaged user opens 30 to 80 screens in a working day)',
      ],
      [
        'Access token lifetime',
        '15 minutes (short enough that a stolen token is nearly worthless, long enough that renewals stay rare; the industry norm is 5 to 60 minutes)',
      ],
      [
        'Refresh token lifetime',
        '168 hours, which is 7 days (how long somebody can close the laptop and come back without signing in again; a product decision, not a technical one)',
      ],
      [
        'Session row',
        '~400 bytes (two 32-byte hashes, two timestamps, a user agent string of around 200 bytes, an address, and row overhead)',
      ],
      [
        'User row',
        '~1.2 KB with profile fields (identifier, email, a 60-byte bcrypt hash, name, role, timestamps, and whatever the product adds)',
      ],
    ],
    estimates: [
      {
        label: 'Sign-in rate: how often the expensive path runs',
        working: [
          '200,000 daily active users x 1.2 sign-ins = 240,000 sign-ins/day',
          '240,000 / 86,400 = 2.8 sign-ins/second average (86,400 is the number of seconds in a day)',
          'peak = 2.8 x 5 = 14 sign-ins/second (5x is the usual ratio of the busiest hour to the daily average when an audience shares working hours; a global consumer app is flatter, nearer 2x)',
        ],
        note: 'This matters because sign-in is the one deliberately slow operation. bcrypt at cost 10 takes roughly 60 milliseconds of processor time, so 14 a second is 14 x 0.06 = 0.84 seconds of work per second, which is under one processor core. That is the number that tells you the cost factor can stay where it is.',
      },
      {
        label: 'Token verification rate: how often the cheap path runs',
        working: [
          '200,000 daily active users x 400 requests = 80,000,000 requests/day',
          '80,000,000 / 86,400 = 926 verifications/second average',
          'peak = 926 x 5 = 4,630 verifications/second (same 5x peaking factor as above)',
        ],
        note: 'Verifying a signature takes a few microseconds, so 4,630 a second is roughly 0.02 seconds of work per second: effectively nothing. Compare it with the session approach, where these same 4,630 requests a second would each be a database read. That gap is the entire argument for stateless access tokens.',
      },
      {
        label: 'Refresh rate: the part that is not free',
        working: [
          'an 8 hour working day / a 15 minute token = 32 renewals per active session',
          '200,000 daily active users x 32 = 6,400,000 renewals/day',
          '6,400,000 / 86,400 = 74 renewals/second average, and 370/second at the same 5x peak',
        ],
        note: 'Every renewal is a database write, because rotation replaces the stored hash. This is the only part of the hot path with state, so it is the number to watch. Note how it moves: halving the access token lifetime to 7 minutes doubles this to 148 and 740. Lengthening it to an hour cuts it to 18. That single assumption is the main lever on write load.',
      },
      {
        label: 'Session storage: how big the revocation table gets',
        working: [
          '200,000 daily active users x 3 devices = 600,000 live sessions (laptop, phone, and one more: a tablet, a work machine, a second browser)',
          '600,000 x 400 bytes = 240 MB',
          'plus indexes on token_hash, user_id and expires_at = ~400 MB (indexes commonly add 50 to 100% on a table this narrow)',
        ],
        note: 'Under a gigabyte, so it stays in the primary database and needs no separate store. It stays that size only because expired rows are swept on a schedule; left alone the table would grow with every sign-in ever made rather than with devices currently signed in.',
      },
      {
        label: 'User storage: the part that grows with signups, not usage',
        working: [
          '1,000,000 users x 1.2 KB = 1.2 GB',
          'with indexes and audit columns = ~2 GB',
        ],
        note: 'Small enough that it never drives a design decision, which is worth saying out loud: dormant accounts cost storage and nothing else. This is why the sums above start from active users rather than registered ones.',
      },
    ],
  },

  highLevel: {
    intro:
      'Sign-in is a write path that ends in two tokens. Everything after it is a read path that touches no storage at all until the access token expires.',
    components: [
      {
        label: 'Auth handlers',
        text: 'the HTTP surface. Register, login, refresh, logout, forgot and reset. They are thin on purpose: they read the request, call a service, and shape the response, holding no logic that would have to be tested through HTTP.',
      },
      {
        label: 'Password service',
        text: 'bcrypt hashing and comparison, plus the strength rules. The only place in the system that ever sees a plain text password, and it never returns one.',
      },
      {
        label: 'JWT service',
        text: 'signs and verifies. Every token carries a user id, a role, a type, an expiry and a jti, which is a unique id for that individual token, so two tokens issued in the same second are still distinguishable from each other.',
      },
      {
        label: 'Session store',
        text: 'one row per signed-in device, found by the SHA-256 hash of the current refresh token. This single table is what turns a stateless token into a revocable one, and it is the reason this design is not purely stateless.',
      },
      {
        label: 'Auth middleware',
        text: 'runs before your handler. It reads the Authorization header or the access cookie, verifies the token, and puts an identity on the request context so everything downstream reads a value instead of parsing a header. It decides who you are and nothing about what you may do.',
      },
      {
        label: 'Lockout counter',
        text: 'failed attempts counted per account rather than per address, so an attack spread across a thousand machines still trips it.',
      },
    ],
    flow: {
      title: 'Sign-in and the first authenticated request',
      nodes: [
        { id: 'client', label: 'Client', sub: 'Web / Mobile', col: 0, row: 0 },
        { id: 'mw', label: 'Auth Middleware', col: 1, row: 0 },
        { id: 'handler', label: 'Auth Handler', col: 2, row: 0, accent: true },
        { id: 'password', label: 'Password Service', col: 3, row: 0 },
        { id: 'lockout', label: 'Lockout Counter', col: 3, row: 1 },
        { id: 'jwt', label: 'JWT Service', col: 2, row: 1, accent: true },
        { id: 'sessions', label: 'Session Store', col: 1, row: 1 },
        { id: 'db', label: 'Database', col: 2, row: 2 },
        { id: 'audit', label: 'Audit Log', col: 3, row: 2 },
      ],
      edges: [
        { from: 'client', to: 'mw', step: 1 },
        { from: 'mw', to: 'handler', step: 2 },
        { from: 'handler', to: 'password', step: 3 },
        { from: 'password', to: 'lockout', step: 4, dashed: true },
        { from: 'handler', to: 'jwt', step: 5 },
        { from: 'jwt', to: 'sessions', step: 6 },
        { from: 'sessions', to: 'db', step: 7, bend: 'v' },
        { from: 'jwt', to: 'db', step: 8 },
        { from: 'db', to: 'audit', step: 9, dashed: true },
      ],
      steps: [
        'The client posts an email and password to /api/v1/auth/login. The middleware lets it straight through, because this route is public: it is the route you use when you have nothing to prove yet.',
        'The handler reads the request and asks the lockout counter whether this account is currently barred. A barred account is refused right here, before any password hashing is paid for, which is what stops a guessing attack from costing processor time.',
        'The password service compares the submitted password against the stored bcrypt hash. A wrong password and an unknown email take the same path and the same amount of time, so the response cannot be used to find out which addresses exist.',
        'A failed comparison increments the counter for that account. Ten failures inside the window bar it for fifteen minutes. An unknown address increments nothing, because otherwise an attacker could deliberately lock out any address they can guess.',
        'On success the JWT service mints two tokens: an access token good for 15 minutes, and a refresh token good for 7 days. Each carries its own jti, the unique per-token id, so they can be told apart later.',
        'The session store records the device: the SHA-256 hash of the refresh token, the browser or app it came from, the address, and the two expiry timestamps.',
        'The row is written. Note what is not written: the refresh token itself. Only its hash is stored, so even somebody holding a dump of this table cannot produce a working token from it.',
        'The response carries both tokens, and also sets the refresh token as a cookie marked HttpOnly, so JavaScript cannot read it, and scoped to the /api/auth path, so the browser does not attach it to ordinary API calls.',
        'The sign-in is written to the audit log, outside the response path so it cannot slow the user down. From here on, each request carries the access token and the middleware verifies it with no storage read at all, which is the whole point of the arrangement.',
      ],
    },
    dataFlow: [
      'Passwords travel over TLS, are hashed the moment they arrive, and are never written to a log, a trace or an error message. A password in a log file is a password in your log search tool, your backups and your screenshots.',
      'Access tokens are stored nowhere on the server. Their validity is entirely a property of the signature and the clock, which is exactly why no server has to be asked about them and why they cannot be cancelled early.',
      'Refresh tokens are stored only as SHA-256 hashes. The server can recognise a token it issued without being able to produce one, the same trick used for passwords and for the same reason.',
      'The hash of the previous refresh token is kept alongside the current one. That one extra column is what makes replay detection possible: it is how the server can tell that a token it already replaced has just been presented again.',
      'The identity placed on the request context is a plain value, not a database handle. A service that needs the user id gets the user id, and does not get the ability to run queries as that user.',
      'Browsers and mobile apps are not the same problem. A mobile app has no cross-site request forgery risk and no JavaScript injection surface, so it holds tokens in the platform keychain and sends them in the Authorization header. The cookie machinery exists for browsers.',
    ],
  },

  stack: [
    ['Token format', 'JWT, HS256, signed with JWT_SECRET'],
    ['Password hashing', 'bcrypt, cost 10 (about 60 ms per hash)'],
    ['Session storage', 'the primary database, one row per device'],
    ['Transport', 'Authorization: Bearer, plus HttpOnly cookies for browsers'],
    ['Social sign-in', 'OAuth2 against Google and GitHub'],
    ['Lockout counter', 'database-backed, so it survives a restart'],
    ['Reset tokens', 'single-use, hashed at rest, short expiry'],
  ],

  dataModel: {
    intro:
      'Three tables carry the whole system: a user, a session per device, and a short-lived token for the password reset flow.',
    entities: [
      {
        name: 'users',
        fields: [
          ['id', 'UUIDv7, which sorts by creation time so the index stays compact as rows are added'],
          ['email', 'unique, and the identifier somebody signs in with'],
          ['password_hash', 'the bcrypt output, which contains the salt and the cost inside it; never the password'],
          ['role', 'ADMIN, EDITOR or USER'],
          ['email_verified_at', 'null until confirmed; gates sign-in when verification is required'],
          ['totp_secret', 'encrypted at rest, null until two-factor is enrolled'],
          ['created_at, updated_at', 'audit columns'],
        ],
        note: 'There is no salt column, which surprises people. bcrypt generates a random salt per password and stores it inside the hash string, so the single column holds the algorithm, the cost, the salt and the hash together.',
      },
      {
        name: 'sessions',
        note: 'One row per signed-in device. This table is the thing that turns a stateless token into a revocable one, and the two timestamp columns are how one session can have both an idle timeout and an absolute one.',
        fields: [
          ['id', 'UUIDv7'],
          ['user_id', 'the owner, indexed'],
          ['token_hash', 'SHA-256 of the current refresh token, unique'],
          ['prev_token_hash', 'SHA-256 of the token this one replaced; the entire replay detection mechanism is this column'],
          ['user_agent, ip', 'what the device list shows a user, so they can recognise their own sessions'],
          ['last_seen_at', 'moved forward on every renewal, which is what drives the idle timeout'],
          ['expires_at', 'the absolute deadline, which never moves no matter how active the session is'],
          ['revoked_at', 'set on sign-out; the row is kept rather than deleted, so a later replay can still be recognised'],
        ],
      },
      {
        name: 'password_reset_tokens',
        fields: [
          ['token_hash', 'the emailed token, hashed, so a database leak yields no working links'],
          ['user_id', 'who it is for'],
          ['expires_at', 'short, measured in minutes rather than days'],
          ['used_at', 'set on first use, which is what makes the link work exactly once'],
        ],
      },
    ],
    storage: [
      'Everything is in the primary relational database. There is no second store to keep consistent, and a session lookup is a single indexed read.',
      'Sessions are swept on a schedule: rows past their absolute expiry are deleted, and revoked rows go once they are older than the retention window. Without the sweep this table grows with every sign-in ever made.',
      'The idle timeout is 7 days and the absolute timeout is 30. An unused sign-in dies in a week; one that is used every day still dies in a month. Two clocks, because one cannot express both rules.',
    ],
  },

  api: {
    groups: [
      {
        title: 'Authentication',
        rows: [
          { method: 'POST', path: '/api/v1/auth/register', what: 'Create an account and sign in' },
          { method: 'POST', path: '/api/v1/auth/login', what: 'Exchange credentials for tokens' },
          { method: 'POST', path: '/api/v1/auth/refresh', what: 'Rotate the refresh token, issue a new access token' },
          { method: 'POST', path: '/api/v1/auth/logout', what: 'Revoke the current session' },
          { method: 'GET', path: '/api/v1/auth/me', what: 'The signed-in user' },
        ],
      },
      {
        title: 'Passwords',
        rows: [
          { method: 'POST', path: '/api/v1/auth/forgot-password', what: 'Email a single-use reset link, answering the same way whether or not the address exists' },
          { method: 'POST', path: '/api/v1/auth/reset-password', what: 'Set a new password from that link' },
          { method: 'POST', path: '/api/v1/auth/change-password', what: 'Change it while signed in, ending every other session' },
        ],
      },
      {
        title: 'Sessions',
        rows: [
          { method: 'GET', path: '/api/v1/auth/sessions', what: 'Every device currently signed in' },
          { method: 'DELETE', path: '/api/v1/auth/sessions/:id', what: 'Sign one device out' },
          { method: 'POST', path: '/api/v1/auth/sessions/revoke-all', what: 'Sign every device out' },
        ],
      },
    ],
    samples: [
      {
        title: 'Sign-in request',
        language: 'json',
        code: `{
  "email": "ada@example.com",
  "password": "correct horse battery staple"
}`,
      },
      {
        title: 'Sign-in response. expires_in is 900 seconds, which is the 15 minute access token',
        language: 'json',
        code: `{
  "data": {
    "user": {
      "id": "01a11e82-d0bc-7367-9c39-5e9873953dde",
      "email": "ada@example.com",
      "first_name": "Ada",
      "role": "USER"
    },
    "tokens": {
      "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
      "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
      "expires_in": 900
    }
  },
  "message": "Signed in"
}`,
      },
      {
        title: 'When two-factor is enrolled, sign-in stops halfway and no real token is issued',
        language: 'json',
        code: `{
  "data": {
    "mfa_required": true,
    "pending_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "methods": ["totp", "backup_code"]
  },
  "message": "A second factor is required"
}`,
      },
    ],
  },

  lowLevel: {
    classes: [
      {
        name: 'AuthHandler',
        file: 'internal/handlers/auth.go',
        what: 'Reads and validates the request, calls the services, shapes the response. Holds no logic that would be worth testing through HTTP.',
        methods: ['Register', 'Login', 'Refresh', 'Logout', 'Me', 'ForgotPassword', 'ResetPassword'],
      },
      {
        name: 'JWTService',
        file: 'internal/services/jwt.go',
        what: 'The only thing in the system that signs or verifies. Issuing and verifying are separate methods, so a service that only needs to verify can exist without ever holding the ability to mint tokens.',
        methods: ['GenerateTokenPair', 'ValidateAccessToken', 'ValidateRefreshToken'],
      },
      {
        name: 'SessionStore',
        file: 'internal/session/session.go',
        what: 'Creates, rotates and revokes device sessions, and decides whether a presented refresh token belongs to a live one.',
        methods: ['CreateSession', 'RotateSession', 'Revoke', 'RevokeAllForUser', 'ListForUser'],
      },
      {
        name: 'PasswordService',
        file: 'internal/password/password.go',
        what: 'Hash, compare, and judge strength. Comparison takes the same time whether or not the user exists, which is a correctness requirement rather than an optimisation.',
        methods: ['Hash', 'Compare', 'Strength'],
      },
      {
        name: 'AuthMiddleware',
        file: 'internal/middleware/auth.go',
        what: 'Turns a token into an identity on the request context, or returns a 401. It does not decide what that identity is allowed to do, and keeping that line clean is what stops "is signed in" from being mistaken for "is allowed".',
        methods: ['RequireAuth', 'OptionalAuth'],
      },
    ],
    principles: [
      {
        label: 'Single responsibility',
        text: 'the handler does HTTP, the JWT service does cryptography, the session store does revocation. No one of them knows about two of those three.',
      },
      {
        label: 'Authentication is not authorization',
        text: 'the middleware answers who, and stops. A separate layer answers what they may touch. These are different questions, and the breaches that come from merging them look like a valid token being treated as permission for a record it was never entitled to.',
      },
      {
        label: 'Dependency inversion',
        text: 'services take interfaces for the clock and the store, which is a plain way of saying the expiry tests can move time forward instead of sleeping for 15 minutes.',
      },
      {
        label: 'Open for extension',
        text: 'password, social and passkey sign-in all end at the same point, issuing a token pair for an existing user record, so adding one never required changing the token layer.',
      },
      {
        label: 'Fail closed',
        text: 'an unparseable token, an unknown session and a revoked session all produce the same 401. There is no path through the code where an error leaves the request authenticated.',
      },
    ],
    patterns: [
      ['Strategy', 'password, social and passkey sign-in, each ending at the same token issuance'],
      ['Chain of responsibility', 'the middleware stack: request id, then auth, then rate limit, then the handler'],
      ['Token rotation', 'every renewal replaces the stored hash and remembers the old one'],
      ['Context carrier', 'the identity rides on context.Context, so a background job started by a request keeps it'],
    ],
  },

  scaling: [
    'Verification is stateless, so the API scales horizontally with no shared session store and no sticky sessions. Adding a replica adds capacity immediately, because a replica needs nothing from the others to check a token.',
    'A database outage degrades this system gradually rather than instantly. Already-signed-in users keep reading for up to one access token lifetime, 15 minutes, because verification reads no storage. What stops is sign-ins, renewals and sign-outs, all of which write.',
    'The refresh endpoint is the only part of the hot path that writes. At 370 writes a second it is comfortable on one primary, and it is the first thing to watch if the access token lifetime is ever shortened.',
    'That trade is worth stating as a rule: shortening the access token lifetime buys a smaller window of damage from a stolen token, and pays for it in writes. Halving it to 7 minutes doubles the write rate on that one endpoint.',
    'Session reads are a single lookup on a unique index. The table stays small because it holds devices currently signed in, not a history of every sign-in.',
    'Bcrypt cost is a deliberate ceiling on sign-in throughput. It is sized so a peak hour fits in a fraction of a core, which means raising it is a capacity decision as much as a security one.',
    'Social sign-in depends on somebody else being up. It sits on its own route with its own timeout, so a slow provider cannot consume the connection pool that ordinary password sign-in needs.',
  ],

  bottlenecks: {
    problems: [
      {
        label: 'Refresh storms',
        text: 'every client with a 15 minute token renews on roughly the same cadence. A deploy that restarts every client at once converts a steady 74 writes a second into a single spike.',
      },
      {
        label: 'Bcrypt under credential stuffing',
        text: 'an attacker sending ten thousand sign-ins a second is asking for ten thousand slow hashes a second. That is a denial of service aimed at the processor rather than the database, and it is caused by the very thing that makes passwords safe.',
      },
      {
        label: 'The sessions table growing',
        text: 'without a sweep, revoked and expired rows accumulate forever, and the unique index on token_hash grows along with them.',
      },
      {
        label: 'Secret rotation',
        text: 'changing JWT_SECRET invalidates every access token at once. With a single key there is no way to roll it without signing everybody out, which means in practice it never gets rotated.',
      },
      {
        label: 'Clock skew',
        text: 'token expiry is a comparison against the local clock. A replica whose clock has drifted will reject valid tokens or accept dead ones, and nothing in the logs will say so.',
      },
    ],
    improvements: [
      {
        label: 'Jitter the renewal',
        text: 'clients renew at a random point inside the last third of the token lifetime rather than at a fixed offset. The same number of renewals, spread across minutes instead of arriving together.',
      },
      {
        label: 'Cost the attack before hashing',
        text: 'the rate limiter and the account lockout both run before the password service, so a stuffing attack is refused without the server paying for a hash.',
      },
      {
        label: 'Sweep on a schedule',
        text: 'a scheduled job deletes sessions past their absolute expiry and revoked rows past the retention window, keeping the table proportional to live devices rather than to history.',
      },
      {
        label: 'Two signing keys',
        text: 'verify against both the current and the previous key, but sign only with the current one. Rotation then takes one access token lifetime instead of signing out every user, which is the difference between a secret that can be rotated and one that cannot.',
      },
      {
        label: 'A small leeway on expiry',
        text: 'allowing a few seconds of clock skew on the expiry comparison removes a whole class of failure that is otherwise invisible until a replica drifts.',
      },
      {
        label: 'Asymmetric signing for many services',
        text: 'with a shared secret, every service that verifies a token could also mint one. Signing with a private key and verifying with the public one means a compromised service can check tokens and not forge them.',
      },
    ],
  },

  seeAlso: [
    { title: 'Authentication reference', href: '/docs/backend/authentication' },
    { title: 'Sessions and devices', href: '/docs/systems/session-management' },
    { title: 'Authorization', href: '/docs/systems/authorization' },
  ],
}
