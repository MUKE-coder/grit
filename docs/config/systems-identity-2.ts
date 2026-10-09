import type { SystemDesign } from './systems-types'

/* Identity, continued: sessions, second factors, passkeys. */

export const SESSIONS: SystemDesign = {
  slug: 'session-management',
  name: 'Session Management',
  tagline:
    'Turning a stateless token into something you can revoke, list by device, and expire two different ways.',
  group: 'Identity',
  packages: ['internal/session', 'internal/models'],

  problem: {
    text: [
      'A signed token is self-contained by design: the server can verify it without looking anything up. That property is also the problem. A token that needs no lookup cannot be cancelled, so "sign out this device", "sign out everywhere", and "that laptop was stolen" have no implementation.',
      'The usual answer is a denylist of revoked tokens, which puts a lookup back in front of every request and grows without bound. The better answer is to keep the access token stateless and short, and to make the long-lived refresh token the thing that is tracked. A session row then represents a logged-in device, and revoking it means the next refresh fails.',
      'That leaves one subtle attack. If a refresh token is copied, both the real client and the attacker can use it. Whoever refreshes second presents a token that has already been rotated away, and that is a signal: the pair has been cloned.',
    ],
    capabilities: [
      'Record one row per logged-in device without ever storing a token that could be replayed from it.',
      'Rotate the refresh token on every use, so a captured one has a short window.',
      'Notice when a token that was already rotated is presented again, and treat the whole session as compromised.',
      'Expire a session two ways: after inactivity, and after an absolute age regardless of activity.',
      'Show a user their devices with enough detail to recognise one, and let them end any of them.',
      'End every session for a user when their password changes.',
    ],
  },

  functional: [
    'Create a session on sign-in, carrying the user agent and address the request arrived with.',
    'Rotate on refresh: write the new token hash, keep the previous one.',
    'Reject and revoke on replay: a presented token matching the previous hash kills the session.',
    'Revoke one session by id, every session for a user, or every session except the current one.',
    'List a user their live sessions, newest activity first, with the current one marked.',
    'Expire idle sessions after seven days and all sessions after thirty, both configurable.',
    'Sweep dead rows on a schedule so the table stays proportional to live devices.',
  ],

  nonFunctional: [
    {
      label: 'Unreplayable at rest',
      text: 'the table stores SHA-256 digests. Somebody with a dump can recognise a token you show them but cannot produce one.',
    },
    {
      label: 'One lookup',
      text: 'validating a refresh is a single hit on a unique index, not a scan. The cost does not grow with the number of sessions.',
    },
    {
      label: 'Bounded size',
      text: 'rows represent live devices, not history. An audit trail of sign-ins belongs in the audit log, which is append-only and has its own retention.',
    },
    {
      label: 'Recognisable',
      text: 'a device list is useless if every row says "Mozilla/5.0". It carries the parsed agent, the address and the last-seen time.',
    },
    {
      label: 'Fail closed',
      text: 'unknown, revoked, idle, expired and replayed all produce the same refusal. There is no state in which an ambiguous session is honoured.',
    },
  ],

  capacity: {
    assumptions: [
      ['Daily active users', '200,000'],
      ['Devices per active user', '3'],
      ['Refresh interval', '15 minutes, while a client is open'],
      ['Session row with indexes', '~650 bytes'],
      ['Idle timeout', '7 days'],
      ['Absolute timeout', '30 days'],
    ],
    estimates: [
      {
        label: 'Live sessions',
        working: ['200,000 DAU x 3 devices = 600,000 rows', '600,000 x 650 bytes = ~390 MB'],
        note: 'Comfortably inside the primary database. The number is driven by devices, not by time, which is the point of sweeping.',
      },
      {
        label: 'Rotation write rate',
        working: [
          '8 active hours / 15 minutes = 32 refreshes per device per day',
          '600,000 devices x 32 = 19,200,000 writes/day',
          '19,200,000 / 86,400 = 222 writes/second average',
          'peak at 5x = 1,110 writes/second',
        ],
        note: 'This is the single heaviest write in the authentication path. It is an update to one row by unique key, which Postgres handles at this rate on modest hardware, but it is the number that decides whether the access token lifetime can be shortened.',
      },
      {
        label: 'Sweep volume',
        working: [
          'sessions created per day = 240,000 sign-ins',
          'at steady state the same number expire per day',
          'a nightly sweep deletes ~240,000 rows',
        ],
        note: 'Small enough for one batched delete. Deleting in chunks rather than one statement keeps the lock short.',
      },
    ],
  },

  highLevel: {
    intro:
      'The session store sits between the token layer and the database. Nothing else writes to the table, which is what keeps the rotation invariant true.',
    components: [
      {
        label: 'Session store',
        text: 'creates, rotates, revokes and lists. The only writer, so the rule that a rotation always records the previous hash lives in one function.',
      },
      {
        label: 'Rotation check',
        text: 'given a presented refresh token, decides whether it belongs to a live session, a replayed one, or nothing at all.',
      },
      {
        label: 'Expiry policy',
        text: 'two clocks per row. last_seen_at drives idle expiry and moves on every refresh; expires_at is fixed at creation and does not.',
      },
      {
        label: 'Device list',
        text: 'the read side. Returns live sessions with the parsed agent and marks whichever row the current request is using.',
      },
      {
        label: 'Sweeper',
        text: 'a scheduled job that deletes rows past their absolute expiry and revoked rows past the retention window.',
      },
    ],
    flow: {
      title: 'Refresh, rotation, and what a replay looks like',
      nodes: [
        { id: 'client', label: 'Client', col: 0, row: 0 },
        { id: 'handler', label: 'Refresh Handler', col: 1, row: 0 },
        { id: 'store', label: 'Session Store', col: 2, row: 0, accent: true },
        { id: 'db', label: 'sessions table', col: 3, row: 0 },
        { id: 'jwt', label: 'JWT Service', col: 2, row: 1, accent: true },
        { id: 'replay', label: 'Replay Detected', col: 3, row: 1 },
        { id: 'revoke', label: 'Revoke Session', col: 3, row: 2 },
        { id: 'audit', label: 'Audit Log', col: 2, row: 2 },
      ],
      edges: [
        { from: 'client', to: 'handler', step: 1 },
        { from: 'handler', to: 'store', step: 2 },
        { from: 'store', to: 'db', step: 3 },
        { from: 'db', to: 'replay', step: 4, dashed: true },
        { from: 'replay', to: 'revoke', step: 5, dashed: true },
        { from: 'store', to: 'jwt', step: 6 },
        { from: 'jwt', to: 'audit', step: 7, dashed: true },
      ],
      steps: [
        'The client posts its refresh token, either in the body or from the HttpOnly cookie scoped to /api/auth.',
        'The handler hands the raw token to the session store. No other component sees it.',
        'The store hashes it and looks for a live session by token_hash: not revoked, not past either deadline.',
        'If nothing matches by token_hash, the store looks for a match on prev_token_hash. A hit there means this token was already rotated away, so two clients hold the pair.',
        'That is treated as theft, not as a stale client. The session is revoked outright, which signs out both the attacker and the real user, who then signs in again.',
        'On a clean match the store writes the new hash, moves the old one to prev_token_hash, and advances last_seen_at. The JWT service mints the new pair.',
        'The rotation is recorded outside the response path, and a detected replay is recorded as a security event rather than an ordinary one.',
      ],
    },
    dataFlow: [
      'The raw refresh token exists in memory for the length of one request and is never written anywhere.',
      'prev_token_hash is kept for exactly one generation. Keeping more would widen the window in which an old token is merely rejected rather than recognised as a replay.',
      'last_seen_at is an update on every refresh, which is why the row is designed to be narrow: the write is frequent.',
      'Revoked rows are kept rather than deleted so a device list can show recent sign-outs and the sweeper can age them out on its own schedule.',
    ],
  },

  stack: [
    ['Storage', 'the primary relational database'],
    ['Token digest', 'SHA-256, stored hex'],
    ['Primary key', 'UUIDv7, time-ordered'],
    ['Lookup index', 'unique on token_hash, secondary on prev_token_hash'],
    ['Cookie', 'HttpOnly, SameSite=Lax, Path=/api/auth'],
    ['Sweeper', 'a cron entry in the scheduling system'],
  ],

  dataModel: {
    entities: [
      {
        name: 'sessions',
        note: 'One row per logged-in device. Narrow on purpose: last_seen_at is written on every refresh.',
        fields: [
          ['id', 'UUIDv7'],
          ['user_id', 'owner, indexed for the device list'],
          ['token_hash', 'SHA-256 of the current refresh token, unique'],
          ['prev_token_hash', 'the previous generation, indexed, for replay detection'],
          ['user_agent', 'as sent, parsed for display'],
          ['ip', 'the address the session was last refreshed from'],
          ['created_at', 'when the device signed in'],
          ['last_seen_at', 'moved on every refresh; drives idle expiry'],
          ['expires_at', 'fixed at creation; the absolute deadline'],
          ['revoked_at', 'set on sign-out, not deleted'],
        ],
      },
    ],
    storage: [
      'One table, no cache in front of it. A refresh is already a write, so caching the read would save nothing and risk serving a revoked session.',
      'The unique index on token_hash is the hot path. The index on prev_token_hash is only touched when the first lookup misses, which is the rare case.',
      'Both timeouts are enforced in the query, not by a background job, so a session is dead the moment it qualifies rather than when the sweeper next runs.',
    ],
  },

  api: {
    groups: [
      {
        title: 'Sessions',
        rows: [
          { method: 'POST', path: '/api/v1/auth/refresh', what: 'Rotate and reissue' },
          { method: 'GET', path: '/api/v1/auth/sessions', what: 'Live devices, newest activity first' },
          { method: 'DELETE', path: '/api/v1/auth/sessions/:id', what: 'End one device' },
          { method: 'POST', path: '/api/v1/auth/sessions/revoke-all', what: 'End every device' },
          { method: 'POST', path: '/api/v1/auth/logout', what: 'End the current device only' },
        ],
      },
    ],
    samples: [
      {
        title: 'The device list',
        language: 'json',
        code: `{
  "data": [
    {
      "id": "01a11e83-a9b0-75e4-9f8b-0ea9507f6994",
      "user_agent": "Chrome 141 on Windows",
      "ip": "102.203.209.219",
      "created_at": "2026-10-02T08:14:09Z",
      "last_seen_at": "2026-10-09T17:02:30Z",
      "current": true
    },
    {
      "id": "01a0f2b1-5c3d-7a21-8e44-1b9cc2f40a77",
      "user_agent": "Safari on iPhone",
      "ip": "41.210.8.3",
      "last_seen_at": "2026-10-07T21:40:11Z",
      "current": false
    }
  ]
}`,
      },
    ],
  },

  lowLevel: {
    classes: [
      {
        name: 'SessionStore',
        file: 'internal/session/session.go',
        what: 'The only writer. Every mutation goes through it, which is what keeps the rule that a rotation records the previous hash from being forgotten at a new call site.',
        methods: ['CreateSession', 'RotateSession', 'Revoke', 'RevokeAllForUser', 'ListForUser', 'Sweep'],
      },
      {
        name: 'RequestMeta',
        file: 'internal/session/session.go',
        what: 'The three things a session needs from the request: the agent, the address and the context. Taking a value rather than a *gin.Context is what lets a background job create a session.',
      },
      {
        name: 'ErrSessionInvalid',
        what: 'One error for every way a session can fail to be usable. Callers cannot accidentally distinguish "revoked" from "expired" in a response and tell an attacker which it was.',
      },
    ],
    principles: [
      {
        label: 'Single writer',
        text: 'the invariant is "a rotation always records the previous hash". It holds because there is exactly one function that can rotate.',
      },
      {
        label: 'Interface segregation',
        text: 'the store takes a clock and a database handle, not a web framework. The expiry tests move the clock instead of sleeping.',
      },
      {
        label: 'Tell, do not ask',
        text: 'callers ask for a rotation and get a result. They do not read the row, decide it is valid, and write it back, which would be a race.',
      },
    ],
    patterns: [
      ['Token rotation', 'a new refresh token on every use, with one generation of memory'],
      ['Sliding plus absolute expiry', 'two clocks, so neither activity nor inactivity alone decides'],
      ['Tombstone', 'revoked rows persist until swept, so the device list can explain itself'],
    ],
  },

  scaling: [
    'The rotation write is a single-row update by unique key. It scales with the primary database and is the reason to think twice before shortening the access token lifetime.',
    'Reads for the device list are per-user and rare. They do not need an index beyond the one on user_id.',
    'The sweeper deletes in batches on a schedule rather than in one statement, so a large backlog does not hold a lock for minutes.',
    'Because validity is computed in the query, a replica lagging behind the primary can briefly honour a session that was just revoked. Revocation therefore reads and writes the primary.',
    'Nothing here is cached. A cache would turn revocation from immediate into eventually consistent, which is the one property this system exists to provide.',
  ],

  bottlenecks: {
    problems: [
      {
        label: 'Write amplification from short tokens',
        text: 'the rotation rate is inversely proportional to the access token lifetime. A 5 minute token triples the write rate on this table.',
      },
      {
        label: 'Replay false positives',
        text: 'a client that retries a refresh after a network timeout can legitimately present a token that already rotated, and gets its session killed for it.',
      },
      {
        label: 'Unbounded growth without a sweeper',
        text: 'if the scheduled sweep stops running, nothing else deletes rows, and the unique index grows with every sign-in ever made.',
      },
      {
        label: 'Device lists that nobody can read',
        text: 'a raw user agent string tells a user nothing, which makes the "end this session" button unusable in practice.',
      },
    ],
    improvements: [
      {
        label: 'Jittered refresh',
        text: 'clients refresh at a random point in the last third of the lifetime, which spreads the write rate instead of synchronising it.',
      },
      {
        label: 'A grace window on rotation',
        text: 'accepting the previous token for a few seconds after a rotation, once, distinguishes a retry from a replay without weakening the detection.',
      },
      {
        label: 'Sweep monitoring',
        text: 'the sweeper reports rows deleted. A sudden zero is the signal that it stopped, long before the table size is a problem.',
      },
      {
        label: 'Parse the agent on write',
        text: 'storing both the raw string and a readable summary means the list is useful without parsing on every read.',
      },
    ],
  },

  seeAlso: [
    { title: 'Authentication', href: '/docs/systems/authentication' },
    { title: 'Multi-factor authentication', href: '/docs/systems/multi-factor-auth' },
  ],
}

export const MFA: SystemDesign = {
  slug: 'multi-factor-auth',
  name: 'Multi-Factor Authentication',
  tagline:
    'A second factor that survives a leaked password, without locking out the user who loses their phone.',
  group: 'Identity',
  packages: ['internal/totp', 'internal/handlers', 'internal/services'],

  problem: {
    text: [
      'Passwords leak. They are reused across sites, phished, and dumped in breaches that have nothing to do with your application. A second factor makes a leaked password insufficient on its own.',
      'The design tension is not the cryptography, which is a published standard and a small amount of code. It is recovery. Any second factor strong enough to stop an attacker is strong enough to lock out a user who drops their phone in a river, and a recovery path weak enough to be convenient is the new weakest link.',
      'The second tension is where in the sign-in flow it goes. Issuing a full token and then asking for a code means the token existed before the second factor was checked. The flow has to stop halfway and issue something that is only good for finishing it.',
    ],
    capabilities: [
      'Enrol an authenticator app by a scannable secret, and confirm enrolment by requiring one working code before switching it on.',
      'Interrupt sign-in after the password is verified but before any usable token is issued.',
      'Accept a time-based code, tolerating a reasonable amount of clock drift but not an unlimited amount.',
      'Accept a single-use backup code when the authenticator is gone, and consume it.',
      'Remember a device the user trusts, for a bounded time, so the factor is not demanded hourly.',
      'Let a user see how many backup codes remain and regenerate the set.',
    ],
  },

  functional: [
    'Generate a secret and an otpauth:// URI, rendered as a QR code in the admin and the mobile app.',
    'Require one valid code before two-factor is marked enabled on the account.',
    'Return a pending token from sign-in when two-factor is enrolled, usable only to complete that sign-in.',
    'Verify a six-digit code against a thirty-second step, searching a bounded window either side.',
    'Issue ten single-use backup codes at enrolment, stored hashed, shown once.',
    'Optionally trust the current device for a period, skipping the prompt on it.',
    'Disable two-factor only after re-authenticating.',
  ],

  nonFunctional: [
    {
      label: 'No usable token before the second factor',
      text: 'the pending token carries one claim: that this user passed the password check. It opens nothing else.',
    },
    {
      label: 'Secrets encrypted at rest',
      text: 'a TOTP secret is a password equivalent. It is stored through the field encryption system, not in plaintext.',
    },
    {
      label: 'Bounded clock tolerance',
      text: 'a window of a few steps either side covers real drift. An unbounded window turns a captured code into a reusable one.',
    },
    {
      label: 'Backup codes are one-shot',
      text: 'each is consumed on use. A code that still works after being used is a password with extra steps.',
    },
    {
      label: 'Recoverable',
      text: 'a user who has lost every factor can be helped by an operator, and that intervention is in the audit log.',
    },
  ],

  capacity: {
    assumptions: [
      ['Users with two-factor enrolled', '15% of 1,000,000 = 150,000'],
      ['Sign-ins per enrolled user per day', '1.2'],
      ['Backup codes per user', '10'],
      ['Trusted device lifetime', '30 days'],
      ['TOTP step', '30 seconds, 6 digits'],
    ],
    estimates: [
      {
        label: 'Code verifications',
        working: [
          '150,000 x 1.2 = 180,000 challenges/day',
          'minus ~60% skipped on a trusted device = 72,000 verifications/day',
          '72,000 / 86,400 = 0.8/second average',
        ],
        note: 'Verification is an HMAC over a counter, repeated across a small window. The cost is irrelevant; the storage and the recovery path are what matter.',
      },
      {
        label: 'Backup code storage',
        working: [
          '150,000 users x 10 codes = 1,500,000 rows',
          'each row is a bcrypt hash plus a used stamp, ~120 bytes',
          '1,500,000 x 120 = ~180 MB',
        ],
        note: 'Codes are hashed with the same function as passwords, so verifying one costs a bcrypt comparison. Ten per user bounds the work when a code is presented.',
      },
      {
        label: 'Trusted devices',
        working: [
          '150,000 users x 2 trusted devices = 300,000 rows',
          'expiring after 30 days, so the table turns over monthly',
        ],
      },
    ],
  },

  highLevel: {
    intro:
      'Two-factor splits sign-in into two requests. The first proves the password and returns a token that can do nothing but finish; the second proves possession and exchanges it for a real pair.',
    components: [
      {
        label: 'TOTP service',
        text: 'generates secrets, builds the otpauth URI, and verifies a code against the current time step and a bounded window either side.',
      },
      {
        label: 'Pending token',
        text: 'a short-lived token with a single claim. It is accepted by exactly one endpoint and nothing else on the API will take it.',
      },
      {
        label: 'Backup code store',
        text: 'ten hashed codes per user, each with a used stamp. A used code is kept so the count shown to the user is honest.',
      },
      {
        label: 'Trusted device register',
        text: 'a cookie-bound record that lets a known device skip the prompt for a bounded period.',
      },
      {
        label: 'Two-factor service',
        text: 'the decisions: is it enrolled, is this device trusted, did this code work, how many codes are left.',
      },
    ],
    flow: {
      title: 'Sign-in interrupted by a second factor',
      nodes: [
        { id: 'client', label: 'Client', col: 0, row: 0 },
        { id: 'login', label: 'Login Handler', col: 1, row: 0 },
        { id: 'password', label: 'Password Service', col: 2, row: 0 },
        { id: 'twofa', label: 'Two-Factor Service', col: 3, row: 0, accent: true },
        { id: 'pending', label: 'Pending Token', col: 1, row: 1 },
        { id: 'verify', label: 'Verify Handler', col: 2, row: 1 },
        { id: 'totp', label: 'TOTP Service', col: 3, row: 1, accent: true },
        { id: 'backup', label: 'Backup Codes', col: 3, row: 2 },
        { id: 'tokens', label: 'Token Pair', col: 2, row: 2 },
      ],
      edges: [
        { from: 'client', to: 'login', step: 1 },
        { from: 'login', to: 'password', step: 2 },
        { from: 'password', to: 'twofa', step: 3 },
        { from: 'twofa', to: 'pending', step: 4, bend: 'v' },
        { from: 'pending', to: 'verify', step: 5 },
        { from: 'verify', to: 'totp', step: 6 },
        { from: 'totp', to: 'backup', step: 7, dashed: true },
        { from: 'verify', to: 'tokens', step: 8 },
      ],
      steps: [
        'The client posts email and password as usual. Nothing about the first request says two-factor is involved.',
        'The password service verifies the hash. A wrong password ends here, with no hint about whether a second factor would have been asked for.',
        'The two-factor service checks whether this account has it enrolled, and whether the request carries a trusted-device cookie that is still valid.',
        'If a factor is needed, sign-in stops. A pending token is issued: short-lived, single-purpose, and accepted by one endpoint. No access token exists at this point.',
        'The client prompts for a code and posts it with the pending token.',
        'The TOTP service recomputes the expected code for the current step and a bounded window either side, to absorb clock drift without accepting an old code indefinitely.',
        'If the submitted value is not a TOTP code, it is checked against the unused backup codes and consumed on a match.',
        'Only now is a real token pair issued and a session created. If the user asked to trust the device, the register gets a row and the response sets the cookie.',
      ],
    },
    dataFlow: [
      'The TOTP secret is written encrypted and decrypted only inside the verification call. It is never returned by an API after enrolment.',
      'Backup codes are shown exactly once, at generation. The server keeps bcrypt hashes and cannot display them again.',
      'The pending token is not a session. Nothing is written to the sessions table until the second factor passes.',
      'A trusted device is bound to a cookie value and an expiry, not to an address, because addresses change and users move.',
    ],
  },

  stack: [
    ['Algorithm', 'TOTP, RFC 6238, SHA-1, 6 digits, 30 second step'],
    ['Enrolment transport', 'otpauth:// URI rendered as a QR code'],
    ['Drift window', 'a bounded number of steps either side of now'],
    ['Backup codes', '10 per user, bcrypt hashed, single use'],
    ['Secret storage', 'the field encryption system, not plaintext'],
    ['Pending token', 'a JWT with a single claim and a short expiry'],
  ],

  dataModel: {
    entities: [
      {
        name: 'users (the two-factor columns)',
        fields: [
          ['totp_secret', 'encrypted at rest, null until enrolled'],
          ['totp_enabled_at', 'null until one working code confirmed enrolment'],
        ],
      },
      {
        name: 'backup_codes',
        fields: [
          ['user_id', 'owner'],
          ['code_hash', 'bcrypt of the code, never the code'],
          ['used_at', 'null until consumed; a used row is kept so the remaining count is honest'],
        ],
      },
      {
        name: 'trusted_devices',
        fields: [
          ['user_id', 'owner'],
          ['token_hash', 'SHA-256 of the cookie value'],
          ['user_agent, ip', 'what the user sees in the list'],
          ['expires_at', 'bounded; a trusted device is not trusted forever'],
        ],
      },
    ],
  },

  api: {
    groups: [
      {
        title: 'Enrolment',
        rows: [
          { method: 'GET', path: '/api/v1/auth/totp/status', what: 'Is it on, how many backup codes remain' },
          { method: 'POST', path: '/api/v1/auth/totp/setup', what: 'Generate a secret and the otpauth URI' },
          { method: 'POST', path: '/api/v1/auth/totp/enable', what: 'Confirm with one working code and switch it on' },
          { method: 'POST', path: '/api/v1/auth/totp/disable', what: 'Turn it off, after re-authenticating' },
        ],
      },
      {
        title: 'Signing in',
        rows: [
          { method: 'POST', path: '/api/v1/auth/totp/verify', what: 'Exchange a pending token and a code for a real pair' },
          { method: 'POST', path: '/api/v1/auth/totp/backup-codes', what: 'Regenerate the set, invalidating the old one' },
        ],
      },
      {
        title: 'Trusted devices',
        rows: [
          { method: 'GET', path: '/api/v1/auth/trusted-devices', what: 'Devices that skip the prompt' },
          { method: 'DELETE', path: '/api/v1/auth/trusted-devices/:id', what: 'Stop trusting one' },
        ],
      },
    ],
    samples: [
      {
        title: 'Completing a sign-in',
        language: 'json',
        code: `{
  "pending_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "code": "418265",
  "trust_device": true
}`,
      },
    ],
  },

  lowLevel: {
    classes: [
      {
        name: 'totp',
        file: 'internal/totp/totp.go',
        what: 'The standard, implemented once. Secret generation, the otpauth URI, and verification against a bounded window.',
        methods: ['GenerateSecret', 'ProvisioningURI', 'Validate', 'GenerateBackupCodes'],
      },
      {
        name: 'TwoFactorService',
        file: 'internal/services/two_factor.go',
        what: 'The decisions, lifted out of the handler so they can be tested without HTTP. Whether a factor is required, whether a code or a backup code matched, how many remain.',
        methods: ['Required', 'VerifyCode', 'ConsumeBackupCode', 'TrustDevice', 'RemainingCodes'],
      },
      {
        name: 'TOTPPendingToken',
        what: 'The half-finished sign-in, as a model rather than an implicit convention. Short-lived and accepted by one endpoint.',
      },
    ],
    principles: [
      {
        label: 'Separate the standard from the policy',
        text: 'the TOTP package knows RFC 6238 and nothing about users. The service knows the policy and nothing about HMAC.',
      },
      {
        label: 'Confirm before enabling',
        text: 'enrolment requires a working code. Without that step, a user who mis-scans the QR locks themselves out at the next sign-in.',
      },
      {
        label: 'Make recovery explicit',
        text: 'backup codes are generated at enrolment, not offered later as an afterthought, because the moment a user needs them is the moment they cannot reach the settings page.',
      },
    ],
    patterns: [
      ['Two-phase authentication', 'a pending token that is good for one thing'],
      ['One-time token', 'backup codes consumed on use and kept as tombstones'],
      ['Time window', 'a bounded search either side of the current step'],
    ],
  },

  scaling: [
    'Verification is pure computation over a handful of steps. It adds nothing measurable to the sign-in path.',
    'Backup code verification is a bcrypt comparison against up to ten hashes, so a wrong code costs ten comparisons. That bounds the work and is the reason the set is small.',
    'Trusted devices remove most of the prompts, which is a user-experience decision with a capacity side effect: the verification rate is a fraction of the sign-in rate.',
    'Nothing here is shared state beyond the database, so it scales with the API.',
    'The drift window is the one knob with a security cost. Widening it to help users with bad clocks lengthens the life of a captured code.',
  ],

  bottlenecks: {
    problems: [
      {
        label: 'Lockout with no recovery',
        text: 'a user who enrols, never saves the backup codes and loses the phone is locked out permanently unless an operator can intervene.',
      },
      {
        label: 'Code replay inside the window',
        text: 'a code is valid for its whole step and the drift window. An attacker who reads one over the user’s shoulder has a real, if short, opportunity.',
      },
      {
        label: 'Trusted devices as a soft spot',
        text: 'a trusted-device cookie is a bearer credential that skips the second factor. On a shared machine it defeats the point.',
      },
      {
        label: 'Phishable by design',
        text: 'TOTP is a shared secret. A convincing proxy page can collect the code and use it within the window, which passkeys solve and TOTP does not.',
      },
    ],
    improvements: [
      {
        label: 'Record used codes',
        text: 'remembering the last accepted step per user rejects the same code twice and closes the shoulder-surfing window.',
      },
      {
        label: 'Nudge the codes',
        text: 'showing the remaining count on the account page, and prompting when it reaches two, catches the lockout before it happens.',
      },
      {
        label: 'Bound and show trusted devices',
        text: 'a short expiry plus a visible list, revocable in one click, keeps the convenience without making it invisible.',
      },
      {
        label: 'Offer passkeys alongside',
        text: 'for the phishing case the answer is a factor that is bound to the origin. Both are supported, and a user can hold either or both.',
      },
    ],
  },

  seeAlso: [
    { title: 'Passkeys', href: '/docs/systems/passkeys' },
    { title: 'Authentication', href: '/docs/systems/authentication' },
  ],
}

export const PASSKEYS: SystemDesign = {
  slug: 'passkeys',
  name: 'Passkeys',
  tagline:
    'A credential that cannot be phished, cannot be reused across sites, and leaves nothing worth stealing in the database.',
  group: 'Identity',
  packages: ['internal/handlers', 'internal/models', 'go-webauthn/webauthn'],

  problem: {
    text: [
      'Every shared-secret factor has the same flaw. The user holds something the server also holds, or can recompute, and anything the user can be persuaded to type can be persuaded out of them by a page that looks right. Two-factor codes narrow the window; they do not close it, because a proxy that relays the code in real time still works.',
      'Public key authentication closes it. The authenticator keeps a private key that never leaves it, the server keeps only the public half, and the browser will only sign a challenge for the origin the credential was registered to. A fake page at a lookalike domain cannot get a usable signature, and the user has no secret to hand over even if they want to.',
      'What makes this a system rather than a library call is everything around the key: a challenge that must be used once, credentials that belong to a user who may have several, a counter that detects a cloned authenticator, and a path for the user whose only passkey is on a device they no longer have.',
    ],
    capabilities: [
      'Register an authenticator against a relying party derived from the application origin.',
      'Issue a challenge that is random, short-lived and single-use.',
      'Verify a signature against a stored public key, and reject one produced for a different origin.',
      'Hold several credentials per user, because people have a laptop and a phone.',
      'Notice a signature counter that goes backwards, which means the authenticator was cloned.',
      'Let a user name, list and remove their credentials, and keep a second factor available if they remove the last one.',
    ],
  },

  functional: [
    'Begin and finish registration, storing the credential id, public key, sign count and transports.',
    'Begin and finish authentication, with the user either named or discovered from the credential.',
    'Derive the relying party id from the configured origin, and refuse to start if that origin is not set correctly.',
    'Store challenges server-side with a short expiry and delete them on use.',
    'List a user their passkeys with a name, the device it was created on and the last use.',
    'Remove a passkey, and warn when it is the last one.',
  ],

  nonFunctional: [
    {
      label: 'Nothing worth stealing',
      text: 'the stored half is public. A full database dump yields no credential that can be used to sign in.',
    },
    {
      label: 'Origin bound',
      text: 'the browser refuses to sign for an origin the credential was not registered to, which is what makes the credential unphishable rather than merely strong.',
    },
    {
      label: 'Single-use challenges',
      text: 'a challenge is generated server-side, stored, and deleted when consumed. A replayed assertion finds nothing to match.',
    },
    {
      label: 'Clone detection',
      text: 'a sign counter that does not advance is evidence the authenticator was copied, and is surfaced rather than ignored.',
    },
    {
      label: 'Never the only door',
      text: 'a user who removes their last passkey still has a password and, if enrolled, a second factor. Authentication does not become unrecoverable.',
    },
  ],

  capacity: {
    assumptions: [
      ['Users with a passkey', '10% of 1,000,000 = 100,000'],
      ['Passkeys per user', '1.6'],
      ['Credential row', '~400 bytes, mostly the public key and transports'],
      ['Challenge lifetime', '2 minutes'],
      ['Sign-ins per passkey user per day', '1.2'],
    ],
    estimates: [
      {
        label: 'Credential storage',
        working: ['100,000 users x 1.6 = 160,000 credentials', '160,000 x 400 bytes = ~64 MB'],
        note: 'Trivial. The storage question for passkeys is never size, it is making sure the credential id is indexed, because discoverable sign-in looks up by it with no user named.',
      },
      {
        label: 'Challenge churn',
        working: [
          '100,000 x 1.2 = 120,000 sign-ins/day, plus abandoned attempts at ~20%',
          '~144,000 challenges/day, each alive for 2 minutes',
          '144,000 / 86,400 x 120 s = ~200 live at any moment',
        ],
        note: 'Two hundred rows. This is the clearest case in the whole framework for a short-lived store rather than a durable table.',
      },
      {
        label: 'Verification cost',
        working: [
          'one ECDSA or RSA signature verification per sign-in',
          '120,000 / 86,400 = 1.4/second average, 7/second at peak',
        ],
        note: 'Verification is roughly 100 microseconds. Compared with bcrypt at 60 milliseconds, passkey sign-in is cheaper for the server than password sign-in by two orders of magnitude.',
      },
    ],
  },

  highLevel: {
    intro:
      'Both ceremonies are two requests: the server issues a challenge, the authenticator signs it, the server verifies. The interesting state is the challenge, which must exist between the two.',
    components: [
      {
        label: 'WebAuthn service',
        text: 'wraps the protocol library. Builds the creation and request options, and verifies what comes back against the stored credential.',
      },
      {
        label: 'Relying party config',
        text: 'the id and origin, derived from the configured application origin. Getting this wrong is the single most common setup failure, so it is validated at startup.',
      },
      {
        label: 'Challenge store',
        text: 'short-lived server-side state tying a challenge to the user or the session that asked for it.',
      },
      {
        label: 'Credential store',
        text: 'the public halves, one row per authenticator, indexed by credential id for discoverable sign-in.',
      },
      {
        label: 'Counter check',
        text: 'compares the sign count in the assertion against the stored one and flags a non-advance.',
      },
    ],
    flow: {
      title: 'Signing in with a passkey',
      nodes: [
        { id: 'client', label: 'Browser', col: 0, row: 0 },
        { id: 'begin', label: 'Begin Handler', col: 1, row: 0 },
        { id: 'webauthn', label: 'WebAuthn Service', col: 2, row: 0, accent: true },
        { id: 'challenge', label: 'Challenge Store', col: 3, row: 0 },
        { id: 'auth', label: 'Authenticator', col: 0, row: 1 },
        { id: 'finish', label: 'Finish Handler', col: 1, row: 1 },
        { id: 'creds', label: 'Credential Store', col: 2, row: 1, accent: true },
        { id: 'counter', label: 'Counter Check', col: 3, row: 1 },
        { id: 'tokens', label: 'Token Pair', col: 2, row: 2 },
      ],
      edges: [
        { from: 'client', to: 'begin', step: 1 },
        { from: 'begin', to: 'webauthn', step: 2 },
        { from: 'webauthn', to: 'challenge', step: 3 },
        { from: 'client', to: 'auth', step: 4 },
        { from: 'auth', to: 'finish', step: 5 },
        { from: 'finish', to: 'creds', step: 6 },
        { from: 'creds', to: 'counter', step: 7 },
        { from: 'creds', to: 'tokens', step: 8 },
      ],
      steps: [
        'The browser asks to begin. It may name a user, or ask for a discoverable credential and let the authenticator decide.',
        'The handler asks the WebAuthn service for request options: the relying party id, the allowed credentials if a user was named, and a fresh random challenge.',
        'The challenge is stored server-side with a short expiry. It is the only thing that makes the second request verifiable, and it is good once.',
        'The browser passes the options to the authenticator, which asks the user for a fingerprint, a face or a PIN, and signs the challenge with the private key. That key does not leave the device.',
        'The assertion comes back: the credential id, the signature, the authenticator data and the client data.',
        'The handler looks up the stored public key by credential id and asks the service to verify. The verification fails if the origin in the client data is not this application, which is what defeats a lookalike page.',
        'The sign count in the assertion is compared with the stored one. A count that did not advance suggests a cloned authenticator and is recorded as a security event.',
        'On success the challenge is deleted, the stored count and last-used stamp are updated, and the ordinary token pair and session are issued. From here it is the same sign-in as any other.',
      ],
    },
    dataFlow: [
      'The private key never reaches the server, and there is no API that could ask for it.',
      'The challenge exists between the two requests and is deleted on use, so an intercepted assertion cannot be replayed.',
      'The relying party id is derived from the configured origin. It is checked at startup rather than at first use, so a misconfiguration fails on deploy instead of on a user.',
      'The credential id is the lookup key for discoverable sign-in, where no username is sent at all.',
    ],
  },

  stack: [
    ['Protocol', 'WebAuthn level 2, via go-webauthn'],
    ['Credential storage', 'the primary database, public key only'],
    ['Challenge storage', 'short-lived server-side state with an expiry'],
    ['Relying party', 'derived from the configured application origin'],
    ['Client API', 'navigator.credentials, with a conditional-UI autofill path'],
  ],

  dataModel: {
    entities: [
      {
        name: 'passkeys',
        fields: [
          ['id', 'UUIDv7'],
          ['user_id', 'owner'],
          ['credential_id', 'the authenticator’s id, unique and indexed for discoverable sign-in'],
          ['public_key', 'the public half, as returned at registration'],
          ['sign_count', 'the last counter seen; a non-advance is suspicious'],
          ['transports', 'usb, nfc, ble, internal, hybrid; used to prompt sensibly'],
          ['name', 'what the user calls it'],
          ['created_at, last_used_at', 'for the management list'],
        ],
      },
      {
        name: 'webauthn_challenges',
        note: 'Deliberately short-lived. Roughly two hundred rows exist at any moment on the traffic above.',
        fields: [
          ['challenge', 'random bytes, the thing to be signed'],
          ['user_id', 'null for a discoverable sign-in'],
          ['kind', 'registration or authentication'],
          ['expires_at', 'minutes, not hours'],
        ],
      },
    ],
  },

  api: {
    groups: [
      {
        title: 'Registration',
        rows: [
          { method: 'POST', path: '/api/v1/auth/passkeys/register/begin', what: 'Creation options and a challenge' },
          { method: 'POST', path: '/api/v1/auth/passkeys/register/finish', what: 'Verify and store the credential' },
        ],
      },
      {
        title: 'Sign-in',
        rows: [
          { method: 'POST', path: '/api/v1/auth/passkeys/login/begin', what: 'Request options and a challenge' },
          { method: 'POST', path: '/api/v1/auth/passkeys/login/finish', what: 'Verify the assertion, issue tokens' },
        ],
      },
      {
        title: 'Management',
        rows: [
          { method: 'GET', path: '/api/v1/auth/passkeys', what: 'The user’s credentials' },
          { method: 'PATCH', path: '/api/v1/auth/passkeys/:id', what: 'Rename one' },
          { method: 'DELETE', path: '/api/v1/auth/passkeys/:id', what: 'Remove one' },
        ],
      },
    ],
  },

  lowLevel: {
    classes: [
      {
        name: 'WebAuthnService',
        what: 'The protocol, wrapped once. Builds options and verifies responses, and is the only place the library’s types appear.',
        methods: ['BeginRegistration', 'FinishRegistration', 'BeginLogin', 'FinishLogin'],
      },
      {
        name: 'PasskeyHandler',
        file: 'internal/handlers/passkey.go',
        what: 'Four endpoints, each thin. The interesting rule it enforces is that a finish request must match a challenge this server issued.',
      },
      {
        name: 'Passkey',
        file: 'internal/models/passkey.go',
        what: 'The stored credential. Implements the library’s credential interface so the service does not need a translation layer.',
      },
    ],
    principles: [
      {
        label: 'Do not reimplement the standard',
        text: 'the protocol is subtle and the failure mode is silent. A reviewed library handles it; the application handles the storage and the policy.',
      },
      {
        label: 'Validate configuration at startup',
        text: 'a wrong relying party id produces a browser error with no server-side trace. Checking the origin on boot turns an invisible failure into a loud one.',
      },
      {
        label: 'Keep the fallback',
        text: 'passkeys are an addition, not a replacement. Removing the last one does not remove the ability to sign in.',
      },
    ],
    patterns: [
      ['Challenge and response', 'server-issued nonce, signed by the holder of the key'],
      ['Two-phase ceremony', 'begin stores state, finish consumes it'],
      ['Adapter', 'the stored model implements the library’s credential interface'],
    ],
  },

  scaling: [
    'Signature verification is around a hundred microseconds, which makes passkey sign-in cheaper for the server than a bcrypt comparison by a wide margin.',
    'The credential table is read by a unique index on credential id. It does not grow with traffic, only with users.',
    'Challenges are the only churn, and the live set is tiny. They are a natural fit for a store with expiry rather than a table that needs sweeping.',
    'Nothing here is shared between replicas except the challenge, so a load balancer does not need sticky sessions provided the challenge store is shared.',
    'The relying party id is tied to the domain. Serving the same application on a second domain means a second set of credentials, which is a product decision, not a scaling one.',
  ],

  bottlenecks: {
    problems: [
      {
        label: 'Origin misconfiguration',
        text: 'the most common failure by a wide margin. A wrong origin produces a browser-side error and nothing in the server log.',
      },
      {
        label: 'Challenge store as a single point',
        text: 'if challenges live only in one replica’s memory, a load balancer that sends the finish request elsewhere breaks every sign-in.',
      },
      {
        label: 'The last-passkey problem',
        text: 'a user who removes their only credential from their only device, with no password set, has no way back in.',
      },
      {
        label: 'Counter false positives',
        text: 'some authenticators do not implement the counter and always report zero, so a naive non-advance check flags every one of them.',
      },
    ],
    improvements: [
      {
        label: 'Check the origin on boot',
        text: 'refusing to start, or warning loudly, when the configured origin cannot produce a valid relying party id moves the failure from a user to a deploy.',
      },
      {
        label: 'Share the challenge store',
        text: 'putting challenges where every replica can see them removes the sticky-session requirement, which is otherwise an invisible constraint.',
      },
      {
        label: 'Warn on the last credential',
        text: 'the management page says what will happen before the last passkey is removed, and points at the password and second-factor settings.',
      },
      {
        label: 'Treat a zero counter as absent',
        text: 'an authenticator that always reports zero is not a clone. The check applies only where the counter has ever advanced.',
      },
    ],
  },

  seeAlso: [
    { title: 'Passkeys reference', href: '/docs/backend/passkeys' },
    { title: 'Multi-factor authentication', href: '/docs/systems/multi-factor-auth' },
  ],
}
