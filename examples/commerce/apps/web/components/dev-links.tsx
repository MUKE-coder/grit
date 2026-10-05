"use client";

// v3.31.49 -- DevLinks renders a grid of the local URLs printed by
// `grit new` (API, GORM Studio, Sentinel, Admin, MinIO, Mailhog, ...)
// directly on the marketing site, so operators don't have to keep the
// terminal output around to find them.
//
// Only renders in development (NODE_ENV !== "production") so production
// marketing pages never leak the internal port map. The check happens
// at module level so the section disappears entirely from the prod
// bundle -- not just hidden behind a class.

import {
  Server,
  Database,
  Shield,
  Activity,
  HardDrive,
  Mail,
  LayoutDashboard,
  FileText,
  ExternalLink,
} from "lucide-react";
import { API_URL } from "@/lib/api-core";

// Where the admin panel is.
//
// A triple project runs it as its own app on 3001. A double has no second app:
// the panel is a route group in THIS app at /admin, and a link to
// http://localhost:3001 pointed at nothing, which is what a double's navbar did
// until v3.235.0. NEXT_PUBLIC_ADMIN_URL still overrides both.
const ADMIN_URL = process.env.NEXT_PUBLIC_ADMIN_URL || "http://localhost:3001";

interface DevLink {
  title: string;
  description: string;
  href: string;
  icon: React.ComponentType<{ className?: string }>;
  group: "app" | "api" | "data" | "ops";
}

const LINKS: DevLink[] = [
  { title: "Admin panel",   description: "Resource CRUD, system hub, dashboard customisation",  href: ADMIN_URL,             icon: LayoutDashboard, group: "app"  },
  { title: "API root",      description: "Health check + every JSON endpoint",                   href: API_URL,               icon: Server,          group: "api"  },
  { title: "API docs",      description: "OpenAPI / Swagger UI",                                 href: API_URL + "/docs",     icon: FileText,        group: "api"  },
  { title: "GORM Studio",   description: "Visual database browser (your tables, no SQL)",        href: API_URL + "/studio",   icon: Database,        group: "data" },
  { title: "Sentinel",      description: "Security + rate-limit dashboard",                      href: API_URL + "/sentinel/ui", icon: Shield,       group: "ops"  },
  { title: "Pulse",         description: "Observability: traces, slow queries, SLO timelines",   href: API_URL + "/pulse/ui",    icon: Activity,     group: "ops"  },
  { title: "MinIO console", description: "Object storage browser (buckets, uploads)",            href: "http://localhost:9003",  icon: HardDrive,    group: "data" },
  { title: "Mailhog",       description: "Email catcher (dev only; intercepts every outbound)",  href: "http://localhost:8025",  icon: Mail,         group: "ops"  },
];

const groupAccent: Record<DevLink["group"], string> = {
  app:  "text-accent bg-accent/10",
  api:  "text-info bg-info/10",
  data: "text-success bg-success/10",
  ops:  "text-warning bg-warning/10",
};

export function DevLinks() {
  if (process.env.NODE_ENV === "production") return null;
  return (
    <section className="py-16 border-t border-border/50">
      <div className="mx-auto max-w-5xl px-6">
        <div className="mb-8">
          <p className="text-xs font-semibold uppercase tracking-wider text-text-muted">
            Local development
          </p>
          <h2 className="mt-1 text-2xl font-bold tracking-tight text-foreground">
            Developer links
          </h2>
          <p className="mt-1 text-sm text-text-secondary">
            All the dashboards and consoles your project ships with — wired to your local ports. Hidden in production.
          </p>
        </div>
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-3">
          {LINKS.map((l) => {
            const Icon = l.icon;
            return (
              <a
                key={l.href}
                href={l.href}
                target="_blank"
                rel="noopener noreferrer"
                className="group rounded-xl border border-border bg-bg-elevated p-4 transition-colors hover:bg-bg-hover"
              >
                <div className="flex items-start justify-between">
                  <span className={"inline-flex h-9 w-9 items-center justify-center rounded-lg " + groupAccent[l.group]}>
                    <Icon className="h-4 w-4" />
                  </span>
                  <ExternalLink className="h-3.5 w-3.5 text-text-muted opacity-0 transition-opacity group-hover:opacity-100" />
                </div>
                <p className="mt-3 text-sm font-semibold text-foreground group-hover:text-accent">{l.title}</p>
                <p className="mt-0.5 text-xs text-text-muted">{l.description}</p>
                <p className="mt-2 truncate text-[11px] font-mono text-text-muted">
                  {l.href.replace(/^https?:\/\//, "")}
                </p>
              </a>
            );
          })}
        </div>
      </div>
    </section>
  );
}
