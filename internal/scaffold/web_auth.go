package scaffold

// Phase 4 of PLAN_FORMS_AND_SHARING.md — `grit add web-auth`.
//
// The web app already ships with /login, /register, /forgot-password
// pages and a useMe() hook. What's missing is a way to mark certain
// customer-facing pages as protected. AddWebAuth scaffolds the two
// pieces that close that gap:
//
//   apps/web/middleware.ts                       — SSR cookie redirect
//   apps/web/components/ProtectedWebRoute.tsx    — client-side guard
//
// Both files are idempotent — re-running grit add web-auth without
// --force skips files that exist and prints a notice.

import (
	"fmt"
	"os"
	"path/filepath"
)

// AddWebAuth runs the web-auth scaffold against an existing Grit
// project. The web app at apps/web/ is expected to exist; the
// command is a no-op otherwise (with a clear error).
func AddWebAuth(root string, force bool) error {
	webRoot := filepath.Join(root, "apps", "web")
	if _, err := os.Stat(webRoot); err != nil {
		return fmt.Errorf("apps/web not found at %s — this command needs a project with a web frontend", webRoot)
	}

	files := []struct {
		path    string
		content string
	}{
		{filepath.Join(webRoot, "middleware.ts"), webMiddlewareTS()},
		{filepath.Join(webRoot, "components", "ProtectedWebRoute.tsx"), webProtectedRouteTSX()},
	}

	for _, f := range files {
		rel, _ := filepath.Rel(root, f.path)

		if _, err := os.Stat(f.path); err == nil && !force {
			fmt.Printf("  • skipped %s (already exists — pass --force to overwrite)\n", rel)
			continue
		}

		if err := os.MkdirAll(filepath.Dir(f.path), 0755); err != nil {
			return fmt.Errorf("creating directory: %w", err)
		}
		if err := os.WriteFile(f.path, []byte(f.content), 0644); err != nil {
			return fmt.Errorf("writing %s: %w", rel, err)
		}
		fmt.Printf("  ✓ wrote %s\n", rel)
	}

	fmt.Println()
	fmt.Println("  Next steps:")
	fmt.Println("    1. Open apps/web/middleware.ts and add protected paths to the matcher.")
	fmt.Println("    2. Or wrap a page client-side with <ProtectedWebRoute>.")
	fmt.Println("    3. See /docs/concepts/protecting-web-pages for both patterns.")
	fmt.Println()
	return nil
}

// webMiddlewareTS — runs on every Next.js request. Cheap: checks for
// the existence of the grit_access HttpOnly cookie. If the cookie is
// missing AND the path matches a protected pattern, redirect to /login
// with the original path as a `next` query param. Otherwise pass
// through.
//
// Why "existence" not "validity"? The middleware can't verify a JWT
// without a network call (it'd add latency to every page request). A
// missing cookie is a guaranteed unauth state. An invalid/expired
// cookie still gets bounced — by the API itself when the page makes
// its first authenticated call. The middleware just stops the
// "blank page, then 401, then redirect" flicker.
func webMiddlewareTS() string {
	return `import { NextResponse, type NextRequest } from "next/server";

// v3.31.22 — apps/web/middleware.ts
//
// SSR cookie gate for protected pages. Edit PROTECTED_PATHS to add
// routes that require sign-in. Anything not listed here remains
// public, regardless of cookie state.
//
// v3.31.42 update: the gate now reads grit_web_session instead of
// grit_access. grit_access is set by the API on the API origin
// (e.g. localhost:8080) and is the same cookie the admin app uses
// after its own login — which meant an admin who's signed in via
// apps/admin could walk straight into apps/web's protected pages
// in the same browser. grit_web_session is set by the web app's
// own login/register flow on the WEB origin, so admin-only
// sessions don't unlock the web's gates. See lib/web-session.ts
// for the full rationale.
//
// What this DOESN'T do:
//   - Verify any JWT (would add an API round-trip to every request).
//     Forged markers without a valid API session still get bounced —
//     useMe() returns null after a 401 from /api/auth/me and
//     ProtectedWebRoute forwards to /login.
//   - Touch any cookie itself. Login/refresh/logout flows own the
//     marker; this middleware only reads.

// Add paths here that require an authenticated visitor.
const PROTECTED_PATHS: string[] = [
  "/account",
  "/account/:path*",
  // "/checkout",
  // "/dashboard",
];

// Paths that should be inaccessible to ALREADY-signed-in users
// (the login form, sign-up). Sends them to /account instead.
const AUTH_PATHS: string[] = [
  "/login",
  "/register",
  "/forgot-password",
];

function matchesAny(pathname: string, patterns: string[]): boolean {
  for (const pat of patterns) {
    if (pat === pathname) return true;
    // Naive prefix match for ":path*" suffixes.
    if (pat.endsWith(":path*")) {
      const prefix = pat.slice(0, -":path*".length);
      if (pathname === prefix.replace(/\/$/, "") || pathname.startsWith(prefix)) {
        return true;
      }
    }
  }
  return false;
}

export function middleware(request: NextRequest) {
  const { pathname } = request.nextUrl;
  const hasSession = request.cookies.has("grit_web_session");

  if (!hasSession && matchesAny(pathname, PROTECTED_PATHS)) {
    const url = request.nextUrl.clone();
    url.pathname = "/login";
    url.searchParams.set("next", pathname);
    return NextResponse.redirect(url);
  }

  if (hasSession && matchesAny(pathname, AUTH_PATHS)) {
    const url = request.nextUrl.clone();
    url.pathname = "/account";
    url.search = "";
    return NextResponse.redirect(url);
  }

  return NextResponse.next();
}

// Matcher: limit middleware execution to the paths it might act on
// — saves Next.js work on every static asset request. Keep this
// in sync with PROTECTED_PATHS + AUTH_PATHS above.
export const config = {
  matcher: [
    "/account/:path*",
    "/login",
    "/register",
    "/forgot-password",
  ],
};
`
}

// webProtectedRouteTSX — client-side guard. Use when a page can't
// be expressed as a middleware match (e.g. role-gated content where
// the role isn't carried in the cookie, or per-page-flag like
// "owner-only"). Shows a spinner while useMe() probes, redirects on
// null user, renders children otherwise.
func webProtectedRouteTSX() string {
	return `"use client";

// v3.31.22 — apps/web/components/ProtectedWebRoute.tsx
//
// Client-side route guard. Wrap a page's contents to enforce
// authentication. Pairs with apps/web/middleware.ts — middleware
// handles the SSR cookie check; this component handles the
// post-hydration "is the cookie actually valid?" probe.
//
// Use this when you need:
//   - Role/permission checks beyond cookie existence (the cookie
//     doesn't carry role; useMe() returns the full user payload).
//   - Per-page protection without editing middleware.ts.
//
// For simple "is the visitor signed in?" pages, middleware.ts is
// faster — it bounces unauthenticated requests before the page
// ever loads.

import { useEffect, type ReactNode } from "react";
import { useRouter, usePathname } from "next/navigation";
import { useMe } from "@/hooks/use-auth";

interface ProtectedWebRouteProps {
  children: ReactNode;
  // Optional role gate. When set, the user's role must match (or
  // include, for arrays). Visitors with a session but the wrong
  // role get bounced to / (the landing page).
  roles?: string | string[];
  // Where to send unauthenticated visitors. Defaults to /login with
  // a ?next= param so they return to the original page after login.
  loginPath?: string;
  // Custom loading view while useMe() is mid-flight.
  fallback?: ReactNode;
}

export function ProtectedWebRoute({
  children,
  roles,
  loginPath = "/login",
  fallback,
}: ProtectedWebRouteProps) {
  const { data: user, isLoading, isError } = useMe();
  const router = useRouter();
  const pathname = usePathname();

  useEffect(() => {
    if (isLoading) return;
    if (isError || user === null) {
      const next = encodeURIComponent(pathname ?? "/");
      router.replace(loginPath + "?next=" + next);
      return;
    }
    if (roles) {
      const allowed = Array.isArray(roles) ? roles : [roles];
      // user.role is the source of truth — falls back to a friendlier
      // redirect (to the landing page) for authenticated-but-wrong-role
      // visitors, since a /login bounce wouldn't help them.
      const userWithRole = user as { role?: string };
      if (!userWithRole.role || !allowed.includes(userWithRole.role)) {
        router.replace("/");
      }
    }
  }, [isLoading, isError, user, roles, router, pathname, loginPath]);

  if (isLoading || !user) {
    return (
      fallback ?? (
        <div className="flex min-h-screen items-center justify-center">
          <div className="h-8 w-8 animate-spin rounded-full border-2 border-slate-300 border-t-slate-900" />
        </div>
      )
    );
  }

  return <>{children}</>;
}
`
}
