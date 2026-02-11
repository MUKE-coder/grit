package scaffold

import (
	"fmt"
	"path/filepath"
)

func writeWebFiles(root string, opts Options) error {
	webRoot := filepath.Join(root, "apps", "web")

	files := map[string]string{
		filepath.Join(webRoot, "package.json"):       webPackageJSON(opts),
		filepath.Join(webRoot, "next.config.ts"):     webNextConfig(),
		filepath.Join(webRoot, "tailwind.config.ts"): webTailwindConfig(),
		filepath.Join(webRoot, "postcss.config.js"):  webPostCSSConfig(),
		filepath.Join(webRoot, "tsconfig.json"):      webTSConfig(),
		filepath.Join(webRoot, "app", "globals.css"): webGlobalCSS(),
		filepath.Join(webRoot, "app", "layout.tsx"):  webRootLayout(opts),
		filepath.Join(webRoot, "app", "page.tsx"):    webLandingPage(opts),
		filepath.Join(webRoot, "lib", "utils.ts"):    webUtils(),
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
    "class-variance-authority": "^0.7.0",
    "clsx": "^2.1.0",
    "framer-motion": "^11.0.0",
    "lucide-react": "^0.303.0",
    "next": "^16.1.6",
    "react": "^19.0.0",
    "react-dom": "^19.0.0",
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
        sans: ["var(--font-dm-sans)", "system-ui", "sans-serif"],
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
  font-family: var(--font-dm-sans), system-ui, sans-serif;
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
`
}

func webRootLayout(opts Options) string {
	return fmt.Sprintf(`import type { Metadata } from "next";
import { DM_Sans, JetBrains_Mono } from "next/font/google";
import "./globals.css";

const dmSans = DM_Sans({
  subsets: ["latin"],
  variable: "--font-dm-sans",
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
      <body className={` + "`" + `${dmSans.variable} ${jetbrainsMono.variable} min-h-screen bg-background font-sans antialiased` + "`" + `}>
        {children}
      </body>
    </html>
  );
}
`, opts.ProjectName)
}

func webLandingPage(opts Options) string {
	return fmt.Sprintf(`"use client";

import { motion } from "framer-motion";
import {
  Zap,
  Database,
  Layout,
  Terminal,
  Layers,
  ArrowRight,
  Code2,
  Server,
  Palette,
  Github,
} from "lucide-react";

const ADMIN_URL = process.env.NEXT_PUBLIC_ADMIN_URL || "http://localhost:3001";

const fadeUp = {
  hidden: { opacity: 0, y: 30 },
  visible: (i: number) => ({
    opacity: 1,
    y: 0,
    transition: { delay: i * 0.1, duration: 0.5 },
  }),
};

const features = [
  {
    icon: Server,
    title: "Go Backend",
    description: "Gin web framework + GORM ORM with auto-migrations, JWT auth, and middleware.",
    color: "text-info",
    bg: "bg-info/10",
  },
  {
    icon: Layout,
    title: "Next.js Frontend",
    description: "React 19 with App Router, TypeScript, Tailwind CSS, and shadcn-style components.",
    color: "text-accent",
    bg: "bg-accent/10",
  },
  {
    icon: Palette,
    title: "Admin Panel",
    description: "Resource-based admin dashboard with CRUD, data tables, forms, and widgets.",
    color: "text-warning",
    bg: "bg-warning/10",
  },
  {
    icon: Database,
    title: "GORM Studio",
    description: "Built-in visual database browser. Inspect tables, run queries, explore data.",
    color: "text-success",
    bg: "bg-success/10",
  },
  {
    icon: Terminal,
    title: "CLI Generator",
    description: "Generate full-stack resources with one command. Models, APIs, types, and UI.",
    color: "text-danger",
    bg: "bg-danger/10",
  },
  {
    icon: Layers,
    title: "Batteries Included",
    description: "Redis caching, S3 storage, email, background jobs, cron, and AI integration.",
    color: "text-info",
    bg: "bg-info/10",
  },
];

const steps = [
  {
    step: "01",
    title: "Scaffold",
    description: "Run grit new my-app to scaffold a full monorepo with Go API, Next.js web, and admin panel.",
    code: "grit new my-project",
  },
  {
    step: "02",
    title: "Generate",
    description: "Generate full-stack resources with models, handlers, schemas, hooks, and admin pages.",
    code: "grit generate resource Post",
  },
  {
    step: "03",
    title: "Ship",
    description: "Docker Compose for local dev, standalone builds for production. Deploy anywhere.",
    code: "docker compose up -d",
  },
];

export default function LandingPage() {
  return (
    <div className="min-h-screen">
      {/* ─── Navbar ─────────────────────────────────────────── */}
      <nav className="sticky top-0 z-50 border-b border-border/50 bg-background/80 backdrop-blur-lg">
        <div className="mx-auto flex h-16 max-w-6xl items-center justify-between px-6">
          <div className="flex items-center gap-2">
            <span className="text-xl font-bold text-accent">G</span>
            <span className="text-xl font-bold text-accent">rit</span>
          </div>
          <div className="hidden md:flex items-center gap-8 text-sm">
            <a href="#features" className="text-text-secondary hover:text-foreground transition-colors">Features</a>
            <a href="#how-it-works" className="text-text-secondary hover:text-foreground transition-colors">How it works</a>
            <a href="#stack" className="text-text-secondary hover:text-foreground transition-colors">Stack</a>
          </div>
          <div className="flex items-center gap-3">
            <a
              href={` + "`" + `${ADMIN_URL}/login` + "`" + `}
              className="text-sm font-medium text-text-secondary hover:text-foreground transition-colors"
            >
              Log In
            </a>
            <a
              href={` + "`" + `${ADMIN_URL}/sign-up` + "`" + `}
              className="rounded-lg bg-accent px-4 py-2 text-sm font-medium text-white hover:bg-accent-hover transition-colors"
            >
              Sign Up
            </a>
          </div>
        </div>
      </nav>

      {/* ─── Hero ───────────────────────────────────────────── */}
      <section className="relative overflow-hidden">
        <div className="absolute inset-0 bg-gradient-to-b from-accent/5 via-transparent to-transparent" />
        <div className="relative mx-auto max-w-6xl px-6 py-24 sm:py-32 lg:py-40">
          <motion.div
            className="mx-auto max-w-3xl text-center"
            initial="hidden"
            animate="visible"
            variants={fadeUp}
            custom={0}
          >
            <div className="mb-6 inline-flex items-center gap-2 rounded-full border border-border bg-bg-secondary px-4 py-1.5 text-sm text-text-secondary">
              <Zap className="h-4 w-4 text-accent" />
              <span>Go + React Meta-Framework</span>
            </div>
            <h1 className="text-4xl font-bold tracking-tight text-foreground sm:text-5xl lg:text-6xl">
              Build full-stack apps{" "}
              <span className="text-accent">with Grit.</span>
            </h1>
            <p className="mt-6 text-lg text-text-secondary leading-relaxed sm:text-xl">
              A production-ready framework combining Go backend, Next.js frontend, and an admin panel.
              Scaffold, generate, and ship — all from the CLI.
            </p>
            <div className="mt-10 flex flex-col sm:flex-row items-center justify-center gap-4">
              <a
                href={` + "`" + `${ADMIN_URL}/sign-up` + "`" + `}
                className="flex items-center gap-2 rounded-lg bg-accent px-6 py-3 text-sm font-semibold text-white hover:bg-accent-hover transition-colors"
              >
                Get Started <ArrowRight className="h-4 w-4" />
              </a>
              <a
                href="https://github.com/MUKE-coder/grit"
                target="_blank"
                rel="noopener noreferrer"
                className="flex items-center gap-2 rounded-lg border border-border bg-bg-secondary px-6 py-3 text-sm font-semibold text-foreground hover:bg-bg-hover transition-colors"
              >
                <Github className="h-4 w-4" /> View on GitHub
              </a>
            </div>
          </motion.div>

          {/* Terminal preview */}
          <motion.div
            className="mx-auto mt-16 max-w-2xl"
            initial={{ opacity: 0, y: 40 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ delay: 0.3, duration: 0.6 }}
          >
            <div className="rounded-xl border border-border bg-bg-secondary shadow-2xl overflow-hidden">
              <div className="flex items-center gap-2 border-b border-border px-4 py-3">
                <div className="h-3 w-3 rounded-full bg-danger/60" />
                <div className="h-3 w-3 rounded-full bg-warning/60" />
                <div className="h-3 w-3 rounded-full bg-success/60" />
                <span className="ml-2 text-xs text-text-muted font-mono">terminal</span>
              </div>
              <div className="p-6 font-mono text-sm space-y-2">
                <p><span className="text-success">$</span> <span className="text-foreground">grit new my-saas</span></p>
                <p className="text-text-muted">
                  Creating project structure...<br />
                  Scaffolding Go API...<br />
                  Scaffolding Next.js web app...<br />
                  Scaffolding admin panel...<br />
                  Setting up shared types...<br />
                  Writing Docker Compose...
                </p>
                <p className="text-success">Done! Project created at ./my-saas</p>
                <p className="mt-2"><span className="text-success">$</span> <span className="text-foreground">cd my-saas && docker compose up -d</span></p>
                <p className="text-text-muted">Starting PostgreSQL, Redis, MinIO...</p>
                <p className="text-success">All services running.</p>
              </div>
            </div>
          </motion.div>
        </div>
      </section>

      {/* ─── Features ───────────────────────────────────────── */}
      <section id="features" className="border-t border-border/50 bg-bg-secondary/30">
        <div className="mx-auto max-w-6xl px-6 py-24">
          <motion.div
            className="text-center mb-16"
            initial="hidden"
            whileInView="visible"
            viewport={{ once: true }}
            variants={fadeUp}
            custom={0}
          >
            <h2 className="text-3xl font-bold text-foreground sm:text-4xl">
              Everything you need
            </h2>
            <p className="mt-4 text-text-secondary text-lg max-w-2xl mx-auto">
              A complete toolkit for building production-ready full-stack applications.
            </p>
          </motion.div>

          <div className="grid gap-6 sm:grid-cols-2 lg:grid-cols-3">
            {features.map((feature, i) => (
              <motion.div
                key={feature.title}
                className="rounded-xl border border-border bg-bg-secondary p-6 hover:border-accent/30 transition-colors"
                initial="hidden"
                whileInView="visible"
                viewport={{ once: true }}
                variants={fadeUp}
                custom={i}
              >
                <div className={` + "`" + `flex h-10 w-10 items-center justify-center rounded-lg ${feature.bg} mb-4` + "`" + `}>
                  <feature.icon className={` + "`" + `h-5 w-5 ${feature.color}` + "`" + `} />
                </div>
                <h3 className="text-lg font-semibold text-foreground">{feature.title}</h3>
                <p className="mt-2 text-sm text-text-secondary leading-relaxed">{feature.description}</p>
              </motion.div>
            ))}
          </div>
        </div>
      </section>

      {/* ─── How It Works ───────────────────────────────────── */}
      <section id="how-it-works" className="border-t border-border/50">
        <div className="mx-auto max-w-6xl px-6 py-24">
          <motion.div
            className="text-center mb-16"
            initial="hidden"
            whileInView="visible"
            viewport={{ once: true }}
            variants={fadeUp}
            custom={0}
          >
            <h2 className="text-3xl font-bold text-foreground sm:text-4xl">
              Three steps to ship
            </h2>
            <p className="mt-4 text-text-secondary text-lg max-w-2xl mx-auto">
              From zero to production-ready in minutes, not weeks.
            </p>
          </motion.div>

          <div className="grid gap-8 lg:grid-cols-3">
            {steps.map((step, i) => (
              <motion.div
                key={step.step}
                className="relative rounded-xl border border-border bg-bg-secondary p-6"
                initial="hidden"
                whileInView="visible"
                viewport={{ once: true }}
                variants={fadeUp}
                custom={i}
              >
                <span className="text-5xl font-black text-accent/10">{step.step}</span>
                <h3 className="mt-2 text-xl font-semibold text-foreground">{step.title}</h3>
                <p className="mt-2 text-sm text-text-secondary leading-relaxed">{step.description}</p>
                <div className="mt-4 rounded-lg bg-bg-tertiary px-4 py-2.5 font-mono text-sm">
                  <span className="text-success">$</span>{" "}
                  <span className="text-foreground">{step.code}</span>
                </div>
              </motion.div>
            ))}
          </div>
        </div>
      </section>

      {/* ─── Tech Stack ─────────────────────────────────────── */}
      <section id="stack" className="border-t border-border/50 bg-bg-secondary/30">
        <div className="mx-auto max-w-6xl px-6 py-24">
          <motion.div
            className="text-center mb-16"
            initial="hidden"
            whileInView="visible"
            viewport={{ once: true }}
            variants={fadeUp}
            custom={0}
          >
            <h2 className="text-3xl font-bold text-foreground sm:text-4xl">
              Built on proven tech
            </h2>
            <p className="mt-4 text-text-secondary text-lg max-w-2xl mx-auto">
              Every layer uses battle-tested, industry-standard technologies.
            </p>
          </motion.div>

          <motion.div
            className="grid grid-cols-2 gap-4 sm:grid-cols-3 lg:grid-cols-4"
            initial="hidden"
            whileInView="visible"
            viewport={{ once: true }}
            variants={fadeUp}
            custom={1}
          >
            {[
              { name: "Go", desc: "Backend language" },
              { name: "Gin", desc: "Web framework" },
              { name: "GORM", desc: "ORM" },
              { name: "PostgreSQL", desc: "Database" },
              { name: "Next.js", desc: "React framework" },
              { name: "TypeScript", desc: "Type safety" },
              { name: "Tailwind CSS", desc: "Styling" },
              { name: "React Query", desc: "Data fetching" },
              { name: "Redis", desc: "Cache + queues" },
              { name: "Docker", desc: "Containerization" },
              { name: "Turborepo", desc: "Monorepo" },
              { name: "Zod", desc: "Validation" },
            ].map((tech) => (
              <div
                key={tech.name}
                className="rounded-lg border border-border bg-bg-secondary p-4 text-center hover:border-accent/30 transition-colors"
              >
                <p className="font-semibold text-foreground">{tech.name}</p>
                <p className="text-xs text-text-muted mt-1">{tech.desc}</p>
              </div>
            ))}
          </motion.div>
        </div>
      </section>

      {/* ─── CTA ────────────────────────────────────────────── */}
      <section className="border-t border-border/50">
        <div className="mx-auto max-w-6xl px-6 py-24">
          <motion.div
            className="rounded-2xl border border-border bg-gradient-to-r from-accent/10 via-bg-secondary to-accent/5 p-12 text-center"
            initial="hidden"
            whileInView="visible"
            viewport={{ once: true }}
            variants={fadeUp}
            custom={0}
          >
            <Code2 className="h-10 w-10 text-accent mx-auto mb-4" />
            <h2 className="text-3xl font-bold text-foreground">
              Ready to build?
            </h2>
            <p className="mt-4 text-text-secondary text-lg max-w-xl mx-auto">
              Get started with Grit and ship your next full-stack app in record time.
            </p>
            <div className="mt-8 flex flex-col sm:flex-row items-center justify-center gap-4">
              <a
                href={` + "`" + `${ADMIN_URL}/sign-up` + "`" + `}
                className="flex items-center gap-2 rounded-lg bg-accent px-8 py-3 text-sm font-semibold text-white hover:bg-accent-hover transition-colors"
              >
                Get Started Free <ArrowRight className="h-4 w-4" />
              </a>
              <a
                href="https://gritframework.dev"
                target="_blank"
                rel="noopener noreferrer"
                className="text-sm font-medium text-accent hover:text-accent-hover transition-colors"
              >
                Read the docs →
              </a>
            </div>
          </motion.div>
        </div>
      </section>

      {/* ─── Footer ─────────────────────────────────────────── */}
      <footer className="border-t border-border/50">
        <div className="mx-auto max-w-6xl px-6 py-12">
          <div className="flex flex-col md:flex-row items-center justify-between gap-6">
            <div className="flex items-center gap-2">
              <span className="text-lg font-bold text-accent">Grit</span>
              <span className="text-text-muted text-sm">— Go + React. Built with Grit.</span>
            </div>
            <div className="flex items-center gap-6 text-sm text-text-secondary">
              <a href="https://gritframework.dev" target="_blank" rel="noopener noreferrer" className="hover:text-foreground transition-colors">
                Docs
              </a>
              <a href="https://github.com/MUKE-coder/grit" target="_blank" rel="noopener noreferrer" className="hover:text-foreground transition-colors">
                GitHub
              </a>
              <a href={` + "`" + `${ADMIN_URL}/login` + "`" + `} className="hover:text-foreground transition-colors">
                Admin
              </a>
            </div>
          </div>
          <div className="mt-8 pt-8 border-t border-border/50 text-center text-xs text-text-muted">
            Built with Grit by MUKE-coder. Open source under MIT License.
          </div>
        </div>
      </footer>
    </div>
  );
}
`, opts.ProjectName)
}

func webUtils() string {
	return `import { clsx, type ClassValue } from "clsx";
import { twMerge } from "tailwind-merge";

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}
`
}
