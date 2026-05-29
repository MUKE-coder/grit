# Stage 1: Build frontend
FROM node:20-alpine AS frontend
WORKDIR /app/frontend
RUN corepack enable && corepack prepare pnpm@latest --activate
COPY frontend/package.json frontend/pnpm-lock.yaml ./
RUN pnpm install --frozen-lockfile
COPY frontend/ ./

# Vite needs env vars at build time — pass via build args
ARG VITE_INTERNAL_API_KEY
ARG VITE_APP_NAME=Kayula Motors
ENV VITE_INTERNAL_API_KEY=$VITE_INTERNAL_API_KEY
ENV VITE_APP_NAME=$VITE_APP_NAME

RUN pnpm build

# Stage 2: Build Go binary
FROM golang:1.24-alpine AS backend
RUN apk add --no-cache git
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
COPY --from=frontend /app/frontend/dist ./frontend/dist
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o stoka .

# Stage 3: Final minimal image
FROM alpine:3.20
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /app
COPY --from=backend /app/stoka .
EXPOSE 8080
CMD ["./stoka"]
