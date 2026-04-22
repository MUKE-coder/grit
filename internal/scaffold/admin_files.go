package scaffold

import (
	"fmt"
	"path/filepath"
)

// adminLoginPageForStyle returns the login page for the given style variant.
func adminLoginPageForStyle(style string) string {
	switch style {
	case "modern":
		return modernLoginPage()
	case "minimal":
		return minimalLoginPage()
	case "glass":
		return glassLoginPage()
	case "centered":
		return centeredLoginPage()
	default:
		return adminLoginPage()
	}
}

// adminSignUpPageForStyle returns the sign-up page for the given style variant.
func adminSignUpPageForStyle(style string) string {
	switch style {
	case "modern":
		return modernSignUpPage()
	case "minimal":
		return minimalSignUpPage()
	case "glass":
		return glassSignUpPage()
	case "centered":
		return centeredSignUpPage()
	default:
		return adminSignUpPage()
	}
}

// adminForgotPasswordPageForStyle returns the forgot password page for the given style variant.
func adminForgotPasswordPageForStyle(style string) string {
	switch style {
	case "modern":
		return modernForgotPasswordPage()
	case "minimal":
		return minimalForgotPasswordPage()
	case "glass":
		return glassForgotPasswordPage()
	case "centered":
		return centeredForgotPasswordPage()
	default:
		return adminForgotPasswordPage()
	}
}

// adminDashboardPageForStyle returns the dashboard page for the given style variant.
func adminDashboardPageForStyle(style string) string {
	switch style {
	case "modern":
		return modernDashboardPage()
	case "minimal":
		return minimalDashboardPage()
	case "glass":
		return glassDashboardPage()
	case "centered":
		// Centered style uses the default dashboard (only auth pages differ)
		return adminDashboardPage()
	default:
		return adminDashboardPage()
	}
}

// socialLoginButtonsJSX returns the shared social login buttons JSX snippet.
// Insert this between the submit button and the "Don't have an account?" or "Already have an account?" text.
func socialLoginButtonsJSX() string {
	return `
          {/* Social Login */}
          <div className="relative my-6">
            <div className="absolute inset-0 flex items-center">
              <div className="w-full border-t border-border" />
            </div>
            <div className="relative flex justify-center text-xs">
              <span className="bg-bg-primary px-3 text-text-muted lg:bg-transparent">or continue with</span>
            </div>
          </div>

          <div className="flex gap-3">
            <a
              href={` + "`" + `${process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080"}/api/auth/oauth/google` + "`" + `}
              className="flex flex-1 items-center justify-center gap-2 rounded-lg border border-border bg-white px-4 py-2.5 text-sm font-medium text-gray-700 hover:bg-gray-50 transition-colors"
            >
              <svg className="h-5 w-5" viewBox="0 0 24 24">
                <path d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92a5.06 5.06 0 0 1-2.2 3.32v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.1z" fill="#4285F4" />
                <path d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z" fill="#34A853" />
                <path d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z" fill="#FBBC05" />
                <path d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z" fill="#EA4335" />
              </svg>
              Google
            </a>
            <a
              href={` + "`" + `${process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080"}/api/auth/oauth/github` + "`" + `}
              className="flex flex-1 items-center justify-center gap-2 rounded-lg border border-border bg-[#24292f] px-4 py-2.5 text-sm font-medium text-white hover:bg-[#2f363d] transition-colors"
            >
              <svg className="h-5 w-5" fill="currentColor" viewBox="0 0 24 24">
                <path fillRule="evenodd" d="M12 2C6.477 2 2 6.484 2 12.017c0 4.425 2.865 8.18 6.839 9.504.5.092.682-.217.682-.483 0-.237-.008-.868-.013-1.703-2.782.605-3.369-1.343-3.369-1.343-.454-1.158-1.11-1.466-1.11-1.466-.908-.62.069-.608.069-.608 1.003.07 1.531 1.032 1.531 1.032.892 1.53 2.341 1.088 2.91.832.092-.647.35-1.088.636-1.338-2.22-.253-4.555-1.113-4.555-4.951 0-1.093.39-1.988 1.029-2.688-.103-.253-.446-1.272.098-2.65 0 0 .84-.27 2.75 1.026A9.564 9.564 0 0 1 12 6.844a9.59 9.59 0 0 1 2.504.337c1.909-1.296 2.747-1.027 2.747-1.027.546 1.379.202 2.398.1 2.651.64.7 1.028 1.595 1.028 2.688 0 3.848-2.339 4.695-4.566 4.943.359.309.678.92.678 1.855 0 1.338-.012 2.419-.012 2.747 0 .268.18.58.688.482A10.02 10.02 0 0 0 22 12.017C22 6.484 17.522 2 12 2z" clipRule="evenodd" />
              </svg>
              GitHub
            </a>
          </div>
`
}

// adminAuthCallbackPage returns the OAuth callback page for the admin panel.
func adminAuthCallbackPage() string {
	return `"use client";

import { useEffect, useRef } from "react";
import { useRouter, useSearchParams } from "next/navigation";
import Cookies from "js-cookie";
import { apiClient } from "@/lib/api-client";
import { useQueryClient } from "@tanstack/react-query";

export default function AuthCallbackPage() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const queryClient = useQueryClient();
  const processed = useRef(false);

  useEffect(() => {
    if (processed.current) return;
    processed.current = true;

    const accessToken = searchParams.get("access_token");
    const refreshToken = searchParams.get("refresh_token");
    const error = searchParams.get("error");

    if (error) {
      router.push("/login?error=" + encodeURIComponent(error));
      return;
    }

    if (accessToken && refreshToken) {
      Cookies.set("access_token", accessToken, { expires: 1 });
      Cookies.set("refresh_token", refreshToken, { expires: 7 });

      // Fetch user to determine redirect
      apiClient
        .get("/api/auth/me")
        .then(({ data }) => {
          const user = data.data;
          queryClient.setQueryData(["me"], user);
          router.push(user.role === "USER" ? "/profile" : "/dashboard");
        })
        .catch(() => {
          router.push("/dashboard");
        });
    } else {
      router.push("/login?error=Authentication+failed");
    }
  }, [searchParams, router, queryClient]);

  return (
    <div className="flex min-h-screen items-center justify-center bg-bg-primary">
      <div className="text-center">
        <div className="inline-flex h-10 w-10 animate-spin items-center justify-center rounded-full border-2 border-accent border-t-transparent" />
        <p className="mt-4 text-sm text-text-secondary">Signing you in...</p>
      </div>
    </div>
  );
}
`
}

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

		// Error pages
		filepath.Join(adminRoot, "app", "error.tsx"):        adminErrorPage(),
		filepath.Join(adminRoot, "app", "not-found.tsx"):    adminNotFoundPage(),
		filepath.Join(adminRoot, "app", "global-error.tsx"): adminGlobalErrorPage(),

		// Auth pages — (auth) route group (style variant)
		filepath.Join(adminRoot, "app", "(auth)", "login", "page.tsx"):           adminLoginPageForStyle(opts.Style),
		filepath.Join(adminRoot, "app", "(auth)", "sign-up", "page.tsx"):         adminSignUpPageForStyle(opts.Style),
		filepath.Join(adminRoot, "app", "(auth)", "forgot-password", "page.tsx"): adminForgotPasswordPageForStyle(opts.Style),
		filepath.Join(adminRoot, "app", "(auth)", "callback", "page.tsx"):        adminAuthCallbackPage(),

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
		filepath.Join(adminRoot, "components", "layout", "page-header.tsx"):  adminPageHeader(),

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
		filepath.Join(adminRoot, "components", "forms", "form-builder.tsx"):      adminFormBuilder(),
		filepath.Join(adminRoot, "components", "forms", "form-modal.tsx"):        adminFormModal(),
		filepath.Join(adminRoot, "components", "forms", "form-page.tsx"):         adminFormPage(),
		filepath.Join(adminRoot, "components", "forms", "form-stepper.tsx"):      adminFormStepper(),
		filepath.Join(adminRoot, "components", "forms", "form-modal-steps.tsx"):  adminFormModalSteps(),
		filepath.Join(adminRoot, "components", "forms", "form-page-steps.tsx"):   adminFormPageSteps(),
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
		filepath.Join(adminRoot, "components", "forms", "fields", "files-field.tsx"):                  adminFilesField(),
		filepath.Join(adminRoot, "components", "forms", "fields", "relationship-select-field.tsx"):      adminRelationshipSelectField(),
		filepath.Join(adminRoot, "components", "forms", "fields", "multi-relationship-select-field.tsx"): adminMultiRelationshipSelectField(),
		filepath.Join(adminRoot, "components", "forms", "fields", "rich-text-field.tsx"):                adminRichTextField(),

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
		filepath.Join(adminRoot, "resources", "blogs.ts"): adminBlogsResource(),

		// Profile components
		filepath.Join(adminRoot, "components", "profile", "delete-account-dialog.tsx"): adminDeleteAccountDialog(),

		// Hooks
		filepath.Join(adminRoot, "hooks", "use-auth.ts"):     adminUseAuth(),
		filepath.Join(adminRoot, "hooks", "use-resource.ts"): adminUseResource(),
		filepath.Join(adminRoot, "hooks", "use-system.ts"):   adminUseSystem(),
		filepath.Join(adminRoot, "hooks", "use-profile.ts"):  adminUseProfile(),

		// Dashboard pages — (dashboard) route group (style variant)
		filepath.Join(adminRoot, "app", "(dashboard)", "dashboard", "page.tsx"):          adminDashboardPageForStyle(opts.Style),
		filepath.Join(adminRoot, "app", "(dashboard)", "profile", "page.tsx"):            adminProfilePage(),
		filepath.Join(adminRoot, "app", "(dashboard)", "resources", "users", "page.tsx"): adminUsersPage(),
		filepath.Join(adminRoot, "app", "(dashboard)", "resources", "blogs", "page.tsx"): adminBlogsPage(),

		// System pages — under (dashboard) route group
		filepath.Join(adminRoot, "app", "(dashboard)", "system", "jobs", "page.tsx"):  adminJobsPage(),
		filepath.Join(adminRoot, "app", "(dashboard)", "system", "files", "page.tsx"): adminFilesPage(),
		filepath.Join(adminRoot, "app", "(dashboard)", "system", "cron", "page.tsx"):  adminCronPage(),
		filepath.Join(adminRoot, "app", "(dashboard)", "system", "mail", "page.tsx"):     adminMailPage(),
		filepath.Join(adminRoot, "app", "(dashboard)", "system", "security", "page.tsx"): adminSecurityPage(),
		filepath.Join(adminRoot, "public", ".gitkeep"):                                   "",
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
    "dev": "rm -rf .next && next dev --port 3001",
    "build": "next build",
    "start": "next start",
    "lint": "next lint",
    "format": "prettier --write .",
    "analyze": "ANALYZE=true next build",
    "test": "vitest run",
    "test:watch": "vitest",
    "test:ui": "vitest --ui"
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
    "@tiptap/extension-link": "^2.1.0",
    "@tiptap/pm": "^2.1.0",
    "@tiptap/react": "^2.1.0",
    "@tiptap/starter-kit": "^2.1.0",
    "react-dropzone": "^14.2.0",
    "react-hook-form": "^7.49.0",
    "recharts": "^2.12.0",
    "sonner": "^1.3.0",
    "tailwind-merge": "^2.2.0",
    "tailwindcss-animate": "^1.0.7",
    "zod": "^3.22.0",
    "@repo/shared": "workspace:*"
  },
  "devDependencies": {
    "@testing-library/jest-dom": "^6.4.0",
    "@testing-library/react": "^16.0.0",
    "@testing-library/user-event": "^14.5.0",
    "@vitejs/plugin-react": "^4.3.0",
    "@types/js-cookie": "^3.0.6",
    "@types/node": "^20.0.0",
    "@types/react": "^19.0.0",
    "@types/react-dom": "^19.0.0",
    "autoprefixer": "^10.4.0",
    "jsdom": "^25.0.0",
    "postcss": "^8.4.0",
    "prettier": "^3.3.0",
    "prettier-plugin-tailwindcss": "^0.6.0",
    "tailwindcss": "^3.4.0",
    "typescript": "^5.3.0",
    "vitest": "^2.0.0"
  }
}
`, opts.ProjectName)
}

func adminNextConfig() string {
	return `import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  output: "standalone",
  reactStrictMode: true,
  transpilePackages: ["@repo/shared"],
  // Uncomment and run "ANALYZE=true pnpm build" to inspect the bundle
  // ...(process.env.ANALYZE === "true"
  //   ? { ...require("@next/bundle-analyzer")({ enabled: true })(nextConfig) }
  //   : {}),
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
        sans: ["var(--font-onest)", "system-ui", "sans-serif"],
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
  font-family: var(--font-onest), system-ui, sans-serif;
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
import { Onest, JetBrains_Mono } from "next/font/google";
import "./globals.css";
import { Providers } from "@/components/shared/providers";

const onest = Onest({
  subsets: ["latin"],
  variable: "--font-onest",
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
      <body className={` + "`" + `${onest.variable} ${jetbrainsMono.variable} min-h-screen bg-background font-sans antialiased` + "`" + `}>
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
  first_name: string;
  last_name: string;
  email: string;
  role: string;
  avatar: string;
  job_title: string;
  bio: string;
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
      router.push(data.data.user.role === "USER" ? "/profile" : "/dashboard");
    },
  });
}

export function useRegister() {
  const router = useRouter();
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (data: {
      first_name: string;
      last_name: string;
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
      router.push(data.data.user.role === "USER" ? "/profile" : "/dashboard");
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

/**
 * Upload a file via presigned URL (browser uploads directly to storage).
 * 1. POST /api/uploads/presign → get presigned PUT URL
 * 2. XHR PUT to presigned URL (direct to R2/S3/MinIO)
 * 3. POST /api/uploads/complete → record in DB
 */
export async function uploadFile(
  file: File,
  _endpoint = "/api/uploads",
  onProgress?: (percent: number) => void
): Promise<{ data: Record<string, unknown>; message: string }> {
  // Step 1: Get presigned URL from API
  const { data: presignRes } = await apiClient.post("/api/uploads/presign", {
    filename: file.name,
    content_type: file.type,
    file_size: file.size,
  });
  const { presigned_url, key } = presignRes.data;

  // Step 2: Upload directly to storage via XHR PUT (bypasses API server)
  await new Promise<void>((resolve, reject) => {
    const xhr = new XMLHttpRequest();
    if (onProgress) {
      xhr.upload.onprogress = (e) => {
        if (e.lengthComputable) onProgress(Math.round((e.loaded / e.total) * 100));
      };
    }
    xhr.onload = () =>
      xhr.status >= 200 && xhr.status < 300
        ? resolve()
        : reject(new Error(` + "`" + `Storage upload failed: ${xhr.status}` + "`" + `));
    xhr.onerror = () => reject(new Error("Network error during upload"));
    xhr.open("PUT", presigned_url);
    xhr.setRequestHeader("Content-Type", file.type);
    xhr.send(file);
  });

  // Step 3: Record the upload in the database
  const { data: completeRes } = await apiClient.post("/api/uploads/complete", {
    key,
    filename: file.name,
    content_type: file.type,
    size: file.size,
  });
  return completeRes;
}
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
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { LoginSchema, type LoginInput } from "@repo/shared/schemas";

const inputClass = "w-full rounded-lg border border-border bg-bg-tertiary px-4 py-3 text-foreground placeholder:text-text-muted focus:border-accent focus:outline-none focus:ring-1 focus:ring-accent transition-colors";
const errorInputClass = "w-full rounded-lg border border-danger bg-bg-tertiary px-4 py-3 text-foreground placeholder:text-text-muted focus:border-danger focus:outline-none focus:ring-1 focus:ring-danger transition-colors";

export default function LoginPage() {
  const [showPassword, setShowPassword] = useState(false);
  const { mutate: login, isPending, error: serverError } = useLogin();

  const { register, handleSubmit, formState: { errors } } = useForm<LoginInput>({
    resolver: zodResolver(LoginSchema),
  });

  const onSubmit = (data: LoginInput) => {
    login(data);
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

          <form onSubmit={handleSubmit(onSubmit)} className="space-y-5">
            {serverError && (
              <div className="rounded-lg bg-danger/10 border border-danger/20 px-4 py-3 text-sm text-danger">
                {(serverError as unknown as { response?: { data?: { error?: { message?: string } } } })?.response?.data?.error?.message || "Invalid credentials"}
              </div>
            )}

            <div className="space-y-2">
              <label htmlFor="email" className="block text-sm font-medium text-text-secondary">
                Email
              </label>
              <input
                id="email"
                type="email"
                {...register("email")}
                className={errors.email ? errorInputClass : inputClass}
                placeholder="you@example.com"
                autoFocus
              />
              {errors.email && <p className="text-sm text-danger">{errors.email.message}</p>}
            </div>

            <div className="space-y-2">
              <label htmlFor="password" className="block text-sm font-medium text-text-secondary">
                Password
              </label>
              <div className="relative">
                <input
                  id="password"
                  type={showPassword ? "text" : "password"}
                  {...register("password")}
                  className={errors.password ? errorInputClass + " pr-12" : inputClass + " pr-12"}
                  placeholder="Enter your password"
                />
                <button
                  type="button"
                  onClick={() => setShowPassword(!showPassword)}
                  className="absolute right-3 top-1/2 -translate-y-1/2 text-text-muted hover:text-text-secondary transition-colors"
                >
                  {showPassword ? <EyeOff className="h-5 w-5" /> : <Eye className="h-5 w-5" />}
                </button>
              </div>
              {errors.password && <p className="text-sm text-danger">{errors.password.message}</p>}
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
` + socialLoginButtonsJSX() + `
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
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { RegisterSchema, type RegisterInput } from "@repo/shared/schemas";

const inputClass = "w-full rounded-lg border border-border bg-bg-tertiary px-4 py-3 text-foreground placeholder:text-text-muted focus:border-accent focus:outline-none focus:ring-1 focus:ring-accent transition-colors";
const errorInputClass = "w-full rounded-lg border border-danger bg-bg-tertiary px-4 py-3 text-foreground placeholder:text-text-muted focus:border-danger focus:outline-none focus:ring-1 focus:ring-danger transition-colors";

export default function SignUpPage() {
  const [showPassword, setShowPassword] = useState(false);
  const [showConfirmPassword, setShowConfirmPassword] = useState(false);
  const { mutate: registerUser, isPending, error: serverError } = useRegister();

  const { register, handleSubmit, formState: { errors } } = useForm<RegisterInput>({
    resolver: zodResolver(RegisterSchema),
  });

  const onSubmit = (data: RegisterInput) => {
    registerUser({
      first_name: data.firstName,
      last_name: data.lastName,
      email: data.email,
      password: data.password,
    });
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

          <form onSubmit={handleSubmit(onSubmit)} className="space-y-5">
            {serverError && (
              <div className="rounded-lg bg-danger/10 border border-danger/20 px-4 py-3 text-sm text-danger">
                {(serverError as unknown as { response?: { data?: { error?: { message?: string } } } })?.response?.data?.error?.message || "Registration failed"}
              </div>
            )}

            <div className="grid grid-cols-2 gap-4">
              <div className="space-y-2">
                <label htmlFor="firstName" className="block text-sm font-medium text-text-secondary">
                  First name
                </label>
                <input
                  id="firstName"
                  type="text"
                  {...register("firstName")}
                  className={errors.firstName ? errorInputClass : inputClass}
                  placeholder="John"
                  autoFocus
                />
                {errors.firstName && <p className="text-sm text-danger">{errors.firstName.message}</p>}
              </div>
              <div className="space-y-2">
                <label htmlFor="lastName" className="block text-sm font-medium text-text-secondary">
                  Last name
                </label>
                <input
                  id="lastName"
                  type="text"
                  {...register("lastName")}
                  className={errors.lastName ? errorInputClass : inputClass}
                  placeholder="Doe"
                />
                {errors.lastName && <p className="text-sm text-danger">{errors.lastName.message}</p>}
              </div>
            </div>

            <div className="space-y-2">
              <label htmlFor="email" className="block text-sm font-medium text-text-secondary">
                Email
              </label>
              <input
                id="email"
                type="email"
                {...register("email")}
                className={errors.email ? errorInputClass : inputClass}
                placeholder="you@example.com"
              />
              {errors.email && <p className="text-sm text-danger">{errors.email.message}</p>}
            </div>

            <div className="space-y-2">
              <label htmlFor="password" className="block text-sm font-medium text-text-secondary">
                Password
              </label>
              <div className="relative">
                <input
                  id="password"
                  type={showPassword ? "text" : "password"}
                  {...register("password")}
                  className={errors.password ? errorInputClass + " pr-12" : inputClass + " pr-12"}
                  placeholder="Min. 8 characters"
                />
                <button
                  type="button"
                  onClick={() => setShowPassword(!showPassword)}
                  className="absolute right-3 top-1/2 -translate-y-1/2 text-text-muted hover:text-text-secondary transition-colors"
                >
                  {showPassword ? <EyeOff className="h-5 w-5" /> : <Eye className="h-5 w-5" />}
                </button>
              </div>
              {errors.password && <p className="text-sm text-danger">{errors.password.message}</p>}
            </div>

            <div className="space-y-2">
              <label htmlFor="confirmPassword" className="block text-sm font-medium text-text-secondary">
                Confirm password
              </label>
              <div className="relative">
                <input
                  id="confirmPassword"
                  type={showConfirmPassword ? "text" : "password"}
                  {...register("confirmPassword")}
                  className={errors.confirmPassword ? errorInputClass + " pr-12" : inputClass + " pr-12"}
                  placeholder="Repeat your password"
                />
                <button
                  type="button"
                  onClick={() => setShowConfirmPassword(!showConfirmPassword)}
                  className="absolute right-3 top-1/2 -translate-y-1/2 text-text-muted hover:text-text-secondary transition-colors"
                >
                  {showConfirmPassword ? <EyeOff className="h-5 w-5" /> : <Eye className="h-5 w-5" />}
                </button>
              </div>
              {errors.confirmPassword && <p className="text-sm text-danger">{errors.confirmPassword.message}</p>}
            </div>

            <button
              type="submit"
              disabled={isPending}
              className="w-full rounded-lg bg-accent py-3 font-medium text-white hover:bg-accent-hover disabled:opacity-50 transition-colors"
            >
              {isPending ? "Creating account..." : "Create Account"}
            </button>
          </form>
` + socialLoginButtonsJSX() + `
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
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { ForgotPasswordSchema, type ForgotPasswordInput } from "@repo/shared/schemas";

const inputClass = "w-full rounded-lg border border-border bg-bg-tertiary px-4 py-3 text-foreground placeholder:text-text-muted focus:border-accent focus:outline-none focus:ring-1 focus:ring-accent transition-colors";
const errorInputClass = "w-full rounded-lg border border-danger bg-bg-tertiary px-4 py-3 text-foreground placeholder:text-text-muted focus:border-danger focus:outline-none focus:ring-1 focus:ring-danger transition-colors";

export default function ForgotPasswordPage() {
  const [sent, setSent] = useState(false);
  const [loading, setLoading] = useState(false);

  const { register, handleSubmit, formState: { errors } } = useForm<ForgotPasswordInput>({
    resolver: zodResolver(ForgotPasswordSchema),
  });

  const onSubmit = async (data: ForgotPasswordInput) => {
    setLoading(true);
    try {
      await apiClient.post("/api/auth/forgot-password", data);
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
            <form onSubmit={handleSubmit(onSubmit)} className="space-y-5">
              <div className="space-y-2">
                <label htmlFor="email" className="block text-sm font-medium text-text-secondary">
                  Email
                </label>
                <input
                  id="email"
                  type="email"
                  {...register("email")}
                  className={errors.email ? errorInputClass : inputClass}
                  placeholder="you@example.com"
                  autoFocus
                />
                {errors.email && <p className="text-sm text-danger">{errors.email.message}</p>}
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

func adminErrorPage() string {
	return `"use client";

import { useEffect } from "react";

export default function Error({
  error,
  reset,
}: {
  error: Error & { digest?: string };
  reset: () => void;
}) {
  useEffect(() => {
    console.error(error);
  }, [error]);

  return (
    <div className="flex min-h-screen items-center justify-center bg-background px-4">
      <div className="w-full max-w-md text-center">
        <div className="mx-auto mb-6 flex h-16 w-16 items-center justify-center rounded-full bg-danger/10 border border-danger/20">
          <svg className="h-8 w-8 text-danger" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
            <path strokeLinecap="round" strokeLinejoin="round" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
          </svg>
        </div>
        <h2 className="mb-2 text-2xl font-bold text-foreground">Something went wrong</h2>
        <p className="mb-6 text-text-secondary">
          An unexpected error occurred. You can try again or go back.
        </p>
        {error.digest && (
          <p className="mb-4 text-xs text-text-muted font-mono">Error ID: {error.digest}</p>
        )}
        <div className="flex gap-3 justify-center">
          <button
            onClick={() => window.history.back()}
            className="inline-flex items-center gap-2 rounded-lg border border-border bg-bg-secondary px-4 py-2.5 text-sm font-medium text-foreground hover:bg-bg-tertiary transition-colors"
          >
            Go Back
          </button>
          <button
            onClick={reset}
            className="inline-flex items-center gap-2 rounded-lg bg-accent px-4 py-2.5 text-sm font-medium text-white hover:bg-accent-hover transition-colors"
          >
            Try Again
          </button>
        </div>
      </div>
    </div>
  );
}
`
}

func adminNotFoundPage() string {
	return `import Link from "next/link";

export default function NotFound() {
  return (
    <div className="flex min-h-screen items-center justify-center bg-background px-4">
      <div className="w-full max-w-md text-center">
        <p className="mb-4 text-7xl font-bold text-accent">404</p>
        <h2 className="mb-2 text-2xl font-bold text-foreground">Page not found</h2>
        <p className="mb-8 text-text-secondary">
          The page you&apos;re looking for doesn&apos;t exist or has been moved.
        </p>
        <div className="flex gap-3 justify-center">
          <Link
            href="/dashboard"
            className="inline-flex items-center gap-2 rounded-lg bg-accent px-4 py-2.5 text-sm font-medium text-white hover:bg-accent-hover transition-colors"
          >
            Dashboard
          </Link>
        </div>
      </div>
    </div>
  );
}
`
}

func adminGlobalErrorPage() string {
	return `"use client";

export default function GlobalError({
  error,
  reset,
}: {
  error: Error & { digest?: string };
  reset: () => void;
}) {
  return (
    <html lang="en" className="dark">
      <body style={{ minHeight: "100vh", backgroundColor: "#0a0a0f", fontFamily: "system-ui, sans-serif", margin: 0 }}>
        <div style={{ display: "flex", minHeight: "100vh", alignItems: "center", justifyContent: "center", padding: "1rem" }}>
          <div style={{ maxWidth: "28rem", textAlign: "center" }}>
            <div style={{ margin: "0 auto 1.5rem", display: "flex", height: "4rem", width: "4rem", alignItems: "center", justifyContent: "center", borderRadius: "9999px", backgroundColor: "rgba(239,68,68,0.1)", border: "1px solid rgba(239,68,68,0.2)" }}>
              <svg style={{ height: "2rem", width: "2rem", color: "#f87171" }} fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
                <path strokeLinecap="round" strokeLinejoin="round" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
              </svg>
            </div>
            <h2 style={{ marginBottom: "0.5rem", fontSize: "1.5rem", fontWeight: 700, color: "#e8e8f0" }}>Application Error</h2>
            <p style={{ marginBottom: "1.5rem", color: "#9090a8" }}>
              A critical error occurred. Please try refreshing the page.
            </p>
            {error.digest && (
              <p style={{ marginBottom: "1rem", fontSize: "0.75rem", color: "#606078", fontFamily: "monospace" }}>Error ID: {error.digest}</p>
            )}
            <button
              onClick={reset}
              style={{ display: "inline-flex", alignItems: "center", gap: "0.5rem", borderRadius: "0.5rem", backgroundColor: "#6c5ce7", padding: "0.625rem 1.25rem", fontSize: "0.875rem", fontWeight: 500, color: "white", border: "none", cursor: "pointer" }}
            >
              Try Again
            </button>
          </div>
        </div>
      </body>
    </html>
  );
}
`
}
