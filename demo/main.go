package main

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

	"gritdemo/internal/cache"
	"gritdemo/internal/config"
	kcron "gritdemo/internal/cron"
	"gritdemo/internal/database"
	"gritdemo/internal/mail"
	"gritdemo/internal/models"
	"gritdemo/internal/routes"
	gritservices "gritdemo/internal/services"
	"gritdemo/internal/storage"
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

	// Auto-migrate all models
	log.Println("Running auto-migrations...")
	if err := models.Migrate(db); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}
	log.Println("Migrations complete")

	// Canonical seed — admin user + Grit Motors business + Main Branch.
	// Idempotent; runs on every boot.
	if err := database.Seed(db); err != nil {
		log.Printf("Warning: seeding failed: %v", err)
	}
	// Demo cohort — staff, categories, products, stock, motorcycles, loan
	// products, borrowers, loans + schedules + repayments, cash sales,
	// POS history, daily-boda. Each step idempotent. Same code path the
	// nightly demo-reset cron uses after wiping mutable rows.
	if cfg.DemoMode {
		if err := database.SeedDemo(db); err != nil {
			log.Printf("Warning: demo seeding failed: %v", err)
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
	// In DEMO_MODE we never wire the real mailer — keeps the live demo from
	// blasting low-stock alerts and invitation emails into anyone's inbox.
	// Handlers that hold a nil *mail.Mailer treat it as "email disabled".
	if cfg.DemoMode {
		log.Println("Email service disabled (DEMO_MODE=true)")
	} else if cfg.ResendAPIKey != "" && cfg.ResendAPIKey != "re_your_api_key" {
		mailer = mail.New(cfg.ResendAPIKey, cfg.MailFrom)
		log.Println("Email service configured")
	}

	// Build services
	var secObsBridge *gritservices.SecObsBridge
	if cfg.SentinelEnabled || cfg.PulseEnabled {
		secObsBridge = gritservices.NewSecObsBridge(cfg)
	}

	svc := &routes.Services{
		Cache:   cacheService,
		Storage: storageService,
		Mailer:  mailer,
		SecObs:  secObsBridge,
	}

	// Background scheduler — nightly overdue flagging and (in demo mode)
	// the demo-reset job that wipes mutable rows + reseeds every 24h.
	scheduler := kcron.New(db).WithDemoMode(cfg.DemoMode)
	if err := scheduler.Register(); err != nil {
		log.Printf("Warning: cron registration failed: %v", err)
	} else {
		scheduler.Start()
		log.Println("Cron scheduler started")
	}
	defer scheduler.Stop()

	// Sentinel + Pulse → in-app notification feed. Runs once a minute on
	// its own goroutine; no-op when the bridge is nil.
	var secObsPoller *gritservices.SecObsPoller
	if secObsBridge != nil {
		secObsPoller = gritservices.NewSecObsPoller(db, secObsBridge)
		secObsPoller.Start()
		log.Println("SecObs poller started")
	}
	defer func() {
		if secObsPoller != nil {
			secObsPoller.Stop()
		}
	}()

	// Setup router
	router := routes.Setup(db, cfg, svc)

	// Serve embedded frontend (SPA fallback)
	feFS, err := fs.Sub(frontendFS, "frontend/dist")
	if err != nil {
		log.Printf("Warning: embedded frontend not available: %v", err)
	} else {
		// Pre-read index.html for SPA fallback (avoids http.FileServer redirect loop)
		indexHTML, err := fs.ReadFile(feFS, "index.html")
		if err != nil {
			log.Printf("Warning: index.html not found in embedded frontend: %v", err)
		}
		fileServer := http.FileServer(http.FS(feFS))
		router.NoRoute(func(c *gin.Context) {
			p := c.Request.URL.Path
			// Try to serve static file (JS, CSS, images, etc.)
			if p != "/" && p != "/index.html" {
				if f, err := feFS.Open(p[1:]); err == nil {
					f.Close()
					fileServer.ServeHTTP(c.Writer, c.Request)
					return
				}
			}
			// SPA fallback: serve index.html directly from memory
			c.Data(http.StatusOK, "text/html; charset=utf-8", indexHTML)
		})
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

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Server shutdown error: %v", err)
	}
	log.Println("Server stopped")
}
