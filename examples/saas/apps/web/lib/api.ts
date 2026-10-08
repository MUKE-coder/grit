// The client. Its address, the /api/v1 rewrite, the CSRF header and the
// idempotency key all live in lib/api-core.ts, which the admin panel uses too,
// so the two clients cannot drift apart again.
import type { AxiosError, InternalAxiosRequestConfig } from "axios";
import { createApiClient } from "@/lib/api-core";

export const api = createApiClient();

// A signed-in session outlives its access token: when a request comes back
// 401, the refresh cookie buys a new one and the request is sent again. The
// web app had no such retry, so every signed-in page failed fifteen minutes
// in, when the access token expired.
//
// Requests that fail together share one refresh. The server rotates refresh
// tokens and treats a second use of the old one as theft, so ten refreshes for
// ten failed requests would sign the person out everywhere.
let refreshing: Promise<void> | null = null;

interface RetriableConfig extends InternalAxiosRequestConfig {
  _retried?: boolean;
}

api.interceptors.response.use(undefined, async (error: AxiosError) => {
  const config = error.config as RetriableConfig | undefined;
  // Sign-in, sign-up, refresh and sign-out answer 401 for reasons a refresh
  // cannot fix.
  const noRetry = /\/auth\/(login|register|refresh|logout)/.test(config?.url ?? "");
  if (error.response?.status !== 401 || !config || config._retried || noRetry) {
    throw error;
  }
  config._retried = true;
  refreshing ??= api
    .post("/api/auth/refresh")
    .then(() => undefined)
    .finally(() => {
      refreshing = null;
    });
  await refreshing;
  return api.request(config);
});

// Re-exported so "@/lib/api" stays the one import a page needs, whether it
// wants the client or just the URL.
export { API_URL, API_VERSION, apiUrl, versionedPath } from "@/lib/api-core";

// v3.31.21: alias kept so generated React Query hooks that import
// { apiClient } from "@/lib/api" resolve symmetrically with apps/admin
// (which exports the same name from its own api-client.ts).
export const apiClient = api;
