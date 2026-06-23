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
		filepath.Join(webRoot, "app", "blog", "page.tsx"):                   webBlogListPage(),
		filepath.Join(webRoot, "app", "blog", "[slug]", "page.tsx"):         webBlogDetailPage(),
		filepath.Join(webRoot, "app", "components", "page.tsx"):             webComponentsPage(opts),
		filepath.Join(webRoot, "app", "components", "[name]", "page.tsx"):   webComponentDetailPage(opts),
		// v3.31.20: public form-share page (Phase 2)
		filepath.Join(webRoot, "app", "forms", "[token]", "page.tsx"):       webPublicFormPage(),
		filepath.Join(webRoot, "public", ".gitkeep"):                         "",

		// Auth pages (v3.28.1: shell-driven, same per-theme treatment as admin)
		filepath.Join(webRoot, "app", "(auth)", "login", "page.tsx"):           webThemedLoginPage(),
		filepath.Join(webRoot, "app", "(auth)", "register", "page.tsx"):        webThemedRegisterPage(),
		filepath.Join(webRoot, "app", "(auth)", "forgot-password", "page.tsx"): webThemedForgotPasswordPage(),
		filepath.Join(webRoot, "app", "(auth)", "callback", "page.tsx"):        webAuthCallback(),
		filepath.Join(webRoot, "hooks", "use-auth.ts"):                         webUseAuth(),
		filepath.Join(webRoot, "lib", "auth-provider.tsx"):                     webAuthProvider(),

		// Theme-aware auth shells (v3.28.1) — same components as admin, copied
		// in so web's bundle can render without reaching across packages.
		filepath.Join(webRoot, "components", "auth", "AuthShell.tsx"):         adminAuthShellDispatcher(),
		filepath.Join(webRoot, "components", "auth", "AtlasAuthShell.tsx"):    adminAtlasAuthShell(),
		filepath.Join(webRoot, "components", "auth", "AuroraAuthShell.tsx"):   adminAuroraAuthShell(),
		filepath.Join(webRoot, "components", "auth", "PulseAuthShell.tsx"):    adminPulseAuthShell(),
		filepath.Join(webRoot, "components", "auth", "SocialAuthButtons.tsx"): adminAuthSocialButtons(),
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
    "dev": "rm -rf .next && next dev --port 3000",
    "build": "next build",
    "start": "next start",
    "lint": "next lint",
    "format": "prettier --write .",
    "test": "vitest run",
    "test:watch": "vitest",
    "test:ui": "vitest --ui"
  },
  "dependencies": {
    "@repo/shared": "workspace:*",
    "@tanstack/react-query": "^5.17.0",
    "axios": "^1.6.0",
    "class-variance-authority": "^0.7.0",
    "clsx": "^2.1.0",
    "lucide-react": "^0.303.0",
    "next": "^16.1.6",
    "react": "^19.0.0",
    "react-dom": "^19.0.0",
    "react-hook-form": "^7.49.0",
    "js-cookie": "^3.0.5",
    "tailwind-merge": "^2.2.0",
    "tailwindcss-animate": "^1.0.7"
  },
  "devDependencies": {
    "@testing-library/jest-dom": "^6.4.0",
    "@testing-library/react": "^16.0.0",
    "@testing-library/user-event": "^14.5.0",
    "@vitejs/plugin-react": "^4.3.0",
    "@types/node": "^20.0.0",
    "@types/js-cookie": "^3.0.6",
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

func webNextConfig() string {
	return `import type { NextConfig } from "next";
import { existsSync, readFileSync } from "node:fs";
import { resolve } from "node:path";

// Hoist the monorepo's root .env into process.env. Next.js auto-loads
// .env only from the package's own directory, so without this the THEME
// and SOCIAL_AUTH_ENABLED values set at the root are invisible to the
// web app. Shell env wins — we only fill in unset keys.
const rootEnv = resolve(process.cwd(), "..", "..", ".env");
if (existsSync(rootEnv)) {
  for (const line of readFileSync(rootEnv, "utf8").split(/\r?\n/)) {
    const m = line.match(/^([A-Z_][A-Z0-9_]*)=(.*)$/i);
    if (!m) continue;
    if (process.env[m[1]] === undefined) process.env[m[1]] = m[2].trim();
  }
}

const nextConfig: NextConfig = {
  output: "standalone",
  reactStrictMode: true,
  // packages/shared ships TypeScript source rather than a built bundle,
  // so Next needs to run it through SWC. Otherwise imports of
  // @repo/shared/types fail with "Cannot find module" at build time.
  transpilePackages: ["@repo/shared"],
  // Mirror THEME + SOCIAL_AUTH_ENABLED from .env into the NEXT_PUBLIC_*
  // namespace so the active theme is visible to server components and
  // the client bundle. Defaults keep new apps booting without env edits.
  env: {
    NEXT_PUBLIC_THEME: process.env.THEME || "atlas",
    NEXT_PUBLIC_SOCIAL_AUTH_ENABLED: process.env.SOCIAL_AUTH_ENABLED || "true",
  },
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
        // v3.28.1: --font-display + --font-mono are set per theme in the
        // root layout (Inter for atlas, Geist for aurora, Onest for pulse).
        sans: ["var(--font-display)", "system-ui", "sans-serif"],
        mono: ["var(--font-mono)", "ui-monospace", "monospace"],
        serif: ["var(--font-serif)", "Georgia", "serif"],
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

/* v3.28.1 — theme-aware CSS variables. Web mirrors admin's variable
 * system so a single THEME=<name> in .env paints both surfaces. */

/* atlas (default) */
:root,
[data-theme="atlas"] {
  --bg-primary: #ffffff;
  --bg-secondary: #f8fafc;
  --bg-tertiary: #f1f5f9;
  --bg-elevated: #ffffff;
  --bg-hover: #f1f5f9;
  --border: #e2e8f0;
  --text-primary: #0f172a;
  --text-secondary: #475569;
  --text-muted: #94a3b8;
  --accent: #2563eb;
  --accent-hover: #1d4ed8;
  --success: #10b981;
  --danger: #ef4444;
  --warning: #f59e0b;
  --info: #0ea5e9;
}

[data-theme="aurora"] {
  --bg-primary: #fafaf9;
  --bg-secondary: #ffffff;
  --bg-tertiary: #f5f5f4;
  --bg-elevated: #ffffff;
  --bg-hover: #f5f5f4;
  --border: #e7e5e4;
  --text-primary: #1c1917;
  --text-secondary: #57534e;
  --text-muted: #a8a29e;
  --accent: #7c3aed;
  --accent-hover: #6d28d9;
  --success: #10b981;
  --danger: #ef4444;
  --warning: #f59e0b;
  --info: #0ea5e9;
}

[data-theme="pulse"] {
  --bg-primary: #fafaf9;
  --bg-secondary: #ffffff;
  --bg-tertiary: #f5f5f4;
  --bg-elevated: #ffffff;
  --bg-hover: #f5f5f4;
  --border: #e5e5e5;
  --text-primary: #0f0f0f;
  --text-secondary: #525252;
  --text-muted: #a3a3a3;
  --accent: #0f0f0f;
  --accent-hover: #1f1f1f;
  --success: #16a34a;
  --danger: #dc2626;
  --warning: #fbbf24;
  --info: #0284c7;
}

[data-theme="midnight"] {
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
  font-family: var(--font-display), system-ui, sans-serif;
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
	// v3.28.1: per-theme font loading + data-theme attribute. Mirrors the
	// admin layout — see the comment there for the trade-off (build-time
	// font import means switching themes without re-scaffolding loses the
	// matching family; fine for now since runtime theme switching is rare).
	var fontImport, fontVars string
	switch opts.Theme {
	case "aurora":
		fontImport = `import { Geist, Geist_Mono } from "next/font/google";

const geist = Geist({ subsets: ["latin"], variable: "--font-display", weight: ["400", "500", "600", "700"] });
const geistMono = Geist_Mono({ subsets: ["latin"], variable: "--font-mono", weight: ["400", "500", "600"] });`
		fontVars = "${geist.variable} ${geistMono.variable}"
	case "pulse":
		fontImport = `import { Onest, DM_Serif_Display, JetBrains_Mono } from "next/font/google";

const onest = Onest({ subsets: ["latin"], variable: "--font-display", weight: ["400", "500", "600", "700"] });
const dmSerif = DM_Serif_Display({ subsets: ["latin"], variable: "--font-serif", weight: ["400"] });
const jetbrainsMono = JetBrains_Mono({ subsets: ["latin"], variable: "--font-mono", weight: ["400", "500", "600"] });`
		fontVars = "${onest.variable} ${dmSerif.variable} ${jetbrainsMono.variable}"
	default: // atlas
		fontImport = `import { Inter, JetBrains_Mono } from "next/font/google";

const inter = Inter({ subsets: ["latin"], variable: "--font-display", weight: ["400", "500", "600", "700"] });
const jetbrainsMono = JetBrains_Mono({ subsets: ["latin"], variable: "--font-mono", weight: ["400", "500", "600"] });`
		fontVars = "${inter.variable} ${jetbrainsMono.variable}"
	}

	return fmt.Sprintf(`import type { Metadata } from "next";
%s
import { Providers } from "@/components/providers";
import { Navbar } from "@/components/navbar";
import { Footer } from "@/components/footer";
import "./globals.css";

export const metadata: Metadata = {
  title: "%s — Go + React. Built with Grit.",
  description: "A full-stack framework that combines Go backend with Next.js frontend. Build fast, ship faster.",
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  const dataTheme = process.env.NEXT_PUBLIC_THEME || "atlas";

  return (
    <html lang="en" data-theme={dataTheme} suppressHydrationWarning>
      <body className={`+"`%s "+`font-sans antialiased`+"`"+`}>
        <Providers>
          <Navbar />
          <main className="min-h-screen">{children}</main>
          <Footer />
        </Providers>
      </body>
    </html>
  );
}`, fontImport, opts.ProjectName, fontVars)
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
              {blogs.map((blog) => (
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

// viteAPIClient is the Vite/TanStack-flavoured API client. It uses
// import.meta.env.VITE_API_URL (the Vite convention) instead of
// process.env.NEXT_PUBLIC_API_URL (which is a Next.js-only prefix and
// undefined in a Vite project). Default is '' so the embedded single-binary
// deploy serves SPA + API from the same origin without configuration.
func viteAPIClient() string {
	return `import axios from "axios";

const API_URL = (import.meta.env.VITE_API_URL as string | undefined) ?? "";

export const api = axios.create({
  baseURL: API_URL,
  headers: {
    "Content-Type": "application/json",
  },
});

// Auto-attach Idempotency-Key on unsafe methods so any mutation gets
// safe-retry semantics for free.
api.interceptors.request.use((config) => {
  const method = (config.method || "get").toUpperCase();
  const unsafe = method === "POST" || method === "PUT" || method === "PATCH" || method === "DELETE";
  if (unsafe && config.headers && !config.headers["Idempotency-Key"]) {
    config.headers["Idempotency-Key"] = crypto.randomUUID();
  }
  return config;
});
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
  // The browser attaches the HttpOnly grit_access / grit_refresh cookies
  // set by /api/auth/login automatically. Without this, axios skips them
  // on cross-origin requests in dev (api on :8080, web on :3000) and the
  // server treats every request as anonymous.
  withCredentials: true,
});

// Echo the grit_csrf cookie into X-CSRF-Token on every state-changing
// request. The cookie is intentionally not HttpOnly — it's the
// double-submit token, paired with the cookie the AutoCSRF middleware
// enforces. Safe-method requests don't need it; the middleware skips
// them and issues / refreshes the cookie as a side effect.
api.interceptors.request.use((config) => {
  if (typeof document !== "undefined") {
    const m = document.cookie.match(/(?:^|; )grit_csrf=([^;]+)/);
    if (m && config.headers) {
      config.headers["X-CSRF-Token"] = decodeURIComponent(m[1]);
    }
  }

  // Auto-attach Idempotency-Key on unsafe methods so any mutation gets
  // safe-retry semantics for free.
  const method = (config.method || "get").toUpperCase();
  const unsafe = method === "POST" || method === "PUT" || method === "PATCH" || method === "DELETE";
  if (unsafe && config.headers && !config.headers["Idempotency-Key"]) {
    config.headers["Idempotency-Key"] = crypto.randomUUID();
  }
  return config;
});
`
}

func webUseBlogsHook() string {
	return `"use client";

import { useQuery } from "@tanstack/react-query";
import type { Blog, PaginatedResponse } from "@repo/shared/types";
import { api } from "@/lib/api";

// Types live in packages/shared — never inline them in a hook. Reasons:
//   - the same Blog shape is consumed by web, admin, and the Go API (via
//     grit sync's generated counterpart). One source of truth.
//   - editing the model in Go and running 'grit sync' updates ALL
//     consumers in one shot; inline duplicates silently drift.
type BlogMeta = PaginatedResponse<Blog>["meta"];

export function usePublicBlogs(page = 1, pageSize = 20) {
  return useQuery({
    queryKey: ["public-blogs", page, pageSize],
    queryFn: async (): Promise<{ blogs: Blog[]; meta: BlogMeta | undefined }> => {
      const { data } = await api.get<PaginatedResponse<Blog>>(
        ` + "`" + `/api/blogs?page=${page}&page_size=${pageSize}` + "`" + `
      );
      return {
        blogs: data.data ?? [],
        meta: data.meta,
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

// Vite/TanStack blog routes import { useBlogs, useBlog } from this hook.
// Aliases so a single file works for both Next.js and Vite scaffolds.
export function useBlogs(page = 1, pageSize = 20) {
  const result = usePublicBlogs(page, pageSize);
  return { ...result, data: result.data?.blogs };
}

export const useBlog = usePublicBlog;
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

func webComponentsPage(_ Options) string {
	return `"use client";

import { useState } from "react";
import Link from "next/link";
import { Copy, Check, Package } from "lucide-react";
import { useQuery } from "@tanstack/react-query";
import { api } from "@/lib/api";

interface RegistryItem {
  name: string;
  type: string;
  title: string;
  description: string;
  category: string;
  registryUrl?: string;
}

interface Registry {
  name: string;
  items: RegistryItem[];
}

function CopyButton({ text }: { text: string }) {
  const [copied, setCopied] = useState(false);
  const copy = async () => {
    await navigator.clipboard.writeText(text);
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  };
  return (
    <button
      onClick={copy}
      className="flex-shrink-0 rounded-md p-1.5 text-text-muted hover:text-foreground hover:bg-bg-hover transition-colors"
      title="Copy"
    >
      {copied ? <Check className="h-4 w-4 text-success" /> : <Copy className="h-4 w-4" />}
    </button>
  );
}

const CATEGORIES = ["All", "marketing", "ecommerce", "layout"];

export default function ComponentsPage() {
  const [activeCategory, setActiveCategory] = useState("All");

  const { data: registry, isLoading } = useQuery<Registry>({
    queryKey: ["grit-ui-registry"],
    queryFn: async () => {
      const { data } = await api.get("/r.json");
      return data;
    },
  });

  const items = (registry?.items || []).filter(
    (item) => activeCategory === "All" || item.category === activeCategory
  );

  const apiUrl = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";

  return (
    <div className="mx-auto max-w-6xl px-6 py-16">
      {/* Hero */}
      <div className="mb-12 max-w-2xl">
        <div className="mb-3 inline-flex items-center gap-2 rounded-full border border-accent/20 bg-accent/10 px-3 py-1 text-xs font-medium text-accent">
          <Package className="h-3.5 w-3.5" />
          Grit UI Registry
        </div>
        <h1 className="text-4xl font-bold tracking-tight">
          Beautiful components for{" "}
          <span className="bg-gradient-to-r from-accent to-purple-400 bg-clip-text text-transparent">
            Grit apps
          </span>
        </h1>
        <p className="mt-3 text-lg text-text-secondary">
          Production-ready UI components. Install with a single command via the
          shadcn CLI.
        </p>
      </div>

      {/* Category pills */}
      <div className="mb-8 flex flex-wrap gap-2">
        {CATEGORIES.map((cat) => (
          <button
            key={cat}
            onClick={() => setActiveCategory(cat)}
            className={"rounded-full border px-4 py-1.5 text-sm font-medium capitalize transition-colors " +
              (activeCategory === cat
                ? "border-accent bg-accent text-white"
                : "border-border bg-bg-elevated text-text-secondary hover:border-accent/40 hover:text-foreground")}
          >
            {cat}
          </button>
        ))}
      </div>

      {/* Grid */}
      {isLoading ? (
        <div className="grid gap-5 sm:grid-cols-2 lg:grid-cols-3">
          {Array.from({ length: 6 }).map((_, i) => (
            <div
              key={i}
              className="rounded-xl border border-border bg-bg-elevated p-5 animate-pulse"
            >
              <div className="mb-3 h-4 w-20 rounded bg-bg-hover" />
              <div className="mb-2 h-5 w-3/4 rounded bg-bg-hover" />
              <div className="h-3 w-full rounded bg-bg-hover" />
              <div className="mt-4 h-8 w-full rounded bg-bg-hover" />
            </div>
          ))}
        </div>
      ) : items.length > 0 ? (
        <div className="grid gap-5 sm:grid-cols-2 lg:grid-cols-3">
          {items.map((item) => {
            const installCmd = ` + "`" + `npx shadcn@latest add ${apiUrl}/r/${item.name}.json` + "`" + `;
            return (
              <div
                key={item.name}
                className="group rounded-xl border border-border bg-bg-elevated p-5 hover:border-accent/40 hover:shadow-lg hover:shadow-accent/5 transition-all duration-300 flex flex-col"
              >
                <div className="mb-3 flex items-center gap-2">
                  <span className="rounded-md bg-bg-hover px-2 py-0.5 text-xs font-mono text-text-muted capitalize">
                    {item.category}
                  </span>
                </div>
                <Link
                  href={` + "`" + `/components/${item.name}` + "`" + `}
                  className="flex-1"
                >
                  <h3 className="font-semibold text-foreground group-hover:text-accent transition-colors">
                    {item.title || item.name}
                  </h3>
                  <p className="mt-1.5 text-sm text-text-secondary line-clamp-2">
                    {item.description}
                  </p>
                </Link>
                <div className="mt-4 flex items-center gap-2 rounded-lg border border-border bg-bg-primary px-3 py-2">
                  <code className="flex-1 truncate text-xs font-mono text-text-secondary">
                    {installCmd}
                  </code>
                  <CopyButton text={installCmd} />
                </div>
              </div>
            );
          })}
        </div>
      ) : (
        <div className="flex flex-col items-center justify-center py-20 text-center">
          <div className="mb-4 flex h-16 w-16 items-center justify-center rounded-2xl bg-bg-elevated border border-border">
            <Package className="h-7 w-7 text-text-muted" />
          </div>
          <h3 className="text-lg font-semibold">No components yet</h3>
          <p className="mt-1 text-sm text-text-muted">
            Components will appear here once added via the admin panel.
          </p>
        </div>
      )}
    </div>
  );
}
`
}

func webComponentDetailPage(_ Options) string {
	return `"use client";

import { useParams } from "next/navigation";
import Link from "next/link";
import { ArrowLeft, Copy, Check, Terminal } from "lucide-react";
import { useQuery } from "@tanstack/react-query";
import { api } from "@/lib/api";
import { useState } from "react";

interface ComponentFile {
  path: string;
  content: string;
  type: string;
  target?: string;
}

interface ComponentDetail {
  name: string;
  type: string;
  title: string;
  description: string;
  category: string;
  files: ComponentFile[];
  dependencies: string[];
  registryDependencies: string[];
}

function CopyButton({ text, label }: { text: string; label?: string }) {
  const [copied, setCopied] = useState(false);
  const copy = async () => {
    await navigator.clipboard.writeText(text);
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  };
  return (
    <button
      onClick={copy}
      className="inline-flex items-center gap-1.5 rounded-lg border border-border bg-bg-elevated px-3 py-1.5 text-sm text-text-secondary hover:text-foreground hover:bg-bg-hover transition-colors"
    >
      {copied ? (
        <Check className="h-3.5 w-3.5 text-success" />
      ) : (
        <Copy className="h-3.5 w-3.5" />
      )}
      {label ?? (copied ? "Copied!" : "Copy")}
    </button>
  );
}

export default function ComponentDetailPage() {
  const params = useParams();
  const name = typeof params.name === "string" ? params.name : "";
  const apiUrl = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";

  const { data: component, isLoading, error } = useQuery<ComponentDetail>({
    queryKey: ["grit-ui-component", name],
    queryFn: async () => {
      const { data } = await api.get(` + "`" + `/r/${name}.json` + "`" + `);
      return data;
    },
    enabled: !!name,
  });

  const installCmd = ` + "`" + `npx shadcn@latest add ${apiUrl}/r/${name}.json` + "`" + `;

  if (isLoading) {
    return (
      <div className="mx-auto max-w-4xl px-6 py-16 animate-pulse">
        <div className="h-4 w-24 rounded bg-bg-hover mb-8" />
        <div className="h-8 w-1/2 rounded bg-bg-hover mb-3" />
        <div className="h-4 w-2/3 rounded bg-bg-hover mb-10" />
        <div className="h-12 rounded-xl bg-bg-hover mb-10" />
        <div className="h-96 rounded-xl bg-bg-hover" />
      </div>
    );
  }

  if (error || !component) {
    return (
      <div className="mx-auto max-w-4xl px-6 py-16 text-center">
        <p className="text-6xl font-bold text-text-muted mb-4">404</p>
        <h1 className="text-xl font-semibold mb-2">Component not found</h1>
        <p className="text-sm text-text-muted mb-6">
          The component &quot;{name}&quot; doesn&apos;t exist in the registry.
        </p>
        <Link
          href="/components"
          className="inline-flex items-center gap-2 rounded-lg bg-accent px-4 py-2 text-sm font-medium text-white hover:bg-accent-hover transition-colors"
        >
          <ArrowLeft className="h-4 w-4" />
          Back to Components
        </Link>
      </div>
    );
  }

  const mainFile = component.files?.[0];

  return (
    <div className="mx-auto max-w-4xl px-6 py-16">
      {/* Back */}
      <Link
        href="/components"
        className="inline-flex items-center gap-1.5 text-sm text-text-secondary hover:text-foreground transition-colors mb-8"
      >
        <ArrowLeft className="h-4 w-4" />
        All Components
      </Link>

      {/* Header */}
      <div className="mb-8">
        <div className="mb-2 flex items-center gap-3">
          <span className="rounded-md bg-bg-elevated border border-border px-2.5 py-0.5 text-xs font-mono text-text-muted capitalize">
            {component.category}
          </span>
        </div>
        <h1 className="text-3xl font-bold tracking-tight">
          {component.title || component.name}
        </h1>
        <p className="mt-2 text-text-secondary">{component.description}</p>
      </div>

      {/* Install command */}
      <div className="mb-8 rounded-xl border border-border bg-bg-elevated p-4">
        <div className="mb-2 flex items-center gap-2 text-xs font-medium text-text-muted uppercase tracking-wide">
          <Terminal className="h-3.5 w-3.5" />
          Install
        </div>
        <div className="flex items-center gap-3">
          <code className="flex-1 text-sm font-mono text-text-secondary truncate">
            {installCmd}
          </code>
          <CopyButton text={installCmd} label="Copy" />
        </div>
      </div>

      {/* Dependencies */}
      {component.dependencies && component.dependencies.length > 0 && (
        <div className="mb-8">
          <h2 className="mb-2 text-sm font-medium text-text-secondary uppercase tracking-wide">
            Dependencies
          </h2>
          <div className="flex flex-wrap gap-2">
            {component.dependencies.map((dep) => (
              <span
                key={dep}
                className="rounded-md border border-border bg-bg-elevated px-2.5 py-1 text-xs font-mono text-text-secondary"
              >
                {dep}
              </span>
            ))}
          </div>
        </div>
      )}

      {/* Source code */}
      {mainFile && (
        <div className="rounded-xl border border-border bg-bg-elevated overflow-hidden">
          <div className="flex items-center justify-between border-b border-border px-4 py-3">
            <span className="text-xs font-mono text-text-muted">
              {mainFile.target || mainFile.path}
            </span>
            <CopyButton text={mainFile.content} label="Copy code" />
          </div>
          <div className="overflow-x-auto">
            <pre className="p-5 text-sm font-mono text-text-secondary leading-relaxed whitespace-pre">
              <code>{mainFile.content}</code>
            </pre>
          </div>
        </div>
      )}
    </div>
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

// ---------------------------------------------------------------------------
// Web Auth Pages
// ---------------------------------------------------------------------------

func webAuthLogin(_ Options) string {
	return `"use client";

import { useState } from "react";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { api } from "@/lib/api";
import { isSocialAuthEnabled } from "@repo/shared/themes";
import { brand } from "@repo/shared/brand";

const inputClass = "w-full rounded-lg border border-border bg-bg-elevated px-4 py-3 text-foreground placeholder:text-text-muted focus:border-accent focus:outline-none focus:ring-1 focus:ring-accent transition-colors";

export default function LoginPage() {
  const router = useRouter();
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [showPassword, setShowPassword] = useState(false);
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);

  const onSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError("");
    setLoading(true);
    try {
      // The API responds with Set-Cookie: grit_access + grit_refresh
      // (HttpOnly). The axios client has withCredentials: true, so the
      // browser stores them automatically. We don't touch tokens in JS.
      await api.post("/api/auth/login", { email, password });
      router.push("/");
    } catch (err: unknown) {
      const msg = (err as { response?: { data?: { error?: { message?: string } } } })?.response?.data?.error?.message;
      setError(msg || "Invalid email or password");
    } finally {
      setLoading(false);
    }
  };

  const API_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";

  return (
    <div className="flex min-h-screen items-center justify-center bg-background px-4">
      <div className="w-full max-w-md space-y-8">
        <div className="text-center">
          <h1 className="text-3xl font-bold text-foreground">Welcome back to {brand.name}</h1>
          <p className="mt-2 text-text-secondary">Sign in to your account</p>
        </div>

        <form onSubmit={onSubmit} className="space-y-5">
          {error && (
            <div className="rounded-lg bg-danger/10 border border-danger/20 px-4 py-3 text-sm text-danger">
              {error}
            </div>
          )}

          <div className="space-y-2">
            <label htmlFor="email" className="block text-sm font-medium text-text-secondary">
              Email
            </label>
            <input
              id="email"
              type="email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              className={inputClass}
              placeholder="you@example.com"
              required
              autoFocus
            />
          </div>

          <div className="space-y-2">
            <label htmlFor="password" className="block text-sm font-medium text-text-secondary">
              Password
            </label>
            <div className="relative">
              <input
                id="password"
                type={showPassword ? "text" : "password"}
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                className={inputClass + " pr-12"}
                placeholder="Enter your password"
                required
              />
              <button
                type="button"
                onClick={() => setShowPassword(!showPassword)}
                className="absolute right-3 top-1/2 -translate-y-1/2 text-text-muted hover:text-text-secondary transition-colors text-sm"
              >
                {showPassword ? "Hide" : "Show"}
              </button>
            </div>
          </div>

          <div className="flex items-center justify-end">
            <Link href="/forgot-password" className="text-sm text-accent hover:text-accent-hover transition-colors">
              Forgot password?
            </Link>
          </div>

          <button
            type="submit"
            disabled={loading}
            className="w-full rounded-lg bg-accent py-3 font-medium text-white hover:bg-accent-hover disabled:opacity-50 transition-colors"
          >
            {loading ? "Signing in..." : "Sign In"}
          </button>
        </form>

        {/* Social Login — hidden when SOCIAL_AUTH_ENABLED=false in .env */}
        {isSocialAuthEnabled() && (<>
        <div className="relative my-6">
          <div className="absolute inset-0 flex items-center">
            <div className="w-full border-t border-border" />
          </div>
          <div className="relative flex justify-center text-xs">
            <span className="bg-background px-3 text-text-muted">or continue with</span>
          </div>
        </div>

        <div className="flex gap-3">
          <a
            href={` + "`" + `${API_URL}/api/auth/oauth/google` + "`" + `}
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
            href={` + "`" + `${API_URL}/api/auth/oauth/github` + "`" + `}
            className="flex flex-1 items-center justify-center gap-2 rounded-lg border border-border bg-[#24292f] px-4 py-2.5 text-sm font-medium text-white hover:bg-[#2f363d] transition-colors"
          >
            <svg className="h-5 w-5" fill="currentColor" viewBox="0 0 24 24">
              <path fillRule="evenodd" d="M12 2C6.477 2 2 6.484 2 12.017c0 4.425 2.865 8.18 6.839 9.504.5.092.682-.217.682-.483 0-.237-.008-.868-.013-1.703-2.782.605-3.369-1.343-3.369-1.343-.454-1.158-1.11-1.466-1.11-1.466-.908-.62.069-.608.069-.608 1.003.07 1.531 1.032 1.531 1.032.892 1.53 2.341 1.088 2.91.832.092-.647.35-1.088.636-1.338-2.22-.253-4.555-1.113-4.555-4.951 0-1.093.39-1.988 1.029-2.688-.103-.253-.446-1.272.098-2.65 0 0 .84-.27 2.75 1.026A9.564 9.564 0 0 1 12 6.844a9.59 9.59 0 0 1 2.504.337c1.909-1.296 2.747-1.027 2.747-1.027.546 1.379.202 2.398.1 2.651.64.7 1.028 1.595 1.028 2.688 0 3.848-2.339 4.695-4.566 4.943.359.309.678.92.678 1.855 0 1.338-.012 2.419-.012 2.747 0 .268.18.58.688.482A10.02 10.02 0 0 0 22 12.017C22 6.484 17.522 2 12 2z" clipRule="evenodd" />
            </svg>
            GitHub
          </a>
        </div>
        </>)}

        <p className="text-center text-sm text-text-secondary">
          Don&apos;t have an account?{" "}
          <Link href="/register" className="text-accent hover:text-accent-hover font-medium transition-colors">
            Create one
          </Link>
        </p>
      </div>
    </div>
  );
}
`
}

func webAuthRegister(_ Options) string {
	return `"use client";

import { useState } from "react";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { api } from "@/lib/api";
import { isSocialAuthEnabled } from "@repo/shared/themes";

const inputClass = "w-full rounded-lg border border-border bg-bg-elevated px-4 py-3 text-foreground placeholder:text-text-muted focus:border-accent focus:outline-none focus:ring-1 focus:ring-accent transition-colors";

export default function RegisterPage() {
  const router = useRouter();
  const [firstName, setFirstName] = useState("");
  const [lastName, setLastName] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [confirmPassword, setConfirmPassword] = useState("");
  const [showPassword, setShowPassword] = useState(false);
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);

  const onSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError("");

    if (password !== confirmPassword) {
      setError("Passwords do not match");
      return;
    }

    if (password.length < 8) {
      setError("Password must be at least 8 characters");
      return;
    }

    setLoading(true);
    try {
      // The API sets HttpOnly grit_access + grit_refresh cookies on the
      // response. The browser stores them automatically (axios withCredentials
      // is on). No tokens touched in JS.
      await api.post("/api/auth/register", {
        first_name: firstName,
        last_name: lastName,
        email,
        password,
      });
      router.push("/");
    } catch (err: unknown) {
      const msg = (err as { response?: { data?: { error?: { message?: string } } } })?.response?.data?.error?.message;
      setError(msg || "Registration failed. Please try again.");
    } finally {
      setLoading(false);
    }
  };

  const API_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";

  return (
    <div className="flex min-h-screen items-center justify-center bg-background px-4 py-12">
      <div className="w-full max-w-md space-y-8">
        <div className="text-center">
          <h1 className="text-3xl font-bold text-foreground">Create account</h1>
          <p className="mt-2 text-text-secondary">Sign up to get started</p>
        </div>

        <form onSubmit={onSubmit} className="space-y-5">
          {error && (
            <div className="rounded-lg bg-danger/10 border border-danger/20 px-4 py-3 text-sm text-danger">
              {error}
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
                value={firstName}
                onChange={(e) => setFirstName(e.target.value)}
                className={inputClass}
                placeholder="John"
                required
                autoFocus
              />
            </div>
            <div className="space-y-2">
              <label htmlFor="lastName" className="block text-sm font-medium text-text-secondary">
                Last name
              </label>
              <input
                id="lastName"
                type="text"
                value={lastName}
                onChange={(e) => setLastName(e.target.value)}
                className={inputClass}
                placeholder="Doe"
                required
              />
            </div>
          </div>

          <div className="space-y-2">
            <label htmlFor="email" className="block text-sm font-medium text-text-secondary">
              Email
            </label>
            <input
              id="email"
              type="email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              className={inputClass}
              placeholder="you@example.com"
              required
            />
          </div>

          <div className="space-y-2">
            <label htmlFor="password" className="block text-sm font-medium text-text-secondary">
              Password
            </label>
            <div className="relative">
              <input
                id="password"
                type={showPassword ? "text" : "password"}
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                className={inputClass + " pr-12"}
                placeholder="Min. 8 characters"
                required
              />
              <button
                type="button"
                onClick={() => setShowPassword(!showPassword)}
                className="absolute right-3 top-1/2 -translate-y-1/2 text-text-muted hover:text-text-secondary transition-colors text-sm"
              >
                {showPassword ? "Hide" : "Show"}
              </button>
            </div>
          </div>

          <div className="space-y-2">
            <label htmlFor="confirmPassword" className="block text-sm font-medium text-text-secondary">
              Confirm password
            </label>
            <input
              id="confirmPassword"
              type="password"
              value={confirmPassword}
              onChange={(e) => setConfirmPassword(e.target.value)}
              className={inputClass}
              placeholder="Repeat your password"
              required
            />
          </div>

          <button
            type="submit"
            disabled={loading}
            className="w-full rounded-lg bg-accent py-3 font-medium text-white hover:bg-accent-hover disabled:opacity-50 transition-colors"
          >
            {loading ? "Creating account..." : "Create Account"}
          </button>
        </form>

        {/* Social Login — hidden when SOCIAL_AUTH_ENABLED=false in .env */}
        {isSocialAuthEnabled() && (<>
        <div className="relative my-6">
          <div className="absolute inset-0 flex items-center">
            <div className="w-full border-t border-border" />
          </div>
          <div className="relative flex justify-center text-xs">
            <span className="bg-background px-3 text-text-muted">or continue with</span>
          </div>
        </div>

        <div className="flex gap-3">
          <a
            href={` + "`" + `${API_URL}/api/auth/oauth/google` + "`" + `}
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
            href={` + "`" + `${API_URL}/api/auth/oauth/github` + "`" + `}
            className="flex flex-1 items-center justify-center gap-2 rounded-lg border border-border bg-[#24292f] px-4 py-2.5 text-sm font-medium text-white hover:bg-[#2f363d] transition-colors"
          >
            <svg className="h-5 w-5" fill="currentColor" viewBox="0 0 24 24">
              <path fillRule="evenodd" d="M12 2C6.477 2 2 6.484 2 12.017c0 4.425 2.865 8.18 6.839 9.504.5.092.682-.217.682-.483 0-.237-.008-.868-.013-1.703-2.782.605-3.369-1.343-3.369-1.343-.454-1.158-1.11-1.466-1.11-1.466-.908-.62.069-.608.069-.608 1.003.07 1.531 1.032 1.531 1.032.892 1.53 2.341 1.088 2.91.832.092-.647.35-1.088.636-1.338-2.22-.253-4.555-1.113-4.555-4.951 0-1.093.39-1.988 1.029-2.688-.103-.253-.446-1.272.098-2.65 0 0 .84-.27 2.75 1.026A9.564 9.564 0 0 1 12 6.844a9.59 9.59 0 0 1 2.504.337c1.909-1.296 2.747-1.027 2.747-1.027.546 1.379.202 2.398.1 2.651.64.7 1.028 1.595 1.028 2.688 0 3.848-2.339 4.695-4.566 4.943.359.309.678.92.678 1.855 0 1.338-.012 2.419-.012 2.747 0 .268.18.58.688.482A10.02 10.02 0 0 0 22 12.017C22 6.484 17.522 2 12 2z" clipRule="evenodd" />
            </svg>
            GitHub
          </a>
        </div>
        </>)}

        <p className="text-center text-sm text-text-secondary">
          Already have an account?{" "}
          <Link href="/login" className="text-accent hover:text-accent-hover font-medium transition-colors">
            Sign in
          </Link>
        </p>
      </div>
    </div>
  );
}
`
}

func webAuthForgotPassword(_ Options) string {
	return `"use client";

import { useState } from "react";
import Link from "next/link";
import { api } from "@/lib/api";

const inputClass = "w-full rounded-lg border border-border bg-bg-elevated px-4 py-3 text-foreground placeholder:text-text-muted focus:border-accent focus:outline-none focus:ring-1 focus:ring-accent transition-colors";

export default function ForgotPasswordPage() {
  const [email, setEmail] = useState("");
  const [sent, setSent] = useState(false);
  const [loading, setLoading] = useState(false);

  const onSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    try {
      await api.post("/api/auth/forgot-password", { email });
      setSent(true);
    } catch {
      setSent(true); // Always show success for security
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="flex min-h-screen items-center justify-center bg-background px-4">
      <div className="w-full max-w-md space-y-8">
        <div className="text-center">
          <h1 className="text-3xl font-bold text-foreground">Reset password</h1>
          <p className="mt-2 text-text-secondary">
            Enter your email and we&apos;ll send you a reset link
          </p>
        </div>

        {sent ? (
          <div className="space-y-6">
            <div className="rounded-lg bg-success/10 border border-success/20 px-4 py-4 text-center">
              <p className="text-sm text-success font-medium">Check your email</p>
              <p className="mt-1 text-xs text-text-secondary">
                If an account exists with that email, we&apos;ve sent a password reset link.
              </p>
            </div>
            <Link
              href="/login"
              className="block w-full rounded-lg border border-border py-3 text-center font-medium text-foreground hover:bg-bg-hover transition-colors"
            >
              Back to sign in
            </Link>
          </div>
        ) : (
          <form onSubmit={onSubmit} className="space-y-5">
            <div className="space-y-2">
              <label htmlFor="email" className="block text-sm font-medium text-text-secondary">
                Email
              </label>
              <input
                id="email"
                type="email"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                className={inputClass}
                placeholder="you@example.com"
                required
                autoFocus
              />
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
          Remember your password?{" "}
          <Link href="/login" className="text-accent hover:text-accent-hover font-medium transition-colors">
            Sign in
          </Link>
        </p>
      </div>
    </div>
  );
}
`
}

func webAuthCallback() string {
	return `"use client";

import { Suspense, useEffect, useRef } from "react";
import { useRouter, useSearchParams } from "next/navigation";
import Cookies from "js-cookie";
import { api } from "@/lib/api";

// useSearchParams forces a client bailout during prerender unless it's
// wrapped in a Suspense boundary. The inner component does the work;
// the page export just provides the boundary.
export default function AuthCallbackPage() {
  return (
    <Suspense fallback={<CallbackSpinner />}>
      <CallbackInner />
    </Suspense>
  );
}

function CallbackInner() {
  const router = useRouter();
  const searchParams = useSearchParams();
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

      api
        .get("/api/auth/me", {
          headers: { Authorization: "Bearer " + accessToken },
        })
        .then(() => {
          router.push("/");
        })
        .catch(() => {
          router.push("/");
        });
    } else {
      router.push("/login?error=Authentication+failed");
    }
  }, [searchParams, router]);

  return <CallbackSpinner />;
}

function CallbackSpinner() {
  return (
    <div className="flex min-h-screen items-center justify-center bg-background">
      <div className="text-center">
        <div className="inline-flex h-10 w-10 animate-spin items-center justify-center rounded-full border-2 border-accent border-t-transparent" />
        <p className="mt-4 text-sm text-text-secondary">Signing you in...</p>
      </div>
    </div>
  );
}
`
}

func webUseAuth() string {
	return `import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useRouter } from "next/navigation";
import type {
  User,
  LoginRequest,
  RegisterRequest,
  AuthResponse,
  ApiResponse,
} from "@repo/shared/types";
import { api } from "@/lib/api";

// Token storage policy (Grit 3.26+):
//   - The API issues HttpOnly cookies (grit_access + grit_refresh) on
//     login/register/refresh and clears them on logout. The browser
//     handles them automatically — JS never reads or writes the access
//     token, so XSS can't exfiltrate it.
//   - The axios client uses withCredentials: true so the browser
//     attaches the cookies on every request.
//   - There is no Bearer-header path here. Mobile/desktop bearer
//     clients live in their own apps; the web app is cookie-only.
//
// Types live in packages/shared/types/user.ts — they're consumed by
// web, admin, and the Go API (via grit sync). Never inline them.

export function useMe() {
  return useQuery<User | null>({
    queryKey: ["me"],
    queryFn: async () => {
      try {
        const { data } = await api.get<ApiResponse<User>>("/api/auth/me");
        return data.data;
      } catch (err: unknown) {
        // 401 is the canonical "not logged in" — return null instead
        // of throwing so guarded pages can read user === null cleanly.
        const e = err as { response?: { status?: number } };
        if (e.response?.status === 401) return null;
        throw err;
      }
    },
    retry: false,
    staleTime: 10 * 60 * 1000,
  });
}

export function useAuth() {
  const { data: user, isLoading, isError } = useMe();
  return {
    user: user ?? null,
    isAuthenticated: !!user,
    isLoading,
    isError,
  };
}

export function useLogin() {
  const router = useRouter();
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (credentials: LoginRequest) => {
      // POST is enough — the API sets HttpOnly cookies on the response.
      // The 'tokens' field is still in the JSON body so native bearer
      // clients work too, but the browser ignores it.
      const { data } = await api.post<ApiResponse<AuthResponse>>(
        "/api/auth/login",
        credentials
      );
      return data;
    },
    onSuccess: (data) => {
      queryClient.setQueryData(["me"], data.data.user);
      router.push("/");
    },
  });
}

export function useRegister() {
  const router = useRouter();
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (payload: RegisterRequest) => {
      const { data } = await api.post<ApiResponse<AuthResponse>>(
        "/api/auth/register",
        payload
      );
      return data;
    },
    onSuccess: (data) => {
      queryClient.setQueryData(["me"], data.data.user);
      router.push("/");
    },
  });
}

export function useLogout() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async () => {
      try {
        await api.post("/api/auth/logout");
      } catch {
        // The cookies are already cleared by the API's Set-Cookie even
        // on a 4xx; just make sure local state is wiped either way.
      }
    },
    onSettled: () => {
      queryClient.clear();
      if (typeof window !== "undefined") {
        window.location.href = "/login";
      }
    },
  });
}
`
}

func webAuthProvider() string {
	return `"use client";

import { createContext, useContext } from "react";
import type { User } from "@repo/shared/types";
import { useMe } from "@/hooks/use-auth";

// User shape comes from packages/shared/types/user.ts — same shape the
// Go API serialises, same shape the admin imports. Never inline it.

interface AuthContextValue {
  user: User | null;
  isAuthenticated: boolean;
  isLoading: boolean;
}

const AuthContext = createContext<AuthContextValue>({
  user: null,
  isAuthenticated: false,
  isLoading: true,
});

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const { data: user, isLoading } = useMe();

  return (
    <AuthContext.Provider
      value={{
        user: user ?? null,
        isAuthenticated: !!user,
        isLoading,
      }}
    >
      {children}
    </AuthContext.Provider>
  );
}

export function useAuthContext() {
  return useContext(AuthContext);
}
`
}

// webPublicFormPage — public-facing page at /forms/[token] that
// renders a generic submission form for any FormShare-exposed resource.
//
// Designed to be minimal but functional:
//   - Fetches /api/public/forms/<token> to confirm the link works +
//     learn the resource_name and whether a password is required.
//   - If password required, shows a gate input. Submitting hits the
//     submit endpoint with the password+empty fields to probe; on a
//     successful probe we move to the form. (Phase 2 hardens this by
//     adding a dedicated /check-password endpoint that issues a
//     short-lived cookie — for v1 the password is sent with each
//     submit so the page itself stays stateless.)
//   - The form itself is a one-size-fits-all key/value editor where
//     visitors fill labelled inputs. The shape is hard-coded to a
//     name+email+message pattern that covers the dominant case (lead
//     forms, contact forms, applications). Resources with different
//     shapes are best exposed via `grit expose form` once Phase 3
//     lands — that command generates a page tailored to the resource's
//     actual fields.
func webPublicFormPage() string {
	return `"use client";

import { useEffect, useState, use } from "react";
import axios from "axios";

const API_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";

interface PageProps {
  params: Promise<{ token: string }>;
}

interface ShareInfo {
  resource_name: string;
  has_password: boolean;
  label: string;
}

export default function PublicFormPage({ params }: PageProps) {
  const { token } = use(params);
  const [info, setInfo] = useState<ShareInfo | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    axios.get(API_URL + "/api/public/forms/" + token)
      .then((res) => setInfo(res.data.data))
      .catch((err) => {
        setError(err?.response?.data?.error?.message || "Link not found or disabled");
      })
      .finally(() => setLoading(false));
  }, [token]);

  if (loading) {
    return (
      <main className="flex min-h-screen items-center justify-center bg-slate-50 p-4">
        <div className="text-sm text-slate-500">Loading…</div>
      </main>
    );
  }

  if (error || !info) {
    return (
      <main className="flex min-h-screen items-center justify-center bg-slate-50 p-4">
        <div className="w-full max-w-md rounded-2xl border border-slate-200 bg-white p-8 text-center shadow-sm">
          <h1 className="text-xl font-semibold text-slate-900">Link unavailable</h1>
          <p className="mt-2 text-sm text-slate-500">{error ?? "Unknown error"}</p>
        </div>
      </main>
    );
  }

  return (
    <main className="flex min-h-screen items-center justify-center bg-slate-50 p-4">
      <div className="w-full max-w-md rounded-2xl border border-slate-200 bg-white p-8 shadow-sm">
        <h1 className="text-2xl font-semibold text-slate-900">
          {info.label || (info.resource_name + " submission")}
        </h1>
        <p className="mt-1 text-sm text-slate-500">
          Fill out the form below to submit a new {info.resource_name}.
        </p>

        <PublicForm token={token} info={info} />
      </div>
    </main>
  );
}

interface PublicFormProps {
  token: string;
  info: ShareInfo;
}

function PublicForm({ token, info }: PublicFormProps) {
  const [password, setPassword] = useState("");
  const [fields, setFields] = useState<Record<string, string>>({});
  const [submitting, setSubmitting] = useState(false);
  const [done, setDone] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const update = (key: string, value: string) => {
    setFields((prev) => ({ ...prev, [key]: value }));
  };

  const onSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);
    setSubmitting(true);
    try {
      await axios.post(API_URL + "/api/public/forms/" + token + "/submit", {
        _password: password,
        fields,
      });
      setDone(true);
    } catch (err) {
      const e = err as { response?: { data?: { error?: { message?: string } } } };
      setError(e?.response?.data?.error?.message || "Submission failed");
    } finally {
      setSubmitting(false);
    }
  };

  if (done) {
    return (
      <div className="mt-6 rounded-xl border border-emerald-200 bg-emerald-50 p-4 text-sm text-emerald-700">
        <p className="font-medium">Thank you</p>
        <p className="mt-1">Your submission was received.</p>
      </div>
    );
  }

  return (
    <form onSubmit={onSubmit} className="mt-6 space-y-4">
      {info.has_password && (
        <Field
          label="Password"
          type="password"
          value={password}
          onChange={setPassword}
          required
          hint="This form is password-protected — ask whoever shared the link."
        />
      )}

      {/* v1 generic shape: name + email + phone + message. Operators
          whose resources need different fields should use the
          ` + "`" + `grit expose form` + "`" + ` command (Phase 3) to scaffold a tailored
          public page. */}
      <Field label="Name"    value={fields.name    ?? ""} onChange={(v) => update("name", v)}    required />
      <Field label="Email"   value={fields.email   ?? ""} onChange={(v) => update("email", v)}   type="email" />
      <Field label="Phone"   value={fields.phone   ?? ""} onChange={(v) => update("phone", v)}   type="tel" />
      <Field label="Message" value={fields.message ?? ""} onChange={(v) => update("message", v)} rows={4} />

      {error && (
        <div className="rounded-lg border border-red-200 bg-red-50 px-3 py-2 text-sm text-red-700">
          {error}
        </div>
      )}

      <button
        type="submit"
        disabled={submitting}
        className="w-full rounded-lg bg-slate-900 px-4 py-2.5 text-sm font-semibold text-white hover:bg-slate-800 disabled:opacity-50"
      >
        {submitting ? "Sending…" : "Submit"}
      </button>
    </form>
  );
}

interface FieldProps {
  label: string;
  value: string;
  onChange: (v: string) => void;
  required?: boolean;
  type?: string;
  rows?: number;
  hint?: string;
}

function Field({ label, value, onChange, required, type, rows, hint }: FieldProps) {
  return (
    <div className="space-y-1.5">
      <label className="block text-sm font-medium text-slate-700">
        {label}
        {required && <span className="ml-1 text-red-500">*</span>}
      </label>
      {rows ? (
        <textarea
          value={value}
          onChange={(e) => onChange(e.target.value)}
          required={required}
          rows={rows}
          className="block w-full rounded-lg border border-slate-300 bg-white px-3 py-2 text-sm text-slate-900 placeholder-slate-400 focus:border-slate-900 focus:outline-none focus:ring-2 focus:ring-slate-900/10"
        />
      ) : (
        <input
          type={type ?? "text"}
          value={value}
          onChange={(e) => onChange(e.target.value)}
          required={required}
          className="block w-full rounded-lg border border-slate-300 bg-white px-3 py-2 text-sm text-slate-900 placeholder-slate-400 focus:border-slate-900 focus:outline-none focus:ring-2 focus:ring-slate-900/10"
        />
      )}
      {hint && <p className="text-xs text-slate-500">{hint}</p>}
    </div>
  );
}
`
}

