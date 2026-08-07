# Deploying Sentex on AWS EC2

Unlike Railway/Render/Fly, EC2 gives you a bare Linux box — there's no
platform-level "connect GitHub" or managed proxy. You run
`docker-compose.prod.yml` almost exactly as-is (real Docker Compose, so
`depends_on: condition: service_completed_successfully` for `migrate`
works natively, same as the Dokploy/Coolify VPS guides), and you build the
public-facing, TLS-terminating, host-based-routing layer yourself using an
AWS **Application Load Balancer (ALB)** + **ACM** — the standard
production pattern on AWS, and the reason this guide looks different from
the plain-VPS ones.

## Architecture at a glance

```
Internet → ALB (443, ACM cert, host-header routing) → EC2 instance (3 host ports)
                                                          ├─ :8080 → api container
                                                          ├─ :3001 → admin container
                                                          └─ :3000 → web container
EC2 instance also runs postgres + redis, reachable only inside the
instance's own Docker network (no host port, no ALB route to them).
```

The ALB replaces Traefik/Dokploy's role from the VPS guides: it terminates
TLS with an AWS-managed certificate and routes each hostname to a different
container port on the same instance.

---

## Part 1 — Network and security groups

1. In the VPC you'll deploy into (the default VPC is fine for a first
   pass), confirm you have **at least two public subnets in two different
   Availability Zones** — the ALB requires this even though your EC2
   instance only runs in one AZ.
2. Create a security group **`sentex-alb-sg`**:
   - Inbound: HTTP (80) from `0.0.0.0/0`, HTTPS (443) from `0.0.0.0/0`.
   - Outbound: all traffic.
3. Create a security group **`sentex-ec2-sg`**:
   - Inbound: SSH (22) from your IP only (not `0.0.0.0/0`).
   - Inbound: custom TCP `8080`, `3000`, `3001` — **source: `sentex-alb-sg`**
     (not an IP range). This is what keeps the app containers unreachable
     from the internet except through the ALB.
   - Outbound: all traffic.

---

## Part 2 — Launch the EC2 instance

1. EC2 → **Launch Instance**.
2. AMI: **Ubuntu Server 24.04 LTS**.
3. Instance type: **t3.medium** (2 vCPU / 4 GB) minimum — building three
   Docker images (two of them Next.js/pnpm builds) on `t3.micro`/`t3.small`
   routinely OOMs. Scale up later once you know your real traffic.
4. Key pair: create or select one — you'll need it for SSH.
5. Network settings: the VPC/subnet from Part 1, **auto-assign public IP:
   enabled** (needed for SSH and for `git`/`docker pull` egress; the app
   ports themselves stay locked to ALB-only per the security group).
6. Security group: attach **`sentex-ec2-sg`**.
7. Storage: bump the root volume to at least **40 GB gp3** — three Docker
   images plus Postgres/Redis data adds up fast.
8. Launch. Once running, allocate an **Elastic IP** and associate it to the
   instance, so the public IP doesn't change on stop/start (useful for SSH
   convenience; the ALB, not this IP, is what your DNS will point at).

---

## Part 3 — Install Docker on the instance

SSH in:
```
ssh -i /path/to/key.pem ubuntu@<instance-public-ip>
```

Install Docker Engine + Compose plugin:
```
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh
sudo usermod -aG docker $USER
newgrp docker
docker compose version   # confirm the plugin is present
```

(Optional but recommended on `t3.medium`) add a swap file so pnpm/Next.js
builds don't get OOM-killed under memory pressure:
```
sudo fallocate -l 2G /swapfile
sudo chmod 600 /swapfile
sudo mkswap /swapfile
sudo swapon /swapfile
echo '/swapfile none swap sw 0 0' | sudo tee -a /etc/fstab
```

---

## Part 4 — Get the code and adjust the Compose file for a bare VM

1. Clone the repo:
   ```
   git clone https://github.com/<you>/sentex.git
   cd sentex
   ```
2. `docker-compose.prod.yml` was written for Dokploy's Traefik + external
   `dokploy-network`. Neither exists here — the ALB replaces Traefik, so
   make a `docker-compose.ec2.yml` copy with:
   - **Remove** the `traefik.*` labels from `api`, `admin`, `web` (nothing
     is watching for them on a bare instance).
   - **Remove** the top-level `dokploy-network` entry and drop it from
     every service's `networks:` list — keep only the internal `sentex`
     network.
   - **Add host port mappings** so the ALB's target groups have something
     to hit:
     ```yaml
     services:
       api:
         ports:
           - "8080:8080"
       admin:
         ports:
           - "3001:3000"
       web:
         ports:
           - "3000:3000"
     ```
     `postgres` and `redis` get **no** `ports:` entry — they stay reachable
     only via the internal `sentex` network at hostnames `postgres`/`redis`,
     exactly as your app's `POSTGRES_HOST=postgres` / `REDIS_URL` values
     already assume.
3. Everything else — `build:`, `env_file:`, `environment:`, `volumes:`,
   `depends_on:`, `healthcheck:`, the `migrate` service's `command:` and
   `restart: "no"` — stays **unchanged**. Real Docker Compose on a real VM
   honors `service_completed_successfully` natively, so `migrate` still
   runs to completion before `api` starts, with no pre-deploy-command
   workaround needed.

## Part 5 — Environment variables

1. Create `.env` in the repo root on the instance (do **not** commit it):
   ```
   nano .env
   ```
2. Paste the contents of your local `.env.production`, with `WEB_DOMAIN`,
   `ADMIN_DOMAIN`, `API_DOMAIN` set to the real hostnames:
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
   `docker compose` reads `.env` in the working directory automatically and
   interpolates every `${VARIABLE}` in the Compose file — same mechanism
   the file's own `env_file: [.env]` entries rely on for `api`/`migrate`.

## Part 6 — First deploy

```
docker compose -f docker-compose.ec2.yml --env-file .env up -d --build
```

Watch it come up:
```
docker compose -f docker-compose.ec2.yml logs -f migrate
docker compose -f docker-compose.ec2.yml ps
```
Confirm `migrate` exits 0, then `api`, `admin`, `web` show as running.

Sanity-check locally on the instance before touching DNS/ALB:
```
curl -I http://localhost:8080/<health endpoint>
curl -I http://localhost:3001
curl -I http://localhost:3000
```

---

## Part 7 — Request the TLS certificate (ACM)

1. Open **AWS Certificate Manager** in the **same region** as your ALB
   will live.
2. **Request a certificate** → Public certificate.
3. Domain names — add all three as Subject Alternative Names on one
   certificate:
   ```
   sentex.gritcms.com
   admin.sentex.gritcms.com
   api.sentex.gritcms.com
   ```
4. Validation method: **DNS validation** (faster and auto-renewing).
5. ACM shows a CNAME record per domain. Create these at your DNS provider.
   If your zone is hosted in Route 53, ACM offers a **Create records in
   Route 53** button that does this for you.
6. Wait for all three domains to show **Issued** (usually a few minutes
   once the CNAMEs resolve).

---

## Part 8 — Target groups

Create three **target groups**, type **Instances**, protocol **HTTP**,
in the same VPC:

| Name | Port | Health check path |
|---|---|---|
| `sentex-api-tg` | 8080 | `/<your api health endpoint>` |
| `sentex-admin-tg` | 3001 | `/` |
| `sentex-web-tg` | 3000 | `/` |

For each: on the **Register targets** step, select your EC2 instance and,
importantly, set the **port override** to that target group's port (they
all point at the same instance, just different ports) before clicking
**Include as pending below** → **Register pending targets**.

---

## Part 9 — Create the Application Load Balancer

1. EC2 → **Load Balancers** → **Create load balancer** → **Application
   Load Balancer**.
2. Scheme: **Internet-facing**. IP type: IPv4.
3. VPC: same as the instance. Mappings: select the two+ public subnets
   from Part 1.
4. Security group: **`sentex-alb-sg`**.
5. Listeners:
   - **HTTP:80** — you'll edit this after creation to redirect to HTTPS
     (step 7 below).
   - **HTTPS:443** — default certificate: the ACM cert from Part 7.
     Default action: pick any target group for now (e.g. `sentex-web-tg`)
     — you'll override per-hostname with rules next.
6. Create the load balancer and wait for its state to become **Active**.
7. On the **HTTP:80** listener, edit the default rule to **Redirect to
   HTTPS://#{host}:443/#{path}?#{query}** (status 301) instead of
   forwarding — this makes plain-HTTP requests upgrade automatically.
8. On the **HTTPS:443** listener, add rules (in order, above the default
   action):
   - **IF** Host header is `api.sentex.gritcms.com` → **THEN** forward to
     `sentex-api-tg`
   - **IF** Host header is `admin.sentex.gritcms.com` → **THEN** forward
     to `sentex-admin-tg`
   - **IF** Host header is `sentex.gritcms.com` → **THEN** forward to
     `sentex-web-tg`
   - Default action (no rule matched): forward to `sentex-web-tg`, or
     return a fixed 404 — your call.

---

## Part 10 — DNS

1. Note the ALB's DNS name (something like
   `sentex-alb-123456789.us-east-1.elb.amazonaws.com`), shown on the load
   balancer's detail page.
2. Create three records at your DNS provider:
   - If hosted in **Route 53**: create **A records with Alias** target =
     the ALB, for all three hostnames (Alias records work for subdomains
     and are free of the CNAME-at-apex restriction, and Route 53 resolves
     them without an extra DNS lookup).
   - If hosted **elsewhere**: create **CNAME** records for all three
     hostnames pointing at the ALB's DNS name (fine here since none of the
     three is a bare apex domain).
3. Wait for propagation, then confirm:
   ```
   dig +short api.sentex.gritcms.com
   dig +short admin.sentex.gritcms.com
   dig +short sentex.gritcms.com
   ```

---

## Part 11 — Verify

1. Visit `https://api.sentex.gritcms.com/<health endpoint>`,
   `https://admin.sentex.gritcms.com`, `https://sentex.gritcms.com` — all
   should load over a valid ACM certificate with no browser warnings.
2. Open dev tools on `admin`/`web`, confirm API calls succeed with no CORS
   errors (double check `CORS_ORIGINS` in `.env` matches exactly).
3. Log in with the seeded demo credentials to confirm `migrate`/`seed`
   populated the database.
4. In the ALB console, confirm all three target groups show their
   registered target as **healthy**.

---

## Part 12 — Ongoing deploys

There's no GitHub webhook wired up out of the box on a raw EC2 box — pick
one:

**Manual** (fine for a solo project):
```
ssh ubuntu@<instance-ip>
cd sentex && git pull
docker compose -f docker-compose.ec2.yml --env-file .env up -d --build
```

**GitHub Actions** (push-to-deploy, closer to the other guides' workflow):
```yaml
name: Deploy to EC2
on:
  push:
    branches: [main]
jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: appleboy/ssh-action@v1
        with:
          host: ${{ secrets.EC2_HOST }}
          username: ubuntu
          key: ${{ secrets.EC2_SSH_KEY }}
          script: |
            cd sentex
            git pull
            docker compose -f docker-compose.ec2.yml --env-file .env up -d --build
```
Store the instance's SSH private key as `EC2_SSH_KEY` and its Elastic IP
as `EC2_HOST` in the repo's GitHub Actions secrets.

`migrate` reruns on every `up -d --build` and is idempotent (per the
Compose file's own comments), so this is safe to run on every deploy.

---

## Notes and troubleshooting

| Symptom | Likely cause |
|---|---|
| ALB target shows "unhealthy" | Health check path returns non-2xx, or the security group doesn't allow the ALB SG on that port — recheck Part 1 step 3. |
| 504 from the ALB | Container isn't listening on the port the target group expects — confirm with `docker compose ps` and `curl localhost:<port>` on the instance. |
| Build gets OOM-killed | Instance too small for concurrent Next.js builds — add swap (Part 3) or size up to `t3.large` temporarily for the first build. |
| CORS errors in browser | `CORS_ORIGINS` in `.env` doesn't exactly match the live domains (scheme + host, no trailing slash) — redeploy `api` after fixing. |
| Cheaper alternative to the ALB | If you don't need AWS-native routing/WAF/autoscaling, you can skip Parts 7–10 entirely and instead run a self-managed reverse proxy (e.g. Caddy) directly on the EC2 instance with automatic Let's Encrypt certs — see the Lightsail guide's Part 6 for the exact pattern, which works identically on a plain EC2 box. |
