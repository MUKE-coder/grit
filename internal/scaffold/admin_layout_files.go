package scaffold

// adminThemeScript is the script that applies the stored theme before the
// browser paints anything.
//
// It replaces two things. The first was a ThemeProvider that stored a "grit-theme"
// key nothing else read and whose only consumer was a navbar no layout rendered.
// The second was DarkModeToggle applying the theme in an effect: the toggle sits
// inside the dashboard chrome, which mounts after /auth/me answers, so a dark
// dashboard painted light first on every single load.
//
// Single-quoted so it can be embedded in a double-quoted TS string with no
// escaping, and written without backticks so it survives the Go raw strings it
// travels through. Keys off the same localStorage key DarkModeToggle writes.
const adminThemeScript = `(function(){try{var m=localStorage.getItem('grit-theme-mode');` +
	`if(m!=='dark'&&m!=='light'){m=window.matchMedia&&window.matchMedia('(prefers-color-scheme: dark)').matches?'dark':'light';}` +
	`var r=document.documentElement;r.setAttribute('data-theme-mode',m);` +
	`r.classList.toggle('dark',m==='dark');r.style.colorScheme=m;}catch(e){}})();`

// adminThemeScriptSource is the TS declaration the admin layouts embed, so the
// standalone app and the panel embedded in apps/web cannot drift apart.
func adminThemeScriptSource() string {
	return `// Applied synchronously, before the first paint, from the <head> below.
// Writes all three signals the stylesheets key off: data-theme-mode for the
// CSS variable cascade, .dark for Tailwind's darkMode: "class", and
// color-scheme so the browser's own scrollbars and inputs match.
const themeScript =
  "` + adminThemeScript + `";`
}

// adminIconMap returns the Lucide icon lookup map.
func adminIconMap() string {
	return `import {
  Users,
  UserCheck,
  FileText,
  File,
  Newspaper,
  BookOpen,
  MessageSquare,
  FolderOpen,
  Tag,
  Package,
  ShoppingCart,
  Receipt,
  CreditCard,
  User,
  UserCircle,
  Briefcase,
  Save,
  AlertTriangle,
  Lock,
  CheckSquare,
  Calendar,
  Paperclip,
  Image,
  Mail,
  Bell,
  Settings,
  Shield,
  ShieldCheck,
  Webhook,
  Building2,
  LayoutDashboard,
  Database,
  Layers,
  ChevronLeft,
  ChevronRight,
  ChevronUp,
  ChevronDown,
  Search,
  Sun,
  Moon,
  MoreHorizontal,
  Plus,
  Pencil,
  Trash2,
  Printer,
  Download,
  X,
  Maximize2,
  Minimize2,
  Check,
  AlertCircle,
  Loader2,
  Eye,
  EyeOff,
  Columns3,
  Inbox,
  Archive,
  ArchiveRestore,
  GripVertical,
  RefreshCw,
  Upload,
  Play,
  ExternalLink,
  LogOut,
  Monitor,
  Smartphone,
  Bold,
  Italic,
  Strikethrough,
  Heading1,
  Heading2,
  Heading3,
  List,
  ListOrdered,
  Quote,
  Code,
  Link,
  Undo,
  Redo,
  Activity,
  TrendingUp,
  TrendingDown,
  ArrowLeft,
  Cpu,
  Zap,
  Globe,
  Menu,
  ArrowUpRight,
  ArrowRight,
  Settings2,
  LayoutGrid,
  Home,
  ArrowDownRight,
  Flag,
  LogIn,
  UserPlus,
  AlignLeft,
  AlignCenter,
  AlignRight,
  AlignJustify,
  Highlighter,
  Palette,
  Table,
  Minus,
  Underline,
  CheckCircle,
  Server,
  HardDrive,
  Clock,
  Gauge,
  Copy,
  KeyRound,
  Fingerprint,
  Unlock,
  FileSpreadsheet,
  Music,
  Bug,
  ClipboardList,
  UploadCloud,
  Contact,
  FolderTree,
  FormInput,
  Gem,
  GitBranch,
  GitCommit,
  HeartHandshake,
  HelpCircle,
  LayoutTemplate,
  MapPin,
  Megaphone,
  MessageCircle,
  Percent,
  Rocket,
  ScrollText,
  Sparkles,
  Star,
  Target,
  Ticket,
  Truck,
  UsersRound,
  Workflow,
  // grit:icons:import
  type LucideIcon,
} from "lucide-react";

export const iconMap: Record<string, LucideIcon> = {
  Users,
  UserCheck,
  FileText,
  Newspaper,
  BookOpen,
  MessageSquare,
  FolderOpen,
  Tag,
  Package,
  ShoppingCart,
  Receipt,
  CreditCard,
  UserCircle,
  Briefcase,
  CheckSquare,
  Calendar,
  Paperclip,
  Image,
  Mail,
  Bell,
  Settings,
  Shield,
  ShieldCheck,
  Webhook,
  Building2,
  LayoutDashboard,
  Database,
  Layers,
  GripVertical,
  RefreshCw,
  Upload,
  AlertCircle,
  Bug,
  ClipboardList,
  UploadCloud,
  Contact,
  File,
  FolderTree,
  FormInput,
  Gem,
  GitBranch,
  GitCommit,
  HeartHandshake,
  HelpCircle,
  LayoutTemplate,
  Lock,
  MapPin,
  Megaphone,
  MessageCircle,
  Percent,
  Rocket,
  ScrollText,
  Sparkles,
  Star,
  Target,
  Ticket,
  TrendingUp,
  Truck,
  UsersRound,
  Workflow,
  // grit:icons:map
};

export function getIcon(name: string): LucideIcon {
  return iconMap[name] ?? FileText;
}

export {
  ChevronLeft,
  ChevronRight,
  ChevronUp,
  ChevronDown,
  Search,
  Sun,
  Moon,
  MoreHorizontal,
  Plus,
  Pencil,
  Trash2,
  Printer,
  Download,
  X,
  Maximize2,
  Minimize2,
  Check,
  AlertCircle,
  Loader2,
  Eye,
  EyeOff,
  Columns3,
  Inbox,
  Archive,
  ArchiveRestore,
  LayoutDashboard,
  Database,
  Briefcase,
  FolderOpen,
  Calendar,
  Mail,
  Bell,
  Settings,
  CreditCard,
  GripVertical,
  RefreshCw,
  Upload,
  File,
  Image,
  User,
  UserCheck,
  UserCircle,
  Save,
  AlertTriangle,
  Lock,
  Play,
  ExternalLink,
  LogOut,
  Monitor,
  Smartphone,
  Bold,
  Italic,
  Strikethrough,
  Heading1,
  Heading2,
  Heading3,
  List,
  ListOrdered,
  Quote,
  Code,
  Link,
  Undo,
  Redo,
  Activity,
  TrendingUp,
  TrendingDown,
  ArrowLeft,
  Shield,
  ShieldCheck,
  Webhook,
  Cpu,
  Zap,
  Globe,
  Menu,
  // v3.31.x — icons referenced by the new dashboard, profile, activity
  // feed, and sidebar chrome.
  Users,
  MessageSquare,
  ShoppingCart,
  FileText,
  ArrowUpRight,
  ArrowRight,
  Settings2,
  LayoutGrid,
  Home,
  ArrowDownRight,
  Flag,
  LogIn,
  UserPlus,
  AlignLeft,
  AlignCenter,
  AlignRight,
  AlignJustify,
  Highlighter,
  Palette,
  Table,
  Minus,
  Underline,
  // v3.31.12 — icons used by /system/health, /system/security,
  // /system/performance pages.
  CheckCircle,
  Server,
  HardDrive,
  Clock,
  Gauge,
  // v3.31.20 — form-shares admin page
  Copy,
  // v3.123 — API keys page + System Hub tile
  KeyRound,
  Unlock,
  // v3.331 — the passkey button on the sign-in page
  Fingerprint,
  // v3.31.31 — type-aware FilePreview icons (excel, audio).
  FileSpreadsheet,
  Music,
  // v3.103.0 — the form field "Generate" button. Note: an icon being in the
  // iconMap above does NOT make it a named export — it must be listed here too.
  Sparkles,
};
`
}

// adminLayoutComponent returns the admin layout with responsive sidebar.
func adminLayoutComponent() string {
	return `"use client";

import { useState, useEffect, useSyncExternalStore } from "react";
import { useRouter } from "next/navigation";
import { useMe } from "@/hooks/use-auth";
import { usePermissions } from "@/hooks/use-permissions";
import { CollapsibleSidebar } from "@/components/chrome/CollapsibleSidebar";
import { SessionWatchdog } from "@/components/chrome/SessionWatchdog";
import { QuickAccess } from "@/components/chrome/QuickAccess";
import { EmailVerifiedBanner } from "@/components/chrome/EmailVerifiedBanner";
import { Menu } from "@/lib/icons";
// grit:layout:imports

// v3.29: navbar is gone — pages now drop a <PageHeader> at the top of
// their JSX to get title/subtitle/search/dark-toggle/bell/user-menu in
// one consistent strip. The dashboard layout only owns sidebar + main.
// Nothing to subscribe to: the snapshot only differs between server and browser.
const subscribeNever = () => () => {};

export function AdminLayout({ children }: { children: React.ReactNode }) {
  const { data: user, isLoading, isError } = useMe();
  // Asked for beside /auth/me, not after it. The sidebar reads the same
  // query, so when it mounts the answer is already in the cache.
  const { permissions, isSuper, isLoading: permsLoading } = usePermissions();
  // False on the server and during hydration, true from the first browser
  // pass. Admin pages are written for the browser: they read localStorage,
  // window and search params while rendering, and prerendering them breaks
  // the build. So they render on that first pass, which is still before
  // /auth/me has answered.
  const inBrowser = useSyncExternalStore(subscribeNever, () => true, () => false);
  const router = useRouter();
  const [sidebarCollapsed, setSidebarCollapsed] = useState(false);
  const [mobileMenuOpen, setMobileMenuOpen] = useState(false);

  useEffect(() => {
    const stored = localStorage.getItem("grit-sidebar-collapsed");
    if (stored === "true") setSidebarCollapsed(true);
  }, []);

  // v3.31.15: redirect on BOTH isError (network/server down) AND
  // user === null (401 from /api/auth/me). The previous version
  // only handled isError and returned null on missing user, which
  // rendered a blank white page when the server was restarted
  // mid-session.
  useEffect(() => {
    if (isLoading) return;
    if (isError || user === null) {
      router.replace("/login");
    }
  }, [isError, user, isLoading, router]);

  // A USER holding no grants has nothing in the admin but their own profile and
  // account. The redirect used to cover only the dashboard, so such a user could
  // open System pages; the API refused their requests, and now the pages do not
  // open either.
  useEffect(() => {
    if (!user || user.role !== "USER" || permsLoading || isSuper || permissions.length > 0) return;
    const path = window.location.pathname;
    if (path.startsWith("/profile") || path.startsWith("/account")) return;
    router.replace("/profile");
  }, [user, router, permsLoading, isSuper, permissions]);

  const toggleSidebar = () => {
    const next = !sidebarCollapsed;
    setSidebarCollapsed(next);
    localStorage.setItem("grit-sidebar-collapsed", String(next));
  };

  // The page renders as soon as it is in the browser, so its queries start
  // alongside /auth/me rather than waiting for it to answer. Until the user is known an
  // overlay covers the page, and the effect above sends a signed-out visitor
  // to the login page. A signed-out page query gets a 401, one refresh attempt
  // and the same redirect from the API client.
  return (
    <div className="min-h-screen">
      {!user && (
        <div
          role="status"
          aria-busy="true"
          aria-label="Loading"
          className="fixed inset-0 z-50 flex items-center justify-center bg-background"
        >
          <div className="h-8 w-8 animate-spin rounded-full border-2 border-accent border-t-transparent" />
        </div>
      )}

      {user && (
        <>
          {/* v3.31.15: warns and refreshes the session before silent expiry. */}
          <SessionWatchdog />

          <CollapsibleSidebar
            user={user}
            collapsed={sidebarCollapsed}
            onToggleCollapsed={toggleSidebar}
            mobileOpen={mobileMenuOpen}
            onMobileClose={() => setMobileMenuOpen(false)}
          />
        </>
      )}

      <div
        className={` + "`" + `flex min-h-screen flex-col transition-all duration-200 ${
          sidebarCollapsed ? "md:ml-16" : "md:ml-64"
        }` + "`" + `}
      >
        {/* Every page of this admin starts with the same eight or so sidebar
            links. Without this, reaching the content by keyboard means tabbing
            past all of them on every page, which is WCAG 2.4.1 and is also
            just tedious. Hidden until focused, then it is the first thing you
            get. */}
        <a
          href="#main"
          className="sr-only rounded-lg bg-accent px-4 py-2 text-sm font-medium text-white focus:not-sr-only focus:absolute focus:left-4 focus:top-4 focus:z-50"
        >
          Skip to content
        </a>

        {/* Mobile menu button — only shown when the sidebar is hidden
            on small screens. PageHeader supplies the rest of the chrome. */}
        <button
          type="button"
          onClick={() => setMobileMenuOpen(true)}
          aria-label="Open menu"
          className="fixed top-3 left-3 z-30 inline-flex h-9 w-9 items-center justify-center rounded-lg border border-border bg-bg-elevated text-text-secondary shadow-sm md:hidden"
        >
          <Menu className="h-4 w-4" />
        </button>

        {/* grit:layout:banner */}
        <EmailVerifiedBanner />
        <main id="main" tabIndex={-1} className="flex-1 px-4 py-6 md:px-8 focus:outline-none">
          {inBrowser ? children : null}
        </main>
      </div>

      {/* Floating quick-access button (configurable) */}
      {user && <QuickAccess />}
    </div>
  );
}
`
}

// adminStatCards is components/chrome/StatCards.tsx: the stat cards a page
// header shows under its title row.
//
// They were part of a second page header, components/layout/page-header.tsx,
// which only resource-page.tsx used while 21 other pages used
// components/chrome/PageHeader.tsx (contact-app review M44). The cards moved
// here, PageHeader gained a stats prop, and the second header is gone.
func adminStatCards() string {
	return `"use client";

import { useQuery } from "@tanstack/react-query";
import { apiClient } from "@/lib/api-client";
import { resourceKeys } from "@/hooks/use-resource";
import { getIcon, TrendingUp, TrendingDown } from "@/lib/icons";

export interface StatCard {
  label: string;
  icon?: string;
  color?: "default" | "success" | "warning" | "danger" | "info";
  /** Either provide a static value... */
  value?: string | number;
  /** ...or an endpoint + field to fetch it from the API (e.g. endpoint: "/api/posts?page_size=1", field: "meta.total") */
  endpoint?: string;
  field?: string;
  /** Optional trend delta shown next to value */
  trend?: { value: number; direction: "up" | "down" };
  /** Shows the placeholder while a static value is still on its way. */
  loading?: boolean;
}

const colorClasses: Record<string, { bg: string; text: string }> = {
  default: { bg: "bg-accent/10", text: "text-accent" },
  success: { bg: "bg-success/10", text: "text-success" },
  warning: { bg: "bg-warning/10", text: "text-warning" },
  danger: { bg: "bg-danger/10", text: "text-danger" },
  info: { bg: "bg-info/10", text: "text-info" },
};

// Reads a dotted path from an object. e.g. getPath(data, "meta.total") → data.meta.total
function getPath(obj: unknown, path: string): unknown {
  return path.split(".").reduce<unknown>(
    (acc, key) => (acc == null ? acc : (acc as Record<string, unknown>)[key]),
    obj,
  );
}

function StatCardItem({ stat }: { stat: StatCard }) {
  const color = colorClasses[stat.color || "default"];
  const Icon = stat.icon ? getIcon(stat.icon) : null;

  // The key starts with the base endpoint (no query string) so a save to the
  // resource, which invalidates [endpoint], reaches this card too.
  const baseEndpoint = stat.endpoint ? stat.endpoint.split("?")[0] : "stat";

  const { data, isLoading } = useQuery({
    queryKey: [...resourceKeys.stats(baseEndpoint), stat.endpoint, stat.field],
    queryFn: async () => {
      if (!stat.endpoint) return null;
      const res = await apiClient.get(stat.endpoint);
      return res.data;
    },
    enabled: !!stat.endpoint,
  });

  const value =
    stat.value !== undefined
      ? stat.value
      : stat.endpoint && stat.field
      ? (getPath(data, stat.field) as string | number | undefined) ?? "—"
      : "—";

  return (
    <div className="rounded-xl border border-border bg-bg-secondary p-5 transition-colors hover:border-border/80">
      <div className="flex items-start justify-between">
        <div className="min-w-0 flex-1">
          <p className="text-[11px] font-semibold uppercase tracking-wider text-text-muted">
            {stat.label}
          </p>
          <div className="mt-2 flex items-baseline gap-2">
            {(isLoading && stat.endpoint) || stat.loading ? (
              <div className="h-7 w-16 rounded bg-bg-hover animate-pulse" />
            ) : (
              <p className="text-2xl font-bold text-foreground tabular-nums">
                {typeof value === "number" ? value.toLocaleString() : value}
              </p>
            )}
            {stat.trend && (
              <span
                className={` + "`" + `flex items-center gap-0.5 text-xs font-medium ${
                  stat.trend.direction === "up" ? "text-success" : "text-danger"
                }` + "`" + `}
              >
                {stat.trend.direction === "up" ? (
                  <TrendingUp className="h-3 w-3" />
                ) : (
                  <TrendingDown className="h-3 w-3" />
                )}
                {stat.trend.value}%
              </span>
            )}
          </div>
        </div>
        {Icon && (
          <div className={` + "`" + `flex h-9 w-9 items-center justify-center rounded-lg ${color.bg}` + "`" + `}>
            <Icon className={` + "`" + `h-4 w-4 ${color.text}` + "`" + `} />
          </div>
        )}
      </div>
    </div>
  );
}

/** The row of stat cards under a page header. */
export function StatCards({ stats }: { stats: StatCard[] }) {
  return (
    <div className="mb-8 grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-4">
      {stats.map((stat, i) => (
        <StatCardItem key={i} stat={stat} />
      ))}
    </div>
  );
}
`
}
