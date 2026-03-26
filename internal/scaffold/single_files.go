package scaffold

import (
	"fmt"
	"path/filepath"
	"strings"
)

// writeSingleMainGo writes the root main.go with go:embed for the single app architecture.
// This replaces the apps/api/cmd/server/main.go with one that also embeds the frontend.
func writeSingleMainGo(root string, opts Options) error {
	mainContent := singleMainGo(opts)
	mainContent = strings.ReplaceAll(mainContent, "{{MODULE}}", opts.Module())
	return writeFile(filepath.Join(root, "cmd", "server", "main.go"), mainContent)
}

// writeSingleFrontendFiles writes the frontend scaffold inside frontend/ for single app.
func writeSingleFrontendFiles(root string, opts Options) error {
	feRoot := filepath.Join(root, "frontend")

	// Use TanStack Router by default for single app (Vite produces static dist/)
	// Next.js can work via `next export` but TanStack/Vite is the natural fit
	files := map[string]string{
		filepath.Join(feRoot, "package.json"):                              singleFrontendPackageJSON(opts),
		filepath.Join(feRoot, "vite.config.ts"):                           singleFrontendViteConfig(),
		filepath.Join(feRoot, "index.html"):                               webTanStackIndexHTML(opts),
		filepath.Join(feRoot, "tailwind.config.ts"):                       webTanStackTailwindConfig(),
		filepath.Join(feRoot, "postcss.config.js"):                        webPostCSSConfig(),
		filepath.Join(feRoot, "tsconfig.json"):                            webTanStackTSConfig(),
		filepath.Join(feRoot, "src", "main.tsx"):                          webTanStackMain(),
		filepath.Join(feRoot, "src", "globals.css"):                       webGlobalCSS(),
		filepath.Join(feRoot, "src", "routes", "__root.tsx"):              webTanStackRootRoute(opts),
		filepath.Join(feRoot, "src", "routes", "index.tsx"):               webTanStackIndexRoute(opts),
		filepath.Join(feRoot, "src", "routes", "blog", "index.tsx"):       webTanStackBlogListRoute(),
		filepath.Join(feRoot, "src", "routes", "blog", "$slug.tsx"):       webTanStackBlogDetailRoute(),
		filepath.Join(feRoot, "src", "components", "navbar.tsx"):          webNavbar(opts),
		filepath.Join(feRoot, "src", "components", "footer.tsx"):          webFooter(opts),
		filepath.Join(feRoot, "src", "components", "providers.tsx"):       webTanStackProviders(),
		filepath.Join(feRoot, "src", "lib", "utils.ts"):                   webUtils(),
		filepath.Join(feRoot, "src", "lib", "api.ts"):                     webAPIClient(),
		filepath.Join(feRoot, "src", "hooks", "use-blogs.ts"):             webUseBlogsHook(),
		filepath.Join(feRoot, "public", ".gitkeep"):                       "",
	}

	for path, content := range files {
		if err := writeFile(path, content); err != nil {
			return fmt.Errorf("writing %s: %w", path, err)
		}
	}

	return nil
}

// writeSingleRootFiles writes single-app specific root files (Makefile, README, .env).
func writeSingleRootFiles(root string, opts Options) error {
	files := map[string]string{
		filepath.Join(root, "Makefile"):  singleMakefile(opts),
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
	"{{MODULE}}/internal/routes"
	"{{MODULE}}/internal/storage"
)

//go:embed frontend/dist/*
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

	// Background jobs (asynq)
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

	// Serve embedded frontend (SPA fallback)
	feFS, err := fs.Sub(frontendFS, "frontend/dist")
	if err != nil {
		log.Printf("Warning: embedded frontend not available: %v", err)
	} else {
		// Serve static files
		fileServer := http.FileServer(http.FS(feFS))
		router.NoRoute(func(c *gin.Context) {
			// Try to serve static file first
			path := c.Request.URL.Path
			if f, err := feFS.Open(path[1:]); err == nil {
				f.Close()
				fileServer.ServeHTTP(c.Writer, c.Request)
				return
			}
			// SPA fallback: serve index.html for all non-API routes
			c.FileFromFS("index.html", http.FS(feFS))
		})
	}

	// Start cron scheduler
	cron.Start(cfg, cacheService)

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
	return fmt.Sprintf(`{
  "name": "%s-frontend",
  "version": "0.1.0",
  "private": true,
  "type": "module",
  "scripts": {
    "dev": "vite",
    "build": "tsc -b && vite build",
    "preview": "vite preview"
  },
  "dependencies": {
    "@tanstack/react-query": "^5.62.0",
    "@tanstack/react-router": "^1.93.0",
    "axios": "^1.7.9",
    "clsx": "^2.1.1",
    "lucide-react": "^0.468.0",
    "react": "^19.0.0",
    "react-dom": "^19.0.0",
    "tailwind-merge": "^2.6.0"
  },
  "devDependencies": {
    "@tanstack/react-router-devtools": "^1.93.0",
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
	@go build -o bin/%s ./cmd/server
	@echo "Done! Binary at bin/%s"

# Clean build artifacts
clean:
	@rm -rf bin/ frontend/dist/
`, opts.ProjectName, opts.ProjectName, opts.ProjectName)
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
