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
		filepath.Join(apiRoot, "go.mod"):                                            apiGoMod(opts),
		filepath.Join(apiRoot, ".gitignore"):                                        apiGitignore(),
		filepath.Join(apiRoot, "cmd", "server", "main.go"):                          apiMainGo(opts),
		filepath.Join(apiRoot, "internal", "config", "config.go"):                   apiConfigGo(),
		filepath.Join(apiRoot, "internal", "database", "database.go"):               apiDatabaseGo(),
		filepath.Join(apiRoot, "internal", "database", "replicas.go"):               apiReplicasGo(),
		filepath.Join(apiRoot, "internal", "cache", "remember.go"):                  apiCacheRememberGo(),
		filepath.Join(apiRoot, "internal", "middleware", "latency.go"):              apiLatencyMiddlewareGo(),
		filepath.Join(apiRoot, "internal", "handlers", "scale.go"):                  apiScaleHandlerGo(),
		filepath.Join(apiRoot, "internal", "database", "for_request.go"):            apiForRequestGo(),
		filepath.Join(apiRoot, "internal", "middleware", "read_your_writes.go"):     apiReadYourWritesGo(),
		filepath.Join(apiRoot, "internal", "database", "dialect.go"):                apiDialectGo(),
		filepath.Join(apiRoot, "internal", "models", "user.go"):                     apiUserModelGo(),
		filepath.Join(apiRoot, "internal", "models", "upload.go"):                   apiUploadModelGo(),
		filepath.Join(apiRoot, "internal", "services", "auth.go"):                   apiAuthServiceGo(),
		filepath.Join(apiRoot, "internal", "models", "session.go"):                  apiSessionModelGo(),
		filepath.Join(apiRoot, "internal", "services", "session.go"):                apiSessionServiceGo(),
		filepath.Join(apiRoot, "internal", "handlers", "session.go"):                apiSessionHandlerGo(),
		filepath.Join(apiRoot, "internal", "models", "sso.go"):                      apiSSOModelGo(),
		filepath.Join(apiRoot, "internal", "services", "sso.go"):                    apiSSOServiceGo(),
		filepath.Join(apiRoot, "internal", "services", "sso_connections_test.go"):   apiSSOServiceTestGo(),
		filepath.Join(apiRoot, "internal", "handlers", "sso.go"):                    apiSSOHandlerGo(),
		filepath.Join(apiRoot, "internal", "models", "saml.go"):                     apiSAMLModelGo(),
		filepath.Join(apiRoot, "internal", "services", "saml.go"):                   apiSAMLServiceGo(),
		filepath.Join(apiRoot, "internal", "handlers", "saml.go"):                   apiSAMLHandlerGo(),
		filepath.Join(apiRoot, "internal", "handlers", "auth.go"):                   apiAuthHandlerGo(),
		filepath.Join(apiRoot, "internal", "middleware", "auth.go"):                 apiAuthMiddlewareGo(),
		filepath.Join(apiRoot, "internal", "middleware", "cors.go"):                 apiCorsMiddlewareGo(),
		filepath.Join(apiRoot, "internal", "middleware", "logger.go"):               apiLoggerMiddlewareGo(),
		filepath.Join(apiRoot, "internal", "middleware", "limits.go"):               middlewareLimitsGo(),
		filepath.Join(apiRoot, "internal", "middleware", "limits_test.go"):          middlewareLimitsTestGo(),
		filepath.Join(apiRoot, "internal", "middleware", "gzip.go"):                 middlewareGzipGo(),
		filepath.Join(apiRoot, "internal", "middleware", "gzip_test.go"):            middlewareGzipTestGo(),
		filepath.Join(apiRoot, "internal", "middleware", "proxies.go"):              middlewareProxiesGo(),
		filepath.Join(apiRoot, "internal", "middleware", "proxies_test.go"):         middlewareProxiesTestGo(),
		filepath.Join(apiRoot, "internal", "middleware", "maintenance.go"):          apiMaintenanceMiddlewareGo(),
		filepath.Join(apiRoot, "internal", "middleware", "idempotency.go"):          apiIdempotencyMiddlewareGo(),
		filepath.Join(apiRoot, "internal", "realtime", "hub.go"):                    apiRealtimeHubGo(),
		filepath.Join(apiRoot, "internal", "realtime", "backplane.go"):              apiRealtimeBackplaneGo(),
		filepath.Join(apiRoot, "internal", "realtime", "backplane_test.go"):         apiRealtimeBackplaneTestGo(),
		filepath.Join(apiRoot, "internal", "realtime", "guard.go"):                  apiRealtimeGuardGo(),
		filepath.Join(apiRoot, "internal", "realtime", "guard_test.go"):             apiRealtimeGuardTestGo(),
		filepath.Join(apiRoot, "internal", "realtime", "channels.go"):               apiRealtimeChannelsGo(),
		filepath.Join(apiRoot, "internal", "realtime", "channels_test.go"):          apiRealtimeChannelsTestGo(),
		filepath.Join(apiRoot, "internal", "realtime", "presence.go"):               apiRealtimePresenceGo(),
		filepath.Join(apiRoot, "internal", "realtime", "presence_test.go"):          apiRealtimePresenceTestGo(),
		filepath.Join(apiRoot, "internal", "realtime", "hub_send_test.go"):          apiRealtimeHubSendTestGo(),
		filepath.Join(apiRoot, "internal", "realtime", "whispers.go"):               apiRealtimeWhispersGo(),
		filepath.Join(apiRoot, "internal", "realtime", "whispers_test.go"):          apiRealtimeWhispersTestGo(),
		filepath.Join(apiRoot, "internal", "realtime", "stats.go"):                  apiRealtimeStatsGo(),
		filepath.Join(apiRoot, "internal", "handlers", "realtime_greeting_test.go"): apiRealtimeGreetingTestGo(),
		filepath.Join(apiRoot, "internal", "handlers", "realtime_channels_test.go"): apiRealtimeChannelsHandlerTestGo(),
		filepath.Join(apiRoot, "internal", "handlers", "realtime_test.go"):          apiRealtimeHandlerTestGo(),
		filepath.Join(apiRoot, "internal", "handlers", "realtime.go"):               apiRealtimeHandlerGo(),
		filepath.Join(apiRoot, "internal", "sync", "registry.go"):                   apiSyncRegistryGo(),
		filepath.Join(apiRoot, "internal", "sync", "policy.go"):                     apiSyncPolicyGo(),
		filepath.Join(apiRoot, "internal", "sync", "softdelete.go"):                 syncSoftDeleteGo(),
		filepath.Join(apiRoot, "internal", "sync", "softdelete_test.go"):            syncSoftDeleteTestGo(),
		filepath.Join(apiRoot, "internal", "settings", "settings.go"):               apiSettingsRegistryGo(),
		filepath.Join(apiRoot, "internal", "settings", "store.go"):                  apiSettingsStoreGo(),
		filepath.Join(apiRoot, "internal", "settings", "defaults.go"):               apiSettingsDefaultsGo(),
		filepath.Join(apiRoot, "internal", "models", "setting.go"):                  apiSettingsModelGo(),
		filepath.Join(apiRoot, "internal", "handlers", "settings.go"):               apiSettingsHandlerGo(),
		filepath.Join(apiRoot, "internal", "services", "event_subscribers.go"):      apiEventsSubscribersGo(),
		filepath.Join(apiRoot, "internal", "handlers", "sync.go"):                   apiSyncHandlerGo(),
		filepath.Join(apiRoot, "internal", "models", "activity_log.go"):             apiActivityLogModelGo(),
		filepath.Join(apiRoot, "internal", "middleware", "activity.go"):             apiActivityMiddlewareGo(),
		filepath.Join(apiRoot, "internal", "middleware", "activity_read_test.go"):   apiActivityReadTestGo(),
		filepath.Join(apiRoot, "internal", "password", "password.go"):               apiPasswordRulesGo(),
		filepath.Join(apiRoot, "internal", "password", "password_test.go"):          apiPasswordRulesTestGo(),
		filepath.Join(apiRoot, "internal", "audit", "audit.go"):                     apiAuditGo(),
		filepath.Join(apiRoot, "internal", "audit", "chain_test.go"):                apiAuditChainTestGo(),
		filepath.Join(apiRoot, "internal", "cluster", "cluster.go"):                 apiClusterGo(),
		filepath.Join(apiRoot, "internal", "cluster", "cluster_test.go"):            apiClusterTestGo(),
		filepath.Join(apiRoot, "internal", "concurrency", "concurrency.go"):         apiConcurrencyGo(),
		filepath.Join(apiRoot, "internal", "concurrency", "concurrency_test.go"):    apiConcurrencyTestGo(),
		filepath.Join(apiRoot, "internal", "webhooks", "verifiers.go"):              apiWebhooksVerifiersGo(),

		// v3.30 — semantic UserActivity log + ticket system
		filepath.Join(apiRoot, "internal", "models", "user_activity.go"):     userActivityModelGo(),
		filepath.Join(apiRoot, "internal", "services", "activity.go"):        userActivityServiceGo(),
		filepath.Join(apiRoot, "internal", "services", "activity_writer.go"): servicesActivityWriterGo(),
		// v3.31.49 — ResolveClientIP honours the X-Public-IP-Hint
		// header sent by the admin/web clients when the TCP peer is
		// loopback, so dev activity logs show the operator's actual
		// public IP instead of "::1".
		filepath.Join(apiRoot, "internal", "services", "clientip.go"):       clientIPHelperGo(),
		filepath.Join(apiRoot, "internal", "services", "ocsf.go"):           apiOCSFServiceGo(),
		filepath.Join(apiRoot, "internal", "services", "ocsf_test.go"):      apiOCSFTestGo(),
		filepath.Join(apiRoot, "internal", "models", "access_review.go"):    apiAccessReviewModelGo(),
		filepath.Join(apiRoot, "internal", "models", "deletion_journal.go"): apiGDPRModelGo(),
		filepath.Join(apiRoot, "internal", "services", "gdpr.go"):           apiGDPRServiceGo(),
		filepath.Join(apiRoot, "internal", "handlers", "gdpr.go"):           apiGDPRHandlerGo(),
		filepath.Join(apiRoot, "internal", "services", "gdpr_test.go"):      apiGDPRTestGo(),
		filepath.Join(apiRoot, "internal", "models", "erasable.go"):         apiErasableModelsGo(),
		filepath.Join(apiRoot, "internal", "erasure", "erasure.go"):         apiErasureGo(),
		filepath.Join(apiRoot, "internal", "erasure", "erasure_test.go"):    apiErasureTestGo(),
		filepath.Join(apiRoot, "internal", "crypto", "field.go"):            apiCryptoFieldGo(),
		filepath.Join(apiRoot, "internal", "crypto", "field_test.go"):       apiCryptoFieldTestGo(),
		filepath.Join(apiRoot, "internal", "crypto", "map_update_test.go"):  apiCryptoMapUpdateTestGo(),
		filepath.Join(apiRoot, "internal", "sanitize", "html.go"):           apiSanitizeHTMLGo(),
		filepath.Join(apiRoot, "internal", "sanitize", "html_test.go"):      apiSanitizeHTMLTestGo(),
		// v3.31.40 — per-user dashboard customisation
		filepath.Join(apiRoot, "internal", "models", "dashboard_layout.go"):   dashboardLayoutModelGo(),
		filepath.Join(apiRoot, "internal", "handlers", "dashboard_layout.go"): strings.ReplaceAll(dashboardLayoutHandlerGo(), "{{MODULE}}", opts.Module()),
		filepath.Join(apiRoot, "internal", "models", "ticket.go"):             ticketModelGo(),
		filepath.Join(apiRoot, "internal", "handlers", "ticket.go"):           ticketHandlerGo(),
		// v3.31.68 — background CSV import job tracking (shared across resources)
		filepath.Join(apiRoot, "internal", "models", "import_job.go"):    importJobModelGo(),
		filepath.Join(apiRoot, "internal", "handlers", "import_job.go"):  importJobHandlerGo(),
		filepath.Join(apiRoot, "internal", "services", "ticket_mail.go"): ticketMailGo(),
		// M29: the ticket rules and queries, out of the handler.
		filepath.Join(apiRoot, "internal", "services", "ticket.go"):  ticketServiceGo(),
		filepath.Join(apiRoot, "internal", "routes", "routes.go"):    apiRoutesGo(),
		filepath.Join(apiRoot, "internal", "routes", "resources.go"): apiRoutesRegistryGo(),
		filepath.Join(apiRoot, "internal", "routes", "apidocs.go"):   apiDocsRoutesGo(),
		filepath.Join(apiRoot, ".air.toml"):                          airConfig(opts),
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
		filepath.Join(apiRoot, "internal", "handlers", "user_test.go"):               apiUserTestGo(),
		filepath.Join(apiRoot, "internal", "handlers", "bench_test.go"):              apiBenchTestGo(),
	}
	// The checks behind the formatted field types (email, tel, json, ...).
	for path, content := range fieldTypesFiles(apiRoot) {
		files[path] = content
	}

	// The packages generated code imports. Their own function because
	// upgrade calls it too; see writeCodegenRuntimeFiles.
	if err := writeCodegenRuntimeFiles(root, opts); err != nil {
		return err
	}

	// Models the framework owns and keeps in step with its own code.
	if err := writeFrameworkOwnedFiles(root, opts); err != nil {
		return err
	}
	// internal/respond, through the same writer the upgrade uses.
	//
	// These four files used to be listed in the map above as well, so that a
	// new project got them and an upgrade got them, from two places. Both
	// lists then had to be kept in step, and twice they were not: codes.go
	// reached only upgraded projects, and validation.go only new ones, each
	// time leaving generated handlers calling a function that was not there.
	// One writer, called from both paths.
	if err := writeRespondFiles(root, opts); err != nil {
		return err
	}
	// The project's own sagas package, created once. The generator injects
	// into it, so it is never rewritten.
	if _, err := writeSagaRegistry(root, opts); err != nil {
		return err
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
` + importJobCreatedByField + `	CreatedAt time.Time ` + "`json:\"created_at\"`" + `
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
func importJobHandlerGo() string { return tmpl("api/handlers/import_job.go") }

func apiGoMod(opts Options) string {
	return fmt.Sprintf(`module %s

// Go 1.26.6, and not an earlier patch. Go 1.26.4 carries eight standard library
// vulnerabilities this code reaches, among them net/http and the encoding/xml
// behind the SAML sign-in path, before anyone has signed in. The go directive is
// also what CI, setup-go and the release workflow install, and an older Go
// downloads this one, so raising it here raises it everywhere.
go 1.26.6

require (
	github.com/MUKE-coder/gin-docs v0.0.0-20260222113017-4d647cb4e7aa
	github.com/MUKE-coder/gorm-studio v1.1.1
	github.com/MUKE-coder/pulse v1.2.0
	github.com/aws/aws-sdk-go-v2 v1.43.0
	github.com/aws/aws-sdk-go-v2/config v1.32.31
	github.com/aws/aws-sdk-go-v2/credentials v1.19.30
	github.com/aws/aws-sdk-go-v2/service/s3 v1.106.0
	// The QR code on the 2FA setup screen. Replaced skip2/go-qrcode, last
	// released in 2020.
	github.com/boombuler/barcode v1.1.0
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
	// Pure Go WebAuthn, so passkeys do not cost the static binary.
	github.com/go-webauthn/webauthn v0.18.0
	github.com/gin-gonic/gin v1.11.0
	// gin's own validator, named directly because internal/respond runs the
	// binding: tags on values that did not arrive as a request body: a grid of
	// rows decoded with encoding/json never passes through gin's binding, so
	// required and min would not fire on that path.
	github.com/go-playground/validator/v10 v10.27.0
	github.com/go-pdf/fpdf v1.4.3
	github.com/golang-jwt/jwt/v5 v5.3.1
	github.com/google/uuid v1.6.0
	// Linked in through goth's gothic, which still asks for v1.6.2 from 2018.
	github.com/gorilla/mux v1.8.1
	github.com/gorilla/sessions v1.4.0
	github.com/gorilla/websocket v1.5.3
	github.com/hibiken/asynq v0.24.1
	github.com/markbates/goth v1.80.0
	github.com/joho/godotenv v1.5.1
	// The HTML sanitiser behind internal/sanitize. Rich text is rendered as HTML
	// by the blog and by anything built on a richtext field, so it is cleaned on
	// the way in rather than trusted on the way out.
	github.com/microcosm-cc/bluemonday v1.0.27
	github.com/redis/go-redis/v9 v9.22.0
	// CVE-2026-54063 and CVE-2026-59161: the worksheet and streaming
	// parsers could be made to allocate without bound, and this is what
	// the CSV/XLSX importer hands user uploads to.
	github.com/xuri/excelize/v2 v2.11.0
	// OpenTelemetry, for internal/tracing. Four modules and no more: the
	// OTLP/HTTP exporter rather than both transports, and a hand-written gin
	// middleware rather than otelgin, because the span has to be named after
	// the route pattern and carry this project's request id, which a
	// general-purpose middleware cannot know about.
	//
	// Inert until OTEL_EXPORTER_OTLP_ENDPOINT is set: with no endpoint the
	// global provider is a no-op and a span costs an allocation.
	go.opentelemetry.io/otel v1.45.0
	// Pulled in by the OTLP exporter even on the HTTP transport: the
	// collector proto package depends on it. Pinned rather than left to
	// whatever the otel modules ask for, which was v1.77.0 and carries two
	// reachable advisories, and v1.83.1 carries CVE-2026-84445.
	google.golang.org/grpc v1.83.2
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp v1.45.0
	go.opentelemetry.io/otel/sdk v1.45.0
	go.opentelemetry.io/otel/trace v1.45.0
	golang.org/x/crypto v0.57.0
	// singleflight, which collapses concurrent cache misses in middleware/cache.go.
	golang.org/x/sync v0.23.0
	// Sentinel now ships a proper /v2 module path, so we track real tags.
	// v2.1.1 is the minimum safe release for WAF.Mode = ModeBlock: v2.1.0
	// fixed the SSRF rule matching "0.0.0.0" inside a Chrome User-Agent
	// (403'ing every Chrome 140/130/120/110 user), and v2.1.1 fixed
	// SQLi_Basic matching a bare "--" inside JWT cookies (roughly one
	// session in ten 403'd at random). v2.2.0 adds ValidateConfig, which
	// Mount runs at startup so dead config shows up in the boot log rather
	// than as a 403 weeks later. Do not downgrade below v2.1.1.
	// v2.2.2 is a security release (client-IP spoofing behind a proxy, a
	// sort_by SQL injection, SSRF bypasses), and v2.5.0 keeps rate limits and
	// lockouts in Redis, so every replica counts against the same numbers.
	github.com/MUKE-coder/sentinel/v2 v2.6.0
	gorm.io/datatypes v1.2.7
	gorm.io/driver/mysql v1.6.0
	gorm.io/driver/postgres v1.6.0
	gorm.io/gorm v1.31.1
	gorm.io/plugin/dbresolver v1.6.2
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
	// CVE-2026-33487: XML Digital Signature validation could be bypassed. On the
	// SAML assertion path that is an authentication bypass, and crewjam/saml
	// v0.5.1 is its newest release and still asks for the vulnerable version.
	github.com/russellhaering/goxmldsig v1.6.0 // CVE-2026-33487
	golang.org/x/image v0.45.0 // GO-2026-5066, -5062, -5032, -5031, -4815, CVE-2026-46603
	golang.org/x/oauth2 v0.27.0 // CVE-2025-22868
	golang.org/x/text v0.39.0 // GO-2026-5970
	// Pulled in by the MySQL driver; govulncheck flags releases before v1.1.1.
	filippo.io/edwards25519 v1.2.0
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
` + apiGitignoreLocalStorage
}

// airMainPackage is the package air rebuilds on every change.
//
// A monorepo keeps its entry point in cmd/server. A single project keeps
// main.go at the top of its module, so there is no cmd/server to build: air
// failed with "directory not found", the API never came up, and `grit start`
// shut the frontend down with it. What the user saw was a Vite dev server
// proxying to nothing and a login that could not reach the API, with the real
// error twenty lines up the log.
func airMainPackage(opts Options) string {
	if opts.Architecture == ArchSingle && !opts.LegacySingleFlat {
		return "."
	}
	return "./cmd/server"
}

func airConfig(opts Options) string {
	// air v1.64+ deprecated `build.bin` in favour of `build.entrypoint`.
	// Both name the BUILT binary that air execs after each rebuild —
	// not the Go source directory. Always use a .exe suffix so Windows
	// CreateProcess accepts the file (no "open with" dialog); Linux
	// + macOS treat .exe as just part of the name. One config, every
	// platform.
	return `root = "."
tmp_dir = "tmp"

[build]
  cmd = "go build -o ./tmp/server.exe ` + airMainPackage(opts) + `"
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
	"errors"
	"crypto/sha256"
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

	"` + "{{MODULE}}" + `/internal/ai"
	"` + "{{MODULE}}" + `/internal/cache"
	"` + "{{MODULE}}" + `/internal/config"
	"` + "{{MODULE}}" + `/internal/cron"
	"` + "{{MODULE}}" + `/internal/database"
	"` + "{{MODULE}}" + `/internal/events"
	"` + "{{MODULE}}" + `/internal/jobs"
	"` + "{{MODULE}}" + `/internal/mail"
	"` + "{{MODULE}}" + `/internal/models"
	"` + "{{MODULE}}" + `/internal/routes"
	"` + "{{MODULE}}" + `/internal/services"
	"` + "{{MODULE}}" + `/internal/storage"
	"` + "{{MODULE}}" + `/internal/tracing"
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

` + mainStorageInit + `
` + mailerInitNew + `
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
` + gothicStoreNew + `	var oauthProviders []goth.Provider
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
` + serverTimeoutFields + `	}

	// Start server in goroutine
	go func() {
		log.Printf("Server starting on port %s", cfg.Port)
` + serverDashboardLogs + `
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
` + shutdownDrain + `
	log.Println("Server exited")
}
`
}

func apiConfigGo() string {
	return `package config

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
` + configStorageFields + `
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
` + configStorageDriverFieldNew + configStorageDisksField + `
` + configMailFieldsNew + `
	CORSOrigins []string

` + configUploadMIMEField + `
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
	TOTPIssuer string
	RequireEmailVerification bool
	LoginMaxAttempts         int
	LoginLockoutWindow       time.Duration

	// Security (Sentinel)
	SentinelEnabled        bool
	SentinelUsername       string
	SentinelPassword       string
	SentinelSecretKey      string
	SentinelAuditKey       string
	// Sentinel v2.0 — CIDRs allowed to send X-Forwarded-For / X-Real-IP.
	// Empty (default) means "ignore those headers entirely" — safe when
	// the app speaks to the public internet directly; populate when
	// you're behind a known reverse proxy (Caddy/Traefik/Cloudflare).
	SentinelTrustedProxies []string
` + configTrustedProxiesField + `
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

	// Serve the API reference at /docs in production (API_DOCS_PUBLIC).
	APIDocsPublic bool
}

// Load reads configuration from environment variables.
func Load() (*Config, error) {
	// Load .env file (ignore error if not found — production uses real env vars)
	// The project root holds one .env, and the binary runs from wherever its
	// module is: the root itself, api/ one level down, or apps/api two. Each
	// candidate is tried and the first that exists wins, because godotenv does
	// not overwrite a variable that is already set.
	//
	// Without the one-level case, a single project's migrate and seed ran from
	// api/, found no .env at all, and fell back to the built-in defaults: the
	// first thing anybody saw was a refusal to start under APP_ENV=production,
	// naming five secrets that were sitting in a file one directory up.
	_ = godotenv.Load()
	_ = godotenv.Load("../.env")
	_ = godotenv.Load("../../.env")

` + configStorageDriverNew + `
	cfg := &Config{
		AppName:     getEnv("APP_NAME", "grit-app"),
		// Production unless told otherwise: a server that forgot APP_ENV is strict.
		AppEnv:      getEnv("APP_ENV", "production"),
		// PORT first: every platform that routes to a container injects it, and
		// a deploy that binds the wrong port goes green with nothing answering.
		// APP_PORT stays the one you set yourself, and wins locally because a
		// platform is not the thing setting PORT there.
		Port:        firstNonEmpty(os.Getenv("PORT"), os.Getenv("APP_PORT"), "8080"),
		AppURL:      resolveAppURL(firstNonEmpty(os.Getenv("PORT"), os.Getenv("APP_PORT"), "8080")),
		DatabaseURL: resolveDatabaseURL(),
		JWTSecret:   getEnv("JWT_SECRET", ""),
		FieldEncryptionKey: getEnv("FIELD_ENCRYPTION_KEY", ""),
		RedisURL:    resolveRedisURL(),

		StorageDriver: storageDriver,
` + configStorageDisksLoad + `
		ResendAPIKey: getEnv("RESEND_API_KEY", ""),
` + configMailLoadNew + `
		// The Wails desktop webview is allowed by middleware.isWailsOrigin (it
		// matches the wails.localhost host on any port), so it needs no entry
		// here — its dev origin includes a configurable port.
		CORSOrigins: strings.Split(getEnv("CORS_ORIGINS", "http://localhost:3000,http://localhost:3001"), ","),

` + configUploadMIMELoad + `
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
		LoginMaxAttempts:   getEnvInt("LOGIN_MAX_ATTEMPTS", 10),
		LoginLockoutWindow: getEnvDuration("LOGIN_LOCKOUT_MINUTES", 15, time.Minute),

		SentinelEnabled:        getEnv("SENTINEL_ENABLED", "true") == "true",
		SentinelUsername:       getEnv("SENTINEL_USERNAME", "admin"),
		SentinelPassword:       getEnv("SENTINEL_PASSWORD", "sentinel"),
		SentinelSecretKey:      getEnv("SENTINEL_SECRET_KEY", "sentinel-secret-change-me"),
		SentinelAuditKey:       getEnv("SENTINEL_AUDIT_KEY", ""),
		` + configSentinelProxiesNew + `

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

` + configLocalStorageCheck + `	// DatabaseURL is always populated by resolveDatabaseURL() — either from
	// the DATABASE_URL env var or built from POSTGRES_* parts. The actual
	// connection attempt in cmd/server/main.go will surface a useful error
	// if the resolved URL points at an unreachable database.

` + configProductionBlock + `	if cfg.JWTSecret == "" {
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

` + configCheckSecretsFunc + `
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

// resolveAppURL reads APP_URL and says so when it names a different port from
// the one the server is about to listen on.
//
// APP_URL is where the outside world reaches this app, so behind a proxy it is
// https://shop.example.com while the server listens on 8080, and that is
// correct. On localhost it is not: a local APP_URL whose port disagrees with
// APP_PORT is always a mistake, and a silent one. Every URL built from it is
// wrong, which with STORAGE_DRIVER=local means every uploaded file's URL points
// at a port where nothing is listening, and the symptom is an image that does
// not load rather than anything that mentions APP_URL.
func resolveAppURL(port string) string {
	appURL := getEnv("APP_URL", "http://localhost:8080")
	u, err := url.Parse(appURL)
	if err != nil || u.Host == "" {
		return appURL
	}
	host, urlPort, err := net.SplitHostPort(u.Host)
	if err != nil || urlPort == port {
		return appURL
	}
	if host != "localhost" && host != "127.0.0.1" && host != "[::1]" {
		// A real hostname: a proxy in front, which is the normal case.
		return appURL
	}
	log.Printf("WARNING: APP_URL is %s but this server listens on port %s. Every URL the app builds, uploaded files included, will point at port %s where nothing is answering. Set APP_URL=http://%s:%s",
		appURL, port, urlPort, host, port)
	return appURL
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

` + configResolveStorageDriverFunc + configStorageDisksFuncs + `// resolveStorage returns the StorageConfig for the active driver.
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
` + configLocalCase + `	default: // minio
		return StorageConfig{
			Endpoint:  resolveMinioEndpoint(),
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

` + configTrustedProxiesFunc + `// splitCSV trims and splits a comma-separated env var. Empty strings
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

` + configMailFuncs
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

	"{{MODULE}}/internal/appendonly"
	"{{MODULE}}/internal/crypto"
	"{{MODULE}}/internal/fieldtypes"
	"{{MODULE}}/internal/paginate"
	"{{MODULE}}/internal/sanitize"
	"{{MODULE}}/internal/sync"
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
		db, err = gorm.Open(sqlite.Open(sqliteDSN(strings.TrimPrefix(dsn, "sqlite://"))), gormCfg)
	case strings.HasPrefix(dsn, "sqlite:"):
		db, err = gorm.Open(sqlite.Open(sqliteDSN(strings.TrimPrefix(dsn, "sqlite:"))), gormCfg)
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

	// Map-based updates to an EncryptedString column are encrypted too. GORM only
	// runs a column type's Value() when the value already has that type, and the
	// generated update and PATCH handlers write maps of plain strings, so without
	// this the first edit to an encrypted column stored plaintext.
	if err := crypto.Install(db); err != nil {
		return nil, fmt.Errorf("installing field encryption for map updates: %w", err)
	}

	// Rich text is sanitised on its way into the database: every field tagged
	// sanitize:"html", through every write that has a model (create, update,
	// PATCH, bulk edit, CSV import, sync push, GORM Studio's row editor).
	if err := sanitize.Install(db); err != nil {
		return nil, fmt.Errorf("installing the HTML sanitiser: %w", err)
	}

	// Append-only tables refuse UPDATE and DELETE through this handle, so GORM
	// Studio's row editor, the CSV importer and every handler are bound by it,
	// not just the one service that remembered. See internal/appendonly.
	if err := appendonly.Install(db); err != nil {
		return nil, fmt.Errorf("installing append-only guard: %w", err)
	}

` + paginateConnectHook + syncConnectHook + fieldTypesConnectHook + `	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

` + dbPoolBlockNew + `
	log.Println("Database connected successfully")
	return db, nil
}

// sqliteDSN adds the two pragmas a SQLite file needs to survive more than one
// writer.
//
// Without them a second writer fails at once with "database is locked", which
// is not a load problem: a background worker ticking while a webhook writes is
// enough, and the error surfaces as a failed request nobody can reproduce by
// hand. busy_timeout makes a writer wait its turn instead of giving up, and
// WAL lets reads carry on while one write is in flight.
//
// Anything the caller set is left alone, and :memory: gets neither: there is
// no file to journal and a memory database is one connection's own.
func sqliteDSN(path string) string {
	if path == "" || strings.Contains(path, ":memory:") || strings.Contains(path, "mode=memory") {
		return path
	}
	if !strings.Contains(path, "busy_timeout") {
		path += pragmaJoin(path) + "_pragma=busy_timeout(5000)"
	}
	if !strings.Contains(path, "journal_mode") {
		path += pragmaJoin(path) + "_pragma=journal_mode(WAL)"
	}
	return path
}

func pragmaJoin(path string) string {
	if strings.Contains(path, "?") {
		return "&"
	}
	return "?"
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
` + sentinelPoolFunc
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

	"{{MODULE}}/internal/appendonly"
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
	CreatedAt       time.Time      ` + "`" + `gorm:"index" json:"created_at"` + "`" + `
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
		&RecoveryContactToken{},
		&RecoveryContact{},
		&Passkey{},
		&WebAuthnSession{},
		&APIKey{},
		// Role/UserRole must migrate before anything authorises a request.
		&Role{},
		&UserRole{},
		&Upload{},
		&Blog{},
		&TwoFactorConfig{},
		&TrustedDevice{},
		&TOTPPendingToken{},
		&MagicLinkToken{},
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
		// The transactional outbox. Registered here like every other table so
		// AutoMigrate creates it and the backup writer includes it; an outbox
		// missing from a backup loses events nobody knows were pending.
		&OutboxMessage{},
` + sagaModelRegistration + `		// grit:models
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
	// And at the database, for whatever does not go through GORM: the SQL
	// editor in Studio, psql, a script somebody writes next year.
	if err := appendonly.InstallTriggers(db); err != nil {
		return fmt.Errorf("installing append-only triggers: %w", err)
	}

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
	CreatedAt    time.Time      ` + "`" + `gorm:"index" json:"created_at"` + "`" + `
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

func apiAuthServiceGo() string { return tmpl("api/services/auth.go") }

func apiAuthHandlerGo() string { return tmpl("api/handlers/auth.go") }

func apiUserHandlerGo() string { return tmpl("api/handlers/user.go") }

func apiAuthMiddlewareGo() string {
	return `package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"` + "{{MODULE}}" + `/internal/authz"
	"` + "{{MODULE}}" + `/internal/models"
	"` + "{{MODULE}}" + `/internal/services"
	"` + "{{MODULE}}" + `/internal/respond"
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
				respond.Fail(c, respond.CodeUnauthorized, "Invalid authorization header format")
				c.Abort()
				return
			}
			token = parts[1]
		}

		if token == "" {
			respond.Fail(c, respond.CodeUnauthorized, "Authentication required")
			c.Abort()
			return
		}

		// An access token, whose session is still live: a refresh token, or a
		// token from a session that logged out, is refused here.
		claims, err := authService.ValidateAccessToken(token)
		if err != nil {
			respond.Fail(c, respond.CodeUnauthorized, "Invalid or expired token")
			c.Abort()
			return
		}

		// Load user from database.
		// Use Where("id = ?") rather than First(&user, id) — GORM's shorthand
		// emits the bare value into the WHERE clause and Postgres rejects UUID
		// primary keys with "trailing junk after numeric literal".
		var user models.User
		if err := db.WithContext(c.Request.Context()).Where("id = ?", claims.UserID).First(&user).Error; err != nil {
			respond.Fail(c, respond.CodeUnauthorized, "User not found")
			c.Abort()
			return
		}

		if !user.Active {
			respond.Fail(c, respond.CodeAccountDisabled, "Your account has been disabled")
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

// Identify is Auth without the wall: it reads the session if there is one and
// lets the request through either way.
//
// It exists for the public pages that are not the same page for everybody. A
// catalogue that marks what the reader already owns, a pricing page that knows
// their current plan, an article with their own comment on it: each has to be
// readable signed out, which means no guard, which with Auth alone means the
// handler cannot tell who is reading even when they are signed in. Every app
// that wants this ends up parsing the Authorization header by hand in a
// handler, and that is how the cookie flow gets forgotten.
//
// It sets exactly what Auth sets, so c.GetString("user_id") reads the same on
// both kinds of route, and nothing else changes: a missing, malformed, expired
// or revoked token is simply an anonymous request, not a 401.
//
// It is not a guard and cannot be used as one. A handler behind Identify must
// treat "there is a user" as information, never as permission: anything that
// must not be served to a stranger belongs behind Auth.
func Identify(db *gorm.DB, authService *services.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Same two places Auth looks, in the same order: the HttpOnly cookie
		// browsers use, then the bearer header native clients use.
		token := ""
		if cookieValue, err := c.Cookie("grit_access"); err == nil && cookieValue != "" {
			token = cookieValue
		} else if authHeader := c.GetHeader("Authorization"); authHeader != "" {
			if parts := strings.SplitN(authHeader, " ", 2); len(parts) == 2 && parts[0] == "Bearer" {
				token = parts[1]
			}
		}
		if token == "" {
			c.Next()
			return
		}

		claims, err := authService.ValidateAccessToken(token)
		if err != nil {
			c.Next()
			return
		}

		var user models.User
		if err := db.WithContext(c.Request.Context()).Where("id = ?", claims.UserID).First(&user).Error; err != nil {
			c.Next()
			return
		}
		if !user.Active {
			c.Next()
			return
		}

		c.Set("user", user)
		c.Set("user_id", user.ID)
		c.Set("user_email", user.Email)
		c.Set("user_role", user.Role)
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
			respond.Fail(c, respond.CodeForbidden, "You do not have permission to perform this action")
			c.Abort()
			return
		}

		userRole, exists := c.Get("user_role")
		if !exists {
			respond.Fail(c, respond.CodeUnauthorized, "Not authenticated")
			c.Abort()
			return
		}

		role, ok := userRole.(string)
		if !ok {
			respond.Fail(c, respond.CodeInternalError, "Invalid user role")
			c.Abort()
			return
		}

		for _, r := range roles {
			if role == r {
				c.Next()
				return
			}
		}

		respond.Fail(c, respond.CodeForbidden, "You do not have permission to access this resource")
		c.Abort()
	}
}

// RequireStaff admits anyone who can do something in the admin: the ADMIN
// role, or any permission at all. It is a gate, not a guard. Every route behind
// it names the permission it needs, and routes that name none stay in the ADMIN
// group, so one that forgets fails closed.
func RequireStaff() gin.HandlerFunc {
	return func(c *gin.Context) {
		if role, _ := c.Get("user_role"); role == models.RoleAdmin {
			c.Next()
			return
		}
		if grants, ok := c.Get("user_grants"); ok {
			if list, ok := grants.([]string); ok && len(list) > 0 {
				c.Next()
				return
			}
		}
		respond.Fail(c, respond.CodeForbidden, "You do not have permission to access this resource")
		c.Abort()
	}
}

// RequirePermissionFor guards a route whose resource is in the URL. The
// dashboard's stats and charts take it as :resource, so the permission they
// need is that resource's, which a route written once cannot name in advance.
func RequirePermissionFor(param, action string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if role, _ := c.Get("user_role"); role == models.RoleAdmin {
			c.Next()
			return
		}
		if grants, ok := c.Get("user_grants"); ok {
			if list, ok := grants.([]string); ok && authz.Granted(list, c.Param(param)+"."+action) {
				c.Next()
				return
			}
		}
		respond.Fail(c, respond.CodeForbidden, "You do not have permission to access this resource")
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
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"{{MODULE}}/internal/respond"
)

// RequestID injects a unique X-Request-ID header into every request and
// stores it in the context for downstream logging and tracing.
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			// #nosec G404 -- a request id has to be unique for tracing, not
			// unpredictable. crypto/rand here would be a syscall on every
			// request for no security gain.
			requestID = fmt.Sprintf("%d-%d", time.Now().UnixNano(), rand.Int63())
		}
		c.Set("request_id", requestID)
		c.Header("X-Request-ID", requestID)
		c.Next()
	}
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
			respond.Fail(c, respond.CodePayloadTooLarge, fmt.Sprintf("Request body exceeds %dMB limit", limit/(1024*1024)))
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
//	        h.DB.WithContext(c.Request.Context()).Model(&models.Shop{}).Preload("Building"),
//	        paginate.Bind(c),
//	        paginate.Config{
//	            Searchable:   []string{"shop_number", "description"},
//	            Sortable:     map[string]bool{"created_at": true, "monthly_rent": true},
//	            DefaultSort:  "created_at",
//	            DefaultOrder: "desc",
//	        },
//	    )
//	    if err != nil {
//	        respond.WriteError(c, err, "Failed to fetch shops")
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

	// CursorMode asks for keyset pagination on this request.
	//
	// Set by ?mode=cursor, and implied by sending a ?cursor= at all, so a
	// client that follows meta.next_cursor never has to say so twice. Offset
	// stays the default because the admin table needs page numbers and a
	// total, which keyset pagination cannot give you.
	CursorMode bool
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

	// Counts names extra totals to return beside the page, from ?counts=:
	// created_7d is rows created in the last 7 days, updated_30d rows updated
	// in the last 30. Each is one COUNT over the same filters and search as
	// Total. The admin's stat cards ask for theirs on the list request rather
	// than sending a request per card.
	Counts []string

	// Series asks for a count per period, from ?series=: created_at:month:12
	// is the last twelve months by the month a row was created. The column is
	// created_at or updated_at, the unit is day, week or month.
	//
	// Raw as it arrived, because what is allowed depends on Config, which
	// ParamsFrom does not have. List parses it.
	Series string

	// Breakdown asks for a count per value of a column, from ?breakdown=:
	// status,role is how many rows hold each status and each role. Only
	// columns in Config.Filterable are answered.
	Breakdown string
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

	// Mode is "cursor" on a keyset response and absent on an offset one.
	//
	// It is here because total, page and pages are all zero in cursor mode,
	// and a client cannot otherwise tell that apart from an empty table. The
	// counts are zero because counting the whole set on every page is the
	// cost keyset pagination exists to avoid; ask for it with
	// Config.IncludeTotal when you genuinely need it.
	Mode string ` + "`" + `json:"mode,omitempty"` + "`" + `

	// Counts answers ?counts=, keyed by the names asked for. See Params.Counts.
	Counts map[string]int64 ` + "`" + `json:"counts,omitempty"` + "`" + `

	// Series answers ?series=: one entry per period that has rows, oldest
	// first. Absent when the request did not ask.
	Series []Bucket ` + "`" + `json:"series,omitempty"` + "`" + `

	// Breakdown answers ?breakdown=: the values of each column asked for, with
	// how many rows hold each, largest first.
	Breakdown map[string][]Slice ` + "`" + `json:"breakdown,omitempty"` + "`" + `
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
		CursorMode:   c.Query("mode") == "cursor" || c.Query("cursor") != "",
		DateField:    dateField,
		DateFrom:     dateFrom,
		DateTo:       dateTo,
		QueryFilters: collectQueryFilters(c),
		Counts:       parseCounts(c.Query("counts")),
		Series:       c.Query("series"),
		Breakdown:    c.Query("breakdown"),
	}
}

// reservedParams are the query keys pagination owns. Everything else is a
// candidate column filter, to be checked against Config.Filterable before it
// is used.
var reservedParams = map[string]bool{
	"page": true, "page_size": true, "search": true,
	"sort_by": true, "sort_order": true, "cursor": true,
	"date_field": true, "date_from": true, "date_to": true,
	"mode": true,
	"created_since": true, "created_from": true, "created_to": true,
	"updated_since": true, "archived": true, "format": true,
	"counts": true, "series": true, "breakdown": true,
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

	// Either the handler insists, or the request asked.
	if cfg.CursorMode || p.CursorMode {
		return listCursor[T](query, p, cfg, sortBy, sortOrder)
	}

	var result Result[T]

	// Count first (before Order/Offset/Limit so Count reflects the whole match).
	if err := countTotal(query, &result.Meta.Total); err != nil {
		return result, err
	}
	if len(p.Counts) > 0 {
		counts, err := countWindows(query, p.Counts)
		if err != nil {
			return result, err
		}
		result.Meta.Counts = counts
	}

	// The insights panel's chart, over this same query: the rows the table is
	// showing, not the whole table. Collapsed by default in the admin, so the
	// GROUP BY runs when somebody opens the panel and not before.
	if spec, ok := parseSeries(p.Series); ok {
		series, err := seriesBuckets(query, spec)
		if err != nil {
			return result, err
		}
		result.Meta.Series = series
	}
	if columns := parseBreakdown(p.Breakdown, cfg.Filterable); len(columns) > 0 {
		breakdown, err := breakdownSlices(query, columns)
		if err != nil {
			return result, err
		}
		result.Meta.Breakdown = breakdown
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

// maxCounts caps ?counts=, so one request cannot ask for any number of COUNTs.
const maxCounts = 8

// parseCounts keeps the well-formed names from ?counts=: created_<N>d or
// updated_<N>d, with N from 1 to 3660. Anything else is dropped rather than
// failing the list; a stat card asking for a count it cannot have shows a dash.
func parseCounts(raw string) []string {
	if raw == "" {
		return nil
	}
	var out []string
	seen := map[string]bool{}
	for _, name := range strings.Split(raw, ",") {
		name = strings.TrimSpace(name)
		if _, _, ok := countWindow(name); !ok || seen[name] {
			continue
		}
		seen[name] = true
		out = append(out, name)
		if len(out) == maxCounts {
			break
		}
	}
	return out
}

// countWindow turns created_7d into the column it filters and how far back.
// The column is one of two literals, never text from the request.
func countWindow(name string) (string, time.Duration, bool) {
	var column, rest string
	switch {
	case strings.HasPrefix(name, "created_"):
		column, rest = "created_at", strings.TrimPrefix(name, "created_")
	case strings.HasPrefix(name, "updated_"):
		column, rest = "updated_at", strings.TrimPrefix(name, "updated_")
	default:
		return "", 0, false
	}
	if !strings.HasSuffix(rest, "d") {
		return "", 0, false
	}
	days, err := strconv.Atoi(strings.TrimSuffix(rest, "d"))
	if err != nil || days < 1 || days > 3660 {
		return "", 0, false
	}
	return column, time.Duration(days) * 24 * time.Hour, true
}

// countWindows runs one COUNT per name over query, which carries the list's
// filters, search and date window but no order or page.
func countWindows(query *gorm.DB, names []string) (map[string]int64, error) {
	// To the minute, so the same request moments later builds the same SQL and
	// the count cache can answer it. Compared in local time, as GORM writes
	// created_at: SQLite compares times as text, and a UTC bound against local
	// timestamps put the window off by the machine's UTC offset.
	now := time.Now().UTC().Truncate(time.Minute)
	out := make(map[string]int64, len(names))
	for _, name := range names {
		column, window, ok := countWindow(name)
		if !ok {
			continue
		}
		var n int64
		if err := countTotal(query.Session(&gorm.Session{}).Where(column+" >= ?", now.Add(-window).In(time.Local)), &n); err != nil {
			return nil, fmt.Errorf("counting %s: %w", name, err)
		}
		out[name] = n
	}
	return out, nil
}

// Counts answers ?counts= for a handler that builds its own list query instead
// of calling List. query is that list's query with its filters and search, before
// order and paging. It returns nil when the request asked for nothing.
func Counts(c *gin.Context, query *gorm.DB) (map[string]int64, error) {
	names := parseCounts(c.Query("counts"))
	if len(names) == 0 {
		return nil, nil
	}
	return countWindows(query, names)
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
			//
			// typedCursorValue matters more than it looks: handing the raw
			// string over compares a timestamp column against text, which
			// SQLite answers true for every row, and page two comes back
			// identical to page one.
			query = query.Where(fmt.Sprintf("(%s, id) %s (?, ?)", sortBy, op),
				typedCursorValue(sortVal), lastID)
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
	result.Meta.Mode = "cursor"
	return result, nil
}

// typedCursorValue turns the cursor's text back into something the column can
// be compared against.
//
// The cursor is text because it travels in a URL, but the column is not: a
// timestamp compared against a string, or an integer compared against '500',
// is a different comparison in every database and the wrong one in SQLite.
// Binding the value at its own type lets the driver do what it does for every
// other parameter.
//
// The order is deliberate. RFC3339 first, because that is what extractCursor
// writes for a time and it is specific enough not to swallow anything else.
// Integers before floats, so an id-like value does not become 1.0. A string
// that is none of those was a string in the column too.
func typedCursorValue(s string) any {
	if t, err := time.Parse(time.RFC3339Nano, s); err == nil {
		return t
	}
	if n, err := strconv.ParseInt(s, 10, 64); err == nil {
		return n
	}
	if f, err := strconv.ParseFloat(s, 64); err == nil {
		return f
	}
	return s
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

	sortField := lookupSortField(rv, sortBy)
	if !sortField.IsValid() {
		// No field to read means no usable cursor. Returning the id alone
		// would encode an empty sort value, and the next page would compare
		// against "" and return the whole table from the top.
		return "", ""
	}

	if t, ok := sortField.Interface().(time.Time); ok {
		return t.Format(time.RFC3339Nano), id
	}
	return fmt.Sprintf("%v", sortField.Interface()), id
}

// lookupSortField finds the struct field behind a column name.
//
// Usually that is one field: "created_at" is CreatedAt. An embedded struct is
// two, because GORM flattens it with a prefix: a money field declared as
// Price money.Money becomes the columns price_amount and price_currency, and
// there is no PriceAmount field to find. So a name that does not resolve
// whole is split at each underscore and tried as a path, longest prefix first,
// which finds Price then Amount.
//
// Getting this wrong is not a crash. FieldByName returns an invalid Value, the
// cursor encodes an empty sort value, and the next page compares against ""
// and starts again from the top: an infinite scroll that repeats its first
// page forever.
func lookupSortField(rv reflect.Value, column string) reflect.Value {
	if f := rv.FieldByName(snakeToPascal(column)); f.IsValid() {
		return f
	}

	parts := strings.Split(column, "_")
	for split := len(parts) - 1; split >= 1; split-- {
		outer := rv.FieldByName(snakeToPascal(strings.Join(parts[:split], "_")))
		if !outer.IsValid() {
			continue
		}
		if outer.Kind() == reflect.Ptr {
			if outer.IsNil() {
				return reflect.Value{}
			}
			outer = outer.Elem()
		}
		if outer.Kind() != reflect.Struct {
			continue
		}
		if inner := outer.FieldByName(snakeToPascal(strings.Join(parts[split:], "_"))); inner.IsValid() {
			return inner
		}
	}
	return reflect.Value{}
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
	"os"

	"github.com/gin-gonic/gin"

	"{{MODULE}}/internal/respond"
)

// Maintenance returns a middleware that checks for a .maintenance file.
// When the file exists, all requests receive a 503 Service Unavailable response.
// Toggle with: grit down (enable) / grit up (disable)
func Maintenance() gin.HandlerFunc {
	return func(c *gin.Context) {
		if _, err := os.Stat(".maintenance"); err == nil {
			respond.Fail(c, respond.CodeMaintenance, "Application is in maintenance mode. Please try again later.")
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
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"strings"
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

` + idempotencyKeyNew + `
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
` + idempotencyScopeFunc + `
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

// Shared is the registry the application built, for code that runs outside a
// request and cannot be handed one.
//
// The background worker is the case this exists for: it runs in its own process
// with no access to routes.Setup's local variable, and the retention sweep needs
// the same list of resources the API has. It is the same pointer rather than a
// second registry, which would be empty and would sweep nothing while reporting
// success.
var Shared *Registry

// NewRegistry returns an empty Registry.
func NewRegistry() *Registry {
	r := &Registry{models: make(map[string]reflect.Type)}
	Shared = r
	return r
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

func apiSyncHandlerGo() string { return tmpl("api/handlers/sync.go") }

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
	// What an audited read served, or what a reseal covered. Empty on every
	// other entry, and left out of the hash when empty, so entries written
	// before these fields existed hash exactly as they did.
	Resource    string ` + "`" + `gorm:"size:100;index" json:"resource,omitempty"` + "`" + `
	ResourceIDs string ` + "`" + `gorm:"type:text" json:"resource_ids,omitempty"` + "`" + `
	RecordCount int    ` + "`" + `json:"record_count,omitempty"` + "`" + `
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
	"errors"
	"fmt"
	"log"
	"sync"
	"sync/atomic"
	"time"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
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
		Resource:          e.Resource,
		ResourceIDs:       e.ResourceIDs,
		RecordCount:       e.RecordCount,
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
	// Appended, and omitted when empty, so every entry written before they
	// existed has the same bytes and old chains still verify.
	Resource    string ` + "`" + `json:"resource,omitempty"` + "`" + `
	ResourceIDs string ` + "`" + `json:"resource_ids,omitempty"` + "`" + `
	RecordCount int    ` + "`" + `json:"record_count,omitempty"` + "`" + `
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

// MaxReadIDs caps how many ids one read entry lists. A page is a few hundred
// rows at most; an export is recorded as a count instead.
const MaxReadIDs = 1000

const readKey = "audit.read"

// ReadMark is what a handler served, for the activity middleware to record.
type ReadMark struct {
	Resource string
	IDs      []string
	Count    int
}

// Read marks the request as having served these rows of resource. The
// activity middleware records it in the chain once the response has gone out
// with a 2xx. Handlers generated with --audit-reads call it; any handler can.
func Read(c *gin.Context, resource string, ids ...string) {
	mark := ReadMark{Resource: resource, Count: len(ids)}
	if len(ids) > MaxReadIDs {
		ids = ids[:MaxReadIDs]
	}
	mark.IDs = append([]string(nil), ids...)
	c.Set(readKey, mark)
}

// ReadCount marks a read too large to list row by row: an export.
func ReadCount(c *gin.Context, resource string, n int) {
	c.Set(readKey, ReadMark{Resource: resource, Count: n})
}

// ReadMarkOf returns what the handler marked, if it marked anything.
func ReadMarkOf(c *gin.Context) (ReadMark, bool) {
	v, ok := c.Get(readKey)
	if !ok {
		return ReadMark{}, false
	}
	mark, ok := v.(ReadMark)
	return mark, ok
}

// Precision is what every supported database keeps of a timestamp: Postgres
// stores microseconds and MySQL, as GORM creates the column, milliseconds.
// The hash covers created_at, so a stamp finer than the column gives a hash
// nobody can recompute from the stored row. Before v3.215.0 every entry carried
// nanoseconds, and the chain failed verification on its first row on both.
const Precision = time.Millisecond

// chainLockKey is the Postgres advisory lock every chain writer takes, so two
// API replicas never read the same latest hash and fork the chain.
const chainLockKey int64 = 0x677269745f617564

var (
	queue     = make(chan models.ActivityLog, 4096)
	startOnce sync.Once
	dropped   atomic.Uint64
)

// Start runs this process's chain writer. Safe to call more than once.
//
// One writer fed by a bounded channel, so a burst of requests never waits on
// the database or spawns a goroutine per entry. It is not what keeps the chain
// whole across processes: every batch takes the chain lock and reads the latest
// hash from the database, so each replica can run one.
func Start(db *gorm.DB) {
	startOnce.Do(func() { go writer(db) })
}

// Enqueue hands an entry to the writer without blocking, and reports false
// when the backlog was full and the entry was dropped. Losing an audit row is
// better than stalling the request path; Dropped says how often it happened.
func Enqueue(entry models.ActivityLog) bool {
	select {
	case queue <- entry:
		return true
	default:
		dropped.Add(1)
		return false
	}
}

// Dropped is how many entries Enqueue has dropped since the process started.
// Sustained growth means the writer cannot keep up.
func Dropped() uint64 { return dropped.Load() }

func writer(db *gorm.DB) {
	for first := range queue {
		batch := []models.ActivityLog{first}
	drain:
		for len(batch) < 256 {
			select {
			case e := <-queue:
				batch = append(batch, e)
			default:
				break drain
			}
		}
		if err := appendBatch(db, batch); err != nil {
			// One bad entry must not take the rest of the batch with it.
			for _, e := range batch {
				if err := appendBatch(db, []models.ActivityLog{e}); err != nil {
					log.Printf("[audit] could not record %s %s: %v", e.Method, e.Path, err)
				}
			}
		}
	}
}

// AppendChained writes one entry and returns once it is stored, for a caller
// that must know it landed: a security event, a reseal. Same lock and same
// stamping as the writer, so the two never fork the chain between them.
func AppendChained(db *gorm.DB, entry *models.ActivityLog) error {
	batch := []models.ActivityLog{*entry}
	if err := appendBatch(db, batch); err != nil {
		return err
	}
	*entry = batch[0]
	return nil
}

// appendBatch chains entries onto the latest stored one, in one transaction
// holding the chain lock. created_at is stamped here rather than by the caller:
// at the precision every database keeps, and strictly after the entry before,
// so VerifyChain's (created_at, id) order is the order the chain was written.
func appendBatch(db *gorm.DB, entries []models.ActivityLog) error {
	return db.Transaction(func(tx *gorm.DB) error {
		head, err := lockAndReadHead(tx)
		if err != nil {
			return err
		}
		prevHash, last := head.Hash, head.CreatedAt
		for i := range entries {
			e := &entries[i]
			clipToColumns(e)
			e.CreatedAt = nextStamp(last)
			canonical, err := Canonical(e)
			if err != nil {
				return fmt.Errorf("canonicalize: %w", err)
			}
			e.PrevHash = prevHash
			e.Hash = ComputeHash(prevHash, canonical)
			prevHash, last = e.Hash, e.CreatedAt
		}
		// One multi-row INSERT for the batch. An INSERT per entry kept the chain
		// lock, which every replica's writer waits on, for 256 round trips.
		return tx.CreateInBatches(entries, 256).Error
	})
}

// lockAndReadHead takes the chain lock and returns the latest entry, or a zero
// entry for an empty log. Postgres takes an advisory lock held until the
// transaction ends; MySQL a locking read, which sees the latest committed row
// once granted; SQLite serialises writers on its own.
func lockAndReadHead(tx *gorm.DB) (models.ActivityLog, error) {
	var head models.ActivityLog
	q := tx.Order("created_at desc, id desc").Limit(1)
	switch tx.Dialector.Name() {
	case "postgres":
		if err := tx.Exec("SELECT pg_advisory_xact_lock(?)", chainLockKey).Error; err != nil {
			return head, fmt.Errorf("taking the audit chain lock: %w", err)
		}
	case "mysql":
		q = q.Clauses(clause.Locking{Strength: "UPDATE"})
	}
	if err := q.Find(&head).Error; err != nil {
		return head, fmt.Errorf("reading the chain head: %w", err)
	}
	return head, nil
}

// nextStamp is now at the chain's precision, and strictly after prev.
func nextStamp(prev time.Time) time.Time {
	now := time.Now().UTC().Truncate(Precision)
	if !prev.IsZero() {
		if floor := prev.UTC().Truncate(Precision).Add(Precision); now.Before(floor) {
			now = floor
		}
	}
	return now
}

// clipToColumns trims what could overflow its column, before hashing, so the
// stored row is the hashed row. An oversized user agent used to fail the insert
// and lose the entry.
func clipToColumns(e *models.ActivityLog) {
	e.Method = clip(e.Method, 10)
	e.Path = clip(e.Path, 500)
	e.IPAddress = clip(e.IPAddress, 45)
	e.UserAgent = clip(e.UserAgent, 500)
	e.Resource = clip(e.Resource, 100)
}

// clip shortens s to at most n bytes without splitting a character, which a
// database would refuse as invalid UTF-8.
func clip(s string, n int) string {
	if len(s) <= n {
		return s
	}
	s = s[:n]
	for len(s) > 0 && !utf8.ValidString(s) {
		s = s[:len(s)-1]
	}
	return s
}

// PruneChunk is how many entries one prune transaction deletes.
const PruneChunk = 5000

// prunePath marks the SECURITY entry a prune appends.
const prunePath = "audit.chain.pruned"

// Prune deletes the entries created before cutoff, oldest first and PruneChunk
// at a time, and returns how many it deleted.
//
// The entries that remain keep their hashes. The oldest of them now chains
// from an entry that is gone, so each chunk appends a SECURITY entry naming the
// hash it chains from, and VerifyChain accepts a log that starts there only
// when that entry exists: old entries deleted any other way still fail
// verification. Each chunk is one transaction holding the chain lock, so the
// log verifies after every commit and a prune that stops part way keeps what
// it did.
//
// Pruning used to rewrite the oldest remaining entry's hash, which broke the
// link from the entry after it, so every prune that deleted anything left a
// log that failed verification.
func Prune(ctx context.Context, db *gorm.DB, cutoff time.Time) (int64, error) {
	var total int64
	for {
		n, err := pruneChunk(db.WithContext(ctx), cutoff)
		total += n
		if err != nil || n < PruneChunk {
			return total, err
		}
	}
}

func pruneChunk(db *gorm.DB, cutoff time.Time) (int64, error) {
	var removed int64
	err := db.Transaction(func(tx *gorm.DB) error {
		if _, err := lockAndReadHead(tx); err != nil {
			return err
		}
		// The chunk ends at its last entry, so it is deleted as an index range
		// rather than a list of 5,000 ids.
		var last models.ActivityLog
		if err := tx.Where("created_at < ?", cutoff).
			Order("created_at asc, id asc").
			Offset(PruneChunk - 1).Limit(1).
			Find(&last).Error; err != nil {
			return fmt.Errorf("finding the end of the chunk: %w", err)
		}
		del := tx.Where("created_at < ?", cutoff)
		if last.ID != "" {
			del = del.Where("created_at < ? OR (created_at = ? AND id <= ?)", last.CreatedAt, last.CreatedAt, last.ID)
		}
		res := del.Delete(&models.ActivityLog{})
		if res.Error != nil {
			return fmt.Errorf("deleting entries: %w", res.Error)
		}
		if res.RowsAffected == 0 {
			return nil
		}
		var first models.ActivityLog
		if err := tx.Order("created_at asc, id asc").Limit(1).Find(&first).Error; err != nil {
			return fmt.Errorf("reading the oldest remaining entry: %w", err)
		}
		removed = res.RowsAffected
		return appendBatch(tx, []models.ActivityLog{{
			Method:      "SECURITY",
			Path:        prunePath,
			Status:      200,
			Resource:    "activity_logs",
			ResourceIDs: first.PrevHash,
			RecordCount: int(removed),
		}})
	})
	if err != nil {
		return 0, err
	}
	return removed, nil
}

// pruneRecorded reports whether a prune recorded that the log starts after the
// entry whose hash is prevHash.
func pruneRecorded(ctx context.Context, db *gorm.DB, prevHash string) (bool, error) {
	var n int64
	err := db.WithContext(ctx).Model(&models.ActivityLog{}).
		Where("method = ? AND path = ? AND resource_ids = ?", "SECURITY", prunePath, prevHash).
		Count(&n).Error
	return n > 0, err
}

// ErrChainIntact is returned by Reseal when the chain verifies.
var ErrChainIntact = errors.New("the chain verifies: there is nothing to reseal")

// ErrNotTheBreak is returned by Reseal when the entry named is not the first
// one that fails verification.
var ErrNotTheBreak = errors.New("that is not the first entry that fails verification")

// Reseal recomputes the chain from its first bad entry, and records that it did.
//
// For a log that fails through no one's tampering: before v3.215.0 every entry
// was hashed with a timestamp finer than Postgres and MySQL store, so those
// chains fail on their first row and always will. A reseal trusts the rows as
// they stand now, which is exactly what makes it dangerous, so it is never
// automatic. The caller names the first bad entry, which is checked against a
// fresh verification, and the reseal appends a SECURITY entry naming who did
// it, from which entry, how many it covered, and a digest of every hash it
// replaced. Changing that entry afterwards breaks the chain like any other.
func Reseal(ctx context.Context, db *gorm.DB, fromID, userID, ip, userAgent string) (int, error) {
	status, err := VerifyChain(ctx, db)
	if err != nil {
		return 0, err
	}
	if status.Valid {
		return 0, ErrChainIntact
	}
	if status.BrokenAtID != fromID {
		return 0, fmt.Errorf("%w: the chain first fails at %s", ErrNotTheBreak, status.BrokenAtID)
	}

	resealed := 0
	err = db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if _, err := lockAndReadHead(tx); err != nil {
			return err
		}
		var from models.ActivityLog
		if err := tx.First(&from, "id = ?", fromID).Error; err != nil {
			return fmt.Errorf("loading %s: %w", fromID, err)
		}
		// The entry before the break is the last one that verified.
		var before models.ActivityLog
		if err := tx.Where("(created_at, id) < (?, ?)", from.CreatedAt, from.ID).
			Order("created_at desc, id desc").Limit(1).Find(&before).Error; err != nil {
			return fmt.Errorf("loading the entry before %s: %w", fromID, err)
		}

		prevHash := before.Hash
		replaced := sha256.New()
		cond, args := "(created_at, id) >= (?, ?)", []interface{}{from.CreatedAt, from.ID}
		for {
			var batch []models.ActivityLog
			if err := tx.Where(cond, args...).Order("created_at asc, id asc").
				Limit(verifyBatchSize).Find(&batch).Error; err != nil {
				return err
			}
			for i := range batch {
				e := &batch[i]
				canonical, err := Canonical(e)
				if err != nil {
					return err
				}
				hash := ComputeHash(prevHash, canonical)
				replaced.Write([]byte(e.Hash))
				if err := tx.Model(&models.ActivityLog{}).Where("id = ?", e.ID).
					Updates(map[string]interface{}{"prev_hash": prevHash, "hash": hash}).Error; err != nil {
					return fmt.Errorf("resealing %s: %w", e.ID, err)
				}
				prevHash = hash
				resealed++
			}
			if len(batch) < verifyBatchSize {
				break
			}
			last := batch[len(batch)-1]
			cond, args = "(created_at, id) > (?, ?)", []interface{}{last.CreatedAt, last.ID}
		}

		return appendBatch(tx, []models.ActivityLog{{
			UserID:        userID,
			Method:        "SECURITY",
			Path:          "audit.chain.resealed",
			Status:        200,
			PayloadDigest: hex.EncodeToString(replaced.Sum(nil)),
			IPAddress:     ip,
			UserAgent:     userAgent,
			Resource:      "activity_logs",
			ResourceIDs:   fromID,
			RecordCount:   resealed,
		}})
	})
	if err != nil {
		return 0, err
	}
	return resealed, nil
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
			if total+i == 0 && e.PrevHash != "" {
				// The oldest entry chains from one that is gone, which is what a
				// prune leaves. It verifies only if the prune recorded that hash.
				recorded, err := pruneRecorded(ctx, db, e.PrevHash)
				if err != nil {
					return ChainStatus{TotalEntries: total}, err
				}
				if !recorded {
					return ChainStatus{
						Valid:      false,
						BrokenAtID: e.ID,
						Got:        e.PrevHash,
						Message:    "the log starts after entries that were deleted, and no prune recorded deleting them",
					}, nil
				}
				prevHash = e.PrevHash
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
	Provider  string ` + "`" + `gorm:"size:50;not null;uniqueIndex:idx_webhook_provider_external,priority:1" json:"provider"` + "`" + `
	EventType string ` + "`" + `gorm:"size:100;index" json:"event_type"` + "`" + `

	// ExternalID is the provider's own event id, and the second half of the
	// idempotency key. A pointer because it has to be NULL when a provider
	// does not supply one: every database allows repeated NULLs in a unique
	// index and none allow a repeated empty string, so storing "" would make
	// two unrelated anonymous events collide and silently drop the second.
	ExternalID *string ` + "`" + `gorm:"size:255;uniqueIndex:idx_webhook_provider_external,priority:2" json:"external_id,omitempty"` + "`" + `
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

// The unique index is declared on the fields above rather than built here.
//
// It used to live in a method returning DDL as a string, which nothing
// called. The column had a plain index, the INSERT never failed, and the
// handler's "duplicate means already processed" branch was unreachable: every
// retried delivery ran the handler again. Declaring it as a tag means the
// migration creates it, and the constraint is where it can be seen.
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

func apiWebhooksHandlerGo() string { return tmpl("api/handlers/webhooks.go") }

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

	// Attributes restrict the flag to subjects whose attributes match:
	// {"business_unit": ["eu", "ke"]} is on only for a subject whose
	// business_unit is eu or ke. Every key listed must match, ignoring case.
	// A subject without the attribute does not match, so a rule on something
	// the app never supplies fails closed. See flags.AttributesFor.
	Attributes map[string][]string ` + "`" + `json:"attributes,omitempty"` + "`" + `
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
	"strings"
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
	db    *gorm.DB
	hub   *realtime.Hub // optional — when set, broadcasts on Refresh
	mu    sync.RWMutex
	flags map[string]*models.FeatureFlag
	stop  chan struct{}
	// exposures feeds one writer. Each flag check used to start a goroutine and
	// an INSERT of its own, which a busy page turned into thousands a second.
	exposures chan models.FlagExposure
}

// New returns an Engine with the cache pre-warmed. Call from
// routes.Setup. hub is optional — pass nil to disable broadcasts.
func New(db *gorm.DB, hub *realtime.Hub) *Engine {
	e := &Engine{
		db:        db,
		hub:       hub,
		flags:     make(map[string]*models.FeatureFlag),
		stop:      make(chan struct{}),
		exposures: make(chan models.FlagExposure, exposureQueueSize),
	}
	if err := e.Refresh(); err != nil {
		log.Printf("[flags] initial refresh failed: %v", err)
	}
	go e.refreshLoop()
	go e.writeExposures()
	setDefault(e)
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

// Subject is who a flag is checked for: a user, and whatever attributes
// the rules can target (a business unit, a region, a plan).
type Subject struct {
	UserID     string
	Attributes map[string]string
}

// AttributesFor supplies the attributes of the user making a request, for
// rules that target them. The default knows the role the auth middleware
// set. Replace it at boot to target anything else your users carry:
//
//	flags.AttributesFor = func(c *gin.Context) map[string]string {
//	    return map[string]string{"business_unit": businessUnitOf(c)}
//	}
var AttributesFor = func(c *gin.Context) map[string]string {
	attrs := map[string]string{}
	if c == nil {
		return attrs
	}
	if v, ok := c.Get("user_role"); ok {
		if role, ok := v.(string); ok && role != "" {
			attrs["role"] = role
		}
	}
	return attrs
}

func subjectFrom(c *gin.Context) Subject {
	return Subject{UserID: userIDFrom(c), Attributes: AttributesFor(c)}
}

// IsEnabled returns true when the flag is on for the current user.
// Always returns false for unknown flags (fail closed).
func (e *Engine) IsEnabled(c *gin.Context, name string) bool {
	return e.evaluate(subjectFrom(c), name) == "enabled"
}

// Variant returns the assigned variant for an A/B flag. For boolean
// flags, returns "enabled" or "disabled". For unknown flags, returns
// the empty string.
func (e *Engine) Variant(c *gin.Context, name string) string {
	return e.evaluate(subjectFrom(c), name)
}

// IsEnabledForUser is the explicit form for backend code that has the
// user_id directly (e.g. cron jobs operating on a specific user). It has
// no attributes, so a flag with attribute rules is off for it: use
// IsEnabledFor with a Subject for those.
func (e *Engine) IsEnabledForUser(userID, name string) bool {
	return e.evaluate(Subject{UserID: userID}, name) == "enabled"
}

// VariantForUser is the explicit form of Variant.
func (e *Engine) VariantForUser(userID, name string) string {
	return e.evaluate(Subject{UserID: userID}, name)
}

// IsEnabledFor checks a flag for a subject built by the caller: a job that
// knows the user and their business unit, say.
func (e *Engine) IsEnabledFor(s Subject, name string) bool {
	return e.evaluate(s, name) == "enabled"
}

// VariantFor is the explicit form of Variant for a subject.
func (e *Engine) VariantFor(s Subject, name string) string {
	return e.evaluate(s, name)
}

// ── Package-level ────────────────────────────────────────────────────────
//
// The engine routes.Setup starts, reachable from any handler, service or job
// without threading it through. Before v3.223.0 the package comment showed
// flags.IsEnabled(c, ...) and no such function existed: the only engine was
// a local variable in routes.Setup, so application code could not check a
// flag at all.

var (
	defaultMu     sync.RWMutex
	defaultEngine *Engine
)

func setDefault(e *Engine) {
	defaultMu.Lock()
	defaultEngine = e
	defaultMu.Unlock()
}

func current() *Engine {
	defaultMu.RLock()
	defer defaultMu.RUnlock()
	return defaultEngine
}

// IsEnabled reports whether a flag is on for the request's user. False
// before the engine has started and for a flag that does not exist: a check
// that cannot be answered fails closed.
func IsEnabled(c *gin.Context, name string) bool {
	if e := current(); e != nil {
		return e.IsEnabled(c, name)
	}
	return false
}

// Variant returns the request user's variant, or "" when it cannot be answered.
func Variant(c *gin.Context, name string) string {
	if e := current(); e != nil {
		return e.Variant(c, name)
	}
	return ""
}

// IsEnabledForUser is IsEnabled for code with a user ID and no request.
func IsEnabledForUser(userID, name string) bool {
	if e := current(); e != nil {
		return e.IsEnabledForUser(userID, name)
	}
	return false
}

// VariantForUser is Variant for code with a user ID and no request.
func VariantForUser(userID, name string) string {
	if e := current(); e != nil {
		return e.VariantForUser(userID, name)
	}
	return ""
}

// IsEnabledFor is IsEnabled for a subject built by the caller.
func IsEnabledFor(s Subject, name string) bool {
	if e := current(); e != nil {
		return e.IsEnabledFor(s, name)
	}
	return false
}

// VariantFor is Variant for a subject built by the caller.
func VariantFor(s Subject, name string) string {
	if e := current(); e != nil {
		return e.VariantFor(s, name)
	}
	return ""
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
func (e *Engine) evaluate(s Subject, name string) string {
	userID := s.UserID
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

	// Attributes restrict the flag to matching subjects, every key listed.
	for key, allowed := range rules.Attributes {
		if !attributeMatches(s.Attributes[key], allowed) {
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

// exposureQueueSize bounds the exposures waiting for the writer. Past it they
// are dropped: an exposure is analytics, never worth slowing a request for.
const exposureQueueSize = 4096

// trackExposure hands the flag check to the exposure writer without blocking.
func (e *Engine) trackExposure(flagID, flagName, userID, variant string) {
	if userID == "" {
		// Anonymous exposures pollute the table without buying us
		// anything (we can't link them to a user later). Skip.
		return
	}
	select {
	case e.exposures <- models.FlagExposure{FlagID: flagID, FlagName: flagName, UserID: userID, Variant: variant}:
	default:
	}
}

// writeExposures writes exposures in batches, and records each user, flag and
// variant once a UTC day: an exposure says who saw which variant, not how many
// times they reloaded the page.
func (e *Engine) writeExposures() {
	seen := map[string]bool{}
	day := ""
	for {
		var first models.FlagExposure
		select {
		case <-e.stop:
			return
		case first = <-e.exposures:
		}
		batch := []models.FlagExposure{first}
	drain:
		for len(batch) < 256 {
			select {
			case x := <-e.exposures:
				batch = append(batch, x)
			default:
				break drain
			}
		}

		today := time.Now().UTC().Format("2006-01-02")
		if today != day || len(seen) > 100000 {
			seen, day = map[string]bool{}, today
		}
		fresh := batch[:0]
		for _, x := range batch {
			k := x.UserID + "|" + x.FlagID + "|" + x.Variant
			if seen[k] {
				continue
			}
			seen[k] = true
			fresh = append(fresh, x)
		}
		if len(fresh) == 0 {
			continue
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		if err := e.db.WithContext(ctx).CreateInBatches(fresh, 256).Error; err != nil {
			log.Printf("[flags] recording %d exposures: %v", len(fresh), err)
		}
		cancel()
	}
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
	if c == nil {
		return ""
	}
	if v, ok := c.Get("user_id"); ok {
		s, _ := v.(string)
		return s
	}
	return ""
}

// attributeMatches reports whether a subject's value is one of a rule's.
// Case-insensitive, because these are typed in by hand. An empty value
// matches nothing.
func attributeMatches(value string, allowed []string) bool {
	if value == "" {
		return false
	}
	for _, a := range allowed {
		if strings.EqualFold(a, value) {
			return true
		}
	}
	return false
}
`
}

func apiFlagsHandlerGo() string { return tmpl("api/handlers/flags.go") }

func apiActivityMiddlewareGo() string {
	return `package middleware

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"
	"strings"
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
	// One chain writer per process, shared with the security-event log.
	audit.Start(db)
	return func(c *gin.Context) {
		switch c.Request.Method {
		case http.MethodGet, http.MethodHead, http.MethodOptions:
			recordRead(c)
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
			// Which record the write touched, so the log answers "who changed
			// this record" as well as "who read it". Empty on routes without one.
			ResourceIDs: c.Param("id"),
		}
		// Non-blocking. The writer stamps created_at and chains the entry;
		// with a full backlog it is dropped rather than stalling the request.
		audit.Enqueue(entry)
	}
}

// recordRead runs a read and, when the handler marked what it served, records it.
//
// Only marked reads: handlers generated with --audit-reads mark the rows they
// returned, and nothing else is recorded, because reads are most of all
// traffic and logging every page load would bury the writes.
func recordRead(c *gin.Context) {
	started := time.Now()
	c.Next()

	mark, ok := audit.ReadMarkOf(c)
	if !ok || c.Writer.Status() < 200 || c.Writer.Status() >= 300 {
		return
	}
	userID, _ := c.Get("user_id")
	uid, _ := userID.(string)
	audit.Enqueue(models.ActivityLog{
		UserID: uid,
		Method: c.Request.Method,
		Path:   c.FullPath(),
		Status: c.Writer.Status(),
		// The query can hold what was searched for, a patient's name typed
		// into a search box, so only its digest is kept.
		PayloadDigest: digestBody([]byte(c.Request.URL.RawQuery)),
		IPAddress:     resolveClientIP(c),
		UserAgent:     c.Request.UserAgent(),
		DurationMS:    time.Since(started).Milliseconds(),
		Resource:      mark.Resource,
		ResourceIDs:   strings.Join(mark.IDs, ","),
		RecordCount:   mark.Count,
	})
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

// AuditDroppedCount returns the number of audit entries dropped because the
// writer's backlog was full. Read it from a /healthz or admin endpoint to spot
// sustained back-pressure.
func AuditDroppedCount() uint64 {
	return audit.Dropped()
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

func apiActivityHandlerGo() string { return tmpl("api/handlers/activity.go") }

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
	"errors"
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

// XLSX writes items as an Excel workbook into w. An export that reads its rows
// in batches should use NewXLSXStream instead, so the rows are never all in
// memory at once.
func XLSX(w io.Writer, items interface{}, opts Options) error {
	sheet, err := NewXLSXStream(opts)
	if err != nil {
		return err
	}
	if err := sheet.Rows(items); err != nil {
		return errors.Join(err, sheet.Close())
	}
	return errors.Join(sheet.Finish(w), sheet.Close())
}

// ErrTooManyRows is returned when an export has more rows than a worksheet can
// hold: 1,048,576, the header included. CSV has no such limit.
var ErrTooManyRows = errors.New("export: more rows than an XLSX worksheet holds; export as CSV instead")

// XLSXStream builds a workbook one batch of rows at a time.
//
// Rows go through excelize's stream writer, which moves them to a temporary
// file once they pass a few megabytes, so memory stays flat however many there
// are. Building the sheet cell by cell held every row in memory: 300,000
// contacts took an API from 80 MB to 1.3 GB.
type XLSXStream struct {
	file   *excelize.File
	stream *excelize.StreamWriter
	opts   Options
	next   int // the worksheet row the next item goes on
	cells  []interface{}
}

// NewXLSXStream starts a workbook and writes its header row. Close it when done,
// which removes its temporary files.
func NewXLSXStream(opts Options) (*XLSXStream, error) {
	f := excelize.NewFile()
	sheet := opts.Sheet
	if sheet == "" {
		sheet = "Sheet1"
	}
	if sheet != "Sheet1" {
		// excelize creates "Sheet1" by default; swap to the requested name.
		if err := f.SetSheetName("Sheet1", sheet); err != nil {
			return nil, errors.Join(err, f.Close())
		}
	}
	stream, err := f.NewStreamWriter(sheet)
	if err != nil {
		return nil, errors.Join(err, f.Close())
	}
	header := make([]interface{}, len(opts.Columns))
	for i, col := range opts.Columns {
		header[i] = col.Header
	}
	if err := stream.SetRow("A1", header); err != nil {
		return nil, errors.Join(err, f.Close())
	}
	return &XLSXStream{file: f, stream: stream, opts: opts, next: 2, cells: make([]interface{}, len(opts.Columns))}, nil
}

// Rows appends items, a slice of structs, below the rows already written.
func (x *XLSXStream) Rows(items interface{}) error {
	v := reflect.ValueOf(items)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Slice {
		return fmt.Errorf("export: items must be a slice, got %T", items)
	}
	for i := 0; i < v.Len(); i++ {
		if x.next > excelize.TotalRows {
			return ErrTooManyRows
		}
		for j, col := range x.opts.Columns {
			x.cells[j] = formatCell(extractField(v.Index(i), col.Field), col.Format)
		}
		cell, err := excelize.CoordinatesToCellName(1, x.next)
		if err != nil {
			return err
		}
		if err := x.stream.SetRow(cell, x.cells); err != nil {
			return err
		}
		x.next++
	}
	return nil
}

// Written is the number of rows added so far, the header not included.
func (x *XLSXStream) Written() int {
	return x.next - 2
}

// Finish completes the sheet and writes the workbook to w.
func (x *XLSXStream) Finish(w io.Writer) error {
	if err := x.stream.Flush(); err != nil {
		return err
	}
	return x.file.Write(w)
}

// Close removes the workbook's temporary files.
func (x *XLSXStream) Close() error {
	return x.file.Close()
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
	// Indexed by byte rather than ranged by rune: intPart is digits from
	// FormatFloat, so the two are the same here, and mixing a rune loop with
	// byte arithmetic is how the first multi-byte character someone passes in
	// gets truncated.
	for i := 0; i < len(intPart); i++ {
		if i > 0 && (len(intPart)-i)%3 == 0 {
			out = append(out, ',')
		}
		out = append(out, intPart[i])
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
	"errors"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Error is the wire shape of every error envelope.
type Error struct {
	Code    string            ` + "`" + `json:"code"` + "`" + `
	Message string            ` + "`" + `json:"message"` + "`" + `
	Details map[string]string ` + "`" + `json:"details,omitempty"` + "`" + `
}

// The named helpers, each one line over Fail.
//
// They used to carry a status and a code side by side, hardcoded here, while
// codes.go carried the same pairs generated from the catalogue. Two tables for
// one fact is how VALIDATION_ERROR came back as 422 from one helper and 400
// from a handler that wrote its own envelope. There is one table now, in
// codes.go, and everything below reads the status out of it.

// BadRequest is 400: a malformed request the client cannot fix without
// changing what it sent.
func BadRequest(c *gin.Context, message string) {
	Fail(c, CodeBadRequest, message)
}

// Unauthorized is 401: missing or invalid credentials.
func Unauthorized(c *gin.Context, message string) {
	if message == "" {
		message = "Authentication required"
	}
	Fail(c, CodeUnauthorized, message)
}

// Forbidden is 403: authenticated, but not allowed.
func Forbidden(c *gin.Context, message string) {
	if message == "" {
		message = "You don't have permission to do that"
	}
	Fail(c, CodeForbidden, message)
}

// NotFound is 404: the entity did not exist, or access rules filtered it out.
func NotFound(c *gin.Context, message string) {
	if message == "" {
		message = "Not found"
	}
	Fail(c, CodeNotFound, message)
}

// Conflict is 409: a unique constraint, or a version conflict.
func Conflict(c *gin.Context, message string) {
	Fail(c, CodeConflict, message)
}

// Validation is 422: the payload was well formed and failed validation. Pass
// per-field messages so the frontend can highlight the fields.
func Validation(c *gin.Context, message string, fields map[string]string) {
	Fail(c, CodeValidationError, message, fields)
}

// RuleError is a business rule the caller broke, phrased for the caller.
//
// Errors coming back from a write are mostly not for the client: a driver
// failure or a constraint violation carries schema details and sometimes SQL,
// so the handler cannot simply echo what it gets. This is how application code
// says "this one is different, the message is the point".
//
// Return it from a GORM hook, a callback or a service method:
//
//	func (e *JournalEntry) BeforeCreate(tx *gorm.DB) error {
//	    if !balanced(e.Lines) {
//	        return respond.Rule("debits and credits do not balance")
//	    }
//	    return nil
//	}
//
// and the caller gets 422 with that sentence instead of an opaque 500.
type RuleError struct{ Message string }

func (e *RuleError) Error() string { return e.Message }

// Rule builds a RuleError. Takes a format string because most rules want to
// quote the values that broke them.
func Rule(format string, args ...interface{}) error {
	return &RuleError{Message: fmt.Sprintf(format, args...)}
}

// IsRule reports whether err is, or wraps, a RuleError.
func IsRule(err error) (*RuleError, bool) {
	var rule *RuleError
	if errors.As(err, &rule) {
		return rule, true
	}
	return nil, false
}

// Coded is an error that knows what it should look like on the wire.
//
// Implement it on a sentinel error when the caller needs to act on it, rather
// than leaving it to become an opaque 500:
//
//	var ErrNoSeatsLeft = seatsError{}
//
//	type seatsError struct{}
//	func (seatsError) Error() string          { return "no seats left on this plan" }
//	func (seatsError) ErrorCode() respond.Code { return respond.CodeConflict }
//
// Anything returned from a service or a GORM hook that implements this is
// answered with the code's documented status. This exists because
// tenant.ErrNoOrganization, which means "say which organization you are acting
// in", arrived as a 500 saying "Failed to fetch deals": respond cannot import
// the tenant package, and an error that can describe itself does not need it to.
type Coded interface {
	error
	ErrorCode() Code
}

// FieldErrors is an error about particular fields, such as a value refused by
// its column's format in internal/fieldtypes. It is answered with 422 and the
// per-field messages in details, so a form can put each one under its input.
type FieldErrors interface {
	error
	FieldErrors() map[string]string
}

// WriteError picks the right response for an error returned by a write.
//
// An error that carries its own code is answered with that code's status. A rule
// the caller broke becomes 422 with its message. A missing row becomes 404. A
// value a unique column already holds becomes 409, naming the field when the
// database says which. Everything else is logged and comes back as an opaque 500, which is what it was
// before, minus the part where the error vanished entirely.
func WriteError(c *gin.Context, err error, fallback string) {
	var fields FieldErrors
	if errors.As(err, &fields) {
		Validation(c, fields.Error(), fields.FieldErrors())
		return
	}
	var coded Coded
	if errors.As(err, &coded) {
		Fail(c, coded.ErrorCode(), coded.Error())
		return
	}
	if rule, ok := IsRule(err); ok {
		Validation(c, rule.Message, nil)
		return
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		NotFound(c, "")
		return
	}
	if field, ok := DuplicateKey(c, err); ok {
		duplicate(c, field)
		return
	}
	ServerError(c, "INTERNAL_ERROR", err, fallback)
}

// DuplicateKey reports whether err is a unique-constraint violation, and the
// column it was on when the driver's message says so.
//
// GORM has gorm.ErrDuplicatedKey, but it only translates the driver's error
// into it when the connection was opened with TranslateError, so the four
// drivers Grit runs on are also matched by message. Without this a second
// contact with the same code came back as a 500 "Failed to create contact",
// which tells the person filling the form nothing they can fix.
func DuplicateKey(c *gin.Context, err error) (string, bool) {
	if err == nil {
		return "", false
	}
	msg := err.Error()
	lower := strings.ToLower(msg)
	found := errors.Is(err, gorm.ErrDuplicatedKey)
	for _, marker := range []string{
		"duplicate key value",      // Postgres
		"unique constraint failed", // SQLite
		"duplicate entry",          // MySQL
		"violation of unique key",  // SQL Server
	} {
		if strings.Contains(lower, marker) {
			found = true
		}
	}
	if !found {
		return "", false
	}
	return duplicateColumn(c, msg), true
}

// sqliteUnique is "UNIQUE constraint failed: contacts.code"; constraintName
// is the quoted index name Postgres, MySQL and SQL Server report.
var (
	sqliteUnique   = regexp.MustCompile(` + "`" + `(?i)unique constraint failed: ([\w.]+)` + "`" + `)
	constraintName = regexp.MustCompile(` + "`" + `(?i)(?:constraint|key|index) ['"]([\w.]+)['"]` + "`" + `)
)

// duplicateColumn names the column when it can be told, and "" when not.
//
// SQLite names it outright. The others name the index, which GORM calls
// idx_<table>_<column> or uni_<table>_<column>; the table is taken from the
// route (/api/contacts/:id is contacts), so the column is what is left. When
// that does not line up the column is not guessed: the 409 still stands, just
// without saying which field.
func duplicateColumn(c *gin.Context, msg string) string {
	if m := sqliteUnique.FindStringSubmatch(msg); m != nil {
		col := m[1]
		if i := strings.LastIndex(col, "."); i >= 0 {
			col = col[i+1:]
		}
		return col
	}
	m := constraintName.FindStringSubmatch(msg)
	if m == nil {
		return ""
	}
	name := m[1]
	if i := strings.LastIndex(name, "."); i >= 0 {
		name = name[i+1:] // MySQL 8: 'contacts.idx_contacts_code'
	}
	table := routeTable(c)
	for _, prefix := range []string{"idx_", "uni_"} {
		if table != "" && strings.HasPrefix(name, prefix+table+"_") {
			return strings.TrimPrefix(name, prefix+table+"_")
		}
	}
	return ""
}

// routeTable is the resource segment of the matched route: the last segment
// that is not a parameter, with dashes as underscores.
func routeTable(c *gin.Context) string {
	if c == nil {
		return ""
	}
	parts := strings.Split(strings.Trim(c.FullPath(), "/"), "/")
	for i := len(parts) - 1; i >= 0; i-- {
		if parts[i] != "" && !strings.HasPrefix(parts[i], ":") && !strings.HasPrefix(parts[i], "*") {
			return strings.ReplaceAll(parts[i], "-", "_")
		}
	}
	return ""
}

// duplicate answers 409, naming the field when it is known so a form can put
// the message under that input.
func duplicate(c *gin.Context, field string) {
	if field == "" {
		Fail(c, CodeConflict, "A record with that value already exists")
		return
	}
	label := strings.ReplaceAll(field, "_", " ")
	Fail(c, CodeConflict, "That "+label+" is already taken", map[string]string{field: "This " + label + " is already taken"})
}

// Internal answers 500 with a generic message. The error is logged with the
// request it failed rather than sent: its text is for the operator, and to a
// client it can describe the schema. It used to be discarded, so a 500 from
// here left no trace anywhere.
func Internal(c *gin.Context, internalErr error) {
	ServerError(c, "INTERNAL_ERROR", internalErr, "Internal server error")
}

// ServerError answers 500 with code and a message safe for the client, and
// logs err with the method, path and request id, so the cause reaches the
// operator and not the caller. The request id is the X-Request-ID header the
// client received, which is how a report of a failure finds its log line.
func ServerError(c *gin.Context, code string, err error, message string) {
	if err != nil {
		_ = c.Error(err)
	}
	log.Printf("[500] %s %s | id=%s | %s: %v", c.Request.Method, c.Request.URL.Path, c.GetString("request_id"), code, err)
	// Not Fail: code is a plain string here, so callers can name the subsystem
	// that broke (DB_ERROR, STORAGE_ERROR) without every one of them being a
	// catalogued constant. The status is 500 either way.
	c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": Error{Code: code, Message: message}})
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
	// Indexed by byte rather than ranged by rune: intPart is digits from
	// FormatFloat, so the two are the same here, and mixing a rune loop with
	// byte arithmetic is how the first multi-byte character someone passes in
	// gets truncated.
	for i := 0; i < len(intPart); i++ {
		if i > 0 && (len(intPart)-i)%3 == 0 {
			out = append(out, ',')
		}
		out = append(out, intPart[i])
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
	"context"
	"encoding/json"
	"log"
	"sync"
	"time"

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

	// ExpiresAt is when the JWT that authorised this connection stops being
	// valid. The handshake is the only time the token is checked, so without
	// this the socket would outlive its own credential and keep streaming to
	// a session that has since been signed out. writePump closes the
	// connection at this instant and the client reconnects with a fresh
	// token, which is the same bound REST already operates under.
` + hubClientEventsFieldNew + `

// disconnectGrace is how long DisconnectUser waits for writePump to send its
// close frame and tear the socket down before forcing it. It matches the write
// deadline writePump already works to.
const disconnectGrace = 10 * time.Second

// Hub manages connected clients. Safe for concurrent use.
type Hub struct {
	mu      sync.RWMutex
` + hubChannelsFieldNew + hubPresenceFieldAdd + hubStatsFieldAdd + `
	// nodeID identifies this process so it can ignore its own messages coming
	// back off the backplane.
	nodeID string

	// backplane is nil for a single-process deployment, which is the default
	// and is correct for most projects. See WithBackplane.
	backplane Backplane
	pub       chan []byte
	cancel    context.CancelFunc
}

// Option configures a Hub at construction.
type Option func(*Hub)

// WithBackplane fans every send out to the other API processes.
//
// Without one the Hub is an in-process registry, so a second replica silently
// halves realtime: a user connected to replica A never sees an event published
// on replica B, the push succeeds into a registry that does not contain them,
// and nothing errors. Pass this whenever more than one process serves
// websockets, including during a rolling deploy where two versions overlap.
//
// Delivery between nodes is best effort. See Backplane.
func WithBackplane(b Backplane) Option {
	return func(h *Hub) { h.backplane = b }
}

// WithRedis is WithBackplane over the Redis this project already runs.
//
// An empty url is a no-op, so the common wiring is one unconditional line:
//
//	realtime.NewHub(realtime.WithRedis(cfg.RedisURL, ""))
//
// A project with no Redis configured then keeps the single-process behaviour
// instead of failing to start.
func WithRedis(redisURL, channel string) Option {
	return func(h *Hub) {
		if redisURL == "" {
			return
		}
		bp, err := RedisBackplane(redisURL, channel)
		if err != nil {
			// Not fatal. Realtime degrades to this process's own clients,
			// which is exactly what every project had before backplanes
			// existed, and the API still serves everything else.
			log.Printf("[realtime] no backplane, staying single-process: %v", err)
			return
		}
		h.backplane = bp
	}
}

// NewHub returns a Hub, single-process unless an option says otherwise.
func NewHub(opts ...Option) *Hub {
	h := &Hub{
		clients: make(map[string]map[*Client]struct{}),
		nodeID:  newNodeID(),
	}
	for _, opt := range opts {
		opt(h)
	}
	if h.backplane != nil {
		h.pub = make(chan []byte, publishBuffer)
		ctx, cancel := context.WithCancel(context.Background())
		h.cancel = cancel
		go h.backplane.Subscribe(ctx, h.receive)
` + hubPresenceStartNew + `	}
	return h
}

// Close stops the backplane subscription. Local clients are unaffected.
func (h *Hub) Close() error {
	if h.cancel != nil {
		h.cancel()
	}
	if h.backplane != nil {
		return h.backplane.Close()
	}
	return nil
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
` + hubPresenceUnregisterNew + `	if set, ok := h.clients[c.UserID]; ok {
		if _, exists := set[c]; exists {
` + hubUnregisterNew + `		}
		if len(set) == 0 {
			delete(h.clients, c.UserID)
		}
	}
}

// DisconnectUser closes every connection a user holds, immediately.
//
// Call this when a session is revoked, a password changes, or an account is
// deactivated. Waiting for the token to expire leaves a revoked device
// receiving live data for the rest of the access-token lifetime, which is not
// what "sign out of all devices" tells the user happened.
func (h *Hub) DisconnectUser(userID string) {
	h.disconnectLocal(userID)
	// And on every other node. A revoked connection on a replica that did not
	// handle the revocation request would otherwise stay open and keep
	// receiving, which is the failure this function exists to prevent.
	h.publish(fanout{Kick: userID})
}

// disconnectLocal closes this process's connections for a user. It is what
// both DisconnectUser and an incoming kick run.
func (h *Hub) disconnectLocal(userID string) {
	h.mu.Lock()
	set, ok := h.clients[userID]
	if !ok {
		h.mu.Unlock()
		return
	}
	conns := make([]*websocket.Conn, 0, len(set))
	for c := range set {
` + hubDisconnectNew + `		// Closing Send makes writePump emit a proper close frame and tear the
		// connection down itself, so the client sees a clean 1000 rather than
		// an abnormal 1006 and can distinguish "you were signed out" from
		// "the network dropped".
		close(c.Send)
	}
	delete(h.clients, userID)
	h.mu.Unlock()

	// Backstop, outside the lock. If writePump is wedged on a dead socket it
	// will not get to its own deferred Close, and a revoked connection that
	// stays open is the thing this function exists to prevent. wsWriteWait is
	// the bound writePump is already operating under.
	go func() {
		time.Sleep(disconnectGrace)
` + hubBackstopNew + `	}()
}

// SendToUser delivers an event to every connection bound to userID.
// If a connection's send buffer is full the message is dropped for that
// connection only — we never block the entire hub on a slow client.
// The slow client will resync on its next REST poll/refetch.
func (h *Hub) SendToUser(userID string, evt Event) {
	h.SendToUsers([]string{userID}, evt)
}

` + hubLocalSendsCounted + `
// SendToUsers fans out to a slice of user IDs.
func (h *Hub) SendToUsers(userIDs []string, evt Event) {
	if len(userIDs) == 0 {
		return
	}
	bytes, err := json.Marshal(evt)
	if err != nil {
		log.Printf("[realtime] marshal: %v", err)
		return
	}
	// This node first, and unconditionally: local clients must not wait on
	// Redis, and must still be served when Redis is down.
	h.deliverLocal(userIDs, bytes)
	// Then once for every other node, whatever the size of the audience. A
	// per-user publish would send the same payload N times over the wire.
	h.publish(fanout{Users: userIDs, Event: bytes})
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
	h.broadcastLocal(bytes)
	h.publish(fanout{All: true, Event: bytes})
}
`
}

func apiRealtimeHandlerGo() string { return tmpl("api/handlers/realtime.go") }

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
	"github.com/MUKE-coder/sentinel/v2/redisstore"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"` + "{{MODULE}}" + `/internal/ai"
	"` + "{{MODULE}}" + `/internal/authz"
	"` + "{{MODULE}}" + `/internal/cache"
	"` + "{{MODULE}}" + `/internal/config"
	"` + "{{MODULE}}" + `/internal/database"
	"` + "{{MODULE}}" + `/internal/events"
	"` + "{{MODULE}}" + `/internal/health"
	"` + "{{MODULE}}" + `/internal/llms"
	"` + "{{MODULE}}" + `/internal/saga"
	"` + "{{MODULE}}" + `/internal/sagas"
	"` + "{{MODULE}}" + `/internal/paginate"
	"` + "{{MODULE}}" + `/internal/handlers"
	"` + "{{MODULE}}" + `/internal/settings"
	"` + "{{MODULE}}" + `/internal/tracing"
	"` + "{{MODULE}}" + `/internal/mail"
	"` + "{{MODULE}}" + `/internal/middleware"
	"` + "{{MODULE}}" + `/internal/models"
	"` + "{{MODULE}}" + `/internal/jobs"
	"` + "{{MODULE}}" + `/internal/realtime"
	"` + "{{MODULE}}" + `/internal/respond"
	"` + "{{MODULE}}" + `/internal/services"
	"` + "{{MODULE}}" + `/internal/storage"
	"` + "{{MODULE}}" + `/internal/flags"
	"` + "{{MODULE}}" + `/internal/sync"
	"` + "{{MODULE}}" + `/internal/webhooks"
	// Imports added by plugins.
	// grit:imports
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
` + eventBusStatusBlock + llmsRegisterFunc + `
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

` + routesAuthLimitsFunc + `
// wafExcludedRoutes lists the paths Sentinel's WAF steps aside for, under the
// live API prefix.
//
// Uploads are here because the WAF rejects any body over its inspection cap
// before the route runs, and a photograph is larger than that cap by design.
// The richtext resources are here because the XSS heuristics flag ordinary
// markup: a blog body is <p> and <strong> and <img> by definition.
//
// Every richtext resource is listed twice, once at its public path and once
// under /admin, because a body is inspected on the way in and the writes are
// what carry the markup: the admin panel's are the /admin ones. With only the
// public path listed, a post whose body quotes a shell command could be
// published and then never edited, every PUT from the admin answered 403.
//
// Exclusion is from body inspection only. These routes still pass through
// auth, RBAC, binding validation, sanitize:"html" and rate limiting.
func wafExcludedRoutes() []string {
	prefix := "/api/" + APIVersion
	paths := []string{
		"/blogs", "/blogs/*",
		"/admin/blogs", "/admin/blogs/*",
		"/posts", "/posts/*",
		"/admin/posts", "/admin/posts/*",
		"/articles", "/articles/*",
		"/admin/articles", "/admin/articles/*",
		"/uploads", "/uploads/*",
` + wafImportExclusion + `		// Public form-share submissions. Auth is the share's bcrypt password
		// (optional) and the token itself; Sentinel rate-limits the path. The
		// subtree match also covers .../submit.
		"/public/forms/*",
` + wafRichtextMarkerBlock + `
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

// Setup wires the application, and these five are the parts of it that are
// configuration rather than wiring.
//
// Setup was 902 lines. Most of that was not the shape of the application, which
// is what you open routes.go to see: it was the Sentinel config literal, the
// Pulse options, the Studio mount and the middleware chain, in the middle of
// it. They are here now, each named for what it mounts, and Setup reads as the
// sequence it always was.
//
// What deliberately stays in Setup: the handler construction and every route
// group with a // grit: marker in it. Those markers are where plugins and
// grit generate resource inject code, and injected code refers to handlers
// Setup builds, so the marker and the variable have to share a scope.

// mountGlobalMiddleware installs the chain every request passes through.
//
// Order matters and is the reason this is one function rather than a list:
// maintenance mode answers before anything else does the work, the request id
// exists before the logger prints it, and CSRF runs after the body limits so an
// oversized body is refused before it is parsed.
// It returns the CORS origin resolver, because two things outside the chain
// read the same list: the realtime socket, which accepts a cookie handshake
// only from an origin CORS allows.
func mountGlobalMiddleware(r *gin.Engine, cfg *config.Config, svc *Services) func() []string {
	r.Use(middleware.Maintenance())
	r.Use(middleware.SecurityHeaders())
` + routesRequestLimitsBlock + `	r.Use(middleware.RequestID())
	// A span per request, continuing one the caller sent. Mounted after
	// RequestID because it replaces the generated id with the trace id when
	// the request is sampled, so a log line and a span share one id, and
	// before RequestMeta so what RequestMeta puts on the context is that id.
	//
	// Unconditional: with no OTEL_EXPORTER_OTLP_ENDPOINT the provider is a
	// no-op and the span costs an allocation. A conditional mount is a line
	// somebody reorders.
	r.Use(tracing.Middleware())
	// The client IP, the user agent and the request id, on the request's
	// context, so a service can read them without taking a *gin.Context.
	r.Use(middleware.RequestMeta())
	// Request timing, for the percentiles grit scale reads. A ring buffer
	// of the last few thousand requests: no allocation, no time series, and the
	// only question it answers is how the service is doing right now.
	r.Use(middleware.Latency())
	// Pins a person's reads to the primary for a few seconds after they write,
	// so their own post is never missing from their own feed. A no-op until
	// DATABASE_REPLICA_URLS is set.
	r.Use(middleware.ReadYourWrites())
	r.Use(middleware.Logger())
	r.Use(gin.Recovery())
	// Origins come from the cors.origins setting when it has a value, and from
	// CORS_ORIGINS otherwise. Resolved per request, so adding a domain in the
	// admin takes effect immediately rather than at the next deploy.
` + routesCORSNew + `	r.Use(middleware.Gzip())

	// CSRF defence — only enforces on cookie-authenticated mutations.
	// Bearer (mobile/desktop) flows pass through with no header required.
	// Pairs with services.AuthService.SetAuthCookies (the HttpOnly cookie
	// path documented in /docs/backend/authentication).
	r.Use(middleware.AutoCSRF())

	// Idempotent retries for unsafe methods. Activates only when the client
	// sends an Idempotency-Key header; cached for 24h on 2xx responses.
	r.Use(middleware.Idempotency(svc.Cache))

	return corsOrigins
}

// mountSentinel mounts the security suite: WAF, rate limiting, auth shield and
// anomaly detection, with its dashboard at /sentinel.
//
// A mount failure is logged rather than fatal. Sentinel refuses to start on a
// misconfiguration, and in development that should not take the API down with
// it.
func mountSentinel(r *gin.Engine, db *gorm.DB, cfg *config.Config, svc *Services) {
	// Mount Sentinel security suite (WAF, rate limiting, auth shield, anomaly detection)
	if cfg.SentinelEnabled {
		// In development, use relaxed rate limits so devs don't get blocked while testing
		isDev := cfg.AppEnv == "development"
		ipLimit := &sentinel.Limit{Requests: 100, Window: 1 * time.Minute}
` + routesLimitsBlock + `
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
		// Keys the audit log's hash chain: an entry edited by someone with the
		// database but not this key fails verification on the dashboard.
		sentinelStorage.AuditKey = cfg.SentinelAuditKey
` + sentinelPoolLine + `

		// Rate limits and AuthShield lockouts counted in Redis, so replicas
		// share them. Counted per process, N replicas gave a client N times
		// every limit and N times the failed logins before a lockout.
		var sentinelCounters sentinel.CounterStore
		if svc.Cache != nil {
			sentinelCounters = redisstore.New(svc.Cache.Client())
		}

		// Sentinel v2 — use MountE so we can recover gracefully on
		// misconfiguration in dev instead of log.Fatalf-ing the host.
		// Mount runs sentinel.ValidateConfig and logs any dead config.
		if err := sentinel.MountE(r, db, sentinel.Config{
			Storage:  sentinelStorage,
			Counters: sentinelCounters,
			// Who made each request, for the Users page, the GDPR report and
			// anomaly detection, which sees nothing without it. Sentinel reads
			// it after the handler chain, by when the auth middleware has put
			// the caller on the context.
			UserExtractor: func(c *gin.Context) *sentinel.UserContext {
				id := c.GetString("user_id")
				if id == "" {
					return nil
				}
				return &sentinel.UserContext{ID: id, Email: c.GetString("user_email"), Role: c.GetString("user_role")}
			},
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
				LoginRoute: "/api/" + APIVersion + "/auth/login",
				// v2.0 CAPTCHA tier sits between soft and hard thresholds.
				// Wire a provider by setting CaptchaProvider in your app code.
			},
			Anomaly: sentinel.AnomalyConfig{Enabled: !isDev},
			Geo:     sentinel.GeoConfig{Enabled: !isDev},
		}); err != nil {
			log.Printf("Warning: Sentinel mount failed: %v", err)
		} else {
			log.Println("Sentinel mounted at /sentinel")
		}
	}
}

// mountStudio mounts the database browser at /studio.
func mountStudio(r *gin.Engine, db *gorm.DB, cfg *config.Config) {
	// Mount GORM Studio
	if cfg.GORMStudioEnabled {
		studioCfg := studio.Config{
			Prefix:     "/studio",
			ReadOnly:   cfg.GORMStudioReadOnly,
			DisableSQL: cfg.GORMStudioDisableSQL,
		}
		if cfg.GORMStudioUsername != "" && cfg.GORMStudioPassword != "" {
			studioCfg.AuthMiddleware = gin.BasicAuth(gin.Accounts{
				cfg.GORMStudioUsername: cfg.GORMStudioPassword,
			})
		}
		studio.Mount(r, db, []interface{}{&models.User{}, &models.Upload{}, &models.Blog{}, /* grit:studio */}, studioCfg)
		log.Println("GORM Studio mounted at /studio")
	}
}

// mountPulse mounts observability at /pulse: request tracing, database
// monitoring, runtime metrics and error tracking.
func mountPulse(r *gin.Engine, db *gorm.DB, cfg *config.Config, svc *Services) {
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
			// Request-body capture is left on, which needs pulse v1.0.1 or
			// later. Before that the error middleware read the body and put
			// back only the first 4 KB, so every request carrying a
			// Content-Length reached the handler truncated: uploads from
			// mobile and curl failed while browsers, which send chunked,
			// did not. v1.2.0 wraps the body instead and passes it through
			// whole. Pin below v1.0.1 and you want
			// pulse.WithRequestBodyCaptureDisabled() back.
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
}

// mountAuthRoutes mounts everything a caller with no session yet can reach:
// the password flow, OAuth, enterprise SSO over OIDC and SAML, and the two
// second factors.
//
// Public by necessity. These are the sign-in, so there is nobody to
// authenticate yet; what makes each safe is its own protocol, a server-side
// challenge for passkeys and a short-lived pending token for TOTP.
func mountAuthRoutes(v1 *gin.RouterGroup, authHandler *handlers.AuthHandler, ssoHandler *handlers.SSOHandler, totpHandler *handlers.TOTPHandler, passkeyHandler *handlers.PasskeyHandler) {
	// Public auth routes
	auth := v1.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
		auth.POST("/refresh", authHandler.Refresh)
		auth.POST("/forgot-password", authHandler.ForgotPassword)
		auth.POST("/reset-password", authHandler.ResetPassword)
		auth.POST("/verify-email", authHandler.VerifyEmail)
		// Signing in with an emailed link. Consuming is a POST from the page
		// the link opens, never the GET that opens it: mail scanners follow
		// every URL in a message and would spend the token first.
		auth.POST("/magic-link", authHandler.RequestMagicLink)
		auth.POST("/magic-link/consume", authHandler.ConsumeMagicLink)
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
	// Passkey sign-in. Public because there is no session yet; the
	// server-side challenge is what makes it safe.
	auth.POST("/passkeys/login/begin", passkeyHandler.BeginLogin)
	auth.POST("/passkeys/login/finish", passkeyHandler.FinishLogin)
	auth.POST("/totp/verify", totpHandler.Verify)
	auth.POST("/totp/backup-codes/verify", totpHandler.VerifyBackupCode)
}

// Setup configures all routes and returns the Gin engine.
func Setup(db *gorm.DB, cfg *config.Config, svc *Services) *gin.Engine {
	if cfg.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

` + routesTrustProxiesBlock + `	corsOrigins := mountGlobalMiddleware(r, cfg, svc)

	mountSentinel(r, db, cfg, svc)

	mountStudio(r, db, cfg)

	// API Documentation (gin-docs — auto-generated from routes + models)
	//
	// The OpenAPI reference. Its 141 route overrides live in apidocs.go, where
	// they are 500 lines of description rather than 500 lines in the middle of
	// the file that wires your application together.
` + routesDocsBlock + llmsRoutesBlock + `
	mountPulse(r, db, cfg, svc)

	// Auth service
	authService := &services.AuthService{
		Secret:        cfg.JWTSecret,
		AccessExpiry:  cfg.JWTAccessExpiry,
		RefreshExpiry: cfg.JWTRefreshExpiry,
		// Sessions are checked on every access token, so logging out, signing out
		// everywhere and changing a password take effect at once.
		DB: db,
	}
	// The refresh cookie is scoped to the auth routes, wherever the version puts
	// them.
	services.RefreshCookiePath = "/api/" + APIVersion + "/auth"

	// Handlers
	authHandler := &handlers.AuthHandler{
		DB:          db,
		AuthService: authService,
		Config:      cfg,
		Mailer:      svc.Mailer,
		Jobs:        svc.Jobs,
	}
	apiKeyHandler := &handlers.APIKeyHandler{DB: db}

	userHandler := &handlers.UserHandler{
		DB:          db,
		AuthService: authService,
	}
` + routesUploadHandlerNew + `	aiHandler := &handlers.AIHandler{
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
		// For the email second factor: the same mailer and queue the auth
		// handler uses, so a code goes out the way every other email does.
		Config: cfg,
		Mailer: svc.Mailer,
		Jobs:   svc.Jobs,
	}
	// The passkey relying party, built once from the origins the frontends
	// actually run on. A deployment with none (CORS_ORIGINS unset or '*')
	// gets a nil service, and every passkey route answers 501 rather than
	// panicking: passkeys are optional, a broken boot is not.
	passkeys, passkeyErr := services.NewPasskeys(db, cfg.AppName, cfg.CORSOrigins)
	if passkeyErr != nil {
		log.Printf("Passkeys disabled: %v", passkeyErr)
		passkeys = nil
	}
	passkeyHandler := handlers.NewPasskeyHandler(db, passkeys, authHandler)
	activityHandler := handlers.NewActivityHandler(db)
	webhookHandler := handlers.NewWebhookHandler(db)
	webhooks.Setup(db)
	// WithRedis is a no-op when REDIS_URL is empty, so a single-process
	// project keeps the in-process hub and a multi-replica one fans out
	// across every instance without a second thing to configure. Without a
	// backplane, a user on replica A never sees an event published on
	// replica B and nothing anywhere reports it.
` + routesRealtimeHubNew + `	// Revoking a session has to close that user's live sockets too. Without
	// this, "sign out of all devices" leaves every open WebSocket streaming.
	services.OnSessionsRevoked = realtimeHub.DisconnectUser

	// The domain event bus. Created before any handler so an emit during
	// startup has somewhere to go, and wired to the audit log, realtime and
	// (when the plugin is installed) outbound webhooks.
	events.Init(4)
	services.RegisterEventSubscribers(db, realtimeHub, nil)
	// Durable subscribers are delivered from the outbox by this relay. Every
	// replica runs one; row claims keep two from delivering a message twice.
	events.StartRelay(db)
` + sagaRuntimeBlock + `
	// Settings: declare, then open the store. Declaring after Init would mean
	// a setting the first cache load never saw.
	settings.RegisterDefaults()
	settings.Init(db)
	settingsHandler := &handlers.SettingsHandler{DB: db}
	flagsEngine := flags.New(db, realtimeHub)
	featureFlagHandler := handlers.NewFeatureFlagHandler(db, flagsEngine)
` + routesRealtimeHandlerNew + `
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
	ticketHandler := &handlers.TicketHandler{DB: db, Mail: svc.Mailer, Jobs: svc.Jobs}
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
	// Never users or uploads. A push is a generic write: a user row carries its
	// own role and email, an upload row the key of any object in the bucket.
	// Both have their own APIs with their own checks, and syncing them let any
	// account make itself ADMIN.
	syncRegistry.Register("blogs", &models.Blog{})
	// grit:sync
	syncHandler := handlers.NewSyncHandler(db, syncRegistry)
	// The bin reads the same registry: every resource that can be synced is a
	// resource whose rows soft-delete, which is the same list.
	trashHandler := handlers.NewTrashHandler(db, syncRegistry)
	// v3.31.68 — shared background CSV import status endpoint
	importJobHandler := &handlers.ImportJobHandler{DB: db}
` + sagaHandlerBlock + `` + routeExplorerHandlerBlock + `	// v3.31.77 — full-database backups (weekly cron + manual + download)
	backupHandler := &handlers.BackupHandler{DB: db, Storage: svc.Storage}
	roleHandler := handlers.NewRoleHandler(db)
	// Which scaling stage this deployment is at, measured rather than guessed.
	// Reads the live request percentiles, the connection count against the
	// ceiling, and what Postgres says about its own slowest queries.
	scaleHandler := &handlers.ScaleHandler{DB: db, Cache: svc.Cache}
	// Permission caches are per process. Share makes a role change on one
	// replica reach every other within a second.
	authz.Share(db)
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
	// Rebuild the SSO registries when another replica changes a connection.
	services.WatchSSO(db, ssoRegistry, samlRegistry)
	// grit:handlers

` + routesLocalFilesBlock + healthQueueStatsSetup + `	// Health check
	// /api/health probes every infrastructure dependency the dashboard's
	// System Health page wants to render. Each probe is bounded by a 500ms
	// timeout so a hung dependency doesn't pile up health requests; failing
	// probes mark themselves down and the overall status downgrades to
	// "degraded" rather than failing the endpoint.
	// Registered at two paths, not one. The frontends' axios client rewrites
	// /api/... to /api/v1/..., so the admin's System Health page asked for
	// /api/v1/health, no route matched, and the unversioned fallback refused to
	// rewrite a path that already names the version: the page read "degraded"
	// with a 404 behind it. /api/health stays for probes, load balancers and the
	// desktop client's heartbeat, which are configured outside this repo.
` + healthCheckBlock + `	r.GET("/api/health", healthCheck)
	r.GET("/api/"+APIVersion+"/health", healthCheck)

` + routesRealtimeRouteNew + `
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

	mountAuthRoutes(v1, authHandler, ssoHandler, totpHandler, passkeyHandler)

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
	// Request middleware added by plugins. It runs after auth, so anything
	// here can read the authenticated user.
	// grit:middleware:protected
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
		// The second factor by email rather than by app. Two steps, like the
		// authenticator: send a code to prove the address works, then confirm
		// it, so nobody turns on a factor they cannot receive.
		// Who has been asking for sign-in links to this account. Read-only,
		// and never the token itself.
		protected.GET("/auth/magic-link/recent", authHandler.RecentMagicLinks)
		protected.POST("/auth/totp/email/send", totpHandler.SendEmailSetupCode)
		protected.POST("/auth/totp/email/enable", totpHandler.EnableEmail)
		protected.POST("/auth/totp/disable", totpHandler.Disable)
		protected.GET("/auth/totp/status", totpHandler.Status)
		protected.POST("/auth/totp/backup-codes", totpHandler.RegenerateBackupCodes)
		protected.DELETE("/auth/totp/trusted-devices", totpHandler.RevokeTrustedDevices)
		protected.POST("/auth/verify-email/send", authHandler.SendVerificationEmail)

		// Recovery contacts. Every write takes the account password, because a
		// recovery address is a second way in and a live session is exactly what
		// somebody on a borrowed laptop already has.
		recoveryHandler := handlers.NewRecoveryHandler(db, svc.Mailer)
		// Passkeys. Registration is behind auth because you add one to an
		// account you are already in; the sign-in pair is public by necessity.
		protected.GET("/auth/passkeys", passkeyHandler.List)
		protected.POST("/auth/passkeys/register/begin", passkeyHandler.BeginRegistration)
		protected.POST("/auth/passkeys/register/finish", passkeyHandler.FinishRegistration)
		protected.PATCH("/auth/passkeys/:id", passkeyHandler.Rename)
		protected.DELETE("/auth/passkeys/:id", passkeyHandler.Delete)
		protected.GET("/auth/security", recoveryHandler.Overview)
		protected.POST("/auth/recovery/email", recoveryHandler.SetEmail)
		protected.POST("/auth/recovery/email/verify", recoveryHandler.VerifyEmail)
		protected.DELETE("/auth/recovery/email", recoveryHandler.ClearEmail)
		protected.POST("/auth/recovery/phone", recoveryHandler.SetPhone)
		protected.POST("/auth/recovery/phone/verify", recoveryHandler.VerifyPhone)
		protected.DELETE("/auth/recovery/phone", recoveryHandler.ClearPhone)
		protected.GET("/api-keys", apiKeyHandler.List)
		protected.POST("/api-keys", apiKeyHandler.Create)
		protected.DELETE("/api-keys/:id", apiKeyHandler.Revoke)
		protected.GET("/auth/totp/trusted-devices", totpHandler.ListTrustedDevices)
		protected.DELETE("/auth/totp/trusted-devices/:id", totpHandler.RevokeTrustedDevice)

		// User routes (authenticated)

		// GDPR right-to-access: a user may export their own data; an admin, anyone's.
		protected.GET("/users/:id/gdpr-export", gdprHandler.Export)

		// File uploads
		protected.POST("/uploads", uploadHandler.Create)
		// The client optimises before it uploads, so it needs the same numbers
		// the server would have used.
		protected.GET("/media/profiles", uploadHandler.Profiles)
		protected.POST("/uploads/presign", uploadHandler.Presign)
		protected.POST("/uploads/complete", uploadHandler.CompleteUpload)
		protected.GET("/uploads", uploadHandler.List)
		protected.GET("/uploads/stats", uploadHandler.Stats)
` + routesUploadDownload + `		protected.DELETE("/uploads/:id", uploadHandler.Delete)

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

	// Staff routes: anyone who holds a permission reaches this group, and each
	// route names the one it needs. Before v3.220.0 all of these sat behind the
	// ADMIN role, so a custom role granted users.view was refused by every admin
	// endpoint. The admin group below stays ADMIN-only, so a route that names no
	// permission, a plugin's included, fails closed there.
	staff := v1.Group("")
	staff.Use(middleware.APIKeyOrAuth(db, middleware.Auth(db, authService)))
	// Plugin middleware for the staff group, between authentication and the gate.
	//
	// This group carries every DELETE and every bulk route, so a plugin that
	// scopes queries has to run here: with the multitenant plugin mounted on the
	// protected group alone, deleting a tenant-owned row answered 500 because no
	// organization was ever resolved.
	//
	// Before the gate, not after, because a plugin can decide what the caller may
	// do: a role held through an organization membership has to be in hand before
	// RequireStaff reads the grants, or the answer is 403 for somebody who is
	// staff of the organization they are acting in.
	// grit:middleware:staff
	staff.Use(middleware.RequireStaff())
	{
		staff.GET("/users", middleware.RequireRole("ADMIN", "perm:users.view"), userHandler.List)
` + userByIDStaffRoute + `		staff.POST("/users", middleware.RequireRole("ADMIN", "perm:users.create"), userHandler.Create)
		staff.PUT("/users/:id", middleware.RequireRole("ADMIN", "perm:users.edit"), userHandler.Update)
		staff.DELETE("/users/:id", middleware.RequireRole("ADMIN", "perm:users.delete"), userHandler.Delete)
		staff.PUT("/users/:id/roles", middleware.RequireRole("ADMIN", "perm:users.edit"), roleHandler.AssignUserRoles)
		staff.POST("/users/:id/unlock", middleware.RequireRole("ADMIN", "perm:users.edit"), userHandler.Unlock)

		// GDPR right-to-erasure: anonymize a user and hard-delete their PII, with
		// the erasure recorded in a tamper-evident deletion journal.
		staff.POST("/users/:id/gdpr-erase", middleware.RequireRole("ADMIN", "perm:users.delete"), gdprHandler.Erase)
		staff.GET("/gdpr/journal", middleware.RequireRole("ADMIN", "perm:audit.view"), gdprHandler.Journal)

		// Activity: the tamper-evident log and its verification, the semantic
		// activity feed, and its OCSF export for SIEMs. Resealing is admin-only.
		staff.GET("/admin/activity", middleware.RequireRole("ADMIN", "perm:audit.view"), activityHandler.List)
		staff.GET("/admin/activity/integrity", middleware.RequireRole("ADMIN", "perm:audit.view"), activityHandler.VerifyIntegrity)
		staff.GET("/user-activity", middleware.RequireRole("ADMIN", "perm:audit.view"), userActivityHandler.List)
		staff.GET("/user-activity/stats", middleware.RequireRole("ADMIN", "perm:audit.view"), userActivityHandler.Stats)
		staff.GET("/audit/ocsf", middleware.RequireRole("ADMIN", "perm:audit.view"), ocsfHandler.Export)

		// Roles, the permission catalog, and access reviews over the grants.
		staff.GET("/permissions", middleware.RequireRole("ADMIN", "perm:roles.view"), roleHandler.Catalog)
		// Admin only, and rightly: it reports connection counts, the slowest
		// queries and their shapes. That is a map of the database to anybody
		// who can read it.
		staff.GET("/scale", middleware.RequireRole("ADMIN", "perm:system.view"), scaleHandler.Report)
		staff.GET("/roles", middleware.RequireRole("ADMIN", "perm:roles.view"), roleHandler.List)
		staff.POST("/roles", middleware.RequireRole("ADMIN", "perm:roles.create"), roleHandler.Create)
		staff.GET("/roles/:id", middleware.RequireRole("ADMIN", "perm:roles.view"), roleHandler.Get)
		staff.PUT("/roles/:id", middleware.RequireRole("ADMIN", "perm:roles.edit"), roleHandler.Update)
		staff.DELETE("/roles/:id", middleware.RequireRole("ADMIN", "perm:roles.delete"), roleHandler.Delete)
		staff.GET("/access-reviews", middleware.RequireRole("ADMIN", "perm:roles.view"), accessReviewHandler.List)
		staff.POST("/access-reviews", middleware.RequireRole("ADMIN", "perm:roles.edit"), accessReviewHandler.Open)
		staff.GET("/access-reviews/:id", middleware.RequireRole("ADMIN", "perm:roles.view"), accessReviewHandler.Get)
		staff.POST("/access-reviews/:id/items/:itemId/decision", middleware.RequireRole("ADMIN", "perm:roles.edit"), accessReviewHandler.Decide)
		staff.POST("/access-reviews/:id/complete", middleware.RequireRole("ADMIN", "perm:roles.edit"), accessReviewHandler.Complete)

		// Operations: jobs and the schedule, and the dashboards over Sentinel,
		// Pulse, webhooks and feature flags. Changing any of those is admin-only.
		staff.GET("/admin/jobs/stats", middleware.RequireRole("ADMIN", "perm:jobs.view"), jobsHandler.Stats)
		staff.GET("/admin/jobs/:status", middleware.RequireRole("ADMIN", "perm:jobs.view"), jobsHandler.ListByStatus)
		staff.POST("/admin/jobs/:id/retry", middleware.RequireRole("ADMIN", "perm:jobs.edit"), jobsHandler.Retry)
		staff.DELETE("/admin/jobs/queue/:queue", middleware.RequireRole("ADMIN", "perm:jobs.edit"), jobsHandler.ClearQueue)
` + sagaRoutesBlock + `` + routeExplorerRoutesBlock + `` + routesMailPreviewAnchor + routesMailPreview + `		staff.GET("/admin/security/summary", middleware.RequireRole("ADMIN", "perm:system.view"), securityHandler.Summary)
		staff.GET("/admin/observability/summary", middleware.RequireRole("ADMIN", "perm:system.view"), observabilityHandler.Summary)

		// The bin. Reading it is a system view; restoring somebody else's
		// delete, or making one permanent, is ADMIN and nothing less: a
		// resource's own delete permission says you may remove a row, not that
		// you may undo another person's removal or put it beyond recovery.
		staff.GET("/admin/trash", middleware.RequireRole("ADMIN", "perm:system.view"), trashHandler.Buckets)
		staff.GET("/admin/trash/:table", middleware.RequireRole("ADMIN", "perm:system.view"), trashHandler.List)
		staff.POST("/admin/trash/:table/:id/restore", middleware.RequireRole("ADMIN"), trashHandler.Restore)
		staff.DELETE("/admin/trash/:table/:id", middleware.RequireRole("ADMIN"), trashHandler.Purge)
		staff.DELETE("/admin/trash/:table", middleware.RequireRole("ADMIN"), trashHandler.Empty)

		// Closed accounts, which are a soft delete the bin cannot show: users
		// are deliberately not in the sync registry it reads, because a row
		// carrying its own role is not something a generic restore should write.
		staff.GET("/admin/deleted-accounts", middleware.RequireRole("ADMIN"), userHandler.DeletedAccounts)
		staff.POST("/admin/deleted-accounts/:id/restore", middleware.RequireRole("ADMIN"), userHandler.RestoreAccount)
		staff.DELETE("/admin/deleted-accounts/:id", middleware.RequireRole("ADMIN"), userHandler.PurgeAccount)
		staff.GET("/admin/webhooks", middleware.RequireRole("ADMIN", "perm:system.view"), webhookHandler.List)
		staff.GET("/admin/flags", middleware.RequireRole("ADMIN", "perm:system.view"), featureFlagHandler.List)
		staff.GET("/admin/flags/:id/exposures", middleware.RequireRole("ADMIN", "perm:system.view"), featureFlagHandler.Exposures)

		// Blog management.
		staff.GET("/admin/blogs", middleware.RequireRole("ADMIN", "perm:blogs.view"), blogHandler.List)
		staff.GET("/admin/blogs/:id", middleware.RequireRole("ADMIN", "perm:blogs.view"), blogHandler.GetByID)
		staff.POST("/admin/blogs", middleware.RequireRole("ADMIN", "perm:blogs.create"), blogHandler.Create)
		staff.PUT("/admin/blogs/:id", middleware.RequireRole("ADMIN", "perm:blogs.edit"), blogHandler.Update)
		staff.DELETE("/admin/blogs/:id", middleware.RequireRole("ADMIN", "perm:blogs.delete"), blogHandler.Delete)

		// Per-resource dashboard stats and charts. The resource is in the URL, so
		// the permission is that resource's view. Only resources registered in
		// the stats and chart dispatchers are reachable at all.
		staff.GET("/admin/dashboard/resource-stats/:resource", middleware.RequirePermissionFor("resource", "view"), resourceStatsHandler.Get)
		staff.GET("/admin/dashboard/chart/:resource", middleware.RequirePermissionFor("resource", "view"), chartHandler.Get)

		// Full-database backups: a weekly cron writes them, and an operator can
		// take one on demand (once a day) and download it through a short-lived
		// pre-signed URL. The settings live at their own path so they do not
		// collide with the /backups/:id wildcard.
		staff.GET("/backups", middleware.RequireRole("ADMIN", "perm:backups.view"), backupHandler.List)
		staff.POST("/backups/generate", middleware.RequireRole("ADMIN", "perm:backups.create"), backupHandler.Generate)
		staff.GET("/backups/:id/download", middleware.RequireRole("ADMIN", "perm:backups.view"), backupHandler.Download)
		staff.GET("/backup-settings", middleware.RequireRole("ADMIN", "perm:backups.view"), backupHandler.GetSettings)
		staff.PUT("/backup-settings", middleware.RequireRole("ADMIN", "perm:backups.edit"), backupHandler.UpdateSettings)
	}

	// Admin routes: the ADMIN role and nothing less. Anything that should be
	// grantable to a custom role belongs in the staff group above, naming its
	// permission.
	admin := v1.Group("")
	admin.Use(middleware.APIKeyOrAuth(db, middleware.Auth(db, authService)))
	// Plugin middleware for the admin group, for the same reasons and in the same
	// order as the staff one.
	// grit:middleware:admin
	admin.Use(middleware.RequireRole("ADMIN"))
	{
		admin.POST("/admin/activity/reseal", activityHandler.Reseal)
		admin.POST("/admin/webhooks/:id/replay", webhookHandler.Replay)
		admin.POST("/admin/flags", featureFlagHandler.Create)
		admin.PUT("/admin/flags/:id", featureFlagHandler.Update)
		admin.DELETE("/admin/flags/:id", featureFlagHandler.Delete)

		// SSO connections. Client secrets are write-only: they go in on create
		// and update and are never returned, so a compromised admin session
		// cannot read a customer's IdP credentials back out.
		admin.GET("/sso/connections", ssoHandler.List)
		admin.POST("/sso/connections", ssoHandler.Create)
		admin.PUT("/sso/connections/:id", ssoHandler.Update)
		admin.DELETE("/sso/connections/:id", ssoHandler.Delete)
		admin.GET("/sso/connections/:id/test", ssoHandler.Test)

		// Public form sharing, its field preview, and the log of submissions.
		admin.GET("/admin/form-shares", formShareHandler.List)
		admin.POST("/admin/form-shares", formShareHandler.Create)
		admin.PATCH("/admin/form-shares/:id", formShareHandler.Update)
		admin.DELETE("/admin/form-shares/:id", formShareHandler.Delete)
		admin.GET("/admin/form-shares/resources", formShareHandler.Resources)
		admin.GET("/admin/form-shares/resources/:resource/fields", formShareHandler.FieldsPreview)
		admin.GET("/admin/form-submissions", formShareHandler.ListSubmissions)

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

	// Every generated resource, each from its own <resource>_routes.go.
	//
	// A resource file registers itself from an init(), so this loop is the
	// only place routes.go mentions them. Adding a resource does not edit this
	// file, and neither does removing one.
	mountResources(&Mount{
		Engine:    r,
		DB:        db,
		Cfg:       cfg,
		Svc:       svc,
		Hub:       realtimeHub,
		Auth:      authService,
		V1:        v1,
		Public:    publicAPI,
		Protected: protected,
		Admin:     admin,
		Staff:     staff,
	})

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

		respond.Fail(c, respond.CodeNotFound, "no route matches " + c.Request.Method + " " + p)
	})
}
`
}
