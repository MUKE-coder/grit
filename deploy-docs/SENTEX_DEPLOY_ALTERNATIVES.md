# Deploying Sentex: Dokploy, Coolify, Render & Fly.io

Four self-contained guides, all against the same `docker-compose.prod.yml` and
monorepo layout (`apps/api`, `apps/admin`, `apps/web`, Postgres, Redis,
`migrate` job) used in the Railway guide. Two of these platforms
(**Dokploy**, **Coolify**) run your Compose file almost as-is; the other two
(**Render**, **Fly.io**) don't run Compose at all and need a translated
config, the same way Railway did. Jump to the section you need — each is
independent and repeats what it needs from the Compose file.

---

# Part 1 — VPS + Dokploy

Your `docker-compose.prod.yml` is already written for Dokploy (see its
header comments), so this is the least amount of translation of the four —
mostly server setup, DNS, and pasting environment variables.

## 1.1 Provision the VPS

1. Spin up a VPS (2 vCPU / 4 GB RAM minimum for four app containers +
   Postgres + Redis; more if traffic is expected). Ubuntu 22.04/24.04 is the
   best-supported OS for Dokploy's install script.
2. Point your registrar/DNS provider's records at the VPS's public IP —
   same three hosts as before:
   ```
   sentex.gritcms.com          → <VPS public IP>
   admin.sentex.gritcms.com    → <VPS public IP>
   api.sentex.gritcms.com      → <VPS public IP>
   ```
   These stay **A records** (not CNAMEs) because Dokploy's Traefik
   terminates TLS directly on your server's IP — exactly what the Compose
   file's header comments describe.
3. SSH into the server as root (or a user with sudo).

## 1.2 Install Dokploy

1. Run the official installer:
   ```
   curl -sSL https://dokploy.com/install.sh | sh
   ```
2. Wait for it to finish — it installs Docker if missing, starts Dokploy's
   own containers, and creates the `dokploy-network` Docker network that
   your Compose file's `networks.dokploy-network.external: true` expects.
3. Open `http://<VPS public IP>:3000` in a browser and create your admin
   account on first load.
4. (Recommended) Under **Settings → Server**, point a domain at the Dokploy
   dashboard itself and enable HTTPS for it, so you're not managing
   infrastructure over plain HTTP long-term. This is separate from your
   three app domains.

## 1.3 Connect GitHub

1. In Dokploy, go to **Settings → Git Providers → GitHub**.
2. Follow the prompts to install the Dokploy GitHub App on your account/org
   and grant it access to the `sentex` repository. (SSH deploy keys are the
   alternative if you'd rather not install a GitHub App — see Dokploy's
   Providers docs for that flow.)

## 1.4 Create the project and Compose application

1. In the Dokploy dashboard, click **Create Project**, name it `sentex`.
2. Inside the project, click **Create Service → Compose**.
3. Under **Source**, choose **GitHub**, select the `sentex` repository and
   the branch you're deploying (e.g. `main`).
4. Set **Compose Path** to:
   ```
   docker-compose.prod.yml
   ```
5. Leave the network settings alone — since your file already declares
   `networks: dokploy-network: external: true` plus its own internal
   `sentex` bridge network, Dokploy will attach correctly without any
   extra configuration.

## 1.5 Set environment variables

1. Go to the **Environment** tab of the Compose service.
2. Paste in the full contents of your `.env.production` file (one
   `KEY=VALUE` per line) — Dokploy writes this to a `.env` file next to
   your compose file on the server and uses it to interpolate every
   `${VARIABLE}` reference in `docker-compose.prod.yml`, for both build
   args (`NEXT_PUBLIC_API_URL`, `THEME`, etc.) and runtime environment
   (`POSTGRES_PASSWORD`, `APP_URL`, `CORS_ORIGINS`, R2 credentials, and so
   on). This is a direct match for `env_file: [.env]` already declared for
   `migrate` and `api` in your Compose file.
3. Double-check `WEB_DOMAIN`, `ADMIN_DOMAIN`, `API_DOMAIN` are set to the
   exact hosts from step 1.1 — these are baked into the frontend bundles at
   build time and drive the Traefik `Host()` rules already written into
   your Compose file's labels.

## 1.6 Deploy

1. Click **Deploy**. Dokploy will:
   - Build `migrate`, `api`, `admin`, `web` from their Dockerfiles.
   - Start `postgres` and `redis` and wait for their healthchecks.
   - Run `migrate` (`./migrate && ./seed`) to completion — this works
     unmodified because Dokploy runs a real `docker compose up`, and plain
     Docker Compose (unlike Railway) natively understands
     `depends_on: condition: service_completed_successfully` and
     `restart: "no"`. You don't need any pre-deploy-command workaround
     here.
   - Start `api`, then `admin` and `web` once `api` has started.
2. Watch the deployment logs in the Dokploy UI. Confirm `migrate` exits 0
   before `api` starts.
3. After roughly 10 seconds, Traefik should finish provisioning Let's
   Encrypt certificates for the three `Host()` rules already defined in
   your Compose labels (`sentex-api`, `sentex-admin`, `sentex-web`
   routers). No separate "Domains" configuration step is required in the
   UI, since the labels already declare everything — that's what the
   Compose file's own header comments mean by "routing works" this way.

## 1.7 Verify

1. Visit `https://api.sentex.gritcms.com/<health endpoint>`.
2. Visit `https://admin.sentex.gritcms.com` and `https://sentex.gritcms.com`
   and confirm no CORS errors in the browser console.
3. Log in with the seeded demo SACCO credentials to confirm `migrate`/`seed`
   actually ran.

## 1.8 Ongoing deploys

1. In the Compose service's **General** tab, enable the GitHub webhook
   ("Auto Deploy" / deploy-on-push) so pushes to your branch redeploy
   automatically — Dokploy re-clones the repo, re-reads the `.env` it
   generated, and reruns `docker compose up -d --build`.
2. `migrate`/`seed` reruns every deploy; it's already idempotent per the
   Compose file's own comments, so this is safe.

---

# Part 2 — VPS + Coolify

Coolify also runs your Compose file close to as-is (it literally runs
`docker compose` under the hood), but it manages its own reverse-proxy
network and strongly warns against custom Compose networks. You'll make a
small, mechanical edit to `docker-compose.prod.yml` for this platform —
everything else (services, builds, volumes, the `migrate` job) stays
exactly as written.

## 2.1 Provision the VPS and install Coolify

1. Spin up a VPS (same sizing guidance as Dokploy). Ubuntu 24.04 or Debian
   13 are the best-supported targets.
2. Point the same three DNS **A records** at this VPS's IP (a different VPS
   than Dokploy's, obviously, if you're comparing platforms — don't point
   both at the same IP at the same time).
3. SSH in and run the official installer:
   ```
   curl -fsSL https://cdn.coollabs.io/coolify/install.sh | bash
   ```
4. Open `http://<VPS public IP>:8000`, create your admin account, and
   (recommended) attach a domain + HTTPS to the Coolify dashboard itself
   under server settings.

## 2.2 Make a Coolify-specific branch/copy of the Compose file

Create `docker-compose.coolify.yml` (or a `coolify` branch — whatever fits
your workflow) with two changes from `docker-compose.prod.yml`:

1. **Remove the `networks:` block at every service**, and **delete the
   top-level `networks:` section entirely** (both the `sentex` bridge and
   the `dokploy-network: external: true` reference). Coolify creates its
   own isolated bridge network per Compose stack and attaches its own
   Traefik to it automatically — defining custom networks alongside that
   causes exactly the kind of intermittent 504/unreachable behavior
   Coolify's docs specifically warn about, because your containers would
   sit on two networks at once and Traefik might pick the wrong one.
2. **Remove the `labels:` blocks** (the `traefik.*` labels) from `api`,
   `admin`, and `web`. You'll set domains through Coolify's UI instead in
   Part 2.4 — its proxy is still Traefik, but its label names/entrypoint
   names don't necessarily match Dokploy's, so hand-rolled labels are more
   likely to conflict with what Coolify generates than to help.

Everything else — `build:`, `image:`, `environment:`, `env_file:`,
`volumes:`, `depends_on:`, `healthcheck:`, `command:`, `restart:` — stays
identical. In particular, **leave the `migrate` service and its
`depends_on: condition: service_completed_successfully` exactly as-is** —
Coolify runs real Docker Compose, so this ordering guarantee works without
any translation, the same as on Dokploy.

Optionally, mark `migrate` as excluded from Coolify's aggregate
healthchecks, since it's meant to exit rather than stay running:

```yaml
services:
  migrate:
    exclude_from_hc: true
    # ...rest unchanged
```

## 2.3 Create the resource in Coolify

1. In the Coolify dashboard, create a **Project**, then a new **Resource**
   inside it.
2. Choose your Git source (Public Repository, or GitHub App / Deploy Key
   for a private repo — set up whichever you haven't already under
   **Sources**).
3. Select the `sentex` repository and branch.
4. When prompted for a Build Pack, change it from the Nixpacks default to
   **Docker Compose**.
5. Set:
   - **Base Directory**: `/` (repo root)
   - **Docker Compose Location**: `docker-compose.coolify.yml` (the file
     from Part 2.2 — match the exact filename/extension you used)
6. Click **Continue**.

## 2.4 Domains

Coolify reads your Compose file's services and lets you assign a domain to
each one directly — no labels needed since you removed them in Part 2.2.

1. On the resource's configuration screen, find the domain field for each
   service and set:
   - `api` → `https://api.sentex.gritcms.com:8080` (append `:8080` because
     that's the *container* port `api` listens on — Coolify's proxy still
     serves the public side on the normal HTTPS port; the `:8080` just
     tells it where to send traffic internally)
   - `admin` → `https://admin.sentex.gritcms.com:3000`
   - `web` → `https://sentex.gritcms.com:3000`
2. Leave `postgres` and `redis` with **no domain assigned** — without a
   domain or a `ports:` mapping, Coolify keeps a service private and
   reachable only over the internal network at `http://postgres:5432` /
   `http://redis:6379`-style hostnames (i.e., exactly the plain service-name
   DNS your Compose file's `POSTGRES_HOST=postgres` / `REDIS_URL` values
   already assume).

## 2.5 Environment variables

1. Coolify auto-detects every `${VARIABLE}` referenced in your Compose file
   (in `environment:`, `env_file:`-driven values you reference, and
   `build.args`) and lists them in the resource's **Environment Variables**
   tab.
2. Fill in the same values you'd put in `.env.production`: `POSTGRES_USER`,
   `POSTGRES_PASSWORD`, `POSTGRES_DB`, `WEB_DOMAIN`, `ADMIN_DOMAIN`,
   `API_DOMAIN`, `THEME`, `SOCIAL_AUTH_ENABLED`, `NEXT_PUBLIC_DEMO_LOGINS`,
   the `R2_*` credentials, and any app secrets.
3. Coolify injects these both as build args (for `admin`/`web`'s
   `NEXT_PUBLIC_*` values) and as runtime environment — matching how the
   Compose file already declares them.

## 2.6 Deploy and verify

1. Click **Deploy**. Watch the build/deploy log stream in the UI.
2. Confirm `migrate` runs and exits cleanly before `api`, `admin`, and
   `web` start (same ordering as Dokploy — real Compose semantics).
3. Give Coolify's Traefik a short moment to issue Let's Encrypt certs for
   the three domains, then visit each in a browser and confirm no CORS
   errors and that the seeded demo login works.
4. Under the resource's **Webhooks/Source** settings, confirm auto-deploy
   on push is enabled if you want GitHub pushes to redeploy automatically.

---

# Part 3 — Render

Render does **not** run `docker-compose.yml` files. Like Railway, it needs
an explicit config — Render's version is a `render.yaml` **Blueprint**
committed to the repo, which becomes the single source of truth for every
service, database, and env var. This part is closer to the Railway guide
than to Dokploy/Coolify.

## 3.1 Write the Blueprint

Create `render.yaml` at the repo root:

```yaml
databases:
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
        value: "false"
```

Notes on the choices above, matching what's in `docker-compose.prod.yml`:

- **`dockerContext`** mirrors the Compose `build.context` for each service:
  `apps/api` for `api` (matches `context: ./apps/api`), and `.` (repo root)
  for `admin`/`web` (matches `context: .`, needed for the pnpm workspace).
  `dockerfilePath` is always relative to the repo root regardless of
  `dockerContext`.
- **`preDeployCommand`** replaces the standalone `migrate` service the same
  way Railway's Pre-Deploy Command did — it runs in a fresh instance,
  before the new `api` version goes live, and aborts the deploy on
  non-zero exit. Wrap it in `sh -c "..."` so `&&` is interpreted by a
  shell, same reasoning as the Railway guide.
- **`type: keyvalue`** is Render's current Redis-compatible managed store
  (runs Valkey; `redis` is a deprecated alias for the same type). Setting
  `ipAllowList: []` keeps it unreachable from the public internet — only
  other services in your Render workspace can reach it.
- **Build args**: Render automatically forwards a Docker service's
  `envVars` into the build as `ARG`s (as long as the Dockerfile declares
  matching `ARG` lines) — so no separate build-args block is needed for
  `admin`/`web`'s `NEXT_PUBLIC_*` values, same requirement as in the
  Railway guide.
- **`sync: false`** on the R2 credentials means Render will prompt you to
  type the actual values into the dashboard the first time the Blueprint
  syncs, rather than committing secrets into `render.yaml`.

Commit and push `render.yaml`.

## 3.2 Deploy the Blueprint

1. In the Render Dashboard, click **New → Blueprint**.
2. Connect your GitHub account if you haven't, then select the `sentex`
   repository.
3. Render detects `render.yaml` at the repo root automatically. Give the
   Blueprint instance a name and pick the branch to track.
4. Render shows a preview of every resource it's about to create
   (`sentex-postgres`, `sentex-redis`, `sentex-api`, `sentex-admin`,
   `sentex-web`). Review it, then fill in the prompted values for the
   `sync: false` variables (your real R2 credentials).
5. Click **Deploy Blueprint**.

## 3.3 Domains

1. Because `domains:` is already set per service in `render.yaml`, Render
   provisions those custom domains automatically as part of the sync — you
   don't need a separate manual step to attach them.
2. Render shows you the CNAME target for each domain (under that service's
   **Settings → Custom Domains**). Create the corresponding CNAME records
   at your DNS provider:
   ```
   api.sentex.gritcms.com    → CNAME → <target shown by Render>
   admin.sentex.gritcms.com  → CNAME → <target shown by Render>
   sentex.gritcms.com        → CNAME → <target shown by Render>
   ```
3. Wait for DNS to propagate and for Render to show each domain as
   verified with an active TLS certificate.

## 3.4 Verify and iterate

1. Check each service's **Logs** tab — confirm the `sentex-api` pre-deploy
   log shows `./migrate && ./seed` exiting 0 before the service starts.
2. Visit all three domains, confirm no CORS errors, and confirm the seeded
   demo login works.
3. Because Blueprints auto-redeploy affected services whenever
   `render.yaml` changes, and each service still auto-deploys on pushes to
   its watched paths, day-to-day pushes to `main` behave like the Railway
   and Dokploy/Coolify GitHub-connected flows — no manual redeploy step
   needed.

---

# Part 4 — Fly.io

Fly.io is the most different of the four: there's no single "project" that
holds multiple services the way Railway/Render/Dokploy/Coolify have one.
Every deployable thing is its own **Fly app** with its own `fly.toml`, and
you orchestrate the monorepo yourself with `flyctl` flags rather than a
platform-level "root directory" setting. Databases are provisioned
separately too (Fly Postgres runs as your own VMs; Redis comes from Fly's
built-in Upstash integration).

## 4.1 Install flyctl and log in

```
curl -L https://fly.io/install.sh | sh
fly auth login
```

## 4.2 Create three `fly.toml` files (one per app)

Fly doesn't read a monorepo config format — each app gets its own TOML file
and you tell `flyctl` which Dockerfile and build context to use for it via
flags. Create these at the **repo root** (keeping them there, rather than
inside `apps/*`, keeps the build context flexible — see 4.4):

**`fly.api.toml`**
```toml
app = "sentex-api"
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
  memory_mb = 512
```

**`fly.admin.toml`**
```toml
app = "sentex-admin"
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
  memory_mb = 512
```

**`fly.web.toml`** — same shape as `fly.admin.toml`, with:
```toml
app = "sentex-web"
```
and its build args matching the `web` service's args from the Compose file
(no `NEXT_PUBLIC_WEB_URL` — `web` doesn't need its own URL as a build arg,
same as in the Compose file):
```toml
[build]
  [build.args]
    NEXT_PUBLIC_API_URL = "https://api.sentex.gritcms.com"
    NEXT_PUBLIC_ADMIN_URL = "https://admin.sentex.gritcms.com"
    THEME = "atlas"
    SOCIAL_AUTH_ENABLED = "false"
```

Notes:

- `release_command` on `sentex-api` is Fly's equivalent of Railway's
  Pre-Deploy Command and Render's `preDeployCommand` — it spins up a
  temporary Machine using the freshly built image, runs
  `./migrate && ./seed`, and only proceeds to deploy the real release if
  it exits 0. Same `sh -c` wrapping reasoning as the other platforms.
- `[build.args]` is how Fly forwards **Docker build arguments** — this maps
  directly to the `build.args` block for `admin`/`web` in
  `docker-compose.prod.yml`. Build args aren't available at runtime, so
  this only covers the `NEXT_PUBLIC_*`/`THEME`/`SOCIAL_AUTH_ENABLED`
  values that Next.js needs baked into the bundle — exactly like the
  Compose file's own comments describe.
- Runtime secrets (Postgres/Redis URLs, R2 credentials, `APP_URL`,
  `CORS_ORIGINS`) are **not** put in `fly.toml` — they go in via
  `fly secrets set` in Part 4.5, the same way Fly handles all sensitive
  runtime config.

## 4.3 Provision Postgres and Redis

```
# Postgres — runs as Fly Machines you own, not a separate managed product
fly postgres create --name sentex-postgres --region jnb

# Attach it to the api app — this auto-creates a DATABASE_URL secret on sentex-api
fly postgres attach sentex-postgres --app sentex-api

# Redis — Fly's built-in Upstash-managed integration
fly redis create
# When prompted: name it sentex-redis, pick the same region (jnb), and
# choose whether to enable eviction based on your caching needs.
```

`fly postgres attach` sets a `DATABASE_URL` secret directly on `sentex-api`
automatically. Since your app code expects discrete
`POSTGRES_HOST`/`POSTGRES_PORT`/`POSTGRES_USER`/`POSTGRES_PASSWORD`/`POSTGRES_DB`
variables rather than a single DSN, parse `DATABASE_URL` into those five
values in application startup code, **or** set the five secrets explicitly
from the credentials Fly prints when it creates the cluster (Part 4.5
shows the explicit-secrets approach, which needs no code change).

`fly redis create` prints a `redis://` connection string — copy it for
Part 4.5.

## 4.4 Deploy each app

Because `admin` and `web` need the repo root as build context (same reason
as every other platform in this doc — pnpm workspace access) while `api`'s
context is just `apps/api`, deploy each with an explicit working directory
and `--dockerfile`/`--config` pair:

```
# api — context is apps/api, matching build.context: ./apps/api in Compose
fly deploy apps/api --config fly.api.toml --dockerfile apps/api/Dockerfile

# admin — context is the repo root, matching build.context: . in Compose
fly deploy . --config fly.admin.toml --dockerfile apps/admin/Dockerfile

# web — same reasoning as admin
fly deploy . --config fly.web.toml --dockerfile apps/web/Dockerfile
```

The first argument to `fly deploy` is the build context sent to Docker; the
`--config` flag tells `flyctl` which `fly.toml` (and therefore which app)
you mean, and `--dockerfile` overrides the default `<context>/Dockerfile`
lookup.

The first `fly deploy apps/api ...` call will prompt to create the
`sentex-api` app if it doesn't exist yet (since `fly.toml` names it but you
haven't run `fly launch` interactively) — accept the prompt, or run
`fly apps create sentex-api` (and the equivalent for `admin`/`web`) ahead
of time if you'd rather do it explicitly.

## 4.5 Set runtime secrets

```
fly secrets set -a sentex-api \
  POSTGRES_HOST="sentex-postgres.flycast" \
  POSTGRES_PORT="5432" \
  POSTGRES_USER="<from fly postgres create output>" \
  POSTGRES_PASSWORD="<from fly postgres create output>" \
  POSTGRES_DB="<your database name>" \
  REDIS_URL="<redis:// URL from fly redis create>" \
  APP_URL="https://api.sentex.gritcms.com" \
  CORS_ORIGINS="https://sentex.gritcms.com,https://admin.sentex.gritcms.com" \
  R2_ACCOUNT_ID="..." \
  R2_ACCESS_KEY_ID="..." \
  R2_SECRET_ACCESS_KEY="..." \
  R2_BUCKET="..." \
  R2_ENDPOINT="..."
```

Setting a secret triggers a new deploy of `sentex-api` automatically (which
also reruns `release_command`, safe since it's idempotent). `admin` and
`web` don't need runtime secrets in this stack — their configuration is
entirely build-time (`[build.args]` in step 4.2).

## 4.6 Attach custom domains

```
fly certs add api.sentex.gritcms.com -a sentex-api
fly certs add admin.sentex.gritcms.com -a sentex-admin
fly certs add sentex.gritcms.com -a sentex-web
```

Each command prints the DNS record(s) to create — typically an **A/AAAA**
pair pointing at Fly's anycast IPs, or a CNAME depending on whether the
hostname is a root domain or subdomain. Create those records at your DNS
provider, then poll status until issued:

```
fly certs check api.sentex.gritcms.com -a sentex-api
```

## 4.7 Verify

1. `fly logs -a sentex-api` — confirm the release command ran
   `./migrate && ./seed` successfully before the app started serving.
2. Visit all three domains, confirm no CORS errors, confirm the seeded
   demo login works.
3. `fly status -a sentex-api` / `-a sentex-admin` / `-a sentex-web` —
   confirm each shows healthy running Machines.

## 4.8 Ongoing deploys (GitHub Actions)

Fly has no native "connect a GitHub repo and auto-deploy on push" toggle
the way the other three platforms do — the standard pattern is a GitHub
Actions workflow that calls `flyctl deploy` for each app:

1. Generate a deploy token: `fly tokens create deploy -a sentex-api` (and
   the same for `sentex-admin`, `sentex-web`, or one org-wide token if you
   prefer — see `fly tokens create org`).
2. Add it as a repo secret, e.g. `FLY_API_TOKEN`.
3. Add `.github/workflows/fly-deploy.yml`:
   ```yaml
   name: Deploy to Fly.io
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
             FLY_API_TOKEN: ${{ secrets.FLY_API_TOKEN }}
     deploy-admin:
       runs-on: ubuntu-latest
       steps:
         - uses: actions/checkout@v4
         - uses: superfly/flyctl-actions/setup-flyctl@master
         - run: flyctl deploy . --config fly.admin.toml --dockerfile apps/admin/Dockerfile --remote-only
           env:
             FLY_API_TOKEN: ${{ secrets.FLY_API_TOKEN }}
     deploy-web:
       runs-on: ubuntu-latest
       steps:
         - uses: actions/checkout@v4
         - uses: superfly/flyctl-actions/setup-flyctl@master
         - run: flyctl deploy . --config fly.web.toml --dockerfile apps/web/Dockerfile --remote-only
           env:
             FLY_API_TOKEN: ${{ secrets.FLY_API_TOKEN }}
   ```
4. Optionally add `paths:` filters per job (mirroring the Watch Paths
   concept from the Railway guide) so a change under `apps/web/**` doesn't
   trigger an unnecessary `sentex-api` rebuild.

---

# Quick comparison

| | Runs your Compose file? | Migration handling | Domain/TLS | GitHub auto-deploy |
|---|---|---|---|---|
| **Dokploy** | Yes, almost unmodified | Native `depends_on: service_completed_successfully` | Traefik labels already in the file | Built-in webhook toggle |
| **Coolify** | Yes, after removing custom `networks:`/labels | Native `depends_on:` (unchanged) | UI **Domains** field per service | Built-in webhook toggle |
| **Render** | No — translated to `render.yaml` | `preDeployCommand` on `api` | `domains:` in Blueprint + CNAME | Built-in (Blueprint sync + per-service auto-deploy) |
| **Fly.io** | No — one `fly.toml` per app | `release_command` on `api` | `fly certs add` + DNS record | Manual — GitHub Actions workflow |
