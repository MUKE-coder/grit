package scaffold

import (
	"path/filepath"
	"strings"

	"github.com/MUKE-coder/grit/v3/internal/manifest"
)

// ─── M36: an erasure needs a reason, and its audit row is checked ────────────
//
// The handler bound the request body with the error thrown away, so a POST with
// no body at all erased an account and wrote "" into the deletion journal as the
// compliance reason. The activity row was written with a bare Create whose
// error, when it failed, was invisible: the caller was told the erasure was
// clean and the one record an auditor asks for was never there.
//
// gdpr.go is not a framework-owned file, so upgrade patches the two functions
// rather than rewriting it. Every step checks the text it is about to replace,
// and a file that no longer looks like the one Grit wrote is left alone with a
// line saying so.

const gdprEraseOldBind = "	var req EraseRequest\n	_ = c.ShouldBindJSON(&req)\n"

const gdprEraseSignature = "func (h *GDPRHandler) Erase(c *gin.Context) {"

const gdprExportGuardOld = `	targetID := c.Param("id")
	callerID, _ := c.Get("user_id")
	callerRole, _ := c.Get("user_role")
	if fmt.Sprint(callerID) != targetID && fmt.Sprint(callerRole) != "ADMIN" {`

const gdprExportGuardNew = `	targetID := c.Param("id")
	callerRole, _ := c.Get("user_role")
	if authz.CurrentUserID(c) != targetID && fmt.Sprint(callerRole) != "ADMIN" {`

// logActivityErrFunc is added to services/activity.go beside LogActivity.
const logActivityErrFunc = `// LogActivityErr is LogActivity for the handful of events where losing the row
// matters more than the noise: a GDPR erasure, a key revocation, anything an
// auditor will later ask you to produce. It writes the same row and hands back
// the error instead of logging it, so the caller can tell the client the audit
// trail is incomplete.
func LogActivityErr(db *gorm.DB, c *gin.Context, args ActivityArgs) error {
	row := activityRow(c, args)
	return db.Create(&row).Error
}
`

// repairGDPRErase applies M36 to an existing project.
func repairGDPRErase(root string, opts Options) error {
	apiRoot := opts.APIRoot(root)
	handler := filepath.Join(apiRoot, "internal", "handlers", "gdpr.go")
	if !fileExists(handler) {
		return nil
	}
	m, err := manifest.Load(root)
	if err != nil {
		return err
	}

	if activity := filepath.Join(apiRoot, "internal", "services", "activity.go"); fileExists(activity) {
		if err := repairSourceFile(root, m, activity, repairLogActivityErrSource); err != nil {
			return err
		}
	}

	module := opts.Module()
	return repairSourceFile(root, m, handler, func(src string) (string, []string, []string) {
		return repairGDPREraseSource(src, module)
	})
}

// repairLogActivityErrSource adds LogActivityErr next to LogActivity.
func repairLogActivityErrSource(src string) (string, []string, []string) {
	if strings.Contains(src, "func LogActivityErr(") {
		return src, nil, nil
	}
	const anchor = "func LogActivity(db *gorm.DB, c *gin.Context, args ActivityArgs) {"
	if !strings.Contains(src, anchor) || !strings.Contains(src, "func activityRow(c *gin.Context, args ActivityArgs)") {
		return src, nil, nil
	}
	// Above LogActivity's doc comment, so the two read as a pair.
	at := strings.Index(src, anchor)
	for at > 0 {
		prev := strings.LastIndex(src[:at-1], "\n") + 1
		if !strings.HasPrefix(strings.TrimLeft(src[prev:at-1], "\t "), "//") {
			break
		}
		at = prev
	}
	return src[:at] + logActivityErrFunc + "\n" + src[at:],
		[]string{"services/activity.go can write an audit row and report a failure to write it"}, nil
}

// repairGDPREraseSource rewrites Erase, and the identity check in Export.
func repairGDPREraseSource(src, module string) (string, []string, []string) {
	if !strings.Contains(src, gdprEraseSignature) {
		return src, nil, nil
	}
	if strings.Contains(src, `binding:"required,min=3,max=500"`) {
		return src, nil, nil
	}
	if !strings.Contains(src, gdprEraseOldBind) {
		return src, nil, []string{"handlers/gdpr.go is not the file Grit wrote: check the error from ShouldBindJSON in Erase, or an erasure runs with no compliance reason"}
	}

	body := strings.ReplaceAll(gdprEraseFunc, "~", "`")
	out, ok := replaceGoFunc(src, gdprEraseSignature, body)
	if !ok {
		return src, nil, []string{"handlers/gdpr.go is not the file Grit wrote: check the error from ShouldBindJSON in Erase, or an erasure runs with no compliance reason"}
	}

	// The request type carries the binding tag that makes the reason required.
	oldReq := "type EraseRequest struct {\n\tReason string `json:\"reason\"`\n}"
	if strings.Contains(out, oldReq) {
		out = strings.Replace(out, oldReq, strings.ReplaceAll(gdprEraseRequestType, "~", "`"), 1)
	}

	changes := []string{"an erasure is refused with 422 unless it carries a compliance reason, and a lost audit row is reported"}
	if strings.Contains(out, gdprExportGuardOld) {
		out = strings.Replace(out, gdprExportGuardOld, gdprExportGuardNew, 1)
		changes = append(changes, "the GDPR export compares the caller with authz.CurrentUserID")
	}

	if strings.Contains(out, "authz.") && !strings.Contains(out, "\""+module+"/internal/authz\"") {
		added := false
		if out, added = addImportGroup(out, module+"/internal/authz"); !added {
			return src, nil, []string{"handlers/gdpr.go has no import block this can add internal/authz to; left alone"}
		}
	}
	if !strings.Contains(out, "\t\"errors\"\n") {
		if withErrors, added := addImportGroup(out, "errors"); added {
			out = withErrors
		}
	}
	for _, path := range []string{"log", module + "/internal/models"} {
		out, _ = dropUnusedImport(out, path)
	}
	return out, changes, nil
}

const gdprEraseRequestType = `type EraseRequest struct {
	// Required, not optional. An erasure is irreversible and its journal row is
	// what an auditor reads a year later; an empty reason tells them nothing
	// about why a person's data was destroyed. The previous handler ignored the
	// bind error, so a request with no body at all, or with the field
	// misspelled, erased an account and recorded no reason for it.
	Reason string ~json:"reason" binding:"required,min=3,max=500"~
}`

const gdprEraseFunc = `// Erase fulfils a right-to-erasure request (admin only). It refuses to let an
// admin erase themselves — that would revoke their own access mid-request and
// orphan the operation.
func (h *GDPRHandler) Erase(c *gin.Context) {
	targetID := c.Param("id")
	callerID := authz.CurrentUserID(c)
	callerEmail, _ := c.Get("user_email")

	if callerID == targetID {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "SELF_ERASE", "message": "you cannot erase your own account here"}})
		return
	}

	var req EraseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respond.Validation(c, "A reason is required for an erasure", map[string]string{
			"reason": "Give the compliance reason for this erasure (3 to 500 characters). It is written to the deletion journal.",
		})
		return
	}

	journal, err := services.EraseUser(h.DB, targetID, callerID, fmt.Sprint(callerEmail), req.Reason)
	if err != nil {
		if errors.Is(err, services.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": gin.H{"code": "NOT_FOUND", "message": "user not found"}})
			return
		}
		respond.ServerError(c, "ERASE_FAILED", err, "Internal server error")
		return
	}

	// Record the erasure in the semantic activity log too, so it shows up in the
	// dashboard and flows out through the OCSF/SIEM export.
	//
	// LogActivityErr rather than a bare Create: the erasure itself has already
	// happened, so a lost audit row must not fail the request, but it is the
	// record an auditor asks for. The response says so rather than reporting a
	// clean erasure that left no trace.
	activityLogged := true
	if err := services.LogActivityErr(h.DB, c, services.ActivityArgs{
		UserID:       callerID,
		Action:       "user.gdpr_erase",
		Severity:     "warn",
		Summary:      fmt.Sprintf("Erased all personal data for user %s (%d records)", targetID, journal.RecordsAffected),
		ResourceType: "user",
		ResourceID:   targetID,
	}); err != nil {
		activityLogged = false
	}

	if !activityLogged {
		c.JSON(http.StatusOK, gin.H{
			"data":    journal,
			"message": "User data erased, but the activity row could not be written",
			"meta":    gin.H{"activity_logged": false},
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": journal, "message": "User data erased"})
}
`
