import Link from 'next/link'
import { CheckCircle2 } from 'lucide-react'
import { SiteHeader } from '@/components/site-header'
import { GridFrame } from '@/components/grid-frame'
import { CodeBlock, Challenge, Note, Tip, Definition, Code, CourseNav, CourseFooter } from '@/components/course-components'
import type { Metadata } from 'next'

export const metadata: Metadata = {
  title: 'React + Vite + Go: Building with TanStack Router — Grit Course',
  description: 'Learn how to use TanStack Router with Vite as the frontend for your Grit application instead of Next.js. File-based routing, data fetching, auth, and production builds.',
}

export default function ReactViteGoCourse() {
  return (
    <div className="relative min-h-screen bg-background isolate">
      <SiteHeader />
      <GridFrame />

      <main className="max-w-3xl mx-auto px-6 py-12">
        {/* Breadcrumb */}
        <div className="flex items-center gap-2 text-sm text-muted-foreground mb-8">
          <Link href="/courses" className="hover:text-foreground transition-colors">Courses</Link>
          <span>/</span>
          <span className="text-foreground">React + Vite + Go</span>
        </div>

        {/* Header */}
        <div className="mb-10">
          <div className="flex items-center gap-2 mb-3">
            <span className="text-xs font-mono font-medium text-primary bg-primary/10 px-2 py-0.5 rounded">Standalone Course</span>
            <span className="text-xs text-muted-foreground">~30 min</span>
            <span className="text-xs text-muted-foreground">•</span>
            <span className="text-xs text-muted-foreground">12 challenges</span>
          </div>
          <h1 className="text-3xl md:text-4xl font-bold text-foreground mb-4">
            React + Vite + Go: Building with TanStack Router
          </h1>
          <p className="text-lg text-muted-foreground leading-relaxed">
            Not every app needs server-side rendering. In this course, you will learn how to use TanStack
            Router with Vite as your frontend instead of Next.js — perfect for dashboards, internal tools,
            and admin panels where SEO doesn{"'"}t matter and speed is everything.
          </p>
        </div>

        <hr className="border-border/40 mb-10" />

        {/* ═══ Section 1: Why TanStack Router? ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Why TanStack Router?</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Throughout the Grit Web course, you used Next.js as your frontend. Next.js is great for
            public-facing websites that need SEO, but it comes with overhead you don{"'"}t always need —
            server-side rendering, server components, hydration. For apps where every user is logged in,
            there{"'"}s a faster alternative.
          </p>

          <Definition term="SPA (Single Page Application)">
            A web application that loads a single HTML page and dynamically updates content as the user
            interacts with it. The browser downloads the entire app upfront, then all navigation happens
            client-side without full page reloads. Gmail, Figma, and Notion are SPAs.
          </Definition>

          <Definition term="Vite">
            A modern build tool for JavaScript/TypeScript projects. It serves code via native ES modules
            during development (instant startup, near-instant hot reload) and uses Rollup for optimized
            production builds. Much faster than webpack-based tools.
          </Definition>

          <Definition term="TanStack Router">
            A fully type-safe router for React applications built by Tanner Linsley (the creator of
            TanStack Query). It supports file-based routing (like Next.js), URL search params as state,
            and deep integration with TypeScript — every route, parameter, and link is type-checked at
            compile time.
          </Definition>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Here{"'"}s how TanStack Router + Vite compares to Next.js:
          </p>

          <div className="overflow-x-auto mb-4">
            <table className="w-full text-sm border border-border/40 rounded-lg overflow-hidden">
              <thead>
                <tr className="bg-muted/20">
                  <th className="text-left px-4 py-3 font-semibold text-foreground border-b border-border/40">Feature</th>
                  <th className="text-left px-4 py-3 font-semibold text-foreground border-b border-border/40">TanStack Router + Vite</th>
                  <th className="text-left px-4 py-3 font-semibold text-foreground border-b border-border/40">Next.js</th>
                </tr>
              </thead>
              <tbody className="text-muted-foreground">
                <tr className="border-b border-border/20">
                  <td className="px-4 py-3 font-medium text-foreground">Rendering</td>
                  <td className="px-4 py-3">Client-side only (SPA)</td>
                  <td className="px-4 py-3">Server + Client (SSR/SSG)</td>
                </tr>
                <tr className="border-b border-border/20">
                  <td className="px-4 py-3 font-medium text-foreground">Dev startup</td>
                  <td className="px-4 py-3">~200ms (Vite HMR)</td>
                  <td className="px-4 py-3">~2-5s (compilation)</td>
                </tr>
                <tr className="border-b border-border/20">
                  <td className="px-4 py-3 font-medium text-foreground">Bundle size</td>
                  <td className="px-4 py-3">Smaller (no SSR runtime)</td>
                  <td className="px-4 py-3">Larger (includes SSR)</td>
                </tr>
                <tr className="border-b border-border/20">
                  <td className="px-4 py-3 font-medium text-foreground">SEO</td>
                  <td className="px-4 py-3">Poor (client-rendered)</td>
                  <td className="px-4 py-3">Excellent (pre-rendered)</td>
                </tr>
                <tr className="border-b border-border/20">
                  <td className="px-4 py-3 font-medium text-foreground">Type safety</td>
                  <td className="px-4 py-3">Full (routes, params, links)</td>
                  <td className="px-4 py-3">Partial</td>
                </tr>
                <tr>
                  <td className="px-4 py-3 font-medium text-foreground">Best for</td>
                  <td className="px-4 py-3">Dashboards, admin panels, internal tools</td>
                  <td className="px-4 py-3">Public websites, blogs, e-commerce</td>
                </tr>
              </tbody>
            </table>
          </div>

          <Tip>
            If your app requires login to use (dashboard, admin panel, CRM), TanStack Router + Vite is
            likely the better choice. If your app has public pages that need to rank on Google (blog,
            landing page, product listings), stick with Next.js.
          </Tip>

          <Challenge number={1} title="When Would You Choose TanStack Router?">
            <p>For each of these projects, decide whether you{"'"}d use TanStack Router + Vite or Next.js, and explain why: (1) a company blog, (2) an employee timesheet tracker, (3) a real estate listing site, (4) a project management dashboard.</p>
          </Challenge>
        </section>

        {/* ═══ Section 2: Scaffold with Vite ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Scaffold with Vite</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Grit supports Vite as an alternative frontend. When you pass the <Code>--vite</Code> flag,
            Grit scaffolds TanStack Router instead of Next.js:
          </p>

          <CodeBlock filename="Terminal">
{`grit new dashboard --double --vite`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            This creates a Double architecture (Go API + frontend) but with a Vite-powered React app
            instead of Next.js. Here{"'"}s what{"'"}s different:
          </p>

          <div className="overflow-x-auto mb-4">
            <table className="w-full text-sm border border-border/40 rounded-lg overflow-hidden">
              <thead>
                <tr className="bg-muted/20">
                  <th className="text-left px-4 py-3 font-semibold text-foreground border-b border-border/40">Next.js (--next)</th>
                  <th className="text-left px-4 py-3 font-semibold text-foreground border-b border-border/40">TanStack Router (--vite)</th>
                </tr>
              </thead>
              <tbody className="text-muted-foreground">
                <tr className="border-b border-border/20">
                  <td className="px-4 py-3"><Code>app/</Code> directory</td>
                  <td className="px-4 py-3"><Code>src/routes/</Code> directory</td>
                </tr>
                <tr className="border-b border-border/20">
                  <td className="px-4 py-3"><Code>next.config.js</Code></td>
                  <td className="px-4 py-3"><Code>vite.config.ts</Code></td>
                </tr>
                <tr className="border-b border-border/20">
                  <td className="px-4 py-3"><Code>layout.tsx</Code> for layouts</td>
                  <td className="px-4 py-3"><Code>__root.tsx</Code> for root layout</td>
                </tr>
                <tr>
                  <td className="px-4 py-3"><Code>page.tsx</Code> for pages</td>
                  <td className="px-4 py-3"><Code>index.tsx</Code> or <Code>about.tsx</Code> for pages</td>
                </tr>
              </tbody>
            </table>
          </div>

          <Challenge number={2} title="Scaffold and Compare">
            <p>Scaffold two projects — one with <Code>--next</Code> and one with <Code>--vite</Code>. Compare the frontend folder structures. How many files are different? Which has more configuration files?</p>
          </Challenge>
        </section>

        {/* ═══ Section 3: File-Based Routing ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">File-Based Routing</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            TanStack Router discovers routes from your file system, just like Next.js. But the naming
            conventions are different:
          </p>

          <Definition term="File-Based Routing">
            The file and folder structure inside your routes directory determines your URL routes. No
            manual route configuration needed. TanStack Router reads the file tree, auto-generates a
            route tree, and provides full type safety for every route. Drop a file in the right folder
            and it becomes a URL.
          </Definition>

          <CodeBlock filename="src/routes/">
{`src/routes/
├── __root.tsx       ← Root layout (navbar, providers)
├── index.tsx        ← / (home page)
├── about.tsx        ← /about
├── products/
│   ├── index.tsx    ← /products (list)
│   ├── $id.tsx      ← /products/:id (detail)
│   └── new.tsx      ← /products/new (create)`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Key naming rules:
          </p>

          <ul className="space-y-2 text-muted-foreground mb-4">
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">__root.tsx</strong> — the root layout that wraps every page (like Next.js{"'"} root <Code>layout.tsx</Code>)</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">index.tsx</strong> — the default page for a directory (like Next.js{"'"} <Code>page.tsx</Code>)</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">$param.tsx</strong> — a dynamic route segment (Next.js uses <Code>[param]</Code>)</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">_layout.tsx</strong> — a nested layout for a route group</li>
          </ul>

          <Note>
            TanStack Router auto-generates a <Code>routeTree.gen.ts</Code> file every time you add or
            remove a route file. This file contains the full type-safe route tree. Never edit it manually —
            it{"'"}s regenerated automatically.
          </Note>

          <Challenge number={3} title="Create a New Route">
            <p>Create a file at <Code>src/routes/settings.tsx</Code> with a simple component. Run the dev server — does TanStack Router pick it up automatically? Can you navigate to <Code>/settings</Code>?</p>
          </Challenge>
        </section>

        {/* ═══ Section 4: Route Parameters ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Route Parameters</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Files starting with <Code>$</Code> create dynamic route parameters. The <Code>$id.tsx</Code> file
            matches any value in that URL segment and makes it available as a typed parameter:
          </p>

          <CodeBlock filename="src/routes/products/$id.tsx">
{`import { createFileRoute } from '@tanstack/react-router'

export const Route = createFileRoute('/products/$id')({
  component: ProductDetail,
})

function ProductDetail() {
  const { id } = Route.useParams()

  return (
    <div>
      <h1>Product {id}</h1>
      {/* Fetch and display product details */}
    </div>
  )
}`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            The <Code>id</Code> parameter is fully type-safe. If you try to access a parameter that
            doesn{"'"}t exist, TypeScript will catch the error at compile time. This is a major advantage
            over Next.js{"'"} <Code>params</Code> prop, which is loosely typed.
          </p>

          <Tip>
            You can have multiple parameters in a single route. A file at <Code>src/routes/teams/$teamId/members/$memberId.tsx</Code> gives you both <Code>teamId</Code> and <Code>memberId</Code> as typed params.
          </Tip>

          <Challenge number={4} title="Create a Dynamic Route">
            <p>Create a route at <Code>src/routes/users/$id.tsx</Code> that displays the user ID from the URL. Navigate to <Code>/users/42</Code> and <Code>/users/hello</Code> — does the parameter work for both?</p>
          </Challenge>
        </section>

        {/* ═══ Section 5: Layouts ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Layouts</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            The <Code>__root.tsx</Code> file is the root layout — it wraps every page in your app.
            This is where you place your navbar, providers (query client, auth context), and global styles:
          </p>

          <CodeBlock filename="src/routes/__root.tsx">
{`import { createRootRoute, Outlet } from '@tanstack/react-router'
import { Navbar } from '../components/navbar'

export const Route = createRootRoute({
  component: RootLayout,
})

function RootLayout() {
  return (
    <div className="min-h-screen bg-background isolate">
      <Navbar />
      <Outlet />
    </div>
  )
}`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            The <Code>Outlet</Code> component renders the current route{"'"}s page — just like the
            <Code> children</Code> prop in Next.js layouts.
          </p>

          <p className="text-muted-foreground leading-relaxed mb-4">
            For nested layouts (e.g., a sidebar that only appears on dashboard pages), create a route
            group with a layout file:
          </p>

          <CodeBlock filename="src/routes/dashboard/">
{`src/routes/
├── dashboard/
│   ├── _layout.tsx      ← Sidebar layout for all /dashboard/* routes
│   ├── index.tsx        ← /dashboard (main dashboard)
│   ├── projects.tsx     ← /dashboard/projects
│   └── settings.tsx     ← /dashboard/settings`}
          </CodeBlock>

          <Challenge number={5} title="Add a Sidebar Layout">
            <p>Create a <Code>dashboard/_layout.tsx</Code> that renders a sidebar on the left and the page content on the right. Add at least 3 routes under <Code>/dashboard/</Code>. Does the sidebar persist across navigation?</p>
          </Challenge>
        </section>

        {/* ═══ Section 6: Data Fetching ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Data Fetching</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Data fetching in a Vite app is identical to what you learned in the Grit Web course.
            TanStack Query (<Code>useQuery</Code>, <Code>useMutation</Code>) works the same way — it{"'"}s
            a React library, not a Next.js feature. The only difference is that everything runs client-side.
          </p>

          <CodeBlock filename="src/routes/products/index.tsx">
{`import { useQuery } from '@tanstack/react-query'
import { api } from '../../lib/api-client'

function ProductList() {
  const { data, isLoading, error } = useQuery({
    queryKey: ['products'],
    queryFn: () => api.get('/api/products'),
  })

  if (isLoading) return <div>Loading...</div>
  if (error) return <div>Error loading products</div>

  return (
    <div>
      <h1>Products</h1>
      {data.data.map((product: any) => (
        <div key={product.id}>{product.name}</div>
      ))}
    </div>
  )
}`}
          </CodeBlock>

          <Note>
            Since there{"'"}s no server-side rendering, the user sees a loading state while data is being
            fetched. This is normal for SPAs. For dashboards and internal tools, this is perfectly acceptable.
            For public pages where you want content visible immediately, you{"'"}d use Next.js with SSR.
          </Note>

          <Challenge number={6} title="Fetch and Display Data">
            <p>Generate a <Code>Product</Code> resource with <Code>grit generate</Code>. Then build a products list page that fetches from <Code>/api/products</Code> and displays them in a table. Does the loading state appear before the data?</p>
          </Challenge>
        </section>

        {/* ═══ Section 7: Navigation ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Navigation</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            TanStack Router provides its own <Code>Link</Code> component and <Code>useNavigate</Code> hook.
            The <Code>Link</Code> component is type-safe — it only allows you to link to routes that actually exist:
          </p>

          <CodeBlock filename="Navigation Examples">
{`import { Link, useNavigate } from '@tanstack/react-router'

// Declarative navigation (in JSX)
<Link to="/products">All Products</Link>

// With parameters — fully type-checked
<Link to="/products/$id" params={{ id: '42' }}>
  View Product
</Link>

// Programmatic navigation (in event handlers)
function CreateButton() {
  const navigate = useNavigate()

  const handleCreate = async () => {
    const product = await createProduct(data)
    navigate({ to: '/products/$id', params: { id: product.id } })
  }

  return <button onClick={handleCreate}>Create</button>
}`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            If you try to link to a route that doesn{"'"}t exist, TypeScript shows an error immediately.
            If you forget a required parameter, TypeScript catches that too. This is a significant
            improvement over Next.js{"'"} <Code>Link</Code> component, which accepts any string.
          </p>

          <Challenge number={7} title="Add Navigation Links">
            <p>Add a navbar with <Code>Link</Code> components to your root layout. Include links to Home, Products, and Dashboard. Then add a button that uses <Code>useNavigate</Code> to redirect to a product detail page after creating a product.</p>
          </Challenge>
        </section>

        {/* ═══ Section 8: Auth in Vite Apps ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Auth in Vite Apps</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Authentication in a Vite SPA works the same way as in Next.js, with one key difference:
            there{"'"}s no server to manage cookies. Instead, you store the JWT token in <Code>localStorage</Code> and
            include it in every API request via an Authorization header.
          </p>

          <CodeBlock filename="src/lib/auth-context.tsx">
{`// Auth context wraps the entire app
function AuthProvider({ children }) {
  const [token, setToken] = useState(localStorage.getItem('token'))
  const [user, setUser] = useState(null)

  const login = async (email: string, password: string) => {
    const res = await api.post('/api/auth/login', { email, password })
    localStorage.setItem('token', res.data.token)
    setToken(res.data.token)
    setUser(res.data.user)
  }

  const logout = () => {
    localStorage.removeItem('token')
    setToken(null)
    setUser(null)
  }

  return (
    <AuthContext.Provider value={{ user, token, login, logout }}>
      {children}
    </AuthContext.Provider>
  )
}`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            To protect routes, check for the token and redirect to <Code>/login</Code> if not authenticated.
            You can do this in a layout component that wraps all protected routes:
          </p>

          <CodeBlock filename="src/routes/dashboard/_layout.tsx">
{`function DashboardLayout() {
  const { user } = useAuth()
  const navigate = useNavigate()

  if (!user) {
    navigate({ to: '/login' })
    return null
  }

  return (
    <div className="flex">
      <Sidebar />
      <main className="flex-1 p-6">
        <Outlet />
      </main>
    </div>
  )
}`}
          </CodeBlock>

          <Challenge number={8} title="Implement a Protected Route">
            <p>Set up an auth context with login/logout. Create a <Code>/login</Code> page and a protected <Code>/dashboard</Code> layout. Verify that visiting <Code>/dashboard</Code> without logging in redirects to <Code>/login</Code>.</p>
          </Challenge>
        </section>

        {/* ═══ Section 9: Building for Production ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Building for Production</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            When you{"'"}re ready to deploy, build the Vite frontend into static files:
          </p>

          <CodeBlock filename="Terminal">
{`cd frontend && pnpm build`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            This outputs optimized files to <Code>frontend/dist/</Code>. Vite tree-shakes unused code,
            minifies everything, and splits the bundle into chunks for optimal loading.
          </p>

          <p className="text-muted-foreground leading-relaxed mb-4">
            For Grit{"'"}s Single architecture, the <Code>dist/</Code> folder gets embedded directly
            into the Go binary via <Code>go:embed</Code>. This means your entire app — API and frontend — is
            a single executable file:
          </p>

          <CodeBlock filename="cmd/server/main.go">
{`//go:embed all:frontend/dist
var frontendFS embed.FS

// Serve frontend from the embedded files
router.NoRoute(gin.WrapH(http.FileServer(http.FS(frontendFS))))`}
          </CodeBlock>

          <Tip>
            A typical Vite + TanStack Router build produces a bundle between 150-300 KB gzipped,
            compared to 500 KB+ for a Next.js app. The smaller bundle means faster initial load for
            your users.
          </Tip>

          <Challenge number={9} title="Build and Measure">
            <p>Run <Code>pnpm build</Code> in the frontend directory. Check the output — what{"'"}s the total bundle size? How many chunks were created? Compare this to a Next.js build if you have one available.</p>
          </Challenge>
        </section>

        {/* ═══ Section 10: Vite Dev Server Proxy ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Vite Dev Server Proxy</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            During development, the Vite dev server runs on one port (typically 5173) and the Go API
            runs on another (typically 8080). The Vite proxy configuration forwards API requests to
            the Go server so you don{"'"}t have to deal with CORS:
          </p>

          <CodeBlock filename="vite.config.ts">
{`import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import { TanStackRouterVite } from '@tanstack/router-plugin/vite'

export default defineConfig({
  plugins: [TanStackRouterVite(), react()],
  server: {
    proxy: {
      '/api': 'http://localhost:8080',
    },
  },
})`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            With this configuration, when your frontend code calls <Code>fetch({'"'}/api/products{'"'})</Code>,
            Vite intercepts the request and forwards it to <Code>http://localhost:8080/api/products</Code>.
            This only applies during development — in production, both frontend and API are served from
            the same origin.
          </p>

          <Challenge number={10} title="Explore the Proxy Config">
            <p>Open <Code>vite.config.ts</Code> in your scaffolded project. Find the proxy configuration. What path prefix does it intercept? What happens if you change the Go server port — what else would you need to update?</p>
          </Challenge>
        </section>

        {/* ═══ Section 11: When to Use Next.js vs Vite ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">When to Use Next.js vs Vite</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Here{"'"}s the decision framework:
          </p>

          <div className="overflow-x-auto mb-4">
            <table className="w-full text-sm border border-border/40 rounded-lg overflow-hidden">
              <thead>
                <tr className="bg-muted/20">
                  <th className="text-left px-4 py-3 font-semibold text-foreground border-b border-border/40">Criteria</th>
                  <th className="text-left px-4 py-3 font-semibold text-foreground border-b border-border/40">Use Next.js</th>
                  <th className="text-left px-4 py-3 font-semibold text-foreground border-b border-border/40">Use Vite + TanStack Router</th>
                </tr>
              </thead>
              <tbody className="text-muted-foreground">
                <tr className="border-b border-border/20">
                  <td className="px-4 py-3 font-medium text-foreground">SEO needed?</td>
                  <td className="px-4 py-3">Yes — blog, e-commerce, landing pages</td>
                  <td className="px-4 py-3">No — dashboards, admin panels</td>
                </tr>
                <tr className="border-b border-border/20">
                  <td className="px-4 py-3 font-medium text-foreground">Public pages?</td>
                  <td className="px-4 py-3">Many public pages</td>
                  <td className="px-4 py-3">Mostly behind login</td>
                </tr>
                <tr className="border-b border-border/20">
                  <td className="px-4 py-3 font-medium text-foreground">Build speed</td>
                  <td className="px-4 py-3">Slower (webpack/turbopack)</td>
                  <td className="px-4 py-3">Faster (Vite/esbuild)</td>
                </tr>
                <tr className="border-b border-border/20">
                  <td className="px-4 py-3 font-medium text-foreground">Bundle size</td>
                  <td className="px-4 py-3">Larger</td>
                  <td className="px-4 py-3">Smaller</td>
                </tr>
                <tr className="border-b border-border/20">
                  <td className="px-4 py-3 font-medium text-foreground">Deploy to</td>
                  <td className="px-4 py-3">Vercel, VPS, Docker</td>
                  <td className="px-4 py-3">Embedded in Go binary, VPS, CDN</td>
                </tr>
                <tr>
                  <td className="px-4 py-3 font-medium text-foreground">Type safety</td>
                  <td className="px-4 py-3">Partial (loose route params)</td>
                  <td className="px-4 py-3">Full (routes, params, links)</td>
                </tr>
              </tbody>
            </table>
          </div>

          <Challenge number={11} title="Make the Right Choice">
            <p>For each project below, decide Next.js or Vite + TanStack Router, and explain your reasoning: (1) a personal blog with comments, (2) a company HR dashboard, (3) an online store with 10,000 products, (4) an internal inventory management tool, (5) a SaaS landing page + app.</p>
          </Challenge>
        </section>

        {/* ═══ Section 12: Summary + Final Challenge ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">What You Learned</h2>

          <ul className="space-y-2 text-muted-foreground mb-6">
            <li className="flex gap-2"><CheckCircle2 className="h-4 w-4 text-green-500 shrink-0 mt-1" /> When to choose TanStack Router + Vite over Next.js</li>
            <li className="flex gap-2"><CheckCircle2 className="h-4 w-4 text-green-500 shrink-0 mt-1" /> How to scaffold a Vite project with <Code>grit new --vite</Code></li>
            <li className="flex gap-2"><CheckCircle2 className="h-4 w-4 text-green-500 shrink-0 mt-1" /> File-based routing with TanStack Router</li>
            <li className="flex gap-2"><CheckCircle2 className="h-4 w-4 text-green-500 shrink-0 mt-1" /> Type-safe route parameters, layouts, and navigation</li>
            <li className="flex gap-2"><CheckCircle2 className="h-4 w-4 text-green-500 shrink-0 mt-1" /> Data fetching with TanStack Query (same as Next.js)</li>
            <li className="flex gap-2"><CheckCircle2 className="h-4 w-4 text-green-500 shrink-0 mt-1" /> JWT auth with localStorage for SPAs</li>
            <li className="flex gap-2"><CheckCircle2 className="h-4 w-4 text-green-500 shrink-0 mt-1" /> Vite dev server proxy for API calls</li>
            <li className="flex gap-2"><CheckCircle2 className="h-4 w-4 text-green-500 shrink-0 mt-1" /> Building for production and embedding in Go</li>
          </ul>

          <Challenge number={12} title="Build a Task Dashboard">
            <p>Build a complete task dashboard from scratch:</p>
            <ol className="list-decimal list-inside space-y-1 mt-2">
              <li>Scaffold with <Code>grit new task-dash --double --vite</Code></li>
              <li>Generate a Task resource: <Code>grit generate resource Task title:string status:string priority:int due_date:date:optional</Code></li>
              <li>Build a task list page with filtering by status</li>
              <li>Build a create/edit task form</li>
              <li>Add an auth guard on all routes</li>
              <li>Build and check the production bundle size</li>
            </ol>
          </Challenge>
        </section>

        {/* Footer */}
        <CourseFooter />

        {/* Navigation */}
        <div className="mt-8">
          <CourseNav
            prev={{ href: '/courses', label: 'All Courses' }}
            next={{ href: '/courses', label: 'More Courses' }}
          />
        </div>
      </main>
    </div>
  )
}
