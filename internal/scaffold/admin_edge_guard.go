package scaffold

// Review 2.15: the standalone admin panel had no edge protection.
//
// A signed-out visitor who typed /system/api-keys was served the whole page.
// The React tree then mounted, asked /api/auth/me, got a 401 and replaced the
// route with /login. Everything behind it was refused by the API, so nothing
// leaked, but the panel downloaded and ran for somebody with no business in
// it, the route names and the sidebar's shape were readable in the response,
// and the first paint was the shell of a dashboard they could not use.
//
// The embedded panel in a double has had a gate since v3.391.0 (adminGate, in
// deploy_hardening_repair.go) because the web app it lives inside has a
// middleware to put it in. apps/admin had no middleware at all.
//
// The gate here reads two cookies, neither of which carries a secret:
//
//   grit_admin_session  set by the admin's own api-client on the admin origin
//                       after any response that carries a session. Works
//                       wherever the panel is hosted, including an admin on
//                       its own domain and an API on another.
//   grit_signed_in      set by the API beside its session cookies. Reaches the
//                       admin only when the two share a host, which is the
//                       development default and the single-domain deploy. It
//                       is accepted so that a project upgrading into this
//                       change does not bounce everybody who is already
//                       signed in to the login page once.
//
// Both are forgeable from devtools, and that is fine: this is not the
// authorisation boundary. It stops a stranger from downloading the panel. The
// API checks the real session on every request, and that is the check that
// decides anything.

// adminSessionLib is apps/admin/lib/admin-session.ts: the cookie's name, and
// the two functions that write it. proxy.ts imports the name from here
// rather than repeating it, so the gate and the thing that opens it cannot
// drift apart.
func adminSessionLib() string {
	return `// grit_admin_session is a marker cookie, with no secret in it, written on
// the ADMIN origin whenever the API hands back a session. proxy.ts reads
// it to tell a signed-in browser from a stranger before any JavaScript runs.
//
// Why a cookie of our own when the API already sets grit_access? Because the
// API sets it on the API's origin. In development that is localhost:8080 and
// the browser does send it to localhost:3001, but an admin on its own domain
// with the API on another never sees it. A marker on this origin works in both
// shapes.
//
// It is not a credential. Anyone can write it from the console, and all that
// buys them is the HTML of a panel whose every request the API refuses. The
// session itself is the HttpOnly pair the API sets, which JavaScript cannot
// read and this file never touches.

export const ADMIN_SESSION_COOKIE = "grit_admin_session";

// Set by the API beside its session cookies. Only arrives when the API shares
// this app's host; proxy.ts accepts it as well so an already-signed-in
// browser is not sent to the login page the first time it loads a panel that
// has just gained this gate.
export const API_SIGNED_IN_COOKIE = "grit_signed_in";

// Seven days, which is the refresh window. The cookie outliving the session is
// harmless: the page loads, /api/auth/me answers 401, and the clear below runs.
const MAX_AGE_SECONDS = 7 * 24 * 60 * 60;

export function setAdminSession(): void {
  if (typeof document === "undefined") return;
  const secure =
    typeof location !== "undefined" && location.protocol === "https:" ? "; Secure" : "";
  document.cookie =
    ADMIN_SESSION_COOKIE +
    "=1; Path=/; Max-Age=" +
    MAX_AGE_SECONDS +
    "; SameSite=Lax" +
    secure;
}

export function clearAdminSession(): void {
  if (typeof document === "undefined") return;
  document.cookie = ADMIN_SESSION_COOKIE + "=; Path=/; Max-Age=0; SameSite=Lax";
}

export function hasAdminSession(): boolean {
  if (typeof document === "undefined") return false;
  return document.cookie
    .split("; ")
    .some((c) => c.startsWith(ADMIN_SESSION_COOKIE + "="));
}
`
}

// adminProxyTS is apps/admin/proxy.ts.
//
// proxy.ts, not middleware.ts: Next 16 renamed the convention and prints a
// deprecation notice for the old name on every dev start. Same file, same
// matcher, and the function it exports is called proxy.
//
// The matcher denies by default: every path except Next's own output and
// anything that looks like a file. A page added later, by grit generate or by
// a plugin, is protected the day it lands rather than the day somebody
// remembers to add it to a list. The auth pages are let through inside the
// function, where the list of them is short and in one place.
func adminProxyTS() string {
	return `import { NextResponse, type NextRequest } from "next/server";
import { ADMIN_SESSION_COOKIE, API_SIGNED_IN_COOKIE } from "@/lib/admin-session";

// The pages a signed-out visitor is supposed to reach: the (auth) route
// group. Everything else in the panel needs a session.
//
// Keep this in step with app/(auth)/. A page missing from here is a page
// nobody can open, which is the one failure of this file that is not
// self-correcting, so there is a test on it: TestEveryAuthPageIsPublic in
// the framework's own suite derives this list from the panel's file map.
const PUBLIC_PATHS = [
  "/login",
  "/sign-up",
  "/forgot-password",
  "/reset-password",
  "/verify-email",
  "/magic-link",
  "/callback",
];

function isPublic(pathname: string): boolean {
  return PUBLIC_PATHS.some((p) => pathname === p || pathname.startsWith(p + "/"));
}

export function proxy(request: NextRequest) {
  const { pathname } = request.nextUrl;
  if (isPublic(pathname)) return NextResponse.next();

  // Neither cookie is a credential, and neither is checked against anything.
  // They answer "has this browser signed in here?", which is enough to decide
  // whether to send the panel. Whether the session is still valid is the API's
  // answer, on every request, and the page asks it as soon as it mounts.
  if (
    request.cookies.has(ADMIN_SESSION_COOKIE) ||
    request.cookies.has(API_SIGNED_IN_COOKIE)
  ) {
    return NextResponse.next();
  }

  const url = request.nextUrl.clone();
  url.pathname = "/login";
  url.search = "";
  return NextResponse.redirect(url);
}

export const config = {
  // Everything but Next's own output and anything with a file extension. The
  // login page's logo is served from public/, so it has to stay reachable
  // without a session or the page a stranger is sent to renders without it.
  // [.] rather than a backslash: this string is read by JavaScript before it
  // is read as a regex, so a single backslash here would be eaten and the
  // pattern would match any character where it means a dot. A character class
  // says the same thing and survives every layer it passes through.
  matcher: ["/((?!_next/|[^?]*[.][A-Za-z0-9]+$).*)"],
};
`
}
