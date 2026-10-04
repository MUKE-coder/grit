import type { NextConfig } from "next";
import { existsSync, readFileSync } from "node:fs";
import { resolve } from "node:path";

// Hoist the monorepo's root .env into process.env. Next.js auto-loads
// .env only from the package's own directory, so without this the THEME
// and SOCIAL_AUTH_ENABLED values set at the root are invisible to the
// admin app. Shell env wins — we only fill in unset keys.
const rootEnv = resolve(process.cwd(), "..", "..", ".env");
if (existsSync(rootEnv)) {
  for (const line of readFileSync(rootEnv, "utf8").split(/\r?\n/)) {
    const m = line.match(/^([A-Z_][A-Z0-9_]*)=(.*)$/i);
    if (!m) continue;
    if (process.env[m[1]] === undefined) process.env[m[1]] = m[2].trim();
  }
}


// --- Security headers -------------------------------------------------------
// Next.js ships no security headers by default. These mirror the Go API's
// middleware.SecurityHeaders so both halves of the app agree; keep them in sync.
//
// CSP caveats (both load-bearing — test before tightening):
//   * script-src needs 'unsafe-inline' — Next inlines its bootstrap + streams
//     the RSC payload via inline <script>. 'self' alone white-screens the app.
//     Locking this down means nonce-based CSP in middleware, which forces every
//     route to render dynamically.
//   * connect-src must include the API origin — in double/triple mode the
//     browser calls the Go API cross-origin, and a missing entry breaks every
//     fetch with a CSP violation.
// A CSP source expression matches paths EXACTLY unless it ends in "/", so
// "http://api.example.com/api/v1" allows that one path and blocks every route
// under it. Setting NEXT_PUBLIC_API_URL with a path is a natural mistake and
// the failure is silent (a console violation, never an HTTP status), so reduce
// whatever is configured to its origin. A value we cannot parse is passed
// through unchanged rather than dropped, which would break every fetch.
function toOrigin(value: string): string {
  try {
    return new URL(value).origin;
  } catch {
    return value;
  }
}

const API_ORIGIN = toOrigin(process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080");
// The realtime socket is the API origin over ws: or wss:. A CSP source matches
// its scheme exactly, so API_ORIGIN does not admit it, and without this the
// browser blocks every live update once the dev-only ws: wss: below is gone.
const API_WS_ORIGIN = API_ORIGIN.replace(/^http/, "ws");
// Browser-facing origin of stored files. Uploads are presigned PUTs made
// directly from the browser to object storage, and stored images are served
// from the same host — both are blocked unless this origin is in connect-src
// and img-src. Defaults to the local MinIO endpoint; in production set
// NEXT_PUBLIC_STORAGE_URL to your S3/R2/B2 public origin
// (e.g. https://cdn.example.com or https://<bucket>.s3.<region>.amazonaws.com).
const STORAGE_ORIGIN = toOrigin(process.env.NEXT_PUBLIC_STORAGE_URL || "http://localhost:9002");
// Where presigned uploads are actually sent, when that is not where the files
// are read back from. A CSP source is an origin, so listing only the read origin
// blocks the PUT: the API signs a perfectly good URL and the browser refuses it,
// which shows up as a console violation and never as an HTTP status. A managed
// bucket does this by default, serving reads from its own CDN host and signing
// writes for the underlying S3 endpoint. Unset, it costs nothing.
const STORAGE_UPLOAD_ORIGIN = process.env.NEXT_PUBLIC_STORAGE_UPLOAD_URL
  ? " " + toOrigin(process.env.NEXT_PUBLIC_STORAGE_UPLOAD_URL)
  : "";
const isDev = process.env.NODE_ENV !== "production";
// Hosts images load from besides this site, the API and storage: the avatars
// Google and GitHub sign-in save, plus NEXT_PUBLIC_IMAGE_ORIGINS (comma-separated
// origins, for a CDN or an identity provider's avatars). img-src used to allow
// every https: host, so markup that got into a page could load, and report to,
// any server.
const IMAGE_ORIGINS = ["https://lh3.googleusercontent.com", "https://avatars.githubusercontent.com"]
  .concat((process.env.NEXT_PUBLIC_IMAGE_ORIGINS || "").split(",").map((s) => s.trim()).filter(Boolean).map((s) => toOrigin(s)))
  .concat(isDev ? ["https://picsum.photos"] : [])
  .join(" ");

// Third-party origins a plugin needs, as [directive, origin]. Each is
// appended to its directive below. grit plugin add stripe puts Stripe.js here,
// and grit plugin remove takes it out again.
const pluginOrigins: [string, string][] = [
  // grit:csp-origins
];

const csp = [
  "default-src 'self'",
  "script-src 'self' 'unsafe-inline'" + (isDev ? " 'unsafe-eval'" : ""),
  "style-src 'self' 'unsafe-inline'",
  "img-src 'self' data: blob: " + API_ORIGIN + " " + STORAGE_ORIGIN + " " + IMAGE_ORIGINS,
  "media-src 'self' blob: " + API_ORIGIN + " " + STORAGE_ORIGIN,
  "font-src 'self' data:",
  // ws:/wss: keep the dev overlay + HMR socket working. api.ipify.org is the
  // public-IP hint the API client fetches so local audit records show a real
  // address instead of ::1 — dev only, and it must be allowed here or the
  // browser logs a CSP violation on every page load.
  "connect-src 'self' " + API_ORIGIN + " " + API_WS_ORIGIN + " " + STORAGE_ORIGIN + STORAGE_UPLOAD_ORIGIN + (isDev ? " ws: wss: https://api.ipify.org" : ""),
  "frame-src 'self'",
  "frame-ancestors 'none'",
  "base-uri 'self'",
  "form-action 'self'",
  "object-src 'none'",
]
  .map((directive) => {
    const name = directive.split(" ")[0];
    const extra = pluginOrigins.filter(([d]) => d === name).map(([, origin]) => origin);
    return extra.length > 0 ? directive + " " + extra.join(" ") : directive;
  })
  .join("; ");

// Hosts next/image is allowed to fetch from. Same source as the CSP above:
// stored uploads come from the storage origin, and picsum.photos is where
// "grit generate resource --faker" points its placeholder images, which is why
// it is here in development and not in production.
const storageURL = new URL(STORAGE_ORIGIN);
const nextImageHosts = [
  {
    protocol: storageURL.protocol.replace(":", "") as "http" | "https",
    hostname: storageURL.hostname,
    port: storageURL.port,
    pathname: "/**",
  },
  ...(isDev
    ? [{ protocol: "https" as const, hostname: "picsum.photos", pathname: "/**" }]
    : []),
];

const securityHeaders = [
  { key: "Content-Security-Policy", value: csp },
  { key: "X-Content-Type-Options", value: "nosniff" },
  { key: "X-Frame-Options", value: "DENY" },
  { key: "Referrer-Policy", value: "strict-origin-when-cross-origin" },
  {
    key: "Permissions-Policy",
    value: "camera=(), microphone=(), geolocation=(), payment=(), usb=()",
  },
  // Browsers ignore HSTS over plain http, so sending it in dev is harmless.
  {
    key: "Strict-Transport-Security",
    value: "max-age=63072000; includeSubDomains; preload",
  },
  { key: "Cross-Origin-Opener-Policy", value: "same-origin" },
];

const nextConfig: NextConfig = {
  output: "standalone",
  reactStrictMode: true,
  // @repo/upload ships raw TypeScript (its package.json points main at src/),
  // so Next has to transpile it. lib/api-client.ts has imported it since
  // v3.31.30 and this list did not mention it: the build fails with "Module
  // parse failed: Unexpected token" on a type-only import, which is what a
  // workspace package that is never compiled looks like.
  transpilePackages: ["@repo/shared", "@repo/upload"],
  // Mirror THEME + SOCIAL_AUTH_ENABLED from .env into the NEXT_PUBLIC_*
  // namespace so server components and the client bundle both see the
  // active theme without a flash of unstyled content. Falls back to the
  // safe defaults when the vars aren't set (atlas, social auth on).
  env: {
    NEXT_PUBLIC_THEME: process.env.THEME || "atlas",
    NEXT_PUBLIC_SOCIAL_AUTH_ENABLED: process.env.SOCIAL_AUTH_ENABLED || "true",
  },
  // Don't advertise the framework + version to attackers.
  poweredByHeader: false,
  // next/image refuses any remote host it was not told about, and it THROWS
  // rather than falling back to a plain <img>, so one uploaded image takes the
  // whole page down with "hostname is not configured". Stored files live on the
  // storage origin, not this one, so that host has to be named here.
  //
  // Derived from the same STORAGE_ORIGIN the CSP uses, so moving storage to a
  // CDN is one env var rather than two places that drift.
  images: {
    remotePatterns: nextImageHosts,
    // Next refuses an upstream image whose host resolves to a private IP,
    // which in development is all of them: the API and MinIO both live on
    // localhost. Without this the optimizer answers 400 and the page shows
    // alt text, with the reason only in the terminal. Never in production,
    // where storage is a real origin and this guard is worth having.
    dangerouslyAllowLocalIP: isDev,
  },
  async headers() {
    return [{ source: "/:path*", headers: securityHeaders }];
  },
  // Uncomment and run "ANALYZE=true pnpm build" to inspect the bundle
  // ...(process.env.ANALYZE === "true"
  //   ? { ...require("@next/bundle-analyzer")({ enabled: true })(nextConfig) }
  //   : {}),
};

export default nextConfig;
