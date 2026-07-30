# Grit UI — ui.gritframework.dev
#
# Built by Dokploy from this directory. Set the app's build context to `ui/`
# so a commit that only touches the framework does not rebuild the site.
#
# Multi-stage: deps are installed once and cached, the build runs against them,
# and the runtime image carries only Next's standalone output — no node_modules,
# no source, no package manager.

# ── deps ─────────────────────────────────────────────────────────────────────
FROM node:22-alpine AS deps
WORKDIR /app
RUN corepack enable && corepack prepare pnpm@11.9.0 --activate
# pnpm-workspace.yaml has to be here, not just in the build context: pnpm 10+
# reads its allowBuilds approvals from that file, and without it sharp's install
# script is treated as unapproved — which is a hard error, not a warning. The
# install otherwise succeeds first, so the failure reads as unrelated to a
# missing file.
COPY package.json pnpm-lock.yaml* pnpm-workspace.yaml* ./
# --no-frozen-lockfile so a first deploy works before a lockfile is committed;
# once ui/pnpm-lock.yaml exists this resolves from it.
RUN pnpm install --no-frozen-lockfile

# ── build ────────────────────────────────────────────────────────────────────
FROM node:22-alpine AS builder
WORKDIR /app
RUN corepack enable && corepack prepare pnpm@11.9.0 --activate
COPY --from=deps /app/node_modules ./node_modules
COPY . .

# Baked into the generated registry URLs, so it must be the public origin.
ARG NEXT_PUBLIC_SITE_URL=https://ui.gritframework.dev
ENV NEXT_PUBLIC_SITE_URL=$NEXT_PUBLIC_SITE_URL
ENV NEXT_TELEMETRY_DISABLED=1

# Zenith's public pair belongs HERE, in the stage that runs the build, and needs
# both an ARG and a matching ENV — Next reads process.env, and a FROM starts a
# fresh stage so an ARG declared earlier does not exist here.
#
# These are compiled into the prerendered HTML. Supplied only at runtime, the
# tracker never appears in the page no matter what the container environment
# says. Both are public by design: they ship in the snippet on every page and
# the site key can only write events.
#
# ZENITH_API_KEY is OPTIONAL here and secret. It reads analytics, which the
# install-count badges need. Supplied as a build arg, the counts are baked into
# the prerendered pages immediately; omitted, the pages simply ship without
# badges and pick the numbers up on the first ISR revalidation (15 minutes) from
# the RUNTIME environment. The runtime variable is the one that matters — set it
# in Dokploy's Environment, not only as a build argument. This stage is discarded
# so nothing here reaches the final image, but a build log is a wider audience
# than a container environment, which is why runtime-only is the safer default.
ARG ZENITH_URL
ARG ZENITH_SITE_KEY
ARG ZENITH_API_KEY
ENV ZENITH_URL=$ZENITH_URL
ENV ZENITH_SITE_KEY=$ZENITH_SITE_KEY
ENV ZENITH_API_KEY=$ZENITH_API_KEY

RUN pnpm build

# ── runtime ──────────────────────────────────────────────────────────────────
FROM node:22-alpine AS runner
WORKDIR /app
ENV NODE_ENV=production
ENV NEXT_TELEMETRY_DISABLED=1
ENV PORT=3100
ENV HOSTNAME=0.0.0.0

# The PUBLIC pair is needed at runtime as well as at build time, and is repeated
# here on purpose. /zenith and the install-count refresh both read ZENITH_URL at
# request time, and a FROM discards the builder stage — so without these the
# dashboard answers 503 in production even with every secret correctly set.
# Verified: the route 503s with all three secrets present and ZENITH_URL absent.
# Safe to bake into a layer; both values already ship in the page HTML.
ARG ZENITH_URL
ARG ZENITH_SITE_KEY
ENV ZENITH_URL=$ZENITH_URL
ENV ZENITH_SITE_KEY=$ZENITH_SITE_KEY

# The runtime secrets are deliberately NOT declared here. Baking them into an
# ENV would write them into an image layer; Dokploy injects them into the running
# container instead. Read at runtime by /zenith and by the ISR revalidation that
# refreshes the install counts:
#
#   ZENITH_API_KEY    reads analytics — counts and dashboard
#   ZENITH_PW_HASH    gates the dashboard
#   ZENITH_JWT_SECRET signs the dashboard session cookie
#
# All three absent → /zenith answers 503 and the badges stay hidden. Neither
# fails the site.

RUN addgroup -g 1001 -S nodejs && adduser -S nextjs -u 1001

COPY --from=builder /app/public ./public
COPY --from=builder --chown=nextjs:nodejs /app/.next/standalone ./
COPY --from=builder --chown=nextjs:nodejs /app/.next/static ./.next/static

USER nextjs
EXPOSE 3100

# The registry is the reason this service exists, so the health check asks for
# it rather than the homepage — a site that renders but serves a broken
# registry.json is not healthy in any way that matters.
HEALTHCHECK --interval=30s --timeout=5s --start-period=15s --retries=3 \
  CMD node -e "fetch('http://127.0.0.1:3100/r/registry.json').then(r=>process.exit(r.ok?0:1)).catch(()=>process.exit(1))"

CMD ["node", "server.js"]
