package scaffold

// adminThemeProvider returns the theme provider component.
func adminThemeProvider() string {
	return `"use client";

import { createContext, useContext, useEffect, useState } from "react";

type Theme = "dark" | "light";

interface ThemeContextValue {
  theme: Theme;
  toggleTheme: () => void;
}

const ThemeContext = createContext<ThemeContextValue>({
  theme: "dark",
  toggleTheme: () => {},
});

export function ThemeProvider({ children }: { children: React.ReactNode }) {
  const [theme, setTheme] = useState<Theme>("dark");

  useEffect(() => {
    const stored = localStorage.getItem("grit-theme") as Theme | null;
    if (stored) {
      setTheme(stored);
      document.documentElement.classList.toggle("light", stored === "light");
    }
  }, []);

  const toggleTheme = () => {
    const next = theme === "dark" ? "light" : "dark";
    setTheme(next);
    localStorage.setItem("grit-theme", next);
    document.documentElement.classList.toggle("light", next === "light");
  };

  return (
    <ThemeContext.Provider value={{ theme, toggleTheme }}>
      {children}
    </ThemeContext.Provider>
  );
}

export function useTheme() {
  return useContext(ThemeContext);
}
`
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
  Download,
  X,
  Check,
  AlertCircle,
  Loader2,
  Eye,
  EyeOff,
  Columns3,
  RefreshCw,
  Upload,
  Play,
  ExternalLink,
  LogOut,
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
  Building2,
  LayoutDashboard,
  Database,
  Layers,
  RefreshCw,
  Upload,
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
  Download,
  X,
  Check,
  AlertCircle,
  Loader2,
  Eye,
  EyeOff,
  Columns3,
  LayoutDashboard,
  Database,
  Briefcase,
  FolderOpen,
  Calendar,
  Mail,
  Bell,
  Settings,
  CreditCard,
  RefreshCw,
  Upload,
  File,
  Image,
  User,
  UserCircle,
  Save,
  AlertTriangle,
  Lock,
  Play,
  ExternalLink,
  LogOut,
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
};
`
}

// adminLayoutComponent returns the admin layout with responsive sidebar.
func adminLayoutComponent() string {
	return `"use client";

import { useState, useEffect } from "react";
import { useRouter } from "next/navigation";
import { useMe } from "@/hooks/use-auth";
import { Sidebar } from "./sidebar";
import { Navbar } from "./navbar";

export function AdminLayout({ children }: { children: React.ReactNode }) {
  const { data: user, isLoading, isError } = useMe();
  const router = useRouter();
  const [sidebarCollapsed, setSidebarCollapsed] = useState(false);
  const [mobileMenuOpen, setMobileMenuOpen] = useState(false);

  useEffect(() => {
    const stored = localStorage.getItem("grit-sidebar-collapsed");
    if (stored === "true") setSidebarCollapsed(true);
  }, []);

  useEffect(() => {
    if (isError) {
      router.push("/login");
    }
  }, [isError, router]);

  // Redirect USER role to profile page
  useEffect(() => {
    if (user && user.role === "USER" && window.location.pathname === "/dashboard") {
      router.replace("/profile");
    }
  }, [user, router]);

  const toggleSidebar = () => {
    const next = !sidebarCollapsed;
    setSidebarCollapsed(next);
    localStorage.setItem("grit-sidebar-collapsed", String(next));
  };

  if (isLoading) {
    return (
      <div className="flex min-h-screen items-center justify-center">
        <div className="h-8 w-8 animate-spin rounded-full border-2 border-accent border-t-transparent" />
      </div>
    );
  }

  if (!user) return null;

  return (
    <div className="flex min-h-screen">
      {/* Mobile overlay */}
      {mobileMenuOpen && (
        <div
          className="fixed inset-0 z-40 bg-black/50 lg:hidden"
          onClick={() => setMobileMenuOpen(false)}
        />
      )}

      <Sidebar
        user={user}
        collapsed={sidebarCollapsed}
        mobileOpen={mobileMenuOpen}
        onMobileClose={() => setMobileMenuOpen(false)}
      />

      <div
        className={` + "`" + `flex flex-1 flex-col transition-all duration-300 ${
          sidebarCollapsed ? "lg:ml-16" : "lg:ml-64"
        }` + "`" + `}
      >
        <Navbar
          user={user}
          onMenuToggle={() => setMobileMenuOpen(true)}
          collapsed={sidebarCollapsed}
          onToggleSidebar={toggleSidebar}
        />
        <main className="flex-1 p-6">{children}</main>
      </div>
    </div>
  );
}
`
}

// adminSidebar returns the enhanced collapsible sidebar with Lucide icons.
func adminSidebar() string {
	return `"use client";

import { useState } from "react";
import Link from "next/link";
import { usePathname } from "next/navigation";
import { resources } from "@/resources";
import { getIcon, ChevronDown, LayoutDashboard } from "@/lib/icons";

// Theme toggle and user menu have moved to the topbar (Navbar component).
// The sidebar now contains only navigation.

interface User {
  first_name: string;
  last_name: string;
  email: string;
  role: string;
  avatar?: string;
}

interface SidebarProps {
  user: User;
  collapsed: boolean;
  mobileOpen: boolean;
  onMobileClose: () => void;
}

export function Sidebar({ user, collapsed, mobileOpen, onMobileClose }: SidebarProps) {
  const pathname = usePathname();

  const isAdmin = user.role === "ADMIN" || user.role === "EDITOR";

  const navItems = isAdmin
    ? [
        { label: "Dashboard", href: "/dashboard", icon: "LayoutDashboard" },
        ...resources.map((r) => ({
          label: r.label?.plural ?? r.name,
          href: ` + "`" + `/resources/${r.slug}` + "`" + `,
          icon: r.icon,
        })),
      ]
    : [];

  // Profile link is always visible
  const profileItem = { label: "Profile", href: "/profile", icon: "UserCircle" };

  const systemItems = isAdmin
    ? [
        { label: "Jobs", href: "/system/jobs", icon: "Briefcase" },
        { label: "Files", href: "/system/files", icon: "FolderOpen" },
        { label: "Cron", href: "/system/cron", icon: "Calendar" },
        { label: "Mail", href: "/system/mail", icon: "Mail" },
        { label: "Security", href: "/system/security", icon: "Shield" },
      ]
    : [];

  const sidebarContent = (
    <>
      {/* Logo */}
      <div className={` + "`" + `flex h-16 items-center border-b border-border px-4 ${collapsed ? "justify-center" : "gap-2 px-6"}` + "`" + `}>
        <span className="text-xl font-bold text-accent">G</span>
        {!collapsed && (
          <>
            <span className="text-xl font-bold text-accent">rit</span>
            <span className="ml-1 rounded-md bg-accent/10 px-2 py-0.5 text-xs font-medium text-accent">
              Admin
            </span>
          </>
        )}
      </div>

      {/* Scrollable content: nav + bottom controls flow together */}
      <div className="flex-1 overflow-y-auto">
        {/* Navigation */}
        <nav className="p-3 space-y-1">
          {navItems.map((item) => {
            const Icon = getIcon(item.icon);
            const isActive = item.href === "/dashboard"
              ? pathname === "/dashboard"
              : pathname.startsWith(item.href);

            return (
              <Link
                key={item.href}
                href={item.href}
                onClick={onMobileClose}
                title={collapsed ? item.label : undefined}
                className={` + "`" + `flex items-center gap-3 rounded-lg px-3 py-2.5 text-sm font-medium transition-colors ${
                  isActive
                    ? "bg-accent/10 text-accent"
                    : "text-text-secondary hover:bg-bg-hover hover:text-foreground"
                } ${collapsed ? "justify-center" : ""}` + "`" + `}
              >
                <Icon className="h-4 w-4 shrink-0" />
                {!collapsed && <span>{item.label}</span>}
              </Link>
            );
          })}

          {/* Profile link */}
          {(() => {
            const ProfileIcon = getIcon(profileItem.icon);
            const isProfileActive = pathname === profileItem.href;
            return (
              <Link
                href={profileItem.href}
                onClick={onMobileClose}
                title={collapsed ? profileItem.label : undefined}
                className={` + "`" + `flex items-center gap-3 rounded-lg px-3 py-2.5 text-sm font-medium transition-colors ${
                  isProfileActive
                    ? "bg-accent/10 text-accent"
                    : "text-text-secondary hover:bg-bg-hover hover:text-foreground"
                } ${collapsed ? "justify-center" : ""}` + "`" + `}
              >
                <ProfileIcon className="h-4 w-4 shrink-0" />
                {!collapsed && <span>{profileItem.label}</span>}
              </Link>
            );
          })()}

          {/* System section */}
          {isAdmin && (
            <>
              {!collapsed && (
                <p className="px-3 mt-6 mb-2 text-xs font-semibold text-text-muted uppercase tracking-wider">
                  System
                </p>
              )}
              {collapsed && <div className="my-3 mx-3 border-t border-border" />}
              {systemItems.map((item) => {
                const Icon = getIcon(item.icon);
                const isActive = pathname.startsWith(item.href);

                return (
                  <Link
                    key={item.href}
                    href={item.href}
                    onClick={onMobileClose}
                    title={collapsed ? item.label : undefined}
                    className={` + "`" + `flex items-center gap-3 rounded-lg px-3 py-2.5 text-sm font-medium transition-colors ${
                      isActive
                        ? "bg-accent/10 text-accent"
                        : "text-text-secondary hover:bg-bg-hover hover:text-foreground"
                    } ${collapsed ? "justify-center" : ""}` + "`" + `}
                  >
                    <Icon className="h-4 w-4 shrink-0" />
                    {!collapsed && <span>{item.label}</span>}
                  </Link>
                );
              })}
            </>
          )}
        </nav>

      </div>
    </>
  );

  return (
    <>
      {/* Desktop sidebar */}
      <aside
        className={` + "`" + `fixed inset-y-0 left-0 z-50 hidden lg:flex flex-col bg-bg-secondary border-r border-border transition-all duration-300 ${
          collapsed ? "w-16" : "w-64"
        }` + "`" + `}
      >
        {sidebarContent}
      </aside>

      {/* Mobile sidebar */}
      <aside
        className={` + "`" + `fixed inset-y-0 left-0 z-50 flex flex-col w-64 bg-bg-secondary border-r border-border lg:hidden transition-transform duration-300 ${
          mobileOpen ? "translate-x-0" : "-translate-x-full"
        }` + "`" + `}
      >
        {sidebarContent}
      </aside>
    </>
  );
}
`
}

// adminNavbar returns the topbar with sidebar-collapse toggle (left) and
// theme/notifications/user cluster (right). Follows GRIT_STYLE_GUIDE §7.7.
func adminNavbar() string {
	return `"use client";

import { useState } from "react";
import { usePathname, useRouter } from "next/navigation";
import { useLogout } from "@/hooks/use-auth";
import { useTheme } from "@/components/shared/theme-provider";
import { getResource } from "@/resources";
import { Search, ChevronLeft, ChevronRight, Sun, Moon, Bell, Activity, Settings, CreditCard, LogOut } from "@/lib/icons";

interface User {
  first_name: string;
  last_name: string;
  email: string;
  role: string;
  avatar?: string;
}

interface NavbarProps {
  user: User;
  onMenuToggle: () => void;
  collapsed: boolean;
  onToggleSidebar: () => void;
}

export function Navbar({ user, onMenuToggle, collapsed, onToggleSidebar }: NavbarProps) {
  const { mutate: logout } = useLogout();
  const { theme, toggleTheme } = useTheme();
  const [userOpen, setUserOpen] = useState(false);
  const router = useRouter();
  const pathname = usePathname();

  const fullName = [user.first_name, user.last_name].filter(Boolean).join(" ") || "User";

  // Build breadcrumbs from pathname
  const segments = pathname.split("/").filter(Boolean);
  const breadcrumbs: { label: string; href: string }[] = [{ label: "Home", href: "/dashboard" }];

  for (let i = 0; i < segments.length; i++) {
    const href = "/" + segments.slice(0, i + 1).join("/");
    const seg = segments[i];

    if (seg === "resources" && segments[i + 1]) {
      const resource = getResource(segments[i + 1]);
      if (resource) {
        breadcrumbs.push({
          label: resource.label?.plural ?? resource.name,
          href: href + "/" + segments[i + 1],
        });
        i++;
        continue;
      }
    }

    breadcrumbs.push({
      label: seg.charAt(0).toUpperCase() + seg.slice(1).replace(/-/g, " "),
      href,
    });
  }

  const handleSignOut = () => {
    setUserOpen(false);
    logout(undefined, { onSuccess: () => router.push("/login") });
  };

  return (
    <header className="sticky top-0 z-30 flex h-16 items-center justify-between border-b border-border bg-background px-4 sm:px-6">
      {/* LEFT CLUSTER: mobile hamburger + collapse toggle + breadcrumbs */}
      <div className="flex items-center gap-3">
        {/* Mobile menu button (hidden on desktop) */}
        <button
          onClick={onMenuToggle}
          className="lg:hidden flex items-center justify-center h-9 w-9 rounded-lg hover:bg-bg-hover text-text-secondary transition-colors"
          aria-label="Open menu"
        >
          <svg className="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 6h16M4 12h16M4 18h16" />
          </svg>
        </button>

        {/* Sidebar collapse toggle (desktop only) */}
        <button
          onClick={onToggleSidebar}
          className="hidden lg:flex items-center justify-center h-9 w-9 rounded-lg hover:bg-bg-hover text-text-secondary transition-colors"
          aria-label={collapsed ? "Expand sidebar" : "Collapse sidebar"}
          title={collapsed ? "Expand sidebar" : "Collapse sidebar"}
        >
          {collapsed ? <ChevronRight className="h-4 w-4" /> : <ChevronLeft className="h-4 w-4" />}
        </button>

        {/* Breadcrumbs */}
        <nav className="hidden sm:flex items-center gap-1.5 text-sm">
          {breadcrumbs.map((crumb, i) => (
            <span key={crumb.href} className="flex items-center gap-1.5">
              {i > 0 && <span className="text-text-muted">/</span>}
              {i === breadcrumbs.length - 1 ? (
                <span className="text-foreground font-medium">{crumb.label}</span>
              ) : (
                <a
                  href={crumb.href}
                  className="text-text-secondary hover:text-foreground transition-colors"
                >
                  {crumb.label}
                </a>
              )}
            </span>
          ))}
        </nav>
      </div>

      {/* RIGHT CLUSTER: search + notifications + theme + user */}
      <div className="flex items-center gap-2">
        {/* Search */}
        <div className="hidden md:flex items-center gap-2 rounded-lg border border-border bg-bg-tertiary px-3 h-9">
          <Search className="h-4 w-4 text-text-muted" />
          <input
            type="text"
            placeholder="Search..."
            className="w-40 bg-transparent text-sm text-foreground placeholder:text-text-muted focus:outline-none"
          />
        </div>

        {/* Notifications */}
        <button
          className="relative flex items-center justify-center h-9 w-9 rounded-lg hover:bg-bg-hover text-text-secondary transition-colors"
          aria-label="Notifications"
          title="Notifications"
        >
          <Bell className="h-4 w-4" />
          {/* Unread dot (uncomment when wired up) */}
          {/* <span className="absolute top-2 right-2 h-2 w-2 rounded-full bg-accent" /> */}
        </button>

        {/* Theme toggle */}
        <button
          onClick={toggleTheme}
          className="flex items-center justify-center h-9 w-9 rounded-lg hover:bg-bg-hover text-text-secondary transition-colors"
          aria-label={theme === "dark" ? "Switch to light mode" : "Switch to dark mode"}
          title={theme === "dark" ? "Light mode" : "Dark mode"}
        >
          {theme === "dark" ? <Sun className="h-4 w-4" /> : <Moon className="h-4 w-4" />}
        </button>

        {/* User menu */}
        <div className="relative">
          <button
            onClick={() => setUserOpen(!userOpen)}
            className="flex items-center gap-2 rounded-lg pl-1 pr-2 py-1 hover:bg-bg-hover transition-colors"
            aria-label="User menu"
          >
            {user.avatar ? (
              <img src={user.avatar} alt="" className="h-8 w-8 rounded-full object-cover" />
            ) : (
              <div className="h-8 w-8 rounded-full bg-accent/20 flex items-center justify-center">
                <span className="text-sm font-medium text-accent">
                  {user.first_name?.charAt(0)?.toUpperCase()}
                </span>
              </div>
            )}
          </button>

          {userOpen && (
            <>
              <div className="fixed inset-0 z-40" onClick={() => setUserOpen(false)} />
              <div className="absolute right-0 top-full mt-2 w-64 rounded-xl border border-border bg-bg-elevated shadow-xl z-50 overflow-hidden">
                {/* Header: name + email */}
                <div className="px-4 py-3 border-b border-border">
                  <p className="text-sm font-semibold text-foreground truncate">{fullName}</p>
                  <p className="text-xs text-text-muted truncate">{user.email}</p>
                </div>

                {/* Menu items */}
                <div className="p-1">
                  <button
                    onClick={() => { setUserOpen(false); router.push("/profile"); }}
                    className="flex w-full items-center gap-2.5 rounded-md px-3 py-2 text-sm text-text-secondary hover:bg-bg-hover hover:text-foreground transition-colors"
                  >
                    <Activity className="h-4 w-4" />
                    User Activity
                  </button>
                  <button
                    onClick={() => { setUserOpen(false); router.push("/profile"); }}
                    className="flex w-full items-center gap-2.5 rounded-md px-3 py-2 text-sm text-text-secondary hover:bg-bg-hover hover:text-foreground transition-colors"
                  >
                    <Settings className="h-4 w-4" />
                    Settings
                  </button>
                  <button
                    onClick={() => { setUserOpen(false); router.push("/system/billing"); }}
                    className="flex w-full items-center justify-between gap-2.5 rounded-md px-3 py-2 text-sm text-text-secondary hover:bg-bg-hover hover:text-foreground transition-colors"
                  >
                    <span className="flex items-center gap-2.5">
                      <CreditCard className="h-4 w-4" />
                      Billing
                    </span>
                  </button>
                </div>

                {/* Sign out */}
                <div className="p-1 border-t border-border">
                  <button
                    onClick={handleSignOut}
                    className="flex w-full items-center gap-2.5 rounded-md px-3 py-2 text-sm text-danger hover:bg-danger/10 transition-colors"
                  >
                    <LogOut className="h-4 w-4" />
                    Log out
                  </button>
                </div>
              </div>
            </>
          )}
        </div>
      </div>
    </header>
  );
}
`
}

// adminPageHeader returns a consistent <PageHeader /> component used by every
// dashboard page (including generated resource pages). Structure:
//   - Breadcrumbs (optional)
//   - Title + Description + Actions row
//   - 4-card stats grid (optional)
// Follows GRIT_STYLE_GUIDE §7.8.
func adminPageHeader() string {
	return `"use client";

import Link from "next/link";
import { useQuery } from "@tanstack/react-query";
import { apiClient } from "@/lib/api-client";
import { getIcon, TrendingUp, TrendingDown } from "@/lib/icons";

export interface BreadcrumbItem {
  label: string;
  href?: string;
}

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
}

export interface PageHeaderProps {
  title: string;
  description?: string;
  breadcrumbs?: BreadcrumbItem[];
  actions?: React.ReactNode;
  stats?: StatCard[];
}

const colorClasses: Record<string, { bg: string; text: string }> = {
  default: { bg: "bg-accent/10", text: "text-accent" },
  success: { bg: "bg-success/10", text: "text-success" },
  warning: { bg: "bg-warning/10", text: "text-warning" },
  danger: { bg: "bg-danger/10", text: "text-danger" },
  info: { bg: "bg-info/10", text: "text-info" },
};

// Reads a dotted path from an object. e.g. getPath(data, "meta.total") → data.meta.total
function getPath(obj: any, path: string): any {
  return path.split(".").reduce((acc, key) => (acc == null ? acc : acc[key]), obj);
}

function StatCardItem({ stat }: { stat: StatCard }) {
  const color = colorClasses[stat.color || "default"];
  const Icon = stat.icon ? getIcon(stat.icon) : null;

  // Fetch value from endpoint if provided; otherwise use static value
  const { data, isLoading } = useQuery({
    queryKey: ["stat", stat.endpoint, stat.field],
    queryFn: async () => {
      if (!stat.endpoint) return null;
      const res = await apiClient.get(stat.endpoint);
      return res.data;
    },
    enabled: !!stat.endpoint,
    staleTime: 30_000,
  });

  const value =
    stat.value !== undefined
      ? stat.value
      : stat.endpoint && stat.field
      ? getPath(data, stat.field) ?? "—"
      : "—";

  return (
    <div className="rounded-xl border border-border bg-bg-secondary p-5 transition-colors hover:border-border/80">
      <div className="flex items-start justify-between">
        <div className="min-w-0 flex-1">
          <p className="text-[11px] font-semibold uppercase tracking-wider text-text-muted">
            {stat.label}
          </p>
          <div className="mt-2 flex items-baseline gap-2">
            {isLoading && stat.endpoint ? (
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

export function PageHeader({ title, description, breadcrumbs, actions, stats }: PageHeaderProps) {
  return (
    <div className="mb-8">
      {/* Breadcrumbs */}
      {breadcrumbs && breadcrumbs.length > 0 && (
        <nav className="mb-3 flex items-center gap-1.5 text-xs">
          {breadcrumbs.map((crumb, i) => (
            <span key={i} className="flex items-center gap-1.5">
              {i > 0 && <span className="text-text-muted">/</span>}
              {crumb.href && i < breadcrumbs.length - 1 ? (
                <Link
                  href={crumb.href}
                  className="text-text-secondary hover:text-foreground transition-colors"
                >
                  {crumb.label}
                </Link>
              ) : (
                <span className="text-foreground font-medium">{crumb.label}</span>
              )}
            </span>
          ))}
        </nav>
      )}

      {/* Title row */}
      <div className="flex flex-col sm:flex-row sm:items-start sm:justify-between gap-4">
        <div className="min-w-0 flex-1">
          <h1 className="text-2xl font-bold text-foreground tracking-tight">{title}</h1>
          {description && (
            <p className="mt-1 text-sm text-text-secondary">{description}</p>
          )}
        </div>
        {actions && <div className="flex items-center gap-2 shrink-0">{actions}</div>}
      </div>

      {/* Stats grid */}
      {stats && stats.length > 0 && (
        <div className="mt-6 grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-4">
          {stats.map((stat, i) => (
            <StatCardItem key={i} stat={stat} />
          ))}
        </div>
      )}
    </div>
  );
}
`
}
