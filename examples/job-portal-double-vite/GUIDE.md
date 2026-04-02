# Build Guide: Job Portal (Double + Vite)

Step-by-step guide to building a job portal with Grit's Double architecture
and TanStack Router. Two apps: a Go API and a React SPA. No separate admin
panel -- ADMIN users manage everything from role-protected pages in the web app.

---

## Step 1: Scaffold the Project

```bash
grit new job-portal --double --vite
cd job-portal
```

This creates a Turborepo monorepo with `apps/api` (Go) and `apps/web`
(React + TanStack Router + Vite). No `apps/admin` directory -- that is
the key difference from the Triple architecture.

Start infrastructure:

```bash
docker compose up -d   # PostgreSQL, Redis, MinIO, Mailhog
```

## Step 2: Generate Resources

```bash
grit generate resource Category --fields "name:string:unique,slug:slug:name"

grit generate resource Company --fields \
  "name:string,slug:slug:name,description:richtext,logo:string,website:string:optional,email:string"

grit generate resource Job --fields \
  "title:string,description:richtext,company_id:belongs_to:Company,location:string,salary_min:float,salary_max:float,type:string,experience:string,active:bool,deadline:date"

grit generate resource Application --fields \
  "job_id:belongs_to:Job,applicant_name:string,applicant_email:string,resume:string,cover_letter:text:optional,status:string"
```

Each command generates:
- Go model, service, handler in `apps/api/internal/`
- Zod schemas and TypeScript types in `packages/shared/`
- React Query hooks in `apps/web/src/hooks/`
- Route injections in `apps/api/internal/routes/routes.go`

## Step 3: Run Migrations and Seed

```bash
grit migrate
grit seed
```

Start both servers:

```bash
# Terminal 1
cd apps/api && go run cmd/server/main.go

# Terminal 2
pnpm install && pnpm dev
```

Web app: http://localhost:5173 | API: http://localhost:8080

## Step 4: Understand the Routing (TanStack Router)

TanStack Router uses file-based routing in `src/routes/`. The key files:

```
src/routes/
├── __root.tsx         # Root layout (navbar, footer, auth context)
├── index.tsx          # "/" — landing page
├── jobs/
│   ├── index.tsx      # "/jobs" — listings with filters
│   └── $id.tsx        # "/jobs/:id" — job detail
├── companies/
│   ├── index.tsx      # "/companies" — directory
│   └── $slug.tsx      # "/companies/:slug" — profile
├── apply/
│   └── $jobId.tsx     # "/apply/:jobId" — application form
├── auth/
│   ├── login.tsx      # "/auth/login"
│   └── register.tsx   # "/auth/register"
└── admin/             # Role-protected management pages
    ├── index.tsx       # "/admin" — dashboard
    ├── jobs.tsx        # "/admin/jobs" — CRUD table
    ├── companies.tsx   # "/admin/companies"
    ├── applications.tsx
    └── categories.tsx
```

Route params use the `$param` convention. A file named `$id.tsx` maps to
`:id` in the URL.

## Step 5: Build the Landing Page

Edit `src/routes/index.tsx`:

```tsx
import { createFileRoute } from "@tanstack/react-router";
import { useQuery } from "@tanstack/react-query";
import { api } from "../lib/api-client";

export const Route = createFileRoute("/")({
  component: HomePage,
});

function HomePage() {
  const { data: featuredJobs } = useQuery({
    queryKey: ["jobs", "featured"],
    queryFn: () => api.get("/api/jobs?active=true&page_size=6"),
  });

  const { data: categories } = useQuery({
    queryKey: ["categories"],
    queryFn: () => api.get("/api/categories"),
  });

  return (
    <div className="min-h-screen bg-background">
      <section className="py-20 text-center">
        <h1 className="text-5xl font-bold text-foreground">
          Find Your Next Opportunity
        </h1>
        <p className="mt-4 text-text-secondary">
          Browse thousands of jobs from top companies
        </p>
        <JobSearchBar />
      </section>
      <section className="mx-auto max-w-6xl px-4 py-12">
        <h2 className="text-2xl font-semibold text-foreground">Categories</h2>
        <div className="mt-6 grid grid-cols-2 gap-4 md:grid-cols-4">
          {categories?.data?.map((cat) => (
            <CategoryCard key={cat.id} category={cat} />
          ))}
        </div>
      </section>
      <section className="mx-auto max-w-6xl px-4 py-12">
        <h2 className="text-2xl font-semibold text-foreground">Featured Jobs</h2>
        <div className="mt-6 grid gap-4 md:grid-cols-2 lg:grid-cols-3">
          {featuredJobs?.data?.map((job) => (
            <JobCard key={job.id} job={job} />
          ))}
        </div>
      </section>
    </div>
  );
}
```

## Step 6: Build Job Listings and Detail Pages

`src/routes/jobs/index.tsx` -- search, filter by location, salary range:

```tsx
export const Route = createFileRoute("/jobs/")({
  component: JobListings,
  validateSearch: (search) => ({
    q: (search.q as string) || "",
    location: (search.location as string) || "",
    type: (search.type as string) || "",
    page: Number(search.page) || 1,
  }),
});

function JobListings() {
  const { q, location, type, page } = Route.useSearch();
  const { data, isLoading } = useQuery({
    queryKey: ["jobs", { q, location, type, page }],
    queryFn: () => api.get("/api/jobs", { params: { q, location, type, page } }),
  });
  // Render search bar, filter sidebar, job cards, pagination
}
```

`src/routes/jobs/$id.tsx` -- job detail with apply button:

```tsx
export const Route = createFileRoute("/jobs/$id")({
  component: JobDetail,
});

function JobDetail() {
  const { id } = Route.useParams();
  const { data: job } = useQuery({
    queryKey: ["jobs", id],
    queryFn: () => api.get("/api/jobs/" + id),
  });
  // Render job title, company, description, salary, apply button
}
```

## Step 7: Build Company Profiles

`src/routes/companies/$slug.tsx`:

```tsx
export const Route = createFileRoute("/companies/$slug")({
  component: CompanyProfile,
});

function CompanyProfile() {
  const { slug } = Route.useParams();
  const { data: company } = useQuery({
    queryKey: ["companies", slug],
    queryFn: () => api.get("/api/companies/slug/" + slug),
  });
  // Render company logo, description, website, open positions
}
```

## Step 8: Build the Application Flow

`src/routes/apply/$jobId.tsx` -- form with resume upload:

```tsx
export const Route = createFileRoute("/apply/$jobId")({
  beforeLoad: ({ context }) => {
    if (!context.auth.user) {
      throw redirect({ to: "/auth/login" });
    }
  },
  component: ApplyPage,
});

function ApplyPage() {
  const { jobId } = Route.useParams();
  const mutation = useMutation({
    mutationFn: (formData: FormData) =>
      api.post("/api/applications", formData),
    onSuccess: () => navigate({ to: "/jobs/" + jobId }),
  });
  // Render form: name, email, resume upload, cover letter, submit
}
```

## Step 9: Protect Admin Routes (No Separate Admin Panel)

This is the core pattern for Double architecture. Create a layout route at
`src/routes/admin.tsx` that guards all `/admin/*` pages:

```tsx
import { createFileRoute, Outlet, redirect } from "@tanstack/react-router";
import { AdminSidebar } from "../components/admin-sidebar";

export const Route = createFileRoute("/admin")({
  beforeLoad: ({ context }) => {
    if (!context.auth.user) {
      throw redirect({ to: "/auth/login" });
    }
    if (context.auth.user.role !== "ADMIN") {
      throw redirect({ to: "/" });
    }
  },
  component: AdminLayout,
});

function AdminLayout() {
  return (
    <div className="flex min-h-screen bg-background">
      <AdminSidebar />
      <main className="flex-1 overflow-auto p-6">
        <Outlet />
      </main>
    </div>
  );
}
```

The `beforeLoad` guard runs before rendering. Non-admin users are redirected
to the home page. No admin content is ever sent to unauthorized users.

## Step 10: Build Admin Management Pages

`src/routes/admin/jobs.tsx` -- full CRUD with DataTable:

```tsx
export const Route = createFileRoute("/admin/jobs")({
  component: AdminJobs,
});

function AdminJobs() {
  const { data, isLoading } = useQuery({
    queryKey: ["admin", "jobs"],
    queryFn: () => api.get("/api/jobs"),
  });

  return (
    <div>
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-bold text-foreground">Manage Jobs</h1>
        <Button onClick={() => setShowForm(true)}>Add Job</Button>
      </div>
      <DataTable
        columns={jobColumns}
        data={data?.data || []}
        onEdit={handleEdit}
        onDelete={handleDelete}
      />
      {showForm && <JobFormDialog onClose={() => setShowForm(false)} />}
    </div>
  );
}
```

`src/routes/admin/applications.tsx` -- review and update status:

```tsx
function AdminApplications() {
  const updateStatus = useMutation({
    mutationFn: ({ id, status }) =>
      api.patch("/api/applications/" + id, { status }),
    onSuccess: () => queryClient.invalidateQueries(["admin", "applications"]),
  });

  // Render DataTable with status dropdown: Pending, Reviewed, Interview, Accepted, Rejected
}
```

Repeat for `admin/companies.tsx` and `admin/categories.tsx`.

## Step 11: Conditional Navigation for Admin Users

In the root layout (`src/routes/__root.tsx`), show an Admin link only for
ADMIN users:

```tsx
function RootLayout() {
  const { user } = useAuth();
  return (
    <div className="min-h-screen bg-background">
      <nav className="flex items-center justify-between border-b border-border px-6 py-3">
        <div className="flex items-center gap-6">
          <Link to="/" className="text-lg font-bold text-foreground">JobPortal</Link>
          <Link to="/jobs" className="text-text-secondary hover:text-foreground">Jobs</Link>
          <Link to="/companies" className="text-text-secondary hover:text-foreground">Companies</Link>
          {user?.role === "ADMIN" && (
            <Link to="/admin" className="text-accent hover:text-accent-hover">Admin</Link>
          )}
        </div>
        <AuthButtons />
      </nav>
      <Outlet />
    </div>
  );
}
```

## Step 12: File Uploads

Company logos and resumes go through the API's upload handler:

```tsx
async function uploadFile(file: File): Promise<string> {
  const formData = new FormData();
  formData.append("file", file);
  const res = await api.post("/api/uploads", formData);
  return res.data.url;
}
```

The API stores files in MinIO (dev) or S3/R2 (prod) and returns a public URL.

## Step 13: Email Notifications

The API sends emails via Resend (configured in `.env`):
- **Welcome email** -- on user registration
- **Application confirmation** -- when a user submits an application
- **Status update** -- when an admin changes application status

These are triggered by background jobs (asynq) so they do not block the API
response.

## Step 14: Build for Production

```bash
# Build everything
turbo build

# Or build individually
cd apps/api && go build -o job-portal-api cmd/server/main.go
cd apps/web && pnpm build   # outputs dist/
```

## Step 15: Deploy

### Docker Compose (recommended)

```bash
docker compose -f docker-compose.prod.yml up -d
```

This runs:
- `api` -- Go binary on port 8080
- `web` -- Nginx serving the Vite build on port 80
- `postgres` -- PostgreSQL 16
- `redis` -- Redis 7
- `minio` -- S3-compatible storage

Two containers for your code (api + web) instead of three (no admin).

### VPS with grit deploy

```bash
grit deploy --host deploy@server.com --domain jobs.example.com
```

### Static + API split

Deploy the API to any VPS and the Vite build (`apps/web/dist/`) to any
static hosting provider (Vercel, Netlify, Cloudflare Pages). Set `VITE_API_URL`
to point to your API.

---

## Summary

The Double architecture removes the separate admin panel and folds management
pages into the web app behind role guards. This gives you:

- **Simpler deployment** -- 2 apps instead of 3
- **Less code** -- no duplicate layouts, auth, or API client setup
- **Same features** -- ADMIN users get full CRUD via `/admin/*` routes
- **TanStack Router** -- type-safe routing, fast Vite builds, SPA

For apps that need a dedicated admin experience with its own layout and
navigation, use the Triple architecture instead (`grit new job-portal --triple --vite`).
