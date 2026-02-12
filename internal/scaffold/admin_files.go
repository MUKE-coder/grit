package scaffold

import (
	"fmt"
	"path/filepath"
)

func writeAdminFiles(root string, opts Options) error {
	adminRoot := filepath.Join(root, "apps", "admin")

	files := map[string]string{
		// Config files
		filepath.Join(adminRoot, "package.json"):     adminPackageJSON(opts),
		filepath.Join(adminRoot, "next.config.ts"):   adminNextConfig(),
		filepath.Join(adminRoot, "tailwind.config.ts"): adminTailwindConfig(),
		filepath.Join(adminRoot, "postcss.config.js"): adminPostCSSConfig(),
		filepath.Join(adminRoot, "tsconfig.json"):    adminTSConfig(),
		filepath.Join(adminRoot, "app", "globals.css"): adminGlobalCSS(),
		filepath.Join(adminRoot, "app", "layout.tsx"): adminRootLayout(opts),

		// Root redirect page
		filepath.Join(adminRoot, "app", "page.tsx"): adminRedirectPage(),

		// Auth pages — (auth) route group
		filepath.Join(adminRoot, "app", "(auth)", "login", "page.tsx"):           adminLoginPage(),
		filepath.Join(adminRoot, "app", "(auth)", "sign-up", "page.tsx"):         adminSignUpPage(),
		filepath.Join(adminRoot, "app", "(auth)", "forgot-password", "page.tsx"): adminForgotPasswordPage(),

		// Dashboard route group layout
		filepath.Join(adminRoot, "app", "(dashboard)", "layout.tsx"): adminDashboardLayout(),

		// Lib
		filepath.Join(adminRoot, "lib", "api-client.ts"):  adminAPIClient(),
		filepath.Join(adminRoot, "lib", "query-client.ts"): adminQueryClient(),
		filepath.Join(adminRoot, "lib", "utils.ts"):        adminUtils(),
		filepath.Join(adminRoot, "lib", "resource.ts"):     adminResourceTypes(),
		filepath.Join(adminRoot, "lib", "icons.ts"):        adminIconMap(),
		filepath.Join(adminRoot, "lib", "formatters.ts"):   adminFormatters(),

		// Shared components
		filepath.Join(adminRoot, "components", "shared", "providers.tsx"):      adminProviders(),
		filepath.Join(adminRoot, "components", "shared", "theme-provider.tsx"): adminThemeProvider(),

		// Layout components
		filepath.Join(adminRoot, "components", "layout", "admin-layout.tsx"): adminLayoutComponent(),
		filepath.Join(adminRoot, "components", "layout", "sidebar.tsx"):      adminSidebar(),
		filepath.Join(adminRoot, "components", "layout", "navbar.tsx"):       adminNavbar(),

		// Table components
		filepath.Join(adminRoot, "components", "tables", "data-table.tsx"):       adminDataTable(),
		filepath.Join(adminRoot, "components", "tables", "column-header.tsx"):    adminColumnHeader(),
		filepath.Join(adminRoot, "components", "tables", "cell-renderers.tsx"):   adminCellRenderers(),
		filepath.Join(adminRoot, "components", "tables", "table-filters.tsx"):    adminTableFilters(),
		filepath.Join(adminRoot, "components", "tables", "table-toolbar.tsx"):    adminTableToolbar(),
		filepath.Join(adminRoot, "components", "tables", "table-pagination.tsx"): adminTablePagination(),
		filepath.Join(adminRoot, "components", "tables", "table-skeleton.tsx"):   adminTableSkeleton(),
		filepath.Join(adminRoot, "components", "tables", "table-empty-state.tsx"): adminTableEmptyState(),

		// Form components
		filepath.Join(adminRoot, "components", "forms", "form-builder.tsx"): adminFormBuilder(),
		filepath.Join(adminRoot, "components", "forms", "form-modal.tsx"):   adminFormModal(),
		filepath.Join(adminRoot, "components", "forms", "fields", "text-field.tsx"):     adminTextField(),
		filepath.Join(adminRoot, "components", "forms", "fields", "textarea-field.tsx"): adminTextareaField(),
		filepath.Join(adminRoot, "components", "forms", "fields", "number-field.tsx"):   adminNumberField(),
		filepath.Join(adminRoot, "components", "forms", "fields", "select-field.tsx"):   adminSelectField(),
		filepath.Join(adminRoot, "components", "forms", "fields", "date-field.tsx"):     adminDateField(),
		filepath.Join(adminRoot, "components", "forms", "fields", "toggle-field.tsx"):   adminToggleField(),
		filepath.Join(adminRoot, "components", "forms", "fields", "checkbox-field.tsx"): adminCheckboxField(),
		filepath.Join(adminRoot, "components", "forms", "fields", "radio-field.tsx"):    adminRadioField(),
		filepath.Join(adminRoot, "components", "forms", "fields", "image-field.tsx"):    adminImageField(),
		filepath.Join(adminRoot, "components", "forms", "fields", "images-field.tsx"):   adminImagesField(),
		filepath.Join(adminRoot, "components", "forms", "fields", "video-field.tsx"):    adminVideoField(),
		filepath.Join(adminRoot, "components", "forms", "fields", "videos-field.tsx"):   adminVideosField(),
		filepath.Join(adminRoot, "components", "forms", "fields", "file-field.tsx"):     adminFileField(),
		filepath.Join(adminRoot, "components", "forms", "fields", "files-field.tsx"):    adminFilesField(),

		// UI components
		filepath.Join(adminRoot, "components", "ui", "dropzone.tsx"):      adminDropzone(),
		filepath.Join(adminRoot, "components", "ui", "confirm-modal.tsx"): adminConfirmModal(),

		// Widget components
		filepath.Join(adminRoot, "components", "widgets", "stats-card.tsx"):      adminStatsCard(),
		filepath.Join(adminRoot, "components", "widgets", "chart-widget.tsx"):    adminChartWidget(),
		filepath.Join(adminRoot, "components", "widgets", "activity-widget.tsx"): adminActivityWidget(),
		filepath.Join(adminRoot, "components", "widgets", "widget-grid.tsx"):     adminWidgetGrid(),

		// Resource components
		filepath.Join(adminRoot, "components", "resource", "resource-page.tsx"): adminResourcePage(),
		filepath.Join(adminRoot, "components", "resource", "view-modal.tsx"):    adminViewModal(),

		// Resource definitions
		filepath.Join(adminRoot, "resources", "index.ts"): adminResourceRegistry(),
		filepath.Join(adminRoot, "resources", "users.ts"): adminUsersResource(),

		// Hooks
		filepath.Join(adminRoot, "hooks", "use-auth.ts"):     adminUseAuth(),
		filepath.Join(adminRoot, "hooks", "use-resource.ts"): adminUseResource(),
		filepath.Join(adminRoot, "hooks", "use-system.ts"):   adminUseSystem(),

		// Dashboard pages — (dashboard) route group
		filepath.Join(adminRoot, "app", "(dashboard)", "dashboard", "page.tsx"):          adminDashboardPage(),
		filepath.Join(adminRoot, "app", "(dashboard)", "resources", "users", "page.tsx"): adminUsersPage(),

		// System pages — under (dashboard) route group
		filepath.Join(adminRoot, "app", "(dashboard)", "system", "jobs", "page.tsx"):  adminJobsPage(),
		filepath.Join(adminRoot, "app", "(dashboard)", "system", "files", "page.tsx"): adminFilesPage(),
		filepath.Join(adminRoot, "app", "(dashboard)", "system", "cron", "page.tsx"):  adminCronPage(),
		filepath.Join(adminRoot, "app", "(dashboard)", "system", "mail", "page.tsx"):  adminMailPage(),
	}

	for path, content := range files {
		if err := writeFile(path, content); err != nil {
			return fmt.Errorf("writing %s: %w", path, err)
		}
	}

	return nil
}

func adminPackageJSON(opts Options) string {
	return fmt.Sprintf(`{
  "name": "@%s/admin",
  "version": "0.1.0",
  "private": true,
  "scripts": {
    "dev": "next dev --port 3001",
    "build": "next build",
    "start": "next start",
    "lint": "next lint"
  },
  "dependencies": {
    "@hookform/resolvers": "^3.3.0",
    "@tanstack/react-query": "^5.17.0",
    "axios": "^1.6.0",
    "class-variance-authority": "^0.7.0",
    "clsx": "^2.1.0",
    "js-cookie": "^3.0.5",
    "lucide-react": "^0.303.0",
    "next": "^16.1.6",
    "react": "^19.0.0",
    "react-dom": "^19.0.0",
    "react-dropzone": "^14.2.0",
    "react-hook-form": "^7.49.0",
    "recharts": "^2.12.0",
    "sonner": "^1.3.0",
    "tailwind-merge": "^2.2.0",
    "tailwindcss-animate": "^1.0.7",
    "zod": "^3.22.0"
  },
  "devDependencies": {
    "@types/js-cookie": "^3.0.6",
    "@types/node": "^20.0.0",
    "@types/react": "^19.0.0",
    "@types/react-dom": "^19.0.0",
    "autoprefixer": "^10.4.0",
    "postcss": "^8.4.0",
    "tailwindcss": "^3.4.0",
    "typescript": "^5.3.0"
  }
}
`, opts.ProjectName)
}

func adminNextConfig() string {
	return `import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  output: "standalone",
  reactStrictMode: true,
};

export default nextConfig;
`
}

func adminTailwindConfig() string {
	return `import type { Config } from "tailwindcss";

const config: Config = {
  darkMode: "class",
  content: [
    "./app/**/*.{ts,tsx}",
    "./components/**/*.{ts,tsx}",
    "./lib/**/*.{ts,tsx}",
  ],
  theme: {
    extend: {
      colors: {
        background: "var(--bg-primary)",
        "bg-secondary": "var(--bg-secondary)",
        "bg-tertiary": "var(--bg-tertiary)",
        "bg-elevated": "var(--bg-elevated)",
        "bg-hover": "var(--bg-hover)",
        border: "var(--border)",
        foreground: "var(--text-primary)",
        "text-secondary": "var(--text-secondary)",
        "text-muted": "var(--text-muted)",
        accent: {
          DEFAULT: "var(--accent)",
          hover: "var(--accent-hover)",
        },
        success: "var(--success)",
        danger: "var(--danger)",
        warning: "var(--warning)",
        info: "var(--info)",
      },
      fontFamily: {
        sans: ["var(--font-dm-sans)", "system-ui", "sans-serif"],
        mono: ["var(--font-jetbrains-mono)", "ui-monospace", "monospace"],
      },
    },
  },
  plugins: [require("tailwindcss-animate")],
};

export default config;
`
}

func adminPostCSSConfig() string {
	return `module.exports = {
  plugins: {
    tailwindcss: {},
    autoprefixer: {},
  },
};
`
}

func adminTSConfig() string {
	return `{
  "compilerOptions": {
    "target": "ES2017",
    "lib": ["dom", "dom.iterable", "esnext"],
    "allowJs": true,
    "skipLibCheck": true,
    "strict": true,
    "noEmit": true,
    "esModuleInterop": true,
    "module": "esnext",
    "moduleResolution": "bundler",
    "resolveJsonModule": true,
    "isolatedModules": true,
    "jsx": "preserve",
    "incremental": true,
    "plugins": [{ "name": "next" }],
    "paths": {
      "@/*": ["./*"]
    }
  },
  "include": ["next-env.d.ts", "**/*.ts", "**/*.tsx", ".next/types/**/*.ts"],
  "exclude": ["node_modules"]
}
`
}

func adminGlobalCSS() string {
	return `@tailwind base;
@tailwind components;
@tailwind utilities;

/* Dark theme (default) */
:root {
  --bg-primary: #0a0a0f;
  --bg-secondary: #111118;
  --bg-tertiary: #1a1a24;
  --bg-elevated: #22222e;
  --bg-hover: #2a2a38;
  --border: #2a2a3a;
  --text-primary: #e8e8f0;
  --text-secondary: #9090a8;
  --text-muted: #606078;
  --accent: #6c5ce7;
  --accent-hover: #7c6cf7;
  --success: #00b894;
  --danger: #ff6b6b;
  --warning: #fdcb6e;
  --info: #74b9ff;
}

/* Light theme */
:root.light {
  --bg-primary: #f8f9fc;
  --bg-secondary: #ffffff;
  --bg-tertiary: #f1f3f8;
  --bg-elevated: #ffffff;
  --bg-hover: #e8eaf0;
  --border: #d8dbe5;
  --text-primary: #1a1a2e;
  --text-secondary: #555570;
  --text-muted: #8888a0;
  --accent: #6c5ce7;
  --accent-hover: #5a4bd6;
  --success: #00b894;
  --danger: #ff6b6b;
  --warning: #e5a800;
  --info: #3b8beb;
}

body {
  background-color: var(--bg-primary);
  color: var(--text-primary);
  font-family: var(--font-dm-sans), system-ui, sans-serif;
}

* {
  border-color: var(--border);
}

::-webkit-scrollbar {
  width: 6px;
  height: 6px;
}

::-webkit-scrollbar-track {
  background: var(--bg-primary);
}

::-webkit-scrollbar-thumb {
  background: var(--bg-hover);
  border-radius: 3px;
}
`
}

func adminRootLayout(opts Options) string {
	return fmt.Sprintf(`import type { Metadata } from "next";
import { DM_Sans, JetBrains_Mono } from "next/font/google";
import "./globals.css";
import { Providers } from "@/components/shared/providers";

const dmSans = DM_Sans({
  subsets: ["latin"],
  variable: "--font-dm-sans",
  weight: ["400", "500", "600", "700"],
});

const jetbrainsMono = JetBrains_Mono({
  subsets: ["latin"],
  variable: "--font-jetbrains-mono",
  weight: ["400", "500", "600"],
});

export const metadata: Metadata = {
  title: "%s Admin",
  description: "Admin panel — Built with Grit",
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="en" className="dark">
      <body className={` + "`" + `${dmSans.variable} ${jetbrainsMono.variable} min-h-screen bg-background font-sans antialiased` + "`" + `}>
        <Providers>{children}</Providers>
      </body>
    </html>
  );
}
`, opts.ProjectName)
}

func adminProviders() string {
	return `"use client";

import { QueryClientProvider } from "@tanstack/react-query";
import { Toaster } from "sonner";
import { queryClient } from "@/lib/query-client";
import { ThemeProvider } from "./theme-provider";

export function Providers({ children }: { children: React.ReactNode }) {
  return (
    <ThemeProvider>
      <QueryClientProvider client={queryClient}>
        {children}
        <Toaster
          position="bottom-right"
          toastOptions={{
            style: {
              background: "var(--bg-elevated)",
              border: "1px solid var(--border)",
              color: "var(--text-primary)",
            },
          }}
        />
      </QueryClientProvider>
    </ThemeProvider>
  );
}
`
}

func adminUseAuth() string {
	return `import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useRouter } from "next/navigation";
import Cookies from "js-cookie";
import { apiClient } from "@/lib/api-client";

interface User {
  id: number;
  name: string;
  email: string;
  role: string;
  avatar: string;
  active: boolean;
}

interface AuthResponse {
  data: {
    user: User;
    tokens: {
      access_token: string;
      refresh_token: string;
      expires_at: number;
    };
  };
}

function storeTokens(tokens: { access_token: string; refresh_token: string }) {
  Cookies.set("access_token", tokens.access_token, { expires: 1 });
  Cookies.set("refresh_token", tokens.refresh_token, { expires: 7 });
}

function clearTokens() {
  Cookies.remove("access_token");
  Cookies.remove("refresh_token");
}

export function useMe() {
  return useQuery<User>({
    queryKey: ["me"],
    queryFn: async () => {
      const { data } = await apiClient.get("/api/auth/me");
      return data.data;
    },
    retry: false,
    staleTime: 10 * 60 * 1000,
  });
}

export function useLogin() {
  const router = useRouter();
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (credentials: { email: string; password: string }) => {
      const { data } = await apiClient.post<AuthResponse>(
        "/api/auth/login",
        credentials
      );
      return data;
    },
    onSuccess: (data) => {
      storeTokens(data.data.tokens);
      queryClient.setQueryData(["me"], data.data.user);
      router.push("/dashboard");
    },
  });
}

export function useRegister() {
  const router = useRouter();
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (data: {
      name: string;
      email: string;
      password: string;
    }) => {
      const { data: response } = await apiClient.post<AuthResponse>(
        "/api/auth/register",
        data
      );
      return response;
    },
    onSuccess: (data) => {
      storeTokens(data.data.tokens);
      queryClient.setQueryData(["me"], data.data.user);
      router.push("/dashboard");
    },
  });
}

export function useLogout() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async () => {
      try {
        await apiClient.post("/api/auth/logout");
      } catch {
        // Ignore
      }
    },
    onSettled: () => {
      clearTokens();
      queryClient.clear();
      if (typeof window !== "undefined") {
        window.location.href = "/login";
      }
    },
  });
}
`
}

func adminAPIClient() string {
	return `import axios from "axios";
import Cookies from "js-cookie";

const API_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";

export const apiClient = axios.create({
  baseURL: API_URL,
  headers: {
    "Content-Type": "application/json",
  },
});

apiClient.interceptors.request.use((config) => {
  const token = Cookies.get("access_token");
  if (token) {
    config.headers.Authorization = ` + "`" + `Bearer ${token}` + "`" + `;
  }
  return config;
});

let isRefreshing = false;
let failedQueue: Array<{
  resolve: (value: unknown) => void;
  reject: (reason: unknown) => void;
}> = [];

const processQueue = (error: unknown) => {
  failedQueue.forEach((promise) => {
    if (error) {
      promise.reject(error);
    } else {
      promise.resolve(undefined);
    }
  });
  failedQueue = [];
};

apiClient.interceptors.response.use(
  (response) => response,
  async (error) => {
    const originalRequest = error.config;

    if (error.response?.status === 401 && !originalRequest._retry) {
      if (isRefreshing) {
        return new Promise((resolve, reject) => {
          failedQueue.push({ resolve, reject });
        }).then(() => apiClient(originalRequest));
      }

      originalRequest._retry = true;
      isRefreshing = true;

      try {
        const refreshToken = Cookies.get("refresh_token");
        if (!refreshToken) throw new Error("No refresh token");

        const { data } = await axios.post(` + "`" + `${API_URL}/api/auth/refresh` + "`" + `, {
          refresh_token: refreshToken,
        });

        const { access_token, refresh_token: newRefreshToken } = data.data.tokens;
        Cookies.set("access_token", access_token, { expires: 1 });
        Cookies.set("refresh_token", newRefreshToken, { expires: 7 });

        processQueue(null);
        return apiClient(originalRequest);
      } catch (refreshError) {
        processQueue(refreshError);
        Cookies.remove("access_token");
        Cookies.remove("refresh_token");
        if (typeof window !== "undefined") {
          window.location.href = "/login";
        }
        return Promise.reject(refreshError);
      } finally {
        isRefreshing = false;
      }
    }

    return Promise.reject(error);
  }
);
`
}

func adminQueryClient() string {
	return `import { QueryClient } from "@tanstack/react-query";

export const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      staleTime: 5 * 60 * 1000,
      retry: 1,
      refetchOnWindowFocus: false,
    },
    mutations: {
      retry: 0,
    },
  },
});
`
}

func adminUtils() string {
	return `import { clsx, type ClassValue } from "clsx";
import { twMerge } from "tailwind-merge";

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}
`
}

// adminRedirectPage returns the root page that redirects to /dashboard or /login.
func adminRedirectPage() string {
	return `"use client";

import { useEffect } from "react";
import { useRouter } from "next/navigation";
import Cookies from "js-cookie";

export default function RootPage() {
  const router = useRouter();

  useEffect(() => {
    const token = Cookies.get("access_token");
    if (token) {
      router.replace("/dashboard");
    } else {
      router.replace("/login");
    }
  }, [router]);

  return (
    <div className="flex min-h-screen items-center justify-center">
      <div className="h-8 w-8 animate-spin rounded-full border-2 border-accent border-t-transparent" />
    </div>
  );
}
`
}

// adminDashboardLayout returns the (dashboard) route group layout with AdminLayout.
func adminDashboardLayout() string {
	return `"use client";

import { AdminLayout } from "@/components/layout/admin-layout";

export default function DashboardGroupLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return <AdminLayout>{children}</AdminLayout>;
}
`
}

// adminLoginPage returns the split-screen login page.
func adminLoginPage() string {
	return `"use client";

import { useState } from "react";
import Link from "next/link";
import { Eye, EyeOff } from "@/lib/icons";
import { useLogin } from "@/hooks/use-auth";

export default function LoginPage() {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [showPassword, setShowPassword] = useState(false);
  const { mutate: login, isPending, error } = useLogin();

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    login({ email, password });
  };

  return (
    <div className="flex min-h-screen">
      {/* Left panel — branding */}
      <div className="hidden lg:flex lg:w-1/2 flex-col justify-between bg-gradient-to-br from-accent/20 via-bg-secondary to-bg-primary p-12">
        <div>
          <span className="text-2xl font-bold text-accent">G</span>
          <span className="text-2xl font-bold text-accent">rit</span>
          <span className="ml-2 rounded-md bg-accent/10 px-2 py-0.5 text-xs font-medium text-accent">
            Admin
          </span>
        </div>
        <div className="space-y-4">
          <h1 className="text-4xl font-bold text-foreground leading-tight">
            Manage everything<br />in one place.
          </h1>
          <p className="text-text-secondary text-lg max-w-md">
            The admin dashboard for your Grit application. Monitor, manage, and control your entire platform.
          </p>
        </div>
        <p className="text-text-muted text-sm">Built with Grit — Go + React framework</p>
      </div>

      {/* Right panel — form */}
      <div className="flex flex-1 items-center justify-center px-6 py-12">
        <div className="w-full max-w-md space-y-8">
          <div className="lg:hidden text-center mb-8">
            <span className="text-2xl font-bold text-accent">Grit</span>
            <span className="ml-2 rounded-md bg-accent/10 px-2 py-0.5 text-xs font-medium text-accent">
              Admin
            </span>
          </div>

          <div>
            <h2 className="text-2xl font-bold text-foreground">Welcome back</h2>
            <p className="mt-2 text-text-secondary">Sign in to your admin account</p>
          </div>

          <form onSubmit={handleSubmit} className="space-y-5">
            {error && (
              <div className="rounded-lg bg-danger/10 border border-danger/20 px-4 py-3 text-sm text-danger">
                {(error as unknown as { response?: { data?: { error?: { message?: string } } } })?.response?.data?.error?.message || "Invalid credentials"}
              </div>
            )}

            <div className="space-y-2">
              <label htmlFor="email" className="block text-sm font-medium text-text-secondary">
                Email
              </label>
              <input
                id="email"
                type="email"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                className="w-full rounded-lg border border-border bg-bg-tertiary px-4 py-3 text-foreground placeholder:text-text-muted focus:border-accent focus:outline-none focus:ring-1 focus:ring-accent transition-colors"
                placeholder="you@example.com"
                required
                autoFocus
              />
            </div>

            <div className="space-y-2">
              <label htmlFor="password" className="block text-sm font-medium text-text-secondary">
                Password
              </label>
              <div className="relative">
                <input
                  id="password"
                  type={showPassword ? "text" : "password"}
                  value={password}
                  onChange={(e) => setPassword(e.target.value)}
                  className="w-full rounded-lg border border-border bg-bg-tertiary px-4 py-3 pr-12 text-foreground placeholder:text-text-muted focus:border-accent focus:outline-none focus:ring-1 focus:ring-accent transition-colors"
                  placeholder="Enter your password"
                  required
                />
                <button
                  type="button"
                  onClick={() => setShowPassword(!showPassword)}
                  className="absolute right-3 top-1/2 -translate-y-1/2 text-text-muted hover:text-text-secondary transition-colors"
                >
                  {showPassword ? <EyeOff className="h-5 w-5" /> : <Eye className="h-5 w-5" />}
                </button>
              </div>
            </div>

            <div className="flex items-center justify-between">
              <label className="flex items-center gap-2 text-sm text-text-secondary cursor-pointer">
                <input type="checkbox" className="h-4 w-4 rounded border-border bg-bg-tertiary accent-accent" />
                Remember me
              </label>
              <Link href="/forgot-password" className="text-sm text-accent hover:text-accent-hover transition-colors">
                Forgot password?
              </Link>
            </div>

            <button
              type="submit"
              disabled={isPending}
              className="w-full rounded-lg bg-accent py-3 font-medium text-white hover:bg-accent-hover disabled:opacity-50 transition-colors"
            >
              {isPending ? "Signing in..." : "Sign In"}
            </button>
          </form>

          <p className="text-center text-sm text-text-secondary">
            Don&apos;t have an account?{" "}
            <Link href="/sign-up" className="text-accent hover:text-accent-hover font-medium transition-colors">
              Create one
            </Link>
          </p>
        </div>
      </div>
    </div>
  );
}
`
}

// adminSignUpPage returns the split-screen sign-up page.
func adminSignUpPage() string {
	return `"use client";

import { useState } from "react";
import Link from "next/link";
import { Eye, EyeOff } from "@/lib/icons";
import { useRegister } from "@/hooks/use-auth";

export default function SignUpPage() {
  const [name, setName] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [confirmPassword, setConfirmPassword] = useState("");
  const [showPassword, setShowPassword] = useState(false);
  const [showConfirmPassword, setShowConfirmPassword] = useState(false);
  const { mutate: register, isPending, error } = useRegister();

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (password !== confirmPassword) return;
    register({ name, email, password });
  };

  return (
    <div className="flex min-h-screen">
      {/* Left panel — branding */}
      <div className="hidden lg:flex lg:w-1/2 flex-col justify-between bg-gradient-to-br from-accent/20 via-bg-secondary to-bg-primary p-12">
        <div>
          <span className="text-2xl font-bold text-accent">G</span>
          <span className="text-2xl font-bold text-accent">rit</span>
          <span className="ml-2 rounded-md bg-accent/10 px-2 py-0.5 text-xs font-medium text-accent">
            Admin
          </span>
        </div>
        <div className="space-y-4">
          <h1 className="text-4xl font-bold text-foreground leading-tight">
            Get started with<br />your admin panel.
          </h1>
          <p className="text-text-secondary text-lg max-w-md">
            Create your account and start managing your application in minutes.
          </p>
        </div>
        <p className="text-text-muted text-sm">Built with Grit — Go + React framework</p>
      </div>

      {/* Right panel — form */}
      <div className="flex flex-1 items-center justify-center px-6 py-12">
        <div className="w-full max-w-md space-y-8">
          <div className="lg:hidden text-center mb-8">
            <span className="text-2xl font-bold text-accent">Grit</span>
            <span className="ml-2 rounded-md bg-accent/10 px-2 py-0.5 text-xs font-medium text-accent">
              Admin
            </span>
          </div>

          <div>
            <h2 className="text-2xl font-bold text-foreground">Create account</h2>
            <p className="mt-2 text-text-secondary">Sign up to get started</p>
          </div>

          <form onSubmit={handleSubmit} className="space-y-5">
            {error && (
              <div className="rounded-lg bg-danger/10 border border-danger/20 px-4 py-3 text-sm text-danger">
                {(error as unknown as { response?: { data?: { error?: { message?: string } } } })?.response?.data?.error?.message || "Registration failed"}
              </div>
            )}

            <div className="space-y-2">
              <label htmlFor="name" className="block text-sm font-medium text-text-secondary">
                Full name
              </label>
              <input
                id="name"
                type="text"
                value={name}
                onChange={(e) => setName(e.target.value)}
                className="w-full rounded-lg border border-border bg-bg-tertiary px-4 py-3 text-foreground placeholder:text-text-muted focus:border-accent focus:outline-none focus:ring-1 focus:ring-accent transition-colors"
                placeholder="John Doe"
                required
                minLength={2}
                autoFocus
              />
            </div>

            <div className="space-y-2">
              <label htmlFor="email" className="block text-sm font-medium text-text-secondary">
                Email
              </label>
              <input
                id="email"
                type="email"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                className="w-full rounded-lg border border-border bg-bg-tertiary px-4 py-3 text-foreground placeholder:text-text-muted focus:border-accent focus:outline-none focus:ring-1 focus:ring-accent transition-colors"
                placeholder="you@example.com"
                required
              />
            </div>

            <div className="space-y-2">
              <label htmlFor="password" className="block text-sm font-medium text-text-secondary">
                Password
              </label>
              <div className="relative">
                <input
                  id="password"
                  type={showPassword ? "text" : "password"}
                  value={password}
                  onChange={(e) => setPassword(e.target.value)}
                  className="w-full rounded-lg border border-border bg-bg-tertiary px-4 py-3 pr-12 text-foreground placeholder:text-text-muted focus:border-accent focus:outline-none focus:ring-1 focus:ring-accent transition-colors"
                  placeholder="Min. 8 characters"
                  required
                  minLength={8}
                />
                <button
                  type="button"
                  onClick={() => setShowPassword(!showPassword)}
                  className="absolute right-3 top-1/2 -translate-y-1/2 text-text-muted hover:text-text-secondary transition-colors"
                >
                  {showPassword ? <EyeOff className="h-5 w-5" /> : <Eye className="h-5 w-5" />}
                </button>
              </div>
            </div>

            <div className="space-y-2">
              <label htmlFor="confirmPassword" className="block text-sm font-medium text-text-secondary">
                Confirm password
              </label>
              <div className="relative">
                <input
                  id="confirmPassword"
                  type={showConfirmPassword ? "text" : "password"}
                  value={confirmPassword}
                  onChange={(e) => setConfirmPassword(e.target.value)}
                  className="w-full rounded-lg border border-border bg-bg-tertiary px-4 py-3 pr-12 text-foreground placeholder:text-text-muted focus:border-accent focus:outline-none focus:ring-1 focus:ring-accent transition-colors"
                  placeholder="Repeat your password"
                  required
                  minLength={8}
                />
                <button
                  type="button"
                  onClick={() => setShowConfirmPassword(!showConfirmPassword)}
                  className="absolute right-3 top-1/2 -translate-y-1/2 text-text-muted hover:text-text-secondary transition-colors"
                >
                  {showConfirmPassword ? <EyeOff className="h-5 w-5" /> : <Eye className="h-5 w-5" />}
                </button>
              </div>
              {password !== confirmPassword && confirmPassword && (
                <p className="text-sm text-danger">Passwords do not match</p>
              )}
            </div>

            <button
              type="submit"
              disabled={isPending || (!!confirmPassword && password !== confirmPassword)}
              className="w-full rounded-lg bg-accent py-3 font-medium text-white hover:bg-accent-hover disabled:opacity-50 transition-colors"
            >
              {isPending ? "Creating account..." : "Create Account"}
            </button>
          </form>

          <p className="text-center text-sm text-text-secondary">
            Already have an account?{" "}
            <Link href="/login" className="text-accent hover:text-accent-hover font-medium transition-colors">
              Sign in
            </Link>
          </p>
        </div>
      </div>
    </div>
  );
}
`
}

// adminForgotPasswordPage returns the split-screen forgot password page.
func adminForgotPasswordPage() string {
	return `"use client";

import { useState } from "react";
import Link from "next/link";
import { apiClient } from "@/lib/api-client";

export default function ForgotPasswordPage() {
  const [email, setEmail] = useState("");
  const [sent, setSent] = useState(false);
  const [loading, setLoading] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    try {
      await apiClient.post("/api/auth/forgot-password", { email });
      setSent(true);
    } catch {
      setSent(true); // Always show success for security
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="flex min-h-screen">
      {/* Left panel — branding */}
      <div className="hidden lg:flex lg:w-1/2 flex-col justify-between bg-gradient-to-br from-accent/20 via-bg-secondary to-bg-primary p-12">
        <div>
          <span className="text-2xl font-bold text-accent">G</span>
          <span className="text-2xl font-bold text-accent">rit</span>
          <span className="ml-2 rounded-md bg-accent/10 px-2 py-0.5 text-xs font-medium text-accent">
            Admin
          </span>
        </div>
        <div className="space-y-4">
          <h1 className="text-4xl font-bold text-foreground leading-tight">
            Reset your<br />password.
          </h1>
          <p className="text-text-secondary text-lg max-w-md">
            Enter your email and we&apos;ll send you a link to get back into your account.
          </p>
        </div>
        <p className="text-text-muted text-sm">Built with Grit — Go + React framework</p>
      </div>

      {/* Right panel — form */}
      <div className="flex flex-1 items-center justify-center px-6 py-12">
        <div className="w-full max-w-md space-y-8">
          <div className="lg:hidden text-center mb-8">
            <span className="text-2xl font-bold text-accent">Grit</span>
          </div>

          <div>
            <h2 className="text-2xl font-bold text-foreground">Forgot password?</h2>
            <p className="mt-2 text-text-secondary">No worries, we&apos;ll send you reset instructions.</p>
          </div>

          {sent ? (
            <div className="rounded-xl bg-bg-secondary border border-border p-8 text-center space-y-4">
              <div className="flex justify-center">
                <div className="h-12 w-12 rounded-full bg-success/10 flex items-center justify-center">
                  <svg className="h-6 w-6 text-success" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M3 8l7.89 5.26a2 2 0 002.22 0L21 8M5 19h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z" />
                  </svg>
                </div>
              </div>
              <h3 className="text-lg font-medium text-foreground">Check your email</h3>
              <p className="text-text-secondary text-sm">
                If an account with that email exists, we&apos;ve sent a password reset link.
              </p>
              <Link
                href="/login"
                className="inline-block text-accent hover:text-accent-hover font-medium text-sm transition-colors"
              >
                Back to login
              </Link>
            </div>
          ) : (
            <form onSubmit={handleSubmit} className="space-y-5">
              <div className="space-y-2">
                <label htmlFor="email" className="block text-sm font-medium text-text-secondary">
                  Email
                </label>
                <input
                  id="email"
                  type="email"
                  value={email}
                  onChange={(e) => setEmail(e.target.value)}
                  className="w-full rounded-lg border border-border bg-bg-tertiary px-4 py-3 text-foreground placeholder:text-text-muted focus:border-accent focus:outline-none focus:ring-1 focus:ring-accent transition-colors"
                  placeholder="you@example.com"
                  required
                  autoFocus
                />
              </div>

              <button
                type="submit"
                disabled={loading}
                className="w-full rounded-lg bg-accent py-3 font-medium text-white hover:bg-accent-hover disabled:opacity-50 transition-colors"
              >
                {loading ? "Sending..." : "Send Reset Link"}
              </button>
            </form>
          )}

          <p className="text-center text-sm text-text-secondary">
            <Link href="/login" className="text-accent hover:text-accent-hover font-medium transition-colors">
              Back to login
            </Link>
          </p>
        </div>
      </div>
    </div>
  );
}
`
}
