import { createUploader, createAxiosTransport } from "@repo/upload";
import { optimizeImage } from "@repo/upload/web";
import { createApiClient } from "@/lib/api-core";

// The API's address, the /api/v1 rewrite, the CSRF header and the idempotency
// key are in lib/api-core.ts, whose copy in the web app is the same file. What
// is left here is what only the admin needs: the public-IP hint and the
// 401-refresh retry.
export { API_URL, API_VERSION, apiUrl, versionedPath } from "@/lib/api-core";

// Auth storage policy (Grit 3.27+):
//   - The API issues HttpOnly grit_access + grit_refresh cookies on
//     login / register / refresh / OAuth callback. The browser stores
//     them automatically; JS never reads or writes the access token.
//   - withCredentials: true tells axios to attach those cookies on every
//     request, including cross-origin dev (admin :3001 → api :8080).
//   - The CSRF token rides a NON-HttpOnly grit_csrf cookie. We echo it
//     into X-CSRF-Token on every state-changing method — the API's
//     AutoCSRF middleware requires that double-submit token for the
//     mutation to pass.
//   - The 401-refresh interceptor below POSTS /api/auth/refresh with no
//     body — the API reads grit_refresh from the cookie and issues a
//     new grit_access via Set-Cookie. JS still never sees a token.
export const apiClient = createApiClient();

// v3.31.49 -- public-IP hint. When the operator runs the admin on
// localhost (the default), the API sees the TCP peer as ::1 and
// logs that in the activity feed. The browser fetches its public IP
// once (cached in sessionStorage for the tab's lifetime) and
// attaches it as X-Public-IP-Hint. The API uses it only when the
// observed peer is loopback -- production traffic from real proxies
// keeps using the trusted X-Forwarded-For path and never honours
// this hint, so it can't be used to spoof audit records.
let publicIPCache: string | null = null;
async function getPublicIPHint(): Promise<string | null> {
  if (publicIPCache) return publicIPCache;
  if (typeof window === "undefined") return null;
  const cached = window.sessionStorage.getItem("grit_public_ip");
  if (cached) {
    publicIPCache = cached;
    return cached;
  }
  try {
    const res = await fetch("https://api.ipify.org?format=json", {
      cache: "no-store",
    });
    if (!res.ok) return null;
    const data = (await res.json()) as { ip?: string };
    if (data.ip) {
      publicIPCache = data.ip;
      window.sessionStorage.setItem("grit_public_ip", data.ip);
      return data.ip;
    }
  } catch {
    // Offline / blocked by an ad-blocker -- fall through; the API
    // will log "::1" as it does today.
  }
  return null;
}
// Kick off the lookup eagerly so the cache is warm by the time the
// first request fires. Fire-and-forget; failures are silent.
//
// Dev only, for two reasons: the hint is only ever honoured when the API sees
// a loopback peer (production goes through X-Forwarded-For), and a production
// build should not reach out to a third party on every page load.
if (typeof window !== "undefined" && process.env.NODE_ENV !== "production") {
  void getPublicIPHint();
}

// The CSRF header and the Idempotency-Key come from createApiClient. This one
// is the admin's alone.
apiClient.interceptors.request.use((config) => {
  // v3.31.49 -- attach the cached public-IP hint when we have one.
  if (publicIPCache && config.headers) {
    config.headers["X-Public-IP-Hint"] = publicIPCache;
  }
  return config;
});

let isRefreshing = false;
let failedQueue: Array<{
  resolve: (value: unknown) => void;
  reject: (reason: unknown) => void;
}> = [];

const processQueue = (error: unknown) => {
  failedQueue.forEach((promise) => {
    if (error) {
      promise.reject(error);
    } else {
      promise.resolve(undefined);
    }
  });
  failedQueue = [];
};

apiClient.interceptors.response.use(
  (response) => response,
  async (error) => {
    const originalRequest = error.config;

    // Skip refresh on the auth endpoints themselves — a wrong password
    // 401-ing into a refresh attempt would loop and wipe the session.
    const url = originalRequest?.url || "";
    const isAuthEndpoint =
      url.includes("/auth/login") ||
      url.includes("/auth/register") ||
      url.includes("/auth/refresh") ||
      url.includes("/auth/me");

    if (error.response?.status === 401 && !originalRequest._retry && !isAuthEndpoint) {
      if (isRefreshing) {
        return new Promise((resolve, reject) => {
          failedQueue.push({ resolve, reject });
        }).then(() => apiClient(originalRequest));
      }

      originalRequest._retry = true;
      isRefreshing = true;

      try {
        // Empty body — the API reads grit_refresh from the HttpOnly
        // cookie that the browser attached automatically, and issues a
        // new grit_access via the Set-Cookie response header.
        await apiClient.post("/api/auth/refresh");

        processQueue(null);
        return apiClient(originalRequest);
      } catch (refreshError) {
        processQueue(refreshError);
        // The cookies are HttpOnly so we can't expire them from JS.
        // Forcing a navigation to /login lets the user re-authenticate;
        // the next successful login replaces the cookies via Set-Cookie.
        if (typeof window !== "undefined") {
          window.location.href = "/login";
        }
        return Promise.reject(refreshError);
      } finally {
        isRefreshing = false;
      }
    }

    return Promise.reject(error);
  }
);

/**
 * The shared uploader: optimise on the client, then upload straight to storage.
 *
 * Built once at module scope so the profile list is fetched once per session
 * rather than per file.
 */
const uploader = createUploader({
  // "/api" here, not "/api/v1": the interceptor above pins the version, so the
  // package never has to know it exists.
  transport: createAxiosTransport(apiClient, "/api"),
  optimize: optimizeImage,
});

/**
 * Upload a file, optimising it first.
 *
 * The bytes go browser to storage through a presigned URL and never touch the
 * API. Optimising before the PUT is what makes that worth having: a 6 MB photo
 * is about 35 KB by the time it is sent, so the upload is faster, the bucket
 * is smaller, and on mobile data the difference is most of the wait.
 *
 * The signature is unchanged, so every existing call site keeps working. What
 * changed is that onProgress now covers the thumbnail as well as the primary,
 * and the returned data carries the renditions.
 */
export async function uploadFile(
  file: File,
  endpoint = "/api/uploads",
  onProgress?: (percent: number) => void,
  accepts?: string[]
): Promise<{ data: Record<string, unknown>; message: string }> {
  // The endpoint carries the field's metadata as query params — see
  // buildUploadEndpoint(). Read them back out instead of dropping them:
  // without accepts, a field declared file:zip could never upload (the API
  // falls back to its global allow-list) and a field declared file:pdf would
  // happily take a PNG.
  let fieldAccepts = accepts;
  let profile: string | undefined;
  if (endpoint.includes("?")) {
    const q = new URLSearchParams(endpoint.slice(endpoint.indexOf("?") + 1));
    if (!fieldAccepts?.length) {
      const a = q.get("accepts");
      if (a) fieldAccepts = a.split(",").map((x) => x.trim()).filter(Boolean);
    }
    profile = q.get("profile") ?? undefined;
  }

  const ref = await uploader.upload(file, file.name, {
    accepts: fieldAccepts,
    profile,
    onProgress: onProgress ? (f) => onProgress(Math.round(f * 100)) : undefined,
  });

  return {
    data: ref as unknown as Record<string, unknown>,
    message: "File uploaded successfully",
  };
}
