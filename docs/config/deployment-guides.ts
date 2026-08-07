/**
 * The long-form deployment guides, one per provider.
 *
 * These are transcribed from the researched source guides rather than
 * summarised. Every step, every table row and every troubleshooting entry is
 * carried across verbatim, because the value of these documents is that
 * somebody actually did the deploy and wrote down what bit them — a paraphrase
 * loses exactly the sentence that saves an afternoon.
 *
 * WHAT GOES IN HERE: the whole walkthrough, in order, including the parts that
 * are specific to somebody else's dashboard.
 *
 * WHAT DOES NOT: invention. If a step is not in the source guide it does not
 * belong here. Where a platform's UI has moved on, the provider's own docs are
 * linked from the page and are the authority.
 *
 * The example app throughout is called `sentex` and its domains are
 * `*.gritcms.com`. Substitute your own project name and domains; the structure
 * is what transfers.
 */

import { ORBITA, AWS_EC2, AWS_LIGHTSAIL } from './deployment-guides-extra'

export type GuideBlock =
  | { kind: 'p'; text: string }
  | { kind: 'h3'; text: string }
  | { kind: 'code'; language: string; code: string }
  | { kind: 'ul'; items: string[] }
  | { kind: 'ol'; items: GuideItem[] }
  | { kind: 'table'; headers: string[]; rows: string[][] }
  | { kind: 'note'; text: string }

export interface GuideItem {
  text: string
  code?: { language: string; code: string }
  sub?: string[]
}

export interface GuideSection {
  heading: string
  blocks: GuideBlock[]
}

export interface DeploymentGuide {
  /** Matches a Provider slug in deployment-providers.ts. */
  slug: string
  title: string
  intro: string[]
  sections: GuideSection[]
}

/* ────────────────────────────────────────────────────────────────────────── */

const RAILWAY: DeploymentGuide = {
  slug: 'railway',
  title: 'Deploying Sentex to Railway (GitHub Deploy)',
  intro: [
    'This guide translates `docker-compose.prod.yml` into a Railway project, deployed from GitHub. Railway does **not** run `docker-compose.yml` files directly — every service in the Compose file becomes its own Railway service, wired together with Railway’s private network and managed databases instead of Traefik + Docker networks. Follow every step in order; later steps assume earlier ones are done.',
  ],
  sections: [
    {
      heading: '0. Service map (Compose → Railway)',
      blocks: [
        {
          kind: 'table',
          headers: ['Compose service', 'Railway equivalent', 'Source'],
          rows: [
            ['postgres', 'Managed Postgres plugin', 'n/a — provisioned by Railway'],
            ['redis', 'Managed Redis plugin', 'n/a — provisioned by Railway'],
            ['migrate', 'Pre-Deploy Command on the api service', 'apps/api'],
            ['api', 'GitHub-deployed service, root dir apps/api', 'apps/api/Dockerfile'],
            ['admin', 'GitHub-deployed service, root dir repo root', 'apps/admin/Dockerfile'],
            ['web', 'GitHub-deployed service, root dir repo root', 'apps/web/Dockerfile'],
            ['traefik labels', 'Railway’s built-in edge (automatic)', 'n/a'],
            ['dokploy-network', 'Railway’s automatic private network', 'n/a'],
          ],
        },
        { kind: 'p', text: 'Reasons for the two differences you’ll notice below:' },
        {
          kind: 'ul',
          items: [
            '**No standalone `migrate` service.** Railway has no `depends_on: condition: service_completed_successfully`, so a one-shot container that "runs once then exits" isn’t a first-class concept the way it is in Compose. The correct Railway feature for this is a **Pre-Deploy Command** on the `api` service — it runs in its own container, before the new `api` version starts serving traffic, and the deploy is aborted if it fails.',
            '**Postgres/Redis become managed plugins, not containers.** Railway’s guidance is to use its managed database services instead of raw `postgres:16-alpine` / `redis:7-alpine` images — you get backups, connection pooling and a dashboard for free, and you no longer need to manage volumes for them yourself.',
          ],
        },
      ],
    },
    {
      heading: 'Part 1 — Prepare the repository',
      blocks: [
        {
          kind: 'ol',
          items: [
            {
              text: 'Confirm you’re on the branch you intend to deploy (commonly `main`), and that it’s pushed to GitHub. Railway builds from GitHub, so uncommitted local changes will not be deployed.',
            },
            {
              text: 'Confirm the three Dockerfiles exist at:',
              sub: ['apps/api/Dockerfile', 'apps/admin/Dockerfile', 'apps/web/Dockerfile'],
            },
            {
              text: 'Open `apps/admin/Dockerfile` and `apps/web/Dockerfile` and confirm they declare `ARG` for every build-time variable the Compose file passes under `build.args` (`NEXT_PUBLIC_API_URL`, `NEXT_PUBLIC_WEB_URL`, `NEXT_PUBLIC_ADMIN_URL`, `THEME`, `SOCIAL_AUTH_ENABLED`, `NEXT_PUBLIC_DEMO_LOGINS`). Railway only forwards a variable into the Docker build if the Dockerfile has a matching `ARG` line — without it, the value you set in the Railway dashboard is silently ignored and the bundle bakes in whatever the Dockerfile’s default is.',
            },
            {
              text: 'Open `apps/api/Dockerfile` and confirm the `./migrate` and `./seed` binaries referenced in the Compose `command:` are actually produced by the build (i.e. they’re part of the same image `api` runs from — not a separate build stage that gets discarded).',
            },
            {
              text: 'Check `.dockerignore` at the repo root. Since `admin` and `web` build with `context: .` (the whole repo), make sure `.dockerignore` excludes `node_modules`, `.next`, `.git`, and `e2e` so the build context upload to Railway stays small and fast.',
            },
            {
              text: 'You will **not** commit `.env` or `.env.production` — their values move into Railway’s Variables UI in Part 5–7. Keep them open locally for reference while you copy values over.',
            },
          ],
        },
      ],
    },
    {
      heading: 'Part 2 — DNS records',
      blocks: [
        {
          kind: 'p',
          text: 'Before deploying, create three DNS **A** records (or a **CNAME** once you know Railway’s target — see note below) pointing at wherever your traffic will land:',
        },
        {
          kind: 'table',
          headers: ['Host', 'Purpose'],
          rows: [
            ['sentex.gritcms.com', 'web — public shell + login'],
            ['admin.sentex.gritcms.com', 'admin — staff admin panel'],
            ['api.sentex.gritcms.com', 'api — Go API'],
          ],
        },
        {
          kind: 'note',
          text: '**Note on Railway specifically:** unlike the Dokploy/Traefik setup in your Compose file, Railway custom domains are attached as a **CNAME** to a Railway-generated target (e.g. `xxxx.up.railway.app`), not an A record to a static IP. You’ll get the exact CNAME target in Part 8, after each service has a Railway-generated domain to point at. It’s fine to leave DNS for last — just don’t skip it, since `WEB_DOMAIN` / `ADMIN_DOMAIN` / `API_DOMAIN` are baked into the frontend bundles at build time and mismatches show up as CORS/CSP failures, exactly as the comments in your Compose file warn.',
        },
        {
          kind: 'p',
          text: 'Storage stays on Cloudflare R2 and needs no DNS record here, same as in the Compose file — only the `R2_*` variables need to be set (Part 5).',
        },
      ],
    },
    {
      heading: 'Part 3 — Create the Railway project',
      blocks: [
        {
          kind: 'ol',
          items: [
            {
              text: 'Go to [railway.com/dashboard](https://railway.com/dashboard) and log in (or sign up) with GitHub — this makes the GitHub connection in the next part one click instead of a separate OAuth flow.',
            },
            { text: 'Click **+ New Project**.' },
            {
              text: 'Choose **Empty Project**. (Don’t use "Deploy from GitHub repo" on this screen — because this repo is a monorepo with multiple deployable services and a non-JS `apps/api` alongside JS apps, Railway’s automatic monorepo importer won’t produce the exact 4-service layout you need. It’s more reliable to add each service by hand in Part 4–7.)',
            },
            {
              text: 'Rename the project (top left, click the project name) to something like `sentex-production` so it’s unambiguous in your Railway dashboard.',
            },
            {
              text: 'If you have more than one Railway environment planned (e.g. `production` vs `staging`), do this whole guide once per environment — Railway environments are separate blank canvases under the same project. This guide assumes a single `production` environment.',
            },
          ],
        },
      ],
    },
    {
      heading: 'Part 4 — Add the Postgres database',
      blocks: [
        {
          kind: 'ol',
          items: [
            { text: 'Inside your empty project canvas, click **+ New**.' },
            { text: 'Select **Database** → **Add PostgreSQL**.' },
            {
              text: 'Railway provisions a Postgres instance immediately and creates a service card named `Postgres` on the canvas. It automatically generates and exposes these variables on that service (you don’t set these — Railway does): `DATABASE_URL`, `DATABASE_PUBLIC_URL`, `PGHOST`, `PGPORT`, `PGUSER`, `PGPASSWORD`, `PGDATABASE`.',
            },
            {
              text: 'Click into the `Postgres` service → **Variables** tab and leave it as is. You’ll pull these values into `api` via reference variables in Part 5 — don’t retype the password anywhere.',
            },
            {
              text: 'No volume setup is needed here — managed databases handle their own storage, unlike the `postgres-data` named volume in your Compose file.',
            },
          ],
        },
      ],
    },
    {
      heading: 'Part 5 — Add the Redis database',
      blocks: [
        {
          kind: 'ol',
          items: [
            { text: 'Click **+ New** again.' },
            { text: 'Select **Database** → **Add Redis**.' },
            {
              text: 'Railway creates a `Redis` service card and exposes `REDIS_URL` and `REDIS_PUBLIC_URL` automatically, the same way it did for Postgres.',
            },
            {
              text: 'Same as Postgres — nothing to configure here; you’ll reference `REDIS_URL` from the `api` service next.',
            },
          ],
        },
      ],
    },
    {
      heading: 'Part 6 — Deploy the api service (and the migration)',
      blocks: [
        { kind: 'h3', text: '6.1 Create the service' },
        {
          kind: 'ol',
          items: [
            { text: 'Click **+ New** → **GitHub Repo**.' },
            {
              text: 'If this is the first time connecting this GitHub account/org, authorize the Railway GitHub App and grant it access to the `sentex` repository (or "All repositories" if you’re comfortable with that).',
            },
            {
              text: 'Select the `sentex` repository, then select the branch you’re deploying (e.g. `main`).',
            },
            {
              text: 'Railway creates a service and immediately tries to build it. It will likely fail or pick the wrong Dockerfile at this point — that’s expected, since it just tried to build from the repo root. Continue to 6.2 before worrying about that first failed build.',
            },
          ],
        },
        { kind: 'h3', text: '6.2 Point it at apps/api' },
        {
          kind: 'ol',
          items: [
            { text: 'Click into the new service → **Settings** tab.' },
            {
              text: 'Under **Source**, find **Root Directory** and set it to:',
              code: { language: 'text', code: 'apps/api' },
            },
            {
              text: 'This matches `build.context: ./apps/api` in the Compose file. With the root directory set, Railway will look for `Dockerfile` inside `apps/api/` automatically — you don’t need to set a separate Dockerfile path for this service.',
            },
            {
              text: 'Rename the service (click the service name at the top) to `api`, so its private-network hostname becomes `api.railway.internal`.',
            },
            {
              text: 'Under **Build**, scroll to **Watch Paths** and add:',
              code: { language: 'text', code: 'apps/api/**\npackages/**' },
            },
            {
              text: 'This stops a commit that only touches `apps/web` or `apps/admin` from triggering an unnecessary `api` rebuild.',
            },
          ],
        },
        { kind: 'h3', text: '6.3 Set environment variables' },
        {
          kind: 'p',
          text: 'Go to the **Variables** tab and add the following. Use **Raw Editor** to paste several at once (`KEY=VALUE` per line). For anything that references another Railway service, use Railway’s reference-variable syntax `${{ServiceName.VARIABLE}}` instead of hardcoding a value — it stays in sync automatically if credentials rotate.',
        },
        {
          kind: 'code',
          language: 'bash',
          code: `APP_ENV=production

# Pulled from the managed Postgres plugin (Part 4) instead of a container
POSTGRES_HOST=\${{Postgres.PGHOST}}
POSTGRES_PORT=\${{Postgres.PGPORT}}
POSTGRES_USER=\${{Postgres.PGUSER}}
POSTGRES_PASSWORD=\${{Postgres.PGPASSWORD}}
POSTGRES_DB=\${{Postgres.PGDATABASE}}

# Pulled from the managed Redis plugin (Part 5) instead of a container
REDIS_URL=\${{Redis.REDIS_URL}}

# Same public-origin / CORS logic as the Compose file, now pointed at
# real domains instead of Docker network hostnames
APP_URL=https://api.sentex.gritcms.com
CORS_ORIGINS=https://sentex.gritcms.com,https://admin.sentex.gritcms.com`,
        },
        {
          kind: 'p',
          text: 'Then add every other application variable your `.env.production` defines that isn’t Postgres/Redis/routing related — for example your R2 storage credentials, since those are unrelated to Compose networking and carry over as plain values:',
        },
        {
          kind: 'code',
          language: 'bash',
          code: `R2_ACCOUNT_ID=...
R2_ACCESS_KEY_ID=...
R2_SECRET_ACCESS_KEY=...
R2_BUCKET=...
R2_ENDPOINT=...`,
        },
        {
          kind: 'p',
          text: 'Copy across any remaining app-specific secrets (JWT signing keys, OAuth credentials if social auth is enabled, mail provider keys, etc.) from your local `.env.production` the same way — one `KEY=VALUE` line per variable.',
        },
        { kind: 'h3', text: '6.4 Replace the migrate job with a Pre-Deploy Command' },
        {
          kind: 'ol',
          items: [
            { text: 'Still in the `api` service, go to **Settings** → **Deploy**.' },
            {
              text: 'Find **Pre-Deploy Command** and set it to:',
              code: { language: 'bash', code: 'sh -c "./migrate && ./seed"' },
            },
            {
              text: 'Wrap it in `sh -c "..."` exactly like this — because `api` builds from a Dockerfile (not Railpack), Railway executes the pre-deploy command directly rather than through a shell, so `&&` needs an explicit shell to be interpreted correctly. Without the `sh -c` wrapper the command fails with an "exec format" style error.',
            },
            {
              text: 'Leave **Custom Start Command** as whatever `apps/api/Dockerfile`’s `CMD`/`ENTRYPOINT` already runs your compiled API binary with — you don’t need to override it, since the Compose file’s `api` service doesn’t override `command:` either.',
            },
            {
              text: 'This reproduces the ordering guarantees from your Compose file (`migrate` → seed → then `api` starts serving) without needing a separate one-shot service: the pre-deploy command runs in its own container on every deploy, and if it exits non-zero the new `api` version never goes live. It is idempotent on your side already (the seed skips when members exist), so redeploys are safe exactly as documented in the Compose file’s header comment.',
            },
          ],
        },
        { kind: 'h3', text: '6.5 Networking' },
        {
          kind: 'ol',
          items: [
            { text: 'Still under **Settings**, find **Networking**.' },
            {
              text: 'Click **Generate Domain** if you want a temporary `*.up.railway.app` URL to test with before the custom domain is live. This is optional — you can also wait and go straight to the custom domain in Part 8.',
            },
            {
              text: 'You do **not** need to expose port 8080 manually the way `expose: ["8080"]` did in Compose — Railway detects the port your app listens on, or you can pin it explicitly under **Settings → Networking → Port**.',
            },
            {
              text: 'Trigger a redeploy (**Deploy** button, top right) now that Root Directory, Variables, and the Pre-Deploy Command are all set. Watch the **Deploy Logs** — you should see the pre-deploy container run `./migrate && ./seed` and exit 0, then the `api` container start.',
            },
          ],
        },
      ],
    },
    {
      heading: 'Part 7 — Deploy the admin service',
      blocks: [
        { kind: 'h3', text: '7.1 Create the service' },
        {
          kind: 'ol',
          items: [
            {
              text: 'Click **+ New** → **GitHub Repo** → select the same `sentex` repository and branch again. Railway allows multiple services from the same repo.',
            },
            { text: 'Rename the service to `admin`.' },
          ],
        },
        { kind: 'h3', text: '7.2 Point it at the right Dockerfile' },
        {
          kind: 'p',
          text: 'Unlike `api`, the Compose file builds `admin` with `context: .` (the repo root) but `dockerfile: apps/admin/Dockerfile` — it needs the repo root as build context so it can see the shared `packages/` workspace and root `pnpm-lock.yaml`, but the Dockerfile itself lives one level down.',
        },
        {
          kind: 'ol',
          items: [
            { text: 'Go to **Settings** → **Source**.' },
            {
              text: 'Leave **Root Directory** **blank** (i.e. the repo root) — do not set it to `apps/admin`, or the build will lose access to `packages/` and `pnpm-lock.yaml` and fail during `pnpm install`.',
            },
            {
              text: 'Under **Build Configuration**, set **Dockerfile Path** to:',
              code: { language: 'text', code: 'apps/admin/Dockerfile' },
            },
            {
              text: 'Under **Build**, set **Watch Paths** to:',
              code: { language: 'text', code: 'apps/admin/**\npackages/**' },
            },
          ],
        },
        { kind: 'h3', text: '7.3 Set build-time variables' },
        {
          kind: 'p',
          text: 'These map directly to the `build.args` block for `admin` in your Compose file. Because `apps/admin/Dockerfile` declares them with `ARG`, setting them as ordinary **Variables** on this service (not anything special) is enough for Railway to forward them into the Docker build automatically — that’s what "Using variables at build time" means in Railway’s Dockerfile docs, as long as the `ARG` line exists in the Dockerfile.',
        },
        {
          kind: 'code',
          language: 'bash',
          code: `NEXT_PUBLIC_API_URL=https://api.sentex.gritcms.com
NEXT_PUBLIC_WEB_URL=https://sentex.gritcms.com
NEXT_PUBLIC_ADMIN_URL=https://admin.sentex.gritcms.com
THEME=atlas
SOCIAL_AUTH_ENABLED=false
NEXT_PUBLIC_DEMO_LOGINS=false`,
        },
        {
          kind: 'p',
          text: 'Adjust `THEME`, `SOCIAL_AUTH_ENABLED`, and `NEXT_PUBLIC_DEMO_LOGINS` to your actual production values — the ones above just mirror the safe defaults called out in the Compose file’s comments (social auth off, demo logins off, since this stack has no OAuth provider configured). Remember these are **build-time**: changing them later requires a new deploy/rebuild, not just a service restart, exactly as your Compose file’s comments warn.',
        },
        { kind: 'h3', text: '7.4 Runtime variables and networking' },
        {
          kind: 'ol',
          items: [
            {
              text: 'Add any runtime-only variables `admin` needs beyond the build args (for example, session secrets, if the admin app reads them at request time rather than build time).',
            },
            {
              text: 'Under **Settings → Networking**, click **Generate Domain** for a temporary test URL, or skip straight to Part 8 for the real domain.',
            },
            {
              text: 'Deploy the service and check **Deploy Logs** for a successful Next.js start on port 3000.',
            },
          ],
        },
      ],
    },
    {
      heading: 'Part 8 — Deploy the web service',
      blocks: [
        { kind: 'p', text: 'Repeat Part 7 with these differences:' },
        {
          kind: 'ol',
          items: [
            { text: 'Rename the service to `web`.' },
            { text: '**Root Directory**: blank (repo root) — same reasoning as `admin`.' },
            { text: '**Dockerfile Path**: `apps/web/Dockerfile`.' },
            {
              text: '**Watch Paths**:',
              code: { language: 'text', code: 'apps/web/**\npackages/**' },
            },
            {
              text: 'Build-time variables (matches the `web` service’s `build.args` — note it has no `NEXT_PUBLIC_WEB_URL` arg of its own, since `web` *is* the web app):',
              code: {
                language: 'bash',
                code: `NEXT_PUBLIC_API_URL=https://api.sentex.gritcms.com
NEXT_PUBLIC_ADMIN_URL=https://admin.sentex.gritcms.com
THEME=atlas
SOCIAL_AUTH_ENABLED=false`,
              },
            },
            {
              text: 'Generate a domain / deploy, and confirm the Deploy Logs show a successful start on port 3000.',
            },
          ],
        },
      ],
    },
    {
      heading: 'Part 9 — Attach your custom domains',
      blocks: [
        {
          kind: 'p',
          text: 'Do this once `api`, `admin`, and `web` have each successfully deployed at least once.',
        },
        {
          kind: 'ol',
          items: [
            { text: 'Open the `api` service → **Settings → Networking → Custom Domain**.' },
            { text: 'Enter `api.sentex.gritcms.com` and click **Add**.' },
            {
              text: 'Railway shows you a CNAME target (something like `xxxx.up.railway.app`). Copy it.',
            },
            {
              text: 'In your DNS provider, create a **CNAME** record: `api.sentex.gritcms.com → xxxx.up.railway.app` (replacing the A record placeholder from Part 2).',
            },
            {
              text: 'Repeat steps 1–4 for `admin` (`admin.sentex.gritcms.com`) and `web` (`sentex.gritcms.com`) on their respective services.',
            },
            {
              text: 'Wait for DNS to propagate, then confirm each domain in the Railway dashboard shows a green "Active"/verified TLS status. Railway provisions and renews the TLS certificate automatically once DNS resolves correctly — you don’t manage Let’s Encrypt yourself the way the Compose file’s Traefik labels did.',
            },
          ],
        },
        {
          kind: 'note',
          text: '**Important:** if you generated a domain, deployed, and *then* changed `NEXT_PUBLIC_API_URL` / `NEXT_PUBLIC_ADMIN_URL` / `NEXT_PUBLIC_WEB_URL` to match the final custom domains, you must trigger a fresh deploy of `admin` and `web` after adding the custom domains — those values are baked into the JS bundle at build time, so the bundle built against the temporary `*.up.railway.app` URL will not automatically pick up the custom domain.',
        },
      ],
    },
    {
      heading: 'Part 10 — Verify the full deployment',
      blocks: [
        {
          kind: 'ol',
          items: [
            {
              text: '**Postgres / Redis**: open each service’s **Data** tab (Postgres) or **Metrics** tab and confirm they show as running/healthy.',
            },
            {
              text: '**api**: open Deploy Logs, confirm you see the pre-deploy command run `./migrate && ./seed` and exit cleanly, then the API start message. Hit `https://api.sentex.gritcms.com/<your health endpoint>` with `curl` and confirm a 200.',
            },
            {
              text: '**admin**: visit `https://admin.sentex.gritcms.com` in a browser, open dev tools → Network tab, and confirm requests go to `api.sentex.gritcms.com` with no CORS errors in the console.',
            },
            { text: '**web**: visit `https://sentex.gritcms.com` and repeat the same CORS/API check.' },
            {
              text: 'Log in with the demo SACCO credentials seeded by the `migrate`/`seed` step to confirm the database was actually populated.',
            },
          ],
        },
      ],
    },
    {
      heading: 'Part 11 — Ongoing deploys',
      blocks: [
        {
          kind: 'p',
          text: 'Because every service is connected via **GitHub Repo** (not the CLI), the day-to-day workflow is now:',
        },
        {
          kind: 'ol',
          items: [
            { text: 'Push a commit to the connected branch (e.g. `main`).' },
            {
              text: 'Railway’s GitHub webhook triggers a build for every service whose **Watch Paths** matched the changed files.',
            },
            {
              text: 'For `api`, the pre-deploy command reruns automatically on every deploy — safe, since your migration/seed are idempotent.',
            },
            {
              text: 'Watch **Deploy Logs** per service if anything looks off; roll back to a previous deployment from the service’s **Deployments** tab if needed (Railway keeps deployment history and supports one-click rollback).',
            },
          ],
        },
        {
          kind: 'p',
          text: 'You generally won’t need the Railway CLI for this project’s day-to-day deploys — it’s most useful for local `railway run` / `railway variables` debugging against the same environment, or `railway logs` to tail a service without opening the dashboard.',
        },
      ],
    },
    {
      heading: 'Troubleshooting',
      blocks: [
        {
          kind: 'table',
          headers: ['Symptom', 'Likely cause'],
          rows: [
            [
              'admin/web build fails at pnpm install, can’t find lockfile',
              'Root Directory was set to apps/admin/apps/web instead of left blank — it needs the repo root as context.',
            ],
            [
              'Build-time NEXT_PUBLIC_* value doesn’t show up in the shipped app',
              'Missing ARG for that variable in the Dockerfile, or the variable was set after the last build (rebuild it).',
            ],
            [
              'CORS errors in browser console',
              'CORS_ORIGINS on api doesn’t exactly match the live admin/web domains (scheme + host, no trailing slash).',
            ],
            [
              'Pre-deploy command fails with an exec/format error',
              'Missing the sh -c "..." wrapper — required for Dockerfile-based services.',
            ],
            [
              'api can’t reach Postgres/Redis',
              'Reference variables (${{Postgres.PGHOST}}, etc.) weren’t set, or the plugin service was renamed after wiring.',
            ],
            [
              'A push to apps/web triggers an api rebuild too',
              'Watch Paths weren’t set (or were set too broadly) on one of the services.',
            ],
            [
              'Custom domain stuck "pending"',
              'DNS CNAME not propagated yet, or it’s pointed at the wrong Railway target — recheck against Part 9.',
            ],
          ],
        },
      ],
    },
  ],
}

/* ────────────────────────────────────────────────────────────────────────── */

const DOKPLOY: DeploymentGuide = {
  slug: 'dokploy',
  title: 'Deploying Sentex: VPS + Dokploy',
  intro: [
    'Your `docker-compose.prod.yml` is already written for Dokploy (see its header comments), so this is the least amount of translation of the four — mostly server setup, DNS, and pasting environment variables.',
  ],
  sections: [
    {
      heading: '1.1 Provision the VPS',
      blocks: [
        {
          kind: 'ol',
          items: [
            {
              text: 'Spin up a VPS (2 vCPU / 4 GB RAM minimum for four app containers + Postgres + Redis; more if traffic is expected). Ubuntu 22.04/24.04 is the best-supported OS for Dokploy’s install script.',
            },
            {
              text: 'Point your registrar/DNS provider’s records at the VPS’s public IP — same three hosts as before:',
              code: {
                language: 'text',
                code: `sentex.gritcms.com          → <VPS public IP>
admin.sentex.gritcms.com    → <VPS public IP>
api.sentex.gritcms.com      → <VPS public IP>`,
              },
            },
            {
              text: 'These stay **A records** (not CNAMEs) because Dokploy’s Traefik terminates TLS directly on your server’s IP — exactly what the Compose file’s header comments describe.',
            },
            { text: 'SSH into the server as root (or a user with sudo).' },
          ],
        },
      ],
    },
    {
      heading: '1.2 Install Dokploy',
      blocks: [
        {
          kind: 'ol',
          items: [
            {
              text: 'Run the official installer:',
              code: { language: 'bash', code: 'curl -sSL https://dokploy.com/install.sh | sh' },
            },
            {
              text: 'Wait for it to finish — it installs Docker if missing, starts Dokploy’s own containers, and creates the `dokploy-network` Docker network that your Compose file’s `networks.dokploy-network.external: true` expects.',
            },
            {
              text: 'Open `http://<VPS public IP>:3000` in a browser and create your admin account on first load.',
            },
            {
              text: '(Recommended) Under **Settings → Server**, point a domain at the Dokploy dashboard itself and enable HTTPS for it, so you’re not managing infrastructure over plain HTTP long-term. This is separate from your three app domains.',
            },
          ],
        },
      ],
    },
    {
      heading: '1.3 Connect GitHub',
      blocks: [
        {
          kind: 'ol',
          items: [
            { text: 'In Dokploy, go to **Settings → Git Providers → GitHub**.' },
            {
              text: 'Follow the prompts to install the Dokploy GitHub App on your account/org and grant it access to the `sentex` repository. (SSH deploy keys are the alternative if you’d rather not install a GitHub App — see Dokploy’s Providers docs for that flow.)',
            },
          ],
        },
      ],
    },
    {
      heading: '1.4 Create the project and Compose application',
      blocks: [
        {
          kind: 'ol',
          items: [
            { text: 'In the Dokploy dashboard, click **Create Project**, name it `sentex`.' },
            { text: 'Inside the project, click **Create Service → Compose**.' },
            {
              text: 'Under **Source**, choose **GitHub**, select the `sentex` repository and the branch you’re deploying (e.g. `main`).',
            },
            {
              text: 'Set **Compose Path** to:',
              code: { language: 'text', code: 'docker-compose.prod.yml' },
            },
            {
              text: 'Leave the network settings alone — since your file already declares `networks: dokploy-network: external: true` plus its own internal `sentex` bridge network, Dokploy will attach correctly without any extra configuration.',
            },
          ],
        },
      ],
    },
    {
      heading: '1.5 Set environment variables',
      blocks: [
        {
          kind: 'ol',
          items: [
            { text: 'Go to the **Environment** tab of the Compose service.' },
            {
              text: 'Paste in the full contents of your `.env.production` file (one `KEY=VALUE` per line) — Dokploy writes this to a `.env` file next to your compose file on the server and uses it to interpolate every `${VARIABLE}` reference in `docker-compose.prod.yml`, for both build args (`NEXT_PUBLIC_API_URL`, `THEME`, etc.) and runtime environment (`POSTGRES_PASSWORD`, `APP_URL`, `CORS_ORIGINS`, R2 credentials, and so on). This is a direct match for `env_file: [.env]` already declared for `migrate` and `api` in your Compose file.',
            },
            {
              text: 'Double-check `WEB_DOMAIN`, `ADMIN_DOMAIN`, `API_DOMAIN` are set to the exact hosts from step 1.1 — these are baked into the frontend bundles at build time and drive the Traefik `Host()` rules already written into your Compose file’s labels.',
            },
          ],
        },
      ],
    },
    {
      heading: '1.6 Deploy',
      blocks: [
        {
          kind: 'ol',
          items: [
            {
              text: 'Click **Deploy**. Dokploy will:',
              sub: [
                'Build `migrate`, `api`, `admin`, `web` from their Dockerfiles.',
                'Start `postgres` and `redis` and wait for their healthchecks.',
                'Run `migrate` (`./migrate && ./seed`) to completion — this works unmodified because Dokploy runs a real `docker compose up`, and plain Docker Compose (unlike Railway) natively understands `depends_on: condition: service_completed_successfully` and `restart: "no"`. You don’t need any pre-deploy-command workaround here.',
                'Start `api`, then `admin` and `web` once `api` has started.',
              ],
            },
            {
              text: 'Watch the deployment logs in the Dokploy UI. Confirm `migrate` exits 0 before `api` starts.',
            },
            {
              text: 'After roughly 10 seconds, Traefik should finish provisioning Let’s Encrypt certificates for the three `Host()` rules already defined in your Compose labels (`sentex-api`, `sentex-admin`, `sentex-web` routers). No separate "Domains" configuration step is required in the UI, since the labels already declare everything — that’s what the Compose file’s own header comments mean by "routing works" this way.',
            },
          ],
        },
      ],
    },
    {
      heading: '1.7 Verify',
      blocks: [
        {
          kind: 'ol',
          items: [
            { text: 'Visit `https://api.sentex.gritcms.com/<health endpoint>`.' },
            {
              text: 'Visit `https://admin.sentex.gritcms.com` and `https://sentex.gritcms.com` and confirm no CORS errors in the browser console.',
            },
            {
              text: 'Log in with the seeded demo SACCO credentials to confirm `migrate`/`seed` actually ran.',
            },
          ],
        },
      ],
    },
    {
      heading: '1.8 Ongoing deploys',
      blocks: [
        {
          kind: 'ol',
          items: [
            {
              text: 'In the Compose service’s **General** tab, enable the GitHub webhook ("Auto Deploy" / deploy-on-push) so pushes to your branch redeploy automatically — Dokploy re-clones the repo, re-reads the `.env` it generated, and reruns `docker compose up -d --build`.',
            },
            {
              text: '`migrate`/`seed` reruns every deploy; it’s already idempotent per the Compose file’s own comments, so this is safe.',
            },
          ],
        },
      ],
    },
  ],
}

/* ────────────────────────────────────────────────────────────────────────── */

const COOLIFY: DeploymentGuide = {
  slug: 'coolify',
  title: 'Deploying Sentex: VPS + Coolify',
  intro: [
    'Coolify also runs your Compose file close to as-is (it literally runs `docker compose` under the hood), but it manages its own reverse-proxy network and strongly warns against custom Compose networks. You’ll make a small, mechanical edit to `docker-compose.prod.yml` for this platform — everything else (services, builds, volumes, the `migrate` job) stays exactly as written.',
  ],
  sections: [
    {
      heading: '2.1 Provision the VPS and install Coolify',
      blocks: [
        {
          kind: 'ol',
          items: [
            {
              text: 'Spin up a VPS (same sizing guidance as Dokploy). Ubuntu 24.04 or Debian 13 are the best-supported targets.',
            },
            {
              text: 'Point the same three DNS **A records** at this VPS’s IP (a different VPS than Dokploy’s, obviously, if you’re comparing platforms — don’t point both at the same IP at the same time).',
            },
            {
              text: 'SSH in and run the official installer:',
              code: {
                language: 'bash',
                code: 'curl -fsSL https://cdn.coollabs.io/coolify/install.sh | bash',
              },
            },
            {
              text: 'Open `http://<VPS public IP>:8000`, create your admin account, and (recommended) attach a domain + HTTPS to the Coolify dashboard itself under server settings.',
            },
          ],
        },
      ],
    },
    {
      heading: '2.2 Make a Coolify-specific branch/copy of the Compose file',
      blocks: [
        {
          kind: 'p',
          text: 'Create `docker-compose.coolify.yml` (or a `coolify` branch — whatever fits your workflow) with two changes from `docker-compose.prod.yml`:',
        },
        {
          kind: 'ol',
          items: [
            {
              text: '**Remove the `networks:` block at every service**, and **delete the top-level `networks:` section entirely** (both the `sentex` bridge and the `dokploy-network: external: true` reference). Coolify creates its own isolated bridge network per Compose stack and attaches its own Traefik to it automatically — defining custom networks alongside that causes exactly the kind of intermittent 504/unreachable behavior Coolify’s docs specifically warn about, because your containers would sit on two networks at once and Traefik might pick the wrong one.',
            },
            {
              text: '**Remove the `labels:` blocks** (the `traefik.*` labels) from `api`, `admin`, and `web`. You’ll set domains through Coolify’s UI instead in Part 2.4 — its proxy is still Traefik, but its label names/entrypoint names don’t necessarily match Dokploy’s, so hand-rolled labels are more likely to conflict with what Coolify generates than to help.',
            },
          ],
        },
        {
          kind: 'p',
          text: 'Everything else — `build:`, `image:`, `environment:`, `env_file:`, `volumes:`, `depends_on:`, `healthcheck:`, `command:`, `restart:` — stays identical. In particular, **leave the `migrate` service and its `depends_on: condition: service_completed_successfully` exactly as-is** — Coolify runs real Docker Compose, so this ordering guarantee works without any translation, the same as on Dokploy.',
        },
        {
          kind: 'p',
          text: 'Optionally, mark `migrate` as excluded from Coolify’s aggregate healthchecks, since it’s meant to exit rather than stay running:',
        },
        {
          kind: 'code',
          language: 'yaml',
          code: `services:
  migrate:
    exclude_from_hc: true
    # ...rest unchanged`,
        },
      ],
    },
    {
      heading: '2.3 Create the resource in Coolify',
      blocks: [
        {
          kind: 'ol',
          items: [
            {
              text: 'In the Coolify dashboard, create a **Project**, then a new **Resource** inside it.',
            },
            {
              text: 'Choose your Git source (Public Repository, or GitHub App / Deploy Key for a private repo — set up whichever you haven’t already under **Sources**).',
            },
            { text: 'Select the `sentex` repository and branch.' },
            {
              text: 'When prompted for a Build Pack, change it from the Nixpacks default to **Docker Compose**.',
            },
            {
              text: 'Set:',
              sub: [
                '**Base Directory**: `/` (repo root)',
                '**Docker Compose Location**: `docker-compose.coolify.yml` (the file from Part 2.2 — match the exact filename/extension you used)',
              ],
            },
            { text: 'Click **Continue**.' },
          ],
        },
      ],
    },
    {
      heading: '2.4 Domains',
      blocks: [
        {
          kind: 'p',
          text: 'Coolify reads your Compose file’s services and lets you assign a domain to each one directly — no labels needed since you removed them in Part 2.2.',
        },
        {
          kind: 'ol',
          items: [
            {
              text: 'On the resource’s configuration screen, find the domain field for each service and set:',
              sub: [
                '`api` → `https://api.sentex.gritcms.com:8080` (append `:8080` because that’s the *container* port `api` listens on — Coolify’s proxy still serves the public side on the normal HTTPS port; the `:8080` just tells it where to send traffic internally)',
                '`admin` → `https://admin.sentex.gritcms.com:3000`',
                '`web` → `https://sentex.gritcms.com:3000`',
              ],
            },
            {
              text: 'Leave `postgres` and `redis` with **no domain assigned** — without a domain or a `ports:` mapping, Coolify keeps a service private and reachable only over the internal network at `http://postgres:5432` / `http://redis:6379`-style hostnames (i.e., exactly the plain service-name DNS your Compose file’s `POSTGRES_HOST=postgres` / `REDIS_URL` values already assume).',
            },
          ],
        },
      ],
    },
    {
      heading: '2.5 Environment variables',
      blocks: [
        {
          kind: 'ol',
          items: [
            {
              text: 'Coolify auto-detects every `${VARIABLE}` referenced in your Compose file (in `environment:`, `env_file:`-driven values you reference, and `build.args`) and lists them in the resource’s **Environment Variables** tab.',
            },
            {
              text: 'Fill in the same values you’d put in `.env.production`: `POSTGRES_USER`, `POSTGRES_PASSWORD`, `POSTGRES_DB`, `WEB_DOMAIN`, `ADMIN_DOMAIN`, `API_DOMAIN`, `THEME`, `SOCIAL_AUTH_ENABLED`, `NEXT_PUBLIC_DEMO_LOGINS`, the `R2_*` credentials, and any app secrets.',
            },
            {
              text: 'Coolify injects these both as build args (for `admin`/`web`’s `NEXT_PUBLIC_*` values) and as runtime environment — matching how the Compose file already declares them.',
            },
          ],
        },
      ],
    },
    {
      heading: '2.6 Deploy and verify',
      blocks: [
        {
          kind: 'ol',
          items: [
            { text: 'Click **Deploy**. Watch the build/deploy log stream in the UI.' },
            {
              text: 'Confirm `migrate` runs and exits cleanly before `api`, `admin`, and `web` start (same ordering as Dokploy — real Compose semantics).',
            },
            {
              text: 'Give Coolify’s Traefik a short moment to issue Let’s Encrypt certs for the three domains, then visit each in a browser and confirm no CORS errors and that the seeded demo login works.',
            },
            {
              text: 'Under the resource’s **Webhooks/Source** settings, confirm auto-deploy on push is enabled if you want GitHub pushes to redeploy automatically.',
            },
          ],
        },
      ],
    },
  ],
}

/* ────────────────────────────────────────────────────────────────────────── */

const RENDER: DeploymentGuide = {
  slug: 'render',
  title: 'Deploying Sentex: Render',
  intro: [
    'Render does **not** run `docker-compose.yml` files. Like Railway, it needs an explicit config — Render’s version is a `render.yaml` **Blueprint** committed to the repo, which becomes the single source of truth for every service, database, and env var. This part is closer to the Railway guide than to Dokploy/Coolify.',
  ],
  sections: [
    {
      heading: '3.1 Write the Blueprint',
      blocks: [
        { kind: 'p', text: 'Create `render.yaml` at the repo root:' },
        {
          kind: 'code',
          language: 'yaml',
          code: `databases:
  - name: sentex-postgres
    databaseName: sentex
    user: sentex
    plan: starter
    region: oregon

services:
  - type: keyvalue
    name: sentex-redis
    plan: starter
    region: oregon
    ipAllowList: [] # private — reachable only from services in this workspace

  - type: web
    name: sentex-api
    runtime: docker
    region: oregon
    plan: starter
    dockerfilePath: apps/api/Dockerfile
    dockerContext: apps/api
    domains:
      - api.sentex.gritcms.com
    preDeployCommand: sh -c "./migrate && ./seed"
    envVars:
      - key: APP_ENV
        value: production
      - key: POSTGRES_HOST
        fromDatabase: { name: sentex-postgres, property: host }
      - key: POSTGRES_PORT
        fromDatabase: { name: sentex-postgres, property: port }
      - key: POSTGRES_USER
        fromDatabase: { name: sentex-postgres, property: user }
      - key: POSTGRES_PASSWORD
        fromDatabase: { name: sentex-postgres, property: password }
      - key: POSTGRES_DB
        fromDatabase: { name: sentex-postgres, property: database }
      - key: REDIS_URL
        fromService: { type: keyvalue, name: sentex-redis, property: connectionString }
      - key: APP_URL
        value: https://api.sentex.gritcms.com
      - key: CORS_ORIGINS
        value: https://sentex.gritcms.com,https://admin.sentex.gritcms.com
      - key: R2_ACCOUNT_ID
        sync: false
      - key: R2_ACCESS_KEY_ID
        sync: false
      - key: R2_SECRET_ACCESS_KEY
        sync: false
      - key: R2_BUCKET
        sync: false
      - key: R2_ENDPOINT
        sync: false

  - type: web
    name: sentex-admin
    runtime: docker
    region: oregon
    plan: starter
    dockerfilePath: apps/admin/Dockerfile
    dockerContext: .
    domains:
      - admin.sentex.gritcms.com
    envVars:
      - key: NEXT_PUBLIC_API_URL
        value: https://api.sentex.gritcms.com
      - key: NEXT_PUBLIC_WEB_URL
        value: https://sentex.gritcms.com
      - key: NEXT_PUBLIC_ADMIN_URL
        value: https://admin.sentex.gritcms.com
      - key: THEME
        value: atlas
      - key: SOCIAL_AUTH_ENABLED
        value: "false"
      - key: NEXT_PUBLIC_DEMO_LOGINS
        value: "false"

  - type: web
    name: sentex-web
    runtime: docker
    region: oregon
    plan: starter
    dockerfilePath: apps/web/Dockerfile
    dockerContext: .
    domains:
      - sentex.gritcms.com
    envVars:
      - key: NEXT_PUBLIC_API_URL
        value: https://api.sentex.gritcms.com
      - key: NEXT_PUBLIC_ADMIN_URL
        value: https://admin.sentex.gritcms.com
      - key: THEME
        value: atlas
      - key: SOCIAL_AUTH_ENABLED
        value: "false"`,
        },
        {
          kind: 'p',
          text: 'Notes on the choices above, matching what’s in `docker-compose.prod.yml`:',
        },
        {
          kind: 'ul',
          items: [
            '**`dockerContext`** mirrors the Compose `build.context` for each service: `apps/api` for `api` (matches `context: ./apps/api`), and `.` (repo root) for `admin`/`web` (matches `context: .`, needed for the pnpm workspace). `dockerfilePath` is always relative to the repo root regardless of `dockerContext`.',
            '**`preDeployCommand`** replaces the standalone `migrate` service the same way Railway’s Pre-Deploy Command did — it runs in a fresh instance, before the new `api` version goes live, and aborts the deploy on non-zero exit. Wrap it in `sh -c "..."` so `&&` is interpreted by a shell, same reasoning as the Railway guide.',
            '**`type: keyvalue`** is Render’s current Redis-compatible managed store (runs Valkey; `redis` is a deprecated alias for the same type). Setting `ipAllowList: []` keeps it unreachable from the public internet — only other services in your Render workspace can reach it.',
            '**Build args**: Render automatically forwards a Docker service’s `envVars` into the build as `ARG`s (as long as the Dockerfile declares matching `ARG` lines) — so no separate build-args block is needed for `admin`/`web`’s `NEXT_PUBLIC_*` values, same requirement as in the Railway guide.',
            '**`sync: false`** on the R2 credentials means Render will prompt you to type the actual values into the dashboard the first time the Blueprint syncs, rather than committing secrets into `render.yaml`.',
          ],
        },
        { kind: 'p', text: 'Commit and push `render.yaml`.' },
      ],
    },
    {
      heading: '3.2 Deploy the Blueprint',
      blocks: [
        {
          kind: 'ol',
          items: [
            { text: 'In the Render Dashboard, click **New → Blueprint**.' },
            {
              text: 'Connect your GitHub account if you haven’t, then select the `sentex` repository.',
            },
            {
              text: 'Render detects `render.yaml` at the repo root automatically. Give the Blueprint instance a name and pick the branch to track.',
            },
            {
              text: 'Render shows a preview of every resource it’s about to create (`sentex-postgres`, `sentex-redis`, `sentex-api`, `sentex-admin`, `sentex-web`). Review it, then fill in the prompted values for the `sync: false` variables (your real R2 credentials).',
            },
            { text: 'Click **Deploy Blueprint**.' },
          ],
        },
      ],
    },
    {
      heading: '3.3 Domains',
      blocks: [
        {
          kind: 'ol',
          items: [
            {
              text: 'Because `domains:` is already set per service in `render.yaml`, Render provisions those custom domains automatically as part of the sync — you don’t need a separate manual step to attach them.',
            },
            {
              text: 'Render shows you the CNAME target for each domain (under that service’s **Settings → Custom Domains**). Create the corresponding CNAME records at your DNS provider:',
              code: {
                language: 'text',
                code: `api.sentex.gritcms.com    → CNAME → <target shown by Render>
admin.sentex.gritcms.com  → CNAME → <target shown by Render>
sentex.gritcms.com        → CNAME → <target shown by Render>`,
              },
            },
            {
              text: 'Wait for DNS to propagate and for Render to show each domain as verified with an active TLS certificate.',
            },
          ],
        },
      ],
    },
    {
      heading: '3.4 Verify and iterate',
      blocks: [
        {
          kind: 'ol',
          items: [
            {
              text: 'Check each service’s **Logs** tab — confirm the `sentex-api` pre-deploy log shows `./migrate && ./seed` exiting 0 before the service starts.',
            },
            {
              text: 'Visit all three domains, confirm no CORS errors, and confirm the seeded demo login works.',
            },
            {
              text: 'Because Blueprints auto-redeploy affected services whenever `render.yaml` changes, and each service still auto-deploys on pushes to its watched paths, day-to-day pushes to `main` behave like the Railway and Dokploy/Coolify GitHub-connected flows — no manual redeploy step needed.',
            },
          ],
        },
      ],
    },
  ],
}

/* ────────────────────────────────────────────────────────────────────────── */

const FLY: DeploymentGuide = {
  slug: 'fly-io',
  title: 'Deploying Sentex: Fly.io',
  intro: [
    'Fly.io is the most different of the four: there’s no single "project" that holds multiple services the way Railway/Render/Dokploy/Coolify have one. Every deployable thing is its own **Fly app** with its own `fly.toml`, and you orchestrate the monorepo yourself with `flyctl` flags rather than a platform-level "root directory" setting. Databases are provisioned separately too (Fly Postgres runs as your own VMs; Redis comes from Fly’s built-in Upstash integration).',
  ],
  sections: [
    {
      heading: '4.1 Install flyctl and log in',
      blocks: [
        {
          kind: 'code',
          language: 'bash',
          code: 'curl -L https://fly.io/install.sh | sh\nfly auth login',
        },
      ],
    },
    {
      heading: '4.2 Create three fly.toml files (one per app)',
      blocks: [
        {
          kind: 'p',
          text: 'Fly doesn’t read a monorepo config format — each app gets its own TOML file and you tell `flyctl` which Dockerfile and build context to use for it via flags. Create these at the **repo root** (keeping them there, rather than inside `apps/*`, keeps the build context flexible — see 4.4):',
        },
        { kind: 'h3', text: 'fly.api.toml' },
        {
          kind: 'code',
          language: 'toml',
          code: `app = "sentex-api"
primary_region = "jnb"   # pick the region closest to your users

[build]

[env]
  APP_ENV = "production"

[deploy]
  release_command = "sh -c './migrate && ./seed'"

[http_service]
  internal_port = 8080
  force_https = true
  auto_stop_machines = false
  auto_start_machines = true
  min_machines_running = 1

[[vm]]
  cpu_kind = "shared"
  cpus = 1
  memory_mb = 512`,
        },
        { kind: 'h3', text: 'fly.admin.toml' },
        {
          kind: 'code',
          language: 'toml',
          code: `app = "sentex-admin"
primary_region = "jnb"

[build]
  [build.args]
    NEXT_PUBLIC_API_URL = "https://api.sentex.gritcms.com"
    NEXT_PUBLIC_WEB_URL = "https://sentex.gritcms.com"
    NEXT_PUBLIC_ADMIN_URL = "https://admin.sentex.gritcms.com"
    THEME = "atlas"
    SOCIAL_AUTH_ENABLED = "false"
    NEXT_PUBLIC_DEMO_LOGINS = "false"

[http_service]
  internal_port = 3000
  force_https = true
  auto_stop_machines = false
  auto_start_machines = true
  min_machines_running = 1

[[vm]]
  cpu_kind = "shared"
  cpus = 1
  memory_mb = 512`,
        },
        {
          kind: 'p',
          text: '**`fly.web.toml`** — same shape as `fly.admin.toml`, with:',
        },
        { kind: 'code', language: 'toml', code: 'app = "sentex-web"' },
        {
          kind: 'p',
          text: 'and its build args matching the `web` service’s args from the Compose file (no `NEXT_PUBLIC_WEB_URL` — `web` doesn’t need its own URL as a build arg, same as in the Compose file):',
        },
        {
          kind: 'code',
          language: 'toml',
          code: `[build]
  [build.args]
    NEXT_PUBLIC_API_URL = "https://api.sentex.gritcms.com"
    NEXT_PUBLIC_ADMIN_URL = "https://admin.sentex.gritcms.com"
    THEME = "atlas"
    SOCIAL_AUTH_ENABLED = "false"`,
        },
        { kind: 'p', text: 'Notes:' },
        {
          kind: 'ul',
          items: [
            '`release_command` on `sentex-api` is Fly’s equivalent of Railway’s Pre-Deploy Command and Render’s `preDeployCommand` — it spins up a temporary Machine using the freshly built image, runs `./migrate && ./seed`, and only proceeds to deploy the real release if it exits 0. Same `sh -c` wrapping reasoning as the other platforms.',
            '`[build.args]` is how Fly forwards **Docker build arguments** — this maps directly to the `build.args` block for `admin`/`web` in `docker-compose.prod.yml`. Build args aren’t available at runtime, so this only covers the `NEXT_PUBLIC_*`/`THEME`/`SOCIAL_AUTH_ENABLED` values that Next.js needs baked into the bundle — exactly like the Compose file’s own comments describe.',
            'Runtime secrets (Postgres/Redis URLs, R2 credentials, `APP_URL`, `CORS_ORIGINS`) are **not** put in `fly.toml` — they go in via `fly secrets set` in Part 4.5, the same way Fly handles all sensitive runtime config.',
          ],
        },
      ],
    },
    {
      heading: '4.3 Provision Postgres and Redis',
      blocks: [
        {
          kind: 'code',
          language: 'bash',
          code: `# Postgres — runs as Fly Machines you own, not a separate managed product
fly postgres create --name sentex-postgres --region jnb

# Attach it to the api app — this auto-creates a DATABASE_URL secret on sentex-api
fly postgres attach sentex-postgres --app sentex-api

# Redis — Fly's built-in Upstash-managed integration
fly redis create
# When prompted: name it sentex-redis, pick the same region (jnb), and
# choose whether to enable eviction based on your caching needs.`,
        },
        {
          kind: 'p',
          text: '`fly postgres attach` sets a `DATABASE_URL` secret directly on `sentex-api` automatically. Since your app code expects discrete `POSTGRES_HOST`/`POSTGRES_PORT`/`POSTGRES_USER`/`POSTGRES_PASSWORD`/`POSTGRES_DB` variables rather than a single DSN, parse `DATABASE_URL` into those five values in application startup code, **or** set the five secrets explicitly from the credentials Fly prints when it creates the cluster (Part 4.5 shows the explicit-secrets approach, which needs no code change).',
        },
        { kind: 'p', text: '`fly redis create` prints a `redis://` connection string — copy it for Part 4.5.' },
      ],
    },
    {
      heading: '4.4 Deploy each app',
      blocks: [
        {
          kind: 'p',
          text: 'Because `admin` and `web` need the repo root as build context (same reason as every other platform in this doc — pnpm workspace access) while `api`’s context is just `apps/api`, deploy each with an explicit working directory and `--dockerfile`/`--config` pair:',
        },
        {
          kind: 'code',
          language: 'bash',
          code: `# api — context is apps/api, matching build.context: ./apps/api in Compose
fly deploy apps/api --config fly.api.toml --dockerfile apps/api/Dockerfile

# admin — context is the repo root, matching build.context: . in Compose
fly deploy . --config fly.admin.toml --dockerfile apps/admin/Dockerfile

# web — same reasoning as admin
fly deploy . --config fly.web.toml --dockerfile apps/web/Dockerfile`,
        },
        {
          kind: 'p',
          text: 'The first argument to `fly deploy` is the build context sent to Docker; the `--config` flag tells `flyctl` which `fly.toml` (and therefore which app) you mean, and `--dockerfile` overrides the default `<context>/Dockerfile` lookup.',
        },
        {
          kind: 'p',
          text: 'The first `fly deploy apps/api ...` call will prompt to create the `sentex-api` app if it doesn’t exist yet (since `fly.toml` names it but you haven’t run `fly launch` interactively) — accept the prompt, or run `fly apps create sentex-api` (and the equivalent for `admin`/`web`) ahead of time if you’d rather do it explicitly.',
        },
      ],
    },
    {
      heading: '4.5 Set runtime secrets',
      blocks: [
        {
          kind: 'code',
          language: 'bash',
          code: `fly secrets set -a sentex-api \\
  POSTGRES_HOST="sentex-postgres.flycast" \\
  POSTGRES_PORT="5432" \\
  POSTGRES_USER="<from fly postgres create output>" \\
  POSTGRES_PASSWORD="<from fly postgres create output>" \\
  POSTGRES_DB="<your database name>" \\
  REDIS_URL="<redis:// URL from fly redis create>" \\
  APP_URL="https://api.sentex.gritcms.com" \\
  CORS_ORIGINS="https://sentex.gritcms.com,https://admin.sentex.gritcms.com" \\
  R2_ACCOUNT_ID="..." \\
  R2_ACCESS_KEY_ID="..." \\
  R2_SECRET_ACCESS_KEY="..." \\
  R2_BUCKET="..." \\
  R2_ENDPOINT="..."`,
        },
        {
          kind: 'p',
          text: 'Setting a secret triggers a new deploy of `sentex-api` automatically (which also reruns `release_command`, safe since it’s idempotent). `admin` and `web` don’t need runtime secrets in this stack — their configuration is entirely build-time (`[build.args]` in step 4.2).',
        },
      ],
    },
    {
      heading: '4.6 Attach custom domains',
      blocks: [
        {
          kind: 'code',
          language: 'bash',
          code: `fly certs add api.sentex.gritcms.com -a sentex-api
fly certs add admin.sentex.gritcms.com -a sentex-admin
fly certs add sentex.gritcms.com -a sentex-web`,
        },
        {
          kind: 'p',
          text: 'Each command prints the DNS record(s) to create — typically an **A/AAAA** pair pointing at Fly’s anycast IPs, or a CNAME depending on whether the hostname is a root domain or subdomain. Create those records at your DNS provider, then poll status until issued:',
        },
        {
          kind: 'code',
          language: 'bash',
          code: 'fly certs check api.sentex.gritcms.com -a sentex-api',
        },
      ],
    },
    {
      heading: '4.7 Verify',
      blocks: [
        {
          kind: 'ol',
          items: [
            {
              text: '`fly logs -a sentex-api` — confirm the release command ran `./migrate && ./seed` successfully before the app started serving.',
            },
            {
              text: 'Visit all three domains, confirm no CORS errors, confirm the seeded demo login works.',
            },
            {
              text: '`fly status -a sentex-api` / `-a sentex-admin` / `-a sentex-web` — confirm each shows healthy running Machines.',
            },
          ],
        },
      ],
    },
    {
      heading: '4.8 Ongoing deploys (GitHub Actions)',
      blocks: [
        {
          kind: 'p',
          text: 'Fly has no native "connect a GitHub repo and auto-deploy on push" toggle the way the other three platforms do — the standard pattern is a GitHub Actions workflow that calls `flyctl deploy` for each app:',
        },
        {
          kind: 'ol',
          items: [
            {
              text: 'Generate a deploy token: `fly tokens create deploy -a sentex-api` (and the same for `sentex-admin`, `sentex-web`, or one org-wide token if you prefer — see `fly tokens create org`).',
            },
            { text: 'Add it as a repo secret, e.g. `FLY_API_TOKEN`.' },
            {
              text: 'Add `.github/workflows/fly-deploy.yml`:',
              code: {
                language: 'yaml',
                code: `name: Deploy to Fly.io
on:
  push:
    branches: [main]
jobs:
  deploy-api:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: superfly/flyctl-actions/setup-flyctl@master
      - run: flyctl deploy apps/api --config fly.api.toml --dockerfile apps/api/Dockerfile --remote-only
        env:
          FLY_API_TOKEN: \${{ secrets.FLY_API_TOKEN }}
  deploy-admin:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: superfly/flyctl-actions/setup-flyctl@master
      - run: flyctl deploy . --config fly.admin.toml --dockerfile apps/admin/Dockerfile --remote-only
        env:
          FLY_API_TOKEN: \${{ secrets.FLY_API_TOKEN }}
  deploy-web:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: superfly/flyctl-actions/setup-flyctl@master
      - run: flyctl deploy . --config fly.web.toml --dockerfile apps/web/Dockerfile --remote-only
        env:
          FLY_API_TOKEN: \${{ secrets.FLY_API_TOKEN }}`,
              },
            },
            {
              text: 'Optionally add `paths:` filters per job (mirroring the Watch Paths concept from the Railway guide) so a change under `apps/web/**` doesn’t trigger an unnecessary `sentex-api` rebuild.',
            },
          ],
        },
      ],
    },
  ],
}

/* ────────────────────────────────────────────────────────────────────────── */

/* ────────────────────────────────────────────────────────────────────────── */

export const DEPLOYMENT_GUIDES: DeploymentGuide[] = [
  ORBITA,
  RAILWAY,
  DOKPLOY,
  COOLIFY,
  RENDER,
  FLY,
  AWS_EC2,
  AWS_LIGHTSAIL,
]

export function getGuide(slug: string) {
  return DEPLOYMENT_GUIDES.find((g) => g.slug === slug)
}

/**
 * The comparison table, carried over from the guides.
 *
 * It uses the Orbita guide's version, which is the superset: it adds Orbita's
 * two paths to the four platforms the alternatives guide compared. The EC2 and
 * Lightsail rows follow from their own guides — both run real Docker Compose on
 * a VM, so `depends_on` ordering works unchanged, and both build the TLS layer
 * themselves rather than getting one from a platform.
 *
 * These four columns are the questions that actually differ between the eight.
 * Everything else is detail you can absorb after choosing.
 */
export const GUIDE_COMPARISON = {
  headers: ['', 'Runs your Compose file?', 'Migration handling', 'Domain/TLS', 'GitHub auto-deploy'],
  rows: [
    [
      'Orbita — Grit fast path',
      'No — derives everything from grit.json; you write ~10 lines of orbita.yaml',
      'cmd/migrate under a Postgres advisory lock, gates cutover',
      'All three hosts from one deploy; Traefik + Let’s Encrypt, derived',
      'Built-in webhook (push to redeploy)',
    ],
    [
      'Orbita — Compose path',
      'Yes, almost unmodified (no stripping networks/labels)',
      'Native depends_on: (unchanged)',
      'One routable web service per app; Traefik + Let’s Encrypt',
      'Built-in webhook',
    ],
    [
      'Dokploy',
      'Yes, almost unmodified',
      'Native depends_on: service_completed_successfully',
      'Traefik labels already in the file',
      'Built-in webhook toggle',
    ],
    [
      'Coolify',
      'Yes, after removing custom networks:/labels',
      'Native depends_on: (unchanged)',
      'UI Domains field per service',
      'Built-in webhook toggle',
    ],
    [
      'Render',
      'No — translated to render.yaml',
      'preDeployCommand on api',
      'domains: in Blueprint + CNAME',
      'Built-in (Blueprint sync + per-service auto-deploy)',
    ],
    [
      'Fly.io',
      'No — one fly.toml per app',
      'release_command on api',
      'fly certs add + DNS record',
      'Manual — GitHub Actions workflow',
    ],
    [
      'Railway',
      'No — one service per Compose service',
      'Pre-Deploy Command on api',
      'Custom Domain + CNAME',
      'Built-in (GitHub Repo connection)',
    ],
    [
      'AWS EC2',
      'Yes, after removing labels and adding host ports',
      'Native depends_on: (unchanged)',
      'ALB + ACM certificate, host-header rules',
      'Manual — GitHub Actions over SSH',
    ],
    [
      'AWS Lightsail',
      'Yes, after removing labels and adding a Caddy service',
      'Native depends_on: (unchanged)',
      'Caddy on the instance, automatic Let’s Encrypt',
      'Manual — GitHub Actions over SSH',
    ],
  ],
}
