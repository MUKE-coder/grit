package scaffold

import (
	"path/filepath"
	"strings"

	"github.com/MUKE-coder/grit/v3/internal/manifest"
)

// Contact-app review L15, L16 and L18, in the files upgrade does not deliver
// whole. services/gdpr.go arrives whole with the erasure files; L17 is in
// service_writes_repair.go, and L19 is an admin page upgrade rewrites.

// ─── L16: the GDPR export streams its activity log ───────────────────────────

const (
	gdprExportCallOld = "\tbundle, err := services.ExportUserData(h.DB, targetID)\n"
	gdprExportCallNew = "\tbundle, err := services.ExportUserData(h.DB.WithContext(c.Request.Context()), targetID)\n"

	gdprExportWriteOld = "\tc.JSON(http.StatusOK, bundle)\n}\n"
	gdprExportWriteNew = `	c.Header("Content-Type", "application/json; charset=utf-8")
	c.Status(http.StatusOK)
	// The activity log is written a page at a time, so the status is on the wire
	// before the last page is read. A failure after that is logged, and leaves a
	// file that does not parse rather than one that looks complete.
	if err := bundle.WriteJSON(c.Writer); err != nil {
		log.Printf("gdpr export: writing the bundle: %v", err)
	}
}
`
)

// ─── L18: the notification bell's queries have indexes ───────────────────────

const (
	notificationUserIDTagOld = "`gorm:\"size:36;index\" json:\"user_id\"`"
	notificationUserIDTagNew = "`gorm:\"size:36;index:idx_notifications_user_created,priority:1;index:idx_notifications_user_read,priority:1\" json:\"user_id\"`"

	notificationReadAtTagOld = "`json:\"read_at\"`"
	notificationReadAtTagNew = "`gorm:\"index:idx_notifications_user_read,priority:2\" json:\"read_at\"`"

	notificationCreatedAtTagOld = "`gorm:\"index\" json:\"created_at\"`"
	notificationCreatedAtTagNew = "`gorm:\"index;index:idx_notifications_user_created,priority:2\" json:\"created_at\"`"

	notificationStructLine = "type Notification struct {\n"
	notificationIndexNote  = `//
// The bell polls two queries: a viewer's newest rows, and a count of their
// unread ones. Each has a composite index that starts at user_id, so neither
// reads every notification the viewer has ever had to answer. grit migrate
// builds them.
`
)

// repairPerfLows applies L15, L16 and L18 to an existing project.
func repairPerfLows(root string, opts Options) error {
	apiRoot := opts.APIRoot(root)
	m, err := manifest.Load(root)
	if err != nil {
		return err
	}
	for _, r := range []struct {
		path string
		fn   func(string) (string, []string, []string)
	}{
		{filepath.Join(apiRoot, "internal", "handlers", "jobs.go"), repairJobsInspectorSource},
		{filepath.Join(apiRoot, "internal", "handlers", "gdpr.go"), repairGDPRExportSource},
		{filepath.Join(apiRoot, "internal", "models", "notification.go"), repairNotificationIndexSource},
	} {
		if !fileExists(r.path) {
			continue
		}
		if err := repairSourceFile(root, m, r.path, r.fn); err != nil {
			return err
		}
	}
	return repairGeneratedServiceWrites(root, opts)
}

// repairJobsInspectorSource builds the jobs handler's inspector once.
func repairJobsInspectorSource(src string) (string, []string, []string) {
	if strings.Contains(src, "inspectorOnce.Do(") {
		return src, nil, nil
	}
	if strings.Count(src, jobsInspectorOld) != 1 || strings.Count(src, "\t\"strconv\"\n") != 1 {
		if strings.Contains(src, "asynq.NewInspector(") {
			return src, nil, []string{"not the file Grit wrote: build the asynq inspector once rather than per request, or every refresh of the jobs screen opens a Redis connection"}
		}
		return src, nil, nil
	}
	out := strings.Replace(src, jobsInspectorOld, jobsInspectorShared, 1)
	out = strings.ReplaceAll(out, "\tdefer inspector.Close()\n", "")
	out = strings.Replace(out, "\t\"strconv\"\n", "\t\"strconv\"\n\t\"sync\"\n", 1)
	return out, []string{"the jobs screen reuses one Redis inspector instead of opening a connection per request"}, nil
}

// repairGDPRExportSource binds the export to the request and streams it.
func repairGDPRExportSource(src string) (string, []string, []string) {
	if strings.Contains(src, "bundle.WriteJSON(c.Writer)") {
		return src, nil, nil
	}
	if strings.Count(src, gdprExportCallOld) != 1 || strings.Count(src, gdprExportWriteOld) != 1 ||
		strings.Count(src, "\t\"fmt\"\n") != 1 {
		if strings.Contains(src, "services.ExportUserData(") {
			return src, nil, []string{"not the file Grit wrote: write the export with bundle.WriteJSON(c.Writer) on the request's context, or it is built in memory whole"}
		}
		return src, nil, nil
	}
	out := strings.Replace(src, gdprExportCallOld, gdprExportCallNew, 1)
	out = strings.Replace(out, gdprExportWriteOld, gdprExportWriteNew, 1)
	if !strings.Contains(out, "\t\"log\"\n") {
		out = strings.Replace(out, "\t\"fmt\"\n", "\t\"fmt\"\n\t\"log\"\n", 1)
	}
	return out, []string{"the GDPR export follows the request's context and writes the activity log a page at a time"}, nil
}

// repairNotificationIndexSource adds the bell's composite indexes.
func repairNotificationIndexSource(src string) (string, []string, []string) {
	if strings.Contains(src, "idx_notifications_user_created") {
		return src, nil, nil
	}
	for _, anchor := range []string{notificationUserIDTagOld, notificationReadAtTagOld, notificationCreatedAtTagOld, notificationStructLine} {
		if strings.Count(src, anchor) != 1 {
			return src, nil, []string{"not the file Grit wrote: index notifications on (user_id, created_at) and (user_id, read_at), or the bell reads every row a viewer has"}
		}
	}
	out := strings.Replace(src, notificationUserIDTagOld, notificationUserIDTagNew, 1)
	out = strings.Replace(out, notificationReadAtTagOld, notificationReadAtTagNew, 1)
	out = strings.Replace(out, notificationCreatedAtTagOld, notificationCreatedAtTagNew, 1)
	out = strings.Replace(out, notificationStructLine, notificationIndexNote+notificationStructLine, 1)
	return out, []string{"notifications are indexed on (user_id, created_at) and (user_id, read_at): run grit migrate to build them"}, nil
}
