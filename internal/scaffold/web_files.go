package scaffold

import (
	"fmt"
	"path/filepath"
)

func writeWebFiles(root string, opts Options) error {
	webRoot := filepath.Join(root, "apps", "web")

	files := map[string]string{
		filepath.Join(webRoot, "package.json"):                        webPackageJSON(opts),
		filepath.Join(webRoot, "next.config.ts"):                      webNextConfig(),
		filepath.Join(webRoot, "tailwind.config.ts"):                  webTailwindConfig(),
		filepath.Join(webRoot, "postcss.config.js"):                   webPostCSSConfig(),
		filepath.Join(webRoot, "tsconfig.json"):                       webTSConfig(),
		filepath.Join(webRoot, "app", "globals.css"):                  webGlobalCSS(),
		filepath.Join(webRoot, "app", "layout.tsx"):                   webRootLayout(opts),
		filepath.Join(webRoot, "app", "page.tsx"):                     webLandingPage(opts),
		filepath.Join(webRoot, "app", "error.tsx"):                    webErrorPage(),
		filepath.Join(webRoot, "app", "not-found.tsx"):                webNotFoundPage(),
		filepath.Join(webRoot, "app", "global-error.tsx"):             webGlobalErrorPage(),
		filepath.Join(webRoot, "lib", "utils.ts"):                     webUtils(),
		filepath.Join(webRoot, "components", "navbar.tsx"):            webNavbar(opts),
		filepath.Join(webRoot, "components", "footer.tsx"):            webFooter(opts),
		filepath.Join(webRoot, "components", "providers.tsx"):         webProviders(),
		filepath.Join(webRoot, "lib", "api.ts"):                       webAPIClient(),
		filepath.Join(webRoot, "hooks", "use-blogs.ts"):               webUseBlogsHook(),
		filepath.Join(webRoot, "app", "blog", "page.tsx"):             webBlogListPage(),
		filepath.Join(webRoot, "app", "blog", "[slug]", "page.tsx"):   webBlogDetailPage(),
		filepath.Join(webRoot, "public", ".gitkeep"):                   "",
	}

	for path, content := range files {
		if err := writeFile(path, content); err != nil {
			return fmt.Errorf("writing %s: %w", path, err)
		}
	}

	return nil
}

func webPackageJSON(opts Options) string {
	return fmt.Sprintf(`{
  "name": "@%s/web",
  "version": "0.1.0",
  "private": true,
  "scripts": {
    "dev": "next dev --port 3000",
    "build": "next build",
    "start": "next start",
    "lint": "next lint"
  },
  "dependencies": {
    "@tanstack/react-query": "^5.17.0",
    "axios": "^1.6.0",
    "class-variance-authority": "^0.7.0",
    "clsx": "^2.1.0",
    "lucide-react": "^0.303.0",
    "next": "^16.1.6",
    "react": "^19.0.0",
    "react-dom": "^19.0.0",
    "react-hook-form": "^7.49.0",
    "tailwind-merge": "^2.2.0",
    "tailwindcss-animate": "^1.0.7"
  },
  "devDependencies": {
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

func webNextConfig() string {
	return `import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  output: "standalone",
  reactStrictMode: true,
};

export default nextConfig;
`
}

func webTailwindConfig() string {
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

func webPostCSSConfig() string {
	return `module.exports = {
  plugins: {
    tailwindcss: {},
    autoprefixer: {},
  },
};
`
}

func webTSConfig() string {
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

func webGlobalCSS() string {
	return `@tailwind base;
@tailwind components;
@tailwind utilities;

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

::-webkit-scrollbar-thumb:hover {
  background: var(--text-muted);
}

/* Blog prose styles for rendered HTML content */
.prose-blog {
  color: var(--text-primary);
  font-size: 1.0625rem;
  line-height: 1.8;
}

.prose-blog h1 {
  font-size: 2rem;
  font-weight: 700;
  margin-top: 2.5rem;
  margin-bottom: 1rem;
  letter-spacing: -0.025em;
  color: var(--text-primary);
}

.prose-blog h2 {
  font-size: 1.5rem;
  font-weight: 700;
  margin-top: 2rem;
  margin-bottom: 0.75rem;
  letter-spacing: -0.025em;
  color: var(--text-primary);
}

.prose-blog h3 {
  font-size: 1.25rem;
  font-weight: 600;
  margin-top: 1.75rem;
  margin-bottom: 0.5rem;
  color: var(--text-primary);
}

.prose-blog h4 {
  font-size: 1.125rem;
  font-weight: 600;
  margin-top: 1.5rem;
  margin-bottom: 0.5rem;
  color: var(--text-primary);
}

.prose-blog p {
  margin-bottom: 1.25rem;
  color: var(--text-secondary);
}

.prose-blog a {
  color: var(--accent);
  text-decoration: underline;
  text-underline-offset: 2px;
  transition: color 0.15s ease;
}

.prose-blog a:hover {
  color: var(--accent-hover);
}

.prose-blog strong {
  font-weight: 600;
  color: var(--text-primary);
}

.prose-blog em {
  font-style: italic;
}

.prose-blog ul {
  list-style: disc;
  padding-left: 1.5rem;
  margin-bottom: 1.25rem;
}

.prose-blog ol {
  list-style: decimal;
  padding-left: 1.5rem;
  margin-bottom: 1.25rem;
}

.prose-blog li {
  margin-bottom: 0.375rem;
  color: var(--text-secondary);
}

.prose-blog li::marker {
  color: var(--text-muted);
}

.prose-blog blockquote {
  border-left: 3px solid var(--accent);
  padding-left: 1rem;
  margin: 1.5rem 0;
  font-style: italic;
  color: var(--text-secondary);
  background: var(--bg-secondary);
  padding: 1rem 1.25rem;
  border-radius: 0 0.5rem 0.5rem 0;
}

.prose-blog pre {
  background: var(--bg-secondary);
  border: 1px solid var(--border);
  border-radius: 0.75rem;
  padding: 1.25rem;
  overflow-x: auto;
  margin: 1.5rem 0;
  font-family: var(--font-jetbrains-mono), ui-monospace, monospace;
  font-size: 0.875rem;
  line-height: 1.7;
}

.prose-blog code {
  font-family: var(--font-jetbrains-mono), ui-monospace, monospace;
  font-size: 0.875em;
  background: var(--bg-elevated);
  padding: 0.15rem 0.4rem;
  border-radius: 0.25rem;
  border: 1px solid var(--border);
  color: var(--accent);
}

.prose-blog pre code {
  background: transparent;
  padding: 0;
  border: none;
  border-radius: 0;
  font-size: inherit;
  color: var(--text-primary);
}

.prose-blog img {
  max-width: 100%;
  height: auto;
  border-radius: 0.75rem;
  margin: 1.5rem 0;
  border: 1px solid var(--border);
}

.prose-blog hr {
  border: none;
  border-top: 1px solid var(--border);
  margin: 2rem 0;
}

.prose-blog table {
  width: 100%;
  border-collapse: collapse;
  margin: 1.5rem 0;
}

.prose-blog th {
  text-align: left;
  padding: 0.75rem 1rem;
  border-bottom: 2px solid var(--border);
  font-weight: 600;
  color: var(--text-primary);
  font-size: 0.875rem;
}

.prose-blog td {
  padding: 0.75rem 1rem;
  border-bottom: 1px solid var(--border);
  color: var(--text-secondary);
  font-size: 0.875rem;
}

.prose-blog tr:last-child td {
  border-bottom: none;
}
`
}

func webRootLayout(opts Options) string {
	return fmt.Sprintf(`import type { Metadata } from "next";
import { Onest, JetBrains_Mono } from "next/font/google";
import { Providers } from "@/components/providers";
import { Navbar } from "@/components/navbar";
import { Footer } from "@/components/footer";
import "./globals.css";

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
  title: "%s — Go + React. Built with Grit.",
  description: "A full-stack framework that combines Go backend with Next.js frontend. Build fast, ship faster.",
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="en" className="dark">
      <body className={`+"`"+`${onest.variable} ${jetbrainsMono.variable} font-sans dark antialiased`+"`"+`}>
        <Providers>
          <Navbar />
          <main className="min-h-screen">{children}</main>
          <Footer />
        </Providers>
      </body>
    </html>
  );
}
`, opts.ProjectName)
}

func webLandingPage(opts Options) string {
	return `"use client";

import { ArrowRight } from "lucide-react";
import Link from "next/link";
import { usePublicBlogs } from "@/hooks/use-blogs";

const DOCS_URL = "https://grit-vert.vercel.app/docs";

export default function HomePage() {
  const { data, isLoading } = usePublicBlogs(1, 3);
  const blogs = data?.blogs || [];

  return (
    <div className="flex flex-col">
      {/* Hero */}
      <section className="flex items-center justify-center">
        <div className="mx-auto max-w-2xl px-6 py-24 text-center">
          <div className="mb-8 inline-flex items-center gap-2 rounded-full border border-border bg-bg-secondary px-4 py-1.5 text-sm text-text-secondary">
            <span className="h-2 w-2 rounded-full bg-success animate-pulse" />
            <span>Go + React Full-Stack Framework</span>
          </div>

          <h1 className="text-5xl font-bold tracking-tight sm:text-6xl lg:text-7xl">
            <span className="text-accent">Grit</span>
          </h1>

          <p className="mt-6 text-lg text-text-secondary leading-relaxed max-w-lg mx-auto">
            The full-stack meta-framework that fuses Go, React, and a
            Filament-like admin panel. Scaffold entire projects, generate
            resources, and ship fast.
          </p>

          <div className="mt-10 flex flex-col sm:flex-row items-center justify-center gap-4">
            <a
              href={` + "`" + `${DOCS_URL}/getting-started/quick-start` + "`" + `}
              target="_blank"
              rel="noopener noreferrer"
              className="flex items-center gap-2 rounded-lg bg-accent px-6 py-3 text-sm font-semibold text-white hover:bg-accent-hover transition-colors"
            >
              Get Started <ArrowRight className="h-4 w-4" />
            </a>
            <a
              href={` + "`" + `${DOCS_URL}` + "`" + `}
              target="_blank"
              rel="noopener noreferrer"
              className="flex items-center gap-2 rounded-lg border border-border bg-bg-secondary px-6 py-3 text-sm font-semibold text-foreground hover:bg-bg-hover transition-colors"
            >
              Read the Docs
            </a>
          </div>

          {/* Terminal snippet */}
          <div className="mt-16 mx-auto max-w-md rounded-xl border border-border bg-bg-secondary shadow-2xl overflow-hidden text-left">
            <div className="flex items-center gap-2 border-b border-border px-4 py-2.5">
              <div className="h-2.5 w-2.5 rounded-full bg-danger/60" />
              <div className="h-2.5 w-2.5 rounded-full bg-warning/60" />
              <div className="h-2.5 w-2.5 rounded-full bg-success/60" />
              <span className="ml-2 text-[11px] text-text-muted font-mono">terminal</span>
            </div>
            <div className="p-5 font-mono text-sm space-y-1.5">
              <p><span className="text-success select-none">$ </span><span className="text-foreground">grit new my-saas</span></p>
              <p><span className="text-success select-none">$ </span><span className="text-foreground">cd my-saas && docker compose up -d</span></p>
              <p><span className="text-success select-none">$ </span><span className="text-foreground">pnpm dev</span></p>
              <p className="text-success pt-1">Ready on http://localhost:3000</p>
            </div>
          </div>
        </div>
      </section>

      {/* Recent Posts */}
      <section className="border-t border-border/50 bg-bg-secondary/30">
        <div className="mx-auto max-w-5xl px-6 py-20">
          <div className="flex items-center justify-between mb-10">
            <div>
              <h2 className="text-2xl font-bold tracking-tight">Recent Posts</h2>
              <p className="mt-1 text-sm text-text-secondary">Latest articles and updates</p>
            </div>
            <Link
              href="/blog"
              className="text-sm text-accent hover:text-accent-hover transition-colors font-medium flex items-center gap-1"
            >
              View all <ArrowRight className="h-3.5 w-3.5" />
            </Link>
          </div>

          {isLoading ? (
            <div className="grid gap-6 sm:grid-cols-2 lg:grid-cols-3">
              {[1, 2, 3].map((i) => (
                <div key={i} className="rounded-xl border border-border bg-bg-elevated overflow-hidden animate-pulse">
                  <div className="h-48 bg-bg-hover" />
                  <div className="p-5 space-y-3">
                    <div className="h-4 bg-bg-hover rounded w-3/4" />
                    <div className="h-3 bg-bg-hover rounded w-full" />
                    <div className="h-3 bg-bg-hover rounded w-2/3" />
                  </div>
                </div>
              ))}
            </div>
          ) : blogs.length > 0 ? (
            <div className="grid gap-6 sm:grid-cols-2 lg:grid-cols-3">
              {blogs.map((blog: { id: number; title: string; slug: string; image: string; excerpt: string; published_at: string | null; created_at: string }) => (
                <Link
                  key={blog.id}
                  href={` + "`" + `/blog/${blog.slug}` + "`" + `}
                  className="group rounded-xl border border-border bg-bg-elevated overflow-hidden hover:border-accent/40 hover:shadow-lg hover:shadow-accent/5 transition-all duration-300"
                >
                  <div className="h-48 bg-bg-hover overflow-hidden">
                    {blog.image ? (
                      <img
                        src={blog.image}
                        alt={blog.title}
                        className="h-full w-full object-cover group-hover:scale-105 transition-transform duration-500"
                      />
                    ) : (
                      <div className="h-full w-full flex items-center justify-center bg-gradient-to-br from-accent/10 to-accent/5">
                        <span className="text-4xl font-bold text-accent/20">{blog.title.charAt(0)}</span>
                      </div>
                    )}
                  </div>
                  <div className="p-5">
                    <p className="text-xs text-text-muted mb-2">
                      {new Date(blog.published_at || blog.created_at).toLocaleDateString("en-US", {
                        month: "short",
                        day: "numeric",
                        year: "numeric",
                      })}
                    </p>
                    <h3 className="font-semibold text-foreground group-hover:text-accent transition-colors line-clamp-2">
                      {blog.title}
                    </h3>
                    {blog.excerpt && (
                      <p className="mt-2 text-sm text-text-secondary line-clamp-2">{blog.excerpt}</p>
                    )}
                  </div>
                </Link>
              ))}
            </div>
          ) : (
            <div className="text-center py-12">
              <p className="text-text-muted text-sm">No blog posts yet. Create your first post in the admin panel.</p>
            </div>
          )}
        </div>
      </section>
    </div>
  );
}
`
}

func webUtils() string {
	return `import { clsx, type ClassValue } from "clsx";
import { twMerge } from "tailwind-merge";

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}
`
}

func webNavbar(opts Options) string {
	return `"use client";

import { useState } from "react";
import Link from "next/link";
import { usePathname } from "next/navigation";
import { Menu, X, Github } from "lucide-react";

const ADMIN_URL = process.env.NEXT_PUBLIC_ADMIN_URL || "http://localhost:3001";
const DOCS_URL = "https://grit-vert.vercel.app/docs";

const navLinks = [
  { href: "/", label: "Home" },
  { href: "/blog", label: "Blog" },
];

export function Navbar() {
  const [mobileOpen, setMobileOpen] = useState(false);
  const pathname = usePathname();

  return (
    <nav className="sticky top-0 z-50 border-b border-border/50 bg-background/80 backdrop-blur-lg">
      <div className="mx-auto flex h-16 max-w-5xl items-center justify-between px-6">
        {/* Logo */}
        <Link href="/" className="flex items-center gap-2.5">
          <div className="flex h-8 w-8 items-center justify-center rounded-lg bg-accent/15 border border-accent/20">
            <span className="text-accent font-mono font-bold text-sm">G</span>
          </div>
          <span className="text-lg font-bold tracking-tight">` + opts.ProjectName + `</span>
        </Link>

        {/* Desktop nav */}
        <div className="hidden md:flex items-center gap-6">
          {navLinks.map((link) => (
            <Link
              key={link.href}
              href={link.href}
              className={` + "`" + `text-sm transition-colors ${
                pathname === link.href
                  ? "text-foreground font-medium"
                  : "text-text-secondary hover:text-foreground"
              }` + "`" + `}
            >
              {link.label}
            </Link>
          ))}
          <a
            href={DOCS_URL}
            target="_blank"
            rel="noopener noreferrer"
            className="text-sm text-text-secondary hover:text-foreground transition-colors"
          >
            Docs
          </a>
          <a
            href="https://github.com/MUKE-coder/grit"
            target="_blank"
            rel="noopener noreferrer"
            className="text-text-secondary hover:text-foreground transition-colors"
          >
            <Github className="h-5 w-5" />
          </a>
          <a
            href={` + "`" + `${ADMIN_URL}/login` + "`" + `}
            className="rounded-lg bg-accent px-4 py-2 text-sm font-medium text-white hover:bg-accent-hover transition-colors"
          >
            Admin
          </a>
        </div>

        {/* Mobile hamburger */}
        <button
          onClick={() => setMobileOpen(!mobileOpen)}
          className="md:hidden p-2 text-text-secondary hover:text-foreground transition-colors"
          aria-label="Toggle menu"
        >
          {mobileOpen ? <X className="h-5 w-5" /> : <Menu className="h-5 w-5" />}
        </button>
      </div>

      {/* Mobile menu */}
      {mobileOpen && (
        <div className="md:hidden border-t border-border/50 bg-background/95 backdrop-blur-lg">
          <div className="mx-auto max-w-5xl px-6 py-4 flex flex-col gap-3">
            {navLinks.map((link) => (
              <Link
                key={link.href}
                href={link.href}
                onClick={() => setMobileOpen(false)}
                className={` + "`" + `text-sm py-2 transition-colors ${
                  pathname === link.href
                    ? "text-foreground font-medium"
                    : "text-text-secondary hover:text-foreground"
                }` + "`" + `}
              >
                {link.label}
              </Link>
            ))}
            <a
              href={DOCS_URL}
              target="_blank"
              rel="noopener noreferrer"
              className="text-sm py-2 text-text-secondary hover:text-foreground transition-colors"
            >
              Docs
            </a>
            <a
              href="https://github.com/MUKE-coder/grit"
              target="_blank"
              rel="noopener noreferrer"
              className="text-sm py-2 text-text-secondary hover:text-foreground transition-colors"
            >
              GitHub
            </a>
            <a
              href={` + "`" + `${ADMIN_URL}/login` + "`" + `}
              className="mt-2 rounded-lg bg-accent px-4 py-2.5 text-sm font-medium text-white text-center hover:bg-accent-hover transition-colors"
            >
              Admin
            </a>
          </div>
        </div>
      )}
    </nav>
  );
}
`
}

func webFooter(opts Options) string {
	return `import Link from "next/link";

const DOCS_URL = "https://grit-vert.vercel.app/docs";

export function Footer() {
  return (
    <footer className="border-t border-border/50 bg-background">
      <div className="mx-auto max-w-5xl px-6 py-8">
        <div className="flex flex-col sm:flex-row items-center justify-between gap-4">
          <div className="flex items-center gap-2 text-sm text-text-muted">
            <span className="font-semibold text-text-secondary">` + opts.ProjectName + `</span>
            <span className="text-border">·</span>
            <span>Built with Grit</span>
          </div>
          <div className="flex items-center gap-6 text-sm text-text-muted">
            <a
              href="https://github.com/MUKE-coder/grit"
              target="_blank"
              rel="noopener noreferrer"
              className="hover:text-foreground transition-colors"
            >
              GitHub
            </a>
            <a
              href={DOCS_URL}
              target="_blank"
              rel="noopener noreferrer"
              className="hover:text-foreground transition-colors"
            >
              Documentation
            </a>
          </div>
        </div>
        <div className="mt-6 pt-6 border-t border-border/30 text-center">
          <p className="text-xs text-text-muted">
            &copy; {new Date().getFullYear()} ` + opts.ProjectName + `. All rights reserved.
          </p>
        </div>
      </div>
    </footer>
  );
}
`
}

func webProviders() string {
	return `"use client";

import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { useState } from "react";

export function Providers({ children }: { children: React.ReactNode }) {
  const [queryClient] = useState(
    () =>
      new QueryClient({
        defaultOptions: {
          queries: {
            staleTime: 60 * 1000,
            retry: 1,
          },
        },
      })
  );

  return (
    <QueryClientProvider client={queryClient}>{children}</QueryClientProvider>
  );
}
`
}

func webAPIClient() string {
	return `import axios from "axios";

const API_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";

export const api = axios.create({
  baseURL: API_URL,
  headers: {
    "Content-Type": "application/json",
  },
});
`
}

func webUseBlogsHook() string {
	return `"use client";

import { useQuery } from "@tanstack/react-query";
import { api } from "@/lib/api";

interface Blog {
  id: number;
  title: string;
  slug: string;
  content: string;
  image: string;
  excerpt: string;
  published: boolean;
  published_at: string | null;
  created_at: string;
}

interface BlogMeta {
  total: number;
  page: number;
  page_size: number;
  pages: number;
}

export function usePublicBlogs(page = 1, pageSize = 20) {
  return useQuery({
    queryKey: ["public-blogs", page, pageSize],
    queryFn: async () => {
      const { data } = await api.get(
        ` + "`" + `/api/blogs?page=${page}&page_size=${pageSize}` + "`" + `
      );
      return {
        blogs: (data.data || []) as Blog[],
        meta: data.meta as BlogMeta | undefined,
      };
    },
  });
}

export function usePublicBlog(slug: string) {
  return useQuery({
    queryKey: ["public-blog", slug],
    queryFn: async () => {
      const { data } = await api.get(` + "`" + `/api/blogs/${slug}` + "`" + `);
      return data.data as Blog;
    },
    enabled: !!slug,
  });
}
`
}

func webBlogListPage() string {
	return `"use client";

import { useState } from "react";
import Link from "next/link";
import { ChevronLeft, ChevronRight } from "lucide-react";
import { usePublicBlogs } from "@/hooks/use-blogs";

export default function BlogListPage() {
  const [page, setPage] = useState(1);
  const { data, isLoading } = usePublicBlogs(page, 9);
  const blogs = data?.blogs || [];
  const meta = data?.meta;

  return (
    <div className="mx-auto max-w-5xl px-6 py-16">
      {/* Header */}
      <div className="mb-12">
        <h1 className="text-4xl font-bold tracking-tight">Blog</h1>
        <p className="mt-2 text-text-secondary">
          Insights, tutorials, and updates from the team.
        </p>
      </div>

      {/* Blog grid */}
      {isLoading ? (
        <div className="grid gap-6 sm:grid-cols-2 lg:grid-cols-3">
          {Array.from({ length: 6 }).map((_, i) => (
            <div
              key={i}
              className="rounded-xl border border-border bg-bg-elevated overflow-hidden animate-pulse"
            >
              <div className="h-52 bg-bg-hover" />
              <div className="p-5 space-y-3">
                <div className="h-3 bg-bg-hover rounded w-1/3" />
                <div className="h-5 bg-bg-hover rounded w-3/4" />
                <div className="h-3 bg-bg-hover rounded w-full" />
                <div className="h-3 bg-bg-hover rounded w-2/3" />
              </div>
            </div>
          ))}
        </div>
      ) : blogs.length > 0 ? (
        <>
          <div className="grid gap-6 sm:grid-cols-2 lg:grid-cols-3">
            {blogs.map((blog) => (
              <Link
                key={blog.id}
                href={` + "`" + `/blog/${blog.slug}` + "`" + `}
                className="group rounded-xl border border-border bg-bg-elevated overflow-hidden hover:border-accent/40 hover:shadow-lg hover:shadow-accent/5 transition-all duration-300"
              >
                <div className="h-52 bg-bg-hover overflow-hidden">
                  {blog.image ? (
                    <img
                      src={blog.image}
                      alt={blog.title}
                      className="h-full w-full object-cover group-hover:scale-105 transition-transform duration-500"
                    />
                  ) : (
                    <div className="h-full w-full flex items-center justify-center bg-gradient-to-br from-accent/10 to-accent/5">
                      <span className="text-5xl font-bold text-accent/20">
                        {blog.title.charAt(0)}
                      </span>
                    </div>
                  )}
                </div>
                <div className="p-5">
                  <p className="text-xs text-text-muted mb-2.5">
                    {new Date(
                      blog.published_at || blog.created_at
                    ).toLocaleDateString("en-US", {
                      month: "short",
                      day: "numeric",
                      year: "numeric",
                    })}
                  </p>
                  <h2 className="font-semibold text-foreground group-hover:text-accent transition-colors line-clamp-2 text-lg leading-snug">
                    {blog.title}
                  </h2>
                  {blog.excerpt && (
                    <p className="mt-2.5 text-sm text-text-secondary line-clamp-3 leading-relaxed">
                      {blog.excerpt}
                    </p>
                  )}
                  <span className="mt-4 inline-block text-xs font-medium text-accent group-hover:text-accent-hover transition-colors">
                    Read more &rarr;
                  </span>
                </div>
              </Link>
            ))}
          </div>

          {/* Pagination */}
          {meta && meta.pages > 1 && (
            <div className="mt-12 flex items-center justify-center gap-2">
              <button
                onClick={() => setPage((p) => Math.max(1, p - 1))}
                disabled={page <= 1}
                className="flex items-center gap-1 rounded-lg border border-border bg-bg-elevated px-3 py-2 text-sm text-text-secondary hover:bg-bg-hover hover:text-foreground transition-colors disabled:opacity-40 disabled:cursor-not-allowed"
              >
                <ChevronLeft className="h-4 w-4" />
                Previous
              </button>
              <div className="flex items-center gap-1 px-3">
                {Array.from({ length: meta.pages }).map((_, i) => (
                  <button
                    key={i + 1}
                    onClick={() => setPage(i + 1)}
                    className={` + "`" + `h-8 w-8 rounded-lg text-sm font-medium transition-colors ${
                      page === i + 1
                        ? "bg-accent text-white"
                        : "text-text-secondary hover:bg-bg-hover hover:text-foreground"
                    }` + "`" + `}
                  >
                    {i + 1}
                  </button>
                ))}
              </div>
              <button
                onClick={() => setPage((p) => Math.min(meta.pages, p + 1))}
                disabled={page >= meta.pages}
                className="flex items-center gap-1 rounded-lg border border-border bg-bg-elevated px-3 py-2 text-sm text-text-secondary hover:bg-bg-hover hover:text-foreground transition-colors disabled:opacity-40 disabled:cursor-not-allowed"
              >
                Next
                <ChevronRight className="h-4 w-4" />
              </button>
            </div>
          )}
        </>
      ) : (
        <div className="text-center py-20">
          <div className="mx-auto mb-4 flex h-16 w-16 items-center justify-center rounded-2xl bg-bg-elevated border border-border">
            <span className="text-2xl text-text-muted">&#9998;</span>
          </div>
          <h3 className="text-lg font-semibold text-foreground">No posts yet</h3>
          <p className="mt-1 text-sm text-text-muted">
            Blog posts will appear here once published from the admin panel.
          </p>
        </div>
      )}
    </div>
  );
}
`
}

func webBlogDetailPage() string {
	return `"use client";

import { useParams } from "next/navigation";
import Link from "next/link";
import { ArrowLeft, Calendar } from "lucide-react";
import { usePublicBlog } from "@/hooks/use-blogs";

export default function BlogDetailPage() {
  const params = useParams();
  const slug = typeof params.slug === "string" ? params.slug : "";
  const { data: blog, isLoading, error } = usePublicBlog(slug);

  if (isLoading) {
    return (
      <div className="mx-auto max-w-3xl px-6 py-16 animate-pulse">
        <div className="h-4 bg-bg-hover rounded w-24 mb-8" />
        <div className="h-8 bg-bg-hover rounded w-3/4 mb-4" />
        <div className="h-4 bg-bg-hover rounded w-1/3 mb-12" />
        <div className="aspect-[2/1] bg-bg-hover rounded-xl mb-12" />
        <div className="space-y-4">
          <div className="h-4 bg-bg-hover rounded w-full" />
          <div className="h-4 bg-bg-hover rounded w-full" />
          <div className="h-4 bg-bg-hover rounded w-5/6" />
          <div className="h-4 bg-bg-hover rounded w-full" />
          <div className="h-4 bg-bg-hover rounded w-4/6" />
        </div>
      </div>
    );
  }

  if (error || !blog) {
    return (
      <div className="mx-auto max-w-3xl px-6 py-16 text-center">
        <div className="mx-auto mb-4 flex h-16 w-16 items-center justify-center rounded-2xl bg-bg-elevated border border-border">
          <span className="text-2xl text-text-muted">404</span>
        </div>
        <h1 className="text-xl font-semibold text-foreground">Post not found</h1>
        <p className="mt-2 text-sm text-text-muted">
          The blog post you're looking for doesn't exist or has been removed.
        </p>
        <Link
          href="/blog"
          className="mt-6 inline-flex items-center gap-2 rounded-lg bg-accent px-4 py-2 text-sm font-medium text-white hover:bg-accent-hover transition-colors"
        >
          <ArrowLeft className="h-4 w-4" />
          Back to Blog
        </Link>
      </div>
    );
  }

  return (
    <article className="mx-auto max-w-3xl px-6 py-16">
      {/* Back link */}
      <Link
        href="/blog"
        className="inline-flex items-center gap-1.5 text-sm text-text-secondary hover:text-foreground transition-colors mb-8"
      >
        <ArrowLeft className="h-4 w-4" />
        Back to Blog
      </Link>

      {/* Title and meta */}
      <header className="mb-10">
        <h1 className="text-3xl sm:text-4xl font-bold tracking-tight leading-tight">
          {blog.title}
        </h1>
        <div className="mt-4 flex items-center gap-2 text-sm text-text-muted">
          <Calendar className="h-4 w-4" />
          <time dateTime={blog.published_at || blog.created_at}>
            {new Date(blog.published_at || blog.created_at).toLocaleDateString(
              "en-US",
              {
                month: "long",
                day: "numeric",
                year: "numeric",
              }
            )}
          </time>
        </div>
      </header>

      {/* Cover image */}
      {blog.image && (
        <div className="mb-12 rounded-xl overflow-hidden border border-border">
          <img
            src={blog.image}
            alt={blog.title}
            className="w-full h-auto object-cover"
          />
        </div>
      )}

      {/* Content */}
      <div
        className="prose-blog"
        dangerouslySetInnerHTML={{ __html: blog.content }}
      />

      {/* Bottom nav */}
      <div className="mt-16 pt-8 border-t border-border/50">
        <Link
          href="/blog"
          className="inline-flex items-center gap-1.5 text-sm text-accent hover:text-accent-hover transition-colors font-medium"
        >
          <ArrowLeft className="h-4 w-4" />
          All posts
        </Link>
      </div>
    </article>
  );
}
`
}

func webErrorPage() string {
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
    <div className="flex min-h-[60vh] items-center justify-center px-4">
      <div className="w-full max-w-md text-center">
        <div className="mx-auto mb-6 flex h-16 w-16 items-center justify-center rounded-full bg-red-500/10 border border-red-500/20">
          <svg className="h-8 w-8 text-red-400" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
            <path strokeLinecap="round" strokeLinejoin="round" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
          </svg>
        </div>
        <h2 className="mb-2 text-2xl font-bold text-foreground">Something went wrong</h2>
        <p className="mb-6 text-muted-foreground">
          An unexpected error occurred. You can try again or go back.
        </p>
        {error.digest && (
          <p className="mb-4 text-xs text-muted-foreground/60 font-mono">Error ID: {error.digest}</p>
        )}
        <div className="flex gap-3 justify-center">
          <button
            onClick={() => window.history.back()}
            className="inline-flex items-center gap-2 rounded-lg border border-border px-4 py-2.5 text-sm font-medium text-foreground hover:bg-accent/50 transition-colors"
          >
            Go Back
          </button>
          <button
            onClick={reset}
            className="inline-flex items-center gap-2 rounded-lg bg-primary px-4 py-2.5 text-sm font-medium text-primary-foreground hover:bg-primary/90 transition-colors"
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

func webNotFoundPage() string {
	return `import Link from "next/link";

export default function NotFound() {
  return (
    <div className="flex min-h-[60vh] items-center justify-center px-4">
      <div className="w-full max-w-md text-center">
        <p className="mb-4 text-7xl font-bold text-primary">404</p>
        <h2 className="mb-2 text-2xl font-bold text-foreground">Page not found</h2>
        <p className="mb-8 text-muted-foreground">
          The page you&apos;re looking for doesn&apos;t exist or has been moved.
        </p>
        <div className="flex gap-3 justify-center">
          <Link
            href="/"
            className="inline-flex items-center gap-2 rounded-lg bg-primary px-4 py-2.5 text-sm font-medium text-primary-foreground hover:bg-primary/90 transition-colors"
          >
            Home
          </Link>
        </div>
      </div>
    </div>
  );
}
`
}

func webGlobalErrorPage() string {
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
