package main

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	gothGithub "github.com/markbates/goth/providers/github"
	"github.com/markbates/goth/providers/google"

	"commerce/apps/api/internal/ai"
	"commerce/apps/api/internal/cache"
	"commerce/apps/api/internal/config"
	"commerce/apps/api/internal/cron"
	"commerce/apps/api/internal/database"
	"commerce/apps/api/internal/events"
	"commerce/apps/api/internal/jobs"
	"commerce/apps/api/internal/mail"
	"commerce/apps/api/internal/models"
	"commerce/apps/api/internal/routes"
	"commerce/apps/api/internal/services"
	"commerce/apps/api/internal/storage"
	"commerce/apps/api/internal/tracing"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// OpenTelemetry, when OTEL_EXPORTER_OTLP_ENDPOINT points somewhere.
	//
	// Started before the database, so a slow connection is a span rather
	// than a gap. A failure here is logged and not fatal: a collector that
	// is down must not stop the API from serving, because observability is
	// how you find out something is wrong and it refusing to start is not
	// something being wrong.
	shutdownTracing, err := tracing.Setup(context.Background(), cfg)
	if err != nil {
		log.Printf("Tracing: %v", err)
	}
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := shutdownTracing(ctx); err != nil {
			log.Printf("Draining traces: %v", err)
		}
	}()
	if tracing.Enabled() {
		log.Printf("Tracing to %s", os.Getenv(tracing.EndpointEnv))
	}

	// Connect to database
	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Stage 6 of scaling, and it is one environment variable: set
	// DATABASE_REPLICA_URLS and reads go to the replicas from the next boot.
	// Registered here rather than inside Connect because migrate and seed use
	// Connect too, and a migration must never be routed to a replica.
	if _, err := database.UseReplicas(db); err != nil {
		log.Fatalf("Read replicas: %v", err)
	}

	// DB_PROVIDER=memory keeps everything in RAM, so the only process that can
	// usefully migrate it is this one: a separate migrate command would build the
	// schema in a
	// process that then exits, and this one would start on an empty database and
	// answer "no such table" to every request. Nothing is at risk either way,
	// because an in-memory database starts empty by definition.
	if strings.HasPrefix(cfg.DatabaseURL, "sqlite:file::memory:") {
		log.Println("DB_PROVIDER=memory: migrating in process, because the schema cannot outlive it")
		if err := models.Migrate(db); err != nil {
			log.Fatalf("Failed to migrate the in-memory database: %v", err)
		}
	}

	// ── Phase 4 Services ─────────────────────────────────────────

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

	// File storage. STORAGE_DRIVER=local keeps files in a directory on this
	// server and serves them from /files; minio, s3, r2 and b2 use a bucket.
	// STORAGE_DISKS adds named disks, such as a private bucket for backups,
	// which code reaches through storage.Disks.Get(name).
	storage.SetPublicPrefixes(cfg.StoragePublicPrefixes)
	var storageService *storage.Storage
	if s, err := storage.Open(cfg.StorageDriver, cfg.Storage, storage.LocalConfig{
		Root:      cfg.Storage.LocalRoot,
		PublicURL: cfg.Storage.PublicURL,
		Secret:    cfg.Storage.URLSecret,
	}); err != nil {
		log.Printf("Warning: Storage unavailable: %v (uploads disabled)", err)
	} else {
		storageService = s
		storage.Disks.SetDefault(s)
		log.Printf("File storage: %s", s.Describe())
	}
	for _, named := range cfg.StorageDisks {
		s, err := storage.Open(named.Driver, named.Storage, storage.LocalConfig{
			Root:      named.Storage.LocalRoot,
			PublicURL: named.Storage.PublicURL,
			Secret:    named.Storage.URLSecret,
		})
		if err != nil {
			log.Printf("Warning: storage disk %q unavailable: %v", named.Name, err)
			continue
		}
		storage.Disks.Add(named.Name, s)
		log.Printf("Storage disk %q: %s", named.Name, s.Describe())
	}

	// Email. MAIL_MAILER picks the transport: smtp, resend, mailgun, postmark,
	// sendgrid, ses, log or failover. Left empty, RESEND_API_KEY means Resend,
	// and development sends to Mailhog, falling back to the log.
	mailer, mailErr := mail.FromConfig(cfg)
	if mailErr != nil {
		log.Fatalf("Email: %v", mailErr)
	}
	if mailer != nil {
		log.Printf("Email: sending with %s", mailer.Driver())
	} else {
		log.Println("Warning: no mail driver configured (emails disabled). Set MAIL_MAILER, or RESEND_API_KEY for Resend.")
	}

	// AI service (Vercel AI Gateway)
	var aiService *ai.AI
	if cfg.AIGatewayAPIKey != "" {
		aiService = ai.New(cfg.AIGatewayAPIKey, cfg.AIGatewayModel, cfg.AIGatewayURL)
		log.Printf("AI service configured via AI Gateway (%s)", cfg.AIGatewayModel)
	}

	// Background jobs (asynq)
	//
	// Only when Redis actually answered. jobs.NewClient parses the URL and builds
	// a client without connecting to anything, so "Job queue connected" used to
	// print on a machine with no Redis at all, two lines after the warning saying
	// Redis was unavailable.
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
	//
	// NOTE: these callback URLs are deliberately NOT versioned, even though the
	// routes now live under /api/v1. The same string is registered in the
	// Google / GitHub console as an authorized redirect URI — a value you
	// control there, not here. Adding "/v1" would stop matching what every
	// existing deployment has registered and break social login on upgrade,
	// which is the exact class of breakage the version prefix exists to avoid.
	// The unversioned path is re-dispatched to the current version by
	// mountLegacyAPIAlias (query string preserved), so these keep working.
	// Gothic keeps the OAuth handshake in a cookie of its own. Its key is derived
	// from the JWT secret rather than being it, so the cookie and the tokens never
	// share a key, and the cookie is HttpOnly, SameSite=Lax, and Secure when the
	// API is served over https.
	gothicKey := sha256.Sum256([]byte("gothic-cookie-store:" + cfg.JWTSecret))
	gothicStore := sessions.NewCookieStore(gothicKey[:])
	gothicStore.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   600,
		HttpOnly: true,
		Secure:   strings.HasPrefix(cfg.AppURL, "https://"),
		SameSite: http.SameSiteLaxMode,
	}
	gothic.Store = gothicStore
	var oauthProviders []goth.Provider
	if cfg.GoogleClientID != "" {
		oauthProviders = append(oauthProviders, google.New(
			cfg.GoogleClientID, cfg.GoogleClientSecret,
			cfg.AppURL+"/api/auth/oauth/google/callback",
		))
		log.Println("Google OAuth2 configured")
	}
	if cfg.GithubClientID != "" {
		oauthProviders = append(oauthProviders, gothGithub.New(
			cfg.GithubClientID, cfg.GithubClientSecret,
			cfg.AppURL+"/api/auth/oauth/github/callback",
		))
		log.Println("GitHub OAuth2 configured")
	}
	if len(oauthProviders) > 0 {
		goth.UseProviders(oauthProviders...)
	}

	// Build services
	var secObsBridge *services.SecObsBridge
	if cfg.SentinelEnabled || cfg.PulseEnabled {
		secObsBridge = services.NewSecObsBridge(cfg)
	}

	svc := &routes.Services{
		Cache:   cacheService,
		Storage: storageService,
		Mailer:  mailer,
		AI:      aiService,
		Jobs:    jobClient,
		SecObs:  secObsBridge,
	}

	// Setup router
	router := routes.Setup(db, cfg, svc)

	// Start the SecObs notification poller (turns Sentinel/Pulse findings
	// into in-app notifications). Runs once a minute on its own goroutine;
	// no-op when the bridge is nil.
	var secObsPoller *services.SecObsPoller
	if secObsBridge != nil {
		secObsPoller = services.NewSecObsPoller(db, secObsBridge)
		secObsPoller.Start()
	}

	// Start background worker
	//
	// Only when Redis answered: the worker polls in a loop, so without Redis it
	// writes an asynq error every second or two, forever. That noise was the worst
	// part of starting a project with no Redis running, and it drowned out the one
	// line that explained it.
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

	// Start cron scheduler, for the same reason and on the same condition.
	var cronScheduler *cron.Scheduler
	if cfg.RedisURL != "" && redisReachable {
		cs, err := cron.New(cfg.RedisURL)
		if err != nil {
			log.Printf("Warning: Cron scheduler failed to start: %v", err)
		} else {
			cronScheduler = cs
			if err := cs.Start(); err != nil {
				log.Printf("Warning: Cron scheduler failed to start: %v", err)
			} else {
				log.Println("Cron scheduler ready: it runs on the replica holding the cron lock")
			}
		}
	}

	// Reap ImportJobs orphaned by a crash or restart. A background CSV import
	// runs in a goroutine that flips the job to completed/failed at the end;
	// if the process dies first, the row is stuck "processing" forever and the
	// client's poll never terminates. This needs no Redis, so it always runs.
	go func() {
		// At boot, ANY processing job is orphaned — its goroutine died with
		// the previous process — so reap them all immediately.
		db.Model(&models.ImportJob{}).Where("status = ?", "processing").
			Updates(map[string]interface{}{
				"status":  "failed",
				"message": "import interrupted by server restart",
			})
		// Thereafter, reap only jobs with no progress for 15 minutes (a stall).
		ticker := time.NewTicker(10 * time.Minute)
		defer ticker.Stop()
		for range ticker.C {
			cutoff := time.Now().Add(-15 * time.Minute)
			db.Model(&models.ImportJob{}).
				Where("status = ? AND updated_at < ?", "processing", cutoff).
				Updates(map[string]interface{}{
					"status":  "failed",
					"message": "import stalled (no progress for 15 minutes)",
				})
		}
	}()

	// Create server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.Port),
		Handler: router,

		// Short, because they are what stops a slow client holding a connection
		// open. Uploads, exports and streams get longer deadlines per route, from
		// middleware.RequestLimits in routes.go.
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      60 * time.Second,
		IdleTimeout:       120 * time.Second,
	}

	// Start server in goroutine
	go func() {
		log.Printf("Server starting on port %s", cfg.Port)
		if cfg.GORMStudioEnabled {
			log.Printf("GORM Studio available at http://localhost:%s/studio", cfg.Port)
		}
		if cfg.AppEnv != "production" || cfg.APIDocsPublic {
			log.Printf("API Documentation at http://localhost:%s/docs", cfg.Port)
		}

		if cfg.PulseEnabled {
			log.Printf("Pulse dashboard at http://localhost:%s/pulse/ui/", cfg.Port)
		}
		if cfg.SentinelEnabled {
			log.Printf("Sentinel dashboard at http://localhost:%s/sentinel/ui", cfg.Port)
		}
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	if secObsPoller != nil {
		secObsPoller.Stop()
	}

	// Stop cron scheduler
	if cronScheduler != nil {
		cronScheduler.Stop()
	}

	// Stop background worker
	if workerStop != nil {
		workerStop()
	}

	// Close job client
	if jobClient != nil {
		jobClient.Close()
	}

	// Close cache connection
	if cacheService != nil {
		cacheService.Close()
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	// After the last request: stop claiming outbox messages this process will
	// not live to deliver, and write the activity rows still queued.
	events.StopRelay()
	services.FlushActivity(ctx)

	log.Println("Server exited")
}
