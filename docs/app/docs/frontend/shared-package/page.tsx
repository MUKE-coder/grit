import Link from 'next/link'
import { ArrowRight, ArrowLeft } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'
import { CodeBlock } from '@/components/code-block'
import { getDocMetadata } from '@/config/docs-metadata'

export const metadata = getDocMetadata('/docs/frontend/shared-package')

export default function SharedPackagePage() {
  return (
    <div className="min-h-screen bg-background isolate">
      <SiteHeader />
      <DocsSidebar />

      <main className="lg:pl-64">
        <div className="container max-w-screen-xl py-10 px-6">
          <div className="max-w-3xl">
            {/* Header */}
            <div className="mb-10">
              <span className="tag-mono text-primary/80 mb-3 block">Frontend</span>
              <h1 className="text-4xl font-bold tracking-tight mb-4">
                Shared Package
              </h1>
              <p className="text-lg text-muted-foreground leading-relaxed">
                The shared package at <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">packages/shared/</code> contains
                Zod validation schemas, TypeScript types, and constants shared between the web app, admin panel,
                and any other frontend app in the monorepo.
              </p>
            </div>

            <div className="prose-grit">
              {/* Package Structure */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Package Structure
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  The shared package is a pnpm workspace member that both <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">apps/web</code> and <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">apps/admin</code> depend
                  on. Import from it using the workspace alias (e.g., <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">@myapp/shared/schemas</code>).
                </p>

                <CodeBlock filename="packages/shared/" code={`packages/shared/
├── package.json
├── tsconfig.json
├── schemas/              # Zod validation schemas
│   ├── user.ts           # User auth schemas (Login, Register, etc.)
│   ├── post.ts           # Generated: CreatePost, UpdatePost schemas
│   ├── blog.ts           # Generated: CreateBlog, UpdateBlog schemas
│   └── index.ts          # Re-exports all schemas
├── types/                # TypeScript interfaces
│   ├── user.ts           # User, AuthResponse interfaces
│   ├── api.ts            # ApiResponse, PaginatedResponse, ApiError
│   ├── upload.ts         # Upload interface
│   ├── post.ts           # Generated: Post interface
│   ├── blog.ts           # Generated: Blog interface
│   └── index.ts          # Re-exports all types
└── constants/
    └── index.ts          # ROLES, API_ROUTES, etc.`} />

                <CodeBlock language="json" filename="packages/shared/package.json" code={`{
  "name": "@myapp/shared",
  "version": "0.8.0",
  "private": true,
  "main": "./index.ts",
  "types": "./index.ts",
  "exports": {
    "./schemas": "./schemas/index.ts",
    "./types": "./types/index.ts",
    "./constants": "./constants/index.ts"
  },
  "dependencies": {
    "zod": "^3.22.0"
  }
}`} />
              </div>

              {/* Zod Schemas */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Zod Schemas
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Each resource has Create and Update Zod schemas. The base User schemas are scaffolded
                  with the project. When you run <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">grit generate resource Post</code>,
                  a new schema file is generated and automatically added to the barrel export.
                </p>

                <h3 className="text-lg font-semibold tracking-tight mb-3 mt-6">User Schemas (Built-in)</h3>

                <CodeBlock language="typescript" filename="packages/shared/schemas/user.ts" code={`import { z } from "zod";

export const LoginSchema = z.object({
  email: z.string().email("Please enter a valid email"),
  password: z.string().min(1, "Password is required"),
});

export const RegisterSchema = z.object({
  name: z.string().min(2, "Name must be at least 2 characters"),
  email: z.string().email("Please enter a valid email"),
  password: z.string().min(8, "Password must be at least 8 characters"),
  confirmPassword: z.string(),
}).refine((data) => data.password === data.confirmPassword, {
  message: "Passwords do not match",
  path: ["confirmPassword"],
});

export const UpdateUserSchema = z.object({
  name: z.string().min(2).optional(),
  email: z.string().email().optional(),
  role: z.enum(["admin", "editor", "user"]).optional(),
  active: z.boolean().optional(),
});

export const ForgotPasswordSchema = z.object({
  email: z.string().email("Please enter a valid email"),
});

export const ResetPasswordSchema = z.object({
  token: z.string().min(1, "Token is required"),
  password: z.string().min(8, "Password must be at least 8 characters"),
});

// Inferred types from schemas
export type LoginInput = z.infer<typeof LoginSchema>;
export type RegisterInput = z.infer<typeof RegisterSchema>;
export type UpdateUserInput = z.infer<typeof UpdateUserSchema>;
export type ForgotPasswordInput = z.infer<typeof ForgotPasswordSchema>;
export type ResetPasswordInput = z.infer<typeof ResetPasswordSchema>;`} />

                <h3 className="text-lg font-semibold tracking-tight mb-3 mt-6">Generated Resource Schemas</h3>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Running <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">grit generate resource Post --fields &quot;title:string,content:text,published:bool&quot;</code> generates:
                </p>

                <CodeBlock language="typescript" filename="packages/shared/schemas/post.ts" code={`import { z } from "zod";

export const CreatePostSchema = z.object({
  title: z.string().min(1, "Title is required"),
  content: z.string().min(1, "Content is required"),
  published: z.boolean(),
});

export const UpdatePostSchema = z.object({
  title: z.string().min(1).optional(),
  content: z.string().min(1).optional(),
  published: z.boolean().optional(),
});

export type CreatePostInput = z.infer<typeof CreatePostSchema>;
export type UpdatePostInput = z.infer<typeof UpdatePostSchema>;`} />
              </div>

              {/* TypeScript Types */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  TypeScript Types
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Type interfaces mirror Go structs. They are auto-generated by <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">grit sync</code> or
                  created alongside schemas by <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">grit generate resource</code>. The
                  API response wrapper types standardize how the frontend handles all API responses.
                </p>

                <CodeBlock language="typescript" filename="packages/shared/types/user.ts" code={`export interface User {
  id: number;
  name: string;
  email: string;
  role: "admin" | "editor" | "user";
  avatar: string;
  active: boolean;
  email_verified_at: string | null;
  created_at: string;
  updated_at: string;
}

export interface LoginRequest {
  email: string;
  password: string;
}

export interface RegisterRequest {
  name: string;
  email: string;
  password: string;
}

export interface AuthResponse {
  user: User;
  tokens: {
    access_token: string;
    refresh_token: string;
    expires_at: number;
  };
}`} />

                <CodeBlock language="typescript" filename="packages/shared/types/api.ts" code={`export interface ApiResponse<T> {
  data: T;
  message?: string;
}

export interface PaginatedResponse<T> {
  data: T[];
  meta: {
    total: number;
    page: number;
    page_size: number;
    pages: number;
  };
}

export interface ApiError {
  error: {
    code: string;
    message: string;
    details?: Record<string, string>;
  };
}`} />

                <CodeBlock language="typescript" filename="packages/shared/types/upload.ts" code={`export interface Upload {
  id: number;
  filename: string;
  original_name: string;
  mime_type: string;
  size: number;
  path: string;
  url: string;
  thumbnail_url?: string;
  user_id: number;
  created_at: string;
  updated_at: string;
}`} />
              </div>

              {/* Constants */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Constants
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Shared constants define roles, API route paths, and configuration values.
                  Using these constants prevents typos and provides type-safe route references
                  across the entire frontend.
                </p>

                <CodeBlock language="typescript" filename="packages/shared/constants/index.ts" code={`export const ROLES = {
  ADMIN: "admin",
  EDITOR: "editor",
  USER: "user",
} as const;

export const API_ROUTES = {
  AUTH: {
    LOGIN: "/api/auth/login",
    REGISTER: "/api/auth/register",
    REFRESH: "/api/auth/refresh",
    LOGOUT: "/api/auth/logout",
    ME: "/api/auth/me",
    FORGOT_PASSWORD: "/api/auth/forgot-password",
    RESET_PASSWORD: "/api/auth/reset-password",
  },
  USERS: {
    LIST: "/api/users",
    GET: (id: number) => \`/api/users/\${id}\`,
    UPDATE: (id: number) => \`/api/users/\${id}\`,
    DELETE: (id: number) => \`/api/users/\${id}\`,
  },
  UPLOADS: {
    CREATE: "/api/uploads",
    LIST: "/api/uploads",
    GET: (id: number) => \`/api/uploads/\${id}\`,
    DELETE: (id: number) => \`/api/uploads/\${id}\`,
  },
  AI: {
    COMPLETE: "/api/ai/complete",
    CHAT: "/api/ai/chat",
    STREAM: "/api/ai/stream",
  },
  ADMIN: {
    JOBS_STATS: "/api/admin/jobs/stats",
    JOBS_LIST: (status: string) => \`/api/admin/jobs/\${status}\`,
    JOBS_RETRY: (id: string) => \`/api/admin/jobs/\${id}/retry\`,
    CRON_TASKS: "/api/admin/cron/tasks",
  },
  HEALTH: "/api/health",
  // grit:api-routes  (marker for code generation)
} as const;`} />
              </div>

              {/* How Schemas Are Used */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  How Schemas Work Across the Stack
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Zod schemas serve as the single source of truth for validation. The same schema validates
                  form input on the frontend and can be referenced when building Go struct tags for backend
                  validation. This ensures consistent error messages and validation rules across the stack.
                </p>

                <CodeBlock filename="validation-flow.txt" code={`┌─────────────────────────────────────────────────┐
│  Go Struct (source of truth)                    │
│  type Post struct {                             │
│    Title   string \`binding:"required"\`           │
│    Content string \`binding:"required"\`           │
│  }                                              │
└─────────────┬───────────────────────────────────┘
              │  grit sync / grit generate
              v
┌─────────────────────────────────────────────────┐
│  Zod Schema (packages/shared/schemas/post.ts)   │
│  CreatePostSchema = z.object({                  │
│    title: z.string().min(1),                    │
│    content: z.string().min(1),                  │
│  })                                             │
└─────────────┬──────────────┬────────────────────┘
              │              │
              v              v
      ┌──────────────┐  ┌──────────────┐
      │  apps/web    │  │  apps/admin  │
      │  form valid. │  │  form valid. │
      └──────────────┘  └──────────────┘`} />

                <p className="text-muted-foreground leading-relaxed mb-4">
                  Example usage in a form component:
                </p>

                <CodeBlock language="tsx" filename="using-schemas-in-forms.tsx" code={`import { CreatePostSchema, type CreatePostInput } from "@myapp/shared/schemas";
import { useCreatePost } from "@/hooks/use-posts";

function CreatePostForm() {
  const createPost = useCreatePost();
  const [form, setForm] = useState<CreatePostInput>({
    title: "",
    content: "",
    published: false,
  });

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();

    // Validate with Zod before sending to API
    const result = CreatePostSchema.safeParse(form);
    if (!result.success) {
      // Display validation errors
      const errors = result.error.flatten().fieldErrors;
      // { title: ["Title is required"], content: [...] }
      return;
    }

    // Schema validation passed -- send to Go API
    createPost.mutate(result.data);
  };
}`} />
              </div>

              {/* Export Structure */}
              <div className="mb-10">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Export Structure
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Each directory has an <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">index.ts</code> barrel file that
                  re-exports everything. The <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">// grit:schemas</code> and <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">// grit:types</code> markers
                  are used by the code generator to inject new exports when a resource is generated.
                </p>

                <CodeBlock language="typescript" filename="packages/shared/schemas/index.ts" code={`export {
  LoginSchema,
  RegisterSchema,
  UpdateUserSchema,
  ForgotPasswordSchema,
  ResetPasswordSchema,
  type LoginInput,
  type RegisterInput,
  type UpdateUserInput,
  type ForgotPasswordInput,
  type ResetPasswordInput,
} from "./user";
// grit:schemas  <-- new exports are injected here`} />

                <CodeBlock language="typescript" filename="packages/shared/types/index.ts" code={`export type {
  User, LoginRequest, RegisterRequest, AuthResponse,
} from "./user";

export type {
  ApiResponse, PaginatedResponse, ApiError,
} from "./api";

export type { Upload } from "./upload";
// grit:types  <-- new type exports are injected here`} />

                <div className="p-4 rounded-lg border border-primary/20 bg-primary/5">
                  <p className="text-sm text-foreground/80">
                    <strong>Marker-based injection:</strong> When you run <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">grit generate resource Invoice</code>,
                    the CLI finds the <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">// grit:schemas</code> marker and inserts the new
                    export line directly above it. This means you never need to manually update barrel files.
                  </p>
                </div>
              </div>
            </div>

            {/* Nav */}
            <div className="flex items-center justify-between pt-6 border-t border-border/30">
              <Button variant="ghost" size="sm" asChild className="text-muted-foreground/60 hover:text-foreground">
                <Link href="/docs/frontend/hooks" className="gap-1.5">
                  <ArrowLeft className="h-3.5 w-3.5" />
                  React Query Hooks
                </Link>
              </Button>
              <Button variant="ghost" size="sm" asChild className="text-muted-foreground/60 hover:text-foreground">
                <Link href="/docs/batteries/storage" className="gap-1.5">
                  File Storage
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
