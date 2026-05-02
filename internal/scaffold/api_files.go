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
	github.com/MUKE-coder/pulse v0.0.0-20260223005903-6f5d6e356231
	github.com/aws/aws-sdk-go-v2 v1.25.0
	github.com/aws/aws-sdk-go-v2/config v1.27.0
	github.com/aws/aws-sdk-go-v2/credentials v1.17.0
	github.com/aws/aws-sdk-go-v2/feature/s3/manager v1.16.0
	github.com/aws/aws-sdk-go-v2/service/s3 v1.51.0
	github.com/disintegration/imaging v1.6.2
	github.com/gin-gonic/gin v1.10.0
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
	github.com/MUKE-coder/sentinel v0.0.0-20260220061042-2d2324be6824
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
	svc := &routes.Services{
		Cache:   cacheService,
		Storage: storageService,
		Mailer:  mailer,
		AI:      aiService,
		Jobs:    jobClient,
	}

	// Setup router
	router := routes.Setup(db, cfg, svc)

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
	StorageDriver string        // "minio", "r2", or "b2"
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
	SentinelEnabled   bool
	SentinelUsername   string
	SentinelPassword  string
	SentinelSecretKey string

	// Observability (Pulse)
	PulseEnabled   bool
	PulseUsername   string
	PulsePassword  string

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
		DatabaseURL: getEnv("DATABASE_URL", ""),
		JWTSecret:   getEnv("JWT_SECRET", ""),
		RedisURL:    getEnv("REDIS_URL", "redis://localhost:6379"),

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

		SentinelEnabled:   getEnv("SENTINEL_ENABLED", "true") == "true",
		SentinelUsername:   getEnv("SENTINEL_USERNAME", "admin"),
		SentinelPassword:  getEnv("SENTINEL_PASSWORD", "sentinel"),
		SentinelSecretKey: getEnv("SENTINEL_SECRET_KEY", "sentinel-secret-change-me"),

		PulseEnabled:  getEnv("PULSE_ENABLED", "true") == "true",
		PulseUsername: getEnv("PULSE_USERNAME", "admin"),
		PulsePassword: getEnv("PULSE_PASSWORD", "pulse"),

		GoogleClientID:     getEnv("GOOGLE_CLIENT_ID", ""),
		GoogleClientSecret: getEnv("GOOGLE_CLIENT_SECRET", ""),
		GithubClientID:     getEnv("GITHUB_CLIENT_ID", ""),
		GithubClientSecret: getEnv("GITHUB_CLIENT_SECRET", ""),
		OAuthFrontendURL:   getEnv("OAUTH_FRONTEND_URL", "http://localhost:3001"),
	}

	if cfg.DatabaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL is required")
	}

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

// resolveStorage returns the StorageConfig for the active driver.
func resolveStorage(driver string) StorageConfig {
	switch driver {
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
			Endpoint:  getEnv("MINIO_ENDPOINT", "http://localhost:9000"),
			AccessKey: getEnv("MINIO_ACCESS_KEY", "minioadmin"),
			SecretKey: getEnv("MINIO_SECRET_KEY", "minioadmin"),
			Bucket:    getEnv("MINIO_BUCKET", "uploads"),
			Region:    getEnv("MINIO_REGION", "us-east-1"),
			UseSSL:    getEnv("MINIO_USE_SSL", "false") == "true",
		}
	}
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
`
}

func apiDatabaseGo() string {
	return `package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Connect establishes a database connection using the provided DSN.
func Connect(dsn string) (*gorm.DB, error) {
	// Use Warn level by default — only logs slow queries and errors.
	// Set DB_LOG_LEVEL=info for verbose SQL logging during debugging.
	// In production, Warn prevents log flooding during AutoMigrate.
	logLevel := logger.Warn
	if os.Getenv("DB_LOG_LEVEL") == "info" {
		logLevel = logger.Info
	} else if os.Getenv("DB_LOG_LEVEL") == "silent" {
		logLevel = logger.Silent
	}

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true, // Avoids prepared statement issues with schema changes
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	// Connection pool settings
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
		// grit:models
	}
}

// Migrate runs database migrations only for tables that don't exist yet.
// It prints which tables were created and which were skipped.
func Migrate(db *gorm.DB) error {
	models := Models()
	migrated := 0

	// Use Silent logger during migration to suppress schema inspection SQL noise
	silentDB := db.Session(&gorm.Session{Logger: logger.Default.LogMode(logger.Silent)})

	for _, model := range models {
		if silentDB.Migrator().HasTable(model) {
			log.Printf("  ✓ %T — already exists, skipping", model)
			continue
		}

		if err := silentDB.AutoMigrate(model); err != nil {
			return fmt.Errorf("migrating %T: %w", model, err)
		}
		log.Printf("  ✓ %T — created", model)
		migrated++
	}

	if migrated == 0 {
		log.Println("All tables are up to date — nothing to migrate.")
	} else {
		log.Printf("Migrated %d table(s).", migrated)
	}

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
	"time"

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

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"user":   user,
			"tokens": tokens,
		},
		"message": "Logged in successfully",
	})
}

// Refresh generates a new access token from a refresh token.
func (h *AuthHandler) Refresh(c *gin.Context) {
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

	claims, err := h.AuthService.ValidateToken(req.RefreshToken)
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

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"tokens": tokens,
		},
		"message": "Token refreshed successfully",
	})
}

// Logout invalidates the user's session.
func (h *AuthHandler) Logout(c *gin.Context) {
	// In a production system, you would blacklist the refresh token in Redis
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
	if err := h.DB.First(&user, id).Error; err != nil {
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
	if err := h.DB.First(&user, id).Error; err != nil {
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
	h.DB.First(&user, id)

	c.JSON(http.StatusOK, gin.H{
		"data":    user,
		"message": "User updated successfully",
	})
}

// Delete soft-deletes a user.
func (h *UserHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	var user models.User
	if err := h.DB.First(&user, id).Error; err != nil {
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
	if err := h.DB.First(&user, userID).Error; err != nil {
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
	if err := h.DB.First(&user, userID).Error; err != nil {
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

	h.DB.First(&user, userID)

	c.JSON(http.StatusOK, gin.H{
		"data":    user,
		"message": "Profile updated successfully",
	})
}

// DeleteProfile soft-deletes the currently authenticated user's account.
func (h *UserHandler) DeleteProfile(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var user models.User
	if err := h.DB.First(&user, userID).Error; err != nil {
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
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": gin.H{
					"code":    "UNAUTHORIZED",
					"message": "Authorization header is required",
				},
			})
			c.Abort()
			return
		}

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

		claims, err := authService.ValidateToken(parts[1])
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

		// Load user from database
		var user models.User
		if err := db.First(&user, claims.UserID).Error; err != nil {
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

// SecurityHeaders adds production security headers to all responses.
func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Header("Permissions-Policy", "camera=(), microphone=(), geolocation=()")
		// HSTS only in production (when behind HTTPS)
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
	"math"
	"strconv"

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
type Params struct {
	Page      int
	PageSize  int
	Search    string
	SortBy    string
	SortOrder string
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
}

// Meta is the pagination envelope, matching Grit's existing response shape.
type Meta struct {
	Total    int64 ` + "`" + `json:"total"` + "`" + `
	Page     int   ` + "`" + `json:"page"` + "`" + `
	PageSize int   ` + "`" + `json:"page_size"` + "`" + `
	Pages    int   ` + "`" + `json:"pages"` + "`" + `
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

// ActivityLog records every successful authenticated mutation. The
// payload digest is a SHA-256 of the request body so we have evidence
// of what was sent without storing PII verbatim. Read-only — no
// updates, no deletes (use a separate retention job to prune old rows).
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

func apiActivityMiddlewareGo() string {
	return `package middleware

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

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
// Insert is fire-and-forget (goroutine + new context) so the response
// path isn't blocked by DB latency. If the DB is down, the audit log
// drops the entry; that's acceptable — we'd rather serve users than
// fail-closed on a logging dependency.
func ActivityLogger(db *gorm.DB) gin.HandlerFunc {
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
		}
		// Fire-and-forget — log failure inside the goroutine, don't
		// block the response.
		go func(e models.ActivityLog) {
			_ = db.Create(&e).Error
		}(entry)
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
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

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

// CSV writes items as a comma-separated stream into w. Sets the right
// Content-Type/Content-Disposition headers if w is an http.ResponseWriter.
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
	"` + "{{MODULE}}" + `/internal/sync"
)

// Services holds all Phase 4 services for dependency injection.
type Services struct {
	Cache   *cache.Cache
	Storage *storage.Storage
	Mailer  *mail.Mailer
	AI      *ai.AI
	Jobs    *jobs.Client
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

		sentinel.Mount(r, db, sentinel.Config{
			Dashboard: sentinel.DashboardConfig{
				Username:  cfg.SentinelUsername,
				Password:  cfg.SentinelPassword,
				SecretKey: cfg.SentinelSecretKey,
			},
			WAF: sentinel.WAFConfig{
				Enabled: true,
				Mode: func() sentinel.WAFMode {
					if isDev { return sentinel.ModeLog }
					return sentinel.ModeBlock
				}(),
			},
			RateLimit: sentinel.RateLimitConfig{
				Enabled: !isDev, // Disabled in development, enabled in production
				ByIP:    ipLimit,
				ByRoute: routeLimits,
			},
			AuthShield: sentinel.AuthShieldConfig{
				Enabled:    !isDev, // Disabled in development
				LoginRoute: "/api/auth/login",
			},
			Anomaly: sentinel.AnomalyConfig{
				Enabled: !isDev, // Disabled in development
			},
			Geo: sentinel.GeoConfig{
				Enabled: !isDev, // Disabled in development
			},
		})
		log.Println("Sentinel security suite mounted at /sentinel")
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
		p := pulse.Mount(r, db, pulse.Config{
			AppName: cfg.AppName,
			DevMode: cfg.IsDevelopment(),
			Dashboard: pulse.DashboardConfig{
				Username: cfg.PulseUsername,
				Password: cfg.PulsePassword,
			},
			Tracing: pulse.TracingConfig{
				ExcludePaths: []string{"/studio/*", "/sentinel/*", "/docs/*", "/pulse/*"},
			},
			Alerts: pulse.AlertConfig{},
			Prometheus: pulse.PrometheusConfig{
				Enabled: true,
			},
		})

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
	realtimeHub := realtime.NewHub()
	realtimeHandler := handlers.NewRealtimeHandler(realtimeHub, authService)
	_ = realtimeHub // available to handlers/services that want to push events

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

		// Activity audit log
		admin.GET("/admin/activity", activityHandler.List)

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

		// grit:routes:admin
	}

	// Custom role-restricted routes
	// grit:routes:custom

	return r
}
`
}
