package scaffold

// The OpenAPI reference lives in its own file.
//
// It was 28KB and 141 route overrides sitting in the middle of routes.go,
// which is the file you open to find out how the application is wired. Route
// descriptions are documentation; they were 38% of a file about routing.
//
// The generator injects into the markers here rather than in routes.go now.
// Projects generated before this still have the markers in routes.go, so the
// injector tries this file first and falls back.

// apiDocsRoutesGo emits internal/routes/apidocs.go.
func apiDocsRoutesGo() string {
	return `package routes

import (
	"log"

	"github.com/MUKE-coder/gin-docs/gindocs"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"` + "{{MODULE}}" + `/internal/ai"
	"` + "{{MODULE}}" + `/internal/config"
	"` + "{{MODULE}}" + `/internal/handlers"
	"` + "{{MODULE}}" + `/internal/models"
	"` + "{{MODULE}}" + `/internal/services"
)

// registerAPIDocs mounts the OpenAPI reference at /docs and describes every
// endpoint worth describing.
//
// Called once from Setup, after the routes exist: gindocs reads the router to
// find paths and status codes, and the overrides below are what attach request
// and response schemas so an operation renders a body rather than "No Body".
func registerAPIDocs(r *gin.Engine, db *gorm.DB, cfg *config.Config) {
	// Route introspection alone gives paths and status codes but no bodies:
	// gindocs only attaches a request or response schema when a route override
	// hands it a concrete type. Without the docs.Route(...) calls below, every
	// operation renders as "No Body" and the copy-paste curl snippets are
	// unusable. The spec is assembled lazily on the first /docs request, so
	// overrides registered after Mount are picked up.
	docs := gindocs.Mount(r, db, gindocs.Config{
		Title:       cfg.AppName + " API",
		Description: "REST API built with [Grit](https://gritframework.dev) — Go + React meta-framework.",
		Version:     "1.0.0",
		UI:          gindocs.UIScalar,
		ScalarTheme: "kepler",
		// Pulse, Sentinel and GORM Studio mount their own dashboards and APIs
		// inside this app. They are not your API, and 111 of their routes in
		// the reference buries the ~140 that are yours.
		ExcludePrefixes: []string{"/pulse", "/sentinel", "/studio", "/docs"},
		Models:      []interface{}{&models.User{}, &models.Upload{}, &models.Blog{} /* grit:docs:models */},
		Auth: gindocs.AuthConfig{
			Type:         gindocs.AuthBearer,
			BearerFormat: "JWT",
		},
	})

	// Core endpoints. The generated summaries are derived from the path, which
	// turns POST /auth/login into "Create a new login" — worth overriding for
	// the handful of routes every reader hits first.
	docs.Route("POST /api/v1/auth/register").
		Summary("Register a new account").
		RequestBody(handlers.RegisterRequest{}).
		Response(201, handlers.AuthResponse{}, "Account created").
		Response(422, handlers.ErrorResponse{}, "Validation failed")
	docs.Route("POST /api/v1/auth/login").
		Summary("Sign in").
		RequestBody(handlers.LoginRequest{}).
		Response(200, handlers.AuthResponse{}, "Signed in").
		Response(401, handlers.ErrorResponse{}, "Invalid credentials")
	docs.Route("POST /api/v1/auth/refresh").
		Summary("Exchange a refresh token for a new access token").
		RequestBody(handlers.RefreshRequest{}).
		Response(200, handlers.AuthResponse{}, "New token pair")
	docs.Route("GET /api/v1/auth/me").
		Summary("The signed-in user").
		Response(200, models.User{}, "The current user")

	// Built-in endpoints. Response schemas come from the models, which is why
	// these are one line each; request bodies appear only where a named type
	// exists, since gindocs reflects over a type and an anonymous struct
	// inside a handler gives it nothing to read.
	docs.Route("POST /api/v1/auth/forgot-password").
		Summary("Send a password-reset link").
		RequestBody(handlers.ForgotPasswordRequest{}).
		Response(200, handlers.MessageResponse{}, "Always the same reply, whether or not the address exists")
	docs.Route("POST /api/v1/auth/reset-password").
		Summary("Set a new password with a reset token").
		RequestBody(handlers.ResetPasswordRequest{}).
		Response(200, handlers.MessageResponse{}, "Password changed; every other session is signed out")
	docs.Route("POST /api/v1/auth/verify-email").
		Summary("Confirm an email address").
		RequestBody(handlers.VerifyEmailRequest{}).
		Response(200, handlers.MessageResponse{}, "Address verified").
		Response(400, handlers.ErrorResponse{}, "The link is invalid or has expired")
	docs.Route("POST /api/v1/auth/verify-email/send").
		Summary("Send a fresh verification link").
		Response(200, handlers.MessageResponse{}, "Sent")
	docs.Route("POST /api/v1/auth/logout").
		Summary("Sign out this device").
		Response(200, handlers.MessageResponse{}, "Signed out")
	docs.Route("GET /api/v1/auth/sessions").
		Summary("Devices signed in to this account").
		Response(200, []models.Session{}, "One row per active session")
	docs.Route("DELETE /api/v1/auth/sessions/:id").
		Summary("Sign one device out").
		Response(200, handlers.MessageResponse{}, "That device is signed out")
	docs.Route("POST /api/v1/auth/sessions/revoke-all").
		Summary("Sign out every other device").
		Response(200, handlers.MessageResponse{}, "All other sessions revoked")
	docs.Route("GET /api/v1/auth/totp/status").
		Summary("Two-factor status for the signed-in user").
		Response(200, handlers.TOTPStatusResponse{}, "Whether 2FA is on, codes left, trusted devices")
	docs.Route("GET /api/v1/auth/totp/trusted-devices").
		Summary("Devices allowed to skip the 2FA prompt").
		Response(200, []models.TrustedDevice{}, "Trusted devices, newest first")

	docs.Route("GET /api/v1/users").
		Summary("List users").
		Response(200, []models.User{}, "A page of users")
	docs.Route("GET /api/v1/users/:id").
		Summary("Get one user").
		Response(200, models.User{}, "The user").
		Response(404, handlers.ErrorResponse{}, "Not found")
	docs.Route("GET /api/v1/users/:id/gdpr-export").
		Summary("Download everything held about a user (GDPR Art. 15)").
		Response(200, handlers.MessageResponse{}, "A JSON attachment: profile, uploads, sessions, activity")
	docs.Route("POST /api/v1/users/:id/gdpr-erase").
		Summary("Erase a user's personal data (GDPR Art. 17)").
		RequestBody(handlers.EraseRequest{}).
		Response(200, handlers.MessageResponse{}, "Personal records deleted, the account anonymised")
	docs.Route("POST /api/v1/users/:id/unlock").
		Summary("Clear a login lockout").
		Response(200, handlers.MessageResponse{}, "The account can sign in again")

	docs.Route("GET /api/v1/profile").
		Summary("The signed-in user's profile").
		Response(200, models.User{}, "The current user")

	docs.Route("GET /api/v1/api-keys").
		Summary("List your API keys").
		Response(200, []models.APIKey{}, "Keys, without their secrets — only hashes are stored")
	docs.Route("POST /api/v1/api-keys").
		Summary("Create an API key").
		RequestBody(handlers.CreateAPIKeyRequest{}).
		Response(201, handlers.IssuedKeyResponse{}, "The key, shown for the only time")
	docs.Route("DELETE /api/v1/api-keys/:id").
		Summary("Revoke an API key").
		Response(200, handlers.MessageResponse{}, "Revoked; requests using it are refused from now on")

	docs.Route("GET /api/v1/notifications").
		Summary("List notifications").
		Response(200, []models.Notification{}, "A page of notifications")

	docs.Route("GET /api/v1/uploads").
		Summary("List uploaded files").
		Response(200, []models.Upload{}, "A page of uploads")
	docs.Route("POST /api/v1/uploads").
		Summary("Upload a file (multipart)").
		Response(201, models.Upload{}, "The stored file").
		Response(400, handlers.ErrorResponse{}, "The file type is not allowed for this field")
	docs.Route("POST /api/v1/uploads/presign").
		Summary("Get a presigned PUT URL for a direct browser upload").
		RequestBody(handlers.PresignRequest{}).
		Response(200, handlers.PresignResponse{}, "A URL to PUT to, and the key to record afterwards")
	docs.Route("POST /api/v1/uploads/complete").
		Summary("Record a file that was uploaded directly to storage").
		RequestBody(handlers.CompleteUploadRequest{}).
		Response(201, models.Upload{}, "The stored file")

	docs.Route("GET /api/v1/roles").
		Summary("List roles").
		Response(200, []models.Role{}, "Every role and its grants")
	docs.Route("GET /api/v1/permissions").
		Summary("The permission catalogue").
		Response(200, handlers.MessageResponse{}, "Modules, groups and features a role can be granted")

	docs.Route("GET /api/v1/backups").
		Summary("List database backups").
		Response(200, []models.Backup{}, "Archives, newest first")
	docs.Route("POST /api/v1/backups/generate").
		Summary("Take a backup now").
		Response(202, models.Backup{}, "The run, which continues in the background").
		Response(429, handlers.ErrorResponse{}, "One was already taken in the last 24 hours")
	docs.Route("GET /api/v1/backup-settings").
		Summary("The backup schedule").
		Response(200, models.BackupSchedule{}, "Frequency and time of day")

	docs.Route("GET /api/v1/access-reviews").
		Summary("List access-review campaigns").
		Response(200, []models.AccessReview{}, "Campaigns, newest first")
	docs.Route("GET /api/v1/access-reviews/:id").
		Summary("One campaign and its items").
		Response(200, models.AccessReview{}, "The campaign, with a row per role assignment")

	docs.Route("GET /api/v1/sso/connections").
		Summary("List SSO connections").
		Response(200, []models.SSOConnection{}, "Connections; secrets are never returned")
	docs.Route("GET /api/v1/gdpr/journal").
		Summary("The tamper-evident deletion journal").
		Response(200, []models.DeletionJournal{}, "Every erasure, hash-chained")

	docs.Route("GET /api/v1/admin/activity").
		Summary("The tamper-evident activity log").
		Response(200, []models.ActivityLog{}, "Every authenticated write, newest first")
	docs.Route("GET /api/v1/admin/activity/integrity").
		Summary("Verify the activity-log hash chain").
		Response(200, handlers.ChainStatusResponse{}, "Valid, or the position and id of the first bad row")
	docs.Route("GET /api/v1/admin/jobs/stats").
		Summary("Background-queue counts").
		Response(200, handlers.MessageResponse{}, "Active, pending, completed, failed and retry totals").
		Response(503, handlers.ErrorResponse{}, "Redis is not configured")
	docs.Route("GET /api/v1/admin/cron/tasks").
		Summary("Registered scheduled tasks").
		Response(200, handlers.MessageResponse{}, "Every cron entry with its schedule")

	// ── Content ────────────────────────────────────────────────────
	docs.Route("GET /api/v1/blogs").
		Summary("List published posts").
		Response(200, []models.Blog{}, "Published posts, newest first")
	docs.Route("GET /api/v1/blogs/:slug").
		Summary("Read one published post").
		Response(200, models.Blog{}, "The post").
		Response(404, handlers.ErrorResponse{}, "No published post with that slug")
	docs.Route("GET /api/v1/admin/blogs").
		Summary("List posts, including drafts").
		Response(200, []models.Blog{}, "Every post, newest first")
	docs.Route("GET /api/v1/admin/blogs/:id").
		Summary("Get one post").
		Response(200, models.Blog{}, "The post").
		Response(404, handlers.ErrorResponse{}, "Not found")
	docs.Route("POST /api/v1/admin/blogs").
		Summary("Create a post").
		RequestBody(handlers.CreateBlogRequest{}).
		Response(201, models.Blog{}, "The new post")
	docs.Route("PUT /api/v1/admin/blogs/:id").
		Summary("Update a post").
		RequestBody(handlers.UpdateBlogRequest{}).
		Response(200, models.Blog{}, "The updated post")
	docs.Route("DELETE /api/v1/admin/blogs/:id").
		Summary("Delete a post").
		Response(200, handlers.MessageResponse{}, "Deleted")

	// ── Support tickets ────────────────────────────────────────────
	docs.Route("GET /api/v1/tickets").
		Summary("List tickets").
		Response(200, []models.Ticket{}, "Tickets visible to you, newest first")
	docs.Route("GET /api/v1/tickets/:id").
		Summary("One ticket and its thread").
		Response(200, models.Ticket{}, "The ticket, with replies").
		Response(404, handlers.ErrorResponse{}, "Not found")
	docs.Route("POST /api/v1/tickets").
		Summary("Open a ticket").
		RequestBody(handlers.CreateTicketRequest{}).
		Response(201, models.Ticket{}, "The new ticket")
	docs.Route("POST /api/v1/tickets/:id/reply").
		Summary("Reply to a ticket").
		RequestBody(handlers.TicketReplyRequest{}).
		Response(201, models.TicketReply{}, "The reply")
	docs.Route("PATCH /api/v1/tickets/:id/assign").
		Summary("Assign a ticket to a staff member").
		RequestBody(handlers.AssignTicketRequest{}).
		Response(200, models.Ticket{}, "The ticket, reassigned")
	docs.Route("PATCH /api/v1/tickets/:id/close").
		Summary("Close a ticket").
		Response(200, models.Ticket{}, "The closed ticket")
	docs.Route("PATCH /api/v1/tickets/:id/reopen").
		Summary("Reopen a closed ticket").
		Response(200, models.Ticket{}, "The reopened ticket")

	// ── Roles ──────────────────────────────────────────────────────
	docs.Route("GET /api/v1/roles/:id").
		Summary("Get one role").
		Response(200, models.Role{}, "The role and its grants").
		Response(404, handlers.ErrorResponse{}, "Not found")
	docs.Route("POST /api/v1/roles").
		Summary("Create a role").
		RequestBody(handlers.RoleRequest{}).
		Response(201, models.Role{}, "The new role")
	docs.Route("PUT /api/v1/roles/:id").
		Summary("Update a role's grants").
		RequestBody(handlers.RoleRequest{}).
		Response(200, models.Role{}, "The updated role")
	docs.Route("DELETE /api/v1/roles/:id").
		Summary("Delete a role").
		Response(200, handlers.MessageResponse{}, "Deleted").
		Response(400, handlers.ErrorResponse{}, "Built-in roles cannot be deleted")

	// ── Users ──────────────────────────────────────────────────────
	docs.Route("POST /api/v1/users").
		Summary("Create a user").
		RequestBody(handlers.CreateUserRequest{}).
		Response(201, models.User{}, "The new user").
		Response(409, handlers.ErrorResponse{}, "That email is already registered")
	docs.Route("PUT /api/v1/users/:id").
		Summary("Update a user").
		RequestBody(handlers.UpdateUserRequest{}).
		Response(200, models.User{}, "The updated user")
	docs.Route("DELETE /api/v1/users/:id").
		Summary("Delete a user").
		Response(200, handlers.MessageResponse{}, "Deleted")
	docs.Route("PUT /api/v1/users/:id/roles").
		Summary("Replace a user's roles").
		RequestBody(handlers.AssignRolesRequest{}).
		Response(200, models.User{}, "The user, with the new roles")

	// ── Uploads ────────────────────────────────────────────────────
	docs.Route("GET /api/v1/uploads/:id").
		Summary("Get one stored file").
		Response(200, models.Upload{}, "The file record").
		Response(404, handlers.ErrorResponse{}, "Not found")
	docs.Route("GET /api/v1/uploads/stats").
		Summary("Storage totals").
		Response(200, handlers.MessageResponse{}, "Count and bytes, by type")
	docs.Route("DELETE /api/v1/uploads/:id").
		Summary("Delete a stored file").
		Response(200, handlers.MessageResponse{}, "Removed from storage and from the database")

	// ── Profile ────────────────────────────────────────────────────
	docs.Route("PUT /api/v1/profile").
		Summary("Update your own profile").
		RequestBody(handlers.UpdateProfileRequest{}).
		Response(200, models.User{}, "The updated profile")
	docs.Route("DELETE /api/v1/profile").
		Summary("Delete your own account").
		Response(200, handlers.MessageResponse{}, "Account deleted and every session signed out")

	// ── SSO ────────────────────────────────────────────────────────
	docs.Route("POST /api/v1/sso/connections").
		Summary("Add an SSO connection").
		RequestBody(handlers.SSOConnectionRequest{}).
		Response(201, models.SSOConnection{}, "The new connection")
	docs.Route("PUT /api/v1/sso/connections/:id").
		Summary("Update an SSO connection").
		RequestBody(handlers.SSOConnectionRequest{}).
		Response(200, models.SSOConnection{}, "The updated connection")
	docs.Route("DELETE /api/v1/sso/connections/:id").
		Summary("Delete an SSO connection").
		Response(200, handlers.MessageResponse{}, "Deleted")
	docs.Route("GET /api/v1/sso/connections/:id/test").
		Summary("Check a connection's discovery document and credentials").
		Response(200, handlers.MessageResponse{}, "What the provider returned").
		Response(502, handlers.ErrorResponse{}, "The provider could not be reached")

	// ── Access reviews ─────────────────────────────────────────────
	docs.Route("POST /api/v1/access-reviews").
		Summary("Open an access review").
		RequestBody(handlers.OpenReviewRequest{}).
		Response(201, models.AccessReview{}, "The new review, with one item per grant")
	docs.Route("POST /api/v1/access-reviews/:id/complete").
		Summary("Close a review and apply its decisions").
		Response(200, models.AccessReview{}, "The completed review")
	docs.Route("POST /api/v1/access-reviews/:id/items/:itemId/decision").
		Summary("Keep or revoke one grant").
		RequestBody(handlers.ReviewDecisionRequest{}).
		Response(200, models.AccessReviewItem{}, "The decided item")

	// ── Notifications, dashboard, activity ─────────────────────────
	docs.Route("POST /api/v1/notifications/:id/read").
		Summary("Mark one notification read").
		Response(200, models.Notification{}, "The notification")
	docs.Route("POST /api/v1/notifications/read-all").
		Summary("Mark every notification read").
		Response(200, handlers.MessageResponse{}, "How many were marked")
	docs.Route("GET /api/v1/dashboard-layout").
		Summary("Your saved dashboard arrangement").
		Response(200, models.DashboardLayout{}, "Widget order and visibility")
	docs.Route("PUT /api/v1/dashboard-layout").
		Summary("Save your dashboard arrangement").
		RequestBody(handlers.DashboardLayoutRequest{}).
		Response(200, models.DashboardLayout{}, "The saved layout")
	docs.Route("GET /api/v1/user-activity").
		Summary("Your own recent activity").
		Response(200, []models.UserActivity{}, "Newest first")
	docs.Route("GET /api/v1/user-activity/stats").
		Summary("Your activity, aggregated").
		Response(200, handlers.MessageResponse{}, "Counts by day and by action")

	// ── Imports, backups, audit, system ────────────────────────────
	docs.Route("GET /api/v1/imports/:id").
		Summary("Progress of a CSV import").
		Response(200, models.ImportJob{}, "Row counts, status and any per-row errors")
	docs.Route("GET /api/v1/backups/:id/download").
		Summary("Download a backup archive").
		Response(200, handlers.MessageResponse{}, "A ZIP attachment: one CSV per table, a dump.sql and a manifest").
		Response(404, handlers.ErrorResponse{}, "That archive is no longer on disk")
	docs.Route("PUT /api/v1/backup-settings").
		Summary("Change the backup schedule").
		RequestBody(handlers.BackupSettingsRequest{}).
		Response(200, models.BackupSchedule{}, "The saved schedule")
	docs.Route("GET /api/v1/audit/ocsf").
		Summary("The activity log as OCSF events").
		Response(200, handlers.MessageResponse{}, "OCSF-shaped records, for a SIEM to ingest")
	docs.Route("GET /api/v1/system/modules").
		Summary("Which optional modules this build has").
		Response(200, handlers.MessageResponse{}, "One flag per module, so a client can hide what is absent")
	docs.Route("GET /api/health").
		Summary("Liveness and dependency check").
		Response(200, handlers.MessageResponse{}, "Database, cache and queue status")

	// ── Offline sync ───────────────────────────────────────────────
	docs.Route("GET /api/v1/sync/pull").
		Summary("Changes since a cursor").
		Response(200, handlers.MessageResponse{}, "Rows created, updated or deleted since the cursor, plus the next cursor")
	docs.Route("POST /api/v1/sync/push").
		Summary("Send changes made while offline").
		RequestBody(handlers.SyncPushRequest{}).
		Response(200, handlers.MessageResponse{}, "Per-change accepted or conflicted, with the winning row")
	docs.Route("GET /api/v1/sync/policy").
		Summary("How each model behaves offline").
		Response(200, handlers.SyncPolicyResponse{}, "Sync mode, conflict strategy, synced fields and offline age limit per model")

	// ── Public form sharing ────────────────────────────────────────
	docs.Route("GET /api/v1/public/forms/:token").
		Summary("Fetch a shared form's fields").
		Response(200, handlers.MessageResponse{}, "The form definition, or a password prompt").
		Response(404, handlers.ErrorResponse{}, "No share with that token, or it has expired")
	docs.Route("POST /api/v1/public/forms/:token/submit").
		Summary("Submit a shared form").
		RequestBody(handlers.PublicFormSubmitRequest{}).
		Response(201, handlers.MessageResponse{}, "Recorded").
		Response(401, handlers.ErrorResponse{}, "The share is password-protected")

	// ── Admin: flags, shares, jobs, webhooks, summaries ────────────
	docs.Route("GET /api/v1/admin/flags").
		Summary("List feature flags").
		Response(200, []models.FeatureFlag{}, "Every flag and its rules")
	docs.Route("POST /api/v1/admin/flags").
		Summary("Create a feature flag").
		RequestBody(handlers.FeatureFlagRequest{}).
		Response(201, models.FeatureFlag{}, "The new flag")
	docs.Route("PUT /api/v1/admin/flags/:id").
		Summary("Update a flag's rules").
		RequestBody(handlers.FeatureFlagRequest{}).
		Response(200, models.FeatureFlag{}, "The updated flag")
	docs.Route("DELETE /api/v1/admin/flags/:id").
		Summary("Delete a feature flag").
		Response(200, handlers.MessageResponse{}, "Deleted")
	docs.Route("GET /api/v1/admin/flags/:id/exposures").
		Summary("Who has seen a flag, and which way it evaluated").
		Response(200, []models.FlagExposure{}, "Newest first")
	docs.Route("GET /api/v1/admin/form-shares").
		Summary("List public form shares").
		Response(200, []models.FormShare{}, "Every share, newest first")
	docs.Route("POST /api/v1/admin/form-shares").
		Summary("Create a public form share").
		RequestBody(handlers.CreateFormShareRequest{}).
		Response(201, models.FormShare{}, "The new share, with its public token")
	docs.Route("PATCH /api/v1/admin/form-shares/:id").
		Summary("Update a share").
		RequestBody(handlers.UpdateFormShareRequest{}).
		Response(200, models.FormShare{}, "The updated share")
	docs.Route("DELETE /api/v1/admin/form-shares/:id").
		Summary("Delete a share").
		Response(200, handlers.MessageResponse{}, "Deleted; the public link stops working")
	docs.Route("GET /api/v1/admin/form-shares/resources").
		Summary("Resources that can be shared publicly").
		Response(200, handlers.MessageResponse{}, "Only what services/form_share_dispatch.go registers")
	docs.Route("GET /api/v1/admin/form-shares/resources/:resource/fields").
		Summary("A shareable resource's fields").
		Response(200, []services.PublicFieldInfo{}, "Field names, types and whether each is required")
	docs.Route("GET /api/v1/admin/form-submissions").
		Summary("Submissions received through public shares").
		Response(200, []models.FormSubmission{}, "Newest first")
	docs.Route("GET /api/v1/admin/jobs/:status").
		Summary("Background jobs in one state").
		Response(200, handlers.MessageResponse{}, "Jobs with that status").
		Response(503, handlers.ErrorResponse{}, "Redis is not configured")
	docs.Route("POST /api/v1/admin/jobs/:id/retry").
		Summary("Retry a failed job").
		Response(200, handlers.MessageResponse{}, "Requeued")
	docs.Route("DELETE /api/v1/admin/jobs/queue/:queue").
		Summary("Empty a queue").
		Response(200, handlers.MessageResponse{}, "How many jobs were discarded")
	docs.Route("GET /api/v1/admin/webhooks").
		Summary("Outbound webhook deliveries").
		Response(200, []models.WebhookEvent{}, "Each attempt with its response status")
	docs.Route("POST /api/v1/admin/webhooks/:id/replay").
		Summary("Send a webhook again").
		Response(202, models.WebhookEvent{}, "The new delivery attempt")
	docs.Route("GET /api/v1/admin/observability/summary").
		Summary("Throughput, latency and error rate").
		Response(200, handlers.MessageResponse{}, "What the observability page charts")
	docs.Route("GET /api/v1/admin/security/summary").
		Summary("Failed logins, lockouts and recent security events").
		Response(200, handlers.MessageResponse{}, "What the security page shows")
	docs.Route("GET /api/v1/admin/dashboard/chart/:resource").
		Summary("A resource's rows over time").
		Response(200, services.ChartResult{}, "Buckets ready to plot")
	docs.Route("GET /api/v1/admin/dashboard/resource-stats/:resource").
		Summary("Totals and trend for one resource").
		Response(200, services.ResourceStats{}, "Count, change against the previous period, and a sparkline")

	// ── Auth: the endpoints a browser redirects to ─────────────────
	// These return redirects rather than JSON, which is worth stating
	// explicitly: a client that follows them expecting a body gets HTML.
	docs.Route("POST /api/v1/auth/totp/setup").
		Summary("Begin two-factor enrolment").
		Response(200, handlers.MessageResponse{}, "A secret, an otpauth:// URI and a QR as a PNG data URI").
		Response(409, handlers.ErrorResponse{}, "Two-factor is already on")
	docs.Route("POST /api/v1/auth/totp/enable").
		Summary("Finish enrolment by proving you can generate a code").
		RequestBody(handlers.EnableTOTPRequest{}).
		Response(200, handlers.MessageResponse{}, "Enabled, with the one-time backup codes").
		Response(400, handlers.ErrorResponse{}, "That code did not match the secret")
	docs.Route("POST /api/v1/auth/totp/disable").
		Summary("Turn two-factor off").
		RequestBody(handlers.DisableTOTPRequest{}).
		Response(200, handlers.MessageResponse{}, "Disabled; trusted devices are revoked too").
		Response(401, handlers.ErrorResponse{}, "Wrong password")
	docs.Route("POST /api/v1/auth/totp/verify").
		Summary("Answer a two-factor challenge").
		RequestBody(handlers.VerifyTOTPRequest{}).
		Response(200, handlers.AuthResponse{}, "Signed in").
		Response(401, handlers.ErrorResponse{}, "Wrong or expired code")
	docs.Route("POST /api/v1/auth/totp/backup-codes/verify").
		Summary("Answer a challenge with a backup code").
		RequestBody(handlers.VerifyBackupCodeRequest{}).
		Response(200, handlers.AuthResponse{}, "Signed in; that code is now spent").
		Response(401, handlers.ErrorResponse{}, "Unknown or already-used code")
	docs.Route("POST /api/v1/auth/totp/backup-codes").
		Summary("Replace your backup codes").
		Response(200, handlers.MessageResponse{}, "A fresh set; the old ones stop working")

	docs.Route("GET /api/v1/auth/permissions").
		Summary("What the signed-in user may do").
		Response(200, handlers.MessageResponse{}, "The permission keys granted by this user's roles")
	docs.Route("DELETE /api/v1/auth/totp/trusted-devices").
		Summary("Revoke every trusted device").
		Response(200, handlers.MessageResponse{}, "All devices revoked; each will be asked for a code next time")
	docs.Route("DELETE /api/v1/auth/totp/trusted-devices/:id").
		Summary("Revoke one trusted device").
		Response(200, handlers.MessageResponse{}, "That device is no longer trusted")
	docs.Route("GET /api/v1/auth/oauth/:provider").
		Summary("Start a social sign-in").
		Response(302, handlers.MessageResponse{}, "Redirects to the provider")
	docs.Route("GET /api/v1/auth/oauth/:provider/callback").
		Summary("Return from a social sign-in").
		Response(302, handlers.MessageResponse{}, "Redirects to the frontend with a token pair")
	docs.Route("GET /api/v1/auth/sso/:slug").
		Summary("Start an OIDC SSO sign-in").
		Response(302, handlers.MessageResponse{}, "Redirects to the identity provider")
	docs.Route("GET /api/v1/auth/sso/:slug/callback").
		Summary("Return from an OIDC SSO sign-in").
		Response(302, handlers.MessageResponse{}, "Redirects to the frontend with a token pair")
	docs.Route("GET /api/v1/auth/saml/:slug").
		Summary("Start a SAML sign-in").
		Response(302, handlers.MessageResponse{}, "Redirects to the identity provider")
	docs.Route("GET /api/v1/auth/saml/:slug/metadata").
		Summary("This service's SAML metadata").
		Response(200, handlers.MessageResponse{}, "XML to hand to the identity provider")
	docs.Route("POST /api/v1/auth/saml/:slug/acs").
		Summary("SAML assertion consumer").
		Response(302, handlers.MessageResponse{}, "Redirects to the frontend with a token pair")
	docs.Route("POST /api/v1/auth/sso/discover").
		Summary("Find the SSO connection for an email domain").
		RequestBody(handlers.SSODiscoverRequest{}).
		Response(200, handlers.MessageResponse{}, "The connection to use, or nothing if the domain is not enrolled")

	// ── AI ─────────────────────────────────────────────────────────
	docs.Route("POST /api/v1/ai/chat").
		Summary("Chat completion").
		RequestBody(handlers.ChatRequest{}).
		Response(200, ai.CompletionResponse{}, "The assistant's reply")
	docs.Route("POST /api/v1/ai/complete").
		Summary("Single-turn completion").
		RequestBody(handlers.CompleteRequest{}).
		Response(200, ai.CompletionResponse{}, "The completion")
	docs.Route("POST /api/v1/ai/stream").
		Summary("Streaming chat completion").
		RequestBody(handlers.ChatRequest{}).
		Response(200, ai.CompletionResponse{}, "Server-sent events, one token per chunk")

	// ── Realtime and inbound webhooks ──────────────────────────────
	docs.Route("GET /api/ws").
		Summary("Realtime WebSocket").
		Response(101, handlers.MessageResponse{}, "Switches protocols; then pushes notification and record events")
	docs.Route("POST /webhooks/:provider").
		Summary("Inbound webhook from a third party").
		Response(200, handlers.MessageResponse{}, "Accepted").
		Response(401, handlers.ErrorResponse{}, "Signature check failed")

	// grit:docs:routes — "grit generate resource" registers each resource here.
	// grit:docs:routes:end
	log.Println("API docs available at /docs")
}
`
}
