package config

import (
	"fmt"
	"log"
	"net"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"golang.org/x/net/publicsuffix"

	"library/apps/api/internal/crypto"
)

// StorageConfig holds credentials for a single S3-compatible provider.
type StorageConfig struct {
	Endpoint  string
	AccessKey string
	SecretKey string
	Bucket    string
	Region    string
	UseSSL    bool

	// PublicURL is the origin a BROWSER loads stored objects from, which is
	// not always the origin the SDK talks to.
	//
	// MinIO serves objects from the same host it takes API calls on, so this
	// can stay empty in development. R2 cannot: its S3 endpoint
	// (<account>.r2.cloudflarestorage.com) only answers SigV4-signed requests,
	// so an <img src> pointed at it gets a 401 — the upload succeeds and
	// nothing ever renders, which looks like a CORS problem and is not one.
	// Set this to the bucket's public origin: an r2.dev subdomain, a custom
	// domain, or a CDN in front of S3.
	//
	// When set, object URLs become <PublicURL>/<key> — public origins are
	// already scoped to one bucket, so the bucket segment is not repeated.
	PublicURL string

	// LocalRoot is the directory STORAGE_DRIVER=local keeps files in
	// (STORAGE_LOCAL_ROOT). URLSecret signs its temporary URLs
	// (STORAGE_URL_SECRET, or JWT_SECRET when that is unset).
	LocalRoot string
	URLSecret string
}

// Config holds all application configuration.
// ModuleFlags switches optional batteries on and off.
//
// A disabled module mounts no routes, registers no workers or cron entries, and
// migrates no tables — so turning one off removes it from the running app and
// the database, not just from view. The code stays in the repo; delete it by
// hand if you want it gone entirely.
type ModuleFlags struct {
	AI        bool // /api/ai/* — chat + completion endpoints
	Jobs      bool // asynq background workers + the Jobs admin page
	Cron      bool // scheduled tasks
	Backup    bool // database backup/restore + the Data & Backup page
	Webhooks  bool // outbound webhook delivery
	Realtime  bool // WebSocket hub
	Files     bool // uploads + the File manager
	Mail      bool // transactional email
	Audit     bool // activity log
	Flags     bool // feature flags
	TwoFactor bool // TOTP / 2FA
}

// Enabled reports whether a module is on, by the name used in the API and the
// admin nav. Unknown names return false so a typo hides the feature rather than
// silently exposing it.
func (m ModuleFlags) Enabled(name string) bool {
	switch name {
	case "ai":
		return m.AI
	case "jobs":
		return m.Jobs
	case "cron":
		return m.Cron
	case "backup":
		return m.Backup
	case "webhooks":
		return m.Webhooks
	case "realtime":
		return m.Realtime
	case "files":
		return m.Files
	case "mail":
		return m.Mail
	case "audit":
		return m.Audit
	case "flags":
		return m.Flags
	case "twofactor":
		return m.TwoFactor
	}
	return false
}

// Map renders the flags for the /api/system/modules endpoint, which the admin
// uses to hide nav entries for modules that are off.
func (m ModuleFlags) Map() map[string]bool {
	return map[string]bool{
		"ai":        m.AI,
		"jobs":      m.Jobs,
		"cron":      m.Cron,
		"backup":    m.Backup,
		"webhooks":  m.Webhooks,
		"realtime":  m.Realtime,
		"files":     m.Files,
		"mail":      m.Mail,
		"audit":     m.Audit,
		"flags":     m.Flags,
		"twofactor": m.TwoFactor,
	}
}

type Config struct {
	AppName     string
	AppEnv      string
	Port        string
	AppURL      string
	DatabaseURL string

	JWTSecret        string
	JWTAccessExpiry  time.Duration
	JWTRefreshExpiry time.Duration

	// FieldEncryptionKey (base64, 32 bytes) enables transparent AES-256-GCM on
	// crypto.EncryptedString columns. Empty = disabled (values stored plaintext).
	FieldEncryptionKey string

	RedisURL string

	// Storage
	StorageDriver string        // "local", "minio", "s3", "r2", or "b2"
	Storage       StorageConfig // Resolved config for the active driver

	// StorageDisks are the named disks in STORAGE_DISKS, such as a private bucket
	// for backups. StoragePublicPrefixes are the key prefixes anyone may read
	// without a signed link (STORAGE_PUBLIC_PREFIXES).
	StorageDisks          []StorageDisk
	StoragePublicPrefixes []string

	ResendAPIKey string
	MailFrom     string

	// Mail picks the transport (MAIL_MAILER) and holds each driver's
	// settings. internal/mail.FromConfig reads it.
	Mail MailConfig

	CORSOrigins []string

	// UploadAllowedMIME are MIME types added to the upload allowlist, from
	// UPLOAD_ALLOWED_MIME (comma separated). They extend the baseline in
	// handlers/upload.go; a field narrows it with accepts.
	UploadAllowedMIME []string

	// Modules turns optional batteries off.
	//
	// Grit ships everything on purpose — the batteries are the point. But not
	// every app wants an AI endpoint or a job queue, and a module you aren't
	// using shouldn't mount routes, start workers, or create tables.
	//
	// All default to TRUE, so an existing app behaves exactly as before. Set
	// MODULE_<NAME>=false in .env to switch one off.
	Modules ModuleFlags

	GORMStudioEnabled  bool
	GORMStudioUsername string
	GORMStudioPassword string

	// Studio's write paths. The SQL editor runs any UPDATE or DELETE it is
	// given through a raw Exec that no GORM callback sees, so production
	// turns it off. internal/appendonly guards tables that must never
	// change even with it on.
	GORMStudioReadOnly   bool
	GORMStudioDisableSQL bool

	// AI (Vercel AI Gateway)
	AIGatewayAPIKey string
	AIGatewayModel  string
	AIGatewayURL    string

	// TOTP (Two-Factor Authentication)
	TOTPIssuer               string
	RequireEmailVerification bool
	LoginMaxAttempts         int
	LoginLockoutWindow       time.Duration

	// Security (Sentinel)
	SentinelEnabled   bool
	SentinelUsername  string
	SentinelPassword  string
	SentinelSecretKey string
	SentinelAuditKey  string
	// Sentinel v2.0 — CIDRs allowed to send X-Forwarded-For / X-Real-IP.
	// Empty (default) means "ignore those headers entirely" — safe when
	// the app speaks to the public internet directly; populate when
	// you're behind a known reverse proxy (Caddy/Traefik/Cloudflare).
	SentinelTrustedProxies []string
	// TrustedProxies are the reverse proxies (IPs or CIDRs) whose X-Forwarded-For,
	// X-Real-IP and X-Forwarded-Proto the API believes. Read from TRUSTED_PROXIES;
	// see trustedProxies for the default.
	TrustedProxies []string

	// Observability (Pulse v1.0)
	PulseEnabled  bool
	PulseUsername string
	PulsePassword string
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

	// Serve the API reference at /docs in production (API_DOCS_PUBLIC).
	APIDocsPublic bool
}

// Load reads configuration from environment variables.
func Load() (*Config, error) {
	// Load .env file (ignore error if not found — production uses real env vars)
	_ = godotenv.Load()
	_ = godotenv.Load("../../.env") // Load from project root when running from apps/api

	storageDriver := resolveStorageDriver()

	cfg := &Config{
		AppName: getEnv("APP_NAME", "grit-app"),
		// Production unless told otherwise: a server that forgot APP_ENV is strict.
		AppEnv: getEnv("APP_ENV", "production"),
		// PORT first: every platform that routes to a container injects it, and
		// a deploy that binds the wrong port goes green with nothing answering.
		// APP_PORT stays the one you set yourself, and wins locally because a
		// platform is not the thing setting PORT there.
		Port:               firstNonEmpty(os.Getenv("PORT"), os.Getenv("APP_PORT"), "8080"),
		AppURL:             getEnv("APP_URL", "http://localhost:8080"),
		DatabaseURL:        resolveDatabaseURL(),
		JWTSecret:          getEnv("JWT_SECRET", ""),
		FieldEncryptionKey: getEnv("FIELD_ENCRYPTION_KEY", ""),
		RedisURL:           resolveRedisURL(),

		StorageDriver: storageDriver,
		Storage:       resolveStorage(storageDriver),

		StorageDisks:          resolveStorageDisks(storageDriver),
		StoragePublicPrefixes: resolveStoragePublicPrefixes(),

		ResendAPIKey: getEnv("RESEND_API_KEY", ""),
		MailFrom:     getEnv("MAIL_FROM", "noreply@localhost"),
		Mail:         loadMailConfig(),

		// The Wails desktop webview is allowed by middleware.isWailsOrigin (it
		// matches the wails.localhost host on any port), so it needs no entry
		// here — its dev origin includes a configurable port.
		CORSOrigins: strings.Split(getEnv("CORS_ORIGINS", "http://localhost:3000,http://localhost:3001"), ","),

		UploadAllowedMIME: splitCSV(getEnv("UPLOAD_ALLOWED_MIME", "")),

		// Optional batteries. Default on, so nothing changes for an existing
		// app; set MODULE_<NAME>=false to switch one off.
		Modules: ModuleFlags{
			AI:        getEnv("MODULE_AI", "true") == "true",
			Jobs:      getEnv("MODULE_JOBS", "true") == "true",
			Cron:      getEnv("MODULE_CRON", "true") == "true",
			Backup:    getEnv("MODULE_BACKUP", "true") == "true",
			Webhooks:  getEnv("MODULE_WEBHOOKS", "true") == "true",
			Realtime:  getEnv("MODULE_REALTIME", "true") == "true",
			Files:     getEnv("MODULE_FILES", "true") == "true",
			Mail:      getEnv("MODULE_MAIL", "true") == "true",
			Audit:     getEnv("MODULE_AUDIT", "true") == "true",
			Flags:     getEnv("MODULE_FLAGS", "true") == "true",
			TwoFactor: getEnv("MODULE_TWOFACTOR", "true") == "true",
		},

		GORMStudioEnabled:  getEnv("GORM_STUDIO_ENABLED", "true") == "true",
		GORMStudioUsername: getEnv("GORM_STUDIO_USERNAME", "admin"),
		GORMStudioPassword: getEnv("GORM_STUDIO_PASSWORD", "studio"),

		GORMStudioReadOnly:   getEnv("GORM_STUDIO_READ_ONLY", "false") == "true",
		GORMStudioDisableSQL: getEnv("GORM_STUDIO_DISABLE_SQL", "false") == "true",

		AIGatewayAPIKey: getEnv("AI_GATEWAY_API_KEY", ""),
		AIGatewayModel:  getEnv("AI_GATEWAY_MODEL", "anthropic/claude-sonnet-4-6"),
		AIGatewayURL:    getEnv("AI_GATEWAY_URL", "https://ai-gateway.vercel.sh/v1"),

		TOTPIssuer: getEnv("TOTP_ISSUER", getEnv("APP_NAME", "grit-app")),
		// Off by default and deliberately so: switching it on for an existing
		// project would lock out every user at once, because they all have a
		// NULL email_verified_at.
		RequireEmailVerification: getEnv("REQUIRE_EMAIL_VERIFICATION", "false") == "true",
		LoginMaxAttempts:         getEnvInt("LOGIN_MAX_ATTEMPTS", 10),
		LoginLockoutWindow:       getEnvDuration("LOGIN_LOCKOUT_MINUTES", 15, time.Minute),

		SentinelEnabled:        getEnv("SENTINEL_ENABLED", "true") == "true",
		SentinelUsername:       getEnv("SENTINEL_USERNAME", "admin"),
		SentinelPassword:       getEnv("SENTINEL_PASSWORD", "sentinel"),
		SentinelSecretKey:      getEnv("SENTINEL_SECRET_KEY", "sentinel-secret-change-me"),
		SentinelAuditKey:       getEnv("SENTINEL_AUDIT_KEY", ""),
		SentinelTrustedProxies: splitCSV(getEnv("SENTINEL_TRUSTED_PROXIES", strings.Join(trustedProxies(), ","))),
		TrustedProxies:         trustedProxies(),

		PulseEnabled:    getEnv("PULSE_ENABLED", "true") == "true",
		PulseUsername:   getEnv("PULSE_USERNAME", "admin"),
		PulsePassword:   getEnv("PULSE_PASSWORD", "pulse"),
		PulseStorage:    getEnv("PULSE_STORAGE", "memory"),
		PulseStorageDSN: getEnv("PULSE_STORAGE_DSN", "pulse.db"),

		GoogleClientID:     getEnv("GOOGLE_CLIENT_ID", ""),
		GoogleClientSecret: getEnv("GOOGLE_CLIENT_SECRET", ""),
		GithubClientID:     getEnv("GITHUB_CLIENT_ID", ""),
		GithubClientSecret: getEnv("GITHUB_CLIENT_SECRET", ""),
		OAuthFrontendURL:   getEnv("OAUTH_FRONTEND_URL", "http://localhost:3001"),
	}

	// The local driver keeps files on this server's disk: a second replica cannot
	// see them, and a redeploy that replaces the container loses them. A
	// single-server deployment with the directory on a volume it backs up says so.
	if cfg.StorageDriver == "local" && cfg.AppEnv == "production" && getEnv("ALLOW_LOCAL_STORAGE_IN_PRODUCTION", "false") != "true" {
		return nil, fmt.Errorf("APP_ENV=production is using STORAGE_DRIVER=local: use minio, s3, r2 or b2, or set ALLOW_LOCAL_STORAGE_IN_PRODUCTION=true if STORAGE_LOCAL_ROOT is on a volume you back up")
	}

	// DatabaseURL is always populated by resolveDatabaseURL() — either from
	// the DATABASE_URL env var or built from POSTGRES_* parts. The actual
	// connection attempt in cmd/server/main.go will surface a useful error
	// if the resolved URL points at an unreachable database.

	// Whether the API reference at /docs is served in production. It maps every
	// route, admin, backup, SSO and GDPR ones included, with a console to call
	// them, so it is not published unless asked for.
	cfg.APIDocsPublic = getEnv("API_DOCS_PUBLIC", "false") == "true"

	if cfg.AppEnv == "production" {
		// GORM Studio is a browser SQL console over every table. A development
		// .env says GORM_STUDIO_ENABLED=true, and .env files get copied to servers
		// whole, so in production it takes its own switch, and is read-only with
		// no SQL editor even then.
		cfg.GORMStudioEnabled = getEnv("GORM_STUDIO_IN_PRODUCTION", "false") == "true"
		cfg.GORMStudioReadOnly = true
		cfg.GORMStudioDisableSQL = true

		// SQLite in production keeps the database inside the container, and a
		// redeploy replaces the container: every deploy began from an empty
		// database while the Postgres the production compose file starts sat
		// unused. A deployment that does want SQLite, on a volume it backs up,
		// says so.
		if strings.HasPrefix(cfg.DatabaseURL, "sqlite:") && getEnv("ALLOW_SQLITE_IN_PRODUCTION", "false") != "true" {
			return nil, fmt.Errorf("APP_ENV=production is using SQLite: set DB_PROVIDER=postgres, or ALLOW_SQLITE_IN_PRODUCTION=true if the database file is on a volume you back up")
		}
	}
	if cfg.AppEnv != "development" {
		if err := checkSecrets(cfg); err != nil {
			return nil, err
		}
	}

	if cfg.JWTSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET is required")
	}
	if len(cfg.JWTSecret) < 32 {
		log.Println("WARNING: JWT_SECRET should be at least 32 characters for security. Generate one with: openssl rand -hex 32")
	}

	// Configure field-level encryption. A malformed key fails fast — running
	// without the encryption you configured is worse than refusing to start.
	if err := crypto.InitFieldKey(cfg.FieldEncryptionKey); err != nil {
		return nil, err
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

	cfg.warnCrossSiteAuth()

	return cfg, nil
}

// warnCrossSiteAuth says so when the API and a frontend are on different sites,
// which breaks cookie authentication without breaking anything that logs.
//
// Sign-in returns 200 and sets the refresh cookie, and every request after it is
// a 401, because a browser will not send a SameSite=Lax cookie on a cross-site
// request. A site is the registrable domain, not the hostname: api.example.com
// and app.example.com are the same site, and two apps on a platform whose domain
// is on the Public Suffix List are not. laravel.cloud, vercel.app, onrender.com
// and up.railway.app are all on it, so api-x.laravel.cloud and web-x.laravel.cloud
// share no registrable domain at all and the cookies are never sent.
//
// A warning and not a refusal: the arrangement is correct behind a same-origin
// proxy that forwards /api to the API, and correct for a client holding the
// access token itself. What is never correct is finding out from a user who
// cannot stay signed in.
func (c *Config) warnCrossSiteAuth() {
	apiSite := registrableDomain(c.AppURL)
	if apiSite == "" {
		return
	}
	for _, origin := range c.CORSOrigins {
		site := registrableDomain(origin)
		if site == "" || site == apiSite {
			continue
		}
		log.Printf("WARNING: the API is on %s and CORS_ORIGINS has %s, which is a different site. "+
			"Browsers do not send SameSite=Lax cookies across sites, so sign-in will return 200 and "+
			"every request after it 401. Put the frontend and the API on one domain, proxy /api from "+
			"the frontend so the browser sees one origin, or have the client send the access token as "+
			"a bearer header. Platform domains like laravel.cloud and vercel.app are themselves public "+
			"suffixes, so two apps under one project are already cross-site.",
			apiSite, site)
	}
}

// registrableDomain is the site a URL belongs to: its public suffix plus one
// label. Empty for localhost, a bare IP, and anything unparseable, none of which
// this check can say anything useful about.
//
// The suffix list is the snapshot compiled into golang.org/x/net, so a platform
// domain added to the list since that release reads as an ordinary domain and
// the warning is not raised. It errs that way on purpose: a missed warning
// costs a search, and a wrong one costs trust in every other line at boot.
func registrableDomain(raw string) string {
	if raw == "" {
		return ""
	}
	if !strings.Contains(raw, "://") {
		raw = "https://" + raw
	}
	u, err := url.Parse(raw)
	if err != nil {
		return ""
	}
	host := u.Hostname()
	if host == "" || host == "localhost" || net.ParseIP(host) != nil {
		return ""
	}
	site, err := publicsuffix.EffectiveTLDPlusOne(host)
	if err != nil {
		return ""
	}
	return site
}

// weakSecrets are values a scaffold, a README or a tutorial once put in these
// settings.
var weakSecrets = map[string]bool{
	"studio": true, "sentinel": true, "pulse": true, "admin": true, "password": true, "secret": true,
	"change-me": true, "change_me": true, "change-me-in-prod": true, "sentinel-secret-change-me": true,
	"minioadmin": true,
}

// checkSecrets refuses to start anywhere but development with a secret that can
// be guessed: a JWT secret short enough to brute-force, which forges any
// account's token, or a dashboard password from a default.
func checkSecrets(cfg *Config) error {
	var weak []string
	if len(cfg.JWTSecret) < 32 || weakSecrets[strings.ToLower(cfg.JWTSecret)] {
		weak = append(weak, "JWT_SECRET")
	}
	for _, s := range []struct {
		name, value string
		used        bool
	}{
		{"GORM_STUDIO_PASSWORD", cfg.GORMStudioPassword, cfg.GORMStudioEnabled},
		{"SENTINEL_PASSWORD", cfg.SentinelPassword, cfg.SentinelEnabled},
		{"SENTINEL_SECRET_KEY", cfg.SentinelSecretKey, cfg.SentinelEnabled},
		{"PULSE_PASSWORD", cfg.PulsePassword, cfg.PulseEnabled},
		// MinIO's root password, when MinIO is the store. S3, R2 and B2 keys are
		// issued by the provider and are not guessable defaults.
		{"MINIO_SECRET_KEY", cfg.Storage.SecretKey, cfg.StorageDriver == "minio"},
	} {
		if s.used && (len(s.value) < 16 || weakSecrets[strings.ToLower(s.value)]) {
			weak = append(weak, s.name)
		}
	}
	if len(weak) > 0 {
		return fmt.Errorf("refusing to start with APP_ENV=%s, because these are missing, too short or a known default: %s. "+
			"Generate each with: openssl rand -hex 32. On a machine only you can reach, APP_ENV=development skips this check",
			cfg.AppEnv, strings.Join(weak, ", "))
	}
	return nil
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
//
// resolveRedisURL decides whether this process talks to Redis at all.
//
// It cannot use getEnv, because getEnv treats an empty value as "unset" and
// hands back the default. That made Redis impossible to turn off: setting
// REDIS_URL= in .env looked like it should disable it and silently did not, so
// the asynq worker and the cron scheduler started anyway, failed to dial, and
// retried in a tight loop. The result was a process burning CPU on reconnects
// with nothing in the logs but a wall of dial errors — on a box with no Redis,
// simply running the API cost real cycles.
//
// So the three cases are distinguished explicitly:
//
//	REDIS_URL unset      → the local default, which is what most dev setups want
//	REDIS_URL=           → no Redis. Cache, jobs, worker and cron all stay off.
//	REDIS_URL=redis://…  → use it
//
// The empty case is a deliberate configuration, not a mistake, so it says so
// once at boot rather than leaving someone to wonder why their jobs never run.
func resolveRedisURL() string {
	v, ok := os.LookupEnv("REDIS_URL")
	if !ok {
		// Built from the port docker-compose binds, so moving REDIS_PORT in
		// .env moves the connection with it. This used to be a fixed 6380,
		// so a project that moved its port kept dialling the old one, which
		// on a machine running a second Grit project is that project's
		// Redis: the two then shared a cache and a job queue.
		return fmt.Sprintf("redis://%s:%s", getEnv("REDIS_HOST", "localhost"), getEnv("REDIS_PORT", "6380"))
	}
	if strings.TrimSpace(v) == "" {
		log.Println("REDIS_URL is empty: cache, background jobs and cron are disabled")
		return ""
	}
	warnPortMismatch("REDIS_URL", v, "REDIS_PORT")
	return v
}

// resolveMinioEndpoint follows the same rule for MinIO: an explicit
// MINIO_ENDPOINT wins, and without one the endpoint follows MINIO_PORT.
func resolveMinioEndpoint() string {
	if v := os.Getenv("MINIO_ENDPOINT"); v != "" {
		warnPortMismatch("MINIO_ENDPOINT", v, "MINIO_PORT")
		return v
	}
	return fmt.Sprintf("http://%s:%s", getEnv("MINIO_HOST", "localhost"), getEnv("MINIO_PORT", "9002"))
}

// warnPortMismatch says so when a URL names localhost on a different port from
// the one docker-compose was told to publish this project's service on.
//
// That combination almost always means the URL is left over from before the
// port moved, and that it now reaches some other project's container. Nothing
// fails: the other Redis answers, the other MinIO stores the file. So it is
// said once, loudly, at boot.
func warnPortMismatch(urlVar, raw, portVar string) {
	want := os.Getenv(portVar)
	if want == "" {
		return
	}
	u, err := url.Parse(raw)
	if err != nil {
		return
	}
	host := u.Hostname()
	if host != "localhost" && host != "127.0.0.1" {
		return
	}
	if got := u.Port(); got != "" && got != want {
		log.Printf("WARNING: %s points at %s:%s, but %s is %s. This project's container is published on %s, so %s is probably another project's. Set %s to port %s, or remove it from .env to follow %s.",
			urlVar, host, got, portVar, want, want, got, urlVar, want, portVar)
	}
}

// resolveDatabaseURL builds the DSN from DB_PROVIDER and that provider's parts.
//
// DATABASE_URL still wins: it is the escape hatch for a managed database whose
// connection string carries options nothing here models (a Neon pooler, an RDS
// proxy, a TLS mode). When both are set and they disagree about the engine, that
// is said once at boot rather than silently resolved, because the parts below
// then describe a database nothing connects to.
//
// DB_PROVIDER is named rather than inferred because the engine is a decision, not
// a detail: before this, the only way to choose one was the prefix of
// DATABASE_URL, which meant a project had POSTGRES_* variables in .env and no
// place at all to put MySQL credentials.
func resolveDatabaseURL() string {
	provider := strings.ToLower(strings.TrimSpace(getEnv("DB_PROVIDER", "postgres")))

	if v := os.Getenv("DATABASE_URL"); v != "" {
		warnProviderMismatch(provider, v)
		return v
	}

	switch provider {
	case "postgres", "postgresql", "pg", "":
		user := getEnv("POSTGRES_USER", "grit")
		pass := getEnv("POSTGRES_PASSWORD", "grit")
		host := getEnv("POSTGRES_HOST", "localhost")
		port := getEnv("POSTGRES_PORT", "5432")
		db := getEnv("POSTGRES_DB", getEnv("APP_NAME", "grit-app"))
		return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
			user, pass, host, port, db, getEnv("POSTGRES_SSLMODE", "disable"))

	case "mysql", "mariadb":
		// go-sql-driver's own DSN shape, not a URL: Connect strips the prefix and
		// hands the rest over as it is.
		user := getEnv("MYSQL_USER", "grit")
		pass := getEnv("MYSQL_PASSWORD", "grit")
		host := getEnv("MYSQL_HOST", "localhost")
		port := getEnv("MYSQL_PORT", "3306")
		db := getEnv("MYSQL_DB", getEnv("APP_NAME", "grit-app"))
		return fmt.Sprintf("mysql:%s:%s@tcp(%s:%s)/%s", user, pass, host, port, db)

	case "sqlite", "sqlite3", "file":
		return "sqlite:" + getEnv("SQLITE_PATH", "./app.db")

	case "memory", ":memory:":
		// Shared cache, not a bare :memory:. GORM pools connections, and a bare
		// in-memory SQLite gives each connection its own empty database: the
		// migration runs on one, the first query lands on another, and the table
		// "does not exist" on a database that was just migrated.
		return "sqlite:file::memory:?cache=shared"

	default:
		log.Fatalf("DB_PROVIDER=%q is not one this app knows. Use postgres, mysql, sqlite or memory, or set DATABASE_URL directly.", provider)
		return ""
	}
}

// warnProviderMismatch says so when DATABASE_URL names a different engine from
// DB_PROVIDER.
//
// Nothing breaks: DATABASE_URL wins and the app runs on whatever it names. What
// misleads is everything else in .env, which now describes a database this
// process never opens.
func warnProviderMismatch(provider, dsn string) {
	engine := "postgres"
	switch {
	case strings.HasPrefix(dsn, "mysql:"):
		engine = "mysql"
	case strings.HasPrefix(dsn, "sqlite:"):
		engine = "sqlite"
	}
	normalised := map[string]string{
		"postgresql": "postgres", "pg": "postgres", "": "postgres",
		"mariadb": "mysql", "sqlite3": "sqlite", "file": "sqlite",
		"memory": "sqlite", ":memory:": "sqlite",
	}
	if n, ok := normalised[provider]; ok {
		provider = n
	}
	if provider != engine {
		log.Printf("WARNING: DB_PROVIDER is %s and DATABASE_URL points at %s. DATABASE_URL wins, so this app is running on %s and the %s settings in .env are not being used.",
			provider, engine, engine, strings.ToUpper(provider))
	}
}

// resolveStorageDriver picks the storage driver from STORAGE_DRIVER.
//
// Outside production, minio with no MINIO_ACCESS_KEY becomes local, so a new
// project stores uploads before Docker is running instead of answering each
// one with STORAGE_UNAVAILABLE. Production never falls back: a server that
// lost its credentials should say so, not start writing to its own disk.
func resolveStorageDriver() string {
	driver := getEnv("STORAGE_DRIVER", "minio")
	if driver == "minio" && getEnv("MINIO_ACCESS_KEY", "") == "" && getEnv("APP_ENV", "production") != "production" {
		log.Println("MinIO has no credentials (MINIO_ACCESS_KEY), so files are kept on the local disk. Set STORAGE_DRIVER=local to make that the choice, or set the MinIO credentials to use MinIO")
		return "local"
	}
	return driver
}

// StorageDisk is one named disk from STORAGE_DISKS.
type StorageDisk struct {
	Name    string
	Driver  string
	Storage StorageConfig
}

// resolveStorageDisks reads the named disks in STORAGE_DISKS, a comma
// separated list of names.
//
// Each name takes STORAGE_DISK_<NAME>_* settings (DRIVER, BUCKET, ENDPOINT,
// ACCESS_KEY, SECRET_KEY, REGION, PUBLIC_URL, ROOT), and whatever it leaves
// unset comes from its driver's own settings. A private bucket for backups on
// the same provider as the uploads needs one line more than its name:
//
//	STORAGE_DISKS=backups
//	STORAGE_DISK_BACKUPS_BUCKET=myapp-backups
func resolveStorageDisks(defaultDriver string) []StorageDisk {
	var disks []StorageDisk
	for _, name := range strings.Split(os.Getenv("STORAGE_DISKS"), ",") {
		name = strings.ToLower(strings.TrimSpace(name))
		if name == "" || name == "default" {
			continue
		}
		if strings.IndexFunc(name, func(r rune) bool {
			return (r < 'a' || r > 'z') && (r < '0' || r > '9') && r != '_' && r != '-'
		}) >= 0 {
			log.Printf("STORAGE_DISKS: %q is not a disk name (letters, digits, _ and -), so it was skipped", name)
			continue
		}
		prefix := "STORAGE_DISK_" + strings.ToUpper(strings.ReplaceAll(name, "-", "_")) + "_"
		driver := getEnv(prefix+"DRIVER", defaultDriver)
		disk := resolveStorage(driver)
		if driver == "local" {
			// Apart from the default disk's files, and served under its route.
			disk.LocalRoot = filepath.Join(filepath.Dir(disk.LocalRoot), name)
			disk.PublicURL = strings.TrimRight(disk.PublicURL, "/") + "/_disks/" + name
		}
		disk.Endpoint = getEnv(prefix+"ENDPOINT", disk.Endpoint)
		disk.AccessKey = getEnv(prefix+"ACCESS_KEY", disk.AccessKey)
		disk.SecretKey = getEnv(prefix+"SECRET_KEY", disk.SecretKey)
		disk.Bucket = getEnv(prefix+"BUCKET", disk.Bucket)
		disk.Region = getEnv(prefix+"REGION", disk.Region)
		disk.PublicURL = getEnv(prefix+"PUBLIC_URL", disk.PublicURL)
		disk.LocalRoot = getEnv(prefix+"ROOT", disk.LocalRoot)
		if driver == defaultDriver && driver != "local" && disk.Bucket == resolveStorage(defaultDriver).Bucket {
			log.Printf("STORAGE_DISKS: %q uses the same bucket as the default disk, set %sBUCKET to keep its files apart", name, prefix)
		}
		disks = append(disks, StorageDisk{Name: name, Driver: driver, Storage: disk})
	}
	return disks
}

// resolveStoragePublicPrefixes reads STORAGE_PUBLIC_PREFIXES: the key
// prefixes anyone may read without a signed link, comma separated. Every
// other key, backups included, is private.
func resolveStoragePublicPrefixes() []string {
	var prefixes []string
	for _, prefix := range strings.Split(getEnv("STORAGE_PUBLIC_PREFIXES", "uploads/,thumbnails/"), ",") {
		if prefix = strings.Trim(strings.TrimSpace(prefix), "/"); prefix != "" {
			prefixes = append(prefixes, prefix+"/")
		}
	}
	return prefixes
}

// resolveStorage returns the StorageConfig for the active driver.
//
// For AWS S3, leave S3_ENDPOINT empty — the AWS SDK will use the
// regional endpoint automatically (s3.<region>.amazonaws.com).
// Credentials fall back to the AWS standard env vars
// AWS_ACCESS_KEY_ID + AWS_SECRET_ACCESS_KEY if you don't set the S3_*
// variants, which is convenient when running on EC2 / ECS / Lambda
// with an IAM role and you'd rather not duplicate keys in .env. The bucket
// and endpoint fall back to AWS_BUCKET and AWS_ENDPOINT_URL for the same
// reason: they are what the AWS SDKs read, and what a managed bucket
// injects, so attaching one needs no configuration here at all.
func resolveStorage(driver string) StorageConfig {
	switch driver {
	case "s3":
		// Empty endpoint = AWS SDK uses the regional default
		// (s3.<region>.amazonaws.com). This also flips the client into
		// virtual-hosted style, which AWS requires for buckets created
		// after Sep 2020.
		return StorageConfig{
			Endpoint:  firstNonEmpty(os.Getenv("S3_ENDPOINT"), os.Getenv("AWS_ENDPOINT_URL")),
			AccessKey: firstNonEmpty(os.Getenv("S3_ACCESS_KEY"), os.Getenv("AWS_ACCESS_KEY_ID")),
			SecretKey: firstNonEmpty(os.Getenv("S3_SECRET_KEY"), os.Getenv("AWS_SECRET_ACCESS_KEY")),
			Bucket:    firstNonEmpty(os.Getenv("S3_BUCKET"), os.Getenv("AWS_BUCKET"), "uploads"),
			Region:    firstNonEmpty(os.Getenv("S3_REGION"), os.Getenv("AWS_REGION"), "us-east-1"),
			UseSSL:    true,
			PublicURL: firstNonEmpty(os.Getenv("S3_PUBLIC_URL"), os.Getenv("STORAGE_PUBLIC_URL")),
		}
	case "r2":
		return StorageConfig{
			Endpoint:  getEnv("R2_ENDPOINT", ""),
			AccessKey: getEnv("R2_ACCESS_KEY", ""),
			SecretKey: getEnv("R2_SECRET_KEY", ""),
			Bucket:    getEnv("R2_BUCKET", "uploads"),
			Region:    getEnv("R2_REGION", "auto"),
			UseSSL:    true,
			PublicURL: firstNonEmpty(os.Getenv("R2_PUBLIC_URL"), os.Getenv("STORAGE_PUBLIC_URL")),
		}
	case "b2":
		return StorageConfig{
			Endpoint:  getEnv("B2_ENDPOINT", ""),
			AccessKey: getEnv("B2_ACCESS_KEY", ""),
			SecretKey: getEnv("B2_SECRET_KEY", ""),
			Bucket:    getEnv("B2_BUCKET", "uploads"),
			Region:    getEnv("B2_REGION", "us-west-004"),
			UseSSL:    true,
			PublicURL: firstNonEmpty(os.Getenv("B2_PUBLIC_URL"), os.Getenv("STORAGE_PUBLIC_URL")),
		}
	case "local":
		return StorageConfig{
			LocalRoot: getEnv("STORAGE_LOCAL_ROOT", "storage/app"),
			URLSecret: firstNonEmpty(os.Getenv("STORAGE_URL_SECRET"), os.Getenv("JWT_SECRET")),
			// The API serves these files itself, from the route routes.Setup mounts.
			PublicURL: firstNonEmpty(os.Getenv("STORAGE_PUBLIC_URL"), strings.TrimRight(getEnv("APP_URL", "http://localhost:8080"), "/")+"/files"),
		}
	default: // minio
		return StorageConfig{
			Endpoint: resolveMinioEndpoint(),
			// No default: minioadmin/minioadmin is the whole bucket to anyone who
			// can reach MinIO. .env carries credentials generated per project.
			AccessKey: getEnv("MINIO_ACCESS_KEY", ""),
			SecretKey: getEnv("MINIO_SECRET_KEY", ""),
			Bucket:    getEnv("MINIO_BUCKET", "uploads"),
			Region:    getEnv("MINIO_REGION", "us-east-1"),
			UseSSL:    getEnv("MINIO_USE_SSL", "false") == "true",
			PublicURL: firstNonEmpty(os.Getenv("MINIO_PUBLIC_URL"), os.Getenv("STORAGE_PUBLIC_URL")),
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

// getEnvInt reads a whole-number env var. A malformed value falls back rather
// than failing the boot: an unparseable LOGIN_MAX_ATTEMPTS should not take the
// API down, and the fallback is the safe direction.
func getEnvInt(key string, fallback int) int {
	if val := os.Getenv(key); val != "" {
		if n, err := strconv.Atoi(val); err == nil {
			return n
		}
		log.Printf("config: %s=%q is not a number, using %d", key, val, fallback)
	}
	return fallback
}

// getEnvDuration reads a whole number of units, named by the caller, which
// keeps the env var name self-describing
// (LOGIN_LOCKOUT_MINUTES=15 rather than a duration string nobody formats
// consistently).
func getEnvDuration(key string, fallback int, unit time.Duration) time.Duration {
	return time.Duration(getEnvInt(key, fallback)) * unit
}

// defaultTrustedProxies are the addresses a reverse proxy on this machine or on
// a private network (a Docker network, a VPC) connects from. A client on the
// internet cannot connect from one, so the proxies of the documented deployments
// are believed and nobody else is.
var defaultTrustedProxies = []string{"127.0.0.0/8", "::1/128", "10.0.0.0/8", "172.16.0.0/12", "192.168.0.0/16", "fc00::/7"}

// trustedProxies reads TRUSTED_PROXIES: a comma-separated list of IPs and CIDRs,
// "none" to trust no proxy, or unset for defaultTrustedProxies. Behind a CDN
// that connects to the API directly (Cloudflare with no proxy of your own), list
// the CDN's ranges, or every request carries the CDN's address.
func trustedProxies() []string {
	raw := strings.TrimSpace(os.Getenv("TRUSTED_PROXIES"))
	if raw == "" {
		return append([]string(nil), defaultTrustedProxies...)
	}
	if strings.EqualFold(raw, "none") {
		return []string{}
	}
	return splitCSV(raw)
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

// MailConfig configures internal/mail. MAIL_MAILER picks the driver, and each
// driver reads only its own keys; .env.example lists them.
type MailConfig struct {
	Mailer               string   // MAIL_MAILER: smtp, resend, mailgun, postmark, sendgrid, ses, log or failover
	FromName             string   // MAIL_FROM_NAME
	Failover             []string // MAIL_FAILOVER, tried in order
	AllowLogInProduction bool     // MAIL_ALLOW_LOG_IN_PRODUCTION
	LogPath              string   // MAIL_LOG_PATH, where the log driver saves messages

	SMTPHost       string
	SMTPPort       string
	SMTPUsername   string
	SMTPPassword   string
	SMTPEncryption string // tls, starttls or none; empty picks from the host and port

	MailgunDomain   string
	MailgunSecret   string
	MailgunEndpoint string // api.mailgun.net, or api.eu.mailgun.net

	PostmarkToken         string
	PostmarkMessageStream string

	SendGridAPIKey string

	SESRegion          string
	SESAccessKeyID     string
	SESSecretAccessKey string
	SESSessionToken    string
}

func loadMailConfig() MailConfig {
	return MailConfig{
		Mailer:               strings.ToLower(getEnv("MAIL_MAILER", "")),
		FromName:             getEnv("MAIL_FROM_NAME", ""),
		Failover:             splitCSV(getEnv("MAIL_FAILOVER", "")),
		AllowLogInProduction: getEnv("MAIL_ALLOW_LOG_IN_PRODUCTION", "false") == "true",
		LogPath:              getEnv("MAIL_LOG_PATH", "storage/mail"),

		SMTPHost: getEnv("SMTP_HOST", "localhost"),
		// Unset, the port follows MAILHOG_SMTP_PORT, so moving Mailhog's port
		// in .env moves the connection with it.
		SMTPPort:       getEnv("SMTP_PORT", getEnv("MAILHOG_SMTP_PORT", "1025")),
		SMTPUsername:   getEnv("SMTP_USERNAME", ""),
		SMTPPassword:   getEnv("SMTP_PASSWORD", ""),
		SMTPEncryption: getEnv("SMTP_ENCRYPTION", ""),

		MailgunDomain:   getEnv("MAILGUN_DOMAIN", ""),
		MailgunSecret:   getEnv("MAILGUN_SECRET", ""),
		MailgunEndpoint: getEnv("MAILGUN_ENDPOINT", "api.mailgun.net"),

		PostmarkToken:         getEnv("POSTMARK_TOKEN", ""),
		PostmarkMessageStream: getEnv("POSTMARK_MESSAGE_STREAM", "outbound"),

		SendGridAPIKey: getEnv("SENDGRID_API_KEY", ""),

		SESRegion:          getEnv("AWS_SES_REGION", getEnv("AWS_REGION", "us-east-1")),
		SESAccessKeyID:     getEnv("AWS_ACCESS_KEY_ID", ""),
		SESSecretAccessKey: getEnv("AWS_SECRET_ACCESS_KEY", ""),
		SESSessionToken:    getEnv("AWS_SESSION_TOKEN", ""),
	}
}
