package scaffold

// nextSecurityHeaders returns the shared security-header block injected into
// every scaffolded Next.js config (web, admin, docs).
//
// Why this exists: the Go API has always sent security headers via
// middleware.SecurityHeaders, but the Next.js frontends sent none — so a
// scaffolded app's public face scored F on securityheaders.com with all six
// headers missing (Strict-Transport-Security, Content-Security-Policy,
// X-Frame-Options, X-Content-Type-Options, Referrer-Policy, Permissions-Policy).
// Next.js sends no security headers by default; you have to opt in.
//
// The policy deliberately mirrors middleware.SecurityHeaders in the Go API so
// the two halves of an app agree. Keep them in sync.
//
// Two CSP details that are load-bearing — do not "tighten" them without testing:
//
//   - script-src needs 'unsafe-inline'. Next.js inlines its bootstrap and
//     streams the RSC payload through inline <script> tags, so 'self' alone
//     white-screens every app. The alternative is nonce-based CSP generated in
//     middleware, which forces every route to render dynamically and throws away
//     static generation — the wrong default for a static-first framework.
//   - connect-src must include the API origin. In double/triple mode the browser
//     calls the Go API on a different host; omitting it silently breaks every
//     fetch with a CSP violation rather than an HTTP error, which is a horrible
//     thing to debug.
//
// Emitted as plain string concatenation (no JS template literals) because this
// lives inside a Go raw string literal, which cannot contain backticks.
func nextSecurityHeaders() string {
	return `
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
const API_ORIGIN = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";
const isDev = process.env.NODE_ENV !== "production";

const csp = [
  "default-src 'self'",
  "script-src 'self' 'unsafe-inline'" + (isDev ? " 'unsafe-eval'" : ""),
  "style-src 'self' 'unsafe-inline'",
  "img-src 'self' data: blob: https:",
  "font-src 'self' data:",
  // ws:/wss: keep the dev overlay + HMR socket working.
  "connect-src 'self' " + API_ORIGIN + (isDev ? " ws: wss:" : ""),
  "frame-ancestors 'none'",
  "base-uri 'self'",
  "form-action 'self'",
  "object-src 'none'",
].join("; ");

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
`
}

// nextSecurityHeadersConfig returns the NextConfig fields that apply the block
// above. Inserted into each app's config object.
func nextSecurityHeadersConfig() string {
	return `  // Don't advertise the framework + version to attackers.
  poweredByHeader: false,
  async headers() {
    return [{ source: "/:path*", headers: securityHeaders }];
  },
`
}
