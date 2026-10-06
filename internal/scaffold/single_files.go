package scaffold

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// writeSingleMainGo writes the single-app main.go at the project root.
//
// main.go sits at the top of api/ rather than under api/cmd/server/, because
// //go:embed is resolved relative to the source file: from cmd/server it would
// look for cmd/server/web/, which nothing writes. From api/ it finds api/web,
// which is where Vite is told to build.
//
// A placeholder index.html goes in first so `go build` works on a fresh clone,
// before anyone has built the frontend. `pnpm build` overwrites it.
func writeSingleMainGo(root string, opts Options) error {
	// A Next single has nothing to embed.
	//
	// Both singles embed their frontend, so both get this entry point. A Vite
	// build is static files; a Next build is too, with output: "export". What
	// differs is only the shape of what lands in api/web: one index.html for the
	// SPA, one HTML file per route for the export, and the handler in here reads
	// both.
	mainContent := singleMainGo(opts)
	api := opts.APIRoot(root)
	mainContent = strings.ReplaceAll(mainContent, "{{MODULE}}", opts.Module())
	if err := writeFile(filepath.Join(api, "main.go"), mainContent); err != nil {
		return err
	}
	// A placeholder so `go build` works on a fresh clone, before anyone has built
	// the frontend: //go:embed fails the build outright on a pattern that matches
	// nothing. Both frontends overwrite it, Vite by building into api/web and
	// Next by copying its export there.
	//
	// It lands inside api/ because //go:embed cannot reach above the directory
	// its source file is in: a main.go in api/ embedding ../out does not compile.
	if err := writeFile(filepath.Join(api, "web", "index.html"), singleFrontendDistPlaceholder()); err != nil {
		return err
	}
	// writeAPIFiles seeds a multi-app cmd/server/main.go (sized for the
	// monorepo). In --single mode the canonical entry point is the root
	// main.go, so the leftover under cmd/server/ is a duplicate `package
	// main` that would break `go build ./...`. Remove it.
	cmdServerMain := filepath.Join(api, "cmd", "server", "main.go")
	if err := os.Remove(cmdServerMain); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("removing duplicate cmd/server/main.go: %w", err)
	}
	// Drop the now-empty cmd/server directory if nothing else lives in it.
	_ = os.Remove(filepath.Join(api, "cmd", "server"))
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
    <p>Run <code>pnpm build</code> (or <code>make build</code>) to produce the real SPA bundle.</p>
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

	return writeBrandLogo(filepath.Join(webAppRoot(root, opts), "public"), "grit_logo.png")
}

// singleSharedMirrorFiles is the shared package, mirrored into the SPA.
//
// Its own map so upgrade can deliver it to a project scaffolded before a file was
// added to it: the panel imports Zod schemas from here as values, and a mirror
// missing one of them is a build that fails on an import of "./money".
func singleSharedMirrorFiles(root string, opts Options) map[string]string {
	shared := filepath.Join(webAppRoot(root, opts), "src", "shared")
	files := map[string]string{
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
		filepath.Join(shared, "schemas", "money.ts"):   sharedMoneySchema(),
		filepath.Join(shared, "types", "money.ts"):     sharedMoneyTypes(),
		filepath.Join(shared, "types", "errors.ts"):    sharedErrorsTS(),
		filepath.Join(shared, "types", "user.ts"):      sharedUserTypes(),
		filepath.Join(shared, "types", "api.ts"):       sharedAPITypes(),
		filepath.Join(shared, "types", "index.ts"):     sharedTypesIndex(),
		filepath.Join(shared, "types", "upload.ts"):    sharedUploadTypes(),
		filepath.Join(shared, "types", "blog.ts"):      sharedBlogTypes(),
		filepath.Join(shared, "types", "file-ref.ts"):  sharedFileRefTypes(),
		filepath.Join(shared, "brand.config.ts"):       sharedBrandConfig(opts),
		filepath.Join(shared, "themes.ts"):             singleSharedThemes(opts),
		filepath.Join(shared, "schemas", "profile.ts"): sharedProfileSchemas(),
	}
	for path, body := range sharedModelTypeFiles(shared) {
		files[path] = body
	}
	for path, body := range sharedFieldFormatFiles(shared) {
		files[path] = body
	}
	return files
}

// singleFrontendOwnFiles is the SPA itself: its app shell, routes, components and
// configuration.
func singleFrontendOwnFiles(root string, opts Options) map[string]string {
	feRoot := webAppRoot(root, opts)
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
	"errors"
	"crypto/sha256"
	"embed"
	"io/fs"
	"log"
	"net/http"
	gopath "path"
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
	"{{MODULE}}/internal/events"
	"{{MODULE}}/internal/jobs"
	"{{MODULE}}/internal/mail"
	"{{MODULE}}/internal/models"
	"{{MODULE}}/internal/routes"
	"{{MODULE}}/internal/services"
	"{{MODULE}}/internal/storage"
)

//go:embed all:web
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

` + singleAutoSeedCommentNew + `	autoSeed := strings.ToLower(os.Getenv("AUTO_SEED"))
	` + singleAutoSeedNew + `
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

` + mainStorageInit + `
` + mailerInitNew + `
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
` + gothicStoreNew + `	var oauthProviders []goth.Provider
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

	// Serve the embedded frontend.
	//
	// Two shapes end up in here. A Vite build is a SPA: one index.html, and the
	// router in the browser resolves every path. A Next export is one HTML file
	// per route, so /admin/dashboard is admin/dashboard.html, and handing back
	// the root index.html for it would serve the wrong page with a 200.
	//
	// The lookups are ordered, and the two HTML steps are misses for a SPA, which
	// has no per-route files: it falls through to index.html exactly as it did
	// before. One handler for both, rather than two that drift apart.
	//
	// index.html is pre-read and served with c.Data rather than through
	// http.FileServer, whose canonical-URL redirect causes ERR_TOO_MANY_REDIRECTS
	// behind Traefik and Cloudflare.
	feFS, err := fs.Sub(frontendFS, "web")
	if err != nil {
		log.Printf("Warning: embedded frontend not available: %v", err)
	} else {
		indexHTML, readErr := fs.ReadFile(feFS, "index.html")
		if readErr != nil {
			log.Printf("Warning: failed to read embedded index.html: %v", readErr)
		}
		fileServer := http.FileServer(http.FS(feFS))

		serveFile := func(c *gin.Context, name string) bool {
			data, readErr := fs.ReadFile(feFS, name)
			if readErr != nil {
				return false
			}
			c.Data(http.StatusOK, "text/html; charset=utf-8", data)
			return true
		}

		router.NoRoute(func(c *gin.Context) {
			if c.Request.Method != http.MethodGet && c.Request.Method != http.MethodHead {
				c.Status(http.StatusNotFound)
				return
			}
			name := strings.TrimPrefix(gopath.Clean(c.Request.URL.Path), "/")

			// A real asset: JS, CSS, an image. FileServer sets the headers and
			// handles range requests.
			if name != "" && name != "index.html" {
				if st, statErr := fs.Stat(feFS, name); statErr == nil && !st.IsDir() {
					fileServer.ServeHTTP(c.Writer, c.Request)
					return
				}
			}

			// A clean URL, which a Next export has written as a file.
			if name != "" && gopath.Ext(name) == "" {
				if serveFile(c, name+".html") {
					return
				}
				if serveFile(c, name+"/index.html") {
					return
				}
			}

			// The SPA, or the export's own root.
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
		Addr:    ":" + port,
		Handler: router,
` + serverTimeoutFields + `	}

	go func() {
		log.Printf("Server starting on :%s", port)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
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
` + shutdownDrain + `	log.Println("Server stopped")
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
		// The single's own aliases are added below, so the base config is enough here.
		webTanStackTSConfig(Options{}),
		`"@/*": ["./src/*"]`,
		`"@/*": ["./src/*"],
      "@repo/upload/web": ["./packages/upload/src/web.ts"],
      "@repo/upload": ["./packages/upload/src/index.ts"],
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
  "name": "%s",
  "version": "0.1.0",
  "private": true,
  "type": "module",
  "scripts": {
    "dev": "vite",
    "build": "tsr generate && tsc -b && vite build",
    "preview": "vite preview",
    "routes:generate": "tsr generate",
    "postinstall": "tsr generate || true",
    "lint": "`+biomeLintScript+`",
    "format": "`+biomeFormatScript+`"
  },
  "dependencies": {
`+fieldInputDependencyLines("    ")+`    "@tanstack/react-query": "^5.62.0",
    "@tanstack/react-router": "^1.93.0",
    "axios": "^1.7.9",
    "clsx": "^2.1.1",
    "dompurify": "^3.4.15",
    "lucide-react": "^0.468.0",
    "react": "19.2.7",
    "react-dom": "19.2.7",
    "tailwind-merge": "^2.6.0",
    "zod": "^3.22.0",
    "@hookform/resolvers": "^3.3.0",
    "@react-pdf/renderer": "^4.1.5",
`+tiptapDependencyLines("    ")+`    "react-dropzone": "^14.2.0",
    "react-hook-form": "^7.49.0",
    "recharts": "^2.12.0",
    "sonner": "^1.3.0",
    "tw-animate-css": "^1.4.0",
    "xlsx": "https://cdn.sheetjs.com/xlsx-0.20.3/xlsx-0.20.3.tgz"
  },
  "devDependencies": {
    `+biomeDevDependency+`,
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
	return `import { defineConfig, loadEnv } from 'vite'
import react from '@vitejs/plugin-react'
import { TanStackRouterVite } from '@tanstack/router-vite-plugin'
import path from 'node:path'

// Where the Go API listens, from .env, so the dev proxy and the server agree.
// loadEnv reads .env from the project root; APP_PORT is the same variable the Go
// binary reads, so changing it in one place moves both.
const env = loadEnv(process.env.NODE_ENV || 'development', process.cwd(), '')
const apiTarget = 'http://localhost:' + (env.APP_PORT || '8080')

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
      '@repo/upload/web': path.resolve(__dirname, './packages/upload/src/web.ts'),
      '@repo/upload': path.resolve(__dirname, './packages/upload/src/index.ts'),
    },
  },
  server: {
    port: 5173,
    proxy: {
      // Everything the Go binary serves. In production it serves this SPA too,
      // so these are all same-origin paths; in dev the SPA is on :5173 and they
      // have to be forwarded, or the admin panel's GORM Studio, Pulse and
      // Sentinel links land on the Vite dev server and 404.
      //
      // The port is read from .env rather than written in, because APP_PORT is
      // a thing people change: a proxy pinned to 8080 while the API listens
      // elsewhere refuses every call with ECONNREFUSED, and the only visible
      // symptom is a login that does nothing.
      '/api': { target: apiTarget, changeOrigin: true },
      '/studio': { target: apiTarget, changeOrigin: true },
      '/pulse': { target: apiTarget, changeOrigin: true },
      '/sentinel': { target: apiTarget, changeOrigin: true },
      '/docs': { target: apiTarget, changeOrigin: true },
    },
  },
  build: {
    // Into the Go tree, which is the only place //go:embed can read from: the
    // directive sits in api/main.go and cannot reach above its own directory.
    outDir: './api/web',
    emptyOutDir: true,
  },
})
`
}

func singleMakefile(opts Options) string {
	return fmt.Sprintf(`# %s: Single App Makefile

.PHONY: dev build migrate seed run clean

# Development: the Vite app at the root, the Go API in api/, in parallel.
dev:
	@echo "Starting development servers..."
	@pnpm dev &
	@cd api && air

# Build production binary. Vite writes into api/web and the binary embeds it,
# so the result is one file with the whole app inside it.
build:
	@echo "Building frontend..."
	@pnpm install && pnpm build
	@echo "Building Go binary..."
	@cd api && go build -o ../bin/%s .
	@echo "Done! Binary at bin/%s"

# One-shot migrate + seed (handy for local resets).
migrate:
	@cd api && go run ./cmd/migrate

seed:
	@cd api && go run ./cmd/seed

# Run the built binary.
run: build
	@./bin/%s

# Clean build artifacts (keeps the placeholder so go build still works).
clean:
	@rm -rf bin/
	@rm -rf api/web/assets
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
node_modules/
.vite/
.next/

# The built SPA, which lives in the Go tree so the binary can embed it.
# The placeholder index.html beside it is checked in on purpose: without it
# go build fails on a fresh clone.
api/web/assets/

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
