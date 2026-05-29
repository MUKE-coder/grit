package config

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
	FrontendURL string
	DatabaseURL string

	JWTSecret        string
	JWTAccessExpiry  time.Duration
	JWTRefreshExpiry time.Duration

	InternalAPIKey string

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

	// Security (Sentinel v2.0)
	SentinelEnabled        bool
	SentinelUsername       string
	SentinelPassword       string
	SentinelSecretKey      string
	// Empty by default — XFF is ignored. Populate when behind a known proxy.
	SentinelTrustedProxies []string

	// Observability (Pulse v1.0)
	PulseEnabled    bool
	PulseUsername    string
	PulsePassword   string
	// "memory" (default) or "sqlite" for the v1.0 persistent backend.
	PulseStorage    string
	PulseStorageDSN string

	// Demo mode — when true, the seeder reseeds the demo cohort,
	// the mailer is short-circuited so no email goes out, and a cron
	// task resets the database every 24 h.
	DemoMode bool

	// OAuth2 Social Login
	GoogleClientID     string
	GoogleClientSecret string
	GithubClientID     string
	GithubClientSecret string
	OAuthFrontendURL   string // Where to redirect after OAuth callback

	// DGateway (Desispay) — mobile money collections for loan repayments
	DGatewayAPIKey         string
	DGatewayBaseURL        string
	DGatewayDefaultProvider string // iotec | relworx
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
		FrontendURL: getEnv("FRONTEND_URL", "http://localhost:5173"),
		DatabaseURL: getEnv("DATABASE_URL", ""),
		JWTSecret:   getEnv("JWT_SECRET", ""),
		InternalAPIKey: getEnv("INTERNAL_API_KEY", ""),
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

		DemoMode: getEnv("DEMO_MODE", "false") == "true",

		GoogleClientID:     getEnv("GOOGLE_CLIENT_ID", ""),
		GoogleClientSecret: getEnv("GOOGLE_CLIENT_SECRET", ""),
		GithubClientID:     getEnv("GITHUB_CLIENT_ID", ""),
		GithubClientSecret: getEnv("GITHUB_CLIENT_SECRET", ""),
		OAuthFrontendURL:   getEnv("OAUTH_FRONTEND_URL", "http://localhost:3001"),

		DGatewayAPIKey:          getEnv("DGATEWAY_API_KEY", ""),
		DGatewayBaseURL:         getEnv("DGATEWAY_BASE_URL", ""),
		DGatewayDefaultProvider: getEnv("DGATEWAY_DEFAULT_PROVIDER", "iotec"),
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

// splitCSV trims and splits a comma-separated env var. Empty entries
// after trimming are dropped.
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
