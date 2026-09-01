/*
 * The capability comparison that sits under the benchmark on the homepage.
 *
 * The benchmark answers "how fast". This answers "how much do you assemble
 * yourself", which is the question Grit actually exists to change.
 *
 * ── Why four states and not a tick ────────────────────────────────────────
 *
 * A tick-or-cross matrix would be dishonest, and obviously so to anyone who
 * knows these frameworks. Django's admin is one of the best in software and it
 * ships in the box. Laravel has first-party queues, first-party observability
 * and first-party auth scaffolding. Encore generates typed clients and has a
 * tracing dashboard. Marking those as absent to make a column look good would
 * be the kind of comparison nobody should trust, including ours.
 *
 * So every cell records how the capability arrives, not whether it exists:
 *
 *   generated  the scaffolder or `grit generate` emits it, wired, with a UI.
 *              Nothing to install, nothing to configure to see it work.
 *   builtin    ships in the framework core. You use it, you do not install it.
 *   official   a first-party package you add and wire yourself.
 *   community  a third-party package, often excellent, but chosen and
 *              maintained by you.
 *   none       out of scope for that framework.
 *
 * Read honestly, Grit's case is not "we have things others lack". Most of these
 * rows are achievable everywhere. The claim is narrower and easier to check:
 * how many of them are already running the first time you open the app.
 *
 * Every cell is a claim about someone else's software, so anything non-obvious
 * carries a `note` naming the package involved. If a note is wrong, the cell is
 * wrong, and it should be fixed rather than defended.
 */

export type Provision = 'generated' | 'builtin' | 'official' | 'community' | 'none'

export interface Cell {
  state: Provision
  /** the package or subsystem involved, shown on hover and to screen readers */
  note?: string
}

export interface Capability {
  id: string
  label: string
  /** what the row actually means, so a reader is not guessing */
  detail: string
  cells: Record<string, Cell>
}

export interface CapabilityGroup {
  title: string
  blurb: string
  rows: Capability[]
}

/** Column order. Matches the benchmark so the two sections read as one story. */
export const COMPARED = [
  { slug: 'grit', name: 'Grit', logo: '/logos/grit.png' },
  { slug: 'laravel', name: 'Laravel', logo: '/logos/laravel.png' },
  { slug: 'django', name: 'Django', logo: '/logos/django.svg' },
  { slug: 'nextjs', name: 'Next.js', logo: '/logos/nextjs.png', invertOnDark: true },
  { slug: 'encore', name: 'Encore.ts', logo: '/logos/encore.png' },
  { slug: 'express', name: 'Express', logo: '/logos/express.png' },
] as const

const G = (note?: string): Cell => ({ state: 'generated', note })
const B = (note?: string): Cell => ({ state: 'builtin', note })
const O = (note?: string): Cell => ({ state: 'official', note })
const C = (note?: string): Cell => ({ state: 'community', note })
const N = (note?: string): Cell => ({ state: 'none', note })

export const CAPABILITY_GROUPS: CapabilityGroup[] = [
  {
    title: 'Identity',
    blurb:
      'The part every product rebuilds and nobody wants to. Grit generates the handlers, the ' +
      'middleware, the database tables and the screens together, so the first run already has a ' +
      'working sign-in.',
    rows: [
      {
        id: 'auth',
        label: 'Email and password auth',
        detail: 'Register, sign in, refresh, forgot and reset password, with the screens.',
        cells: {
          grit: G('Scaffolded handlers, JWT with rotating refresh tokens, and the pages'),
          laravel: O('Breeze or Fortify, installed separately'),
          django: B('django.contrib.auth, with views you template yourself'),
          nextjs: C('Auth.js, formerly NextAuth'),
          encore: C('Auth handlers are a hook you implement'),
          express: C('Passport or a hand-rolled equivalent'),
        },
      },
      {
        id: 'rbac',
        label: 'Roles and permissions',
        detail:
          'Roles as database rows rather than an enum, per-resource permissions registered by ' +
          'the generator, and a UI to edit them without a deploy.',
        cells: {
          grit: G('Every generated resource registers its own permissions'),
          laravel: C('spatie/laravel-permission is the de facto choice'),
          django: B('Groups and permissions, surfaced in the admin'),
          nextjs: C(),
          encore: C(),
          express: C(),
        },
      },
      {
        id: 'twofa',
        label: 'Two-factor with backup codes',
        detail: 'TOTP enrolment, a QR code, ten hashed single-use recovery codes.',
        cells: {
          grit: G('In the generated admin, not a wiring exercise'),
          laravel: O('Fortify'),
          django: C('django-otp'),
          nextjs: C(),
          encore: N(),
          express: C(),
        },
      },
      {
        id: 'passkeys',
        label: 'Passkeys',
        detail:
          'WebAuthn sign-in with a fingerprint, face or device PIN. The private key never ' +
          'leaves the authenticator, so there is nothing for a breach to leak and nothing for ' +
          'a phishing page to collect.',
        cells: {
          grit: B('Built in, pure Go, with the management UI'),
          laravel: O('Laravel Passkeys, or a community package'),
          django: C('django-passkeys or similar'),
          nextjs: C('SimpleWebAuthn, wired yourself'),
          encore: N(),
          express: C('SimpleWebAuthn, wired yourself'),
        },
      },
      {
        id: 'sso',
        label: 'Enterprise SSO (OIDC and SAML)',
        detail: 'Connections stored as rows, so onboarding a customer is a form and not a deploy.',
        cells: {
          grit: G('OIDC and SAML 2.0, one connection per customer'),
          laravel: C('Socialite covers OAuth; SAML is third party'),
          django: C('python-social-auth or djangosaml2'),
          nextjs: C('Auth.js covers OIDC; SAML is separate'),
          encore: N(),
          express: C(),
        },
      },
      {
        id: 'sessions',
        label: 'Server-side sessions with device revoke',
        detail:
          'Every refresh token backed by a row, replay detection, and a screen listing the ' +
          'devices signed in.',
        cells: {
          grit: G('Rotation with replay detection, idle and absolute timeouts'),
          laravel: B('Sessions are core; per-device revoke is yours to build'),
          django: B('Sessions are core; per-device revoke is yours to build'),
          nextjs: C(),
          encore: N(),
          express: C(),
        },
      },
      {
        id: 'recovery',
        label: 'Account recovery',
        detail:
          'A verified second address that gets somebody back in when the primary is gone, ' +
          'with the code hashed, single use, expiring, and capped at five guesses.',
        cells: {
          grit: B('Built in, and every change to it takes the account password'),
          laravel: N('Password reset only, to the address you may have lost'),
          django: N('Password reset only'),
          nextjs: C('Whatever your auth provider offers'),
          encore: N(),
          express: N(),
        },
      },
    ],
  },
  {
    title: 'Security and compliance',
    blurb:
      'The rows most frameworks leave to you, because they are product decisions rather than ' +
      'framework ones. Grit takes a position and ships a screen.',
    rows: [
      {
        id: 'sentinel',
        label: 'Security dashboard',
        detail:
          'Blocked requests, suspicious agents, rate-limit trips and lockouts, visible in the ' +
          'app rather than in a log aggregator you have to buy.',
        cells: {
          grit: G('Sentinel, in the generated admin'),
          laravel: N(),
          django: N(),
          nextjs: N(),
          encore: N(),
          express: N(),
        },
      },
      {
        id: 'audit',
        label: 'Tamper-evident audit log',
        detail:
          'A hash-chained record of every mutation, with one button that replays the chain and ' +
          'names the first row that fails.',
        cells: {
          grit: G('SHA-256 chain, plus OCSF export for a SIEM'),
          laravel: C('owen-it/laravel-auditing and similar'),
          django: C('django-auditlog and similar'),
          nextjs: N(),
          encore: N(),
          express: N(),
        },
      },
      {
        id: 'gdpr',
        label: 'GDPR export and erasure',
        detail: 'A subject access export and a delete that actually cascades.',
        cells: {
          grit: G('Endpoints and admin screens'),
          laravel: C(),
          django: C(),
          nextjs: N(),
          encore: N(),
          express: N(),
        },
      },
      {
        id: 'hardening',
        label: 'Hardened by default',
        detail:
          'Security headers, CSRF, body-size limits, rate limiting and lockout wired before you ' +
          'write a line.',
        cells: {
          grit: G('OWASP-aligned defaults in the scaffold'),
          laravel: B('CSRF and headers are core; rate limiting is core'),
          django: B('CSRF, XSS and clickjacking middleware are core'),
          nextjs: C('Headers are yours to configure'),
          encore: B('Validation and CORS handled by the runtime'),
          express: C('helmet, express-rate-limit, and so on'),
        },
      },
    ],
  },
  {
    title: 'Operations',
    blurb:
      'Everything that turns a demo into something you can run. Laravel is genuinely strong here ' +
      'and the table says so.',
    rows: [
      {
        id: 'jobs',
        label: 'Background jobs with a queue UI',
        detail: 'A worker, retries, and a screen showing what failed and why.',
        cells: {
          grit: G('asynq, with the dashboard scaffolded'),
          laravel: B('Queues are core; Horizon adds the UI as a first-party package'),
          django: C('Celery, plus Flower for the UI'),
          nextjs: N(),
          encore: B('Pub/Sub and cron are part of the framework'),
          express: C('BullMQ and similar'),
        },
      },
      {
        id: 'observability',
        label: 'Request tracing and metrics in-app',
        detail: 'Latency, error rates and slow queries, without adding a vendor.',
        cells: {
          grit: G('Pulse, in the generated admin'),
          laravel: O('Telescope'),
          django: C('django-silk or debug-toolbar, development oriented'),
          nextjs: N(),
          encore: B('Tracing is a core part of the platform'),
          express: C(),
        },
      },
      {
        id: 'storage',
        label: 'File storage with image processing',
        detail: 'S3-compatible uploads, presigned URLs, and thumbnails generated for you.',
        cells: {
          grit: G('S3, R2 or MinIO, with a dropzone in the admin'),
          laravel: B('Flysystem is core; image processing is a package'),
          django: B('Storage backends are core; processing is a package'),
          nextjs: C(),
          encore: B('Object storage is a framework primitive'),
          express: C('multer, and an SDK'),
        },
      },
      {
        id: 'backup',
        label: 'Scheduled database backup',
        detail: 'A cron entry, an off-site target, and a restore path you can test.',
        cells: {
          grit: G('Scheduled, with a screen'),
          laravel: C('spatie/laravel-backup'),
          django: C('django-dbbackup'),
          nextjs: N(),
          encore: N('Managed by the cloud provider'),
          express: N(),
        },
      },
      {
        id: 'deploy',
        label: 'Deploy command',
        detail: 'One command from a working tree to a running server with TLS.',
        cells: {
          grit: B('grit deploy: SSH, systemd, Caddy with automatic TLS'),
          laravel: C('Forge and Envoyer are paid services'),
          django: N(),
          nextjs: C('Trivial on Vercel, yours anywhere else'),
          encore: B('Deploys to Encore Cloud or your own AWS and GCP'),
          express: N(),
        },
      },
      {
        id: 'hosting',
        label: 'One-click managed hosting',
        detail:
          'Push, and someone else runs it. Grit deploys to a server you own, which is cheaper ' +
          'and more portable, and is more work than a git push to a platform that knows your ' +
          'framework.',
        cells: {
          grit: C('Any VPS or container host; no platform is tailored to it'),
          laravel: C('Forge and Vapor are paid services'),
          django: C(),
          nextjs: B('Vercel is built by the same team and it shows'),
          encore: B('Encore Cloud is part of the product'),
          express: C(),
        },
      },
      {
        id: 'i18n',
        label: 'Internationalisation',
        detail:
          'Translation catalogues, locale negotiation and pluralisation. Grit formats numbers ' +
          'and dates by locale but ships no translation system, which is a real gap next to the ' +
          'two frameworks that have had one for twenty years.',
        cells: {
          grit: O('grit add i18n: next-intl plus translated API messages. The generated admin chrome is not translated yet, so this is a foundation rather than a finished feature'),
          laravel: B('Translation files, helpers and pluralisation are core'),
          django: B('gettext, locale middleware and translated admin'),
          nextjs: C('next-intl or similar'),
          encore: N(),
          express: C('i18next and similar'),
        },
      },
    ],
  },
  {
    title: 'Developer experience',
    blurb:
      'The row that explains the benchmark result as well: less code on the request path because ' +
      'less code was written by hand.',
    rows: [
      {
        id: 'scaffold',
        label: 'Full-stack resource generation',
        detail:
          'One command emitting the model, migration, service, handler, validation schema, ' +
          'TypeScript types, data-fetching hooks and an admin screen.',
        cells: {
          grit: G('grit generate resource, every layer at once'),
          laravel: O('make:model -mcr covers the backend layers'),
          django: N('startapp gives you empty files'),
          nextjs: N(),
          encore: N(),
          express: N(),
        },
      },
      {
        id: 'typedclient',
        label: 'Typed client from the backend',
        detail: 'Backend types crossing the language boundary without being retyped by hand.',
        cells: {
          grit: G('Go structs to TypeScript types, Zod schemas and React Query hooks'),
          laravel: C('Typescript transformers and similar'),
          django: C('drf-spectacular plus a generator'),
          nextjs: N('Same language, so the question does not arise'),
          encore: B('Generates typed clients from your API definitions'),
          express: C('Via an OpenAPI generator you wire up'),
        },
      },
      {
        id: 'admin',
        label: 'Admin panel',
        detail: 'Tables, filters, forms, bulk actions and detail views over your own models.',
        cells: {
          grit: G('Generated per resource, in four themes'),
          laravel: C('Filament is third party; Nova is first party and paid'),
          django: B('The Django admin, and it is still the benchmark'),
          nextjs: N(),
          encore: N(),
          express: N(),
        },
      },
      {
        id: 'admincustom',
        label: 'Customising the admin',
        detail:
          'Override one cell, swap the whole table, or replace the page, without forking the ' +
          'generated file. Regenerating never eats the overrides.',
        cells: {
          grit: G('A .custom.tsx overlay per resource, separate from the generated definition'),
          laravel: O('Filament, which is its own framework to learn'),
          django: B('ModelAdmin subclassing, in Python templates'),
          nextjs: N('There is no admin to customise'),
          encore: N(),
          express: N(),
        },
      },
      {
        id: 'dbbrowser',
        label: 'Database browser in the app',
        detail: 'Browse and edit rows without leaving the running application.',
        cells: {
          grit: B('GORM Studio, mounted at /studio'),
          laravel: C('Tinker is a REPL rather than a browser'),
          django: B('The admin covers much of this'),
          nextjs: N(),
          encore: B('A local development dashboard'),
          express: N(),
        },
      },
      {
        id: 'apidocs',
        label: 'API documentation from the code',
        detail: 'An OpenAPI document and a browsable reference, kept current by the generator.',
        cells: {
          grit: G('Emitted as routes are generated, served in-app'),
          laravel: C('Scribe or L5-Swagger'),
          django: C('drf-spectacular'),
          nextjs: N(),
          encore: B('Generated from the API definitions'),
          express: C('swagger-jsdoc and similar'),
        },
      },
      {
        id: 'graphql',
        label: 'GraphQL API',
        detail:
          'Grit is REST and OpenAPI only. If GraphQL is a requirement, the mature options are ' +
          'elsewhere and this is the wrong framework for the job.',
        cells: {
          grit: N('REST and OpenAPI only'),
          laravel: C('Lighthouse'),
          django: C('Strawberry or Graphene'),
          nextjs: C('Apollo or similar'),
          encore: N(),
          express: C('Apollo Server'),
        },
      },
    ],
  },
  {
    title: 'Reach',
    blurb:
      'One Go backend and one set of generated clients. This is the row with the least company: ' +
      'nothing else here scaffolds past the browser.',
    rows: [
      {
        id: 'web',
        label: 'Web frontend',
        detail: 'A public site or app, typed against the same API.',
        cells: {
          grit: G('Next.js or Vite, your choice at scaffold time'),
          laravel: B('Blade, with Inertia or Livewire'),
          django: B('Templates'),
          nextjs: B('This is what Next.js is'),
          encore: N('Bring your own'),
          express: N('Bring your own'),
        },
      },
      {
        id: 'templating',
        label: 'Server-rendered templating',
        detail:
          'A page rendered by the backend with no JavaScript build. Grit is an API with React ' +
          'clients by design, so for a content site with a little interactivity Blade or Django ' +
          'templates are simply less machinery.',
        cells: {
          grit: N('API plus React by design'),
          laravel: B('Blade, with Livewire for interactivity'),
          django: B('The template language is core'),
          nextjs: B('Server components render on the server'),
          encore: N(),
          express: C('ejs, pug and similar'),
        },
      },
      {
        id: 'mobile',
        label: 'Mobile app',
        detail: 'A React Native client sharing the API types.',
        cells: {
          grit: G('Expo, scaffolded with the shared types'),
          laravel: N(),
          django: N(),
          nextjs: N(),
          encore: N(),
          express: N(),
        },
      },
      {
        id: 'desktop',
        label: 'Desktop app',
        detail: 'A native window, offline-first with a local database and sync.',
        cells: {
          grit: G('Wails, with an offline SQLite store and a sync engine'),
          laravel: C('NativePHP is young and third party'),
          django: N(),
          nextjs: C('Electron or Tauri, wired by you'),
          encore: N(),
          express: N(),
        },
      },
    ],
  },
]

export const PROVISION_LABEL: Record<Provision, string> = {
  generated: 'Generated for you',
  builtin: 'Built into the framework',
  official: 'Official package, you wire it',
  community: 'Third-party package',
  none: 'Not provided',
}

export const PROVISION_SHORT: Record<Provision, string> = {
  generated: 'Generated',
  builtin: 'Built in',
  official: 'Official add-on',
  community: 'Third party',
  none: 'None',
}

/** How many rows arrive without you installing or wiring anything. */
export function readyCount(slug: string): number {
  return CAPABILITY_GROUPS.flatMap((g) => g.rows).filter((r) => {
    const s = r.cells[slug]?.state
    return s === 'generated' || s === 'builtin'
  }).length
}

export const TOTAL_ROWS = CAPABILITY_GROUPS.reduce((n, g) => n + g.rows.length, 0)
