"use client";

import { useState } from "react";
import Link from "next/link";
import { Menu, X, Github, Shield } from "lucide-react";

import { MegaMenu, MegaMenuMobile } from "@/components/mega-menu";
import { UserMenu } from "@/components/UserMenu";

const DOCS_URL = "https://gritframework.dev/docs";
// Where the admin panel is.
//
// A triple project runs it as its own app on 3001. A double has no second app:
// the panel is a route group in THIS app at /admin, and a link to
// http://localhost:3001 pointed at nothing, which is what a double's navbar did
// until v3.235.0. NEXT_PUBLIC_ADMIN_URL still overrides both.
const ADMIN_URL = process.env.NEXT_PUBLIC_ADMIN_URL || "http://localhost:3001";

// Grit's docs and repository. Links for whoever is building this app, not for
// the people visiting it, so they render in development only -- the same rule
// DevLinks follows on the landing page.
const DEV = process.env.NODE_ENV !== "production";

export function Navbar() {
  const [mobileOpen, setMobileOpen] = useState(false);

  return (
    <nav className="sticky top-0 z-50 border-b border-border/50 bg-background/80 backdrop-blur-lg">
      <div className="mx-auto flex h-16 max-w-5xl items-center justify-between px-6">
        <Link href="/" className="flex items-center gap-2.5">
          <div className="flex h-8 w-8 items-center justify-center rounded-lg bg-accent/15 border border-accent/20">
            <span className="text-accent font-mono font-bold text-sm">S</span>
          </div>
          <span className="text-lg font-bold tracking-tight">saas</span>
        </Link>

        <div className="hidden md:flex items-center gap-6">
          <MegaMenu />
          {DEV && (
            <>
              <a
                href={DOCS_URL}
                target="_blank"
                rel="noopener noreferrer"
                className="text-sm text-text-secondary hover:text-foreground transition-colors"
              >
                Docs
              </a>
              <a
                href="https://github.com/MUKE-coder/grit"
                target="_blank"
                rel="noopener noreferrer"
                aria-label="Grit on GitHub (opens in a new tab)"
                className="text-text-secondary hover:text-foreground transition-colors"
              >
                <Github className="h-5 w-5" aria-hidden="true" />
              </a>
            </>
          )}
          {/* v3.31.49 -- Admin CTA (always visible, even with auth). */}
          <a
            href={ADMIN_URL}
            target="_blank"
            rel="noopener noreferrer"
            className="inline-flex items-center gap-1.5 rounded-lg border border-border bg-bg-tertiary px-3 py-1.5 text-sm font-medium text-text-secondary hover:bg-bg-hover hover:text-foreground transition-colors"
          >
            <Shield className="h-3.5 w-3.5" />
            Admin
          </a>
          {/* grit:nav:account-desktop */}
          <UserMenu />
        </div>

        <button
          onClick={() => setMobileOpen(!mobileOpen)}
          className="md:hidden p-2 text-text-secondary hover:text-foreground transition-colors"
          aria-label="Toggle menu"
        >
          {mobileOpen ? <X className="h-5 w-5" /> : <Menu className="h-5 w-5" />}
        </button>
      </div>

      {mobileOpen && (
        <div className="md:hidden border-t border-border/50 bg-background/95 backdrop-blur-lg">
          <div className="mx-auto max-w-5xl px-6 py-4 flex flex-col gap-3">
            <MegaMenuMobile onNavigate={() => setMobileOpen(false)} />
            {DEV && (
              <>
                <a
                  href={DOCS_URL}
                  target="_blank"
                  rel="noopener noreferrer"
                  className="text-sm py-2 text-text-secondary hover:text-foreground transition-colors"
                >
                  Docs
                </a>
                <a
                  href="https://github.com/MUKE-coder/grit"
                  target="_blank"
                  rel="noopener noreferrer"
                  className="text-sm py-2 text-text-secondary hover:text-foreground transition-colors"
                >
                  GitHub
                </a>
              </>
            )}
            <a
              href={ADMIN_URL}
              target="_blank"
              rel="noopener noreferrer"
              className="inline-flex items-center gap-1.5 rounded-lg border border-border bg-bg-tertiary px-3 py-2 text-sm font-medium text-text-secondary hover:bg-bg-hover hover:text-foreground transition-colors"
            >
              <Shield className="h-3.5 w-3.5" />
              Admin
            </a>
            {/* grit:nav:account-mobile */}
            <div className="mt-2 border-t border-border/50 pt-3">
              <UserMenu />
            </div>
          </div>
        </div>
      )}
    </nav>
  );
}
