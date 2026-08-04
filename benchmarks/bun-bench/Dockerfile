# Bun's own image, production mode.
#
# Bun.serve is single-threaded per process, so one process on a 4-CPU container
# would leave three cores idle while Go and PHP use all four. start.ts forks one
# process per CPU sharing the socket via reusePort — the arrangement Bun's docs
# recommend, and the one their own published benchmarks use.
FROM oven/bun:1-alpine

ENV NODE_ENV=production
WORKDIR /app

COPY package.json ./
RUN bun install --production --frozen-lockfile 2>/dev/null || bun install --production

COPY server.ts start.ts schema.ts ./

EXPOSE 8080
CMD ["bun", "run", "start.ts"]
