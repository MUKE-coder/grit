# Grit Motors

End-to-end stock, motorcycle-sales, and loans platform — Go + React, built on [Grit](https://gritframework.dev).

Two business segments share one platform:
- **Spares** — multi-branch parts inventory + POS
- **Motorcycles** — unit inventory (per number plate), cash sales, loan sales with auto-generated repayment schedules and mobile-money collections via DGateway

## Quick Start

```bash
# 1. Install Air for Go hot reloading
go install github.com/air-verse/air@latest

# 2. Start infrastructure (PostgreSQL, Redis, MinIO, Mailhog)
docker compose up -d

# 3. Install frontend dependencies
pnpm install

# 4. Start all services (API auto-reloads on file changes)
pnpm dev
```

## Project Structure

```
stoka/
├── apps/
│   ├── api/          # Go backend (Gin + GORM)
│   ├── web/          # Next.js frontend
│   └── admin/        # Next.js admin panel
├── packages/
│   └── shared/       # Shared types, schemas, constants
├── docker-compose.yml
└── turbo.json
```

## Services

| Service       | URL                          |
|---------------|------------------------------|
| API           | http://localhost:8080         |
| GORM Studio   | http://localhost:8080/studio  |
| Web App       | http://localhost:3000         |
| Admin Panel   | http://localhost:3001         |
| PostgreSQL    | localhost:5432               |
| Redis         | localhost:6379               |
| MinIO Console | http://localhost:9001         |
| Mailhog       | http://localhost:8025         |

## Development

```bash
# Run Go API with hot reload
cd apps/api && air

# Run Next.js web app
cd apps/web && pnpm dev

# Run admin panel
cd apps/admin && pnpm dev

# Run all services via Turborepo
pnpm dev
```

## No Docker? No Problem

If you can't run Docker, use cloud services instead:

```bash
cp .env.cloud.example .env
```

Then fill in your keys for:
- **[Neon](https://neon.tech)** — PostgreSQL (free tier)
- **[Upstash](https://upstash.com)** — Redis (free tier)
- **[Cloudflare R2](https://dash.cloudflare.com)** — File storage (free tier)
- **[Resend](https://resend.com)** — Email (free tier)

No Docker needed — just your API keys and ``go run``.

## Tech Stack

- **Backend:** Go + Gin + GORM
- **Frontend:** Next.js 14+ (App Router) + React + TypeScript
- **Styling:** Tailwind CSS + shadcn/ui
- **Database:** PostgreSQL
- **Cache:** Redis
- **Monorepo:** Turborepo + pnpm
- **Validation:** Zod (shared schemas)
- **Data Fetching:** React Query (TanStack Query)

---

*Built with Grit v3.27.0*
