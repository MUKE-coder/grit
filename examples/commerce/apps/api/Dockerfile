# Build stage
#
# IMAGE_BACKEND selects the image optimisation backend.
#
#   purego  (default) no system libraries, static binary, JPEG + lossless WebP
#   vips              libvips: lossy WebP and AVIF, and files roughly 4x
#                     smaller on a photograph, at the cost of cgo and a
#                     runtime dependency on libvips
#
#   docker build --build-arg IMAGE_BACKEND=vips -f Dockerfile.api .
#
# Docker is the right place to opt in: the image controls the environment, so
# libvips is guaranteed present at build and at run. "grit deploy" cannot do
# this, because it cross-compiles to linux/amd64 with CGO_ENABLED=0 from your
# machine, and cgo cannot cross-compile without a target toolchain.
ARG IMAGE_BACKEND=purego

FROM golang:1.26.6-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build binary
ARG IMAGE_BACKEND
RUN if [ "$IMAGE_BACKEND" = "vips" ]; then       apk --no-cache add vips-dev build-base pkgconfig &&       CGO_ENABLED=1 GOOS=linux go build -tags vips -o /app/server ./cmd/server ;     else       CGO_ENABLED=0 GOOS=linux go build -o /app/server ./cmd/server ;     fi

# Run stage
FROM alpine:3.24
ARG IMAGE_BACKEND

RUN apk --no-cache add ca-certificates tzdata

# The vips build links against libvips, so the runtime image needs it too. The
# default build is static and adds nothing.
RUN if [ "$IMAGE_BACKEND" = "vips" ]; then apk --no-cache add vips; fi

# Packages a plugin needs at runtime go here, such as ffmpeg for grit plugin
# add video.
# grit:runtime-packages

# Non-root runtime user — Sentinel/Pulse open embedded SQLite stores under
# /app, so chown before USER or those fail with the misleading
# "out of memory (14)" SQLITE_CANTOPEN error.
RUN addgroup -S app && adduser -S -G app app

WORKDIR /app

COPY --from=builder /app/server .

RUN chown -R app:app /app

USER app

EXPOSE 8080

# Unhealthy when the API stops answering. /api/health answers 200 while the
# server is serving, and reports each dependency in its body.
HEALTHCHECK --interval=30s --timeout=5s --start-period=30s --retries=3 \
  CMD wget -q -O /dev/null "http://127.0.0.1:${APP_PORT:-8080}/api/health" || exit 1

CMD ["./server"]
