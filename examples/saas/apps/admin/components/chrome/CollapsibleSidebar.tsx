"use client";

import { useEffect, useRef, useState } from "react";
import Link from "next/link";
import { usePathname } from "next/navigation";
import { resources } from "@/resources";
import { brand } from "@repo/shared/brand";
import { useLogout } from "@/hooks/use-auth";
import { usePermissions } from "@/hooks/use-permissions";
import { useT } from "@/lib/i18n";
import {
  getIcon,
  ChevronDown,
  ChevronLeft,
  ChevronRight,
  LayoutDashboard,
  LayoutGrid,
  Activity,
  MessageSquare,
  Bell,
  Settings,
  TrendingUp,
  Shield,
  ShieldCheck,
  HardDrive,
  UserCircle,
  UserCheck,
  Webhook,
  User as UserIcon,
  LogOut,
} from "@/lib/icons";
import type { User } from "@repo/shared/types";

// GRIT_CLI_VERSION is the scaffold version that generated this file.
// Surfaced in the sidebar footer so the user can quickly see what Grit
// release their dashboard was built from. Update by re-scaffolding or
// hand-bumping if you carry framework patches locally.
const GRIT_CLI_VERSION = "v3.392.0";

// Internal nav block — pages that exist for every Grit app regardless of
// which resources were generated. Kept out of the resources registry so
// developers don't accidentally remove them when editing resources.ts.
// NavEntry describes a plugin-injectable System link. As of v3.102.0 the
// built-in operational surfaces (Health, Performance, Security, Roles, …) no
// longer live in the rail — they're grouped into tabs inside the System Hub
// (/system). SYSTEM_NAV therefore holds ONLY what Grit plugins inject
// (Webhooks, Impersonate, …); the hub reads it and lists those under an
// "Extensions" tab. Exported so the hub can import it.
export type NavEntry = {
  href: string;
  label: string;
  iconKey: string;
  adminOnly: boolean;
  /** Permission required to see this item, e.g. "roles.view". Optional. */
  requires?: string;
  /** Optional module this item belongs to; hidden when that module is off. */
  module?: string;
};

// Plugins inject their System link here (see internal/plugin/*.go — they anchor
// on the injection marker below). Empty by default. That marker MUST stay
// inside this array literal, or "grit plugin add" can't wire its nav link.
export const SYSTEM_NAV: readonly NavEntry[] = [
  // grit:nav:system
] as const;

// The system pages that earn a place in the rail.
//
// Four, not thirty. The test is whether somebody opens it while working rather
// than while configuring: Operations when something feels slow, Security when
// something looks wrong, Backups before anything irreversible, and Account
// because it is where your own password and sessions are. Everything else is
// one click further on, in the hub.
const SYSTEM_GROUP = "__system";

const SYSTEM_LINKS: ReadonlyArray<{
  href: string;
  label: string;
  i18nKey: string;
  icon: React.ReactNode;
}> = [
  { href: "/system/performance", label: "Operations", i18nKey: "nav.operations", icon: <TrendingUp className="h-5 w-5" /> },
  { href: "/system/security", label: "Security", i18nKey: "nav.security", icon: <ShieldCheck className="h-5 w-5" /> },
  { href: "/system/backups", label: "Backups", i18nKey: "nav.backups", icon: <HardDrive className="h-5 w-5" /> },
  { href: "/account", label: "Account", i18nKey: "nav.account", icon: <UserCircle className="h-5 w-5" /> },
];

// Icon lookup for plugin-injected SYSTEM_NAV entries, keyed by iconKey.
// Exported so the System Hub's Extensions tab renders the same icons. Keep
// every entry even if the default (empty) nav doesn't use it: plugins
// reference these keys, and listing them keeps the icon imports live.
export const INTERNAL_ICON: Record<string, React.ReactNode> = {
  Activity: <Activity className="h-5 w-5" />,
  ActivityIcon: <Activity className="h-5 w-5" />,
  MessageSquare: <MessageSquare className="h-5 w-5" />,
  Bell: <Bell className="h-5 w-5" />,
  Settings: <Settings className="h-5 w-5" />,
  TrendingUp: <TrendingUp className="h-5 w-5" />,
  Shield: <Shield className="h-5 w-5" />,
  ShieldCheck: <ShieldCheck className="h-5 w-5" />,
  UserCheck: <UserCheck className="h-5 w-5" />,
  LayoutGrid: <LayoutGrid className="h-5 w-5" />,
  Webhook: <Webhook className="h-5 w-5" />,
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
  // A group is open when the reader said so, and otherwise when it holds the
  // page being shown. Worked out during render: copying it into state from an
  // effect rendered the sidebar twice per navigation, and compared the path
  // with "/" + slug, which no resource page matches, so a group never opened.
  const [manualGroups, setManualGroups] = useState<Record<string, boolean>>({});
  const activeGroup = resources.find(
    (r) => r.group && pathname.startsWith("/resources/" + r.slug)
  )?.group;
  // A resource group opens when you are inside it. The system group opens by
  // default, because the four links in it were promoted out of the hub to save
  // a click and a collapsed group gives that click straight back.
  const isGroupOpen = (key: string) =>
    manualGroups[key] ?? (key === SYSTEM_GROUP || key === activeGroup);

  const isAdmin = user.role === "ADMIN" || user.role === "EDITOR";
  const toggle = (key: string) => setManualGroups((p) => ({ ...p, [key]: !isGroupOpen(key) }));

  const { can, isLoading: permsLoading } = usePermissions();
  const t = useT();

  // The System Hub fronts operations pages that each need a permission. Shown
  // to an ADMIN, or to anyone holding one of them: by role name alone it put an
  // EDITOR in front of pages the API then refused.
  const showSystem =
    user.role === "ADMIN" ||
    (!permsLoading &&
      ["system.view", "jobs.view", "backups.view", "audit.view", "roles.view"].some((p) => can(p)));

  // Permission gating, on top of the adminOnly check.
  //
  // Only applied once permissions have loaded: can() fails closed while the
  // request is in flight, so filtering during load would briefly blank items an
  // admin is entitled to. This is a UX layer either way — every route is
  // enforced server-side.
  //
  // Every generated resource registers a "<slug>.view" permission (feature key
  // == slug), so a user who was granted, say, only Categories and Products sees
  // exactly those two in the nav — not the whole catalog. Super-admins short-
  // circuit inside can(), so they still see everything.
  const visibleResources = resources.filter(
    (r) => !r.hidden && (!r.adminOnly || isAdmin) && (permsLoading || can(r.slug + ".view"))
  );

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
            label={t("nav.dashboard", "Dashboard")}
            active={pathname === "/dashboard"}
            collapsed={collapsed}
            onClick={onMobileClose}
          />

          {/* The resources, under a heading of their own so the system group
              below reads as a different kind of thing rather than more of the
              same list. Suppressed when there are none, because a heading over
              nothing is worse than no heading. */}
          {!collapsed && groups._root.length > 0 && (
            <p className="-mx-2 mt-3 border-t border-border px-4 pb-0.5 pt-3 text-xs font-semibold uppercase tracking-wider text-text-muted">
              {t("nav.content", "Content")}
            </p>
          )}

          {groups._root.map((r) => {
            const Icon = getIcon(r.icon);
            return (
              <SidebarLink
                key={r.slug}
                href={"/resources/" + r.slug}
                icon={<Icon className="h-5 w-5" />}
                label={t("resources." + r.slug + ".plural", r.label?.plural ?? r.name)}
                active={pathname.startsWith("/resources/" + r.slug)}
                collapsed={collapsed}
                onClick={onMobileClose}
              />
            );
          })}

          {Object.entries(groups)
            .filter(([k]) => k !== "_root")
            .map(([groupName, items]) => (
              <div key={groupName} className={collapsed ? "pt-1" : "-mx-2 mt-3 border-t border-border px-2 pt-3"}>
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
                        (isGroupOpen(groupName) ? "rotate-0" : "-rotate-90")
                      }
                    />
                  </button>
                )}
                {(collapsed || isGroupOpen(groupName)) && items.map((r) => {
                  const Icon = getIcon(r.icon);
                  return (
                    <SidebarLink
                      key={r.slug}
                      href={"/resources/" + r.slug}
                      icon={<Icon className="h-5 w-5" />}
                      label={t("resources." + r.slug + ".plural", r.label?.plural ?? r.name)}
                      active={pathname.startsWith("/resources/" + r.slug)}
                      collapsed={collapsed}
                      onClick={onMobileClose}
                    />
                  );
                })}
              </div>
            ))}

          {/* System Hub — a single entry point for every operational surface
              (Health, Performance, Security, Roles, GDPR, Access Reviews, plus
              anything plugins add). They're grouped into tabs inside the hub at
              /system rather than crowding the rail. Admin-only, matching the
              surfaces it fronts. Highlights for /system/* and /settings/*. */}
          {/* System — the four surfaces that get opened daily, then the hub
              that fronts the other thirty. Promoting these is only affordable
              because the group header keeps them from reading as a flat list:
              before this, Operations was two clicks from anywhere. */}
          {showSystem && (
            <div className={collapsed ? "pt-2" : "-mx-2 mt-3 border-t border-border px-2 pt-3"}>
              {!collapsed && (
                <button
                  type="button"
                  onClick={() => toggle(SYSTEM_GROUP)}
                  aria-expanded={isGroupOpen(SYSTEM_GROUP)}
                  className="flex w-full items-center justify-between rounded px-2 py-1.5 text-xs font-semibold uppercase tracking-wider text-text-muted hover:text-text-secondary"
                >
                  <span>{t("nav.system", "System")}</span>
                  <ChevronDown
                    className={
                      "h-3.5 w-3.5 transition-transform " +
                      (isGroupOpen(SYSTEM_GROUP) ? "rotate-0" : "-rotate-90")
                    }
                  />
                </button>
              )}
              {(collapsed || isGroupOpen(SYSTEM_GROUP)) && (
                <>
                  {SYSTEM_LINKS.map((link) => (
                    <SidebarLink
                      key={link.href}
                      href={link.href}
                      icon={link.icon}
                      label={t(link.i18nKey, link.label)}
                      active={pathname.startsWith(link.href)}
                      collapsed={collapsed}
                      onClick={onMobileClose}
                    />
                  ))}
                  <SidebarLink
                    href="/system"
                    icon={<LayoutGrid className="h-5 w-5" />}
                    label={t("nav.systemHub", "System Hub")}
                    active={pathname === "/system" || pathname.startsWith("/settings")}
                    collapsed={collapsed}
                    onClick={onMobileClose}
                  />
                </>
              )}
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
            <img src={user.avatar} alt={fullName} width={36} height={36} decoding="async" className="h-full w-full object-cover" />
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
              href="/account"
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
        height={32}
        decoding="async"
        className={collapsed ? "h-8 w-8 object-contain" : "h-8 w-auto object-contain"}
      />
    );
  }

  return (
    <span className="inline-flex h-8 w-8 items-center justify-center rounded-lg bg-accent text-accent-fg font-bold text-sm">
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
