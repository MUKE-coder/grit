package scaffold

import (
	"fmt"
	"path/filepath"
)

func writeDockerFiles(root string, opts Options) error {
	files := map[string]string{
		filepath.Join(root, "docker-compose.yml"):      dockerCompose(opts),
		filepath.Join(root, "docker-compose.prod.yml"): dockerComposeProd(opts),
		filepath.Join(root, "apps", "api", "Dockerfile"): dockerfileAPI(),
	}

	if opts.ShouldIncludeWeb() {
		files[filepath.Join(root, "apps", "web", "Dockerfile")] = dockerfileNextJS("web")
	}
	if opts.ShouldIncludeAdmin() {
		files[filepath.Join(root, "apps", "admin", "Dockerfile")] = dockerfileNextJS("admin")
	}
	if opts.ShouldIncludeDocs() {
		files[filepath.Join(root, "apps", "docs", "Dockerfile")] = dockerfileNextJS("docs")
	}

	if !opts.APIOnly {
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
	return fmt.Sprintf(`services:
  postgres:
    image: postgres:16-alpine
    container_name: %s-postgres
    restart: unless-stopped
    ports:
      - "5432:5432"
    environment:
      POSTGRES_USER: grit
      POSTGRES_PASSWORD: grit
      POSTGRES_DB: %s
    volumes:
      - postgres-data:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U grit"]
      interval: 5s
      timeout: 5s
      retries: 5

  redis:
    image: redis:7-alpine
    container_name: %s-redis
    restart: unless-stopped
    ports:
      - "6379:6379"
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
      - "9000:9000"
      - "9001:9001"
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
      - "1025:1025"
      - "8025:8025"

volumes:
  postgres-data:
  redis-data:
  minio-data:
`, opts.ProjectName, opts.ProjectName, opts.ProjectName, opts.ProjectName, opts.ProjectName)
}

func dockerComposeProd(opts Options) string {
	name := opts.ProjectName

	result := fmt.Sprintf(`services:
  api:
    build:
      context: ./apps/api
      dockerfile: Dockerfile
    container_name: %s-api
    restart: unless-stopped
    ports:
      - "8080:8080"
    environment:
      APP_ENV: production
      DATABASE_URL: postgres://grit:grit@postgres:5432/%s?sslmode=disable
      REDIS_URL: redis://redis:6379
      JWT_SECRET: ${JWT_SECRET}
    depends_on:
      postgres:
        condition: service_healthy
      redis:
        condition: service_healthy
`, name, name)

	if opts.ShouldIncludeWeb() {
		result += fmt.Sprintf(`
  web:
    build:
      context: .
      dockerfile: apps/web/Dockerfile
    container_name: %s-web
    restart: unless-stopped
    ports:
      - "3000:3000"
    environment:
      NEXT_PUBLIC_API_URL: http://api:8080
`, name)
	}

	if opts.ShouldIncludeAdmin() {
		result += fmt.Sprintf(`
  admin:
    build:
      context: .
      dockerfile: apps/admin/Dockerfile
    container_name: %s-admin
    restart: unless-stopped
    ports:
      - "3001:3000"
    environment:
      NEXT_PUBLIC_API_URL: http://api:8080
`, name)
	}

	if opts.ShouldIncludeDocs() {
		result += fmt.Sprintf(`
  docs:
    build:
      context: .
      dockerfile: apps/docs/Dockerfile
    container_name: %s-docs
    restart: unless-stopped
    ports:
      - "3002:3002"
`, name)
	}

	result += fmt.Sprintf(`
  postgres:
    image: postgres:16-alpine
    container_name: %s-postgres
    restart: unless-stopped
    environment:
      POSTGRES_USER: grit
      POSTGRES_PASSWORD: grit
      POSTGRES_DB: %s
    volumes:
      - postgres-data:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U grit"]
      interval: 5s
      timeout: 5s
      retries: 5

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

volumes:
  postgres-data:
  redis-data:
`, name, name, name)

	return result
}

func dockerfileAPI() string {
	return `# Build stage
FROM golang:1.23-alpine AS builder

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

WORKDIR /app

COPY --from=builder /app/server .

EXPOSE 8080

CMD ["./server"]
`
}

func dockerfileNextJS(app string) string {
	return fmt.Sprintf(`# Build stage
FROM node:20-alpine AS base

RUN corepack enable

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
