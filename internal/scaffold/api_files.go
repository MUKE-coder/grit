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
		filepath.Join(apiRoot, "go.mod"):                                       apiGoMod(opts),
		filepath.Join(apiRoot, ".gitignore"):                                   apiGitignore(),
		filepath.Join(apiRoot, "cmd", "server", "main.go"):                     apiMainGo(opts),
		filepath.Join(apiRoot, "internal", "config", "config.go"):              apiConfigGo(),
		filepath.Join(apiRoot, "internal", "database", "database.go"):          apiDatabaseGo(),
		filepath.Join(apiRoot, "internal", "database", "dialect.go"):           apiDialectGo(),
		filepath.Join(apiRoot, "internal", "models", "user.go"):                apiUserModelGo(),
		filepath.Join(apiRoot, "internal", "models", "upload.go"):              apiUploadModelGo(),
		filepath.Join(apiRoot, "internal", "services", "auth.go"):              apiAuthServiceGo(),
		filepath.Join(apiRoot, "internal", "models", "session.go"):             apiSessionModelGo(),
		filepath.Join(apiRoot, "internal", "services", "session.go"):           apiSessionServiceGo(),
		filepath.Join(apiRoot, "internal", "handlers", "session.go"):           apiSessionHandlerGo(),
		filepath.Join(apiRoot, "internal", "ids", "ids.go"):                    apiIDsGo(),
		filepath.Join(apiRoot, "internal", "ids", "ids_test.go"):               apiIDsTestGo(),
		filepath.Join(apiRoot, "internal", "models", "sso.go"):                 apiSSOModelGo(),
		filepath.Join(apiRoot, "internal", "services", "sso.go"):               apiSSOServiceGo(),
		filepath.Join(apiRoot, "internal", "handlers", "sso.go"):               apiSSOHandlerGo(),
		filepath.Join(apiRoot, "internal", "models", "saml.go"):                apiSAMLModelGo(),
		filepath.Join(apiRoot, "internal", "services", "saml.go"):              apiSAMLServiceGo(),
		filepath.Join(apiRoot, "internal", "handlers", "saml.go"):              apiSAMLHandlerGo(),
		filepath.Join(apiRoot, "internal", "handlers", "auth.go"):              apiAuthHandlerGo(),
		filepath.Join(apiRoot, "internal", "handlers", "user.go"):              apiUserHandlerGo(),
		filepath.Join(apiRoot, "internal", "middleware", "auth.go"):            apiAuthMiddlewareGo(),
		filepath.Join(apiRoot, "internal", "middleware", "cors.go"):            apiCorsMiddlewareGo(),
		filepath.Join(apiRoot, "internal", "middleware", "logger.go"):          apiLoggerMiddlewareGo(),
		filepath.Join(apiRoot, "internal", "middleware", "maintenance.go"):     apiMaintenanceMiddlewareGo(),
		filepath.Join(apiRoot, "internal", "middleware", "idempotency.go"):     apiIdempotencyMiddlewareGo(),
		filepath.Join(apiRoot, "internal", "paginate", "paginate.go"):          apiPaginateGo(),
		filepath.Join(apiRoot, "internal", "realtime", "hub.go"):               apiRealtimeHubGo(),
		filepath.Join(apiRoot, "internal", "handlers", "realtime.go"):          apiRealtimeHandlerGo(),
		filepath.Join(apiRoot, "internal", "sync", "registry.go"):              apiSyncRegistryGo(),
		filepath.Join(apiRoot, "internal", "sync", "policy.go"):                apiSyncPolicyGo(),
		filepath.Join(apiRoot, "internal", "events", "events.go"):              apiEventsGo(),
		filepath.Join(apiRoot, "internal", "workflow", "workflow.go"):          apiWorkflowGo(),
		filepath.Join(apiRoot, "internal", "settings", "settings.go"):          apiSettingsRegistryGo(),
		filepath.Join(apiRoot, "internal", "settings", "store.go"):             apiSettingsStoreGo(),
		filepath.Join(apiRoot, "internal", "settings", "defaults.go"):          apiSettingsDefaultsGo(),
		filepath.Join(apiRoot, "internal", "models", "setting.go"):             apiSettingsModelGo(),
		filepath.Join(apiRoot, "internal", "handlers", "settings.go"):          apiSettingsHandlerGo(),
		filepath.Join(apiRoot, "internal", "services", "event_subscribers.go"): apiEventsSubscribersGo(),
		filepath.Join(apiRoot, "internal", "handlers", "sync.go"):              apiSyncHandlerGo(),
		filepath.Join(apiRoot, "internal", "models", "activity_log.go"):        apiActivityLogModelGo(),
		filepath.Join(apiRoot, "internal", "middleware", "activity.go"):        apiActivityMiddlewareGo(),
		filepath.Join(apiRoot, "internal", "handlers", "activity.go"):          apiActivityHandlerGo(),
		filepath.Join(apiRoot, "internal", "export", "export.go"):              apiExportGo(),
		filepath.Join(apiRoot, "internal", "respond", "respond.go"):            apiRespondGo(),
		filepath.Join(apiRoot, "internal", "pdf", "pdf.go"):                    apiPDFGo(),
		filepath.Join(apiRoot, "internal", "pdf", "invoice.go"):                apiPDFInvoiceGo(),
		filepath.Join(apiRoot, "internal", "pdf", "record.go"):                 apiPDFRecordGo(),
		filepath.Join(apiRoot, "internal", "audit", "audit.go"):                apiAuditGo(),
		filepath.Join(apiRoot, "internal", "models", "webhook_event.go"):       apiWebhookEventModelGo(),
		filepath.Join(apiRoot, "internal", "webhooks", "webhooks.go"):          apiWebhooksGo(),
		filepath.Join(apiRoot, "internal", "webhooks", "verifiers.go"):         apiWebhooksVerifiersGo(),
		filepath.Join(apiRoot, "internal", "handlers", "webhooks.go"):          apiWebhooksHandlerGo(),
		filepath.Join(apiRoot, "internal", "models", "feature_flag.go"):        apiFeatureFlagModelGo(),
		filepath.Join(apiRoot, "internal", "flags", "flags.go"):                apiFlagsGo(),
		filepath.Join(apiRoot, "internal", "handlers", "flags.go"):             apiFlagsHandlerGo(),

		// v3.30 — semantic UserActivity log + ticket system
		filepath.Join(apiRoot, "internal", "models", "user_activity.go"): userActivityModelGo(),
		filepath.Join(apiRoot, "internal", "services", "activity.go"):    userActivityServiceGo(),
		// v3.31.49 — ResolveClientIP honours the X-Public-IP-Hint
		// header sent by the admin/web clients when the TCP peer is
		// loopback, so dev activity logs show the operator's actual
		// public IP instead of "::1".
		filepath.Join(apiRoot, "internal", "services", "clientip.go"):           clientIPHelperGo(),
		filepath.Join(apiRoot, "internal", "handlers", "user_activity.go"):      userActivityHandlerGo(),
		filepath.Join(apiRoot, "internal", "services", "ocsf.go"):               apiOCSFServiceGo(),
		filepath.Join(apiRoot, "internal", "handlers", "ocsf.go"):               apiOCSFHandlerGo(),
		filepath.Join(apiRoot, "internal", "services", "ocsf_test.go"):          apiOCSFTestGo(),
		filepath.Join(apiRoot, "internal", "models", "access_review.go"):        apiAccessReviewModelGo(),
		filepath.Join(apiRoot, "internal", "services", "access_review.go"):      apiAccessReviewServiceGo(),
		filepath.Join(apiRoot, "internal", "handlers", "access_review.go"):      apiAccessReviewHandlerGo(),
		filepath.Join(apiRoot, "internal", "services", "access_review_test.go"): apiAccessReviewTestGo(),
		filepath.Join(apiRoot, "internal", "models", "deletion_journal.go"):     apiGDPRModelGo(),
		filepath.Join(apiRoot, "internal", "services", "gdpr.go"):               apiGDPRServiceGo(),
		filepath.Join(apiRoot, "internal", "handlers", "gdpr.go"):               apiGDPRHandlerGo(),
		filepath.Join(apiRoot, "internal", "services", "gdpr_test.go"):          apiGDPRTestGo(),
		filepath.Join(apiRoot, "internal", "crypto", "field.go"):                apiCryptoFieldGo(),
		filepath.Join(apiRoot, "internal", "crypto", "field_test.go"):           apiCryptoFieldTestGo(),
		// v3.31.40 — per-user dashboard customisation
		filepath.Join(apiRoot, "internal", "models", "dashboard_layout.go"):   dashboardLayoutModelGo(),
		filepath.Join(apiRoot, "internal", "handlers", "dashboard_layout.go"): strings.ReplaceAll(dashboardLayoutHandlerGo(), "{{MODULE}}", opts.Module()),
		filepath.Join(apiRoot, "internal", "models", "ticket.go"):             ticketModelGo(),
		filepath.Join(apiRoot, "internal", "handlers", "ticket.go"):           ticketHandlerGo(),
		// v3.31.68 — background CSV import job tracking (shared across resources)
		filepath.Join(apiRoot, "internal", "models", "import_job.go"):    importJobModelGo(),
		filepath.Join(apiRoot, "internal", "handlers", "import_job.go"):  importJobHandlerGo(),
		filepath.Join(apiRoot, "internal", "services", "ticket_mail.go"): ticketMailGo(),
		filepath.Join(apiRoot, "internal", "routes", "routes.go"):        apiRoutesGo(),
		filepath.Join(apiRoot, "internal", "routes", "apidocs.go"):       apiDocsRoutesGo(),
		filepath.Join(apiRoot, ".air.toml"):                              airConfig(),
		// Test files — give the generated API a working test suite out of the box
		filepath.Join(apiRoot, "internal", "handlers", "auth_test.go"):               apiAuthTestGo(),
		filepath.Join(apiRoot, "internal", "handlers", "sso_test.go"):                apiSSOTestGo(),
		filepath.Join(apiRoot, "internal", "models", "bool_flags_test.go"):           apiBoolFlagTestGo(),
		filepath.Join(apiRoot, "internal", "handlers", "saml_test.go"):               apiSAMLTestGo(),
		filepath.Join(apiRoot, "internal", "services", "session_test.go"):            apiSessionTestGo(),
		filepath.Join(apiRoot, "internal", "models", "password_reset.go"):            apiPasswordResetModelGo(),
		filepath.Join(apiRoot, "internal", "services", "password_reset.go"):          apiPasswordResetServiceGo(),
		filepath.Join(apiRoot, "internal", "services", "password_reset_test.go"):     apiPasswordResetTestGo(),
		filepath.Join(apiRoot, "internal", "models", "email_verification.go"):        apiEmailVerifyModelGo(),
		filepath.Join(apiRoot, "internal", "services", "email_verification.go"):      apiEmailVerifyServiceGo(),
		filepath.Join(apiRoot, "internal", "services", "email_verification_test.go"): apiEmailVerifyTestGo(),
		filepath.Join(apiRoot, "internal", "models", "api_key.go"):                   apiAPIKeyModelGo(),
		filepath.Join(apiRoot, "internal", "services", "api_key.go"):                 apiAPIKeyServiceGo(),
		filepath.Join(apiRoot, "internal", "services", "api_key_test.go"):            apiAPIKeyTestGo(),
		filepath.Join(apiRoot, "internal", "middleware", "api_key.go"):               apiAPIKeyMiddlewareGo(),
		filepath.Join(apiRoot, "internal", "handlers", "api_key.go"):                 apiAPIKeyHandlerGo(),
		filepath.Join(apiRoot, "internal", "database", "api_keys_seeder.go"):         apiAPIKeySeederGo(),
		filepath.Join(apiRoot, "internal", "handlers", "user_test.go"):               apiUserTestGo(),
		filepath.Join(apiRoot, "internal", "handlers", "bench_test.go"):              apiBenchTestGo(),
	}

	for path, content := range files {
		// Replace module placeholder
		content = strings.ReplaceAll(content, "{{MODULE}}", module)
		// And the project name, which the settings defaults use so app.name
		// starts as something recognisable rather than a placeholder.
		content = strings.ReplaceAll(content, "{{PROJECT}}", opts.ProjectName)
		if err := writeFile(path, content); err != nil {
			return fmt.Errorf("writing %s: %w", path, err)
		}
	}

	return nil
}

// importJobModelGo is the shared ImportJob model. A single table tracks every
// resource's background CSV imports; the per-resource Import handler creates a
// row, processes the file in a goroutine, and updates the counters as it goes.
func importJobModelGo() string {
	return `package models

import (
	"time"

	"{{MODULE}}/internal/ids"
	"gorm.io/gorm"
)

// ImportJob records the progress of a background CSV import started by
// POST /<resource>/import. The handler inserts one row, launches a goroutine
// that streams the file, and updates Processed/Created/Skipped/Failed as it
// runs. Clients poll GET /imports/:id to drive a live progress bar, then read
// the final counts and per-row errors when Status is "completed". Errors holds
// up to the first 50 row failures as a JSON array string.
type ImportJob struct {
	ID        string    ` + "`gorm:\"primarykey;size:36\" json:\"id\"`" + `
	Resource  string    ` + "`gorm:\"size:255;index\" json:\"resource\"`" + `
	Status    string    ` + "`gorm:\"size:20;index\" json:\"status\"`" + ` // processing | completed | failed
	Total     int       ` + "`json:\"total\"`" + `
	Processed int       ` + "`json:\"processed\"`" + `
	Created   int       ` + "`json:\"created\"`" + `
	Skipped   int       ` + "`json:\"skipped\"`" + `
	Failed    int       ` + "`json:\"failed\"`" + `
	Errors    string    ` + "`gorm:\"type:text\" json:\"-\"`" + `
	Message   string    ` + "`gorm:\"size:500\" json:\"message\"`" + `
	CreatedAt time.Time ` + "`json:\"created_at\"`" + `
	UpdatedAt time.Time ` + "`json:\"updated_at\"`" + `
}

// BeforeCreate assigns a UUID so the job id is opaque in poll URLs.
func (m *ImportJob) BeforeCreate(tx *gorm.DB) error {
	if m.ID == "" {
		m.ID = ids.New()
	}
	if m.Status == "" {
		m.Status = "processing"
	}
	return nil
}
`
}

// importJobHandlerGo is the shared status endpoint clients poll while a
// background import runs. It is resource-agnostic — the resource field on the
// job says which table it targeted.
func importJobHandlerGo() string {
	return `package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"{{MODULE}}/internal/models"
)

// ImportJobHandler serves the progress/result of background CSV imports.
type ImportJobHandler struct {
	DB *gorm.DB
}

// GetByID returns a single import job. Poll this while Status is "processing"
// to drive a progress bar (processed/total), then read created/skipped/failed
// and the per-row errors once Status is "completed".
func (h *ImportJobHandler) GetByID(c *gin.Context) {
	var job models.ImportJob
	if err := h.DB.First(&job, "id = ?", c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": gin.H{"code": "NOT_FOUND", "message": "Import job not found"},
		})
		return
	}

	rowErrors := []map[string]interface{}{}
	if job.Errors != "" {
		_ = json.Unmarshal([]byte(job.Errors), &rowErrors)
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"id":        job.ID,
			"resource":  job.Resource,
			"status":    job.Status,
			"total":     job.Total,
			"processed": job.Processed,
			"created":   job.Created,
			"skipped":   job.Skipped,
			"failed":    job.Failed,
			"errors":    rowErrors,
			"message":   job.Message,
		},
	})
}
`
}

func apiGoMod(opts Options) string {
	return fmt.Sprintf(`module %s

go 1.21

require (
	github.com/MUKE-coder/gin-docs v0.0.0-20260222113017-4d647cb4e7aa
	github.com/MUKE-coder/gorm-studio v1.0.1
	// Pinned to the v1.0.0 commit on main.
	github.com/MUKE-coder/pulse v0.0.0-20260529025319-478cdfa8ce5f
	github.com/aws/aws-sdk-go-v2 v1.43.0
	github.com/aws/aws-sdk-go-v2/config v1.32.31
	github.com/aws/aws-sdk-go-v2/credentials v1.19.30
	github.com/aws/aws-sdk-go-v2/service/s3 v1.106.0
	github.com/brianvoe/gofakeit/v7 v7.15.0
	// SAML 2.0 service provider for enterprise SSO. OIDC covers every modern
	// IdP and needs no library, but SAML is still what a lot of enterprise
	// procurement asks for, and it cannot be hand-rolled safely — assertion
	// signature verification, audience restriction and clock-skew handling are
	// exactly the places a DIY implementation becomes an auth bypass.
	github.com/crewjam/saml v0.5.1
	// Pure-Go WebP, so image optimisation does not cost the static
	// cross-compiled binary. Lossless (VP8L) only: there is no pure-Go
	// lossy WebP encoder, which is why lossy compression targets JPEG.
	github.com/HugoSmits86/nativewebp v1.3.0
	github.com/disintegration/imaging v1.6.2
	github.com/gin-gonic/gin v1.11.0
	github.com/go-pdf/fpdf v1.4.3
	github.com/golang-jwt/jwt/v5 v5.3.1
	github.com/google/uuid v1.6.0
	github.com/gorilla/sessions v1.4.0
	github.com/gorilla/websocket v1.5.3
	github.com/hibiken/asynq v0.24.1
	github.com/markbates/goth v1.80.0
	github.com/joho/godotenv v1.5.1
	github.com/redis/go-redis/v9 v9.6.3
	github.com/skip2/go-qrcode v0.0.0-20200617195104-da1b6568686e
	github.com/xuri/excelize/v2 v2.10.0
	golang.org/x/crypto v0.53.0
	// Sentinel now ships a proper /v2 module path, so we track real tags.
	// v2.1.1 is the minimum safe release for WAF.Mode = ModeBlock: v2.1.0
	// fixed the SSRF rule matching "0.0.0.0" inside a Chrome User-Agent
	// (403'ing every Chrome 140/130/120/110 user), and v2.1.1 fixed
	// SQLi_Basic matching a bare "--" inside JWT cookies (roughly one
	// session in ten 403'd at random). v2.2.0 adds ValidateConfig, which
	// Mount runs at startup so dead config shows up in the boot log rather
	// than as a 403 weeks later. Do not downgrade below v2.1.1.
	github.com/MUKE-coder/sentinel/v2 v2.2.1
	gorm.io/datatypes v1.2.7
	gorm.io/driver/mysql v1.6.0
	gorm.io/driver/postgres v1.6.0
	gorm.io/gorm v1.31.1
)

require (
	github.com/stretchr/testify v1.11.1
	github.com/glebarez/sqlite v1.11.0
)

// Security floors for transitive dependencies. These are not imported directly;
// they are pinned because a dependency pulls in a version with a known,
// reachable vulnerability and MVS would otherwise settle on it. Each one was
// confirmed with govulncheck against a freshly scaffolded project.
//
// Raise these, never lower them. Re-check with:
//   go install golang.org/x/vuln/cmd/govulncheck@latest && govulncheck ./...
require (
	github.com/jackc/pgx/v5 v5.9.2 // GO-2026-5004
	github.com/quic-go/quic-go v0.59.1 // GO-2026-5676, GO-2025-4233
	golang.org/x/image v0.43.0 // GO-2026-5066, -5062, -5032, -5031, -4815
	golang.org/x/text v0.39.0 // GO-2026-5970
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
	// air v1.64+ deprecated `build.bin` in favour of `build.entrypoint`.
	// Both name the BUILT binary that air execs after each rebuild —
	// not the Go source directory. Always use a .exe suffix so Windows
	// CreateProcess accepts the file (no "open with" dialog); Linux
	// + macOS treat .exe as just part of the name. One config, every
	// platform.
	return `root = "."
tmp_dir = "tmp"

[build]
  cmd = "go build -o ./tmp/server.exe ./cmd/server"
  entrypoint = "./tmp/server.exe"
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
	"` + "{{MODULE}}" + `/internal/models"
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
	//
	// NOTE: these callback URLs are deliberately NOT versioned, even though the
	// routes now live under /api/v1. The same string is registered in the
	// Google / GitHub console as an authorized redirect URI — a value you
	// control there, not here. Adding "/v1" would stop matching what every
	// existing deployment has registered and break social login on upgrade,
	// which is the exact class of breakage the version prefix exists to avoid.
	// The unversioned path is re-dispatched to the current version by
	// mountLegacyAPIAlias (query string preserved), so these keep working.
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
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"

	"{{MODULE}}/internal/crypto"
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
	StorageDriver string        // "minio", "s3", "r2", or "b2"
	Storage       StorageConfig // Resolved config for the active driver

	ResendAPIKey string
	MailFrom     string

	CORSOrigins []string

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

	// AI (Vercel AI Gateway)
	AIGatewayAPIKey string
	AIGatewayModel  string
	AIGatewayURL    string

	// TOTP (Two-Factor Authentication)
	TOTPIssuer string
	RequireEmailVerification bool
	LoginMaxAttempts         int
	LoginLockoutWindow       time.Duration

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
		FieldEncryptionKey: getEnv("FIELD_ENCRYPTION_KEY", ""),
		RedisURL:    resolveRedisURL(),

		StorageDriver: storageDriver,
		Storage:       resolveStorage(storageDriver),

		ResendAPIKey: getEnv("RESEND_API_KEY", ""),
		MailFrom:     getEnv("MAIL_FROM", "noreply@localhost"),

		// The Wails desktop webview is allowed by middleware.isWailsOrigin (it
		// matches the wails.localhost host on any port), so it needs no entry
		// here — its dev origin includes a configurable port.
		CORSOrigins: strings.Split(getEnv("CORS_ORIGINS", "http://localhost:3000,http://localhost:3001"), ","),

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

		AIGatewayAPIKey: getEnv("AI_GATEWAY_API_KEY", ""),
		AIGatewayModel:  getEnv("AI_GATEWAY_MODEL", "anthropic/claude-sonnet-4-6"),
		AIGatewayURL:    getEnv("AI_GATEWAY_URL", "https://ai-gateway.vercel.sh/v1"),

		TOTPIssuer: getEnv("TOTP_ISSUER", getEnv("APP_NAME", "grit-app")),
		// Off by default and deliberately so: switching it on for an existing
		// project would lock out every user at once, because they all have a
		// NULL email_verified_at.
		RequireEmailVerification: getEnv("REQUIRE_EMAIL_VERIFICATION", "false") == "true",
		LoginMaxAttempts:   getEnvInt("LOGIN_MAX_ATTEMPTS", 10),
		LoginLockoutWindow: getEnvDuration("LOGIN_LOCKOUT_MINUTES", 15) * time.Minute,

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
		return "redis://localhost:6380"
	}
	if strings.TrimSpace(v) == "" {
		log.Println("REDIS_URL is empty: cache, background jobs and cron are disabled")
		return ""
	}
	return v
}

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
	default: // minio
		return StorageConfig{
			Endpoint:  getEnv("MINIO_ENDPOINT", "http://localhost:9002"),
			AccessKey: getEnv("MINIO_ACCESS_KEY", "minioadmin"),
			SecretKey: getEnv("MINIO_SECRET_KEY", "minioadmin"),
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

// getEnvDuration reads a whole number of units; the caller multiplies by the
// unit it means, which keeps the env var name self-describing
// (LOGIN_LOCKOUT_MINUTES=15 rather than a duration string nobody formats
// consistently).
func getEnvDuration(key string, fallback int) time.Duration {
	return time.Duration(getEnvInt(key, fallback))
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

// apiDialectGo emits internal/database/dialect.go.
// APIDialectGo is exported so grit generate resource can add this file to a
// project that predates it. The generated handlers depend on it.
func APIDialectGo() string { return apiDialectGo() }

func apiDialectGo() string {
	return `package database

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
)

// SupportsReturning reports whether the connected dialect can hand back the
// written row from an INSERT or UPDATE.
//
// Postgres and SQLite can. MySQL cannot, and this is the important part: it
// does not error when asked. The write succeeds, the RETURNING clause is
// dropped, and the struct comes back with every database-assigned default
// still at its zero value. A handler that skipped its reload on the strength
// of RETURNING would then answer 201 with a half-empty record.
func SupportsReturning(db *gorm.DB) bool {
	switch db.Dialector.Name() {
	case "postgres", "sqlite":
		return true
	default:
		return false
	}
}

// Write returns a session for a single-statement write: no wrapping
// transaction, and RETURNING where the dialect has it.
//
// Skipping the transaction is safe only because the caller has already
// established there is exactly one statement. The generator decides that from
// the resource definition, where it can be known rather than assumed.
func Write(db *gorm.DB) *gorm.DB {
	tx := db.Session(&gorm.Session{SkipDefaultTransaction: true})
	if SupportsReturning(db) {
		tx = tx.Clauses(clause.Returning{})
	}
	return tx
}

// TableCount returns how many tables the connected database holds, or 0 when
// the dialect cannot be asked. Purely informational: it feeds the "tables: N"
// figure on the health card.
//
// Three dialects need three different questions. information_schema exists on
// Postgres and MySQL but not on SQLite, and the two that have it spell "this
// database" differently. Asking the Postgres question everywhere logged a red
// SQL error on every health poll of a SQLite project, and quietly returned 0
// on MySQL, where current_schema() does not exist either.
//
// Errors are swallowed with the logger silenced, because a missing tooltip
// figure is not worth a stack of scary log lines on a healthy server.
func TableCount(db *gorm.DB) int {
	var query string
	switch db.Dialector.Name() {
	case "postgres":
		query = "SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = current_schema()"
	case "mysql":
		query = "SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = DATABASE()"
	case "sqlite":
		query = "SELECT COUNT(*) FROM sqlite_master WHERE type = 'table' AND name NOT LIKE 'sqlite_%'"
	default:
		return 0
	}

	var count int
	quiet := db.Session(&gorm.Session{Logger: db.Logger.LogMode(logger.Silent)})
	if err := quiet.Raw(query).Scan(&count).Error; err != nil {
		return 0
	}
	return count
}
`
}

func apiDatabaseGo() string {
	return `package database

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/glebarez/sqlite"
	"gorm.io/driver/mysql"
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
	// Two GORM knobs that get recommended a lot. Both are off by default here,
	// and both defaults were measured rather than assumed — k6, 50 concurrent
	// writers, 4 CPUs, three 20-second runs each, median req/s on inserts:
	//
	//	off / off        740     what ships
	//	PrepareStmt      738     no difference
	//	+ skip the tx  1,294     +75%
	//
	// PrepareStmt caches a prepared statement per connection so a query is
	// planned once instead of per request. On this workload it measured as
	// nothing — the cache is mutex-guarded and under concurrency the contention
	// cancels out the saved planning. It also breaks against a connection pooler
	// in transaction mode (pgbouncer, RDS Proxy), because server-side prepared
	// statements do not survive a pooler that hands each transaction a different
	// backend. No measured gain, real downsides, so: opt in with
	// DB_PREPARED_STATEMENTS=true if your own numbers disagree.
	gormCfg := &gorm.Config{
		Logger:      logger.Default.LogMode(logLevel),
		PrepareStmt: os.Getenv("DB_PREPARED_STATEMENTS") == "true",
	}

	// Skipping the default transaction is the one that actually pays — GORM
	// wraps every Create, Update and Delete in an implicit transaction, so a
	// single-row insert costs BEGIN + INSERT + COMMIT where one round trip would
	// do. Turning it off was worth 75% here.
	//
	// It is still off by default, and that is a correctness decision rather than
	// a cautious one. The resource generator emits models with relations, and
	// saving an invoice with its line items is several INSERTs that GORM's
	// implicit transaction is currently what makes atomic — the generated
	// handler does not open its own. Without it, a failure halfway through
	// leaves an invoice holding some of its lines, with nothing logged and
	// nobody the wiser until the totals stop adding up.
	//
	// So: if your resources are flat, DB_SKIP_DEFAULT_TRANSACTION=true is close
	// to free throughput. If you generate anything with line items, leave it
	// alone until the generated handlers wrap their own writes.
	if os.Getenv("DB_SKIP_DEFAULT_TRANSACTION") == "true" {
		gormCfg.SkipDefaultTransaction = true
	}

	var (
		db  *gorm.DB
		err error
	)

	switch {
	case strings.HasPrefix(dsn, "sqlite://"):
		db, err = gorm.Open(sqlite.Open(strings.TrimPrefix(dsn, "sqlite://")), gormCfg)
	case strings.HasPrefix(dsn, "sqlite:"):
		db, err = gorm.Open(sqlite.Open(strings.TrimPrefix(dsn, "sqlite:")), gormCfg)
	case strings.HasPrefix(dsn, "mysql://"), strings.HasPrefix(dsn, "mysql:"):
		// go-sql-driver wants "user:pass@tcp(host:port)/db", not a URL, so the
		// scheme is stripped rather than parsed. parseTime is not optional:
		// without it DATETIME columns arrive as []byte and every time.Time
		// field on every model fails to scan.
		my := strings.TrimPrefix(strings.TrimPrefix(dsn, "mysql://"), "mysql:")
		if !strings.Contains(my, "parseTime=") {
			sep := "?"
			if strings.Contains(my, "?") {
				sep = "&"
			}
			my += sep + "parseTime=true&loc=UTC"
		}
		db, err = gorm.Open(mysql.Open(my), gormCfg)
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
	//
	// Idle defaults to Open, and that default matters more than it looks. When
	// idle is lower, a request past the idle limit returns its connection to a
	// full pool, so the connection is CLOSED — and the next request opens a new
	// one, which makes Postgres fork a backend process. Under concurrency that
	// is a connection storm, and it surfaces as database CPU rather than as
	// anything you would think to look for in the application.
	//
	// Measured with k6 at 50 VUs, 4 CPUs per container, single-row reads:
	// idle=10 gave ~810 req/s with Postgres pinned near 840% while the API used
	// 196%; idle=100 gave ~2,720 req/s with both around 300%. Same binary, same
	// query.
	//
	// Both are tunable because the right answer depends on the workload. If
	// your queries are heavy enough to saturate the database — an unindexed
	// COUNT over a large table on every request, say — a smaller pool acts as
	// admission control and can measure faster, because queueing in the app is
	// cheaper than thrashing in Postgres. Start here, then measure.
	maxOpen := getEnvInt("DB_MAX_OPEN_CONNS", 100)
	if maxOpen < 1 {
		maxOpen = 1
	}
	maxIdle := getEnvInt("DB_MAX_IDLE_CONNS", maxOpen)
	if maxIdle < 1 || maxIdle > maxOpen {
		maxIdle = maxOpen
	}
	sqlDB.SetMaxIdleConns(maxIdle)
	sqlDB.SetMaxOpenConns(maxOpen)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)
	sqlDB.SetConnMaxIdleTime(10 * time.Minute)

	log.Println("Database connected successfully")
	return db, nil
}

// getEnvInt reads a whole-number env var. A malformed value falls back rather
// than failing the boot — a typo in DB_MAX_OPEN_CONNS should not stop the app
// from starting.
func getEnvInt(key string, fallback int) int {
	if raw := os.Getenv(key); raw != "" {
		if n, err := strconv.Atoi(raw); err == nil {
			return n
		}
		log.Printf("warning: %s=%q is not a number, using %d", key, raw, fallback)
	}
	return fallback
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

	"{{MODULE}}/internal/ids"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"{{MODULE}}/internal/crypto"
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
	Bio             crypto.EncryptedString ` + "`" + `gorm:"type:text" json:"bio"` + "`" + `
	// No gorm default on this bool. GORM omits zero-valued fields from an
	// INSERT when the column carries a default, so default:true made
	// Active:false unstorable on create — an admin creating a deactivated user
	// silently got an active one. Every create path sets this explicitly.
	Active          bool           ` + "`" + `gorm:"" json:"active"` + "`" + `
	Provider        string         ` + "`" + `gorm:"size:50;default:'local'" json:"provider"` + "`" + `
	GoogleID        string         ` + "`" + `gorm:"size:255" json:"-"` + "`" + `
	GithubID        string         ` + "`" + `gorm:"size:255" json:"-"` + "`" + `
	EmailVerifiedAt *time.Time     ` + "`" + `json:"email_verified_at"` + "`" + `

	// Per-account brute-force protection. Sentinel rate-limits by IP, which
	// does nothing against attempts spread across many addresses at one
	// account — the shape of every credential-stuffing run. FailedLoginCount
	// is reset on any successful sign-in.
	FailedLoginCount int        ` + "`" + `gorm:"default:0" json:"-"` + "`" + `
	LockedUntil      *time.Time ` + "`" + `json:"locked_until,omitempty"` + "`" + `
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
		u.ID = ids.New()
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
	tx.Statement.SetColumn("version", gorm.Expr("version + 1"))
	return nil
}

// BeforeCreate generates a UUID for uploads.
func (u *Upload) BeforeCreate(tx *gorm.DB) error {
	if u.ID == "" {
		u.ID = ids.New()
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
		// Server-side refresh sessions — must exist before anything logs in.
		&Session{},
		&PasswordResetToken{},
		&EmailVerificationToken{},
		&APIKey{},
		// Role/UserRole must migrate before anything authorises a request.
		&Role{},
		&UserRole{},
		&Upload{},
		&Blog{},
		&TwoFactorConfig{},
		&TrustedDevice{},
		&TOTPPendingToken{},
		&ActivityLog{},
		&WebhookEvent{},
		&FeatureFlag{},
		&FlagExposure{},
		&Notification{},
		// v3.30
		&UserActivity{},
		&AccessReview{},
		&AccessReviewItem{},
		&DeletionJournal{},
		&Ticket{},
		&TicketReply{},
		// v3.31.20 — public form sharing (Phase 2)
		&FormShare{},
		// v3.31.25 — audit log for public submissions
		&FormSubmission{},
		// v3.31.40 — per-user dashboard customisation
		&DashboardLayout{},
		// v3.31.68 — background CSV import job tracking
		&ImportJob{},
		// v3.31.77 — full-database backup index
		&Backup{},
		// backup schedule (period + time-of-day for automatic backups)
		&BackupSchedule{},
		// enterprise SSO: one OIDC connection per customer, plus the external
		// identities linking their users to local accounts
		&SSOConnection{},
		&UserIdentity{},
		// the service provider's own signing keypair, generated on first use
		&SAMLKeypair{},
		&Setting{},
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
	log.Printf("DATABASE MIGRATION: %d model(s) registered", len(models))
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
		log.Printf("  ~ %T: added %d column(s): %s", model, len(added), strings.Join(added, ", "))
		altered++
		columnsAdded += len(added)
	}

	log.Println(thinSep)
	log.Printf("Migration done: %d table(s) created, %d altered (+%d column(s)), %d unchanged.",
		created, altered, columnsAdded, unchanged)

	// Seed the default roles here rather than in database.Seed(): authorization
	// must work on a freshly migrated database, without anyone remembering to
	// run "grit seed". SeedRoles is idempotent and never overwrites an existing
	// role's grants.
	if err := SeedRoles(db); err != nil {
		return fmt.Errorf("seeding default roles: %w", err)
	}

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
	Path         string         ` + "`" + `gorm:"size:500;not null;index" json:"path"` + "`" + `
	URL          string         ` + "`" + `gorm:"size:500" json:"url"` + "`" + `
	ThumbnailURL string         ` + "`" + `gorm:"size:500" json:"thumbnail_url"` + "`" + `
	UserID       string         ` + "`" + `gorm:"size:36;index;not null" json:"user_id"` + "`" + `
	User         User           ` + "`" + `gorm:"foreignKey:UserID" json:"-"` + "`" + `
	Version      int            ` + "`" + `gorm:"not null;default:1" json:"version"` + "`" + `
	// v3.31.33 -- claimed_at is set when a parent record references this
	// upload's path/key via a FileRef column. NULL means abandoned, and
	// the daily orphan cleanup cron deletes the S3 object + DB row when
	// the upload is older than 24h and still unclaimed.
	ClaimedAt    *time.Time     ` + "`" + `gorm:"index" json:"claimed_at,omitempty"` + "`" + `
	CreatedAt    time.Time      ` + "`" + `json:"created_at"` + "`" + `
	UpdatedAt    time.Time      ` + "`" + `json:"updated_at"` + "`" + `
	DeletedAt    gorm.DeletedAt ` + "`" + `gorm:"index" json:"-"` + "`" + `
}

// BeforeUpdate increments Version on every server-side write so offline
// clients can detect that a record they edited has moved on.
func (u *Upload) BeforeUpdate(tx *gorm.DB) error {
	tx.Statement.SetColumn("version", gorm.Expr("version + 1"))
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

	// Every token gets a unique jti. Without it, two tokens minted for the same
	// user in the same second are byte-identical — same claims, same
	// second-resolution exp, same key — so two different devices would share one
	// refresh token and could not be told apart or revoked independently.
	jti, err := GenerateResetToken()
	if err != nil {
		return "", 0, fmt.Errorf("generating token id: %w", err)
	}

	claims := &Claims{
		UserID: userID,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        jti,
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
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/markbates/goth/gothic"
	"gorm.io/gorm"

	"golang.org/x/crypto/bcrypt"

	"` + "{{MODULE}}" + `/internal/config"
	"` + "{{MODULE}}" + `/internal/mail"
	"` + "{{MODULE}}" + `/internal/models"
	"` + "{{MODULE}}" + `/internal/services"
	"` + "{{MODULE}}" + `/internal/totp"
)

// AuthHandler handles authentication endpoints.
type AuthHandler struct {
	DB          *gorm.DB
	AuthService *services.AuthService
	Config      *config.Config
	// Mailer is optional. Without it the reset link is logged instead of sent,
	// which is what you want in dev and must never be what happens in prod —
	// ForgotPassword refuses to log the token when APP_ENV is production.
	Mailer *mail.Mailer
}

// AuthResponse documents the body returned by register, login and refresh.
//
// The handlers emit gin.H rather than this struct, so nothing enforces the two
// agree — if you change what an auth handler writes, change this with it. It
// exists because a reference that says "No Body" is worse than no reference.
type AuthResponse struct {
	Data struct {
		Tokens struct {
			AccessToken  string ` + "`" + `json:"access_token"` + "`" + `
			RefreshToken string ` + "`" + `json:"refresh_token"` + "`" + `
			ExpiresAt    int64  ` + "`" + `json:"expires_at"` + "`" + `
		} ` + "`" + `json:"tokens"` + "`" + `
		User models.User ` + "`" + `json:"user"` + "`" + `
	} ` + "`" + `json:"data"` + "`" + `
	Message string ` + "`" + `json:"message"` + "`" + `
}

// The types below exist so the API reference can show a body instead of
// "No Body". The handlers emit gin.H maps, so nothing enforces that these stay
// in step — if you change what a handler writes, change its type here too.
// They are documentation with a compiler attached, which is still better than
// prose nobody updates.

// MessageResponse is the plain acknowledgement shape.
type MessageResponse struct {
	Message string ` + "`" + `json:"message"` + "`" + `
}

// IssuedKeyResponse is returned once, when an API key is created.
type IssuedKeyResponse struct {
	Data struct {
		Key   models.APIKey ` + "`" + `json:"key"` + "`" + `
		Token string        ` + "`" + `json:"token"` + "`" + `
	} ` + "`" + `json:"data"` + "`" + `
	Message string ` + "`" + `json:"message"` + "`" + `
}

// TOTPStatusResponse describes the caller's two-factor state.
type TOTPStatusResponse struct {
	Data struct {
		Enabled              bool  ` + "`" + `json:"enabled"` + "`" + `
		BackupCodesRemaining int   ` + "`" + `json:"backup_codes_remaining"` + "`" + `
		TrustedDevices       int64 ` + "`" + `json:"trusted_devices"` + "`" + `
	} ` + "`" + `json:"data"` + "`" + `
}

// PresignResponse carries the URL a browser PUTs to, and the key to send back
// to /uploads/complete afterwards.
type PresignResponse struct {
	Data struct {
		PresignedURL string ` + "`" + `json:"presigned_url"` + "`" + `
		Key          string ` + "`" + `json:"key"` + "`" + `
	} ` + "`" + `json:"data"` + "`" + `
}

// ChainStatusResponse is the activity-log integrity verdict.
type ChainStatusResponse struct {
	Valid        bool   ` + "`" + `json:"valid"` + "`" + `
	TotalEntries int    ` + "`" + `json:"total_entries"` + "`" + `
	BrokenAt     int    ` + "`" + `json:"broken_at,omitempty"` + "`" + `
	BrokenAtID   string ` + "`" + `json:"broken_at_id,omitempty"` + "`" + `
	Expected     string ` + "`" + `json:"expected,omitempty"` + "`" + `
	Got          string ` + "`" + `json:"got,omitempty"` + "`" + `
	Message      string ` + "`" + `json:"message,omitempty"` + "`" + `
}

// ErrorResponse is the error envelope every endpoint uses.
type ErrorResponse struct {
	Error struct {
		Code    string            ` + "`" + `json:"code"` + "`" + `
		Message string            ` + "`" + `json:"message"` + "`" + `
		Details map[string]string ` + "`" + `json:"details,omitempty"` + "`" + `
	} ` + "`" + `json:"error"` + "`" + `
}

type RegisterRequest struct {
	FirstName  string ` + "`" + `json:"first_name" binding:"required,min=2"` + "`" + `
	LastName   string ` + "`" + `json:"last_name" binding:"required,min=2"` + "`" + `
	Email      string ` + "`" + `json:"email" binding:"required,email"` + "`" + `
	Password   string ` + "`" + `json:"password" binding:"required,min=8"` + "`" + `
	MACAddress string ` + "`" + `json:"mac_address"` + "`" + ` // optional — provided by client if available
}

type LoginRequest struct {
	Email    string ` + "`" + `json:"email" binding:"required,email"` + "`" + `
	Password string ` + "`" + `json:"password" binding:"required"` + "`" + `
}

type RefreshRequest struct {
	RefreshToken string ` + "`" + `json:"refresh_token" binding:"required"` + "`" + `
}

type ForgotPasswordRequest struct {
	Email string ` + "`" + `json:"email" binding:"required,email"` + "`" + `
}

type ResetPasswordRequest struct {
	Token    string ` + "`" + `json:"token" binding:"required"` + "`" + `
	Password string ` + "`" + `json:"password" binding:"required,min=8"` + "`" + `
}

// Register creates a new user account.
func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
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

	// Off the request path: signup should not wait on SMTP, and a mail failure
	// must not fail an account that was created successfully.
	go h.deliverVerificationEmail(user)

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

	// Record the refresh token as a server-side session so it can be revoked.
	if _, err := services.CreateSession(h.DB, c, user.ID, tokens.RefreshToken); err != nil {
		log.Printf("auth: failed to record session for %s: %v", user.ID, err)
	}

	// Set HttpOnly auth cookies for browser clients.
	h.AuthService.SetAuthCookies(c, tokens)

	// v3.30.1: emit a semantic activity row so /system/activity reflects
	// the signup. Non-blocking — a logging failure won't fail the
	// register request.
	services.LogRegister(h.DB, c, user.ID, user.Email)

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
	var req LoginRequest
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
		// v3.30.1: unknown email is the most common brute-force fingerprint;
		// surface it in /system/activity as "warn" severity so operators
		// can spot credential-stuffing spikes.
		services.LogLoginFailed(h.DB, c, req.Email)
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": gin.H{
				"code":    "INVALID_CREDENTIALS",
				"message": "Invalid email or password",
			},
		})
		return
	}

	if !user.Active {
		services.LogActivity(h.DB, c, services.ActivityArgs{
			Action:       "auth.login_blocked",
			Severity:     "warn",
			Summary:      "Sign-in blocked for disabled account " + user.Email,
			ResourceType: "user",
			ResourceID:   user.ID,
		})
		c.JSON(http.StatusForbidden, gin.H{
			"error": gin.H{
				"code":    "ACCOUNT_DISABLED",
				"message": "Your account has been disabled",
			},
		})
		return
	}

	// Locked accounts are refused before the password is even compared, so a
	// lockout cannot be probed by timing the comparison.
	if user.LockedUntil != nil && time.Now().Before(*user.LockedUntil) {
		remaining := time.Until(*user.LockedUntil).Round(time.Minute)
		if remaining < time.Minute {
			remaining = time.Minute
		}
		services.LogActivity(h.DB, c, services.ActivityArgs{
			Action:       "auth.login_locked",
			Severity:     "warn",
			Summary:      "Sign-in refused: account is temporarily locked",
			ResourceType: "user",
			ResourceID:   user.ID,
		})
		c.JSON(http.StatusTooManyRequests, gin.H{
			"error": gin.H{
				"code":    "ACCOUNT_LOCKED",
				"message": fmt.Sprintf("Too many failed attempts. Try again in about %d minute(s), or reset your password.", int(remaining.Minutes())),
			},
		})
		return
	}

	// Opt-in gate. Social and SSO sign-ins are unaffected — the IdP already
	// proved the address, and those paths set EmailVerifiedAt on first login.
	if h.Config.RequireEmailVerification && user.EmailVerifiedAt == nil && user.Password != "" {
		c.JSON(http.StatusForbidden, gin.H{
			"error": gin.H{
				"code":    "EMAIL_NOT_VERIFIED",
				"message": "Confirm your email address before signing in. Check your inbox for the link.",
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
		// Wrong password on a real account — distinct from "unknown email"
		// because Sentinel's brute-force heuristics weight these higher.
		services.LogLoginFailed(h.DB, c, req.Email)
		h.registerFailedLogin(&user)
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": gin.H{
				"code":    "INVALID_CREDENTIALS",
				"message": "Invalid email or password",
			},
		})
		return
	}

	// Any successful password check clears the counter, including one that
	// still has 2FA ahead of it — the password was correct, which is what this
	// counter measures.
	if user.FailedLoginCount > 0 || user.LockedUntil != nil {
		h.DB.Model(&models.User{}).Where("id = ?", user.ID).
			Updates(map[string]interface{}{"failed_login_count": 0, "locked_until": nil})
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
	//
	// Record the refresh token as a server-side session so this device can be
	// listed and revoked later.
	if _, err := services.CreateSession(h.DB, c, user.ID, tokens.RefreshToken); err != nil {
		log.Printf("auth: failed to record session for %s: %v", user.ID, err)
	}
	h.AuthService.SetAuthCookies(c, tokens)

	// v3.30.1: successful sign-in lands in /system/activity at info
	// severity. IP + user-agent come from the request context inside
	// LogLogin so brute-force investigation has the full pair.
	services.LogLogin(h.DB, c, user.ID, user.Email)

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
		var req RefreshRequest
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

	// Re-verify the account on every refresh. A stateless refresh token is
	// otherwise valid for its full lifetime even after the user is deleted or
	// deactivated; re-loading the user closes that window and lets a role
	// change take effect on the next refresh (partial revocation without a
	// server-side token store).
	var user models.User
	if err := h.DB.First(&user, "id = ?", claims.UserID).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": gin.H{
				"code":    "INVALID_TOKEN",
				"message": "Account no longer exists",
			},
		})
		return
	}
	if !user.Active {
		c.JSON(http.StatusForbidden, gin.H{
			"error": gin.H{
				"code":    "ACCOUNT_DISABLED",
				"message": "This account has been disabled",
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

	// Refresh the HttpOnly cookies so the new access token lands in the
	// browser without any JS handling. The bearer JSON path is unchanged
	// for native clients.
	//
	// Rotate the session. This is where revocation actually bites: a session
	// that was revoked, idled out, aged past its absolute limit, or whose token
	// was replayed after rotation has no live row, and the refresh is refused.
	if _, err := services.RotateSession(h.DB, c, refreshToken, tokens.RefreshToken); err != nil {
		h.AuthService.ClearAuthCookies(c)
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": gin.H{
				"code":    "SESSION_REVOKED",
				"message": "This session is no longer valid. Please sign in again.",
			},
		})
		return
	}
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
	// v3.30.1: read the user out of context BEFORE clearing cookies so
	// the activity row carries the right email. The auth middleware set
	// "user" on the gin context when the request came in.
	var actorID, actorEmail string
	if v, ok := c.Get("user"); ok {
		if u, ok := v.(models.User); ok {
			actorID = u.ID
			actorEmail = u.Email
		}
	}

	// Revoke the server-side session BEFORE clearing cookies — once they're
	// gone we can no longer identify which session to kill. This is what makes
	// logout real: the refresh token is dead immediately, not merely forgotten
	// by this browser.
	if rt, err := c.Cookie("grit_refresh"); err == nil && rt != "" {
		if err := services.RevokeSessionByToken(h.DB, rt); err != nil {
			log.Printf("auth: failed to revoke session on logout: %v", err)
		}
	}

	h.AuthService.ClearAuthCookies(c)

	if actorID != "" {
		services.LogLogout(h.DB, c, actorID, actorEmail)
	}
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
	var req ForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": gin.H{
				"code":    "VALIDATION_ERROR",
				"message": err.Error(),
			},
		})
		return
	}

	// One response for every outcome. Any variation — a different message, a
	// different status, a measurably different latency — turns this endpoint
	// into an oracle for which email addresses hold accounts.
	const genericResponse = "If an account with that email exists, a password reset link has been sent"

	var user models.User
	if err := h.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"message": genericResponse})
		return
	}

	// Everything past the lookup — minting the token, storing it, delivering the
	// link — runs off the request path. Both branches then do the same work
	// before answering (parse, one indexed SELECT), so a registered address does
	// not take measurably longer to respond than an unregistered one. Identical
	// wording with a distinguishable response time is still an oracle.
	//
	// c.ClientIP() is read here: the gin context must not be touched once the
	// handler has returned.
	go h.deliverPasswordReset(user, c.ClientIP())

	c.JSON(http.StatusOK, gin.H{"message": genericResponse})
}

// deliverPasswordReset issues a reset token and sends the link. It runs in its
// own goroutine, so it owns its context and reports failures only to the log —
// there is no caller left to tell, and telling the original one would have
// confirmed the address exists.
func (h *AuthHandler) deliverPasswordReset(user models.User, clientIP string) {
	token, err := services.GenerateResetToken()
	if err != nil {
		log.Printf("password reset: generating token for %s: %v", user.Email, err)
		return
	}

	if _, err := services.CreatePasswordResetToken(h.DB, user.ID, token, clientIP); err != nil {
		log.Printf("password reset: storing token for %s: %v", user.Email, err)
		return
	}

	resetURL := strings.TrimSuffix(h.Config.OAuthFrontendURL, "/") + "/reset-password?token=" + url.QueryEscape(token)

	if h.Mailer != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		if err := h.Mailer.Send(ctx, mail.SendOptions{
			To:       user.Email,
			Subject:  "Reset your password",
			Template: "password-reset",
			Data: map[string]interface{}{
				"AppName":  h.Config.AppName,
				"Title":    "Reset your password",
				"Message":  "We received a request to reset your password. This link expires in one hour and can only be used once. If you didn't ask for this, you can ignore this email.",
				"ResetURL": resetURL,
				"Year":     time.Now().Year(),
			},
		}); err != nil {
			log.Printf("password reset: sending email to %s: %v", user.Email, err)
		}
		return
	}

	if h.Config.AppEnv == "production" {
		// No mailer in production means nobody can complete a reset. Say so
		// loudly rather than printing a working token into the log — a live
		// reset link in a log file is a credential.
		log.Printf("password reset: NO MAILER CONFIGURED: %s cannot receive a reset link. Set RESEND_API_KEY.", user.Email)
		return
	}

	// Dev convenience only, and only outside production.
	log.Printf("password reset link for %s: %s", user.Email, resetURL)
}


// Unlock clears a lockout early. Waiting out the window is the normal path;
// this exists for the support call that follows a user locking themselves out
// five minutes before a demo.
func (h *UserHandler) Unlock(c *gin.Context) {
	id := c.Param("id")

	res := h.DB.Model(&models.User{}).Where("id = ?", id).
		Updates(map[string]interface{}{"locked_until": nil, "failed_login_count": 0})
	if res.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{"code": "INTERNAL_ERROR", "message": "Failed to unlock the account"},
		})
		return
	}
	if res.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": gin.H{"code": "NOT_FOUND", "message": "User not found"},
		})
		return
	}

	services.LogActivity(h.DB, c, services.ActivityArgs{
		Action:       "user.unlock",
		Severity:     "warn",
		Summary:      "Account lockout cleared by an administrator",
		ResourceType: "user",
		ResourceID:   id,
	})

	c.JSON(http.StatusOK, gin.H{"message": "Account unlocked"})
}

// registerFailedLogin counts a wrong password against the account and locks it
// once the threshold is reached.
//
// Only wrong-password-on-a-real-account is counted. Counting unknown emails
// would let anyone lock an address they can guess, which turns a defence into
// a denial-of-service tool.
//
// The increment is a single UPDATE rather than read-modify-write, so parallel
// attempts cannot each read the same count and overwrite one another.
func (h *AuthHandler) registerFailedLogin(user *models.User) {
	max := h.Config.LoginMaxAttempts
	if max <= 0 {
		return // lockout disabled
	}

	if err := h.DB.Model(&models.User{}).
		Where("id = ?", user.ID).
		UpdateColumn("failed_login_count", gorm.Expr("failed_login_count + 1")).Error; err != nil {
		log.Printf("lockout: incrementing failed_login_count for %s: %v", user.ID, err)
		return
	}

	var fresh models.User
	if err := h.DB.Select("id", "failed_login_count").First(&fresh, "id = ?", user.ID).Error; err != nil {
		return
	}
	if fresh.FailedLoginCount < max {
		return
	}

	until := time.Now().Add(h.Config.LoginLockoutWindow)
	if err := h.DB.Model(&models.User{}).
		Where("id = ?", user.ID).
		Updates(map[string]interface{}{"locked_until": until, "failed_login_count": 0}).Error; err != nil {
		log.Printf("lockout: locking %s: %v", user.ID, err)
		return
	}
	log.Printf("lockout: %s locked until %s after %d failed attempts", user.Email, until.Format(time.RFC3339), max)
}

// SendVerificationEmail issues a fresh verification link for the signed-in
// user. Authenticated on purpose: an unauthenticated "send a link to this
// address" endpoint is a spam cannon aimed at whoever you name.
func (h *AuthHandler) SendVerificationEmail(c *gin.Context) {
	userID := c.GetString("user_id")

	var user models.User
	if err := h.DB.First(&user, "id = ?", userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": gin.H{"code": "NOT_FOUND", "message": "User not found"},
		})
		return
	}

	if user.EmailVerifiedAt != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{"code": "ALREADY_VERIFIED", "message": "This email is already verified"},
		})
		return
	}

	go h.deliverVerificationEmail(user)

	c.JSON(http.StatusOK, gin.H{
		"message": "Verification email sent. The link is valid for 48 hours.",
	})
}

// The token from a verification link.
type VerifyEmailRequest struct {
	Token string ` + "`" + `json:"token" binding:"required"` + "`" + `
}

// VerifyEmail consumes a verification token. Public — the user clicks this
// from their mail client, where they are usually not signed in.

func (h *AuthHandler) VerifyEmail(c *gin.Context) {
	var req VerifyEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": gin.H{"code": "VALIDATION_ERROR", "message": err.Error()},
		})
		return
	}

	if _, err := services.ConsumeEmailVerificationToken(h.DB, req.Token); err != nil {
		// One message for expired, spent, unknown and address-changed. Telling
		// them apart tells an attacker which tokens once existed.
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_TOKEN",
				"message": "That verification link is invalid or has expired. Request a new one.",
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Email verified"})
}

// deliverVerificationEmail mints a token and sends the link, off the request
// path so a slow SMTP call cannot hold the response open.
func (h *AuthHandler) deliverVerificationEmail(user models.User) {
	token, err := services.GenerateVerificationToken()
	if err != nil {
		log.Printf("email verification: generating token for %s: %v", user.Email, err)
		return
	}

	if _, err := services.CreateEmailVerificationToken(h.DB, user.ID, user.Email, token); err != nil {
		log.Printf("email verification: storing token for %s: %v", user.Email, err)
		return
	}

	verifyURL := strings.TrimSuffix(h.Config.OAuthFrontendURL, "/") + "/verify-email?token=" + url.QueryEscape(token)

	if h.Mailer != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		if err := h.Mailer.Send(ctx, mail.SendOptions{
			To:       user.Email,
			Subject:  "Confirm your email address",
			Template: "email-verification",
			Data: map[string]interface{}{
				"AppName":   h.Config.AppName,
				"Title":     "Confirm your email address",
				"Message":   "Click the button below to confirm this address. The link expires in 48 hours and can only be used once.",
				"VerifyURL": verifyURL,
				"Year":      time.Now().Year(),
			},
		}); err != nil {
			log.Printf("email verification: sending to %s: %v", user.Email, err)
		}
		return
	}

	if h.Config.AppEnv == "production" {
		log.Printf("email verification: NO MAILER CONFIGURED: %s cannot receive a link. Set RESEND_API_KEY.", user.Email)
		return
	}

	log.Printf("email verification link for %s: %s", user.Email, verifyURL)
}

// ResetPassword resets a user's password with a valid token.
func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var req ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": gin.H{
				"code":    "VALIDATION_ERROR",
				"message": err.Error(),
			},
		})
		return
	}

	// Consume first. The token is single-use and burning it before doing any
	// work means a failure later can't leave a still-valid token behind.
	userID, err := services.ConsumePasswordResetToken(h.DB, req.Token)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_TOKEN",
				"message": "This reset link is invalid or has expired. Request a new one.",
			},
		})
		return
	}

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

	if err := h.DB.Model(&models.User{}).Where("id = ?", userID).
		Update("password", string(hashedPassword)).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to update password",
			},
		})
		return
	}

	// The reason someone resets a password is to evict whoever they think is in
	// their account. Leaving that person's session alive would defeat the entire
	// exercise, so every device is signed out — including any the attacker holds.
	if err := services.RevokeAllUserSessions(h.DB, userID, ""); err != nil {
		log.Printf("password reset: revoking sessions for %s: %v", userID, err)
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Password reset successfully. Please sign in with your new password.",
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

	// Record the refresh token as a server-side session, so an OAuth login is
	// listed and revocable exactly like a password login.
	if _, err := services.CreateSession(h.DB, c, user.ID, tokens.RefreshToken); err != nil {
		log.Printf("OAuth: failed to record session for %s: %v", user.ID, err)
	}

	// Set HttpOnly auth cookies BEFORE redirecting so the browser stores
	// them as part of this same response. The callback page then just
	// navigates — no tokens in URL, no tokens in JS, no XSS exposure.
	h.AuthService.SetAuthCookies(c, tokens)

	// Redirect to frontend callback. No query params — tokens travel as
	// HttpOnly Set-Cookie headers on this 307 response.
	redirectURL := fmt.Sprintf("%s/auth/callback", h.Config.OAuthFrontendURL)
	c.Redirect(http.StatusTemporaryRedirect, redirectURL)
}
`
}

func apiUserHandlerGo() string {
	return `package handlers

import (
	"log"
	"math"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"` + "{{MODULE}}" + `/internal/authz"
	"` + "{{MODULE}}" + `/internal/crypto"
	"` + "{{MODULE}}" + `/internal/models"
	"` + "{{MODULE}}" + `/internal/services"
)

// UserHandler handles user management endpoints.
type UserHandler struct {
	DB *gorm.DB
	// AuthService is used by DeleteProfile to clear the HttpOnly auth
	// cookies as part of the soft-delete response. Optional — if nil,
	// the cookies just won't be cleared and the client's next request
	// will 401 normally.
	AuthService *services.AuthService
}

// A new user.
type CreateUserRequest struct {
	FirstName string ` + "`" + `json:"first_name" binding:"required"` + "`" + `
	LastName  string ` + "`" + `json:"last_name" binding:"required"` + "`" + `
	Email     string ` + "`" + `json:"email" binding:"required,email"` + "`" + `
	Password  string ` + "`" + `json:"password" binding:"required,min=6"` + "`" + `
	Role      string ` + "`" + `json:"role"` + "`" + `
	Avatar    string ` + "`" + `json:"avatar"` + "`" + `
	JobTitle  string ` + "`" + `json:"job_title"` + "`" + `
	Active    *bool  ` + "`" + `json:"active"` + "`" + `
}

// Create creates a new user (admin only).

func (h *UserHandler) Create(c *gin.Context) {
	var req CreateUserRequest

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
		// LOWER(...) LIKE LOWER(...) rather than ILIKE: ILIKE is Postgres-only
		// and this API also runs on SQLite.
		query = query.Where("LOWER(first_name) LIKE LOWER(?) OR LOWER(last_name) LIKE LOWER(?) OR LOWER(email) LIKE LOWER(?)", "%"+search+"%", "%"+search+"%", "%"+search+"%")
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

// syncUserRoleAssignment makes the user_roles table reflect a single role name.
//
// The admin UI edits a user's role as one string, while authorization resolves
// through the many-to-many user_roles table. This bridges the two so the simple
// dropdown keeps working and actually takes effect.
//
// Assign several roles to one user with PUT /api/users/:id/roles instead — that
// endpoint is the full many-to-many path and this helper is its one-role case.
//
// A role name with no matching row (a custom legacy string) clears the
// assignments and lets grant resolution fall back to users.role, rather than
// failing the update.
func syncUserRoleAssignment(db *gorm.DB, userID, roleName string) error {
	var role models.Role
	err := db.Where("name = ?", roleName).First(&role).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}

	txErr := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("user_id = ?", userID).Delete(&models.UserRole{}).Error; err != nil {
			return err
		}
		if role.ID == "" {
			return nil // unknown name — fall back to the legacy string
		}
		return tx.Create(&models.UserRole{UserID: userID, RoleID: role.ID}).Error
	})
	if txErr != nil {
		return txErr
	}

	// Permissions just changed for this user; drop the cached grants.
	authz.Invalidate()
	return nil
}

// Changes to a user.
type UpdateUserRequest struct {
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

	var req UpdateUserRequest

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
		updates["bio"] = crypto.EncryptedString(req.Bio)
	}
	if req.Active != nil {
		updates["active"] = *req.Active
	}

	// Keep role ASSIGNMENTS in step with the role string.
	//
	// Grant resolution prefers the user_roles table and only falls back to
	// users.role. Without this, changing the Role dropdown for a user who
	// already has an assignment would update the string and change nothing
	// about what they can actually do — a silent no-op, and a nasty one to
	// debug.
	if req.Role != "" {
		if err := syncUserRoleAssignment(h.DB, user.ID, req.Role); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": gin.H{
					"code":    "INTERNAL_ERROR",
					"message": "Failed to update role assignment",
				},
			})
			return
		}
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

// Changes to your own profile.
type UpdateProfileRequest struct {
	FirstName string ` + "`" + `json:"first_name"` + "`" + `
	LastName  string ` + "`" + `json:"last_name"` + "`" + `
	Email     string ` + "`" + `json:"email"` + "`" + `
	Password  string ` + "`" + `json:"password"` + "`" + `
	Avatar    string ` + "`" + `json:"avatar"` + "`" + `
	JobTitle  string ` + "`" + `json:"job_title"` + "`" + `
	Bio       string ` + "`" + `json:"bio"` + "`" + `
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

	var req UpdateProfileRequest

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
	passwordChanged := false
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
		passwordChanged = true
	}
	if req.Avatar != "" {
		updates["avatar"] = req.Avatar
	}
	if req.JobTitle != "" {
		updates["job_title"] = req.JobTitle
	}
	if req.Bio != "" {
		updates["bio"] = crypto.EncryptedString(req.Bio)
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

	// Changing a password must invalidate every logged-in device — that is the
	// whole point of changing it after a suspected compromise.
	//
	// The caller is then re-issued a brand-new session rather than spared: the
	// grit_refresh cookie is scoped to /api/auth, so a PUT /api/profile never
	// carries it and there is no way to recognise "this device" here. Revoking
	// everything and minting a fresh pair is both simpler and stricter — the old
	// token is dead even for the caller, and they stay signed in.
	if passwordChanged {
		if err := services.RevokeAllUserSessions(h.DB, user.ID, ""); err != nil {
			// Log it; the password DID change, so failing the request now would be
			// misleading. Sessions still die at their idle/absolute timeout.
			log.Printf("failed to revoke sessions after password change for user %s: %v", user.ID, err)
		} else if h.AuthService != nil {
			pair, terr := h.AuthService.GenerateTokenPair(user.ID, user.Email, user.Role)
			if terr != nil {
				log.Printf("failed to re-issue tokens after password change for user %s: %v", user.ID, terr)
			} else if _, serr := services.CreateSession(h.DB, c, user.ID, pair.RefreshToken); serr != nil {
				log.Printf("failed to open a session after password change for user %s: %v", user.ID, serr)
			} else {
				h.AuthService.SetAuthCookies(c, pair)
			}
		}
	}

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

	// Soft-delete leaves the JWT valid in theory; the auth middleware
	// would still 401 on the next request because the user row is
	// excluded by the default scope. We still expire the HttpOnly auth
	// cookies on the way out so the next /api/* call from this browser
	// doesn't even attempt — saves a round trip and a confusing 401 in
	// the dev console.
	if h.AuthService != nil {
		h.AuthService.ClearAuthCookies(c)
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

	"` + "{{MODULE}}" + `/internal/authz"
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
		c.Set("user_email", user.Email)
		c.Set("user_role", user.Role)

		// Resolve the caller's permission grants once per request so route
		// guards and handlers don't each hit the database. authz.GrantsFor is
		// cached and invalidated on role changes, so this is usually free.
		// A failure here is not fatal: the request continues with no grants and
		// role-name checks still apply, which fails closed rather than 500ing
		// every route the moment the roles table has a problem.
		if grants, err := authz.GrantsFor(db, user.ID); err == nil {
			c.Set("user_grants", grants)
		}

		c.Next()
	}
}

// RequireRole guards a route by role name, permission, or both.
//
// Each argument is either a legacy role name ("ADMIN") or a permission key
// prefixed with "perm:" ("perm:users.delete"). Access is granted if ANY
// argument matches — so the two styles can be mixed during a migration:
//
//	protected.Use(middleware.RequireRole("ADMIN", "perm:users.delete"))
//
// The signature is unchanged on purpose: every existing RequireRole("ADMIN")
// call site keeps working untouched, and permissions can be adopted route by
// route instead of in one breaking sweep.
func RequireRole(rolesOrPerms ...string) gin.HandlerFunc {
	// Split once at construction rather than per request.
	var roles, perms []string
	for _, arg := range rolesOrPerms {
		if strings.HasPrefix(arg, "perm:") {
			perms = append(perms, strings.TrimPrefix(arg, "perm:"))
			continue
		}
		roles = append(roles, arg)
	}

	return func(c *gin.Context) {
		// Permission check first — it's the model we want callers to move to.
		if len(perms) > 0 {
			if grants, ok := c.Get("user_grants"); ok {
				if list, ok := grants.([]string); ok {
					for _, p := range perms {
						if authz.Granted(list, p) {
							c.Next()
							return
						}
					}
				}
			}
		}

		// No permission matched; fall back to the legacy role names.
		if len(roles) == 0 {
			c.JSON(http.StatusForbidden, gin.H{
				"error": gin.H{
					"code":    "FORBIDDEN",
					"message": "You do not have permission to perform this action",
				},
			})
			c.Abort()
			return
		}

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
	"net/url"

	"github.com/gin-gonic/gin"
)

// isWailsOrigin reports whether the request came from the Wails desktop
// webview, whose origin is not stably enumerable:
//
//	wails dev   (Windows)     http://wails.localhost:34115   <- port from wails.json
//	wails build (Windows)     http://wails.localhost
//	wails build (mac/linux)   wails://wails
//
// The dev-server port is configurable, so pinning exact origins in
// CORS_ORIGINS is fragile: change the port and the desktop login silently
// starts failing with an opaque "Network Error". Match on the host instead.
//
// Safe by construction: "wails.localhost" is a virtual host the webview
// resolves internally, so a page on the public internet cannot be served
// from it and cannot forge this origin. Every other origin still has to be
// in the explicit CORS_ORIGINS allowlist.
func isWailsOrigin(origin string) bool {
	if origin == "wails://wails" {
		return true
	}
	u, err := url.Parse(origin)
	if err != nil {
		return false
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return false
	}
	return u.Hostname() == "wails.localhost"
}

// CORS creates a CORS middleware with a fixed allowlist.
//
// Kept for callers that genuinely have a static list. Anything user-facing
// should use CORSDynamic, so adding a domain does not need a redeploy.
func CORS(allowedOrigins []string) gin.HandlerFunc {
	return CORSDynamic(func() []string { return allowedOrigins })
}

// CORSDynamic resolves the allowlist per request.
//
// Per request rather than at construction, because the point of putting
// origins in settings is that somebody adds a domain at 9pm and it works.
// Capturing the slice at boot would mean the setting existed and did nothing
// until the next deploy, which is worse than not offering it.
//
// The cost is a map build per request over a list with single digits of
// entries, against a settings store that is already cached in memory. That is
// not the thing to optimise.
func CORSDynamic(resolve func() []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")

		allowed := false
		for _, candidate := range resolve() {
			if candidate == origin && origin != "" {
				allowed = true
				break
			}
		}
		if allowed || isWailsOrigin(origin) {
			c.Header("Access-Control-Allow-Origin", origin)
		}

		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		// X-CSRF-Token + Idempotency-Key are injected by the web and admin
		// axios clients on every unsafe method. Without them in the allowed
		// list, the browser's preflight strips the headers and the request
		// either fails the AutoCSRF check or replays without an idempotency
		// guarantee. Authorization stays for native bearer clients.
		// X-API-Key is not optional here. A storefront calls the public
		// endpoints with it, cross-origin, and a header missing from this list
		// is stripped by the browser during preflight: the request fails in
		// every browser and works perfectly under curl, which is the worst
		// shape a bug can have.
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization, X-API-Key, X-CSRF-Token, Idempotency-Key, X-Public-IP-Hint")
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
// page-clamping, sort whitelisting, and search-clause construction live
// in exactly one place. Addresses issue #14.
// APIPaginateGo is exported so grit generate resource can bring an older
// project's paginate package forward. A generated public handler declares
// RangeFilterable, which a paginate.Config that predates it does not have, and
// the project then fails to compile on a field name.
func APIPaginateGo() string { return apiPaginateGo() }

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
	"gorm.io/gorm/schema"
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
	// QueryFilters are raw query-string params that are NOT reserved words.
	// Untrusted: kept apart from Filters (which handlers set in Go) so the
	// whitelist can be applied to these and only these.
	QueryFilters map[string]string

	// v3.31.34 — date filter. DateField is the column (default
	// "created_at"); DateFrom/DateTo are inclusive bounds. When both
	// are zero values, no date filter is applied. Set via the
	// ?created_from=/?created_to= query params (or the legacy
	// ?created_since=Nd shortcut used by the stat cards).
	DateField string
	DateFrom  time.Time
	DateTo    time.Time
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
	Searchable   []string        // columns included in case-insensitive search
	Sortable     map[string]bool // whitelist for sort_by values

	// Filterable whitelists columns that may be filtered from the query
	// string, so ?status=pending becomes WHERE status = 'pending'.
	//
	// A whitelist and not a free-for-all, because the column name is
	// interpolated into the SQL: without it, ?"1=1 OR x"= would be a query
	// the caller wrote. Anything not listed here is ignored rather than
	// rejected, so an unknown param is never an error.
	Filterable map[string]bool

	// RangeFilterable whitelists numeric columns that accept a window from
	// the query string: ?price_min=50&price_max=200 becomes
	// WHERE price >= 50 AND price <= 200.
	//
	// Separate from Filterable because the two answer different questions and
	// a column often wants one and not the other. Equality on a price is
	// almost never what a caller means; a range on a status is meaningless.
	// Same whitelist reasoning as above: the column name reaches the WHERE
	// clause, and the bound is parameterised.
	//
	// Either bound may be omitted. A bound that does not parse as a number is
	// ignored rather than rejected, so a hand-typed URL degrades to a wider
	// result set instead of an error page.
	RangeFilterable map[string]bool

	// InFilterable whitelists columns that accept a comma-separated list:
	// ?category_id=a,b,c becomes WHERE category_id IN (a,b,c).
	//
	// Only for id columns, and that is a deliberate restriction rather than an
	// accident of naming. Splitting on commas is wrong for anything a person
	// types, because "Smith, John" is one value; it is safe for ids, where a
	// comma never appears. So this is opt-in per column and never inferred from
	// the value.
	//
	// The case it exists for: a category tree. "Products in Electronics" means
	// Electronics and every category under it, which is a list of ids the client
	// already has from the category it fetched.
	InFilterable map[string]bool

	DefaultSort  string // fallback sort column (defaults to "created_at")
	DefaultOrder string // fallback sort order (defaults to "desc")

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
	// No omitempty on the counts. Zero is an answer: an empty result set that
	// reports {"page":1,"page_size":20} and no total leaves every client doing
	// meta.total with undefined, which renders as a blank stat card rather
	// than a nought and turns arithmetic into NaN.
	Total      int64  ` + "`" + `json:"total"` + "`" + `
	Page       int    ` + "`" + `json:"page"` + "`" + `
	PageSize   int    ` + "`" + `json:"page_size"` + "`" + `
	Pages      int    ` + "`" + `json:"pages"` + "`" + `
	NextCursor string ` + "`" + `json:"next_cursor,omitempty"` + "`" + `
	HasMore    bool   ` + "`" + `json:"has_more,omitempty"` + "`" + `
}

// Result wraps the paginated data in the canonical { data, meta } envelope.
type Result[T any] struct {
	Data []T  ` + "`" + `json:"data"` + "`" + `
	Meta Meta ` + "`" + `json:"meta"` + "`" + `
}

// coerceFilterValue turns a query-string value into something the column can
// actually be compared against.
//
// Only booleans need this, and only because the two databases disagree about
// what to do with a string. Postgres reads WHERE active = 'true' as a boolean
// and answers correctly; MySQL stores the column as tinyint(1), coerces the
// non-numeric string to 0, and quietly returns nothing. Neither errors, so the
// bug is a filter that silently matches no rows.
//
// Numbers stay strings on purpose: both databases coerce those identically,
// and a varchar column holding "12" would break if this second-guessed it.
func coerceFilterValue(val string, dataType schema.DataType) any {
	if dataType != schema.Bool {
		return val
	}
	switch strings.ToLower(val) {
	case "true", "1", "yes", "on":
		return true
	case "false", "0", "no", "off":
		return false
	}
	return val
}


// Bind reads page / page_size / search / sort_by / sort_order from the Gin
// context, clamps them, and returns a normalized Params.
//
// v3.31.34 — also parses the date-filter query params:
//   ?created_from=2026-01-01&created_to=2026-12-31
//   ?created_since=7d   (legacy shortcut: last 7 days)
//   ?date_field=published_at   (override the default "created_at" column)
//
// Both _from and _to are inclusive. Dates without time components are
// snapped to the start (00:00) for _from and end (23:59:59) for _to.
func Bind(c *gin.Context) Params {
	page, _ := strconv.Atoi(c.DefaultQuery("page", strconv.Itoa(DefaultPage)))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", strconv.Itoa(DefaultPageSize)))

	if page < 1 {
		page = DefaultPage
	}
	if pageSize < 1 || pageSize > MaxPageSize {
		pageSize = DefaultPageSize
	}

	dateField := c.Query("date_field")
	if dateField == "" {
		dateField = "created_at"
	}
	dateFrom, dateTo := parseDateRange(c)

	return Params{
		Page:         page,
		PageSize:     pageSize,
		Search:       c.Query("search"),
		SortBy:       c.Query("sort_by"),
		SortOrder:    c.Query("sort_order"),
		Cursor:       c.Query("cursor"),
		DateField:    dateField,
		DateFrom:     dateFrom,
		DateTo:       dateTo,
		QueryFilters: collectQueryFilters(c),
	}
}

// reservedParams are the query keys pagination owns. Everything else is a
// candidate column filter, to be checked against Config.Filterable before it
// is used.
var reservedParams = map[string]bool{
	"page": true, "page_size": true, "search": true,
	"sort_by": true, "sort_order": true, "cursor": true,
	"date_field": true, "date_from": true, "date_to": true,
	"created_since": true, "created_from": true, "created_to": true,
	"updated_since": true, "archived": true, "format": true,
}

func collectQueryFilters(c *gin.Context) map[string]string {
	out := map[string]string{}
	for key, values := range c.Request.URL.Query() {
		if reservedParams[key] || len(values) == 0 || values[0] == "" {
			continue
		}
		out[key] = values[0]
	}
	if len(out) == 0 {
		return nil
	}
	return out
}

// parseDateRange reads the three supported date-window query params and
// returns the resolved (from, to) bounds. Zero values mean "no bound".
//
// Precedence: explicit created_from/created_to wins over created_since.
// This lets a UI date picker override a stat card's "last 7 days" link
// without surprising clobber.
func parseDateRange(c *gin.Context) (time.Time, time.Time) {
	var from, to time.Time
	if since := c.Query("created_since"); since != "" {
		if d, ok := parseRelativeDuration(since); ok {
			from = time.Now().Add(-d)
		}
	}
	if s := c.Query("created_from"); s != "" {
		if t, err := parseDateInput(s, false); err == nil {
			from = t
		}
	}
	if s := c.Query("created_to"); s != "" {
		if t, err := parseDateInput(s, true); err == nil {
			to = t
		}
	}
	return from, to
}

// parseRelativeDuration parses "7d", "30d", "12h", "1w" into a
// time.Duration. Used by the stat-card shortcut ?created_since=7d so
// hand-written URLs stay short. Returns ok=false on unrecognised input
// (caller falls back to no bound rather than failing the request).
func parseRelativeDuration(s string) (time.Duration, bool) {
	if len(s) < 2 {
		return 0, false
	}
	unit := s[len(s)-1]
	nStr := s[:len(s)-1]
	n, err := strconv.Atoi(nStr)
	if err != nil || n < 0 {
		return 0, false
	}
	switch unit {
	case 'h':
		return time.Duration(n) * time.Hour, true
	case 'd':
		return time.Duration(n) * 24 * time.Hour, true
	case 'w':
		return time.Duration(n) * 7 * 24 * time.Hour, true
	case 'm':
		// month = 30 days. Good enough for stats; calendar-accurate
		// month math isn't worth the dep.
		return time.Duration(n) * 30 * 24 * time.Hour, true
	}
	return 0, false
}

// parseDateInput parses an ISO date or datetime string. If endOfDay is
// true and the input is a bare date, it snaps to 23:59:59.999 so the
// _to bound is inclusive of the whole day the user picked.
func parseDateInput(s string, endOfDay bool) (time.Time, error) {
	if t, err := time.Parse(time.RFC3339, s); err == nil {
		return t, nil
	}
	if t, err := time.Parse("2006-01-02", s); err == nil {
		if endOfDay {
			return t.Add(24*time.Hour - time.Nanosecond), nil
		}
		return t, nil
	}
	return time.Time{}, fmt.Errorf("invalid date %q", s)
}

// isSafeDateColumn reports whether the client-supplied date_field column is
// safe to interpolate into a WHERE fragment. Only bare identifiers that are
// either a known timestamp column or an explicitly-declared sortable column
// are allowed; everything else is rejected so the caller falls back to
// "created_at". This is the guard behind the date-filter injection fix.
func isSafeDateColumn(col string, cfg Config) bool {
	if col == "" {
		return false
	}
	// Reject anything that isn't a plain snake_case identifier up front —
	// no spaces, parens, quotes, or SQL operators can survive this.
	for _, r := range col {
		if !(r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' || r >= '0' && r <= '9' || r == '_') {
			return false
		}
	}
	if col == "created_at" || col == "updated_at" {
		return true
	}
	return cfg.Sortable[col]
}

// List runs the query with search / sort / filters / pagination applied and
// returns a typed Result[T]. The caller is expected to have already set the
// model and any relevant Preload() chains on the passed-in *gorm.DB.
//
// Invariants enforced:
//   - page >= 1, 1 <= page_size <= MaxPageSize
//   - sort_by must be in cfg.Sortable, else cfg.DefaultSort (or DefaultSortColumn)
//   - sort_order must be "asc" or "desc", else cfg.DefaultOrder (or DefaultSortOrder)
//   - search is applied case-insensitively across cfg.Searchable columns (nothing if empty)
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

	// Equality filters set by the handler in Go. Trusted: the column names
	// are literals in our own source.
	for col, val := range p.Filters {
		query = query.Where(col+" = ?", val)
	}

	// Equality filters from the query string (?status=pending). Untrusted, so
	// every column is checked against the whitelist before it reaches the SQL.
	// This is what makes the admin's filter dropdowns and tab strips work:
	// before it existed, Bind never collected them and the comment above
	// promised a feature the code did not have.
	if len(p.QueryFilters) > 0 {
		// The model's schema is parsed once so each value can be bound as the
		// column's real type. See coerceFilterValue for why that matters.
		var zero T
		stmt := &gorm.Statement{DB: query}
		parsed := stmt.Parse(&zero) == nil

		for col, val := range p.QueryFilters {
			if !cfg.Filterable[col] {
				continue
			}
			var bound any = val
			if parsed && stmt.Schema != nil {
				if f := stmt.Schema.LookUpField(col); f != nil {
					bound = coerceFilterValue(val, f.DataType)
				}
			}
			query = query.Where(col+" = ?", bound)
		}

		// Id lists. Same untrusted map, same whitelist rule.
		for col := range cfg.InFilterable {
			raw, ok := p.QueryFilters[col]
			if !ok || raw == "" {
				continue
			}
			parts := strings.Split(raw, ",")
			values := make([]string, 0, len(parts))
			for _, part := range parts {
				if trimmed := strings.TrimSpace(part); trimmed != "" {
					values = append(values, trimmed)
				}
			}
			if len(values) == 0 {
				continue
			}
			query = query.Where(col+" IN ?", values)
		}

		// Range windows. Read from the same untrusted QueryFilters map, so a
		// column has to be declared RangeFilterable before "price_min" can
		// become a WHERE on price.
		for col := range cfg.RangeFilterable {
			for suffix, op := range map[string]string{"_min": ">=", "_max": "<="} {
				raw, ok := p.QueryFilters[col+suffix]
				if !ok || raw == "" {
					continue
				}
				bound, err := strconv.ParseFloat(raw, 64)
				if err != nil {
					// A bound nobody can parse widens the window rather than
					// failing the request. ?price_min=cheap is a typo, not an
					// attack, and an error page is a worse answer than results.
					continue
				}
				query = query.Where(col+" "+op+" ?", bound)
			}
		}
	}

	// v3.31.34 — date-window filter. DateField defaults to "created_at"
	// in Bind() so we only need to apply when at least one bound is set.
	//
	// SECURITY (v3.31.84): DateField arrives straight from the ?date_field=
	// query param, so it MUST be whitelisted before it touches a WHERE
	// fragment — GORM treats the condition string as raw SQL. We allow the
	// two always-present timestamp columns plus anything the resource
	// already declared sortable (a strict, developer-controlled set), and
	// fall back to "created_at" on anything else. This closes the
	// date_field SQL-injection vector reachable by any authenticated user.
	dateField := p.DateField
	if !isSafeDateColumn(dateField, cfg) {
		dateField = "created_at"
	}
	if !p.DateFrom.IsZero() {
		query = query.Where(dateField+" >= ?", p.DateFrom)
	}
	if !p.DateTo.IsZero() {
		query = query.Where(dateField+" <= ?", p.DateTo)
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

// buildSearchClause builds "LOWER(col1) LIKE LOWER(?) OR ..." with the
// same wildcard-wrapped search term repeated as each arg.
func buildSearchClause(cols []string, term string) (string, []any) {
	clause := ""
	args := make([]any, 0, len(cols))
	wild := "%" + term + "%"
	for i, col := range cols {
		if i > 0 {
			clause += " OR "
		}
		clause += "LOWER(" + col + ") LIKE LOWER(?)"
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
	"gorm.io/gorm/clause"

	"` + "{{MODULE}}" + `/internal/services"
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

// Changes made while offline.
type SyncPushRequest struct {
	Changes []PushChange ` + "`" + `json:"changes"` + "`" + `
}

// Push handles POST /api/sync/push. Each change is applied
// independently — one conflict does not abort the rest of the batch.

func (h *SyncHandler) Push(c *gin.Context) {
	var req SyncPushRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "INVALID_BODY", "message": err.Error()}})
		return
	}

	results := make([]PushResult, len(req.Changes))
	for i, ch := range req.Changes {
		results[i] = h.applyChange(c, ch)
	}
	c.JSON(http.StatusOK, gin.H{"results": results})
}

// syncIdentifier picks a human-friendly label for the semantic activity feed
// from a change payload (name / title / slug / email), falling back to the id.
func syncIdentifier(data map[string]interface{}, id string) string {
	for _, k := range []string{"name", "title", "slug", "email"} {
		if v, ok := data[k]; ok {
			if s, ok := v.(string); ok && s != "" {
				return s
			}
		}
	}
	return id
}

func (h *SyncHandler) applyChange(c *gin.Context, ch PushChange) PushResult {
	proto, err := h.Registry.New(ch.Model)
	if err != nil {
		return PushResult{OK: false, Code: "UNKNOWN_MODEL", Message: err.Error()}
	}
	policy := h.Registry.PolicyFor(ch.Model)
	// Enforced before the payload reaches a decoder: a field declared
	// local_only is a promise it never leaves the device, and a promise kept
	// only by well-behaved clients is not one.
	ch.Data = policy.StripLocalOnly(ch.Data)
	// The Go struct name (e.g. "Category") is the nicest entity label for the
	// activity feed — offline edits should read the same as online ones.
	entityType := reflect.TypeOf(proto).Elem().Name()

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
		// Mirror the online handler: emit a semantic activity row so offline
		// creates surface in /system/activity, not just the raw audit log.
		services.LogCreate(h.DB, c, entityType, syncIdentifier(ch.Data, ch.ID), ch.ID, "")
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
			// The versions disagree. What happens next is the resource's
			// declared policy, decided here rather than in the client, because
			// a rule an old build can ignore is not a rule.
			switch policy.Conflict {
			case sync.ConflictServerWins:
				// The client's change is dropped and it is told to take the
				// server row. Reported distinctly from a conflict so the
				// client can apply it without asking anyone.
				return PushResult{
					OK:            false,
					Code:          "SERVER_WINS",
					Message:       fmt.Sprintf("server v%d kept over client v%d", serverVersion, ch.Version),
					ServerVersion: serverVersion,
					ServerData:    current,
				}
			case sync.ConflictClientWins:
				// Fall through and overwrite. The version check was protecting
				// nothing for this resource, and the author said so.
			default:
				return PushResult{
					OK:            false,
					Code:          "VERSION_CONFLICT",
					Message:       fmt.Sprintf("client had v%d, server has v%d", ch.Version, serverVersion),
					ServerVersion: serverVersion,
					ServerData:    current,
				}
			}
		}
		// Apply the update.
		//
		// Decode the client payload into a fresh, typed model struct rather than
		// calling .Updates(ch.Data) directly. A raw map[string]interface{} hands
		// nested values (a FileRef image, a FileRefs slice, a belongs-to relation
		// object) straight to the DB driver, which cannot encode a Go map into a
		// json column ("cannot find encode plan for OID 0") — the update fails and
		// the offline outbox entry gets stuck forever. Decoding first routes those
		// fields through their driver.Valuer implementations, exactly like create.
		obj := proto
		if err := decodeInto(obj, ch.Data); err != nil {
			return PushResult{OK: false, Code: "DECODE_ERROR", Message: err.Error()}
		}
		setField(obj, "ID", ch.ID)
		// Seed Version with the server's value so the BeforeUpdate hook bumps it to
		// serverVersion+1 regardless of what the client sent in the payload.
		setIntField(obj, "Version", serverVersion)
		// Save writes every column (so cleared/zeroed fields persist) and runs the
		// BeforeUpdate hook. Omit associations so the nested relation object is not
		// upserted, and CreatedAt so the client can't rewind the original timestamp.
		if err := h.DB.Omit(clause.Associations, "CreatedAt").Save(obj).Error; err != nil {
			return PushResult{OK: false, Code: "UPDATE_FAILED", Message: err.Error()}
		}
		newVersion := getIntField(obj, "Version")
		services.LogUpdate(h.DB, c, entityType, syncIdentifier(ch.Data, ch.ID), ch.ID, services.DiffSummary(ch.Data))
		return PushResult{OK: true, NewVersion: newVersion}

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
		services.LogDelete(h.DB, c, entityType, ch.ID, ch.ID)
		return PushResult{OK: true}

	default:
		return PushResult{OK: false, Code: "INVALID_OP", Message: "op must be create, update, or delete"}
	}
}

// SyncPolicyResponse is what GET /api/sync/policy answers with.
type SyncPolicyResponse struct {
	Models map[string]sync.Policy ` + "`" + `json:"models"` + "`" + `
}

// Policy handles GET /api/sync/policy.
//
// Clients configure themselves from this rather than from a generated copy.
// A copy is a second thing to keep in sync, and an offline client running
// last month's build against this month's conflict rules is exactly the
// silent failure this whole feature exists to prevent.
func (h *SyncHandler) Policy(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"data": SyncPolicyResponse{Models: h.Registry.Policies()}})
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
	policy := h.Registry.PolicyFor(model)
	if policy.Mode == sync.ModeOnlineOnly {
		// Answering with an empty page would look like "nothing has changed",
		// and the client would mirror an empty table forever.
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{
			"code":    "NOT_SYNCABLE",
			"message": model + " is registered online_only and is not mirrored",
		}})
		return
	}

	// Build a slice of the right type via reflection.
	sliceType := reflect.SliceOf(reflect.TypeOf(proto).Elem())
	results := reflect.New(sliceType)

	// Effective change time = the LATER of updated_at and deleted_at. A soft
	// delete only sets deleted_at, so ordering/cursoring on updated_at alone
	// would never carry the delete to offline clients (they'd keep a ghost
	// row forever). We order + cursor on the effective time and mark deleted
	// rows with "_deleted": true so the client can drop them from its mirror.
	effExpr := "MAX(updated_at, COALESCE(deleted_at, updated_at))"
	if h.DB.Dialector.Name() == "postgres" {
		effExpr = "GREATEST(updated_at, COALESCE(deleted_at, updated_at))"
	}

	// Unscoped so soft-deleted rows are included (they're the tombstones).
	q := h.DB.Unscoped().Model(proto)
	if sinceStr != "" {
		t, err := time.Parse(time.RFC3339Nano, sinceStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "INVALID_SINCE", "message": err.Error()}})
			return
		}
		q = q.Where("updated_at > ? OR deleted_at > ?", t, t)
	}
	if err := q.Order(effExpr + " asc").Limit(limit).Find(results.Interface()).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"code": "INTERNAL_ERROR", "message": err.Error()}})
		return
	}

	rs := results.Elem()
	rows := make([]map[string]interface{}, 0, rs.Len())
	cursor := sinceStr
	var maxEff time.Time
	for i := 0; i < rs.Len(); i++ {
		item := rs.Index(i).Addr().Interface()
		b, merr := json.Marshal(item)
		if merr != nil {
			continue
		}
		var m map[string]interface{}
		if uerr := json.Unmarshal(b, &m); uerr != nil {
			continue
		}
		m["_deleted"] = isSyncDeleted(item)
		rows = append(rows, policy.Projects(m))
		if eff, ok := effectiveSyncTime(item); ok && eff.After(maxEff) {
			maxEff = eff
		}
	}
	if !maxEff.IsZero() {
		cursor = maxEff.Format(time.RFC3339Nano)
	}

	c.JSON(http.StatusOK, gin.H{
		"data":   rows,
		"cursor": cursor,
		"count":  len(rows),
	})
}

// isSyncDeleted reports whether a model row is soft-deleted (a tombstone).
func isSyncDeleted(obj interface{}) bool {
	v := reflect.ValueOf(obj)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	f := v.FieldByName("DeletedAt")
	if !f.IsValid() {
		return false
	}
	if d, ok := f.Interface().(gorm.DeletedAt); ok {
		return d.Valid
	}
	return false
}

// effectiveSyncTime returns the later of a row's UpdatedAt and DeletedAt — the
// timestamp the pull cursor advances on so both edits and deletes are carried.
func effectiveSyncTime(obj interface{}) (time.Time, bool) {
	v := reflect.ValueOf(obj)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	var eff time.Time
	if f := v.FieldByName("UpdatedAt"); f.IsValid() {
		if t, ok := f.Interface().(time.Time); ok {
			eff = t
		}
	}
	if f := v.FieldByName("DeletedAt"); f.IsValid() {
		if d, ok := f.Interface().(gorm.DeletedAt); ok && d.Valid && d.Time.After(eff) {
			eff = d.Time
		}
	}
	return eff, !eff.IsZero()
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

// setIntField sets an int field on a struct via reflection. Used to seed
// Version before an update save so the BeforeUpdate hook bumps from the
// server's value rather than whatever the client happened to send.
func setIntField(obj interface{}, name string, value int) {
	v := reflect.ValueOf(obj).Elem()
	f := v.FieldByName(name)
	if f.IsValid() && f.CanSet() && f.CanInt() {
		f.SetInt(int64(value))
	}
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

	"{{MODULE}}/internal/ids"
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
		a.ID = ids.New()
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

	"{{MODULE}}/internal/ids"
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
	// No explicit type: datatypes.JSON maps to jsonb on Postgres and json on
	// MySQL by itself. Naming jsonb here fails AutoMigrate on MySQL, which has
	// no such type.
	Payload      datatypes.JSON ` + "`" + `json:"payload"` + "`" + `
	Status       string         ` + "`" + `gorm:"size:20;index;not null;default:pending" json:"status"` + "`" + `
	HandlerError string         ` + "`" + `gorm:"type:text" json:"handler_error,omitempty"` + "`" + `
	RetryCount   int            ` + "`" + `gorm:"not null;default:0" json:"retry_count"` + "`" + `
	ProcessedAt  *time.Time     ` + "`" + `json:"processed_at,omitempty"` + "`" + `
	CreatedAt    time.Time      ` + "`" + `gorm:"index" json:"created_at"` + "`" + `
}

func (w *WebhookEvent) BeforeCreate(tx *gorm.DB) error {
	if w.ID == "" {
		w.ID = ids.New()
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

	"{{MODULE}}/internal/ids"
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
	Rules       datatypes.JSON ` + "`" + `json:"rules"` + "`" + `
	CreatedAt   time.Time      ` + "`" + `json:"created_at"` + "`" + `
	UpdatedAt   time.Time      ` + "`" + `json:"updated_at"` + "`" + `
	Version     int            ` + "`" + `gorm:"not null;default:1" json:"version"` + "`" + `
}

func (f *FeatureFlag) BeforeCreate(tx *gorm.DB) error {
	if f.ID == "" {
		f.ID = ids.New()
	}
	return nil
}

func (f *FeatureFlag) BeforeUpdate(tx *gorm.DB) error {
	tx.Statement.SetColumn("version", gorm.Expr("version + 1"))
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
		e.ID = ids.New()
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

// FeatureFlagRequest is the request shape for create/update. Rules is taken
// as a structured object — the handler encodes to JSON before hitting
// the DB so the wire format stays consistent.
type FeatureFlagRequest struct {
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
	var body FeatureFlagRequest
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

	var body FeatureFlagRequest
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
	"strings"
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
		//
		// Skip multipart/form-data (file uploads): buffering the whole file
		// into memory is pointless for an audit digest, and re-reading it here
		// can leave the handler's ParseMultipartForm with nothing to parse
		// ("No file provided"). Uploads are logged by path/actor, not payload.
		var bodyBytes []byte
		if c.Request.Body != nil &&
			!strings.HasPrefix(c.GetHeader("Content-Type"), "multipart/form-data") {
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
			IPAddress:     resolveClientIP(c),
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

// v3.31.49 -- mirror of services.ResolveClientIP. Inlined here
// (rather than imported) because middleware is a leaf dep that the
// services package itself relies on through the request chain;
// duplicating ten lines avoids the cycle and keeps the audit path
// allocation-free.
func resolveClientIP(c *gin.Context) string {
	ip := c.ClientIP()
	if ip == "::1" || ip == "127.0.0.1" || ip == "0.0.0.0" {
		if hint := strings.TrimSpace(c.GetHeader("X-Public-IP-Hint")); hint != "" {
			if len(hint) > 64 {
				hint = hint[:64]
			}
			return hint
		}
	}
	return ip
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
// CURRENT page only. For a footer that repeats on every page (and page
// numbers), use RunningFooter instead.
func (d *Doc) Footer(text string) {
	d.SetY(-25)
	d.SetFont("Helvetica", "I", 9)
	d.SetTextColor(d.Muted[0], d.Muted[1], d.Muted[2])
	d.CellFormat(0, 5, text, "", 1, "C", false, 0, "")
}

// RunningHeader repeats a header band on EVERY page: the title on the
// left, an optional right-aligned line (document number, date), and a
// thin accent rule beneath. Call it before writing body content — fpdf
// invokes the callback as each page is added, including the first.
//
//	d.RunningHeader("INVOICE", "INV-202605-0001")
func (d *Doc) RunningHeader(title, right string) {
	draw := func() {
		d.SetY(10)
		d.SetFont("Helvetica", "B", 10)
		d.SetTextColor(d.Accent[0], d.Accent[1], d.Accent[2])
		d.CellFormat(0, 5, title, "", 0, "L", false, 0, "")
		if right != "" {
			d.SetFont("Helvetica", "", 9)
			d.SetTextColor(d.Muted[0], d.Muted[1], d.Muted[2])
			d.CellFormat(0, 5, right, "", 0, "R", false, 0, "")
		}
		d.Ln(7)
		d.SetDrawColor(d.Accent[0], d.Accent[1], d.Accent[2])
		d.SetLineWidth(0.4)
		leftM, _, rightM, _ := d.GetMargins()
		w, _ := d.GetPageSize()
		y := d.GetY()
		d.Line(leftM, y, w-rightM, y)
		d.Ln(4)
	}
	d.SetHeaderFunc(draw)
	// New() already added page 1 before any header func existed, so fpdf
	// never invoked the callback for it. Draw it once now; the callback
	// covers every page added from here on.
	if d.PageNo() == 1 {
		draw()
	}
}

// RunningFooter repeats a footer on EVERY page: text on the left and
// "Page N of M" on the right. The page count uses fpdf's page-number
// alias, substituted when the document is finalized.
//
//	d.RunningFooter("Generated 2 Jun 2026 · Acme Ltd")
func (d *Doc) RunningFooter(text string) {
	d.AliasNbPages("")
	d.SetFooterFunc(func() {
		d.SetY(-15)
		d.SetFont("Helvetica", "I", 8)
		d.SetTextColor(d.Muted[0], d.Muted[1], d.Muted[2])
		d.CellFormat(0, 5, text, "", 0, "L", false, 0, "")
		d.CellFormat(0, 5, fmt.Sprintf("Page %d of {nb}", d.PageNo()), "", 0, "R", false, 0, "")
	})
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

// apiPDFRecordGo emits internal/pdf/record.go — the generic, data-driven
// renderer behind every generated resource's GET /:id/pdf endpoint. It knows
// nothing about any specific model: handlers describe the document (title,
// key/value fields, tables, totals) and this lays it out with a repeating
// header, footer and page numbers.
func apiPDFRecordGo() string {
	return `package pdf

import (
	"fmt"
	"reflect"
	"strings"
	"time"
)

// Value formats a single model field for printing: times as "2 Jan 2006",
// booleans as Yes/No, nil as an em dash, everything else via %v. Keeps the
// generated handlers free of per-type formatting noise.
func Value(v any) string {
	switch t := v.(type) {
	case nil:
		return "—"
	case time.Time:
		if t.IsZero() {
			return "—"
		}
		return t.Format("2 Jan 2006")
	case *time.Time:
		if t == nil || t.IsZero() {
			return "—"
		}
		return t.Format("2 Jan 2006")
	case bool:
		if t {
			return "Yes"
		}
		return "No"
	case string:
		if strings.TrimSpace(t) == "" {
			return "—"
		}
		return t
	}
	s := fmt.Sprintf("%v", v)
	if strings.TrimSpace(s) == "" {
		return "—"
	}
	return s
}

// Display renders a related record (a belongs_to association) as a human
// label. It reflects for the first of Name / Title / Subject / Label / Email /
// Number / Code that exists and is non-empty, falling back to the ID. Taking
// "any" means a handler can pass any association without the renderer knowing
// that model's shape.
func Display(v any) string {
	rv := reflect.ValueOf(v)
	for rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			return "—"
		}
		rv = rv.Elem()
	}
	if rv.Kind() != reflect.Struct {
		return Value(v)
	}
	for _, name := range []string{"Name", "Title", "Subject", "Label", "Email", "Number", "Code"} {
		f := rv.FieldByName(name)
		if f.IsValid() && f.Kind() == reflect.String && strings.TrimSpace(f.String()) != "" {
			return f.String()
		}
	}
	if f := rv.FieldByName("ID"); f.IsValid() && f.Kind() == reflect.String {
		if s := f.String(); s != "" {
			return s
		}
	}
	return "—"
}

// Field is one label/value pair in a record's detail grid.
type Field struct {
	Label string
	Value string
}

// Section is a titled table inside a record — an invoice's line items, an
// order's shipments, a booking's guests.
type Section struct {
	Title   string
	Headers []string
	Rows    [][]string
	// Aligns is per-column: "L", "C" or "R". Empty means all left.
	Aligns []string
	// Widths is per-column in mm. Empty means evenly distributed.
	Widths []float64
}

// Record is a whole printable document. Build one in a handler from your
// model and hand it to RenderRecord — nothing here is model-specific, so the
// same shape prints an invoice, a receipt, a work order or a patient chart.
type Record struct {
	// Title is the big word at the top ("INVOICE", "ORDER").
	Title string
	// Subtitle sits under the title — usually the record's identifier.
	Subtitle string
	// Brand is the app/company name shown in the repeating header.
	Brand string
	// Fields render as a two-column detail grid under the header.
	Fields []Field
	// Sections render in order as titled tables.
	Sections []Section
	// Totals renders a right-aligned stack after the sections.
	Totals []TotalLine
	// Notes is free text at the end.
	Notes string
	// FooterNote sits bottom-left on every page, opposite the page number.
	FooterNote string
}

// RenderRecord lays a Record out as PDF bytes: a repeating header band
// (brand + identifier) and footer (note + "Page N of M") on every page, the
// title block, a two-up field grid, each section as a table, then totals and
// notes. Long tables page-break naturally and the header/footer follow.
func RenderRecord(r Record) ([]byte, error) {
	d := New()

	header := r.Brand
	if header == "" {
		header = r.Title
	}
	d.RunningHeader(header, r.Subtitle)
	footer := r.FooterNote
	if footer == "" {
		footer = r.Title
	}
	d.RunningFooter(footer)

	d.Header(r.Title, r.Subtitle)

	// Detail grid, two pairs per row so a record with many columns stays
	// compact instead of running one-per-line down the page.
	for i := 0; i < len(r.Fields); i += 2 {
		if i+1 < len(r.Fields) {
			d.TwoColumnKV(
				r.Fields[i].Label, r.Fields[i].Value,
				r.Fields[i+1].Label, r.Fields[i+1].Value,
			)
			continue
		}
		d.KV(r.Fields[i].Label, r.Fields[i].Value)
	}

	for _, s := range r.Sections {
		if len(s.Rows) == 0 {
			continue
		}
		if s.Title != "" {
			d.Ln(3)
			d.SetFont("Helvetica", "B", 11)
			d.SetTextColor(0, 0, 0)
			d.CellFormat(0, 6, s.Title, "", 1, "L", false, 0, "")
			d.Ln(1)
		}
		aligns := s.Aligns
		if len(aligns) == 0 {
			aligns = make([]string, len(s.Headers))
			for i := range aligns {
				aligns[i] = "L"
			}
		}
		widths := s.Widths
		if len(widths) == 0 {
			widths = evenWidths(d, len(s.Headers))
		}
		d.Table(s.Headers, s.Rows, widths, aligns)
	}

	if len(r.Totals) > 0 {
		d.Totals(r.Totals)
	}
	if strings.TrimSpace(r.Notes) != "" {
		d.Notes(r.Notes)
	}

	return d.Bytes()
}

// evenWidths splits the printable width evenly across n columns.
func evenWidths(d *Doc, n int) []float64 {
	if n <= 0 {
		return nil
	}
	left, _, right, _ := d.GetMargins()
	pageW, _ := d.GetPageSize()
	each := (pageW - left - right) / float64(n)
	out := make([]float64, n)
	for i := range out {
		out[i] = each
	}
	return out
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
	"strings"
	"time"
	"unicode"

	"github.com/MUKE-coder/gorm-studio/studio"
	"github.com/MUKE-coder/pulse/pulse"
	sentinel "github.com/MUKE-coder/sentinel/v2"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"` + "{{MODULE}}" + `/internal/ai"
	"` + "{{MODULE}}" + `/internal/cache"
	"` + "{{MODULE}}" + `/internal/config"
	"` + "{{MODULE}}" + `/internal/database"
	"` + "{{MODULE}}" + `/internal/events"
	"` + "{{MODULE}}" + `/internal/handlers"
	"` + "{{MODULE}}" + `/internal/settings"
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

// splitOrigins parses the cors.origins setting.
//
// Newlines or commas, because the admin renders a textarea and people type
// both. Blank lines and stray whitespace are dropped rather than becoming an
// origin nothing can ever match.
func splitOrigins(raw string) []string {
	fields := strings.FieldsFunc(raw, func(r rune) bool {
		// No escapes here on purpose. unicode.IsSpace covers newline, carriage
		// return, tab and space in one call, rather than four rune literals a
		// shell heredoc can eat on the way in.
		return r == ',' || unicode.IsSpace(r)
	})
	out := make([]string, 0, len(fields))
	for _, f := range fields {
		if f != "" {
			out = append(out, f)
		}
	}
	return out
}

// eventBusStatus summarises the domain event bus for /api/health.
//
// Nil-safe so a project whose routes.go predates events.Init still answers
// the health check rather than panicking on it.
func eventBusStatus() interface{} {
	bus := events.Default()
	if bus == nil {
		return map[string]interface{}{"ok": false, "configured": false}
	}
	s := bus.Stats()
	return map[string]interface{}{
		"ok":          s.Dropped == 0,
		"configured":  true,
		"subscribers": s.Subscribers,
		"queued":      s.Queued,
		"capacity":    s.Capacity,
		"dropped":     s.Dropped,
	}
}

// APIVersion is the version segment every /api route is served under, so the
// public surface is /api/v1/... rather than /api/....
//
// Why a prefix at all: once anything outside this repo calls your API — a
// mobile build you can't force-update, a partner integration, a customer's
// script — you can no longer change a response shape without breaking them.
// A version in the path gives you somewhere to put the new shape. When that
// day comes, add a v2 group next to v1 and leave v1 answering the old way
// until consumers have moved; delete it when your logs say nobody's left.
//
// Unversioned /api/... requests are rewritten to this version (see
// mountLegacyAPIAlias), so existing clients keep working after an upgrade.
// That alias is a courtesy for the transition, not a second API: it always
// points at whatever APIVersion currently is, so a client that never adopts
// the prefix will eventually be dragged onto a version it wasn't written for.
const APIVersion = "v1"

// wafExcludedRoutes lists the paths Sentinel's WAF steps aside for, under the
// live API prefix.
//
// Uploads are here because the WAF rejects any body over its inspection cap
// before the route runs, and a photograph is larger than that cap by design.
// The richtext resources are here because the XSS heuristics flag ordinary
// markup: a blog body is <p> and <strong> and <img> by definition.
//
// Exclusion is from body inspection only. These routes still pass through
// auth, RBAC, binding validation and rate limiting.
func wafExcludedRoutes() []string {
	prefix := "/api/" + APIVersion
	paths := []string{
		"/blogs", "/blogs/*",
		"/posts", "/posts/*",
		"/articles", "/articles/*",
		"/uploads", "/uploads/*",
		// Public form-share submissions. Auth is the share's bcrypt password
		// (optional) and the token itself; Sentinel rate-limits the path. The
		// subtree match also covers .../submit.
		"/public/forms/*",
	}
	out := make([]string, 0, len(paths))
	for _, p := range paths {
		out = append(out, prefix+p)
	}
	return out
}

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
	// Origins come from the cors.origins setting when it has a value, and from
	// CORS_ORIGINS otherwise. Resolved per request, so adding a domain in the
	// admin takes effect immediately rather than at the next deploy.
	r.Use(middleware.CORSDynamic(func() []string {
		if stored := settings.String(context.Background(), "cors.origins"); strings.TrimSpace(stored) != "" {
			return splitOrigins(stored)
		}
		return cfg.CORSOrigins
	}))
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

		// Sentinel persists its security data (threat log, blocked IPs,
		// audit trail) through its own storage adapter, NOT the *gorm.DB we
		// pass in. Left unset it silently falls back to a local sentinel.db
		// SQLite file — which is ephemeral inside a container, so every
		// redeploy would drop the threat log and the blocked-IP list. Point
		// it at the same database the app uses when that's Postgres.
		sentinelStorage := sentinel.StorageConfig{Driver: sentinel.SQLite, DSN: "sentinel.db"}
		if !strings.HasPrefix(cfg.DatabaseURL, "sqlite:") {
			sentinelStorage = sentinel.StorageConfig{Driver: sentinel.Postgres, DSN: cfg.DatabaseURL}
		}

		// Sentinel v2 — use MountE so we can recover gracefully on
		// misconfiguration in dev instead of log.Fatalf-ing the host.
		// Mount runs sentinel.ValidateConfig and logs any dead config.
		if err := sentinel.MountE(r, db, sentinel.Config{
			Storage: sentinelStorage,
			Dashboard: sentinel.DashboardConfig{
				Username:               cfg.SentinelUsername,
				Password:               cfg.SentinelPassword,
				SecretKey:              cfg.SentinelSecretKey,
				// Sentinel refuses default credentials in gin.ReleaseMode;
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
				// 1 MB cap covers richtext admin payloads — Tiptap blog
				// bodies with embedded inline images comfortably exceed
				// the prior 64 KB ceiling. Bump higher if your content
				// embeds large base64 images.
				MaxBodyBytes:          1 * 1024 * 1024,
				RejectOversizedBody:   true,
				// Authenticated admin write endpoints handle their own
				// HTML/richtext payloads via Tiptap. The WAF's XSS detection
				// otherwise flags every <p>/<strong>/<img> tag in a blog
				// body as a payload. These routes still pass through auth
				// + RBAC + binding validation; WAF is just stepped aside
				// for their bodies.
				//
				// IMPORTANT: the WAF matches these against the real request
				// path (c.Request.URL.Path), NOT gin's route template. Gin
				// params like "/api/blogs/:id" therefore match only the
				// literal string ":id" and never "/api/blogs/123" — they
				// were silent dead config. Use "/*" (a subtree match) so the
				// id/token routes are actually excluded.
				//
				// They are built from APIVersion for the same reason. Every
				// entry here was once written as a literal "/api/blogs", while
				// the router mounts "/api/" + APIVersion, so not one of them
				// ever matched. Two things were broken by that and neither
				// announced itself: an upload over MaxBodyBytes was rejected
				// with 413 before the handler saw it, and richtext bodies were
				// never actually stepped aside, so in production (ModeBlock) a
				// blog post containing markup could be refused as an XSS
				// payload. Deriving the prefix means the next version bump
				// cannot quietly disable all of it again.
				ExcludeRoutes: wafExcludedRoutes(),
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
			log.Println("Sentinel v2.2.0 mounted at /sentinel")
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
	//
	// The OpenAPI reference. Its 141 route overrides live in apidocs.go, where
	// they are 500 lines of description rather than 500 lines in the middle of
	// the file that wires your application together.
	registerAPIDocs(r, db, cfg)

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
			// CRITICAL: Pulse's error middleware captures a request-body snippet
			// (MaxBodySize, default 4096) for error context, but restores ONLY
			// that snippet to the request — it discards everything past 4096
			// bytes. That truncates EVERY request carrying a Content-Length
			// (mobile / native / curl clients; browsers dodge it by sending
			// chunked), silently breaking file uploads and any large JSON POST.
			// Disable body capture so the full body always reaches the handler.
			pulse.WithRequestBodyCaptureDisabled(),
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
		Mailer:      svc.Mailer,
	}
	apiKeyHandler := &handlers.APIKeyHandler{DB: db}

	userHandler := &handlers.UserHandler{
		DB:          db,
		AuthService: authService,
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
	totpHandler := &handlers.TOTPHandler{
		DB:          db,
		AuthService: authService,
		Issuer:      cfg.TOTPIssuer,
	}
	activityHandler := handlers.NewActivityHandler(db)
	webhookHandler := handlers.NewWebhookHandler(db)
	webhooks.Setup(db)
	realtimeHub := realtime.NewHub()

	// The domain event bus. Created before any handler so an emit during
	// startup has somewhere to go, and wired to the audit log, realtime and
	// (when the plugin is installed) outbound webhooks.
	events.Init(4)
	services.RegisterEventSubscribers(db, realtimeHub, nil)

	// Settings: declare, then open the store. Declaring after Init would mean
	// a setting the first cache load never saw.
	settings.RegisterDefaults()
	settings.Init(db)
	settingsHandler := &handlers.SettingsHandler{DB: db}
	flagsEngine := flags.New(db, realtimeHub)
	featureFlagHandler := handlers.NewFeatureFlagHandler(db, flagsEngine)
	realtimeHandler := handlers.NewRealtimeHandler(realtimeHub, authService)
	_ = realtimeHub // available to handlers/services that want to push events

	// In-app Security + Observability dashboards — read from Sentinel/Pulse APIs
	// over loopback. notificationHandler powers the admin bell.
	notificationHandler := &handlers.NotificationHandler{DB: db}
	securityHandler := &handlers.SecurityHandler{Bridge: svc.SecObs}
	observabilityHandler := &handlers.ObservabilityHandler{Bridge: svc.SecObs}

	// v3.30 — semantic activity log + ticket system. Mailer is optional;
	// when nil the ticket handler skips email-out and only writes the row
	// + admin notifications.
	userActivityHandler := &handlers.UserActivityHandler{DB: db}
	ocsfHandler := handlers.NewOCSFHandler(db, cfg.AppName)
	accessReviewHandler := handlers.NewAccessReviewHandler(db)
	gdprHandler := handlers.NewGDPRHandler(db)
	ticketHandler := &handlers.TicketHandler{DB: db, Mail: svc.Mailer}
	// v3.31.20 — public form sharing (Phase 2)
	formShareHandler := &handlers.FormShareHandler{DB: db}
	// v3.31.40 — per-user dashboard customisation
	dashboardLayoutHandler := &handlers.DashboardLayoutHandler{DB: db}
	// v3.31.44 — per-resource dashboard stats (Total + sparkline + Latest N)
	resourceStatsHandler := &handlers.ResourceStatsHandler{DB: db}
	// v3.31.47 — Preset Chart builder
	chartHandler := &handlers.ChartHandler{DB: db}

	// Sync registry — list every model that should be syncable from
	// offline-first desktop clients. The resource generator injects
	// new resources at the marker below.
	syncRegistry := sync.NewRegistry()
	syncRegistry.Register("users", &models.User{})
	syncRegistry.Register("uploads", &models.Upload{})
	syncRegistry.Register("blogs", &models.Blog{})
	// grit:sync
	syncHandler := handlers.NewSyncHandler(db, syncRegistry)
	// v3.31.68 — shared background CSV import status endpoint
	importJobHandler := &handlers.ImportJobHandler{DB: db}
	// v3.31.77 — full-database backups (weekly cron + manual + download)
	backupHandler := &handlers.BackupHandler{DB: db, Storage: svc.Storage}
	roleHandler := handlers.NewRoleHandler(db)
	sessionHandler := handlers.NewSessionHandler(db)

	// Enterprise SSO. Providers are built once here (each one performs OIDC
	// discovery against the customer's IdP) and rebuilt whenever an admin saves
	// a connection, so adding a customer never needs a restart. A connection
	// whose discovery fails is logged and skipped — one broken IdP must not
	// stop everyone else signing in.
	ssoRegistry := services.NewSSORegistry(cfg.AppURL)
	for _, err := range ssoRegistry.Reload(db) {
		log.Printf("sso: %v", err)
	}
	samlRegistry := services.NewSAMLRegistry(cfg.AppURL)
	for _, err := range samlRegistry.Reload(db) {
		log.Printf("saml: %v", err)
	}
	ssoHandler := handlers.NewSSOHandler(db, authService, cfg, ssoRegistry, samlRegistry)
	// grit:handlers

	// Health check
	// /api/health probes every infrastructure dependency the dashboard's
	// System Health page wants to render. Each probe is bounded by a 500ms
	// timeout so a hung dependency doesn't pile up health requests; failing
	// probes mark themselves down and the overall status downgrades to
	// "degraded" rather than failing the endpoint.
	r.GET("/api/health", func(c *gin.Context) {
		type compStatus struct {
			OK         bool   ` + "`" + `json:"ok"` + "`" + `
			LatencyMS  int64  ` + "`" + `json:"latency_ms,omitempty"` + "`" + `
			Tables     int    ` + "`" + `json:"tables,omitempty"` + "`" + `
			QueueKeys  int    ` + "`" + `json:"queue_keys,omitempty"` + "`" + `
			Configured bool   ` + "`" + `json:"configured,omitempty"` + "`" + `
			Error      string ` + "`" + `json:"error,omitempty"` + "`" + `
		}

		// Database ping + table count. We probe with a 500ms deadline so a
		// blocked write loop can't hang the health check.
		dbStatus := compStatus{OK: true}
		dbStart := time.Now()
		if sqlDB, err := db.DB(); err == nil {
			ctx, cancel := context.WithTimeout(c.Request.Context(), 500*time.Millisecond)
			defer cancel()
			if err := sqlDB.PingContext(ctx); err != nil {
				dbStatus.OK = false
				dbStatus.Error = err.Error()
			}
		}
		dbStatus.LatencyMS = time.Since(dbStart).Milliseconds()
		if dbStatus.OK {
			// Best-effort table count. Dialect-aware, and 0 rather than an
			// error when the database cannot be asked: a missing tooltip
			// figure is not a health problem.
			dbStatus.Tables = database.TableCount(db)
		}

		// Redis ping. Reuse the same cache client the rest of the app uses
		// rather than opening a new connection — that way "Redis healthy"
		// on the dashboard means the same Redis the cache + jobs use.
		redisStatus := compStatus{}
		if svc.Cache != nil {
			redisStart := time.Now()
			ctx, cancel := context.WithTimeout(c.Request.Context(), 500*time.Millisecond)
			defer cancel()
			if err := svc.Cache.Client().Ping(ctx).Err(); err != nil {
				redisStatus.OK = false
				redisStatus.Error = err.Error()
			} else {
				redisStatus.OK = true
			}
			redisStatus.LatencyMS = time.Since(redisStart).Milliseconds()
		}

		// Background-jobs queue — count active asynq keys as a liveness
		// signal. If asynq isn't wired (Jobs == nil), report unconfigured
		// rather than "down" so the dashboard distinguishes the cases.
		jobsStatus := compStatus{}
		if svc.Jobs != nil && svc.Cache != nil {
			ctx, cancel := context.WithTimeout(c.Request.Context(), 500*time.Millisecond)
			defer cancel()
			n, err := svc.Cache.Client().Eval(ctx,
				"local total = 0\nfor _, k in ipairs(redis.call('keys', 'asynq:*')) do total = total + 1 end\nreturn total",
				[]string{}).Int()
			if err == nil {
				jobsStatus.OK = true
				jobsStatus.QueueKeys = n
			} else {
				// Fall back to a simple ping so a "no keys yet" install still
				// reports OK rather than down.
				if perr := svc.Cache.Client().Ping(ctx).Err(); perr == nil {
					jobsStatus.OK = true
				}
			}
		}

		// Email is "configured" when Resend key is set + non-default. The
		// dashboard treats unconfigured as "—" not "down".
		mailStatus := compStatus{
			Configured: cfg.ResendAPIKey != "" && cfg.ResendAPIKey != "re_your_api_key",
			OK:         cfg.ResendAPIKey != "" && cfg.ResendAPIKey != "re_your_api_key",
		}

		// Overall status — ok if every wired-up component is up. Components
		// that aren't configured (e.g. Redis off in a single-binary dev
		// run) don't drag the overall status down.
		overall := "ok"
		if !dbStatus.OK || (svc.Cache != nil && !redisStatus.OK) {
			overall = "degraded"
		}

		c.JSON(http.StatusOK, gin.H{
			"status":   overall,
			"version":  "0.1.0",
			"database": dbStatus,
			"redis":    redisStatus,
			"api":      compStatus{OK: true},
			"jobs":     jobsStatus,
			"email":    mailStatus,
			// The event bus reports itself. Dropped rising is the only signal
			// from outside that a subscriber is too slow or the queue too
			// small, and "did my webhook fire" deserves a better answer than
			// reading logs.
			"events": eventBusStatus(),
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


	// ── API version ──────────────────────────────────────────────────────
	// Every /api route hangs off this group, so the whole surface is served
	// under /api/v1. When a breaking change is unavoidable, add a v2 group
	// beside it and keep v1 serving the old shape until consumers migrate —
	// that's the entire point of the prefix.
	//
	// Unversioned /api/... requests are rewritten to the current version by
	// the fallback at the bottom of this file, so older clients (and the
	// generated frontends) keep working untouched.
	v1 := r.Group("/api/" + APIVersion)

	// Public blog routes (no auth required)
	blogs := v1.Group("/blogs")
	{
		blogs.GET("", blogHandler.ListPublished)
		blogs.GET("/:slug", blogHandler.GetBySlug)
	}

	// Public API surface, for clients with no logged-in user: a storefront, a
	// mobile app, a public directory.
	//
	// Guarded by an API key rather than open. That is not secrecy, because a
	// publishable key ships inside your app where anyone can read it. It buys
	// identification, a rate-limit bucket per key, per-endpoint and per-origin
	// narrowing, and the ability to turn one client off without a deploy.
	//
	// Resources land here through: grit generate resource <Name> --public
	publicAPI := v1.Group("/public")
	publicAPI.Use(middleware.RequireAPIKey(db, svc.Cache))
	// Response caching, and only here.
	//
	// The cache key is the URL, nothing else. On a public endpoint that is
	// exactly right: every caller gets the same answer, so one cached copy
	// serves all of them and a catalogue page stops hitting Postgres on every
	// visit. On a protected endpoint the same key would serve one user's data
	// to another, which is why this middleware is mounted on this group and
	// nowhere else.
	//
	// The TTL is read once at boot rather than per request. A cache lifetime is
	// not something anybody changes at 9pm, and re-reading it on the hot path
	// of a cached response would cost more than it saves.
	if svc.Cache != nil {
		ttl := time.Duration(settings.Int(context.Background(), "cache.public_ttl_seconds")) * time.Second
		if ttl <= 0 {
			ttl = 60 * time.Second
		}
		publicAPI.Use(middleware.CacheResponse(svc.Cache, ttl))
		log.Printf("Public endpoints cached for %s", ttl)
	}
	{
		// grit:routes:public
	}

	// Public auth routes
	auth := v1.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
		auth.POST("/refresh", authHandler.Refresh)
		auth.POST("/forgot-password", authHandler.ForgotPassword)
		auth.POST("/reset-password", authHandler.ResetPassword)
		auth.POST("/verify-email", authHandler.VerifyEmail)
	}

	// OAuth2 social login
	oauth := auth.Group("/oauth")
	{
		oauth.GET("/:provider", authHandler.OAuthBegin)
		oauth.GET("/:provider/callback", authHandler.OAuthCallback)
	}

	// Enterprise SSO (OIDC). Public by design — these ARE the login flow.
	// Discover tells the login form whether an address belongs to a connection;
	// the other two are the redirect out to the IdP and the return trip.
	//
	// Like the OAuth callbacks above, /callback is registered in the customer's
	// IdP console, so its unversioned path must keep working — see the note on
	// APIVersion.
	sso := auth.Group("/sso")
	{
		sso.POST("/discover", ssoHandler.Discover)
		sso.GET("/:slug", ssoHandler.Begin)
		sso.GET("/:slug/callback", ssoHandler.Callback)
	}

	// SAML 2.0. /metadata is what the customer uploads to their IdP and /acs is
	// where that IdP POSTs the signed assertion — both get registered on their
	// side, so like the OAuth callbacks these unversioned paths must keep
	// working across API version bumps.
	samlGroup := auth.Group("/saml")
	{
		samlGroup.GET("/:slug/metadata", ssoHandler.SAMLMetadata)
		samlGroup.GET("/:slug", ssoHandler.SAMLBegin)
		samlGroup.POST("/:slug/acs", ssoHandler.SAMLACS)
	}

	// TOTP verification (public — uses pending tokens, not JWT)
	auth.POST("/totp/verify", totpHandler.Verify)
	auth.POST("/totp/backup-codes/verify", totpHandler.VerifyBackupCode)

	// Protected routes
	protected := v1.Group("")
	// Accepts an API key OR the usual JWT. With no key header present this
	// delegates straight to middleware.Auth, so browser sessions behave
	// exactly as before; with one, it sets the same context values so every
	// downstream handler and RequireRole check works unchanged.
	protected.Use(middleware.APIKeyOrAuth(db, middleware.Auth(db, authService)))
	// Activity logger writes one row per successful authenticated mutation.
	// Records who/what/when/where for audit. Read-only — see admin/activity.
	protected.Use(middleware.ActivityLogger(db))
	{
		protected.GET("/auth/me", authHandler.Me)
		// The caller's own permissions, for the frontend can() helper and nav
		// gating. Any authenticated user may read their own — it tells them
		// nothing they can't already discover by clicking.
		protected.GET("/auth/permissions", roleHandler.MyPermissions)

		// Which optional modules are enabled. The admin reads this to hide nav
		// entries for modules that are switched off — a dead link to a route
		// that no longer exists is worse than no link.
		protected.GET("/system/modules", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"data": cfg.Modules.Map()})
		})
		protected.POST("/auth/logout", authHandler.Logout)

		// Active sessions — see every signed-in device and revoke one or all.
		protected.GET("/auth/sessions", sessionHandler.List)
		protected.DELETE("/auth/sessions/:id", sessionHandler.Revoke)
		protected.POST("/auth/sessions/revoke-all", sessionHandler.RevokeAll)

		// Two-Factor Authentication (TOTP)
		protected.POST("/auth/totp/setup", totpHandler.Setup)
		protected.POST("/auth/totp/enable", totpHandler.Enable)
		protected.POST("/auth/totp/disable", totpHandler.Disable)
		protected.GET("/auth/totp/status", totpHandler.Status)
		protected.POST("/auth/totp/backup-codes", totpHandler.RegenerateBackupCodes)
		protected.DELETE("/auth/totp/trusted-devices", totpHandler.RevokeTrustedDevices)
		protected.POST("/auth/verify-email/send", authHandler.SendVerificationEmail)
		protected.GET("/api-keys", apiKeyHandler.List)
		protected.POST("/api-keys", apiKeyHandler.Create)
		protected.DELETE("/api-keys/:id", apiKeyHandler.Revoke)
		protected.GET("/auth/totp/trusted-devices", totpHandler.ListTrustedDevices)
		protected.DELETE("/auth/totp/trusted-devices/:id", totpHandler.RevokeTrustedDevice)

		// User routes (authenticated)
		protected.GET("/users/:id", userHandler.GetByID)

		// GDPR right-to-access: a user may export their own data; an admin, anyone's.
		protected.GET("/users/:id/gdpr-export", gdprHandler.Export)

		// File uploads
		protected.POST("/uploads", uploadHandler.Create)
		protected.POST("/uploads/presign", uploadHandler.Presign)
		protected.POST("/uploads/complete", uploadHandler.CompleteUpload)
		protected.GET("/uploads", uploadHandler.List)
		protected.GET("/uploads/stats", uploadHandler.Stats)
		protected.GET("/uploads/:id", uploadHandler.GetByID)
		protected.DELETE("/uploads/:id", uploadHandler.Delete)

		// Offline-first sync — desktop clients call these to flush their
		// local outbox and pull server-side updates.
		protected.POST("/sync/push", syncHandler.Push)
		protected.GET("/sync/pull", syncHandler.Pull)
		protected.GET("/sync/policy", syncHandler.Policy)

		// Reading settings is open to any authenticated caller, because a
		// screen needs app.name to render its header. Writing is admin-only
		// and mounted with the other admin routes below.
		protected.GET("/settings", settingsHandler.List)

		// AI — only mounted when the module is enabled, so an app that
		// doesn't use it exposes no AI surface at all (MODULE_AI=false).
		if cfg.Modules.AI {
			protected.POST("/ai/complete", aiHandler.Complete)
			protected.POST("/ai/chat", aiHandler.Chat)
			protected.POST("/ai/stream", aiHandler.Stream)
		}


		// In-app notification bell — every authenticated user. Pulls
		// from a single Notification table that the SecObs poller
		// writes into when Sentinel/Pulse fires a high-severity event.
		protected.GET("/notifications", notificationHandler.List)
		protected.POST("/notifications/:id/read", notificationHandler.MarkRead)
		protected.POST("/notifications/read-all", notificationHandler.MarkAllRead)

		// v3.31.40 — per-user dashboard layout customisation.
		protected.GET("/dashboard-layout", dashboardLayoutHandler.Get)
		protected.PUT("/dashboard-layout", dashboardLayoutHandler.Put)

		// v3.30 — tickets. Any authenticated user can open + reply; the
		// handler scopes List/Get visibility to the caller unless they're
		// ADMIN/EDITOR (then they see the full queue).
		protected.POST("/tickets", ticketHandler.Create)
		protected.GET("/tickets", ticketHandler.List)
		protected.GET("/tickets/:id", ticketHandler.Get)
		protected.POST("/tickets/:id/reply", ticketHandler.Reply)
		protected.PATCH("/tickets/:id/close", ticketHandler.Close)
		protected.PATCH("/tickets/:id/reopen", ticketHandler.Reopen)
		protected.PATCH("/tickets/:id/assign", ticketHandler.Assign) // admin-gated inside the handler

		// v3.31.68 — poll a background CSV import's progress/result.
		protected.GET("/imports/:id", importJobHandler.GetByID)

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
	admin := v1.Group("")
	admin.Use(middleware.APIKeyOrAuth(db, middleware.Auth(db, authService)))
	admin.Use(middleware.RequireRole("ADMIN"))
	{
		admin.GET("/users", userHandler.List)
		admin.POST("/users", userHandler.Create)
		admin.PUT("/users/:id", userHandler.Update)
		admin.DELETE("/users/:id", userHandler.Delete)

		// Activity audit log + tamper-evident chain verification
		admin.GET("/admin/activity", activityHandler.List)
		admin.GET("/admin/activity/integrity", activityHandler.VerifyIntegrity)

		// v3.30 — semantic user activity dashboard (action + IP + severity).
		// Separate from /admin/activity above which is the HTTP audit log.
		admin.GET("/user-activity", userActivityHandler.List)
		admin.GET("/user-activity/stats", userActivityHandler.Stats)

		// OCSF audit export — the semantic activity log in the vendor-neutral
		// shape SIEMs ingest. Cursor-paginated NDJSON; poll to resume.
		admin.GET("/audit/ocsf", ocsfHandler.Export)

		// Access reviews (recertification) — snapshot every grant, certify or
		// revoke each, sign off. Admin-only; revocations hit the audit trail.
		admin.GET("/access-reviews", accessReviewHandler.List)
		admin.POST("/access-reviews", accessReviewHandler.Open)
		admin.GET("/access-reviews/:id", accessReviewHandler.Get)
		admin.POST("/access-reviews/:id/items/:itemId/decision", accessReviewHandler.Decide)
		admin.POST("/access-reviews/:id/complete", accessReviewHandler.Complete)

		// GDPR right-to-erasure: anonymize a user + hard-delete their PII, recorded
		// in a tamper-evident deletion journal. Admin-only; the journal is verifiable.
		admin.POST("/users/:id/gdpr-erase", gdprHandler.Erase)
		admin.GET("/gdpr/journal", gdprHandler.Journal)

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
		admin.GET("/admin/blogs/:id", blogHandler.GetByID)
		admin.POST("/admin/blogs", blogHandler.Create)
		admin.PUT("/admin/blogs/:id", blogHandler.Update)
		admin.DELETE("/admin/blogs/:id", blogHandler.Delete)


		// In-app Security dashboard — aggregates Sentinel APIs into one
		// envelope so the React page does a single round-trip. Operators
		// who want to dig deeper open /sentinel/ui directly.
		admin.GET("/admin/security/summary", securityHandler.Summary)
		// In-app Observability dashboard — same pattern against Pulse.
		// Operators who want a flame graph or the full SLO timeline open
		// /pulse/ui directly.
		admin.GET("/admin/observability/summary", observabilityHandler.Summary)

		// v3.31.20 — public form sharing admin
		// SSO connections — admin only. Client secrets are write-only: they go
		// in on create/update and are never returned, so a compromised admin
		// session can't read a customer's IdP credentials back out.
		admin.GET("/sso/connections", middleware.RequireRole("ADMIN"), ssoHandler.List)
		admin.POST("/sso/connections", middleware.RequireRole("ADMIN"), ssoHandler.Create)
		admin.PUT("/sso/connections/:id", middleware.RequireRole("ADMIN"), ssoHandler.Update)
		admin.DELETE("/sso/connections/:id", middleware.RequireRole("ADMIN"), ssoHandler.Delete)
		admin.GET("/sso/connections/:id/test", middleware.RequireRole("ADMIN"), ssoHandler.Test)

		admin.GET("/admin/form-shares", formShareHandler.List)
		admin.POST("/admin/form-shares", formShareHandler.Create)
		admin.PATCH("/admin/form-shares/:id", formShareHandler.Update)
		admin.DELETE("/admin/form-shares/:id", formShareHandler.Delete)
		// v3.31.50 — dropdown source + field preview for the New
		// Share / Edit Share modal. Both read-only.
		admin.GET("/admin/form-shares/resources", formShareHandler.Resources)
		admin.GET("/admin/form-shares/resources/:resource/fields", formShareHandler.FieldsPreview)
		// v3.31.25 — audit log of public submissions
		admin.GET("/admin/form-submissions", formShareHandler.ListSubmissions)

		// v3.31.44 — per-resource dashboard stats: Total + 30-day
		// sparkline + Latest N. Dispatched server-side; only resources
		// registered in services/resource_stats_dispatch.go are reachable.
		admin.GET("/admin/dashboard/resource-stats/:resource", resourceStatsHandler.Get)

		// v3.31.47 — Preset Chart builder. Same dispatch boundary;
		// only resources registered in chart_dispatch.go reachable.
		admin.GET("/admin/dashboard/chart/:resource", chartHandler.Get)

		// v3.31.77 — full-database backups. Weekly cron writes them; an
		// operator can also take one on demand (rate-limited to 1/24h) and
		// download it via a short-lived pre-signed URL straight from storage.
		admin.GET("/backups", backupHandler.List)
		admin.POST("/backups/generate", backupHandler.Generate)
		admin.GET("/backups/:id/download", backupHandler.Download)
		// Separate path (not /backups/settings) so it doesn't collide with the
		// /backups/:id wildcard segment in Gin's router.
		admin.GET("/backup-settings", backupHandler.GetSettings)
		admin.PUT("/backup-settings", backupHandler.UpdateSettings)

		// Roles & permissions. Guarded by permission as well as the group's
		// ADMIN role, so a custom role can be given role-management rights
		// without being made a full admin.
		admin.GET("/permissions", roleHandler.Catalog)
		admin.GET("/roles", middleware.RequireRole("ADMIN", "perm:roles.view"), roleHandler.List)
		admin.POST("/roles", middleware.RequireRole("ADMIN", "perm:roles.create"), roleHandler.Create)
		admin.GET("/roles/:id", middleware.RequireRole("ADMIN", "perm:roles.view"), roleHandler.Get)
		admin.PUT("/roles/:id", middleware.RequireRole("ADMIN", "perm:roles.edit"), roleHandler.Update)
		admin.DELETE("/roles/:id", middleware.RequireRole("ADMIN", "perm:roles.delete"), roleHandler.Delete)
		admin.PUT("/users/:id/roles", middleware.RequireRole("ADMIN", "perm:users.edit"), roleHandler.AssignUserRoles)
		admin.POST("/users/:id/unlock", middleware.RequireRole("ADMIN", "perm:users.edit"), userHandler.Unlock)

		// Writing settings. Per-setting permissions are checked inside the
		// handler, because which permission applies depends on which setting
		// is being changed and a route can only know one.
		admin.PUT("/settings", settingsHandler.Update)
		admin.DELETE("/settings/:key", settingsHandler.Reset)

		// grit:routes:admin
	}

	// Public form-sharing endpoints. NO auth, NO CSRF — Sentinel rate
	// limits each token aggressively. The dispatch service is the
	// security boundary (whitelists which resources are reachable).
	publicForms := v1.Group("/public/forms")
	{
		publicForms.GET("/:token", formShareHandler.PublicGet)
		publicForms.POST("/:token/submit", formShareHandler.PublicSubmit)
	}

	// Custom role-restricted routes
	// grit:routes:custom

	mountLegacyAPIAlias(r)

	return r
}

// mountLegacyAPIAlias keeps unversioned /api/... paths working by re-dispatching
// them to /api/<APIVersion>/... .
//
// It runs as the 404 fallback rather than as middleware because Gin resolves the
// route before middleware executes — by the time a handler could rewrite the
// path, the routing decision is already made. Landing here means no route
// matched, so the only cost is on requests that were going to 404 anyway.
//
// /api/ws is deliberately excluded: a WebSocket upgrade re-dispatched through
// HandleContext does not survive reliably, and a transport endpoint isn't part
// of the REST surface being versioned.
func mountLegacyAPIAlias(r *gin.Engine) {
	versioned := "/api/" + APIVersion + "/"

	r.NoRoute(func(c *gin.Context) {
		p := c.Request.URL.Path

		if strings.HasPrefix(p, "/api/") &&
			!strings.HasPrefix(p, versioned) &&
			p != "/api/ws" {
			c.Request.URL.Path = "/api/" + APIVersion + strings.TrimPrefix(p, "/api")
			// Tell the caller they're on a deprecated path. Harmless to
			// ignore, but it shows up in their logs before v2 forces the issue.
			c.Header("Deprecation", "true")
			c.Header("Link", "</api/"+APIVersion+">; rel=\"successor-version\"")
			r.HandleContext(c)
			return
		}

		c.JSON(http.StatusNotFound, gin.H{
			"error": gin.H{
				"code":    "NOT_FOUND",
				"message": "no route matches " + c.Request.Method + " " + p,
			},
		})
	})
}
`
}
