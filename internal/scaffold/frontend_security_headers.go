package scaffold

import "regexp"

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
// toPlainJavaScript removes the TypeScript from a config that is not TypeScript.
//
// The docs app's config is next.config.mjs, which fumadocs expects and which is
// plain JavaScript. Every shared block spliced into it is written for the .ts
// configs apps/web and apps/admin have, and one annotation in a .mjs file is a
// syntax error that stops Next loading the config at all: the documentation site
// in a --full project had never built.
//
// A transform rather than a second copy of each block, so the two cannot drift
// apart on a policy they are both meant to enforce. Patterns rather than a list
// of exact strings, because the exact-string version was wrong three times in a
// row: each fix revealed the next construct, and the fourth was found by a
// syntax checker rather than by reading.
func toPlainJavaScript(config string) string {
	// "value: string): string {" and the like, in a function signature.
	config = tsParamType.ReplaceAllString(config, "($1)")
	config = tsReturnType.ReplaceAllString(config, ") {")
	// "const pluginOrigins: [string, string][] = ["
	config = tsConstType.ReplaceAllString(config, "$1 =")
	// "as const", and assertions onto a type or a union of them.
	config = tsAssertion.ReplaceAllString(config, "")
	return config
}

var (
	// (value: string) -> (value)
	tsParamType = regexp.MustCompile(`\((\w+): [\w\[\]|" ]+\)`)
	// ): string { -> ) {
	tsReturnType = regexp.MustCompile(`\): [\w\[\]|" ]+ \{`)
	// const x: T = -> const x =
	tsConstType = regexp.MustCompile(`(const \w+): [^=\n]+ =`)
	// " as const", ` as "http" | "https"`, " as SomeType"
	tsAssertion = regexp.MustCompile(` as (?:const|[\w"]+(?: \| [\w"]+)*)`)
)

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

` + nextAPIOriginLine + nextAPIWSOrigin + `// Browser-facing origin of stored files. Uploads are presigned PUTs made
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
` + nextIsDevLine + nextImageOrigins + `
` + cspPluginOrigins + `const csp = [
  "default-src 'self'",
  "script-src 'self' 'unsafe-inline'" + (isDev ? " 'unsafe-eval'" : ""),
  "style-src 'self' 'unsafe-inline'",
  ` + cspImgSrcTight + `
  ` + cspMediaSrc + `
  "font-src 'self' data:",
  // ws:/wss: keep the dev overlay + HMR socket working. api.ipify.org is the
  // public-IP hint the API client fetches so local audit records show a real
  // address instead of ::1 — dev only, and it must be allowed here or the
  // browser logs a CSP violation on every page load.
  ` + nextConnectSrcUpload + ` ? " ws: wss: https://api.ipify.org" : ""),
  ` + cspFrameSrc + `
  "frame-ancestors 'none'",
  "base-uri 'self'",
  "form-action 'self'",
  "object-src 'none'",
` + cspJoinWithPlugins + `
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
  // Both hosts: picsum.photos 302s to fastly.picsum.photos, and next/image
  // checks the host it is finally fetching from.
  ...(isDev
    ? [
        { protocol: "https" as const, hostname: "picsum.photos", pathname: "/**" },
        { protocol: "https" as const, hostname: "fastly.picsum.photos", pathname: "/**" },
      ]
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
`
}

// nextSecurityHeadersConfig returns the NextConfig fields that apply the block
// above. Inserted into each app's config object.
// staticExportConfig is the same block for an app with no server to run it.
//
// Two things cannot survive an export, and Next says so during the build rather
// than at runtime:
//
//   - headers(). There is no server to set them. The Go binary that serves these
//     files sets the same ones in middleware.SecurityHeaders, which is where they
//     belonged anyway: a header set by the frontend host does not protect the API
//     on the same origin.
//   - The image optimizer, for the same reason. Without unoptimized, next/image
//     emits URLs pointing at /_next/image, which the export does not contain, and
//     every image 404s.
func staticExportConfig() string {
	return `  // Don't advertise the framework + version to attackers.
  poweredByHeader: false,
  // No optimizer in an export: there is no server to run one. next/image would
  // otherwise emit /_next/image URLs that are not in the output, and every
  // image on the site would 404.
  images: {
    remotePatterns: nextImageHosts,
    unoptimized: true,
  },
  // No headers() either. The Go binary serving these files sets the same ones
  // in middleware.SecurityHeaders, which is the only place that can set them
  // for the API on the same origin.
`
}

func nextSecurityHeadersConfig(opts Options) string {
	if opts.StaticExport() {
		return staticExportConfig()
	}
	return `  // Don't advertise the framework + version to attackers.
  poweredByHeader: false,
  // next/image refuses any remote host it was not told about, and it THROWS
  // rather than falling back to a plain <img>, so one uploaded image takes the
  // whole page down with "hostname is not configured". Stored files live on the
  // storage origin, not this one, so that host has to be named here.
  //
  // Derived from the same STORAGE_ORIGIN the CSP uses, so moving storage to a
  // CDN is one env var rather than two places that drift.
` + nextImagesNew + `
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
// envDir is the JavaScript expression for the directory holding .env: the
// project root in a single, two levels up in a monorepo app.
func viteSecurityHeaders(envDir string) string {
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
//
// envDir, not cwd: in a monorepo the app runs from apps/web or apps/admin and
// there is no .env there, so every value read here came back undefined and the
// defaults below were all anyone ever got.
const viteMode = process.env.NODE_ENV || 'development'
const viteEnv = loadEnv(viteMode, ` + envDir + `, '')
// Where the Go API listens. APP_PORT is the same variable the binary reads, so
// moving it moves the dev proxy with it.
const apiTarget = viteEnv.VITE_API_URL || 'http://localhost:' + (viteEnv.APP_PORT || '8080')
// A CSP source expression matches paths EXACTLY unless it ends in '/', so a
// value carrying a path ('http://host/api/v1') allows that one path and blocks
// every route under it, silently, as a console violation rather than an HTTP
// status. Reduce whatever is configured to its origin, and pass a value we
// cannot parse through unchanged rather than dropping it.
function toOrigin(value: string): string {
  try {
    return new URL(value).origin
  } catch {
    return value
  }
}

const API_ORIGIN = toOrigin(apiTarget)
// The browser bundle reads import.meta.env.VITE_API_URL, and Vite only exposes
// a variable to the client when it is set with the VITE_ prefix. A project
// configures the API with APP_PORT in the root .env, which is not one, so the
// client fell back to localhost:8080 while the policy below authorised the
// real port: every request blocked, as a console violation with no HTTP
// status. Handing the client the value already resolved above means the proxy,
// the policy and the bundle cannot disagree.
const clientEnv = {
  'import.meta.env.VITE_API_URL': JSON.stringify(apiTarget),
}
// Browser-facing storage origin — presigned uploads PUT here directly and
// stored images load from it. Defaults to local MinIO; set VITE_STORAGE_URL
// to your S3/R2/B2 public origin in production.
const STORAGE_ORIGIN = toOrigin(viteEnv.VITE_STORAGE_URL || 'http://localhost:9002')
// Where presigned uploads are actually sent, when that is not where the files
// are read back from. A CSP source is an origin, so listing only the read origin
// blocks the PUT: the API signs a perfectly good URL and the browser refuses it,
// which shows up as a console violation and never as an HTTP status. A managed
// bucket does this by default, serving reads from its own CDN host and signing
// writes for the underlying S3 endpoint. Unset, it costs nothing.
const STORAGE_UPLOAD_ORIGIN = viteEnv.VITE_STORAGE_UPLOAD_URL
  ? ' ' + toOrigin(viteEnv.VITE_STORAGE_UPLOAD_URL)
  : ''
` + viteIsDevLine + viteImageOrigins + `
const csp = [
  "default-src 'self'",
  "script-src 'self' 'unsafe-inline'",
  "style-src 'self' 'unsafe-inline' https://fonts.googleapis.com",
  ` + cspImgSrcTight + `
  ` + cspMediaSrc + `
  "font-src 'self' data: https://fonts.gstatic.com",
  // api.ipify.org is the dev-only public-IP hint the API client fetches so
  // local audit records show a real address instead of ::1.
  "connect-src 'self' ws: wss: " + API_ORIGIN + " " + STORAGE_ORIGIN + STORAGE_UPLOAD_ORIGIN + (isDev ? ' https://api.ipify.org' : ''),
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
FROM ` + nodeImage + ` AS base

` + pnpmPinNew + `
# Install dependencies
FROM base AS deps
WORKDIR /app

# pnpm-lock.yaml is written by the first pnpm install, so a project that has
# not run one yet does not have it. The glob keeps this COPY working either
# way: Docker only fails when nothing matches, and package.json always does.
COPY package.json pnpm-workspace.yaml pnpm-lock.yaml* ./
COPY apps/` + app + `/package.json ./apps/` + app + `/
# The whole workspace, not a hand-listed subset. Both apps depend on
# @repo/upload as well as @repo/shared, and naming members here meant pnpm
# stopped with ERR_PNPM_WORKSPACE_PKG_NOT_FOUND the first time one was
# added. Costs the install cache when packages/ changes; the builder stage
# copies everything a few lines later anyway.
COPY packages ./packages

# Frozen when there is a lockfile to freeze to, which is what a repo that
# commits its lock wants, and a plain install when there is not, so a
# freshly scaffolded project still builds an image.
RUN if [ -f pnpm-lock.yaml ]; then pnpm install --frozen-lockfile; else pnpm install; fi

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
FROM ` + nginxImage + ` AS runner

COPY --from=builder /app/apps/` + app + `/dist /usr/share/nginx/html
COPY apps/` + app + `/nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 3000

` + viteHealthcheck + `CMD ["nginx", "-g", "daemon off;"]
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
    add_header Content-Security-Policy       "default-src 'self'; script-src 'self'; style-src 'self' 'unsafe-inline' https://fonts.googleapis.com; img-src 'self' data: blob: https:; media-src 'self' blob: https:; font-src 'self' data: https://fonts.gstatic.com; ` + nginxConnectSrcNew + ` frame-ancestors 'none'; base-uri 'self'; form-action 'self'; object-src 'none'" always;

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
