# Job Portal -- API Only Build Guide

A step-by-step guide to building a headless job portal API with Grit.
No frontend, no Node.js -- just Go, PostgreSQL, and a clean REST API.

## Step 1: Scaffold the Project

```bash
grit new job-portal --api
cd job-portal
```

The `--api` flag tells Grit to skip all frontend scaffolding. You get:

```
job-portal/
├── apps/api/                # Go backend -- the only app
│   ├── cmd/server/main.go
│   └── internal/
│       ├── config/
│       ├── database/
│       ├── models/          # User, Upload (built-in)
│       ├── handlers/        # Auth, upload handlers
│       ├── services/        # Auth, storage, email services
│       ├── middleware/       # JWT auth, CORS, rate limit, cache
│       └── routes/          # Route registration
├── .env
├── docker-compose.yml
└── go.mod
```

No `packages/`, no `pnpm-workspace.yaml`, no `turbo.json`. Clean.

## Step 2: Start Infrastructure

```bash
docker compose up -d
```

This starts PostgreSQL, Redis, MinIO (S3-compatible storage), and Mailhog
(email testing). The Go API connects to these on startup.

## Step 3: Run the API

```bash
cd apps/api && go run cmd/server/main.go
```

Open http://localhost:8080/docs to see the interactive API documentation.
This is your primary way to explore and test the API.

Open http://localhost:8080/studio to browse the database with GORM Studio.

## Step 4: Generate Resources

Generate the four resources that make up the job portal:

```bash
grit generate resource Category \
  --fields "name:string:unique,slug:slug:name"

grit generate resource Company \
  --fields "name:string,slug:slug:name,description:richtext,logo:string,website:string:optional,email:string"

grit generate resource Job \
  --fields "title:string,description:richtext,company_id:belongs_to:Company,location:string,salary_min:float,salary_max:float,type:string,experience:string,active:bool,deadline:date"

grit generate resource Application \
  --fields "job_id:belongs_to:Job,applicant_name:string,applicant_email:string,resume:string,cover_letter:text:optional,status:string"
```

With `--api`, each `grit generate resource` creates only Go files:
- `internal/models/<name>.go` -- GORM model struct
- `internal/services/<name>_service.go` -- CRUD business logic
- `internal/handlers/<name>_handler.go` -- REST endpoint handlers
- Injects routes into `internal/routes/routes.go`
- Injects model into `internal/database/database.go` for auto-migration

No Zod schemas, no TypeScript types, no React hooks, no admin pages.

## Step 5: Run Migrations and Seed

```bash
grit migrate
grit seed
```

This creates all tables and seeds an admin user:
- Email: `admin@example.com`
- Password: `password`

## Step 6: Explore the API

Restart the server to pick up the new resources:

```bash
cd apps/api && go run cmd/server/main.go
```

Open http://localhost:8080/docs -- you should see all endpoints for Jobs,
Companies, Applications, and Categories alongside the built-in Auth endpoints.

## Step 7: Authentication Flow

### Register

```bash
curl -s -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Alice Johnson",
    "email": "alice@example.com",
    "password": "securePass123!"
  }' | jq .
```

### Login

```bash
TOKEN=$(curl -s -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "alice@example.com",
    "password": "securePass123!"
  }' | jq -r '.data.access_token')

echo $TOKEN
```

### Get Current User

```bash
curl -s http://localhost:8080/api/auth/me \
  -H "Authorization: Bearer $TOKEN" | jq .
```

### Login as Admin

```bash
ADMIN_TOKEN=$(curl -s -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@example.com",
    "password": "password"
  }' | jq -r '.data.access_token')
```

## Step 8: CRUD Lifecycle

### Create a category

```bash
curl -s -X POST http://localhost:8080/api/categories \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -d '{"name": "Engineering"}' | jq .
```

### Create a company

```bash
COMPANY_ID=$(curl -s -X POST http://localhost:8080/api/companies \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -d '{
    "name": "Acme Corp",
    "description": "<p>Building the future.</p>",
    "logo": "",
    "website": "https://acme.example.com",
    "email": "jobs@acme.example.com"
  }' | jq -r '.data.id')

echo "Company ID: $COMPANY_ID"
```

### Create a job

```bash
JOB_ID=$(curl -s -X POST http://localhost:8080/api/jobs \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -d '{
    "title": "Senior Go Developer",
    "description": "<p>Build APIs with Grit.</p>",
    "company_id": "'$COMPANY_ID'",
    "location": "Remote",
    "salary_min": 120000,
    "salary_max": 180000,
    "type": "Full-time",
    "experience": "Senior",
    "active": true,
    "deadline": "2026-06-01"
  }' | jq -r '.data.id')

echo "Job ID: $JOB_ID"
```

### List jobs with pagination

```bash
curl -s "http://localhost:8080/api/jobs?page=1&page_size=10&sort=-created_at" | jq .
```

### Get a single job

```bash
curl -s "http://localhost:8080/api/jobs/$JOB_ID" | jq .
```

### Update a job

```bash
curl -s -X PUT "http://localhost:8080/api/jobs/$JOB_ID" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -d '{"salary_max": 200000}' | jq .
```

### Delete a job

```bash
curl -s -X DELETE "http://localhost:8080/api/jobs/$JOB_ID" \
  -H "Authorization: Bearer $ADMIN_TOKEN" | jq .
```

## Step 9: File Upload via Presigned URL

Grit uses presigned URLs for file uploads. The client never sends files to the
Go API -- it uploads directly to S3/MinIO.

```bash
# 1. Request a presigned URL
UPLOAD=$(curl -s -X POST http://localhost:8080/api/uploads/presign \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "filename": "resume.pdf",
    "content_type": "application/pdf"
  }')

UPLOAD_URL=$(echo $UPLOAD | jq -r '.data.upload_url')
FILE_URL=$(echo $UPLOAD | jq -r '.data.file_url')

# 2. Upload the file directly to MinIO/S3
curl -X PUT "$UPLOAD_URL" \
  -H "Content-Type: application/pdf" \
  --data-binary @resume.pdf

# 3. Use the file URL in your application
curl -s -X POST http://localhost:8080/api/applications \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "job_id": "'$JOB_ID'",
    "applicant_name": "Alice Johnson",
    "applicant_email": "alice@example.com",
    "resume": "'$FILE_URL'",
    "cover_letter": "I am excited to apply..."
  }' | jq .
```

## Step 10: Application Status Workflow

Applications move through these statuses:

```
Pending -> Reviewed -> Interview -> Accepted
                                 -> Rejected
```

```bash
# List all applications (admin only)
curl -s http://localhost:8080/api/applications \
  -H "Authorization: Bearer $ADMIN_TOKEN" | jq .

# Update status to "Interview"
curl -s -X PUT "http://localhost:8080/api/applications/<app-id>" \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"status": "Interview"}' | jq .
```

Status update sends an email notification to the applicant (check Mailhog at
http://localhost:8025 in development).

## Step 11: Two-Factor Authentication

```bash
# Setup TOTP (returns secret + QR URL)
curl -s -X POST http://localhost:8080/api/auth/totp/setup \
  -H "Authorization: Bearer $TOKEN" | jq .

# Scan the QR code with your authenticator app, then verify
curl -s -X POST http://localhost:8080/api/auth/totp/verify \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"code": "123456"}' | jq .
```

After enabling 2FA, login requires a TOTP code in addition to the password.

## Step 12: Testing with API Docs

The interactive API documentation at http://localhost:8080/docs lets you:
- Try every endpoint with a visual form
- Set the `Authorization` header once and reuse it
- See request/response schemas
- Download the OpenAPI spec for code generation

If you prefer Postman, import the OpenAPI spec from
`http://localhost:8080/docs/openapi.json`.

## Step 13: Deployment

### Option 1: grit deploy

The fastest path to production:

```bash
grit deploy --host deploy@server.com --domain api.jobs.example.com
```

This:
1. Cross-compiles the Go binary for Linux
2. Copies it to your server via SSH
3. Creates a systemd service
4. Configures Caddy for automatic HTTPS

### Option 2: Docker Compose

```bash
# Copy .env.example to .env and fill in production values
cp .env.example .env
# Edit .env with real DATABASE_URL, JWT_SECRET, etc.

docker compose -f docker-compose.prod.yml up -d
```

### Option 3: Bare binary

```bash
cd apps/api
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server cmd/server/main.go
```

The output is a single static binary. Copy it anywhere, set environment
variables, and run it. Zero dependencies at runtime.

## Summary

What we built:
- A complete REST API with 4 resources (Job, Company, Application, Category)
- JWT authentication with OAuth and TOTP 2FA
- Presigned URL file uploads
- Role-based access control
- Paginated, sortable, filterable listings
- Email notifications
- Interactive API documentation
- Production-ready deployment

All without writing a single line of frontend code.
