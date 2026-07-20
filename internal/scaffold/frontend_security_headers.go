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
  // ws:/wss: keep the dev overlay + HMR socket working. api.ipify.org is the
  // public-IP hint the API client fetches so local audit records show a real
  // address instead of ::1 — dev only, and it must be allowed here or the
  // browser logs a CSP violation on every page load.
  "connect-src 'self' " + API_ORIGIN + (isDev ? " ws: wss: https://api.ipify.org" : ""),
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

// viteSecurityHeaders returns the shared security-header block for the Vite
// (TanStack) web + admin configs. Single source for both, mirroring the Next.js
// helper above and the Go API's middleware.SecurityHeaders.
//
// Vite only applies these when IT serves the bytes — dev (`grit start`) and
// `vite preview`. Production serves dist/ from nginx, which sets the same set
// in apps/<app>/nginx.conf (see viteNginxConf); keep the two in agreement.
//
// Load-bearing details:
//   - script-src needs 'unsafe-inline' here because React Refresh injects an
//     inline preamble in dev. The production build emits NO inline scripts, so
//     nginx.conf can and does use a stricter script-src 'self'.
//   - style-src/font-src must allow Google Fonts. The admin's index.html loads
//     the active theme's fonts from fonts.googleapis.com; CSP blocks them
//     silently, and the app quietly falls back to system-ui and stops matching
//     the Next.js admin. (The web app doesn't use them today — allowing the
//     origin anyway keeps one policy for both instead of two that drift.)
//   - connect-src must include the API origin or every fetch is blocked.
func viteSecurityHeaders() string {
	return `
// Security headers — mirrors the Next.js apps and the Go API's
// middleware.SecurityHeaders. Applied by the dev + preview servers below;
// production serves dist/ via nginx, which sets the same set in nginx.conf.
//
//   * script-src allows 'unsafe-inline' for React Refresh's dev preamble. The
//     production build has no inline scripts, so nginx uses script-src 'self'.
//   * style-src/font-src allow Google Fonts — index.html loads the theme's
//     fonts from fonts.googleapis.com, and CSP blocks them silently.
//   * connect-src must include the API origin or every fetch is blocked.
//
// VITE_API_URL is read through loadEnv, not process.env: Vite does not load
// .env files into process.env for the config file itself, so reading it
// directly always saw undefined and pinned the CSP to localhost:8080 — which
// then blocked every request for anyone who moved the API.
const viteMode = process.env.NODE_ENV || 'development'
const viteEnv = loadEnv(viteMode, process.cwd(), '')
const API_ORIGIN = viteEnv.VITE_API_URL || 'http://localhost:8080'
const isDev = viteMode !== 'production'

const csp = [
  "default-src 'self'",
  "script-src 'self' 'unsafe-inline'",
  "style-src 'self' 'unsafe-inline' https://fonts.googleapis.com",
  "img-src 'self' data: blob: https:",
  "font-src 'self' data: https://fonts.gstatic.com",
  // api.ipify.org is the dev-only public-IP hint the API client fetches so
  // local audit records show a real address instead of ::1.
  "connect-src 'self' ws: wss: " + API_ORIGIN + (isDev ? ' https://api.ipify.org' : ''),
  "frame-ancestors 'none'",
  "base-uri 'self'",
  "form-action 'self'",
  "object-src 'none'",
].join('; ')

const securityHeaders = {
  'Content-Security-Policy': csp,
  'X-Content-Type-Options': 'nosniff',
  'X-Frame-Options': 'DENY',
  'Referrer-Policy': 'strict-origin-when-cross-origin',
  'Permissions-Policy': 'camera=(), microphone=(), geolocation=(), payment=(), usb=()',
  'Strict-Transport-Security': 'max-age=63072000; includeSubDomains; preload',
  'Cross-Origin-Opener-Policy': 'same-origin',
}
`
}

// dockerfileVite builds a Vite (TanStack) app and serves the static bundle with
// nginx. app is "web" or "admin".
//
// Exists because docker_files.go used to hand every frontend the Next.js
// Dockerfile regardless of --vite, so a Vite app's image build failed outright:
// the runner stage copies .next/standalone and runs `node server.js`, but a Vite
// build emits dist/ and there is no server. `grit new --vite` produced an app
// that could not be containerised at all.
//
// The output is static, so there's no Node at runtime — nginx serves dist/ and
// applies the security headers + SPA history fallback (see viteNginxConf).
func dockerfileVite(app string) string {
	return `# Build stage
FROM node:22-alpine AS base

# Pin pnpm — pnpm@latest started resolving to pnpm 11 which needs Node 22's
# node:sqlite builtin. Pinning here avoids surprise breakage on rebuilds.
RUN corepack enable && corepack prepare pnpm@9.15.0 --activate

# Install dependencies
FROM base AS deps
WORKDIR /app

COPY package.json pnpm-lock.yaml pnpm-workspace.yaml ./
COPY apps/` + app + `/package.json ./apps/` + app + `/
COPY packages/shared/package.json ./packages/shared/

RUN pnpm install --frozen-lockfile

# Build
FROM base AS builder
WORKDIR /app

COPY --from=deps /app/node_modules ./node_modules
COPY --from=deps /app/apps/` + app + `/node_modules ./apps/` + app + `/node_modules
COPY --from=deps /app/packages/shared/node_modules ./packages/shared/node_modules
COPY . .

# Vite inlines env at BUILD time (unlike Next.js, which reads it at runtime),
# so these must be set here — setting them on the running container does
# nothing, the values are already baked into the bundle.
ARG VITE_API_URL
ARG VITE_THEME
ENV VITE_API_URL=$VITE_API_URL
ENV VITE_THEME=$VITE_THEME

RUN pnpm --filter ` + app + ` build

# Run — nginx serves the built SPA. No Node in the runtime image.
FROM nginx:1.27-alpine AS runner

COPY --from=builder /app/apps/` + app + `/dist /usr/share/nginx/html
COPY apps/` + app + `/nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 3000

CMD ["nginx", "-g", "daemon off;"]
`
}

// viteNginxConf returns the nginx config that serves a built Vite SPA in
// production. app is "web" or "admin".
//
// nginx gotcha worth knowing: add_header does NOT merge across levels — a
// location block with its own add_header drops every header inherited from the
// server block. That's why the asset cache rule below uses `expires` (which
// sets Cache-Control directly) instead of add_header: it keeps the security
// headers intact for /assets/* rather than silently dropping them there.
func viteNginxConf(app string) string {
	return `# Serves the built Vite SPA (apps/` + app + `/dist) in production.
#
# Security headers mirror the app's vite.config.ts, the Next.js apps, and the
# Go API's middleware.SecurityHeaders. script-src is stricter here than in dev:
# a production Vite build emits no inline scripts, so 'unsafe-inline' isn't
# needed. style-src/font-src still allow Google Fonts because index.html loads
# the active theme's fonts from fonts.googleapis.com.
server {
    listen       3000;
    server_name  _;
    root         /usr/share/nginx/html;
    index        index.html;

    add_header X-Content-Type-Options        "nosniff" always;
    add_header X-Frame-Options               "DENY" always;
    add_header Referrer-Policy               "strict-origin-when-cross-origin" always;
    add_header Permissions-Policy            "camera=(), microphone=(), geolocation=(), payment=(), usb=()" always;
    add_header Strict-Transport-Security     "max-age=63072000; includeSubDomains; preload" always;
    add_header Cross-Origin-Opener-Policy    "same-origin" always;
    add_header Content-Security-Policy       "default-src 'self'; script-src 'self'; style-src 'self' 'unsafe-inline' https://fonts.googleapis.com; img-src 'self' data: blob: https:; font-src 'self' data: https://fonts.gstatic.com; connect-src 'self' https: http:; frame-ancestors 'none'; base-uri 'self'; form-action 'self'; object-src 'none'" always;

    gzip              on;
    gzip_vary         on;
    gzip_min_length   1024;
    gzip_types        text/plain text/css application/javascript application/json image/svg+xml;

    # Hashed filenames are immutable. Uses "expires" rather than add_header so
    # the security headers above are NOT dropped for this location.
    location /assets/ {
        expires 1y;
        access_log off;
    }

    # SPA history fallback: TanStack Router owns the routes client-side, so any
    # deep link (/system/health, /resources/users) must return index.html
    # instead of a 404.
    location / {
        try_files $uri $uri/ /index.html;
    }
}
`
}
