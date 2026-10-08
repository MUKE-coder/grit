"use client";

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
        className={`flex min-h-screen flex-col transition-all duration-200 ${
          sidebarCollapsed ? "md:ml-16" : "md:ml-64"
        }`}
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
