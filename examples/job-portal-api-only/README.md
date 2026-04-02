# Job Portal -- API Only (Headless)

A headless job portal API built with Grit using the **API Only** architecture:
- `apps/api` -- Go API (Gin + GORM + PostgreSQL)
- No frontend -- bring your own client (mobile app, SPA, third-party integration)

## Architecture

This is the **leanest** Grit configuration. A pure Go API with no frontend apps,
no Node.js, no pnpm, no Turborepo. The API speaks JSON and serves interactive
documentation at `/docs`. Everything else -- mobile apps, SPAs, CLI tools -- is
built separately and talks to this API over HTTP.

Best for: mobile backends, microservices, headless CMS, third-party API providers,
teams where frontend and backend are separate repos.

## Quick Start

### Prerequisites
- Go 1.21+
- Docker (for PostgreSQL, Redis, MinIO, Mailhog)

That's it. No Node.js, no pnpm, no frontend toolchain.

### Setup

1. Install Grit CLI:
   ```bash
   go install github.com/MUKE-coder/grit/v3/cmd/grit@latest
   ```

2. Create the project:
   ```bash
   grit new job-portal --api
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

   With `--api`, only Go files are generated (model, service, handler, routes).
   No Zod schemas, no TypeScript types, no React hooks, no admin pages.

5. Run migrations and seed:
   ```bash
   grit migrate
   grit seed
   ```

6. Start the API:
   ```bash
   cd apps/api && go run cmd/server/main.go
   ```

7. Open:
   - API Docs: http://localhost:8080/docs
   - GORM Studio: http://localhost:8080/studio
   - Mailhog: http://localhost:8025

## Project Structure

```
job-portal/
├── apps/
│   └── api/                    # Go backend (the only app)
│       ├── cmd/server/main.go  # Entry point
│       └── internal/
│           ├── config/         # Environment + app config
│           ├── database/       # GORM setup, migrations
│           ├── models/         # Job, Company, Application, Category, User
│           ├── handlers/       # REST handlers
│           ├── services/       # Business logic
│           ├── middleware/      # Auth, CORS, cache, rate limiting
│           └── routes/         # Route registration
├── .env                        # Environment variables
├── .env.example                # Template for other developers
├── docker-compose.yml          # Dev services (Postgres, Redis, MinIO, Mailhog)
└── docker-compose.prod.yml     # Production deployment
```

No `packages/shared`, no `apps/web`, no `apps/admin`, no `pnpm-workspace.yaml`,
no `turbo.json`, no `node_modules`. Just Go.

## Features Demonstrated

### Authentication
- Email/password registration and login
- OAuth2 social login (Google, GitHub)
- JWT access tokens (15 min) + refresh tokens (7 days)
- TOTP two-factor authentication with backup codes
- Role-based access: ADMIN (full access), EDITOR (manage jobs), USER (apply)

### REST API
- Full CRUD for Jobs, Companies, Applications, Categories
- Server-side pagination, sorting, filtering via query params
- Presigned URL file uploads for resumes and company logos
- Application status workflow: Pending -> Reviewed -> Interview -> Accepted/Rejected
- Interactive API documentation at `/docs` (Scalar UI)

### Emails
- Welcome email on registration
- Application confirmation email
- Application status update notifications

### Dashboard Stats
- `GET /api/stats` -- Total jobs, companies, applications, users
- `GET /api/stats/applications` -- Applications by status breakdown
- `GET /api/stats/jobs` -- Jobs by category breakdown

### Deployment
See the Deployment section below.

## API Reference

### Authentication

| Method | Endpoint | Description | Auth |
|--------|----------|-------------|------|
| `POST` | `/api/auth/register` | Register new user | -- |
| `POST` | `/api/auth/login` | Login (returns JWT) | -- |
| `POST` | `/api/auth/refresh` | Refresh access token | Refresh token |
| `GET` | `/api/auth/me` | Get current user | Bearer |
| `POST` | `/api/auth/totp/setup` | Enable 2FA | Bearer |
| `POST` | `/api/auth/totp/verify` | Verify TOTP code | Bearer |
| `GET` | `/api/auth/google` | Start Google OAuth | -- |
| `GET` | `/api/auth/google/callback` | Google OAuth callback | -- |
| `GET` | `/api/auth/github` | Start GitHub OAuth | -- |
| `GET` | `/api/auth/github/callback` | GitHub OAuth callback | -- |

### Jobs

| Method | Endpoint | Description | Auth |
|--------|----------|-------------|------|
| `GET` | `/api/jobs` | List jobs (paginated) | -- |
| `POST` | `/api/jobs` | Create job | ADMIN, EDITOR |
| `GET` | `/api/jobs/:id` | Get job by ID | -- |
| `PUT` | `/api/jobs/:id` | Update job | ADMIN, EDITOR |
| `DELETE` | `/api/jobs/:id` | Delete job | ADMIN |

### Companies

| Method | Endpoint | Description | Auth |
|--------|----------|-------------|------|
| `GET` | `/api/companies` | List companies | -- |
| `POST` | `/api/companies` | Create company | ADMIN, EDITOR |
| `GET` | `/api/companies/:id` | Get company by ID | -- |
| `PUT` | `/api/companies/:id` | Update company | ADMIN, EDITOR |

### Applications

| Method | Endpoint | Description | Auth |
|--------|----------|-------------|------|
| `POST` | `/api/applications` | Submit application | Bearer |
| `GET` | `/api/applications` | List applications | ADMIN |
| `GET` | `/api/applications/:id` | Get application | ADMIN, Owner |
| `PUT` | `/api/applications/:id` | Update status | ADMIN |

### Categories

| Method | Endpoint | Description | Auth |
|--------|----------|-------------|------|
| `GET` | `/api/categories` | List categories | -- |
| `POST` | `/api/categories` | Create category | ADMIN |

### Uploads

| Method | Endpoint | Description | Auth |
|--------|----------|-------------|------|
| `POST` | `/api/uploads/presign` | Get presigned upload URL | Bearer |

### Stats

| Method | Endpoint | Description | Auth |
|--------|----------|-------------|------|
| `GET` | `/api/stats` | Dashboard overview stats | ADMIN |
| `GET` | `/api/stats/applications` | Applications by status | ADMIN |
| `GET` | `/api/stats/jobs` | Jobs by category | ADMIN |

### System

| Method | Endpoint | Description | Auth |
|--------|----------|-------------|------|
| `GET` | `/v3/health` | Health check | -- |
| `GET` | `/docs` | Interactive API documentation | -- |
| `GET` | `/studio` | GORM Studio (database browser) | -- |

## Testing via curl

### Register a user

```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Alice Johnson",
    "email": "alice@example.com",
    "password": "securePass123!"
  }'
```

Response:
```json
{
  "data": {
    "id": "uuid-here",
    "name": "Alice Johnson",
    "email": "alice@example.com",
    "role": "USER"
  },
  "message": "Registration successful"
}
```

### Login

```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "alice@example.com",
    "password": "securePass123!"
  }'
```

Response:
```json
{
  "data": {
    "access_token": "eyJhbG...",
    "refresh_token": "eyJhbG...",
    "user": { "id": "uuid-here", "name": "Alice Johnson", "role": "USER" }
  },
  "message": "Login successful"
}
```

### Create a job (requires ADMIN or EDITOR)

```bash
curl -X POST http://localhost:8080/api/jobs \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <access_token>" \
  -d '{
    "title": "Senior Go Developer",
    "description": "<p>We are looking for an experienced Go developer...</p>",
    "company_id": "company-uuid",
    "location": "Remote",
    "salary_min": 120000,
    "salary_max": 180000,
    "type": "Full-time",
    "experience": "Senior",
    "active": true,
    "deadline": "2026-06-01"
  }'
```

### List jobs with pagination and filters

```bash
# Page 1, 10 per page, sorted by created_at desc
curl "http://localhost:8080/api/jobs?page=1&page_size=10&sort=-created_at"

# Filter by location and type
curl "http://localhost:8080/api/jobs?location=Remote&type=Full-time"

# Search by title
curl "http://localhost:8080/api/jobs?search=golang"
```

### Submit an application with resume

```bash
# Step 1: Get a presigned upload URL for the resume
curl -X POST http://localhost:8080/api/uploads/presign \
  -H "Authorization: Bearer <access_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "filename": "alice-resume.pdf",
    "content_type": "application/pdf"
  }'
# Response: { "data": { "upload_url": "https://...", "file_url": "https://..." } }

# Step 2: Upload the file to the presigned URL
curl -X PUT "<upload_url>" \
  -H "Content-Type: application/pdf" \
  --data-binary @alice-resume.pdf

# Step 3: Submit the application with the file URL
curl -X POST http://localhost:8080/api/applications \
  -H "Authorization: Bearer <access_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "job_id": "job-uuid",
    "applicant_name": "Alice Johnson",
    "applicant_email": "alice@example.com",
    "resume": "<file_url from step 1>",
    "cover_letter": "I am excited to apply for this position..."
  }'
```

### Update application status (ADMIN)

```bash
curl -X PUT http://localhost:8080/api/applications/<id> \
  -H "Authorization: Bearer <admin_access_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "status": "Interview"
  }'
```

### Enable TOTP 2FA

```bash
# Step 1: Get TOTP secret and QR code
curl -X POST http://localhost:8080/api/auth/totp/setup \
  -H "Authorization: Bearer <access_token>"
# Response: { "data": { "secret": "JBSWY3DPEHPK3PXP", "qr_url": "otpauth://..." } }

# Step 2: Verify with code from authenticator app
curl -X POST http://localhost:8080/api/auth/totp/verify \
  -H "Authorization: Bearer <access_token>" \
  -H "Content-Type: application/json" \
  -d '{ "code": "123456" }'
```

### Refresh access token

```bash
curl -X POST http://localhost:8080/api/auth/refresh \
  -H "Content-Type: application/json" \
  -d '{ "refresh_token": "eyJhbG..." }'
```

## When to Choose API Only

Choose `--api` when:

- **Mobile backend** -- Your frontend is a native iOS/Android app or React Native
- **Microservice** -- This API is one service in a larger system
- **Headless CMS** -- Content is consumed by multiple clients (web, mobile, IoT)
- **Third-party API** -- Other companies or developers will integrate with your API
- **Separate frontend team** -- Frontend lives in a different repo with its own stack
- **CLI/Desktop client** -- The client is not a web browser

Do NOT choose `--api` when:
- You want an admin panel to manage content -- use `--single`, `--double`, or `--triple`
- You want a public website -- use `--single` or `--triple`
- You want everything in one repo -- use `--double` or `--triple`

## Environment Variables

See `.env.example` for the full list.

Key variables:
- `DATABASE_URL` -- PostgreSQL connection string
- `JWT_SECRET` -- Must be changed in production
- `STORAGE_ENDPOINT` -- MinIO (dev) or S3/R2 (prod)
- `RESEND_API_KEY` -- For sending emails
- `GOOGLE_CLIENT_ID` -- For OAuth (optional)
- `GITHUB_CLIENT_ID` -- For OAuth (optional)

## Deployment

### Option 1: grit deploy (VPS)
```bash
grit deploy --host deploy@server.com --domain api.jobs.example.com
```

This builds the Go binary, copies it to your server, and sets up a systemd
service behind Caddy. One command, done.

### Option 2: Docker Compose
```bash
docker compose -f docker-compose.prod.yml up -d
```

This starts PostgreSQL, Redis, MinIO, and the Go API in production mode.
Put Caddy or Nginx in front for TLS.

### Option 3: Build and ship the binary
```bash
cd apps/api
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server cmd/server/main.go
scp server deploy@server.com:/opt/job-portal/
```

The Go binary is a single file with zero dependencies. Copy it anywhere.

## LLM Instructions

If you are an AI assistant working on this project:
1. Read `.claude/skills/grit/SKILL.md` for Grit conventions
2. Use `grit generate resource` for new resources
3. Never delete `// grit:*` marker comments
4. Follow the API response format: `{ data, message }` / `{ data, meta }` / `{ error: { code, message } }`
5. Report bugs at https://github.com/MUKE-coder/grit/issues

## License

MIT
