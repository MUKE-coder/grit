import type { SystemDesign } from './systems-types'

/* Identity: who the caller is. Authentication, sessions, MFA, passkeys. */

export const AUTHENTICATION: SystemDesign = {
  slug: 'authentication',
  name: 'Authentication',
  tagline:
    'Proving who the caller is, on every request, without asking them for a password more than once.',
  group: 'Identity',
  packages: ['internal/handlers', 'internal/password', 'internal/middleware', 'internal/crypto'],

  problem: {
    text: [
      'Every endpoint past the public ones needs to know which user is calling. The naive answer is to send the password with each request, which means storing it somewhere on the client and handing it to the server a thousand times a day. The second naive answer is a session cookie looked up in a database on every request, which puts a database read in front of everything the application does.',
      'What is wanted instead is a short-lived credential the server can verify without a lookup, backed by a long-lived one that can be revoked. That is the split between an access token and a refresh token, and almost every decision below falls out of it.',
      'The hard parts are not the signing. They are what happens when a token is stolen, what happens when a user signs out on one device, and what happens when an attacker tries ten thousand passwords against one account from a thousand addresses.',
    ],
    capabilities: [
      'Register a user, hash the password with a function that is slow on purpose, and never store the plaintext.',
      'Issue a credential a request can carry, and verify it without reading the database.',
      'Renew that credential without a second sign-in, and notice when a renewal is a replay.',
      'Revoke one device without signing the user out everywhere, and revoke everywhere when a password changes.',
      'Resist a password-guessing attack spread across many IP addresses, without locking out a user whose address an attacker can guess.',
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
      text: 'checking an access token must not touch the database. It is an HMAC verification and a clock comparison, which is what makes it affordable on every request.',
    },
    {
      label: 'Blast radius',
      text: 'a stolen access token expires in 15 minutes. A stolen refresh token is usable once, and using it after the real client already has alerts the server.',
    },
    {
      label: 'Storage safety',
      text: 'a dump of the database must not let anyone sign in. Passwords are bcrypt hashes and refresh tokens are stored as SHA-256 digests, so neither column is replayable.',
    },
    {
      label: 'Availability',
      text: 'sign-in may depend on the database. Verifying a token may not, so an application stays readable through a database blip for as long as its access tokens live.',
    },
    {
      label: 'Honest errors',
      text: 'a wrong password and an unknown address answer the same way, in the same time, so the endpoint is not a membership oracle.',
    },
  ],

  capacity: {
    assumptions: [
      ['Registered users', '1,000,000'],
      ['Daily active users', '200,000'],
      ['Sign-ins per active user per day', '1.2'],
      ['API requests per active user per day', '400'],
      ['Access token lifetime', '15 minutes'],
      ['Refresh token lifetime', '168 hours (7 days)'],
      ['Session row', '~400 bytes'],
      ['User row', '~1.2 KB with profile fields'],
    ],
    estimates: [
      {
        label: 'Sign-in rate',
        working: [
          '200,000 DAU x 1.2 sign-ins = 240,000 sign-ins/day',
          '240,000 / 86,400 s = 2.8 sign-ins/second average',
          'peak at 5x average = 14 sign-ins/second',
        ],
        note: 'Bcrypt at cost 10 takes roughly 60 ms of CPU. Fourteen a second is under one core of hashing, which is why the cost factor can stay where it is.',
      },
      {
        label: 'Token verification rate',
        working: [
          '200,000 DAU x 400 requests = 80,000,000 requests/day',
          '80,000,000 / 86,400 s = 926 verifications/second average',
          'peak at 5x = 4,630 verifications/second',
        ],
        note: 'An HMAC-SHA256 verification is a few microseconds. At 4,630 a second this is noise, which is the entire argument for stateless access tokens.',
      },
      {
        label: 'Refresh rate',
        working: [
          'a 15 minute token over an 8 hour working session = 32 refreshes',
          '200,000 DAU x 32 = 6,400,000 refreshes/day',
          '6,400,000 / 86,400 = 74 refreshes/second average, 370 at peak',
        ],
        note: 'Every refresh is a database write, because rotation replaces the stored hash. This is the one part of the hot path that is not stateless, and it is the number to watch.',
      },
      {
        label: 'Session storage',
        working: [
          '200,000 DAU x 3 devices = 600,000 live sessions',
          '600,000 x 400 bytes = 240 MB',
          'plus indexes on token_hash, user_id, expires_at = ~400 MB',
        ],
        note: 'Small enough to stay in the primary database. Expired rows are swept on a schedule rather than left to accumulate.',
      },
      {
        label: 'User storage',
        working: ['1,000,000 users x 1.2 KB = 1.2 GB', 'with indexes and audit columns = ~2 GB'],
      },
    ],
  },

  highLevel: {
    intro:
      'Sign-in is a write path that ends in two tokens. Everything after it is a read path that touches no storage at all until the access token expires.',
    components: [
      {
        label: 'Auth handlers',
        text: 'the HTTP surface. Register, login, refresh, logout, forgot and reset. Thin: they bind, call a service, and shape the envelope.',
      },
      {
        label: 'Password service',
        text: 'bcrypt hashing and comparison, plus the strength rules. The only place that sees a plaintext password, and it never returns one.',
      },
      {
        label: 'JWT service',
        text: 'signs and verifies. Every token carries a user id, a role, a type, an expiry and a unique jti, so two tokens issued in the same second are still different tokens.',
      },
      {
        label: 'Session store',
        text: 'one row per logged-in device, keyed by the SHA-256 of the current refresh token. This is what makes a stateless token revocable.',
      },
      {
        label: 'Auth middleware',
        text: 'reads the Authorization header or the access cookie, verifies, and puts an identity on the request context. Everything downstream reads the context.',
      },
      {
        label: 'Lockout counter',
        text: 'failed attempts per account, not per address, so an attack spread across a botnet still trips it.',
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
        'The client posts credentials to /api/v1/auth/login. The middleware lets it through: this route is public.',
        'The handler binds the request and asks the lockout counter whether this account is currently barred. A barred account is refused here, before any hashing is paid for.',
        'The password service compares the submitted password against the stored bcrypt hash. A wrong password and an unknown address take the same path and the same time.',
        'A failed comparison increments the counter for that account. Ten failures inside the window bar it for fifteen minutes. An unknown address increments nothing, or an attacker could lock out any address they can guess.',
        'On success the JWT service mints an access token (15 minutes) and a refresh token (7 days), each with its own jti.',
        'The session store records the device: the SHA-256 of the refresh token, the user agent, the address, and the two expiry stamps.',
        'The session row is written. The refresh token itself is never stored, so this table cannot be replayed.',
        'The response carries both tokens, and sets the refresh token as an HttpOnly cookie scoped to /api/auth so no other route can read it.',
        'The sign-in is recorded in the audit log, outside the response path. A later request carries the access token and the middleware verifies it with no storage read at all.',
      ],
    },
    dataFlow: [
      'Passwords travel over TLS, are hashed on arrival, and are never written to a log, a trace or an error message.',
      'Access tokens are not stored anywhere on the server. Their validity is entirely a property of the signature and the clock.',
      'Refresh tokens are stored only as SHA-256 digests. The server can recognise a token it issued without being able to produce one.',
      'The previous refresh token hash is kept alongside the current one, which is what makes replay detection possible.',
      'The identity on the request context is a value, not a database handle. A service that needs the user id does not get the ability to query as them.',
    ],
  },

  stack: [
    ['Token format', 'JWT, HS256, signed with JWT_SECRET'],
    ['Password hashing', 'bcrypt, cost 10'],
    ['Session storage', 'the primary database, one row per device'],
    ['Transport', 'Authorization: Bearer, plus HttpOnly cookies for browsers'],
    ['Social sign-in', 'OAuth2 against Google and GitHub'],
    ['Lockout counter', 'database-backed, so it survives a restart'],
    ['Reset tokens', 'single-use, hashed at rest, short expiry'],
  ],

  dataModel: {
    intro:
      'Three tables carry the whole system. A user, a session per device, and a short-lived token for the reset flow.',
    entities: [
      {
        name: 'users',
        fields: [
          ['id', 'UUIDv7, time-ordered so the index stays dense'],
          ['email', 'unique, the sign-in identifier'],
          ['password_hash', 'bcrypt output, never the password'],
          ['role', 'ADMIN, EDITOR or USER'],
          ['email_verified_at', 'null until confirmed; gates sign-in when verification is required'],
          ['totp_secret', 'encrypted at rest, null until two-factor is enrolled'],
          ['created_at, updated_at', 'audit columns'],
        ],
      },
      {
        name: 'sessions',
        note: 'One row per logged-in device. This table is what turns a stateless token into a revocable one.',
        fields: [
          ['id', 'UUIDv7'],
          ['user_id', 'the owner, indexed'],
          ['token_hash', 'SHA-256 of the current refresh token, unique'],
          ['prev_token_hash', 'SHA-256 of the token this one replaced, for replay detection'],
          ['user_agent, ip', 'what the device list shows a user'],
          ['last_seen_at', 'moved on every refresh; drives the idle timeout'],
          ['expires_at', 'the absolute deadline, regardless of activity'],
          ['revoked_at', 'set on sign-out; a revoked row is kept, not deleted'],
        ],
      },
      {
        name: 'password_reset_tokens',
        fields: [
          ['token_hash', 'the emailed token, hashed'],
          ['user_id', 'who it is for'],
          ['expires_at', 'short, measured in minutes'],
          ['used_at', 'set on first use, so the link works once'],
        ],
      },
    ],
    storage: [
      'Everything is in the primary relational database. There is no second store to keep consistent, and a session read is a single indexed lookup.',
      'Sessions are swept on a schedule: rows past their absolute expiry are deleted, rows revoked more than a retention window ago go with them.',
      'The idle timeout is seven days and the absolute timeout is thirty. An unused login dies in a week; an active one still dies in a month.',
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
          { method: 'POST', path: '/api/v1/auth/forgot-password', what: 'Email a single-use reset link' },
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
        title: 'Sign-in response',
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
        title: 'When two-factor is enrolled, sign-in stops halfway',
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
        what: 'Binds and validates the request, calls the services, shapes the envelope. Holds no logic worth testing through HTTP.',
        methods: ['Register', 'Login', 'Refresh', 'Logout', 'Me', 'ForgotPassword', 'ResetPassword'],
      },
      {
        name: 'JWTService',
        file: 'internal/services/jwt.go',
        what: 'The only thing that signs or verifies. Issuing and verifying are separate methods so a verifier can exist without a signing key.',
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
        what: 'Hash, compare, and judge strength. Comparison takes the same time whether or not the user exists.',
        methods: ['Hash', 'Compare', 'Strength'],
      },
      {
        name: 'AuthMiddleware',
        file: 'internal/middleware/auth.go',
        what: 'Turns a token into an identity on the request context, or a 401. Does not decide what that identity may do.',
        methods: ['RequireAuth', 'OptionalAuth'],
      },
    ],
    principles: [
      {
        label: 'Single responsibility',
        text: 'the handler does HTTP, the JWT service does cryptography, the session store does revocation. Nothing knows two of those three.',
      },
      {
        label: 'Dependency inversion',
        text: 'services take interfaces for the clock and the store, which is what lets the expiry tests run without sleeping.',
      },
      {
        label: 'Open for extension',
        text: 'social sign-in and passkeys both end at the same point, issuing a token pair for an existing user record, so neither needed the token layer changed.',
      },
      {
        label: 'Fail closed',
        text: 'an unparseable token, an unknown session and a revoked session all produce the same 401. There is no path where an error leaves the request authenticated.',
      },
    ],
    patterns: [
      ['Strategy', 'password, social and passkey sign-in, each ending at the same token issuance'],
      ['Chain of responsibility', 'the middleware stack: request id, then auth, then rate limit, then the handler'],
      ['Token rotation', 'every refresh replaces the stored hash and remembers the old one'],
      ['Context carrier', 'the identity rides on context.Context, so a background job started by a request keeps it'],
    ],
  },

  scaling: [
    'Verification is stateless, so the API scales horizontally with no shared session store and no sticky sessions. Adding a replica adds capacity immediately.',
    'The refresh endpoint is the only part of the hot path that writes. At 370 writes a second it is comfortable on one primary, and it is the first thing to watch if the access token lifetime is shortened.',
    'Shortening the access token lifetime trades verification cost for refresh cost: halving it to 7 minutes doubles the write rate on that one endpoint.',
    'Session reads are a single lookup on a unique index. The table stays small because it holds live devices, not history.',
    'Bcrypt cost is a deliberate ceiling on sign-in throughput. It is sized so a peak hour fits in a fraction of a core, and raising it is a capacity decision, not just a security one.',
    'Social sign-in depends on a third party. It is on its own route with its own timeout, so a slow provider cannot consume the connection pool that password sign-in needs.',
  ],

  bottlenecks: {
    problems: [
      {
        label: 'Refresh storms',
        text: 'every client with a 15 minute token refreshes on roughly the same cadence. A deploy that restarts every client at once turns a steady 74 writes a second into a spike.',
      },
      {
        label: 'Bcrypt under credential stuffing',
        text: 'an attacker who sends ten thousand sign-ins a second is asking for ten thousand bcrypt hashes a second, which is a denial of service aimed at the CPU rather than the database.',
      },
      {
        label: 'The sessions table growing',
        text: 'without a sweep, revoked and expired rows accumulate forever, and the unique index on token_hash grows with them.',
      },
      {
        label: 'Secret rotation',
        text: 'changing JWT_SECRET invalidates every access token at once. With a single key there is no way to roll it without signing everybody out.',
      },
      {
        label: 'Clock skew',
        text: 'token expiry is a comparison against the local clock. A replica whose clock has drifted rejects valid tokens or accepts dead ones.',
      },
    ],
    improvements: [
      {
        label: 'Jitter the refresh',
        text: 'clients refresh at a random point inside the last third of the token lifetime rather than at a fixed offset, which spreads the spike across minutes.',
      },
      {
        label: 'Cost the attack before hashing',
        text: 'the rate limiter and the account lockout both run before the password service, so a stuffing attack is rejected without paying for a hash.',
      },
      {
        label: 'Sweep on a schedule',
        text: 'a cron job deletes sessions past their absolute expiry and revoked rows past the retention window, keeping the table proportional to live devices.',
      },
      {
        label: 'Two signing keys',
        text: 'verify against both the current and the previous key, sign with the current one. Rotation then takes one access-token lifetime instead of signing everybody out.',
      },
      {
        label: 'A small leeway on expiry',
        text: 'allowing a few seconds of skew on the expiry comparison removes a class of failure that is otherwise invisible until a replica drifts.',
      },
    ],
  },

  seeAlso: [
    { title: 'Authentication reference', href: '/docs/backend/authentication' },
    { title: 'Sessions and devices', href: '/docs/systems/session-management' },
    { title: 'Authorization', href: '/docs/systems/authorization' },
  ],
}
