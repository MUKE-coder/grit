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

		// UI components
		filepath.Join(adminRoot, "components", "ui", "dropzone.tsx"): adminDropzone(),

		// Widget components
		filepath.Join(adminRoot, "components", "widgets", "stats-card.tsx"):      adminStatsCard(),
		filepath.Join(adminRoot, "components", "widgets", "chart-widget.tsx"):    adminChartWidget(),
		filepath.Join(adminRoot, "components", "widgets", "activity-widget.tsx"): adminActivityWidget(),
		filepath.Join(adminRoot, "components", "widgets", "widget-grid.tsx"):     adminWidgetGrid(),

		// Resource components
		filepath.Join(adminRoot, "components", "resource", "resource-page.tsx"): adminResourcePage(),

		// Resource definitions
		filepath.Join(adminRoot, "resources", "index.ts"): adminResourceRegistry(),
		filepath.Join(adminRoot, "resources", "users.ts"): adminUsersResource(),

		// Hooks
		filepath.Join(adminRoot, "hooks", "use-auth.ts"):     adminUseAuth(),
		filepath.Join(adminRoot, "hooks", "use-resource.ts"): adminUseResource(),
		filepath.Join(adminRoot, "hooks", "use-system.ts"):   adminUseSystem(),

		// Pages
		filepath.Join(adminRoot, "app", "page.tsx"):                       adminDashboardPage(),
		filepath.Join(adminRoot, "app", "resources", "users", "page.tsx"): adminUsersPage(),

		// System pages
		filepath.Join(adminRoot, "app", "system", "jobs", "page.tsx"):  adminJobsPage(),
		filepath.Join(adminRoot, "app", "system", "files", "page.tsx"): adminFilesPage(),
		filepath.Join(adminRoot, "app", "system", "cron", "page.tsx"):  adminCronPage(),
		filepath.Join(adminRoot, "app", "system", "mail", "page.tsx"):  adminMailPage(),
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
import { AdminLayout } from "@/components/layout/admin-layout";

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
        <Providers>
          <AdminLayout>{children}</AdminLayout>
        </Providers>
      </body>
    </html>
  );
}
`, opts.ProjectName)
}

func adminProviders() string {
	return `"use client";

import { QueryClientProvider } from "@tanstack/react-query";
import { queryClient } from "@/lib/query-client";
import { ThemeProvider } from "./theme-provider";

export function Providers({ children }: { children: React.ReactNode }) {
  return (
    <ThemeProvider>
      <QueryClientProvider client={queryClient}>
        {children}
      </QueryClientProvider>
    </ThemeProvider>
  );
}
`
}

func adminUseAuth() string {
	return `import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
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
