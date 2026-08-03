/**
 * Deployment targets, one entry per hosting provider.
 *
 * Data rather than nine near-identical page files: every provider needs the same
 * spine — what it costs, what it runs, how the build works, what breaks — and a
 * shared shape is the only way that stays true after the fourth one is added.
 *
 * WHAT GOES IN HERE: the parts that come from Grit and change slowly — the
 * Dockerfile it generates, the port it listens on, the health path, the migrate
 * command, the environment it needs.
 *
 * WHAT DOES NOT: click-by-click walkthroughs of somebody else's dashboard. Those
 * rot within a release and leave the docs confidently wrong, which is worse than
 * being brief. Every page links out for the parts that move.
 */

export type ProviderKind = 'paas' | 'self-hosted' | 'vps' | 'container'

export interface ProviderStep {
  title: string
  body: string
  code?: { language: string; code: string }
}

export interface Provider {
  slug: string
  name: string
  /** One line for the index card. */
  tagline: string
  kind: ProviderKind
  /** Rough monthly floor for a small production app, in USD. */
  costFrom: string
  /** How much operational work you are signing up for. */
  effort: 'Lowest' | 'Low' | 'Medium' | 'High'
  /** Who this is genuinely the right answer for. */
  bestFor: string
  /** And who it is not. Written plainly — steering someone away is a service. */
  notFor: string
  managedPostgres: boolean
  managedRedis: boolean
  /** Persistent disk available, which decides whether SQLite or local uploads work. */
  persistentDisk: boolean
  docsUrl: string
  steps: ProviderStep[]
  gotchas: string[]
}

/* Shared fragments. Repeating them per provider is how five copies drift and
   three end up wrong. */

const MIGRATE_NOTE =
  'Grit does not auto-migrate on boot in production — a process that rewrites the schema every time it restarts is a bad idea when the platform can restart it for its own reasons. Run migrations as an explicit step.'

const HEALTHCHECK = '/api/health'

export const DEPLOYMENT_PROVIDERS: Provider[] = [
  /* ── PaaS ──────────────────────────────────────────────────────────── */
  {
    slug: 'fly-io',
    name: 'Fly.io',
    tagline: 'Runs the Docker image close to your users, with a real disk if you need one.',
    kind: 'paas',
    costFrom: '~$5',
    effort: 'Low',
    bestFor:
      'A single Go binary you want in several regions without running servers. The best fit for Grit’s single-binary mode.',
    notFor:
      'Teams that want a click-through dashboard for everything — Fly is CLI-first and expects you to read a TOML file.',
    managedPostgres: true,
    managedRedis: true,
    persistentDisk: true,
    docsUrl: 'https://fly.io/docs/',
    steps: [
      {
        title: 'Launch without deploying',
        body: 'Let Fly detect the Dockerfile and write a fly.toml, but stop before it ships anything — the defaults need two changes first.',
        code: { language: 'bash', code: 'fly launch --no-deploy' },
      },
      {
        title: 'Point the health check at the API',
        body: `Fly checks \`/\` by default. Grit serves its health endpoint at \`${HEALTHCHECK}\`, so without this the machine is marked unhealthy and cycled forever while the app is running perfectly.`,
        code: {
          language: 'toml',
          code: `[http_service]
  internal_port = 8080
  force_https = true
  auto_stop_machines = "suspend"
  auto_start_machines = true
  min_machines_running = 1

  [[http_service.checks]]
    interval = "15s"
    timeout = "3s"
    grace_period = "10s"
    method = "GET"
    path = "${HEALTHCHECK}"`,
        },
      },
      {
        title: 'Attach Postgres and Redis',
        body: 'Attaching sets DATABASE_URL for you. Redis comes from Upstash through Fly and gives you REDIS_URL.',
        code: {
          language: 'bash',
          code: `fly postgres create --name my-app-db
fly postgres attach my-app-db

fly redis create`,
        },
      },
      {
        title: 'Set the secrets',
        body: 'Secrets are encrypted and injected at runtime. Setting them triggers a redeploy, so do it before the first one.',
        code: {
          language: 'bash',
          code: `fly secrets set \\
  JWT_SECRET="$(openssl rand -base64 32)" \\
  APP_ENV=production \\
  CORS_ORIGINS="https://your-domain.com"`,
        },
      },
      {
        title: 'Deploy, then migrate',
        body: MIGRATE_NOTE,
        code: {
          language: 'bash',
          code: `fly deploy
fly ssh console -C "/app/migrate"`,
        },
      },
    ],
    gotchas: [
      'Machines suspend on idle by default. The first request after a quiet period pays the wake-up cost — set `min_machines_running = 1` for anything user-facing.',
      'A volume is attached to one machine in one region. Scale past one machine and each gets its own empty disk, so uploads must go to S3/R2 rather than local storage.',
      'Fly Postgres is an unmanaged Postgres you own, not a managed service. Backups are your job — `fly postgres` will not do point-in-time recovery for you.',
    ],
  },
  {
    slug: 'railway',
    name: 'Railway',
    tagline: 'Push to Git, get a URL. The least ceremony of anything here.',
    kind: 'paas',
    costFrom: '~$5',
    effort: 'Lowest',
    bestFor:
      'Getting something in front of people today. Postgres and Redis are one click and wire themselves into your service.',
    notFor:
      'Cost-sensitive workloads that run hot — usage billing is generous until it is not, and there is no hard spend cap.',
    managedPostgres: true,
    managedRedis: true,
    persistentDisk: true,
    docsUrl: 'https://docs.railway.com/',
    steps: [
      {
        title: 'Create the service from your repo',
        body: 'Railway finds the Dockerfile and builds it. For a monorepo, set the service root to the API directory so it does not try to build the whole tree.',
        code: { language: 'text', code: 'Root directory:  apps/api\nBuilder:         Dockerfile' },
      },
      {
        title: 'Add Postgres and Redis',
        body: 'Add them from the same project canvas. Railway exposes them as variables you reference rather than copy — so a rotated password propagates instead of silently breaking the app.',
        code: {
          language: 'bash',
          code: `DATABASE_URL=\${{Postgres.DATABASE_URL}}
REDIS_URL=\${{Redis.REDIS_URL}}`,
        },
      },
      {
        title: 'Bind to the port Railway gives you',
        body: 'Railway injects PORT. Grit reads APP_PORT, so map one to the other rather than hardcoding 8080 — the assigned port is not stable.',
        code: { language: 'bash', code: 'APP_PORT=${{PORT}}' },
      },
      {
        title: 'Run migrations',
        body: `${MIGRATE_NOTE} Railway’s one-off command runs in the same environment as the service, so it sees the same DATABASE_URL.`,
        code: { language: 'bash', code: 'railway run /app/migrate' },
      },
    ],
    gotchas: [
      'The generated Dockerfile builds only the API. A monorepo with a separate web app needs a second service, not a second process in the same one.',
      'Sleeping is not available on all plans; an always-on service bills continuously even with no traffic.',
      'Set a usage alert on day one. Railway bills by consumption and a runaway job queue is an expensive way to find that out.',
    ],
  },
  {
    slug: 'render',
    name: 'Render',
    tagline: 'Managed Postgres with real backups, and a free tier you can demo on.',
    kind: 'paas',
    costFrom: '$0 (free tier) / ~$7 production',
    effort: 'Lowest',
    bestFor:
      'Teams that want managed Postgres with point-in-time recovery without thinking about it, and a blueprint file checked into the repo.',
    notFor:
      'Anything latency-sensitive on the free tier — free services spin down after inactivity and cold starts run to tens of seconds.',
    managedPostgres: true,
    managedRedis: true,
    persistentDisk: true,
    docsUrl: 'https://render.com/docs',
    steps: [
      {
        title: 'Describe the stack in render.yaml',
        body: 'Committing the blueprint means the environment is reviewable in a pull request instead of living in a dashboard nobody can diff.',
        code: {
          language: 'yaml',
          code: `services:
  - type: web
    name: my-app-api
    runtime: docker
    dockerfilePath: ./apps/api/Dockerfile
    healthCheckPath: ${HEALTHCHECK}
    envVars:
      - key: APP_ENV
        value: production
      - key: JWT_SECRET
        generateValue: true
      - key: DATABASE_URL
        fromDatabase:
          name: my-app-db
          property: connectionString

databases:
  - name: my-app-db
    plan: basic-256mb`,
        },
      },
      {
        title: 'Run migrations before the first boot',
        body: `${MIGRATE_NOTE} On Render that is a pre-deploy command, which runs after the build and before traffic moves over.`,
        code: { language: 'text', code: 'Pre-Deploy Command:  /app/migrate' },
      },
    ],
    gotchas: [
      'Free Postgres instances expire after 90 days. Fine for a demo, quietly fatal for anything you forgot about.',
      'Health check failures roll the deploy back rather than leaving it half-live — good behaviour, but it means a wrong `healthCheckPath` looks like "my deploy never finishes".',
      'Free web services sleep. The first request wakes them, and the wake is slow enough that a health check can time out first.',
    ],
  },

  /* ── Self-hosted PaaS ──────────────────────────────────────────────── */
  {
    slug: 'dokploy',
    name: 'Dokploy',
    tagline: 'A Heroku-like panel on your own server. Grit’s own sites run on it.',
    kind: 'self-hosted',
    costFrom: 'Cost of the VPS (~$5)',
    effort: 'Medium',
    bestFor:
      'Running several apps on one box with a UI, automatic TLS and Git deploys — without paying per service.',
    notFor:
      'Anyone who does not want to own an operating system. You are the one patching it.',
    managedPostgres: false,
    managedRedis: false,
    persistentDisk: true,
    docsUrl: 'https://docs.dokploy.com/',
    steps: [
      {
        title: 'Install on a fresh VPS',
        body: 'One command on a clean Ubuntu box. Give it 2 GB of RAM minimum — Dokploy plus a Go API plus Postgres will not fit comfortably in 1 GB.',
        code: { language: 'bash', code: 'curl -sSL https://dokploy.com/install.sh | sh' },
      },
      {
        title: 'Create the application',
        body: 'Point it at your repository and set the build context to the API directory. In a monorepo, set the watch path too, or every commit to the frontend rebuilds the backend.',
        code: {
          language: 'text',
          code: `Build type:           Dockerfile
Docker file:          apps/api/Dockerfile
Docker context path:  apps/api
Watch paths:          apps/api/**`,
        },
      },
      {
        title: 'Know which variables are build-time',
        body: 'Anything read while the frontend compiles — NEXT_PUBLIC_*, and anything baked into prerendered HTML — must be a Build Argument. Set only at runtime, it is simply absent from the built output, and nothing warns you.',
        code: {
          language: 'text',
          code: `Build arguments:  NEXT_PUBLIC_API_URL
Environment:      DATABASE_URL, REDIS_URL, JWT_SECRET`,
        },
      },
      {
        title: 'Add the domain',
        body: 'Dokploy provisions TLS through Traefik. If Cloudflare sits in front with the orange cloud on, set SSL/TLS to Full (strict) — Flexible produces a redirect loop that looks like an application bug.',
      },
    ],
    gotchas: [
      'Postgres and Redis run as containers you own. Backups, upgrades and disk pressure are yours — see the backup page.',
      'The server is a single point of failure. Fine for internal tools, a real decision for anything customer-facing.',
      'Dokploy updates itself in place. Snapshot the volume before a major upgrade.',
    ],
  },
  {
    slug: 'coolify',
    name: 'Coolify',
    tagline: 'Open-source self-hosted PaaS. Similar shape to Dokploy, larger ecosystem.',
    kind: 'self-hosted',
    costFrom: 'Cost of the VPS (~$5)',
    effort: 'Medium',
    bestFor:
      'Self-hosting with a big library of one-click services alongside your app, and multi-server support once one box is not enough.',
    notFor: 'Minimal setups — it carries more moving parts than a plain Docker Compose file.',
    managedPostgres: false,
    managedRedis: false,
    persistentDisk: true,
    docsUrl: 'https://coolify.io/docs/',
    steps: [
      {
        title: 'Install',
        body: 'Same shape as Dokploy: one script on a clean VPS, then everything else through the panel.',
        code: { language: 'bash', code: 'curl -fsSL https://cdn.coollabs.io/coolify/install.sh | bash' },
      },
      {
        title: 'Add the application',
        body: 'Choose the Dockerfile build pack and set the base directory to the API. Coolify reads the Dockerfile as-is, so the image Grit generates needs no changes.',
        code: {
          language: 'text',
          code: `Build pack:       Dockerfile
Base directory:   /apps/api
Ports exposed:    8080`,
        },
      },
      {
        title: 'Attach databases and set the health check',
        body: `Add Postgres and Redis as resources in the same project, then copy their internal connection strings into the app. Set the health check path to \`${HEALTHCHECK}\` so a failed boot is caught rather than served.`,
      },
    ],
    gotchas: [
      'Internal connection strings use the container name, not localhost. Using localhost is the single most common Coolify support question.',
      'Coolify keeps every build image by default and will fill the disk. Set a cleanup schedule on day one.',
    ],
  },

  /* ── Plain servers ─────────────────────────────────────────────────── */
  {
    slug: 'vps',
    name: 'VPS (systemd + Caddy)',
    tagline: 'No platform at all. `grit deploy` builds, uploads and runs it as a service.',
    kind: 'vps',
    costFrom: '~$5',
    effort: 'Medium',
    bestFor:
      'The cheapest and fastest production setup for a single-binary Grit app. No Docker, no daemon, no control plane — one binary under systemd behind Caddy.',
    notFor:
      'Teams that need rollbacks, several environments, or more than one machine without building that themselves.',
    managedPostgres: false,
    managedRedis: false,
    persistentDisk: true,
    docsUrl: '/docs/infrastructure/deployment',
    steps: [
      {
        title: 'Let the CLI do it',
        body: '`grit deploy` cross-compiles for linux/amd64, uploads over SSH, writes a systemd unit and configures Caddy with automatic TLS. It is the shortest path from a laptop to a live URL that Grit has.',
        code: { language: 'bash', code: 'grit deploy --host your-server.com --user deploy --domain app.example.com' },
      },
      {
        title: 'Or do it by hand',
        body: 'The full manual walkthrough — user creation, hardening, Postgres, systemd unit, Caddyfile, log rotation — is its own page.',
      },
    ],
    gotchas: [
      'A cross-compiled binary that uses cgo will not run on the server. Grit’s default SQLite driver is pure Go for exactly this reason; adding a cgo dependency breaks the build silently until it fails on the box.',
      'There is no rollback. Keep the previous binary next to the current one and know how to swap it back before you need to.',
    ],
  },
  {
    slug: 'docker-compose',
    name: 'Docker Compose',
    tagline: 'The whole stack — API, Postgres, Redis, MinIO — on one machine.',
    kind: 'container',
    costFrom: 'Cost of the host',
    effort: 'Medium',
    bestFor:
      'Staging environments, on-premise installs, and anywhere you need the whole stack reproduced exactly as it runs locally.',
    notFor: 'Zero-downtime deploys. `compose up` stops the old container before the new one is ready.',
    managedPostgres: false,
    managedRedis: false,
    persistentDisk: true,
    docsUrl: '/docs/infrastructure/docker',
    steps: [
      {
        title: 'Use the production compose file',
        body: 'Every Grit project ships `docker-compose.prod.yml` alongside the development one. It differs in the ways that matter: no bind mounts, no exposed database ports, restart policies on.',
        code: {
          language: 'bash',
          code: 'docker compose -f docker-compose.prod.yml up -d --build',
        },
      },
      {
        title: 'Run migrations in the running container',
        body: MIGRATE_NOTE,
        code: { language: 'bash', code: 'docker compose -f docker-compose.prod.yml exec api /app/migrate' },
      },
    ],
    gotchas: [
      'The development compose file publishes Postgres on a host port so you can connect a GUI. The production one must not — that is how a database ends up on the public internet.',
      'Named volumes survive `down`, but not `down -v`. That single flag is the difference between a restart and losing the database.',
    ],
  },
]

export function getProvider(slug: string): Provider | undefined {
  return DEPLOYMENT_PROVIDERS.find((p) => p.slug === slug)
}

export const PROVIDER_KIND_LABEL: Record<ProviderKind, string> = {
  paas: 'Managed platform',
  'self-hosted': 'Self-hosted platform',
  vps: 'Plain server',
  container: 'Containers',
}
