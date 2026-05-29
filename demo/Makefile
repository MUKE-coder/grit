# stoka — Single App Makefile

.PHONY: dev build clean

# Development: run Go API + Vite frontend in parallel
dev:
	@echo "Starting development servers..."
	@cd frontend && pnpm dev &
	@air

# Build production binary (embeds frontend)
build:
	@echo "Building frontend..."
	@cd frontend && pnpm install && pnpm build
	@echo "Building Go binary..."
	@go build -o bin/stoka .
	@echo "Done! Binary at bin/stoka"

# Clean build artifacts
clean:
	@rm -rf bin/ tmp/ frontend/dist/
