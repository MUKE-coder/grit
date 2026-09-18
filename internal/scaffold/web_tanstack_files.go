package scaffold

import (
	"fmt"
	"path/filepath"
)

func writeWebTanStackFiles(root string, opts Options) error {
	webRoot := filepath.Join(root, "apps", "web")

	files := map[string]string{
		filepath.Join(webRoot, "package.json"):   webTanStackPackageJSON(opts),
		filepath.Join(webRoot, "vite.config.ts"): webTanStackViteConfig(opts),
		filepath.Join(webRoot, "index.html"):     webTanStackIndexHTML(opts),
		// .cjs (not .js) because package.json sets "type": "module" and PostCSS
		// config still uses CommonJS module.exports.
		filepath.Join(webRoot, "postcss.config.cjs"):          postCSSConfigFor(webRoot),
		filepath.Join(webRoot, "src", "vite-env.d.ts"):        viteEnvTypes(),
		filepath.Join(webRoot, "tsconfig.json"):               webTanStackTSConfig(opts),
		filepath.Join(webRoot, "src", "main.tsx"):             webTanStackMain(),
		filepath.Join(webRoot, "src", "globals.css"):          webGlobalCSS(),
		filepath.Join(webRoot, "src", "routes", "__root.tsx"): webTanStackRootRoute(opts),
		// The public site, as a pathless layout route: the URLs are still / and
		// /blog, and _site.tsx draws the navbar and footer, so a section with its
		// own chrome (the admin panel, the customer area) cannot inherit theirs.
		filepath.Join(webRoot, "src", "routes", "_site.tsx"):                  singleSiteLayoutRoute(),
		filepath.Join(webRoot, "src", "routes", "_site", "index.tsx"):         siteRouteID(webTanStackIndexRoute(opts)),
		filepath.Join(webRoot, "src", "routes", "_site", "blog", "index.tsx"): siteRouteID(webTanStackBlogListRoute()),
		filepath.Join(webRoot, "src", "routes", "_site", "blog", "$slug.tsx"): siteRouteID(webTanStackBlogDetailRoute()),
		filepath.Join(webRoot, "src", "components", "navbar.tsx"):             nextToTanStack(webNavbar(opts)),
		filepath.Join(webRoot, "src", "components", "footer.tsx"):             nextToTanStack(webFooter(opts)),
		filepath.Join(webRoot, "src", "components", "providers.tsx"):          webTanStackProviders(),
		filepath.Join(webRoot, "src", "lib", "next-compat.tsx"):               adminNextCompatShim(),
		filepath.Join(webRoot, "src", "lib", "utils.ts"):                      webUtils(),
		filepath.Join(webRoot, "components.json"):                             viteComponentsJSON(),
		filepath.Join(webRoot, "src", "lib", "api.ts"):                        viteAPIClient(),
		filepath.Join(webRoot, "src", "hooks", "use-blogs.ts"):                nextToTanStack(webUseBlogsHook()),
		filepath.Join(webRoot, "public", ".gitkeep"):                          "",
	}

	for path, content := range files {
		// adminHref fills in where the panel lives, as writeWebFiles does for the
		// Next.js app. Without it the navbar of every --vite web app carried the
		// literal {{ADMIN_HREF}}: a link to a page called that.
		if err := writeFile(path, adminHref(content, opts)); err != nil {
			return fmt.Errorf("writing %s: %w", path, err)
		}
	}

	return nil
}

func webTanStackPackageJSON(opts Options) string {
	return fmt.Sprintf(`{
  "name": "@%s/web",
  "version": "0.1.0",
  "private": true,
  "type": "module",
  "scripts": {
    "dev": "vite",
    "build": "vite build",
    "typecheck": "tsc -b",
    "preview": "vite preview",
    "lint": "eslint ."
  },
  "dependencies": {
    "@tanstack/react-query": "^5.62.0",
    "@tanstack/react-router": "^1.93.0",
    "axios": "^1.7.9",
    "clsx": "^2.1.1",
    "dompurify": "^3.4.15",
    "lucide-react": "^0.468.0",
    "react": "19.2.7",
    "react-dom": "19.2.7",
    "tailwind-merge": "^2.6.0",
    "@repo/shared": "workspace:*",
    "@repo/upload": "workspace:*"`+webAdminDependencies(opts)+viteHostDependencies(opts)+`
  },
  "devDependencies": {
    "vitest": "^2.1.0",
    "jsdom": "^25.0.0",
    "@testing-library/react": "^16.1.0",
    "@testing-library/jest-dom": "^6.4.0",
    "@testing-library/user-event": "^14.5.0",
    "@tanstack/react-router-devtools": "^1.93.0",
    "@tanstack/router-vite-plugin": "^1.93.0",
    "@types/react": "^19.0.0",
    "@types/node": "^22.0.0",
    "@types/react-dom": "^19.0.0",
    "@tailwindcss/postcss": "^4.1.13",
    "tw-animate-css": "^1.4.0",
    "postcss": "^8.4.49",
    "tailwindcss": "^4.1.13",
    "typescript": "~5.7.0",
    "vite": "^6.0.0",
    "@vitejs/plugin-react": "^4.3.4"
  }
}`, opts.ProjectName)
}

func webTanStackViteConfig(opts Options) string {
	return `import { defineConfig, loadEnv } from 'vite'
import react from '@vitejs/plugin-react'
import { TanStackRouterVite } from '@tanstack/router-vite-plugin'
import path from 'path'
` + viteSecurityHeaders() + `
export default defineConfig({
  plugins: [
    TanStackRouterVite(),
    react(),
  ],
  resolve: {
    alias: {
      '@': path.resolve(__dirname, './src'),` + viteHostAliases(opts) + `
    },
  },
  preview: {
    headers: securityHeaders,
  },
  server: {
    headers: securityHeaders,
    port: 3000,
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

func webTanStackIndexHTML(opts Options) string {
	return fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <link rel="icon" type="image/svg+xml" href="/vite.svg" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <title>%s</title>
  </head>
  <body class="min-h-screen bg-background text-foreground antialiased">
    <div id="root"></div>
    <script type="module" src="/src/main.tsx"></script>
  </body>
</html>
`, opts.ProjectName)
}

func webTanStackTailwindConfig() string {
	return `import type { Config } from 'tailwindcss'

const config: Config = {
  darkMode: 'class',
  content: [
    './index.html',
    './src/**/*.{ts,tsx}',
  ],
  theme: {
    extend: {
      colors: {
        background: '#0a0a0f',
        foreground: '#e8e8f0',
        border: '#2a2a3a',
        accent: {
          DEFAULT: '#6c5ce7',
          hover: '#7c6cf7',
        },
        muted: {
          DEFAULT: '#1a1a24',
          foreground: '#9090a8',
        },
        card: {
          DEFAULT: '#22222e',
          foreground: '#e8e8f0',
        },
        destructive: {
          DEFAULT: '#ff6b6b',
          foreground: '#e8e8f0',
        },
        success: '#00b894',
        warning: '#fdcb6e',
        info: '#74b9ff',
      },
      fontFamily: {
        sans: ['Onest', 'system-ui', 'sans-serif'],
        mono: ['JetBrains Mono', 'monospace'],
      },
    },
  },
  plugins: [],
}

export default config
`
}

func webTanStackTSConfig(opts Options) string {
	return `{
  "compilerOptions": {
    "target": "ES2020",
    "useDefineForClassFields": true,
    "lib": ["ES2020", "DOM", "DOM.Iterable"],
    "module": "ESNext",
    "skipLibCheck": true,
    "moduleResolution": "bundler",
    "allowImportingTsExtensions": true,
    "isolatedModules": true,
    "moduleDetection": "force",
    "noEmit": true,
    "jsx": "react-jsx",
    "strict": true,
    "noUnusedLocals": true,
    "noUnusedParameters": true,
    "noFallthroughCasesInSwitch": true,
    "paths": {
      "@/*": ["./src/*"]` + hostTSConfigPaths(opts) + `
    },
    "baseUrl": "."
  },
  "include": ["src", "src/vite-env.d.ts"]
}
`
}

func webTanStackMain() string {
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

// webTanStackRootRoute is the document's outlet and nothing else.
//
// Every section of a Vite app owns its layout: routes/_site draws the navbar and
// footer, routes/admin the panel's chrome, routes/_auth and routes/account their
// own. A root that drew the site chrome would put it above all of them, which is
// what shipped the admin panel inside the marketing navbar.
func webTanStackRootRoute(opts Options) string {
	_ = opts
	return `import { createRootRoute, Outlet } from '@tanstack/react-router'

export const Route = createRootRoute({
  component: () => <Outlet />,
})
`
}

func webTanStackIndexRoute(opts Options) string {
	return fmt.Sprintf(`import { createFileRoute, Link } from '@tanstack/react-router'

export const Route = createFileRoute('/')({
  component: HomePage,
})

function HomePage() {
  return (
    <div className="relative">
      {/* Hero */}
      <section className="relative py-24 px-6">
        <div className="max-w-5xl mx-auto text-center">
          <span className="inline-flex items-center rounded-full bg-accent/10 border border-accent/20 px-4 py-1.5 text-sm font-medium text-accent mb-6">
            Go + React. Built with Grit.
          </span>
          <h1 className="text-5xl md:text-7xl font-bold tracking-tight mb-6">
            <span className="text-foreground">Build faster with</span>{' '}
            <span className="bg-gradient-to-r from-accent to-purple-400 bg-clip-text text-transparent">
              %s
            </span>
          </h1>
          <p className="text-xl text-muted-foreground max-w-2xl mx-auto mb-10 leading-relaxed">
            A full-stack application powered by Go, React, and the Grit framework.
            Production-ready from day one.
          </p>
          <div className="flex items-center justify-center gap-4">
            <Link
              to="/blog"
              className="inline-flex items-center px-6 py-3 rounded-lg bg-accent text-white font-medium hover:bg-accent-hover transition-colors"
            >
              Read the Blog
            </Link>
            <a
              href="/api/health"
              target="_blank"
              rel="noopener noreferrer"
              className="inline-flex items-center px-6 py-3 rounded-lg border border-border text-foreground font-medium hover:bg-muted transition-colors"
            >
              API Health Check
            </a>
          </div>
        </div>
      </section>

      {/* Features */}
      <section className="py-20 px-6 border-t border-border/30">
        <div className="max-w-5xl mx-auto">
          <h2 className="text-3xl font-bold text-center mb-12">What{"'"}s Included</h2>
          <div className="grid md:grid-cols-3 gap-6">
            {[
              { title: 'Go API', desc: 'Gin + GORM with JWT auth, RBAC, file uploads, and background jobs.' },
              { title: 'React Frontend', desc: 'TanStack Router + React Query + Tailwind CSS. Fast and lightweight.' },
              { title: 'Full-Stack DX', desc: 'Shared types, one-command resource generation, hot reload everywhere.' },
            ].map((f) => (
              <div key={f.title} className="rounded-xl border border-border/40 bg-card/50 p-6">
                <h3 className="text-lg font-semibold mb-2">{f.title}</h3>
                <p className="text-sm text-muted-foreground leading-relaxed">{f.desc}</p>
              </div>
            ))}
          </div>
        </div>
      </section>
    </div>
  )
}
`, opts.ProjectName)
}

func webTanStackBlogListRoute() string {
	return `import { createFileRoute, Link } from '@tanstack/react-router'
import { useBlogs } from '@/hooks/use-blogs'

export const Route = createFileRoute('/blog/')({
  component: BlogListPage,
})

function BlogListPage() {
  const { data: blogs, isLoading } = useBlogs()

  return (
    <div className="max-w-4xl mx-auto py-16 px-6">
      <h1 className="text-4xl font-bold mb-2">Blog</h1>
      <p className="text-muted-foreground mb-10">Latest articles and updates.</p>

      {isLoading ? (
        <div className="space-y-4">
          {[1, 2, 3].map((i) => (
            <div key={i} className="h-32 rounded-xl bg-card/50 animate-pulse" />
          ))}
        </div>
      ) : !blogs?.length ? (
        <p className="text-muted-foreground">No posts yet. Check back soon!</p>
      ) : (
        <div className="space-y-6">
          {blogs.map((blog) => (
            <Link
              key={blog.id}
              to="/blog/$slug"
              params={{ slug: blog.slug }}
              className="block rounded-xl border border-border/40 bg-card/50 p-6 hover:border-accent/30 transition-colors"
            >
              <h2 className="text-xl font-semibold mb-2">{blog.title}</h2>
              <p className="text-sm text-muted-foreground line-clamp-2">{blog.excerpt || blog.content?.substring(0, 150)}</p>
              <span className="text-xs text-muted-foreground/50 mt-3 block">
                {new Date(blog.created_at).toLocaleDateString()}
              </span>
            </Link>
          ))}
        </div>
      )}
    </div>
  )
}
`
}

func webTanStackBlogDetailRoute() string {
	return `import { createFileRoute, Link } from '@tanstack/react-router'
import { useQuery } from '@tanstack/react-query'
import { api } from '@/lib/api'
import { ArrowLeft } from 'lucide-react'
import DOMPurify from 'dompurify'

export const Route = createFileRoute('/blog/$slug')({
  component: BlogDetailPage,
})

function BlogDetailPage() {
  const { slug } = Route.useParams()
  const { data: blog, isLoading, error } = useQuery({
    queryKey: ['blog', slug],
    queryFn: async () => {
      const res = await api.get("/api/blogs/" + slug)
      return res.data.data
    },
  })

  if (isLoading) {
    return (
      <div className="max-w-3xl mx-auto py-16 px-6">
        <div className="h-8 w-48 bg-card/50 animate-pulse rounded mb-4" />
        <div className="h-4 w-full bg-card/50 animate-pulse rounded mb-2" />
        <div className="h-4 w-3/4 bg-card/50 animate-pulse rounded" />
      </div>
    )
  }

  if (error || !blog) {
    return (
      <div className="max-w-3xl mx-auto py-16 px-6 text-center">
        <h1 className="text-2xl font-bold mb-4">Post not found</h1>
        <Link to="/blog" className="text-accent hover:underline">Back to blog</Link>
      </div>
    )
  }

  return (
    <div className="max-w-3xl mx-auto py-16 px-6">
      <Link to="/blog" className="inline-flex items-center gap-1.5 text-sm text-muted-foreground hover:text-foreground mb-8">
        <ArrowLeft className="h-4 w-4" /> Back to blog
      </Link>
      <h1 className="text-4xl font-bold mb-4">{blog.title}</h1>
      <span className="text-sm text-muted-foreground/50 block mb-8">
        {new Date(blog.created_at).toLocaleDateString()}
      </span>
      <div className="prose prose-invert max-w-none" dangerouslySetInnerHTML={{ __html: DOMPurify.sanitize(blog.content, { USE_PROFILES: { html: true }, ADD_ATTR: ['target'] }) }} />
    </div>
  )
}
`
}

func webTanStackProviders() string {
	return `import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { ReactNode, useState } from 'react'

export function Providers({ children }: { children: ReactNode }) {
  const [queryClient] = useState(() => new QueryClient({
    defaultOptions: {
      queries: {
        staleTime: 5 * 60 * 1000,
        retry: 1,
      },
    },
  }))

  return (
    <QueryClientProvider client={queryClient}>
      {children}
    </QueryClientProvider>
  )
}
`
}

// viteEnvTypes emits src/vite-env.d.ts.
//
// Without the vite/client reference, import.meta.env is untyped and any file
// reading it fails to compile with "Property 'env' does not exist on type
// 'ImportMeta'". The api client reads VITE_API_URL, so this is required, not
// cosmetic.
func viteEnvTypes() string {
	return `/// <reference types="vite/client" />

interface ImportMetaEnv {
	readonly VITE_API_URL?: string
	readonly VITE_THEME?: string
}

interface ImportMeta {
	readonly env: ImportMetaEnv
}
`
}

// viteHostAliases are the extra aliases a Vite app needs when the admin panel
// lives inside it.
//
// Vite resolves imports itself and does not read tsconfig paths, so both files
// have to be told. @repo/upload is aliased to its source because the package
// ships raw TypeScript: the Next.js apps list it in transpilePackages for the
// same reason.
func viteHostAliases(opts Options) string {
	if !opts.ShouldEmbedAdminInSPA() {
		return ""
	}
	// A single project has no workspace, so its copy of the shared package is
	// mirrored into the SPA and aliased there instead; see singleFrontendViteConfig.
	if opts.Architecture == ArchSingle {
		return ""
	}
	return `
      // The admin panel's own code, which lives in this app under src/admin-panel.
      '@admin': path.resolve(__dirname, './src/admin-panel'),
      // The upload package ships raw TypeScript, so Vite compiles it as source
      // rather than resolving a build. The panel's api-client builds its uploader
      // from it.
      '@repo/upload/web': path.resolve(__dirname, '../../packages/upload/src/web.ts'),
      '@repo/upload': path.resolve(__dirname, '../../packages/upload/src/index.ts'),`
}

// hostTSConfigPaths is the same set for the typechecker.
func hostTSConfigPaths(opts Options) string {
	if !opts.ShouldEmbedAdminInSPA() || opts.Architecture == ArchSingle {
		return ""
	}
	return `,
      "@admin/*": ["./src/admin-panel/*"],
      "@repo/upload/web": ["../../packages/upload/src/web.ts"],
      "@repo/upload": ["../../packages/upload/src/index.ts"]`
}

// viteHostDependencies is what the panel needs that the Next.js web app already
// had and this one did not.
//
// webAdminDependencies is the set the two web apps share; react-hook-form and its
// resolvers sit in the Next app's own dependency list rather than in there, so the
// first build of a --double --vite panel failed on its login screen importing
// react-hook-form.
func viteHostDependencies(opts Options) string {
	if !opts.ShouldEmbedAdminInSPA() || opts.Architecture == ArchSingle {
		return ""
	}
	return `,
    "@hookform/resolvers": "^3.3.0",
    "react-hook-form": "^7.49.0"`
}
