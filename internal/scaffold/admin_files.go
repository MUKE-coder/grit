package scaffold

import (
	"fmt"
	"path/filepath"
)

func writeAdminFiles(root string, opts Options) error {
	adminRoot := filepath.Join(root, "apps", "admin")

	files := map[string]string{
		filepath.Join(adminRoot, "package.json"):                                  adminPackageJSON(opts),
		filepath.Join(adminRoot, "next.config.ts"):                                adminNextConfig(),
		filepath.Join(adminRoot, "tailwind.config.ts"):                            adminTailwindConfig(),
		filepath.Join(adminRoot, "postcss.config.js"):                             adminPostCSSConfig(),
		filepath.Join(adminRoot, "tsconfig.json"):                                 adminTSConfig(),
		filepath.Join(adminRoot, "app", "globals.css"):                            adminGlobalCSS(),
		filepath.Join(adminRoot, "app", "layout.tsx"):                             adminRootLayout(opts),
		filepath.Join(adminRoot, "app", "page.tsx"):                               adminDashboardPage(),
		filepath.Join(adminRoot, "app", "resources", "users", "page.tsx"):         adminUsersPage(),
		filepath.Join(adminRoot, "components", "layout", "admin-layout.tsx"):      adminLayoutComponent(),
		filepath.Join(adminRoot, "components", "layout", "sidebar.tsx"):           adminSidebar(),
		filepath.Join(adminRoot, "components", "layout", "navbar.tsx"):            adminNavbar(),
		filepath.Join(adminRoot, "components", "widgets", "stats-card.tsx"):       adminStatsCard(),
		filepath.Join(adminRoot, "components", "tables", "data-table.tsx"):        adminDataTable(),
		filepath.Join(adminRoot, "hooks", "use-auth.ts"):                          adminUseAuth(),
		filepath.Join(adminRoot, "hooks", "use-users.ts"):                         adminUseUsers(),
		filepath.Join(adminRoot, "lib", "api-client.ts"):                          adminAPIClient(),
		filepath.Join(adminRoot, "lib", "query-client.ts"):                        adminQueryClient(),
		filepath.Join(adminRoot, "lib", "utils.ts"):                               adminUtils(),
		filepath.Join(adminRoot, "components", "shared", "providers.tsx"):         adminProviders(),
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
    "@tanstack/react-query": "^5.17.0",
    "axios": "^1.6.0",
    "class-variance-authority": "^0.7.0",
    "clsx": "^2.1.0",
    "js-cookie": "^3.0.5",
    "lucide-react": "^0.303.0",
    "next": "14.1.0",
    "react": "^18.2.0",
    "react-dom": "^18.2.0",
    "sonner": "^1.3.0",
    "tailwind-merge": "^2.2.0",
    "tailwindcss-animate": "^1.0.7",
    "zod": "^3.22.0"
  },
  "devDependencies": {
    "@types/js-cookie": "^3.0.6",
    "@types/node": "^20.0.0",
    "@types/react": "^18.2.0",
    "@types/react-dom": "^18.2.0",
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
        sans: ["DM Sans", "system-ui", "sans-serif"],
        mono: ["JetBrains Mono", "monospace"],
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

@import url("https://fonts.googleapis.com/css2?family=DM+Sans:wght@400;500;600;700&family=JetBrains+Mono:wght@400;500;600&display=swap");

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

body {
  background-color: var(--bg-primary);
  color: var(--text-primary);
  font-family: "DM Sans", system-ui, sans-serif;
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
import "./globals.css";
import { Providers } from "@/components/shared/providers";
import { AdminLayout } from "@/components/layout/admin-layout";

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
      <body className="min-h-screen bg-background font-sans antialiased">
        <Providers>
          <AdminLayout>{children}</AdminLayout>
        </Providers>
      </body>
    </html>
  );
}
`, opts.ProjectName)
}

func adminDashboardPage() string {
	return `"use client";

import { StatsCard } from "@/components/widgets/stats-card";

const stats = [
  { label: "Total Users", value: "128", change: "+12%", icon: "👥", gradient: "from-accent/20 to-accent/5" },
  { label: "Active Users", value: "96", change: "+8%", icon: "✅", gradient: "from-success/20 to-success/5" },
  { label: "New This Month", value: "24", change: "+24%", icon: "📈", gradient: "from-info/20 to-info/5" },
  { label: "Revenue", value: "$12,450", change: "+18%", icon: "💰", gradient: "from-warning/20 to-warning/5" },
];

export default function AdminDashboard() {
  return (
    <div className="space-y-8">
      <div>
        <h1 className="text-2xl font-bold text-foreground">Dashboard</h1>
        <p className="text-text-secondary mt-1">Overview of your application</p>
      </div>

      <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-4">
        {stats.map((stat) => (
          <StatsCard key={stat.label} {...stat} />
        ))}
      </div>

      <div className="grid gap-6 lg:grid-cols-2">
        <div className="rounded-xl border border-border bg-bg-secondary p-6">
          <h2 className="text-lg font-semibold text-foreground mb-4">Recent Activity</h2>
          <div className="space-y-4">
            {[1, 2, 3, 4, 5].map((i) => (
              <div key={i} className="flex items-center gap-3 text-sm">
                <div className="h-2 w-2 rounded-full bg-accent" />
                <span className="text-text-secondary">User activity placeholder #{i}</span>
                <span className="ml-auto text-text-muted text-xs">Just now</span>
              </div>
            ))}
          </div>
        </div>

        <div className="rounded-xl border border-border bg-bg-secondary p-6">
          <h2 className="text-lg font-semibold text-foreground mb-4">Quick Actions</h2>
          <div className="grid gap-3 sm:grid-cols-2">
            <a
              href="/resources/users"
              className="rounded-lg border border-border bg-bg-tertiary p-4 hover:bg-bg-hover transition-colors"
            >
              <h3 className="font-medium text-foreground">Manage Users</h3>
              <p className="text-xs text-text-muted mt-1">View and manage users</p>
            </a>
            <a
              href="http://localhost:8080/studio"
              target="_blank"
              rel="noopener noreferrer"
              className="rounded-lg border border-border bg-bg-tertiary p-4 hover:bg-bg-hover transition-colors"
            >
              <h3 className="font-medium text-foreground">GORM Studio</h3>
              <p className="text-xs text-text-muted mt-1">Browse database</p>
            </a>
          </div>
        </div>
      </div>
    </div>
  );
}
`
}

func adminUsersPage() string {
	return `"use client";

import { useState } from "react";
import { DataTable } from "@/components/tables/data-table";
import { useUsers, useDeleteUser } from "@/hooks/use-users";

export default function UsersPage() {
  const [page, setPage] = useState(1);
  const [search, setSearch] = useState("");
  const { data, isLoading } = useUsers({ page, search, pageSize: 20 });
  const { mutate: deleteUser } = useDeleteUser();

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold text-foreground">Users</h1>
          <p className="text-text-secondary mt-1">Manage user accounts</p>
        </div>
      </div>

      <div className="rounded-xl border border-border bg-bg-secondary">
        <div className="p-4 border-b border-border">
          <input
            type="text"
            value={search}
            onChange={(e) => {
              setSearch(e.target.value);
              setPage(1);
            }}
            placeholder="Search users..."
            className="w-full max-w-sm rounded-lg border border-border bg-bg-tertiary px-4 py-2 text-sm text-foreground placeholder:text-text-muted focus:border-accent focus:outline-none focus:ring-1 focus:ring-accent"
          />
        </div>

        <DataTable
          data={data?.data || []}
          isLoading={isLoading}
          columns={[
            { key: "id", label: "ID", width: "80px" },
            { key: "name", label: "Name" },
            { key: "email", label: "Email" },
            {
              key: "role",
              label: "Role",
              render: (value: string) => (
                <span
                  className={` + "`" + `inline-flex items-center rounded-full px-2.5 py-0.5 text-xs font-medium ${
                    value === "admin"
                      ? "bg-accent/10 text-accent"
                      : value === "editor"
                      ? "bg-info/10 text-info"
                      : "bg-bg-hover text-text-secondary"
                  }` + "`" + `}
                >
                  {value}
                </span>
              ),
            },
            {
              key: "active",
              label: "Active",
              render: (value: boolean) => (
                <span
                  className={` + "`" + `inline-flex items-center rounded-full px-2.5 py-0.5 text-xs font-medium ${
                    value
                      ? "bg-success/10 text-success"
                      : "bg-danger/10 text-danger"
                  }` + "`" + `}
                >
                  {value ? "Active" : "Inactive"}
                </span>
              ),
            },
            {
              key: "created_at",
              label: "Created At",
              render: (value: string) =>
                new Date(value).toLocaleDateString("en-US", {
                  year: "numeric",
                  month: "short",
                  day: "numeric",
                }),
            },
          ]}
          actions={(row) => (
            <div className="flex gap-2">
              <button
                onClick={() => deleteUser(row.id)}
                className="text-xs text-danger hover:text-danger/80"
              >
                Delete
              </button>
            </div>
          )}
        />

        {data?.meta && (
          <div className="flex items-center justify-between border-t border-border p-4">
            <p className="text-sm text-text-muted">
              Showing {data.data.length} of {data.meta.total} users
            </p>
            <div className="flex gap-2">
              <button
                onClick={() => setPage(Math.max(1, page - 1))}
                disabled={page <= 1}
                className="rounded-lg border border-border bg-bg-tertiary px-3 py-1.5 text-sm text-text-secondary hover:bg-bg-hover disabled:opacity-50"
              >
                Previous
              </button>
              <span className="flex items-center px-3 text-sm text-text-muted">
                Page {page} of {data.meta.pages}
              </span>
              <button
                onClick={() => setPage(Math.min(data.meta.pages, page + 1))}
                disabled={page >= data.meta.pages}
                className="rounded-lg border border-border bg-bg-tertiary px-3 py-1.5 text-sm text-text-secondary hover:bg-bg-hover disabled:opacity-50"
              >
                Next
              </button>
            </div>
          </div>
        )}
      </div>
    </div>
  );
}
`
}

func adminLayoutComponent() string {
	return `"use client";

import { useEffect } from "react";
import { useRouter } from "next/navigation";
import { useMe } from "@/hooks/use-auth";
import { Sidebar } from "./sidebar";
import { Navbar } from "./navbar";

export function AdminLayout({ children }: { children: React.ReactNode }) {
  const { data: user, isLoading, isError } = useMe();
  const router = useRouter();

  useEffect(() => {
    if (isError) {
      router.push("/login");
    }
  }, [isError, router]);

  if (isLoading) {
    return (
      <div className="flex min-h-screen items-center justify-center">
        <div className="h-8 w-8 animate-spin rounded-full border-2 border-accent border-t-transparent" />
      </div>
    );
  }

  if (!user) return null;

  return (
    <div className="flex min-h-screen">
      <Sidebar user={user} />
      <div className="flex flex-1 flex-col lg:ml-64">
        <Navbar user={user} />
        <main className="flex-1 p-6">{children}</main>
      </div>
    </div>
  );
}
`
}

func adminSidebar() string {
	return `"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";

interface User {
  name: string;
  email: string;
  role: string;
}

const navItems = [
  { label: "Dashboard", href: "/", icon: "📊" },
  { label: "Users", href: "/resources/users", icon: "👥" },
];

export function Sidebar({ user }: { user: User }) {
  const pathname = usePathname();

  return (
    <aside className="fixed inset-y-0 left-0 z-50 hidden w-64 bg-bg-secondary border-r border-border lg:block">
      <div className="flex h-16 items-center gap-2 border-b border-border px-6">
        <span className="text-xl font-bold text-accent">Grit</span>
        <span className="rounded-md bg-accent/10 px-2 py-0.5 text-xs font-medium text-accent">
          Admin
        </span>
      </div>

      <nav className="p-4 space-y-1">
        {navItems.map((item) => (
          <Link
            key={item.href}
            href={item.href}
            className={` + "`" + `flex items-center gap-3 rounded-lg px-3 py-2.5 text-sm font-medium transition-colors ${
              pathname === item.href
                ? "bg-accent/10 text-accent"
                : "text-text-secondary hover:bg-bg-hover hover:text-foreground"
            }` + "`" + `}
          >
            <span>{item.icon}</span>
            {item.label}
          </Link>
        ))}
      </nav>

      <div className="absolute bottom-0 left-0 right-0 border-t border-border p-4">
        <div className="flex items-center gap-3">
          <div className="h-8 w-8 rounded-full bg-accent/20 flex items-center justify-center">
            <span className="text-sm font-medium text-accent">
              {user.name?.charAt(0)?.toUpperCase()}
            </span>
          </div>
          <div className="flex-1 min-w-0">
            <p className="text-sm font-medium text-foreground truncate">{user.name}</p>
            <p className="text-xs text-text-muted truncate">{user.role}</p>
          </div>
        </div>
      </div>
    </aside>
  );
}
`
}

func adminNavbar() string {
	return `"use client";

import { useState } from "react";
import { useLogout } from "@/hooks/use-auth";

interface User {
  name: string;
  email: string;
}

export function Navbar({ user }: { user: User }) {
  const { mutate: logout } = useLogout();
  const [dropdownOpen, setDropdownOpen] = useState(false);

  return (
    <header className="sticky top-0 z-30 flex h-16 items-center justify-between border-b border-border bg-bg-primary px-6">
      <div className="flex-1" />

      <div className="relative">
        <button
          onClick={() => setDropdownOpen(!dropdownOpen)}
          className="flex items-center gap-3 rounded-lg px-3 py-2 hover:bg-bg-hover transition-colors"
        >
          <div className="h-8 w-8 rounded-full bg-accent/20 flex items-center justify-center">
            <span className="text-sm font-medium text-accent">
              {user.name?.charAt(0)?.toUpperCase()}
            </span>
          </div>
          <span className="text-sm font-medium text-foreground hidden sm:block">
            {user.name}
          </span>
        </button>

        {dropdownOpen && (
          <>
            <div className="fixed inset-0 z-40" onClick={() => setDropdownOpen(false)} />
            <div className="absolute right-0 top-full mt-2 w-48 rounded-lg border border-border bg-bg-elevated shadow-lg z-50">
              <div className="px-4 py-3 border-b border-border">
                <p className="text-sm font-medium text-foreground">{user.name}</p>
                <p className="text-xs text-text-muted">{user.email}</p>
              </div>
              <button
                onClick={() => logout()}
                className="w-full px-4 py-2.5 text-left text-sm text-danger hover:bg-bg-hover transition-colors"
              >
                Sign Out
              </button>
            </div>
          </>
        )}
      </div>
    </header>
  );
}
`
}

func adminStatsCard() string {
	return `interface StatsCardProps {
  label: string;
  value: string;
  change?: string;
  icon: string;
  gradient: string;
}

export function StatsCard({ label, value, change, icon, gradient }: StatsCardProps) {
  return (
    <div className={` + "`" + `rounded-xl border border-border bg-gradient-to-br ${gradient} p-6` + "`" + `}>
      <div className="flex items-center justify-between">
        <span className="text-2xl">{icon}</span>
        {change && (
          <span className="text-xs font-medium text-success">{change}</span>
        )}
      </div>
      <div className="mt-4">
        <p className="text-3xl font-bold text-foreground">{value}</p>
        <p className="text-sm text-text-secondary mt-1">{label}</p>
      </div>
    </div>
  );
}
`
}

func adminDataTable() string {
	return `interface Column {
  key: string;
  label: string;
  width?: string;
  render?: (value: any, row: any) => React.ReactNode;
}

interface DataTableProps {
  data: any[];
  columns: Column[];
  isLoading?: boolean;
  actions?: (row: any) => React.ReactNode;
}

export function DataTable({ data, columns, isLoading, actions }: DataTableProps) {
  if (isLoading) {
    return (
      <div className="p-8 text-center">
        <div className="inline-block h-6 w-6 animate-spin rounded-full border-2 border-accent border-t-transparent" />
        <p className="mt-2 text-sm text-text-muted">Loading...</p>
      </div>
    );
  }

  if (data.length === 0) {
    return (
      <div className="p-8 text-center">
        <p className="text-text-muted">No data found</p>
      </div>
    );
  }

  return (
    <div className="overflow-x-auto">
      <table className="w-full">
        <thead>
          <tr className="border-b border-border">
            {columns.map((col) => (
              <th
                key={col.key}
                className="px-4 py-3 text-left text-xs font-medium text-text-muted uppercase tracking-wider"
                style={col.width ? { width: col.width } : undefined}
              >
                {col.label}
              </th>
            ))}
            {actions && (
              <th className="px-4 py-3 text-right text-xs font-medium text-text-muted uppercase tracking-wider w-[100px]">
                Actions
              </th>
            )}
          </tr>
        </thead>
        <tbody>
          {data.map((row, idx) => (
            <tr
              key={row.id || idx}
              className="border-b border-border/50 hover:bg-bg-hover/50 transition-colors"
            >
              {columns.map((col) => (
                <td key={col.key} className="px-4 py-3 text-sm text-foreground">
                  {col.render ? col.render(row[col.key], row) : row[col.key]}
                </td>
              ))}
              {actions && (
                <td className="px-4 py-3 text-right text-sm">
                  {actions(row)}
                </td>
              )}
            </tr>
          ))}
        </tbody>
      </table>
    </div>
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

func adminUseUsers() string {
	return `import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { apiClient } from "@/lib/api-client";

interface User {
  id: number;
  name: string;
  email: string;
  role: string;
  active: boolean;
  created_at: string;
}

interface UsersResponse {
  data: User[];
  meta: {
    total: number;
    page: number;
    page_size: number;
    pages: number;
  };
}

interface UseUsersParams {
  page?: number;
  pageSize?: number;
  search?: string;
}

export function useUsers({ page = 1, pageSize = 20, search = "" }: UseUsersParams = {}) {
  return useQuery<UsersResponse>({
    queryKey: ["users", { page, pageSize, search }],
    queryFn: async () => {
      const params = new URLSearchParams({
        page: String(page),
        page_size: String(pageSize),
      });
      if (search) {
        params.set("search", search);
      }
      const { data } = await apiClient.get(` + "`" + `/api/users?${params}` + "`" + `);
      return data;
    },
  });
}

export function useDeleteUser() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (id: number) => {
      await apiClient.delete(` + "`" + `/api/users/${id}` + "`" + `);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["users"] });
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

func adminProviders() string {
	return `"use client";

import { QueryClientProvider } from "@tanstack/react-query";
import { queryClient } from "@/lib/query-client";

export function Providers({ children }: { children: React.ReactNode }) {
  return (
    <QueryClientProvider client={queryClient}>
      {children}
    </QueryClientProvider>
  );
}
`
}
