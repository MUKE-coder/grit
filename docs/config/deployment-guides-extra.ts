import type { DeploymentGuide } from './deployment-guides'

/**
 * The Orbita, EC2 and Lightsail walkthroughs.
 *
 * Split out of deployment-guides.ts purely for file size — the guides are long
 * transcriptions and one module holding all eight is unpleasant to open. The
 * types live in deployment-guides.ts and are imported type-only here, so this
 * file has no runtime dependency on that one and the import graph stays a
 * straight line rather than a cycle.
 *
 * Same rule as the sibling file: transcribed, not summarised. If a step is not
 * in the source guide it does not belong here.
 */

export const ORBITA: DeploymentGuide = {
  slug: 'orbita',
  title: 'Deploying Sentex on Orbita',
  intro: [
    'A self-contained guide for deploying Sentex to **Orbita** — the self-hosted, multi-tenant PaaS built by the same team as Grit. It targets the same monorepo layout (`apps/api`, `apps/admin`, `apps/web`, Postgres, Redis, `migrate` job) and the same three hosts used in the Railway, Dokploy, Coolify, Render and Fly.io guides.',
    'Orbita is different from the other four platforms in one important way: because Sentex is a **Grit app**, Orbita recognises it and deploys the whole thing with almost no configuration. It reads `grit.json`, works out that this is a three-service app, builds the Dockerfiles Grit already ships, provisions Postgres and Redis for you, runs your migrations under a lock, and wires up all three domains — from a manifest that’s about ten lines long. You don’t translate your Compose file the way Render and Fly.io needed, and you don’t hand-edit it the way Coolify needed.',
    'You install Orbita **once** (Part 1), then deploy Sentex to it. There are two ways to deploy, and they produce the same result:',
    '**Part 2 — the Grit fast path** (recommended). Orbita derives everything from `grit.json`. This is the path that makes Orbita worth using for a Grit app.',
    '**Part 3 — the Docker Compose path** (parity with Dokploy/Coolify). Keep `docker-compose.prod.yml` as your source of truth and let Orbita run it as a Swarm stack. Use this only if you specifically want the Compose file to stay authoritative.',
    'Jump to the part you need. Part 1 is a prerequisite for both.',
  ],
  sections: [
    {
      heading: 'Part 1 — Stand up Orbita (once)',
      blocks: [
        {
          kind: 'p',
          text: 'This installs Orbita on a VPS: it becomes your dashboard, your build server, and your reverse proxy. You do this once, then deploy as many apps (and as many tenants) as you like onto it. If you already have an Orbita server running, skip to Part 2.',
        },
        {
          kind: 'p',
          text: 'The install brings **Docker, Docker Swarm, PostgreSQL, Redis, and Traefik** with it — there is nothing to set up beforehand.',
        },
        { kind: 'h3', text: '1.1 Provision the VPS' },
        {
          kind: 'ol',
          items: [
            {
              text: 'Spin up a VPS. **2 vCPU / 4 GB RAM** is a sensible minimum — Orbita itself is tiny (~50 MB idle), but *builds* are hungry, and you’ll be building three Next.js/Go services plus running Postgres and Redis. 4 GB gives the builds room; go higher if you expect real traffic or multiple tenants.',
            },
            {
              text: 'Use a **fresh Ubuntu 22.04 or 24.04** box. Orbita’s Traefik needs to own ports **80** and **443**, so don’t install it on a server that already runs nginx, Apache, or another panel — the installer will stop if it finds one.',
            },
            {
              text: 'Note the server’s **public IP** and the **root password** (or the SSH key) the provider gave you. You’ll need them in 1.3.',
            },
          ],
        },
        { kind: 'p', text: 'Any of Contabo, Hetzner, DigitalOcean, or Vultr work well.' },
        { kind: 'h3', text: '1.2 Point DNS at the server' },
        {
          kind: 'p',
          text: 'Orbita’s Traefik terminates TLS directly on your server’s IP, so these are **A records** (not CNAMEs) — exactly like Dokploy in Part 1 of the alternatives guide. You need **four** hosts: one for the Orbita dashboard itself, and the three Sentex app hosts.',
        },
        {
          kind: 'code',
          language: 'text',
          code: `orbita.gritcms.com          → <VPS public IP>     # the Orbita dashboard
sentex.gritcms.com          → <VPS public IP>     # web
admin.sentex.gritcms.com    → <VPS public IP>     # admin
api.sentex.gritcms.com      → <VPS public IP>     # api`,
        },
        {
          kind: 'note',
          text: 'If your DNS is behind Cloudflare, set each record to **DNS only** (grey cloud), not Proxied — Orbita fetches the Let’s Encrypt certificate itself, and the orange-cloud proxy gets in the way of that first handshake. You can turn it back on afterwards.',
        },
        {
          kind: 'p',
          text: 'Confirm the dashboard record resolves before installing — certificates can’t be issued until it does:',
        },
        {
          kind: 'code',
          language: 'bash',
          code: 'dig orbita.gritcms.com +short      # must print your VPS IP',
        },
        { kind: 'h3', text: '1.3 Harden the server' },
        {
          kind: 'p',
          text: 'Never skip this. A fresh VPS with root SSH open is a target within minutes. SSH in as root, then run the hardening script:',
        },
        {
          kind: 'code',
          language: 'bash',
          code: `ssh root@<VPS public IP>

curl -sSL https://raw.githubusercontent.com/MUKE-coder/vps-harden/main/vps-harden.sh -o vps-harden.sh
chmod +x vps-harden.sh
sudo ./vps-harden.sh --no-dokploy`,
        },
        { kind: 'p', text: 'It asks a few plain questions:' },
        {
          kind: 'ul',
          items: [
            '**A username** for your everyday account — type `deploy`.',
            '**An SSH port** — press Enter to keep the default.',
            '**Your SSH public key** — paste your `~/.ssh/id_ed25519.pub`, or leave it blank and the script generates a key and tells you where it saved it.',
            '**A password for the account** — leave it blank and the script generates a strong one and **prints it once at the end**. Save that password.',
          ],
        },
        {
          kind: 'p',
          text: 'When it finishes you have a `deploy` user with its own password *and* passwordless `sudo`, your key installed, root and password SSH logins disabled, a firewall (UFW + ufw-docker), Fail2ban, kernel hardening, and a 0–100 security score. `--no-dokploy` tells it not to install Dokploy, since Orbita is your platform here.',
        },
        {
          kind: 'note',
          text: 'Before you disconnect: open a **second terminal** and confirm `ssh deploy@<VPS public IP>` works, so you don’t lock yourself out.',
        },
        { kind: 'h3', text: '1.4 Install Orbita' },
        {
          kind: 'p',
          text: 'From here you’re logged in as `deploy` (the `deploy` user reaches Docker through `sudo` — it’s deliberately not in the `docker` group, which is root-equivalent). Run the one-line installer, passing your dashboard domain and an email for Let’s Encrypt:',
        },
        {
          kind: 'code',
          language: 'bash',
          code: `curl -sSL https://raw.githubusercontent.com/MUKE-coder/orbita/main/install.sh \\
  | sudo ORBITA_DOMAIN=orbita.gritcms.com ORBITA_ACME_EMAIL=you@gritcms.com bash -s -- --yes`,
        },
        {
          kind: 'p',
          text: '(If you skipped the domain in 1.2 and want to trial on the IP, drop the two env vars and Orbita comes up on `http://<VPS public IP>:8080` with no TLS.)',
        },
        {
          kind: 'p',
          text: 'In order, the installer installs Docker and starts Swarm, checks ports 80/443/8080, generates secrets into `/opt/orbita/.env`, pulls the Orbita image, starts all four services (`orbita`, `orbita-postgres`, `orbita-redis`, `orbita-traefik`), opens the firewall for the ports it needs, and waits for a healthy `/health` before printing your dashboard URL.',
        },
        { kind: 'p', text: 'Verify:' },
        {
          kind: 'code',
          language: 'bash',
          code: `curl -s http://localhost:8080/health          # want {"status":"ok", ...}
cd /opt/orbita && sudo docker compose ps       # all four services "Up"`,
        },
        {
          kind: 'note',
          text: '**Back up `/opt/orbita/.env`.** It holds `ENCRYPTION_MASTER_KEY`, from which every organisation’s encryption key is derived. Lose it and every stored secret (including the Sentex env you’re about to upload) is unrecoverable. Copy it somewhere safe before you put real data in.',
        },
        { kind: 'h3', text: '1.5 Create your super-admin and organisation' },
        { kind: 'p', text: 'Open the dashboard:' },
        {
          kind: 'ul',
          items: [
            '**With a domain:** `https://orbita.gritcms.com`',
            '**IP only:** `http://<VPS public IP>:8080`',
          ],
        },
        {
          kind: 'p',
          text: 'Click **Register** and create your account **immediately** — the first person to register becomes the **super-admin** with full control of the box. Once that account exists, public sign-up **closes automatically**; nobody else can walk in and register. (Later teammates join by invitation, or through an account you create for them under **Admin**.)',
        },
        {
          kind: 'p',
          text: 'Then create an **organisation** — your top-level workspace. Everything lives inside one, and each org is fully isolated: its own Docker network, its own encryption key, its own resource quota. Name it something like `gritcms` (this becomes the org **slug**, used to namespace networks and volumes).',
        },
        {
          kind: 'p',
          text: 'If you’re running Sentex for a client and want to hand *them* the org, use **Admin → Organisations → New tenant** instead: it creates the org, sizes it (CPU/RAM/disk/app limits), and creates the client’s login in one step, showing you a generated password to hand over. They set their own password at first sign-in. For deploying your own app, a plain organisation is fine.',
        },
        { kind: 'h3', text: '1.6 Connect GitHub' },
        {
          kind: 'p',
          text: 'Sentex is a private repo, so Orbita needs a token to clone it and to register the auto-deploy webhook.',
        },
        {
          kind: 'ol',
          items: [
            { text: 'In the dashboard, go to **Settings → Git Connections**.' },
            {
              text: 'Add a **GitHub** connection with a Personal Access Token that has the **`repo`** and **`admin:repo_hook`** scopes. `repo` lets Orbita clone the private `sentex` repository; `admin:repo_hook` lets it install the push-to-deploy webhook so future commits redeploy automatically.',
            },
          ],
        },
        { kind: 'p', text: 'That’s the whole platform set up. Everything below is per-app.' },
      ],
    },
    {
      heading: 'Part 2 — Deploy Sentex the Grit way (recommended)',
      blocks: [
        {
          kind: 'p',
          text: 'This is the path that makes Orbita worth using for a Grit app. You write a short `orbita.yaml`, and Orbita derives the rest from `grit.json`.',
        },
        { kind: 'h3', text: '2.1 What Orbita derives from grit.json (you write none of this)' },
        {
          kind: 'p',
          text: 'A Grit app has a known shape, declared in `grit.json` at the repo root. Sentex has `apps/api`, `apps/admin`, and `apps/web`, which is Grit’s **triple** architecture. From that single fact, Orbita works out the entire deployment — this is the table you would otherwise have hand-written as a Compose file, `render.yaml`, or three `fly.toml` files in the other guides:',
        },
        {
          kind: 'table',
          headers: ['From grit.json', 'Orbita derives for Sentex'],
          rows: [
            [
              '`architecture: triple`',
              'Three containers: `api`, `web`, `admin` (plus `docs` if the repo has it)',
            ],
            [
              'The Dockerfiles Grit ships',
              'Builds `api` from `apps/api`, and the Next.js apps from the repo root — the exact `build.context` values you set by hand in `docker-compose.prod.yml`',
            ],
            [
              'Ports',
              '`8080` for the API, `3000` for the Next.js apps — no `:port` mapping to write',
            ],
            [
              '`NEXT_PUBLIC_API_URL`',
              'Baked into the `admin` and `web` bundles at build time from your **api** domain — the single most error-prone value in every other guide, derived here',
            ],
          ],
        },
        {
          kind: 'p',
          text: 'Orbita does **not** generate a Dockerfile and does **not** fall back to Nixpacks for a Grit app — it reuses the correct multi-stage Dockerfiles Grit already ships, the same ones the other platforms build.',
        },
        { kind: 'h3', text: '2.2 Write orbita.yaml' },
        {
          kind: 'p',
          text: 'Create `orbita.yaml` at the Sentex repo root. This is the whole deploy config — compare it to the ~90-line `docker-compose.prod.yml`, the `render.yaml` Blueprint, or the three `fly.toml` files:',
        },
        {
          kind: 'code',
          language: 'yaml',
          code: `app: sentex
repo: <your-org>/sentex          # GitHub owner/name of the Sentex repo
branch: main

addons:                          # provisioned in this org's isolated network;
  - postgres                     #   connection URLs injected into the app env
  - redis

domains:
  web:   sentex.gritcms.com
  admin: admin.sentex.gritcms.com
  api:   api.sentex.gritcms.com

migrate: true                    # run the migrations under an advisory lock (default true)

env:
  from: .env.production          # local file; values are encrypted into Orbita, never committed`,
        },
        { kind: 'p', text: 'Notes on the choices, mapped to `docker-compose.prod.yml`:' },
        {
          kind: 'ul',
          items: [
            '**`addons: [postgres, redis]`** replaces the `postgres` and `redis` services you declared in Compose. Orbita provisions managed instances inside your org’s private Docker network and injects their connection URLs (e.g. `DATABASE_URL`, `REDIS_URL`) into the app — the values Grit’s code reads. You don’t set `POSTGRES_HOST` / `POSTGRES_PASSWORD` / `REDIS_URL` yourself; that’s the point. (If Sentex uses Orbita’s object storage instead of Cloudflare R2, add `minio` to the list and Orbita injects `MINIO_*` / `STORAGE_DRIVER` too. If you’re keeping R2, leave `minio` off and put the `R2_*` values in `.env.production` — see 2.3.)',
            '**`domains`** replaces the three Traefik `Host()` labels in your Compose file. Bare hostnames only — no scheme, port, or path. Orbita creates the routers and fetches Let’s Encrypt certs for all three.',
            '**`migrate: true`** replaces the standalone `migrate` service and its `depends_on: condition: service_completed_successfully`. Orbita runs `cmd/migrate` in a one-off container **before cutover**, under a Postgres advisory lock. This is stronger than the Compose ordering guarantee: a non-zero exit **aborts the deploy** and the previous version keeps serving (see 2.8).',
          ],
        },
        {
          kind: 'p',
          text: 'Two optional Grit toggles, both **on by default**, so you only add them to turn something off:',
        },
        {
          kind: 'code',
          language: 'yaml',
          code: `observability: true    # Pulse    — latency/SQL/error tracing on the API
security: true         # Sentinel — WAF, rate limiting, anomaly detection on the API
studio: false          # GORM Studio — off by default; it edits live data`,
        },
        { kind: 'h3', text: '2.3 Provide the environment values' },
        {
          kind: 'p',
          text: 'Everything that isn’t derived comes from `.env.production` — the same file you’d paste into Dokploy or fill into Coolify, minus the Postgres/Redis values Orbita now supplies. That means the build-time and app-secret values Sentex needs:',
        },
        {
          kind: 'code',
          language: 'bash',
          code: `THEME=atlas
SOCIAL_AUTH_ENABLED=false
NEXT_PUBLIC_DEMO_LOGINS=false

# Cloudflare R2 (only if you're keeping R2 instead of Orbita's minio addon)
R2_ACCOUNT_ID=...
R2_ACCESS_KEY_ID=...
R2_SECRET_ACCESS_KEY=...
R2_BUCKET=...
R2_ENDPOINT=...

# ...any other app secrets Sentex reads at runtime`,
        },
        {
          kind: 'p',
          text: 'You do **not** put `NEXT_PUBLIC_API_URL`, `DATABASE_URL`, `REDIS_URL`, `POSTGRES_*`, or the domain variables here — Orbita derives those from your `domains` and `addons`. When Orbita reads `env.from`, it encrypts every value into the org’s key at rest; the file is never committed and never leaves your machine in plaintext.',
        },
        {
          kind: 'p',
          text: 'Now pick a route: **2.4** (dashboard) or **2.5** (CLI). They do the same thing.',
        },
        { kind: 'h3', text: '2.4 Route A — deploy from the dashboard' },
        { kind: 'p', text: 'No CLI needed; everything happens in the browser.' },
        {
          kind: 'ol',
          items: [
            {
              text: 'In your organisation, create a **Project** (e.g. *Sentex*) and an **Environment** (e.g. *production*) inside it. Apps live under project → environment.',
            },
            {
              text: 'Click **Create App → Source: Git Repository**. Pick the GitHub connection from 1.6, then the `sentex` repo and the `main` branch.',
            },
            {
              text: 'Because the repo has a `grit.json`, Orbita recognises it as a Grit app and uses the fast path — you don’t choose a builder or a Dockerfile. Confirm the three derived domains match your DNS from 1.2: `sentex.gritcms.com`, `admin.sentex.gritcms.com`, `api.sentex.gritcms.com`.',
            },
            {
              text: 'Open the app’s **Environment** tab and paste the contents of `.env.production` (one `KEY=VALUE` per line). Mark the `R2_*` keys and any secrets as **secret** so they’re encrypted at rest and never shown again.',
            },
            { text: 'Click **Deploy**.' },
          ],
        },
        {
          kind: 'p',
          text: 'Watch the **Deployments** tab. Orbita builds all three services, provisions Postgres and Redis, runs the migration, and cuts over only if it succeeds. Skip to 2.6 for what you’re watching.',
        },
        { kind: 'h3', text: '2.5 Route B — deploy with the CLI' },
        {
          kind: 'p',
          text: 'The `orbita` CLI is optional, but it’s the tidiest way to deploy from your machine and keep `orbita.yaml` as the source of truth. Today it’s **built from source** — there’s no `curl | sh` installer yet, and `go install` doesn’t work because of the module path — so build it once from a clone (requires Go 1.25+):',
        },
        {
          kind: 'code',
          language: 'bash',
          code: `git clone https://github.com/MUKE-coder/orbita.git
cd orbita
make build-cli
sudo mv ./orbita /usr/local/bin/orbita
orbita --help`,
        },
        {
          kind: 'p',
          text: 'It won’t clash with Grit’s own `grit` binary — different repo, different name.',
        },
        {
          kind: 'p',
          text: 'Then, from the **Sentex** project directory (the one with `grit.json` and the `orbita.yaml` you wrote in 2.2):',
        },
        {
          kind: 'code',
          language: 'bash',
          code: `# Register your server with the CLI (once). Prompts for the admin email +
# password you created in 1.5, mints a deploy token, saves the host as "prod".
orbita login https://orbita.gritcms.com

# Store a GitHub token (repo + admin:repo_hook) so Orbita can push/clone (once).
orbita github-auth

# Preview the plan without changing anything — highly recommended first run.
orbita deploy --plan --host prod`,
        },
        {
          kind: 'p',
          text: 'The plan prints exactly what it will create, so you can confirm the mode and domains before anything happens:',
        },
        {
          kind: 'code',
          language: 'text',
          code: `▸ Plan (dry run — nothing will be changed)
  App:       sentex
  Mode:      triple
  Migrate:   true
  Addons:    postgres, redis

  create  sentex-api    → api.sentex.gritcms.com
  create  sentex-web    → sentex.gritcms.com
  create  sentex-admin  → admin.sentex.gritcms.com`,
        },
        { kind: 'p', text: 'When it looks right, deploy for real:' },
        { kind: 'code', language: 'bash', code: 'orbita deploy --host prod' },
        {
          kind: 'p',
          text: '(If you don’t yet have an Orbita server at all, `orbita init` collapses all of Part 1 — harden, install, admin account, host registration — into one interactive command from your machine. Use it *instead of* Part 1, not as well.)',
        },
        { kind: 'h3', text: '2.6 What happens, in order' },
        {
          kind: 'p',
          text: 'Whichever route you used, a deploy runs these steps — this is the pipeline the other four platforms make you assemble by hand:',
        },
        {
          kind: 'ol',
          items: [
            {
              text: '**Detect** — `grit.json` at the repo root marks it a Grit app; `architecture: triple` picks the three-service strategy.',
            },
            {
              text: '**Ensure the repo** — Orbita confirms it can reach the `sentex` repo with your token (and, over the CLI, pushes your current commit).',
            },
            {
              text: '**Reconcile** — org, project, environment, the `postgres` + `redis` addons, your encrypted env, and the three domains. Idempotent — safe to re-run.',
            },
            {
              text: '**Build** — `api`, `web`, and `admin` from the Dockerfiles Grit ships, with `NEXT_PUBLIC_API_URL` baked into the two Next.js bundles from your api domain.',
            },
            {
              text: '**Migrate** — `cmd/migrate` in a one-off container, under a Postgres advisory lock so two concurrent deploys can’t race.',
            },
            {
              text: '**Cut over** — only if the migration exited 0. The previous images are kept for instant rollback.',
            },
            {
              text: '**Route** — Traefik serves all three domains over HTTPS. Certs are issued on the first request to each host.',
            },
          ],
        },
        { kind: 'h3', text: '2.7 Verify' },
        {
          kind: 'ol',
          items: [
            {
              text: 'Visit `https://api.sentex.gritcms.com/<health endpoint>` and confirm it answers.',
            },
            {
              text: 'Visit `https://sentex.gritcms.com` and `https://admin.sentex.gritcms.com` and confirm there are **no CORS errors** in the browser console — if there are, the API domain baked into the frontend bundle doesn’t match `api`’s real domain; recheck `domains.api` in `orbita.yaml`.',
            },
            {
              text: 'Log in with the seeded demo SACCO credentials to confirm the migration and seed actually ran.',
            },
          ],
        },
        { kind: 'p', text: 'From the CLI you can also stream logs and confirm the migration:' },
        {
          kind: 'code',
          language: 'bash',
          code: `orbita logs -f --host prod                       # all services
orbita logs --host prod --service migrate        # just the migration job`,
        },
        { kind: 'h3', text: '2.8 Migrations gate the cutover (troubleshooting)' },
        {
          kind: 'p',
          text: 'Orbita runs your migrations **before** it cuts over, under an advisory lock. A non-zero exit stops the deploy and leaves the previous version serving — you never end up on a schema-mismatched image. If a deploy fails at the migrate step, that’s why.',
        },
        {
          kind: 'p',
          text: 'The most common cause with a Grit app is **`go.sum` not being committed**, so `go run ./cmd/migrate` can’t resolve modules inside the one-off container. Commit it — real Grit apps ship it — and redeploy. Check the migrate log:',
        },
        { kind: 'code', language: 'bash', code: 'orbita logs --host prod --service migrate' },
        { kind: 'h3', text: '2.9 Batteries included (Pulse, Sentinel, Studio)' },
        {
          kind: 'p',
          text: 'Because Sentex is a Grit app, Orbita mounts these on the API by default — no setup:',
        },
        {
          kind: 'ul',
          items: [
            '**Pulse** — latency, SQL, and error tracing → `https://api.sentex.gritcms.com/pulse/ui`',
            '**Sentinel** — WAF, rate limiting, anomaly detection → `https://api.sentex.gritcms.com/sentinel/ui`',
            '**GORM Studio** — off by default because it edits live data. Turn it on with `studio: true` in `orbita.yaml` only when you need it.',
          ],
        },
        { kind: 'h3', text: '2.10 Ongoing deploys' },
        {
          kind: 'p',
          text: 'Because you connected GitHub in 1.6, Orbita installed a push-to-deploy webhook when it created the app. Every push to `main` now:',
        },
        {
          kind: 'ul',
          items: [
            're-clones the repo and rebuilds the changed services,',
            'reruns `cmd/migrate` under the lock (idempotent, so this is safe every time),',
            'cuts over only if the migration succeeds.',
          ],
        },
        {
          kind: 'p',
          text: 'No manual redeploy step — the same GitHub-connected flow as Dokploy and Coolify. To revert a bad deploy, `orbita rollback --host prod` (or the **Rollback** button on a previous deployment in the dashboard) swaps back to the previous image instantly, since Orbita keeps it.',
        },
      ],
    },
    {
      heading: 'Part 3 — Alternative: run your Compose file on Orbita',
      blocks: [
        {
          kind: 'p',
          text: 'Use this **only** if you specifically want `docker-compose.prod.yml` to stay the source of truth — for example, to keep one Compose file working identically across Dokploy, Coolify, and Orbita. For a Grit app, Part 2 is simpler and gives you migrations-under-a-lock, provisioned addons, and the observability mounts that this path does not. This path treats Sentex as a generic multi-service stack, not as a Grit app.',
        },
        {
          kind: 'p',
          text: 'Orbita deploys a Compose file as a **Docker Swarm stack**. It runs the file essentially as-is — you do **not** strip networks or labels the way Coolify required, and you do **not** translate it to another format the way Render and Fly.io required.',
        },
        { kind: 'h3', text: '3.1 Create the app from Docker Compose' },
        {
          kind: 'ol',
          items: [
            { text: 'In your project/environment, click **Create App → Source: Docker Compose**.' },
            {
              text: 'Point it at the compose file:',
              sub: [
                '**From your Git repo** (recommended, so pushes redeploy): pick the GitHub connection, the `sentex` repo and branch, and set the compose file path to `docker-compose.prod.yml`.',
                '**Or paste it inline** — but note a pasted file can’t use `build:` (there’s no source tree to build from), so it must reference prebuilt images. Sentex builds from source, so use the Git option.',
              ],
            },
            {
              text: 'Set the **web service** to the service that serves your primary domain — `web` for Sentex. This is the service Orbita routes your app domain to; the others stay private to the stack, reachable by their compose service name (`api`, `postgres`, `redis`) exactly as they are locally.',
            },
            {
              text: 'Set the **port** to the web service’s container port — `3000` for the Sentex `web` service. Port is required for Compose apps, because that’s what the domain routes to.',
            },
          ],
        },
        { kind: 'h3', text: '3.2 Domains' },
        {
          kind: 'p',
          text: 'Add your domains under the app’s **Domains** tab. Only the nominated **web service** is routable from a single Compose app, so:',
        },
        {
          kind: 'ul',
          items: [
            'Add `sentex.gritcms.com` → routes to the `web` service you nominated.',
            'To give `api` and `admin` their own domains, the clean approach on this path is to **deploy each as its own app** (three Compose apps, or better, use the Grit fast path in Part 2 which does all three at once). A single Compose app exposes one routable service.',
          ],
        },
        {
          kind: 'p',
          text: 'This is the main reason Part 2 is preferable for a three-domain app like Sentex — the Grit fast path routes all three hosts from one deployment.',
        },
        { kind: 'h3', text: '3.3 Environment' },
        {
          kind: 'p',
          text: 'Open the app’s **Environment** tab and paste your `.env.production` values. Orbita injects them into **every** service in the stack (so a worker gets the same `DATABASE_URL` the web tier does), and they’re encrypted at rest. A service’s own `environment:` block in the compose file still wins if it sets the same key. `${VAR}` references in the compose file are interpolated from these values too, matching `env_file: [.env]`.',
        },
        { kind: 'h3', text: '3.4 Deploy and verify' },
        {
          kind: 'ol',
          items: [
            {
              text: 'Click **Deploy**. Orbita builds the services that declare `build:`, then runs `docker stack deploy` for the whole file.',
            },
            {
              text: 'Watch the deploy log. Because this is a real Swarm deploy of your Compose file, the `migrate` service’s `depends_on: service_completed_successfully` ordering works unmodified — same as Dokploy/Coolify.',
            },
            { text: 'Visit your web domain and confirm it serves.' },
          ],
        },
        { kind: 'h3', text: '3.5 Limits worth knowing on the Compose path' },
        {
          kind: 'ul',
          items: [
            '**Only the web service is routable** per Compose app (see 3.2).',
            '**No rollback** for Compose apps — a Compose deploy has no single image to revert to. Redeploy the previous commit instead. (The Grit path in Part 2 *does* support instant rollback.)',
            '**`build:` needs a Git repo** — pasted YAML must use prebuilt images.',
            'Stopping, starting, or deleting the app applies to **every** service in the stack.',
          ],
        },
      ],
    },
    {
      heading: 'Quick comparison',
      blocks: [
        {
          kind: 'p',
          text: 'Why Orbita’s Grit path is the shortest of them all for Sentex: it’s the only one that already knows what a Grit app is. The others need you to describe a three-service app in their own dialect (Compose, Blueprint, or three TOMLs); Orbita reads the same `grit.json` your app already ships and derives the rest — addons, ports, build contexts, the API URL baked into the frontends, and migrations under a lock — from that.',
        },
      ],
    },
  ],
}

export const AWS_EC2: DeploymentGuide = {
  slug: 'aws-ec2',
  title: 'Deploying Sentex on AWS EC2',
  intro: [
    'Unlike Railway/Render/Fly, EC2 gives you a bare Linux box — there’s no platform-level "connect GitHub" or managed proxy. You run `docker-compose.prod.yml` almost exactly as-is (real Docker Compose, so `depends_on: condition: service_completed_successfully` for `migrate` works natively, same as the Dokploy/Coolify VPS guides), and you build the public-facing, TLS-terminating, host-based-routing layer yourself using an AWS **Application Load Balancer (ALB)** + **ACM** — the standard production pattern on AWS, and the reason this guide looks different from the plain-VPS ones.',
  ],
  sections: [
    {
      heading: 'Architecture at a glance',
      blocks: [
        {
          kind: 'code',
          language: 'text',
          code: `Internet → ALB (443, ACM cert, host-header routing) → EC2 instance (3 host ports)
                                                          ├─ :8080 → api container
                                                          ├─ :3001 → admin container
                                                          └─ :3000 → web container
EC2 instance also runs postgres + redis, reachable only inside the
instance's own Docker network (no host port, no ALB route to them).`,
        },
        {
          kind: 'p',
          text: 'The ALB replaces Traefik/Dokploy’s role from the VPS guides: it terminates TLS with an AWS-managed certificate and routes each hostname to a different container port on the same instance.',
        },
      ],
    },
    {
      heading: 'Part 1 — Network and security groups',
      blocks: [
        {
          kind: 'ol',
          items: [
            {
              text: 'In the VPC you’ll deploy into (the default VPC is fine for a first pass), confirm you have **at least two public subnets in two different Availability Zones** — the ALB requires this even though your EC2 instance only runs in one AZ.',
            },
            {
              text: 'Create a security group **`sentex-alb-sg`**:',
              sub: [
                'Inbound: HTTP (80) from `0.0.0.0/0`, HTTPS (443) from `0.0.0.0/0`.',
                'Outbound: all traffic.',
              ],
            },
            {
              text: 'Create a security group **`sentex-ec2-sg`**:',
              sub: [
                'Inbound: SSH (22) from your IP only (not `0.0.0.0/0`).',
                'Inbound: custom TCP `8080`, `3000`, `3001` — **source: `sentex-alb-sg`** (not an IP range). This is what keeps the app containers unreachable from the internet except through the ALB.',
                'Outbound: all traffic.',
              ],
            },
          ],
        },
      ],
    },
    {
      heading: 'Part 2 — Launch the EC2 instance',
      blocks: [
        {
          kind: 'ol',
          items: [
            { text: 'EC2 → **Launch Instance**.' },
            { text: 'AMI: **Ubuntu Server 24.04 LTS**.' },
            {
              text: 'Instance type: **t3.medium** (2 vCPU / 4 GB) minimum — building three Docker images (two of them Next.js/pnpm builds) on `t3.micro`/`t3.small` routinely OOMs. Scale up later once you know your real traffic.',
            },
            { text: 'Key pair: create or select one — you’ll need it for SSH.' },
            {
              text: 'Network settings: the VPC/subnet from Part 1, **auto-assign public IP: enabled** (needed for SSH and for `git`/`docker pull` egress; the app ports themselves stay locked to ALB-only per the security group).',
            },
            { text: 'Security group: attach **`sentex-ec2-sg`**.' },
            {
              text: 'Storage: bump the root volume to at least **40 GB gp3** — three Docker images plus Postgres/Redis data adds up fast.',
            },
            {
              text: 'Launch. Once running, allocate an **Elastic IP** and associate it to the instance, so the public IP doesn’t change on stop/start (useful for SSH convenience; the ALB, not this IP, is what your DNS will point at).',
            },
          ],
        },
      ],
    },
    {
      heading: 'Part 3 — Install Docker on the instance',
      blocks: [
        { kind: 'p', text: 'SSH in:' },
        {
          kind: 'code',
          language: 'bash',
          code: 'ssh -i /path/to/key.pem ubuntu@<instance-public-ip>',
        },
        { kind: 'p', text: 'Install Docker Engine + Compose plugin:' },
        {
          kind: 'code',
          language: 'bash',
          code: `curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh
sudo usermod -aG docker $USER
newgrp docker
docker compose version   # confirm the plugin is present`,
        },
        {
          kind: 'p',
          text: '(Optional but recommended on `t3.medium`) add a swap file so pnpm/Next.js builds don’t get OOM-killed under memory pressure:',
        },
        {
          kind: 'code',
          language: 'bash',
          code: `sudo fallocate -l 2G /swapfile
sudo chmod 600 /swapfile
sudo mkswap /swapfile
sudo swapon /swapfile
echo '/swapfile none swap sw 0 0' | sudo tee -a /etc/fstab`,
        },
      ],
    },
    {
      heading: 'Part 4 — Get the code and adjust the Compose file for a bare VM',
      blocks: [
        {
          kind: 'ol',
          items: [
            {
              text: 'Clone the repo:',
              code: {
                language: 'bash',
                code: 'git clone https://github.com/<you>/sentex.git\ncd sentex',
              },
            },
            {
              text: '`docker-compose.prod.yml` was written for Dokploy’s Traefik + external `dokploy-network`. Neither exists here — the ALB replaces Traefik, so make a `docker-compose.ec2.yml` copy with:',
              sub: [
                '**Remove** the `traefik.*` labels from `api`, `admin`, `web` (nothing is watching for them on a bare instance).',
                '**Remove** the top-level `dokploy-network` entry and drop it from every service’s `networks:` list — keep only the internal `sentex` network.',
                '**Add host port mappings** so the ALB’s target groups have something to hit.',
              ],
              code: {
                language: 'yaml',
                code: `services:
  api:
    ports:
      - "8080:8080"
  admin:
    ports:
      - "3001:3000"
  web:
    ports:
      - "3000:3000"`,
              },
            },
            {
              text: '`postgres` and `redis` get **no** `ports:` entry — they stay reachable only via the internal `sentex` network at hostnames `postgres`/`redis`, exactly as your app’s `POSTGRES_HOST=postgres` / `REDIS_URL` values already assume.',
            },
            {
              text: 'Everything else — `build:`, `env_file:`, `environment:`, `volumes:`, `depends_on:`, `healthcheck:`, the `migrate` service’s `command:` and `restart: "no"` — stays **unchanged**. Real Docker Compose on a real VM honors `service_completed_successfully` natively, so `migrate` still runs to completion before `api` starts, with no pre-deploy-command workaround needed.',
            },
          ],
        },
      ],
    },
    {
      heading: 'Part 5 — Environment variables',
      blocks: [
        {
          kind: 'ol',
          items: [
            {
              text: 'Create `.env` in the repo root on the instance (do **not** commit it):',
              code: { language: 'bash', code: 'nano .env' },
            },
            {
              text: 'Paste the contents of your local `.env.production`, with `WEB_DOMAIN`, `ADMIN_DOMAIN`, `API_DOMAIN` set to the real hostnames:',
              code: {
                language: 'bash',
                code: `WEB_DOMAIN=sentex.gritcms.com
ADMIN_DOMAIN=admin.sentex.gritcms.com
API_DOMAIN=api.sentex.gritcms.com
POSTGRES_USER=...
POSTGRES_PASSWORD=...
POSTGRES_DB=...
THEME=atlas
SOCIAL_AUTH_ENABLED=false
NEXT_PUBLIC_DEMO_LOGINS=false
R2_ACCOUNT_ID=...
R2_ACCESS_KEY_ID=...
R2_SECRET_ACCESS_KEY=...
R2_BUCKET=...
R2_ENDPOINT=...`,
              },
            },
            {
              text: '`docker compose` reads `.env` in the working directory automatically and interpolates every `${VARIABLE}` in the Compose file — same mechanism the file’s own `env_file: [.env]` entries rely on for `api`/`migrate`.',
            },
          ],
        },
      ],
    },
    {
      heading: 'Part 6 — First deploy',
      blocks: [
        {
          kind: 'code',
          language: 'bash',
          code: 'docker compose -f docker-compose.ec2.yml --env-file .env up -d --build',
        },
        { kind: 'p', text: 'Watch it come up:' },
        {
          kind: 'code',
          language: 'bash',
          code: `docker compose -f docker-compose.ec2.yml logs -f migrate
docker compose -f docker-compose.ec2.yml ps`,
        },
        { kind: 'p', text: 'Confirm `migrate` exits 0, then `api`, `admin`, `web` show as running.' },
        { kind: 'p', text: 'Sanity-check locally on the instance before touching DNS/ALB:' },
        {
          kind: 'code',
          language: 'bash',
          code: `curl -I http://localhost:8080/<health endpoint>
curl -I http://localhost:3001
curl -I http://localhost:3000`,
        },
      ],
    },
    {
      heading: 'Part 7 — Request the TLS certificate (ACM)',
      blocks: [
        {
          kind: 'ol',
          items: [
            {
              text: 'Open **AWS Certificate Manager** in the **same region** as your ALB will live.',
            },
            { text: '**Request a certificate** → Public certificate.' },
            {
              text: 'Domain names — add all three as Subject Alternative Names on one certificate:',
              code: {
                language: 'text',
                code: `sentex.gritcms.com
admin.sentex.gritcms.com
api.sentex.gritcms.com`,
              },
            },
            { text: 'Validation method: **DNS validation** (faster and auto-renewing).' },
            {
              text: 'ACM shows a CNAME record per domain. Create these at your DNS provider. If your zone is hosted in Route 53, ACM offers a **Create records in Route 53** button that does this for you.',
            },
            {
              text: 'Wait for all three domains to show **Issued** (usually a few minutes once the CNAMEs resolve).',
            },
          ],
        },
      ],
    },
    {
      heading: 'Part 8 — Target groups',
      blocks: [
        {
          kind: 'p',
          text: 'Create three **target groups**, type **Instances**, protocol **HTTP**, in the same VPC:',
        },
        {
          kind: 'table',
          headers: ['Name', 'Port', 'Health check path'],
          rows: [
            ['sentex-api-tg', '8080', '/<your api health endpoint>'],
            ['sentex-admin-tg', '3001', '/'],
            ['sentex-web-tg', '3000', '/'],
          ],
        },
        {
          kind: 'p',
          text: 'For each: on the **Register targets** step, select your EC2 instance and, importantly, set the **port override** to that target group’s port (they all point at the same instance, just different ports) before clicking **Include as pending below** → **Register pending targets**.',
        },
      ],
    },
    {
      heading: 'Part 9 — Create the Application Load Balancer',
      blocks: [
        {
          kind: 'ol',
          items: [
            {
              text: 'EC2 → **Load Balancers** → **Create load balancer** → **Application Load Balancer**.',
            },
            { text: 'Scheme: **Internet-facing**. IP type: IPv4.' },
            {
              text: 'VPC: same as the instance. Mappings: select the two+ public subnets from Part 1.',
            },
            { text: 'Security group: **`sentex-alb-sg`**.' },
            {
              text: 'Listeners:',
              sub: [
                '**HTTP:80** — you’ll edit this after creation to redirect to HTTPS (step 7 below).',
                '**HTTPS:443** — default certificate: the ACM cert from Part 7. Default action: pick any target group for now (e.g. `sentex-web-tg`) — you’ll override per-hostname with rules next.',
              ],
            },
            { text: 'Create the load balancer and wait for its state to become **Active**.' },
            {
              text: 'On the **HTTP:80** listener, edit the default rule to **Redirect to HTTPS://#{host}:443/#{path}?#{query}** (status 301) instead of forwarding — this makes plain-HTTP requests upgrade automatically.',
            },
            {
              text: 'On the **HTTPS:443** listener, add rules (in order, above the default action):',
              sub: [
                '**IF** Host header is `api.sentex.gritcms.com` → **THEN** forward to `sentex-api-tg`',
                '**IF** Host header is `admin.sentex.gritcms.com` → **THEN** forward to `sentex-admin-tg`',
                '**IF** Host header is `sentex.gritcms.com` → **THEN** forward to `sentex-web-tg`',
                'Default action (no rule matched): forward to `sentex-web-tg`, or return a fixed 404 — your call.',
              ],
            },
          ],
        },
      ],
    },
    {
      heading: 'Part 10 — DNS',
      blocks: [
        {
          kind: 'ol',
          items: [
            {
              text: 'Note the ALB’s DNS name (something like `sentex-alb-123456789.us-east-1.elb.amazonaws.com`), shown on the load balancer’s detail page.',
            },
            {
              text: 'Create three records at your DNS provider:',
              sub: [
                'If hosted in **Route 53**: create **A records with Alias** target = the ALB, for all three hostnames (Alias records work for subdomains and are free of the CNAME-at-apex restriction, and Route 53 resolves them without an extra DNS lookup).',
                'If hosted **elsewhere**: create **CNAME** records for all three hostnames pointing at the ALB’s DNS name (fine here since none of the three is a bare apex domain).',
              ],
            },
            {
              text: 'Wait for propagation, then confirm:',
              code: {
                language: 'bash',
                code: `dig +short api.sentex.gritcms.com
dig +short admin.sentex.gritcms.com
dig +short sentex.gritcms.com`,
              },
            },
          ],
        },
      ],
    },
    {
      heading: 'Part 11 — Verify',
      blocks: [
        {
          kind: 'ol',
          items: [
            {
              text: 'Visit `https://api.sentex.gritcms.com/<health endpoint>`, `https://admin.sentex.gritcms.com`, `https://sentex.gritcms.com` — all should load over a valid ACM certificate with no browser warnings.',
            },
            {
              text: 'Open dev tools on `admin`/`web`, confirm API calls succeed with no CORS errors (double check `CORS_ORIGINS` in `.env` matches exactly).',
            },
            {
              text: 'Log in with the seeded demo credentials to confirm `migrate`/`seed` populated the database.',
            },
            {
              text: 'In the ALB console, confirm all three target groups show their registered target as **healthy**.',
            },
          ],
        },
      ],
    },
    {
      heading: 'Part 12 — Ongoing deploys',
      blocks: [
        {
          kind: 'p',
          text: 'There’s no GitHub webhook wired up out of the box on a raw EC2 box — pick one:',
        },
        { kind: 'p', text: '**Manual** (fine for a solo project):' },
        {
          kind: 'code',
          language: 'bash',
          code: `ssh ubuntu@<instance-ip>
cd sentex && git pull
docker compose -f docker-compose.ec2.yml --env-file .env up -d --build`,
        },
        {
          kind: 'p',
          text: '**GitHub Actions** (push-to-deploy, closer to the other guides’ workflow):',
        },
        {
          kind: 'code',
          language: 'yaml',
          code: `name: Deploy to EC2
on:
  push:
    branches: [main]
jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: appleboy/ssh-action@v1
        with:
          host: \${{ secrets.EC2_HOST }}
          username: ubuntu
          key: \${{ secrets.EC2_SSH_KEY }}
          script: |
            cd sentex
            git pull
            docker compose -f docker-compose.ec2.yml --env-file .env up -d --build`,
        },
        {
          kind: 'p',
          text: 'Store the instance’s SSH private key as `EC2_SSH_KEY` and its Elastic IP as `EC2_HOST` in the repo’s GitHub Actions secrets.',
        },
        {
          kind: 'p',
          text: '`migrate` reruns on every `up -d --build` and is idempotent (per the Compose file’s own comments), so this is safe to run on every deploy.',
        },
      ],
    },
    {
      heading: 'Notes and troubleshooting',
      blocks: [
        {
          kind: 'table',
          headers: ['Symptom', 'Likely cause'],
          rows: [
            [
              'ALB target shows "unhealthy"',
              'Health check path returns non-2xx, or the security group doesn’t allow the ALB SG on that port — recheck Part 1 step 3.',
            ],
            [
              '504 from the ALB',
              'Container isn’t listening on the port the target group expects — confirm with docker compose ps and curl localhost:<port> on the instance.',
            ],
            [
              'Build gets OOM-killed',
              'Instance too small for concurrent Next.js builds — add swap (Part 3) or size up to t3.large temporarily for the first build.',
            ],
            [
              'CORS errors in browser',
              'CORS_ORIGINS in .env doesn’t exactly match the live domains (scheme + host, no trailing slash) — redeploy api after fixing.',
            ],
            [
              'Cheaper alternative to the ALB',
              'If you don’t need AWS-native routing/WAF/autoscaling, you can skip Parts 7–10 entirely and instead run a self-managed reverse proxy (e.g. Caddy) directly on the EC2 instance with automatic Let’s Encrypt certs — see the Lightsail guide’s Part 6 for the exact pattern, which works identically on a plain EC2 box.',
            ],
          ],
        },
      ],
    },
  ],
}

export const AWS_LIGHTSAIL: DeploymentGuide = {
  slug: 'aws-lightsail',
  title: 'Deploying Sentex on AWS Lightsail',
  intro: [
    'Lightsail is AWS’s simplified VPS product — same underlying EC2/Docker mechanics as the plain-VM guides, but with a flat monthly price, a simplified firewall UI, a static IP included, and one-click snapshots. The one thing to know going in: **Lightsail’s own load balancer cannot do host-based routing** (it forwards everything to one instance port, no `Host:` header rules) — so instead of an AWS load balancer in front, this guide runs a small **Caddy** reverse-proxy container on the instance itself, which does host-based routing *and* gets you free, auto-renewing HTTPS with zero manual certificate steps. It’s the same job Traefik does in the Dokploy guide, just self-managed instead of platform-managed.',
  ],
  sections: [
    {
      heading: 'Architecture at a glance',
      blocks: [
        {
          kind: 'code',
          language: 'text',
          code: `Internet → Lightsail static IP → Caddy container (80/443, auto HTTPS)
                                    ├─ Host: api.sentex.gritcms.com   → api:8080
                                    ├─ Host: admin.sentex.gritcms.com → admin:3000
                                    └─ Host: sentex.gritcms.com       → web:3000
postgres + redis: no public/host ports, reachable only on the
instance's internal Docker network.`,
        },
      ],
    },
    {
      heading: 'Part 1 — Create the instance',
      blocks: [
        {
          kind: 'ol',
          items: [
            { text: 'Lightsail console → **Create instance**.' },
            {
              text: 'Platform: **Linux/Unix**. Blueprint: **OS Only → Ubuntu 24.04 LTS** (skip the app blueprints — you’re bringing your own Docker Compose stack).',
            },
            {
              text: 'Instance plan: pick **at least the 4 GB RAM / 2 vCPU plan**. The 2 GB plan can OOM while building two Next.js apps and a Go API back-to-back in the same `docker compose build` run.',
            },
            { text: 'Name it (e.g. `sentex-prod`), choose your region, and create it.' },
            { text: 'Wait for the instance state to show **Running**.' },
          ],
        },
      ],
    },
    {
      heading: 'Part 2 — Static IP',
      blocks: [
        {
          kind: 'ol',
          items: [
            {
              text: 'In the instance’s **Networking** tab, click **Create static IP**, attach it to `sentex-prod`.',
            },
            {
              text: 'Note the static IP — this is what your DNS records will point at. (Skip Part 2 entirely if you plan to put a Lightsail load balancer in front — see the note at the end of this guide — but for the recommended Caddy-on-instance setup below, you want the static IP.)',
            },
          ],
        },
      ],
    },
    {
      heading: 'Part 3 — Firewall (networking tab)',
      blocks: [
        {
          kind: 'p',
          text: 'Lightsail’s instance firewall is a simplified security group. Configure:',
        },
        {
          kind: 'table',
          headers: ['Application', 'Protocol', 'Port', 'Source'],
          rows: [
            ['SSH', 'TCP', '22', 'Restrict to your IP'],
            ['HTTP', 'TCP', '80', 'Anywhere'],
            ['HTTPS', 'TCP', '443', 'Anywhere'],
          ],
        },
        {
          kind: 'p',
          text: 'Do **not** open 8080/3000/3001 — the app containers stay internal, only reachable through Caddy on 80/443, same principle as locking the EC2 security group to ALB-only in the EC2 guide.',
        },
      ],
    },
    {
      heading: 'Part 4 — Install Docker',
      blocks: [
        { kind: 'p', text: 'Connect via the Lightsail browser SSH button, or:' },
        {
          kind: 'code',
          language: 'bash',
          code: 'ssh -i /path/to/your-lightsail-key.pem ubuntu@<static-ip>',
        },
        { kind: 'p', text: 'Install Docker Engine + Compose plugin:' },
        {
          kind: 'code',
          language: 'bash',
          code: `curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh
sudo usermod -aG docker $USER
newgrp docker
docker compose version`,
        },
        {
          kind: 'p',
          text: '(Recommended on the 4 GB plan) add 2 GB of swap the same way as the EC2 guide, so concurrent builds don’t get killed under memory pressure:',
        },
        {
          kind: 'code',
          language: 'bash',
          code: `sudo fallocate -l 2G /swapfile
sudo chmod 600 /swapfile
sudo mkswap /swapfile
sudo swapon /swapfile
echo '/swapfile none swap sw 0 0' | sudo tee -a /etc/fstab`,
        },
      ],
    },
    {
      heading: 'Part 5 — Get the code',
      blocks: [
        {
          kind: 'code',
          language: 'bash',
          code: 'git clone https://github.com/<you>/sentex.git\ncd sentex',
        },
      ],
    },
    {
      heading: 'Part 6 — Add the Caddy reverse proxy and adjust the Compose file',
      blocks: [
        {
          kind: 'p',
          text: 'Create `docker-compose.lightsail.yml` as a copy of `docker-compose.prod.yml` with these changes:',
        },
        {
          kind: 'ol',
          items: [
            {
              text: '**Remove** the `traefik.*` labels from `api`, `admin`, `web` — nothing is watching for them here.',
            },
            {
              text: '**Remove** the `dokploy-network` entry (top-level and from every service’s `networks:` list) — keep only the internal `sentex` network. Add `caddy` to that same network so it can reach the app containers by service name.',
            },
            {
              text: '**Add** a `caddy` service:',
              code: {
                language: 'yaml',
                code: `services:
  caddy:
    image: caddy:2-alpine
    container_name: sentex-caddy
    restart: unless-stopped
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./Caddyfile:/etc/caddy/Caddyfile:ro
      - caddy-data:/data
      - caddy-config:/config
    networks: [sentex]

volumes:
  caddy-data:
  caddy-config:`,
              },
            },
            {
              text: 'Create `Caddyfile` in the repo root:',
              code: {
                language: 'text',
                code: `api.sentex.gritcms.com {
    reverse_proxy api:8080
}

admin.sentex.gritcms.com {
    reverse_proxy admin:3000
}

sentex.gritcms.com {
    reverse_proxy web:3000
}`,
              },
            },
            {
              text: 'Caddy requests and renews a Let’s Encrypt certificate for each site block automatically the first time it starts and each domain resolves — no `certbot`, no cert-renewal cron job, no ACM console. This is the single biggest reason to reach for Caddy over plain nginx here.',
            },
          ],
        },
        {
          kind: 'p',
          text: 'Everything else in the Compose file — `build:`, `env_file:`, `environment:`, `volumes:`, `depends_on:`, `healthcheck:`, and the `migrate` service’s `command: ["sh", "-c", "./migrate && ./seed"]` with `restart: "no"` — stays **unchanged**. Real Docker Compose on a real VM honors `depends_on: condition: service_completed_successfully` natively, so `migrate` still blocks `api` from starting until it exits 0.',
        },
      ],
    },
    {
      heading: 'Part 7 — Environment variables',
      blocks: [
        { kind: 'code', language: 'bash', code: 'nano .env' },
        { kind: 'p', text: 'Paste your production values, same as the EC2 guide:' },
        {
          kind: 'code',
          language: 'bash',
          code: `WEB_DOMAIN=sentex.gritcms.com
ADMIN_DOMAIN=admin.sentex.gritcms.com
API_DOMAIN=api.sentex.gritcms.com
POSTGRES_USER=...
POSTGRES_PASSWORD=...
POSTGRES_DB=...
THEME=atlas
SOCIAL_AUTH_ENABLED=false
NEXT_PUBLIC_DEMO_LOGINS=false
R2_ACCOUNT_ID=...
R2_ACCESS_KEY_ID=...
R2_SECRET_ACCESS_KEY=...
R2_BUCKET=...
R2_ENDPOINT=...`,
        },
      ],
    },
    {
      heading: 'Part 8 — DNS (before you deploy)',
      blocks: [
        {
          kind: 'p',
          text: 'Caddy’s automatic HTTPS validates ownership over the internet (HTTP-01 challenge on port 80), so DNS has to resolve **before** the containers start, unlike the EC2/ALB flow where you could issue the cert independently via DNS validation first.',
        },
        {
          kind: 'p',
          text: 'Create three **A records** at your DNS provider, pointing at the Lightsail static IP from Part 2:',
        },
        {
          kind: 'code',
          language: 'text',
          code: `sentex.gritcms.com          → <static IP>
admin.sentex.gritcms.com    → <static IP>
api.sentex.gritcms.com      → <static IP>`,
        },
        {
          kind: 'p',
          text: 'Wait for propagation (`dig +short sentex.gritcms.com` should return the static IP) before moving to Part 9, or Caddy’s first certificate request will fail and retry on a backoff.',
        },
      ],
    },
    {
      heading: 'Part 9 — First deploy',
      blocks: [
        {
          kind: 'code',
          language: 'bash',
          code: 'docker compose -f docker-compose.lightsail.yml --env-file .env up -d --build',
        },
        { kind: 'p', text: 'Watch `migrate` complete, then watch Caddy issue certificates:' },
        {
          kind: 'code',
          language: 'bash',
          code: `docker compose -f docker-compose.lightsail.yml logs -f migrate
docker compose -f docker-compose.lightsail.yml logs -f caddy`,
        },
        {
          kind: 'p',
          text: 'You’re looking for lines like `certificate obtained successfully` for each of the three domains in the Caddy log.',
        },
      ],
    },
    {
      heading: 'Part 10 — Verify',
      blocks: [
        {
          kind: 'ol',
          items: [
            { text: 'Visit all three domains over `https://` — valid Let’s Encrypt cert, no warnings.' },
            { text: 'Confirm no CORS errors in the browser console on `admin`/`web`.' },
            { text: 'Log in with the seeded demo credentials to confirm `migrate`/`seed` ran.' },
            {
              text: '`docker compose -f docker-compose.lightsail.yml ps` — all containers should show healthy/running, `migrate` should show `Exited (0)`.',
            },
          ],
        },
      ],
    },
    {
      heading: 'Part 11 — Backups',
      blocks: [
        {
          kind: 'p',
          text: 'Lightsail’s headline convenience feature: turn on **Automatic snapshots** under the instance’s **Snapshots** tab. This snapshots the entire instance disk daily (including your Postgres data directory’s Docker volume, since it lives on the instance’s own disk) — a much lower-effort baseline than configuring `pg_dump` cron jobs yourself, though for a real production system you’ll still want logical Postgres backups (`pg_dump`) in addition, since a full-disk snapshot restores the *whole instance* to a point in time, not just the database.',
        },
      ],
    },
    {
      heading: 'Part 12 — Ongoing deploys',
      blocks: [
        { kind: 'p', text: 'Same two options as the EC2 guide — manual:' },
        {
          kind: 'code',
          language: 'bash',
          code: `ssh ubuntu@<static-ip>
cd sentex && git pull
docker compose -f docker-compose.lightsail.yml --env-file .env up -d --build`,
        },
        { kind: 'p', text: 'or GitHub Actions:' },
        {
          kind: 'code',
          language: 'yaml',
          code: `name: Deploy to Lightsail
on:
  push:
    branches: [main]
jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: appleboy/ssh-action@v1
        with:
          host: \${{ secrets.LIGHTSAIL_STATIC_IP }}
          username: ubuntu
          key: \${{ secrets.LIGHTSAIL_SSH_KEY }}
          script: |
            cd sentex
            git pull
            docker compose -f docker-compose.lightsail.yml --env-file .env up -d --build`,
        },
      ],
    },
    {
      heading: 'Alternatives worth knowing about',
      blocks: [
        {
          kind: 'p',
          text: '**Lightsail Load Balancer instead of Caddy.** Lightsail does offer a managed load balancer with free auto-renewing certificates (up to 9 alternate domains on one cert) — but it forwards *all* traffic to a single instance port with no host-header routing. To use it here you’d still need a reverse proxy on the instance doing the host-based split; the load balancer would just replace Caddy’s TLS termination (forwarding plain HTTP to instance port 80) while an nginx/Caddy container on the instance did the routing over HTTP internally. For a single-instance deployment like this one, that’s extra cost and moving parts for the same result the Caddy-only setup already gives you for free — worth it mainly if you later add a *second* instance and want AWS to load-balance between them.',
        },
        {
          kind: 'p',
          text: '**Lightsail Container Service.** A separate, fully-managed alternative to running your own instance — closer in spirit to Railway/Render than to this guide. It runs container images (not a `docker-compose.yml`) behind a managed HTTPS endpoint, has its own concept of "public endpoint" per deployment, and pairs with Lightsail’s own managed database offering instead of a `postgres` container. If you’d rather not manage an instance at all, that’s a different (and shorter) path than this guide — but it means giving up the one-file-does-everything simplicity of deploying `docker-compose.prod.yml` more or less as written.',
        },
        {
          kind: 'p',
          text: '**An AWS ALB in front of a Lightsail instance.** Also possible (via VPC peering between the Lightsail and EC2 networking planes), which gets you the exact host-based-routing pattern from the EC2 guide while keeping your compute on cheaper Lightsail pricing. More moving parts than either guide alone — only worth it if you specifically want ALB features (WAF, weighted routing, etc.) without leaving Lightsail entirely.',
        },
      ],
    },
    {
      heading: 'Notes and troubleshooting',
      blocks: [
        {
          kind: 'table',
          headers: ['Symptom', 'Likely cause'],
          rows: [
            [
              'Caddy never issues a certificate',
              'DNS wasn’t pointing at the static IP yet when Caddy first started — fix DNS, then docker compose restart caddy.',
            ],
            [
              'Exited (1) on migrate',
              'Check docker compose logs migrate — usually a bad POSTGRES_* value in .env, or Postgres not yet healthy (shouldn’t happen given the healthcheck + depends_on, but worth checking docker compose logs postgres too).',
            ],
            [
              'Build gets killed partway through',
              'Instance too small — bump to the 4 GB+ plan or add swap (Part 4).',
            ],
            [
              '502 from Caddy',
              'The upstream container isn’t listening yet or crashed — docker compose ps and docker compose logs api (or admin/web).',
            ],
            [
              'CORS errors in browser',
              'CORS_ORIGINS in .env doesn’t exactly match the live domains — fix and redeploy api.',
            ],
          ],
        },
      ],
    },
  ],
}
