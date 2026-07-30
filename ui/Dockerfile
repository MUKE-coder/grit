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
# ZENITH_API_KEY is deliberately absent: it is a secret, it is only read by a
# dashboard route this site does not mount, and tracking does not need it.
ARG ZENITH_URL
ARG ZENITH_SITE_KEY
ENV ZENITH_URL=$ZENITH_URL
ENV ZENITH_SITE_KEY=$ZENITH_SITE_KEY

RUN pnpm build

# ── runtime ──────────────────────────────────────────────────────────────────
FROM node:22-alpine AS runner
WORKDIR /app
ENV NODE_ENV=production
ENV NEXT_TELEMETRY_DISABLED=1
ENV PORT=3100
ENV HOSTNAME=0.0.0.0

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
