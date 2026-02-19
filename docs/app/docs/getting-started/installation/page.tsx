import Link from "next/link";
import { ArrowRight } from "lucide-react";
import { Button } from "@/components/ui/button";
import { SiteHeader } from "@/components/site-header";
import { DocsSidebar } from "@/components/docs-sidebar";
import { CodeBlock } from '@/components/code-block'

export default function InstallationPage() {
  return (
    <div className="min-h-screen bg-background">
      <SiteHeader />
      <DocsSidebar />

      <main className="lg:pl-64">
        <div className="container max-w-screen-xl py-10 px-6">
          <div className="max-w-3xl">
            {/* Header */}
            <div className="mb-10">
              <span className="tag-mono text-primary/80 mb-3 block">
                Getting Started
              </span>
              <h1 className="text-4xl font-bold tracking-tight mb-4">
                Installation
              </h1>
              <p className="text-lg text-muted-foreground leading-relaxed">
                Everything you need to install and configure before creating
                your first Grit project. This guide covers system requirements,
                tool installation, and the Grit CLI setup.
              </p>
            </div>

            <div className="prose-grit">
              {/* System Requirements */}
              <h2>System Requirements</h2>
              <p>
                Grit runs on macOS, Linux, and Windows. You need the following
                tools installed:
              </p>
            </div>

            <div className="mb-10">
              <table className="w-full border-collapse">
                <thead>
                  <tr className="border-b border-border">
                    <th className="text-left text-sm font-semibold text-foreground pb-3 pr-4">
                      Tool
                    </th>
                    <th className="text-left text-sm font-semibold text-foreground pb-3 pr-4">
                      Minimum Version
                    </th>
                    <th className="text-left text-sm font-semibold text-foreground pb-3 pr-4">
                      Required For
                    </th>
                    <th className="text-left text-sm font-semibold text-foreground pb-3">
                      Check Command
                    </th>
                  </tr>
                </thead>
                <tbody>
                  {[
                    {
                      tool: "Go",
                      version: "1.21+",
                      purpose: "Backend API + CLI",
                      check: "go version",
                    },
                    {
                      tool: "Node.js",
                      version: "18+",
                      purpose: "Frontend apps",
                      check: "node --version",
                    },
                    {
                      tool: "pnpm",
                      version: "8+",
                      purpose: "Package management",
                      check: "pnpm --version",
                    },
                    {
                      tool: "Docker",
                      version: "Latest",
                      purpose: "PostgreSQL, Redis, MinIO",
                      check: "docker --version",
                    },
                    {
                      tool: "Docker Compose",
                      version: "V2",
                      purpose: "Service orchestration",
                      check: "docker compose version",
                    },
                    {
                      tool: "Git",
                      version: "Any",
                      purpose: "Version control",
                      check: "git --version",
                    },
                  ].map((row) => (
                    <tr key={row.tool} className="border-b border-border/50">
                      <td className="text-sm text-foreground py-3 pr-4 font-medium">
                        {row.tool}
                      </td>
                      <td className="text-sm text-muted-foreground py-3 pr-4 font-mono">
                        {row.version}
                      </td>
                      <td className="text-sm text-muted-foreground py-3 pr-4">
                        {row.purpose}
                      </td>
                      <td className="text-sm text-muted-foreground py-3 font-mono text-primary/60">
                        {row.check}
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>

            <div className="prose-grit">
              {/* Installing Go */}
              <h2>Installing Go</h2>
              <p>
                Go is required for both the Grit CLI and the backend API.
                Download the latest version from the official website or use a
                version manager.
              </p>
            </div>

            {/* macOS */}
            <div className="mb-4">
              <p className="text-sm font-semibold mb-2 text-foreground/80">
                macOS (Homebrew)
              </p>
              <div className="rounded-xl border border-border/40 bg-card/80 overflow-hidden">
                <div className="flex items-center gap-2 px-4 py-2.5 border-b border-border/30 bg-accent/30">
                  <div className="flex items-center gap-1.5">
                    <div className="h-2.5 w-2.5 rounded-full bg-[#ff5f57]" />
                    <div className="h-2.5 w-2.5 rounded-full bg-[#febc2e]" />
                    <div className="h-2.5 w-2.5 rounded-full bg-[#28c840]" />
                  </div>
                  <span className="ml-2 text-[11px] font-mono text-muted-foreground/40">
                    terminal
                  </span>
                </div>
                <div className="p-5 font-mono text-sm">
                  <span className="text-primary/50 select-none">$ </span>
                  <span className="text-foreground/80">brew install go</span>
                </div>
              </div>
            </div>

            {/* Linux */}
            <div className="mb-4">
              <p className="text-sm font-semibold mb-2 text-foreground/80">
                Linux (Ubuntu/Debian)
              </p>
              <div className="rounded-xl border border-border/40 bg-card/80 overflow-hidden">
                <div className="flex items-center gap-2 px-4 py-2.5 border-b border-border/30 bg-accent/30">
                  <div className="flex items-center gap-1.5">
                    <div className="h-2.5 w-2.5 rounded-full bg-[#ff5f57]" />
                    <div className="h-2.5 w-2.5 rounded-full bg-[#febc2e]" />
                    <div className="h-2.5 w-2.5 rounded-full bg-[#28c840]" />
                  </div>
                  <span className="ml-2 text-[11px] font-mono text-muted-foreground/40">
                    terminal
                  </span>
                </div>
                <div className="p-5 font-mono text-sm space-y-1">
                  <div>
                    <span className="text-primary/50 select-none">$ </span>
                    <span className="text-foreground/80">
                      wget https://go.dev/dl/go1.22.0.linux-amd64.tar.gz
                    </span>
                  </div>
                  <div>
                    <span className="text-primary/50 select-none">$ </span>
                    <span className="text-foreground/80">
                      sudo tar -C /usr/local -xzf go1.22.0.linux-amd64.tar.gz
                    </span>
                  </div>
                  <div>
                    <span className="text-primary/50 select-none">$ </span>
                    <span className="text-foreground/80">
                      export PATH=$PATH:/usr/local/go/bin
                    </span>
                  </div>
                </div>
              </div>
            </div>

            {/* Windows */}
            <div className="mb-8">
              <p className="text-sm font-semibold mb-2 text-foreground/80">
                Windows
              </p>
              <div className="prose-grit">
                <p>
                  Download the MSI installer from{" "}
                  <a href="https://go.dev/dl/" target="_blank" rel="noreferrer">
                    go.dev/dl
                  </a>{" "}
                  and run it. The installer adds Go to your system PATH
                  automatically. Alternatively, use{" "}
                  <code>winget install GoLang.Go</code> or{" "}
                  <code>scoop install go</code>.
                </p>
              </div>
            </div>

            <div className="prose-grit">
              {/* Installing Node.js */}
              <h2>Installing Node.js</h2>
              <p>
                Node.js 18 or later is required for the Next.js frontend apps.
                We recommend using a version manager like <strong>nvm</strong>{" "}
                or <strong>fnm</strong>.
              </p>
            </div>

            <div className="mb-4">
              <p className="text-sm font-semibold mb-2 text-foreground/80">
                Using nvm (recommended)
              </p>
              <div className="rounded-xl border border-border/40 bg-card/80 overflow-hidden">
                <div className="flex items-center gap-2 px-4 py-2.5 border-b border-border/30 bg-accent/30">
                  <div className="flex items-center gap-1.5">
                    <div className="h-2.5 w-2.5 rounded-full bg-[#ff5f57]" />
                    <div className="h-2.5 w-2.5 rounded-full bg-[#febc2e]" />
                    <div className="h-2.5 w-2.5 rounded-full bg-[#28c840]" />
                  </div>
                  <span className="ml-2 text-[11px] font-mono text-muted-foreground/40">
                    terminal
                  </span>
                </div>
                <div className="p-5 font-mono text-sm space-y-1">
                  <div>
                    <span className="text-primary/50 select-none">$ </span>
                    <span className="text-foreground/80">nvm install 20</span>
                  </div>
                  <div>
                    <span className="text-primary/50 select-none">$ </span>
                    <span className="text-foreground/80">nvm use 20</span>
                  </div>
                </div>
              </div>
            </div>

            <div className="mb-8">
              <p className="text-sm font-semibold mb-2 text-foreground/80">
                Using Homebrew (macOS)
              </p>
              <div className="rounded-xl border border-border/40 bg-card/80 overflow-hidden">
                <div className="flex items-center gap-2 px-4 py-2.5 border-b border-border/30 bg-accent/30">
                  <div className="flex items-center gap-1.5">
                    <div className="h-2.5 w-2.5 rounded-full bg-[#ff5f57]" />
                    <div className="h-2.5 w-2.5 rounded-full bg-[#febc2e]" />
                    <div className="h-2.5 w-2.5 rounded-full bg-[#28c840]" />
                  </div>
                  <span className="ml-2 text-[11px] font-mono text-muted-foreground/40">
                    terminal
                  </span>
                </div>
                <div className="p-5 font-mono text-sm">
                  <span className="text-primary/50 select-none">$ </span>
                  <span className="text-foreground/80">
                    brew install node@20
                  </span>
                </div>
              </div>
            </div>

            <div className="prose-grit">
              {/* Installing pnpm */}
              <h2>Installing pnpm</h2>
              <p>
                Grit uses pnpm for frontend package management. It is faster,
                more disk-efficient, and enforces stricter dependency resolution
                than npm or yarn.
              </p>
            </div>

            <div className="mb-8">
              <div className="rounded-xl border border-border/40 bg-card/80 overflow-hidden">
                <div className="flex items-center gap-2 px-4 py-2.5 border-b border-border/30 bg-accent/30">
                  <div className="flex items-center gap-1.5">
                    <div className="h-2.5 w-2.5 rounded-full bg-[#ff5f57]" />
                    <div className="h-2.5 w-2.5 rounded-full bg-[#febc2e]" />
                    <div className="h-2.5 w-2.5 rounded-full bg-[#28c840]" />
                  </div>
                  <span className="ml-2 text-[11px] font-mono text-muted-foreground/40">
                    terminal
                  </span>
                </div>
                <div className="p-5 font-mono text-sm space-y-1">
                  <div className="text-muted-foreground/40 text-xs mb-2">
                    # Using npm (easiest)
                  </div>
                  <div>
                    <span className="text-primary/50 select-none">$ </span>
                    <span className="text-foreground/80">
                      npm install -g pnpm
                    </span>
                  </div>
                  <div className="text-muted-foreground/40 text-xs mb-2 mt-4">
                    # Or using corepack (built into Node.js)
                  </div>
                  <div>
                    <span className="text-primary/50 select-none">$ </span>
                    <span className="text-foreground/80">corepack enable</span>
                  </div>
                  <div>
                    <span className="text-primary/50 select-none">$ </span>
                    <span className="text-foreground/80">
                      corepack prepare pnpm@latest --activate
                    </span>
                  </div>
                </div>
              </div>
            </div>

            <div className="prose-grit">
              {/* Installing Docker */}
              <h2>Installing Docker</h2>
              <p>
                Docker is used to run PostgreSQL, Redis, MinIO, and Mailhog in
                development. Install Docker Desktop for your platform:
              </p>
              <ul>
                <li>
                  <strong>macOS:</strong>{" "}
                  <a
                    href="https://docs.docker.com/desktop/install/mac-install/"
                    target="_blank"
                    rel="noreferrer"
                  >
                    Docker Desktop for Mac
                  </a>
                </li>
                <li>
                  <strong>Windows:</strong>{" "}
                  <a
                    href="https://docs.docker.com/desktop/install/windows-install/"
                    target="_blank"
                    rel="noreferrer"
                  >
                    Docker Desktop for Windows
                  </a>{" "}
                  (WSL2 backend recommended)
                </li>
                <li>
                  <strong>Linux:</strong>{" "}
                  <a
                    href="https://docs.docker.com/engine/install/"
                    target="_blank"
                    rel="noreferrer"
                  >
                    Docker Engine
                  </a>{" "}
                  +{" "}
                  <a
                    href="https://docs.docker.com/compose/install/"
                    target="_blank"
                    rel="noreferrer"
                  >
                    Docker Compose V2
                  </a>
                </li>
              </ul>
              <p>After installation, verify Docker is running:</p>
            </div>

            <div className="mb-8">
              <div className="rounded-xl border border-border/40 bg-card/80 overflow-hidden">
                <div className="flex items-center gap-2 px-4 py-2.5 border-b border-border/30 bg-accent/30">
                  <div className="flex items-center gap-1.5">
                    <div className="h-2.5 w-2.5 rounded-full bg-[#ff5f57]" />
                    <div className="h-2.5 w-2.5 rounded-full bg-[#febc2e]" />
                    <div className="h-2.5 w-2.5 rounded-full bg-[#28c840]" />
                  </div>
                  <span className="ml-2 text-[11px] font-mono text-muted-foreground/40">
                    terminal
                  </span>
                </div>
                <div className="p-5 font-mono text-sm space-y-1">
                  <div>
                    <span className="text-primary/50 select-none">$ </span>
                    <span className="text-foreground/80">docker --version</span>
                  </div>
                  <div>
                    <span className="text-primary/50 select-none">$ </span>
                    <span className="text-foreground/80">
                      docker compose version
                    </span>
                  </div>
                </div>
              </div>
            </div>

            <div className="prose-grit">
              {/* Installing the Grit CLI */}
              <h2>Installing the Grit CLI</h2>
              <p>
                The Grit CLI is a Go binary that you install globally. It
                provides the <code>grit</code> command for scaffolding projects,
                generating resources, running migrations, and syncing types.
              </p>
            </div>

            <div className="mb-4">
              <div className="rounded-xl border border-border/40 bg-card/80 overflow-hidden glow-purple-sm">
                <div className="flex items-center gap-2 px-4 py-2.5 border-b border-border/30 bg-accent/30">
                  <div className="flex items-center gap-1.5">
                    <div className="h-2.5 w-2.5 rounded-full bg-[#ff5f57]" />
                    <div className="h-2.5 w-2.5 rounded-full bg-[#febc2e]" />
                    <div className="h-2.5 w-2.5 rounded-full bg-[#28c840]" />
                  </div>
                  <span className="ml-2 text-[11px] font-mono text-muted-foreground/40">
                    terminal
                  </span>
                </div>
                <div className="p-5 font-mono text-sm">
                  <span className="text-primary/50 select-none">$ </span>
                  <span className="text-foreground/80">
                    go install github.com/MUKE-coder/grit/cmd/grit@latest
                  </span>
                </div>
              </div>
            </div>

            <div className="prose-grit">
              <p>
                This downloads, compiles, and installs the <code>grit</code>{" "}
                binary into your <code>$GOPATH/bin</code> directory (usually{" "}
                <code>~/go/bin</code>). Make sure this directory is in your
                system <code>PATH</code>.
              </p>

              {/* Verify */}
              <h3>Verifying the Installation</h3>
              <p>
                Run the following command to verify everything is installed
                correctly:
              </p>
            </div>

            <div className="mb-8">
              <div className="rounded-xl border border-border/40 bg-card/80 overflow-hidden">
                <div className="flex items-center gap-2 px-4 py-2.5 border-b border-border/30 bg-accent/30">
                  <div className="flex items-center gap-1.5">
                    <div className="h-2.5 w-2.5 rounded-full bg-[#ff5f57]" />
                    <div className="h-2.5 w-2.5 rounded-full bg-[#febc2e]" />
                    <div className="h-2.5 w-2.5 rounded-full bg-[#28c840]" />
                  </div>
                  <span className="ml-2 text-[11px] font-mono text-muted-foreground/40">
                    terminal
                  </span>
                </div>
                <div className="p-5 font-mono text-sm space-y-1">
                  <div>
                    <span className="text-primary/50 select-none">$ </span>
                    <span className="text-foreground/80">grit --help</span>
                  </div>
                  <div className="text-muted-foreground/40 text-xs mt-3 space-y-1">
                    <div className="text-primary/40">
                      {"   ____      _  _   "}
                    </div>
                    <div className="text-primary/40">
                      {"  / ___|_ __(_)| |_ "}
                    </div>
                    <div className="text-primary/40">
                      {" | |  _| '__| || __|"}
                    </div>
                    <div className="text-primary/40">
                      {" | |_| | |  | || |_ "}
                    </div>
                    <div className="text-primary/40">
                      {"  \\____|_|  |_| \\__|"}
                    </div>
                    <div className="mt-2">Go + React. Built with Grit.</div>
                  </div>
                </div>
              </div>
            </div>

            <div className="prose-grit">
              {/* Creating First Project */}
              <h2>Creating Your First Project</h2>
              <p>
                With all tools installed, you are ready to create your first
                Grit project:
              </p>
            </div>

            <div className="mb-4">
              <div className="rounded-xl border border-border/40 bg-card/80 overflow-hidden glow-purple-sm">
                <div className="flex items-center gap-2 px-4 py-2.5 border-b border-border/30 bg-accent/30">
                  <div className="flex items-center gap-1.5">
                    <div className="h-2.5 w-2.5 rounded-full bg-[#ff5f57]" />
                    <div className="h-2.5 w-2.5 rounded-full bg-[#febc2e]" />
                    <div className="h-2.5 w-2.5 rounded-full bg-[#28c840]" />
                  </div>
                  <span className="ml-2 text-[11px] font-mono text-muted-foreground/40">
                    terminal
                  </span>
                </div>
                <div className="p-5 font-mono text-sm">
                  <span className="text-primary/50 select-none">$ </span>
                  <span className="text-foreground/80">
                    grit new my-saas-app
                  </span>
                </div>
              </div>
            </div>

            <div className="prose-grit">
              {/* Project Flags */}
              <h2>Project Flags</h2>
              <p>
                The <code>grit new</code> command supports flags to customize
                what gets scaffolded:
              </p>
            </div>

            <div className="mb-8">
              <table className="w-full border-collapse">
                <thead>
                  <tr className="border-b border-border">
                    <th className="text-left text-sm font-semibold text-foreground pb-3 pr-4">
                      Flag
                    </th>
                    <th className="text-left text-sm font-semibold text-foreground pb-3 pr-4">
                      Description
                    </th>
                    <th className="text-left text-sm font-semibold text-foreground pb-3">
                      What Gets Created
                    </th>
                  </tr>
                </thead>
                <tbody>
                  {[
                    {
                      flag: "(default)",
                      desc: "Full-stack project",
                      creates:
                        "Go API + Next.js web + admin panel + shared types + Docker",
                    },
                    {
                      flag: "--api",
                      desc: "Backend only",
                      creates:
                        "Go API with auth, Docker setup, no frontend apps",
                    },
                    {
                      flag: "--full",
                      desc: "Everything + mobile",
                      creates: "Default + Expo mobile app + docs site",
                    },
                  ].map((row) => (
                    <tr key={row.flag} className="border-b border-border/50">
                      <td className="text-sm text-foreground py-3 pr-4 font-mono text-primary/70">
                        {row.flag}
                      </td>
                      <td className="text-sm text-muted-foreground py-3 pr-4">
                        {row.desc}
                      </td>
                      <td className="text-sm text-muted-foreground py-3">
                        {row.creates}
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>

            <div className="prose-grit">
              <p>Examples:</p>
            </div>

            <div className="mb-8">
              <div className="rounded-xl border border-border/40 bg-card/80 overflow-hidden">
                <div className="flex items-center gap-2 px-4 py-2.5 border-b border-border/30 bg-accent/30">
                  <div className="flex items-center gap-1.5">
                    <div className="h-2.5 w-2.5 rounded-full bg-[#ff5f57]" />
                    <div className="h-2.5 w-2.5 rounded-full bg-[#febc2e]" />
                    <div className="h-2.5 w-2.5 rounded-full bg-[#28c840]" />
                  </div>
                  <span className="ml-2 text-[11px] font-mono text-muted-foreground/40">
                    terminal
                  </span>
                </div>
                <div className="p-5 font-mono text-sm space-y-2">
                  <div className="text-muted-foreground/40 text-xs">
                    # Full-stack (default)
                  </div>
                  <div>
                    <span className="text-primary/50 select-none">$ </span>
                    <span className="text-foreground/80">grit new my-crm</span>
                  </div>
                  <div className="text-muted-foreground/40 text-xs mt-3">
                    # API only (no frontend)
                  </div>
                  <div>
                    <span className="text-primary/50 select-none">$ </span>
                    <span className="text-foreground/80">
                      grit new my-api --api
                    </span>
                  </div>
                  <div className="text-muted-foreground/40 text-xs mt-3">
                    # Full + mobile + docs
                  </div>
                  <div>
                    <span className="text-primary/50 select-none">$ </span>
                    <span className="text-foreground/80">
                      grit new my-startup --full
                    </span>
                  </div>
                </div>
              </div>
            </div>

            <div className="prose-grit">
              {/* Troubleshooting */}
              <h2>Troubleshooting</h2>

              <h3>
                <code>grit: command not found</code>
              </h3>
              <p>
                This means <code>$GOPATH/bin</code> is not in your system{" "}
                <code>PATH</code>. Add the following to your shell profile (
                <code>~/.zshrc</code>, <code>~/.bashrc</code>, or{" "}
                <code>~/.profile</code>):
              </p>
            </div>

            <div className="mb-4">
              <CodeBlock filename="~/.zshrc" code={`export PATH=$PATH:$(go env GOPATH)/bin`} />
            </div>

            <div className="prose-grit">
              <p>
                Then reload your shell: <code>source ~/.zshrc</code>. On
                Windows, make sure <code>%GOPATH%\\bin</code> is in your system
                PATH environment variable.
              </p>

              <h3>Docker daemon not running</h3>
              <p>
                If <code>docker compose up</code> fails with a connection error,
                make sure Docker Desktop is running. On Linux, start the Docker
                daemon with <code>sudo systemctl start docker</code>.
              </p>

              <h3>Port conflicts</h3>
              <p>
                If ports 5432, 6379, 8080, 3000, or 3001 are already in use,
                edit the <code>.env</code> file and{" "}
                <code>docker-compose.yml</code> to use different ports. Common
                conflicts:
              </p>
              <ul>
                <li>Port 5432 -- another PostgreSQL instance is running</li>
                <li>
                  Port 3000 -- another Next.js or React dev server is running
                </li>
                <li>Port 8080 -- another API server or proxy is running</li>
              </ul>

              <h3>pnpm install fails</h3>
              <p>
                If <code>pnpm install</code> fails with peer dependency errors,
                try deleting the lockfile and <code>node_modules</code> folders
                and running again:
              </p>
            </div>

            <div className="mb-8">
              <div className="rounded-xl border border-border/40 bg-card/80 overflow-hidden">
                <div className="flex items-center gap-2 px-4 py-2.5 border-b border-border/30 bg-accent/30">
                  <div className="flex items-center gap-1.5">
                    <div className="h-2.5 w-2.5 rounded-full bg-[#ff5f57]" />
                    <div className="h-2.5 w-2.5 rounded-full bg-[#febc2e]" />
                    <div className="h-2.5 w-2.5 rounded-full bg-[#28c840]" />
                  </div>
                  <span className="ml-2 text-[11px] font-mono text-muted-foreground/40">
                    terminal
                  </span>
                </div>
                <div className="p-5 font-mono text-sm space-y-1">
                  <div>
                    <span className="text-primary/50 select-none">$ </span>
                    <span className="text-foreground/80">
                      rm -rf node_modules pnpm-lock.yaml apps/*/node_modules
                      packages/*/node_modules
                    </span>
                  </div>
                  <div>
                    <span className="text-primary/50 select-none">$ </span>
                    <span className="text-foreground/80">pnpm install</span>
                  </div>
                </div>
              </div>
            </div>

            <div className="prose-grit">
              <h3>Cloud-only setup (no Docker)</h3>
              <p>
                If you cannot run Docker, Grit includes a{" "}
                <code>.env.cloud.example</code> file that configures the project
                to use cloud services instead:
              </p>
              <ul>
                <li>
                  <strong>Database:</strong>{" "}
                  <a href="https://neon.tech" target="_blank" rel="noreferrer">
                    Neon
                  </a>{" "}
                  (free serverless PostgreSQL)
                </li>
                <li>
                  <strong>Redis:</strong>{" "}
                  <a
                    href="https://upstash.com"
                    target="_blank"
                    rel="noreferrer"
                  >
                    Upstash
                  </a>{" "}
                  (free serverless Redis)
                </li>
                <li>
                  <strong>Storage:</strong>{" "}
                  <a
                    href="https://dash.cloudflare.com"
                    target="_blank"
                    rel="noreferrer"
                  >
                    Cloudflare R2
                  </a>{" "}
                  or{" "}
                  <a
                    href="https://www.backblaze.com/b2"
                    target="_blank"
                    rel="noreferrer"
                  >
                    Backblaze B2
                  </a>
                </li>
              </ul>
              <p>Copy the cloud example and fill in your credentials:</p>
            </div>

            <div className="mb-8">
              <div className="rounded-xl border border-border/40 bg-card/80 overflow-hidden">
                <div className="flex items-center gap-2 px-4 py-2.5 border-b border-border/30 bg-accent/30">
                  <div className="flex items-center gap-1.5">
                    <div className="h-2.5 w-2.5 rounded-full bg-[#ff5f57]" />
                    <div className="h-2.5 w-2.5 rounded-full bg-[#febc2e]" />
                    <div className="h-2.5 w-2.5 rounded-full bg-[#28c840]" />
                  </div>
                  <span className="ml-2 text-[11px] font-mono text-muted-foreground/40">
                    terminal
                  </span>
                </div>
                <div className="p-5 font-mono text-sm">
                  <span className="text-primary/50 select-none">$ </span>
                  <span className="text-foreground/80">
                    cp .env.cloud.example .env
                  </span>
                </div>
              </div>
            </div>

            {/* Nav */}
            <div className="flex flex-wrap gap-3 mt-12 pt-6 border-t border-border/30">
              <Button
                variant="outline"
                asChild
                className="border-border/60 bg-transparent hover:bg-accent/50"
              >
                <Link href="/docs/getting-started/quick-start">
                  Quick Start
                </Link>
              </Button>
              <Button asChild className="glow-purple-sm ml-auto">
                <Link href="/docs/getting-started/project-structure">
                  Project Structure
                  <ArrowRight className="ml-1.5 h-4 w-4" />
                </Link>
              </Button>
            </div>
          </div>
        </div>
      </main>
    </div>
  );
}
