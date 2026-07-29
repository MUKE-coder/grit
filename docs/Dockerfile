FROM node:20-alpine AS base

# Install pnpm
RUN corepack enable && corepack prepare pnpm@9 --activate

# --- Dependencies ---
FROM base AS deps
WORKDIR /app
COPY package.json pnpm-lock.yaml ./
RUN pnpm install --frozen-lockfile

# --- Build ---
FROM base AS builder
WORKDIR /app
COPY --from=deps /app/node_modules ./node_modules
COPY . .

# Zenith's public pair has to be present for `pnpm run build`, not just at
# runtime. The root layout renders <Analytics config={ZENITH_PUBLIC} /> and
# nearly every page here is prerendered, so whatever these resolve to during
# the build is baked into the HTML permanently.
#
# config/zenith.ts falls back to a placeholder site key rather than throwing,
# which makes getting this wrong invisible: the image builds, the site serves,
# and every page ships zk_the_public_site_key while the dashboard stays empty.
# Setting them only in the runner stage below does nothing for those pages.
ARG ZENITH_URL
ARG ZENITH_SITE_KEY
ENV ZENITH_URL=$ZENITH_URL
ENV ZENITH_SITE_KEY=$ZENITH_SITE_KEY

RUN pnpm run build

# --- Production ---
FROM base AS runner
WORKDIR /app

ENV NODE_ENV=production
ENV NEXT_TELEMETRY_DISABLED=1

RUN addgroup --system --gid 1001 nodejs
RUN adduser --system --uid 1001 nextjs

# Copy public assets
COPY --from=builder /app/public ./public

# Copy standalone output
COPY --from=builder --chown=nextjs:nodejs /app/.next/standalone ./
COPY --from=builder --chown=nextjs:nodejs /app/.next/static ./.next/static

USER nextjs

EXPOSE 3000
ENV PORT=3000
ENV HOSTNAME="0.0.0.0"

# Repeated at runtime because the /zenith dashboard route reads the same pair
# through ZENITH_CONFIG when it serves a request. The build-time copy above is
# what puts the tracker on the pages; this is what lets the dashboard query it.
ARG ZENITH_URL
ARG ZENITH_SITE_KEY
ENV ZENITH_URL=$ZENITH_URL
ENV ZENITH_SITE_KEY=$ZENITH_SITE_KEY

# ZENITH_API_KEY, ZENITH_PW_HASH and ZENITH_JWT_SECRET are deliberately NOT
# declared here. They are secrets, and an ARG would persist them in the image
# history for anyone who pulls it. Set them as environment variables on the
# deployment instead. Until all three are present, zenithDashboardReady() is
# false and /zenith answers 503 rather than starting up half-configured.

CMD ["node", "server.js"]
