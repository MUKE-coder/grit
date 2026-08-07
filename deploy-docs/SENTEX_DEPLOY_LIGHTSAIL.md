# Deploying Sentex on AWS Lightsail

Lightsail is AWS's simplified VPS product — same underlying EC2/Docker
mechanics as the plain-VM guides, but with a flat monthly price, a
simplified firewall UI, a static IP included, and one-click snapshots.
The one thing to know going in: **Lightsail's own load balancer cannot do
host-based routing** (it forwards everything to one instance port, no
`Host:` header rules) — so instead of an AWS load balancer in front, this
guide runs a small **Caddy** reverse-proxy container on the instance
itself, which does host-based routing *and* gets you free, auto-renewing
HTTPS with zero manual certificate steps. It's the same job Traefik does
in the Dokploy guide, just self-managed instead of platform-managed.

## Architecture at a glance

```
Internet → Lightsail static IP → Caddy container (80/443, auto HTTPS)
                                    ├─ Host: api.sentex.gritcms.com   → api:8080
                                    ├─ Host: admin.sentex.gritcms.com → admin:3000
                                    └─ Host: sentex.gritcms.com       → web:3000
postgres + redis: no public/host ports, reachable only on the
instance's internal Docker network.
```

---

## Part 1 — Create the instance

1. Lightsail console → **Create instance**.
2. Platform: **Linux/Unix**. Blueprint: **OS Only → Ubuntu 24.04 LTS**
   (skip the app blueprints — you're bringing your own Docker Compose
   stack).
3. Instance plan: pick **at least the 4 GB RAM / 2 vCPU plan**. The 2 GB
   plan can OOM while building two Next.js apps and a Go API back-to-back
   in the same `docker compose build` run.
4. Name it (e.g. `sentex-prod`), choose your region, and create it.
5. Wait for the instance state to show **Running**.

## Part 2 — Static IP

1. In the instance's **Networking** tab, click **Create static IP**,
   attach it to `sentex-prod`.
2. Note the static IP — this is what your DNS records will point at.
   (Skip Part 2 entirely if you plan to put a Lightsail load balancer in
   front — see the note at the end of this guide — but for the
   recommended Caddy-on-instance setup below, you want the static IP.)

## Part 3 — Firewall (networking tab)

Lightsail's instance firewall is a simplified security group. Configure:

| Application | Protocol | Port | Source |
|---|---|---|---|
| SSH | TCP | 22 | Restrict to your IP |
| HTTP | TCP | 80 | Anywhere |
| HTTPS | TCP | 443 | Anywhere |

Do **not** open 8080/3000/3001 — the app containers stay internal, only
reachable through Caddy on 80/443, same principle as locking the EC2
security group to ALB-only in the EC2 guide.

## Part 4 — Install Docker

Connect via the Lightsail browser SSH button, or:
```
ssh -i /path/to/your-lightsail-key.pem ubuntu@<static-ip>
```

Install Docker Engine + Compose plugin:
```
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh
sudo usermod -aG docker $USER
newgrp docker
docker compose version
```

(Recommended on the 4 GB plan) add 2 GB of swap the same way as the EC2
guide, so concurrent builds don't get killed under memory pressure:
```
sudo fallocate -l 2G /swapfile
sudo chmod 600 /swapfile
sudo mkswap /swapfile
sudo swapon /swapfile
echo '/swapfile none swap sw 0 0' | sudo tee -a /etc/fstab
```

## Part 5 — Get the code

```
git clone https://github.com/<you>/sentex.git
cd sentex
```

## Part 6 — Add the Caddy reverse proxy and adjust the Compose file

Create `docker-compose.lightsail.yml` as a copy of
`docker-compose.prod.yml` with these changes:

1. **Remove** the `traefik.*` labels from `api`, `admin`, `web` — nothing
   is watching for them here.
2. **Remove** the `dokploy-network` entry (top-level and from every
   service's `networks:` list) — keep only the internal `sentex` network.
   Add `caddy` to that same network so it can reach the app containers by
   service name.
3. **Add** a `caddy` service:
   ```yaml
   services:
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
     caddy-config:
   ```
4. Create `Caddyfile` in the repo root:
   ```
   api.sentex.gritcms.com {
       reverse_proxy api:8080
   }

   admin.sentex.gritcms.com {
       reverse_proxy admin:3000
   }

   sentex.gritcms.com {
       reverse_proxy web:3000
   }
   ```
   Caddy requests and renews a Let's Encrypt certificate for each site
   block automatically the first time it starts and each domain resolves
   — no `certbot`, no cert-renewal cron job, no ACM console. This is the
   single biggest reason to reach for Caddy over plain nginx here.

Everything else in the Compose file — `build:`, `env_file:`,
`environment:`, `volumes:`, `depends_on:`, `healthcheck:`, and the
`migrate` service's `command: ["sh", "-c", "./migrate && ./seed"]` with
`restart: "no"` — stays **unchanged**. Real Docker Compose on a real VM
honors `depends_on: condition: service_completed_successfully` natively,
so `migrate` still blocks `api` from starting until it exits 0.

## Part 7 — Environment variables

```
nano .env
```
Paste your production values, same as the EC2 guide:
```
WEB_DOMAIN=sentex.gritcms.com
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
R2_ENDPOINT=...
```

## Part 8 — DNS (before you deploy)

Caddy's automatic HTTPS validates ownership over the internet (HTTP-01
challenge on port 80), so DNS has to resolve **before** the containers
start, unlike the EC2/ALB flow where you could issue the cert independently
via DNS validation first.

Create three **A records** at your DNS provider, pointing at the Lightsail
static IP from Part 2:
```
sentex.gritcms.com          → <static IP>
admin.sentex.gritcms.com    → <static IP>
api.sentex.gritcms.com      → <static IP>
```
Wait for propagation (`dig +short sentex.gritcms.com` should return the
static IP) before moving to Part 9, or Caddy's first certificate request
will fail and retry on a backoff.

## Part 9 — First deploy

```
docker compose -f docker-compose.lightsail.yml --env-file .env up -d --build
```

Watch `migrate` complete, then watch Caddy issue certificates:
```
docker compose -f docker-compose.lightsail.yml logs -f migrate
docker compose -f docker-compose.lightsail.yml logs -f caddy
```
You're looking for lines like `certificate obtained successfully` for each
of the three domains in the Caddy log.

## Part 10 — Verify

1. Visit all three domains over `https://` — valid Let's Encrypt cert, no
   warnings.
2. Confirm no CORS errors in the browser console on `admin`/`web`.
3. Log in with the seeded demo credentials to confirm `migrate`/`seed` ran.
4. `docker compose -f docker-compose.lightsail.yml ps` — all containers
   should show healthy/running, `migrate` should show `Exited (0)`.

## Part 11 — Backups

Lightsail's headline convenience feature: turn on **Automatic snapshots**
under the instance's **Snapshots** tab. This snapshots the entire instance
disk daily (including your Postgres data directory's Docker volume, since
it lives on the instance's own disk) — a much lower-effort baseline than
configuring `pg_dump` cron jobs yourself, though for a real production
system you'll still want logical Postgres backups (`pg_dump`) in addition,
since a full-disk snapshot restores the *whole instance* to a point in
time, not just the database.

## Part 12 — Ongoing deploys

Same two options as the EC2 guide — manual:
```
ssh ubuntu@<static-ip>
cd sentex && git pull
docker compose -f docker-compose.lightsail.yml --env-file .env up -d --build
```

or GitHub Actions:
```yaml
name: Deploy to Lightsail
on:
  push:
    branches: [main]
jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: appleboy/ssh-action@v1
        with:
          host: ${{ secrets.LIGHTSAIL_STATIC_IP }}
          username: ubuntu
          key: ${{ secrets.LIGHTSAIL_SSH_KEY }}
          script: |
            cd sentex
            git pull
            docker compose -f docker-compose.lightsail.yml --env-file .env up -d --build
```

---

## Alternatives worth knowing about

**Lightsail Load Balancer instead of Caddy.** Lightsail does offer a
managed load balancer with free auto-renewing certificates (up to 9
alternate domains on one cert) — but it forwards *all* traffic to a single
instance port with no host-header routing. To use it here you'd still need
a reverse proxy on the instance doing the host-based split; the load
balancer would just replace Caddy's TLS termination (forwarding plain
HTTP to instance port 80) while an nginx/Caddy container on the instance
did the routing over HTTP internally. For a single-instance deployment
like this one, that's extra cost and moving parts for the same result the
Caddy-only setup already gives you for free — worth it mainly if you later
add a *second* instance and want AWS to load-balance between them.

**Lightsail Container Service.** A separate, fully-managed alternative to
running your own instance — closer in spirit to Railway/Render than to
this guide. It runs container images (not a `docker-compose.yml`) behind
a managed HTTPS endpoint, has its own concept of "public endpoint" per
deployment, and pairs with Lightsail's own managed database offering
instead of a `postgres` container. If you'd rather not manage an instance
at all, that's a different (and shorter) path than this guide — but it
means giving up the one-file-does-everything simplicity of deploying
`docker-compose.prod.yml` more or less as written.

**An AWS ALB in front of a Lightsail instance.** Also possible (via VPC
peering between the Lightsail and EC2 networking planes), which gets you
the exact host-based-routing pattern from the EC2 guide while keeping your
compute on cheaper Lightsail pricing. More moving parts than either guide
alone — only worth it if you specifically want ALB features (WAF,
weighted routing, etc.) without leaving Lightsail entirely.

---

## Notes and troubleshooting

| Symptom | Likely cause |
|---|---|
| Caddy never issues a certificate | DNS wasn't pointing at the static IP yet when Caddy first started — fix DNS, then `docker compose restart caddy`. |
| `Exited (1)` on `migrate` | Check `docker compose logs migrate` — usually a bad `POSTGRES_*` value in `.env`, or Postgres not yet healthy (shouldn't happen given the healthcheck + `depends_on`, but worth checking `docker compose logs postgres` too). |
| Build gets killed partway through | Instance too small — bump to the 4 GB+ plan or add swap (Part 4). |
| 502 from Caddy | The upstream container isn't listening yet or crashed — `docker compose ps` and `docker compose logs api` (or `admin`/`web`). |
| CORS errors in browser | `CORS_ORIGINS` in `.env` doesn't exactly match the live domains — fix and redeploy `api`. |
