"use client";

// marketing-navbar-01 — Marketing site navbar (Vercel.com / Linear.app-inspired)
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

export default MarketingNavbar;
