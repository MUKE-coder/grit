package scaffold

import (
	"path/filepath"
	"strings"

	"github.com/MUKE-coder/grit/v3/internal/manifest"
)

// Contact-app review L29, L31 and L32, in the files upgrade does not deliver
// whole.
//
// L29: three handlers took the signed-in user out of the gin context with a
// bare type assertion. handlers/totp.go and handlers/upload.go arrive whole,
// so dashboard_layout.go is repaired here, and upload.go only where the
// manifest guard held an edited copy back. L32's ticket
// constants live in models/ticket.go, which upgrade does not write, and the
// ticket service and handler it does write refer to them.
//
// L31: handlers/upload.go arrives whole through writeStorageFiles, and no longer
// reads UPLOAD_ALLOWED_MIME itself, so config.go and routes.go are repaired to
// hand it the list. Without that, a project that set the variable would lose
// the types it added.
//
// L30 and L33 have no repair. Dead code left in an older project is harmless,
// and the files of it upgrade owns (webhooks.go, sso.go, saml.go) arrive whole
// behind the manifest guard.

// ─── L29: the caller is checked, not asserted ────────────────────────────────

const (
	dashboardUserIDOld = "\tuserID, ok := c.Get(\"user_id\")\n\tif !ok {\n"
	dashboardUserIDNew = "\tuserID := authz.CurrentUserID(c)\n\tif userID == \"\" {\n"

	uploadCreateGuardOld = "func (h *UploadHandler) Create(c *gin.Context) {\n" +
		"\tif h.Storage == nil {\n" +
		"\t\trespond.Fail(c, respond.CodeStorageUnavailable, \"File storage is not configured\")\n" +
		"\t\treturn\n" +
		"\t}\n"
	uploadCreateGuardNew = uploadCreateGuardOld + `
	// Read before the body is. The route is behind the auth middleware, and a
	// handler that assumed so, rather than checking, panicked on every request
	// the day it was mounted anywhere else.
	userID := authz.CurrentUserID(c)
	if userID == "" {
		respond.Fail(c, respond.CodeUnauthorized, "Not signed in")
		return
	}
`
	uploadUserIDOld = "\tuserID, _ := c.Get(\"user_id\")\n\n\tupload := models.Upload{\n"
	uploadUserIDNew = "\tupload := models.Upload{\n"
)

// ─── L31: the upload allowlist is built from config ──────────────────────────

const (
	configUploadMIMEField = `	// UploadAllowedMIME are MIME types added to the upload allowlist, from
	// UPLOAD_ALLOWED_MIME (comma separated). They extend the baseline in
	// handlers/upload.go; a field narrows it with accepts.
	UploadAllowedMIME []string
`
	configUploadMIMELoad = "\t\tUploadAllowedMIME: splitCSV(getEnv(\"UPLOAD_ALLOWED_MIME\", \"\")),\n"

	configCORSFieldAnchor = "\tCORSOrigins []string\n\n"
	configCORSLoadAnchor  = "\t\tCORSOrigins: strings.Split(getEnv(\"CORS_ORIGINS\", \"http://localhost:3000,http://localhost:3001\"), \",\"),\n\n"

	routesUploadHandlerOld = "\tuploadHandler := &handlers.UploadHandler{\n" +
		"\t\tDB:      db,\n" +
		"\t\tStorage: svc.Storage,\n" +
		"\t\tJobs:    svc.Jobs,\n" +
		"\t}\n"
	routesUploadHandlerNew = "\tuploadHandler := &handlers.UploadHandler{\n" +
		"\t\tDB:          db,\n" +
		"\t\tStorage:     svc.Storage,\n" +
		"\t\tJobs:        svc.Jobs,\n" +
		"\t\tAllowedMIME: handlers.UploadMIMEAllowlist(cfg.UploadAllowedMIME),\n" +
		"\t}\n"
)

// repairUploadMIMEConfigSource adds UploadAllowedMIME to config.Config.
func repairUploadMIMEConfigSource(src string) (string, []string, []string) {
	if strings.Contains(src, "UploadAllowedMIME") || !strings.Contains(src, "type Config struct {") {
		return src, nil, nil
	}
	if strings.Count(src, configCORSFieldAnchor) != 1 || strings.Count(src, configCORSLoadAnchor) != 1 ||
		!strings.Contains(src, "func splitCSV(") {
		return src, nil, []string{uploadMIMEManualWarning}
	}
	out := strings.Replace(src, configCORSFieldAnchor, configCORSFieldAnchor+configUploadMIMEField+"\n", 1)
	out = strings.Replace(out, configCORSLoadAnchor, configCORSLoadAnchor+configUploadMIMELoad+"\n", 1)
	return out, []string{"UPLOAD_ALLOWED_MIME is read into Config.UploadAllowedMIME"}, nil
}

// repairUploadMIMERoutesSource hands the allowlist to the upload handler.
func repairUploadMIMERoutesSource(src string) (string, []string, []string) {
	if strings.Contains(src, "UploadMIMEAllowlist(") || !strings.Contains(src, "&handlers.UploadHandler{") {
		return src, nil, nil
	}
	if strings.Count(src, routesUploadHandlerOld) != 1 {
		return src, nil, []string{uploadMIMEManualWarning}
	}
	out := strings.Replace(src, routesUploadHandlerOld, routesUploadHandlerNew, 1)
	return out, []string{"the upload handler takes its allowlist from config, built once at startup"}, nil
}

const uploadMIMEManualWarning = "not the file Grit wrote: handlers/upload.go no longer reads UPLOAD_ALLOWED_MIME itself, so until the upload handler is built with AllowedMIME: handlers.UploadMIMEAllowlist(cfg.UploadAllowedMIME), with UploadAllowedMIME read into config, only the built-in types upload"

// ─── L32: the ticket vocabulary is named once ────────────────────────────────

const ticketConstantsBlock = `
// A ticket's status and priority. The column is a string, so these are the
// values it holds; the admin panel offers the same four priorities.
const (
	TicketStatusOpen   = "open"
	TicketStatusClosed = "closed"

	TicketPriorityLow      = "low"
	TicketPriorityMedium   = "medium"
	TicketPriorityHigh     = "high"
	TicketPriorityCritical = "critical"
)
`

const (
	ticketDefaultsOld = "\tif t.Status == \"\" {\n\t\tt.Status = \"open\"\n\t}\n\tif t.Priority == \"\" {\n\t\tt.Priority = \"medium\"\n\t}\n"
	ticketDefaultsNew = "\tif t.Status == \"\" {\n\t\tt.Status = TicketStatusOpen\n\t}\n\tif t.Priority == \"\" {\n\t\tt.Priority = TicketPriorityMedium\n\t}\n"
	ticketStructLine  = "type Ticket struct {\n"
)

// repairCodeQualityLows applies L29, L31 and L32 to an existing project.
func repairCodeQualityLows(root string, opts Options) error {
	apiRoot := opts.APIRoot(root)
	m, err := manifest.Load(root)
	if err != nil {
		return err
	}
	for _, r := range []struct {
		path string
		fn   func(string) (string, []string, []string)
	}{
		{filepath.Join(apiRoot, "internal", "handlers", "dashboard_layout.go"), repairDashboardLayoutUserSource},
		{filepath.Join(apiRoot, "internal", "handlers", "upload.go"), repairUploadUserSource},
		{filepath.Join(apiRoot, "internal", "models", "ticket.go"), repairTicketConstantsSource},
	} {
		if !fileExists(r.path) {
			continue
		}
		if err := repairSourceFile(root, m, r.path, r.fn); err != nil {
			return err
		}
	}

	// L31. Only once upload.go is the one that takes its allowlist from the
	// handler: routes.go naming UploadMIMEAllowlist over an older upload.go
	// would not compile. The config comes first, since routes.go reads it.
	if !fileContains(filepath.Join(apiRoot, "internal", "handlers", "upload.go"), "func UploadMIMEAllowlist(") {
		return nil
	}
	config := filepath.Join(apiRoot, "internal", "config", "config.go")
	if err := repairSourceFile(root, m, config, repairUploadMIMEConfigSource); err != nil {
		return err
	}
	if !fileContains(config, "UploadAllowedMIME") {
		return nil
	}
	routes := filepath.Join(apiRoot, "internal", "routes", "routes.go")
	if !fileExists(routes) {
		return nil
	}
	return repairSourceFile(root, m, routes, repairUploadMIMERoutesSource)
}

// repairDashboardLayoutUserSource reads the caller through authz.CurrentUserID.
func repairDashboardLayoutUserSource(src string) (string, []string, []string) {
	if !strings.Contains(src, "userID.(string)") {
		return src, nil, nil
	}
	module := authModule(src)
	if strings.Count(src, dashboardUserIDOld) != 2 || module == "" {
		return src, nil, []string{"not the file Grit wrote: read the caller with authz.CurrentUserID(c) and answer 401 when it is empty, or the route panics when mounted without the auth middleware"}
	}
	out := strings.ReplaceAll(src, dashboardUserIDOld, dashboardUserIDNew)
	out = strings.ReplaceAll(out, "userID.(string)", "userID")
	out, ok := addImportBefore(out, module+"/internal/authz", "\t\""+module+"/internal/models\"\n")
	if !ok {
		return src, nil, []string{"could not add the internal/authz import: read the caller with authz.CurrentUserID(c) by hand"}
	}
	return out, []string{"the dashboard layout answers 401 instead of panicking when there is no signed-in user"}, nil
}

// repairUploadUserSource checks the caller before an upload is read.
func repairUploadUserSource(src string) (string, []string, []string) {
	if !strings.Contains(src, "userID.(string)") {
		return src, nil, nil
	}
	if strings.Count(src, uploadCreateGuardOld) != 1 || strings.Count(src, uploadUserIDOld) != 1 ||
		strings.Count(src, "userID.(string)") != 1 || !strings.Contains(src, "/internal/authz\"\n") {
		return src, nil, []string{"not the file Grit wrote: read the caller with authz.CurrentUserID(c) at the top of Create and answer 401 when it is empty, or an upload panics when the route is mounted without the auth middleware"}
	}
	out := strings.Replace(src, uploadCreateGuardOld, uploadCreateGuardNew, 1)
	out = strings.Replace(out, uploadUserIDOld, uploadUserIDNew, 1)
	out = strings.Replace(out, "userID.(string)", "userID", 1)
	return out, []string{"an upload with no signed-in user is refused with 401 before it is stored, instead of panicking after"}, nil
}

// repairTicketConstantsSource names the ticket statuses and priorities. The
// ticket service and handler upgrade writes use them, so the block is added
// even to a model somebody has edited, as long as it does not declare them.
func repairTicketConstantsSource(src string) (string, []string, []string) {
	if strings.Contains(src, "TicketStatusOpen") {
		return src, nil, nil
	}
	start := strings.Index(src, "\nimport (\n")
	end := -1
	if start >= 0 {
		end = strings.Index(src[start:], "\n)\n")
	}
	if strings.Count(src, ticketStructLine) != 1 || end < 0 {
		return src, nil, []string{"no Ticket struct here: declare TicketStatusOpen, TicketStatusClosed and the four TicketPriority constants in package models, which services/ticket.go uses"}
	}
	// After the imports, so the Ticket struct keeps its doc comment.
	at := start + end + len("\n)\n")
	out := src[:at] + ticketConstantsBlock + src[at:]
	// The literals in BeforeCreate are Grit's, and only replaced when they are.
	out = strings.Replace(out, ticketDefaultsOld, ticketDefaultsNew, 1)
	return out, []string{"ticket statuses and priorities are named constants in models"}, nil
}
