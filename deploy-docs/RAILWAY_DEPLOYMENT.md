# Deploying Sentex to Railway (GitHub Deploy)

This guide translates `docker-compose.prod.yml` into a Railway project, deployed
from GitHub. Railway does **not** run `docker-compose.yml` files directly —
every service in the Compose file becomes its own Railway service, wired
together with Railway's private network and managed databases instead of
Traefik + Docker networks. Follow every step in order; later steps assume
earlier ones are done.

## 0. Service map (Compose → Railway)

| Compose service | Railway equivalent                              | Source                          |
| ---------------- | ------------------------------------------------ | -------------------------------- |
| `postgres`        | **Managed Postgres** plugin                       | n/a — provisioned by Railway     |
| `redis`            | **Managed Redis** plugin                          | n/a — provisioned by Railway     |
| `migrate`          | **Pre-Deploy Command** on the `api` service       | `apps/api`                       |
| `api`              | GitHub-deployed service, root dir `apps/api`      | `apps/api/Dockerfile`            |
| `admin`            | GitHub-deployed service, root dir repo root       | `apps/admin/Dockerfile`          |
| `web`              | GitHub-deployed service, root dir repo root       | `apps/web/Dockerfile`            |
| `traefik` labels   | Railway's built-in edge (automatic)               | n/a                               |
| `dokploy-network`  | Railway's automatic private network               | n/a                               |

Reasons for the two differences you'll notice below:

- **No standalone `migrate` service.** Railway has no `depends_on:
  condition: service_completed_successfully`, so a one-shot container that
  "runs once then exits" isn't a first-class concept the way it is in
  Compose. The correct Railway feature for this is a **Pre-Deploy Command**
  on the `api` service — it runs in its own container, before the new `api`
  version starts serving traffic, and the deploy is aborted if it fails.
- **Postgres/Redis become managed plugins, not containers.** Railway's
  guidance is to use its managed database services instead of raw
  `postgres:16-alpine` / `redis:7-alpine` images — you get backups,
  connection pooling and a dashboard for free, and you no longer need to
  manage volumes for them yourself.

---

## Part 1 — Prepare the repository

1. Confirm you're on the branch you intend to deploy (commonly `main`), and
   that it's pushed to GitHub. Railway builds from GitHub, so uncommitted
   local changes will not be deployed.
2. Confirm the three Dockerfiles exist at:
   - `apps/api/Dockerfile`
   - `apps/admin/Dockerfile`
   - `apps/web/Dockerfile`
3. Open `apps/admin/Dockerfile` and `apps/web/Dockerfile` and confirm they
   declare `ARG` for every build-time variable the Compose file passes under
   `build.args` (`NEXT_PUBLIC_API_URL`, `NEXT_PUBLIC_WEB_URL`,
   `NEXT_PUBLIC_ADMIN_URL`, `THEME`, `SOCIAL_AUTH_ENABLED`,
   `NEXT_PUBLIC_DEMO_LOGINS`). Railway only forwards a variable into the
   Docker build if the Dockerfile has a matching `ARG` line — without it,
   the value you set in the Railway dashboard is silently ignored and the
   bundle bakes in whatever the Dockerfile's default is.
4. Open `apps/api/Dockerfile` and confirm the `./migrate` and `./seed`
   binaries referenced in the Compose `command:` are actually produced by
   the build (i.e. they're part of the same image `api` runs from — not a
   separate build stage that gets discarded).
5. Check `.dockerignore` at the repo root. Since `admin` and `web` build
   with `context: .` (the whole repo), make sure `.dockerignore` excludes
   `node_modules`, `.next`, `.git`, and `e2e` so the build context upload to
   Railway stays small and fast.
6. You will **not** commit `.env` or `.env.production` — their values move
   into Railway's Variables UI in Part 5–7. Keep them open locally for
   reference while you copy values over.

---

## Part 2 — DNS records

Before deploying, create three DNS **A** records (or a **CNAME** once you
know Railway's target — see note below) pointing at wherever your traffic
will land:

| Host                          | Purpose                     |
| ------------------------------ | ---------------------------- |
| `sentex.gritcms.com`           | `web` — public shell + login |
| `admin.sentex.gritcms.com`     | `admin` — staff admin panel  |
| `api.sentex.gritcms.com`       | `api` — Go API               |

**Note on Railway specifically:** unlike the Dokploy/Traefik setup in your
Compose file, Railway custom domains are attached as a **CNAME** to a
Railway-generated target (e.g. `xxxx.up.railway.app`), not an A record to a
static IP. You'll get the exact CNAME target in Part 8, after each service
has a Railway-generated domain to point at. It's fine to leave DNS for last
— just don't skip it, since `WEB_DOMAIN` / `ADMIN_DOMAIN` / `API_DOMAIN`
are baked into the frontend bundles at build time and mismatches show up as
CORS/CSP failures, exactly as the comments in your Compose file warn.

Storage stays on Cloudflare R2 and needs no DNS record here, same as in the
Compose file — only the `R2_*` variables need to be set (Part 5).

---

## Part 3 — Create the Railway project

1. Go to [railway.com/dashboard](https://railway.com/dashboard) and log in
   (or sign up) with GitHub — this makes the GitHub connection in the next
   part one click instead of a separate OAuth flow.
2. Click **+ New Project**.
3. Choose **Empty Project**. (Don't use "Deploy from GitHub repo" on this
   screen — because this repo is a monorepo with multiple deployable
   services and a non-JS `apps/api` alongside JS apps, Railway's automatic
   monorepo importer won't produce the exact 4-service layout you need. It's
   more reliable to add each service by hand in Part 4–7.)
4. Rename the project (top left, click the project name) to something like
   `sentex-production` so it's unambiguous in your Railway dashboard.
5. If you have more than one Railway environment planned (e.g.
   `production` vs `staging`), do this whole guide once per environment —
   Railway environments are separate blank canvases under the same project.
   This guide assumes a single `production` environment.

---

## Part 4 — Add the Postgres database

1. Inside your empty project canvas, click **+ New**.
2. Select **Database** → **Add PostgreSQL**.
3. Railway provisions a Postgres instance immediately and creates a service
   card named `Postgres` on the canvas. It automatically generates and
   exposes these variables on that service (you don't set these — Railway
   does): `DATABASE_URL`, `DATABASE_PUBLIC_URL`, `PGHOST`, `PGPORT`,
   `PGUSER`, `PGPASSWORD`, `PGDATABASE`.
4. Click into the `Postgres` service → **Variables** tab and leave it as is.
   You'll pull these values into `api` via reference variables in Part 5 —
   don't retype the password anywhere.
5. No volume setup is needed here — managed databases handle their own
   storage, unlike the `postgres-data` named volume in your Compose file.

---

## Part 5 — Add the Redis database

1. Click **+ New** again.
2. Select **Database** → **Add Redis**.
3. Railway creates a `Redis` service card and exposes `REDIS_URL` and
   `REDIS_PUBLIC_URL` automatically, the same way it did for Postgres.
4. Same as Postgres — nothing to configure here; you'll reference
   `REDIS_URL` from the `api` service next.

---

## Part 6 — Deploy the `api` service (and the migration)

### 6.1 Create the service

1. Click **+ New** → **GitHub Repo**.
2. If this is the first time connecting this GitHub account/org, authorize
   the Railway GitHub App and grant it access to the `sentex` repository
   (or "All repositories" if you're comfortable with that).
3. Select the `sentex` repository, then select the branch you're deploying
   (e.g. `main`).
4. Railway creates a service and immediately tries to build it. It will
   likely fail or pick the wrong Dockerfile at this point — that's expected,
   since it just tried to build from the repo root. Continue to 6.2 before
   worrying about that first failed build.

### 6.2 Point it at `apps/api`

1. Click into the new service → **Settings** tab.
2. Under **Source**, find **Root Directory** and set it to:
   ```
   apps/api
   ```
   This matches `build.context: ./apps/api` in the Compose file. With the
   root directory set, Railway will look for `Dockerfile` inside
   `apps/api/` automatically — you don't need to set a separate Dockerfile
   path for this service.
3. Rename the service (click the service name at the top) to `api`, so its
   private-network hostname becomes `api.railway.internal`.
4. Under **Build**, scroll to **Watch Paths** and add:
   ```
   apps/api/**
   packages/**
   ```
   This stops a commit that only touches `apps/web` or `apps/admin` from
   triggering an unnecessary `api` rebuild.

### 6.3 Set environment variables

Go to the **Variables** tab and add the following. Use **Raw Editor** to
paste several at once (`KEY=VALUE` per line). For anything that references
another Railway service, use Railway's reference-variable syntax
`${{ServiceName.VARIABLE}}` instead of hardcoding a value — it stays in
sync automatically if credentials rotate.

```
APP_ENV=production

# Pulled from the managed Postgres plugin (Part 4) instead of a container
POSTGRES_HOST=${{Postgres.PGHOST}}
POSTGRES_PORT=${{Postgres.PGPORT}}
POSTGRES_USER=${{Postgres.PGUSER}}
POSTGRES_PASSWORD=${{Postgres.PGPASSWORD}}
POSTGRES_DB=${{Postgres.PGDATABASE}}

# Pulled from the managed Redis plugin (Part 5) instead of a container
REDIS_URL=${{Redis.REDIS_URL}}

# Same public-origin / CORS logic as the Compose file, now pointed at
# real domains instead of Docker network hostnames
APP_URL=https://api.sentex.gritcms.com
CORS_ORIGINS=https://sentex.gritcms.com,https://admin.sentex.gritcms.com
```

Then add every other application variable your `.env.production` defines
that isn't Postgres/Redis/routing related — for example your R2 storage
credentials, since those are unrelated to Compose networking and carry over
as plain values:

```
R2_ACCOUNT_ID=...
R2_ACCESS_KEY_ID=...
R2_SECRET_ACCESS_KEY=...
R2_BUCKET=...
R2_ENDPOINT=...
```

Copy across any remaining app-specific secrets (JWT signing keys, OAuth
credentials if social auth is enabled, mail provider keys, etc.) from your
local `.env.production` the same way — one `KEY=VALUE` line per variable.

### 6.4 Replace the `migrate` job with a Pre-Deploy Command

1. Still in the `api` service, go to **Settings** → **Deploy**.
2. Find **Pre-Deploy Command** and set it to:
   ```
   sh -c "./migrate && ./seed"
   ```
   Wrap it in `sh -c "..."` exactly like this — because `api` builds from a
   Dockerfile (not Railpack), Railway executes the pre-deploy command
   directly rather than through a shell, so `&&` needs an explicit shell to
   be interpreted correctly. Without the `sh -c` wrapper the command fails
   with an "exec format" style error.
3. Leave **Custom Start Command** as whatever `apps/api/Dockerfile`'s
   `CMD`/`ENTRYPOINT` already runs your compiled API binary with — you don't
   need to override it, since the Compose file's `api` service doesn't
   override `command:` either.
4. This reproduces the ordering guarantees from your Compose file
   (`migrate` → seed → then `api` starts serving) without needing a
   separate one-shot service: the pre-deploy command runs in its own
   container on every deploy, and if it exits non-zero the new `api`
   version never goes live. It is idempotent on your side already (the seed
   skips when members exist), so redeploys are safe exactly as documented
   in the Compose file's header comment.

### 6.5 Networking

1. Still under **Settings**, find **Networking**.
2. Click **Generate Domain** if you want a temporary `*.up.railway.app` URL
   to test with before the custom domain is live. This is optional — you
   can also wait and go straight to the custom domain in Part 8.
3. You do **not** need to expose port 8080 manually the way `expose:
   ["8080"]` did in Compose — Railway detects the port your app listens on,
   or you can pin it explicitly under **Settings → Networking → Port**.
4. Trigger a redeploy (**Deploy** button, top right) now that Root
   Directory, Variables, and the Pre-Deploy Command are all set. Watch the
   **Deploy Logs** — you should see the pre-deploy container run
   `./migrate && ./seed` and exit 0, then the `api` container start.

---

## Part 7 — Deploy the `admin` service

### 7.1 Create the service

1. Click **+ New** → **GitHub Repo** → select the same `sentex` repository
   and branch again. Railway allows multiple services from the same repo.
2. Rename the service to `admin`.

### 7.2 Point it at the right Dockerfile

Unlike `api`, the Compose file builds `admin` with `context: .` (the repo
root) but `dockerfile: apps/admin/Dockerfile` — it needs the repo root as
build context so it can see the shared `packages/` workspace and root
`pnpm-lock.yaml`, but the Dockerfile itself lives one level down.

1. Go to **Settings** → **Source**.
2. Leave **Root Directory** **blank** (i.e. the repo root) — do not set it
   to `apps/admin`, or the build will lose access to `packages/` and
   `pnpm-lock.yaml` and fail during `pnpm install`.
3. Under **Build Configuration**, set **Dockerfile Path** to:
   ```
   apps/admin/Dockerfile
   ```
4. Under **Build**, set **Watch Paths** to:
   ```
   apps/admin/**
   packages/**
   ```

### 7.3 Set build-time variables

These map directly to the `build.args` block for `admin` in your Compose
file. Because `apps/admin/Dockerfile` declares them with `ARG`, setting
them as ordinary **Variables** on this service (not anything special) is
enough for Railway to forward them into the Docker build automatically —
that's what "Using variables at build time" means in Railway's Dockerfile
docs, as long as the `ARG` line exists in the Dockerfile.

```
NEXT_PUBLIC_API_URL=https://api.sentex.gritcms.com
NEXT_PUBLIC_WEB_URL=https://sentex.gritcms.com
NEXT_PUBLIC_ADMIN_URL=https://admin.sentex.gritcms.com
THEME=atlas
SOCIAL_AUTH_ENABLED=false
NEXT_PUBLIC_DEMO_LOGINS=false
```

Adjust `THEME`, `SOCIAL_AUTH_ENABLED`, and `NEXT_PUBLIC_DEMO_LOGINS` to your
actual production values — the ones above just mirror the safe defaults
called out in the Compose file's comments (social auth off, demo logins
off, since this stack has no OAuth provider configured). Remember these are
**build-time**: changing them later requires a new deploy/rebuild, not just
a service restart, exactly as your Compose file's comments warn.

### 7.4 Runtime variables and networking

1. Add any runtime-only variables `admin` needs beyond the build args (for
   example, session secrets, if the admin app reads them at request time
   rather than build time).
2. Under **Settings → Networking**, click **Generate Domain** for a
   temporary test URL, or skip straight to Part 8 for the real domain.
3. Deploy the service and check **Deploy Logs** for a successful Next.js
   start on port 3000.

---

## Part 8 — Deploy the `web` service

Repeat Part 7 with these differences:

1. Rename the service to `web`.
2. **Root Directory**: blank (repo root) — same reasoning as `admin`.
3. **Dockerfile Path**: `apps/web/Dockerfile`.
4. **Watch Paths**:
   ```
   apps/web/**
   packages/**
   ```
5. Build-time variables (matches the `web` service's `build.args` — note it
   has no `NEXT_PUBLIC_WEB_URL` arg of its own, since `web` *is* the web
   app):
   ```
   NEXT_PUBLIC_API_URL=https://api.sentex.gritcms.com
   NEXT_PUBLIC_ADMIN_URL=https://admin.sentex.gritcms.com
   THEME=atlas
   SOCIAL_AUTH_ENABLED=false
   ```
6. Generate a domain / deploy, and confirm the Deploy Logs show a
   successful start on port 3000.

---

## Part 9 — Attach your custom domains

Do this once `api`, `admin`, and `web` have each successfully deployed at
least once.

1. Open the `api` service → **Settings → Networking → Custom Domain**.
2. Enter `api.sentex.gritcms.com` and click **Add**.
3. Railway shows you a CNAME target (something like
   `xxxx.up.railway.app`). Copy it.
4. In your DNS provider, create a **CNAME** record:
   `api.sentex.gritcms.com → xxxx.up.railway.app`
   (replacing the A record placeholder from Part 2).
5. Repeat steps 1–4 for `admin` (`admin.sentex.gritcms.com`) and `web`
   (`sentex.gritcms.com`) on their respective services.
6. Wait for DNS to propagate, then confirm each domain in the Railway
   dashboard shows a green "Active"/verified TLS status. Railway
   provisions and renews the TLS certificate automatically once DNS
   resolves correctly — you don't manage Let's Encrypt yourself the way the
   Compose file's Traefik labels did.

**Important:** if you generated a domain, deployed, and *then* changed
`NEXT_PUBLIC_API_URL` / `NEXT_PUBLIC_ADMIN_URL` / `NEXT_PUBLIC_WEB_URL` to
match the final custom domains, you must trigger a fresh deploy of `admin`
and `web` after adding the custom domains — those values are baked into the
JS bundle at build time, so the bundle built against the temporary
`*.up.railway.app` URL will not automatically pick up the custom domain.

---

## Part 10 — Verify the full deployment

1. **Postgres / Redis**: open each service's **Data** tab (Postgres) or
   **Metrics** tab and confirm they show as running/healthy.
2. **api**: open Deploy Logs, confirm you see the pre-deploy command run
   `./migrate && ./seed` and exit cleanly, then the API start message. Hit
   `https://api.sentex.gritcms.com/<your health endpoint>` with `curl` and
   confirm a 200.
3. **admin**: visit `https://admin.sentex.gritcms.com` in a browser, open
   dev tools → Network tab, and confirm requests go to
   `api.sentex.gritcms.com` with no CORS errors in the console.
4. **web**: visit `https://sentex.gritcms.com` and repeat the same CORS/API
   check.
5. Log in with the demo SACCO credentials seeded by the `migrate`/`seed`
   step to confirm the database was actually populated.

---

## Part 11 — Ongoing deploys

Because every service is connected via **GitHub Repo** (not the CLI), the
day-to-day workflow is now:

1. Push a commit to the connected branch (e.g. `main`).
2. Railway's GitHub webhook triggers a build for every service whose
   **Watch Paths** matched the changed files.
3. For `api`, the pre-deploy command reruns automatically on every deploy
   — safe, since your migration/seed are idempotent.
4. Watch **Deploy Logs** per service if anything looks off; roll back to a
   previous deployment from the service's **Deployments** tab if needed
   (Railway keeps deployment history and supports one-click rollback).

You generally won't need the Railway CLI for this project's day-to-day
deploys — it's most useful for local `railway run` / `railway variables`
debugging against the same environment, or `railway logs` to tail a
service without opening the dashboard.

---

## Troubleshooting

| Symptom                                                              | Likely cause                                                                                                   |
| ----------------------------------------------------------------------| ------------------------------------------------------------------------------------------------------------- |
| `admin`/`web` build fails at `pnpm install`, can't find lockfile      | Root Directory was set to `apps/admin`/`apps/web` instead of left blank — it needs the repo root as context.  |
| Build-time `NEXT_PUBLIC_*` value doesn't show up in the shipped app   | Missing `ARG` for that variable in the Dockerfile, or the variable was set after the last build (rebuild it). |
| CORS errors in browser console                                       | `CORS_ORIGINS` on `api` doesn't exactly match the live `admin`/`web` domains (scheme + host, no trailing slash).|
| Pre-deploy command fails with an exec/format error                   | Missing the `sh -c "..."` wrapper — required for Dockerfile-based services.                                    |
| `api` can't reach Postgres/Redis                                     | Reference variables (`${{Postgres.PGHOST}}`, etc.) weren't set, or the plugin service was renamed after wiring.|
| A push to `apps/web` triggers an `api` rebuild too                   | Watch Paths weren't set (or were set too broadly) on one of the services.                                      |
| Custom domain stuck "pending"                                        | DNS CNAME not propagated yet, or it's pointed at the wrong Railway target — recheck against Part 9.            |
