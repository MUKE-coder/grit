"use client";

// app-navbar-01 — App top navbar (Vercel / Linear-inspired)
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

export default AppNavbar;
