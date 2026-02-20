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
        onToggle={toggleSidebar}
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

import { useState, useRef, useEffect, useCallback } from "react";
import { createPortal } from "react-dom";
import Link from "next/link";
import { usePathname, useRouter } from "next/navigation";
import { resources } from "@/resources";
import { getIcon, ChevronLeft, ChevronRight, ChevronDown, Sun, Moon, LayoutDashboard, UserCircle, LogOut } from "@/lib/icons";
import { useTheme } from "@/components/shared/theme-provider";
import { useLogout } from "@/hooks/use-auth";

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
  onToggle: () => void;
  mobileOpen: boolean;
  onMobileClose: () => void;
}

function UserMenu({ user, collapsed, fullName }: { user: User; collapsed: boolean; fullName: string }) {
  const [open, setOpen] = useState(false);
  const triggerRef = useRef<HTMLDivElement>(null);
  const menuRef = useRef<HTMLDivElement>(null);
  const [pos, setPos] = useState({ bottom: 0, left: 0, width: 0 });
  const router = useRouter();
  const { mutate: logout } = useLogout();

  const updatePosition = useCallback(() => {
    if (triggerRef.current) {
      const rect = triggerRef.current.getBoundingClientRect();
      setPos({ bottom: window.innerHeight - rect.top + 4, left: rect.left, width: rect.width });
    }
  }, []);

  useEffect(() => {
    if (!open) return;
    updatePosition();
    function handleClickOutside(e: MouseEvent) {
      const target = e.target as Node;
      if (
        triggerRef.current && !triggerRef.current.contains(target) &&
        menuRef.current && !menuRef.current.contains(target)
      ) {
        setOpen(false);
      }
    }
    document.addEventListener("mousedown", handleClickOutside);
    return () => document.removeEventListener("mousedown", handleClickOutside);
  }, [open, updatePosition]);

  const handleSignOut = () => {
    setOpen(false);
    logout(undefined, {
      onSuccess: () => router.push("/login"),
    });
  };

  const avatar = user.avatar ? (
    <img src={user.avatar} alt="" className="h-8 w-8 rounded-full object-cover" />
  ) : (
    <div className="h-8 w-8 rounded-full bg-accent/20 flex items-center justify-center shrink-0">
      <span className="text-sm font-medium text-accent">{user.first_name?.charAt(0)?.toUpperCase()}</span>
    </div>
  );

  const menu = open ? createPortal(
    <div
      ref={menuRef}
      className="fixed z-[9999] w-56 rounded-lg border border-border bg-bg-elevated shadow-xl"
      style={{ bottom: pos.bottom, left: pos.left, backgroundColor: "var(--bg-elevated, #22222e)" }}
    >
      <div className="p-3 border-b border-border">
        <p className="text-sm font-medium text-foreground truncate">{fullName}</p>
        <p className="text-xs text-text-muted truncate">{user.email}</p>
      </div>
      <div className="p-1">
        <button
          onClick={() => { setOpen(false); router.push("/profile"); }}
          className="flex w-full items-center gap-2 rounded-md px-3 py-2 text-sm text-text-secondary hover:bg-bg-hover hover:text-foreground transition-colors"
        >
          <UserCircle className="h-4 w-4" />
          Profile
        </button>
        <button
          onClick={handleSignOut}
          className="flex w-full items-center gap-2 rounded-md px-3 py-2 text-sm text-danger hover:bg-danger/10 transition-colors"
        >
          <LogOut className="h-4 w-4" />
          Sign Out
        </button>
      </div>
    </div>,
    document.body
  ) : null;

  return (
    <div>
      <div
        ref={triggerRef}
        onClick={() => { if (!open) updatePosition(); setOpen(!open); }}
        className={` + "`" + `flex cursor-pointer items-center gap-3 rounded-lg px-3 py-2 text-sm hover:bg-bg-hover transition-colors ${collapsed ? "justify-center" : ""}` + "`" + `}
      >
        {avatar}
        {!collapsed && (
          <>
            <div className="flex-1 min-w-0">
              <p className="text-sm font-medium text-foreground truncate">{fullName}</p>
              <p className="text-xs text-text-muted truncate">{user.role}</p>
            </div>
            <ChevronDown className={` + "`" + `h-3 w-3 text-text-muted transition-transform ${open ? "rotate-180" : ""}` + "`" + `} />
          </>
        )}
      </div>
      {menu}
    </div>
  );
}

export function Sidebar({ user, collapsed, onToggle, mobileOpen, onMobileClose }: SidebarProps) {
  const pathname = usePathname();
  const { theme, toggleTheme } = useTheme();

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

  const fullName = [user.first_name, user.last_name].filter(Boolean).join(" ") || "User";

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

      {/* Navigation */}
      <nav className="flex-1 overflow-y-auto p-3 space-y-1">
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

      {/* Bottom section */}
      <div className="border-t border-border p-3 space-y-2">
        {/* Theme toggle */}
        <button
          onClick={toggleTheme}
          className={` + "`" + `flex w-full items-center gap-3 rounded-lg px-3 py-2 text-sm text-text-secondary hover:bg-bg-hover hover:text-foreground transition-colors ${collapsed ? "justify-center" : ""}` + "`" + `}
          title={collapsed ? (theme === "dark" ? "Light mode" : "Dark mode") : undefined}
        >
          {theme === "dark" ? (
            <Sun className="h-4 w-4 shrink-0" />
          ) : (
            <Moon className="h-4 w-4 shrink-0" />
          )}
          {!collapsed && <span>{theme === "dark" ? "Light Mode" : "Dark Mode"}</span>}
        </button>

        {/* Collapse toggle (desktop only) */}
        <button
          onClick={onToggle}
          className="hidden lg:flex w-full items-center gap-3 rounded-lg px-3 py-2 text-sm text-text-secondary hover:bg-bg-hover hover:text-foreground transition-colors justify-center"
        >
          {collapsed ? (
            <ChevronRight className="h-4 w-4" />
          ) : (
            <ChevronLeft className="h-4 w-4" />
          )}
        </button>

        {/* User avatar menu */}
        <UserMenu user={user} collapsed={collapsed} fullName={fullName} />
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

// adminNavbar returns the enhanced navbar with breadcrumbs and search.
func adminNavbar() string {
	return `"use client";

import { useState } from "react";
import { usePathname } from "next/navigation";
import { useLogout } from "@/hooks/use-auth";
import { getResource } from "@/resources";
import { Search } from "@/lib/icons";

interface User {
  first_name: string;
  last_name: string;
  email: string;
  role: string;
}

export function Navbar({ user, onMenuToggle }: { user: User; onMenuToggle: () => void }) {
  const { mutate: logout } = useLogout();
  const [dropdownOpen, setDropdownOpen] = useState(false);
  const pathname = usePathname();

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
        i++; // skip the next segment
        continue;
      }
    }

    breadcrumbs.push({
      label: seg.charAt(0).toUpperCase() + seg.slice(1).replace(/-/g, " "),
      href,
    });
  }

  return (
    <header className="sticky top-0 z-30 flex h-16 items-center justify-between border-b border-border bg-background px-6">
      <div className="flex items-center gap-4">
        {/* Mobile menu button */}
        <button
          onClick={onMenuToggle}
          className="lg:hidden p-2 rounded-lg hover:bg-bg-hover text-text-secondary"
        >
          <svg className="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 6h16M4 12h16M4 18h16" />
          </svg>
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

      <div className="flex items-center gap-3">
        {/* Search */}
        <div className="hidden md:flex items-center gap-2 rounded-lg border border-border bg-bg-tertiary px-3 py-1.5">
          <Search className="h-4 w-4 text-text-muted" />
          <input
            type="text"
            placeholder="Search..."
            className="w-40 bg-transparent text-sm text-foreground placeholder:text-text-muted focus:outline-none"
          />
        </div>

        {/* User menu */}
        <div className="relative">
          <button
            onClick={() => setDropdownOpen(!dropdownOpen)}
            className="flex items-center gap-3 rounded-lg px-3 py-2 hover:bg-bg-hover transition-colors"
          >
            <div className="h-8 w-8 rounded-full bg-accent/20 flex items-center justify-center">
              <span className="text-sm font-medium text-accent">
                {user.first_name?.charAt(0)?.toUpperCase()}
              </span>
            </div>
            <span className="text-sm font-medium text-foreground hidden sm:block">
              {user.first_name}
            </span>
          </button>

          {dropdownOpen && (
            <>
              <div className="fixed inset-0 z-40" onClick={() => setDropdownOpen(false)} />
              <div className="absolute right-0 top-full mt-2 w-48 rounded-lg border border-border bg-bg-elevated shadow-lg z-50">
                <div className="px-4 py-3 border-b border-border">
                  <p className="text-sm font-medium text-foreground">{user.first_name} {user.last_name}</p>
                  <p className="text-xs text-text-muted">{user.email}</p>
                </div>
                <button
                  onClick={() => logout()}
                  className="w-full px-4 py-2.5 text-left text-sm text-danger hover:bg-bg-hover transition-colors"
                >
                  Sign Out
                </button>
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
