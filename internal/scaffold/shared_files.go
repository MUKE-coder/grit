package scaffold

import (
	"fmt"
	"path/filepath"
)

func writeSharedFiles(root string, opts Options) error {
	sharedRoot := filepath.Join(root, "packages", "shared")

	files := map[string]string{
		filepath.Join(sharedRoot, "package.json"):        sharedPackageJSON(opts),
		filepath.Join(sharedRoot, "tsconfig.json"):       sharedTSConfig(),
		filepath.Join(sharedRoot, "schemas", "user.ts"):  sharedUserSchema(),
		filepath.Join(sharedRoot, "schemas", "index.ts"): sharedSchemasIndex(),
		filepath.Join(sharedRoot, "types", "user.ts"):    sharedUserTypes(),
		filepath.Join(sharedRoot, "types", "api.ts"):     sharedAPITypes(),
		filepath.Join(sharedRoot, "types", "index.ts"):   sharedTypesIndex(),
		filepath.Join(sharedRoot, "constants", "index.ts"): sharedConstants(),
		filepath.Join(sharedRoot, "types", "upload.ts"):    sharedUploadTypes(),
		filepath.Join(sharedRoot, "schemas", "blog.ts"):    sharedBlogSchema(),
		filepath.Join(sharedRoot, "types", "blog.ts"):      sharedBlogTypes(),
	}

	for path, content := range files {
		if err := writeFile(path, content); err != nil {
			return fmt.Errorf("writing %s: %w", path, err)
		}
	}

	return nil
}

func sharedPackageJSON(opts Options) string {
	return `{
  "name": "@repo/shared",
  "version": "0.1.0",
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
  },
  "devDependencies": {
    "typescript": "^5.3.0"
  }
}
`
}

func sharedTSConfig() string {
	return `{
  "compilerOptions": {
    "target": "ES2020",
    "module": "ESNext",
    "moduleResolution": "bundler",
    "strict": true,
    "esModuleInterop": true,
    "skipLibCheck": true,
    "declaration": true,
    "declarationMap": true,
    "sourceMap": true,
    "outDir": "dist"
  },
  "include": ["**/*.ts"],
  "exclude": ["node_modules", "dist"]
}
`
}

func sharedUserSchema() string {
	return `import { z } from "zod";

export const LoginSchema = z.object({
  email: z.string().email("Please enter a valid email"),
  password: z.string().min(1, "Password is required"),
});

export const RegisterSchema = z.object({
  firstName: z.string().min(2, "First name must be at least 2 characters"),
  lastName: z.string().min(2, "Last name must be at least 2 characters"),
  email: z.string().email("Please enter a valid email"),
  password: z.string().min(8, "Password must be at least 8 characters"),
  confirmPassword: z.string(),
}).refine((data) => data.password === data.confirmPassword, {
  message: "Passwords do not match",
  path: ["confirmPassword"],
});

export const UpdateUserSchema = z.object({
  firstName: z.string().min(2).optional(),
  lastName: z.string().min(2).optional(),
  email: z.string().email().optional(),
  role: z.enum(["ADMIN", "EDITOR", "USER"]).optional(), // grit:role-enum
  jobTitle: z.string().optional(),
  bio: z.string().optional(),
  active: z.boolean().optional(),
  provider: z.string().optional(),
});

export const ForgotPasswordSchema = z.object({
  email: z.string().email("Please enter a valid email"),
});

export const ResetPasswordSchema = z.object({
  token: z.string().min(1, "Token is required"),
  password: z.string().min(8, "Password must be at least 8 characters"),
});

export type LoginInput = z.infer<typeof LoginSchema>;
export type RegisterInput = z.infer<typeof RegisterSchema>;
export type UpdateUserInput = z.infer<typeof UpdateUserSchema>;
export type ForgotPasswordInput = z.infer<typeof ForgotPasswordSchema>;
export type ResetPasswordInput = z.infer<typeof ResetPasswordSchema>;
`
}

func sharedSchemasIndex() string {
	return `export {
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
export {
  BlogSchema,
  CreateBlogSchema,
  UpdateBlogSchema,
  type CreateBlogInput,
  type UpdateBlogInput,
} from "./blog";
// grit:schemas
`
}

func sharedUserTypes() string {
	return `export interface User {
  id: number;
  first_name: string;
  last_name: string;
  email: string;
  role: "ADMIN" | "EDITOR" | "USER"; // grit:role-union
  avatar: string;
  job_title: string;
  bio: string;
  active: boolean;
  provider: string;
  email_verified_at: string | null;
  created_at: string;
  updated_at: string;
}

export interface LoginRequest {
  email: string;
  password: string;
}

export interface RegisterRequest {
  first_name: string;
  last_name: string;
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
}
`
}

func sharedAPITypes() string {
	return `export interface ApiResponse<T> {
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
}
`
}

func sharedTypesIndex() string {
	return `export type {
  User,
  LoginRequest,
  RegisterRequest,
  AuthResponse,
} from "./user";

export type {
  ApiResponse,
  PaginatedResponse,
  ApiError,
} from "./api";

export type { Upload } from "./upload";
export type { Blog } from "./blog";
// grit:types
`
}

func sharedConstants() string {
	return `export const ROLES = {
  ADMIN: "ADMIN",
  EDITOR: "EDITOR",
  USER: "USER",
  // grit:role-constants
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
    OAUTH: {
      GOOGLE: "/api/auth/oauth/google",
      GITHUB: "/api/auth/oauth/github",
    },
  },
  USERS: {
    LIST: "/api/users",
    GET: (id: number) => ` + "`" + `/api/users/${id}` + "`" + `,
    UPDATE: (id: number) => ` + "`" + `/api/users/${id}` + "`" + `,
    DELETE: (id: number) => ` + "`" + `/api/users/${id}` + "`" + `,
  },
  UPLOADS: {
    CREATE: "/api/uploads",
    LIST: "/api/uploads",
    GET: (id: number) => ` + "`" + `/api/uploads/${id}` + "`" + `,
    DELETE: (id: number) => ` + "`" + `/api/uploads/${id}` + "`" + `,
  },
  AI: {
    COMPLETE: "/api/ai/complete",
    CHAT: "/api/ai/chat",
    STREAM: "/api/ai/stream",
  },
  ADMIN: {
    JOBS_STATS: "/api/admin/jobs/stats",
    JOBS_LIST: (status: string) => ` + "`" + `/api/admin/jobs/${status}` + "`" + `,
    JOBS_RETRY: (id: string) => ` + "`" + `/api/admin/jobs/${id}/retry` + "`" + `,
    JOBS_CLEAR: (queue: string) => ` + "`" + `/api/admin/jobs/queue/${queue}` + "`" + `,
    CRON_TASKS: "/api/admin/cron/tasks",
  },
  PROFILE: {
    GET: "/api/profile",
    UPDATE: "/api/profile",
    DELETE: "/api/profile",
  },
  BLOGS: {
    LIST: "/api/blogs",
    GET: (slug: string) => ` + "`" + `/api/blogs/${slug}` + "`" + `,
    ADMIN_LIST: "/api/admin/blogs",
    CREATE: "/api/admin/blogs",
    UPDATE: (id: number) => ` + "`" + `/api/admin/blogs/${id}` + "`" + `,
    DELETE: (id: number) => ` + "`" + `/api/admin/blogs/${id}` + "`" + `,
  },
  HEALTH: "/api/health",
  // grit:api-routes
} as const;
`
}

func sharedUploadTypes() string {
	return `export interface Upload {
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
}
`
}

func sharedBlogSchema() string {
	return `import { z } from "zod";

export const BlogSchema = z.object({
  id: z.number(),
  title: z.string(),
  slug: z.string(),
  content: z.string(),
  image: z.string().nullable(),
  excerpt: z.string().nullable(),
  published: z.boolean(),
  published_at: z.string().nullable(),
  created_at: z.string(),
  updated_at: z.string(),
});

export const CreateBlogSchema = z.object({
  title: z.string().min(1, "Title is required"),
  content: z.string().optional(),
  image: z.string().optional(),
  excerpt: z.string().optional(),
  published: z.boolean().optional(),
});

export const UpdateBlogSchema = CreateBlogSchema.partial();

export type CreateBlogInput = z.infer<typeof CreateBlogSchema>;
export type UpdateBlogInput = z.infer<typeof UpdateBlogSchema>;
`
}

func sharedBlogTypes() string {
	return `export interface Blog {
  id: number;
  title: string;
  slug: string;
  content: string;
  image: string | null;
  excerpt: string | null;
  published: boolean;
  published_at: string | null;
  created_at: string;
  updated_at: string;
}
`
}
