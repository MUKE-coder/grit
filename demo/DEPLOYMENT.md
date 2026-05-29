# Grit Motors — Deployment Guide (Dokploy + VPS)

## Architecture

```
┌──────────────────────────────────────┐
│              Dokploy VPS             │
│                                      │
│  ┌────────────┐  ┌────────────────┐  │
│  │  PostgreSQL │  │     Redis      │  │
│  │   :5432     │  │     :6379      │  │
│  └──────┬─────┘  └───────┬────────┘  │
│         │                │           │
│  ┌──────┴────────────────┴────────┐  │
│  │      Grit Motors App         │  │
│  │    Go binary + embedded SPA    │  │
│  │          :8080                 │  │
│  └────────────────────────────────┘  │
│              │                       │
│         Reverse Proxy (Traefik)      │
│         gritdemo.yourdomain.com         │
└──────────────────────────────────────┘
           │
    Cloudflare R2 (images)
    Resend (emails)
```

Single binary serves both API and frontend. PostgreSQL and Redis run as Docker services.

---

## Step 1: Prepare Your VPS

Ensure Dokploy is running on your VPS. You need:
- Docker & Docker Compose
- A domain pointed to your VPS IP (e.g., `gritdemo.yourdomain.com`)

---

## Step 2: Generate Secrets

Run these on your local machine or VPS to generate secure values:

```bash
# JWT Secret (64 chars)
openssl rand -hex 32

# Internal API Key
echo "stk_prod_$(openssl rand -hex 16)"

# Database password
openssl rand -base64 24

# Redis password
openssl rand -base64 24
```

---

## Step 3: Configure Environment

Edit `.env.production` with your real values:

```env
# REQUIRED — change all of these:
APP_URL=https://gritdemo.yourdomain.com
FRONTEND_URL=https://gritdemo.yourdomain.com
DATABASE_URL=postgres://gritdemo:YOUR_DB_PASSWORD@postgres:5432/gritdemo?sslmode=disable
JWT_SECRET=YOUR_64_CHAR_SECRET
INTERNAL_API_KEY=stk_prod_YOUR_KEY
VITE_INTERNAL_API_KEY=stk_prod_YOUR_KEY    # Must match INTERNAL_API_KEY
R2_ENDPOINT=https://YOUR_ACCOUNT_ID.r2.cloudflarestorage.com
R2_ACCESS_KEY=YOUR_KEY
R2_SECRET_KEY=YOUR_SECRET
R2_BUCKET=gritdemo-images
RESEND_API_KEY=re_YOUR_KEY
MAIL_FROM=noreply@yourdomain.com
REDIS_URL=redis://:YOUR_REDIS_PASSWORD@redis:6379
CORS_ORIGINS=https://gritdemo.yourdomain.com
POSTGRES_PASSWORD=YOUR_DB_PASSWORD          # Must match DATABASE_URL
REDIS_PASSWORD=YOUR_REDIS_PASSWORD          # Must match REDIS_URL
```

**Important:** `INTERNAL_API_KEY` and `VITE_INTERNAL_API_KEY` must be identical. The `VITE_` version is baked into the frontend at build time.

---

## Step 4: Deploy on Dokploy

### Option A: Deploy from GitHub (Recommended)

1. In Dokploy, create a new **Compose** application
2. Connect your GitHub repo: `https://github.com/MUKE-coder/grit.git`
3. Set the compose file path to `docker-compose.yml`
4. Add all environment variables from `.env.production` in Dokploy's environment settings
5. Set the domain: `gritdemo.yourdomain.com`
6. Enable auto-SSL (Let's Encrypt via Traefik)
7. Deploy

### Option B: Deploy via Docker Compose directly on VPS

```bash
# Clone the repo
git clone https://github.com/MUKE-coder/grit.git
cd grit/demo

# Copy and edit production env
cp .env.production .env.production.local
nano .env.production.local
# Fill in all real values

# Rename for docker-compose to pick up
cp .env.production.local .env.production

# Build and start
docker compose up -d --build

# Check logs
docker compose logs -f gritdemo
```

---

## Step 5: Run Migrations

After the first deploy, run migrations to create database tables:

```bash
# Via docker exec
docker compose exec gritdemo ./gritdemo migrate

# Or if you have the binary locally:
# DATABASE_URL=postgres://... go run ./cmd/migrate
```

If using Dokploy, you can run this as a one-time command in the container terminal.

**Alternative:** The app auto-migrates on startup if tables don't exist (via GORM AutoMigrate in the seed command). To seed demo data:

```bash
docker compose exec gritdemo ./gritdemo seed
```

---

## Step 6: Verify

1. Visit `https://gritdemo.yourdomain.com` — should show the login page
2. Register a new account or use seeded credentials:
   - Admin: `admin@nakawafashion.com` / `password123`
3. Check the health endpoint: `https://gritdemo.yourdomain.com/api/health`

---

## Step 7: Domain & SSL (Dokploy)

In Dokploy:
1. Go to your application → **Domains**
2. Add domain: `gritdemo.yourdomain.com`
3. Port: `8080`
4. Enable HTTPS (Dokploy handles Let's Encrypt automatically via Traefik)

Make sure your DNS has an **A record** pointing `gritdemo.yourdomain.com` to your VPS IP.

---

## Updating

### From GitHub (Dokploy)
Push to `main` → trigger redeploy in Dokploy (or enable auto-deploy)

### Manual
```bash
cd grit/demo
git pull
docker compose up -d --build
```

---

## Monitoring & Logs

```bash
# All services
docker compose logs -f

# Just the app
docker compose logs -f gritdemo

# Database
docker compose logs -f postgres

# Check running containers
docker compose ps
```

---

## Backup

### Database
```bash
# Dump
docker compose exec postgres pg_dump -U gritdemo gritdemo > backup_$(date +%Y%m%d).sql

# Restore
cat backup.sql | docker compose exec -T postgres psql -U gritdemo gritdemo
```

### Volumes
```bash
# Find volume paths
docker volume inspect gritdemo_postgres_data
docker volume inspect gritdemo_redis_data
```

---

## Troubleshooting

| Issue | Fix |
|-------|-----|
| `connection refused` on DB | Check `DATABASE_URL` uses `postgres` (service name), not `localhost` |
| Frontend shows "Invalid API key" | Ensure `VITE_INTERNAL_API_KEY` matches `INTERNAL_API_KEY` and rebuild |
| CORS errors | Add your domain to `CORS_ORIGINS` |
| Images not loading | Check R2 credentials and bucket name |
| Emails not sending | Verify `RESEND_API_KEY` and `MAIL_FROM` domain is verified in Resend |
| Container keeps restarting | Check `docker compose logs gritdemo` for the error |

---

## Production Checklist

- [ ] Generated unique `JWT_SECRET` (64+ chars)
- [ ] Generated unique `INTERNAL_API_KEY`
- [ ] `VITE_INTERNAL_API_KEY` matches `INTERNAL_API_KEY`
- [ ] Strong `POSTGRES_PASSWORD` and `REDIS_PASSWORD`
- [ ] `APP_URL` and `FRONTEND_URL` set to actual domain
- [ ] `CORS_ORIGINS` set to actual domain
- [ ] R2 credentials configured
- [ ] Resend API key configured, sender domain verified
- [ ] `GORM_STUDIO_ENABLED=false` (or strong password)
- [ ] SSL/HTTPS enabled via Dokploy/Traefik
- [ ] DNS A record pointing to VPS IP
- [ ] Migrations run successfully
- [ ] Can login and create a sale
