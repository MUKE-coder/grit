package scaffold

// v3.29 chrome components for the admin dashboard.
//
//   components/chrome/DarkModeToggle.tsx — localStorage-backed two-mode toggle
//   components/chrome/PageHeader.tsx     — title/subtitle/search/actions
//   components/chrome/UserMenu.tsx       — avatar + dropdown (activity/settings/billing/logout)
//   components/chrome/NotificationBell.tsx — bell + unread badge stub (full system in
//                                            the notifications bundle)
//   components/chrome/CollapsibleSidebar.tsx — top-mounted collapse + 2 logos
//
// These render inside the (dashboard) layout. Pages import PageHeader at
// the top of their JSX to get a consistent header strip with chrome
// affordances pre-wired.

// adminDarkModeToggleComponent generates a self-contained light/dark toggle
// that flips data-theme-mode on <html> (independent of the theme name) and
// persists choice to localStorage. The CSS rules under [data-theme-mode]
// already live in globals.css; the toggle is just the input device.
func adminDarkModeToggleComponent() string {
	return `"use client";

import { useEffect, useState } from "react";
import { Moon, Sun } from "@/lib/icons";

/**
 * Two-mode light/dark toggle. Persists to localStorage("grit-theme-mode").
 * The chosen mode is written to <html data-theme-mode="light|dark">; the
 * data-theme attribute set by the root layout stays untouched so a user
 * can be on the Atlas palette in dark mode without losing brand colors.
 *
 * Server-render returns null to avoid a hydration mismatch; the body
 * paints once we know what's in localStorage.
 */
export function DarkModeToggle({ className = "" }: { className?: string }) {
  const [mode, setMode] = useState<"light" | "dark" | null>(null);

  useEffect(() => {
    const stored = (typeof window !== "undefined"
      ? window.localStorage.getItem("grit-theme-mode")
      : null) as "light" | "dark" | null;
    const initial: "light" | "dark" = stored || "light";
    setMode(initial);
    document.documentElement.setAttribute("data-theme-mode", initial);
  }, []);

  const flip = () => {
    const next: "light" | "dark" = mode === "dark" ? "light" : "dark";
    setMode(next);
    document.documentElement.setAttribute("data-theme-mode", next);
    try {
      window.localStorage.setItem("grit-theme-mode", next);
    } catch {
      // private browsing or quota — non-fatal, the in-memory mode still flips.
    }
  };

  if (mode === null) return null;

  return (
    <button
      type="button"
      onClick={flip}
      aria-label={mode === "dark" ? "Switch to light mode" : "Switch to dark mode"}
      className={"inline-flex h-9 w-9 items-center justify-center rounded-lg border border-border bg-bg-elevated text-text-secondary hover:bg-bg-hover transition-colors " + className}
    >
      {mode === "dark" ? <Sun className="h-4 w-4" /> : <Moon className="h-4 w-4" />}
    </button>
  );
}
`
}

// adminUserMenuComponent — avatar + dropdown with name, email, activity,
// settings, billing, logout. Built off the existing useMe + useLogout
// hooks from v3.27 so it slots into both new and upgraded apps.
func adminUserMenuComponent() string {
	return `"use client";

import Link from "next/link";
import { useEffect, useRef, useState } from "react";
import { useMe, useLogout } from "@/hooks/use-auth";
import { Activity, Settings, CreditCard, LogOut } from "@/lib/icons";

/**
 * Avatar button + dropdown. Click outside to close. Click items to
 * navigate (Link) or perform an action (Logout). The dropdown links
 * default to routes Grit ships — change the hrefs to fit your app.
 */
export function UserMenu() {
  const { data: user } = useMe();
  const { mutate: logout } = useLogout();
  const [open, setOpen] = useState(false);
  const ref = useRef<HTMLDivElement>(null);

  useEffect(() => {
    if (!open) return;
    const onClick = (e: MouseEvent) => {
      if (ref.current && !ref.current.contains(e.target as Node)) setOpen(false);
    };
    document.addEventListener("mousedown", onClick);
    return () => document.removeEventListener("mousedown", onClick);
  }, [open]);

  if (!user) return null;

  const initials = ((user.first_name?.[0] || "") + (user.last_name?.[0] || "")).toUpperCase() || "U";
  const fullName = [user.first_name, user.last_name].filter(Boolean).join(" ") || "User";

  return (
    <div ref={ref} className="relative">
      <button
        type="button"
        onClick={() => setOpen((v) => !v)}
        aria-label="Open user menu"
        className="h-9 w-9 overflow-hidden rounded-full border border-border bg-bg-elevated"
      >
        {user.avatar ? (
          <img src={user.avatar} alt={fullName} className="h-full w-full object-cover" />
        ) : (
          <span className="flex h-full w-full items-center justify-center text-sm font-semibold text-text-secondary">
            {initials}
          </span>
        )}
      </button>

      {open && (
        <div className="absolute right-0 z-50 mt-2 w-64 overflow-hidden rounded-xl border border-border bg-bg-elevated shadow-xl">
          <div className="border-b border-border px-4 py-3">
            <p className="text-sm font-semibold text-foreground truncate">{fullName}</p>
            <p className="text-xs text-text-muted truncate">{user.email}</p>
          </div>
          <nav className="py-1 text-sm">
            <Link
              href="/system/activity"
              onClick={() => setOpen(false)}
              className="flex items-center gap-3 px-4 py-2.5 text-text-secondary hover:bg-bg-hover hover:text-foreground"
            >
              <Activity className="h-4 w-4" />
              User Activity
            </Link>
            <Link
              href="/profile"
              onClick={() => setOpen(false)}
              className="flex items-center gap-3 px-4 py-2.5 text-text-secondary hover:bg-bg-hover hover:text-foreground"
            >
              <Settings className="h-4 w-4" />
              Settings
            </Link>
            <Link
              href="/system/billing"
              onClick={() => setOpen(false)}
              className="flex items-center justify-between px-4 py-2.5 text-text-secondary hover:bg-bg-hover hover:text-foreground"
            >
              <span className="flex items-center gap-3">
                <CreditCard className="h-4 w-4" />
                Billing
              </span>
              {user.role === "ADMIN" && (
                <span className="rounded bg-accent/10 px-1.5 py-0.5 text-[10px] font-semibold uppercase text-accent">
                  Admin
                </span>
              )}
            </Link>
            <button
              type="button"
              onClick={() => { setOpen(false); logout(); }}
              className="flex w-full items-center gap-3 px-4 py-2.5 text-left text-text-secondary hover:bg-bg-hover hover:text-foreground"
            >
              <LogOut className="h-4 w-4" />
              Log out
            </button>
          </nav>
        </div>
      )}
    </div>
  );
}
`
}

// adminNotificationBellComponent — full bell + dropdown. Polls the
// existing GET /api/notifications (returns {data, unread}) every 60s so
// the count and the list stay aligned with one request. Opens a panel
// with the most recent notifications + mark-read affordances.
func adminNotificationBellComponent() string {
	return `"use client";

import Link from "next/link";
import { useEffect, useRef, useState } from "react";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { Bell, Check, AlertCircle, AlertTriangle, Activity } from "@/lib/icons";
import { apiClient } from "@/lib/api-client";

interface Notification {
  id: string;
  source: "sentinel" | "pulse" | "system";
  severity: "critical" | "high" | "medium" | "low" | "info";
  title: string;
  body: string;
  link: string;
  count: number;
  read_at: string | null;
  created_at: string;
}

interface ListResponse {
  data: Notification[];
  unread: number;
}

const severityColor: Record<Notification["severity"], string> = {
  critical: "text-danger",
  high: "text-warning",
  medium: "text-info",
  low: "text-text-secondary",
  info: "text-text-secondary",
};

const sourceIcon: Record<Notification["source"], typeof AlertCircle> = {
  sentinel: AlertTriangle,
  pulse: Activity,
  system: AlertCircle,
};

/**
 * Bell + dropdown. The bell badge polls /api/notifications every 60s for
 * fresh data. Opening the dropdown reveals the most recent items with a
 * mark-read button per row + a mark-all-read footer. Clicking a row
 * navigates to the linked surface (Sentinel finding, Pulse trace, etc.).
 */
export function NotificationBell() {
  const [open, setOpen] = useState(false);
  const ref = useRef<HTMLDivElement>(null);
  const queryClient = useQueryClient();

  const { data } = useQuery<ListResponse>({
    queryKey: ["notifications", "list"],
    queryFn: async () => {
      try {
        const { data } = await apiClient.get<ListResponse>("/api/notifications");
        return data;
      } catch {
        return { data: [], unread: 0 };
      }
    },
    refetchInterval: 60_000,
    retry: false,
    staleTime: 30_000,
  });

  const markRead = useMutation({
    mutationFn: async (id: string) => apiClient.post("/api/notifications/" + id + "/read"),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ["notifications", "list"] }),
  });

  const markAllRead = useMutation({
    mutationFn: async () => apiClient.post("/api/notifications/read-all"),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ["notifications", "list"] }),
  });

  useEffect(() => {
    if (!open) return;
    const onClick = (e: MouseEvent) => {
      if (ref.current && !ref.current.contains(e.target as Node)) setOpen(false);
    };
    document.addEventListener("mousedown", onClick);
    return () => document.removeEventListener("mousedown", onClick);
  }, [open]);

  const unread = data?.unread || 0;
  const items = data?.data || [];

  return (
    <div ref={ref} className="relative">
      <button
        type="button"
        onClick={() => setOpen((v) => !v)}
        aria-label={"Notifications" + (unread > 0 ? " (" + unread + " unread)" : "")}
        className="relative inline-flex h-9 w-9 items-center justify-center rounded-lg border border-border bg-bg-elevated text-text-secondary hover:bg-bg-hover transition-colors"
      >
        <Bell className="h-4 w-4" />
        {unread > 0 && (
          <span className="absolute -top-1 -right-1 flex h-4 min-w-4 items-center justify-center rounded-full bg-danger px-1 text-[10px] font-semibold text-white">
            {unread > 99 ? "99+" : unread}
          </span>
        )}
      </button>

      {open && (
        <div className="absolute right-0 z-50 mt-2 w-96 max-w-[calc(100vw-2rem)] overflow-hidden rounded-xl border border-border bg-bg-elevated shadow-xl">
          <header className="flex items-center justify-between border-b border-border px-4 py-3">
            <h3 className="text-sm font-semibold text-foreground">Notifications</h3>
            {unread > 0 && (
              <button
                type="button"
                onClick={() => markAllRead.mutate()}
                disabled={markAllRead.isPending}
                className="text-xs font-medium text-accent hover:text-accent-hover disabled:opacity-50"
              >
                Mark all read
              </button>
            )}
          </header>

          <div className="max-h-96 overflow-y-auto">
            {items.length === 0 ? (
              <p className="px-4 py-8 text-center text-sm text-text-muted">
                You're all caught up
              </p>
            ) : (
              <ul className="divide-y divide-border">
                {items.map((n) => {
                  const Icon = sourceIcon[n.source] || AlertCircle;
                  const unreadRow = !n.read_at;
                  return (
                    <li key={n.id} className={unreadRow ? "bg-bg-secondary/50" : ""}>
                      <div className="flex gap-3 px-4 py-3">
                        <Icon className={"mt-0.5 h-4 w-4 shrink-0 " + severityColor[n.severity]} />
                        <div className="min-w-0 flex-1">
                          <Link
                            href={n.link || "#"}
                            onClick={() => { setOpen(false); if (unreadRow) markRead.mutate(n.id); }}
                            className="block"
                          >
                            <p className="text-sm font-medium text-foreground truncate">
                              {n.title}
                              {n.count > 1 && (
                                <span className="ml-1 text-xs text-text-muted">×{n.count}</span>
                              )}
                            </p>
                            {n.body && (
                              <p className="mt-0.5 text-xs text-text-secondary line-clamp-2">{n.body}</p>
                            )}
                            <p className="mt-1 text-[10px] uppercase tracking-wide text-text-muted">
                              {n.source} · {new Date(n.created_at).toLocaleString()}
                            </p>
                          </Link>
                        </div>
                        {unreadRow && (
                          <button
                            type="button"
                            onClick={() => markRead.mutate(n.id)}
                            aria-label="Mark read"
                            className="shrink-0 rounded p-1 text-text-muted hover:bg-bg-hover hover:text-foreground"
                          >
                            <Check className="h-3.5 w-3.5" />
                          </button>
                        )}
                      </div>
                    </li>
                  );
                })}
              </ul>
            )}
          </div>

          <footer className="border-t border-border px-4 py-2 text-center">
            <Link
              href="/system/notifications"
              onClick={() => setOpen(false)}
              className="text-xs font-medium text-accent hover:text-accent-hover"
            >
              View all notifications
            </Link>
          </footer>
        </div>
      )}
    </div>
  );
}
`
}

// adminPageHeaderComponent — every dashboard page imports <PageHeader> and
// drops it at the top of its JSX. Title + subtitle live on the left; a
// search box centres; refresh / dark / + / bell / user menu sit on the
// right. Custom action buttons can be passed via the actions prop.
func adminPageHeaderComponent() string {
	return `"use client";

import { useQueryClient } from "@tanstack/react-query";
import type { ReactNode } from "react";
import { RefreshCw, Search } from "@/lib/icons";
import { DarkModeToggle } from "./DarkModeToggle";
import { UserMenu } from "./UserMenu";
import { NotificationBell } from "./NotificationBell";

interface PageHeaderProps {
  /** Page title. Required. */
  title: string;
  /** Optional short description shown under the title. */
  subtitle?: string;
  /** When set, renders a search input centred between title and actions. */
  searchPlaceholder?: string;
  /** Search value (controlled). */
  searchValue?: string;
  /** Search change handler. */
  onSearchChange?: (value: string) => void;
  /** Extra action buttons rendered before the always-on chrome (dark, bell, user). */
  actions?: ReactNode;
  /** React Query keys to invalidate when the refresh button is pressed. */
  refreshKeys?: string[];
  /** Hide the refresh button entirely. */
  hideRefresh?: boolean;
}

/**
 * Standard dashboard page header. Sticks the title + subtitle on the
 * left, an optional search input in the middle, and the right-hand chrome
 * (refresh / dark toggle / custom actions / bell / user menu) on the
 * right. Layout collapses cleanly on mobile by stacking and hiding the
 * search label.
 */
export function PageHeader({
  title,
  subtitle,
  searchPlaceholder,
  searchValue,
  onSearchChange,
  actions,
  refreshKeys,
  hideRefresh,
}: PageHeaderProps) {
  const queryClient = useQueryClient();

  // Refresh defaults to invalidating every query on the page. Pages with
  // hot keys (jobs, files, sentinel) can scope by passing refreshKeys.
  const onRefresh = () => {
    if (refreshKeys && refreshKeys.length > 0) {
      refreshKeys.forEach((k) => queryClient.invalidateQueries({ queryKey: [k] }));
    } else {
      queryClient.invalidateQueries();
    }
  };

  return (
    <header className="mb-6 flex flex-col gap-4 md:flex-row md:items-center md:justify-between">
      {/* Title block */}
      <div className="min-w-0 flex-shrink-0">
        <h1 className="text-2xl font-bold text-foreground tracking-tight">{title}</h1>
        {subtitle && <p className="mt-1 text-sm text-text-secondary">{subtitle}</p>}
      </div>

      {/* Search */}
      {searchPlaceholder && (
        <div className="relative w-full md:max-w-md md:flex-1 md:mx-6">
          <Search className="pointer-events-none absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-text-muted" />
          <input
            type="search"
            value={searchValue ?? ""}
            onChange={(e) => onSearchChange?.(e.target.value)}
            placeholder={searchPlaceholder}
            className="w-full rounded-lg border border-border bg-bg-elevated py-2 pl-9 pr-3 text-sm text-foreground placeholder:text-text-muted focus:border-accent focus:outline-none focus:ring-1 focus:ring-accent"
          />
        </div>
      )}

      {/* Chrome actions */}
      <div className="flex items-center justify-end gap-2">
        {!hideRefresh && (
          <button
            type="button"
            onClick={onRefresh}
            aria-label="Refresh"
            className="inline-flex h-9 w-9 items-center justify-center rounded-lg border border-border bg-bg-elevated text-text-secondary hover:bg-bg-hover transition-colors"
          >
            <RefreshCw className="h-4 w-4" />
          </button>
        )}
        <DarkModeToggle />
        {actions}
        <NotificationBell />
        <UserMenu />
      </div>
    </header>
  );
}
`
}

// adminCollapsibleSidebarComponent — refactored sidebar with two-logo
// brand area and a top-mounted collapse chevron. Replaces the old sidebar
// that put the chevron in the navbar middle. Collapsed state shrinks to
// icon-only nav rail.
func adminCollapsibleSidebarComponent() string {
	return `"use client";

import { useEffect, useState } from "react";
import Link from "next/link";
import { usePathname } from "next/navigation";
import { resources } from "@/resources";
import { brand } from "@repo/shared/brand";
import {
  getIcon,
  ChevronDown,
  ChevronLeft,
  ChevronRight,
  LayoutDashboard,
} from "@/lib/icons";
import type { User } from "@repo/shared/types";

interface SidebarProps {
  user: User;
  collapsed: boolean;
  onToggleCollapsed: () => void;
  mobileOpen: boolean;
  onMobileClose: () => void;
}

/**
 * Collapsible left sidebar.
 *
 * Layout:
 *   ┌──────────────────────┐
 *   │  [logo]      [<]     │  ← brand row + top-mounted collapse chevron
 *   ├──────────────────────┤
 *   │  > Dashboard         │
 *   │  > Resource A        │  ← nav items (icon-only when collapsed)
 *   │  ...                 │
 *   └──────────────────────┘
 *
 * Brand row swaps between brand.logo.image (expanded) and brand.logo.mark
 * (collapsed). Both fall back to the brand.logo.text character on a
 * coloured square if no image is provided.
 */
export function CollapsibleSidebar({
  user,
  collapsed,
  onToggleCollapsed,
  mobileOpen,
  onMobileClose,
}: SidebarProps) {
  const pathname = usePathname();
  const [expandedGroups, setExpandedGroups] = useState<Record<string, boolean>>({});

  const isAdmin = user.role === "ADMIN" || user.role === "EDITOR";

  const toggle = (key: string) => setExpandedGroups((p) => ({ ...p, [key]: !p[key] }));

  // Auto-expand any group whose child matches the current route.
  useEffect(() => {
    resources.forEach((r) => {
      if (r.group && pathname.startsWith("/" + r.slug)) {
        setExpandedGroups((p) => ({ ...p, [r.group as string]: true }));
      }
    });
  }, [pathname]);

  // Hide resources flagged as adminOnly when the user isn't ADMIN/EDITOR.
  const visibleResources = resources.filter(
    (r) => !r.adminOnly || isAdmin
  );

  // Group resources by their group key; ungrouped resources go to "_root".
  const groups: Record<string, typeof resources> = { _root: [] };
  for (const r of visibleResources) {
    const key = r.group || "_root";
    if (!groups[key]) groups[key] = [];
    groups[key].push(r);
  }

  return (
    <>
      {/* Mobile backdrop */}
      {mobileOpen && (
        <div
          className="fixed inset-0 z-30 bg-black/50 md:hidden"
          onClick={onMobileClose}
        />
      )}

      <aside
        className={
          "fixed top-0 left-0 z-40 flex h-screen flex-col border-r border-border bg-bg-secondary transition-all duration-200 " +
          (collapsed ? "w-16 " : "w-64 ") +
          (mobileOpen ? "translate-x-0 " : "-translate-x-full md:translate-x-0 ")
        }
      >
        {/* Brand row */}
        <div className="relative flex h-16 items-center border-b border-border px-3">
          <Link href="/dashboard" className="flex flex-1 items-center gap-2 overflow-hidden">
            <BrandMark collapsed={collapsed} />
            {!collapsed && (
              <span className="truncate text-base font-bold text-foreground">
                {brand.name}
              </span>
            )}
          </Link>

          {/* Top-mounted collapse chevron. Floats over the right edge so
              the click target stays accessible whether collapsed or not. */}
          <button
            type="button"
            onClick={onToggleCollapsed}
            aria-label={collapsed ? "Expand sidebar" : "Collapse sidebar"}
            className="absolute -right-3 top-1/2 hidden h-6 w-6 -translate-y-1/2 items-center justify-center rounded-full border border-border bg-bg-elevated text-text-secondary shadow-sm hover:text-foreground md:flex"
          >
            {collapsed ? (
              <ChevronRight className="h-3.5 w-3.5" />
            ) : (
              <ChevronLeft className="h-3.5 w-3.5" />
            )}
          </button>
        </div>

        {/* Nav */}
        <nav className="flex-1 space-y-1 overflow-y-auto px-2 py-3">
          {/* Always-on Dashboard link */}
          <SidebarLink
            href="/dashboard"
            icon={<LayoutDashboard className="h-5 w-5" />}
            label="Dashboard"
            active={pathname === "/dashboard"}
            collapsed={collapsed}
            onClick={onMobileClose}
          />

          {/* Ungrouped resources */}
          {groups._root.map((r) => (
            <SidebarLink
              key={r.slug}
              href={"/" + r.slug}
              icon={getIcon(r.icon, "h-5 w-5")}
              label={r.label}
              active={pathname.startsWith("/" + r.slug)}
              collapsed={collapsed}
              onClick={onMobileClose}
            />
          ))}

          {/* Grouped resources */}
          {Object.entries(groups)
            .filter(([k]) => k !== "_root")
            .map(([groupName, items]) => (
              <div key={groupName} className="pt-1">
                {!collapsed && (
                  <button
                    type="button"
                    onClick={() => toggle(groupName)}
                    className="flex w-full items-center justify-between rounded px-2 py-1.5 text-xs font-semibold uppercase tracking-wider text-text-muted hover:text-text-secondary"
                  >
                    <span>{groupName}</span>
                    <ChevronDown
                      className={
                        "h-3.5 w-3.5 transition-transform " +
                        (expandedGroups[groupName] ? "rotate-0" : "-rotate-90")
                      }
                    />
                  </button>
                )}
                {(collapsed || expandedGroups[groupName]) && items.map((r) => (
                  <SidebarLink
                    key={r.slug}
                    href={"/" + r.slug}
                    icon={getIcon(r.icon, "h-5 w-5")}
                    label={r.label}
                    active={pathname.startsWith("/" + r.slug)}
                    collapsed={collapsed}
                    onClick={onMobileClose}
                  />
                ))}
              </div>
            ))}
        </nav>

        {/* Footer — user role chip when expanded */}
        {!collapsed && (
          <div className="border-t border-border px-3 py-3">
            <div className="rounded-lg border border-accent/20 bg-accent/5 px-3 py-2">
              <p className="text-xs font-semibold text-foreground">{user.role}</p>
              <p className="text-[10px] text-text-muted">
                {user.role === "ADMIN" ? "Comped — not billed" : "Member"}
              </p>
            </div>
          </div>
        )}
      </aside>
    </>
  );
}

function BrandMark({ collapsed }: { collapsed: boolean }) {
  const src = collapsed ? (brand.logo.mark || brand.logo.image) : (brand.logo.image || brand.logo.mark);

  if (src) {
    return (
      <img
        src={src}
        alt={brand.name}
        className={collapsed ? "h-8 w-8 object-contain" : "h-8 w-auto object-contain"}
      />
    );
  }

  return (
    <span className="inline-flex h-8 w-8 items-center justify-center rounded-lg bg-accent text-white font-bold text-sm">
      {brand.logo.text}
    </span>
  );
}

interface LinkProps {
  href: string;
  icon: React.ReactNode;
  label: string;
  active: boolean;
  collapsed: boolean;
  onClick?: () => void;
}

function SidebarLink({ href, icon, label, active, collapsed, onClick }: LinkProps) {
  return (
    <Link
      href={href}
      onClick={onClick}
      aria-label={collapsed ? label : undefined}
      title={collapsed ? label : undefined}
      className={
        "flex items-center gap-3 rounded-lg px-2.5 py-2 text-sm font-medium transition-colors " +
        (active
          ? "bg-accent/10 text-accent"
          : "text-text-secondary hover:bg-bg-hover hover:text-foreground")
      }
    >
      <span className="flex h-5 w-5 shrink-0 items-center justify-center">{icon}</span>
      {!collapsed && <span className="truncate">{label}</span>}
    </Link>
  );
}
`
}

// adminGlobalCSSDarkModeAddon returns the CSS block that the dark-mode
// toggle activates. We append it to the per-theme blocks so toggling
// data-theme-mode="dark" in browsers swaps every variable to its darker
// counterpart while keeping the active theme's hue intent. Lives at the
// end of globals.css (called from adminGlobalCSS via concatenation).
func adminDarkModeCSSAddon() string {
	return `

/* v3.29 — Dark mode toggle. data-theme-mode is set by DarkModeToggle on
 * the <html> root. We invert the surface variables but keep the active
 * theme's hue/contrast intent. Atlas dark = slate; Aurora dark = warm
 * stone; Pulse dark = pure-black. */

[data-theme-mode="dark"] {
  --bg-primary: #0a0a0f;
  --bg-secondary: #111118;
  --bg-tertiary: #1a1a24;
  --bg-elevated: #22222e;
  --bg-hover: #2a2a38;
  --border: #2a2a3a;
  --text-primary: #e8e8f0;
  --text-secondary: #9090a8;
  --text-muted: #606078;
}

[data-theme="atlas"][data-theme-mode="dark"] {
  --bg-primary: #0f172a;
  --bg-secondary: #1e293b;
  --bg-tertiary: #334155;
  --bg-elevated: #1e293b;
  --bg-hover: #334155;
  --border: #334155;
  --text-primary: #f1f5f9;
  --text-secondary: #cbd5e1;
  --text-muted: #94a3b8;
  --accent: #60a5fa;
  --accent-hover: #93c5fd;
}

[data-theme="aurora"][data-theme-mode="dark"] {
  --bg-primary: #1c1917;
  --bg-secondary: #292524;
  --bg-tertiary: #44403c;
  --bg-elevated: #292524;
  --bg-hover: #44403c;
  --border: #44403c;
  --text-primary: #fafaf9;
  --text-secondary: #d6d3d1;
  --text-muted: #a8a29e;
  --accent: #a78bfa;
  --accent-hover: #c4b5fd;
}

[data-theme="pulse"][data-theme-mode="dark"] {
  --bg-primary: #0a0a0a;
  --bg-secondary: #171717;
  --bg-tertiary: #262626;
  --bg-elevated: #171717;
  --bg-hover: #262626;
  --border: #404040;
  --text-primary: #fafafa;
  --text-secondary: #d4d4d4;
  --text-muted: #a3a3a3;
  --accent: #fbbf24;
  --accent-hover: #fcd34d;
}
`
}
