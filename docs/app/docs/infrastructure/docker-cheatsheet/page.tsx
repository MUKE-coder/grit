import Link from "next/link";
import { ArrowLeft, ArrowRight } from "lucide-react";
import { Button } from "@/components/ui/button";
import { SiteHeader } from "@/components/site-header";
import { DocsSidebar } from "@/components/docs-sidebar";
import { getDocMetadata } from '@/config/docs-metadata'

export const metadata = getDocMetadata('/docs/infrastructure/docker-cheatsheet')

export default function DockerCheatSheetPage() {
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
                Infrastructure
              </span>
              <h1 className="text-4xl font-bold tracking-tight mb-4">
                Docker Cheat Sheet
              </h1>
              <p className="text-lg text-muted-foreground leading-relaxed">
                Quick reference for all Docker commands you&apos;ll need when
                working with Grit projects. Keep this page bookmarked &mdash; it
                covers starting services, viewing logs, database operations,
                Redis, MinIO, cleanup, and more.
              </p>
            </div>

            <div className="prose-grit">
              {/* Default Services & Ports */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Default Services &amp; Ports
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Every Grit project ships with these four infrastructure
                  services via Docker Compose:
                </p>
                <div className="rounded-lg border border-border/30 bg-card/30 overflow-hidden">
                  <table className="w-full text-sm">
                    <thead>
                      <tr className="border-b border-border/30 bg-accent/20">
                        <th className="text-left px-4 py-2.5 font-medium text-foreground/80">
                          Service
                        </th>
                        <th className="text-left px-4 py-2.5 font-medium text-foreground/80">
                          Port
                        </th>
                        <th className="text-left px-4 py-2.5 font-medium text-foreground/80">
                          Dashboard
                        </th>
                      </tr>
                    </thead>
                    <tbody className="text-muted-foreground">
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-medium text-foreground/90">
                          PostgreSQL
                        </td>
                        <td className="px-4 py-2.5 font-mono text-xs">5432</td>
                        <td className="px-4 py-2.5 text-muted-foreground/50">
                          &mdash;
                        </td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-medium text-foreground/90">
                          Redis
                        </td>
                        <td className="px-4 py-2.5 font-mono text-xs">6379</td>
                        <td className="px-4 py-2.5 text-muted-foreground/50">
                          &mdash;
                        </td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-medium text-foreground/90">
                          MinIO
                        </td>
                        <td className="px-4 py-2.5 font-mono text-xs">
                          9000 (API), 9001 (Console)
                        </td>
                        <td className="px-4 py-2.5">
                          <code className="text-xs font-mono text-primary/70 bg-primary/5 px-1.5 py-0.5 rounded-md border border-primary/10">
                            http://localhost:9001
                          </code>
                        </td>
                      </tr>
                      <tr>
                        <td className="px-4 py-2.5 font-medium text-foreground/90">
                          Mailhog
                        </td>
                        <td className="px-4 py-2.5 font-mono text-xs">
                          1025 (SMTP), 8025 (UI)
                        </td>
                        <td className="px-4 py-2.5">
                          <code className="text-xs font-mono text-primary/70 bg-primary/5 px-1.5 py-0.5 rounded-md border border-primary/10">
                            http://localhost:8025
                          </code>
                        </td>
                      </tr>
                    </tbody>
                  </table>
                </div>
              </div>

              {/* Starting & Stopping Services */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Starting &amp; Stopping Services
                </h2>
                <div className="space-y-3">
                  {[
                    {
                      cmd: "docker compose up -d",
                      desc: "Start all services in the background",
                    },
                    {
                      cmd: "docker stop $(docker ps -q)",
                      desc: "Stop all containers even for other projects",
                    },
                    {
                      cmd: "docker compose down",
                      desc: "Stop all services and remove containers",
                    },
                    {
                      cmd: "docker compose restart",
                      desc: "Restart all services",
                    },
                    {
                      cmd: "docker compose up -d postgres",
                      desc: "Start a specific service only",
                    },
                    {
                      cmd: "docker compose stop redis",
                      desc: "Stop a specific service without removing it",
                    },
                  ].map((item) => (
                    <div
                      key={item.cmd}
                      className="rounded-xl border border-border/40 bg-card/80 overflow-hidden"
                    >
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
                      <div className="px-5 py-3 font-mono text-sm">
                        <div>
                          <span className="text-primary/50 select-none">
                            ${" "}
                          </span>
                          <span className="text-foreground/80">{item.cmd}</span>
                        </div>
                        <p className="text-xs text-muted-foreground/50 mt-1">
                          {item.desc}
                        </p>
                      </div>
                    </div>
                  ))}
                </div>
              </div>

              {/* Checking Service Status */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Checking Service Status
                </h2>
                <div className="space-y-3">
                  {[
                    { cmd: "docker ps", desc: "List all running containers" },
                    {
                      cmd: "docker ps -a",
                      desc: "List all containers including stopped ones",
                    },
                    {
                      cmd: "docker compose ps",
                      desc: "List services for the current project",
                    },
                    {
                      cmd: "docker stats",
                      desc: "Live resource usage (CPU, memory, network I/O)",
                    },
                  ].map((item) => (
                    <div
                      key={item.cmd}
                      className="rounded-xl border border-border/40 bg-card/80 overflow-hidden"
                    >
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
                      <div className="px-5 py-3 font-mono text-sm">
                        <div>
                          <span className="text-primary/50 select-none">
                            ${" "}
                          </span>
                          <span className="text-foreground/80">{item.cmd}</span>
                        </div>
                        <p className="text-xs text-muted-foreground/50 mt-1">
                          {item.desc}
                        </p>
                      </div>
                    </div>
                  ))}
                </div>
              </div>

              {/* Viewing Logs */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Viewing Logs
                </h2>
                <div className="space-y-3">
                  {[
                    {
                      cmd: "docker compose logs",
                      desc: "Show all service logs",
                    },
                    {
                      cmd: "docker compose logs -f",
                      desc: "Follow logs live (stream in real time)",
                    },
                    {
                      cmd: "docker compose logs postgres",
                      desc: "Show logs for a specific service",
                    },
                    {
                      cmd: "docker compose logs --tail 50 redis",
                      desc: "Show only the last 50 lines for a service",
                    },
                  ].map((item) => (
                    <div
                      key={item.cmd}
                      className="rounded-xl border border-border/40 bg-card/80 overflow-hidden"
                    >
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
                      <div className="px-5 py-3 font-mono text-sm">
                        <div>
                          <span className="text-primary/50 select-none">
                            ${" "}
                          </span>
                          <span className="text-foreground/80">{item.cmd}</span>
                        </div>
                        <p className="text-xs text-muted-foreground/50 mt-1">
                          {item.desc}
                        </p>
                      </div>
                    </div>
                  ))}
                </div>
              </div>

              {/* Database Operations */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Database Operations
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Common commands for working with the PostgreSQL container
                  directly. Replace{" "}
                  <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">
                    grit_dev
                  </code>{" "}
                  with your actual database name if you changed it.
                </p>
                <div className="space-y-3">
                  {[
                    {
                      cmd: "docker compose exec postgres psql -U grit grit_dev",
                      desc: "Open an interactive psql shell",
                    },
                    {
                      cmd: "docker compose exec postgres pg_dump -U grit grit_dev > backup.sql",
                      desc: "Backup the entire database to a SQL file",
                    },
                    {
                      cmd: "docker compose exec -T postgres psql -U grit grit_dev < backup.sql",
                      desc: "Restore database from a SQL backup file",
                    },
                    {
                      cmd: 'docker compose exec postgres psql -U grit -c "\\dt" grit_dev',
                      desc: "List all tables in the database",
                    },
                  ].map((item) => (
                    <div
                      key={item.cmd}
                      className="rounded-xl border border-border/40 bg-card/80 overflow-hidden"
                    >
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
                      <div className="px-5 py-3 font-mono text-sm">
                        <div>
                          <span className="text-primary/50 select-none">
                            ${" "}
                          </span>
                          <span className="text-foreground/80">{item.cmd}</span>
                        </div>
                        <p className="text-xs text-muted-foreground/50 mt-1">
                          {item.desc}
                        </p>
                      </div>
                    </div>
                  ))}
                </div>
              </div>

              {/* Redis Operations */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Redis Operations
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Commands for inspecting and managing the Redis cache,
                  sessions, and job queue data.
                </p>
                <div className="space-y-3">
                  {[
                    {
                      cmd: "docker compose exec redis redis-cli",
                      desc: "Open an interactive Redis CLI session",
                    },
                    {
                      cmd: 'docker compose exec redis redis-cli KEYS "*"',
                      desc: "List all keys stored in Redis",
                    },
                    {
                      cmd: "docker compose exec redis redis-cli FLUSHALL",
                      desc: "Clear all Redis data (cache, sessions, queues)",
                    },
                    {
                      cmd: "docker compose exec redis redis-cli INFO memory",
                      desc: "Show Redis memory usage statistics",
                    },
                  ].map((item) => (
                    <div
                      key={item.cmd}
                      className="rounded-xl border border-border/40 bg-card/80 overflow-hidden"
                    >
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
                      <div className="px-5 py-3 font-mono text-sm">
                        <div>
                          <span className="text-primary/50 select-none">
                            ${" "}
                          </span>
                          <span className="text-foreground/80">{item.cmd}</span>
                        </div>
                        <p className="text-xs text-muted-foreground/50 mt-1">
                          {item.desc}
                        </p>
                      </div>
                    </div>
                  ))}
                </div>
                <div className="p-4 rounded-lg border border-yellow-500/20 bg-yellow-500/5 mt-4">
                  <p className="text-sm text-muted-foreground leading-relaxed">
                    <strong className="text-yellow-500/90">Warning:</strong>{" "}
                    <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">
                      FLUSHALL
                    </code>{" "}
                    deletes everything in Redis including active sessions and
                    pending job queue entries. Use with caution in shared
                    environments.
                  </p>
                </div>
              </div>

              {/* MinIO (File Storage) */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  MinIO (File Storage)
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  MinIO provides S3-compatible file storage for development.
                  Access the web console to manage buckets and files visually.
                </p>
                <div className="space-y-4">
                  <div className="p-4 rounded-lg border border-border/30 bg-card/30">
                    <div className="flex items-center justify-between mb-1.5">
                      <h3 className="text-sm font-semibold">MinIO Console</h3>
                      <code className="text-xs font-mono text-primary/70 bg-primary/5 px-2 py-0.5 rounded-md border border-primary/10">
                        http://localhost:9001
                      </code>
                    </div>
                    <p className="text-xs text-muted-foreground/70 leading-relaxed">
                      Web-based file browser. Create buckets, upload files, and
                      manage access policies. Default credentials:{" "}
                      <code className="text-xs font-mono bg-accent/50 px-1 py-0.5 rounded">
                        minioadmin
                      </code>{" "}
                      /{" "}
                      <code className="text-xs font-mono bg-accent/50 px-1 py-0.5 rounded">
                        minioadmin
                      </code>
                    </p>
                  </div>
                </div>
                <div className="mt-4 rounded-xl border border-border/40 bg-card/80 overflow-hidden">
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
                  <div className="px-5 py-3 font-mono text-sm">
                    <div>
                      <span className="text-primary/50 select-none">$ </span>
                      <span className="text-foreground/80">
                        docker compose exec minio mc alias set local
                        http://localhost:9000 minioadmin minioadmin
                      </span>
                    </div>
                    <p className="text-xs text-muted-foreground/50 mt-1">
                      Set up the MinIO client alias for CLI operations
                    </p>
                  </div>
                </div>
              </div>

              {/* Cleanup & Reset */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Cleanup &amp; Reset
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Commands for cleaning up containers, images, and volumes. Use
                  these when you need a fresh start or want to reclaim disk
                  space.
                </p>
                <div className="space-y-3">
                  {[
                    {
                      cmd: "docker compose down -v",
                      desc: "Stop all services and delete volumes (RESETS ALL DATA)",
                    },
                    {
                      cmd: "docker system prune",
                      desc: "Remove all unused containers, networks, and dangling images",
                    },
                    {
                      cmd: "docker volume ls",
                      desc: "List all Docker volumes",
                    },
                    {
                      cmd: "docker volume rm <volume-name>",
                      desc: "Remove a specific volume by name",
                    },
                  ].map((item) => (
                    <div
                      key={item.cmd}
                      className="rounded-xl border border-border/40 bg-card/80 overflow-hidden"
                    >
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
                      <div className="px-5 py-3 font-mono text-sm">
                        <div>
                          <span className="text-primary/50 select-none">
                            ${" "}
                          </span>
                          <span className="text-foreground/80">{item.cmd}</span>
                        </div>
                        <p className="text-xs text-muted-foreground/50 mt-1">
                          {item.desc}
                        </p>
                      </div>
                    </div>
                  ))}
                </div>
                <div className="p-4 rounded-lg border border-red-500/20 bg-red-500/5 mt-4">
                  <p className="text-sm text-muted-foreground leading-relaxed">
                    <strong className="text-red-500/90">Danger:</strong>{" "}
                    <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">
                      docker compose down -v
                    </code>{" "}
                    permanently deletes all volume data including your database,
                    Redis cache, and uploaded files. Make sure to back up any
                    important data first.
                  </p>
                </div>
              </div>

              {/* Port Conflicts */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Port Conflicts
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  If a service fails to start because a port is already in use,
                  these commands help diagnose and resolve the issue.
                </p>
                <div className="space-y-3">
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
                    <div className="px-5 py-3 font-mono text-sm">
                      <div>
                        <span className="text-primary/50 select-none">$ </span>
                        <span className="text-foreground/80">{`docker ps --format "table {{.Names}}\\t{{.Ports}}"`}</span>
                      </div>
                      <p className="text-xs text-muted-foreground/50 mt-1">
                        Show which containers are using which ports
                      </p>
                    </div>
                  </div>
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
                    <div className="px-5 py-3 font-mono text-sm">
                      <div>
                        <span className="text-primary/50 select-none">$ </span>
                        <span className="text-foreground/80">
                          docker stop $(docker ps -q)
                        </span>
                      </div>
                      <p className="text-xs text-muted-foreground/50 mt-1">
                        Stop all running containers (nuclear option)
                      </p>
                    </div>
                  </div>
                </div>
                <p className="text-sm text-muted-foreground/60 mt-3">
                  For more detailed troubleshooting steps including OS-specific
                  port debugging, see the{" "}
                  <Link
                    href="/docs/getting-started/troubleshooting"
                    className="text-primary hover:underline"
                  >
                    Troubleshooting
                  </Link>{" "}
                  page.
                </p>
              </div>

              {/* Production Docker */}
              <div className="mb-10">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Production Docker
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Grit ships with a separate production compose file that builds
                  and runs the full application stack. These commands use{" "}
                  <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">
                    docker-compose.prod.yml
                  </code>{" "}
                  instead of the default development file.
                </p>
                <div className="space-y-3">
                  {[
                    {
                      cmd: "docker compose -f docker-compose.prod.yml up -d",
                      desc: "Start all production services in the background",
                    },
                    {
                      cmd: "docker compose -f docker-compose.prod.yml logs -f api",
                      desc: "Follow production API logs",
                    },
                  ].map((item) => (
                    <div
                      key={item.cmd}
                      className="rounded-xl border border-border/40 bg-card/80 overflow-hidden"
                    >
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
                      <div className="px-5 py-3 font-mono text-sm">
                        <div>
                          <span className="text-primary/50 select-none">
                            ${" "}
                          </span>
                          <span className="text-foreground/80">{item.cmd}</span>
                        </div>
                        <p className="text-xs text-muted-foreground/50 mt-1">
                          {item.desc}
                        </p>
                      </div>
                    </div>
                  ))}
                </div>
                <p className="text-sm text-muted-foreground/60 mt-3">
                  For full production deployment guides including cloud
                  providers and CI/CD, see the{" "}
                  <Link
                    href="/docs/infrastructure/deployment"
                    className="text-primary hover:underline"
                  >
                    Deployment
                  </Link>{" "}
                  page.
                </p>
              </div>

              {/* Nav */}
              <div className="flex items-center justify-between pt-6 border-t border-border/30">
                <Button
                  variant="ghost"
                  size="sm"
                  asChild
                  className="text-muted-foreground/60 hover:text-foreground"
                >
                  <Link
                    href="/docs/infrastructure/database"
                    className="gap-1.5"
                  >
                    <ArrowLeft className="h-3.5 w-3.5" />
                    Database &amp; Migrations
                  </Link>
                </Button>
                <Button
                  variant="ghost"
                  size="sm"
                  asChild
                  className="text-muted-foreground/60 hover:text-foreground"
                >
                  <Link
                    href="/docs/getting-started/troubleshooting"
                    className="gap-1.5"
                  >
                    Troubleshooting
                    <ArrowRight className="h-3.5 w-3.5" />
                  </Link>
                </Button>
              </div>
            </div>
          </div>
        </div>
      </main>
    </div>
  );
}
