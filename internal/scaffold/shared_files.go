package scaffold

import (
	"fmt"
	"path/filepath"
)

func writeSharedFiles(root string, opts Options) error {
	sharedRoot := filepath.Join(root, "packages", "shared")

	files := map[string]string{
		filepath.Join(sharedRoot, "package.json"):          sharedPackageJSON(opts),
		filepath.Join(sharedRoot, "tsconfig.json"):         sharedTSConfig(),
		filepath.Join(sharedRoot, "schemas", "user.ts"):    sharedUserSchema(),
		filepath.Join(sharedRoot, "schemas", "index.ts"):   sharedSchemasIndex(),
		filepath.Join(sharedRoot, "types", "user.ts"):      sharedUserTypes(),
		filepath.Join(sharedRoot, "types", "api.ts"):       sharedAPITypes(),
		filepath.Join(sharedRoot, "types", "index.ts"):     sharedTypesIndex(),
		filepath.Join(sharedRoot, "constants", "index.ts"): sharedConstants(),
		filepath.Join(sharedRoot, "types", "upload.ts"):    sharedUploadTypes(),
		filepath.Join(sharedRoot, "schemas", "blog.ts"):    sharedBlogSchema(),
		filepath.Join(sharedRoot, "types", "blog.ts"):      sharedBlogTypes(),

		// v3.31.30: FileRef — canonical shape of a stored file. Used by
		// resource schemas with :file: / :files: fields and by the
		// frontend FileField component to validate upload responses.
		filepath.Join(sharedRoot, "schemas", "file-ref.ts"): sharedFileRefSchema(),
		filepath.Join(sharedRoot, "types", "file-ref.ts"):   sharedFileRefTypes(),

		// v3.28: brand identity + theme tokens — single source of truth for
		// logo, brand name, hero copy, social links, and the 3 theme palettes
		// (atlas/aurora/pulse). Auth pages and dashboards in apps/admin and
		// apps/web import these so a rebrand is one file, not a grep + edit.
		filepath.Join(sharedRoot, "brand.config.ts"): sharedBrandConfig(opts),
		filepath.Join(sharedRoot, "themes.ts"):       sharedThemes(),
	}

	for path, content := range files {
		if err := writeFile(path, content); err != nil {
			return fmt.Errorf("writing %s: %w", path, err)
		}
	}

	return nil
}

func sharedPackageJSON(opts Options) string {
	_ = opts
	return `{
  "name": "@repo/shared",
  "version": "0.1.0",
  "private": true,
  "main": "./index.ts",
  "types": "./index.ts",
  "exports": {
    "./schemas": "./schemas/index.ts",
    "./types": "./types/index.ts",
    "./constants": "./constants/index.ts",
    "./brand": "./brand.config.ts",
    "./themes": "./themes.ts"
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
  macAddress: z.string().optional(), // optional — passed by client if available
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
export { FileRefSchema, type FileRef } from "./file-ref";
// grit:schemas
`
}

func sharedUserTypes() string {
	return `export interface User {
  id: string;
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
  ip_address: string;
  mac_address: string;
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
  mac_address?: string; // optional — provided by client if obtainable
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

// apiErrorMessage extracts a human-readable message from any error
// that came out of the API client (axios or fetch). The standard error
// envelope is { error: { code, message, details? } }; this walks the
// usual axios.error.response.data.error.message chain plus a few
// fallbacks so toast.error(apiErrorMessage(err)) is always meaningful.
//
// Use it everywhere instead of err.message directly:
//
//   import { apiErrorMessage } from "@grit/shared/types/api";
//   ...
//   onError: (err) => toast.error(apiErrorMessage(err))
export function apiErrorMessage(err: unknown, fallback = "Something went wrong"): string {
  if (!err) return fallback;
  if (typeof err === "string") return err;

  // Standard envelope — preferred path
  const e = err as {
    response?: { data?: { error?: { message?: string; code?: string } } };
    message?: string;
  };
  const envMsg = e.response?.data?.error?.message;
  if (envMsg) return envMsg;

  // Axios error.message ("Network Error", "timeout of 5000ms exceeded")
  if (e.message) return e.message;

  // Last resort
  return fallback;
}

// apiErrorCode returns the standard envelope code (VALIDATION_ERROR,
// NOT_FOUND, ...) when present. Useful for branching on specific codes
// (e.g. "VERSION_CONFLICT" -> open conflict dialog).
export function apiErrorCode(err: unknown): string | null {
  if (!err) return null;
  const e = err as { response?: { data?: { error?: { code?: string } } } };
  return e.response?.data?.error?.code ?? null;
}

// apiErrorFields surfaces per-field validation errors when the API
// returns details: { fieldName: "error message" }. Used to highlight
// individual inputs in a form.
export function apiErrorFields(err: unknown): Record<string, string> {
  if (!err) return {};
  const e = err as { response?: { data?: { error?: { details?: Record<string, string> } } } };
  return e.response?.data?.error?.details ?? {};
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

export {
  apiErrorMessage,
  apiErrorCode,
  apiErrorFields,
} from "./api";

export type { Upload } from "./upload";
export type { Blog } from "./blog";
export type { FileRef } from "./file-ref";
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
    GET: (id: string) => ` + "`" + `/api/users/${id}` + "`" + `,
    UPDATE: (id: string) => ` + "`" + `/api/users/${id}` + "`" + `,
    DELETE: (id: string) => ` + "`" + `/api/users/${id}` + "`" + `,
  },
  UPLOADS: {
    CREATE: "/api/uploads",
    LIST: "/api/uploads",
    GET: (id: string) => ` + "`" + `/api/uploads/${id}` + "`" + `,
    DELETE: (id: string) => ` + "`" + `/api/uploads/${id}` + "`" + `,
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
    UPDATE: (id: string) => ` + "`" + `/api/admin/blogs/${id}` + "`" + `,
    DELETE: (id: string) => ` + "`" + `/api/admin/blogs/${id}` + "`" + `,
  },
  HEALTH: "/api/health",
  // grit:api-routes
} as const;
`
}

func sharedUploadTypes() string {
	return `export interface Upload {
  id: string;
  filename: string;
  original_name: string;
  mime_type: string;
  size: number;
  path: string;
  url: string;
  thumbnail_url?: string;
  user_id: string;
  created_at: string;
  updated_at: string;
}
`
}

func sharedBlogSchema() string {
	return `import { z } from "zod";

export const BlogSchema = z.object({
  id: z.string(),
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
  id: string;
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

// v3.31.30 — FileRef shared shape.
//
// Every resource with a :file: / :files: field stores this exact JSON
// in its column. The frontend uploads via POST /api/uploads and gets a
// FileRef back, then submits the parent form with the FileRef
// embedded — no multipart wrangling at the resource endpoint.

func sharedFileRefSchema() string {
	return `import { z } from "zod";

// FileRef — canonical shape of a stored file. The API's POST /api/uploads
// returns this; resource forms embed it in their submit body.
//
// Fields width / height / duration / thumbnail_url are populated by the
// server when the source format makes them cheap to compute (images get
// dimensions, audio gets duration). They're optional because not every
// upload has them — a PDF has no width, a CSV has no thumbnail.
export const FileRefSchema = z.object({
  url: z.string().url(),
  key: z.string().min(1),
  name: z.string().min(1),
  mime: z.string().min(1),
  size: z.number().int().nonnegative(),
  width: z.number().int().positive().optional(),
  height: z.number().int().positive().optional(),
  duration: z.number().int().nonnegative().optional(),
  thumbnail_url: z.string().url().optional(),
});

export type FileRef = z.infer<typeof FileRefSchema>;
`
}

func sharedFileRefTypes() string {
	return `// FileRef — re-export of the Zod-inferred type for code that only
// needs the type, not the schema. The schema lives in
// schemas/file-ref.ts; importing the type from here avoids pulling in
// Zod just to get a type definition.

export type FileRef = {
  url: string;
  key: string;
  name: string;
  mime: string;
  size: number;
  width?: number;
  height?: number;
  duration?: number;
  thumbnail_url?: string;
};
`
}
