import type { UploadTransport } from "./types";

/**
 * Transports: the small amount of glue between the uploader and your API
 * client.
 *
 * There are two methods that talk to your API and one that does not. The PUT
 * goes to a presigned URL, which carries its own authorisation in the query
 * string, and attaching an Authorization header to that makes S3 reject the
 * request as a signature mismatch. That is the entire reason put() is separate
 * from post() rather than another call on the same client.
 */

interface AxiosLike {
  get(url: string, config?: unknown): Promise<{ data: unknown }>;
  post(url: string, body?: unknown, config?: unknown): Promise<{ data: unknown }>;
}

/**
 * Adapts the generated apiClient (axios) used by the admin and web apps.
 *
 *   import { apiClient } from "@/lib/api-client";
 *   const transport = createAxiosTransport(apiClient, "/api");
 *
 * basePath is prepended to every API call. Grit's client writes endpoints as
 * "/api/users" and pins the version in an interceptor, so "/api" here becomes
 * "/api/v1/uploads/presign" on the wire without this package knowing the
 * version exists.
 */
export function createAxiosTransport(client: AxiosLike, basePath = ""): UploadTransport {
  const at = (path: string) => basePath + path;
  return {
    get: async <T,>(path: string) => (await client.get(at(path))).data as T,
    post: async <T,>(path: string, body: unknown) =>
      (await client.post(at(path), body)).data as T,
    put: (url, body, contentType, onProgress) =>
      putWithProgress(url, body, contentType, onProgress),
  };
}

/**
 * A transport built on fetch, for Expo and anywhere without axios.
 *
 * getToken is called per request rather than captured once, so a token
 * refreshed mid-session is picked up instead of a stale one being reused until
 * the app restarts.
 */
export function createFetchTransport(
  baseUrl: string,
  getToken?: () => string | null | Promise<string | null>,
): UploadTransport {
  const headers = async (): Promise<Record<string, string>> => {
    const h: Record<string, string> = { "Content-Type": "application/json" };
    const token = await getToken?.();
    if (token) h.Authorization = "Bearer " + token;
    return h;
  };

  const json = async <T,>(res: Response): Promise<T> => {
    if (!res.ok) {
      const body = await res.text();
      throw new Error("Request failed (" + res.status + "): " + body.slice(0, 200));
    }
    return (await res.json()) as T;
  };

  return {
    get: async <T,>(path: string) =>
      json<T>(await fetch(baseUrl + path, { headers: await headers(), credentials: "include" })),
    post: async <T,>(path: string, body: unknown) =>
      json<T>(
        await fetch(baseUrl + path, {
          method: "POST",
          headers: await headers(),
          credentials: "include",
          body: JSON.stringify(body),
        }),
      ),
    put: (url, body, contentType, onProgress) =>
      putWithProgress(url, body, contentType, onProgress),
  };
}

/**
 * PUT to a presigned URL, reporting progress.
 *
 * XMLHttpRequest rather than fetch, because fetch still cannot report upload
 * progress in any browser. On a phone uploading over mobile data, a progress
 * bar that does not move is the difference between waiting and force-quitting.
 *
 * No credentials and no Authorization header: the signature is the auth, and
 * anything extra invalidates it.
 */
function putWithProgress(
  url: string,
  body: Blob,
  contentType: string,
  onProgress?: (fraction: number) => void,
): Promise<void> {
  if (typeof XMLHttpRequest === "undefined") {
    return fetch(url, { method: "PUT", body, headers: { "Content-Type": contentType } }).then(
      (res) => {
        if (!res.ok) throw new Error("Upload failed (" + res.status + ")");
      },
    );
  }

  return new Promise<void>((resolve, reject) => {
    const xhr = new XMLHttpRequest();
    xhr.open("PUT", url, true);
    xhr.setRequestHeader("Content-Type", contentType);
    xhr.upload.onprogress = (e) => {
      if (e.lengthComputable && onProgress) onProgress(e.loaded / e.total);
    };
    xhr.onload = () => {
      if (xhr.status >= 200 && xhr.status < 300) resolve();
      else reject(new Error("Upload failed (" + xhr.status + "): " + xhr.responseText.slice(0, 200)));
    };
    xhr.onerror = () => reject(new Error("Upload failed: network error"));
    xhr.ontimeout = () => reject(new Error("Upload timed out"));
    xhr.send(body);
  });
}
