package scaffold

import (
	"go/format"
	"strings"
	"testing"
)

func formatPerfLowGo(t *testing.T, name, src string) string {
	t.Helper()
	out, err := format.Source([]byte(strings.ReplaceAll(src, "{{MODULE}}", "example.com/app")))
	if err != nil {
		t.Fatalf("%s is not valid Go: %v", name, err)
	}
	return string(out)
}

// checkPerfLowRepair reconstructs the file Grit wrote before from the current
// template, repairs it, and expects the template back, once.
func checkPerfLowRepair(t *testing.T, name, current, old string, repair func(string) (string, []string, []string)) {
	t.Helper()
	if old == current {
		t.Fatalf("%s: the reconstruction did not change the template", name)
	}
	out, fixed, warn := repair(formatPerfLowGo(t, name, old))
	if len(fixed) != 1 || len(warn) != 0 {
		t.Fatalf("%s was not repaired (%v %v)", name, fixed, warn)
	}
	if formatPerfLowGo(t, name, out) != formatPerfLowGo(t, name, current) {
		t.Errorf("repairing %s did not produce the template:\n%s", name, formatPerfLowGo(t, name, out))
	}
	again, fixed, warn := repair(formatPerfLowGo(t, name, out))
	if len(fixed) != 0 || len(warn) != 0 || formatPerfLowGo(t, name, again) != formatPerfLowGo(t, name, out) {
		t.Errorf("the %s repair is not idempotent", name)
	}
}

// L15: the jobs handler builds its asynq inspector once.
func TestJobsHandlerSharesOneInspector(t *testing.T) {
	current := jobsHandlerGo()
	if strings.Count(current, "asynq.NewInspector(") != 1 || strings.Contains(current, "inspector.Close()") ||
		!strings.Contains(current, "h.inspectorOnce.Do(") {
		t.Fatal("the jobs handler does not build one inspector for its lifetime")
	}
	old := strings.Replace(current, jobsInspectorShared, jobsInspectorOld, 1)
	old = strings.Replace(old, "\t\"sync\"\n", "", 1)
	const unavailable = "\"Job queue not available\")\n\t\treturn\n\t}\n"
	old = strings.ReplaceAll(old, unavailable, unavailable+"\tdefer inspector.Close()\n")
	if strings.Count(old, "defer inspector.Close()") != 4 {
		t.Fatalf("the reconstruction has %d deferred closes, want 4", strings.Count(old, "defer inspector.Close()"))
	}
	checkPerfLowRepair(t, "handlers/jobs.go", current, old, repairJobsInspectorSource)
}

// L16: the export is bound to the request, its errors are returned, and the
// activity log is written a page at a time.
func TestGDPRExportStreamsTheActivityLog(t *testing.T) {
	service := apiGDPRServiceGo()
	for _, want := range []string{
		"func (e *UserExport) WriteJSON(w io.Writer) error {",
		"const exportActivityPage = 500",
		`if err := db.Where("user_id = ?", userID).Find(&out.Uploads).Error; err != nil {`,
	} {
		if !strings.Contains(service, want) {
			t.Errorf("services/gdpr.go is missing %q", want)
		}
	}
	if strings.Contains(service, "Find(&out.Activity)") {
		t.Error("the export still reads the whole activity log into memory")
	}
	formatPerfLowGo(t, "services/gdpr.go", service)
	formatPerfLowGo(t, "services/gdpr_test.go", apiGDPRTestGo())

	current := apiGDPRHandlerGo()
	old := strings.Replace(current, gdprExportCallNew, gdprExportCallOld, 1)
	old = strings.Replace(old, gdprExportWriteNew, gdprExportWriteOld, 1)
	old = strings.Replace(old, "\t\"log\"\n", "", 1)
	checkPerfLowRepair(t, "handlers/gdpr.go", current, old, repairGDPRExportSource)
}

// L18: the bell's two queries each have a composite index.
func TestNotificationsIndexedForTheBell(t *testing.T) {
	current := notificationModelGo()
	for _, want := range []string{"idx_notifications_user_created,priority:1", "idx_notifications_user_read,priority:2", "idx_notifications_user_created,priority:2"} {
		if !strings.Contains(current, want) {
			t.Errorf("models/notification.go is missing %q", want)
		}
	}
	old := strings.Replace(current, notificationIndexNote, "", 1)
	old = strings.Replace(old, notificationUserIDTagNew, notificationUserIDTagOld, 1)
	old = strings.Replace(old, notificationReadAtTagNew, notificationReadAtTagOld, 1)
	old = strings.Replace(old, notificationCreatedAtTagNew, notificationCreatedAtTagOld, 1)
	checkPerfLowRepair(t, "models/notification.go", current, old, repairNotificationIndexSource)
}

// L19: the panel's root redirects on the server, in both shapes.
func TestAdminRootRedirectsOnTheServer(t *testing.T) {
	page := adminRedirectPage()
	if strings.Contains(page, "use client") || strings.Contains(page, "useMe") || !strings.Contains(page, `redirect("/dashboard");`) {
		t.Errorf("the admin root page is not a server redirect:\n%s", page)
	}
	if embedded := embedAdminContent(page, adminRoutePrefixes); !strings.Contains(embedded, `redirect("/admin/dashboard");`) {
		t.Errorf("inside the web app the root does not redirect to /admin/dashboard:\n%s", embedded)
	}
}
