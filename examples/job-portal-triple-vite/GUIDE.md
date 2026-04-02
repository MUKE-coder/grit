# Building a Job Portal with Grit -- Triple + TanStack Router (Vite)

This guide walks through every step of building a production-ready job portal
using Grit's Triple architecture with Vite and TanStack Router. By the end
you will have a Go API, a Vite-powered public website, and a Vite-powered
admin panel -- all sharing types through a monorepo.

This is the same job portal as the `job-portal-triple-next` example, but
using TanStack Router instead of Next.js. The API is identical. The
differences are all in the frontend apps.

---

## Table of Contents

1. [Why Triple + TanStack Router](#1-why-triple--tanstack-router)
2. [Prerequisites](#2-prerequisites)
3. [Scaffold the Project](#3-scaffold-the-project)
4. [Explore the Generated Structure](#4-explore-the-generated-structure)
5. [Start Development Services](#5-start-development-services)
6. [Generate the Category Resource](#6-generate-the-category-resource)
7. [Generate the Company Resource](#7-generate-the-company-resource)
8. [Generate the Job Resource](#8-generate-the-job-resource)
9. [Generate the Application Resource](#9-generate-the-application-resource)
10. [Run Migrations and Seed Data](#10-run-migrations-and-seed-data)
11. [Configure Authentication](#11-configure-authentication)
12. [Set Up OAuth Providers](#12-set-up-oauth-providers)
13. [Enable Two-Factor Authentication](#13-enable-two-factor-authentication)
14. [Customize the Landing Page](#14-customize-the-landing-page)
15. [Build the Job Listings Page](#15-build-the-job-listings-page)
16. [Build the Application Form](#16-build-the-application-form)
17. [Configure File Uploads](#17-configure-file-uploads)
18. [Set Up Email Notifications](#18-set-up-email-notifications)
19. [Admin Panel Customizations](#19-admin-panel-customizations)
20. [Add Dashboard Analytics](#20-add-dashboard-analytics)
21. [Export Data Tables](#21-export-data-tables)
22. [Testing](#22-testing)
23. [Production Deployment](#23-production-deployment)
24. [Wails Desktop Compatibility](#24-wails-desktop-compatibility)
25. [Next Steps](#25-next-steps)

---

## 1. Why Triple + TanStack Router

TanStack Router gives you fully type-safe routing with the fastest possible
development experience through Vite's HMR. Compared to the Next.js version:

- **No SSR complexity** -- Everything is client-side. Simpler mental model.
- **Faster dev builds** -- Vite's hot module replacement is near-instant.
- **Type-safe params** -- `Route.useParams()` returns typed parameters,
  catching routing bugs at compile time.
- **Simpler deployment** -- Static files can be served by the Go API itself,
  a CDN, or any static file server. No Node.js runtime needed in production.
- **Desktop ready** -- Switch to `createHashHistory()` and the app works
  inside a Wails desktop shell.

The trade-off is no server-side rendering. If your job listings need to
be indexed by Google, either use the Next.js version or add a prerendering
step (see section 23).

## 2. Prerequisites

Install the required tools:

```bash
# Go
go version   # 1.21+

# Node.js and pnpm
node --version   # 18+
pnpm --version   # 8+

# Docker
docker --version
docker compose version

# Grit CLI
go install github.com/MUKE-coder/grit/v3/cmd/grit@latest
grit version
```

## 3. Scaffold the Project

```bash
grit new job-portal --triple --vite
```

The `--vite` flag selects TanStack Router for both frontend apps. This is
the key difference from `grit new job-portal --triple --next`.

What gets generated:
- `apps/api` -- Go backend (identical to the Next.js version)
- `apps/web` -- Vite + TanStack Router + React Query
- `apps/admin` -- Vite + TanStack Router + Grit admin components
- `packages/shared` -- Shared TypeScript types and Zod schemas
- `docker-compose.yml` -- PostgreSQL, Redis, MinIO, Mailhog

## 4. Explore the Generated Structure

```bash
cd job-portal
tree -L 3 --dirsfirst
```

Key differences from the Next.js version:

```
apps/web/
├── src/
│   ├── routes/              # TanStack Router file-based routes
│   │   ├── __root.tsx       # Root layout (navbar, footer)
│   │   ├── index.tsx        # "/" landing page
│   │   └── jobs/
│   │       ├── index.tsx    # "/jobs" listings
│   │       └── $id.tsx      # "/jobs/:id" detail (typed $id param)
│   ├── components/          # UI components
│   ├── hooks/               # React Query hooks
│   ├── lib/                 # API client, utils
│   └── main.tsx             # App entry point + router setup
├── index.html               # SPA shell
├── vite.config.ts           # Dev proxy + build config
└── package.json
```

Notice:
- `index.html` at the root (Vite entry point, not in `public/`)
- `src/routes/` instead of `app/` (TanStack Router convention)
- `$id.tsx` instead of `[id]/page.tsx` (dollar sign for dynamic params)
- `__root.tsx` instead of `layout.tsx` (root layout)
- `vite.config.ts` instead of `next.config.js`

## 5. Start Development Services

```bash
# Start PostgreSQL, Redis, MinIO, and Mailhog
docker compose up -d

# Verify
docker compose ps
```

Start the dev servers:

```bash
# Terminal 1: API (same as Next.js version)
cd apps/api && go run cmd/server/main.go

# Terminal 2: All frontends via pnpm
pnpm install && pnpm dev
```

Vite serves the web app at http://localhost:5173 and the admin at
http://localhost:5174. Note the different ports from Next.js (3000/3001).

The `vite.config.ts` in each frontend app includes a dev proxy:

```ts
export default defineConfig({
  server: {
    port: 5173,   // 5174 for admin
    proxy: {
      '/v3': 'http://localhost:8080',
    },
  },
})
```

This proxies API requests so the frontend can call `/v3/jobs` without
worrying about CORS in development.

## 6. Generate the Category Resource

```bash
grit generate resource Category --fields "name:string:unique,slug:slug:name"
```

This generates the same API code as the Next.js version. On the frontend
side, the difference is in the generated route files:

- **Next.js**: `app/(dashboard)/categories/page.tsx`
- **Vite**: `src/routes/categories/index.tsx`

The component code inside is the same -- DataTable, FormBuilder, React Query
hooks. Only the routing wrapper differs.

## 7. Generate the Company Resource

```bash
grit generate resource Company \
  --fields "name:string,slug:slug:name,description:richtext,logo:string,website:string:optional,email:string"
```

The generated company detail page uses TanStack Router's typed params:

```tsx
// apps/web/src/routes/companies/$slug.tsx
import { createFileRoute } from '@tanstack/react-router'

export const Route = createFileRoute('/companies/$slug')({
  component: CompanyDetailPage,
})

function CompanyDetailPage() {
  const { slug } = Route.useParams()  // slug is typed as string
  // fetch company by slug...
}
```

Compare with Next.js where you would use `params.slug` from the page props
or `useParams()` from `next/navigation`.

## 8. Generate the Job Resource

```bash
grit generate resource Job \
  --fields "title:string,description:richtext,company_id:belongs_to:Company,location:string,salary_min:float,salary_max:float,type:string,experience:string,active:bool,deadline:date"
```

The job detail page demonstrates TanStack Router's `Link` component:

```tsx
// Navigating to a job detail page
import { Link } from '@tanstack/react-router'

<Link to="/jobs/$id" params={{ id: job.id }}>
  {job.title}
</Link>
```

Compare with Next.js:

```tsx
// Next.js equivalent
import Link from 'next/link'

<Link href={`/jobs/${job.id}`}>
  {job.title}
</Link>
```

TanStack Router's version is type-safe -- if you typo the route path or
forget a required param, TypeScript catches it at compile time.

## 9. Generate the Application Resource

```bash
grit generate resource Application \
  --fields "job_id:belongs_to:Job,applicant_name:string,applicant_email:string,resume:string,cover_letter:text:optional,status:string"
```

The application form page at `src/routes/apply/$jobId.tsx` uses
`Route.useParams()` to get the typed job ID:

```tsx
export const Route = createFileRoute('/apply/$jobId')({
  component: ApplyPage,
})

function ApplyPage() {
  const { jobId } = Route.useParams()
  // Load job details, show application form...
}
```

## 10. Run Migrations and Seed Data

```bash
grit migrate
grit seed
```

This is identical to the Next.js version. The API is the same -- only the
frontend framework differs.

## 11. Configure Authentication

Auth configuration is identical to the Next.js version. Set in `.env`:

```env
JWT_SECRET=your-random-32-byte-secret
JWT_ACCESS_TTL=15m
JWT_REFRESH_TTL=168h
```

The frontend auth flow uses the same React hooks but with TanStack Router
for redirects:

```tsx
// Redirect after login
const navigate = useNavigate()
navigate({ to: '/dashboard' })

// Compare with Next.js:
// const router = useRouter()
// router.push('/dashboard')
```

Protected routes use TanStack Router's `beforeLoad` guard:

```tsx
export const Route = createFileRoute('/dashboard')({
  beforeLoad: ({ context }) => {
    if (!context.auth.isAuthenticated) {
      throw redirect({ to: '/login' })
    }
  },
  component: DashboardPage,
})
```

## 12. Set Up OAuth Providers

OAuth setup is identical to the Next.js version -- it is handled entirely
by the Go API. The frontend just redirects to:

```
GET /v3/auth/google/login
GET /v3/auth/github/login
```

After the OAuth flow completes, the API redirects back to the frontend
with tokens. Configure the redirect URL in `.env`:

```env
GOOGLE_CLIENT_ID=your-client-id
GOOGLE_CLIENT_SECRET=your-client-secret
GOOGLE_REDIRECT_URL=http://localhost:8080/v3/auth/google/callback

GITHUB_CLIENT_ID=your-client-id
GITHUB_CLIENT_SECRET=your-client-secret
GITHUB_REDIRECT_URL=http://localhost:8080/v3/auth/github/callback
```

## 13. Enable Two-Factor Authentication

TOTP 2FA is identical to the Next.js version. The API endpoints are the
same. The only frontend difference is using `useNavigate()` instead of
`useRouter()` for redirects after 2FA verification.

```env
TOTP_ISSUER=JobPortal
```

## 14. Customize the Landing Page

Open `apps/web/src/routes/index.tsx`. This is the landing page.

Unlike Next.js, there are no server components. All data fetching uses
React Query:

```tsx
import { createFileRoute } from '@tanstack/react-router'
import { useQuery } from '@tanstack/react-query'

export const Route = createFileRoute('/')({
  component: LandingPage,
})

function LandingPage() {
  const { data: featuredJobs } = useQuery({
    queryKey: ['jobs', 'featured'],
    queryFn: () => api.get('/v3/jobs?active=true&limit=6'),
  })

  const { data: categories } = useQuery({
    queryKey: ['categories'],
    queryFn: () => api.get('/v3/categories'),
  })

  return (
    <div>
      <HeroSection />
      <FeaturedJobs jobs={featuredJobs} />
      <CategoryCards categories={categories} />
      <CTASection />
    </div>
  )
}
```

Key sections to customize:
- **Hero**: Search bar with location and keyword inputs
- **Featured Jobs**: Grid of active jobs with company logos
- **Categories**: Cards showing job categories with counts
- **CTA**: Call-to-action for employers to post jobs

## 15. Build the Job Listings Page

The generated `apps/web/src/routes/jobs/index.tsx` provides a paginated
job list. Filters are passed as search params:

```tsx
import { createFileRoute } from '@tanstack/react-router'

type JobSearchParams = {
  q?: string
  location?: string
  type?: string
  salary_min?: number
  salary_max?: number
  page?: number
}

export const Route = createFileRoute('/jobs/')({
  validateSearch: (search: Record<string, unknown>): JobSearchParams => ({
    q: search.q as string,
    location: search.location as string,
    type: search.type as string,
    salary_min: Number(search.salary_min) || undefined,
    salary_max: Number(search.salary_max) || undefined,
    page: Number(search.page) || 1,
  }),
  component: JobListingsPage,
})

function JobListingsPage() {
  const search = Route.useSearch()  // Fully typed search params
  // Use search.q, search.location, etc. to fetch and filter
}
```

TanStack Router's `validateSearch` gives you type-safe search parameters --
something Next.js does not provide out of the box.

## 16. Build the Application Form

The application form at `src/routes/apply/$jobId.tsx` works the same as the
Next.js version but uses TanStack Router navigation:

```tsx
const navigate = useNavigate()

async function onSubmit(data: ApplicationFormData) {
  await api.post('/v3/applications', data)
  navigate({ to: '/jobs/$id', params: { id: jobId } })
}
```

Form validation uses the same Zod schema from `packages/shared`.

## 17. Configure File Uploads

File upload configuration is identical to the Next.js version. The presigned
URL flow is the same:

1. `POST /v3/upload/presign` -- Get presigned URL
2. `PUT` to presigned URL -- Upload file directly to storage
3. Submit form with the storage key

```env
STORAGE_ENDPOINT=localhost:9000
STORAGE_ACCESS_KEY=minioadmin
STORAGE_SECRET_KEY=minioadmin
STORAGE_BUCKET=job-portal
STORAGE_USE_SSL=false
```

The React hooks for file upload are shared between Next.js and Vite versions
through `packages/shared`.

## 18. Set Up Email Notifications

Email setup is identical to the Next.js version -- it is handled entirely
by the Go API.

Three emails for this project:
1. **Welcome Email** -- On registration
2. **Application Confirmation** -- When someone applies
3. **Status Update** -- When admin changes application status

Development: Mailhog at http://localhost:8025
Production: Resend API

```env
RESEND_API_KEY=re_xxxxxxxxxxxx
RESEND_FROM=JobPortal <noreply@jobs.example.com>
```

## 19. Admin Panel Customizations

The admin panel at `apps/admin` uses TanStack Router for navigation. Each
resource has routes in `src/routes/`:

```
src/routes/
├── __root.tsx           # Sidebar layout
├── index.tsx            # Dashboard
├── jobs/
│   ├── index.tsx        # Job list (DataTable)
│   └── $id.tsx          # Job edit (FormBuilder)
├── companies/
│   ├── index.tsx
│   └── $id.tsx
├── applications/
│   ├── index.tsx
│   └── $id.tsx
└── categories/
    ├── index.tsx
    └── $id.tsx
```

### Application Status Workflow

The application detail page at `src/routes/applications/$id.tsx` shows the
status workflow with a select dropdown:

- **Pending** (yellow) -- Default status
- **Reviewed** (blue) -- Admin has seen it
- **Interview** (purple) -- Scheduled for interview
- **Accepted** (green) -- Offer extended
- **Rejected** (red) -- Not selected

### Admin Navigation

TanStack Router's `Link` component highlights the active route automatically:

```tsx
<Link
  to="/jobs"
  activeProps={{ className: 'bg-accent text-white' }}
>
  Jobs
</Link>
```

## 20. Add Dashboard Analytics

The admin dashboard at `src/routes/index.tsx` shows four stats cards:

1. **Total Jobs** -- Count of all jobs
2. **Active Jobs** -- Count where `active = true`
3. **Total Applications** -- Count of all applications
4. **Companies** -- Count of all companies

Below the stats, a "Recent Applications" table shows the last 10
applications. Data fetching uses React Query (same as the Next.js version).

## 21. Export Data Tables

Export works identically to the Next.js version. The API endpoints are the
same:

```
GET /v3/jobs/export?format=csv&filter[active]=true
GET /v3/applications/export?format=xlsx
```

The DataTable component triggers a download when the export button is clicked.

## 22. Testing

### API Tests

```bash
cd apps/api
go test ./...
```

Identical to the Next.js version.

### Frontend Tests

```bash
# Run all frontend tests
pnpm test

# Run web tests only
pnpm --filter web test

# Run admin tests only
pnpm --filter admin test
```

Tests use Vitest (not Jest) since we are using Vite. The testing APIs are
nearly identical but Vitest is faster because it reuses Vite's transform
pipeline.

### End-to-End Tests

```bash
pnpm e2e
```

Playwright tests work the same way. The only difference is the dev server
URLs (5173/5174 instead of 3000/3001).

## 23. Production Deployment

### Option A: grit deploy (recommended)

The simplest deployment. Since Vite builds to static files, the Go API
serves them directly. Only one container needed (no Node.js runtime).

```bash
cp .env.example .env.production
# Edit .env.production with real values

grit deploy \
  --host deploy@server.com \
  --domain jobs.example.com \
  --env .env.production
```

This command:
1. Builds both Vite apps (`pnpm build`)
2. Copies `dist/` output into the Go API's static file directory
3. Builds the Go binary
4. Deploys a single Docker container that serves API + frontend
5. Configures Caddy for automatic HTTPS

This is simpler than the Next.js version because there is no separate
Node.js container for the frontend.

### Option B: Docker Compose

```bash
docker compose -f docker-compose.prod.yml up -d
```

The production compose file uses a multi-stage Dockerfile:
1. Stage 1: Build Vite apps (Node.js)
2. Stage 2: Build Go binary
3. Stage 3: Copy binary + static files into Alpine image

Final image contains: Go binary + `web-dist/` + `admin-dist/`.
The API serves `web-dist/` at `/` and `admin-dist/` at `/admin`.

### Option C: CDN + VPS

Deploy static files to a CDN and the API to a VPS:

```bash
# API on VPS
grit deploy --host deploy@server.com --domain api.jobs.example.com

# Build web
cd apps/web
VITE_API_URL=https://api.jobs.example.com pnpm build

# Deploy to Cloudflare Pages
npx wrangler pages deploy dist --project-name=job-portal

# Build and deploy admin
cd apps/admin
VITE_ADMIN_API_URL=https://api.jobs.example.com pnpm build
npx wrangler pages deploy dist --project-name=job-portal-admin
```

Other CDN options: Vercel (static mode), Netlify, AWS CloudFront.

### SEO Without SSR

Since Vite produces a client-side SPA, search engines may not index content
rendered by JavaScript. Options:

1. **Prerendering** -- Use `vite-plugin-ssr` or `react-snap` to prerender
   key pages (job listings, company profiles) at build time.
2. **Meta tags** -- The Go API can serve an HTML page with `<meta>` tags
   for social sharing. The SPA takes over after the initial load.
3. **Sitemap** -- Generate a `sitemap.xml` from the API with all job URLs.

### Production Checklist

- [ ] Change `JWT_SECRET` to a random 32+ byte string
- [ ] Set real database credentials
- [ ] Configure Resend API key for email delivery
- [ ] Set up S3/R2 for file storage
- [ ] Configure OAuth redirect URLs for production domain
- [ ] Enable rate limiting (`RATE_LIMIT_MAX=100`)
- [ ] Set `LOG_LEVEL=info`
- [ ] Set up database backups (pg_dump cron)
- [ ] Add prerendering for SEO if needed
- [ ] Configure CDN caching for static assets

## 24. Wails Desktop Compatibility

One major advantage of the Vite version: it can be embedded in a Wails
desktop application with minimal changes.

To enable desktop compatibility, switch to hash-based routing in
`apps/web/src/main.tsx`:

```tsx
import { createRouter, createHashHistory } from '@tanstack/react-router'
import { routeTree } from './routeTree.gen'

const router = createRouter({
  routeTree,
  history: createHashHistory(),
})
```

Hash routing uses `/#/jobs/123` instead of `/jobs/123`, which works with
Wails' embedded webview that serves files from the filesystem.

You can also scaffold a standalone Wails desktop app with:

```bash
grit new-desktop my-desktop-app
```

See the `grit new-desktop` documentation for details.

## 25. Next Steps

With the core job portal running, consider these enhancements:

- **Job alerts**: Let users subscribe to email alerts for new jobs matching
  their criteria. Use Redis pub/sub and a background worker.
- **Applicant tracking**: Add interview scheduling, notes, and rating to
  the Application resource.
- **Employer accounts**: Add a role between USER and ADMIN that lets
  companies manage their own job listings.
- **AI features**: Use the AI Gateway to generate job descriptions, parse
  resumes, and match candidates to jobs.
- **Analytics**: Track page views, application conversion rates, and
  search patterns.
- **Desktop app**: Wrap the admin panel in Wails for an installable
  desktop application.

---

Built with [Grit Framework](https://github.com/MUKE-coder/grit).
