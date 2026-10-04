# Build stage
FROM node:22.23-alpine AS base

# The pnpm that package.json names in packageManager, which CI installs too.
# The image ran pnpm 9, which does not read onlyBuiltDependencies, so every
# dependency's install script ran during the build. Raise the two together.
RUN corepack enable && corepack prepare pnpm@10.33.4 --activate

# Install dependencies
FROM base AS deps
WORKDIR /app

# pnpm-lock.yaml is written by the first pnpm install, so a project that has
# not run one yet does not have it. The glob keeps this COPY working either
# way: Docker only fails when nothing matches, and package.json always does.
COPY package.json pnpm-workspace.yaml pnpm-lock.yaml* ./
COPY apps/admin/package.json ./apps/admin/
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
COPY --from=deps /app/apps/admin/node_modules ./apps/admin/node_modules
COPY --from=deps /app/packages/shared/node_modules ./packages/shared/node_modules
COPY . .

ARG NEXT_PUBLIC_API_URL
ENV NEXT_PUBLIC_API_URL=$NEXT_PUBLIC_API_URL

RUN pnpm --filter admin build

# Run. A fresh Node image rather than the base stage, and without the package
# managers Node ships: the server is node and the standalone output, so pnpm,
# corepack, npm and yarn in the image that serves traffic are only something
# more to patch.
FROM node:22.23-alpine AS runner
RUN rm -rf /usr/local/lib/node_modules/npm /usr/local/lib/node_modules/corepack \
    /usr/local/bin/npm /usr/local/bin/npx /usr/local/bin/corepack \
    /usr/local/bin/yarn /usr/local/bin/yarnpkg /opt/yarn-*
WORKDIR /app

ENV NODE_ENV=production

RUN addgroup --system --gid 1001 nodejs
RUN adduser --system --uid 1001 nextjs

COPY --from=builder /app/apps/admin/.next/standalone ./
COPY --from=builder /app/apps/admin/.next/static ./apps/admin/.next/static
COPY --from=builder /app/apps/admin/public ./apps/admin/public

# chown before USER so the runtime user can write to /app (Next.js writes
# build manifests + cache here at startup).
RUN chown -R nextjs:nodejs /app

USER nextjs

EXPOSE 3000

ENV PORT=3000
ENV HOSTNAME="0.0.0.0"

# Unhealthy when the server stops answering. Any response counts: a page that
# fails because the API is down is the API's health check to report.
HEALTHCHECK --interval=30s --timeout=5s --start-period=20s --retries=3 \
  CMD node -e "fetch('http://127.0.0.1:'+(process.env.PORT||3000)+'/').then(()=>process.exit(0),()=>process.exit(1))"

CMD ["node", "apps/admin/server.js"]
