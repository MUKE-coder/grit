import Link from 'next/link'
import { ArrowLeft, ArrowRight } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'
import { CodeBlock } from '@/components/code-block'
import { getDocMetadata } from '@/config/docs-metadata'

export const metadata = getDocMetadata('/docs/backend/response-format')

export default function ResponseFormatPage() {
  return (
    <div className="min-h-screen bg-background">
      <SiteHeader />
      <DocsSidebar />
      <main className="lg:pl-64">
        <div className="container max-w-screen-xl py-10 px-6">
          <div className="max-w-3xl">
            {/* Header */}
            <div className="mb-10">
              <span className="tag-mono text-primary/80 mb-3 block">Backend (Go API)</span>
              <h1 className="text-4xl font-bold tracking-tight mb-4">
                API Response Format
              </h1>
              <p className="text-lg text-muted-foreground leading-relaxed">
                All Grit API endpoints follow a consistent response format. This makes it
                predictable for frontend consumers and ensures error handling is uniform
                across the entire application.
              </p>
            </div>

            <div className="prose-grit">
              {/* ── Success (Single Item) ─────────────────────────────── */}
              <h2 id="success-single">Success Response (Single Item)</h2>
              <p>
                When an endpoint returns a single resource, the response wraps it in a
                <code>data</code> field. An optional <code>message</code> field provides
                a human-readable description of what happened.
              </p>
              <CodeBlock filename="GET /api/users/1 -- 200 OK" code={`{
    "data": {
        "id": 1,
        "name": "John Doe",
        "email": "john@example.com",
        "role": "admin",
        "avatar": "",
        "active": true,
        "email_verified_at": null,
        "created_at": "2026-02-11T10:00:00Z",
        "updated_at": "2026-02-11T10:00:00Z"
    }
}`} />

              <p>
                For create and update operations, include a <code>message</code> field:
              </p>
              <CodeBlock filename="POST /api/posts -- 201 Created" code={`{
    "data": {
        "id": 42,
        "title": "Getting Started with Grit",
        "slug": "getting-started-with-grit",
        "body": "Grit is a full-stack meta-framework...",
        "published": false,
        "author_id": 1,
        "created_at": "2026-02-11T14:30:00Z",
        "updated_at": "2026-02-11T14:30:00Z"
    },
    "message": "Post created successfully"
}`} />

              <p>In Go, this looks like:</p>
              <CodeBlock filename="handler.go" code={`// Single item with message
c.JSON(http.StatusCreated, gin.H{
    "data":    post,
    "message": "Post created successfully",
})

// Single item without message
c.JSON(http.StatusOK, gin.H{
    "data": user,
})`} />

              {/* ── Success (List) ─────────────────────────────── */}
              <h2 id="success-list">Success Response (List with Pagination)</h2>
              <p>
                List endpoints return an array of resources in the <code>data</code> field and
                pagination metadata in the <code>meta</code> field.
              </p>
              <CodeBlock filename="GET /api/users?page=2&page_size=10 -- 200 OK" code={`{
    "data": [
        {
            "id": 11,
            "name": "Alice Smith",
            "email": "alice@example.com",
            "role": "user",
            "active": true,
            "created_at": "2026-02-10T09:00:00Z",
            "updated_at": "2026-02-10T09:00:00Z"
        },
        {
            "id": 12,
            "name": "Bob Johnson",
            "email": "bob@example.com",
            "role": "editor",
            "active": true,
            "created_at": "2026-02-10T10:30:00Z",
            "updated_at": "2026-02-10T10:30:00Z"
        }
    ],
    "meta": {
        "total": 57,
        "page": 2,
        "page_size": 10,
        "pages": 6
    }
}`} />

              {/* ── Pagination Meta ─────────────────────────────── */}
              <h3 id="pagination-meta">Pagination Meta Structure</h3>
              <div className="rounded-lg border border-border/30 bg-card/30 overflow-hidden mb-6">
                <table className="w-full text-sm">
                  <thead>
                    <tr className="border-b border-border/30 bg-accent/20">
                      <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Field</th>
                      <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Type</th>
                      <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Description</th>
                    </tr>
                  </thead>
                  <tbody className="text-muted-foreground">
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5 font-mono text-xs">total</td>
                      <td className="px-4 py-2.5 font-mono text-xs">integer</td>
                      <td className="px-4 py-2.5">Total number of records matching the query (before pagination)</td>
                    </tr>
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5 font-mono text-xs">page</td>
                      <td className="px-4 py-2.5 font-mono text-xs">integer</td>
                      <td className="px-4 py-2.5">Current page number (1-based)</td>
                    </tr>
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5 font-mono text-xs">page_size</td>
                      <td className="px-4 py-2.5 font-mono text-xs">integer</td>
                      <td className="px-4 py-2.5">Number of records per page (max 100)</td>
                    </tr>
                    <tr>
                      <td className="px-4 py-2.5 font-mono text-xs">pages</td>
                      <td className="px-4 py-2.5 font-mono text-xs">integer</td>
                      <td className="px-4 py-2.5">Total number of pages (ceil(total / page_size))</td>
                    </tr>
                  </tbody>
                </table>
              </div>
              <p>In Go:</p>
              <CodeBlock filename="handler.go" code={`pages := int(math.Ceil(float64(total) / float64(pageSize)))

c.JSON(http.StatusOK, gin.H{
    "data": users,
    "meta": gin.H{
        "total":     total,
        "page":      page,
        "page_size": pageSize,
        "pages":     pages,
    },
})`} />

              {/* ── Error Response ─────────────────────────────── */}
              <h2 id="error-response">Error Response</h2>
              <p>
                All errors follow the same envelope format with an <code>error</code> object
                containing a machine-readable <code>code</code>, a human-readable <code>message</code>,
                and optional <code>details</code> for field-level validation errors.
              </p>
              <CodeBlock filename="422 Unprocessable Entity" code={`{
    "error": {
        "code": "VALIDATION_ERROR",
        "message": "Key: 'Email' Error:Field validation for 'Email' failed on the 'required' tag",
        "details": {
            "email": "This field is required",
            "password": "Must be at least 8 characters"
        }
    }
}`} />

              <p>Simple error (no field details):</p>
              <CodeBlock filename="404 Not Found" code={`{
    "error": {
        "code": "NOT_FOUND",
        "message": "User not found"
    }
}`} />

              <p>In Go:</p>
              <CodeBlock filename="handler.go" code={`// Simple error
c.JSON(http.StatusNotFound, gin.H{
    "error": gin.H{
        "code":    "NOT_FOUND",
        "message": "User not found",
    },
})

// Error with field details
c.JSON(http.StatusUnprocessableEntity, gin.H{
    "error": gin.H{
        "code":    "VALIDATION_ERROR",
        "message": err.Error(),
        "details": gin.H{
            "email":    "This field is required",
            "password": "Must be at least 8 characters",
        },
    },
})`} />

              {/* ── Delete / Action Response ─────────────────────────────── */}
              <h2 id="action-response">Action Response (Delete, Logout, etc.)</h2>
              <p>
                For operations that do not return a resource (like delete or logout),
                return only a <code>message</code> field:
              </p>
              <CodeBlock filename="DELETE /api/users/5 -- 200 OK" code={`{
    "message": "User deleted successfully"
}`} />

              {/* ── HTTP Status Codes ─────────────────────────────── */}
              <h2 id="status-codes">HTTP Status Codes</h2>
              <p>
                Grit uses standard HTTP status codes consistently across all endpoints:
              </p>
              <div className="rounded-lg border border-border/30 bg-card/30 overflow-hidden mb-6">
                <table className="w-full text-sm">
                  <thead>
                    <tr className="border-b border-border/30 bg-accent/20">
                      <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Code</th>
                      <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Name</th>
                      <th className="text-left px-4 py-2.5 font-medium text-foreground/80">When Used</th>
                    </tr>
                  </thead>
                  <tbody className="text-muted-foreground">
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5 font-mono text-xs text-green-400">200</td>
                      <td className="px-4 py-2.5">OK</td>
                      <td className="px-4 py-2.5">Successful GET, PUT, DELETE requests</td>
                    </tr>
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5 font-mono text-xs text-green-400">201</td>
                      <td className="px-4 py-2.5">Created</td>
                      <td className="px-4 py-2.5">Successful POST that creates a resource</td>
                    </tr>
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5 font-mono text-xs text-yellow-400">400</td>
                      <td className="px-4 py-2.5">Bad Request</td>
                      <td className="px-4 py-2.5">Malformed request body or invalid parameters</td>
                    </tr>
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5 font-mono text-xs text-yellow-400">401</td>
                      <td className="px-4 py-2.5">Unauthorized</td>
                      <td className="px-4 py-2.5">Missing, invalid, or expired JWT token</td>
                    </tr>
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5 font-mono text-xs text-yellow-400">403</td>
                      <td className="px-4 py-2.5">Forbidden</td>
                      <td className="px-4 py-2.5">Authenticated but lacks permission (wrong role)</td>
                    </tr>
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5 font-mono text-xs text-yellow-400">404</td>
                      <td className="px-4 py-2.5">Not Found</td>
                      <td className="px-4 py-2.5">Resource does not exist</td>
                    </tr>
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5 font-mono text-xs text-orange-400">409</td>
                      <td className="px-4 py-2.5">Conflict</td>
                      <td className="px-4 py-2.5">Duplicate entry (e.g., email already registered)</td>
                    </tr>
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5 font-mono text-xs text-orange-400">422</td>
                      <td className="px-4 py-2.5">Unprocessable Entity</td>
                      <td className="px-4 py-2.5">Validation errors on request fields</td>
                    </tr>
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5 font-mono text-xs text-orange-400">429</td>
                      <td className="px-4 py-2.5">Too Many Requests</td>
                      <td className="px-4 py-2.5">Rate limit exceeded</td>
                    </tr>
                    <tr>
                      <td className="px-4 py-2.5 font-mono text-xs text-red-400">500</td>
                      <td className="px-4 py-2.5">Internal Server Error</td>
                      <td className="px-4 py-2.5">Unexpected server-side error</td>
                    </tr>
                  </tbody>
                </table>
              </div>

              {/* ── Error Codes ─────────────────────────────── */}
              <h2 id="error-codes">Error Codes</h2>
              <p>
                Error codes are machine-readable strings that the frontend can use to
                display localized messages or take programmatic action. They are always
                SCREAMING_SNAKE_CASE.
              </p>
              <div className="rounded-lg border border-border/30 bg-card/30 overflow-hidden mb-6">
                <table className="w-full text-sm">
                  <thead>
                    <tr className="border-b border-border/30 bg-accent/20">
                      <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Error Code</th>
                      <th className="text-left px-4 py-2.5 font-medium text-foreground/80">HTTP Status</th>
                      <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Description</th>
                    </tr>
                  </thead>
                  <tbody className="text-muted-foreground">
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5 font-mono text-xs">VALIDATION_ERROR</td>
                      <td className="px-4 py-2.5 font-mono text-xs">422</td>
                      <td className="px-4 py-2.5">One or more request fields failed validation</td>
                    </tr>
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5 font-mono text-xs">NOT_FOUND</td>
                      <td className="px-4 py-2.5 font-mono text-xs">404</td>
                      <td className="px-4 py-2.5">The requested resource does not exist</td>
                    </tr>
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5 font-mono text-xs">UNAUTHORIZED</td>
                      <td className="px-4 py-2.5 font-mono text-xs">401</td>
                      <td className="px-4 py-2.5">Authentication is required or the token is invalid</td>
                    </tr>
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5 font-mono text-xs">FORBIDDEN</td>
                      <td className="px-4 py-2.5 font-mono text-xs">403</td>
                      <td className="px-4 py-2.5">Authenticated but insufficient permissions (role check failed)</td>
                    </tr>
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5 font-mono text-xs">INTERNAL_ERROR</td>
                      <td className="px-4 py-2.5 font-mono text-xs">500</td>
                      <td className="px-4 py-2.5">An unexpected server-side error occurred</td>
                    </tr>
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5 font-mono text-xs">CONFLICT</td>
                      <td className="px-4 py-2.5 font-mono text-xs">409</td>
                      <td className="px-4 py-2.5">A resource with the same unique key already exists</td>
                    </tr>
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5 font-mono text-xs">INVALID_CREDENTIALS</td>
                      <td className="px-4 py-2.5 font-mono text-xs">401</td>
                      <td className="px-4 py-2.5">Email/password combination is incorrect</td>
                    </tr>
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5 font-mono text-xs">INVALID_TOKEN</td>
                      <td className="px-4 py-2.5 font-mono text-xs">401</td>
                      <td className="px-4 py-2.5">The refresh token is invalid or expired</td>
                    </tr>
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5 font-mono text-xs">EMAIL_EXISTS</td>
                      <td className="px-4 py-2.5 font-mono text-xs">409</td>
                      <td className="px-4 py-2.5">An account with this email already exists</td>
                    </tr>
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5 font-mono text-xs">ACCOUNT_DISABLED</td>
                      <td className="px-4 py-2.5 font-mono text-xs">403</td>
                      <td className="px-4 py-2.5">The user account has been deactivated</td>
                    </tr>
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5 font-mono text-xs">TOKEN_ERROR</td>
                      <td className="px-4 py-2.5 font-mono text-xs">500</td>
                      <td className="px-4 py-2.5">Failed to generate JWT tokens</td>
                    </tr>
                    <tr>
                      <td className="px-4 py-2.5 font-mono text-xs">RATE_LIMITED</td>
                      <td className="px-4 py-2.5 font-mono text-xs">429</td>
                      <td className="px-4 py-2.5">Too many requests from the same IP address</td>
                    </tr>
                  </tbody>
                </table>
              </div>

              {/* ── Full Examples ─────────────────────────────── */}
              <h2 id="full-examples">Full JSON Examples</h2>

              <h3 id="example-register">Register (Success)</h3>
              <CodeBlock filename="POST /api/auth/register -- 201 Created" code={`{
    "data": {
        "user": {
            "id": 1,
            "name": "John Doe",
            "email": "john@example.com",
            "role": "user",
            "avatar": "",
            "active": true,
            "email_verified_at": null,
            "created_at": "2026-02-11T10:00:00Z",
            "updated_at": "2026-02-11T10:00:00Z"
        },
        "tokens": {
            "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
            "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
            "expires_at": 1707649200
        }
    },
    "message": "User registered successfully"
}`} />

              <h3 id="example-register-error">Register (Email Taken)</h3>
              <CodeBlock filename="POST /api/auth/register -- 409 Conflict" code={`{
    "error": {
        "code": "EMAIL_EXISTS",
        "message": "A user with this email already exists"
    }
}`} />

              <h3 id="example-validation">Validation Error</h3>
              <CodeBlock filename="POST /api/posts -- 422 Unprocessable Entity" code={`{
    "error": {
        "code": "VALIDATION_ERROR",
        "message": "Key: 'Title' Error:Field validation for 'Title' failed on the 'required' tag"
    }
}`} />

              <h3 id="example-unauthorized">Unauthorized (Missing Token)</h3>
              <CodeBlock filename="GET /api/posts -- 401 Unauthorized" code={`{
    "error": {
        "code": "UNAUTHORIZED",
        "message": "Authorization header is required"
    }
}`} />

              <h3 id="example-forbidden">Forbidden (Insufficient Role)</h3>
              <CodeBlock filename="DELETE /api/users/3 -- 403 Forbidden" code={`{
    "error": {
        "code": "FORBIDDEN",
        "message": "You do not have permission to access this resource"
    }
}`} />

              <h3 id="example-list">Paginated List</h3>
              <CodeBlock filename="GET /api/users?page=1&page_size=2&search=john -- 200 OK" code={`{
    "data": [
        {
            "id": 1,
            "name": "John Doe",
            "email": "john@example.com",
            "role": "admin",
            "avatar": "",
            "active": true,
            "email_verified_at": null,
            "created_at": "2026-02-11T10:00:00Z",
            "updated_at": "2026-02-11T10:00:00Z"
        },
        {
            "id": 15,
            "name": "Johnny Appleseed",
            "email": "johnny@example.com",
            "role": "user",
            "avatar": "",
            "active": true,
            "email_verified_at": "2026-02-11T12:00:00Z",
            "created_at": "2026-02-11T11:00:00Z",
            "updated_at": "2026-02-11T11:00:00Z"
        }
    ],
    "meta": {
        "total": 3,
        "page": 1,
        "page_size": 2,
        "pages": 2
    }
}`} />

              <h3 id="example-server-error">Internal Server Error</h3>
              <CodeBlock filename="POST /api/posts -- 500 Internal Server Error" code={`{
    "error": {
        "code": "INTERNAL_ERROR",
        "message": "Failed to create post"
    }
}`} />

              {/* ── Frontend Consumption ─────────────────────────────── */}
              <h2 id="frontend-consumption">Consuming on the Frontend</h2>
              <p>
                Because the format is consistent, your React Query hooks can use a single
                error handler and response parser:
              </p>
              <CodeBlock language="typescript" filename="apps/web/hooks/use-posts.ts" code={`import { useQuery, useMutation } from '@tanstack/react-query';
import api from '@/lib/api-client';

interface PaginatedResponse<T> {
  data: T[];
  meta: {
    total: number;
    page: number;
    page_size: number;
    pages: number;
  };
}

interface ApiError {
  error: {
    code: string;
    message: string;
    details?: Record<string, string>;
  };
}

export function usePosts(page = 1, pageSize = 20) {
  return useQuery({
    queryKey: ['posts', page, pageSize],
    queryFn: async () => {
      const { data } = await api.get<PaginatedResponse<Post>>(
        \`/api/posts?page=\${page}&page_size=\${pageSize}\`
      );
      return data;
    },
  });
}

export function useCreatePost() {
  return useMutation({
    mutationFn: async (post: CreatePostInput) => {
      const { data } = await api.post('/api/posts', post);
      return data.data; // unwrap the "data" envelope
    },
    onError: (error: any) => {
      const apiError = error.response?.data as ApiError;
      // apiError.error.code === "VALIDATION_ERROR"
      // apiError.error.message === "..."
      // apiError.error.details?.title === "..."
    },
  });
}`} />

              {/* ── Summary ─────────────────────────────── */}
              <h2 id="summary">Response Format Summary</h2>
              <div className="rounded-lg border border-border/30 bg-card/30 overflow-hidden mb-6">
                <table className="w-full text-sm">
                  <thead>
                    <tr className="border-b border-border/30 bg-accent/20">
                      <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Scenario</th>
                      <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Shape</th>
                    </tr>
                  </thead>
                  <tbody className="text-muted-foreground">
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5">Single resource</td>
                      <td className="px-4 py-2.5 font-mono text-xs">{`{ "data": { ... }, "message"?: "..." }`}</td>
                    </tr>
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5">List (paginated)</td>
                      <td className="px-4 py-2.5 font-mono text-xs">{`{ "data": [...], "meta": { total, page, page_size, pages } }`}</td>
                    </tr>
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5">Action (delete, logout)</td>
                      <td className="px-4 py-2.5 font-mono text-xs">{`{ "message": "..." }`}</td>
                    </tr>
                    <tr>
                      <td className="px-4 py-2.5">Error</td>
                      <td className="px-4 py-2.5 font-mono text-xs">{`{ "error": { "code": "...", "message": "...", "details"?: { ... } } }`}</td>
                    </tr>
                  </tbody>
                </table>
              </div>
            </div>

            {/* Nav */}
            <div className="flex items-center justify-between pt-6 mt-10 border-t border-border/30">
              <Button variant="ghost" size="sm" asChild className="text-muted-foreground/60 hover:text-foreground">
                <Link href="/docs/backend/authentication" className="gap-1.5">
                  <ArrowLeft className="h-3.5 w-3.5" />
                  Authentication
                </Link>
              </Button>
              <Button variant="ghost" size="sm" asChild className="text-muted-foreground/60 hover:text-foreground">
                <Link href="/docs/admin/overview" className="gap-1.5">
                  Admin Panel
                  <ArrowRight className="h-3.5 w-3.5" />
                </Link>
              </Button>
            </div>
          </div>
        </div>
      </main>
    </div>
  )
}
