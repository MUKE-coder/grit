package scaffold

import "strings"

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

// adminDarkModeToggleComponent — robust light/dark switcher. v3.31 makes
// the toggle defensive by writing THREE signals on every flip:
//   1. <html data-theme-mode="dark">  — our CSS variable cascade
//   2. <html class="dark">            — Tailwind's darkMode: "class" hook
//   3. <html style="color-scheme: dark"> — browser chrome (scrollbars, native inputs)
// Any one of these is sufficient to repaint the dashboard; carrying all
// three means the toggle still works when one path breaks (e.g. a custom
// CSS file forgets to consume data-theme-mode).
func adminDarkModeToggleComponent() string {
	return `"use client";

import { useEffect, useState } from "react";
import { Moon, Sun } from "@/lib/icons";

type Mode = "light" | "dark";

// applyMode writes every signal a downstream consumer might key off so
// the visual swap happens regardless of which mechanism a stylesheet is
// using. Idempotent — safe to call on every render.
function applyMode(mode: Mode) {
  const root = document.documentElement;
  root.setAttribute("data-theme-mode", mode);
  root.classList.toggle("dark", mode === "dark");
  root.style.colorScheme = mode;
}

/**
 * Two-mode light/dark toggle. Persists to localStorage("grit-theme-mode").
 *
 * The button stays mounted in SSR (initial render returns the light icon)
 * to avoid a layout jump when the client picks up the stored preference.
 * The actual mode is settled in useEffect after hydration; before that
 * point, the button is decorative.
 */
export function DarkModeToggle({ className = "" }: { className?: string }) {
  const [mode, setMode] = useState<Mode>("light");
  const [hydrated, setHydrated] = useState(false);

  useEffect(() => {
    const stored = (typeof window !== "undefined"
      ? (window.localStorage.getItem("grit-theme-mode") as Mode | null)
      : null);
    // Prefer stored choice; fall back to OS-level preference; default light.
    const osDark = typeof window !== "undefined"
      && window.matchMedia
      && window.matchMedia("(prefers-color-scheme: dark)").matches;
    const initial: Mode = stored || (osDark ? "dark" : "light");
    setMode(initial);
    applyMode(initial);
    setHydrated(true);
  }, []);

  const flip = () => {
    const next: Mode = mode === "dark" ? "light" : "dark";
    setMode(next);
    applyMode(next);
    try {
      window.localStorage.setItem("grit-theme-mode", next);
    } catch {
      // Private browsing / storage quota — the in-memory mode still flips,
      // it just won't survive a reload. Non-fatal.
    }
  };

  return (
    <button
      type="button"
      onClick={flip}
      aria-label={mode === "dark" ? "Switch to light mode" : "Switch to dark mode"}
      aria-pressed={mode === "dark"}
      suppressHydrationWarning
      className={"inline-flex h-9 w-9 items-center justify-center rounded-lg border border-border bg-bg-elevated text-text-secondary hover:bg-bg-hover transition-colors " + className}
    >
      {hydrated && mode === "dark" ? <Sun className="h-4 w-4" /> : <Moon className="h-4 w-4" />}
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
import { Activity, User as UserIcon, LogOut } from "@/lib/icons";

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
        className="h-9 w-9 overflow-hidden rounded-full ring-2 ring-accent/40 hover:ring-accent transition-colors bg-bg-elevated"
      >
        {user.avatar ? (
          <img src={user.avatar} alt={fullName} className="h-full w-full object-cover" />
        ) : (
          <span className="flex h-full w-full items-center justify-center text-sm font-semibold text-foreground">
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
              <UserIcon className="h-4 w-4" />
              Profile
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
    // v3.31.6: PageHeader is now sticky — pinned to the top of the
    // scrollable main area with a solid background + bottom border so
    // long page content scrolls behind it. -mx-4 md:-mx-8 cancels the
    // main's px-* padding so the bg + border stretch to the edges, and
    // px-* inside brings the content back inside the original gutter.
    <header className="sticky top-0 z-20 -mx-4 mb-6 border-b border-border bg-bg-primary/90 backdrop-blur supports-[backdrop-filter]:bg-bg-primary/75 md:-mx-8">
      <div className="flex flex-col gap-4 px-4 py-4 md:flex-row md:items-center md:justify-between md:px-8">
        {/* Title block — min-w-0 + flex-shrink lets the title wrap
            cleanly when long subtitles share the row with action chrome. */}
        <div className="min-w-0 md:flex-1">
          <h1 className="text-2xl font-bold text-foreground tracking-tight truncate">{title}</h1>
          {subtitle && <p className="mt-1 text-sm text-text-secondary md:line-clamp-2">{subtitle}</p>}
        </div>

        {/* Search */}
        {searchPlaceholder && (
          <div className="relative w-full md:max-w-xs md:flex-1">
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

        {/* Chrome actions — shrink-0 + whitespace-nowrap on the row so
            action buttons (e.g. "Open full Pulse") don't wrap mid-label. */}
        <div className="flex shrink-0 items-center justify-end gap-2 whitespace-nowrap">
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
func adminCollapsibleSidebarComponent(opts Options) string {
	// {{GRIT_VERSION}} is substituted at scaffold time so the sidebar
	// footer truthfully reports the CLI version that generated the app —
	// helps when an upgrade later changes the scaffold and a user is
	// figuring out which template they're running.
	body := `"use client";

import { useEffect, useRef, useState } from "react";
import Link from "next/link";
import { usePathname } from "next/navigation";
import { resources } from "@/resources";
import { brand } from "@repo/shared/brand";
import { useLogout } from "@/hooks/use-auth";
import {
  getIcon,
  ChevronDown,
  ChevronLeft,
  ChevronRight,
  LayoutDashboard,
  Activity,
  MessageSquare,
  Bell,
  Settings,
  TrendingUp,
  Shield,
  User as UserIcon,
  LogOut,
} from "@/lib/icons";
import type { User } from "@repo/shared/types";

// GRIT_CLI_VERSION is the scaffold version that generated this file.
// Surfaced in the sidebar footer so the user can quickly see what Grit
// release their dashboard was built from. Update by re-scaffolding or
// hand-bumping if you carry framework patches locally.
const GRIT_CLI_VERSION = "v{{GRIT_VERSION}}";

// Internal nav block — pages that exist for every Grit app regardless of
// which resources were generated. Kept out of the resources registry so
// developers don't accidentally remove them when editing resources.ts.
const INTERNAL_NAV = [
  { href: "/system/activity",      label: "Activity",      iconKey: "Activity",       adminOnly: false },
  { href: "/system/support",       label: "Support",       iconKey: "MessageSquare",  adminOnly: false },
  { href: "/system/notifications", label: "Notifications", iconKey: "Bell",           adminOnly: false },
] as const;

// v3.31.5: dedicated SYSTEM section for admin-only operational surfaces.
// Health / Performance / Security live here so they're one click away
// during an incident — the System hub at /system still aggregates every
// surface for the broader browse case.
const SYSTEM_NAV = [
  { href: "/system/health",       label: "System Health", iconKey: "ActivityIcon", adminOnly: true },
  { href: "/system/performance",  label: "Performance",   iconKey: "TrendingUp",   adminOnly: true },
  { href: "/system/security",     label: "Security",      iconKey: "Shield",       adminOnly: true },
  { href: "/system",              label: "System Hub",    iconKey: "Settings",     adminOnly: true },
] as const;

const INTERNAL_ICON: Record<string, React.ReactNode> = {
  Activity: <Activity className="h-5 w-5" />,
  ActivityIcon: <Activity className="h-5 w-5" />,
  MessageSquare: <MessageSquare className="h-5 w-5" />,
  Bell: <Bell className="h-5 w-5" />,
  Settings: <Settings className="h-5 w-5" />,
  TrendingUp: <TrendingUp className="h-5 w-5" />,
  Shield: <Shield className="h-5 w-5" />,
};

interface SidebarProps {
  user: User;
  collapsed: boolean;
  onToggleCollapsed: () => void;
  mobileOpen: boolean;
  onMobileClose: () => void;
}

/**
 * Collapsible left sidebar. Three vertical zones:
 *
 *   ┌──────────────────────┐
 *   │  [logo]      [<]     │  ← brand + collapse chevron (fixed)
 *   ├──────────────────────┤
 *   │  > Dashboard         │
 *   │  > Resource A        │  ← scrollable nav (overflow-y-auto)
 *   │  > Internal section  │
 *   │  ...                 │
 *   ├──────────────────────┤
 *   │  [avatar] Name +     │  ← rich user menu + version (sticky bottom)
 *   │           email      │
 *   └──────────────────────┘
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

  useEffect(() => {
    resources.forEach((r) => {
      if (r.group && pathname.startsWith("/" + r.slug)) {
        setExpandedGroups((p) => ({ ...p, [r.group as string]: true }));
      }
    });
  }, [pathname]);

  const visibleResources = resources.filter((r) => !r.adminOnly || isAdmin);
  const visibleInternal = INTERNAL_NAV.filter((r) => !r.adminOnly || isAdmin);
  const visibleSystem = SYSTEM_NAV.filter((r) => !r.adminOnly || isAdmin);

  const groups: Record<string, typeof resources> = { _root: [] };
  for (const r of visibleResources) {
    const key = r.group || "_root";
    if (!groups[key]) groups[key] = [];
    groups[key].push(r);
  }

  return (
    <>
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
        {/* Brand row — fixed at top */}
        <div className="relative flex h-16 items-center border-b border-border px-3 shrink-0">
          <Link href="/dashboard" className="flex flex-1 items-center gap-2 overflow-hidden">
            <BrandMark collapsed={collapsed} />
            {!collapsed && (
              <span className="truncate text-base font-bold text-foreground">
                {brand.name}
              </span>
            )}
          </Link>

          <button
            type="button"
            onClick={onToggleCollapsed}
            aria-label={collapsed ? "Expand sidebar" : "Collapse sidebar"}
            className="absolute -right-3 top-1/2 hidden h-6 w-6 -translate-y-1/2 items-center justify-center rounded-full border border-border bg-bg-elevated text-text-secondary shadow-sm hover:text-foreground md:flex"
          >
            {collapsed ? <ChevronRight className="h-3.5 w-3.5" /> : <ChevronLeft className="h-3.5 w-3.5" />}
          </button>
        </div>

        {/* Nav — scrollable middle zone. flex-1 + min-h-0 makes it scroll
            instead of pushing the footer off-screen when there are many
            items. Custom scrollbar styles keep the rail subtle. */}
        <nav className="flex-1 min-h-0 space-y-1 overflow-y-auto px-2 py-3 [scrollbar-width:thin] [scrollbar-color:var(--border)_transparent]">
          <SidebarLink
            href="/dashboard"
            icon={<LayoutDashboard className="h-5 w-5" />}
            label="Dashboard"
            active={pathname === "/dashboard"}
            collapsed={collapsed}
            onClick={onMobileClose}
          />

          {groups._root.map((r) => {
            const Icon = getIcon(r.icon);
            return (
              <SidebarLink
                key={r.slug}
                href={"/resources/" + r.slug}
                icon={<Icon className="h-5 w-5" />}
                label={r.label?.plural ?? r.name}
                active={pathname.startsWith("/resources/" + r.slug)}
                collapsed={collapsed}
                onClick={onMobileClose}
              />
            );
          })}

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
                {(collapsed || expandedGroups[groupName]) && items.map((r) => {
                  const Icon = getIcon(r.icon);
                  return (
                    <SidebarLink
                      key={r.slug}
                      href={"/resources/" + r.slug}
                      icon={<Icon className="h-5 w-5" />}
                      label={r.label?.plural ?? r.name}
                      active={pathname.startsWith("/resources/" + r.slug)}
                      collapsed={collapsed}
                      onClick={onMobileClose}
                    />
                  );
                })}
              </div>
            ))}

          {/* Internal section — Activity / Support / Notifications. */}
          <div className="pt-3">
            {!collapsed && (
              <p className="px-2 py-1 text-xs font-semibold uppercase tracking-wider text-text-muted">
                Internal
              </p>
            )}
            {visibleInternal.map((r) => (
              <SidebarLink
                key={r.href}
                href={r.href}
                icon={INTERNAL_ICON[r.iconKey]}
                label={r.label}
                active={pathname.startsWith(r.href)}
                collapsed={collapsed}
                onClick={onMobileClose}
              />
            ))}
          </div>

          {/* System section — admin-only operational surfaces (Health,
              Performance, Security, hub). Active-state uses the full path
              so /system/security highlights without /system also lighting up. */}
          {visibleSystem.length > 0 && (
            <div className="pt-3">
              {!collapsed && (
                <p className="px-2 py-1 text-xs font-semibold uppercase tracking-wider text-text-muted">
                  System
                </p>
              )}
              {visibleSystem.map((r) => (
                <SidebarLink
                  key={r.href}
                  href={r.href}
                  icon={INTERNAL_ICON[r.iconKey]}
                  label={r.label}
                  active={r.href === "/system" ? pathname === "/system" : pathname.startsWith(r.href)}
                  collapsed={collapsed}
                  onClick={onMobileClose}
                />
              ))}
            </div>
          )}
        </nav>

        {/* Footer — sticky bottom: rich user menu + grit version */}
        <SidebarUserMenu user={user} collapsed={collapsed} />
      </aside>
    </>
  );
}

// SidebarUserMenu — sits at the bottom of the sidebar. Click to pop a
// menu with profile / billing / activity / logout. When collapsed shows
// just the avatar + indicator dot. Mirrors the top-right UserMenu but
// renders with the user's name/email visible inline.
function SidebarUserMenu({ user, collapsed }: { user: User; collapsed: boolean }) {
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

  const initials = ((user.first_name?.[0] || "") + (user.last_name?.[0] || "")).toUpperCase() || "U";
  const fullName = [user.first_name, user.last_name].filter(Boolean).join(" ") || "User";

  return (
    <div ref={ref} className="relative border-t border-border shrink-0">
      <button
        type="button"
        onClick={() => setOpen((v) => !v)}
        aria-label="Open account menu"
        className={
          "flex w-full items-center gap-3 px-3 py-3 text-left transition-colors hover:bg-bg-hover " +
          (collapsed ? "justify-center" : "")
        }
      >
        <span
          className={
            "inline-flex h-9 w-9 shrink-0 items-center justify-center rounded-full overflow-hidden ring-2 ring-accent/30 bg-bg-elevated text-sm font-semibold text-foreground"
          }
        >
          {user.avatar ? (
            <img src={user.avatar} alt={fullName} className="h-full w-full object-cover" />
          ) : initials}
        </span>
        {!collapsed && (
          <span className="min-w-0 flex-1">
            <span className="block truncate text-sm font-semibold text-foreground">{fullName}</span>
            <span className="block truncate text-xs text-text-muted">{user.email}</span>
          </span>
        )}
        {!collapsed && (
          <ChevronDown className={"h-3.5 w-3.5 text-text-muted transition-transform " + (open ? "rotate-180" : "")} />
        )}
      </button>

      {open && (
        <div className="absolute bottom-full left-2 right-2 mb-2 overflow-hidden rounded-xl border border-border bg-bg-elevated shadow-xl">
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
              <UserIcon className="h-4 w-4" />
              Profile
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
          {!collapsed && (
            <div className="border-t border-border px-4 py-2 text-center">
              <p className="text-[10px] uppercase tracking-wide text-text-muted">Grit {GRIT_CLI_VERSION}</p>
            </div>
          )}
        </div>
      )}
    </div>
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
	return strings.ReplaceAll(body, "{{GRIT_VERSION}}", opts.Version)
}

// adminDarkModeCSSAddon — dark-mode variable overrides keyed off BOTH
// the .dark class (Tailwind's canonical selector) and data-theme-mode
// (our explicit attribute). Carrying both selectors means the toggle
// still produces a visible repaint even if a stylesheet only listens
// to one of them. v3.31: dropped the per-theme dark refinements that
// were silently winning specificity over the base dark block; the
// flat dark palette is good enough and removes a footgun.
func adminDarkModeCSSAddon() string {
	return `

/* v3.31 — Dark mode. The DarkModeToggle component writes three signals
 * to <html> on every flip: class="dark", data-theme-mode="dark", and
 * style="color-scheme: dark". We honour the two pertinent to CSS here.
 * The active data-theme (atlas / aurora / pulse) keeps its accent hue;
 * only the surface + text variables flip. */

.dark,
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

/* Per-theme accent lifts so blue/purple/black don't disappear on the
 * dark canvas. These come after the base dark block on purpose — the
 * surface tokens above stay; only the brand colours refresh. */

.dark[data-theme="atlas"],
[data-theme-mode="dark"][data-theme="atlas"] {
  --accent: #60a5fa;
  --accent-hover: #93c5fd;
}

.dark[data-theme="aurora"],
[data-theme-mode="dark"][data-theme="aurora"] {
  --accent: #a78bfa;
  --accent-hover: #c4b5fd;
}

.dark[data-theme="pulse"],
[data-theme-mode="dark"][data-theme="pulse"] {
  --accent: #fbbf24;
  --accent-hover: #fcd34d;
}

/* Form controls + selection still need the dark hint or the browser
 * renders them with the OS light scheme regardless of our vars. */
.dark input,
.dark textarea,
.dark select,
[data-theme-mode="dark"] input,
[data-theme-mode="dark"] textarea,
[data-theme-mode="dark"] select {
  color-scheme: dark;
}
`
}
