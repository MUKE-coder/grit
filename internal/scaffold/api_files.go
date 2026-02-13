package scaffold

import (
	"fmt"
	"path/filepath"
	"strings"
)

func writeAPIFiles(root string, opts Options) error {
	apiRoot := filepath.Join(root, "apps", "api")
	module := opts.ProjectName + "/apps/api"

	files := map[string]string{
		filepath.Join(apiRoot, "go.mod"):                              apiGoMod(opts),
		filepath.Join(apiRoot, ".gitignore"):                          apiGitignore(),
		filepath.Join(apiRoot, "cmd", "server", "main.go"):           apiMainGo(opts),
		filepath.Join(apiRoot, "internal", "config", "config.go"):    apiConfigGo(),
		filepath.Join(apiRoot, "internal", "database", "database.go"): apiDatabaseGo(),
		filepath.Join(apiRoot, "internal", "models", "user.go"):      apiUserModelGo(),
		filepath.Join(apiRoot, "internal", "models", "upload.go"):    apiUploadModelGo(),
		filepath.Join(apiRoot, "internal", "services", "auth.go"):    apiAuthServiceGo(),
		filepath.Join(apiRoot, "internal", "handlers", "auth.go"):    apiAuthHandlerGo(),
		filepath.Join(apiRoot, "internal", "handlers", "user.go"):    apiUserHandlerGo(),
		filepath.Join(apiRoot, "internal", "middleware", "auth.go"):  apiAuthMiddlewareGo(),
		filepath.Join(apiRoot, "internal", "middleware", "cors.go"):  apiCorsMiddlewareGo(),
		filepath.Join(apiRoot, "internal", "middleware", "logger.go"): apiLoggerMiddlewareGo(),
		filepath.Join(apiRoot, "internal", "routes", "routes.go"):    apiRoutesGo(),
		filepath.Join(apiRoot, "internal", "docs", "docs.go"):       apiScalarHandlerGo(opts),
		filepath.Join(apiRoot, "internal", "docs", "openapi.go"):    apiOpenAPISpecGo(opts),
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
	return fmt.Sprintf(`module %s/apps/api

go 1.21

require (
	github.com/MarceloPetrucio/go-scalar-api-reference v0.0.0-20240521013641-ce5d2efe0e06
	github.com/MUKE-coder/gorm-studio v0.0.0-20260210080608-2fb07b651336
	github.com/aws/aws-sdk-go-v2 v1.25.0
	github.com/aws/aws-sdk-go-v2/config v1.27.0
	github.com/aws/aws-sdk-go-v2/credentials v1.17.0
	github.com/aws/aws-sdk-go-v2/feature/s3/manager v1.16.0
	github.com/aws/aws-sdk-go-v2/service/s3 v1.51.0
	github.com/disintegration/imaging v1.6.2
	github.com/gin-gonic/gin v1.10.0
	github.com/golang-jwt/jwt/v5 v5.2.0
	github.com/hibiken/asynq v0.24.1
	github.com/joho/godotenv v1.5.1
	github.com/redis/go-redis/v9 v9.4.0
	golang.org/x/crypto v0.23.0
	gorm.io/driver/postgres v1.5.11
	gorm.io/gorm v1.25.12
)
`, opts.ProjectName)
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

	// AI service
	var aiService *ai.AI
	if cfg.AIAPIKey != "" {
		aiService = ai.New(cfg.AIProvider, cfg.AIAPIKey, cfg.AIModel)
		log.Printf("AI service configured (%s)", cfg.AIProvider)
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

	GORMStudioEnabled bool

	// AI
	AIProvider string // "claude", "openai", or "gemini"
	AIAPIKey   string
	AIModel    string
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

		GORMStudioEnabled: getEnv("GORM_STUDIO_ENABLED", "true") == "true",

		AIProvider: getEnv("AI_PROVIDER", "claude"),
		AIAPIKey:   getEnv("AI_API_KEY", ""),
		AIModel:    getEnv("AI_MODEL", "claude-sonnet-4-5-20250929"),
	}

	if cfg.DatabaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL is required")
	}

	if cfg.JWTSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET is required")
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

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Connect establishes a database connection using the provided DSN.
func Connect(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true, // Avoids prepared statement issues with schema changes
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
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

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Role constants
const (
	RoleAdmin  = "ADMIN"
	RoleEditor = "EDITOR"
	RoleUser   = "USER"
)

// User represents a user in the system.
type User struct {
	ID              uint           ` + "`" + `gorm:"primarykey" json:"id"` + "`" + `
	FirstName       string         ` + "`" + `gorm:"size:255;not null" json:"first_name" binding:"required"` + "`" + `
	LastName        string         ` + "`" + `gorm:"size:255;not null" json:"last_name" binding:"required"` + "`" + `
	Email           string         ` + "`" + `gorm:"size:255;uniqueIndex;not null" json:"email" binding:"required,email"` + "`" + `
	Password        string         ` + "`" + `gorm:"size:255;not null" json:"-"` + "`" + `
	Role            string         ` + "`" + `gorm:"size:20;default:USER" json:"role"` + "`" + `
	Avatar          string         ` + "`" + `gorm:"size:500" json:"avatar"` + "`" + `
	JobTitle        string         ` + "`" + `gorm:"size:255" json:"job_title"` + "`" + `
	Bio             string         ` + "`" + `gorm:"type:text" json:"bio"` + "`" + `
	Active          bool           ` + "`" + `gorm:"default:true" json:"active"` + "`" + `
	EmailVerifiedAt *time.Time     ` + "`" + `json:"email_verified_at"` + "`" + `
	CreatedAt       time.Time      ` + "`" + `json:"created_at"` + "`" + `
	UpdatedAt       time.Time      ` + "`" + `json:"updated_at"` + "`" + `
	DeletedAt       gorm.DeletedAt ` + "`" + `gorm:"index" json:"-"` + "`" + `
}

// BeforeCreate hashes the password before saving.
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		u.Password = string(hashedPassword)
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
		// grit:models
	}
}

// Migrate runs database migrations only for tables that don't exist yet.
// It prints which tables were created and which were skipped.
func Migrate(db *gorm.DB) error {
	models := Models()
	migrated := 0

	for _, model := range models {
		if db.Migrator().HasTable(model) {
			log.Printf("  ✓ %T — already exists, skipping", model)
			continue
		}

		if err := db.AutoMigrate(model); err != nil {
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
	ID           uint           ` + "`" + `gorm:"primarykey" json:"id"` + "`" + `
	Filename     string         ` + "`" + `gorm:"size:255;not null" json:"filename"` + "`" + `
	OriginalName string         ` + "`" + `gorm:"size:255;not null" json:"original_name"` + "`" + `
	MimeType     string         ` + "`" + `gorm:"size:100;not null" json:"mime_type"` + "`" + `
	Size         int64          ` + "`" + `gorm:"not null" json:"size"` + "`" + `
	Path         string         ` + "`" + `gorm:"size:500;not null" json:"path"` + "`" + `
	URL          string         ` + "`" + `gorm:"size:500" json:"url"` + "`" + `
	ThumbnailURL string         ` + "`" + `gorm:"size:500" json:"thumbnail_url"` + "`" + `
	UserID       uint           ` + "`" + `gorm:"index;not null" json:"user_id"` + "`" + `
	User         User           ` + "`" + `gorm:"foreignKey:UserID" json:"-"` + "`" + `
	CreatedAt    time.Time      ` + "`" + `json:"created_at"` + "`" + `
	UpdatedAt    time.Time      ` + "`" + `json:"updated_at"` + "`" + `
	DeletedAt    gorm.DeletedAt ` + "`" + `gorm:"index" json:"-"` + "`" + `
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
	UserID uint   ` + "`" + `json:"user_id"` + "`" + `
	Email  string ` + "`" + `json:"email"` + "`" + `
	Role   string ` + "`" + `json:"role"` + "`" + `
	jwt.RegisteredClaims
}

// GenerateTokenPair creates a new access + refresh token pair.
func (s *AuthService) GenerateTokenPair(userID uint, email, role string) (*TokenPair, error) {
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

func (s *AuthService) generateToken(userID uint, email, role string, expiry time.Duration) (string, int64, error) {
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
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"golang.org/x/crypto/bcrypt"

	"` + "{{MODULE}}" + `/internal/models"
	"` + "{{MODULE}}" + `/internal/services"
)

// AuthHandler handles authentication endpoints.
type AuthHandler struct {
	DB          *gorm.DB
	AuthService *services.AuthService
}

type registerRequest struct {
	FirstName string ` + "`" + `json:"first_name" binding:"required,min=2"` + "`" + `
	LastName  string ` + "`" + `json:"last_name" binding:"required,min=2"` + "`" + `
	Email     string ` + "`" + `json:"email" binding:"required,email"` + "`" + `
	Password  string ` + "`" + `json:"password" binding:"required,min=8"` + "`" + `
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
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Password:  req.Password,
		Role:      models.RoleUser,
		Active:    true,
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

	if !user.CheckPassword(req.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": gin.H{
				"code":    "INVALID_CREDENTIALS",
				"message": "Invalid email or password",
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
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// Logger creates a structured JSON logging middleware.
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()
		method := c.Request.Method
		clientIP := c.ClientIP()

		if query != "" {
			path = path + "?" + query
		}

		log.Printf("[%d] %s %s | %s | %v",
			status,
			method,
			path,
			clientIP,
			latency,
		)
	}
}
`
}

func apiRoutesGo() string {
	return `package routes

import (
	"log"
	"net/http"

	"github.com/MUKE-coder/gorm-studio/studio"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"` + "{{MODULE}}" + `/internal/ai"
	"` + "{{MODULE}}" + `/internal/cache"
	"` + "{{MODULE}}" + `/internal/config"
	"` + "{{MODULE}}" + `/internal/docs"
	"` + "{{MODULE}}" + `/internal/handlers"
	"` + "{{MODULE}}" + `/internal/mail"
	"` + "{{MODULE}}" + `/internal/middleware"
	"` + "{{MODULE}}" + `/internal/models"
	"` + "{{MODULE}}" + `/internal/jobs"
	"` + "{{MODULE}}" + `/internal/services"
	"` + "{{MODULE}}" + `/internal/storage"
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
	r.Use(middleware.Logger())
	r.Use(gin.Recovery())
	r.Use(middleware.CORS(cfg.CORSOrigins))

	// Mount GORM Studio
	if cfg.GORMStudioEnabled {
		studio.Mount(r, db, []interface{}{&models.User{}, &models.Upload{}, /* grit:studio */}, studio.Config{
			Prefix: "/studio",
		})
		log.Println("GORM Studio mounted at /studio")
	}

	// API Documentation (Scalar)
	docs.RegisterRoutes(r)
	log.Println("API docs available at /docs")

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
	// grit:handlers

	// Health check
	r.GET("/api/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"version": "0.1.0",
		})
	})

	// Public auth routes
	auth := r.Group("/api/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
		auth.POST("/refresh", authHandler.Refresh)
		auth.POST("/forgot-password", authHandler.ForgotPassword)
		auth.POST("/reset-password", authHandler.ResetPassword)
	}

	// Protected routes
	protected := r.Group("/api")
	protected.Use(middleware.Auth(db, authService))
	{
		protected.GET("/auth/me", authHandler.Me)
		protected.POST("/auth/logout", authHandler.Logout)

		// User routes (authenticated)
		protected.GET("/users/:id", userHandler.GetByID)

		// File uploads
		protected.POST("/uploads", uploadHandler.Create)
		protected.GET("/uploads", uploadHandler.List)
		protected.GET("/uploads/:id", uploadHandler.GetByID)
		protected.DELETE("/uploads/:id", uploadHandler.Delete)

		// AI
		protected.POST("/ai/complete", aiHandler.Complete)
		protected.POST("/ai/chat", aiHandler.Chat)
		protected.POST("/ai/stream", aiHandler.Stream)

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

		// Admin system routes
		admin.GET("/admin/jobs/stats", jobsHandler.Stats)
		admin.GET("/admin/jobs/:status", jobsHandler.ListByStatus)
		admin.POST("/admin/jobs/:id/retry", jobsHandler.Retry)
		admin.DELETE("/admin/jobs/queue/:queue", jobsHandler.ClearQueue)
		admin.GET("/admin/cron/tasks", cronHandler.ListTasks)

		// grit:routes:admin
	}

	// Custom role-restricted routes
	// grit:routes:custom

	return r
}
`
}
