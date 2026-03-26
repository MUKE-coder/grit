import Link from 'next/link'
import { ArrowLeft, ArrowRight } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'
import { CodeBlock, StepWithCode } from '@/components/code-block'
import { getDocMetadata } from '@/config/docs-metadata'

export const metadata = getDocMetadata('/docs/frontend/tanstack-router')

export default function TanStackRouterPage() {
  return (
    <div className="min-h-screen bg-background">
      <SiteHeader />
      <DocsSidebar />

      <main className="lg:pl-64">
        <div className="max-w-4xl mx-auto py-12 px-6 lg:px-8">
          <div className="mb-14">
            <p className="text-sm font-mono font-medium text-primary mb-3 tracking-wide uppercase">
              Frontend
            </p>
            <h1 className="text-3xl lg:text-4xl font-bold tracking-tight text-foreground mb-4">
              TanStack Router (Vite)
            </h1>
            <p className="text-lg text-muted-foreground leading-relaxed max-w-2xl">
              When you choose TanStack Router as your frontend framework, Grit scaffolds a
              Vite-powered React SPA with file-based routing, React Query, and Tailwind CSS.
              Fast builds, small bundles, no Node.js server needed.
            </p>
          </div>

          {/* Why TanStack */}
          <div className="mb-14">
            <h2 className="text-xl font-semibold text-foreground mb-6">Why TanStack Router?</h2>
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              {[
                { title: 'Instant HMR', desc: 'Vite provides sub-50ms hot module replacement. Changes appear instantly in the browser.' },
                { title: 'Small bundles', desc: 'No server runtime overhead. The production output is static HTML + JS that any CDN can serve.' },
                { title: 'Type-safe routing', desc: 'TanStack Router provides fully type-safe route params, search params, and loaders.' },
                { title: 'File-based routes', desc: 'Routes auto-discovered by @tanstack/router-vite-plugin. No manual route registry needed.' },
              ].map((item) => (
                <div key={item.title} className="rounded-lg border border-border/40 bg-accent/20 p-5">
                  <h3 className="text-sm font-semibold text-foreground mb-1.5">{item.title}</h3>
                  <p className="text-sm text-muted-foreground leading-relaxed">{item.desc}</p>
                </div>
              ))}
            </div>
          </div>

          {/* Project structure */}
          <div className="mb-14">
            <h2 className="text-xl font-semibold text-foreground mb-4">Project structure</h2>
            <p className="text-muted-foreground mb-6">
              TanStack Router apps use <code className="text-primary bg-accent/30 px-1.5 py-0.5 rounded text-[13px]">src/routes/</code> for
              file-based routing instead of Next.js{"'"}s <code className="text-primary bg-accent/30 px-1.5 py-0.5 rounded text-[13px]">app/</code> directory.
            </p>
            <CodeBlock language="bash" filename="apps/web/ (TanStack Router)" code={`apps/web/
├── src/
│   ├── routes/
│   │   ├── __root.tsx       # Root layout (Navbar + Footer)
│   │   ├── index.tsx        # Home page (/)
│   │   └── blog/
│   │       ├── index.tsx    # Blog list (/blog)
│   │       └── $slug.tsx    # Blog detail (/blog/:slug)
│   ├── components/
│   │   ├── navbar.tsx
│   │   └── footer.tsx
│   ├── hooks/
│   │   └── use-blogs.ts
│   ├── lib/
│   │   ├── api.ts           # Axios client
│   │   └── utils.ts
│   ├── main.tsx             # Entry point
│   └── globals.css
├── index.html
├── vite.config.ts           # TanStack Router plugin + API proxy
├── tailwind.config.ts
└── package.json`} />
          </div>

          {/* Key differences */}
          <div className="mb-14">
            <h2 className="text-xl font-semibold text-foreground mb-4">Key differences from Next.js</h2>
            <div className="overflow-x-auto rounded-lg border border-border/40">
              <table className="w-full text-sm">
                <thead>
                  <tr className="border-b border-border/40 bg-accent/20">
                    <th className="px-4 py-3 text-left font-medium text-foreground/80">Aspect</th>
                    <th className="px-4 py-3 text-left font-medium text-foreground/80">Next.js</th>
                    <th className="px-4 py-3 text-left font-medium text-foreground/80">TanStack Router</th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-border/30">
                  {[
                    ['Routing', 'app/ directory convention', 'src/routes/ via Vite plugin'],
                    ['Layouts', 'layout.tsx', '__root.tsx + _layout.tsx'],
                    ['Build tool', 'Next.js (webpack/turbopack)', 'Vite'],
                    ['SSR', 'Built-in', 'SPA only (no SSR)'],
                    ['"use client"', 'Required for client components', 'Not needed (everything is client)'],
                    ['Dev server', 'next dev (:3000)', 'vite dev (:3000)'],
                    ['Output', '.next/', 'dist/'],
                    ['Params', 'useParams() from next/navigation', 'Route.useParams()'],
                    ['Navigation', '<Link> from next/link', '<Link> from @tanstack/react-router'],
                  ].map(([aspect, next, tanstack]) => (
                    <tr key={aspect} className="hover:bg-accent/20 transition-colors">
                      <td className="px-4 py-3 font-medium text-foreground">{aspect}</td>
                      <td className="px-4 py-3 text-muted-foreground">{next}</td>
                      <td className="px-4 py-3 text-primary/80">{tanstack}</td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          </div>

          {/* Route examples */}
          <div className="mb-14">
            <h2 className="text-xl font-semibold text-foreground mb-6">Route examples</h2>

            <CodeBlock language="tsx" filename="src/routes/__root.tsx" code={`import { createRootRoute, Outlet } from '@tanstack/react-router'
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
})`} highlightLines={[1, 5, 6]} />

            <CodeBlock language="tsx" filename="src/routes/blog/$slug.tsx" code={`import { createFileRoute, Link } from '@tanstack/react-router'
import { useQuery } from '@tanstack/react-query'
import { api } from '@/lib/api'

export const Route = createFileRoute('/blog/$slug')({
  component: BlogDetailPage,
})

function BlogDetailPage() {
  const { slug } = Route.useParams()
  const { data: blog } = useQuery({
    queryKey: ['blog', slug],
    queryFn: () => api.get('/api/blogs/' + slug).then(r => r.data.data),
  })

  return <h1>{blog?.title}</h1>
}`} highlightLines={[5, 10]} />
          </div>

          {/* Admin panel */}
          <div className="mb-14">
            <h2 className="text-xl font-semibold text-foreground mb-4">Admin panel with TanStack Router</h2>
            <p className="text-muted-foreground mb-6">
              When you choose TanStack Router, the admin panel also uses it. Auth and dashboard
              are handled via layout routes with <code className="text-primary bg-accent/30 px-1.5 py-0.5 rounded text-[13px]">beforeLoad</code> guards.
            </p>
            <CodeBlock language="tsx" filename="src/routes/_dashboard.tsx" code={`import { createFileRoute, Outlet, redirect } from '@tanstack/react-router'
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
})`} highlightLines={[5, 6, 7, 8, 9]} />
          </div>

          {/* Nav */}
          <div className="flex items-center justify-between pt-8 border-t border-border/40">
            <Button variant="ghost" size="sm" asChild className="text-muted-foreground/70 hover:text-foreground">
              <Link href="/docs/frontend/web-app" className="gap-1.5">
                <ArrowLeft className="h-3.5 w-3.5" />
                Web App
              </Link>
            </Button>
            <Button variant="ghost" size="sm" asChild className="text-muted-foreground/70 hover:text-foreground">
              <Link href="/docs/frontend/hooks" className="gap-1.5">
                React Query Hooks
                <ArrowRight className="h-3.5 w-3.5" />
              </Link>
            </Button>
          </div>
        </div>
      </main>
    </div>
  )
}
