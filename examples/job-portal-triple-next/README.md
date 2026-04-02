# Job Portal -- Triple Architecture + Next.js

A full-featured job portal built with Grit using the **Triple** architecture:
- `apps/api` -- Go API (Gin + GORM + PostgreSQL)
- `apps/web` -- Next.js public website (job listings, company profiles, apply)
- `apps/admin` -- Next.js admin panel (manage jobs, companies, applications)

## Architecture

This is the **most feature-rich** Grit configuration. Three apps share types
and schemas via `packages/shared`. The admin panel provides a full back-office
for managing content.

Best for: SaaS platforms, marketplaces, content-heavy apps.

## Quick Start

### Prerequisites
- Go 1.21+
- Node.js 18+
- pnpm 8+
- Docker (for PostgreSQL, Redis, MinIO, Mailhog)

### Setup

1. Install Grit CLI:
   ```bash
   go install github.com/MUKE-coder/grit/v3/cmd/grit@latest
   ```

2. Create the project:
   ```bash
   grit new job-portal --triple --next
   ```

3. Start Docker services:
   ```bash
   cd job-portal
   docker compose up -d
   ```

4. Generate resources:
   ```bash
   grit generate resource Category --fields "name:string:unique,slug:slug:name"
   grit generate resource Company --fields "name:string,slug:slug:name,description:richtext,logo:string,website:string:optional,email:string"
   grit generate resource Job --fields "title:string,description:richtext,company_id:belongs_to:Company,location:string,salary_min:float,salary_max:float,type:string,experience:string,active:bool,deadline:date"
   grit generate resource Application --fields "job_id:belongs_to:Job,applicant_name:string,applicant_email:string,resume:string,cover_letter:text:optional,status:string"
   ```

5. Run migrations and seed:
   ```bash
   grit migrate
   grit seed
   ```

6. Start development servers:
   ```bash
   # Terminal 1: API
   cd apps/api && go run cmd/server/main.go

   # Terminal 2: All frontends
   pnpm install && pnpm dev
   ```

7. Open:
   - Web: http://localhost:3000
   - Admin: http://localhost:3001
   - API Docs: http://localhost:8080/docs
   - GORM Studio: http://localhost:8080/studio
   - Mailhog: http://localhost:8025

## Project Structure

```
job-portal/
├── apps/
│   ├── api/                    # Go backend
│   │   ├── cmd/server/main.go  # Entry point
│   │   └── internal/
│   │       ├── models/         # Job, Company, Application, Category, User
│   │       ├── handlers/       # REST handlers
│   │       ├── services/       # Business logic
│   │       ├── middleware/     # Auth, CORS, cache, rate limiting
│   │       └── routes/         # Route registration
│   ├── web/                    # Next.js public site
│   │   └── app/
│   │       ├── page.tsx        # Landing page (hero + featured jobs)
│   │       ├── jobs/           # Job listings + detail pages
│   │       ├── companies/      # Company profiles
│   │       └── apply/          # Application form
│   └── admin/                  # Next.js admin panel
│       └── resources/
│           ├── jobs.ts         # Job management (DataTable + FormBuilder)
│           ├── companies.ts    # Company management
│           ├── applications.ts # Application review + status updates
│           └── categories.ts   # Category management
├── packages/shared/            # Shared types + schemas
│   ├── schemas/                # Zod validation
│   └── types/                  # TypeScript interfaces
├── .env                        # Environment variables
├── .env.example                # Template for other developers
├── docker-compose.yml          # Dev services
└── docker-compose.prod.yml     # Production deployment
```

## Features Demonstrated

### Authentication
- Email/password registration and login
- OAuth2 social login (Google, GitHub)
- JWT access tokens (15 min) + refresh tokens (7 days)
- TOTP two-factor authentication with backup codes
- Role-based access: ADMIN (full access), EDITOR (manage jobs), USER (apply)

### Data Management (Admin Panel)
- DataTable with server-side pagination, sorting, filtering
- FormBuilder with rich text editor for job descriptions
- File upload for company logos and resumes
- Application status workflow: Pending -> Reviewed -> Interview -> Accepted/Rejected
- Export to CSV/Excel

### Public Website
- Landing page with hero, featured jobs, and category cards
- Job listings with search, location filter, and salary range filter
- Company profiles with job listings
- Application form with resume upload
- SEO-optimized pages (Next.js SSR)

### Emails
- Welcome email on registration
- Application confirmation email
- Application status update notifications

### Deployment
See the Deployment section below.

## Environment Variables

See `.env.example` for the full list.

Key variables:
- `DATABASE_URL` -- PostgreSQL connection string
- `JWT_SECRET` -- Must be changed in production
- `STORAGE_ENDPOINT` -- MinIO (dev) or S3/R2 (prod)
- `RESEND_API_KEY` -- For sending emails
- `AI_GATEWAY_API_KEY` -- For AI features (optional)
- `GOOGLE_CLIENT_ID` -- For OAuth (optional)
- `GITHUB_CLIENT_ID` -- For OAuth (optional)

## Deployment

### Option 1: grit deploy (VPS)
```bash
grit deploy --host deploy@server.com --domain jobs.example.com
```

### Option 2: Docker Compose
```bash
docker compose -f docker-compose.prod.yml up -d
```

### Option 3: Vercel (frontend) + VPS (API)
```bash
# Deploy API
grit deploy --host deploy@server.com --domain api.jobs.example.com

# Deploy web to Vercel
cd apps/web && vercel --env NEXT_PUBLIC_API_URL=https://api.jobs.example.com
```

## LLM Instructions

If you are an AI assistant working on this project:
1. Read `.claude/skills/grit/SKILL.md` for Grit conventions
2. Use `grit generate resource` for new resources
3. Never delete `// grit:*` marker comments
4. Follow the API response format: `{ data, message }` / `{ data, meta }` / `{ error: { code, message } }`
5. Use React Query for data fetching, Zod for validation
6. Report bugs at https://github.com/MUKE-coder/grit/issues

## License

MIT
