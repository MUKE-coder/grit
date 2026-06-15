package scaffold

import (
	"fmt"
	"path/filepath"
	"strings"
)

func writeAPIFiles(root string, opts Options) error {
	apiRoot := opts.APIRoot(root)
	module := opts.Module()

	files := map[string]string{
		filepath.Join(apiRoot, "go.mod"):                              apiGoMod(opts),
		filepath.Join(apiRoot, ".gitignore"):                          apiGitignore(),
		filepath.Join(apiRoot, "cmd", "server", "main.go"):           apiMainGo(opts),
		filepath.Join(apiRoot, "internal", "config", "config.go"):    apiConfigGo(),
		filepath.Join(apiRoot, "internal", "database", "database.go"): apiDatabaseGo(),
		filepath.Join(apiRoot, "internal", "models", "user.go"):               apiUserModelGo(),
		filepath.Join(apiRoot, "internal", "models", "upload.go"):             apiUploadModelGo(),
		filepath.Join(apiRoot, "internal", "models", "ui_component.go"):       apiUIComponentModelGo(),
		filepath.Join(apiRoot, "internal", "models", "seed_ui_components.go"):          apiUIComponentSeedGo(),
		filepath.Join(apiRoot, "internal", "models", "seed_ui_components_extended.go"): apiUIComponentSeedExtendedGo(),
		filepath.Join(apiRoot, "internal", "handlers", "ui_registry.go"):      apiUIRegistryHandlerGo(),
		filepath.Join(apiRoot, "internal", "services", "auth.go"):    apiAuthServiceGo(),
		filepath.Join(apiRoot, "internal", "handlers", "auth.go"):    apiAuthHandlerGo(),
		filepath.Join(apiRoot, "internal", "handlers", "user.go"):    apiUserHandlerGo(),
		filepath.Join(apiRoot, "internal", "middleware", "auth.go"):        apiAuthMiddlewareGo(),
		filepath.Join(apiRoot, "internal", "middleware", "cors.go"):        apiCorsMiddlewareGo(),
		filepath.Join(apiRoot, "internal", "middleware", "logger.go"):      apiLoggerMiddlewareGo(),
		filepath.Join(apiRoot, "internal", "middleware", "maintenance.go"): apiMaintenanceMiddlewareGo(),
		filepath.Join(apiRoot, "internal", "middleware", "idempotency.go"): apiIdempotencyMiddlewareGo(),
		filepath.Join(apiRoot, "internal", "paginate", "paginate.go"): apiPaginateGo(),
		filepath.Join(apiRoot, "internal", "realtime", "hub.go"):       apiRealtimeHubGo(),
		filepath.Join(apiRoot, "internal", "handlers", "realtime.go"):  apiRealtimeHandlerGo(),
		filepath.Join(apiRoot, "internal", "sync", "registry.go"):      apiSyncRegistryGo(),
		filepath.Join(apiRoot, "internal", "handlers", "sync.go"):      apiSyncHandlerGo(),
		filepath.Join(apiRoot, "internal", "models", "activity_log.go"):       apiActivityLogModelGo(),
		filepath.Join(apiRoot, "internal", "middleware", "activity.go"):       apiActivityMiddlewareGo(),
		filepath.Join(apiRoot, "internal", "handlers", "activity.go"):         apiActivityHandlerGo(),
		filepath.Join(apiRoot, "internal", "export", "export.go"):             apiExportGo(),
		filepath.Join(apiRoot, "internal", "respond", "respond.go"):           apiRespondGo(),
		filepath.Join(apiRoot, "internal", "pdf", "pdf.go"):                   apiPDFGo(),
		filepath.Join(apiRoot, "internal", "pdf", "invoice.go"):               apiPDFInvoiceGo(),
		filepath.Join(apiRoot, "internal", "audit", "audit.go"):               apiAuditGo(),
		filepath.Join(apiRoot, "internal", "models", "webhook_event.go"):      apiWebhookEventModelGo(),
		filepath.Join(apiRoot, "internal", "webhooks", "webhooks.go"):         apiWebhooksGo(),
		filepath.Join(apiRoot, "internal", "webhooks", "verifiers.go"):        apiWebhooksVerifiersGo(),
		filepath.Join(apiRoot, "internal", "handlers", "webhooks.go"):         apiWebhooksHandlerGo(),
		filepath.Join(apiRoot, "internal", "models", "feature_flag.go"):       apiFeatureFlagModelGo(),
		filepath.Join(apiRoot, "internal", "flags", "flags.go"):               apiFlagsGo(),
		filepath.Join(apiRoot, "internal", "handlers", "flags.go"):            apiFlagsHandlerGo(),
		filepath.Join(apiRoot, "internal", "routes", "routes.go"):    apiRoutesGo(),
		filepath.Join(apiRoot, ".air.toml"):                          airConfig(),
		// Test files — give the generated API a working test suite out of the box
		filepath.Join(apiRoot, "internal", "handlers", "auth_test.go"):  apiAuthTestGo(),
		filepath.Join(apiRoot, "internal", "handlers", "user_test.go"):  apiUserTestGo(),
		filepath.Join(apiRoot, "internal", "handlers", "bench_test.go"): apiBenchTestGo(),
	}

	for path, content := range files {
		// Replace module placeholder
		content = strings.ReplaceAll(content, "{{MODULE}}", module)
		if err := writeFile(path, content); err != nil {
			return fmt.Errorf("writing %s: %w", path, err)
		}
	}

	return nil
}

func apiGoMod(opts Options) string {
	return fmt.Sprintf(`module %s

go 1.21

require (
	github.com/MUKE-coder/gin-docs v0.0.0-20260222113017-4d647cb4e7aa
	github.com/MUKE-coder/gorm-studio v1.0.1
	// Pinned to the v1.0.0 commit on main.
	github.com/MUKE-coder/pulse v0.0.0-20260529025319-478cdfa8ce5f
	github.com/aws/aws-sdk-go-v2 v1.25.0
	github.com/aws/aws-sdk-go-v2/config v1.27.0
	github.com/aws/aws-sdk-go-v2/credentials v1.17.0
	github.com/aws/aws-sdk-go-v2/feature/s3/manager v1.16.0
	github.com/aws/aws-sdk-go-v2/service/s3 v1.51.0
	github.com/disintegration/imaging v1.6.2
	github.com/gin-gonic/gin v1.10.0
	github.com/go-pdf/fpdf v0.9.0
	github.com/golang-jwt/jwt/v5 v5.2.0
	github.com/google/uuid v1.6.0
	github.com/gorilla/sessions v1.4.0
	github.com/gorilla/websocket v1.5.3
	github.com/hibiken/asynq v0.24.1
	github.com/markbates/goth v1.80.0
	github.com/joho/godotenv v1.5.1
	github.com/redis/go-redis/v9 v9.4.0
	github.com/xuri/excelize/v2 v2.8.1
	golang.org/x/crypto v0.23.0
	// Pinned to the v2.0.1 commit on main. Upstream is tagged v2.0.1
	// but the module path lacks the /v2 suffix Go modules requires for
	// major versions >= 2, so the tag fails go.sum verification. A
	// pseudo-version side-steps the rule; the underlying code is the
	// same as the tag. Advance the pin once sentinel re-tags with /v2.
	github.com/MUKE-coder/sentinel v0.0.0-20260529033414-0e945440db7f
	gorm.io/datatypes v1.2.7
	gorm.io/driver/postgres v1.5.11
	gorm.io/gorm v1.25.12
)

require (
	github.com/stretchr/testify v1.9.0
	github.com/glebarez/sqlite v1.11.0
)
`, opts.Module())
}

func apiGitignore() string {
	return `# Binary
*.exe
*.exe~
*.dll
*.so
*.dylib
tmp/

# Environment
.env
.env.local

# IDE
.vscode/
.idea/
`
}

func airConfig() string {
	return `root = "."
tmp_dir = "tmp"

[build]
  bin = "./tmp/server"
  cmd = "go build -o ./tmp/server ./cmd/server"
  delay = 1000
  exclude_dir = ["tmp", "vendor", "node_modules"]
  exclude_regex = ["_test.go"]
  include_ext = ["go", "toml", "yaml"]
  kill_delay = "0s"
  send_interrupt = false
  stop_on_error = true

[log]
  time = false

[color]
  build = "yellow"
  main = "magenta"
  runner = "green"
  watcher = "cyan"
`
}

func apiMainGo(opts Options) string {
	return `package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	gothGithub "github.com/markbates/goth/providers/github"
	"github.com/markbates/goth/providers/google"

	"` + "{{MODULE}}" + `/internal/ai"
	"` + "{{MODULE}}" + `/internal/cache"
	"` + "{{MODULE}}" + `/internal/config"
	"` + "{{MODULE}}" + `/internal/cron"
	"` + "{{MODULE}}" + `/internal/database"
	"` + "{{MODULE}}" + `/internal/jobs"
	"` + "{{MODULE}}" + `/internal/mail"
	"` + "{{MODULE}}" + `/internal/routes"
	"` + "{{MODULE}}" + `/internal/services"
	"` + "{{MODULE}}" + `/internal/storage"
)

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

	// ── Phase 4 Services ─────────────────────────────────────────

	// Redis cache
	var cacheService *cache.Cache
	if cfg.RedisURL != "" {
		c, err := cache.New(cfg.RedisURL)
		if err != nil {
			log.Printf("Warning: Redis unavailable: %v (caching disabled)", err)
		} else {
			cacheService = c
			log.Println("Redis cache connected")
		}
	}

	// File storage (S3-compatible)
	var storageService *storage.Storage
	if cfg.Storage.Endpoint != "" && cfg.Storage.AccessKey != "" {
		s, err := storage.New(cfg.Storage)
		if err != nil {
			log.Printf("Warning: Storage unavailable: %v (uploads disabled)", err)
		} else {
			storageService = s
			log.Println("File storage connected")
		}
	}

	// Email (Resend)
	var mailer *mail.Mailer
	if cfg.ResendAPIKey != "" && cfg.ResendAPIKey != "re_your_api_key" {
		mailer = mail.New(cfg.ResendAPIKey, cfg.MailFrom)
		log.Println("Email service configured")
	} else {
		log.Println("Warning: Resend API key not set (emails disabled)")
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
	var cronScheduler *cron.Scheduler
	if cfg.RedisURL != "" {
		cs, err := cron.New(cfg.RedisURL)
		if err != nil {
			log.Printf("Warning: Cron scheduler failed to start: %v", err)
		} else {
			cronScheduler = cs
			if err := cs.Start(); err != nil {
				log.Printf("Warning: Cron scheduler failed to start: %v", err)
			} else {
				log.Println("Cron scheduler started")
			}
		}
	}

	// Create server
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.Port),
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in goroutine
	go func() {
		log.Printf("Server starting on port %s", cfg.Port)
		log.Printf("GORM Studio available at http://localhost:%s/studio", cfg.Port)
		log.Printf("API Documentation at http://localhost:%s/docs", cfg.Port)
		if cfg.PulseEnabled {
			log.Printf("Pulse dashboard at http://localhost:%s/pulse/ui/", cfg.Port)
		}
		if cfg.SentinelEnabled {
			log.Printf("Sentinel dashboard at http://localhost:%s/sentinel/ui", cfg.Port)
		}
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
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

	log.Println("Server exited")
}
`
}

func apiConfigGo() string {
	return `package config

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

// StorageConfig holds credentials for a single S3-compatible provider.
type StorageConfig struct {
	Endpoint  string
	AccessKey string
	SecretKey string
	Bucket    string
	Region    string
	UseSSL    bool
}

// Config holds all application configuration.
type Config struct {
	AppName     string
	AppEnv      string
	Port        string
	AppURL      string
	DatabaseURL string

	JWTSecret        string
	JWTAccessExpiry  time.Duration
	JWTRefreshExpiry time.Duration

	RedisURL string

	// Storage
	StorageDriver string        // "minio", "s3", "r2", or "b2"
	Storage       StorageConfig // Resolved config for the active driver

	ResendAPIKey string
	MailFrom     string

	CORSOrigins []string

	GORMStudioEnabled  bool
	GORMStudioUsername string
	GORMStudioPassword string

	// AI (Vercel AI Gateway)
	AIGatewayAPIKey string
	AIGatewayModel  string
	AIGatewayURL    string

	// TOTP (Two-Factor Authentication)
	TOTPIssuer string

	// Security (Sentinel)
	SentinelEnabled        bool
	SentinelUsername       string
	SentinelPassword       string
	SentinelSecretKey      string
	// Sentinel v2.0 — CIDRs allowed to send X-Forwarded-For / X-Real-IP.
	// Empty (default) means "ignore those headers entirely" — safe when
	// the app speaks to the public internet directly; populate when
	// you're behind a known reverse proxy (Caddy/Traefik/Cloudflare).
	SentinelTrustedProxies []string

	// Observability (Pulse v1.0)
	PulseEnabled    bool
	PulseUsername    string
	PulsePassword   string
	// Pulse v1.0 storage. Defaults to in-memory ring buffer (no disk).
	// Set PULSE_STORAGE=sqlite + PULSE_STORAGE_DSN=pulse.db to enable
	// the new persistent backend (WAL, busy_timeout=5s, survives restart).
	PulseStorage    string // "memory" (default) | "sqlite"
	PulseStorageDSN string // path for sqlite, e.g. "pulse.db" or ":memory:"

	// OAuth2 Social Login
	GoogleClientID     string
	GoogleClientSecret string
	GithubClientID     string
	GithubClientSecret string
	OAuthFrontendURL   string // Where to redirect after OAuth callback
}

// Load reads configuration from environment variables.
func Load() (*Config, error) {
	// Load .env file (ignore error if not found — production uses real env vars)
	_ = godotenv.Load()
	_ = godotenv.Load("../../.env") // Load from project root when running from apps/api

	storageDriver := getEnv("STORAGE_DRIVER", "minio")

	cfg := &Config{
		AppName:     getEnv("APP_NAME", "grit-app"),
		AppEnv:      getEnv("APP_ENV", "development"),
		Port:        getEnv("APP_PORT", "8080"),
		AppURL:      getEnv("APP_URL", "http://localhost:8080"),
		DatabaseURL: resolveDatabaseURL(),
		JWTSecret:   getEnv("JWT_SECRET", ""),
		RedisURL:    getEnv("REDIS_URL", "redis://localhost:6380"),

		StorageDriver: storageDriver,
		Storage:       resolveStorage(storageDriver),

		ResendAPIKey: getEnv("RESEND_API_KEY", ""),
		MailFrom:     getEnv("MAIL_FROM", "noreply@localhost"),

		CORSOrigins: strings.Split(getEnv("CORS_ORIGINS", "http://localhost:3000,http://localhost:3001"), ","),

		GORMStudioEnabled:  getEnv("GORM_STUDIO_ENABLED", "true") == "true",
		GORMStudioUsername: getEnv("GORM_STUDIO_USERNAME", "admin"),
		GORMStudioPassword: getEnv("GORM_STUDIO_PASSWORD", "studio"),

		AIGatewayAPIKey: getEnv("AI_GATEWAY_API_KEY", ""),
		AIGatewayModel:  getEnv("AI_GATEWAY_MODEL", "anthropic/claude-sonnet-4-6"),
		AIGatewayURL:    getEnv("AI_GATEWAY_URL", "https://ai-gateway.vercel.sh/v1"),

		TOTPIssuer: getEnv("TOTP_ISSUER", getEnv("APP_NAME", "grit-app")),

		SentinelEnabled:        getEnv("SENTINEL_ENABLED", "true") == "true",
		SentinelUsername:       getEnv("SENTINEL_USERNAME", "admin"),
		SentinelPassword:       getEnv("SENTINEL_PASSWORD", "sentinel"),
		SentinelSecretKey:      getEnv("SENTINEL_SECRET_KEY", "sentinel-secret-change-me"),
		SentinelTrustedProxies: splitCSV(getEnv("SENTINEL_TRUSTED_PROXIES", "")),

		PulseEnabled:    getEnv("PULSE_ENABLED", "true") == "true",
		PulseUsername:    getEnv("PULSE_USERNAME", "admin"),
		PulsePassword:   getEnv("PULSE_PASSWORD", "pulse"),
		PulseStorage:    getEnv("PULSE_STORAGE", "memory"),
		PulseStorageDSN: getEnv("PULSE_STORAGE_DSN", "pulse.db"),

		GoogleClientID:     getEnv("GOOGLE_CLIENT_ID", ""),
		GoogleClientSecret: getEnv("GOOGLE_CLIENT_SECRET", ""),
		GithubClientID:     getEnv("GITHUB_CLIENT_ID", ""),
		GithubClientSecret: getEnv("GITHUB_CLIENT_SECRET", ""),
		OAuthFrontendURL:   getEnv("OAUTH_FRONTEND_URL", "http://localhost:3001"),
	}

	// DatabaseURL is always populated by resolveDatabaseURL() — either from
	// the DATABASE_URL env var or built from POSTGRES_* parts. The actual
	// connection attempt in cmd/server/main.go will surface a useful error
	// if the resolved URL points at an unreachable database.

	if cfg.JWTSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET is required")
	}
	if len(cfg.JWTSecret) < 32 {
		log.Println("WARNING: JWT_SECRET should be at least 32 characters for security. Generate one with: openssl rand -hex 32")
	}

	// Parse durations
	accessExpiry, err := time.ParseDuration(getEnv("JWT_ACCESS_EXPIRY", "15m"))
	if err != nil {
		return nil, fmt.Errorf("invalid JWT_ACCESS_EXPIRY: %w", err)
	}
	cfg.JWTAccessExpiry = accessExpiry

	refreshExpiry, err := time.ParseDuration(getEnv("JWT_REFRESH_EXPIRY", "168h"))
	if err != nil {
		return nil, fmt.Errorf("invalid JWT_REFRESH_EXPIRY: %w", err)
	}
	cfg.JWTRefreshExpiry = refreshExpiry

	return cfg, nil
}

// IsDevelopment returns true if the app is running in development mode.
func (c *Config) IsDevelopment() bool {
	return c.AppEnv == "development"
}

// resolveDatabaseURL returns the connection string for the database.
//
// Single source of truth: edit POSTGRES_USER / POSTGRES_PASSWORD /
// POSTGRES_DB / POSTGRES_HOST / POSTGRES_PORT in .env and both
// docker-compose.yml and this function read the SAME values, so they
// can't drift.
//
// Resolution order:
//
//  1. If DATABASE_URL is set, use it verbatim — that's the escape hatch
//     for external Postgres (Neon, Supabase, RDS) or SQLite. It wins over
//     the POSTGRES_* parts so a one-line override is enough to swap.
//  2. Otherwise build postgres://USER:PASS@HOST:PORT/DB?sslmode=disable
//     from the parts above. Defaults match docker-compose.yml's
//     ${VAR:-grit} fallbacks so a fresh project boots even before the
//     user touches .env.
func resolveDatabaseURL() string {
	if v := os.Getenv("DATABASE_URL"); v != "" {
		return v
	}
	user := getEnv("POSTGRES_USER", "grit")
	pass := getEnv("POSTGRES_PASSWORD", "grit")
	host := getEnv("POSTGRES_HOST", "localhost")
	port := getEnv("POSTGRES_PORT", "5432")
	db := getEnv("POSTGRES_DB", getEnv("APP_NAME", "grit-app"))
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		user, pass, host, port, db)
}

// resolveStorage returns the StorageConfig for the active driver.
//
// For AWS S3, leave S3_ENDPOINT empty — the AWS SDK will use the
// regional endpoint automatically (s3.<region>.amazonaws.com).
// Credentials fall back to the AWS standard env vars
// AWS_ACCESS_KEY_ID + AWS_SECRET_ACCESS_KEY if you don't set the S3_*
// variants, which is convenient when running on EC2 / ECS / Lambda
// with an IAM role and you'd rather not duplicate keys in .env.
func resolveStorage(driver string) StorageConfig {
	switch driver {
	case "s3":
		// Empty endpoint = AWS SDK uses the regional default
		// (s3.<region>.amazonaws.com). This also flips the client into
		// virtual-hosted style, which AWS requires for buckets created
		// after Sep 2020.
		return StorageConfig{
			Endpoint:  getEnv("S3_ENDPOINT", ""),
			AccessKey: firstNonEmpty(os.Getenv("S3_ACCESS_KEY"), os.Getenv("AWS_ACCESS_KEY_ID")),
			SecretKey: firstNonEmpty(os.Getenv("S3_SECRET_KEY"), os.Getenv("AWS_SECRET_ACCESS_KEY")),
			Bucket:    getEnv("S3_BUCKET", "uploads"),
			Region:    firstNonEmpty(os.Getenv("S3_REGION"), os.Getenv("AWS_REGION"), "us-east-1"),
			UseSSL:    true,
		}
	case "r2":
		return StorageConfig{
			Endpoint:  getEnv("R2_ENDPOINT", ""),
			AccessKey: getEnv("R2_ACCESS_KEY", ""),
			SecretKey: getEnv("R2_SECRET_KEY", ""),
			Bucket:    getEnv("R2_BUCKET", "uploads"),
			Region:    getEnv("R2_REGION", "auto"),
			UseSSL:    true,
		}
	case "b2":
		return StorageConfig{
			Endpoint:  getEnv("B2_ENDPOINT", ""),
			AccessKey: getEnv("B2_ACCESS_KEY", ""),
			SecretKey: getEnv("B2_SECRET_KEY", ""),
			Bucket:    getEnv("B2_BUCKET", "uploads"),
			Region:    getEnv("B2_REGION", "us-west-004"),
			UseSSL:    true,
		}
	default: // minio
		return StorageConfig{
			Endpoint:  getEnv("MINIO_ENDPOINT", "http://localhost:9002"),
			AccessKey: getEnv("MINIO_ACCESS_KEY", "minioadmin"),
			SecretKey: getEnv("MINIO_SECRET_KEY", "minioadmin"),
			Bucket:    getEnv("MINIO_BUCKET", "uploads"),
			Region:    getEnv("MINIO_REGION", "us-east-1"),
			UseSSL:    getEnv("MINIO_USE_SSL", "false") == "true",
		}
	}
}

// firstNonEmpty returns the first non-empty string in vals, or "" if all
// are empty. Useful for letting S3_* override AWS_* with a graceful
// fallback.
func firstNonEmpty(vals ...string) string {
	for _, v := range vals {
		if v != "" {
			return v
		}
	}
	return ""
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

// splitCSV trims and splits a comma-separated env var. Empty strings
// after trimming are dropped so "a, ,b" yields ["a","b"].
func splitCSV(s string) []string {
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if p = strings.TrimSpace(p); p != "" {
			out = append(out, p)
		}
	}
	return out
}
`
}

func apiDatabaseGo() string {
	return `package database

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/glebarez/sqlite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Connect establishes a database connection using the provided DSN.
//
// Driver is chosen by DSN shape:
//   - "sqlite://path" or "sqlite:path"  → SQLite (file or :memory:)
//   - anything else                     → Postgres
//
// Examples:
//   DATABASE_URL=sqlite:./bench.db
//   DATABASE_URL=sqlite::memory:
//   DATABASE_URL=postgres://user:pass@host:5432/db?sslmode=disable
func Connect(dsn string) (*gorm.DB, error) {
	logLevel := logger.Warn
	if os.Getenv("DB_LOG_LEVEL") == "info" {
		logLevel = logger.Info
	} else if os.Getenv("DB_LOG_LEVEL") == "silent" {
		logLevel = logger.Silent
	}
	gormCfg := &gorm.Config{Logger: logger.Default.LogMode(logLevel)}

	var (
		db  *gorm.DB
		err error
	)

	switch {
	case strings.HasPrefix(dsn, "sqlite://"):
		db, err = gorm.Open(sqlite.Open(strings.TrimPrefix(dsn, "sqlite://")), gormCfg)
	case strings.HasPrefix(dsn, "sqlite:"):
		db, err = gorm.Open(sqlite.Open(strings.TrimPrefix(dsn, "sqlite:")), gormCfg)
	default:
		db, err = gorm.Open(postgres.New(postgres.Config{
			DSN:                  dsn,
			PreferSimpleProtocol: true,
		}), gormCfg)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	// Connection pool settings. SQLite ignores most of these — single-writer
	// semantics mean MaxOpenConns above 1 only helps concurrent reads, and
	// SQLite serialises writes internally. Postgres uses every knob.
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)
	sqlDB.SetConnMaxIdleTime(10 * time.Minute)

	log.Println("Database connected successfully")
	return db, nil
}
`
}

func apiUserModelGo() string {
	return `package models

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Role constants
const (
	RoleAdmin  = "ADMIN"
	RoleEditor = "EDITOR"
	RoleUser   = "USER"
	// grit:roles
)

// User represents a user in the system.
type User struct {
	ID              string         ` + "`" + `gorm:"primarykey;size:36" json:"id"` + "`" + `
	FirstName       string         ` + "`" + `gorm:"size:255;not null" json:"first_name" binding:"required"` + "`" + `
	LastName        string         ` + "`" + `gorm:"size:255;not null" json:"last_name" binding:"required"` + "`" + `
	Email           string         ` + "`" + `gorm:"size:255;uniqueIndex;not null" json:"email" binding:"required,email"` + "`" + `
	Password        string         ` + "`" + `gorm:"size:255" json:"-"` + "`" + `
	Role            string         ` + "`" + `gorm:"size:20;default:USER" json:"role"` + "`" + `
	Avatar          string         ` + "`" + `gorm:"size:500" json:"avatar"` + "`" + `
	JobTitle        string         ` + "`" + `gorm:"size:255" json:"job_title"` + "`" + `
	Bio             string         ` + "`" + `gorm:"type:text" json:"bio"` + "`" + `
	Active          bool           ` + "`" + `gorm:"default:true" json:"active"` + "`" + `
	Provider        string         ` + "`" + `gorm:"size:50;default:'local'" json:"provider"` + "`" + `
	GoogleID        string         ` + "`" + `gorm:"size:255" json:"-"` + "`" + `
	GithubID        string         ` + "`" + `gorm:"size:255" json:"-"` + "`" + `
	EmailVerifiedAt *time.Time     ` + "`" + `json:"email_verified_at"` + "`" + `
	IPAddress       string         ` + "`" + `gorm:"size:45" json:"ip_address"` + "`" + `
	MACAddress      string         ` + "`" + `gorm:"size:50" json:"mac_address"` + "`" + `
	Version         int            ` + "`" + `gorm:"not null;default:1" json:"version"` + "`" + `
	CreatedAt       time.Time      ` + "`" + `json:"created_at"` + "`" + `
	UpdatedAt       time.Time      ` + "`" + `json:"updated_at"` + "`" + `
	DeletedAt       gorm.DeletedAt ` + "`" + `gorm:"index" json:"-"` + "`" + `
}

// BeforeCreate generates a UUID and hashes the password before saving.
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == "" {
		u.ID = uuid.New().String()
	}
	if u.Version == 0 {
		u.Version = 1
	}
	if u.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		u.Password = string(hashedPassword)
	}
	return nil
}

// BeforeUpdate increments Version so offline clients can detect that
// a record they edited has moved on. Pair with the Idempotency-Key
// middleware + /api/sync/push for safe write replay.
func (u *User) BeforeUpdate(tx *gorm.DB) error {
	u.Version++
	return nil
}

// BeforeCreate generates a UUID for uploads.
func (u *Upload) BeforeCreate(tx *gorm.DB) error {
	if u.ID == "" {
		u.ID = uuid.New().String()
	}
	return nil
}

// CheckPassword compares the given password with the stored hash.
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

// Models returns the ordered list of all models for migration.
// Models with no foreign key dependencies come first.
func Models() []interface{} {
	return []interface{}{
		&User{},
		&Upload{},
		&Blog{},
		&UIComponent{},
		&TwoFactorConfig{},
		&TrustedDevice{},
		&TOTPPendingToken{},
		&ActivityLog{},
		&WebhookEvent{},
		&FeatureFlag{},
		&FlagExposure{},
		&Notification{},
		// grit:models
	}
}

// Migrate runs AutoMigrate for every registered model. For tables that
// already exist, GORM ALTERs them to add missing columns — we snapshot
// the column set before and after so the deploy log surfaces exactly
// what changed. Silent migrations are gone: if a column you expected
// didn't land, the diff makes it obvious.
//
//	================================================================
//	DATABASE MIGRATION — 8 model(s) registered
//	================================================================
//	  + created models.Building
//	  ~ models.User — added 2 column(s): is_vip, vip_notes
//	----------------------------------------------------------------
//	Migration done — 1 table(s) created, 1 altered (+2 column(s)), 6 unchanged.
//	================================================================
func Migrate(db *gorm.DB) error {
	models := Models()
	separator := strings.Repeat("=", 64)
	thinSep := strings.Repeat("-", 64)

	log.Println(separator)
	log.Printf("DATABASE MIGRATION — %d model(s) registered", len(models))
	log.Println(separator)

	// Silent logger keeps the schema-inspection SQL noise out of the diff log.
	silentDB := db.Session(&gorm.Session{Logger: logger.Default.LogMode(logger.Silent)})
	mig := silentDB.Migrator()

	created := 0
	altered := 0
	columnsAdded := 0
	unchanged := 0

	for _, model := range models {
		existed := mig.HasTable(model)

		var before map[string]bool
		if existed {
			before = make(map[string]bool)
			cols, err := mig.ColumnTypes(model)
			if err == nil {
				for _, c := range cols {
					before[c.Name()] = true
				}
			}
		}

		if err := silentDB.AutoMigrate(model); err != nil {
			return fmt.Errorf("migrating %T: %w", model, err)
		}

		if !existed {
			log.Printf("  + created %T", model)
			created++
			continue
		}

		// Diff columns to surface anything AutoMigrate added.
		after, err := mig.ColumnTypes(model)
		if err != nil {
			unchanged++
			continue
		}
		var added []string
		for _, c := range after {
			if !before[c.Name()] {
				added = append(added, c.Name())
			}
		}
		if len(added) == 0 {
			unchanged++
			continue
		}
		log.Printf("  ~ %T — added %d column(s): %s", model, len(added), strings.Join(added, ", "))
		altered++
		columnsAdded += len(added)
	}

	log.Println(thinSep)
	log.Printf("Migration done — %d table(s) created, %d altered (+%d column(s)), %d unchanged.",
		created, altered, columnsAdded, unchanged)
	log.Println(separator)
	return nil
}
`
}

func apiUploadModelGo() string {
	return `package models

import (
	"time"

	"gorm.io/gorm"
)

// Upload represents a file uploaded to storage.
type Upload struct {
	ID           string         ` + "`" + `gorm:"primarykey;size:36" json:"id"` + "`" + `
	Filename     string         ` + "`" + `gorm:"size:255;not null" json:"filename"` + "`" + `
	OriginalName string         ` + "`" + `gorm:"size:255;not null" json:"original_name"` + "`" + `
	MimeType     string         ` + "`" + `gorm:"size:100;not null" json:"mime_type"` + "`" + `
	Size         int64          ` + "`" + `gorm:"not null" json:"size"` + "`" + `
	Path         string         ` + "`" + `gorm:"size:500;not null" json:"path"` + "`" + `
	URL          string         ` + "`" + `gorm:"size:500" json:"url"` + "`" + `
	ThumbnailURL string         ` + "`" + `gorm:"size:500" json:"thumbnail_url"` + "`" + `
	UserID       string         ` + "`" + `gorm:"size:36;index;not null" json:"user_id"` + "`" + `
	User         User           ` + "`" + `gorm:"foreignKey:UserID" json:"-"` + "`" + `
	Version      int            ` + "`" + `gorm:"not null;default:1" json:"version"` + "`" + `
	CreatedAt    time.Time      ` + "`" + `json:"created_at"` + "`" + `
	UpdatedAt    time.Time      ` + "`" + `json:"updated_at"` + "`" + `
	DeletedAt    gorm.DeletedAt ` + "`" + `gorm:"index" json:"-"` + "`" + `
}

// BeforeUpdate increments Version on every server-side write so offline
// clients can detect that a record they edited has moved on.
func (u *Upload) BeforeUpdate(tx *gorm.DB) error {
	u.Version++
	return nil
}
`
}

func apiAuthServiceGo() string {
	return `package services

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// AuthService handles JWT token operations.
type AuthService struct {
	Secret        string
	AccessExpiry  time.Duration
	RefreshExpiry time.Duration
}

// TokenPair holds access and refresh tokens.
type TokenPair struct {
	AccessToken  string ` + "`" + `json:"access_token"` + "`" + `
	RefreshToken string ` + "`" + `json:"refresh_token"` + "`" + `
	ExpiresAt    int64  ` + "`" + `json:"expires_at"` + "`" + `
}

// Claims represents JWT claims.
type Claims struct {
	UserID string ` + "`" + `json:"user_id"` + "`" + `
	Email  string ` + "`" + `json:"email"` + "`" + `
	Role   string ` + "`" + `json:"role"` + "`" + `
	jwt.RegisteredClaims
}

// GenerateTokenPair creates a new access + refresh token pair.
func (s *AuthService) GenerateTokenPair(userID string, email, role string) (*TokenPair, error) {
	accessToken, expiresAt, err := s.generateToken(userID, email, role, s.AccessExpiry)
	if err != nil {
		return nil, fmt.Errorf("generating access token: %w", err)
	}

	refreshToken, _, err := s.generateToken(userID, email, role, s.RefreshExpiry)
	if err != nil {
		return nil, fmt.Errorf("generating refresh token: %w", err)
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt,
	}, nil
}

// ValidateToken parses and validates a JWT token.
func (s *AuthService) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.Secret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("parsing token: %w", err)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}

// GenerateResetToken creates a random hex token for password resets.
func GenerateResetToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("generating reset token: %w", err)
	}
	return hex.EncodeToString(bytes), nil
}

func (s *AuthService) generateToken(userID string, email, role string, expiry time.Duration) (string, int64, error) {
	expiresAt := time.Now().Add(expiry)

	claims := &Claims{
		UserID: userID,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.Secret))
	if err != nil {
		return "", 0, err
	}

	return tokenString, expiresAt.Unix(), nil
}

// SetAuthCookies writes the token pair as HttpOnly cookies so the browser
// holds the credentials out of JavaScript's reach. The native mobile and
// desktop clients keep using the Authorization: Bearer header, which is
// why the JSON body still includes the tokens — both paths work.
//
// Cookie names: grit_access (sent on every request) and grit_refresh
// (scoped to /api/auth so it isn't sent everywhere). Both are HttpOnly,
// Secure when on HTTPS, and SameSite=Lax so CSRF surface is limited to
// top-level navigations. The CSRF middleware adds defence in depth.
//
// Reference: docs/backend/authentication §"Token Storage on the Frontend".
func (s *AuthService) SetAuthCookies(c *gin.Context, pair *TokenPair) {
	secure := isRequestHTTPS(c)
	accessSeconds := int(s.AccessExpiry / time.Second)
	refreshSeconds := int(s.RefreshExpiry / time.Second)

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("grit_access", pair.AccessToken, accessSeconds, "/", "", secure, true)
	c.SetCookie("grit_refresh", pair.RefreshToken, refreshSeconds, "/api/auth", "", secure, true)
}

// ClearAuthCookies expires both auth cookies. Call this from the Logout
// handler so a stolen browser session is cut off as soon as the user
// signs out.
func (s *AuthService) ClearAuthCookies(c *gin.Context) {
	secure := isRequestHTTPS(c)
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("grit_access", "", -1, "/", "", secure, true)
	c.SetCookie("grit_refresh", "", -1, "/api/auth", "", secure, true)
}

// isRequestHTTPS returns true when the request is on HTTPS (directly or
// via a trusted proxy that set X-Forwarded-Proto=https). We use it to flip
// the Secure cookie flag so the browser refuses to send these cookies
// over an unencrypted hop.
func isRequestHTTPS(c *gin.Context) bool {
	if c.Request.TLS != nil {
		return true
	}
	if strings.EqualFold(c.GetHeader("X-Forwarded-Proto"), "https") {
		return true
	}
	return false
}
`
}

func apiAuthHandlerGo() string {
	return `package handlers

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/markbates/goth/gothic"
	"gorm.io/gorm"

	"golang.org/x/crypto/bcrypt"

	"` + "{{MODULE}}" + `/internal/config"
	"` + "{{MODULE}}" + `/internal/models"
	"` + "{{MODULE}}" + `/internal/services"
	"` + "{{MODULE}}" + `/internal/totp"
)

// AuthHandler handles authentication endpoints.
type AuthHandler struct {
	DB          *gorm.DB
	AuthService *services.AuthService
	Config      *config.Config
}

type registerRequest struct {
	FirstName  string ` + "`" + `json:"first_name" binding:"required,min=2"` + "`" + `
	LastName   string ` + "`" + `json:"last_name" binding:"required,min=2"` + "`" + `
	Email      string ` + "`" + `json:"email" binding:"required,email"` + "`" + `
	Password   string ` + "`" + `json:"password" binding:"required,min=8"` + "`" + `
	MACAddress string ` + "`" + `json:"mac_address"` + "`" + ` // optional — provided by client if available
}

type loginRequest struct {
	Email    string ` + "`" + `json:"email" binding:"required,email"` + "`" + `
	Password string ` + "`" + `json:"password" binding:"required"` + "`" + `
}

type refreshRequest struct {
	RefreshToken string ` + "`" + `json:"refresh_token" binding:"required"` + "`" + `
}

type forgotPasswordRequest struct {
	Email string ` + "`" + `json:"email" binding:"required,email"` + "`" + `
}

type resetPasswordRequest struct {
	Token    string ` + "`" + `json:"token" binding:"required"` + "`" + `
	Password string ` + "`" + `json:"password" binding:"required,min=8"` + "`" + `
}

// Register creates a new user account.
func (h *AuthHandler) Register(c *gin.Context) {
	var req registerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": gin.H{
				"code":    "VALIDATION_ERROR",
				"message": err.Error(),
			},
		})
		return
	}

	// Check email uniqueness
	var existingUser models.User
	if err := h.DB.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{
			"error": gin.H{
				"code":    "EMAIL_EXISTS",
				"message": "A user with this email already exists",
			},
		})
		return
	}

	user := models.User{
		FirstName:  req.FirstName,
		LastName:   req.LastName,
		Email:      req.Email,
		Password:   req.Password,
		Role:       models.RoleUser,
		Active:     true,
		IPAddress:  c.ClientIP(),
		MACAddress: req.MACAddress,
	}

	if err := h.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to create user",
			},
		})
		return
	}

	tokens, err := h.AuthService.GenerateTokenPair(user.ID, user.Email, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "TOKEN_ERROR",
				"message": "Failed to generate tokens",
			},
		})
		return
	}

	// Set HttpOnly auth cookies for browser clients.
	h.AuthService.SetAuthCookies(c, tokens)

	c.JSON(http.StatusCreated, gin.H{
		"data": gin.H{
			"user":   user,
			"tokens": tokens,
		},
		"message": "User registered successfully",
	})
}

// Login authenticates a user and returns tokens.
func (h *AuthHandler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": gin.H{
				"code":    "VALIDATION_ERROR",
				"message": err.Error(),
			},
		})
		return
	}

	var user models.User
	if err := h.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": gin.H{
				"code":    "INVALID_CREDENTIALS",
				"message": "Invalid email or password",
			},
		})
		return
	}

	if !user.Active {
		c.JSON(http.StatusForbidden, gin.H{
			"error": gin.H{
				"code":    "ACCOUNT_DISABLED",
				"message": "Your account has been disabled",
			},
		})
		return
	}

	if user.Password == "" {
		provider := user.Provider
		if provider == "" || provider == "local" {
			provider = "social login"
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "SOCIAL_AUTH_ONLY",
				"message": fmt.Sprintf("This account uses %s. Please sign in with your social account.", provider),
			},
		})
		return
	}

	if !user.CheckPassword(req.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": gin.H{
				"code":    "INVALID_CREDENTIALS",
				"message": "Invalid email or password",
			},
		})
		return
	}

	// Check if user has TOTP enabled
	var totpConfig models.TwoFactorConfig
	if err := h.DB.Where("user_id = ? AND enabled = ?", user.ID, true).First(&totpConfig).Error; err == nil {
		// TOTP is enabled — check for trusted device
		if !IsTrustedDevice(c, h.DB, user.ID) {
			// Generate a short-lived pending token for TOTP verification
			pendingToken, err := totp.GeneratePendingToken()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": gin.H{"code": "TOKEN_ERROR", "message": "Failed to create verification session"},
				})
				return
			}

			// Store hashed pending token in DB
			h.DB.Create(&models.TOTPPendingToken{
				UserID:    user.ID,
				TokenHash: totp.HashToken(pendingToken),
				ExpiresAt: time.Now().Add(totp.PendingTokenExpiry),
			})

			c.JSON(http.StatusOK, gin.H{
				"data": gin.H{
					"totp_required": true,
					"pending_token": pendingToken,
				},
				"message": "Two-factor authentication required",
			})
			return
		}
	}

	tokens, err := h.AuthService.GenerateTokenPair(user.ID, user.Email, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "TOKEN_ERROR",
				"message": "Failed to generate tokens",
			},
		})
		return
	}

	// Set HttpOnly auth cookies for browser clients. Native mobile/desktop
	// clients ignore them and continue to use the Bearer header from the
	// tokens object below — both flows work.
	h.AuthService.SetAuthCookies(c, tokens)

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"user":   user,
			"tokens": tokens,
		},
		"message": "Logged in successfully",
	})
}

// Refresh generates a new access token from a refresh token. The token is
// read from the grit_refresh cookie first (web client) and falls back to
// the JSON body (mobile/desktop bearer clients) — so a single endpoint
// supports both flows.
func (h *AuthHandler) Refresh(c *gin.Context) {
	refreshToken := ""
	if cookieValue, err := c.Cookie("grit_refresh"); err == nil && cookieValue != "" {
		refreshToken = cookieValue
	} else {
		var req refreshRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"error": gin.H{
					"code":    "VALIDATION_ERROR",
					"message": err.Error(),
				},
			})
			return
		}
		refreshToken = req.RefreshToken
	}

	claims, err := h.AuthService.ValidateToken(refreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": gin.H{
				"code":    "INVALID_TOKEN",
				"message": "Invalid or expired refresh token",
			},
		})
		return
	}

	tokens, err := h.AuthService.GenerateTokenPair(claims.UserID, claims.Email, claims.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "TOKEN_ERROR",
				"message": "Failed to generate tokens",
			},
		})
		return
	}

	// Refresh the HttpOnly cookies so the new access token lands in the
	// browser without any JS handling. The bearer JSON path is unchanged
	// for native clients.
	h.AuthService.SetAuthCookies(c, tokens)

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"tokens": tokens,
		},
		"message": "Token refreshed successfully",
	})
}

// Logout invalidates the user's session. Cookies are cleared immediately;
// native bearer clients should also drop their stored tokens client-side.
func (h *AuthHandler) Logout(c *gin.Context) {
	h.AuthService.ClearAuthCookies(c)
	// In a production system, you'd also blacklist the refresh token in Redis
	// so a leaked token can't be reused before its natural expiry.
	c.JSON(http.StatusOK, gin.H{
		"message": "Logged out successfully",
	})
}

// Me returns the current authenticated user.
func (h *AuthHandler) Me(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": gin.H{
				"code":    "UNAUTHORIZED",
				"message": "Not authenticated",
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": user,
	})
}

// ForgotPassword initiates a password reset.
func (h *AuthHandler) ForgotPassword(c *gin.Context) {
	var req forgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": gin.H{
				"code":    "VALIDATION_ERROR",
				"message": err.Error(),
			},
		})
		return
	}

	var user models.User
	if err := h.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		// Return success even if email not found (security)
		c.JSON(http.StatusOK, gin.H{
			"message": "If an account with that email exists, a password reset link has been sent",
		})
		return
	}

	token, err := services.GenerateResetToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to generate reset token",
			},
		})
		return
	}

	// For Phase 1, just log the token (email integration comes in Phase 4)
	log.Printf("Password reset token for %s: %s", user.Email, token)

	c.JSON(http.StatusOK, gin.H{
		"message": "If an account with that email exists, a password reset link has been sent",
	})
}

// ResetPassword resets a user's password with a valid token.
func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var req resetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": gin.H{
				"code":    "VALIDATION_ERROR",
				"message": err.Error(),
			},
		})
		return
	}

	// Phase 1: simplified reset (in production, validate the token against stored tokens)
	// For now, this is a placeholder that demonstrates the API contract
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to hash password",
			},
		})
		return
	}
	_ = hashedPassword

	c.JSON(http.StatusOK, gin.H{
		"message": "Password reset successfully",
	})
}

// OAuthBegin redirects the user to the OAuth provider's consent screen.
func (h *AuthHandler) OAuthBegin(c *gin.Context) {
	provider := c.Param("provider")

	// Gothic reads provider from query string, not URL params
	q := c.Request.URL.Query()
	q.Set("provider", provider)
	c.Request.URL.RawQuery = q.Encode()

	gothic.BeginAuthHandler(c.Writer, c.Request)
}

// OAuthCallback completes the OAuth flow, finds or creates the user, and redirects with JWT tokens.
func (h *AuthHandler) OAuthCallback(c *gin.Context) {
	provider := c.Param("provider")

	q := c.Request.URL.Query()
	q.Set("provider", provider)
	c.Request.URL.RawQuery = q.Encode()

	gothUser, err := gothic.CompleteUserAuth(c.Writer, c.Request)
	if err != nil {
		log.Printf("OAuth callback error: %v", err)
		redirectURL := fmt.Sprintf("%s/login?error=%s", h.Config.OAuthFrontendURL, url.QueryEscape("Authentication failed. Please try again."))
		c.Redirect(http.StatusTemporaryRedirect, redirectURL)
		return
	}

	// Find or create user by email
	var user models.User
	result := h.DB.Where("email = ?", gothUser.Email).First(&user)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			// Create new user from OAuth data
			now := time.Now()
			user = models.User{
				FirstName:       gothUser.FirstName,
				LastName:        gothUser.LastName,
				Email:           gothUser.Email,
				Avatar:          gothUser.AvatarURL,
				Provider:        provider,
				Active:          true,
				EmailVerifiedAt: &now,
				IPAddress:       c.ClientIP(),
			}

			if provider == "google" {
				user.GoogleID = gothUser.UserID
			} else if provider == "github" {
				user.GithubID = gothUser.UserID
			}

			// If name is empty, try to use NickName
			if user.FirstName == "" && gothUser.NickName != "" {
				user.FirstName = gothUser.NickName
			}
			if user.FirstName == "" {
				user.FirstName = "User"
			}
			if user.LastName == "" {
				user.LastName = ""
			}

			if err := h.DB.Create(&user).Error; err != nil {
				log.Printf("OAuth: failed to create user: %v", err)
				redirectURL := fmt.Sprintf("%s/login?error=%s", h.Config.OAuthFrontendURL, url.QueryEscape("Failed to create account."))
				c.Redirect(http.StatusTemporaryRedirect, redirectURL)
				return
			}
		} else {
			log.Printf("OAuth: database error: %v", result.Error)
			redirectURL := fmt.Sprintf("%s/login?error=%s", h.Config.OAuthFrontendURL, url.QueryEscape("Something went wrong."))
			c.Redirect(http.StatusTemporaryRedirect, redirectURL)
			return
		}
	} else {
		// Link OAuth provider to existing account
		updates := map[string]interface{}{}
		if provider == "google" && user.GoogleID == "" {
			updates["google_id"] = gothUser.UserID
		} else if provider == "github" && user.GithubID == "" {
			updates["github_id"] = gothUser.UserID
		}
		if user.Avatar == "" && gothUser.AvatarURL != "" {
			updates["avatar"] = gothUser.AvatarURL
		}
		if user.Provider == "local" {
			updates["provider"] = provider
		}

		if len(updates) > 0 {
			h.DB.Model(&user).Updates(updates)
		}
	}

	if !user.Active {
		redirectURL := fmt.Sprintf("%s/login?error=%s", h.Config.OAuthFrontendURL, url.QueryEscape("Your account has been disabled."))
		c.Redirect(http.StatusTemporaryRedirect, redirectURL)
		return
	}

	// Generate JWT tokens
	tokens, err := h.AuthService.GenerateTokenPair(user.ID, user.Email, user.Role)
	if err != nil {
		log.Printf("OAuth: failed to generate tokens: %v", err)
		redirectURL := fmt.Sprintf("%s/login?error=%s", h.Config.OAuthFrontendURL, url.QueryEscape("Failed to sign in."))
		c.Redirect(http.StatusTemporaryRedirect, redirectURL)
		return
	}

	// Redirect to frontend with tokens
	redirectURL := fmt.Sprintf("%s/auth/callback?access_token=%s&refresh_token=%s",
		h.Config.OAuthFrontendURL,
		url.QueryEscape(tokens.AccessToken),
		url.QueryEscape(tokens.RefreshToken),
	)
	c.Redirect(http.StatusTemporaryRedirect, redirectURL)
}
`
}

func apiUserHandlerGo() string {
	return `package handlers

import (
	"math"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"` + "{{MODULE}}" + `/internal/models"
)

// UserHandler handles user management endpoints.
type UserHandler struct {
	DB *gorm.DB
}

// Create creates a new user (admin only).
func (h *UserHandler) Create(c *gin.Context) {
	var req struct {
		FirstName string ` + "`" + `json:"first_name" binding:"required"` + "`" + `
		LastName  string ` + "`" + `json:"last_name" binding:"required"` + "`" + `
		Email     string ` + "`" + `json:"email" binding:"required,email"` + "`" + `
		Password  string ` + "`" + `json:"password" binding:"required,min=6"` + "`" + `
		Role      string ` + "`" + `json:"role"` + "`" + `
		Avatar    string ` + "`" + `json:"avatar"` + "`" + `
		JobTitle  string ` + "`" + `json:"job_title"` + "`" + `
		Active    *bool  ` + "`" + `json:"active"` + "`" + `
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": gin.H{
				"code":    "VALIDATION_ERROR",
				"message": err.Error(),
			},
		})
		return
	}

	// Check email uniqueness
	var existing models.User
	if err := h.DB.Where("email = ?", req.Email).First(&existing).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{
			"error": gin.H{
				"code":    "EMAIL_EXISTS",
				"message": "A user with this email already exists",
			},
		})
		return
	}

	user := models.User{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Password:  req.Password,
		Role:      req.Role,
		Avatar:    req.Avatar,
		JobTitle:  req.JobTitle,
		Active:    true,
	}

	if req.Active != nil {
		user.Active = *req.Active
	}
	if user.Role == "" {
		user.Role = models.RoleUser
	}

	if err := h.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to create user",
			},
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data":    user,
		"message": "User created successfully",
	})
}

// List returns a paginated list of users.
func (h *UserHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	search := c.Query("search")
	sortBy := c.DefaultQuery("sort_by", "created_at")
	sortOrder := c.DefaultQuery("sort_order", "desc")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	// Validate sort order
	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = "desc"
	}

	// Validate sort column
	allowedSorts := map[string]bool{
		"id": true, "first_name": true, "last_name": true, "email": true, "role": true, "created_at": true,
	}
	if !allowedSorts[sortBy] {
		sortBy = "created_at"
	}

	query := h.DB.Model(&models.User{})

	// Search
	if search != "" {
		query = query.Where("first_name ILIKE ? OR last_name ILIKE ? OR email ILIKE ?", "%"+search+"%", "%"+search+"%", "%"+search+"%")
	}

	// Count total
	var total int64
	query.Count(&total)

	// Fetch paginated results
	var users []models.User
	offset := (page - 1) * pageSize
	if err := query.Order(sortBy + " " + sortOrder).Offset(offset).Limit(pageSize).Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to fetch users",
			},
		})
		return
	}

	pages := int(math.Ceil(float64(total) / float64(pageSize)))

	c.JSON(http.StatusOK, gin.H{
		"data": users,
		"meta": gin.H{
			"total":     total,
			"page":      page,
			"page_size": pageSize,
			"pages":     pages,
		},
	})
}

// GetByID returns a single user by ID.
func (h *UserHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	var user models.User
	if err := h.DB.Where("id = ?", id).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": gin.H{
				"code":    "NOT_FOUND",
				"message": "User not found",
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": user,
	})
}

// Update modifies an existing user.
func (h *UserHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var user models.User
	if err := h.DB.Where("id = ?", id).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": gin.H{
				"code":    "NOT_FOUND",
				"message": "User not found",
			},
		})
		return
	}

	var req struct {
		FirstName string ` + "`" + `json:"first_name"` + "`" + `
		LastName  string ` + "`" + `json:"last_name"` + "`" + `
		Email     string ` + "`" + `json:"email"` + "`" + `
		Password  string ` + "`" + `json:"password"` + "`" + `
		Role      string ` + "`" + `json:"role"` + "`" + `
		Avatar    string ` + "`" + `json:"avatar"` + "`" + `
		JobTitle  string ` + "`" + `json:"job_title"` + "`" + `
		Bio       string ` + "`" + `json:"bio"` + "`" + `
		Active    *bool  ` + "`" + `json:"active"` + "`" + `
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": gin.H{
				"code":    "VALIDATION_ERROR",
				"message": err.Error(),
			},
		})
		return
	}

	updates := map[string]interface{}{}
	if req.FirstName != "" {
		updates["first_name"] = req.FirstName
	}
	if req.LastName != "" {
		updates["last_name"] = req.LastName
	}
	if req.Email != "" {
		updates["email"] = req.Email
	}
	if req.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": gin.H{
					"code":    "INTERNAL_ERROR",
					"message": "Failed to hash password",
				},
			})
			return
		}
		updates["password"] = string(hashedPassword)
	}
	if req.Role != "" {
		updates["role"] = req.Role
	}
	if req.Avatar != "" {
		updates["avatar"] = req.Avatar
	}
	if req.JobTitle != "" {
		updates["job_title"] = req.JobTitle
	}
	if req.Bio != "" {
		updates["bio"] = req.Bio
	}
	if req.Active != nil {
		updates["active"] = *req.Active
	}

	if err := h.DB.Model(&user).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to update user",
			},
		})
		return
	}

	// Reload to get updated values
	h.DB.Where("id = ?", id).First(&user)

	c.JSON(http.StatusOK, gin.H{
		"data":    user,
		"message": "User updated successfully",
	})
}

// Delete soft-deletes a user.
func (h *UserHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	var user models.User
	if err := h.DB.Where("id = ?", id).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": gin.H{
				"code":    "NOT_FOUND",
				"message": "User not found",
			},
		})
		return
	}

	if err := h.DB.Delete(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to delete user",
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User deleted successfully",
	})
}

// GetProfile returns the currently authenticated user's profile.
func (h *UserHandler) GetProfile(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var user models.User
	if err := h.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": gin.H{
				"code":    "NOT_FOUND",
				"message": "User not found",
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": user,
	})
}

// UpdateProfile updates the currently authenticated user's profile.
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var user models.User
	if err := h.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": gin.H{
				"code":    "NOT_FOUND",
				"message": "User not found",
			},
		})
		return
	}

	var req struct {
		FirstName string ` + "`" + `json:"first_name"` + "`" + `
		LastName  string ` + "`" + `json:"last_name"` + "`" + `
		Email     string ` + "`" + `json:"email"` + "`" + `
		Password  string ` + "`" + `json:"password"` + "`" + `
		Avatar    string ` + "`" + `json:"avatar"` + "`" + `
		JobTitle  string ` + "`" + `json:"job_title"` + "`" + `
		Bio       string ` + "`" + `json:"bio"` + "`" + `
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": gin.H{
				"code":    "VALIDATION_ERROR",
				"message": err.Error(),
			},
		})
		return
	}

	updates := map[string]interface{}{}
	if req.FirstName != "" {
		updates["first_name"] = req.FirstName
	}
	if req.LastName != "" {
		updates["last_name"] = req.LastName
	}
	if req.Email != "" {
		updates["email"] = req.Email
	}
	if req.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": gin.H{
					"code":    "INTERNAL_ERROR",
					"message": "Failed to hash password",
				},
			})
			return
		}
		updates["password"] = string(hashedPassword)
	}
	if req.Avatar != "" {
		updates["avatar"] = req.Avatar
	}
	if req.JobTitle != "" {
		updates["job_title"] = req.JobTitle
	}
	if req.Bio != "" {
		updates["bio"] = req.Bio
	}

	if err := h.DB.Model(&user).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to update profile",
			},
		})
		return
	}

	h.DB.Where("id = ?", userID).First(&user)

	c.JSON(http.StatusOK, gin.H{
		"data":    user,
		"message": "Profile updated successfully",
	})
}

// DeleteProfile soft-deletes the currently authenticated user's account.
func (h *UserHandler) DeleteProfile(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var user models.User
	if err := h.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": gin.H{
				"code":    "NOT_FOUND",
				"message": "User not found",
			},
		})
		return
	}

	if err := h.DB.Delete(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to delete account",
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Account deleted successfully",
	})
}
`
}

func apiAuthMiddlewareGo() string {
	return `package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"` + "{{MODULE}}" + `/internal/models"
	"` + "{{MODULE}}" + `/internal/services"
)

// Auth creates a JWT authentication middleware.
func Auth(db *gorm.DB, authService *services.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Resolve the access token. The HttpOnly cookie path is the
		// recommended flow for browser clients — JS never sees the token,
		// so XSS cannot exfiltrate it. The Authorization: Bearer header
		// path is the fallback for native mobile / desktop clients that
		// can't or don't want to use cookies.
		token := ""
		if cookieValue, err := c.Cookie("grit_access"); err == nil && cookieValue != "" {
			token = cookieValue
		} else if authHeader := c.GetHeader("Authorization"); authHeader != "" {
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || parts[0] != "Bearer" {
				c.JSON(http.StatusUnauthorized, gin.H{
					"error": gin.H{
						"code":    "UNAUTHORIZED",
						"message": "Invalid authorization header format",
					},
				})
				c.Abort()
				return
			}
			token = parts[1]
		}

		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": gin.H{
					"code":    "UNAUTHORIZED",
					"message": "Authentication required",
				},
			})
			c.Abort()
			return
		}

		claims, err := authService.ValidateToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": gin.H{
					"code":    "UNAUTHORIZED",
					"message": "Invalid or expired token",
				},
			})
			c.Abort()
			return
		}

		// Load user from database.
		// Use Where("id = ?") rather than First(&user, id) — GORM's shorthand
		// emits the bare value into the WHERE clause and Postgres rejects UUID
		// primary keys with "trailing junk after numeric literal".
		var user models.User
		if err := db.Where("id = ?", claims.UserID).First(&user).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": gin.H{
					"code":    "UNAUTHORIZED",
					"message": "User not found",
				},
			})
			c.Abort()
			return
		}

		if !user.Active {
			c.JSON(http.StatusForbidden, gin.H{
				"error": gin.H{
					"code":    "ACCOUNT_DISABLED",
					"message": "Your account has been disabled",
				},
			})
			c.Abort()
			return
		}

		c.Set("user", user)
		c.Set("user_id", user.ID)
		c.Set("user_role", user.Role)
		c.Next()
	}
}

// RequireRole creates a middleware that checks if the user has one of the required roles.
func RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("user_role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": gin.H{
					"code":    "UNAUTHORIZED",
					"message": "Not authenticated",
				},
			})
			c.Abort()
			return
		}

		role, ok := userRole.(string)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": gin.H{
					"code":    "INTERNAL_ERROR",
					"message": "Invalid user role",
				},
			})
			c.Abort()
			return
		}

		for _, r := range roles {
			if role == r {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, gin.H{
			"error": gin.H{
				"code":    "FORBIDDEN",
				"message": "You do not have permission to access this resource",
			},
		})
		c.Abort()
	}
}
`
}

func apiCorsMiddlewareGo() string {
	return `package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// CORS creates a CORS middleware with the given allowed origins.
func CORS(allowedOrigins []string) gin.HandlerFunc {
	originsMap := make(map[string]bool)
	for _, origin := range allowedOrigins {
		originsMap[origin] = true
	}

	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")

		if originsMap[origin] {
			c.Header("Access-Control-Allow-Origin", origin)
		}

		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Max-Age", "86400")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
`
}

func apiLoggerMiddlewareGo() string {
	return `package middleware

import (
	"compress/gzip"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// RequestID injects a unique X-Request-ID header into every request and
// stores it in the context for downstream logging and tracing.
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = fmt.Sprintf("%d-%d", time.Now().UnixNano(), rand.Int63())
		}
		c.Set("request_id", requestID)
		c.Header("X-Request-ID", requestID)
		c.Next()
	}
}

// Gzip compresses responses using gzip encoding when the client supports it.
func Gzip() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !strings.Contains(c.GetHeader("Accept-Encoding"), "gzip") {
			c.Next()
			return
		}

		gz, err := gzip.NewWriterLevel(c.Writer, gzip.BestSpeed)
		if err != nil {
			c.Next()
			return
		}
		defer gz.Close()

		c.Header("Content-Encoding", "gzip")
		c.Header("Vary", "Accept-Encoding")
		c.Writer = &gzipResponseWriter{ResponseWriter: c.Writer, Writer: gz}
		c.Next()
	}
}

type gzipResponseWriter struct {
	gin.ResponseWriter
	Writer *gzip.Writer
}

func (g *gzipResponseWriter) Write(data []byte) (int, error) {
	return g.Writer.Write(data)
}

func (g *gzipResponseWriter) WriteString(s string) (int, error) {
	return g.Writer.Write([]byte(s))
}

// SecurityHeaders adds production security headers to every response.
//
// Coverage against OWASP Top 10:2025 — A02 Security Misconfiguration,
// A05 Injection (XSS hardening via CSP), A04 Cryptographic Failures
// (HSTS forces TLS), plus mitigations for clickjacking, MIME sniffing,
// referrer leakage, and Spectre-class cross-origin attacks.
//
// CSP is deliberately strict-by-default. The scaffold's SPA serves /api
// from the same origin, so 'self' covers the normal case. Customise
// CSPDirectives via config when adding a CDN / inline scripts.
func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Header("Permissions-Policy", "camera=(), microphone=(), geolocation=(), payment=(), usb=()")
		// Spectre-class defence: isolate this origin from cross-origin reads
		// and require explicit opt-in for cross-origin embedders.
		c.Header("Cross-Origin-Opener-Policy", "same-origin")
		c.Header("Cross-Origin-Resource-Policy", "same-origin")
		// Content-Security-Policy — strict default, blocks inline script
		// (XSS A05 hardening). Skip on /docs and /studio which serve
		// vendored UIs that rely on inline styles.
		path := c.Request.URL.Path
		if !strings.HasPrefix(path, "/docs") && !strings.HasPrefix(path, "/studio") && !strings.HasPrefix(path, "/sentinel") && !strings.HasPrefix(path, "/pulse") {
			c.Header("Content-Security-Policy",
				"default-src 'self'; "+
					"script-src 'self'; "+
					"style-src 'self' 'unsafe-inline'; "+
					"img-src 'self' data: blob: https:; "+
					"font-src 'self' data:; "+
					"connect-src 'self'; "+
					"frame-ancestors 'none'; "+
					"base-uri 'self'; "+
					"form-action 'self'; "+
					"object-src 'none'")
		}
		// HSTS only when actually on HTTPS (don't break dev on http://).
		if c.Request.TLS != nil || c.GetHeader("X-Forwarded-Proto") == "https" {
			c.Header("Strict-Transport-Security", "max-age=63072000; includeSubDomains; preload")
		}
		c.Next()
	}
}

// MaxBodySize limits the request body to prevent abuse.
func MaxBodySize(limit int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.ContentLength > limit {
			c.AbortWithStatusJSON(http.StatusRequestEntityTooLarge, gin.H{
				"error": gin.H{
					"code":    "PAYLOAD_TOO_LARGE",
					"message": fmt.Sprintf("Request body exceeds %dMB limit", limit/(1024*1024)),
				},
			})
			return
		}
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, limit)
		c.Next()
	}
}

// Logger creates a structured logging middleware with request ID correlation.
// Silently skips internal dashboard paths to keep the terminal readable.
func Logger() gin.HandlerFunc {
	// Paths that generate noise and aren't useful to see in dev logs
	skipPrefixes := []string{
		"/studio/",
		"/pulse/",
		"/pulse",
		"/sentinel/",
		"/docs/",
		"/docs",
		"/r.json",
		"/r/",
		"/api/health",
		"/favicon.ico",
	}

	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path

		// Skip noisy internal paths
		for _, prefix := range skipPrefixes {
			if strings.HasPrefix(path, prefix) || path == prefix {
				c.Next()
				return
			}
		}

		query := c.Request.URL.RawQuery

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()
		method := c.Request.Method
		clientIP := c.ClientIP()
		requestID, _ := c.Get("request_id")

		if query != "" {
			path = path + "?" + query
		}

		log.Printf("[%d] %s %s | %s | %v | id=%v",
			status,
			method,
			path,
			clientIP,
			latency,
			requestID,
		)
	}
}
`
}

// apiPaginateGo returns the generic pagination/sort/search helper.
// Every generated resource's List endpoint uses paginate.List so that
// page-clamping, sort whitelisting, and search ILIKE construction live
// in exactly one place. Addresses issue #14.
func apiPaginateGo() string {
	return `// Package paginate provides a generic list/sort/search/paginate helper
// used by every resource's List endpoint. The goal: one source of truth
// for page clamping, sort whitelisting, and search construction so that
// new resources don't drift on the boilerplate. Works with any GORM model.
//
// Usage (handler side):
//
//	func (h *ShopHandler) List(c *gin.Context) {
//	    res, err := paginate.List[models.Shop](
//	        h.DB.Model(&models.Shop{}).Preload("Building"),
//	        paginate.Bind(c),
//	        paginate.Config{
//	            Searchable:   []string{"shop_number", "description"},
//	            Sortable:     map[string]bool{"created_at": true, "monthly_rent": true},
//	            DefaultSort:  "created_at",
//	            DefaultOrder: "desc",
//	        },
//	    )
//	    if err != nil {
//	        c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"code": "INTERNAL_ERROR", "message": err.Error()}})
//	        return
//	    }
//	    c.JSON(http.StatusOK, res)
//	}
package paginate

import (
	"encoding/base64"
	"fmt"
	"math"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Defaults applied when the request query is empty or out of range.
const (
	DefaultPage         = 1
	DefaultPageSize     = 20
	MaxPageSize         = 100
	DefaultSortColumn   = "created_at"
	DefaultSortOrder    = "desc"
)

// Params is the normalized query state for a List request.
// Produced by Bind(c). Filters is free-form extra WHERE col = val clauses.
// Cursor (when present) drives cursor-mode pagination — see Config.CursorMode.
type Params struct {
	Page      int
	PageSize  int
	Search    string
	SortBy    string
	SortOrder string
	Cursor    string // opaque base64 from a previous Result.Meta.NextCursor
	Filters   map[string]any
}

// With returns a copy of Params with an additional filter applied.
// Empty string values are ignored so handlers can pipe c.Query() directly.
//
//	paginate.Bind(c).With("building_id", c.Query("building_id"))
func (p Params) With(key string, value any) Params {
	if s, ok := value.(string); ok && s == "" {
		return p
	}
	if value == nil {
		return p
	}
	if p.Filters == nil {
		p.Filters = map[string]any{key: value}
		return p
	}
	// Copy the map so we don't mutate the caller's Params.
	copied := make(map[string]any, len(p.Filters)+1)
	for k, v := range p.Filters {
		copied[k] = v
	}
	copied[key] = value
	p.Filters = copied
	return p
}

// Config describes which columns the caller has declared searchable / sortable
// for a particular resource. Anything not in Sortable falls back to DefaultSort.
type Config struct {
	Searchable   []string        // columns included in ILIKE search
	Sortable     map[string]bool // whitelist for sort_by values
	DefaultSort  string          // fallback sort column (defaults to "created_at")
	DefaultOrder string          // fallback sort order (defaults to "desc")

	// CursorMode opts into cursor-based pagination (default is offset/page).
	// When true, the response carries Meta.NextCursor + Meta.HasMore instead
	// of Page/Pages/Total. Cursor is opaque base64 of (sort_value, id) so
	// pages stay stable when rows insert mid-pagination.
	CursorMode bool

	// IncludeTotal asks cursor mode to also run COUNT(*). Slow on big
	// tables — leave off unless your UI shows a "X of Y" indicator.
	IncludeTotal bool
}

// Meta is the pagination envelope, matching Grit's existing response shape.
// Cursor mode populates NextCursor + HasMore; offset mode populates
// Page + Pages. Total is shared (always set in offset mode; opt-in in
// cursor mode via Config.IncludeTotal).
type Meta struct {
	Total      int64  ` + "`" + `json:"total,omitempty"` + "`" + `
	Page       int    ` + "`" + `json:"page,omitempty"` + "`" + `
	PageSize   int    ` + "`" + `json:"page_size"` + "`" + `
	Pages      int    ` + "`" + `json:"pages,omitempty"` + "`" + `
	NextCursor string ` + "`" + `json:"next_cursor,omitempty"` + "`" + `
	HasMore    bool   ` + "`" + `json:"has_more,omitempty"` + "`" + `
}

// Result wraps the paginated data in the canonical { data, meta } envelope.
type Result[T any] struct {
	Data []T  ` + "`" + `json:"data"` + "`" + `
	Meta Meta ` + "`" + `json:"meta"` + "`" + `
}

// Bind reads page / page_size / search / sort_by / sort_order from the Gin
// context, clamps them, and returns a normalized Params.
func Bind(c *gin.Context) Params {
	page, _ := strconv.Atoi(c.DefaultQuery("page", strconv.Itoa(DefaultPage)))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", strconv.Itoa(DefaultPageSize)))

	if page < 1 {
		page = DefaultPage
	}
	if pageSize < 1 || pageSize > MaxPageSize {
		pageSize = DefaultPageSize
	}

	return Params{
		Page:      page,
		PageSize:  pageSize,
		Search:    c.Query("search"),
		SortBy:    c.Query("sort_by"),
		SortOrder: c.Query("sort_order"),
		Cursor:    c.Query("cursor"),
	}
}

// List runs the query with search / sort / filters / pagination applied and
// returns a typed Result[T]. The caller is expected to have already set the
// model and any relevant Preload() chains on the passed-in *gorm.DB.
//
// Invariants enforced:
//   - page >= 1, 1 <= page_size <= MaxPageSize
//   - sort_by must be in cfg.Sortable, else cfg.DefaultSort (or DefaultSortColumn)
//   - sort_order must be "asc" or "desc", else cfg.DefaultOrder (or DefaultSortOrder)
//   - search is applied as ILIKE across cfg.Searchable columns (nothing if empty)
func List[T any](query *gorm.DB, p Params, cfg Config) (Result[T], error) {
	// Normalize sort_by against the whitelist.
	sortBy := p.SortBy
	if sortBy == "" || !cfg.Sortable[sortBy] {
		sortBy = cfg.DefaultSort
		if sortBy == "" {
			sortBy = DefaultSortColumn
		}
	}

	// Normalize sort_order.
	sortOrder := p.SortOrder
	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = cfg.DefaultOrder
		if sortOrder == "" {
			sortOrder = DefaultSortOrder
		}
	}

	// Apply equality filters (e.g. ?status=active&building_id=...).
	for col, val := range p.Filters {
		query = query.Where(col+" = ?", val)
	}

	// Apply search across configured columns.
	if p.Search != "" && len(cfg.Searchable) > 0 {
		clause, args := buildSearchClause(cfg.Searchable, p.Search)
		query = query.Where(clause, args...)
	}

	if cfg.CursorMode {
		return listCursor[T](query, p, cfg, sortBy, sortOrder)
	}

	var result Result[T]

	// Count first (before Order/Offset/Limit so Count reflects the whole match).
	if err := query.Count(&result.Meta.Total).Error; err != nil {
		return result, err
	}

	// Then fetch the page.
	offset := (p.Page - 1) * p.PageSize
	if err := query.
		Order(sortBy + " " + sortOrder).
		Offset(offset).
		Limit(p.PageSize).
		Find(&result.Data).Error; err != nil {
		return result, err
	}

	result.Meta.Page = p.Page
	result.Meta.PageSize = p.PageSize
	result.Meta.Pages = 0
	if result.Meta.Total > 0 && p.PageSize > 0 {
		result.Meta.Pages = int(math.Ceil(float64(result.Meta.Total) / float64(p.PageSize)))
	}

	return result, nil
}

// listCursor implements cursor-based pagination. The cursor is an
// opaque base64 of (sort_value, id) so pages stay stable when rows
// insert mid-pagination. We fetch PageSize+1 rows to detect HasMore
// without a separate count query.
func listCursor[T any](query *gorm.DB, p Params, cfg Config, sortBy, sortOrder string) (Result[T], error) {
	var result Result[T]

	if cfg.IncludeTotal {
		countQuery := query.Session(&gorm.Session{})
		if err := countQuery.Count(&result.Meta.Total).Error; err != nil {
			return result, err
		}
	}

	if p.Cursor != "" {
		sortVal, lastID, err := decodeCursor(p.Cursor)
		if err == nil {
			op := "<"
			if sortOrder == "asc" {
				op = ">"
			}
			// Postgres tuple comparison: (sort_col, id) < (val, id).
			// Works on SQLite too. The id tiebreaker keeps the cursor
			// stable when sort_value collides on multiple rows.
			query = query.Where(fmt.Sprintf("(%s, id) %s (?, ?)", sortBy, op), sortVal, lastID)
		}
	}

	limit := p.PageSize + 1
	if err := query.
		Order(sortBy + " " + sortOrder).
		Order("id " + sortOrder).
		Limit(limit).
		Find(&result.Data).Error; err != nil {
		return result, err
	}

	if len(result.Data) > p.PageSize {
		result.Data = result.Data[:p.PageSize]
		result.Meta.HasMore = true
	}

	if len(result.Data) > 0 {
		last := result.Data[len(result.Data)-1]
		sortVal, id := extractCursor(last, sortBy)
		if id != "" {
			result.Meta.NextCursor = encodeCursor(sortVal, id)
		}
	}

	result.Meta.PageSize = p.PageSize
	return result, nil
}

// EncodeCursor / DecodeCursor are exported for handlers that build
// custom cursors (e.g. nested resource links).
func EncodeCursor(sortValue, id string) string { return encodeCursor(sortValue, id) }
func DecodeCursor(s string) (string, string, error) { return decodeCursor(s) }

func encodeCursor(sortVal, id string) string {
	return base64.URLEncoding.EncodeToString([]byte(sortVal + "|" + id))
}

func decodeCursor(s string) (string, string, error) {
	b, err := base64.URLEncoding.DecodeString(s)
	if err != nil {
		return "", "", fmt.Errorf("invalid cursor: %w", err)
	}
	parts := strings.SplitN(string(b), "|", 2)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid cursor format")
	}
	return parts[0], parts[1], nil
}

// extractCursor reflects on the last row to pull out the sort field
// + ID. The sort field is stored as snake_case (matching the column),
// so we convert to PascalCase for the Go struct field lookup.
func extractCursor(item interface{}, sortBy string) (string, string) {
	rv := reflect.ValueOf(item)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	if rv.Kind() != reflect.Struct {
		return "", ""
	}

	idVal := rv.FieldByName("ID")
	if !idVal.IsValid() || idVal.Kind() != reflect.String {
		return "", ""
	}
	id := idVal.String()

	goFieldName := snakeToPascal(sortBy)
	sortField := rv.FieldByName(goFieldName)
	if !sortField.IsValid() {
		return "", id
	}

	if t, ok := sortField.Interface().(time.Time); ok {
		return t.Format(time.RFC3339Nano), id
	}
	return fmt.Sprintf("%v", sortField.Interface()), id
}

// snakeToPascal turns "created_at" into "CreatedAt".
func snakeToPascal(s string) string {
	parts := strings.Split(s, "_")
	for i, p := range parts {
		if p == "" {
			continue
		}
		parts[i] = strings.ToUpper(p[:1]) + p[1:]
	}
	return strings.Join(parts, "")
}

// buildSearchClause builds "col1 ILIKE ? OR col2 ILIKE ? OR ..." with the
// same wildcard-wrapped search term repeated as each arg.
func buildSearchClause(cols []string, term string) (string, []any) {
	clause := ""
	args := make([]any, 0, len(cols))
	wild := "%" + term + "%"
	for i, col := range cols {
		if i > 0 {
			clause += " OR "
		}
		clause += col + " ILIKE ?"
		args = append(args, wild)
	}
	return clause, args
}
`
}

func apiMaintenanceMiddlewareGo() string {
	return `package middleware

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

// Maintenance returns a middleware that checks for a .maintenance file.
// When the file exists, all requests receive a 503 Service Unavailable response.
// Toggle with: grit down (enable) / grit up (disable)
func Maintenance() gin.HandlerFunc {
	return func(c *gin.Context) {
		if _, err := os.Stat(".maintenance"); err == nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"error": gin.H{
					"code":    "MAINTENANCE",
					"message": "Application is in maintenance mode. Please try again later.",
				},
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
`
}

func apiIdempotencyMiddlewareGo() string {
	return `package middleware

import (
	"bytes"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"` + "{{MODULE}}" + `/internal/cache"
)

// IdempotencyTTL is how long a stored idempotent response is replayed.
// 24h matches Stripe's published behavior and is plenty long for client
// retries while keeping Redis pressure bounded.
const IdempotencyTTL = 24 * time.Hour

// IdempotencyHeader is the header clients set to opt into idempotent retries.
const IdempotencyHeader = "Idempotency-Key"

// Idempotency is a middleware that gives clients safe retry semantics for
// unsafe methods (POST/PUT/PATCH/DELETE). When a request carries an
// Idempotency-Key header, the first successful response (any 2xx) is cached
// and any subsequent request with the same key replays the cached response
// instead of re-executing the handler.
//
// Skipped when:
//   - cacheService is nil (Redis unavailable)
//   - request method is GET/HEAD/OPTIONS (already idempotent)
//   - Idempotency-Key header is missing or empty
//
// Cache key is namespaced per HTTP method + path so the same key reused across
// different endpoints does not collide. The cached payload includes status +
// content type + body, so replay returns a byte-for-byte identical response.
//
// Errors (5xx) are intentionally NOT cached so transient failures can be
// retried with the same key; only 2xx responses are stored.
func Idempotency(cacheService *cache.Cache) gin.HandlerFunc {
	return func(c *gin.Context) {
		if cacheService == nil {
			c.Next()
			return
		}

		switch c.Request.Method {
		case http.MethodGet, http.MethodHead, http.MethodOptions:
			c.Next()
			return
		}

		key := c.GetHeader(IdempotencyHeader)
		if key == "" {
			c.Next()
			return
		}

		cacheKey := "idem:" + c.Request.Method + ":" + c.FullPath() + ":" + key

		// Replay if we've seen this key before.
		var cached idempotentResponse
		found, err := cacheService.Get(c.Request.Context(), cacheKey, &cached)
		if err == nil && found {
			c.Header("Idempotent-Replayed", "true")
			c.Data(cached.Status, cached.ContentType, cached.Body)
			c.Abort()
			return
		}

		// Capture the live response so we can store it after the handler runs.
		writer := &idempotencyCapture{ResponseWriter: c.Writer, buf: bytes.NewBuffer(nil)}
		c.Writer = writer

		c.Next()

		// Only cache 2xx — let clients retry on 4xx/5xx with the same key.
		if writer.status >= 200 && writer.status < 300 {
			resp := idempotentResponse{
				Status:      writer.status,
				ContentType: writer.Header().Get("Content-Type"),
				Body:        writer.buf.Bytes(),
			}
			_ = cacheService.Set(c.Request.Context(), cacheKey, resp, IdempotencyTTL)
		}
	}
}

type idempotentResponse struct {
	Status      int    ` + "`" + `json:"status"` + "`" + `
	ContentType string ` + "`" + `json:"content_type"` + "`" + `
	Body        []byte ` + "`" + `json:"body"` + "`" + `
}

type idempotencyCapture struct {
	gin.ResponseWriter
	buf    *bytes.Buffer
	status int
}

func (w *idempotencyCapture) Write(b []byte) (int, error) {
	w.buf.Write(b)
	return w.ResponseWriter.Write(b)
}

func (w *idempotencyCapture) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}
`
}

func apiSyncRegistryGo() string {
	return `// Package sync owns the model registry used by the offline-first
// /api/sync/push and /api/sync/pull endpoints.
//
// Every model that should be syncable from a desktop client must:
//   1. Have an ID string (UUID) primary key.
//   2. Have a Version int field.
//   3. Have CreatedAt / UpdatedAt timestamps.
//   4. Have a BeforeUpdate hook that increments Version.
//   5. Be registered with Register("table_name", &models.X{}).
//
// The handler uses reflection to decode push payloads into the
// registered struct type, run a versioned update, and detect conflicts
// when the client's version doesn't match what's on disk.
package sync

import (
	"fmt"
	"reflect"
	"sync"
)

// Registry holds the syncable model types keyed by their plural snake_case
// name (e.g. "buildings"). Population happens at app boot from routes.Setup.
type Registry struct {
	mu     sync.RWMutex
	models map[string]reflect.Type
}

// NewRegistry returns an empty Registry.
func NewRegistry() *Registry {
	return &Registry{models: make(map[string]reflect.Type)}
}

// Register adds a model under its plural-snake table name. proto must be
// a pointer to a zero-value struct (e.g. &models.Building{}).
func (r *Registry) Register(table string, proto interface{}) {
	r.mu.Lock()
	defer r.mu.Unlock()
	t := reflect.TypeOf(proto)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	r.models[table] = t
}

// New returns a new pointer to a zero-value model struct for the given
// table, or an error if the table isn't registered.
func (r *Registry) New(table string) (interface{}, error) {
	r.mu.RLock()
	t, ok := r.models[table]
	r.mu.RUnlock()
	if !ok {
		return nil, fmt.Errorf("sync: unknown table %q", table)
	}
	return reflect.New(t).Interface(), nil
}

// Tables lists every registered table name. Used by /api/sync/pull when
// the client asks for the full set of types.
func (r *Registry) Tables() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]string, 0, len(r.models))
	for k := range r.models {
		out = append(out, k)
	}
	return out
}
`
}

func apiSyncHandlerGo() string {
	return `package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"` + "{{MODULE}}" + `/internal/sync"
)

// SyncHandler implements /api/sync/push and /api/sync/pull. The push
// endpoint applies a batch of client changes with per-change version
// checking; the pull endpoint streams server-side updates since a
// caller-supplied cursor.
type SyncHandler struct {
	DB       *gorm.DB
	Registry *sync.Registry
}

// NewSyncHandler wires the handler to the database + model registry.
func NewSyncHandler(db *gorm.DB, reg *sync.Registry) *SyncHandler {
	return &SyncHandler{DB: db, Registry: reg}
}

// PushChange is one entry in a /api/sync/push batch. Op is one of
// "create" / "update" / "delete". Version is the version the client
// believes the server has — mismatches surface as VERSION_CONFLICT.
type PushChange struct {
	Op      string                 ` + "`" + `json:"op"` + "`" + `
	Model   string                 ` + "`" + `json:"model"` + "`" + `
	ID      string                 ` + "`" + `json:"id"` + "`" + `
	Version int                    ` + "`" + `json:"version"` + "`" + `
	Data    map[string]interface{} ` + "`" + `json:"data"` + "`" + `
}

// PushResult is the per-change result returned in the same order as
// the input batch. On VERSION_CONFLICT, ServerVersion + ServerData
// carry the current server state so the client can build a merge UI.
type PushResult struct {
	OK            bool        ` + "`" + `json:"ok"` + "`" + `
	Code          string      ` + "`" + `json:"code,omitempty"` + "`" + `
	Message       string      ` + "`" + `json:"message,omitempty"` + "`" + `
	ServerVersion int         ` + "`" + `json:"server_version,omitempty"` + "`" + `
	ServerData    interface{} ` + "`" + `json:"server_data,omitempty"` + "`" + `
	NewVersion    int         ` + "`" + `json:"new_version,omitempty"` + "`" + `
}

// Push handles POST /api/sync/push. Each change is applied
// independently — one conflict does not abort the rest of the batch.
func (h *SyncHandler) Push(c *gin.Context) {
	var req struct {
		Changes []PushChange ` + "`" + `json:"changes"` + "`" + `
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "INVALID_BODY", "message": err.Error()}})
		return
	}

	results := make([]PushResult, len(req.Changes))
	for i, ch := range req.Changes {
		results[i] = h.applyChange(ch)
	}
	c.JSON(http.StatusOK, gin.H{"results": results})
}

func (h *SyncHandler) applyChange(ch PushChange) PushResult {
	proto, err := h.Registry.New(ch.Model)
	if err != nil {
		return PushResult{OK: false, Code: "UNKNOWN_MODEL", Message: err.Error()}
	}

	switch ch.Op {
	case "create":
		// Decode the client payload into a fresh model struct and insert.
		// We trust the client-supplied ID (UUID) so the local outbox can
		// keep referring to the same row after the server insert.
		obj := proto
		if err := decodeInto(obj, ch.Data); err != nil {
			return PushResult{OK: false, Code: "DECODE_ERROR", Message: err.Error()}
		}
		setField(obj, "ID", ch.ID)
		if err := h.DB.Create(obj).Error; err != nil {
			return PushResult{OK: false, Code: "CREATE_FAILED", Message: err.Error()}
		}
		return PushResult{OK: true, NewVersion: 1}

	case "update":
		// Versioned update: load current row, compare versions, update if match.
		current := proto
		if err := h.DB.First(current, "id = ?", ch.ID).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return PushResult{OK: false, Code: "NOT_FOUND", Message: "row was deleted on the server"}
			}
			return PushResult{OK: false, Code: "INTERNAL_ERROR", Message: err.Error()}
		}
		serverVersion := getIntField(current, "Version")
		if serverVersion != ch.Version {
			return PushResult{
				OK:            false,
				Code:          "VERSION_CONFLICT",
				Message:       fmt.Sprintf("client had v%d, server has v%d", ch.Version, serverVersion),
				ServerVersion: serverVersion,
				ServerData:    current,
			}
		}
		// Versions match — apply the update. The BeforeUpdate hook will
		// bump Version on save so the client knows what to remember.
		if err := h.DB.Model(current).Updates(ch.Data).Error; err != nil {
			return PushResult{OK: false, Code: "UPDATE_FAILED", Message: err.Error()}
		}
		return PushResult{OK: true, NewVersion: serverVersion + 1}

	case "delete":
		current := proto
		if err := h.DB.First(current, "id = ?", ch.ID).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				// Already gone — treat as success so the outbox can clear.
				return PushResult{OK: true}
			}
			return PushResult{OK: false, Code: "INTERNAL_ERROR", Message: err.Error()}
		}
		serverVersion := getIntField(current, "Version")
		if ch.Version != 0 && serverVersion != ch.Version {
			return PushResult{
				OK:            false,
				Code:          "VERSION_CONFLICT",
				Message:       "row was modified after the client's last sync",
				ServerVersion: serverVersion,
				ServerData:    current,
			}
		}
		if err := h.DB.Delete(current, "id = ?", ch.ID).Error; err != nil {
			return PushResult{OK: false, Code: "DELETE_FAILED", Message: err.Error()}
		}
		return PushResult{OK: true}

	default:
		return PushResult{OK: false, Code: "INVALID_OP", Message: "op must be create, update, or delete"}
	}
}

// Pull handles GET /api/sync/pull?since=<rfc3339>&model=<table>. Returns
// every row in the requested table with UpdatedAt > since. The client
// uses the response's cursor as the next ?since value.
func (h *SyncHandler) Pull(c *gin.Context) {
	model := c.Query("model")
	if model == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "MISSING_MODEL", "message": "?model is required"}})
		return
	}
	sinceStr := c.DefaultQuery("since", "")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "500"))
	if limit < 1 || limit > 5000 {
		limit = 500
	}

	proto, err := h.Registry.New(model)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "UNKNOWN_MODEL", "message": err.Error()}})
		return
	}

	// Build a slice of the right type via reflection so we can return
	// proper struct values (not gin.H maps).
	sliceType := reflect.SliceOf(reflect.TypeOf(proto).Elem())
	results := reflect.New(sliceType)

	q := h.DB.Model(proto)
	if sinceStr != "" {
		t, err := time.Parse(time.RFC3339Nano, sinceStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "INVALID_SINCE", "message": err.Error()}})
			return
		}
		q = q.Where("updated_at > ?", t)
	}
	if err := q.Order("updated_at asc").Limit(limit).Find(results.Interface()).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"code": "INTERNAL_ERROR", "message": err.Error()}})
		return
	}

	// Cursor = last UpdatedAt in the page so the client picks up there
	// next time. Empty when nothing came back.
	cursor := sinceStr
	rs := results.Elem()
	if rs.Len() > 0 {
		last := rs.Index(rs.Len() - 1).Addr().Interface()
		if t, ok := getTimeField(last, "UpdatedAt"); ok {
			cursor = t.Format(time.RFC3339Nano)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"data":   results.Elem().Interface(),
		"cursor": cursor,
		"count":  rs.Len(),
	})
}

// decodeInto round-trips a map through JSON into the target struct so
// gorm field tags + types are respected. Cheap; the maps are small.
func decodeInto(target interface{}, data map[string]interface{}) error {
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, target)
}

// setField sets a string field on a struct via reflection. Used for ID.
func setField(obj interface{}, name, value string) {
	v := reflect.ValueOf(obj).Elem()
	f := v.FieldByName(name)
	if f.IsValid() && f.CanSet() && f.Kind() == reflect.String {
		f.SetString(value)
	}
}

func getIntField(obj interface{}, name string) int {
	v := reflect.ValueOf(obj)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	f := v.FieldByName(name)
	if !f.IsValid() {
		return 0
	}
	return int(f.Int())
}

func getTimeField(obj interface{}, name string) (time.Time, bool) {
	v := reflect.ValueOf(obj)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	f := v.FieldByName(name)
	if !f.IsValid() {
		return time.Time{}, false
	}
	t, ok := f.Interface().(time.Time)
	return t, ok
}
`
}

func apiActivityLogModelGo() string {
	return `package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ActivityLog records every successful authenticated mutation, with a
// tamper-evident hash chain — each row's Hash is SHA-256 of (PrevHash
// || canonical(this_row)). Mutating any row breaks the chain on the
// next VerifyChain pass.
//
// The payload digest is a SHA-256 of the request body so we have
// evidence of what was sent without storing PII verbatim. Read-only —
// no updates, no deletes (use a separate retention job to prune old
// rows; deletion still breaks the chain so it must rebuild from a
// safe checkpoint).
type ActivityLog struct {
	ID            string    ` + "`" + `gorm:"primarykey;size:36" json:"id"` + "`" + `
	UserID        string    ` + "`" + `gorm:"size:36;index" json:"user_id"` + "`" + `
	Method        string    ` + "`" + `gorm:"size:10" json:"method"` + "`" + `
	Path          string    ` + "`" + `gorm:"size:500;index" json:"path"` + "`" + `
	Status        int       ` + "`" + `json:"status"` + "`" + `
	PayloadDigest string    ` + "`" + `gorm:"size:64" json:"payload_digest"` + "`" + ` // sha256 hex
	IPAddress     string    ` + "`" + `gorm:"size:45" json:"ip_address"` + "`" + `
	UserAgent     string    ` + "`" + `gorm:"size:500" json:"user_agent"` + "`" + `
	DurationMS    int64     ` + "`" + `json:"duration_ms"` + "`" + `
	PrevHash      string    ` + "`" + `gorm:"size:64" json:"prev_hash"` + "`" + ` // hex sha256, "" for the genesis row
	Hash          string    ` + "`" + `gorm:"size:64;uniqueIndex" json:"hash"` + "`" + ` // hex sha256(prev_hash || canonical)
	CreatedAt     time.Time ` + "`" + `gorm:"index" json:"created_at"` + "`" + `
}

func (a *ActivityLog) BeforeCreate(tx *gorm.DB) error {
	if a.ID == "" {
		a.ID = uuid.New().String()
	}
	return nil
}
`
}

func apiAuditGo() string {
	return `// Package audit owns the tamper-evident hash chain over the activity log.
//
// Each row's Hash = SHA-256(PrevHash || canonical(row)) where canonical
// is a stable JSON serialization of the audit-relevant fields. Any
// mutation to a row breaks every Hash from that row forward, which
// VerifyChain detects.
//
// Insert is serialized via a row-level FOR UPDATE lock on the latest
// row inside the same transaction that does the INSERT — concurrent
// inserts queue cleanly without forking the chain. Verification walks
// the chain in created_at + id order; ties broken by id.
//
// What this defends against:
//   - Direct SQL UPDATE / DELETE on activity_logs (most common attack
//     vector — DBA covering tracks).
//   - Out-of-band insertion of forged history.
//
// What this does NOT defend against:
//   - Compromise of the running server itself (an attacker with code
//     execution can rewrite the whole chain). External anchoring
//     (publishing the daily root hash to a public ledger) is the
//     follow-up — see #48.
package audit

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"` + "{{MODULE}}" + `/internal/models"
)

// Canonical returns the stable JSON bytes of an entry for hashing.
// We exclude ID / PrevHash / Hash from the canonical form: ID is
// random and uncorrelated with content; PrevHash + Hash are derived
// values, not inputs to the hash.
func Canonical(e *models.ActivityLog) ([]byte, error) {
	c := canonicalEntry{
		UserID:        e.UserID,
		Method:        e.Method,
		Path:          e.Path,
		Status:        e.Status,
		PayloadDigest: e.PayloadDigest,
		IPAddress:     e.IPAddress,
		UserAgent:     e.UserAgent,
		DurationMS:    e.DurationMS,
		// Use unix-nano so the canonical bytes are stable across tz
		// changes / TIMESTAMPTZ formatting differences.
		CreatedAtUnixNano: e.CreatedAt.UTC().UnixNano(),
	}
	return json.Marshal(c)
}

// canonicalEntry's field order is the wire format for hashing —
// reorder ONLY in a major version bump (verify breaks otherwise).
type canonicalEntry struct {
	UserID            string ` + "`" + `json:"user_id"` + "`" + `
	Method            string ` + "`" + `json:"method"` + "`" + `
	Path              string ` + "`" + `json:"path"` + "`" + `
	Status            int    ` + "`" + `json:"status"` + "`" + `
	PayloadDigest     string ` + "`" + `json:"payload_digest"` + "`" + `
	IPAddress         string ` + "`" + `json:"ip_address"` + "`" + `
	UserAgent         string ` + "`" + `json:"user_agent"` + "`" + `
	DurationMS        int64  ` + "`" + `json:"duration_ms"` + "`" + `
	CreatedAtUnixNano int64  ` + "`" + `json:"created_at_unix_nano"` + "`" + `
}

// ComputeHash returns hex(sha256(prevHash || canonical)) — the prev
// hash is included as a hex string (not raw bytes) so the input is
// trivially auditable: cat prev_hash | xxd; cat canonical.json.
func ComputeHash(prevHash string, canonical []byte) string {
	h := sha256.New()
	h.Write([]byte(prevHash))
	h.Write(canonical)
	return hex.EncodeToString(h.Sum(nil))
}

// AppendChained inserts a new ActivityLog with PrevHash + Hash filled
// in. Intended for ad-hoc / one-off audit writes from app code (NOT
// the hot-path middleware — that uses the buffered worker pattern).
//
// Concurrency note: this function takes a row-level FOR UPDATE lock
// on the latest row to serialize concurrent callers. Use sparingly;
// for any high-throughput audit source, route through the middleware's
// channel writer instead.
func AppendChained(db *gorm.DB, entry *models.ActivityLog) error {
	return db.Transaction(func(tx *gorm.DB) error {
		var prev models.ActivityLog
		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Order("created_at desc, id desc").
			Limit(1).
			First(&prev).Error
		if err != nil && err != gorm.ErrRecordNotFound {
			return err
		}

		canonical, err := Canonical(entry)
		if err != nil {
			return fmt.Errorf("canonicalize: %w", err)
		}
		entry.PrevHash = prev.Hash
		entry.Hash = ComputeHash(prev.Hash, canonical)
		return tx.Create(entry).Error
	})
}

// ChainStatus is the result of VerifyChain.
type ChainStatus struct {
	Valid        bool   ` + "`" + `json:"valid"` + "`" + `
	TotalEntries int    ` + "`" + `json:"total_entries"` + "`" + `
	BrokenAtID   string ` + "`" + `json:"broken_at_id,omitempty"` + "`" + `
	BrokenAt     int    ` + "`" + `json:"broken_at,omitempty"` + "`" + ` // zero-indexed position
	Expected     string ` + "`" + `json:"expected,omitempty"` + "`" + `
	Got          string ` + "`" + `json:"got,omitempty"` + "`" + `
	Message      string ` + "`" + `json:"message,omitempty"` + "`" + `
}

// VerifyChain walks the entire activity log in (created_at, id) order
// and recomputes every Hash. The first mismatch is reported with the
// position and offending row's ID — everything before that position
// is trustworthy.
//
// Memory-bounded: iterates in batches of verifyBatchSize so a 100M-row
// log doesn't OOM the process. Honours context cancellation so the
// caller can attach a deadline (the admin endpoint should pass
// c.Request.Context() with a 30s timeout).
//
// Cost is O(n) — about a second per million rows on a warm cache.
// Wire to a nightly cron + a /api/admin/activity/integrity endpoint.
const verifyBatchSize = 1000

func VerifyChain(ctx context.Context, db *gorm.DB) (ChainStatus, error) {
	prevHash := ""
	total := 0
	var lastCreatedAt time.Time
	var lastID string

	for {
		select {
		case <-ctx.Done():
			return ChainStatus{TotalEntries: total}, ctx.Err()
		default:
		}

		var batch []models.ActivityLog
		q := db.Order("created_at asc, id asc").Limit(verifyBatchSize)
		if total > 0 {
			// Cursor on (created_at, id) so we don't re-read rows
			// already verified in the previous batch.
			q = q.Where("(created_at, id) > (?, ?)", lastCreatedAt, lastID)
		}
		if err := q.Find(&batch).Error; err != nil {
			return ChainStatus{TotalEntries: total}, err
		}
		if len(batch) == 0 {
			break
		}

		for i := range batch {
			e := &batch[i]
			canonical, err := Canonical(e)
			if err != nil {
				return ChainStatus{TotalEntries: total}, err
			}
			expected := ComputeHash(prevHash, canonical)
			if expected != e.Hash {
				return ChainStatus{
					Valid:        false,
					TotalEntries: total + i,
					BrokenAtID:   e.ID,
					BrokenAt:     total + i,
					Expected:     expected,
					Got:          e.Hash,
					Message:      "hash mismatch — row was modified, deleted, or inserted out of order",
				}, nil
			}
			if e.PrevHash != prevHash {
				return ChainStatus{
					Valid:        false,
					TotalEntries: total + i,
					BrokenAtID:   e.ID,
					BrokenAt:     total + i,
					Expected:     prevHash,
					Got:          e.PrevHash,
					Message:      "prev_hash mismatch — chain link broken",
				}, nil
			}
			prevHash = e.Hash
		}

		last := &batch[len(batch)-1]
		lastCreatedAt = last.CreatedAt
		lastID = last.ID
		total += len(batch)

		if len(batch) < verifyBatchSize {
			break // last page
		}
	}

	return ChainStatus{
		Valid:        true,
		TotalEntries: total,
	}, nil
}
`
}

func apiWebhookEventModelGo() string {
	return `package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// WebhookEvent persists every webhook the API receives. ExternalID is
// the provider's own event ID — we use it as the idempotency key, so
// duplicate deliveries (Stripe retries, partner pings) become no-ops.
//
// Status lifecycle:
//   pending   — received + verified, handler not yet run
//   processed — handler returned nil
//   failed    — handler returned an error; HandlerError holds the message
//   skipped   — duplicate ExternalID — handler was bypassed
type WebhookEvent struct {
	ID           string         ` + "`" + `gorm:"primarykey;size:36" json:"id"` + "`" + `
	Provider     string         ` + "`" + `gorm:"size:50;index;not null" json:"provider"` + "`" + `
	EventType    string         ` + "`" + `gorm:"size:100;index" json:"event_type"` + "`" + `
	ExternalID   string         ` + "`" + `gorm:"size:255;index" json:"external_id"` + "`" + ` // provider's event id
	Payload      datatypes.JSON ` + "`" + `gorm:"type:jsonb" json:"payload"` + "`" + `
	Status       string         ` + "`" + `gorm:"size:20;index;not null;default:pending" json:"status"` + "`" + `
	HandlerError string         ` + "`" + `gorm:"type:text" json:"handler_error,omitempty"` + "`" + `
	RetryCount   int            ` + "`" + `gorm:"not null;default:0" json:"retry_count"` + "`" + `
	ProcessedAt  *time.Time     ` + "`" + `json:"processed_at,omitempty"` + "`" + `
	CreatedAt    time.Time      ` + "`" + `gorm:"index" json:"created_at"` + "`" + `
}

func (w *WebhookEvent) BeforeCreate(tx *gorm.DB) error {
	if w.ID == "" {
		w.ID = uuid.New().String()
	}
	return nil
}

// Composite unique index on (provider, external_id) gives us
// idempotent receipt: a duplicate delivery from the same provider
// with the same event id fails the INSERT, which we treat as
// "already processed".
func (WebhookEvent) Indexes() string {
	return "CREATE UNIQUE INDEX IF NOT EXISTS idx_webhook_events_provider_external_id ON webhook_events(provider, external_id) WHERE external_id <> ''"
}
`
}

func apiWebhooksGo() string {
	return `// Package webhooks is the receive-side framework for inbound
// webhooks (Stripe, GitHub, WhatsApp, Twilio, Slack, anything that
// pings you). The shape:
//
//   webhooks.Register("stripe", webhooks.Provider{
//       SecretEnv: "STRIPE_WEBHOOK_SECRET",
//       Verify:    webhooks.StripeVerifier,
//       Extract:   webhooks.StripeExtractor,
//   })
//
//   webhooks.On("stripe", "invoice.paid", func(ctx context.Context, e *models.WebhookEvent) error {
//       // process the event…
//       return nil
//   })
//
// At app boot, call webhooks.Setup(db) once. The HTTP handler is
// already wired in routes.go at POST /webhooks/:provider — it does:
//   1. Look up the provider config (404 if unknown)
//   2. Read raw body + headers
//   3. Verify signature via Provider.Verify
//   4. Extract event type + external id via Provider.Extract
//   5. INSERT into webhook_events (unique on provider+external_id —
//      duplicate delivery becomes a no-op, status=skipped)
//   6. Run the registered handler for (provider, event_type)
//   7. Update event row with processed/failed status
//
// Failed handlers stay in the table with status=failed; the admin
// endpoint POST /api/admin/webhooks/:id/replay re-runs the handler.
package webhooks

import (
	"context"
	"fmt"
	"sync"

	"gorm.io/gorm"

	"` + "{{MODULE}}" + `/internal/models"
)

// VerifyFunc validates a request's signature. Returns an error if the
// payload was tampered with or the signature is missing/invalid.
type VerifyFunc func(secret string, body []byte, headers map[string]string) error

// ExtractFunc pulls (eventType, externalID) from a verified payload.
// EventType drives handler dispatch; ExternalID drives idempotency.
type ExtractFunc func(body []byte, headers map[string]string) (eventType string, externalID string, err error)

// Handler is the user-defined function that processes a verified +
// deduplicated webhook. Errors are persisted to webhook_events.handler_error.
type Handler func(ctx context.Context, e *models.WebhookEvent) error

// Provider is the per-source configuration.
type Provider struct {
	SecretEnv string      // env var holding the signing secret
	Verify    VerifyFunc  // signature verifier (StripeVerifier, GitHubVerifier, HMACVerifier, etc.)
	Extract   ExtractFunc // event type + external id extractor
}

var (
	mu        sync.RWMutex
	providers = map[string]Provider{}
	handlers  = map[string]map[string]Handler{} // provider → eventType → handler
	db        *gorm.DB
)

// Setup wires the package to the project's *gorm.DB. Call once at app
// boot from routes.Setup or main.
func Setup(database *gorm.DB) {
	mu.Lock()
	defer mu.Unlock()
	db = database
}

// Register adds a provider configuration. Call from package init() or
// from a setup function in your handlers package.
func Register(name string, p Provider) {
	mu.Lock()
	defer mu.Unlock()
	providers[name] = p
	if _, ok := handlers[name]; !ok {
		handlers[name] = map[string]Handler{}
	}
}

// On binds a handler to (provider, eventType). Use the empty string
// "" as eventType to register a catch-all handler — it runs for any
// event from this provider that doesn't have a specific handler.
func On(provider, eventType string, h Handler) {
	mu.Lock()
	defer mu.Unlock()
	if _, ok := handlers[provider]; !ok {
		handlers[provider] = map[string]Handler{}
	}
	handlers[provider][eventType] = h
}

// LookupProvider returns the Provider config for name.
func LookupProvider(name string) (Provider, bool) {
	mu.RLock()
	defer mu.RUnlock()
	p, ok := providers[name]
	return p, ok
}

// Dispatch finds a handler for (provider, eventType). Falls back to
// the catch-all "" handler if no specific match. Returns nil if no
// handler is registered (the event is still persisted, just unprocessed).
func Dispatch(ctx context.Context, e *models.WebhookEvent) error {
	mu.RLock()
	pmap, ok := handlers[e.Provider]
	mu.RUnlock()
	if !ok {
		return nil
	}
	mu.RLock()
	h, exact := pmap[e.EventType]
	if !exact {
		h = pmap[""] // catch-all
	}
	mu.RUnlock()
	if h == nil {
		return nil
	}
	return h(ctx, e)
}

// DB returns the registered *gorm.DB or nil if Setup hasn't been called.
// Used by the HTTP handler — exposed so admin endpoints can re-use it.
func DB() *gorm.DB {
	mu.RLock()
	defer mu.RUnlock()
	return db
}

// IsDuplicateError reports whether err looks like a unique-constraint
// violation on (provider, external_id). Postgres + SQLite both surface
// these distinctly, but the message format varies — check substrings.
func IsDuplicateError(err error) bool {
	if err == nil {
		return false
	}
	s := err.Error()
	return contains(s, "duplicate key") ||
		contains(s, "UNIQUE constraint") ||
		contains(s, "duplicate entry")
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && (indexOf(s, sub) >= 0)
}

func indexOf(s, sub string) int {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return i
		}
	}
	return -1
}

// MissingProviderError is returned when an unregistered provider is hit.
type MissingProviderError struct{ Name string }

func (e MissingProviderError) Error() string {
	return fmt.Sprintf("webhooks: provider %q not registered", e.Name)
}
`
}

func apiWebhooksVerifiersGo() string {
	return `package webhooks

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// HMACVerifier returns a VerifyFunc that validates a hex-encoded
// HMAC-SHA256 signature found in the named header. Most simple
// providers (custom partners, self-rolled webhooks) use this scheme.
//
//	webhooks.Register("partner", webhooks.Provider{
//	    SecretEnv: "PARTNER_WEBHOOK_SECRET",
//	    Verify:    webhooks.HMACVerifier("X-Signature"),
//	    Extract:   webhooks.JSONFieldExtractor("type", "id"),
//	})
func HMACVerifier(header string) VerifyFunc {
	return func(secret string, body []byte, headers map[string]string) error {
		if secret == "" {
			return fmt.Errorf("webhooks: signing secret is empty")
		}
		got := headers[header]
		if got == "" {
			got = headers[strings.ToLower(header)]
		}
		if got == "" {
			return fmt.Errorf("webhooks: missing signature header %q", header)
		}
		mac := hmac.New(sha256.New, []byte(secret))
		mac.Write(body)
		expected := hex.EncodeToString(mac.Sum(nil))
		if !hmac.Equal([]byte(got), []byte(expected)) {
			return fmt.Errorf("webhooks: signature mismatch")
		}
		return nil
	}
}

// StripeVerifier validates Stripe's "Stripe-Signature" header, which
// has the form "t=<unix>,v1=<hex>" where v1 = HMAC-SHA256 of
// "<timestamp>.<payload>" using the webhook signing secret. Tolerance
// of 5 minutes guards against replay.
//
// See https://stripe.com/docs/webhooks/signatures
func StripeVerifier(secret string, body []byte, headers map[string]string) error {
	const tolerance = 5 * time.Minute
	if secret == "" {
		return fmt.Errorf("webhooks: stripe secret is empty")
	}
	header := headers["Stripe-Signature"]
	if header == "" {
		header = headers["stripe-signature"]
	}
	if header == "" {
		return fmt.Errorf("webhooks: missing Stripe-Signature header")
	}

	var ts int64
	var sigs []string
	for _, part := range strings.Split(header, ",") {
		kv := strings.SplitN(strings.TrimSpace(part), "=", 2)
		if len(kv) != 2 {
			continue
		}
		switch kv[0] {
		case "t":
			ts, _ = strconv.ParseInt(kv[1], 10, 64)
		case "v1":
			sigs = append(sigs, kv[1])
		}
	}
	if ts == 0 || len(sigs) == 0 {
		return fmt.Errorf("webhooks: malformed Stripe-Signature header")
	}
	if time.Since(time.Unix(ts, 0)) > tolerance {
		return fmt.Errorf("webhooks: stripe timestamp outside tolerance")
	}

	signed := strconv.FormatInt(ts, 10) + "." + string(body)
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(signed))
	expected := hex.EncodeToString(mac.Sum(nil))
	for _, s := range sigs {
		if hmac.Equal([]byte(s), []byte(expected)) {
			return nil
		}
	}
	return fmt.Errorf("webhooks: stripe signature mismatch")
}

// GitHubVerifier validates GitHub's "X-Hub-Signature-256" header,
// which is "sha256=<hex>" — HMAC-SHA256 of the raw body using the
// webhook secret.
func GitHubVerifier(secret string, body []byte, headers map[string]string) error {
	if secret == "" {
		return fmt.Errorf("webhooks: github secret is empty")
	}
	header := headers["X-Hub-Signature-256"]
	if header == "" {
		header = headers["x-hub-signature-256"]
	}
	if header == "" {
		return fmt.Errorf("webhooks: missing X-Hub-Signature-256 header")
	}
	prefix := "sha256="
	if !strings.HasPrefix(header, prefix) {
		return fmt.Errorf("webhooks: unexpected X-Hub-Signature-256 format")
	}
	got := header[len(prefix):]
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(body)
	expected := hex.EncodeToString(mac.Sum(nil))
	if !hmac.Equal([]byte(got), []byte(expected)) {
		return fmt.Errorf("webhooks: github signature mismatch")
	}
	return nil
}

// JSONFieldExtractor returns an ExtractFunc that pulls type + id from
// top-level JSON fields in the body. Stripe-style payloads use
// JSONFieldExtractor("type", "id") — the most common shape.
func JSONFieldExtractor(typeField, idField string) ExtractFunc {
	return func(body []byte, headers map[string]string) (string, string, error) {
		var raw map[string]interface{}
		if err := json.Unmarshal(body, &raw); err != nil {
			return "", "", fmt.Errorf("decoding payload: %w", err)
		}
		t, _ := raw[typeField].(string)
		id, _ := raw[idField].(string)
		return t, id, nil
	}
}

// StripeExtractor pulls (type, id) from Stripe's standard
// { "type": "...", "id": "evt_..." } envelope.
var StripeExtractor = JSONFieldExtractor("type", "id")

// GitHubExtractor reads the event type from the "X-GitHub-Event"
// header and the delivery ID from "X-GitHub-Delivery".
func GitHubExtractor(body []byte, headers map[string]string) (string, string, error) {
	t := headers["X-GitHub-Event"]
	if t == "" {
		t = headers["x-github-event"]
	}
	id := headers["X-GitHub-Delivery"]
	if id == "" {
		id = headers["x-github-delivery"]
	}
	return t, id, nil
}
`
}

func apiWebhooksHandlerGo() string {
	return `package handlers

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/datatypes"
	"gorm.io/gorm"

	"` + "{{MODULE}}" + `/internal/models"
	"` + "{{MODULE}}" + `/internal/paginate"
	"` + "{{MODULE}}" + `/internal/webhooks"
)

// WebhookHandler is the universal entry point for inbound webhooks.
// One handler, one route — POST /webhooks/:provider routes by the
// provider path param and dispatches to whatever was registered.
type WebhookHandler struct {
	DB *gorm.DB
}

func NewWebhookHandler(db *gorm.DB) *WebhookHandler {
	return &WebhookHandler{DB: db}
}

// Receive is mounted at POST /webhooks/:provider. It:
//
//  1. Looks up the provider in the registry (404 if unknown).
//  2. Reads the raw body + collects headers.
//  3. Calls Provider.Verify — 401 on signature mismatch.
//  4. Calls Provider.Extract to get (event_type, external_id).
//  5. Inserts a WebhookEvent (unique on provider+external_id — a
//     duplicate becomes status=skipped and we 200 immediately).
//  6. Calls webhooks.Dispatch in the request context.
//  7. Updates status=processed or status=failed with HandlerError.
//
// Always returns 200 to the provider on a verified+stored event so
// they don't retry forever — handler failures are surfaced via the
// admin replay endpoint.
func (h *WebhookHandler) Receive(c *gin.Context) {
	providerName := c.Param("provider")
	provider, ok := webhooks.LookupProvider(providerName)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{
			"error": gin.H{"code": "UNKNOWN_PROVIDER", "message": "no webhook provider registered for " + providerName},
		})
		return
	}

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{"code": "READ_BODY_FAILED", "message": err.Error()},
		})
		return
	}

	headers := flattenHeaders(c.Request.Header)
	secret := os.Getenv(provider.SecretEnv)
	if err := provider.Verify(secret, body, headers); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": gin.H{"code": "INVALID_SIGNATURE", "message": err.Error()},
		})
		return
	}

	eventType, externalID, err := provider.Extract(body, headers)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{"code": "EXTRACT_FAILED", "message": err.Error()},
		})
		return
	}

	event := models.WebhookEvent{
		Provider:   providerName,
		EventType:  eventType,
		ExternalID: externalID,
		Payload:    datatypes.JSON(body),
		Status:     "pending",
	}
	if err := h.DB.Create(&event).Error; err != nil {
		// Duplicate (provider, external_id) — already processed.
		// Return 200 so the provider doesn't retry, and skip the handler.
		if webhooks.IsDuplicateError(err) {
			c.JSON(http.StatusOK, gin.H{"status": "skipped", "reason": "duplicate"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{"code": "PERSIST_FAILED", "message": err.Error()},
		})
		return
	}

	// Dispatch in the request context so handlers can attach DB
	// timeouts / cancellation. Failures are recorded but never bubble
	// up — the provider already got a 200 once we persisted.
	if dispatchErr := webhooks.Dispatch(c.Request.Context(), &event); dispatchErr != nil {
		now := time.Now()
		h.DB.Model(&event).Updates(map[string]interface{}{
			"status":        "failed",
			"handler_error": dispatchErr.Error(),
			"processed_at":  &now,
		})
		c.JSON(http.StatusOK, gin.H{"status": "received", "id": event.ID, "handler": "failed"})
		return
	}
	now := time.Now()
	h.DB.Model(&event).Updates(map[string]interface{}{
		"status":       "processed",
		"processed_at": &now,
	})
	c.JSON(http.StatusOK, gin.H{"status": "processed", "id": event.ID})
}

// List returns the recent webhook events with the standard paginate envelope.
//
//	GET /api/admin/webhooks?provider=stripe&status=failed
func (h *WebhookHandler) List(c *gin.Context) {
	q := h.DB.Model(&models.WebhookEvent{})
	params := paginate.Bind(c).
		With("provider", c.Query("provider")).
		With("status", c.Query("status"))

	res, err := paginate.List[models.WebhookEvent](q, params, paginate.Config{
		Sortable:     map[string]bool{"created_at": true, "status": true, "provider": true, "event_type": true},
		DefaultSort:  "created_at",
		DefaultOrder: "desc",
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{"code": "INTERNAL_ERROR", "message": err.Error()},
		})
		return
	}
	c.JSON(http.StatusOK, res)
}

// Replay re-runs the handler for an existing webhook event. Used to
// recover from a transient handler failure or a deploy that fixed a
// bug. Increments retry_count + records the new outcome.
//
//	POST /api/admin/webhooks/:id/replay
func (h *WebhookHandler) Replay(c *gin.Context) {
	id := c.Param("id")
	var event models.WebhookEvent
	if err := h.DB.First(&event, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": gin.H{"code": "NOT_FOUND", "message": "webhook event not found"},
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{"code": "INTERNAL_ERROR", "message": err.Error()},
		})
		return
	}

	dispatchErr := webhooks.Dispatch(c.Request.Context(), &event)
	now := time.Now()
	// Atomic retry_count increment via gorm.Expr — two concurrent
	// replays of the same event are safe (each adds 1 instead of
	// both reading the same baseline and writing the same +1 result).
	updates := map[string]interface{}{
		"retry_count":  gorm.Expr("retry_count + ?", 1),
		"processed_at": &now,
	}
	if dispatchErr == nil {
		updates["status"] = "processed"
		updates["handler_error"] = ""
	} else {
		updates["status"] = "failed"
		updates["handler_error"] = dispatchErr.Error()
	}
	if err := h.DB.Model(&event).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{"code": "INTERNAL_ERROR", "message": err.Error()},
		})
		return
	}
	// Re-read to get the post-increment count; the original event.RetryCount
	// is stale after the gorm.Expr update.
	_ = h.DB.Select("retry_count").First(&event, "id = ?", id).Error
	if dispatchErr != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":        "failed",
			"handler_error": dispatchErr.Error(),
			"retry_count":   event.RetryCount,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "processed", "retry_count": event.RetryCount})
}

// flattenHeaders turns http.Header (multi-value) into a single-value
// map for the Verify / Extract callbacks. Keeps the framework API
// simple — nobody needs the multi-value form for webhook signing.
func flattenHeaders(h http.Header) map[string]string {
	out := make(map[string]string, len(h))
	for k, v := range h {
		if len(v) > 0 {
			out[k] = v[0]
		}
	}
	return out
}

// Dispatch is exposed so app code can fire a synthetic event in tests.
var _ = context.Background
var _ = fmt.Sprint
`
}

func apiFeatureFlagModelGo() string {
	return `package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// FeatureFlag is one rollout switch. Two flavors:
//   - Boolean (Variants empty) → IsEnabled returns true/false
//   - A/B (Variants set)       → Variant returns one of the listed
//                                strings, sticky per (user, flag).
//
// Rules JSON shape (FlagRules): rollout_percentage, allowlist_user_ids,
// blocklist_user_ids, enabled_from, enabled_until, variants. The
// percentage and variant assignment both bucket users by
// SHA-256(user_id || ":" || flag_name) % 100 so the same user always
// lands in the same slot for a given flag — no flicker between sessions.
type FeatureFlag struct {
	ID          string         ` + "`" + `gorm:"primarykey;size:36" json:"id"` + "`" + `
	Name        string         ` + "`" + `gorm:"size:100;uniqueIndex;not null" json:"name"` + "`" + ` // e.g. "new_dashboard"
	Description string         ` + "`" + `gorm:"type:text" json:"description"` + "`" + `
	Enabled     bool           ` + "`" + `gorm:"not null;default:false" json:"enabled"` + "`" + ` // master switch — false short-circuits all rules
	Rules       datatypes.JSON ` + "`" + `gorm:"type:jsonb" json:"rules"` + "`" + `
	CreatedAt   time.Time      ` + "`" + `json:"created_at"` + "`" + `
	UpdatedAt   time.Time      ` + "`" + `json:"updated_at"` + "`" + `
	Version     int            ` + "`" + `gorm:"not null;default:1" json:"version"` + "`" + `
}

func (f *FeatureFlag) BeforeCreate(tx *gorm.DB) error {
	if f.ID == "" {
		f.ID = uuid.New().String()
	}
	return nil
}

func (f *FeatureFlag) BeforeUpdate(tx *gorm.DB) error {
	f.Version++
	return nil
}

// FlagRules is the structured form of FeatureFlag.Rules. Use
// (*FeatureFlag).ParsedRules() to decode; (*FeatureFlag).SetRules() to
// encode + assign.
type FlagRules struct {
	RolloutPercentage int        ` + "`" + `json:"rollout_percentage,omitempty"` + "`" + ` // 0..100; 0 = off, 100 = full rollout
	AllowlistUserIDs  []string   ` + "`" + `json:"allowlist_user_ids,omitempty"` + "`" + `  // when non-empty, ONLY these users get the flag
	BlocklistUserIDs  []string   ` + "`" + `json:"blocklist_user_ids,omitempty"` + "`" + `  // always-deny set; runs before allowlist + percentage
	EnabledFrom       *time.Time ` + "`" + `json:"enabled_from,omitempty"` + "`" + `        // before this, flag is off (date window)
	EnabledUntil      *time.Time ` + "`" + `json:"enabled_until,omitempty"` + "`" + `       // after this, flag is off
	Variants          []string   ` + "`" + `json:"variants,omitempty"` + "`" + `            // when set, A/B mode — Variant() returns one of these
}

// ParsedRules decodes the Rules JSON. Returns a zero FlagRules on
// missing or malformed JSON — callers shouldn't error out for
// misconfigured flags; they should fail closed (return false).
func (f *FeatureFlag) ParsedRules() FlagRules {
	var r FlagRules
	if len(f.Rules) > 0 {
		_ = json.Unmarshal(f.Rules, &r)
	}
	return r
}

// SetRules encodes a FlagRules and assigns it. Errors propagate.
func (f *FeatureFlag) SetRules(r FlagRules) error {
	b, err := json.Marshal(r)
	if err != nil {
		return err
	}
	f.Rules = b
	return nil
}

// FlagExposure records that a user was checked against a flag and what
// outcome they got. Used by the admin UI to show rollout health
// ("4,231 users saw variant_a, 4,189 saw variant_b") and to power
// downstream A/B analytics joins.
//
// Insert is fire-and-forget — exposure tracking should never block a
// flag check. We persist async in a goroutine.
type FlagExposure struct {
	ID        string    ` + "`" + `gorm:"primarykey;size:36" json:"id"` + "`" + `
	FlagID    string    ` + "`" + `gorm:"size:36;index;not null" json:"flag_id"` + "`" + `
	FlagName  string    ` + "`" + `gorm:"size:100;index" json:"flag_name"` + "`" + ` // denormalized for join-free analytics
	UserID    string    ` + "`" + `gorm:"size:36;index" json:"user_id"` + "`" + `
	Variant   string    ` + "`" + `gorm:"size:50" json:"variant"` + "`" + ` // "enabled" / "disabled" / "control" / "variant_a" / etc.
	CreatedAt time.Time ` + "`" + `gorm:"index" json:"created_at"` + "`" + `
}

func (e *FlagExposure) BeforeCreate(tx *gorm.DB) error {
	if e.ID == "" {
		e.ID = uuid.New().String()
	}
	return nil
}
`
}

func apiFlagsGo() string {
	return `// Package flags is the feature flag + A/B testing engine.
//
// At a glance:
//
//   if flags.IsEnabled(c, "new_dashboard") {
//       // … render the new dashboard
//   }
//
//   switch flags.Variant(c, "checkout_redesign") {
//   case "control":   /* old flow */
//   case "variant_a": /* new flow */
//   case "variant_b": /* alternate new flow */
//   }
//
// Mechanics:
//   - All flags are loaded into an in-memory map at boot. A background
//     goroutine refreshes every 30s. Flag checks never hit the DB.
//   - Bucketing: SHA-256(user_id || ":" || flag_name) % 100. Sticky
//     per (user, flag) — a user always gets the same bucket for a
//     given flag, so variant assignment doesn't flicker across sessions.
//   - Anonymous users (empty user_id) bucket on a random per-request
//     value, which is effectively random. For sticky anonymous flags
//     pass a stable identifier (session ID, device ID).
//   - Exposure tracking is fire-and-forget — flag checks never block
//     on the DB.
//   - When a flag is created/updated/deleted, the engine refreshes
//     immediately and broadcasts a "flag.updated" realtime event so
//     subscribed clients can refetch.
package flags

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"log"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"` + "{{MODULE}}" + `/internal/models"
	"` + "{{MODULE}}" + `/internal/realtime"
)

// DefaultRefreshInterval is how often the engine pulls fresh flag
// state from the DB. 30s is a reasonable middle ground — admin
// changes propagate quickly without hammering the DB.
const DefaultRefreshInterval = 30 * time.Second

// Engine owns the in-memory flag cache. One per process.
type Engine struct {
	db        *gorm.DB
	hub       *realtime.Hub // optional — when set, broadcasts on Refresh
	mu        sync.RWMutex
	flags     map[string]*models.FeatureFlag
	stop      chan struct{}
}

// New returns an Engine with the cache pre-warmed. Call from
// routes.Setup. hub is optional — pass nil to disable broadcasts.
func New(db *gorm.DB, hub *realtime.Hub) *Engine {
	e := &Engine{
		db:    db,
		hub:   hub,
		flags: make(map[string]*models.FeatureFlag),
		stop:  make(chan struct{}),
	}
	if err := e.Refresh(); err != nil {
		log.Printf("[flags] initial refresh failed: %v", err)
	}
	go e.refreshLoop()
	return e
}

// Stop terminates the background refresh goroutine. Call on graceful
// shutdown to avoid leaking goroutines in tests.
func (e *Engine) Stop() {
	close(e.stop)
}

// Refresh pulls all flags from the DB and replaces the cache. Called
// every DefaultRefreshInterval and immediately after admin writes.
func (e *Engine) Refresh() error {
	var rows []models.FeatureFlag
	if err := e.db.Find(&rows).Error; err != nil {
		return err
	}
	next := make(map[string]*models.FeatureFlag, len(rows))
	for i := range rows {
		f := rows[i]
		next[f.Name] = &f
	}
	e.mu.Lock()
	e.flags = next
	e.mu.Unlock()
	return nil
}

// RefreshAndBroadcast refreshes the cache and (if a hub was provided)
// emits a "flag.updated" realtime event so subscribed clients can
// refetch. Call after admin writes.
func (e *Engine) RefreshAndBroadcast(flagName string) {
	if err := e.Refresh(); err != nil {
		log.Printf("[flags] refresh after change failed: %v", err)
	}
	if e.hub != nil {
		e.hub.Broadcast(realtime.Event{
			Type:    "flag.updated",
			Payload: map[string]interface{}{"name": flagName},
		})
	}
}

func (e *Engine) refreshLoop() {
	t := time.NewTicker(DefaultRefreshInterval)
	defer t.Stop()
	for {
		select {
		case <-t.C:
			if err := e.Refresh(); err != nil {
				log.Printf("[flags] periodic refresh failed: %v", err)
			}
		case <-e.stop:
			return
		}
	}
}

// IsEnabled returns true when the flag is on for the current user.
// Always returns false for unknown flags (fail closed).
func (e *Engine) IsEnabled(c *gin.Context, name string) bool {
	return e.evaluate(userIDFrom(c), name) == "enabled"
}

// Variant returns the assigned variant for an A/B flag. For boolean
// flags, returns "enabled" or "disabled". For unknown flags, returns
// the empty string.
func (e *Engine) Variant(c *gin.Context, name string) string {
	return e.evaluate(userIDFrom(c), name)
}

// IsEnabledForUser is the explicit form for backend code that has the
// user_id directly (e.g. cron jobs operating on a specific user).
func (e *Engine) IsEnabledForUser(userID, name string) bool {
	return e.evaluate(userID, name) == "enabled"
}

// VariantForUser is the explicit form of Variant.
func (e *Engine) VariantForUser(userID, name string) string {
	return e.evaluate(userID, name)
}

// evaluate is the core decision routine. Returns:
//   ""           — unknown flag
//   "disabled"   — flag exists but rules deny the user
//   "enabled"    — boolean flag passed; user is in the rollout
//   "<variant>"  — A/B flag passed; the user's bucket maps to this variant
//
// Lock discipline: the read lock is held only long enough to copy the
// flag struct + ID. All decision logic (date checks, allowlist scans,
// bucketing) runs unlocked. Under sustained read load this turns the
// flag check into a near-zero-contention path.
func (e *Engine) evaluate(userID, name string) string {
	e.mu.RLock()
	cached, ok := e.flags[name]
	if !ok {
		e.mu.RUnlock()
		return ""
	}
	flagID := cached.ID
	enabled := cached.Enabled
	rulesJSON := cached.Rules
	e.mu.RUnlock()

	if !enabled {
		return "disabled"
	}

	// Decode rules outside the lock — JSON parsing is the slowest
	// part of the flag check and we don't want it serializing.
	flagForParse := models.FeatureFlag{Rules: rulesJSON}
	rules := flagForParse.ParsedRules()

	// Date window — out-of-window short-circuits before bucketing.
	now := time.Now()
	if rules.EnabledFrom != nil && now.Before(*rules.EnabledFrom) {
		return "disabled"
	}
	if rules.EnabledUntil != nil && now.After(*rules.EnabledUntil) {
		return "disabled"
	}

	// Blocklist always wins.
	for _, b := range rules.BlocklistUserIDs {
		if b == userID {
			return "disabled"
		}
	}

	// Allowlist (when non-empty) restricts to the listed users.
	// Skip the percentage roll for allowlisted users — they always
	// see it, that's the point.
	allowlistMode := len(rules.AllowlistUserIDs) > 0
	allowed := false
	for _, a := range rules.AllowlistUserIDs {
		if a == userID {
			allowed = true
			break
		}
	}
	if allowlistMode && !allowed {
		return "disabled"
	}

	bucket := bucketFor(userID, name)

	// A/B mode — assign variant by bucket.
	if len(rules.Variants) > 0 {
		v := rules.Variants[bucket%len(rules.Variants)]
		e.trackExposure(flagID, name, userID, v)
		return v
	}

	// Boolean mode — percentage rollout. Allowlisted users always
	// pass; everyone else is gated by the percentage.
	if allowed || bucket < rules.RolloutPercentage {
		e.trackExposure(flagID, name, userID, "enabled")
		return "enabled"
	}
	e.trackExposure(flagID, name, userID, "disabled")
	return "disabled"
}

// trackExposure records the flag check asynchronously. Never blocks
// the request path; logs failures.
func (e *Engine) trackExposure(flagID, flagName, userID, variant string) {
	if userID == "" {
		// Anonymous exposures pollute the table without buying us
		// anything (we can't link them to a user later). Skip.
		return
	}
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		err := e.db.WithContext(ctx).Create(&models.FlagExposure{
			FlagID:   flagID,
			FlagName: flagName,
			UserID:   userID,
			Variant:  variant,
		}).Error
		if err != nil {
			log.Printf("[flags] exposure insert failed: %v", err)
		}
	}()
}

// bucketFor hashes (userID || ":" || flagName) and returns the bucket
// 0..99. Same input always produces the same bucket — that's what
// makes the assignment sticky.
//
// We use SHA-256 (not Go's default hash) because it's stable across
// process restarts + Go versions. FNV would be faster but Grit isn't
// running flag checks in a hot loop — sub-microsecond cost is fine.
func bucketFor(userID, name string) int {
	if userID == "" {
		// Anonymous users get a uniform random bucket. We avoid
		// UnixNano%100 because nanosecond timing is biased toward
		// recent buckets under high QPS. crypto/rand gives us a
		// uniform draw without that artifact.
		var b [4]byte
		if _, err := rand.Read(b[:]); err != nil {
			// rand should never fail on a healthy OS; if it does,
			// fall back to bucket 0 so behavior is deterministic.
			return 0
		}
		return int(binary.BigEndian.Uint32(b[:]) % 100)
	}
	h := sha256.Sum256([]byte(userID + ":" + name))
	return int(binary.BigEndian.Uint32(h[:4]) % 100)
}

// userIDFrom reads "user_id" from the gin context (set by the auth
// middleware). Empty string for anonymous requests.
func userIDFrom(c *gin.Context) string {
	if v, ok := c.Get("user_id"); ok {
		s, _ := v.(string)
		return s
	}
	return ""
}
`
}

func apiFlagsHandlerGo() string {
	return `package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"` + "{{MODULE}}" + `/internal/flags"
	"` + "{{MODULE}}" + `/internal/models"
	"` + "{{MODULE}}" + `/internal/paginate"
)

// FeatureFlagHandler exposes admin-side CRUD over feature flags +
// exposure analytics. Mounted under admin/* (admin role required).
type FeatureFlagHandler struct {
	DB     *gorm.DB
	Engine *flags.Engine
}

func NewFeatureFlagHandler(db *gorm.DB, engine *flags.Engine) *FeatureFlagHandler {
	return &FeatureFlagHandler{DB: db, Engine: engine}
}

// List returns all flags with the standard paginate envelope.
//
//	GET /api/admin/flags
func (h *FeatureFlagHandler) List(c *gin.Context) {
	q := h.DB.Model(&models.FeatureFlag{})
	res, err := paginate.List[models.FeatureFlag](q, paginate.Bind(c), paginate.Config{
		Searchable:   []string{"name", "description"},
		Sortable:     map[string]bool{"name": true, "created_at": true, "enabled": true},
		DefaultSort:  "name",
		DefaultOrder: "asc",
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{"code": "INTERNAL_ERROR", "message": err.Error()},
		})
		return
	}
	c.JSON(http.StatusOK, res)
}

// flagPayload is the request shape for create/update. Rules is taken
// as a structured object — the handler encodes to JSON before hitting
// the DB so the wire format stays consistent.
type flagPayload struct {
	Name        string            ` + "`" + `json:"name"` + "`" + `
	Description string            ` + "`" + `json:"description"` + "`" + `
	Enabled     bool              ` + "`" + `json:"enabled"` + "`" + `
	Rules       models.FlagRules  ` + "`" + `json:"rules"` + "`" + `
}

// Create adds a new flag. Name must be unique.
//
//	POST /api/admin/flags
//	{ "name": "new_dashboard", "enabled": true, "rules": { "rollout_percentage": 25 } }
func (h *FeatureFlagHandler) Create(c *gin.Context) {
	var body flagPayload
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": gin.H{"code": "VALIDATION_ERROR", "message": err.Error()},
		})
		return
	}
	if body.Name == "" {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": gin.H{"code": "VALIDATION_ERROR", "message": "name is required"},
		})
		return
	}

	flag := models.FeatureFlag{
		Name:        body.Name,
		Description: body.Description,
		Enabled:     body.Enabled,
	}
	if err := flag.SetRules(body.Rules); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{"code": "INTERNAL_ERROR", "message": err.Error()},
		})
		return
	}
	if err := h.DB.Create(&flag).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{"code": "INTERNAL_ERROR", "message": err.Error()},
		})
		return
	}

	h.Engine.RefreshAndBroadcast(flag.Name)
	c.JSON(http.StatusCreated, gin.H{"data": flag})
}

// Update modifies an existing flag.
//
//	PUT /api/admin/flags/:id
func (h *FeatureFlagHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var flag models.FeatureFlag
	if err := h.DB.First(&flag, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": gin.H{"code": "NOT_FOUND", "message": "flag not found"},
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{"code": "INTERNAL_ERROR", "message": err.Error()},
		})
		return
	}

	var body flagPayload
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": gin.H{"code": "VALIDATION_ERROR", "message": err.Error()},
		})
		return
	}

	// Name is immutable post-create — too easy to break consumers.
	flag.Description = body.Description
	flag.Enabled = body.Enabled
	if err := flag.SetRules(body.Rules); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{"code": "INTERNAL_ERROR", "message": err.Error()},
		})
		return
	}
	if err := h.DB.Save(&flag).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{"code": "INTERNAL_ERROR", "message": err.Error()},
		})
		return
	}

	h.Engine.RefreshAndBroadcast(flag.Name)
	c.JSON(http.StatusOK, gin.H{"data": flag})
}

// Delete removes a flag. The cache refreshes immediately so app code
// stops seeing it on the next check.
//
//	DELETE /api/admin/flags/:id
func (h *FeatureFlagHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	var flag models.FeatureFlag
	if err := h.DB.First(&flag, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": gin.H{"code": "NOT_FOUND", "message": "flag not found"},
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{"code": "INTERNAL_ERROR", "message": err.Error()},
		})
		return
	}
	if err := h.DB.Delete(&flag).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{"code": "INTERNAL_ERROR", "message": err.Error()},
		})
		return
	}
	h.Engine.RefreshAndBroadcast(flag.Name)
	c.JSON(http.StatusOK, gin.H{"message": "flag deleted"})
}

// Exposures returns aggregate counts per variant for one flag —
// powers the rollout-health view in the admin UI.
//
//	GET /api/admin/flags/:id/exposures
//	→ { "data": [{ "variant": "enabled", "count": 4231 }, ...] }
func (h *FeatureFlagHandler) Exposures(c *gin.Context) {
	id := c.Param("id")
	type bucket struct {
		Variant string ` + "`" + `json:"variant"` + "`" + `
		Count   int64  ` + "`" + `json:"count"` + "`" + `
	}
	var rows []bucket
	if err := h.DB.Model(&models.FlagExposure{}).
		Select("variant, COUNT(DISTINCT user_id) as count").
		Where("flag_id = ?", id).
		Group("variant").
		Order("count desc").
		Scan(&rows).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{"code": "INTERNAL_ERROR", "message": err.Error()},
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": rows})
}
`
}

func apiActivityMiddlewareGo() string {
	return `package middleware

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"` + "{{MODULE}}" + `/internal/audit"
	"` + "{{MODULE}}" + `/internal/models"
)

// ActivityLogger records every successful authenticated mutation
// (POST/PUT/PATCH/DELETE) into models.ActivityLog. Skips:
//   - safe methods (GET/HEAD/OPTIONS)
//   - non-2xx responses (errors aren't audit-relevant)
//   - unauthenticated requests (no user_id ⇒ nothing to attribute)
//
// The payload digest is a SHA-256 hash of the request body — enough to
// prove "this exact payload was sent" without persisting plain-text
// passwords / secrets / PII. Buffered in memory, so MaxBodySize earlier
// in the chain still bounds it.
//
// Insert is fire-and-forget via a bounded channel + single writer
// goroutine. The single-writer design eliminates lock contention on
// the hash chain — only one goroutine ever appends — and the bounded
// channel caps memory + goroutine count under traffic spikes.
func ActivityLogger(db *gorm.DB) gin.HandlerFunc {
	auditOnce.Do(func() { go startAuditWorker(db) })
	return func(c *gin.Context) {
		switch c.Request.Method {
		case http.MethodGet, http.MethodHead, http.MethodOptions:
			c.Next()
			return
		}

		// Capture the body so we can hash it after the handler runs.
		// gin reads from c.Request.Body, so we tee it through a
		// bytes.Buffer and put a fresh ReadCloser back.
		var bodyBytes []byte
		if c.Request.Body != nil {
			bodyBytes, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}

		started := time.Now()
		c.Next()

		// Only log successful mutations — failed ones can be diagnosed
		// from request logs without polluting the audit trail.
		if c.Writer.Status() < 200 || c.Writer.Status() >= 300 {
			return
		}

		userID, _ := c.Get("user_id")
		uid, _ := userID.(string)
		if uid == "" {
			return // unauthenticated — nothing to audit
		}

		entry := models.ActivityLog{
			UserID:        uid,
			Method:        c.Request.Method,
			Path:          c.FullPath(),
			Status:        c.Writer.Status(),
			PayloadDigest: digestBody(bodyBytes),
			IPAddress:     c.ClientIP(),
			UserAgent:     c.Request.UserAgent(),
			DurationMS:    time.Since(started).Milliseconds(),
			CreatedAt:     time.Now(), // explicit — Canonical hashes this field
		}
		// Non-blocking enqueue. Channel is bounded so a runaway request
		// rate can't spawn unbounded goroutines or exhaust the DB pool.
		// On overflow we drop — better to lose an audit row than to
		// stall the request path or OOM the process.
		select {
		case auditChan <- entry:
		default:
			auditDropped.Add(1)
		}
	}
}

// auditChan is the bounded backlog for the single audit writer. 4096
// is enough to absorb a few-second burst (10k req/s for 0.4s) without
// blocking. The single-worker design also removes the need for a
// row-level FOR UPDATE lock on every write — chain integrity comes
// for free from sequential writes.
var (
	auditChan    = make(chan models.ActivityLog, 4096)
	auditOnce    sync.Once
	auditDropped atomicCounter
)

// auditDropped is exported via the integrity endpoint so ops can
// monitor when the audit channel saturates (signal to scale or
// reduce log noise).
type atomicCounter struct {
	mu sync.Mutex
	n  uint64
}

func (c *atomicCounter) Add(n uint64) {
	c.mu.Lock()
	c.n += n
	c.mu.Unlock()
}

// AuditDroppedCount returns the number of audit entries dropped due
// to channel saturation. Read this from a /healthz or admin endpoint
// to detect sustained back-pressure.
func AuditDroppedCount() uint64 {
	auditDropped.mu.Lock()
	defer auditDropped.mu.Unlock()
	return auditDropped.n
}

// startAuditWorker drains auditChan and writes each entry to the
// database with the hash chain attached. Single goroutine — no lock
// contention, no goroutine explosion, deterministic ordering.
//
// On boot the worker reads the latest persisted hash so the chain
// continues across restarts.
func startAuditWorker(db *gorm.DB) {
	var prev models.ActivityLog
	prevHash := ""
	if err := db.Order("created_at desc, id desc").Limit(1).First(&prev).Error; err == nil {
		prevHash = prev.Hash
	}

	for entry := range auditChan {
		canonical, err := audit.Canonical(&entry)
		if err != nil {
			log.Printf("[audit] canonicalize failed: %v", err)
			continue
		}
		entry.PrevHash = prevHash
		entry.Hash = audit.ComputeHash(prevHash, canonical)
		if err := db.Create(&entry).Error; err != nil {
			log.Printf("[audit] insert failed: %v", err)
			// Don't advance prevHash on failure — the next successful
			// write should chain off the last persisted row.
			continue
		}
		prevHash = entry.Hash
	}
}

func digestBody(b []byte) string {
	if len(b) == 0 {
		return ""
	}
	sum := sha256.Sum256(b)
	return hex.EncodeToString(sum[:])
}
`
}

func apiActivityHandlerGo() string {
	return `package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"` + "{{MODULE}}" + `/internal/audit"
	"` + "{{MODULE}}" + `/internal/models"
	"` + "{{MODULE}}" + `/internal/paginate"
)

// ActivityHandler exposes the audit log as a paginated, filterable
// list. Mounted under admin/* in routes.go.
type ActivityHandler struct {
	DB *gorm.DB
}

func NewActivityHandler(db *gorm.DB) *ActivityHandler {
	return &ActivityHandler{DB: db}
}

// List returns activity log entries, newest first. Supports filtering
// by user_id, method, and path prefix via query params.
func (h *ActivityHandler) List(c *gin.Context) {
	q := h.DB.Model(&models.ActivityLog{}).Order("created_at desc")
	params := paginate.Bind(c).
		With("user_id", c.Query("user_id")).
		With("method", c.Query("method"))

	if pathPrefix := c.Query("path"); pathPrefix != "" {
		q = q.Where("path LIKE ?", pathPrefix+"%")
	}

	res, err := paginate.List[models.ActivityLog](q, params, paginate.Config{
		Sortable: map[string]bool{
			"created_at": true,
			"status":     true,
			"method":     true,
		},
		DefaultSort:  "created_at",
		DefaultOrder: "desc",
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{"code": "INTERNAL_ERROR", "message": err.Error()},
		})
		return
	}
	c.JSON(http.StatusOK, res)
}

// VerifyIntegrity walks the entire activity log and verifies every
// row's Hash matches what we'd compute now. A mismatch means a row
// was modified, deleted, or inserted out of order — the response
// pinpoints which row broke the chain.
//
// Bounded by a 60-second deadline so a runaway scan can't hold the
// connection forever — if you have hundreds of millions of rows,
// run this from a cron job instead of an HTTP request.
//
//	GET /api/admin/activity/integrity
//	→ { "valid": true, "total_entries": 12345 }
//	→ { "valid": false, "broken_at": 47, "broken_at_id": "uuid",
//	    "expected": "abc...", "got": "def...",
//	    "message": "hash mismatch — row was modified..." }
func (h *ActivityHandler) VerifyIntegrity(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 60*time.Second)
	defer cancel()
	status, err := audit.VerifyChain(ctx, h.DB)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{"code": "INTERNAL_ERROR", "message": err.Error()},
		})
		return
	}
	c.JSON(http.StatusOK, status)
}
`
}

func apiExportGo() string {
	return `// Package export streams resource data out as CSV or XLSX. Used by
// auto-generated /<resource>/export endpoints — handlers reuse the
// List service layer to fetch rows, then call CSV(w, items, opts) or
// XLSX(w, items, opts) directly into the response writer.
//
// Column.Field uses Go-side struct field names with dot-notation for
// associations: "Tenant.Name", "Owner.Email", etc. Empty values render
// as empty strings.
//
// Format strings:
//   ""                — Sprintf %v (default)
//   "date:..."        — time.Time.Format(layout) — layout follows after the colon
//   "datetime"        — RFC3339-friendly date+time
//   "currency:CCC"    — formatted as "CCC 1,234.56"
//   "bool"            — "Yes" / "No"
package export

import (
	"encoding/csv"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"
)

// Column describes one output column.
type Column struct {
	Header string // human-readable column header
	Field  string // Go struct field path, e.g. "Tenant.Name"
	Format string // optional formatter — see package doc
}

// Options controls how items are rendered.
type Options struct {
	Columns []Column
	Sheet   string // XLSX only — defaults to "Sheet1"
}

// CSV writes items as a comma-separated stream into w. Includes the
// header row. For streaming exports (write headers once, then many
// batches) call CSV() for the first batch and CSVRows() for the rest.
func CSV(w io.Writer, items interface{}, opts Options) error {
	cw := csv.NewWriter(w)
	defer cw.Flush()

	headers := make([]string, len(opts.Columns))
	for i, col := range opts.Columns {
		headers[i] = col.Header
	}
	if err := cw.Write(headers); err != nil {
		return err
	}
	return writeCSVRows(cw, items, opts)
}

// CSVRows writes items WITHOUT a header row — used by streaming
// exports for batches after the first one (the header was already
// written by the initial CSV() call).
func CSVRows(w io.Writer, items interface{}, opts Options) error {
	cw := csv.NewWriter(w)
	defer cw.Flush()
	return writeCSVRows(cw, items, opts)
}

func writeCSVRows(cw *csv.Writer, items interface{}, opts Options) error {
	v := reflect.ValueOf(items)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Slice {
		return fmt.Errorf("export: items must be a slice, got %T", items)
	}

	for i := 0; i < v.Len(); i++ {
		row := make([]string, len(opts.Columns))
		for j, col := range opts.Columns {
			row[j] = formatCell(extractField(v.Index(i), col.Field), col.Format)
		}
		if err := cw.Write(row); err != nil {
			return err
		}
	}
	return nil
}

// XLSX writes items as an Excel workbook into w.
func XLSX(w io.Writer, items interface{}, opts Options) error {
	f := excelize.NewFile()
	defer f.Close()

	sheet := opts.Sheet
	if sheet == "" {
		sheet = "Sheet1"
	}
	if sheet != "Sheet1" {
		// excelize creates "Sheet1" by default; swap to the requested name.
		_ = f.SetSheetName("Sheet1", sheet)
	}

	for i, col := range opts.Columns {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		_ = f.SetCellValue(sheet, cell, col.Header)
	}

	v := reflect.ValueOf(items)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Slice {
		return fmt.Errorf("export: items must be a slice, got %T", items)
	}

	for i := 0; i < v.Len(); i++ {
		for j, col := range opts.Columns {
			cell, _ := excelize.CoordinatesToCellName(j+1, i+2)
			val := formatCell(extractField(v.Index(i), col.Field), col.Format)
			_ = f.SetCellValue(sheet, cell, val)
		}
	}

	return f.Write(w)
}

// extractField walks a dot-path through a struct. Returns the zero
// value if any segment is missing.
func extractField(v reflect.Value, path string) interface{} {
	if path == "" {
		return nil
	}
	for v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface {
		if v.IsNil() {
			return nil
		}
		v = v.Elem()
	}
	parts := strings.Split(path, ".")
	for _, p := range parts {
		if v.Kind() != reflect.Struct {
			return nil
		}
		f := v.FieldByName(p)
		if !f.IsValid() {
			return nil
		}
		v = f
		for v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface {
			if v.IsNil() {
				return nil
			}
			v = v.Elem()
		}
	}
	if !v.IsValid() {
		return nil
	}
	return v.Interface()
}

func formatCell(v interface{}, format string) string {
	if v == nil {
		return ""
	}
	if format == "" {
		return fmt.Sprintf("%v", v)
	}

	// "currency:UGX"
	if strings.HasPrefix(format, "currency:") {
		ccy := strings.TrimPrefix(format, "currency:")
		switch n := v.(type) {
		case float64:
			return ccy + " " + thousands(n)
		case float32:
			return ccy + " " + thousands(float64(n))
		case int:
			return ccy + " " + thousands(float64(n))
		case int64:
			return ccy + " " + thousands(float64(n))
		}
		return fmt.Sprintf("%v", v)
	}

	// "date:2006-01-02"
	if strings.HasPrefix(format, "date:") {
		layout := strings.TrimPrefix(format, "date:")
		if t, ok := v.(time.Time); ok {
			return t.Format(layout)
		}
	}

	if format == "datetime" {
		if t, ok := v.(time.Time); ok {
			return t.Format("2006-01-02 15:04")
		}
	}

	if format == "bool" {
		if b, ok := v.(bool); ok {
			if b {
				return "Yes"
			}
			return "No"
		}
	}

	return fmt.Sprintf("%v", v)
}

// thousands formats a float with thousands separators and 2 decimals.
func thousands(f float64) string {
	s := strconv.FormatFloat(f, 'f', 2, 64)
	parts := strings.SplitN(s, ".", 2)
	intPart := parts[0]
	neg := strings.HasPrefix(intPart, "-")
	if neg {
		intPart = intPart[1:]
	}
	var out []byte
	for i, c := range intPart {
		if i > 0 && (len(intPart)-i)%3 == 0 {
			out = append(out, ',')
		}
		out = append(out, byte(c))
	}
	result := string(out) + "." + parts[1]
	if neg {
		return "-" + result
	}
	return result
}
`
}

func apiRespondGo() string {
	return `// Package respond is the standard error/response envelope for handlers.
// Use these instead of writing c.JSON(500, gin.H{"error": err.Error()})
// inline so error shapes stay consistent and the frontend's
// apiErrorMessage() helper has a single shape to walk.
package respond

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Error is the wire shape of every error envelope.
type Error struct {
	Code    string            ` + "`" + `json:"code"` + "`" + `
	Message string            ` + "`" + `json:"message"` + "`" + `
	Details map[string]string ` + "`" + `json:"details,omitempty"` + "`" + `
}

// fail writes the standard error envelope at the given status code.
func fail(c *gin.Context, status int, code, message string, details ...map[string]string) {
	body := gin.H{"error": Error{Code: code, Message: message}}
	if len(details) > 0 {
		body = gin.H{"error": Error{Code: code, Message: message, Details: details[0]}}
	}
	c.AbortWithStatusJSON(status, body)
}

// 400 — malformed request that the client can't possibly fix without
// changing what it sent.
func BadRequest(c *gin.Context, message string) {
	fail(c, http.StatusBadRequest, "BAD_REQUEST", message)
}

// 401 — missing or invalid credentials.
func Unauthorized(c *gin.Context, message string) {
	if message == "" {
		message = "Authentication required"
	}
	fail(c, http.StatusUnauthorized, "UNAUTHORIZED", message)
}

// 403 — authenticated but not allowed.
func Forbidden(c *gin.Context, message string) {
	if message == "" {
		message = "You don't have permission to do that"
	}
	fail(c, http.StatusForbidden, "FORBIDDEN", message)
}

// 404 — entity didn't exist (or is filtered out by access rules).
func NotFound(c *gin.Context, message string) {
	if message == "" {
		message = "Not found"
	}
	fail(c, http.StatusNotFound, "NOT_FOUND", message)
}

// 409 — conflict (e.g. unique constraint, version conflict).
func Conflict(c *gin.Context, message string) {
	fail(c, http.StatusConflict, "CONFLICT", message)
}

// 422 — payload was well-formed but failed validation. Pass per-field
// errors via details map so the frontend can highlight them.
func Validation(c *gin.Context, message string, fields map[string]string) {
	fail(c, http.StatusUnprocessableEntity, "VALIDATION_ERROR", message, fields)
}

// 500 — server fault. Don't echo the raw error; log it and return a
// generic message so we don't leak internals.
func Internal(c *gin.Context, internalErr error) {
	msg := "Internal server error"
	if internalErr != nil {
		// In dev you may want the actual message. For now keep it
		// opaque; logger middleware records the full err.
		_ = internalErr
	}
	fail(c, http.StatusInternalServerError, "INTERNAL_ERROR", msg)
}

// OK writes 200 with { data, message? }.
func OK(c *gin.Context, data interface{}, message ...string) {
	body := gin.H{"data": data}
	if len(message) > 0 && message[0] != "" {
		body["message"] = message[0]
	}
	c.JSON(http.StatusOK, body)
}

// Created writes 201 with { data, message? }.
func Created(c *gin.Context, data interface{}, message ...string) {
	body := gin.H{"data": data}
	if len(message) > 0 && message[0] != "" {
		body["message"] = message[0]
	}
	c.JSON(http.StatusCreated, body)
}
`
}

func apiPDFGo() string {
	return `// Package pdf is a tiny styled-PDF builder backed by go-pdf/fpdf.
//
// The package exports two layers:
//
//   1) Doc primitives — Header, KV, Table, Totals, Notes, Footer — that
//      apply Grit's default styling (Helvetica, 20mm margins, blue
//      accent, A4 portrait). Compose them to build any document.
//
//   2) Pre-built templates — RenderInvoice (in invoice.go) — for the
//      common business-app cases. Copy + adapt these for receipts,
//      leases, statements, etc.
//
// When the helpers don't fit, the embedded *fpdf.Fpdf gives you the
// full underlying API. Call d.Bytes() at the end to finalize.
package pdf

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/go-pdf/fpdf"
)

// Doc wraps fpdf.Fpdf with section helpers + Grit's default colors.
// Mutate Accent on the returned Doc to retheme.
type Doc struct {
	*fpdf.Fpdf
	Accent [3]int // RGB; default Grit blue (30, 126, 245)
	Muted  [3]int // RGB; default neutral gray (110, 110, 110)
}

// New returns a fresh A4 portrait document with Grit's default styling.
// Adds the first page automatically — call d.AddPage() for additional
// pages.
func New() *Doc {
	pdf := fpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(20, 18, 20)
	pdf.AddPage()
	pdf.SetFont("Helvetica", "", 10)
	return &Doc{
		Fpdf:   pdf,
		Accent: [3]int{30, 126, 245},
		Muted:  [3]int{110, 110, 110},
	}
}

// Header writes the standard top-of-document title bar — accent-colored
// title in 22pt + a smaller secondary line below in muted gray.
//
//	d.Header("INVOICE", "INV-202605-0001")
//	d.Header("RECEIPT", "RCT-202605-0042")
func (d *Doc) Header(title, subtitle string) {
	d.SetFont("Helvetica", "B", 22)
	d.SetTextColor(d.Accent[0], d.Accent[1], d.Accent[2])
	d.CellFormat(0, 10, strings.ToUpper(title), "", 1, "L", false, 0, "")
	if subtitle != "" {
		d.SetTextColor(d.Muted[0], d.Muted[1], d.Muted[2])
		d.SetFont("Helvetica", "", 10)
		d.CellFormat(0, 5, subtitle, "", 1, "L", false, 0, "")
	}
	d.SetTextColor(0, 0, 0)
	d.Ln(6)
}

// KV writes a "label: value" pair. Label is bold + small caps style;
// value is regular weight on the next line. Used for "Bill To",
// "Issue Date", "Reference Number", etc.
func (d *Doc) KV(label, value string) {
	d.SetFont("Helvetica", "B", 9)
	d.SetTextColor(d.Muted[0], d.Muted[1], d.Muted[2])
	d.CellFormat(0, 5, strings.ToUpper(label), "", 1, "L", false, 0, "")
	d.SetFont("Helvetica", "", 10)
	d.SetTextColor(0, 0, 0)
	d.CellFormat(0, 5, value, "", 1, "L", false, 0, "")
	d.Ln(3)
}

// TwoColumnKV writes two KV pairs side by side — useful for fitting
// "BILL TO" + "ISSUE DATE" or "FROM" + "TO" on one row.
func (d *Doc) TwoColumnKV(leftLabel, leftValue, rightLabel, rightValue string) {
	d.SetFont("Helvetica", "B", 9)
	d.SetTextColor(d.Muted[0], d.Muted[1], d.Muted[2])
	d.CellFormat(95, 5, strings.ToUpper(leftLabel), "", 0, "L", false, 0, "")
	d.CellFormat(0, 5, strings.ToUpper(rightLabel), "", 1, "L", false, 0, "")
	d.SetFont("Helvetica", "", 10)
	d.SetTextColor(0, 0, 0)
	d.CellFormat(95, 5, leftValue, "", 0, "L", false, 0, "")
	d.CellFormat(0, 5, rightValue, "", 1, "L", false, 0, "")
	d.Ln(3)
}

// Table writes a styled table. headers + rows are matched by index.
// colWidths are in mm — pass 0 for the last column to fill remaining
// width. Header row gets a light gray background; data rows are plain.
//
//	d.Table(
//	    []string{"DESCRIPTION", "QTY", "UNIT", "TOTAL"},
//	    [][]string{
//	        {"Office rent — June", "1", "1,500,000", "1,500,000"},
//	        {"Service charge",      "1",   "120,000",   "120,000"},
//	    },
//	    []float64{105, 15, 25, 0},
//	    []string{"L", "R", "R", "R"},
//	)
func (d *Doc) Table(headers []string, rows [][]string, colWidths []float64, aligns []string) {
	if len(headers) == 0 || len(colWidths) != len(headers) {
		return
	}
	if len(aligns) != len(headers) {
		// Default all-left if alignment slice is malformed.
		aligns = make([]string, len(headers))
		for i := range aligns {
			aligns[i] = "L"
		}
	}

	// Header row
	d.SetFillColor(244, 244, 245)
	d.SetFont("Helvetica", "B", 9)
	d.SetTextColor(120, 120, 120)
	for i, h := range headers {
		end := 0
		if i == len(headers)-1 {
			end = 1
		}
		d.CellFormat(colWidths[i], 7, h, "", end, aligns[i], true, 0, "")
	}

	// Data rows
	d.SetTextColor(0, 0, 0)
	d.SetFont("Helvetica", "", 10)
	for _, row := range rows {
		for i, cell := range row {
			if i >= len(colWidths) {
				break
			}
			end := 0
			if i == len(row)-1 {
				end = 1
			}
			d.CellFormat(colWidths[i], 6, cell, "", end, aligns[i], false, 0, "")
		}
	}
}

// TotalLine is one entry in a Totals stack.
type TotalLine struct {
	Label string
	Value string // pre-formatted with currency + thousands separators
	Bold  bool   // bold + accent color (for the grand total line)
}

// Totals writes a right-aligned totals stack. The last "Bold" line
// gets accent coloring + a slightly larger size — typically used for
// the grand total or outstanding balance.
func (d *Doc) Totals(lines []TotalLine) {
	for _, line := range lines {
		d.CellFormat(120, 6, "", "", 0, "L", false, 0, "")
		d.SetFont("Helvetica", "", 10)
		d.SetTextColor(d.Muted[0], d.Muted[1], d.Muted[2])
		d.CellFormat(20, 6, line.Label, "", 0, "R", false, 0, "")
		if line.Bold {
			d.SetFont("Helvetica", "B", 11)
			d.SetTextColor(d.Accent[0], d.Accent[1], d.Accent[2])
		} else {
			d.SetFont("Helvetica", "", 10)
			d.SetTextColor(0, 0, 0)
		}
		d.CellFormat(0, 6, line.Value, "", 1, "R", false, 0, "")
	}
	d.SetTextColor(0, 0, 0)
}

// Notes writes a "NOTES" header + the body text wrapped to page width.
// Skipped silently when text is empty.
func (d *Doc) Notes(text string) {
	if text == "" {
		return
	}
	d.Ln(4)
	d.SetFont("Helvetica", "B", 10)
	d.SetTextColor(0, 0, 0)
	d.CellFormat(0, 5, "NOTES", "", 1, "L", false, 0, "")
	d.SetFont("Helvetica", "", 10)
	d.SetTextColor(82, 82, 91)
	d.MultiCell(0, 5, text, "", "L", false)
	d.Ln(4)
}

// Footer writes a centered italic footer 25mm from the bottom of the
// page. Use for "Generated 2 Jun 2026 14:30" or terms of service URLs.
func (d *Doc) Footer(text string) {
	d.SetY(-25)
	d.SetFont("Helvetica", "I", 9)
	d.SetTextColor(d.Muted[0], d.Muted[1], d.Muted[2])
	d.CellFormat(0, 5, text, "", 1, "C", false, 0, "")
}

// Bytes finalizes the document and returns the PDF bytes. Call this
// once at the very end — the underlying fpdf is not reusable after.
func (d *Doc) Bytes() ([]byte, error) {
	var buf bytes.Buffer
	if err := d.Output(&buf); err != nil {
		return nil, fmt.Errorf("pdf output: %w", err)
	}
	return buf.Bytes(), nil
}
`
}

func apiPDFInvoiceGo() string {
	return `package pdf

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Invoice is the data shape RenderInvoice consumes. Build one from
// your domain model in the handler and pass it through.
type Invoice struct {
	Number    string    // "INV-202605-0001"
	IssueDate time.Time
	DueDate   time.Time
	BillTo    Party
	From      Party     // your company — optional, shown in the header area
	Items     []LineItem
	Subtotal  float64
	Tax       float64   // tax amount (not rate)
	Total     float64
	Paid      float64   // amount already paid; if > 0, an "Outstanding" line is added
	Currency  string    // "UGX", "USD", etc. — prefixed to every amount
	Notes     string    // free-text footer notes
	Status    string    // shown in the document footer ("paid", "overdue", "draft")
}

// Party is a name + free-text contact block (phone, email, address).
type Party struct {
	Name    string
	Contact string
}

// LineItem is one row in the invoice's items table.
type LineItem struct {
	Description string
	Quantity    float64
	UnitPrice   float64
	Total       float64
}

// RenderInvoice returns the invoice as PDF bytes ready to stream to
// the response writer. Composition over inheritance: it's just a Doc
// with the section helpers called in order — copy this file as a
// starting point for receipts / leases / statements / quotes.
//
//	GET /api/invoices/:id/pdf
//	    inv := h.Service.GetByID(c.Param("id"))
//	    bytes, _ := pdf.RenderInvoice(toInvoice(inv))
//	    c.Data(200, "application/pdf", bytes)
func RenderInvoice(inv Invoice) ([]byte, error) {
	d := New()

	d.Header("INVOICE", inv.Number)
	d.TwoColumnKV("BILL TO", inv.BillTo.Name, "ISSUE DATE", inv.IssueDate.Format("2 Jan 2006"))
	if inv.BillTo.Contact != "" {
		d.SetFont("Helvetica", "", 10)
		d.SetTextColor(d.Muted[0], d.Muted[1], d.Muted[2])
		d.CellFormat(95, 5, inv.BillTo.Contact, "", 0, "L", false, 0, "")
		d.SetTextColor(0, 0, 0)
		d.SetFont("Helvetica", "B", 9)
		d.CellFormat(0, 5, "DUE DATE", "", 1, "L", false, 0, "")
		d.CellFormat(95, 5, "", "", 0, "L", false, 0, "")
		d.SetFont("Helvetica", "", 10)
		d.CellFormat(0, 5, inv.DueDate.Format("2 Jan 2006"), "", 1, "L", false, 0, "")
	}
	d.Ln(8)

	// Items table
	rows := make([][]string, len(inv.Items))
	for i, it := range inv.Items {
		rows[i] = []string{
			it.Description,
			strconv.FormatFloat(it.Quantity, 'f', -1, 64),
			formatAmount(it.UnitPrice),
			formatAmount(it.Total),
		}
	}
	d.Table(
		[]string{"DESCRIPTION", "QTY", "UNIT", "TOTAL"},
		rows,
		[]float64{105, 15, 25, 0},
		[]string{"L", "R", "R", "R"},
	)
	d.Ln(4)

	// Totals
	totals := []TotalLine{
		{Label: "Subtotal", Value: inv.Currency + " " + formatAmount(inv.Subtotal)},
	}
	if inv.Tax > 0 {
		totals = append(totals, TotalLine{Label: "Tax", Value: inv.Currency + " " + formatAmount(inv.Tax)})
	}
	totals = append(totals, TotalLine{Label: "Total", Value: inv.Currency + " " + formatAmount(inv.Total), Bold: true})
	if inv.Paid > 0 {
		totals = append(totals,
			TotalLine{Label: "Paid", Value: inv.Currency + " " + formatAmount(inv.Paid)},
			TotalLine{Label: "Outstanding", Value: inv.Currency + " " + formatAmount(inv.Total - inv.Paid), Bold: true},
		)
	}
	d.Totals(totals)

	d.Notes(inv.Notes)

	footer := fmt.Sprintf("Generated %s", time.Now().Format("2 Jan 2006 15:04"))
	if inv.Status != "" {
		footer += " · Status: " + inv.Status
	}
	d.Footer(footer)

	return d.Bytes()
}

// formatAmount renders 1234567.89 as "1,234,567.89" — matches the
// thousands-separator convention used by the export package.
func formatAmount(n float64) string {
	s := strconv.FormatFloat(n, 'f', 2, 64)
	parts := strings.SplitN(s, ".", 2)
	intPart := parts[0]
	neg := strings.HasPrefix(intPart, "-")
	if neg {
		intPart = intPart[1:]
	}
	var out []byte
	for i, c := range intPart {
		if i > 0 && (len(intPart)-i)%3 == 0 {
			out = append(out, ',')
		}
		out = append(out, byte(c))
	}
	result := string(out) + "." + parts[1]
	if neg {
		return "-" + result
	}
	return result
}
`
}

func apiRealtimeHubGo() string {
	return `// Package realtime is a tiny WebSocket fan-out hub. One Hub per process;
// each authenticated user can have multiple connections (e.g. desktop +
// mobile + web). The hub owns the registry and exposes safe SendToUser /
// SendToUsers / Broadcast helpers that handlers call from anywhere.
//
// Wire format on the websocket is a JSON envelope:
//
//	{ "type": "<topic>", "payload": { ... } }
//
// Topics are caller-defined strings. Suggested namespacing:
//
//   chat.message.new       — payload is a chat message
//   notification.new       — payload is a notification
//   system.connected       — server greeting on first connect
//   resource.<name>.<verb> — e.g. building.created, lease.expired
package realtime

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

// Event is the envelope every WS message uses on the wire.
type Event struct {
	Type    string      ` + "`" + `json:"type"` + "`" + `
	Payload interface{} ` + "`" + `json:"payload"` + "`" + `
}

// Client is one open WebSocket connection bound to a user.
type Client struct {
	UserID string
	Conn   *websocket.Conn
	Send   chan []byte
}

// Hub manages connected clients. Safe for concurrent use.
type Hub struct {
	mu      sync.RWMutex
	clients map[string]map[*Client]struct{} // userID -> set of connections
}

// NewHub returns an empty Hub.
func NewHub() *Hub {
	return &Hub{clients: make(map[string]map[*Client]struct{})}
}

// Register adds a client to the hub. A user can have multiple registered
// clients (different devices); each gets its own slot.
func (h *Hub) Register(c *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	set, ok := h.clients[c.UserID]
	if !ok {
		set = make(map[*Client]struct{})
		h.clients[c.UserID] = set
	}
	set[c] = struct{}{}
	log.Printf("[realtime] client registered user=%s total=%d", c.UserID, len(set))
}

// Unregister removes a client and closes its Send channel. Safe to call
// once per client (e.g. from the read pump's defer).
func (h *Hub) Unregister(c *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if set, ok := h.clients[c.UserID]; ok {
		if _, exists := set[c]; exists {
			delete(set, c)
			close(c.Send)
		}
		if len(set) == 0 {
			delete(h.clients, c.UserID)
		}
	}
}

// SendToUser delivers an event to every connection bound to userID.
// If a connection's send buffer is full the message is dropped for that
// connection only — we never block the entire hub on a slow client.
// The slow client will resync on its next REST poll/refetch.
func (h *Hub) SendToUser(userID string, evt Event) {
	bytes, err := json.Marshal(evt)
	if err != nil {
		log.Printf("[realtime] marshal: %v", err)
		return
	}
	h.mu.RLock()
	set := h.clients[userID]
	targets := make([]*Client, 0, len(set))
	for c := range set {
		targets = append(targets, c)
	}
	h.mu.RUnlock()
	for _, c := range targets {
		select {
		case c.Send <- bytes:
		default:
			log.Printf("[realtime] dropping message for slow client user=%s", userID)
		}
	}
}

// SendToUsers fans out to a slice of user IDs.
func (h *Hub) SendToUsers(userIDs []string, evt Event) {
	for _, uid := range userIDs {
		h.SendToUser(uid, evt)
	}
}

// Broadcast delivers an event to every connected client, regardless of
// user. Use sparingly — for system-wide announcements, maintenance
// notices, etc.
func (h *Hub) Broadcast(evt Event) {
	bytes, err := json.Marshal(evt)
	if err != nil {
		log.Printf("[realtime] marshal: %v", err)
		return
	}
	h.mu.RLock()
	targets := make([]*Client, 0)
	for _, set := range h.clients {
		for c := range set {
			targets = append(targets, c)
		}
	}
	h.mu.RUnlock()
	for _, c := range targets {
		select {
		case c.Send <- bytes:
		default:
			log.Printf("[realtime] dropping broadcast for slow client user=%s", c.UserID)
		}
	}
}
`
}

func apiRealtimeHandlerGo() string {
	return `package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"

	"` + "{{MODULE}}" + `/internal/realtime"
	"` + "{{MODULE}}" + `/internal/services"
)

const (
	wsWriteWait      = 10 * time.Second
	wsPongWait       = 60 * time.Second
	wsPingPeriod     = (wsPongWait * 9) / 10
	wsMaxMessageSize = 1024 // we don't expect clients to send anything large
)

// upgrader allows any origin — desktop clients use Wails (file://) and
// the API is mounted behind CORS that already restricts origins for
// regular HTTP traffic.
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

// RealtimeHandler upgrades an HTTP request to a WebSocket and registers
// it with the hub. Authentication uses a query-string JWT (?token=...)
// because browsers can't set custom Authorization headers on WebSocket
// handshakes — there is no other portable way to pass the JWT.
type RealtimeHandler struct {
	Hub  *realtime.Hub
	Auth *services.AuthService
}

// NewRealtimeHandler wires the handler to the global Hub and AuthService.
func NewRealtimeHandler(hub *realtime.Hub, auth *services.AuthService) *RealtimeHandler {
	return &RealtimeHandler{Hub: hub, Auth: auth}
}

// Connect upgrades the request to a WebSocket connection.
//
//   GET /api/ws?token=<jwt>
func (h *RealtimeHandler) Connect(c *gin.Context) {
	tokenStr := c.Query("token")
	if tokenStr == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": gin.H{"code": "MISSING_TOKEN", "message": "?token query is required"}})
		return
	}
	claims, err := h.Auth.ValidateToken(tokenStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": gin.H{"code": "INVALID_TOKEN", "message": err.Error()}})
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("[ws] upgrade error: %v", err)
		return
	}

	client := &realtime.Client{
		UserID: claims.UserID,
		Conn:   conn,
		Send:   make(chan []byte, 32),
	}
	h.Hub.Register(client)

	// Greeting so the client knows the link is live.
	greeting, _ := json.Marshal(realtime.Event{
		Type:    "system.connected",
		Payload: gin.H{"user_id": claims.UserID},
	})
	select {
	case client.Send <- greeting:
	default:
	}

	go writePump(client)
	go readPump(h.Hub, client)
}

// readPump pumps messages from the client → hub. We don't currently
// accept commands from clients (mutations go through the REST API), so
// this loop just services ping/pong and cleans up on disconnect.
func readPump(hub *realtime.Hub, c *realtime.Client) {
	defer func() {
		hub.Unregister(c)
		_ = c.Conn.Close()
	}()
	c.Conn.SetReadLimit(wsMaxMessageSize)
	_ = c.Conn.SetReadDeadline(time.Now().Add(wsPongWait))
	c.Conn.SetPongHandler(func(string) error {
		_ = c.Conn.SetReadDeadline(time.Now().Add(wsPongWait))
		return nil
	})
	for {
		if _, _, err := c.Conn.ReadMessage(); err != nil {
			return
		}
	}
}

// writePump pumps messages from the hub → client and emits keepalive pings.
func writePump(c *realtime.Client) {
	ticker := time.NewTicker(wsPingPeriod)
	defer func() {
		ticker.Stop()
		_ = c.Conn.Close()
	}()
	for {
		select {
		case msg, ok := <-c.Send:
			_ = c.Conn.SetWriteDeadline(time.Now().Add(wsWriteWait))
			if !ok {
				_ = c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.Conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				return
			}
		case <-ticker.C:
			_ = c.Conn.SetWriteDeadline(time.Now().Add(wsWriteWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
`
}

func apiRoutesGo() string {
	return `package routes

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/MUKE-coder/gin-docs/gindocs"
	"github.com/MUKE-coder/gorm-studio/studio"
	"github.com/MUKE-coder/pulse/pulse"
	"github.com/MUKE-coder/sentinel"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"` + "{{MODULE}}" + `/internal/ai"
	"` + "{{MODULE}}" + `/internal/cache"
	"` + "{{MODULE}}" + `/internal/config"
	"` + "{{MODULE}}" + `/internal/handlers"
	"` + "{{MODULE}}" + `/internal/mail"
	"` + "{{MODULE}}" + `/internal/middleware"
	"` + "{{MODULE}}" + `/internal/models"
	"` + "{{MODULE}}" + `/internal/jobs"
	"` + "{{MODULE}}" + `/internal/realtime"
	"` + "{{MODULE}}" + `/internal/services"
	"` + "{{MODULE}}" + `/internal/storage"
	"` + "{{MODULE}}" + `/internal/flags"
	"` + "{{MODULE}}" + `/internal/sync"
	"` + "{{MODULE}}" + `/internal/webhooks"
)

// Services holds all Phase 4 services for dependency injection.
type Services struct {
	Cache   *cache.Cache
	Storage *storage.Storage
	Mailer  *mail.Mailer
	AI      *ai.AI
	Jobs    *jobs.Client
	// SecObsBridge talks to Sentinel + Pulse over loopback so the
	// in-app Security/Observability dashboards can show summary cards
	// without iframing. Nil when Sentinel/Pulse are both disabled.
	SecObs  *services.SecObsBridge
}

// Setup configures all routes and returns the Gin engine.
func Setup(db *gorm.DB, cfg *config.Config, svc *Services) *gin.Engine {
	if cfg.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()

	// Global middleware
	r.Use(middleware.Maintenance())
	r.Use(middleware.SecurityHeaders())
	r.Use(middleware.MaxBodySize(10 << 20)) // 10MB max request body
	r.Use(middleware.RequestID())
	r.Use(middleware.Logger())
	r.Use(gin.Recovery())
	r.Use(middleware.CORS(cfg.CORSOrigins))
	r.Use(middleware.Gzip())

	// CSRF defence — only enforces on cookie-authenticated mutations.
	// Bearer (mobile/desktop) flows pass through with no header required.
	// Pairs with services.AuthService.SetAuthCookies (the HttpOnly cookie
	// path documented in /docs/backend/authentication).
	r.Use(middleware.AutoCSRF())

	// Idempotent retries for unsafe methods. Activates only when the client
	// sends an Idempotency-Key header; cached for 24h on 2xx responses.
	r.Use(middleware.Idempotency(svc.Cache))

	// Mount Sentinel security suite (WAF, rate limiting, auth shield, anomaly detection)
	if cfg.SentinelEnabled {
		// In development, use relaxed rate limits so devs don't get blocked while testing
		isDev := cfg.AppEnv == "development"
		ipLimit := &sentinel.Limit{Requests: 100, Window: 1 * time.Minute}
		routeLimits := map[string]sentinel.Limit{
			"/api/auth/login":    {Requests: 5, Window: 15 * time.Minute},
			"/api/auth/register": {Requests: 3, Window: 15 * time.Minute},
		}
		if isDev {
			ipLimit = &sentinel.Limit{Requests: 1000, Window: 1 * time.Minute}
			routeLimits = map[string]sentinel.Limit{
				"/api/auth/login":    {Requests: 100, Window: 1 * time.Minute},
				"/api/auth/register": {Requests: 100, Window: 1 * time.Minute},
			}
		}

		// Sentinel v2.0.1 — use MountE so we can recover gracefully on
		// misconfiguration in dev instead of log.Fatalf-ing the host.
		if err := sentinel.MountE(r, db, sentinel.Config{
			Dashboard: sentinel.DashboardConfig{
				Username:               cfg.SentinelUsername,
				Password:               cfg.SentinelPassword,
				SecretKey:              cfg.SentinelSecretKey,
				// v2.0 refuses default credentials in gin.ReleaseMode;
				// opt-in only for dev so prod can't ship forgeable JWTs.
				AllowInsecureDefaults:  isDev,
			},
			WAF: sentinel.WAFConfig{
				Enabled: true,
				Mode: func() sentinel.WAFMode {
					if isDev { return sentinel.ModeLog }
					return sentinel.ModeBlock
				}(),
				// v2.0 X-Forwarded-For trust closed. Empty list = ignore
				// XFF entirely (the safe default). Operators behind a known
				// reverse proxy should populate via SENTINEL_TRUSTED_PROXIES.
				TrustedProxies:        cfg.SentinelTrustedProxies,
				// v2.0 body-inspection cap; rejects bodies larger than
				// MaxBodyBytes outright so attackers can't pad past it.
				MaxBodyBytes:          64 * 1024,
				RejectOversizedBody:   true,
			},
			RateLimit: sentinel.RateLimitConfig{
				Enabled: !isDev,
				ByIP:    ipLimit,
				ByRoute: routeLimits,
			},
			AuthShield: sentinel.AuthShieldConfig{
				Enabled:    !isDev,
				LoginRoute: "/api/auth/login",
				// v2.0 CAPTCHA tier sits between soft and hard thresholds.
				// Wire a provider by setting CaptchaProvider in your app code.
			},
			Anomaly: sentinel.AnomalyConfig{Enabled: !isDev},
			Geo:     sentinel.GeoConfig{Enabled: !isDev},
		}); err != nil {
			log.Printf("Warning: Sentinel mount failed: %v", err)
		} else {
			log.Println("Sentinel v2.0 mounted at /sentinel")
		}
	}

	// Mount GORM Studio
	if cfg.GORMStudioEnabled {
		studioCfg := studio.Config{
			Prefix: "/studio",
		}
		if cfg.GORMStudioUsername != "" && cfg.GORMStudioPassword != "" {
			studioCfg.AuthMiddleware = gin.BasicAuth(gin.Accounts{
				cfg.GORMStudioUsername: cfg.GORMStudioPassword,
			})
		}
		studio.Mount(r, db, []interface{}{&models.User{}, &models.Upload{}, &models.Blog{}, /* grit:studio */}, studioCfg)
		log.Println("GORM Studio mounted at /studio")
	}

	// API Documentation (gin-docs — auto-generated from routes + models)
	gindocs.Mount(r, db, gindocs.Config{
		Title:       cfg.AppName + " API",
		Description: "REST API built with [Grit](https://gritframework.dev) — Go + React meta-framework.",
		Version:     "1.0.0",
		UI:          gindocs.UIScalar,
		ScalarTheme: "kepler",
		Models:      []interface{}{&models.User{}, &models.Upload{}, &models.Blog{}},
		Auth: gindocs.AuthConfig{
			Type:         gindocs.AuthBearer,
			BearerFormat: "JWT",
		},
	})
	log.Println("API docs available at /docs")

	// Mount Pulse observability (request tracing, DB monitoring, runtime metrics, error tracking)
	if cfg.PulseEnabled {
		// Pulse v1.0 uses functional options + a context. The context
		// drives clean shutdown of the dashboard's WebSocket + background
		// samplers; we hand it the request context so a server shutdown
		// also unwinds Pulse.
		pulseOpts := []pulse.Option{
			pulse.WithAppName(cfg.AppName),
			pulse.WithCredentials(cfg.PulseUsername, cfg.PulsePassword),
			pulse.WithExcludePaths("/studio/*", "/sentinel/*", "/docs/*", "/pulse/*"),
			pulse.WithPrometheus(),
		}
		if cfg.IsDevelopment() {
			pulseOpts = append(pulseOpts, pulse.WithDevMode())
		}
		// Pulse v1.0 SQLite-backed storage — request/query/error data
		// survives a restart. Stay on the in-memory ring buffer for peak
		// write throughput.
		if cfg.PulseStorage == "sqlite" && cfg.PulseStorageDSN != "" {
			pulseOpts = append(pulseOpts, pulse.WithSQLite(cfg.PulseStorageDSN))
		}
		p := pulse.Mount(context.Background(), r, db, pulseOpts...)

		// Register health checks for connected services
		if svc.Cache != nil {
			p.AddHealthCheck(pulse.HealthCheck{
				Name:     "redis",
				Type:     "redis",
				Critical: false,
				CheckFunc: func(ctx context.Context) error {
					return svc.Cache.Client().Ping(ctx).Err()
				},
			})
		}

		log.Println("Pulse observability mounted at /pulse")
	}

	// Auth service
	authService := &services.AuthService{
		Secret:        cfg.JWTSecret,
		AccessExpiry:  cfg.JWTAccessExpiry,
		RefreshExpiry: cfg.JWTRefreshExpiry,
	}

	// Handlers
	authHandler := &handlers.AuthHandler{
		DB:          db,
		AuthService: authService,
		Config:      cfg,
	}
	userHandler := &handlers.UserHandler{
		DB: db,
	}
	uploadHandler := &handlers.UploadHandler{
		DB:      db,
		Storage: svc.Storage,
		Jobs:    svc.Jobs,
	}
	aiHandler := &handlers.AIHandler{
		AI: svc.AI,
	}
	jobsHandler := &handlers.JobsHandler{
		RedisURL: cfg.RedisURL,
	}
	cronHandler := &handlers.CronHandler{}
	blogHandler := handlers.NewBlogHandler(db)
	uiRegistryHandler := handlers.NewUIRegistryHandler(db, cfg.AppURL)
	totpHandler := &handlers.TOTPHandler{
		DB:          db,
		AuthService: authService,
		Issuer:      cfg.TOTPIssuer,
	}
	activityHandler := handlers.NewActivityHandler(db)
	webhookHandler := handlers.NewWebhookHandler(db)
	webhooks.Setup(db)
	realtimeHub := realtime.NewHub()
	flagsEngine := flags.New(db, realtimeHub)
	featureFlagHandler := handlers.NewFeatureFlagHandler(db, flagsEngine)
	realtimeHandler := handlers.NewRealtimeHandler(realtimeHub, authService)
	_ = realtimeHub // available to handlers/services that want to push events

	// In-app Security + Observability dashboards — read from Sentinel/Pulse APIs
	// over loopback. notificationHandler powers the admin bell.
	notificationHandler := &handlers.NotificationHandler{DB: db}
	securityHandler := &handlers.SecurityHandler{Bridge: svc.SecObs}
	observabilityHandler := &handlers.ObservabilityHandler{Bridge: svc.SecObs}

	// Sync registry — list every model that should be syncable from
	// offline-first desktop clients. The resource generator injects
	// new resources at the marker below.
	syncRegistry := sync.NewRegistry()
	syncRegistry.Register("users", &models.User{})
	syncRegistry.Register("uploads", &models.Upload{})
	syncRegistry.Register("blogs", &models.Blog{})
	// grit:sync
	syncHandler := handlers.NewSyncHandler(db, syncRegistry)
	// grit:handlers

	// Health check
	r.GET("/api/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"version": "0.1.0",
		})
	})

	// WebSocket: realtime hub. Auth via ?token=<jwt> on the handshake
	// because browsers can't set custom headers on WS upgrade.
	r.GET("/api/ws", realtimeHandler.Connect)

	// Public webhook receiver — no auth on the route itself; each
	// provider's signature verification is the real auth boundary.
	// POST /webhooks/:provider routes to whatever was registered via
	// webhooks.Register(...) at app boot.
	r.POST("/webhooks/:provider", webhookHandler.Receive)

	// Public Grit UI component registry (shadcn-compatible)
	r.GET("/r.json", uiRegistryHandler.GetRegistry)
	r.GET("/r/:name", uiRegistryHandler.GetComponent)

	// Public blog routes (no auth required)
	blogs := r.Group("/api/blogs")
	{
		blogs.GET("", blogHandler.ListPublished)
		blogs.GET("/:slug", blogHandler.GetBySlug)
	}

	// Public auth routes
	auth := r.Group("/api/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
		auth.POST("/refresh", authHandler.Refresh)
		auth.POST("/forgot-password", authHandler.ForgotPassword)
		auth.POST("/reset-password", authHandler.ResetPassword)
	}

	// OAuth2 social login
	oauth := auth.Group("/oauth")
	{
		oauth.GET("/:provider", authHandler.OAuthBegin)
		oauth.GET("/:provider/callback", authHandler.OAuthCallback)
	}

	// TOTP verification (public — uses pending tokens, not JWT)
	auth.POST("/totp/verify", totpHandler.Verify)
	auth.POST("/totp/backup-codes/verify", totpHandler.VerifyBackupCode)

	// Protected routes
	protected := r.Group("/api")
	protected.Use(middleware.Auth(db, authService))
	// Activity logger writes one row per successful authenticated mutation.
	// Records who/what/when/where for audit. Read-only — see admin/activity.
	protected.Use(middleware.ActivityLogger(db))
	{
		protected.GET("/auth/me", authHandler.Me)
		protected.POST("/auth/logout", authHandler.Logout)

		// Two-Factor Authentication (TOTP)
		protected.POST("/auth/totp/setup", totpHandler.Setup)
		protected.POST("/auth/totp/enable", totpHandler.Enable)
		protected.POST("/auth/totp/disable", totpHandler.Disable)
		protected.GET("/auth/totp/status", totpHandler.Status)
		protected.POST("/auth/totp/backup-codes", totpHandler.RegenerateBackupCodes)
		protected.DELETE("/auth/totp/trusted-devices", totpHandler.RevokeTrustedDevices)

		// User routes (authenticated)
		protected.GET("/users/:id", userHandler.GetByID)

		// File uploads
		protected.POST("/uploads", uploadHandler.Create)
		protected.POST("/uploads/presign", uploadHandler.Presign)
		protected.POST("/uploads/complete", uploadHandler.CompleteUpload)
		protected.GET("/uploads", uploadHandler.List)
		protected.GET("/uploads/:id", uploadHandler.GetByID)
		protected.DELETE("/uploads/:id", uploadHandler.Delete)

		// Offline-first sync — desktop clients call these to flush their
		// local outbox and pull server-side updates.
		protected.POST("/sync/push", syncHandler.Push)
		protected.GET("/sync/pull", syncHandler.Pull)

		// AI
		protected.POST("/ai/complete", aiHandler.Complete)
		protected.POST("/ai/chat", aiHandler.Chat)
		protected.POST("/ai/stream", aiHandler.Stream)

		// Grit UI component registry (authenticated browse)
		protected.GET("/ui-components", uiRegistryHandler.ListComponents)
		protected.GET("/ui-components/:name", uiRegistryHandler.GetComponentDetail)

		// In-app notification bell — every authenticated user. Pulls
		// from a single Notification table that the SecObs poller
		// writes into when Sentinel/Pulse fires a high-severity event.
		protected.GET("/notifications", notificationHandler.List)
		protected.POST("/notifications/:id/read", notificationHandler.MarkRead)
		protected.POST("/notifications/read-all", notificationHandler.MarkAllRead)

		// grit:routes:protected
	}

	// Profile routes (any authenticated user)
	profile := protected.Group("/profile")
	{
		profile.GET("", userHandler.GetProfile)
		profile.PUT("", userHandler.UpdateProfile)
		profile.DELETE("", userHandler.DeleteProfile)
	}

	// Admin routes
	admin := r.Group("/api")
	admin.Use(middleware.Auth(db, authService))
	admin.Use(middleware.RequireRole("ADMIN"))
	{
		admin.GET("/users", userHandler.List)
		admin.POST("/users", userHandler.Create)
		admin.PUT("/users/:id", userHandler.Update)
		admin.DELETE("/users/:id", userHandler.Delete)

		// Activity audit log + tamper-evident chain verification
		admin.GET("/admin/activity", activityHandler.List)
		admin.GET("/admin/activity/integrity", activityHandler.VerifyIntegrity)

		// Webhook receiver admin (review + replay failed events)
		admin.GET("/admin/webhooks", webhookHandler.List)
		admin.POST("/admin/webhooks/:id/replay", webhookHandler.Replay)

		// Feature flags + A/B testing
		admin.GET("/admin/flags", featureFlagHandler.List)
		admin.POST("/admin/flags", featureFlagHandler.Create)
		admin.PUT("/admin/flags/:id", featureFlagHandler.Update)
		admin.DELETE("/admin/flags/:id", featureFlagHandler.Delete)
		admin.GET("/admin/flags/:id/exposures", featureFlagHandler.Exposures)

		// Admin system routes
		admin.GET("/admin/jobs/stats", jobsHandler.Stats)
		admin.GET("/admin/jobs/:status", jobsHandler.ListByStatus)
		admin.POST("/admin/jobs/:id/retry", jobsHandler.Retry)
		admin.DELETE("/admin/jobs/queue/:queue", jobsHandler.ClearQueue)
		admin.GET("/admin/cron/tasks", cronHandler.ListTasks)

		// Blog management (admin)
		admin.GET("/admin/blogs", blogHandler.List)
		admin.POST("/admin/blogs", blogHandler.Create)
		admin.PUT("/admin/blogs/:id", blogHandler.Update)
		admin.DELETE("/admin/blogs/:id", blogHandler.Delete)

		// Grit UI component registry (admin management)
		admin.POST("/admin/ui-components", uiRegistryHandler.CreateComponent)
		admin.PUT("/admin/ui-components/:name", uiRegistryHandler.UpdateComponent)
		admin.DELETE("/admin/ui-components/:name", uiRegistryHandler.DeleteComponent)

		// In-app Security dashboard — aggregates Sentinel APIs into one
		// envelope so the React page does a single round-trip. Operators
		// who want to dig deeper open /sentinel/ui directly.
		admin.GET("/admin/security/summary", securityHandler.Summary)
		// In-app Observability dashboard — same pattern against Pulse.
		// Operators who want a flame graph or the full SLO timeline open
		// /pulse/ui directly.
		admin.GET("/admin/observability/summary", observabilityHandler.Summary)

		// grit:routes:admin
	}

	// Custom role-restricted routes
	// grit:routes:custom

	return r
}
`
}
