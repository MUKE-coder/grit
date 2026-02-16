import Link from 'next/link'
import { ArrowRight, ArrowLeft } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'

export default function WebAppPage() {
  return (
    <div className="min-h-screen bg-background">
      <SiteHeader />
      <DocsSidebar />

      <main className="lg:pl-64">
        <div className="container max-w-screen-xl py-10 px-6">
          <div className="max-w-3xl">
            {/* Header */}
            <div className="mb-10">
              <span className="tag-mono text-primary/80 mb-3 block">Frontend</span>
              <h1 className="text-4xl font-bold tracking-tight mb-4">
                Web App
              </h1>
              <p className="text-lg text-muted-foreground leading-relaxed">
                Grit scaffolds a full Next.js 14+ frontend with App Router, Tailwind CSS, shadcn/ui,
                React Query, and a pre-built authentication flow -- all wired to your Go API out of the box.
              </p>
            </div>

            <div className="prose-grit">
              {/* Architecture Overview */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Architecture Overview
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  The web app lives at <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">apps/web/</code> inside the monorepo. It uses Next.js App Router for file-based routing, Tailwind CSS
                  with shadcn/ui for the design system, and React Query for all data fetching. The dark theme
                  is enabled by default, matching the premium aesthetic across the entire Grit stack.
                </p>

                <div className="rounded-xl border border-border/40 bg-card/80 overflow-hidden">
                  <div className="flex items-center gap-2 px-4 py-2.5 border-b border-border/30 bg-accent/30">
                    <span className="text-[11px] font-mono text-muted-foreground/40">apps/web/</span>
                  </div>
                  <pre className="p-5 text-sm font-mono text-foreground/80 overflow-x-auto">{`apps/web/
├── app/
│   ├── layout.tsx              # Root layout with providers
│   ├── page.tsx                # Landing / home page
│   ├── globals.css             # Tailwind + dark theme variables
│   ├── (auth)/                 # Auth route group (public)
│   │   ├── login/page.tsx
│   │   ├── register/page.tsx
│   │   └── forgot-password/page.tsx
│   ├── (dashboard)/            # Protected route group
│   │   ├── layout.tsx          # Sidebar + auth guard
│   │   └── dashboard/page.tsx  # Main dashboard
│   └── blog/                   # Public blog pages
│       ├── page.tsx            # Blog listing page
│       └── [slug]/page.tsx     # Blog detail page
├── components/
│   ├── ui/                     # shadcn/ui primitives
│   ├── navbar.tsx              # Site navigation bar
│   ├── footer.tsx              # Site footer
│   ├── providers.tsx           # QueryClient + theme providers
│   └── shared/
│       └── providers.tsx       # QueryClient + theme providers
├── hooks/
│   ├── use-auth.ts             # Auth hooks (login, register, etc.)
│   ├── use-users.ts            # Generated CRUD hooks
│   └── use-blogs.ts            # Blog data fetching hooks
├── lib/
│   ├── api.ts                  # API helper functions
│   ├── api-client.ts           # Axios instance with JWT interceptor
│   ├── query-client.ts         # React Query client config
│   └── utils.ts                # Utility functions (cn, etc.)
├── next.config.ts
├── tailwind.config.ts
├── tsconfig.json
└── package.json`}</pre>
                </div>
              </div>

              {/* App Router Structure */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  App Router Structure
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Grit uses Next.js route groups to separate public and protected routes.
                  The <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">(auth)</code> group
                  holds login, register, and forgot-password pages. The <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">(dashboard)</code> group
                  wraps all authenticated pages with a sidebar layout and auth guard.
                </p>

                <h3 className="text-lg font-semibold tracking-tight mb-3 mt-6">Root Layout</h3>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  The root layout wraps the entire app with the Providers component, which sets up React Query
                  and the theme. The DM Sans font is loaded via Google Fonts and JetBrains Mono is used
                  for code elements.
                </p>

                <div className="rounded-xl border border-border/40 bg-card/80 overflow-hidden mb-6">
                  <div className="flex items-center gap-2 px-4 py-2.5 border-b border-border/30 bg-accent/30">
                    <span className="text-[11px] font-mono text-muted-foreground/40">apps/web/app/layout.tsx</span>
                  </div>
                  <pre className="p-5 text-sm font-mono text-foreground/80 overflow-x-auto">{`import type { Metadata } from "next";
import { Providers } from "@/components/shared/providers";
import "./globals.css";

export const metadata: Metadata = {
  title: "My App",
  description: "Built with Grit",
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="en" className="dark">
      <body>
        <Providers>{children}</Providers>
      </body>
    </html>
  );
}`}</pre>
                </div>

                <h3 className="text-lg font-semibold tracking-tight mb-3 mt-6">Providers Component</h3>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  The Providers component initializes React Query with a shared QueryClient instance
                  and wraps children with the QueryClientProvider. Toast notifications use Sonner.
                </p>

                <div className="rounded-xl border border-border/40 bg-card/80 overflow-hidden mb-6">
                  <div className="flex items-center gap-2 px-4 py-2.5 border-b border-border/30 bg-accent/30">
                    <span className="text-[11px] font-mono text-muted-foreground/40">apps/web/components/shared/providers.tsx</span>
                  </div>
                  <pre className="p-5 text-sm font-mono text-foreground/80 overflow-x-auto">{`"use client";

import { QueryClientProvider } from "@tanstack/react-query";
import { Toaster } from "sonner";
import { queryClient } from "@/lib/query-client";

export function Providers({ children }: { children: React.ReactNode }) {
  return (
    <QueryClientProvider client={queryClient}>
      {children}
      <Toaster theme="dark" position="top-right" />
    </QueryClientProvider>
  );
}`}</pre>
                </div>
              </div>

              {/* Auth Pages */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Authentication Pages
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Grit ships three pre-built auth pages with the dark theme, form validation via Zod schemas
                  from the shared package, and API integration via the auth hooks. All pages are fully responsive.
                </p>

                <div className="rounded-lg border border-border/30 bg-card/30 overflow-hidden mb-6">
                  <table className="w-full text-sm">
                    <thead>
                      <tr className="border-b border-border/30 bg-accent/20">
                        <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Route</th>
                        <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Page</th>
                        <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Description</th>
                      </tr>
                    </thead>
                    <tbody className="text-muted-foreground">
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-mono text-xs">/login</td>
                        <td className="px-4 py-2.5">Login</td>
                        <td className="px-4 py-2.5">Email + password, JWT token storage</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-mono text-xs">/register</td>
                        <td className="px-4 py-2.5">Register</td>
                        <td className="px-4 py-2.5">Name + email + password + confirm</td>
                      </tr>
                      <tr>
                        <td className="px-4 py-2.5 font-mono text-xs">/forgot-password</td>
                        <td className="px-4 py-2.5">Forgot Password</td>
                        <td className="px-4 py-2.5">Email input, triggers reset email</td>
                      </tr>
                    </tbody>
                  </table>
                </div>

                <p className="text-muted-foreground leading-relaxed mb-4">
                  Here is the login page pattern. It uses the <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">useLogin()</code> hook
                  from <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">hooks/use-auth.ts</code> and
                  validates with the <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">LoginSchema</code> from the shared package.
                </p>

                <div className="rounded-xl border border-border/40 bg-card/80 overflow-hidden">
                  <div className="flex items-center gap-2 px-4 py-2.5 border-b border-border/30 bg-accent/30">
                    <span className="text-[11px] font-mono text-muted-foreground/40">apps/web/app/(auth)/login/page.tsx</span>
                  </div>
                  <pre className="p-5 text-sm font-mono text-foreground/80 overflow-x-auto">{`"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import Link from "next/link";
import { LoginSchema, type LoginInput } from "@myapp/shared/schemas";
import { useLogin } from "@/hooks/use-auth";

export default function LoginPage() {
  const router = useRouter();
  const login = useLogin();
  const [form, setForm] = useState<LoginInput>({
    email: "", password: "",
  });
  const [errors, setErrors] = useState<Record<string, string>>({});

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setErrors({});

    const result = LoginSchema.safeParse(form);
    if (!result.success) {
      const fieldErrors: Record<string, string> = {};
      result.error.issues.forEach((issue) => {
        fieldErrors[issue.path[0] as string] = issue.message;
      });
      setErrors(fieldErrors);
      return;
    }

    login.mutate(form, {
      onSuccess: () => router.push("/dashboard"),
      onError: (err) => setErrors({ form: "Invalid credentials" }),
    });
  };

  return (
    <div className="min-h-screen flex items-center justify-center">
      <form onSubmit={handleSubmit} className="w-full max-w-sm space-y-4">
        <h1 className="text-2xl font-bold">Sign In</h1>
        {/* Email and password inputs with error display */}
        {/* Submit button with loading state */}
        <p className="text-sm text-text-muted">
          Don&apos;t have an account?{" "}
          <Link href="/register" className="text-accent">Register</Link>
        </p>
      </form>
    </div>
  );
}`}</pre>
                </div>
              </div>

              {/* Protected Dashboard */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Protected Dashboard Layout
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  The <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">(dashboard)</code> route group wraps all
                  authenticated pages. Its layout component checks for a valid JWT token and redirects
                  unauthenticated users to the login page. It renders a sidebar for navigation and a top
                  navbar with the user avatar and logout button.
                </p>

                <div className="rounded-xl border border-border/40 bg-card/80 overflow-hidden mb-6">
                  <div className="flex items-center gap-2 px-4 py-2.5 border-b border-border/30 bg-accent/30">
                    <span className="text-[11px] font-mono text-muted-foreground/40">apps/web/app/(dashboard)/layout.tsx</span>
                  </div>
                  <pre className="p-5 text-sm font-mono text-foreground/80 overflow-x-auto">{`"use client";

import { useEffect } from "react";
import { useRouter } from "next/navigation";
import { useMe } from "@/hooks/use-auth";
import Cookies from "js-cookie";

export default function DashboardLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  const router = useRouter();
  const { data: user, isLoading, isError } = useMe();

  useEffect(() => {
    const token = Cookies.get("access_token");
    if (!token) {
      router.replace("/login");
    }
  }, [router]);

  useEffect(() => {
    if (isError) {
      Cookies.remove("access_token");
      Cookies.remove("refresh_token");
      router.replace("/login");
    }
  }, [isError, router]);

  if (isLoading) {
    return <div className="min-h-screen bg-background" />;
  }

  return (
    <div className="min-h-screen bg-background flex">
      {/* Sidebar navigation */}
      <aside className="hidden lg:flex w-64 flex-col border-r border-border">
        {/* Logo, nav links, user section */}
      </aside>

      {/* Main content */}
      <main className="flex-1">
        {/* Top navbar with user avatar */}
        <div className="p-6">{children}</div>
      </main>
    </div>
  );
}`}</pre>
                </div>

                <p className="text-muted-foreground leading-relaxed mb-4">
                  The dashboard page itself shows a welcome message, stats cards (total users, active sessions, etc.),
                  and a recent activity feed. Stats are fetched from the Go API using React Query.
                </p>

                <div className="rounded-xl border border-border/40 bg-card/80 overflow-hidden">
                  <div className="flex items-center gap-2 px-4 py-2.5 border-b border-border/30 bg-accent/30">
                    <span className="text-[11px] font-mono text-muted-foreground/40">apps/web/app/(dashboard)/dashboard/page.tsx</span>
                  </div>
                  <pre className="p-5 text-sm font-mono text-foreground/80 overflow-x-auto">{`"use client";

import { useMe } from "@/hooks/use-auth";

export default function DashboardPage() {
  const { data: user } = useMe();

  return (
    <div>
      <h1 className="text-2xl font-bold mb-6">
        Welcome back, {user?.name}
      </h1>

      {/* Stats cards grid */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4 mb-8">
        <StatsCard label="Total Users" value="1,234" />
        <StatsCard label="Active Sessions" value="56" />
        <StatsCard label="Revenue" value="$12,345" />
        <StatsCard label="Growth" value="+12.3%" />
      </div>

      {/* Recent activity */}
      <div className="rounded-xl border border-border bg-bg-secondary p-6">
        <h2 className="text-lg font-semibold mb-4">Recent Activity</h2>
        {/* Activity list */}
      </div>
    </div>
  );
}`}</pre>
                </div>
              </div>

              {/* API Client */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  API Client Setup
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  The API client is an Axios instance configured to point at the Go backend. It automatically
                  injects the JWT access token from cookies into every request, and handles token refresh
                  when a 401 response is received.
                </p>

                <div className="rounded-xl border border-border/40 bg-card/80 overflow-hidden mb-6">
                  <div className="flex items-center gap-2 px-4 py-2.5 border-b border-border/30 bg-accent/30">
                    <span className="text-[11px] font-mono text-muted-foreground/40">apps/web/lib/api-client.ts</span>
                  </div>
                  <pre className="p-5 text-sm font-mono text-foreground/80 overflow-x-auto">{`import axios from "axios";
import Cookies from "js-cookie";

const API_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";

export const apiClient = axios.create({
  baseURL: API_URL,
  headers: { "Content-Type": "application/json" },
});

// Inject JWT token into every request
apiClient.interceptors.request.use((config) => {
  const token = Cookies.get("access_token");
  if (token) {
    config.headers.Authorization = \`Bearer \${token}\`;
  }
  return config;
});

// Handle 401 by refreshing the token
apiClient.interceptors.response.use(
  (response) => response,
  async (error) => {
    const originalRequest = error.config;

    if (error.response?.status === 401 && !originalRequest._retry) {
      originalRequest._retry = true;

      try {
        const refreshToken = Cookies.get("refresh_token");
        const { data } = await axios.post(
          \`\${API_URL}/api/auth/refresh\`,
          { refresh_token: refreshToken }
        );

        Cookies.set("access_token", data.data.access_token);
        Cookies.set("refresh_token", data.data.refresh_token);

        originalRequest.headers.Authorization =
          \`Bearer \${data.data.access_token}\`;
        return apiClient(originalRequest);
      } catch {
        Cookies.remove("access_token");
        Cookies.remove("refresh_token");
        window.location.href = "/login";
      }
    }

    return Promise.reject(error);
  }
);`}</pre>
                </div>
              </div>

              {/* React Query Setup */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  React Query Setup
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  The React Query client is configured with sensible defaults: no automatic refetching on
                  window focus in development, stale time of 30 seconds, and retry of 1 attempt on failure.
                </p>

                <div className="rounded-xl border border-border/40 bg-card/80 overflow-hidden">
                  <div className="flex items-center gap-2 px-4 py-2.5 border-b border-border/30 bg-accent/30">
                    <span className="text-[11px] font-mono text-muted-foreground/40">apps/web/lib/query-client.ts</span>
                  </div>
                  <pre className="p-5 text-sm font-mono text-foreground/80 overflow-x-auto">{`import { QueryClient } from "@tanstack/react-query";

export const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      staleTime: 30 * 1000,         // 30 seconds
      retry: 1,
      refetchOnWindowFocus: false,
    },
    mutations: {
      retry: 0,
    },
  },
});`}</pre>
                </div>
              </div>

              {/* Tailwind + shadcn/ui */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Tailwind CSS + shadcn/ui
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  The web app uses Tailwind CSS with a custom dark theme and shadcn/ui components. All CSS
                  variables are defined in <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">globals.css</code> and
                  mapped to Tailwind utility classes in <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">tailwind.config.ts</code>.
                </p>

                <div className="rounded-xl border border-border/40 bg-card/80 overflow-hidden mb-6">
                  <div className="flex items-center gap-2 px-4 py-2.5 border-b border-border/30 bg-accent/30">
                    <span className="text-[11px] font-mono text-muted-foreground/40">apps/web/app/globals.css</span>
                  </div>
                  <pre className="p-5 text-sm font-mono text-foreground/80 overflow-x-auto">{`:root {
  --bg-primary:    #0a0a0f;
  --bg-secondary:  #111118;
  --bg-tertiary:   #1a1a24;
  --bg-elevated:   #22222e;
  --bg-hover:      #2a2a38;
  --border:        #2a2a3a;
  --text-primary:  #e8e8f0;
  --text-secondary:#9090a8;
  --text-muted:    #606078;
  --accent:        #6c5ce7;
  --accent-hover:  #7c6cf7;
  --success:       #00b894;
  --danger:        #ff6b6b;
  --warning:       #fdcb6e;
  --info:          #74b9ff;
}`}</pre>
                </div>

                <div className="rounded-xl border border-border/40 bg-card/80 overflow-hidden">
                  <div className="flex items-center gap-2 px-4 py-2.5 border-b border-border/30 bg-accent/30">
                    <span className="text-[11px] font-mono text-muted-foreground/40">apps/web/tailwind.config.ts</span>
                  </div>
                  <pre className="p-5 text-sm font-mono text-foreground/80 overflow-x-auto">{`import type { Config } from "tailwindcss";

const config: Config = {
  darkMode: "class",
  content: [
    "./app/**/*.{ts,tsx}",
    "./components/**/*.{ts,tsx}",
  ],
  theme: {
    extend: {
      colors: {
        background: "var(--bg-primary)",
        "bg-secondary": "var(--bg-secondary)",
        "bg-tertiary": "var(--bg-tertiary)",
        border: "var(--border)",
        foreground: "var(--text-primary)",
        accent: {
          DEFAULT: "var(--accent)",
          hover: "var(--accent-hover)",
        },
        success: "var(--success)",
        danger: "var(--danger)",
      },
      fontFamily: {
        sans: ["DM Sans", "system-ui", "sans-serif"],
        mono: ["JetBrains Mono", "monospace"],
      },
    },
  },
  plugins: [require("tailwindcss-animate")],
};

export default config;`}</pre>
                </div>
              </div>

              {/* Dark Theme Default */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Dark Theme by Default
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Grit uses a dark-first design philosophy. The <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">dark</code> class
                  is applied to the root <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">&lt;html&gt;</code> element in the root layout.
                  The color palette is inspired by premium tools like Linear, Vercel Dashboard, and Raycast.
                </p>

                <div className="rounded-lg border border-border/30 bg-card/30 overflow-hidden">
                  <table className="w-full text-sm">
                    <thead>
                      <tr className="border-b border-border/30 bg-accent/20">
                        <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Token</th>
                        <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Value</th>
                        <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Usage</th>
                      </tr>
                    </thead>
                    <tbody className="text-muted-foreground">
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-mono text-xs">--bg-primary</td>
                        <td className="px-4 py-2.5 font-mono text-xs">#0a0a0f</td>
                        <td className="px-4 py-2.5">Page background</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-mono text-xs">--bg-secondary</td>
                        <td className="px-4 py-2.5 font-mono text-xs">#111118</td>
                        <td className="px-4 py-2.5">Cards, sidebar</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-mono text-xs">--accent</td>
                        <td className="px-4 py-2.5 font-mono text-xs">#6c5ce7</td>
                        <td className="px-4 py-2.5">Buttons, links, highlights</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-mono text-xs">--text-primary</td>
                        <td className="px-4 py-2.5 font-mono text-xs">#e8e8f0</td>
                        <td className="px-4 py-2.5">Headings, body text</td>
                      </tr>
                      <tr>
                        <td className="px-4 py-2.5 font-mono text-xs">--danger</td>
                        <td className="px-4 py-2.5 font-mono text-xs">#ff6b6b</td>
                        <td className="px-4 py-2.5">Error states, delete actions</td>
                      </tr>
                    </tbody>
                  </table>
                </div>
              </div>

              {/* Running the Web App */}
              <div className="mb-10">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Running the Web App
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Start the web app in development mode. Make sure the Go API is running first.
                </p>

                <div className="rounded-xl border border-border/40 bg-card/80 overflow-hidden glow-purple-sm mb-4">
                  <div className="flex items-center gap-2 px-4 py-2.5 border-b border-border/30 bg-accent/30">
                    <div className="flex items-center gap-1.5">
                      <div className="h-2.5 w-2.5 rounded-full bg-[#ff5f57]" />
                      <div className="h-2.5 w-2.5 rounded-full bg-[#febc2e]" />
                      <div className="h-2.5 w-2.5 rounded-full bg-[#28c840]" />
                    </div>
                    <span className="ml-2 text-[11px] font-mono text-muted-foreground/40">terminal</span>
                  </div>
                  <div className="p-5 font-mono text-sm space-y-2">
                    <div>
                      <span className="text-primary/50 select-none">$ </span>
                      <span className="text-foreground/80">cd apps/web && pnpm dev</span>
                    </div>
                  </div>
                </div>

                <p className="text-sm text-muted-foreground/60">
                  Or use Turborepo to start all apps at once from the project root:
                </p>

                <div className="rounded-xl border border-border/40 bg-card/80 overflow-hidden mt-3">
                  <div className="flex items-center gap-2 px-4 py-2.5 border-b border-border/30 bg-accent/30">
                    <div className="flex items-center gap-1.5">
                      <div className="h-2.5 w-2.5 rounded-full bg-[#ff5f57]" />
                      <div className="h-2.5 w-2.5 rounded-full bg-[#febc2e]" />
                      <div className="h-2.5 w-2.5 rounded-full bg-[#28c840]" />
                    </div>
                    <span className="ml-2 text-[11px] font-mono text-muted-foreground/40">terminal</span>
                  </div>
                  <div className="p-5 font-mono text-sm">
                    <div>
                      <span className="text-primary/50 select-none">$ </span>
                      <span className="text-foreground/80">turbo dev</span>
                    </div>
                  </div>
                </div>
              </div>
            </div>

            {/* Nav */}
            <div className="flex items-center justify-between pt-6 border-t border-border/30">
              <Button variant="ghost" size="sm" asChild className="text-muted-foreground/60 hover:text-foreground">
                <Link href="/docs/admin/widgets" className="gap-1.5">
                  <ArrowLeft className="h-3.5 w-3.5" />
                  Dashboard & Widgets
                </Link>
              </Button>
              <Button variant="ghost" size="sm" asChild className="text-muted-foreground/60 hover:text-foreground">
                <Link href="/docs/frontend/hooks" className="gap-1.5">
                  React Query Hooks
                  <ArrowRight className="h-3.5 w-3.5" />
                </Link>
              </Button>
            </div>
          </div>
        </div>
      </main>
    </div>
  )
}
