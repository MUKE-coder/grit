# Building a Job Portal with Grit -- Triple + Next.js

This guide walks through every step of building a production-ready job portal
using Grit's Triple architecture. By the end you will have a Go API, a
Next.js public website, and a Next.js admin panel -- all sharing types and
schemas through a monorepo.

---

## Table of Contents

1. [Why Triple + Next.js](#1-why-triple--nextjs)
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
24. [Next Steps](#24-next-steps)

---

## 1. Why Triple + Next.js

The Triple architecture gives you three separate apps that share code through
a monorepo. Compared to Double (API + single frontend), Triple adds a
dedicated admin panel. This matters for job portals because:

- **Public site** needs SEO (Next.js SSR), fast page loads, and a clean UI.
- **Admin panel** needs data tables, form builders, and role-based access --
  concerns that would clutter the public codebase.
- **API** stays clean and serves both frontends through the same endpoints.

Next.js is chosen over Vite/TanStack Router here because the public job
listings benefit from server-side rendering for SEO and social sharing.

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
grit new job-portal --triple --next
```

This command creates a pnpm monorepo with:
- `apps/api` -- Go backend with Gin, GORM, and PostgreSQL
- `apps/web` -- Next.js 14 with App Router
- `apps/admin` -- Next.js 14 with Grit's admin framework
- `packages/shared` -- Shared TypeScript types and Zod schemas
- `docker-compose.yml` -- PostgreSQL, Redis, MinIO, Mailhog

The `--triple` flag tells Grit to scaffold three apps. The `--next` flag
selects Next.js for both frontend apps (the alternative is `--vite` for
TanStack Router).

## 4. Explore the Generated Structure

```bash
cd job-portal
tree -L 3 --dirsfirst
```

Key files to note:
- `apps/api/cmd/server/main.go` -- API entry point
- `apps/api/internal/models/user.go` -- User model (always generated)
- `apps/web/app/page.tsx` -- Landing page
- `apps/admin/app/(dashboard)/layout.tsx` -- Admin sidebar layout
- `packages/shared/types/index.ts` -- Shared TypeScript interfaces
- `.env` -- Pre-filled with dev defaults
- `pnpm-workspace.yaml` -- Monorepo configuration

## 5. Start Development Services

```bash
# Start PostgreSQL, Redis, MinIO, and Mailhog
docker compose up -d

# Verify all services are running
docker compose ps
```

You should see four containers running. MinIO provides S3-compatible storage
for file uploads, and Mailhog catches all outgoing emails at
http://localhost:8025.

## 6. Generate the Category Resource

Categories are the simplest resource -- just a name and an auto-generated slug.
Generate it first because other resources will reference categories later.

```bash
grit generate resource Category --fields "name:string:unique,slug:slug:name"
```

This generates:
- **API**: model, handler, service, routes, migration
- **Admin**: resource definition with DataTable and FormBuilder
- **Shared**: TypeScript type and Zod schema

The `slug:slug:name` field type tells Grit to auto-generate a URL-safe slug
from the `name` field. The `:unique` modifier adds a database unique
constraint.

## 7. Generate the Company Resource

```bash
grit generate resource Company \
  --fields "name:string,slug:slug:name,description:richtext,logo:string,website:string:optional,email:string"
```

Important field choices:
- `description:richtext` -- Generates a Tiptap rich text editor in forms
- `logo:string` -- Stores the S3 key; the admin form renders a file upload
- `website:string:optional` -- The `:optional` modifier makes the field
  nullable in the database and optional in Zod schemas

## 8. Generate the Job Resource

```bash
grit generate resource Job \
  --fields "title:string,description:richtext,company_id:belongs_to:Company,location:string,salary_min:float,salary_max:float,type:string,experience:string,active:bool,deadline:date"
```

Key points:
- `company_id:belongs_to:Company` -- Creates a foreign key to the companies
  table. GORM will add the relationship, and the API handler automatically
  preloads the Company when fetching jobs.
- `salary_min:float,salary_max:float` -- Two floats for salary range filtering
- `active:bool` -- Toggle jobs on/off without deleting them
- `deadline:date` -- Date-only field (no time component)

## 9. Generate the Application Resource

```bash
grit generate resource Application \
  --fields "job_id:belongs_to:Job,applicant_name:string,applicant_email:string,resume:string,cover_letter:text:optional,status:string"
```

The `resume:string` field stores the S3 key for the uploaded PDF. The
`status:string` field will hold values like "pending", "reviewed",
"interview", "accepted", or "rejected" -- we will add validation for this
in the service layer.

## 10. Run Migrations and Seed Data

```bash
# Apply all migrations
grit migrate

# Seed default data (admin user, sample categories)
grit seed
```

Grit runs GORM AutoMigrate for each model, then executes the seed file at
`apps/api/internal/database/seed.go`. The default seed creates an admin user
with email `admin@example.com` and password `password`.

## 11. Configure Authentication

Grit scaffolds a complete auth system. Open `.env` and set:

```env
JWT_SECRET=your-random-32-byte-secret
JWT_ACCESS_TTL=15m
JWT_REFRESH_TTL=168h
```

The auth flow works like this:
1. `POST /v3/auth/register` -- Create account, receive tokens
2. `POST /v3/auth/login` -- Authenticate, receive tokens
3. `POST /v3/auth/refresh` -- Exchange refresh token for new access token
4. Access token goes in the `Authorization: Bearer <token>` header

The admin panel handles this automatically through its built-in auth provider.

## 12. Set Up OAuth Providers

### Google

1. Go to https://console.cloud.google.com/apis/credentials
2. Create an OAuth 2.0 Client ID
3. Set the redirect URI to `http://localhost:8080/v3/auth/google/callback`
4. Add credentials to `.env`:

```env
GOOGLE_CLIENT_ID=your-client-id
GOOGLE_CLIENT_SECRET=your-client-secret
```

### GitHub

1. Go to https://github.com/settings/developers
2. Create a new OAuth App
3. Set the callback URL to `http://localhost:8080/v3/auth/github/callback`
4. Add credentials to `.env`:

```env
GITHUB_CLIENT_ID=your-client-id
GITHUB_CLIENT_SECRET=your-client-secret
```

Both providers follow the same flow: redirect to provider, receive callback
with code, exchange code for user info, create or link account, return JWT.

## 13. Enable Two-Factor Authentication

TOTP 2FA is built in. The flow:

1. User calls `POST /v3/auth/2fa/enable` -- receives a QR code URI
2. User scans QR code with an authenticator app (Google Authenticator, Authy)
3. User calls `POST /v3/auth/2fa/verify` with the 6-digit code to confirm
4. On subsequent logins, after email/password, user must provide the TOTP code
5. Backup codes are generated during setup for account recovery

Set the issuer name in `.env`:

```env
TOTP_ISSUER=JobPortal
```

## 14. Customize the Landing Page

Open `apps/web/app/page.tsx`. The generated landing page includes a hero
section. Customize it for the job portal:

The landing page uses Next.js server components by default, so data fetching
happens on the server. The `getFeaturedJobs` function calls the API internally
(server-to-server, no CORS needed).

Key sections to customize:
- **Hero**: Search bar with location and keyword inputs
- **Featured Jobs**: Grid of active jobs with company logos
- **Categories**: Cards showing job categories with counts
- **CTA**: Call-to-action for employers to post jobs

## 15. Build the Job Listings Page

The generated `apps/web/app/jobs/page.tsx` provides a paginated job list.
Enhance it with filters:

- **Search**: Full-text search on title and description
- **Location**: Dropdown or text filter
- **Salary Range**: Min/max slider
- **Job Type**: Full-time, Part-time, Contract, Remote
- **Experience**: Entry, Mid, Senior, Lead

These filters map to query parameters: `/jobs?q=react&location=remote&type=full-time`.
The API handler already supports these through Grit's built-in query parser.

Job detail pages live at `apps/web/app/jobs/[id]/page.tsx` and use
`generateMetadata` for SEO -- the job title and company name appear in
social sharing cards.

## 16. Build the Application Form

The application form at `apps/web/app/jobs/[id]/apply/page.tsx` collects:
- Applicant name and email
- Resume (PDF upload)
- Cover letter (optional rich text)

The form uses React Hook Form with the shared Zod schema from
`packages/shared/schemas/application.ts`. File upload uses a presigned URL
flow (see next section).

## 17. Configure File Uploads

Grit uses presigned URLs for file uploads. The flow:

1. Client requests a presigned URL: `POST /v3/upload/presign`
2. API generates a presigned PUT URL from MinIO/S3
3. Client uploads the file directly to storage
4. Client sends the storage key to the API with the form data

Configure storage in `.env`:

```env
STORAGE_ENDPOINT=localhost:9000
STORAGE_ACCESS_KEY=minioadmin
STORAGE_SECRET_KEY=minioadmin
STORAGE_BUCKET=job-portal
STORAGE_USE_SSL=false
```

For production, swap MinIO for AWS S3 or Cloudflare R2 by changing the
endpoint and credentials. The upload code stays the same.

Two upload contexts in this project:
- **Company logos**: Image files, max 2 MB, stored under `logos/`
- **Resumes**: PDF files, max 10 MB, stored under `resumes/`

## 18. Set Up Email Notifications

Grit generates email templates and a sending service. For development,
Mailhog catches all emails at http://localhost:8025.

Three emails for this project:

1. **Welcome Email** -- Sent on registration. Template at
   `apps/api/internal/emails/welcome.go`.

2. **Application Confirmation** -- Sent when someone applies. Includes the
   job title, company name, and a link to track status.

3. **Status Update** -- Sent when an admin changes application status.
   Includes the new status and any notes.

For production, configure Resend:

```env
RESEND_API_KEY=re_xxxxxxxxxxxx
RESEND_FROM=JobPortal <noreply@jobs.example.com>
```

## 19. Admin Panel Customizations

The admin panel at `apps/admin` provides CRUD for all resources. Each
resource definition in `resources/` configures the DataTable columns and
FormBuilder fields.

### Application Status Workflow

Edit `apps/admin/resources/applications.ts` to add a status column with
color-coded badges:

- **Pending** (yellow) -- Default status
- **Reviewed** (blue) -- Admin has seen the application
- **Interview** (purple) -- Scheduled for interview
- **Accepted** (green) -- Offer extended
- **Rejected** (red) -- Not selected

The status field renders as a select dropdown in the form and a colored badge
in the table.

### Company Logo Preview

Edit `apps/admin/resources/companies.ts` to render the logo as an image
preview in the table and a file upload in the form.

## 20. Add Dashboard Analytics

The admin dashboard at `apps/admin/app/(dashboard)/page.tsx` shows stats
cards. Add these four:

1. **Total Jobs** -- Count of all jobs
2. **Active Jobs** -- Count where `active = true`
3. **Total Applications** -- Count of all applications
4. **Companies** -- Count of all companies

The API provides a dashboard endpoint at `GET /v3/dashboard/stats` that
returns these counts. The admin panel fetches them with React Query.

Below the stats cards, add a "Recent Applications" DataTable showing the
last 10 applications with columns: applicant name, job title, company,
status, and date.

## 21. Export Data Tables

Every DataTable in the admin panel supports export. The export button
triggers a server-side CSV/Excel generation:

```
GET /v3/jobs/export?format=csv&filter[active]=true
GET /v3/applications/export?format=xlsx
```

The API uses the `excelize` library for Excel export. Grit generates the
export handler for each resource automatically.

## 22. Testing

### API Tests

```bash
cd apps/api
go test ./...
```

Grit generates test files for each handler. The tests use a test database
(SQLite in-memory by default) and cover CRUD operations, authentication,
and authorization.

### Frontend Tests

```bash
# Run all frontend tests
pnpm test

# Run web tests only
pnpm --filter web test

# Run admin tests only
pnpm --filter admin test
```

### End-to-End Tests

```bash
pnpm e2e
```

Grit generates Playwright tests for critical flows: registration, login,
job listing, and application submission.

## 23. Production Deployment

### Option A: grit deploy (recommended)

The fastest path to production. Grit SSHs into your server, builds Docker
images, and starts containers with Caddy for automatic HTTPS.

```bash
# Prepare the .env for production
cp .env.example .env.production
# Edit .env.production with real values

# Deploy
grit deploy \
  --host deploy@server.com \
  --domain jobs.example.com \
  --env .env.production
```

This command:
1. Copies the project to the server
2. Builds Docker images for API, web, and admin
3. Starts all services with `docker-compose.prod.yml`
4. Configures Caddy with automatic TLS

### Option B: Docker Compose manually

```bash
scp -r . deploy@server.com:~/job-portal/
ssh deploy@server.com
cd job-portal
cp .env.example .env
# Edit .env with production values
docker compose -f docker-compose.prod.yml up -d
```

### Option C: Vercel + VPS split

Deploy the Next.js apps to Vercel for zero-config frontend hosting, and
the Go API to a VPS:

```bash
# API on VPS
grit deploy --host deploy@server.com --domain api.jobs.example.com

# Web on Vercel
cd apps/web
vercel --env NEXT_PUBLIC_API_URL=https://api.jobs.example.com

# Admin on Vercel
cd apps/admin
vercel --env NEXT_PUBLIC_ADMIN_API_URL=https://api.jobs.example.com
```

### Production Checklist

- [ ] Change `JWT_SECRET` to a random 32+ byte string
- [ ] Set real database credentials (not the dev defaults)
- [ ] Configure Resend API key for email delivery
- [ ] Set up S3/R2 for file storage (or keep MinIO behind a CDN)
- [ ] Configure OAuth redirect URLs for production domain
- [ ] Enable rate limiting (`RATE_LIMIT_MAX=100`)
- [ ] Set `LOG_LEVEL=info` (not debug)
- [ ] Set up database backups (pg_dump cron)
- [ ] Configure monitoring (Prometheus + Grafana or similar)

## 24. Next Steps

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

---

Built with [Grit Framework](https://github.com/MUKE-coder/grit).
