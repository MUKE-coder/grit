# Next.js in production: `next build` then the standalone server. Not `next
# dev`, which recompiles on demand and would make this meaningless.
#
# Node is single-threaded, so server.js is wrapped by cluster.js to fork one
# worker per CPU — the same allowance every other framework in this benchmark
# gets on a 4-CPU container.
FROM node:22-alpine AS builder

WORKDIR /app
COPY package.json ./
COPY prisma ./prisma
RUN npm install --no-audit --no-fund && npx prisma generate

COPY . .
ENV NEXT_TELEMETRY_DISABLED=1
RUN npm run build

FROM node:22-alpine AS runner

ENV NODE_ENV=production
ENV NEXT_TELEMETRY_DISABLED=1
WORKDIR /app

# The standalone output already contains the trimmed node_modules it needs.
COPY --from=builder /app/.next/standalone ./
COPY --from=builder /app/.next/static ./.next/static
# Next's standalone trace does not follow Prisma's generated client or its
# native query engine binary, so copy them in explicitly. Without this the
# server boots fine and then 500s on the first query.
COPY --from=builder /app/node_modules/.prisma ./node_modules/.prisma
COPY --from=builder /app/node_modules/@prisma ./node_modules/@prisma
COPY cluster.js ./

EXPOSE 8080
ENV PORT=8080
ENV HOSTNAME=0.0.0.0

CMD ["node", "cluster.js"]
