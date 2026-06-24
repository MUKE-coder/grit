package scaffold

// web_chrome_files.go — v3.31.42 scaffold templates for the web
// app's chrome surface: AppChrome wrapper, UserMenu component, and
// the web-session marker library. These three pieces work together
// to give the customer-facing web app a sane auth UX:
//
//   - AppChrome hides navbar + footer on auth and public-form-share
//     pages so they render full-bleed like the admin's auth pages.
//   - UserMenu shows Login/Sign up CTAs for anonymous visitors and
//     an avatar dropdown for signed-in users.
//   - web-session.ts manages a non-HttpOnly marker cookie on the web
//     origin so the middleware can tell admin sessions apart from
//     web sessions when both apps run on the same browser.

func webAppChrome() string {
	return `"use client";

// v3.31.42 -- AppChrome conditionally renders Navbar + Footer based on
// the current pathname. The auth pages and the public form-share page
// each have their own standalone layout (AuthShell / form layout) and
// should not double up on chrome. Rendering this inside the root
// layout keeps it server-friendly: the root layout stays a server
// component, only this small wrapper is client.

import { usePathname } from "next/navigation";
import { Navbar } from "@/components/navbar";
import { Footer } from "@/components/footer";

// Pathname prefixes that opt out of Navbar + Footer.
// - /login, /register, /forgot-password, /callback come from the (auth)
//   route group. AuthShell already provides its own full-bleed layout.
// - /forms/<token> is the public form-share page. It needs to look like
//   a stand-alone form, not part of the marketing site.
const CHROMELESS_PREFIXES = [
  "/login",
  "/register",
  "/forgot-password",
  "/callback",
  "/forms/",
];

export function AppChrome({ children }: { children: React.ReactNode }) {
  const pathname = usePathname() ?? "";
  const chromeless = CHROMELESS_PREFIXES.some((p) =>
    pathname === p || pathname.startsWith(p)
  );

  if (chromeless) {
    return <main className="min-h-screen">{children}</main>;
  }

  return (
    <>
      <Navbar />
      <main className="min-h-screen">{children}</main>
      <Footer />
    </>
  );
}
`
}

func webUserMenu() string {
	return `"use client";

// v3.31.42 -- UserMenu shows different chrome in the navbar
// depending on whether the visitor is signed in:
//   - signed in : avatar dropdown with Account + Sign out
//   - signed out: Log in + Sign up buttons
//   - loading   : a small placeholder so the navbar doesn't shift
//                 layout once useMe() resolves
//
// The signed-in / signed-out decision uses the same useMe() that
// every other auth-aware page reads, so it picks up login + logout
// immediately via React Query's cache.

import { useEffect, useRef, useState } from "react";
import Link from "next/link";
import { useMe, useLogout } from "@/hooks/use-auth";
import { ChevronDown, LogOut, User as UserIcon } from "lucide-react";

export function UserMenu() {
  const { data: user, isLoading } = useMe();
  const logout = useLogout();
  const [open, setOpen] = useState(false);
  const rootRef = useRef<HTMLDivElement>(null);

  // Close on outside click. Attached only while the dropdown is open
  // so we don't pay the cost on every page.
  useEffect(() => {
    if (!open) return;
    function handle(e: MouseEvent) {
      if (rootRef.current && !rootRef.current.contains(e.target as Node)) {
        setOpen(false);
      }
    }
    document.addEventListener("mousedown", handle);
    return () => document.removeEventListener("mousedown", handle);
  }, [open]);

  if (isLoading) {
    // Reserve roughly the same width as the signed-out CTA pair so
    // the navbar doesn't jump on auth resolve.
    return <div className="h-8 w-32 rounded-md bg-bg-elevated/40 animate-pulse" />;
  }

  if (!user) {
    return (
      <div className="flex items-center gap-2">
        <Link
          href="/login"
          className="rounded-lg border border-border bg-transparent px-3 py-1.5 text-sm font-medium text-foreground hover:bg-bg-hover transition-colors"
        >
          Log in
        </Link>
        <Link
          href="/register"
          className="rounded-lg bg-accent px-3 py-1.5 text-sm font-medium text-white hover:bg-accent-hover transition-colors"
        >
          Sign up
        </Link>
      </div>
    );
  }

  const firstName = (user as { first_name?: string }).first_name || "";
  const lastName = (user as { last_name?: string }).last_name || "";
  const email = (user as { email?: string }).email || "";
  const fullName =
    [firstName, lastName].filter(Boolean).join(" ") || email || "Account";
  const initials =
    ((firstName[0] || "") + (lastName[0] || "")).toUpperCase() ||
    email.slice(0, 2).toUpperCase() ||
    "U";

  return (
    <div ref={rootRef} className="relative">
      <button
        type="button"
        onClick={() => setOpen((v) => !v)}
        className="flex items-center gap-2 rounded-lg border border-border bg-bg-elevated/60 px-2 py-1.5 text-sm text-foreground hover:bg-bg-hover transition-colors"
      >
        <span className="inline-flex h-7 w-7 items-center justify-center rounded-full bg-accent/15 text-accent text-xs font-semibold">
          {initials}
        </span>
        <span className="hidden sm:inline max-w-[100px] truncate">{fullName}</span>
        <ChevronDown className="h-3.5 w-3.5 text-text-muted" />
      </button>

      {open && (
        <div className="absolute right-0 top-full z-50 mt-1 w-56 rounded-lg border border-border bg-bg-elevated shadow-lg">
          <div className="border-b border-border px-3 py-2.5">
            <p className="truncate text-sm font-medium text-foreground">{fullName}</p>
            {email && (
              <p className="truncate text-xs text-text-muted">{email}</p>
            )}
          </div>
          <div className="p-1">
            <Link
              href="/account"
              onClick={() => setOpen(false)}
              className="flex items-center gap-2 rounded-md px-2.5 py-2 text-sm text-foreground hover:bg-bg-hover transition-colors"
            >
              <UserIcon className="h-3.5 w-3.5 text-text-muted" />
              Account
            </Link>
            <button
              type="button"
              onClick={() => {
                setOpen(false);
                logout.mutate();
              }}
              className="flex w-full items-center gap-2 rounded-md px-2.5 py-2 text-left text-sm text-foreground hover:bg-bg-hover transition-colors"
            >
              <LogOut className="h-3.5 w-3.5 text-text-muted" />
              Sign out
            </button>
          </div>
        </div>
      )}
    </div>
  );
}
`
}

func webSessionLib() string {
	return `// v3.31.42 -- grit_web_session is a non-HttpOnly marker cookie set
// on the WEB origin (e.g. localhost:3000) after a successful login
// or registration through the web app's own auth flow.
//
// Why a separate cookie when the API already sets grit_access on
// the API origin? Because the API origin (localhost:8080) is shared
// between the admin and web apps -- the browser attaches grit_access
// to every cross-origin call to the API from either app. Without
// this marker, an admin who's logged in via apps/admin can hit any
// "protected" page in apps/web (the API call succeeds, useMe()
// returns the admin user, ProtectedWebRoute is happy). That's a
// real surprise vs the user's expectation that web-side auth and
// admin-side auth are independent.
//
// The marker fixes this without weakening the actual session
// security:
//   - Middleware checks for grit_web_session at the edge. No
//     marker on the web origin -> redirect to /login. Admins
//     who never signed in through the web app don't have it.
//   - The marker is forgeable from devtools (any client cookie is),
//     but useMe() still calls the API which validates the real
//     grit_access JWT. Forging the marker without a valid session
//     only gets you past the marker check; the API still 401s.
//
// SameSite=Lax + Path=/ matches the rest of the web app cookies and
// stays attached on top-level navigations from links.

const WEB_SESSION_COOKIE = "grit_web_session";
const MAX_AGE_SECONDS = 7 * 24 * 60 * 60; // 7 days, matches refresh window

export function setWebSessionMarker(): void {
  if (typeof document === "undefined") return;
  document.cookie =
    WEB_SESSION_COOKIE +
    "=1; Path=/; Max-Age=" +
    MAX_AGE_SECONDS +
    "; SameSite=Lax";
}

export function clearWebSessionMarker(): void {
  if (typeof document === "undefined") return;
  document.cookie = WEB_SESSION_COOKIE + "=; Path=/; Max-Age=0; SameSite=Lax";
}

export function hasWebSessionMarker(): boolean {
  if (typeof document === "undefined") return false;
  return document.cookie.split("; ").some((c) => c.startsWith(WEB_SESSION_COOKIE + "="));
}
`
}
