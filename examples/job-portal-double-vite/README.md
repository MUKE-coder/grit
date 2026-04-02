# Job Portal -- Double Architecture + TanStack Router (Vite)

A full-featured job portal built with Grit using the **Double** architecture:
- `apps/api` -- Go API (Gin + GORM + PostgreSQL)
- `apps/web` -- React SPA with TanStack Router (Vite)

**No separate admin panel.** ADMIN users manage jobs, companies, and applications
from role-protected pages inside the same web app.

## Why Double?

| | Triple (Web + Admin + API) | Double (Web + API) |
|---|---|---|
| Apps to deploy | 3 | **2** |
| Admin panel | Separate Next.js app | **Built into the web app** |
| Best for | Large SaaS, marketplace | **Blogs, portfolios, simpler apps** |
| Complexity | Higher | **Lower** |
| Shared code | packages/shared | packages/shared |

Choose Double when a dedicated admin panel is overkill. ADMIN users see extra
navigation items and management pages inside the same web app. Regular users
never see those routes.

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
   grit new job-portal --double --vite
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

   # Terminal 2: Web (Vite)
   pnpm install && pnpm dev
   ```

7. Open:
   - Web: http://localhost:5173
   - API: http://localhost:8080
   - GORM Studio: http://localhost:8080/studio
   - Mailhog: http://localhost:8025

## Project Structure

```
job-portal/
├── apps/
│   ├── api/                       # Go backend
│   │   ├── cmd/server/main.go     # Entry point
│   │   └── internal/
│   │       ├── config/            # Environment config
│   │       ├── database/          # GORM setup + migrations
│   │       ├── models/            # Job, Company, Application, Category, User
│   │       ├── handlers/          # REST handlers
│   │       ├── services/          # Business logic
│   │       ├── middleware/        # Auth, CORS, cache, rate limiting
│   │       └── routes/            # Route registration
│   └── web/                       # React SPA (TanStack Router)
│       ├── index.html
│       ├── vite.config.ts
│       ├── tailwind.config.ts
│       ├── package.json
│       └── src/
│           ├── main.tsx           # App entry point
│           ├── routes/            # File-based routing (TanStack Router)
│           │   ├── __root.tsx     # Root layout
│           │   ├── index.tsx      # Landing page
│           │   ├── jobs/
│           │   │   ├── index.tsx  # Job listings
│           │   │   └── $id.tsx    # Job detail
│           │   ├── companies/
│           │   │   ├── index.tsx  # Company directory
│           │   │   └── $slug.tsx  # Company profile
│           │   ├── apply/
│           │   │   └── $jobId.tsx # Application form
│           │   ├── auth/
│           │   │   ├── login.tsx
│           │   │   └── register.tsx
│           │   └── admin/         # Role-protected management pages
│           │       ├── index.tsx  # Admin dashboard
│           │       ├── jobs.tsx   # Job CRUD
│           │       ├── companies.tsx
│           │       ├── applications.tsx
│           │       └── categories.tsx
│           ├── components/        # Shared UI components
│           ├── hooks/             # React Query hooks
│           └── lib/               # API client, utils
├── packages/shared/               # Shared types + Zod schemas
│   ├── schemas/
│   └── types/
├── .env                           # Environment variables
├── .env.example                   # Template
├── docker-compose.yml             # Dev services
├── docker-compose.prod.yml        # Production deployment
├── turbo.json                     # Turborepo config
├── pnpm-workspace.yaml
└── package.json
```

## How Admin Pages Work (No Separate Admin Panel)

In the Double architecture, management pages live inside the web app under
`src/routes/admin/`. These routes are protected by a role check:

```tsx
// src/routes/admin.tsx (layout route)
import { createFileRoute, redirect } from "@tanstack/react-router";
import { useAuth } from "../hooks/use-auth";

export const Route = createFileRoute("/admin")({
  beforeLoad: ({ context }) => {
    if (context.auth.user?.role !== "ADMIN") {
      throw redirect({ to: "/" });
    }
  },
  component: AdminLayout,
});

function AdminLayout() {
  return (
    <div className="flex min-h-screen bg-background">
      <AdminSidebar />
      <main className="flex-1 p-6">
        <Outlet />
      </main>
    </div>
  );
}
```

The navigation bar conditionally renders an "Admin" link for ADMIN users:

```tsx
{user?.role === "ADMIN" && (
  <Link to="/admin" className="text-text-secondary hover:text-foreground">
    Admin
  </Link>
)}
```

This is simpler than running a separate admin app but gives ADMIN users full
management capabilities.

## Features Demonstrated

### Authentication
- Email/password registration and login
- OAuth2 social login (Google, GitHub)
- JWT access tokens (15 min) + refresh tokens (7 days)
- TOTP two-factor authentication with backup codes
- Role-based access: ADMIN (full access), EDITOR (manage jobs), USER (apply)

### Data Management (Admin Pages in Web App)
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
- Client-side routing (TanStack Router SPA)

### Emails
- Welcome email on registration
- Application confirmation email
- Application status update notifications

### TanStack Router Specifics
- File-based routing in `src/routes/`
- Type-safe route params (`$id`, `$slug`, `$jobId`)
- `beforeLoad` guards for auth and role checks
- `__root.tsx` for the root layout
- No `"use client"` directives needed (SPA by default)
- Vite for fast HMR and small production bundles

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

This starts two containers (api + web) behind nginx, plus PostgreSQL, Redis,
and MinIO. No admin container needed.

### Option 3: Static hosting (frontend) + VPS (API)
```bash
# Deploy API
grit deploy --host deploy@server.com --domain api.jobs.example.com

# Build and deploy web to any static host (Vercel, Netlify, Cloudflare Pages)
cd apps/web && pnpm build
# Upload dist/ to your static host
```

Since the frontend is a Vite SPA, it can be served from any static hosting
provider. Set `VITE_API_URL` to your API domain.

## LLM Instructions

If you are an AI assistant working on this project:
1. Read `.claude/skills/grit/SKILL.md` for Grit conventions
2. Use `grit generate resource` for new resources
3. Never delete `// grit:*` marker comments
4. Follow the API response format: `{ data, message }` / `{ data, meta }` / `{ error: { code, message } }`
5. Use React Query for data fetching, Zod for validation
6. Admin pages go in `src/routes/admin/` with role guards
7. Report bugs at https://github.com/MUKE-coder/grit/issues

## License

MIT
