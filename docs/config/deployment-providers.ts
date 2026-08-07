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
  'Grit does not auto-migrate on boot in production: a process that rewrites the schema every time it restarts is a bad idea when the platform can restart it for its own reasons. Run migrations as an explicit step.'

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
      'Teams that want a click-through dashboard for everything: Fly is CLI-first and expects you to read a TOML file.',
    managedPostgres: true,
    managedRedis: true,
    persistentDisk: true,
    docsUrl: 'https://fly.io/docs/',
    steps: [
      {
        title: 'Launch without deploying',
        body: 'Let Fly detect the Dockerfile and write a fly.toml, but stop before it ships anything: the defaults need two changes first.',
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
      'Machines suspend on idle by default. The first request after a quiet period pays the wake-up cost: set `min_machines_running = 1` for anything user-facing.',
      'A volume is attached to one machine in one region. Scale past one machine and each gets its own empty disk, so uploads must go to S3/R2 rather than local storage.',
      'Fly Postgres is an unmanaged Postgres you own, not a managed service. Backups are your job: `fly postgres` will not do point-in-time recovery for you.',
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
      'Cost-sensitive workloads that run hot: usage billing is generous until it is not, and there is no hard spend cap.',
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
        body: 'Add them from the same project canvas. Railway exposes them as variables you reference rather than copy, so a rotated password propagates instead of silently breaking the app.',
        code: {
          language: 'bash',
          code: `DATABASE_URL=\${{Postgres.DATABASE_URL}}
REDIS_URL=\${{Redis.REDIS_URL}}`,
        },
      },
      {
        title: 'Bind to the port Railway gives you',
        body: 'Railway injects PORT. Grit reads APP_PORT, so map one to the other rather than hardcoding 8080: the assigned port is not stable.',
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
      'Railway imports a compose file, it does not run one. Your services are translated into its own model, so service-to-service hostnames become <service>.railway.internal rather than the plain compose service name, and every connection string has to be updated.',
      'depends_on does not gate startup here. Compose can wait for a healthcheck before starting the API; Railway starts everything at once, so the API must retry its first database connection rather than assume Postgres is up.',
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
      'Anything latency-sensitive on the free tier: free services spin down after inactivity and cold starts run to tens of seconds.',
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
      'Health check failures roll the deploy back rather than leaving it half-live: good behaviour, but it means a wrong `healthCheckPath` looks like "my deploy never finishes".',
      'Free web services sleep. The first request wakes them, and the wake is slow enough that a health check can time out first.',
    ],
  },

  /* ── Self-hosted PaaS ──────────────────────────────────────────────── */
  {
    slug: 'orbita',
    name: 'Orbita',
    tagline:
      'A self-hosted PaaS that already knows what a Grit app is. It reads grit.json and derives the rest.',
    kind: 'self-hosted',
    costFrom: 'Cost of the VPS (~$5)',
    effort: 'Low',
    bestFor:
      'Grit apps specifically. It is the only target that reads grit.json, so a three-service deploy is about ten lines of orbita.yaml rather than a Compose file, a Blueprint or three TOMLs. Multi-tenant, so one box can host several clients in isolated organisations.',
    notFor:
      'Non-Grit apps that need the fast path — those fall back to the Compose route, where only one service per app is routable. Also not for anyone who does not want to own an operating system.',
    managedPostgres: true,
    managedRedis: true,
    persistentDisk: true,
    docsUrl: 'https://github.com/MUKE-coder/orbita',
    steps: [
      {
        title: 'Install it once on a fresh VPS',
        body: 'The installer brings Docker, Swarm, Postgres, Redis and Traefik with it, so there is nothing to set up beforehand. Give it 2 vCPU and 4 GB: Orbita idles around 50 MB, but building three services does not.',
        code: {
          language: 'bash',
          code: `curl -sSL https://raw.githubusercontent.com/MUKE-coder/orbita/main/install.sh \\
  | sudo ORBITA_DOMAIN=orbita.example.com ORBITA_ACME_EMAIL=you@example.com bash -s -- --yes`,
        },
      },
      {
        title: 'Register first, immediately',
        body: 'The first account to register becomes super-admin, and public sign-up closes the moment it exists. Leaving that gap open on a public IP is the one genuinely dangerous minute in the install.',
      },
      {
        title: 'Write orbita.yaml',
        body: 'This is the whole deploy config. Addons replace the Postgres and Redis services you would declare in Compose, domains replace the Traefik labels, and migrate replaces the one-shot migrate container.',
        code: {
          language: 'yaml',
          code: `app: storefront
repo: your-org/storefront
branch: main

addons: [postgres, redis]

domains:
  web:   example.com
  admin: admin.example.com
  api:   api.example.com

migrate: true`,
        },
      },
      {
        title: 'Deploy, and look at the plan first',
        body: 'The dry run prints exactly what it will create before anything happens, which is the cheapest way to catch a wrong domain.',
        code: {
          language: 'bash',
          code: `orbita deploy --plan --host prod
orbita deploy --host prod`,
        },
      },
    ],
    gotchas: [
      'Back up /opt/orbita/.env. It holds ENCRYPTION_MASTER_KEY, from which every organisation key is derived — lose it and every stored secret is unrecoverable.',
      'Migrations run before cutover under a Postgres advisory lock, and a non-zero exit aborts the deploy. The usual cause of a failure there is go.sum not being committed, so the one-off container cannot resolve modules.',
      'On the Compose path only one service per app is routable, which is why a three-domain app wants the Grit fast path instead.',
      'If your DNS is behind Cloudflare, set the records to DNS only for the first certificate — the proxy interferes with the ACME handshake.',
    ],
  },
  {
    slug: 'dokploy',
    name: 'Dokploy',
    tagline: 'A Heroku-like panel on your own server. Grit’s own sites run on it.',
    kind: 'self-hosted',
    costFrom: 'Cost of the VPS (~$5)',
    effort: 'Medium',
    bestFor:
      'Running several apps on one box with a UI, automatic TLS and Git deploys: without paying per service.',
    notFor:
      'Anyone who does not want to own an operating system. You are the one patching it.',
    managedPostgres: false,
    managedRedis: false,
    persistentDisk: true,
    docsUrl: 'https://docs.dokploy.com/',
    steps: [
      {
        title: 'Install on a fresh VPS',
        body: 'One command on a clean Ubuntu box. Give it 2 GB of RAM minimum: Dokploy plus a Go API plus Postgres will not fit comfortably in 1 GB.',
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
        body: 'Anything read while the frontend compiles (NEXT_PUBLIC_*, and anything baked into prerendered HTML) must be a Build Argument. Set only at runtime, it is simply absent from the built output, and nothing warns you.',
        code: {
          language: 'text',
          code: `Build arguments:  NEXT_PUBLIC_API_URL
Environment:      DATABASE_URL, REDIS_URL, JWT_SECRET`,
        },
      },
      {
        title: 'Or deploy the whole stack as Compose',
        body: 'The steps above deploy the API as a single Dockerfile application. To bring up API, web, Postgres and Redis together, add a Compose service instead and point it at the production compose file. Choose Compose rather than Stack: Stack targets Docker Swarm and does not support the build key, and three of the services build from a Dockerfile.',
        code: {
          language: 'text',
          code: `Add Service:    Compose
Provider:       GitHub -> your repository
Compose Path:   ./docker-compose.prod.yml
Compose Type:   docker-compose`,
        },
      },
      {
        title: 'Give the containers their environment',
        body: 'Dokploy writes dashboard variables into a .env beside the compose file, and does not inject them into containers. That file only reaches a container if the compose file asks for it, so either keep the ${VAR} substitution in every environment block or add env_file to each service. Skipping this is why a stack boots with an empty DATABASE_URL despite the dashboard being full.',
        code: {
          language: 'yaml',
          code: `services:
  api:
    env_file:
      - .env        # or keep \${VAR} in the environment block
    environment:
      DATABASE_URL: \${DATABASE_URL}`,
        },
      },
      {
        title: 'Attach domains per service',
        body: 'On a Compose deployment the Domains tab asks which service a domain belongs to, so web and api each get their own. Dokploy generates the Traefik labels; you do not write them. Preview Compose shows exactly what it will inject before you deploy, which is worth reading once.',
      },
      {
        title: 'Add the domain and TLS',
        body: 'Dokploy provisions TLS through Traefik. If Cloudflare sits in front with the orange cloud on, set SSL/TLS to Full (strict): Flexible produces a redirect loop that looks like an application bug.',
      },
    ],
    gotchas: [
      'Postgres and Redis run as containers you own. Backups, upgrades and disk pressure are yours: see the backup page.',
      'The server is a single point of failure. Fine for internal tools, a real decision for anything customer-facing.',
      'Dokploy updates itself in place. Snapshot the volume before a major upgrade.',
      'Remove explicit container_name values from the compose file. Dokploy suffixes names so two projects can share a server, and an explicit container_name overrides that, which means a staging copy of the same stack on the same box collides with production.',
      'Keep using ${VAR} substitution in the environment blocks. Dashboard variables are written to a .env beside the compose file, so they reach a running container through Compose substitution rather than being injected directly.',
      'Auto-deploy re-clones the repository on every deploy, which wipes the working directory. Anything written there at runtime, uploads included, is gone on the next push. Use a named volume or Dokploy File Mounts, never a path inside the repo.',
      'A custom deploy command replaces the default rather than adding to it. The default is docker compose -p <name> -f <path> up -d --build --remove-orphans, so anything you drop from it stops happening.',
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
    notFor: 'Minimal setups: it carries more moving parts than a plain Docker Compose file.',
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
      'Coolify runs the real Compose engine, which is why it needs the fewest edits: depends_on with condition: service_healthy behaves exactly as it does locally, and services still resolve each other by compose service name.',
    ],
  },
  {
    slug: 'aws-ec2',
    name: 'AWS EC2',
    tagline:
      'A bare Linux box plus an Application Load Balancer doing TLS and host-based routing.',
    kind: 'vps',
    costFrom: '~$30 (t3.medium + ALB)',
    effort: 'High',
    bestFor:
      'Teams already on AWS who want the standard production pattern there: ACM certificates, ALB routing, and everything inside a VPC they control.',
    notFor:
      'A first deploy, or a solo project. You are assembling the load balancer, the certificate, three target groups and the DNS yourself, and the ALB alone costs more per month than most of the PaaS options.',
    managedPostgres: false,
    managedRedis: false,
    persistentDisk: true,
    docsUrl: 'https://docs.aws.amazon.com/elasticloadbalancing/latest/application/introduction.html',
    steps: [
      {
        title: 'Two security groups, not one',
        body: 'The instance group allows 8080/3000/3001 only from the ALB group — as a source, not an IP range. That is the single control keeping your containers off the public internet.',
      },
      {
        title: 'Size the instance for the build, not the traffic',
        body: 'Building two Next.js apps and a Go API back to back OOMs on t3.micro and t3.small. Start at t3.medium with 40 GB of gp3, and add swap.',
        code: {
          language: 'bash',
          code: `sudo fallocate -l 2G /swapfile
sudo chmod 600 /swapfile
sudo mkswap /swapfile && sudo swapon /swapfile`,
        },
      },
      {
        title: 'Adjust the Compose file for a bare VM',
        body: 'Drop the Traefik labels and the dokploy-network, then add host port mappings so the ALB target groups have something to reach. Postgres and Redis get no ports at all.',
        code: {
          language: 'yaml',
          code: `services:
  api:
    ports: ["8080:8080"]
  admin:
    ports: ["3001:3000"]
  web:
    ports: ["3000:3000"]`,
        },
      },
      {
        title: 'One certificate, three hostnames',
        body: 'Request a single ACM certificate with all three hosts as subject alternative names, validated by DNS so it renews itself.',
      },
    ],
    gotchas: [
      'An unhealthy target group is almost always the health check path returning non-2xx, or the instance security group not allowing the ALB security group on that port.',
      'depends_on: condition: service_completed_successfully works here unchanged — this is real Docker Compose on a real VM, so the migrate job needs no pre-deploy workaround.',
      'There is no GitHub webhook on a bare instance. Either redeploy over SSH by hand or add a GitHub Actions workflow that does it for you.',
      'If you do not need ALB features, a Caddy container on the instance does host routing and automatic certificates for nothing — see the Lightsail guide.',
    ],
  },
  {
    slug: 'aws-lightsail',
    name: 'AWS Lightsail',
    tagline:
      'AWS at a flat monthly price, with a static IP and snapshots, fronted by Caddy for automatic HTTPS.',
    kind: 'vps',
    costFrom: '~$24 (4 GB plan)',
    effort: 'Medium',
    bestFor:
      'A predictable AWS bill without the VPC and ALB assembly. Automatic daily snapshots of the whole instance disk are the lowest-effort backup baseline of any target here.',
    notFor:
      'Anyone needing host-based routing from an AWS load balancer: Lightsail’s own balancer cannot do it, which is why this runs Caddy on the instance instead.',
    managedPostgres: false,
    managedRedis: false,
    persistentDisk: true,
    docsUrl: 'https://docs.aws.amazon.com/lightsail/',
    steps: [
      {
        title: 'Take the 4 GB plan',
        body: 'The 2 GB plan OOMs building two Next.js apps and a Go API in one compose build. Add 2 GB of swap on top.',
      },
      {
        title: 'Static IP, then firewall',
        body: 'Attach a static IP so DNS does not break on restart, then open only 22 (from your IP), 80 and 443. The app ports stay closed — traffic reaches them through Caddy.',
      },
      {
        title: 'Add a Caddy service and a Caddyfile',
        body: 'Caddy does the host-based routing the Lightsail load balancer cannot, and fetches and renews Let’s Encrypt certificates on its own. No certbot, no renewal cron, no ACM console.',
        code: {
          language: 'text',
          code: `api.example.com {
    reverse_proxy api:8080
}

admin.example.com {
    reverse_proxy admin:3000
}

example.com {
    reverse_proxy web:3000
}`,
        },
      },
      {
        title: 'Point DNS before the first deploy',
        body: 'Caddy validates over HTTP-01 on port 80, so the records must already resolve when the container starts. Deploy first and the certificate request fails and backs off.',
      },
    ],
    gotchas: [
      'If Caddy never issues a certificate, DNS was not resolving to the static IP when it first started. Fix the records, then restart the caddy container.',
      'Automatic snapshots restore the whole instance to a point in time, not just the database. Keep pg_dump backups as well for anything you would hate to roll back wholesale.',
      'A 502 from Caddy means the upstream container is not listening yet or has crashed — check the app container logs, not Caddy’s.',
      'Lightsail Container Service is a different product entirely: it runs images rather than a Compose file, and is closer to Railway than to this guide.',
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
