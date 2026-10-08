import { NextResponse, type NextRequest } from "next/server";

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

// ── The admin panel ──────────────────────────────────────────────────────
//
// The admin checks the session in the browser, so a visitor who never signed
// in downloaded every admin page before being sent to the login. The API sets
// grit_signed_in, a marker with no secret in it, beside its session cookies. A
// browser only sends it here when the web app and the API share a host
// (localhost in development, or one domain behind a proxy); on separate hosts
// it never arrives, so this gate stands aside and the page's own check still
// applies. The API authorises every request either way.
const ADMIN_PUBLIC_PATHS = [
  "/admin/login",
  "/admin/sign-up",
  "/admin/forgot-password",
  "/admin/reset-password",
  "/admin/verify-email",
  "/admin/callback",
];

function apiHostname(): string {
  try {
    return new URL(process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080").hostname;
  } catch {
    return "";
  }
}

function adminGate(request: NextRequest): NextResponse | null {
  const { pathname, hostname } = request.nextUrl;
  if (pathname !== "/admin" && !pathname.startsWith("/admin/")) return null;
  if (ADMIN_PUBLIC_PATHS.some((p) => pathname === p || pathname.startsWith(p + "/"))) return null;
  if (apiHostname() !== hostname) return null;
  if (request.cookies.has("grit_signed_in") || request.cookies.has("grit_access")) return null;
  const url = request.nextUrl.clone();
  url.pathname = "/admin/login";
  url.search = "";
  return NextResponse.redirect(url);
}

export function middleware(request: NextRequest) {
  const gate = adminGate(request);
  if (gate) return gate;

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
    "/admin",
    "/admin/:path*",
    "/account/:path*",
    "/login",
    "/register",
    "/forgot-password",
  ],
};
