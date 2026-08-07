# Deploying Sentex on Orbita

A self-contained guide for deploying Sentex to **Orbita** — the self-hosted,
multi-tenant PaaS built by the same team as Grit. It targets the same monorepo
layout (`apps/api`, `apps/admin`, `apps/web`, Postgres, Redis, `migrate` job)
and the same three hosts used in the Railway, Dokploy, Coolify, Render and
Fly.io guides.

Orbita is different from the other four platforms in one important way: because
Sentex is a **Grit app**, Orbita recognises it and deploys the whole thing with
almost no configuration. It reads `grit.json`, works out that this is a
three-service app, builds the Dockerfiles Grit already ships, provisions
Postgres and Redis for you, runs your migrations under a lock, and wires up all
three domains — from a manifest that's about ten lines long. You don't translate
your Compose file the way Render and Fly.io needed, and you don't hand-edit it
the way Coolify needed.

You install Orbita **once** (Part 1), then deploy Sentex to it. There are two
ways to deploy, and they produce the same result:

- **Part 2 — the Grit fast path** (recommended). Orbita derives everything from
  `grit.json`. This is the path that makes Orbita worth using for a Grit app.
- **Part 3 — the Docker Compose path** (parity with Dokploy/Coolify). Keep
  `docker-compose.prod.yml` as your source of truth and let Orbita run it as a
  Swarm stack. Use this only if you specifically want the Compose file to stay
  authoritative.

Jump to the part you need. Part 1 is a prerequisite for both.

---

# Part 1 — Stand up Orbita (once)

This installs Orbita on a VPS: it becomes your dashboard, your build server, and
your reverse proxy. You do this once, then deploy as many apps (and as many
tenants) as you like onto it. If you already have an Orbita server running, skip
to [Part 2](#part-2--deploy-sentex-the-grit-way-recommended).

The install brings **Docker, Docker Swarm, PostgreSQL, Redis, and Traefik** with
it — there is nothing to set up beforehand.

## 1.1 Provision the VPS

1. Spin up a VPS. **2 vCPU / 4 GB RAM** is a sensible minimum — Orbita itself is
   tiny (~50 MB idle), but *builds* are hungry, and you'll be building three
   Next.js/Go services plus running Postgres and Redis. 4 GB gives the builds
   room; go higher if you expect real traffic or multiple tenants.
2. Use a **fresh Ubuntu 22.04 or 24.04** box. Orbita's Traefik needs to own
   ports **80** and **443**, so don't install it on a server that already runs
   nginx, Apache, or another panel — the installer will stop if it finds one.
3. Note the server's **public IP** and the **root password** (or the SSH key)
   the provider gave you. You'll need them in 1.3.

Any of Contabo, Hetzner, DigitalOcean, or Vultr work well.

## 1.2 Point DNS at the server

Orbita's Traefik terminates TLS directly on your server's IP, so these are
**A records** (not CNAMEs) — exactly like Dokploy in Part 1 of the alternatives
guide. You need **four** hosts: one for the Orbita dashboard itself, and the
three Sentex app hosts.

```
orbita.gritcms.com          → <VPS public IP>     # the Orbita dashboard
sentex.gritcms.com          → <VPS public IP>     # web
admin.sentex.gritcms.com    → <VPS public IP>     # admin
api.sentex.gritcms.com      → <VPS public IP>     # api
```

If your DNS is behind Cloudflare, set each record to **DNS only** (grey cloud),
not Proxied — Orbita fetches the Let's Encrypt certificate itself, and the
orange-cloud proxy gets in the way of that first handshake. You can turn it back
on afterwards.

Confirm the dashboard record resolves before installing — certificates can't be
issued until it does:

```
dig orbita.gritcms.com +short      # must print your VPS IP
```

## 1.3 Harden the server

Never skip this. A fresh VPS with root SSH open is a target within minutes. SSH
in as root, then run the hardening script:

```
ssh root@<VPS public IP>

curl -sSL https://raw.githubusercontent.com/MUKE-coder/vps-harden/main/vps-harden.sh -o vps-harden.sh
chmod +x vps-harden.sh
sudo ./vps-harden.sh --no-dokploy
```

It asks a few plain questions:

- **A username** for your everyday account — type `deploy`.
- **An SSH port** — press Enter to keep the default.
- **Your SSH public key** — paste your `~/.ssh/id_ed25519.pub`, or leave it
  blank and the script generates a key and tells you where it saved it.
- **A password for the account** — leave it blank and the script generates a
  strong one and **prints it once at the end**. Save that password.

When it finishes you have a `deploy` user with its own password *and*
passwordless `sudo`, your key installed, root and password SSH logins disabled,
a firewall (UFW + ufw-docker), Fail2ban, kernel hardening, and a 0–100 security
score. `--no-dokploy` tells it not to install Dokploy, since Orbita is your
platform here.

Before you disconnect: open a **second terminal** and confirm
`ssh deploy@<VPS public IP>` works, so you don't lock yourself out.

## 1.4 Install Orbita

From here you're logged in as `deploy` (the `deploy` user reaches Docker through
`sudo` — it's deliberately not in the `docker` group, which is root-equivalent).
Run the one-line installer, passing your dashboard domain and an email for
Let's Encrypt:

```
curl -sSL https://raw.githubusercontent.com/MUKE-coder/orbita/main/install.sh \
  | sudo ORBITA_DOMAIN=orbita.gritcms.com ORBITA_ACME_EMAIL=you@gritcms.com bash -s -- --yes
```

(If you skipped the domain in 1.2 and want to trial on the IP, drop the two env
vars and Orbita comes up on `http://<VPS public IP>:8080` with no TLS.)

In order, the installer installs Docker and starts Swarm, checks ports
80/443/8080, generates secrets into `/opt/orbita/.env`, pulls the Orbita image,
starts all four services (`orbita`, `orbita-postgres`, `orbita-redis`,
`orbita-traefik`), opens the firewall for the ports it needs, and waits for a
healthy `/health` before printing your dashboard URL.

Verify:

```
curl -s http://localhost:8080/health          # want {"status":"ok", ...}
cd /opt/orbita && sudo docker compose ps       # all four services "Up"
```

**Back up `/opt/orbita/.env`.** It holds `ENCRYPTION_MASTER_KEY`, from which
every organisation's encryption key is derived. Lose it and every stored secret
(including the Sentex env you're about to upload) is unrecoverable. Copy it
somewhere safe before you put real data in.

## 1.5 Create your super-admin and organisation

Open the dashboard:

- **With a domain:** `https://orbita.gritcms.com`
- **IP only:** `http://<VPS public IP>:8080`

Click **Register** and create your account **immediately** — the first person to
register becomes the **super-admin** with full control of the box. Once that
account exists, public sign-up **closes automatically**; nobody else can walk in
and register. (Later teammates join by invitation, or through an account you
create for them under **Admin**.)

Then create an **organisation** — your top-level workspace. Everything lives
inside one, and each org is fully isolated: its own Docker network, its own
encryption key, its own resource quota. Name it something like `gritcms` (this
becomes the org **slug**, used to namespace networks and volumes).

If you're running Sentex for a client and want to hand *them* the org, use
**Admin → Organisations → New tenant** instead: it creates the org, sizes it
(CPU/RAM/disk/app limits), and creates the client's login in one step, showing
you a generated password to hand over. They set their own password at first
sign-in. For deploying your own app, a plain organisation is fine.

## 1.6 Connect GitHub

Sentex is a private repo, so Orbita needs a token to clone it and to register the
auto-deploy webhook.

1. In the dashboard, go to **Settings → Git Connections**.
2. Add a **GitHub** connection with a Personal Access Token that has the **`repo`**
   and **`admin:repo_hook`** scopes. `repo` lets Orbita clone the private
   `sentex` repository; `admin:repo_hook` lets it install the push-to-deploy
   webhook so future commits redeploy automatically.

That's the whole platform set up. Everything below is per-app.

---

# Part 2 — Deploy Sentex the Grit way (recommended)

This is the path that makes Orbita worth using for a Grit app. You write a short
`orbita.yaml`, and Orbita derives the rest from `grit.json`.

## 2.1 What Orbita derives from `grit.json` (you write none of this)

A Grit app has a known shape, declared in `grit.json` at the repo root. Sentex
has `apps/api`, `apps/admin`, and `apps/web`, which is Grit's **triple**
architecture. From that single fact, Orbita works out the entire deployment —
this is the table you would otherwise have hand-written as a Compose file,
`render.yaml`, or three `fly.toml` files in the other guides:

| From `grit.json` | Orbita derives for Sentex |
|---|---|
| `architecture: triple` | Three containers: `api`, `web`, `admin` (plus `docs` if the repo has it) |
| The Dockerfiles Grit ships | Builds `api` from `apps/api`, and the Next.js apps from the repo root — the exact `build.context` values you set by hand in `docker-compose.prod.yml` |
| Ports | `8080` for the API, `3000` for the Next.js apps — no `:port` mapping to write |
| `NEXT_PUBLIC_API_URL` | Baked into the `admin` and `web` bundles at build time from your **api** domain — the single most error-prone value in every other guide, derived here |

Orbita does **not** generate a Dockerfile and does **not** fall back to Nixpacks
for a Grit app — it reuses the correct multi-stage Dockerfiles Grit already
ships, the same ones the other platforms build.

## 2.2 Write `orbita.yaml`

Create `orbita.yaml` at the Sentex repo root. This is the whole deploy config —
compare it to the ~90-line `docker-compose.prod.yml`, the `render.yaml`
Blueprint, or the three `fly.toml` files:

```yaml
app: sentex
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
  from: .env.production          # local file; values are encrypted into Orbita, never committed
```

Notes on the choices, mapped to `docker-compose.prod.yml`:

- **`addons: [postgres, redis]`** replaces the `postgres` and `redis` services
  you declared in Compose. Orbita provisions managed instances inside your org's
  private Docker network and injects their connection URLs (e.g. `DATABASE_URL`,
  `REDIS_URL`) into the app — the values Grit's code reads. You don't set
  `POSTGRES_HOST` / `POSTGRES_PASSWORD` / `REDIS_URL` yourself; that's the point.
  (If Sentex uses Orbita's object storage instead of Cloudflare R2, add `minio`
  to the list and Orbita injects `MINIO_*` / `STORAGE_DRIVER` too. If you're
  keeping R2, leave `minio` off and put the `R2_*` values in `.env.production` —
  see 2.3.)
- **`domains`** replaces the three Traefik `Host()` labels in your Compose file.
  Bare hostnames only — no scheme, port, or path. Orbita creates the routers and
  fetches Let's Encrypt certs for all three.
- **`migrate: true`** replaces the standalone `migrate` service and its
  `depends_on: condition: service_completed_successfully`. Orbita runs
  `cmd/migrate` in a one-off container **before cutover**, under a Postgres
  advisory lock. This is stronger than the Compose ordering guarantee: a
  non-zero exit **aborts the deploy** and the previous version keeps serving (see
  2.8).

Two optional Grit toggles, both **on by default**, so you only add them to turn
something off:

```yaml
observability: true    # Pulse    — latency/SQL/error tracing on the API
security: true         # Sentinel — WAF, rate limiting, anomaly detection on the API
studio: false          # GORM Studio — off by default; it edits live data
```

## 2.3 Provide the environment values

Everything that isn't derived comes from `.env.production` — the same file you'd
paste into Dokploy or fill into Coolify, minus the Postgres/Redis values Orbita
now supplies. That means the build-time and app-secret values Sentex needs:

```
THEME=atlas
SOCIAL_AUTH_ENABLED=false
NEXT_PUBLIC_DEMO_LOGINS=false

# Cloudflare R2 (only if you're keeping R2 instead of Orbita's minio addon)
R2_ACCOUNT_ID=...
R2_ACCESS_KEY_ID=...
R2_SECRET_ACCESS_KEY=...
R2_BUCKET=...
R2_ENDPOINT=...

# ...any other app secrets Sentex reads at runtime
```

You do **not** put `NEXT_PUBLIC_API_URL`, `DATABASE_URL`, `REDIS_URL`,
`POSTGRES_*`, or the domain variables here — Orbita derives those from your
`domains` and `addons`. When Orbita reads `env.from`, it encrypts every value
into the org's key at rest; the file is never committed and never leaves your
machine in plaintext.

Now pick a route: **2.4** (dashboard) or **2.5** (CLI). They do the same thing.

## 2.4 Route A — deploy from the dashboard

No CLI needed; everything happens in the browser.

1. In your organisation, create a **Project** (e.g. *Sentex*) and an
   **Environment** (e.g. *production*) inside it. Apps live under
   project → environment.
2. Click **Create App → Source: Git Repository**. Pick the GitHub connection
   from 1.6, then the `sentex` repo and the `main` branch.
3. Because the repo has a `grit.json`, Orbita recognises it as a Grit app and
   uses the fast path — you don't choose a builder or a Dockerfile. Confirm the
   three derived domains match your DNS from 1.2:
   `sentex.gritcms.com`, `admin.sentex.gritcms.com`, `api.sentex.gritcms.com`.
4. Open the app's **Environment** tab and paste the contents of
   `.env.production` (one `KEY=VALUE` per line). Mark the `R2_*` keys and any
   secrets as **secret** so they're encrypted at rest and never shown again.
5. Click **Deploy**.

Watch the **Deployments** tab. Orbita builds all three services, provisions
Postgres and Redis, runs the migration, and cuts over only if it succeeds. Skip
to [2.6](#26-what-happens-in-order) for what you're watching.

## 2.5 Route B — deploy with the CLI

The `orbita` CLI is optional, but it's the tidiest way to deploy from your
machine and keep `orbita.yaml` as the source of truth. Today it's **built from
source** — there's no `curl | sh` installer yet, and `go install` doesn't work
because of the module path — so build it once from a clone (requires Go 1.25+):

```
git clone https://github.com/MUKE-coder/orbita.git
cd orbita
make build-cli
sudo mv ./orbita /usr/local/bin/orbita
orbita --help
```

It won't clash with Grit's own `grit` binary — different repo, different name.

Then, from the **Sentex** project directory (the one with `grit.json` and the
`orbita.yaml` you wrote in 2.2):

```
# Register your server with the CLI (once). Prompts for the admin email +
# password you created in 1.5, mints a deploy token, saves the host as "prod".
orbita login https://orbita.gritcms.com

# Store a GitHub token (repo + admin:repo_hook) so Orbita can push/clone (once).
orbita github-auth

# Preview the plan without changing anything — highly recommended first run.
orbita deploy --plan --host prod
```

The plan prints exactly what it will create, so you can confirm the mode and
domains before anything happens:

```
▸ Plan (dry run — nothing will be changed)
  App:       sentex
  Mode:      triple
  Migrate:   true
  Addons:    postgres, redis

  create  sentex-api    → api.sentex.gritcms.com
  create  sentex-web    → sentex.gritcms.com
  create  sentex-admin  → admin.sentex.gritcms.com
```

When it looks right, deploy for real:

```
orbita deploy --host prod
```

(If you don't yet have an Orbita server at all, `orbita init` collapses all of
Part 1 — harden, install, admin account, host registration — into one
interactive command from your machine. Use it *instead of* Part 1, not as well.)

## 2.6 What happens, in order

Whichever route you used, a deploy runs these steps — this is the pipeline the
other four platforms make you assemble by hand:

1. **Detect** — `grit.json` at the repo root marks it a Grit app; `architecture:
   triple` picks the three-service strategy.
2. **Ensure the repo** — Orbita confirms it can reach the `sentex` repo with your
   token (and, over the CLI, pushes your current commit).
3. **Reconcile** — org, project, environment, the `postgres` + `redis` addons,
   your encrypted env, and the three domains. Idempotent — safe to re-run.
4. **Build** — `api`, `web`, and `admin` from the Dockerfiles Grit ships, with
   `NEXT_PUBLIC_API_URL` baked into the two Next.js bundles from your api domain.
5. **Migrate** — `cmd/migrate` in a one-off container, under a Postgres advisory
   lock so two concurrent deploys can't race.
6. **Cut over** — only if the migration exited 0. The previous images are kept
   for instant rollback.
7. **Route** — Traefik serves all three domains over HTTPS. Certs are issued on
   the first request to each host.

## 2.7 Verify

1. Visit `https://api.sentex.gritcms.com/<health endpoint>` and confirm it
   answers.
2. Visit `https://sentex.gritcms.com` and `https://admin.sentex.gritcms.com` and
   confirm there are **no CORS errors** in the browser console — if there are,
   the API domain baked into the frontend bundle doesn't match `api`'s real
   domain; recheck `domains.api` in `orbita.yaml`.
3. Log in with the seeded demo SACCO credentials to confirm the migration and
   seed actually ran.

From the CLI you can also stream logs and confirm the migration:

```
orbita logs -f --host prod                       # all services
orbita logs --host prod --service migrate        # just the migration job
```

## 2.8 Migrations gate the cutover (troubleshooting)

Orbita runs your migrations **before** it cuts over, under an advisory lock. A
non-zero exit stops the deploy and leaves the previous version serving — you
never end up on a schema-mismatched image. If a deploy fails at the migrate
step, that's why.

The most common cause with a Grit app is **`go.sum` not being committed**, so
`go run ./cmd/migrate` can't resolve modules inside the one-off container. Commit
it — real Grit apps ship it — and redeploy. Check the migrate log:

```
orbita logs --host prod --service migrate
```

## 2.9 Batteries included (Pulse, Sentinel, Studio)

Because Sentex is a Grit app, Orbita mounts these on the API by default — no
setup:

- **Pulse** — latency, SQL, and error tracing → `https://api.sentex.gritcms.com/pulse/ui`
- **Sentinel** — WAF, rate limiting, anomaly detection → `https://api.sentex.gritcms.com/sentinel/ui`
- **GORM Studio** — off by default because it edits live data. Turn it on with
  `studio: true` in `orbita.yaml` only when you need it.

## 2.10 Ongoing deploys

Because you connected GitHub in 1.6, Orbita installed a push-to-deploy webhook
when it created the app. Every push to `main` now:

- re-clones the repo and rebuilds the changed services,
- reruns `cmd/migrate` under the lock (idempotent, so this is safe every time),
- cuts over only if the migration succeeds.

No manual redeploy step — the same GitHub-connected flow as Dokploy and Coolify.
To revert a bad deploy, `orbita rollback --host prod` (or the **Rollback** button
on a previous deployment in the dashboard) swaps back to the previous image
instantly, since Orbita keeps it.

---

# Part 3 — Alternative: run your Compose file on Orbita

Use this **only** if you specifically want `docker-compose.prod.yml` to stay the
source of truth — for example, to keep one Compose file working identically
across Dokploy, Coolify, and Orbita. For a Grit app, Part 2 is simpler and gives
you migrations-under-a-lock, provisioned addons, and the observability mounts
that this path does not. This path treats Sentex as a generic multi-service
stack, not as a Grit app.

Orbita deploys a Compose file as a **Docker Swarm stack**. It runs the file
essentially as-is — you do **not** strip networks or labels the way Coolify
required, and you do **not** translate it to another format the way Render and
Fly.io required.

## 3.1 Create the app from Docker Compose

1. In your project/environment, click **Create App → Source: Docker Compose**.
2. Point it at the compose file:
   - **From your Git repo** (recommended, so pushes redeploy): pick the GitHub
     connection, the `sentex` repo and branch, and set the compose file path to
     `docker-compose.prod.yml`.
   - **Or paste it inline** — but note a pasted file can't use `build:` (there's
     no source tree to build from), so it must reference prebuilt images. Sentex
     builds from source, so use the Git option.
3. Set the **web service** to the service that serves your primary domain —
   `web` for Sentex. This is the service Orbita routes your app domain to; the
   others stay private to the stack, reachable by their compose service name
   (`api`, `postgres`, `redis`) exactly as they are locally.
4. Set the **port** to the web service's container port — `3000` for the Sentex
   `web` service. Port is required for Compose apps, because that's what the
   domain routes to.

## 3.2 Domains

Add your domains under the app's **Domains** tab. Only the nominated **web
service** is routable from a single Compose app, so:

- Add `sentex.gritcms.com` → routes to the `web` service you nominated.
- To give `api` and `admin` their own domains, the clean approach on this path is
  to **deploy each as its own app** (three Compose apps, or better, use the Grit
  fast path in Part 2 which does all three at once). A single Compose app exposes
  one routable service.

This is the main reason Part 2 is preferable for a three-domain app like Sentex —
the Grit fast path routes all three hosts from one deployment.

## 3.3 Environment

Open the app's **Environment** tab and paste your `.env.production` values. Orbita
injects them into **every** service in the stack (so a worker gets the same
`DATABASE_URL` the web tier does), and they're encrypted at rest. A service's own
`environment:` block in the compose file still wins if it sets the same key.
`${VAR}` references in the compose file are interpolated from these values too,
matching `env_file: [.env]`.

## 3.4 Deploy and verify

1. Click **Deploy**. Orbita builds the services that declare `build:`, then runs
   `docker stack deploy` for the whole file.
2. Watch the deploy log. Because this is a real Swarm deploy of your Compose
   file, the `migrate` service's `depends_on: service_completed_successfully`
   ordering works unmodified — same as Dokploy/Coolify.
3. Visit your web domain and confirm it serves.

## 3.5 Limits worth knowing on the Compose path

- **Only the web service is routable** per Compose app (see 3.2).
- **No rollback** for Compose apps — a Compose deploy has no single image to
  revert to. Redeploy the previous commit instead. (The Grit path in Part 2
  *does* support instant rollback.)
- **`build:` needs a Git repo** — pasted YAML must use prebuilt images.
- Stopping, starting, or deleting the app applies to **every** service in the
  stack.

---

# Quick comparison

How Orbita's two paths sit next to the four platforms in the alternatives guide:

| | Runs your Compose file? | Migration handling | Domain/TLS | GitHub auto-deploy |
|---|---|---|---|---|
| **Orbita — Grit fast path** | No — derives everything from `grit.json`; you write ~10 lines of `orbita.yaml` | `cmd/migrate` under a **Postgres advisory lock**, gates cutover | All three hosts from one deploy; Traefik + Let's Encrypt, derived | Built-in webhook (push to redeploy) |
| **Orbita — Compose path** | Yes, almost unmodified (no stripping networks/labels) | Native `depends_on:` (unchanged) | One routable web service per app; Traefik + Let's Encrypt | Built-in webhook |
| **Dokploy** | Yes, almost unmodified | Native `depends_on: service_completed_successfully` | Traefik labels already in the file | Built-in webhook toggle |
| **Coolify** | Yes, after removing custom `networks:`/labels | Native `depends_on:` (unchanged) | UI **Domains** field per service | Built-in webhook toggle |
| **Render** | No — translated to `render.yaml` | `preDeployCommand` on `api` | `domains:` in Blueprint + CNAME | Built-in (Blueprint sync + per-service auto-deploy) |
| **Fly.io** | No — one `fly.toml` per app | `release_command` on `api` | `fly certs add` + DNS record | Manual — GitHub Actions workflow |

Why Orbita's Grit path is the shortest of them all for Sentex: it's the only one
that already knows what a Grit app is. The others need you to describe a
three-service app in their own dialect (Compose, Blueprint, or three TOMLs);
Orbita reads the same `grit.json` your app already ships and derives the rest —
addons, ports, build contexts, the API URL baked into the frontends, and
migrations under a lock — from that.
