package scaffold

// Layout & Auth UI Components for Grit Registry
// Layout inspired by Vercel, Linear, GitHub, Notion docs
// Auth inspired by Clerk, Auth0, Supabase Auth UI
// All components use Grit's dark theme CSS variable classes

// ── Layout Components ──────────────────────────────────────────────────────────

func layoutUIAppNavbar01() string {
	return `// app-navbar-01 — App top navbar (Vercel / Linear-inspired)
"use client";
import { useState } from "react";
export function AppNavbar({
  logo = "G",
  appName = "Grit",
  navItems = [
    { label: "Dashboard", href: "/dashboard", active: true },
    { label: "Projects", href: "/projects" },
    { label: "Team", href: "/team" },
    { label: "Settings", href: "/settings" },
  ],
  user = { name: "Sarah Chen", email: "sarah@example.com", initials: "SC" },
  onSignOut,
}: {
  logo?: string; appName?: string;
  navItems?: Array<{ label: string; href: string; active?: boolean }>;
  user?: { name: string; email: string; initials: string };
  onSignOut?: () => void;
}) {
  const [menuOpen, setMenuOpen] = useState(false);
  return (
    <nav className="h-14 border-b border-border bg-background/80 backdrop-blur-sm flex items-center px-4 gap-4 sticky top-0 z-40">
      <div className="flex items-center gap-2.5 flex-shrink-0">
        <div className="w-7 h-7 rounded-lg bg-accent flex items-center justify-center">
          <span className="text-white text-xs font-bold">{logo}</span>
        </div>
        <span className="text-sm font-semibold text-foreground">{appName}</span>
      </div>
      <div className="flex items-center gap-1 flex-1 overflow-x-auto">
        {navItems.map((item) => (
          <a
            key={item.href} href={item.href}
            className={"px-3 py-1.5 rounded-lg text-sm font-medium transition-colors whitespace-nowrap " + (item.active ? "bg-bg-elevated text-foreground" : "text-text-secondary hover:text-foreground hover:bg-bg-hover")}
          >
            {item.label}
          </a>
        ))}
      </div>
      <div className="relative flex-shrink-0">
        <button onClick={() => setMenuOpen(!menuOpen)} className="w-8 h-8 rounded-full bg-accent/20 flex items-center justify-center hover:bg-accent/30 transition-colors">
          <span className="text-xs font-bold text-accent">{user.initials}</span>
        </button>
        {menuOpen && (
          <div className="absolute right-0 top-full mt-2 w-48 bg-bg-elevated border border-border rounded-xl shadow-xl overflow-hidden z-50">
            <div className="px-3 py-3 border-b border-border">
              <p className="text-sm font-medium text-foreground">{user.name}</p>
              <p className="text-xs text-text-muted truncate">{user.email}</p>
            </div>
            <div className="p-1">
              <a href="/profile" className="flex items-center gap-2 px-2.5 py-2 rounded-lg text-sm text-text-secondary hover:text-foreground hover:bg-bg-hover transition-colors">Profile</a>
              <a href="/settings" className="flex items-center gap-2 px-2.5 py-2 rounded-lg text-sm text-text-secondary hover:text-foreground hover:bg-bg-hover transition-colors">Settings</a>
              <button onClick={onSignOut} className="w-full flex items-center gap-2 px-2.5 py-2 rounded-lg text-sm text-danger hover:bg-danger/10 transition-colors">Sign out</button>
            </div>
          </div>
        )}
      </div>
    </nav>
  );
}
`
}

func layoutUIMarketingNavbar01() string {
	return `// marketing-navbar-01 — Marketing site navbar (Vercel.com / Linear.app-inspired)
"use client";
import { useState } from "react";
export function MarketingNavbar({
  logo = "G",
  brand = "Grit",
  links = [
    { label: "Features", href: "#features" },
    { label: "Pricing", href: "#pricing" },
    { label: "Docs", href: "/docs" },
    { label: "Blog", href: "/blog" },
  ],
  ctaLabel = "Get started",
  ctaHref = "/signup",
  loginHref = "/login",
}: {
  logo?: string; brand?: string;
  links?: Array<{ label: string; href: string }>;
  ctaLabel?: string; ctaHref?: string; loginHref?: string;
}) {
  const [mobileOpen, setMobileOpen] = useState(false);
  return (
    <header className="sticky top-0 z-50 border-b border-border bg-background/80 backdrop-blur-md">
      <div className="max-w-6xl mx-auto px-4 h-16 flex items-center gap-8">
        <a href="/" className="flex items-center gap-2 flex-shrink-0">
          <div className="w-8 h-8 rounded-lg bg-accent flex items-center justify-center">
            <span className="text-white text-sm font-bold">{logo}</span>
          </div>
          <span className="text-base font-bold text-foreground">{brand}</span>
        </a>
        <nav className="hidden md:flex items-center gap-1 flex-1">
          {links.map((link) => (
            <a key={link.href} href={link.href} className="px-3 py-1.5 text-sm text-text-secondary hover:text-foreground transition-colors rounded-lg hover:bg-bg-hover">{link.label}</a>
          ))}
        </nav>
        <div className="hidden md:flex items-center gap-3 ml-auto">
          <a href={loginHref} className="text-sm text-text-secondary hover:text-foreground transition-colors">Sign in</a>
          <a href={ctaHref} className="text-sm font-semibold px-4 py-2 rounded-lg bg-accent hover:bg-accent/90 text-white transition-colors">{ctaLabel}</a>
        </div>
        <button onClick={() => setMobileOpen(!mobileOpen)} className="md:hidden ml-auto p-2 rounded-lg hover:bg-bg-hover transition-colors text-text-secondary">
          <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 6h16M4 12h16M4 18h16" /></svg>
        </button>
      </div>
      {mobileOpen && (
        <div className="md:hidden border-t border-border bg-bg-elevated px-4 py-4 flex flex-col gap-2">
          {links.map((link) => (
            <a key={link.href} href={link.href} className="px-3 py-2.5 text-sm text-text-secondary hover:text-foreground rounded-lg hover:bg-bg-hover transition-colors">{link.label}</a>
          ))}
          <div className="flex items-center gap-2 pt-2 border-t border-border mt-2">
            <a href={loginHref} className="flex-1 text-center text-sm py-2 rounded-lg border border-border text-text-secondary hover:text-foreground transition-colors">Sign in</a>
            <a href={ctaHref} className="flex-1 text-center text-sm py-2 rounded-lg bg-accent hover:bg-accent/90 text-white font-semibold transition-colors">{ctaLabel}</a>
          </div>
        </div>
      )}
    </header>
  );
}
`
}

func layoutUISidebarNav01() string {
	return `// sidebar-nav-01 — App sidebar with icon+label nav (Linear / Notion-inspired)
export function SidebarNav({
  items = [
    { label: "Dashboard", href: "/dashboard", icon: "home", active: true },
    { label: "Projects", href: "/projects", icon: "folder" },
    { label: "Issues", href: "/issues", icon: "circle", badge: 4 },
    { label: "Team", href: "/team", icon: "users" },
    { label: "Analytics", href: "/analytics", icon: "chart" },
    { label: "Settings", href: "/settings", icon: "settings" },
  ],
}: {
  items?: Array<{ label: string; href: string; icon: string; active?: boolean; badge?: number }>;
}) {
  const iconPaths: Record<string, string> = {
    home: "M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6",
    folder: "M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-6l-2-2H5a2 2 0 00-2 2z",
    circle: "M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z",
    users: "M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0z",
    chart: "M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z",
    settings: "M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z M15 12a3 3 0 11-6 0 3 3 0 016 0z",
  };
  return (
    <aside className="w-56 h-full bg-background border-r border-border flex flex-col py-4 px-2 gap-1">
      {items.map((item) => (
        <a
          key={item.href} href={item.href}
          className={"flex items-center gap-2.5 px-3 py-2 rounded-lg text-sm font-medium transition-colors " + (item.active ? "bg-bg-elevated text-foreground" : "text-text-secondary hover:text-foreground hover:bg-bg-hover")}
        >
          <svg className="w-4 h-4 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5} d={iconPaths[item.icon] ?? iconPaths["circle"]} />
          </svg>
          <span className="flex-1 truncate">{item.label}</span>
          {item.badge != null && item.badge > 0 && (
            <span className="text-xs bg-accent text-white rounded-full min-w-[1.25rem] h-5 flex items-center justify-center px-1 font-medium">{item.badge > 99 ? "99+" : item.badge}</span>
          )}
        </a>
      ))}
    </aside>
  );
}
`
}

func layoutUIBreadcrumb01() string {
	return `// breadcrumb-01 — Breadcrumb navigation (Vercel dashboard-inspired)
export function Breadcrumb({
  items = [
    { label: "Projects", href: "/projects" },
    { label: "grit-app", href: "/projects/grit-app" },
    { label: "Settings" },
  ],
}: {
  items?: Array<{ label: string; href?: string }>;
}) {
  return (
    <nav className="flex items-center gap-1 text-sm">
      {items.map((item, i) => (
        <div key={i} className="flex items-center gap-1">
          {i > 0 && (
            <svg className="w-3.5 h-3.5 text-text-muted" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" /></svg>
          )}
          {item.href ? (
            <a href={item.href} className="text-text-muted hover:text-foreground transition-colors">{item.label}</a>
          ) : (
            <span className="text-foreground font-medium">{item.label}</span>
          )}
        </div>
      ))}
    </nav>
  );
}
`
}

func layoutUIPageHeader01() string {
	return `// page-header-01 — Page header with title + description + action (Vercel / Linear-inspired)
export function PageHeader({
  title = "Team Settings",
  description = "Manage your team members, roles, and workspace preferences.",
  action,
  badge,
  breadcrumb,
}: {
  title?: string; description?: string;
  action?: React.ReactNode; badge?: string; breadcrumb?: React.ReactNode;
}) {
  return (
    <div className="flex flex-col gap-3 pb-6 border-b border-border">
      {breadcrumb && <div>{breadcrumb}</div>}
      <div className="flex items-start justify-between gap-4">
        <div className="flex-1 min-w-0">
          <div className="flex items-center gap-3">
            <h1 className="text-2xl font-bold text-foreground">{title}</h1>
            {badge && <span className="text-xs bg-accent/10 text-accent border border-accent/20 px-2 py-0.5 rounded-full font-medium">{badge}</span>}
          </div>
          {description && <p className="text-sm text-text-secondary mt-1 max-w-2xl">{description}</p>}
        </div>
        {action && <div className="flex-shrink-0">{action}</div>}
      </div>
    </div>
  );
}
`
}

func layoutUIMarketingFooter01() string {
	return `// marketing-footer-01 — Marketing site footer with columns (Vercel-inspired)
export function MarketingFooter({
  brand = "Grit",
  tagline = "The full-stack Go framework.",
  columns = [
    { title: "Product", links: [{ label: "Features", href: "#features" }, { label: "Pricing", href: "#pricing" }, { label: "Changelog", href: "/changelog" }] },
    { title: "Resources", links: [{ label: "Documentation", href: "/docs" }, { label: "Blog", href: "/blog" }, { label: "Community", href: "/community" }] },
    { title: "Company", links: [{ label: "About", href: "/about" }, { label: "Careers", href: "/careers" }, { label: "Contact", href: "/contact" }] },
    { title: "Legal", links: [{ label: "Privacy", href: "/privacy" }, { label: "Terms", href: "/terms" }] },
  ],
}: {
  brand?: string; tagline?: string;
  columns?: Array<{ title: string; links: Array<{ label: string; href: string }> }>;
}) {
  return (
    <footer className="border-t border-border bg-background py-12 px-4">
      <div className="max-w-6xl mx-auto">
        <div className="grid grid-cols-2 md:grid-cols-5 gap-8">
          <div className="col-span-2 md:col-span-1">
            <div className="flex items-center gap-2 mb-2">
              <div className="w-7 h-7 rounded-lg bg-accent flex items-center justify-center">
                <span className="text-white text-xs font-bold">{brand[0]}</span>
              </div>
              <span className="font-bold text-foreground">{brand}</span>
            </div>
            <p className="text-sm text-text-muted">{tagline}</p>
          </div>
          {columns.map((col) => (
            <div key={col.title}>
              <h4 className="text-xs font-semibold text-foreground uppercase tracking-wider mb-4">{col.title}</h4>
              <ul className="flex flex-col gap-2.5">
                {col.links.map((link) => (
                  <li key={link.href}>
                    <a href={link.href} className="text-sm text-text-muted hover:text-foreground transition-colors">{link.label}</a>
                  </li>
                ))}
              </ul>
            </div>
          ))}
        </div>
        <div className="border-t border-border mt-10 pt-6 flex flex-col sm:flex-row items-center justify-between gap-2">
          <p className="text-xs text-text-muted">&copy; {new Date().getFullYear()} {brand}. All rights reserved.</p>
          <p className="text-xs text-text-muted">Built with{" "}<a href="https://gritframework.dev" className="text-accent hover:underline">{brand}</a></p>
        </div>
      </div>
    </footer>
  );
}
`
}

func layoutUIAlertBanner01() string {
	return `// alert-banner-01 — Alert/warning/success/error banner (common UI pattern)
export function AlertBanner({
  variant = "info",
  title,
  message = "Something needs your attention.",
  dismissible = true,
  onDismiss,
  action,
}: {
  variant?: "info" | "success" | "warning" | "error";
  title?: string; message?: string; dismissible?: boolean;
  onDismiss?: () => void;
  action?: { label: string; onClick: () => void };
}) {
  const config = {
    info: { bg: "bg-accent/10", border: "border-accent/30", icon: "text-accent", path: "M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" },
    success: { bg: "bg-success/10", border: "border-success/30", icon: "text-success", path: "M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" },
    warning: { bg: "bg-yellow-500/10", border: "border-yellow-500/30", icon: "text-yellow-400", path: "M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" },
    error: { bg: "bg-danger/10", border: "border-danger/30", icon: "text-danger", path: "M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z" },
  };
  const c = config[variant];
  return (
    <div className={"flex items-start gap-3 p-4 rounded-xl border " + c.bg + " " + c.border}>
      <svg className={"w-5 h-5 flex-shrink-0 mt-0.5 " + c.icon} fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d={c.path} />
      </svg>
      <div className="flex-1 min-w-0">
        {title && <p className="text-sm font-semibold text-foreground mb-0.5">{title}</p>}
        <p className="text-sm text-text-secondary">{message}</p>
        {action && (
          <button onClick={action.onClick} className={"mt-2 text-sm font-semibold hover:underline " + c.icon}>{action.label}</button>
        )}
      </div>
      {dismissible && (
        <button onClick={onDismiss} className="p-1 rounded text-text-muted hover:text-foreground transition-colors flex-shrink-0">
          <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" /></svg>
        </button>
      )}
    </div>
  );
}
`
}

func layoutUIModalShell01() string {
	return `// modal-shell-01 — Modal dialog shell (Headless UI / Radix-inspired pattern)
"use client";
import { useEffect } from "react";
export function ModalShell({
  title = "Confirm action",
  description,
  children,
  footer,
  onClose,
  size = "md",
}: {
  title?: string; description?: string;
  children?: React.ReactNode; footer?: React.ReactNode;
  onClose?: () => void; size?: "sm" | "md" | "lg" | "xl";
}) {
  const sizeClasses = { sm: "max-w-sm", md: "max-w-md", lg: "max-w-lg", xl: "max-w-2xl" };
  useEffect(() => {
    const handler = (e: KeyboardEvent) => { if (e.key === "Escape") onClose?.(); };
    window.addEventListener("keydown", handler);
    return () => window.removeEventListener("keydown", handler);
  }, [onClose]);
  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center p-4">
      <div className="absolute inset-0 bg-background/60 backdrop-blur-sm" onClick={onClose} />
      <div className={"relative w-full " + sizeClasses[size] + " bg-bg-elevated border border-border rounded-2xl shadow-2xl overflow-hidden"}>
        <div className="flex items-start justify-between p-6 border-b border-border">
          <div>
            <h2 className="text-base font-semibold text-foreground">{title}</h2>
            {description && <p className="text-sm text-text-secondary mt-1">{description}</p>}
          </div>
          <button onClick={onClose} className="p-1 rounded-lg text-text-muted hover:text-foreground hover:bg-bg-hover transition-colors -mr-1">
            <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" /></svg>
          </button>
        </div>
        {children && <div className="p-6">{children}</div>}
        {footer && <div className="px-6 py-4 border-t border-border flex items-center justify-end gap-3">{footer}</div>}
      </div>
    </div>
  );
}
`
}

func layoutUICodeBlock01() string {
	return `// code-block-01 — Code display block (Vercel docs / GitHub-inspired)
"use client";
import { useState } from "react";
export function CodeBlock({
  code = "grit new my-app --template=saas",
  language = "bash",
  filename,
  showLineNumbers = false,
}: {
  code?: string; language?: string; filename?: string; showLineNumbers?: boolean;
}) {
  const [copied, setCopied] = useState(false);
  const copy = () => {
    navigator.clipboard.writeText(code).catch(() => {});
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  };
  const lines = code.split("\n");
  return (
    <div className="bg-background border border-border rounded-xl overflow-hidden">
      <div className="flex items-center justify-between px-4 py-2.5 border-b border-border bg-bg-elevated">
        <div className="flex items-center gap-2">
          {filename ? (
            <span className="text-xs text-text-secondary font-mono">{filename}</span>
          ) : (
            <span className="text-xs text-text-muted uppercase tracking-wider">{language}</span>
          )}
        </div>
        <button onClick={copy} className={"flex items-center gap-1.5 text-xs px-2 py-1 rounded transition-colors " + (copied ? "text-success" : "text-text-muted hover:text-foreground")}>
          {copied ? (
            <><svg className="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" /></svg>Copied</>
          ) : (
            <><svg className="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M8 5H6a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2v-1M8 5a2 2 0 002 2h2a2 2 0 002-2M8 5a2 2 0 012-2h2a2 2 0 012 2" /></svg>Copy</>
          )}
        </button>
      </div>
      <pre className="p-4 overflow-x-auto text-sm font-mono">
        {showLineNumbers ? (
          <table className="w-full border-collapse">
            <tbody>
              {lines.map((line, i) => (
                <tr key={i}>
                  <td className="pr-4 text-right text-text-muted select-none w-8">{i + 1}</td>
                  <td className="text-text-secondary">{line}</td>
                </tr>
              ))}
            </tbody>
          </table>
        ) : (
          <code className="text-text-secondary">{code}</code>
        )}
      </pre>
    </div>
  );
}
`
}

func layoutUIStatCard01() string {
	return `// stat-card-01 — Simple stat card (Stripe / Vercel dashboard-inspired)
export function StatCard({
  label = "Total Revenue",
  value = "$48,352",
  subtext = "+12% from last month",
  icon = "dollar",
  trend = "up",
}: {
  label?: string; value?: string; subtext?: string;
  icon?: "dollar" | "users" | "chart" | "rocket";
  trend?: "up" | "down" | "neutral";
}) {
  const iconPaths: Record<string, string> = {
    dollar: "M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z",
    users: "M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0z",
    chart: "M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z",
    rocket: "M15.59 14.37a6 6 0 01-5.84 7.38v-4.82m5.84-2.56a14.98 14.98 0 006.16-12.12A14.98 14.98 0 009.63 8.41m5.96 5.96a14.926 14.926 0 01-5.96 2.91m0 0a6 6 0 01-7.38-5.84m7.38 5.84A14.928 14.928 0 003.67 9.63",
  };
  const trendColor = trend === "up" ? "text-success" : trend === "down" ? "text-danger" : "text-text-muted";
  return (
    <div className="bg-bg-elevated border border-border rounded-xl p-5 flex flex-col gap-3">
      <div className="flex items-center justify-between">
        <span className="text-sm text-text-secondary">{label}</span>
        <div className="w-8 h-8 rounded-lg bg-accent/10 flex items-center justify-center">
          <svg className="w-4 h-4 text-accent" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5} d={iconPaths[icon] ?? iconPaths["chart"]} />
          </svg>
        </div>
      </div>
      <p className="text-2xl font-bold text-foreground tracking-tight">{value}</p>
      {subtext && <p className={"text-xs font-medium " + trendColor}>{trend === "up" ? "↑ " : trend === "down" ? "↓ " : ""}{subtext}</p>}
    </div>
  );
}
`
}

func layoutUISectionHeader01() string {
	return `// section-header-01 — Marketing section header (centered, Vercel landing-inspired)
export function SectionHeader({
  badge,
  title = "Everything you need to ship",
  subtitle = "Grit gives you the batteries-included stack so you can focus on building your product, not infrastructure.",
  align = "center",
}: {
  badge?: string; title?: string; subtitle?: string;
  align?: "left" | "center";
}) {
  const isCenter = align === "center";
  return (
    <div className={"flex flex-col gap-4 " + (isCenter ? "items-center text-center" : "items-start")}>
      {badge && (
        <span className="text-xs font-semibold border border-accent/30 bg-accent/10 text-accent px-3 py-1 rounded-full">
          {badge}
        </span>
      )}
      <h2 className="text-3xl sm:text-4xl font-bold text-foreground leading-tight max-w-2xl">
        {title}
      </h2>
      {subtitle && (
        <p className={"text-text-secondary leading-relaxed " + (isCenter ? "max-w-xl" : "max-w-2xl")}>
          {subtitle}
        </p>
      )}
    </div>
  );
}
`
}

func layoutUITabNav01() string {
	return `// tab-nav-01 — Tab navigation bar (GitHub repo tabs-inspired)
"use client";
import { useState } from "react";
export function TabNav({
  tabs = [
    { id: "overview", label: "Overview" },
    { id: "code", label: "Code", badge: 0 },
    { id: "issues", label: "Issues", badge: 4 },
    { id: "pulls", label: "Pull Requests", badge: 2 },
    { id: "settings", label: "Settings" },
  ],
  defaultTab = "overview",
  onChange,
}: {
  tabs?: Array<{ id: string; label: string; badge?: number }>;
  defaultTab?: string; onChange?: (id: string) => void;
}) {
  const [active, setActive] = useState(defaultTab);
  return (
    <div className="border-b border-border flex items-end gap-0 overflow-x-auto">
      {tabs.map((tab) => (
        <button
          key={tab.id}
          onClick={() => { setActive(tab.id); onChange?.(tab.id); }}
          className={"flex items-center gap-1.5 px-4 py-3 text-sm font-medium border-b-2 -mb-px transition-colors whitespace-nowrap " + (active === tab.id ? "border-accent text-foreground" : "border-transparent text-text-muted hover:text-foreground hover:border-border")}
        >
          {tab.label}
          {tab.badge != null && tab.badge > 0 && (
            <span className={"text-xs rounded-full px-1.5 py-0.5 min-w-[1.25rem] text-center " + (active === tab.id ? "bg-accent/20 text-accent" : "bg-bg-hover text-text-muted")}>{tab.badge}</span>
          )}
        </button>
      ))}
    </div>
  );
}
`
}

func layoutUILoadingSkeleton01() string {
	return `// loading-skeleton-01 — Skeleton loading placeholder (common pattern)
export function Skeleton({ className = "" }: { className?: string }) {
  return <div className={"animate-pulse bg-bg-hover rounded-lg " + className} />;
}

// Pre-built skeleton layouts
export function CardSkeleton() {
  return (
    <div className="bg-bg-elevated border border-border rounded-xl p-5 flex flex-col gap-4">
      <div className="flex items-center gap-3">
        <Skeleton className="w-10 h-10 rounded-full" />
        <div className="flex-1 flex flex-col gap-2">
          <Skeleton className="h-4 w-32" />
          <Skeleton className="h-3 w-24" />
        </div>
      </div>
      <Skeleton className="h-3 w-full" />
      <Skeleton className="h-3 w-4/5" />
      <Skeleton className="h-3 w-2/3" />
      <div className="flex gap-2 mt-2">
        <Skeleton className="h-8 flex-1 rounded-lg" />
        <Skeleton className="h-8 w-20 rounded-lg" />
      </div>
    </div>
  );
}

export function TableRowSkeleton({ cols = 4 }: { cols?: number }) {
  return (
    <div className="flex items-center gap-4 py-3 px-4">
      <Skeleton className="w-4 h-4 rounded" />
      {Array.from({ length: cols }).map((_, i) => (
        <Skeleton key={i} className={"h-4 " + (i === 0 ? "flex-1" : "w-24")} />
      ))}
    </div>
  );
}
`
}

func layoutUIToastNotification01() string {
	return `// toast-notification-01 — Toast notification (Sonner / React Hot Toast-inspired)
export function Toast({
  message = "Changes saved successfully",
  variant = "success",
  action,
  onDismiss,
}: {
  message?: string; variant?: "success" | "error" | "info" | "warning";
  action?: { label: string; onClick: () => void }; onDismiss?: () => void;
}) {
  const config = {
    success: { icon: "text-success bg-success/20", path: "M5 13l4 4L19 7" },
    error: { icon: "text-danger bg-danger/20", path: "M6 18L18 6M6 6l12 12" },
    info: { icon: "text-accent bg-accent/20", path: "M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" },
    warning: { icon: "text-yellow-400 bg-yellow-500/20", path: "M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" },
  };
  const c = config[variant];
  return (
    <div className="flex items-center gap-3 bg-bg-elevated border border-border rounded-xl px-4 py-3 shadow-lg min-w-64 max-w-sm">
      <div className={"w-7 h-7 rounded-full flex items-center justify-center flex-shrink-0 " + c.icon}>
        <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d={c.path} />
        </svg>
      </div>
      <span className="text-sm text-foreground flex-1">{message}</span>
      {action && (
        <button onClick={action.onClick} className="text-xs text-accent hover:underline font-medium flex-shrink-0">{action.label}</button>
      )}
      {onDismiss && (
        <button onClick={onDismiss} className="p-0.5 rounded text-text-muted hover:text-foreground transition-colors flex-shrink-0">
          <svg className="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" /></svg>
        </button>
      )}
    </div>
  );
}
`
}

func layoutUIDataList01() string {
	return `// data-list-01 — Key-value data list (Stripe details panel-inspired)
export function DataList({
  items = [
    { label: "Status", value: "Active", valueClass: "text-success" },
    { label: "Plan", value: "Pro Monthly" },
    { label: "Next billing", value: "Apr 1, 2025" },
    { label: "Amount", value: "$29.00" },
  ],
}: {
  items?: Array<{ label: string; value: string; valueClass?: string }>;
}) {
  return (
    <dl className="divide-y divide-border">
      {items.map((item) => (
        <div key={item.label} className="flex items-center justify-between py-3 gap-4">
          <dt className="text-sm text-text-secondary">{item.label}</dt>
          <dd className={"text-sm font-medium " + (item.valueClass ?? "text-foreground")}>{item.value}</dd>
        </div>
      ))}
    </dl>
  );
}
`
}

func layoutUIErrorPage01() string {
	return `// error-page-01 — Error page (404 / 500) (Vercel / Linear-inspired)
export function ErrorPage({
  code = "404",
  title = "Page not found",
  description = "The page you're looking for doesn't exist or has been moved.",
  homeHref = "/",
  backLabel = "Go back",
}: {
  code?: string; title?: string; description?: string;
  homeHref?: string; backLabel?: string;
}) {
  return (
    <div className="min-h-[60vh] flex flex-col items-center justify-center px-4 text-center">
      <div className="text-8xl font-black text-bg-hover mb-4 select-none">{code}</div>
      <h1 className="text-2xl font-bold text-foreground mb-2">{title}</h1>
      <p className="text-text-secondary mb-8 max-w-md">{description}</p>
      <div className="flex items-center gap-3">
        <button onClick={() => window.history.back()} className="px-4 py-2 rounded-lg border border-border text-text-secondary hover:text-foreground hover:bg-bg-hover transition-colors text-sm">
          {backLabel}
        </button>
        <a href={homeHref} className="px-4 py-2 rounded-lg bg-accent hover:bg-accent/90 text-white text-sm font-semibold transition-colors">
          Back to home
        </a>
      </div>
    </div>
  );
}
`
}

func layoutUICardGrid01() string {
	return `// card-grid-01 — Responsive card grid wrapper with header
export function CardGrid({
  title,
  description,
  action,
  columns = 3,
  children,
}: {
  title?: string; description?: string; action?: React.ReactNode;
  columns?: 2 | 3 | 4; children?: React.ReactNode;
}) {
  const colClass = { 2: "sm:grid-cols-2", 3: "sm:grid-cols-2 lg:grid-cols-3", 4: "sm:grid-cols-2 lg:grid-cols-4" };
  return (
    <section className="flex flex-col gap-6">
      {(title || action) && (
        <div className="flex items-start justify-between gap-4">
          <div>
            {title && <h2 className="text-lg font-semibold text-foreground">{title}</h2>}
            {description && <p className="text-sm text-text-secondary mt-0.5">{description}</p>}
          </div>
          {action}
        </div>
      )}
      <div className={"grid grid-cols-1 gap-4 " + colClass[columns]}>{children}</div>
    </section>
  );
}
`
}

func layoutUIFloatingAction01() string {
	return `// floating-action-01 — Floating action button (Notion "New page"-inspired)
export function FloatingActionButton({
  label = "New",
  icon = "plus",
  onClick,
  position = "bottom-right",
}: {
  label?: string; icon?: "plus" | "edit" | "message";
  onClick?: () => void; position?: "bottom-right" | "bottom-left";
}) {
  const iconPaths: Record<string, string> = {
    plus: "M12 4v16m8-8H4",
    edit: "M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z",
    message: "M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z",
  };
  const posClass = position === "bottom-right" ? "bottom-6 right-6" : "bottom-6 left-6";
  return (
    <button
      onClick={onClick}
      className={"fixed " + posClass + " flex items-center gap-2 px-4 py-3 bg-accent hover:bg-accent/90 text-white rounded-full shadow-lg hover:shadow-accent/30 transition-all font-semibold text-sm hover:scale-105 active:scale-95 z-40"}
    >
      <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d={iconPaths[icon] ?? iconPaths["plus"]} />
      </svg>
      {label}
    </button>
  );
}
`
}

// ── Extra Auth Components ───────────────────────────────────────────────────────

func authUITwoFactorSetup01() string {
	return `// two-factor-setup-01 — 2FA setup card with QR + backup codes (GitHub / Clerk-inspired)
"use client";
import { useState } from "react";
export function TwoFactorSetup({
  qrCodeUrl = "",
  secret = "JBSWY3DPEHPK3PXP",
  backupCodes = ["a1b2-c3d4", "e5f6-g7h8", "i9j0-k1l2", "m3n4-o5p6", "q7r8-s9t0", "u1v2-w3x4"],
  onComplete,
}: {
  qrCodeUrl?: string; secret?: string; backupCodes?: string[];
  onComplete?: (code: string) => void;
}) {
  const [step, setStep] = useState<"qr" | "verify" | "backup">("qr");
  const [code, setCode] = useState("");
  return (
    <div className="bg-bg-elevated border border-border rounded-2xl p-6 w-full max-w-md flex flex-col gap-5">
      <div className="flex items-center gap-3">
        {["qr", "verify", "backup"].map((s, i) => (
          <div key={s} className="flex items-center gap-2">
            <div className={"w-6 h-6 rounded-full flex items-center justify-center text-xs font-bold " + (s === step ? "bg-accent text-white" : ["qr","verify","backup"].indexOf(step) > i ? "bg-success text-white" : "bg-bg-hover text-text-muted border border-border")}>{i + 1}</div>
            {i < 2 && <div className={"h-0.5 w-8 " + (["qr","verify","backup"].indexOf(step) > i ? "bg-success" : "bg-border")} />}
          </div>
        ))}
      </div>
      {step === "qr" && (
        <div className="flex flex-col gap-4">
          <h3 className="text-base font-semibold text-foreground">Scan QR Code</h3>
          <p className="text-sm text-text-secondary">Scan this code with your authenticator app (Google Authenticator, Authy, etc.)</p>
          <div className="w-40 h-40 bg-white rounded-xl flex items-center justify-center self-center">
            {qrCodeUrl ? <img src={qrCodeUrl} alt="QR Code" className="w-full h-full" /> : <span className="text-xs text-gray-400">QR Code</span>}
          </div>
          <div className="flex flex-col gap-1">
            <p className="text-xs text-text-muted">Can't scan? Enter this code manually:</p>
            <code className="text-sm font-mono text-foreground bg-background border border-border rounded-lg px-3 py-2 select-all">{secret}</code>
          </div>
          <button onClick={() => setStep("verify")} className="w-full py-2.5 rounded-lg bg-accent hover:bg-accent/90 text-white font-semibold transition-colors">Continue</button>
        </div>
      )}
      {step === "verify" && (
        <div className="flex flex-col gap-4">
          <h3 className="text-base font-semibold text-foreground">Verify Setup</h3>
          <p className="text-sm text-text-secondary">Enter the 6-digit code from your authenticator app to confirm setup.</p>
          <input
            value={code} onChange={(e) => setCode(e.target.value.slice(0, 6))}
            placeholder="000000" maxLength={6}
            className="text-center text-2xl font-mono tracking-widest bg-background border border-border rounded-xl px-4 py-3 text-foreground outline-none focus:border-accent transition-colors"
          />
          <button onClick={() => { if (code.length === 6) { onComplete?.(code); setStep("backup"); } }} disabled={code.length < 6} className="w-full py-2.5 rounded-lg bg-accent hover:bg-accent/90 text-white font-semibold transition-colors disabled:opacity-50">Verify</button>
        </div>
      )}
      {step === "backup" && (
        <div className="flex flex-col gap-4">
          <h3 className="text-base font-semibold text-foreground">Save backup codes</h3>
          <p className="text-sm text-text-secondary">Store these codes somewhere safe. Each can be used once if you lose access to your device.</p>
          <div className="grid grid-cols-2 gap-2 bg-background border border-border rounded-xl p-4">
            {backupCodes.map((c) => <code key={c} className="text-xs font-mono text-text-secondary">{c}</code>)}
          </div>
          <button onClick={() => navigator.clipboard.writeText(backupCodes.join("\n")).catch(() => {})} className="w-full py-2 rounded-lg border border-border text-text-secondary hover:text-foreground hover:bg-bg-hover text-sm transition-colors">Copy all codes</button>
        </div>
      )}
    </div>
  );
}
`
}

func authUIProfileAvatar01() string {
	return `// profile-avatar-01 — User profile avatar with upload (Clerk / GitHub-inspired)
"use client";
import { useRef, useState } from "react";
export function ProfileAvatar({
  name = "Sarah Chen",
  email = "sarah@example.com",
  avatarUrl,
  onUpload,
}: {
  name?: string; email?: string; avatarUrl?: string;
  onUpload?: (file: File) => void;
}) {
  const [preview, setPreview] = useState(avatarUrl ?? "");
  const inputRef = useRef<HTMLInputElement>(null);
  const initials = name.split(" ").map((n) => n[0]).join("").toUpperCase().slice(0, 2);
  const handleFile = (file: File | undefined) => {
    if (!file) return;
    const url = URL.createObjectURL(file);
    setPreview(url);
    onUpload?.(file);
  };
  return (
    <div className="flex items-center gap-4">
      <div className="relative group cursor-pointer" onClick={() => inputRef.current?.click()}>
        <div className="w-20 h-20 rounded-full overflow-hidden bg-accent/20 flex items-center justify-center border-2 border-border">
          {preview ? (
            <img src={preview} alt={name} className="w-full h-full object-cover" />
          ) : (
            <span className="text-2xl font-bold text-accent">{initials}</span>
          )}
        </div>
        <div className="absolute inset-0 rounded-full bg-background/70 flex items-center justify-center opacity-0 group-hover:opacity-100 transition-opacity">
          <svg className="w-6 h-6 text-foreground" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5} d="M3 9a2 2 0 012-2h.93a2 2 0 001.664-.89l.812-1.22A2 2 0 0110.07 4h3.86a2 2 0 011.664.89l.812 1.22A2 2 0 0018.07 7H19a2 2 0 012 2v9a2 2 0 01-2 2H5a2 2 0 01-2-2V9z M15 13a3 3 0 11-6 0 3 3 0 016 0z" /></svg>
        </div>
        <input ref={inputRef} type="file" accept="image/*" className="hidden" onChange={(e) => handleFile(e.target.files?.[0])} />
      </div>
      <div>
        <p className="text-sm font-semibold text-foreground">{name}</p>
        <p className="text-xs text-text-muted">{email}</p>
        <button onClick={() => inputRef.current?.click()} className="text-xs text-accent hover:underline mt-1.5">Upload photo</button>
        {preview && <button onClick={() => { setPreview(""); }} className="text-xs text-text-muted hover:text-danger ml-2">Remove</button>}
      </div>
    </div>
  );
}
`
}

func authUIAccountDeletion01() string {
	return `// account-deletion-01 — Account deletion confirmation (GitHub / Linear-inspired)
"use client";
import { useState } from "react";
export function AccountDeletion({
  workspaceName = "Acme Corp",
  onDelete,
}: {
  workspaceName?: string; onDelete?: () => void;
}) {
  const [input, setInput] = useState("");
  const confirmed = input === workspaceName;
  return (
    <div className="border border-danger/40 bg-danger/5 rounded-xl p-5 flex flex-col gap-4">
      <div className="flex items-start gap-3">
        <svg className="w-5 h-5 text-danger flex-shrink-0 mt-0.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
        </svg>
        <div>
          <h3 className="text-sm font-semibold text-danger">Delete workspace</h3>
          <p className="text-xs text-text-secondary mt-0.5">This action cannot be undone. All projects, data, and team members will be permanently removed.</p>
        </div>
      </div>
      <div className="flex flex-col gap-1.5">
        <label className="text-xs text-text-secondary">
          Type <strong className="text-foreground font-mono">{workspaceName}</strong> to confirm
        </label>
        <input
          value={input} onChange={(e) => setInput(e.target.value)}
          placeholder={workspaceName}
          className="bg-background border border-border rounded-lg px-3 py-2 text-sm text-foreground outline-none focus:border-danger transition-colors"
        />
      </div>
      <button
        onClick={onDelete} disabled={!confirmed}
        className="w-full py-2.5 rounded-lg bg-danger hover:bg-danger/90 text-white font-semibold text-sm transition-colors disabled:opacity-40 disabled:cursor-not-allowed"
      >
        Delete {workspaceName}
      </button>
    </div>
  );
}
`
}

func authUIOAuthButton01() string {
	return `// oauth-button-01 — OAuth provider sign-in buttons (Clerk-inspired)
export function OAuthButton({
  provider = "github",
  label,
  onClick,
}: {
  provider?: "github" | "google" | "discord" | "twitter";
  label?: string; onClick?: () => void;
}) {
  const config: Record<string, { label: string; color: string; icon: JSX.Element }> = {
    github: {
      label: "Continue with GitHub",
      color: "hover:bg-[#24292e] hover:text-white hover:border-[#24292e]",
      icon: <svg className="w-5 h-5" fill="currentColor" viewBox="0 0 24 24"><path d="M12 0C5.37 0 0 5.37 0 12c0 5.31 3.435 9.795 8.205 11.385.6.105.825-.255.825-.57 0-.285-.015-1.23-.015-2.235-3.015.555-3.795-.735-4.035-1.41-.135-.345-.72-1.41-1.23-1.695-.42-.225-1.02-.78-.015-.795.945-.015 1.62.87 1.845 1.23 1.08 1.815 2.805 1.305 3.495.99.105-.78.42-1.305.765-1.605-2.67-.3-5.46-1.335-5.46-5.925 0-1.305.465-2.385 1.23-3.225-.12-.3-.54-1.53.12-3.18 0 0 1.005-.315 3.3 1.23.96-.27 1.98-.405 3-.405s2.04.135 3 .405c2.295-1.56 3.3-1.23 3.3-1.23.66 1.65.24 2.88.12 3.18.765.84 1.23 1.905 1.23 3.225 0 4.605-2.805 5.625-5.475 5.925.435.375.81 1.095.81 2.22 0 1.605-.015 2.895-.015 3.3 0 .315.225.69.825.57A12.02 12.02 0 0024 12c0-6.63-5.37-12-12-12z"/></svg>,
    },
    google: {
      label: "Continue with Google",
      color: "hover:bg-[#4285f4] hover:text-white hover:border-[#4285f4]",
      icon: <svg className="w-5 h-5" viewBox="0 0 24 24"><path fill="#4285F4" d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z"/><path fill="#34A853" d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z"/><path fill="#FBBC05" d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z"/><path fill="#EA4335" d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z"/></svg>,
    },
    discord: {
      label: "Continue with Discord",
      color: "hover:bg-[#5865F2] hover:text-white hover:border-[#5865F2]",
      icon: <svg className="w-5 h-5" fill="currentColor" viewBox="0 0 24 24"><path d="M20.317 4.37a19.791 19.791 0 00-4.885-1.515.074.074 0 00-.079.037c-.21.375-.444.864-.608 1.25a18.27 18.27 0 00-5.487 0 12.64 12.64 0 00-.617-1.25.077.077 0 00-.079-.037A19.736 19.736 0 003.677 4.37a.07.07 0 00-.032.027C.533 9.046-.32 13.58.099 18.057a.082.082 0 00.031.057 19.9 19.9 0 005.993 3.03.078.078 0 00.084-.028c.462-.63.874-1.295 1.226-1.994a.076.076 0 00-.041-.106 13.107 13.107 0 01-1.872-.892.077.077 0 01-.008-.128 10.2 10.2 0 00.372-.292.074.074 0 01.077-.01c3.928 1.793 8.18 1.793 12.062 0a.074.074 0 01.078.01c.12.098.246.198.373.292a.077.077 0 01-.006.127 12.299 12.299 0 01-1.873.892.077.077 0 00-.041.107c.36.698.772 1.362 1.225 1.993a.076.076 0 00.084.028 19.839 19.839 0 006.002-3.03.077.077 0 00.032-.054c.5-5.177-.838-9.674-3.549-13.66a.061.061 0 00-.031-.03z"/></svg>,
    },
    twitter: {
      label: "Continue with X",
      color: "hover:bg-black hover:text-white hover:border-black",
      icon: <svg className="w-5 h-5" fill="currentColor" viewBox="0 0 24 24"><path d="M18.244 2.25h3.308l-7.227 8.26 8.502 11.24H16.17l-4.714-6.231-5.401 6.231H2.747l7.73-8.835L1.254 2.25H8.08l4.259 5.63L18.245 2.25zm-1.161 17.52h1.833L7.084 4.126H5.117z"/></svg>,
    },
  };
  const c = config[provider] ?? config["github"];
  return (
    <button onClick={onClick} className={"w-full flex items-center justify-center gap-3 py-2.5 px-4 rounded-xl border border-border bg-bg-elevated text-text-secondary text-sm font-medium transition-all " + c.color}>
      {c.icon}
      {label ?? c.label}
    </button>
  );
}
`
}
