// v3.31.42 -- grit_web_session is a non-HttpOnly marker cookie set
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
