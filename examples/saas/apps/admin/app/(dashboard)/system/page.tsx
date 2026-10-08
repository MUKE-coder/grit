"use client";

import { useState } from "react";
import Link from "next/link";
import { useModules } from "@/hooks/use-modules";
import { PageHeader } from "@/components/chrome/PageHeader";
import { SYSTEM_NAV, INTERNAL_ICON } from "@/components/chrome/CollapsibleSidebar";
import {
  Activity, Bell, Calendar, Database, FileText, Mail,
  MessageSquare, Shield, ShieldCheck, TrendingUp, Upload, Link as LinkIcon,
  UserCheck, Settings, LayoutGrid, KeyRound, Workflow, Server, Trash2,
} from "@/lib/icons";

// Every operational surface is grouped under one of these tabs. The whole
// point of the hub is to keep the sidebar to just resources + one "System
// Hub" link, so ALL of these live here rather than in the rail.
type Category =
  | "Operations"
  | "Security & Access"
  | "Data & Files"
  | "Communication"
  | "Settings";

const CATEGORIES: Category[] = [
  "Operations",
  "Security & Access",
  "Data & Files",
  "Communication",
  "Settings",
];

// Plugin-injected system links (Webhooks, Impersonate, …) that aren't built-in
// tiles show under this extra tab, which only appears when something injected.
const EXTENSIONS = "Extensions";

interface SystemTile {
  href: string;
  /** Optional module; the tile is hidden when that module is disabled. */
  module?: string;
  category: Category;
  title: string;
  description: string;
  icon: React.ReactNode;
}

const TILES: SystemTile[] = [
  // ── Operations ──────────────────────────────────────────────────────────
  { href: "/system/health",        category: "Operations", title: "System Health",   description: "Real-time infrastructure status — Postgres, Redis, API, jobs, email.",    icon: <Activity className="h-5 w-5" /> },
  { href: "/system/deleted-accounts", category: "Security & Access", title: "Deleted accounts",  description: "Accounts somebody closed. Restore one, or remove it for good.", icon: <UserCheck className="h-5 w-5" /> },
  { href: "/system/trash",         category: "Data & Files", title: "Trash",                description: "Deleted records, restorable for 30 days before they go for good.", icon: <Trash2 className="h-5 w-5" /> },
  { href: "/system/performance",   category: "Operations", title: "Operations",      description: "Latency, traffic, errors, saturation, slow routes and scaling readiness.", icon: <TrendingUp className="h-5 w-5" /> },
  { href: "/system/jobs",          category: "Operations", title: "Background Jobs",  description: "Queue depth, in-flight workers, dead-letter queue.",                      icon: <Database className="h-5 w-5" /> , module: "jobs" },
  { href: "/system/sagas",         category: "Operations", title: "Sagas",           description: "Multi-step processes, and the ones whose undo could not be completed.", icon: <Workflow className="h-5 w-5" /> },
  { href: "/system/routes",        category: "Operations", title: "Routes",          description: "Every endpoint this API serves, searchable, with a curl for each.",     icon: <Server className="h-5 w-5" /> },
  { href: "/system/cron",          category: "Operations", title: "Cron Schedules",  description: "Recurring jobs, next-run times, run history.",                            icon: <Calendar className="h-5 w-5" /> , module: "cron" },
  { href: "/system/activity",      category: "Operations", title: "User Activity",   description: "Auth events, writes, operator actions with IP + severity.",               icon: <Activity className="h-5 w-5" /> , module: "audit" },
  // ── Security & Access ───────────────────────────────────────────────────
  { href: "/account",               category: "Security & Access", title: "Account",             description: "Your password, two-factor, passkeys and the devices you are signed in on.",  icon: <UserCheck className="h-5 w-5" /> },
  { href: "/system/security",       category: "Security & Access", title: "Security",            description: "Sentinel summary — banned IPs, rate-limit pressure, recent threats.",       icon: <Shield className="h-5 w-5" /> },
  { href: "/system/roles",          category: "Security & Access", title: "Roles & permissions", description: "Define what each role can see and do, and assign roles to users.",           icon: <ShieldCheck className="h-5 w-5" /> },
  { href: "/system/access-reviews", category: "Security & Access", title: "Access Reviews",      description: "Point-in-time snapshots of who holds which role — for audits and attestations.", icon: <UserCheck className="h-5 w-5" /> },
  { href: "/system/gdpr",           category: "Security & Access", title: "GDPR",                description: "Export or erase a user's data, with a tamper-evident deletion journal.",     icon: <Shield className="h-5 w-5" /> },
  { href: "/system/api-keys",       category: "Security & Access", title: "API keys",            description: "Long-lived credentials for scripts and integrations, instead of a login.", icon: <KeyRound className="h-5 w-5" /> },
  { href: "/system/audit",          category: "Security & Access", title: "Audit log",           description: "Every authenticated write, hash-chained so a changed row can be detected.", icon: <ShieldCheck className="h-5 w-5" /> },
  { href: "/system/sso",            category: "Security & Access", title: "Single sign-on",      description: "Let a customer's team sign in with their own identity provider — OIDC or SAML.",     icon: <ShieldCheck className="h-5 w-5" /> },
  // ── Data & Files ────────────────────────────────────────────────────────
  { href: "/system/backups",     category: "Data & Files", title: "Data & Backup",       description: "Create, schedule and restore database backups.",                                     icon: <Database className="h-5 w-5" /> , module: "backup" },
  { href: "/system/files",       category: "Data & Files", title: "File Storage",        description: "Browse uploads, manage retention, audit usage.",                                      icon: <Upload className="h-5 w-5" /> , module: "files" },
  { href: "/system/form-shares", category: "Data & Files", title: "Public form sharing", description: "Token-gated public submission links. Generate, enable/disable, view submissions.",     icon: <LinkIcon className="h-5 w-5" /> },
  // ── Communication ───────────────────────────────────────────────────────
  { href: "/system/mail",          category: "Communication", title: "Mail Preview",   description: "Every email template, rendered by the API.",       icon: <Mail className="h-5 w-5" /> , module: "mail" },
  { href: "/system/support",       category: "Communication", title: "Support",        description: "Incoming tickets, threads, assignments, closures.", icon: <MessageSquare className="h-5 w-5" /> },
  { href: "/system/notifications", category: "Communication", title: "Notifications",  description: "Recent system + Sentinel + Pulse notifications.",   icon: <Bell className="h-5 w-5" /> },
  // ── Settings ────────────────────────────────────────────────────────────
  { href: "/settings/dashboard",   category: "Settings", title: "Dashboard settings", description: "Configure the admin dashboard — widgets, stats, and layout.", icon: <Settings className="h-5 w-5" /> },
];

function Tile({ href, icon, title, description }: { href: string; icon: React.ReactNode; title: string; description: string }) {
  return (
    <Link
      href={href}
      className="group rounded-xl border border-border bg-bg-elevated p-5 transition-colors hover:bg-bg-hover hover:border-accent/30"
    >
      <div className="mb-3 inline-flex h-10 w-10 items-center justify-center rounded-lg bg-accent/10 text-accent">
        {icon}
      </div>
      <p className="text-base font-semibold text-foreground group-hover:text-accent">{title}</p>
      <p className="mt-1 text-sm text-text-secondary">{description}</p>
    </Link>
  );
}

export default function SystemHubPage() {
  // Hide tiles for modules that are switched off — a tile linking to a route
  // the server no longer mounts just 404s.
  const { moduleEnabled } = useModules();
  const tiles = TILES.filter((t) => !t.module || moduleEnabled(t.module));

  // Anything a plugin injected into SYSTEM_NAV that isn't one of our built-in
  // tiles (matched by href) becomes an "Extensions" entry.
  const builtinHrefs = new Set(TILES.map((t) => t.href));
  const extensions = SYSTEM_NAV.filter(
    (n) => !builtinHrefs.has(n.href) && n.href !== "/system" && (!n.module || moduleEnabled(n.module))
  );

  const tabs: string[] = [...CATEGORIES, ...(extensions.length ? [EXTENSIONS] : [])];
  const [active, setActive] = useState<string>(CATEGORIES[0]);
  const current = tabs.includes(active) ? active : CATEGORIES[0];

  return (
    <div>
      <PageHeader
        title="System"
        subtitle="Every operational surface for this app, grouped. Pick a tab, then a tile."
      />

      {/* Category tabs */}
      <div className="mb-6 flex flex-wrap gap-1 border-b border-border">
        {tabs.map((t) => (
          <button
            key={t}
            type="button"
            onClick={() => setActive(t)}
            className={
              "relative px-4 py-2.5 text-sm font-medium transition-colors " +
              (current === t ? "text-accent" : "text-text-secondary hover:text-foreground")
            }
          >
            {t}
            {current === t && (
              <span className="absolute inset-x-2 -bottom-px h-0.5 rounded-full bg-accent" />
            )}
          </button>
        ))}
      </div>

      <div className="grid grid-cols-1 gap-4 md:grid-cols-2 xl:grid-cols-3">
        {current === EXTENSIONS
          ? extensions.map((n) => (
              <Tile
                key={n.href}
                href={n.href}
                icon={INTERNAL_ICON[n.iconKey] ?? <LayoutGrid className="h-5 w-5" />}
                title={n.label}
                description="Added by a Grit plugin."
              />
            ))
          : tiles
              .filter((t) => t.category === current)
              .map((t) => (
                <Tile key={t.href} href={t.href} icon={t.icon} title={t.title} description={t.description} />
              ))}
      </div>

      <FileText className="hidden" />
    </div>
  );
}
