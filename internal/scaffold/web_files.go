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
	return `import { ArrowRight, Github } from "lucide-react";

const ADMIN_URL = process.env.NEXT_PUBLIC_ADMIN_URL || "http://localhost:3001";
const DOCS_URL = "https://grit-vert.vercel.app/docs";

export default function HomePage() {
  return (
    <div className="min-h-screen flex flex-col">
      {/* Navbar */}
      <nav className="border-b border-border/50 bg-background/80 backdrop-blur-lg">
        <div className="mx-auto flex h-16 max-w-5xl items-center justify-between px-6">
          <div className="flex items-center gap-2.5">
            <div className="flex h-8 w-8 items-center justify-center rounded-lg bg-accent/15 border border-accent/20">
              <span className="text-accent font-mono font-bold text-sm">G</span>
            </div>
            <span className="text-lg font-bold tracking-tight">Grit</span>
          </div>
          <div className="flex items-center gap-4">
            <a
              href={` + "`" + `${DOCS_URL}` + "`" + `}
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
        </div>
      </nav>

      {/* Hero */}
      <main className="flex-1 flex items-center justify-center">
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
      </main>

      {/* Footer */}
      <footer className="border-t border-border/50 py-6">
        <div className="mx-auto max-w-5xl px-6 flex items-center justify-between">
          <span className="text-xs text-text-muted">
            Built with Grit — Go + React. Open source under MIT.
          </span>
          <div className="flex items-center gap-4 text-xs text-text-muted">
            <a href={` + "`" + `${DOCS_URL}` + "`" + `} target="_blank" rel="noopener noreferrer" className="hover:text-foreground transition-colors">
              Docs
            </a>
            <a href="https://github.com/MUKE-coder/grit" target="_blank" rel="noopener noreferrer" className="hover:text-foreground transition-colors">
              GitHub
            </a>
          </div>
        </div>
      </footer>
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
