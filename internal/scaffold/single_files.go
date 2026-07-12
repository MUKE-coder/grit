package scaffold

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// writeSingleMainGo writes the single-app main.go at the project root.
//
// We deliberately put main.go at the root (not under cmd/server/) so that the
// //go:embed all:frontend/dist directive resolves to <root>/frontend/dist —
// the same path that "pnpm build" emits. Putting main.go under cmd/server/
// makes go embed look for cmd/server/frontend/dist/* (relative to source
// file) which doesn't exist on a fresh clone and breaks `go build`.
//
// We also drop a tiny placeholder index.html into frontend/dist so that
// `go build` works on a fresh clone before the frontend has been built.
// `pnpm build` simply overwrites it.
func writeSingleMainGo(root string, opts Options) error {
	mainContent := singleMainGo(opts)
	mainContent = strings.ReplaceAll(mainContent, "{{MODULE}}", opts.Module())
	if err := writeFile(filepath.Join(root, "main.go"), mainContent); err != nil {
		return err
	}
	if err := writeFile(filepath.Join(root, "frontend", "dist", "index.html"), singleFrontendDistPlaceholder()); err != nil {
		return err
	}
	// writeAPIFiles seeds a multi-app cmd/server/main.go (sized for the
	// monorepo). In --single mode the canonical entry point is the root
	// main.go, so the leftover under cmd/server/ is a duplicate `package
	// main` that would break `go build ./...`. Remove it.
	cmdServerMain := filepath.Join(root, "cmd", "server", "main.go")
	if err := os.Remove(cmdServerMain); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("removing duplicate cmd/server/main.go: %w", err)
	}
	// Drop the now-empty cmd/server directory if nothing else lives in it.
	_ = os.Remove(filepath.Join(root, "cmd", "server"))
	return nil
}

// singleFrontendDistPlaceholder is the minimal index.html committed alongside
// the scaffold so `go build` succeeds on a fresh clone (before `pnpm build`).
// `pnpm build` overwrites this with the real Vite output.
func singleFrontendDistPlaceholder() string {
	return `<!doctype html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <title>Frontend not built</title>
  </head>
  <body>
    <p>Run <code>pnpm --filter ./frontend build</code> (or <code>make build</code>) to produce the real SPA bundle.</p>
  </body>
</html>
`
}

// writeSingleFrontendFiles writes the frontend scaffold inside frontend/ for single app.
func writeSingleFrontendFiles(root string, opts Options) error {
	feRoot := filepath.Join(root, "frontend")

	// Use TanStack Router by default for single app (Vite produces static dist/)
	// Next.js can work via `next export` but TanStack/Vite is the natural fit
	files := map[string]string{
		filepath.Join(feRoot, "package.json"):       singleFrontendPackageJSON(opts),
		filepath.Join(feRoot, "vite.config.ts"):     singleFrontendViteConfig(),
		filepath.Join(feRoot, "index.html"):         webTanStackIndexHTML(opts),
		filepath.Join(feRoot, "tailwind.config.ts"): webTanStackTailwindConfig(),
		// .cjs (not .js) because package.json sets "type": "module" and PostCSS
		// config still uses CommonJS module.exports.
		filepath.Join(feRoot, "postcss.config.cjs"):                 webPostCSSConfig(),
		filepath.Join(feRoot, "tsconfig.json"):                      webTanStackTSConfig(),
		filepath.Join(feRoot, "src", "main.tsx"):                    webTanStackMain(),
		filepath.Join(feRoot, "src", "vite-env.d.ts"):               singleViteEnvTypes(),
		filepath.Join(feRoot, "src", "globals.css"):                 webGlobalCSS(),
		filepath.Join(feRoot, "src", "routes", "__root.tsx"):        webTanStackRootRoute(opts),
		filepath.Join(feRoot, "src", "routes", "index.tsx"):         webTanStackIndexRoute(opts),
		filepath.Join(feRoot, "src", "routes", "blog", "index.tsx"): webTanStackBlogListRoute(),
		filepath.Join(feRoot, "src", "routes", "blog", "$slug.tsx"): webTanStackBlogDetailRoute(),
		// Vite-flavoured navbar/footer (use TanStack Router's <Link> + useRouterState),
		// not the Next.js variants from web_files.go which import next/link.
		filepath.Join(feRoot, "src", "components", "navbar.tsx"):    singleViteNavbar(opts),
		filepath.Join(feRoot, "src", "components", "footer.tsx"):    singleViteFooter(opts),
		filepath.Join(feRoot, "src", "components", "providers.tsx"): webTanStackProviders(),
		filepath.Join(feRoot, "src", "lib", "utils.ts"):             webUtils(),
		// viteAPIClientWithAuth (instead of viteAPIClient) so the axios
		// instance auto-attaches Authorization from the stored token AND
		// transparently refreshes on 401 — both gaps called out in the
		// real-world deployment review.
		filepath.Join(feRoot, "src", "lib", "api.ts"):         viteAPIClientWithAuth(),
		filepath.Join(feRoot, "src", "lib", "auth.ts"):        singleAuthLib(),
		filepath.Join(feRoot, "src", "hooks", "use-blogs.ts"): webUseBlogsHook(),
		// One file exports both ConfirmProvider (mount once at root) and the
		// useConfirm() hook. Import via: import { ConfirmProvider, useConfirm } from "@/hooks/use-confirm"
		filepath.Join(feRoot, "src", "hooks", "use-confirm.tsx"):                 singleUseConfirmHook(),
		filepath.Join(feRoot, "src", "components", "money-input.tsx"):            singleMoneyInput(),
		filepath.Join(feRoot, "src", "components", "combobox.tsx"):               singleCombobox(),
		filepath.Join(feRoot, "src", "components", "session-expiry-monitor.tsx"): singleSessionExpiryMonitor(),
		filepath.Join(feRoot, "src", "components", "status-badge.tsx"):           singleStatusBadge(),
		filepath.Join(feRoot, "src", "components", "stats-row.tsx"):              singleStatsRow(),
		filepath.Join(feRoot, "public", ".gitkeep"):                              "",
	}

	for path, content := range files {
		if err := writeFile(path, content); err != nil {
			return fmt.Errorf("writing %s: %w", path, err)
		}
	}

	if err := writeBrandLogo(filepath.Join(feRoot, "public"), "grit_logo.png"); err != nil {
		return err
	}

	return nil
}

// writeSingleRootFiles writes single-app specific root files (Makefile, README, .env).
func writeSingleRootFiles(root string, opts Options) error {
	files := map[string]string{
		filepath.Join(root, "Makefile"):   singleMakefile(opts),
		filepath.Join(root, ".gitignore"): singleGitignore(),
	}

	for path, content := range files {
		if err := writeFile(path, content); err != nil {
			return fmt.Errorf("writing %s: %w", path, err)
		}
	}

	return nil
}

func singleMainGo(opts Options) string {
	return `package main

import (
	"context"
	"embed"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	gothGithub "github.com/markbates/goth/providers/github"
	"github.com/markbates/goth/providers/google"

	"{{MODULE}}/internal/ai"
	"{{MODULE}}/internal/cache"
	"{{MODULE}}/internal/config"
	"{{MODULE}}/internal/cron"
	"{{MODULE}}/internal/database"
	"{{MODULE}}/internal/jobs"
	"{{MODULE}}/internal/mail"
	"{{MODULE}}/internal/models"
	"{{MODULE}}/internal/routes"
	"{{MODULE}}/internal/storage"
)

//go:embed all:frontend/dist
var frontendFS embed.FS

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Connect to database
	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Auto-migrate (idempotent — GORM AutoMigrate is additive).
	// Disable with AUTO_MIGRATE=false for environments that run migrations
	// out-of-band via "./migrate".
	if strings.ToLower(os.Getenv("AUTO_MIGRATE")) != "false" {
		log.Println("Running auto-migrations...")
		if err := models.Migrate(db); err != nil {
			log.Fatalf("Auto-migration failed: %v", err)
		}
	}

	// First-boot seed: only runs when the users table is empty. Off by default
	// in production unless AUTO_SEED=true is set explicitly.
	autoSeed := strings.ToLower(os.Getenv("AUTO_SEED"))
	seedEnabled := autoSeed == "true" || (autoSeed != "false" && cfg.AppEnv != "production")
	if seedEnabled {
		var userCount int64
		db.Model(&models.User{}).Count(&userCount)
		if userCount == 0 {
			log.Println("Empty database detected — running first-boot seed...")
			if err := database.Seed(db); err != nil {
				log.Printf("Warning: first-boot seed failed: %v", err)
			} else {
				log.Println("First-boot seed completed")
			}
		}
	}

	// Redis cache
	var cacheService *cache.Cache
	if cfg.RedisURL != "" {
		c, err := cache.New(cfg.RedisURL)
		if err != nil {
			log.Printf("Warning: Redis unavailable: %v", err)
		} else {
			cacheService = c
			log.Println("Redis cache connected")
		}
	}

	// S3-compatible storage
	var storageService *storage.Storage
	s, err := storage.New(cfg.Storage)
	if err != nil {
		log.Printf("Warning: Storage unavailable: %v", err)
	} else {
		storageService = s
		log.Println("Storage configured")
	}

	// Email (Resend)
	var mailer *mail.Mailer
	if cfg.ResendAPIKey != "" && cfg.ResendAPIKey != "re_your_api_key" {
		mailer = mail.New(cfg.ResendAPIKey, cfg.MailFrom)
		log.Println("Email service configured")
	}

	// AI service (Vercel AI Gateway)
	var aiService *ai.AI
	if cfg.AIGatewayAPIKey != "" {
		aiService = ai.New(cfg.AIGatewayAPIKey, cfg.AIGatewayModel, cfg.AIGatewayURL)
		log.Printf("AI service configured via AI Gateway (%s)", cfg.AIGatewayModel)
	}

	// Background jobs (asynq) — client (enqueue side)
	var jobClient *jobs.Client
	if cfg.RedisURL != "" {
		jc, err := jobs.NewClient(cfg.RedisURL)
		if err != nil {
			log.Printf("Warning: Job queue unavailable: %v", err)
		} else {
			jobClient = jc
			log.Println("Job queue connected")
		}
	}

	// OAuth2 social login providers
	gothic.Store = sessions.NewCookieStore([]byte(cfg.JWTSecret))
	var oauthProviders []goth.Provider
	if cfg.GoogleClientID != "" {
		oauthProviders = append(oauthProviders, google.New(
			cfg.GoogleClientID, cfg.GoogleClientSecret,
			cfg.AppURL+"/api/auth/oauth/google/callback",
		))
	}
	if cfg.GithubClientID != "" {
		oauthProviders = append(oauthProviders, gothGithub.New(
			cfg.GithubClientID, cfg.GithubClientSecret,
			cfg.AppURL+"/api/auth/oauth/github/callback",
		))
	}
	if len(oauthProviders) > 0 {
		goth.UseProviders(oauthProviders...)
	}

	// Build services
	svc := &routes.Services{
		Cache:   cacheService,
		Storage: storageService,
		Mailer:  mailer,
		AI:      aiService,
		Jobs:    jobClient,
	}

	// Setup router
	router := routes.Setup(db, cfg, svc)

	// Serve embedded frontend (SPA fallback).
	// We pre-read index.html once and serve it via c.Data() to avoid the
	// canonical-URL redirect rule in http.FileServer, which causes
	// ERR_TOO_MANY_REDIRECTS behind reverse proxies (Traefik / Cloudflare).
	feFS, err := fs.Sub(frontendFS, "frontend/dist")
	if err != nil {
		log.Printf("Warning: embedded frontend not available: %v", err)
	} else {
		indexHTML, readErr := fs.ReadFile(feFS, "index.html")
		if readErr != nil {
			log.Printf("Warning: failed to read embedded index.html: %v", readErr)
		}
		fileServer := http.FileServer(http.FS(feFS))
		router.NoRoute(func(c *gin.Context) {
			path := c.Request.URL.Path
			// Try to serve a real static asset first.
			if path != "/" && path != "/index.html" {
				if f, err := feFS.Open(path[1:]); err == nil {
					f.Close()
					fileServer.ServeHTTP(c.Writer, c.Request)
					return
				}
			}
			// SPA fallback: hand back the pre-read index.html.
			if indexHTML == nil {
				c.Status(http.StatusNotFound)
				return
			}
			c.Data(http.StatusOK, "text/html; charset=utf-8", indexHTML)
		})
	}

	// Start asynq worker (consumes the queue jobClient enqueues to).
	var workerStop func()
	if cfg.RedisURL != "" {
		stop, err := jobs.StartWorker(cfg.RedisURL, jobs.WorkerDeps{
			DB:      db,
			Mailer:  mailer,
			Storage: storageService,
			Cache:   cacheService,
		})
		if err != nil {
			log.Printf("Warning: Background worker failed to start: %v", err)
		} else {
			workerStop = stop
			log.Println("Background worker started")
		}
	}

	// Start cron scheduler
	cronScheduler, err := cron.Start(cfg, cacheService)
	if err != nil {
		log.Printf("Warning: Cron scheduler failed to start: %v", err)
	}

	// Start server
	port := cfg.Port
	if port == "" {
		port = "8080"
	}

	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("Server starting on :%s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down...")

	if cronScheduler != nil {
		cronScheduler.Stop()
	}
	if workerStop != nil {
		workerStop()
	}
	if jobClient != nil {
		_ = jobClient.Close()
	}
	if cacheService != nil {
		_ = cacheService.Close()
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Server shutdown error: %v", err)
	}
	log.Println("Server stopped")
}
`
}

func singleFrontendPackageJSON(opts Options) string {
	// postinstall + "routes:generate" wire @tanstack/router-cli so that
	// routeTree.gen.ts exists even before `pnpm dev` has been run once.
	// Without it, `tsc --noEmit` fails on a fresh clone because the
	// generated file isn't there. The Vite plugin keeps it regenerated
	// during dev/build.
	return fmt.Sprintf(`{
  "name": "%s-frontend",
  "version": "0.1.0",
  "private": true,
  "type": "module",
  "scripts": {
    "dev": "vite",
    "build": "tsr generate && tsc -b && vite build",
    "preview": "vite preview",
    "routes:generate": "tsr generate",
    "postinstall": "tsr generate || true"
  },
  "dependencies": {
    "@tanstack/react-query": "^5.62.0",
    "@tanstack/react-router": "^1.93.0",
    "axios": "^1.7.9",
    "clsx": "^2.1.1",
    "lucide-react": "^0.468.0",
    "react": "19.2.7",
    "react-dom": "19.2.7",
    "tailwind-merge": "^2.6.0"
  },
  "devDependencies": {
    "@tanstack/react-router-devtools": "^1.93.0",
    "@tanstack/router-cli": "^1.93.0",
    "@tanstack/router-vite-plugin": "^1.93.0",
    "@types/react": "^19.0.0",
    "@types/react-dom": "^19.0.0",
    "autoprefixer": "^10.4.20",
    "postcss": "^8.4.49",
    "tailwindcss": "^3.4.17",
    "typescript": "~5.7.0",
    "vite": "^6.0.0",
    "@vitejs/plugin-react": "^4.3.4"
  }
}`, opts.ProjectName)
}

func singleFrontendViteConfig() string {
	return `import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import { TanStackRouterVite } from '@tanstack/router-vite-plugin'
import path from 'path'

export default defineConfig({
  plugins: [
    TanStackRouterVite(),
    react(),
  ],
  resolve: {
    alias: {
      '@': path.resolve(__dirname, './src'),
    },
  },
  server: {
    port: 5173,
    proxy: {
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      },
    },
  },
  build: {
    outDir: 'dist',
  },
})
`
}

func singleMakefile(opts Options) string {
	return fmt.Sprintf(`# %s — Single App Makefile

.PHONY: dev build migrate seed run clean

# Development: run Go API + Vite frontend in parallel.
dev:
	@echo "Starting development servers..."
	@cd frontend && pnpm dev &
	@air

# Build production binary (embeds frontend/dist via //go:embed).
# main.go lives at the project root so the embed path resolves correctly
# without any pre-build copy step.
build:
	@echo "Building frontend..."
	@cd frontend && pnpm install && pnpm build
	@echo "Building Go binary..."
	@go build -o bin/%s .
	@echo "Done! Binary at bin/%s"

# One-shot migrate + seed (handy for local resets).
migrate:
	@go run ./cmd/migrate

seed:
	@go run ./cmd/seed

# Run the built binary.
run: build
	@./bin/%s

# Clean build artifacts (keeps the dist placeholder so go build still works).
clean:
	@rm -rf bin/
	@rm -rf frontend/dist/assets
`, opts.ProjectName, opts.ProjectName, opts.ProjectName, opts.ProjectName)
}

func singleGitignore() string {
	return `# Go
bin/
*.exe
*.exe~
*.dll
*.so
*.dylib
*.test
*.out

# Frontend
frontend/node_modules/
frontend/dist/
frontend/.vite/

# Environment
.env
.env.local

# IDE
.idea/
.vscode/
*.swp
*.swo

# OS
.DS_Store
Thumbs.db

# Air
tmp/
`
}
