package scaffold

import "strings"

// ─── M38: one module owns the API's address ─────────────────────────────────
//
// The web app's lib/api.ts and the admin panel's lib/api-client.ts were two
// copies of the same client: the same /api/v1 rewrite, the same CSRF
// double-submit header, the same Idempotency-Key. Two copies of a rule drift,
// and these had: only one of them skipped /api/ws.
//
// Worse than the duplication was what surrounded it. The expression
//
//	process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080"
//
// was written out in eleven files, and the ones that used it built their URL by
// hand: an <a href> to the OAuth redirect, the public form page's axios calls,
// the SAML URLs in the SSO screen. None of those went through the client, so
// none of them got the version prefix, and they reached the API's transitional
// unversioned alias instead. That alias exists for old clients. The day it goes,
// public forms, social sign-in and every SSO connection go with it, and nothing
// says so at build time.
//
// So: one module per app, with the same contents in both, exporting API_URL,
// apiUrl(path) for code that cannot use the client, and createApiClient() for
// the clients themselves.

// apiCoreTS is lib/api-core.ts, identical in the web app and the admin panel.
func apiCoreTS() string {
	src := `import axios, { type AxiosInstance, type CreateAxiosDefaults } from "axios";

// The API's address, written down once.
//
// Every module that needs it imports it from here. Before this, the same
// environment-variable expression was copied into eleven files, so "where does
// the frontend think the API is" had eleven answers and changing the fallback
// changed one of them.
export const API_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";

// The API is served under a version prefix (/api/v1/...). Endpoints are written
// as "/api/..." throughout the app and pinned to the version here, so moving to
// v2 is a one-line change instead of a find-and-replace across ~200 strings.
export const API_VERSION = "v1";

// versionedPath puts an "/api/..." path under the current version.
//
// Idempotent: a path that already names the version comes back unchanged, so a
// retried request cannot end up at /api/v1/v1/users. /api/ws is left alone
// because a WebSocket upgrade is not part of the versioned REST surface.
export function versionedPath(path: string): string {
  if (
    !path.startsWith("/api/") ||
    path === "/api/ws" ||
    path.startsWith("/api/" + API_VERSION + "/")
  ) {
    return path;
  }
  return "/api/" + API_VERSION + path.slice("/api".length);
}

// apiUrl builds an absolute URL for the code that cannot go through the axios
// client: an <a href> the browser navigates to, a fetch in a server component,
// a form action.
//
// Use it instead of API_URL + "/api/...". Written by hand, those paths skip the
// version rewrite and land on the API's unversioned alias, which is there for
// clients you cannot update and is not a place your own frontend should be.
export function apiUrl(path: string): string {
  return API_URL + versionedPath(path);
}

// createApiClient builds the axios instance this app talks to the API with.
//
// Both frontends call it, so the version rewrite, the CSRF header and the
// idempotency key are defined once. Anything app-specific (the admin's
// 401-refresh retry, its public-IP hint) is added by the caller afterwards.
export function createApiClient(config: CreateAxiosDefaults = {}): AxiosInstance {
  const client = axios.create({
    baseURL: API_URL,
    headers: { "Content-Type": "application/json" },
    // The browser attaches the HttpOnly grit_access / grit_refresh cookies set
    // by /api/auth/login automatically. Without this, axios skips them on
    // cross-origin requests in dev (api on :8080, web on :3000) and the server
    // treats every request as anonymous.
    withCredentials: true,
    ...config,
  });

  client.interceptors.request.use((cfg) => {
    if (cfg.url) cfg.url = versionedPath(cfg.url);
    return cfg;
  });

  client.interceptors.request.use((cfg) => {
    // Echo the grit_csrf cookie into X-CSRF-Token on every state-changing
    // request. The cookie is intentionally not HttpOnly: it is the
    // double-submit token, paired with the cookie the AutoCSRF middleware
    // checks. Safe methods don't need it; the middleware skips them and issues
    // or refreshes the cookie as a side effect.
    if (typeof document !== "undefined") {
      const m = document.cookie.match(/(?:^|; )grit_csrf=([^;]+)/);
      if (m && cfg.headers) {
        cfg.headers["X-CSRF-Token"] = decodeURIComponent(m[1]);
      }
    }

    // Auto-attach Idempotency-Key on unsafe methods so any mutation gets
    // safe-retry semantics for free. A retry that replays the same config
    // object reuses the key, and the server answers with the first 2xx it
    // cached for (method, path, key).
    const method = (cfg.method || "get").toUpperCase();
    const unsafe =
      method === "POST" || method === "PUT" || method === "PATCH" || method === "DELETE";
    if (unsafe && cfg.headers && !cfg.headers["Idempotency-Key"]) {
      cfg.headers["Idempotency-Key"] = crypto.randomUUID();
    }
    return cfg;
  });

  return client;
}

// getApiErrorMessage pulls the message out of whatever a failed request threw.
//
// The API answers an error as { error: { code, message, details } }, so the
// message is four levels down an axios error, and every page that wanted it
// wrote the same cast:
//
//	(err as { response?: { data?: { error?: { message?: string } } } })
//	  ?.response?.data?.error?.message
//
// That appeared 27 times in 21 files, in at least four spellings, two of them
// through ~any~. Here instead, with the fallback the caller wants when the
// server said nothing useful: a network failure, a 502 from a proxy, a thrown
// string.
export function getApiErrorMessage(err: unknown, fallback: string): string {
  const message = (err as { response?: { data?: { error?: { message?: string } } } })
    ?.response?.data?.error?.message;
  if (typeof message === "string" && message.trim()) return message;
  if (err instanceof Error && err.message) return err.message;
  return fallback;
}
`
	return strings.ReplaceAll(src, "~", "`")
}
