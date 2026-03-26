package scaffold

import (
	"fmt"
	"path/filepath"
	"strings"
)

// stripUseClient removes "use client" directives from component code.
// TanStack Router apps don't need this since everything is client-side.
func stripUseClient(code string) string {
	code = strings.Replace(code, "\"use client\"\n\n", "", 1)
	code = strings.Replace(code, "'use client'\n\n", "", 1)
	code = strings.Replace(code, "\"use client\";\n\n", "", 1)
	code = strings.Replace(code, "'use client';\n\n", "", 1)
	code = strings.Replace(code, "\"use client\"\n", "", 1)
	code = strings.Replace(code, "'use client'\n", "", 1)
	return code
}

func writeAdminTanStackFiles(root string, opts Options) error {
	adminRoot := filepath.Join(root, "apps", "admin")

	files := map[string]string{
		// Config files
		filepath.Join(adminRoot, "package.json"):         adminTanStackPackageJSON(opts),
		filepath.Join(adminRoot, "vite.config.ts"):       adminTanStackViteConfig(),
		filepath.Join(adminRoot, "index.html"):           adminTanStackIndexHTML(opts),
		filepath.Join(adminRoot, "tailwind.config.ts"):   adminTailwindConfig(),
		filepath.Join(adminRoot, "postcss.config.js"):    adminPostCSSConfig(),
		filepath.Join(adminRoot, "tsconfig.json"):        adminTanStackTSConfig(),
		filepath.Join(adminRoot, "src", "main.tsx"):      adminTanStackMain(),
		filepath.Join(adminRoot, "src", "globals.css"):   adminGlobalCSS(),

		// Routes
		filepath.Join(adminRoot, "src", "routes", "__root.tsx"):                          adminTanStackRootRoute(),
		filepath.Join(adminRoot, "src", "routes", "index.tsx"):                           adminTanStackRedirectRoute(),
		filepath.Join(adminRoot, "src", "routes", "_auth.tsx"):                           adminTanStackAuthLayout(),
		filepath.Join(adminRoot, "src", "routes", "_auth", "login.tsx"):                  adminTanStackLoginRoute(opts.Style),
		filepath.Join(adminRoot, "src", "routes", "_auth", "sign-up.tsx"):                adminTanStackSignUpRoute(opts.Style),
		filepath.Join(adminRoot, "src", "routes", "_auth", "forgot-password.tsx"):        adminTanStackForgotPasswordRoute(opts.Style),
		filepath.Join(adminRoot, "src", "routes", "_dashboard.tsx"):                      adminTanStackDashboardLayout(),
		filepath.Join(adminRoot, "src", "routes", "_dashboard", "dashboard.tsx"):         adminTanStackDashboardRoute(opts.Style),
		filepath.Join(adminRoot, "src", "routes", "_dashboard", "profile.tsx"):           adminTanStackProfileRoute(),
		filepath.Join(adminRoot, "src", "routes", "_dashboard", "resources", "users.tsx"): adminTanStackUsersRoute(),
		filepath.Join(adminRoot, "src", "routes", "_dashboard", "resources", "blogs.tsx"): adminTanStackBlogsRoute(),
		filepath.Join(adminRoot, "src", "routes", "_dashboard", "system", "jobs.tsx"):     adminTanStackSystemRoute("jobs"),
		filepath.Join(adminRoot, "src", "routes", "_dashboard", "system", "files.tsx"):    adminTanStackSystemRoute("files"),
		filepath.Join(adminRoot, "src", "routes", "_dashboard", "system", "cron.tsx"):     adminTanStackSystemRoute("cron"),
		filepath.Join(adminRoot, "src", "routes", "_dashboard", "system", "mail.tsx"):     adminTanStackSystemRoute("mail"),
		filepath.Join(adminRoot, "src", "routes", "_dashboard", "system", "security.tsx"): adminTanStackSystemRoute("security"),

		// Lib (same as Next.js versions)
		filepath.Join(adminRoot, "src", "lib", "api-client.ts"):  adminAPIClient(),
		filepath.Join(adminRoot, "src", "lib", "query-client.ts"): adminQueryClient(),
		filepath.Join(adminRoot, "src", "lib", "utils.ts"):        adminUtils(),
		filepath.Join(adminRoot, "src", "lib", "resource.ts"):     adminResourceTypes(),
		filepath.Join(adminRoot, "src", "lib", "icons.ts"):        adminIconMap(),
		filepath.Join(adminRoot, "src", "lib", "formatters.ts"):   adminFormatters(),

		// Hooks (same as Next.js versions)
		filepath.Join(adminRoot, "src", "hooks", "use-auth.ts"):     adminUseAuth(),
		filepath.Join(adminRoot, "src", "hooks", "use-resource.ts"): adminUseResource(),
		filepath.Join(adminRoot, "src", "hooks", "use-system.ts"):   adminUseSystem(),
		filepath.Join(adminRoot, "src", "hooks", "use-profile.ts"):  adminUseProfile(),

		// Shared components
		filepath.Join(adminRoot, "src", "components", "shared", "providers.tsx"):      stripUseClient(adminProviders()),
		filepath.Join(adminRoot, "src", "components", "shared", "theme-provider.tsx"): stripUseClient(adminThemeProvider()),

		// Layout components (reuse with stripped "use client")
		filepath.Join(adminRoot, "src", "components", "layout", "admin-layout.tsx"): stripUseClient(adminLayoutComponent()),
		filepath.Join(adminRoot, "src", "components", "layout", "sidebar.tsx"):      stripUseClient(adminSidebar()),
		filepath.Join(adminRoot, "src", "components", "layout", "navbar.tsx"):       stripUseClient(adminNavbar()),

		// Table components (pure React — strip "use client")
		filepath.Join(adminRoot, "src", "components", "tables", "data-table.tsx"):        stripUseClient(adminDataTable()),
		filepath.Join(adminRoot, "src", "components", "tables", "column-header.tsx"):     stripUseClient(adminColumnHeader()),
		filepath.Join(adminRoot, "src", "components", "tables", "cell-renderers.tsx"):    stripUseClient(adminCellRenderers()),
		filepath.Join(adminRoot, "src", "components", "tables", "table-filters.tsx"):     stripUseClient(adminTableFilters()),
		filepath.Join(adminRoot, "src", "components", "tables", "table-toolbar.tsx"):     stripUseClient(adminTableToolbar()),
		filepath.Join(adminRoot, "src", "components", "tables", "table-pagination.tsx"):  stripUseClient(adminTablePagination()),
		filepath.Join(adminRoot, "src", "components", "tables", "table-skeleton.tsx"):    stripUseClient(adminTableSkeleton()),
		filepath.Join(adminRoot, "src", "components", "tables", "table-empty-state.tsx"): stripUseClient(adminTableEmptyState()),

		// Form components (pure React — strip "use client")
		filepath.Join(adminRoot, "src", "components", "forms", "form-builder.tsx"):      stripUseClient(adminFormBuilder()),
		filepath.Join(adminRoot, "src", "components", "forms", "form-modal.tsx"):        stripUseClient(adminFormModal()),
		filepath.Join(adminRoot, "src", "components", "forms", "form-page.tsx"):         stripUseClient(adminFormPage()),
		filepath.Join(adminRoot, "src", "components", "forms", "form-stepper.tsx"):      stripUseClient(adminFormStepper()),
		filepath.Join(adminRoot, "src", "components", "forms", "form-modal-steps.tsx"):  stripUseClient(adminFormModalSteps()),
		filepath.Join(adminRoot, "src", "components", "forms", "form-page-steps.tsx"):   stripUseClient(adminFormPageSteps()),
		filepath.Join(adminRoot, "src", "components", "forms", "fields", "text-field.tsx"):                      stripUseClient(adminTextField()),
		filepath.Join(adminRoot, "src", "components", "forms", "fields", "textarea-field.tsx"):                  stripUseClient(adminTextareaField()),
		filepath.Join(adminRoot, "src", "components", "forms", "fields", "number-field.tsx"):                    stripUseClient(adminNumberField()),
		filepath.Join(adminRoot, "src", "components", "forms", "fields", "select-field.tsx"):                    stripUseClient(adminSelectField()),
		filepath.Join(adminRoot, "src", "components", "forms", "fields", "date-field.tsx"):                      stripUseClient(adminDateField()),
		filepath.Join(adminRoot, "src", "components", "forms", "fields", "toggle-field.tsx"):                    stripUseClient(adminToggleField()),
		filepath.Join(adminRoot, "src", "components", "forms", "fields", "checkbox-field.tsx"):                  stripUseClient(adminCheckboxField()),
		filepath.Join(adminRoot, "src", "components", "forms", "fields", "radio-field.tsx"):                     stripUseClient(adminRadioField()),
		filepath.Join(adminRoot, "src", "components", "forms", "fields", "image-field.tsx"):                     stripUseClient(adminImageField()),
		filepath.Join(adminRoot, "src", "components", "forms", "fields", "images-field.tsx"):                    stripUseClient(adminImagesField()),
		filepath.Join(adminRoot, "src", "components", "forms", "fields", "video-field.tsx"):                     stripUseClient(adminVideoField()),
		filepath.Join(adminRoot, "src", "components", "forms", "fields", "videos-field.tsx"):                    stripUseClient(adminVideosField()),
		filepath.Join(adminRoot, "src", "components", "forms", "fields", "file-field.tsx"):                      stripUseClient(adminFileField()),
		filepath.Join(adminRoot, "src", "components", "forms", "fields", "files-field.tsx"):                     stripUseClient(adminFilesField()),
		filepath.Join(adminRoot, "src", "components", "forms", "fields", "relationship-select-field.tsx"):       stripUseClient(adminRelationshipSelectField()),
		filepath.Join(adminRoot, "src", "components", "forms", "fields", "multi-relationship-select-field.tsx"): stripUseClient(adminMultiRelationshipSelectField()),
		filepath.Join(adminRoot, "src", "components", "forms", "fields", "rich-text-field.tsx"):                 stripUseClient(adminRichTextField()),

		// UI components
		filepath.Join(adminRoot, "src", "components", "ui", "dropzone.tsx"):      stripUseClient(adminDropzone()),
		filepath.Join(adminRoot, "src", "components", "ui", "confirm-modal.tsx"): stripUseClient(adminConfirmModal()),

		// Widget components
		filepath.Join(adminRoot, "src", "components", "widgets", "stats-card.tsx"):      stripUseClient(adminStatsCard()),
		filepath.Join(adminRoot, "src", "components", "widgets", "chart-widget.tsx"):    stripUseClient(adminChartWidget()),
		filepath.Join(adminRoot, "src", "components", "widgets", "activity-widget.tsx"): stripUseClient(adminActivityWidget()),
		filepath.Join(adminRoot, "src", "components", "widgets", "widget-grid.tsx"):     stripUseClient(adminWidgetGrid()),

		// Resource components
		filepath.Join(adminRoot, "src", "components", "resource", "resource-page.tsx"): stripUseClient(adminResourcePage()),
		filepath.Join(adminRoot, "src", "components", "resource", "view-modal.tsx"):    stripUseClient(adminViewModal()),

		// Resource definitions (same as Next.js)
		filepath.Join(adminRoot, "src", "resources", "index.ts"): adminResourceRegistry(),
		filepath.Join(adminRoot, "src", "resources", "users.ts"): adminUsersResource(),
		filepath.Join(adminRoot, "src", "resources", "blogs.ts"): adminBlogsResource(),

		// Profile
		filepath.Join(adminRoot, "src", "components", "profile", "delete-account-dialog.tsx"): stripUseClient(adminDeleteAccountDialog()),

		filepath.Join(adminRoot, "public", ".gitkeep"): "",
	}

	for path, content := range files {
		if err := writeFile(path, content); err != nil {
			return fmt.Errorf("writing %s: %w", path, err)
		}
	}

	return nil
}

func adminTanStackPackageJSON(opts Options) string {
	return fmt.Sprintf(`{
  "name": "@%s/admin",
  "version": "0.1.0",
  "private": true,
  "type": "module",
  "scripts": {
    "dev": "vite --port 3001",
    "build": "tsc -b && vite build",
    "preview": "vite preview",
    "lint": "eslint .",
    "test": "vitest run",
    "test:watch": "vitest"
  },
  "dependencies": {
    "@hookform/resolvers": "^3.9.1",
    "@tanstack/react-query": "^5.62.0",
    "@tanstack/react-router": "^1.93.0",
    "@tanstack/react-table": "^8.20.6",
    "axios": "^1.7.9",
    "clsx": "^2.1.1",
    "lucide-react": "^0.468.0",
    "react": "^19.0.0",
    "react-dom": "^19.0.0",
    "react-hook-form": "^7.54.1",
    "recharts": "^2.15.0",
    "tailwind-merge": "^2.6.0",
    "zod": "^3.24.1"
  },
  "devDependencies": {
    "@tanstack/react-router-devtools": "^1.93.0",
    "@tanstack/router-vite-plugin": "^1.93.0",
    "@types/react": "^19.0.0",
    "@types/react-dom": "^19.0.0",
    "autoprefixer": "^10.4.20",
    "postcss": "^8.4.49",
    "tailwindcss": "^3.4.17",
    "typescript": "~5.7.0",
    "vite": "^6.0.0",
    "@vitejs/plugin-react": "^4.3.4",
    "vitest": "^2.1.0",
    "@testing-library/react": "^16.1.0"
  }
}`, opts.ProjectName)
}

func adminTanStackViteConfig() string {
	return `import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import { TanStackRouterVite } from '@tanstack/router-vite-plugin'
import path from 'path'

export default defineConfig({
  plugins: [
    TanStackRouterVite(),
    react(),
  ],
  resolve: {
    alias: {
      '@': path.resolve(__dirname, './src'),
    },
  },
  server: {
    port: 3001,
    proxy: {
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      },
    },
  },
})
`
}

func adminTanStackIndexHTML(opts Options) string {
	return fmt.Sprintf(`<!DOCTYPE html>
<html lang="en" class="dark">
  <head>
    <meta charset="UTF-8" />
    <link rel="icon" type="image/svg+xml" href="/vite.svg" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <title>%s — Admin</title>
  </head>
  <body class="min-h-screen bg-background text-foreground antialiased">
    <div id="root"></div>
    <script type="module" src="/src/main.tsx"></script>
  </body>
</html>
`, opts.ProjectName)
}

func adminTanStackTSConfig() string {
	return webTanStackTSConfig() // Same config
}

func adminTanStackMain() string {
	return `import React from 'react'
import ReactDOM from 'react-dom/client'
import { RouterProvider, createRouter } from '@tanstack/react-router'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { routeTree } from './routeTree.gen'
import './globals.css'

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      staleTime: 5 * 60 * 1000,
      retry: 1,
    },
  },
})

const router = createRouter({ routeTree })

declare module '@tanstack/react-router' {
  interface Register {
    router: typeof router
  }
}

ReactDOM.createRoot(document.getElementById('root')!).render(
  <React.StrictMode>
    <QueryClientProvider client={queryClient}>
      <RouterProvider router={router} />
    </QueryClientProvider>
  </React.StrictMode>,
)
`
}

func adminTanStackRootRoute() string {
	return `import { createRootRoute, Outlet } from '@tanstack/react-router'

export const Route = createRootRoute({
  component: () => <Outlet />,
})
`
}

func adminTanStackRedirectRoute() string {
	return `import { createFileRoute, redirect } from '@tanstack/react-router'

export const Route = createFileRoute('/')({
  beforeLoad: () => {
    throw redirect({ to: '/dashboard' })
  },
})
`
}

func adminTanStackAuthLayout() string {
	return `import { createFileRoute, Outlet, redirect } from '@tanstack/react-router'

export const Route = createFileRoute('/_auth')({
  beforeLoad: () => {
    const token = localStorage.getItem('access_token')
    if (token) {
      throw redirect({ to: '/dashboard' })
    }
  },
  component: () => (
    <div className="min-h-screen flex items-center justify-center bg-background">
      <Outlet />
    </div>
  ),
})
`
}

func adminTanStackLoginRoute(style string) string {
	// Reuse the style-specific login page but wrap in TanStack route
	return `import { createFileRoute, Link, useNavigate } from '@tanstack/react-router'
import { useState } from 'react'
import { useForm } from 'react-hook-form'
import { api } from '@/lib/api-client'

export const Route = createFileRoute('/_auth/login')({
  component: LoginPage,
})

function LoginPage() {
  const navigate = useNavigate()
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)
  const { register, handleSubmit, formState: { errors } } = useForm()

  const onSubmit = async (data: any) => {
    setLoading(true)
    setError('')
    try {
      const res = await api.post('/api/auth/login', data)
      if (res.data.data?.totp_required) {
        // 2FA required — redirect to TOTP page (store pending token)
        localStorage.setItem('totp_pending', res.data.data.pending_token)
        // TODO: navigate to TOTP verification page
        return
      }
      localStorage.setItem('access_token', res.data.data.tokens.access_token)
      localStorage.setItem('refresh_token', res.data.data.tokens.refresh_token)
      navigate({ to: '/dashboard' })
    } catch (err: any) {
      setError(err.response?.data?.error?.message || 'Login failed')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="w-full max-w-md mx-auto p-8">
      <div className="text-center mb-8">
        <h1 className="text-3xl font-bold">Welcome back</h1>
        <p className="text-muted-foreground mt-2">Sign in to your account</p>
      </div>

      <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
        {error && (
          <div className="p-3 rounded-lg bg-destructive/10 border border-destructive/20 text-sm text-destructive">
            {error}
          </div>
        )}

        <div>
          <label className="block text-sm font-medium mb-1.5">Email</label>
          <input
            type="email"
            {...register('email', { required: 'Email is required' })}
            className="w-full px-3 py-2 rounded-lg border border-border bg-background text-foreground focus:outline-none focus:ring-2 focus:ring-accent"
            placeholder="you@example.com"
          />
          {errors.email && <p className="text-xs text-destructive mt-1">{String(errors.email.message)}</p>}
        </div>

        <div>
          <label className="block text-sm font-medium mb-1.5">Password</label>
          <input
            type="password"
            {...register('password', { required: 'Password is required' })}
            className="w-full px-3 py-2 rounded-lg border border-border bg-background text-foreground focus:outline-none focus:ring-2 focus:ring-accent"
            placeholder="••••••••"
          />
          {errors.password && <p className="text-xs text-destructive mt-1">{String(errors.password.message)}</p>}
        </div>

        <div className="flex items-center justify-between">
          <Link to="/forgot-password" className="text-sm text-accent hover:underline">
            Forgot password?
          </Link>
        </div>

        <button
          type="submit"
          disabled={loading}
          className="w-full py-2.5 rounded-lg bg-accent text-white font-medium hover:bg-accent-hover transition-colors disabled:opacity-50"
        >
          {loading ? 'Signing in...' : 'Sign in'}
        </button>

        <p className="text-center text-sm text-muted-foreground">
          Don{"'"}t have an account?{' '}
          <Link to="/sign-up" className="text-accent hover:underline">Sign up</Link>
        </p>
      </form>
    </div>
  )
}
`
}

func adminTanStackSignUpRoute(style string) string {
	return `import { createFileRoute, Link, useNavigate } from '@tanstack/react-router'
import { useState } from 'react'
import { useForm } from 'react-hook-form'
import { api } from '@/lib/api-client'

export const Route = createFileRoute('/_auth/sign-up')({
  component: SignUpPage,
})

function SignUpPage() {
  const navigate = useNavigate()
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)
  const { register, handleSubmit, formState: { errors } } = useForm()

  const onSubmit = async (data: any) => {
    setLoading(true)
    setError('')
    try {
      const res = await api.post('/api/auth/register', data)
      localStorage.setItem('access_token', res.data.data.tokens.access_token)
      localStorage.setItem('refresh_token', res.data.data.tokens.refresh_token)
      navigate({ to: '/dashboard' })
    } catch (err: any) {
      setError(err.response?.data?.error?.message || 'Registration failed')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="w-full max-w-md mx-auto p-8">
      <div className="text-center mb-8">
        <h1 className="text-3xl font-bold">Create account</h1>
        <p className="text-muted-foreground mt-2">Get started with your new account</p>
      </div>

      <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
        {error && (
          <div className="p-3 rounded-lg bg-destructive/10 border border-destructive/20 text-sm text-destructive">
            {error}
          </div>
        )}

        <div className="grid grid-cols-2 gap-4">
          <div>
            <label className="block text-sm font-medium mb-1.5">First Name</label>
            <input type="text" {...register('first_name', { required: 'Required' })} className="w-full px-3 py-2 rounded-lg border border-border bg-background text-foreground focus:outline-none focus:ring-2 focus:ring-accent" />
          </div>
          <div>
            <label className="block text-sm font-medium mb-1.5">Last Name</label>
            <input type="text" {...register('last_name', { required: 'Required' })} className="w-full px-3 py-2 rounded-lg border border-border bg-background text-foreground focus:outline-none focus:ring-2 focus:ring-accent" />
          </div>
        </div>

        <div>
          <label className="block text-sm font-medium mb-1.5">Email</label>
          <input type="email" {...register('email', { required: 'Email is required' })} className="w-full px-3 py-2 rounded-lg border border-border bg-background text-foreground focus:outline-none focus:ring-2 focus:ring-accent" placeholder="you@example.com" />
        </div>

        <div>
          <label className="block text-sm font-medium mb-1.5">Password</label>
          <input type="password" {...register('password', { required: 'Password is required', minLength: { value: 8, message: 'Min 8 characters' } })} className="w-full px-3 py-2 rounded-lg border border-border bg-background text-foreground focus:outline-none focus:ring-2 focus:ring-accent" placeholder="••••••••" />
          {errors.password && <p className="text-xs text-destructive mt-1">{String(errors.password.message)}</p>}
        </div>

        <button type="submit" disabled={loading} className="w-full py-2.5 rounded-lg bg-accent text-white font-medium hover:bg-accent-hover transition-colors disabled:opacity-50">
          {loading ? 'Creating account...' : 'Create account'}
        </button>

        <p className="text-center text-sm text-muted-foreground">
          Already have an account? <Link to="/login" className="text-accent hover:underline">Sign in</Link>
        </p>
      </form>
    </div>
  )
}
`
}

func adminTanStackForgotPasswordRoute(style string) string {
	return `import { createFileRoute, Link } from '@tanstack/react-router'
import { useState } from 'react'
import { useForm } from 'react-hook-form'
import { api } from '@/lib/api-client'

export const Route = createFileRoute('/_auth/forgot-password')({
  component: ForgotPasswordPage,
})

function ForgotPasswordPage() {
  const [sent, setSent] = useState(false)
  const [loading, setLoading] = useState(false)
  const { register, handleSubmit } = useForm()

  const onSubmit = async (data: any) => {
    setLoading(true)
    try {
      await api.post('/api/auth/forgot-password', data)
      setSent(true)
    } catch {} finally {
      setLoading(false)
    }
  }

  if (sent) {
    return (
      <div className="w-full max-w-md mx-auto p-8 text-center">
        <h1 className="text-3xl font-bold mb-4">Check your email</h1>
        <p className="text-muted-foreground mb-6">We sent a password reset link to your email address.</p>
        <Link to="/login" className="text-accent hover:underline">Back to login</Link>
      </div>
    )
  }

  return (
    <div className="w-full max-w-md mx-auto p-8">
      <div className="text-center mb-8">
        <h1 className="text-3xl font-bold">Reset password</h1>
        <p className="text-muted-foreground mt-2">Enter your email to receive a reset link</p>
      </div>
      <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
        <div>
          <label className="block text-sm font-medium mb-1.5">Email</label>
          <input type="email" {...register('email', { required: true })} className="w-full px-3 py-2 rounded-lg border border-border bg-background text-foreground focus:outline-none focus:ring-2 focus:ring-accent" placeholder="you@example.com" />
        </div>
        <button type="submit" disabled={loading} className="w-full py-2.5 rounded-lg bg-accent text-white font-medium hover:bg-accent-hover transition-colors disabled:opacity-50">
          {loading ? 'Sending...' : 'Send reset link'}
        </button>
        <p className="text-center text-sm text-muted-foreground">
          <Link to="/login" className="text-accent hover:underline">Back to login</Link>
        </p>
      </form>
    </div>
  )
}
`
}

func adminTanStackDashboardLayout() string {
	return `import { createFileRoute, Outlet, redirect } from '@tanstack/react-router'
import { AdminLayout } from '@/components/layout/admin-layout'

export const Route = createFileRoute('/_dashboard')({
  beforeLoad: () => {
    const token = localStorage.getItem('access_token')
    if (!token) {
      throw redirect({ to: '/login' })
    }
  },
  component: () => (
    <AdminLayout>
      <Outlet />
    </AdminLayout>
  ),
})
`
}

func adminTanStackDashboardRoute(style string) string {
	return `import { createFileRoute } from '@tanstack/react-router'
import { StatsCard } from '@/components/widgets/stats-card'
import { useAuth } from '@/hooks/use-auth'
import { Users, FileText, Upload, Activity } from 'lucide-react'

export const Route = createFileRoute('/_dashboard/dashboard')({
  component: DashboardPage,
})

function DashboardPage() {
  const { user } = useAuth()

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold">Dashboard</h1>
        <p className="text-muted-foreground">Welcome back, {user?.first_name || 'Admin'}</p>
      </div>

      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
        <StatsCard title="Total Users" value="--" icon={<Users className="h-5 w-5" />} />
        <StatsCard title="Blog Posts" value="--" icon={<FileText className="h-5 w-5" />} />
        <StatsCard title="Uploads" value="--" icon={<Upload className="h-5 w-5" />} />
        <StatsCard title="Active Sessions" value="--" icon={<Activity className="h-5 w-5" />} />
      </div>
    </div>
  )
}
`
}

func adminTanStackProfileRoute() string {
	return `import { createFileRoute } from '@tanstack/react-router'
import { useAuth } from '@/hooks/use-auth'

export const Route = createFileRoute('/_dashboard/profile')({
  component: ProfilePage,
})

function ProfilePage() {
  const { user } = useAuth()

  return (
    <div className="max-w-2xl">
      <h1 className="text-2xl font-bold mb-6">Profile</h1>
      <div className="rounded-xl border border-border/40 bg-card/50 p-6 space-y-4">
        <div>
          <label className="text-sm text-muted-foreground">Name</label>
          <p className="font-medium">{user?.first_name} {user?.last_name}</p>
        </div>
        <div>
          <label className="text-sm text-muted-foreground">Email</label>
          <p className="font-medium">{user?.email}</p>
        </div>
        <div>
          <label className="text-sm text-muted-foreground">Role</label>
          <p className="font-medium">{user?.role}</p>
        </div>
      </div>
    </div>
  )
}
`
}

func adminTanStackUsersRoute() string {
	return `import { createFileRoute } from '@tanstack/react-router'
import { ResourcePage } from '@/components/resource/resource-page'
import { usersResource } from '@/resources/users'

export const Route = createFileRoute('/_dashboard/resources/users')({
  component: () => <ResourcePage resource={usersResource} />,
})
`
}

func adminTanStackBlogsRoute() string {
	return `import { createFileRoute } from '@tanstack/react-router'
import { ResourcePage } from '@/components/resource/resource-page'
import { blogsResource } from '@/resources/blogs'

export const Route = createFileRoute('/_dashboard/resources/blogs')({
  component: () => <ResourcePage resource={blogsResource} />,
})
`
}

func adminTanStackSystemRoute(page string) string {
	titles := map[string]string{
		"jobs":     "Background Jobs",
		"files":    "File Manager",
		"cron":     "Cron Scheduler",
		"mail":     "Mail Preview",
		"security": "Security",
	}
	title := titles[page]
	if title == "" {
		title = page
	}

	return fmt.Sprintf(`import { createFileRoute } from '@tanstack/react-router'

export const Route = createFileRoute('/_dashboard/system/%s')({
  component: %sPage,
})

function %sPage() {
  return (
    <div>
      <h1 className="text-2xl font-bold mb-6">%s</h1>
      <div className="rounded-xl border border-border/40 bg-card/50 p-6">
        <p className="text-muted-foreground">System page content will be loaded here.</p>
      </div>
    </div>
  )
}
`, page, strings.Title(page), strings.Title(page), title)
}
