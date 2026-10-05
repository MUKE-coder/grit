"use client";

import { useQueryClient } from "@tanstack/react-query";
import type { ReactNode } from "react";
import Link from "next/link";
import { usePathname } from "next/navigation";
import { RefreshCw, Search, ArrowLeft } from "@/lib/icons";
import { DarkModeToggle } from "./DarkModeToggle";
import { UserMenu } from "./UserMenu";
import { NotificationBell } from "./NotificationBell";
import { StatCards, type StatCard } from "./StatCards";
import { inputClasses } from "@/components/ui/input";

// Queries that belong to the layout rather than to any page: the signed-in
// user, their permissions and the notification bell. The refresh button leaves
// them alone; they have their own lifetimes.
const CHROME_QUERY_ROOTS = ["me", "my-permissions", "notifications"];

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
  /** Override the auto-derived back link. Pass null to suppress it entirely. */
  backHref?: string | null;
  /** Label for the back link. Defaults to "Back to System Hub". */
  backLabel?: string;
  /** Stat cards shown under the header, e.g. a resource page's totals. */
  stats?: StatCard[];
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
  backHref,
  backLabel,
  stats,
}: PageHeaderProps) {
  const queryClient = useQueryClient();
  const pathname = usePathname();

  // Every operational surface lives under the System Hub, but each page is
  // its own route with no shared layout — so a sub-page would otherwise be a
  // dead end. Derive the back link from the path: any /system/* or
  // /settings/* page (but not the hub itself) gets "Back to System Hub"
  // automatically, including pages added later by plugins. Pass backHref to
  // point somewhere else, or backHref={null} to suppress it.
  const autoBack =
    pathname !== "/system" && (pathname.startsWith("/system/") || pathname.startsWith("/settings"));
  const backTo =
    backHref === null
      ? null
      : backHref
        ? { href: backHref, label: backLabel ?? "Back" }
        : autoBack
          ? { href: "/system", label: backLabel ?? "Back to System Hub" }
          : null;

  // Refresh reloads what the page shows, not the chrome around it. With no
  // filter it also refetched the signed-in user, their permissions and the
  // notification list, three requests on every click that the page never
  // asked for, and a permissions refetch briefly re-gates the sidebar.
  // Pages with hot keys (jobs, files, sentinel) can scope by passing refreshKeys.
  const onRefresh = () => {
    if (refreshKeys && refreshKeys.length > 0) {
      for (const k of refreshKeys) queryClient.invalidateQueries({ queryKey: [k] });
    } else {
      queryClient.invalidateQueries({
        predicate: (query) => !CHROME_QUERY_ROOTS.includes(String(query.queryKey[0])),
      });
    }
  };

  return (
    <>
    {/* v3.31.6: PageHeader is now sticky, pinned to the top of the
        scrollable main area with a solid background + bottom border so
        long page content scrolls behind it. -mx-4 md:-mx-8 cancels the
        main's px-* padding so the bg + border stretch to the edges, and
        px-* inside brings the content back inside the original gutter. */}
    <header className="sticky top-0 z-20 -mx-4 mb-6 border-b border-border bg-background/90 backdrop-blur supports-[backdrop-filter]:bg-background/75 md:-mx-8">
      <div className="flex flex-col gap-4 px-4 py-4 md:flex-row md:items-center md:justify-between md:px-8">
        {/* Title block — min-w-0 + flex-shrink lets the title wrap
            cleanly when long subtitles share the row with action chrome. */}
        <div className="min-w-0 md:flex-1">
          {backTo && (
            <Link
              href={backTo.href}
              className="mb-1.5 inline-flex items-center gap-1.5 text-xs font-medium text-text-muted transition-colors hover:text-accent"
            >
              <ArrowLeft className="h-3.5 w-3.5" />
              {backTo.label}
            </Link>
          )}
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
              className={inputClasses({ className: "pl-9 pr-3" })}
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
    {/* The cards are content, not chrome: they scroll away under the header. */}
    {stats && stats.length > 0 && <StatCards stats={stats} />}
    </>
  );
}
