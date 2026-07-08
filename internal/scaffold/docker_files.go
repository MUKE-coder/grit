package scaffold

import (
	"fmt"
	"path/filepath"
)

func writeDockerFiles(root string, opts Options) error {
	files := map[string]string{
		filepath.Join(root, "docker-compose.yml"): dockerCompose(opts),
	}

	if opts.Architecture == ArchSingle {
		// Single-app: one binary (Go server + embedded SPA), one Dockerfile at
		// the project root. We deliberately skip apps/api/Dockerfile + the
		// multi-app docker-compose.prod.yml because both expect a monorepo
		// layout (apps/api/go.mod, apps/web, ...) that doesn't exist in
		// --single. PaaS tools like Dokploy auto-detect Dockerfiles, so any
		// leftover apps/api/Dockerfile gets picked up and fails the build.
		files[filepath.Join(root, "Dockerfile")] = dockerfileSingle()
	} else {
		files[filepath.Join(root, "docker-compose.prod.yml")] = dockerComposeProd(opts)
		files[filepath.Join(root, "apps", "api", "Dockerfile")] = dockerfileAPI()
		if opts.ShouldIncludeWeb() {
			files[filepath.Join(root, "apps", "web", "Dockerfile")] = dockerfileNextJS("web")
		}
		if opts.ShouldIncludeAdmin() {
			files[filepath.Join(root, "apps", "admin", "Dockerfile")] = dockerfileNextJS("admin")
		}
		if opts.ShouldIncludeDocs() {
			files[filepath.Join(root, "apps", "docs", "Dockerfile")] = dockerfileNextJS("docs")
		}
	}

	if opts.Architecture != ArchAPI {
		files[filepath.Join(root, ".dockerignore")] = dockerIgnore()
	}

	for path, content := range files {
		if err := writeFile(path, content); err != nil {
			return fmt.Errorf("writing %s: %w", path, err)
		}
	}

	return nil
}

func dockerCompose(opts Options) string {
	return fmt.Sprintf(`# Local development infrastructure.
#
# Every port is bound to 127.0.0.1 — that means these services are
# reachable from YOUR machine only, not from your LAN, not from the
# wider internet. Public binding (the Docker default) on a laptop on
# a coffee-shop wifi exposes Postgres with grit:grit credentials to
# the room. Don't do that.
#
# If you need to share infra with a teammate, use a VPN (Tailscale,
# Wireguard) and forward 127.0.0.1:5432 from your machine — never
# rebind to 0.0.0.0.
services:
  postgres:
    image: postgres:16-alpine
    container_name: %s-postgres
    restart: unless-stopped
    ports:
      # Host port 5434 (not the Postgres-default 5432) avoids two common
      # collisions: a Postgres service installed on the host, and Windows
      # WinNAT reserving 5432 inside its dynamic Hyper-V port range
      # ("An attempt was made to access a socket in a way forbidden..."
      # — see `+"`netsh int ipv4 show excludedportrange protocol=tcp`"+`).
      # The container still listens on 5432 inside the Docker network.
      - "127.0.0.1:5434:5432"
    # Credentials come from .env (POSTGRES_USER / POSTGRES_PASSWORD / POSTGRES_DB).
    # Edit them ONLY in .env — never duplicate them here. The :- syntax
    # provides a fallback so 'docker compose up' still works if the env
    # var is somehow missing.
    environment:
      POSTGRES_USER: ${POSTGRES_USER:-grit}
      POSTGRES_PASSWORD: ${POSTGRES_PASSWORD:-grit}
      POSTGRES_DB: ${POSTGRES_DB:-%s}
    volumes:
      - postgres-data:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U ${POSTGRES_USER:-grit}"]
      interval: 5s
      timeout: 5s
      retries: 5

  redis:
    image: redis:7-alpine
    container_name: %s-redis
    restart: unless-stopped
    ports:
      # Host 6380 (not the Redis-default 6379) avoids native Redis installs
      # — Memurai on Windows, brew/apt Redis on macOS/Linux, WSL Redis.
      # Container still listens on 6379 inside the Docker network.
      - "127.0.0.1:6380:6379"
    volumes:
      - redis-data:/data
    healthcheck:
      test: ["CMD", "redis-cli", "ping"]
      interval: 5s
      timeout: 5s
      retries: 5

  minio:
    image: minio/minio
    container_name: %s-minio
    restart: unless-stopped
    ports:
      # Host 9002 / 9003 (not the MinIO-default 9000 / 9001) avoids
      # collisions with Portainer (9000) and a handful of monitoring
      # stacks that grab 9000/9001. Container still listens on 9000/9001
      # inside the Docker network.
      #
      # Bound to all interfaces (not 127.0.0.1) so a phone/emulator on your LAN
      # can load uploaded images: stored URLs point at this host:9002 and the
      # Expo app rewrites "localhost" to your dev IP (apps/expo/lib/images.ts).
      - "9002:9000"
      - "9003:9001"
    environment:
      MINIO_ROOT_USER: minioadmin
      MINIO_ROOT_PASSWORD: minioadmin
    volumes:
      - minio-data:/data
    command: server /data --console-address ":9001"

  mailhog:
    image: mailhog/mailhog
    container_name: %s-mailhog
    restart: unless-stopped
    ports:
      - "127.0.0.1:1025:1025"
      - "127.0.0.1:8025:8025"

volumes:
  postgres-data:
  redis-data:
  minio-data:
`, opts.ProjectName, opts.ProjectName, opts.ProjectName, opts.ProjectName, opts.ProjectName)
}

func dockerComposeProd(opts Options) string {
	name := opts.ProjectName

	result := fmt.Sprintf(`# Production stack.
#
# Security posture: nothing in this file uses `+"`ports:`"+` — only
# `+"`expose:`"+`. That means none of these services bind to the public host
# interface. Traffic reaches the API and the front-ends ONLY through
# your reverse proxy (Traefik, Caddy, nginx, or a managed PaaS like
# Dokploy/Coolify), which lives on the same Docker network and routes
# domain.com → `+"`api:8080`"+`, app.domain.com → `+"`admin:3000`"+`, etc.
#
# Postgres and Redis have NO host binding at all — they're reachable
# only by containers on the `+"`%s`"+` network. That's the property you
# want: a successful host-level compromise still has to pop a container
# to reach the database.
#
# MinIO is included as an option but most production deployments use
# Cloudflare R2 or AWS S3 (see STORAGE_* in .env). Delete the minio
# service if you're on managed object storage.
#
# To configure your reverse proxy:
#   - Point your apex domain at the api container's expose port (8080)
#   - Point app.domain.com at the admin container (3000)
#   - Point www.domain.com at the web container (3000)
#   - Use Let's Encrypt or your PaaS' auto-TLS for HTTPS — never serve
#     the API over plain HTTP in production.

services:
  api:
    build:
      context: ./apps/api
      dockerfile: Dockerfile
    container_name: %s-api
    restart: unless-stopped
    expose:
      - "8080"
    env_file:
      - .env
    # POSTGRES_USER / POSTGRES_PASSWORD / POSTGRES_DB come from .env
    # (env_file above). We override POSTGRES_HOST so the API hits the
    # postgres container on the Docker network, and POSTGRES_PORT back
    # to 5432 because the postgres container listens on 5432 internally
    # — only the dev host port mapping uses 5434. The Go config builds
    # DATABASE_URL from these parts at startup.
    environment:
      APP_ENV: production
      POSTGRES_HOST: postgres
      POSTGRES_PORT: "5432"
      REDIS_URL: redis://redis:6379
      MINIO_ENDPOINT: http://minio:9000
    depends_on:
      postgres:
        condition: service_healthy
      redis:
        condition: service_healthy
    networks:
      - %s
`, name, name, name)

	if opts.ShouldIncludeWeb() {
		result += fmt.Sprintf(`
  web:
    build:
      context: .
      dockerfile: apps/web/Dockerfile
      args:
        NEXT_PUBLIC_API_URL: ${API_URL:-http://localhost:8080}
    container_name: %s-web
    restart: unless-stopped
    expose:
      - "3000"
    networks:
      - %s
`, name, name)
	}

	if opts.ShouldIncludeAdmin() {
		result += fmt.Sprintf(`
  admin:
    build:
      context: .
      dockerfile: apps/admin/Dockerfile
      args:
        NEXT_PUBLIC_API_URL: ${API_URL:-http://localhost:8080}
    container_name: %s-admin
    restart: unless-stopped
    expose:
      - "3000"
    networks:
      - %s
`, name, name)
	}

	if opts.ShouldIncludeDocs() {
		result += fmt.Sprintf(`
  docs:
    build:
      context: .
      dockerfile: apps/docs/Dockerfile
    container_name: %s-docs
    restart: unless-stopped
    expose:
      - "3002"
    networks:
      - %s
`, name, name)
	}

	result += fmt.Sprintf(`
  postgres:
    image: postgres:16-alpine
    container_name: %s-postgres
    restart: unless-stopped
    env_file:
      - .env
    # Identical credentials to the api service above — they come from the
    # SAME .env block, so they can't drift. The :-grit fallbacks below
    # protect against a missing env var blocking startup.
    environment:
      POSTGRES_USER: ${POSTGRES_USER:-grit}
      POSTGRES_PASSWORD: ${POSTGRES_PASSWORD:-grit}
      POSTGRES_DB: ${POSTGRES_DB:-%s}
    volumes:
      - postgres-data:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U ${POSTGRES_USER:-grit} -d ${POSTGRES_DB:-%s}"]
      interval: 5s
      timeout: 5s
      retries: 5
    networks:
      - %s

  redis:
    image: redis:7-alpine
    container_name: %s-redis
    restart: unless-stopped
    volumes:
      - redis-data:/data
    healthcheck:
      test: ["CMD", "redis-cli", "ping"]
      interval: 5s
      timeout: 5s
      retries: 5
    networks:
      - %s

  minio:
    image: minio/minio
    container_name: %s-minio
    restart: unless-stopped
    env_file:
      - .env
    environment:
      MINIO_ROOT_USER: ${MINIO_ACCESS_KEY:-minioadmin}
      MINIO_ROOT_PASSWORD: ${MINIO_SECRET_KEY:-minioadmin}
    volumes:
      - minio-data:/data
    networks:
      - %s
    command: server /data --console-address ":9001"

networks:
  %s:
    driver: bridge

volumes:
  postgres-data:
  redis-data:
  minio-data:
`, name, name, name, name, name, name, name, name, name)

	return result
}

// dockerfileSingle is the Dockerfile for --single architecture: one binary
// that bundles the Go server and an embedded SPA. Multi-stage:
//
//  1. node:22-alpine   builds the frontend bundle to /app/frontend/dist
//  2. golang:1.24-alpine builds the Go binary with //go:embed picking up
//     the frontend output from the project root
//  3. alpine:3.19      runtime, drops to a non-root user
//
// Notes that paid for themselves debugging real deploys:
//   - pnpm is pinned to a version compatible with Node 22 (pnpm@latest
//     started pulling pnpm 11 which assumes node:sqlite from Node 22.x and
//     broke node:20 base images).
//   - chown /app to the runtime user BEFORE the USER directive. Otherwise
//     Sentinel/Pulse fail to open their embedded SQLite stores and report
//     the misleading "out of memory (14)" — which is actually
//     SQLITE_CANTOPEN for permission-denied.
func dockerfileSingle() string {
	return `# syntax=docker/dockerfile:1.7
# ---------- Stage 1: build the SPA bundle ----------
FROM node:22-alpine AS frontend
WORKDIR /app/frontend

# Pin pnpm — pnpm@latest started resolving to pnpm 11 which needs Node 22's
# node:sqlite builtin; pinning avoids surprise breakage from that drift.
RUN corepack enable && corepack prepare pnpm@9.15.0 --activate

COPY frontend/package.json frontend/pnpm-lock.yaml* ./
RUN pnpm install --no-frozen-lockfile

COPY frontend/ ./
RUN pnpm build

# ---------- Stage 2: build the Go binary ----------
FROM golang:1.24-alpine AS gobuild
WORKDIR /app

RUN apk add --no-cache git
COPY go.mod go.sum ./
RUN go mod download

# Copy source, then drop the freshly-built SPA bundle into the location the
# //go:embed all:frontend/dist directive expects (project root).
COPY . .
COPY --from=frontend /app/frontend/dist ./frontend/dist

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /out/server .

# ---------- Stage 3: minimal runtime ----------
FROM alpine:3.19
RUN apk add --no-cache ca-certificates tzdata

# Create non-root user.
RUN addgroup -S app && adduser -S -G app app

WORKDIR /app
COPY --from=gobuild /out/server ./server

# IMPORTANT: chown happens before the USER directive. Sentinel/Pulse open
# embedded SQLite stores relative to CWD; without write access here they
# crash with "unable to open database file: out of memory (14)".
RUN chown -R app:app /app

USER app

EXPOSE 8080
ENV PORT=8080

CMD ["./server"]
`
}

func dockerfileAPI() string {
	return `# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build binary
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server ./cmd/server

# Run stage
FROM alpine:3.19

RUN apk --no-cache add ca-certificates tzdata

# Non-root runtime user — Sentinel/Pulse open embedded SQLite stores under
# /app, so chown before USER or those fail with the misleading
# "out of memory (14)" SQLITE_CANTOPEN error.
RUN addgroup -S app && adduser -S -G app app

WORKDIR /app

COPY --from=builder /app/server .

RUN chown -R app:app /app

USER app

EXPOSE 8080

CMD ["./server"]
`
}

func dockerfileNextJS(app string) string {
	return fmt.Sprintf(`# Build stage
FROM node:22-alpine AS base

# Pin pnpm — pnpm@latest started resolving to pnpm 11 which needs Node 22's
# node:sqlite builtin. Pinning here avoids surprise breakage on rebuilds.
RUN corepack enable && corepack prepare pnpm@9.15.0 --activate

# Install dependencies
FROM base AS deps
WORKDIR /app

COPY package.json pnpm-lock.yaml pnpm-workspace.yaml ./
COPY apps/%s/package.json ./apps/%s/
COPY packages/shared/package.json ./packages/shared/

RUN pnpm install --frozen-lockfile

# Build
FROM base AS builder
WORKDIR /app

COPY --from=deps /app/node_modules ./node_modules
COPY --from=deps /app/apps/%s/node_modules ./apps/%s/node_modules
COPY --from=deps /app/packages/shared/node_modules ./packages/shared/node_modules
COPY . .

ARG NEXT_PUBLIC_API_URL
ENV NEXT_PUBLIC_API_URL=$NEXT_PUBLIC_API_URL

RUN pnpm --filter %s build

# Run
FROM base AS runner
WORKDIR /app

ENV NODE_ENV=production

RUN addgroup --system --gid 1001 nodejs
RUN adduser --system --uid 1001 nextjs

COPY --from=builder /app/apps/%s/.next/standalone ./
COPY --from=builder /app/apps/%s/.next/static ./apps/%s/.next/static
COPY --from=builder /app/apps/%s/public ./apps/%s/public

# chown before USER so the runtime user can write to /app (Next.js writes
# build manifests + cache here at startup).
RUN chown -R nextjs:nodejs /app

USER nextjs

EXPOSE 3000

ENV PORT=3000
ENV HOSTNAME="0.0.0.0"

CMD ["node", "apps/%s/server.js"]
`, app, app, app, app, app, app, app, app, app, app, app)
}

func dockerIgnore() string {
	return `node_modules
.next
.turbo
dist
*.log
.env
.env.local
.git
`
}
