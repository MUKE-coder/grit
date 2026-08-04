# Node in production mode, on Prisma.
#
# One process would leave three of four CPUs idle while Go schedules goroutines
# across all of them, so cluster.js forks one worker per CPU — the same
# allowance every framework in this benchmark gets.
FROM node:22-alpine

ENV NODE_ENV=production
WORKDIR /app

# Prisma's generate step needs the schema present, and it runs on postinstall.
COPY package.json ./
COPY prisma ./prisma
# --include=dev because the prisma CLI is a devDependency and generate needs it.
RUN npm install --include=dev --no-audit --no-fund \
    && npx prisma generate \
    && npm prune --omit=dev

COPY server.js cluster.js ./

EXPOSE 8080
CMD ["node", "cluster.js"]
