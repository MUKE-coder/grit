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
  Info,
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
  BarChart3,
  Pause,
  RotateCcw,
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
  Circle,
  CircleDashed,
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
  Info,
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
  Info,
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
  BarChart3,
  Pause,
  RotateCcw,
  TrendingUp,
  TrendingDown,
  ArrowLeft,
  Shield,
  ShieldCheck,
  Bug,
  Target,
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
  // v3.336.0 — the scaling readiness panel on /system/observability.
  Circle,
  CircleDashed,
  // v3.355.0 — the saga runs screen. Note again: being in the iconMap above
  // does NOT make an icon importable. It has to be listed here too.
  Workflow,
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
import { TopBar } from "@/components/chrome/TopBar";
import { CommandPalette } from "@/components/chrome/CommandPalette";
import { EmailVerifiedBanner } from "@/components/chrome/EmailVerifiedBanner";
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
    router.replace("/account");
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
        className={` + "`" + `flex min-h-screen flex-col transition-grit duration-200 ${
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

        {/* The panel's own chrome: where you are, how the system is doing,
            refresh, theme, notifications, and the menu button on small
            screens. One strip for the whole admin, instead of a copy of it
            inside every page's header competing with that page's own title
            and its one real action. */}
        <TopBar onOpenMenu={() => setMobileMenuOpen(true)} />

        {/* grit:layout:banner */}
        <EmailVerifiedBanner />
        {/* No bottom padding reserved any more. It was clearing the floating
            quick-access button, which offered one action from on top of the
            last row of whatever was underneath it. Cmd+K offers every action
            and covers nothing. */}
        <main id="main" tabIndex={-1} className="flex-1 px-4 py-6 pb-12 md:px-8 focus:outline-none">
          {inBrowser ? children : null}
        </main>
      </div>

      {user && <CommandPalette />}
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
import { StatStrip } from "@/components/ui/stat-strip";

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
  /**
   * One line under the value saying what the number means.
   *
   * This is the part that makes a number worth putting on a page. "4" is a
   * fact; "4, one of them inactive" is something somebody can act on. Without
   * it a stat strip is four numbers nobody can do anything with.
   */
  hint?: string;
  /** Colours the hint when it reports a state rather than a count. */
  hintTone?: "muted" | "good" | "warn" | "bad";
  /** Shows the placeholder while a static value is still on its way. */
  loading?: boolean;
  /**
   * The rows this card counted. Give it one and the card becomes a button
   * that narrows the table to them.
   *
   * Only the cards that stand for a subset get it. "Total" and "This Week"
   * have nothing to narrow to and stay as they were.
   */
  filter?: { field: string; value: string };
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

function StatCardItem({
  stat,
  onFilter,
  active,
  reserveHint,
}: {
  stat: StatCard;
  onFilter?: (filter: StatFilter) => void;
  active?: boolean;
  /**
   * Whether any cell in this strip has a hint.
   *
   * The line is reserved for all of them or for none. Reserved for all, the
   * four values share a baseline even when only one cell has something to
   * say. Reserved for none, a strip whose page has not adopted hints yet is
   * compact instead of carrying four blank lines of dead space.
   */
  reserveHint?: boolean;
}) {
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

  const hintClass =
    stat.hintTone === "good"
      ? "text-success"
      : stat.hintTone === "warn"
        ? "text-warning"
        : stat.hintTone === "bad"
          ? "text-danger"
          : "text-text-muted";

  const body = (
    <>
      {/* The label, in the mono eyebrow the panel uses for machine-ish words.
          The icon is inline with it rather than a tile floating on the right:
          at this size a 36px tinted square is the loudest thing in the cell,
          and it carries nothing the label does not already say. */}
      <p className="flex items-center gap-1.5 font-mono text-[11px] uppercase tracking-wider text-text-muted">
        {Icon && <Icon className={"h-3.5 w-3.5 shrink-0 " + color.text} aria-hidden="true" />}
        <span className="truncate">{stat.label}</span>
      </p>

      <div className="mt-2 flex items-baseline gap-2">
        {(isLoading && stat.endpoint) || stat.loading ? (
          <div className="h-[30px] w-16 animate-pulse rounded bg-bg-hover" />
        ) : (
          <p className="text-[30px] font-semibold leading-none tracking-[-0.02em] text-foreground tabular-nums">
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

      {/* What the number means. Kept even when this cell has nothing to say,
          as long as some cell in the strip does, so the values share one
          baseline instead of sitting at two heights. */}
      {(stat.hint || reserveHint) && (
        <p className={"mt-1.5 truncate text-[13px] " + hintClass}>{stat.hint ?? "\u00a0"}</p>
      )}
    </>
  );

  const base = "block h-full px-5 py-4 text-left transition-grit ";

  // A cell with nothing to narrow to stays a div, and loses the hover state
  // with it: a surface that lights up under the pointer is a promise, and
  // this one had nothing behind it.
  if (!stat.filter || !onFilter) {
    return <div className={base}>{body}</div>;
  }

  return (
    <button
      type="button"
      onClick={() => onFilter(stat.filter!)}
      aria-pressed={active}
      className={
        base +
        "w-full cursor-pointer focus:outline-none focus-visible:ring-2 focus-visible:ring-inset focus-visible:ring-accent " +
        (active ? "bg-accent-soft" : "hover:bg-bg-tertiary")
      }
    >
      {body}
      <span className="sr-only">
        {active ? ". Showing only these. Activate to show all." : ". Activate to show only these."}
      </span>
    </button>
  );
}

export interface StatFilter {
  field: string;
  value: string;
}

/**
 * The strip of numbers under a page header.
 *
 * One bordered card divided into cells, not four bordered cards with gaps
 * between them. Four cards is four of everything: four borders, four
 * backgrounds, four shadows, each announcing itself as a separate object. The
 * numbers belong together and are read together, and a strip says so with one
 * border and three hairlines.
 *
 * Two-up under 900px, one-up under 520px, where four columns is four columns
 * of nothing.
 */
export function StatCards({
  stats,
  onFilter,
  isFilterActive,
}: {
  stats: StatCard[];
  /** Given, every cell carrying a filter becomes a button. */
  onFilter?: (filter: StatFilter) => void;
  /** Whether the table is currently narrowed to that cell. */
  isFilterActive?: (filter: StatFilter) => boolean;
}) {
  if (stats.length === 0) return null;

  const reserveHint = stats.some((stat) => Boolean(stat.hint));

  return (
    <StatStrip className="mb-8">
      {stats.map((stat, i) => (
        <StatCardItem
          key={i}
          stat={stat}
          onFilter={onFilter}
          active={stat.filter ? (isFilterActive?.(stat.filter) ?? false) : false}
          reserveHint={reserveHint}
        />
      ))}
    </StatStrip>
  );
}
`
}
