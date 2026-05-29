# syntax=docker/dockerfile:1.7
#
# Grit Motors — demo deployment Dockerfile.
#
# Three fixes vs the original kayula version (all of which broke on
# Dokploy):
#   1. Pin pnpm to v9.15.0 — pnpm@latest now resolves to v11+ which
#      requires Node v22 (uses node:sqlite). Pinned matches the
#      pnpm-lock.yaml v9.0 schema we ship.
#   2. Final binary name renamed stoka → gritdemo to match the module.
#   3. /app chowned to the non-root user before USER so Sentinel/Pulse
#      can create their SQLite stores ("out of memory (14)" is actually
#      SQLITE_CANTOPEN for permission-denied; not obvious from the log).

# ── Stage 1: build the frontend ──────────────────────────────────────
FROM node:20-alpine AS frontend
WORKDIR /app/frontend

# Pinned pnpm version — pnpm@latest pulls v11+ which needs Node 22.
RUN corepack enable && corepack prepare pnpm@9.15.0 --activate

COPY frontend/package.json frontend/pnpm-lock.yaml ./
RUN pnpm install --frozen-lockfile
COPY frontend/ ./

# Vite needs env vars at build time — pass via build args.
ARG VITE_INTERNAL_API_KEY
ARG VITE_APP_NAME="Grit Motors"
ENV VITE_INTERNAL_API_KEY=$VITE_INTERNAL_API_KEY
ENV VITE_APP_NAME=$VITE_APP_NAME

RUN pnpm build

# ── Stage 2: build the Go binary ─────────────────────────────────────
FROM golang:1.24-alpine AS backend
RUN apk add --no-cache git
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
COPY --from=frontend /app/frontend/dist ./frontend/dist
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o gritdemo .

# ── Stage 3: minimal runtime image ───────────────────────────────────
FROM alpine:3.20
RUN apk --no-cache add ca-certificates tzdata && \
    addgroup -S grit && adduser -S -G grit grit
WORKDIR /app
COPY --from=backend /app/gritdemo .
# Sentinel + Pulse write SQLite stores into the working dir under their
# own process user. Chown /app to grit BEFORE switching user so the
# stores can be created at runtime.
RUN chown -R grit:grit /app
USER grit
EXPOSE 8080
CMD ["./gritdemo"]
