package scaffold

import (
	"fmt"
	"path/filepath"
)

func writeWebTanStackFiles(root string, opts Options) error {
	webRoot := filepath.Join(root, "apps", "web")

	files := map[string]string{
		filepath.Join(webRoot, "package.json"):       webTanStackPackageJSON(opts),
		filepath.Join(webRoot, "vite.config.ts"):     webTanStackViteConfig(),
		filepath.Join(webRoot, "index.html"):         webTanStackIndexHTML(opts),
		filepath.Join(webRoot, "tailwind.config.ts"): webTanStackTailwindConfig(),
		// .cjs (not .js) because package.json sets "type": "module" and PostCSS
		// config still uses CommonJS module.exports.
		filepath.Join(webRoot, "postcss.config.cjs"):                 webPostCSSConfig(),
		filepath.Join(webRoot, "tsconfig.json"):                      webTanStackTSConfig(),
		filepath.Join(webRoot, "src", "main.tsx"):                    webTanStackMain(),
		filepath.Join(webRoot, "src", "globals.css"):                 webGlobalCSS(),
		filepath.Join(webRoot, "src", "routes", "__root.tsx"):        webTanStackRootRoute(opts),
		filepath.Join(webRoot, "src", "routes", "index.tsx"):         webTanStackIndexRoute(opts),
		filepath.Join(webRoot, "src", "routes", "blog", "index.tsx"): webTanStackBlogListRoute(),
		filepath.Join(webRoot, "src", "routes", "blog", "$slug.tsx"): webTanStackBlogDetailRoute(),
		filepath.Join(webRoot, "src", "components", "navbar.tsx"):    webNavbar(opts),
		filepath.Join(webRoot, "src", "components", "footer.tsx"):    webFooter(opts),
		filepath.Join(webRoot, "src", "components", "providers.tsx"): webTanStackProviders(),
		filepath.Join(webRoot, "src", "lib", "utils.ts"):             webUtils(),
		filepath.Join(webRoot, "src", "lib", "api.ts"):               viteAPIClient(),
		filepath.Join(webRoot, "src", "hooks", "use-blogs.ts"):       webUseBlogsHook(),
		filepath.Join(webRoot, "public", ".gitkeep"):                 "",
	}

	for path, content := range files {
		if err := writeFile(path, content); err != nil {
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
    "build": "tsc -b && vite build",
    "preview": "vite preview",
    "lint": "eslint ."
  },
  "dependencies": {
    "@tanstack/react-query": "^5.62.0",
    "@tanstack/react-router": "^1.93.0",
    "axios": "^1.7.9",
    "clsx": "^2.1.1",
    "lucide-react": "^0.468.0",
    "react": "^19.0.0",
    "react-dom": "^19.0.0",
    "tailwind-merge": "^2.6.0"
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
    "@vitejs/plugin-react": "^4.3.4"
  }
}`, opts.ProjectName)
}

func webTanStackViteConfig() string {
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

func webTanStackTSConfig() string {
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
      "@/*": ["./src/*"]
    },
    "baseUrl": "."
  },
  "include": ["src"]
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

func webTanStackRootRoute(opts Options) string {
	return fmt.Sprintf(`import { createRootRoute, Outlet } from '@tanstack/react-router'
import { Navbar } from '@/components/navbar'
import { Footer } from '@/components/footer'

export const Route = createRootRoute({
  component: () => (
    <div className="min-h-screen bg-background flex flex-col">
      <Navbar />
      <main className="flex-1">
        <Outlet />
      </main>
      <Footer />
    </div>
  ),
})
`)
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
          {blogs.map((blog: any) => (
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
      <div className="prose prose-invert max-w-none" dangerouslySetInnerHTML={{ __html: blog.content }} />
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
