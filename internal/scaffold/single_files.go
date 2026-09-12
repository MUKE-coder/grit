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
	// Use TanStack Router by default for single app (Vite produces static dist/)
	// Next.js can work via `next export` but TanStack/Vite is the natural fit
	// The shared Zod schemas and TS types, mirrored locally. Every other
	// architecture gets these as the packages/shared workspace package; a
	// single-binary app has no workspace, so they live under src/shared and
	// tsconfig aliases @repo/shared/* onto them (see singleFrontendTSConfig).
	files := singleSharedMirrorFiles(root, opts)
	for path, content := range singleFrontendOwnFiles(root, opts) {
		files[path] = content
	}

	for path, content := range files {
		if err := writeFile(path, content); err != nil {
			return fmt.Errorf("writing %s: %w", path, err)
		}
	}

	return writeBrandLogo(filepath.Join(filepath.Join(root, "frontend"), "public"), "grit_logo.png")
}

// singleSharedMirrorFiles is the shared package, mirrored into the SPA.
//
// Its own map so upgrade can deliver it to a project scaffolded before a file was
// added to it: the panel imports Zod schemas from here as values, and a mirror
// missing one of them is a build that fails on an import of "./money".
func singleSharedMirrorFiles(root string, opts Options) map[string]string {
	shared := filepath.Join(root, "frontend", "src", "shared")
	return map[string]string{
		filepath.Join(shared, "schemas", "user.ts"):     sharedUserSchema(),
		filepath.Join(shared, "schemas", "index.ts"):    sharedSchemasIndex(),
		filepath.Join(shared, "schemas", "blog.ts"):     sharedBlogSchema(),
		filepath.Join(shared, "schemas", "file-ref.ts"): sharedFileRefSchema(),
		// money and errors are exported by the barrels above, and were not
		// mirrored: a single project's schemas/index.ts re-exported ./money and
		// types/index.ts re-exported ./money and ./errors from files that were not
		// there. Nothing failed because the SPA's own use of this package is
		// type-only, which esbuild erases; the first value import from it, which
		// the admin panel brings, could not resolve.
		filepath.Join(shared, "schemas", "money.ts"):  sharedMoneySchema(),
		filepath.Join(shared, "types", "money.ts"):    sharedMoneyTypes(),
		filepath.Join(shared, "types", "errors.ts"):   sharedErrorsTS(),
		filepath.Join(shared, "types", "user.ts"):     sharedUserTypes(),
		filepath.Join(shared, "types", "api.ts"):      sharedAPITypes(),
		filepath.Join(shared, "types", "index.ts"):    sharedTypesIndex(),
		filepath.Join(shared, "types", "upload.ts"):   sharedUploadTypes(),
		filepath.Join(shared, "types", "blog.ts"):     sharedBlogTypes(),
		filepath.Join(shared, "types", "file-ref.ts"): sharedFileRefTypes(),
		filepath.Join(shared, "brand.config.ts"):      sharedBrandConfig(opts),
		filepath.Join(shared, "themes.ts"):            singleSharedThemes(opts),
	}
}

// singleFrontendOwnFiles is the SPA itself: its app shell, routes, components and
// configuration.
func singleFrontendOwnFiles(root string, opts Options) map[string]string {
	feRoot := filepath.Join(root, "frontend")
	return map[string]string{
		filepath.Join(feRoot, "package.json"):   singleFrontendPackageJSON(opts),
		filepath.Join(feRoot, "vite.config.ts"): singleFrontendViteConfig(),
		filepath.Join(feRoot, "index.html"):     webTanStackIndexHTML(opts),
		// .cjs (not .js) because package.json sets "type": "module" and PostCSS
		// config still uses CommonJS module.exports.
		filepath.Join(feRoot, "postcss.config.cjs"):          postCSSConfigFor(feRoot),
		filepath.Join(feRoot, "tsconfig.json"):               singleFrontendTSConfig(),
		filepath.Join(feRoot, "src", "main.tsx"):             webTanStackMain(),
		filepath.Join(feRoot, "src", "vite-env.d.ts"):        singleViteEnvTypes(),
		filepath.Join(feRoot, "src", "globals.css"):          webGlobalCSS(),
		filepath.Join(feRoot, "src", "routes", "__root.tsx"): webTanStackRootRoute(opts),
		// The public site, as a pathless layout route: the URLs are still / and
		// /blog, and the navbar and footer come from _site.tsx rather than from
		// the root deciding which paths deserve them.
		filepath.Join(feRoot, "src", "routes", "_site.tsx"):                  singleSiteLayoutRoute(),
		filepath.Join(feRoot, "src", "routes", "_site", "index.tsx"):         siteRouteID(webTanStackIndexRoute(opts)),
		filepath.Join(feRoot, "src", "routes", "_site", "blog", "index.tsx"): siteRouteID(webTanStackBlogListRoute()),
		filepath.Join(feRoot, "src", "routes", "_site", "blog", "$slug.tsx"): siteRouteID(webTanStackBlogDetailRoute()),
		// Vite-flavoured navbar/footer (use TanStack Router's <Link> + useRouterState),
		// not the Next.js variants from web_files.go which import next/link.
		filepath.Join(feRoot, "src", "components", "navbar.tsx"):    singleViteNavbar(opts),
		filepath.Join(feRoot, "src", "components", "footer.tsx"):    singleViteFooter(opts),
		filepath.Join(feRoot, "src", "components", "providers.tsx"): webTanStackProviders(),
		filepath.Join(feRoot, "src", "lib", "utils.ts"):             webUtils(),
		filepath.Join(feRoot, "components.json"):                    viteComponentsJSON(),
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
			log.Println("Empty database detected: running first-boot seed...")
			if err := database.Seed(db); err != nil {
				log.Printf("Warning: first-boot seed failed: %v", err)
			} else {
				log.Println("First-boot seed completed")
			}
		}
	}

	// Redis cache
	//
	// The driver's own logger goes first: without it, a project started with no
	// Redis running prints a wall of identical pool failures from inside go-redis
	// before any line of ours, and they look like a crash rather than a missing
	// optional service.
	cache.QuietDriverLogs()

	var cacheService *cache.Cache
	redisReachable := false
	if cfg.RedisURL != "" {
		c, err := cache.New(cfg.RedisURL)
		if err != nil {
			log.Printf("Redis is not reachable at %s: caching, background jobs and cron are off. Start it, or set REDIS_URL= in .env to run without it. (%v)", cfg.RedisURL, err)
		} else {
			cacheService = c
			redisReachable = true
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
	//
	// Only when Redis actually answered. jobs.NewClient parses the URL and builds
	// a client without connecting to anything, so "Job queue connected" used to
	// print on a machine with no Redis at all, right after the line saying Redis
	// was unreachable.
	var jobClient *jobs.Client
	if cfg.RedisURL != "" && redisReachable {
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
	//
	// Also only when Redis answered: the worker polls in a loop, so without Redis
	// it writes an asynq error every second or two, forever. That noise was the
	// worst part of starting a project with no Redis running, and it drowned out
	// the one line that explained it.
	var workerStop func()
	if cfg.RedisURL != "" && redisReachable {
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

	// Start cron scheduler, on the same condition and for the same reason.
	var cronScheduler *cron.Scheduler
	if cfg.RedisURL != "" && redisReachable {
		cs, cronErr := cron.Start(cfg, cacheService)
		if cronErr != nil {
			log.Printf("Warning: Cron scheduler failed to start: %v", cronErr)
		} else {
			cronScheduler = cs
		}
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

// singleSharedThemes is sharedThemes() with the two Next.js env reads swapped
// for Vite's. process is a Node global: it doesn't typecheck in a Vite app and
// is undefined in the browser, so the shipped source would silently fall back
// to the default theme even when THEME was set at scaffold time. The chosen
// theme is baked in as the fallback so the app renders correctly with no env
// at all. Panics rather than emitting a subtly wrong file if sharedThemes()
// ever drifts out from under these replacements.
func singleSharedThemes(opts Options) string {
	src := sharedThemes()

	replacements := [][2]string{
		{
			`typeof process !== "undefined" ? process.env.NEXT_PUBLIC_THEME : undefined`,
			fmt.Sprintf(`import.meta.env.VITE_THEME ?? %q`, opts.Theme),
		},
		{
			`typeof process !== "undefined"
    ? process.env.NEXT_PUBLIC_SOCIAL_AUTH_ENABLED
    : undefined`,
			`import.meta.env.VITE_SOCIAL_AUTH_ENABLED`,
		},
	}
	for _, r := range replacements {
		if !strings.Contains(src, r[0]) {
			panic("singleSharedThemes: sharedThemes() no longer contains: " + r[0])
		}
		src = strings.Replace(src, r[0], r[1], 1)
	}

	// The doc comments still describe the Next.js env names.
	src = strings.ReplaceAll(src, "process.env.NEXT_PUBLIC_THEME", "import.meta.env.VITE_THEME")
	src = strings.ReplaceAll(src, "NEXT_PUBLIC_SOCIAL_AUTH_ENABLED", "VITE_SOCIAL_AUTH_ENABLED")
	return src
}

// singleFrontendTSConfig is the Vite web tsconfig plus an alias for
// @repo/shared. A single-binary app has no pnpm workspace, so there is no
// packages/shared to resolve against; the shared types are mirrored into
// frontend/src/shared instead. Keeping the import specifier identical means
// the scaffolded hooks and anything `grit generate resource` emits work
// unchanged across all architectures.
func singleFrontendTSConfig() string {
	return strings.Replace(
		webTanStackTSConfig(),
		`"@/*": ["./src/*"]`,
		`"@/*": ["./src/*"],
      "@repo/upload/web": ["../packages/upload/src/web.ts"],
      "@repo/upload": ["../packages/upload/src/index.ts"],
      "@repo/shared/brand": ["./src/shared/brand.config.ts"],
      "@repo/shared/*": ["./src/shared/*"],
      "@admin/*": ["./src/admin-panel/*"]`,
		1,
	)
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
    "tailwind-merge": "^2.6.0",
    "zod": "^3.22.0",
    "@hookform/resolvers": "^3.3.0",
    "@react-pdf/renderer": "^4.1.5",
    "@tiptap/extension-color": "^2.1.0",
    "@tiptap/extension-highlight": "^2.1.0",
    "@tiptap/extension-image": "^2.1.0",
    "@tiptap/extension-link": "^2.1.0",
    "@tiptap/extension-placeholder": "^2.1.0",
    "@tiptap/extension-table": "^2.1.0",
    "@tiptap/extension-table-cell": "^2.1.0",
    "@tiptap/extension-table-header": "^2.1.0",
    "@tiptap/extension-table-row": "^2.1.0",
    "@tiptap/extension-text-align": "^2.1.0",
    "@tiptap/extension-text-style": "^2.1.0",
    "@tiptap/extension-underline": "^2.1.0",
    "@tiptap/pm": "^2.1.0",
    "@tiptap/react": "^2.1.0",
    "@tiptap/starter-kit": "^2.1.0",
    "react-dropzone": "^14.2.0",
    "react-hook-form": "^7.49.0",
    "recharts": "^2.12.0",
    "sonner": "^1.3.0",
    "tw-animate-css": "^1.4.0",
    "xlsx": "^0.18.5"
  },
  "devDependencies": {
    "@tanstack/react-router-devtools": "^1.93.0",
    "@tanstack/router-cli": "^1.93.0",
    "@tanstack/router-vite-plugin": "^1.93.0",
    "@types/react": "^19.0.0",
    "@types/react-dom": "^19.0.0",
    "@tailwindcss/postcss": "^4.1.13",
    "tw-animate-css": "^1.4.0",
    "postcss": "^8.4.49",
    "tailwindcss": "^4.1.13",
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
      // The admin panel's own code. tsconfig knows this alias too; Vite resolves
      // imports itself and needs telling separately, which is why the first build
      // of the embedded panel failed on every @admin import.
      '@admin': path.resolve(__dirname, './src/admin-panel'),
      // The mirrored shared package. tsconfig has aliased this since single
      // shipped, and Vite did not: the SPA's own use of it is type-only, which
      // esbuild erases before Rollup ever tries to resolve it, so nothing failed
      // until the admin panel imported a Zod schema from it as a value.
      // "./brand" is a subpath the shared package's exports map points at
      // brand.config.ts. An alias onto the directory resolves it to a file that
      // does not exist, which is what the admin sidebar's logo import hit.
      '@repo/shared/brand': path.resolve(__dirname, './src/shared/brand.config.ts'),
      '@repo/shared': path.resolve(__dirname, './src/shared'),
      // The upload package lives at the project root, outside the SPA, and is
      // raw TypeScript: Vite compiles it as source rather than resolving a build.
      // The admin panel's api-client builds its uploader from it.
      '@repo/upload/web': path.resolve(__dirname, '../packages/upload/src/web.ts'),
      '@repo/upload': path.resolve(__dirname, '../packages/upload/src/index.ts'),
    },
  },
  server: {
    port: 5173,
    proxy: {
      // Everything the Go binary serves. In production it serves this SPA too,
      // so these are all same-origin paths; in dev the SPA is on :5173 and they
      // have to be forwarded, or the admin panel's GORM Studio, Pulse and
      // Sentinel links land on the Vite dev server and 404.
      '/api': { target: 'http://localhost:8080', changeOrigin: true },
      '/studio': { target: 'http://localhost:8080', changeOrigin: true },
      '/pulse': { target: 'http://localhost:8080', changeOrigin: true },
      '/sentinel': { target: 'http://localhost:8080', changeOrigin: true },
      '/docs': { target: 'http://localhost:8080', changeOrigin: true },
    },
  },
  build: {
    outDir: 'dist',
  },
})
`
}

func singleMakefile(opts Options) string {
	return fmt.Sprintf(`# %s: Single App Makefile

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

// singleSiteLayoutRoute is the public site's layout: the navbar and the footer.
//
// Pathless, so it adds nothing to the URL. Its children are the landing page and
// the blog; the panel, the auth pages and the customer area are siblings with
// layouts of their own, which is why none of them can end up wearing this one.
func singleSiteLayoutRoute() string {
	return `import { createFileRoute, Outlet } from '@tanstack/react-router'

import { Navbar } from '@/components/navbar'
import { Footer } from '@/components/footer'

export const Route = createFileRoute('/_site')({
  component: SiteLayout,
})

function SiteLayout() {
  return (
    <div className="min-h-screen bg-background flex flex-col">
      <Navbar />
      <main className="flex-1">
        <Outlet />
      </main>
      <Footer />
    </div>
  )
}
`
}

// siteRouteID moves a route file into the _site section.
//
// TanStack checks that a file route's id equals its path, so a page that moves
// into a layout route has to say so: '/' becomes '/_site/', '/blog/$slug' becomes
// '/_site/blog/$slug'. The pages themselves are unchanged, which is the point:
// the same templates serve the monorepo web app, where there is no _site.
func siteRouteID(content string) string {
	return routeIDPattern.ReplaceAllStringFunc(content, func(match string) string {
		id := routeIDPattern.FindStringSubmatch(match)[1]
		if strings.HasPrefix(id, "/_site") {
			return match
		}
		return "createFileRoute('/_site" + id + "')"
	})
}
