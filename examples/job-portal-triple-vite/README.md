# Job Portal -- Triple Architecture + TanStack Router (Vite)

A full-featured job portal built with Grit using the **Triple** architecture:
- `apps/api` -- Go API (Gin + GORM + PostgreSQL)
- `apps/web` -- Vite + TanStack Router public website (job listings, company profiles, apply)
- `apps/admin` -- Vite + TanStack Router admin panel (manage jobs, companies, applications)

## Architecture

Three apps sharing types and schemas via `packages/shared`, using Vite for
blazing-fast HMR and TanStack Router for type-safe client-side routing.

Best for: SaaS platforms, internal tools, and projects where SEO is handled
by meta tags / prerendering rather than SSR.

## Key Differences from the Next.js Version

If you are choosing between this and the `job-portal-triple-next` example,
here are the practical differences:

| Aspect | Next.js | TanStack Router (Vite) |
|--------|---------|----------------------|
| Routing | `app/` directory (file-based) | `src/routes/` (file-based, type-safe) |
| Rendering | SSR + RSC by default | Pure SPA (client-side only) |
| Dev speed | Slower cold start | Fast HMR via Vite |
| SEO | Built-in SSR/SSG | Requires prerendering or meta-only |
| Navigation | `<Link>` from `next/link`, `useRouter()` | `<Link>` from `@tanstack/react-router`, `useNavigate()` |
| Params | `params` prop or `useParams()` | `Route.useParams()` (fully typed) |
| Build output | `.next/` directory | `dist/` directory (static files) |
| Server components | Yes | No -- everything is client-side |
| API proxy | `next.config.js` rewrites | `vite.config.ts` proxy |
| Deployment | Needs Node.js runtime | Static files -- serve from Go API or CDN |
| Wails compat | Not compatible | Use `createHashHistory()` for desktop apps |

### When to Choose Vite + TanStack Router

- Your public pages do not need server-side rendering
- You want the fastest possible development experience
- You plan to embed the frontend in a Wails desktop app later
- You prefer fully type-safe routing with zero runtime cost
- You want simpler deployment (static files, no Node.js server)

### When to Choose Next.js

- Job listings need SEO (Google indexing, social sharing cards)
- You need server components for heavy data fetching
- You want ISR (Incremental Static Regeneration) for cached pages

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
   grit new job-portal --triple --vite
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
   - Web: http://localhost:5173
   - Admin: http://localhost:5174
   - API Docs: http://localhost:8080/docs
   - GORM Studio: http://localhost:8080/studio
   - Mailhog: http://localhost:8025

## Project Structure

```
job-portal/
├── apps/
│   ├── api/                        # Go backend (identical to Next.js version)
│   │   ├── cmd/server/main.go
│   │   └── internal/
│   │       ├── models/
│   │       ├── handlers/
│   │       ├── services/
│   │       ├── middleware/
│   │       └── routes/
│   ├── web/                        # Vite + TanStack Router public site
│   │   ├── src/
│   │   │   ├── routes/
│   │   │   │   ├── __root.tsx      # Root layout
│   │   │   │   ├── index.tsx       # Landing page (hero + featured jobs)
│   │   │   │   ├── jobs/
│   │   │   │   │   ├── index.tsx   # Job listings
│   │   │   │   │   └── $id.tsx     # Job detail (typed $id param)
│   │   │   │   ├── companies/
│   │   │   │   │   └── $slug.tsx   # Company profile
│   │   │   │   └── apply/
│   │   │   │       └── $jobId.tsx  # Application form
│   │   │   ├── components/         # Shared UI components
│   │   │   ├── hooks/              # React Query hooks
│   │   │   └── lib/                # API client, utils
│   │   ├── vite.config.ts          # Dev proxy + build config
│   │   └── index.html              # SPA entry point
│   └── admin/                      # Vite + TanStack Router admin panel
│       └── src/
│           ├── routes/
│           │   ├── __root.tsx      # Admin layout (sidebar + header)
│           │   ├── index.tsx       # Dashboard (stats + recent apps)
│           │   ├── jobs/           # Job management
│           │   ├── companies/      # Company management
│           │   ├── applications/   # Application review
│           │   └── categories/     # Category management
│           └── resources/
│               ├── jobs.ts
│               ├── companies.ts
│               ├── applications.ts
│               └── categories.ts
├── packages/shared/                # Shared types + schemas
│   ├── schemas/
│   └── types/
├── .env
├── .env.example
├── docker-compose.yml
└── docker-compose.prod.yml
```

### Routing Differences

TanStack Router uses file-based routing in `src/routes/`:

```
src/routes/
├── __root.tsx          # Root layout (wraps all routes)
├── index.tsx           # "/" -- landing page
├── jobs/
│   ├── index.tsx       # "/jobs" -- listings
│   └── $id.tsx         # "/jobs/:id" -- detail (typed param)
├── companies/
│   └── $slug.tsx       # "/companies/:slug"
└── apply/
    └── $jobId.tsx      # "/apply/:jobId"
```

Parameters use the `$` prefix and are fully typed:

```tsx
// In $id.tsx
import { Route } from './jobs/$id'

function JobDetailPage() {
  const { id } = Route.useParams()  // id is typed as string
  // ...
}
```

Navigation uses `<Link>` from TanStack Router:

```tsx
import { Link, useNavigate } from '@tanstack/react-router'

// Declarative
<Link to="/jobs/$id" params={{ id: job.id }}>View Job</Link>

// Imperative
const navigate = useNavigate()
navigate({ to: '/jobs/$id', params: { id: job.id } })
```

### Vite Dev Proxy

API requests are proxied through Vite in development. In `vite.config.ts`:

```ts
export default defineConfig({
  server: {
    proxy: {
      '/v3': 'http://localhost:8080',
    },
  },
})
```

This means frontend code calls `/v3/jobs` and Vite forwards it to the Go API.
In production the Go API serves the static files directly, so no proxy is needed.

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
- Client-side rendering with React Query for data fetching

### Emails
- Welcome email on registration
- Application confirmation email
- Application status update notifications

### Wails Desktop Compatibility

This Vite-based frontend can be embedded in a Wails desktop app. Switch to
hash-based routing in `src/main.tsx`:

```tsx
import { createHashHistory } from '@tanstack/react-router'

const router = createRouter({
  routeTree,
  history: createHashHistory(),
})
```

## Environment Variables

See `.env.example` for the full list. The variables are identical to the
Next.js version, except the frontend variables use `VITE_` prefix instead
of `NEXT_PUBLIC_`:

```env
VITE_API_URL=http://localhost:8080
VITE_APP_NAME=Job Portal
VITE_GOOGLE_CLIENT_ID=
```

## Deployment

### Option 1: grit deploy (VPS)

Since Vite builds to static files, the Go API can serve them directly.
No separate Node.js container is needed.

```bash
grit deploy --host deploy@server.com --domain jobs.example.com
```

### Option 2: Docker Compose

```bash
docker compose -f docker-compose.prod.yml up -d
```

The production compose file builds the Vite apps into static files and
copies them into the Go API container. The API serves them at `/` (web)
and `/admin` (admin panel). This means only one container serves both
the API and the frontend -- simpler than the Next.js version.

### Option 3: CDN + VPS

```bash
# Deploy API
grit deploy --host deploy@server.com --domain api.jobs.example.com

# Build and deploy web to CDN (e.g., Cloudflare Pages)
cd apps/web
pnpm build
npx wrangler pages deploy dist --project-name=job-portal

# Build and deploy admin to CDN
cd apps/admin
pnpm build
npx wrangler pages deploy dist --project-name=job-portal-admin
```

## LLM Instructions

If you are an AI assistant working on this project:
1. Read `.claude/skills/grit/SKILL.md` for Grit conventions
2. Use `grit generate resource` for new resources
3. Never delete `// grit:*` marker comments
4. Follow the API response format: `{ data, message }` / `{ data, meta }` / `{ error: { code, message } }`
5. Use React Query for data fetching, Zod for validation
6. Routes go in `src/routes/`, use `$param` naming for dynamic segments
7. Use `Route.useParams()` for typed params, `useNavigate()` for navigation
8. Report bugs at https://github.com/MUKE-coder/grit/issues

## License

MIT
